package transport

import (
	"context"
	"fmt"
	"log"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// wykonajBezpiecznie oddaje żądanie rdzeniowi i pilnuje, by pojedyncze
// wywołanie nie mogło unieruchomić procesu.
//
// Dwa zachowania fail-open:
//   - rdzeń niepodłączony — żądanie dostaje zdarzenie `*.unknown` z kontraktu,
//     połączenie żyje dalej;
//   - załamanie obsługi — wraca odpowiedź z kodem internal_error zamiast
//     przerwania procesu; kolejne żądania są przyjmowane.
//
// Straż bramki (`bramka.go`) nie jest fail-open i ma pierwszeństwo przed
// obydwoma: poza pętlą zwrotną żądanie z gniazda, które nie przedstawiło tokenu,
// nie dochodzi do rdzenia — wraca odmowa `not_authenticated`. Na pętli zwrotnej
// straż jest wyłączona.
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
	// Straż stoi tutaj, a nie w pętli odbioru, bo to jedyne miejsce, przez które
	// żądanie z gniazda przechodzi do rdzenia.
	if !straz.przepusc(rdzen, zadanie.Komenda, ujscie) {
		dziennik.Printf("transport: komenda %s bez przejścia przez bramkę — odmowa", zadanie.Komenda)
		return odmowaBezBramki(zadanie)
	}
	return rdzen.Obsluz(kontekst, zadanie, ujscie)
}

// kopertaBleduStruktury odpowiada na komunikat, którego nie da się odczytać jako
// koperty kontraktu. Nie jest to nieznana komenda — tam wraca `*.unknown` — lecz
// komunikat niepoprawny strukturalnie, więc odpowiedź niesie kod
// validation_failed. Połączenie pozostaje otwarte.
func kopertaBleduStruktury(err error) protocol.Koperta {
	pusta := protocol.Koperta{Type: shared.EventConnectionUnknown, Timestamp: protocol.Teraz()}
	blad := protocol.NowyBlad(shared.ErrorCodeValidationFailed, "komunikat niezgodny z kopertą kontraktu: "+err.Error())
	return protocol.KopertaBledu(pusta, blad)
}
