// Dwanaście komend zakresu działania eksperta: Permissions Center,
// dopełnienie Skills i Connectors Managera, podgląd wersji i przypisania od
// strony eksperta. Rozstrzyga zakres działania, nie bramkę dostępu do
// aplikacji. Stan wyjściowy to pełny dostęp.
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

type adapterZakresuEksperta struct {
	zakres       dane.RepozytoriumZakresuAgenta
	biblioteka   dane.RepozytoriumAgentow
	historia     dane.RepozytoriumWersjiAgenta
	moduly       dane.RepozytoriumModulow
	punkty       dane.RepozytoriumPunktowDostepu
	rozstrzygacz *konfig.Rozstrzygacz
}

var _ ZakresEksperta = (*adapterZakresuEksperta)(nil)

func NowyPortZakresuEksperta(zakres dane.RepozytoriumZakresuAgenta,
	biblioteka dane.RepozytoriumAgentow, historia dane.RepozytoriumWersjiAgenta,
	moduly dane.RepozytoriumModulow, punkty dane.RepozytoriumPunktowDostepu,
	rozstrzygacz *konfig.Rozstrzygacz) *adapterZakresuEksperta {

	return &adapterZakresuEksperta{zakres: zakres, biblioteka: biblioteka, historia: historia,
		moduly: moduly, punkty: punkty, rozstrzygacz: rozstrzygacz}
}

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

func zakresTechnicznyZnany(zakres shared.IsolationTechnicalScope) bool {
	for _, pozycja := range punktyTechniczne {
		if pozycja.zakres == zakres {
			return true
		}
	}
	return false
}

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
	// Dziedziczenie liczy się wyłącznie na żądanie okna; bez wskazania okna
	// wraca sama wartość wyjściowa.
	if idOkna := strings.TrimSpace(wartoscTekstu(z.WindowId)); idOkna != "" && a.rozstrzygacz != nil {
		nadpisujacy := a.nadpiszZasiegiem(ctx, idOkna, polityka.TechnicalSwitches)
		if nadpisujacy != nil {
			polityka.OverriddenBy = nadpisujacy
		}
	}
	return shared.AgentPolicyGetResponse{Policy: polityka}, nil
}

func (a *adapterZakresuEksperta) nadpiszZasiegiem(ctx context.Context, idOkna string,
	przelaczniki []shared.IsolationTechnicalSwitch) *shared.ConfigScope {

	kontekst := konfig.Kontekst{Okno: idOkna, KontoOperatora: dane.KontoOperatora(ctx)}
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

func (a *adapterZakresuEksperta) EkspertPoZmianie(ctx context.Context,
	idEksperta string) (shared.Agent, error) {

	return a.ekspertZakresu(ctx, idEksperta)
}
