package uruchomienie

import (
	"context"

	"danacoconsole/server/internal/konfiguracja"
)

// Wedlug uruchamia tor interfejsu, tor wykonawczy albo oba naraz zależnie od
// roli procesu i wraca po zamknięciu kontekstu albo zakończeniu toru
// wiodącego; rola nierozpoznana pracuje jak rola łącząca oba tory.
func Wedlug(kontekst context.Context, rola konfiguracja.Rola, rdzen Rdzen, o Otoczenie) error {
	switch rola {
	case konfiguracja.RolaHub:
		o.zapisz("rola hub: tor interfejsu bez toru wykonawczego")
		return torInterfejsu(kontekst, rdzen)
	case konfiguracja.RolaAgent:
		o.zapisz("rola agent: tor wykonawczy bez serwera interfejsu")
		return TorWykonawczy(kontekst, rdzen, o)
	case konfiguracja.RolaWszystko:
		o.zapisz("rola all: tor interfejsu i tor wykonawczy")
		return obaTory(kontekst, rdzen, o)
	default:
		o.zapisz("rola %q nierozpoznana: praca jak dla roli %s", rola, konfiguracja.RolaWszystko)
		return obaTory(kontekst, rdzen, o)
	}
}

// torInterfejsu oddaje sterowanie warstwie nasłuchu rdzenia. Rdzeń niepodany
// nie przerywa procesu — tor czeka na zatrzymanie.
func torInterfejsu(kontekst context.Context, rdzen Rdzen) error {
	if rdzen == nil {
		<-kontekst.Done()
		return nil
	}
	return rdzen.Uruchom(kontekst)
}

// obaTory prowadzi tor wykonawczy w tle, a tor interfejsu na pierwszym planie.
// Wyczerpanie strumienia wejścia kończy wyłącznie tor wykonawczy: proces roli
// `all` uruchamiany bez podłączonego wejścia ma dalej obsługiwać interfejs.
func obaTory(kontekst context.Context, rdzen Rdzen, o Otoczenie) error {
	tlo, zatrzymajTlo := context.WithCancel(kontekst)
	defer zatrzymajTlo()

	go func() {
		if err := TorWykonawczy(tlo, rdzen, o); err != nil {
			o.zapisz("tor wykonawczy zakończony: %v", err)
			return
		}
		o.zapisz("tor wykonawczy zakończony: wyczerpane wejście")
	}()

	return torInterfejsu(kontekst, rdzen)
}
