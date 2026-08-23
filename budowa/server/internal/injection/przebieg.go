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
// fragment prowenancji, uruchamia program, podaje wypowiedź na wejście
// i czyta strumień do końca tury.
//
// Kolejność jest zamierzona — prowenancja idzie PRZED uruchomieniem procesu,
// więc odbiorca zna warunki wywołania nawet wtedy, gdy proces w ogóle nie
// wystartuje.
func wykonajPrzebieg(kontekst context.Context, z Zapytanie, konto Konto, proba int, powod string, na chan<- Fragment) (przebieg, error) {
	// Treść pól --settings i --mcp-config trzeba najpierw zapisać na dysk, bo
	// program `claude` oczekuje tam ścieżek plików, nie napisów JSON. Pliki żyją
	// tylko przez ten przebieg; sprzątamy je po zakończeniu tury.
	zmaterializowane, sprzataj, err := zmaterializujUstawienia(z.Ustawienia)
	if err != nil {
		return przebieg{}, err
	}
	defer sprzataj()

	// argv składamy z materializowanych ustawień (ze ścieżkami), a prowenancję
	// z pierwotnych — pole „settings" ma pokazywać treść przekazaną kanałowi,
	// a wiersz argv realną komendę ze ścieżką pliku.
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
	// Zawiadomienie idzie ZARAZ po starcie, przed podaniem wejścia: gdyby tura
	// padła w połowie, proces i tak jest już objęty uchwytem sesji.
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

// skroc przycina treść diagnostyczną do rozmiaru czytelnego w komunikacie.
func skroc(tresc string) string {
	const limit = 2000
	if len(tresc) <= limit {
		return tresc
	}
	return tresc[:limit] + "…"
}
