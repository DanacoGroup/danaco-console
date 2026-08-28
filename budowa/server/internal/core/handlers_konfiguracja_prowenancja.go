// Plik wpina komendy config.explain.get oraz config.window.open rodziny
// config.*, rozszerzając port ustawień o prowenancję wywołania i otwarcie
// okna konfiguracji.
package core

import (
	"context"

	"danacoconsole/shared"
)

// ProwenancjaKonfiguracji jest portem rodziny `config.*` rozszerzonym o dwie
// komendy, które sięgają po DROGĘ TURY zamiast po rejestr ustawień.
type ProwenancjaKonfiguracji interface {
	Ustawienia

	// ProwenancjaWywolania obsługuje config.explain.get, zwracając wiersz
	// wywołania i sumę kontrolną.
	ProwenancjaWywolania(ctx context.Context,
		z shared.ConfigExplainGetRequest) (shared.ConfigExplainGetResponse, error)

	// OtworzOknoKonfiguracji obsługuje config.window.open, rozstrzygając
	// zakres startowy okna.
	OtworzOknoKonfiguracji(ctx context.Context,
		z shared.ConfigWindowOpenRequest) (shared.ConfigWindowOpenResponse, error)
}

// zarejestrujProwenancjeKonfiguracji wpina dwie komendy rodziny `config.*`.
//
// Port pusty nie rejestruje niczego: obie komendy odpowiedzą wtedy
// `config.unknown`, a pozostałe siedem rodziny pracuje bez zmian.
func zarejestrujProwenancjeKonfiguracji(r *Rejestr, p ProwenancjaKonfiguracji) {
	if r == nil || p == nil {
		return
	}

	r.Zarejestruj(shared.CommandConfigExplainGet, obsluz(p.ProwenancjaWywolania))
	r.Zarejestruj(shared.CommandConfigWindowOpen, obsluz(p.OtworzOknoKonfiguracji))
}
