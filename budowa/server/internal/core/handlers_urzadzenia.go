// Odpowiedzialność pliku: wpięcie rodziny `device.*` do rejestru komend.
//
// Unieważnienie rozgłasza `device.changed` do WSZYSTKICH połączonych urządzeń,
// nie tylko do tego, które je zleciło. Operator odbierający dostęp maszynie
// stojącej obok ma zobaczyć skutek na obu ekranach naraz — inaczej drugi ekran
// pokazywałby dostęp, którego już nie ma, aż do następnej komendy.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Urzadzenia jest portem rodziny `device.*`.
type Urzadzenia interface {
	// WykazUrzadzen obsługuje `device.list`.
	WykazUrzadzen(ctx context.Context, z shared.DeviceListRequest) (shared.DeviceListResponse, error)
	// UniewaznijUrzadzenie obsługuje `device.revoke`.
	UniewaznijUrzadzenie(ctx context.Context,
		z shared.DeviceRevokeRequest) (shared.DeviceRevokeResponse, error)
}

// Zgodność adaptera z portem sprawdzana jest przy kompilacji.
var _ Urzadzenia = (*adapterUrzadzen)(nil)

// zWiezia przyjmuje więź połączeń z bramką.
//
// Więź powstaje w kompozycji, a port składa się piętro niżej, w montażu — stąd
// dołożenie przez asercję zamiast argumentu konstruktora. Port, który więzi nie
// przyjmie, po prostu nie oznaczy bieżącego urządzenia; wykaz działa dalej.
type zWiezia interface {
	przyjmijWiez(w *wiezBramki)
}

func (a *adapterUrzadzen) przyjmijWiez(w *wiezBramki) { a.wiez = w }

// zarejestrujUrzadzenia wpina dwie komendy rodziny `device.*`.
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
			// Wykaz do zdarzenia bierze się z tego samego źródła, co odpowiedź
			// `device.list` — pozostałe ekrany dostają stan po zmianie, a nie
			// polecenie „odpytaj jeszcze raz".
			wykaz, bladWykazu := u.WykazUrzadzen(ctx, shared.DeviceListRequest{})
			if bladWykazu == nil {
				e.wyslij(shared.EventDeviceChanged, "", shared.DeviceChangedEvent{
					Devices:  wykaz.Devices,
					DeviceId: &z.DeviceId,
				})
			}
			return odpowiedz, err
		}))
}
