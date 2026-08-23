package models

import (
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// Zasiegi niosą kontekst, w którym powstało zapytanie: środowisko platformy,
// projekt, sesja oraz okno komunikacji. Okno jest zasięgiem najwęższym i to ono
// rozstrzyga parametry wykonania. Kanał nie rozstrzyga
// konfiguracji sam — dostaje zasięgi po to, by móc je przekazać dalej
// i odnotować w prowenancji.
type Zasiegi struct {
	Srodowisko string `json:"environmentId,omitempty"`
	Projekt    string `json:"projectId,omitempty"`
	Sesja      string `json:"sessionId,omitempty"`
	Okno       string `json:"windowId"`
}

// Nakladka jest warstwową treścią systemową wywołania, uporządkowaną wg
// krytyczności: konstytucja → profil/rola → ekspertyza zadaniowa.
// Każda warstwa pochodzi z konfiguracji; pakiet nie zna żadnej jej treści.
type Nakladka struct {
	Konstytucja string `json:"constitution,omitempty"`
	ProfilRoli  string `json:"roleProfile,omitempty"`
	Ekspertyza  string `json:"expertise,omitempty"`
	// Tryb rozstrzyga, czy nakładka zastępuje prompt fabryczny w całości, czy
	// dopisuje się do niego. Wartości: injection.TrybZastap / injection.TrybDopisz;
	// pusty znaczy tryb dopisania.
	Tryb string `json:"overlayMode,omitempty"`
}

// PromptSystemowy skleja warstwy w jeden prompt systemowy z zachowaniem
// kolejności krytyczności. Brak warstwy nie jest błędem — nakładka pusta daje
// pusty prompt i wywołanie idzie dalej.
func (n Nakladka) PromptSystemowy() string {
	warstwy := make([]string, 0, 3)
	for _, warstwa := range []string{n.Konstytucja, n.ProfilRoli, n.Ekspertyza} {
		if strings.TrimSpace(warstwa) != "" {
			warstwy = append(warstwy, warstwa)
		}
	}
	return strings.Join(warstwy, "\n\n")
}

// Skrot zwraca skrót diagnostyczny nakładki.
func (n Nakladka) Skrot() string {
	return SkrotTresci(n.Konstytucja, n.ProfilRoli, n.Ekspertyza)
}

// RolaWypowiedzi rozstrzyga, czyja jest wypowiedź w historii rozmowy: pytającego
// czy modelu. Wartości pokrywają się z rolami, którymi posługuje się większość
// punktów końcowych rozmowy (pole "role" wiadomości).
const (
	RolaUzytkownika = "user"
	RolaAsystenta   = "assistant"
)

// WiadomoscHistorii jest jedną wypowiedzią wcześniejszej tury tej samej rozmowy.
// Historia jest strukturą wewnętrzną pakietu, nie kontraktem: adapter kanału
// składa z niej pamięć wywołania, a wypełnia ją warstwa rozmowy z dziennika
// tury. Rola pusta znaczy wypowiedź pytającego — pamięć bez roli nie może
// wywrócić wywołania.
type WiadomoscHistorii struct {
	Rola  string `json:"role,omitempty"`
	Tresc string `json:"content"`
}

// RolaLub zwraca rolę wypowiedzi, a przy jej braku rolę pytającego.
func (w WiadomoscHistorii) RolaLub() string {
	if strings.TrimSpace(w.Rola) != "" {
		return w.Rola
	}
	return RolaUzytkownika
}

// Role obrazu wejściowego. Rola rozstrzyga, czym obraz jest dla wywołania:
// materiałem, na którym model pracuje, czy maską wskazującą obszar pracy.
// Rozdzielenie jest konieczne, bo dwie czynności warsztatu fotografii
// (uzupełnienie ubytku i rozszerzenie kadru) wysyłają OBA naraz, a punkty
// końcowe przyjmują je osobnymi polami. Rola pusta znaczy materiał — obraz bez
// roli nie może wywrócić wywołania.
const (
	RolaObrazuMaterial = "obraz"
	RolaObrazuMaska    = "maska"
)

// ObrazWejsciowy jest obrazem WEJŚCIOWYM wywołania kanału — materiałem, który
// jedzie DO modelu, nie wynikiem, który z niego wraca.
//
// Bez tego pola kanał obrazowy umiał wyłącznie wygenerować obraz z samego
// polecenia tekstowego, więc cztery czynności warsztatu fotografii modułu Design
// (powiększenie, odcięcie tła, uzupełnienie ubytku i rozszerzenie kadru) nie
// miały jak dojść do wariantu neuronowego, choć mają go lepszy od rachunku na
// pikselach.
//
// Bajty jadą WYŁĄCZNIE wprost, zapisem base64 — bez wariantu adresowego, jaki
// ma TrescObrazu w drugą stronę. Powód jest jednostronny: odpowiedź dostawcy
// bywa odsyłaczem, bo to dostawca trzyma plik; materiał wejściowy leży
// w magazynie rdzenia i punkt końcowy nie ma jak po niego sięgnąć. Adres
// przyjęty tutaj byłby adresem, którego druga strona nie odczyta.
type ObrazWejsciowy struct {
	// Rola nazywa przeznaczenie obrazu — RolaObrazuMaterial albo RolaObrazuMaska.
	Rola string `json:"role,omitempty"`
	// TypTresci nazywa rodzaj bajtów (np. „image/png").
	TypTresci string `json:"mimeType,omitempty"`
	// Base64 niesie bajty obrazu w zapisie base64.
	Base64 string `json:"contentBase64"`
	// Nazwa jest nazwą pliku podawaną przy wysyłce wieloczęściowej. Punkty
	// końcowe rozpoznają po niej format, gdy nie dostaną typu treści, więc pusta
	// nazwa nie jest brakiem — adapter dokłada wtedy własną.
	Nazwa string `json:"fileName,omitempty"`
}

// RolaLub zwraca rolę obrazu, a przy jej braku rolę materiału.
func (o ObrazWejsciowy) RolaLub() string {
	if strings.TrimSpace(o.Rola) != "" {
		return o.Rola
	}
	return RolaObrazuMaterial
}

// Zapytanie jest jednym wywołaniem kanału modelu. Niesie treść użytkownika,
// kontekst zasięgów oraz komplet parametrów wykonania okna, dzięki czemu adapter
// nie sięga po konfigurację na własną rękę.
type Zapytanie struct {
	Zasiegi   Zasiegi `json:"scopes"`
	Wiadomosc string  `json:"messageId"`
	Tresc     string  `json:"content"`

	// Historia niesie wcześniejsze wypowiedzi tej samej rozmowy w kolejności
	// nadania, bez bieżącej wypowiedzi (ta jedzie w Tresc). Kanał bezstanowy
	// (api) składa z niej pamięć wywołania; kanał z własną pamięcią (cli
	// przez Wznowienie) może ją pominąć. Puste znaczy turę bez pamięci.
	Historia []WiadomoscHistorii `json:"history,omitempty"`

	// ObrazyWejsciowe niosą materiał wizualny wywołania. Puste znaczy wywołanie
	// z samego polecenia — tak pracuje generowanie obrazu od zera i tak pracują
	// wszystkie kanały tekstowe, które to pole po prostu pomijają. Kolejność jest
	// kolejnością nadania: kanał obrazowy bierze pierwszy materiał i pierwszą
	// maskę, bo punkty końcowe edycji przyjmują po jednym z każdej roli.
	ObrazyWejsciowe []ObrazWejsciowy `json:"inputImages,omitempty"`

	Kanal             string `json:"channelId,omitempty"`
	Model             string `json:"model,omitempty"`
	ModelZapasowy     string `json:"fallbackModel,omitempty"`
	NakladRozumowania string `json:"effort,omitempty"`
	Konto             string `json:"account,omitempty"`
	Wznowienie        string `json:"resume,omitempty"`
	// PulapKosztuUSD jest górną granicą kosztu tej tury w dolarach; pochodzi
	// z nastawy `pulap_kosztu_usd`. Zero znaczy brak pułapu — kanał nie
	// dopisze wtedy przełącznika --max-budget-usd.
	PulapKosztuUSD float64 `json:"maxBudgetUsd,omitempty"`

	// KatalogSesji jest katalogiem roboczym sesji ustalonym z konfiguracji
	// Operatora (klucze `katalog.roboczy.*`). Nie jest nadaniem dostępu — mówi,
	// gdzie model zostawia własne pliki, a nie do czego sięga.
	KatalogSesji string `json:"sessionDir,omitempty"`
	// KonfiguracjaMCP niesie tekst konfiguracji mostów MCP okna albo ścieżkę do
	// pliku z nią. Pusty znaczy: okno nie ma nadanych mostów.
	KonfiguracjaMCP string `json:"mcpConfig,omitempty"`
	// DodatkoweMCP niesie konfiguracje mostów MCP wyliczone z obszaru mcp
	// konfiguracji sesji. Jadą osobnymi przełącznikami --mcp-config obok
	// KonfiguracjaMCP okna, więc nadania okna i wiązania sesji nie przykrywają
	// się nawzajem. Puste znaczy: obszar mcp nie dał żadnego serwera.
	DodatkoweMCP []string `json:"sessionMcpConfigs,omitempty"`
	// PlikUstawien niesie napis JSON pliku ustawień sesji (--settings) złożony
	// z obszarów tools i permissions konfiguracji sesji. Pusty znaczy: obszary
	// nie dały żadnej reguły, więc przełącznika nie ma.
	PlikUstawien string `json:"settingsFile,omitempty"`
	// Srodowisko niesie zmienne środowiskowe wyliczone z obszaru environment
	// (oraz provider) konfiguracji sesji. Dokładają się do środowiska procesu
	// kanału. Puste znaczy: obszar nie dał żadnej zmiennej.
	Srodowisko map[string]string `json:"sessionEnv,omitempty"`

	KatalogiRobocze     []string              `json:"workingDirs,omitempty"`
	SrodowiskoWykonania shared.ExecutionEnv   `json:"executionEnv,omitempty"`
	TrybUprawnien       shared.PermissionMode `json:"permissionMode,omitempty"`
	RolaOkna            shared.WindowRole     `json:"windowRole,omitempty"`

	Nakladka   Nakladka          `json:"overlay"`
	Ustawienia map[string]string `json:"settings,omitempty"`
}

// Okno zwraca identyfikator okna, do którego należy zapytanie.
func (z Zapytanie) Okno() string {
	return z.Zasiegi.Okno
}

// ObrazRoli zwraca pierwszy obraz wejściowy wskazanej roli. Drugi obraz tej
// samej roli jest pomijany, a nie sklejany: punkty końcowe edycji przyjmują po
// jednym materiale i po jednej masce, więc sklejenie byłoby wymyślaniem kształtu
// żądania.
func (z Zapytanie) ObrazRoli(rola string) (ObrazWejsciowy, bool) {
	for _, obraz := range z.ObrazyWejsciowe {
		if obraz.RolaLub() == rola && strings.TrimSpace(obraz.Base64) != "" {
			return obraz, true
		}
	}
	return ObrazWejsciowy{}, false
}

// MaMaterialWejsciowy mówi, czy wywołanie niesie obraz, na którym model ma
// pracować. Sama maska materiałem nie jest — maska bez materiału nie wskazuje
// niczego, więc kanał ma wtedy pójść drogą generowania od zera.
func (z Zapytanie) MaMaterialWejsciowy() bool {
	_, jest := z.ObrazRoli(RolaObrazuMaterial)
	return jest
}

// WybranyModel rozstrzyga model wywołania: wygrywa wskazanie okna, w jego braku
// model zapisany w wierszu rejestru kanałów. Brak obu daje pusty napis — kanał
// używa wtedy własnej wartości domyślnej, a nie odmawia pracy.
func (z Zapytanie) WybranyModel(d Definicja) string {
	if strings.TrimSpace(z.Model) != "" {
		return z.Model
	}
	return d.Model
}

// WybraneKonto rozstrzyga konto wywołania: wygrywa wskazanie zapytania (obszar
// account konfiguracji sesji), a w jego braku obowiązuje konto powiązane
// z wierszem kanału kolumną kanal_modelu.konto_id. Powiązanie kanału z
// kontem staje się dzięki temu widoczne w wywołaniu nawet bez konfiguracji
// sesji. Brak obu daje pusty napis — tożsamość bierze się wtedy z otoczenia
// procesu, a kanał nie odmawia pracy.
func (z Zapytanie) WybraneKonto(d Definicja) string {
	if strings.TrimSpace(z.Konto) != "" {
		return z.Konto
	}
	if d.KontoId != 0 {
		return strconv.FormatInt(d.KontoId, 10)
	}
	return ""
}
