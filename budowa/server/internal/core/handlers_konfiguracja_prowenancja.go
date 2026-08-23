// Odpowiedzialność pliku: wpięcie dwóch komend rodziny `config.*`:
// `config.explain.get` i `config.window.open`.
//
// Kontrakt niesie w tej rodzinie dziewięć komend. Pozostałe siedem wpina
// `handlers_config.go` (`config.get`, `config.set`, `config.reset`) oraz
// `handlers_sesja_konfiguracja.go` (cztery komendy jednolitego modelu
// konfiguracji sesji).
//
// PORT JEST ROZSZERZENIEM, NIE DRUGIM PORTEM. Konfiguracja
// ma w rdzeniu jedną bramę; gdyby prowenancja weszła osobnym polem `Porty`,
// rodzina `config.*` miałaby dwa źródła prawdy i dwa miejsca montażu. Dlatego
// ProwenancjaKonfiguracji osadza port Ustawienia — tak samo, jak port pamięci
// osadza port przestrzeni roboczej.
//
// ZDARZENIA TU NIE MA I NIE UDAJEMY GO. `config.changed` opisuje ZMIANĘ wartości
// ustawienia. Żadna z tych dwóch komend nic nie zapisuje: pierwsza liczy
// prowenancję na świeżo, druga rozstrzyga zakres. Rozgłaszanie zmiany, której
// nie było, wprowadziłoby w błąd każde otwarte okno konfiguracji.
package core

import (
	"context"

	"danacoconsole/shared"
)

// ProwenancjaKonfiguracji jest portem rodziny `config.*` rozszerzonym o dwie
// komendy, które sięgają po DROGĘ TURY zamiast po rejestr ustawień.
type ProwenancjaKonfiguracji interface {
	Ustawienia

	// ProwenancjaWywolania obsługuje `config.explain.get`: oddaje wiersz
	// wywołania, prompt systemowy, ustawienia procesu i sumę kontrolną
	// konstytucji dla najbliższej tury wskazanego okna.
	ProwenancjaWywolania(ctx context.Context,
		z shared.ConfigExplainGetRequest) (shared.ConfigExplainGetResponse, error)

	// OtworzOknoKonfiguracji obsługuje `config.window.open`: rozstrzyga zakres,
	// na którym okno konfiguracji staje na wejściu.
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
