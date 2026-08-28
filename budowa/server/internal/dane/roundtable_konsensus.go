// Odpowiedzialność pliku: wersje stanowiska, zdania odrębne i przekazanie stanowiska do modułu docelowego; wersja jest wpisem, nie licznikiem.
package dane

import (
	"context"
	"fmt"
)

// WersjaStanowiskaDebaty to jedna redakcja stanowiska końcowego, zapisana osobno od poprzednich redakcji.
type WersjaStanowiskaDebaty struct {
	Kod        string
	Stanowisko string
	Wersja     int
	Tresc      string
	Utworzono  string
}

// ZdanieOdrebneDebaty to zdanie uczestnika, który nie dołączył do wypracowanego konsensusu tej debaty.
type ZdanieOdrebneDebaty struct {
	Kod        string
	Okno       string
	Stanowisko string
	Uczestnik  string
	Tresc      string
	Utworzono  string
}

// PrzekazanieDebaty to ślad przekazania stanowiska do modułu docelowego, wraz z chwilą tego przekazania.
type PrzekazanieDebaty struct {
	Kod        string
	Okno       string
	Stanowisko string
	Modul      string
	Artefakt   string
	Utworzono  string
}

// RepozytoriumDebatyKonsensusu jest częścią kontraktu obszaru odpowiadającą za
// stanowisko końcowe poza samą jego treścią.
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
	// Powtórny zapis tej samej wersji nie jest błędem: stanowisko odczytywane
	// wielokrotnie bez zmiany treści nie podbija licznika, więc wersja bieżąca
	// bywa zapisywana ponownie. Wpis zostaje ten, który był pierwszy.
	zapiszWersjeStanowiskaDebaty = `INSERT INTO debata_stanowisko_wersja
	                                (identyfikator_zewnetrzny, stanowisko, wersja, tresc)
	                                VALUES (?, ?, ?, ?)
	                                ON CONFLICT(stanowisko, wersja) DO NOTHING`

	pobierzWersjeStanowiskaDebaty = `SELECT identyfikator_zewnetrzny, stanowisko, wersja, tresc,
	                                        utworzono
	                                 FROM debata_stanowisko_wersja WHERE stanowisko = ?
	                                 ORDER BY wersja ASC`

	zapiszZdanieOdrebneDebaty = `INSERT INTO debata_zdanie_odrebne
	                             (identyfikator_zewnetrzny, okno, stanowisko, uczestnik, tresc)
	                             VALUES (?, ?, ?, ?, ?)
	                             ON CONFLICT(stanowisko, uczestnik) DO UPDATE SET
	                                 tresc = excluded.tresc,
	                                 utworzono = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzZdanieOdrebneDebaty = `SELECT identyfikator_zewnetrzny, okno, stanowisko, uczestnik,
	                                     tresc, utworzono
	                              FROM debata_zdanie_odrebne
	                              WHERE stanowisko = ? AND uczestnik = ?`

	pobierzZdaniaOdrebneDebaty = `SELECT identyfikator_zewnetrzny, okno, stanowisko, uczestnik,
	                                     tresc, utworzono
	                              FROM debata_zdanie_odrebne WHERE stanowisko = ? ORDER BY id ASC`

	zapiszPrzekazanieDebaty = `INSERT INTO debata_przekazanie
	                           (identyfikator_zewnetrzny, okno, stanowisko, modul, artefakt)
	                           VALUES (?, ?, ?, ?, ?)`
)

// ZapiszWersjeStanowiskaDebaty utrwala jedną redakcję stanowiska jako osobny wpis w dzienniku redakcji.
func (r *repozytoriumRoundtable) ZapiszWersjeStanowiskaDebaty(ctx context.Context,
	wersja WersjaStanowiskaDebaty) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWersjeStanowiskaDebaty)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, wersja.Kod, wersja.Stanowisko, wersja.Wersja,
		wersja.Tresc); err != nil {
		return fmt.Errorf("dane: nie można zapisać wersji stanowiska %q: %w", wersja.Stanowisko, err)
	}
	return nil
}

// WersjeStanowiskaDebaty zwraca wszystkie redakcje stanowiska od najstarszej do najnowszej jej wersji.
func (r *repozytoriumRoundtable) WersjeStanowiskaDebaty(ctx context.Context,
	stanowisko string) ([]WersjaStanowiskaDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWersjeStanowiskaDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, stanowisko)
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

// ZapiszZdanieOdrebneDebaty utrwala zdanie odrębne uczestnika, który nie dołączył do tego konsensusu debaty.
func (r *repozytoriumRoundtable) ZapiszZdanieOdrebneDebaty(ctx context.Context,
	zdanie ZdanieOdrebneDebaty) (ZdanieOdrebneDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZdanieOdrebneDebaty)
	if err != nil {
		return ZdanieOdrebneDebaty{}, err
	}
	if _, err := polecenie.ExecContext(ctx, zdanie.Kod, zdanie.Okno, zdanie.Stanowisko,
		zdanie.Uczestnik, zdanie.Tresc); err != nil {
		return ZdanieOdrebneDebaty{}, fmt.Errorf("dane: nie można zapisać zdania odrębnego %q: %w",
			zdanie.Kod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzZdanieOdrebneDebaty)
	if err != nil {
		return ZdanieOdrebneDebaty{}, err
	}
	var zapisane ZdanieOdrebneDebaty
	if err := odczyt.QueryRowContext(ctx, zdanie.Stanowisko, zdanie.Uczestnik).Scan(&zapisane.Kod,
		&zapisane.Okno, &zapisane.Stanowisko, &zapisane.Uczestnik, &zapisane.Tresc,
		&zapisane.Utworzono); err != nil {
		return ZdanieOdrebneDebaty{}, fmt.Errorf("dane: nieczytelny wiersz zdania odrębnego: %w", err)
	}
	return zapisane, nil
}

// ZdaniaOdrebneDebaty zwraca zdania odrębne podpisane wobec stanowiska, w kolejności ich zapisania w bazie.
func (r *repozytoriumRoundtable) ZdaniaOdrebneDebaty(ctx context.Context,
	stanowisko string) ([]ZdanieOdrebneDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZdaniaOdrebneDebaty)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, stanowisko)
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

// ZapiszPrzekazanieDebaty odnotowuje przekazanie stanowiska do wskazanego modułu docelowego tej debaty.
func (r *repozytoriumRoundtable) ZapiszPrzekazanieDebaty(ctx context.Context,
	przekazanie PrzekazanieDebaty) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPrzekazanieDebaty)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, przekazanie.Kod, przekazanie.Okno,
		przekazanie.Stanowisko, przekazanie.Modul, przekazanie.Artefakt); err != nil {
		return fmt.Errorf("dane: nie można zapisać przekazania stanowiska %q: %w",
			przekazanie.Kod, err)
	}
	return nil
}
