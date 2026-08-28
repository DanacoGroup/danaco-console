package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujRozmowe wpina domenę wiadomości okna komunikacji; treść modelu
// wraca osobno strumieniem, a message.stop odpowiada zawsze niezależnie od
// stanu pętli.
func zarejestrujRozmowe(r *Rejestr, rozmowa Rozmowa, e *emiter) {
	if r == nil || rozmowa == nil {
		return
	}

	r.Zarejestruj(shared.CommandMessageSend,
		obsluz(func(ctx context.Context, z shared.MessageSendRequest) (shared.MessageSendResponse, error) {
			w, err := rozmowa.Wyslij(ctx, z)
			if err == nil {
				e.wiadomosc(ctx, shared.ChangeKindCreated, w.Message)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandMessageStop, obsluz(rozmowa.Zatrzymaj))

	r.Zarejestruj(shared.CommandMessageList, obsluz(rozmowa.Wykaz))
}
