// Plik czyni przerwanie tury jawnym: wysłanie wiadomości do okna, które
// właśnie odpowiada, odmawia zamiast anulować turę po cichu. Tu mieszka też
// wykaz tur w biegu i wszystko, co go dotyka.
package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// odmowaTuryWBiegu jest odpowiedzią na `message.send` skierowany do okna, które
// właśnie odpowiada. Kod `conflict` jest tu właściwy: żądanie jest poprawne,
// a odmawia mu stan okna — nie treść, nie uprawnienie i nie brak bytu.
func odmowaTuryWBiegu(idOkna string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"Okno "+idOkna+" prowadzi turę i nie przyjmie kolejnej wiadomości. "+
			"Przerwanie odpowiedzi jest osobną decyzją Operatora, więc wysłanie "+
			"jej nie przerywa. Zatrzymaj odpowiedź przyciskiem zatrzymania "+
			"(komenda message.stop) albo zaczekaj na jej koniec, a potem wyślij "+
			"wiadomość ponownie."))
}

// zajmijBieg zajmuje okno pod nową turę i oddaje prawdę, gdy się to udało;
// sprawdzenie i zajęcie idą pod jednym zamkiem, nie dwoma wywołaniami, żeby
// uniknąć szczeliny wyścigu.
func (a *adapterRozmowy) zajmijBieg(idOkna string, anuluj context.CancelFunc) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, biegnie := a.biegnace[idOkna]; biegnie {
		return false
	}
	a.biegnace[idOkna] = anuluj
	return true
}

// Zatrzymaj przerywa turę okna. Odpowiada zawsze — przycisk zatrzymania jest
// czynny bez względu na stan tury; brak tury w biegu daje Stopped
// równe fałsz, nie błąd.
func (a *adapterRozmowy) Zatrzymaj(_ context.Context, z shared.MessageStopRequest) (shared.MessageStopResponse, error) {
	a.mu.Lock()
	anuluj, biegnie := a.biegnace[z.WindowId]
	delete(a.biegnace, z.WindowId)
	a.mu.Unlock()
	if biegnie {
		anuluj()
	}
	if a.petla != nil {
		a.petla.Zatrzymaj(z.WindowId) // przycisk Operatora wstrzymuje też bieg naprawczy
	}
	idWiadomosci := ""
	if z.MessageId != nil {
		idWiadomosci = *z.MessageId
	}
	return shared.MessageStopResponse{MessageId: idWiadomosci, Stopped: biegnie}, nil
}

// PrzerwijTure przerywa turę okna i mówi, czy jakaś biegła; ma dać biegnącemu
// procesowi sygnał, że ma się domknąć.
func (a *adapterRozmowy) PrzerwijTure(idOkna string) bool {
	a.mu.Lock()
	anuluj, biegnie := a.biegnace[idOkna]
	delete(a.biegnace, idOkna)
	a.mu.Unlock()
	if biegnie {
		anuluj()
	}
	return biegnie
}

// zapomnijBieg usuwa zakończoną turę z wykazu biegnących, aby wykaz niósł
// wyłącznie tury nadal aktywne.
func (a *adapterRozmowy) zapomnijBieg(idOkna string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.biegnace, idOkna)
}
