// Plik wpina dwie komendy rodziny queue.* — wykaz i wiązanie kolejek — jako rozszerzenie portu
// Kolejki, bo silnik kolejek pętli sesyjnej i MultitaskingAI jest jeden.
package core

import (
	"context"

	"danacoconsole/shared"
)

// wiazaneKolejki jest portem dwóch komend rodziny `queue.*` — portem
// kolejek rozszerzonym o wykaz i wiązanie.
type wiazaneKolejki interface {
	Kolejki

	// Wykaz obsługuje `queue.list`.
	Wykaz(ctx context.Context, z shared.QueueListRequest) (shared.QueueListResponse, error)
	// Zwiaz obsługuje `queue.link`.
	Zwiaz(ctx context.Context, z shared.QueueLinkRequest) (shared.QueueLinkResponse, error)
}

// zarejestrujWiazaniaKolejek wpina komendy wykazu i wiązania kolejek. Brak portu kolejek zostawia
// obie komendy nieznane, tak samo jak każdą inną domenę bez portu.
func zarejestrujWiazaniaKolejek(r *Rejestr, kolejki Kolejki, e *emiter) {
	if r == nil || kolejki == nil {
		return
	}
	wiazane, ok := kolejki.(wiazaneKolejki)
	if !ok {
		zarejestrujOdmoweWiazanKolejek(r)
		return
	}

	r.Zarejestruj(shared.CommandQueueList, obsluz(wiazane.Wykaz))

	r.Zarejestruj(shared.CommandQueueLink,
		obsluz(func(ctx context.Context, z shared.QueueLinkRequest) (shared.QueueLinkResponse, error) {
			w, err := wiazane.Zwiaz(ctx, z)
			if err == nil {
				e.kolejka(shared.ChangeKindUpdated, w.Queue)
			}
			return w, err
		}))
}

// zarejestrujOdmoweWiazanKolejek wpina obie komendy rodziny queue.* jako odmowę usterki montażu rdzenia.
func zarejestrujOdmoweWiazanKolejek(r *Rejestr) {
	const powod = "kolejki: port kolejek nie niesie wykazu ani wiązania kolejki"

	r.Zarejestruj(shared.CommandQueueList,
		obsluz(func(context.Context, shared.QueueListRequest) (shared.QueueListResponse, error) {
			return shared.QueueListResponse{}, bladMontazuKolejek(powod)
		}))
	r.Zarejestruj(shared.CommandQueueLink,
		obsluz(func(context.Context, shared.QueueLinkRequest) (shared.QueueLinkResponse, error) {
			return shared.QueueLinkResponse{}, bladMontazuKolejek(powod)
		}))
}
