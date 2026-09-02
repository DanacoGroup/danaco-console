// Szablony materiału modułu Design (szablon_materialu_design i jego warstwy oraz strony); prompty są osobnym bytem.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// SzablonMaterialuDesignu opisuje jeden szablon materiału przypisany do okna, wraz z rozmiarem i opisem.
type SzablonMaterialuDesignu struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          string
	Rodzaj         string
	Szerokosc      float64
	Wysokosc       float64
	Opis           *string
	Zaktualizowano string
}

const (
	kolumnySzablonuMaterialuDesignu = `id, identyfikator_zewnetrzny, okno, nazwa, rodzaj,
	                                   szerokosc, wysokosc, opis, zaktualizowano`

	zapiszSzablonMaterialuDesignu = `INSERT INTO szablon_materialu_design
	                                 (identyfikator_zewnetrzny, okno, nazwa, rodzaj, szerokosc,
	                                  wysokosc, opis, zaktualizowano, konto_id)
	                                 VALUES (?, ?, ?, ?, ?, ?, ?,
	                                         strftime('%Y-%m-%dT%H:%M:%fZ','now'), ` + WskazanieKonta + `)
	                                 ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                     nazwa = excluded.nazwa,
	                                     rodzaj = excluded.rodzaj,
	                                     szerokosc = excluded.szerokosc,
	                                     wysokosc = excluded.wysokosc,
	                                     opis = excluded.opis,
	                                     zaktualizowano = excluded.zaktualizowano
	                                 WHERE ` + WarunekKonta

	pobierzSzablonMaterialuDesignu = `SELECT ` + kolumnySzablonuMaterialuDesignu +
		` FROM szablon_materialu_design WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaSzablonowMaterialuDesignu = `SELECT ` + kolumnySzablonuMaterialuDesignu +
		` FROM szablon_materialu_design WHERE okno = ? AND (? = '' OR rodzaj = ?) AND ` + WarunekKonta +
		`  ORDER BY zaktualizowano DESC, id DESC`

	usunWarstwySzablonuMaterialuDesignu = `DELETE FROM warstwa_szablonu_materialu_design
	                                       WHERE szablon_id = ?`

	wstawWarstweSzablonuMaterialuDesignu = `INSERT INTO warstwa_szablonu_materialu_design
	                                        (identyfikator_zewnetrzny, szablon_id, zasob_id, x, y,
	                                         szerokosc, wysokosc, kolejnosc, zablokowana, adnotacja)
	                                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	listaWarstwSzablonuMaterialuDesignu = `SELECT id, identyfikator_zewnetrzny, szablon_id, zasob_id,
	                                              x, y, szerokosc, wysokosc, kolejnosc, zablokowana,
	                                              adnotacja
	                                       FROM warstwa_szablonu_materialu_design
	                                       WHERE szablon_id = ? ORDER BY kolejnosc, id`
)

// ZapiszSzablonMaterialuDesignu zakłada szablon albo nadpisuje zastany po kodzie i podmienia komplet jego warstw.
func (r *repozytoriumDesignu) ZapiszSzablonMaterialuDesignu(ctx context.Context,
	szablon SzablonMaterialuDesignu, warstwy []WarstwaKompozycji) (SzablonMaterialuDesignu, error) {

	if szablon.Kod == "" {
		return SzablonMaterialuDesignu{}, fmt.Errorf("dane: szablon materiału design bez identyfikatora")
	}
	if szablon.Okno == "" {
		return SzablonMaterialuDesignu{}, fmt.Errorf(
			"dane: szablon materiału design %q bez okna", szablon.Kod)
	}

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszSzablonMaterialuDesignu)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, szablon.Kod, szablon.Okno, szablon.Nazwa,
			szablon.Rodzaj, szablon.Szerokosc, szablon.Wysokosc,
			tekstDoKolumny(szablon.Opis), KontoOperatora(ctx), KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można zapisać szablonu materiału design %q: %w",
				szablon.Kod, err)
		}

		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzSzablonMaterialuDesignu)
		if err != nil {
			return err
		}
		zapisany, err := odczytajSzablonMaterialuDesignu(odczyt.QueryRowContext(ctx, szablon.Kod, KontoOperatora(ctx)))
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać zapisanego szablonu materiału design %q: %w",
				szablon.Kod, err)
		}

		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunWarstwySzablonuMaterialuDesignu)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, zapisany.ID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić warstw szablonu materiału design %q: %w",
				szablon.Kod, err)
		}

		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWarstweSzablonuMaterialuDesignu)
		if err != nil {
			return err
		}
		for numer, warstwa := range warstwy {
			if warstwa.Kod == "" {
				return fmt.Errorf("dane: warstwa numer %d szablonu materiału design %q bez identyfikatora",
					numer, szablon.Kod)
			}
			kolejnosc := warstwa.Kolejnosc
			if kolejnosc == 0 {
				kolejnosc = numer + 1
			}
			_, err := wstawienie.ExecContext(ctx, warstwa.Kod, zapisany.ID,
				tekstDoKolumny(warstwa.ZasobID), warstwa.X, warstwa.Y,
				warstwa.Szerokosc, warstwa.Wysokosc, kolejnosc,
				liczbaLogiczna(warstwa.Zablokowana), tekstDoKolumny(warstwa.Adnotacja))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać warstwy %q szablonu materiału design %q: %w",
					warstwa.Kod, szablon.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return SzablonMaterialuDesignu{}, err
	}
	return r.SzablonMaterialuDesignuPoKodzie(ctx, szablon.Kod)
}

// SzablonMaterialuDesignuPoKodzie zwraca szablon po kodzie; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumDesignu) SzablonMaterialuDesignuPoKodzie(ctx context.Context,
	kod string) (SzablonMaterialuDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSzablonMaterialuDesignu)
	if err != nil {
		return SzablonMaterialuDesignu{}, err
	}
	szablon, err := odczytajSzablonMaterialuDesignu(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return SzablonMaterialuDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return SzablonMaterialuDesignu{}, fmt.Errorf(
			"dane: nieczytelny wiersz szablonu materiału design %q: %w", kod, err)
	}
	return szablon, nil
}

// SzablonyMaterialuDesignu zwraca szablony okna; rodzaj pusty znaczy „wszystkie rodzaje".
func (r *repozytoriumDesignu) SzablonyMaterialuDesignu(ctx context.Context,
	okno, rodzaj string) ([]SzablonMaterialuDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaSzablonowMaterialuDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, rodzaj, rodzaj, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać szablonów materiału design okna %q: %w",
			okno, err)
	}
	defer wiersze.Close()

	lista := []SzablonMaterialuDesignu{}
	for wiersze.Next() {
		szablon, err := odczytajSzablonMaterialuDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz szablonu materiału design okna %q: %w",
				okno, err)
		}
		lista = append(lista, szablon)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt szablonów materiału design okna %q: %w",
			okno, err)
	}
	return lista, nil
}

// WarstwySzablonuMaterialuDesignu zwraca warstwy szablonu w kolejności wyrysu.
func (r *repozytoriumDesignu) WarstwySzablonuMaterialuDesignu(ctx context.Context,
	szablonID int64) ([]WarstwaKompozycji, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWarstwSzablonuMaterialuDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, szablonID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać warstw szablonu materiału design %d: %w",
			szablonID, err)
	}
	defer wiersze.Close()

	lista := []WarstwaKompozycji{}
	for wiersze.Next() {
		var warstwa WarstwaKompozycji
		var zasobID, adnotacja sql.NullString
		var x, y, szerokosc, wysokosc sql.NullFloat64
		var zablokowana int
		err := wiersze.Scan(&warstwa.ID, &warstwa.Kod, &warstwa.KompozycjaID, &zasobID,
			&x, &y, &szerokosc, &wysokosc, &warstwa.Kolejnosc, &zablokowana, &adnotacja)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz warstwy szablonu materiału design %d: %w",
				szablonID, err)
		}
		warstwa.ZasobID = tekstZKolumny(zasobID)
		warstwa.Adnotacja = tekstZKolumny(adnotacja)
		warstwa.X = liczbaRzeczywistaZKolumny(x)
		warstwa.Y = liczbaRzeczywistaZKolumny(y)
		warstwa.Szerokosc = liczbaRzeczywistaZKolumny(szerokosc)
		warstwa.Wysokosc = liczbaRzeczywistaZKolumny(wysokosc)
		warstwa.Zablokowana = zablokowana == 1
		lista = append(lista, warstwa)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt warstw szablonu materiału design %d: %w",
			szablonID, err)
	}
	return lista, nil
}

func odczytajSzablonMaterialuDesignu(wiersz skaner) (SzablonMaterialuDesignu, error) {
	var szablon SzablonMaterialuDesignu
	var opis sql.NullString
	err := wiersz.Scan(&szablon.ID, &szablon.Kod, &szablon.Okno, &szablon.Nazwa,
		&szablon.Rodzaj, &szablon.Szerokosc, &szablon.Wysokosc, &opis, &szablon.Zaktualizowano)
	if err != nil {
		return SzablonMaterialuDesignu{}, err
	}
	szablon.Opis = tekstZKolumny(opis)
	return szablon, nil
}

// StronaSzablonuMaterialuDesignu to jedna strona publikacji wielostronicowej; szablon bez stron jest jednostronicowy.
type StronaSzablonuMaterialuDesignu struct {
	ID        int64
	Kod       string
	SzablonID int64
	Numer     int
	Nazwa     *string
}

const (
	kolumnyStronySzablonuDesignu = `id, identyfikator_zewnetrzny, szablon_id, numer, nazwa`

	usunStronySzablonuDesignu = `DELETE FROM strona_szablonu_materialu_design WHERE szablon_id = ?`

	wstawStroneSzablonuDesignu = `INSERT INTO strona_szablonu_materialu_design
	                              (identyfikator_zewnetrzny, szablon_id, numer, nazwa)
	                              VALUES (?, ?, ?, ?)`

	listaStronSzablonuDesignu = `SELECT ` + kolumnyStronySzablonuDesignu +
		` FROM strona_szablonu_materialu_design WHERE szablon_id = ? ORDER BY numer, id`

	wstawWarstweStronySzablonuDesignu = `INSERT INTO warstwa_strony_szablonu_design
	                                     (identyfikator_zewnetrzny, strona_id, zasob_id, x, y,
	                                      szerokosc, wysokosc, kolejnosc, zablokowana, adnotacja)
	                                     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Pusta wartość `utworzono` dopasowuje liczbę kolumn do wspólnego odczytywacza warstw kompozycji.
	listaWarstwStronySzablonuDesignu = `SELECT id, identyfikator_zewnetrzny, strona_id, zasob_id,
	                                           x, y, szerokosc, wysokosc, kolejnosc, zablokowana,
	                                           adnotacja, '' AS utworzono
	                                    FROM warstwa_strony_szablonu_design
	                                    WHERE strona_id = ? ORDER BY kolejnosc, id`
)

// ZapiszStronySzablonuMaterialuDesignu podmienia komplet stron szablonu wraz z ich warstwami w jednej transakcji.
func (r *repozytoriumDesignu) ZapiszStronySzablonuMaterialuDesignu(ctx context.Context,
	szablonID int64, strony []StronaSzablonuMaterialuDesignu,
	warstwy map[string][]WarstwaKompozycji) error {

	if szablonID == 0 {
		return fmt.Errorf("dane: strony szablonu materiału design bez szablonu")
	}
	transakcja, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("dane: nie można otworzyć transakcji stron szablonu design: %w", err)
	}
	defer func() { _ = transakcja.Rollback() }()

	if _, err := transakcja.ExecContext(ctx, usunStronySzablonuDesignu, szablonID); err != nil {
		return fmt.Errorf("dane: nie można usunąć stron szablonu design %d: %w", szablonID, err)
	}
	for _, strona := range strony {
		if strona.Kod == "" {
			return fmt.Errorf("dane: strona szablonu design bez identyfikatora")
		}
		wynik, err := transakcja.ExecContext(ctx, wstawStroneSzablonuDesignu,
			strona.Kod, szablonID, strona.Numer, tekstDoKolumny(strona.Nazwa))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać strony szablonu design %q: %w",
				strona.Kod, err)
		}
		stronaID, err := wynik.LastInsertId()
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać klucza strony szablonu design %q: %w",
				strona.Kod, err)
		}
		for _, warstwa := range warstwy[strona.Kod] {
			if _, err := transakcja.ExecContext(ctx, wstawWarstweStronySzablonuDesignu,
				warstwa.Kod, stronaID, tekstDoKolumny(warstwa.ZasobID),
				ulamekDoKolumnyDesignu(warstwa.X), ulamekDoKolumnyDesignu(warstwa.Y),
				ulamekDoKolumnyDesignu(warstwa.Szerokosc), ulamekDoKolumnyDesignu(warstwa.Wysokosc),
				warstwa.Kolejnosc, warstwa.Zablokowana,
				tekstDoKolumny(warstwa.Adnotacja)); err != nil {

				return fmt.Errorf("dane: nie można zapisać warstwy %q strony %q: %w",
					warstwa.Kod, strona.Kod, err)
			}
		}
	}
	if err := transakcja.Commit(); err != nil {
		return fmt.Errorf("dane: nie można zamknąć transakcji stron szablonu design: %w", err)
	}
	return nil
}

// StronySzablonuMaterialuDesignu oddaje strony szablonu w kolejności numerów,
// zgodnie z porządkiem publikacji wielostronicowej.
func (r *repozytoriumDesignu) StronySzablonuMaterialuDesignu(ctx context.Context,
	szablonID int64) ([]StronaSzablonuMaterialuDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaStronSzablonuDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, szablonID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać stron szablonu design %d: %w",
			szablonID, err)
	}
	defer wiersze.Close()

	strony := []StronaSzablonuMaterialuDesignu{}
	for wiersze.Next() {
		var strona StronaSzablonuMaterialuDesignu
		var nazwa sql.NullString
		if err := wiersze.Scan(&strona.ID, &strona.Kod, &strona.SzablonID,
			&strona.Numer, &nazwa); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz strony szablonu design: %w", err)
		}
		strona.Nazwa = tekstZKolumny(nazwa)
		strony = append(strony, strona)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt stron szablonu design: %w", err)
	}
	return strony, nil
}

// WarstwyStronySzablonuMaterialuDesignu oddaje warstwy jednej strony w kolejności wyrysu.
func (r *repozytoriumDesignu) WarstwyStronySzablonuMaterialuDesignu(ctx context.Context,
	stronaID int64) ([]WarstwaKompozycji, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWarstwStronySzablonuDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, stronaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać warstw strony szablonu design %d: %w",
			stronaID, err)
	}
	defer wiersze.Close()

	warstwy := []WarstwaKompozycji{}
	for wiersze.Next() {
		warstwa, err := odczytajWarstweKompozycjiDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz warstwy strony szablonu design: %w", err)
		}
		warstwy = append(warstwy, warstwa)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt warstw strony szablonu design: %w", err)
	}
	return warstwy, nil
}
