// Odpowiedzialność pliku: port i wpięcie dwóch komend rodziny `launcher.*` —
// skrótu globalnego otwierającego wywoływacz poleceń.
//
// Rodzina jest przekrojowa: wywoływacz otwiera się z DOWOLNEGO miejsca
// platformy, a nie z jednego okna. Nastawa jest więc własnością rdzenia, choć
// samo przechwycenie klawiszy należy do powłoki programu okiennego — adapter
// (`adapter_wywolywacz.go`) rozstrzyga ten podział i mówi o nim wprost
// w odpowiedzi.
//
// Zdarzeń rodzina nie ma: zmiana skrótu jest zmianą nastawy, a o zmianach
// nastaw mówi rodzina `config.*`. Osobne zdarzenie byłoby drugą drogą tej samej
// wiadomości.
//
// Port niewypełniony nie rejestruje niczego: obie komendy odpowiedzą wtedy
// `launcher.unknown`, a pozostałe domeny pracują bez zmian.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Wywolywacz jest portem rodziny `launcher.*`.
type Wywolywacz interface {
	// SkrotWywolywacza obsługuje `launcher.hotkey.get`.
	SkrotWywolywacza(ctx context.Context, z shared.LauncherHotkeyGetRequest) (shared.LauncherHotkeyGetResponse, error)
	// ZapiszSkrotWywolywacza obsługuje `launcher.hotkey.set`.
	ZapiszSkrotWywolywacza(ctx context.Context, z shared.LauncherHotkeySetRequest) (shared.LauncherHotkeySetResponse, error)
}

// Adapter wypełnia port w całości.
var _ Wywolywacz = (*adapterWywolywacza)(nil)

// zarejestrujWywolywacz wpina dwie komendy rodziny `launcher.*`.
func zarejestrujWywolywacz(r *Rejestr, w Wywolywacz) {
	if r == nil || w == nil {
		return
	}
	r.Zarejestruj(shared.CommandLauncherHotkeyGet, obsluz(w.SkrotWywolywacza))
	r.Zarejestruj(shared.CommandLauncherHotkeySet, obsluz(w.ZapiszSkrotWywolywacza))
}
