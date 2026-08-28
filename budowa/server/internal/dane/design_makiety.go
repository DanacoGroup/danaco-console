// Plik definiuje obszar makiet modułu Design: ramki, przynależność warstw do
// ramek, więzy responsywne i siatki układu, jako część kontraktu
// RepozytoriumDesignu zadeklarowanego w design.go.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// RamkaDesignu to wiersz tabeli `ramka_design` — ekran makiety. SiatkaJSON
// i UkladJSON niosą zapis kontraktu bez rozkładania go w warstwie danych;
// składa go i rozkłada adapter, bo to on zna kontrakt.
type RamkaDesignu struct {
	ID                int64
	Kod               string
	KompozycjaID      int64
	Nazwa             string
	Szerokosc         float64
	Wysokosc          float64
	X                 *float64
	Y                 *float64
	NastawaUrzadzenia *string
	SiatkaJSON        *string
	UkladJSON         *string
	Zaktualizowano    string
}

// WiezRamkiDesignu to wiersz tabeli `wiez_ramki_design` — zachowanie jednej
// warstwy przy zmianie rozmiaru ramki. Kotwica pozioma i pionowa stoją razem,
// bo przeliczenie bierze obie naraz.
type WiezRamkiDesignu struct {
	WarstwaKod string
	Poziomo    string
	Pionowo    string
}

const (
	kolumnyRamkiDesignu = `id, identyfikator_zewnetrzny, kompozycja_id, nazwa, szerokosc, wysokosc,
	                       x, y, nastawa_urzadzenia, siatka_json, uklad_json, zaktualizowano`

	zapiszRamkeDesignu = `INSERT INTO ramka_design
	                      (identyfikator_zewnetrzny, kompozycja_id, nazwa, szerokosc, wysokosc,
	                       x, y, nastawa_urzadzenia, siatka_json, uklad_json, zaktualizowano)
	                      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	                              strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                      ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                          nazwa = excluded.nazwa,
	                          szerokosc = excluded.szerokosc,
	                          wysokosc = excluded.wysokosc,
	                          x = excluded.x,
	                          y = excluded.y,
	                          nastawa_urzadzenia = excluded.nastawa_urzadzenia,
	                          siatka_json = excluded.siatka_json,
	                          uklad_json = excluded.uklad_json,
	                          zaktualizowano = excluded.zaktualizowano`

	pobierzRamkeDesignu = `SELECT ` + kolumnyRamkiDesignu +
		` FROM ramka_design WHERE identyfikator_zewnetrzny = ?`

	listaRamekDesignu = `SELECT ` + kolumnyRamkiDesignu +
		` FROM ramka_design WHERE kompozycja_id = ? ORDER BY id`

	usunRamkeDesignu = `DELETE FROM ramka_design WHERE identyfikator_zewnetrzny = ?`

	listaWarstwRamkiDesignu = `SELECT warstwa_kod FROM warstwa_ramki_design
	                           WHERE ramka_id = ? ORDER BY kolejnosc, id`

	wstawWarstweDoRamkiDesignu = `INSERT INTO warstwa_ramki_design (ramka_id, warstwa_kod, kolejnosc)
	                              VALUES (?, ?, ?)
	                              ON CONFLICT(warstwa_kod) DO UPDATE SET
	                                  ramka_id = excluded.ramka_id,
	                                  kolejnosc = excluded.kolejnosc`

	zwolnijWarstwyRamkiDesignu = `DELETE FROM warstwa_ramki_design WHERE ramka_id = ?`

	listaWiezowRamkiDesignu = `SELECT warstwa_kod, poziomo, pionowo FROM wiez_ramki_design
	                           WHERE ramka_id = ? ORDER BY id`

	zapiszWiezRamkiDesignu = `INSERT INTO wiez_ramki_design (ramka_id, warstwa_kod, poziomo, pionowo)
	                          VALUES (?, ?, ?, ?)
	                          ON CONFLICT(ramka_id, warstwa_kod) DO UPDATE SET
	                              poziomo = excluded.poziomo,
	                              pionowo = excluded.pionowo`

	// Zapis więzu oddaje liczbę wierszy naprawdę zmienionych, ustaloną
	// porównaniem ze stanem zastanym.

	pobierzSiatkeKompozycjiDesignu = `SELECT siatka_json FROM siatka_kompozycji_design
	                                  WHERE kompozycja_id = ?`

	zapiszSiatkeKompozycjiDesignu = `INSERT INTO siatka_kompozycji_design
	                                 (kompozycja_id, siatka_json, zaktualizowano)
	                                 VALUES (?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                                 ON CONFLICT(kompozycja_id) DO UPDATE SET
	                                     siatka_json = excluded.siatka_json,
	                                     zaktualizowano = excluded.zaktualizowano`

	pobierzWarstweKompozycjiDesignuPoKodzie = `SELECT id, identyfikator_zewnetrzny, kompozycja_id,
	                                                  zasob_id, x, y, szerokosc, wysokosc, kolejnosc,
	                                                  zablokowana, adnotacja, utworzono
	                                           FROM warstwa_kompozycji_design
	                                           WHERE identyfikator_zewnetrzny = ?`

	przestawWarstweKompozycjiDesignu = `UPDATE warstwa_kompozycji_design
	                                    SET x = ?, y = ?, szerokosc = ?, wysokosc = ?
	                                    WHERE identyfikator_zewnetrzny = ?`

	dolozWarstweKompozycjiDesignu = `INSERT INTO warstwa_kompozycji_design
	                                 (identyfikator_zewnetrzny, kompozycja_id, zasob_id, x, y,
	                                  szerokosc, wysokosc, kolejnosc, zablokowana, adnotacja)
	                                 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	najwyzszaKolejnoscWarstwyDesignu = `SELECT COALESCE(MAX(kolejnosc), 0)
	                                    FROM warstwa_kompozycji_design WHERE kompozycja_id = ?`
)

// ZapiszRamkeDesignu zakłada ramkę albo nadpisuje zastaną po identyfikatorze
// zewnętrznym, zwracając stan ramki po zapisie.
func (r *repozytoriumDesignu) ZapiszRamkeDesignu(ctx context.Context,
	ramka RamkaDesignu) (RamkaDesignu, error) {

	if ramka.Kod == "" {
		return RamkaDesignu{}, fmt.Errorf("dane: ramka design bez identyfikatora")
	}
	if ramka.KompozycjaID == 0 {
		return RamkaDesignu{}, fmt.Errorf("dane: ramka design %q bez kompozycji", ramka.Kod)
	}
	if ramka.Nazwa == "" {
		return RamkaDesignu{}, fmt.Errorf("dane: ramka design %q bez nazwy", ramka.Kod)
	}

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszRamkeDesignu)
	if err != nil {
		return RamkaDesignu{}, err
	}
	_, err = polecenie.ExecContext(ctx, ramka.Kod, ramka.KompozycjaID, ramka.Nazwa,
		ramka.Szerokosc, ramka.Wysokosc, ulamekDoKolumnyDesignu(ramka.X),
		ulamekDoKolumnyDesignu(ramka.Y), tekstDoKolumny(ramka.NastawaUrzadzenia),
		tekstDoKolumny(ramka.SiatkaJSON), tekstDoKolumny(ramka.UkladJSON))
	if err != nil {
		return RamkaDesignu{}, fmt.Errorf("dane: nie można zapisać ramki design %q: %w", ramka.Kod, err)
	}
	return r.RamkaDesignuPoKodzie(ctx, ramka.Kod)
}

// RamkaDesignuPoKodzie zwraca ramkę o wskazanym identyfikatorze zewnętrznym,
// zwracając błąd ErrBrakWiersza, gdy ramka nie istnieje.
func (r *repozytoriumDesignu) RamkaDesignuPoKodzie(ctx context.Context,
	kod string) (RamkaDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzRamkeDesignu)
	if err != nil {
		return RamkaDesignu{}, err
	}
	ramka, err := odczytajRamkeDesignu(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return RamkaDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return RamkaDesignu{}, fmt.Errorf("dane: nieczytelny wiersz ramki design %q: %w", kod, err)
	}
	return ramka, nil
}

// RamkiDesignu zwraca ramki wskazanej kompozycji w kolejności założenia, od
// pierwszej dodanej do ostatniej.
func (r *repozytoriumDesignu) RamkiDesignu(ctx context.Context,
	kompozycjaID int64) ([]RamkaDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaRamekDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kompozycjaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać ramek design kompozycji %d: %w",
			kompozycjaID, err)
	}
	defer wiersze.Close()

	lista := []RamkaDesignu{}
	for wiersze.Next() {
		ramka, err := odczytajRamkeDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz ramki design kompozycji %d: %w",
				kompozycjaID, err)
		}
		lista = append(lista, ramka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt ramek design kompozycji %d: %w",
			kompozycjaID, err)
	}
	return lista, nil
}

// UsunRamkeDesignu usuwa ramkę i mówi, czy wiersz istniał. Warstwy ramki
// zostają na płótnie — kaskada schematu sięga wyłącznie przynależności i więzów
// (migracja 317), bo usunięcie ramki nie jest usunięciem pracy, która w niej
// leżała.
func (r *repozytoriumDesignu) UsunRamkeDesignu(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunRamkeDesignu)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć ramki design %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można ustalić skutku usunięcia ramki design %q: %w", kod, err)
	}
	return zmienione > 0, nil
}

// WarstwyRamkiDesignu zwraca kody warstw należących do ramki, uporządkowane
// według kolejności ich przypisania.
func (r *repozytoriumDesignu) WarstwyRamkiDesignu(ctx context.Context, ramkaID int64) ([]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaWarstwRamkiDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, ramkaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać warstw ramki design %d: %w", ramkaID, err)
	}
	defer wiersze.Close()

	lista := []string{}
	for wiersze.Next() {
		var kod string
		if err := wiersze.Scan(&kod); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz warstwy ramki design %d: %w", ramkaID, err)
		}
		lista = append(lista, kod)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt warstw ramki design %d: %w", ramkaID, err)
	}
	return lista, nil
}

// PrzypiszWarstwyDoRamkiDesignu wiąże warstwy z ramką w podanej kolejności.
// Warstwa należąca dotąd do innej ramki przechodzi do wskazanej — warunek
// UNIQUE schematu na kodzie warstwy pilnuje, że należy najwyżej do jednej.
func (r *repozytoriumDesignu) PrzypiszWarstwyDoRamkiDesignu(ctx context.Context,
	ramkaID int64, warstwy []string) error {

	if len(warstwy) == 0 {
		return nil
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWarstweDoRamkiDesignu)
		if err != nil {
			return err
		}
		for numer, kod := range warstwy {
			if kod == "" {
				continue
			}
			if _, err := polecenie.ExecContext(ctx, ramkaID, kod, numer+1); err != nil {
				return fmt.Errorf("dane: nie można przypisać warstwy %q do ramki design %d: %w",
					kod, ramkaID, err)
			}
		}
		return nil
	})
}

// ZwolnijWarstwyRamkiDesignu zdejmuje przynależność warstw do ramki i oddaje
// kody zwolnionych. Warstwy zostają na płótnie.
func (r *repozytoriumDesignu) ZwolnijWarstwyRamkiDesignu(ctx context.Context,
	ramkaID int64) ([]string, error) {

	zwolnione, err := r.WarstwyRamkiDesignu(ctx, ramkaID)
	if err != nil {
		return nil, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zwolnijWarstwyRamkiDesignu)
	if err != nil {
		return nil, err
	}
	if _, err := polecenie.ExecContext(ctx, ramkaID); err != nil {
		return nil, fmt.Errorf("dane: nie można zwolnić warstw ramki design %d: %w", ramkaID, err)
	}
	return zwolnione, nil
}

// WiezyRamkiDesignu zwraca więzy responsywne ramki, opisujące zachowanie
// każdej warstwy przy zmianie rozmiaru.
func (r *repozytoriumDesignu) WiezyRamkiDesignu(ctx context.Context,
	ramkaID int64) ([]WiezRamkiDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaWiezowRamkiDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, ramkaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać więzów ramki design %d: %w", ramkaID, err)
	}
	defer wiersze.Close()

	lista := []WiezRamkiDesignu{}
	for wiersze.Next() {
		var wiez WiezRamkiDesignu
		if err := wiersze.Scan(&wiez.WarstwaKod, &wiez.Poziomo, &wiez.Pionowo); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz więzu ramki design %d: %w", ramkaID, err)
		}
		lista = append(lista, wiez)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt więzów ramki design %d: %w", ramkaID, err)
	}
	return lista, nil
}

// ZapiszWiezyRamkiDesignu utrwala więzy ramki w jednej transakcji, nadpisując
// więzy zastane dla tych samych warstw.
func (r *repozytoriumDesignu) ZapiszWiezyRamkiDesignu(ctx context.Context, ramkaID int64,
	wiezy []WiezRamkiDesignu) error {

	if len(wiezy) == 0 {
		return nil
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszWiezRamkiDesignu)
		if err != nil {
			return err
		}
		for _, wiez := range wiezy {
			if wiez.WarstwaKod == "" {
				return fmt.Errorf("dane: więz ramki design %d bez warstwy", ramkaID)
			}
			if _, err := polecenie.ExecContext(ctx, ramkaID, wiez.WarstwaKod,
				wiez.Poziomo, wiez.Pionowo); err != nil {
				return fmt.Errorf("dane: nie można zapisać więzu warstwy %q ramki design %d: %w",
					wiez.WarstwaKod, ramkaID, err)
			}
		}
		return nil
	})
}

// SiatkaKompozycjiDesignu zwraca zapis siatki obowiązującej całą kompozycję.
// Brak wiersza wraca jako ErrBrakWiersza — kompozycja bez siatki to nie to samo,
// co kompozycja z siatką pustą.
func (r *repozytoriumDesignu) SiatkaKompozycjiDesignu(ctx context.Context,
	kompozycjaID int64) (string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSiatkeKompozycjiDesignu)
	if err != nil {
		return "", err
	}
	var zapis string
	err = polecenie.QueryRowContext(ctx, kompozycjaID).Scan(&zapis)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrBrakWiersza
	}
	if err != nil {
		return "", fmt.Errorf("dane: nieczytelna siatka kompozycji design %d: %w", kompozycjaID, err)
	}
	return zapis, nil
}

// ZapiszSiatkeKompozycjiDesignu utrwala siatkę obowiązującą całą kompozycję,
// nadpisując zapis zastany dla tej kompozycji.
func (r *repozytoriumDesignu) ZapiszSiatkeKompozycjiDesignu(ctx context.Context,
	kompozycjaID int64, siatkaJSON string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSiatkeKompozycjiDesignu)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, kompozycjaID, siatkaJSON); err != nil {
		return fmt.Errorf("dane: nie można zapisać siatki kompozycji design %d: %w", kompozycjaID, err)
	}
	return nil
}

// WarstwaKompozycjiDesignuPoKodzie zwraca warstwę po identyfikatorze
// zewnętrznym. Układ automatyczny i więzy pracują na warstwach wskazanych
// kodem, a bez tego odczytu nie miałyby czego przeliczyć.
func (r *repozytoriumDesignu) WarstwaKompozycjiDesignuPoKodzie(ctx context.Context,
	kod string) (WarstwaKompozycji, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWarstweKompozycjiDesignuPoKodzie)
	if err != nil {
		return WarstwaKompozycji{}, err
	}
	warstwa, err := odczytajWarstweKompozycjiDesignu(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return WarstwaKompozycji{}, ErrBrakWiersza
	}
	if err != nil {
		return WarstwaKompozycji{}, fmt.Errorf(
			"dane: nieczytelny wiersz warstwy kompozycji design %q: %w", kod, err)
	}
	return warstwa, nil
}

// PrzestawWarstweKompozycjiDesignu zapisuje nowe położenie i rozmiar jednej
// warstwy. Osobno od zapisu pełnego planszy — powód w nagłówku pliku.
func (r *repozytoriumDesignu) PrzestawWarstweKompozycjiDesignu(ctx context.Context,
	warstwa WarstwaKompozycji) error {

	if warstwa.Kod == "" {
		return fmt.Errorf("dane: przestawienie warstwy kompozycji design bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, przestawWarstweKompozycjiDesignu)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, ulamekDoKolumnyDesignu(warstwa.X),
		ulamekDoKolumnyDesignu(warstwa.Y), ulamekDoKolumnyDesignu(warstwa.Szerokosc),
		ulamekDoKolumnyDesignu(warstwa.Wysokosc), warstwa.Kod)
	if err != nil {
		return fmt.Errorf("dane: nie można przestawić warstwy kompozycji design %q: %w",
			warstwa.Kod, err)
	}
	return nil
}

// DolozWarstweKompozycjiDesignu dokłada jedną warstwę do planszy, nie ruszając
// pozostałych. Kolejność zerowa znaczy „na wierzch".
func (r *repozytoriumDesignu) DolozWarstweKompozycjiDesignu(ctx context.Context,
	warstwa WarstwaKompozycji) (WarstwaKompozycji, error) {

	if warstwa.Kod == "" {
		return WarstwaKompozycji{}, fmt.Errorf("dane: warstwa kompozycji design bez identyfikatora")
	}
	if warstwa.KompozycjaID == 0 {
		return WarstwaKompozycji{}, fmt.Errorf("dane: warstwa kompozycji design %q bez kompozycji",
			warstwa.Kod)
	}

	kolejnosc := warstwa.Kolejnosc
	if kolejnosc == 0 {
		najwyzsza, err := r.zapytania.przygotuj(ctx, najwyzszaKolejnoscWarstwyDesignu)
		if err != nil {
			return WarstwaKompozycji{}, err
		}
		var szczyt int
		if err := najwyzsza.QueryRowContext(ctx, warstwa.KompozycjaID).Scan(&szczyt); err != nil {
			return WarstwaKompozycji{}, fmt.Errorf(
				"dane: nie można ustalić kolejności warstwy kompozycji design %q: %w", warstwa.Kod, err)
		}
		kolejnosc = szczyt + 1
	}

	polecenie, err := r.zapytania.przygotuj(ctx, dolozWarstweKompozycjiDesignu)
	if err != nil {
		return WarstwaKompozycji{}, err
	}
	_, err = polecenie.ExecContext(ctx, warstwa.Kod, warstwa.KompozycjaID,
		tekstDoKolumny(warstwa.ZasobID), ulamekDoKolumnyDesignu(warstwa.X),
		ulamekDoKolumnyDesignu(warstwa.Y), ulamekDoKolumnyDesignu(warstwa.Szerokosc),
		ulamekDoKolumnyDesignu(warstwa.Wysokosc), kolejnosc, warstwa.Zablokowana,
		tekstDoKolumny(warstwa.Adnotacja))
	if err != nil {
		return WarstwaKompozycji{}, fmt.Errorf(
			"dane: nie można dołożyć warstwy kompozycji design %q: %w", warstwa.Kod, err)
	}
	return r.WarstwaKompozycjiDesignuPoKodzie(ctx, warstwa.Kod)
}

// odczytajRamkeDesignu składa strukturę ramki z jednego wiersza wyniku
// zapytania, zamieniając kolumny nullowalne na wskaźniki.
func odczytajRamkeDesignu(wiersz skaner) (RamkaDesignu, error) {
	var ramka RamkaDesignu
	var x, y sql.NullFloat64
	var nastawa, siatka, uklad sql.NullString
	err := wiersz.Scan(&ramka.ID, &ramka.Kod, &ramka.KompozycjaID, &ramka.Nazwa,
		&ramka.Szerokosc, &ramka.Wysokosc, &x, &y, &nastawa, &siatka, &uklad,
		&ramka.Zaktualizowano)
	if err != nil {
		return RamkaDesignu{}, err
	}
	ramka.X = ulamekZKolumnyDesignu(x)
	ramka.Y = ulamekZKolumnyDesignu(y)
	ramka.NastawaUrzadzenia = tekstZKolumny(nastawa)
	ramka.SiatkaJSON = tekstZKolumny(siatka)
	ramka.UkladJSON = tekstZKolumny(uklad)
	return ramka, nil
}

// odczytajWarstweKompozycjiDesignu składa warstwę z jednego wiersza wyniku
// zapytania, zamieniając kolumny nullowalne na wskaźniki.
func odczytajWarstweKompozycjiDesignu(wiersz skaner) (WarstwaKompozycji, error) {
	var warstwa WarstwaKompozycji
	var zasob, adnotacja sql.NullString
	var x, y, szerokosc, wysokosc sql.NullFloat64
	err := wiersz.Scan(&warstwa.ID, &warstwa.Kod, &warstwa.KompozycjaID, &zasob,
		&x, &y, &szerokosc, &wysokosc, &warstwa.Kolejnosc, &warstwa.Zablokowana,
		&adnotacja, &warstwa.Utworzono)
	if err != nil {
		return WarstwaKompozycji{}, err
	}
	warstwa.ZasobID = tekstZKolumny(zasob)
	warstwa.X = ulamekZKolumnyDesignu(x)
	warstwa.Y = ulamekZKolumnyDesignu(y)
	warstwa.Szerokosc = ulamekZKolumnyDesignu(szerokosc)
	warstwa.Wysokosc = ulamekZKolumnyDesignu(wysokosc)
	warstwa.Adnotacja = tekstZKolumny(adnotacja)
	return warstwa, nil
}

// ulamekDoKolumnyDesignu przekłada wskaźnik na wartość kolumny dopuszczającej
// NULL: nil znaczy „nieustawione", nie „zero".
func ulamekDoKolumnyDesignu(wartosc *float64) any {
	if wartosc == nil {
		return nil
	}
	return *wartosc
}

// ulamekZKolumnyDesignu przekłada kolumnę dopuszczającą NULL na wskaźnik
// ułamka, zwracając nil dla wartości nieustawionej.
func ulamekZKolumnyDesignu(kolumna sql.NullFloat64) *float64 {
	if !kolumna.Valid {
		return nil
	}
	wartosc := kolumna.Float64
	return &wartosc
}
