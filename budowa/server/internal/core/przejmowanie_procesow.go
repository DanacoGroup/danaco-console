package core

import "sync"

// przejmowanieProcesow jest pośredniczką między kanałem modelu a rejestrem procesów sesji, potrzebną, bo rejestr kanałów powstaje przed nadzorcą sesji trzymającym ten rejestr.
type przejmowanieProcesow struct {
	mu       sync.RWMutex
	przejmij func(idOkna string, pid int) error
}

// Zwiaz wpina do pośredniczki rzeczywistego odbiorcę przejęcia procesu, dostępnego dopiero po złożeniu nadzorcy sesji.
func (p *przejmowanieProcesow) Zwiaz(przejmij func(idOkna string, pid int) error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.przejmij = przejmij
}

// Haczyk zwraca funkcję przejęcia procesu dla jednego okna; niepowodzenie przejęcia jest pochłaniane, bo przerwana rozmowa boli bardziej niż utracone sprzątanie.
func (p *przejmowanieProcesow) Haczyk(idOkna string) func(pid int) {
	if p == nil {
		return nil
	}
	return func(pid int) {
		p.mu.RLock()
		przejmij := p.przejmij
		p.mu.RUnlock()
		if przejmij == nil || idOkna == "" {
			return
		}
		_ = przejmij(idOkna, pid)
	}
}
