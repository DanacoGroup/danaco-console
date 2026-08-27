// Odpowiedzialność pliku: miara podobieństwa dwóch wypowiedzi — jedyna
// arytmetyka, na której stoi wykrywanie zgody, klastry opinii, scalanie
// powtórzeń, zbieżność stanowisk i dryf.
package core

import (
	"strings"
	"unicode"
)

// slowaNieznaczace — słowa, które w polszczyźnie łączą zdanie, zamiast nieść
// jego treść. Zostawione w zbiorze zawyżałyby podobieństwo dwóch dowolnych
// wypowiedzi, bo „nie", „że" i „w" pada w każdej.
var slowaNieznaczace = map[string]struct{}{
	"i": {}, "a": {}, "w": {}, "z": {}, "o": {}, "na": {}, "do": {}, "od": {}, "po": {},
	"za": {}, "to": {}, "że": {}, "nie": {}, "się": {}, "jest": {}, "są": {}, "był": {},
	"była": {}, "być": {}, "ale": {}, "lub": {}, "oraz": {}, "przez": {}, "dla": {},
	"jako": {}, "tym": {}, "ten": {}, "ta": {}, "te": {}, "tego": {}, "przy": {},
	"jeśli": {}, "gdy": {}, "już": {}, "tylko": {}, "bardzo": {}, "może": {}, "można": {},
}

// dlugoscSlowaZnaczacego — poniżej tej długości słowo pomija się nawet wtedy,
// gdy nie ma go w wykazie: skróty jedno- i dwuliterowe niosą zbyt mało, żeby
// świadczyć o wspólnym stanowisku.
const dlugoscSlowaZnaczacego = 3

// podobienstwoTekstow oddaje miarę zbieżności dwóch wypowiedzi w zakresie
// od zera do jedności, współczynnikiem Jaccarda.
func podobienstwoTekstow(pierwszy, drugi string) float64 {
	zbiorPierwszy := slowaZnaczace(pierwszy)
	zbiorDrugi := slowaZnaczace(drugi)
	if len(zbiorPierwszy) == 0 || len(zbiorDrugi) == 0 {
		return 0
	}
	wspolne := 0
	for slowo := range zbiorPierwszy {
		if _, jest := zbiorDrugi[slowo]; jest {
			wspolne++
		}
	}
	suma := len(zbiorPierwszy) + len(zbiorDrugi) - wspolne
	if suma == 0 {
		return 0
	}
	return float64(wspolne) / float64(suma)
}

// slowaZnaczace rozbija wypowiedź na zbiór słów niosących treść, do
// wyliczenia miary podobieństwa dwóch wypowiedzi.
func slowaZnaczace(tekst string) map[string]struct{} {
	zbior := make(map[string]struct{}, 16)
	for _, slowo := range strings.FieldsFunc(strings.ToLower(tekst), func(znak rune) bool {
		return !unicode.IsLetter(znak) && !unicode.IsDigit(znak)
	}) {
		if len([]rune(slowo)) < dlugoscSlowaZnaczacego {
			continue
		}
		if _, nieznaczace := slowaNieznaczace[slowo]; nieznaczace {
			continue
		}
		zbior[slowo] = struct{}{}
	}
	return zbior
}

// wspolneSlowaZnaczace oddaje słowa obecne w obu wypowiedziach. Służy nazwaniu
// punktu zgody: panel ma pokazać, CO uczestnicy powiedzieli tak samo, a nie
// samą liczbę mówiącą, że coś powiedzieli.
func wspolneSlowaZnaczace(pierwszy, drugi string) []string {
	zbiorDrugi := slowaZnaczace(drugi)
	wspolne := make([]string, 0, 8)
	widziane := make(map[string]struct{}, 8)
	// Kolejność bierze się z wypowiedzi pierwszej: przebieg po mapie w Go
	// jest losowy.
	for _, slowo := range strings.FieldsFunc(strings.ToLower(pierwszy), func(znak rune) bool {
		return !unicode.IsLetter(znak) && !unicode.IsDigit(znak)
	}) {
		if len([]rune(slowo)) < dlugoscSlowaZnaczacego {
			continue
		}
		if _, nieznaczace := slowaNieznaczace[slowo]; nieznaczace {
			continue
		}
		if _, jest := zbiorDrugi[slowo]; !jest {
			continue
		}
		if _, powtorzone := widziane[slowo]; powtorzone {
			continue
		}
		widziane[slowo] = struct{}{}
		wspolne = append(wspolne, slowo)
	}
	return wspolne
}
