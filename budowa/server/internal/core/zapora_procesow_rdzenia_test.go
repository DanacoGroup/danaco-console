package core

import (
	"io/fs"
	"os"
	"path/filepath"
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
// z pominięciem `zewnetrzne.Wolaj`, wraz z powodem. Klucz jest ścieżką liczoną
// od `internal/`, bo zapora obejmuje cały ten katalog, a sama nazwa pliku nie
// mówi, w którym pakiecie wywołanie stoi. Wykaz jest ZAMKNIĘTY: dopisanie do
// niego pliku wymaga powodu tej samej wagi, co pozycje poniżej, i jest zmianą
// rozstrzygnięcia, nie porządkowaniem listy.
//
//   - `core/adapter_modul_extension_protokol.go` i
//     `core/adapter_modul_extension_integracje.go` uruchamiają program WSKAZANY
//     PRZEZ ROZSZERZENIE, a nie program produktu. `zewnetrzne.Wolaj` sprawdza
//     obecność narzędzia z deklaracji rdzenia — tu deklaracji nie ma, bo program
//     przychodzi z manifestu rozszerzenia.
//   - `core/adapter_modul_developer_narzedzia.go` pyta narzędzie o jego własną
//     wersję ścieżką już rozwiązaną; to sonda obecności, nie czynność Operatora.
//   - `injection/rozruch.go` JEST drogą, do której zapora odsyła: `zewnetrzne.Wolaj`
//     idzie portem `session.Uruchamiacz`, a ten port kończy się tutaj. Pozycja nie
//     jest wyjątkiem od reguły, tylko jej dnem — bez niej reguła nie ma się o co
//     oprzeć, a `zewnetrzne/wolanie.go` nazywa ten plik jedynym `exec.Command`
//     w drzewie.
//   - `zdalne/pliki.go` przenosi plik programem `scp` i nie ma dziś drzwi, przez
//     które mógłby przejść. Arsenał (`zewnetrzne.Wolaj`) zbiera całe wyjście do
//     pamięci pod obowiązkową granicą czasu — przenosiny pliku dowolnego rozmiaru
//     nie mają uczciwej granicy, a ich wyjście nie jest wynikiem do zebrania.
//     Spawner platformy stoi po stronie kanału, a zależność biegnie od kanału do
//     toru i nigdy odwrotnie (`zdalne/polecenie.go`), więc pakiet `zdalne` go nie
//     zaimportuje. Pozycja stoi tu po to, żeby brak drzwi był widoczny zamiast
//     niewidoczny — nagłówek `zdalne/zdalne.go` głosi, że pakiet procesów nie
//     uruchamia, a ten plik je uruchamia.
var plikiWlasnegoExecuUzasadnione = map[string]string{
	"core/adapter_modul_extension_protokol.go":   "program z manifestu rozszerzenia, nie z deklaracji rdzenia",
	"core/adapter_modul_extension_integracje.go": "program z manifestu rozszerzenia, nie z deklaracji rdzenia",
	"core/adapter_modul_developer_narzedzia.go":  "sonda wersji narzędzia po ścieżce już rozwiązanej",
	"injection/rozruch.go":                       "jedyny spawner platformy — port session.Uruchamiacz, którym idzie zewnetrzne.Wolaj, kończy się tutaj",
	"zdalne/pliki.go":                            "przenosiny scp bez drzwi: arsenał zbiera wyjście pod obowiązkową granicą czasu, a spawner platformy leży po stronie kanału, której pakiet zdalne nie importuje",
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

// korzenZapory wskazuje katalog, który zapora przegląda: całe `internal/`,
// licząc od katalogu pakietu core.
//
// Nie sam `internal/core`. Reguła jednej drogi do procesu jest regułą rdzenia,
// a nie regułą jednego pakietu: wywołanie przeniesione o katalog dalej wychodzi
// spod niej, choć szkodę robi tę samą. Zapora czytająca własny katalog nie
// widziałaby ani przenosin plików torem zdalnym, ani żadnego następnego
// wywołania spoza tego katalogu.
const korzenZapory = ".."

// TestRdzenUruchamiaProcesyJednaDroga pilnuje, żeby każdy proces rdzenia szedł
// przez `zewnetrzne.Wolaj`, poza zamkniętym wykazem pozycji uzasadnionych.
//
// Własny `exec.Command` pomija trzy rzeczy naraz: sprawdzenie obecności programu,
// bramę izolacji okna i objęcie drzewa procesów. Pierwsza zamienia brak programu
// w niezrozumiały błąd zamiast zdania nazywającego brak, druga wypuszcza czynność
// poza zasięg okna, trzecia zostawia sieroty po granicy czasu.
func TestRdzenUruchamiaProcesyJednaDroga(t *testing.T) {
	sprawdzonych, pozaRdzeniem := 0, 0
	wolajace := map[string]bool{}

	err := filepath.WalkDir(korzenZapory, func(sciezka string, wpis fs.DirEntry, blad error) error {
		if blad != nil {
			return blad
		}
		if wpis.IsDir() {
			return nil
		}
		nazwa := wpis.Name()
		if !strings.HasSuffix(nazwa, ".go") || strings.HasSuffix(nazwa, "_test.go") {
			return nil
		}
		wzgledna, blad := filepath.Rel(korzenZapory, sciezka)
		if blad != nil {
			return blad
		}
		wzgledna = filepath.ToSlash(wzgledna)

		tresc, blad := os.ReadFile(sciezka)
		if blad != nil {
			return blad
		}
		sprawdzonych++
		if !strings.HasPrefix(wzgledna, "core/") {
			pozaRdzeniem++
		}
		if !zawieraWolanieExecu(string(tresc)) {
			return nil
		}
		wolajace[wzgledna] = true

		powod, uzasadnione := plikiWlasnegoExecuUzasadnione[wzgledna]
		if !uzasadnione {
			t.Errorf("plik %s uruchamia proces własnym exec; jedyną drogą jest zewnetrzne.Wolaj "+
				"— ona sprawdza obecność programu, nakłada bramę izolacji okna i obejmuje drzewo "+
				"procesów granicą czasu", wzgledna)
			return nil
		}
		if powod == "" {
			t.Errorf("plik %s stoi w wykazie bez powodu; wykaz bez powodów jest listą "+
				"wyjątków, która rośnie sama", wzgledna)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("nie można przejrzeć %s: %v", korzenZapory, err)
	}

	// Dwie sondy zamiast jednej, bo dwa różne załamania dają ten sam cichy zielony
	// wynik: obchód, który nie dotknął niczego, i obchód, który zawęził się z
	// powrotem do własnego katalogu.
	if sprawdzonych == 0 {
		t.Fatal("zapora nie przejrzała ani jednego pliku — przestała czegokolwiek pilnować")
	}
	if pozaRdzeniem == 0 {
		t.Fatal("zapora nie wyszła poza internal/core — wywołanie przeniesione o katalog dalej " +
			"wychodziłoby spod reguły, choć szkodę robi tę samą")
	}

	// Pozycja wykazu, która procesu już nie uruchamia, jest wyjątkiem bez
	// przedmiotu — a wykaz z takimi pozycjami przestaje być zamknięty, bo nikt
	// nie wie, które z nich jeszcze coś znaczą.
	for sciezka := range plikiWlasnegoExecuUzasadnione {
		if !wolajace[sciezka] {
			t.Errorf("wykaz uzasadnionych niesie %s, a ten plik nie uruchamia procesu — "+
				"pozycję zdejmuje się wraz z wywołaniem, którego dotyczyła; jeśli plik "+
				"zmienił nazwę, popraw ścieżkę zamiast usuwać pozycję", sciezka)
		}
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
		// Warsztat kodu: osiem programów, po które sięgają komendy modułów
		// Developer i Terminal. Bez wpisu w wykazie sonda startowa milczałaby
		// o ich braku, a Operator dowiadywałby się o nim dopiero z funkcji,
		// która odmawia.
		{"ruff", "analiza plików Pythona"},
		{"semgrep", "poszerzenie skanu kodu o reguły semantyczne"},
		{"ast-grep", "wyszukanie i zamiana po składni"},
		{"jscpd", "powtórzenia w TypeScripcie i JavaScripcie"},
		{"dupl", "powtórzenia w plikach Go"},
		{"typos", "literówki w treści repozytorium"},
		{"stylelint", "analiza arkuszy CSS"},
		{"typescript-language-server", "warstwa językowa TypeScriptu"},
		// Interpretery kart Terminala. `terminal.script.lint` woła je wprost:
		// node orzeka o składni skryptu karty node, python3 o składni skryptu
		// karty python na maszynie bez Ruffa. Obecność interpretera na maszynie
		// deweloperskiej jest oczywista i właśnie dlatego brak wpisu przetrwał —
		// u Operatora ta sama funkcja odmawiałaby bez ostrzeżenia przy starcie.
		{"node", "orzeczenie o składni skryptu karty node"},
		{"python3", "orzeczenie o składni skryptu karty python bez Ruffa"},
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
