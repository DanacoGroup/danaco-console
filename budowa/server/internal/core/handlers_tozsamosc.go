// Odpowiedzialność pliku: wpięcie obszaru zasad i tożsamości modelu.
//
// Tożsamość jest zamieniana, nie dołączana. Silnik nakładki to umie —
// `internal/injection` ma TrybZastap (`--system-prompt`) obok TrybDopisz
// (`--append-system-prompt`), a jego test wymusza ich rozdział. Ta rodzina
// komend nie powtarza silnika; daje mu sterowanie z danych: katalog kategorii,
// treść per oś i nakładkę obowiązującą wraz z trybem.
//
// Wartością domyślną klucza `tozsamosc.tryb_domyslny` jest TrybZastap.
package core

import (
	"context"

	"danacoconsole/shared"
)

// KluczTozsamoscTrybDomyslny wskazuje ustawienie rozstrzygające tryb podania
// nakładki, gdy kategoria nie rozstrzyga go sama.
const KluczTozsamoscTrybDomyslny = "tozsamosc.tryb_domyslny"

// Tozsamosc jest portem obszaru zasad i tożsamości modelu.
type Tozsamosc interface {
	Kategorie(ctx context.Context, z shared.IdentityCategoryListRequest) (shared.IdentityCategoryListResponse, error)
	Dokumenty(ctx context.Context, z shared.IdentityDocumentGetRequest) (shared.IdentityDocumentGetResponse, error)
	Zapisz(ctx context.Context, z shared.IdentityDocumentSetRequest) (shared.IdentityDocumentSetResponse, error)
	Usun(ctx context.Context, z shared.IdentityDocumentRemoveRequest) (shared.IdentityDocumentRemoveResponse, error)
	Obowiazujaca(ctx context.Context, z shared.IdentityEffectiveGetRequest) (shared.IdentityEffectiveGetResponse, error)
}

// zarejestrujTozsamosc wpina pięć komend obszaru tożsamości modelu.
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

// tozsamosc rozgłasza zmianę treści kategorii zasad. Treść obowiązuje platformę,
// model albo konto — nie sesję, więc zdarzenie idzie bez jej wskazania.
//
// Zdarzenie niesie treść, bo okno konfiguracji otwarte na drugim urządzeniu ma
// ją pokazać bez dopytywania. Nie jest to sekret: prompt systemowy jest zasadą
// pracy modelu, nie poświadczeniem.
func (e *emiter) tozsamosc(zmiana shared.ChangeKind, d shared.IdentityDocument) {
	e.wyslij(shared.EventIdentityChanged, "", shared.IdentityChangedEvent{Change: zmiana, Document: d})
}
