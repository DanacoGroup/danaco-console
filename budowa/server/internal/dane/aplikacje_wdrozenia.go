// Obszar Deployment modułu Apps: zapis i odczyt przebiegów wdrożenia; tabela niesie ślad zlecenia i jego stan.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// WdrozenieApp to jeden przebieg zlecenia `apps.deployment.run` wraz z jego stanem i wynikiem.
type WdrozenieApp struct {
	Kod            string
	OknoKod        string
	Srodowisko     shared.AppDeployEnvironment
	Strategia      shared.AppDeployStrategy
	Stan           shared.AppDeployStatus
	Wersja         *string
	NotatkiWydania *string
	Adres          *string
	LogOdwolanie   *string
	CofnieteDoKodu *string
	Rozpoczeto     string
	Zakonczono     *string
}

const (
	wstawWdrozenieApp = `INSERT INTO wdrozenie_apps
	                     (kod, okno, srodowisko, strategia, stan, wersja,
	                      notatki_wydania, adres, log_odwolanie, cofniete_do_kodu, zakonczono, konto_id)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                     ON CONFLICT(kod) DO UPDATE SET
	                       srodowisko       = excluded.srodowisko,
	                       strategia        = excluded.strategia,
	                       stan             = excluded.stan,
	                       wersja           = excluded.wersja,
	                       notatki_wydania  = excluded.notatki_wydania,
	                       adres            = excluded.adres,
	                       log_odwolanie    = excluded.log_odwolanie,
	                       cofniete_do_kodu = excluded.cofniete_do_kodu,
	                       zakonczono       = excluded.zakonczono
	                     WHERE ` + WarunekKonta + `
	                     RETURNING kod, okno, srodowisko, strategia, stan, wersja,
	                       notatki_wydania, adres, log_odwolanie, cofniete_do_kodu,
	                       rozpoczeto, zakonczono`

	kolumnyWdrozeniaApp = `kod, okno, srodowisko, strategia, stan, wersja,
	                       notatki_wydania, adres, log_odwolanie, cofniete_do_kodu,
	                       rozpoczeto, zakonczono`

	pobierzWdrozenieApp = `SELECT ` + kolumnyWdrozeniaApp + `
	                       FROM wdrozenie_apps
	                       WHERE kod = ? AND ` + WarunekKonta

	// Zawężenie do środowiska pustym łańcuchem znaczy „bez zawężenia" — jedno zapytanie zamiast dwóch sklejanych.
	warunekWdrozenApp = ` WHERE okno = ? AND (? = '' OR srodowisko = ?) AND ` + WarunekKonta

	pobierzWdrozeniaApp = `SELECT ` + kolumnyWdrozeniaApp + `
	                       FROM wdrozenie_apps` + warunekWdrozenApp + `
	                       ORDER BY rozpoczeto DESC, id DESC
	                       LIMIT CASE WHEN ? > 0 THEN ? ELSE -1 END`

	// Bez LIMIT-u: `total` kontraktu opisuje rozmiar historii, nie rozmiar strony.
	policzWdrozeniaApp = `SELECT COUNT(*) FROM wdrozenie_apps` + warunekWdrozenApp
)

// ZapiszWdrozenie zakłada wiersz przebiegu wdrożenia albo odświeża jego stan, gdy kod przebiegu już istnieje.
func (r *repozytoriumAplikacji) ZapiszWdrozenie(ctx context.Context,
	wdrozenie WdrozenieApp) (WdrozenieApp, error) {

	if wdrozenie.Kod == "" || wdrozenie.OknoKod == "" {
		return WdrozenieApp{}, fmt.Errorf("dane: wdrożenie bez identyfikatora albo okna")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawWdrozenieApp)
	if err != nil {
		return WdrozenieApp{}, err
	}
	stan := wdrozenie.Stan
	if stan == "" {
		stan = shared.AppDeployStatusPending
	}
	wiersz := polecenie.QueryRowContext(ctx, wdrozenie.Kod, wdrozenie.OknoKod, wdrozenie.Srodowisko,
		wdrozenie.Strategia, stan, wdrozenie.Wersja, wdrozenie.NotatkiWydania, wdrozenie.Adres,
		wdrozenie.LogOdwolanie, wdrozenie.CofnieteDoKodu, wdrozenie.Zakonczono,
		KontoOperatora(ctx), KontoOperatora(ctx))
	zapisane, err := odczytajWdrozenieApp(wiersz)
	if err != nil {
		return WdrozenieApp{}, fmt.Errorf("dane: nie można zapisać wdrożenia %q: %w", wdrozenie.Kod, err)
	}
	return zapisane, nil
}

// Wdrozenie zwraca przebieg wdrożenia po kodzie; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumAplikacji) Wdrozenie(ctx context.Context, kod string) (WdrozenieApp, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWdrozenieApp)
	if err != nil {
		return WdrozenieApp{}, err
	}
	wdrozenie, err := odczytajWdrozenieApp(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return WdrozenieApp{}, ErrBrakWiersza
	}
	if err != nil {
		return WdrozenieApp{}, fmt.Errorf("dane: nieczytelne wdrożenie %q: %w", kod, err)
	}
	return wdrozenie, nil
}

// Wdrozenia zwraca przebiegi wdrożenia okna od najnowszego wraz z liczbą wszystkich spełniających warunki; limit niedodatni znaczy wykaz pełny.
func (r *repozytoriumAplikacji) Wdrozenia(ctx context.Context, okno string,
	srodowisko *shared.AppDeployEnvironment, limit int) ([]WdrozenieApp, int, error) {

	zawezenie := ""
	if srodowisko != nil {
		zawezenie = string(*srodowisko)
	}

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWdrozeniaApp)
	if err != nil {
		return nil, 0, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, zawezenie, zawezenie, KontoOperatora(ctx), limit, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać wdrożeń okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	wdrozenia := make([]WdrozenieApp, 0, 8)
	for wiersze.Next() {
		wdrozenie, err := odczytajWdrozenieApp(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz wdrożenia: %w", err)
		}
		wdrozenia = append(wdrozenia, wdrozenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt wdrożeń okna %q: %w", okno, err)
	}

	licznik, err := r.zapytania.przygotuj(ctx, policzWdrozeniaApp)
	if err != nil {
		return nil, 0, err
	}
	razem := 0
	if err := licznik.QueryRowContext(ctx, okno, zawezenie, zawezenie, KontoOperatora(ctx)).Scan(&razem); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć wdrożeń okna %q: %w", okno, err)
	}
	return wdrozenia, razem, nil
}

func odczytajWdrozenieApp(wiersz interface{ Scan(...any) error }) (WdrozenieApp, error) {
	var wdrozenie WdrozenieApp
	err := wiersz.Scan(&wdrozenie.Kod, &wdrozenie.OknoKod, &wdrozenie.Srodowisko, &wdrozenie.Strategia,
		&wdrozenie.Stan, &wdrozenie.Wersja, &wdrozenie.NotatkiWydania, &wdrozenie.Adres,
		&wdrozenie.LogOdwolanie, &wdrozenie.CofnieteDoKodu, &wdrozenie.Rozpoczeto, &wdrozenie.Zakonczono)
	return wdrozenie, err
}
