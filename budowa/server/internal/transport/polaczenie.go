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

// bladPolaczenieZamkniete oznacza wysyłkę do urządzenia, które się rozłączyło.
// Jest błędem pojedynczego wywołania, nie stanem sesji.
var bladPolaczenieZamkniete = errors.New("transport: połączenie zamknięte")

// bladKolejkaPelna oznacza urządzenie wolniejsze niż strumień. Rozgłoszenie
// pomija wtedy jedno urządzenie i żyje dalej; wysyłka adresowana kończy się
// zamknięciem połączenia, bo porzucona odpowiedź na komendę zostawia po drugiej
// stronie obietnicę, która nie rozstrzygnie się nigdy.
var bladKolejkaPelna = errors.New("transport: kolejka wyjściowa pełna")

const (
	// czasNaMiejsceWKolejce jest chwilą czekania na miejsce w kolejce wyjściowej
	// przy wysyłce adresowanej. Krótką: urządzenie, które przez tyle czasu nie
	// odebrało ani jednej ramki, nie odbierze i następnej.
	czasNaMiejsceWKolejce = 2 * time.Second
	// czasNaDosylke ogranicza zapis ramki wysyłanej z pominięciem kolejki, tuż
	// przed zamknięciem połączenia. Odmowa ma dojść do urządzenia przed
	// zerwaniem, inaczej po drugiej stronie zostaje samo zerwanie bez powodu.
	czasNaDosylke = time.Second
)

// Polaczenie jest jednym kanałem WebSocket urządzenia i realizuje interfejs Ujscie widziany przez rdzeń.
type Polaczenie struct {
	id       string
	gniazdo  *websocket.Conn
	wyjscie  chan []byte
	kontekst context.Context
	zakoncz  context.CancelFunc
	dziennik *log.Logger

	zamek     sync.RWMutex
	konto     string
	tozsamosc Tozsamosc
	zamkniete bool
	// przeszlaBramke pamięta, że gniazdo wykonało już komendę spoza wykazu
	// wejścia. Późniejsza odmowa bramki znaczy wtedy sesję, która przestała
	// nadawać — i gniazdo idzie do zamknięcia, a nie do kolejnej odmowy.
	przeszlaBramke bool
}

// Funkcja nowePolaczenie owija świeżo przyjęte gniazdo WebSocket w strukturę gotową do pracy z rdzeniem.
func nowePolaczenie(rodzic context.Context, id, konto string, tozsamosc Tozsamosc, gniazdo *websocket.Conn, pojemnosc int, dziennik *log.Logger) *Polaczenie {
	kontekst, zakoncz := context.WithCancel(rodzic)
	return &Polaczenie{
		id:        id,
		gniazdo:   gniazdo,
		wyjscie:   make(chan []byte, pojemnosc),
		kontekst:  kontekst,
		zakoncz:   zakoncz,
		dziennik:  dziennik,
		konto:     konto,
		tozsamosc: tozsamosc,
	}
}

// Metoda Id zwraca identyfikator tego połączenia, nadany mu przy jego nawiązaniu z danym urządzeniem klienckim.
func (p *Polaczenie) Id() string {
	return p.id
}

// Metoda Konto zwraca konto urządzenia przypisane obecnie do tego konkretnego połączenia gniazda WebSocket.
func (p *Polaczenie) Konto() string {
	p.zamek.RLock()
	defer p.zamek.RUnlock()
	return p.konto
}

// PrzypiszKonto zmienia konto połączenia. Puste konto jest pomijane — rdzeń nie
// odbiera urządzeniu przynależności przez przeoczenie.
func (p *Polaczenie) PrzypiszKonto(konto string) {
	if konto == "" {
		return
	}
	p.zamek.Lock()
	p.konto = konto
	p.zamek.Unlock()
}

// Metoda Tozsamosc zwraca to, co o drugiej stronie gniazda wiadomo w tej chwili, licząc od nawiązania połączenia.
func (p *Polaczenie) Tozsamosc() Tozsamosc {
	p.zamek.RLock()
	defer p.zamek.RUnlock()
	return p.tozsamosc
}

// PrzedstawKlienta zapisuje identyfikator klienta odczytany z powitania.
//
// Identyfikator pusty jest pomijany — powitanie bez `clientId` nie odbiera
// połączeniu tożsamości przedstawionej przy nawiązaniu (wzorem PrzypiszKonto).
func (p *Polaczenie) PrzedstawKlienta(id string) {
	if id == "" {
		return
	}
	p.zamek.Lock()
	p.tozsamosc.IdKlienta = id
	p.zamek.Unlock()
}

// Wyslij koduje kopertę i stawia ją w kolejce wyjściowej. Nie czeka na zapis do
// gniazda, więc wolne urządzenie nie zatrzymuje rdzenia; czeka wyłącznie na
// miejsce w kolejce, bo tędy idą odpowiedzi na komendy i fragmenty domykające
// strumień — te dwie ramki porzucone milczeniem zostawiają wołającego bez
// rozstrzygnięcia.
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

// wyslijZCzekaniem stawia ramkę w kolejce, dając urządzeniu chwilę na zwolnienie miejsca.
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

// Metoda wyslijBajty stawia już gotową ramkę bajtów w kolejce wyjściowej tego samego połączenia gniazda.
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

// Zamknij kończy pracę połączenia. Wywołanie powtórne nic nie zmienia.
// Zamknięcie dotyczy wyłącznie kanału urządzenia: sesje, okna i procesy rdzenia
// biegną dalej.
func (p *Polaczenie) Zamknij(powod string) {
	p.ZamknijKodem(websocket.StatusNormalClosure, powod)
}

// ZamknijKodem kończy pracę połączenia kodem zamknięcia kontraktu WebSocket.
// Kod niesie rozróżnienie, którego powód tekstowy nie niesie: urządzenie ma
// wiedzieć, czy zerwanie jest zwykłym końcem pracy, czy odmową.
func (p *Polaczenie) ZamknijKodem(kod websocket.StatusCode, powod string) {
	p.zamek.Lock()
	if p.zamkniete {
		p.zamek.Unlock()
		return
	}
	p.zamkniete = true
	p.zamek.Unlock()

	// Ramka zamknięcia idzie przed zakończeniem kontekstu: anulowanie kontekstu
	// odczytu zamyka gniazdo TCP, a urządzenie dostałoby EOF zamiast kodu.
	if err := p.gniazdo.Close(kod, powod); err != nil {
		_ = p.gniazdo.CloseNow()
	}
	p.zakoncz()
	p.dziennik.Printf("transport: rozłączenie %s (%s)", p.id, powod)
}

// wyslijWprost zapisuje kopertę do gniazda z pominięciem kolejki wyjściowej.
// Idzie tędy wyłącznie ramka poprzedzająca zamknięcie: ramka odłożona do kolejki
// nie zdążyłaby dojść, bo zamknięcie kończy pętlę wysyłki razem z kontekstem
// połączenia. Biblioteka gniazda pilnuje wyłączności zapisu, więc równoległa
// pętla wysyłki nie miesza ramek.
func (p *Polaczenie) wyslijWprost(k protocol.Koperta) error {
	dane, err := protocol.Zakoduj(k)
	if err != nil {
		return err
	}
	kontekst, koniec := context.WithTimeout(p.kontekst, czasNaDosylke)
	defer koniec()
	return p.gniazdo.Write(kontekst, websocket.MessageText, dane)
}

// oznaczPrzejscieBramki zapisuje, że gniazdo wykonało komendę spoza wykazu wejścia, czyli że jego sesja bramki nadawała.
func (p *Polaczenie) oznaczPrzejscieBramki() {
	p.zamek.Lock()
	p.przeszlaBramke = true
	p.zamek.Unlock()
}

// przeszlaPrzezBramke mówi, czy gniazdo pracowało już na sesji bramki.
func (p *Polaczenie) przeszlaPrzezBramke() bool {
	p.zamek.RLock()
	defer p.zamek.RUnlock()
	return p.przeszlaBramke
}
