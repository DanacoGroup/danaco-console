// Domknięcie stanu podagenta, żeby wiersz nie został w stanie biegnącym,
// którego nikt nie prowadzi; zapis stanu końcowego jest uparty, ponawiany
// kilka razy.
package podagenci

import (
	"context"
	"fmt"
	"time"

	"danacoconsole/server/internal/dane"
)

// WyjasnienieDomkniecia trafia w pole wynik wiersza domkniętego, ale wyłącznie
// gdy pole jest puste, po polsku, bo czyta je Operator.
const WyjasnienieDomkniecia = "Praca nie jest prowadzona przez żaden proces, " +
	"a wiersz stał w stanie biegnącym — stan domknięto, żeby nie kłamał. " +
	"Powołaj podagenta na nowo, jeśli zadanie ma być dokończone."

// ProbyZapisuStanu to liczba podejść do zapisu stanu. Sześć podejść z odstępem
// podwajanym od 20 ms daje łącznie 620 ms czekania między pierwszym a ostatnim —
// rywalizacja o zapis w bazie trwa milisekundy, więc jest to zapas z nawiązką.
const ProbyZapisuStanu = 6

// odstepPierwszejProby otwiera stopniowe podwajanie odstępów między kolejnymi
// próbami zapisu tego stanu.
const odstepPierwszejProby = 20 * time.Millisecond

// StanPodagenta jest wąskim kontraktem zapisu — tyle z repozytorium
// podagentów, ile domknięcie naprawdę pisze.
// `dane.RepozytoriumPodagentow` wypełnia go bez przeróbek.
type StanPodagenta interface {
	UstawStan(ctx context.Context, kod, stan string, wynik *string) error
}

// ZapiszStan zapisuje stan podagenta, ponawiając podejście po niepowodzeniu;
// kontekst zerwany kończy ponawianie od razu.
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

// DomknijNieczynnych przestawia na stan zatrzymany każdy wiersz z wykazu,
// który stoi w stanie niekońcowym, po stwierdzeniu bezczynności przez
// wołającego.
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
// jeszcze nie ma, żeby nie nadpisać pracy już oddanej.
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
