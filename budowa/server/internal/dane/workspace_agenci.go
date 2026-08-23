// Odpowiedzialność pliku: przypisanie eksperta do projektu — okno Agent Manager
// (tabela `przypisanie_agenta_projektu`).
//
// Projekt ma najwyżej jednego wykonawcę domyślnego. Wskazanie nowego zdejmuje
// wskazanie z poprzedniego w tej samej transakcji, bo inaczej pasek promptu nie
// miałby jednoznacznego adresata zlecenia.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// PrzypisanieAgenta to wiersz przypisania eksperta do projektu.
type PrzypisanieAgenta struct {
	ProjektID         int64
	AgentKod          string
	Rola              *string
	DomyslnyWykonawca bool
	Przypisano        string
}

const (
	kolumnyPrzypisania = `projekt_id, agent_kod, rola, domyslny_wykonawca, przypisano`

	zapiszPrzypisanieAgenta = `INSERT INTO przypisanie_agenta_projektu
	                           (projekt_id, agent_kod, rola, domyslny_wykonawca)
	                           VALUES (?, ?, ?, ?)
	                           ON CONFLICT(projekt_id, agent_kod) DO UPDATE SET
	                               rola = excluded.rola,
	                               domyslny_wykonawca = excluded.domyslny_wykonawca,
	                               przypisano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	zdejmijDomyslnegoWykonawce = `UPDATE przypisanie_agenta_projektu
	                              SET domyslny_wykonawca = 0
	                              WHERE projekt_id = ? AND agent_kod <> ?`

	pobierzPrzypisanieAgenta = `SELECT ` + kolumnyPrzypisania + `
	                            FROM przypisanie_agenta_projektu
	                            WHERE projekt_id = ? AND agent_kod = ?`

	listaPrzypisanAgentow = `SELECT ` + kolumnyPrzypisania + `
	                         FROM przypisanie_agenta_projektu
	                         WHERE projekt_id = ?
	                         ORDER BY domyslny_wykonawca DESC, agent_kod`
)

// PrzypiszAgenta zapisuje przypisanie i zwraca jego stan po zapisie.
func (r *repozytoriumPrzestrzeniRoboczej) PrzypiszAgenta(ctx context.Context,
	przypisanie PrzypisanieAgenta) (PrzypisanieAgenta, error) {

	if przypisanie.AgentKod == "" {
		return PrzypisanieAgenta{}, fmt.Errorf("dane: przypisanie bez wskazania eksperta")
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		_, err := transakcja.ExecContext(ctx, zapiszPrzypisanieAgenta, przypisanie.ProjektID,
			przypisanie.AgentKod, tekstDoKolumny(przypisanie.Rola),
			liczbaLogiczna(przypisanie.DomyslnyWykonawca))
		if err != nil {
			return fmt.Errorf("dane: nie można przypisać eksperta %q: %w", przypisanie.AgentKod, err)
		}
		if !przypisanie.DomyslnyWykonawca {
			return nil
		}
		_, err = transakcja.ExecContext(ctx, zdejmijDomyslnegoWykonawce,
			przypisanie.ProjektID, przypisanie.AgentKod)
		if err != nil {
			return fmt.Errorf("dane: nie można zdjąć poprzedniego wykonawcy domyślnego: %w", err)
		}
		return nil
	})
	if err != nil {
		return PrzypisanieAgenta{}, err
	}
	return r.przypisanieAgenta(ctx, przypisanie.ProjektID, przypisanie.AgentKod)
}

// PrzypisaniaAgentow zwraca ekspertów projektu; domyślny wykonawca jest pierwszy.
func (r *repozytoriumPrzestrzeniRoboczej) PrzypisaniaAgentow(ctx context.Context,
	projektID int64) ([]PrzypisanieAgenta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaPrzypisanAgentow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, projektID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać ekspertów projektu %d: %w", projektID, err)
	}
	defer wiersze.Close()

	lista := []PrzypisanieAgenta{}
	for wiersze.Next() {
		przypisanie, err := odczytajPrzypisanie(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz przypisania eksperta: %w", err)
		}
		lista = append(lista, przypisanie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt ekspertów projektu %d: %w", projektID, err)
	}
	return lista, nil
}

// przypisanieAgenta zwraca jedno przypisanie po zapisie.
func (r *repozytoriumPrzestrzeniRoboczej) przypisanieAgenta(ctx context.Context,
	projektID int64, agent string) (PrzypisanieAgenta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPrzypisanieAgenta)
	if err != nil {
		return PrzypisanieAgenta{}, err
	}
	przypisanie, err := odczytajPrzypisanie(polecenie.QueryRowContext(ctx, projektID, agent))
	if err != nil {
		return PrzypisanieAgenta{}, fmt.Errorf("dane: nieczytelne przypisanie eksperta %q: %w", agent, err)
	}
	return przypisanie, nil
}

// odczytajPrzypisanie składa strukturę z jednego wiersza wyniku.
func odczytajPrzypisanie(wiersz skaner) (PrzypisanieAgenta, error) {
	var przypisanie PrzypisanieAgenta
	var rola sql.NullString
	var domyslny int
	err := wiersz.Scan(&przypisanie.ProjektID, &przypisanie.AgentKod, &rola, &domyslny,
		&przypisanie.Przypisano)
	if err != nil {
		return PrzypisanieAgenta{}, err
	}
	przypisanie.Rola = tekstZKolumny(rola)
	przypisanie.DomyslnyWykonawca = domyslny == 1
	return przypisanie, nil
}
