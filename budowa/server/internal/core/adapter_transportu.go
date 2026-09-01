// Plik jest stykiem rdzenia z warstwą transportu: wnosi do kontekstu żądania tożsamość i ujście gniazda, przypisuje gniazdu konto rozpoznane przez bramkę i oddaje serwer transportu jako nasłuch oraz nadajnik zdarzeń.
package core

import (
	"context"
	"strconv"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/transport"
)

// wejscieTransportu podaje rdzeń warstwie transportu w kształcie, którego transport oczekuje,
// a rdzeniowi podaje transport jako nadajnik i nasłuch. Zależność idzie w jedną stronę:
// transport nie zna rdzenia.
type wejscieTransportu struct {
	rdzen *Rdzen
}

// Obsluz wypełnia interfejs transport.Rdzen. Odpowiedź wraca wynikiem, a strumień i zdarzenia
// idą rozgłoszeniem do wszystkich urządzeń konta, bo synchronizacja wielourządzeniowa nie ma
// osobnego protokołu.
func (w wejscieTransportu) Obsluz(kontekst context.Context, zadanie protocol.Request, ujscie transport.Ujscie) protocol.Koperta {
	if ujscie != nil {
		kontekst = zUjsciem(zPolaczeniem(kontekst, ujscie.Tozsamosc()), ujscie)
	}
	return w.rdzen.WykonajZadanie(kontekst, zadanie)
}

// Przylaczono wypełnia nieobowiązkowe rozszerzenie transport.ObserwatorPolaczen.
// Konto przedstawione przy nawiązaniu jest słowem klienta, a nie ustaleniem
// rdzenia, więc świeże gniazdo wraca na konto domyślne: przynależność do konta
// nadaje wyłącznie bramka — `connection.hello` z ważnym tokenem albo udane
// wejście.
func (w wejscieTransportu) Przylaczono(ujscie transport.Ujscie) {
	if ujscie == nil {
		return
	}
	ujscie.PrzypiszKonto(transport.KontoDomyslne)
}

// Odlaczono zdejmuje więź rozłączonego urządzenia. Bez tego mapa więzi rosłaby
// o wpis na każde nawiązanie przez całe życie procesu.
func (w wejscieTransportu) Odlaczono(ujscie transport.Ujscie) {
	if ujscie == nil || w.rdzen == nil {
		return
	}
	w.rdzen.wiez.Rozwiaz(ujscie.Id())
}

// kluczUjscia niesie ujście gniazda, z którego przyszło żądanie. Wpis stoi obok
// kluczPolaczenia, a nie zamiast niego: tożsamość jest kopią faktów o wołającym,
// ujście jest żywym kanałem, któremu rdzeń przypisuje konto.
const kluczUjscia kluczKontekstu = "danaco:ujscie"

// zUjsciem dokłada do kontekstu ujście gniazda. Wpina to `wejscieTransportu.Obsluz` —
// jedyne miejsce, przez które przechodzi każde żądanie z gniazda.
func zUjsciem(ctx context.Context, ujscie transport.Ujscie) context.Context {
	if ctx == nil || ujscie == nil {
		return ctx
	}
	return context.WithValue(ctx, kluczUjscia, ujscie)
}

// ujscieZKontekstu oddaje ujście wołającego. Zero znaczy żądanie spoza gniazda —
// tak wygląda bieg wewnętrzny rdzenia i wywołanie ze sprawdzianu; to nie jest błąd.
func ujscieZKontekstu(ctx context.Context) transport.Ujscie {
	if ctx == nil {
		return nil
	}
	ujscie, jest := ctx.Value(kluczUjscia).(transport.Ujscie)
	if !jest {
		return nil
	}
	return ujscie
}

// przypiszKontoGniazda wiąże gniazdo z kontem, któremu wydano sesję bramki.
// Transport adresuje rozgłoszenia kontem połączenia, a wie o nim tylko tyle,
// ile mu rdzeń powie — bez tego wywołania konto gniazda zostaje takie, jakim
// przedstawił je klient. Skrót pusty odsyła gniazdo na konto domyślne, bo
// połączenie bez rozpoznanej sesji nie ma należeć do żadnego konta.
func przypiszKontoGniazda(ctx context.Context, konta RozpoznanieKontaSesji, skrot string) {
	ujscie := ujscieZKontekstu(ctx)
	if ujscie == nil {
		return
	}
	ujscie.PrzypiszKonto(nazwaKontaGniazda(ctx, konta, skrot))
}

// nazwaKontaGniazda przekłada sesję bramki na nazwę konta, którą posługuje się
// transport. Nazwą jest identyfikator konta z bazy — innej nazwy konta transport
// nie zna. Sesja nierozpoznana daje konto domyślne, nie napis pusty: pustego
// konta PrzypiszKonto nie przyjmuje, a gniazdo musi mieć jak zejść z konta
// przypisanego wcześniej.
func nazwaKontaGniazda(ctx context.Context, konta RozpoznanieKontaSesji, skrot string) string {
	if konta == nil || skrot == "" {
		return transport.KontoDomyslne
	}
	kontoId, err := konta.KontoSesjiBramki(ctx, skrot)
	if err != nil || kontoId == 0 {
		return transport.KontoDomyslne
	}
	return strconv.FormatInt(kontoId, 10)
}

// nasluchTransportu podaje serwer transportu jako nasłuch rdzenia oraz jako nadajnik zdarzeń
// rozgłaszanych do urządzeń konta.
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

// Rozglos wypełnia port Nadajnik rdzenia. Konto adresata zostaje puste — czyli
// rozgłoszenie idzie do wszystkich połączeń rdzenia — bo port Nadajnik niesie
// samą kopertę, a koperta kontraktu nie ma pola adresata. Zawężenie do konta
// wymaga konta w porcie Nadajnik i w emiterze zdarzeń.
func (n nasluchTransportu) Rozglos(k protocol.Koperta) {
	n.serwer.Rozglos("", k)
}
