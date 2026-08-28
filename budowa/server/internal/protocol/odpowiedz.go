package protocol

import (
	"encoding/json"
	"fmt"

	"danacoconsole/shared"
)

// Status rozstrzyga, czy odpowiedź niesie wynik czy błąd. Wartości wyznacza
// wyliczenie EnvelopeStatus kontraktu.
type Status = shared.EnvelopeStatus

// Odpowiedz jest wynikiem wykonania komendy przed zapakowaniem w kopertę:
// wykonawca komendy zwraca ją warstwie transportu, bez znaczników JSON.
type Odpowiedz struct {
	Status Status
	Wynik  json.RawMessage
	Blad   *Blad
}

// Sukces buduje odpowiedź udaną z podanego wyniku wywołania; wynik nil
// oznacza potwierdzenie bez treści.
func Sukces(wynik any) (Odpowiedz, error) {
	o := Odpowiedz{Status: shared.EnvelopeStatusOk}
	if wynik == nil {
		return o, nil
	}
	surowy, err := json.Marshal(wynik)
	if err != nil {
		return Odpowiedz{}, fmt.Errorf("protocol: kodowanie wyniku: %w", err)
	}
	o.Wynik = surowy
	return o, nil
}

// Porazka buduje odpowiedź błędną wprost z gotowego błędu tego kontraktu
// przekazanego przez wywołującego.
func Porazka(b Blad) Odpowiedz {
	return Odpowiedz{Status: shared.EnvelopeStatusError, Blad: &b}
}

// PorazkaKodem buduje odpowiedź błędną wprost z kodu kontraktu i komunikatu —
// skrót dla najczęstszego przypadku.
func PorazkaKodem(kod KodBledu, komunikat string) Odpowiedz {
	return Porazka(NowyBlad(kod, komunikat))
}

// KopertaOdpowiedzi pakuje odpowiedź w kopertę zwrotną: status i błąd trafiają
// do pól odpowiedzi, wynik do ładunku koperty.
func KopertaOdpowiedzi(zadanie Koperta, o Odpowiedz) Koperta {
	return Koperta{
		Type:      zadanie.Type,
		Id:        zadanie.Id,
		SessionId: zadanie.SessionId,
		Payload:   o.Wynik,
		Timestamp: Teraz(),
		Status:    wskaznik(o.Status),
		Error:     o.Blad,
	}
}

// KopertaBledu pakuje dany błąd w kopertę zwrotną, budując odpowiedź błędną
// z tego samego błędu wprost.
func KopertaBledu(zadanie Koperta, b Blad) Koperta {
	return KopertaOdpowiedzi(zadanie, Porazka(b))
}
