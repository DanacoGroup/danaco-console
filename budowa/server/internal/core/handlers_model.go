// Odpowiedzialność pliku: wpięcie rodziny `model.*`. Kontrakt niesie w niej
// dziś jedną komendę — `model.channel.set`.
//
// Port jest rozszerzeniem portu okien, nie drugim portem. Kanał modelu obsługuje
// okno komunikacji, a oknem w rdzeniu włada `Okna`. Rodzina `model.*`
// weszła do kontraktu osobno i osobno się wpina — dokładnie tak, jak `memory.*`
// wpina się osobno od `workspace.*`, jadąc na tej samej maszynerii.
//
// Zdarzenie jest tym samym, które rozgłasza `window.update`. Zmiana kanału jest
// zmianą okna, więc panel sterowania odświeża się z tej jednej subskrypcji
// `window.changed`; drugiego zdarzenia dla tego samego faktu kontrakt nie ma
// i rdzeń go sobie nie wymyśla. Żądanie wskazujące kartę sesji dotyka wielu okien
// naraz i rozgłasza tyle zmian, ile okien naprawdę zmieniło kanał.
package core

import (
	"context"

	"danacoconsole/shared"
)

// KanalModelu jest portem rodziny `model.*` — portem okien
// rozszerzonym o wybór kanału modelu.
type KanalModelu interface {
	Okna

	// UstawKanalModelu obsługuje `model.channel.set`. Drugi wynik niesie okna,
	// których kanał się zmienił — to one idą do rozgłoszenia.
	UstawKanalModelu(ctx context.Context, z shared.ModelChannelSetRequest) (
		shared.ModelChannelSetResponse, []shared.Window, error)
}

// zarejestrujKanalModelu wpina komendę rodziny `model.*`.
//
// Wpięcie w `Zloz` idzie przez asercję `p.Okna.(KanalModelu)` — celowo twardą,
// wzorem rodziny `memory.*`. Gdyby port okien przestał kiedyś nieść wybór
// kanału, rdzeń pada głośno przy montażu, zamiast po cichu zostawić komendę
// nieznaną. Cicha nieobecność uchwytu jest usterką, którą widać dopiero na
// uruchomionym produkcie.
func zarejestrujKanalModelu(r *Rejestr, m KanalModelu, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandModelChannelSet,
		obsluz(func(ctx context.Context, z shared.ModelChannelSetRequest) (
			shared.ModelChannelSetResponse, error) {

			odpowiedz, zmienione, err := m.UstawKanalModelu(ctx, z)
			if err != nil {
				return shared.ModelChannelSetResponse{}, err
			}
			for _, okno := range zmienione {
				e.okno(ctx, shared.ChangeKindUpdated, okno)
			}
			return odpowiedz, nil
		}))
}
