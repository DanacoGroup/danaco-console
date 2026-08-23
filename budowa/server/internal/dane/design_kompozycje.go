// Obszar kompozycji Design Board (tabele `kompozycja_design`
// i `warstwa_kompozycji_design`) — część `RepozytoriumDesignu` zadeklarowanego
// w `design.go`. Prompt strukturalny leży w `design.go`, zasoby i etykiety
// Assets Panel w `design_zasoby.go`.
//
// Zapis jest zawsze pełny. `design.board.update` nadsyła całą listę warstw na
// nowo — kontrakt (`DesignBoardUpdateRequest.Layers`) nie ma trybu częściowej
// zmiany. `ZapiszKompozycje` usuwa więc warstwy kompozycji i wstawia przysłany
// komplet od nowa w jednej transakcji, wzorem `ZapiszKroki`
// z `automations_kroki.go`; inaczej usunięcie warstwy w oknie nie usunęłoby jej
// w bazie i kompozycja rozeszłaby się z tym, co widać na ekranie.
//
// Brak `BoardId` zakłada kompozycję nową jedną ścieżką zapisu: wzorem
// `ZapiszPrompt` w `design.go` zapis idzie po identyfikatorze zewnętrznym
// z klauzulą ON CONFLICT, a brak `BoardId` rozstrzyga wywołujący, generując
// nowy identyfikator przed wywołaniem tej metody. Repozytorium zna wyłącznie
// tryb „załóż albo nadpisz".
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// KompozycjaDesignu to wiersz tabeli `kompozycja_design` — kompozycja Design
// Board bez warstw; warstwy leżą w osobnej tabeli, opisanej przez
// `WarstwaKompozycji`.
type KompozycjaDesignu struct {
	ID             int64
	Kod            string
	Okno           string
	Nazwa          *string
	Zaktualizowano string
}

// WarstwaKompozycji to wiersz tabeli `warstwa_kompozycji_design`. ZasobID jest
// identyfikatorem zewnętrznym zasobu (TEXT), nie kluczem obcym — warstwa może
// wskazywać zasób usunięty z Assets Panel po zapisie kompozycji.
type WarstwaKompozycji struct {
	ID           int64
	Kod          string
	KompozycjaID int64
	ZasobID      *string
	X            *float64
	Y            *float64
	Szerokosc    *float64
	Wysokosc     *float64
	Kolejnosc    int
	Zablokowana  bool
	Adnotacja    *string
	Utworzono    string
}

const (
	kolumnyKompozycjiDesign = `id, identyfikator_zewnetrzny, okno, nazwa, zaktualizowano`

	// Zapis zakłada kompozycję albo nadpisuje zastaną po identyfikatorze
	// zewnętrznym — brak `BoardId` w żądaniu (kontrakt: `DesignBoardUpdateRequest.BoardId`
	// opcjonalny) rozstrzyga wywołujący, nadając nowy identyfikator przed
	// wywołaniem, więc repozytorium ma jedną ścieżkę zapisu.
	zapiszKompozycjeDesign = `INSERT INTO kompozycja_design
	                          (identyfikator_zewnetrzny, okno, nazwa, zaktualizowano)
	                          VALUES (?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                          ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                              nazwa = excluded.nazwa,
	                              zaktualizowano = excluded.zaktualizowano`

	pobierzKompozycjeDesign = `SELECT ` + kolumnyKompozycjiDesign + ` FROM kompozycja_design
	                           WHERE identyfikator_zewnetrzny = ?`

	// Wykaz kompozycji okna (`design.board.list`). Porządek jest ten sam, co
	// w wykazie zasobów — od ostatnio zmienianej, żeby plansza porzucona
	// najpóźniej stała na górze listy. `id` rozstrzyga remis, bo
	// `zaktualizowano` ma rozdzielczość milisekundy i dwa zapisy z jednej pętli
	// okna potrafią w nią trafić razem.
	listaKompozycjiDesign = `SELECT ` + kolumnyKompozycjiDesign + ` FROM kompozycja_design
	                         WHERE okno = ? ORDER BY zaktualizowano DESC, id DESC`

	usunWarstwyKompozycjiDesign = `DELETE FROM warstwa_kompozycji_design WHERE kompozycja_id = ?`

	wstawWarstweKompozycjiDesign = `INSERT INTO warstwa_kompozycji_design
	                                (identyfikator_zewnetrzny, kompozycja_id, zasob_id, x, y,
	                                 szerokosc, wysokosc, kolejnosc, zablokowana, adnotacja)
	                                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	listaWarstwKompozycjiDesign = `SELECT id, identyfikator_zewnetrzny, kompozycja_id, zasob_id,
	                                       x, y, szerokosc, wysokosc, kolejnosc, zablokowana,
	                                       adnotacja, utworzono
	                               FROM warstwa_kompozycji_design WHERE kompozycja_id = ?
	                               ORDER BY kolejnosc, id`
)

// ZapiszKompozycje zakłada kompozycję albo nadpisuje zastaną po identyfikatorze
// zewnętrznym, po czym podmienia komplet jej warstw w jednej transakcji —
// wykaz pusty zostawia kompozycję bez warstw, co jest stanem poprawnym (płótno
// wyczyszczone w Design Board).
func (r *repozytoriumDesignu) ZapiszKompozycje(ctx context.Context,
	kompozycja KompozycjaDesignu, warstwy []WarstwaKompozycji) (KompozycjaDesignu, error) {

	if kompozycja.Kod == "" {
		return KompozycjaDesignu{}, fmt.Errorf("dane: kompozycja design bez identyfikatora")
	}
	if kompozycja.Okno == "" {
		return KompozycjaDesignu{}, fmt.Errorf("dane: kompozycja design %q bez okna", kompozycja.Kod)
	}

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszKompozycjeDesign)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, kompozycja.Kod, kompozycja.Okno,
			tekstDoKolumny(kompozycja.Nazwa)); err != nil {
			return fmt.Errorf("dane: nie można zapisać kompozycji design %q: %w", kompozycja.Kod, err)
		}

		// Kompozycja mogła dopiero powstać w tej transakcji — ID trzeba odczytać
		// przed podmianą warstw, bo `warstwa_kompozycji_design.kompozycja_id`
		// wymaga klucza wewnętrznego, nie identyfikatora zewnętrznego.
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzKompozycjeDesign)
		if err != nil {
			return err
		}
		zapisana, err := odczytajKompozycjeDesign(odczyt.QueryRowContext(ctx, kompozycja.Kod))
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać zapisanej kompozycji design %q: %w",
				kompozycja.Kod, err)
		}

		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunWarstwyKompozycjiDesign)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, zapisana.ID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić warstw kompozycji design %q: %w",
				kompozycja.Kod, err)
		}

		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWarstweKompozycjiDesign)
		if err != nil {
			return err
		}
		for numer, warstwa := range warstwy {
			if warstwa.Kod == "" {
				return fmt.Errorf("dane: warstwa numer %d kompozycji design %q bez identyfikatora",
					numer, kompozycja.Kod)
			}
			kolejnosc := warstwa.Kolejnosc
			if kolejnosc == 0 {
				kolejnosc = numer + 1
			}
			_, err := wstawienie.ExecContext(ctx, warstwa.Kod, zapisana.ID,
				tekstDoKolumny(warstwa.ZasobID), warstwa.X, warstwa.Y,
				warstwa.Szerokosc, warstwa.Wysokosc, kolejnosc,
				liczbaLogiczna(warstwa.Zablokowana), tekstDoKolumny(warstwa.Adnotacja))
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać warstwy %q kompozycji design %q: %w",
					warstwa.Kod, kompozycja.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return KompozycjaDesignu{}, err
	}
	return r.Kompozycja(ctx, kompozycja.Kod)
}

// Kompozycja zwraca kompozycję Design Board o wskazanym kodzie. Brak wiersza
// wraca jako ErrBrakWiersza — warstwa wyższa odróżnia „nie ma” od „odczyt się
// nie powiódł”.
func (r *repozytoriumDesignu) Kompozycja(ctx context.Context, kod string) (KompozycjaDesignu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKompozycjeDesign)
	if err != nil {
		return KompozycjaDesignu{}, err
	}
	kompozycja, err := odczytajKompozycjeDesign(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return KompozycjaDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return KompozycjaDesignu{}, fmt.Errorf("dane: nieczytelny wiersz kompozycji design %q: %w", kod, err)
	}
	return kompozycja, nil
}

// Kompozycje zwraca wszystkie kompozycje wskazanego okna, od ostatnio
// zmienionej. Bez granicy strony: kontrakt `design.board.list` nie niesie ani
// `limit`, ani `total` liczonego osobno od długości wykazu — `Total` odpowiedzi
// jest liczbą zwróconych kompozycji, więc przycięcie wykazu tutaj rozjechałoby
// obie liczby naraz.
//
// Okno puste oddaje wykaz pusty, a nie wszystkie kompozycje: brak wskazania
// okna sprawdza wołający, a zapytanie i tak porównuje kolumnę z pustym
// tekstem, którego żadne okno nie nosi.
func (r *repozytoriumDesignu) Kompozycje(ctx context.Context, okno string) ([]KompozycjaDesignu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKompozycjiDesign)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kompozycji design okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []KompozycjaDesignu{}
	for wiersze.Next() {
		kompozycja, err := odczytajKompozycjeDesign(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kompozycji design okna %q: %w", okno, err)
		}
		lista = append(lista, kompozycja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kompozycji design okna %q: %w", okno, err)
	}
	return lista, nil
}

// odczytajKompozycjeDesign składa strukturę z jednego wiersza wyniku.
func odczytajKompozycjeDesign(wiersz skaner) (KompozycjaDesignu, error) {
	var kompozycja KompozycjaDesignu
	var nazwa sql.NullString
	err := wiersz.Scan(&kompozycja.ID, &kompozycja.Kod, &kompozycja.Okno, &nazwa,
		&kompozycja.Zaktualizowano)
	if err != nil {
		return KompozycjaDesignu{}, err
	}
	kompozycja.Nazwa = tekstZKolumny(nazwa)
	return kompozycja, nil
}

// Warstwy zwraca warstwy kompozycji w zapisanej kolejności renderowania.
func (r *repozytoriumDesignu) Warstwy(ctx context.Context, kompozycjaID int64) ([]WarstwaKompozycji, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaWarstwKompozycjiDesign)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kompozycjaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać warstw kompozycji design %d: %w", kompozycjaID, err)
	}
	defer wiersze.Close()

	lista := []WarstwaKompozycji{}
	for wiersze.Next() {
		var warstwa WarstwaKompozycji
		var zasobID, adnotacja sql.NullString
		var x, y, szerokosc, wysokosc sql.NullFloat64
		var zablokowana int
		err := wiersze.Scan(&warstwa.ID, &warstwa.Kod, &warstwa.KompozycjaID, &zasobID,
			&x, &y, &szerokosc, &wysokosc, &warstwa.Kolejnosc, &zablokowana,
			&adnotacja, &warstwa.Utworzono)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz warstwy kompozycji design %d: %w",
				kompozycjaID, err)
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
		return nil, fmt.Errorf("dane: przerwany odczyt warstw kompozycji design %d: %w", kompozycjaID, err)
	}
	return lista, nil
}
