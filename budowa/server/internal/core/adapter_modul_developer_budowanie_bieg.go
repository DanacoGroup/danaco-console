// Odpowiedzialność pliku: bieg przebiegu budowania od startu do domknięcia —
// uruchomienie przez port warstwy kanału, przejęcie drzewa potomstwa, pompa logu
// rozsyłająca wiersze zdarzeniem i zapis wyniku.
//
// Obserwator zakończenia pracuje poza żądaniem: rozłączenie klienta w połowie
// kompilacji nie przerywa kompilacji, więc obserwator ma własną gorutynę
// i własny kontekst, a nie kontekst komendy.
//
// Log idzie wierszami, nie blokami. Kontrakt niesie w zdarzeniu pole `logLine`
// — jeden wiersz — więc pompa skleja odczyty w pełne wiersze zamiast rozsyłać
// surowe porcje odczytu. Wiersz przerwany w połowie bufora trafiłby do Build
// Output jako dwa wiersze i rozbiłby rozpoznawanie zgłoszeń kompilatora.
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

// czasNaDomknieciePrzebiegu jest chwilą, przez którą komenda przerwania czeka
// na obserwatora. Krótka z zamysłem: odpowiedź ma być szybka, a stan końcowy
// i tak dojdzie zdarzeniem `developer.build.changed`.
const czasNaDomknieciePrzebiegu = 750 * time.Millisecond

// najdluzszyWierszLogu chroni przed wyjściem bez znaku końca wiersza — paskiem
// postępu, który potrafi rosnąć w nieskończoność.
const najdluzszyWierszLogu = 64 * 1024

// uruchomBudowanie startuje przebieg i podpina pompy logu oraz obserwatora.
func (a *adapterDevelopera) uruchomBudowanie(ctx context.Context, okno session.Okno,
	polecenie session.Polecenie, zadanie string, argumenty []string) (*przebiegBudowania, error) {

	if a.uruchamiacz == nil {
		return nil, bladWykonaniaDevelopera("rdzeń nie ma uruchamiacza procesów")
	}
	przebieg := &przebiegBudowania{
		kod:         nowyIdentyfikator(przedrostekBudowania),
		oknoKod:     okno.Id,
		idSesji:     okno.IdSesji,
		zadanie:     zadanie,
		argumenty:   argumenty,
		stan:        shared.BuildStatusRunning,
		ogon:        make([]string, 0, 64),
		uruchomiono: time.Now().UTC(),
		koniec:      make(chan struct{}),
	}
	if !a.rejestr.Zajmij(przebieg) {
		return nil, bladZadaniaDevelopera(
			"w oknie " + okno.Id + " trwa już budowanie — przerwij je, zanim uruchomisz następne")
	}

	uchwyt, err := a.uruchamiacz.UruchomProces(okno, polecenie)
	if err != nil {
		a.rejestr.Zwolnij(przebieg)
		return nil, bladWykonaniaDevelopera(
			"nie można uruchomić zadania " + zadanie + " w " + polecenie.Katalog + ": " + err.Error())
	}
	drzewo, err := session.PrzejmijDrzewo(uchwyt.Pid())
	if err != nil {
		// Proces już biegnie, a uchwytu drzewa nie ma — zostawienie go tak
		// znaczyłoby sierotę poza rejestrem rdzenia.
		_ = uchwyt.Ubij()
		_ = uchwyt.Czekaj()
		a.rejestr.Zwolnij(przebieg)
		return nil, bladWykonaniaDevelopera("nie można objąć drzewa procesu budowania: " + err.Error())
	}
	przebieg.uchwyt, przebieg.drzewo = uchwyt, drzewo

	a.zapiszPrzebieg(ctx, przebieg)
	// Rozgłoszenie idzie przed pompami logu: Build Output rozpoznaje wiersze po
	// identyfikatorze przebiegu, więc gdyby pierwszy wiersz wyprzedził zdarzenie
	// `created`, okno odrzuciłoby początek logu jako cudzy.
	a.rozglosBudowanie(shared.ChangeKindCreated, przebieg, "")

	gotowe := make(chan struct{}, 2)
	go a.pompujLog(przebieg, uchwyt.Wyjscie(), gotowe)
	go a.pompujLog(przebieg, uchwyt.Diagnostyka(), gotowe)
	go a.pilnujBudowania(przebieg, gotowe)
	return przebieg, nil
}

// pompujLog czyta strumień procesu wierszami i rozsyła je zdarzeniem przyrostu.
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
			a.rozglosBudowanie(shared.ChangeKindUpdated, przebieg, wiersz)
		}
		if err != nil {
			// Koniec potoku jest normalnym końcem odczytu, a błąd odczytu
			// dotyczy tego jednego potoku — przebieg domknie czekający na
			// zakończenie procesu.
			return
		}
	}
}

// czytajWiersz zwraca jeden wiersz bez znaków końca. Wiersz dłuższy od granicy
// zostaje oddany w kawałkach zamiast rosnąć w pamięci bez końca.
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

// pilnujBudowania czeka na zakończenie procesu i domyka przebieg. Czekanie na
// pompy logu jest konieczne: wiersz odczytany po rozgłoszeniu stanu końcowego
// dotarłby do Build Output po zamknięciu przebiegu i zostałby odrzucony.
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
	a.rozglosBudowanie(shared.ChangeKindUpdated, przebieg, podsumowaniePrzebiegu(stan, kodWyjscia))
}

// wynikBudowania przekłada wynik oczekiwania na stan i kod wyjścia. Kod różny
// od zera jest wynikiem budowania, nie usterką rdzenia.
func wynikBudowania(blad error) (shared.BuildStatus, *int) {
	if blad == nil {
		zero := 0
		return shared.BuildStatusSucceeded, &zero
	}
	var zakonczenie *exec.ExitError
	if errors.As(blad, &zakonczenie) {
		kod := zakonczenie.ExitCode()
		if kod < 0 {
			// Kod ujemny znaczy zakończenie sygnałem — budowanie zostało
			// przerwane z zewnątrz, a nie zawiodło na własnym wyniku.
			return shared.BuildStatusStopped, nil
		}
		return shared.BuildStatusFailed, &kod
	}
	return shared.BuildStatusFailed, nil
}

// CzekajNaKoniec czeka na domknięcie przebiegu nie dłużej niż podany czas.
func (p *przebiegBudowania) CzekajNaKoniec(najdluzej time.Duration) {
	select {
	case <-p.koniec:
	case <-time.After(najdluzej):
	}
}

// zapiszPrzebieg odkłada uruchomiony przebieg do dziennika. Nieudany zapis nie
// zatrzymuje budowania, które już biegnie.
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

// rozglosBudowanie oddaje przyrost obsługiwaczowi, który rozsyła
// `developer.build.changed`. Brak podpięcia nie zmienia pracy modułu.
func (a *adapterDevelopera) rozglosBudowanie(zmiana shared.ChangeKind, przebieg *przebiegBudowania,
	wiersz string) {

	if a.przyrost == nil {
		return
	}
	a.przyrost(zmiana, budowanieKontraktu(przebieg), wiersz)
}
