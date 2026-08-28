// Odpowiedzialność pliku: żywotność podagenta w bazie — przynależność jego
// pracy do konkretnego uruchomienia rdzenia i zamknięcie sierot po awarii albo restarcie.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// powodOsierocenia to powód zakończenia zapisywany sierocie — wartość ze słownika więzu sprawdzającego.
const powodOsierocenia = "osierocony"

// ZywotnoscPodagentow jest kontraktem żywotności podagenta: kto prowadzi jego
// pracę i co się z nią dzieje, gdy prowadzącego zabrakło. Wchodzi w skład
// `RepozytoriumPodagentow`, więc port podagentów rdzenia dostaje go tą samą
// metodą `Zestaw.Podagenci()`.
type ZywotnoscPodagentow interface {
	// OznaczProwadzenie przypisuje wskazanych podagentów uruchomieniu rdzenia, które podjęło ich pracę.
	OznaczProwadzenie(ctx context.Context, kody []string, uruchomienie string) error
	// ZamknijOsierocone zamyka pracę porzuconą przez uruchomienia wcześniejsze i oddaje wykaz zamkniętych.
	ZamknijOsierocone(ctx context.Context, uruchomienie, wyjasnienie string) ([]Podagent, error)
}

const (
	// Znacznik pusty łapie się w to samo sito co znacznik obcy: pracy nie prowadzi wtedy nikt w tym rdzeniu.
	osieroceniPodagenci = `SELECT ` + kolumnyPodagenta + zrodloPodagenta +
		` WHERE p.stan IN ('pending','running')
		    AND p.uruchomienie_rdzenia <> ?
		  ORDER BY p.id`

	zamknijOsieroconych = `UPDATE podagent
	                          SET stan = 'failed',
	                              powod_zakonczenia = '` + powodOsierocenia + `',
	                              wynik = COALESCE(wynik, ?),
	                              zakonczono = COALESCE(zakonczono,
	                                  strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                        WHERE stan IN ('pending','running')
	                          AND uruchomienie_rdzenia <> ?`

	oznaczProwadzeniePodagenta = `UPDATE podagent
	                                 SET uruchomienie_rdzenia = ?,
	                                     oznaka_zycia = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                               WHERE identyfikator_zewnetrzny = ?`
)

// OznaczProwadzenie zapisuje, które uruchomienie rdzenia prowadzi pracę wskazanych podagentów w jednej
// transakcji na całe powołanie.
func (r *repozytoriumPodagentow) OznaczProwadzenie(ctx context.Context,
	kody []string, uruchomienie string) error {

	if len(kody) == 0 {
		return nil
	}
	if uruchomienie == "" {
		return fmt.Errorf("dane: oznaczenie prowadzenia bez znacznika uruchomienia rdzenia")
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, oznaczProwadzeniePodagenta)
		if err != nil {
			return err
		}
		for _, kod := range kody {
			if kod == "" {
				return fmt.Errorf("dane: oznaczenie prowadzenia bez identyfikatora podagenta")
			}
			if _, err := polecenie.ExecContext(ctx, uruchomienie, kod); err != nil {
				return fmt.Errorf("dane: nie można oznaczyć prowadzenia podagenta %q: %w", kod, err)
			}
		}
		return nil
	})
}

// ZamknijOsierocone odnajduje i zamyka podagentów porzuconych przez uruchomienie rdzenia, którego już nie ma.
func (r *repozytoriumPodagentow) ZamknijOsierocone(ctx context.Context,
	uruchomienie, wyjasnienie string) ([]Podagent, error) {

	if uruchomienie == "" {
		return nil, fmt.Errorf("dane: zamykanie sierot bez znacznika uruchomienia rdzenia")
	}
	var osieroceni []Podagent
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, osieroceniPodagenci)
		if err != nil {
			return err
		}
		wiersze, err := odczyt.QueryContext(ctx, uruchomienie)
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać osieroconych podagentów: %w", err)
		}
		if osieroceni, err = zbierzPodagentow(wiersze); err != nil {
			return err
		}
		if len(osieroceni) == 0 {
			return nil
		}
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zamknijOsieroconych)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, wyjasnienie, uruchomienie); err != nil {
			return fmt.Errorf("dane: nie można zamknąć osieroconych podagentów: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return osieroceni, nil
}
