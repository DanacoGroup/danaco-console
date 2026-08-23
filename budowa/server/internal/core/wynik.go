package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// wynik zamienia rezultat czynności domeny na odpowiedź protokołu.
// Wynik nieserializowalny nie jest odmową wykonania: czynność już się wykonała,
// więc klient dostaje błąd wewnętrzny, a nie ciszę.
func wynik(v any) protocol.Odpowiedz {
	o, err := protocol.Sukces(v)
	if err != nil {
		return protocol.PorazkaKodem(shared.ErrorCodeInternalError, err.Error())
	}
	return o
}

// porazka zamienia błąd czynności domeny na odpowiedź błędną.
//
// Kod błędu wybiera warstwa, która zna przyczynę: pakiet domeny zwraca błąd
// niosący kod kontraktu przez protocol.JakoError, a rdzeń go przenosi bez
// tłumaczenia. Błąd bez kodu jest błędem wewnętrznym rdzenia — rdzeń nie
// zgaduje za domenę, czy to brak bytu, konflikt, czy usterka.
//
// Zerwany kontekst wywołania jest osobnym przypadkiem: żądanie nie zostało
// wykonane, ale ponowienie ma sens, więc idzie kodem ponawialnym.
func porazka(err error) protocol.Odpowiedz {
	if err == nil {
		return protocol.PorazkaKodem(shared.ErrorCodeInternalError, "rdzeń: czynność nie zwróciła ani wyniku, ani przyczyny")
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return protocol.PorazkaKodem(shared.ErrorCodeChannelUnavailable, err.Error())
	}
	return protocol.Porazka(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladNiepoprawnegoLadunku odpowiada na ładunek, którego nie da się odczytać
// w kształcie wyznaczonym przez kontrakt.
func bladNiepoprawnegoLadunku(err error) protocol.Odpowiedz {
	return protocol.PorazkaKodem(shared.ErrorCodeValidationFailed, err.Error())
}
