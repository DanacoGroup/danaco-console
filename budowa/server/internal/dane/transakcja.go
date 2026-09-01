// Odpowiedzialność pliku: wykonanie zapisu obejmującego wiele tabel w jednej
// transakcji. Zapis albo dochodzi do skutku w całości, albo wcale — schemat nie
// zostaje w stanie połowicznym.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"modernc.org/sqlite"
)

// Stałe ponawiania otwarcia transakcji: liczba prób, odstęp pierwszy i pułap odstępu.
const (
	probOtwarcia   = 5
	odstepOtwarcia = 25 * time.Millisecond
	odstepMaksimum = 400 * time.Millisecond
)

// kodyZajetosci to podstawowe kody SQLite oznaczające blokadę zapisu trzymaną
// przez inne połączenie: SQLITE_BUSY i SQLITE_LOCKED.
var kodyZajetosci = map[int]bool{5: true, 6: true}

// wTransakcji otwiera transakcję, wykonuje pracę i zatwierdza ją. Błąd pracy
// powoduje wycofanie i kończy wyłącznie bieżące wywołanie.
//
// Transakcja otwiera się poleceniem BEGIN IMMEDIATE — tak nastawia sterownik
// parametr `_txlock` z DSN złożonego w `store.zbudujDSN`. Blokada zapisu wzięta
// od razu zamienia zgubiony zapis na czekanie, ale przenosi zajętość bazy z ciszy
// na otwarcie transakcji, dlatego samo otwarcie jest tu ponawiane. Ponowienie
// obejmuje wyłącznie otwarcie: praca nie zdążyła się jeszcze wykonać ani raz,
// więc powtórzenie nie dubluje niczego.
func wTransakcji(ctx context.Context, db *sql.DB, praca func(transakcja *sql.Tx) error) error {
	transakcja, err := otworzTransakcje(ctx, db)
	if err != nil {
		return err
	}
	if err := praca(transakcja); err != nil {
		if wycofanie := transakcja.Rollback(); wycofanie != nil && wycofanie != sql.ErrTxDone {
			return fmt.Errorf("%w (wycofanie nieudane: %v)", err, wycofanie)
		}
		return err
	}
	if err := transakcja.Commit(); err != nil {
		return fmt.Errorf("dane: nie można zatwierdzić transakcji: %w", err)
	}
	return nil
}

// otworzTransakcje bierze blokadę zapisu, ponawiając próbę przy bazie zajętej
// przez inne połączenie. Odstęp rośnie dwukrotnie do pułapu, żeby kilku pisarzy
// naraz nie uderzało w bazę w tym samym takcie.
func otworzTransakcje(ctx context.Context, db *sql.DB) (*sql.Tx, error) {
	odstep := odstepOtwarcia
	var ostatni error
	for proba := 0; proba < probOtwarcia; proba++ {
		transakcja, err := db.BeginTx(ctx, nil)
		if err == nil {
			return transakcja, nil
		}
		ostatni = err
		if !bazaZajeta(err) {
			return nil, fmt.Errorf("dane: nie można otworzyć transakcji: %w", err)
		}
		if proba == probOtwarcia-1 {
			break
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("dane: nie można otworzyć transakcji: %w", ctx.Err())
		case <-time.After(odstep):
		}
		if odstep < odstepMaksimum {
			odstep *= 2
		}
	}
	return nil, fmt.Errorf("dane: baza zajęta po %d próbach otwarcia transakcji: %w", probOtwarcia, ostatni)
}

// bazaZajeta rozpoznaje odmowę wynikającą z blokady zapisu trzymanej przez inne
// połączenie — jedyną, którą wolno przeczekać.
func bazaZajeta(err error) bool {
	var bladSterownika *sqlite.Error
	if !errors.As(err, &bladSterownika) {
		return false
	}
	// SQLite dokłada do kodu podstawowego rozszerzenie na starszych bitach
	// (SQLITE_BUSY_SNAPSHOT to 517, czyli 5 z rozszerzeniem 2), stąd maska.
	return kodyZajetosci[bladSterownika.Code()&0xff]
}
