// Odpowiedzialność pliku: wpięcie rejestru kont modeli i kont programów code
// CLI — tworzenie i dodawanie kont z okna konfiguracji.
//
// KONTO A KANAŁ. Konto jest profilem uwierzytelnienia, kanał — definicją
// rozmowy z modelem. Jeden kanał wskazuje konto preferowane, jedno konto
// obsługuje wiele kanałów, a pula rotacji bierze konta tego samego rodzaju.
// Dlatego skasowanie konta ODŁĄCZA kanały, ale ich nie kasuje.
//
// Poświadczenie wchodzi żądaniem i nie wychodzi żadną drogą: ani
// wykazem, ani odpowiedzią na zapis, ani zdarzeniem zmiany. Odpowiedź niesie
// wyłącznie znacznik `hasCredential`.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Konta jest portem rejestru kont.
type Konta interface {
	Dodaj(ctx context.Context, z shared.AccountAddRequest) (shared.AccountAddResponse, error)
	Wykaz(ctx context.Context, z shared.AccountListRequest) (shared.AccountListResponse, error)
	Zmien(ctx context.Context, z shared.AccountUpdateRequest) (shared.AccountUpdateResponse, error)
	Usun(ctx context.Context, z shared.AccountRemoveRequest) (shared.AccountRemoveResponse, error)
	UstawDomyslne(ctx context.Context, z shared.AccountDefaultSetRequest) (shared.AccountDefaultSetResponse, error)
}

// zarejestrujKonta wpina pięć komend rejestru kont.
func zarejestrujKonta(r *Rejestr, konta Konta, e *emiter) {
	if r == nil || konta == nil {
		return
	}
	r.Zarejestruj(shared.CommandAccountList, obsluz(konta.Wykaz))

	r.Zarejestruj(shared.CommandAccountAdd,
		obsluz(func(ctx context.Context, z shared.AccountAddRequest) (shared.AccountAddResponse, error) {
			w, err := konta.Dodaj(ctx, z)
			if err == nil {
				e.konto(shared.ChangeKindCreated, w.Account)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAccountUpdate,
		obsluz(func(ctx context.Context, z shared.AccountUpdateRequest) (shared.AccountUpdateResponse, error) {
			w, err := konta.Zmien(ctx, z)
			if err == nil {
				e.konto(shared.ChangeKindUpdated, w.Account)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAccountRemove,
		obsluz(func(ctx context.Context, z shared.AccountRemoveRequest) (shared.AccountRemoveResponse, error) {
			w, err := konta.Usun(ctx, z)
			if err == nil && w.Removed {
				e.konto(shared.ChangeKindDeleted, shared.Account{Id: z.AccountId})
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAccountDefaultSet,
		obsluz(func(ctx context.Context, z shared.AccountDefaultSetRequest) (shared.AccountDefaultSetResponse, error) {
			w, err := konta.UstawDomyslne(ctx, z)
			if err == nil {
				e.konto(shared.ChangeKindUpdated, w.Account)
			}
			return w, err
		}))
}

// konto rozgłasza zmianę konta. Konto nie należy do sesji, więc zdarzenie idzie
// bez jej wskazania. Struktura Account nie ma pola na poświadczenie, więc
// zdarzenie nie ma jak go wynieść.
func (e *emiter) konto(zmiana shared.ChangeKind, k shared.Account) {
	e.wyslij(shared.EventAccountChanged, "", shared.AccountChangedEvent{Change: zmiana, Account: k})
}
