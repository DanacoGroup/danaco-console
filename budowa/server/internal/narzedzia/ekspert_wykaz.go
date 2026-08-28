// Plik składa wykaz narzędzi z podzbioru wskazanego przez eksperta, kodem
// nazywającym narzędzie albo grupę, z zawężeniem typu fail-open przy braku
// rozpoznania.
package narzedzia

import (
	"fmt"
	"strings"

	"danacoconsole/server/internal/injection"
)

// WykazEksperta jest wynikiem doboru: pozycje, kody nierozpoznane oraz prawda
// o tym, czy zawężenie w ogóle weszło i dlaczego.
type WykazEksperta struct {
	// Pozycje są wykazem podawanym modelowi.
	Pozycje []Narzedzie
	// Nierozpoznane niesie kody eksperta bez dopasowania do narzędzia ani grupy.
	Nierozpoznane []string
	// Zawezony mówi, czy wykaz jest podzbiorem wskazanym przez eksperta, czy
	// wykazem okna w całości.
	Zawezony bool
	// Powod jest zdaniem tłumaczącym brak zawężenia; pusty przy zawężeniu udanym.
	Powod string
}

// ZlozWykazEksperta zawęża wykaz okna do pozycji wskazanych przez definicję,
// zachowując kolejność źródła, nigdy nie dokładając pozycji spoza niego.
func ZlozWykazEksperta(definicja DefinicjaEksperta, zrodlo []Narzedzie) WykazEksperta {
	kody := definicja.Kody()
	if len(kody) == 0 {
		return WykazEksperta{
			Pozycje:  zrodlo,
			Zawezony: false,
			Powod: fmt.Sprintf("ekspert %q nie wskazał ani jednego kodu w polach skillIds i connectorIds"+
				" — nie ma czym zawęzić wykazu, więc idzie wykaz okna w całości", definicja.Kod),
		}
	}
	wybrane, nierozpoznane := przesiej(kody, zrodlo)
	if len(wybrane) == 0 {
		return WykazEksperta{
			Pozycje:       zrodlo,
			Nierozpoznane: nierozpoznane,
			Zawezony:      false,
			Powod: fmt.Sprintf("ekspert %q wskazał %d kodów i ŻADEN nie nazywa narzędzia ani grupy"+
				" — zawężenia nie da się wyprowadzić, więc idzie wykaz okna w całości", definicja.Kod, len(kody)),
		}
	}
	return WykazEksperta{Pozycje: wybrane, Nierozpoznane: nierozpoznane, Zawezony: true}
}

// przesiej wybiera pozycje wskazane kodami i oddaje kody nierozpoznane, jednym
// przejściem po źródle, z zachowaniem jego kolejności.
func przesiej(kody []string, zrodlo []Narzedzie) (wybrane []Narzedzie, nierozpoznane []string) {
	nazwy := make(map[string]bool, len(zrodlo))
	grupy := make(map[string]bool, len(zrodlo))
	for _, pozycja := range zrodlo {
		nazwy[pozycja.Nazwa] = true
		if pozycja.Grupa != "" {
			grupy[pozycja.Grupa] = true
		}
	}
	wskazaneNazwy := make(map[string]bool, len(kody))
	wskazaneGrupy := make(map[string]bool, len(kody))
	nierozpoznane = make([]string, 0, len(kody))
	for _, kod := range kody {
		switch {
		case nazwy[kod]:
			wskazaneNazwy[kod] = true
		case grupy[kod]:
			wskazaneGrupy[kod] = true
		default:
			nierozpoznane = append(nierozpoznane, kod)
		}
	}
	wybrane = make([]Narzedzie, 0, len(zrodlo))
	for _, pozycja := range zrodlo {
		if wskazaneNazwy[pozycja.Nazwa] || wskazaneGrupy[pozycja.Grupa] {
			wybrane = append(wybrane, pozycja)
		}
	}
	return wybrane, nierozpoznane
}

// RozbijDolozenia czyta wartość przełącznika dołożeń narzędzi na nazwy,
// odsiewając nazwy puste i powtórzone, w kolejności dokładania.
func RozbijDolozenia(wartosc string) []string {
	widziane := map[string]bool{}
	nazwy := []string{}
	for _, nazwa := range strings.Split(wartosc, injection.RozdzielnikDolozen) {
		nazwa = strings.TrimSpace(nazwa)
		if nazwa == "" || widziane[nazwa] {
			continue
		}
		widziane[nazwa] = true
		nazwy = append(nazwy, nazwa)
	}
	return nazwy
}

// ZDolozeniami dokłada do wykazu pozycje wskazane doraźnie w sesji,
// rozpoznawane tymi samymi drogami co dobór eksperta, i dopisuje pozycje
// nierozpoznane.
func ZDolozeniami(wynik WykazEksperta, nazwy []string, zrodlo []Narzedzie) WykazEksperta {
	if len(nazwy) == 0 {
		return wynik
	}
	obecne := make(map[string]bool, len(wynik.Pozycje))
	for _, pozycja := range wynik.Pozycje {
		obecne[pozycja.Nazwa] = true
	}
	dolozone, nierozpoznane := przesiej(nazwy, zrodlo)
	for _, pozycja := range dolozone {
		if obecne[pozycja.Nazwa] {
			continue
		}
		obecne[pozycja.Nazwa] = true
		wynik.Pozycje = append(wynik.Pozycje, pozycja)
	}
	wynik.Nierozpoznane = append(wynik.Nierozpoznane, nierozpoznane...)
	return wynik
}

// WykazBezEksperta składa wynik dla przypadku, w którym definicji nie ma:
// rdzeń nie odpowiedział, odmówił albo nie zna kodu.
func WykazBezEksperta(kod string, przyczyna error, zrodlo []Narzedzie) WykazEksperta {
	return WykazEksperta{
		Pozycje:  zrodlo,
		Zawezony: false,
		Powod: fmt.Sprintf("definicji eksperta %q nie udało się odczytać (%v)"+
			" — zawężenia nie ma na czym oprzeć, więc idzie wykaz okna w całości", kod, przyczyna),
	}
}
