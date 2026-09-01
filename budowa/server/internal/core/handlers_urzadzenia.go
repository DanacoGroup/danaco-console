// Plik wpina rodzinę device.* do rejestru komend. Unieważnienie rozgłasza zmianę do wszystkich
// połączonych urządzeń, nie tylko do tego, które je zleciło.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Urzadzenia jest portem rodziny device.* obsługującym wykaz i unieważnienie urządzeń tego konta użytkownika.
type Urzadzenia interface {
	// WykazUrzadzen obsługuje `device.list`.
	WykazUrzadzen(ctx context.Context, z shared.DeviceListRequest) (shared.DeviceListResponse, error)
	// UniewaznijUrzadzenie obsługuje `device.revoke`.
	UniewaznijUrzadzenie(ctx context.Context,
		z shared.DeviceRevokeRequest) (shared.DeviceRevokeResponse, error)
}

// Zgodność adaptera z portem sprawdzana jest przy kompilacji, bez próby wykonania kodu rdzenia platformy.
var _ Urzadzenia = (*adapterUrzadzen)(nil)

// zWiezia przyjmuje więź połączeń z bramką, dołożoną asercją zamiast argumentu konstruktora, bo
// więź powstaje w kompozycji, a port składa się piętro niżej, w montażu.
type zWiezia interface {
	przyjmijWiez(w *wiezBramki)
}

func (a *adapterUrzadzen) przyjmijWiez(w *wiezBramki) { a.wiez = w }

// zarejestrujUrzadzenia wpina dwie komendy rodziny device.* obsługujące wykaz i unieważnienie urządzenia.
func zarejestrujUrzadzenia(r *Rejestr, u Urzadzenia, e *emiter, wiez *wiezBramki) {
	if r == nil || u == nil {
		return
	}
	if przyjmujacy, niesie := u.(zWiezia); niesie {
		przyjmujacy.przyjmijWiez(wiez)
	}

	r.Zarejestruj(shared.CommandDeviceList,
		obsluz(func(ctx context.Context, z shared.DeviceListRequest) (shared.DeviceListResponse, error) {
			return u.WykazUrzadzen(ctx, z)
		}))

	r.Zarejestruj(shared.CommandDeviceRevoke,
		obsluz(func(ctx context.Context, z shared.DeviceRevokeRequest) (shared.DeviceRevokeResponse, error) {
			odpowiedz, err := u.UniewaznijUrzadzenie(ctx, z)
			if err != nil {
				return odpowiedz, err
			}
			// Wykaz do zdarzenia bierze się z tego samego źródła co odpowiedź wykazu urządzeń.
			wykaz, bladWykazu := u.WykazUrzadzen(ctx, shared.DeviceListRequest{})
			if bladWykazu == nil {
				e.wyslijDoKonta(ctx, shared.EventDeviceChanged, "", shared.DeviceChangedEvent{
					Devices:  wykaz.Devices,
					DeviceId: &z.DeviceId,
				})
			}
			return odpowiedz, err
		}))
}
