// Odpowiedzialność pliku: wpięcie dwóch komend rodziny `monitor.*`.
//
// Port jest rozszerzeniem portu nawigacji, a nie drugim portem: monitor
// procesów jest oknem warstwy wspólnej, którą obsługuje Nawigacja, więc
// telemetria ma w rdzeniu jednego właściciela. Wpięcie idzie osobno, bo rodzina
// `monitor.*` weszła do kontraktu osobno — tak samo jak `memory.*` rozszerza
// port modułu Workspace.
//
// Brak portu nie jest tu ciszą. Pozostałe rodziny przy porcie niewypełnionym
// nie rejestrują niczego i rdzeń odpowiada `*.unknown`; komendy monitora
// rejestrują się zawsze, a port bez telemetrii odpowiada `internal_error`
// z nazwą brakującego bytu. Odpowiedź ma brzmieć „monitor nie jest wpięty",
// a nie „nie znam takiej komendy" ani, najgorzej, pusty wykaz procesów.
//
// Zdarzeń ta rodzina nie rozgłasza: telemetria postępu ma w rdzeniu jednego
// producenta (`telemetria.go`, zdarzenie `progress.changed`). Odczyt stanu
// niczego nie zmienia, a zapis obserwacji zmienia pamięć rdzenia, dla której
// kontrakt zdarzenia nie ma. Dlatego `zarejestrujMonitor` nie bierze emitera.
package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// monitorProcesow jest portem rodziny `monitor.*` — portem nawigacji
// rozszerzonym o dwie komendy monitora procesów.
type monitorProcesow interface {
	Nawigacja

	// StanProcesow obsługuje `monitor.status`.
	StanProcesow(ctx context.Context, z shared.MonitorStatusRequest) (shared.MonitorStatusResponse, error)
	// ObserwujProcesy obsługuje `monitor.subscribe`.
	ObserwujProcesy(ctx context.Context, z shared.MonitorSubscribeRequest) (shared.MonitorSubscribeResponse, error)
}

// zarejestrujMonitor wpina dwie komendy rodziny `monitor.*`.
//
// Asercja jest miękka, a nie twarda jak przy `memory.*`, bo odmowa ma dojść do
// pytającego: rdzeń, który padłby przy montażu, powiedziałby to wyłącznie temu,
// kto czyta dziennik startu.
func zarejestrujMonitor(r *Rejestr, n Nawigacja) {
	if r == nil {
		return
	}
	m, niesie := n.(monitorProcesow)
	if !niesie {
		zarejestrujMonitorNiewpiety(r)
		return
	}
	r.Zarejestruj(shared.CommandMonitorStatus, obsluz(m.StanProcesow))
	r.Zarejestruj(shared.CommandMonitorSubscribe, obsluz(m.ObserwujProcesy))
}

// zarejestrujMonitorNiewpiety wpina odmowę na miejsce obsługiwaczy monitora.
// Komenda zostaje znana rdzeniowi i odpowiada wprost, czego brakuje.
func zarejestrujMonitorNiewpiety(r *Rejestr) {
	r.Zarejestruj(shared.CommandMonitorStatus,
		obsluz(func(context.Context, shared.MonitorStatusRequest) (shared.MonitorStatusResponse, error) {
			return shared.MonitorStatusResponse{}, bladMonitoraNiewpietego()
		}))
	r.Zarejestruj(shared.CommandMonitorSubscribe,
		obsluz(func(context.Context, shared.MonitorSubscribeRequest) (shared.MonitorSubscribeResponse, error) {
			return shared.MonitorSubscribeResponse{}, bladMonitoraNiewpietego()
		}))
}

// bladMonitoraNiewpietego składa odmowę portu, który nie niesie monitora.
func bladMonitoraNiewpietego() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"monitor procesów: port nawigacji nie niesie telemetrii postępu — "+
			"stanu procesów nie ma skąd odczytać"))
}
