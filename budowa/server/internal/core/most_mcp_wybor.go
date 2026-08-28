// Plik dobiera i porządkuje nadania przekazywane generatorowi wpisu mcpServers oraz zawęża
// korzenie, bo reguły doboru odpowiadają na inne pytanie niż kształt wpisu.
package core

import (
	"sort"
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// KorzenieNadania zwraca korzenie obowiązujące dla nadania: podzbiór wskazany
// w nadaniu, a przy jego braku — komplet korzeni punktu. Korzenie nadania są
// zawężane do korzeni punktu, bo nadanie nie sięga poza punkt.
func KorzenieNadania(nadanie NadanieMostu) []string {
	if len(nadanie.Nadanie.Roots) == 0 {
		return niepusteKorzenie(nadanie.Punkt.Roots)
	}
	wskazane := niepusteKorzenie(nadanie.Nadanie.Roots)
	korzeniePunktu := niepusteKorzenie(nadanie.Punkt.Roots)
	if len(korzeniePunktu) == 0 {
		// Punkt bez korzeni obejmuje cały system plików maszyny, więc korzenie nadania wchodzą w całości.
		return wskazane
	}
	dozwolone := make(map[string]struct{}, len(korzeniePunktu))
	for _, korzen := range korzeniePunktu {
		dozwolone[korzen] = struct{}{}
	}
	wybrane := make([]string, 0, len(wskazane))
	for _, korzen := range wskazane {
		if _, jest := dozwolone[korzen]; jest {
			wybrane = append(wybrane, korzen)
		}
	}
	return wybrane
}

// niepusteKorzenie przycina korzenie katalogów po obu stronach i odrzuca wpisy puste z całego tego wykazu.
func niepusteKorzenie(korzenie []string) []string {
	wynik := make([]string, 0, len(korzenie))
	for _, korzen := range korzenie {
		if przyciety := strings.TrimSpace(korzen); przyciety != "" {
			wynik = append(wynik, przyciety)
		}
	}
	return wynik
}

// uporzadkowaneNadania odsiewa nadania nieczynne i punkty niebędące maszyną,
// a resztę porządkuje po kolejności nadania. Porządek jest rozstrzygający dla
// przyrostków rozróżniających klucze, więc ten sam zbiór zawsze daje te same
// klucze.
func uporzadkowaneNadania(nadania []NadanieMostu) []NadanieMostu {
	wybrane := make([]NadanieMostu, 0, len(nadania))
	for _, nadanie := range nadania {
		if nadanie.Punkt.Kind != shared.AccessPointKindMcpBridge {
			continue
		}
		if !nadanie.Punkt.Enabled || !nadanie.Nadanie.Enabled {
			continue
		}
		wybrane = append(wybrane, nadanie)
	}
	sort.SliceStable(wybrane, func(lewy, prawy int) bool {
		if wybrane[lewy].Nadanie.Order != wybrane[prawy].Nadanie.Order {
			return wybrane[lewy].Nadanie.Order < wybrane[prawy].Nadanie.Order
		}
		return wybrane[lewy].Nadanie.Id < wybrane[prawy].Nadanie.Id
	})
	return wybrane
}

// kluczUnikalny zwraca klucz niezajęty w mapie. Dwie maszyny o tej samej nazwie
// dostają klucze rozróżnione przyrostkiem, bo klucz powtórzony wyparłby wpis
// poprzedni i okno straciłoby jeden z dostępów.
func kluczUnikalny(klucz string, zajete map[string]WpisMostu) string {
	if _, jest := zajete[klucz]; !jest {
		return klucz
	}
	for numer := 2; ; numer++ {
		kandydat := klucz + "-" + strconv.Itoa(numer)
		if _, jest := zajete[kandydat]; !jest {
			return kandydat
		}
	}
}
