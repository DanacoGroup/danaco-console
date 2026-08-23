package injection

import (
	"encoding/json"
	"time"

	"danacoconsole/shared"
)

// Prowenancja opisuje, co dokładnie poszło do modelu: wiersz argumentów, prompt
// systemowy, ustawienia i skrót nakładki. Jest wyłącznie zapisem warunków
// wywołania — niczego nie dopuszcza ani nie wstrzymuje.
type Prowenancja struct {
	Program         string   `json:"program"`
	Argv            []string `json:"argv"`
	PromptSystemowy string   `json:"systemPrompt"`
	TrybNakladki    string   `json:"overlayMode"`
	WarstwyNakladki []string `json:"overlayLayers"`
	SkrotNakladki   string   `json:"overlayHash"`
	Model           string   `json:"model"`
	ModelZapasowy   string   `json:"fallbackModel"`
	Naklad          string   `json:"effort"`
	TrybUprawnien   string   `json:"permissionMode"`
	Katalogi        []string `json:"directories"`
	KatalogRoboczy  string   `json:"workingDirectory"`
	PlikUstawien    string   `json:"settings"`
	KonfiguracjaMCP []string `json:"mcpConfig"`
	Wznowienie      string   `json:"resume"`
	Konto           string   `json:"account"`
	KatalogKonta    string   `json:"accountConfigDir"`
	// ZrodloKonta nazywa drogę, którą tożsamość weszła do wywołania: wskazanie
	// konfiguracji, kolejność puli, rotacja po wyczerpaniu albo tożsamość
	// otoczenia.
	ZrodloKonta string    `json:"accountOrigin,omitempty"`
	Proba       int       `json:"attempt"`
	PowodProby  string    `json:"attemptReason,omitempty"`
	Chwila      time.Time `json:"at"`
}

// ZlozProwenancje zbiera prowenancję jednej próby wywołania. Próba druga
// i dalsze wynikają z rotacji konta; powód zmiany niesie pole PowodProby.
func ZlozProwenancje(u Ustawienia, n Nakladka, argv []string, konto Konto, proba int, powod string) Prowenancja {
	return Prowenancja{
		Program:         u.Program,
		Argv:            argv,
		PromptSystemowy: n.Prompt(),
		TrybNakladki:    trybNakladki(n),
		WarstwyNakladki: n.NazwyWarstw(),
		SkrotNakladki:   n.Skrot(),
		Model:           u.Model,
		ModelZapasowy:   u.ModelZapasowy,
		Naklad:          u.Naklad,
		TrybUprawnien:   string(u.TrybUprawnien),
		Katalogi:        u.Katalogi,
		KatalogRoboczy:  u.KatalogRoboczy,
		PlikUstawien:    u.PlikUstawien,
		KonfiguracjaMCP: u.KonfiguracjaMCP,
		Wznowienie:      u.Wznowienie,
		Konto:           konto.Kod,
		KatalogKonta:    konto.KatalogKonfiguracji,
		ZrodloKonta:     zrodloKonta(u, konto, proba),
		Proba:           proba,
		PowodProby:      powod,
		Chwila:          time.Now().UTC(),
	}
}

// zrodloKonta nazywa drogę, którą tożsamość weszła do wywołania.
func zrodloKonta(u Ustawienia, konto Konto, proba int) string {
	switch {
	case u.Konto != "":
		return "wskazanie konfiguracji"
	case konto.Kod == "":
		return "tożsamość otoczenia"
	case proba > 1:
		return "rotacja po wyczerpaniu"
	default:
		return "kolejność puli"
	}
}

func trybNakladki(n Nakladka) string {
	if n.Tryb == TrybZastap {
		return TrybZastap
	}
	return TrybDopisz
}

// FragmentProwenancji pakuje prowenancję we fragment strumienia. Kanał wysyła
// go przed jakimkolwiek fragmentem tekstu, żeby zapis rozmowy niósł najpierw
// warunki wywołania, a dopiero potem jego wynik.
func FragmentProwenancji(okno, wiadomosc string, p Prowenancja) Fragment {
	fragment := shared.StreamChunkEvent{
		WindowId:  okno,
		MessageId: wiadomosc,
		Kind:      RodzajProwenancja,
	}
	if surowe, err := json.Marshal(p); err == nil {
		fragment.Data = surowe
	}
	return Fragment{Chunk: fragment}
}
