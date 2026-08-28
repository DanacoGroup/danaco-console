// Odnajdywanie ma dawać pierwszeństwo pakietowi produktu przed ścieżką
// systemu, a sprawdziany tego pliku mierzą właśnie to rozstrzygnięcie.
package zewnetrzne

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// programPakietu zakłada w katalogu tymczasowym pozorny pakiet produktu
// z jednym programem w środku i przestawia na niego katalog bieżący.
func programPakietu(t *testing.T, nazwa string) string {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("sprawdzian ustawia prawo wykonania, którego Windows nie zna w tej postaci")
	}
	korzen := t.TempDir()
	katalog := filepath.Join(korzen, katalogPomocnikow)
	if err := os.MkdirAll(katalog, 0o755); err != nil {
		t.Fatalf("nie można założyć katalogu pomocników: %v", err)
	}
	sciezka := filepath.Join(katalog, nazwa)
	if err := os.WriteFile(sciezka, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("nie można zapisać programu pozornego: %v", err)
	}
	t.Chdir(korzen)
	return sciezka
}

// TestOdnajdywanieWidziProgramZPakietu wykazuje drogę, której wcześniej nie
// było: program leżący w pakiecie produktu, poza ścieżką wyszukiwania systemu.
func TestOdnajdywanieWidziProgramZPakietu(t *testing.T) {
	nazwa := "narzedzie-tylko-w-pakiecie"
	oczekiwana := programPakietu(t, nazwa)

	sciezka, jest := Odnajdz(Narzedzie{Nazwa: "Narzędzie pakietu", Program: nazwa})
	if !jest {
		t.Fatalf("program z pakietu nie został odnaleziony; szukano jako %q", nazwa)
	}
	if sciezka != oczekiwana {
		t.Fatalf("odnaleziono %q, a program pakietu leży w %q", sciezka, oczekiwana)
	}
	if !Stoi(Narzedzie{Program: nazwa}) {
		t.Fatal("sonda obecności nie widzi programu, który odnajdywanie właśnie wskazało")
	}
}

// TestPakietMaPierwszenstwoPrzedSciezkaSystemu pilnuje kolejności. Wersja
// dołożona do pakietu jest tą, którą sprawdzono przed wydaniem; wersja zastana
// na maszynie bywa starsza albo okrojona i nie jest niczyją obietnicą.
func TestPakietMaPierwszenstwoPrzedSciezkaSystemu(t *testing.T) {
	// `sh` stoi na ścieżce systemu każdej maszyny, na której się uruchomi.
	oczekiwana := programPakietu(t, "sh")

	sciezka, jest := Odnajdz(Narzedzie{Nazwa: "Powłoka", Program: "sh"})
	if !jest {
		t.Fatal("powłoka nie została odnaleziona ani w pakiecie, ani na ścieżce")
	}
	if sciezka != oczekiwana {
		t.Fatalf("odnaleziono %q, a pierwszeństwo ma mieć pakiet produktu (%q)",
			sciezka, oczekiwana)
	}
}

// TestKatalogNieUdajeProgramu pilnuje, żeby wpis o nazwie programu, który jest
// katalogiem, nie został uznany za program. Uruchomienie skończyłoby się
// odmową systemu w miejscu, w którym rdzeń zapowiedział już powodzenie.
func TestKatalogNieUdajeProgramu(t *testing.T) {
	korzen := t.TempDir()
	if err := os.MkdirAll(filepath.Join(korzen, katalogPomocnikow, "pozorny"), 0o755); err != nil {
		t.Fatalf("nie można założyć katalogu pozornego: %v", err)
	}
	t.Chdir(korzen)

	if _, jest := Odnajdz(Narzedzie{Nazwa: "Pozorny", Program: "pozorny"}); jest {
		t.Fatal("katalog o nazwie programu został uznany za program")
	}
}

// TestSciezkaBezwzglednaNiePodlegaSzukaniu pilnuje wskazania wprost: gdy montaż
// albo Operator podał ścieżkę bezwzględną, rdzeń ma jej użyć, a nie szukać
// programu o tej nazwie gdzie indziej.
func TestSciezkaBezwzglednaNiePodlegaSzukaniu(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sprawdzian ustawia prawo wykonania, którego Windows nie zna w tej postaci")
	}
	korzen := t.TempDir()
	wskazany := filepath.Join(korzen, "wskazany-wprost")
	if err := os.WriteFile(wskazany, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("nie można zapisać programu wskazanego: %v", err)
	}

	sciezka, jest := Odnajdz(Narzedzie{Nazwa: "Wskazany", Program: wskazany})
	if !jest || sciezka != wskazany {
		t.Fatalf("wskazanie bezwzględne nie zostało użyte wprost: %q, obecny=%v", sciezka, jest)
	}

	brakujacy := filepath.Join(korzen, "nie-ma-takiego-pliku")
	if _, jest := Odnajdz(Narzedzie{Nazwa: "Brak", Program: brakujacy}); jest {
		t.Fatal("wskazanie bezwzględne na plik nieistniejący zameldowano jako obecne")
	}
}
