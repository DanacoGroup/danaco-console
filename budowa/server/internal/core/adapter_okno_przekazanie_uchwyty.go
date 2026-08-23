// Wpięcie komend obszaru window.* — portu `PrzekazanieOkna`, którego jedyna
// implementacja `adapterPrzekazaniaOkna` leży w dwóch plikach:
// `adapter_okno_przekazanie.go` (`window.handoff`) i `adapter_okno_akcja.go`
// (`window.action`).
//
// Komenda `action.list` nie jest tu dublowana — rejestruje ją `zarejestrujAkcje`
// w `kompozycja.go`. Podwójna rejestracja tej samej nazwy komendy nadpisałaby
// jeden obsługiwacz drugim bez ostrzeżenia.
package core

import (
	"context"

	"danacoconsole/shared"
)

// PrzekazanieOkna jest portem obszaru window.*.
type PrzekazanieOkna interface {
	Przekaz(ctx context.Context, z shared.WindowHandoffRequest) (shared.WindowHandoffResponse, error)
	WykonajAkcje(ctx context.Context, z shared.WindowActionRequest) (shared.WindowActionResponse, error)
}

// zarejestrujPrzekazanieOkna wpina dwie komendy obszaru window.*.
func zarejestrujPrzekazanieOkna(r *Rejestr, m PrzekazanieOkna, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandWindowHandoff, obsluz(m.Przekaz))
	r.Zarejestruj(shared.CommandWindowAction, obsluz(m.WykonajAkcje))

	// Kontrakt nie ma zdarzenia rozgłaszającego wykonanie akcji okna ani nowe
	// przekazanie — `e` zostaje w sygnaturze dla zgodności z pozostałymi
	// funkcjami `zarejestruj<Modul>`.
}
