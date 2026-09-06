package shared_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Świeżość generatu kontraktu; sprawdzian puszcza generator na kopii i porównuje wynik z drzewem.

// zrodloKontraktu odwzorowuje wyłącznie listę artefaktów z contract.json.
type zrodloKontraktu struct {
	Kontrakt struct {
		Artefakty struct {
			TypeScript     string `json:"typescript"`
			Go             string `json:"go"`
			RejestrKlienta string `json:"rejestrKlienta"`
		} `json:"artefakty"`
	} `json:"kontrakt"`
}

// TestGeneratOdpowiadaZrodluKontraktu porównuje wygenerowane bindingi z tymi,
// które leżą w repozytorium.
func TestGeneratOdpowiadaZrodluKontraktu(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("brak node w PATH — sprawdzian świeżości generatu wymaga uruchamiacza generatora")
	}

	// Generator zapisuje część artefaktów poza shared/, więc kopia musi odwzorować
	// układ katalogów repozytorium, a nie sam katalog shared/.
	korzen := t.TempDir()
	katalogRoboczy := filepath.Join(korzen, "shared")
	if err := os.MkdirAll(katalogRoboczy, 0o755); err != nil {
		t.Fatalf("nie można założyć %s: %v", katalogRoboczy, err)
	}
	skopiujPlik(t, "contract.json", filepath.Join(katalogRoboczy, "contract.json"))
	skopiujKatalog(t, "gen", filepath.Join(katalogRoboczy, "gen"))

	artefakty := artefaktyKontraktu(t)
	for _, artefakt := range artefakty {
		katalog := filepath.Dir(filepath.Join(katalogRoboczy, artefakt))
		if err := os.MkdirAll(katalog, 0o755); err != nil {
			t.Fatalf("nie można założyć %s: %v", katalog, err)
		}
	}

	polecenie := exec.Command(node, filepath.Join(katalogRoboczy, "gen", "generate.mjs"))
	polecenie.Dir = katalogRoboczy
	if wyjscie, err := polecenie.CombinedOutput(); err != nil {
		t.Fatalf("generator kontraktu nie zakończył pracy: %v\n%s", err, wyjscie)
	}

	for _, artefakt := range artefakty {
		t.Run(filepath.Base(artefakt), func(t *testing.T) {
			wDrzewie := wczytaj(t, artefakt)
			swiezy := wczytaj(t, filepath.Join(katalogRoboczy, artefakt))
			if string(wDrzewie) == string(swiezy) {
				return
			}
			t.Errorf("%s rozjechał się z contract.json (%d B w drzewie, %d B po wygenerowaniu). "+
				"Puść `node shared/gen/generate.mjs` i zatwierdź wynik razem ze zmianą kontraktu.",
				artefakt, len(wDrzewie), len(swiezy))
		})
	}
}

// artefaktyKontraktu zwraca ścieżki wytworów zadeklarowane w źródle prawdy;
// dopisanie kolejnego artefaktu w contract.json obejmuje sprawdzian bez zmiany kodu.
func artefaktyKontraktu(t *testing.T) []string {
	t.Helper()
	var zrodlo zrodloKontraktu
	if err := json.Unmarshal(wczytaj(t, "contract.json"), &zrodlo); err != nil {
		t.Fatalf("nie można odczytać listy artefaktów z contract.json: %v", err)
	}
	sciezki := []string{zrodlo.Kontrakt.Artefakty.TypeScript, zrodlo.Kontrakt.Artefakty.Go, zrodlo.Kontrakt.Artefakty.RejestrKlienta}
	for _, sciezka := range sciezki {
		if sciezka == "" {
			t.Fatal("contract.json nie deklaruje kompletu artefaktów — sprawdzian nie ma czego porównać")
		}
	}
	return sciezki
}

// Funkcja wczytaj czyta wskazany plik z dysku i przerywa cały sprawdzian, gdy tego pliku tam wciąż nie ma.
func wczytaj(t *testing.T, sciezka string) []byte {
	t.Helper()
	tresc, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać %s: %v", sciezka, err)
	}
	return tresc
}

// Funkcja skopiujPlik przenosi jeden wskazany plik do katalogu roboczego tego samego uruchomionego sprawdzianu.
func skopiujPlik(t *testing.T, zrodlo, cel string) {
	t.Helper()
	tresc := wczytaj(t, zrodlo)
	if err := os.WriteFile(cel, tresc, 0o644); err != nil {
		t.Fatalf("nie można zapisać %s: %v", cel, err)
	}
}

// skopiujKatalog przenosi płaski katalog generatora. Głębszego drzewa nie
// obsługuje z zamysłu — `gen/` jest płaski i ma taki zostać.
func skopiujKatalog(t *testing.T, zrodlo, cel string) {
	t.Helper()
	if err := os.MkdirAll(cel, 0o755); err != nil {
		t.Fatalf("nie można założyć %s: %v", cel, err)
	}
	wpisy, err := os.ReadDir(zrodlo)
	if err != nil {
		t.Fatalf("nie można odczytać katalogu %s: %v", zrodlo, err)
	}
	for _, wpis := range wpisy {
		if wpis.IsDir() {
			t.Fatalf("katalog %s zawiera podkatalog %q — sprawdzian kopiuje wyłącznie płaskie drzewo",
				zrodlo, wpis.Name())
		}
		skopiujPlik(t, filepath.Join(zrodlo, wpis.Name()), filepath.Join(cel, wpis.Name()))
	}
}
