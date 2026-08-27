// Plik zapisuje katalog urządzeń: założenie, aktualizację i rozpoznanie
// maszyny bieżącej w tabeli urzadzenie.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	wstawUrzadzenie = `INSERT INTO urzadzenie
	                   (nazwa, identyfikator_sprzetowy, nazwa_hosta, system_operacyjny,
	                    wersja_klienta, zaufane, ostatnio_widziane)
	                   VALUES (?, ?, ?, ?, ?, ?, ?)`

	aktualizujUrzadzenie = `UPDATE urzadzenie
	                        SET nazwa = ?, nazwa_hosta = ?, system_operacyjny = ?,
	                            wersja_klienta = ?, zaufane = ?, ostatnio_widziane = ?
	                        WHERE id = ?`

	zdejmijOznaczenieBiezacego = `UPDATE urzadzenie SET biezace = 0
	                              WHERE biezace = 1 AND identyfikator_sprzetowy <> ?`

	// zapiszUrzadzenieBiezace wstawia albo odświeża wiersz maszyny bieżącej
	// jednym poleceniem, zostawiając kolumny nazwa i zaufane nietknięte przy
	// odświeżeniu.
	zapiszUrzadzenieBiezace = `INSERT INTO urzadzenie
	                           (nazwa, identyfikator_sprzetowy, nazwa_hosta, system_operacyjny,
	                            wersja_klienta, zaufane, biezace, ostatnio_widziane)
	                           VALUES (?, ?, ?, ?, ?, ?, 1,
	                                   strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                           ON CONFLICT(identyfikator_sprzetowy) DO UPDATE SET
	                               nazwa_hosta = excluded.nazwa_hosta,
	                               system_operacyjny = excluded.system_operacyjny,
	                               ostatnio_widziane = excluded.ostatnio_widziane,
	                               biezace = 1`
)

// Dodaj wpisuje urządzenie do katalogu i zwraca jego klucz wiersza — ten sam,
// którym punkt dostępu rodzaju `localDirectory` wskazuje maszynę.
func (r *repozytoriumUrzadzen) Dodaj(ctx context.Context, urzadzenie Urzadzenie) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, wstawUrzadzenie)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, nazwaUrzadzenia(urzadzenie),
		urzadzenie.IdentyfikatorSprzetowy, urzadzenie.NazwaHosta, urzadzenie.SystemOperacyjny,
		urzadzenie.WersjaKlienta, liczbaLogiczna(urzadzenie.Zaufane),
		tekstDoKolumny(urzadzenie.OstatnioWidziane))
	if err != nil {
		return 0, fmt.Errorf("dane: nie można dodać urządzenia %q: %w",
			urzadzenie.IdentyfikatorSprzetowy, err)
	}
	id, err := wynik.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("dane: nieznany identyfikator zapisanego urządzenia: %w", err)
	}
	return id, nil
}

// Aktualizuj zapisuje zmienione urządzenie. Identyfikator sprzętowy pozostaje
// stały — jest rozpoznaniem maszyny, a nie polem do edycji.
func (r *repozytoriumUrzadzen) Aktualizuj(ctx context.Context, urzadzenie Urzadzenie) error {
	polecenie, err := r.zapytania.przygotuj(ctx, aktualizujUrzadzenie)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, nazwaUrzadzenia(urzadzenie), urzadzenie.NazwaHosta,
		urzadzenie.SystemOperacyjny, urzadzenie.WersjaKlienta,
		liczbaLogiczna(urzadzenie.Zaufane), tekstDoKolumny(urzadzenie.OstatnioWidziane),
		urzadzenie.ID)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać urządzenia %d: %w", urzadzenie.ID, err)
	}
	return sprawdzTrafienie(wynik, "urzadzenie", urzadzenie.ID)
}

// ZapewnijBiezace zakłada albo odświeża wiersz maszyny, na której działa rdzeń.
// Obie zmiany idą w jednej transakcji: między zdjęciem oznaczenia z poprzedniej
// maszyny a nadaniem go bieżącej baza nie może zostać bez maszyny bieżącej ani
// z dwiema naraz.
func (r *repozytoriumUrzadzen) ZapewnijBiezace(ctx context.Context,
	urzadzenie Urzadzenie) (Urzadzenie, error) {

	if urzadzenie.IdentyfikatorSprzetowy == "" {
		return Urzadzenie{}, fmt.Errorf("dane: urządzenie bieżące bez identyfikatora sprzętowego")
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zdjecie, err := r.zapytania.wTransakcji(ctx, transakcja, zdejmijOznaczenieBiezacego)
		if err != nil {
			return err
		}
		if _, err := zdjecie.ExecContext(ctx, urzadzenie.IdentyfikatorSprzetowy); err != nil {
			return fmt.Errorf("dane: nie można zdjąć oznaczenia maszyny bieżącej: %w", err)
		}
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszUrzadzenieBiezace)
		if err != nil {
			return err
		}
		_, err = zapis.ExecContext(ctx, nazwaUrzadzenia(urzadzenie),
			urzadzenie.IdentyfikatorSprzetowy, urzadzenie.NazwaHosta,
			urzadzenie.SystemOperacyjny, urzadzenie.WersjaKlienta,
			liczbaLogiczna(urzadzenie.Zaufane))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać urządzenia bieżącego %q: %w",
				urzadzenie.IdentyfikatorSprzetowy, err)
		}
		return nil
	})
	if err != nil {
		return Urzadzenie{}, err
	}
	return r.PoIdentyfikatorze(ctx, urzadzenie.IdentyfikatorSprzetowy)
}

// nazwaUrzadzenia pilnuje kolumny NOT NULL `nazwa`. Brak napisu do pokazania nie
// może zablokować zapisu maszyny — nazwą zastępczą jest wtedy nazwa hosta, a gdy
// i jej brak, identyfikator sprzętowy.
func nazwaUrzadzenia(urzadzenie Urzadzenie) string {
	if urzadzenie.Nazwa != "" {
		return urzadzenie.Nazwa
	}
	if urzadzenie.NazwaHosta != "" {
		return urzadzenie.NazwaHosta
	}
	return urzadzenie.IdentyfikatorSprzetowy
}
