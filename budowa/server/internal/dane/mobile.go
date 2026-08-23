// Odpowiedzialność pliku: dwa rachunki, których potrzebuje stan platformy
// oddawany mobilnemu centrum dowodzenia (`mobile.status.get`) — liczba okien
// komunikacji otwartych i liczba kolejek czynnych.
//
// Warstwa mobilna nie ma własnej tabeli. Proces mobilny jest wpisem rejestru
// telemetrii postępu trzymanego w pamięci; sesje, okna, procesy i kolejki mają
// swoich właścicieli, a ten plik dokłada wyłącznie odczyt.
//
// Metody siedzą na repozytoriach okien i kolejek, bo tabela ma jednego
// właściciela. Osobne repozytorium mobilne z własnymi zdaniami SELECT nad
// `okno_komunikacji` i `kolejka` byłoby drugim czytelnikiem cudzych tabel
// i rozjechałoby się z właścicielem przy pierwszej zmianie słownika stanów.
// Osobny jest wyłącznie plik, żeby widać było, po co te rachunki powstały.
//
// Czynność bytu mierzy się tu stanem spoza stanów końcowych: sesja czynna to
// `czynna` albo `wstrzymana`, nigdy `zakonczona` ani `archiwalna`
// (`adapter_nawigacja.go`). Kolejka idzie tą samą miarą: czynna jest
// `bezczynna`, `pracuje` i `wstrzymana`, końcowe są `zatrzymana` i
// `wyczerpana`. Okno komunikacji zna dwa stany, więc liczą się wiersze
// `otwarte`.
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
	                       WHERE stan IN ('bezczynna','pracuje','wstrzymana')`
)

// LiczbaOtwartych zwraca liczbę okien komunikacji o stanie `otwarte`.
//
// Zero jest wynikiem poprawnym: platforma bez otwartego okna pracuje dalej.
// Usterka odczytu wychodzi błędem, a nie zerem — zero policzone i zero
// nieprzeczytane to dwie różne odpowiedzi.
func (r *repozytoriumOkien) LiczbaOtwartych(ctx context.Context) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, policzOknaOtwarte)
	if err != nil {
		return 0, err
	}
	var razem int
	if err := polecenie.QueryRowContext(ctx).Scan(&razem); err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć okien otwartych: %w", err)
	}
	return razem, nil
}

// LiczbaCzynnych zwraca liczbę kolejek poza stanem końcowym.
func (r *repozytoriumKolejek) LiczbaCzynnych(ctx context.Context) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, policzKolejkiCzynne)
	if err != nil {
		return 0, err
	}
	var razem int
	if err := polecenie.QueryRowContext(ctx).Scan(&razem); err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć kolejek czynnych: %w", err)
	}
	return razem, nil
}
