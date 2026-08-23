package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujRozmowe wpina domenę wiadomości okna komunikacji.
//
// Odpowiedź na message.send potwierdza przyjęcie wiadomości; treść modelu wraca
// osobno, strumieniem stream.chunk przez Nadajnik. Rozdzielenie jest
// celowe: strumień trwa dłużej niż wykonanie komendy, a klient musi dostać
// potwierdzenie od razu.
//
// message.stop odpowiada zawsze — przycisk zatrzymania jest czynny bez względu
// na stan pętli. Gdy nie było czego zatrzymywać, wraca Stopped równe
// fałsz, a nie błąd.
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
