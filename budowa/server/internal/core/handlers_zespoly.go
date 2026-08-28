// Plik wpina cztery komendy obszaru team.* obsługujące zespoły ekspertów, czyli nazwane składy
// biblioteki modułu Agents. Jedno zdarzenie na cały obszar, tak samo jak przy agentach.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Zespoly jest portem trwałości zespołów ekspertów: nazwanych składów biblioteki modułu Agents tego konta.
type Zespoly interface {
	Zapisz(ctx context.Context, z shared.TeamSaveRequest) (shared.TeamSaveResponse, error)
	Wczytaj(ctx context.Context, z shared.TeamLoadRequest) (shared.TeamLoadResponse, error)
	Wykaz(ctx context.Context, z shared.TeamListRequest) (shared.TeamListResponse, error)
	Skopiuj(ctx context.Context, z shared.TeamDuplicateRequest) (shared.TeamDuplicateResponse, error)
}

// zarejestrujZespoly wpina cztery komendy obszaru team.* w rejestrze rdzenia tej platformy dla konta użytkownika.
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
				// Kopia jest założeniem, nie zmianą źródła; zdarzenie niesie zespół nowy.
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
