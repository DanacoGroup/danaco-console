package uruchomienie

import (
	"context"

	"danacoconsole/server/internal/konfiguracja"
)

// Wedlug uruchamia tory właściwe dla roli procesu i wraca po zamknięciu
// kontekstu albo po zakończeniu toru wiodącego:
//
//	hub    — wyłącznie tor interfejsu: nasłuch transportu na porcie rdzenia;
//	         strumienie procesu pozostają nietknięte;
//	agent  — wyłącznie tor wykonawczy: żądania kontraktu ze strumienia wejścia,
//	         odpowiedzi na strumień wyjścia, żaden port nie jest zajmowany;
//	all    — oba tory naraz; torem wiodącym jest interfejs.
//
// Rola spoza katalogu ról nie zatrzymuje procesu: schodzi na zachowanie roli
// `all`, bo brak rozpoznanego ustawienia ma dawać pracę, nie odmowę.
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
