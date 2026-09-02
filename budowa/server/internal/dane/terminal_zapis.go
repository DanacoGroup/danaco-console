// Zapis obszaru Terminal: karta powłoki, wpis procesu, domknięcie procesu i osierocenie po biegu rdzenia.
package dane

import (
	"context"
	"fmt"

	"danacoconsole/shared"
)

const (
	wstawKarteTerminala = `INSERT INTO terminal_karta
	                       (kod, okno_kod, powloka, tytul, katalog_roboczy, stan,
	                        cel_zdalny, port_zdalny, host_kod, konto_id)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                       ON CONFLICT(kod) DO UPDATE SET
	                         okno_kod        = excluded.okno_kod,
	                         powloka         = excluded.powloka,
	                         tytul           = excluded.tytul,
	                         katalog_roboczy = excluded.katalog_roboczy,
	                         stan            = excluded.stan,
	                         cel_zdalny      = excluded.cel_zdalny,
	                         port_zdalny     = excluded.port_zdalny,
	                         host_kod        = excluded.host_kod
	                       WHERE ` + WarunekKonta

	wstawProcesTerminala = `INSERT INTO terminal_proces
	                        (kod, karta_kod, okno_kod, pid, pid_nadrzedny, polecenie,
	                         inicjator, stan, kod_wyjscia, zakonczono, konto_id)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                        ON CONFLICT(kod) DO UPDATE SET
	                          pid           = excluded.pid,
	                          pid_nadrzedny = excluded.pid_nadrzedny,
	                          stan          = excluded.stan,
	                          kod_wyjscia   = excluded.kod_wyjscia,
	                          zakonczono    = excluded.zakonczono
	                        WHERE ` + WarunekKonta

	domknijProcesTerminala = `UPDATE terminal_proces
	                          SET stan = ?, kod_wyjscia = ?,
	                              zakonczono = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                          WHERE kod = ? AND stan = 'running' AND ` + WarunekKonta

	osierocProcesyTerminala = `UPDATE terminal_proces
	                           SET stan = 'stopped',
	                               zakonczono = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                           WHERE stan = 'running'`
)

// ZapiszKarte zakłada kartę powłoki terminala albo, przy istniejącym kodzie, odświeża jej profil i stan.
func (r *repozytoriumTerminala) ZapiszKarte(ctx context.Context, karta KartaTerminala) error {
	if karta.Kod == "" || karta.OknoKod == "" {
		return fmt.Errorf("dane: karta terminala bez identyfikatora karty albo okna")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawKarteTerminala)
	if err != nil {
		return err
	}
	stan := karta.Stan
	if stan == "" {
		stan = shared.TerminalSessionStatusRunning
	}
	if _, err := polecenie.ExecContext(ctx, karta.Kod, karta.OknoKod, karta.Powloka,
		karta.Tytul, karta.KatalogRoboczy, stan,
		karta.CelZdalny, karta.PortZdalny, karta.HostKod,
		KontoOperatora(ctx), KontoOperatora(ctx)); err != nil {
		return fmt.Errorf("dane: nie można zapisać karty terminala %q: %w", karta.Kod, err)
	}
	return nil
}

// ZapiszProces wpisuje proces terminala albo, przy istniejącym kodzie, odświeża jego stan i kod wyjścia.
func (r *repozytoriumTerminala) ZapiszProces(ctx context.Context, proces ProcesTerminala) error {
	if proces.Kod == "" || proces.OknoKod == "" {
		return fmt.Errorf("dane: proces terminala bez identyfikatora procesu albo okna")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawProcesTerminala)
	if err != nil {
		return err
	}
	stan := proces.Stan
	if stan == "" {
		stan = shared.TerminalProcessStatusRunning
	}
	inicjator := proces.Inicjator
	if inicjator == "" {
		inicjator = shared.ProcessInitiatorOperator
	}
	if _, err := polecenie.ExecContext(ctx, proces.Kod, proces.KartaKod, proces.OknoKod,
		proces.Pid, proces.PidNadrzedny, proces.Polecenie, inicjator, stan,
		proces.KodWyjscia, proces.Zakonczono, KontoOperatora(ctx), KontoOperatora(ctx)); err != nil {
		return fmt.Errorf("dane: nie można zapisać procesu %q: %w", proces.Kod, err)
	}
	return nil
}

// ZakonczProces domyka proces terminala stanem końcowym i kodem wyjścia, po jego identyfikatorze.
func (r *repozytoriumTerminala) ZakonczProces(ctx context.Context, kod string,
	stan shared.TerminalProcessStatus, kodWyjscia *int64) error {

	if kod == "" {
		return fmt.Errorf("dane: domknięcie procesu bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, domknijProcesTerminala)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, stan, kodWyjscia, kod, KontoOperatora(ctx)); err != nil {
		return fmt.Errorf("dane: nie można domknąć procesu %q: %w", kod, err)
	}
	return nil
}

// OsierociProcesy przestawia na stopped procesy pozostałe z poprzedniego biegu rdzenia i zwraca ich liczbę.
func (r *repozytoriumTerminala) OsierociProcesy(ctx context.Context) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, osierocProcesyTerminala)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można osierocić procesów terminala: %w", err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return zmienione, nil
}
