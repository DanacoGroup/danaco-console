package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

// Tryby podania nakładki — te same napisy, którymi posługuje się kanał główny
// (injection.TrybDopisz / injection.TrybZastap). Pakiet models trzyma własne
// stałe, żeby nie wiązać się z pakietem injection, ale wartości muszą się
// zgadzać, bo klient czyta jeden wspólny kształt prowenancji.
const (
	// trybNakladkiDopisz dokłada nakładkę do promptu programu.
	trybNakladkiDopisz = "dopisz"
	// trybNakladkiZastap podmienia prompt w całości.
	trybNakladkiZastap = "zastap"
)

// Nazwy warstw nakładki w kolejności krytyczności — zgodne z nazwami
// warstw kanału głównego, żeby klient złożył jeden słownik warstw z każdego
// kanału.
const (
	warstwaKonstytucja = "konstytucja"
	warstwaProfil      = "profil"
	warstwaEkspertyza  = "ekspertyza"
)

// Prowenancja opisuje, co dokładnie poszło do modelu: wiersz wywołania, prompt
// systemowy złożony z warstw nakładki, ustawienia przekazane kanałowi oraz skrót
// nakładki. Jest przejrzystością wywołania, nie bramką — skrót służy wyłącznie
// diagnostyce i nigdy nie dopuszcza ani nie blokuje wywołania.
//
// Nazwy pól JSON szkieletu są wspólne z injection.Prowenancja (kanał główny
// CLI). Klient rozmowy (client/src/rozmowa/prowenancja.ts) parsuje wyłącznie
// ten jeden kształt, więc kanały sieciowe i echo emitują prowenancję o tych
// samych nazwach pól co kanał główny. Pola specyficzne kanałów bez procesu
// (kanał, adapter, adres, ustawienia wywołania, środowisko) jadą jako
// dodatkowe, poza wspólnym szkieletem, i nie kolidują z jego nazwami.
type Prowenancja struct {
	// ── Szkielet wspólny z injection.Prowenancja ─────────────────────────────
	// Program — plik wykonywalny albo klucz adaptera kanału bez procesu.
	Program string `json:"program"`
	// Argv — wiersz wywołania procesu; puste dla kanałów bez procesu.
	Argv []string `json:"argv"`
	// PromptSystemowy — pełna treść nakładki przekazana modelowi.
	PromptSystemowy string `json:"systemPrompt"`
	// TrybNakladki — dopisanie albo zastąpienie promptu (dopisz/zastap).
	TrybNakladki string `json:"overlayMode"`
	// WarstwyNakladki — nazwy warstw nakładki w kolejności złożenia.
	WarstwyNakladki []string `json:"overlayLayers"`
	// SkrotNakladki — skrót nakładki, sha256 w zapisie szesnastkowym.
	SkrotNakladki string `json:"overlayHash"`
	// Model — identyfikator modelu przekazany kanałowi.
	Model string `json:"model"`
	// ModelZapasowy — model rezerwowy wskazany oknu.
	ModelZapasowy string `json:"fallbackModel"`
	// Naklad — nakład rozumowania modelu.
	Naklad string `json:"effort"`
	// TrybUprawnien — tryb uprawnień okna.
	TrybUprawnien string `json:"permissionMode"`
	// Katalogi — katalogi robocze udostępnione modelowi w tym wywołaniu.
	Katalogi []string `json:"directories"`
	// KatalogRoboczy — katalog, w którym model zostawia własne pliki.
	KatalogRoboczy string `json:"workingDirectory"`
	// PlikUstawien — napis JSON pliku ustawień sesji przekazany kanałowi. Nazwa
	// pola JSON to „settings" — tak samo jak w kanale głównym.
	PlikUstawien string `json:"settings"`
	// KonfiguracjaMCP — wskazane konfiguracje mostów MCP.
	KonfiguracjaMCP []string `json:"mcpConfig"`
	// Wznowienie — identyfikator wznawianej rozmowy kanału.
	Wznowienie string `json:"resume"`
	// Konto — konto albo profil poświadczeń użyty w wywołaniu; nigdy sekret.
	Konto string `json:"account"`
	// KatalogKonta — katalog konfiguracji konta; puste dla kanałów bez procesu.
	KatalogKonta string `json:"accountConfigDir"`
	// Proba — numer próby wywołania. Kanały bez rotacji kont mają jedną próbę.
	Proba int `json:"attempt"`
	// PowodProby — powód ponowienia; wypełniony przy rotacji konta.
	PowodProby string `json:"attemptReason,omitempty"`
	// Chwila — chwila złożenia prowenancji, zapis ISO 8601.
	Chwila time.Time `json:"at"`

	// ── Pola specyficzne kanałów sieciowych/echo (dodatkowe) ─────────────────
	// Kanal — kod kanału z rejestru, który wykonuje wywołanie.
	Kanal string `json:"channel,omitempty"`
	// Adapter — klucz adaptera obsługującego kanał (echo, cli, api, …).
	Adapter string `json:"adapter,omitempty"`
	// Adres — punkt końcowy wywołania sieciowego; puste dla kanałów bez sieci.
	Adres string `json:"endpoint,omitempty"`
	// Ustawienia — parametry wywołania kanału; nigdy nie zawierają sekretu,
	// wyłącznie odwołania do danych dostępowych. Osobna nazwa pola
	// („callSettings"), by nie kolidować ze szkieletowym „settings".
	Ustawienia map[string]string `json:"callSettings,omitempty"`
	// SrodowiskoWykonania — gdzie model pracuje; niezależne od umiejscowienia
	// rdzenia.
	SrodowiskoWykonania string `json:"executionEnv,omitempty"`
}

// ProwenancjaZapytania składa prowenancję z tego, co niesie zapytanie i definicja
// kanału. Wypełnia wspólny szkielet nazwami pól kanału głównego, żeby klient
// odczytał prowenancję z każdego kanału jednakowo. Wywołujący dokłada jedynie to,
// co zna sam adapter: argv albo adres.
func ProwenancjaZapytania(z Zapytanie, d Definicja) Prowenancja {
	return Prowenancja{
		Program:             d.KluczAdaptera(),
		PromptSystemowy:     z.Nakladka.PromptSystemowy(),
		TrybNakladki:        trybNakladki(z.Nakladka),
		WarstwyNakladki:     nazwyWarstw(z.Nakladka),
		SkrotNakladki:       z.Nakladka.Skrot(),
		Model:               z.WybranyModel(d),
		ModelZapasowy:       z.ModelZapasowy,
		Naklad:              z.NakladRozumowania,
		TrybUprawnien:       string(z.TrybUprawnien),
		Katalogi:            z.KatalogiRobocze,
		KatalogRoboczy:      katalogRoboczy(z),
		PlikUstawien:        z.PlikUstawien,
		KonfiguracjaMCP:     konfiguracjaMCP(z),
		Wznowienie:          z.Wznowienie,
		Konto:               z.WybraneKonto(d),
		Proba:               1,
		Chwila:              time.Now().UTC(),
		Kanal:               d.Kod,
		Adapter:             d.KluczAdaptera(),
		Ustawienia:          z.Ustawienia,
		SrodowiskoWykonania: string(z.SrodowiskoWykonania),
	}
}

// trybNakladki normalizuje tryb nakładki do jednej z dwóch wartości szkieletu.
// Pusty tryb znaczy dopisanie.
func trybNakladki(n Nakladka) string {
	if n.Tryb == trybNakladkiZastap {
		return trybNakladkiZastap
	}
	return trybNakladkiDopisz
}

// nazwyWarstw wylicza nazwy niepustych warstw nakładki w kolejności
// krytyczności — do słownika warstw prowenancji. Warstwa pusta nie wchodzi do
// wykazu, tak samo jak nie wchodzi do promptu.
func nazwyWarstw(n Nakladka) []string {
	nazwy := make([]string, 0, 3)
	pary := []struct {
		nazwa string
		tresc string
	}{
		{warstwaKonstytucja, n.Konstytucja},
		{warstwaProfil, n.ProfilRoli},
		{warstwaEkspertyza, n.Ekspertyza},
	}
	for _, para := range pary {
		if strings.TrimSpace(para.tresc) != "" {
			nazwy = append(nazwy, para.nazwa)
		}
	}
	return nazwy
}

// katalogRoboczy rozstrzyga katalog, w którym model zostawia własne pliki:
// wygrywa katalog sesji, a w jego braku pierwszy z katalogów roboczych. To
// odpowiednik pola workingDirectory kanału głównego.
func katalogRoboczy(z Zapytanie) string {
	if strings.TrimSpace(z.KatalogSesji) != "" {
		return z.KatalogSesji
	}
	if len(z.KatalogiRobocze) > 0 {
		return z.KatalogiRobocze[0]
	}
	return ""
}

// konfiguracjaMCP składa wykaz konfiguracji mostów MCP w jeden wycinek: najpierw
// konfiguracja okna, potem konfiguracje wyliczone z obszaru mcp sesji. To
// odpowiednik pola mcpConfig kanału głównego.
func konfiguracjaMCP(z Zapytanie) []string {
	komplet := make([]string, 0, 1+len(z.DodatkoweMCP))
	if strings.TrimSpace(z.KonfiguracjaMCP) != "" {
		komplet = append(komplet, z.KonfiguracjaMCP)
	}
	komplet = append(komplet, z.DodatkoweMCP...)
	if len(komplet) == 0 {
		return nil
	}
	return komplet
}

// SkrotTresci liczy skrót diagnostyczny z podanych części. Części puste pomija,
// aby dopisanie pustej warstwy nakładki nie zmieniało skrótu.
func SkrotTresci(czesci ...string) string {
	niepuste := make([]string, 0, len(czesci))
	for _, czesc := range czesci {
		if strings.TrimSpace(czesc) != "" {
			niepuste = append(niepuste, czesc)
		}
	}
	if len(niepuste) == 0 {
		return ""
	}
	suma := sha256.Sum256([]byte(strings.Join(niepuste, "\x00")))
	return hex.EncodeToString(suma[:])
}
