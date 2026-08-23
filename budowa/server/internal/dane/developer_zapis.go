// Odpowiedzialność pliku: zapis obszaru Developer — założenie migawki pliku,
// wpis przebiegu budowania, domknięcie przebiegu wynikiem oraz osierocenie
// przebiegów zostawionych przez poprzedni bieg rdzenia.
package dane

import (
	"context"
	"fmt"

	"danacoconsole/shared"
)

const (
	wstawWersjePliku = `INSERT INTO developer_wersja_pliku
	                    (kod, okno_kod, sciezka, tresc, rozmiar)
	                    VALUES (?, ?, ?, ?, ?)`

	wstawPrzebiegBudowania = `INSERT INTO developer_budowanie
	                          (kod, okno_kod, zadanie, argumenty, stan, kod_wyjscia, log, zakonczono)
	                          VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	                          ON CONFLICT(kod) DO UPDATE SET
	                            zadanie     = excluded.zadanie,
	                            argumenty   = excluded.argumenty,
	                            stan        = excluded.stan,
	                            kod_wyjscia = excluded.kod_wyjscia,
	                            log         = excluded.log,
	                            zakonczono  = excluded.zakonczono`

	// Domknięcie nie rusza wiersza już domkniętego: pierwszy prawdziwy kod
	// wyjścia nie ma prawa zostać nadpisany przez późniejsze przerwanie.
	domknijPrzebiegBudowania = `UPDATE developer_budowanie
	                            SET stan = ?, kod_wyjscia = ?, log = ?,
	                                zakonczono = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                            WHERE kod = ? AND stan = 'running'`

	osierocPrzebiegiBudowania = `UPDATE developer_budowanie
	                             SET stan = 'stopped',
	                                 zakonczono = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                             WHERE stan = 'running'`
)

// ZapiszWersje zakłada migawkę treści pliku. Migawka jest wpisem historii,
// więc powstaje zawsze nowym wierszem — nadpisanie poprzedniej migawki
// odebrałoby jej jedyny sens.
func (r *repozytoriumDevelopera) ZapiszWersje(ctx context.Context, wersja WersjaPliku) error {
	if wersja.Kod == "" || wersja.OknoKod == "" || wersja.Sciezka == "" {
		return fmt.Errorf("dane: wersja pliku bez identyfikatora, okna albo ścieżki")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawWersjePliku)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, wersja.Kod, wersja.OknoKod, wersja.Sciezka,
		wersja.Tresc, wersja.Rozmiar); err != nil {
		return fmt.Errorf("dane: nie można zapisać wersji pliku %q: %w", wersja.Sciezka, err)
	}
	return nil
}

// ZapiszPrzebieg wpisuje przebieg budowania albo odświeża jego stan.
func (r *repozytoriumDevelopera) ZapiszPrzebieg(ctx context.Context, przebieg PrzebiegBudowania) error {
	if przebieg.Kod == "" || przebieg.OknoKod == "" {
		return fmt.Errorf("dane: przebieg budowania bez identyfikatora przebiegu albo okna")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawPrzebiegBudowania)
	if err != nil {
		return err
	}
	stan := przebieg.Stan
	if stan == "" {
		stan = shared.BuildStatusRunning
	}
	if _, err := polecenie.ExecContext(ctx, przebieg.Kod, przebieg.OknoKod, przebieg.Zadanie,
		przebieg.Argumenty, stan, przebieg.KodWyjscia, przebieg.Log,
		przebieg.Zakonczono); err != nil {
		return fmt.Errorf("dane: nie można zapisać przebiegu budowania %q: %w", przebieg.Kod, err)
	}
	return nil
}

// ZakonczPrzebieg domyka wiersz przebiegu stanem końcowym, kodem wyjścia
// i ogonem logu.
func (r *repozytoriumDevelopera) ZakonczPrzebieg(ctx context.Context, kod string,
	stan shared.BuildStatus, kodWyjscia *int64, log string) error {

	if kod == "" {
		return fmt.Errorf("dane: domknięcie przebiegu budowania bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, domknijPrzebiegBudowania)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, stan, kodWyjscia, log, kod); err != nil {
		return fmt.Errorf("dane: nie można domknąć przebiegu budowania %q: %w", kod, err)
	}
	return nil
}

// OsierocPrzebiegi przestawia przebiegi poprzedniego biegu rdzenia na `stopped`.
func (r *repozytoriumDevelopera) OsierocPrzebiegi(ctx context.Context) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, osierocPrzebiegiBudowania)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można osierocić przebiegów budowania: %w", err)
	}
	// Osierocenie już się wykonało; nieudany odczyt liczby wierszy znaczy, że
	// nie wiadomo ilu — a nie że zeru.
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nie można odczytać liczby osieroconych przebiegów budowania: %w", err)
	}
	return zmienione, nil
}
