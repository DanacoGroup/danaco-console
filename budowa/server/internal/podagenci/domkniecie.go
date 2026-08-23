// Domknięcie stanu podagenta — żeby wiersz nie został w stanie biegnącym,
// którego nikt nie prowadzi.
//
// Praca podagenta biegnie goroutine, a jej odwołanie leży w wykazie prac pod
// kodem podagenta. Kiedy praca się kończy — także gdy kończy błędem — wykaz
// zwalnia się przez `defer`. Zapis stanu końcowego idzie po tym i sam też
// potrafi paść, choćby na `database is locked`. Zostaje wtedy wiersz w stanie
// niekońcowym bez uchwytu do pracy, a `subagent.stop` odwołuje wyłącznie po
// uchwycie: uchwytu nie ma, więc melduje `notRunning` i wiersza nie rusza.
// Podagent stoi w `pending`/`running` do końca życia bazy.
//
// Odpowiedzią jest domknięcie, a nie zapisywanie `stopped` zawsze, bo
// `notRunning` niesie potrzebną wiedzę: Operator zatrzymujący wielu podagentów
// ma wiedzieć, których zdążył zatrzymać, a którzy skończyli sami. Domknięcie
// robi rzecz węższą — bierze wyłącznie wiersze niekońcowe, których pracy nikt
// nie prowadzi, i przestawia je na stan końcowy. Wiersz zakończony zostaje
// nietknięty, a jego kod nadal wraca w `notRunning`.
//
// Zapis jest uparty, nie jednokrotny: nieudany wraca tu jeszcze kilka razy
// z rosnącym odstępem. Rywalizacja o zapis w SQLite jest chwilowa, a jedna
// nieudana próba zamieniłaby ją w stan trwale nieprawdziwy.
package podagenci

import (
	"context"
	"fmt"
	"time"

	"danacoconsole/server/internal/dane"
)

// WyjasnienieDomkniecia trafia w pole `wynik` wiersza domkniętego, ale wyłącznie
// gdy pole jest puste — warunek stawia `wyjasnienieGdyPusto`, nie zapytanie.
// Praca oddana przed urwaniem jest ważniejsza niż wyjaśnienie, dlaczego się
// urwała. Zdanie jest po polsku, bo czyta je Operator w panelu zadań w tle.
const WyjasnienieDomkniecia = "Praca nie jest prowadzona przez żaden proces, " +
	"a wiersz stał w stanie biegnącym — stan domknięto, żeby nie kłamał. " +
	"Powołaj podagenta na nowo, jeśli zadanie ma być dokończone."

// ProbyZapisuStanu to liczba podejść do zapisu stanu. Sześć podejść z odstępem
// podwajanym od 20 ms daje łącznie 620 ms czekania między pierwszym a ostatnim —
// rywalizacja o zapis w bazie trwa milisekundy, więc jest to zapas z nawiązką.
const ProbyZapisuStanu = 6

// odstepPierwszejProby otwiera podwajanie odstępów.
const odstepPierwszejProby = 20 * time.Millisecond

// StanPodagenta jest wąskim kontraktem zapisu — tyle z repozytorium
// podagentów, ile domknięcie naprawdę pisze.
// `dane.RepozytoriumPodagentow` wypełnia go bez przeróbek.
type StanPodagenta interface {
	UstawStan(ctx context.Context, kod, stan string, wynik *string) error
}

// ZapiszStan zapisuje stan podagenta, ponawiając podejście po niepowodzeniu.
//
// Kontekst zerwany kończy ponawianie od razu: zatrzymany rdzeń nie ma po co
// czekać na bazę, a zapis pod kontekstem odwołanym i tak nie ma prawa się udać.
// Dlatego wołający podaje kontekst żywy — życie rdzenia — a nie kontekst
// odwołanej pracy.
//
// Trwałość pusta kończy się błędem nazywającym jej brak, a nie meldunkiem
// udanego zapisu: nie ma gdzie zapisać, więc nie ma czego ponawiać.
func ZapiszStan(ctx context.Context, trwalosc StanPodagenta,
	kod, stan string, wynik *string) error {

	if trwalosc == nil {
		return fmt.Errorf("podagenci: zapis stanu podagenta %q bez trwałości", kod)
	}
	if kod == "" || stan == "" {
		return fmt.Errorf("podagenci: zapis stanu bez podagenta albo bez stanu")
	}
	odstep := odstepPierwszejProby
	var ostatni error
	for proba := 1; proba <= ProbyZapisuStanu; proba++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("podagenci: zapis stanu podagenta %q przerwany: %w", kod, err)
		}
		ostatni = trwalosc.UstawStan(ctx, kod, stan, wynik)
		if ostatni == nil {
			return nil
		}
		if proba == ProbyZapisuStanu {
			break
		}
		zegar := time.NewTimer(odstep)
		select {
		case <-ctx.Done():
			zegar.Stop()
			return fmt.Errorf("podagenci: zapis stanu podagenta %q przerwany: %w", kod, ctx.Err())
		case <-zegar.C:
		}
		odstep *= 2
	}
	return fmt.Errorf("podagenci: stan %q podagenta %q nie zapisał się w %d podejściach: %w",
		stan, kod, ProbyZapisuStanu, ostatni)
}

// Domkniecie oddaje wynik przeglądu wierszy — co domknięto, a co było już
// zamknięte. Rozdział jest potrzebny wołającemu: pierwsze jest zmianą, którą
// trzeba rozgłosić, drugie samym potwierdzeniem stanu.
type Domkniecie struct {
	// Domkniete to kody wierszy przestawionych na stan końcowy tym przeglądem.
	Domkniete []string
	// JuzKoncowe to kody wierszy, które stanu końcowego już się trzymały.
	JuzKoncowe []string
}

// DomknijNieczynnych przestawia na `stopped` każdy wiersz z wykazu, który stoi
// w stanie niekońcowym — bo wołający właśnie stwierdził, że pracy tych
// podagentów nikt nie prowadzi.
//
// Woła się po stwierdzeniu bezczynności, nie zamiast niego. Wykaz podaje
// wołający i to on odpowiada za to, że praca tych podagentów naprawdę nie
// biegnie: przy `subagent.stop` są to kody, dla których odwołanie pracy nie
// miało czego odwołać. Sam przegląd żywotności nie mierzy.
//
// Błąd jednego wiersza nie przerywa pozostałych — domknięcie części wykazu jest
// lepsze niż odmowa domknięcia całości. Pierwszy napotkany błąd wraca po
// przejściu całego wykazu.
func DomknijNieczynnych(ctx context.Context, trwalosc StanPodagenta,
	wiersze []dane.Podagent) (Domkniecie, error) {

	wynik := Domkniecie{Domkniete: []string{}, JuzKoncowe: []string{}}
	var pierwszyBlad error
	for _, podagent := range wiersze {
		if StanKoncowy(podagent.Stan) {
			wynik.JuzKoncowe = append(wynik.JuzKoncowe, podagent.Kod)
			continue
		}
		if err := ZapiszStan(ctx, trwalosc, podagent.Kod,
			dane.StanPodagentaZatrzymany, wyjasnienieGdyPusto(podagent)); err != nil {
			if pierwszyBlad == nil {
				pierwszyBlad = err
			}
			continue
		}
		wynik.Domkniete = append(wynik.Domkniete, podagent.Kod)
	}
	return wynik, pierwszyBlad
}

// wyjasnienieGdyPusto oddaje wyjaśnienie wyłącznie dla wiersza, który wyniku
// jeszcze nie ma.
//
// `ustawStanPodagenta` (`dane/orkiestracja_podagenci.go`) pisze
// `wynik = COALESCE(?, wynik)`, a COALESCE bierze pierwszą wartość niepustą —
// wyjaśnienie podane niepuste nadpisałoby więc pracę, którą podagent zdążył
// oddać. Skoro zapytanie warunku nie stawia, stawia go wołający: pustkę
// rozstrzyga wiersz, który już jest w ręku, więc nie kosztuje to ani jednego
// odczytu więcej.
func wyjasnienieGdyPusto(podagent dane.Podagent) *string {
	if podagent.Wynik != nil && *podagent.Wynik != "" {
		return nil
	}
	wyjasnienie := WyjasnienieDomkniecia
	return &wyjasnienie
}

// ZdanieODomknieciu składa meldunek przeglądu do dziennika rdzenia. Zero
// domkniętych też jest meldunkiem — cisza nie odróżnia „nie było czego
// domykać" od „nie zawołano przeglądu".
func ZdanieODomknieciu(wynik Domkniecie) string {
	if len(wynik.Domkniete) == 0 {
		return fmt.Sprintf("domknięcie stanu: wierszy kłamiących brak "+
			"(sprawdzono %d zakończonych)", len(wynik.JuzKoncowe))
	}
	return fmt.Sprintf("domknięcie stanu: %d podagentów stało w stanie biegnącym "+
		"bez prowadzonej pracy — przestawiono na %q (%v)",
		len(wynik.Domkniete), dane.StanPodagentaZatrzymany, wynik.Domkniete)
}
