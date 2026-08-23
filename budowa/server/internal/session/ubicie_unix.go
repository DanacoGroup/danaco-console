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

// Drzewo procesów okna na systemach uniksowych opiera się na grupie procesów.
// Proces okna zakłada własną grupę, potomstwo grupę dziedziczy, więc sygnał
// wysłany do ujemnego identyfikatora grupy kończy całe drzewo naraz.
//
// pid to zapamiętany identyfikator grupy procesów (przez Setpgid równy pidowi
// procesu okna) utrwalony w chwili przejęcia, kiedy jest znany i dodatni.
// Ubicie posługuje się tym polem, a nie os.Process.Pid — to drugie zeruje
// Release na wartość -1, co daje zarazem wyścig danych i policzenie
// -p.Pid = 1 (init) albo Kill(-1) (rozgłoszenie do wszystkich procesów). Pole
// zapisuje się raz, przed jakąkolwiek współbieżnością, i tylko czyta później,
// więc Release go nie tyka.
type drzewoProcesow struct {
	pid int
}

// atrybutyProcesu zakłada procesowi okna własną grupę procesów.
func atrybutyProcesu() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}

// przygotujDrzewo nie ma na tej platformie nic do przygotowania — grupa
// powstaje wraz z procesem.
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

// ubij wysyła sygnał zakończenia całej grupie procesów okna. Proces okna dostaje
// własną grupę już przy starcie (atrybutyProcesu → Setpgid), więc sygnał do
// ujemnego identyfikatora obejmuje całe drzewo. Gdy grupy nie ma (proces zdążył
// ją zmienić albo już jej nie ma), sygnał trafia wprost do samego procesu, żeby
// nie zostawić go przy życiu. ESRCH oznacza, że nie ma już czego ubijać.
//
// Posługuje się zapamiętanym d.pid, nie os.Process.Pid (Release zeruje go na
// -1). Wartość pid <= 1 znaczy brak prawidłowego procesu do ubicia i nie idzie
// wtedy żaden sygnał. Straż pid > 1 pilnuje zarazem sygnału do grupy
// (-pid < -1) i sygnału bezpośredniego (pid > 1), żeby żaden nie wyrodził się
// w Kill(1) — init — ani w Kill(-1) — rozgłoszenie do wszystkich procesów.
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

// zwolnij nie ma na tej platformie uchwytu do oddania.
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
