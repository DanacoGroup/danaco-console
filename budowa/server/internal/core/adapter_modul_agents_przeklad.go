// Odpowiedzialność pliku: przekład między wierszem biblioteki ekspertów
// a strukturami kontraktu oraz sprawdziany kształtu żądań obszaru `agent.*`.
// Identyfikatorem kontraktu jest kod wiersza, nie numer wewnętrzny bazy.
package core

import (
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// przedrostekEksperta znakuje kod eksperta nadany przez rdzeń przy założeniu
// nowego wiersza biblioteki.
const przedrostekEksperta = "ag-"

// przedrostekKonektora znakuje kod konektora nadany przez rdzeń przy założeniu
// nowego wiersza konektora.
const przedrostekKonektora = "ak-"

// transportyKanalu wylicza drogi wywołania dostawcy z kontraktu. Wartość spoza
// zbioru jest pomyłką klienta, nie nowym transportem.
var transportyKanalu = map[shared.ProviderTransport]struct{}{
	shared.ProviderTransportCli:   {},
	shared.ProviderTransportApi:   {},
	shared.ProviderTransportSdk:   {},
	shared.ProviderTransportSsh:   {},
	shared.ProviderTransportLocal: {},
}

// rodzajeKonektora wylicza rodzaje konektora z kontraktu, dopuszczalne przy
// sprawdzianie kształtu żądania.
var rodzajeKonektora = map[shared.AgentConnectorKind]struct{}{
	shared.AgentConnectorKindMcp:    {},
	shared.AgentConnectorKindPlugin: {},
	shared.AgentConnectorKindApi:    {},
}

// poziomyPamieci wylicza cztery poziomy pamięci z kontraktu. Wyłączenie pamięci
// jest zbiorem pustym, nie piątą wartością — inaczej dałoby się zgłosić żądanie
// sprzeczne wobec pustego zbioru.
var poziomyPamieci = map[shared.MemoryLevel]struct{}{
	shared.MemoryLevelGlobal:      {},
	shared.MemoryLevelProject:     {},
	shared.MemoryLevelSession:     {},
	shared.MemoryLevelEnvironment: {},
}

// poziomyPamieciWyjsciowe to komplet czterech poziomów — wartość, z którą staje
// ekspert świeżo założony, bo stanem wyjściowym platformy jest pełny dostęp.
var poziomyPamieciWyjsciowe = []string{
	shared.MemoryLevelGlobal,
	shared.MemoryLevelProject,
	shared.MemoryLevelSession,
	shared.MemoryLevelEnvironment,
}

// widocznosciEksperta wylicza zasięgi widoczności z kontraktu, dopuszczalne
// przy sprawdzianie kształtu żądania.
var widocznosciEksperta = map[shared.AgentVisibility]struct{}{
	shared.AgentVisibilityGlobal:  {},
	shared.AgentVisibilityProject: {},
}

// grupyUprawnien wylicza cztery grupy zakresu uprawnień z kontraktu, poza
// grupą modułów opisaną osobno.
var grupyUprawnien = map[shared.AgentPermissionGroup]struct{}{
	shared.AgentPermissionGroupFiles:        {},
	shared.AgentPermissionGroupNetwork:      {},
	shared.AgentPermissionGroupProcesses:    {},
	shared.AgentPermissionGroupIntegrations: {},
	// Piąta grupa zakresu: moduły i zasoby platformy, poza tym wyliczeniem.
	shared.AgentPermissionGroupModules: {},
}

// ekspertKontraktu przekłada wiersz biblioteki ekspertów na strukturę Agent
// kontraktu, wraz z polami zawsze obecnymi.
func ekspertKontraktu(a dane.Agent) shared.Agent {
	ekspert := shared.Agent{
		Id:           a.Kod,
		Name:         a.Nazwa,
		Description:  wskaznikTekstu(a.Opis),
		SystemPrompt: wskaznikTekstu(a.InstrukcjeSystemowe),
		ChannelId:    a.KanalKod,
		Model:        a.Model,
		SkillIds:     a.Umiejetnosci,
		ConnectorIds: a.Konektory,
		Permissions:  uprawnieniaKontraktu(a.Uprawnienia),
		DisplayName:  wskaznikTekstu(a.ImieWlasne),
		Favicon:      wskaznikTekstu(a.Favikon),
		// Tryb oddawany zawsze, także gdy jest domyślny, by nie zostawiać domysłu.
		Mode: trybKontraktu(a.TrybNakladki),
		// Poziomy pamięci oddawane zawsze, także puste — pusty wycinek jest
		// odpowiedzią pamięć wyłączona.
		MemoryLevels: poziomyKontraktu(a.PoziomyPamieci),
		Visibility:   widocznoscKontraktu(a.Widocznosc),
		Enabled:      a.Aktywny,
		// Moduły zastosowania oddawane zawsze, także puste: lista pusta znaczy
		// brak ograniczenia.
		ModuleCodes: kodyModulowKontraktu(a.ModulyZastosowania),
		// Granica podagentów oddawana zawsze: zero znaczy Subagent Network wyłączony.
		SubagentLimit: a.LimitPodagentow,
	}
	wersja := a.Wersja
	ekspert.Version = &wersja
	if chwila, ok := chwilaZapisu(a.Utworzono); ok {
		ekspert.CreatedAt = chwila
	}
	if chwila, ok := chwilaZapisu(a.Zaktualizowano); ok {
		ekspert.UpdatedAt = chwila
	}
	return ekspert
}

// kodyModulowKontraktu przekłada moduły zastosowania na wycinek kodów.
// Wycinek pusty wraca jako tablica pusta, nie jako nil, tak jak inne pola wynikowe.
func kodyModulowKontraktu(kody []string) []string {
	wynik := make([]string, 0, len(kody))
	wynik = append(wynik, kody...)
	return wynik
}

// poziomyKontraktu przekłada poziomy pamięci z bazy na wycinek kontraktu.
// Wycinek pusty wraca jako tablica pusta, nie jako nil, tak jak inne pola wynikowe.
func poziomyKontraktu(poziomy []string) []shared.MemoryLevel {
	wynik := make([]shared.MemoryLevel, 0, len(poziomy))
	for _, poziom := range poziomy {
		wynik = append(wynik, shared.MemoryLevel(poziom))
	}
	return wynik
}

// widocznoscKontraktu przekłada zasięg widoczności z bazy. Wartość nieznana
// czyta się jako `global`: ekspert ma się pokazać, a nie zniknąć z biblioteki.
func widocznoscKontraktu(widocznosc string) shared.AgentVisibility {
	wartosc := shared.AgentVisibility(strings.TrimSpace(widocznosc))
	if _, jest := widocznosciEksperta[wartosc]; !jest {
		return shared.AgentVisibilityGlobal
	}
	return wartosc
}

// sprawdzPoziomyPamieci pilnuje, żeby każdy podany poziom należał do kontraktu.
// Wycinek pusty jest żądaniem poprawnym — znaczy wyłączenie pamięci.
func sprawdzPoziomyPamieci(poziomy []shared.MemoryLevel) ([]string, error) {
	wynik := make([]string, 0, len(poziomy))
	for _, poziom := range poziomy {
		if _, jest := poziomyPamieci[poziom]; !jest {
			return nil, bladZadaniaEksperta("poziom pamięci " + string(poziom) +
				" nie należy do kontraktu")
		}
		wynik = append(wynik, string(poziom))
	}
	return wynik, nil
}

// sprawdzWidocznosc pilnuje, żeby zgłoszony zasięg widoczności należał
// do zbioru rozpoznawanego przez kontrakt.
func sprawdzWidocznosc(widocznosc *shared.AgentVisibility) error {
	if widocznosc == nil {
		return nil
	}
	if _, jest := widocznosciEksperta[*widocznosc]; !jest {
		return bladZadaniaEksperta("widoczność " + string(*widocznosc) +
			" nie należy do kontraktu")
	}
	return nil
}

// uprawnieniaKontraktu przekłada wiersze uprawnień na strukturę kontraktu.
// Zakres pusty zostaje pominięty — kontrakt opisuje cała grupa brakiem wartości.
func uprawnieniaKontraktu(wiersze []dane.UprawnienieAgenta) []shared.AgentPermission {
	if len(wiersze) == 0 {
		return nil
	}
	uprawnienia := make([]shared.AgentPermission, 0, len(wiersze))
	for _, wiersz := range wiersze {
		uprawnienia = append(uprawnienia, shared.AgentPermission{
			Group:   shared.AgentPermissionGroup(wiersz.Grupa),
			Scope:   wskaznikTekstu(wiersz.Zakres),
			Granted: wiersz.Przyznane,
		})
	}
	return uprawnienia
}

// konektorKontraktu przekłada wiersz konektora na strukturę AgentConnector.
// Punkt dostępu wychodzi kodem, tym samym, którym posługują się komendy
// `access.point.*`.
func konektorKontraktu(k dane.KonektorAgenta, kodPunktu string) shared.AgentConnector {
	konektor := shared.AgentConnector{
		Id:            k.Kod,
		AgentId:       k.AgentKod,
		Name:          k.Nazwa,
		Kind:          shared.AgentConnectorKind(k.Rodzaj),
		AccessPointId: wskaznikTekstu(kodPunktu),
		Enabled:       k.Aktywny,
	}
	if tresc := strings.TrimSpace(k.Konfiguracja); tresc != "" && tresc != "{}" {
		konektor.Config = json.RawMessage(tresc)
	}
	if chwila, ok := chwilaZapisu(k.Utworzono); ok {
		konektor.CreatedAt = chwila
	}
	return konektor
}

// warstwyKontraktu przekłada warstwy promptu eksperta. Porządek przychodzi
// z bazy ułożony wg krytyczności i nie jest tu układany po raz drugi.
func warstwyKontraktu(wiersze []dane.WarstwaAgenta) []shared.AgentLayer {
	if len(wiersze) == 0 {
		return nil
	}
	warstwy := make([]shared.AgentLayer, 0, len(wiersze))
	for _, wiersz := range wiersze {
		warstwa := shared.AgentLayer{
			Layer:   shared.IdentityLayer(wiersz.Warstwa),
			Content: wiersz.Tresc,
			Enabled: wiersz.Aktywna,
		}
		if chwila, ok := chwilaZapisu(wiersz.Zaktualizowano); ok {
			warstwa.UpdatedAt = int(chwila)
		}
		warstwy = append(warstwy, warstwa)
	}
	return warstwy
}

// wtyczkaKontraktu przekłada wiersz wtyczki na strukturę AgentPlugin, kodem
// trwałym, nie numerem wiersza — tak samo jak konektor.
func wtyczkaKontraktu(w dane.WtyczkaAgenta) shared.AgentPlugin {
	wtyczka := shared.AgentPlugin{
		Id:      w.Kod,
		AgentId: w.AgentKod,
		Name:    w.Nazwa,
		Source:  w.Zrodlo,
		Version: w.Wersja,
		Enabled: w.Aktywna,
	}
	if chwila, ok := chwilaZapisu(w.Utworzono); ok {
		wtyczka.CreatedAt = int(chwila)
	}
	return wtyczka
}

// kodyWtyczek zbiera kody wtyczek do pola `Agent.pluginIds`. Ekspert niesie
// wskazania, a nie treść — pełna struktura wraca wynikiem `agent.plugin.add`.
func kodyWtyczek(wiersze []dane.WtyczkaAgenta) []string {
	if len(wiersze) == 0 {
		return nil
	}
	kody := make([]string, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kody = append(kody, wiersz.Kod)
	}
	return kody
}

// bladZadaniaEksperta odmawia wykonania komendy o niepoprawnej treści
// żądania obszaru `agent.*`, kodem błędu kontraktu.
func bladZadaniaEksperta(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "moduł Agents: "+powod))
}

// sprawdzTransport pilnuje, żeby zgłoszona droga wywołania dostawcy należała
// do zbioru rozpoznawanego przez kontrakt.
func sprawdzTransport(transport *shared.ProviderTransport) error {
	if transport == nil {
		return nil
	}
	if _, jest := transportyKanalu[*transport]; !jest {
		return bladZadaniaEksperta("transport " + string(*transport) + " nie należy do kontraktu")
	}
	return nil
}

// sprawdzParametry pilnuje, żeby parametry wywołania konektora były poprawnym
// zapisem JSON, sprawdzanym przed przekazaniem do dostawcy.
func sprawdzParametry(surowe json.RawMessage, nazwaPola string) (string, error) {
	tresc := strings.TrimSpace(string(surowe))
	if tresc == "" {
		return "", nil
	}
	if !json.Valid([]byte(tresc)) {
		return "", bladZadaniaEksperta(nazwaPola + " nie są poprawnym JSON")
	}
	return tresc, nil
}

// trybKontraktu przekłada tryb nałożenia z bazy na wartość kontraktu. Wartość
// spoza katalogu czyta się jako dołączenie, nigdy jako zastąpienie promptu.
func trybKontraktu(tryb string) *shared.IdentityMode {
	wartosc := shared.IdentityMode(shared.IdentityModeDOLACZ)
	if shared.IdentityMode(strings.TrimSpace(tryb)) == shared.IdentityModeZASTAP {
		wartosc = shared.IdentityMode(shared.IdentityModeZASTAP)
	}
	return &wartosc
}
