package core

import (
	"bytes"
	"strings"
	"testing"

	"danacoconsole/server/internal/mowa"
	"danacoconsole/server/internal/zewnetrzne"
)

// Wykaz wypisany maszynowo jest jedyną drogą, którą prowizjonowanie serwera
// (`scripts/arsenal-serwera.sh`) dowiaduje się, co postawić. Szkoda, której te
// sprawdziany pilnują, ma trzy postacie:
//
//	wykaz okrojony     — program wołany przez rdzeń nie trafia do prowizjonowania
//	                     i funkcja odmawia Operatorowi na serwerze;
//	warstwa nieznana   — pozycja wypada z każdej warstwy skryptu i cicho nie
//	                     zostaje postawiona;
//	silnik kontenerów  — wstrzymany decyzją Właściciela, wchodzi do warstwy
//	                     stawianej milcząco.
//
// Sprawdziany nie mierzą, ile programów stoi na tej maszynie — mierzą mechanizm
// wyprowadzenia wykazu.

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
			// Puste dopuszczamy tylko w podpowiedzi pakietu (pole trzecie,
			// indeks 2) — deklaracja bez podpowiedzi jest legalna, choć wykaz
			// zależności pilnuje osobno, żeby jej nie było.
			if numer != 2 && strings.TrimSpace(pole) == "" {
				t.Errorf("wiersz %q ma puste pole numer %d", wiersz, numer+1)
			}
		}
	}
}

// TestWydrukNiesiePolskiPakietTesseractu pilnuje wymagania Właściciela
// wypowiedzianego wprost: rozpoznanie pisma ma działać po polsku. Sam
// `tesseract-ocr` daje program bez języka polskiego, więc prowizjonowanie
// postawiłoby OCR, który na polskim skanie zwraca miał.
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

// TestSilnikKontenerowStoiWWarstwieDecyzyjnej pilnuje decyzji Właściciela:
// silnik kontenerów został wstrzymany i nie może wejść do warstwy, którą
// prowizjonowanie stawia bez pytania. Obie jego deklaracje (warsztat Developera
// i moduł Terminal) muszą trafić do warstwy decyzyjnej.
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

// TestKazdaPozycjaMaZnanaWarstwe pilnuje, żeby nowe narzędzie dopisane do wykazu
// nie wypadło z prowizjonowania. Skrypt stawia warstwami; warstwa nieznana
// znaczy pozycję, której nikt nie postawi i nikt tego nie zauważy.
func TestKazdaPozycjaMaZnanaWarstwe(t *testing.T) {
	znane := map[string]bool{
		WarstwaObowiazkowa:  true,
		WarstwaWarsztatGo:   true,
		WarstwaWarsztatNpm:  true,
		WarstwaSnap:         true,
		WarstwaModelRecznie: true,
		WarstwaDecyzyjna:    true,
	}
	for _, pozycja := range ZaleznosciZewnetrzne() {
		if warstwa := WarstwaZaleznosci(pozycja); !znane[warstwa] {
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
		// Podpowiedź warsztatu Go potrafi nieść adres GitHuba. Gdyby reguła
		// pytała najpierw o GitHuba, golangci-lint stałby się krokiem ręcznym.
		{"warsztat Go z adresem GitHuba", "golangci-lint",
			"go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest", WarstwaWarsztatGo},
		{"warsztat npm", "eslint", "npm i -g eslint", WarstwaWarsztatNpm},
		{"snap", "kubectl", "kubectl (snap)", WarstwaSnap},
		{"wydanie z GitHuba", "realesrgan-ncnn-vulkan",
			"wydanie z github.com/xinntao/Real-ESRGAN/releases", WarstwaModelRecznie},
		{"środowisko pythonowe", "rembg", "rembg[cli] w osobnym środowisku pythonowym", WarstwaModelRecznie},
		{"silnik kontenerów", "docker", "docker.io albo podman", WarstwaDecyzyjna},
		{"silnik kontenerów pod podmanem", "podman", "podman", WarstwaDecyzyjna},
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
// flagi z pojedynczym. Obie postacie mają wchodzić, a zwykły start rdzenia nie
// może wpaść w tryb wykazu.
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
// Program bywa w wykazie dwa razy (silnik kontenerów) — bierzemy pierwszy,
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
