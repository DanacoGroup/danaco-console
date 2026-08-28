package core

import (
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// odpowiedzNieznanej buduje odpowiedź fail-open na komendę nieznaną rdzeniowi: połączenie i sesja pozostają czynne, a klient otrzymuje zdarzenie obszaru żądanego typu z nazwą nierozpoznanej komendy.
func odpowiedzNieznanej(z protocol.Request) protocol.Koperta {
	return protocol.OdpowiedzNieznanej(z)
}

// kopertaNieczytelna zwraca odpowiedź na komunikat nieodkodowalny strukturalnie, wskazując zdarzenie obszaru połączenia zamiast zrywać je z powodu pojedynczego błędnego bajtu.
func kopertaNieczytelna(err error) protocol.Koperta {
	naglowek := protocol.Koperta{
		Type:      shared.ZdarzenieNieznanej(""),
		Timestamp: protocol.Teraz(),
	}
	return protocol.KopertaOdpowiedzi(naglowek, bladNiepoprawnegoLadunku(err))
}
