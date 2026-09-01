// Komenda `terminal.process.suspend` — wstrzymanie i wznowienie procesu
// rejestru rdzenia, oddające procesor bez utraty postępu. Stan procesu nie
// zmienia się na wstrzymany: zgłasza to pole `supported` odpowiedzi.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// WstrzymajProces obsługuje `terminal.process.suspend`, wstrzymując albo
// wznawiając drzewo procesu wskazane identyfikatorem.
func (a *adapterTerminala) WstrzymajProces(ctx context.Context,
	z shared.TerminalProcessSuspendRequest) (shared.TerminalProcessSuspendResponse, error) {

	kod := strings.TrimSpace(z.ProcessId)
	if kod == "" {
		return shared.TerminalProcessSuspendResponse{}, bladZadaniaTerminala(
			"wstrzymanie procesu wymaga jego wskazania")
	}
	proces, jest := a.rejestr.Proces(kod)
	if !jest {
		return shared.TerminalProcessSuspendResponse{}, a.procesZDziennika(ctx, kod)
	}
	if stan, _, _ := proces.Migawka(); stan != shared.TerminalProcessStatusRunning {
		return shared.TerminalProcessSuspendResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeConflict,
			"moduł Terminal: proces "+kod+" jest już zakończony (stan "+string(stan)+
				"), więc nie ma czego wstrzymać ani wznowić"))
	}

	wznowienie := z.Resume != nil && *z.Resume
	wspierane, err := proces.Wstrzymaj(wznowienie)
	if err != nil {
		return shared.TerminalProcessSuspendResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError, "moduł Terminal: nie można zmienić biegu procesu "+
				kod+": "+err.Error()))
	}
	if wspierane {
		// Zmiana biegu jest zmianą stanu widoczną w Process Monitorze.
		a.rozglos(ctx, shared.ChangeKindUpdated, proces)
	}
	return shared.TerminalProcessSuspendResponse{
		Process:   procesKontraktu(proces),
		Supported: wspierane,
	}, nil
}

// Wstrzymaj zatrzymuje albo wznawia drzewo procesu i zapamiętuje jego bieg.
// Powtórzenie czynności jest ciche i udane, idzie do systemu jak pierwsze
// wywołanie, bez skrótu po zapamiętanym biegu.
func (p *procesTerminala) Wstrzymaj(wznowienie bool) (bool, error) {
	var wspierane bool
	var err error
	if wznowienie {
		wspierane, err = p.drzewo.Wznow()
	} else {
		wspierane, err = p.drzewo.Wstrzymaj()
	}
	if err != nil || !wspierane {
		return wspierane, err
	}
	p.mu.Lock()
	p.wstrzymany = !wznowienie
	p.mu.Unlock()
	return true, nil
}
