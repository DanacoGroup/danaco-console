// Wpięcie dwóch komend rodziny `queue.*` — `queue.list` i `queue.link`.
//
// Osobno od `handlers_queue.go`, który wpina `queue.create` i `queue.action`
// przez port `Kolejki`. Port poniżej jest rozszerzeniem portu `Kolejki`, nie
// drugim portem: silnik kolejek pętli sesyjnej i MultitaskingAI jest jeden.
//
// Zdarzenie rozgłasza wyłącznie `queue.link`. Wiązanie zmienia kolejkę, więc
// idzie tym samym `queue.changed`, co założenie i działanie — rodzajem zmiany
// `updated`. `queue.list` niczego nie zmienia i niczego nie rozgłasza; wykaz
// rozgłoszony jako zmiana byłby zdarzeniem bez faktu.
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

// zarejestrujWiazaniaKolejek wpina `queue.list` i `queue.link`.
//
// Brak portu kolejek zostawia obie komendy nieznane — tak samo jak każdą inną
// domenę bez portu. Port kolejek bez tych dwóch czynności to co innego:
// kolejki są, a rdzeń nie umie ich oddać. Wtedy komendy zostają wpięte
// i odmawiają wprost kodem `internal_error`, bo odpowiedź „nieznana komenda"
// wskazywałaby na brak kolejek, a nie na usterkę montażu.
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

// zarejestrujOdmoweWiazanKolejek wpina obie komendy jako odmowę montażu.
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
