package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujUstawienia wpina domenę konfiguracji warstwowej dziewięciu
// poziomów zasięgu, od zasięgu aplikacji po okno komunikacji, i rozgłasza
// zmianę.
func zarejestrujUstawienia(r *Rejestr, ustawienia Ustawienia, e *emiter) {
	if r == nil || ustawienia == nil {
		return
	}

	// Obszary jednolitego modelu konfiguracji sesji idą tym samym portem,
	// jedną bramą do konfiguracji.
	zarejestrujKonfiguracjeSesji(r, ustawienia, e)

	r.Zarejestruj(shared.CommandConfigGet, obsluz(ustawienia.Odczytaj))

	r.Zarejestruj(shared.CommandConfigSet,
		obsluz(func(ctx context.Context, z shared.ConfigSetRequest) (shared.ConfigSetResponse, error) {
			w, err := ustawienia.Zapisz(ctx, z)
			if err == nil {
				e.ustawienie(ctx, shared.ChangeKindUpdated, w.Entry)
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
				e.ustawienie(ctx, shared.ChangeKindDeleted, wpis)
			}
			return w, nil
		}))
}
