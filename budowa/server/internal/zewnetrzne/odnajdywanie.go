// Odpowiedzialność pliku: jedna droga odnajdywania binarium arsenału —
// najpierw w pakiecie produktu, dopiero potem na ścieżce wyszukiwania systemu.
package zewnetrzne

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// katalogPomocnikow to nazwa katalogu, do którego pakowanie wnosi programy
// towarzyszące. Ta sama nazwa stoi w `tauri.conf.json` jako cel zasobów.
const katalogPomocnikow = "pomocniki"

// Odnajdz wskazuje plik wykonywalny narzędzia oraz mówi, czy w ogóle go widać.
//
// Kolejność jest rozstrzygnięciem, nie wygodą: wskazanie bezwzględne Operatora,
// potem pakiet produktu, na końcu ścieżka systemu.
func Odnajdz(n Narzedzie) (string, bool) {
	program := strings.TrimSpace(n.Program)
	if program == "" {
		return "", false
	}

	// Ścieżka wskazana wprost nie podlega szukaniu, tylko sprawdzeniu.
	if filepath.IsAbs(program) {
		if wykonywalny(program) {
			return program, true
		}
		return program, false
	}

	for _, miejsce := range miejscaPakietu(program) {
		if wykonywalny(miejsce) {
			return miejsce, true
		}
	}

	if zeSciezki, err := exec.LookPath(program); err == nil {
		return zeSciezki, true
	}
	return program, false
}

// miejscaPakietu składa ścieżki, pod którymi program mógłby leżeć w pakiecie.
// Rozszerzenie `.exe` dokłada się na Windowsie. Płaski układ idzie pierwszy,
// własny katalog jest drogą zapasową.
func miejscaPakietu(program string) []string {
	nazwy := []string{program}
	if runtime.GOOS == "windows" && !strings.EqualFold(filepath.Ext(program), ".exe") {
		nazwy = append(nazwy, program+".exe")
	}

	var korzenie []string
	if biezace, err := os.Executable(); err == nil {
		korzenie = append(korzenie, filepath.Dir(biezace))
	}
	if katalog, err := os.Getwd(); err == nil {
		korzenie = append(korzenie, katalog)
	}

	// Nazwa katalogu własnego jest nazwą programu bez rozszerzenia.
	katalogWlasny := strings.TrimSuffix(program, filepath.Ext(program))
	if !strings.EqualFold(filepath.Ext(program), ".exe") {
		katalogWlasny = program
	}

	var miejsca []string
	for _, korzen := range korzenie {
		for _, nazwa := range nazwy {
			miejsca = append(miejsca, filepath.Join(korzen, katalogPomocnikow, nazwa))
		}
		for _, nazwa := range nazwy {
			miejsca = append(miejsca, filepath.Join(korzen, katalogPomocnikow, katalogWlasny, nazwa))
		}
	}
	return miejsca
}

// wykonywalny odpowiada, czy pod ścieżką stoi plik zwykły, który system zgodzi
// się uruchomić. Katalog o nazwie programu nie jest programem.
func wykonywalny(sciezka string) bool {
	opis, err := os.Stat(sciezka)
	if err != nil || opis.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}
	return opis.Mode().Perm()&0o111 != 0
}
