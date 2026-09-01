// Plik wpina dwanaście komend rodziny queue.* dotyczących zleceń kolejki, jej polityki, zadań martwych i głębokości; rozgłoszenie idzie queue.changed rodzajem updated, nie created.
package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// zlecenioweKolejki jest portem dwunastu komend rodziny `queue.*` — portem
// kolejek rozszerzonym o zlecenia, politykę, zadania martwe i głębokość.
type zlecenioweKolejki interface {
	Kolejki

	// Kolejka oddaje kolejkę po identyfikatorze, nośnik odpowiedzi czynności, których wynik to kolejka.
	Kolejka(ctx context.Context, id string) (shared.Queue, error)

	DodajZlecenie(ctx context.Context, z shared.QueueItemEnqueueRequest) (shared.QueueItemEnqueueResponse, error)
	ZdejmijZlecenie(ctx context.Context, z shared.QueueItemDequeueRequest) (shared.QueueItemDequeueResponse, error)
	OdlozZlecenie(ctx context.Context, z shared.QueueItemDelayRequest) (shared.QueueItemDelayResponse, error)
	PodzielZlecenie(ctx context.Context, z shared.QueueItemSplitRequest) (shared.QueueItemSplitResponse, error)
	ScalZlecenia(ctx context.Context, z shared.QueueItemMergeRequest) (shared.QueueItemMergeResponse, error)
	SkierujZlecenie(ctx context.Context, z shared.QueueItemRouteRequest) (shared.QueueItemRouteResponse, error)
	RozgalezZlecenie(ctx context.Context, z shared.QueueItemBranchRequest) (shared.QueueItemBranchResponse, error)
	UwarunkujZlecenie(ctx context.Context, z shared.QueueItemConditionRequest) (shared.QueueItemConditionResponse, error)
	WykazZlecen(ctx context.Context, z shared.QueueItemListRequest) (shared.QueueItemListResponse, error)
	UstawPolitykeKolejki(ctx context.Context, z shared.QueuePolicySetRequest) (shared.QueuePolicySetResponse, error)
	WykazZadanMartwych(ctx context.Context, z shared.QueueDeadListRequest) (shared.QueueDeadListResponse, error)
	GlebokoscKolejki(ctx context.Context, z shared.QueueDepthGetRequest) (shared.QueueDepthGetResponse, error)
}

// zarejestrujZleceniaKolejek wpina dwanaście komend zleceń kolejki do rejestru komend tego rdzenia całego.
func zarejestrujZleceniaKolejek(r *Rejestr, kolejki Kolejki, e *emiter) {
	if r == nil || kolejki == nil {
		return
	}
	zleceniowe, ok := kolejki.(zlecenioweKolejki)
	if !ok {
		zarejestrujOdmoweZlecenKolejek(r)
		return
	}

	r.Zarejestruj(shared.CommandQueueItemList, obsluz(zleceniowe.WykazZlecen))
	r.Zarejestruj(shared.CommandQueueDeadList, obsluz(zleceniowe.WykazZadanMartwych))
	r.Zarejestruj(shared.CommandQueueDepthGet, obsluz(zleceniowe.GlebokoscKolejki))

	r.Zarejestruj(shared.CommandQueuePolicySet,
		obsluz(func(ctx context.Context, z shared.QueuePolicySetRequest) (shared.QueuePolicySetResponse, error) {
			w, err := zleceniowe.UstawPolitykeKolejki(ctx, z)
			if err == nil {
				e.kolejka(ctx, shared.ChangeKindUpdated, w.Queue)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandQueueItemEnqueue,
		obsluz(func(ctx context.Context, z shared.QueueItemEnqueueRequest) (shared.QueueItemEnqueueResponse, error) {
			w, err := zleceniowe.DodajZlecenie(ctx, z)
			// Duplikat niczego nie zmienił: zdarzenie odbite idempotencją nie ma faktu, więc się nie rozgłasza.
			if err == nil && !w.Duplicate {
				rozglosKolejkeZlecenia(ctx, zleceniowe, e, z.QueueId)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandQueueItemDequeue,
		obsluz(func(ctx context.Context, z shared.QueueItemDequeueRequest) (shared.QueueItemDequeueResponse, error) {
			w, err := zleceniowe.ZdejmijZlecenie(ctx, z)
			if err == nil && w.Removed {
				e.kolejka(ctx, shared.ChangeKindUpdated, w.Queue)
			}
			return w, err
		}))

	zarejestrujCzynnosciZlecen(r, zleceniowe, e)
}

// zarejestrujCzynnosciZlecen wpina sześć czynności, których wynikiem jest
// zlecenie (albo ich wykaz), a nie kolejka — rozgłoszenie dobiera więc kolejkę
// osobno, po identyfikatorze z żądania.
func zarejestrujCzynnosciZlecen(r *Rejestr, m zlecenioweKolejki, e *emiter) {
	r.Zarejestruj(shared.CommandQueueItemDelay,
		obsluz(func(ctx context.Context, z shared.QueueItemDelayRequest) (shared.QueueItemDelayResponse, error) {
			w, err := m.OdlozZlecenie(ctx, z)
			if err == nil {
				rozglosKolejkeZlecenia(ctx, m, e, z.QueueId)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandQueueItemSplit,
		obsluz(func(ctx context.Context, z shared.QueueItemSplitRequest) (shared.QueueItemSplitResponse, error) {
			w, err := m.PodzielZlecenie(ctx, z)
			if err == nil {
				rozglosKolejkeZlecenia(ctx, m, e, z.QueueId)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandQueueItemMerge,
		obsluz(func(ctx context.Context, z shared.QueueItemMergeRequest) (shared.QueueItemMergeResponse, error) {
			w, err := m.ScalZlecenia(ctx, z)
			if err == nil {
				rozglosKolejkeZlecenia(ctx, m, e, z.QueueId)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandQueueItemRoute,
		obsluz(func(ctx context.Context, z shared.QueueItemRouteRequest) (shared.QueueItemRouteResponse, error) {
			w, err := m.SkierujZlecenie(ctx, z)
			if err == nil {
				// Skierowanie zmienia dwie kolejki: źródłowa traci zlecenie, docelowa je zyskuje, obie się odświeżają.
				rozglosKolejkeZlecenia(ctx, m, e, z.QueueId)
				if kod := wartoscTekstu(z.TargetQueueId); kod != "" && kod != z.QueueId {
					rozglosKolejkeZlecenia(ctx, m, e, kod)
				}
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandQueueItemBranch,
		obsluz(func(ctx context.Context, z shared.QueueItemBranchRequest) (shared.QueueItemBranchResponse, error) {
			w, err := m.RozgalezZlecenie(ctx, z)
			if err == nil {
				rozglosKolejkeZlecenia(ctx, m, e, z.QueueId)
				// Tor skierowany do innej kolejki zmienia także ją, inaczej Queue Manager pokazuje stan nieaktualny.
				for _, tor := range z.Branches {
					if kod := wartoscTekstu(tor.TargetQueueId); kod != "" && kod != z.QueueId {
						rozglosKolejkeZlecenia(ctx, m, e, kod)
					}
				}
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandQueueItemCondition,
		obsluz(func(ctx context.Context, z shared.QueueItemConditionRequest) (shared.QueueItemConditionResponse, error) {
			w, err := m.UwarunkujZlecenie(ctx, z)
			if err == nil {
				rozglosKolejkeZlecenia(ctx, m, e, z.QueueId)
			}
			return w, err
		}))
}

// rozglosKolejkeZlecenia dobiera kolejkę po identyfikatorze i rozgłasza jej
// zmianę. Nieudany dobór kończy wyłącznie rozgłoszenie — komenda już się
// powiodła i odmawianie jej z powodu zdarzenia byłoby odwróceniem porządku.
func rozglosKolejkeZlecenia(ctx context.Context, m zlecenioweKolejki, e *emiter, idKolejki string) {
	kolejka, err := m.Kolejka(ctx, idKolejki)
	if err != nil {
		return
	}
	e.kolejka(ctx, shared.ChangeKindUpdated, kolejka)
}

// zarejestrujOdmoweZlecenKolejek wpina wszystkie dwanaście komend jako odmowę montażu, tym samym powodem co przy wiązaniach: kolejki są, a rdzeń nie umie oddać ich zleceń.
func zarejestrujOdmoweZlecenKolejek(r *Rejestr) {
	const powod = "kolejki: port kolejek nie niesie zleceń ani polityki kolejki"

	for _, typ := range []shared.MessageType{
		shared.CommandQueueItemEnqueue,
		shared.CommandQueueItemDequeue,
		shared.CommandQueueItemDelay,
		shared.CommandQueueItemSplit,
		shared.CommandQueueItemMerge,
		shared.CommandQueueItemRoute,
		shared.CommandQueueItemBranch,
		shared.CommandQueueItemCondition,
		shared.CommandQueueItemList,
		shared.CommandQueuePolicySet,
		shared.CommandQueueDeadList,
		shared.CommandQueueDepthGet,
	} {
		r.Zarejestruj(typ, func(context.Context, protocol.Request) protocol.Odpowiedz {
			return porazka(bladMontazuKolejek(powod))
		})
	}
}
