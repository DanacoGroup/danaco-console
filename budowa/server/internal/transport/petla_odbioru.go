package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/coder/websocket"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// petlaOdbioru czyta komunikaty urządzenia i kieruje je do rdzenia.
//
// kontekstRdzenia jest kontekstem serwera, nie połączenia: rozłączenie klienta
// nie przerywa pracy już rozpoczętej przez rdzeń — sesja, okno i proces biegną
// dalej, a wynik trafi do pozostałych urządzeń konta rozgłoszeniem.
//
// zrodloRdzenia jest dostawcą bieżącej realizacji obsługi komend, a nie
// wartością zamrożoną w chwili nawiązania. Rdzeń podłącza się po utworzeniu
// serwera (PodlaczRdzen); połączenie nawiązane, zanim to nastąpi, musiałoby
// wtedy trzymać nil-rdzeń na całe swoje życie. Pobranie rdzenia dopiero przy
// obsłudze komunikatu sprawia, że komendy z takiego połączenia zaczynają być
// obsługiwane, gdy tylko rdzeń zostanie podłączony.
// straz jest rozstrzygnięciem o wystawieniu nasłuchu ustalonym raz, przy
// normalizacji ustawień (`bramka.go`). Idzie przez pętlę odbioru do wykonania,
// bo tam stoi jedyne wejście żądania do rdzenia.
func (p *Polaczenie) petlaOdbioru(kontekstRdzenia context.Context, zrodloRdzenia func() Rdzen, rejestr *protocol.RejestrKomend, praca *sync.WaitGroup, straz straznikBramki) {
	for {
		_, dane, err := p.gniazdo.Read(p.kontekst)
		if err != nil {
			// Powód bierze się z błędu, nie ze stałej. Odczyt kończy się nie
			// tylko odejściem urządzenia: przekroczenie LimitOdczytu zamyka
			// gniazdo od strony rdzenia (kodem 1009), tak samo błąd ramkowania
			// i zerwanie sieci. Linia dziennika jest jedynym trwałym śladem po
			// takim zdarzeniu, więc ma nazwać, co naprawdę zaszło.
			p.Zamknij(powodRozlaczenia(err))
			return
		}
		p.przyjmij(kontekstRdzenia, zrodloRdzenia, rejestr, praca, straz, dane)
	}
}

// powodRozlaczenia nazywa koniec odczytu i nie zmyśla winnego. Odczyt kończy
// się na cztery sposoby: przekroczenie LimitOdczytu zamyka gniazdo od strony
// rdzenia (biblioteka odsyła kod 1009), zatrzymanie rdzenia zamyka je z woli
// procesu, zerwanie sieci nie jest niczyją decyzją, a odejście urządzenia zamyka
// je od jego strony. Linia dziennika jest jedynym trwałym śladem po rozłączeniu,
// więc niesie to, co naprawdę zaszło, razem ze zdaniem biblioteki jako
// szczegółem.
func powodRozlaczenia(err error) string {
	if status := websocket.CloseStatus(err); status != -1 {
		return fmt.Sprintf("kanał zamknięty przez urządzenie (kod %d)", status)
	}
	if errors.Is(err, context.Canceled) {
		return "kanał zamknięty przez rdzeń (zatrzymanie)"
	}
	return "odczyt przerwany: " + err.Error()
}

// przedstawZPowitania wyjmuje `clientId` z ładunku powitania i dokłada go do
// tożsamości połączenia.
//
// Dlaczego tu, a nie w rdzeniu. Tożsamość jest własnością połączenia i mieszka
// przy nim (`tozsamosc.go`); rdzeń ją czyta, a nie zapisuje. Gdyby zapisywał,
// musiałby dostać ujście do ręki — a ujście do rdzenia świadomie nie idzie
// (`core/adapter_transportu.go`: druga droga wyjścia obok Nadajnika).
// Odczyt jest czysty: transport bierze pole, którego kształt i tak zna
// z kontraktu, i nie rozstrzyga o nim niczego.
//
// Odczyt dzieje się przed oddaniem żądania rdzeniowi, więc zdarzenia rozgłoszone
// przez samo powitanie znają już klienta. Ładunek nieczytelny albo pole puste
// zostawia tożsamość nietkniętą — powitanie ma się udać zawsze.
func (p *Polaczenie) przedstawZPowitania(zadanie protocol.Request) {
	if zadanie.Komenda != shared.CommandConnectionHello {
		return
	}
	var powitanie shared.ConnectionHelloRequest
	if err := json.Unmarshal(zadanie.Ladunek, &powitanie); err != nil {
		return
	}
	p.PrzedstawKlienta(powitanie.ClientId)
}

// przyjmij rozpoznaje komunikat i oddaje go rdzeniowi.
//
// Każde żądanie idzie osobnym biegiem, więc komenda długotrwała (message.send)
// nie zatrzymuje odczytu i komenda przerywająca (message.stop) dociera w trakcie
// jej wykonania. Odpowiedź niesie identyfikator żądania, więc kolejność
// odpowiedzi nie ma znaczenia dla korelacji.
func (p *Polaczenie) przyjmij(kontekstRdzenia context.Context, zrodloRdzenia func() Rdzen, rejestr *protocol.RejestrKomend, praca *sync.WaitGroup, straz straznikBramki, dane []byte) {
	zadanie, err := protocol.OdkodujZadanie(dane, rejestr)
	if err != nil {
		p.dziennik.Printf("transport: komunikat %s odrzucony: %v", p.id, err)
		if err := p.Wyslij(kopertaBleduStruktury(err)); err != nil {
			p.dziennik.Printf("transport: odesłanie błędu do %s nieudane: %v", p.id, err)
		}
		return
	}
	p.przedstawZPowitania(zadanie)
	praca.Add(1)
	go func() {
		defer praca.Done()
		// Rdzeń pobierany dopiero tutaj — w chwili obsługi komunikatu, nie
		// nawiązania połączenia — więc odzwierciedla stan po każdym PodlaczRdzen.
		rdzen := zrodloRdzenia()
		odpowiedz := wykonajBezpiecznie(kontekstRdzenia, rdzen, zadanie, p, straz, p.dziennik)
		if odpowiedz.Type == "" {
			return
		}
		if err := p.Wyslij(odpowiedz); err != nil {
			p.dziennik.Printf("transport: odpowiedź %s do %s nieodesłana: %v", zadanie.Komenda, p.id, err)
		}
	}()
}
