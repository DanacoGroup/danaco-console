// Plik deklaruje zdolności adaptera dostawcy dla komendy config.capabilities.get i pola capabilities konfiguracji obowiązującej. Żaden adapter jeszcze nie wystawia wykazu odwzorowań, więc rdzeń oddaje deklarację pustą ze wskazanym powodem.
package core

import (
	"context"

	"danacoconsole/shared"
)

// idAdapteraCli nazywa adapter powierzchni CLI w deklaracji zdolności
// (AdapterCapabilities.AdapterId). Mieszka tu, bo to jedyne miejsce, które go
// czyta.
const idAdapteraCli = "cli.claude-code"

// ZdolnosciAdaptera zwraca deklarację zdolności adaptera dostawcy dla transportu wskazanego w żądaniu albo domyślnego kanału CLI.
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
		// Powód pustki jedzie do klienta, żeby wiedział, dlaczego wykaz pól jest pusty, a nie niezbadany.
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
