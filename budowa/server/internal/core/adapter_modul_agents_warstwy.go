// Tożsamość własna eksperta: warstwy jego promptu (agent.layer.set,
// agent.layer.remove) i jego wtyczki (agent.plugin.add, agent.plugin.remove)
// — cztery czynności portu WarstwyEksperta. Repozytorium warstw jest
// opcjonalne.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// warstwyTozsamosci wylicza warstwy nakładki z kontraktu. Ekspert bierze ten sam
// słownik co tożsamość modelu — drugiego zbioru nazw w systemie nie ma.
var warstwyTozsamosci = map[shared.IdentityLayer]struct{}{
	shared.IdentityLayerConstitution: {},
	shared.IdentityLayerProfile:      {},
	shared.IdentityLayerExpertise:    {},
}

// ZWarstwami wpina repozytorium tożsamości własnej eksperta. Bez niego cztery
// komendy warstw i wtyczek odmawiają, a reszta modułu Agents pracuje bez zmiany.
func (a *adapterAgentow) ZWarstwami(warstwy dane.RepozytoriumWarstwAgenta) *adapterAgentow {
	a.warstwy = warstwy
	return a
}

// UstawWarstwe zapisuje treść jednej warstwy promptu eksperta (`agent.layer.set`).
// Zapis powtórzony nadpisuje warstwę — warstwa jest stanem, nie zdarzeniem.
func (a *adapterAgentow) UstawWarstwe(ctx context.Context,
	z shared.AgentLayerSetRequest) (shared.AgentLayerSetResponse, error) {

	if a == nil || a.warstwy == nil {
		return shared.AgentLayerSetResponse{}, bladBrakuWarstwEksperta(shared.CommandAgentLayerSet)
	}
	if err := sprawdzWarstweEksperta(z.Layer); err != nil {
		return shared.AgentLayerSetResponse{}, err
	}
	// Trybu warstwy się nie wybiera — jest instrukcją dopisywaną do promptu,
	// nie zamiennikiem.
	const tryb = string(shared.IdentityModeDOLACZ)

	aktywna := true
	if z.Enabled != nil {
		aktywna = *z.Enabled
	}
	if err := a.warstwy.UstawWarstwe(ctx, z.AgentId, string(z.Layer), z.Content, tryb, aktywna); err != nil {
		return shared.AgentLayerSetResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	ekspert, err := a.EkspertPelny(ctx, z.AgentId)
	if err != nil {
		return shared.AgentLayerSetResponse{}, err
	}
	return shared.AgentLayerSetResponse{Agent: ekspert}, nil
}

// UsunWarstwe zdejmuje warstwę z promptu eksperta (`agent.layer.remove`). Brak
// warstwy nie jest odmową: prompt bez niej to ten sam stan, do którego komenda
// dąży.
func (a *adapterAgentow) UsunWarstwe(ctx context.Context,
	z shared.AgentLayerRemoveRequest) (shared.AgentLayerRemoveResponse, error) {

	if a == nil || a.warstwy == nil {
		return shared.AgentLayerRemoveResponse{}, bladBrakuWarstwEksperta(shared.CommandAgentLayerRemove)
	}
	if err := sprawdzWarstweEksperta(z.Layer); err != nil {
		return shared.AgentLayerRemoveResponse{}, err
	}
	if _, err := a.warstwy.UsunWarstwe(ctx, z.AgentId, string(z.Layer)); err != nil {
		return shared.AgentLayerRemoveResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	ekspert, err := a.EkspertPelny(ctx, z.AgentId)
	if err != nil {
		return shared.AgentLayerRemoveResponse{}, err
	}
	return shared.AgentLayerRemoveResponse{Agent: ekspert}, nil
}

// DodajWtyczke przypisuje ekspertowi wtyczkę (`agent.plugin.add`). Wtyczka jest
// katalogiem rozszerzeń powłoki, nie drogą do usługi — dlatego nie ma tu ani
// wskazania punktu dostępu, ani sprawdzenia mostu, które ma konektor.
func (a *adapterAgentow) DodajWtyczke(ctx context.Context,
	z shared.AgentPluginAddRequest) (shared.AgentPluginAddResponse, error) {

	if a == nil || a.warstwy == nil {
		return shared.AgentPluginAddResponse{}, bladBrakuWarstwEksperta(shared.CommandAgentPluginAdd)
	}
	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.AgentPluginAddResponse{}, bladZadaniaEksperta("wtyczka wymaga nazwy")
	}
	zapisana, err := a.warstwy.DodajWtyczke(ctx, z.AgentId, nazwa, z.Source, z.Version)
	if err != nil {
		return shared.AgentPluginAddResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	return shared.AgentPluginAddResponse{Plugin: wtyczkaKontraktu(zapisana)}, nil
}

// UsunWtyczke odłącza wtyczkę od eksperta (`agent.plugin.remove`). Wtyczka, której
// nie ma, kończy się tym samym stanem — wynik niesie eksperta, nie odmowę.
func (a *adapterAgentow) UsunWtyczke(ctx context.Context,
	z shared.AgentPluginRemoveRequest) (shared.AgentPluginRemoveResponse, error) {

	if a == nil || a.warstwy == nil {
		return shared.AgentPluginRemoveResponse{}, bladBrakuWarstwEksperta(shared.CommandAgentPluginRemove)
	}
	kod := strings.TrimSpace(z.PluginId)
	if kod == "" {
		return shared.AgentPluginRemoveResponse{}, bladZadaniaEksperta("odłączenie wymaga wskazania wtyczki")
	}
	if _, err := a.warstwy.UsunWtyczke(ctx, z.AgentId, kod); err != nil {
		return shared.AgentPluginRemoveResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	ekspert, err := a.EkspertPelny(ctx, z.AgentId)
	if err != nil {
		return shared.AgentPluginRemoveResponse{}, err
	}
	return shared.AgentPluginRemoveResponse{Agent: ekspert}, nil
}

// WykazWtyczek oddaje wtyczki eksperta wraz z definicją (agent.plugin.list).
// Osobna komenda, bo Agent.pluginIds niesie same identyfikatory, bez nazwy,
// źródła i wersji.
func (a *adapterAgentow) WykazWtyczek(ctx context.Context,
	z shared.AgentPluginListRequest) (shared.AgentPluginListResponse, error) {

	if a == nil || a.warstwy == nil {
		return shared.AgentPluginListResponse{}, bladBrakuWarstwEksperta(shared.CommandAgentPluginList)
	}
	kod := strings.TrimSpace(z.AgentId)
	if kod == "" {
		return shared.AgentPluginListResponse{}, bladZadaniaEksperta("wykaz wtyczek wymaga wskazania eksperta")
	}
	wiersze, err := a.warstwy.Wtyczki(ctx, kod)
	if err != nil {
		return shared.AgentPluginListResponse{}, bladWskazania(err, "ekspert", kod)
	}
	wykaz := make([]shared.AgentPlugin, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wykaz = append(wykaz, wtyczkaKontraktu(wiersz))
	}
	return shared.AgentPluginListResponse{Plugins: wykaz}, nil
}

// EkspertPelny oddaje eksperta wraz z warstwami i wtyczkami. Bez wpiętego
// repozytorium warstw oddaje to samo, co Pobierz — eksperta bez warstw i bez
// wtyczek.
func (a *adapterAgentow) EkspertPelny(ctx context.Context, idEksperta string) (shared.Agent, error) {
	ekspert, err := a.Pobierz(ctx, idEksperta)
	if err != nil {
		return shared.Agent{}, err
	}
	if a.warstwy == nil {
		return ekspert, nil
	}
	warstwy, err := a.warstwy.Warstwy(ctx, idEksperta)
	if err != nil {
		return shared.Agent{}, bladWskazania(err, "ekspert", idEksperta)
	}
	wtyczki, err := a.warstwy.Wtyczki(ctx, idEksperta)
	if err != nil {
		return shared.Agent{}, bladWskazania(err, "ekspert", idEksperta)
	}
	ekspert.Layers = warstwyKontraktu(warstwy)
	ekspert.PluginIds = kodyWtyczek(wtyczki)
	return ekspert, nil
}

// sprawdzWarstweEksperta pilnuje, żeby nazwa warstwy należała do kontraktu
// tożsamości modelu i platformy.
func sprawdzWarstweEksperta(warstwa shared.IdentityLayer) error {
	if _, jest := warstwyTozsamosci[warstwa]; jest {
		return nil
	}
	return bladZadaniaEksperta("warstwa " + string(warstwa) +
		" nie należy do kontraktu; dopuszczalne są constitution, profile i expertise")
}

// bladBrakuWarstwEksperta odmawia komendy warstw, nazywając komendę, powód
// i miejsce, w którym brak repozytorium się usuwa.
func bladBrakuWarstwEksperta(komenda shared.MessageType) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"serwer: komenda "+string(komenda)+" odmawia wykonania, ponieważ repozytorium warstw "+
			"eksperta nie jest wpięte do adaptera agentów; wpięcie zakłada się w pliku "+
			"server/internal/core/montaz_porty.go przy wywołaniu nowyAdapterAgentow(s.repozytoria.Agenci), "+
			"dopisując ogniwo .ZWarstwami(s.repozytoria.WarstwyAgenta)"))
}

// zapiszTozsamosc utrwala imię własne i favikon eksperta. Osobna droga zapisu,
// bo aktualizujAgenta nie zna kolumn imie_wlasne i favikon.
func (a *adapterAgentow) zapiszTozsamosc(
	ctx context.Context, kod string, imie, favikon *string,
) error {
	if imie == nil && favikon == nil {
		return nil
	}
	if a == nil || a.warstwy == nil {
		return nil
	}
	biezacy, err := a.repozytorium.PoKodzie(ctx, kod)
	if err != nil {
		return bladWskazania(err, "ekspert", kod)
	}
	noweImie := biezacy.ImieWlasne
	if imie != nil {
		noweImie = strings.TrimSpace(*imie)
	}
	nowyFavikon := biezacy.Favikon
	if favikon != nil {
		nowyFavikon = strings.TrimSpace(*favikon)
	}
	if err := a.warstwy.UstawTozsamosc(ctx, kod, noweImie, nowyFavikon); err != nil {
		return bladWskazania(err, "ekspert", kod)
	}
	return nil
}

// tozsamosciWykazu zbiera warstwy i wtyczki wszystkich ekspertów pod wykaz.
// Dwa zapytania na wywołanie, nie dwa na eksperta.
func (a *adapterAgentow) tozsamosciWykazu(
	ctx context.Context,
) (map[string][]dane.WarstwaAgenta, map[string][]dane.WtyczkaAgenta) {
	if a == nil || a.warstwy == nil {
		return nil, nil
	}
	warstwy, err := a.warstwy.WarstwyWszystkich(ctx)
	if err != nil {
		warstwy = nil
	}
	wtyczki, err := a.warstwy.WtyczkiWszystkich(ctx)
	if err != nil {
		wtyczki = nil
	}
	return warstwy, wtyczki
}

// zapiszTrybNakladki utrwala tryb nałożenia instrukcji eksperta (Agent.mode)
// tą samą drogą co tożsamość własna eksperta.
func (a *adapterAgentow) zapiszTrybNakladki(
	ctx context.Context, kod string, tryb *shared.IdentityMode,
) error {
	if tryb == nil {
		return nil
	}
	if a == nil || a.warstwy == nil {
		return nil
	}
	if err := a.warstwy.UstawTrybNakladki(ctx, kod, string(*tryb)); err != nil {
		return bladWskazania(err, "ekspert", kod)
	}
	return nil
}
