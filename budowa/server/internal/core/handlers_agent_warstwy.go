// Plik wpina pięć komend tożsamości własnej eksperta wraz z portem WarstwyEksperta, który
// je wypełnia i wtapia w port Agenci zamiast podawać osobno.
package core

import (
	"context"

	"danacoconsole/shared"
)

// WarstwyEksperta jest portem tożsamości własnej eksperta: warstw jego promptu
// i jego wtyczek. Wtopiony jest w port Agenci — wypełnia go
// ten sam adapter, który obsługuje bibliotekę ekspertów.
type WarstwyEksperta interface {
	UstawWarstwe(ctx context.Context, z shared.AgentLayerSetRequest) (shared.AgentLayerSetResponse, error)
	UsunWarstwe(ctx context.Context, z shared.AgentLayerRemoveRequest) (shared.AgentLayerRemoveResponse, error)
	DodajWtyczke(ctx context.Context, z shared.AgentPluginAddRequest) (shared.AgentPluginAddResponse, error)
	UsunWtyczke(ctx context.Context, z shared.AgentPluginRemoveRequest) (shared.AgentPluginRemoveResponse, error)
	// WykazWtyczek oddaje wtyczki eksperta wraz z ich definicją, nie tylko identyfikatorami.
	WykazWtyczek(ctx context.Context, z shared.AgentPluginListRequest) (shared.AgentPluginListResponse, error)
	// EkspertPelny oddaje eksperta wraz z warstwami i wtyczkami, dla rozgłoszenia po dodaniu
	// wtyczki.
	EkspertPelny(ctx context.Context, idEksperta string) (shared.Agent, error)
}

// zarejestrujWarstwyAgenta wpina pięć komend tożsamości własnej eksperta:
// cztery zapisujące i jedną odczytującą.
// Wywołuje ją `zarejestrujAgentow` — rejestr i emiter są te same.
func zarejestrujWarstwyAgenta(r *Rejestr, warstwy WarstwyEksperta, e *emiter) {
	if r == nil || warstwy == nil {
		return
	}

	r.Zarejestruj(shared.CommandAgentLayerSet,
		obsluz(func(ctx context.Context, z shared.AgentLayerSetRequest) (shared.AgentLayerSetResponse, error) {
			w, err := warstwy.UstawWarstwe(ctx, z)
			if err == nil {
				e.agent(shared.ChangeKindUpdated, w.Agent)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentLayerRemove,
		obsluz(func(ctx context.Context, z shared.AgentLayerRemoveRequest) (shared.AgentLayerRemoveResponse, error) {
			w, err := warstwy.UsunWarstwe(ctx, z)
			if err == nil {
				e.agent(shared.ChangeKindUpdated, w.Agent)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentPluginAdd,
		obsluz(func(ctx context.Context, z shared.AgentPluginAddRequest) (shared.AgentPluginAddResponse, error) {
			w, err := warstwy.DodajWtyczke(ctx, z)
			if err == nil {
				rozglosEkspertaPelnego(ctx, warstwy, e, z.AgentId)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentPluginRemove,
		obsluz(func(ctx context.Context, z shared.AgentPluginRemoveRequest) (shared.AgentPluginRemoveResponse, error) {
			w, err := warstwy.UsunWtyczke(ctx, z)
			if err == nil {
				e.agent(shared.ChangeKindUpdated, w.Agent)
			}
			return w, err
		}))

	// Odczyt bez zdarzenia: wykaz niczego nie zmienia, więc nie ma czego rozgłaszać.
	r.Zarejestruj(shared.CommandAgentPluginList, obsluz(warstwy.WykazWtyczek))
}

// rozglosEkspertaPelnego dobiera eksperta wraz z warstwami i rozgłasza jego zmianę, wzorem
// rozglosEksperta.
func rozglosEkspertaPelnego(ctx context.Context, warstwy WarstwyEksperta, e *emiter, idEksperta string) {
	ekspert, err := warstwy.EkspertPelny(ctx, idEksperta)
	if err != nil {
		return
	}
	e.agent(shared.ChangeKindUpdated, ekspert)
}
