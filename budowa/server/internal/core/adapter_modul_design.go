// Moduł Design składa typ adaptera, konstruktor, budowę promptu strukturalnego
// i wykaz zasobów design.asset.list; wytworzenie, kanał obrazowy, wniesienie,
// kompozycje i etykiety leżą w osobnych plikach modułu wedle odpowiedzialności.
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
	// kanaly jest rejestrem kanałów modelu wpiętym montażem; generowanie
	// czyta stąd adapter obrazów.
	kanaly *models.Rejestr
	// magazyn jest miejscem na bajty zasobów wniesionych upload; baza trzyma
	// wyłącznie odwołanie do nich.
	magazyn *magazynTresciBiblioteki
	// obecnosc trzyma kursory współpracy na kompozycjach; rejestr żyje
	// w pamięci i ginie razem z procesem.
	obecnosc *rejestrObecnosciDesignu
	// sejf jest magazynem sekretów rdzenia; moduł czyta stąd klucze baz
	// zdjęciowych i nic nie zapisuje.
	sejf SejfPoswiadczen
	// uruchamiacz, rozstrzygacz i katalogRoboczy służą odczytowi napisów ze
	// zrzutu ekranu komendą import.
	uruchamiacz    session.Uruchamiacz
	rozstrzygacz   *konfig.Rozstrzygacz
	katalogRoboczy *KatalogRoboczy
}

// ZOdczytemPisma wpina uruchamiacz procesów wraz z bramą izolacji — drogę,
// którą design.mockup.import woła program rozpoznający pismo. Zależność jest
// osobna od pozostałych: moduł Design nie startuje procesów do niczego innego.
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

// nowyAdapterDesignu wiąże port z repozytorium modułu i wpina magazyn treści
// zasobów oparty o domyślny katalog danych rdzenia. Obowiązujący katalog
// wchodzi montażem przez ZKatalogiemDanych; wartość domyślna zostaje dla
// wywołania bez montażu.
func nowyAdapterDesignu(repozytorium dane.RepozytoriumDesignu) *adapterDesignu {
	return &adapterDesignu{
		repozytorium: repozytorium,
		magazyn:      magazynZasobowDesignu(konfiguracja.KatalogDanychDomyslny()),
		obecnosc:     nowyRejestrObecnosciDesignu(),
	}
}

// ZKatalogiemDanych przestawia magazyn treści zasobów na katalog wskazany
// konfiguracją procesu. Bez tego wywołania zasoby lądowałyby w katalogu
// domyślnym, a baza tam, gdzie wskazuje konfiguracja, co gubi treść zasobów
// przy przeniesieniu.
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
// polecenia, jaką podaje się silnikowi generującemu obrazy. Ta sama treść
// jedzie do kanału obrazowego i do odmowy; pola puste nie wchodzą do treści.
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

// dopiszCeche dokłada jedną cechę promptu strukturalnego do polecenia
// generowania, gdy wartość cechy została faktycznie wskazana w żądaniu.
func dopiszCeche(b *strings.Builder, nazwa string, wartosc *string) {
	if wartosc == nil || strings.TrimSpace(*wartosc) == "" {
		return
	}
	b.WriteString("; ")
	b.WriteString(nazwa)
	b.WriteString(": ")
	b.WriteString(strings.TrimSpace(*wartosc))
}

// Zasoby zwraca stronę zasobów spełniających filtr wraz z liczbą wszystkich
// spełniających: Total niesie liczbę spełniających warunki, nie długość
// zwróconej strony. Pole PromptId zostaje w tej odpowiedzi puste zawsze.
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

// zasobKontraktu składa DesignAsset kontraktu z wiersza repozytorium. PromptId
// niesie identyfikator zewnętrzny promptu, prawdziwy albo pusty, nigdy
// odgadnięty. Uri wychodzi jako ścieżka względna magazynu, nie jako ścieżka
// na dysku.
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

// sprawdzWyliczenieDesignu odrzuca wartość spoza wyliczenia kontraktu przed
// dotknięciem bazy, dla dowolnego wyliczenia obszaru Design. Wykaz
// dopuszczalnych wartości pochodzi z kontraktu i nie jest tutaj przepisywany.
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

// bladWskazaniaDesignu nazywa brak danych w żądaniu jako usterkę wołającego,
// nie usterkę rdzenia obsługującego żądanie.
func bladWskazaniaDesignu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		fmt.Sprintf("moduł Design: %s", powod)))
}

// bladNieznanegoBytuDesignu nazywa byt modułu, którego rdzeń nie zna: kolekcję,
// kompozycję, wersję, adnotację, zestaw żetonów. Kod not_found odróżnia
// wskazanie nieaktualne od usterki rdzenia.
func bladNieznanegoBytuDesignu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound, "moduł Design: "+powod))
}

// sprawdzRodzajZasobu odrzuca rodzaj zasobu spoza kontraktu przed dotknięciem
// magazynu, żeby pomyłka wołającego wracała jako pomyłka wołającego, a nie
// jako awaria rdzenia nadająca się do ponowienia. Odmowa wymienia dopuszczalne
// rodzaje.
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
