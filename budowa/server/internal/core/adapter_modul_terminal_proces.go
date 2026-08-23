// Odpowiedzialność pliku: jeden proces terminala widziany od strony jego cyklu
// życia — migawka stanu, domknięcie wynikiem, zakończenie sygnałem i zwolnienie
// uchwytów systemowych.
//
// Stan procesu czytają trzy wątki naraz: obsługiwacz komendy, pompa wyjścia
// i obserwator zakończenia. Dlatego każdy odczyt idzie migawką pod zamkiem,
// a nie wprost po polach — inaczej Process Monitor pokazywałby stan wpisany
// w połowie.
package core

import (
	"time"

	"danacoconsole/shared"
)

// Migawka oddaje stan procesu w postaci odpornej na równoległą zmianę.
func (p *procesTerminala) Migawka() (shared.TerminalProcessStatus, *int, time.Time) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.kodWyjscia == nil {
		return p.stan, nil, p.zakonczono
	}
	kod := *p.kodWyjscia
	return p.stan, &kod, p.zakonczono
}

// Domknij zapisuje stan końcowy procesu. Zwraca fałsz, gdy proces był już
// domknięty — pierwszy prawdziwy wynik nie ma prawa zostać nadpisany przez
// późniejsze ubicie.
func (p *procesTerminala) Domknij(stan shared.TerminalProcessStatus, kodWyjscia *int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stan != shared.TerminalProcessStatusRunning {
		return false
	}
	p.stan, p.kodWyjscia, p.zakonczono = stan, kodWyjscia, time.Now().UTC()
	return true
}

// Zakoncz kończy proces. Sygnał wymuszony obejmuje całe drzewo potomstwa
// jednym uchwytem systemowym; sygnał łagodny kończy sam proces polecenia
// i zostawia jego potomstwo przy życiu.
//
// Rozróżnienie nie odwzorowuje pary SIGTERM/SIGKILL, ponieważ program pracuje
// także na Windows, gdzie sygnału łagodnego dla obcego procesu nie ma.
func (p *procesTerminala) Zakoncz(wymuszony bool) error {
	if wymuszony {
		return p.drzewo.Ubij()
	}
	if p.uchwyt == nil {
		return p.drzewo.Ubij()
	}
	return p.uchwyt.Ubij()
}

// Zwolnij oddaje uchwyty systemowe po zakończeniu procesu.
func (p *procesTerminala) Zwolnij() {
	p.drzewo.Zwolnij()
}
