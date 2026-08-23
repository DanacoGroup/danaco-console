// Odpowiedzialność pliku: wpięcie czterech komend obszaru `terminal.*` —
// modułu Terminal wraz z jego trzema oknami operacyjnymi (Terminal Tabs,
// Process Monitor, Output Console).
//
// Jedno zdarzenie na cały moduł. Kontrakt daje modułowi wyłącznie
// `terminal.process.changed`, więc każda zmiana stanu procesu — uruchomienie,
// zakończenie, ubicie — rozgłasza się procesem po zmianie. Process Monitor
// odświeża się z jednej subskrypcji, a nie z odpytywania.
//
// Wyjście procesu nie idzie tą drogą. Strumień Output Console jedzie zdarzeniem
// `stream.chunk`, bo trwa dłużej niż wykonanie komendy i bo kontrakt
// nie ma zapowiadanego przez wykaz `terminal.output.stream`.
// Składa go adapter_modul_terminal_strumien.go.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Terminal jest portem modułu Terminal.
type Terminal interface {
	OtworzKarte(ctx context.Context, z shared.TerminalSessionOpenRequest) (shared.TerminalSessionOpenResponse, error)
	WykonajPolecenie(ctx context.Context, z shared.TerminalCommandExecRequest) (shared.TerminalCommandExecResponse, error)
	WykazProcesow(ctx context.Context, z shared.TerminalProcessListRequest) (shared.TerminalProcessListResponse, error)
	ZakonczProces(ctx context.Context, z shared.TerminalProcessKillRequest) (shared.TerminalProcessKillResponse, error)
	// PodepnijRozgloszenie oddaje adapterowi drogę do zdarzenia zmiany procesu.
	// Proces kończy się poza wykonaniem komendy — czasem godziny później — więc
	// rozgłoszenie nie może iść wyłącznie z obsługiwacza żądania.
	PodepnijRozgloszenie(rozglos func(shared.ChangeKind, shared.TerminalProcess))
}

// zarejestrujTerminal wpina cztery komendy modułu Terminal.
func zarejestrujTerminal(r *Rejestr, t Terminal, e *emiter) {
	if r == nil || t == nil {
		return
	}
	t.PodepnijRozgloszenie(e.procesTerminala)

	// Rozgłoszenia po `terminal.command.exec` i `terminal.process.kill` nie ma
	// tutaj z zamysłem: zdarzenie nadaje adapter, bo tylko on wie, kiedy proces
	// naprawdę zmienił stan. Podwójne rozgłoszenie z obsługiwacza dałoby
	// Process Monitorowi ten sam wiersz dwa razy.
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

// PodepnijRozgloszenie wypełnia port: adapter zapamiętuje drogę do zdarzenia.
func (a *adapterTerminala) PodepnijRozgloszenie(rozglos func(shared.ChangeKind, shared.TerminalProcess)) {
	a.zmiana = rozglos
}
