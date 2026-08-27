// Plik wpina port i dwie komendy rodziny launcher.*: odczyt i zapis skrótu globalnego
// otwierającego wywoływacz poleceń z dowolnego miejsca platformy.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Wywolywacz jest portem rodziny launcher.* obsługującym odczyt i zapis skrótu globalnego tej platformy.
type Wywolywacz interface {
	// SkrotWywolywacza obsługuje `launcher.hotkey.get`.
	SkrotWywolywacza(ctx context.Context, z shared.LauncherHotkeyGetRequest) (shared.LauncherHotkeyGetResponse, error)
	// ZapiszSkrotWywolywacza obsługuje `launcher.hotkey.set`.
	ZapiszSkrotWywolywacza(ctx context.Context, z shared.LauncherHotkeySetRequest) (shared.LauncherHotkeySetResponse, error)
}

// Adapter wypełnia port w całości, bez pozostawionej metody niezaimplementowanej w rdzeniu tej platformy.
var _ Wywolywacz = (*adapterWywolywacza)(nil)

// zarejestrujWywolywacz wpina dwie komendy rodziny launcher.* w rejestrze rdzenia tej platformy konta.
func zarejestrujWywolywacz(r *Rejestr, w Wywolywacz) {
	if r == nil || w == nil {
		return
	}
	r.Zarejestruj(shared.CommandLauncherHotkeyGet, obsluz(w.SkrotWywolywacza))
	r.Zarejestruj(shared.CommandLauncherHotkeySet, obsluz(w.ZapiszSkrotWywolywacza))
}
