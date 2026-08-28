// Odpowiedzialność pliku: kosz sesji, przeniesienie do kosza, przywrócenie i czyszczenie po terminie, jako odwrotna strona tej samej tabeli.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// SesjaWKoszu to wiersz kosza: tyle, ile potrzeba, by pomyłkę odnaleźć
// i cofnąć — pełny wiersz sesji ma repozytorium sesji po przywróceniu.
type SesjaWKoszu struct {
	ID int64
	// Identyfikator rdzenia; pusty dla wiersza założonego poza rdzeniem.
	Identyfikator string
	Tytul         string
	UsunietoO     string
}

// RepozytoriumKoszaSesji jest kontraktem całego obszaru kosza sesji dla wszystkich warstw wyższych produktu.
type RepozytoriumKoszaSesji interface {
	// PrzeniesDoKosza stawia znacznik na sesji żywej.
	PrzeniesDoKosza(ctx context.Context, id int64) error
	// Przywroc czyści znacznik; fałsz mówi, że sesja w koszu nie leżała.
	Przywroc(ctx context.Context, id int64) (bool, error)
	// Lista zwraca zawartość kosza, od najświeższego wrzucenia.
	Lista(ctx context.Context) ([]SesjaWKoszu, error)
	// UsunPrzeterminowane kasuje trwale sesje leżące w koszu dłużej niż do granicy i zwraca ich liczbę.
	UsunPrzeterminowane(ctx context.Context, granica string) (int, error)
}

const (
	doKoszaSesje = `UPDATE sesja
	                SET usunieto_o = strftime('%Y-%m-%dT%H:%M:%fZ','now'),
	                    zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                WHERE id = ? AND usunieto_o IS NULL`

	zKoszaSesje = `UPDATE sesja
	               SET usunieto_o = NULL,
	                   zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	               WHERE id = ? AND usunieto_o IS NOT NULL`

	listaKosza = `SELECT id, COALESCE(identyfikator_zewnetrzny, ''), tytul, usunieto_o
	              FROM sesja WHERE usunieto_o IS NOT NULL
	              ORDER BY usunieto_o DESC, id DESC`

	// Bloki wiadomości okien sesji przeterminowanych — usuwane przed wierszami
	// sesji, póki łańcuch sesja → okno jeszcze istnieje.
	usunBlokiPrzeterminowane = `DELETE FROM blok_wiadomosci
	                            WHERE okno_kod IN (
	                                SELECT o.identyfikator_zewnetrzny
	                                FROM okno_komunikacji o
	                                JOIN sesja s ON s.id = o.sesja_id
	                                WHERE s.usunieto_o IS NOT NULL AND s.usunieto_o < ?
	                                  AND o.identyfikator_zewnetrzny IS NOT NULL)`

	usunSesjePrzeterminowane = `DELETE FROM sesja
	                            WHERE usunieto_o IS NOT NULL AND usunieto_o < ?`
)

type repozytoriumKoszaSesji struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumKoszaSesji(z *zapytania, db *sql.DB) *repozytoriumKoszaSesji {
	return &repozytoriumKoszaSesji{zapytania: z, db: db}
}

// PrzeniesDoKosza stawia znacznik kosza. Sesja już leżąca w koszu nie jest
// trafieniem — druga próba wraca błędem braku wiersza, żeby wywołujący nie
// wziął powtórki za skutek.
func (r *repozytoriumKoszaSesji) PrzeniesDoKosza(ctx context.Context, id int64) error {
	polecenie, err := r.zapytania.przygotuj(ctx, doKoszaSesje)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("dane: nie można przenieść sesji %d do kosza: %w", id, err)
	}
	return sprawdzTrafienie(wynik, "sesja", id)
}

// Przywroc czyści znacznik kosza. Sesja spoza kosza daje fałsz bez błędu —
// przywracanie zbiorcze nie może paść przez pozycję, która w koszu nie leży.
func (r *repozytoriumKoszaSesji) Przywroc(ctx context.Context, id int64) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zKoszaSesje)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, id)
	if err != nil {
		return false, fmt.Errorf("dane: nie można przywrócić sesji %d z kosza: %w", id, err)
	}
	liczba, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznana liczba zmienionych wierszy sesji %d: %w", id, err)
	}
	return liczba > 0, nil
}

// Lista zwraca całą zawartość kosza sesji dla tej warstwy wyższej; kosz pusty daje pustą listę bez błędu.
func (r *repozytoriumKoszaSesji) Lista(ctx context.Context) ([]SesjaWKoszu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKosza)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kosza sesji: %w", err)
	}
	defer wiersze.Close()

	lista := []SesjaWKoszu{}
	for wiersze.Next() {
		var pozycja SesjaWKoszu
		if err := wiersze.Scan(&pozycja.ID, &pozycja.Identyfikator,
			&pozycja.Tytul, &pozycja.UsunietoO); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kosza sesji: %w", err)
		}
		lista = append(lista, pozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kosza sesji: %w", err)
	}
	return lista, nil
}

// UsunPrzeterminowane kasuje trwale sesje po terminie kosza: najpierw bloki
// wiadomości ich okien, potem wiersze sesji (kaskada zabiera resztę). Jedna
// transakcja — zapis nie zostaje w stanie połowicznym.
func (r *repozytoriumKoszaSesji) UsunPrzeterminowane(ctx context.Context, granica string) (int, error) {
	usuniete := 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if _, err := transakcja.ExecContext(ctx, usunBlokiPrzeterminowane, granica); err != nil {
			return fmt.Errorf("dane: czyszczenie bloków kosza: %w", err)
		}
		wynik, err := transakcja.ExecContext(ctx, usunSesjePrzeterminowane, granica)
		if err != nil {
			return fmt.Errorf("dane: czyszczenie kosza sesji: %w", err)
		}
		liczba, err := wynik.RowsAffected()
		if err != nil {
			return fmt.Errorf("dane: nieznana liczba wierszy czyszczenia kosza: %w", err)
		}
		usuniete = int(liczba)
		return nil
	})
	if err != nil {
		return 0, err
	}
	return usuniete, nil
}
