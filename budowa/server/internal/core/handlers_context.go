package core

import (
	"context"

	"danacoconsole/shared"
)

// ZajetoscKontekstu jest portem `context.usage.get`, portem pomiaru zajętości
// okna kontekstu rozmowy, osobnym od portu `Przenoszenie`.
type ZajetoscKontekstu interface {
	ZajetoscKontekstu(ctx context.Context, z shared.ContextUsageGetRequest) (shared.ContextUsageGetResponse, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ ZajetoscKontekstu = (*adapterZajetosciKontekstu)(nil)

// zarejestrujZajetoscKontekstu wpina `context.usage.get`.
//
// Komenda niczego nie rozgłasza: pomiar jest odczytem, a zmianę zajętości
// wywołuje tura, o której mówi już rodzina `message.*`.
func zarejestrujZajetoscKontekstu(r *Rejestr, z ZajetoscKontekstu) {
	if r == nil || z == nil {
		return
	}
	r.Zarejestruj(shared.CommandContextUsageGet, obsluz(z.ZajetoscKontekstu))
}

// zarejestrujPrzenoszenie wpina przekazanie kompletu kontekstu między
// modułami jedną komendą, obejmującą polecenie, dokumenty, projekt, agentów,
// historię, źródła wiedzy i parametry wykonania.
func zarejestrujPrzenoszenie(r *Rejestr, przenoszenie Przenoszenie, e *emiter) {
	if r == nil || przenoszenie == nil {
		return
	}

	r.Zarejestruj(shared.CommandContextTransfer,
		obsluz(func(ctx context.Context, z shared.ContextTransferRequest) (shared.ContextTransferResponse, error) {
			w, err := przenoszenie.Przenies(ctx, z)
			if err == nil && w.Transferred {
				e.okno(ctx, shared.ChangeKindUpdated, w.Window)
			}
			return w, err
		}))
}
