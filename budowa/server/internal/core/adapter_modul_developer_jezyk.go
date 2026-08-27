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

	if _, jednorazowy := serwerJezykaPliku(sciezka); !jednorazowy {
		// Pusty wykaz przy `serverAvailable: false` mówi „nie miałem czym
		// sprawdzić" — i to jest tu prawdą. Wołanie serwera Go po plik
		// TypeScriptu oddawałoby pusty wykaz przy `true`, czyli zdanie
		// „sprawdziłem i nie ma", którego rdzeń nie ma prawa powiedzieć.
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
//
// ── Dlaczego programów jest kilka, a pole `linterAvailable` jedno ────────────
// Repozytorium bywa wielojęzyczne, a analizator zna jeden język: `golangci-lint`
// czyta Go, `Ruff` Pythona, `Stylelint` arkusze CSS, a `typos` szuka literówek
// niezależnie od języka. Zgłoszenia idą do jednego wykazu, bo kontrakt niesie
// w każdym z nich pole `source` — czytelnik widzi, który program je wystawił.
//
// Pole `linterAvailable` jest jedno i pochodzi z czasu, gdy program też był
// jeden. Znaczy tu: CHOĆ JEDEN program odpowiedział. Fałsz zostaje więc tym,
// czym był — stanem, w którym pusty wykaz zgłoszeń kłamałby, bo nie sprawdzono
// niczym. Które programy milczały, tej odpowiedzi powiedzieć nie da się bez
// zmiany kontraktu.
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

// wyjscieAnalizy jest kształtem odpowiedzi `golangci-lint run` w postaci JSON.
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

// serwerJezykaPliku dobiera serwer języka po rozszerzeniu pliku i mówi, czy
// rdzeń ma czym go zapytać JEDNYM wywołaniem.
//
// ── Dlaczego TypeScript wraca z „nie ma czym" ────────────────────────────────
// Warstwa językowa tego modułu pyta serwer pojedynczym wywołaniem
// (`gopls definition …`) i odbiera odpowiedź z jego wyjścia; powód tego wyboru
// stoi w nagłówku pliku. `gopls` taki tryb ma. Serwer języka TypeScriptu go NIE
// ma: rozmawia wyłącznie sesją protokołu LSP na strumieniu wejścia — uzgodnienie,
// otwarcie dokumentu, pytanie, zamknięcie — a jedyna droga rdzenia do procesu
// (`zewnetrzne.Wolaj`) z zamysłu na wejście procesu nie pisze, bo pisanie
// i czytanie naraz jest tą klasą zakleszczeń, której ten pakiet ma nie mieć.
//
// Dlatego pliki TypeScriptu dostają odpowiedź nazywającą brak zamiast wyniku
// z serwera Go, który tego pliku nie rozumie. Sam program jest zadeklarowany
// w wykazie zależności (`zaleznosci_zewnetrzne.go`), więc sonda startowa mówi
// o nim Operatorowi — deklaracja opisuje zakres, który bez niego nie działa,
// a nie obietnicę, że rdzeń już go woła.
func serwerJezykaPliku(sciezka string) (zewnetrzne.Narzedzie, bool) {
	switch strings.ToLower(filepath.Ext(sciezka)) {
	case ".ts", ".tsx", ".mts", ".cts", ".js", ".jsx", ".mjs", ".cjs":
		return narzedzieSerweraTypeScript, false
	default:
		return narzedzieGopls, true
	}
}

// ── Analizatory repozytorium ────────────────────────────────────────────────

// analizaPozaGo opisuje jedno wywołanie programu analizy: co uruchomić, na czym
// i jak odczytać wynik.
type analizaPozaGo struct {
	narzedzie zewnetrzne.Narzedzie
	argumenty []string
	// czytaj przekłada wyjście programu na zgłoszenia kontraktu. Dostaje oba
	// strumienie, bo jedne programy piszą wynik na wyjście, inne na diagnostykę.
	czytaj func(wyjscie, diagnostyka, korzen string) []shared.DeveloperDiagnostic
}

// rozszerzeniaPythona i rozszerzeniaArkuszy nazywają pliki, które mają swój
// analizator. Wykaz jest wzięty z tego, co program naprawdę czyta.
var (
	rozszerzeniaGo      = []string{".go"}
	rozszerzeniaPythona = []string{".py", ".pyi"}
	rozszerzeniaArkuszy = []string{".css", ".scss", ".less"}
)

// analizyRepozytorium dobiera programy, które w tym żądaniu mają co sprawdzić.
//
// ── Dlaczego dobór idzie po rozszerzeniu, a nie po tym, co stoi na maszynie ──
// Program uruchomiony tam, gdzie nie ma ani jednego pliku jego języka, nie
// milczy — ODMAWIA. `golangci-lint` w repozytorium bez plików Go kończy się
// błędem „no go files to analyze", a odmowa jednego programu nie ma prawa
// zabrać analizy pozostałym: repozytorium samego Pythona albo samego frontendu
// zostałoby wtedy bez analizy w ogóle.
//
// Żądanie bez wskazania plików obejmuje całe drzewo i wtedy każdy program
// dostaje swój katalog.
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
		// `--output.json.path stdout` jest zapisem wydania drugiego programu;
		// zapis `--out-format=json` z wydania pierwszego został z niego zdjęty
		// i program odmawia wtedy uruchomienia, nazywając nieznany parametr.
		analizy = append(analizy, analizaPozaGo{
			narzedzie: narzedzieGolangciLint,
			argumenty: append([]string{"run", "--output.json.path", "stdout",
				"--issues-exit-code=0"}, cele...),
			czytaj: func(wyjscie, _, korzen string) []shared.DeveloperDiagnostic {
				return zgloszeniaAnalizy(wyjscie, korzen)
			},
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

	// Stylelint nie ma wbudowanego zestawu reguł: uruchomienie bez konfiguracji
	// nie jest „analizą z ustawieniami domyślnymi", tylko brakiem analizy wraz
	// z odmową programu. Zestaw reguł jest rozstrzygnięciem repozytorium, a nie
	// rdzenia — więc arkusze analizuje się wyłącznie tam, gdzie repozytorium
	// swój zestaw niesie.
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

// przeprowadzAnalize uruchamia jeden program analizy i mówi, czy ZMIERZYŁ.
//
// Kod wyjścia różny od zera jest tu WYNIKIEM, nie usterką: Ruff, Stylelint
// i `typos` kończą się niezerowo dokładnie wtedy, gdy mają co zgłosić. Miarą
// pomiaru jest więc ODCZYTANA ODPOWIEDŹ, a nie kod wyjścia ani strumień, na
// który program ją napisał — program obecny, który przewrócił się, nie oddając
// nic czytelnego, NIE zmierzył niczego i jego milczenie nie ma prawa wyglądać
// jak brak zastrzeżeń.
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

// sciezkiORozszerzeniu odsiewa wskazania o żądanych rozszerzeniach.
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

// wyjscieRuff jest kształtem odpowiedzi `ruff check --output-format json`.
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

// zgloszeniaRuff przekłada wynik Ruffa na zgłoszenia kontraktu.
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

// wyjscieStylelinta jest kształtem odpowiedzi `stylelint --formatter json`.
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

// zgloszeniaStylelinta przekłada wynik Stylelinta na zgłoszenia kontraktu.
//
// Wykaz czyta się z OBU strumieni, bo Stylelint pisze go na strumień
// diagnostyczny, kiedy uwagi ma (kończy się wtedy kodem 2), a na wyjście — kiedy
// ich nie ma. Zmierzone na tej maszynie; szukanie wykazu tam, gdzie jest, a nie
// tam, gdzie wypadałoby, żeby był.
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

// wyjscieLiterowki jest kształtem jednego wiersza `typos --format json`.
//
// `typos` pisze po jednym zapisie JSON na wiersz, a nie jedną tablicę — i pisze
// tam także wiersze innego rodzaju niż literówka (na przykład pominięty plik
// binarny). Rodzaj jest więc czytany, a nie zakładany.
type wyjscieLiterowki struct {
	Rodzaj       string   `json:"type"`
	Sciezka      string   `json:"path"`
	Wiersz       int      `json:"line_num"`
	Przesuniecie int      `json:"byte_offset"`
	Literowka    string   `json:"typo"`
	Poprawki     []string `json:"corrections"`
}

// zgloszeniaLiterowek przekłada wynik `typos` na zgłoszenia kontraktu.
//
// Waga jest ostrzeżeniem, nie błędem: literówka w identyfikatorze bywa nazwą
// celowo skróconą, a program nie ma jak tego rozstrzygnąć. Podniesienie wagi
// kazałoby Operatorowi poprawiać rzeczy, których nikt tak nie oznaczył.
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

	if serwer, jednorazowy := serwerJezykaPliku(sciezka); !jednorazowy {
		return shared.DeveloperRefactorApplyResponse{}, bladZasobuDevelopera(
			"rdzeń nie ma czym przeprowadzić refaktoryzacji semantycznej pliku " +
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
