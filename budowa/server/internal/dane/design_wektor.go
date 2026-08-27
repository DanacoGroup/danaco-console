// Obszar grafiki wektorowej modułu Design: ścieżki, symbole i ich członkowie,
// zapisywani zawsze pełnym stanem zamiast różnicy.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// SciezkaWektorowaDesignu to wiersz tabeli sciezka_wektorowa_design; WezlyJSON,
// WypelnienieJSON i ObrysJSON niosą zapis kontraktu bez rozkładania.
type SciezkaWektorowaDesignu struct {
	ID              int64
	Kod             string
	KompozycjaID    int64
	WarstwaKod      *string
	Nazwa           *string
	WezlyJSON       string
	Zamknieta       bool
	WypelnienieJSON *string
	ObrysJSON       *string
	Kolejnosc       int
	Zaktualizowano  string
}

// SymbolDesignu to wiersz tabeli `symbol_design`. Członkowie leżą osobno
// i wchodzą osobnym odczytem; `Liczba` jest ich licznikiem — a zarazem liczbą
// miejsc, w których symbol stoi na planszy.
type SymbolDesignu struct {
	ID             int64
	Kod            string
	KompozycjaID   int64
	Nazwa          string
	Liczba         int
	Zaktualizowano string
}

// CzlonekSymbolyDesignu to jeden element składowy symbolu: ścieżka albo
// warstwa, wskazana identyfikatorem zewnętrznym.
type CzlonekSymbolyDesignu struct {
	Rodzaj    string
	Kod       string
	Kolejnosc int
}

// Rodzaje członków symbolu — wartości kolumny `czlonek_symbolu_design.rodzaj`,
// pilnowane warunkiem CHECK schematu.
const (
	CzlonekSymboluSciezka = "sciezka"
	CzlonekSymboluWarstwa = "warstwa"
)

const (
	kolumnySciezkiWektorowejDesignu = `id, identyfikator_zewnetrzny, kompozycja_id, warstwa_kod,
	                                   nazwa, wezly_json, zamknieta, wypelnienie_json, obrys_json,
	                                   kolejnosc, zaktualizowano`

	zapiszSciezkeWektorowaDesignu = `INSERT INTO sciezka_wektorowa_design
	                                 (identyfikator_zewnetrzny, kompozycja_id, warstwa_kod, nazwa,
	                                  wezly_json, zamknieta, wypelnienie_json, obrys_json, kolejnosc,
	                                  zaktualizowano)
	                                 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,
	                                         strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                                 ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                     warstwa_kod = excluded.warstwa_kod,
	                                     nazwa = excluded.nazwa,
	                                     wezly_json = excluded.wezly_json,
	                                     zamknieta = excluded.zamknieta,
	                                     wypelnienie_json = excluded.wypelnienie_json,
	                                     obrys_json = excluded.obrys_json,
	                                     kolejnosc = excluded.kolejnosc,
	                                     zaktualizowano = excluded.zaktualizowano`

	pobierzSciezkeWektorowaDesignu = `SELECT ` + kolumnySciezkiWektorowejDesignu +
		` FROM sciezka_wektorowa_design WHERE identyfikator_zewnetrzny = ?`

	listaSciezekWektorowychDesignu = `SELECT ` + kolumnySciezkiWektorowejDesignu +
		` FROM sciezka_wektorowa_design WHERE kompozycja_id = ? ORDER BY kolejnosc, id`

	listaSciezekWektorowychDesignuWarstwy = `SELECT ` + kolumnySciezkiWektorowejDesignu +
		` FROM sciezka_wektorowa_design WHERE kompozycja_id = ? AND warstwa_kod = ?
		  ORDER BY kolejnosc, id`

	usunSciezkeWektorowaDesignu = `DELETE FROM sciezka_wektorowa_design
	                               WHERE identyfikator_zewnetrzny = ?`

	// Kolejność nowej ścieżki bierze się z tego, co już leży na planszy: ścieżka
	// dorysowana staje NAD poprzednimi, bo tak działa rysowanie.
	najwyzszaKolejnoscSciezkiDesignu = `SELECT COALESCE(MAX(kolejnosc), 0)
	                                    FROM sciezka_wektorowa_design WHERE kompozycja_id = ?`

	kolumnySymbolyDesignu = `s.id, s.identyfikator_zewnetrzny, s.kompozycja_id, s.nazwa,
	                         (SELECT COUNT(*) FROM czlonek_symbolu_design c WHERE c.symbol_id = s.id),
	                         s.zaktualizowano`

	zapiszSymbolDesignu = `INSERT INTO symbol_design
	                       (identyfikator_zewnetrzny, kompozycja_id, nazwa, zaktualizowano)
	                       VALUES (?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                       ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                           nazwa = excluded.nazwa,
	                           zaktualizowano = excluded.zaktualizowano`

	pobierzSymbolDesignu = `SELECT ` + kolumnySymbolyDesignu +
		` FROM symbol_design s WHERE s.identyfikator_zewnetrzny = ?`

	listaSymboliDesignu = `SELECT ` + kolumnySymbolyDesignu +
		` FROM symbol_design s WHERE s.kompozycja_id = ? ORDER BY s.id`

	usunCzlonkowSymbolyDesignu = `DELETE FROM czlonek_symbolu_design WHERE symbol_id = ?`

	wstawCzlonkaSymbolyDesignu = `INSERT INTO czlonek_symbolu_design
	                              (symbol_id, rodzaj, czlonek_kod, kolejnosc)
	                              VALUES (?, ?, ?, ?)
	                              ON CONFLICT(symbol_id, rodzaj, czlonek_kod) DO UPDATE SET
	                                  kolejnosc = excluded.kolejnosc`

	listaCzlonkowSymbolyDesignu = `SELECT rodzaj, czlonek_kod, kolejnosc
	                               FROM czlonek_symbolu_design WHERE symbol_id = ?
	                               ORDER BY kolejnosc, id`
)

// ZapiszSciezkeWektorowaDesignu zakłada ścieżkę albo nadpisuje zastaną po
// identyfikatorze zewnętrznym. Kolejność zerowa znaczy „na wierzch" — wiersz
// dostaje wtedy numer o jeden wyższy od najwyższego na planszy.
func (r *repozytoriumDesignu) ZapiszSciezkeWektorowaDesignu(ctx context.Context,
	sciezka SciezkaWektorowaDesignu) (SciezkaWektorowaDesignu, error) {

	if sciezka.Kod == "" {
		return SciezkaWektorowaDesignu{}, fmt.Errorf("dane: ścieżka wektorowa design bez identyfikatora")
	}
	if sciezka.KompozycjaID == 0 {
		return SciezkaWektorowaDesignu{}, fmt.Errorf(
			"dane: ścieżka wektorowa design %q bez kompozycji", sciezka.Kod)
	}
	if sciezka.WezlyJSON == "" {
		return SciezkaWektorowaDesignu{}, fmt.Errorf(
			"dane: ścieżka wektorowa design %q bez węzłów", sciezka.Kod)
	}

	kolejnosc := sciezka.Kolejnosc
	if kolejnosc == 0 {
		najwyzsza, err := r.zapytania.przygotuj(ctx, najwyzszaKolejnoscSciezkiDesignu)
		if err != nil {
			return SciezkaWektorowaDesignu{}, err
		}
		var szczyt int
		if err := najwyzsza.QueryRowContext(ctx, sciezka.KompozycjaID).Scan(&szczyt); err != nil {
			return SciezkaWektorowaDesignu{}, fmt.Errorf(
				"dane: nie można ustalić kolejności ścieżki wektorowej design %q: %w", sciezka.Kod, err)
		}
		kolejnosc = szczyt + 1
	}

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSciezkeWektorowaDesignu)
	if err != nil {
		return SciezkaWektorowaDesignu{}, err
	}
	_, err = polecenie.ExecContext(ctx, sciezka.Kod, sciezka.KompozycjaID,
		tekstDoKolumny(sciezka.WarstwaKod), tekstDoKolumny(sciezka.Nazwa), sciezka.WezlyJSON,
		sciezka.Zamknieta, tekstDoKolumny(sciezka.WypelnienieJSON),
		tekstDoKolumny(sciezka.ObrysJSON), kolejnosc)
	if err != nil {
		return SciezkaWektorowaDesignu{}, fmt.Errorf(
			"dane: nie można zapisać ścieżki wektorowej design %q: %w", sciezka.Kod, err)
	}
	return r.SciezkaWektorowaDesignuPoKodzie(ctx, sciezka.Kod)
}

// SciezkaWektorowaDesignuPoKodzie zwraca ścieżkę o wskazanym identyfikatorze
// zewnętrznym. Brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumDesignu) SciezkaWektorowaDesignuPoKodzie(ctx context.Context,
	kod string) (SciezkaWektorowaDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSciezkeWektorowaDesignu)
	if err != nil {
		return SciezkaWektorowaDesignu{}, err
	}
	sciezka, err := odczytajSciezkeWektorowaDesignu(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return SciezkaWektorowaDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return SciezkaWektorowaDesignu{}, fmt.Errorf(
			"dane: nieczytelny wiersz ścieżki wektorowej design %q: %w", kod, err)
	}
	return sciezka, nil
}

// SciezkiWektoroweDesignu zwraca ścieżki kompozycji w kolejności rysowania.
// Wskazanie warstwy zawęża wykaz do jednej warstwy.
func (r *repozytoriumDesignu) SciezkiWektoroweDesignu(ctx context.Context,
	kompozycjaID int64, warstwaKod *string) ([]SciezkaWektorowaDesignu, error) {

	zapytanie := listaSciezekWektorowychDesignu
	argumenty := []any{kompozycjaID}
	if warstwaKod != nil && *warstwaKod != "" {
		zapytanie = listaSciezekWektorowychDesignuWarstwy
		argumenty = append(argumenty, *warstwaKod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf(
			"dane: nie można odczytać ścieżek wektorowych design kompozycji %d: %w", kompozycjaID, err)
	}
	defer wiersze.Close()

	lista := []SciezkaWektorowaDesignu{}
	for wiersze.Next() {
		sciezka, err := odczytajSciezkeWektorowaDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf(
				"dane: nieczytelny wiersz ścieżki wektorowej design kompozycji %d: %w", kompozycjaID, err)
		}
		lista = append(lista, sciezka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf(
			"dane: przerwany odczyt ścieżek wektorowych design kompozycji %d: %w", kompozycjaID, err)
	}
	return lista, nil
}

// UsunSciezkeWektorowaDesignu usuwa ścieżkę i mówi, czy wiersz istniał.
// Kontrakt (`DesignVectorPathRemoveResponse.Removed`) pyta o to wprost, a nie
// o powodzenie polecenia SQL, które usuwa zero wierszy równie pomyślnie jak
// jeden.
func (r *repozytoriumDesignu) UsunSciezkeWektorowaDesignu(ctx context.Context,
	kod string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, usunSciezkeWektorowaDesignu)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć ścieżki wektorowej design %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf(
			"dane: nie można ustalić skutku usunięcia ścieżki wektorowej design %q: %w", kod, err)
	}
	return zmienione > 0, nil
}

// ZapiszSymbolDesignu zakłada symbol albo nadpisuje zastany i podmienia komplet
// jego członków w jednej transakcji. Oddaje symbol po zapisie wraz z liczbą
// członków — to liczba miejsc, do których zmiana definicji doszła.
func (r *repozytoriumDesignu) ZapiszSymbolDesignu(ctx context.Context, symbol SymbolDesignu,
	czlonkowie []CzlonekSymbolyDesignu) (SymbolDesignu, error) {

	if symbol.Kod == "" {
		return SymbolDesignu{}, fmt.Errorf("dane: symbol design bez identyfikatora")
	}
	if symbol.KompozycjaID == 0 {
		return SymbolDesignu{}, fmt.Errorf("dane: symbol design %q bez kompozycji", symbol.Kod)
	}
	if symbol.Nazwa == "" {
		return SymbolDesignu{}, fmt.Errorf("dane: symbol design %q bez nazwy", symbol.Kod)
	}

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszSymbolDesignu)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, symbol.Kod, symbol.KompozycjaID, symbol.Nazwa); err != nil {
			return fmt.Errorf("dane: nie można zapisać symbolu design %q: %w", symbol.Kod, err)
		}

		// Symbol mógł powstać dopiero w tej transakcji: klucz wiersza czytamy przed podmianą jego członków.
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzSymbolDesignu)
		if err != nil {
			return err
		}
		zapisany, err := odczytajSymbolDesignu(odczyt.QueryRowContext(ctx, symbol.Kod))
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać zapisanego symbolu design %q: %w", symbol.Kod, err)
		}

		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunCzlonkowSymbolyDesignu)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, zapisany.ID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić członków symbolu design %q: %w", symbol.Kod, err)
		}

		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawCzlonkaSymbolyDesignu)
		if err != nil {
			return err
		}
		for numer, czlonek := range czlonkowie {
			if czlonek.Kod == "" {
				return fmt.Errorf("dane: członek numer %d symbolu design %q bez identyfikatora",
					numer, symbol.Kod)
			}
			kolejnosc := czlonek.Kolejnosc
			if kolejnosc == 0 {
				kolejnosc = numer + 1
			}
			if _, err := wstawienie.ExecContext(ctx, zapisany.ID, czlonek.Rodzaj,
				czlonek.Kod, kolejnosc); err != nil {
				return fmt.Errorf("dane: nie można zapisać członka %q symbolu design %q: %w",
					czlonek.Kod, symbol.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return SymbolDesignu{}, err
	}
	return r.SymbolDesignuPoKodzie(ctx, symbol.Kod)
}

// SymbolDesignuPoKodzie zwraca symbol o wskazanym identyfikatorze zewnętrznym, zapisany w tabeli symbol_design.
func (r *repozytoriumDesignu) SymbolDesignuPoKodzie(ctx context.Context,
	kod string) (SymbolDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSymbolDesignu)
	if err != nil {
		return SymbolDesignu{}, err
	}
	symbol, err := odczytajSymbolDesignu(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return SymbolDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return SymbolDesignu{}, fmt.Errorf("dane: nieczytelny wiersz symbolu design %q: %w", kod, err)
	}
	return symbol, nil
}

// SymboleDesignu zwraca wykaz symboli kompozycji design wraz z liczbą ich członków zapisanych naprawdę.
func (r *repozytoriumDesignu) SymboleDesignu(ctx context.Context,
	kompozycjaID int64) ([]SymbolDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaSymboliDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kompozycjaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać symboli design kompozycji %d: %w",
			kompozycjaID, err)
	}
	defer wiersze.Close()

	lista := []SymbolDesignu{}
	for wiersze.Next() {
		symbol, err := odczytajSymbolDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz symbolu design kompozycji %d: %w",
				kompozycjaID, err)
		}
		lista = append(lista, symbol)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt symboli design kompozycji %d: %w",
			kompozycjaID, err)
	}
	return lista, nil
}

// CzlonkowieSymbolyDesignu zwraca elementy składowe symbolu design w zapisanej kolejności wykazu członków.
func (r *repozytoriumDesignu) CzlonkowieSymbolyDesignu(ctx context.Context,
	symbolID int64) ([]CzlonekSymbolyDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaCzlonkowSymbolyDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, symbolID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać członków symbolu design %d: %w", symbolID, err)
	}
	defer wiersze.Close()

	lista := []CzlonekSymbolyDesignu{}
	for wiersze.Next() {
		var czlonek CzlonekSymbolyDesignu
		if err := wiersze.Scan(&czlonek.Rodzaj, &czlonek.Kod, &czlonek.Kolejnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz członka symbolu design %d: %w", symbolID, err)
		}
		lista = append(lista, czlonek)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt członków symbolu design %d: %w", symbolID, err)
	}
	return lista, nil
}

// odczytajSciezkeWektorowaDesignu składa strukturę ścieżki wektorowej z jednego wiersza wyniku zapytania.
func odczytajSciezkeWektorowaDesignu(wiersz skaner) (SciezkaWektorowaDesignu, error) {
	var sciezka SciezkaWektorowaDesignu
	var warstwa, nazwa, wypelnienie, obrys sql.NullString
	err := wiersz.Scan(&sciezka.ID, &sciezka.Kod, &sciezka.KompozycjaID, &warstwa, &nazwa,
		&sciezka.WezlyJSON, &sciezka.Zamknieta, &wypelnienie, &obrys, &sciezka.Kolejnosc,
		&sciezka.Zaktualizowano)
	if err != nil {
		return SciezkaWektorowaDesignu{}, err
	}
	sciezka.WarstwaKod = tekstZKolumny(warstwa)
	sciezka.Nazwa = tekstZKolumny(nazwa)
	sciezka.WypelnienieJSON = tekstZKolumny(wypelnienie)
	sciezka.ObrysJSON = tekstZKolumny(obrys)
	return sciezka, nil
}

// odczytajSymbolDesignu składa strukturę symbolu design z jednego wiersza wyniku zapytania bazy danych.
func odczytajSymbolDesignu(wiersz skaner) (SymbolDesignu, error) {
	var symbol SymbolDesignu
	err := wiersz.Scan(&symbol.ID, &symbol.Kod, &symbol.KompozycjaID, &symbol.Nazwa,
		&symbol.Liczba, &symbol.Zaktualizowano)
	if err != nil {
		return SymbolDesignu{}, err
	}
	return symbol, nil
}
