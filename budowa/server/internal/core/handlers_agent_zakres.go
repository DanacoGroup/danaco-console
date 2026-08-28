// Plik wpina dwanaście komend zakresu działania eksperta wraz z portem ZakresEksperta,
// który je wypełnia. Port jest osobny, bo wtopienie w port Agenci rozdęłoby go ponad
// czytelność.
package core

import (
	"context"

	"danacoconsole/shared"
)

// ZakresEksperta jest portem zakresu działania eksperta: Permissions Center
// w całości oraz cztery dopełnienia okien Skills, Connectors i historii wersji.
type ZakresEksperta interface {
	UsunUmiejetnosc(ctx context.Context, z shared.AgentSkillRemoveRequest) (shared.AgentSkillRemoveResponse, error)
	WykazKonektorow(ctx context.Context, z shared.AgentConnectorListRequest) (shared.AgentConnectorListResponse, error)
	UsunKonektor(ctx context.Context, z shared.AgentConnectorRemoveRequest) (shared.AgentConnectorRemoveResponse, error)
	KonfigurujKonektor(ctx context.Context, z shared.AgentConnectorConfigureRequest) (shared.AgentConnectorConfigureResponse, error)
	PodgladWersji(ctx context.Context, z shared.AgentVersionGetRequest) (shared.AgentVersionGetResponse, error)
	WykazPrzypisan(ctx context.Context, z shared.AgentAssignmentListRequest) (shared.AgentAssignmentListResponse, error)
	UstawModuly(ctx context.Context, z shared.AgentModulesSetRequest) (shared.AgentModulesSetResponse, error)
	OdczytajIzolacje(ctx context.Context, z shared.AgentIsolationGetRequest) (shared.AgentIsolationGetResponse, error)
	ZapiszIzolacje(ctx context.Context, z shared.AgentIsolationSetRequest) (shared.AgentIsolationSetResponse, error)
	UstawPodagentow(ctx context.Context, z shared.AgentSubagentSetRequest) (shared.AgentSubagentSetResponse, error)
	UsunUprawnienie(ctx context.Context, z shared.AgentPermissionRemoveRequest) (shared.AgentPermissionRemoveResponse, error)
	PolitykaEksperta(ctx context.Context, z shared.AgentPolicyGetRequest) (shared.AgentPolicyGetResponse, error)
	// EkspertPoZmianie oddaje eksperta rozgłoszeniu po komendach, których wynik eksperta nie
	// niesie.
	EkspertPoZmianie(ctx context.Context, idEksperta string) (shared.Agent, error)
}

// zarejestrujZakresEksperta wpina dwanaście komend zakresu działania eksperta w rejestr
// komend rdzenia.
func zarejestrujZakresEksperta(r *Rejestr, zakres ZakresEksperta, e *emiter) {
	if r == nil || zakres == nil {
		return
	}

	// Pięć odczytów bez zdarzenia.
	r.Zarejestruj(shared.CommandAgentConnectorList, obsluz(zakres.WykazKonektorow))
	r.Zarejestruj(shared.CommandAgentVersionGet, obsluz(zakres.PodgladWersji))
	r.Zarejestruj(shared.CommandAgentAssignmentList, obsluz(zakres.WykazPrzypisan))
	r.Zarejestruj(shared.CommandAgentIsolationGet, obsluz(zakres.OdczytajIzolacje))
	r.Zarejestruj(shared.CommandAgentPolicyGet, obsluz(zakres.PolitykaEksperta))

	// Cztery zmiany, których wynik niesie eksperta — rozgłoszenie idzie z niego.
	r.Zarejestruj(shared.CommandAgentSkillRemove,
		obsluz(func(ctx context.Context, z shared.AgentSkillRemoveRequest) (shared.AgentSkillRemoveResponse, error) {
			w, err := zakres.UsunUmiejetnosc(ctx, z)
			if err == nil {
				e.agent(shared.ChangeKindUpdated, w.Agent)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentConnectorRemove,
		obsluz(func(ctx context.Context, z shared.AgentConnectorRemoveRequest) (shared.AgentConnectorRemoveResponse, error) {
			w, err := zakres.UsunKonektor(ctx, z)
			if err == nil {
				e.agent(shared.ChangeKindUpdated, w.Agent)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentModulesSet,
		obsluz(func(ctx context.Context, z shared.AgentModulesSetRequest) (shared.AgentModulesSetResponse, error) {
			w, err := zakres.UstawModuly(ctx, z)
			if err == nil {
				e.agent(shared.ChangeKindUpdated, w.Agent)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentSubagentSet,
		obsluz(func(ctx context.Context, z shared.AgentSubagentSetRequest) (shared.AgentSubagentSetResponse, error) {
			w, err := zakres.UstawPodagentow(ctx, z)
			if err == nil {
				e.agent(shared.ChangeKindUpdated, w.Agent)
			}
			return w, err
		}))

	// Trzy zmiany, których wynik eksperta nie niesie — ekspert dobierany osobno.
	r.Zarejestruj(shared.CommandAgentConnectorConfigure,
		obsluz(func(ctx context.Context, z shared.AgentConnectorConfigureRequest) (shared.AgentConnectorConfigureResponse, error) {
			w, err := zakres.KonfigurujKonektor(ctx, z)
			if err == nil {
				rozglosEkspertaZakresu(ctx, zakres, e, z.AgentId)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentIsolationSet,
		obsluz(func(ctx context.Context, z shared.AgentIsolationSetRequest) (shared.AgentIsolationSetResponse, error) {
			w, err := zakres.ZapiszIzolacje(ctx, z)
			if err == nil {
				rozglosEkspertaZakresu(ctx, zakres, e, z.AgentId)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentPermissionRemove,
		obsluz(func(ctx context.Context, z shared.AgentPermissionRemoveRequest) (shared.AgentPermissionRemoveResponse, error) {
			w, err := zakres.UsunUprawnienie(ctx, z)
			if err == nil {
				rozglosEkspertaZakresu(ctx, zakres, e, z.AgentId)
			}
			return w, err
		}))
}

// rozglosEkspertaZakresu dobiera eksperta i rozgłasza jego zmianę. Nieudany
// dobór kończy wyłącznie rozgłoszenie — komenda już się powiodła, więc jej
// wynik nie zależy od powodzenia tego odczytu.
func rozglosEkspertaZakresu(ctx context.Context, zakres ZakresEksperta, e *emiter, idEksperta string) {
	ekspert, err := zakres.EkspertPoZmianie(ctx, idEksperta)
	if err != nil {
		return
	}
	e.agent(shared.ChangeKindUpdated, ekspert)
}
