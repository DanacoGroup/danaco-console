// Plik sprawdza wykaz wypisany maszynowo, jedyną drogę, którą prowizjonowanie
// serwera dowiaduje się, co postawić: wykaz okrojony, warstwę nieznaną
// i pozycję wstrzymaną wchodzącą do warstwy stawianej milcząco.
package core

import (
	"bytes"
	"strings"
	"testing"

	"danacoconsole/server/internal/mowa"
	"danacoconsole/server/internal/zewnetrzne"
)

// TestWydrukWykazuNiesieKazdaPozycje pilnuje kompletności: wydruk ma mieć
// dokładnie tyle wierszy danych, ile wykaz ma pozycji. Wiersz pominięty to
// program, którego prowizjonowanie nie postawi.
func TestWydrukWykazuNiesieKazdaPozycje(t *testing.T) {
	var bufor bytes.Buffer
	if err := WypiszWykazZaleznosci(&bufor); err != nil {
		t.Fatalf("wypisanie wykazu odmówiło: %v", err)
	}
	wykaz := ZaleznosciZewnetrzne()
	var dane []string
	for _, wiersz := range strings.Split(strings.TrimRight(bufor.String(), "\n"), "\n") {
		if !strings.HasPrefix(wiersz, "#") {
			dane = append(dane, wiersz)
		}
	}
	if len(dane) != len(wykaz) {
		t.Fatalf("wydruk niesie %d wierszy danych, a wykaz ma %d pozycji",
			len(dane), len(wykaz))
	}
	for _, wiersz := range dane {
		pola := strings.Split(wiersz, "\t")
		if len(pola) != 6 {
			t.Errorf("wiersz %q ma %d pól, a postać wykazu to sześć pól rozdzielonych tabulacją",
				wiersz, len(pola))
			continue
		}
		for numer, pole := range pola {
			// Puste pole jest dopuszczalne tylko w podpowiedzi pakietu (indeks 2).
			if numer != 2 && strings.TrimSpace(pole) == "" {
				t.Errorf("wiersz %q ma puste pole numer %d", wiersz, numer+1)
			}
		}
	}
}

// TestWydrukNiesiePolskiPakietTesseractu pilnuje, że rozpoznanie pisma działa
// po polsku: sam `tesseract-ocr` daje program bez języka polskiego, więc
// prowizjonowanie postawiłoby OCR, który na polskim skanie zwraca miał.
func TestWydrukNiesiePolskiPakietTesseractu(t *testing.T) {
	var bufor bytes.Buffer
	if err := WypiszWykazZaleznosci(&bufor); err != nil {
		t.Fatalf("wypisanie wykazu odmówiło: %v", err)
	}
	wiersz := wierszProgramu(bufor.String(), "tesseract")
	if wiersz == "" {
		t.Fatal("w wydruku nie ma Tesseractu — rozpoznanie pisma nie zostanie postawione")
	}
	if !strings.Contains(wiersz, "tesseract-ocr-pol") {
		t.Errorf("wiersz Tesseractu nie niesie polskiego pakietu językowego: %q", wiersz)
	}
	if !strings.HasPrefix(wiersz, WarstwaObowiazkowa+"\t") {
		t.Errorf("Tesseract nie stoi w warstwie obowiązkowej: %q", wiersz)
	}
}

// TestSilnikKontenerowStoiWWarstwieDecyzyjnej pilnuje, że wstrzymany silnik
// kontenerów nie wchodzi do warstwy, którą prowizjonowanie stawia bez pytania:
// obie jego deklaracje muszą trafić do warstwy decyzyjnej.
func TestSilnikKontenerowStoiWWarstwieDecyzyjnej(t *testing.T) {
	policzone := 0
	for _, pozycja := range ZaleznosciZewnetrzne() {
		if pozycja.Narzedzie.Program != "docker" && pozycja.Narzedzie.Program != "podman" {
			continue
		}
		policzone++
		if warstwa := WarstwaZaleznosci(pozycja); warstwa != WarstwaDecyzyjna {
			t.Errorf("pozycja %q trafiła do warstwy %q, a silnik kontenerów jest wstrzymany",
				pozycja.Narzedzie.Nazwa, warstwa)
		}
	}
	if policzone == 0 {
		t.Fatal("w wykazie nie ma silnika kontenerów — sprawdzian nie ma czego pilnować")
	}
}

// TestKazdaPozycjaMaZnanaWarstwe pilnuje umowy z jedynym odbiorcą wykazu:
// arsenal-serwera.sh rozdziela pozycje po nazwach warstw wpisanych na sztywno.
// Nazwy stoją tu literałem, nie stałą tego pliku, żeby nowa warstwa dodana
// w rdzeniu bez gałęzi w skrypcie zepsuła ten sprawdzian, a nie zniknęła.
func TestKazdaPozycjaMaZnanaWarstwe(t *testing.T) {
	znaneOdbiorcy := map[string]bool{
		"obowiazkowa-apt": true,
		"warsztat-go":     true,
		"warsztat-npm":    true,
		"snap":            true,
		"model-recznie":   true,
		"decyzyjna":       true,
	}
	for _, pozycja := range ZaleznosciZewnetrzne() {
		if warstwa := WarstwaZaleznosci(pozycja); !znaneOdbiorcy[warstwa] {
			t.Errorf("pozycja %q dostała warstwę %q, której skrypt prowizjonowania nie zna",
				pozycja.Narzedzie.Nazwa, warstwa)
		}
	}
}

// TestWarstwaRozpoznajePostaciPodpowiedzi mierzy samą regułę na deklaracjach
// ułożonych wprost — bez tego sprawdzian warstw mówiłby tylko o dzisiejszym
// wykazie, a nie o regule, która przyjmie następne narzędzie.
func TestWarstwaRozpoznajePostaciPodpowiedzi(t *testing.T) {
	przypadki := []struct {
		nazwa    string
		program  string
		pakiet   string
		oczekuje string
	}{
		{"pakiet dystrybucji", "pandoc", "pandoc", WarstwaObowiazkowa},
		{"dwa pakiety dystrybucji", "tesseract", "tesseract-ocr tesseract-ocr-pol", WarstwaObowiazkowa},
		{"warsztat Go", "gopls", "go install golang.org/x/tools/gopls@latest", WarstwaWarsztatGo},
		// Podpowiedź warsztatu Go potrafi nieść też adres GitHuba.
		{"warsztat Go z adresem GitHuba", "golangci-lint",
			"go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest", WarstwaWarsztatGo},
		{"warsztat npm", "eslint", "npm i -g eslint", WarstwaWarsztatNpm},
		{"snap", "kubectl", "kubectl (snap)", WarstwaSnap},
		{"wydanie z GitHuba", "realesrgan-ncnn-vulkan",
			"wydanie z github.com/xinntao/Real-ESRGAN/releases", WarstwaModelRecznie},
		{"środowisko pythonowe", "rembg", "rembg[cli] w osobnym środowisku pythonowym", WarstwaModelRecznie},
		{"silnik kontenerów", "docker", "docker.io albo podman", WarstwaDecyzyjna},
		{"silnik kontenerów pod podmanem", "podman", "podman", WarstwaDecyzyjna},
		// Podpowiedź zapisana zdaniem albo poleceniem innego menedżera pakietów zostaje
		// w warstwie obowiązkowej — arsenal-serwera.sh rozpoznaje jej postać i wyprowadza
		// z niej krok apt, krok pip albo krok ręczny, tak samo, jak dla nazw pakietów.
		{"podpowiedź zdaniem — środowisko Javy", "java",
			"środowisko uruchomieniowe Javy (default-jre) wraz z wydaniem Apache Tika w /opt/tika",
			WarstwaObowiazkowa},
		{"podpowiedź zdaniem — plik z wydania projektu", "typst",
			"typst (jeden plik wykonywalny z wydania projektu)", WarstwaObowiazkowa},
		{"podpowiedź zdaniem — pakiet wraz ze słownikiem", "hunspell",
			"hunspell wraz ze słownikiem języka (hunspell-pl, hunspell-en-us)", WarstwaObowiazkowa},
		{"polecenie pip install", "ruff", "pip install ruff", WarstwaObowiazkowa},
		// cargo install zawiera podnapis „go install" („cargo” kończy się na „go”),
		// a mimo to nie jest poleceniem Go.
		{"polecenie cargo install", "typos", "cargo install typos-cli", WarstwaObowiazkowa},
	}
	for _, przypadek := range przypadki {
		t.Run(przypadek.nazwa, func(t *testing.T) {
			pozycja := ZaleznoscZewnetrzna{Narzedzie: zewnetrzne.Narzedzie{
				Nazwa: przypadek.nazwa, Program: przypadek.program, Pakiet: przypadek.pakiet,
			}}
			if warstwa := WarstwaZaleznosci(pozycja); warstwa != przypadek.oczekuje {
				t.Errorf("warstwa %q, oczekiwano %q", warstwa, przypadek.oczekuje)
			}
		})
	}
}

// TestWarstwaPodpowiedziZdaniemZostajeObowiazkowa pilnuje, żeby pozycja z podpowiedzią zapisaną
// zdaniem albo poleceniem menedżera pakietów innego niż go/npm zostawała w jedynej warstwie,
// którą skrypt prowizjonowania rozdziela dalej po kształcie pola — nie w warstwie osobnej, dla
// której skrypt nie ma gałęzi i pozycja wpadłaby donikąd.
func TestWarstwaPodpowiedziZdaniemZostajeObowiazkowa(t *testing.T) {
	oczekiwane := map[string]bool{
		"java": true, "hunspell": true, "vale": true, "typst": true, "ruff": true, "semgrep": true,
	}
	znalezione := map[string]bool{}
	for _, pozycja := range ZaleznosciZewnetrzne() {
		program := pozycja.Narzedzie.Program
		if !oczekiwane[program] {
			continue
		}
		znalezione[program] = true
		if warstwa := WarstwaZaleznosci(pozycja); warstwa != WarstwaObowiazkowa {
			t.Errorf("pozycja %q z podpowiedzią %q trafiła do warstwy %q zamiast obowiązkowej",
				pozycja.Narzedzie.Nazwa, pozycja.Narzedzie.Pakiet, warstwa)
		}
	}
	for program := range oczekiwane {
		if !znalezione[program] {
			t.Errorf("w wykazie nie ma programu %q — sprawdzian nie ma czego pilnować", program)
		}
	}
}

// TestWarstwaCargoInstallNieTrafiaDoWarsztatuGo pilnuje usterki, w której reguła pytała o podnapis
// „go install", a ten stoi wewnątrz „cargo install typos-cli" („cargo” kończy się na „go”, dalej
// idzie spacja i „install”); pozycja wpadała do warstwy warsztatu Go i zostałaby wykonana
// poleceniem go install zamiast cargo install.
func TestWarstwaCargoInstallNieTrafiaDoWarsztatuGo(t *testing.T) {
	policzone := 0
	for _, pozycja := range ZaleznosciZewnetrzne() {
		if !strings.HasPrefix(strings.TrimSpace(pozycja.Narzedzie.Pakiet), "cargo install ") {
			continue
		}
		policzone++
		if warstwa := WarstwaZaleznosci(pozycja); warstwa == WarstwaWarsztatGo {
			t.Errorf("pozycja %q z poleceniem %q trafiła do warsztatu Go — zostałaby wykonana go install",
				pozycja.Narzedzie.Nazwa, pozycja.Narzedzie.Pakiet)
		}
	}
	if policzone == 0 {
		t.Fatal("w wykazie nie ma pozycji cargo install — sprawdzian nie ma czego pilnować")
	}
}

// TestWydrukMeldujeStanZgodnieZSonda pilnuje, żeby kolumna obecności nie
// rozjechała się z sondą: prowizjonowanie czyta ją, decydując, co dołożyć.
func TestWydrukMeldujeStanZgodnieZSonda(t *testing.T) {
	var bufor bytes.Buffer
	if err := WypiszWykazZaleznosci(&bufor); err != nil {
		t.Fatalf("wypisanie wykazu odmówiło: %v", err)
	}
	for _, pozycja := range ZaleznosciZewnetrzne() {
		wiersz := wierszProgramu(bufor.String(), pozycja.Narzedzie.Program)
		if wiersz == "" {
			t.Errorf("programu %q nie ma w wydruku", pozycja.Narzedzie.Program)
			continue
		}
		pola := strings.Split(wiersz, "\t")
		if len(pola) < 4 {
			continue
		}
		oczekiwane := "nie"
		if pozycja.Stoi {
			oczekiwane = "tak"
		}
		if pola[3] != oczekiwane {
			t.Errorf("program %q: wydruk melduje obecność %q, sonda %q",
				pozycja.Narzedzie.Program, pola[3], oczekiwane)
		}
	}
}

// TestWykazMowyNiesieRozmiarModeluIMiejscaArsenalu pilnuje warstwy mowy: bez
// rozmiaru modelu prowizjonowanie nie pobierze wag z góry i pierwsze użycie
// mikrofonu u Operatora zaczeka na sieć.
func TestWykazMowyNiesieRozmiarModeluIMiejscaArsenalu(t *testing.T) {
	var bufor bytes.Buffer
	if err := WypiszWykazMowy(&bufor); err != nil {
		t.Fatalf("wypisanie arsenału mowy odmówiło: %v", err)
	}
	wartosci := make(map[string]string)
	for _, wiersz := range strings.Split(strings.TrimRight(bufor.String(), "\n"), "\n") {
		if strings.HasPrefix(wiersz, "#") {
			continue
		}
		pola := strings.SplitN(wiersz, "\t", 2)
		if len(pola) != 2 {
			t.Errorf("wiersz %q nie ma postaci klucz<TAB>wartość", wiersz)
			continue
		}
		if _, powtorzony := wartosci[pola[0]]; powtorzony {
			t.Errorf("klucz %q pada dwa razy — skrypt czyta pierwszy i cicho gubi drugi", pola[0])
		}
		wartosci[pola[0]] = pola[1]
	}
	if wartosci["rozpoznanie.model-domyslny"] != mowa.ModelDomyslny {
		t.Errorf("rozmiar modelu w wykazie mowy to %q, a ustawienia rdzenia mówią %q",
			wartosci["rozpoznanie.model-domyslny"], mowa.ModelDomyslny)
	}
	if wartosci["rozpoznanie.ustawienie-katalogu-modeli"] != mowa.KluczKatalogModeli {
		t.Error("wykaz mowy nie mówi, którym ustawieniem wskazać rdzeniowi katalog wag")
	}
	if wartosci["synteza.piper-arsenal"] != piperArsenalu {
		t.Errorf("miejsce arsenału pipera w wykazie to %q, a rdzeń szuka w %q",
			wartosci["synteza.piper-arsenal"], piperArsenalu)
	}
	if wartosci["synteza.piper-glosy"] != glosyPiperaArsenalu {
		t.Error("wykaz mowy nie wskazuje katalogu głosów pipera")
	}
	if wartosci["synteza.espeak-program"] != silnikEspeak {
		t.Error("wykaz mowy nie nazywa syntezatora zapasowego")
	}
}

// TestZnacznikiWykazuRozpoznawaneWObuPostaciach pilnuje wejścia do trybów:
// prowizjonowanie woła rdzeń z podwójnym minusem, a reszta projektu zapisuje
// flagi pojedynczym. Obie postacie mają wchodzić.
func TestZnacznikiWykazuRozpoznawaneWObuPostaciach(t *testing.T) {
	if !ZadanoWykazZaleznosci([]string{"--wykaz-zaleznosci"}) ||
		!ZadanoWykazZaleznosci([]string{"-wykaz-zaleznosci"}) {
		t.Error("znacznik wykazu zależności nie jest rozpoznawany w obu postaciach")
	}
	if !ZadanoWykazMowy([]string{"--wykaz-mowy"}) || !ZadanoWykazMowy([]string{"-wykaz-mowy"}) {
		t.Error("znacznik wykazu mowy nie jest rozpoznawany w obu postaciach")
	}
	if ZadanoWykazZaleznosci([]string{"--rola", "all"}) || ZadanoWykazMowy([]string{"--rola", "all"}) {
		t.Error("zwykły start rdzenia wpadł w tryb wykazu — proces wypisałby wykaz zamiast pracować")
	}
	if ZadanoWykazZaleznosci(nil) || ZadanoWykazMowy(nil) {
		t.Error("brak argumentów uznano za żądanie wykazu")
	}
}

// wierszProgramu odnajduje w wydruku wiersz danych opisujący wskazany program.
// Program bywa w wykazie dwa razy (silnik kontenerów), brany jest pierwszy,
// bo sprawdziany pytają o warstwę i obecność, a te są dla obu takie same.
func wierszProgramu(wydruk, program string) string {
	for _, wiersz := range strings.Split(wydruk, "\n") {
		if strings.HasPrefix(wiersz, "#") {
			continue
		}
		pola := strings.Split(wiersz, "\t")
		if len(pola) >= 2 && pola[1] == program {
			return wiersz
		}
	}
	return ""
}
