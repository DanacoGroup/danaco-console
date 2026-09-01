// Plik wpina sześć komend rodziny schedule.* jako rozszerzenie portu Automatyki, bo harmonogram
// ma w rdzeniu jednego właściciela wspólnego z obszarem automation.*.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Harmonogramy jest portem rodziny `schedule.*` — portem modułu
// Automations rozszerzonym o odczyt harmonogramów.
type Harmonogramy interface {
	Automatyki

	// HarmonogramyZadania obsługuje `schedule.get`.
	HarmonogramyZadania(ctx context.Context, z shared.ScheduleGetRequest) (shared.ScheduleGetResponse, error)
}

// HarmonogramyNadzoru rozszerza port `Harmonogramy` o pięć czynności okna
// Scheduler, których stan MUSI przeżyć restart rdzenia: okna wykonania,
// uruchomienie wsteczne, historia wyzwoleń, nadzór obecności uruchomień
// i adres wejściowy wyzwalacza webhook.
type HarmonogramyNadzoru interface {
	Harmonogramy

	UstawOknaWykonania(ctx context.Context, z shared.ScheduleWindowSetRequest) (shared.ScheduleWindowSetResponse, error)
	UruchomWstecznie(ctx context.Context, z shared.ScheduleBackfillRunRequest) (shared.ScheduleBackfillRunResponse, error)
	HistoriaWyzwolen(ctx context.Context, z shared.ScheduleTriggerHistoryRequest) (shared.ScheduleTriggerHistoryResponse, error)
	UstawNadzorUruchomien(ctx context.Context, z shared.ScheduleHeartbeatSetRequest) (shared.ScheduleHeartbeatSetResponse, error)
	AdresWebhooka(ctx context.Context, z shared.ScheduleWebhookEndpointGetRequest) (shared.ScheduleWebhookEndpointGetResponse, error)
}

// zarejestrujHarmonogramy wpina sześć komend rodziny schedule.*. Zmiana harmonogramu rozgłasza
// zdarzenie zmiany powiązania automatyki, bo harmonogram jest bytem wyzwalającym automatykę.
func zarejestrujHarmonogramy(r *Rejestr, nadzor HarmonogramyNadzoru, e *emiter) {
	if r == nil || nadzor == nil {
		return
	}

	r.Zarejestruj(shared.CommandScheduleGet, obsluz(nadzor.HarmonogramyZadania))

	r.Zarejestruj(shared.CommandScheduleTriggerHistory, obsluz(nadzor.HistoriaWyzwolen))
	r.Zarejestruj(shared.CommandScheduleWebhookEndpointGet, obsluz(nadzor.AdresWebhooka))

	r.Zarejestruj(shared.CommandScheduleWindowSet,
		obsluz(func(ctx context.Context, z shared.ScheduleWindowSetRequest) (shared.ScheduleWindowSetResponse, error) {
			odpowiedz, err := nadzor.UstawOknaWykonania(ctx, z)
			if err == nil {
				e.powiazanieAutomatyki(shared.ChangeKindUpdated, odpowiedz.Schedule.WorkflowId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandScheduleHeartbeatSet,
		obsluz(func(ctx context.Context, z shared.ScheduleHeartbeatSetRequest) (shared.ScheduleHeartbeatSetResponse, error) {
			odpowiedz, err := nadzor.UstawNadzorUruchomien(ctx, z)
			if err == nil {
				e.powiazanieAutomatyki(shared.ChangeKindUpdated, odpowiedz.Schedule.WorkflowId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandScheduleBackfillRun, obsluz(nadzor.UruchomWstecznie))
}
