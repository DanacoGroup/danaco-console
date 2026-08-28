// Plik niesie słownik cyklu życia zlecenia kolejki: wartości kolumn stanu i werdyktu weryfikacji
// oraz tabelę przejść jednego kroku pętli koordynator–wykonawca.
package core

import "danacoconsole/server/internal/dane"

const (
	stanPozycjiOczekuje      = "oczekuje"
	stanPozycjiPrzydzielona  = "przydzielona"
	stanPozycjiWykonywana    = "wykonywana"
	stanPozycjiDoWeryfikacji = "do_weryfikacji"
	stanPozycjiUkonczona     = "ukonczona"
	stanPozycjiBledna        = "bledna"
	stanPozycjiAnulowana     = "anulowana"
)

// Werdykt weryfikacji. Schemat zna jeszcze wartość 'odrzucone'; rdzeń jej nie
// wystawia, bo żadne działanie kontraktu jej nie wywołuje, a werdykt bez
// czynności Operatora byłby oceną wymyśloną przez silnik.
const (
	werdyktPrzyjete  = "przyjete"
	werdyktDoPoprawy = "do_poprawy"
)

// stanyKoncowePozycji wylicza stany, po których pozycja nie wraca do pracy sama.
// Powrót jest możliwy wyłącznie biegiem naprawczym wywołanym przez Operatora —
// ten nie ma limitu obiegów ani warunku wstępnego.
var stanyKoncowePozycji = map[string]struct{}{
	stanPozycjiUkonczona: {}, stanPozycjiBledna: {}, stanPozycjiAnulowana: {},
}

// krokNaprzod to tabela przejść pozycji o jeden krok pętli. Stan bez wpisu jest stanem końcowym —
// krok naprzód nic wtedy nie zmienia.
var krokNaprzod = map[string]string{
	stanPozycjiOczekuje:      stanPozycjiWykonywana,
	stanPozycjiPrzydzielona:  stanPozycjiWykonywana,
	stanPozycjiWykonywana:    stanPozycjiDoWeryfikacji,
	stanPozycjiDoWeryfikacji: stanPozycjiUkonczona,
}

// werdyktKroku podaje werdykt zapisywany razem ze stanem docelowym kroku.
// Zamknięcie pozycji niesie werdykt przyjęcia, bo krok naprzód po weryfikacji
// jest właśnie aktem przyjęcia wyniku; pozostałe kroki werdyktu nie mają.
func werdyktKroku(stan string) *string {
	if stan != stanPozycjiUkonczona {
		return nil
	}
	werdykt := werdyktPrzyjete
	return &werdykt
}

// czyStanKoncowyPozycji mówi, czy pozycja w tym stanie nie wróci już do pracy silnika kolejki rdzenia.
func czyStanKoncowyPozycji(stan string) bool {
	_, koncowy := stanyKoncowePozycji[stan]
	return koncowy
}

// pierwszaCzynna wskazuje pozycję, na której stoi kolejka: pierwszą w zapisanej
// kolejności, która nie jest zamknięta. Kolejka pusta i kolejka wyczerpana dają
// brak wskazania, a nie błąd.
func pierwszaCzynna(pozycje []dane.Pozycja) *dane.Pozycja {
	for i := range pozycje {
		if !czyStanKoncowyPozycji(pozycje[i].Stan) {
			return &pozycje[i]
		}
	}
	return nil
}
