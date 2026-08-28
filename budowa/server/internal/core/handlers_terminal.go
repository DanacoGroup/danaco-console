// Plik wpina cztery komendy obszaru terminal.* obsługujące moduł Terminal wraz z jego trzema
// oknami operacyjnymi: kartami terminala, monitorem procesów i konsolą wyjścia.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Terminal jest portem modułu Terminal obsługującym komendy uruchomienia i sterowania procesem powłoki.
type Terminal interface {
	OtworzKarte(ctx context.Context, z shared.TerminalSessionOpenRequest) (shared.TerminalSessionOpenResponse, error)
	WykonajPolecenie(ctx context.Context, z shared.TerminalCommandExecRequest) (shared.TerminalCommandExecResponse, error)
	WykazProcesow(ctx context.Context, z shared.TerminalProcessListRequest) (shared.TerminalProcessListResponse, error)
	ZakonczProces(ctx context.Context, z shared.TerminalProcessKillRequest) (shared.TerminalProcessKillResponse, error)
	// PodepnijRozgloszenie oddaje adapterowi drogę do zdarzenia zmiany procesu.
	PodepnijRozgloszenie(rozglos func(shared.ChangeKind, shared.TerminalProcess))
}

// zarejestrujTerminal wpina cztery komendy modułu Terminal obsługujące karty, procesy i wyjście konsoli.
func zarejestrujTerminal(r *Rejestr, t Terminal, e *emiter) {
	if r == nil || t == nil {
		return
	}
	t.PodepnijRozgloszenie(e.procesTerminala)

	// Rozgłoszenia po wykonaniu komendy i zabiciu procesu nie ma tutaj: zdarzenie nadaje adapter.
	r.Zarejestruj(shared.CommandTerminalSessionOpen, obsluz(t.OtworzKarte))
	r.Zarejestruj(shared.CommandTerminalCommandExec, obsluz(t.WykonajPolecenie))
	r.Zarejestruj(shared.CommandTerminalProcessList, obsluz(t.WykazProcesow))
	r.Zarejestruj(shared.CommandTerminalProcessKill, obsluz(t.ZakonczProces))
}

// procesTerminala rozgłasza zmianę procesu rejestru rdzenia. Sesja komunikatu
// bywa pusta: proces okna, którego sesji rdzeń już nie zna, rozgłasza się bez
// niej, zamiast nie rozgłaszać się wcale.
func (e *emiter) procesTerminala(zmiana shared.ChangeKind, p shared.TerminalProcess) {
	e.wyslij(shared.EventTerminalProcessChanged, "",
		shared.TerminalProcessChangedEvent{Change: zmiana, Process: p})
}

// PodepnijRozgloszenie wypełnia port: adapter zapamiętuje drogę do zdarzenia zmiany procesu w terminalu.
func (a *adapterTerminala) PodepnijRozgloszenie(rozglos func(shared.ChangeKind, shared.TerminalProcess)) {
	a.zmiana = rozglos
}
