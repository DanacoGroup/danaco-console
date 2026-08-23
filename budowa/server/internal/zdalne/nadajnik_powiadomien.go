// Odpowiedzialność pliku: port własny wołania urządzenia. Pakiet opisuje tu
// kształt nadajnika, którego potrzebuje, i nie bierze interfejsu z pakietu
// transportu ani z rdzenia — dokładnie tak, jak `core` opisuje u siebie port
// `Nadajnik` zamiast importować `transport`.
//
// Port nie niesie koperty. Nośnik `transport.Serwer.Rozglos` przyjmuje
// `protocol.Koperta`, ale koperta niesie `type: MessageType` z kontraktu,
// a kontrakt nie ma dziś ani jednego zdarzenia powiadomienia. Zbudowanie koperty
// z typem spoza kontraktu wysłałoby po WebSocket komunikat, którego klient nie
// zna — czyli budzik, który nigdzie nie dzwoni.
//
// Dlatego port bierze sam fakt do przekazania, a złożenie koperty należy do tego,
// kto zna kontrakt: adaptera w rdzeniu, wpiętego po wniesieniu zdarzenia
// `notification.raised`. Ten pakiet prowadzi kolejkę i rozstrzyga, co i kiedy ma
// polecieć; nie rozstrzyga, jakim napisem się to nazywa na łączu.
package zdalne

import (
	"context"

	"danacoconsole/server/internal/dane"
)

// Cel jest jednym adresem wołania — wierszem rejestracji `urzadzenie_powiadomien`.
type Cel struct {
	// Rejestracja jest kluczem wiersza `urzadzenie_powiadomien`; to nim ślad
	// doręczenia wskazuje drogę, którą powiadomienie naprawdę poszło.
	Rejestracja int64
	// Kanal jest dziś zawsze `polaczenie` — jedyna droga, którą rdzeń umie.
	Kanal string
	// KluczKanalu dla kanału `polaczenie` jest identyfikatorem klienta,
	// tym samym, który transport niesie jako `Tozsamosc.IdKlienta`.
	KluczKanalu string
	// Etykieta jest napisem Operatora o tym urządzeniu; może być pusta.
	Etykieta string
}

// Nadajnik doręcza powiadomienie pod jeden adres.
//
// Zwraca prawdę wyłącznie wtedy, gdy urządzenie kopertę przyjęło. „Wysłałem
// i nie wiem" ma zwracać fałsz: kolejka zapisuje `dostarczone` na podstawie tej
// odpowiedzi, więc odpowiedź niepewna zamieniłaby się w bazie w twierdzenie
// pewne — i Operator zobaczyłby „doszło" tam, gdzie nie doszło.
// Ponowienie kosztuje jeden takt; nieprawda w bazie kosztuje zaufanie do kanału.
type Nadajnik interface {
	Zawolaj(ctx context.Context, cel Cel, p dane.Powiadomienie) bool
}
