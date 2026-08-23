// Odpowiedzialność pliku: pamięć projektu okna Context Memory (tabela
// `wpis_pamieci_projektu`).
//
// Wpis niesie poziom zasięgu współdzielenia. Zasięg węższy albo równy projektowi
// widzi wyłącznie ten projekt; zasięg szerszy (globalny, środowisko, moduł, para
// modułów) jest ustaleniem wspólnym i wchodzi do pamięci innych projektów na
// żądanie `includeShared`. Kolejność rozstrzygania poziomów należy do pakietu
// `internal/konfig` — tu leży wyłącznie zapis i odczyt.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// WpisPamieciProjektu to wiersz pamięci projektu.
type WpisPamieciProjektu struct {
	ID             int64
	ProjektID      int64
	ProjektKod     string
	Identyfikator  string
	Tresc          string
	Przypiety      bool
	Pochodzenie    shared.MemoryEntryOrigin
	Poziom         shared.ConfigScope
	KluczZasiegu   string
	Utworzono      string
	Zaktualizowano string
}

const (
	kolumnyWpisuPamieci = `w.id, p.kod, w.projekt_id, w.identyfikator_zewnetrzny, w.tresc,
	                       w.przypiety, w.pochodzenie, z.kod, w.klucz_zasiegu,
	                       w.utworzono, w.zaktualizowano`

	zrodloWpisuPamieci = ` FROM wpis_pamieci_projektu w
	                       JOIN projekt p ON p.id = w.projekt_id
	                       JOIN poziom_zasiegu z ON z.id = w.poziom_zasiegu_id`

	zapiszWpisPamieci = `INSERT INTO wpis_pamieci_projektu
	                     (projekt_id, identyfikator_zewnetrzny, tresc, przypiety, pochodzenie,
	                      poziom_zasiegu_id, klucz_zasiegu)
	                     VALUES (?, ?, ?, ?, ?, (SELECT id FROM poziom_zasiegu WHERE kod = ?), ?)
	                     ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                         tresc = excluded.tresc,
	                         przypiety = excluded.przypiety,
	                         pochodzenie = excluded.pochodzenie,
	                         poziom_zasiegu_id = excluded.poziom_zasiegu_id,
	                         klucz_zasiegu = excluded.klucz_zasiegu,
	                         zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzWpisPamieci = `SELECT ` + kolumnyWpisuPamieci + zrodloWpisuPamieci +
		` WHERE w.identyfikator_zewnetrzny = ?`

	// Wpisy przypięte idą na początek wykazu; w obrębie grupy decyduje czas
	// ostatniej zmiany.
	listaWpisowPamieci = `SELECT ` + kolumnyWpisuPamieci + zrodloWpisuPamieci +
		` WHERE w.projekt_id = ?
		  ORDER BY w.przypiety DESC, w.zaktualizowano DESC
		  LIMIT ?`

	listaWpisowWspoldzielonych = `SELECT ` + kolumnyWpisuPamieci + zrodloWpisuPamieci +
		` WHERE w.projekt_id <> ? AND z.pierwszenstwo < (
		      SELECT pierwszenstwo FROM poziom_zasiegu WHERE kod = 'projekt')
		  ORDER BY w.przypiety DESC, w.zaktualizowano DESC
		  LIMIT ?`

	usunWpisPamieci = `DELETE FROM wpis_pamieci_projektu WHERE identyfikator_zewnetrzny = ?`

	przestawZasiegWpisuPamieci = `UPDATE wpis_pamieci_projektu
	                              SET poziom_zasiegu_id = (SELECT id FROM poziom_zasiegu WHERE kod = ?),
	                                  klucz_zasiegu = ?,
	                                  zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                              WHERE identyfikator_zewnetrzny = ?`
)

// ZapiszWpisPamieci zakłada wpis albo nadpisuje wpis o tym samym identyfikatorze
// i zwraca stan po zapisie.
func (r *repozytoriumPrzestrzeniRoboczej) ZapiszWpisPamieci(ctx context.Context,
	wpis WpisPamieciProjektu) (WpisPamieciProjektu, error) {

	if wpis.Identyfikator == "" {
		return WpisPamieciProjektu{}, fmt.Errorf("dane: wpis pamięci bez identyfikatora")
	}
	poziom, err := poziomZasieguNaBaze(wpis.Poziom)
	if err != nil {
		return WpisPamieciProjektu{}, err
	}
	pochodzenie := string(wpis.Pochodzenie)
	if pochodzenie == "" {
		pochodzenie = shared.MemoryEntryOriginOperator
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWpisPamieci)
	if err != nil {
		return WpisPamieciProjektu{}, err
	}
	_, err = polecenie.ExecContext(ctx, wpis.ProjektID, wpis.Identyfikator, wpis.Tresc,
		liczbaLogiczna(wpis.Przypiety), pochodzenie, poziom, wpis.KluczZasiegu)
	if err != nil {
		return WpisPamieciProjektu{}, fmt.Errorf("dane: nie można zapisać wpisu pamięci %q: %w",
			wpis.Identyfikator, err)
	}
	return r.wpisPamieci(ctx, wpis.Identyfikator)
}

// WpisyPamieci zwraca wpisy jednego projektu, przypięte na początku.
func (r *repozytoriumPrzestrzeniRoboczej) WpisyPamieci(ctx context.Context,
	projektID int64, limit int) ([]WpisPamieciProjektu, error) {

	return r.wpisy(ctx, listaWpisowPamieci, "pamięci projektu", projektID, granicaWykazu(limit))
}

// WpisyPamieciWspoldzielone zwraca ustalenia innych projektów zapisane na
// poziomie szerszym niż projekt — te obowiązują wspólnie.
func (r *repozytoriumPrzestrzeniRoboczej) WpisyPamieciWspoldzielone(ctx context.Context,
	projektID int64, limit int) ([]WpisPamieciProjektu, error) {

	return r.wpisy(ctx, listaWpisowWspoldzielonych, "pamięci współdzielonej",
		projektID, granicaWykazu(limit))
}

// WpisPamieciPoIdentyfikatorze zwraca jeden wpis pamięci. Drugi wynik mówi, czy
// wpis istnieje: brak wiersza nie jest awarią odczytu, lecz stanem, który warstwa
// wyższa zamienia na odmowę z kodem `not_found`. Osobny wynik logiczny jest
// potrzebny, bo `wpisPamieci` zawija `sql.ErrNoRows` we własny błąd.
func (r *repozytoriumPrzestrzeniRoboczej) WpisPamieciPoIdentyfikatorze(ctx context.Context,
	identyfikator string) (WpisPamieciProjektu, bool, error) {

	if identyfikator == "" {
		return WpisPamieciProjektu{}, false, fmt.Errorf("dane: wpis pamięci bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWpisPamieci)
	if err != nil {
		return WpisPamieciProjektu{}, false, err
	}
	wpis, err := odczytajWpisPamieci(polecenie.QueryRowContext(ctx, identyfikator))
	if errors.Is(err, sql.ErrNoRows) {
		return WpisPamieciProjektu{}, false, nil
	}
	if err != nil {
		return WpisPamieciProjektu{}, false,
			fmt.Errorf("dane: nieczytelny wpis pamięci %q: %w", identyfikator, err)
	}
	return wpis, true, nil
}

// UsunWpisPamieci kasuje wpis pamięci. Drugi wynik mówi, czy wiersz naprawdę
// zniknął — usunięcie wpisu, którego nie było, nie jest sukcesem.
func (r *repozytoriumPrzestrzeniRoboczej) UsunWpisPamieci(ctx context.Context,
	identyfikator string) (bool, error) {

	if identyfikator == "" {
		return false, fmt.Errorf("dane: wpis pamięci bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, usunWpisPamieci)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, identyfikator)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć wpisu pamięci %q: %w", identyfikator, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek usunięcia wpisu pamięci %q: %w",
			identyfikator, err)
	}
	return zmienione > 0, nil
}

// PrzestawZasiegWpisuPamieci zmienia poziom zasięgu wpisu i byt tego poziomu,
// nie ruszając treści. Służy sprowadzeniu ustalenia wspólnego z powrotem do
// projektu, który je niesie (komenda `memory.detach`).
func (r *repozytoriumPrzestrzeniRoboczej) PrzestawZasiegWpisuPamieci(ctx context.Context,
	identyfikator string, poziom shared.ConfigScope,
	kluczZasiegu string) (WpisPamieciProjektu, error) {

	kod, err := poziomZasieguNaBaze(poziom)
	if err != nil {
		return WpisPamieciProjektu{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, przestawZasiegWpisuPamieci)
	if err != nil {
		return WpisPamieciProjektu{}, err
	}
	if _, err := polecenie.ExecContext(ctx, kod, kluczZasiegu, identyfikator); err != nil {
		return WpisPamieciProjektu{}, fmt.Errorf(
			"dane: nie można przestawić zasięgu wpisu pamięci %q: %w", identyfikator, err)
	}
	return r.wpisPamieci(ctx, identyfikator)
}

// wpisPamieci zwraca jeden wpis po identyfikatorze kontraktu.
func (r *repozytoriumPrzestrzeniRoboczej) wpisPamieci(ctx context.Context,
	identyfikator string) (WpisPamieciProjektu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWpisPamieci)
	if err != nil {
		return WpisPamieciProjektu{}, err
	}
	wpis, err := odczytajWpisPamieci(polecenie.QueryRowContext(ctx, identyfikator))
	if err != nil {
		return WpisPamieciProjektu{}, fmt.Errorf("dane: nieczytelny wpis pamięci %q: %w",
			identyfikator, err)
	}
	return wpis, nil
}

// wpisy wykonuje zapytanie zwracające wiele wierszy pamięci.
func (r *repozytoriumPrzestrzeniRoboczej) wpisy(ctx context.Context, zapytanie, opis string,
	argumenty ...any) ([]WpisPamieciProjektu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać %s: %w", opis, err)
	}
	defer wiersze.Close()

	lista := []WpisPamieciProjektu{}
	for wiersze.Next() {
		wpis, err := odczytajWpisPamieci(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz %s: %w", opis, err)
		}
		lista = append(lista, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt %s: %w", opis, err)
	}
	return lista, nil
}

// granicaWykazu przekłada granicę żądania na argument zapytania. Brak granicy
// znaczy wykaz pełny, nie wykaz pusty.
func granicaWykazu(limit int) int {
	if limit <= 0 {
		return -1
	}
	return limit
}

// odczytajWpisPamieci składa strukturę z jednego wiersza wyniku.
func odczytajWpisPamieci(wiersz skaner) (WpisPamieciProjektu, error) {
	var wpis WpisPamieciProjektu
	var przypiety int
	var pochodzenie, poziom string
	err := wiersz.Scan(&wpis.ID, &wpis.ProjektKod, &wpis.ProjektID, &wpis.Identyfikator, &wpis.Tresc,
		&przypiety, &pochodzenie, &poziom, &wpis.KluczZasiegu, &wpis.Utworzono, &wpis.Zaktualizowano)
	if err != nil {
		return WpisPamieciProjektu{}, err
	}
	wpis.Przypiety = przypiety == 1
	wpis.Pochodzenie = shared.MemoryEntryOrigin(pochodzenie)
	if wpis.Poziom, err = poziomZasieguZBazy(poziom); err != nil {
		return WpisPamieciProjektu{}, err
	}
	return wpis, nil
}
