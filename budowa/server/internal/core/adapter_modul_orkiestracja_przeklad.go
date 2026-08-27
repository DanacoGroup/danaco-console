// Odpowiedzialność pliku: przekład wiersza podagenta na strukturę
// `Subagent` kontraktu, odwzorowanie stanu pozycji kolejki na
// `SubagentStatus` oraz kody odmów rodziny `subagent.*`.
package core

import (
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// podagenciKontraktu przekłada wykaz wierszy na wykaz kontraktu. Lista pusta
// wychodzi jako lista pusta, nigdy jako `null` — panel ma dostać wykaz, choćby
// bez pozycji.
func podagenciKontraktu(wiersze []dane.Podagent) []shared.Subagent {
	podagenci := make([]shared.Subagent, 0, len(wiersze))
	for _, wiersz := range wiersze {
		podagenci = append(podagenci, podagentKontraktu(wiersz))
	}
	return podagenci
}

// podagentKontraktu składa strukturę `Subagent` z wiersza tabeli. Pola
// czasu wychodzą tylko wtedy, gdy chwila naprawdę nastąpiła.
func podagentKontraktu(wiersz dane.Podagent) shared.Subagent {
	podagent := shared.Subagent{
		Id:       wiersz.Kod,
		WindowId: identyfikatorWiersza(wiersz.OknoKod, wiersz.OknoID),
		Task:     wiersz.Zadanie,
		Status:   shared.SubagentStatus(wiersz.Stan),
		Name:     wiersz.Nazwa,
		Result:   wiersz.Wynik,
	}
	if wiersz.SesjaKod != nil && *wiersz.SesjaKod != "" {
		podagent.SessionId = wiersz.SesjaKod
	}
	if wiersz.PozycjaKolejkiID != nil {
		pozycja := identyfikatorWiersza(nil, *wiersz.PozycjaKolejkiID)
		podagent.QueueItemId = &pozycja
	}
	if wiersz.Rozpoczeto != nil && *wiersz.Rozpoczeto != "" {
		chwila := chwilaBazy(*wiersz.Rozpoczeto)
		podagent.StartedAt = &chwila
	}
	if wiersz.Zakonczono != nil && *wiersz.Zakonczono != "" {
		chwila := chwilaBazy(*wiersz.Zakonczono)
		podagent.FinishedAt = &chwila
	}
	return podagent
}

// stanPodagentaZPozycji odwzorowuje słownik stanów pozycji kolejki na
// wyliczenie `SubagentStatus` kontraktu. Stan nieznany zostaje „w biegu”,
// nie „błędny”.
func stanPodagentaZPozycji(stanPozycji string) string {
	switch stanPozycji {
	case stanPozycjiOczekuje, stanPozycjiPrzydzielona:
		return dane.StanPodagentaOczekuje
	case stanPozycjiDoWeryfikacji, stanPozycjiUkonczona:
		return dane.StanPodagentaUkonczony
	case stanPozycjiBledna:
		return dane.StanPodagentaBledny
	case stanPozycjiAnulowana:
		return dane.StanPodagentaZatrzymany
	default:
		return dane.StanPodagentaWBiegu
	}
}

// bladPodagentow znakuje usterkę kodem kontraktu, żeby panel pokazał powód,
// a nie samo „nie udało się". Błąd, któremu kod już nadano, przechodzi tędy bez
// zmiany kodu; dopiero usterka bez kodu staje się usterką wewnętrzną rdzenia.
func bladPodagentow(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaPodagenta nazywa brak danych w żądaniu — to błąd klienta,
// nie rdzenia, więc kod odmowy jest inny niż przy usterce.
func bladWskazaniaPodagenta(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"podagenci: "+powod))
}

// bladNieznanegoOknaPodagenta odróżnia „okna nie ma" od „odczyt się nie powiódł".
// Panel pokazuje wtedy inny komunikat i inaczej podpowiada Operatorowi.
func bladNieznanegoOknaPodagenta(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"podagenci: okno wykonawcy nie istnieje: "+kod))
	}
	return bladPodagentow(err)
}

// bladNieznanychPodagentow nazywa wskazanie, któremu nie odpowiada żaden wiersz.
// Odmowa niesie wskazane identyfikatory, bo bez nich Operator nie wie, który
// z nich był chybiony.
func bladNieznanychPodagentow(kody []string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"podagenci: żaden ze wskazanych podagentów nie istnieje: "+strings.Join(kody, ", ")))
}
