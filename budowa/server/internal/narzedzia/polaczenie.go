// Gniazdo do rdzenia niesie tę samą kopertę kontraktu, co okno interfejsu.
package narzedzia

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/coder/websocket"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/transport"
	"danacoconsole/shared"
)

// Protokół MCP po stdio przyjmuje żądania wierszami — wywołania idą pojedynczo.
type Polaczenie struct {
	adres  string
	zamek  sync.Mutex
	czynne *czynneGniazdo
}

// Wskazanie oczekiwanego żądania ma własny zamek: bieg odczytu nie może brać
// zamka połączenia, bo trzyma go wywołanie czekające na ten bieg.
type czynneGniazdo struct {
	gniazdo           *websocket.Conn
	koperty           chan protocol.Koperta
	blad              chan error
	zatrzymaj         context.CancelFunc
	zamekOczekiwania  sync.Mutex
	oczekiwaneZadanie string
	oczekiwanaKomenda shared.MessageType
}

func Polacz(adres, okno string, zasieg Zasieg, poswiadczenie string) *Polaczenie {
	return &Polaczenie{adres: zAdresemTozsamosci(adres, okno, zasieg, poswiadczenie)}
}

// Poświadczenie jedzie osobnym parametrem, bo tożsamość go nie niesie.
func zAdresemTozsamosci(adres, okno string, zasieg Zasieg, poswiadczenie string) string {
	cel, err := url.Parse(adres)
	if err != nil {
		return adres
	}
	zapytanie := cel.Query()
	for nazwa, wartosc := range transport.ParametryTozsamosci(transport.Tozsamosc{
		IdKlienta: NazwaBinarium,
		Rodzaj:    transport.RodzajNarzedzi,
		Zasieg:    string(zasieg),
		IdOkna:    okno,
	}) {
		zapytanie.Set(nazwa, wartosc)
	}
	if poswiadczenie != "" {
		zapytanie.Set(transport.ParametrPoswiadczenia, poswiadczenie)
	}
	cel.RawQuery = zapytanie.Encode()
	return cel.String()
}

// Żądanie z gniazda zastanego idzie raz jeszcze gniazdem zestawionym od nowa.
func (p *Polaczenie) Wykonaj(kontekst context.Context, zadanie protocol.Koperta) (protocol.Koperta, error) {
	p.zamek.Lock()
	defer p.zamek.Unlock()

	bajty, err := protocol.Zakoduj(zadanie)
	if err != nil {
		return protocol.Koperta{}, err
	}
	odpowiedz, zastane, err := p.podejscie(kontekst, zadanie, bajty)
	if err != nil && zastane && ponowienieMaSens(kontekst, err) {
		odpowiedz, _, err = p.podejscie(kontekst, zadanie, bajty)
	}
	return odpowiedz, err
}

func (p *Polaczenie) Zamknij() {
	p.zamek.Lock()
	defer p.zamek.Unlock()
	if p.czynne == nil {
		return
	}
	// Pożegnanie czyta ten sam bieg, więc idzie przed jego zatrzymaniem.
	p.czynne.gniazdo.Close(websocket.StatusNormalClosure, "koniec pracy serwera narzędzi")
	p.czynne.zatrzymaj()
	p.czynne = nil
}

// Drugi wynik mówi, czy szło gniazdem zastanym.
func (p *Polaczenie) podejscie(kontekst context.Context, zadanie protocol.Koperta,
	bajty []byte) (protocol.Koperta, bool, error) {

	zastane := p.czynne != nil
	czynne, err := p.nawiazane(kontekst)
	if err != nil {
		return protocol.Koperta{}, zastane, err
	}
	czynne.czekaNa(zadanie)
	defer czynne.przestanCzekac()
	if err := czynne.gniazdo.Write(kontekst, websocket.MessageText, bajty); err != nil {
		return protocol.Koperta{}, zastane, p.zerwane("zapis żądania", err)
	}
	odpowiedz, err := p.czekajNaOdpowiedz(kontekst, czynne, zadanie)
	return odpowiedz, zastane, err
}

// Kontekst wywołującego nie jest stanem gniazda.
func ponowienieMaSens(kontekst context.Context, err error) bool {
	return kontekst.Err() == nil &&
		!errors.Is(err, context.Canceled) &&
		!errors.Is(err, context.DeadlineExceeded)
}

// Kontekst wygasły kończy tutaj samo wywołanie; gniazdo odkłada dopiero zerwane,
// wołane przy zapisie żądania i przy błędzie odczytu.
func (p *Polaczenie) czekajNaOdpowiedz(kontekst context.Context, czynne *czynneGniazdo,
	zadanie protocol.Koperta) (protocol.Koperta, error) {

	for {
		// Koperta doręczona wraca także wtedy, gdy kontekst właśnie wygasł.
		select {
		case koperta := <-czynne.koperty:
			if jestOdpowiedzia(koperta, zadanie) {
				return koperta, nil
			}
			continue
		default:
		}
		select {
		case <-kontekst.Done():
			return protocol.Koperta{}, fmt.Errorf("oczekiwanie na odpowiedź przerwane: %w", kontekst.Err())
		case err := <-czynne.blad:
			return protocol.Koperta{}, p.zerwane("odczyt odpowiedzi", err)
		case koperta := <-czynne.koperty:
			if jestOdpowiedzia(koperta, zadanie) {
				return koperta, nil
			}
		}
	}
}

// Pole stanu niesie odpowiedź, a żądanie go nie ma.
func jestOdpowiedzia(koperta, zadanie protocol.Koperta) bool {
	return koperta.Id == zadanie.Id && koperta.Type == zadanie.Type && koperta.Status != nil
}

func (p *Polaczenie) nawiazane(kontekst context.Context) (*czynneGniazdo, error) {
	if p.czynne != nil {
		return p.czynne, nil
	}
	gniazdo, _, err := websocket.Dial(kontekst, p.adres, nil)
	if err != nil {
		return nil, fmt.Errorf("serwer pod adresem %s nie odpowiada: %w", p.adres, err)
	}
	gniazdo.SetReadLimit(transport.LimitOdczytu)
	// Bieg odczytu żyje tak długo jak gniazdo, więc nie bierze kontekstu wywołania.
	zycie, zatrzymaj := context.WithCancel(context.Background())
	czynne := &czynneGniazdo{
		gniazdo: gniazdo,
		// Oczekiwane żądanie jest jedno naraz, więc kanał niesie jedną odpowiedź.
		koperty:   make(chan protocol.Koperta, 1),
		blad:      make(chan error, 1),
		zatrzymaj: zatrzymaj,
	}
	go czynne.czytaj(zycie)
	p.czynne = czynne
	return czynne, nil
}

func (p *Polaczenie) zerwane(czynnosc string, err error) error {
	if p.czynne != nil {
		p.czynne.gniazdo.CloseNow()
		p.czynne.zatrzymaj()
		p.czynne = nil
	}
	return fmt.Errorf("połączenie z serwerem zerwane przy czynności %q: %w", czynnosc, err)
}

// Biblioteka odsyła pong wyłącznie z pętli czytającej, a rdzeń zamyka gniazdo
// po pierwszym pingu bez odpowiedzi.
func (c *czynneGniazdo) czytaj(kontekst context.Context) {
	for {
		_, bajty, err := c.gniazdo.Read(kontekst)
		if err != nil {
			// Kanał ma miejsce na ten jeden błąd, więc bieg nie wisi.
			c.blad <- err
			return
		}
		var koperta protocol.Koperta
		if err := json.Unmarshal(bajty, &koperta); err != nil {
			continue // ramka nie w kształcie koperty
		}
		c.podaj(koperta)
	}
}

func (c *czynneGniazdo) czekaNa(zadanie protocol.Koperta) {
	select {
	case <-c.koperty:
	default:
	}
	c.zamekOczekiwania.Lock()
	c.oczekiwaneZadanie = zadanie.Id
	c.oczekiwanaKomenda = zadanie.Type
	c.zamekOczekiwania.Unlock()
}

func (c *czynneGniazdo) przestanCzekac() {
	c.zamekOczekiwania.Lock()
	c.oczekiwaneZadanie = ""
	c.oczekiwanaKomenda = ""
	c.zamekOczekiwania.Unlock()
}

// Rozgłoszenia rdzenia odpadają w miejscu: gniazdo nieczytane nie odpowiada na ping.
func (c *czynneGniazdo) podaj(koperta protocol.Koperta) {
	c.zamekOczekiwania.Lock()
	oczekiwana := c.oczekiwaneZadanie != "" &&
		koperta.Id == c.oczekiwaneZadanie &&
		koperta.Type == c.oczekiwanaKomenda &&
		koperta.Status != nil
	c.zamekOczekiwania.Unlock()
	if !oczekiwana {
		return
	}
	select {
	case c.koperty <- koperta:
	default:
	}
}
