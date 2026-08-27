package repozytorium

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"

	"danacoconsole/shared"
)

// Silnikiem wyszukiwania i zamiany jest regexp biblioteki standardowej, nie
// ripgrep; pliki ignorowane przez repozytorium są pomijane.
//
// ZadanieSzukania opisuje, czego i gdzie szukać: wzorzec, tryb dopasowania,
// wzorce ścieżek i granicę wyniku.
type ZadanieSzukania struct {
	Wzorzec            string
	Wyrazenie          bool
	RozrozniajWielkosc bool
	Wlacz              []string
	Wylacz             []string
	Granica            int
}

// GranicaDomyslna jest górną granicą trafień, gdy żądanie jej nie podaje, bo
// odpowiedź idzie kopertą przez gniazdo.
const GranicaDomyslna = 2000

// GranicaPlikuBajtow pomija pliki większe od tej granicy.
//
// Plik stumegabajtowy w repozytorium jest zwykle zrzutem danych albo artefaktem
// budowania, a nie kodem; jego przeszukanie kosztuje czas, którego wynik i tak
// nikomu nie służy.
const GranicaPlikuBajtow = 4 << 20

// Szukaj przechodzi repozytorium i oddaje trafienia w kolejności plików, wraz
// z liczbą trafień przed przycięciem granicą.
func Szukaj(katalog string, zadanie ZadanieSzukania) ([]shared.DeveloperGrepMatch, int, bool, error) {
	wzorzec, err := zlozWzorzec(zadanie.Wzorzec, zadanie.Wyrazenie, zadanie.RozrozniajWielkosc)
	if err != nil {
		return nil, 0, false, err
	}
	korzen, pomijaj, err := obszarSzukania(katalog)
	if err != nil {
		return nil, 0, false, err
	}
	granica := zadanie.Granica
	if granica <= 0 {
		granica = GranicaDomyslna
	}

	trafienia := make([]shared.DeveloperGrepMatch, 0)
	wszystkich := 0

	err = przejdzPliki(korzen, pomijaj, zadanie.Wlacz, zadanie.Wylacz,
		func(wzgledna string, tresc string) error {
			for numer, wiersz := range strings.Split(tresc, "\n") {
				miejsca := wzorzec.FindAllStringIndex(wiersz, -1)
				for _, miejsce := range miejsca {
					wszystkich++
					if len(trafienia) >= granica {
						continue
					}
					kolumna := miejsce[0] + 1
					trafienia = append(trafienia, shared.DeveloperGrepMatch{
						Path:   wzgledna,
						Line:   numer + 1,
						Column: &kolumna,
						Text:   wiersz,
					})
				}
			}
			return nil
		})
	if err != nil {
		return nil, 0, false, err
	}
	return trafienia, wszystkich, wszystkich > len(trafienia), nil
}

// ZadanieZamiany opisuje daną zamianę w tym repozytorium: wzorzec, tekst
// zamiany, ścieżki i tryb podglądu.
type ZadanieZamiany struct {
	Wzorzec   string
	Zamiana   string
	Wyrazenie bool
	Sciezki   []string
	Podglad   bool
}

// Zamien wykonuje zamianę wzorca w repozytorium, oddając wykaz zmian razem
// z odpowiedzią o tym, czy je zapisała.
func Zamien(katalog string, zadanie ZadanieZamiany) ([]shared.DeveloperTextEdit, []string, bool, error) {
	wzorzec, err := zlozWzorzec(zadanie.Wzorzec, zadanie.Wyrazenie, true)
	if err != nil {
		return nil, nil, false, err
	}
	korzen, pomijaj, err := obszarSzukania(katalog)
	if err != nil {
		return nil, nil, false, err
	}
	wybrane := map[string]bool{}
	for _, sciezka := range zadanie.Sciezki {
		wybrane[filepath.ToSlash(strings.TrimSpace(sciezka))] = true
	}

	zmiany := make([]shared.DeveloperTextEdit, 0)
	zmienione := make([]string, 0)

	err = przejdzPliki(korzen, pomijaj, nil, nil, func(wzgledna, tresc string) error {
		if len(wybrane) > 0 && !wybrane[wzgledna] {
			return nil
		}
		wiersze := strings.Split(tresc, "\n")
		dotkniety := false
		for numer, wiersz := range wiersze {
			if !wzorzec.MatchString(wiersz) {
				continue
			}
			var nowy string
			if zadanie.Wyrazenie {
				nowy = wzorzec.ReplaceAllString(wiersz, zadanie.Zamiana)
			} else {
				// Wzorzec dosłowny zamienia się dosłownie: znacznik grupy
				// w treści wstawianej ma zostać, nie zniknąć.
				nowy = wzorzec.ReplaceAllLiteralString(wiersz, zadanie.Zamiana)
			}
			if nowy == wiersz {
				continue
			}
			poczatek := 1
			koniec := len(wiersz) + 1
			zmiany = append(zmiany, shared.DeveloperTextEdit{
				Path:        wzgledna,
				StartLine:   numer + 1,
				StartColumn: &poczatek,
				EndLine:     numer + 1,
				EndColumn:   &koniec,
				NewText:     nowy,
			})
			wiersze[numer] = nowy
			dotkniety = true
		}
		if !dotkniety {
			return nil
		}
		zmienione = append(zmienione, wzgledna)
		if zadanie.Podglad {
			return nil
		}
		pelna := filepath.Join(korzen, filepath.FromSlash(wzgledna))
		opis, err := os.Stat(pelna)
		prawa := os.FileMode(0o644)
		if err == nil {
			prawa = opis.Mode().Perm()
		}
		if err := os.WriteFile(pelna, []byte(strings.Join(wiersze, "\n")), prawa); err != nil {
			return fmt.Errorf("zapis pliku %s: %w", wzgledna, err)
		}
		return nil
	})
	if err != nil {
		return nil, nil, false, err
	}
	sort.Strings(zmienione)
	return zmiany, zmienione, !zadanie.Podglad, nil
}

// zlozWzorzec buduje wyrażenie szukające.
//
// Wzorzec dosłowny przechodzi przez `QuoteMeta`, więc kropka i gwiazdka znaczą
// w nim siebie same. Operator, który szuka `config.json`, nie szuka „config"
// z dowolnym znakiem w środku.
func zlozWzorzec(wzorzec string, wyrazenie, rozrozniajWielkosc bool) (*regexp.Regexp, error) {
	if strings.TrimSpace(wzorzec) == "" {
		return nil, fmt.Errorf("wzorzec pusty nie zawęża niczego")
	}
	tresc := wzorzec
	if !wyrazenie {
		tresc = regexp.QuoteMeta(wzorzec)
	}
	if !rozrozniajWielkosc {
		tresc = "(?i)" + tresc
	}
	zlozony, err := regexp.Compile(tresc)
	if err != nil {
		return nil, fmt.Errorf("wzorzec %q nie jest poprawnym wyrażeniem: %w", wzorzec, err)
	}
	return zlozony, nil
}

// obszarSzukania oddaje korzeń przeszukania oraz regułę pomijania; katalog
// bez repozytorium też podlega przeszukaniu.
func obszarSzukania(katalog string) (string, gitignore.Matcher, error) {
	korzen := katalog
	if repo, err := Otworz(katalog); err == nil {
		if znaleziony, err := Korzen(repo); err == nil {
			korzen = znaleziony
		}
	}
	if strings.TrimSpace(korzen) == "" {
		return "", nil, fmt.Errorf("katalog roboczy okna nie jest wskazany")
	}
	if opis, err := os.Stat(korzen); err != nil || !opis.IsDir() {
		return "", nil, fmt.Errorf("katalog roboczy %s nie istnieje", korzen)
	}
	return korzen, regulyPomijania(korzen), nil
}

// regulyPomijania czyta plik ignorowanych ścieżek korzenia repozytorium, nie
// wszystkich plików całego drzewa.
func regulyPomijania(korzen string) gitignore.Matcher {
	bajty, err := os.ReadFile(filepath.Join(korzen, ".gitignore"))
	if err != nil {
		return gitignore.NewMatcher(nil)
	}
	wzorce := make([]gitignore.Pattern, 0)
	for _, wiersz := range strings.Split(string(bajty), "\n") {
		wiersz = strings.TrimSpace(wiersz)
		if wiersz == "" || strings.HasPrefix(wiersz, "#") {
			continue
		}
		wzorce = append(wzorce, gitignore.ParsePattern(wiersz, nil))
	}
	return gitignore.NewMatcher(wzorce)
}

// przejdzPliki woła obsługę dla każdego pliku tekstowego obszaru, pomijając
// pliki binarne i nadmiernie duże.
func przejdzPliki(korzen string, pomijaj gitignore.Matcher, wlacz, wylacz []string,
	obsluz func(wzgledna, tresc string) error) error {

	return filepath.WalkDir(korzen, func(sciezka string, wpis fs.DirEntry, err error) error {
		if err != nil {
			// Katalog bez prawa odczytu nie zatrzymuje przeszukania całości;
			// wynik ma być niepełny, nie pusty.
			if wpis != nil && wpis.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		wzgledna, err := filepath.Rel(korzen, sciezka)
		if err != nil {
			return nil
		}
		wzgledna = filepath.ToSlash(wzgledna)
		if wzgledna == "." {
			return nil
		}
		czesci := strings.Split(wzgledna, "/")

		if wpis.IsDir() {
			if wpis.Name() == ".git" || pomijaj.Match(czesci, true) {
				return filepath.SkipDir
			}
			return nil
		}
		if pomijaj.Match(czesci, false) {
			return nil
		}
		if !przechodziGloby(wzgledna, wlacz, wylacz) {
			return nil
		}
		opis, err := wpis.Info()
		if err != nil || opis.Size() > GranicaPlikuBajtow {
			return nil
		}
		bajty, err := os.ReadFile(sciezka)
		if err != nil {
			return nil
		}
		if czyBinarna(string(bajty)) {
			return nil
		}
		return obsluz(wzgledna, string(bajty))
	})
}

// przechodziGloby rozstrzyga, czy plik mieści się we wzorcach włączających
// i poza wzorcami wyłączającymi.
func przechodziGloby(wzgledna string, wlacz, wylacz []string) bool {
	for _, wzorzec := range wylacz {
		if pasujeGlob(wzgledna, wzorzec) {
			return false
		}
	}
	if len(wlacz) == 0 {
		return true
	}
	for _, wzorzec := range wlacz {
		if pasujeGlob(wzgledna, wzorzec) {
			return true
		}
	}
	return false
}

// pasujeGlob dopasowuje wzorzec do ścieżki albo do samej nazwy pliku,
// dwutorowo, z obsługą podwójnej gwiazdki.
func pasujeGlob(wzgledna, wzorzec string) bool {
	wzorzec = strings.TrimSpace(wzorzec)
	if wzorzec == "" {
		return false
	}
	if gwiazdki := strings.Index(wzorzec, "**"); gwiazdki >= 0 {
		przedrostek := strings.TrimSuffix(wzorzec[:gwiazdki], "/")
		return przedrostek == "" || wzgledna == przedrostek ||
			strings.HasPrefix(wzgledna, przedrostek+"/")
	}
	if pasuje, err := filepath.Match(wzorzec, wzgledna); err == nil && pasuje {
		return true
	}
	pasuje, err := filepath.Match(wzorzec, filepath.Base(wzgledna))
	return err == nil && pasuje
}

// PlikiDoPrzejrzenia oddaje ścieżki bezwzględne plików tekstowych obszaru,
// z poszanowaniem reguł pomijania repozytorium i z pominięciem plików
// binarnych oraz nadmiernie dużych.
func PlikiDoPrzejrzenia(katalog string, wlacz []string) ([]string, error) {
	korzen, pomijaj, err := obszarSzukania(katalog)
	if err != nil {
		return nil, err
	}
	sciezki := make([]string, 0, 256)
	err = przejdzPliki(korzen, pomijaj, wlacz, nil, func(wzgledna, _ string) error {
		sciezki = append(sciezki, filepath.Join(korzen, filepath.FromSlash(wzgledna)))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(sciezki)
	return sciezki, nil
}
