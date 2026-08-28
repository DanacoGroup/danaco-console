package core

import (
	"os"
	"strings"
	"testing"
)

// Zapora jest napisana grubo — na wystąpienie nazwy w pliku, bez rozróżniania kodu od komentarza.

// silnikiSpozaInstalki wylicza nazwy programów, których w obszarze Design być nie może; biblioteki Go wkompilowane w binarium są dozwolone wszędzie.
var silnikiSpozaInstalki = []string{
	"vips",
	"imagemagick",
	"potrace",
	"inkscape",
	"fontforge",
	"fonttools",
	"graphicsmagick",
}

// wolaniaProcesu wylicza drogi uruchomienia procesu z pliku rdzenia, sprawdzane przez tę zaporę wprost.
var wolaniaProcesu = []string{
	"exec.Command",
	"exec.CommandContext",
	"zewnetrzne.Wolaj",
	"os/exec",
}

// wolaniaWlasneProcesu to drogi, które zapora zabrania zawsze i wszędzie: uruchomienie procesu poza pakietem zewnetrzne pomija sprawdzenie programu i izolację okna.
var wolaniaWlasneProcesu = []string{
	"exec.Command",
	"exec.CommandContext",
	"os/exec",
}

// plikiDesignuZOdczytemPisma to jedyny wyjątek tej zapory: Tesseract rozpoznaje pismo w design.mockup.import, bo czytnika liter w czystym Go nie ma, a biblioteka Go dla obrazu i wektora istnieje wszędzie indziej.
var plikiDesignuZOdczytemPisma = map[string]string{
	"adapter_modul_design_makiety_zrzut.go": "tesseract",
}

// TestWarsztatFotografiiNieWolaProcesu pilnuje, żeby żaden plik rodzin design.photo.* i design.print.* nie uruchamiał procesu; wykaz przedrostków obejmuje cały obszar Design.
func TestWarsztatFotografiiNieWolaProcesu(t *testing.T) {
	wpisy, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("nie można przejrzeć rdzenia: %v", err)
	}
	sprawdzonych := 0
	for _, wpis := range wpisy {
		nazwa := wpis.Name()
		if wpis.IsDir() || !strings.HasSuffix(nazwa, ".go") ||
			strings.HasSuffix(nazwa, "_test.go") {
			continue
		}
		if !strings.HasPrefix(nazwa, "adapter_modul_design") {
			continue
		}
		tresc, err := os.ReadFile(nazwa)
		if err != nil {
			t.Fatalf("nie można odczytać %s: %v", nazwa, err)
		}
		sprawdzonych++

		if program, wyjatek := plikiDesignuZOdczytemPisma[nazwa]; wyjatek {
			// Wyjątek nazwany: proces wolno wołać tylko drogą pakietu zewnetrzne, dla programu z wykazu.
			for _, wolanie := range wolaniaWlasneProcesu {
				if strings.Contains(string(tresc), wolanie) {
					t.Errorf("plik %s uruchamia proces własną drogą (%s); wyjątek na odczyt pisma "+
						"otwiera wyłącznie zewnetrzne.Wolaj, bo tylko ona sprawdza obecność programu, "+
						"nakłada bramę izolacji i obejmuje drzewo procesów", nazwa, wolanie)
				}
			}
			maly := strings.ToLower(string(tresc))
			if strings.Contains(maly, "zewnetrzne.wolaj") && !strings.Contains(maly, program) {
				t.Errorf("plik %s woła proces, ale nie wymienia programu %s; wyjątek tej zapory "+
					"jest nazwany jednym programem — inny program znaczy, że wyjątek posłużył "+
					"do czegoś, na co go nie dano", nazwa, program)
			}
			continue
		}

		for _, wolanie := range wolaniaProcesu {
			if strings.Contains(string(tresc), wolanie) {
				t.Errorf("plik %s woła proces (%s); obróbka obrazu i wydania drukarskie idą "+
					"WYŁĄCZNIE bibliotekami wkompilowanymi w binarium — u Operatora stoi cienka "+
					"instalka i program spoza niej jest tam odmową, nie funkcją", nazwa, wolanie)
			}
		}
	}
	// Zapora, która nie przejrzała ani jednego pliku, jest zaporą, która nie broni niczego.
	if sprawdzonych == 0 {
		t.Fatal("zapora nie znalazła ani jednego pliku obszaru Design; sprawdź, czy nazwy " +
			"plików nie zmieniły przedrostka — inaczej ta zapora przestała czegokolwiek pilnować")
	}
}

// TestObszarDesignNieWymieniaSilnikowObrazuSpozaInstalki pilnuje, żeby nazwy silników obrazu i konturów nie pojawiły się w plikach obszaru Design żadną drogą, także przez komentarz.
func TestObszarDesignNieWymieniaSilnikowObrazuSpozaInstalki(t *testing.T) {
	wpisy, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("nie można przejrzeć rdzenia: %v", err)
	}
	sprawdzonych := 0
	for _, wpis := range wpisy {
		nazwa := wpis.Name()
		if wpis.IsDir() || !strings.HasSuffix(nazwa, ".go") ||
			strings.HasSuffix(nazwa, "_test.go") {
			continue
		}
		if !strings.HasPrefix(nazwa, "adapter_modul_design") {
			continue
		}
		tresc, err := os.ReadFile(nazwa)
		if err != nil {
			t.Fatalf("nie można odczytać %s: %v", nazwa, err)
		}
		sprawdzonych++

		maly := strings.ToLower(string(tresc))
		for _, silnik := range silnikiSpozaInstalki {
			if strings.Contains(maly, silnik) {
				t.Errorf("plik %s wymienia silnik %s; instalka Operatora go nie niesie, więc "+
					"funkcja od niego zależna byłaby u Operatora odmową — a nazwa w treści pliku "+
					"jest wskazówką ku tej szkodzie", nazwa, silnik)
			}
		}
	}
	// Zapora, która nie przejrzała ani jednego pliku, jest zaporą, która nie broni niczego.
	if sprawdzonych == 0 {
		t.Fatal("zapora nie znalazła ani jednego pliku obszaru Design; sprawdź, czy nazwy " +
			"plików nie zmieniły przedrostka — inaczej ta zapora przestała czegokolwiek pilnować")
	}
}
