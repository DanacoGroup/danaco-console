// Pakiet uruchomienie rozgałęzia pracę procesu według roli: rola rozstrzyga,
// które tory ruszają. Tor interfejsu otwiera nasłuch transportu, tor wykonawczy
// przyjmuje żądania kontraktu strumieniami procesu bez zajmowania portu.
package uruchomienie

import (
	"context"
	"io"
	"log"
	"strings"
)

// Rdzen to jedyne, czego tory potrzebują od złożonego rdzenia. Interfejs stoi
// po stronie odbiorcy, więc pakiet nie importuje core i nie tworzy cyklu.
type Rdzen interface {
	// Uruchom oddaje sterowanie warstwie nasłuchu i wraca po zamknięciu kontekstu.
	Uruchom(kontekst context.Context) error
	// WykonajSurowe wykonuje jedno żądanie kontraktu w bajtach i zwraca bajty odpowiedzi.
	WykonajSurowe(kontekst context.Context, dane []byte) []byte
}

// Otoczenie niesie strumienie toru wykonawczego oraz dziennik procesu.
// Każde pole jest opcjonalne: brak strumienia albo dziennika nie wstrzymuje
// uruchomienia.
type Otoczenie struct {
	// Wejscie dostarcza żądania kontraktu, po jednym w wierszu. Nil oznacza wejście puste.
	Wejscie io.Reader
	// Wyjscie przyjmuje odpowiedzi, po jednej w wierszu. Nil kieruje je do kosza.
	Wyjscie io.Writer
	// Dziennik przyjmuje zapisy diagnostyczne. Nil wycisza je w całości.
	Dziennik *log.Logger
}

// wejscie zwraca strumień żądań podany w polu Wejscie albo pusty strumień
// zastępczy, gdy pole go nie wskazuje.
func (o Otoczenie) wejscie() io.Reader {
	if o.Wejscie == nil {
		return strings.NewReader("")
	}
	return o.Wejscie
}

// wyjscie zwraca strumień odpowiedzi podany w polu Wyjscie albo kosz zastępczy,
// gdy pole go nie wskazuje.
func (o Otoczenie) wyjscie() io.Writer {
	if o.Wyjscie == nil {
		return io.Discard
	}
	return o.Wyjscie
}

// zapisz odnotowuje zdarzenie toru w dzienniku procesu. Brak dziennika nie
// zmienia zachowania funkcji.
func (o Otoczenie) zapisz(wzorzec string, argumenty ...any) {
	if o.Dziennik == nil {
		return
	}
	o.Dziennik.Printf(wzorzec, argumenty...)
}
