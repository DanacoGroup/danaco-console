package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujKolejki wpina domenę kolejek jednego silnika pętli
// koordynator–wykonawca. Ten sam silnik obsługuje pętlę sesyjną
// i MultitaskingAI — rdzeń kieruje komendy kolejek w jedno miejsce, żeby druga
// implementacja nie miała gdzie powstać.
//
// Działanie „powtórz" nie ma limitu obiegów. Rdzeń nie zlicza prób
// i nie odmawia po którejś z kolei — przerwanie należy do użytkownika.
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
