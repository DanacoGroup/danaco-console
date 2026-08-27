// Plik wpina rodzinę model.* — dziś jedną komendę wyboru kanału modelu, jako
// rozszerzenie portu okien rozgłaszające zmianę tym samym zdarzeniem co
// aktualizacja okna.
package core

import (
	"context"

	"danacoconsole/shared"
)

// KanalModelu jest portem rodziny model.* — portem okien rozszerzonym
// o wybór kanału modelu dla okna komunikacji.
type KanalModelu interface {
	Okna

	// UstawKanalModelu obsługuje model.channel.set, zwracając okna, których
	// kanał się zmienił.
	UstawKanalModelu(ctx context.Context, z shared.ModelChannelSetRequest) (
		shared.ModelChannelSetResponse, []shared.Window, error)
}

// zarejestrujKanalModelu wpina komendę rodziny model.* przez twardą asercję
// portu okien, żeby rdzeń padł głośno przy montażu, gdyby port przestał
// nieść wybór kanału.
func zarejestrujKanalModelu(r *Rejestr, m KanalModelu, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandModelChannelSet,
		obsluz(func(ctx context.Context, z shared.ModelChannelSetRequest) (
			shared.ModelChannelSetResponse, error) {

			odpowiedz, zmienione, err := m.UstawKanalModelu(ctx, z)
			if err != nil {
				return shared.ModelChannelSetResponse{}, err
			}
			for _, okno := range zmienione {
				e.okno(ctx, shared.ChangeKindUpdated, okno)
			}
			return odpowiedz, nil
		}))
}
