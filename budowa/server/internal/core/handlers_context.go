package core

import (
	"context"

	"danacoconsole/shared"
)

// ZajetoscKontekstu jest portem `context.usage.get` — pomiaru zajętości okna
// kontekstu rozmowy.
//
// Port osobny od `Przenoszenie`, mimo wspólnego przedrostka `context.`.
// Przeniesienie kompletu kontekstu między oknami i pomiar zajętości okna nie
// mają ze sobą nic wspólnego poza słowem w nazwie: pierwsze zakłada okno,
// drugie liczy żetony tokenizatorem. Wspólny port związałby ich dostępność
// w jedno „jest albo nie ma".
type ZajetoscKontekstu interface {
	ZajetoscKontekstu(ctx context.Context, z shared.ContextUsageGetRequest) (shared.ContextUsageGetResponse, error)
}

// Adapter wypełnia port w całości.
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

// zarejestrujPrzenoszenie wpina przekazanie kompletu kontekstu między modułami
// jedną komendą: polecenie, dokumenty, projekt, agenci, historia,
// źródła wiedzy i parametry wykonania idą razem.
//
// Jedna komenda, nie ścieżka per moduł — dlatego rdzeń nie rozgałęzia się tu na
// moduł docelowy. Przeniesienie kończy się oknem docelowym, więc rdzeń rozgłasza
// zmianę tego okna: klient dowiaduje się o nowym oknie tą samą drogą, co przy
// window.create.
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
