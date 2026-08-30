// Plik wpina sześć komend obszaru workspace.* obsługujące moduł Workspace wraz z jego pięcioma
// oknami operacyjnymi: pulpitem projektu, panelem instrukcji, pamięcią kontekstu, biblioteką
// i menedżerem agentów.
package core

import (
	"context"

	"danacoconsole/shared"
)

// PrzestrzenRobocza jest portem modułu Workspace obsługującym pulpit, instrukcje, pamięć i bibliotekę.
type PrzestrzenRobocza interface {
	Pulpit(ctx context.Context, z shared.WorkspaceDashboardGetRequest) (shared.WorkspaceDashboardGetResponse, error)
	ZapiszInstrukcje(ctx context.Context, z shared.WorkspaceInstructionsSetRequest) (shared.WorkspaceInstructionsSetResponse, error)
	ZapiszWpisPamieci(ctx context.Context, z shared.WorkspaceContextSetRequest) (shared.WorkspaceContextSetResponse, error)
	WpisyPamieci(ctx context.Context, z shared.WorkspaceContextGetRequest) (shared.WorkspaceContextGetResponse, error)
	Biblioteka(ctx context.Context, z shared.WorkspaceLibraryListRequest) (shared.WorkspaceLibraryListResponse, error)
	PrzypiszAgenta(ctx context.Context, z shared.WorkspaceAgentAssignRequest) (shared.WorkspaceAgentAssignResponse, error)
	// Projekt oddaje projekt po zmianie, dla komend, których wynik projektu nie niesie.
	Projekt(ctx context.Context, idProjektu string) (shared.WorkspaceProject, error)
	// Projekty oddaje wykaz projektów konta dla lewego panelu ramy.
	Projekty(ctx context.Context, z shared.ProjectListRequest) (shared.ProjectListResponse, error)
}

// zarejestrujPrzestrzenRobocza wpina sześć komend modułu Workspace w rejestrze rdzenia tej platformy konta.
func zarejestrujPrzestrzenRobocza(r *Rejestr, w PrzestrzenRobocza, e *emiter) {
	if r == nil || w == nil {
		return
	}

	r.Zarejestruj(shared.CommandProjectList, obsluz(w.Projekty))

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

// rodzajZmianyProjektu odróżnia projekt założony właśnie teraz od zastanego. Rozróżnienie bierze
// się z bytu, nie z pamięci adaptera: nowy projekt ma czas utworzenia równy czasowi ostatniej zmiany.
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
