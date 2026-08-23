// Odpowiedzialność pliku: wpięcie czterech komend obszaru `team.*` — zespołów
// ekspertów, czyli nazwanych składów biblioteki modułu Agents.
//
// Jedno zdarzenie na cały obszar, tak samo jak `agent.changed` w
// `handlers_agenci.go`. Kontrakt daje zespołom wyłącznie `team.changed`, więc
// zapisanie nowego składu, zmiana istniejącego i skopiowanie rozgłaszają się
// tą samą drogą, różniąc się wyłącznie rodzajem zmiany (ChangeKind). Okno
// składu odświeża się z jednej subskrypcji.
//
// Rodzaj zmiany rozstrzyga żądanie, nie odpowiedź. `team.save` bez
// identyfikatora zakłada zespół (`created`), z identyfikatorem zmienia
// istniejący (`updated`). Odpowiedź w obu przypadkach niesie ten sam kształt
// `Team`, więc po niej samej rozróżnić się tego nie da.
//
// Odczyt nie rozgłasza. `team.load` i `team.list` niczego nie zmieniają, więc
// nie mają czego ogłaszać — zdarzenie po odczycie byłoby zawiadomieniem
// o zmianie, której nie było.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Zespoly jest portem trwałości zespołów ekspertów.
type Zespoly interface {
	Zapisz(ctx context.Context, z shared.TeamSaveRequest) (shared.TeamSaveResponse, error)
	Wczytaj(ctx context.Context, z shared.TeamLoadRequest) (shared.TeamLoadResponse, error)
	Wykaz(ctx context.Context, z shared.TeamListRequest) (shared.TeamListResponse, error)
	Skopiuj(ctx context.Context, z shared.TeamDuplicateRequest) (shared.TeamDuplicateResponse, error)
}

// zarejestrujZespoly wpina cztery komendy obszaru `team.*`.
func zarejestrujZespoly(r *Rejestr, zespoly Zespoly, e *emiter) {
	if r == nil || zespoly == nil {
		return
	}
	r.Zarejestruj(shared.CommandTeamList, obsluz(zespoly.Wykaz))
	r.Zarejestruj(shared.CommandTeamLoad, obsluz(zespoly.Wczytaj))

	r.Zarejestruj(shared.CommandTeamSave,
		obsluz(func(ctx context.Context, z shared.TeamSaveRequest) (shared.TeamSaveResponse, error) {
			w, err := zespoly.Zapisz(ctx, z)
			if err == nil {
				e.zespol(rodzajZapisuZespolu(z), w.Team)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandTeamDuplicate,
		obsluz(func(ctx context.Context, z shared.TeamDuplicateRequest) (shared.TeamDuplicateResponse, error) {
			w, err := zespoly.Skopiuj(ctx, z)
			if err == nil {
				// Kopia jest założeniem, nie zmianą źródła. Zdarzenie niesie
				// zespół nowy — źródło się nie zmieniło i nie ma o czym mówić.
				e.zespol(shared.ChangeKindCreated, w.Team)
			}
			return w, err
		}))
}

// rodzajZapisuZespolu rozstrzyga rodzaj zmiany po treści żądania: brak
// identyfikatora znaczy założenie, obecny — zmianę istniejącego.
func rodzajZapisuZespolu(z shared.TeamSaveRequest) shared.ChangeKind {
	if z.TeamId == nil || *z.TeamId == "" {
		return shared.ChangeKindCreated
	}
	return shared.ChangeKindUpdated
}

// zespol rozgłasza zmianę zespołu ekspertów. Zespół jest bytem własnym
// Operatora, nie bytem sesji, więc zdarzenie idzie bez jej wskazania — tak samo
// jak zmiana eksperta.
func (e *emiter) zespol(zmiana shared.ChangeKind, z shared.Team) {
	e.wyslij(shared.EventTeamChanged, "", shared.TeamChangedEvent{Change: zmiana, Team: z})
}
