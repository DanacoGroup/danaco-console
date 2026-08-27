// Odpowiedzialność pliku: silniki zewnętrzne korekty językowej modułu Translate
// — trzy programy, po które sięga `translate.proofread.run`, wraz z przekładem
// ich wyjścia na ustalenia kontraktu. Samą komendę prowadzi
// `adapter_modul_tlumaczenie_korekta.go`; ten plik odpowiada wyłącznie na
// pytanie „co jeszcze widać w tej treści i czym to zmierzono".
//
// ── Dlaczego programy, a nie model ──────────────────────────────────────────
// Reguły wbudowane komendy są mechaniczne i przez to sprawdzalne: to samo
// wywołanie na tej samej treści daje te same ustalenia o tych samych
// identyfikatorach, więc `proofread.apply` ma co zastosować. Trzy programy
// poniżej mają dokładnie tę własność — są słownikami i regułami, nie zgadywaniem
// — i dlatego wchodzą tu, gdzie wywołanie modelu wejść nie może.
//
// ── Podział pracy między trzema ─────────────────────────────────────────────
//   - LANGUAGETOOL prowadzi całość: gramatykę, ortografię, interpunkcję,
//     typografię i styl, każdą regułą nazwaną i z propozycją poprawki. Jest
//     drogą pierwszą, bo jako jedyny widzi zdanie, a nie samo słowo.
//   - HUNSPELL prowadzi samą ortografię i wchodzi WYŁĄCZNIE wtedy, gdy
//     LanguageToola nie ma. Uruchomione razem podwoiłyby każdą literówkę:
//     ortografię LanguageToola liczy ta sama rodzina słowników morfologicznych.
//     To jest ten sam układ pierwszeństwa, co `piper` przed `espeak-ng` przy
//     syntezie mowy i Ruff przed interpreterem Pythona w module Terminal.
//   - VALE prowadzi styl prozy: powtórzenia i terminy. Wchodzi obok
//     LanguageToola, nie zamiast — mierzy tekst jako całość, a nie zdanie,
//     i nie ma z tamtym wspólnych reguł.
//
// ── Wybór hunspella, nie enchanta ───────────────────────────────────────────
// `enchant-2` jest pośrednikiem, nie słownikiem: na maszynie budowy wypisuje
// trzech dostawców (hunspell, aspell, hspell) i dla polskiej treści oddaje
// wynik znak w znak taki sam jak hunspell wołany wprost — bo woła właśnie jego.
// Pośrednik dokłada dwie rzeczy i obie są kosztem: warstwę, która nie mówi,
// który dostawca odpowiedział, oraz zestaw modułów wtyczkowych do spakowania
// obok programu. Hunspell jest jednym plikiem wykonywalnym z plikami słownika
// obok, czyli układem, który pakowanie produktu (`pomocniki/<program>`) niesie
// wprost, a jego słowniki są tymi samymi plikami, które ma już LibreOffice
// stojący w wykazie zależności.
//
// ── Dlaczego wiersz poleceń, a nie serwer LanguageToola ─────────────────────
// LanguageTool umie chodzić jako usługa i wtedy nie płaci startem maszyny
// wirtualnej przy każdym wywołaniu. Rdzeń tej drogi dziś nie ma: jedyna brama
// do procesu (`zewnetrzne.Wolaj`) prowadzi uruchomienie OD startu DO końca,
// z obowiązkową granicą czasu i ubiciem całego drzewa procesów — czyli robi
// dokładnie to, czego usłudze długożyjącej robić nie wolno. Do tego usługa
// potrzebowałaby portu, a maszyna wdrożenia ma porty zajęte przez rzeczy
// niezwiązane z produktem. Wiersz poleceń działa dziś i idzie tą samą bramą,
// co każdy inny program rdzenia.
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
	// z wczytaniem reguł i słownika morfologicznego języka. Na maszynie budowy
	// całość mieści się w kilku sekundach; granica zostawia zapas na maszynę
	// wolniejszą i pierwszy przebieg bez rozgrzanej pamięci podręcznej plików.
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

	// nazwaKonfiguracjiVale jest nazwą pliku, bez którego vale nie rusza:
	// program nie ma wbudowanego zestawu reguł i przy braku konfiguracji kończy
	// się błędem `E100`. Rdzeń zakłada ją sam w katalogu roboczym czynności —
	// treść panelu nie leży w żadnym repozytorium, więc konfiguracji nie ma skąd
	// wziąć.
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

// narzedzieVale opisuje analizator prozy.
var narzedzieVale = zewnetrzne.Narzedzie{
	Nazwa: "vale", Program: "vale", Pakiet: "vale (jeden plik wykonywalny z wydania projektu)",
}

// ustalenieSilnika niesie jedno ustalenie zewnętrznego silnika w postaci, którą
// komenda odkłada w bazie. Propozycja jest CAŁĄ treścią panelu po poprawce —
// tak samo jak przy regułach wbudowanych, bo `proofread.apply` wstawia ją
// w miejsce treści, a nie w miejsce słowa.
type ustalenieSilnika struct {
	rodzaj     shared.ProofreadCheckKind
	waga       shared.ProofreadSeverity
	segment    string
	szczegol   string
	propozycja string
}

// ustaleniaSilnikow zbiera ustalenia wszystkich silników, które na tej maszynie
// stoją.
//
// Silnik nieobecny NIE JEST tu odmową: jego brak jest stanem maszyny, o którym
// sonda startowa powiedziała Operatorowi przy uruchomieniu rdzenia, a komenda
// ma oddać to, co da się zmierzyć. Silnik obecny, który zawiódł, jest czymś
// innym — to zdarzenie w trakcie czynności i wraca odmową, bo cisza w jego
// miejscu byłaby brakiem pomiaru podanym jako brak zastrzeżeń.
func (a *adapterTlumaczenia) ustaleniaSilnikow(ctx context.Context,
	jezyk, tresc string) ([]ustalenieSilnika, error) {

	if a.uruchamiacz == nil {
		// Rdzeń złożony bez warstwy kanału nie ma czym wystartować programu.
		// Reguły wbudowane komendy pracują dalej — to jest mniej niż pełnia,
		// ale nie brama postawiona przed resztą.
		return nil, nil
	}
	// Słownik jest słownikiem JEDNEGO języka, więc oba silniki pisowni pracują
	// tylko wtedy, gdy język panelu da się nazwać oznaczeniem. Vale pracuje
	// niezależnie: mierzy powtórzenia i terminy, a te widać bez słownika.
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

// korektaLanguageToolem przeprowadza jeden przebieg i przekłada jego wynik.
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
			// Położenie, którego nie da się potwierdzić w treści, jest brakiem
			// pomiaru, nie ustaleniem: wstawienie poprawki pod zgadniętym
			// numerem znaku popsułoby panel.
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

// rodzajKorektyLanguageToola przekłada kategorię reguły na rodzaj kontroli
// kontraktu.
//
// Przekład idzie po KATEGORII, a nie po treści komunikatu: kategoria jest
// polem wyjścia programu, komunikat jest zdaniem po polsku, które zmieni się
// przy pierwszym poprawionym tłumaczeniu reguł. Kategoria nierozpoznana wraca
// jako gramatyka — to najszerszy z ośmiu rodzajów kontraktu i ten sam, którym
// reguły wbudowane komendy znakują powtórzone słowo. Kategoria nierozpoznana
// nie jest przy tym pominięciem: ustalenie i tak dochodzi do Operatora, bo
// zgubione byłoby stratą, a źle nazwane jest wciąż prawdziwe.
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

// wagaKorektyLanguageToola nadaje ustaleniu wagę. LanguageTool wagi nie podaje,
// a kontrakt jej wymaga, więc rdzeń rozstrzyga sam i robi to jedną zasadą:
// błędem jest to, co jest faktem o słowie (słowa nie ma w słowniku),
// ostrzeżeniem to, co jest orzeczeniem reguły o zdaniu. Ta sama zasada dzieli
// reguły wbudowane komendy — odstęp przed przecinkiem jest tam błędem, a
// cudzysłów prosty wskazówką.
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
// słownik. Pustka znaczy „tego języka rdzeń nie umie nazwać" — nie „domyślny".
//
// Kontrakt trzyma język panelu jako napis dowolny („Jezyk docelowy"), więc
// stoją w nim i oznaczenia (`pl`, `en-GB`), i nazwy własne („angielski").
// Rozstrzygnięcie jest trzystopniowe i każdy stopień ma powód:
//
//  1. Napis o KSZTAŁCIE oznaczenia języka idzie do programu bez zmiany. Rdzeń
//     nie ma i nie ma prawa mieć tabeli języków LanguageToola — program zna ich
//     sześćdziesiąt kilka, wykaz zmienia się z wydaniami, a oznaczenie
//     nieznane program odrzuca sam, wymieniając wszystkie, które zna. Sprawdzamy
//     więc KSZTAŁT, a nie przynależność do wykazu.
//  2. Nazwa własna języka sprowadza się do oznaczenia PODSTAWOWEGO, bez odmiany
//     krajowej. „Angielski" nie mówi, czy chodzi o pisownię brytyjską, czy
//     amerykańską, a te różnią się tysiącami słów — dopisanie odmiany byłoby
//     rozstrzygnięciem za Operatora. Zestaw nazw jest ten sam, który rdzeń
//     rozpoznaje przy wskazaniu języka rozpoznania pisma
//     (`jezykRozpoznaniaDokumentu`).
//  3. Napis, który nie jest ani oznaczeniem, ani znaną nazwą, nie idzie
//     nigdzie. Wysłany na chybił trafił dałby korektę polskiego panelu regułami
//     innego języka — wynik wyglądałby jak korekta i byłby zmyśleniem.
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
// dowolna liczba członów po myślniku (`pl`, `pl-PL`, `de-DE-x-simple-language`).
// To jest kształt oznaczeń, których używa LanguageTool i którymi opisuje się
// język w BCP-47.
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

// korektaHunspellem sprawdza ortografię treści panelu.
//
// Wyjście idzie w postaci ispellowej: wiersz `& słowo licz przesunięcie: a, b,
// c` niesie słowo spoza słownika wraz z propozycjami, wiersz `#` to słowo
// spoza słownika bez propozycji, a `*`, `+` i `-` znaczą słowo znane.
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
			// Podmiana wchodzi WYŁĄCZNIE przy jednym wystąpieniu słowa.
			// Hunspell podaje przesunięcie względem wiersza, a nie treści, więc
			// przy drugim wystąpieniu rdzeń nie wie, o które chodzi — i wtedy
			// ustalenie idzie bez propozycji, bo kontrakt mówi wprost, że jej
			// brak znaczy „korekta jej nie podała".
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

// niepoprawneSlowoHunspella czyta jeden wiersz wyjścia ispellowego.
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

// slownikHunspella dobiera nazwę słownika do języka panelu.
//
// Rozstrzyga PYTANIEM O STAN, nie tabelą w kodzie: `hunspell -D` wypisuje
// słowniki, które na tej maszynie stoją, i to jest jedyna prawda o tym, co
// program otworzy. Tabela języków w rdzeniu byłaby drugą prawdą — a przy
// oznaczeniu bez kraju (`pl` zamiast `pl_PL`) i tak nie miałaby czego wybrać.
// Ta sama droga rozstrzyga brak skanera (`odmowaSkanuSane`): pytanie o wykaz,
// a nie czytanie diagnostyki obcego programu.
func (a *adapterTlumaczenia) slownikHunspella(ctx context.Context,
	katalog, jezyk string) (string, error) {

	oznaczenie := strings.ReplaceAll(strings.TrimSpace(jezyk), "-", "_")
	if oznaczenie == "" {
		return "", bladWskazaniaTlumaczenia("panel nie ma języka — nie ma czym dobrać " +
			"słownika ortograficznego")
	}

	// Pusty plik jako materiał: `-D` wypisuje wykaz i kończy pracę, ale bez
	// wskazania pliku czekałby na wejście, którego port uruchamiacza nie
	// wystawia.
	pusty := filepath.Join(katalog, "wykaz-slownikow.txt")
	if err := os.WriteFile(pusty, nil, 0o600); err != nil {
		return "", bladKorektyZewnetrznej("nie można założyć pliku pomocniczego: " + err.Error())
	}
	okno, zasady, obszar := a.zasiegProgramowTlumaczenia()
	wynik, blad := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedzieHunspella, []string{"-D", pusty}, katalog, granicaKorektyPisowni)

	// Kod wyjścia tego wywołania NIE ROZSTRZYGA i nie jest tu czytany: `-D` jest
	// trybem diagnostycznym, który kończy się kodem 1 zawsze — także wtedy, gdy
	// wypisał komplet słowników. Odpowiedzią na pytanie „co ten program ma" jest
	// sam wykaz, a nie kod wyjścia; wykaz pusty jest jedynym stanem, który
	// znaczy „nie ma czym mierzyć". Diagnostyka pierwotna wchodzi wtedy do
	// odmowy, żeby powód nie zginął.
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
	// Oznaczenie bez kraju dopasowujemy do słownika po samym języku — ale
	// wyłącznie wtedy, gdy kandydat jest DOKŁADNIE JEDEN. Dwa kandydaty to
	// wybór, którego rdzeń nie ma prawa dokonać za Operatora: `en_US` i `en_GB`
	// różnią się pisownią tysięcy słów.
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
			// Każdy następny nagłówek sekcji kończy wykaz — po nim idzie
			// „LOADED DICTIONARY" wraz ze ścieżkami plików `.aff` i `.dic`,
			// które nazwami słowników nie są.
			wSekcji = false
		case wSekcji && tresc != "":
			nazwy = append(nazwy, filepath.Base(tresc))
		}
	}
	return nazwy
}

// ── vale ────────────────────────────────────────────────────────────────────

// zgloszenieVale opisuje odczytywaną część wyjścia `--output=JSON`.
type zgloszenieVale struct {
	Check   string `json:"Check"`
	Message string `json:"Message"`
	Match   string `json:"Match"`
	Line    int    `json:"Line"`
}

// konfiguracjaVale jest treścią pliku `.vale.ini` zakładanego na jedno
// wywołanie.
//
// Zestaw reguł jest wbudowany w program (`Vale`) — powtórzenia i terminy — więc
// czynność nie wymaga żadnego pakietu stylów z sieci. Kontrola pisowni tego
// zestawu jest WYŁĄCZONA i to jest rozstrzygnięcie, nie przeoczenie:
// `Vale.Spelling` ma słownik wyłącznie angielski i na polskim panelu zgłasza
// każde polskie słowo jako błąd. Ortografię prowadzi w tej komendzie
// LanguageTool albo hunspell — oba ze słownikiem języka panelu.
const konfiguracjaVale = "MinAlertLevel = suggestion\n\n[*]\nBasedOnStyles = Vale\nVale.Spelling = NO\n"

// korektaVale mierzy prozę panelu.
func (a *adapterTlumaczenia) korektaVale(ctx context.Context,
	katalog, tresc string) ([]ustalenieSilnika, error) {

	konfiguracja := filepath.Join(katalog, nazwaKonfiguracjiVale)
	if err := os.WriteFile(konfiguracja, []byte(konfiguracjaVale), 0o600); err != nil {
		return nil, bladKorektyZewnetrznej("nie można założyć konfiguracji analizatora prozy: " +
			err.Error())
	}
	// `--no-exit` zdejmuje z programu kod wyjścia niezerowy przy znalezionych
	// zastrzeżeniach. Bez niego każdy udany pomiar wracałby przez arsenał jako
	// niepowodzenie programu — czyli znaleziona usterka tekstu wyglądałaby jak
	// usterka rdzenia.
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
				// Reguła wyłączona konfiguracją; gdyby wydanie programu
				// przestało tę nastawę honorować, ustalenia i tak tu nie wejdą.
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
	// Kolejność wyjścia programu jest kolejnością mapy plików, więc porządek
	// ustaleń nie byłby powtarzalny między wywołaniami. Ustalenia dostają
	// identyfikatory w kolejności, w jakiej tu stoją, a `proofread.apply`
	// adresuje je tymi identyfikatorami.
	sort.SliceStable(ustalenia, func(i, j int) bool {
		return ustalenia[i].szczegol < ustalenia[j].szczegol
	})
	return ustalenia, nil
}

// ── Wspólne ─────────────────────────────────────────────────────────────────

// wolajKorekte przeprowadza jedno uruchomienie programu korekty w katalogu
// roboczym czynności.
//
// Zasięg platformy składa `zasiegProgramowTlumaczenia` — ta sama trójka okno–zasady–obszar,
// którą jedzie synteza mowy tego modułu. Katalog roboczy jest tu podany WPROST,
// nie zostawiony bramie: vale szuka swojej konfiguracji w katalogu, z którego
// ruszył, więc uruchomienie gdzie indziej kończyłoby się błędem braku
// konfiguracji, choć plik leży przygotowany.
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

// fragmentTresci wycina fragment treści wskazany przez program i POTWIERDZA, że
// wskazanie mieści się w treści.
//
// Potwierdzenie jest tu warunkiem, nie ostrożnością: numer znaku podany przez
// obcy program liczy się w jego własnej jednostce, a wstawienie poprawki pod
// numerem, którego rdzeń nie sprawdził, popsułoby treść panelu w miejscu
// wybranym przypadkiem.
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
