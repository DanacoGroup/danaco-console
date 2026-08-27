// Plik prowadzi historię wersji definicji automatyki wraz z rozdziałem wersji roboczej od opublikowanej i budżetami
// czasu przebiegu; wersja jest migawką martwą, więc kroki leżą tu jako zapis strukturalny, nie jako drugi komplet wierszy kroków.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WersjaAutomatyki to wiersz tabeli `wersja_automatyki` — migawka definicji automatyki z chwili jej zapisu.
type WersjaAutomatyki struct {
	AutomatykaID int64
	Wersja       int
	Kroki        string
	Opublikowana bool
	Autor        *string
	Utworzono    string
}

const (
	kolumnyWersjiAutomatyki = `automatyka_id, wersja, kroki, opublikowana, autor, utworzono`

	// Zapis wersji jest idempotentny po parze (automatyka, numer): ten sam
	// zapis definicji powtórzony nie zakłada drugiej migawki tego numeru.
	zapiszWersjeAutomatyki = `INSERT INTO wersja_automatyki
	                          (automatyka_id, wersja, kroki, autor)
	                          VALUES (?, ?, ?, ?)
	                          ON CONFLICT(automatyka_id, wersja) DO UPDATE SET
	                              kroki = excluded.kroki,
	                              autor = excluded.autor`

	listaWersjiAutomatyki = `SELECT ` + kolumnyWersjiAutomatyki + ` FROM wersja_automatyki
	                         WHERE automatyka_id = ? ORDER BY wersja DESC LIMIT ?`

	pobierzWersjeAutomatyki = `SELECT ` + kolumnyWersjiAutomatyki + ` FROM wersja_automatyki
	                           WHERE automatyka_id = ? AND wersja = ?`

	// Publikacja jest jedna na automatykę: zdjęcie znacznika ze wszystkich
	// wersji i nadanie go jednej idzie w tej samej transakcji, żeby historia
	// nie pokazała przez chwilę dwóch wersji opublikowanych naraz.
	zdejmijPublikacjeWersji = `UPDATE wersja_automatyki SET opublikowana = 0 WHERE automatyka_id = ?`

	nadajPublikacjeWersji = `UPDATE wersja_automatyki SET opublikowana = 1
	                         WHERE automatyka_id = ? AND wersja = ?`

	ustawWersjeOpublikowana = `UPDATE automatyka
	                           SET wersja_opublikowana = ?,
	                               zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                           WHERE id = ?`

	ustawUdostepnienieAutomatyki = `UPDATE automatyka
	                                SET udostepniona = ?,
	                                    zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                                WHERE id = ?`

	ustawBudzetyAutomatyki = `UPDATE automatyka
	                          SET budzet_przebiegu_sekundy = ?, budzet_kroku_sekundy = ?,
	                              regula_budzetu = ?,
	                              zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                          WHERE id = ?`
)

// ZapiszWersjeAutomatyki dopisuje migawkę definicji automatyki wraz z jej pełnym kompletem kroków wykonania.
func (r *repozytoriumAutomatyk) ZapiszWersjeAutomatyki(ctx context.Context,
	wersja WersjaAutomatyki) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWersjeAutomatyki)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, wersja.AutomatykaID, wersja.Wersja, wersja.Kroki,
		tekstDoKolumny(wersja.Autor))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać wersji %d automatyki %d: %w",
			wersja.Wersja, wersja.AutomatykaID, err)
	}
	return nil
}

// WersjeAutomatyki zwraca migawki definicji automatyki od najnowszej wprost z bazy danych repozytorium.
func (r *repozytoriumAutomatyk) WersjeAutomatyki(ctx context.Context,
	automatykaID int64, limit int) ([]WersjaAutomatyki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWersjiAutomatyki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, automatykaID, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wersji automatyki %d: %w", automatykaID, err)
	}
	defer wiersze.Close()

	lista := []WersjaAutomatyki{}
	for wiersze.Next() {
		wersja, err := odczytajWersjeAutomatyki(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wersji automatyki: %w", err)
		}
		lista = append(lista, wersja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wersji automatyki: %w", err)
	}
	return lista, nil
}

// WersjaAutomatykiNumer zwraca jedną migawkę. Brak wiersza wraca jako
// ErrBrakWiersza — przywrócenie wersji, której nie ma, jest błędem żądania.
func (r *repozytoriumAutomatyk) WersjaAutomatykiNumer(ctx context.Context,
	automatykaID int64, numer int) (WersjaAutomatyki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWersjeAutomatyki)
	if err != nil {
		return WersjaAutomatyki{}, err
	}
	wersja, err := odczytajWersjeAutomatyki(polecenie.QueryRowContext(ctx, automatykaID, numer))
	if errors.Is(err, sql.ErrNoRows) {
		return WersjaAutomatyki{}, ErrBrakWiersza
	}
	if err != nil {
		return WersjaAutomatyki{}, fmt.Errorf("dane: nieczytelna wersja %d automatyki %d: %w",
			numer, automatykaID, err)
	}
	return wersja, nil
}

// OpublikujWersjeAutomatyki nadaje znacznik publikacji jednej wersji i zapisuje
// jej numer przy automatyce. Obie zmiany idą w jednej transakcji: automatyka
// wskazująca wersję, która o publikacji nie wie, byłaby dwiema prawdami o tym
// samym fakcie.
func (r *repozytoriumAutomatyk) OpublikujWersjeAutomatyki(ctx context.Context,
	automatykaID int64, numer int) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		for _, para := range []struct {
			zapytanie string
			argumenty []any
		}{
			{zdejmijPublikacjeWersji, []any{automatykaID}},
			{nadajPublikacjeWersji, []any{automatykaID, numer}},
			{ustawWersjeOpublikowana, []any{numer, automatykaID}},
		} {
			polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, para.zapytanie)
			if err != nil {
				return err
			}
			if _, err := polecenie.ExecContext(ctx, para.argumenty...); err != nil {
				return fmt.Errorf("dane: nie można opublikować wersji %d automatyki %d: %w",
					numer, automatykaID, err)
			}
		}
		return nil
	})
}

// UstawUdostepnienieAutomatyki przełącza udostępnienie automatyki w organizacji i zwraca stan po zmianie.
func (r *repozytoriumAutomatyk) UstawUdostepnienieAutomatyki(ctx context.Context,
	automatykaID int64, udostepniona bool) error {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawUdostepnienieAutomatyki)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, liczbaLogiczna(udostepniona), automatykaID); err != nil {
		return fmt.Errorf("dane: nie można zmienić udostępnienia automatyki %d: %w", automatykaID, err)
	}
	return nil
}

// UstawBudzetyAutomatyki zapisuje budżet czasu przebiegu i kroku wraz z regułą
// alarmowania powiadamianą przy przekroczeniu. Zero znaczy „bez granicy”.
func (r *repozytoriumAutomatyk) UstawBudzetyAutomatyki(ctx context.Context, automatykaID int64,
	budzetPrzebiegu, budzetKroku int, regula *string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, ustawBudzetyAutomatyki)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, budzetPrzebiegu, budzetKroku,
		tekstDoKolumny(regula), automatykaID)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać budżetów automatyki %d: %w", automatykaID, err)
	}
	return nil
}

// odczytajWersjeAutomatyki składa migawkę wprost z jednego wiersza wyniku zapytania do bazy danych repozytorium.
func odczytajWersjeAutomatyki(wiersz skaner) (WersjaAutomatyki, error) {
	var wersja WersjaAutomatyki
	var autor sql.NullString
	var opublikowana int
	err := wiersz.Scan(&wersja.AutomatykaID, &wersja.Wersja, &wersja.Kroki,
		&opublikowana, &autor, &wersja.Utworzono)
	if err != nil {
		return WersjaAutomatyki{}, err
	}
	wersja.Opublikowana = opublikowana == 1
	wersja.Autor = tekstZKolumny(autor)
	return wersja, nil
}
