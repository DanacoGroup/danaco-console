package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujWiazanieSesji wpina więź klienta z kartą sesji: ogniskowanie i powiązanie połączenia
// z sesją trwającą na rdzeniu.
func zarejestrujWiazanieSesji(r *Rejestr, w WiazanieSesji, e *emiter) {
	if r == nil || w == nil {
		return
	}

	r.Zarejestruj(shared.CommandSessionFocus,
		obsluz(func(ctx context.Context, z shared.SessionFocusRequest) (shared.SessionFocusResponse, error) {
			wynik, err := w.Ogniskuj(ctx, z)
			if err == nil {
				e.ognisko(z.ClientId, wynik)
			}
			return wynik, err
		}))

	r.Zarejestruj(shared.CommandSessionBind, obsluz(w.Powiaz))
}

// ognisko rozgłasza zmianę ogniska karty sesji, niosąc klienta, na którym zmiana nastąpiła, do
// pozostałych urządzeń konta.
func (e *emiter) ognisko(idKlienta string, w shared.SessionFocusResponse) {
	e.wyslij(shared.EventSessionFocusChanged, w.SessionId, shared.SessionFocusChangedEvent{
		SessionId:         w.SessionId,
		ClientId:          idKlienta,
		WindowId:          w.WindowId,
		PreviousSessionId: w.PreviousSessionId,
		FocusedAt:         w.FocusedAt,
	})
}
