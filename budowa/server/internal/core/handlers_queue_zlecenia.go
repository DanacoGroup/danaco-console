// Wpięcie dwunastu komend rodziny `queue.*` dotyczących ZLECEŃ kolejki, jej
// polityki, zadań martwych i głębokości w czasie.
//
// Osobno od `handlers_queue.go` (cykl życia kolejki) i `handlers_queue_wiazania.go`
// (wykaz i wiązanie), bo osobny jest przedmiot: tam bytem jest kolejka, tutaj
// zlecenie. Port jest rozszerzeniem portu `Kolejki`, nie drugim portem —
// silnik kolejek pętli sesyjnej, MultitaskingAI i modułu Automations jest jeden.
//
// Rozgłoszenie idzie `queue.changed` wszędzie tam, gdzie zmienia się zawartość
// albo polityka kolejki: Queue Manager rysuje wykaz zleceń i wskaźnik
// głębokości z tego, co o kolejce wie, więc dołożone, zdjęte, odłożone,
// podzielone, scalone, skierowane, rozgałęzione i uwarunkowane zlecenie czyni
// jego obraz nieaktualnym. Trzy odczyty — wykaz zleceń, zadania martwe
// i głębokość — nie rozgłaszają niczego.
//
// Rodzajem zmiany jest `updated`, nie `created`: bytem zdarzenia
// `queue.changed` jest KOLEJKA, a kolejka przy dołożeniu zlecenia nie powstaje,
// tylko się zmienia. `created` opisywałoby założenie samej kolejki i tak jest
// używane w `queue.create`.
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

	// Kolejka oddaje kolejkę po identyfikatorze — nośnik odpowiedzi tych
	// czynności, których wynikiem jest kolejka po zmianie.
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

// zarejestrujZleceniaKolejek wpina dwanaście komend zleceń kolejki.
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
				e.kolejka(shared.ChangeKindUpdated, w.Queue)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandQueueItemEnqueue,
		obsluz(func(ctx context.Context, z shared.QueueItemEnqueueRequest) (shared.QueueItemEnqueueResponse, error) {
			w, err := zleceniowe.DodajZlecenie(ctx, z)
			// Duplikat niczego nie zmienił, więc niczego nie rozgłasza:
			// zdarzenie po wywołaniu odbitym idempotencją byłoby zdarzeniem
			// bez faktu.
			if err == nil && !w.Duplicate {
				rozglosKolejkeZlecenia(ctx, zleceniowe, e, z.QueueId)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandQueueItemDequeue,
		obsluz(func(ctx context.Context, z shared.QueueItemDequeueRequest) (shared.QueueItemDequeueResponse, error) {
			w, err := zleceniowe.ZdejmijZlecenie(ctx, z)
			if err == nil && w.Removed {
				e.kolejka(shared.ChangeKindUpdated, w.Queue)
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
				// Skierowanie zmienia DWIE kolejki: źródłową traci zlecenie,
				// docelowa je zyskuje. Obie muszą się odświeżyć, więc obie
				// dostają zdarzenie.
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
				// Tor skierowany do innej kolejki zmienia także ją — bez tego
				// Queue Manager kolejki docelowej pokazywałby stan sprzed
				// rozgałęzienia aż do ręcznego odświeżenia.
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
	e.kolejka(shared.ChangeKindUpdated, kolejka)
}

// zarejestrujOdmoweZlecenKolejek wpina wszystkie dwanaście komend jako odmowę
// montażu. Powód ten sam, co przy wiązaniach kolejek: kolejki są, a rdzeń nie
// umie oddać ich zleceń — odpowiedź „nieznana komenda" wskazywałaby na brak
// kolejek, a nie na usterkę montażu.
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
