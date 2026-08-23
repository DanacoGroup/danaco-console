package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujOkna wpina domenę okna komunikacji — bytu pośredniego między sesją
// a wiadomością.
//
// Okno niesie moduł, kanał modelu, listę katalogów roboczych, zasięg wykonania,
// tryb uprawnień i rolę w pętli koordynator–wykonawca. Wiele okien jednej sesji
// biegnie równolegle, każde z własnym modelem i własnym katalogiem, dlatego
// każda komenda wskazuje okno wprost, a nie przez sesję.
func zarejestrujOkna(r *Rejestr, okna Okna, e *emiter) {
	if r == nil || okna == nil {
		return
	}

	r.Zarejestruj(shared.CommandWindowCreate,
		obsluz(func(ctx context.Context, z shared.WindowCreateRequest) (shared.WindowCreateResponse, error) {
			w, err := okna.Utworz(ctx, z)
			if err == nil {
				e.okno(ctx, shared.ChangeKindCreated, w.Window)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandWindowList, obsluz(okna.Wykaz))

	r.Zarejestruj(shared.CommandWindowUpdate,
		obsluz(func(ctx context.Context, z shared.WindowUpdateRequest) (shared.WindowUpdateResponse, error) {
			w, err := okna.Zmien(ctx, z)
			if err == nil {
				e.okno(ctx, shared.ChangeKindUpdated, w.Window)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandWindowClose,
		obsluz(func(ctx context.Context, z shared.WindowCloseRequest) (shared.WindowCloseResponse, error) {
			w, err := okna.Zamknij(ctx, z)
			if err == nil {
				e.okno(ctx, shared.ChangeKindUpdated, w.Window)
			}
			return w, err
		}))
}
