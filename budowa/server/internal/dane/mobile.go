// Odpowiedzialność pliku: dwa rachunki, których potrzebuje stan platformy oddawany mobilnemu centrum
// dowodzenia — liczba okien komunikacji otwartych i liczba kolejek czynnych.
package dane

import (
	"context"
	"fmt"
)

const (
	// policzOknaOtwarte liczy okna komunikacji, które nie zostały zamknięte.
	// Słownik kolumny jest dwuwartościowy, więc warunek wymienia wartość
	// obecną, a nie wyklucza wartości nieobecnych.
	policzOknaOtwarte = `SELECT COUNT(*) FROM okno_komunikacji WHERE stan = 'otwarte'`

	// policzKolejkiCzynne liczy kolejki poza stanem końcowym. Stany wymienione
	// wprost, bo CHECK schematu też je wymienia wprost — nowy stan kolejki ma
	// zatrzymać się na tym miejscu i wymusić decyzję, czy jest końcowy, zamiast
	// wpaść do rachunku milcząco.
	policzKolejkiCzynne = `SELECT COUNT(*) FROM kolejka
	                       WHERE stan IN ('bezczynna','pracuje','wstrzymana')
	                         AND ` + WarunekKonta
)

// LiczbaOtwartych zwraca liczbę okien komunikacji o stanie otwartym; zero jest wynikiem poprawnym, usterka
// odczytu wychodzi błędem.
func (r *repozytoriumOkien) LiczbaOtwartych(ctx context.Context) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, policzOknaOtwarte)
	if err != nil {
		return 0, err
	}
	var razem int
	if err := polecenie.QueryRowContext(ctx, KontoOperatora(ctx)).Scan(&razem); err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć okien otwartych: %w", err)
	}
	return razem, nil
}

// LiczbaCzynnych zwraca liczbę kolejek pozostających poza stanem końcowym ich cyklu pracy tej platformy.
func (r *repozytoriumKolejek) LiczbaCzynnych(ctx context.Context) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, policzKolejkiCzynne)
	if err != nil {
		return 0, err
	}
	var razem int
	if err := polecenie.QueryRowContext(ctx, KontoOperatora(ctx)).Scan(&razem); err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć kolejek czynnych: %w", err)
	}
	return razem, nil
}
