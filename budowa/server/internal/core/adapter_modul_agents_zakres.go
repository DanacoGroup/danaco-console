// Odpowiedzialność pliku: dwanaście komend zakresu działania eksperta —
// Permissions Center w całości (`agent.permission.remove`, `agent.modules.set`,
// `agent.isolation.get`, `agent.isolation.set`, `agent.subagent.set`,
// `agent.policy.get`), dopełnienie Skills i Connectors Managera
// (`agent.skill.remove`, `agent.connector.list`, `agent.connector.remove`,
// `agent.connector.configure`), podgląd wersji (`agent.version.get`)
// i przypisania widziane od strony eksperta (`agent.assignment.list`).
//
// PERMISSIONS CENTER NIE JEST BRAMKĄ DOSTĘPU DO APLIKACJI — i nie jest tu nią
// ani przez chwilę. Rozstrzyga zakres działania eksperta, którego Operator już
// uruchomił jako wykonawcę (agents.md rozdz. 9.1). Stanem wyjściowym jest pełny
// dostęp operacyjny: brak modułów znaczy „wszędzie", brak przełącznika izolacji
// znaczy „nieodcięty", brak wpisu uprawnienia znaczy „przyznane".
//
// ZAWĘŻENIE ZAPISANE JEST ZAWĘŻENIEM EGZEKWOWANYM. Zapis, którego rdzeń nie
// czyta przy wykonaniu, byłby gorszy niż jego brak: Operator widziałby
// ograniczenie, którego nikt nie pilnuje. Miejsca odmowy stoją
// w `straz_eksperta.go` — nałożenie eksperta na okno modułu poza jego zakresem
// i powołanie podagentów przez eksperta z wyłączonym Subagent Network.
//
// Wyliczenie polityki efektywnej niczego nie rozstrzyga i niczego nie zapisuje:
// jest przezroczystością, jak mówi kontrakt `agent.policy.get`. Dziedziczenie
// z poziomów zasięgu bierze ten sam rozstrzygacz, którym jedzie okno
// konfiguracji punktów izolacji — drugiej macierzy izolacji nie ma.
package core

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/shared"
)

// adapterZakresuEksperta wypełnia port ZakresEksperta.
type adapterZakresuEksperta struct {
	zakres     dane.RepozytoriumZakresuAgenta
	biblioteka dane.RepozytoriumAgentow
	historia   dane.RepozytoriumWersjiAgenta
	// moduly jest katalogiem modułów platformy. Służy wyłącznie sprawdzeniu
	// wskazania w `agent.modules.set` — rdzeń nie ma własnego wykazu modułów.
	moduly dane.RepozytoriumModulow
	// punkty są katalogiem punktów dostępu; konektor rodzaju `mcp` wskazuje
	// most z tego katalogu.
	punkty dane.RepozytoriumPunktowDostepu
	// rozstrzygacz odpowiada na pytanie o politykę izolacji obowiązującą
	// w oknie. Niewpięty daje politykę bez dziedziczenia — samą wartość
	// wyjściową eksperta, i tak mówi wtedy pole `overriddenBy`.
	rozstrzygacz *konfig.Rozstrzygacz
}

var _ ZakresEksperta = (*adapterZakresuEksperta)(nil)

// NowyPortZakresuEksperta wiąże port z repozytoriami zakresu działania.
func NowyPortZakresuEksperta(zakres dane.RepozytoriumZakresuAgenta,
	biblioteka dane.RepozytoriumAgentow, historia dane.RepozytoriumWersjiAgenta,
	moduly dane.RepozytoriumModulow, punkty dane.RepozytoriumPunktowDostepu,
	rozstrzygacz *konfig.Rozstrzygacz) *adapterZakresuEksperta {

	return &adapterZakresuEksperta{zakres: zakres, biblioteka: biblioteka, historia: historia,
		moduly: moduly, punkty: punkty, rozstrzygacz: rozstrzygacz}
}

// ── agent.skill.remove ───────────────────────────────────────────────────────

// UsunUmiejetnosc zdejmuje umiejętność z definicji eksperta. Brak przypisania
// znaczy definicję bez niego, nie błąd — tak stanowi kontrakt.
func (a *adapterZakresuEksperta) UsunUmiejetnosc(ctx context.Context,
	z shared.AgentSkillRemoveRequest) (shared.AgentSkillRemoveResponse, error) {

	if a == nil || a.zakres == nil {
		return shared.AgentSkillRemoveResponse{}, bladBrakuKatalogu("ekspertów")
	}
	if _, err := a.zakres.UsunUmiejetnoscAgenta(ctx, z.AgentId, strings.TrimSpace(z.SkillId)); err != nil {
		return shared.AgentSkillRemoveResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	ekspert, err := a.ekspertZakresu(ctx, z.AgentId)
	if err != nil {
		return shared.AgentSkillRemoveResponse{}, err
	}
	return shared.AgentSkillRemoveResponse{Agent: ekspert}, nil
}

// ── agent.connector.list ─────────────────────────────────────────────────────

// WykazKonektorow oddaje konektory eksperta wraz z ich definicją. Bez tej
// komendy `Agent.connectorIds` niesie same identyfikatory i okno nie zna ani
// nazwy, ani rodzaju, ani punktu dostępu.
func (a *adapterZakresuEksperta) WykazKonektorow(ctx context.Context,
	z shared.AgentConnectorListRequest) (shared.AgentConnectorListResponse, error) {

	if a == nil || a.zakres == nil {
		return shared.AgentConnectorListResponse{}, bladBrakuKatalogu("ekspertów")
	}
	wiersze, punkty, err := a.zakres.KonektoryAgenta(ctx, z.AgentId)
	if err != nil {
		return shared.AgentConnectorListResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	konektory := make([]shared.AgentConnector, 0, len(wiersze))
	for numer, wiersz := range wiersze {
		konektory = append(konektory, konektorKontraktu(wiersz, punkty[numer]))
	}
	return shared.AgentConnectorListResponse{Connectors: konektory}, nil
}

// ── agent.connector.remove ───────────────────────────────────────────────────

// UsunKonektor odłącza konektor od eksperta. Zakres wykorzystania zapisany
// w uprawnieniach ZOSTAJE — tak stanowi kontrakt, na wypadek ponownego
// podłączenia tego samego rozszerzenia (agents.md rozdz. 9.7).
func (a *adapterZakresuEksperta) UsunKonektor(ctx context.Context,
	z shared.AgentConnectorRemoveRequest) (shared.AgentConnectorRemoveResponse, error) {

	if a == nil || a.zakres == nil {
		return shared.AgentConnectorRemoveResponse{}, bladBrakuKatalogu("ekspertów")
	}
	if _, err := a.zakres.UsunKonektorAgenta(ctx, z.AgentId, strings.TrimSpace(z.ConnectorId)); err != nil {
		return shared.AgentConnectorRemoveResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	ekspert, err := a.ekspertZakresu(ctx, z.AgentId)
	if err != nil {
		return shared.AgentConnectorRemoveResponse{}, err
	}
	return shared.AgentConnectorRemoveResponse{Agent: ekspert}, nil
}

// ── agent.connector.configure ────────────────────────────────────────────────

// KonfigurujKonektor zapisuje konfigurację instancji konektora: adres serwera
// i parametry połączenia właściwe temu egzemplarzowi rozszerzenia.
func (a *adapterZakresuEksperta) KonfigurujKonektor(ctx context.Context,
	z shared.AgentConnectorConfigureRequest) (shared.AgentConnectorConfigureResponse, error) {

	if a == nil || a.zakres == nil {
		return shared.AgentConnectorConfigureResponse{}, bladBrakuKatalogu("ekspertów")
	}
	var punkt *int64
	if z.AccessPointId != nil {
		kod := strings.TrimSpace(*z.AccessPointId)
		numer, err := a.punktZakresu(ctx, kod)
		if err != nil {
			return shared.AgentConnectorConfigureResponse{}, err
		}
		punkt = numer
	}
	var konfiguracja *string
	if z.Config != nil {
		tresc, err := sprawdzParametry(z.Config, "konfiguracja integracji")
		if err != nil {
			return shared.AgentConnectorConfigureResponse{}, err
		}
		konfiguracja = &tresc
	}
	zapisany, kodPunktu, err := a.zakres.ZapiszKonektorAgenta(ctx, z.AgentId,
		strings.TrimSpace(z.ConnectorId), punkt, konfiguracja, z.Enabled)
	if err != nil {
		return shared.AgentConnectorConfigureResponse{}, bladWskazania(err, "konektor", z.ConnectorId)
	}
	return shared.AgentConnectorConfigureResponse{Connector: konektorKontraktu(zapisany, kodPunktu)}, nil
}

// punktZakresu przekłada kod punktu dostępu na numer wiersza. Kod pusty zdejmuje
// wskazanie mostu — konfiguracja instancji ma prawo je odebrać.
func (a *adapterZakresuEksperta) punktZakresu(ctx context.Context, kod string) (*int64, error) {
	if kod == "" {
		return nil, nil
	}
	if a.punkty == nil {
		return nil, bladBrakuKatalogu("punktów dostępu")
	}
	punkt, err := a.punkty.PoKodzie(ctx, kod)
	if err != nil {
		return nil, bladWskazania(err, "punkt dostępu", kod)
	}
	numer := punkt.ID
	return &numer, nil
}

// ── agent.version.get ────────────────────────────────────────────────────────

// PodgladWersji oddaje tożsamość utrwaloną w jednej wersji eksperta. Bez niej
// przywrócenie byłoby jedynym sposobem zobaczenia, co w wersji stało — czyli
// obejrzenie historii wymagałoby jej zmiany.
func (a *adapterZakresuEksperta) PodgladWersji(ctx context.Context,
	z shared.AgentVersionGetRequest) (shared.AgentVersionGetResponse, error) {

	if a == nil || a.historia == nil {
		return shared.AgentVersionGetResponse{}, bladBrakuKatalogu("historii ekspertów")
	}
	if strings.TrimSpace(z.VersionId) == "" {
		return shared.AgentVersionGetResponse{}, bladZadaniaEksperta("wskazanie czytanej wersji jest puste")
	}
	historia, err := a.historia.Wersje(ctx, z.AgentId)
	if err != nil {
		return shared.AgentVersionGetResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	for _, wiersz := range historia {
		if strconv.FormatInt(wiersz.Identyfikator, 10) != z.VersionId {
			continue
		}
		return shared.AgentVersionGetResponse{
			Version:  wersjaEkspertaKontraktu(z.AgentId, wiersz),
			Snapshot: migawkaWersjiKontraktu(wiersz),
		}, nil
	}
	return shared.AgentVersionGetResponse{},
		bladZadaniaEksperta("wersja " + z.VersionId + " nie należy do historii tego eksperta")
}

// migawkaWersjiKontraktu składa tożsamość utrwaloną w wersji.
//
// Pole `layers` zostaje puste i to jest odpowiedź, nie brak: `agent.layer.set`
// nie podnosi licznika wersji i nie zakłada migawki, więc warstwy wpisane tu
// byłyby warstwami BIEŻĄCYMI udającymi historię (migracja 281).
func migawkaWersjiKontraktu(w dane.WersjaAgenta) shared.AgentVersionSnapshot {
	return shared.AgentVersionSnapshot{
		Name:         w.Nazwa,
		Description:  wskaznikTekstu(w.Opis),
		SystemPrompt: wskaznikTekstu(w.InstrukcjeSystemowe),
		DisplayName:  wskaznikTekstu(w.ImieWlasne),
		Favicon:      wskaznikTekstu(w.Favikon),
		Mode:         trybWarstwyKontraktu(w.TrybNakladki),
		MemoryLevels: poziomyKontraktu(w.PoziomyPamieci),
		Visibility:   widocznoscKontraktu(w.Widocznosc),
	}
}

// ── agent.assignment.list ────────────────────────────────────────────────────

// WykazPrzypisan oddaje przypisania eksperta: projekty i role. Bez tej komendy
// biblioteka nie ma skąd wziąć liczby — `workspace.agent.assign` zapisuje
// przynależność, ale nikt nie odczytuje jej OD STRONY EKSPERTA.
func (a *adapterZakresuEksperta) WykazPrzypisan(ctx context.Context,
	z shared.AgentAssignmentListRequest) (shared.AgentAssignmentListResponse, error) {

	if a == nil || a.zakres == nil {
		return shared.AgentAssignmentListResponse{Assignments: []shared.AgentAssignment{}}, nil
	}
	rodzaj := ""
	if z.Kind != nil {
		rodzaj = string(*z.Kind)
		if rodzaj != dane.RodzajPrzypisaniaProjekt && rodzaj != dane.RodzajPrzypisaniaRola {
			return shared.AgentAssignmentListResponse{},
				bladZadaniaEksperta("rodzaj przypisania " + rodzaj + " nie należy do kontraktu")
		}
	}
	wiersze, err := a.zakres.PrzypisaniaEkspertow(ctx, strings.TrimSpace(wartoscTekstu(z.AgentId)), rodzaj)
	if err != nil {
		return shared.AgentAssignmentListResponse{}, err
	}
	przypisania := make([]shared.AgentAssignment, 0, len(wiersze))
	for _, wiersz := range wiersze {
		przypisanie := shared.AgentAssignment{
			AgentId:  wiersz.AgentKod,
			Kind:     shared.AgentAssignmentKind(wiersz.Rodzaj),
			TargetId: wiersz.CelKod,
		}
		przypisanie.TargetName = wskaznikTekstu(wiersz.CelNazwa)
		przypisanie.Role = wskaznikTekstu(wiersz.Rola)
		if wiersz.Rodzaj == dane.RodzajPrzypisaniaProjekt {
			domyslny := wiersz.DomyslnyWykonawca
			przypisanie.DefaultExecutor = &domyslny
		}
		if chwila, ok := chwilaZapisu(wiersz.Przypisano); ok {
			przypisanie.AssignedAt = chwila
		}
		przypisania = append(przypisania, przypisanie)
	}
	return shared.AgentAssignmentListResponse{Assignments: przypisania, Total: len(przypisania)}, nil
}

// ── agent.modules.set ────────────────────────────────────────────────────────

// UstawModuly zapisuje moduły zastosowania eksperta.
//
// LISTA PUSTA ZNACZY BRAK OGRANICZENIA i jest stanem wyjściowym — nie jest
// żądaniem odebrania ekspertowi wszystkiego. Ta sama zasada co przy poziomach
// pamięci, gdzie zbiór pusty też niesie znaczenie własne.
func (a *adapterZakresuEksperta) UstawModuly(ctx context.Context,
	z shared.AgentModulesSetRequest) (shared.AgentModulesSetResponse, error) {

	if a == nil || a.zakres == nil {
		return shared.AgentModulesSetResponse{}, bladBrakuKatalogu("ekspertów")
	}
	kody, err := a.kodyModulowZadania(ctx, z.ModuleCodes)
	if err != nil {
		return shared.AgentModulesSetResponse{}, err
	}
	if err := a.zakres.UstawModulyAgenta(ctx, z.AgentId, kody); err != nil {
		return shared.AgentModulesSetResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	ekspert, err := a.ekspertZakresu(ctx, z.AgentId)
	if err != nil {
		return shared.AgentModulesSetResponse{}, err
	}
	return shared.AgentModulesSetResponse{Agent: ekspert}, nil
}

// kodyModulowZadania sprawdza wskazania wobec katalogu modułów platformy.
//
// Kod spoza katalogu jest odmawiany, bo zawężenie do modułu, którego nie ma,
// nie zawęża niczego — a Operator zobaczyłby je w oknie jako obowiązujące.
// Przedrostek `module.` z rodziny komend jest zdejmowany: katalog trzyma sam
// kod modułu i to on jest prawdą.
func (a *adapterZakresuEksperta) kodyModulowZadania(ctx context.Context,
	wskazania []string) ([]string, error) {

	if len(wskazania) == 0 {
		return nil, nil
	}
	znane := map[string]struct{}{}
	if a.moduly != nil {
		wiersze, err := a.moduly.Lista(ctx)
		if err != nil {
			return nil, err
		}
		for _, wiersz := range wiersze {
			znane[wiersz.Kod] = struct{}{}
		}
	}
	kody := make([]string, 0, len(wskazania))
	for _, wskazanie := range wskazania {
		kod := strings.TrimPrefix(strings.TrimSpace(wskazanie), "module.")
		if kod == "" {
			continue
		}
		if len(znane) > 0 {
			if _, jest := znane[kod]; !jest {
				return nil, bladZadaniaEksperta("moduł " + kod + " nie należy do katalogu platformy")
			}
		}
		kody = append(kody, kod)
	}
	return kody, nil
}

// ── agent.isolation.get / agent.isolation.set ────────────────────────────────

// OdczytajIzolacje oddaje osiem zakresów izolacji technicznej eksperta.
// Komplet ośmiu wychodzi ZAWSZE, także gdy żaden nie został zmieniony —
// stanem wyjściowym jest `isolated: false` dla każdego zakresu.
func (a *adapterZakresuEksperta) OdczytajIzolacje(ctx context.Context,
	z shared.AgentIsolationGetRequest) (shared.AgentIsolationGetResponse, error) {

	if a == nil || a.zakres == nil {
		return shared.AgentIsolationGetResponse{}, bladBrakuKatalogu("ekspertów")
	}
	zapisane, err := a.zapisaneZakresyEksperta(ctx, z.AgentId)
	if err != nil {
		return shared.AgentIsolationGetResponse{}, err
	}
	return shared.AgentIsolationGetResponse{
		AgentId:  z.AgentId,
		Switches: przelacznikiEksperta(zapisane),
	}, nil
}

// ZapiszIzolacje zapisuje wskazane zakresy izolacji technicznej jako wartość
// wyjściową właściwą temu ekspertowi. Zakres pominięty w wykazie zostaje bez
// zmiany — okno przestawia jeden suwak, a nie komplet ośmiu.
func (a *adapterZakresuEksperta) ZapiszIzolacje(ctx context.Context,
	z shared.AgentIsolationSetRequest) (shared.AgentIsolationSetResponse, error) {

	if a == nil || a.zakres == nil {
		return shared.AgentIsolationSetResponse{}, bladBrakuKatalogu("ekspertów")
	}
	wpisy := make([]dane.PrzelacznikIzolacjiAgenta, 0, len(z.Switches))
	for _, przelacznik := range z.Switches {
		if !zakresTechnicznyZnany(przelacznik.Scope) {
			return shared.AgentIsolationSetResponse{},
				bladZadaniaEksperta("zakres izolacji " + string(przelacznik.Scope) +
					" nie należy do ośmiu zakresów kontraktu")
		}
		wpisy = append(wpisy, dane.PrzelacznikIzolacjiAgenta{
			Zakres: string(przelacznik.Scope), Odciety: przelacznik.Isolated,
		})
	}
	if err := a.zakres.UstawIzolacjeAgenta(ctx, z.AgentId, wpisy); err != nil {
		return shared.AgentIsolationSetResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	zapisane, err := a.zapisaneZakresyEksperta(ctx, z.AgentId)
	if err != nil {
		return shared.AgentIsolationSetResponse{}, err
	}
	return shared.AgentIsolationSetResponse{Switches: przelacznikiEksperta(zapisane)}, nil
}

// zapisaneZakresyEksperta czyta zapisane przełączniki jako mapę zakres → stan.
func (a *adapterZakresuEksperta) zapisaneZakresyEksperta(ctx context.Context,
	kodAgenta string) (map[string]bool, error) {

	if _, err := a.biblioteka.PoKodzie(ctx, kodAgenta); err != nil {
		return nil, bladWskazania(err, "ekspert", kodAgenta)
	}
	wiersze, err := a.zakres.IzolacjaAgenta(ctx, kodAgenta)
	if err != nil {
		return nil, bladWskazania(err, "ekspert", kodAgenta)
	}
	zapisane := make(map[string]bool, len(wiersze))
	for _, wiersz := range wiersze {
		zapisane[wiersz.Zakres] = wiersz.Odciety
	}
	return zapisane, nil
}

// przelacznikiEksperta składa komplet ośmiu w kolejności kontraktu.
// Kolejność bierze `punktyTechniczne` — ta sama, którą pokazuje okno
// konfiguracji punktów izolacji, bo macierz jest jedna.
func przelacznikiEksperta(zapisane map[string]bool) []shared.IsolationTechnicalSwitch {
	przelaczniki := make([]shared.IsolationTechnicalSwitch, 0, len(punktyTechniczne))
	for _, pozycja := range punktyTechniczne {
		przelaczniki = append(przelaczniki, shared.IsolationTechnicalSwitch{
			Scope:    pozycja.zakres,
			Isolated: zapisane[string(pozycja.zakres)],
		})
	}
	return przelaczniki
}

// zakresTechnicznyZnany sprawdza wskazanie wobec ośmiu zakresów kontraktu.
func zakresTechnicznyZnany(zakres shared.IsolationTechnicalScope) bool {
	for _, pozycja := range punktyTechniczne {
		if pozycja.zakres == zakres {
			return true
		}
	}
	return false
}

// ── agent.subagent.set ───────────────────────────────────────────────────────

// UstawPodagentow zapisuje dostępność Subagent Network i górną liczbę
// jednoczesnych podagentów.
//
// Zero w kolumnie znaczy „wyłączony" i jest jedynym zapisem wyłączenia, więc
// `enabled: false` zapisuje zero, a `enabled: true` bez granicy przywraca
// piętnaście — maksimum techniczne platformy, to samo, które zna
// `subagent.spawn`.
func (a *adapterZakresuEksperta) UstawPodagentow(ctx context.Context,
	z shared.AgentSubagentSetRequest) (shared.AgentSubagentSetResponse, error) {

	if a == nil || a.zakres == nil {
		return shared.AgentSubagentSetResponse{}, bladBrakuKatalogu("ekspertów")
	}
	granica := 0
	if z.Enabled {
		biezaca, err := a.zakres.GranicaPodagentow(ctx, z.AgentId)
		if err != nil {
			return shared.AgentSubagentSetResponse{}, bladWskazania(err, "ekspert", z.AgentId)
		}
		granica = biezaca
		if granica == 0 {
			granica = granicaPodagentow
		}
		if z.Limit != nil {
			if *z.Limit < 1 || *z.Limit > granicaPodagentow {
				return shared.AgentSubagentSetResponse{},
					bladZadaniaEksperta("granica podagentów jest liczbą od 1 do " +
						strconv.Itoa(granicaPodagentow))
			}
			granica = *z.Limit
		}
	}
	if err := a.zakres.UstawGranicePodagentow(ctx, z.AgentId, granica); err != nil {
		return shared.AgentSubagentSetResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	ekspert, err := a.ekspertZakresu(ctx, z.AgentId)
	if err != nil {
		return shared.AgentSubagentSetResponse{}, err
	}
	return shared.AgentSubagentSetResponse{Agent: ekspert}, nil
}

// ── agent.permission.remove ──────────────────────────────────────────────────

// UsunUprawnienie zdejmuje wpisy uprawnień, przywracając stan „brak ustawienia
// = wartość domyślna". `agent.permission.set` wyłącznie USTAWIA wartość wpisu,
// więc bez tej komendy wiersz zakresu zostawałby w wykazie na zawsze.
//
// Grupa pominięta zdejmuje wpisy wszystkich grup — tak działa „Resetuj do
// pełnego dostępu" jednym wywołaniem (agents.md rozdz. 9.4).
func (a *adapterZakresuEksperta) UsunUprawnienie(ctx context.Context,
	z shared.AgentPermissionRemoveRequest) (shared.AgentPermissionRemoveResponse, error) {

	if a == nil || a.zakres == nil {
		return shared.AgentPermissionRemoveResponse{}, bladBrakuKatalogu("ekspertów")
	}
	grupa := ""
	if z.Group != nil {
		if _, jest := grupyUprawnien[*z.Group]; !jest {
			return shared.AgentPermissionRemoveResponse{},
				bladZadaniaEksperta("grupa zakresu " + string(*z.Group) + " nie należy do kontraktu")
		}
		grupa = string(*z.Group)
	}
	zakres := strings.TrimSpace(wartoscTekstu(z.Scope))
	if zakres != "" && grupa == "" {
		return shared.AgentPermissionRemoveResponse{},
			bladZadaniaEksperta("zakres szczegółowy bez grupy nie wskazuje wpisu — podaj grupę")
	}
	zdjete, err := a.zakres.UsunUprawnieniaAgenta(ctx, z.AgentId, grupa, zakres)
	if err != nil {
		return shared.AgentPermissionRemoveResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	ekspert, err := a.ekspertZakresu(ctx, z.AgentId)
	if err != nil {
		return shared.AgentPermissionRemoveResponse{}, err
	}
	return shared.AgentPermissionRemoveResponse{
		Permissions: uprawnieniaEfektywne(ekspert.Permissions),
		Removed:     zdjete,
	}, nil
}

// uprawnieniaEfektywne oddaje wykaz uprawnień zawsze jako tablicę, nie nil:
// kontrakt oznacza pole jako wymagane, a `null` byłby trzecim stanem obok
// „są zawężenia" i „nie ma".
func uprawnieniaEfektywne(wpisy []shared.AgentPermission) []shared.AgentPermission {
	if wpisy == nil {
		return []shared.AgentPermission{}
	}
	return wpisy
}

// ── agent.policy.get ─────────────────────────────────────────────────────────

// PolitykaEksperta zwraca politykę efektywną: wynikowy zestaw ustawień po
// uwzględnieniu czterech grup zakresu. Transparentność, nie bramka — komenda
// niczego nie zapisuje i niczego nie rozstrzyga.
func (a *adapterZakresuEksperta) PolitykaEksperta(ctx context.Context,
	z shared.AgentPolicyGetRequest) (shared.AgentPolicyGetResponse, error) {

	if a == nil || a.zakres == nil {
		return shared.AgentPolicyGetResponse{}, bladBrakuKatalogu("ekspertów")
	}
	ekspert, err := a.ekspertZakresu(ctx, z.AgentId)
	if err != nil {
		return shared.AgentPolicyGetResponse{}, err
	}
	zapisane, err := a.zapisaneZakresyEksperta(ctx, z.AgentId)
	if err != nil {
		return shared.AgentPolicyGetResponse{}, err
	}
	polityka := shared.AgentPolicy{
		AgentId:           z.AgentId,
		WindowId:          z.WindowId,
		Permissions:       uprawnieniaEfektywne(ekspert.Permissions),
		TechnicalSwitches: przelacznikiEksperta(zapisane),
		ModuleCodes:       ekspert.ModuleCodes,
		SubagentEnabled:   ekspert.SubagentLimit > 0,
		SubagentLimit:     ekspert.SubagentLimit,
	}
	// Dziedziczenie liczy się wyłącznie na żądanie okna. Bez wskazania okna
	// odpowiedzią jest sama wartość wyjściowa eksperta i pole `overriddenBy`
	// zostaje puste — tak stanowi kontrakt.
	if idOkna := strings.TrimSpace(wartoscTekstu(z.WindowId)); idOkna != "" && a.rozstrzygacz != nil {
		nadpisujacy := a.nadpiszZasiegiem(idOkna, polityka.TechnicalSwitches)
		if nadpisujacy != nil {
			polityka.OverriddenBy = nadpisujacy
		}
	}
	return shared.AgentPolicyGetResponse{Policy: polityka}, nil
}

// nadpiszZasiegiem nakłada na przełączniki eksperta reguły zapisane na
// poziomach zasięgu i oddaje poziom, który nadpisał wartość wyjściową.
//
// Rozstrzygacz jest ten sam, którym jedzie okno konfiguracji punktów izolacji —
// macierz izolacji technicznej jest jedna i drugiej być nie może. Poziom
// zwracany jest ten NAJWĘŻSZY spośród nadpisujących, bo to on obowiązuje.
func (a *adapterZakresuEksperta) nadpiszZasiegiem(idOkna string,
	przelaczniki []shared.IsolationTechnicalSwitch) *shared.ConfigScope {

	kontekst := konfig.Kontekst{Okno: idOkna}
	var nadpisujacy *shared.ConfigScope
	for numer, pozycja := range punktyTechniczne {
		if numer >= len(przelaczniki) {
			break
		}
		wynik := a.rozstrzygacz.Rozstrzygnij(kontekst, pozycja.punkt.klucz)
		if wynik.Pochodzenie != konfig.PochodzenieZapis {
			continue
		}
		przelaczniki[numer].Isolated = odcietyPunkt(pozycja.punkt, wynik.Wartosc)
		poziom := wynik.Poziom
		nadpisujacy = &poziom
	}
	return nadpisujacy
}

// ── wspólne ──────────────────────────────────────────────────────────────────

// ekspertZakresu dobiera eksperta w kształcie kontraktu.
func (a *adapterZakresuEksperta) ekspertZakresu(ctx context.Context, kod string) (shared.Agent, error) {
	if a.biblioteka == nil {
		return shared.Agent{}, bladBrakuKatalogu("ekspertów")
	}
	wiersz, err := a.biblioteka.PoKodzie(ctx, kod)
	if err != nil {
		return shared.Agent{}, bladWskazania(err, "ekspert", kod)
	}
	ekspert := ekspertKontraktu(wiersz)
	sort.Strings(ekspert.ModuleCodes)
	return ekspert, nil
}

// EkspertPoZmianie oddaje eksperta rozgłoszeniu po komendach, których wynik
// eksperta nie niesie.
func (a *adapterZakresuEksperta) EkspertPoZmianie(ctx context.Context,
	idEksperta string) (shared.Agent, error) {

	return a.ekspertZakresu(ctx, idEksperta)
}
