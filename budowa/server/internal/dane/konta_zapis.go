// Odpowiedzialność pliku: zapis katalogu kont — założenie konta, zmiana wiersza,
// odwołanie do poświadczenia i utrwalenie stanu rotacji.
//
// Odwołanie do poświadczenia ma w tym pakiecie dokładnie dwie drogi —
// `UstawPoswiadczenie` (wejście) i `OdwolaniePoswiadczenia` (wyjście dla warstwy,
// która musi je rozwiązać w magazynie sekretów). Żaden odczyt wykazu ani żadna
// odpowiedź kontraktu tą kolumną nie jedzie.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const (
	// znacznikZmiany podnosi datę ostatniej zmiany przy każdym zapisie wiersza.
	znacznikZmiany = `zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	nastepnaKolejnoscKonta = `SELECT COALESCE(MAX(kolejnosc), 0) + 1 FROM konto WHERE rodzaj = ?`

	wstawKonto = `INSERT INTO konto
	              (nazwa, rodzaj, dostawca, identyfikator_zewnetrzny, model_domyslny,
	               adres_bazowy, katalog_konfiguracji, poswiadczenie_odwolanie,
	               stan, aktywne, kolejnosc)
	              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// aktualizujKonto nie rusza poświadczenia, stanu rotacji ani oznaczenia konta
	// domyślnego — każde z nich ma własną, jawnie nazwaną drogę zapisu.
	aktualizujKonto = `UPDATE konto
	                   SET nazwa = ?, rodzaj = ?, dostawca = ?, identyfikator_zewnetrzny = ?,
	                       model_domyslny = ?, adres_bazowy = ?, katalog_konfiguracji = ?,
	                       aktywne = ?, kolejnosc = ?, ` + znacznikZmiany + `
	                   WHERE id = ?`

	zapiszPoswiadczenieKonta = `UPDATE konto SET poswiadczenie_odwolanie = ?, ` +
		znacznikZmiany + ` WHERE id = ?`

	odczytajPoswiadczenieKonta = `SELECT poswiadczenie_odwolanie FROM konto WHERE id = ?`

	oznaczStanKonta = `UPDATE konto SET stan = ?, wyczerpane_do = ?, ` +
		znacznikZmiany + ` WHERE id = ?`

	kanalyKonta = `SELECT id FROM kanal_modelu WHERE konto_id = ?`

	usunKonto = `DELETE FROM konto WHERE id = ?`
)

// Dodaj zakłada konto. Poświadczenie wchodzi osobnym argumentem, żeby nie dało
// się go zapisać przypadkiem razem z resztą wiersza. Kolejność niepodana (zero
// albo mniej) zostaje nadana jako następna w obrębie rodzaju — pula rotacji
// dostaje porządek bez pytania Operatora o liczbę.
func (r *repozytoriumKont) Dodaj(ctx context.Context, konto Konto,
	odwolaniePoswiadczenia *string) (int64, error) {
	rodzaj, stan, err := wartosciKonta(konto)
	if err != nil {
		return 0, err
	}
	var id int64
	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		kolejnosc, err := kolejnoscKonta(ctx, r, transakcja, konto, rodzaj)
		if err != nil {
			return err
		}
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawKonto)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, konto.Nazwa, rodzaj, konto.Dostawca,
			tekstDoKolumny(konto.IdentyfikatorZewnetrzny), tekstDoKolumny(konto.ModelDomyslny),
			tekstDoKolumny(konto.AdresBazowy), tekstDoKolumny(konto.KatalogKonfiguracji),
			tekstDoKolumny(odwolaniePoswiadczenia), stan, liczbaLogiczna(konto.Aktywne), kolejnosc)
		if err != nil {
			return fmt.Errorf("dane: nie można dodać konta %q: %w", konto.Nazwa, err)
		}
		id, err = wynik.LastInsertId()
		return err
	})
	return id, err
}

// Aktualizuj zapisuje zmieniony wiersz katalogu.
func (r *repozytoriumKont) Aktualizuj(ctx context.Context, konto Konto) error {
	rodzaj, _, err := wartosciKonta(konto)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, aktualizujKonto)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, konto.Nazwa, rodzaj, konto.Dostawca,
		tekstDoKolumny(konto.IdentyfikatorZewnetrzny), tekstDoKolumny(konto.ModelDomyslny),
		tekstDoKolumny(konto.AdresBazowy), tekstDoKolumny(konto.KatalogKonfiguracji),
		liczbaLogiczna(konto.Aktywne), konto.Kolejnosc, konto.ID)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać konta %d: %w", konto.ID, err)
	}
	return sprawdzTrafienie(wynik, "konto", konto.ID)
}

// UstawPoswiadczenie zapisuje odwołanie do danych dostępowych — nazwę wpisu
// w magazynie sekretów albo ścieżkę profilu. Wartość pusta kasuje odwołanie.
func (r *repozytoriumKont) UstawPoswiadczenie(ctx context.Context, id int64, odwolanie *string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPoswiadczenieKonta)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, tekstDoKolumny(odwolanie), id)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać odwołania poświadczenia konta %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "konto", id)
}

// OdwolaniePoswiadczenia zwraca odwołanie do poświadczenia. Wywołuje je
// wyłącznie warstwa rozwiązująca sekret przy uruchomieniu kanału; brak
// odwołania jest stanem normalnym i wraca jako ErrBrakWiersza.
func (r *repozytoriumKont) OdwolaniePoswiadczenia(ctx context.Context, id int64) (string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, odczytajPoswiadczenieKonta)
	if err != nil {
		return "", err
	}
	var odwolanie sql.NullString
	err = polecenie.QueryRowContext(ctx, id).Scan(&odwolanie)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%w: konto %d", ErrBrakWiersza, id)
	}
	if err != nil {
		return "", fmt.Errorf("dane: nie można odczytać odwołania poświadczenia konta %d: %w", id, err)
	}
	if !odwolanie.Valid || strings.TrimSpace(odwolanie.String) == "" {
		return "", fmt.Errorf("%w: poświadczenie konta %d", ErrBrakWiersza, id)
	}
	return odwolanie.String, nil
}

// OznaczStan utrwala zdatność konta w rotacji. `doChwili` niesie chwilę
// odnowienia limitu w zapisie ISO 8601; nil kasuje ją, bo stan inny niż
// wyczerpanie żadnej chwili nie ma.
func (r *repozytoriumKont) OznaczStan(ctx context.Context, id int64, stan StanKonta,
	doChwili *string) error {
	wartosc, err := stanKontaNaBaze(stan)
	if err != nil {
		return err
	}
	if wartosc != string(StanKontaWyczerpane) {
		doChwili = nil
	}
	polecenie, err := r.zapytania.przygotuj(ctx, oznaczStanKonta)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, wartosc, tekstDoKolumny(doChwili), id)
	if err != nil {
		return fmt.Errorf("dane: nie można oznaczyć stanu konta %d: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "konto", id)
}

// wartosciKonta przekłada pola wyliczeniowe wiersza na wartości kolumn.
func wartosciKonta(konto Konto) (rodzaj string, stan string, err error) {
	if rodzaj, err = rodzajKontaNaBaze(konto.Rodzaj); err != nil {
		return "", "", err
	}
	if stan, err = stanKontaNaBaze(konto.Stan); err != nil {
		return "", "", err
	}
	return rodzaj, stan, nil
}

// kolejnoscKonta zwraca kolejność wskazaną przez Operatora albo następną wolną
// w obrębie rodzaju.
func kolejnoscKonta(ctx context.Context, r *repozytoriumKont, transakcja *sql.Tx,
	konto Konto, rodzaj string) (int, error) {
	if konto.Kolejnosc > 0 {
		return konto.Kolejnosc, nil
	}
	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, nastepnaKolejnoscKonta)
	if err != nil {
		return 0, err
	}
	var kolejnosc int
	if err := polecenie.QueryRowContext(ctx, rodzaj).Scan(&kolejnosc); err != nil {
		return 0, fmt.Errorf("dane: nie można ustalić kolejności konta rodzaju %q: %w", rodzaj, err)
	}
	return kolejnosc, nil
}
