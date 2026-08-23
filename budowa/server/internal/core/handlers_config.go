package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujUstawienia wpina domenę konfiguracji warstwowej dziewięciu
// poziomów zasięgu — od zasięgu aplikacji, najszerszego, po okno komunikacji,
// najwęższe i wygrywające z pozostałymi.
//
// Rdzeń nie rozstrzyga tu pierwszeństwa poziomów ani nie zna katalogu ustawień —
// to należy do warstwy konfiguracji. Rdzeń wyłącznie kieruje komendę i rozgłasza
// zmianę.
//
// config.reset przywraca wartość domyślną. Skutkiem jest usunięcie ustawienia
// z poziomu, więc zmiana idzie jako usunięcie wpisu; brak ustawienia znaczy
// wartość domyślną, nigdy blokadę.
func zarejestrujUstawienia(r *Rejestr, ustawienia Ustawienia, e *emiter) {
	if r == nil || ustawienia == nil {
		return
	}

	// Obszary jednolitego modelu konfiguracji sesji idą tym samym portem —
	// jedna brama do konfiguracji, nie dwie.
	zarejestrujKonfiguracjeSesji(r, ustawienia, e)

	r.Zarejestruj(shared.CommandConfigGet, obsluz(ustawienia.Odczytaj))

	r.Zarejestruj(shared.CommandConfigSet,
		obsluz(func(ctx context.Context, z shared.ConfigSetRequest) (shared.ConfigSetResponse, error) {
			w, err := ustawienia.Zapisz(ctx, z)
			if err == nil {
				e.ustawienie(shared.ChangeKindUpdated, w.Entry)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandConfigReset,
		obsluz(func(ctx context.Context, z shared.ConfigResetRequest) (shared.ConfigResetResponse, error) {
			w, err := ustawienia.Przywroc(ctx, z)
			if err != nil {
				return w, err
			}
			for _, wpis := range w.Entries {
				e.ustawienie(shared.ChangeKindDeleted, wpis)
			}
			return w, nil
		}))
}
