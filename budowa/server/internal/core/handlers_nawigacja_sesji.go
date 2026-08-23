package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujWiazanieSesji wpina więź klienta z kartą sesji: ogniskowanie
// (`session.focus`) i powiązanie połączenia z sesją trwającą na rdzeniu
// (`session.bind`).
//
// Powiązanie następuje po uzgodnieniu połączenia: klient albo wchodzi na stronę
// główną, albo wraca do sesji, która przetrwała jego rozłączenie wraz
// z procesami okien.
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

// ognisko rozgłasza zmianę ogniska karty sesji.
//
// Zdarzenie jest nośnikiem synchronizacji ogniska między urządzeniami konta,
// dlatego niesie klienta, na którym zmiana nastąpiła: pozostałe urządzenia mają
// rozpoznać własną zmianę i cudzą, a nie przejąć ognisko po sobie nawzajem.
func (e *emiter) ognisko(idKlienta string, w shared.SessionFocusResponse) {
	e.wyslij(shared.EventSessionFocusChanged, w.SessionId, shared.SessionFocusChangedEvent{
		SessionId:         w.SessionId,
		ClientId:          idKlienta,
		WindowId:          w.WindowId,
		PreviousSessionId: w.PreviousSessionId,
		FocusedAt:         w.FocusedAt,
	})
}
