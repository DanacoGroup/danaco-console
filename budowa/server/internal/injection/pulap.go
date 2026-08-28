// Repozytorium rozpoznaje pułap kosztu, odróżniając turę wstrzymaną na
// nastawie Operatora od tury, która padła.
package injection

import (
	"context"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// zakonczPulapem domyka strumień odmową wstrzymania na pułapie kosztu,
// niosącą przerwane zadanie, wydaną kwotę i nastawę do zmiany.
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

// wzorcePulapu rozpoznają przerwanie na pułapie kosztu w treści komunikatu,
// wedle wzorców angielskich zapożyczonych z programu zewnętrznego.
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

// pulap opisuje turę wstrzymaną na pułapie kosztu okna, wraz z powodem
// rozpoznania i wydaną dotąd kwotą.
type pulap struct {
	// Powod nazywa źródło rozpoznania — jedzie do dziennika i do odmowy.
	Powod string
	// Wydano jest kosztem tury odczytanym ze zdarzenia result; zero znaczy
	// brak podanego kosztu.
	Wydano float64
}

// rozpoznajPulap ustala, czy tura stanęła na pułapie kosztu, sprawdzając
// kolejno podtyp zdarzenia, jego tekst i wyjście diagnostyczne.
func rozpoznajPulap(o obserwacja, bledy string) (pulap, bool) {
	if o.Tura == nil {
		return pulap{}, false
	}
	if o.Tura.Podtyp == podtypPulapu {
		return pulap{Powod: "podtyp zdarzenia kończącego turę: " + podtypPulapu,
			Wydano: o.Tura.Koszt}, true
	}
	if !o.Tura.Blad {
		// Tura zakończona powodzeniem nie stanęła na pułapie, nawet gdy jej
		// treść o budżecie wspomina.
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

// wzorzecPulapu szuka w treści komunikatu jednego ze zdefiniowanych
// wcześniej wzorców pułapu, zwracając go.
func wzorzecPulapu(tresc string) (string, bool) {
	male := strings.ToLower(tresc)
	for _, wzorzec := range wzorcePulapu {
		if strings.Contains(male, wzorzec) {
			return wzorzec, true
		}
	}
	return "", false
}
