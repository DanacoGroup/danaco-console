// Plik obsługuje cztery komendy odczytu okna Build Output: listę przebiegów
// budowania, log przebiegu, wynik testów oraz pokrycie kodu, wraz z rozbiorem
// wyjścia testów i profilu pokrycia.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// najwiecejPrzebiegowWykazu jest domyślną głębokością historii budowań. Build
// Output pokazuje listę, a nie archiwum — pełny dziennik idzie po wyraźny limit.
const najwiecejPrzebiegowWykazu = 50

// WykazBudowan obsługuje `developer.build.list`.
//
// Przebieg czynny w tej chwili wchodzi do wykazu z pamięci rdzenia, a nie
// z dziennika: w dzienniku ma jeszcze stan sprzed pierwszego wiersza logu,
// a Operator ma zobaczyć to, co dzieje się teraz.
func (a *adapterDevelopera) WykazBudowan(ctx context.Context,
	z shared.DeveloperBuildListRequest) (shared.DeveloperBuildListResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperBuildListResponse{}, err
	}
	if a.repozytorium == nil {
		return shared.DeveloperBuildListResponse{}, bladZasobuDevelopera(
			"serwer nie ma dziennika przebiegów budowania")
	}

	stan := ""
	if z.Status != nil {
		stan = string(*z.Status)
	}
	limit := najwiecejPrzebiegowWykazu
	if z.Limit != nil && *z.Limit > 0 {
		limit = *z.Limit
	}
	wiersze, err := a.repozytorium.Przebiegi(ctx, okno.Id, stan, limit)
	if err != nil {
		return shared.DeveloperBuildListResponse{}, bladWykonaniaDevelopera(
			"nie można odczytać historii budowań okna " + okno.Id + ": " + err.Error())
	}

	czynny, jestCzynny := a.rejestr.Przebieg(okno.Id)
	budowania := make([]shared.DeveloperBuild, 0, len(wiersze))
	for _, wiersz := range wiersze {
		if jestCzynny && wiersz.Kod == czynny.kod {
			budowania = append(budowania, budowanieKontraktu(czynny))
			continue
		}
		budowania = append(budowania, budowanieZDziennika(wiersz))
	}
	razem := len(budowania)
	return shared.DeveloperBuildListResponse{Builds: budowania, Total: &razem}, nil
}

// LogBudowania obsługuje komendę odczytu logu przebiegu budowania.
//
// Przebieg czynny czyta się z pamięci, domknięty — z dziennika. Odpowiedź mówi
// wprost, czy log jest przycięty i od którego wiersza, żeby okno nie musiało
// zgadywać.
func (a *adapterDevelopera) LogBudowania(ctx context.Context,
	z shared.DeveloperBuildLogGetRequest) (shared.DeveloperBuildLogGetResponse, error) {

	kod := strings.TrimSpace(z.BuildId)
	if kod == "" {
		return shared.DeveloperBuildLogGetResponse{}, bladZadaniaDevelopera(
			"odczyt logu wymaga wskazania przebiegu budowania")
	}

	log, przyciety, err := a.trescLoguPrzebiegu(ctx, kod)
	if err != nil {
		return shared.DeveloperBuildLogGetResponse{}, err
	}
	wiersze := []string{}
	if log != "" {
		wiersze = strings.Split(log, "\n")
	}

	odpowiedz := shared.DeveloperBuildLogGetResponse{Truncated: przyciety}
	if z.FromLine != nil && *z.FromLine > 0 {
		od := *z.FromLine - 1
		if od >= len(wiersze) {
			wiersze = []string{}
		} else {
			wiersze = wiersze[od:]
		}
		odpowiedz.TruncatedFrom = z.FromLine
	}
	if z.Tail != nil && *z.Tail > 0 && len(wiersze) > *z.Tail {
		odpowiedz.TruncatedFrom = wskaznikLiczby(len(wiersze) - *z.Tail + 1)
		wiersze = wiersze[len(wiersze)-*z.Tail:]
		odpowiedz.Truncated = true
	}
	odpowiedz.Lines = wiersze
	return odpowiedz, nil
}

// trescLoguPrzebiegu oddaje log przebiegu wraz z informacją o przycięciu:
// sięga do pamięci rdzenia dla przebiegu czynnego, a do dziennika dla
// przebiegu już domkniętego.
func (a *adapterDevelopera) trescLoguPrzebiegu(ctx context.Context,
	kod string) (string, bool, error) {

	if przebieg, jest := a.rejestr.PrzebiegPoKodzie(kod); jest {
		return przebieg.Ogon(), przebieg.OgonPrzyciety(), nil
	}
	if a.repozytorium == nil {
		return "", false, bladZasobuDevelopera(
			"przebieg " + kod + " nie biegnie, a serwer nie ma dziennika przebiegów")
	}
	wiersz, err := a.repozytorium.Przebieg(ctx, kod)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return "", false, bladZasobuDevelopera("nie ma przebiegu budowania " + kod)
	}
	if err != nil {
		return "", false, bladWykonaniaDevelopera(
			"nie można odczytać przebiegu " + kod + ": " + err.Error())
	}
	if wiersz.Log == nil {
		return "", false, nil
	}
	// Dziennik trzyma ogon, nie całość: pusty znacznik przycięcia znaczyłby,
	// że to jest cały log.
	return *wiersz.Log, true, nil
}

// WynikTestow obsługuje komendę odczytu wyniku testów przebiegu budowania
// i liczy podsumowanie po wszystkich wynikach, niezależnie od zawężenia
// odpowiedzi stanem testu.
func (a *adapterDevelopera) WynikTestow(ctx context.Context,
	z shared.DeveloperTestResultGetRequest) (shared.DeveloperTestResultGetResponse, error) {

	kod := strings.TrimSpace(z.BuildId)
	if kod == "" {
		return shared.DeveloperTestResultGetResponse{}, bladZadaniaDevelopera(
			"odczyt wyniku testów wymaga wskazania przebiegu budowania")
	}
	if a.repozytorium == nil {
		return shared.DeveloperTestResultGetResponse{}, bladZasobuDevelopera(
			"serwer nie ma dziennika przebiegów, więc nie ma skąd wziąć wyniku testów")
	}

	stan := ""
	if z.Status != nil {
		stan = string(*z.Status)
	}
	// Podsumowanie liczy się po wszystkich wynikach, nie po zawężonych
	// stanem.
	wszystkie, err := a.repozytorium.WynikiTestow(ctx, kod, "")
	if err != nil {
		return shared.DeveloperTestResultGetResponse{}, bladWykonaniaDevelopera(
			"nie można odczytać wyników testów przebiegu " + kod + ": " + err.Error())
	}

	wyniki := make([]shared.DeveloperTestResult, 0, len(wszystkie))
	var przeszly, nieprzeszly, pominiete int
	var czasRazem int64
	for _, wiersz := range wszystkie {
		switch wiersz.Stan {
		case shared.TestStatusPassed:
			przeszly++
		case shared.TestStatusFailed:
			nieprzeszly++
		default:
			pominiete++
		}
		if wiersz.CzasMs != nil {
			czasRazem += *wiersz.CzasMs
		}
		if stan != "" && string(wiersz.Stan) != stan {
			continue
		}
		wyniki = append(wyniki, shared.DeveloperTestResult{
			Suite:      wiersz.Zestaw,
			Name:       wiersz.Nazwa,
			Status:     wiersz.Stan,
			DurationMs: wiersz.CzasMs,
			Message:    wiersz.Tresc,
			Path:       wiersz.Sciezka,
			Line:       liczbaZDuzej(wiersz.Wiersz),
		})
	}

	odpowiedz := shared.DeveloperTestResultGetResponse{
		Results: wyniki,
		Passed:  przeszly,
		Failed:  nieprzeszly,
		Skipped: pominiete,
	}
	if czasRazem > 0 {
		odpowiedz.DurationMs = &czasRazem
	}
	return odpowiedz, nil
}

// Pokrycie obsługuje komendę odczytu pokrycia kodu przebiegu budowania i liczy
// procent zbiorczy po całym przebiegu, niezależnie od zawężenia odpowiedzi do
// jednego pliku.
func (a *adapterDevelopera) Pokrycie(ctx context.Context,
	z shared.DeveloperCoverageGetRequest) (shared.DeveloperCoverageGetResponse, error) {

	kod := strings.TrimSpace(z.BuildId)
	if kod == "" {
		return shared.DeveloperCoverageGetResponse{}, bladZadaniaDevelopera(
			"odczyt pokrycia wymaga wskazania przebiegu budowania")
	}
	if a.repozytorium == nil {
		return shared.DeveloperCoverageGetResponse{}, bladZasobuDevelopera(
			"serwer nie ma dziennika przebiegów, więc nie ma skąd wziąć pokrycia")
	}

	sciezka := ""
	if z.Path != nil {
		sciezka = strings.TrimSpace(*z.Path)
	}
	// Procent zbiorczy liczy się po całym przebiegu, nawet gdy odpowiedź jest
	// zawężona do jednego pliku.
	wszystkie, err := a.repozytorium.Pokrycie(ctx, kod, "")
	if err != nil {
		return shared.DeveloperCoverageGetResponse{}, bladWykonaniaDevelopera(
			"nie można odczytać pokrycia przebiegu " + kod + ": " + err.Error())
	}

	pliki := make([]shared.DeveloperCoverage, 0, len(wszystkie))
	var instrukcje, pokryte int64
	for _, wiersz := range wszystkie {
		instrukcje += wiersz.Instrukcje
		pokryte += wiersz.Pokryte
		if sciezka != "" && wiersz.Sciezka != sciezka {
			continue
		}
		pliki = append(pliki, shared.DeveloperCoverage{
			Path:           wiersz.Sciezka,
			Statements:     int(wiersz.Instrukcje),
			Covered:        int(wiersz.Pokryte),
			Percent:        int(wiersz.Procent),
			UncoveredLines: wierszeZTekstu(wiersz.WierszeBezPokrycia),
		})
	}

	procent := 0
	if instrukcje > 0 {
		procent = int((pokryte * 100) / instrukcje)
	} else if len(wszystkie) > 0 {
		// Pomiar bez profilu nie zna instrukcji, więc średnia po plikach jest
		// jedyną dostępną wartością.
		var suma int64
		for _, wiersz := range wszystkie {
			suma += wiersz.Procent
		}
		procent = int(suma / int64(len(wszystkie)))
	}
	return shared.DeveloperCoverageGetResponse{Files: pliki, Percent: procent}, nil
}

// wierszeZTekstu rozbiera zapis numerów wierszy rozdzielonych przecinkiem na
// wykaz liczb całkowitych, pomijając człony, które liczbą nie są.
func wierszeZTekstu(zapis *string) []int {
	if zapis == nil || strings.TrimSpace(*zapis) == "" {
		return nil
	}
	czesci := strings.Split(*zapis, ",")
	wiersze := make([]int, 0, len(czesci))
	for _, czesc := range czesci {
		if numer, err := strconv.Atoi(strings.TrimSpace(czesc)); err == nil {
			wiersze = append(wiersze, numer)
		}
	}
	return wiersze
}

// liczbaZDuzej sprowadza wskaźnik liczby długiej do wskaźnika liczby
// kontraktu, zachowując pustą wartość, gdy wskaźnik źródłowy jest pusty.
func liczbaZDuzej(wartosc *int64) *int {
	if wartosc == nil {
		return nil
	}
	liczba := int(*wartosc)
	return &liczba
}

// ── Rozbiór wyjścia testów ───────────────────────────────────────────────────

// zdarzenieTestuGo jest kształtem jednego wiersza strumienia zdarzeń
// narzędzia testowego w postaci maszynowej, niosącym czynność, pakiet, nazwę
// testu i wyjście.
type zdarzenieTestuGo struct {
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Elapsed float64 `json:"Elapsed"`
	Output  string  `json:"Output"`
}

// czyWierszTestu rozpoznaje wiersz, który niesie wynik testu albo pokrycie.
//
// Rozpoznanie jest jawnym sitem, a nie zachowaniem całego logu: wyjście
// kompilatora bywa dużo obszerniejsze od wyjścia testów.
func czyWierszTestu(wiersz string) bool {
	tresc := strings.TrimSpace(wiersz)
	switch {
	case strings.HasPrefix(tresc, `{"Time":`), strings.HasPrefix(tresc, `{"Action":`):
		return true
	case strings.HasPrefix(tresc, "--- PASS:"), strings.HasPrefix(tresc, "--- FAIL:"),
		strings.HasPrefix(tresc, "--- SKIP:"):
		return true
	case strings.Contains(tresc, "coverage:"):
		return true
	default:
		return false
	}
}

// wynikiTestowZLogu składa wyniki testów z zebranych wierszy.
//
// Najpierw sprawdza się postać maszynowa, bo niesie czas, pakiet i treść
// niepowodzenia; gdy jej nie ma, zostaje postać czytelna, która niesie mniej,
// lecz nadal niesie wynik.
func wynikiTestowZLogu(wiersze []string) []dane.WynikTestu {
	jesliMaszynowe := wynikiTestowZPostaciMaszynowej(wiersze)
	if len(jesliMaszynowe) > 0 {
		return jesliMaszynowe
	}
	return wynikiTestowZPostaciCzytelnej(wiersze)
}

// wynikiTestowZPostaciMaszynowej rozbiera strumień zdarzeń testowych
// w postaci maszynowej i składa z nich wyniki testów wraz z czasem trwania
// i treścią niepowodzenia.
func wynikiTestowZPostaciMaszynowej(wiersze []string) []dane.WynikTestu {
	wyniki := make([]dane.WynikTestu, 0, 32)
	// Treść niepowodzenia przychodzi zdarzeniami wyjścia przed zdarzeniem
	// porażki.
	tresci := map[string]*strings.Builder{}

	for _, wiersz := range wiersze {
		tresc := strings.TrimSpace(wiersz)
		if !strings.HasPrefix(tresc, "{") {
			continue
		}
		var zdarzenie zdarzenieTestuGo
		if err := json.Unmarshal([]byte(tresc), &zdarzenie); err != nil || zdarzenie.Test == "" {
			continue
		}
		klucz := zdarzenie.Package + "/" + zdarzenie.Test

		if zdarzenie.Action == "output" {
			budowany, jest := tresci[klucz]
			if !jest {
				budowany = &strings.Builder{}
				tresci[klucz] = budowany
			}
			if budowany.Len() < 4000 {
				budowany.WriteString(zdarzenie.Output)
			}
			continue
		}

		stan, wynikowe := stanTestuZCzynnosci(zdarzenie.Action)
		if !wynikowe {
			continue
		}
		wynik := dane.WynikTestu{Nazwa: zdarzenie.Test, Stan: stan}
		if zdarzenie.Package != "" {
			wynik.Zestaw = wskaznikTekstu(zdarzenie.Package)
		}
		if zdarzenie.Elapsed > 0 {
			czas := int64(zdarzenie.Elapsed * 1000)
			wynik.CzasMs = &czas
		}
		if stan == shared.TestStatusFailed {
			if budowany, jest := tresci[klucz]; jest {
				wynik.Tresc = wskaznikTekstu(strings.TrimSpace(budowany.String()))
				if sciezka, numer, jest := miejsceNiepowodzeniaTestu(budowany.String()); jest {
					wynik.Sciezka, wynik.Wiersz = wskaznikTekstu(sciezka), &numer
				}
			}
		}
		wyniki = append(wyniki, wynik)
	}
	return wyniki
}

// stanTestuZCzynnosci przekłada czynność zdarzenia na stan testu. Czynności
// pośrednie (`run`, `pause`, `cont`) nie są wynikiem i nie zakładają wiersza.
func stanTestuZCzynnosci(czynnosc string) (shared.TestStatus, bool) {
	switch czynnosc {
	case "pass":
		return shared.TestStatusPassed, true
	case "fail":
		return shared.TestStatusFailed, true
	case "skip":
		return shared.TestStatusSkipped, true
	default:
		return "", false
	}
}

// wynikiTestowZPostaciCzytelnej rozbiera wiersze dziennika testów w postaci
// czytelnej dla człowieka i składa z nich wyniki testów wraz z czasem
// trwania, gdy jest podany.
func wynikiTestowZPostaciCzytelnej(wiersze []string) []dane.WynikTestu {
	wyniki := make([]dane.WynikTestu, 0, 32)
	for _, wiersz := range wiersze {
		tresc := strings.TrimSpace(wiersz)
		var stan shared.TestStatus
		switch {
		case strings.HasPrefix(tresc, "--- PASS:"):
			stan = shared.TestStatusPassed
		case strings.HasPrefix(tresc, "--- FAIL:"):
			stan = shared.TestStatusFailed
		case strings.HasPrefix(tresc, "--- SKIP:"):
			stan = shared.TestStatusSkipped
		default:
			continue
		}
		pola := strings.Fields(tresc)
		if len(pola) < 3 {
			continue
		}
		wynik := dane.WynikTestu{Nazwa: pola[2], Stan: stan}
		if len(pola) >= 4 {
			if czas, jest := czasTestuZNawiasu(pola[3]); jest {
				wynik.CzasMs = &czas
			}
		}
		wyniki = append(wyniki, wynik)
	}
	return wyniki
}

// czasTestuZNawiasu rozbiera zapis czasu testu w nawiasie, wyrażony
// w sekundach, na liczbę całkowitą milisekund.
func czasTestuZNawiasu(zapis string) (int64, bool) {
	tresc := strings.TrimSuffix(strings.TrimPrefix(zapis, "("), ")")
	tresc = strings.TrimSuffix(tresc, "s")
	sekundy, err := strconv.ParseFloat(tresc, 64)
	if err != nil {
		return 0, false
	}
	return int64(sekundy * 1000), true
}

// miejsceNiepowodzeniaTestu wyłuskuje z treści niepowodzenia pierwsze wskazanie
// `plik.go:123`. To ono pozwala Build Output przejść z komunikatu do wiersza
// kodu — bez niego Operator szuka miejsca ręcznie.
func miejsceNiepowodzeniaTestu(tresc string) (string, int64, bool) {
	for _, wiersz := range strings.Split(tresc, "\n") {
		pole := strings.TrimSpace(wiersz)
		dwukropek := strings.LastIndex(pole, ".go:")
		if dwukropek < 0 {
			continue
		}
		reszta := pole[dwukropek+4:]
		koniec := strings.IndexFunc(reszta, func(znak rune) bool {
			return znak < '0' || znak > '9'
		})
		if koniec == 0 {
			continue
		}
		if koniec < 0 {
			koniec = len(reszta)
		}
		numer, err := strconv.ParseInt(reszta[:koniec], 10, 64)
		if err != nil {
			continue
		}
		poczatek := strings.LastIndexAny(pole[:dwukropek], " \t")
		return strings.TrimSpace(pole[poczatek+1 : dwukropek+3]), numer, true
	}
	return "", 0, false
}

// ── Rozbiór pokrycia ─────────────────────────────────────────────────────────

// pokrycieZPrzebiegu składa pomiar pokrycia z profilu, gdy jest dostępny,
// a w przeciwnym razie z wierszy logu niosących podsumowanie pokrycia
// pakietu.
func pokrycieZPrzebiegu(katalog string, argumenty, wiersze []string) []dane.PokryciePliku {
	if profil, jest := profilPokryciaZArgumentow(katalog, argumenty); jest {
		if pomiar := pokrycieZProfilu(profil); len(pomiar) > 0 {
			return pomiar
		}
	}
	return pokrycieZWierszyLogu(wiersze)
}

// profilPokryciaZArgumentow odnajduje w parametrach zadania budowania ścieżkę
// profilu pokrycia i sprowadza ją do ścieżki bezwzględnej względem katalogu
// zadania.
func profilPokryciaZArgumentow(katalog string, argumenty []string) (string, bool) {
	for i, argument := range argumenty {
		wartosc := ""
		switch {
		case strings.HasPrefix(argument, "-coverprofile="):
			wartosc = strings.TrimPrefix(argument, "-coverprofile=")
		case strings.HasPrefix(argument, "--coverprofile="):
			wartosc = strings.TrimPrefix(argument, "--coverprofile=")
		case argument == "-coverprofile" || argument == "--coverprofile":
			if i+1 < len(argumenty) {
				wartosc = argumenty[i+1]
			}
		}
		if wartosc == "" {
			continue
		}
		if !filepath.IsAbs(wartosc) {
			wartosc = filepath.Join(katalog, wartosc)
		}
		return wartosc, true
	}
	return "", false
}

// pokrycieZProfilu czyta profil pokrycia Go.
//
// Wiersz profilu ma postać `plik:od.kol,do.kol liczbaInstrukcji liczbaTrafien`.
// Instrukcje sumują się po pliku, a bloki z zerem trafień dają wykaz wierszy
// niepokrytych — to on zasila nakładkę w edytorze.
func pokrycieZProfilu(sciezka string) []dane.PokryciePliku {
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		return nil
	}
	type licznik struct {
		instrukcje, pokryte int64
		bezPokrycia         []int
	}
	pomiary := map[string]*licznik{}
	kolejnosc := make([]string, 0, 16)

	for _, wiersz := range strings.Split(string(bajty), "\n") {
		tresc := strings.TrimSpace(wiersz)
		if tresc == "" || strings.HasPrefix(tresc, "mode:") {
			continue
		}
		pola := strings.Fields(tresc)
		if len(pola) != 3 {
			continue
		}
		dwukropek := strings.LastIndexByte(pola[0], ':')
		if dwukropek < 0 {
			continue
		}
		plik, zakres := pola[0][:dwukropek], pola[0][dwukropek+1:]
		instrukcje, err1 := strconv.ParseInt(pola[1], 10, 64)
		trafienia, err2 := strconv.ParseInt(pola[2], 10, 64)
		if err1 != nil || err2 != nil {
			continue
		}
		pomiar, jest := pomiary[plik]
		if !jest {
			pomiar = &licznik{}
			pomiary[plik] = pomiar
			kolejnosc = append(kolejnosc, plik)
		}
		pomiar.instrukcje += instrukcje
		if trafienia > 0 {
			pomiar.pokryte += instrukcje
		} else {
			pomiar.bezPokrycia = append(pomiar.bezPokrycia, wierszeZakresuProfilu(zakres)...)
		}
	}

	sort.Strings(kolejnosc)
	pokrycie := make([]dane.PokryciePliku, 0, len(kolejnosc))
	for _, plik := range kolejnosc {
		pomiar := pomiary[plik]
		procent := int64(0)
		if pomiar.instrukcje > 0 {
			procent = (pomiar.pokryte * 100) / pomiar.instrukcje
		}
		wiersz := dane.PokryciePliku{
			Sciezka:    plik,
			Instrukcje: pomiar.instrukcje,
			Pokryte:    pomiar.pokryte,
			Procent:    procent,
		}
		if len(pomiar.bezPokrycia) > 0 {
			wiersz.WierszeBezPokrycia = wskaznikTekstu(zapisWierszy(pomiar.bezPokrycia))
		}
		pokrycie = append(pokrycie, wiersz)
	}
	return pokrycie
}

// wierszeZakresuProfilu rozbiera zapis zakresu bloku profilu pokrycia na pełny
// wykaz numerów wierszy, które ten blok obejmuje.
func wierszeZakresuProfilu(zakres string) []int {
	czesci := strings.SplitN(zakres, ",", 2)
	if len(czesci) != 2 {
		return nil
	}
	od, err1 := strconv.Atoi(strings.SplitN(czesci[0], ".", 2)[0])
	doWiersza, err2 := strconv.Atoi(strings.SplitN(czesci[1], ".", 2)[0])
	if err1 != nil || err2 != nil || doWiersza < od {
		return nil
	}
	// Blok bez pokrycia bywa długi, lecz nakładka edytora potrzebuje każdego
	// jego wiersza.
	wiersze := make([]int, 0, doWiersza-od+1)
	for numer := od; numer <= doWiersza; numer++ {
		wiersze = append(wiersze, numer)
	}
	return wiersze
}

// zapisWierszy składa numery wierszy niepokrytych w jeden tekst kolumny bazy,
// rozdzielając je przecinkiem, w postaci odwrotnej do wierszeZTekstu.
func zapisWierszy(wiersze []int) string {
	zapis := make([]string, 0, len(wiersze))
	for _, numer := range wiersze {
		zapis = append(zapis, strconv.Itoa(numer))
	}
	return strings.Join(zapis, ",")
}

// pokrycieZWierszyLogu czyta wiersze podsumowania testu niosące procent
// pokrycia pakietu.
//
// Pomiar bez profilu zna procent, a nie zna instrukcji ani wierszy, więc
// zapisuje się sam procent, a instrukcje zostają zerowe.
func pokrycieZWierszyLogu(wiersze []string) []dane.PokryciePliku {
	pokrycie := make([]dane.PokryciePliku, 0, 8)
	widziane := map[string]bool{}
	for _, wiersz := range wiersze {
		miejsce := strings.Index(wiersz, "coverage:")
		if miejsce < 0 {
			continue
		}
		reszta := strings.TrimSpace(wiersz[miejsce+len("coverage:"):])
		procent, jest := procentZZapisu(reszta)
		if !jest {
			continue
		}
		pakiet := pakietZWierszaPokrycia(wiersz[:miejsce])
		if pakiet == "" || widziane[pakiet] {
			continue
		}
		widziane[pakiet] = true
		pokrycie = append(pokrycie, dane.PokryciePliku{Sciezka: pakiet, Procent: procent})
	}
	return pokrycie
}

// procentZZapisu rozbiera zapis procentu pokrycia pakietu na liczbę całkowitą,
// odrzucając ułamkową część wartości.
func procentZZapisu(zapis string) (int64, bool) {
	pola := strings.Fields(zapis)
	if len(pola) == 0 {
		return 0, false
	}
	liczba, err := strconv.ParseFloat(strings.TrimSuffix(pola[0], "%"), 64)
	if err != nil {
		return 0, false
	}
	return int64(liczba), true
}

// pakietZWierszaPokrycia wyjmuje nazwę pakietu z początku wiersza
// podsumowania testu, pomijając wynik przebiegu i czas jego trwania.
func pakietZWierszaPokrycia(poczatek string) string {
	pola := strings.Fields(poczatek)
	for _, pole := range pola {
		if pole == "ok" || pole == "?" || strings.HasSuffix(pole, "s") &&
			strings.ContainsAny(pole, "0123456789") {
			continue
		}
		return pole
	}
	return ""
}

// odlozPomiarPrzebiegu zapisuje wynik testów i pokrycie domkniętego
// przebiegu, wywołane raz przy domknięciu, kiedy pełne wyjście jeszcze jest
// w zasięgu rdzenia; niepowodzenie zapisu nie zatrzymuje niczego, bo wynik
// budowania jest już znany.
func (a *adapterDevelopera) odlozPomiarPrzebiegu(przebieg *przebiegBudowania) {
	wiersze := przebieg.WierszeTestow()
	if len(wiersze) == 0 {
		return
	}
	kontekst := context.Background()

	if wyniki := wynikiTestowZLogu(wiersze); len(wyniki) > 0 {
		_ = a.repozytorium.ZapiszWynikiTestow(kontekst, przebieg.kod, wyniki)
	}

	katalog := ""
	if okno, err := a.oknoDevelopera(przebieg.oknoKod); err == nil {
		if korzenie := a.korzenieOkna(okno); len(korzenie) > 0 {
			katalog = korzenie[0]
		}
	}
	if pomiar := pokrycieZPrzebiegu(katalog, przebieg.argumenty, wiersze); len(pomiar) > 0 {
		_ = a.repozytorium.ZapiszPokrycie(kontekst, przebieg.kod, pomiar)
	}
}
