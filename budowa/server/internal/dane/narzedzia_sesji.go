// Odpowiedzialność pliku: dołożenia narzędzi żyjące w stanie sesji — zapis, odczyt i zdjęcie
// (tabela `narzedzie_sesji`).
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// NarzedzieSesji to jedno dołożenie: pozycja wykazu odpisana w chwili dołożenia, niesiona jako odpis, nie odwołanie.
type NarzedzieSesji struct {
	ID int64
	// NazwaPelna niesie przedrostek źródła narzędzia; jest tożsamością dołożenia w obrębie sesji.
	NazwaPelna string
	// NazwaSkrocona to nazwa bez przedrostka źródła — ta, którą się wpisuje.
	NazwaSkrocona string
	// Opis pełnym zdaniem: po co to jest i kiedy użyć.
	Opis string
	// Rodzaj niesie wartość kontraktu wprost (shared.SlashEntryKind).
	Rodzaj string
	// Grupa po przeznaczeniu — ta sama, którą niesie wykaz.
	Grupa string
	// ZrodloPozycji to przedrostek nazwy pełnej: skąd pozycja pochodzi.
	ZrodloPozycji string
	// Zrodlo niesie wartość kontraktu wprost (shared.SessionToolSource): czyja
	// ręka dołożyła.
	Zrodlo string
	// Dolozono to milisekundy epoki; podaje je warstwa wyższa, nie baza.
	Dolozono int64
}

// RepozytoriumNarzedziSesji jest kontraktem dołożeń jednej sesji; identyfikator sesji jest tu kluczem
// wiersza sesji.
type RepozytoriumNarzedziSesji interface {
	// Narzedzia oddaje dołożenia sesji w kolejności dokładania; sesja bez dołożeń oddaje wykaz pusty.
	Narzedzia(ctx context.Context, sesjaID int64) ([]NarzedzieSesji, error)
	// Doloz zapisuje dołożenie i oddaje je wraz z nadanym identyfikatorem oraz znacznikiem powtórzenia.
	Doloz(ctx context.Context, sesjaID int64, narzedzie NarzedzieSesji) (NarzedzieSesji, bool, error)
	// Zdejmij kasuje dołożenie po nazwie pełnej; fałsz znaczy „nie było czego zdejmować".
	Zdejmij(ctx context.Context, sesjaID int64, nazwaPelna string) (bool, error)
	// ZdejmijWszystkie kasuje wszystkie dołożenia sesji i oddaje ich liczbę usuniętych wierszy.
	ZdejmijWszystkie(ctx context.Context, sesjaID int64) (int, error)
}

const (
	kolumnyNarzedziaSesji = `id, nazwa_pelna, nazwa_skrocona, opis, rodzaj, grupa,
	                         zrodlo_pozycji, zrodlo, dolozono`

	wykazNarzedziSesji = `SELECT ` + kolumnyNarzedziaSesji + `
	                      FROM narzedzie_sesji
	                      WHERE sesja_id = ?
	                      ORDER BY dolozono, id`

	narzedzieSesjiPoNazwie = `SELECT ` + kolumnyNarzedziaSesji + `
	                          FROM narzedzie_sesji
	                          WHERE sesja_id = ? AND nazwa_pelna = ?`

	wstawNarzedzieSesji = `INSERT INTO narzedzie_sesji
	                           (sesja_id, nazwa_pelna, nazwa_skrocona, opis, rodzaj,
	                            grupa, zrodlo_pozycji, zrodlo, dolozono)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	skasujNarzedzieSesji = `DELETE FROM narzedzie_sesji
	                        WHERE sesja_id = ? AND nazwa_pelna = ?`

	skasujNarzedziaSesji = `DELETE FROM narzedzie_sesji WHERE sesja_id = ?`
)

type repozytoriumNarzedziSesji struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumNarzedziSesji(z *zapytania, db *sql.DB) *repozytoriumNarzedziSesji {
	return &repozytoriumNarzedziSesji{zapytania: z, db: db}
}

// NarzedziaSesji oddaje repozytorium dołożeń narzędzi nad tą samą bazą danych, co reszta zestawu repozytoriów.
func (z *Zestaw) NarzedziaSesji() RepozytoriumNarzedziSesji {
	if z == nil || z.zapytania == nil || z.zapytania.db == nil {
		return nil
	}
	return noweRepozytoriumNarzedziSesji(z.zapytania, z.zapytania.db)
}

// Narzedzia czyta wszystkie dołożenia jednej sesji w kolejności ich dokładania do stanu tej sesji rozmowy.
func (r *repozytoriumNarzedziSesji) Narzedzia(ctx context.Context,
	sesjaID int64) ([]NarzedzieSesji, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wykazNarzedziSesji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, sesjaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać dołożeń sesji %d: %w", sesjaID, err)
	}
	defer wiersze.Close()

	wykaz := make([]NarzedzieSesji, 0, 8)
	for wiersze.Next() {
		narzedzie, err := odczytajNarzedzieSesji(wiersze.Scan)
		if err != nil {
			return nil, fmt.Errorf("dane: uszkodzone dołożenie sesji %d: %w", sesjaID, err)
		}
		wykaz = append(wykaz, narzedzie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt dołożeń sesji %d: %w", sesjaID, err)
	}
	return wykaz, nil
}

// Doloz zapisuje dołożenie; odczyt i zapis idą jedną transakcją, bo razem odpowiadają na jedno pytanie
// o istnienie wiersza.
func (r *repozytoriumNarzedziSesji) Doloz(ctx context.Context, sesjaID int64,
	narzedzie NarzedzieSesji) (NarzedzieSesji, bool, error) {

	var zapisane NarzedzieSesji
	juzBylo := false
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, narzedzieSesjiPoNazwie)
		if err != nil {
			return err
		}
		zastane, err := odczytajNarzedzieSesji(
			odczyt.QueryRowContext(ctx, sesjaID, narzedzie.NazwaPelna).Scan)
		if err == nil {
			// Powtórzenie nie jest błędem i nie nadpisuje zastanego wiersza ani czasu jego dołożenia.
			zapisane, juzBylo = zastane, true
			return nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("dane: nie można sprawdzić dołożenia %q sesji %d: %w",
				narzedzie.NazwaPelna, sesjaID, err)
		}
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, wstawNarzedzieSesji)
		if err != nil {
			return err
		}
		wynik, err := zapis.ExecContext(ctx, sesjaID, narzedzie.NazwaPelna,
			narzedzie.NazwaSkrocona, narzedzie.Opis, narzedzie.Rodzaj, narzedzie.Grupa,
			narzedzie.ZrodloPozycji, narzedzie.Zrodlo, narzedzie.Dolozono)
		if err != nil {
			return fmt.Errorf("dane: nie można dołożyć %q do sesji %d: %w",
				narzedzie.NazwaPelna, sesjaID, err)
		}
		id, err := wynik.LastInsertId()
		if err != nil {
			return fmt.Errorf("dane: nieznany identyfikator dołożenia %q sesji %d: %w",
				narzedzie.NazwaPelna, sesjaID, err)
		}
		zapisane = narzedzie
		zapisane.ID = id
		return nil
	})
	if err != nil {
		return NarzedzieSesji{}, false, err
	}
	return zapisane, juzBylo, nil
}

// Zdejmij kasuje jedno dołożenie narzędzia po jego pełnej nazwie w obrębie tej wskazanej sesji rozmowy.
func (r *repozytoriumNarzedziSesji) Zdejmij(ctx context.Context, sesjaID int64,
	nazwaPelna string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, skasujNarzedzieSesji)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, sesjaID, nazwaPelna)
	if err != nil {
		return false, fmt.Errorf("dane: nie można zdjąć dołożenia %q sesji %d: %w",
			nazwaPelna, sesjaID, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek zdjęcia dołożenia %q sesji %d: %w",
			nazwaPelna, sesjaID, err)
	}
	return zmienione > 0, nil
}

// ZdejmijWszystkie kasuje wszystkie dołożenia narzędzi danej sesji rozmowy naraz, jednym poleceniem zapisu.
func (r *repozytoriumNarzedziSesji) ZdejmijWszystkie(ctx context.Context,
	sesjaID int64) (int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, skasujNarzedziaSesji)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, sesjaID)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zdjąć dołożeń sesji %d: %w", sesjaID, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieznany skutek zdjęcia dołożeń sesji %d: %w", sesjaID, err)
	}
	return int(zmienione), nil
}

// odczytajNarzedzieSesji składa wiersz z jednego skanu. Kolejność kolumn jest
// jedna (`kolumnyNarzedziaSesji`) i obsługuje oba odczyty — wykaz i pojedynczy
// wiersz — więc dopisanie kolumny nie rozjedzie jednego z nich.
func odczytajNarzedzieSesji(skan func(...any) error) (NarzedzieSesji, error) {
	var n NarzedzieSesji
	err := skan(&n.ID, &n.NazwaPelna, &n.NazwaSkrocona, &n.Opis, &n.Rodzaj, &n.Grupa,
		&n.ZrodloPozycji, &n.Zrodlo, &n.Dolozono)
	if err != nil {
		return NarzedzieSesji{}, err
	}
	return n, nil
}
