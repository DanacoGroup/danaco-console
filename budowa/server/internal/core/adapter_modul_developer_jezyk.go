// Plik obsługuje cztery komendy warstwy językowej Code Editora:
// `developer.symbol.navigate`, `developer.format.run`, `developer.lint.get`
// i `developer.refactor.apply`, wołając programy serwera języka na maszynie.
package core

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

const (
	// czasSerweraJezyka jest granicą jednego pytania do serwera języka. Pierwsze
	// pytanie w nieznanym repozytorium ładuje jego graf odwołań, więc granica
	// jest liczona w dziesiątkach sekund, nie w sekundach.
	czasSerweraJezyka = 90 * time.Second
	// czasFormatowania jest granicą czasową formatowania jednego pliku
	// pojedynczym wywołaniem programu formatującego.
	czasFormatowania = 20 * time.Second
	// czasAnalizyStatycznej jest granicą czasową jednego przebiegu programu
	// analizy statycznej po całym repozytorium.
	czasAnalizyStatycznej = 180 * time.Second
	// najwiecejZgloszenAnalizy chroni Dev Tools przed wykazem, którego nikt nie
	// przeczyta, gdy linter zgłosi tysiące uwag w nieuporządkowanym repozytorium.
	najwiecejZgloszenAnalizy = 500
)

// NawigujDoSymbolu obsługuje komendę `developer.symbol.navigate`: przejście do
// definicji, implementacji, wystąpień albo wykaz symboli dokumentu.
func (a *adapterDevelopera) NawigujDoSymbolu(ctx context.Context,
	z shared.DeveloperSymbolNavigateRequest) (shared.DeveloperSymbolNavigateResponse, error) {

	okno, sciezka, err := a.plikOkna(z.WindowId, z.Path)
	if err != nil {
		return shared.DeveloperSymbolNavigateResponse{}, err
	}
	if z.Line < 1 {
		return shared.DeveloperSymbolNavigateResponse{}, bladZadaniaDevelopera(
			"nawigacja po symbolu wymaga wiersza liczonego od jedynki")
	}
	kolumna := z.Column
	if kolumna < 1 {
		kolumna = 1
	}

	miejsce := sciezka + ":" + strconv.Itoa(z.Line) + ":" + strconv.Itoa(kolumna)
	var argumenty []string
	switch z.Kind {
	case shared.SymbolNavigationKindDefinition:
		argumenty = []string{"definition", "-json", miejsce}
	case shared.SymbolNavigationKindImplementation:
		argumenty = []string{"implementation", miejsce}
	case shared.SymbolNavigationKindReferences:
		argumenty = []string{"references", miejsce}
	case shared.SymbolNavigationKindDocumentSymbol:
		argumenty = []string{"symbols", sciezka}
	default:
		return shared.DeveloperSymbolNavigateResponse{}, bladZadaniaDevelopera(
			"nieznany rodzaj nawigacji po symbolu: " + string(z.Kind))
	}

	if _, jednorazowy := serwerJezykaPliku(sciezka); !jednorazowy {
		// Pusty wykaz przy `serverAvailable: false` znaczy brak narzędzia, nie
		// brak wyniku.
		return shared.DeveloperSymbolNavigateResponse{
			Symbols: []shared.DeveloperSymbol{}, ServerAvailable: false}, nil
	}

	wynik, err := a.wolajNarzedzieWarsztatu(ctx, okno, narzedzieGopls, argumenty,
		filepath.Dir(sciezka), czasSerweraJezyka)
	if err != nil {
		if brakNarzedziaWarsztatu(err) {
			return shared.DeveloperSymbolNavigateResponse{
				Symbols: []shared.DeveloperSymbol{}, ServerAvailable: false}, nil
		}
		// Niezerowy kod serwera języka bez symbolu pod kursorem nie jest
		// awarią.
		return shared.DeveloperSymbolNavigateResponse{
			Symbols: []shared.DeveloperSymbol{}, ServerAvailable: true}, nil
	}

	symbole := symboleZOdpowiedziSerwera(z.Kind, string(wynik.Wyjscie), sciezka)
	return shared.DeveloperSymbolNavigateResponse{Symbols: symbole, ServerAvailable: true}, nil
}

// symbolGopls jest kształtem odpowiedzi `gopls definition -json`, z którego
// rdzeń czyta położenie i opis znalezionego symbolu.
type symbolGopls struct {
	Span struct {
		URI   string `json:"uri"`
		Start struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"start"`
	} `json:"span"`
	Description string `json:"description"`
}

// symboleZOdpowiedziSerwera przekłada wyjście serwera języka na kontrakt:
// dokument JSON dla definicji albo wiersze tekstu dla wystąpień i symboli
// dokumentu.
func symboleZOdpowiedziSerwera(rodzaj shared.SymbolNavigationKind, wyjscie,
	sciezkaPliku string) []shared.DeveloperSymbol {

	symbole := make([]shared.DeveloperSymbol, 0, 8)
	if rodzaj == shared.SymbolNavigationKindDefinition {
		var opis symbolGopls
		if err := json.Unmarshal([]byte(strings.TrimSpace(wyjscie)), &opis); err == nil &&
			opis.Span.URI != "" {
			symbole = append(symbole, shared.DeveloperSymbol{
				Name:   pierwszeSlowoOpisu(opis.Description),
				Path:   sciezkaZOdwolania(opis.Span.URI),
				Line:   opis.Span.Start.Line,
				Column: wskaznikLiczby(opis.Span.Start.Column),
			})
			return symbole
		}
	}
	for _, wiersz := range strings.Split(wyjscie, "\n") {
		tresc := strings.TrimSpace(wiersz)
		if tresc == "" {
			continue
		}
		if symbol, jest := symbolZWiersza(tresc, sciezkaPliku); jest {
			symbole = append(symbole, symbol)
		}
	}
	return symbole
}

// symbolZWiersza rozbiera wiersz postaci `plik:wiersz:kolumna[-…] opis` albo
// `Nazwa Rodzaj wiersz:kolumna-…`, którymi serwer języka odpowiada na pytania
// o wystąpienia i o symbole dokumentu.
func symbolZWiersza(wiersz, sciezkaPliku string) (shared.DeveloperSymbol, bool) {
	czesci := strings.SplitN(wiersz, ":", 4)
	if len(czesci) >= 3 {
		if numer, err := strconv.Atoi(czesci[1]); err == nil {
			symbol := shared.DeveloperSymbol{
				Name: strings.TrimSpace(strings.Join(czesci[3:], ":")),
				Path: czesci[0],
				Line: numer,
			}
			if kolumna, err := strconv.Atoi(oberwijZakres(czesci[2])); err == nil {
				symbol.Column = wskaznikLiczby(kolumna)
			}
			if symbol.Name == "" {
				symbol.Name = filepath.Base(czesci[0])
			}
			return symbol, true
		}
	}

	// Wykaz symboli dokumentu: `Nazwa Rodzaj wiersz:kolumna-wiersz:kolumna`.
	pola := strings.Fields(wiersz)
	if len(pola) < 3 {
		return shared.DeveloperSymbol{}, false
	}
	polozenie := strings.SplitN(pola[len(pola)-1], ":", 2)
	numer, err := strconv.Atoi(polozenie[0])
	if err != nil {
		return shared.DeveloperSymbol{}, false
	}
	symbol := shared.DeveloperSymbol{
		Name: pola[0],
		Kind: wskaznikTekstu(pola[1]),
		Path: sciezkaPliku,
		Line: numer,
	}
	if len(polozenie) == 2 {
		if kolumna, err := strconv.Atoi(oberwijZakres(polozenie[1])); err == nil {
			symbol.Column = wskaznikLiczby(kolumna)
		}
	}
	return symbol, true
}

// oberwijZakres zostawia z zapisu `12-15` samo `12`, czyli początek zakresu
// zwróconego przez serwer języka.
func oberwijZakres(wartosc string) string {
	if myslnik := strings.IndexByte(wartosc, '-'); myslnik > 0 {
		return wartosc[:myslnik]
	}
	return wartosc
}

// pierwszeSlowoOpisu bierze z opisu zwróconego przez serwer języka pierwsze
// albo drugie słowo jako nazwę odnalezionego symbolu.
func pierwszeSlowoOpisu(opis string) string {
	pola := strings.Fields(opis)
	if len(pola) == 0 {
		return "symbol"
	}
	if len(pola) > 1 && (pola[0] == "func" || pola[0] == "type" || pola[0] == "var" ||
		pola[0] == "const") {
		return pola[1]
	}
	return pola[0]
}

// sciezkaZOdwolania zdejmuje przedrostek schematu `file://` z odwołania do
// pliku zwróconego przez serwer języka w opisie symbolu.
func sciezkaZOdwolania(odwolanie string) string {
	return strings.TrimPrefix(odwolanie, "file://")
}

// Formatuj obsługuje komendę `developer.format.run`: formatuje treść
// z żądania albo treść pliku programem właściwym dla rozszerzenia.
func (a *adapterDevelopera) Formatuj(ctx context.Context,
	z shared.DeveloperFormatRunRequest) (shared.DeveloperFormatRunResponse, error) {

	okno, sciezka, err := a.plikOkna(z.WindowId, z.Path)
	if err != nil {
		return shared.DeveloperFormatRunResponse{}, err
	}
	zapisz := z.WriteToDisk != nil && *z.WriteToDisk
	if zapisz {
		if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "zapis sformatowanego pliku"); err != nil {
			return shared.DeveloperFormatRunResponse{}, err
		}
	}

	przed := ""
	if z.Content != nil {
		przed = *z.Content
	} else {
		bajty, err := os.ReadFile(sciezka)
		if err != nil {
			return shared.DeveloperFormatRunResponse{}, bladZasobuDevelopera(
				"nie można odczytać pliku " + z.Path + " do formatowania: " + err.Error())
		}
		przed = string(bajty)
	}

	narzedzie, argumenty, jest := formaterPliku(sciezka)
	if !jest {
		return shared.DeveloperFormatRunResponse{}, bladZadaniaDevelopera(
			"dla pliku " + filepath.Base(sciezka) + " serwer nie ma formatera; " +
				"formatowane są pliki Go (gofmt, goimports) oraz obsługiwane przez Prettier")
	}

	// Treść idzie do formatera plikiem tymczasowym — port uruchamiacza nie
	// podaje wejścia procesowi.
	roboczy, sprzataj, err := plikRoboczyFormatowania(sciezka, przed)
	if err != nil {
		return shared.DeveloperFormatRunResponse{}, err
	}
	defer sprzataj()

	wynik, err := a.wolajNarzedzieWarsztatu(ctx, okno, narzedzie,
		append(argumenty, roboczy), filepath.Dir(sciezka), czasFormatowania)
	if err != nil {
		if brakNarzedziaWarsztatu(err) {
			return shared.DeveloperFormatRunResponse{}, bladZasobuDevelopera(
				"serwer nie ma programu " + narzedzie.Nazwa + " — formatowanie pliku " +
					filepath.Base(sciezka) + " nie ma czym się wykonać; naprawa: " +
					narzedzie.Pakiet)
		}
		return shared.DeveloperFormatRunResponse{}, bladWykonaniaDevelopera(
			"formater " + narzedzie.Nazwa + " odmówił: " + skrocDiagnostyke(wynik.Diagnostyka, err))
	}

	po := string(wynik.Wyjscie)
	if strings.TrimSpace(po) == "" {
		po = przed
	}
	po = zawezDoZakresu(przed, po, z.StartLine, z.EndLine)

	if zapisz && po != przed {
		if err := os.WriteFile(sciezka, []byte(po), 0o644); err != nil {
			return shared.DeveloperFormatRunResponse{}, bladWykonaniaDevelopera(
				"nie można zapisać sformatowanego pliku " + z.Path + ": " + err.Error())
		}
	}

	opis, err := os.Stat(sciezka)
	if err != nil {
		return shared.DeveloperFormatRunResponse{}, bladZasobuDevelopera(
			"plik " + z.Path + " zniknął w trakcie formatowania")
	}
	plik := opisPliku(sciezka, opis)
	plik.Content = wskaznikTekstu(po)
	return shared.DeveloperFormatRunResponse{
		File:      plik,
		Changed:   po != przed,
		Formatter: wskaznikTekstu(narzedzie.Program),
	}, nil
}

// formaterPliku dobiera program formatujący po rozszerzeniu pliku spośród
// programów zainstalowanych na serwerze.
func formaterPliku(sciezka string) (zewnetrzne.Narzedzie, []string, bool) {
	switch strings.ToLower(filepath.Ext(sciezka)) {
	case ".go":
		if zewnetrzne.Stoi(narzedzieGoimports) {
			return narzedzieGoimports, nil, true
		}
		return narzedzieGofmt, nil, true
	case ".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs", ".json", ".css", ".scss", ".less",
		".html", ".md", ".yaml", ".yml", ".vue":
		return narzedziePrettier, []string{"--stdin-filepath", sciezka}, true
	default:
		return zewnetrzne.Narzedzie{}, nil, false
	}
}

// plikRoboczyFormatowania odkłada treść bufora w pliku tymczasowym obok
// pliku źródłowego, na czas jednego wywołania formatera.
func plikRoboczyFormatowania(sciezka, tresc string) (string, func(), error) {
	roboczy := filepath.Join(filepath.Dir(sciezka),
		"."+filepath.Base(sciezka)+".danaco-format"+filepath.Ext(sciezka))
	if err := os.WriteFile(roboczy, []byte(tresc), 0o600); err != nil {
		return "", func() {}, bladWykonaniaDevelopera(
			"nie można odłożyć treści do formatowania: " + err.Error())
	}
	return roboczy, func() { _ = os.Remove(roboczy) }, nil
}

// zawezDoZakresu składa wynik formatowania zaznaczenia z treści sprzed
// formatowania i wyniku formatera dla wskazanego zakresu wierszy.
func zawezDoZakresu(przed, po string, odWiersza, doWiersza *int) string {
	if odWiersza == nil || doWiersza == nil {
		return po
	}
	wierszePrzed := strings.Split(przed, "\n")
	wierszePo := strings.Split(po, "\n")
	od, do := *odWiersza, *doWiersza
	if od < 1 || do < od || do > len(wierszePrzed) || len(wierszePo) < do {
		// Formatowanie przesunęło wiersze — zakres przestał się zgadzać, więc
		// wraca pełny wynik formatera.
		return po
	}
	zlozony := make([]string, 0, len(wierszePrzed))
	zlozony = append(zlozony, wierszePrzed[:od-1]...)
	zlozony = append(zlozony, wierszePo[od-1:do]...)
	zlozony = append(zlozony, wierszePrzed[do:]...)
	return strings.Join(zlozony, "\n")
}

// AnalizaStatyczna obsługuje komendę `developer.lint.get`: uruchamia programy
// analizy statycznej właściwe dla plików repozytorium i zwraca ich zgłoszenia.
func (a *adapterDevelopera) AnalizaStatyczna(ctx context.Context,
	z shared.DeveloperLintGetRequest) (shared.DeveloperLintGetResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperLintGetResponse{}, err
	}
	korzenie := a.korzenieOkna(okno)
	if len(korzenie) == 0 {
		return shared.DeveloperLintGetResponse{}, bladZadaniaDevelopera(
			"okno " + z.WindowId + " nie ma katalogu roboczego, więc nie ma czego analizować")
	}
	korzen := korzenie[0]

	zgloszenia := make([]shared.DeveloperDiagnostic, 0, 32)
	odpowiedzialo := false
	for _, analiza := range analizyRepozytorium(korzen, z.Paths) {
		uwagi, zmierzono := a.przeprowadzAnalize(ctx, okno, korzen, analiza)
		if !zmierzono {
			continue
		}
		odpowiedzialo = true
		zgloszenia = append(zgloszenia, uwagi...)
	}

	if !odpowiedzialo {
		return shared.DeveloperLintGetResponse{
			Diagnostics: []shared.DeveloperDiagnostic{}, LinterAvailable: false}, nil
	}
	sort.SliceStable(zgloszenia, func(i, j int) bool {
		if zgloszenia[i].Path != zgloszenia[j].Path {
			return zgloszenia[i].Path < zgloszenia[j].Path
		}
		return zgloszenia[i].Line < zgloszenia[j].Line
	})

	granica := najwiecejZgloszenAnalizy
	if z.Limit != nil && *z.Limit > 0 && *z.Limit < granica {
		granica = *z.Limit
	}
	przyciete := false
	if len(zgloszenia) > granica {
		zgloszenia, przyciete = zgloszenia[:granica], true
	}
	odpowiedz := shared.DeveloperLintGetResponse{Diagnostics: zgloszenia, LinterAvailable: true}
	if przyciete {
		odpowiedz.Truncated = wskaznikPrawdy(true)
	}
	return odpowiedz, nil
}

// wyjscieAnalizy jest kształtem odpowiedzi `golangci-lint run` w postaci
// dokumentu JSON z wykazem zgłoszeń.
type wyjscieAnalizy struct {
	Issues []struct {
		FromLinter string `json:"FromLinter"`
		Text       string `json:"Text"`
		Severity   string `json:"Severity"`
		Pos        struct {
			Filename string `json:"Filename"`
			Line     int    `json:"Line"`
			Column   int    `json:"Column"`
		} `json:"Pos"`
		Replacement *struct{} `json:"Replacement"`
		SourceLines []string  `json:"SourceLines"`
	} `json:"Issues"`
}

// zgloszeniaAnalizy przekłada wyjście programu `golangci-lint` na wykaz
// zgłoszeń kontraktu, uzupełniony ścieżkami bezwzględnymi.
func zgloszeniaAnalizy(wyjscie, korzen string) []shared.DeveloperDiagnostic {
	zgloszenia := make([]shared.DeveloperDiagnostic, 0, 32)
	var odpowiedz wyjscieAnalizy
	if err := json.Unmarshal([]byte(strings.TrimSpace(wyjscie)), &odpowiedz); err != nil {
		return zgloszenia
	}
	for _, uwaga := range odpowiedz.Issues {
		sciezka := uwaga.Pos.Filename
		if !filepath.IsAbs(sciezka) {
			sciezka = filepath.Join(korzen, sciezka)
		}
		zgloszenie := shared.DeveloperDiagnostic{
			Path:     sciezka,
			Line:     uwaga.Pos.Line,
			Severity: wagaAnalizy(uwaga.Severity),
			Message:  uwaga.Text,
			Source:   wskaznikTekstu(uwaga.FromLinter),
		}
		if uwaga.Pos.Column > 0 {
			zgloszenie.Column = wskaznikLiczby(uwaga.Pos.Column)
		}
		if uwaga.Replacement != nil {
			zgloszenie.FixAvailable = wskaznikPrawdy(true)
		}
		zgloszenia = append(zgloszenia, zgloszenie)
	}
	sort.SliceStable(zgloszenia, func(i, j int) bool {
		if zgloszenia[i].Path != zgloszenia[j].Path {
			return zgloszenia[i].Path < zgloszenia[j].Path
		}
		return zgloszenia[i].Line < zgloszenia[j].Line
	})
	return zgloszenia
}

// wagaAnalizy przekłada wagę lintera na wagę kontraktu. Linter, który wagi nie
// podał, zgłasza ostrzeżenie — nie błąd: podniesienie wagi bez podstawy kazałoby
// Operatorowi naprawiać rzeczy, których nikt tak nie oznaczył.
func wagaAnalizy(waga string) shared.ProblemSeverity {
	switch strings.ToLower(strings.TrimSpace(waga)) {
	case "error":
		return shared.ProblemSeverityError
	case "info", "information":
		return shared.ProblemSeverityInfo
	default:
		return shared.ProblemSeverityWarning
	}
}

// serwerJezykaPliku dobiera serwer języka po rozszerzeniu pliku i mówi, czy
// rdzeń ma czym go zapytać pojedynczym wywołaniem. Rozstrzygnięcie trwałe:
// `gopls` ma tryb wiersza poleceń obok trybu LSP, a
// `typescript-language-server` mówi wyłącznie sesją protokołu LSP przez
// stdin/stdout — warstwa językowa pyta serwery pojedynczym wywołaniem i tej
// sesji nie ma czym otworzyć. Wywołujący nazywa odmowę: NawigujDoSymbolu
// przez pole kontraktu `ServerAvailable: false`, Refaktoryzuj zdaniem błędu.
func serwerJezykaPliku(sciezka string) (zewnetrzne.Narzedzie, bool) {
	switch strings.ToLower(filepath.Ext(sciezka)) {
	case ".ts", ".tsx", ".mts", ".cts", ".js", ".jsx", ".mjs", ".cjs":
		return narzedzieSerweraTypeScript, false
	default:
		return narzedzieGopls, true
	}
}

// ── Analizatory repozytorium ────────────────────────────────────────────────

// analizaPozaGo opisuje jedno wywołanie programu analizy: program, argumenty
// i sposób odczytu wyniku z jego strumieni.
type analizaPozaGo struct {
	narzedzie zewnetrzne.Narzedzie
	argumenty []string
	// czytaj przekłada wyjście programu na zgłoszenia kontraktu, czytając oba
	// jego strumienie wyjściowe.
	czytaj func(wyjscie, diagnostyka, korzen string) []shared.DeveloperDiagnostic
}

// rozszerzeniaPythona i rozszerzeniaArkuszy nazywają pliki, które mają swój
// analizator. Wykaz jest wzięty z tego, co program naprawdę czyta.
var (
	rozszerzeniaGo          = []string{".go"}
	rozszerzeniaPythona     = []string{".py", ".pyi"}
	rozszerzeniaArkuszy     = []string{".css", ".scss", ".less"}
	rozszerzeniaTypeScriptu = []string{".ts", ".tsx", ".mts", ".cts", ".js", ".jsx", ".mjs", ".cjs"}
)

// analizyRepozytorium dobiera programy analizy, które w danym żądaniu mają
// w repozytorium co sprawdzić, po rozszerzeniach wskazanych plików.
func analizyRepozytorium(korzen string, sciezki []string) []analizaPozaGo {
	wskazane := make([]string, 0, len(sciezki))
	for _, sciezka := range sciezki {
		if tresc := strings.TrimSpace(sciezka); tresc != "" {
			wskazane = append(wskazane, tresc)
		}
	}
	caleDrzewo := len(wskazane) == 0

	analizy := make([]analizaPozaGo, 0, 4)

	plikiGo := sciezkiORozszerzeniu(wskazane, rozszerzeniaGo)
	if caleDrzewo || len(plikiGo) > 0 {
		cele := plikiGo
		if caleDrzewo {
			cele = []string{"./..."}
		}
		// `--output.json.path stdout` zastępuje usunięty w nowszym wydaniu
		// zapis `--out-format=json`.
		analizy = append(analizy, analizaPozaGo{
			narzedzie: narzedzieGolangciLint,
			argumenty: append([]string{"run", "--output.json.path", "stdout",
				"--issues-exit-code=0"}, cele...),
			czytaj: func(wyjscie, _, korzen string) []shared.DeveloperDiagnostic {
				return zgloszeniaAnalizy(wyjscie, korzen)
			},
		})
		analizy = append(analizy, analizaPozaGo{
			narzedzie: narzedzieStaticcheck,
			argumenty: append([]string{"-f", "json"}, cele...),
			czytaj:    zgloszeniaStaticcheck,
		})
	}

	skrypty := sciezkiORozszerzeniu(wskazane, rozszerzeniaTypeScriptu)
	if caleDrzewo || len(skrypty) > 0 {
		cele := skrypty
		if caleDrzewo {
			cele = []string{"."}
		}
		analizy = append(analizy, analizaPozaGo{
			narzedzie: narzedzieEslint,
			argumenty: append([]string{"--format", "json",
				"--no-error-on-unmatched-pattern"}, cele...),
			czytaj: zgloszeniaEslint,
		})
	}

	pythony := sciezkiORozszerzeniu(wskazane, rozszerzeniaPythona)
	if caleDrzewo || len(pythony) > 0 {
		cele := pythony
		if caleDrzewo {
			cele = []string{"."}
		}
		analizy = append(analizy, analizaPozaGo{
			narzedzie: narzedzieRuff,
			argumenty: append([]string{"check", "--output-format", "json",
				"--no-cache", "--force-exclude"}, cele...),
			czytaj: zgloszeniaRuff,
		})
	}

	// Stylelint bez zestawu reguł odmawia analizy — arkusze analizuje się
	// z konfiguracją repozytorium.
	arkusze := sciezkiORozszerzeniu(wskazane, rozszerzeniaArkuszy)
	if (caleDrzewo || len(arkusze) > 0) && repozytoriumNiesieKonfiguracje(korzen,
		konfiguracjeStylelinta) {
		cele := arkusze
		if caleDrzewo {
			cele = []string{"**/*.css", "**/*.scss", "**/*.less"}
		}
		analizy = append(analizy, analizaPozaGo{
			narzedzie: narzedzieStylelint,
			argumenty: append([]string{"--formatter", "json",
				"--allow-empty-input"}, cele...),
			czytaj: zgloszeniaStylelinta,
		})
	}

	// Literówka nie ma języka, więc `typos` dostaje to, co wskazano, bez
	// odsiewania po rozszerzeniu.
	celeLiterowek := wskazane
	if caleDrzewo {
		celeLiterowek = []string{"."}
	}
	analizy = append(analizy, analizaPozaGo{
		narzedzie: narzedzieTypos,
		argumenty: append([]string{"--format", "json"}, celeLiterowek...),
		czytaj:    zgloszeniaLiterowek,
	})

	return analizy
}

// przeprowadzAnalize uruchamia jeden program analizy dla repozytorium i mówi,
// czy dostarczył odczytywalny wynik.
func (a *adapterDevelopera) przeprowadzAnalize(ctx context.Context, okno session.Okno,
	korzen string, analiza analizaPozaGo) ([]shared.DeveloperDiagnostic, bool) {

	wynik, err := a.wolajNarzedzieWarsztatu(ctx, okno, analiza.narzedzie,
		analiza.argumenty, korzen, czasAnalizyStatycznej)
	if brakNarzedziaWarsztatu(err) {
		return nil, false
	}
	uwagi := analiza.czytaj(string(wynik.Wyjscie), wynik.Diagnostyka, korzen)
	if err != nil && len(uwagi) == 0 {
		return nil, false
	}
	return uwagi, true
}

// sciezkiORozszerzeniu odsiewa ze wskazanych ścieżek repozytorium wyłącznie
// te, których rozszerzenie należy do żądanego zestawu.
func sciezkiORozszerzeniu(sciezki, rozszerzenia []string) []string {
	wybrane := make([]string, 0, len(sciezki))
	for _, sciezka := range sciezki {
		koncowka := strings.ToLower(filepath.Ext(sciezka))
		for _, rozszerzenie := range rozszerzenia {
			if koncowka == rozszerzenie {
				wybrane = append(wybrane, sciezka)
				break
			}
		}
	}
	return wybrane
}

// konfiguracjeStylelinta wylicza nazwy, pod którymi Stylelint szuka swojego
// zestawu reguł w korzeniu repozytorium.
var konfiguracjeStylelinta = []string{
	".stylelintrc", ".stylelintrc.json", ".stylelintrc.yml", ".stylelintrc.yaml",
	".stylelintrc.js", ".stylelintrc.cjs", ".stylelintrc.mjs",
	"stylelint.config.js", "stylelint.config.cjs", "stylelint.config.mjs",
}

// repozytoriumNiesieKonfiguracje orzeka, czy w korzeniu repozytorium leży
// którakolwiek z wymienionych konfiguracji.
func repozytoriumNiesieKonfiguracje(korzen string, nazwy []string) bool {
	for _, nazwa := range nazwy {
		if _, err := os.Stat(filepath.Join(korzen, nazwa)); err == nil {
			return true
		}
	}
	return false
}

// wyjscieRuff jest kształtem pojedynczego zgłoszenia odpowiedzi programu
// `ruff check --output-format json`.
type wyjscieRuff struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
	Severity string `json:"severity"`
	Location struct {
		Row    int `json:"row"`
		Column int `json:"column"`
	} `json:"location"`
	Fix *struct {
		Message string `json:"message"`
	} `json:"fix"`
}

// zgloszeniaRuff przekłada wynik programu `Ruff` na wykaz zgłoszeń kontraktu
// wraz z kodem reguły i kolumną.
func zgloszeniaRuff(wyjscie, _, korzen string) []shared.DeveloperDiagnostic {
	zgloszenia := make([]shared.DeveloperDiagnostic, 0, 16)
	var uwagi []wyjscieRuff
	if err := json.Unmarshal([]byte(strings.TrimSpace(wyjscie)), &uwagi); err != nil {
		return nil
	}
	for _, uwaga := range uwagi {
		zgloszenie := shared.DeveloperDiagnostic{
			Path:     sciezkaWzgledemKorzenia(uwaga.Filename, korzen),
			Line:     uwaga.Location.Row,
			Severity: wagaAnalizy(uwaga.Severity),
			Message:  uwaga.Message,
			Source:   wskaznikTekstu(narzedzieRuff.Program),
		}
		if uwaga.Code != "" {
			zgloszenie.Code = wskaznikTekstu(uwaga.Code)
		}
		if uwaga.Location.Column > 0 {
			zgloszenie.Column = wskaznikLiczby(uwaga.Location.Column)
		}
		if uwaga.Fix != nil {
			zgloszenie.FixAvailable = wskaznikPrawdy(true)
		}
		zgloszenia = append(zgloszenia, zgloszenie)
	}
	return zgloszenia
}

// wyjscieStaticcheck jest kształtem pojedynczego wiersza wyjścia programu
// `staticcheck -f json` — jeden dokument JSON na wiersz, nie tablica.
type wyjscieStaticcheck struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Location struct {
		File   string `json:"file"`
		Line   int    `json:"line"`
		Column int    `json:"column"`
	} `json:"location"`
}

// zgloszeniaStaticcheck przekłada wynik programu `staticcheck` na wykaz
// zgłoszeń kontraktu wraz z kodem reguły i kolumną.
func zgloszeniaStaticcheck(wyjscie, _, korzen string) []shared.DeveloperDiagnostic {
	zgloszenia := make([]shared.DeveloperDiagnostic, 0, 16)
	for _, wiersz := range strings.Split(wyjscie, "\n") {
		tresc := strings.TrimSpace(wiersz)
		if tresc == "" {
			continue
		}
		var uwaga wyjscieStaticcheck
		if err := json.Unmarshal([]byte(tresc), &uwaga); err != nil || uwaga.Location.File == "" {
			continue
		}
		zgloszenie := shared.DeveloperDiagnostic{
			Path:     sciezkaWzgledemKorzenia(uwaga.Location.File, korzen),
			Line:     uwaga.Location.Line,
			Severity: wagaAnalizy(uwaga.Severity),
			Message:  uwaga.Message,
			Source:   wskaznikTekstu(narzedzieStaticcheck.Program),
		}
		if uwaga.Code != "" {
			zgloszenie.Code = wskaznikTekstu(uwaga.Code)
		}
		if uwaga.Location.Column > 0 {
			zgloszenie.Column = wskaznikLiczby(uwaga.Location.Column)
		}
		zgloszenia = append(zgloszenia, zgloszenie)
	}
	return zgloszenia
}

// wyjscieEslint jest kształtem pojedynczego pliku odpowiedzi programu
// `eslint --format json` wraz z jego zgłoszeniami.
type wyjscieEslint struct {
	FilePath string `json:"filePath"`
	Messages []struct {
		RuleId   *string         `json:"ruleId"`
		Severity int             `json:"severity"`
		Message  string          `json:"message"`
		Line     int             `json:"line"`
		Column   int             `json:"column"`
		Fix      json.RawMessage `json:"fix"`
	} `json:"messages"`
}

// zgloszeniaEslint przekłada wynik programu `ESLint` na wykaz zgłoszeń
// kontraktu. Program bez pliku nastaw w repozytorium odmawia analizy przed
// wypisaniem jakiegokolwiek JSON-u na wyjściu — odmowa wraca jednym
// zgłoszeniem NAZYWAJĄCYM przyczynę, nie cichą pustą listą.
func zgloszeniaEslint(wyjscie, diagnostyka, korzen string) []shared.DeveloperDiagnostic {
	zgloszenia := make([]shared.DeveloperDiagnostic, 0, 16)
	var pliki []wyjscieEslint
	if err := json.Unmarshal([]byte(strings.TrimSpace(wyjscie)), &pliki); err != nil {
		if tresc := strings.TrimSpace(diagnostyka); tresc != "" {
			zgloszenia = append(zgloszenia, shared.DeveloperDiagnostic{
				Path:     korzen,
				Line:     1,
				Severity: shared.ProblemSeverityError,
				Message:  "ESLint odmówił analizy: " + skrocDiagnostyke(tresc, nil),
				Source:   wskaznikTekstu(narzedzieEslint.Program),
			})
		}
		return zgloszenia
	}
	for _, plik := range pliki {
		for _, uwaga := range plik.Messages {
			wiersz := uwaga.Line
			if wiersz <= 0 {
				wiersz = 1
			}
			zgloszenie := shared.DeveloperDiagnostic{
				Path:     sciezkaWzgledemKorzenia(plik.FilePath, korzen),
				Line:     wiersz,
				Severity: wagaEslint(uwaga.Severity),
				Message:  uwaga.Message,
				Source:   wskaznikTekstu(narzedzieEslint.Program),
			}
			if uwaga.RuleId != nil && *uwaga.RuleId != "" {
				zgloszenie.Code = wskaznikTekstu(*uwaga.RuleId)
			}
			if uwaga.Column > 0 {
				zgloszenie.Column = wskaznikLiczby(uwaga.Column)
			}
			if len(uwaga.Fix) > 0 {
				zgloszenie.FixAvailable = wskaznikPrawdy(true)
			}
			zgloszenia = append(zgloszenia, zgloszenie)
		}
	}
	return zgloszenia
}

// wagaEslint przekłada wagę zgłoszenia ESLint (1 ostrzeżenie, 2 błąd) na wagę
// kontraktu.
func wagaEslint(waga int) shared.ProblemSeverity {
	if waga >= 2 {
		return shared.ProblemSeverityError
	}
	return shared.ProblemSeverityWarning
}

// wyjscieStylelinta jest kształtem pojedynczego pliku odpowiedzi programu
// `stylelint --formatter json` wraz z jego ostrzeżeniami.
type wyjscieStylelinta struct {
	Source   string `json:"source"`
	Warnings []struct {
		Line     int    `json:"line"`
		Column   int    `json:"column"`
		Rule     string `json:"rule"`
		Severity string `json:"severity"`
		Text     string `json:"text"`
	} `json:"warnings"`
}

// zgloszeniaStylelinta przekłada wynik programu `Stylelint` na wykaz zgłoszeń
// kontraktu, czytany z obu strumieni wyjścia.
func zgloszeniaStylelinta(wyjscie, diagnostyka, korzen string) []shared.DeveloperDiagnostic {
	zgloszenia := make([]shared.DeveloperDiagnostic, 0, 16)
	var pliki []wyjscieStylelinta
	if err := json.Unmarshal([]byte(strings.TrimSpace(wyjscie)), &pliki); err != nil {
		if err := json.Unmarshal([]byte(strings.TrimSpace(diagnostyka)), &pliki); err != nil {
			return nil
		}
	}
	for _, plik := range pliki {
		for _, uwaga := range plik.Warnings {
			zgloszenie := shared.DeveloperDiagnostic{
				Path:     sciezkaWzgledemKorzenia(plik.Source, korzen),
				Line:     uwaga.Line,
				Severity: wagaAnalizy(uwaga.Severity),
				Message:  uwaga.Text,
				Source:   wskaznikTekstu(narzedzieStylelint.Program),
			}
			if uwaga.Rule != "" {
				zgloszenie.Code = wskaznikTekstu(uwaga.Rule)
			}
			if uwaga.Column > 0 {
				zgloszenie.Column = wskaznikLiczby(uwaga.Column)
			}
			zgloszenia = append(zgloszenia, zgloszenie)
		}
	}
	return zgloszenia
}

// wyjscieLiterowki jest kształtem jednego wiersza wyjścia programu
// `typos --format json`, zapisanego jako osobny dokument JSON.
type wyjscieLiterowki struct {
	Rodzaj       string   `json:"type"`
	Sciezka      string   `json:"path"`
	Wiersz       int      `json:"line_num"`
	Przesuniecie int      `json:"byte_offset"`
	Literowka    string   `json:"typo"`
	Poprawki     []string `json:"corrections"`
}

// zgloszeniaLiterowek przekłada wynik programu `typos` na wykaz zgłoszeń
// kontraktu, pomijając wiersze innego rodzaju niż literówka.
func zgloszeniaLiterowek(wyjscie, _, korzen string) []shared.DeveloperDiagnostic {
	zgloszenia := make([]shared.DeveloperDiagnostic, 0, 16)
	for _, wiersz := range strings.Split(wyjscie, "\n") {
		tresc := strings.TrimSpace(wiersz)
		if tresc == "" {
			continue
		}
		var uwaga wyjscieLiterowki
		if err := json.Unmarshal([]byte(tresc), &uwaga); err != nil {
			continue
		}
		if uwaga.Rodzaj != "typo" {
			continue
		}
		zdanie := "literówka „" + uwaga.Literowka + "”"
		if len(uwaga.Poprawki) > 0 {
			zdanie += "; poprawnie: " + strings.Join(uwaga.Poprawki, ", ")
		}
		zgloszenie := shared.DeveloperDiagnostic{
			Path:     sciezkaWzgledemKorzenia(uwaga.Sciezka, korzen),
			Line:     uwaga.Wiersz,
			Severity: shared.ProblemSeverityWarning,
			Message:  zdanie,
			Source:   wskaznikTekstu(narzedzieTypos.Program),
		}
		// `byte_offset` liczy się od zera w obrębie wiersza, a kontrakt liczy
		// kolumny od jedynki.
		zgloszenie.Column = wskaznikLiczby(uwaga.Przesuniecie + 1)
		zgloszenia = append(zgloszenia, zgloszenie)
	}
	return zgloszenia
}

// sciezkaWzgledemKorzenia dopełnia wskazanie względne do ścieżki bezwzględnej.
// Programy analizy oddają raz jedno, raz drugie, a okno ma dostawać jedną postać.
func sciezkaWzgledemKorzenia(sciezka, korzen string) string {
	tresc := strings.TrimSpace(sciezka)
	if tresc == "" || filepath.IsAbs(tresc) {
		return tresc
	}
	return filepath.Join(korzen, tresc)
}

// Refaktoryzuj obsługuje komendę `developer.refactor.apply`: przeprowadza
// refaktoryzację semantyczną serwerem języka, domyślnie w trybie podglądu.
func (a *adapterDevelopera) Refaktoryzuj(ctx context.Context,
	z shared.DeveloperRefactorApplyRequest) (shared.DeveloperRefactorApplyResponse, error) {

	okno, sciezka, err := a.plikOkna(z.WindowId, z.Path)
	if err != nil {
		return shared.DeveloperRefactorApplyResponse{}, err
	}
	zastosuj := z.Preview != nil && !*z.Preview
	if zastosuj {
		if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "zastosowanie refaktoryzacji"); err != nil {
			return shared.DeveloperRefactorApplyResponse{}, err
		}
	}

	argumenty, err := argumentyRefaktoryzacji(z, sciezka, zastosuj)
	if err != nil {
		return shared.DeveloperRefactorApplyResponse{}, err
	}

	if serwer, jednorazowy := serwerJezykaPliku(sciezka); !jednorazowy {
		return shared.DeveloperRefactorApplyResponse{}, bladZasobuDevelopera(
			"serwer nie ma czym przeprowadzić refaktoryzacji semantycznej pliku " +
				filepath.Base(sciezka) + ": serwer języka tego pliku (" + serwer.Program +
				") rozmawia wyłącznie sesją protokołu LSP, a warstwa językowa modułu " +
				"pyta serwery pojedynczym wywołaniem. Droga, która działa: zmiana nazwy " +
				"po składni komendą `developer.grep.replace` wzorcem z metazmienną")
	}

	wynik, err := a.wolajNarzedzieWarsztatu(ctx, okno, narzedzieGopls, argumenty,
		filepath.Dir(sciezka), czasSerweraJezyka)
	if err != nil {
		if brakNarzedziaWarsztatu(err) {
			return shared.DeveloperRefactorApplyResponse{}, bladZasobuDevelopera(
				"serwer nie ma programu gopls — refaktoryzacja semantyczna nie ma czym się " +
					"wykonać; naprawa: " + narzedzieGopls.Pakiet)
		}
		return shared.DeveloperRefactorApplyResponse{}, bladWykonaniaDevelopera(
			"refaktoryzacja odmówiła: " + skrocDiagnostyke(wynik.Diagnostyka, err))
	}

	zmiany, sciezki := zmianyRefaktoryzacji(string(wynik.Wyjscie), sciezka)
	return shared.DeveloperRefactorApplyResponse{
		Edits:        zmiany,
		Applied:      zastosuj,
		ChangedPaths: sciezki,
	}, nil
}

// argumentyRefaktoryzacji składa argumenty wywołania serwera języka dla
// żądanego rodzaju refaktoryzacji.
func argumentyRefaktoryzacji(z shared.DeveloperRefactorApplyRequest, sciezka string,
	zastosuj bool) ([]string, error) {

	kolumna := z.Column
	if kolumna < 1 {
		kolumna = 1
	}
	miejsce := sciezka + ":" + strconv.Itoa(z.Line) + ":" + strconv.Itoa(kolumna)

	switch z.Kind {
	case shared.RefactorKindRename:
		nowa := ""
		if z.NewName != nil {
			nowa = strings.TrimSpace(*z.NewName)
		}
		if nowa == "" {
			return nil, bladZadaniaDevelopera("zmiana nazwy symbolu wymaga nowej nazwy")
		}
		argumenty := []string{"rename"}
		if zastosuj {
			argumenty = append(argumenty, "-w")
		} else {
			argumenty = append(argumenty, "-d")
		}
		return append(argumenty, miejsce, nowa), nil

	case shared.RefactorKindOrganizeImports:
		argumenty := []string{"imports"}
		if zastosuj {
			argumenty = append(argumenty, "-w")
		} else {
			argumenty = append(argumenty, "-d")
		}
		return append(argumenty, sciezka), nil

	case shared.RefactorKindExtractFunction, shared.RefactorKindExtractVariable,
		shared.RefactorKindInline, shared.RefactorKindMove:
		argumenty := []string{"codeaction", "-kind", kodCzynnosciRefaktoryzacji(z.Kind)}
		if zastosuj {
			argumenty = append(argumenty, "-w")
		} else {
			argumenty = append(argumenty, "-d")
		}
		return append(argumenty, miejsce), nil

	default:
		return nil, bladZadaniaDevelopera("nieznany rodzaj refaktoryzacji: " + string(z.Kind))
	}
}

// kodCzynnosciRefaktoryzacji przekłada rodzaj refaktoryzacji z kontraktu na
// kod czynności protokołu LSP.
func kodCzynnosciRefaktoryzacji(rodzaj shared.RefactorKind) string {
	switch rodzaj {
	case shared.RefactorKindExtractFunction:
		return "refactor.extract.function"
	case shared.RefactorKindExtractVariable:
		return "refactor.extract.variable"
	case shared.RefactorKindInline:
		return "refactor.inline"
	default:
		return "refactor.move"
	}
}

// zmianyRefaktoryzacji rozbiera różnicę zunifikowaną oddaną przez serwer
// języka na wykaz zmian tekstu i ścieżek plików.
func zmianyRefaktoryzacji(roznica, sciezkaDomyslna string) ([]shared.DeveloperTextEdit, []string) {
	zmiany := make([]shared.DeveloperTextEdit, 0, 8)
	sciezki := make([]string, 0, 4)
	widziane := map[string]bool{}

	biezaca := sciezkaDomyslna
	for _, wiersz := range strings.Split(roznica, "\n") {
		switch {
		case strings.HasPrefix(wiersz, "+++ "):
			biezaca = sciezkaZNaglowkaRoznicy(wiersz)
			if biezaca != "" && !widziane[biezaca] {
				widziane[biezaca] = true
				sciezki = append(sciezki, biezaca)
			}
		case strings.HasPrefix(wiersz, "@@"):
			if zmiana, jest := zmianaZFragmentu(wiersz, biezaca); jest {
				zmiany = append(zmiany, zmiana)
			}
		}
	}
	return zmiany, sciezki
}

// sciezkaZNaglowkaRoznicy wyjmuje ścieżkę pliku z nagłówka różnicy postaci
// `+++ b/plik.go`, zdejmując przedrostki katalogów.
func sciezkaZNaglowkaRoznicy(wiersz string) string {
	tresc := strings.TrimSpace(strings.TrimPrefix(wiersz, "+++ "))
	if tabulator := strings.IndexByte(tresc, '\t'); tabulator > 0 {
		tresc = tresc[:tabulator]
	}
	return strings.TrimPrefix(strings.TrimPrefix(tresc, "b/"), "a/")
}

// zmianaZFragmentu rozbiera nagłówek fragmentu różnicy postaci
// `@@ -12,3 +12,5 @@` na początek i długość zakresu zmiany.
func zmianaZFragmentu(wiersz, sciezka string) (shared.DeveloperTextEdit, bool) {
	pola := strings.Fields(wiersz)
	if len(pola) < 3 || !strings.HasPrefix(pola[1], "-") {
		return shared.DeveloperTextEdit{}, false
	}
	zakres := strings.SplitN(strings.TrimPrefix(pola[1], "-"), ",", 2)
	poczatek, err := strconv.Atoi(zakres[0])
	if err != nil {
		return shared.DeveloperTextEdit{}, false
	}
	dlugosc := 1
	if len(zakres) == 2 {
		if liczba, err := strconv.Atoi(zakres[1]); err == nil {
			dlugosc = liczba
		}
	}
	return shared.DeveloperTextEdit{
		Path:      sciezka,
		StartLine: poczatek,
		EndLine:   poczatek + dlugosc,
		NewText:   "",
	}, true
}

// plikOkna sprowadza wskazanie klienta do ścieżki bezwzględnej w obszarze
// roboczym okna developera i sprawdza jej istnienie.
func (a *adapterDevelopera) plikOkna(oknoKod, wskazanie string) (session.Okno, string, error) {
	okno, err := a.oknoDevelopera(oknoKod)
	if err != nil {
		return session.Okno{}, "", err
	}
	sciezka, err := sciezkaWObszarze(a.korzenieOkna(okno), wskazanie, plikIstnieje)
	if err != nil {
		return session.Okno{}, "", err
	}
	return okno, sciezka, nil
}

// plikIstnieje mówi, czy pod wskazaną ścieżką w systemie plików stoi
// jakikolwiek plik, dostępny do odczytu.
func plikIstnieje(sciezka string) bool {
	_, err := os.Stat(sciezka)
	return err == nil
}

// skrocDiagnostyke składa zdanie o niepowodzeniu programu z jego diagnostyki,
// a gdy jej nie ma — z treści błędu uruchomienia.
func skrocDiagnostyke(diagnostyka string, err error) string {
	tresc := strings.TrimSpace(diagnostyka)
	if tresc == "" && err != nil {
		tresc = err.Error()
	}
	const najdluzej = 2000
	if len(tresc) > najdluzej {
		return tresc[:najdluzej] + "…"
	}
	return tresc
}
