package transport

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/coder/websocket"

	"danacoconsole/server/internal/protocol"
)

// bladPolaczenieZamkniete jest błędem pojedynczej wysyłki, nie stanem sesji.
var bladPolaczenieZamkniete = errors.New("transport: połączenie zamknięte")

// Porzucona odpowiedź zostawia obietnicę bez rozstrzygnięcia: gniazdo idzie do zamknięcia.
var bladKolejkaPelna = errors.New("transport: kolejka wyjściowa pełna")

const (
	czasNaMiejsceWKolejce = 2 * time.Second
	czasNaDosylke         = time.Second
	czasNaZamkniecie      = 500 * time.Millisecond
)

type Polaczenie struct {
	id       string
	gniazdo  *websocket.Conn
	wyjscie  chan []byte
	kontekst context.Context
	zakoncz  context.CancelFunc
	dziennik *log.Logger

	// Semafor biegów obsługi: pojemność kanału jest granicą biegów naraz.
	biegi chan struct{}

	zamek     sync.RWMutex
	konto     string
	tozsamosc Tozsamosc
	zamkniete bool
	// Późniejsza odmowa bramki znaczy wtedy sesję, która przestała nadawać.
	przeszlaBramke bool
}

func nowePolaczenie(rodzic context.Context, id, konto string, tozsamosc Tozsamosc, gniazdo *websocket.Conn, pojemnosc int, dziennik *log.Logger) *Polaczenie {
	kontekst, zakoncz := context.WithCancel(rodzic)
	return &Polaczenie{
		id:        id,
		gniazdo:   gniazdo,
		wyjscie:   make(chan []byte, pojemnosc),
		biegi:     make(chan struct{}, biegiNaPolaczenie),
		kontekst:  kontekst,
		zakoncz:   zakoncz,
		dziennik:  dziennik,
		konto:     konto,
		tozsamosc: tozsamosc,
	}
}

func (p *Polaczenie) Id() string {
	return p.id
}

func (p *Polaczenie) Konto() string {
	p.zamek.RLock()
	defer p.zamek.RUnlock()
	return p.konto
}

func (p *Polaczenie) PrzypiszKonto(konto string) {
	if konto == "" {
		return
	}
	p.zamek.Lock()
	p.konto = konto
	p.zamek.Unlock()
}

func (p *Polaczenie) Tozsamosc() Tozsamosc {
	p.zamek.RLock()
	defer p.zamek.RUnlock()
	return p.tozsamosc
}

func (p *Polaczenie) PrzedstawKlienta(id string) {
	if id == "" {
		return
	}
	p.zamek.Lock()
	p.tozsamosc.IdKlienta = id
	p.zamek.Unlock()
}

// Czekanie na miejsce dotyczy odpowiedzi i fragmentów domykających strumień.
func (p *Polaczenie) Wyslij(k protocol.Koperta) error {
	dane, err := protocol.Zakoduj(k)
	if err != nil {
		return err
	}
	if err := p.wyslijZCzekaniem(dane); err != nil {
		if errors.Is(err, bladKolejkaPelna) {
			p.ZamknijKodem(websocket.StatusTryAgainLater,
				"urządzenie nie odbiera — kolejka wyjściowa pełna przez "+czasNaMiejsceWKolejce.String())
		}
		return err
	}
	return nil
}

func (p *Polaczenie) wyslijZCzekaniem(dane []byte) error {
	if err := p.wyslijBajty(dane); !errors.Is(err, bladKolejkaPelna) {
		return err
	}
	czekanie := time.NewTimer(czasNaMiejsceWKolejce)
	defer czekanie.Stop()
	select {
	case p.wyjscie <- dane:
		return nil
	case <-p.kontekst.Done():
		return bladPolaczenieZamkniete
	case <-czekanie.C:
		return bladKolejkaPelna
	}
}

func (p *Polaczenie) wyslijBajty(dane []byte) error {
	p.zamek.RLock()
	zamkniete := p.zamkniete
	p.zamek.RUnlock()
	if zamkniete {
		return bladPolaczenieZamkniete
	}
	select {
	case p.wyjscie <- dane:
		return nil
	case <-p.kontekst.Done():
		return bladPolaczenieZamkniete
	default:
		return bladKolejkaPelna
	}
}

func (p *Polaczenie) Zamknij(powod string) {
	p.ZamknijKodem(websocket.StatusNormalClosure, powod)
}

// Kod zamknięcia odróżnia zwykły koniec pracy od odmowy; powód tekstowy nie.
func (p *Polaczenie) ZamknijKodem(kod websocket.StatusCode, powod string) {
	p.zamek.Lock()
	if p.zamkniete {
		p.zamek.Unlock()
		return
	}
	p.zamkniete = true
	p.zamek.Unlock()

	// Anulowanie kontekstu zamyka gniazdo TCP, a urządzenie dostałoby EOF zamiast kodu.
	umowione := make(chan error, 1)
	go func() { umowione <- p.gniazdo.Close(kod, powod) }()
	select {
	case err := <-umowione:
		if err != nil {
			_ = p.gniazdo.CloseNow()
		}
	case <-time.After(czasNaZamkniecie):
		_ = p.gniazdo.CloseNow()
	}
	p.zakoncz()
	p.dziennik.Printf("transport: rozłączenie %s (%s)", p.id, powod)
}

// Ramka przed zamknięciem omija kolejkę: zamknięcie kończy pętlę wysyłki.
func (p *Polaczenie) wyslijWprost(k protocol.Koperta) error {
	dane, err := protocol.Zakoduj(k)
	if err != nil {
		return err
	}
	kontekst, koniec := context.WithTimeout(p.kontekst, czasNaDosylke)
	defer koniec()
	return p.gniazdo.Write(kontekst, websocket.MessageText, dane)
}

func (p *Polaczenie) oznaczPrzejscieBramki() {
	p.zamek.Lock()
	p.przeszlaBramke = true
	p.zamek.Unlock()
}

func (p *Polaczenie) przeszlaPrzezBramke() bool {
	p.zamek.RLock()
	defer p.zamek.RUnlock()
	return p.przeszlaBramke
}
