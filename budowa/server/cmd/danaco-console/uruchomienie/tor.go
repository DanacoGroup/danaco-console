// Pakiet uruchomienie rozgałęzia pracę procesu według roli.
//
// Rola nie jest etykietą: rozstrzyga, które tory procesu ruszają. Tor interfejsu
// otwiera nasłuch transportu, tor wykonawczy przyjmuje żądania kontraktu
// strumieniami procesu i nie zajmuje żadnego portu. Punkt wejścia pozostaje
// kompozycją — cała decyzja mieszka tutaj.
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
	// Uruchom oddaje sterowanie warstwie nasłuchu i wraca po zamknięciu
	// kontekstu — jest torem interfejsu.
	Uruchom(kontekst context.Context) error
	// WykonajSurowe wykonuje jedno żądanie kontraktu podane w bajtach i zwraca
	// bajty odpowiedzi — jest wejściem toru wykonawczego.
	WykonajSurowe(kontekst context.Context, dane []byte) []byte
}

// Otoczenie niesie strumienie toru wykonawczego oraz dziennik procesu.
// Każde pole jest opcjonalne: brak strumienia albo dziennika nie wstrzymuje
// uruchomienia.
type Otoczenie struct {
	// Wejscie dostarcza żądania kontraktu, po jednym w wierszu. Nil oznacza
	// wejście puste — tor wykonawczy kończy pracę od razu, nie zrywając toru
	// interfejsu.
	Wejscie io.Reader
	// Wyjscie przyjmuje odpowiedzi, po jednej w wierszu. Nil kieruje je do kosza.
	Wyjscie io.Writer
	// Dziennik przyjmuje zapisy diagnostyczne. Nil wycisza je w całości.
	Dziennik *log.Logger
}

// wejscie zwraca strumień żądań albo strumień pusty, gdy Operator go nie wskazał.
func (o Otoczenie) wejscie() io.Reader {
	if o.Wejscie == nil {
		return strings.NewReader("")
	}
	return o.Wejscie
}

// wyjscie zwraca strumień odpowiedzi albo kosz, gdy Operator go nie wskazał.
func (o Otoczenie) wyjscie() io.Writer {
	if o.Wyjscie == nil {
		return io.Discard
	}
	return o.Wyjscie
}

// zapisz odnotowuje zdarzenie toru. Brak dziennika nie zmienia zachowania.
func (o Otoczenie) zapisz(wzorzec string, argumenty ...any) {
	if o.Dziennik == nil {
		return
	}
	o.Dziennik.Printf(wzorzec, argumenty...)
}
