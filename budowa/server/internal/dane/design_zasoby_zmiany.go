// Zmiany stanu zasobu Assets Panel, których nie da się wyrazić pełnym zapisem
// wiersza: oznaczenie ulubionego (`design.asset.favorite.set`) i usunięcie
// zasobu (`design.asset.remove`). Odczyt, zapis pełny i etykiety leżą
// w `design_zasoby.go`, kontrakt całego obszaru w `design.go`.
//
// Obie metody są jednym poleceniem, bez transakcji. `UstawEtykietyZasobu`
// prowadzi transakcję, bo składa się z dwóch poleceń (usuń zastane, wstaw
// nadesłane) i stan pośredni byłby widoczny; tutaj polecenie jest jedno,
// a SQLite wykonuje pojedynczy UPDATE lub DELETE niepodzielnie.
//
// Usunięcie nie rusza bajtów w magazynie. Blob leży pod sumą swojej zawartości
// (`adapter_modul_library_magazyn.go`), więc dwa zasoby o tej samej treści
// dzielą jeden plik, a skasowanie bloba przy usunięciu jednego z nich
// odebrałoby treść drugiemu. Rozstrzyga o tym warstwa rdzenia; tutaj znika sam
// wiersz.
package dane

import (
	"context"
	"fmt"
)

// UstawUlubionyZasobu przestawia oznaczenie ulubionego zasobu wskazanego
// kluczem wiersza; wołający ma go z `Zasob`, więc przekład kodu na wiersz
// i odróżnienie zasobu nieznanego zachodzą przed tym wywołaniem.
//
// Brak wiersza nie jest tu błędem. Sprawdzanie `RowsAffected` po raz drugi
// zamieniłoby wyścig — usunięcie zasobu w międzyczasie — w usterkę wewnętrzną
// zamiast w „nic się nie zmieniło”.
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

// UsunZasob usuwa zasób o wskazanym identyfikatorze zewnętrznym i mówi, czy
// jakikolwiek wiersz naprawdę zniknął.
//
// Liczba wierszy jest tu odpowiedzią, nie diagnostyką.
// `DesignAssetRemoveResponse.Removed` pyta, czy zasób istniał i został
// usunięty; usunięcie zera wierszy jest dla SQLite powodzeniem, więc bez
// `RowsAffected` komenda meldowałaby usunięcie zasobu, którego nie było.
// Powtórne usunięcie tego samego kodu oddaje `false`, a nie błąd.
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
	// Sterownik SQLite drzewa zna liczbę zmienionych wierszy; gdyby jej nie
	// podał, milczące „usunięto” byłoby zmyśleniem skutku — stąd błąd, nie
	// domyślne `true`.
	wierszy, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek usunięcia zasobu design %q: %w", kod, err)
	}
	return wierszy > 0, nil
}
