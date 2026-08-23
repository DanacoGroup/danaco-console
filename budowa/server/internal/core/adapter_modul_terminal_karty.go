// Odpowiedzialność pliku: cykl życia karty powłoki poza jej otwarciem —
// komendy `terminal.session.close` i `terminal.session.list` — oraz adres celu
// karty zdalnej (pola `remoteTarget`, `remotePort` i `hostId` żądania otwarcia).
//
// ── Dlaczego zamknięcie karty jest czynnością rdzenia, a nie widoku ──────────
// Do tej pory zamknięcie karty żyło wyłącznie w kliencie: znikała zakładka,
// a powłoka i jej procesy biegły dalej, o czym rdzeń nie wiedział nic. Skutkiem
// było to, że `terminal.process.list` pokazywał procesy karty, której Operator
// już nie widzi, a `Przygotuj` po restarcie odtwarzał karty zamknięte tygodnie
// wcześniej. Zamknięcie ma więc wiersz i ma stan.
//
// ── Dlaczego wykaz kart składa się z dwóch źródeł ───────────────────────────
// Tak samo jak wykaz procesów (`adapter_modul_terminal_monitor.go`): karta
// czynna żyje w rejestrze pamięci, a karta zakończona zostaje w dzienniku bazy.
// Bez rejestru zniknęłyby karty właśnie otwarte na rdzeniu bez bazy, bez
// dziennika — karty zamknięte, o które kontrakt pyta polem `includeExited`.
package core

import (
	"context"
	"sort"
	"strings"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ZamknijKarte obsługuje `terminal.session.close`.
//
// Procesy karty kończą się WYŁĄCZNIE na wyraźne żądanie (`force`). Domyślne
// zamknięcie zostawia je biegnące, bo karta jest profilem powłoki, a nie
// właścicielem pracy: zamknięcie zakładki nie ma prawa przerwać budowania, które
// trwa trzecią minutę. Kontrakt oddaje wykaz zakończonych, żeby Operator wiedział,
// co dokładnie zatrzymał.
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
			// Niepowodzenie jednego zakończenia nie wstrzymuje pozostałych ani
			// samego zamknięcia karty: proces, którego nie udało się ubić, nie
			// wchodzi do wykazu i tym samym nie jest zgłoszony jako zatrzymany.
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

// WykazKart obsługuje `terminal.session.list`.
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

// kartaPrzepuszczona sprawdza kartę trzema zawężeniami żądania.
//
// Zawężenie stanem wpisane wprost bierze pierwszeństwo nad `includeExited`:
// pytanie o stan `exited` ma oddać karty zamknięte także wtedy, gdy pola
// `includeExited` nie podano — inaczej odpowiedź byłaby zawsze pusta.
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
