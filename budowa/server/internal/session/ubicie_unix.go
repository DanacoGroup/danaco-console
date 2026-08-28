//go:build !windows

package session

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// odstepDogladu wyznacza rytm sprawdzania, czy przejęty proces jeszcze żyje.
// Przejętego procesu nie da się `wait()` (nie jest dzieckiem tego procesu), więc
// jego zakończenie rozpoznaje się odpytywaniem sygnałem zerowym.
const odstepDogladu = 200 * time.Millisecond

// drzewoProcesow reprezentuje drzewo procesu okna na systemach uniksowych, oparte na grupie procesów systemowych.
type drzewoProcesow struct {
	pid int
}

// Funkcja atrybutyProcesu zakłada procesowi okna jego własną grupę procesów, oddzielną od procesu rdzenia.
func atrybutyProcesu() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}

// Funkcja przygotujDrzewo nie ma na tej platformie nic do przygotowania, ponieważ grupa powstaje wraz z procesem.
func przygotujDrzewo() (*drzewoProcesow, error) {
	return &drzewoProcesow{}, nil
}

// przejmij utrwala identyfikator grupy procesów w chwili, gdy pid jest jeszcze
// znany i dodatni. Samego przypisania do grupy nie ma co wykonywać — proces
// należy do własnej grupy od utworzenia (Setpgid).
func (d *drzewoProcesow) przejmij(p *os.Process) error {
	if p != nil {
		d.pid = p.Pid
	}
	return nil
}

// Metoda ubij wysyła sygnał zakończenia całej grupie procesów okna, obejmując tym samym całe drzewo potomstwa.
func (d *drzewoProcesow) ubij(_ *os.Process) error {
	pid := d.pid
	if pid <= 1 {
		return nil
	}
	if err := syscall.Kill(-pid, syscall.SIGKILL); err != nil {
		if err == syscall.ESRCH {
			return nil
		}
		if bezposredni := syscall.Kill(pid, syscall.SIGKILL); bezposredni != nil {
			if bezposredni == syscall.ESRCH {
				return nil
			}
			return fmt.Errorf("session: nie można zakończyć grupy ani procesu %d: %w", pid, bezposredni)
		}
	}
	return nil
}

// Metoda zwolnij nie ma na tej platformie żadnego uchwytu systemowego do oddania po zakończeniu procesu.
func (d *drzewoProcesow) zwolnij() {}

// czekaj blokuje wywołującego do faktycznego zakończenia przejętego procesu.
// Sygnał zerowy nie robi procesowi nic — służy wyłącznie sprawdzeniu, czy proces
// jeszcze istnieje. Błąd (ESRCH) znaczy, że proces zniknął.
func (d *drzewoProcesow) czekaj(p *os.Process) {
	if p == nil || d.pid <= 1 {
		return
	}
	for {
		if err := syscall.Kill(d.pid, 0); err != nil {
			return
		}
		time.Sleep(odstepDogladu)
	}
}
