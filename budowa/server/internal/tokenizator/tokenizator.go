// Pakiet tokenizator liczy żetony tekstu tak, jak liczy je model — po to, żeby
// zajętość okna kontekstu (`context.usage.get`) była POMIAREM, a nie
// oszacowaniem.
//
// ── Dlaczego nie „znaki podzielone przez cztery" ────────────────────────────
// Przybliżenie po długości tekstu myli się o kilkadziesiąt procent na języku
// polskim (znaki diakrytyczne rozpadają się na osobne żetony), a jeszcze
// bardziej na kodzie źródłowym i na zapisie strukturalnym. Licznik, który myli
// się o kilkadziesiąt procent, jest gorszy niż jego brak: Operator widzi pasek
// „w normie" i traci turę na przepełnieniu okna.
//
// ── Dlaczego biblioteka wkompilowana, a nie program ─────────────────────────
// Słownik BPE jest wkompilowany w binarium rdzenia razem z pakietem. Nie ma tu
// ani jednego procesu potomnego, ani jednego pobrania z sieci w trakcie
// żądania: tokenizator ma działać na maszynie odciętej od świata tak samo, jak
// na maszynie deweloperskiej.
//
// ── Czego ten pakiet nie robi ───────────────────────────────────────────────
// Nie udaje, że zna podział na żetony każdego modelu świata. Zna rodziny
// słowników, które naprawdę ma; model spoza nich dostaje słownik najbliższy
// wraz z jawnym powiedzeniem, którym słownikiem policzono (`Nazwa`). Nazwa
// wchodzi do odpowiedzi kontraktu (`ContextUsage.tokenizer`), więc czytelnik
// wie, czym zmierzono — a nie dostaje liczby bez świadka.
package tokenizator

import (
	"strings"
	"sync"

	"github.com/tiktoken-go/tokenizer"
)

// Licznik jest jednym słownikiem BPE wraz z jego nazwą.
type Licznik struct {
	nazwa     string
	kodowanie tokenizer.Codec
}

// Nazwa oddaje nazwę słownika, którym policzono — wchodzi do odpowiedzi wprost.
func (l Licznik) Nazwa() string { return l.nazwa }

// Policz zwraca liczbę żetonów tekstu.
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

// slownikDomyslny jest brany dla modeli spoza wykazu — w tym dla rodziny
// Claude, której słownika nikt nie publikuje.
//
// To jest wybór świadomy, a nie zaniedbanie: `cl100k_base` jest najbliższym
// dostępnym podziałem dla modeli tej klasy, a odpowiedź MÓWI, że policzono
// właśnie nim. Czytelnik dostaje liczbę wraz z nazwą miary, więc wie, na ile
// jest wiążąca — inaczej niż przy liczbie podanej bez słowa.
const slownikDomyslny = tokenizer.Cl100kBase

// DlaModelu zwraca licznik właściwy nazwie modelu.
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

// dlaKodowania buduje albo oddaje z pamięci licznik jednego słownika.
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
