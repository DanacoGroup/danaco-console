package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujNawigacje wpina ścieżkę wejścia do platformy: stronę główną, środowiska, moduły,
// przestrzeń roboczą i stan okna komunikacji, gdzie każde ogniwo odpowiada osobno.
func zarejestrujNawigacje(r *Rejestr, n Nawigacja, e *emiter) {
	if r == nil || n == nil {
		return
	}

	r.Zarejestruj(shared.CommandHomeEnter, obsluz(n.StronaGlowna))

	r.Zarejestruj(shared.CommandEnvironmentList, obsluz(n.Srodowiska))

	r.Zarejestruj(shared.CommandEnvironmentEnter, obsluz(n.WejdzDoSrodowiska))

	r.Zarejestruj(shared.CommandModuleList, obsluz(n.Moduly))

	// Wejście do modułu zmienia okno komunikacji; zdarzenie rozgłasza zmianę do innych połączeń konta.
	r.Zarejestruj(shared.CommandWorkspaceEnter,
		obsluz(func(ctx context.Context, z shared.WorkspaceEnterRequest) (shared.WorkspaceEnterResponse, error) {
			w, err := n.WejdzDoPrzestrzeni(ctx, z)
			if err == nil && w.Window.Id != "" {
				e.okno(ctx, rodzajZmianyOkna(z, w), w.Window)
			}
			return w, err
		}))

	// Stan okna należy do ścieżki nawigacji: pytanie brzmi „czy mogę zlecić pracę temu oknu".
	r.Zarejestruj(shared.CommandWindowStateGet, obsluz(n.StanOkna))
}

// rodzajZmianyOkna rozstrzyga, czy wejście do modułu przestawiło okno wskazane,
// czy założyło nowe. Rozróżnienie bierze się ze wskazania w żądaniu: okno
// wskazane i zwrócone pod tym samym identyfikatorem to zmiana, każdy inny
// wynik to okno założone.
func rodzajZmianyOkna(z shared.WorkspaceEnterRequest, w shared.WorkspaceEnterResponse) shared.ChangeKind {
	if z.WindowId != nil && *z.WindowId == w.Window.Id {
		return shared.ChangeKindUpdated
	}
	return shared.ChangeKindCreated
}
