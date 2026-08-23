package shared_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Świeżość generatu kontraktu.
//
// `contract.json` jest źródłem prawdy, ale w kompilacji nie bierze udziału:
// bierze w niej udział `contract.go`, a w budowaniu klienta `contract.ts`.
// Zmiana źródła bez puszczenia generatora rozjeżdża je po cichu i obie strony
// kompilują się dalej — rdzeń zna nazwę, której klient nie zna, albo odwrotnie.
// Rozjazd wychodzi dopiero na gnieździe, u Operatora.
//
// Sprawdzian puszcza generator na kopii i porównuje wynik z tym, co leży
// w drzewie. Kopia, a nie katalog źródłowy: generator zapisuje artefakty na
// dysk, więc puszczony na miejscu nadpisałby pliki sprawdzanego drzewa.

// TestGeneratOdpowiadaZrodluKontraktu porównuje wygenerowane bindingi z tymi,
// które leżą w repozytorium.
func TestGeneratOdpowiadaZrodluKontraktu(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("brak node w PATH — sprawdzian świeżości generatu wymaga uruchamiacza generatora")
	}

	katalogRoboczy := t.TempDir()
	skopiujPlik(t, "contract.json", filepath.Join(katalogRoboczy, "contract.json"))
	skopiujKatalog(t, "gen", filepath.Join(katalogRoboczy, "gen"))

	polecenie := exec.Command(node, filepath.Join(katalogRoboczy, "gen", "generate.mjs"))
	polecenie.Dir = katalogRoboczy
	if wyjscie, err := polecenie.CombinedOutput(); err != nil {
		t.Fatalf("generator kontraktu nie zakończył pracy: %v\n%s", err, wyjscie)
	}

	for _, artefakt := range []string{"contract.go", "contract.ts"} {
		t.Run(artefakt, func(t *testing.T) {
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

// wczytaj czyta plik i przerywa sprawdzian, gdy pliku nie ma.
func wczytaj(t *testing.T, sciezka string) []byte {
	t.Helper()
	tresc, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać %s: %v", sciezka, err)
	}
	return tresc
}

// skopiujPlik przenosi jeden plik do katalogu roboczego sprawdzianu.
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
