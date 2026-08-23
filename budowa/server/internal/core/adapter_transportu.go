package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/transport"
)

// wejscieTransportu podaje rdzeń warstwie transportu w kształcie, którego
// transport oczekuje, i odwrotnie: podaje rdzeniowi transport jako nadajnik
// i nasłuch.
//
// Zależność idzie w jedną stronę. Transport nie zna rdzenia: zna
// wyłącznie własne interfejsy `Rdzen`, `Ujscie` i `Rozglosnik`, więc da się go
// wymienić i uruchomić bez rdzenia. Rdzeń zna transport, bo to w rdzeniu stoi
// montaż — `montaz.go` składa serwer i podaje mu to wejście.
//
// Prawdą wartą pilnowania nie jest brak obu zależności (ten plik i `montaz.go`
// importują transport), tylko brak cyklu: dopisanie w transporcie importu
// z `core` wywraca kompilację i o to właśnie chodzi.
type wejscieTransportu struct {
	rdzen *Rdzen
}

// Obsluz wypełnia interfejs transport.Rdzen. Odpowiedź wraca wynikiem, a strumień
// i zdarzenia idą rozgłoszeniem do wszystkich urządzeń konta — synchronizacja
// wielourządzeniowa nie ma osobnego protokołu.
//
// Z ujścia bierzemy jedno: jego tożsamość. Wchodzi do kontekstu żądania cała
// (`transport.Tozsamosc`), bo bez niej rdzeń nie umiał odpowiedzieć na pytanie
// „kto stoi po drugiej stronie tego gniazda" — a od tego zależy powitanie,
// wyłączenie bieżącej sesji ze zmiany hasła (`wiez_polaczenia.go`) oraz sprawca
// wpisywany w zdarzenia zmiany (`sprawca.go`). Poszerzenie tego jednego wpisu
// jest tu istotą: drugi wpis obok byłby drugą prawdą o tym samym gnieździe.
// Samo ujście dalej nie idzie: rdzeń ma nie umieć odesłać czegokolwiek do
// jednego urządzenia z pominięciem rozgłoszenia, bo to byłaby druga droga
// wyjścia obok Nadajnika.
func (w wejscieTransportu) Obsluz(kontekst context.Context, zadanie protocol.Request, ujscie transport.Ujscie) protocol.Koperta {
	if ujscie != nil {
		kontekst = zPolaczeniem(kontekst, ujscie.Tozsamosc())
	}
	return w.rdzen.WykonajZadanie(kontekst, zadanie)
}

// Przylaczono wypełnia nieobowiązkowe rozszerzenie transport.ObserwatorPolaczen.
// Nowe połączenie nie jest jeszcze z niczym związane — więź zakłada dopiero
// `connection.hello` z tokenem albo udane wejście przez bramkę.
func (w wejscieTransportu) Przylaczono(transport.Ujscie) {}

// Odlaczono zdejmuje więź rozłączonego urządzenia. Bez tego mapa więzi rosłaby
// o wpis na każde nawiązanie przez całe życie procesu.
func (w wejscieTransportu) Odlaczono(ujscie transport.Ujscie) {
	if ujscie == nil || w.rdzen == nil {
		return
	}
	w.rdzen.wiez.Rozwiaz(ujscie.Id())
}

// nasluchTransportu podaje serwer transportu jako nasłuch rdzenia oraz jako
// nadajnik zdarzeń.
type nasluchTransportu struct {
	serwer *transport.Serwer
}

// Sluchaj otwiera nasłuch i pracuje do zamknięcia kontekstu, po czym zamyka
// serwer. Rozłączenie klienta nie kończy pracy rdzenia.
func (n nasluchTransportu) Sluchaj(kontekst context.Context) error {
	if err := n.serwer.Uruchom(kontekst); err != nil {
		return err
	}
	<-kontekst.Done()
	return n.serwer.Zamknij()
}

// Rozglos wypełnia port Nadajnik rdzenia. Puste konto oznacza rozgłoszenie do
// wszystkich połączeń rdzenia — w fazie budowy uwierzytelnianie jest wyłączone,
// więc urządzenia pracują na koncie domyślnym.
func (n nasluchTransportu) Rozglos(k protocol.Koperta) {
	n.serwer.Rozglos("", k)
}
