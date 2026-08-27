// Plik obsługuje trzy komendy zakładki zależności i bezpieczeństwa w Dev
// Tools: `developer.dependency.list`, `developer.scan.run`
// i `developer.scan.result.list`.
package core

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/repozytorium"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

const (
	// przedrostekSkanu znakuje identyfikator przebiegu skanowania w wykazie
	// przebiegów i w odwołaniach do niego ze znalezisk.
	przedrostekSkanu = "scan-"
	// przedrostekZnaleziska znakuje identyfikator pojedynczego spostrzeżenia
	// dowolnego rodzaju skanu w wykazie znalezisk.
	przedrostekZnaleziska = "find-"
	// najwiecejZnaleziskWykazu jest domyślną głębokością wykazu znalezisk
	// zwracaną, gdy żądanie nie ogranicza liczby wyników jawnie.
	najwiecejZnaleziskWykazu = 500
)

// WykazZaleznosci obsługuje komendę `developer.dependency.list`: czyta drzewo
// zależności z manifestu repozytorium wskazanego albo wykrytego automatycznie.
func (a *adapterDevelopera) WykazZaleznosci(_ context.Context,
	z shared.DeveloperDependencyListRequest) (shared.DeveloperDependencyListResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperDependencyListResponse{}, err
	}
	korzenie := a.korzenieOkna(okno)
	if len(korzenie) == 0 {
		return shared.DeveloperDependencyListResponse{}, bladZadaniaDevelopera(
			"okno " + z.WindowId + " nie ma katalogu roboczego, więc nie ma gdzie szukać manifestu")
	}

	sciezka := ""
	if z.Manifest != nil && strings.TrimSpace(*z.Manifest) != "" {
		wskazana, err := sciezkaWObszarze(korzenie, *z.Manifest, plikIstnieje)
		if err != nil {
			return shared.DeveloperDependencyListResponse{}, err
		}
		sciezka = wskazana
	} else {
		wykryta, jest := manifestRepozytorium(korzenie[0])
		if !jest {
			return shared.DeveloperDependencyListResponse{}, bladZasobuDevelopera(
				"w katalogu roboczym okna nie ma manifestu zależności " +
					"(go.mod, package.json, requirements.txt ani Cargo.toml)")
		}
		sciezka = wykryta
	}

	zaleznosci, err := zaleznosciZManifestu(sciezka)
	if err != nil {
		return shared.DeveloperDependencyListResponse{}, err
	}
	if z.OutdatedOnly != nil && *z.OutdatedOnly {
		// Zawężenie „tylko przestarzałe" oddaje pozycje bez przypiętej wersji.
		zawezone := make([]shared.DependencyNode, 0, len(zaleznosci))
		for _, pozycja := range zaleznosci {
			if wersjaNieprzypieta(pozycja.Version) {
				zawezone = append(zawezone, pozycja)
			}
		}
		zaleznosci = zawezone
	}
	if z.Depth != nil && *z.Depth == 1 {
		bezposrednie := make([]shared.DependencyNode, 0, len(zaleznosci))
		for _, pozycja := range zaleznosci {
			if pozycja.Direct {
				bezposrednie = append(bezposrednie, pozycja)
			}
		}
		zaleznosci = bezposrednie
	}
	return shared.DeveloperDependencyListResponse{
		Dependencies: zaleznosci,
		Manifest:     sciezka,
	}, nil
}

// manifestRepozytorium odnajduje manifest w katalogu roboczym, sprawdzając
// nazwy w ustalonej kolejności: `go.mod`, `package.json`, `requirements.txt`,
// `Cargo.toml`.
func manifestRepozytorium(korzen string) (string, bool) {
	for _, nazwa := range []string{"go.mod", "package.json", "requirements.txt", "Cargo.toml"} {
		sciezka := filepath.Join(korzen, nazwa)
		if plikIstnieje(sciezka) {
			return sciezka, true
		}
	}
	return "", false
}

// zaleznosciZManifestu czyta manifest wskazanego rodzaju i oddaje jego
// zależności w kształcie wspólnym dla wszystkich obsługiwanych manifestów.
func zaleznosciZManifestu(sciezka string) ([]shared.DependencyNode, error) {
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		return nil, bladZasobuDevelopera(
			"nie można odczytać manifestu " + sciezka + ": " + err.Error())
	}
	tresc := string(bajty)

	switch strings.ToLower(filepath.Base(sciezka)) {
	case "go.mod":
		return zaleznosciGo(tresc), nil
	case "package.json":
		return zaleznosciNode(tresc)
	case "requirements.txt":
		return zaleznosciPythona(tresc), nil
	case "cargo.toml":
		return zaleznosciRusta(tresc), nil
	default:
		return nil, bladZadaniaDevelopera(
			"plik " + filepath.Base(sciezka) + " nie jest znanym manifestem zależności")
	}
}

// zaleznosciGo czyta `go.mod`.
//
// Rozróżnienie zależności bezpośredniej od pośredniej bierze się ze znacznika
// `// indirect`, który zapisuje tam narzędzie języka — to jest ta sama wiedza,
// którą ma `go list`, zapisana w pliku.
func zaleznosciGo(tresc string) []shared.DependencyNode {
	zaleznosci := make([]shared.DependencyNode, 0, 32)
	wBloku := false
	for _, wiersz := range strings.Split(tresc, "\n") {
		pole := strings.TrimSpace(wiersz)
		switch {
		case pole == "require (":
			wBloku = true
			continue
		case wBloku && pole == ")":
			wBloku = false
			continue
		case strings.HasPrefix(pole, "require "):
			pole = strings.TrimPrefix(pole, "require ")
		case !wBloku:
			continue
		}
		if pole == "" || strings.HasPrefix(pole, "//") {
			continue
		}
		posrednia := strings.Contains(pole, "// indirect")
		if komentarz := strings.Index(pole, "//"); komentarz >= 0 {
			pole = strings.TrimSpace(pole[:komentarz])
		}
		czesci := strings.Fields(pole)
		if len(czesci) < 2 {
			continue
		}
		zaleznosci = append(zaleznosci, shared.DependencyNode{
			Name:    czesci[0],
			Version: czesci[1],
			Direct:  !posrednia,
		})
	}
	return uporzadkujZaleznosci(zaleznosci)
}

// zaleznosciNode czyta `package.json` i oddaje jego zależności produkcyjne,
// deweloperskie i towarzyszące wraz z licencją pakietu.
func zaleznosciNode(tresc string) ([]shared.DependencyNode, error) {
	var manifest struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
		PeerZaleznosci  map[string]string `json:"peerDependencies"`
		License         string            `json:"license"`
	}
	if err := json.Unmarshal([]byte(tresc), &manifest); err != nil {
		return nil, bladZadaniaDevelopera(
			"package.json nie jest poprawnym dokumentem JSON: " + err.Error())
	}
	zaleznosci := make([]shared.DependencyNode, 0, 64)
	dopisz := func(wpisy map[string]string, rodzic string) {
		for nazwa, wersja := range wpisy {
			pozycja := shared.DependencyNode{Name: nazwa, Version: wersja, Direct: rodzic == ""}
			if rodzic != "" {
				pozycja.ParentPath = wskaznikTekstu(rodzic)
			}
			zaleznosci = append(zaleznosci, pozycja)
		}
	}
	dopisz(manifest.Dependencies, "")
	// Zależności deweloperskie i towarzyszące stoją pod własnym rodzicem.
	dopisz(manifest.DevDependencies, "devDependencies")
	dopisz(manifest.PeerZaleznosci, "peerDependencies")
	return uporzadkujZaleznosci(zaleznosci), nil
}

// zaleznosciPythona czyta `requirements.txt` i oddaje wskazane w nim
// zależności wraz z operatorem wersji, gdy manifest go niesie.
func zaleznosciPythona(tresc string) []shared.DependencyNode {
	zaleznosci := make([]shared.DependencyNode, 0, 32)
	for _, wiersz := range strings.Split(tresc, "\n") {
		pole := strings.TrimSpace(wiersz)
		if pole == "" || strings.HasPrefix(pole, "#") || strings.HasPrefix(pole, "-") {
			continue
		}
		if komentarz := strings.Index(pole, " #"); komentarz >= 0 {
			pole = strings.TrimSpace(pole[:komentarz])
		}
		nazwa, wersja := pole, ""
		for _, znak := range []string{"==", ">=", "<=", "~=", "!=", ">", "<"} {
			if miejsce := strings.Index(pole, znak); miejsce > 0 {
				nazwa = strings.TrimSpace(pole[:miejsce])
				wersja = strings.TrimSpace(pole[miejsce:])
				break
			}
		}
		if wersja == "" {
			wersja = "dowolna"
		}
		zaleznosci = append(zaleznosci, shared.DependencyNode{
			Name: nazwa, Version: wersja, Direct: true,
		})
	}
	return uporzadkujZaleznosci(zaleznosci)
}

// zaleznosciRusta czyta sekcje zależności pliku `Cargo.toml`, rozróżniając
// zależności zwykłe od pozostałych po nazwie sekcji.
func zaleznosciRusta(tresc string) []shared.DependencyNode {
	zaleznosci := make([]shared.DependencyNode, 0, 32)
	sekcja := ""
	for _, wiersz := range strings.Split(tresc, "\n") {
		pole := strings.TrimSpace(wiersz)
		if strings.HasPrefix(pole, "[") && strings.HasSuffix(pole, "]") {
			sekcja = strings.Trim(pole, "[]")
			continue
		}
		if !strings.HasSuffix(sekcja, "dependencies") || pole == "" ||
			strings.HasPrefix(pole, "#") {
			continue
		}
		rowne := strings.Index(pole, "=")
		if rowne <= 0 {
			continue
		}
		nazwa := strings.TrimSpace(pole[:rowne])
		wartosc := strings.TrimSpace(pole[rowne+1:])
		wersja := strings.Trim(wartosc, `"`)
		if strings.HasPrefix(wartosc, "{") {
			// Z postaci rozbudowanej `{ version = "1.0", ... }` bierzemy wersję.
			wersja = wersjaZTabeliCargo(wartosc)
		}
		if wersja == "" {
			wersja = "dowolna"
		}
		zaleznosci = append(zaleznosci, shared.DependencyNode{
			Name:    nazwa,
			Version: wersja,
			Direct:  sekcja == "dependencies",
		})
	}
	return uporzadkujZaleznosci(zaleznosci)
}

// wersjaZTabeliCargo wyjmuje wartość pola `version` z zapisu rozbudowanego
// zależności, ograniczoną parą cudzysłowów.
func wersjaZTabeliCargo(wartosc string) string {
	miejsce := strings.Index(wartosc, "version")
	if miejsce < 0 {
		return ""
	}
	reszta := wartosc[miejsce:]
	pierwszy := strings.Index(reszta, `"`)
	if pierwszy < 0 {
		return ""
	}
	drugi := strings.Index(reszta[pierwszy+1:], `"`)
	if drugi < 0 {
		return ""
	}
	return reszta[pierwszy+1 : pierwszy+1+drugi]
}

// uporzadkujZaleznosci ustala kolejność wykazu: najpierw bezpośrednie, potem
// alfabetycznie. Kolejność mapy Go jest losowa, a wykaz zmieniający kolejność
// przy każdym odczycie nie da się porównać z poprzednim.
func uporzadkujZaleznosci(zaleznosci []shared.DependencyNode) []shared.DependencyNode {
	sort.SliceStable(zaleznosci, func(i, j int) bool {
		if zaleznosci[i].Direct != zaleznosci[j].Direct {
			return zaleznosci[i].Direct
		}
		return zaleznosci[i].Name < zaleznosci[j].Name
	})
	return zaleznosci
}

// wersjaNieprzypieta mówi, czy zapis wersji dopuszcza podmianę treści
// zależności bez zmiany manifestu.
func wersjaNieprzypieta(wersja string) bool {
	tresc := strings.TrimSpace(wersja)
	if tresc == "" || tresc == "dowolna" || tresc == "*" || tresc == "latest" {
		return true
	}
	return strings.HasPrefix(tresc, "^") || strings.HasPrefix(tresc, "~") ||
		strings.HasPrefix(tresc, ">") || strings.HasPrefix(tresc, "<")
}

// ── Skanowanie ──────────────────────────────────────────────────────────────

// regulySekretow rozpoznają kształty kluczy i haseł w treści plików: token
// dostawcy, klucz prywatny, klucz dostępowy chmury, hasło w adresie połączenia.
var regulySekretow = []struct {
	nazwa   string
	wzorzec *regexp.Regexp
	waga    shared.ProblemSeverity
}{
	{"klucz prywatny", regexp.MustCompile(`-----BEGIN (RSA |EC |OPENSSH |PGP )?PRIVATE KEY-----`),
		shared.ProblemSeverityError},
	{"token dostępowy GitHub", regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{36,}`),
		shared.ProblemSeverityError},
	{"klucz dostępowy AWS", regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		shared.ProblemSeverityError},
	{"klucz dostawcy modelu", regexp.MustCompile(`sk-[A-Za-z0-9_\-]{32,}`),
		shared.ProblemSeverityError},
	{"hasło w adresie połączenia",
		regexp.MustCompile(`(?i)(postgres|postgresql|mysql|mongodb|redis|amqp)://[^:\s/]+:[^@\s]+@`),
		shared.ProblemSeverityError},
	{"hasło wpisane w kod",
		regexp.MustCompile(`(?i)(password|passwd|secret|api[_-]?key)\s*[:=]\s*["'][^"'\s]{8,}["']`),
		shared.ProblemSeverityWarning},
}

// regulyKodu rozpoznają wzorce, które w Go i TypeScripcie bywają źródłem
// podatności: sklejanie zapytań, wykonanie tekstu jako kodu i tym podobne.
var regulyKodu = []struct {
	nazwa   string
	wzorzec *regexp.Regexp
	waga    shared.ProblemSeverity
	opis    string
}{
	{"sklejanie zapytania SQL",
		regexp.MustCompile(`(?i)"[^"\n]*\b(SELECT|INSERT|UPDATE|DELETE)\b[^"\n]*"\s*\+`),
		shared.ProblemSeverityError,
		"treść polecenia SQL powstaje przez sklejenie tekstu z wartością; " +
			"wartość z zewnątrz wchodzi wtedy do treści polecenia — naprawa: parametr zapytania"},
	{"wykonanie tekstu jako kodu",
		regexp.MustCompile(`(?i)\b(eval|new Function)\s*\(`),
		shared.ProblemSeverityWarning,
		"tekst wykonywany jako kod wykonuje także tekst, który przyszedł z zewnątrz"},
	{"wstawianie treści bez oczyszczenia",
		regexp.MustCompile(`\.innerHTML\s*=`),
		shared.ProblemSeverityWarning,
		"treść wstawiona jako HTML wykona zawarte w niej znaczniki; " +
			"naprawa: textContent albo oczyszczenie treści"},
	{"pominięte sprawdzenie certyfikatu",
		regexp.MustCompile(`InsecureSkipVerify\s*:\s*true`),
		shared.ProblemSeverityError,
		"połączenie szyfrowane bez sprawdzenia certyfikatu nie chroni przed podszyciem"},
	{"pominięty błąd",
		regexp.MustCompile(`(?m)^\s*_\s*(,\s*_\s*)*=\s*\w+\.(Close|Write|Exec)\(`),
		shared.ProblemSeverityInfo,
		"wynik czynności jest odrzucany; niepowodzenie przejdzie niezauważone"},
}

// UruchomSkan obsługuje komendę `developer.scan.run`: przeprowadza skan
// wskazanych rodzajów w całości przed odpowiedzią i zapisuje jego wynik.
func (a *adapterDevelopera) UruchomSkan(ctx context.Context,
	z shared.DeveloperScanRunRequest) (shared.DeveloperScanRunResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperScanRunResponse{}, err
	}
	if a.repozytorium == nil {
		return shared.DeveloperScanRunResponse{}, bladZasobuDevelopera(
			"rdzeń nie ma miejsca na wyniki skanowania")
	}
	korzenie := a.korzenieOkna(okno)
	if len(korzenie) == 0 {
		return shared.DeveloperScanRunResponse{}, bladZadaniaDevelopera(
			"okno " + z.WindowId + " nie ma katalogu roboczego, więc nie ma czego skanować")
	}
	rodzaje := rodzajeSkanu(z.Kinds)
	if len(rodzaje) == 0 {
		return shared.DeveloperScanRunResponse{}, bladZadaniaDevelopera(
			"skan wymaga wskazania choć jednego rodzaju: dependencies, secrets, code albo licenses")
	}

	kod := nowyIdentyfikator(przedrostekSkanu)
	zapisRodzajow := make([]string, 0, len(rodzaje))
	for _, rodzaj := range rodzaje {
		zapisRodzajow = append(zapisRodzajow, string(rodzaj))
	}
	poczatek := time.Now().UTC()
	if err := a.repozytorium.ZapiszSkan(ctx, danePrzebieguSkanu(kod, okno.Id,
		strings.Join(zapisRodzajow, ","), shared.BuildStatusRunning, nil, nil)); err != nil {
		return shared.DeveloperScanRunResponse{}, bladWykonaniaDevelopera(
			"nie można założyć przebiegu skanowania: " + err.Error())
	}

	znaleziska, err := a.przeprowadzSkan(ctx, okno, korzenie[0], kod, rodzaje, z.Paths)
	stan := shared.BuildStatus(shared.BuildStatusSucceeded)
	if err != nil {
		stan = shared.BuildStatusFailed
	}
	if len(znaleziska) > 0 {
		if zapis := a.repozytorium.ZapiszZnaleziska(ctx, wierszeZnalezisk(znaleziska)); zapis != nil {
			return shared.DeveloperScanRunResponse{}, bladWykonaniaDevelopera(
				"nie można zapisać znalezisk skanowania: " + zapis.Error())
		}
	}
	liczba := int64(len(znaleziska))
	koniec := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	if zapis := a.repozytorium.ZapiszSkan(ctx, danePrzebieguSkanu(kod, okno.Id,
		strings.Join(zapisRodzajow, ","), stan, &liczba, &koniec)); zapis != nil {
		return shared.DeveloperScanRunResponse{}, bladWykonaniaDevelopera(
			"nie można domknąć przebiegu skanowania: " + zapis.Error())
	}
	if err != nil {
		return shared.DeveloperScanRunResponse{}, err
	}

	ile := len(znaleziska)
	zakonczono := koniec
	return shared.DeveloperScanRunResponse{Scan: shared.ScanRun{
		Id:           kod,
		WindowId:     okno.Id,
		Kinds:        rodzaje,
		Status:       stan,
		FindingCount: &ile,
		StartedAt:    poczatek.UnixMilli(),
		FinishedAt:   wskaznikDuzejChwili(chwilaBazy(zakonczono)),
	}}, nil
}

// wskaznikDuzejChwili oddaje wskaźnik chwili — kontrakt niesie ją jako wartość
// opcjonalną, a przebieg domknięty zawsze ją ma.
func wskaznikDuzejChwili(chwila int64) *int64 { return &chwila }

// rodzajeSkanu odsiewa z żądanych rodzajów skanu powtórzenia i wartości spoza
// słownika kontraktu, zachowując kolejność pierwszego wystąpienia.
func rodzajeSkanu(zadane []shared.ScanKind) []shared.ScanKind {
	widziane := map[shared.ScanKind]bool{}
	rodzaje := make([]shared.ScanKind, 0, len(zadane))
	for _, rodzaj := range zadane {
		switch rodzaj {
		case shared.ScanKindDependencies, shared.ScanKindSecrets, shared.ScanKindCode,
			shared.ScanKindLicenses:
			if !widziane[rodzaj] {
				widziane[rodzaj] = true
				rodzaje = append(rodzaje, rodzaj)
			}
		}
	}
	return rodzaje
}

// przeprowadzSkan wykonuje wskazane rodzaje skanu w katalogu roboczym i oddaje
// zebrane znaleziska niezależnie od rodzaju, który je wystawił.
func (a *adapterDevelopera) przeprowadzSkan(ctx context.Context, okno session.Okno,
	korzen, skanKod string, rodzaje []shared.ScanKind,
	sciezki []string) ([]daneZnaleziska, error) {

	znaleziska := make([]daneZnaleziska, 0, 32)
	trescioweRodzaje := false
	skanKodu := false
	for _, rodzaj := range rodzaje {
		switch rodzaj {
		case shared.ScanKindSecrets:
			trescioweRodzaje = true
		case shared.ScanKindCode:
			trescioweRodzaje = true
			skanKodu = true
		case shared.ScanKindDependencies, shared.ScanKindLicenses:
			znaleziska = append(znaleziska,
				znaleziskaZManifestu(korzen, skanKod, rodzaj)...)
		}
	}
	if skanKodu {
		znaleziska = append(znaleziska,
			a.znaleziskaProgramow(ctx, okno, korzen, skanKod, sciezki)...)
	}
	if !trescioweRodzaje {
		return znaleziska, nil
	}

	// Przejście po plikach idzie tą samą drogą, co wyszukiwanie w repozytorium.
	pliki, err := repozytorium.PlikiDoPrzejrzenia(korzen, sciezki)
	if err != nil {
		return znaleziska, bladWykonaniaDevelopera(
			"nie można przejrzeć plików repozytorium: " + err.Error())
	}
	for _, sciezka := range pliki {
		trescPliku, err := os.ReadFile(sciezka)
		if err != nil {
			continue
		}
		for _, rodzaj := range rodzaje {
			switch rodzaj {
			case shared.ScanKindSecrets:
				znaleziska = append(znaleziska,
					znaleziskaWTresci(skanKod, sciezka, string(trescPliku), rodzaj)...)
			case shared.ScanKindCode:
				znaleziska = append(znaleziska,
					znaleziskaWKodzie(skanKod, sciezka, string(trescPliku))...)
			}
		}
	}
	return znaleziska, nil
}

// znaleziskaWTresci przepuszcza treść pliku wiersz po wierszu przez reguły
// sekretów i zwraca znaleziska z numerem wiersza dopasowania.
func znaleziskaWTresci(skanKod, sciezka, tresc string,
	rodzaj shared.ScanKind) []daneZnaleziska {

	znaleziska := make([]daneZnaleziska, 0, 2)
	czytnik := bufio.NewScanner(strings.NewReader(tresc))
	czytnik.Buffer(make([]byte, 0, 64*1024), 1<<20)
	numer := 0
	for czytnik.Scan() {
		numer++
		wiersz := czytnik.Text()
		for _, regula := range regulySekretow {
			if !regula.wzorzec.MatchString(wiersz) {
				continue
			}
			linia := int64(numer)
			znaleziska = append(znaleziska, daneZnaleziska{
				kod:     nowyIdentyfikator(przedrostekZnaleziska),
				skanKod: skanKod,
				rodzaj:  rodzaj,
				waga:    regula.waga,
				tytul:   regula.nazwa + " w pliku " + filepath.Base(sciezka),
				// Treść dopasowania nie wchodzi do znaleziska.
				opis:    "reguła rozpoznała kształt sekretu; wartości nie zapisano",
				sciezka: sciezka,
				wiersz:  &linia,
				regula:  regula.nazwa,
			})
		}
	}
	return znaleziska
}

// znaleziskaWKodzie przepuszcza treść pliku przez reguły wzorców podatności,
// ograniczone do rozszerzeń języków, które te reguły rozpoznają.
func znaleziskaWKodzie(skanKod, sciezka, tresc string) []daneZnaleziska {
	rozszerzenie := strings.ToLower(filepath.Ext(sciezka))
	switch rozszerzenie {
	case ".go", ".ts", ".tsx", ".js", ".jsx", ".py", ".rs", ".java", ".rb", ".php":
	default:
		return nil
	}

	znaleziska := make([]daneZnaleziska, 0, 2)
	for numer, wiersz := range strings.Split(tresc, "\n") {
		for _, regula := range regulyKodu {
			if !regula.wzorzec.MatchString(wiersz) {
				continue
			}
			linia := int64(numer + 1)
			znaleziska = append(znaleziska, daneZnaleziska{
				kod:     nowyIdentyfikator(przedrostekZnaleziska),
				skanKod: skanKod,
				rodzaj:  shared.ScanKindCode,
				waga:    regula.waga,
				tytul:   regula.nazwa,
				opis:    regula.opis,
				sciezka: sciezka,
				wiersz:  &linia,
				regula:  regula.nazwa,
			})
		}
	}
	return znaleziska
}

// znaleziskaZManifestu orzeka o zależnościach bez przypiętej wersji albo
// o brakujących licencjach wyłącznie z treści manifestu repozytorium.
func znaleziskaZManifestu(korzen, skanKod string, rodzaj shared.ScanKind) []daneZnaleziska {
	sciezka, jest := manifestRepozytorium(korzen)
	if !jest {
		return nil
	}
	zaleznosci, err := zaleznosciZManifestu(sciezka)
	if err != nil {
		return nil
	}

	znaleziska := make([]daneZnaleziska, 0, 8)
	for _, pozycja := range zaleznosci {
		switch rodzaj {
		case shared.ScanKindDependencies:
			if !wersjaNieprzypieta(pozycja.Version) {
				continue
			}
			znaleziska = append(znaleziska, daneZnaleziska{
				kod:     nowyIdentyfikator(przedrostekZnaleziska),
				skanKod: skanKod,
				rodzaj:  rodzaj,
				waga:    shared.ProblemSeverityWarning,
				tytul:   "zależność " + pozycja.Name + " bez przypiętej wersji",
				opis: "wersja " + pozycja.Version + " dopuszcza podmianę treści zależności " +
					"bez zmiany w manifeście — dwa budowania tego samego commitu mogą dać " +
					"różny wynik",
				sciezka: sciezka,
				regula:  "wersja nieprzypięta",
				pakiet:  pozycja.Name,
			})
		case shared.ScanKindLicenses:
			if pozycja.License != nil && *pozycja.License != "" {
				continue
			}
			znaleziska = append(znaleziska, daneZnaleziska{
				kod:     nowyIdentyfikator(przedrostekZnaleziska),
				skanKod: skanKod,
				rodzaj:  rodzaj,
				waga:    shared.ProblemSeverityInfo,
				tytul:   "licencja zależności " + pozycja.Name + " nie jest znana z manifestu",
				opis: "manifest nie niesie licencji tej zależności; ustalenie jej wymaga " +
					"odczytu pakietu, którego serwer bez pobrania zależności nie ma",
				sciezka: sciezka,
				regula:  "licencja nieznana",
				pakiet:  pozycja.Name,
			})
		}
	}
	return znaleziska
}

// WykazZnalezisk obsługuje komendę `developer.scan.result.list`: zwraca
// znaleziska wskazanego przebiegu skanowania albo wszystkie znaleziska okna.
func (a *adapterDevelopera) WykazZnalezisk(ctx context.Context,
	z shared.DeveloperScanResultListRequest) (shared.DeveloperScanResultListResponse, error) {

	if a.repozytorium == nil {
		return shared.DeveloperScanResultListResponse{}, bladZasobuDevelopera(
			"rdzeń nie ma wyników skanowania")
	}
	filtr := filtrZnaleziskZadania(z)
	if filtr.SkanKod == "" && filtr.OknoKod == "" {
		return shared.DeveloperScanResultListResponse{}, bladZadaniaDevelopera(
			"wykaz znalezisk wymaga wskazania przebiegu skanowania albo okna")
	}
	if filtr.OknoKod != "" {
		if _, err := a.oknoDevelopera(filtr.OknoKod); err != nil {
			return shared.DeveloperScanResultListResponse{}, err
		}
	}

	wiersze, err := a.repozytorium.Znaleziska(ctx, filtr)
	if err != nil {
		return shared.DeveloperScanResultListResponse{}, bladWykonaniaDevelopera(
			"nie można odczytać znalezisk: " + err.Error())
	}
	znaleziska := make([]shared.ScanFinding, 0, len(wiersze))
	for _, wiersz := range wiersze {
		znaleziska = append(znaleziska, shared.ScanFinding{
			Id:           wiersz.Kod,
			Kind:         wiersz.Rodzaj,
			Severity:     wiersz.Waga,
			Title:        wiersz.Tytul,
			Description:  wiersz.Opis,
			Path:         wiersz.Sciezka,
			Line:         liczbaZDuzej(wiersz.Wiersz),
			RuleId:       wiersz.Regula,
			CveId:        wiersz.Cve,
			PackageName:  wiersz.Pakiet,
			FixedVersion: wiersz.WersjaNaprawy,
		})
	}
	razem := len(znaleziska)
	return shared.DeveloperScanResultListResponse{Findings: znaleziska, Total: &razem}, nil
}

// daneZnaleziska jest spostrzeżeniem skanu w postaci wygodnej do składania:
// reguły skanu wypełniają pola wartościami wprost, a nie wskaźnikami.
type daneZnaleziska struct {
	kod           string
	skanKod       string
	rodzaj        shared.ScanKind
	waga          shared.ProblemSeverity
	tytul         string
	opis          string
	sciezka       string
	wiersz        *int64
	regula        string
	cve           string
	pakiet        string
	wersjaNaprawy string
}

// wierszZnaleziska przekłada spostrzeżenie na wiersz bazy. Pole puste zostaje
// puste, a nie zapisane jako pusty tekst — brak wartości i wartość pusta to
// w odczycie dwie różne rzeczy.
func (z daneZnaleziska) wierszZnaleziska() dane.ZnaleziskoSkanu {
	wiersz := dane.ZnaleziskoSkanu{
		Kod:     z.kod,
		SkanKod: z.skanKod,
		Rodzaj:  z.rodzaj,
		Waga:    z.waga,
		Tytul:   z.tytul,
		Wiersz:  z.wiersz,
	}
	if z.opis != "" {
		wiersz.Opis = wskaznikTekstu(z.opis)
	}
	if z.sciezka != "" {
		wiersz.Sciezka = wskaznikTekstu(z.sciezka)
	}
	if z.regula != "" {
		wiersz.Regula = wskaznikTekstu(z.regula)
	}
	if z.cve != "" {
		wiersz.Cve = wskaznikTekstu(z.cve)
	}
	if z.pakiet != "" {
		wiersz.Pakiet = wskaznikTekstu(z.pakiet)
	}
	if z.wersjaNaprawy != "" {
		wiersz.WersjaNaprawy = wskaznikTekstu(z.wersjaNaprawy)
	}
	return wiersz
}

// wierszeZnalezisk przekłada cały zbiór spostrzeżeń skanu na wiersze bazy,
// wywołując dla każdej pozycji jej metodę przekładu.
func wierszeZnalezisk(znaleziska []daneZnaleziska) []dane.ZnaleziskoSkanu {
	wiersze := make([]dane.ZnaleziskoSkanu, 0, len(znaleziska))
	for _, znalezisko := range znaleziska {
		wiersze = append(wiersze, znalezisko.wierszZnaleziska())
	}
	return wiersze
}

// danePrzebieguSkanu składa wiersz nagłówka przebiegu skanowania z jego
// stanem, rodzajami i, gdy przebieg jest domknięty, liczbą znalezisk.
func danePrzebieguSkanu(kod, oknoKod, rodzaje string, stan shared.BuildStatus,
	znalezisk *int64, zakonczono *string) dane.PrzebiegSkanu {

	return dane.PrzebiegSkanu{
		Kod:        kod,
		OknoKod:    oknoKod,
		Rodzaje:    rodzaje,
		Stan:       stan,
		Znalezisk:  znalezisk,
		Zakonczono: zakonczono,
	}
}

// filtrZnaleziskZadania składa zawężenie wykazu znalezisk z pól żądania:
// przebiegu, okna, rodzaju, wagi i granicy liczby wyników.
func filtrZnaleziskZadania(z shared.DeveloperScanResultListRequest) dane.FiltrZnalezisk {
	filtr := dane.FiltrZnalezisk{Limit: najwiecejZnaleziskWykazu}
	if z.ScanId != nil {
		filtr.SkanKod = strings.TrimSpace(*z.ScanId)
	}
	if z.WindowId != nil {
		filtr.OknoKod = strings.TrimSpace(*z.WindowId)
	}
	if z.Kind != nil {
		filtr.Rodzaj = string(*z.Kind)
	}
	if z.Severity != nil {
		filtr.Waga = string(*z.Severity)
	}
	if z.Limit != nil && *z.Limit > 0 && *z.Limit < filtr.Limit {
		filtr.Limit = *z.Limit
	}
	return filtr
}

// Programy zewnętrzne niżej poszerzają skan kodu wyrażeniami regularnymi.

// czasProgramuSkanu jest granicą czasową jednego przejścia programu
// analizującego po całym repozytorium podczas skanu kodu.
const czasProgramuSkanu = 180 * time.Second

// Próg powtórzenia rdzeń nie narzuca — obowiązuje próg własny programu.

// konfiguracjeSemgrepa wylicza nazwy, pod którymi Semgrep szuka zestawu reguł
// w korzeniu repozytorium.
var konfiguracjeSemgrepa = []string{
	".semgrep.yml", ".semgrep.yaml", "semgrep.yml", "semgrep.yaml", ".semgrep",
}

// znaleziskaProgramow zbiera spostrzeżenia programów poszerzających skan
// kodu: Semgrepa oraz wykrywaczy powtórzeń dla TypeScriptu, JavaScriptu i Go.
func (a *adapterDevelopera) znaleziskaProgramow(ctx context.Context, okno session.Okno,
	korzen, skanKod string, sciezki []string) []daneZnaleziska {

	znaleziska := make([]daneZnaleziska, 0, 8)
	znaleziska = append(znaleziska,
		a.znaleziskaSemgrepa(ctx, okno, korzen, skanKod, sciezki)...)
	znaleziska = append(znaleziska,
		a.znaleziskaPowtorzenTypeScriptu(ctx, okno, korzen, skanKod)...)
	znaleziska = append(znaleziska,
		a.znaleziskaPowtorzenGo(ctx, okno, korzen, skanKod)...)
	return znaleziska
}

// wyjscieSemgrepa jest kształtem odpowiedzi `semgrep --json` z wykazem
// wyników wraz z ich położeniem i wagą.
type wyjscieSemgrepa struct {
	Results []struct {
		CheckId string `json:"check_id"`
		Path    string `json:"path"`
		Start   struct {
			Line int `json:"line"`
		} `json:"start"`
		Extra struct {
			Message  string `json:"message"`
			Severity string `json:"severity"`
		} `json:"extra"`
	} `json:"results"`
}

// znaleziskaSemgrepa przeprowadza skan semantyczny zestawem reguł repozytorium
// — nie zestawem z rejestru — i oddaje pustą listę bez własnego zestawu reguł.
func (a *adapterDevelopera) znaleziskaSemgrepa(ctx context.Context, okno session.Okno,
	korzen, skanKod string, sciezki []string) []daneZnaleziska {

	if !repozytoriumNiesieKonfiguracje(korzen, konfiguracjeSemgrepa) {
		return nil
	}
	argumenty := []string{"scan", "--json", "--quiet", "--metrics=off",
		"--disable-version-check", "--config", "."}
	for _, sciezka := range sciezki {
		if tresc := strings.TrimSpace(sciezka); tresc != "" {
			argumenty = append(argumenty, tresc)
		}
	}

	wynik, err := a.wolajNarzedzieWarsztatu(ctx, okno, narzedzieSemgrep, argumenty,
		korzen, czasProgramuSkanu)
	if err != nil && len(wynik.Wyjscie) == 0 {
		return nil
	}
	var odpowiedz wyjscieSemgrepa
	if json.Unmarshal([]byte(strings.TrimSpace(string(wynik.Wyjscie))), &odpowiedz) != nil {
		return nil
	}

	znaleziska := make([]daneZnaleziska, 0, len(odpowiedz.Results))
	for _, uwaga := range odpowiedz.Results {
		wiersz := int64(uwaga.Start.Line)
		znaleziska = append(znaleziska, daneZnaleziska{
			kod:     nowyIdentyfikator(przedrostekZnaleziska),
			skanKod: skanKod,
			rodzaj:  shared.ScanKindCode,
			waga:    wagaSemgrepa(uwaga.Extra.Severity),
			tytul:   uwaga.CheckId,
			opis:    uwaga.Extra.Message,
			sciezka: sciezkaWzgledemKorzenia(uwaga.Path, korzen),
			wiersz:  &wiersz,
			regula:  uwaga.CheckId,
		})
	}
	return znaleziska
}

// wagaSemgrepa przekłada wagę Semgrepa na wagę kontraktu. Waga nierozpoznana
// zostaje ostrzeżeniem — rdzeń nie podnosi wagi, której program tak nie nazwał.
func wagaSemgrepa(waga string) shared.ProblemSeverity {
	switch strings.ToUpper(strings.TrimSpace(waga)) {
	case "ERROR":
		return shared.ProblemSeverityError
	case "INFO":
		return shared.ProblemSeverityInfo
	default:
		return shared.ProblemSeverityWarning
	}
}

// wyjsciePowtorzen jest kształtem raportu `jscpd --reporters json` z wykazem
// par plików niosących ten sam powtórzony fragment.
type wyjsciePowtorzen struct {
	Duplicates []struct {
		Format    string `json:"format"`
		Lines     int    `json:"lines"`
		FirstFile struct {
			Name  string `json:"name"`
			Start int    `json:"start"`
			End   int    `json:"end"`
		} `json:"firstFile"`
		SecondFile struct {
			Name  string `json:"name"`
			Start int    `json:"start"`
			End   int    `json:"end"`
		} `json:"secondFile"`
	} `json:"duplicates"`
}

// znaleziskaPowtorzenTypeScriptu szuka powtórzonych fragmentów w plikach
// TypeScriptu i JavaScriptu programem jscpd, z raportem w katalogu tymczasowym.
func (a *adapterDevelopera) znaleziskaPowtorzenTypeScriptu(ctx context.Context,
	okno session.Okno, korzen, skanKod string) []daneZnaleziska {

	katalog, err := os.MkdirTemp("", "danaco-powtorzenia-*")
	if err != nil {
		return nil
	}
	defer func() { _ = os.RemoveAll(katalog) }()

	// Powtórzenia w plikach Go liczy program `dupl` osobno, niżej.
	_, wolanie := a.wolajNarzedzieWarsztatu(ctx, okno, narzedzieJscpd, []string{
		"--reporters", "json", "--output", katalog, "--silent",
		"--format", "typescript,javascript,jsx,tsx", ".",
	}, korzen, czasProgramuSkanu)

	bajty, odczyt := os.ReadFile(filepath.Join(katalog, "jscpd-report.json"))
	if odczyt != nil {
		// Raportu nie ma — droga podstawowa skanu i tak już zmierzyła swoje.
		_ = wolanie
		return nil
	}
	var raport wyjsciePowtorzen
	if json.Unmarshal(bajty, &raport) != nil {
		return nil
	}

	znaleziska := make([]daneZnaleziska, 0, len(raport.Duplicates))
	for _, powtorzenie := range raport.Duplicates {
		wiersz := int64(powtorzenie.FirstFile.Start)
		znaleziska = append(znaleziska, daneZnaleziska{
			kod:     nowyIdentyfikator(przedrostekZnaleziska),
			skanKod: skanKod,
			rodzaj:  shared.ScanKindCode,
			waga:    shared.ProblemSeverityInfo,
			tytul:   "powtórzony fragment " + strconv.Itoa(powtorzenie.Lines) + " wierszy",
			opis: "ten sam fragment stoi w " + powtorzenie.FirstFile.Name + ":" +
				strconv.Itoa(powtorzenie.FirstFile.Start) + " oraz w " +
				powtorzenie.SecondFile.Name + ":" +
				strconv.Itoa(powtorzenie.SecondFile.Start),
			sciezka: sciezkaWzgledemKorzenia(powtorzenie.FirstFile.Name, korzen),
			wiersz:  &wiersz,
			regula:  narzedzieJscpd.Program,
		})
	}
	return znaleziska
}

// wzorzecPowtorzeniaGo rozbiera wiersz odpowiedzi `dupl -plumbing` postaci
// `plik:od-do: duplicate of plik:od-do` na obie strony powtórzenia.
var wzorzecPowtorzeniaGo = regexp.MustCompile(
	`^(.+):(\d+)-(\d+): duplicate of (.+):(\d+)-(\d+)$`)

// znaleziskaPowtorzenGo szuka powtórzonych fragmentów w plikach Go programem
// dupl, zapisując z każdej pary powtórzenia wyłącznie jedną stronę.
func (a *adapterDevelopera) znaleziskaPowtorzenGo(ctx context.Context, okno session.Okno,
	korzen, skanKod string) []daneZnaleziska {

	wynik, err := a.wolajNarzedzieWarsztatu(ctx, okno, narzedzieDupl,
		[]string{"-plumbing", "."}, korzen, czasProgramuSkanu)
	if err != nil && len(wynik.Wyjscie) == 0 {
		return nil
	}

	znaleziska := make([]daneZnaleziska, 0, 8)
	widziane := map[string]bool{}
	for _, wiersz := range strings.Split(string(wynik.Wyjscie), "\n") {
		czesci := wzorzecPowtorzeniaGo.FindStringSubmatch(strings.TrimSpace(wiersz))
		if czesci == nil {
			continue
		}
		para := czesci[1] + ":" + czesci[2] + "|" + czesci[4] + ":" + czesci[5]
		odwrotna := czesci[4] + ":" + czesci[5] + "|" + czesci[1] + ":" + czesci[2]
		if widziane[para] || widziane[odwrotna] {
			continue
		}
		widziane[para] = true

		poczatek, blad := strconv.Atoi(czesci[2])
		if blad != nil {
			continue
		}
		koniec, blad := strconv.Atoi(czesci[3])
		if blad != nil {
			continue
		}
		numer := int64(poczatek)
		znaleziska = append(znaleziska, daneZnaleziska{
			kod:     nowyIdentyfikator(przedrostekZnaleziska),
			skanKod: skanKod,
			rodzaj:  shared.ScanKindCode,
			waga:    shared.ProblemSeverityInfo,
			tytul: "powtórzony fragment " + strconv.Itoa(koniec-poczatek+1) +
				" wierszy",
			opis: "ten sam fragment stoi w " + czesci[1] + ":" + czesci[2] +
				" oraz w " + czesci[4] + ":" + czesci[5],
			sciezka: sciezkaWzgledemKorzenia(czesci[1], korzen),
			wiersz:  &numer,
			regula:  narzedzieDupl.Program,
		})
	}
	return znaleziska
}
