// Plik wpina rejestr kont modeli i kont programów wiersza poleceń, obsługując
// tworzenie, zmianę, usuwanie oraz ustawienie konta domyślnego z okna
// konfiguracji.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Konta jest portem rejestru kont: profili uwierzytelnienia współdzielonych
// przez kanały rozmowy z modelem oraz przez pulę rotacji tego samego rodzaju.
type Konta interface {
	Dodaj(ctx context.Context, z shared.AccountAddRequest) (shared.AccountAddResponse, error)
	Wykaz(ctx context.Context, z shared.AccountListRequest) (shared.AccountListResponse, error)
	Zmien(ctx context.Context, z shared.AccountUpdateRequest) (shared.AccountUpdateResponse, error)
	Usun(ctx context.Context, z shared.AccountRemoveRequest) (shared.AccountRemoveResponse, error)
	UstawDomyslne(ctx context.Context, z shared.AccountDefaultSetRequest) (shared.AccountDefaultSetResponse, error)
}

// zarejestrujKonta wpina pięć komend rejestru kont i rozgłasza zdarzenie
// zmiany konta po każdym udanym zapisie.
func zarejestrujKonta(r *Rejestr, konta Konta, e *emiter) {
	if r == nil || konta == nil {
		return
	}
	r.Zarejestruj(shared.CommandAccountList, obsluz(konta.Wykaz))

	r.Zarejestruj(shared.CommandAccountAdd,
		obsluz(func(ctx context.Context, z shared.AccountAddRequest) (shared.AccountAddResponse, error) {
			w, err := konta.Dodaj(ctx, z)
			if err == nil {
				e.konto(ctx, shared.ChangeKindCreated, w.Account)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAccountUpdate,
		obsluz(func(ctx context.Context, z shared.AccountUpdateRequest) (shared.AccountUpdateResponse, error) {
			w, err := konta.Zmien(ctx, z)
			if err == nil {
				e.konto(ctx, shared.ChangeKindUpdated, w.Account)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAccountRemove,
		obsluz(func(ctx context.Context, z shared.AccountRemoveRequest) (shared.AccountRemoveResponse, error) {
			w, err := konta.Usun(ctx, z)
			if err == nil && w.Removed {
				e.konto(ctx, shared.ChangeKindDeleted, shared.Account{Id: z.AccountId})
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAccountDefaultSet,
		obsluz(func(ctx context.Context, z shared.AccountDefaultSetRequest) (shared.AccountDefaultSetResponse, error) {
			w, err := konta.UstawDomyslne(ctx, z)
			if err == nil {
				e.konto(ctx, shared.ChangeKindUpdated, w.Account)
			}
			return w, err
		}))
}

// konto rozgłasza zmianę konta. Konto nie należy do sesji, więc zdarzenie idzie
// bez jej wskazania. Struktura Account nie ma pola na poświadczenie, więc
// zdarzenie nie ma jak go wynieść.
func (e *emiter) konto(ctx context.Context, zmiana shared.ChangeKind, k shared.Account) {
	e.wyslijDoKonta(ctx, shared.EventAccountChanged, "", shared.AccountChangedEvent{Change: zmiana, Account: k})
}
