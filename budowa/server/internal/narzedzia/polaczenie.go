// Odpowiedzialność pliku: gniazdo do rdzenia Danaco Console.
//
// Serwer narzędzi jest dla rdzenia zwykłym urządzeniem: łączy się tym samym
// gniazdem WebSocket i tą samą kopertą kontraktu, co okno interfejsu. Drugiego
// wejścia do rdzenia nie ma i nie powstaje tutaj.
//
// Połączenie jest leniwe. Proces modelu uruchamia serwery MCP na starcie
// rozmowy, więc odmowa startu przy niedostępnym rdzeniu zabrałaby modelowi
// wszystkie narzędzia na całą turę. Serwer wstaje zawsze, a gniazdo zestawia
// się przy pierwszym wywołaniu; nieudane zestawienie wraca do modelu treścią
// błędu i nie przeszkadza próbie następnej.
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

// Polacz przygotowuje gniazdo pod wskazanym adresem. Nie nawiązuje go: adres
// bywa nieosiągalny w chwili startu procesu modelu.
//
// Serwer narzędzi przedstawia się przy nawiązaniu, nie powitaniem — bo
// powitania nie wysyła: jest klientem wołającym komendy, a nie oknem
// interfejsu. Adres niesie więc rodzaj klienta, rolę okna i samo okno
// (`transport/tozsamosc.go`), dzięki czemu rdzeń wie, czy po drugiej stronie
// stoi klawiatura Operatora (asystent), czy zwykły model roboczy. Bez tego
// rdzeń widziałby wyłącznie identyfikator gniazda i obu rąk nie odróżniał.
//
// Przedstawienie nie jest uprawnieniem. Rdzeń niczego na nim nie warunkuje —
// wpisuje je do pola opisowego zdarzenia. Zasięg narzędzi rozstrzyga nadal
// wyłącznie przełącznik `--zasieg` czytany przy uruchomieniu z wpisu MCP
// ułożonego przez rdzeń (`zasieg_roli.go`), a nie ten napis.
func Polacz(adres, okno string, zasieg Zasieg) *Polaczenie {
	return &Polaczenie{adres: zAdresemTozsamosci(adres, okno, zasieg)}
}

// zAdresemTozsamosci dokłada do adresu gniazda parametry tożsamości.
//
// Adres nieczytelny zostaje adresem dotychczasowym: gniazdo bez tożsamości jest
// gorsze od gniazda z tożsamością, ale nieporównanie lepsze od braku narzędzi
// przez cały czas życia procesu modelu.
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

// Wykonaj wysyła żądanie i czeka na odpowiedź o tym samym identyfikatorze.
//
// Gniazdo niesie także rozgłoszenia rdzenia (zdarzenia i fragmenty strumienia
// innych okien), dlatego odpowiedź rozpoznaje się po trzech rzeczach naraz:
// identyfikatorze żądania, nazwie komendy i obecności pola stanu. Rozgłoszenie
// nie ma stanu i nosi nazwę zdarzenia, więc nie da się go wziąć za odpowiedź.
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

// Zamknij kończy gniazdo. Wywołanie na połączeniu nienawiązanym nic nie robi.
func (p *Polaczenie) Zamknij() {
	p.zamek.Lock()
	defer p.zamek.Unlock()
	if p.gniazdo == nil {
		return
	}
	p.gniazdo.Close(websocket.StatusNormalClosure, "koniec pracy serwera narzędzi")
	p.gniazdo = nil
}

// czekajNaOdpowiedz czyta ramki aż do odpowiedzi na to jedno żądanie.
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

// nawiazane zwraca gniazdo czynne, zestawiając je przy pierwszym użyciu.
func (p *Polaczenie) nawiazane(kontekst context.Context) (*websocket.Conn, error) {
	if p.gniazdo != nil {
		return p.gniazdo, nil
	}
	gniazdo, _, err := websocket.Dial(kontekst, p.adres, nil)
	if err != nil {
		return nil, fmt.Errorf("rdzeń pod adresem %s nie odpowiada: %w", p.adres, err)
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
	return fmt.Errorf("połączenie z rdzeniem zerwane przy czynności %q: %w", czynnosc, err)
}
