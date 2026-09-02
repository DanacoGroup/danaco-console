// Odpowiedzialność pliku: bieg przebiegu budowania od startu do domknięcia — uruchomienie przez port kanału, przejęcie drzewa potomstwa, pompa logu i zapis wyniku.
package core

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os/exec"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// czasNaDomknieciePrzebiegu ogranicza czekanie komendy przerwania; stan końcowy dochodzi zdarzeniem developer.build.changed.
const czasNaDomknieciePrzebiegu = 750 * time.Millisecond

// najdluzszyWierszLogu tnie wyjście bez znaku końca wiersza (pasek postępu).
const najdluzszyWierszLogu = 64 * 1024

func (a *adapterDevelopera) uruchomBudowanie(ctx context.Context, okno session.Okno,
	polecenie session.Polecenie, zadanie string, argumenty []string) (*przebiegBudowania, error) {

	if a.uruchamiacz == nil {
		return nil, bladWykonaniaDevelopera("serwer nie ma uruchamiacza procesów")
	}
	przebieg := &przebiegBudowania{
		kod:         nowyIdentyfikator(przedrostekBudowania),
		oknoKod:     okno.Id,
		idSesji:     okno.IdSesji,
		zadanie:     zadanie,
		argumenty:   argumenty,
		kontekst:    zKontemZadania(context.Background(), ctx),
		stan:        shared.BuildStatusRunning,
		ogon:        make([]string, 0, 64),
		uruchomiono: time.Now().UTC(),
		koniec:      make(chan struct{}),
	}
	if !a.rejestr.Zajmij(przebieg) {
		return nil, bladZadaniaDevelopera(
			"w oknie " + okno.Id + " trwa już budowanie — przerwij je, zanim uruchomisz następne")
	}

	uchwyt, err := a.uruchamiacz.UruchomProces(ctx, okno, polecenie)
	if err != nil {
		a.rejestr.Zwolnij(przebieg)
		return nil, bladWykonaniaDevelopera(
			"nie można uruchomić zadania " + zadanie + " w " + polecenie.Katalog + ": " + err.Error())
	}
	drzewo, err := session.PrzejmijDrzewo(uchwyt.Pid())
	if err != nil {
		_ = uchwyt.Ubij()
		_ = uchwyt.Czekaj()
		a.rejestr.Zwolnij(przebieg)
		return nil, bladWykonaniaDevelopera("nie można objąć drzewa procesu budowania: " + err.Error())
	}
	przebieg.uchwyt, przebieg.drzewo = uchwyt, drzewo

	a.zapiszPrzebieg(ctx, przebieg)
	// Rozgłoszenie idzie przed pompami logu: wiersz nie może wyprzedzić zdarzenia created.
	a.rozglosBudowanie(ctx, shared.ChangeKindCreated, przebieg, "")

	gotowe := make(chan struct{}, 2)
	go a.pompujLog(przebieg, uchwyt.Wyjscie(), gotowe)
	go a.pompujLog(przebieg, uchwyt.Diagnostyka(), gotowe)
	go a.pilnujBudowania(przebieg, gotowe)
	return przebieg, nil
}

func (a *adapterDevelopera) pompujLog(przebieg *przebiegBudowania, zrodlo io.Reader,
	gotowe chan<- struct{}) {

	defer func() { gotowe <- struct{}{} }()
	if zrodlo == nil {
		return
	}
	czytnik := bufio.NewReaderSize(zrodlo, 8*1024)
	for {
		wiersz, err := czytajWiersz(czytnik)
		if wiersz != "" {
			przebieg.Dopisz(wiersz)
			a.rozglosBudowanie(przebieg.kontekst, shared.ChangeKindUpdated, przebieg, wiersz)
		}
		if err != nil {
			return
		}
	}
}

func czytajWiersz(czytnik *bufio.Reader) (string, error) {
	var budowany strings.Builder
	for {
		czesc, przedluzenie, err := czytnik.ReadLine()
		budowany.Write(czesc)
		if err != nil {
			return budowany.String(), err
		}
		if !przedluzenie || budowany.Len() >= najdluzszyWierszLogu {
			return budowany.String(), nil
		}
	}
}

// Pompy logu kończą przed rozgłoszeniem stanu końcowego: wiersz po zamknięciu przebiegu Build Output odrzuca.
func (a *adapterDevelopera) pilnujBudowania(przebieg *przebiegBudowania, gotowe <-chan struct{}) {
	blad := przebieg.uchwyt.Czekaj()
	<-gotowe
	<-gotowe

	stan, kodWyjscia := wynikBudowania(blad)
	if !przebieg.Domknij(stan, kodWyjscia) {
		stan, kodWyjscia, _, _ = przebieg.Migawka()
	}
	przebieg.Zwolnij()
	a.rejestr.Zwolnij(przebieg)

	if a.repozytorium != nil {
		_ = a.repozytorium.ZakonczPrzebieg(context.Background(), przebieg.kod, stan,
			wskaznikDuzej(kodWyjscia), przebieg.Ogon())
		a.odlozPomiarPrzebiegu(przebieg)
	}
	a.rozglosBudowanie(przebieg.kontekst, shared.ChangeKindUpdated, przebieg, podsumowaniePrzebiegu(stan, kodWyjscia))
}

func wynikBudowania(blad error) (shared.BuildStatus, *int) {
	if blad == nil {
		zero := 0
		return shared.BuildStatusSucceeded, &zero
	}
	var zakonczenie *exec.ExitError
	if errors.As(blad, &zakonczenie) {
		kod := zakonczenie.ExitCode()
		if kod < 0 {
			return shared.BuildStatusStopped, nil
		}
		return shared.BuildStatusFailed, &kod
	}
	return shared.BuildStatusFailed, nil
}

func (p *przebiegBudowania) CzekajNaKoniec(najdluzej time.Duration) {
	select {
	case <-p.koniec:
	case <-time.After(najdluzej):
	}
}

func (a *adapterDevelopera) zapiszPrzebieg(ctx context.Context, przebieg *przebiegBudowania) {
	if a.repozytorium == nil {
		return
	}
	_ = a.repozytorium.ZapiszPrzebieg(ctx, dane.PrzebiegBudowania{
		Kod:       przebieg.kod,
		OknoKod:   przebieg.oknoKod,
		Zadanie:   przebieg.zadanie,
		Argumenty: wskaznikTekstu(strings.Join(przebieg.argumenty, "\n")),
		Stan:      shared.BuildStatusRunning,
	})
}

func (a *adapterDevelopera) rozglosBudowanie(ctx context.Context, zmiana shared.ChangeKind, przebieg *przebiegBudowania,
	wiersz string) {

	if a.przyrost == nil {
		return
	}
	a.przyrost(ctx, zmiana, budowanieKontraktu(przebieg), wiersz)
}
