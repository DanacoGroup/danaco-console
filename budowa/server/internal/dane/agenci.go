// Odpowiedzialność pliku: odczyt biblioteki ekspertów modułu Agents (tabele
// `agent`, `agent_umiejetnosc`, `agent_konektor`, `agent_uprawnienie` z migracji
// 037–038). Zapis leży w `agenci_zapis.go`, powiązania w `agenci_powiazania.go`.
//
// Ekspert jest komponentem własnym. Nie jest sesją ani oknem: żyje dłużej niż
// jedno i drugie, a okna modułu Agents wyłącznie go opisują.
//
// Odczyt wykazu nie mnoży zapytań: trzy tabele podrzędne czytamy w całości raz
// na wywołanie i grupujemy po numerze eksperta, zamiast pytać o nie osobno dla
// każdego wiersza wykazu. Biblioteka ekspertów jest zbiorem rzędu dziesiątek,
// więc koszt jest stały i przewidywalny.
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
	// ImieWlasne to imię, które widzi Operator; `Nazwa` zostaje napisem
	// technicznym, po którym idzie porządek wykazu (migracja 073).
	ImieWlasne string
	// Favikon niesie odwołanie do znaku graficznego, nie jego bajty.
	Favikon string
	// UstawieniaJSON niesie treść pliku `--settings` nakładanego przez eksperta:
	// zaczepy (hooks) i reguły narzędzi (migracja 085). Pusty znaczy „ekspert nie
	// nakłada własnych ustawień" — okno rusza z ustawieniami sesji.
	//
	// Napis w całości, nie rozebrany na kolumny: kształt tego pliku należy do
	// CLI, nie do platformy, więc rdzeń go przenosi, a nie wersjonuje.
	UstawieniaJSON string
	// TrybNakladki rozstrzyga, czy instrukcja eksperta dopisuje się do globalnego
	// promptu systemowego ('DOLACZ' — opcja domyślna), czy staje się nim
	// ('ZASTAP' — odstępstwo od ustawień domyślnych, oznaczane przez Operatora
	// świadomie). Wartości kontraktu wprost (`shared.IdentityMode`), bo katalog
	// jest jeden i nie powtarzamy go tu drugi raz. Kolumna z migracji
	// 081; pusty napis nie występuje, bo kolumna jest NOT NULL z wartością
	// domyślną.
	TrybNakladki string
	// Widocznosc rozstrzyga, czy ekspert pokazuje się wszędzie ('global'), czy
	// wyłącznie w projektach, do których jest przypisany ('project') — wartości
	// kontraktu wprost (`shared.AgentVisibility`), kolumna z migracji 121.
	// Przynależności do projektu to pole nie niesie: tę zapisuje tabela
	// `przypisanie_agenta_projektu` (migracja 035) i tylko ona.
	Widocznosc string
	// PoziomyPamieci niesie poziomy pamięci z definicji eksperta (tabela
	// `agent_pamiec_poziom`, migracja 121) w porządku alfabetycznym.
	//
	// Wycinek pusty znaczy pamięć wyłączoną w całości — to jedyny zapis
	// wyłączenia i dlatego wyłączenie nie ma własnej wartości poziomu. Wartość
	// piąta obok czterech pozwoliłaby ułożyć stan sprzeczny („sesja" i „wyłączona"
	// naraz), którego nikt nie umiałby rozstrzygnąć bez zgadywania.
	PoziomyPamieci []string
	// LimitPodagentow to górna liczba jednoczesnych podagentów eksperta
	// (kolumna z migracji 276). Zero znaczy Subagent Network wyłączony — jedyny
	// zapis wyłączenia; wartość wyjściowa 15 jest maksimum technicznym platformy.
	LimitPodagentow int
	// ModulyZastosowania niesie kody modułów, w których ekspert pojawia się jako
	// wykonawca doraźny (tabela `agent_modul_zastosowania`, migracja 276).
	// Wycinek pusty znaczy BRAK OGRANICZENIA, a nie brak odpowiedzi.
	ModulyZastosowania []string
	// Umiejetnosci niesie kody umiejętności w porządku alfabetycznym.
	Umiejetnosci []string
	// Konektory niesie kody konektorów w porządku nazw.
	Konektory []string
	// Uprawnienia niesie cztery grupy zakresu wraz z zakresami szczegółowymi.
	Uprawnienia []UprawnienieAgenta
}

// UprawnienieAgenta to wiersz `agent_uprawnienie`. Zakres pusty obejmuje całą
// grupę.
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

// FiltrAgentow zawęża wykaz biblioteki. Pole puste znaczy „bez zawężenia”.
type FiltrAgentow struct {
	Fraza        string
	TylkoAktywne bool
	// KodProjektu zawęża wykaz do ekspertów WIDOCZNYCH w tym projekcie:
	// wszystkich `global` oraz tych `project`, które zostały do niego przypisane
	// (`przypisanie_agenta_projektu`, migracja 035). Pusty znaczy bibliotekę
	// w całości — okno Agent Buildera musi widzieć także ekspertów projektowych,
	// bo inaczej nie dałoby się ich poprawić.
	KodProjektu string
	// Granica ogranicza liczbę zwróconych wierszy; zero i wartości ujemne
	// znaczą „bez granicy”. Liczba wszystkich spełniających warunki wraca
	// osobno, żeby okno wiedziało, ile pozycji ucięto.
	Granica int
}

// RepozytoriumAgentow jest kontraktem biblioteki ekspertów.
type RepozytoriumAgentow interface {
	Lista(ctx context.Context, filtr FiltrAgentow) ([]Agent, int, error)
	PoKodzie(ctx context.Context, kod string) (Agent, error)
	Dodaj(ctx context.Context, agent Agent) (Agent, error)
	Aktualizuj(ctx context.Context, agent Agent) (Agent, error)
	Usun(ctx context.Context, kod string) (bool, error)
	// UstawPoziomyPamieci zastępuje komplet poziomów pamięci eksperta. Wycinek
	// pusty jest żądaniem wyłączenia pamięci, nie brakiem żądania — rozróżnienie
	// „pominięto pole" od „podano listę pustą" należy do warstwy wyższej.
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
	// Wartość `global` idzie parametrem, a nie literałem w zapytaniu: katalog
	// widoczności ma jedno miejsce — kontrakt.
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

// odczytajAgenta składa strukturę z jednego wiersza wyniku.
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
