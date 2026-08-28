// Moduł Library — udostępnienia odnośnikiem i nasłuchy zewnętrzne:
// udostępnienie, wykaz, odwołanie oraz nasłuch webhook — ustawienie, wykaz i
// usunięcie. Token jest jawny i losowy kryptograficznie.
package core

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// dlugoscTokenuUdostepnienia — 32 bajty losowości, czyli 64 znaki zapisu
// szesnastkowego. Tyle wystarcza, żeby zgadywanie tokenu nie było drogą wejścia.
const dlugoscTokenuUdostepnienia = 32

// WystawUdostepnienie obsługuje `library.share.create`: wystawia token dostępu
// do zasobu albo kolekcji i sprawdza, że wskazany byt istnieje, zanim odnośnik
// powstanie.
func (a *adapterBiblioteki) WystawUdostepnienie(ctx context.Context,
	z shared.LibraryShareCreateRequest) (shared.LibraryShareCreateResponse, error) {

	cel := strings.TrimSpace(z.TargetId)
	if cel == "" {
		return shared.LibraryShareCreateResponse{}, bladWskazaniaBiblioteki(
			"udostępnienie bez wskazania zasobu albo kolekcji")
	}
	// Byt udostępniany jest sprawdzany wprost, żeby odnośnik nie prowadził do
	// zasobu, którego nie ma.
	if z.Scope == shared.LibraryShareScopeCollection {
		wykaz, _, err := a.repozytorium.KolekcjeWykaz(ctx, nil, nil, 0)
		if err != nil {
			return shared.LibraryShareCreateResponse{}, bladBiblioteki(err)
		}
		if !kolekcjaIstniejeBiblioteki(wykaz, cel) {
			return shared.LibraryShareCreateResponse{}, bladNieznanejKolekcji(cel, dane.ErrBrakWiersza)
		}
	} else if _, err := a.plik(ctx, cel); err != nil {
		return shared.LibraryShareCreateResponse{}, err
	}

	token, err := tokenUdostepnieniaBiblioteki()
	if err != nil {
		return shared.LibraryShareCreateResponse{}, bladBiblioteki(err)
	}
	var wygasa *string
	if z.ExpiresInSeconds != nil && *z.ExpiresInSeconds > 0 {
		chwila := time.Now().UTC().Add(time.Duration(*z.ExpiresInSeconds) * time.Second).
			Format(formatZnacznikaBazy)
		wygasa = &chwila
	}
	zapisane, err := a.repozytorium.ZapiszUdostepnienie(ctx, dane.UdostepnienieBiblioteki{
		Kod:    nowyIdentyfikator(przedrostekUdostepnieniaBiblioteki),
		Zasieg: zasiegUdostepnieniaBazy(z.Scope), CelKod: cel, Token: token, Wygasa: wygasa,
	})
	if err != nil {
		return shared.LibraryShareCreateResponse{}, bladBiblioteki(err)
	}
	a.odnotuj(ctx, shared.LibraryAuditActionExport, wskazanieBiblioteki(cel),
		"wystawienie odnośnika udostępnienia")
	return shared.LibraryShareCreateResponse{
		Share: udostepnienieKontraktuBiblioteki(zapisane),
		Url:   adresUdostepnienia(zapisane),
	}, nil
}

// WykazUdostepnien obsługuje `library.share.list`, oddając udostępnienia
// zasobu wskazanego żądaniem, opcjonalnie ograniczone do udostępnień czynnych.
func (a *adapterBiblioteki) WykazUdostepnien(ctx context.Context,
	z shared.LibraryShareListRequest) (shared.LibraryShareListResponse, error) {

	tylkoCzynne := z.ActiveOnly != nil && *z.ActiveOnly
	wiersze, err := a.repozytorium.Udostepnienia(ctx, z.TargetId, tylkoCzynne)
	if err != nil {
		return shared.LibraryShareListResponse{}, bladBiblioteki(err)
	}
	udostepnienia := make([]shared.LibraryShare, 0, len(wiersze))
	for _, wiersz := range wiersze {
		udostepnienia = append(udostepnienia, udostepnienieKontraktuBiblioteki(wiersz))
	}
	return shared.LibraryShareListResponse{Shares: udostepnienia, Total: len(udostepnienia)}, nil
}

// OdwolajUdostepnienie obsługuje `library.share.revoke`, unieważniając
// udostępnienie wskazanego kodu i odnotowując zmianę w dzienniku audytu.
func (a *adapterBiblioteki) OdwolajUdostepnienie(ctx context.Context,
	z shared.LibraryShareRevokeRequest) (shared.LibraryShareRevokeResponse, error) {

	kod := strings.TrimSpace(z.ShareId)
	if kod == "" {
		return shared.LibraryShareRevokeResponse{}, bladWskazaniaBiblioteki(
			"odwołanie bez wskazania udostępnienia")
	}
	odwolane, err := a.repozytorium.OdwolajUdostepnienie(ctx, kod)
	if err != nil {
		return shared.LibraryShareRevokeResponse{}, bladBiblioteki(err)
	}
	if odwolane {
		a.odnotuj(ctx, shared.LibraryAuditActionChange, nil, "odwołanie udostępnienia "+kod)
	}
	return shared.LibraryShareRevokeResponse{Revoked: odwolane}, nil
}

// UstawWebhook obsługuje `library.webhook.set`: wymaga adresu HTTP albo HTTPS
// i co najmniej jednego zdarzenia, na które nasłuch ma reagować.
func (a *adapterBiblioteki) UstawWebhook(ctx context.Context,
	z shared.LibraryWebhookSetRequest) (shared.LibraryWebhookSetResponse, error) {

	adres := strings.TrimSpace(z.Webhook.Url)
	if adres == "" {
		return shared.LibraryWebhookSetResponse{}, bladWskazaniaBiblioteki(
			"nasłuch bez adresu nie ma dokąd zgłaszać zdarzeń")
	}
	if !strings.HasPrefix(adres, "http://") && !strings.HasPrefix(adres, "https://") {
		return shared.LibraryWebhookSetResponse{}, bladWskazaniaBiblioteki(
			"adres nasłuchu musi być adresem HTTP albo HTTPS; podano: " + adres)
	}
	if len(z.Webhook.Events) == 0 {
		return shared.LibraryWebhookSetResponse{}, bladWskazaniaBiblioteki(
			"nasłuch bez ani jednego zdarzenia nigdy by nie zadziałał")
	}
	zdarzenia := make([]string, 0, len(z.Webhook.Events))
	for _, zdarzenie := range z.Webhook.Events {
		zdarzenia = append(zdarzenia, zdarzenieWebhookaBazy(zdarzenie))
	}
	kod := strings.TrimSpace(z.Webhook.Id)
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekWebhookaBiblioteki)
	}
	zapisany, err := a.repozytorium.ZapiszWebhook(ctx, dane.WebhookBiblioteki{
		Kod: kod, Adres: adres, Sekret: z.Webhook.Secret,
		Czynny: z.Webhook.Enabled, Zdarzenia: zdarzenia,
	})
	if err != nil {
		return shared.LibraryWebhookSetResponse{}, bladBiblioteki(err)
	}
	a.odnotuj(ctx, shared.LibraryAuditActionChange, nil, "zapis nasłuchu "+kod)
	return shared.LibraryWebhookSetResponse{Webhook: webhookKontraktuBiblioteki(zapisany)}, nil
}

// WykazWebhookow obsługuje `library.webhook.list`, oddając nasłuchy
// zewnętrzne, opcjonalnie ograniczone do nasłuchów czynnych.
func (a *adapterBiblioteki) WykazWebhookow(ctx context.Context,
	z shared.LibraryWebhookListRequest) (shared.LibraryWebhookListResponse, error) {

	tylkoCzynne := z.EnabledOnly != nil && *z.EnabledOnly
	wiersze, err := a.repozytorium.Webhooki(ctx, tylkoCzynne)
	if err != nil {
		return shared.LibraryWebhookListResponse{}, bladBiblioteki(err)
	}
	nasluchy := make([]shared.LibraryWebhook, 0, len(wiersze))
	for _, wiersz := range wiersze {
		nasluchy = append(nasluchy, webhookKontraktuBiblioteki(wiersz))
	}
	return shared.LibraryWebhookListResponse{Webhooks: nasluchy, Total: len(nasluchy)}, nil
}

// UsunWebhook obsługuje `library.webhook.remove`, usuwając nasłuch wskazanego
// kodu i odnotowując usunięcie w dzienniku audytu.
func (a *adapterBiblioteki) UsunWebhook(ctx context.Context,
	z shared.LibraryWebhookRemoveRequest) (shared.LibraryWebhookRemoveResponse, error) {

	kod := strings.TrimSpace(z.WebhookId)
	if kod == "" {
		return shared.LibraryWebhookRemoveResponse{}, bladWskazaniaBiblioteki(
			"usunięcie bez wskazania nasłuchu")
	}
	usuniety, err := a.repozytorium.UsunWebhook(ctx, kod)
	if err != nil {
		return shared.LibraryWebhookRemoveResponse{}, bladBiblioteki(err)
	}
	if usuniety {
		a.odnotuj(ctx, shared.LibraryAuditActionChange, nil, "usunięcie nasłuchu "+kod)
	}
	return shared.LibraryWebhookRemoveResponse{Removed: usuniety}, nil
}

// tokenUdostepnieniaBiblioteki losuje token dostępu generatorem losowości
// kryptograficznej, zapisany szesnastkowo w liczbie znaków
// `dlugoscTokenuUdostepnienia`.
func tokenUdostepnieniaBiblioteki() (string, error) {
	bajty := make([]byte, dlugoscTokenuUdostepnienia)
	if _, err := rand.Read(bajty); err != nil {
		return "", err
	}
	return hex.EncodeToString(bajty), nil
}

// adresUdostepnienia składa adres, pod którym zasób jest osiągalny, jako
// ścieżkę względną niosącą token udostępnienia.
func adresUdostepnienia(udostepnienie dane.UdostepnienieBiblioteki) string {
	return "/biblioteka/udostepnienie/" + udostepnienie.Token
}

// kolekcjaIstniejeBiblioteki sprawdza obecność kolekcji wskazanego kodu w
// podanym wykazie kolekcji biblioteki.
func kolekcjaIstniejeBiblioteki(wykaz []dane.KolekcjaBiblioteki, kod string) bool {
	for _, kolekcja := range wykaz {
		if kolekcja.Kod == kod {
			return true
		}
	}
	return false
}

// udostepnienieKontraktuBiblioteki przenosi wiersz udostępnienia na kontrakt,
// dokładając chwilę wygaśnięcia i odwołania, gdy wiersz je niesie.
func udostepnienieKontraktuBiblioteki(wiersz dane.UdostepnienieBiblioteki) shared.LibraryShare {
	udostepnienie := shared.LibraryShare{
		Id: wiersz.Kod, Scope: zasiegUdostepnieniaKontraktu(wiersz.Zasieg),
		TargetId: wiersz.CelKod, Token: wiersz.Token,
		CreatedAt: chwilaBazy(wiersz.Utworzono),
	}
	if wiersz.Wygasa != nil && *wiersz.Wygasa != "" {
		chwila := chwilaBazy(*wiersz.Wygasa)
		udostepnienie.ExpiresAt = &chwila
	}
	if wiersz.Odwolano != nil && *wiersz.Odwolano != "" {
		chwila := chwilaBazy(*wiersz.Odwolano)
		udostepnienie.RevokedAt = &chwila
	}
	return udostepnienie
}

// webhookKontraktuBiblioteki przenosi wiersz nasłuchu na kontrakt, dokładając
// zdarzenia subskrybowane i chwilę ostatniego zgłoszenia, gdy wiersz ją
// niesie.
func webhookKontraktuBiblioteki(wiersz dane.WebhookBiblioteki) shared.LibraryWebhook {
	nasluch := shared.LibraryWebhook{
		Id: wiersz.Kod, Url: wiersz.Adres, Secret: wiersz.Sekret, Enabled: wiersz.Czynny,
		CreatedAt: chwilaBazy(wiersz.Utworzono),
	}
	for _, zdarzenie := range wiersz.Zdarzenia {
		nasluch.Events = append(nasluch.Events, zdarzenieWebhookaKontraktu(zdarzenie))
	}
	if wiersz.OstatnieZgloszenie != nil && *wiersz.OstatnieZgloszenie != "" {
		chwila := chwilaBazy(*wiersz.OstatnieZgloszenie)
		nasluch.LastDeliveryAt = &chwila
	}
	return nasluch
}

// zasiegUdostepnieniaBazy i zasiegUdostepnieniaKontraktu przekładają zakres
// udostępnienia (odwzorowanie: `udostepnienie_biblioteki.zasieg`).
func zasiegUdostepnieniaBazy(zasieg shared.LibraryShareScope) string {
	if zasieg == shared.LibraryShareScopeCollection {
		return "kolekcja"
	}
	return "plik"
}

func zasiegUdostepnieniaKontraktu(zasieg string) shared.LibraryShareScope {
	if zasieg == "kolekcja" {
		return shared.LibraryShareScopeCollection
	}
	return shared.LibraryShareScopeFile
}
