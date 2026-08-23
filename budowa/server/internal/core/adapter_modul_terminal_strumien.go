// Odpowiedzialność pliku: wyjście procesu terminala jako strumień fragmentów
// kontraktu — nośnik okna Output Console.
//
// Wyjście jedzie wspólnym strumieniem `stream.chunk`, który kontrakt opisuje
// jako jedną drogę dla wszystkich kanałów: `windowId` wskazuje okno terminala,
// `messageId` — identyfikator procesu z rejestru, a numer fragmentu i znacznik
// końca żyją w kopercie. Dzięki temu Output Console rozdziela wyjście po
// procesach i po kartach, nie zakładając drugiego protokołu.
//
// Odczyt jest blokowy: surowe bajty trafiają do bufora 32 KiB i wysyłane jest
// tyle, ile przyszło. Odczyt po liniach zawiesiłby się na wyjściu bez znaku
// końca linii (pasek postępu), a fragment na linię zamieniłby jedno `find /`
// w setki tysięcy kopert. Odczyt blokowy sam skleja napływ: im szybciej proces
// pisze, tym większe porcje wracają z jednego odczytu.
package core

import (
	"io"
	"strconv"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// rozmiarBuforaWyjscia jest porcją odczytu z potoku procesu.
const rozmiarBuforaWyjscia = 32 * 1024

// nadawcaWyjscia rozsyła fragmenty wyjścia procesów terminala.
type nadawcaWyjscia struct {
	nadajnik Nadajnik
}

// nowyNadawcaWyjscia opakowuje nadajnik zdarzeń transportu.
func nowyNadawcaWyjscia(nadajnik Nadajnik) *nadawcaWyjscia {
	return &nadawcaWyjscia{nadajnik: nadajnik}
}

// Pompuj czyta wyjście procesu do wyczerpania i rozsyła je fragmentami.
// Wywołuje się to w osobnej gorutynie — jedną dla wyjścia zwykłego, jedną dla
// diagnostycznego.
func (n *nadawcaWyjscia) Pompuj(proces *procesTerminala, zrodlo io.Reader, rodzaj shared.ChunkKind) {
	if zrodlo == nil {
		return
	}
	bufor := make([]byte, rozmiarBuforaWyjscia)
	for {
		odczytane, err := zrodlo.Read(bufor)
		if odczytane > 0 {
			n.Fragment(proces, rodzaj, string(bufor[:odczytane]), false)
		}
		if err != nil {
			// Koniec potoku jest normalnym końcem odczytu, a błąd odczytu
			// dotyczy tego jednego potoku — proces i tak domknie go czekający
			// na jego zakończenie.
			return
		}
	}
}

// Fragment wysyła jedną porcję wyjścia. Fragment ostatni zamyka strumień
// procesu znacznikiem `done` — bez niego Output Console czekałby na ciąg dalszy,
// którego nigdy nie będzie.
func (n *nadawcaWyjscia) Fragment(proces *procesTerminala, rodzaj shared.ChunkKind,
	tresc string, ostatni bool) {

	if n == nil || n.nadajnik == nil || proces == nil {
		return
	}
	fragment := protocol.ChunkTekstu(proces.oknoKod, proces.kod, tresc)
	fragment.Kind = rodzaj
	numer := int(proces.numerFragmentu.Add(1))
	koperta, err := protocol.KopertaFragmentu(proces.kod, proces.idSesji, numer, ostatni, fragment)
	if err != nil {
		return
	}
	n.nadajnik.Rozglos(koperta)
}

// Domknij wysyła fragment zamykający strumień procesu wraz z podsumowaniem
// zakończenia. Output Console pokazuje tę treść jako ostatni wiersz przebiegu.
func (n *nadawcaWyjscia) Domknij(proces *procesTerminala, stan shared.TerminalProcessStatus,
	kodWyjscia *int, powod string) {

	var rodzaj shared.ChunkKind = shared.ChunkKindText
	if stan != shared.TerminalProcessStatusFinished {
		rodzaj = shared.ChunkKindError
	}
	n.Fragment(proces, rodzaj, podsumowanieZakonczenia(stan, kodWyjscia, powod), true)
}

// podsumowanieZakonczenia składa ostatni wiersz przebiegu.
func podsumowanieZakonczenia(stan shared.TerminalProcessStatus, kodWyjscia *int, powod string) string {
	tresc := "\n[proces " + string(stan)
	if kodWyjscia != nil {
		tresc += ", kod wyjścia " + strconv.Itoa(*kodWyjscia)
	}
	if powod != "" {
		tresc += ", " + powod
	}
	return tresc + "]\n"
}
