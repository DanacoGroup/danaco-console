// Odpowiedzialność pliku: licznik wykazu — ile pozycji i ile bajtów ładuje ten
// zasięg, zanim padnie pierwsze słowo zadania.
//
// ── DLACZEGO NIE SAM LICZNIK POZYCJI ────────────────────────────────────────
// Operatorowi sama liczba pozycji nie mówi nic; liczba bajtów i rząd żetonów
// okna kontekstu mówi wszystko. Limit, w który się uderza,
// jest limitem okna kontekstu, a okno kontekstu liczy się w żetonach, nie
// w narzędziach. Licznik podający wyłącznie pozycje kazałby Operatorowi
// przeliczać je w głowie na koszt — czyli dowiedziałby się o przekroczeniu
// dopiero z cichej degradacji, przed którą ten cały zasięg powstał.
//
// ── DLACZEGO KSZTAŁT PROTOKOŁU MIESZKA TUTAJ ────────────────────────────────
// Bajty mają być bajtami tej samej odpowiedzi, którą dostanie model — inaczej
// pomiar jest oszacowaniem podanym jako pomiar. Dlatego trzy pola,
// z których protokół składa `tools/list`, stoją w jednym miejscu: tutaj.
// Warstwa protokołu (`cmd/danaco-narzedzia/stdio`) bierze je stąd, zamiast
// składać drugi raz po swojemu. Pakiet nadal nie zna ramki JSON-RPC ani metod —
// zna wyłącznie kształt danych jednej pozycji, i to jest cena za to, że pomiar
// nie kłamie.
//
// ── ŻETONY SĄ PRZELICZENIEM, NIE POMIAREM — I TAK SĄ NAZWANE ────────────────
// Licznik oddaje bajty, bo bajty umie policzyć dokładnie. Przelicznika na żetony
// tu nie ma: zależy od tokenizatora kanału modelu, którego ten proces nie zna,
// a liczba podana jako dokładna, a wyprowadzona z założenia, byłaby drugą prawdą.
// Rząd wielkości podaje się w zdaniu dziennika, nie w polu struktury.
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

// PomiarGrupy jest udziałem jednej grupy w wykazie.
type PomiarGrupy struct {
	// Nazwa jest nazwą grupy (`grupa.go`).
	Nazwa string
	// Pozycji liczy narzędzia tej grupy.
	Pozycji int
	// Bajtow liczy wagę tych narzędzi w złożonej odpowiedzi.
	Bajtow int
}

// PomiarWykazu jest odpowiedzią na pytanie „ile ten zasięg ładuje".
type PomiarWykazu struct {
	// Pozycji liczy narzędzia w wykazie.
	Pozycji int
	// Bajtow jest wagą złożonej odpowiedzi `tools/list` — dokładnie tej, którą
	// dostanie model.
	Bajtow int
	// Grupy niesie rozbicie na grupy, w porządku alfabetycznym nazw. Suma pozycji
	// grup równa się `Pozycji`; suma bajtów grup jest mniejsza od `Bajtow`
	// o narzut listy (nawiasy i przecinki), i to jest powiedziane, a nie ukryte.
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

// Zmierz liczy pozycje i bajty wykazu wraz z rozbiciem na grupy.
//
// Wykaz, którego nie da się zakodować, daje bajty zerowe przy prawdziwej liczbie
// pozycji. Kłamstwem byłoby oddanie liczby zmyślonej; zerem jest tu widoczne,
// że wagi nie zmierzono. Przypadek jest teoretyczny — schemat powstaje
// z kontraktu i koduje się zawsze — więc nie rozrasta się w osobną drogę błędu.
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

// wagaWykazu koduje wykaz w kształcie protokołu i zwraca liczbę bajtów.
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
