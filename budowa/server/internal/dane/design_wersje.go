// Obszar nazwanych wersji kompozycji Design Board (tabele
// `wersja_kompozycji_design`, `warstwa_wersji_kompozycji_design`, migracja 233)
// — część `RepozytoriumDesignu` zadeklarowanego w `design.go`.
//
// Wersja jest migawką układu, nie odwołaniem do warstw żywych. Warstwy zapisuje
// się kolumna w kolumnę, bo `design.board.update` usuwa je i wstawia od nowa
// przy każdym zapisie — odwołanie wskazywałoby wtedy wiersze, których już nie ma,
// a wersja przestałaby opisywać cokolwiek dokładnie wtedy, gdy jest potrzebna.
//
// Wykaz wersji warstw NIE czyta (kontrakt: `versions` bez `layers`), stąd
// `liczba_warstw` utrwalona w wierszu wersji. Warstwy wchodzą wyłącznie przy
// przywróceniu, osobnym odczytem (`WarstwyWersjiKompozycjiDesignu`).
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WersjaKompozycjiDesignu to wiersz tabeli `wersja_kompozycji_design`.
// KompozycjaID jest kluczem wiersza kompozycji, nie jej identyfikatorem
// zewnętrznym — przekład na kod robi `KompozycjaDesignuPoKluczu`.
type WersjaKompozycjiDesignu struct {
	ID           int64
	Kod          string
	KompozycjaID int64
	Nazwa        *string
	Uzasadnienie *string
	LiczbaWarstw int
	Utworzono    string
}

const (
	kolumnyWersjiKompozycjiDesignu = `id, identyfikator_zewnetrzny, kompozycja_id, nazwa,
	                                  uzasadnienie, liczba_warstw, utworzono`

	zapiszWersjeKompozycjiDesignu = `INSERT INTO wersja_kompozycji_design
	                                 (identyfikator_zewnetrzny, kompozycja_id, nazwa,
	                                  uzasadnienie, liczba_warstw)
	                                 VALUES (?, ?, ?, ?, ?)`

	pobierzWersjeKompozycjiDesignu = `SELECT ` + kolumnyWersjiKompozycjiDesignu +
		` FROM wersja_kompozycji_design WHERE identyfikator_zewnetrzny = ?`

	listaWersjiKompozycjiDesignu = `SELECT ` + kolumnyWersjiKompozycjiDesignu +
		` FROM wersja_kompozycji_design WHERE kompozycja_id = ?
		  ORDER BY utworzono DESC, id DESC LIMIT ?`

	liczbaWersjiKompozycjiDesignu = `SELECT COUNT(*) FROM wersja_kompozycji_design
	                                 WHERE kompozycja_id = ?`

	wstawWarstweWersjiKompozycjiDesignu = `INSERT INTO warstwa_wersji_kompozycji_design
	                                       (wersja_id, identyfikator_warstwy, zasob_id, x, y,
	                                        szerokosc, wysokosc, kolejnosc, zablokowana, adnotacja)
	                                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	listaWarstwWersjiKompozycjiDesignu = `SELECT id, identyfikator_warstwy, zasob_id, x, y,
	                                             szerokosc, wysokosc, kolejnosc, zablokowana,
	                                             adnotacja
	                                      FROM warstwa_wersji_kompozycji_design
	                                      WHERE wersja_id = ? ORDER BY kolejnosc, id`

	pobierzKompozycjeDesignuPoKluczu = `SELECT ` + kolumnyKompozycjiDesign +
		` FROM kompozycja_design WHERE id = ?`
)

// ZapiszWersjeKompozycjiDesignu utrwala migawkę układu wraz z warstwami
// w jednej transakcji — wersja bez warstw, które miała opisać, nie jest wersją
// niczego, więc niepowodzenie zapisu warstwy cofa również wiersz wersji.
//
// Wersja jest zawsze nowa. Nadpisania nie ma i nie ma być: wersja to zapis
// stanu z konkretnej chwili, a nadpisanie oznaczałoby, że stan sprzed godziny
// właśnie się zmienił.
func (r *repozytoriumDesignu) ZapiszWersjeKompozycjiDesignu(ctx context.Context,
	wersja WersjaKompozycjiDesignu, warstwy []WarstwaKompozycji) (WersjaKompozycjiDesignu, error) {

	if wersja.Kod == "" {
		return WersjaKompozycjiDesignu{}, fmt.Errorf("dane: wersja kompozycji design bez identyfikatora")
	}
	if wersja.KompozycjaID == 0 {
		return WersjaKompozycjiDesignu{}, fmt.Errorf(
			"dane: wersja kompozycji design %q bez kompozycji", wersja.Kod)
	}

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszWersjeKompozycjiDesignu)
		if err != nil {
			return err
		}
		wynik, err := zapis.ExecContext(ctx, wersja.Kod, wersja.KompozycjaID,
			tekstDoKolumny(wersja.Nazwa), tekstDoKolumny(wersja.Uzasadnienie), len(warstwy))
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać wersji kompozycji design %q: %w", wersja.Kod, err)
		}
		klucz, err := wynik.LastInsertId()
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać klucza wersji kompozycji design %q: %w",
				wersja.Kod, err)
		}

		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWarstweWersjiKompozycjiDesignu)
		if err != nil {
			return err
		}
		for numer, warstwa := range warstwy {
			kolejnosc := warstwa.Kolejnosc
			if kolejnosc == 0 {
				kolejnosc = numer + 1
			}
			_, err := wstawienie.ExecContext(ctx, klucz, warstwa.Kod,
				tekstDoKolumny(warstwa.ZasobID), warstwa.X, warstwa.Y,
				warstwa.Szerokosc, warstwa.Wysokosc, kolejnosc,
				liczbaLogiczna(warstwa.Zablokowana), tekstDoKolumny(warstwa.Adnotacja))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać warstwy %q wersji kompozycji design %q: %w",
					warstwa.Kod, wersja.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return WersjaKompozycjiDesignu{}, err
	}
	return r.WersjaKompozycjiDesignuPoKodzie(ctx, wersja.Kod)
}

// WersjaKompozycjiDesignuPoKodzie zwraca wersję o wskazanym identyfikatorze
// zewnętrznym. Brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumDesignu) WersjaKompozycjiDesignuPoKodzie(ctx context.Context,
	kod string) (WersjaKompozycjiDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWersjeKompozycjiDesignu)
	if err != nil {
		return WersjaKompozycjiDesignu{}, err
	}
	wersja, err := odczytajWersjeKompozycjiDesignu(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return WersjaKompozycjiDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return WersjaKompozycjiDesignu{}, fmt.Errorf(
			"dane: nieczytelny wiersz wersji kompozycji design %q: %w", kod, err)
	}
	return wersja, nil
}

// WersjeKompozycjiDesignu zwraca wersje kompozycji od najświeższej wraz
// z liczbą wszystkich — wykaz bywa przycięty limitem, licznik nie.
func (r *repozytoriumDesignu) WersjeKompozycjiDesignu(ctx context.Context,
	kompozycjaID int64, limit int) ([]WersjaKompozycjiDesignu, int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWersjiKompozycjiDesignu)
	if err != nil {
		return nil, 0, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kompozycjaID, granicaWykazu(limit))
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać wersji kompozycji design %d: %w",
			kompozycjaID, err)
	}
	defer wiersze.Close()

	lista := []WersjaKompozycjiDesignu{}
	for wiersze.Next() {
		wersja, err := odczytajWersjeKompozycjiDesignu(wiersze)
		if err != nil {
			return nil, 0, fmt.Errorf("dane: nieczytelny wiersz wersji kompozycji design %d: %w",
				kompozycjaID, err)
		}
		lista = append(lista, wersja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt wersji kompozycji design %d: %w",
			kompozycjaID, err)
	}

	polecenieLiczby, err := r.zapytania.przygotuj(ctx, liczbaWersjiKompozycjiDesignu)
	if err != nil {
		return nil, 0, err
	}
	var razem int
	if err := polecenieLiczby.QueryRowContext(ctx, kompozycjaID).Scan(&razem); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć wersji kompozycji design %d: %w",
			kompozycjaID, err)
	}
	return lista, razem, nil
}

// WarstwyWersjiKompozycjiDesignu zwraca warstwy migawki w kolejności
// renderowania. Struktura jest ta sama co dla warstw żywych
// (`WarstwaKompozycji`), bo przywrócenie przepisuje je wprost.
func (r *repozytoriumDesignu) WarstwyWersjiKompozycjiDesignu(ctx context.Context,
	wersjaID int64) ([]WarstwaKompozycji, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWarstwWersjiKompozycjiDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, wersjaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać warstw wersji kompozycji design %d: %w",
			wersjaID, err)
	}
	defer wiersze.Close()

	lista := []WarstwaKompozycji{}
	for wiersze.Next() {
		var warstwa WarstwaKompozycji
		var zasobID, adnotacja sql.NullString
		var x, y, szerokosc, wysokosc sql.NullFloat64
		var zablokowana int
		err := wiersze.Scan(&warstwa.ID, &warstwa.Kod, &zasobID, &x, &y, &szerokosc,
			&wysokosc, &warstwa.Kolejnosc, &zablokowana, &adnotacja)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz warstwy wersji kompozycji design %d: %w",
				wersjaID, err)
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
		return nil, fmt.Errorf("dane: przerwany odczyt warstw wersji kompozycji design %d: %w",
			wersjaID, err)
	}
	return lista, nil
}

// KompozycjaDesignuPoKluczu zwraca kompozycję po kluczu wiersza — przekład
// potrzebny wersji, która wskazuje kompozycję kluczem obcym.
func (r *repozytoriumDesignu) KompozycjaDesignuPoKluczu(ctx context.Context,
	id int64) (KompozycjaDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKompozycjeDesignuPoKluczu)
	if err != nil {
		return KompozycjaDesignu{}, err
	}
	kompozycja, err := odczytajKompozycjeDesign(polecenie.QueryRowContext(ctx, id))
	if errors.Is(err, sql.ErrNoRows) {
		return KompozycjaDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return KompozycjaDesignu{}, fmt.Errorf("dane: nieczytelny wiersz kompozycji design %d: %w", id, err)
	}
	return kompozycja, nil
}

// odczytajWersjeKompozycjiDesignu składa strukturę z jednego wiersza wyniku.
func odczytajWersjeKompozycjiDesignu(wiersz skaner) (WersjaKompozycjiDesignu, error) {
	var wersja WersjaKompozycjiDesignu
	var nazwa, uzasadnienie sql.NullString
	err := wiersz.Scan(&wersja.ID, &wersja.Kod, &wersja.KompozycjaID, &nazwa,
		&uzasadnienie, &wersja.LiczbaWarstw, &wersja.Utworzono)
	if err != nil {
		return WersjaKompozycjiDesignu{}, err
	}
	wersja.Nazwa = tekstZKolumny(nazwa)
	wersja.Uzasadnienie = tekstZKolumny(uzasadnienie)
	return wersja, nil
}
