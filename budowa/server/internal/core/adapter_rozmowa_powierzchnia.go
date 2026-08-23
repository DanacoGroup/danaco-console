// Odpowiedzialność pliku: przekład obszarów tools, permissions, hooks, skills,
// environment, provider i mcp konfiguracji sesji na trzy powierzchnie procesu
// kanału, które warstwa injection już potrafi złożyć:
//
//   - napis --settings (plik ustawień sesji) — reguły uprawnień i narzędzi,
//     sekcja hooks (zaczepy cyklu życia) oraz odmowa narzędzia Skill, gdy obszar
//     skills jest wyłączony;
//   - zmienne środowiskowe procesu — obszar environment oraz adres dostawcy;
//   - osobne --mcp-config — wiązania serwerów MCP opisane wprost.
//
// Ten plik jest jedynym miejscem w drzewie, które zna kształt pliku ustawień
// dostawcy i nazwy zmiennych środowiskowych. Model konfiguracji pozostaje
// dziedzinowy — dopiero tutaj staje się „permissions.allow” czy
// „ANTHROPIC_BASE_URL”.
package core

import (
	"encoding/json"
	"strconv"

	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// Nazwy zmiennych środowiskowych, które program kanału rozumie. Trzymane w
// jednym miejscu, bo to jedyne miejsce w drzewie, w którym wolno je znać.
const (
	zmiennaAdresDostawcy     = "ANTHROPIC_BASE_URL"
	zmiennaGranicaOdpowiedzi = "CLAUDE_CODE_MAX_OUTPUT_TOKENS"
	zmiennaStrefaCzasowa     = "TZ"
	zmiennaJezyka            = "LANG"
	zmiennaJezykaWymuszona   = "LC_ALL"
)

const (
	// typZaczepuPolecenie jest jedynym rodzajem zaczepu, który powierzchnia
	// niesie: polecenie powłoki wywoływane w punkcie cyklu życia. Program kanału
	// rozumie ten napis dosłownie w polu type wpisu hooks.
	typZaczepuPolecenie = "command"
	// nazwaNarzedziaUmiejetnosc jest nazwą narzędzia, którym model sięga po
	// umiejętności. Wyłączenie obszaru skills odmawia właśnie tego narzędzia —
	// tą samą drogą, którą obszar tools odmawia narzędzi imiennych.
	nazwaNarzedziaUmiejetnosc = "Skill"
)

// plikUstawienCLI jest kształtem napisu --settings ograniczonym do tego, co
// przekład wypełnia z konfiguracji sesji. Pole puste znika z JSON (omitempty),
// więc pusty obszar nie daje pliku.
type plikUstawienCLI struct {
	Permissions *uprawnieniaCLI               `json:"permissions,omitempty"`
	Hooks       map[string][]grupaZaczepowCLI `json:"hooks,omitempty"`
}

// grupaZaczepowCLI jest jedną pozycją listy zaczepów pod kluczem zdarzenia w
// sekcji hooks pliku ustawień. Zaczepy o tym samym zdarzeniu i tym samym
// zawężeniu (matcher) zbierają się w jednej grupie, dokładnie jak w pliku
// ustawień dostawcy.
type grupaZaczepowCLI struct {
	Matcher string               `json:"matcher,omitempty"`
	Hooks   []zaczepPoleceniaCLI `json:"hooks"`
}

// zaczepPoleceniaCLI jest jednym poleceniem zaczepu. Granica czasu jest w
// SEKUNDACH (pole timeout pliku ustawień), gdy tymczasem kontrakt trzyma ją w
// milisekundach — przeliczenie zapada w hooksZKonfiguracji.
type zaczepPoleceniaCLI struct {
	Type    string `json:"type"`
	Command string `json:"command"`
	Timeout *int   `json:"timeout,omitempty"`
}

// uprawnieniaCLI odwzorowuje obszar permissions oraz tools na sekcję permissions
// pliku ustawień dostawcy. Trybu domyślnego tu nie ma: tryb uprawnień jedzie
// przełącznikiem --permission-mode (z.TrybUprawnien, przelozUprawnienia), którego
// słownik kontrakt potwierdza dosłownie — powtórzenie go w pliku groziłoby
// rozejściem słownictwa i dwoma źródłami tej samej decyzji.
type uprawnieniaCLI struct {
	Allow []string `json:"allow,omitempty"`
	Deny  []string `json:"deny,omitempty"`
	Ask   []string `json:"ask,omitempty"`
}

// wpisSerweraMCP jest jedną pozycją mapy mcpServers napisu --mcp-config. Kształt
// jest szerszy niż WpisMostu (most_mcp.go): niesie i program stdio, i adres
// serwera sieciowego, bo obszar mcp konfiguracji sesji zna oba transporty.
type wpisSerweraMCP struct {
	Type    string            `json:"type,omitempty"`
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	URL     string            `json:"url,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

// konfiguracjaSerwerowMCP jest kompletem wiązań serwerów MCP sesji.
type konfiguracjaSerwerowMCP struct {
	McpServers map[string]wpisSerweraMCP `json:"mcpServers"`
}

// przelozPowierzchnieProcesu składa plik ustawień, zmienne środowiska oraz
// konfigurację MCP z konfiguracji sesji i wpisuje je do zapytania. Powierzchnia
// pusta nie nadpisuje niczego — brak reguł znaczy brak przełącznika.
func przelozPowierzchnieProcesu(k shared.SessionConfig, z *models.Zapytanie) {
	if plik := plikUstawienZKonfiguracji(k); plik != "" {
		z.PlikUstawien = plik
	}
	if srodowisko := srodowiskoZKonfiguracji(k); len(srodowisko) > 0 {
		z.Srodowisko = srodowisko
	}
	if mcp := mcpZKonfiguracji(k); mcp != "" {
		z.DodatkoweMCP = append(z.DodatkoweMCP, mcp)
	}
}

// plikUstawienZKonfiguracji buduje napis --settings z obszarów permissions,
// tools, skills oraz hooks. Reguły narzędzi i wyłączenie obszaru skills dokładają
// się do reguł uprawnień; zaczepy jadą osobną sekcją hooks. Brak reguły i brak
// zaczepu daje napis pusty, czyli brak przełącznika.
func plikUstawienZKonfiguracji(k shared.SessionConfig) string {
	plik := plikUstawienCLI{
		Permissions: uprawnieniaZKonfiguracji(k),
		Hooks:       hooksZKonfiguracji(k),
	}
	if plik.Permissions == nil && len(plik.Hooks) == 0 {
		return ""
	}
	tresc, err := json.Marshal(plik)
	if err != nil {
		return ""
	}
	return string(tresc)
}

// uprawnieniaZKonfiguracji zbiera reguły uprawnień, narzędzi oraz umiejętności.
// Zwraca nil, gdy żaden obszar nie dał reguły — pusty obiekt uprawnień nie
// miałby po co jechać.
func uprawnieniaZKonfiguracji(k shared.SessionConfig) *uprawnieniaCLI {
	uprawnienia := uprawnieniaCLI{}
	if k.Permissions != nil {
		uprawnienia.Allow = append(uprawnienia.Allow, k.Permissions.AllowRules...)
		uprawnienia.Deny = append(uprawnienia.Deny, k.Permissions.DenyRules...)
		uprawnienia.Ask = append(uprawnienia.Ask, k.Permissions.AskRules...)
	}
	dodajRegulyNarzedzi(k.Tools, &uprawnienia)
	dodajRegulyUmiejetnosci(k.Skills, &uprawnienia)
	if len(uprawnienia.Allow) == 0 && len(uprawnienia.Deny) == 0 && len(uprawnienia.Ask) == 0 {
		return nil
	}
	return &uprawnienia
}

// dodajRegulyUmiejetnosci przekłada obszar skills na regułę uprawnień. Jedyny
// przekład, który powierzchnia pliku ustawień unosi bez atrapy: wyłączenie
// obszaru wprost (Enabled == false) odmawia narzędzia Skill w sekcji deny —
// tą samą drogą, którą obszar tools odmawia narzędzi imiennych. Obszar włączony
// albo nieokreślony nie dokłada reguły: umiejętności pozostają wtedy
// dostępne, jak przed wpięciem. Dopuszczanie imienne (AllowedSkillIds), katalogi
// wyszukiwania (Directories) i samowykrywanie (AutoDiscovery) nie mają pola na
// tej powierzchni i jadą do ryzyk — nie ma tu dla nich cichej atrapy.
func dodajRegulyUmiejetnosci(skills *shared.SessionConfigSkills, uprawnienia *uprawnieniaCLI) {
	if skills == nil || skills.Enabled == nil || *skills.Enabled {
		return
	}
	uprawnienia.Deny = append(uprawnienia.Deny, nazwaNarzedziaUmiejetnosc)
}

// hooksZKonfiguracji buduje sekcję hooks pliku ustawień z obszaru hooks. Obszar
// wyłączony wprost (Enabled == false) nie daje żadnego zaczepu; obszar włączony
// albo nieokreślony przenosi zaczepy czynne. Zaczep bez zdarzenia albo bez
// polecenia jest niekompletny i nie jedzie. Zaczepy o tym samym
// zdarzeniu i zawężeniu zbierają się w jednej grupie. Brak zaczepów daje nil,
// więc sekcja znika z JSON.
func hooksZKonfiguracji(k shared.SessionConfig) map[string][]grupaZaczepowCLI {
	if k.Hooks == nil {
		return nil
	}
	if k.Hooks.Enabled != nil && !*k.Hooks.Enabled {
		return nil
	}
	wynik := map[string][]grupaZaczepowCLI{}
	for _, zaczep := range k.Hooks.Hooks {
		if !zaczep.Enabled || zaczep.Event == "" || zaczep.Command == "" {
			continue
		}
		polecenie := zaczepPoleceniaCLI{Type: typZaczepuPolecenie, Command: zaczep.Command}
		if granica := granicaSekund(zaczep.TimeoutMs); granica != nil {
			polecenie.Timeout = granica
		}
		wpiszZaczep(wynik, zaczep.Event, wartoscTekstu(zaczep.Matcher), polecenie)
	}
	if len(wynik) == 0 {
		return nil
	}
	return wynik
}

// wpiszZaczep dokłada polecenie do grupy o danym zdarzeniu i zawężeniu, zakładając
// grupę, gdy jeszcze jej nie ma. Dzięki temu wiele zaczepów jednego zdarzenia
// i jednego zawężenia trafia do wspólnej listy hooks.
func wpiszZaczep(wynik map[string][]grupaZaczepowCLI, zdarzenie, matcher string,
	polecenie zaczepPoleceniaCLI) {
	grupy := wynik[zdarzenie]
	for i := range grupy {
		if grupy[i].Matcher == matcher {
			grupy[i].Hooks = append(grupy[i].Hooks, polecenie)
			wynik[zdarzenie] = grupy
			return
		}
	}
	wynik[zdarzenie] = append(grupy, grupaZaczepowCLI{
		Matcher: matcher,
		Hooks:   []zaczepPoleceniaCLI{polecenie},
	})
}

// granicaSekund przelicza granicę czasu zaczepu z milisekund kontraktu na
// sekundy pliku ustawień. Wartość niedodatnia nie daje granicy (nil); wartość
// dodatnia poniżej sekundy zaokrągla w górę do jednej sekundy — pole timeout nie
// wyraża ułamka, a granica poniżej pełnej sekundy nie może zejść do zera i
// zamienić się w brak granicy.
func granicaSekund(timeoutMs *int) *int {
	if timeoutMs == nil || *timeoutMs <= 0 {
		return nil
	}
	sekundy := *timeoutMs / 1000
	if sekundy < 1 {
		sekundy = 1
	}
	return &sekundy
}

// dodajRegulyNarzedzi przekłada obszar tools na reguły uprawnień. Lista
// dopuszczonych liczy się przy polityce allowList, odmówionych — przy denyList;
// polityki all i none nie wnoszą reguł imiennych.
func dodajRegulyNarzedzi(tools *shared.SessionConfigTools, uprawnienia *uprawnieniaCLI) {
	if tools == nil || tools.Policy == nil {
		return
	}
	switch *tools.Policy {
	case shared.ToolPolicyAllowList:
		uprawnienia.Allow = append(uprawnienia.Allow, tools.AllowedToolNames...)
	case shared.ToolPolicyDenyList:
		uprawnienia.Deny = append(uprawnienia.Deny, tools.DeniedToolNames...)
	}
}

// srodowiskoZKonfiguracji składa zmienne środowiskowe procesu z obszaru
// environment oraz z adresu i granicy odpowiedzi obszarów provider i model.
// Zmienna tajna (SecretRef) nie wchodzi: jej treść zna wyłącznie sejf, którego
// ta droga nie ma wpiętego — wpisanie nazwy bez wartości byłoby atrapą.
func srodowiskoZKonfiguracji(k shared.SessionConfig) map[string]string {
	srodowisko := map[string]string{}
	if k.Environment != nil {
		for _, zmienna := range k.Environment.Variables {
			if !zmienna.Enabled || zmienna.Name == "" || zmienna.Value == nil {
				continue
			}
			srodowisko[zmienna.Name] = *zmienna.Value
		}
		ustawNiepusteSrodowisko(srodowisko, zmiennaStrefaCzasowa, wartoscTekstu(k.Environment.TimeZone))
		if jezyk := wartoscTekstu(k.Environment.LocaleTag); jezyk != "" {
			srodowisko[zmiennaJezyka] = jezyk
			srodowisko[zmiennaJezykaWymuszona] = jezyk
		}
	}
	if k.Provider != nil {
		ustawNiepusteSrodowisko(srodowisko, zmiennaAdresDostawcy, wartoscTekstu(k.Provider.EndpointUrl))
	}
	if k.Model != nil && k.Model.MaxOutputTokens != nil {
		srodowisko[zmiennaGranicaOdpowiedzi] = strconv.Itoa(*k.Model.MaxOutputTokens)
	}
	return srodowisko
}

// mcpZKonfiguracji buduje napis --mcp-config z wiązań obszaru mcp opisanych
// wprost (transport stdio z programem albo sse/http z adresem). Wiązania
// wskazujące punkt dostępu (AccessPointId) pomija: ich adres i poświadczenie
// żyją w rejestrze punktów dostępu, którego ta droga nie rozstrzyga —
// jadą one drogą nadań okna (mosty). Brak wiązań opisanych wprost daje napis
// pusty, czyli brak przełącznika.
func mcpZKonfiguracji(k shared.SessionConfig) string {
	if k.Mcp == nil || len(k.Mcp.Servers) == 0 {
		return ""
	}
	serwery := map[string]wpisSerweraMCP{}
	for _, wiazanie := range k.Mcp.Servers {
		if !wiazanie.Enabled || wiazanie.Name == "" || wiazanie.AccessPointId != nil {
			continue
		}
		if wpis, poprawny := wpisSerweraZWiazania(wiazanie); poprawny {
			serwery[wiazanie.Name] = wpis
		}
	}
	if len(serwery) == 0 {
		return ""
	}
	tresc, err := json.Marshal(konfiguracjaSerwerowMCP{McpServers: serwery})
	if err != nil {
		return ""
	}
	return string(tresc)
}

// wpisSerweraZWiazania przekłada jedno wiązanie opisane wprost na pozycję mapy
// mcpServers. Transport stdio wymaga programu, sse i http — adresu; wiązanie
// bez wymaganego pola nie daje wpisu.
func wpisSerweraZWiazania(wiazanie shared.McpServerBinding) (wpisSerweraMCP, bool) {
	switch wiazanie.Transport {
	case shared.McpTransportStdio:
		polecenie := wartoscTekstu(wiazanie.Command)
		if polecenie == "" {
			return wpisSerweraMCP{}, false
		}
		return wpisSerweraMCP{
			Command: polecenie,
			Args:    wiazanie.Arguments,
			Env:     srodowiskoWiazaniaMCP(wiazanie.Environment),
		}, true
	case shared.McpTransportSse, shared.McpTransportHttp:
		adres := wartoscTekstu(wiazanie.EndpointUrl)
		if adres == "" {
			return wpisSerweraMCP{}, false
		}
		return wpisSerweraMCP{Type: string(wiazanie.Transport), URL: adres}, true
	default:
		return wpisSerweraMCP{}, false
	}
}

// srodowiskoWiazaniaMCP zbiera jawne zmienne środowiska serwera MCP. Zmienna
// tajna (SecretRef) nie wchodzi z tego samego powodu co w środowisku procesu:
// jej treści ta droga nie zna. Brak zmiennych daje nil.
func srodowiskoWiazaniaMCP(zmienne []shared.EnvironmentVariable) map[string]string {
	wynik := map[string]string{}
	for _, zmienna := range zmienne {
		if zmienna.Enabled && zmienna.Name != "" && zmienna.Value != nil {
			wynik[zmienna.Name] = *zmienna.Value
		}
	}
	if len(wynik) == 0 {
		return nil
	}
	return wynik
}

// ustawNiepusteSrodowisko wpisuje zmienną, gdy jej wartość jest niepusta.
func ustawNiepusteSrodowisko(srodowisko map[string]string, nazwa, wartosc string) {
	if wartosc != "" {
		srodowisko[nazwa] = wartosc
	}
}
