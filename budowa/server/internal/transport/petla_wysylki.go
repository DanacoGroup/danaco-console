package transport

import (
	"github.com/coder/websocket"
)

// petlaWysylki zapisuje do gniazda ramki odłożone przez Wyslij jako jedyny pisarz tego gniazda, dzięki czemu kolejność ramek odpowiada kolejności wysyłek, a pętla kończy się wraz z kontekstem połączenia albo pierwszym błędem zapisu.
func (p *Polaczenie) petlaWysylki() {
	for {
		select {
		case <-p.kontekst.Done():
			return
		case dane := <-p.wyjscie:
			if err := p.gniazdo.Write(p.kontekst, websocket.MessageText, dane); err != nil {
				p.dziennik.Printf("transport: zapis do %s przerwany: %v", p.id, err)
				p.Zamknij("błąd zapisu")
				return
			}
		}
	}
}
