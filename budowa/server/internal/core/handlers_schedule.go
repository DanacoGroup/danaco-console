// Plik wpina sześć komend rodziny schedule.* jako rozszerzenie portu Automatyki, bo harmonogram
// ma w rdzeniu jednego właściciela wspólnego z obszarem automation.*.
package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
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
func zarejestrujHarmonogramy(r *Rejestr, m Harmonogramy, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandScheduleGet, obsluz(m.HarmonogramyZadania))

	nadzor, ok := m.(HarmonogramyNadzoru)
	if !ok {
		zarejestrujOdmoweNadzoruHarmonogramow(r)
		return
	}

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

// zarejestrujOdmoweNadzoruHarmonogramow wpina pięć komend nadzoru jako odmowę
// montażu — port harmonogramów jest, lecz nie niesie ich czynności.
func zarejestrujOdmoweNadzoruHarmonogramow(r *Rejestr) {
	const powod = "moduł Automations: port harmonogramów nie niesie nadzoru ani okien wykonania"

	for _, typ := range []shared.MessageType{
		shared.CommandScheduleWindowSet,
		shared.CommandScheduleBackfillRun,
		shared.CommandScheduleTriggerHistory,
		shared.CommandScheduleHeartbeatSet,
		shared.CommandScheduleWebhookEndpointGet,
	} {
		r.Zarejestruj(typ, func(context.Context, protocol.Request) protocol.Odpowiedz {
			return porazka(bladNiedostepnegoSilnika(powod))
		})
	}
}
