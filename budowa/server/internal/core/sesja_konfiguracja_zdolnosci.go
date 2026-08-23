// Odpowiedzialność pliku: deklaracja zdolności adaptera dostawcy dla komendy
// `config.capabilities.get` oraz dla pola `capabilities` konfiguracji
// obowiązującej.
//
// Stan faktyczny, nie życzeniowy. Deklaracja zdolności może pochodzić wyłącznie
// od adaptera, który tłumaczy model konfiguracji na powierzchnię dostawcy.
// Żaden adapter takiego wykazu jeszcze nie wystawia — nie ma ani funkcji
// wejściowej, ani tablicy odwzorowań pole → powierzchnia. Rdzeń oddaje więc
// deklarację pustą:
//
//   - wykaz obszarów i wykaz pól są puste — rdzeń nie zna ani jednego
//     rozstrzygnięcia adaptera, więc żadnego nie zgłasza;
//   - probedAt pozostaje zerowe — nic nie zostało sprawdzone, a czas
//     sprawdzenia, którego nie było, byłby zapisem czynności niewykonanej;
//   - channelId i accountId pozostają puste — nie sprawdzono żadnego kanału
//     ani konta.
//
// Wypełnienie tych wykazów zgadywanką byłoby gorsze od pustki: okno
// konfiguracji pokazałoby Operatorowi, że pole jedzie do dostawcy albo że nie
// jedzie, na podstawie niczego. Pusta prawda jest tu jedyną dopuszczalną
// odpowiedzią, a brak deklaracji niczego nie wstrzymuje — zapis
// konfiguracji idzie dalej.
//
// Identyfikator adaptera jest stałą tego pakietu, bo powierzchnia CLI nie ma
// w drzewie osobnego pakietu — przekład konfiguracji sesji na wejście procesu
// jedzie warstwą injection na drodze tury (adapter_rozmowa_*).
// Transport spoza kanału CLI nie ma w drzewie odpowiednika, więc jego
// identyfikator pozostaje pusty.
package core

import (
	"context"

	"danacoconsole/shared"
)

// idAdapteraCli nazywa adapter powierzchni CLI w deklaracji zdolności
// (AdapterCapabilities.AdapterId). Mieszka tu, bo to jedyne miejsce, które go
// czyta.
const idAdapteraCli = "cli.claude-code"

// ZdolnosciAdaptera zwraca deklarację zdolności adaptera dostawcy.
func (a *adapterUstawienOsi) ZdolnosciAdaptera(_ context.Context,
	z shared.ConfigCapabilitiesGetRequest) (shared.ConfigCapabilitiesGetResponse, error) {

	return shared.ConfigCapabilitiesGetResponse{Capabilities: deklaracjaZdolnosci(z.Transport)}, nil
}

// deklaracjaZdolnosci składa deklarację dla wskazanego transportu. Transport
// pominięty znaczy kanał CLI — jedyny transport, dla którego adapter w ogóle
// istnieje.
func deklaracjaZdolnosci(transport *shared.ProviderTransport) shared.AdapterCapabilities {
	wybrany := shared.ProviderTransport(shared.ProviderTransportCli)
	if transport != nil && *transport != "" {
		wybrany = *transport
	}
	return shared.AdapterCapabilities{
		AdapterId: identyfikatorAdaptera(wybrany),
		Transport: wybrany,
		Areas:     []shared.SessionConfigAreaCapability{},
		Fields:    []shared.SessionConfigFieldCapability{},
		// Powód pustki jedzie do klienta, nie zostaje w komentarzu. Wykaz pusty
		// bez powodu jest nieodróżnialny od wykazu, którego nikt nie policzył —
		// a Operator ma prawo wiedzieć, dlaczego okno konfiguracji nie mówi mu,
		// które pola dostawca zignoruje. Pola `unsupported` nie wpisujemy, bo
		// znaczą „dostawca nie ma odpowiednika", a nie „nie zbadano".
		Reason: wskaznikTekstu("adapter dostawcy nie wystawia jeszcze wykazu odwzorowań pole → powierzchnia; zdolności nie zostały zbadane"),
	}
}

// identyfikatorAdaptera nazywa adapter obsługujący transport. Transport bez
// adaptera w drzewie nie dostaje nazwy zastępczej — pusta nazwa mówi wprost,
// że adaptera nie ma.
func identyfikatorAdaptera(transport shared.ProviderTransport) string {
	if transport == shared.ProviderTransportCli {
		return idAdapteraCli
	}
	return ""
}
