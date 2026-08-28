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
	r.Zarejestruj(shared.CommandConnectionHello,
		obsluz(func(ctx context.Context, z shared.ConnectionHelloRequest) (shared.ConnectionHelloResponse, error) {
			zdolnosci.zapamietaj(z.ClientId, z.Capabilities)
			odpowiedz := shared.ConnectionHelloResponse{
				ServerVersion:   WersjaRdzenia,
				ProtocolVersion: shared.ProtocolVersion,
				Commands:        r.Nazwy(),
			}
			// Nastawa nieznana zostaje milczeniem: pole puste znaczy brak
			// wiedzy, nie brak wymogu.
			if n != nil {
				if wymog, wskazana := n.WymogLogowania(ctx); wskazana {
					odpowiedz.LoginRequired = wskaznikPrawdy(wymog)
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
			odpowiedz.Authenticated = wskaznikPrawdy(zwiazPowitanie(ctx, u, wiez, z.Token))
			return odpowiedz, nil
		}))
}

// zwiazPowitanie rozpoznaje token i wiąże z nim połączenie, zwracając prawdę
// wyłącznie wtedy, gdy sesja istnieje, nie jest unieważniona i nie wygasła.
func zwiazPowitanie(ctx context.Context, u Uwierzytelnianie, wiez *wiezBramki, token *string) bool {
	if token == nil || *token == "" {
		wiez.Rozwiaz(polaczenieZKontekstu(ctx))
		return false
	}
	skrot, wazna, err := u.RozpoznajSesjeBramki(ctx, *token)
	if err != nil || !wazna {
		wiez.Rozwiaz(polaczenieZKontekstu(ctx))
		return false
	}
	wiez.Zwiaz(polaczenieZKontekstu(ctx), skrot)
	return true
}
