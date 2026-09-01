// Plik czyni przerwanie tury jawnym: wysłanie wiadomości do okna, które właśnie odpowiada, odmawia zamiast anulować turę po cichu. Tu mieszka też wykaz tur w biegu i wszystko, co go dotyka.
package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// odmowaTuryWBiegu jest odpowiedzią na `message.send` skierowany do okna, które właśnie odpowiada. Kod `conflict` jest tu właściwy: żądanie jest poprawne, a odmawia mu stan okna — nie treść, nie uprawnienie i nie brak bytu.
func odmowaTuryWBiegu(idOkna string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"Okno "+idOkna+" prowadzi turę i nie przyjmie kolejnej wiadomości. "+
			"Przerwanie odpowiedzi jest osobną decyzją Operatora, więc wysłanie "+
			"jej nie przerywa. Zatrzymaj odpowiedź przyciskiem zatrzymania "+
			"(komenda message.stop) albo zaczekaj na jej koniec, a potem wyślij "+
			"wiadomość ponownie."))
}

// biegTury jest wpisem tury w wykazie biegnących. Wpis ma tożsamość wskaźnika: tura kończąca się po zatrzymaniu wykreśla wyłącznie wpis własny, nie wpis tury, która zajęła okno po niej.
type biegTury struct {
	anuluj context.CancelFunc
}

// zajmijBieg zajmuje okno pod nową turę i oddaje jej wpis; zero znaczy okno zajęte. Sprawdzenie i zajęcie idą pod jednym zamkiem, nie dwoma wywołaniami, żeby uniknąć szczeliny wyścigu.
func (a *adapterRozmowy) zajmijBieg(idOkna string, anuluj context.CancelFunc) *biegTury {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, biegnie := a.biegnace[idOkna]; biegnie {
		return nil
	}
	bieg := &biegTury{anuluj: anuluj}
	a.biegnace[idOkna] = bieg
	return bieg
}

// Zatrzymaj przerywa turę okna. Odpowiada zawsze — przycisk zatrzymania jest czynny bez względu na stan tury; brak tury w biegu daje Stopped równe fałsz, nie błąd.
func (a *adapterRozmowy) Zatrzymaj(_ context.Context, z shared.MessageStopRequest) (shared.MessageStopResponse, error) {
	biegnie := a.PrzerwijTure(z.WindowId)
	if a.petla != nil {
		a.petla.Zatrzymaj(z.WindowId) // przycisk Operatora wstrzymuje też bieg naprawczy
	}
	idWiadomosci := ""
	if z.MessageId != nil {
		idWiadomosci = *z.MessageId
	}
	return shared.MessageStopResponse{MessageId: idWiadomosci, Stopped: biegnie}, nil
}

// PrzerwijTure przerywa turę okna i mówi, czy jakaś biegła; ma dać biegnącemu procesowi sygnał, że ma się domknąć. Wpis schodzi od razu, żeby okno przyjęło kolejną wiadomość, zanim tura przerwana się domknie.
func (a *adapterRozmowy) PrzerwijTure(idOkna string) bool {
	a.mu.Lock()
	bieg, biegnie := a.biegnace[idOkna]
	delete(a.biegnace, idOkna)
	a.mu.Unlock()
	if biegnie {
		bieg.anuluj()
	}
	return biegnie
}

// zapomnijBieg usuwa zakończoną turę z wykazu biegnących, jeżeli wpis okna nadal jest jej własny. Wpis cudzy należy do tury zajętej po zatrzymaniu i zostaje.
func (a *adapterRozmowy) zapomnijBieg(idOkna string, bieg *biegTury) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if zastany, jest := a.biegnace[idOkna]; jest && zastany == bieg {
		delete(a.biegnace, idOkna)
	}
}
