// Odpowiedzialność pliku: sprawca zdarzenia — rozstrzygnięcie, czyja ręka
// wykonała czynność, i milczenie tam, gdzie ręki nie widać. Operator ma widzieć
// nie tylko skutek, ale i rękę, żeby odróżnić okno założone przez siebie od okna
// założonego za niego przez asystenta.
//
// Rękę rozpoznajemy po faktach gniazda, nie po treści żądania i nie po
// przedrostku napisu:
//
//	operator  — gniazdo, które przedstawiło się identyfikatorem klienta
//	            w powitaniu i nie jest serwerem narzędzi modelu. Tak wygląda
//	            okno interfejsu Operatora i nic innego tak nie wygląda.
//	assistant — serwer narzędzi okna o roli klawiatury. Rola przychodzi z wpisu
//	            MCP ułożonego przez rdzeń dla okna modułu Assistant
//	            (`adapter_modul_asystent_sterowanie.go` → `--zasieg klawiatura`),
//	            więc jest faktem rdzenia, nie deklaracją modelu.
//	model     — serwer narzędzi każdego innego okna, czyli model roboczy.
//	core      — czynność powołana przez sam rdzeń: przemiatanie, harmonogram,
//	            odtworzenie stanu. Tego nie da się wyprowadzić z gniazda, bo
//	            gniazda tam nie ma — więc te miejsca znaczą się same
//	            (`zSprawcaRdzenia`), zamiast być domyślane z pustki.
//
// Brak gniazda nie znaczy `core`. Poza pracą własną rdzenia bez gniazda woła się
// też z próby, z sondy stdio i z biegu wewnętrznego adaptera — nazwanie tego
// wszystkiego rdzeniem byłoby zgadywaniem. Kontrakt mówi o polu `actor` wprost:
// „brak znaczy, że rdzeń nie potrafił tego rozstrzygnąć", więc brak jest tu
// odpowiedzią, a nie luką. Gniazdo bez identyfikatora klienta też zostaje bez
// sprawcy: połączenie, które jeszcze się nie przywitało, może być czymkolwiek.
//
// Sprawca jest opisem i niczego nie rozstrzyga o tym, czy wolno. W tym pliku nie
// ma ani jednej odmowy, a wynik nie wchodzi do żadnego warunku poza wypełnieniem
// pola zdarzenia. Strażą rdzenia zostaje bramka wejścia (`transport/bramka.go`).
package core

import (
	"context"

	"danacoconsole/server/internal/narzedzia"
	"danacoconsole/server/internal/transport"
	"danacoconsole/shared"
)

// kluczSprawcyRdzenia znaczy kontekst pracy własnej rdzenia.
const kluczSprawcyRdzenia kluczKontekstu = "danaco:sprawca-rdzenia"

// zSprawcaRdzenia znakuje kontekst biegu powołanego przez rdzeń, nie przez
// wołającego: przemiatania, harmonogramu, odtworzenia stanu po restarcie.
//
// Znak stawia się ręcznie, w miejscu, które wie, czym jest. Wywiedzenie `core`
// z samego braku gniazda byłoby nieprawdziwe — brak gniazda ma więcej powodów
// niż jeden (zob. nagłówek pliku).
func zSprawcaRdzenia(ctx context.Context) context.Context {
	if ctx == nil {
		return ctx
	}
	return context.WithValue(ctx, kluczSprawcyRdzenia, true)
}

// rdzenJestSprawca odczytuje ten znak.
func rdzenJestSprawca(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	znak, jest := ctx.Value(kluczSprawcyRdzenia).(bool)
	return jest && znak
}

// sprawca oddaje parę pól kontraktu wspólną sześciu zdarzeniom: rodzaj sprawcy
// i identyfikator jego klienta.
//
// Oba pola opisują jedną rzecz i mają milczeć razem. Rodzaj bez klienta jest
// zdaniem pełnym („to był model"); klient bez rodzaju byłby napisem bez
// znaczenia.
func sprawca(ctx context.Context) (*shared.ActorKind, *string) {
	if rdzenJestSprawca(ctx) {
		// Rdzeń nie jest klientem i nie ma identyfikatora klienta — pole zostaje
		// puste, bo pustka znaczy tu „nie dotyczy", a nie „nie wiadomo".
		return rodzajSprawcy(shared.ActorKindCore), nil
	}
	tozsamosc := tozsamoscZKontekstu(ctx)
	if tozsamosc.Narzedzia() {
		return rodzajSprawcy(rodzajNarzedzi(tozsamosc)), klientSprawcy(tozsamosc.IdKlienta)
	}
	if tozsamosc.IdKlienta != "" {
		return rodzajSprawcy(shared.ActorKindOperator), klientSprawcy(tozsamosc.IdKlienta)
	}
	return nil, nil
}

// rodzajNarzedzi rozdziela dwie ręce pracujące tą samą drogą narzędzi: klawiaturę
// Operatora i model roboczy.
//
// Rozdziela je rola okna, którą rdzeń nadał wpisowi MCP, zanim proces modelu
// wystartował — a nie cokolwiek, o co model mógłby poprosić w trakcie. Zasięg
// nieznany daje `model`: to jest prawda węższa, ale prawda, i tak samo czyta go
// sam serwer narzędzi (`RozpoznajZasieg`).
func rodzajNarzedzi(t transport.Tozsamosc) string {
	if narzedzia.RozpoznajZasieg(t.Zasieg) == narzedzia.ZasiegKlawiatury {
		return shared.ActorKindAssistant
	}
	return shared.ActorKindModel
}

// rodzajSprawcy przenosi wartość wyliczenia do pola opcjonalnego kontraktu.
func rodzajSprawcy(rodzaj string) *shared.ActorKind {
	wartosc := shared.ActorKind(rodzaj)
	return &wartosc
}

// klientSprawcy przenosi identyfikator klienta do pola opcjonalnego. Napis pusty
// nie zakłada pola — „nie wiadomo" ma wyglądać na brak, nie na pustą wartość.
func klientSprawcy(id string) *string {
	if id == "" {
		return nil
	}
	return &id
}
