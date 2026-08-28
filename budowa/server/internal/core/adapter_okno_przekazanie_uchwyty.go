// Plik wpina komendy obszaru window.* portu PrzekazanieOkna, którego jedyna
// implementacja leży w dwóch plikach obsługujących osobno przekazanie okna
// i akcję panelu okna.
package core

import (
	"context"

	"danacoconsole/shared"
)

// PrzekazanieOkna jest portem obszaru window.*, obsługującym przekazanie okna
// i wykonanie akcji panelu okna.
type PrzekazanieOkna interface {
	Przekaz(ctx context.Context, z shared.WindowHandoffRequest) (shared.WindowHandoffResponse, error)
	WykonajAkcje(ctx context.Context, z shared.WindowActionRequest) (shared.WindowActionResponse, error)
}

// zarejestrujPrzekazanieOkna wpina dwie komendy obszaru window.* w rejestr
// obsługiwaczy komend rdzenia.
func zarejestrujPrzekazanieOkna(r *Rejestr, m PrzekazanieOkna, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandWindowHandoff, obsluz(m.Przekaz))
	r.Zarejestruj(shared.CommandWindowAction, obsluz(m.WykonajAkcje))

	// Kontrakt nie ma zdarzenia dla akcji ani przekazania okna, emiter zostaje niewykorzystany.
}
