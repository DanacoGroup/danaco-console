package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Zapora rodzin `studio.pdf.*` i `studio.security.*`: żadnego programu z zewnątrz.
//
// ── Powód, żeby następny wykonawca nie musiał go odtwarzać z rozmowy ────────
// Warsztat dokumentu stał raz na `qpdf` i na Ghostscripcie. Sprawdziany świeciły
// zielono, bo maszyna deweloperska ma oba doinstalowane ręcznie — a instalka
// produktu ich nie niesie i nieść nie ma. To nie przeoczenie, to decyzja
// Właściciela. U Operatora każda z tych czynności odmawiałaby za każdym razem,
// a rdzeń meldowałby, że komenda jest obsłużona.
//
// Dlatego te dwie rodziny pracują wyłącznie bibliotekami wkompilowanymi
// w binarium: `pdfcpu` dla dokumentu, `crypto` ze standardowej biblioteki dla
// podpisu i szyfrowania. Zapora patrzy na pliki tych rodzin i nie przepuszcza
// ani `zewnetrzne.Wolaj`, ani `exec.Command`.
//
// Zapora działa w obie strony: gdy plik rodziny zniknie albo zmieni nazwę,
// kończy się niepowodzeniem zamiast cicho przestać czegokolwiek pilnować.

// plikiBezProcesow wymienia pliki, których ta zasada dotyczy.
var plikiBezProcesow = []string{
	"adapter_studio_pdf.go",
	"adapter_studio_bezpieczenstwo.go",
}

// wywolaniaZabronione wymienia sposoby uruchomienia cudzego programu.
var wywolaniaZabronione = []string{
	"zewnetrzne.Wolaj",
	"exec.Command",
	"exec.CommandContext",
	"os/exec",
}

// TestWarsztatDokumentuNieUruchamiaProcesow pilnuje zasady bezwzględnej.
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
