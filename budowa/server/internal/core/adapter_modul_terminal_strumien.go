// Odpowiedzialność pliku: wyjście procesu terminala jako strumień fragmentów kontraktu — nośnik okna Output Console. Wyjście jedzie wspólnym strumieniem `stream.chunk`, opisanym kontraktem jako jedna droga dla wszystkich kanałów.
package core

import (
	"io"
	"strconv"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// rozmiarBuforaWyjscia jest porcją odczytu z potoku wyjścia procesu terminala uruchomionego przez kartę.
const rozmiarBuforaWyjscia = 32 * 1024

// nadawcaWyjscia rozsyła fragmenty wyjścia procesów terminala zdarzeniami transportu okna komunikacji.
type nadawcaWyjscia struct {
	nadajnik Nadajnik
}

// nowyNadawcaWyjscia opakowuje nadajnik zdarzeń transportu nadawcą wyjścia procesów terminala montowanym przy starcie.
func nowyNadawcaWyjscia(nadajnik Nadajnik) *nadawcaWyjscia {
	return &nadawcaWyjscia{nadajnik: nadajnik}
}

// Pompuj czyta wyjście procesu do wyczerpania i rozsyła je fragmentami. Wywołuje się to w osobnej gorutynie — jedną dla wyjścia zwykłego, jedną dla diagnostycznego.
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
			// Koniec potoku jest normalnym końcem odczytu; proces i tak domknie go czekający na zakończenie.
			return
		}
	}
}

// Fragment wysyła jedną porcję wyjścia do urządzeń konta, które proces zamówiło. Fragment ostatni zamyka strumień procesu znacznikiem `done` — bez niego Output Console czekałby na ciąg dalszy, którego nigdy nie będzie.
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
	n.nadajnik.Rozglos(kontoAdresata(proces.kontekst), koperta)
}

// Domknij wysyła fragment zamykający strumień procesu wraz z podsumowaniem zakończenia. Output Console pokazuje tę treść jako ostatni wiersz przebiegu.
func (n *nadawcaWyjscia) Domknij(proces *procesTerminala, stan shared.TerminalProcessStatus,
	kodWyjscia *int, powod string) {

	var rodzaj shared.ChunkKind = shared.ChunkKindText
	if stan != shared.TerminalProcessStatusFinished {
		rodzaj = shared.ChunkKindError
	}
	n.Fragment(proces, rodzaj, podsumowanieZakonczenia(stan, kodWyjscia, powod), true)
}

// podsumowanieZakonczenia składa ostatni wiersz przebiegu procesu, widoczny jako podsumowanie w oknie.
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
