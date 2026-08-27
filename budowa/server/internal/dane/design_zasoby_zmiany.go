// Warstwa danych obsługuje zmiany zasobu Assets Panel niewyrażalne pełnym
// zapisem wiersza: oznaczenie ulubionego i usunięcie zasobu, każde jednym
// poleceniem SQL bez transakcji.
package dane

import (
	"context"
	"fmt"
)

// UstawUlubionyZasobu przestawia oznaczenie ulubionego zasobu wskazanego
// kluczem wiersza; brak wiersza do zmiany nie jest tu błędem.
func (r *repozytoriumDesignu) UstawUlubionyZasobu(ctx context.Context,
	zasobID int64, ulubiony bool) error {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawUlubionyZasobuDesign)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, liczbaLogiczna(ulubiony), zasobID); err != nil {
		return fmt.Errorf("dane: nie można przestawić ulubionego zasobu design %d: %w", zasobID, err)
	}
	return nil
}

// UsunZasob usuwa zasób o wskazanym identyfikatorze zewnętrznym i zwraca, czy
// jakikolwiek wiersz naprawdę zniknął.
func (r *repozytoriumDesignu) UsunZasob(ctx context.Context, kod string) (bool, error) {
	if kod == "" {
		return false, fmt.Errorf("dane: usunięcie zasobu design bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, usunZasobDesign)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć zasobu design %q: %w", kod, err)
	}
	// Sterownik SQLite zna liczbę zmienionych wierszy; brak jej wartości jest tu błędem.
	wierszy, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek usunięcia zasobu design %q: %w", kod, err)
	}
	return wierszy > 0, nil
}
