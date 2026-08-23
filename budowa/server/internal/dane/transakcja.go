// Odpowiedzialność pliku: wykonanie zapisu obejmującego wiele tabel w jednej
// transakcji. Zapis albo dochodzi do skutku w całości, albo wcale — schemat nie
// zostaje w stanie połowicznym.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// wTransakcji otwiera transakcję, wykonuje pracę i zatwierdza ją. Błąd pracy
// powoduje wycofanie i kończy wyłącznie bieżące wywołanie.
func wTransakcji(ctx context.Context, db *sql.DB, praca func(transakcja *sql.Tx) error) error {
	transakcja, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("dane: nie można otworzyć transakcji: %w", err)
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
