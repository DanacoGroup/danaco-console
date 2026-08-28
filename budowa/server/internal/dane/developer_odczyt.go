// Warstwa danych obsługuje odczyt obszaru Developer: migawki pliku edytora
// i ostatni przebieg budowania okna, ze wspólnym parametrem limitu wykazu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	kolumnyWersjiPliku = `kod, okno_kod, sciezka, tresc, rozmiar, utworzono`

	pobierzOstatniaWersje = `SELECT ` + kolumnyWersjiPliku + `
	                         FROM developer_wersja_pliku
	                         WHERE okno_kod = ? AND sciezka = ?
	                         ORDER BY utworzono DESC, id DESC
	                         LIMIT 1`

	pobierzWersjePliku = `SELECT ` + kolumnyWersjiPliku + `
	                      FROM developer_wersja_pliku
	                      WHERE okno_kod = ? AND sciezka = ?
	                      ORDER BY utworzono DESC, id DESC
	                      LIMIT CASE WHEN ? > 0 THEN ? ELSE -1 END`

	kolumnyPrzebieguBudowania = `kod, okno_kod, zadanie, argumenty, stan, kod_wyjscia, log,
	                    uruchomiono, zakonczono`

	pobierzOstatniPrzebieg = `SELECT ` + kolumnyPrzebieguBudowania + `
	                          FROM developer_budowanie
	                          WHERE okno_kod = ?
	                          ORDER BY uruchomiono DESC, id DESC
	                          LIMIT 1`
)

type repozytoriumDevelopera struct {
	zapytania *zapytania
}

func noweRepozytoriumDevelopera(z *zapytania) *repozytoriumDevelopera {
	return &repozytoriumDevelopera{zapytania: z}
}

// OstatniaWersja zwraca najnowszą migawkę pliku. Brak migawki wraca jako
// ErrBrakWiersza — rdzeń odróżnia „plik nie ma wersji” od „odczyt zawiódł”.
func (r *repozytoriumDevelopera) OstatniaWersja(ctx context.Context,
	oknoKod, sciezka string) (WersjaPliku, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzOstatniaWersje)
	if err != nil {
		return WersjaPliku{}, err
	}
	wersja, err := odczytajWersjePliku(polecenie.QueryRowContext(ctx, oknoKod, sciezka))
	if errors.Is(err, sql.ErrNoRows) {
		return WersjaPliku{}, ErrBrakWiersza
	}
	if err != nil {
		return WersjaPliku{}, fmt.Errorf("dane: nieczytelna wersja pliku %q: %w", sciezka, err)
	}
	return wersja, nil
}

// Wersje zwraca migawki pliku wskazanej ścieżki okna od najnowszej; wartość
// limitu niedodatnia zwraca wykaz pełny.
func (r *repozytoriumDevelopera) Wersje(ctx context.Context,
	oknoKod, sciezka string, limit int) ([]WersjaPliku, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWersjePliku)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, oknoKod, sciezka, limit, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wersji pliku %q: %w", sciezka, err)
	}
	defer wiersze.Close()

	wersje := make([]WersjaPliku, 0, 8)
	for wiersze.Next() {
		wersja, err := odczytajWersjePliku(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wersji pliku: %w", err)
		}
		wersje = append(wersje, wersja)
	}
	return wersje, wiersze.Err()
}

// OstatniPrzebieg zwraca najnowszy przebieg budowania wskazanego okna,
// odczytany z tabeli developer_budowanie.
func (r *repozytoriumDevelopera) OstatniPrzebieg(ctx context.Context,
	oknoKod string) (PrzebiegBudowania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzOstatniPrzebieg)
	if err != nil {
		return PrzebiegBudowania{}, err
	}
	przebieg, err := odczytajPrzebiegBudowania(polecenie.QueryRowContext(ctx, oknoKod))
	if errors.Is(err, sql.ErrNoRows) {
		return PrzebiegBudowania{}, ErrBrakWiersza
	}
	if err != nil {
		return PrzebiegBudowania{}, fmt.Errorf("dane: nieczytelny przebieg budowania okna %q: %w", oknoKod, err)
	}
	return przebieg, nil
}

// odczytajWersjePliku składa migawkę WersjaPliku z jednego wiersza wyniku
// zapytania SQL o wersji pliku.
func odczytajWersjePliku(wiersz interface{ Scan(...any) error }) (WersjaPliku, error) {
	var wersja WersjaPliku
	err := wiersz.Scan(&wersja.Kod, &wersja.OknoKod, &wersja.Sciezka, &wersja.Tresc,
		&wersja.Rozmiar, &wersja.Utworzono)
	return wersja, err
}

// odczytajPrzebiegBudowania składa przebieg budowania PrzebiegBudowania
// z jednego wiersza wyniku zapytania.
func odczytajPrzebiegBudowania(wiersz interface{ Scan(...any) error }) (PrzebiegBudowania, error) {
	var przebieg PrzebiegBudowania
	err := wiersz.Scan(&przebieg.Kod, &przebieg.OknoKod, &przebieg.Zadanie, &przebieg.Argumenty,
		&przebieg.Stan, &przebieg.KodWyjscia, &przebieg.Log,
		&przebieg.Uruchomiono, &przebieg.Zakonczono)
	return przebieg, err
}
