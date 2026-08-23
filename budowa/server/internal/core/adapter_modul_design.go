// Odpowiedzialność pliku: moduł Design — typ adaptera, konstruktor, złożenie
// promptu strukturalnego i wykaz zasobów (`design.asset.list`). Wytworzenie
// zasobu (`design.asset.generate`) leży w `adapter_modul_design_generowanie.go`,
// wybór kanału obrazowego w `_kanal.go`, wniesienie w `_wgranie.go`, kompozycje
// i etykiety w swoich plikach — plik wedle odpowiedzialności.
//
// Generowanie oddaje bajty obrazu, nie samą kopertę: komenda albo zapisuje
// zasób o prawdziwej treści, albo odmawia zdaniem nazywającym brak rzeczywisty,
// składany przy każdym wywołaniu osobno. Obrazu zastępczego, zasobu bez bajtów
// ani powodzenia bez treści tu nie ma — każda droga bez bajtów kończy się
// błędem.
package core

import (
	"context"
	"fmt"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// adapterDesignu wypełnia port Design. Trzy zależności, bo droga generowania
// potrzebuje wszystkich trzech: wiersza w bazie, kanału do silnika i miejsca
// na bajty.
type adapterDesignu struct {
	repozytorium dane.RepozytoriumDesignu
	// kanaly jest rejestrem kanałów modelu rdzenia, wpiętym przez montaż
	// (`ZKanalami`). Generowanie sięga wyłącznie po wiersze z adapterem
	// „obrazy", bo tylko one oddają fragment `image`; wybór i jego odmowy stoją
	// w `adapter_modul_design_kanal.go`.
	kanaly *models.Rejestr
	// magazyn jest miejscem na bajty zasobów wniesionych komendą
	// `design.asset.upload`. Baza trzyma wyłącznie odwołanie, więc treść musi
	// mieć gdzie leżeć — szczegóły w `adapter_modul_design_wgranie.go`. Od tego
	// pola zależą obie drogi zasobu do modułu: wniesienie i generowanie.
	magazyn *magazynTresciBiblioteki
	// obecnosc trzyma kursory współpracy na kompozycjach. Rejestr żyje
	// w pamięci i ginie razem z procesem — powód w nagłówku
	// `adapter_modul_design_obecnosc.go`.
	obecnosc *rejestrObecnosciDesignu
	// sejf jest magazynem sekretów rdzenia. Moduł czyta z niego WYŁĄCZNIE klucze
	// baz zdjęciowych (`design.stock.*`) i nigdy nic tam nie zapisuje: klucz
	// dostawcy wchodzi do sejfu drogą kont i punktów dostępu, a moduł Design ma
	// go tylko odczytać. Brak sejfu nie odbiera modułowi funkcji — dostawcy bez
	// klucza pracują dalej, a dostawcy z kluczem wracają w `providersFailed`
	// (`adapter_modul_design_bazy_zdjeciowe.go`).
	sejf SejfPoswiadczen
	// uruchamiacz, rozstrzygacz i katalogRoboczy służą JEDNEJ czynności:
	// odczytowi treści napisów ze zrzutu ekranu (`design.mockup.import`).
	// Czytnika liter w czystym Go nie ma, więc rozpoznanie pisma idzie programem
	// pakietu serwera — drogą pakietu `zewnetrzne`, jedyną w drzewie, która
	// sprawdza obecność programu i nakłada bramę izolacji. Uruchomienie stoi
	// w `adapter_modul_design_makiety_zrzut.go` i to JEDEN plik obszaru Design
	// objęty nazwanym wyjątkiem zapory (`plikiDesignuZOdczytemPisma`); nazwy
	// wywołania nie ma prawa być nawet w tym komentarzu, bo zapora czyta treść
	// pliku, nie składnię. Brak tych zależności nie odbiera modułowi układu
	// makiety: rozpoznanie odmawia wtedy zdaniem nazywającym brak, a obszary
	// wychodzą jak dotąd.
	uruchamiacz    session.Uruchamiacz
	rozstrzygacz   *konfig.Rozstrzygacz
	katalogRoboczy *KatalogRoboczy
}

// ZOdczytemPisma wpina uruchamiacz procesów wraz z bramą izolacji — drogę,
// którą `design.mockup.import` woła program rozpoznający pismo.
//
// Zależność jest OSOBNA od pozostałych i nazwana od czynności, nie od warstwy:
// moduł Design nie startuje procesów do niczego innego i ma się nie dać
// rozszerzyć „bo uruchamiacz już jest".
func (a *adapterDesignu) ZOdczytemPisma(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterDesignu {

	a.uruchamiacz, a.rozstrzygacz, a.katalogRoboczy = uruchamiacz, rozstrzygacz, katalog
	return a
}

// ZSejfem wpina magazyn sekretów rdzenia. Tędy `design.stock.search`
// i `design.stock.import` czytają klucze dostawców darmowych — bez nowej
// zależności modułu i bez drugiego sejfu nad tym samym plikiem.
func (a *adapterDesignu) ZSejfem(sejf SejfPoswiadczen) *adapterDesignu {
	a.sejf = sejf
	return a
}

// nowyAdapterDesignu wiąże (część) portu z repozytorium modułu i wpina magazyn
// treści zasobów oparty o katalog danych rdzenia — wzorem `nowyAdapterBiblioteki`.
//
// Katalog danych jest tu domyślny, a obowiązujący wchodzi montażem przez
// `ZKatalogiemDanych` — tą samą zmienną, którą jadą sejf poświadczeń i magazyn
// biblioteki, żeby o katalogu była jedna prawda. Wartość domyślna zostaje dla
// wywołania bez montażu, aby konstruktor nigdy nie oddał adaptera bez magazynu.
func nowyAdapterDesignu(repozytorium dane.RepozytoriumDesignu) *adapterDesignu {
	return &adapterDesignu{
		repozytorium: repozytorium,
		magazyn:      magazynZasobowDesignu(konfiguracja.KatalogDanychDomyslny()),
		obecnosc:     nowyRejestrObecnosciDesignu(),
	}
}

// ZKatalogiemDanych przestawia magazyn treści zasobów na katalog wskazany
// konfiguracją procesu (przełącznik `-dane`, zmienna `DANACO_KATALOG_DANYCH`).
// Bez tego wywołania zasoby lądowałyby w katalogu domyślnym, a baza tam, gdzie
// wskazuje konfiguracja; wiersz wskazywałby wtedy blob spoza swojego katalogu
// danych i przeniesienie rdzenia gubiłoby treść zasobów.
func (a *adapterDesignu) ZKatalogiemDanych(katalog string) *adapterDesignu {
	if katalog != "" {
		a.magazyn = magazynZasobowDesignu(katalog)
	}
	return a
}

// ZKanalami przyjmuje rejestr kanałów modelu od montażu i zwraca adapter, żeby
// montaż wiązał zależność w łańcuchu. Tędy generowanie sięga po kanał obrazowy —
// bez nowej zależności i bez drugiego wpięcia.
func (a *adapterDesignu) ZKanalami(kanaly *models.Rejestr) *adapterDesignu {
	a.kanaly = kanaly
	return a
}

// zlozPolecenieObrazu składa pola promptu strukturalnego w jedną gotową treść
// polecenia — taką, jaką podaje się silnikowi generującemu obrazy.
//
// Ta sama treść jedzie dwiema drogami: do kanału obrazowego jako polecenie
// generowania, a przy każdej odmowie w `error.details.polecenie`. Dzięki drugiej
// pracy Prompt Buildera odmowa nie kasuje.
//
// Pola puste nie wchodzą: kontrakt opisuje brak wartości brakiem pola, więc
// „styl: " bez stylu byłby zmyśleniem cechy, której nikt nie wskazał.
func zlozPolecenieObrazu(p shared.DesignPrompt) string {
	var b strings.Builder
	b.WriteString(p.Subject)
	dopiszCeche(&b, "styl", p.Style)
	dopiszCeche(&b, "kompozycja", p.Composition)
	dopiszCeche(&b, "oświetlenie", p.Lighting)
	dopiszCeche(&b, "paleta", p.Palette)
	dopiszCeche(&b, "proporcje kadru", p.AspectRatio)
	dopiszCeche(&b, "bez", p.Exclusions)
	return b.String()
}

// dopiszCeche dokłada jedną cechę promptu do polecenia, gdy została wskazana.
func dopiszCeche(b *strings.Builder, nazwa string, wartosc *string) {
	if wartosc == nil || strings.TrimSpace(*wartosc) == "" {
		return
	}
	b.WriteString("; ")
	b.WriteString(nazwa)
	b.WriteString(": ")
	b.WriteString(strings.TrimSpace(*wartosc))
}

// Zasoby zwraca stronę zasobów Assets Panel spełniających filtr wraz z liczbą
// wszystkich spełniających: `Total` niesie liczbę spełniających warunki, nie
// długość zwróconej strony — repozytorium liczy je osobno.
//
// Pole `PromptId` zostaje puste zawsze — patrz `zasobKontraktu`.
func (a *adapterDesignu) Zasoby(ctx context.Context,
	z shared.DesignAssetListRequest) (shared.DesignAssetListResponse, error) {

	filtr := dane.FiltrZasobow{
		Okno:          z.WindowId,
		Etykiety:      z.Tags,
		TylkoUlubione: z.FavoriteOnly != nil && *z.FavoriteOnly,
		Limit:         wartoscLiczby(z.Limit),
	}
	if z.Kind != nil {
		rodzaj := string(*z.Kind)
		filtr.Rodzaj = &rodzaj
	}

	wiersze, razem, err := a.repozytorium.Zasoby(ctx, filtr)
	if err != nil {
		return shared.DesignAssetListResponse{}, bladDesignu(err)
	}
	zasoby := make([]shared.DesignAsset, 0, len(wiersze))
	for _, wiersz := range wiersze {
		etykiety, err := a.repozytorium.EtykietyZasobu(ctx, wiersz.ID)
		if err != nil {
			return shared.DesignAssetListResponse{}, bladDesignu(err)
		}
		zasoby = append(zasoby, zasobKontraktu(wiersz, etykiety))
	}
	return shared.DesignAssetListResponse{Assets: zasoby, Total: &razem}, nil
}

// zasobKontraktu składa DesignAsset kontraktu z wiersza repozytorium.
//
// `PromptId` niesie identyfikator ZEWNĘTRZNY promptu. Repozytorium czyta go
// podzapytaniem obok wiersza zasobu (`dane/design_zasoby.go`), więc pole jest
// prawdziwe albo puste — nigdy odgadnięte. Puste znaczy „ten zasób nie powstał
// z promptu": zasób wniesiony przez Operatora promptu nie ma i mieć nie może.
// Drugą stronę tego powiązania oddaje `design.prompt.history.list`.
//
// `Uri` wychodzi jako ścieżka względna magazynu, nie jako ścieżka na dysku.
// Wiersz trzyma ścieżkę bezwzględną bloba, bo tą ścieżką rdzeń otwiera plik, ale
// wypuszczenie jej odsłania układ katalogów maszyny. Przełożenie stoi tutaj,
// w jedynym miejscu składania `DesignAsset` z wiersza, więc obejmuje wszystkich
// wołających naraz i nie da się go ominąć nową drogą. Powód wyboru postaci:
// `odwolanieMagazynu` w `adapter_modul_library_magazyn.go`.
//
// Puste odwołanie znaczy „nie mam czego podać": wiersz, którego `uri` nie leży
// w magazynie zasobów, nie dostaje pola zastępczego — kontrakt czyni je
// niewymaganym właśnie po to.
func zasobKontraktu(z dane.ZasobDesignu, etykiety []string) shared.DesignAsset {
	var odwolanie *string
	if z.URI != nil {
		if wzgledne := odwolanieZasobuDesignu(*z.URI); wzgledne != "" {
			odwolanie = &wzgledne
		}
	}
	return shared.DesignAsset{
		Id: z.Kod, WindowId: z.Okno, Name: z.Nazwa, Kind: shared.DesignAssetKind(z.Rodzaj),
		Format: z.Format, Uri: odwolanie, PromptId: z.PromptKod,
		VariantOfAssetId: z.WariantZasobuID,
		Tags:             etykiety, Favorite: &z.Ulubiony, Width: z.Szerokosc, Height: z.Wysokosc,
		CreatedAt: chwilaBazy(z.Utworzono),
	}
}

// sprawdzWyliczenieDesignu odrzuca wartość spoza wyliczenia kontraktu PRZED
// dotknięciem bazy — ten sam rozstrzyg i ten sam powód, co
// `sprawdzRodzajZasobu`, tylko dla dowolnego wyliczenia obszaru.
//
// Wykaz dopuszczalnych wartości pochodzi z kontraktu (`shared.WartosciDesign*`)
// i NIE jest tutaj przepisywany: wartość dołożona do kontraktu wchodzi do
// sprawdzenia sama, bez zmiany tego pliku. Rodziny komend obszaru Design mają
// kilkanaście wyliczeń (rodzaj węzła, operacja logiczna, kotwica więzu,
// wyzwalacz przejścia, harmonia barw, norma druku…), a osobne sprawdzenie na
// każde z nich byłoby tym samym zdaniem napisanym kilkanaście razy.
func sprawdzWyliczenieDesignu[T ~string](komenda, pole string, wartosc T, dopuszczalne []T) error {
	for _, znana := range dopuszczalne {
		if wartosc == znana {
			return nil
		}
	}
	nazwy := make([]string, 0, len(dopuszczalne))
	for _, znana := range dopuszczalne {
		nazwy = append(nazwy, string(znana))
	}
	return bladWskazaniaDesignu(fmt.Sprintf(
		"komenda %s z wartością %q pola %s, której kontrakt nie zna; wartości dopuszczalne: %s",
		komenda, string(wartosc), pole, strings.Join(nazwy, ", ")))
}

// bladDesignu znakuje usterkę kodem kontraktu, żeby okno modułu pokazało
// powód, nie samo „nie udało się”.
func bladDesignu(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaDesignu nazywa brak danych w żądaniu — usterka wołającego, nie
// rdzenia.
func bladWskazaniaDesignu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		fmt.Sprintf("moduł Design: %s", powod)))
}

// bladNieznanegoBytuDesignu nazywa byt modułu, którego rdzeń nie zna: kolekcję,
// kompozycję, wersję, adnotację, zestaw żetonów. Kod `not_found` odróżnia
// wskazanie nieaktualne (okno ma odświeżyć swój wykaz) od usterki rdzenia —
// tak samo, jak `bladNieznanegoZasobuDesignu` robi to dla zasobu.
func bladNieznanegoBytuDesignu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound, "moduł Design: "+powod))
}

// sprawdzRodzajZasobu odrzuca rodzaj spoza kontraktu PRZED dotknięciem magazynu.
//
// Bez tego sprawdzenia wartość ósma szła wprost do kolumny, odbijała się od
// warunku CHECK schematu i wracała jako `internal_error` z `retryable: true` —
// czyli jako awaria rdzenia, którą wolno ponawiać. Żądanie z rodzajem spoza
// wykazu nie uda się przy ŻADNYM ponowieniu, więc klient z pętlą ponowień
// powtarzał je bez końca. To jest pomyłka wołającego i ma wracać jako pomyłka
// wołającego — rdzeń zna to rozróżnienie i stosuje je w module Library.
//
// Odmowa wymienia dopuszczalne rodzaje, zamiast cytować warunek schematu wraz
// z nazwą kolumny: Operator ma dostać zdanie o swoim żądaniu, nie o bazie.
//
// Wykaz pochodzi z kontraktu i nie jest tu przepisywany. Wartość dołożona
// do kontraktu wchodzi więc do sprawdzenia sama, bez zmiany tego pliku.
func sprawdzRodzajZasobu(komenda string, rodzaj shared.DesignAssetKind) error {
	dopuszczalne := shared.WartosciDesignAssetKind()
	for _, znany := range dopuszczalne {
		if rodzaj == znany {
			return nil
		}
	}
	nazwy := make([]string, 0, len(dopuszczalne))
	for _, znany := range dopuszczalne {
		nazwy = append(nazwy, string(znany))
	}
	return bladWskazaniaDesignu(fmt.Sprintf(
		"komenda %s z rodzajem zasobu %q, którego kontrakt nie zna; rodzaje dopuszczalne: %s",
		komenda, string(rodzaj), strings.Join(nazwy, ", ")))
}
