package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujSesje wpina domenę sesji: zakładanie, wykaz, otwarcie z oknami, zamknięcie i
// usunięcie, jako byt wspólny dla plików, pamięci, projektu i agentów.
func zarejestrujSesje(r *Rejestr, sesje Sesje, e *emiter) {
	if r == nil || sesje == nil {
		return
	}

	r.Zarejestruj(shared.CommandSessionCreate,
		obsluz(func(ctx context.Context, z shared.SessionCreateRequest) (shared.SessionCreateResponse, error) {
			w, err := sesje.Utworz(ctx, z)
			if err == nil {
				e.sesja(ctx, shared.ChangeKindCreated, w.Session)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandSessionList, obsluz(sesje.Wykaz))

	r.Zarejestruj(shared.CommandSessionOpen, obsluz(sesje.Otworz))

	r.Zarejestruj(shared.CommandSessionClose,
		obsluz(func(ctx context.Context, z shared.SessionCloseRequest) (shared.SessionCloseResponse, error) {
			w, err := sesje.Zamknij(ctx, z)
			if err == nil {
				e.sesja(ctx, shared.ChangeKindUpdated, w.Session)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandSessionDelete,
		obsluz(func(ctx context.Context, z shared.SessionDeleteRequest) (shared.SessionDeleteResponse, error) {
			w, err := sesje.Usun(ctx, z)
			if err == nil {
				// Zdarzenie niesie wyłącznie identyfikator, osobno za każdą usuniętą sesję.
				for _, identyfikator := range w.DeletedIds {
					e.sesja(ctx, shared.ChangeKindDeleted, shared.Session{Id: identyfikator})
				}
			}
			return w, err
		}))

	zarejestrujHistorieSesji(r, sesje, e)
}

// zarejestrujHistorieSesji wpina czynności Operatora na wykazie sesji: nazwę, projekt, archiwum,
// kopię, wznowienie i zatrzymanie, osobno od bieżącej pracy nad sesją.
func zarejestrujHistorieSesji(r *Rejestr, sesje Sesje, e *emiter) {
	r.Zarejestruj(shared.CommandSessionRename,
		obsluz(func(ctx context.Context, z shared.SessionRenameRequest) (shared.SessionRenameResponse, error) {
			w, err := sesje.ZmienNazwe(ctx, z)
			if err == nil {
				e.sesja(ctx, shared.ChangeKindUpdated, w.Session)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandSessionArchive,
		obsluz(func(ctx context.Context, z shared.SessionArchiveRequest) (shared.SessionArchiveResponse, error) {
			w, err := sesje.Archiwizuj(ctx, z)
			if err == nil {
				for _, identyfikator := range w.ArchivedIds {
					e.sesja(ctx, shared.ChangeKindUpdated, shared.Session{
						Id: identyfikator, Status: shared.SessionStatusArchived,
					})
				}
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandSessionRestore,
		obsluz(func(ctx context.Context, z shared.SessionRestoreRequest) (shared.SessionRestoreResponse, error) {
			w, err := sesje.Przywroc(ctx, z)
			if err == nil {
				for _, identyfikator := range w.RestoredIds {
					e.sesja(ctx, shared.ChangeKindUpdated, shared.Session{
						Id: identyfikator, Status: shared.SessionStatusActive,
					})
				}
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandSessionArchiveList,
		obsluz(sesje.WykazArchiwum))

	r.Zarejestruj(shared.CommandSessionProjectSet,
		obsluz(sesje.PrzypiszProjekt))

	r.Zarejestruj(shared.CommandSessionProjectClear,
		obsluz(sesje.OdepnijProjekt))

	r.Zarejestruj(shared.CommandSessionResume,
		obsluz(func(ctx context.Context, z shared.SessionResumeRequest) (shared.SessionResumeResponse, error) {
			w, err := sesje.Wznow(ctx, z)
			if err == nil {
				e.sesja(ctx, shared.ChangeKindUpdated, w.Session)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandSessionStop,
		obsluz(sesje.Zatrzymaj))

	r.Zarejestruj(shared.CommandSessionCopy,
		obsluz(func(ctx context.Context, z shared.SessionCopyRequest) (shared.SessionCopyResponse, error) {
			w, err := sesje.Kopiuj(ctx, z)
			if err == nil {
				e.sesja(ctx, shared.ChangeKindCreated, w.Session)
			}
			return w, err
		}))
}
