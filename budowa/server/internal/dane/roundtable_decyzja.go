// Odpowiedzialność pliku: macierz decyzyjna wariantów i kryteriów
// (`store/migracja_197_roundtable_macierz.sql`).
//
// Wyniku ważonego tu nie ma: liczy go rdzeń przy odczycie. Zmiana wagi jednego
// kryterium przestawia wynik każdego wariantu naraz, więc kolumna z wynikiem
// rozjechałaby się z ocenami przy pierwszym pominięciu przeliczenia.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// MacierzDebaty to macierz decyzyjna okna wraz z kryteriami i wariantami.
type MacierzDebaty struct {
	Kod            string
	Okno           string
	Nazwa          string
	Zaktualizowano string
	Kryteria       []KryteriumMacierzyDebaty
	Warianty       []WariantMacierzyDebaty
}

// KryteriumMacierzyDebaty to jedno kryterium wraz z wagą.
type KryteriumMacierzyDebaty struct {
	Kod       string
	Macierz   string
	Nazwa     string
	Waga      float64
	Kolejnosc int
}

// WariantMacierzyDebaty to wariant wraz z ocenami w kryteriach, zapisanymi jako
// mapa „kryterium → ocena".
type WariantMacierzyDebaty struct {
	Kod       string
	Macierz   string
	Etykieta  string
	OcenyJson string
	Kolejnosc int
}

// RepozytoriumDebatyDecyzji jest częścią kontraktu obszaru odpowiadającą za
// macierz decyzyjną.
type RepozytoriumDebatyDecyzji interface {
	ZapiszMacierzDebaty(ctx context.Context, macierz MacierzDebaty) (MacierzDebaty, error)
	MacierzDebatyPoKodzie(ctx context.Context, kod string) (MacierzDebaty, error)
	OstatniaMacierzDebaty(ctx context.Context, okno string) (MacierzDebaty, error)
}

const (
	zapiszMacierzDebaty = `INSERT INTO debata_macierz (identyfikator_zewnetrzny, okno, nazwa)
	                       VALUES (?, ?, ?)
	                       ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                           nazwa = excluded.nazwa,
	                           zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	usunKryteriaMacierzyDebaty = `DELETE FROM debata_macierz_kryterium WHERE macierz = ?`
	usunWariantyMacierzyDebaty = `DELETE FROM debata_macierz_wariant WHERE macierz = ?`

	zapiszKryteriumMacierzyDebaty = `INSERT INTO debata_macierz_kryterium
	                                 (identyfikator_zewnetrzny, macierz, nazwa, waga, kolejnosc)
	                                 VALUES (?, ?, ?, ?, ?)`

	zapiszWariantMacierzyDebaty = `INSERT INTO debata_macierz_wariant
	                               (identyfikator_zewnetrzny, macierz, etykieta, oceny_json, kolejnosc)
	                               VALUES (?, ?, ?, ?, ?)`

	pobierzMacierzDebaty = `SELECT identyfikator_zewnetrzny, okno, nazwa, zaktualizowano
	                        FROM debata_macierz WHERE identyfikator_zewnetrzny = ?`

	pobierzOstatniaMacierzDebaty = `SELECT identyfikator_zewnetrzny, okno, nazwa, zaktualizowano
	                                FROM debata_macierz WHERE okno = ? ORDER BY id DESC LIMIT 1`

	pobierzKryteriaMacierzyDebaty = `SELECT identyfikator_zewnetrzny, macierz, nazwa, waga, kolejnosc
	                                 FROM debata_macierz_kryterium WHERE macierz = ?
	                                 ORDER BY kolejnosc ASC, id ASC`

	pobierzWariantyMacierzyDebaty = `SELECT identyfikator_zewnetrzny, macierz, etykieta, oceny_json,
	                                        kolejnosc
	                                 FROM debata_macierz_wariant WHERE macierz = ?
	                                 ORDER BY kolejnosc ASC, id ASC`
)

// ZapiszMacierzDebaty zakłada albo zmienia macierz wraz z całą zawartością
// w jednej transakcji: macierz z nowymi kryteriami i starymi wariantami byłaby
// tabelą, w której kolumny nie odpowiadają wierszom.
func (r *repozytoriumRoundtable) ZapiszMacierzDebaty(ctx context.Context,
	macierz MacierzDebaty) (MacierzDebaty, error) {

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		naglowek, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszMacierzDebaty)
		if err != nil {
			return err
		}
		if _, err := naglowek.ExecContext(ctx, macierz.Kod, macierz.Okno, macierz.Nazwa); err != nil {
			return fmt.Errorf("dane: nie można zapisać macierzy decyzyjnej %q: %w", macierz.Kod, err)
		}
		for _, polecenieUsuwajace := range []string{usunKryteriaMacierzyDebaty, usunWariantyMacierzyDebaty} {
			wyczysc, err := r.zapytania.wTransakcji(ctx, transakcja, polecenieUsuwajace)
			if err != nil {
				return err
			}
			if _, err := wyczysc.ExecContext(ctx, macierz.Kod); err != nil {
				return fmt.Errorf("dane: nie można wyczyścić macierzy decyzyjnej %q: %w",
					macierz.Kod, err)
			}
		}
		wstawKryterium, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszKryteriumMacierzyDebaty)
		if err != nil {
			return err
		}
		for pozycja, kryterium := range macierz.Kryteria {
			if _, err := wstawKryterium.ExecContext(ctx, kryterium.Kod, macierz.Kod, kryterium.Nazwa,
				kryterium.Waga, pozycja+1); err != nil {
				return fmt.Errorf("dane: nie można zapisać kryterium macierzy %q: %w",
					kryterium.Kod, err)
			}
		}
		wstawWariant, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszWariantMacierzyDebaty)
		if err != nil {
			return err
		}
		for pozycja, wariant := range macierz.Warianty {
			if _, err := wstawWariant.ExecContext(ctx, wariant.Kod, macierz.Kod, wariant.Etykieta,
				wariant.OcenyJson, pozycja+1); err != nil {
				return fmt.Errorf("dane: nie można zapisać wariantu macierzy %q: %w", wariant.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return MacierzDebaty{}, err
	}
	return r.MacierzDebatyPoKodzie(ctx, macierz.Kod)
}

// MacierzDebatyPoKodzie zwraca macierz wraz z kryteriami i wariantami.
func (r *repozytoriumRoundtable) MacierzDebatyPoKodzie(ctx context.Context,
	kod string) (MacierzDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzMacierzDebaty)
	if err != nil {
		return MacierzDebaty{}, err
	}
	return r.zlozMacierzDebaty(ctx, polecenie.QueryRowContext(ctx, kod))
}

// OstatniaMacierzDebaty zwraca ostatnio założoną macierz okna.
func (r *repozytoriumRoundtable) OstatniaMacierzDebaty(ctx context.Context,
	okno string) (MacierzDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzOstatniaMacierzDebaty)
	if err != nil {
		return MacierzDebaty{}, err
	}
	return r.zlozMacierzDebaty(ctx, polecenie.QueryRowContext(ctx, okno))
}

// zlozMacierzDebaty czyta nagłówek macierzy i dobiera do niego zawartość.
func (r *repozytoriumRoundtable) zlozMacierzDebaty(ctx context.Context,
	wiersz *sql.Row) (MacierzDebaty, error) {

	var macierz MacierzDebaty
	err := wiersz.Scan(&macierz.Kod, &macierz.Okno, &macierz.Nazwa, &macierz.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return MacierzDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return MacierzDebaty{}, fmt.Errorf("dane: nieczytelny wiersz macierzy decyzyjnej: %w", err)
	}

	kryteria, err := r.zapytania.przygotuj(ctx, pobierzKryteriaMacierzyDebaty)
	if err != nil {
		return MacierzDebaty{}, err
	}
	wierszeKryteriow, err := kryteria.QueryContext(ctx, macierz.Kod)
	if err != nil {
		return MacierzDebaty{}, fmt.Errorf("dane: nie można odczytać kryteriów macierzy %q: %w",
			macierz.Kod, err)
	}
	defer wierszeKryteriow.Close()
	for wierszeKryteriow.Next() {
		var kryterium KryteriumMacierzyDebaty
		if err := wierszeKryteriow.Scan(&kryterium.Kod, &kryterium.Macierz, &kryterium.Nazwa,
			&kryterium.Waga, &kryterium.Kolejnosc); err != nil {
			return MacierzDebaty{}, fmt.Errorf("dane: nieczytelny wiersz kryterium macierzy: %w", err)
		}
		macierz.Kryteria = append(macierz.Kryteria, kryterium)
	}
	if err := wierszeKryteriow.Err(); err != nil {
		return MacierzDebaty{}, err
	}

	warianty, err := r.zapytania.przygotuj(ctx, pobierzWariantyMacierzyDebaty)
	if err != nil {
		return MacierzDebaty{}, err
	}
	wierszeWariantow, err := warianty.QueryContext(ctx, macierz.Kod)
	if err != nil {
		return MacierzDebaty{}, fmt.Errorf("dane: nie można odczytać wariantów macierzy %q: %w",
			macierz.Kod, err)
	}
	defer wierszeWariantow.Close()
	for wierszeWariantow.Next() {
		var wariant WariantMacierzyDebaty
		if err := wierszeWariantow.Scan(&wariant.Kod, &wariant.Macierz, &wariant.Etykieta,
			&wariant.OcenyJson, &wariant.Kolejnosc); err != nil {
			return MacierzDebaty{}, fmt.Errorf("dane: nieczytelny wiersz wariantu macierzy: %w", err)
		}
		macierz.Warianty = append(macierz.Warianty, wariant)
	}
	return macierz, wierszeWariantow.Err()
}
