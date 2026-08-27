// Odpowiedzialność pliku: cykl życia karty powłoki poza jej otwarciem — terminal.session.close i .list — oraz adres celu karty zdalnej żądania otwarcia.
package core

import (
	"context"
	"sort"
	"strings"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ZamknijKarte obsługuje terminal.session.close; procesy karty kończą się wyłącznie na wyraźne żądanie force, inaczej zostają biegnące.
func (a *adapterTerminala) ZamknijKarte(ctx context.Context,
	z shared.TerminalSessionCloseRequest) (shared.TerminalSessionCloseResponse, error) {

	karta, err := a.kartaZadania(z.SessionId)
	if err != nil {
		return shared.TerminalSessionCloseResponse{}, err
	}
	if karta.stan == shared.TerminalSessionStatusExited {
		return shared.TerminalSessionCloseResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeConflict, "moduł Terminal: karta "+karta.kod+" jest już zamknięta"))
	}

	zatrzymane := make([]string, 0, 4)
	if z.Force != nil && *z.Force {
		for _, proces := range a.rejestr.Procesy() {
			if proces.kartaKod != karta.kod {
				continue
			}
			if stan, _, _ := proces.Migawka(); stan != shared.TerminalProcessStatusRunning {
				continue
			}
			// Niepowodzenie jednego nie wstrzymuje pozostałych; nieubity proces nie wchodzi do wykazu.
			if err := proces.Zakoncz(true); err != nil {
				continue
			}
			proces.CzekajNaKoniec(czasNaDomkniecie)
			zatrzymane = append(zatrzymane, proces.kod)
		}
		sort.Strings(zatrzymane)
	}

	karta.stan = shared.TerminalSessionStatusExited
	a.rejestr.ZapiszKarte(karta)
	a.zapiszKarte(ctx, karta)

	odpowiedz := shared.TerminalSessionCloseResponse{Session: kartaKontraktu(karta)}
	if len(zatrzymane) > 0 {
		odpowiedz.StoppedProcessIds = zatrzymane
	}
	return odpowiedz, nil
}

// WykazKart obsługuje terminal.session.list, zwracając karty z rejestru pamięci i dziennika bazy od najnowszej.
func (a *adapterTerminala) WykazKart(ctx context.Context,
	z shared.TerminalSessionListRequest) (shared.TerminalSessionListResponse, error) {

	oknoKod := strings.TrimSpace(wartoscTekstu(z.WindowId))
	stan := shared.TerminalSessionStatus("")
	if z.Status != nil {
		stan = *z.Status
	}
	zZakonczonymi := z.IncludeExited != nil && *z.IncludeExited

	wykaz := make([]shared.TerminalSession, 0, 8)
	widziane := make(map[string]struct{}, 8)
	for _, karta := range a.rejestr.Karty() {
		widok := kartaKontraktu(karta)
		if !kartaPrzepuszczona(widok, oknoKod, stan, zZakonczonymi) {
			continue
		}
		widziane[widok.Id] = struct{}{}
		wykaz = append(wykaz, widok)
	}
	if a.repozytorium != nil {
		wiersze, err := a.repozytorium.Karty(ctx)
		if err != nil {
			return shared.TerminalSessionListResponse{}, err
		}
		for _, wiersz := range wiersze {
			if _, jest := widziane[wiersz.Kod]; jest {
				continue
			}
			widok := kartaWierszaKontraktu(wiersz)
			if !kartaPrzepuszczona(widok, oknoKod, stan, zZakonczonymi) {
				continue
			}
			wykaz = append(wykaz, widok)
		}
	}
	// Najnowsza karta stoi pierwsza — kontrakt mówi „od najnowszej”.
	sort.SliceStable(wykaz, func(i, j int) bool { return wykaz[i].CreatedAt > wykaz[j].CreatedAt })
	return shared.TerminalSessionListResponse{Sessions: wykaz, Total: len(wykaz)}, nil
}

// kartaPrzepuszczona sprawdza kartę trzema zawężeniami żądania, ze stanem wpisanym wprost mającym pierwszeństwo nad includeExited.
func kartaPrzepuszczona(karta shared.TerminalSession, oknoKod string,
	stan shared.TerminalSessionStatus, zZakonczonymi bool) bool {

	if oknoKod != "" && karta.WindowId != oknoKod {
		return false
	}
	if stan != "" {
		return karta.Status == stan
	}
	return zZakonczonymi || karta.Status != shared.TerminalSessionStatusExited
}
