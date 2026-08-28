//go:build !windows

package session

import (
	"fmt"
	"syscall"
)

// wstrzymaj zatrzymuje całą grupę procesów sygnałem SIGSTOP wysyłanym do ujemnego identyfikatora grupy, obejmując tym samym wywołaniem całe drzewo potomstwa procesu okna.
func (d *drzewoProcesow) wstrzymaj() (bool, error) {
	return d.sygnalDrzewa(syscall.SIGSTOP, "wstrzymać")
}

// wznow podejmuje pracę wstrzymanej grupy procesów sygnałem SIGCONT, przywracając jej dostęp do czasu procesora dokładnie tam, gdzie została zatrzymana.
func (d *drzewoProcesow) wznow() (bool, error) {
	return d.sygnalDrzewa(syscall.SIGCONT, "wznowić")
}

// sygnalDrzewa wysyła podany sygnał całej grupie procesów, a przy jej braku kieruje ten sam sygnał bezpośrednio do samego procesu.
func (d *drzewoProcesow) sygnalDrzewa(sygnal syscall.Signal, czynnosc string) (bool, error) {
	pid := d.pid
	if pid <= 1 {
		return true, nil
	}
	if err := syscall.Kill(-pid, sygnal); err != nil {
		if err == syscall.ESRCH {
			return true, nil
		}
		if bezposredni := syscall.Kill(pid, sygnal); bezposredni != nil {
			if bezposredni == syscall.ESRCH {
				return true, nil
			}
			return true, fmt.Errorf("session: nie można %s grupy ani procesu %d: %w",
				czynnosc, pid, bezposredni)
		}
	}
	return true, nil
}
