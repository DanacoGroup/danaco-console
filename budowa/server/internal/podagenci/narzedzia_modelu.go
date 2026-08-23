// Pakiet podagenci mierzy drogę narzędzia modelu dla rodziny `subagent.*`
// oraz ocenia żywotność procesu z rejestru procesów sesji.
//
// Odpowiedzialność tego pliku: droga narzędzia modelu — czy powołanie podagenta
// jest dziś dostępne modelowi jako wywołanie narzędzia w trakcie tury.
//
// Jak narzędzia docierają do rdzenia. Wykaz narzędzi
// modelu powstaje WYŁĄCZNIE z sekcji `narzedzia.pozycje` kontraktu
// (`shared/contract.json`), z której generator (`shared/gen/narzedzia.mjs`)
// wytwarza `shared.NarzedziaModelu()` i `shared.KomendyNarzedzi`. Serwer
// narzędzi (`server/internal/narzedzia`, binarium `danaco-narzedzia`) czyta ten
// jeden wykaz i mówi wprost: „dopisanie komendy do sekcji narzedzia kontraktu
// powiększa ten serwer bez zmiany choćby jednej linii kodu". Wpis `danaco`
// w konfiguracji MCP okna składa `core/most_narzedzi.go` — osobno dla każdego
// okna, z jego identyfikatorem.
//
// Wynikają z tego dwie rzeczy:
//
//  1. NIE MA drugiego miejsca rejestracji narzędzia. Jedyną drogą jest pozycja
//     w `narzedzia.pozycje` kontraktu — a `contract.*` jest plikiem zakazanym
//     dla tego pakietu, więc dopisanie pozycji idzie zgłoszeniem wpięcia
//     z kotwicą (PozycjeWpiecia niżej niesie dokładną treść zgłoszenia).
//  2. Budowanie tu własnego wykazu narzędzi byłoby drugą prawdą o wykazie
//     — dlatego ten plik wyłącznie CZYTA `shared.KomendyNarzedzi`
//     i odpowiada, czy komendy rodziny `subagent.*` już w nim są.
//
// Nazwa spoza wykazu nie znika po cichu: `narzedzia/zastrzezenia.go` oddaje
// modelowi czytelną odmowę „komenda stoi poza wykazem narzędzi" — granica
// przesuwa się w chwili, w której przesunie ją kontrakt.
package podagenci

import (
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// komendyPowolania wylicza komendy rodziny `subagent.*`, którymi model prowadzi
// podagentów: powołanie, wykaz i zbiór wyników. Kolejność jest kolejnością
// pracy orkiestratora.
func komendyPowolania() []shared.MessageType {
	return []shared.MessageType{
		shared.CommandSubagentSpawn,
		shared.CommandSubagentList,
		shared.CommandSubagentResultCollect,
	}
}

// drogaNarzedzia mówi o jednej komendzie rodziny, czy kontrakt wystawia ją
// modelowi jako narzędzie — i pod jaką nazwą.
type drogaNarzedzia struct {
	Komenda shared.MessageType
	// Narzedzie jest nazwą z wykazu kontraktu (np. `danaco_subagent_spawn`).
	// Puste dokładnie wtedy, gdy komenda stoi poza wykazem narzędzi.
	Narzedzie string
	Wpieta    bool
}

// drogiNarzedzi mierzy drogę narzędzia modelu dla całej rodziny `subagent.*`.
//
// Sprawdzenie idzie po `shared.KomendyNarzedzi` — jedynym odwzorowaniu nazwa
// narzędzia → komenda, wytworzonym z kontraktu. Nazwy narzędzia ten pakiet nie
// składa sam (przedrostek i separator należą do kontraktu); odnajduje ją po
// wartości komendy, więc zmiana konwencji nazw niczego tu nie psuje.
func drogiNarzedzi() []drogaNarzedzia {
	drogi := make([]drogaNarzedzia, 0, len(komendyPowolania()))
	for _, komenda := range komendyPowolania() {
		droga := drogaNarzedzia{Komenda: komenda}
		for nazwa, cel := range shared.KomendyNarzedzi {
			if cel == komenda {
				droga.Narzedzie, droga.Wpieta = nazwa, true
				break
			}
		}
		drogi = append(drogi, droga)
	}
	return drogi
}

// ZdanieODrodze składa jedno zdanie do dziennika rdzenia: które komendy rodziny
// są narzędziami modelu, a które czekają na pozycję w kontrakcie. Zdanie mówi
// wprost, GDZIE leży brak — bez tego meldunek nie prowadzi do naprawy.
func ZdanieODrodze() string {
	wpiete := []string{}
	czekaja := []string{}
	for _, droga := range drogiNarzedzi() {
		if droga.Wpieta {
			wpiete = append(wpiete, fmt.Sprintf("%s→%s", droga.Komenda, droga.Narzedzie))
			continue
		}
		czekaja = append(czekaja, string(droga.Komenda))
	}
	if len(czekaja) == 0 {
		return "droga narzędzia modelu: komplet rodziny subagent.* w wykazie narzędzi (" +
			strings.Join(wpiete, ", ") + ")"
	}
	return "droga narzędzia modelu: komendy " + strings.Join(czekaja, ", ") +
		" stoją poza wykazem narzędzi kontraktu — powołanie podagenta przez model czeka " +
		"na pozycje w `narzedzia.pozycje` (shared/contract.json) i przebieg generatora; " +
		"do tej chwili model dostaje czytelną odmowę serwera narzędzi"
}

// pozycjaWpiecia jest treścią JEDNEJ pozycji sekcji `narzedzia.pozycje`
// kontraktu — dokładnie w kształcie, w jakim czyta ją generator
// (`shared/gen/narzedzia.mjs`): komenda plus zdanie zastosowania, z którego
// powstaje druga połowa opisu narzędzia.
type pozycjaWpiecia struct {
	Komenda      string `json:"komenda"`
	Zastosowanie string `json:"zastosowanie"`
}

// PozycjeWpiecia niesie komplet pozycji zgłaszanych integratorowi do wpięcia
// w `narzedzia.pozycje`. Treść stoi w kodzie, a nie w meldunku, żeby dowód
// z uruchomienia i zgłoszenie wpięcia czytały JEDNO źródło.
//
// Zdania zastosowania mówią modelowi, KIEDY sięgnąć po narzędzie — wzorem
// pozycji istniejących. Powołanie: podagent to zadanie w tle pod oknem
// wykonawcy, powoływane w trakcie tury.
func PozycjeWpiecia() []pozycjaWpiecia {
	return []pozycjaWpiecia{
		{Komenda: string(shared.CommandSubagentSpawn),
			Zastosowanie: "Uzyj, aby powolac do pietnastu podagentow pracujacych w tle pod Twoim oknem wykonawcy; powolanie nie przerywa Twojej tury, a wyniki zbierzesz narzedziem subagent.result.collect"},
		{Komenda: string(shared.CommandSubagentList),
			Zastosowanie: "Uzyj, aby poznac stan podagentow powolanych pod oknem wykonawcy albo w calej karcie sesji, zanim siegniesz po ich wyniki"},
		{Komenda: string(shared.CommandSubagentResultCollect),
			Zastosowanie: "Uzyj, aby zebrac wyniki pracy podagentow; waitForAll kaze odpowiedzi czekac na zakonczenie wszystkich objetych zbieraniem"},
	}
}
