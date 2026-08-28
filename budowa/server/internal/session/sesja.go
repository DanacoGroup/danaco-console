package session

import (
	"time"

	"danacoconsole/shared"
)

// Sesja jest przestrzenią wspólną: plikami, pamięcią, projektem i agentami.
// Nie jest jednostką wykonania — parametry wykonania niesie okno komunikacji,
// a sesja tylko gromadzi okna.
type Sesja struct {
	// Id sesji.
	Id string
	// Tytul sesji; pusty jest dopuszczalny.
	Tytul string
	// IdProjektu wiąże sesję z projektem; pusty oznacza sesję bez projektu.
	IdProjektu string
	// Stan sesji ze słownika kontraktu (shared.SessionStatus).
	Stan shared.SessionStatus
	// IdOkien wylicza okna sesji w kolejności otwarcia; jedna sesja ma wiele
	// okien.
	IdOkien []string
	// Utworzono i Zaktualizowano znakują cykl życia sesji.
	Utworzono      time.Time
	Zaktualizowano time.Time
}

// nowaSesja zakłada sesję czynną. Brak tytułu ani projektu niczego nie blokuje
// — sesja bez nazwy jest poprawną sesją roboczą.
func nowaSesja(tytul, idProjektu string) *Sesja {
	teraz := time.Now().UTC()
	return &Sesja{
		Id:             nowyIdentyfikator(przedrostekSesji),
		Tytul:          tytul,
		IdProjektu:     idProjektu,
		Stan:           shared.SessionStatusActive,
		IdOkien:        nil,
		Utworzono:      teraz,
		Zaktualizowano: teraz,
	}
}

// Kopia zwraca niezależny odpis sesji. Rejestr wydaje wyłącznie odpisy, więc
// wywołujący nie może przez wspólny wycinek zmienić stanu rejestru.
func (s Sesja) Kopia() Sesja {
	odpis := s
	if s.IdOkien != nil {
		odpis.IdOkien = append([]string(nil), s.IdOkien...)
	}
	return odpis
}

// Metoda CzyCzynna mówi, czy sesja przyjmuje pracę i pozostaje aktywna w tym bieżącym rejestrze sesji rdzenia.
func (s Sesja) CzyCzynna() bool {
	return s.Stan == shared.SessionStatusActive
}

// Metoda dopiszOkno dokłada okno do wykazu sesji, pilnując braku powtórzeń identyfikatora tego samego okna.
func (s *Sesja) dopiszOkno(idOkna string) {
	for _, istniejace := range s.IdOkien {
		if istniejace == idOkna {
			return
		}
	}
	s.IdOkien = append(s.IdOkien, idOkna)
	s.Zaktualizowano = time.Now().UTC()
}

// usunOkno wykreśla okno z wykazu sesji. Sesja trwa dalej także bez okien —
// zamknięcie okna nie kończy sesji.
func (s *Sesja) usunOkno(idOkna string) {
	pozostale := s.IdOkien[:0]
	for _, istniejace := range s.IdOkien {
		if istniejace != idOkna {
			pozostale = append(pozostale, istniejace)
		}
	}
	s.IdOkien = pozostale
	s.Zaktualizowano = time.Now().UTC()
}
