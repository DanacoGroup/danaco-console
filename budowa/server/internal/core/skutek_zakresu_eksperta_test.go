package core

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"danacoconsole/shared"
)

// bazaZakresuSprawdzianu otwiera drugie połączenie do bazy rdzenia, żeby
// zmierzyć skutek niezależnie od odpowiedzi komendy.
func bazaZakresuSprawdzianu(t *testing.T, katalog string) *sql.DB {
	t.Helper()

	polaczenie, err := sql.Open("sqlite", filepath.Join(katalog, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy sprawdzianu: %v", err)
	}
	t.Cleanup(func() { _ = polaczenie.Close() })
	return polaczenie
}

// ekspertZakresuSprawdzianu zakłada nowego eksperta o podanej nazwie w rdzeniu
// i oddaje jego kod, wykorzystywany dalej do zapytań o zakres działania.
func ekspertZakresuSprawdzianu(t *testing.T, zmontowany *Zmontowany,
	zycie context.Context, nazwa string) string {
	t.Helper()

	var wynik shared.AgentCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentCreate,
		shared.AgentCreateRequest{Name: nazwa}, &wynik)
	return wynik.Agent.Id
}

// liczbaWierszyZakresu wykonuje jedno zapytanie SQL i oddaje liczbę wierszy
// spełniających jego warunek.
func liczbaWierszyZakresu(t *testing.T, baza *sql.DB, zapytanie string, argumenty ...any) int {
	t.Helper()

	var liczba int
	if err := baza.QueryRow(zapytanie, argumenty...).Scan(&liczba); err != nil {
		t.Fatalf("zapytanie sprawdzianu nie powiodło się: %v", err)
	}
	return liczba
}

// TestSkutekModulowEkspertaOdmawiaNalozenia sprawdza, czy `agent.modules.set`
// zapisuje moduły zastosowania w bazie i czy ten zapis powstrzymuje nałożenie
// eksperta na okno modułu spoza jego zakresu.
func TestSkutekModulowEkspertaOdmawiaNalozenia(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	ekspert := ekspertZakresuSprawdzianu(t, zmontowany, zycie, "Ekspert kodu")

	var zapis shared.AgentModulesSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentModulesSet,
		shared.AgentModulesSetRequest{AgentId: ekspert, ModuleCodes: []string{"terminal"}}, &zapis)

	// Skutek w bazie, nie w kopercie.
	liczba := liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM agent_modul_zastosowania m JOIN agent a ON a.id = m.agent_id
		  WHERE a.kod = ? AND m.kod_modulu = 'terminal'`, ekspert)
	if liczba != 1 {
		t.Fatalf("moduł zastosowania nie doszedł do bazy: wierszy %d", liczba)
	}
	if len(zapis.Agent.ModuleCodes) != 1 || zapis.Agent.ModuleCodes[0] != "terminal" {
		t.Errorf("odpowiedź nie niesie modułów zastosowania: %v", zapis.Agent.ModuleCodes)
	}

	// Odmowa: ten sam ekspert nakładany na okno modułu spoza zakresu.
	straz := NowaStrazEksperta(zmontowany.dane.Agenci, zmontowany.dane.Moduly)
	obcy := numerModuluSprawdzianu(t, baza, "studio")
	if err := straz.SprawdzModulOkna(zycie, ekspert, obcy); err == nil {
		t.Fatal("straż wpuściła eksperta do modułu spoza jego zakresu — zawężenie jest napisem, nie regułą")
	}
	wlasny := numerModuluSprawdzianu(t, baza, "terminal")
	if err := straz.SprawdzModulOkna(zycie, ekspert, wlasny); err != nil {
		t.Fatalf("straż odmówiła w module, który JEST w zakresie eksperta: %v", err)
	}

	// Lista pusta znaczy brak ograniczenia — i straż ma wtedy milczeć.
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentModulesSet,
		shared.AgentModulesSetRequest{AgentId: ekspert, ModuleCodes: []string{}}, nil)
	if err := straz.SprawdzModulOkna(zycie, ekspert, obcy); err != nil {
		t.Fatalf("lista pusta zawęziła zamiast zdjąć ograniczenie: %v", err)
	}
}

// numerModuluSprawdzianu odczytuje z katalogu modułów klucz wiersza
// odpowiadający podanemu kodowi modułu i oddaje go wywołującemu.
func numerModuluSprawdzianu(t *testing.T, baza *sql.DB, kod string) int64 {
	t.Helper()

	var id int64
	if err := baza.QueryRow(`SELECT id FROM modul WHERE kod = ?`, kod).Scan(&id); err != nil {
		t.Fatalf("moduł %q nie istnieje w katalogu: %v", kod, err)
	}
	return id
}

// TestSkutekUprawnieniaModuluOdmawiaNalozenia mierzy drugą drogę zawężenia
// modułowego: wpis uprawnienia grupy `modules` z odebranym dostępem, i sprawdza,
// że ten wpis powstrzymuje nałożenie eksperta na moduł.
func TestSkutekUprawnieniaModuluOdmawiaNalozenia(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	ekspert := ekspertZakresuSprawdzianu(t, zmontowany, zycie, "Ekspert kontrolny")

	zakres := "studio"
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentPermissionSet,
		shared.AgentPermissionSetRequest{
			AgentId: ekspert, Group: shared.AgentPermissionGroupModules,
			Granted: false, Scope: &zakres,
		}, nil)

	odebrane := liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM agent_uprawnienie u JOIN agent a ON a.id = u.agent_id
		  WHERE a.kod = ? AND u.grupa = 'modules' AND u.zakres = 'studio' AND u.przyznane = 0`, ekspert)
	if odebrane != 1 {
		t.Fatalf("odebranie modułu nie doszło do bazy: wierszy %d", odebrane)
	}

	straz := NowaStrazEksperta(zmontowany.dane.Agenci, zmontowany.dane.Moduly)
	if err := straz.SprawdzModulOkna(zycie, ekspert, numerModuluSprawdzianu(t, baza, "studio")); err == nil {
		t.Fatal("straż wpuściła eksperta do modułu, który mu odebrano")
	}

	// agent.permission.remove przywraca stan bez ustawienia, więc straż
	// przestaje odmawiać dostępu.
	grupa := shared.AgentPermissionGroup(shared.AgentPermissionGroupModules)
	var zdjecie shared.AgentPermissionRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentPermissionRemove,
		shared.AgentPermissionRemoveRequest{AgentId: ekspert, Group: &grupa, Scope: &zakres}, &zdjecie)
	if zdjecie.Removed != 1 {
		t.Errorf("zdjęto %d wpisów zamiast jednego", zdjecie.Removed)
	}
	zostalo := liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM agent_uprawnienie u JOIN agent a ON a.id = u.agent_id
		  WHERE a.kod = ? AND u.grupa = 'modules'`, ekspert)
	if zostalo != 0 {
		t.Fatalf("wpis uprawnienia został w bazie po zdjęciu: wierszy %d", zostalo)
	}
	if err := straz.SprawdzModulOkna(zycie, ekspert, numerModuluSprawdzianu(t, baza, "studio")); err != nil {
		t.Fatalf("straż odmawia mimo zdjęcia zawężenia: %v", err)
	}
}

// TestSkutekSubagentNetworkOdmawiaPowolania mierzy zawężenie czwartej grupy
// zakresu: wyłączony Subagent Network ma POWSTRZYMAĆ powołanie podagentów.
func TestSkutekSubagentNetworkOdmawiaPowolania(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	ekspert := ekspertZakresuSprawdzianu(t, zmontowany, zycie, "Ekspert samotny")

	// Stan wyjściowy: piętnaście, maksimum techniczne platformy.
	granica := liczbaWierszyZakresu(t, baza, `SELECT limit_podagentow FROM agent WHERE kod = ?`, ekspert)
	if granica != 15 {
		t.Fatalf("stan wyjściowy granicy podagentów wynosi %d zamiast 15", granica)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentSubagentSet,
		shared.AgentSubagentSetRequest{AgentId: ekspert, Enabled: false}, nil)
	po := liczbaWierszyZakresu(t, baza, `SELECT limit_podagentow FROM agent WHERE kod = ?`, ekspert)
	if po != 0 {
		t.Fatalf("wyłączenie Subagent Network nie doszło do bazy: granica %d", po)
	}

	straz := NowaStrazEksperta(zmontowany.dane.Agenci, zmontowany.dane.Moduly)
	wartosc, dotyczy := straz.GranicaPodagentowEksperta(zycie, ekspert)
	if !dotyczy || wartosc != 0 {
		t.Fatalf("straż nie widzi wyłączenia: granica=%d dotyczy=%v", wartosc, dotyczy)
	}

	// Powrót z granicą własną, węższą niż platformowa.
	trzy := 3
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentSubagentSet,
		shared.AgentSubagentSetRequest{AgentId: ekspert, Enabled: true, Limit: &trzy}, nil)
	wartosc, _ = straz.GranicaPodagentowEksperta(zycie, ekspert)
	if wartosc != 3 {
		t.Fatalf("granica własna eksperta nie obowiązuje: %d", wartosc)
	}
}

// TestSkutekIzolacjiEkspertaWBazie mierzy zapis ośmiu zakresów izolacji
// technicznej i to, że odczyt oddaje komplet ośmiu także wtedy, gdy zapisany
// jest jeden.
func TestSkutekIzolacjiEkspertaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	ekspert := ekspertZakresuSprawdzianu(t, zmontowany, zycie, "Ekspert autonomiczny")

	var przed shared.AgentIsolationGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentIsolationGet,
		shared.AgentIsolationGetRequest{AgentId: ekspert}, &przed)
	if len(przed.Switches) != 8 {
		t.Fatalf("odczyt oddał %d przełączników zamiast ośmiu", len(przed.Switches))
	}
	for _, przelacznik := range przed.Switches {
		if przelacznik.Isolated {
			t.Fatalf("stan wyjściowy odcina zakres %s — a nie powinien odcinać żadnego", przelacznik.Scope)
		}
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentIsolationSet,
		shared.AgentIsolationSetRequest{AgentId: ekspert, Switches: []shared.IsolationTechnicalSwitch{
			{Scope: shared.IsolationTechnicalScopeNetworkAccess, Isolated: true},
		}}, nil)

	odciete := liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM agent_izolacja_techniczna i JOIN agent a ON a.id = i.agent_id
		  WHERE a.kod = ? AND i.zakres = ? AND i.odciety = 1`,
		ekspert, string(shared.IsolationTechnicalScopeNetworkAccess))
	if odciete != 1 {
		t.Fatalf("odcięcie dostępu sieciowego nie doszło do bazy: wierszy %d", odciete)
	}

	var po shared.AgentIsolationGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentIsolationGet,
		shared.AgentIsolationGetRequest{AgentId: ekspert}, &po)
	if len(po.Switches) != 8 {
		t.Fatalf("odczyt po zapisie oddał %d przełączników zamiast ośmiu", len(po.Switches))
	}
	odcietych := 0
	for _, przelacznik := range po.Switches {
		if przelacznik.Isolated {
			odcietych++
			if przelacznik.Scope != shared.IsolationTechnicalScopeNetworkAccess {
				t.Errorf("odcięto zakres, którego żądanie nie ruszało: %s", przelacznik.Scope)
			}
		}
	}
	if odcietych != 1 {
		t.Fatalf("po zapisie odciętych jest %d zamiast jednego", odcietych)
	}
}

// TestSkutekKonektorowEkspertaWBazie mierzy trzy komendy Connectors Managera:
// wykaz z definicją, konfigurację instancji i odłączenie.
func TestSkutekKonektorowEkspertaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	ekspert := ekspertZakresuSprawdzianu(t, zmontowany, zycie, "Ekspert integracji")

	var dodanie shared.AgentConnectorAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentConnectorAdd,
		shared.AgentConnectorAddRequest{
			AgentId: ekspert, Name: "Słownik terminologiczny",
			Kind: shared.AgentConnectorKindApi,
		}, &dodanie)

	var wykaz shared.AgentConnectorListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentConnectorList,
		shared.AgentConnectorListRequest{AgentId: ekspert}, &wykaz)
	if len(wykaz.Connectors) != 1 || wykaz.Connectors[0].Name != "Słownik terminologiczny" {
		t.Fatalf("wykaz konektorów nie niesie definicji: %+v", wykaz.Connectors)
	}

	wylaczony := false
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentConnectorConfigure,
		shared.AgentConnectorConfigureRequest{
			AgentId: ekspert, ConnectorId: dodanie.Connector.Id, Enabled: &wylaczony,
		}, nil)
	czynne := liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM agent_konektor WHERE kod = ? AND aktywny = 1`, dodanie.Connector.Id)
	if czynne != 0 {
		t.Fatalf("konfiguracja instancji nie zeszła do bazy: konektor nadal czynny")
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentConnectorRemove,
		shared.AgentConnectorRemoveRequest{AgentId: ekspert, ConnectorId: dodanie.Connector.Id}, nil)
	zostalo := liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM agent_konektor WHERE kod = ?`, dodanie.Connector.Id)
	if zostalo != 0 {
		t.Fatalf("konektor został w bazie po odłączeniu: wierszy %d", zostalo)
	}
}

// TestSkutekUmiejetnosciEkspertaWBazie mierzy zapis umiejętności w bazie po
// `agent.skill.add` oraz jej usunięcie po `agent.skill.remove`.
func TestSkutekUmiejetnosciEkspertaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	ekspert := ekspertZakresuSprawdzianu(t, zmontowany, zycie, "Ekspert redakcji")

	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentSkillAdd,
		shared.AgentSkillAddRequest{AgentId: ekspert, SkillId: "korekta"}, nil)
	if liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM agent_umiejetnosc u JOIN agent a ON a.id = u.agent_id
		  WHERE a.kod = ? AND u.kod = 'korekta'`, ekspert) != 1 {
		t.Fatal("umiejętność nie doszła do bazy")
	}

	var zdjecie shared.AgentSkillRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentSkillRemove,
		shared.AgentSkillRemoveRequest{AgentId: ekspert, SkillId: "korekta"}, &zdjecie)
	if liczbaWierszyZakresu(t, baza,
		`SELECT COUNT(*) FROM agent_umiejetnosc u JOIN agent a ON a.id = u.agent_id
		  WHERE a.kod = ? AND u.kod = 'korekta'`, ekspert) != 0 {
		t.Fatal("umiejętność została w bazie po zdjęciu")
	}
	// Powtórzone zdjęcie nie jest odmową, bo brak przypisania oznacza
	// definicję bez niego.
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentSkillRemove,
		shared.AgentSkillRemoveRequest{AgentId: ekspert, SkillId: "korekta"}, nil)
}

// TestSkutekPodgladuWersjiEksperta mierzy `agent.version.get`: czy migawka
// niesie tożsamość utrwaloną w TEJ wersji, a nie tożsamość bieżącą.
func TestSkutekPodgladuWersjiEksperta(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	ekspert := ekspertZakresuSprawdzianu(t, zmontowany, zycie, "Ekspert pierwotny")

	var historiaPrzed shared.AgentVersionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentVersionList,
		shared.AgentVersionListRequest{AgentId: ekspert}, &historiaPrzed)
	if len(historiaPrzed.Versions) == 0 {
		t.Fatal("ekspert po założeniu nie ma ani jednej wersji")
	}
	pierwsza := historiaPrzed.Versions[len(historiaPrzed.Versions)-1].Id

	nowaNazwa := "Ekspert poprawiony"
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentUpdate,
		shared.AgentUpdateRequest{AgentId: ekspert, Name: &nowaNazwa}, nil)

	var podglad shared.AgentVersionGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentVersionGet,
		shared.AgentVersionGetRequest{AgentId: ekspert, VersionId: pierwsza}, &podglad)
	if podglad.Snapshot.Name != "Ekspert pierwotny" {
		t.Fatalf("migawka niesie nazwę bieżącą zamiast utrwalonej: %q", podglad.Snapshot.Name)
	}
	if len(podglad.Snapshot.MemoryLevels) != 4 {
		t.Errorf("migawka niesie %d poziomów pamięci zamiast czterech wyjściowych",
			len(podglad.Snapshot.MemoryLevels))
	}
	if podglad.Snapshot.Visibility != shared.AgentVisibilityGlobal {
		t.Errorf("migawka niesie widoczność %q zamiast wyjściowej", podglad.Snapshot.Visibility)
	}
}

// TestSkutekPolitykiEksperta mierzy `agent.policy.get`: czy polityka efektywna
// odzwierciedla wszystkie cztery grupy zakresu.
func TestSkutekPolitykiEksperta(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	ekspert := ekspertZakresuSprawdzianu(t, zmontowany, zycie, "Ekspert polityki")

	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentModulesSet,
		shared.AgentModulesSetRequest{AgentId: ekspert, ModuleCodes: []string{"terminal"}}, nil)
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentIsolationSet,
		shared.AgentIsolationSetRequest{AgentId: ekspert, Switches: []shared.IsolationTechnicalSwitch{
			{Scope: shared.IsolationTechnicalScopeFileAccess, Isolated: true},
		}}, nil)
	piec := 5
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentSubagentSet,
		shared.AgentSubagentSetRequest{AgentId: ekspert, Enabled: true, Limit: &piec}, nil)

	var polityka shared.AgentPolicyGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentPolicyGet,
		shared.AgentPolicyGetRequest{AgentId: ekspert}, &polityka)

	if len(polityka.Policy.ModuleCodes) != 1 || polityka.Policy.ModuleCodes[0] != "terminal" {
		t.Errorf("polityka nie niesie modułów zastosowania: %v", polityka.Policy.ModuleCodes)
	}
	if len(polityka.Policy.TechnicalSwitches) != 8 {
		t.Errorf("polityka niesie %d zakresów izolacji zamiast ośmiu",
			len(polityka.Policy.TechnicalSwitches))
	}
	if !polityka.Policy.SubagentEnabled || polityka.Policy.SubagentLimit != 5 {
		t.Errorf("polityka niesie Subagent Network jako %v/%d",
			polityka.Policy.SubagentEnabled, polityka.Policy.SubagentLimit)
	}
	if polityka.Policy.OverriddenBy != nil {
		t.Errorf("polityka bez wskazania okna nie ma prawa wskazywać nadpisania: %v",
			*polityka.Policy.OverriddenBy)
	}
}

// TestSkutekPrzypisanEkspertaWBazie mierzy `agent.assignment.list` od strony
// eksperta: przypisanie zapisane przez Agent Managera modułu Workspace ma być
// widoczne stąd, bo inaczej biblioteka nie ma skąd wziąć licznika.
func TestSkutekPrzypisanEkspertaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaZakresuSprawdzianu(t, katalog)
	ekspert := ekspertZakresuSprawdzianu(t, zmontowany, zycie, "Ekspert projektowy")

	// Przypisanie zakładamy wprost w bazie: sprawdzian mierzy tu odczyt po
	// stronie eksperta.
	if _, err := baza.Exec(
		`INSERT INTO projekt (kod, nazwa) VALUES ('projekt-sprawdzianu', 'Projekt sprawdzianu')`); err != nil {
		t.Fatalf("nie można założyć projektu sprawdzianu: %v", err)
	}
	if _, err := baza.Exec(
		`INSERT INTO przypisanie_agenta_projektu (projekt_id, agent_kod, rola, domyslny_wykonawca)
		 VALUES ((SELECT id FROM projekt WHERE kod = 'projekt-sprawdzianu'), ?, 'redaktor', 1)`,
		ekspert); err != nil {
		t.Fatalf("nie można zapisać przypisania: %v", err)
	}

	var wykaz shared.AgentAssignmentListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAgentAssignmentList,
		shared.AgentAssignmentListRequest{AgentId: &ekspert}, &wykaz)
	if wykaz.Total != 1 || len(wykaz.Assignments) != 1 {
		t.Fatalf("wykaz przypisań liczy %d pozycji zamiast jednej", wykaz.Total)
	}
	przypisanie := wykaz.Assignments[0]
	if przypisanie.Kind != shared.AgentAssignmentKindProject ||
		przypisanie.TargetId != "projekt-sprawdzianu" {
		t.Errorf("przypisanie wskazuje %q/%q", przypisanie.Kind, przypisanie.TargetId)
	}
	if przypisanie.TargetName == nil || *przypisanie.TargetName != "Projekt sprawdzianu" {
		t.Errorf("przypisanie nie niesie nazwy projektu: %v", przypisanie.TargetName)
	}
}
