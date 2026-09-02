// Plik zapisuje bibliotekę ekspertów: założenie, zmianę tożsamości i usunięcie; nowo założony ekspert otrzymuje od razu
// pełny komplet czterech uprawnień przyznanych, a kolumna wersji rośnie z każdą zmianą tożsamości w tym samym poleceniu zapisu.
package dane

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

const (
	wstawAgenta = `INSERT INTO agent
	               (kod, nazwa, opis, instrukcje_systemowe, kanal_kod, model, transport,
	                parametry_json, wersja, aktywny, widocznosc, konto_id)
	               VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?, ` + WskazanieKonta + `)`

	aktualizujAgenta = `UPDATE agent
	                    SET nazwa = ?, opis = ?, instrukcje_systemowe = ?, kanal_kod = ?, model = ?,
	                        transport = ?, parametry_json = ?, aktywny = ?, widocznosc = ?,
	                        wersja = wersja + 1,
	                        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
	                    WHERE kod = ? AND ` + WarunekKonta

	usunAgenta = `DELETE FROM agent WHERE kod = ? AND ` + WarunekKonta

	// Poziomy pamięci zapisuje się kompletem, nie po jednym: żądanie niesie
	// zbiór, więc zdjęcie wszystkich i wstawienie podanych zapisuje zbiór pusty
	// (pamięć wyłączona) tą samą drogą co każdy inny.
	zdejmijPoziomyPamieci = `DELETE FROM agent_pamiec_poziom WHERE agent_id = ?`

	zapiszPoziomPamieci = `INSERT INTO agent_pamiec_poziom (agent_id, poziom) VALUES (?, ?)`
)

// grupyUprawnienWyjsciowe wylicza cztery grupy zakresu kontraktu. Lista stoi
// tutaj, a nie w migracji, bo wiersze zakłada się per ekspert, nie raz na bazę.
var grupyUprawnienWyjsciowe = []string{
	shared.AgentPermissionGroupFiles,
	shared.AgentPermissionGroupNetwork,
	shared.AgentPermissionGroupProcesses,
	shared.AgentPermissionGroupIntegrations,
}

func (r *repozytoriumAgentow) Dodaj(ctx context.Context, agent Agent) (Agent, error) {
	parametry, err := parametryAgenta(agent)
	if err != nil {
		return Agent{}, err
	}
	if strings.TrimSpace(agent.Nazwa) == "" {
		return Agent{}, fmt.Errorf("dane: ekspert %q wymaga nazwy", agent.Kod)
	}
	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawAgenta)
		if err != nil {
			return err
		}
		wynik, err := polecenie.ExecContext(ctx, agent.Kod, agent.Nazwa, agent.Opis,
			agent.InstrukcjeSystemowe, tekstDoKolumny(agent.KanalKod), tekstDoKolumny(agent.Model),
			tekstDoKolumny(agent.Transport), parametry, liczbaLogiczna(agent.Aktywny),
			widocznoscKolumny(agent.Widocznosc), KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można założyć eksperta %q: %w", agent.Kod, err)
		}
		numer, err := wynik.LastInsertId()
		if err != nil {
			return fmt.Errorf("dane: nieznany numer założonego eksperta %q: %w", agent.Kod, err)
		}
		if err := r.zasiejUprawnienia(ctx, transakcja, numer); err != nil {
			return err
		}
		return r.zapiszPoziomy(ctx, transakcja, numer, agent.PoziomyPamieci)
	})
	if err != nil {
		return Agent{}, err
	}
	return r.PoKodzie(ctx, agent.Kod)
}

func (r *repozytoriumAgentow) Aktualizuj(ctx context.Context, agent Agent) (Agent, error) {
	parametry, err := parametryAgenta(agent)
	if err != nil {
		return Agent{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, aktualizujAgenta)
	if err != nil {
		return Agent{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, agent.Nazwa, agent.Opis, agent.InstrukcjeSystemowe,
		tekstDoKolumny(agent.KanalKod), tekstDoKolumny(agent.Model), tekstDoKolumny(agent.Transport),
		parametry, liczbaLogiczna(agent.Aktywny), widocznoscKolumny(agent.Widocznosc), agent.Kod,
		KontoOperatora(ctx))
	if err != nil {
		return Agent{}, fmt.Errorf("dane: nie można zapisać eksperta %q: %w", agent.Kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return Agent{}, fmt.Errorf("dane: nieznana liczba zapisanych ekspertów %q: %w", agent.Kod, err)
	}
	if zmienione == 0 {
		return Agent{}, odmowaEksperta(ctx, r.zapytania, agent.Kod)
	}
	return r.PoKodzie(ctx, agent.Kod)
}

// Usun kasuje eksperta wraz z powiązaniami — klucze obce tabel podrzędnych mają
// klauzulę ON DELETE CASCADE.
func (r *repozytoriumAgentow) Usun(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunAgenta)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć eksperta %q: %w", kod, err)
	}
	liczba, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznana liczba usuniętych ekspertów %q: %w", kod, err)
	}
	if liczba == 0 {
		if err := odmowaEksperta(ctx, r.zapytania, kod); errors.Is(err, ErrKolizjaWiersza) {
			return false, err
		}
	}
	return liczba > 0, nil
}

func (r *repozytoriumAgentow) zasiejUprawnienia(ctx context.Context, transakcja *sql.Tx, agentID int64) error {
	polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszUprawnienieAgenta)
	if err != nil {
		return err
	}
	for _, grupa := range grupyUprawnienWyjsciowe {
		if _, err := polecenie.ExecContext(ctx, agentID, grupa, "", 1); err != nil {
			return fmt.Errorf("dane: nie można zapisać wyjściowego uprawnienia %q: %w", grupa, err)
		}
	}
	return nil
}

func (r *repozytoriumAgentow) UstawPoziomyPamieci(ctx context.Context,
	kodAgenta string, poziomy []string) (Agent, error) {

	biezacy, err := r.PoKodzie(ctx, kodAgenta)
	if err != nil {
		return Agent{}, err
	}
	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		return r.zapiszPoziomy(ctx, transakcja, biezacy.ID, poziomy)
	})
	if err != nil {
		return Agent{}, err
	}
	return r.PoKodzie(ctx, kodAgenta)
}

func (r *repozytoriumAgentow) zapiszPoziomy(ctx context.Context, transakcja *sql.Tx,
	agentID int64, poziomy []string) error {

	zdejmij, err := r.zapytania.wTransakcji(ctx, transakcja, zdejmijPoziomyPamieci)
	if err != nil {
		return err
	}
	if _, err := zdejmij.ExecContext(ctx, agentID); err != nil {
		return fmt.Errorf("dane: nie można zdjąć poziomów pamięci eksperta: %w", err)
	}
	if len(poziomy) == 0 {
		return nil
	}
	wstaw, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszPoziomPamieci)
	if err != nil {
		return err
	}
	zapisane := map[string]struct{}{}
	for _, poziom := range poziomy {
		if _, juz := zapisane[poziom]; juz {
			continue
		}
		zapisane[poziom] = struct{}{}
		if _, err := wstaw.ExecContext(ctx, agentID, poziom); err != nil {
			return fmt.Errorf("dane: nie można zapisać poziomu pamięci %q: %w", poziom, err)
		}
	}
	return nil
}

func widocznoscKolumny(widocznosc string) string {
	if strings.TrimSpace(widocznosc) == "" {
		return shared.AgentVisibilityGlobal
	}
	return widocznosc
}

func parametryAgenta(agent Agent) (string, error) {
	parametry := strings.TrimSpace(agent.ParametryJSON)
	if parametry == "" {
		return "{}", nil
	}
	if !json.Valid([]byte(parametry)) {
		return "", fmt.Errorf("dane: parametry wywołania eksperta %q nie są poprawnym JSON", agent.Kod)
	}
	return parametry, nil
}
