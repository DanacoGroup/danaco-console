package core

import (
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// odpowiedzNieznanej buduje odpowiedź na komendę, której rdzeń nie obsługuje.
//
// To jest ścieżka fail-open. Nieznana nazwa nie jest błędem
// protokołu i nie ma prawa niczego zamknąć:
//
//   - połączenie zostaje otwarte,
//   - sesja pracuje dalej,
//   - kolejne żądania są przyjmowane,
//   - klient dostaje zdarzenie `<obszar>.unknown` z nazwą, której nie
//     rozpoznano, więc widzi przyczynę zamiast ciszy.
//
// Nazwę zdarzenia wskazuje kontrakt, po obszarze żądanego typu — rdzeń nie
// wybiera jej sam i nie ma tu ani jednego literału nazwy.
func odpowiedzNieznanej(z protocol.Request) protocol.Koperta {
	return protocol.OdpowiedzNieznanej(z)
}

// kopertaNieczytelna odpowiada na komunikat, którego nie da się odkodować.
//
// Komunikat niepoprawny strukturalnie nie ma typu ani identyfikatora, więc nie
// da się zbudować dla niego żądania ani skorelować odpowiedzi. Odpowiedź idzie
// mimo to — zerwanie połączenia byłoby karą za jeden zepsuty bajt.
//
// Typ odpowiedzi wskazuje kontrakt: nierozpoznany komunikat bez obszaru należy
// do obszaru połączenia, więc wraca jego zdarzeniem `*.unknown`. Kod błędu mówi
// klientowi, że rzecz nie w nieznanej nazwie, lecz w niepoprawnej treści.
func kopertaNieczytelna(err error) protocol.Koperta {
	naglowek := protocol.Koperta{
		Type:      shared.ZdarzenieNieznanej(""),
		Timestamp: protocol.Teraz(),
	}
	return protocol.KopertaOdpowiedzi(naglowek, bladNiepoprawnegoLadunku(err))
}
