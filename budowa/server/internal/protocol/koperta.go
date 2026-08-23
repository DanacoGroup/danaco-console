// Pakiet protocol dokłada zachowanie do kontraktu komunikatów Danaco Console:
// kodowanie i dekodowanie koperty, rozpoznanie komendy, budowę odpowiedzi
// oraz fragmenty strumienia.
//
// Pakiet nie definiuje ani jednej nazwy, kodu, wartości wyliczenia ani kształtu
// komunikatu. Wszystkie pochodzą z pakietu shared wytworzonego z pliku
// shared/contract.json — jedynego źródła prawdy. Zmiana kontraktu
// przerywa kompilację tego pakietu, zamiast rozjeżdżać się z nim po cichu.
package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"danacoconsole/shared"
)

// Koperta jest jedynym kształtem komunikatu przesyłanego po WebSocket w obie
// strony. Alias, nie kopia: pola i znaczniki JSON wyznacza
// wyłącznie shared.Envelope, więc warstwa protokołu nie ma jak rozejść się
// z kopertą kontraktu.
type Koperta = shared.Envelope

// bladPustegoTypu sygnalizuje kopertę bez pola type — komunikat jest
// niepoprawny strukturalnie, co jest czym innym niż nieznana komenda.
var bladPustegoTypu = errors.New("protocol: koperta bez pola type")

// NowaKoperta składa kopertę i serializuje ładunek. Ładunek nil oznacza
// komunikat bez treści właściwej. Typ podaje wywołujący, biorąc go ze stałych
// pakietu shared.
func NowaKoperta(typ shared.MessageType, id, idSesji string, ladunek any) (Koperta, error) {
	k := Koperta{
		Type:      typ,
		Id:        id,
		SessionId: wskaznikTekstu(idSesji),
		Timestamp: Teraz(),
	}
	if ladunek == nil {
		return k, nil
	}
	surowy, err := json.Marshal(ladunek)
	if err != nil {
		return Koperta{}, fmt.Errorf("protocol: kodowanie ładunku %s: %w", typ, err)
	}
	k.Payload = surowy
	return k, nil
}

// Zakoduj zamienia kopertę na bajty JSON gotowe do wysłania.
func Zakoduj(k Koperta) ([]byte, error) {
	dane, err := json.Marshal(k)
	if err != nil {
		return nil, fmt.Errorf("protocol: kodowanie koperty %s: %w", k.Type, err)
	}
	return dane, nil
}

// Odkoduj czyta kopertę z bajtów JSON. Nie ocenia, czy typ jest znany —
// rozpoznanie komendy należy do rejestru (rozpoznanie.go).
func Odkoduj(dane []byte) (Koperta, error) {
	var k Koperta
	if err := json.Unmarshal(dane, &k); err != nil {
		return Koperta{}, fmt.Errorf("protocol: dekodowanie koperty: %w", err)
	}
	if k.Type == "" {
		return k, bladPustegoTypu
	}
	return k, nil
}

// LadunekDo rozpakowuje ładunek koperty do wskazanej struktury — zwykle do
// typu żądania albo zdarzenia wytworzonego z kontraktu. Pusty ładunek nie jest
// błędem: cel pozostaje nietknięty.
//
// Funkcja, nie metoda: Koperta jest typem kontraktu, więc zachowanie dokłada
// się obok niego, a nie w jego definicji.
func LadunekDo(k Koperta, cel any) error {
	if len(k.Payload) == 0 {
		return nil
	}
	if err := json.Unmarshal(k.Payload, cel); err != nil {
		return fmt.Errorf("protocol: dekodowanie ładunku %s: %w", k.Type, err)
	}
	return nil
}

// IdSesji odczytuje sesję koperty. Kontrakt pozostawia to pole puste
// w powitaniu, więc brak sesji jest wartością, nie błędem.
func IdSesji(k Koperta) string {
	return wartoscTekstu(k.SessionId)
}

// Numer odczytuje numer fragmentu strumienia. Zero oznacza komunikat spoza
// strumienia, bo kontrakt liczy fragmenty od jedynki.
func Numer(k Koperta) int {
	if k.Seq == nil {
		return 0
	}
	return *k.Seq
}

// Ostatni informuje, czy koperta zamyka strumień fragmentów.
func Ostatni(k Koperta) bool {
	return k.Done != nil && *k.Done
}

// Teraz zwraca znacznik czasu protokołu w milisekundach epoki uniksowej —
// jednostkę wyznacza pole timestamp koperty kontraktu.
func Teraz() int64 {
	return time.Now().UnixMilli()
}

// wskaznikTekstu zamienia napis na pole opcjonalne koperty: pusty napis znaczy
// brak pola, dzięki czemu znacznik omitempty kontraktu działa zgodnie z opisem.
func wskaznikTekstu(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

// wartoscTekstu odczytuje pole opcjonalne, zwracając pusty napis dla braku.
func wartoscTekstu(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// wskaznik przenosi wartość do pola opcjonalnego koperty.
func wskaznik[T any](v T) *T {
	return &v
}
