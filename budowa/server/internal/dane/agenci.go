// Plik odczytuje bibliotekę ekspertów modułu Agents z tabel agent, agent_umiejetnosc, agent_konektor
// oraz agent_uprawnienie, grupując wyniki trzech tabel podrzędnych po numerze eksperta w jednym przebiegu zapytań.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// Agent to wiersz tabeli `agent` wraz z powiązaniami. `Kod` jest identyfikatorem
// trwałym i odpowiada polu `Agent.id` kontraktu.
type Agent struct {
	ID                  int64
	Kod                 string
	Nazwa               string
	Opis                string
	InstrukcjeSystemowe string
	KanalKod            *string
	Model               *string
	Transport           *string
	ParametryJSON       string
	Wersja              int
	Aktywny             bool
	Utworzono           string
	Zaktualizowano      string
	// ImieWlasne to imię widoczne dla Operatora; `Nazwa` pozostaje napisem technicznym porządku wykazu.
	ImieWlasne string
	Favikon string
	// UstawieniaJSON niesie treść pliku `--settings` eksperta; pusty znaczy brak ustawień własnych.
	UstawieniaJSON string
	// TrybNakladki rozstrzyga, czy instrukcja eksperta dopisuje się do promptu, czy go zastępuje.
	TrybNakladki string
	// Widocznosc rozstrzyga, czy ekspert jest widoczny wszędzie, czy tylko w przypisanych projektach.
	Widocznosc string
	// PoziomyPamieci niesie poziomy pamięci eksperta w porządku alfabetycznym; pusty wyłącza pamięć.
	PoziomyPamieci []string
	// LimitPodagentow to górna liczba jednoczesnych podagentów eksperta; zero wyłącza Subagent Network.
	LimitPodagentow int
	// ModulyZastosowania niesie kody modułów, w których ekspert działa jako wykonawca doraźny.
	ModulyZastosowania []string
	// Umiejetnosci niesie kody umiejętności w porządku alfabetycznym.
	Umiejetnosci []string
	// Konektory niesie kody konektorów w porządku nazw.
	Konektory []string
	// Uprawnienia niesie cztery grupy zakresu wraz z zakresami szczegółowymi.
	Uprawnienia []UprawnienieAgenta
}

// UprawnienieAgenta to wiersz tabeli `agent_uprawnienie`, w którym zakres pusty obejmuje całą grupę uprawnień eksperta.
type UprawnienieAgenta struct {
	Grupa     string
	Zakres    string
	Przyznane bool
}

// KonektorAgenta to wiersz `agent_konektor`. `PunktDostepuID` wskazuje most
// z katalogu punktów dostępu — ten sam, z którego rdzeń składa `mcpServers`.
type KonektorAgenta struct {
	ID             int64
	Kod            string
	AgentID        int64
	AgentKod       string
	Nazwa          string
	Rodzaj         string
	PunktDostepuID *int64
	Konfiguracja   string
	Aktywny        bool
	Utworzono      string
}

// FiltrAgentow zawęża wykaz biblioteki ekspertów według projektu i limitu wierszy; pole puste oznacza brak zawężenia.
type FiltrAgentow struct {
	Fraza        string
	TylkoAktywne bool
	// KodProjektu zawęża wykaz do ekspertów widocznych w tym projekcie; pusty zwraca całą bibliotekę.
	KodProjektu string
	// Granica ogranicza liczbę zwróconych wierszy; zero i wartości ujemne oznaczają brak granicy.
	Granica int
}

// RepozytoriumAgentow jest kontraktem biblioteki ekspertów, obejmującym odczyt wykazu, odczyt pojedynczego eksperta oraz zapis poziomów pamięci.
type RepozytoriumAgentow interface {
	Lista(ctx context.Context, filtr FiltrAgentow) ([]Agent, int, error)
	PoKodzie(ctx context.Context, kod string) (Agent, error)
	Dodaj(ctx context.Context, agent Agent) (Agent, error)
	Aktualizuj(ctx context.Context, agent Agent) (Agent, error)
	Usun(ctx context.Context, kod string) (bool, error)
	// UstawPoziomyPamieci zastępuje komplet poziomów pamięci eksperta; wycinek pusty wyłącza pamięć.
	UstawPoziomyPamieci(ctx context.Context, kodAgenta string, poziomy []string) (Agent, error)
	DodajUmiejetnosc(ctx context.Context, kodAgenta, kodUmiejetnosci string) (Agent, error)
	DodajKonektor(ctx context.Context, konektor KonektorAgenta) (KonektorAgenta, error)
	UstawUprawnienie(ctx context.Context, kodAgenta string, uprawnienie UprawnienieAgenta) (Agent, error)
}

const (
	kolumnyAgenta = `id, kod, nazwa, opis, instrukcje_systemowe, kanal_kod, model, transport,
	                 parametry_json, wersja, aktywny, utworzono, zaktualizowano,
	                 imie_wlasne, favikon, ustawienia_json, tryb_nakladki, widocznosc,
	                 limit_podagentow`

	// Wykaz czynnych pomija archiwum (migracja 102). Bez tego warunku ekspert
	// odłożony nadal wisiałby na liście, a archiwizacja byłaby znacznikiem bez
	// skutku — czyli atrapą. Archiwum ma własną komendę `agent.archive.list`.
	listaAgentow = `SELECT ` + kolumnyAgenta + ` FROM agent
	                WHERE zarchiwizowano_o IS NULL
	                  AND (? = 0 OR aktywny = 1)
	                  AND (? = '' OR lower(nazwa) LIKE ? OR lower(opis) LIKE ?)
	                  AND (? = '' OR widocznosc = ? OR EXISTS (
	                          SELECT 1 FROM przypisanie_agenta_projektu pap
	                            JOIN projekt p ON p.id = pap.projekt_id
	                           WHERE pap.agent_kod = agent.kod AND p.kod = ?))
	                ORDER BY nazwa, kod`

	agentPoKodzie = `SELECT ` + kolumnyAgenta + ` FROM agent WHERE kod = ?`
)

type repozytoriumAgentow struct {
	zapytania *zapytania
	db        *sql.DB
}

var _ RepozytoriumAgentow = (*repozytoriumAgentow)(nil)

func noweRepozytoriumAgentow(z *zapytania, db *sql.DB) *repozytoriumAgentow {
	return &repozytoriumAgentow{zapytania: z, db: db}
}

// Lista zwraca ekspertów spełniających warunki wraz z ich liczbą przed ucięciem
// granicą. Biblioteka pusta nie jest błędem — okno pokazuje wtedy stan pusty
// i zachętę do założenia pierwszego eksperta.
func (r *repozytoriumAgentow) Lista(ctx context.Context, filtr FiltrAgentow) ([]Agent, int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaAgentow)
	if err != nil {
		return nil, 0, err
	}
	fraza := strings.ToLower(strings.TrimSpace(filtr.Fraza))
	wzorzec := "%" + fraza + "%"
	// Wartość `global` idzie parametrem, nie literałem: widoczność ma jedno miejsce w kontrakcie.
	wiersze, err := polecenie.QueryContext(ctx, liczbaLogiczna(filtr.TylkoAktywne), fraza,
		wzorzec, wzorzec, filtr.KodProjektu, shared.AgentVisibilityGlobal, filtr.KodProjektu)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać biblioteki ekspertów: %w", err)
	}
	defer wiersze.Close()

	wszyscy := []Agent{}
	for wiersze.Next() {
		agent, err := odczytajAgenta(wiersze)
		if err != nil {
			return nil, 0, err
		}
		wszyscy = append(wszyscy, agent)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt biblioteki ekspertów: %w", err)
	}
	if err := r.dolaczPowiazania(ctx, wszyscy); err != nil {
		return nil, 0, err
	}
	razem := len(wszyscy)
	if filtr.Granica > 0 && filtr.Granica < razem {
		wszyscy = wszyscy[:filtr.Granica]
	}
	return wszyscy, razem, nil
}

// PoKodzie zwraca jednego eksperta wraz z powiązaniami. Brak wiersza wraca jako
// ErrBrakWiersza, żeby warstwa wyższa odróżniła „nie ma” od „odczyt padł”.
func (r *repozytoriumAgentow) PoKodzie(ctx context.Context, kod string) (Agent, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, agentPoKodzie)
	if err != nil {
		return Agent{}, err
	}
	agent, err := odczytajAgenta(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return Agent{}, fmt.Errorf("dane: ekspert %q nie istnieje: %w", kod, ErrBrakWiersza)
	}
	if err != nil {
		return Agent{}, fmt.Errorf("dane: nie można odczytać eksperta %q: %w", kod, err)
	}
	jeden := []Agent{agent}
	if err := r.dolaczPowiazania(ctx, jeden); err != nil {
		return Agent{}, err
	}
	return jeden[0], nil
}

// odczytajAgenta składa strukturę agenta z jednego wiersza wyniku zapytania, odczytując kolejno wszystkie jego kolumny.
func odczytajAgenta(wiersz skaner) (Agent, error) {
	var agent Agent
	var kanal, model, transport sql.NullString
	var aktywny int
	err := wiersz.Scan(&agent.ID, &agent.Kod, &agent.Nazwa, &agent.Opis,
		&agent.InstrukcjeSystemowe, &kanal, &model, &transport, &agent.ParametryJSON,
		&agent.Wersja, &aktywny, &agent.Utworzono, &agent.Zaktualizowano,
		&agent.ImieWlasne, &agent.Favikon, &agent.UstawieniaJSON, &agent.TrybNakladki,
		&agent.Widocznosc, &agent.LimitPodagentow)
	if err != nil {
		return Agent{}, err
	}
	agent.KanalKod = tekstZKolumny(kanal)
	agent.Model = tekstZKolumny(model)
	agent.Transport = tekstZKolumny(transport)
	agent.Aktywny = aktywny == 1
	return agent, nil
}
