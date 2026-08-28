package transport

import (
	"context"
	"fmt"
	"log"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// wykonajBezpiecznie oddaje żądanie rdzeniowi w sposób odporny na awarię pojedynczego wywołania, zwracając zdarzenie nieznanej komendy przy braku rdzenia i odpowiedź błędu wewnętrznego przy załamaniu obsługi.
func wykonajBezpiecznie(kontekst context.Context, rdzen Rdzen, zadanie protocol.Request, ujscie Ujscie, straz straznikBramki, dziennik *log.Logger) (odpowiedz protocol.Koperta) {
	defer func() {
		if przyczyna := recover(); przyczyna != nil {
			dziennik.Printf("transport: obsługa %s załamana: %v", zadanie.Komenda, przyczyna)
			blad := protocol.NowyBlad(shared.ErrorCodeInternalError, fmt.Sprintf("obsługa komendy przerwana: %v", przyczyna))
			odpowiedz = protocol.KopertaBledu(zadanie.Koperta(), blad)
		}
	}()
	if rdzen == nil {
		return protocol.OdpowiedzNieznanej(zadanie)
	}
	// Straż stoi tutaj, bo to jedyne miejsce, przez które żądanie z gniazda przechodzi do rdzenia.
	if !straz.przepusc(rdzen, zadanie.Komenda, ujscie) {
		dziennik.Printf("transport: komenda %s bez przejścia przez bramkę — odmowa", zadanie.Komenda)
		return odmowaBezBramki(zadanie)
	}
	return rdzen.Obsluz(kontekst, zadanie, ujscie)
}

// kopertaBleduStruktury odpowiada na komunikat, którego nie da się odczytać jako koperty kontraktu, zwracając odpowiedź z kodem błędu walidacji przy zachowaniu otwartego połączenia.
func kopertaBleduStruktury(err error) protocol.Koperta {
	pusta := protocol.Koperta{Type: shared.EventConnectionUnknown, Timestamp: protocol.Teraz()}
	blad := protocol.NowyBlad(shared.ErrorCodeValidationFailed, "komunikat niezgodny z kopertą kontraktu: "+err.Error())
	return protocol.KopertaBledu(pusta, blad)
}
