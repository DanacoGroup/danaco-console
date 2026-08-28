// Licznik wykazu mierzy, ile pozycji i ile bajtów kształtu protokołu
// tools/list ładuje dany zasięg, zanim padnie pierwsze słowo zadania.
package narzedzia

import (
	"encoding/json"
	"sort"
)

// Nazwy pól pozycji wykazu w kształcie `tools/list`. Stoją stałymi, bo padają
// i przy składaniu odpowiedzi, i przy pomiarze — literał powtórzony jest
// literałem, który da się rozjechać.
const (
	polePozycjiNazwa   = "name"
	polePozycjiOpis    = "description"
	polePozycjiSchemat = "inputSchema"
)

// PomiarGrupy jest udziałem jednej grupy narzędzi w wykazie, wyrażonym liczbą
// pozycji i wagą w bajtach.
type PomiarGrupy struct {
	// Nazwa jest nazwą grupy (`grupa.go`).
	Nazwa string
	// Pozycji liczy narzędzia tej grupy.
	Pozycji int
	// Bajtow liczy wagę tych narzędzi w złożonej odpowiedzi.
	Bajtow int
}

// PomiarWykazu jest odpowiedzią na pytanie, ile dany zasięg ładuje: liczbą
// pozycji, bajtów i rozbiciem na grupy.
type PomiarWykazu struct {
	// Pozycji liczy narzędzia w wykazie.
	Pozycji int
	// Bajtow jest wagą złożonej odpowiedzi `tools/list` — dokładnie tej, którą
	// dostanie model.
	Bajtow int
	// Grupy niesie rozbicie na grupy, alfabetycznie; suma bajtów grup jest
	// mniejsza o narzut listy.
	Grupy []PomiarGrupy
}

// PozycjeWykazu składa wykaz w kształt danych `tools/list`.
//
// Jedyne miejsce, w którym te trzy pola powstają — bierze je stąd i warstwa
// protokołu, i licznik, więc zmiana kształtu przestawia obie naraz.
func PozycjeWykazu(wykaz []Narzedzie) []map[string]any {
	pozycje := make([]map[string]any, 0, len(wykaz))
	for _, narzedzie := range wykaz {
		pozycje = append(pozycje, map[string]any{
			polePozycjiNazwa:   narzedzie.Nazwa,
			polePozycjiOpis:    narzedzie.Opis,
			polePozycjiSchemat: narzedzie.Schemat,
		})
	}
	return pozycje
}

// Zmierz liczy pozycje i bajty wykazu wraz z rozbiciem na grupy; wykaz nie do
// zakodowania daje bajty zerowe przy liczbie pozycji prawdziwej.
func Zmierz(wykaz []Narzedzie) PomiarWykazu {
	pomiar := PomiarWykazu{Pozycji: len(wykaz), Bajtow: wagaWykazu(wykaz)}
	pozycji := map[string]int{}
	bajtow := map[string]int{}
	for _, narzedzie := range wykaz {
		pozycji[narzedzie.Grupa]++
		bajtow[narzedzie.Grupa] += wagaPozycji(narzedzie)
	}
	nazwy := make([]string, 0, len(pozycji))
	for nazwa := range pozycji {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)
	pomiar.Grupy = make([]PomiarGrupy, 0, len(nazwy))
	for _, nazwa := range nazwy {
		pomiar.Grupy = append(pomiar.Grupy, PomiarGrupy{
			Nazwa:   nazwa,
			Pozycji: pozycji[nazwa],
			Bajtow:  bajtow[nazwa],
		})
	}
	return pomiar
}

// wagaWykazu koduje wykaz w kształcie protokołu tools/list i zwraca liczbę
// bajtów zakodowanej odpowiedzi.
func wagaWykazu(wykaz []Narzedzie) int {
	bajty, err := json.Marshal(PozycjeWykazu(wykaz))
	if err != nil {
		return 0
	}
	return len(bajty)
}

// wagaPozycji liczy bajty jednej pozycji, bez nawiasów listy. Dzięki temu suma
// bajtów grup jest sumą samych narzędzi, a różnica wobec `PomiarWykazu.Bajtow`
// jest dokładnie narzutem listy — wielkością znaną, a nie rozjazdem.
func wagaPozycji(narzedzie Narzedzie) int {
	pozycje := PozycjeWykazu([]Narzedzie{narzedzie})
	bajty, err := json.Marshal(pozycje[0])
	if err != nil {
		return 0
	}
	return len(bajty)
}
