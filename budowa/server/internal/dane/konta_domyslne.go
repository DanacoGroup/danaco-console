// Odpowiedzialność pliku: wskazanie konta domyślnego. Domyślne jest dokładnie
// jedno na rodzaj — pilnuje tego indeks częściowy `idx_konto_domyslne_rodzaj`
// w schemacie, a nie warunek w kodzie. Repozytorium ma jedynie zdjąć
// oznaczenie z poprzedniego i nadać je nowemu w jednej transakcji, żeby indeks
// nigdy nie zobaczył dwóch kont domyślnych naraz.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	rodzajKontaPoId = `SELECT rodzaj FROM konto WHERE id = ?`

	poprzednieKontoDomyslne = `SELECT id FROM konto WHERE rodzaj = ? AND domyslne = 1`

	zdejmijOznaczenieDomyslnego = `UPDATE konto SET domyslne = 0, ` + znacznikZmiany + `
	                               WHERE rodzaj = ? AND domyslne = 1`

	nadajOznaczenieDomyslnego = `UPDATE konto SET domyslne = 1, ` + znacznikZmiany + `
	                             WHERE id = ?`
)

// UstawDomyslne czyni wskazane konto domyślnym w obrębie jego rodzaju i zwraca
// identyfikator konta, które oznaczenie utraciło. Brak poprzedniego konta
// domyślnego daje nil — nie jest to błąd.
func (r *repozytoriumKont) UstawDomyslne(ctx context.Context, id int64) (*int64, error) {
	var poprzednie *int64
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		rodzaj, err := rodzajKontaWTransakcji(ctx, r, transakcja, id)
		if err != nil {
			return err
		}
		poprzednie, err = poprzednieDomyslneWTransakcji(ctx, r, transakcja, rodzaj)
		if err != nil {
			return err
		}
		if err := wykonajWTransakcji(ctx, r, transakcja, zdejmijOznaczenieDomyslnego, rodzaj); err != nil {
			return fmt.Errorf("dane: nie można zdjąć oznaczenia konta domyślnego rodzaju %q: %w",
				rodzaj, err)
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, nadajOznaczenieDomyslnego)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, id)
		if err != nil {
			return fmt.Errorf("dane: nie można wskazać konta domyślnego %d: %w", id, err)
		}
		return sprawdzTrafienie(wynik, "konto", id)
	})
	if err != nil {
		return nil, err
	}
	if poprzednie != nil && *poprzednie == id {
		return nil, nil
	}
	return poprzednie, nil
}

// rodzajKontaWTransakcji odczytuje rodzaj konta w toczącej się transakcji.
func rodzajKontaWTransakcji(ctx context.Context, r *repozytoriumKont, transakcja *sql.Tx,
	id int64) (string, error) {
	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, rodzajKontaPoId)
	if err != nil {
		return "", err
	}
	var rodzaj string
	err = polecenie.QueryRowContext(ctx, id).Scan(&rodzaj)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%w: konto %d", ErrBrakWiersza, id)
	}
	if err != nil {
		return "", fmt.Errorf("dane: nie można odczytać rodzaju konta %d: %w", id, err)
	}
	return rodzaj, nil
}

// poprzednieDomyslneWTransakcji zwraca konto domyślne rodzaju albo nil.
func poprzednieDomyslneWTransakcji(ctx context.Context, r *repozytoriumKont, transakcja *sql.Tx,
	rodzaj string) (*int64, error) {
	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, poprzednieKontoDomyslne)
	if err != nil {
		return nil, err
	}
	var id int64
	err = polecenie.QueryRowContext(ctx, rodzaj).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać konta domyślnego rodzaju %q: %w", rodzaj, err)
	}
	return &id, nil
}

// wykonajWTransakcji uruchamia polecenie zapisu, którego liczba trafionych
// wierszy nie ma znaczenia — zdjęcie oznaczenia z pustego zbioru jest poprawne.
func wykonajWTransakcji(ctx context.Context, r *repozytoriumKont, transakcja *sql.Tx,
	zapytanie string, argumenty ...any) error {
	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, zapytanie)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, argumenty...)
	return err
}
