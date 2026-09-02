// Plik prowadzi ustalenia badania i ich powiązania ze źródłami; źródła leżą w badania.go, a zapis powiązań jest
// wymianą kompletu — ZapiszUstalenie usuwa zastane powiązania i wstawia je od nowa w jednej transakcji.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// UstalenieBadania to wiersz tabeli `ustalenie_badania`. Powiązane źródła nie
// są polem tej struktury — leżą w osobnej tabeli złącznikowej i zwraca je
// `ZrodlaUstalenia`.
type UstalenieBadania struct {
	ID             int64
	Kod            string
	Okno           string
	Tresc          *string
	TrescOdwolanie *string
	Stan           shared.ResearchFindingStatus
	Utworzono      string
	Zaktualizowano string
}

const (
	kolumnyUstaleniaBadania = `id, identyfikator_zewnetrzny, okno, tresc, tresc_odwolanie,
	                           stan, utworzono, zaktualizowano`

	// Więz UNIQUE na `identyfikator_zewnetrzny` obejmuje całą tabelę, więc
	// warunek konta w gałęzi DO UPDATE zostawia wiersz cudzy nietknięty.
	zapiszUstalenieBadania = `INSERT INTO ustalenie_badania
	                          (identyfikator_zewnetrzny, okno, tresc, tresc_odwolanie, stan, konto_id)
	                          VALUES (?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                          ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                              tresc = excluded.tresc,
	                              tresc_odwolanie = excluded.tresc_odwolanie,
	                              stan = excluded.stan,
	                              zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                          WHERE ` + WarunekKonta

	pobierzUstalenieBadania = `SELECT ` + kolumnyUstaleniaBadania + ` FROM ustalenie_badania
	                           WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	pobierzUstaleniaBadaniaOkna = `SELECT ` + kolumnyUstaleniaBadania + ` FROM ustalenie_badania
	                               WHERE okno = ? AND ` + WarunekKonta + `
	                               ORDER BY zaktualizowano DESC, id DESC`

	// Klucz `ustalenie_id` pochodzi z odczytu zawężonego kontem.
	usunZrodlaUstalenia = `DELETE FROM zrodlo_ustalenia_badania WHERE ustalenie_id = ?`

	// idZrodlaPoKodzie odnajduje wiersz źródła po kodzie zewnętrznym — powiązanie
	// w tabeli złącznikowej trzyma klucz liczbowy, kontrakt oddaje kod tekstowy.
	idZrodlaPoKodzie = `SELECT id FROM zrodlo_badania WHERE identyfikator_zewnetrzny = ?`

	wstawZrodloUstalenia = `INSERT INTO zrodlo_ustalenia_badania (ustalenie_id, zrodlo_id)
	                        VALUES (?, ?)
	                        ON CONFLICT(ustalenie_id, zrodlo_id) DO NOTHING`

	pobierzZrodlaUstalenia = `SELECT ` + kolumnyZrodlaBadania + ` FROM zrodlo_badania
	                          WHERE id IN (
	                              SELECT zrodlo_id FROM zrodlo_ustalenia_badania
	                              WHERE ustalenie_id = ?
	                          )
	                          ORDER BY pozyskano_o DESC, id DESC`
)

// ZapiszUstalenie zakłada wiersz ustalenia albo nadpisuje zastany, po czym wymienia komplet powiązań ze źródłami
// w jednej transakcji; kod źródła bez odpowiadającego wiersza jest pomijany, nie wywraca zapisu ustalenia.
func (r *repozytoriumBadan) ZapiszUstalenie(ctx context.Context, ustalenie UstalenieBadania,
	kodyZrodel []string) (UstalenieBadania, error) {

	if ustalenie.Kod == "" {
		return UstalenieBadania{}, fmt.Errorf("dane: ustalenie badania bez identyfikatora")
	}
	if ustalenie.Okno == "" {
		return UstalenieBadania{}, fmt.Errorf("dane: ustalenie badania %q bez okna", ustalenie.Kod)
	}
	stan := string(ustalenie.Stan)
	if stan == "" {
		stan = string(shared.ResearchFindingStatusOpen)
	}

	var ustalenieID int64
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszUstalenieBadania)
		if err != nil {
			return err
		}
		wynik, err := zapis.ExecContext(ctx, ustalenie.Kod, ustalenie.Okno,
			tekstDoKolumny(ustalenie.Tresc), tekstDoKolumny(ustalenie.TrescOdwolanie), stan,
			KontoOperatora(ctx), KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać ustalenia badania %q: %w", ustalenie.Kod, err)
		}
		// Kod zewnętrzny zajęty przez wiersz konta obcego daje zero zmienionych
		// wierszy; dalszy odczyt oddałby „brak wiersza” zamiast powodu odmowy.
		if err := sprawdzTrafienieZapisu(wynik, "ustalenie badania", ustalenie.Kod); err != nil {
			return err
		}

		id, err := r.identyfikatorUstalenia(ctx, transakcja, ustalenie.Kod)
		if err != nil {
			return err
		}
		ustalenieID = id

		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunZrodlaUstalenia)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, ustalenieID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić źródeł ustalenia %d: %w", ustalenieID, err)
		}

		wyszukanie, err := r.zapytania.wTransakcji(ctx, transakcja, idZrodlaPoKodzie)
		if err != nil {
			return err
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawZrodloUstalenia)
		if err != nil {
			return err
		}
		for _, kodZrodla := range kodyZrodel {
			var zrodloID int64
			err := wyszukanie.QueryRowContext(ctx, kodZrodla).Scan(&zrodloID)
			if errors.Is(err, sql.ErrNoRows) {
				// Kod bez wiersza — pomijany, nie wywraca zapisu ustalenia.
				continue
			}
			if err != nil {
				return fmt.Errorf("dane: nie można odnaleźć źródła %q dla ustalenia %d: %w",
					kodZrodla, ustalenieID, err)
			}
			if _, err := wstawienie.ExecContext(ctx, ustalenieID, zrodloID); err != nil {
				return fmt.Errorf("dane: nie można powiązać źródła %q z ustaleniem %d: %w",
					kodZrodla, ustalenieID, err)
			}
		}
		return nil
	})
	if err != nil {
		return UstalenieBadania{}, err
	}
	return r.Ustalenie(ctx, ustalenie.Kod)
}

// identyfikatorUstalenia odczytuje klucz liczbowy ustalenia w obrębie
// bieżącej transakcji — insert/update powyżej nie zwraca go wprost przy
// konflikcie (ON CONFLICT DO UPDATE nie niesie LastInsertId na wierszu
// istniejącym).
func (r *repozytoriumBadan) identyfikatorUstalenia(ctx context.Context, transakcja *sql.Tx, kod string) (int64, error) {
	odczyt, err := r.zapytania.wTransakcji(ctx, transakcja,
		`SELECT id FROM ustalenie_badania WHERE identyfikator_zewnetrzny = ? AND `+WarunekKonta)
	if err != nil {
		return 0, err
	}
	var id int64
	if err := odczyt.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&id); err != nil {
		return 0, fmt.Errorf("dane: nie można odczytać identyfikatora ustalenia %q: %w", kod, err)
	}
	return id, nil
}

// Ustalenie zwraca ustalenie badania o wskazanym kodzie zewnętrznym; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumBadan) Ustalenie(ctx context.Context, kod string) (UstalenieBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzUstalenieBadania)
	if err != nil {
		return UstalenieBadania{}, err
	}
	ustalenie, err := odczytajUstalenieBadania(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return UstalenieBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return UstalenieBadania{}, fmt.Errorf("dane: nieczytelny wiersz ustalenia badania %q: %w", kod, err)
	}
	return ustalenie, nil
}

// Ustalenia zwraca ustalenia badania okna, posortowane od ostatnio zmienionych, wprost z bazy danych repozytorium.
func (r *repozytoriumBadan) Ustalenia(ctx context.Context, okno string) ([]UstalenieBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzUstaleniaBadaniaOkna)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać ustaleń badania okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []UstalenieBadania{}
	for wiersze.Next() {
		ustalenie, err := odczytajUstalenieBadania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz ustalenia badania: %w", err)
		}
		lista = append(lista, ustalenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt ustaleń badania okna %q: %w", okno, err)
	}
	return lista, nil
}

// ZrodlaUstalenia zwraca źródła powiązane z ustaleniem — komplet wynikły
// z ostatniego `ZapiszUstalenie`, bez kodów pominiętych przy zapisie.
func (r *repozytoriumBadan) ZrodlaUstalenia(ctx context.Context, ustalenieID int64) ([]ZrodloBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZrodlaUstalenia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, ustalenieID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać źródeł ustalenia %d: %w", ustalenieID, err)
	}
	defer wiersze.Close()

	lista := []ZrodloBadania{}
	for wiersze.Next() {
		zrodlo, err := odczytajZrodloBadania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz źródła ustalenia %d: %w", ustalenieID, err)
		}
		lista = append(lista, zrodlo)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt źródeł ustalenia %d: %w", ustalenieID, err)
	}
	return lista, nil
}

// odczytajUstalenieBadania składa strukturę ustalenia wprost z jednego wiersza wyniku zapytania do bazy.
func odczytajUstalenieBadania(wiersz skaner) (UstalenieBadania, error) {
	var ustalenie UstalenieBadania
	var tresc, trescOdwolanie sql.NullString
	var stan string
	err := wiersz.Scan(&ustalenie.ID, &ustalenie.Kod, &ustalenie.Okno, &tresc, &trescOdwolanie,
		&stan, &ustalenie.Utworzono, &ustalenie.Zaktualizowano)
	if err != nil {
		return UstalenieBadania{}, err
	}
	ustalenie.Stan = shared.ResearchFindingStatus(stan)
	ustalenie.Tresc = tekstZKolumny(tresc)
	ustalenie.TrescOdwolanie = tekstZKolumny(trescOdwolanie)
	return ustalenie, nil
}
