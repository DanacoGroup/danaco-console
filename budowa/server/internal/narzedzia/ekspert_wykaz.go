// Odpowiedzialność pliku: złożenie wykazu z podzbioru wskazanego przez eksperta.
//
// ── JAK CZYTA SIĘ KOD WSKAZANY PRZEZ EKSPERTA ───────────────────────────────
// `Agent.SkillIds` i `Agent.ConnectorIds` są w kontrakcie listami napisów bez
// narzuconego słownika — Operator wpisuje kod ręcznie (`client/src/moduly/
// agents/okno-skills-manager.ts`). Kod rozpoznaje się więc dwiema drogami, obie
// sprawdzalne wprost wobec kontraktu i żadna nie wymagająca wykazu własnego:
//
//	(a) nazwa narzędzia — `danaco_agent_list` wskazuje jedną pozycję wykazu;
//	(b) nazwa grupy — `agent` wskazuje wszystkie pozycje tego obszaru (`grupa.go`).
//
// Grupa jest tu jednostką doboru: narzędzia są pogrupowane „do czego służą"
// właśnie po to, żeby dobór dawał się wykonać jednym wskazaniem zamiast dwustu.
//
// ── KOD NIEROZPOZNANY NIE JEST POŁYKANY ─────────────────────────────────────
// Kod, którego nie da się rozpoznać żadną z dwóch dróg, nie jest połykany po
// cichu: wraca osobną listą `Nierozpoznane` i idzie stamtąd do dziennika oraz
// do licznika (`licznik.go`).
//
// ── KIEDY ZAWĘŻENIE W OGÓLE WCHODZI ─────────────────────────────────────────
// Zawężenie wchodzi wtedy i tylko wtedy, gdy jest czym zawęzić — czyli gdy
// rozpoznano co najmniej jeden kod. Każdy inny przypadek (rdzeń milczy, kodu nie
// ma, żadnego kodu nie rozpoznano) zostawia wykaz okna i melduje powód: pole
// `Powod` niesie zdanie, `Zawezony` niesie fakt, a `main.go` i rozdzielnia kładą
// je na wyjście diagnostyczne.
//
// Wykaz okna (fail-open) wybrano zamiast wykazu pustego: usterka jednej drogi
// nie zabiera modelowi wszystkich narzędzi na całą turę, a okno z ekspertem,
// którego rdzeń chwilowo nie potwierdził, pracuje dalej. Pełny wykaz kosztuje
// jednak trzydzieści kilka tysięcy żetonów, dlatego koszt jest mierzony
// i meldowany przy każdym takim przypadku. Wykaz pusty odpada, bo model bez ani
// jednego narzędzia nie mówi „nie znam eksperta", tylko po prostu nie działa.
// Rozstrzyga to jedna funkcja — `ZlozWykazEksperta` niżej.
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
	// Nierozpoznane niesie kody eksperta, które nie nazwały ani narzędzia, ani
	// grupy. Kolejność jest kolejnością definicji eksperta.
	Nierozpoznane []string
	// Zawezony mówi, czy wykaz jest podzbiorem wskazanym przez eksperta (prawda),
	// czy wykazem okna oddanym w całości (fałsz).
	Zawezony bool
	// Powod jest zdaniem tłumaczącym brak zawężenia. Pusty przy zawężeniu — nie
	// ma czego tłumaczyć, gdy dobór się udał.
	Powod string
}

// ZlozWykazEksperta zawęża wykaz okna do pozycji wskazanych przez definicję.
//
// `zrodlo` jest wykazem, który dostałoby to okno bez eksperta — funkcja go
// wyłącznie przesiewa i nigdy nie dokłada pozycji spoza niego. Zasięg eksperta
// nie może dać modelowi ani jednego narzędzia, którego okno by nie dostało:
// dobór jest zawężeniem, nie nadaniem.
//
// Kolejność zostaje kolejnością źródła, nie kolejnością kodów eksperta: ten sam
// kontrakt i ten sam ekspert mają dawać ten sam wykaz, a kody bywają przestawiane
// w bibliotece bez znaczenia dla doboru.
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

// przesiej wybiera pozycje wskazane kodami i oddaje kody nierozpoznane.
//
// Jedno przejście po źródle na koniec, a nie sklejanie wyników kod po kodzie:
// dwa kody wskazujące tę samą pozycję (nazwa narzędzia i jego grupa naraz) mają
// dać ją raz, a kolejność ma zostać kolejnością źródła.
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

// RozbijDolozenia czyta wartość przełącznika `--dolozenia` na nazwy.
//
// Nazwa przełącznika i rozdzielnik pochodzą z `injection`, a nie z literałów
// zapisanych tutaj. Wiersz uruchomienia składa strona rdzenia
// (`injection/zestaw_narzedzi.go`), a czyta go ta; `flag.Parse` idzie na
// domyślnym `flag.CommandLine`, czyli z `ExitOnError`, więc rozjazd nazwy
// o jeden znak nie dałby „tury bez kilku pozycji", tylko ubity proces serwera
// narzędzi i turę bez całego wykazu. Dwie kopie tej nazwy byłyby więc drugą
// prawdą z ceną nieproporcjonalną do zysku.
//
// Odsiew jest ten sam, co po stronie składającej: nazwy puste odpadają, nazwa
// powtórzona zostaje na pierwszej pozycji. Kolejność jest kolejnością dokładania.
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

// ZDolozeniami dokłada do wykazu pozycje wskazane doraźnie w sesji.
//
// Dołożenie dokłada i nigdy nie zawęża — Operator dorzuca narzędzie do pracy,
// którą prowadzi, poza definicją eksperta i bez jej ruszania. Rozpoznanie nazwy
// idzie tymi samymi dwiema drogami, co dobór eksperta (nazwa narzędzia albo
// nazwa grupy), bo dwa różne rozpoznania tej samej nazwy byłyby dwiema prawdami.
//
// Pozycje dołożone idą na koniec listy, tak jak rozszerzenie roli okna
// (`wykaz.go`): kolejność wykazu zostaje kolejnością kontraktu, a dołożenie ma
// być widoczne jako dołożenie.
//
// Nazwa nierozpoznana dopisuje się do `Nierozpoznane` także wtedy, gdy wykaz nie
// był zawężony i dokładać nie było czego. To nie jest nadgorliwość: nazwa, która
// nie jest ani narzędziem, ani grupą, nie zadziała nigdy — ani dziś, ani pod
// ekspertem — a Operator widziałby ją na wykazie dołożeń sesji i nie widziałby
// jej w turze.
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
//
// Osobna funkcja, a nie `WykazEksperta{}` składane w miejscu wywołania, bo to
// jest przypadek, o którym najłatwiej zapomnieć powiedzieć — a milczenie o nim
// jest dokładnie tą cichą degradacją, przed którą cały ten zasięg powstał.
// Powód wchodzi tu zawsze, więc nie da się złożyć niewiedzy bez jej opisania.
func WykazBezEksperta(kod string, przyczyna error, zrodlo []Narzedzie) WykazEksperta {
	return WykazEksperta{
		Pozycje:  zrodlo,
		Zawezony: false,
		Powod: fmt.Sprintf("definicji eksperta %q nie udało się odczytać (%v)"+
			" — zawężenia nie ma na czym oprzeć, więc idzie wykaz okna w całości", kod, przyczyna),
	}
}
