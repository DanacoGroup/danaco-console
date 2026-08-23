package session

import (
	"fmt"
	"os"
	"syscall"
)

// Przejęcie drzewa procesu uruchomionego poza rejestrem procesów okna.
//
// Kanał modelu startuje własny proces (internal/injection), bo tylko on zna
// wiersz poleceń, rotację kont i kształt strumienia. Ubicie drzewa procesów
// pozostaje jednak w jednym miejscu — tutaj.
// Bez przejęcia po ubiciu okna zostałyby wnuki procesu kanału: serwery MCP,
// powłoki narzędziowe i inne potomstwo uruchomione przez model.
//
// Warstwa kanału używa obu części naraz:
//
//	polecenie.SysProcAttr = session.AtrybutyDrzewa()   // przed uruchomieniem
//	drzewo, err := session.PrzejmijDrzewo(proces.Pid()) // zaraz po uruchomieniu
//	defer drzewo.Zwolnij()
//	drzewo.Ubij()                                       // przy zamknięciu okna
//
// Pominięcie AtrybutyDrzewa nie wywraca kanału: na Windows przejęcie zadziała
// mimo to, a na systemach uniksowych ubicie obejmie sam proces zamiast całej
// grupy (brak elementu opcjonalnego nie blokuje uruchomienia).
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

// Ubij kończy przejęty proces wraz z całym jego potomstwem.
func (d *DrzewoProcesu) Ubij() error {
	if d == nil {
		return nil
	}
	return d.drzewo.ubij(d.proces)
}

// Czekaj blokuje wywołującego do faktycznego zakończenia przejętego procesu
// (Windows: czekanie na uchwyt; systemy uniksowe: odpytywanie sygnałem zerowym).
// Wraca natychmiast, gdy procesu już nie ma.
//
// Czekanie nie rusza uchwytów oddawanych przez Zwolnij, więc wolno je prowadzić
// równolegle z Ubij: obserwator obudzi się wtedy, gdy ubicie zrobi swoje.
func (d *DrzewoProcesu) Czekaj() {
	if d == nil {
		return
	}
	d.drzewo.czekaj(d.proces)
}

// Zwolnij oddaje uchwyty systemowe po zakończeniu procesu.
func (d *DrzewoProcesu) Zwolnij() {
	if d == nil {
		return
	}
	d.drzewo.zwolnij()
	_ = d.proces.Release()
}
