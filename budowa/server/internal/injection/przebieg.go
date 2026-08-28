package injection

import (
	"context"
	"fmt"
)

// przebieg jest wynikiem jednego uruchomienia programu w ramach tury.
// Rotacja kont sprawia, że tura może mieć więcej niż jeden przebieg.
type przebieg struct {
	Obserwacja obserwacja
	Bledy      string
	Pid        int
}

// wykonajPrzebieg przeprowadza jedno wywołanie: składa argumenty, wysyła
// fragment prowenancji, uruchamia program i czyta strumień do końca tury.
func wykonajPrzebieg(kontekst context.Context, z Zapytanie, konto Konto, proba int, powod string, na chan<- Fragment) (przebieg, error) {
	// Treść pól --settings i --mcp-config trzeba najpierw zapisać na dysk
	// jako pliki tymczasowe.
	zmaterializowane, sprzataj, err := zmaterializujUstawienia(z.Ustawienia)
	if err != nil {
		return przebieg{}, err
	}
	defer sprzataj()

	// argv składamy z materializowanych ustawień, a prowenancję z pierwotnych,
	// przed materializacją.
	argv := Argumenty(zmaterializowane, z.Nakladka)
	prowenancja := ZlozProwenancje(z.Ustawienia, z.Nakladka, argv, konto, proba, powod)
	if err := wyslij(kontekst, na, FragmentProwenancji(z.IdOkna, z.IdWiadomosci, prowenancja)); err != nil {
		return przebieg{}, err
	}

	proces, err := Uruchom(kontekst, zmaterializowane, argv, konto.KatalogKonfiguracji)
	if err != nil {
		return przebieg{}, err
	}
	wynik := przebieg{Pid: proces.Pid()}
	// Zawiadomienie idzie zaraz po starcie procesu, przed podaniem wejścia na
	// strumień.
	if z.NaStartProcesu != nil {
		z.NaStartProcesu(wynik.Pid)
	}

	if err := proces.Wyslij(z.Tekst); err != nil {
		wynik.Bledy = proces.Bledy()
		return wynik, err
	}
	if err := proces.ZamknijWejscie(); err != nil {
		wynik.Bledy = proces.Bledy()
		return wynik, err
	}

	obserwowane, bladCzytania := czytajStrumien(kontekst, proces.Wyjscie(), z, na)
	wynik.Obserwacja = obserwowane
	bladProcesu := proces.Czekaj()
	wynik.Bledy = proces.Bledy()

	if bladCzytania != nil {
		return wynik, bladCzytania
	}
	if obserwowane.Tura == nil {
		return wynik, bladBezTury(bladProcesu, wynik.Bledy)
	}
	return wynik, nil
}

// bladBezTury opisuje przebieg zakończony bez zdarzenia kończącego turę.
// Treść wyjścia diagnostycznego jedzie razem z błędem, bo bez niej przyczyna
// zostałaby po stronie procesu i nie doszłaby do Operatora.
func bladBezTury(bladProcesu error, bledy string) error {
	if bladProcesu != nil {
		if bledy != "" {
			return fmt.Errorf("%w; wyjście diagnostyczne: %s", bladProcesu, skroc(bledy))
		}
		return bladProcesu
	}
	if bledy != "" {
		return fmt.Errorf("injection: strumień skończył się bez zdarzenia %s; wyjście diagnostyczne: %s", TypResult, skroc(bledy))
	}
	return fmt.Errorf("injection: strumień skończył się bez zdarzenia %s", TypResult)
}

// skroc przycina treść diagnostyczną do rozmiaru czytelnego w komunikacie
// błędu zwracanego wywołującemu.
func skroc(tresc string) string {
	const limit = 2000
	if len(tresc) <= limit {
		return tresc
	}
	return tresc[:limit] + "…"
}
