// Odpowiedzialność pliku: wniesienie zasobu do Assets Panel
// (`design.asset.upload`) — magazyn bajtów zasobu, rozstrzygnięcie treści
// żądania i rozpoznanie formatu oraz wymiarów z nagłówka pliku. Metoda stoi na
// typie `*adapterDesignu` zadeklarowanym w `adapter_modul_design.go`; plik
// osobny wedle odpowiedzialności, jak `adapter_modul_library_tresc.go`
// w module Library.
//
// Drogi zasobu do modułu są dwie: wniesienie przez Operatora i generowanie
// kanałem. Obie odkładają bajty tym samym magazynem
// (`adapter_modul_design_generowanie.go`) — jeden magazyn, jedna miara formatu
// i wymiarów, żadnej drugiej prawdy.
//
// Treść zawsze ląduje w magazynie rdzenia; `sourcePath` jest źródłem bajtów,
// nie miejscem składowania. Odwołanie do cudzego pliku jest wskaźnikiem na treść
// żywą, a zasób Assets Panelu — jak wersja w bibliotece — ma być treścią
// zamrożoną: warstwa kompozycji ułożona na obrazie nie ma prawa zacząć leżeć na
// innym obrazie dlatego, że Operator posprzątał katalog pobrań. Ścieżka źródłowa
// nie jest tu nigdzie zapisywana: tabela `zasob_design` nie ma kolumny
// `sciezka`, a `uri` wskazuje blob.
//
// Magazyn jest ten sam co w bibliotece, tylko nad innym katalogiem. Typ
// `magazynTresciBiblioteki` (`adapter_modul_library_magazyn.go`) niesie
// wszystko, czego ten moduł potrzebuje — blob pod sumą sha256, zapis atomowy
// (temp w katalogu docelowym, `Sync`, `Rename`), dedup po nazwie — a drugi taki
// magazyn byłby drugą prawdą o tym, jak rdzeń trzyma bajty poza bazą.
// Tak samo pożycza go magazyn załączników rozmowy (`adapter_rozmowa_zalaczniki.go`).
// Cena jest jedna i widoczna: teksty odmów tego typu mówią „magazyn treści
// biblioteki", więc wołający ubiera je we własne zdanie modułu Design.
//
// Wymiary bierzemy z nagłówka pliku albo nie bierzemy wcale.
// `image.DecodeConfig` czyta sam nagłówek (nie rozpakowuje obrazu), więc PNG,
// JPEG i GIF oddają prawdziwe piksele za cenę kilkudziesięciu bajtów odczytu.
// Format spoza tej trójki (SVG, WEBP, plik wideo) zostawia `width` i `height`
// puste, bo pusto znaczy „nie wiem", a zgadnięta liczba wygląda w panelu
// identycznie jak zmierzona i nie da się jej odróżnić.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"image"
	// Rejestracja dekoderów nagłówka: bez tych trzech importów
	// `image.DecodeConfig` nie rozpozna żadnego formatu i każdy zasób zostałby
	// bez wymiarów. Import pusty, bo używamy wyłącznie skutku ubocznego
	// rejestracji — tak przewiduje pakiet `image`.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// przedrostekZasobuDesign znakuje identyfikatory zewnętrzne zasobów —
	// obok `plansza-` i `warstwa-` z `adapter_modul_design_kompozycje.go`.
	// Znakuje wiersze obu dróg — wniesionej i wygenerowanej — bo jedne i drugie
	// są tym samym bytem: zasobem z bajtami w magazynie.
	przedrostekZasobuDesign = "zasob-"

	// przedrostekPromptuDesign znakuje identyfikatory zewnętrzne promptów
	// wydanych — tych, które wchodzą do `design.prompt.history.list` i które
	// zasób wskazuje polem `promptId`.
	przedrostekPromptuDesign = "prompt-"

	// Magazyn zasobów leży w katalogu danych rdzenia, obok magazynu biblioteki
	// (`biblioteka/tresc`) i sejfu poświadczeń — przeżywa restart tak samo jak
	// wiersz w bazie, który go wskazuje. Osobne podkatalogi, bo moduły nie
	// dzielą stanu i skasowanie zasobów Designu nie ma prawa ruszyć
	// treści biblioteki.
	podkatalogDesignu        = "design"
	podkatalogZasobowDesignu = "zasoby"
)

// magazynZasobowDesignu składa magazyn bajtów zasobów nad katalogiem danych
// rdzenia. Katalog nie musi istnieć — powstaje przy pierwszym zapisie.
func magazynZasobowDesignu(katalogDanych string) *magazynTresciBiblioteki {
	if strings.TrimSpace(katalogDanych) == "" {
		return nil
	}
	return &magazynTresciBiblioteki{
		katalog: filepath.Join(katalogDanych, podkatalogDesignu, podkatalogZasobowDesignu),
	}
}

// korzenZasobowDesignu jest korzeniem magazynu zasobów liczonym od katalogu
// danych — ta sama para podkatalogów, którą składa `magazynZasobowDesignu`,
// więc postać odwołania i miejsce zapisu nie mają jak się rozjechać.
const korzenZasobowDesignu = podkatalogDesignu + "/" + podkatalogZasobowDesignu

// odwolanieZasobuDesignu oddaje `uri` zasobu w postaci, którą wolno wypuścić
// z rdzenia: ścieżkę względną magazynu zasobów, nie ścieżkę na dysku Operatora.
// Postać jest ta sama, co dla treści biblioteki — powód i cena stoją w nagłówku
// `odwolanieMagazynu` (`adapter_modul_library_magazyn.go`), bo to jedna decyzja
// dla całego produktu, a nie dwie zbieżne.
//
// Baza zostaje przy ścieżce bezwzględnej. Kolumna `uri` wiersza `zasob_design`
// jest bookkeepingiem rdzenia: czytają ją narzędzia modelu
// (`adapter_narzedzia_obraz.go`, `_media.go`, `_archiwum.go`, `_dokument.go`)
// i wysyłka poczty (`adapter_modul_poczta_wysylka.go`), otwierając plik wprost.
// Przełożenie jest więc granicą kontraktu, nie zmianą schematu — zmiana
// schematu ruszyłaby pięciu czytelników i wymagała migracji, a wyciek dotyczy
// wyłącznie tego, co opuszcza rdzeń.
func odwolanieZasobuDesignu(sciezka string) string {
	return odwolanieMagazynu(sciezka, korzenZasobowDesignu)
}

// WniesZasob utrwala bajty zasobu, zakłada jego wiersz i oddaje zasób odczytany
// z bazy wraz z etykietami — obsługuje `design.asset.upload`.
//
// Kolejność jest zamierzona: najpierw bajty, potem wiersz. Wiersz zasobu
// powstaje dopiero po utrwaleniu treści — inaczej Assets Panel pokazywałby
// kafelek, za którym nie ma nic, a brak ujawniłby się dopiero przy próbie
// obejrzenia zasobu. Ten sam porządek co `Wgraj` w module Library.
//
// Etykiety idą po zapisie wiersza, bo `etykieta_zasobu_design.zasob_id`
// wskazuje klucz wiersza, którego przed zapisem nie ma. Ich niepowodzenie jest
// odmową całej komendy: zasób wniesiony z etykietami, które nie weszły,
// wypadłby z własnego filtra w panelu i wyglądałby na zgubiony.
func (a *adapterDesignu) WniesZasob(ctx context.Context,
	z shared.DesignAssetUploadRequest) (shared.DesignAssetUploadResponse, error) {

	if z.WindowId == "" {
		return shared.DesignAssetUploadResponse{}, bladWskazaniaDesignu(
			"komenda design.asset.upload bez wskazania okna")
	}
	if strings.TrimSpace(string(z.Kind)) == "" {
		return shared.DesignAssetUploadResponse{}, bladWskazaniaDesignu(
			"komenda design.asset.upload bez rodzaju zasobu")
	}
	// Sprawdzenie idzie PRZED wniesieniem treści: bajty wchodzą do magazynu
	// wcześniej niż wiersz zasobu, więc odmowa po zapisie zostawiałaby blob,
	// którego nic już nie wskazuje, a magazyn Designu nie ma sprzątania.
	if err := sprawdzRodzajZasobu("design.asset.upload", z.Kind); err != nil {
		return shared.DesignAssetUploadResponse{}, err
	}

	odwolanie, err := a.trescZasobu(z.ContentBase64, z.SourcePath)
	if err != nil {
		return shared.DesignAssetUploadResponse{}, err
	}

	// Format i wymiary rozpoznajemy z tego, co naprawdę leży w magazynie, a nie
	// z żądania — jedna droga dla obu sposobów wniesienia i miara z bajtów
	// utrwalonych, nie z bajtów zapowiedzianych.
	format, szerokosc, wysokosc := rozpoznajObrazZasobu(odwolanie)
	// Wskazanie Operatora uzupełnia to, czego nagłówek nie powiedział, i nie
	// nadpisuje tego, co powiedział: zmierzone bije zadeklarowane, bo tylko
	// jedno z nich opisuje treść leżącą w magazynie.
	format = pierwszyTekst(format, z.Format)
	szerokosc = pierwszaLiczba(szerokosc, z.Width)
	wysokosc = pierwszaLiczba(wysokosc, z.Height)

	zapisany, err := a.repozytorium.ZapiszZasob(ctx, dane.ZasobDesignu{
		Kod:       nowyIdentyfikator(przedrostekZasobuDesign),
		Okno:      z.WindowId,
		Nazwa:     z.Name,
		Rodzaj:    string(z.Kind),
		Format:    format,
		URI:       &odwolanie,
		Szerokosc: szerokosc,
		Wysokosc:  wysokosc,
	})
	if err != nil {
		return shared.DesignAssetUploadResponse{}, bladDesignu(err)
	}

	if len(z.Tags) > 0 {
		if err := a.repozytorium.UstawEtykietyZasobu(ctx, zapisany.ID, z.Tags); err != nil {
			return shared.DesignAssetUploadResponse{}, bladDesignu(err)
		}
	}
	etykiety, err := a.repozytorium.EtykietyZasobu(ctx, zapisany.ID)
	if err != nil {
		return shared.DesignAssetUploadResponse{}, bladDesignu(err)
	}
	return shared.DesignAssetUploadResponse{Asset: zasobKontraktu(zapisany, etykiety)}, nil
}

// trescZasobu utrwala treść żądania w magazynie i oddaje odwołanie do bajtów —
// ścieżkę bloba, która staje się `uri` zasobu.
//
// Ścieżka ma pierwszeństwo przed treścią, gdy przyszły obie: bajty spod ścieżki
// są tym, co Operator naprawdę wskazał, a base64 bywa wtedy podglądem złożonym
// przez okno (ten sam rozstrzyg co w Library).
//
// Żądanie bez jednego i drugiego jest odmową — inaczej niż w Library, i to
// świadomie: tam wgranie bez treści zakłada wersję-znacznik, a tu zasób
// wizualny bez treści wizualnej nie ma czego pokazać (nagłówek
// `adapter_modul_design.go`).
func (a *adapterDesignu) trescZasobu(trescBase64, sciezkaZrodlowa *string) (string, error) {
	if a.magazyn == nil {
		return "", bladZapisuZasobuDesignu(
			"magazyn treści nie jest wpięty — nie ma gdzie odłożyć bajtów zasobu")
	}

	if !bezWartosci(sciezkaZrodlowa) {
		// Wciągamy bajty, nie dowiązujemy pliku — patrz nagłówek pliku.
		odwolanie, _, _, err := a.magazyn.ZapiszZePliku(*sciezkaZrodlowa)
		if err != nil {
			return "", bladZapisuZasobuDesignu("nie można wciągnąć treści spod ścieżki: " + err.Error())
		}
		return odwolanie, nil
	}

	if bezWartosci(trescBase64) {
		return "", bladWskazaniaDesignu(
			"komenda design.asset.upload bez treści zasobu: brakuje pola contentBase64 albo sourcePath — " +
				"wniesienie nie ma skąd wziąć bajtów; obraz wytworzony przez rdzeń zamawia się " +
				"komendą design.asset.generate, wskazując kanał obrazowy")
	}
	// Rozbiór base64 idzie tu, a nie pomocnikiem Library (`trescZadania`),
	// wyłącznie dla odmowy: tamten mówi „moduł Library" w zdaniu, które czytałby
	// Operator Assets Panelu, a moduł, którego nie wołał, nie ma prawa być
	// autorem jego odmowy. Rachunek jest ten sam — sha256 z bajtów, ta sama
	// suma, którą magazyn robi nazwą bloba.
	bajty, err := base64.StdEncoding.DecodeString(*trescBase64)
	if err != nil {
		return "", bladWskazaniaDesignu("treść zasobu nie jest poprawnym base64: " + err.Error())
	}
	if len(bajty) == 0 {
		return "", bladWskazaniaDesignu("treść zasobu jest pusta — zasób bez bajtów nie ma czego pokazać")
	}
	suma := sha256.Sum256(bajty)
	odwolanie, err := a.magazyn.Zapisz(bajty, hex.EncodeToString(suma[:]))
	if err != nil {
		return "", bladZapisuZasobuDesignu("nie można utrwalić treści zasobu: " + err.Error())
	}
	return odwolanie, nil
}

// rozpoznajObrazZasobu czyta nagłówek utrwalonego pliku i oddaje format wraz
// z wymiarami. Formatu nierozpoznanego nie zgaduje: trójka pustych wartości
// znaczy „plik nie jest obrazem, który rdzeń potrafi zmierzyć", a nie „obraz
// zerowy". Nieudany odczyt też kończy się pustkami — wniesienie już się
// powiodło, a brak wymiarów nie jest powodem, żeby odebrać Operatorowi zasób,
// którego bajty leżą utrwalone.
func rozpoznajObrazZasobu(sciezka string) (*string, *int, *int) {
	plik, err := os.Open(sciezka)
	if err != nil {
		return nil, nil, nil
	}
	defer plik.Close()

	opis, format, err := image.DecodeConfig(plik)
	if err != nil {
		return nil, nil, nil
	}
	szerokosc, wysokosc := opis.Width, opis.Height
	return &format, &szerokosc, &wysokosc
}

// pierwszyTekst oddaje wartość zmierzoną, a gdy jej nie ma — zadeklarowaną
// przez Operatora. Wartość pusta jest tu brakiem, nie tekstem: „format: " bez
// formatu byłby cechą, której nikt nie wskazał.
func pierwszyTekst(zmierzona, zadeklarowana *string) *string {
	if zmierzona != nil && strings.TrimSpace(*zmierzona) != "" {
		return zmierzona
	}
	if zadeklarowana != nil && strings.TrimSpace(*zadeklarowana) != "" {
		return zadeklarowana
	}
	return nil
}

// pierwszaLiczba działa jak `pierwszyTekst` dla wymiarów. Liczba niedodatnia
// jest odrzucana po obu stronach — obraz o szerokości zero nie istnieje, a
// zapisany w wierszu wyglądałby na zmierzony.
func pierwszaLiczba(zmierzona, zadeklarowana *int) *int {
	if zmierzona != nil && *zmierzona > 0 {
		return zmierzona
	}
	if zadeklarowana != nil && *zadeklarowana > 0 {
		return zadeklarowana
	}
	return nil
}

// bladZapisuZasobuDesignu nazywa niepowodzenie utrwalenia bajtów: nośnik pełny,
// brak praw do katalogu danych, ścieżka źródłowa nieczytelna, magazyn niewpięty.
// To awaria zapisu po stronie rdzenia, nie wina Operatora, więc kod jest
// wewnętrzny; komenda odmawia w całości, bo zasób, którego bajtów nie ma nigdzie,
// nie ma prawa trafić do panelu jako wniesiony.
func bladZapisuZasobuDesignu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError, "moduł Design: "+powod))
}
