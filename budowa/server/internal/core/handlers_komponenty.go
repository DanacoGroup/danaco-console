// Plik wpina pięć komend rodziny `component.*`, obsługujących komponenty
// własne Strefy 2 Strony głównej, i rozgłasza jedno zdarzenie
// `component.changed`.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Komponenty jest portem komponentów własnych Strefy 2 Strony głównej,
// obsługującym wykaz, utworzenie, zmianę, usunięcie i przypisanie.
type Komponenty interface {
	Wykaz(ctx context.Context, z shared.ComponentListRequest) (shared.ComponentListResponse, error)
	Utworz(ctx context.Context, z shared.ComponentCreateRequest) (shared.ComponentCreateResponse, error)
	Zmien(ctx context.Context, z shared.ComponentUpdateRequest) (shared.ComponentUpdateResponse, error)
	Usun(ctx context.Context, z shared.ComponentDeleteRequest) (shared.ComponentDeleteResponse, error)
	Przypisz(ctx context.Context, z shared.ComponentAssignRequest) (shared.ComponentAssignResponse, error)
}

// zarejestrujKomponenty wpina pięć komend rodziny `component.*`. Rejestr jest
// mapą nazwa→obsługiwacz, więc kolejność wpięcia nie znaczy nic dla działania.
func zarejestrujKomponenty(r *Rejestr, k Komponenty, e *emiter) {
	if r == nil || k == nil {
		return
	}

	r.Zarejestruj(shared.CommandComponentCreate,
		obsluz(func(ctx context.Context, z shared.ComponentCreateRequest) (shared.ComponentCreateResponse, error) {
			odpowiedz, err := k.Utworz(ctx, z)
			if err == nil {
				e.komponent(ctx, shared.ChangeKindCreated, odpowiedz.Component)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandComponentUpdate,
		obsluz(func(ctx context.Context, z shared.ComponentUpdateRequest) (shared.ComponentUpdateResponse, error) {
			odpowiedz, err := k.Zmien(ctx, z)
			if err == nil {
				e.komponent(ctx, shared.ChangeKindUpdated, odpowiedz.Component)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandComponentDelete,
		obsluz(func(ctx context.Context, z shared.ComponentDeleteRequest) (shared.ComponentDeleteResponse, error) {
			odpowiedz, err := k.Usun(ctx, z)
			if err == nil && odpowiedz.Deleted {
				e.komponent(ctx, shared.ChangeKindDeleted, shared.Component{Id: z.ComponentId})
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandComponentList, obsluz(k.Wykaz))

	r.Zarejestruj(shared.CommandComponentAssign,
		obsluz(func(ctx context.Context, z shared.ComponentAssignRequest) (shared.ComponentAssignResponse, error) {
			odpowiedz, err := k.Przypisz(ctx, z)
			if err == nil && odpowiedz.Assigned {
				e.komponent(ctx, shared.ChangeKindUpdated, odpowiedz.Component)
			}
			return odpowiedz, err
		}))
}

// komponent rozgłasza zmianę komponentu własnego. Komponent nie jest bytem
// jednej karty sesji — jest kaflem Strony głównej, więc zdarzenie idzie bez
// wskazania sesji, tak samo jak `workspace.project.changed`.
func (e *emiter) komponent(ctx context.Context, zmiana shared.ChangeKind, k shared.Component) {
	e.wyslijDoKonta(ctx, shared.EventComponentChanged, "",
		shared.ComponentChangedEvent{Change: zmiana, Component: k})
}
