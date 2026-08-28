// Plik wpina dwie komendy rodziny monitor.* jako rozszerzenie portu
// nawigacji; port bez telemetrii odpowiada odmową zamiast milczeć.
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

// bladMonitoraNiewpietego składa odmowę portu nawigacji, który nie niesie
// monitora procesów, zamiast zwracać pustą odpowiedź.
func bladMonitoraNiewpietego() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"monitor procesów: port nawigacji nie niesie telemetrii postępu — "+
			"stanu procesów nie ma skąd odczytać"))
}
