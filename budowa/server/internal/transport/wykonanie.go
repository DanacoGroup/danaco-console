package transport

import (
	"context"
	"fmt"
	"log"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// wykonajBezpiecznie oddaje żądanie rdzeniowi w sposób odporny na awarię pojedynczego wywołania, zwracając zdarzenie nieznanej komendy przy braku rdzenia i odpowiedź błędu wewnętrznego przy załamaniu obsługi.
// Bramkę rozstrzyga dopuscZadanie przed powołaniem biegu, więc tutaj żądanie jest już dopuszczone; parametr dopuszczenia zostaje do czasu zdjęcia go z wywołania w petla_odbioru.go.
func wykonajBezpiecznie(kontekst context.Context, rdzen Rdzen, zadanie protocol.Request, ujscie Ujscie, _ dopuszczenieBramki, dziennik *log.Logger) (odpowiedz protocol.Koperta) {
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
	return rdzen.Obsluz(kontekst, zadanie, ujscie)
}

// kopertaBleduStruktury odpowiada na komunikat, którego nie da się odczytać jako koperty kontraktu, zwracając odpowiedź z kodem błędu walidacji przy zachowaniu otwartego połączenia.
func kopertaBleduStruktury(err error) protocol.Koperta {
	pusta := protocol.Koperta{Type: shared.EventConnectionUnknown, Timestamp: protocol.Teraz()}
	blad := protocol.NowyBlad(shared.ErrorCodeValidationFailed, "komunikat niezgodny z kopertą kontraktu: "+err.Error())
	return protocol.KopertaBledu(pusta, blad)
}
