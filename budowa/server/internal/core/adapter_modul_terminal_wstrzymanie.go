// Komenda `terminal.process.suspend` — wstrzymanie i wznowienie procesu
// rejestru rdzenia.
//
// ── Czego brakowało ─────────────────────────────────────────────────────────
// Rdzeń umiał proces wyłącznie ZAKOŃCZYĆ, więc jedyną odpowiedzią na zadanie,
// które zajęło maszynę, było wyrzucenie wykonanej przez nie pracy. Wstrzymanie
// oddaje procesor bez utraty postępu (`session/wstrzymanie.go`).
//
// ── Dlaczego stan procesu nie zmienia się na „wstrzymany" ───────────────────
// Bo takiego stanu nie ma w kontrakcie: `TerminalProcessStatus` zna `running`,
// `finished`, `failed` i `stopped`. Proces wstrzymany JEST wciąż uruchomiony —
// ma PID, pamięć i otwarte pliki — więc `running` jest o nim prawdą, a `stopped`
// byłoby nieprawdą, bo ten stan oznacza w tym module zakończenie sygnałem.
// Wstrzymania nie zgłaszamy więc stanem, którego kontrakt nie ma; zgłasza je
// pole `supported` odpowiedzi wraz z powtarzalnością samej czynności.
//
// ── Dlaczego odpowiedź ma `supported`, a nie odmowę ─────────────────────────
// Windows nie zna wstrzymania obcego drzewa procesów. „System tego nie umie" to
// co innego niż „czynność zawiodła", i kontrakt mówi to wprost: fałsz znaczy, że
// proces został nietknięty. Odmowa w tym miejscu kazałaby klientowi zgadywać,
// czy proces jednak nie stanął.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// WstrzymajProces obsługuje `terminal.process.suspend`.
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
		// Zmiana biegu jest zmianą stanu procesu widoczną w Process Monitorze,
		// więc idzie tą samą drogą co uruchomienie i zakończenie.
		a.rozglos(shared.ChangeKindUpdated, proces)
	}
	return shared.TerminalProcessSuspendResponse{
		Process:   procesKontraktu(proces),
		Supported: wspierane,
	}, nil
}

// Wstrzymaj zatrzymuje albo wznawia drzewo procesu i zapamiętuje jego bieg.
//
// Powtórzenie czynności jest ciche i udane, i dlatego idzie do systemu tak samo
// jak pierwsze wywołanie, bez skrótu po zapamiętanym biegu: SIGSTOP dla procesu
// już wstrzymanego i SIGCONT dla biegnącego nie robią nic, a odpowiedź o wsparciu
// platformy ma pochodzić od platformy, nie od pamięci rdzenia. Skrót
// odpowiadałby „niewspierane” na wznowienie procesu, którego nikt nie wstrzymał.
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
