package models

import (
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// Zasiegi niosą kontekst, w którym powstało zapytanie: środowisko platformy,
// projekt, sesja oraz okno komunikacji. Okno jest zasięgiem najwęższym i
// rozstrzyga parametry wykonania.
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
	// Tryb rozstrzyga, czy nakładka zastępuje prompt fabryczny, czy się dopisuje;
	// pusty znaczy dopisanie.
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

// Skrot zwraca skrót diagnostyczny nakładki, złożony ze wszystkich jej warstw
// treści systemowej, do celów dziennika.
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
// Rola pusta znaczy wypowiedź pytającego — pamięć bez roli nie może wywrócić
// wywołania.
type WiadomoscHistorii struct {
	Rola  string `json:"role,omitempty"`
	Tresc string `json:"content"`
}

// RolaLub zwraca rolę wypowiedzi historii rozmowy, a przy jej braku zwraca rolę
// pytającego jako wartość domyślną.
func (w WiadomoscHistorii) RolaLub() string {
	if strings.TrimSpace(w.Rola) != "" {
		return w.Rola
	}
	return RolaUzytkownika
}

// Role obrazu wejściowego. Rola rozstrzyga, czym obraz jest dla wywołania:
// materiałem, na którym model pracuje, czy maską wskazującą obszar pracy.
// Rola pusta znaczy materiał.
const (
	RolaObrazuMaterial = "obraz"
	RolaObrazuMaska    = "maska"
)

// ObrazWejsciowy jest obrazem wejściowym wywołania kanału — materiałem, który
// jedzie do modelu, nie wynikiem, który z niego wraca. Bajty jadą wyłącznie
// wprost, zapisem base64.
type ObrazWejsciowy struct {
	// Rola nazywa przeznaczenie obrazu — RolaObrazuMaterial albo RolaObrazuMaska.
	Rola string `json:"role,omitempty"`
	// TypTresci nazywa rodzaj bajtów (np. „image/png").
	TypTresci string `json:"mimeType,omitempty"`
	// Base64 niesie bajty obrazu w zapisie base64.
	Base64 string `json:"contentBase64"`
	// Nazwa jest nazwą pliku wysyłki wieloczęściowej; pusta nie jest brakiem,
	// adapter dokłada własną.
	Nazwa string `json:"fileName,omitempty"`
}

// RolaLub zwraca rolę obrazu wejściowego, a przy jej braku zwraca rolę
// materiału jako wartość domyślną.
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

	// Historia niesie wcześniejsze wypowiedzi rozmowy w kolejności nadania; puste
	// znaczy turę bez pamięci.
	Historia []WiadomoscHistorii `json:"history,omitempty"`

	// ObrazyWejsciowe niosą materiał wizualny w kolejności nadania; puste znaczy
	// wywołanie z polecenia.
	ObrazyWejsciowe []ObrazWejsciowy `json:"inputImages,omitempty"`

	Kanal             string `json:"channelId,omitempty"`
	Model             string `json:"model,omitempty"`
	ModelZapasowy     string `json:"fallbackModel,omitempty"`
	NakladRozumowania string `json:"effort,omitempty"`
	Konto             string `json:"account,omitempty"`
	Wznowienie        string `json:"resume,omitempty"`
	// PulapKosztuUSD jest górną granicą kosztu tej tury w dolarach; zero znaczy
	// brak pułapu.
	PulapKosztuUSD float64 `json:"maxBudgetUsd,omitempty"`

	// KatalogSesji jest katalogiem roboczym sesji; nie nadaje dostępu, wskazuje
	// miejsce zapisu plików.
	KatalogSesji string `json:"sessionDir,omitempty"`
	// KonfiguracjaMCP niesie tekst konfiguracji mostów MCP okna albo ścieżkę do
	// pliku z nią.
	KonfiguracjaMCP string `json:"mcpConfig,omitempty"`
	// DodatkoweMCP niesie konfiguracje mostów MCP sesji, osobno od
	// KonfiguracjaMCP okna.
	DodatkoweMCP []string `json:"sessionMcpConfigs,omitempty"`
	// PlikUstawien niesie napis JSON pliku ustawień sesji złożony z reguł
	// narzędzi i uprawnień sesji.
	PlikUstawien string `json:"settingsFile,omitempty"`
	// Srodowisko niesie zmienne środowiskowe sesji dokładane do środowiska
	// procesu kanału.
	Srodowisko map[string]string `json:"sessionEnv,omitempty"`

	KatalogiRobocze     []string              `json:"workingDirs,omitempty"`
	SrodowiskoWykonania shared.ExecutionEnv   `json:"executionEnv,omitempty"`
	TrybUprawnien       shared.PermissionMode `json:"permissionMode,omitempty"`
	RolaOkna            shared.WindowRole     `json:"windowRole,omitempty"`

	Nakladka   Nakladka          `json:"overlay"`
	Ustawienia map[string]string `json:"settings,omitempty"`
}

// Okno zwraca identyfikator okna komunikacji, do którego to zapytanie należy
// w obrębie bieżącej sesji roboczej.
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

// WybraneKonto rozstrzyga konto wywołania: wygrywa wskazanie zapytania, w jego
// braku obowiązuje konto powiązane z wierszem kanału. Brak obu daje pusty
// napis — kanał nie odmawia pracy.
func (z Zapytanie) WybraneKonto(d Definicja) string {
	if strings.TrimSpace(z.Konto) != "" {
		return z.Konto
	}
	if d.KontoDostawcyId != 0 {
		return strconv.FormatInt(d.KontoDostawcyId, 10)
	}
	return ""
}
