// Odpowiedzialność pliku: rozpoznanie pułapu kosztu — odróżnienie tury
// wstrzymanej na nastawie Operatora od tury, która padła.
//
// Obie sytuacje wymagają od Operatora innego działania:
//
//	PUŁAP   — tura zatrzymała się na kwocie ustawionej w oknie. Konto jest
//	          sprawne, limit dostawcy nietknięty, kanał czynny; naprawą jest
//	          podniesienie pułapu albo zawężenie zadania.
//	AWARIA  — tura padła: kanał odmówił, proces zginął, strumień się urwał.
//	          Powtórzenie bez zmiany warunków zwykle daje ten sam skutek.
//
// Rozpoznanie stoi obok `wyczerpanie.go`, a nie w nim. Wyczerpanie jest granicą
// dostawcy (limit konta, HTTP 429, `rate_limit_event`) i uruchamia rotację
// kont. Pułap jest granicą nastawy okna i rotacji uruchamiać nie może: kolejne
// konto wydałoby tę samą kwotę, której nastawa zabrania.
//
// Wzorce rozpoznania są zachowawcze. Rozpoznanie nietrafione zostawia turę
// w drodze awarii, a zgłoszenie trafia do wykazu nierozstrzygniętych meldunku
// pakietu.
package injection

import (
	"context"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// zakonczPulapem domyka strumień odmową wstrzymania na pułapie kosztu. Odmowa
// niesie trzy rzeczy: co zostało przerwane, jaka kwota została wydana wobec
// pułapu i którą nastawą Operator to zmieni (`pulap_kosztu_usd`, Ustawienia →
// Modele → Pułap kosztu okna).
//
// Kod odmowy to `permission_denied`, nie `rate_limited`: kod `rate_limited`
// niesie w kontrakcie ponawialność, a ponowienie tury na tym samym pułapie da
// ten sam wynik. Kanał, konto i sesja zostają czynne — wstrzymana jest jedna
// tura.
func zakonczPulapem(kontekst context.Context, na chan<- Fragment, z Zapytanie, p pulap) {
	tresc := fmt.Sprintf(
		"Tura wstrzymana na pułapie kosztu — nie dokończono jej, bo koszt dobił "+
			"do granicy ustawionej przez Operatora (%s). Wydano %.2f USD przy pułapie "+
			"%.2f USD. To NIE jest awaria: kanał, konto i sesja pozostają czynne. "+
			"Zmień to nastawą „Pułap kosztu okna\" (klucz pulap_kosztu_usd) — podnieś "+
			"kwotę albo wpisz 0, żeby znieść pułap; możesz też podzielić zadanie na "+
			"mniejsze tury.",
		p.Powod, p.Wydano, z.Ustawienia.PulapKosztuUSD)
	zakonczBledem(kontekst, na, z, shared.ErrorCodePermissionDenied, tresc)
}

// wzorcePulapu rozpoznają przerwanie na pułapie kosztu w treści komunikatu.
//
// Wzorce są angielskie, bo pochodzą z komunikatów programu zewnętrznego.
// Każdy zawiera słowo „budget", które odróżnia tę granicę od granicy dostawcy
// („usage limit", „rate limit" — zob. wyczerpanie.go). Sam wzorzec „limit"
// złapałby wyczerpanie i zawrócił turę z rotacji kont.
var wzorcePulapu = []string{
	"max budget",
	"budget exceeded",
	"budget limit",
	"max_budget",
	"budget_exceeded",
}

// podtypPulapu jest podtypem zdarzenia `result`, którym program zewnętrzny
// melduje wyczerpanie budżetu wywołania. Podtyp jest rozpoznaniem pewniejszym
// niż tekst komunikatu, dlatego sprawdza się go jako pierwszy.
const podtypPulapu = "error_max_budget"

// pulap opisuje turę wstrzymaną na pułapie kosztu.
type pulap struct {
	// Powod nazywa źródło rozpoznania — jedzie do dziennika i do odmowy.
	Powod string
	// Wydano jest kosztem tury odczytanym ze zdarzenia `result` (`costUsd`).
	// Zero znaczy „program kosztu nie podał", a nie „nic nie kosztowało".
	Wydano float64
}

// rozpoznajPulap ustala, czy tura stanęła na pułapie kosztu.
//
// Źródła sprawdzane są w kolejności wiarygodności, tej samej co
// w rozpoznajWyczerpanie: podtyp zdarzenia kończącego turę, potem jego tekst,
// na końcu wyjście diagnostyczne procesu.
//
// Tura bez zdarzenia `result` nie jest pułapem: pułap przerywa turę wewnątrz
// programu, więc program zdąża zgłosić przerwanie. Strumień urwany bez
// zdarzenia kończącego jest awarią.
func rozpoznajPulap(o obserwacja, bledy string) (pulap, bool) {
	if o.Tura == nil {
		return pulap{}, false
	}
	if o.Tura.Podtyp == podtypPulapu {
		return pulap{Powod: "podtyp zdarzenia kończącego turę: " + podtypPulapu,
			Wydano: o.Tura.Koszt}, true
	}
	if !o.Tura.Blad {
		// Tura zakończona powodzeniem nie stanęła na pułapie, choćby jej treść
		// o budżecie wspominała: o stanie tury rozstrzyga zdarzenie kończące,
		// a nie tekst wygenerowany przez model.
		return pulap{}, false
	}
	if wzorzec, jest := wzorzecPulapu(o.Tura.Tekst); jest {
		return pulap{Powod: "podsumowanie tury: " + wzorzec, Wydano: o.Tura.Koszt}, true
	}
	if wzorzec, jest := wzorzecPulapu(bledy); jest {
		return pulap{Powod: "wyjście diagnostyczne procesu: " + wzorzec,
			Wydano: o.Tura.Koszt}, true
	}
	return pulap{}, false
}

// wzorzecPulapu szuka w treści któregoś ze wzorców pułapu.
func wzorzecPulapu(tresc string) (string, bool) {
	male := strings.ToLower(tresc)
	for _, wzorzec := range wzorcePulapu {
		if strings.Contains(male, wzorzec) {
			return wzorzec, true
		}
	}
	return "", false
}
