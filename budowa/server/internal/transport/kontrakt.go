// Pakiet transport jest warstwą wejścia rdzenia Danaco Console: nasłuchuje na
// WebSocket, utrzymuje wiele połączeń równocześnie, rozgłasza zdarzenia do
// wszystkich urządzeń konta i serwuje pliki klienta.
//
// Transport **nie zna rdzenia**. Zna wyłącznie interfejs Rdzen, który rdzeń
// realizuje, oraz kontrakt komunikatów z pakietu shared.
// W tym pakiecie nie ma ani jednego importu pakietu wykonawczego — dzięki temu
// zamiana rdzenia nie dotyka transportu, a transport da się uruchomić w teście
// z atrapą rdzenia po stronie testu.
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
	// PrzypiszKonto wiąże połączenie z kontem po rozpoznaniu urządzenia.
	// Do czasu przypisania obowiązuje konto z parametru nawiązania
	// (uwierzytelnianie jest jedyną kontrolą dostępu i w fazie budowy nie działa).
	PrzypiszKonto(konto string)
	// Tozsamosc zwraca fakty o drugiej stronie gniazda (`tozsamosc.go`): kto to
	// jest i czym jest. Rdzeń rozstrzyga z nich SPRAWCĘ zdarzenia; transport
	// żadnego rozstrzygnięcia na nich nie opiera.
	Tozsamosc() Tozsamosc
	// Wyslij kieruje kopertę do tego jednego połączenia. Błąd dotyczy wyłącznie
	// tego wywołania: nie zrywa sesji, nie blokuje kolejnych prób.
	Wyslij(k protocol.Koperta) error
}

// Rdzen jest jedynym wejściem transportu do warstwy wykonawczej. Transport
// dostarcza rozpoznane żądanie oraz ujście połączenia; rdzeń zwraca kopertę do
// odesłania.
//
// Koperta o pustym polu Type oznacza „odpowiedź pójdzie osobno" — transport nic
// wtedy nie odsyła, a rdzeń sam korzysta z Ujscie albo Rozglosnik. Kontekst jest
// kontekstem **serwera**, nie połączenia: rozłączenie klienta nie przerywa
// rozpoczętej pracy rdzenia.
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

// Rozgłaszanie zdarzeń do wszystkich połączeń konta nie ma tutaj własnego
// interfejsu. Rdzeń opisuje tę drogę PORTEM WŁASNYM (core: Nadajnik) i sięga
// wprost po metodę Rozglos serwera — interfejs po stronie transportu byłby
// drugą deklaracją tej samej rzeczy.
