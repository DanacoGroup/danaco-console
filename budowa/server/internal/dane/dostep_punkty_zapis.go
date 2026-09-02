// Odpowiedzialność pliku: zapis katalogu punktów dostępu. Zapis dotyka trzech
// tabel (`punkt_dostepu`, `korzen_punktu_dostepu`, `argument_trybu_mostu`), więc
// idzie w jednej transakcji.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// ErrBrakUrzadzeniaBiezacego oznacza, że katalog lokalny nie ma na czym stanąć:
// żądanie nie wskazało urządzenia, a rdzeń nie ma w katalogu wiersza swojej maszyny.
var ErrBrakUrzadzeniaBiezacego = errors.New("dane: brak urządzenia bieżącego")

// ErrNieznaneUrzadzenie oznacza wskazanie urządzenia, którego w katalogu nie ma.
// Bez tego rozpoznania wskazanie rozbijałoby się dopiero o więz klucza obcego,
// czyli wracało jako awaria zapisu zamiast odmowy merytorycznej.
var ErrNieznaneUrzadzenie = errors.New("dane: nieznane urządzenie")

const (
	wstawPunktDostepu = `INSERT INTO punkt_dostepu
	                     (kod, nazwa, opis, rodzaj, urzadzenie_id, host, port, uzytkownik,
	                      sciezka_klucza, polecenie_startu, nazwa_mostu, poswiadczenie_odwolanie,
	                      tryb_domyslny, stan, aktywny, kolejnosc, konto_id)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ` +
		WskazanieKonta + `)`

	aktualizujPunktDostepu = `UPDATE punkt_dostepu
	                          SET nazwa = ?, opis = ?, rodzaj = ?, urzadzenie_id = ?, host = ?,
	                              port = ?, uzytkownik = ?, sciezka_klucza = ?,
	                              polecenie_startu = ?, nazwa_mostu = ?,
	                              poswiadczenie_odwolanie = ?, tryb_domyslny = ?, aktywny = ?,
	                              kolejnosc = ?,
	                              zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                          WHERE id = ? AND ` + WarunekKonta

	usunPunktDostepu = `DELETE FROM punkt_dostepu WHERE id = ? AND ` + WarunekKonta

	// Wskazanie urządzenia jest sprawdzane przed zapisem, bo więz klucza obcego
	// zgłasza pomyłkę dopiero jako awarię INSERT-a.
	istnieniUrzadzenia = `SELECT id FROM urzadzenie WHERE id = ?`

	// Maszyna, na której działa rdzeń. Oznaczenie `biezace` w schemacie zwalnia
	// warstwę proponującą urządzenie dla katalogu lokalnego z rozpoznawania
	// maszyny po raz drugi.
	urzadzenieBiezaceDlaPunktu = `SELECT id FROM urzadzenie WHERE biezace = 1`

	zapiszStanPunktu = `UPDATE punkt_dostepu
	                    SET stan = ?, sprawdzono = ?,
	                        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                    WHERE id = ? AND ` + WarunekKonta
)

// Dodaj wpisuje punkt dostępu do katalogu wraz z jego korzeniami i słownictwem trybu używanym przez most.
func (r *repozytoriumPunktowDostepu) Dodaj(ctx context.Context, punkt PunktDostepu) (int64, error) {
	wartosci, err := wartosciPunktu(punkt)
	if err != nil {
		return 0, err
	}
	var id int64
	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		urzadzenie, err := r.urzadzeniePunktu(ctx, transakcja, punkt)
		if err != nil {
			return err
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawPunktDostepu)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, punkt.Kod, punkt.Nazwa, punkt.Opis,
			wartosci.rodzaj, liczbaDoKolumny(urzadzenie), punkt.Host, portPunktu(punkt),
			punkt.Uzytkownik, punkt.SciezkaKlucza, punkt.PolecenieStartu, punkt.NazwaMostu,
			tekstDoKolumny(punkt.PoswiadczenieOdwolanie), wartosci.tryb, wartosci.stan,
			liczbaLogiczna(punkt.Aktywny), punkt.Kolejnosc, KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można dodać punktu dostępu %q: %w", punkt.Kod, err)
		}
		if id, err = wynik.LastInsertId(); err != nil {
			return fmt.Errorf("dane: nieznany identyfikator zapisanego punktu dostępu: %w", err)
		}
		return r.zapiszListy(ctx, transakcja, id, punkt)
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Aktualizuj zapisuje zmieniony punkt. Kod pozostaje stały — jest identyfikatorem
// punktu w nadaniach i w kontrakcie.
func (r *repozytoriumPunktowDostepu) Aktualizuj(ctx context.Context, punkt PunktDostepu) error {
	wartosci, err := wartosciPunktu(punkt)
	if err != nil {
		return err
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		urzadzenie, err := r.urzadzeniePunktu(ctx, transakcja, punkt)
		if err != nil {
			return err
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, aktualizujPunktDostepu)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, punkt.Nazwa, punkt.Opis, wartosci.rodzaj,
			liczbaDoKolumny(urzadzenie), punkt.Host, portPunktu(punkt), punkt.Uzytkownik,
			punkt.SciezkaKlucza, punkt.PolecenieStartu, punkt.NazwaMostu,
			tekstDoKolumny(punkt.PoswiadczenieOdwolanie), wartosci.tryb,
			liczbaLogiczna(punkt.Aktywny), punkt.Kolejnosc, punkt.ID, KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać punktu dostępu %d: %w", punkt.ID, err)
		}
		if err := sprawdzTrafienie(wynik, "punkt_dostepu", punkt.ID); err != nil {
			return err
		}
		return r.zapiszListy(ctx, transakcja, punkt.ID, punkt)
	})
}

// Usun kasuje punkt razem z jego korzeniami, słownictwem trybu i nadaniami —
// więzy klucza obcego kasują wiersze podrzędne kaskadą.
func (r *repozytoriumPunktowDostepu) Usun(ctx context.Context, id int64) error {
	polecenie, err := r.zapytania.przygotuj(ctx, usunPunktDostepu)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, id, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można usunąć punktu dostępu %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "punkt_dostepu", id)
}

// ZapiszWynikSprawdzenia odnotowuje wynik próby sięgnięcia do punktu. Stan
// nierozpoznany zostaje `unknown` i nie wygasza kontrolki.
func (r *repozytoriumPunktowDostepu) ZapiszWynikSprawdzenia(ctx context.Context, id int64,
	stan shared.AccessPointStatus, sprawdzono string) error {

	kolumna, err := stanPunktuNaBaze(stan)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszStanPunktu)
	if err != nil {
		return err
	}
	var czas any
	if sprawdzono != "" {
		czas = sprawdzono
	}
	wynik, err := polecenie.ExecContext(ctx, kolumna, czas, id, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać stanu punktu dostępu %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "punkt_dostepu", id)
}

// urzadzeniePunktu rozstrzyga, na które urządzenie wskazuje zapisywany punkt dostępu, zanim polecenie zapisu trafi do bazy.
func (r *repozytoriumPunktowDostepu) urzadzeniePunktu(ctx context.Context,
	transakcja *sql.Tx, punkt PunktDostepu) (*int64, error) {

	if punkt.UrzadzenieID != nil {
		id, err := r.jedenIdentyfikator(ctx, transakcja, istnieniUrzadzenia, *punkt.UrzadzenieID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("dane: urządzenie %d nie jest w katalogu: %w",
				*punkt.UrzadzenieID, ErrNieznaneUrzadzenie)
		}
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz urządzenia %d: %w",
				*punkt.UrzadzenieID, err)
		}
		return &id, nil
	}
	if punkt.Rodzaj != shared.AccessPointKindLocalDirectory {
		return nil, nil
	}
	id, err := r.jedenIdentyfikator(ctx, transakcja, urzadzenieBiezaceDlaPunktu)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBrakUrzadzeniaBiezacego
	}
	if err != nil {
		return nil, fmt.Errorf("dane: nieczytelny wiersz urządzenia bieżącego: %w", err)
	}
	return &id, nil
}

// jedenIdentyfikator odczytuje pojedynczy klucz wiersza w transakcji zapisu.
// Brak wiersza wraca nietknięty jako sql.ErrNoRows — dopiero wołający wie, czy
// oznacza pomyłkę wskazania, czy brak rozpoznania maszyny.
func (r *repozytoriumPunktowDostepu) jedenIdentyfikator(ctx context.Context, transakcja *sql.Tx,
	zapytanie string, argumenty ...any) (int64, error) {

	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, zapytanie)
	if err != nil {
		return 0, err
	}
	var id int64
	if err := polecenie.QueryRowContext(ctx, argumenty...).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// zapiszListy wymienia obie listy podrzędne punktu dostępu w ramach jednej transakcji jego zapisu do bazy.
func (r *repozytoriumPunktowDostepu) zapiszListy(ctx context.Context, transakcja *sql.Tx,
	punktID int64, punkt PunktDostepu) error {

	err := zapiszKorzenie(ctx, r.zapytania, transakcja, usunKorzeniePunktu, wstawKorzenPunktu,
		punktID, punkt.Korzenie, "punktu dostępu")
	if err != nil {
		return err
	}
	return zapiszArgumentyTrybu(ctx, r.zapytania, transakcja, punktID, punkt.ArgumentyTrybu)
}

// portPunktu pilnuje wartości domyślnej portu SSH — brak ustawienia nie może
// dawać portu zerowego, bo takiego połączenia nie da się nawiązać.
func portPunktu(punkt PunktDostepu) int {
	if punkt.Port > 0 {
		return punkt.Port
	}
	return domyslnyPortMostu
}

// domyslnyPortMostu odpowiada wartości domyślnej kolumny portu mostu ustawionej w schemacie bazy danych.
const domyslnyPortMostu = 22
