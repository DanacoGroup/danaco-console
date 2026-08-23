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

// Wyszukiwanie i zamiana w całym repozytorium.
//
// Silnikiem jest `regexp` biblioteki standardowej, a nie `ripgrep`: wyszukiwanie
// globalne jest czynnością, bez której Code Editor przestaje być edytorem kodu,
// więc nie może zależeć od programu, którego instalka nie niesie. RE2 nie ma
// wstecznych odwołań, i to jest jedyna różnica, którą Operator zobaczy —
// wyrażenie z `\1` dostanie odmowę nazwaną, a nie ciche zero trafień.
//
// Pliki ignorowane przez repozytorium są pomijane. Wyszukiwanie, które wchodzi
// w `node_modules` i katalog budowania, oddaje tysiące trafień w kodzie, którego
// Operator nie pisał — i jest wtedy bezużyteczne, choć formalnie poprawne.

// ZadanieSzukania opisuje, czego i gdzie szukać.
type ZadanieSzukania struct {
	Wzorzec            string
	Wyrazenie          bool
	RozrozniajWielkosc bool
	Wlacz              []string
	Wylacz             []string
	Granica            int
}

// GranicaDomyslna jest górną granicą trafień, gdy żądanie jej nie podaje.
//
// Granica istnieje, bo odpowiedź idzie kopertą przez gniazdo: wyszukanie litery
// „e" w dużym repozytorium dałoby zbiór, którego klient nie postawi na ekranie,
// a rdzeń trzymałby go w pamięci w całości.
const GranicaDomyslna = 2000

// GranicaPlikuBajtow pomija pliki większe od tej granicy.
//
// Plik stumegabajtowy w repozytorium jest zwykle zrzutem danych albo artefaktem
// budowania, a nie kodem; jego przeszukanie kosztuje czas, którego wynik i tak
// nikomu nie służy.
const GranicaPlikuBajtow = 4 << 20

// Szukaj przechodzi repozytorium i oddaje trafienia w kolejności plików.
//
// Trzecia wartość mówi, czy wynik przycięto granicą; druga — ile trafień
// znaleziono przed przycięciem. Sam wykaz przycięty bez tej liczby kazałby
// Operatorowi zgadywać, czy zobaczył wszystko.
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

// ZadanieZamiany opisuje zamianę w repozytorium.
type ZadanieZamiany struct {
	Wzorzec   string
	Zamiana   string
	Wyrazenie bool
	Sciezki   []string
	Podglad   bool
}

// Zamien wykonuje zamianę wzorca w repozytorium.
//
// Podgląd jest stanem domyślnym pracy eksperckiej: zamiana masowa dotyka wielu
// plików naraz i „Cofnij" po niej nie istnieje. Dlatego czynność oddaje wykaz
// zmian razem z odpowiedzią o tym, czy je zapisała — a nie samo słowo
// „wykonano", z którego Operator nie wyczyta, co się stało z jego kodem.
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
				// Wzorzec dosłowny zamienia się dosłownie: `$1` w treści
				// wstawianej ma zostać `$1`, a nie zniknąć jako odwołanie do
				// grupy, której w takim wzorcu nie ma.
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

// obszarSzukania oddaje korzeń przeszukania oraz regułę pomijania.
//
// Katalog bez repozytorium też podlega przeszukaniu: Operator otwiera w oknie
// także katalogi, których nie wersjonuje, a odmowa wyszukiwania byłaby wtedy
// odmową bez powodu. Reguły pomijania są wtedy puste poza katalogiem `.git`.
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

// regulyPomijania czyta `.gitignore` korzenia repozytorium.
//
// Czytany jest plik korzenia, nie wszystkie pliki drzewa: to on niesie reguły
// katalogów budowania i zależności, czyli te, których pominięcie decyduje
// o użyteczności wyniku. Reguły podkatalogów zawężają wynik dodatkowo i ich brak
// oznacza wyłącznie kilka trafień więcej, nigdy mniej.
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

// przejdzPliki woła obsługę dla każdego pliku tekstowego obszaru.
func przejdzPliki(korzen string, pomijaj gitignore.Matcher, wlacz, wylacz []string,
	obsluz func(wzgledna, tresc string) error) error {

	return filepath.WalkDir(korzen, func(sciezka string, wpis fs.DirEntry, err error) error {
		if err != nil {
			// Katalog bez prawa odczytu nie zatrzymuje przeszukania całości:
			// wynik ma być niepełny o ten katalog, a nie pusty.
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
// i poza wyłączającymi.
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

// pasujeGlob dopasowuje wzorzec do ścieżki albo do samej nazwy pliku.
//
// Dopasowanie idzie dwutorowo, bo Operator pisze i `*.go`, i `server/**`.
// `filepath.Match` nie zna `**`, więc wzorzec z podwójną gwiazdką sprowadza się
// do przedrostka ścieżki — to jest znaczenie, którego Operator się spodziewa.
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

// PlikiDoPrzejrzenia oddaje ścieżki bezwzględne plików tekstowych obszaru —
// z poszanowaniem reguł pomijania repozytorium i z pominięciem plików
// binarnych oraz nadmiernie dużych.
//
// Wystawione poza pakiet dla skanowania bezpieczeństwa i jakości: skan przechodzi
// dokładnie ten sam zbiór plików, co wyszukiwanie, więc jedno przejście ma
// stanowić jedną prawdę o tym, co należy do repozytorium. Skan zgłaszający
// sekret w katalogu pobranych zależności byłby wykazem, którego nikt nie czyta.
//
// Wzorce włączające zawężają wynik tak samo, jak w `Szukaj`; pusty wykaz znaczy
// całe drzewo.
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
