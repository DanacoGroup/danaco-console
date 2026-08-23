// Plik wpina cztery komendy rodziny `notification.*` — centrum powiadomień.
//
// Rodzina ma sam odczyt i zmianę stanu. Komendy zgłaszającej zdarzenie nie ma:
// powiadomienia zgłasza platforma tam, gdzie coś zaszło (`adapter_centrum_
// powiadomien.go` — `Zglos`), a komenda w kontrakcie pozwalałaby klientowi
// wpisać do rejestru zdarzenie, które nigdy nie zaszło.
//
// Zdarzeń `notification.raised` i `notification.changed` nie rozgłasza ten plik,
// tylko adapter — bo rozgłasza je także droga wewnętrzna, która przez rejestr
// komend nie przechodzi. Rozgłoszenie w dwóch miejscach dałoby przy zgłoszeniu
// dwie koperty o jednym zdarzeniu.
package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujCentrumPowiadomien wpina rodzinę `notification.*`.
//
// Port niewypełniony nie rejestruje niczego: komendy odpowiedzą wtedy
// `notification.unknown`, a pozostałe domeny pracują bez zmian.
func zarejestrujCentrumPowiadomien(r *Rejestr, centrum CentrumPowiadomien) {
	if r == nil || centrum == nil {
		return
	}

	r.Zarejestruj(shared.CommandNotificationList,
		obsluz(func(ctx context.Context, z shared.NotificationListRequest) (shared.NotificationListResponse, error) {
			return centrum.Wykaz(ctx, z)
		}))

	r.Zarejestruj(shared.CommandNotificationAcknowledge,
		obsluz(func(ctx context.Context, z shared.NotificationAcknowledgeRequest) (shared.NotificationAcknowledgeResponse, error) {
			return centrum.Odczytaj(ctx, z)
		}))

	r.Zarejestruj(shared.CommandNotificationResolve,
		obsluz(func(ctx context.Context, z shared.NotificationResolveRequest) (shared.NotificationResolveResponse, error) {
			return centrum.Zamknij(ctx, z)
		}))

	r.Zarejestruj(shared.CommandNotificationSnooze,
		obsluz(func(ctx context.Context, z shared.NotificationSnoozeRequest) (shared.NotificationSnoozeResponse, error) {
			return centrum.Odloz(ctx, z)
		}))
}
