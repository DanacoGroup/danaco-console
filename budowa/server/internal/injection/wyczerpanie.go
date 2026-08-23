package injection

import (
	"strings"
	"time"
)

// wzorceWyczerpania rozpoznają wyczerpanie limitu w treści komunikatu. Ścieżka
// pierwsza i pewniejsza jest inna — pole rate_limit_info zdarzenia
// rate_limit_event; wzorce łapią przypadki, w których program powiedział
// o limicie samym tekstem.
var wzorceWyczerpania = []string{
	"usage limit reached",
	"rate limit exceeded",
	"rate_limit_error",
	"quota exceeded",
	"429",
}

// Wyczerpanie opisuje rozpoznane wyczerpanie limitu konta.
type Wyczerpanie struct {
	// Powod jest krótkim opisem źródła rozpoznania — jedzie do dziennika
	// i do fragmentu prowenancji następnej próby.
	Powod string
	// DoChwili jest momentem odnowienia limitu; zero znaczy „nieznany".
	DoChwili time.Time
}

// rozpoznajWyczerpanie ustala, czy tura padła na limicie konta. Bierze pod
// uwagę stan limitu podany przez program, podsumowanie tury oraz wyjście
// diagnostyczne procesu — w tej kolejności wiarygodności.
func rozpoznajWyczerpanie(o obserwacja, bledy string) (Wyczerpanie, bool) {
	if o.Limit != nil && o.Limit.Status == "rejected" {
		return Wyczerpanie{
			Powod:    "limit konta odrzucił wywołanie (" + o.Limit.Rodzaj + ")",
			DoChwili: chwilaOdnowienia(o.Limit.ResetsAt),
		}, true
	}
	if o.Tura != nil && o.Tura.Blad {
		if w, jest := wzorzecWTresci(o.Tura.Tekst, "podsumowanie tury"); jest {
			return w, true
		}
	}
	if w, jest := wzorzecWTresci(bledy, "wyjście diagnostyczne procesu"); jest {
		return w, true
	}
	return Wyczerpanie{}, false
}

func wzorzecWTresci(tresc, zrodlo string) (Wyczerpanie, bool) {
	male := strings.ToLower(tresc)
	for _, wzorzec := range wzorceWyczerpania {
		if strings.Contains(male, wzorzec) {
			return Wyczerpanie{Powod: zrodlo + ": " + wzorzec}, true
		}
	}
	return Wyczerpanie{}, false
}

func chwilaOdnowienia(sekundy int64) time.Time {
	if sekundy <= 0 {
		return time.Time{}
	}
	return time.Unix(sekundy, 0)
}
