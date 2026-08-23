// Tożsamość własna eksperta: warstwy jego promptu (`agent.layer.set`,
// `agent.layer.remove`) i jego wtyczki (`agent.plugin.add`,
// `agent.plugin.remove`) — cztery czynności portu WarstwyEksperta.
//
// Repozytorium warstw jest opcjonalne. Adapter agentów składa się w
// `montaz_porty.go` bez ogniwa `.ZWarstwami(...)`, więc pole `warstwy` bywa
// `nil`; każda z czterech komend odmawia wtedy, nazywając komendę, powód
// i miejsce wpięcia, zamiast panikować albo udawać zapis.
//
// `agent.list` i `agent.create` nie przechodzą tędy — jadą dalej
// `ekspertKontraktu` i przy `warstwy == nil` oddają eksperta bez warstw
// i bez wtyczek. Nazwy warstw sprawdzane są wobec wartości kontraktu
// (`shared.IdentityLayer`, `shared.IdentityMode`) przed zapisem, żeby powodem
// odmowy było zdanie po polsku, a nie naruszony warunek CHECK.
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
	// Trybu warstwy eksperta się nie wybiera. Warstwa jest instrukcją DOPISYWANĄ
	// do promptu systemowego, nie jego zamiennikiem, więc kontrakt nie ma dla
	// niej pola `mode` — tryb jest tu stałą, nie parametrem. Tożsamość osi
	// (`identity.document.set`) tryb wybiera i domyślnie zastępuje, bo jest
	// konfiguracją platformy, a nie warstwą nałożoną na pojedyncze wywołanie.
	//
	// Wartość musi pochodzić z kontraktu, nie z silnika nakładki: kolumna
	// `agent_warstwa.tryb` ma warunek CHECK na `ZASTAP` albo `DOLACZ`, a stała
	// silnika `injection.TrybDopisz` jest napisem `"dopisz"`. Przekład między
	// jednym a drugim robi `trybSilnika`.
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

// WykazWtyczek oddaje wtyczki eksperta wraz z definicją (`agent.plugin.list`).
// Osobna komenda jest potrzebna, bo `Agent.pluginIds` niesie same
// identyfikatory — bez nazwy nadanej wtyczce, bez źródła i bez wersji.
//
// Ekspert bez wtyczek oddaje tablicę pustą: brak wtyczek jest poprawnym stanem,
// a nie awarią odczytu. Odmowa idzie wyłącznie wtedy, gdy eksperta o wskazanym
// kodzie nie ma w katalogu.
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
// repozytorium warstw oddaje dokładnie to, co `Pobierz` — eksperta bez warstw
// i bez wtyczek. Odczyt nie odmawia z powodu niewpiętego katalogu: czytać nie ma
// czego, a to inny przypadek niż zapis, który nigdzie nie trafia.
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

// sprawdzWarstweEksperta pilnuje, żeby nazwa warstwy należała do kontraktu.
func sprawdzWarstweEksperta(warstwa shared.IdentityLayer) error {
	if _, jest := warstwyTozsamosci[warstwa]; jest {
		return nil
	}
	return bladZadaniaEksperta("warstwa " + string(warstwa) +
		" nie należy do kontraktu; dopuszczalne są constitution, profile i expertise")
}

// bladBrakuWarstwEksperta odmawia komendy warstw, nazywając komendę, powód —
// niewpięte repozytorium — i miejsce, w którym brak się usuwa. Milczenie
// udawałoby zapis, a panika przewracałaby proces z powodu brakującej linii
// montażu.
func bladBrakuWarstwEksperta(komenda shared.MessageType) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"rdzeń: komenda "+string(komenda)+" odmawia wykonania, ponieważ repozytorium warstw "+
			"eksperta nie jest wpięte do adaptera agentów; wpięcie zakłada się w pliku "+
			"server/internal/core/montaz_porty.go przy wywołaniu nowyAdapterAgentow(s.repozytoria.Agenci), "+
			"dopisując ogniwo .ZWarstwami(s.repozytoria.WarstwyAgenta)"))
}

// zapiszTozsamosc utrwala imię własne i favikon eksperta.
//
// Osobna droga zapisu, bo zapytanie `aktualizujAgenta` nie zna kolumn
// `imie_wlasne` i `favikon`; osobna czynność repozytorium dotyka wyłącznie tych
// dwóch kolumn, zamiast zmieniać drogę, którą jadą wszystkie pozostałe pola.
//
// Pominięte pole nie jest polem pustym: `nil` zostawia wartość zastaną, pusty
// napis czyści ją i zostaje zapisany. Zlanie obu w jedno kasowałoby imię przy
// każdej zmianie samego opisu.
//
// Brak wpiętego repozytorium nie jest tu odmową — ekspert bez imienia własnego
// jest ekspertem, bo tożsamość niesie pole `nazwa`.
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
//
// Dwa zapytania na wywołanie, nie dwa na eksperta: odczyt po jednym dałby przy
// stu ekspertach dwieście zapytań na jedno otwarcie biblioteki. Wzorzec ten sam
// co `dolaczPowiazania` w warstwie danych.
//
// Błąd odczytu nie wywraca wykazu — ekspert, którego warstw nie udało się
// odczytać, wchodzi do wykazu bez warstw, tak samo jak ekspert, który ich nie
// ma. Wykaz ekspertów ma się pokazać także wtedy, gdy tożsamość jest chwilowo
// nieczytelna.
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

// zapiszTrybNakladki utrwala tryb nałożenia instrukcji eksperta (`Agent.mode`).
//
// Tą samą drogą co tożsamość własna i z tego samego powodu: `aktualizujAgenta`
// nie zna kolumny `tryb_nakladki`.
//
// Pominięte pole zostawia wartość zastaną — `nil` nie znaczy powrotu do
// domyślnego trybu; raz oznaczone odstępstwo nie ma prawa zniknąć przy zmianie
// samego opisu eksperta. Wartość spoza katalogu odrzuca warunek CHECK na
// kolumnie; drugiej listy dopuszczonych trybów tu nie ma.
//
// Brak wpiętego repozytorium nie jest tu odmową: ekspert bez zapisanego trybu
// dopisuje się do promptu globalnego.
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
