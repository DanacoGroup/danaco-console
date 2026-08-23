// Odpowiedzialność pliku: obecność Operatorów na kompozycji Design Board
// (`design.presence.report`) — kursory współpracy z opracowania modułu.
// Metody stoją na `*adapterDesignu` (`adapter_modul_design.go`).
//
// ── Obecność jest ULOTNA i nie ma wiersza w bazie ───────────────────────────
// Kontrakt mówi to wprost, a powód jest prosty: położenie kursora sprzed
// godziny nie jest wiedzą o niczym. Rejestr żyje w pamięci rdzenia i ginie
// razem z procesem — tak samo, jak ginie sesja, w której ten kursor był.
// Tabeli dla obecności nie ma i nie ma być: wiersze zapisywane dziesięć razy na
// sekundę na klienta byłyby zapisem do dysku, którego nikt nigdy nie odczyta.
//
// ── Zgłoszenie stare przestaje być obecnością ───────────────────────────────
// Klient, który zamknął kartę bez zgłoszenia odejścia, zostawia po sobie wpis.
// Bez terminu ważności jego kursor wisiałby na cudzej kanwie do restartu
// rdzenia — czyli kłamał o obecności kogoś, kogo nie ma. Wpis starszy niż
// `terminObecnosciDesignu` wypada przy najbliższym zgłoszeniu; sprzątanie
// dzieje się przy odczycie, bez własnego budzika, bo bez ruchu na kompozycji
// nie ma też komu tego kursora pokazywać.
package core

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"danacoconsole/shared"
)

// terminObecnosciDesignu jest czasem, po którym zgłoszenie przestaje znaczyć
// obecność. Trzydzieści sekund: dłużej niż odstęp między ruchami kursora
// w normalnej pracy, krócej niż czas, po którym kursor duch zaczyna mylić
// pozostałych.
const terminObecnosciDesignu = 30 * time.Second

// rejestrObecnosciDesignu trzyma obecnych na kompozycjach. Klucz zewnętrzny to
// kompozycja, wewnętrzny — klient; jeden klient jest obecny na jednej
// kompozycji raz, a kolejne zgłoszenie przestawia jego położenie, zamiast
// dokładać drugi kursor.
type rejestrObecnosciDesignu struct {
	zamek       sync.Mutex
	kompozycje  map[string]map[string]shared.DesignPresence
	terazTestem func() time.Time
}

// nowyRejestrObecnosciDesignu zakłada pusty rejestr.
func nowyRejestrObecnosciDesignu() *rejestrObecnosciDesignu {
	return &rejestrObecnosciDesignu{kompozycje: map[string]map[string]shared.DesignPresence{}}
}

// teraz oddaje chwilę bieżącą. Pole zastępowalne, żeby sprawdzian terminu nie
// musiał czekać trzydziestu sekund realnego czasu.
func (r *rejestrObecnosciDesignu) teraz() time.Time {
	if r.terazTestem != nil {
		return r.terazTestem()
	}
	return time.Now()
}

// zglos odnotowuje obecność klienta na kompozycji albo ją zdejmuje i oddaje
// obecnych po zmianie, uporządkowanych po kliencie — porządek stały, żeby dwa
// kolejne zgłoszenia nie przestawiały kursorów na ekranie bez powodu.
func (r *rejestrObecnosciDesignu) zglos(kompozycja string, obecnosc shared.DesignPresence,
	odchodzi bool) []shared.DesignPresence {

	r.zamek.Lock()
	defer r.zamek.Unlock()

	wpisy, znane := r.kompozycje[kompozycja]
	if !znane {
		wpisy = map[string]shared.DesignPresence{}
		r.kompozycje[kompozycja] = wpisy
	}
	if odchodzi {
		delete(wpisy, obecnosc.ClientId)
	} else {
		wpisy[obecnosc.ClientId] = obecnosc
	}

	granica := r.teraz().Add(-terminObecnosciDesignu).UnixMilli()
	obecni := make([]shared.DesignPresence, 0, len(wpisy))
	for klient, wpis := range wpisy {
		if wpis.SeenAt < granica {
			delete(wpisy, klient)
			continue
		}
		obecni = append(obecni, wpis)
	}
	if len(wpisy) == 0 {
		// Kompozycja bez obecnych nie ma po co zajmować miejsca w rejestrze:
		// mapa rosnąca o wpis na każdą kiedykolwiek otwartą tablicę byłaby
		// wyciekiem pamięci rozłożonym na tygodnie.
		delete(r.kompozycje, kompozycja)
	}
	sort.Slice(obecni, func(i, j int) bool { return obecni[i].ClientId < obecni[j].ClientId })
	return obecni
}

// ZglosObecnosc odnotowuje obecność i położenie kursora Operatora na
// kompozycji — obsługuje `design.presence.report`.
//
// Kompozycję sprawdzamy w bazie mimo ulotności zgłoszenia: obecność na tablicy,
// której nie ma, rozgłaszałaby kursory na kompozycji, do której nikt nigdy nie
// zajrzy, a literówka w identyfikatorze wyglądałaby jak cisza współpracowników.
func (a *adapterDesignu) ZglosObecnosc(ctx context.Context,
	z shared.DesignPresenceReportRequest) (shared.DesignPresenceReportResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignPresenceReportResponse{}, bladWskazaniaDesignu(
			"komenda design.presence.report bez wskazania kompozycji")
	}
	if _, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId)); err != nil {
		return shared.DesignPresenceReportResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}

	_, klient := sprawca(ctx)
	if klient == nil || strings.TrimSpace(*klient) == "" {
		return shared.DesignPresenceReportResponse{}, bladWskazaniaDesignu(
			"zgłoszenie obecności bez rozpoznanego klienta — kursor bez tożsamości nie da się " +
				"odróżnić od cudzego ani zdjąć przy odejściu")
	}
	if a.obecnosc == nil {
		a.obecnosc = nowyRejestrObecnosciDesignu()
	}

	obecni := a.obecnosc.zglos(strings.TrimSpace(z.BoardId), shared.DesignPresence{
		ClientId:         strings.TrimSpace(*klient),
		X:                z.X,
		Y:                z.Y,
		SelectedLayerIds: z.SelectedLayerIds,
		SeenAt:           time.Now().UnixMilli(),
	}, z.Leaving != nil && *z.Leaving)

	return shared.DesignPresenceReportResponse{Participants: obecni}, nil
}
