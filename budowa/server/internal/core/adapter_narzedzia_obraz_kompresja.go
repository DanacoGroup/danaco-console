// Dogniatanie zapisu obrazu po rachunku wkompilowanym: `image.convert` oddaje
// plik mniejszy o tyle, ile potrafią zdjąć optipng, jpegoptim, pngquant i cwebp.
//
// ── Dlaczego to jest osobny krok, a nie inny koder ───────────────────────────
// Kodery Go zapisują obraz poprawnie, ale nie szukają najlepszego zapisu:
// `image/png` bierze jeden filtr na wiersz i jeden przebieg deflate, `image/jpeg`
// zapisuje domyślne tablice Huffmana, `nativewebp` nie stroi predyktorów.
// Wymienione programy robią dokładnie jedną rzecz — przeliczają TEN SAM obraz
// na krótszy strumień bajtów — i robią to lepiej, bo na to je napisano.
//
// ── Ulepszenie, nie warunek ─────────────────────────────────────────────────
// Rodzina `image.*` liczy się biblioteką wkompilowaną i ma się liczyć dalej na
// maszynie, na której żaden z tych czterech programów nie stoi. Dlatego ten plik
// NIE ODMAWIA nigdy: brak programu, niezerowy kod wyjścia, wynik pusty albo
// wynik większy od źródła — każdy z tych przypadków oddaje bajty wejściowe bez
// zmiany. Obraz ma się zapisać także wtedy, gdy nie ma czym go dogniatać, a
// odmowa w tym miejscu zamieniłaby ulepszenie w wymóg wobec wdrożenia.
//
// Z tego samego powodu nie ma tu odmowy nazywającej brak, jaką niesie
// `zewnetrzne.BrakNarzedzia`: brak programu dogniatającego nie jest czymś, o
// czym Operator ma się dowiedzieć w chwili zapisu obrazu — dowiaduje się przy
// starcie, z sondy wykazu zależności.
//
// ── Bezstratnie znaczy bezstratnie ──────────────────────────────────────────
// Trzy z czterech programów przeliczają zapis bez ruszania pikseli: optipng
// szuka filtrów i dłuższego deflate, jpegoptim przelicza tablice Huffmana,
// cwebp w trybie `-lossless` stroi predyktory WebP. `pngquant` jest inny —
// sprowadza obraz do palety, więc PIKSELE ZMIENIA. Wchodzi wyłącznie wtedy, gdy
// żądanie wprost prosi o zapis stratny (`lossless: false`); przy braku pola
// i przy `lossless: true` nie jest w ogóle wołany. Pomylenie tych dwóch rzeczy
// oddałoby model prosząc o zapis bezstratny obraz o zmienionych barwach.
//
// Z tej samej strony patrzy druga reguła: dogniatanie nie dokłada POKOLENIA
// kompresji stratnej. WEBP stratny powstaje już programem, bo kodera stratnego
// w Go nie ma, więc przy `lossless: false` ten plik zostawia go nietkniętym —
// drugi przebieg kodera stratnego odjąłby jakość, nie bajty.
package core

import (
	"context"
	"os"
	"path/filepath"

	"danacoconsole/server/internal/zewnetrzne"
)

// narzedzieOptipng dogniata PNG bez zmiany pikseli.
func narzedzieOptipng() zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{Nazwa: "OptiPNG", Program: "optipng", Pakiet: "optipng"}
}

// narzedzieJpegoptim przelicza tablice Huffmana JPEG-a bez zmiany pikseli.
func narzedzieJpegoptim() zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{Nazwa: "jpegoptim", Program: "jpegoptim", Pakiet: "jpegoptim"}
}

// narzedziePngquant sprowadza PNG do palety. Zapis STRATNY — patrz nagłówek.
func narzedziePngquant() zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{Nazwa: "pngquant", Program: "pngquant", Pakiet: "pngquant"}
}

// narzedzieCwebp zapisuje WebP z pełnym strojeniem predyktorów.
func narzedzieCwebp() zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{Nazwa: "cwebp", Program: "cwebp", Pakiet: "webp"}
}

// stopienOptipng jest poziomem poszukiwań OptiPNG. `-o2` schodzi z rozmiaru
// zauważalnie, a mieści się w ułamku sekundy; `-o7` bywa o promile lepsze
// i trwa dziesiątki sekund, czego zapis pojedynczego obrazu nie zniesie.
const stopienOptipng = "-o2"

// jakoscPngquant jest zakresem, poniżej którego pngquant woli nie oddać nic.
// Dolna granica jest progiem odrzucenia, górna — celem; program kończy się
// wtedy kodem 99, który ten plik czyta jako „nie da się bez utraty jakości"
// i zostawia zapis bezstratny.
const jakoscPngquant = "65-95"

// dogniecZapisObrazu oddaje bajty krótsze od podanych albo podane bez zmiany.
//
// Nie zwraca błędu żadną drogą — to jest istota tego kroku (patrz nagłówek).
// Wołający nie ma tu czego obsłużyć: zapis obrazu już się udał, a ten krok może
// go wyłącznie skrócić.
func (a *adapterNarzedziObrazu) dogniecZapisObrazu(ctx context.Context, bajty []byte,
	format string, bezstratnie *bool) []byte {

	if len(bajty) == 0 || a.uruchamiacz == nil {
		return bajty
	}
	stratnieWolno := bezstratnie != nil && !*bezstratnie

	switch format {
	case "png":
		krotsze := a.przezPlik(ctx, bajty, "png", narzedzieOptipng(),
			func(wejscie, wyjscie string) []string {
				return []string{stopienOptipng, "-quiet", "-out", wyjscie, wejscie}
			})
		if stratnieWolno {
			// Paleta idzie PO optipng: pngquant oddaje plik palety, który
			// optipng jeszcze skraca, a odwrotna kolejność marnuje pierwszy
			// przebieg.
			krotsze = a.przezPlik(ctx, krotsze, "png", narzedziePngquant(),
				func(wejscie, wyjscie string) []string {
					return []string{"--quality=" + jakoscPngquant, "--speed", "3",
						"--force", "--output", wyjscie, wejscie}
				})
		}
		return krotsze
	case "jpeg", "jpg":
		return a.przezPlik(ctx, bajty, "jpg", narzedzieJpegoptim(),
			func(wejscie, wyjscie string) []string {
				// `--dest` żąda katalogu, nie pliku, i zachowuje nazwę źródła —
				// dlatego nazwa wejścia i wyjścia jest ta sama, a różni je
				// katalog.
				return []string{"-q", "--strip-none", "--dest", filepath.Dir(wyjscie), wejscie}
			})
	case "webp":
		if stratnieWolno {
			// WEBP stratny powstaje już programem (`policzProgramem`), bo kodera
			// stratnego w Go nie ma. Ponowne przepuszczenie gotowego pliku przez
			// koder stratny byłoby DRUGĄ stratą na tych samych pikselach —
			// dogniecenie ma skracać zapis, a nie dokładać pokolenie kompresji.
			return bajty
		}
		return a.przezPlik(ctx, bajty, "webp", narzedzieCwebp(),
			func(wejscie, wyjscie string) []string {
				return []string{"-quiet", "-lossless", "-z", "9", wejscie, "-o", wyjscie}
			})
	}
	// Formaty bez programu dogniatającego (avif, tiff, gif) wychodzą takie,
	// jakie przyszły. Milczenie jest tu właściwe: nie ma czego zgłaszać.
	return bajty
}

// przezPlik przeprowadza jedno dogniecenie: kładzie bajty w katalogu przebiegu,
// woła program i oddaje wynik, jeśli ten naprawdę jest krótszy.
//
// Pliki pośrednie są konieczne z tego samego powodu, co przy silnikach
// neuronowych (`pracowniaObrazu`): optipng, pngquant i cwebp żądają ścieżki
// wyniku i nie umieją pisać na standardowe wyjście. jpegoptim by umiał, ale idzie
// tą samą drogą — jedna droga zamiast dwóch jest tu warta jednego zapisu na dysk.
func (a *adapterNarzedziObrazu) przezPlik(ctx context.Context, bajty []byte,
	rozszerzenie string, narzedzie zewnetrzne.Narzedzie,
	argumenty func(wejscie, wyjscie string) []string) []byte {

	if !zewnetrzne.Stoi(narzedzie) {
		return bajty
	}
	katalog, err := os.MkdirTemp("", "danaco-obraz-dognie-")
	if err != nil {
		return bajty
	}
	defer func() { _ = os.RemoveAll(katalog) }()

	// Wejście i wyjście mają tę samą nazwę w dwóch katalogach, bo jpegoptim
	// wskazuje wynik katalogiem, a nie nazwą pliku; pozostałym trzem programom
	// jest to obojętne.
	podkatalog := filepath.Join(katalog, "wynik")
	if err := os.Mkdir(podkatalog, 0o700); err != nil {
		return bajty
	}
	wejscie := filepath.Join(katalog, "obraz."+rozszerzenie)
	wyjscie := filepath.Join(podkatalog, "obraz."+rozszerzenie)
	if err := os.WriteFile(wejscie, bajty, 0o600); err != nil {
		return bajty
	}

	okno, zasady, obszar := a.zasiegNarzedzi()
	if _, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedzie, argumenty(wejscie, wyjscie), katalog, granicaNarzedziObrazu); err != nil {
		// Kod niezerowy znaczy tu najczęściej „nie umiem tego skrócić"
		// (pngquant kończy tak zapis, którego nie da się sprowadzić do palety
		// w zadanej jakości). Zapis pierwotny jest wtedy właściwą odpowiedzią.
		return bajty
	}

	krotsze, err := os.ReadFile(wyjscie)
	if err != nil || len(krotsze) == 0 || len(krotsze) >= len(bajty) {
		// Program, który zakończył się powodzeniem i nie zostawił pliku
		// krótszego, nie miał czego skrócić. Oddanie jego wyniku mimo to
		// powiększyłoby zasób w imię jego zmniejszenia.
		return bajty
	}
	return krotsze
}
