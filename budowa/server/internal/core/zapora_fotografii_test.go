package core

import (
	"os"
	"strings"
	"testing"
)

// Zapora warsztatu fotografii i części drukarskiej modułu Design.
//
// ── DLACZEGO TA ZAPORA ISTNIEJE ─────────────────────────────────────────────
// Obróbkę zdjęć i wydania drukarskie robi się zwykle programami zewnętrznymi.
// Sięgnięcie po nie jest tu jednym `exec.Command` i wygląda w kodzie niewinnie,
// a kosztuje CAŁĄ FUNKCJĘ u Operatora: arsenał produktu stoi wkompilowany
// w binarium serwera, a u Operatora leży cienka instalka — samo okno. Funkcja
// zależna od programu, którego instalka nie niesie, jest u niego odmową, nie
// funkcją. Na maszynie deweloperskiej, gdzie te programy bywają doinstalowane
// ręcznie, sprawdzian takiej funkcji świeciłby zielono i nikt by się nie
// dowiedział.
//
// ── CZEGO PILNUJE ───────────────────────────────────────────────────────────
//  1. Pliki rodzin `design.photo.*` i `design.print.*` nie wołają procesu:
//     ani przez `exec.Command`, ani przez pomocnika drzewa `zewnetrzne.Wolaj`.
//  2. ŻADEN plik obszaru Design nie wymienia nazw silników obrazu
//     i obrysowywania konturów, po które sięga się przy takiej pracy — także
//     w komentarzu. Nazwa w komentarzu jest wskazówką dla następnego wykonawcy,
//     a wskazówka w tę stronę jest wskazówką ku szkodzie.
//
// ── CZEGO NIE WOLNO ZROBIĆ Z TĄ ZAPORĄ ──────────────────────────────────────
// Nie wolno jej osłabić ani przestawić. Gdy nowa funkcja potrzebuje czegoś,
// czego rdzeń nie ma, drogą jest biblioteka Go wkompilowana przez `go.mod` albo
// ZGŁOSZENIE braku — nigdy program w miejsce biblioteki, która istnieje. Ten sam
// warunek pilnuje warsztatu dokumentu Studio (`zapora_warsztatu_pdf_test.go`)
// i z tego samego powodu.
//
// Wyjątek wolno dopisać w JEDNYM przypadku: gdy dla danej pracy nie istnieje
// żadna biblioteka czysto-Go, a program jest składnikiem pakietu serwera i stoi
// w wykazie zależności. Dziś taki wyjątek jest jeden — odczyt liter przy
// wniesieniu makiety (`plikiDesignuZOdczytemPisma`) — i jest ograniczony do
// jednego pliku oraz jednego programu. Wyjątek bez tych dwóch ograniczeń nie
// jest wyjątkiem, tylko zdjęciem zapory.
//
// Zapora jest napisana GRUBO — na wystąpienie nazwy w treści pliku, bez
// rozróżniania kodu od komentarza. Rozróżnienie wymagałoby rozbioru składni,
// a wtedy zapora sama stałaby się kodem, który może się mylić.

// silnikiSpozaInstalki wylicza nazwy, których w obszarze Design być nie może.
//
// Wykaz jest nazwami PROGRAMÓW, nie bibliotek: `pdfcpu`, `tdewolff/canvas`,
// `disintegration/imaging` i `x/image` są bibliotekami Go wkompilowanymi
// w binarium i wolno ich używać wszędzie. Zakaz dotyczy tego, co trzeba by
// URUCHOMIĆ jako proces.
var silnikiSpozaInstalki = []string{
	"vips",
	"imagemagick",
	"potrace",
	"inkscape",
	"fontforge",
	"fonttools",
	"graphicsmagick",
}

// wolaniaProcesu wylicza drogi uruchomienia procesu z pliku rdzenia.
var wolaniaProcesu = []string{
	"exec.Command",
	"exec.CommandContext",
	"zewnetrzne.Wolaj",
	"os/exec",
}

// wolaniaWlasneProcesu to drogi, które zapora zabrania ZAWSZE i wszędzie,
// niezależnie od wyjątków niżej: uruchomienie procesu poza pakietem `zewnetrzne`
// pomija sprawdzenie obecności programu, bramę izolacji okna i objęcie drzewa
// procesów granicą czasu.
var wolaniaWlasneProcesu = []string{
	"exec.Command",
	"exec.CommandContext",
	"os/exec",
}

// plikiDesignuZOdczytemPisma to JEDYNY wyjątek tej zapory, nazwany wraz
// z programem, którego dotyczy.
//
// ── Dlaczego wyjątek istnieje ───────────────────────────────────────────────
// Reguła produktu ma dwie połowy. Pierwsza: jest biblioteka Go — bierz
// bibliotekę, a `exec` jest regresem; tego pilnuje cała ta zapora. Druga: NIE MA
// biblioteki Go — program jest składnikiem pakietu serwera. Rozpoznanie pisma
// jest właśnie drugim przypadkiem: czytnika liter w czystym Go nie ma, a
// Tesseract stoi w wykazie zależności pakietu serwera (`zaleznosci_zewnetrzne.go`)
// i jest wołany tą samą drogą, co w Studiu i Translate.
//
// `design.mockup.import` bez odczytu liter rozpoznaje układ obszarów i to, które
// z nich są liniami tekstu, ale treści napisów nie czyta. Zakaz procesu bez tego
// wyjątku nie chronił tu instalki — odbierał funkcję, której nie ma czym zastąpić.
//
// ── Czego wyjątek NIE otwiera ───────────────────────────────────────────────
// Obowiązuje jeden plik i jeden program. Własny `exec` zostaje zabroniony i tu.
// Silniki obrazu z `silnikiSpozaInstalki` zostają zabronione w całym obszarze
// Design, ten plik włącznie — bo dla nich biblioteka Go istnieje.
var plikiDesignuZOdczytemPisma = map[string]string{
	"adapter_modul_design_makiety_zrzut.go": "tesseract",
}

// TestWarsztatFotografiiNieWolaProcesu pilnuje, żeby żaden plik rodzin
// `design.photo.*` i `design.print.*` nie uruchamiał procesu.
//
// Pliki rozpoznaje się po nazwie, bo tak leżą w drzewie: warsztat fotografii
// w `adapter_modul_design_fotografia*.go`, część drukarska
// w `adapter_modul_design_druk*.go`, a rachunek wektora, ikon i wykresów —
// w plikach, które te dwie rodziny wołają. Wykaz przedrostków obejmuje CAŁY
// obszar Design, bo droga do procesu wprowadzona w pliku sąsiednim byłaby tą
// samą szkodą.
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
			// Wyjątek nazwany: proces wolno wołać, ale wyłącznie drogą pakietu
			// `zewnetrzne` i wyłącznie dla programu wpisanego w wykaz wyjątku.
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
	// Zapora, która nie przejrzała ani jednego pliku, jest zaporą, która nie
	// broni niczego — a wygląda w wyniku identycznie jak zapora spełniona.
	if sprawdzonych == 0 {
		t.Fatal("zapora nie znalazła ani jednego pliku obszaru Design; sprawdź, czy nazwy " +
			"plików nie zmieniły przedrostka — inaczej ta zapora przestała czegokolwiek pilnować")
	}
}

// TestObszarDesignNieWymieniaSilnikowObrazuSpozaInstalki pilnuje, żeby nazwy
// silników obrazu i obrysowywania konturów nie pojawiły się w plikach obszaru
// Design żadną drogą — także przez komentarz, który wskazywałby następnemu
// wykonawcy tę drogę jako możliwą.
//
// Zakres jest obszarem Design (`adapter_modul_design*.go`), nie całym rdzeniem,
// i jest to POMIAR granicy odpowiedzialności, nie ustępstwo: warsztat obrazu
// modułu Library i narzędzia obrazowe modelu mają własne nagłówki, w których te
// nazwy stoją jako zapis decyzji ich autorów. Poszerzenie tej zapory na cały
// rdzeń wymagałoby przepisania cudzych plików — a to jest osobna praca
// i osobna zgoda. Zapora obszaru Design jest warunkiem, który TEN obszar
// spełnia w całości i który da się utrzymać.
func TestObszarDesignNieWymieniaSilnikowObrazuSpozaInstalki(t *testing.T) {
	wpisy, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("nie można przejrzeć rdzenia: %v", err)
	}
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
		maly := strings.ToLower(string(tresc))
		for _, silnik := range silnikiSpozaInstalki {
			if strings.Contains(maly, silnik) {
				t.Errorf("plik %s wymienia silnik %s; instalka Operatora go nie niesie, więc "+
					"funkcja od niego zależna byłaby u Operatora odmową — a nazwa w treści pliku "+
					"jest wskazówką ku tej szkodzie", nazwa, silnik)
			}
		}
	}
}
