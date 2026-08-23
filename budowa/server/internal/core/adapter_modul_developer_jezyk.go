// Odpowiedzialność pliku: cztery komendy warstwy językowej Code Editora —
// `developer.symbol.navigate` (przejdź do definicji, implementacji, wystąpień,
// symbole dokumentu), `developer.format.run` (formatowanie), `developer.lint.get`
// (analiza statyczna) i `developer.refactor.apply` (refaktoryzacje).
//
// ── Skąd bierze się wiedza o kodzie ──────────────────────────────────────────
// Tu jako jedyne w module wołamy programy serwera, bo program JEST tą wiedzą.
// `gopls` zna typy repozytorium i jego graf odwołań; napisanie tego od nowa w Go
// znaczyłoby napisanie drugiego kompilatora. `gofmt` i `prettier` znają styl,
// `golangci-lint` zna reguły. Wszystkie stoją na serwerze i wchodzą do sondy
// zależności (`zaleznosci_zewnetrzne.go`).
//
// ── Czego brak programu NIE robi ─────────────────────────────────────────────
// Nie zamienia odpowiedzi w odmowę. Kontrakt niesie `serverAvailable`
// i `linterAvailable` właśnie po to: pusty wykaz symboli przy `serverAvailable:
// false` znaczy „nie miałem czym sprawdzić”, a przy `true` — „sprawdziłem
// i nie ma”. To są dwa różne zdania i okno ma je rozróżniać, bo pierwsze
// naprawia się instalacją na serwerze, a drugie — poprawką w kodzie.
//
// ── Jeden serwer języka na wywołanie, a nie utrzymywana sesja ────────────────
// `gopls` woła się tu w trybie jednorazowym (`gopls definition …`), a nie jako
// utrzymywany proces mostkujący LSP na WebSocket. Powód: komenda kontraktu jest
// pytaniem i odpowiedzią, bez stanu między nimi, a proces utrzymywany na okno
// wymagałby własnego rejestru uchwytów, własnego sprzątania i własnego
// przerywania — trzeciego takiego mechanizmu w module obok budowania
// i debugowania. Most LSP z rozdz. 7.1 opracowania jest rozszerzeniem tej drogi,
// a nie jej zaprzeczeniem: kontrakt pozostaje ten sam.
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
	// czasFormatowania jest granicą formatowania jednego pliku.
	czasFormatowania = 20 * time.Second
	// czasAnalizyStatycznej jest granicą przebiegu lintera po repozytorium.
	czasAnalizyStatycznej = 180 * time.Second
	// najwiecejZgloszenAnalizy chroni Dev Tools przed wykazem, którego nikt nie
	// przeczyta, gdy linter zgłosi tysiące uwag w nieuporządkowanym repozytorium.
	najwiecejZgloszenAnalizy = 500
)

// NawigujDoSymbolu obsługuje `developer.symbol.navigate`.
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

	wynik, err := a.wolajNarzedzieWarsztatu(ctx, okno, narzedzieGopls, argumenty,
		filepath.Dir(sciezka), czasSerweraJezyka)
	if err != nil {
		if brakNarzedziaWarsztatu(err) {
			return shared.DeveloperSymbolNavigateResponse{
				Symbols: []shared.DeveloperSymbol{}, ServerAvailable: false}, nil
		}
		// Serwer języka odpowiada niezerowym kodem także wtedy, gdy pod
		// kursorem nie ma symbolu. To nie jest awaria: pusty wykaz przy
		// `serverAvailable: true` mówi dokładnie to, co zaszło.
		return shared.DeveloperSymbolNavigateResponse{
			Symbols: []shared.DeveloperSymbol{}, ServerAvailable: true}, nil
	}

	symbole := symboleZOdpowiedziSerwera(z.Kind, string(wynik.Wyjscie), sciezka)
	return shared.DeveloperSymbolNavigateResponse{Symbols: symbole, ServerAvailable: true}, nil
}

// symbolGopls jest kształtem odpowiedzi `gopls definition -json`.
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

// symboleZOdpowiedziSerwera przekłada wyjście serwera języka na kontrakt.
//
// Dwa kształty wyjścia, bo `gopls` mówi dwoma językami: `definition -json` daje
// dokument JSON, a `references` i `symbols` — wiersze tekstu. Rozbiór idzie po
// rodzaju pytania, a nie po zgadywaniu z treści.
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

// oberwijZakres zostawia z `12-15` samo `12`.
func oberwijZakres(wartosc string) string {
	if myslnik := strings.IndexByte(wartosc, '-'); myslnik > 0 {
		return wartosc[:myslnik]
	}
	return wartosc
}

// pierwszeSlowoOpisu bierze z opisu serwera nazwę symbolu.
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

// sciezkaZOdwolania zdejmuje przedrostek `file://` z odwołania serwera języka.
func sciezkaZOdwolania(odwolanie string) string {
	return strings.TrimPrefix(odwolanie, "file://")
}

// Formatuj obsługuje `developer.format.run`.
//
// Formater dobiera się po rozszerzeniu pliku, bo styl należy do języka, nie do
// Operatora. Treść w żądaniu wygrywa z treścią na dysku: Code Editor formatuje
// bufor, którego jeszcze nie zapisano, i to on jest przedmiotem czynności.
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

	// Treść idzie do formatera plikiem tymczasowym w katalogu roboczym okna,
	// a nie strumieniem wejścia: port uruchamiacza podaje procesowi wyjście
	// i diagnostykę, lecz nie wejście, a formatery czytające ze standardowego
	// wejścia i tak nie znają wtedy ścieżki pliku, więc gubią reguły
	// repozytorium zależne od położenia.
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

// formaterPliku dobiera program formatujący po rozszerzeniu pliku.
//
// Dla Go pierwszeństwo ma `goimports`: robi to samo co `gofmt`, a dodatkowo
// porządkuje wykaz importów, którego ręczne pilnowanie jest najczęstszym
// powodem niekompilującego się pliku po refaktoryzacji. Gdy go na serwerze nie
// ma, zostaje `gofmt`.
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

// plikRoboczyFormatowania odkłada treść bufora obok pliku źródłowego.
//
// Obok, a nie w katalogu tymczasowym systemu: reguły formatowania zależą od
// położenia w repozytorium (`.editorconfig`, `.prettierrc`, moduł Go), więc plik
// przeniesiony gdzie indziej sformatowałby się wedle innych reguł niż ten,
// którego dotyczy czynność.
func plikRoboczyFormatowania(sciezka, tresc string) (string, func(), error) {
	roboczy := filepath.Join(filepath.Dir(sciezka),
		"."+filepath.Base(sciezka)+".danaco-format"+filepath.Ext(sciezka))
	if err := os.WriteFile(roboczy, []byte(tresc), 0o600); err != nil {
		return "", func() {}, bladWykonaniaDevelopera(
			"nie można odłożyć treści do formatowania: " + err.Error())
	}
	return roboczy, func() { _ = os.Remove(roboczy) }, nil
}

// zawezDoZakresu składa wynik formatowania zaznaczenia.
//
// Formatery pracują na całym pliku, a kontrakt dopuszcza zakres wierszy. Wynik
// powstaje przez wzięcie z formatowanej treści wyłącznie zakresu wskazanego
// przez Operatora — reszta pliku zostaje nietknięta, bo Operator prosił
// o zaznaczenie, a nie o cały plik.
func zawezDoZakresu(przed, po string, odWiersza, doWiersza *int) string {
	if odWiersza == nil || doWiersza == nil {
		return po
	}
	wierszePrzed := strings.Split(przed, "\n")
	wierszePo := strings.Split(po, "\n")
	od, do := *odWiersza, *doWiersza
	if od < 1 || do < od || do > len(wierszePrzed) || len(wierszePo) < do {
		// Formatowanie przesunęło wiersze tak, że zakres przestał się zgadzać —
		// zawężenie dałoby wtedy plik posklejany z dwóch różnych stanów.
		return po
	}
	zlozony := make([]string, 0, len(wierszePrzed))
	zlozony = append(zlozony, wierszePrzed[:od-1]...)
	zlozony = append(zlozony, wierszePo[od-1:do]...)
	zlozony = append(zlozony, wierszePrzed[do:]...)
	return strings.Join(zlozony, "\n")
}

// AnalizaStatyczna obsługuje `developer.lint.get`.
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
	cele := celeAnalizy(korzenie[0], z.Paths)

	wynik, err := a.wolajNarzedzieWarsztatu(ctx, okno, narzedzieGolangciLint,
		append([]string{"run", "--out-format=json", "--issues-exit-code=0"}, cele...),
		korzenie[0], czasAnalizyStatycznej)
	if err != nil {
		if brakNarzedziaWarsztatu(err) {
			return shared.DeveloperLintGetResponse{
				Diagnostics: []shared.DeveloperDiagnostic{}, LinterAvailable: false}, nil
		}
		return shared.DeveloperLintGetResponse{}, bladWykonaniaDevelopera(
			"analiza statyczna odmówiła: " + skrocDiagnostyke(wynik.Diagnostyka, err))
	}

	zgloszenia := zgloszeniaAnalizy(string(wynik.Wyjscie), korzenie[0])
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

// celeAnalizy składa wskazania dla lintera. Puste żądanie znaczy całe drzewo
// repozytorium.
func celeAnalizy(korzen string, sciezki []string) []string {
	cele := make([]string, 0, len(sciezki))
	for _, sciezka := range sciezki {
		if tresc := strings.TrimSpace(sciezka); tresc != "" {
			cele = append(cele, tresc)
		}
	}
	if len(cele) == 0 {
		return []string{"./..."}
	}
	return cele
}

// wyjscieAnalizy jest kształtem odpowiedzi `golangci-lint run --out-format=json`.
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

// zgloszeniaAnalizy przekłada wyjście lintera na zgłoszenia kontraktu.
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

// Refaktoryzuj obsługuje `developer.refactor.apply`.
//
// Podgląd jest domyślny. Refaktoryzacja semantyczna dotyka wielu plików naraz,
// a Operator ma prawo zobaczyć różnicę, zanim ją przyjmie — dlatego zmiana
// zapisuje się na dysk wyłącznie na wyraźne `preview: false`.
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

// argumentyRefaktoryzacji składa wywołanie serwera języka dla żądanego rodzaju.
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

// kodCzynnosciRefaktoryzacji przekłada rodzaj kontraktu na kod czynności LSP.
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

// zmianyRefaktoryzacji rozbiera różnicę zunifikowaną oddaną przez serwer języka.
//
// Różnica jest jedynym kształtem, w którym serwer mówi o zmianie wielu plików
// naraz; kontrakt niesie ją jako wykaz zmian tekstu, więc rozbiór idzie po
// nagłówkach `--- / +++` i fragmentach `@@`.
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

// sciezkaZNaglowkaRoznicy wyjmuje ścieżkę z wiersza `+++ b/plik.go\t…`.
func sciezkaZNaglowkaRoznicy(wiersz string) string {
	tresc := strings.TrimSpace(strings.TrimPrefix(wiersz, "+++ "))
	if tabulator := strings.IndexByte(tresc, '\t'); tabulator > 0 {
		tresc = tresc[:tabulator]
	}
	return strings.TrimPrefix(strings.TrimPrefix(tresc, "b/"), "a/")
}

// zmianaZFragmentu rozbiera nagłówek `@@ -12,3 +12,5 @@` na zakres zmiany.
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

// plikOkna sprowadza wskazanie klienta do ścieżki bezwzględnej w obszarze okna.
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

// plikIstnieje mówi, czy pod ścieżką coś stoi.
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
