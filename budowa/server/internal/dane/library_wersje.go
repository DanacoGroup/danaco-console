// Odpowiedzialność pliku: wersje pliku repozytorium wiedzy (tabela
// `wersja_pliku_biblioteki`) — trwałość Versioning Panelu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WersjaPlikuBiblioteki to wiersz tabeli `wersja_pliku_biblioteki`, niosący treść i sumę kontrolną wersji.
type WersjaPlikuBiblioteki struct {
	ID             int64
	Kod            string
	PlikID         int64
	Etykieta       *string
	Autor          *string
	RozmiarBajtow  *int64
	SumaKontrolna  *string
	TrescOdwolanie *string
	Utworzono      string
}

const (
	kolumnyWersji = `id, identyfikator_zewnetrzny, plik_id, etykieta, autor,
	                 rozmiar_bajtow, suma_kontrolna, tresc_odwolanie, utworzono`

	wstawWersjePlikuBiblioteki = `INSERT INTO wersja_pliku_biblioteki
	                    (identyfikator_zewnetrzny, plik_id, etykieta, autor,
	                     rozmiar_bajtow, suma_kontrolna, tresc_odwolanie)
	                    VALUES (?, ?, ?, ?, ?, ?, ?)`

	pobierzWersjePlikuBiblioteki = `SELECT ` + kolumnyWersji + ` FROM wersja_pliku_biblioteki
	                      WHERE identyfikator_zewnetrzny = ?`

	listaWersjiPliku = `SELECT ` + kolumnyWersji + ` FROM wersja_pliku_biblioteki
	                    WHERE plik_id = ? ORDER BY utworzono DESC, id DESC`

	// wersjaDoPrzywrocenia szuka wersji wyłącznie w obrębie wskazanego pliku —
	// kod wersji z innego pliku nie nadpisze bieżącego wskaźnika przy pomyłce
	// identyfikatora.
	wersjaDoPrzywrocenia = `SELECT ` + kolumnyWersji + ` FROM wersja_pliku_biblioteki
	                        WHERE identyfikator_zewnetrzny = ? AND plik_id = ?`

	przywrocBiezacaWersje = `UPDATE plik_biblioteki
	                         SET wersja_biezaca_id = ?, suma_kontrolna = ?,
	                             tresc_odwolanie = ?, rozmiar_bajtow = ?,
	                             zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                         WHERE id = ?`

	kodPlikuPoId = `SELECT identyfikator_zewnetrzny FROM plik_biblioteki WHERE id = ?`
)

// ZapiszWersje dokłada wiersz historii wersji pliku. Nie przestawia wskaźnika
// bieżącej wersji na pliku — robi to wyłącznie `PrzywrocWersje` albo zapis
// pliku, żeby jedna droga zmiany „co jest bieżące" nie rozjechała się z drugą.
func (r *repozytoriumBiblioteki) ZapiszWersje(ctx context.Context, plikID int64, wersja WersjaPlikuBiblioteki) (WersjaPlikuBiblioteki, error) {
	if wersja.Kod == "" || plikID == 0 {
		return WersjaPlikuBiblioteki{}, fmt.Errorf("dane: wersja pliku bez identyfikatora albo bez pliku")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawWersjePlikuBiblioteki)
	if err != nil {
		return WersjaPlikuBiblioteki{}, err
	}
	_, err = polecenie.ExecContext(ctx, wersja.Kod, plikID, tekstDoKolumny(wersja.Etykieta),
		tekstDoKolumny(wersja.Autor), liczbaDoKolumny(wersja.RozmiarBajtow),
		tekstDoKolumny(wersja.SumaKontrolna), tekstDoKolumny(wersja.TrescOdwolanie))
	if err != nil {
		return WersjaPlikuBiblioteki{}, fmt.Errorf("dane: nie można zapisać wersji %q pliku %d: %w", wersja.Kod, plikID, err)
	}
	return r.jednaWersja(ctx, pobierzWersjePlikuBiblioteki, wersja.Kod, "wersja "+wersja.Kod)
}

// Wersje zwraca wersje pliku od najnowszej — Versioning Panel pokazuje
// historię odwróconą, żeby ostatnia zmiana była na wierzchu.
func (r *repozytoriumBiblioteki) Wersje(ctx context.Context, plikID int64) ([]WersjaPlikuBiblioteki, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaWersjiPliku)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, plikID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wersji pliku %d: %w", plikID, err)
	}
	defer wiersze.Close()

	lista := []WersjaPlikuBiblioteki{}
	for wiersze.Next() {
		wersja, err := odczytajWersje(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wersji pliku %d: %w", plikID, err)
		}
		lista = append(lista, wersja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wersji pliku %d: %w", plikID, err)
	}
	return lista, nil
}

// PrzywrocWersje przestawia plik na wskazaną wersję historyczną: wskaźnik bieżącej wersji, sumę kontrolną,
// odwołanie do treści i rozmiar.
func (r *repozytoriumBiblioteki) PrzywrocWersje(ctx context.Context, plikID int64, kodWersji string) (PlikBiblioteki, error) {
	if kodWersji == "" || plikID == 0 {
		return PlikBiblioteki{}, fmt.Errorf("dane: przywrócenie wersji bez identyfikatora albo bez pliku")
	}
	wersja, err := r.jednaWersja(ctx, wersjaDoPrzywrocenia, []any{kodWersji, plikID},
		fmt.Sprintf("wersja %q pliku %d", kodWersji, plikID))
	if err != nil {
		return PlikBiblioteki{}, err
	}

	var kodPliku string
	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		aktualizuj, err := r.zapytania.wTransakcji(ctx, transakcja, przywrocBiezacaWersje)
		if err != nil {
			return err
		}
		if _, err := aktualizuj.ExecContext(ctx, wersja.ID, tekstDoKolumny(wersja.SumaKontrolna),
			tekstDoKolumny(wersja.TrescOdwolanie), liczbaDoKolumny(wersja.RozmiarBajtow), plikID); err != nil {
			return fmt.Errorf("dane: nie można przywrócić wersji %q na pliku %d: %w", kodWersji, plikID, err)
		}
		kod, err := r.zapytania.wTransakcji(ctx, transakcja, kodPlikuPoId)
		if err != nil {
			return err
		}
		if err := kod.QueryRowContext(ctx, plikID).Scan(&kodPliku); err != nil {
			return fmt.Errorf("dane: nie można odczytać kodu pliku %d po przywróceniu: %w", plikID, err)
		}
		return nil
	})
	if err != nil {
		return PlikBiblioteki{}, err
	}
	return r.Plik(ctx, kodPliku)
}

// jednaWersja wykonuje odczyt pojedynczego wiersza wersji, zasilany jednym
// albo wieloma argumentami zapytania.
func (r *repozytoriumBiblioteki) jednaWersja(ctx context.Context, zapytanie string, argumenty any, opis string) (WersjaPlikuBiblioteki, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return WersjaPlikuBiblioteki{}, err
	}
	lista, jest := argumenty.([]any)
	if !jest {
		lista = []any{argumenty}
	}
	wersja, err := odczytajWersje(polecenie.QueryRowContext(ctx, lista...))
	if errors.Is(err, sql.ErrNoRows) {
		return WersjaPlikuBiblioteki{}, ErrBrakWiersza
	}
	if err != nil {
		return WersjaPlikuBiblioteki{}, fmt.Errorf("dane: nieczytelny %s: %w", opis, err)
	}
	return wersja, nil
}

// odczytajWersje składa pełną strukturę wersji pliku z jednego wiersza wyniku zapytania do bazy danych.
func odczytajWersje(wiersz skaner) (WersjaPlikuBiblioteki, error) {
	var wersja WersjaPlikuBiblioteki
	var etykieta, autor, sumaKontrolna, trescOdwolanie sql.NullString
	var rozmiarBajtow sql.NullInt64
	err := wiersz.Scan(&wersja.ID, &wersja.Kod, &wersja.PlikID, &etykieta, &autor,
		&rozmiarBajtow, &sumaKontrolna, &trescOdwolanie, &wersja.Utworzono)
	if err != nil {
		return WersjaPlikuBiblioteki{}, err
	}
	wersja.Etykieta, wersja.Autor = tekstZKolumny(etykieta), tekstZKolumny(autor)
	wersja.RozmiarBajtow = liczbaZKolumny(rozmiarBajtow)
	wersja.SumaKontrolna, wersja.TrescOdwolanie = tekstZKolumny(sumaKontrolna), tekstZKolumny(trescOdwolanie)
	return wersja, nil
}
