// Zestawy żetonów systemu projektowego (zestaw_zetonow_design, zeton_design); zapis zestawu jest zawsze pełny.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ZestawZetonowDesignu to wiersz `zestaw_zetonow_design`; żetony wchodzą osobnym odczytem, `Liczba` je liczy.
type ZestawZetonowDesignu struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	Motyw          *string
	Liczba         int
	Zaktualizowano string
}

// ZetonDesignu to wiersz `zeton_design` — jedna rola systemu projektowego wraz z jej wartością.
type ZetonDesignu struct {
	Nazwa      string
	Rodzaj     string
	Wartosc    string
	Opis       *string
	OdsylaczDo *string
	Kolejnosc  int
}

const (
	kolumnyZestawuZetonowDesignu = `z.id, z.identyfikator_zewnetrzny, z.okno, z.nazwa, z.motyw,
	                                (SELECT COUNT(*) FROM zeton_design t WHERE t.zestaw_id = z.id),
	                                z.zaktualizowano`

	zapiszZestawZetonowDesignu = `INSERT INTO zestaw_zetonow_design
	                              (identyfikator_zewnetrzny, okno, nazwa, motyw, zaktualizowano, konto_id)
	                              VALUES (?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'), ` + WskazanieKonta + `)
	                              ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                  nazwa = excluded.nazwa,
	                                  motyw = excluded.motyw,
	                                  zaktualizowano = excluded.zaktualizowano
	                              WHERE ` + WarunekKonta

	pobierzZestawZetonowDesignu = `SELECT ` + kolumnyZestawuZetonowDesignu +
		` FROM zestaw_zetonow_design z WHERE z.identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaZestawowZetonowDesignu = `SELECT ` + kolumnyZestawuZetonowDesignu +
		` FROM zestaw_zetonow_design z WHERE z.okno = ? AND ` + WarunekKonta +
		`  ORDER BY z.zaktualizowano DESC, z.id DESC`

	listaZestawowZetonowDesignuJeden = `SELECT ` + kolumnyZestawuZetonowDesignu +
		` FROM zestaw_zetonow_design z WHERE z.okno = ? AND ` + WarunekKonta +
		` AND z.identyfikator_zewnetrzny = ?`

	usunZetonyZestawuDesignu = `DELETE FROM zeton_design WHERE zestaw_id = ?`

	wstawZetonDesignu = `INSERT INTO zeton_design
	                     (zestaw_id, nazwa, rodzaj, wartosc, opis, odsylacz_do, kolejnosc)
	                     VALUES (?, ?, ?, ?, ?, ?, ?)
	                     ON CONFLICT(zestaw_id, nazwa) DO UPDATE SET
	                         rodzaj = excluded.rodzaj,
	                         wartosc = excluded.wartosc,
	                         opis = excluded.opis,
	                         odsylacz_do = excluded.odsylacz_do,
	                         kolejnosc = excluded.kolejnosc`

	listaZetonowZestawuDesignu = `SELECT nazwa, rodzaj, wartosc, opis, odsylacz_do, kolejnosc
	                              FROM zeton_design WHERE zestaw_id = ? ORDER BY kolejnosc, nazwa`
)

// ZapiszZestawZetonowDesignu zakłada zestaw albo nadpisuje zastany po kodzie i podmienia komplet żetonów w transakcji.
func (r *repozytoriumDesignu) ZapiszZestawZetonowDesignu(ctx context.Context,
	zestaw ZestawZetonowDesignu, zetony []ZetonDesignu) (ZestawZetonowDesignu, error) {

	if zestaw.Kod == "" {
		return ZestawZetonowDesignu{}, fmt.Errorf("dane: zestaw żetonów design bez identyfikatora")
	}
	if zestaw.Okno == "" {
		return ZestawZetonowDesignu{}, fmt.Errorf("dane: zestaw żetonów design %q bez okna", zestaw.Kod)
	}
	if zestaw.Nazwa == "" {
		return ZestawZetonowDesignu{}, fmt.Errorf("dane: zestaw żetonów design %q bez nazwy", zestaw.Kod)
	}

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszZestawZetonowDesignu)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, zestaw.Kod, zestaw.Okno, zestaw.Nazwa,
			tekstDoKolumny(zestaw.Motyw), KontoOperatora(ctx), KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można zapisać zestawu żetonów design %q: %w", zestaw.Kod, err)
		}

		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzZestawZetonowDesignu)
		if err != nil {
			return err
		}
		zapisany, err := odczytajZestawZetonowDesignu(odczyt.QueryRowContext(ctx, zestaw.Kod, KontoOperatora(ctx)))
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać zapisanego zestawu żetonów design %q: %w",
				zestaw.Kod, err)
		}

		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunZetonyZestawuDesignu)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, zapisany.ID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić żetonów zestawu design %q: %w", zestaw.Kod, err)
		}

		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawZetonDesignu)
		if err != nil {
			return err
		}
		for numer, zeton := range zetony {
			if zeton.Nazwa == "" {
				return fmt.Errorf("dane: żeton numer %d zestawu design %q bez nazwy roli",
					numer, zestaw.Kod)
			}
			kolejnosc := zeton.Kolejnosc
			if kolejnosc == 0 {
				kolejnosc = numer + 1
			}
			_, err := wstawienie.ExecContext(ctx, zapisany.ID, zeton.Nazwa, zeton.Rodzaj,
				zeton.Wartosc, tekstDoKolumny(zeton.Opis), tekstDoKolumny(zeton.OdsylaczDo), kolejnosc)
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać żetonu %q zestawu design %q: %w",
					zeton.Nazwa, zestaw.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return ZestawZetonowDesignu{}, err
	}
	return r.ZestawZetonowDesignuPoKodzie(ctx, zestaw.Kod)
}

// ZestawZetonowDesignuPoKodzie zwraca zestaw po kodzie; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumDesignu) ZestawZetonowDesignuPoKodzie(ctx context.Context,
	kod string) (ZestawZetonowDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZestawZetonowDesignu)
	if err != nil {
		return ZestawZetonowDesignu{}, err
	}
	zestaw, err := odczytajZestawZetonowDesignu(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return ZestawZetonowDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return ZestawZetonowDesignu{}, fmt.Errorf(
			"dane: nieczytelny wiersz zestawu żetonów design %q: %w", kod, err)
	}
	return zestaw, nil
}

// ZestawyZetonowDesignu zwraca zestawy okna; wskazanie kodu zawęża wykaz do jednego zestawu.
func (r *repozytoriumDesignu) ZestawyZetonowDesignu(ctx context.Context,
	okno string, kod *string) ([]ZestawZetonowDesignu, error) {

	zapytanie := listaZestawowZetonowDesignu
	argumenty := []any{okno, KontoOperatora(ctx)}
	if kod != nil && *kod != "" {
		zapytanie = listaZestawowZetonowDesignuJeden
		argumenty = append(argumenty, *kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zestawów żetonów design okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []ZestawZetonowDesignu{}
	for wiersze.Next() {
		zestaw, err := odczytajZestawZetonowDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zestawu żetonów design okna %q: %w", okno, err)
		}
		lista = append(lista, zestaw)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zestawów żetonów design okna %q: %w", okno, err)
	}
	return lista, nil
}

// ZetonyZestawuDesignu zwraca żetony zestawu wskazanego kluczem wiersza
// w zapisanej kolejności, odczytane z tabeli zeton_design.
func (r *repozytoriumDesignu) ZetonyZestawuDesignu(ctx context.Context,
	zestawID int64) ([]ZetonDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZetonowZestawuDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, zestawID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać żetonów zestawu design %d: %w", zestawID, err)
	}
	defer wiersze.Close()

	lista := []ZetonDesignu{}
	for wiersze.Next() {
		var zeton ZetonDesignu
		var opis, odsylacz sql.NullString
		if err := wiersze.Scan(&zeton.Nazwa, &zeton.Rodzaj, &zeton.Wartosc, &opis,
			&odsylacz, &zeton.Kolejnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz żetonu zestawu design %d: %w", zestawID, err)
		}
		zeton.Opis = tekstZKolumny(opis)
		zeton.OdsylaczDo = tekstZKolumny(odsylacz)
		lista = append(lista, zeton)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt żetonów zestawu design %d: %w", zestawID, err)
	}
	return lista, nil
}

// odczytajZestawZetonowDesignu składa strukturę ZestawZetonowDesignu z jednego
// wiersza wyniku zapytania, w tym pole motywu dopuszczające wartość pustą.
func odczytajZestawZetonowDesignu(wiersz skaner) (ZestawZetonowDesignu, error) {
	var zestaw ZestawZetonowDesignu
	var motyw sql.NullString
	err := wiersz.Scan(&zestaw.ID, &zestaw.Kod, &zestaw.Okno, &zestaw.Nazwa, &motyw,
		&zestaw.Liczba, &zestaw.Zaktualizowano)
	if err != nil {
		return ZestawZetonowDesignu{}, err
	}
	zestaw.Motyw = tekstZKolumny(motyw)
	return zestaw, nil
}
