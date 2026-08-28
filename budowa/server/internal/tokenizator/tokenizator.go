// Pakiet tokenizator liczy żetony tekstu tak, jak liczy je model, żeby zajętość okna kontekstu była pomiarem, a nie oszacowaniem.
package tokenizator

import (
	"strings"
	"sync"

	"github.com/tiktoken-go/tokenizer"
)

// Licznik jest jednym słownikiem BPE wraz z jego nazwą, używanym do liczenia żetonów w każdym podanym tekście.
type Licznik struct {
	nazwa     string
	kodowanie tokenizer.Codec
}

// Metoda Nazwa oddaje nazwę słownika, którym policzono tekst; wchodzi do odpowiedzi kontraktu wprost dla czytelnika.
func (l Licznik) Nazwa() string { return l.nazwa }

// Metoda Policz zwraca liczbę żetonów danego tekstu policzoną słownikiem BPE przypisanym do tego licznika.
func (l Licznik) Policz(tekst string) int {
	if l.kodowanie == nil || tekst == "" {
		return 0
	}
	zetony, _, err := l.kodowanie.Encode(tekst)
	if err != nil {
		return 0
	}
	return len(zetony)
}

// pamiec trzyma słowniki już zbudowane. Budowa słownika to rozpakowanie
// kilkunastu megabajtów tablic — raz na proces, nie raz na żądanie.
var (
	zamek  sync.Mutex
	pamiec = map[tokenizer.Encoding]Licznik{}
)

// slownikiModeli wiąże przedrostek nazwy modelu ze słownikiem, którym ten model
// naprawdę liczy. Dopasowanie idzie po przedrostku, bo nazwy wersji rosną
// (`gpt-4o-2024-…`), a słownik zostaje ten sam.
var slownikiModeli = []struct {
	Przedrostek string
	Kodowanie   tokenizer.Encoding
}{
	{"gpt-4o", tokenizer.O200kBase},
	{"o1", tokenizer.O200kBase},
	{"o3", tokenizer.O200kBase},
	{"gpt-4", tokenizer.Cl100kBase},
	{"gpt-3.5", tokenizer.Cl100kBase},
	{"text-embedding", tokenizer.Cl100kBase},
}

// slownikDomyslny jest brany dla modeli spoza wykazu, w tym dla modeli, których słownika nikt nie publikuje.
const slownikDomyslny = tokenizer.Cl100kBase

// Funkcja DlaModelu zwraca licznik właściwy nazwie modelu wskazanego jako argument tego wywołania funkcji.
func DlaModelu(model string) (Licznik, error) {
	nazwa := strings.ToLower(strings.TrimSpace(model))
	kodowanie := slownikDomyslny
	for _, wpis := range slownikiModeli {
		if strings.HasPrefix(nazwa, wpis.Przedrostek) {
			kodowanie = wpis.Kodowanie
			break
		}
	}
	return dlaKodowania(kodowanie)
}

// Funkcja dlaKodowania buduje albo oddaje z pamięci podręcznej licznik jednego wskazanego słownika kodowania.
func dlaKodowania(kodowanie tokenizer.Encoding) (Licznik, error) {
	zamek.Lock()
	defer zamek.Unlock()
	if licznik, jest := pamiec[kodowanie]; jest {
		return licznik, nil
	}
	kodek, err := tokenizer.Get(kodowanie)
	if err != nil {
		return Licznik{}, err
	}
	licznik := Licznik{nazwa: string(kodowanie), kodowanie: kodek}
	pamiec[kodowanie] = licznik
	return licznik, nil
}
