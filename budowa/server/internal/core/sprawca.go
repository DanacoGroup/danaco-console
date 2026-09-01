// Plik rozstrzyga sprawcę zdarzenia po faktach gniazda: operator, assistant, model albo core, milcząc, gdy ręki nie widać. Sprawca jest opisem i niczego nie rozstrzyga o uprawnieniach.
package core

import (
	"context"

	"danacoconsole/server/internal/narzedzia"
	"danacoconsole/server/internal/transport"
	"danacoconsole/shared"
)

// kluczSprawcyRdzenia znaczy kontekst pracy własnej rdzenia, wprowadzony do kontekstu funkcją zSprawcaRdzenia.
const kluczSprawcyRdzenia kluczKontekstu = "danaco:sprawca-serwera"

// zSprawcaRdzenia znakuje kontekst biegu powołanego przez rdzeń, nie przez wołającego: przemiatania, harmonogramu, odtworzenia stanu po restarcie. Znak stawia się ręcznie, w miejscu, które wie, czym jest.
func zSprawcaRdzenia(ctx context.Context) context.Context {
	if ctx == nil {
		return ctx
	}
	return context.WithValue(ctx, kluczSprawcyRdzenia, true)
}

// rdzenJestSprawca odczytuje znak kontekstu ustawiony funkcją zSprawcaRdzenia, mówiąc, czy sprawcą jest sam rdzeń.
func rdzenJestSprawca(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	znak, jest := ctx.Value(kluczSprawcyRdzenia).(bool)
	return jest && znak
}

// sprawca oddaje parę pól kontraktu wspólną sześciu zdarzeniom: rodzaj sprawcy i identyfikator jego klienta. Oba pola mają milczeć razem — rodzaj bez klienta jest zdaniem pełnym, klient bez rodzaju byłby napisem bez znaczenia.
func sprawca(ctx context.Context) (*shared.ActorKind, *string) {
	if rdzenJestSprawca(ctx) {
		// Rdzeń nie jest klientem i nie ma identyfikatora klienta — pole zostaje puste, znacząc nie dotyczy.
		return rodzajSprawcy(shared.ActorKindCore), nil
	}
	tozsamosc := tozsamoscZKontekstu(ctx)
	// Rodzaj narzędzi liczy się po sprawdzeniu poświadczenia przez transport; napis z zapytania bez poświadczenia nie wchodzi do tożsamości.
	if tozsamosc.Narzedzia() {
		return rodzajSprawcy(rodzajNarzedzi(tozsamosc)), klientSprawcy(tozsamosc.IdKlienta)
	}
	// Identyfikator klienta bez sprawdzonego poświadczenia pochodzi z powitania okna, nie z nawiązania.
	if tozsamosc.IdKlienta != "" {
		return rodzajSprawcy(shared.ActorKindOperator), klientSprawcy(tozsamosc.IdKlienta)
	}
	return nil, nil
}

// rodzajNarzedzi rozdziela dwie ręce pracujące tą samą drogą narzędzi: klawiaturę Operatora i model roboczy, według roli okna nadanej wpisowi MCP przez rdzeń wraz z poświadczeniem, zanim proces modelu wystartował. Zasięg nieznany daje model.
func rodzajNarzedzi(t transport.Tozsamosc) string {
	if narzedzia.RozpoznajZasieg(t.Zasieg) == narzedzia.ZasiegKlawiatury {
		return shared.ActorKindAssistant
	}
	return shared.ActorKindModel
}

// rodzajSprawcy przenosi wartość wyliczenia rodzaju sprawcy do pola opcjonalnego kontraktu typu ActorKind.
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
