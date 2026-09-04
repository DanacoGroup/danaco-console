// Plik zapisuje powiązania eksperta: przypisanie umiejętności, podłączenie konektora oraz ustalenie uprawnienia; każdy
// zapis oddaje ekspertowi zwrotnie zmieniony wiersz, więc warstwa wyższa nie wykonuje drugiego odczytu.
package dane

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	wstawUmiejetnoscAgenta = `INSERT INTO agent_umiejetnosc (agent_id, kod) VALUES (?, ?)
	                          ON CONFLICT(agent_id, kod) DO NOTHING`

	wstawKonektorAgenta = `INSERT INTO agent_konektor
	                       (kod, agent_id, nazwa, rodzaj, punkt_dostepu_id, konfiguracja, aktywny)
	                       VALUES (?, ?, ?, ?, ?, ?, ?)`

	zapiszUprawnienieAgenta = `INSERT INTO agent_uprawnienie (agent_id, grupa, zakres, przyznane)
	                           VALUES (?, ?, ?, ?)
	                           ON CONFLICT(agent_id, grupa, zakres) DO UPDATE SET przyznane = excluded.przyznane`
)

// Konektor nie ma wskazania konta; granicę niesie korzeń `agent` po stronie złączenia.
var konektorPoKodzie = `SELECT k.id, k.kod, k.agent_id, a.kod, k.nazwa, k.rodzaj,
                               k.punkt_dostepu_id, k.konfiguracja, k.aktywny, k.utworzono
                          FROM agent_konektor k JOIN agent a ON a.id = k.agent_id
                         WHERE k.kod = ? AND ` + warunekKontaEksperta

// DodajUmiejetnosc przypisuje ekspertowi umiejętność. Powtórzone przypisanie tej
// samej umiejętności nie jest błędem — kończy się tym samym stanem.
func (r *repozytoriumAgentow) DodajUmiejetnosc(ctx context.Context, kodAgenta, kodUmiejetnosci string) (Agent, error) {
	agent, err := r.PoKodzie(ctx, kodAgenta)
	if err != nil {
		return Agent{}, err
	}
	kod := strings.TrimSpace(kodUmiejetnosci)
	if kod == "" {
		return Agent{}, fmt.Errorf("dane: umiejętność eksperta %q wymaga kodu", kodAgenta)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawUmiejetnoscAgenta)
	if err != nil {
		return Agent{}, err
	}
	if _, err := polecenie.ExecContext(ctx, agent.ID, kod); err != nil {
		return Agent{}, fmt.Errorf("dane: nie można przypisać umiejętności %q ekspertowi %q: %w", kod, kodAgenta, err)
	}
	return r.PoKodzie(ctx, kodAgenta)
}

// DodajKonektor podłącza ekspertowi konektor i oddaje zapisany wiersz konektora wraz z nadanym numerem.
func (r *repozytoriumAgentow) DodajKonektor(ctx context.Context, konektor KonektorAgenta) (KonektorAgenta, error) {
	agent, err := r.PoKodzie(ctx, konektor.AgentKod)
	if err != nil {
		return KonektorAgenta{}, err
	}
	konfiguracja, err := konfiguracjaKonektora(konektor)
	if err != nil {
		return KonektorAgenta{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawKonektorAgenta)
	if err != nil {
		return KonektorAgenta{}, err
	}
	_, err = polecenie.ExecContext(ctx, konektor.Kod, agent.ID, konektor.Nazwa, konektor.Rodzaj,
		liczbaDoKolumny(konektor.PunktDostepuID), konfiguracja, liczbaLogiczna(konektor.Aktywny))
	if err != nil {
		return KonektorAgenta{}, fmt.Errorf("dane: nie można podłączyć konektora %q: %w", konektor.Nazwa, err)
	}
	return r.konektorPoKodzie(ctx, konektor.Kod)
}

// UstawUprawnienie przyznaje albo odbiera uprawnienie eksperta w jednej grupie zakresu i oddaje zmieniony wiersz.
func (r *repozytoriumAgentow) UstawUprawnienie(ctx context.Context, kodAgenta string,
	uprawnienie UprawnienieAgenta) (Agent, error) {

	agent, err := r.PoKodzie(ctx, kodAgenta)
	if err != nil {
		return Agent{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszUprawnienieAgenta)
	if err != nil {
		return Agent{}, err
	}
	_, err = polecenie.ExecContext(ctx, agent.ID, uprawnienie.Grupa,
		strings.TrimSpace(uprawnienie.Zakres), liczbaLogiczna(uprawnienie.Przyznane))
	if err != nil {
		return Agent{}, fmt.Errorf("dane: nie można zapisać uprawnienia %q eksperta %q: %w",
			uprawnienie.Grupa, kodAgenta, err)
	}
	return r.PoKodzie(ctx, kodAgenta)
}

// konektorPoKodzie odczytuje jeden konektor wraz z kodem jego eksperta na podstawie kodu konektora eksperta.
func (r *repozytoriumAgentow) konektorPoKodzie(ctx context.Context, kod string) (KonektorAgenta, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, konektorPoKodzie)
	if err != nil {
		return KonektorAgenta{}, err
	}
	var wpis KonektorAgenta
	var punkt sql.NullInt64
	var aktywny int
	err = polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)).Scan(&wpis.ID, &wpis.Kod, &wpis.AgentID, &wpis.AgentKod,
		&wpis.Nazwa, &wpis.Rodzaj, &punkt, &wpis.Konfiguracja, &aktywny, &wpis.Utworzono)
	if err != nil {
		return KonektorAgenta{}, fmt.Errorf("dane: nie można odczytać konektora %q: %w", kod, err)
	}
	wpis.PunktDostepuID = liczbaZKolumny(punkt)
	wpis.Aktywny = aktywny == 1
	return wpis, nil
}

// konfiguracjaKonektora pilnuje, żeby kolumna niosła poprawny JSON. Pusta
// wartość znaczy brak konfiguracji, nie błąd.
func konfiguracjaKonektora(konektor KonektorAgenta) (string, error) {
	tresc := strings.TrimSpace(konektor.Konfiguracja)
	if tresc == "" {
		return "{}", nil
	}
	if !json.Valid([]byte(tresc)) {
		return "", fmt.Errorf("dane: konfiguracja konektora %q nie jest poprawnym JSON", konektor.Nazwa)
	}
	return tresc, nil
}
