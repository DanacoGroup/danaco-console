// Odpowiedzialność pliku: powiązanie konta z kanałem modelu oraz usunięcie konta wraz z odłączeniem kanałów.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// Usun kasuje konto i zwraca identyfikatory kanałów, które utraciły powiązanie.
// Odczyt i skasowanie idą jedną transakcją, żeby wykaz odpowiadał stanowi
// sprzed skasowania.
func (r *repozytoriumKont) Usun(ctx context.Context, id int64) ([]int64, error) {
	odlaczone := []int64{}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		wykaz, err := kanalyKontaWTransakcji(ctx, r, transakcja, id)
		if err != nil {
			return err
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunKonto)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, id)
		if err != nil {
			return fmt.Errorf("dane: nie można usunąć konta %d: %w", id, err)
		}
		if err := sprawdzTrafienie(wynik, "konto", id); err != nil {
			return err
		}
		odlaczone = wykaz
		return nil
	})
	if err != nil {
		return nil, err
	}
	return odlaczone, nil
}

// kanalyKontaWTransakcji zwraca identyfikatory kanałów modelu wskazujących na dane konto w bazie danych.
func kanalyKontaWTransakcji(ctx context.Context, r *repozytoriumKont, transakcja *sql.Tx,
	id int64) ([]int64, error) {
	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, kanalyKonta)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kanałów konta %d: %w", id, err)
	}
	defer wiersze.Close()

	wykaz := []int64{}
	for wiersze.Next() {
		var kanal int64
		if err := wiersze.Scan(&kanal); err != nil {
			return nil, fmt.Errorf("dane: uszkodzony wiersz kanału konta %d: %w", id, err)
		}
		wykaz = append(wykaz, kanal)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kanałów konta %d: %w", id, err)
	}
	return wykaz, nil
}
