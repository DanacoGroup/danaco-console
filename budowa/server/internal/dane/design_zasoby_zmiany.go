// Odpowiedzialność pliku: zmiany zasobu Assets Panel jednym poleceniem SQL —
// oznaczenie ulubionego i usunięcie zasobu, każde w granicy konta.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *repozytoriumDesignu) UstawUlubionyZasobu(ctx context.Context,
	zasobID int64, ulubiony bool) error {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawUlubionyZasobuDesign)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, liczbaLogiczna(ulubiony), zasobID, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można przestawić ulubionego zasobu design %d: %w", zasobID, err)
	}
	return r.trafienieZasobuDesign(ctx, wynik, zastanyZasobDesignPoId, zasobID,
		fmt.Sprint(zasobID))
}

// UsunZasob oddaje, czy wiersz zniknął; brak wiersza nie jest odmową, wiersz cudzego konta jest.
func (r *repozytoriumDesignu) UsunZasob(ctx context.Context, kod string) (bool, error) {
	if kod == "" {
		return false, fmt.Errorf("dane: usunięcie zasobu design bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, usunZasobDesign)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć zasobu design %q: %w", kod, err)
	}
	err = r.trafienieZasobuDesign(ctx, wynik, zastanyZasobDesignPoKodzie, kod, kod)
	if errors.Is(err, ErrBrakWiersza) {
		return false, nil
	}
	return err == nil, err
}

// trafienieZasobuDesign odróżnia zero zmienionych wierszy przy wierszu cudzego konta
// (ErrKolizjaWiersza) od zera przy braku wiersza (ErrBrakWiersza).
func (r *repozytoriumDesignu) trafienieZasobuDesign(ctx context.Context, wynik sql.Result,
	zastany string, klucz any, wskazanie string) error {

	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nieznany skutek zapisu zasobu design %q: %w", wskazanie, err)
	}
	if zmienione > 0 {
		return nil
	}
	sonda, err := r.zapytania.przygotuj(ctx, zastany)
	if err != nil {
		return err
	}
	var jeden int
	err = sonda.QueryRowContext(ctx, klucz).Scan(&jeden)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("dane: zasób design %q nie istnieje: %w", wskazanie, ErrBrakWiersza)
	}
	if err != nil {
		return fmt.Errorf("dane: nie można sprawdzić zasobu design %q: %w", wskazanie, err)
	}
	return sprawdzTrafienieZapisu(wynik, "zasób design", wskazanie)
}
