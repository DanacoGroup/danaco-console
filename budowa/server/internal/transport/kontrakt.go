// Pakiet transport jest warstwą wejścia rdzenia, nasłuchującą na WebSocket i rozgłaszającą zdarzenia do urządzeń konta.
package transport

import (
	"context"

	"danacoconsole/server/internal/protocol"
)

// Ujscie jest kanałem wyjściowym pojedynczego połączenia. Realizuje je
// transport, a używa rdzeń — stąd rdzeń odsyła fragmenty strumienia i zdarzenia
// do konkretnego urządzenia bez wiedzy o gnieździe, kolejce ani bibliotece.
type Ujscie interface {
	// Id zwraca identyfikator połączenia nadany przy nawiązaniu.
	Id() string
	// Konto zwraca konto urządzenia — adresata rozgłoszeń.
	Konto() string
	// PrzypiszKonto wiąże połączenie z kontem po rozpoznaniu urządzenia zgłaszającego się w gnieździe.
	PrzypiszKonto(konto string)
	// Tozsamosc zwraca fakty o drugiej stronie gniazda; rdzeń rozstrzyga z nich sprawcę zdarzenia.
	Tozsamosc() Tozsamosc
	// Wyslij kieruje kopertę do tego jednego połączenia; błąd dotyczy wyłącznie tego wywołania.
	Wyslij(k protocol.Koperta) error
}

// Rdzen jest jedynym wejściem transportu do warstwy wykonawczej; transport dostarcza żądanie i ujście, a rdzeń zwraca kopertę.
type Rdzen interface {
	Obsluz(kontekst context.Context, zadanie protocol.Request, ujscie Ujscie) protocol.Koperta
}

// ObserwatorPolaczen jest rozszerzeniem nieobowiązkowym. Gdy rdzeń go realizuje,
// transport zawiadamia go o przyłączeniu i odłączeniu urządzenia. Brak
// realizacji nie jest błędem ani brakiem funkcji.
type ObserwatorPolaczen interface {
	Przylaczono(ujscie Ujscie)
	Odlaczono(ujscie Ujscie)
}

// Rozgłaszanie zdarzeń do wszystkich połączeń konta nie ma tutaj własnego interfejsu transportu.
