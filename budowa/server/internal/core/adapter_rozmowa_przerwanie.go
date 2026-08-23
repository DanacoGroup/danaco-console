// Odpowiedzialność pliku: przerwanie tury jest jawne. Wysłanie wiadomości do
// okna, które właśnie odpowiada, odmawia zamiast anulować turę po cichu.
// Kontrakt zna `message.stop` jako zatrzymanie odpowiedzi; klient, który chce
// przerwać i wysłać od razu, wysyła parę: najpierw `message.stop`, potem
// `message.send`. Tu też mieszka wykaz tur w biegu (`biegnace`) i wszystko,
// co go dotyka.
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

// zajmijBieg zajmuje okno pod nową turę i oddaje prawdę, gdy się to udało.
// Sprawdzenie i zajęcie idą pod jednym zamkiem, nie dwoma wywołaniami: odczyt
// osobny od zapisu zostawiłby szczelinę, w której dwie wiadomości nadane w tej
// samej chwili obie zobaczyłyby okno wolne i obie ruszyłyby turę. Okno zajęte
// nie jest przerywane — wywołujący dostaje fałsz i odmawia (odmowaTuryWBiegu
// wyżej).
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

// PrzerwijTure przerywa turę okna i mówi, czy jakaś biegła.
//
// Wydzielone z Zatrzymaj, bo zamykanie okna nie jest przyciskiem Operatora:
// nie dotyczy pętli naprawczej i nie odpowiada kontraktem. Ma zrobić jedną
// rzecz — dać biegnącemu procesowi sygnał, że ma się domknąć.
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

// zapomnijBieg usuwa zakończoną turę z wykazu biegnących.
func (a *adapterRozmowy) zapomnijBieg(idOkna string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.biegnace, idOkna)
}
