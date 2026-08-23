// Odpowiedzialność pliku: grupa pozycji wykazu — po co dane narzędzie jest.
//
// Grupą jest obszar nazwy komendy: człon przed pierwszym `shared.SeparatorObszaru`
// (`agent.list` → `agent`, `session.window.create` → `session`). Kontrakt tę
// notację stanowi i sam się nią posługuje (`shared.ZdarzenieNieznanej` tnie
// nazwę dokładnie tak samo), więc grupa jest tu wyprowadzona, a nie zapisana.
// Wykaz własny „narzędzie → grupa" byłby drugą prawdą, która rozjedzie się przy
// pierwszym narzędziu dopisanym do kontraktu.
//
// Nazwa grupy jest kodem obszaru, nie zdaniem „do czego służy". Zdania mieszkają
// w `shared/contract.json`, w sekcji `obszary` (każda pozycja ma `nazwa` i
// `opis`), ale generator Go nie emituje ich do `shared/contract.go`. Nie
// dopisujemy ich tutaj z pamięci: opis zmyślony byłby zdaniem o kontrakcie,
// którego kontrakt nie mówi — do czasu udostępnienia opisów grupa niesie sam
// kod obszaru.
package narzedzia

import (
	"sort"
	"strings"

	"danacoconsole/shared"
)

// grupaKomendy wyprowadza grupę pozycji wykazu z nazwy jej komendy.
//
// Komenda bez separatora daje grupę równą całej nazwie — tak samo rozstrzyga
// `shared.ZdarzenieNieznanej`, i to jest jedyny powód, dla którego wynik
// `strings.Cut` o znalezieniu separatora jest tu pominięty: dwa różne
// rozstrzygnięcia tej samej notacji byłyby rozjazdem czekającym na okazję.
func grupaKomendy(komenda shared.MessageType) string {
	obszar, _, _ := strings.Cut(string(komenda), shared.SeparatorObszaru)
	return obszar
}

// Grupy zwraca nazwy grup występujących w wykazie, w porządku alfabetycznym,
// bez powtórzeń.
//
// Kolejność alfabetyczna, a nie kolejność kontraktu: grupy są gałęziami drzewa
// wyboru (`client/src/komponenty/menu-drzewo.ts`), a gałąź szuka się okiem po
// nazwie. Kolejność kontraktu zostaje tam, gdzie ma znaczenie — wewnątrz wykazu
// pozycji (`WykazZasiegu`).
func Grupy(wykaz []Narzedzie) []string {
	widziane := make(map[string]bool, len(wykaz))
	nazwy := make([]string, 0, len(wykaz))
	for _, pozycja := range wykaz {
		if pozycja.Grupa == "" || widziane[pozycja.Grupa] {
			continue
		}
		widziane[pozycja.Grupa] = true
		nazwy = append(nazwy, pozycja.Grupa)
	}
	sort.Strings(nazwy)
	return nazwy
}

// WGrupie zwraca pozycje wykazu należące do wskazanej grupy, w kolejności
// wykazu. Grupa nieznana daje pustkę — orzekanie, czy to pomyłka, należy do
// tego, kto grupę wskazał (`ekspert_wykaz.go` melduje kod nierozpoznany).
func WGrupie(wykaz []Narzedzie, grupa string) []Narzedzie {
	pozycje := make([]Narzedzie, 0, len(wykaz))
	for _, pozycja := range wykaz {
		if pozycja.Grupa == grupa {
			pozycje = append(pozycje, pozycja)
		}
	}
	return pozycje
}
