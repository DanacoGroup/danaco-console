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
	// Gniazdo bez nikogo po drugiej stronie nie zgłasza tego odczytem.
	odstepPingu          = 20 * time.Second
	czasNaOdpowiedzPingu = 10 * time.Second
	// Bieg na każdy komunikat rośnie bez granicy; zajęta wstrzymuje odczyt.
	biegiNaPolaczenie = 32
)

// petlaPingu biegnie obok pętli odbioru, bo odpowiedź na ping czyta właśnie ta pętla.
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

func (p *Polaczenie) petlaOdbioru(kontekstRdzenia context.Context, zrodloRdzenia func() Rdzen, rejestr *protocol.RejestrKomend, praca *sync.WaitGroup, dopuszczenie dopuszczenieBramki) {
	if dopuszczenie.wymagana && !p.przeszlaPrzezBramke() {
		p.gniazdo.SetReadLimit(LimitOdczytuPrzedBramka)
	}
	for {
		_, dane, err := p.gniazdo.Read(p.kontekst)
		if err != nil {
			p.Zamknij(powodRozlaczenia(err))
			return
		}
		p.przyjmij(kontekstRdzenia, zrodloRdzenia, rejestr, praca, dopuszczenie, dane)
	}
}

func powodRozlaczenia(err error) string {
	if status := websocket.CloseStatus(err); status != -1 {
		return fmt.Sprintf("kanał zamknięty przez urządzenie (kod %d)", status)
	}
	if errors.Is(err, context.Canceled) {
		return "kanał zamknięty przez serwer (zatrzymanie)"
	}
	return "odczyt przerwany: " + err.Error()
}

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
	rdzen := zrodloRdzenia()
	if rdzen != nil && !p.dopuscZadanie(rdzen, dopuszczenie, zadanie) {
		return
	}
	if !p.zajmijBieg() {
		return
	}
	praca.Add(1)
	go func() {
		defer praca.Done()
		defer p.zwolnijBieg()
		odpowiedz := wykonajBezpiecznie(kontekstRdzenia, rdzen, zadanie, p, dopuszczenie, p.dziennik)
		if odpowiedz.Type == "" {
			return
		}
		if err := p.Wyslij(odpowiedz); err != nil {
			p.dziennik.Printf("transport: odpowiedź %s do %s nieodesłana: %v", zadanie.Komenda, p.id, err)
		}
	}()
}

// zajmijBieg czeka w pętli odbioru celowo: wstrzymany odczyt dławi nadawcę.
func (p *Polaczenie) zajmijBieg() bool {
	select {
	case p.biegi <- struct{}{}:
		return true
	case <-p.kontekst.Done():
		return false
	}
}

func (p *Polaczenie) zwolnijBieg() {
	<-p.biegi
}

// dopuscZadanie rozstrzyga bramkę przed powołaniem biegu obsługi.
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
	// Odmowa po komendach spoza wejścia znaczy sesję unieważnioną.
	if err := p.wyslijWprost(odmowaBezBramki(zadanie)); err != nil {
		p.dziennik.Printf("transport: odmowa bramki do %s nieodesłana: %v", p.id, err)
	}
	p.ZamknijKodem(websocket.StatusPolicyViolation, "sesja bramki przestała nadawać")
	return false
}
