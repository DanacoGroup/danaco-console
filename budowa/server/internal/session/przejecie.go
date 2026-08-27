package session

import (
	"fmt"
	"os"
	"syscall"
)

// DrzewoProcesu reprezentuje przejęcie drzewa procesu uruchomionego poza rejestrem procesów okna, przez kanał modelu poza tym pakietem.
type DrzewoProcesu struct {
	drzewo *drzewoProcesow
	proces *os.Process
}

// AtrybutyDrzewa zwraca atrybuty uruchomienia, po których proces staje się
// korzeniem własnego drzewa: grupą procesów na systemach uniksowych, procesem
// wstrzymanym do czasu przejęcia na Windows.
func AtrybutyDrzewa() *syscall.SysProcAttr {
	return atrybutyProcesu()
}

// PrzejmijDrzewo bierze biegnący proces wraz z jego przyszłym potomstwem pod
// jeden uchwyt systemowy. Na Windows kończy zarazem wstrzymanie nadane przez
// AtrybutyDrzewa, więc wywołuje się je zaraz po uruchomieniu procesu.
func PrzejmijDrzewo(pid int) (*DrzewoProcesu, error) {
	proces, err := os.FindProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("session: nie można odnaleźć procesu %d: %w", pid, err)
	}
	drzewo, err := przygotujDrzewo()
	if err != nil {
		_ = proces.Release()
		return nil, err
	}
	if err := drzewo.przejmij(proces); err != nil {
		drzewo.zwolnij()
		_ = proces.Release()
		return nil, err
	}
	return &DrzewoProcesu{drzewo: drzewo, proces: proces}, nil
}

// Metoda Ubij kończy przejęty proces wraz z całym jego potomstwem, korzystając z uchwytu systemowego drzewa.
func (d *DrzewoProcesu) Ubij() error {
	if d == nil {
		return nil
	}
	return d.drzewo.ubij(d.proces)
}

// Metoda Czekaj blokuje wywołującego do faktycznego zakończenia przejętego procesu i wraca natychmiast, gdy procesu już nie ma.
func (d *DrzewoProcesu) Czekaj() {
	if d == nil {
		return
	}
	d.drzewo.czekaj(d.proces)
}

// Metoda Zwolnij oddaje uchwyty systemowe przejętego drzewa procesu po jego faktycznym zakończeniu pracy.
func (d *DrzewoProcesu) Zwolnij() {
	if d == nil {
		return
	}
	d.drzewo.zwolnij()
	_ = d.proces.Release()
}
