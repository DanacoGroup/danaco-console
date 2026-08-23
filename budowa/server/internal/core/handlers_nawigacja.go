package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujNawigacje wpina ścieżkę wejścia do platformy: stronę główną,
// środowiska, moduły, przestrzeń roboczą i stan okna komunikacji.
//
// Ciąg jest jeden: `home.enter` → `environment.list`
// / `environment.enter` → `module.list` → `workspace.enter` →
// `window.state.get`. Każde ogniwo odpowiada osobno, więc klient może wejść
// w dowolnym miejscu — na przykład wprost do środowiska zapamiętanego z
// poprzedniej pracy.
//
// Brak podłączonej domeny nie wywraca rdzenia: komendy nawigacji odpowiedzą
// wtedy `home.unknown`, `environment.unknown`, `module.unknown` albo
// `workspace.unknown`, a pozostałe domeny pracują dalej.
func zarejestrujNawigacje(r *Rejestr, n Nawigacja, e *emiter) {
	if r == nil || n == nil {
		return
	}

	r.Zarejestruj(shared.CommandHomeEnter, obsluz(n.StronaGlowna))

	r.Zarejestruj(shared.CommandEnvironmentList, obsluz(n.Srodowiska))

	r.Zarejestruj(shared.CommandEnvironmentEnter, obsluz(n.WejdzDoSrodowiska))

	r.Zarejestruj(shared.CommandModuleList, obsluz(n.Moduly))

	// Wejście do modułu zmienia okno komunikacji: to samo okno dostaje inny
	// moduł albo powstaje okno nowe. Bez rozgłoszenia tej zmiany okno rozmowy
	// otwarte w innym połączeniu tego konta stałoby przy module poprzednim —
	// pasek narzędzi promptu, panel akcji i kontekst pokazywałyby moduł, z
	// którego Operator już wyszedł. Zdarzenie jest jedynym nośnikiem zmiany
	// (zob. emiter), więc rozgłasza je uchwyt, nie adapter.
	r.Zarejestruj(shared.CommandWorkspaceEnter,
		obsluz(func(ctx context.Context, z shared.WorkspaceEnterRequest) (shared.WorkspaceEnterResponse, error) {
			w, err := n.WejdzDoPrzestrzeni(ctx, z)
			if err == nil && w.Window.Id != "" {
				e.okno(ctx, rodzajZmianyOkna(z, w), w.Window)
			}
			return w, err
		}))

	// Stan okna należy do ścieżki nawigacji, nie do domeny okien: pytanie brzmi
	// „czy mogę zlecić pracę temu oknu", a nie „zmień jego parametry".
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
