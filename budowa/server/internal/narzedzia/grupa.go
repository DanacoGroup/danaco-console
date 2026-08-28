// Grupa pozycji wykazu jest obszarem nazwy komendy: człon nazwy przed
// pierwszym separatorem obszaru, wyprowadzony z kontraktu, nie zapisany
// osobnym wykazem.
package narzedzia

import (
	"sort"
	"strings"

	"danacoconsole/shared"
)

// grupaKomendy wyprowadza grupę pozycji wykazu z nazwy jej komendy; komenda
// bez separatora daje grupę równą całej nazwie.
func grupaKomendy(komenda shared.MessageType) string {
	obszar, _, _ := strings.Cut(string(komenda), shared.SeparatorObszaru)
	return obszar
}

// Grupy zwraca nazwy grup występujących w wykazie, w porządku alfabetycznym,
// bez powtórzeń, do gałęzi drzewa wyboru.
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
