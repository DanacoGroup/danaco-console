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

// Metoda petlaOdbioru czyta komunikaty urządzenia i kieruje je do rdzenia w kontekście serwera, nie połączenia.
func (p *Polaczenie) petlaOdbioru(kontekstRdzenia context.Context, zrodloRdzenia func() Rdzen, rejestr *protocol.RejestrKomend, praca *sync.WaitGroup, straz straznikBramki) {
	for {
		_, dane, err := p.gniazdo.Read(p.kontekst)
		if err != nil {
			// Powód bierze się z błędu, nie ze stałej, bo linia dziennika ma nazwać, co naprawdę zaszło.
			p.Zamknij(powodRozlaczenia(err))
			return
		}
		p.przyjmij(kontekstRdzenia, zrodloRdzenia, rejestr, praca, straz, dane)
	}
}

// Funkcja powodRozlaczenia nazywa koniec odczytu i nie zmyśla winnego, niosąc prawdziwy powód rozłączenia gniazda.
func powodRozlaczenia(err error) string {
	if status := websocket.CloseStatus(err); status != -1 {
		return fmt.Sprintf("kanał zamknięty przez urządzenie (kod %d)", status)
	}
	if errors.Is(err, context.Canceled) {
		return "kanał zamknięty przez rdzeń (zatrzymanie)"
	}
	return "odczyt przerwany: " + err.Error()
}

// Metoda przedstawZPowitania wyjmuje identyfikator klienta z ładunku powitania i dokłada go do tożsamości połączenia.
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

// Metoda przyjmij rozpoznaje komunikat i oddaje go rdzeniowi, uruchamiając obsługę każdego żądania osobnym biegiem.
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
		// Rdzeń pobierany dopiero tutaj, przy obsłudze komunikatu, odzwierciedla stan po podłączeniu rdzenia.
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
