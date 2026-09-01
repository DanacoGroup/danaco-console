// Plik wpina komendy wiedzy i współpracy modułu Workspace: notatki i wiki, graf wiedzy, tablicę
// wizualną, materiały biblioteki, wyszukiwanie, oś czasu i komentarze, oraz trzy czynności na
// projekcie: stan, odłączenie eksperta i historię instrukcji.
package core

import (
	"context"

	"danacoconsole/shared"
)

// WiedzaProjektu jest portem obszaru wiedzy i współpracy projektu: notatek, grafu i biblioteki materiałów.
type WiedzaProjektu interface {
	ZapiszNotatke(ctx context.Context, z shared.WorkspaceNoteSaveRequest) (shared.WorkspaceNoteSaveResponse, error)
	Notatka(ctx context.Context, z shared.WorkspaceNoteGetRequest) (shared.WorkspaceNoteGetResponse, error)
	Notatki(ctx context.Context, z shared.WorkspaceNoteListRequest) (shared.WorkspaceNoteListResponse, error)
	UsunNotatke(ctx context.Context, z shared.WorkspaceNoteDeleteRequest) (shared.WorkspaceNoteDeleteResponse, error)
	DrzewoNotatek(ctx context.Context, z shared.WorkspaceNoteTreeGetRequest) (shared.WorkspaceNoteTreeGetResponse, error)
	OdnosnikiWsteczne(ctx context.Context, z shared.WorkspaceNoteBacklinkListRequest) (shared.WorkspaceNoteBacklinkListResponse, error)
	GrafWiedzy(ctx context.Context, z shared.WorkspaceKnowledgeGraphGetRequest) (shared.WorkspaceKnowledgeGraphGetResponse, error)
	Kanwa(ctx context.Context, z shared.WorkspaceCanvasGetRequest) (shared.WorkspaceCanvasGetResponse, error)
	ZapiszKanwe(ctx context.Context, z shared.WorkspaceCanvasSaveRequest) (shared.WorkspaceCanvasSaveResponse, error)
	WydobadzTekst(ctx context.Context, z shared.WorkspaceLibraryTextExtractRequest) (shared.WorkspaceLibraryTextExtractResponse, error)
	Duplikaty(ctx context.Context, z shared.WorkspaceLibraryDuplicateListRequest) (shared.WorkspaceLibraryDuplicateListResponse, error)
	ScalDuplikaty(ctx context.Context, z shared.WorkspaceLibraryDuplicateMergeRequest) (shared.WorkspaceLibraryDuplicateMergeResponse, error)
	SzukajWProjekcie(ctx context.Context, z shared.WorkspaceSearchProjectRequest) (shared.WorkspaceSearchProjectResponse, error)
	OsCzasu(ctx context.Context, z shared.WorkspaceActivityListRequest) (shared.WorkspaceActivityListResponse, error)
	DolozKomentarz(ctx context.Context, z shared.WorkspaceCommentAddRequest) (shared.WorkspaceCommentAddResponse, error)
	Komentarze(ctx context.Context, z shared.WorkspaceCommentListRequest) (shared.WorkspaceCommentListResponse, error)
	UsunKomentarz(ctx context.Context, z shared.WorkspaceCommentDeleteRequest) (shared.WorkspaceCommentDeleteResponse, error)
	UstawStanProjektu(ctx context.Context, z shared.WorkspaceProjectStatusSetRequest) (shared.WorkspaceProjectStatusSetResponse, error)
	OdlaczAgenta(ctx context.Context, z shared.WorkspaceAgentUnassignRequest) (shared.WorkspaceAgentUnassignResponse, error)
	WersjeInstrukcji(ctx context.Context, z shared.WorkspaceInstructionsVersionListRequest) (shared.WorkspaceInstructionsVersionListResponse, error)
	PrzywrocInstrukcje(ctx context.Context, z shared.WorkspaceInstructionsVersionRestoreRequest) (shared.WorkspaceInstructionsVersionRestoreResponse, error)
	Projekt(ctx context.Context, idProjektu string) (shared.WorkspaceProject, error)
}

// zarejestrujWiedzeProjektu wpina dwadzieścia jeden komend obszaru wiedzy w rejestrze rdzenia platformy.
func zarejestrujWiedzeProjektu(r *Rejestr, w WiedzaProjektu, e *emiter) {
	if r == nil || w == nil {
		return
	}

	r.Zarejestruj(shared.CommandWorkspaceNoteSave,
		obsluz(func(ctx context.Context, z shared.WorkspaceNoteSaveRequest) (shared.WorkspaceNoteSaveResponse, error) {
			odpowiedz, err := w.ZapiszNotatke(ctx, z)
			if err == nil {
				rozglosProjektWiedzy(ctx, w, e, z.ProjectId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceNoteGet, obsluz(w.Notatka))
	r.Zarejestruj(shared.CommandWorkspaceNoteList, obsluz(w.Notatki))
	r.Zarejestruj(shared.CommandWorkspaceNoteDelete, obsluz(w.UsunNotatke))
	r.Zarejestruj(shared.CommandWorkspaceNoteTreeGet, obsluz(w.DrzewoNotatek))
	r.Zarejestruj(shared.CommandWorkspaceNoteBacklinkList, obsluz(w.OdnosnikiWsteczne))
	r.Zarejestruj(shared.CommandWorkspaceKnowledgeGraphGet, obsluz(w.GrafWiedzy))
	r.Zarejestruj(shared.CommandWorkspaceCanvasGet, obsluz(w.Kanwa))

	r.Zarejestruj(shared.CommandWorkspaceCanvasSave,
		obsluz(func(ctx context.Context, z shared.WorkspaceCanvasSaveRequest) (shared.WorkspaceCanvasSaveResponse, error) {
			odpowiedz, err := w.ZapiszKanwe(ctx, z)
			if err == nil {
				rozglosProjektWiedzy(ctx, w, e, z.ProjectId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceLibraryTextExtract, obsluz(w.WydobadzTekst))
	r.Zarejestruj(shared.CommandWorkspaceLibraryDuplicateList, obsluz(w.Duplikaty))

	r.Zarejestruj(shared.CommandWorkspaceLibraryDuplicateMerge,
		obsluz(func(ctx context.Context, z shared.WorkspaceLibraryDuplicateMergeRequest) (shared.WorkspaceLibraryDuplicateMergeResponse, error) {
			odpowiedz, err := w.ScalDuplikaty(ctx, z)
			if err == nil {
				rozglosProjektWiedzy(ctx, w, e, z.ProjectId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceSearchProject, obsluz(w.SzukajWProjekcie))
	r.Zarejestruj(shared.CommandWorkspaceActivityList, obsluz(w.OsCzasu))

	r.Zarejestruj(shared.CommandWorkspaceCommentAdd,
		obsluz(func(ctx context.Context, z shared.WorkspaceCommentAddRequest) (shared.WorkspaceCommentAddResponse, error) {
			odpowiedz, err := w.DolozKomentarz(ctx, z)
			if err == nil {
				rozglosProjektWiedzy(ctx, w, e, z.ProjectId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceCommentList, obsluz(w.Komentarze))
	r.Zarejestruj(shared.CommandWorkspaceCommentDelete, obsluz(w.UsunKomentarz))

	r.Zarejestruj(shared.CommandWorkspaceProjectStatusSet,
		obsluz(func(ctx context.Context, z shared.WorkspaceProjectStatusSetRequest) (shared.WorkspaceProjectStatusSetResponse, error) {
			odpowiedz, err := w.UstawStanProjektu(ctx, z)
			if err == nil {
				e.projekt(ctx, shared.ChangeKindUpdated, odpowiedz.Project)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceAgentUnassign,
		obsluz(func(ctx context.Context, z shared.WorkspaceAgentUnassignRequest) (shared.WorkspaceAgentUnassignResponse, error) {
			odpowiedz, err := w.OdlaczAgenta(ctx, z)
			if err == nil && odpowiedz.Unassigned {
				rozglosProjektWiedzy(ctx, w, e, z.ProjectId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandWorkspaceInstructionsVersionList, obsluz(w.WersjeInstrukcji))

	r.Zarejestruj(shared.CommandWorkspaceInstructionsVersionRestore,
		obsluz(func(ctx context.Context, z shared.WorkspaceInstructionsVersionRestoreRequest) (shared.WorkspaceInstructionsVersionRestoreResponse, error) {
			odpowiedz, err := w.PrzywrocInstrukcje(ctx, z)
			if err == nil {
				rozglosProjektWiedzy(ctx, w, e, z.ProjectId)
			}
			return odpowiedz, err
		}))
}

// rozglosProjektWiedzy dobiera projekt po zmianie wiedzy i rozgłasza ją oknom biblioteki tego projektu.
func rozglosProjektWiedzy(ctx context.Context, w WiedzaProjektu, e *emiter, idProjektu string) {
	if idProjektu == "" {
		return
	}
	projekt, err := w.Projekt(ctx, idProjektu)
	if err != nil {
		return
	}
	e.projekt(ctx, shared.ChangeKindUpdated, projekt)
}
