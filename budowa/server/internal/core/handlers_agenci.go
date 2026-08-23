// Odpowiedzialność pliku: wpięcie ośmiu komend obszaru `agent.*` — modułu
// Agents wraz z jego pięcioma oknami operacyjnymi (Agent Builder, Model
// Configuration, Skills Manager, Connectors Manager, Permissions Center).
//
// Cały obszar ma jedno zdarzenie. Kontrakt daje modułowi Agents wyłącznie
// `agent.changed`, więc zmiana modelu bazowego, przypisanie umiejętności,
// podłączenie konektora i zmiana uprawnienia rozgłaszają się tak samo jak
// zmiana tożsamości: rodzajem `updated` wraz z ekspertem po zmianie. Okna
// modułu odświeżają się z jednej subskrypcji, a nie z pięciu.
//
// Dwie komendy oddają co innego niż eksperta: `agent.connector.add` oddaje
// konektor, `agent.permission.set` — wykaz uprawnień. Zdarzenie ma nieść
// eksperta, więc obsługiwacz dobiera go portem `Pobierz`. Nieudany dobór nie
// wywraca komendy — zmiana już zaszła, gaśnie wyłącznie rozgłoszenie.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Agenci jest portem biblioteki ekspertów.
type Agenci interface {
	Utworz(ctx context.Context, z shared.AgentCreateRequest) (shared.AgentCreateResponse, error)
	Zmien(ctx context.Context, z shared.AgentUpdateRequest) (shared.AgentUpdateResponse, error)
	Wykaz(ctx context.Context, z shared.AgentListRequest) (shared.AgentListResponse, error)
	Usun(ctx context.Context, z shared.AgentDeleteRequest) (shared.AgentDeleteResponse, error)
	UstawModel(ctx context.Context, z shared.AgentModelSetRequest) (shared.AgentModelSetResponse, error)
	DodajUmiejetnosc(ctx context.Context, z shared.AgentSkillAddRequest) (shared.AgentSkillAddResponse, error)
	DodajKonektor(ctx context.Context, z shared.AgentConnectorAddRequest) (shared.AgentConnectorAddResponse, error)
	UstawUprawnienie(ctx context.Context, z shared.AgentPermissionSetRequest) (shared.AgentPermissionSetResponse, error)
	// Pobierz oddaje jednego eksperta. Służy rozgłoszeniu zmiany po komendach,
	// których wynik eksperta nie niesie.
	Pobierz(ctx context.Context, idEksperta string) (shared.Agent, error)
	// WarstwyEksperta dokłada tożsamość własną eksperta — warstwy promptu
	// i wtyczki. Osobnym portem być nie może: nowy port znaczyłby nowe pole
	// w Portach i nową linię w `kompozycja.go`, a moduł Agents wchodzi do
	// rejestru jednym wywołaniem (`handlers_agent_warstwy.go`).
	WarstwyEksperta
}

// zarejestrujAgentow wpina osiem komend modułu Agents i dokłada pięć komend
// tożsamości własnej eksperta z `handlers_agent_warstwy.go`.
func zarejestrujAgentow(r *Rejestr, agenci Agenci, e *emiter) {
	if r == nil || agenci == nil {
		return
	}
	r.Zarejestruj(shared.CommandAgentList, obsluz(agenci.Wykaz))

	r.Zarejestruj(shared.CommandAgentCreate,
		obsluz(func(ctx context.Context, z shared.AgentCreateRequest) (shared.AgentCreateResponse, error) {
			w, err := agenci.Utworz(ctx, z)
			if err == nil {
				e.agent(shared.ChangeKindCreated, w.Agent)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentUpdate,
		obsluz(func(ctx context.Context, z shared.AgentUpdateRequest) (shared.AgentUpdateResponse, error) {
			w, err := agenci.Zmien(ctx, z)
			if err == nil {
				e.agent(shared.ChangeKindUpdated, w.Agent)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentDelete,
		obsluz(func(ctx context.Context, z shared.AgentDeleteRequest) (shared.AgentDeleteResponse, error) {
			w, err := agenci.Usun(ctx, z)
			if err == nil && w.Deleted {
				e.agent(shared.ChangeKindDeleted, shared.Agent{Id: z.AgentId})
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentModelSet,
		obsluz(func(ctx context.Context, z shared.AgentModelSetRequest) (shared.AgentModelSetResponse, error) {
			w, err := agenci.UstawModel(ctx, z)
			if err == nil {
				e.agent(shared.ChangeKindUpdated, w.Agent)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentSkillAdd,
		obsluz(func(ctx context.Context, z shared.AgentSkillAddRequest) (shared.AgentSkillAddResponse, error) {
			w, err := agenci.DodajUmiejetnosc(ctx, z)
			if err == nil {
				e.agent(shared.ChangeKindUpdated, w.Agent)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentConnectorAdd,
		obsluz(func(ctx context.Context, z shared.AgentConnectorAddRequest) (shared.AgentConnectorAddResponse, error) {
			w, err := agenci.DodajKonektor(ctx, z)
			if err == nil {
				rozglosEksperta(ctx, agenci, e, z.AgentId)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentPermissionSet,
		obsluz(func(ctx context.Context, z shared.AgentPermissionSetRequest) (shared.AgentPermissionSetResponse, error) {
			w, err := agenci.UstawUprawnienie(ctx, z)
			if err == nil {
				rozglosEksperta(ctx, agenci, e, z.AgentId)
			}
			return w, err
		}))

	// Pięć komend tożsamości własnej eksperta idzie tym samym rejestrem
	// i tym samym emiterem — treść w `handlers_agent_warstwy.go`.
	zarejestrujWarstwyAgenta(r, agenci, e)
}

// rozglosEksperta dobiera eksperta i rozgłasza jego zmianę. Nieudany dobór
// kończy wyłącznie rozgłoszenie — komenda już się powiodła.
func rozglosEksperta(ctx context.Context, agenci Agenci, e *emiter, idEksperta string) {
	ekspert, err := agenci.Pobierz(ctx, idEksperta)
	if err != nil {
		return
	}
	e.agent(shared.ChangeKindUpdated, ekspert)
}

// agent rozgłasza zmianę eksperta. Ekspert jest komponentem własnym, nie bytem
// sesji, więc zdarzenie idzie bez jej wskazania — tak samo jak zmiana konta.
func (e *emiter) agent(zmiana shared.ChangeKind, a shared.Agent) {
	e.wyslij(shared.EventAgentChanged, "", shared.AgentChangedEvent{Change: zmiana, Agent: a})
}
