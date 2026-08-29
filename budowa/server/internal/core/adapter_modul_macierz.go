package core

import (
	"context"
	"fmt"

	"danacoconsole/server/internal/dane"
)

// adapterMacierzy wypełnia port Macierz odczytem tabeli `srodowisko_modul`.
// Odwzorowanie „moduł → środowiska, w których jest widoczny” składa się
// jednym złączeniem, bez pętli po środowiskach.
type adapterMacierzy struct {
	macierz dane.RepozytoriumMacierzy
}

// nowyAdapterMacierzy wiąże port Macierz z repozytorium macierzy
// widoczności modułów w środowiskach pracy.
func nowyAdapterMacierzy(macierz dane.RepozytoriumMacierzy) *adapterMacierzy {
	return &adapterMacierzy{macierz: macierz}
}

// ── macierz widoczności ─────────────────────────────────────────────────────

// KodySrodowisk zwraca kody środowisk, w których widoczny jest każdy moduł,
// jednym złączeniem tabel bazy.
func (a *adapterMacierzy) KodySrodowisk(ctx context.Context) (map[int64][]string, error) {
	if err := a.sprawdzKatalog(); err != nil {
		return nil, err
	}
	kody, err := a.macierz.KodySrodowiskModulow(ctx)
	if err != nil {
		return nil, fmt.Errorf("serwer: nie można odczytać macierzy widoczności: %w", err)
	}
	return kody, nil
}

// sprawdzKatalog zwraca błąd, gdy repozytorium nie jest wpięte. Pusta
// macierz to nie to samo co macierz nieodczytana.
func (a *adapterMacierzy) sprawdzKatalog() error {
	if a.macierz == nil {
		return fmt.Errorf("serwer: repozytorium macierzy widoczności nie jest wpięte")
	}
	return nil
}
