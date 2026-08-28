// Plik wybiera cel działania i wyprowadza stan kolejki, dopełniając silnik wykonania pytaniami
// pobocznymi, które da się sprawdzić bez dotykania bazy.
package core

import (
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// celDzialania wskazuje pozycję, której dotyczy działanie. Pozycja wskazana
// w żądaniu, a nieznana kolejce, jest odmową: milczące przejście na inną
// pozycję wykonałoby co innego, niż zlecono.
func celDzialania(pozycje []dane.Pozycja, idPozycji *string) (*dane.Pozycja, error) {
	if !czyWskazano(idPozycji) {
		return pierwszaCzynna(pozycje), nil
	}
	if id, err := strconv.ParseInt(strings.TrimSpace(*idPozycji), 10, 64); err == nil {
		for i := range pozycje {
			if pozycje[i].ID == id {
				return &pozycje[i], nil
			}
		}
	}
	return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"pozycja kolejki nie istnieje: "+*idPozycji))
}

// celNaprawy wskazuje pozycję biegu naprawczego. Kolejka wyczerpana powtarza
// pozycję ostatnią — powtórzenie kroku nie ma warunku wstępnego.
func celNaprawy(pozycje []dane.Pozycja, cel *dane.Pozycja) *dane.Pozycja {
	if cel != nil || len(pozycje) == 0 {
		return cel
	}
	return &pozycje[len(pozycje)-1]
}

// przerywane wylicza pozycje objęte przerwaniem: wskazaną pozycję albo wszystkie pozycje czynne kolejki.
func przerywane(pozycje []dane.Pozycja, cel *dane.Pozycja, wskazano bool) []dane.Pozycja {
	if wskazano {
		if cel == nil {
			return nil
		}
		return []dane.Pozycja{*cel}
	}
	czynne := make([]dane.Pozycja, 0, len(pozycje))
	for _, pozycja := range pozycje {
		if !czyStanKoncowyPozycji(pozycja.Stan) {
			czynne = append(czynne, pozycja)
		}
	}
	return czynne
}

// czyWskazano mówi, czy żądanie wskazało konkretną pozycję kolejki identyfikatorem wprost w treści żądania.
func czyWskazano(idPozycji *string) bool {
	return idPozycji != nil && strings.TrimSpace(*idPozycji) != ""
}

// stanKolejkiPoDzialaniu wyprowadza stan kolejki z działania i z pozycji po zmianie. Wstrzymanie,
// zatrzymanie i opróżnienie nazywa wprost samo działanie; przy pracy rozstrzygają pozycje.
func stanKolejkiPoDzialaniu(dzialanie shared.QueueAction, pozycje []dane.Pozycja) shared.QueueStatus {
	switch dzialanie {
	case shared.QueueActionPause:
		return shared.QueueStatusPaused
	case shared.QueueActionStop:
		return shared.QueueStatusStopped
	case shared.QueueActionClear:
		return shared.QueueStatusIdle
	}
	if len(pozycje) > 0 && pierwszaCzynna(pozycje) == nil {
		return shared.QueueStatusDone
	}
	return shared.QueueStatusRunning
}
