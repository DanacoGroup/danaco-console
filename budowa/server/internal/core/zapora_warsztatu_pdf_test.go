package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Zapora studio.pdf.* i studio.security.*: bez programu zewnętrznego, tylko biblioteki wkompilowane.

// plikiBezProcesow wymienia pliki, których dotyczy zakaz wywoływania jakiegokolwiek programu zewnętrznego.
var plikiBezProcesow = []string{
	"adapter_studio_pdf.go",
	"adapter_studio_bezpieczenstwo.go",
}

// wywolaniaZabronione wymienia sposoby uruchomienia cudzego programu, których żaden plik nie może zawierać.
var wywolaniaZabronione = []string{
	"zewnetrzne.Wolaj",
	"exec.Command",
	"exec.CommandContext",
	"os/exec",
}

// TestWarsztatDokumentuNieUruchamiaProcesow pilnuje zasady bezwzględnej: żadnego wywołania programu zewnętrznego.
func TestWarsztatDokumentuNieUruchamiaProcesow(t *testing.T) {
	for _, nazwa := range plikiBezProcesow {
		sciezka := filepath.Join(".", nazwa)
		tresc, err := os.ReadFile(sciezka)
		if err != nil {
			t.Fatalf("pliku %s nie ma w drzewie — zapora przestałaby czegokolwiek "+
				"pilnować: %v", nazwa, err)
		}
		for _, wywolanie := range wywolaniaZabronione {
			if strings.Contains(string(tresc), wywolanie) {
				t.Fatalf("plik %s sięga po %s; rodziny studio.pdf.* i studio.security.* "+
					"pracują wyłącznie biblioteką wkompilowaną w rdzeń, bo instalka nie "+
					"niesie programów zewnętrznych dla dokumentu i kryptografii",
					nazwa, wywolanie)
			}
		}
	}
}

// TestWarsztatDokumentuNieWymieniaZdjetychProgramow pilnuje, żeby `qpdf`
// i Ghostscript nie wróciły do rdzenia żadną drogą — także przez sondę
// zależności startowych, która nazwałaby je zależnością produktu.
func TestWarsztatDokumentuNieWymieniaZdjetychProgramow(t *testing.T) {
	wpisy, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("nie można przejrzeć rdzenia: %v", err)
	}
	zdjete := []string{"qpdf", "ghostscript"}
	for _, wpis := range wpisy {
		nazwa := wpis.Name()
		if wpis.IsDir() || !strings.HasSuffix(nazwa, ".go") ||
			strings.HasSuffix(nazwa, "_test.go") {
			continue
		}
		tresc, err := os.ReadFile(nazwa)
		if err != nil {
			t.Fatalf("nie można odczytać %s: %v", nazwa, err)
		}
		for _, program := range zdjete {
			if strings.Contains(strings.ToLower(string(tresc)), program) {
				t.Fatalf("plik %s wymienia program %s; został zdjęty z rdzenia, bo "+
					"instalka go nie niesie", nazwa, program)
			}
		}
	}
}
