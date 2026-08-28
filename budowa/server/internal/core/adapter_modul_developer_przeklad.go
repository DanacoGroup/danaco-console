// Odpowiedzialność pliku: przekład bytów modułu Developer na kształt kontraktu
// — plik edytora, przebieg budowania (z pamięci i z dziennika) oraz
// rozpoznanie zgłoszenia kompilatora w wierszu logu.
package core

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// jezykBinarny znakuje treść binarną, której Code Editor nie pokaże jako
// tekstu czytelnego dla człowieka.
const jezykBinarny = "binarny"

// jezykiPlikow wiąże rozszerzenie z nazwą języka podświetlania. Wykaz jest
// danymi: nowy język to nowy wiersz, nie nowy warunek.
var jezykiPlikow = map[string]string{
	".go": "go", ".ts": "typescript", ".tsx": "typescript", ".js": "javascript",
	".jsx": "javascript", ".mjs": "javascript", ".cjs": "javascript",
	".json": "json", ".sql": "sql", ".md": "markdown", ".css": "css",
	".html": "html", ".htm": "html", ".yml": "yaml", ".yaml": "yaml",
	".sh": "shell", ".ps1": "powershell", ".py": "python", ".rs": "rust",
	".toml": "toml", ".xml": "xml", ".txt": "plaintext", ".env": "dotenv",
}

// wzorceZgloszen rozpoznają miejsce w kodzie w wierszu logu budowania.
// Pierwszy wzorzec obsługuje postać `plik:wiersz:kolumna: treść` (go, gcc,
// eslint), drugi — `plik(wiersz,kolumna): rodzaj KOD: treść` (tsc, msvc).
var wzorceZgloszen = []*regexp.Regexp{
	regexp.MustCompile(`^\s*([^\s:]+(?::[^\s:]+)?):(\d+):(\d+):\s*(.+)$`),
	regexp.MustCompile(`^\s*(.+?)\((\d+),(\d+)\):\s*(.+)$`),
}

// opisPliku składa plik kontraktu z jego opisu na dysku, dobierając język
// podświetlania po rozszerzeniu.
func opisPliku(sciezka string, opis os.FileInfo) shared.DeveloperFile {
	rozmiar := opis.Size()
	return shared.DeveloperFile{
		Path:      sciezka,
		Language:  jezykPliku(sciezka),
		SizeBytes: &rozmiar,
		UpdatedAt: opis.ModTime().UTC().UnixMilli(),
	}
}

// jezykPliku rozpoznaje język po rozszerzeniu. Rozszerzenie nieznane zostawia
// pole puste — zgadywanie języka po treści dałoby edytorowi podświetlanie
// mylące, a kontrakt dopuszcza brak wskazania.
func jezykPliku(sciezka string) *string {
	jezyk, jest := jezykiPlikow[strings.ToLower(filepath.Ext(sciezka))]
	if !jest {
		return nil
	}
	return &jezyk
}

// budowanieKontraktu składa przebieg budowania z pamięci rdzenia na byt
// kontraktu widoczny oknu Developer.
func budowanieKontraktu(p *przebiegBudowania) shared.DeveloperBuild {
	stan, kodWyjscia, zgloszenia, zakonczono := p.Migawka()
	budowanie := shared.DeveloperBuild{
		Id:        p.kod,
		WindowId:  p.oknoKod,
		Task:      p.zadanie,
		Status:    stan,
		ExitCode:  kodWyjscia,
		Problems:  zgloszenia,
		LogRef:    wskaznikTekstu(p.kod),
		StartedAt: p.uruchomiono.UnixMilli(),
	}
	if !zakonczono.IsZero() {
		koniec := zakonczono.UnixMilli()
		budowanie.FinishedAt = &koniec
	}
	return budowanie
}

// budowanieZDziennika składa przebieg z wiersza dziennika. Zgłoszenia wracają
// z zachowanego ogona logu, bo dziennik trzyma log, a nie ich wykaz — jedno
// źródło prawdy zamiast dwóch.
func budowanieZDziennika(wiersz dane.PrzebiegBudowania) shared.DeveloperBuild {
	budowanie := shared.DeveloperBuild{
		Id:        wiersz.Kod,
		WindowId:  wiersz.OknoKod,
		Task:      wiersz.Zadanie,
		Status:    wiersz.Stan,
		Problems:  zgloszeniaZLogu(wartoscTekstu(wiersz.Log)),
		LogRef:    wskaznikTekstu(wiersz.Kod),
		StartedAt: chwilaBazy(wiersz.Uruchomiono),
	}
	if wiersz.KodWyjscia != nil {
		kod := int(*wiersz.KodWyjscia)
		budowanie.ExitCode = &kod
	}
	if wiersz.Zakonczono != nil {
		koniec := chwilaBazy(*wiersz.Zakonczono)
		budowanie.FinishedAt = &koniec
	}
	return budowanie
}

// zgloszeniaZLogu wyławia zgłoszenia kompilatora z zachowanego ogona logu
// budowania Developera rdzenia.
func zgloszeniaZLogu(log string) []shared.BuildProblem {
	if log == "" {
		return nil
	}
	zgloszenia := make([]shared.BuildProblem, 0, 8)
	for _, wiersz := range strings.Split(log, "\n") {
		if zgloszenie, jest := rozpoznajZgloszenie(wiersz); jest {
			zgloszenia = append(zgloszenia, zgloszenie)
		}
	}
	if len(zgloszenia) == 0 {
		return nil
	}
	return zgloszenia
}

// rozpoznajZgloszenie czyta z wiersza logu budowania miejsce w kodzie i wagę
// zgłoszenia kompilatora budowy.
func rozpoznajZgloszenie(wiersz string) (shared.BuildProblem, bool) {
	for _, wzorzec := range wzorceZgloszen {
		czesci := wzorzec.FindStringSubmatch(wiersz)
		if czesci == nil {
			continue
		}
		numer, err := strconv.Atoi(czesci[2])
		if err != nil {
			continue
		}
		kolumna, err := strconv.Atoi(czesci[3])
		if err != nil {
			continue
		}
		sciezka, tresc := czesci[1], strings.TrimSpace(czesci[4])
		return shared.BuildProblem{
			Path:     &sciezka,
			Line:     &numer,
			Column:   &kolumna,
			Severity: wagaZgloszenia(tresc),
			Message:  tresc,
		}, true
	}
	return shared.BuildProblem{}, false
}

// wagaZgloszenia rozpoznaje wagę po słowie wiodącym. Wiersz bez rozpoznanego
// słowa liczy się jako błąd: narzędzia budowania piszą tak głównie o błędach,
// a zaniżona waga ukryłaby powód niepowodzenia.
func wagaZgloszenia(tresc string) shared.ProblemSeverity {
	nizsza := strings.ToLower(tresc)
	switch {
	case strings.HasPrefix(nizsza, "warning"), strings.HasPrefix(nizsza, "ostrzeżenie"):
		return shared.ProblemSeverityWarning
	case strings.HasPrefix(nizsza, "note"), strings.HasPrefix(nizsza, "info"):
		return shared.ProblemSeverityInfo
	default:
		return shared.ProblemSeverityError
	}
}

// podsumowaniePrzebiegu składa ostatni wiersz logu budowania jako podsumowanie
// całego przebiegu budowy.
func podsumowaniePrzebiegu(stan shared.BuildStatus, kodWyjscia *int) string {
	tresc := "[budowanie " + string(stan)
	if kodWyjscia != nil {
		tresc += ", kod wyjścia " + strconv.Itoa(*kodWyjscia)
	}
	return tresc + "]"
}

// itoa skraca zapis liczby całkowitej w komunikatach błędów modułu Developer
// bez importu pakietu strconv.
func itoa(wartosc int) string {
	return strconv.Itoa(wartosc)
}
