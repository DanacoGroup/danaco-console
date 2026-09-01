// Odpowiedzialność pliku: wykaz urządzeń powiązanych z kontem właściciela (`device.list`)
// i unieważnienie tokenu wskazanego urządzenia (`device.revoke`), złożone z sesji bramki.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// adapterUrzadzen wypełnia port Urzadzenia, złożony z sesji bramki bez osobnej tabeli urządzeń tego konta.
type adapterUrzadzen struct {
	repozytorium dane.RepozytoriumUwierzytelnienia
	wiez         *wiezBramki
	// rozlaczanie zrywa gniazda urządzenia po unieważnieniu; zerowe zostawia je bramce.
	rozlaczanie RozlaczanieSesji
}

func (a *adapterUrzadzen) przyjmijRozlaczanie(r RozlaczanieSesji) { a.rozlaczanie = r }

func nowyAdapterUrzadzen(repozytorium dane.RepozytoriumUwierzytelnienia) *adapterUrzadzen {
	return &adapterUrzadzen{repozytorium: repozytorium}
}

// WykazUrzadzen obsługuje `device.list` i oddaje wykaz urządzeń wynikający wprost z sesji bramki tego konta.
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

// UniewaznijUrzadzenie obsługuje `device.revoke`; unieważnienie własnego urządzenia jest
// dozwolone i nie jest pomyłką, bo klient sam ostrzega, gdy dotyczy urządzenia bieżącego.
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
	// Gniazda urządzenia zrywa się od razu; inaczej unieważnienie skutkuje dopiero przy jego następnej komendzie.
	if a.rozlaczanie != nil && zamkniete > 0 {
		a.rozlaczanie.RozlaczPoUniewaznieniu(ctx, kontoAdresata(ctx), "sesja bramki urządzenia unieważniona")
	}
	// Zero zamkniętych sesji nie jest odmową: skutek żądany przez Operatora i tak obowiązuje.
	return shared.DeviceRevokeResponse{Revoked: zamkniete > 0}, nil
}

// urzadzenieBiezace odczytuje urządzenie sesji, którą stoi to połączenie; pusty wynik
// znaczy, że rdzeń nie wie, zamiast zgadywać cudzą maszynę jako własną.
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
