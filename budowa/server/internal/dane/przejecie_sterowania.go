// Odpowiedzialność pliku: trwały ślad przejęć i oddań bezpośredniego
// sterowania zleceniem przez Operatora.
//
// Ślad leży w istniejącym dzienniku akcji okna, bez osobnej tabeli: przejęcie
// sterowania jest akcją wykonaną na oknie, a te mają jeden dziennik
// `log_akcji_okna`, pisany metodą `RepozytoriumPrzekazan.ZapiszAkcje` przez
// komendę `window.action`. Kolumna `akcja_id` jest w tej tabeli wolnym tekstem
// — bez ograniczenia CHECK i bez klucza obcego do katalogu akcji — więc oba
// rodzaje wpisu mieszczą się w niej bez zmiany schematu. Druga tabela na to
// samo zdarzenie byłaby drugą prawdą o jednym fakcie.
//
// Plik dokłada wyłącznie odczyt zawężony do tych dwóch rodzajów wpisu; zapis
// idzie istniejącą `ZapiszAkcje`. Odczyt musi być osobny, bo `AkcjeOkna` oddaje
// dziennik okna w całości i z limitem: okno z setką zwykłych akcji zepchnęłoby
// przejęcia poza limit, czyli historia sterowania milczałaby dokładnie tam,
// gdzie jest najbardziej potrzebna.
//
// Parametry i wynik zostają surowe, tak samo jak w `przekazanie_okna_akcje.go`:
// kształt obu kolumn zależy od rodzaju akcji, którego warstwa danych nie zna.
// Ten plik ich nie rozbiera — oddaje `AkcjaOkna` w tej samej postaci, co
// dziennik akcji, a rozbiór robi rdzeń
// (`core/adapter_przejecie_sterowania.go`).
package dane

import (
	"context"
	"fmt"
)

// Rodzaje wpisu dziennika akcji, którymi znaczy się przejęcie sterowania.
//
// Wartości są nazwami śladu w kolumnie `akcja_id`, a nie katalogiem kontraktu:
// `log_akcji_okna` przyjmuje wolny tekst, tak samo jak przy `window.action`,
// gdzie wpisem jest identyfikator akcji katalogu. Brzmią jak komendy, które ten
// ślad wytwarzają, żeby dziennik nazywał zdarzenie tym samym słowem, którym
// Operator go zażądał.
const (
	// AkcjaPrzejeciaSterowania — Operator przejął prowadzenie zlecenia.
	AkcjaPrzejeciaSterowania = "control.takeover"
	// AkcjaOddaniaSterowania — Operator oddał prowadzenie Koordynatorowi.
	AkcjaOddaniaSterowania = "control.release"
)

// RepozytoriumSteru jest wąskim kontraktem odczytu śladu sterowania.
//
// Jest osobny od `RepozytoriumPrzekazan`, choć stoi na tym samym typie:
// interfejs obszaru window.* deklaruje w całości `przekazanie_okna.go`, a wąski
// interfejs pozwala rdzeniowi sięgnąć po tę jedną metodę asercją typu, bez
// dopisywania linii do szerokiego kontraktu. Wzorem jest `wiazaneKolejki` —
// rdzeń sięga tak po zdolność, której szeroki port nie ogłasza.
type RepozytoriumSteru interface {
	// SladySterowania zwraca przejęcia i oddania sterowania nad wskazanym
	// oknem, od najnowszego. Okno bez ani jednego przejęcia oddaje wykaz pusty,
	// nie błąd.
	SladySterowania(ctx context.Context, okno string, limit int) ([]AkcjaOkna, error)
}

// Typ repozytorium obszaru window.* wypełnia ten kontrakt. Gdyby metoda znikła
// albo zmieniła kształt, kompilacja stanie tutaj, a nie na martwej komendzie.
var _ RepozytoriumSteru = (*repozytoriumPrzekazan)(nil)

const pobierzSladySterowania = `SELECT log.id, ok.identyfikator_zewnetrzny, log.akcja_id, log.parametry,
                                       log.wynik, log.utworzono
                                FROM log_akcji_okna log
                                JOIN okno_komunikacji ok ON ok.id = log.okno_komunikacji_id
                                WHERE ok.identyfikator_zewnetrzny = ?
                                  AND log.akcja_id IN (?, ?)
                                ORDER BY log.utworzono DESC, log.id DESC
                                LIMIT ?`

// SladySterowania oddaje historię przejęć i oddań sterowania nad oknem.
//
// Limit niedodatni schodzi na domyślny — ten sam, co w dzienniku akcji. Odczyt
// bez granicy byłby pułapką wydajności przy oknie prowadzonym tygodniami.
func (r *repozytoriumPrzekazan) SladySterowania(ctx context.Context, okno string, limit int) ([]AkcjaOkna, error) {
	if okno == "" {
		return nil, fmt.Errorf("dane: ślad sterowania bez identyfikatora okna")
	}
	if limit <= 0 {
		limit = limitDomyslnyAkcjiOkna
	}

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSladySterowania)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno,
		AkcjaPrzejeciaSterowania, AkcjaOddaniaSterowania, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać śladu sterowania okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []AkcjaOkna{}
	for wiersze.Next() {
		akcja, err := odczytajAkcjeOkna(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny ślad sterowania okna %q: %w", okno, err)
		}
		lista = append(lista, akcja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt śladu sterowania okna %q: %w", okno, err)
	}
	return lista, nil
}
