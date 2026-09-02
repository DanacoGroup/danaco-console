// Odpowiedzialność pliku: bieg procesu terminala od startu do domknięcia.
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

func (a *adapterTerminala) uruchom(ctx context.Context, karta *kartaTerminala, okno session.Okno,
	polecenie session.Polecenie, tresc string, inicjator shared.ProcessInitiator) (*procesTerminala, error) {

	if a.uruchamiacz == nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Terminal: serwer nie ma uruchamiacza procesów"))
	}
	uchwyt, err := a.uruchamiacz.UruchomProces(ctx, okno, polecenie)
	if err != nil {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Terminal: nie można uruchomić polecenia w powłoce "+string(karta.powloka)+": "+err.Error()))
	}
	drzewo, err := session.PrzejmijDrzewo(uchwyt.Pid())
	if err != nil {
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
		kontekst:     zKontemZadania(context.Background(), ctx),
		stan:         shared.TerminalProcessStatusRunning,
		uruchomiono:  time.Now().UTC(),
		uchwyt:       uchwyt,
		drzewo:       drzewo,
		koniec:       make(chan struct{}),
	}
	a.rejestr.ZapiszProces(proces)
	a.zapiszProces(ctx, proces)

	// Rozgłoszenie idzie przed pompą wyjścia, przed odrzuceniem fragmentu jako cudzego.
	a.rozglos(ctx, shared.ChangeKindCreated, proces)

	go a.wyjscie.Pompuj(proces, uchwyt.Wyjscie(), shared.ChunkKindText)
	go a.wyjscie.Pompuj(proces, uchwyt.Diagnostyka(), shared.ChunkKindError)
	return proces, nil
}

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

func (a *adapterTerminala) domknij(proces *procesTerminala, blad error, powod string) {
	stan, kodWyjscia := wynikZakonczenia(blad)
	if !proces.Domknij(stan, kodWyjscia) {
		stan, kodWyjscia, _ = proces.Migawka()
	}
	proces.Zwolnij()
	close(proces.koniec)
	a.rejestr.Przytnij()

	if a.repozytorium != nil {
		_ = a.repozytorium.ZakonczProces(proces.kontekst, proces.kod, stan, wskaznikDuzej(kodWyjscia))
	}
	a.wyjscie.Domknij(proces, stan, kodWyjscia, powod)
	a.rozglos(proces.kontekst, shared.ChangeKindUpdated, proces)
}

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

func (a *adapterTerminala) rozglos(ctx context.Context, zmiana shared.ChangeKind, proces *procesTerminala) {
	if a.zmiana == nil {
		return
	}
	a.zmiana(ctx, zmiana, procesKontraktu(proces))
}

func (p *procesTerminala) CzekajNaKoniec(najdluzej time.Duration) {
	select {
	case <-p.koniec:
	case <-time.After(najdluzej):
	}
}

func granicaCzasu(milisekundy *int) time.Duration {
	if milisekundy == nil || *milisekundy <= 0 {
		return granicaCzasuDomyslna
	}
	return time.Duration(*milisekundy) * time.Millisecond
}

func inicjatorZadania(inicjator *shared.ProcessInitiator) shared.ProcessInitiator {
	if inicjator == nil || *inicjator == "" {
		return shared.ProcessInitiatorOperator
	}
	return *inicjator
}
