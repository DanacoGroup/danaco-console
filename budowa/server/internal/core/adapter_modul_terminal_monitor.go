// Odpowiedzialność pliku: dwie komendy okna Process Monitor — terminal.process.list i .kill, na wykazie procesów czynnego rejestru i dziennika bazy.
package core

import (
	"context"
	"errors"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// WykazProcesow obsługuje `terminal.process.list`. Wykaz składa się z procesów
// czynnych rejestru i z dziennika przebiegów zakończonych.
func (a *adapterTerminala) WykazProcesow(ctx context.Context,
	z shared.TerminalProcessListRequest) (shared.TerminalProcessListResponse, error) {

	filtr := filtrProcesow(z)
	wykaz := make([]shared.TerminalProcess, 0, 16)
	widziane := make(map[string]struct{}, 16)

	for _, proces := range a.rejestr.Procesy() {
		widok := procesKontraktu(proces)
		if !filtrPrzepuszcza(filtr, widok) {
			continue
		}
		widziane[widok.Id] = struct{}{}
		wykaz = append(wykaz, widok)
	}
	if a.repozytorium != nil {
		wiersze, err := a.repozytorium.Procesy(ctx, filtr)
		if err != nil {
			return shared.TerminalProcessListResponse{}, err
		}
		for _, wiersz := range wiersze {
			if _, jest := widziane[wiersz.Kod]; jest {
				continue
			}
			wykaz = append(wykaz, procesWierszaKontraktu(wiersz))
		}
	}
	// Najnowszy przebieg stoi pierwszy — Process Monitor czyta od góry.
	sort.SliceStable(wykaz, func(i, j int) bool { return wykaz[i].StartedAt > wykaz[j].StartedAt })
	return shared.TerminalProcessListResponse{Processes: wykaz}, nil
}

// ZakonczProces obsługuje terminal.process.kill, kończąc proces terminala i czekając krótko na stan końcowy.
func (a *adapterTerminala) ZakonczProces(ctx context.Context,
	z shared.TerminalProcessKillRequest) (shared.TerminalProcessKillResponse, error) {

	kod := strings.TrimSpace(z.ProcessId)
	if kod == "" {
		return shared.TerminalProcessKillResponse{}, bladZadaniaTerminala("zakończenie procesu wymaga jego wskazania")
	}
	proces, jest := a.rejestr.Proces(kod)
	if !jest {
		return shared.TerminalProcessKillResponse{}, a.procesZDziennika(ctx, kod)
	}
	if stan, _, _ := proces.Migawka(); stan != shared.TerminalProcessStatusRunning {
		return shared.TerminalProcessKillResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeConflict,
			"moduł Terminal: proces "+kod+" jest już zakończony (stan "+string(stan)+")"))
	}
	wymuszony := z.Force != nil && *z.Force
	if err := proces.Zakoncz(wymuszony); err != nil {
		return shared.TerminalProcessKillResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError, "moduł Terminal: nie można zakończyć procesu "+kod+": "+err.Error()))
	}
	// Stan końcowy nadaje obserwator zakończenia; krótkie oczekiwanie zapobiega odpowiedzi running.
	proces.CzekajNaKoniec(czasNaDomkniecie)
	return shared.TerminalProcessKillResponse{Process: procesKontraktu(proces)}, nil
}

// procesZDziennika odpowiada na zakończenie procesu, którego nie ma już
// w rejestrze: sięga do dziennika i zwraca konflikt ze stanem końcowym zamiast
// komunikatu o braku procesu.
func (a *adapterTerminala) procesZDziennika(ctx context.Context, kod string) error {
	if a.repozytorium == nil {
		return bladBrakuProcesu(kod)
	}
	wiersz, err := a.repozytorium.Proces(ctx, kod)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return bladBrakuProcesu(kod)
	}
	if err != nil {
		return err
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"moduł Terminal: proces "+kod+" jest już zakończony (stan "+string(wiersz.Stan)+")"))
}
