package protocol

import (
	"encoding/json"
	"fmt"

	"danacoconsole/shared"
)

// Zasieg niesie kontekst wykonania komendy: okno komunikacji jest bytem
// pośrednim między sesją a wiadomością, najwęższym poziomem zasięgu
// konfiguracji.
type Zasieg struct {
	Srodowisko string `json:"environmentId,omitempty"`
	Projekt    string `json:"projectId,omitempty"`
	Sesja      string `json:"sessionId,omitempty"`
	Okno       string `json:"windowId,omitempty"`
}

// Request jest komendą przygotowaną do wykonania: nazwą, zasięgiem i ładunkiem.
// Znana odróżnia komendę z kontraktu od nazwy zastępczej `*.unknown`;
// TypZadany zachowuje typ, który rzeczywiście przyszedł od klienta.
type Request struct {
	Komenda       shared.MessageType
	TypZadany     shared.MessageType
	Id            string
	Znana         bool
	Zasieg        Zasieg
	Ladunek       json.RawMessage
	ZnacznikCzasu int64
}

// trescNieznanejKomendy jest ładunkiem zdarzenia unknown w kształcie
// kontraktu, niosącym typ nierozpoznany.
type trescNieznanejKomendy = shared.UnknownCommandPayload

// powodNieznanej wyjaśnia klientowi, dlaczego typ nie został rozpoznany.
// Kontrakt zostawia pole reason wolnym opisem, więc treść należy do rdzenia.
const powodNieznanej = "typ spoza kontraktu rdzenia"

// zasiegZLadunku wyjmuje pola zasięgu z ładunku komendy. Ładunek bez zasięgu
// ani ładunek innego kształtu nie są błędem — wraca zasięg pusty.
func zasiegZLadunku(ladunek json.RawMessage) Zasieg {
	if len(ladunek) == 0 {
		return Zasieg{}
	}
	var z Zasieg
	if err := json.Unmarshal(ladunek, &z); err != nil {
		return Zasieg{}
	}
	return z
}

// ZbudujZadanie zamienia kopertę na żądanie gotowe do skierowania.
// Nieznana komenda nie przerywa przetwarzania — dostaje nazwę `*.unknown`
// z kontraktu i znacznik Znana równy false.
func ZbudujZadanie(k Koperta, rejestr *RejestrKomend) Request {
	komenda, znana := rejestr.Rozpoznaj(k.Type)
	zasieg := zasiegZLadunku(k.Payload)
	if zasieg.Sesja == "" {
		zasieg.Sesja = IdSesji(k)
	}
	return Request{
		Komenda:       komenda,
		TypZadany:     k.Type,
		Id:            k.Id,
		Znana:         znana,
		Zasieg:        zasieg,
		Ladunek:       k.Payload,
		ZnacznikCzasu: k.Timestamp,
	}
}

// OdkodujZadanie łączy dekodowanie koperty z rozpoznaniem komendy — jest
// wejściem warstwy transportu. Błąd oznacza wyłącznie komunikat niepoprawny
// strukturalnie; nieznana komenda błędem nie jest.
func OdkodujZadanie(dane []byte, rejestr *RejestrKomend) (Request, error) {
	k, err := Odkoduj(dane)
	if err != nil {
		return Request{}, err
	}
	return ZbudujZadanie(k, rejestr), nil
}

// Koperta odtwarza kopertę żądania — do zbudowania odpowiedzi na komendę,
// która przyszła już rozłożona na części.
func (r Request) Koperta() Koperta {
	return Koperta{
		Type:      r.Komenda,
		Id:        r.Id,
		SessionId: wskaznikTekstu(r.Zasieg.Sesja),
		Payload:   r.Ladunek,
		Timestamp: r.ZnacznikCzasu,
	}
}

// LadunekDo rozpakowuje ładunek tego żądania wprost do struktury komendy
// pochodzącej z tego samego kontraktu.
func (r Request) LadunekDo(cel any) error {
	if len(r.Ladunek) == 0 {
		return nil
	}
	if err := json.Unmarshal(r.Ladunek, cel); err != nil {
		return fmt.Errorf("protocol: dekodowanie ładunku %s: %w", r.Komenda, err)
	}
	return nil
}

// OdpowiedzNieznanej buduje odpowiedź na komendę spoza kontraktu: zdarzenie
// `<obszar>.unknown` z typem, którego rdzeń nie rozpoznał. Nie jest to błąd —
// połączenie nie jest zrywane, sesja nie jest blokowana, kolejne żądania są
// przyjmowane.
func OdpowiedzNieznanej(r Request) Koperta {
	tresc := trescNieznanejKomendy{
		RequestedType: r.TypZadany,
		RequestId:     wskaznikTekstu(r.Id),
		Reason:        wskaznikTekstu(powodNieznanej),
	}
	k, err := NowaKoperta(r.Komenda, r.Id, r.Zasieg.Sesja, tresc)
	if err != nil {
		return zeStanemOdmowy(Koperta{
			Type:      r.Komenda,
			Id:        r.Id,
			SessionId: wskaznikTekstu(r.Zasieg.Sesja),
			Timestamp: Teraz(),
		}, r.TypZadany)
	}
	return zeStanemOdmowy(k, r.TypZadany)
}

// zeStanemOdmowy dokłada kopercie odmowy stan i błąd: ten sam identyfikator
// żądania, stan error i kod not_found.
func zeStanemOdmowy(k Koperta, typZadany shared.MessageType) Koperta {
	stan := shared.EnvelopeStatus(shared.EnvelopeStatusError)
	k.Status = &stan
	blad := NowyBlad(shared.ErrorCodeNotFound,
		"rdzeń nie ma uchwytu komendy "+string(typZadany))
	k.Error = &blad
	return k
}
