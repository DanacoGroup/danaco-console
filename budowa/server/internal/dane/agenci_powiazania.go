// Odpowiedzialność pliku: odczyt powiązań eksperta — umiejętności, konektorów
// i uprawnień — oraz dołączenie ich do wierszy biblioteki. Zapis powiązania leży
// w `agenci_powiazania_zapis.go`.
//
// Cztery tabele podrzędne czytamy raz na wywołanie i grupujemy po numerze
// eksperta, zamiast pytać o nie osobno dla każdego wiersza wykazu.
package dane

import (
	"context"
	"fmt"
)

const (
	umiejetnosciAgentow = `SELECT agent_id, kod FROM agent_umiejetnosc ORDER BY agent_id, kod`

	konektoryAgentow = `SELECT agent_id, kod FROM agent_konektor ORDER BY agent_id, nazwa, kod`

	uprawnieniaAgentow = `SELECT agent_id, grupa, zakres, przyznane FROM agent_uprawnienie
	                      ORDER BY agent_id, grupa, zakres`

	// Poziomy pamięci mają ten sam kształt dwóch kolumn co umiejętności
	// i konektory, więc idą tą samą drogą. Brak wierszy dla eksperta znaczy
	// pamięć wyłączoną, a nie „nie odczytano".
	poziomyPamieciAgentow = `SELECT agent_id, poziom FROM agent_pamiec_poziom
	                         ORDER BY agent_id, poziom`

	// Moduły zastosowania mają ten sam kształt dwóch kolumn. Brak wierszy dla
	// eksperta znaczy brak ograniczenia — dostępność we wszystkich modułach.
	modulyZastosowaniaAgentow = `SELECT agent_id, kod_modulu FROM agent_modul_zastosowania
	                             ORDER BY agent_id, kod_modulu`
)

// dolaczPowiazania uzupełnia wiersze ekspertów o umiejętności, konektory,
// uprawnienia i poziomy pamięci. Cztery zapytania na całe wywołanie, nie cztery
// na eksperta.
func (r *repozytoriumAgentow) dolaczPowiazania(ctx context.Context, agenci []Agent) error {
	if len(agenci) == 0 {
		return nil
	}
	umiejetnosci, err := r.napisyPowiazane(ctx, umiejetnosciAgentow, "umiejętności")
	if err != nil {
		return err
	}
	konektory, err := r.napisyPowiazane(ctx, konektoryAgentow, "konektorów")
	if err != nil {
		return err
	}
	uprawnienia, err := r.uprawnienia(ctx)
	if err != nil {
		return err
	}
	poziomy, err := r.napisyPowiazane(ctx, poziomyPamieciAgentow, "poziomów pamięci")
	if err != nil {
		return err
	}
	moduly, err := r.napisyPowiazane(ctx, modulyZastosowaniaAgentow, "modułów zastosowania")
	if err != nil {
		return err
	}
	for i := range agenci {
		agenci[i].ModulyZastosowania = moduly[agenci[i].ID]
		agenci[i].Umiejetnosci = umiejetnosci[agenci[i].ID]
		agenci[i].Konektory = konektory[agenci[i].ID]
		agenci[i].Uprawnienia = uprawnienia[agenci[i].ID]
		agenci[i].PoziomyPamieci = poziomy[agenci[i].ID]
	}
	return nil
}

// napisyPowiazane czyta tabelę podrzędną o kształcie (agent_id, kod).
func (r *repozytoriumAgentow) napisyPowiazane(ctx context.Context,
	zapytanie, obszar string) (map[int64][]string, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać %s ekspertów: %w", obszar, err)
	}
	defer wiersze.Close()

	zebrane := map[int64][]string{}
	for wiersze.Next() {
		var agentID int64
		var kod string
		if err := wiersze.Scan(&agentID, &kod); err != nil {
			return nil, fmt.Errorf("dane: przerwany odczyt %s ekspertów: %w", obszar, err)
		}
		zebrane[agentID] = append(zebrane[agentID], kod)
	}
	return zebrane, wiersze.Err()
}

// uprawnienia czyta cztery grupy zakresu wszystkich ekspertów.
func (r *repozytoriumAgentow) uprawnienia(ctx context.Context) (map[int64][]UprawnienieAgenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, uprawnieniaAgentow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać uprawnień ekspertów: %w", err)
	}
	defer wiersze.Close()

	zebrane := map[int64][]UprawnienieAgenta{}
	for wiersze.Next() {
		var agentID int64
		var wpis UprawnienieAgenta
		var przyznane int
		if err := wiersze.Scan(&agentID, &wpis.Grupa, &wpis.Zakres, &przyznane); err != nil {
			return nil, fmt.Errorf("dane: przerwany odczyt uprawnień ekspertów: %w", err)
		}
		wpis.Przyznane = przyznane == 1
		zebrane[agentID] = append(zebrane[agentID], wpis)
	}
	return zebrane, wiersze.Err()
}
