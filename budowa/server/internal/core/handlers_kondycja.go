// Plik wpina sześć komend rodziny `health.*`: definicje sond, ich przebiegi,
// serię wyników i dostępność, bez rozgłaszania zdarzeń.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Kondycja jest portem rodziny `health.*`. Mówi wyłącznie typami kontraktu;
// gniazda, adresy, programy sondujące i wiersze serii leżą po drugiej stronie
// adaptera.
type Kondycja interface {
	// ZapiszSonde obsługuje `health.probe.save`.
	ZapiszSonde(ctx context.Context, z shared.HealthProbeSaveRequest) (shared.HealthProbeSaveResponse, error)
	// WykazSond obsługuje `health.probe.list`.
	WykazSond(ctx context.Context, z shared.HealthProbeListRequest) (shared.HealthProbeListResponse, error)
	// UsunSonde obsługuje `health.probe.remove`.
	UsunSonde(ctx context.Context, z shared.HealthProbeRemoveRequest) (shared.HealthProbeRemoveResponse, error)
	// WykonajSonde obsługuje `health.probe.run` — pomiar na żądanie.
	WykonajSonde(ctx context.Context, z shared.HealthProbeRunRequest) (shared.HealthProbeRunResponse, error)
	// WykazWynikow obsługuje `health.result.list`.
	WykazWynikow(ctx context.Context, z shared.HealthResultListRequest) (shared.HealthResultListResponse, error)
	// Dostepnosc obsługuje `health.uptime.get`.
	Dostepnosc(ctx context.Context, z shared.HealthUptimeGetRequest) (shared.HealthUptimeGetResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ Kondycja = (*adapterKondycji)(nil)

// zarejestrujKondycje wpina sześć komend rodziny `health.*`, obsługujących
// sondy, ich przebiegi, serię wyników i dostępność.
func zarejestrujKondycje(r *Rejestr, k Kondycja) {
	if r == nil || k == nil {
		return
	}

	r.Zarejestruj(shared.CommandHealthProbeSave, obsluz(k.ZapiszSonde))
	r.Zarejestruj(shared.CommandHealthProbeList, obsluz(k.WykazSond))
	r.Zarejestruj(shared.CommandHealthProbeRemove, obsluz(k.UsunSonde))
	r.Zarejestruj(shared.CommandHealthProbeRun, obsluz(k.WykonajSonde))
	r.Zarejestruj(shared.CommandHealthResultList, obsluz(k.WykazWynikow))
	r.Zarejestruj(shared.CommandHealthUptimeGet, obsluz(k.Dostepnosc))
}
