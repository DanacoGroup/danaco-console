// Dogniatanie zapisu obrazu po rachunku wkompilowanym: `image.convert` oddaje
// plik mniejszy o tyle, ile potrafią zdjąć optipng, jpegoptim, pngquant
// i cwebp. Ten plik nigdy nie odmawia: brak programu oddaje bajty wejściowe.
package core

import (
	"context"
	"os"
	"path/filepath"

	"danacoconsole/server/internal/zewnetrzne"
)

// narzedzieOptipng dogniata PNG bez zmiany pikseli, szukając filtrów i dłuższego
// przebiegu deflate niż koder Go wkompilowany w rachunek podstawowy.
func narzedzieOptipng() zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{Nazwa: "OptiPNG", Program: "optipng", Pakiet: "optipng"}
}

// narzedzieJpegoptim przelicza tablice Huffmana JPEG-a bez zmiany pikseli,
// tam gdzie koder Go wkompilowany zapisuje tablice domyślne.
func narzedzieJpegoptim() zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{Nazwa: "jpegoptim", Program: "jpegoptim", Pakiet: "jpegoptim"}
}

// narzedziePngquant sprowadza PNG do palety: zapis jest stratny, wołany
// wyłącznie na wprost żądaną prośbę o stratę jakości obrazu.
func narzedziePngquant() zewnetrzne.Narzedzie {
	return zewnetrzne.Narzedzie{Nazwa: "pngquant", Program: "pngquant", Pakiet: "pngquant"}
}

// narzedzieCwebp zapisuje WebP z pełnym strojeniem predyktorów, w trybie
// bezstratnym, bez ruszania pikseli obrazu źródłowego.
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

// dogniecZapisObrazu oddaje bajty krótsze od podanych albo podane bez zmiany,
// nie zwracając błędu żadną drogą: zapis obrazu już się udał wcześniej.
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
			// Paleta idzie po optipng: pngquant oddaje plik, który optipng jeszcze skraca.

			// Odwrotna kolejność marnuje pierwszy przebieg.
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
				// `--dest` żąda katalogu, nie pliku, i zachowuje nazwę źródła.
				return []string{"-q", "--strip-none", "--dest", filepath.Dir(wyjscie), wejscie}
			})
	case "webp":
		if stratnieWolno {
			// WEBP stratny powstaje już programem, bo kodera stratnego w Go nie ma.

			// Ponowne przepuszczenie przez koder stratny byłoby drugą stratą pikseli.
			return bajty
		}
		return a.przezPlik(ctx, bajty, "webp", narzedzieCwebp(),
			func(wejscie, wyjscie string) []string {
				return []string{"-quiet", "-lossless", "-z", "9", wejscie, "-o", wyjscie}
			})
	}
	// Formaty bez programu dogniatającego (avif, tiff, gif) wychodzą jak przyszły.
	return bajty
}

// przezPlik przeprowadza jedno dogniecenie: kładzie bajty w katalogu przebiegu,
// woła program i oddaje wynik, jeśli ten naprawdę jest krótszy.
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

	// Wejście i wyjście mają tę samą nazwę w dwóch katalogach.

	// Jpegoptim wskazuje wynik katalogiem, a nie nazwą pliku, inne programy nie.
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
		// Kod niezerowy znaczy tu najczęściej „nie umiem tego skrócić".

		// Zapis pierwotny jest wtedy właściwą odpowiedzią.
		return bajty
	}

	krotsze, err := os.ReadFile(wyjscie)
	if err != nil || len(krotsze) == 0 || len(krotsze) >= len(bajty) {
		// Program, który zakończył się powodzeniem, a nie zostawił krótszego pliku.

		// Nie miał czego skrócić — oddanie wyniku powiększyłoby zasób.
		return bajty
	}
	return krotsze
}
