// Odpowiedzialność pliku: wpięcie komend huba planowania modułu Workspace —
// zadania, tablica kanban, oś czasu i kalendarz.
//
// Port jest osobny od `PrzestrzenRobocza`, choć wypełnia go ten sam adapter.
// Powód jest ten sam, co przy pamięci projektu: rodzina liczy trzynaście komend
// i wpisanie ich do portu wspólnego zrobiłoby z niego wykaz wszystkiego, co
// moduł umie, zamiast wykazu tego, czym jest projekt.
//
// Zdarzeniem modułu jest jedno `workspace.project.changed` — każda zmiana
// planu rozgłasza się projektem. Okna huba odświeżają się z tej jednej
// subskrypcji, a nie z trzynastu.
package core

import (
	"context"

	"danacoconsole/shared"
)

// PlanowanieProjektu jest portem huba planowania.
type PlanowanieProjektu interface {
	ZalozZadanie(ctx context.Context, z shared.WorkspaceTaskCreateRequest) (shared.WorkspaceTaskCreateResponse, error)
	ZmienZadanie(ctx context.Context, z shared.WorkspaceTaskUpdateRequest) (shared.WorkspaceTaskUpdateResponse, error)
	UsunZadanie(ctx context.Context, z shared.WorkspaceTaskDeleteRequest) (shared.WorkspaceTaskDeleteResponse, error)
	Zadania(ctx context.Context, z shared.WorkspaceTaskListRequest) (shared.WorkspaceTaskListResponse, error)
	PrzeniesZadanie(ctx context.Context, z shared.WorkspaceTaskMoveRequest) (shared.WorkspaceTaskMoveResponse, error)
	ZalozZaleznosc(ctx context.Context, z shared.WorkspaceTaskDependencySetRequest) (shared.WorkspaceTaskDependencySetResponse, error)
	ZniesZaleznosc(ctx context.Context, z shared.WorkspaceTaskDependencyRemoveRequest) (shared.WorkspaceTaskDependencyRemoveResponse, error)
	Tablica(ctx context.Context, z shared.WorkspaceBoardGetRequest) (shared.WorkspaceBoardGetResponse, error)
	Harmonogram(ctx context.Context, z shared.WorkspaceScheduleGetRequest) (shared.WorkspaceScheduleGetResponse, error)
	Kalendarz(ctx context.Context, z shared.WorkspaceCalendarGetRequest) (shared.WorkspaceCalendarGetResponse, error)
	WciagnijKalendarz(ctx context.Context, z shared.WorkspaceCalendarImportRequest) (shared.WorkspaceCalendarImportResponse, error)
	ZapiszKalendarz(ctx context.Context, z shared.WorkspaceCalendarExportRequest) (shared.WorkspaceCalendarExportResponse, error)
	// Projekt służy rozgłoszeniu zmiany po komendach, których wynik projektu
	// nie niesie.
	Projekt(ctx context.Context, idProjektu string) (shared.WorkspaceProject, error)
}

// zarejestrujPlanowanieProjektu wpina dwanaście komend huba planowania.
func zarejestrujPlanowanieProjektu(r *Rejestr, p PlanowanieProjektu, e *emiter) {
	if r == nil || p == nil {
		return
	}

	r.Zarejestruj(shared.CommandWorkspaceTaskCreate,
		obsluz(func(ctx context.Context, z shared.WorkspaceTaskCreateRequest) (shared.WorkspaceTaskCreateResponse, error) {
			odpowiedz, err := p.ZalozZadanie(ctx, z)
			if err == nil {
				rozglosProjektPlanowania(ctx, p, e, z.ProjectId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceTaskUpdate,
		obsluz(func(ctx context.Context, z shared.WorkspaceTaskUpdateRequest) (shared.WorkspaceTaskUpdateResponse, error) {
			odpowiedz, err := p.ZmienZadanie(ctx, z)
			if err == nil {
				rozglosProjektPlanowania(ctx, p, e, odpowiedz.Task.ProjectId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceTaskDelete,
		obsluz(func(ctx context.Context, z shared.WorkspaceTaskDeleteRequest) (shared.WorkspaceTaskDeleteResponse, error) {
			return p.UsunZadanie(ctx, z)
		}))

	r.Zarejestruj(shared.CommandWorkspaceTaskList, obsluz(p.Zadania))

	r.Zarejestruj(shared.CommandWorkspaceTaskMove,
		obsluz(func(ctx context.Context, z shared.WorkspaceTaskMoveRequest) (shared.WorkspaceTaskMoveResponse, error) {
			odpowiedz, err := p.PrzeniesZadanie(ctx, z)
			if err == nil {
				rozglosProjektPlanowania(ctx, p, e, odpowiedz.Task.ProjectId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceTaskDependencySet,
		obsluz(func(ctx context.Context, z shared.WorkspaceTaskDependencySetRequest) (shared.WorkspaceTaskDependencySetResponse, error) {
			odpowiedz, err := p.ZalozZaleznosc(ctx, z)
			if err == nil {
				rozglosProjektPlanowania(ctx, p, e, z.ProjectId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceTaskDependencyRemove, obsluz(p.ZniesZaleznosc))
	r.Zarejestruj(shared.CommandWorkspaceBoardGet, obsluz(p.Tablica))
	r.Zarejestruj(shared.CommandWorkspaceScheduleGet, obsluz(p.Harmonogram))
	r.Zarejestruj(shared.CommandWorkspaceCalendarGet, obsluz(p.Kalendarz))

	r.Zarejestruj(shared.CommandWorkspaceCalendarImport,
		obsluz(func(ctx context.Context, z shared.WorkspaceCalendarImportRequest) (shared.WorkspaceCalendarImportResponse, error) {
			odpowiedz, err := p.WciagnijKalendarz(ctx, z)
			if err == nil && odpowiedz.Imported > 0 {
				rozglosProjektPlanowania(ctx, p, e, z.ProjectId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceCalendarExport, obsluz(p.ZapiszKalendarz))
}

// rozglosProjektPlanowania dobiera projekt po zmianie i rozgłasza ją. Nieudany
// dobór kończy wyłącznie rozgłoszenie — komenda już się powiodła.
func rozglosProjektPlanowania(ctx context.Context, p PlanowanieProjektu, e *emiter, idProjektu string) {
	if idProjektu == "" {
		return
	}
	projekt, err := p.Projekt(ctx, idProjektu)
	if err != nil {
		return
	}
	e.projekt(shared.ChangeKindUpdated, projekt)
}
