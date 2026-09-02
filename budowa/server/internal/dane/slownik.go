// Odpowiedzialność pliku: terminy słownika modułu Translate, wraz z ich zapisem w jednej transakcji przez wstawienie albo aktualizację po kodzie.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// TerminSlownika to wiersz tabeli termin_slownika, odpowiednik terminu stosowany przy generowaniu tłumaczeń, wraz z jego wartościami sterującymi.
type TerminSlownika struct {
	ID           int64
	Kod          string
	Zrodlo       string
	Jezyk        string
	Cel          *string
	NieTlumaczyc bool
	Uwaga        *string
	// Stan i Dziedzina bywają puste; termin zastany nikogo o stan nie pytał wcześniej.
	Stan           *string
	Dziedzina      *string
	Zaktualizowano int64
}

const (
	kolumnyTerminuSlownika = `id, identyfikator_zewnetrzny, zrodlo, jezyk, cel,
	                          nie_tlumaczyc, uwaga, stan, dziedzina, zaktualizowano`

	zapiszTerminSlownika = `INSERT INTO termin_slownika
	                        (identyfikator_zewnetrzny, zrodlo, jezyk, cel,
	                         nie_tlumaczyc, uwaga, stan, dziedzina, zaktualizowano, konto_id)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                        ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                            zrodlo = excluded.zrodlo,
	                            jezyk = excluded.jezyk,
	                            cel = excluded.cel,
	                            nie_tlumaczyc = excluded.nie_tlumaczyc,
	                            uwaga = excluded.uwaga,
	                            stan = excluded.stan,
	                            dziedzina = excluded.dziedzina,
	                            zaktualizowano = excluded.zaktualizowano
	                        WHERE ` + WarunekKonta

	pobierzTerminSlownika = `SELECT ` + kolumnyTerminuSlownika + ` FROM termin_slownika
	                         WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaTerminowSlownika = `SELECT ` + kolumnyTerminuSlownika + ` FROM termin_slownika
	                         WHERE ` + WarunekKonta + `
	                         ORDER BY jezyk, zrodlo`
)

// ZapiszTerminy zapisuje cały nadesłany wykaz terminów w jednej transakcji i zwraca terminy odczytane z bazy po zapisie, nie przepisane żądanie.
func (r *repozytoriumTlumaczen) ZapiszTerminy(ctx context.Context,
	terminy []TerminSlownika) ([]TerminSlownika, error) {

	kody := make([]string, 0, len(terminy))
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszTerminSlownika)
		if err != nil {
			return err
		}
		teraz := time.Now().UnixMilli()
		for _, termin := range terminy {
			if termin.Kod == "" {
				return fmt.Errorf("dane: termin słownika bez identyfikatora")
			}
			if termin.Zrodlo == "" {
				return fmt.Errorf("dane: termin słownika %q bez treści źródłowej", termin.Kod)
			}
			if termin.Jezyk == "" {
				return fmt.Errorf("dane: termin słownika %q bez języka odpowiednika", termin.Kod)
			}
			zaktualizowano := termin.Zaktualizowano
			if zaktualizowano == 0 {
				zaktualizowano = teraz
			}
			_, err := zapis.ExecContext(ctx, termin.Kod, termin.Zrodlo, termin.Jezyk,
				tekstDoKolumny(termin.Cel), liczbaLogiczna(termin.NieTlumaczyc),
				tekstDoKolumny(termin.Uwaga), tekstDoKolumny(termin.Stan),
				tekstDoKolumny(termin.Dziedzina), zaktualizowano,
				KontoOperatora(ctx), KontoOperatora(ctx))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać terminu słownika %q: %w", termin.Kod, err)
			}
			kody = append(kody, termin.Kod)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	zapisane := make([]TerminSlownika, 0, len(kody))
	for _, kod := range kody {
		termin, err := r.Termin(ctx, kod)
		if err != nil {
			return nil, err
		}
		zapisane = append(zapisane, termin)
	}
	return zapisane, nil
}

// Terminy zwraca cały słownik po języku i treści źródłowej; zasila `glossary.apply` i `glossary.occurrences`.
func (r *repozytoriumTlumaczen) Terminy(ctx context.Context) ([]TerminSlownika, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaTerminowSlownika)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać terminów słownika: %w", err)
	}
	defer wiersze.Close()

	lista := []TerminSlownika{}
	for wiersze.Next() {
		termin, err := odczytajTerminSlownika(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz terminu słownika: %w", err)
		}
		lista = append(lista, termin)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt terminów słownika: %w", err)
	}
	return lista, nil
}

// Termin zwraca termin słownika po kodzie; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumTlumaczen) Termin(ctx context.Context, kod string) (TerminSlownika, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzTerminSlownika)
	if err != nil {
		return TerminSlownika{}, err
	}
	termin, err := odczytajTerminSlownika(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return TerminSlownika{}, ErrBrakWiersza
	}
	if err != nil {
		return TerminSlownika{}, fmt.Errorf("dane: nieczytelny wiersz terminu słownika %q: %w", kod, err)
	}
	return termin, nil
}

// odczytajTerminSlownika składa strukturę terminu z jednego wiersza wyniku zapytania, kolumna po kolumnie.
func odczytajTerminSlownika(wiersz skaner) (TerminSlownika, error) {
	var termin TerminSlownika
	var cel, uwaga, stan, dziedzina sql.NullString
	var nieTlumaczyc int
	err := wiersz.Scan(&termin.ID, &termin.Kod, &termin.Zrodlo, &termin.Jezyk, &cel,
		&nieTlumaczyc, &uwaga, &stan, &dziedzina, &termin.Zaktualizowano)
	if err != nil {
		return TerminSlownika{}, err
	}
	termin.Cel = tekstZKolumny(cel)
	termin.NieTlumaczyc = nieTlumaczyc != 0
	termin.Uwaga = tekstZKolumny(uwaga)
	termin.Stan = tekstZKolumny(stan)
	termin.Dziedzina = tekstZKolumny(dziedzina)
	return termin, nil
}
