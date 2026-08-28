package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujKolejki wpina domenę kolejek jednego silnika pętli koordynator–wykonawca, wspólnego
// dla pętli sesyjnej i MultitaskingAI, żeby druga implementacja nie miała gdzie powstać.
func zarejestrujKolejki(r *Rejestr, kolejki Kolejki, e *emiter) {
	if r == nil || kolejki == nil {
		return
	}

	r.Zarejestruj(shared.CommandQueueCreate,
		obsluz(func(ctx context.Context, z shared.QueueCreateRequest) (shared.QueueCreateResponse, error) {
			w, err := kolejki.Utworz(ctx, z)
			if err == nil {
				e.kolejka(shared.ChangeKindCreated, w.Queue)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandQueueAction,
		obsluz(func(ctx context.Context, z shared.QueueActionRequest) (shared.QueueActionResponse, error) {
			w, err := kolejki.Wykonaj(ctx, z)
			if err == nil {
				e.kolejka(shared.ChangeKindUpdated, w.Queue)
			}
			return w, err
		}))
}
