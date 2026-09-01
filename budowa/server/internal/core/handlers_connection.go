package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujPolaczenie wpina powitanie klienta, uzgadniające wersję
// protokołu, wiążące token sesji bramki i zwracające zdolności obsługiwane
// przez rdzeń.
func zarejestrujPolaczenie(r *Rejestr, u Uwierzytelnianie, wiez *wiezBramki, n NastawyAplikacji,
	zdolnosci *zdolnosciKlientow) {
	if r == nil {
		return
	}
	// Konto sesji bramki oddaje ten sam adapter, który sesje wydaje; port
	// Uwierzytelnianie go nie wymienia, bo rozpoznanie konta służy adresowaniu
	// rozgłoszeń, nie powitaniu.
	konta, _ := u.(RozpoznanieKontaSesji)
	r.Zarejestruj(shared.CommandConnectionHello,
		obsluz(func(ctx context.Context, z shared.ConnectionHelloRequest) (shared.ConnectionHelloResponse, error) {
			zdolnosci.zapamietaj(z.ClientId, z.Capabilities)
			odpowiedz := shared.ConnectionHelloResponse{
				ServerVersion:   WersjaRdzenia,
				ProtocolVersion: shared.ProtocolVersion,
				Commands:        r.Nazwy(),
			}
			// Pole niesie wyłącznie wymóg włączony nastawą. Nastawa zdejmująca
			// nie zdejmuje bramki (o tym rozstrzyga adres nasłuchu w transporcie),
			// więc pole zostaje puste — brak wiedzy, nie brak wymogu.
			if n != nil {
				if wymog, wskazana := n.WymogLogowania(ctx); wskazana && wymog {
					odpowiedz.LoginRequired = wskaznikPrawdy(true)
				}
			}
			// Rdzeń bez wpiętej bramki mówi o tym milczeniem — oba pola
			// zostają puste.
			if u == nil {
				return odpowiedz, nil
			}
			if zalozona, err := u.BramkaZalozona(ctx); err == nil {
				odpowiedz.GatewayConfigured = &zalozona
			}
			odpowiedz.Authenticated = wskaznikPrawdy(zwiazPowitanie(ctx, u, wiez, konta, z.Token))
			return odpowiedz, nil
		}))
}

// zwiazPowitanie rozpoznaje token i wiąże z nim połączenie, zwracając prawdę
// wyłącznie wtedy, gdy sesja istnieje, nie jest unieważniona i nie wygasła.
// Konto gniazda idzie za więzią w obie strony: powitanie bez ważnego tokenu
// odsyła gniazdo na konto domyślne, żeby połączenie rozwiązane przestało
// dostawać rozgłoszenia konta, na którym stało wcześniej.
func zwiazPowitanie(ctx context.Context, u Uwierzytelnianie, wiez *wiezBramki,
	konta RozpoznanieKontaSesji, token *string) bool {

	if token == nil || *token == "" {
		wiez.Rozwiaz(polaczenieZKontekstu(ctx))
		przypiszKontoGniazda(ctx, konta, "")
		return false
	}
	skrot, wazna, err := u.RozpoznajSesjeBramki(ctx, *token)
	if err != nil || !wazna {
		wiez.Rozwiaz(polaczenieZKontekstu(ctx))
		przypiszKontoGniazda(ctx, konta, "")
		return false
	}
	wiez.Zwiaz(polaczenieZKontekstu(ctx), skrot)
	przypiszKontoGniazda(ctx, konta, skrot)
	return true
}
