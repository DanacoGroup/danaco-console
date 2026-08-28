// Plik rejestruje cztery komendy rodziny `notification.*`, obsługujące odczyt
// i zmianę stanu centrum powiadomień, bez komendy zgłaszającej zdarzenie.
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
