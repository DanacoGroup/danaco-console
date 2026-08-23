package core

import (
	"os"
	"strings"
	"testing"
)

// Zapora procesów CAŁEGO rdzenia — trzecia po zaporze warsztatu dokumentu
// (`zapora_warsztatu_pdf_test.go`) i zaporze fotografii
// (`zapora_fotografii_test.go`), i pierwsza, która nie ogranicza się do jednego
// obszaru.
//
// ── DLACZEGO POWSTAŁA I CZEGO NIE POWTARZA ──────────────────────────────────
// Dwie starsze zapory zabraniają procesu w swoich obszarach, bo tam każdą pracę
// wykonuje w całości biblioteka Go. Tej reguły nie da się rozciągnąć na rdzeń
// bez kłamstwa: rozpoznanie pisma, archiwa, mowa, nagrania i silniki neuronowe
// nie mają biblioteki czysto-Go, więc ich programy są **składnikiem pakietu
// serwera** i wolno je wołać. Granicą nie jest „proces czy biblioteka", a:
//
//  1. czy program jest wołany JEDYNĄ dozwoloną drogą (`zewnetrzne.Wolaj`), która
//     sprawdza obecność, obejmuje drzewo procesów i odmawia zdaniem nazywającym
//     brak — a nie własnym `exec.Command`, który żadnej z tych rzeczy nie robi;
//  2. czy pakiet serwera ten program NIESIE, czyli czy stoi w wykazie
//     zależności (`zaleznosci_zewnetrzne.go`). Program wołany bez wpisu w wykazie
//     to cichy wymóg wobec wdrożenia: u Operatora funkcja odmawia, a przy starcie
//     nikt nie powiedział, czego brakuje.
//
// Trzecia rzecz, której zapora pilnuje, jest odwrotna do dwóch pierwszych:
// rachunek, który JUŻ jest wkompilowany, nie ma prawa wrócić do procesu. Cztery
// czynności obrazu modelu (`image.inspect`, `image.transform`, `image.adjust`,
// `image.convert`) liczyły się kiedyś programem, choć biblioteka wystarcza —
// i właśnie dlatego ta zapora wymienia ich pliki po nazwie.

// plikiRachunkuObrazuRdzenia wylicza pliki, w których obraz liczy się WYŁĄCZNIE
// biblioteką wkompilowaną. Proces w którymkolwiek z nich jest regresem: te same
// piksele policzy `imaging`, `x/image` i `nativewebp` bez uruchamiania niczego.
var plikiRachunkuObrazuRdzenia = []string{
	"adapter_narzedzia_obraz_wkompilowany.go",
	"adapter_narzedzia_obraz_raster.go",
	"adapter_narzedzia_obraz_warstwy.go",
	"adapter_narzedzia_obraz_wektor.go",
	"adapter_narzedzia_obraz_wektor_slad.go",
	"adapter_narzedzia_obraz_zlozenie.go",
	"adapter_modul_library_odciski.go",
}

// plikiWlasnegoExecuUzasadnione wylicza pliki, w których proces uruchamia się
// z pominięciem `zewnetrzne.Wolaj`, wraz z powodem. Wykaz jest ZAMKNIĘTY:
// dopisanie do niego pliku wymaga powodu tej samej wagi, co dwa poniżej, i jest
// zmianą rozstrzygnięcia, nie porządkowaniem listy.
//
//   - `adapter_modul_extension_protokol.go` i `adapter_modul_extension_integracje.go`
//     uruchamiają program WSKAZANY PRZEZ ROZSZERZENIE, a nie program produktu.
//     `zewnetrzne.Wolaj` sprawdza obecność narzędzia z deklaracji rdzenia —
//     tu deklaracji nie ma, bo program przychodzi z manifestu rozszerzenia.
//   - `adapter_modul_developer_narzedzia.go` pyta narzędzie o jego własną wersję
//     ścieżką już rozwiązaną; to sonda obecności, nie czynność Operatora.
var plikiWlasnegoExecuUzasadnione = map[string]string{
	"adapter_modul_extension_protokol.go":   "program z manifestu rozszerzenia, nie z deklaracji rdzenia",
	"adapter_modul_extension_integracje.go": "program z manifestu rozszerzenia, nie z deklaracji rdzenia",
	"adapter_modul_developer_narzedzia.go":  "sonda wersji narzędzia po ścieżce już rozwiązanej",
}

// TestRachunekObrazuRdzeniaNieWolaProcesu pilnuje, żeby pliki liczące obraz
// biblioteką wkompilowaną nie sięgnęły po proces żadną drogą.
func TestRachunekObrazuRdzeniaNieWolaProcesu(t *testing.T) {
	for _, nazwa := range plikiRachunkuObrazuRdzenia {
		tresc, err := os.ReadFile(nazwa)
		if err != nil {
			// Plik zniknął albo zmienił nazwę — zapora przestałaby wtedy pilnować
			// obszaru, o którym myśli, że go pilnuje.
			t.Fatalf("nie można odczytać %s: %v; jeśli plik zmienił nazwę, popraw wykaz "+
				"zapory, a nie usuwaj z niego pozycji", nazwa, err)
		}
		for _, wolanie := range wolaniaProcesu {
			if strings.Contains(string(tresc), wolanie) {
				t.Errorf("plik %s woła proces (%s); rachunek na pikselach robią tu biblioteki "+
					"wkompilowane w binarium — program wołany tam, gdzie biblioteka wystarcza, "+
					"kosztuje uruchomienie procesu i wiąże funkcję z cudzym wydaniem", nazwa, wolanie)
			}
		}
	}
}

// TestCzteryCzynnosciObrazuLiczaSieWkompilowane pilnuje, że program pakietu
// serwera został w czynnościach obrazu WYJĄTKIEM, a nie drogą podstawową.
//
// Mierzy dwie rzeczy naraz: że wołanie programu stoi wyłącznie w dwóch
// funkcjach drogi zapasowej i że każda z czterech czynności przechodzi przez
// rachunek wkompilowany. Sprawdzian liczy funkcje, a nie samą liczbę wystąpień,
// bo trzecie wołanie dopisane do istniejącej funkcji jest tą samą szkodą.
func TestCzteryCzynnosciObrazuLiczaSieWkompilowane(t *testing.T) {
	const plik = "adapter_narzedzia_obraz_czynnosci.go"
	tresc, err := os.ReadFile(plik)
	if err != nil {
		t.Fatalf("nie można odczytać %s: %v", plik, err)
	}

	funkcjeZapasowe := map[string]bool{
		"zbadajProgramem": false,
		"policzProgramem": false,
	}
	biezaca := ""
	for _, wiersz := range strings.Split(string(tresc), "\n") {
		if strings.HasPrefix(wiersz, "func ") {
			biezaca = nazwaFunkcjiZWiersza(wiersz)
		}
		if !strings.Contains(wiersz, "wolajImageMagick(") {
			continue
		}
		if _, dozwolona := funkcjeZapasowe[biezaca]; !dozwolona {
			t.Errorf("funkcja %s woła program pakietu serwera; wolno to wyłącznie drodze "+
				"zapasowej (zbadajProgramem, policzProgramem) dla AVIF-a i WEBP-a stratnego — "+
				"reszta rodziny image.* liczy się biblioteką wkompilowaną", biezaca)
			continue
		}
		funkcjeZapasowe[biezaca] = true
	}
	for nazwa, obecna := range funkcjeZapasowe {
		if !obecna {
			t.Errorf("droga zapasowa %s przestała wołać program; AVIF i WEBP stratny nie mają "+
				"kodera w Go, więc bez tej drogi te wyjścia zniknęły z produktu — jeśli koder "+
				"czysto-Go już istnieje, zdejmij tę pozycję z zapory wraz z drogą", nazwa)
		}
	}

	for _, rachunek := range []string{
		"zbadajObrazArsenalu(",
		"przeksztalcObrazArsenalu(",
		"poprawObrazArsenalu(",
		"zakodujObrazArsenalu(",
	} {
		if !strings.Contains(string(tresc), rachunek) {
			t.Errorf("%s nie woła %s — czynność przestała liczyć się biblioteką wkompilowaną",
				plik, rachunek)
		}
	}
}

// TestRdzenUruchamiaProcesyJednaDroga pilnuje, żeby każdy proces rdzenia szedł
// przez `zewnetrzne.Wolaj`, poza zamkniętym wykazem pozycji uzasadnionych.
//
// Własny `exec.Command` pomija trzy rzeczy naraz: sprawdzenie obecności programu,
// bramę izolacji okna i objęcie drzewa procesów. Pierwsza zamienia brak programu
// w niezrozumiały błąd zamiast zdania nazywającego brak, druga wypuszcza czynność
// poza zasięg okna, trzecia zostawia sieroty po granicy czasu.
func TestRdzenUruchamiaProcesyJednaDroga(t *testing.T) {
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
		tresc, err := os.ReadFile(nazwa)
		if err != nil {
			t.Fatalf("nie można odczytać %s: %v", nazwa, err)
		}
		sprawdzonych++
		if !zawieraWolanieExecu(string(tresc)) {
			continue
		}
		if powod, uzasadnione := plikiWlasnegoExecuUzasadnione[nazwa]; uzasadnione {
			if powod == "" {
				t.Errorf("plik %s stoi w wykazie bez powodu; wykaz bez powodów jest listą "+
					"wyjątków, która rośnie sama", nazwa)
			}
			continue
		}
		t.Errorf("plik %s uruchamia proces własnym exec; jedyną drogą jest zewnetrzne.Wolaj "+
			"— ona sprawdza obecność programu, nakłada bramę izolacji okna i obejmuje drzewo "+
			"procesów granicą czasu", nazwa)
	}
	if sprawdzonych == 0 {
		t.Fatal("zapora nie przejrzała ani jednego pliku rdzenia — przestała czegokolwiek pilnować")
	}
}

// TestProgramyRdzeniaStojaWWykazieZaleznosci pilnuje, żeby każdy program wołany
// przez rdzeń był zadeklarowany w wykazie zależności pakietu serwera.
//
// To jest właściwa miara gotowości funkcji opartej o program: nie „czy w kodzie
// jest exec", a „czy pakiet serwera to niesie i czy Operator dowie się o braku
// przy starcie". Program wołany bez wpisu w wykazie jest cichym wymogiem wobec
// wdrożenia — na maszynie deweloperskiej, gdzie ktoś doinstalował go ręcznie,
// wygląda jak funkcja gotowa.
func TestProgramyRdzeniaStojaWWykazieZaleznosci(t *testing.T) {
	wykaz := zaleznosciZewnetrzne()
	if len(wykaz) == 0 {
		t.Fatal("wykaz zależności pakietu serwera jest pusty — nie ma czego pilnować")
	}
	zadeklarowane := map[string]struct{}{}
	for _, pozycja := range wykaz {
		zadeklarowane[strings.ToLower(strings.TrimSpace(pozycja.Narzedzie.Program))] = struct{}{}
	}

	// ImageMagick musi w wykazie stać, bo dwa wyjścia rodziny `image.*` bez niego
	// nie powstaną. Sprawdzamy to wprost: wpis wniesiony wraz z rachunkiem
	// wkompilowanym łatwo usunąć „przy porządkach", a wtedy brak programu wróci
	// do bycia niespodzianką w chwili czynności.
	if _, stoi := zadeklarowane["magick"]; !stoi {
		t.Error("wykaz zależności nie zna programu magick; AVIF i WEBP stratny idą tą drogą, " +
			"więc bez wpisu Operator nie dowie się przy starcie, czego serwerowi brakuje")
	}

	for _, narzedzie := range []struct {
		program string
		zakres  string
	}{
		{"tesseract", "rozpoznanie pisma"},
		{"ffmpeg", "zamiana formatu nagrania"},
		{"7z", "archiwa"},
	} {
		if _, stoi := zadeklarowane[narzedzie.program]; !stoi {
			t.Errorf("wykaz zależności nie zna programu %s (%s); ten obszar nie ma biblioteki "+
				"czysto-Go, więc program jest składnikiem pakietu serwera i musi być w wykazie",
				narzedzie.program, narzedzie.zakres)
		}
	}
}

// zawieraWolanieExecu orzeka, czy plik uruchamia proces własną drogą. Sam import
// `os/exec` nie wystarczy: pakiet bywa importowany dla typu błędu albo
// `exec.LookPath`, a to nie jest uruchomienie.
func zawieraWolanieExecu(tresc string) bool {
	return strings.Contains(tresc, "exec.Command(") ||
		strings.Contains(tresc, "exec.CommandContext(")
}

// nazwaFunkcjiZWiersza wyciąga nazwę funkcji z wiersza deklaracji — także wtedy,
// gdy funkcja jest metodą (`func (a *typ) Nazwa(`).
func nazwaFunkcjiZWiersza(wiersz string) string {
	reszta := strings.TrimPrefix(wiersz, "func ")
	if strings.HasPrefix(reszta, "(") {
		if koniec := strings.Index(reszta, ")"); koniec >= 0 {
			reszta = strings.TrimSpace(reszta[koniec+1:])
		}
	}
	if nawias := strings.Index(reszta, "("); nawias >= 0 {
		reszta = reszta[:nawias]
	}
	return strings.TrimSpace(reszta)
}
