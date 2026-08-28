// Zasięg okna wywołania: serwer narzędzi należy do jednego okna rozmowy
// i uzupełnia pole okna pominięte przez model identyfikatorem tego okna,
// nigdy cudzym.
package narzedzia

// poleOkna nazywa pole treści żądania niosące okno rozmowy. Nazwa pochodzi
// z kontraktu — pojawia się w deklaracjach narzędzi, a pakiet wyłącznie ją
// rozpoznaje.
const poleOkna = "windowId"

// zZasiegiemOkna zwraca argumenty uzupełnione o okno serwera.
//
// Zwraca kopię: mapa przychodzi z odkodowanego wywołania i nie ma powodu, by
// rozdzielnia zmieniała cudzą strukturę w miejscu.
func zZasiegiemOkna(argumenty map[string]any, parametry []string, okno string) map[string]any {
	uzupelnione := make(map[string]any, len(argumenty)+1)
	for nazwa, wartosc := range argumenty {
		uzupelnione[nazwa] = wartosc
	}
	if okno == "" || !niesiePoleOkna(parametry) || !brakWskazaniaOkna(argumenty) {
		return uzupelnione
	}
	uzupelnione[poleOkna] = okno
	return uzupelnione
}

// niesiePoleOkna mówi, czy dane narzędzie w ogóle przyjmuje okno rozmowy jako
// jedno z pól treści żądania.
func niesiePoleOkna(parametry []string) bool {
	for _, nazwa := range parametry {
		if nazwa == poleOkna {
			return true
		}
	}
	return false
}

// brakWskazaniaOkna rozstrzyga, czy model okna nie wskazał; wartość innego
// rodzaju niż napis zostaje nietknięta.
func brakWskazaniaOkna(argumenty map[string]any) bool {
	wartosc, jest := argumenty[poleOkna]
	if !jest {
		return true
	}
	napis, jestNapisem := wartosc.(string)
	return jestNapisem && napis == ""
}
