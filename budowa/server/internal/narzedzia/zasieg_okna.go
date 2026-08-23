// Odpowiedzialność pliku: zasięg okna wywołania.
//
// Serwer narzędzi należy do jednego okna rozmowy. Wpis `danaco` w konfiguracji
// MCP powstaje osobno dla każdego okna i niesie jego identyfikator, więc
// narzędzie zawsze wie, z którego okna przyszło wywołanie. Pole `windowId`
// pominięte przez model dostaje identyfikator tego okna — nigdy cudzy i nigdy
// zgadnięty.
//
// Wskazanie jawne zostaje. Kiedy model podaje `windowId` sam, wartość idzie do
// rdzenia bez zmiany: pętla koordynator–wykonawca polega na tym, że okno
// koordynatora wysyła wiadomość do okna wykonawcy, a opisy narzędzi
// w kontrakcie mówią to wprost. Podmienianie wskazania jawnego na własne okno
// zamknęłoby tę pętlę i rozminęło serwer z kontraktem, który go opisuje.
//
// Pole rozpoznaje się po nazwie z deklaracji kontraktu, nie po własnym wykazie
// komend okna: narzędzie bez pola `windowId` przechodzi nietknięte.
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

// niesiePoleOkna mówi, czy narzędzie w ogóle przyjmuje okno rozmowy.
func niesiePoleOkna(parametry []string) bool {
	for _, nazwa := range parametry {
		if nazwa == poleOkna {
			return true
		}
	}
	return false
}

// brakWskazaniaOkna rozstrzyga, czy model okna nie wskazał.
//
// Wartość innego rodzaju niż napis zostaje nietknięta: jest wskazaniem wadliwym,
// a orzekanie o kształcie treści żądania należy do rdzenia, nie do rozdzielni.
// Podmiana takiej wartości na własne okno ukryłaby pomyłkę modelu.
func brakWskazaniaOkna(argumenty map[string]any) bool {
	wartosc, jest := argumenty[poleOkna]
	if !jest {
		return true
	}
	napis, jestNapisem := wartosc.(string)
	return jestNapisem && napis == ""
}
