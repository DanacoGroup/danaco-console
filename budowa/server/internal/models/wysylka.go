package models

import (
	"context"
	"fmt"
	"strings"
)

// Wyslij kieruje zapytanie do kanału wskazanego przez okno komunikacji i nadaje
// jego strumień do ujścia; jest jedyną drogą, którą warstwy wyższe sięgają po
// kanał, nie znając ani adapterów, ani dostawców.
func (r *Rejestr) Wyslij(ctx context.Context, z Zapytanie, u Ujscie) error {
	klucz := strings.TrimSpace(z.Kanal)
	if klucz == "" {
		return NadajBlad(ctx, u, z, fmt.Errorf("models: okno %s bez wskazanego kanału modelu", z.Okno()))
	}
	kanal, jest := r.Kanal(klucz)
	if !jest {
		return NadajBlad(ctx, u, z, fmt.Errorf("models: kanał %q nie istnieje w rejestrze albo jest nieczynny", klucz))
	}
	if err := kanal.Wyslij(ctx, z, u); err != nil {
		return NadajBlad(ctx, u, z, err)
	}
	return nil
}
