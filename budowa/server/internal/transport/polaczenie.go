package transport

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/coder/websocket"

	"danacoconsole/server/internal/protocol"
)

// bladPolaczenieZamkniete oznacza wysyłkę do urządzenia, które się rozłączyło.
// Jest błędem pojedynczego wywołania, nie stanem sesji.
var bladPolaczenieZamkniete = errors.New("transport: połączenie zamknięte")

// bladKolejkaPelna oznacza urządzenie wolniejsze niż strumień. Wysyłka zostaje
// odrzucona, połączenie żyje dalej — nadawca decyduje, czy powtórzyć.
var bladKolejkaPelna = errors.New("transport: kolejka wyjściowa pełna")

// Polaczenie jest jednym kanałem WebSocket urządzenia. Realizuje interfejs
// Ujscie, więc rdzeń widzi wyłącznie identyfikator, konto i wysyłkę koperty.
//
// Zapis do gniazda prowadzi jedna pętla (petla_wysylki.go), odczyt druga
// (petla_odbioru.go). Dzięki temu wysyłka z wielu miejsc rdzenia równocześnie
// nie wymaga blokady na gnieździe i nigdy nie miesza ramek.
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
}

// nowePolaczenie owija świeżo przyjęte gniazdo.
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

// Id zwraca identyfikator połączenia.
func (p *Polaczenie) Id() string {
	return p.id
}

// Konto zwraca konto urządzenia.
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

// Tozsamosc zwraca to, co o drugiej stronie gniazda wiadomo w tej chwili.
//
// „W tej chwili" jest tu istotne: identyfikator klienta dochodzi dopiero
// z powitaniem, więc żądanie wcześniejsze widzi tożsamość uboższą. To jest
// prawda o stanie, a nie brak do naprawienia — rdzeń zamilknie wtedy o sprawcy
// zamiast go zmyślić.
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
// gniazda, więc wolne urządzenie nie zatrzymuje rdzenia.
func (p *Polaczenie) Wyslij(k protocol.Koperta) error {
	dane, err := protocol.Zakoduj(k)
	if err != nil {
		return err
	}
	return p.wyslijBajty(dane)
}

// wyslijBajty stawia gotową ramkę w kolejce wyjściowej.
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
	p.zamek.Lock()
	if p.zamkniete {
		p.zamek.Unlock()
		return
	}
	p.zamkniete = true
	p.zamek.Unlock()

	p.zakoncz()
	if err := p.gniazdo.Close(websocket.StatusNormalClosure, powod); err != nil {
		_ = p.gniazdo.CloseNow()
	}
	p.dziennik.Printf("transport: rozłączenie %s (%s)", p.id, powod)
}
