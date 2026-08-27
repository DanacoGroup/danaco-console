// Pakiet obsługuje silniki zewnętrzne korekty językowej modułu Translate:
// LanguageTool, hunspell i Vale, po które sięga `translate.proofread.run`,
// wraz z przekładem ich wyjścia na ustalenia kontraktu tej komendy.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

const (
	// granicaKorektyJezykowej obejmuje start maszyny wirtualnej Javy wraz
	// z wczytaniem reguł i słownika morfologicznego LanguageToola.
	granicaKorektyJezykowej = 3 * time.Minute
	// granicaKorektyPisowni jest krótsza, bo hunspell i vale są programami
	// natywnymi bez własnego środowiska uruchomieniowego.
	granicaKorektyPisowni = 60 * time.Second

	// zmiennaLanguageToola wskazuje katalog wydania LanguageToola —
	// pierwszeństwo przed miejscem typowym, tą samą zasadą co `DANACO_PIPER`
	// przy syntezie mowy.
	zmiennaLanguageToola       = "DANACO_LANGUAGETOOL"
	katalogLanguageToolaTypowy = "/opt/languagetool"
	// archiwumLanguageToola to nazwa archiwum wiersza poleceń. Wydanie niesie
	// obok niego archiwum serwera i archiwum biblioteki — nazwa rozstrzyga,
	// które z trzech uruchamiamy.
	archiwumLanguageToola = "languagetool-commandline.jar"

	// nazwaKonfiguracjiVale jest nazwą pliku, bez którego program vale nie
	// rusza; rdzeń zakłada go sam w katalogu roboczym czynności.
	nazwaKonfiguracjiVale = ".vale.ini"

	// nazwaTresciKorekty jest nazwą pliku, którym treść panelu jedzie do
	// programów. Wszystkie trzy czytają PLIK, a port `session.Uruchamiacz` nie
	// wystawia wejścia procesu — tak samo jak przy syntezie mowy.
	nazwaTresciKorekty = "panel.md"
)

// narzedzieLanguageToola opisuje maszynę wirtualną, którą LanguageTool się
// uruchamia — rozdział programu od archiwum jest ten sam, co przy Apache Tice
// (`adapter_narzedzia_dokument_tekst.go`) i z tego samego powodu.
var narzedzieLanguageToola = zewnetrzne.Narzedzie{
	Nazwa:   "LanguageTool (uruchamiany środowiskiem Javy)",
	Program: "java",
	Pakiet: "środowisko uruchomieniowe Javy (default-jre) wraz z wydaniem LanguageToola " +
		"w " + katalogLanguageToolaTypowy + " albo w katalogu wskazanym zmienną " +
		zmiennaLanguageToola,
}

// narzedzieHunspella opisuje słownik ortograficzny. Pakiet wymienia i program,
// i słowniki — sam program bez pliku słownika nie sprawdzi niczego.
var narzedzieHunspella = zewnetrzne.Narzedzie{
	Nazwa:   "hunspell",
	Program: "hunspell",
	Pakiet:  "hunspell wraz ze słownikiem języka (hunspell-pl, hunspell-en-us)",
}

// narzedzieVale opisuje analizator prozy wołany jednym plikiem wykonywalnym,
// bez zależności zewnętrznych do zainstalowania osobno.
var narzedzieVale = zewnetrzne.Narzedzie{
	Nazwa: "vale", Program: "vale", Pakiet: "vale (jeden plik wykonywalny z wydania projektu)",
}

// ustalenieSilnika niesie jedno ustalenie zewnętrznego silnika w postaci,
// którą komenda odkłada w bazie danych czynności korekty.
type ustalenieSilnika struct {
	rodzaj     shared.ProofreadCheckKind
	waga       shared.ProofreadSeverity
	segment    string
	szczegol   string
	propozycja string
}

// ustaleniaSilnikow zbiera ustalenia wszystkich silników zewnętrznych, które
// na tej maszynie stoją; silnik nieobecny nie jest odmową czynności.
func (a *adapterTlumaczenia) ustaleniaSilnikow(ctx context.Context,
	jezyk, tresc string) ([]ustalenieSilnika, error) {

	if a.uruchamiacz == nil {
		// Rdzeń złożony bez warstwy kanału nie ma czym wystartować programu.
		return nil, nil
	}
	// Oba silniki pisowni pracują tylko wtedy, gdy język panelu ma oznaczenie.
	oznaczenie := jezykKorekty(jezyk)
	stoiLanguageTool := oznaczenie != "" &&
		zewnetrzne.Stoi(narzedzieLanguageToola) && archiwumKorekty() != ""
	stoiHunspell := oznaczenie != "" && zewnetrzne.Stoi(narzedzieHunspella)
	stoiVale := zewnetrzne.Stoi(narzedzieVale)
	if !stoiLanguageTool && !stoiHunspell && !stoiVale {
		return nil, nil
	}

	katalog, err := os.MkdirTemp("", "danaco-korekta-")
	if err != nil {
		return nil, bladKorektyZewnetrznej("nie można założyć katalogu roboczego korekty: " + err.Error())
	}
	defer func() { _ = os.RemoveAll(katalog) }()

	plikTresci := filepath.Join(katalog, nazwaTresciKorekty)
	if err := os.WriteFile(plikTresci, []byte(tresc), 0o600); err != nil {
		return nil, bladKorektyZewnetrznej("nie można odłożyć treści panelu do korekty: " + err.Error())
	}

	var ustalenia []ustalenieSilnika
	switch {
	case stoiLanguageTool:
		zLanguageToola, err := a.korektaLanguageToolem(ctx, katalog, plikTresci, oznaczenie, tresc)
		if err != nil {
			return nil, err
		}
		ustalenia = append(ustalenia, zLanguageToola...)
	case stoiHunspell:
		zHunspella, err := a.korektaHunspellem(ctx, katalog, plikTresci, oznaczenie, tresc)
		if err != nil {
			return nil, err
		}
		ustalenia = append(ustalenia, zHunspella...)
	}
	if stoiVale {
		zVale, err := a.korektaVale(ctx, katalog, tresc)
		if err != nil {
			return nil, err
		}
		ustalenia = append(ustalenia, zVale...)
	}
	return ustalenia, nil
}

// ── LanguageTool ────────────────────────────────────────────────────────────

// odpowiedzLanguageToola opisuje tę część wyjścia `--json`, którą rdzeń czyta.
// Pola nieodczytane są pominięte celowo — wykaz przepisany w całości
// rozjeżdżałby się z każdym wydaniem programu, a nic z niego nie wynika.
type odpowiedzLanguageToola struct {
	Matches []struct {
		Message string `json:"message"`
		Offset  int    `json:"offset"`
		Length  int    `json:"length"`
		Rule    struct {
			ID        string `json:"id"`
			IssueType string `json:"issueType"`
			Category  struct {
				ID string `json:"id"`
			} `json:"category"`
		} `json:"rule"`
		Replacements []struct {
			Value string `json:"value"`
		} `json:"replacements"`
	} `json:"matches"`
}

// korektaLanguageToolem przeprowadza jeden przebieg programu LanguageTool
// i przekłada jego wynik JSON na ustalenia silnika.
func (a *adapterTlumaczenia) korektaLanguageToolem(ctx context.Context,
	katalog, plikTresci, jezyk, tresc string) ([]ustalenieSilnika, error) {

	wyjscie, err := a.wolajKorekte(ctx, katalog, narzedzieLanguageToola, []string{
		"-jar", archiwumKorekty(), "--language", jezyk,
		"--encoding", "UTF-8", "--json", plikTresci,
	}, granicaKorektyJezykowej)
	if err != nil {
		return nil, err
	}

	var odpowiedz odpowiedzLanguageToola
	if err := json.Unmarshal(wyjscie, &odpowiedz); err != nil {
		return nil, bladKorektyZewnetrznej("LanguageTool oddał wynik, którego nie da się " +
			"rozłożyć jako JSON: " + err.Error())
	}

	znaki := []rune(tresc)
	ustalenia := make([]ustalenieSilnika, 0, len(odpowiedz.Matches))
	for _, trafienie := range odpowiedz.Matches {
		fragment, jest := fragmentTresci(znaki, trafienie.Offset, trafienie.Length)
		if !jest {
			// Położenie, którego nie da się potwierdzić w treści, jest brakiem pomiaru.
			continue
		}
		propozycja := ""
		if len(trafienie.Replacements) > 0 {
			propozycja = string(znaki[:trafienie.Offset]) +
				trafienie.Replacements[0].Value +
				string(znaki[trafienie.Offset+trafienie.Length:])
		}
		ustalenia = append(ustalenia, ustalenieSilnika{
			rodzaj:     rodzajKorektyLanguageToola(trafienie.Rule.Category.ID),
			waga:       wagaKorektyLanguageToola(trafienie.Rule.IssueType),
			segment:    fragment,
			szczegol:   trafienie.Message + " (reguła " + trafienie.Rule.ID + ")",
			propozycja: propozycja,
		})
	}
	return ustalenia, nil
}

// rodzajKorektyLanguageToola przekłada kategorię reguły programu na rodzaj
// kontroli kontraktu; kategoria nierozpoznana wraca jako rodzaj gramatyczny.
func rodzajKorektyLanguageToola(kategoria string) shared.ProofreadCheckKind {
	switch strings.ToUpper(strings.TrimSpace(kategoria)) {
	case "TYPOS":
		return shared.ProofreadCheckKindSpelling
	case "PUNCTUATION":
		return shared.ProofreadCheckKindPunctuation
	case "TYPOGRAPHY":
		return shared.ProofreadCheckKindTypography
	case "STYLE", "REDUNDANCY", "PLAIN_ENGLISH", "WORDINESS", "COLLOQUIALISMS":
		return shared.ProofreadCheckKindStyle
	default:
		return shared.ProofreadCheckKindGrammar
	}
}

// wagaKorektyLanguageToola nadaje ustaleniu wagę, której program sam nie
// podaje: błędem jest fakt o słowie, ostrzeżeniem — orzeczenie o zdaniu.
func wagaKorektyLanguageToola(rodzajZgloszenia string) shared.ProofreadSeverity {
	if strings.EqualFold(strings.TrimSpace(rodzajZgloszenia), "misspelling") {
		return shared.ProofreadSeverityError
	}
	return shared.ProofreadSeverityWarning
}

// archiwumKorekty wskazuje archiwum wiersza poleceń LanguageToola albo pustkę,
// gdy wydania nie ma. Szukanie obejmuje katalog wskazany i jeden poziom niżej,
// bo wydanie rozpakowuje się do katalogu z numerem wersji w nazwie.
func archiwumKorekty() string {
	katalog := katalogLanguageToolaTypowy
	if wskazany := strings.TrimSpace(os.Getenv(zmiennaLanguageToola)); wskazany != "" {
		katalog = wskazany
	}
	kandydaci := []string{filepath.Join(katalog, archiwumLanguageToola)}
	zagniezdzone, _ := filepath.Glob(filepath.Join(katalog, "*", archiwumLanguageToola))
	sort.Strings(zagniezdzone)
	kandydaci = append(kandydaci, zagniezdzone...)
	for _, kandydat := range kandydaci {
		if opis, err := os.Stat(kandydat); err == nil && !opis.IsDir() {
			return kandydat
		}
	}
	return ""
}

// jezykKorekty sprowadza język panelu do oznaczenia, którym da się zawołać
// słownik. Pustka znaczy, że rdzeń tego języka nie umie nazwać, nie „domyślny".
func jezykKorekty(jezyk string) string {
	nazwa := strings.TrimSpace(jezyk)
	if czyOznaczenieJezyka(nazwa) {
		return nazwa
	}
	switch strings.ToLower(nazwa) {
	case "polski", "polish":
		return "pl"
	case "angielski", "english":
		return "en"
	}
	return ""
}

// czyOznaczenieJezyka orzeka o kształcie: dwie albo trzy litery, po nich
// dowolna liczba członów po myślniku (`pl`, `pl-PL`, `de-DE-x-simple-language`),
// czyli kształcie oznaczeń, którego oczekuje LanguageTool.
func czyOznaczenieJezyka(napis string) bool {
	if napis == "" {
		return false
	}
	for numer, czlon := range strings.Split(napis, "-") {
		if czlon == "" {
			return false
		}
		if numer == 0 && (len(czlon) < 2 || len(czlon) > 3) {
			return false
		}
		if numer > 0 && len(czlon) > 8 {
			return false
		}
		for _, znak := range czlon {
			litera := znak >= 'a' && znak <= 'z' || znak >= 'A' && znak <= 'Z'
			cyfra := znak >= '0' && znak <= '9'
			if !litera && !(numer > 0 && cyfra) {
				return false
			}
		}
	}
	return true
}

// ── hunspell ────────────────────────────────────────────────────────────────

// korektaHunspellem sprawdza ortografię treści panelu i przekłada wyjście
// w postaci ispellowej programu na ustalenia silnika.
func (a *adapterTlumaczenia) korektaHunspellem(ctx context.Context,
	katalog, plikTresci, jezyk, tresc string) ([]ustalenieSilnika, error) {

	slownik, err := a.slownikHunspella(ctx, katalog, jezyk)
	if err != nil {
		return nil, err
	}
	wyjscie, err := a.wolajKorekte(ctx, katalog, narzedzieHunspella, []string{
		"-d", slownik, "-i", "UTF-8", "-a", plikTresci,
	}, granicaKorektyPisowni)
	if err != nil {
		return nil, err
	}

	ustalenia := []ustalenieSilnika{}
	for _, wiersz := range strings.Split(string(wyjscie), "\n") {
		slowo, propozycje, jest := niepoprawneSlowoHunspella(wiersz)
		if !jest {
			continue
		}
		propozycja := ""
		if len(propozycje) > 0 && strings.Count(tresc, slowo) == 1 {
			// Podmiana wchodzi wyłącznie przy jednym wystąpieniu słowa w treści.
			propozycja = strings.Replace(tresc, slowo, propozycje[0], 1)
		}
		ustalenia = append(ustalenia, ustalenieSilnika{
			rodzaj:     shared.ProofreadCheckKindSpelling,
			waga:       shared.ProofreadSeverityError,
			segment:    slowo,
			szczegol:   "słowa „" + slowo + "” nie ma w słowniku " + slownik,
			propozycja: propozycja,
		})
	}
	return ustalenia, nil
}

// niepoprawneSlowoHunspella czyta jeden wiersz wyjścia ispellowego programu
// i oddaje słowo spoza słownika wraz z propozycjami, gdy takie niesie.
func niepoprawneSlowoHunspella(wiersz string) (string, []string, bool) {
	tresc := strings.TrimRight(wiersz, "\r")
	switch {
	case strings.HasPrefix(tresc, "& "):
		glowa, ogon, znaleziono := strings.Cut(tresc[2:], ": ")
		if !znaleziono {
			return "", nil, false
		}
		pola := strings.Fields(glowa)
		if len(pola) == 0 {
			return "", nil, false
		}
		propozycje := []string{}
		for _, propozycja := range strings.Split(ogon, ", ") {
			if oczyszczona := strings.TrimSpace(propozycja); oczyszczona != "" {
				propozycje = append(propozycje, oczyszczona)
			}
		}
		return pola[0], propozycje, true
	case strings.HasPrefix(tresc, "# "):
		pola := strings.Fields(tresc[2:])
		if len(pola) == 0 {
			return "", nil, false
		}
		return pola[0], nil, true
	default:
		return "", nil, false
	}
}

// slownikHunspella dobiera nazwę słownika do języka panelu pytaniem o stan
// maszyny (`hunspell -D`), nie tabelą wpisaną w kod rdzenia.
func (a *adapterTlumaczenia) slownikHunspella(ctx context.Context,
	katalog, jezyk string) (string, error) {

	oznaczenie := strings.ReplaceAll(strings.TrimSpace(jezyk), "-", "_")
	if oznaczenie == "" {
		return "", bladWskazaniaTlumaczenia("panel nie ma języka — nie ma czym dobrać " +
			"słownika ortograficznego")
	}

	// Plik pusty jest materiałem: bez wskazania pliku `-D` czekałby na wejście.
	pusty := filepath.Join(katalog, "wykaz-slownikow.txt")
	if err := os.WriteFile(pusty, nil, 0o600); err != nil {
		return "", bladKorektyZewnetrznej("nie można założyć pliku pomocniczego: " + err.Error())
	}
	okno, zasady, obszar := a.zasiegProgramowTlumaczenia()
	wynik, blad := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedzieHunspella, []string{"-D", pusty}, katalog, granicaKorektyPisowni)

	// Kod wyjścia tego wywołania nie rozstrzyga niczego i nie jest tu czytany.
	dostepne := slownikiZWykazu(wynik.Diagnostyka)
	if len(dostepne) == 0 {
		powod := "hunspell stoi na tej maszynie, ale nie wypisał ani jednego słownika; " +
			"naprawa: doinstalować pakiet słownika języka panelu (hunspell-pl, hunspell-en-us)"
		if blad != nil {
			powod += ". Diagnostyka wywołania: " + blad.Error()
		}
		return "", bladKorektyZewnetrznej(powod)
	}
	for _, nazwa := range dostepne {
		if strings.EqualFold(nazwa, oznaczenie) {
			return nazwa, nil
		}
	}
	// Oznaczenie bez kraju dopasowuje się po samym języku, wyłącznie gdy kandydat jest jeden.
	podstawa, _, _ := strings.Cut(oznaczenie, "_")
	pasujace := []string{}
	for _, nazwa := range dostepne {
		if czesc, _, _ := strings.Cut(nazwa, "_"); strings.EqualFold(czesc, podstawa) {
			pasujace = append(pasujace, nazwa)
		}
	}
	if len(pasujace) == 1 {
		return pasujace[0], nil
	}
	if len(pasujace) > 1 {
		return "", bladKorektyZewnetrznej("język panelu „" + jezyk + "” nie wskazuje odmiany, " +
			"a hunspell ma ich na tej maszynie kilka: " + strings.Join(pasujace, ", ") +
			"; naprawa: podać język panelu wraz z krajem, na przykład " + pasujace[0])
	}
	return "", bladKorektyZewnetrznej("hunspell nie ma słownika dla języka panelu „" + jezyk +
		"”; słowniki, które stoją na tej maszynie: " + strings.Join(dostepne, ", ") +
		"; naprawa: doinstalować pakiet słownika tego języka")
}

// slownikiZWykazu czyta nazwy słowników z wykazu `hunspell -D`. Wykaz idzie
// diagnostyką programu, a nie wyjściem — stąd czytamy go osobno.
func slownikiZWykazu(wykaz string) []string {
	nazwy := []string{}
	wSekcji := false
	for _, wiersz := range strings.Split(wykaz, "\n") {
		tresc := strings.TrimSpace(wiersz)
		switch {
		case strings.HasPrefix(tresc, "AVAILABLE DICTIONARIES"):
			wSekcji = true
		case strings.HasSuffix(tresc, ":"):
			// Każdy następny nagłówek sekcji kończy wykaz nazw słowników.
			wSekcji = false
		case wSekcji && tresc != "":
			nazwy = append(nazwy, filepath.Base(tresc))
		}
	}
	return nazwy
}

// ── vale ────────────────────────────────────────────────────────────────────

// zgloszenieVale opisuje odczytywaną część wyjścia wywołania z opcją
// `--output=JSON`: sprawdzenie, komunikat, dopasowanie i wiersz.
type zgloszenieVale struct {
	Check   string `json:"Check"`
	Message string `json:"Message"`
	Match   string `json:"Match"`
	Line    int    `json:"Line"`
}

// konfiguracjaVale jest treścią pliku `.vale.ini` zakładanego na jedno
// wywołanie, z kontrolą pisowni programu wyłączoną rozstrzygnięciem.
const konfiguracjaVale = "MinAlertLevel = suggestion\n\n[*]\nBasedOnStyles = Vale\nVale.Spelling = NO\n"

// korektaVale mierzy prozę panelu programem Vale i przekłada jego wyjście
// JSON na ustalenia silnika, posortowane dla powtarzalności identyfikatorów.
func (a *adapterTlumaczenia) korektaVale(ctx context.Context,
	katalog, tresc string) ([]ustalenieSilnika, error) {

	konfiguracja := filepath.Join(katalog, nazwaKonfiguracjiVale)
	if err := os.WriteFile(konfiguracja, []byte(konfiguracjaVale), 0o600); err != nil {
		return nil, bladKorektyZewnetrznej("nie można założyć konfiguracji analizatora prozy: " +
			err.Error())
	}
	// `--no-exit` zdejmuje z programu kod wyjścia niezerowy przy zastrzeżeniach.
	wyjscie, err := a.wolajKorekte(ctx, katalog, narzedzieVale, []string{
		"--no-exit", "--output=JSON", nazwaTresciKorekty,
	}, granicaKorektyPisowni)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(wyjscie))) == 0 {
		return nil, nil
	}

	var wedlugPliku map[string][]zgloszenieVale
	if err := json.Unmarshal(wyjscie, &wedlugPliku); err != nil {
		return nil, bladKorektyZewnetrznej("vale oddał wynik, którego nie da się rozłożyć " +
			"jako JSON: " + err.Error())
	}
	ustalenia := []ustalenieSilnika{}
	for _, zgloszenia := range wedlugPliku {
		for _, zgloszenie := range zgloszenia {
			if strings.EqualFold(zgloszenie.Check, "vale.spelling") {
				// Reguła wyłączona konfiguracją nie wejdzie do ustaleń mimo to.
				continue
			}
			ustalenia = append(ustalenia, ustalenieSilnika{
				rodzaj:   shared.ProofreadCheckKindStyle,
				waga:     shared.ProofreadSeverityHint,
				segment:  zgloszenie.Match,
				szczegol: zgloszenie.Message + " (reguła " + zgloszenie.Check + ")",
			})
		}
	}
	// Kolejność wyjścia programu jest kolejnością mapy plików, stąd sortowanie.
	sort.SliceStable(ustalenia, func(i, j int) bool {
		return ustalenia[i].szczegol < ustalenia[j].szczegol
	})
	return ustalenia, nil
}

// ── Wspólne ─────────────────────────────────────────────────────────────────

// wolajKorekte przeprowadza jedno uruchomienie programu korekty w katalogu
// roboczym czynności, podanym wprost zamiast zostawionym bramie wywołania.
func (a *adapterTlumaczenia) wolajKorekte(ctx context.Context, katalog string,
	narzedzie zewnetrzne.Narzedzie, argumenty []string, granica time.Duration) ([]byte, error) {

	okno, zasady, obszar := a.zasiegProgramowTlumaczenia()
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedzie, argumenty, katalog, granica)
	if err != nil {
		return nil, bladNarzedziaKorekty(err)
	}
	return wynik.Wyjscie, nil
}

// fragmentTresci wycina fragment treści wskazany przez program zewnętrzny
// i potwierdza, że wskazanie rzeczywiście mieści się w granicach treści.
func fragmentTresci(znaki []rune, przesuniecie, dlugosc int) (string, bool) {
	if przesuniecie < 0 || dlugosc <= 0 || przesuniecie+dlugosc > len(znaki) {
		return "", false
	}
	return string(znaki[przesuniecie : przesuniecie+dlugosc]), true
}

// bladKorektyZewnetrznej znakuje zaplecze korekty jako niedostępne: żądanie było
// poprawne, rdzeń nie jest zepsuty, brakuje czegoś w instalacji i zdanie mówi
// czego. Ten sam kod i to samo uzasadnienie, co przy braku syntezatora mowy.
func bladKorektyZewnetrznej(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Translate: "+powod))
}

// bladNarzedziaKorekty przekłada odmowy arsenału na kody kontraktu — tym samym
// rozstrzygnięciem, co `bladNarzedziaDokumentu` i silnik syntezy mowy.
func bladNarzedziaKorekty(err error) error {
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return bladKorektyZewnetrznej(brak.Error())
	}
	if errors.Is(err, session.ErrIzolacja) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
			"moduł Translate: punkt izolacji zatrzymał uruchomienie programu korekty: "+err.Error()))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Translate: "+err.Error()))
}
