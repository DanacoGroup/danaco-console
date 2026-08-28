// Plik utrzymuje ocenę Operatora, rubryki oceny, werdykty modeli-sędziów
// i ranking akumulowany między sesjami debaty.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// OcenaDebaty to ocena wypowiedzi debaty postawiona wprost przez Operatora w oknie komunikacji Roundtable.
type OcenaDebaty struct {
	Kod       string
	Okno      string
	Rodzaj    string
	Wypowiedz string
	Uczestnik string
	Gwiazdki  int
	Wskazana  string
	Utworzono string
}

// RubrykaDebaty to nazwana rubryka oceny debaty wraz z jej kryteriami, każdym niosącym własną wagę oceny.
type RubrykaDebaty struct {
	Kod       string
	Okno      string
	Nazwa     string
	Utworzono string
	Kryteria  []KryteriumRubrykiDebaty
}

// KryteriumRubrykiDebaty to jedno kryterium rubryki wraz z wagą wchodzącą do sumy jedności całej rubryki.
type KryteriumRubrykiDebaty struct {
	Kod       string
	Rubryka   string
	Nazwa     string
	Waga      float64
	Opis      *string
	Kolejnosc int
}

// WerdyktDebaty to ocena wystawiona przez model-sędziego wraz z uzasadnieniem i punktacją kryteriów rubryki.
type WerdyktDebaty struct {
	Kod          string
	Okno         string
	Rubryka      string
	Sedzia       string
	Wypowiedz    string
	Uczestnik    string
	PunktyJson   string
	Wynik        float64
	Uzasadnienie string
	Utworzono    string
}

// PozycjaRankinguDebaty to jedna tożsamość w rankingu wraz z jej punktacją i liczbą pojedynków rozegranych.
type PozycjaRankinguDebaty struct {
	KluczTozsamosci string
	Nazwa           string
	Zakres          string
	Okno            string
	Algorytm        string
	Punktacja       float64
	Odchylenie      float64
	Pojedynki       int
	Wygrane         int
	Zaktualizowano  string
}

// RepozytoriumDebatyOceny jest częścią kontraktu obszaru Roundtable odpowiadającą za ocenę i ranking debat.
type RepozytoriumDebatyOceny interface {
	ZapiszOceneDebaty(ctx context.Context, ocena OcenaDebaty) (OcenaDebaty, error)
	OcenyDebaty(ctx context.Context, okno string) ([]OcenaDebaty, error)

	ZapiszRubrykeDebaty(ctx context.Context, rubryka RubrykaDebaty) (RubrykaDebaty, error)
	RubrykaDebatyPoKodzie(ctx context.Context, kod string) (RubrykaDebaty, error)
	RubrykiDebaty(ctx context.Context, okno string) ([]RubrykaDebaty, error)

	ZapiszWerdyktDebaty(ctx context.Context, werdykt WerdyktDebaty) error
	WerdyktyDebaty(ctx context.Context, okno string) ([]WerdyktDebaty, error)

	ZapiszPozycjeRankinguDebaty(ctx context.Context, pozycja PozycjaRankinguDebaty) error
	PozycjaRankinguDebaty(ctx context.Context, klucz, zakres, okno,
		algorytm string) (PozycjaRankinguDebaty, error)
	RankingDebaty(ctx context.Context, zakres, okno, algorytm string,
		limit int) ([]PozycjaRankinguDebaty, error)
}

const (
	kolumnyOcenyDebaty = `identyfikator_zewnetrzny, okno, rodzaj, wypowiedz, uczestnik,
	                      gwiazdki, wskazana, utworzono`

	zapiszOceneDebaty = `INSERT INTO debata_ocena
	                     (identyfikator_zewnetrzny, okno, rodzaj, wypowiedz, uczestnik,
	                      gwiazdki, wskazana)
	                     VALUES (?, ?, ?, ?, ?, ?, ?)`

	pobierzOcenePoKodzieDebaty = `SELECT ` + kolumnyOcenyDebaty + `
	                              FROM debata_ocena WHERE identyfikator_zewnetrzny = ?`

	pobierzOcenyDebaty = `SELECT ` + kolumnyOcenyDebaty + `
	                      FROM debata_ocena WHERE okno = ? ORDER BY id ASC`

	zapiszRubrykeDebaty = `INSERT INTO debata_rubryka (identyfikator_zewnetrzny, okno, nazwa)
	                       VALUES (?, ?, ?)
	                       ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                           nazwa = excluded.nazwa, okno = excluded.okno`

	usunKryteriaRubrykiDebaty = `DELETE FROM debata_kryterium WHERE rubryka = ?`

	zapiszKryteriumRubrykiDebaty = `INSERT INTO debata_kryterium
	                                (identyfikator_zewnetrzny, rubryka, nazwa, waga, opis, kolejnosc)
	                                VALUES (?, ?, ?, ?, ?, ?)`

	pobierzRubrykeDebaty = `SELECT identyfikator_zewnetrzny, okno, nazwa, utworzono
	                        FROM debata_rubryka WHERE identyfikator_zewnetrzny = ?`

	// Okno puste w żądaniu oddaje same rubryki wspólne; wskazanie okna oddaje
	// wspólne i te należące do wskazanego okna.
	pobierzRubrykiDebaty = `SELECT identyfikator_zewnetrzny, okno, nazwa, utworzono
	                        FROM debata_rubryka WHERE okno = '' OR okno = ?
	                        ORDER BY id ASC`

	pobierzKryteriaRubrykiDebaty = `SELECT identyfikator_zewnetrzny, rubryka, nazwa, waga, opis,
	                                       kolejnosc
	                                FROM debata_kryterium WHERE rubryka = ?
	                                ORDER BY kolejnosc ASC, id ASC`

	zapiszWerdyktDebaty = `INSERT INTO debata_werdykt
	                       (identyfikator_zewnetrzny, okno, rubryka, sedzia, wypowiedz, uczestnik,
	                        punkty_json, wynik, uzasadnienie)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzWerdyktyDebaty = `SELECT identyfikator_zewnetrzny, okno, rubryka, sedzia, wypowiedz,
	                                uczestnik, punkty_json, wynik, uzasadnienie, utworzono
	                         FROM debata_werdykt WHERE okno = ? ORDER BY id ASC`

	kolumnyRankinguDebaty = `klucz_tozsamosci, nazwa, zakres, okno, algorytm, punktacja,
	                         odchylenie, pojedynki, wygrane, zaktualizowano`

	zapiszRankingDebaty = `INSERT INTO debata_ranking
	                       (klucz_tozsamosci, nazwa, zakres, okno, algorytm, punktacja,
	                        odchylenie, pojedynki, wygrane)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	                       ON CONFLICT(klucz_tozsamosci, zakres, okno, algorytm) DO UPDATE SET
	                           nazwa = excluded.nazwa,
	                           punktacja = excluded.punktacja,
	                           odchylenie = excluded.odchylenie,
	                           pojedynki = excluded.pojedynki,
	                           wygrane = excluded.wygrane,
	                           zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzPozycjeRankinguDebaty = `SELECT ` + kolumnyRankinguDebaty + `
	                                FROM debata_ranking
	                                WHERE klucz_tozsamosci = ? AND zakres = ? AND okno = ?
	                                  AND algorytm = ?`

	pobierzRankingDebaty = `SELECT ` + kolumnyRankinguDebaty + `
	                        FROM debata_ranking
	                        WHERE zakres = ? AND okno = ? AND algorytm = ?
	                        ORDER BY punktacja DESC, pojedynki DESC
	                        LIMIT (CASE WHEN ? > 0 THEN ? ELSE -1 END)`
)

// ZapiszOceneDebaty dopisuje ocenę Operatora do dziennika ocen wypowiedzi debaty w danym oknie Roundtable.
func (r *repozytoriumRoundtable) ZapiszOceneDebaty(ctx context.Context,
	ocena OcenaDebaty) (OcenaDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszOceneDebaty)
	if err != nil {
		return OcenaDebaty{}, err
	}
	if _, err := polecenie.ExecContext(ctx, ocena.Kod, ocena.Okno, ocena.Rodzaj, ocena.Wypowiedz,
		ocena.Uczestnik, ocena.Gwiazdki, ocena.Wskazana); err != nil {
		return OcenaDebaty{}, fmt.Errorf("dane: nie można zapisać oceny debaty %q: %w", ocena.Kod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzOcenePoKodzieDebaty)
	if err != nil {
		return OcenaDebaty{}, err
	}
	return odczytajOceneDebaty(odczyt.QueryRowContext(ctx, ocena.Kod))
}

// OcenyDebaty zwraca wszystkie oceny Operatora postawione w oknie debaty, w kolejności ich zapisu do dziennika.
func (r *repozytoriumRoundtable) OcenyDebaty(ctx context.Context, okno string) ([]OcenaDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzOcenyDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać ocen debaty okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	oceny := make([]OcenaDebaty, 0, 8)
	for wiersze.Next() {
		ocena, err := odczytajOceneDebaty(wiersze)
		if err != nil {
			return nil, err
		}
		oceny = append(oceny, ocena)
	}
	return oceny, wiersze.Err()
}

// ZapiszRubrykeDebaty zakłada albo zmienia rubrykę wraz z kryteriami w jednej
// transakcji. Kryteria wymienia się w całości: rubryka to zestaw wag sumujący
// się do jedności, więc dopisywanie po jednym zostawiałoby ją w stanie
// niesumującym się do niczego.
func (r *repozytoriumRoundtable) ZapiszRubrykeDebaty(ctx context.Context,
	rubryka RubrykaDebaty) (RubrykaDebaty, error) {

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		naglowek, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszRubrykeDebaty)
		if err != nil {
			return err
		}
		if _, err := naglowek.ExecContext(ctx, rubryka.Kod, rubryka.Okno, rubryka.Nazwa); err != nil {
			return fmt.Errorf("dane: nie można zapisać rubryki debaty %q: %w", rubryka.Kod, err)
		}
		wyczysc, err := r.zapytania.wTransakcji(ctx, transakcja, usunKryteriaRubrykiDebaty)
		if err != nil {
			return err
		}
		if _, err := wyczysc.ExecContext(ctx, rubryka.Kod); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić kryteriów rubryki %q: %w", rubryka.Kod, err)
		}
		wstaw, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszKryteriumRubrykiDebaty)
		if err != nil {
			return err
		}
		for pozycja, kryterium := range rubryka.Kryteria {
			if _, err := wstaw.ExecContext(ctx, kryterium.Kod, rubryka.Kod, kryterium.Nazwa,
				kryterium.Waga, kryterium.Opis, pozycja+1); err != nil {
				return fmt.Errorf("dane: nie można zapisać kryterium rubryki %q: %w",
					kryterium.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return RubrykaDebaty{}, err
	}
	return r.RubrykaDebatyPoKodzie(ctx, rubryka.Kod)
}

// RubrykaDebatyPoKodzie zwraca rubrykę oceny debaty wraz z jej kryteriami po jej kodzie zewnętrznym wprost.
func (r *repozytoriumRoundtable) RubrykaDebatyPoKodzie(ctx context.Context,
	kod string) (RubrykaDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzRubrykeDebaty)
	if err != nil {
		return RubrykaDebaty{}, err
	}
	var rubryka RubrykaDebaty
	err = polecenie.QueryRowContext(ctx, kod).Scan(&rubryka.Kod, &rubryka.Okno, &rubryka.Nazwa,
		&rubryka.Utworzono)
	if errors.Is(err, sql.ErrNoRows) {
		return RubrykaDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return RubrykaDebaty{}, fmt.Errorf("dane: nieczytelny wiersz rubryki debaty %q: %w", kod, err)
	}
	rubryka.Kryteria, err = r.kryteriaRubrykiDebaty(ctx, kod)
	if err != nil {
		return RubrykaDebaty{}, err
	}
	return rubryka, nil
}

// RubrykiDebaty zwraca rubryki wspólne oraz te należące do wskazanego okna, wraz z kompletem ich kryteriów.
func (r *repozytoriumRoundtable) RubrykiDebaty(ctx context.Context,
	okno string) ([]RubrykaDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzRubrykiDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać rubryk debaty: %w", err)
	}
	defer wiersze.Close()

	rubryki := make([]RubrykaDebaty, 0, 8)
	for wiersze.Next() {
		var rubryka RubrykaDebaty
		if err := wiersze.Scan(&rubryka.Kod, &rubryka.Okno, &rubryka.Nazwa,
			&rubryka.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz rubryki debaty: %w", err)
		}
		rubryki = append(rubryki, rubryka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, err
	}
	for i := range rubryki {
		kryteria, err := r.kryteriaRubrykiDebaty(ctx, rubryki[i].Kod)
		if err != nil {
			return nil, err
		}
		rubryki[i].Kryteria = kryteria
	}
	return rubryki, nil
}

// kryteriaRubrykiDebaty czyta kryteria jednej rubryki w kolejności ich zapisanej pozycji wśród kryteriów.
func (r *repozytoriumRoundtable) kryteriaRubrykiDebaty(ctx context.Context,
	rubryka string) ([]KryteriumRubrykiDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKryteriaRubrykiDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, rubryka)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kryteriów rubryki %q: %w", rubryka, err)
	}
	defer wiersze.Close()

	kryteria := make([]KryteriumRubrykiDebaty, 0, 8)
	for wiersze.Next() {
		var kryterium KryteriumRubrykiDebaty
		if err := wiersze.Scan(&kryterium.Kod, &kryterium.Rubryka, &kryterium.Nazwa,
			&kryterium.Waga, &kryterium.Opis, &kryterium.Kolejnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kryterium rubryki: %w", err)
		}
		kryteria = append(kryteria, kryterium)
	}
	return kryteria, wiersze.Err()
}

// ZapiszWerdyktDebaty dopisuje ocenę wystawioną przez sędziego do dziennika werdyktów debaty Roundtable.
func (r *repozytoriumRoundtable) ZapiszWerdyktDebaty(ctx context.Context, werdykt WerdyktDebaty) error {
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWerdyktDebaty)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, werdykt.Kod, werdykt.Okno, werdykt.Rubryka,
		werdykt.Sedzia, werdykt.Wypowiedz, werdykt.Uczestnik, werdykt.PunktyJson, werdykt.Wynik,
		werdykt.Uzasadnienie); err != nil {
		return fmt.Errorf("dane: nie można zapisać werdyktu %q: %w", werdykt.Kod, err)
	}
	return nil
}

// WerdyktyDebaty zwraca werdykty sędziów wystawione w oknie debaty, w kolejności ich zapisu do dziennika.
func (r *repozytoriumRoundtable) WerdyktyDebaty(ctx context.Context,
	okno string) ([]WerdyktDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWerdyktyDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać werdyktów okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	werdykty := make([]WerdyktDebaty, 0, 8)
	for wiersze.Next() {
		var werdykt WerdyktDebaty
		if err := wiersze.Scan(&werdykt.Kod, &werdykt.Okno, &werdykt.Rubryka, &werdykt.Sedzia,
			&werdykt.Wypowiedz, &werdykt.Uczestnik, &werdykt.PunktyJson, &werdykt.Wynik,
			&werdykt.Uzasadnienie, &werdykt.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz werdyktu: %w", err)
		}
		werdykty = append(werdykty, werdykt)
	}
	return werdykty, wiersze.Err()
}

// ZapiszPozycjeRankinguDebaty utrwala punktację tożsamości po pojedynku, zakładając albo nadpisując wiersz.
func (r *repozytoriumRoundtable) ZapiszPozycjeRankinguDebaty(ctx context.Context,
	pozycja PozycjaRankinguDebaty) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszRankingDebaty)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, pozycja.KluczTozsamosci, pozycja.Nazwa, pozycja.Zakres,
		pozycja.Okno, pozycja.Algorytm, pozycja.Punktacja, pozycja.Odchylenie, pozycja.Pojedynki,
		pozycja.Wygrane); err != nil {
		return fmt.Errorf("dane: nie można zapisać pozycji rankingu %q: %w",
			pozycja.KluczTozsamosci, err)
	}
	return nil
}

// PozycjaRankinguDebaty zwraca punktację jednej tożsamości w danym zakresie, oknie i algorytmie rankingu.
func (r *repozytoriumRoundtable) PozycjaRankinguDebaty(ctx context.Context,
	klucz, zakres, okno, algorytm string) (PozycjaRankinguDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPozycjeRankinguDebaty)
	if err != nil {
		return PozycjaRankinguDebaty{}, err
	}
	pozycja, err := odczytajPozycjeRankinguDebaty(
		polecenie.QueryRowContext(ctx, klucz, zakres, okno, algorytm))
	if errors.Is(err, sql.ErrNoRows) {
		return PozycjaRankinguDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return PozycjaRankinguDebaty{}, fmt.Errorf("dane: nieczytelny wiersz rankingu %q: %w",
			klucz, err)
	}
	return pozycja, nil
}

// RankingDebaty zwraca ranking od najwyższej punktacji, zawężony zakresem, oknem i algorytmem liczenia.
func (r *repozytoriumRoundtable) RankingDebaty(ctx context.Context, zakres, okno, algorytm string,
	limit int) ([]PozycjaRankinguDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzRankingDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, zakres, okno, algorytm, limit, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać rankingu debaty: %w", err)
	}
	defer wiersze.Close()

	ranking := make([]PozycjaRankinguDebaty, 0, 8)
	for wiersze.Next() {
		pozycja, err := odczytajPozycjeRankinguDebaty(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz rankingu debaty: %w", err)
		}
		ranking = append(ranking, pozycja)
	}
	return ranking, wiersze.Err()
}

// odczytajOceneDebaty składa ocenę Operatora z jednego wiersza wyniku danego zapytania do bazy danych.
func odczytajOceneDebaty(wiersz interface{ Scan(...any) error }) (OcenaDebaty, error) {
	var ocena OcenaDebaty
	err := wiersz.Scan(&ocena.Kod, &ocena.Okno, &ocena.Rodzaj, &ocena.Wypowiedz, &ocena.Uczestnik,
		&ocena.Gwiazdki, &ocena.Wskazana, &ocena.Utworzono)
	if errors.Is(err, sql.ErrNoRows) {
		return OcenaDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return OcenaDebaty{}, fmt.Errorf("dane: nieczytelny wiersz oceny debaty: %w", err)
	}
	return ocena, nil
}

// odczytajPozycjeRankinguDebaty składa pozycję rankingu tożsamości z jednego wiersza wyniku zapytania.
func odczytajPozycjeRankinguDebaty(wiersz interface{ Scan(...any) error }) (PozycjaRankinguDebaty, error) {
	var pozycja PozycjaRankinguDebaty
	err := wiersz.Scan(&pozycja.KluczTozsamosci, &pozycja.Nazwa, &pozycja.Zakres, &pozycja.Okno,
		&pozycja.Algorytm, &pozycja.Punktacja, &pozycja.Odchylenie, &pozycja.Pojedynki,
		&pozycja.Wygrane, &pozycja.Zaktualizowano)
	return pozycja, err
}
