package injection

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// Nazwy warstw nakładki uporządkowane wg krytyczności:
// konstytucja stoi najwyżej, ekspertyza zadaniowa najniżej.
const (
	WarstwaKonstytucja = "konstytucja"
	WarstwaProfil      = "profil"
	WarstwaEkspertyza  = "ekspertyza"
)

// Tryby podania nakładki programowi `claude` — dopisanie do promptu własnego
// programu albo jego pełne zastąpienie.
const (
	// TrybDopisz dokłada nakładkę do promptu własnego programu
	// (--append-system-prompt). Narzędzia i zachowanie powłoki zostają.
	TrybDopisz = "dopisz"
	// TrybZastap podmienia prompt programu w całości, przełącznikiem
	// --system-prompt, tracąc prompt oryginalny.
	TrybZastap = "zastap"
)

// kolejnoscWarstw daje pierwszeństwo warstwom nazwanym; warstwa o nazwie
// spoza katalogu ląduje za nimi, w kolejności podania.
var kolejnoscWarstw = map[string]int{
	WarstwaKonstytucja: 1,
	WarstwaProfil:      2,
	WarstwaEkspertyza:  3,
}

// Warstwa jest jednym poziomem nakładki. Treść pochodzi wyłącznie
// z konfiguracji — pakiet nie zna ani jednego zdania promptu.
type Warstwa struct {
	Nazwa string
	Tresc string
}

// Nakladka jest złożeniem warstw promptu systemowego w jeden gotowy do
// podania programowi prompt tekstowy.
type Nakladka struct {
	// Tryb rozstrzyga, czy nakładka dopisuje się do promptu, czy go zastępuje.
	Tryb string
	// Warstwy w dowolnej kolejności podania; Prompt układa je wg krytyczności.
	Warstwy []Warstwa
}

// NowaNakladka składa nakładkę, pomijając warstwy puste. Brak warstwy nie jest
// błędem i niczego nie wstrzymuje — kanał ruszy z samą powłoką.
func NowaNakladka(tryb string, warstwy ...Warstwa) Nakladka {
	n := Nakladka{Tryb: tryb}
	for _, w := range warstwy {
		if strings.TrimSpace(w.Tresc) == "" {
			continue
		}
		n.Warstwy = append(n.Warstwy, w)
	}
	return n
}

// Uporzadkowane zwraca warstwy ułożone wg krytyczności. Sortowanie
// jest stabilne, więc warstwy nienazwane zachowują kolejność podania.
func (n Nakladka) Uporzadkowane() []Warstwa {
	w := make([]Warstwa, len(n.Warstwy))
	copy(w, n.Warstwy)
	sort.SliceStable(w, func(i, j int) bool {
		return pozycjaWarstwy(w[i].Nazwa) < pozycjaWarstwy(w[j].Nazwa)
	})
	return w
}

func pozycjaWarstwy(nazwa string) int {
	if p, jest := kolejnoscWarstw[nazwa]; jest {
		return p
	}
	return len(kolejnoscWarstw) + 1
}

// Prompt skleja warstwy w jeden prompt systemowy. Pusta nakładka daje pusty
// napis, a wtedy kanał nie podaje żadnego przełącznika promptu.
func (n Nakladka) Prompt() string {
	uporzadkowane := n.Uporzadkowane()
	czesci := make([]string, 0, len(uporzadkowane))
	for _, w := range uporzadkowane {
		czesci = append(czesci, strings.TrimSpace(w.Tresc))
	}
	return strings.Join(czesci, "\n\n")
}

// NazwyWarstw zwraca wykaz nazw warstw nakładki w kolejności ich złożenia,
// przeznaczony do prowenancji.
func (n Nakladka) NazwyWarstw() []string {
	uporzadkowane := n.Uporzadkowane()
	nazwy := make([]string, 0, len(uporzadkowane))
	for _, w := range uporzadkowane {
		nazwy = append(nazwy, w.Nazwa)
	}
	return nazwy
}

// Skrot jest skrótem SHA-256 złożonego promptu. Służy wyłącznie diagnostyce
// prowenancji: niczego nie dopuszcza i niczego nie blokuje.
func (n Nakladka) Skrot() string {
	prompt := n.Prompt()
	if prompt == "" {
		return ""
	}
	suma := sha256.Sum256([]byte(prompt))
	return hex.EncodeToString(suma[:])
}

// Przelacznik zwraca przełącznik programu właściwy trybowi nakładki oraz
// informację, czy w ogóle jest co podawać.
func (n Nakladka) Przelacznik() (string, string, bool) {
	prompt := n.Prompt()
	if prompt == "" {
		return "", "", false
	}
	if n.Tryb == TrybZastap {
		return "--system-prompt", prompt, true
	}
	return "--append-system-prompt", prompt, true
}
