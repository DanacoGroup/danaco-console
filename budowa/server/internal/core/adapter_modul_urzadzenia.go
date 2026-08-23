// Odpowiedzialność pliku: wykaz urządzeń powiązanych z kontem właściciela
// (`device.list`) i unieważnienie tokenu wskazanego urządzenia (`device.revoke`).
//
// ── SKĄD BIERZE SIĘ WYKAZ ────────────────────────────────────────────────────
// Urządzeniem konta jest to, które kiedykolwiek weszło przez bramkę — a to
// wiedzą sesje bramki. Osobnej tabeli urządzeń nie ma z rozmysłem: wymagałaby
// sprzątania wierszy, których nic już nie dotyczy, i rozjeżdżałaby się z prawdą
// przy pierwszym unieważnieniu, o którym ktoś zapomniałby ją powiadomić.
//
// ── DLACZEGO „TO URZĄDZENIE" JEST POLEM DANYCH ──────────────────────────────
// Wiersz własnego urządzenia wygląda w wykazie tak samo jak każdy inny, więc bez
// oznaczenia Operator odbiera dostęp sobie i traci go w tej samej chwili.
// Rozstrzygnięcie należy do rdzenia, nie do klienta: klient zna identyfikator,
// który sam nadał, ale nie wie, którą sesją stoi połączenie.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// adapterUrzadzen wypełnia port Urzadzenia.
type adapterUrzadzen struct {
	repozytorium dane.RepozytoriumUwierzytelnienia
	wiez         *wiezBramki
}

func nowyAdapterUrzadzen(repozytorium dane.RepozytoriumUwierzytelnienia) *adapterUrzadzen {
	return &adapterUrzadzen{repozytorium: repozytorium}
}

// ── device.list ──────────────────────────────────────────────────────────────

// WykazUrzadzen obsługuje `device.list`.
func (a *adapterUrzadzen) WykazUrzadzen(ctx context.Context,
	_ shared.DeviceListRequest) (shared.DeviceListResponse, error) {

	teraz := time.Now().UnixMilli()
	wiersze, err := a.repozytorium.UrzadzeniaKonta(ctx, teraz)
	if err != nil {
		return shared.DeviceListResponse{}, err
	}
	biezace := a.urzadzenieBiezace(ctx)

	wykaz := make([]shared.Device, 0, len(wiersze))
	for _, wiersz := range wiersze {
		ostatnio := wiersz.OstatnioWidziane
		wykaz = append(wykaz, shared.Device{
			DeviceId:   wiersz.Kod,
			LastSeenAt: &ostatnio,
			Current:    biezace != "" && wiersz.Kod == biezace,
			HasToken:   wiersz.MaToken,
		})
	}
	return shared.DeviceListResponse{Devices: wykaz}, nil
}

// ── device.revoke ────────────────────────────────────────────────────────────

// UniewaznijUrzadzenie obsługuje `device.revoke`.
//
// Unieważnienie własnego urządzenia jest dozwolone i nie jest pomyłką: Operator
// bywa przy cudzej maszynie i zamyka na niej swój dostęp świadomie. Klient wie,
// które urządzenie jest bieżące (pole `current` wykazu), więc ostrzeżenie należy
// do ekranu, a nie do odmowy rdzenia — blokada tutaj byłaby bramkowaniem.
func (a *adapterUrzadzen) UniewaznijUrzadzenie(ctx context.Context,
	z shared.DeviceRevokeRequest) (shared.DeviceRevokeResponse, error) {

	urzadzenie := strings.TrimSpace(z.DeviceId)
	if urzadzenie == "" {
		return shared.DeviceRevokeResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"unieważnienie bez wskazania urządzenia")
	}
	zamkniete, err := a.repozytorium.UniewaznijSesjeUrzadzenia(ctx, urzadzenie, time.Now().UnixMilli())
	if err != nil {
		return shared.DeviceRevokeResponse{}, err
	}
	// Zero zamkniętych sesji NIE jest odmową: urządzenie mogło już nie mieć
	// ważnego tokenu, a skutek żądany przez Operatora — „to urządzenie nie ma
	// dostępu" — i tak obowiązuje.
	return shared.DeviceRevokeResponse{Revoked: zamkniete > 0}, nil
}

// urzadzenieBiezace odczytuje urządzenie sesji, którą stoi to połączenie.
//
// Pusty wynik znaczy „rdzeń nie wie" i wykaz nie oznacza wtedy żadnego wiersza
// jako bieżącego. Zgadywanie po ostatnim wejściu byłoby gorsze niż milczenie:
// wskazałoby cudzą maszynę jako własną.
func (a *adapterUrzadzen) urzadzenieBiezace(ctx context.Context) string {
	if a.wiez == nil {
		return ""
	}
	skrot := a.wiez.SkrotKontekstu(ctx)
	if skrot == "" {
		return ""
	}
	sesja, err := a.repozytorium.SesjaBramkiPoSkrocie(ctx, skrot)
	if err != nil || sesja.UrzadzenieKod == nil {
		return ""
	}
	return *sesja.UrzadzenieKod
}
