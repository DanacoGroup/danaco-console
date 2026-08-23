// Odpowiedzialność pliku: model bazowy eksperta i parametry jego wywołania
// (okno Model Configuration, komenda `agent.model.set`).
//
// Kanał wskazany przez okno jest sprawdzany w rejestrze kanałów rdzenia — tym
// samym, który obsługuje `channel.list` i okna rozmowy. Adapter nie zna żadnego
// dostawcy z nazwy; zna wyłącznie kod wiersza rejestru.
//
// Sprawdzenie jest możliwe, nie obowiązkowe. Rejestr pusty znaczy „rejestr
// jeszcze nie wstał”, a nie „kanał nie istnieje”, więc odmowa zapisu w takim
// stanie zablokowałaby konfigurację eksperta z powodu leżącego poza nią. Odmowa
// zapada dopiero wtedy, gdy rejestr coś zna, a wskazania w nim nie ma.
package core

import (
	"context"
	"strings"

	"danacoconsole/shared"
)

// UstawModel ustala kanał, model bazowy, transport i parametry wywołania.
func (a *adapterAgentow) UstawModel(ctx context.Context,
	z shared.AgentModelSetRequest) (shared.AgentModelSetResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AgentModelSetResponse{}, bladBrakuKatalogu("ekspertów")
	}
	kod := strings.TrimSpace(z.ChannelId)
	if kod == "" {
		return shared.AgentModelSetResponse{}, bladZadaniaEksperta("konfiguracja modelu wymaga kanału")
	}
	if err := sprawdzTransport(z.Transport); err != nil {
		return shared.AgentModelSetResponse{}, err
	}
	parametry, err := sprawdzParametry(z.Parameters, "parametry wywołania")
	if err != nil {
		return shared.AgentModelSetResponse{}, err
	}
	kanal, err := a.kanalWskazany(&kod)
	if err != nil {
		return shared.AgentModelSetResponse{}, err
	}
	biezacy, err := a.repozytorium.PoKodzie(ctx, z.AgentId)
	if err != nil {
		return shared.AgentModelSetResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	biezacy.KanalKod = kanal
	biezacy.Model = a.modelWskazany(z.Model, kanal)
	biezacy.Transport = transportWiersza(z.Transport)
	biezacy.ParametryJSON = parametry
	zapisany, err := a.repozytorium.Aktualizuj(ctx, biezacy)
	if err != nil {
		return shared.AgentModelSetResponse{}, err
	}
	return shared.AgentModelSetResponse{Agent: ekspertKontraktu(zapisany)}, nil
}

// kanalWskazany sprawdza kod kanału w rejestrze rdzenia. Wskazanie puste znaczy
// „ekspert bez kanału bazowego” i przechodzi bez sprawdzenia.
func (a *adapterAgentow) kanalWskazany(kod *string) (*string, error) {
	wskazanie := strings.TrimSpace(wartoscTekstu(kod))
	if wskazanie == "" {
		return nil, nil
	}
	if a.kanaly == nil || len(a.kanaly.Wykaz()) == 0 {
		return &wskazanie, nil
	}
	if _, jest := a.kanaly.Definicja(wskazanie); !jest {
		return nil, bladZadaniaEksperta("kanał " + wskazanie + " nie istnieje w rejestrze kanałów modelu")
	}
	return &wskazanie, nil
}

// modelWskazany dobiera identyfikator modelu. Wskazanie własne ma pierwszeństwo;
// w jego braku model bierze się z wiersza kanału, żeby okno nie musiało
// przepisywać wartości, którą rejestr już zna.
func (a *adapterAgentow) modelWskazany(model, kanal *string) *string {
	if wskazany := strings.TrimSpace(wartoscTekstu(model)); wskazany != "" {
		return &wskazany
	}
	if a.kanaly == nil || kanal == nil {
		return nil
	}
	definicja, jest := a.kanaly.Definicja(*kanal)
	if !jest {
		return nil
	}
	return wskaznikTekstu(strings.TrimSpace(definicja.Model))
}

// transportWiersza przenosi drogę wywołania do kolumny. Brak wskazania znaczy
// „transport kanału”, więc kolumna zostaje pusta zamiast zgadywać wartość.
func transportWiersza(transport *shared.ProviderTransport) *string {
	if transport == nil {
		return nil
	}
	wartosc := string(*transport)
	return &wartosc
}
