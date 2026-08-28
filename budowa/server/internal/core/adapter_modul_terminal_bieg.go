// Odpowiedzialność pliku: bieg procesu terminala od startu do domknięcia —
// uruchomienie, przejęcie drzewa potomstwa, pompa wyjścia, granica czasu
// i zapis wyniku. Obserwator zakończenia pracuje poza żądaniem, zdarzeniem.
package core

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// uruchom startuje proces przez port warstwy kanału, obejmuje go drzewem
// potomstwa i podpina pompę wyjścia.
func (a *adapterTerminala) uruchom(ctx context.Context, karta *kartaTerminala, okno session.Okno,
	polecenie session.Polecenie, tresc string, inicjator shared.ProcessInitiator) (*procesTerminala, error) {

	if a.uruchamiacz == nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Terminal: rdzeń nie ma uruchamiacza procesów"))
	}
	uchwyt, err := a.uruchamiacz.UruchomProces(okno, polecenie)
	if err != nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Terminal: nie można uruchomić polecenia w powłoce "+string(karta.powloka)+": "+err.Error()))
	}
	drzewo, err := session.PrzejmijDrzewo(uchwyt.Pid())
	if err != nil {
		// Proces już biegnie, a uchwytu drzewa nie ma — zostawienie znaczyłoby sierotę.
		_ = uchwyt.Ubij()
		_ = uchwyt.Czekaj()
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Terminal: nie można objąć drzewa procesu: "+err.Error()))
	}

	proces := &procesTerminala{
		kod:          nowyIdentyfikator(przedrostekProcesuTerminala),
		kartaKod:     karta.kod,
		oknoKod:      karta.oknoKod,
		idSesji:      okno.IdSesji,
		polecenie:    tresc,
		inicjator:    inicjator,
		pid:          uchwyt.Pid(),
		pidNadrzedny: os.Getpid(),
		stan:         shared.TerminalProcessStatusRunning,
		uruchomiono:  time.Now().UTC(),
		uchwyt:       uchwyt,
		drzewo:       drzewo,
		koniec:       make(chan struct{}),
	}
	a.rejestr.ZapiszProces(proces)
	a.zapiszProces(ctx, proces)

	// Rozgłoszenie idzie przed pompą wyjścia, przed odrzuceniem fragmentu jako cudzego.
	a.rozglos(shared.ChangeKindCreated, proces)

	go a.wyjscie.Pompuj(proces, uchwyt.Wyjscie(), shared.ChunkKindText)
	go a.wyjscie.Pompuj(proces, uchwyt.Diagnostyka(), shared.ChunkKindError)
	return proces, nil
}

// pilnujZakonczenia czeka na koniec procesu we własnej gorutynie, domyka jego
// stan i rozgłasza zmianę. Granica czasu większa od zera ubija drzewo procesu
// po jej upływie.
func (a *adapterTerminala) pilnujZakonczenia(proces *procesTerminala, granica time.Duration) {
	zakonczenie := make(chan error, 1)
	go func() { zakonczenie <- proces.uchwyt.Czekaj() }()

	go func() {
		var powod string
		var blad error
		if granica > 0 {
			zegar := time.NewTimer(granica)
			defer zegar.Stop()
			select {
			case blad = <-zakonczenie:
			case <-zegar.C:
				powod = "przekroczona granica czasu " + granica.String()
				_ = proces.Zakoncz(true)
				blad = <-zakonczenie
			}
		} else {
			blad = <-zakonczenie
		}
		a.domknij(proces, blad, powod)
	}()
}

// domknij ustala stan końcowy procesu, odkłada go do dziennika, zamyka strumień
// wyjścia i rozgłasza zmianę.
func (a *adapterTerminala) domknij(proces *procesTerminala, blad error, powod string) {
	stan, kodWyjscia := wynikZakonczenia(blad)
	if !proces.Domknij(stan, kodWyjscia) {
		stan, kodWyjscia, _ = proces.Migawka()
	}
	proces.Zwolnij()
	close(proces.koniec)
	a.rejestr.Przytnij()

	if a.repozytorium != nil {
		_ = a.repozytorium.ZakonczProces(context.Background(), proces.kod, stan, wskaznikDuzej(kodWyjscia))
	}
	a.wyjscie.Domknij(proces, stan, kodWyjscia, powod)
	a.rozglos(shared.ChangeKindUpdated, proces)
}

// wynikZakonczenia przekłada wynik oczekiwania na stan i kod wyjścia. Kod
// różny od zera jest wynikiem polecenia, nie usterką rdzenia.
func wynikZakonczenia(blad error) (shared.TerminalProcessStatus, *int) {
	if blad == nil {
		zero := 0
		return shared.TerminalProcessStatusFinished, &zero
	}
	var zakonczenie *exec.ExitError
	if errors.As(blad, &zakonczenie) {
		kod := zakonczenie.ExitCode()
		if kod < 0 {
			// Kod ujemny znaczy zakończenie sygnałem, nie zawodem własnego wyniku.
			return shared.TerminalProcessStatusStopped, nil
		}
		return shared.TerminalProcessStatusFailed, &kod
	}
	return shared.TerminalProcessStatusFailed, nil
}

// zapiszProces odkłada uruchomiony proces do dziennika. Nieudany zapis nie
// zatrzymuje procesu, który już biegnie.
func (a *adapterTerminala) zapiszProces(ctx context.Context, proces *procesTerminala) {
	if a.repozytorium == nil {
		return
	}
	pid, rodzic := int64(proces.pid), int64(proces.pidNadrzedny)
	_ = a.repozytorium.ZapiszProces(ctx, dane.ProcesTerminala{
		Kod:          proces.kod,
		KartaKod:     wskaznikTekstu(proces.kartaKod),
		OknoKod:      proces.oknoKod,
		Pid:          &pid,
		PidNadrzedny: &rodzic,
		Polecenie:    proces.polecenie,
		Inicjator:    proces.inicjator,
		Stan:         shared.TerminalProcessStatusRunning,
	})
}

// rozglos oddaje zmianę procesu obsługiwaczowi, który rozsyła
// `terminal.process.changed`. Brak podpięcia nie zmienia pracy modułu.
func (a *adapterTerminala) rozglos(zmiana shared.ChangeKind, proces *procesTerminala) {
	if a.zmiana == nil {
		return
	}
	a.zmiana(zmiana, procesKontraktu(proces))
}

// CzekajNaKoniec czeka na domknięcie procesu nie dłużej niż podany czas,
// po jego upływie oddając ostatni znany stan.
func (p *procesTerminala) CzekajNaKoniec(najdluzej time.Duration) {
	select {
	case <-p.koniec:
	case <-time.After(najdluzej):
	}
}

// granicaCzasu czyta granicę z żądania; wartość niedodatnia oddaje granicę
// domyślną (`granicaCzasuDomyslna`).
func granicaCzasu(milisekundy *int) time.Duration {
	if milisekundy == nil || *milisekundy <= 0 {
		return granicaCzasuDomyslna
	}
	return time.Duration(*milisekundy) * time.Millisecond
}

// inicjatorZadania czyta inicjatora z żądania; brak wskazania znaczy Operatora,
// wartość domyślną inicjatora poleceń terminala.
func inicjatorZadania(inicjator *shared.ProcessInitiator) shared.ProcessInitiator {
	if inicjator == nil || *inicjator == "" {
		return shared.ProcessInitiatorOperator
	}
	return *inicjator
}
