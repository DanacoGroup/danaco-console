package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// odstepPingu jest odstępem między pingami kanału. Gniazdo, po którego
	// drugiej stronie nie ma już nikogo (uśpiony laptop, zerwana sieć), nie
	// zgłasza tego odczytem — bez pingu trzymałoby dwie gorutyny i wpis
	// w rejestrze przez całe życie procesu.
	odstepPingu = 20 * time.Second
	// czasNaOdpowiedzPingu jest granicą czekania na odpowiedź. Przekroczenie
	// znaczy urządzenie nieodpowiadające, nie urządzenie wolne.
	czasNaOdpowiedzPingu = 10 * time.Second
)

// petlaPingu pilnuje, czy po drugiej stronie gniazda ktoś jeszcze jest, i zamyka
// połączenie po pierwszym pingu bez odpowiedzi. Biegnie obok pętli odbioru, bo
// odpowiedź na ping czyta właśnie ta pętla.
func (p *Polaczenie) petlaPingu() {
	tykanie := time.NewTicker(odstepPingu)
	defer tykanie.Stop()
	for {
		select {
		case <-p.kontekst.Done():
			return
		case <-tykanie.C:
			kontekst, koniec := context.WithTimeout(p.kontekst, czasNaOdpowiedzPingu)
			err := p.gniazdo.Ping(kontekst)
			koniec()
			if err != nil {
				p.Zamknij("brak odpowiedzi na ping w " + czasNaOdpowiedzPingu.String())
				return
			}
		}
	}
}

// Metoda petlaOdbioru czyta komunikaty urządzenia i kieruje je do rdzenia w kontekście serwera, nie połączenia.
func (p *Polaczenie) petlaOdbioru(kontekstRdzenia context.Context, zrodloRdzenia func() Rdzen, rejestr *protocol.RejestrKomend, praca *sync.WaitGroup, dopuszczenie dopuszczenieBramki) {
	if dopuszczenie.wymagana && !p.przeszlaPrzezBramke() {
		p.gniazdo.SetReadLimit(LimitOdczytuPrzedBramka)
	}
	for {
		_, dane, err := p.gniazdo.Read(p.kontekst)
		if err != nil {
			// Powód bierze się z błędu, nie ze stałej, bo linia dziennika ma nazwać, co naprawdę zaszło.
			p.Zamknij(powodRozlaczenia(err))
			return
		}
		p.przyjmij(kontekstRdzenia, zrodloRdzenia, rejestr, praca, dopuszczenie, dane)
	}
}

// Funkcja powodRozlaczenia nazywa koniec odczytu i nie zmyśla winnego, niosąc prawdziwy powód rozłączenia gniazda.
func powodRozlaczenia(err error) string {
	if status := websocket.CloseStatus(err); status != -1 {
		return fmt.Sprintf("kanał zamknięty przez urządzenie (kod %d)", status)
	}
	if errors.Is(err, context.Canceled) {
		return "kanał zamknięty przez serwer (zatrzymanie)"
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
func (p *Polaczenie) przyjmij(kontekstRdzenia context.Context, zrodloRdzenia func() Rdzen, rejestr *protocol.RejestrKomend, praca *sync.WaitGroup, dopuszczenie dopuszczenieBramki, dane []byte) {
	zadanie, err := protocol.OdkodujZadanie(dane, rejestr)
	if err != nil {
		p.dziennik.Printf("transport: komunikat %s odrzucony: %v", p.id, err)
		if err := p.Wyslij(kopertaBleduStruktury(err)); err != nil {
			p.dziennik.Printf("transport: odesłanie błędu do %s nieudane: %v", p.id, err)
		}
		return
	}
	p.przedstawZPowitania(zadanie)
	// Rdzeń pobierany dopiero tutaj, przy obsłudze komunikatu, odzwierciedla stan po podłączeniu rdzenia.
	rdzen := zrodloRdzenia()
	if rdzen != nil && !p.dopuscZadanie(rdzen, dopuszczenie, zadanie) {
		return
	}
	praca.Add(1)
	go func() {
		defer praca.Done()
		odpowiedz := wykonajBezpiecznie(kontekstRdzenia, rdzen, zadanie, p, dopuszczenie, p.dziennik)
		if odpowiedz.Type == "" {
			return
		}
		if err := p.Wyslij(odpowiedz); err != nil {
			p.dziennik.Printf("transport: odpowiedź %s do %s nieodesłana: %v", zadanie.Komenda, p.id, err)
		}
	}()
}

// dopuscZadanie rozstrzyga bramkę przed powołaniem biegu obsługi. Bieg powołany
// przed sprawdzeniem jest pracą rdzenia wykonaną na rzecz gniazda, które bramki
// nie przeszło — przy ramce do 16 MiB i biegu na każdy komunikat wystarcza to do
// zajęcia maszyny bez jednego logowania.
func (p *Polaczenie) dopuscZadanie(rdzen Rdzen, dopuszczenie dopuszczenieBramki, zadanie protocol.Request) bool {
	if dopuszczenie.przepusc(rdzen, zadanie.Komenda, p) {
		if _, wejscie := komendyWejscia[zadanie.Komenda]; dopuszczenie.wymagana && !wejscie {
			p.oznaczPrzejscieBramki()
			p.gniazdo.SetReadLimit(LimitOdczytu)
		}
		return true
	}
	p.dziennik.Printf("transport: komenda %s z %s bez przejścia przez bramkę — odmowa", zadanie.Komenda, p.id)
	if !p.przeszlaPrzezBramke() {
		if err := p.Wyslij(odmowaBezBramki(zadanie)); err != nil {
			p.dziennik.Printf("transport: odmowa bramki do %s nieodesłana: %v", p.id, err)
		}
		return false
	}
	// Gniazdo wykonywało już komendy, więc odmowa znaczy sesję unieważnioną albo
	// wygasłą. Dostęp odbiera się wtedy razem z kanałem: sesja zamknięta gdzie
	// indziej ma zerwać to gniazdo, a nie czekać na jego kolejną komendę.
	if err := p.wyslijWprost(odmowaBezBramki(zadanie)); err != nil {
		p.dziennik.Printf("transport: odmowa bramki do %s nieodesłana: %v", p.id, err)
	}
	p.ZamknijKodem(websocket.StatusPolicyViolation, "sesja bramki przestała nadawać")
	return false
}
