// Odpowiedzialność pliku: stan konta w rotacji i odczyt puli rotacji.
//
// Algorytm rotacji — które konto jest bieżące, kiedy przejść na następne, jak
// długo pamiętać wyczerpanie — mieszka wyłącznie w `injection.PulaKont`, żeby
// nie powstały dwie implementacje. Ten plik daje katalogowi dwie rzeczy:
//  1. `KontaRotacji` — uporządkowaną listę kont, które w ogóle wolno wziąć;
//  2. `OznaczStan` — trwały ślad tego, co pula już rozpoznała, oraz miejsce na
//     decyzję Operatora o zawieszeniu konta.
//
// Repozytorium nie wygasza wyczerpania po czasie i nie wybiera konta bieżącego —
// robi to pula, która jedyna zna chwilę wywołania.
package dane

import (
	"context"
	"fmt"

	"danacoconsole/shared"
)

// StanKonta jest zdatnością konta w rotacji. Kontrakt tego wyliczenia nie zna —
// jest wewnętrznym stanem katalogu, odpowiadającym kolumnie `konto.stan`.
// Słownik poniżej jest jedynym miejscem, w którym te wartości występują
// w kodzie Go.
type StanKonta string

const (
	// StanKontaAktywne — konto zdatne do wzięcia przez pulę.
	StanKontaAktywne StanKonta = "aktywne"
	// StanKontaWyczerpane — pula rozpoznała wyczerpanie limitu; konto zostaje
	// w rotacji, bo limit sam się odnawia, a chwilę odnowienia niesie
	// `wyczerpane_do`.
	StanKontaWyczerpane StanKonta = "wyczerpane"
	// StanKontaZawieszone — Operator wyłączył konto z rotacji do odwołania.
	StanKontaZawieszone StanKonta = "zawieszone"
)

// wartosciBazyStanKonta i wartosciKontraktuStanKonta odwzorowują CHECK kolumny
// `konto.stan`. Wzorzec jest ten sam co w `shared` — przekład w jednym miejscu,
// literały nigdzie indziej.
var wartosciBazyStanKonta = map[StanKonta]string{
	StanKontaAktywne:    "aktywne",
	StanKontaWyczerpane: "wyczerpane",
	StanKontaZawieszone: "zawieszone",
}

var wartosciKontraktuStanKonta = map[string]StanKonta{
	"aktywne":    StanKontaAktywne,
	"wyczerpane": StanKontaWyczerpane,
	"zawieszone": StanKontaZawieszone,
}

// stanKontaNaBaze przekłada stan na wartość kolumny; pusty stan znaczy „aktywne".
func stanKontaNaBaze(stan StanKonta) (string, error) {
	return naBaze(wartosciBazyStanKonta, stan, StanKontaAktywne, "konto.stan")
}

func stanKontaZBazy(kolumna string) (StanKonta, error) {
	return zBazy(wartosciKontraktuStanKonta, kolumna, "konto.stan")
}

// rodzajKontaNaBaze i rodzajKontaZBazy sięgają po słowniki kontraktu wskazane
// wprost dla kolumny `konto.rodzaj`. Brak wskazania rodzaju znaczy
// konto programu code CLI — kanał główny jest kanałem domyślnym.
func rodzajKontaNaBaze(rodzaj shared.AccountKind) (string, error) {
	return naBaze(shared.WartosciBazyAccountKind, rodzaj, shared.AccountKindCli, "konto.rodzaj")
}

func rodzajKontaZBazy(kolumna string) (shared.AccountKind, error) {
	return zBazy(shared.WartosciKontraktuAccountKind, kolumna, "konto.rodzaj")
}

// listaKontRotacji zwraca konta zdatne do rotacji jednego rodzaju. Konto
// domyślne idzie pierwsze, dalej rozstrzyga kolejność nadana przez Operatora.
// Konta wyczerpane zostają na liście — o ich pominięciu rozstrzyga pula, która
// zna chwilę wywołania i chwilę odnowienia limitu.
const listaKontRotacji = `SELECT ` + kolumnyKonta + ` FROM konto
	WHERE rodzaj = ? AND aktywne = 1 AND stan <> '` + string(StanKontaZawieszone) + `'
	ORDER BY domyslne DESC, kolejnosc, id`

// KontaRotacji zwraca pulę kont wskazanego rodzaju w kolejności rotacji.
// Pusty wynik nie jest błędem: rdzeń rusza bez kont, a brak konta zgłasza
// dopiero próba rozmowy.
func (r *repozytoriumKont) KontaRotacji(ctx context.Context, rodzaj shared.AccountKind) ([]Konto, error) {
	wartosc, err := rodzajKontaNaBaze(rodzaj)
	if err != nil {
		return nil, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, listaKontRotacji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, wartosc)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać puli kont rodzaju %q: %w", wartosc, err)
	}
	defer wiersze.Close()
	return zbierzKonta(wiersze)
}
