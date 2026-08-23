// Odpowiedzialność pliku: przygotowanie zapytań SQL i ich pamięć podręczna.
// Każde zapytanie repozytoriów przechodzi przez `sql.Stmt` — sterownik dostaje
// gotowy plan, a wartości wyłącznie jako parametry, nigdy jako sklejony tekst.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

// zapytania przechowuje przygotowane polecenia. Klucz mapy to treść zapytania,
// dzięki czemu repozytorium podaje SQL, a nie identyfikator z osobnego rejestru.
type zapytania struct {
	db     *sql.DB
	mutex  sync.Mutex
	pamiec map[string]*sql.Stmt
}

// noweZapytania zakłada pamięć podręczną nad otwartą pulą połączeń.
func noweZapytania(db *sql.DB) *zapytania {
	return &zapytania{db: db, pamiec: map[string]*sql.Stmt{}}
}

// przygotuj zwraca polecenie przygotowane dla podanego zapytania. Pierwsze
// wywołanie przygotowuje je w bazie, kolejne korzystają z pamięci podręcznej.
func (z *zapytania) przygotuj(ctx context.Context, tekst string) (*sql.Stmt, error) {
	z.mutex.Lock()
	polecenie, jest := z.pamiec[tekst]
	z.mutex.Unlock()
	if jest {
		return polecenie, nil
	}
	polecenie, err := z.db.PrepareContext(ctx, tekst)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można przygotować zapytania %q: %w", skrot(tekst), err)
	}
	z.mutex.Lock()
	defer z.mutex.Unlock()
	if istniejace, jest := z.pamiec[tekst]; jest {
		// Inny wątek zdążył przygotować to samo zapytanie — zwracamy jego wersję,
		// a własną zamykamy, żeby nie zostawiać osieroconego polecenia.
		polecenie.Close()
		return istniejace, nil
	}
	z.pamiec[tekst] = polecenie
	return polecenie, nil
}

// wTransakcji zwraca to samo przygotowane polecenie związane z transakcją.
func (z *zapytania) wTransakcji(ctx context.Context, transakcja *sql.Tx, tekst string) (*sql.Stmt, error) {
	polecenie, err := z.przygotuj(ctx, tekst)
	if err != nil {
		return nil, err
	}
	return transakcja.StmtContext(ctx, polecenie), nil
}

// zamknij zwalnia wszystkie przygotowane polecenia. Błąd pierwszego zamknięcia
// nie przerywa zwalniania pozostałych.
func (z *zapytania) zamknij() error {
	z.mutex.Lock()
	defer z.mutex.Unlock()
	var pierwszy error
	for tekst, polecenie := range z.pamiec {
		if err := polecenie.Close(); err != nil && pierwszy == nil {
			pierwszy = fmt.Errorf("dane: nie można zamknąć zapytania %q: %w", skrot(tekst), err)
		}
		delete(z.pamiec, tekst)
	}
	return pierwszy
}

// skrot skraca treść zapytania w komunikacie błędu do pierwszego wiersza.
func skrot(tekst string) string {
	const granica = 60
	if len(tekst) <= granica {
		return tekst
	}
	return tekst[:granica] + "…"
}
