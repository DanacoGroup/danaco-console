// Plik obsługuje powołanie podagentów (`subagent.spawn`): zakładanie wierszy podagentów, zasilenie silnika kolejek ich pracą i puszczenie jej w tle, każdego na własnej kolejce.
package core

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/podagenci"
	"danacoconsole/shared"
)

// granicaPodagentow to górna granica jednego powołania, wprost z kontraktu — żądanie ponad nią przycina się do niej.
const granicaPodagentow = 15

// przedrostekPodagenta znakuje identyfikator nadany przez rdzeń — podagent
// wychodzi kontraktem pod nim, a nie pod kluczem wiersza.
const przedrostekPodagenta = "podagent-"

// rodzajKolejkiPodagenta jest rodzajem kolejki niosącej pracę podagenta. Podagent jest zadaniem w tle pod oknem wykonawcy, więc jego kolejka nosi rodzaj `multitasking`, odróżniający ją od kolejki etapu.
const rodzajKolejkiPodagenta = "multitasking"

// Powolaj zakłada podagentów okna wykonawcy i puszcza ich pracę w tle. Powołuje model, nie okno — żądanie przychodzi tą samą kopertą niezależnie od nadawcy wywołania.
func (a *adapterPodagentow) Powolaj(ctx context.Context,
	z shared.SubagentSpawnRequest) (shared.SubagentSpawnResponse, error) {

	a.drogaNarzedzia.Do(func() { a.zapisz("podagenci: %s", podagenci.ZdanieODrodze()) })
	if z.WindowId == "" || z.Task == "" {
		return shared.SubagentSpawnResponse{}, bladWskazaniaPodagenta(
			"powołanie bez wskazania okna wykonawcy albo bez zadania")
	}
	okno, err := a.okna.PoIdentyfikatorze(ctx, z.WindowId)
	if err != nil {
		return shared.SubagentSpawnResponse{}, bladNieznanegoOknaPodagenta(z.WindowId, err)
	}
	// Zakres eksperta czytany przed powołaniem: bez Subagent Network powołanie nie idzie dalej.
	ile := a.liczbaPowolania(z.Count)
	if a.straz != nil {
		kod := wartoscTekstu(okno.AgentKod)
		if granica, dotyczy := a.straz.GranicaPodagentowEksperta(ctx, kod); dotyczy {
			if granica <= 0 {
				return shared.SubagentSpawnResponse{}, odmowaPodagentowEksperta(kod)
			}
			if ile > granica {
				a.zapisz("subagent.spawn: ekspert %s ma granicę %d podagentów — powołuję %d zamiast %d",
					kod, granica, granica, ile)
				ile = granica
			}
		}
	}
	powolani, err := a.repozytorium.ZalozPodagentow(ctx,
		a.wierszePowolania(ctx, z, okno, ile))
	if err != nil {
		return shared.SubagentSpawnResponse{}, bladPodagentow(err)
	}
	// Oznaczenie prowadzenia idzie przed puszczeniem pracy w tle, inaczej podagent zostałby sierotą.
	if err := a.repozytorium.OznaczProwadzenie(ctx, kodyPodagentow(powolani), a.uruchomienie); err != nil {
		a.zapisz("podagenci: prowadzenie powołania nie zostało oznaczone: %v", err)
	}
	// Rozgłoszenie idzie przed pracą w tle, żeby panel widział kolejność zgodną z faktami.
	a.rozglosPowolanie(powolani)
	for _, podagent := range powolani {
		a.puscWTle(podagent, okno)
	}
	return shared.SubagentSpawnResponse{Subagents: podagenciKontraktu(powolani)}, nil
}

// liczbaPowolania rozstrzyga, ilu podagentów powołać. Brak wskazania i wskazanie poniżej jedynki znaczą jednego; ponad granicę przycina się i zostawia ślad w dzienniku.
func (a *adapterPodagentow) liczbaPowolania(wskazanie *int) int {
	if wskazanie == nil || *wskazanie < 1 {
		return 1
	}
	if *wskazanie > granicaPodagentow {
		a.zapisz("subagent.spawn: żądano %d podagentów, granica kontraktu wynosi %d — powołuję %d",
			*wskazanie, granicaPodagentow, granicaPodagentow)
		return granicaPodagentow
	}
	return *wskazanie
}

// wierszePowolania składa wiersze podagentów jednego powołania.
//
// Nazwa numeruje się, gdy podagentów jest więcej niż jeden: podagenci
// o identycznej nazwie byliby w panelu nie do rozróżnienia. Powołanie pojedyncze
// zostaje przy nazwie podanej.
func (a *adapterPodagentow) wierszePowolania(ctx context.Context,
	z shared.SubagentSpawnRequest, okno dane.Okno, ile int) []dane.Podagent {

	bieg := a.biegOkna(ctx, okno)
	nazwa := wartoscTekstu(z.Name)
	wiersze := make([]dane.Podagent, 0, ile)
	for numer := 1; numer <= ile; numer++ {
		podagent := dane.Podagent{
			Kod:     nowyIdentyfikator(przedrostekPodagenta),
			OknoID:  okno.ID,
			BiegID:  bieg,
			Zadanie: z.Task,
			Stan:    dane.StanPodagentaOczekuje,
		}
		if nazwa != "" {
			pelna := nazwa
			if ile > 1 {
				pelna = fmt.Sprintf("%s %d", nazwa, numer)
			}
			podagent.Nazwa = &pelna
		}
		wiersze = append(wiersze, podagent)
	}
	return wiersze
}

// biegOkna odnajduje bieg orkiestracji, w którym pracuje okno wykonawcy.
// Najpierw pyta o bieg jego koordynatora, potem o bieg prowadzony przez samo
// okno. Brak biegu znaczy podagenta powołanego w zwykłej rozmowie i jest stanem
// zwykłym, nie usterką.
func (a *adapterPodagentow) biegOkna(ctx context.Context, okno dane.Okno) *int64 {
	if a.biegi == nil {
		return nil
	}
	kandydaci := []int64{}
	if okno.OknoKoordynatoraID != nil {
		kandydaci = append(kandydaci, *okno.OknoKoordynatoraID)
	}
	kandydaci = append(kandydaci, okno.ID)
	for _, oknoID := range kandydaci {
		bieg, err := a.biegi.BiegOknaKoordynatora(ctx, oknoID)
		if err == nil {
			id := bieg.ID
			return &id
		}
		if !errors.Is(err, dane.ErrBrakWiersza) {
			a.zapisz("subagent.spawn: nie można odczytać biegu okna %d: %v", oknoID, err)
			return nil
		}
	}
	return nil
}

// puscWTle uruchamia pracę podagenta silnikiem kolejek, nie czekając na wynik.
//
// Silnik niewpięty nie udaje wykonania: podagent zostaje wtedy w stanie
// `pending`, czyli powołany i nierozpoczęty, zamiast oznaczonego jako zrobiony
// bez wykonawcy.
func (a *adapterPodagentow) puscWTle(podagent dane.Podagent, okno dane.Okno) {
	if a.kolejki == nil {
		a.zapisz("subagent.spawn: podagent %s powołany bez wpiętego silnika kolejek — zostaje w stanie %q",
			podagent.Kod, dane.StanPodagentaOczekuje)
		return
	}
	// Nadzór mówi o procesie orkiestratora, nie samej pozycji podagenta.
	if a.nadzor != nil {
		a.zapisz("subagent.spawn: podagent %s pod oknem %s — orkiestrator: %s",
			podagent.Kod, wartoscTekstu(okno.IdentyfikatorZewnetrzny),
			a.nadzor(wartoscTekstu(okno.IdentyfikatorZewnetrzny)).Opis())
	}
	zycie := a.zycie
	if zycie == nil {
		zycie = context.Background()
	}
	// Kontekst własny na podagenta daje `subagent.stop` uchwyt do jednej pracy.
	zycie, odwolaj := context.WithCancel(zycie)
	a.zapamietajPrace(podagent.Kod, odwolaj)
	go func() {
		// Odwołanie zwalnia się zawsze, inaczej wykaz prac rósłby z każdym powołaniem.
		defer odwolaj()
		defer a.zapomnijPrace(podagent.Kod)
		if err := a.wykonaj(zycie, podagent, okno); err != nil {
			a.zapisz("subagent.spawn: praca podagenta %s nie doszła do skutku: %v", podagent.Kod, err)
			komunikat := err.Error()
			if err := a.ustawStanPodagenta(zycie, podagent.Kod,
				dane.StanPodagentaBledny, &komunikat); err != nil {
				a.zapisz("subagent.spawn: nie można zapisać niepowodzenia podagenta %s: %v",
					podagent.Kod, err)
			}
		}
	}()
}

// wykonaj zakłada kolejkę podagenta, wiąże ją z jego wierszem i prowadzi pracę przez silnik, a po jej zakończeniu przepisuje stan pozycji na stan podagenta. Pozycja i wiązanie idą przed wejściem w stan `running`.
func (a *adapterPodagentow) wykonaj(ctx context.Context, podagent dane.Podagent, okno dane.Okno) error {
	kolejkaID, err := a.kolejki.repozytorium.UtworzKolejke(ctx, dane.Kolejka{
		Nazwa: podagent.Kod, Rodzaj: rodzajKolejkiPodagenta, Stan: shared.QueueStatusIdle,
		SesjaID: &okno.SesjaID, OknoKoordynatoraID: &okno.ID,
	})
	if err != nil {
		return err
	}
	pozycjaID, err := a.kolejki.repozytorium.DodajPozycje(ctx, dane.Pozycja{
		KolejkaID: kolejkaID, OknoWykonawcyID: &okno.ID,
		Tytul: tytulPodagenta(podagent), TrescZlecenia: &podagent.Zadanie,
		Stan: stanPozycjiOczekuje,
	})
	if err != nil {
		return err
	}
	if err := a.repozytorium.PrzypiszPozycje(ctx, podagent.Kod, pozycjaID); err != nil {
		return err
	}
	if err := a.ustawStanPodagenta(ctx, podagent.Kod, dane.StanPodagentaWBiegu, nil); err != nil {
		return err
	}
	if _, err := a.kolejki.Wykonaj(ctx, shared.QueueActionRequest{
		QueueId: strconv.FormatInt(kolejkaID, 10), Action: shared.QueueActionStart,
	}); err != nil {
		return err
	}
	return a.przepiszStanPozycji(ctx, podagent.Kod, kolejkaID, pozycjaID)
}

// przepiszStanPozycji przenosi stan pozycji kolejki na stan podagenta. Automat stanu jest jeden: cyklem życia zlecenia rządzi silnik kolejek, a wiersz podagenta jest jego odwzorowaniem dla panelu.
func (a *adapterPodagentow) przepiszStanPozycji(ctx context.Context,
	kod string, kolejkaID, pozycjaID int64) error {

	pozycje, err := a.kolejki.repozytorium.ListaPozycji(ctx, kolejkaID)
	if err != nil {
		return err
	}
	for _, pozycja := range pozycje {
		if pozycja.ID != pozycjaID {
			continue
		}
		return a.ustawStanPodagenta(ctx, kod, stanPodagentaZPozycji(pozycja.Stan), nil)
	}
	return nil
}

// tytulPodagenta nazywa pozycję kolejki w Queue Managerze. Nazwa własna, gdy
// jest; identyfikator, gdy jej nie ma — pozycja bez tytułu byłaby w wykazie
// nie do wskazania.
func tytulPodagenta(podagent dane.Podagent) string {
	if podagent.Nazwa != nil && *podagent.Nazwa != "" {
		return *podagent.Nazwa
	}
	return podagent.Kod
}
