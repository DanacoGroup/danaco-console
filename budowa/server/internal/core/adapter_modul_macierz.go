package core

import (
	"context"
	"fmt"

	"danacoconsole/server/internal/dane"
)

// adapterMacierzy wypełnia port Macierz odczytem tabeli `srodowisko_modul`.
//
// Odwzorowanie „moduł → środowiska, w których jest widoczny" składa się jednym
// złączeniem, bez pętli po środowiskach: czytają je `home.enter`,
// `environment.list`, `environment.enter` i `module.list`, czyli każde wejście
// do pracy.
//
// Adapter zwraca błąd zwykły, nie protokolarny. Nie obsługuje żadnej komendy
// (patrz `handlers_macierz.go`), a jego czytelnikiem jest nawigacja, która sama
// zamienia błąd na odpowiedź protokołu.

type adapterMacierzy struct {
	macierz dane.RepozytoriumMacierzy
}

// nowyAdapterMacierzy wiąże port z repozytorium macierzy.
func nowyAdapterMacierzy(macierz dane.RepozytoriumMacierzy) *adapterMacierzy {
	return &adapterMacierzy{macierz: macierz}
}

// ── macierz widoczności ─────────────────────────────────────────────────────

// KodySrodowisk zwraca kody środowisk, w których widoczny jest każdy moduł.
func (a *adapterMacierzy) KodySrodowisk(ctx context.Context) (map[int64][]string, error) {
	if err := a.sprawdzKatalog(); err != nil {
		return nil, err
	}
	kody, err := a.macierz.KodySrodowiskModulow(ctx)
	if err != nil {
		return nil, fmt.Errorf("rdzeń: nie można odczytać macierzy widoczności: %w", err)
	}
	return kody, nil
}

// sprawdzKatalog zwraca błąd, gdy repozytorium nie jest wpięte.
//
// Pusta macierz to nie to samo co macierz nieodczytana: pierwsza znaczy „żaden
// moduł nie stoi w nawigacji", druga znaczy błąd złożenia rdzenia. Zwrócenie
// pustki zamiast błędu skasowałoby całą boczną nawigację bez wyjaśnienia.
func (a *adapterMacierzy) sprawdzKatalog() error {
	if a.macierz == nil {
		return fmt.Errorf("rdzeń: repozytorium macierzy widoczności nie jest wpięte")
	}
	return nil
}
