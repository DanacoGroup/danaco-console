package core

import (
	"context"

	"danacoconsole/shared"
)

// zarejestrujPolaczenie wpina powitanie klienta.
//
// Powitanie uzgadnia wersje i oddaje klientowi wykaz komend rzeczywiście
// obsługiwanych przez rdzeń — nie wykaz z kontraktu. Klient dowiaduje się
// z niego, co ta wersja rdzenia potrafi, zamiast zgadywać po odpowiedziach.
//
// Powitanie niesie token i jest jedyną drogą, którą token sesji bramki wchodzi
// do rdzenia: transportem jest jedno gniazdo, więc nie ma nagłówka na każdym
// żądaniu ani ciasteczka.
//
// Uzgodnienie nie jest bramą. Token niepasujący, wygasły albo pusty nie odrzuca
// powitania i nie zamyka połączenia — zmienia wyłącznie jedno pole odpowiedzi.
// Klient czyta `authenticated` i sam rozstrzyga, czy pokazać okno logowania.
//
// `gatewayConfigured` odróżnia dwie różne pustki. „Bramka istnieje, ale nie
// jesteś zalogowany" prowadzi do okna logowania; „bramki nie ma wcale" prowadzi
// do okna pierwszej rejestracji. Bez tego pola klient musiałby zgadywać po
// odmowie `auth.login`, czyli wyprowadzać stan z błędu.
//
// `loginRequired` jest odczytem przy nawiązaniu. Warstwa nasłuchu składa straż
// bramki raz na połączenie, więc chwilą, w której nastawa poziomu `aplikacja`
// ma znaczenie, jest właśnie powitanie. Rdzeń czyta ją tym samym rozstrzygaczem,
// którym idzie każde inne ustawienie platformy (`nastawy_aplikacji.go`), więc
// zmianę zapisaną komendą `config.set` widać od następnego połączenia, bez
// restartu rdzenia. Pole niczego nie odmawia: mówi klientowi, czy okno logowania
// ma się w ogóle pokazać.
//
// Deklaracja zdolności klienta (`capabilities`) jest tu zapamiętywana, bo
// powitanie to JEDYNA chwila, w której klient o sobie mówi. Korzysta z niej
// `launcher.hotkey.*`: skrót globalny przechwytuje powłoka programu okiennego,
// więc rdzeń musi wiedzieć, czy po drugiej stronie w ogóle stoi ktoś, kto to
// potrafi — inaczej obiecywałby skrót, który nikogo nie obudzi.
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
			// Nastawa nieznana zostaje milczeniem. Pole puste znaczy „rdzeń nie
			// wie”, co jest czymś innym niż „nie wymaga” — wpisanie tu fałszu
			// kazałoby klientowi ukryć okno logowania na nasłuchu, który wymogu
			// akurat pilnuje.
			if n != nil {
				if wymog, wskazana := n.WymogLogowania(ctx); wskazana {
					odpowiedz.LoginRequired = wskaznikPrawdy(wymog)
				}
			}
			// Rdzeń bez wpiętej bramki mówi o niej prawdę przez milczenie:
			// oba pola zostają puste, bo „nie wiadomo" to nie to samo co „nie".
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

// zwiazPowitanie rozpoznaje token i wiąże z nim połączenie. Zwraca prawdę
// wyłącznie wtedy, gdy sesja istnieje, nie jest unieważniona i nie wygasła.
//
// Błąd odczytu nie jest odmową ani potwierdzeniem — jest fałszem, tak samo jak
// token nieznany. Powitanie ma się udać zawsze; jedyną szkodą z niedostępnej
// bazy jest to, że Operator zobaczy okno logowania.
//
// Odpowiedź i więź mówią to samo: powitanie nierozpoznane zrywa poprzednią więź
// tego połączenia, zamiast zostawić ją nietkniętą. Więź nietknięta kazałaby
// rdzeniowi odpowiadać `authenticated: false`, a mimo to uważać gniazdo za
// związane z sesją sprzed powitania i wyłączać ją ze zmiany hasła — czyli
// oszczędzać sesję, której wołający nie przedstawił. Zerwanie nie jest bramką:
// niczego nie odrzuca i połączenia nie zamyka, tylko sprowadza stan do „nie
// wiadomo, kto woła".
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
