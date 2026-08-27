// Odpowiedzialność pliku: jeden proces terminala widziany od strony jego cyklu życia — migawka stanu, domknięcie wynikiem, zakończenie sygnałem, zwolnienie uchwytów.
package core

import (
	"time"

	"danacoconsole/shared"
)

// Migawka oddaje stan procesu w postaci odpornej na równoległą zmianę, czytaną trzema wątkami naraz pod zamkiem.
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

// Zakoncz kończy proces: sygnał wymuszony obejmuje całe drzewo potomstwa jednym uchwytem, sygnał łagodny kończy sam proces polecenia.
func (p *procesTerminala) Zakoncz(wymuszony bool) error {
	if wymuszony {
		return p.drzewo.Ubij()
	}
	if p.uchwyt == nil {
		return p.drzewo.Ubij()
	}
	return p.uchwyt.Ubij()
}

// Zwolnij oddaje uchwyty systemowe po zakończeniu procesu, zamykając zasoby powiązane z jego cyklem życia.
func (p *procesTerminala) Zwolnij() {
	p.drzewo.Zwolnij()
}
