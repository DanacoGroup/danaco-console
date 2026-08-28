// Repozytorium obsługuje tabelę `polaczenie_prototypu_design` z migracji 319,
// przechowującą połączenia między ramkami prototypu modułu Design, adresowane
// identyfikatorem zewnętrznym.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// PolaczeniePrototypuDesignu to wiersz tabeli `polaczenie_prototypu_design` —
// jedno przejście między ekranami makiety.
type PolaczeniePrototypuDesignu struct {
	ID             int64
	Kod            string
	KompozycjaID   int64
	RamkaOdKod     string
	RamkaDoKod     string
	Wyzwalacz      string
	Przejscie      string
	CzasMs         *int64
	WarstwaKod     *string
	Zaktualizowano string
}

const (
	kolumnyPolaczeniaPrototypuDesignu = `id, identyfikator_zewnetrzny, kompozycja_id, ramka_od_kod,
	                                     ramka_do_kod, wyzwalacz, przejscie, czas_ms, warstwa_kod,
	                                     zaktualizowano`

	zapiszPolaczeniePrototypuDesignu = `INSERT INTO polaczenie_prototypu_design
	                                    (identyfikator_zewnetrzny, kompozycja_id, ramka_od_kod,
	                                     ramka_do_kod, wyzwalacz, przejscie, czas_ms, warstwa_kod,
	                                     zaktualizowano)
	                                    VALUES (?, ?, ?, ?, ?, ?, ?, ?,
	                                            strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                                    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                        ramka_od_kod = excluded.ramka_od_kod,
	                                        ramka_do_kod = excluded.ramka_do_kod,
	                                        wyzwalacz = excluded.wyzwalacz,
	                                        przejscie = excluded.przejscie,
	                                        czas_ms = excluded.czas_ms,
	                                        warstwa_kod = excluded.warstwa_kod,
	                                        zaktualizowano = excluded.zaktualizowano`

	pobierzPolaczeniePrototypuDesignu = `SELECT ` + kolumnyPolaczeniaPrototypuDesignu +
		` FROM polaczenie_prototypu_design WHERE identyfikator_zewnetrzny = ?`

	listaPolaczenPrototypuDesignu = `SELECT ` + kolumnyPolaczeniaPrototypuDesignu +
		` FROM polaczenie_prototypu_design WHERE kompozycja_id = ? ORDER BY id`

	usunPolaczeniePrototypuDesignu = `DELETE FROM polaczenie_prototypu_design
	                                  WHERE identyfikator_zewnetrzny = ?`
)

// ZapiszPolaczeniePrototypuDesignu zakłada połączenie albo nadpisuje zastane po
// identyfikatorze zewnętrznym.
func (r *repozytoriumDesignu) ZapiszPolaczeniePrototypuDesignu(ctx context.Context,
	polaczenie PolaczeniePrototypuDesignu) (PolaczeniePrototypuDesignu, error) {

	if polaczenie.Kod == "" {
		return PolaczeniePrototypuDesignu{}, fmt.Errorf(
			"dane: połączenie prototypu design bez identyfikatora")
	}
	if polaczenie.KompozycjaID == 0 {
		return PolaczeniePrototypuDesignu{}, fmt.Errorf(
			"dane: połączenie prototypu design %q bez kompozycji", polaczenie.Kod)
	}
	if polaczenie.RamkaOdKod == "" || polaczenie.RamkaDoKod == "" {
		return PolaczeniePrototypuDesignu{}, fmt.Errorf(
			"dane: połączenie prototypu design %q bez obu ramek", polaczenie.Kod)
	}

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPolaczeniePrototypuDesignu)
	if err != nil {
		return PolaczeniePrototypuDesignu{}, err
	}
	_, err = polecenie.ExecContext(ctx, polaczenie.Kod, polaczenie.KompozycjaID,
		polaczenie.RamkaOdKod, polaczenie.RamkaDoKod, polaczenie.Wyzwalacz,
		polaczenie.Przejscie, liczbaDoKolumny(polaczenie.CzasMs),
		tekstDoKolumny(polaczenie.WarstwaKod))
	if err != nil {
		return PolaczeniePrototypuDesignu{}, fmt.Errorf(
			"dane: nie można zapisać połączenia prototypu design %q: %w", polaczenie.Kod, err)
	}
	return r.PolaczeniePrototypuDesignuPoKodzie(ctx, polaczenie.Kod)
}

// PolaczeniePrototypuDesignuPoKodzie zwraca połączenie z tabeli
// `polaczenie_prototypu_design` wskazane identyfikatorem zewnętrznym.
func (r *repozytoriumDesignu) PolaczeniePrototypuDesignuPoKodzie(ctx context.Context,
	kod string) (PolaczeniePrototypuDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPolaczeniePrototypuDesignu)
	if err != nil {
		return PolaczeniePrototypuDesignu{}, err
	}
	polaczenie, err := odczytajPolaczeniePrototypuDesignu(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return PolaczeniePrototypuDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return PolaczeniePrototypuDesignu{}, fmt.Errorf(
			"dane: nieczytelny wiersz połączenia prototypu design %q: %w", kod, err)
	}
	return polaczenie, nil
}

// PolaczeniaPrototypuDesignu zwraca komplet połączeń przypisanych do wskazanej
// kompozycji, uporządkowany rosnąco według identyfikatora.
func (r *repozytoriumDesignu) PolaczeniaPrototypuDesignu(ctx context.Context,
	kompozycjaID int64) ([]PolaczeniePrototypuDesignu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaPolaczenPrototypuDesignu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kompozycjaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać połączeń prototypu design kompozycji %d: %w",
			kompozycjaID, err)
	}
	defer wiersze.Close()

	lista := []PolaczeniePrototypuDesignu{}
	for wiersze.Next() {
		polaczenie, err := odczytajPolaczeniePrototypuDesignu(wiersze)
		if err != nil {
			return nil, fmt.Errorf(
				"dane: nieczytelny wiersz połączenia prototypu design kompozycji %d: %w",
				kompozycjaID, err)
		}
		lista = append(lista, polaczenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt połączeń prototypu design kompozycji %d: %w",
			kompozycjaID, err)
	}
	return lista, nil
}

// UsunPolaczeniePrototypuDesignu usuwa połączenie z tabeli
// `polaczenie_prototypu_design` i zwraca informację, czy wiersz istniał przed usunięciem.
func (r *repozytoriumDesignu) UsunPolaczeniePrototypuDesignu(ctx context.Context,
	kod string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, usunPolaczeniePrototypuDesignu)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć połączenia prototypu design %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf(
			"dane: nie można ustalić skutku usunięcia połączenia prototypu design %q: %w", kod, err)
	}
	return zmienione > 0, nil
}

// odczytajPolaczeniePrototypuDesignu składa strukturę PolaczeniePrototypuDesignu
// z jednego wiersza wyniku zapytania SQL.
func odczytajPolaczeniePrototypuDesignu(wiersz skaner) (PolaczeniePrototypuDesignu, error) {
	var polaczenie PolaczeniePrototypuDesignu
	var czas sql.NullInt64
	var warstwa sql.NullString
	err := wiersz.Scan(&polaczenie.ID, &polaczenie.Kod, &polaczenie.KompozycjaID,
		&polaczenie.RamkaOdKod, &polaczenie.RamkaDoKod, &polaczenie.Wyzwalacz,
		&polaczenie.Przejscie, &czas, &warstwa, &polaczenie.Zaktualizowano)
	if err != nil {
		return PolaczeniePrototypuDesignu{}, err
	}
	polaczenie.CzasMs = liczbaZKolumny(czas)
	polaczenie.WarstwaKod = tekstZKolumny(warstwa)
	return polaczenie, nil
}
