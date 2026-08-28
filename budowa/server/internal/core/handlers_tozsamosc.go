// Plik wpina obszar zasad i tożsamości modelu. Tożsamość jest zamieniana, nie dołączana: ta
// rodzina komend daje silnikowi nakładki sterowanie z danych — katalog kategorii, treść per oś
// i nakładkę obowiązującą wraz z trybem.
package core

import (
	"context"

	"danacoconsole/shared"
)

// KluczTozsamoscTrybDomyslny wskazuje ustawienie rozstrzygające tryb podania
// nakładki, gdy kategoria nie rozstrzyga go sama.
const KluczTozsamoscTrybDomyslny = "tozsamosc.tryb_domyslny"

// Tozsamosc jest portem obszaru zasad i tożsamości modelu obsługującym katalog kategorii i treść każdej osi.
type Tozsamosc interface {
	Kategorie(ctx context.Context, z shared.IdentityCategoryListRequest) (shared.IdentityCategoryListResponse, error)
	Dokumenty(ctx context.Context, z shared.IdentityDocumentGetRequest) (shared.IdentityDocumentGetResponse, error)
	Zapisz(ctx context.Context, z shared.IdentityDocumentSetRequest) (shared.IdentityDocumentSetResponse, error)
	Usun(ctx context.Context, z shared.IdentityDocumentRemoveRequest) (shared.IdentityDocumentRemoveResponse, error)
	Obowiazujaca(ctx context.Context, z shared.IdentityEffectiveGetRequest) (shared.IdentityEffectiveGetResponse, error)
}

// zarejestrujTozsamosc wpina pięć komend obszaru tożsamości modelu w rejestrze rdzenia platformy konta.
func zarejestrujTozsamosc(r *Rejestr, tozsamosc Tozsamosc, e *emiter) {
	if r == nil || tozsamosc == nil {
		return
	}
	r.Zarejestruj(shared.CommandIdentityCategoryList, obsluz(tozsamosc.Kategorie))
	r.Zarejestruj(shared.CommandIdentityDocumentGet, obsluz(tozsamosc.Dokumenty))
	r.Zarejestruj(shared.CommandIdentityEffectiveGet, obsluz(tozsamosc.Obowiazujaca))

	r.Zarejestruj(shared.CommandIdentityDocumentSet,
		obsluz(func(ctx context.Context, z shared.IdentityDocumentSetRequest) (shared.IdentityDocumentSetResponse, error) {
			w, err := tozsamosc.Zapisz(ctx, z)
			if err == nil {
				e.tozsamosc(shared.ChangeKindUpdated, w.Document)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandIdentityDocumentRemove,
		obsluz(func(ctx context.Context, z shared.IdentityDocumentRemoveRequest) (shared.IdentityDocumentRemoveResponse, error) {
			w, err := tozsamosc.Usun(ctx, z)
			if err == nil && w.Removed {
				e.tozsamosc(shared.ChangeKindDeleted, shared.IdentityDocument{Id: z.DocumentId})
			}
			return w, err
		}))
}

// tozsamosc rozgłasza zmianę treści kategorii zasad, obowiązującej platformę, model albo konto,
// bez wskazania sesji; okno konfiguracji na drugim urządzeniu ma ją pokazać bez dopytywania.
func (e *emiter) tozsamosc(zmiana shared.ChangeKind, d shared.IdentityDocument) {
	e.wyslij(shared.EventIdentityChanged, "", shared.IdentityChangedEvent{Change: zmiana, Document: d})
}
