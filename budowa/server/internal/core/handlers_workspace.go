// Odpowiedzialność pliku: wpięcie sześciu komend obszaru `workspace.*` — modułu
// Workspace wraz z jego pięcioma oknami operacyjnymi (Project Dashboard,
// Instructions Panel, Context Memory, Project Library, Agent Manager).
//
// `workspace.enter` nie należy do tego obszaru. Ta komenda przeładowuje
// przestrzeń roboczą karty sesji na dowolny moduł i wpina ją nawigacja
// (handlers_nawigacja.go) — nazwa jest wspólna, obszar nie.
//
// Jedno zdarzenie na cały moduł. Kontrakt daje modułowi wyłącznie
// `workspace.project.changed`, więc każda zmiana stanu projektu — instrukcje,
// wpis pamięci, przypisanie eksperta — rozgłasza się projektem po zmianie.
// Okna modułu odświeżają się z jednej subskrypcji, a nie z czterech. Odczyty
// (`dashboard.get`, `context.get`, `library.list`) niczego nie rozgłaszają poza
// jednym przypadkiem: wejście na pulpit projektu, którego jeszcze nie było,
// zakłada go — a założenie bytu jest zmianą (ChangeKind `created`).
package core

import (
	"context"

	"danacoconsole/shared"
)

// PrzestrzenRobocza jest portem modułu Workspace.
type PrzestrzenRobocza interface {
	Pulpit(ctx context.Context, z shared.WorkspaceDashboardGetRequest) (shared.WorkspaceDashboardGetResponse, error)
	ZapiszInstrukcje(ctx context.Context, z shared.WorkspaceInstructionsSetRequest) (shared.WorkspaceInstructionsSetResponse, error)
	ZapiszWpisPamieci(ctx context.Context, z shared.WorkspaceContextSetRequest) (shared.WorkspaceContextSetResponse, error)
	WpisyPamieci(ctx context.Context, z shared.WorkspaceContextGetRequest) (shared.WorkspaceContextGetResponse, error)
	Biblioteka(ctx context.Context, z shared.WorkspaceLibraryListRequest) (shared.WorkspaceLibraryListResponse, error)
	PrzypiszAgenta(ctx context.Context, z shared.WorkspaceAgentAssignRequest) (shared.WorkspaceAgentAssignResponse, error)
	// Projekt oddaje projekt po zmianie. Służy rozgłoszeniu zdarzenia po
	// komendach, których wynik projektu nie niesie.
	Projekt(ctx context.Context, idProjektu string) (shared.WorkspaceProject, error)
}

// zarejestrujPrzestrzenRobocza wpina sześć komend modułu Workspace.
func zarejestrujPrzestrzenRobocza(r *Rejestr, w PrzestrzenRobocza, e *emiter) {
	if r == nil || w == nil {
		return
	}

	r.Zarejestruj(shared.CommandWorkspaceDashboardGet,
		obsluz(func(ctx context.Context, z shared.WorkspaceDashboardGetRequest) (shared.WorkspaceDashboardGetResponse, error) {
			odpowiedz, err := w.Pulpit(ctx, z)
			if err == nil {
				e.projekt(rodzajZmianyProjektu(odpowiedz.Dashboard.Project), odpowiedz.Dashboard.Project)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceInstructionsSet,
		obsluz(func(ctx context.Context, z shared.WorkspaceInstructionsSetRequest) (shared.WorkspaceInstructionsSetResponse, error) {
			odpowiedz, err := w.ZapiszInstrukcje(ctx, z)
			if err == nil {
				rozglosProjekt(ctx, w, e, z.ProjectId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceContextSet,
		obsluz(func(ctx context.Context, z shared.WorkspaceContextSetRequest) (shared.WorkspaceContextSetResponse, error) {
			odpowiedz, err := w.ZapiszWpisPamieci(ctx, z)
			if err == nil {
				rozglosProjekt(ctx, w, e, z.ProjectId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceContextGet, obsluz(w.WpisyPamieci))

	r.Zarejestruj(shared.CommandWorkspaceLibraryList, obsluz(w.Biblioteka))

	r.Zarejestruj(shared.CommandWorkspaceAgentAssign,
		obsluz(func(ctx context.Context, z shared.WorkspaceAgentAssignRequest) (shared.WorkspaceAgentAssignResponse, error) {
			odpowiedz, err := w.PrzypiszAgenta(ctx, z)
			if err == nil {
				rozglosProjekt(ctx, w, e, z.ProjectId)
			}
			return odpowiedz, err
		}))
}

// rodzajZmianyProjektu odróżnia projekt założony właśnie teraz od zastanego.
// Rozróżnienie bierze się z samego bytu, nie z pamięci adaptera: projekt
// dopiero założony ma czas utworzenia równy czasowi ostatniej zmiany, a każda
// późniejsza czynność ten drugi przesuwa.
func rodzajZmianyProjektu(p shared.WorkspaceProject) shared.ChangeKind {
	if p.CreatedAt == p.UpdatedAt {
		return shared.ChangeKindCreated
	}
	return shared.ChangeKindUpdated
}

// rozglosProjekt dobiera projekt po zmianie i rozgłasza ją. Nieudany dobór
// kończy wyłącznie rozgłoszenie — komenda już się powiodła.
func rozglosProjekt(ctx context.Context, w PrzestrzenRobocza, e *emiter, idProjektu string) {
	projekt, err := w.Projekt(ctx, idProjektu)
	if err != nil {
		return
	}
	e.projekt(shared.ChangeKindUpdated, projekt)
}

// projekt rozgłasza zmianę projektu przestrzeni roboczej. Projekt jest
// komponentem własnym, nie bytem jednej karty sesji, więc zdarzenie idzie bez
// jej wskazania — pracuje nad nim wiele kart naraz.
func (e *emiter) projekt(zmiana shared.ChangeKind, p shared.WorkspaceProject) {
	e.wyslij(shared.EventWorkspaceProjectChanged, "",
		shared.WorkspaceProjectChangedEvent{Change: zmiana, Project: p})
}
