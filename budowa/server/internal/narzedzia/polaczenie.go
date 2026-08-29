// Gniazdo do rdzenia Danaco Console jest dla rdzenia zwykłym urządzeniem,
// łączącym się tym samym gniazdem WebSocket i kopertą kontraktu, co okno
// interfejsu; zestawia się leniwie, przy pierwszym wywołaniu.
package narzedzia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/coder/websocket"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/transport"
)

// Polaczenie prowadzi jedno gniazdo do rdzenia. Wywołania narzędzi idą po nim
// pojedynczo — protokół MCP po stdio przyjmuje żądania wierszami, a rozdzielnia
// obsługuje je po kolei.
type Polaczenie struct {
	adres   string
	zamek   sync.Mutex
	gniazdo *websocket.Conn
}

// Polacz przygotowuje gniazdo pod wskazanym adresem, nie nawiązując go,
// i dokłada do adresu parametry tożsamości opisujące klienta wobec rdzenia.
func Polacz(adres, okno string, zasieg Zasieg) *Polaczenie {
	return &Polaczenie{adres: zAdresemTozsamosci(adres, okno, zasieg)}
}

// zAdresemTozsamosci dokłada do adresu gniazda parametry tożsamości; adres
// nieczytelny zostaje adresem dotychczasowym, bez tożsamości.
func zAdresemTozsamosci(adres, okno string, zasieg Zasieg) string {
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
	cel.RawQuery = zapytanie.Encode()
	return cel.String()
}

// Wykonaj wysyła żądanie i czeka na odpowiedź o tym samym identyfikatorze,
// rozpoznawaną też po nazwie komendy i obecności pola stanu.
func (p *Polaczenie) Wykonaj(kontekst context.Context, zadanie protocol.Koperta) (protocol.Koperta, error) {
	p.zamek.Lock()
	defer p.zamek.Unlock()

	gniazdo, err := p.nawiazane(kontekst)
	if err != nil {
		return protocol.Koperta{}, err
	}
	bajty, err := protocol.Zakoduj(zadanie)
	if err != nil {
		return protocol.Koperta{}, err
	}
	if err := gniazdo.Write(kontekst, websocket.MessageText, bajty); err != nil {
		return protocol.Koperta{}, p.zerwane("zapis żądania", err)
	}
	return p.czekajNaOdpowiedz(kontekst, gniazdo, zadanie)
}

// Zamknij kończy gniazdo połączenia z rdzeniem; wywołanie na połączeniu
// nienawiązanym nic nie robi tutaj.
func (p *Polaczenie) Zamknij() {
	p.zamek.Lock()
	defer p.zamek.Unlock()
	if p.gniazdo == nil {
		return
	}
	p.gniazdo.Close(websocket.StatusNormalClosure, "koniec pracy serwera narzędzi")
	p.gniazdo = nil
}

// czekajNaOdpowiedz czyta ramki gniazda aż do odpowiedzi na to jedno żądanie,
// pomijając rozgłoszenia rdzenia po drodze.
func (p *Polaczenie) czekajNaOdpowiedz(kontekst context.Context, gniazdo *websocket.Conn,
	zadanie protocol.Koperta) (protocol.Koperta, error) {

	for {
		_, bajty, err := gniazdo.Read(kontekst)
		if err != nil {
			return protocol.Koperta{}, p.zerwane("odczyt odpowiedzi", err)
		}
		var koperta protocol.Koperta
		if err := json.Unmarshal(bajty, &koperta); err != nil {
			continue // ramka nie w kształcie koperty — nie jest odpowiedzią na to żądanie
		}
		if koperta.Id == zadanie.Id && koperta.Type == zadanie.Type && koperta.Status != nil {
			return koperta, nil
		}
	}
}

// nawiazane zwraca gniazdo czynne połączenia z rdzeniem, zestawiając je przy
// pierwszym użyciu tego połączenia.
func (p *Polaczenie) nawiazane(kontekst context.Context) (*websocket.Conn, error) {
	if p.gniazdo != nil {
		return p.gniazdo, nil
	}
	gniazdo, _, err := websocket.Dial(kontekst, p.adres, nil)
	if err != nil {
		return nil, fmt.Errorf("serwer pod adresem %s nie odpowiada: %w", p.adres, err)
	}
	gniazdo.SetReadLimit(transport.LimitOdczytu)
	p.gniazdo = gniazdo
	return gniazdo, nil
}

// zerwane odkłada gniazdo zerwane, żeby wywołanie następne zestawiło je od nowa,
// i zwraca czytelny opis dla modelu.
func (p *Polaczenie) zerwane(czynnosc string, err error) error {
	if p.gniazdo != nil {
		p.gniazdo.CloseNow()
		p.gniazdo = nil
	}
	return fmt.Errorf("połączenie z serwerem zerwane przy czynności %q: %w", czynnosc, err)
}
