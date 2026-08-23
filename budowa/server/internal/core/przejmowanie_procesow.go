package core

import "sync"

// przejmowanieProcesow jest pośredniczką między kanałem modelu a rejestrem
// procesów sesji.
//
// Istnieje z powodu kolejności montażu: rejestr kanałów powstaje przed nadzorcą
// sesji, a to nadzorca ma rejestr procesów. Kanał musi więc dostać haczyk
// wcześniej, niż istnieje jego odbiorca. Pośredniczka rozwiązuje to jawnie —
// jedna wąska rzecz, którą widać w montażu — zamiast przestawiania kolejności
// całego składania rdzenia albo zmiennej pakietowej.
//
// Do chwili związania haczyk jest bezczynny: tura biegnie, tylko proces nie jest
// jeszcze obejmowany uchwytem. To stan przejściowy trwający ułamek montażu, a nie
// tryb pracy.
type przejmowanieProcesow struct {
	mu       sync.RWMutex
	przejmij func(idOkna string, pid int) error
}

// Zwiaz wpina rzeczywistego odbiorcę po złożeniu nadzorcy.
func (p *przejmowanieProcesow) Zwiaz(przejmij func(idOkna string, pid int) error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.przejmij = przejmij
}

// Haczyk zwraca funkcję dla jednego okna. Kanał wkłada ją do zapytania, a
// injection woła ją zaraz po starcie procesu tury.
//
// Niepowodzenie przejęcia jest pochłaniane: proces odpowiada i tura ma się
// odbyć. Utrata sprzątania boli, ale przerwana rozmowa boli bardziej.
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
