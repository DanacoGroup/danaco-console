// Odpowiedzialność pliku: wpięcie pięciu komend rodziny `component.*` —
// komponentów własnych Strefy 2 Strony głównej — oraz rozgłoszenie zdarzenia
// `component.changed`. Adapter rodziny leży w `adapter_modul_komponenty.go`,
// a schemat rejestru w `store/migracja_059_komponenty_strony_glownej.sql`.
//
// Komend jest pięć, choć rodzina liczy sześć pozycji kontraktu:
// `component.changed` jest zdarzeniem, nie komendą — stoi w dziale `zdarzenia`
// kontraktu i niesie `payload`, nie parę żądanie/wynik. Nie rejestruje się go
// więc w rejestrze komend; wychodzi nadajnikiem po komendzie, która zmieniła
// stan.
//
// Na całą rodzinę przypada jedno zdarzenie, wzorem `agent.changed`
// i `workspace.project.changed`: założenie, zmiana, przypisanie i usunięcie
// rozgłaszają się tym samym zdarzeniem, różniąc się polem `change`. Strefa 2
// odświeża się z jednej subskrypcji.
//
// Usunięcie rozgłasza komponent z samym identyfikatorem, bo po usunięciu nie
// ma już czego dobrać z rejestru — tak samo robi `agent.changed` przy
// `agent.delete` (`handlers_agenci.go`). Kontrakt wymaga pola `component`
// w ładunku, więc idzie tam identyfikator bytu, który zniknął, a nie pusty
// kształt bez tożsamości.
//
// Przypisanie rozgłasza się tylko wtedy, gdy doszło do skutku. Powtórzone
// `component.assign` oddaje `assigned: false` i niczego nie zmienia, a zdarzenie
// zmiany po czynności, która nic nie zmieniła, byłoby fałszywym powiadomieniem.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Komponenty jest portem komponentów własnych Strefy 2 Strony głównej.
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
				e.komponent(shared.ChangeKindCreated, odpowiedz.Component)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandComponentUpdate,
		obsluz(func(ctx context.Context, z shared.ComponentUpdateRequest) (shared.ComponentUpdateResponse, error) {
			odpowiedz, err := k.Zmien(ctx, z)
			if err == nil {
				e.komponent(shared.ChangeKindUpdated, odpowiedz.Component)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandComponentDelete,
		obsluz(func(ctx context.Context, z shared.ComponentDeleteRequest) (shared.ComponentDeleteResponse, error) {
			odpowiedz, err := k.Usun(ctx, z)
			if err == nil && odpowiedz.Deleted {
				e.komponent(shared.ChangeKindDeleted, shared.Component{Id: z.ComponentId})
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandComponentList, obsluz(k.Wykaz))

	r.Zarejestruj(shared.CommandComponentAssign,
		obsluz(func(ctx context.Context, z shared.ComponentAssignRequest) (shared.ComponentAssignResponse, error) {
			odpowiedz, err := k.Przypisz(ctx, z)
			if err == nil && odpowiedz.Assigned {
				e.komponent(shared.ChangeKindUpdated, odpowiedz.Component)
			}
			return odpowiedz, err
		}))
}

// komponent rozgłasza zmianę komponentu własnego. Komponent nie jest bytem
// jednej karty sesji — jest kaflem Strony głównej, więc zdarzenie idzie bez
// wskazania sesji, tak samo jak `workspace.project.changed`.
func (e *emiter) komponent(zmiana shared.ChangeKind, k shared.Component) {
	e.wyslij(shared.EventComponentChanged, "",
		shared.ComponentChangedEvent{Change: zmiana, Component: k})
}
