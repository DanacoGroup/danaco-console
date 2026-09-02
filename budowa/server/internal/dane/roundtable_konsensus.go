// Odpowiedzialność pliku: wersje stanowiska, zdania odrębne i przekazanie stanowiska do modułu docelowego; wersja jest wpisem, nie licznikiem.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type WersjaStanowiskaDebaty struct {
	Kod        string
	Stanowisko string
	Wersja     int
	Tresc      string
	Utworzono  string
}

type ZdanieOdrebneDebaty struct {
	Kod        string
	Okno       string
	Stanowisko string
	Uczestnik  string
	Tresc      string
	Utworzono  string
}

type PrzekazanieDebaty struct {
	Kod        string
	Okno       string
	Stanowisko string
	Modul      string
	Artefakt   string
	Utworzono  string
}

type RepozytoriumDebatyKonsensusu interface {
	StanowiskoPoKodzie(ctx context.Context, kod string) (StanowiskoDebaty, error)
	RedagujStanowisko(ctx context.Context, stanowisko StanowiskoDebaty) (StanowiskoDebaty, error)

	ZapiszWersjeStanowiskaDebaty(ctx context.Context, wersja WersjaStanowiskaDebaty) error
	WersjeStanowiskaDebaty(ctx context.Context, stanowisko string) ([]WersjaStanowiskaDebaty, error)

	ZapiszZdanieOdrebneDebaty(ctx context.Context,
		zdanie ZdanieOdrebneDebaty) (ZdanieOdrebneDebaty, error)
	ZdaniaOdrebneDebaty(ctx context.Context, stanowisko string) ([]ZdanieOdrebneDebaty, error)

	ZapiszPrzekazanieDebaty(ctx context.Context, przekazanie PrzekazanieDebaty) error
}

const (
	// Wersja i przekazanie nie mają własnego konta: granica dochodzi przez korzeń debata_stanowisko (konto_id od migracji 484).
	// Powtórny zapis tej samej wersji nie jest błędem — stanowisko odczytane bez zmiany treści nie podbija licznika.
	zapiszWersjeStanowiskaDebaty = `INSERT INTO debata_stanowisko_wersja
	                                (identyfikator_zewnetrzny, stanowisko, wersja, tresc)
	                                SELECT ?, ?, ?, ? FROM debata_stanowisko
	                                 WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta + `
	                                ON CONFLICT(stanowisko, wersja) DO NOTHING`

	pobierzWersjeStanowiskaDebaty = `SELECT identyfikator_zewnetrzny, stanowisko, wersja, tresc,
	                                        utworzono
	                                 FROM debata_stanowisko_wersja WHERE stanowisko = ?
	                                   AND EXISTS (SELECT 1 FROM debata_stanowisko
	                                                WHERE debata_stanowisko.identyfikator_zewnetrzny
	                                                        = debata_stanowisko_wersja.stanowisko
	                                                  AND ` + WarunekKonta + `)
	                                 ORDER BY wersja ASC`

	granicaStanowiskaDebaty = `SELECT ` + WarunekKonta + ` FROM debata_stanowisko
	                           WHERE identyfikator_zewnetrzny = ?`

	zapiszZdanieOdrebneDebaty = `INSERT INTO debata_zdanie_odrebne
	                             (identyfikator_zewnetrzny, okno, stanowisko, uczestnik, tresc,
	                              konto_id)
	                             VALUES (?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                             ON CONFLICT(stanowisko, uczestnik) DO UPDATE SET
	                                 tresc = excluded.tresc,
	                                 utworzono = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                             WHERE ` + WarunekKonta

	pobierzZdanieOdrebneDebaty = `SELECT identyfikator_zewnetrzny, okno, stanowisko, uczestnik,
	                                     tresc, utworzono
	                              FROM debata_zdanie_odrebne
	                              WHERE stanowisko = ? AND uczestnik = ?
	                                AND ` + WarunekKonta

	pobierzZdaniaOdrebneDebaty = `SELECT identyfikator_zewnetrzny, okno, stanowisko, uczestnik,
	                                     tresc, utworzono
	                              FROM debata_zdanie_odrebne WHERE stanowisko = ?
	                                AND ` + WarunekKonta + ` ORDER BY id ASC`

	zapiszPrzekazanieDebaty = `INSERT INTO debata_przekazanie
	                           (identyfikator_zewnetrzny, okno, stanowisko, modul, artefakt)
	                           SELECT ?, ?, ?, ?, ? FROM debata_stanowisko
	                            WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

func (r *repozytoriumRoundtable) ZapiszWersjeStanowiskaDebaty(ctx context.Context,
	wersja WersjaStanowiskaDebaty) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWersjeStanowiskaDebaty)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, wersja.Kod, wersja.Stanowisko, wersja.Wersja,
		wersja.Tresc, wersja.Stanowisko, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać wersji stanowiska %q: %w", wersja.Stanowisko, err)
	}
	return r.trafienieStanowiska(ctx, wynik, wersja.Stanowisko)
}

func (r *repozytoriumRoundtable) WersjeStanowiskaDebaty(ctx context.Context,
	stanowisko string) ([]WersjaStanowiskaDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWersjeStanowiskaDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, stanowisko, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wersji stanowiska %q: %w", stanowisko, err)
	}
	defer wiersze.Close()

	wersje := make([]WersjaStanowiskaDebaty, 0, 8)
	for wiersze.Next() {
		var wersja WersjaStanowiskaDebaty
		if err := wiersze.Scan(&wersja.Kod, &wersja.Stanowisko, &wersja.Wersja, &wersja.Tresc,
			&wersja.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wersji stanowiska: %w", err)
		}
		wersje = append(wersje, wersja)
	}
	return wersje, wiersze.Err()
}

func (r *repozytoriumRoundtable) ZapiszZdanieOdrebneDebaty(ctx context.Context,
	zdanie ZdanieOdrebneDebaty) (ZdanieOdrebneDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZdanieOdrebneDebaty)
	if err != nil {
		return ZdanieOdrebneDebaty{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, zdanie.Kod, zdanie.Okno, zdanie.Stanowisko,
		zdanie.Uczestnik, zdanie.Tresc, KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return ZdanieOdrebneDebaty{}, fmt.Errorf("dane: nie można zapisać zdania odrębnego %q: %w",
			zdanie.Kod, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "zdanie odrębne", zdanie.Kod); err != nil {
		return ZdanieOdrebneDebaty{}, err
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzZdanieOdrebneDebaty)
	if err != nil {
		return ZdanieOdrebneDebaty{}, err
	}
	var zapisane ZdanieOdrebneDebaty
	wiersz := odczyt.QueryRowContext(ctx, zdanie.Stanowisko, zdanie.Uczestnik, KontoOperatora(ctx))
	if err := wiersz.Scan(&zapisane.Kod, &zapisane.Okno, &zapisane.Stanowisko,
		&zapisane.Uczestnik, &zapisane.Tresc, &zapisane.Utworzono); err != nil {
		return ZdanieOdrebneDebaty{}, fmt.Errorf("dane: nieczytelny wiersz zdania odrębnego: %w", err)
	}
	return zapisane, nil
}

func (r *repozytoriumRoundtable) ZdaniaOdrebneDebaty(ctx context.Context,
	stanowisko string) ([]ZdanieOdrebneDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZdaniaOdrebneDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, stanowisko, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zdań odrębnych %q: %w", stanowisko, err)
	}
	defer wiersze.Close()

	zdania := make([]ZdanieOdrebneDebaty, 0, 4)
	for wiersze.Next() {
		var zdanie ZdanieOdrebneDebaty
		if err := wiersze.Scan(&zdanie.Kod, &zdanie.Okno, &zdanie.Stanowisko, &zdanie.Uczestnik,
			&zdanie.Tresc, &zdanie.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zdania odrębnego: %w", err)
		}
		zdania = append(zdania, zdanie)
	}
	return zdania, wiersze.Err()
}

func (r *repozytoriumRoundtable) ZapiszPrzekazanieDebaty(ctx context.Context,
	przekazanie PrzekazanieDebaty) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPrzekazanieDebaty)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, przekazanie.Kod, przekazanie.Okno,
		przekazanie.Stanowisko, przekazanie.Modul, przekazanie.Artefakt,
		przekazanie.Stanowisko, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać przekazania stanowiska %q: %w",
			przekazanie.Kod, err)
	}
	return r.trafienieStanowiska(ctx, wynik, przekazanie.Stanowisko)
}

// Zero zmienionych wierszy przy zapisie dziecka stanowiska ma trzy powody: powtórka wersji (DO NOTHING),
// brak korzenia (ErrBrakWiersza) albo korzeń konta cudzego (ErrKolizjaWiersza); rozstrzyga odczyt korzenia.
func (r *repozytoriumRoundtable) trafienieStanowiska(ctx context.Context, wynik sql.Result,
	stanowisko string) error {

	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nieznana liczba zapisanych wierszy stanowiska %q: %w", stanowisko, err)
	}
	if zmienione > 0 {
		return nil
	}
	odczyt, err := r.zapytania.przygotuj(ctx, granicaStanowiskaDebaty)
	if err != nil {
		return err
	}
	var wGranicyKonta bool
	err = odczyt.QueryRowContext(ctx, KontoOperatora(ctx), stanowisko).Scan(&wGranicyKonta)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrBrakWiersza
	}
	if err != nil {
		return fmt.Errorf("dane: nie można odczytać stanowiska %q: %w", stanowisko, err)
	}
	if !wGranicyKonta {
		return fmt.Errorf("dane: stanowisko %q należy do innego konta: %w", stanowisko, ErrKolizjaWiersza)
	}
	return nil
}
