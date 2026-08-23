//go:build !windows

package session

import (
	"fmt"
	"syscall"
)

// wstrzymaj zatrzymuje całą grupę procesów sygnałem SIGSTOP.
//
// Sygnał idzie do UJEMNEGO identyfikatora grupy, tak samo jak ubicie: proces
// okna dostaje własną grupę już przy starcie (atrybutyProcesu → Setpgid), więc
// jedno wywołanie obejmuje całe drzewo. Wstrzymanie samego korzenia zostawiłoby
// biegnące potomstwo, czyli tę część pracy, która zwykle zajmuje maszynę.
//
// Straż `pid > 1` pilnuje tego samego, co przy ubiciu: sygnał do -1 byłby
// rozgłoszeniem do wszystkich procesów systemu, a do 1 — sygnałem do init.
// ESRCH znaczy, że nie ma już czego wstrzymywać, i nie jest błędem.
func (d *drzewoProcesow) wstrzymaj() (bool, error) {
	return d.sygnalDrzewa(syscall.SIGSTOP, "wstrzymać")
}

// wznow podejmuje pracę wstrzymanej grupy sygnałem SIGCONT.
func (d *drzewoProcesow) wznow() (bool, error) {
	return d.sygnalDrzewa(syscall.SIGCONT, "wznowić")
}

// sygnalDrzewa wysyła sygnał całej grupie, a przy jej braku samemu procesowi.
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
