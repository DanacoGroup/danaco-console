// Plik obsługuje `terminal.script.lint`: analiza statyczna i formatowanie
// treści skryptu, osobno dla każdej powłoki, programami zewnętrznymi
// zmierzonymi na maszynie rdzenia.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// granicaCzasuAnalizy ogranicza jedno uruchomienie programu analizy. Analiza
// czyta treść jednego skryptu, więc sekundy wystarczają z zapasem; program,
// który utknął, nie ma prawa trzymać żądania gniazda.
const granicaCzasuAnalizy = 20 * time.Second

// analizatorPowloki opisuje program analizy jednej powłoki: narzędzie,
// rozszerzenie pliku, argumenty i sposób odczytu wyniku.
type analizatorPowloki struct {
	// narzedzie jest programem uruchamianym; jego nazwa wchodzi do odpowiedzi.
	narzedzie zewnetrzne.Narzedzie
	// rozszerzenie pliku tymczasowego — programy rozpoznają po nim język.
	rozszerzenie string
	// argumenty składa wiersz analizy dla wskazanego pliku.
	argumenty func(sciezka string) []string
	// czytaj przekłada wyjście programu na uwagi kontraktu, biorąc oba
	// strumienie: wyjście i diagnostykę.
	czytaj func(wyjscie, diagnostyka string) []shared.TerminalLintFinding
	// formatowanie składa wiersz formatowania; brak znaczy, że program tej
	// powłoki formatowania nie umie.
	formatowanie func(sciezka string) (zewnetrzne.Narzedzie, []string)
	// formatZPliku mówi, że program formatujący poprawia plik zamiast pisać
	// wynik na wyjście.
	formatZPliku bool
}

// narzedzieShellCheck, narzedzieShfmt i narzedziePowerShell są programami
// analizy zmierzonymi na maszynie rdzenia. Deklaracje stoją tutaj, a sonda
// startowa bierze je stąd (`zaleznosci_zewnetrzne.go`).
var (
	narzedzieShellCheck = zewnetrzne.Narzedzie{
		Nazwa: "ShellCheck", Program: "shellcheck", Pakiet: "shellcheck",
	}
	narzedzieShfmt = zewnetrzne.Narzedzie{
		Nazwa: "shfmt", Program: "shfmt", Pakiet: "shfmt",
	}
	narzedziePowerShell = zewnetrzne.Narzedzie{
		Nazwa:   "PowerShell z modułem PSScriptAnalyzer",
		Program: "pwsh",
		Pakiet:  "powershell (snap) wraz z modułem PSScriptAnalyzer",
	}
	narzedzieNode = zewnetrzne.Narzedzie{
		Nazwa: "Node.js", Program: "node", Pakiet: "nodejs",
	}
	narzedziePython = zewnetrzne.Narzedzie{
		Nazwa: "Python", Program: "python3", Pakiet: "python3",
	}
)

// analizatory jest wykazem POMIERZONYM i jawnie niepełnym: powłoka, dla której
// nie ma tu pozycji, nie ma analizatora i odpowiedź mówi to wprost.
var analizatory = map[shared.TerminalShell]analizatorPowloki{
	shared.TerminalShellBash: {
		narzedzie:    narzedzieShellCheck,
		rozszerzenie: ".sh",
		argumenty: func(sciezka string) []string {
			return []string{"-f", "json1", "-s", "bash", sciezka}
		},
		czytaj: func(wyjscie, _ string) []shared.TerminalLintFinding {
			return uwagiShellCheck(wyjscie)
		},
		formatowanie: func(sciezka string) (zewnetrzne.Narzedzie, []string) {
			return narzedzieShfmt, []string{"-i", "2", sciezka}
		},
	},
	shared.TerminalShellPowershell: {
		narzedzie:    narzedziePowerShell,
		rozszerzenie: ".ps1",
		argumenty: func(sciezka string) []string {
			return []string{"-NoLogo", "-NoProfile", "-NonInteractive", "-Command",
				"Invoke-ScriptAnalyzer -Path '" + sciezka +
					"' | Select-Object Line,Column,Severity,Message,RuleName | ConvertTo-Json -Compress -AsArray"}
		},
		czytaj: func(wyjscie, _ string) []shared.TerminalLintFinding {
			return uwagiPowerShell(wyjscie)
		},
		formatowanie: func(sciezka string) (zewnetrzne.Narzedzie, []string) {
			return narzedziePowerShell, []string{"-NoLogo", "-NoProfile", "-NonInteractive",
				"-Command", "Invoke-Formatter -ScriptDefinition (Get-Content -Raw '" + sciezka + "')"}
		},
	},
	shared.TerminalShellNode: {
		narzedzie:    narzedzieNode,
		rozszerzenie: ".js",
		argumenty:    func(sciezka string) []string { return []string{"--check", sciezka} },
		czytaj: func(_, diagnostyka string) []shared.TerminalLintFinding {
			return uwagiInterpretera(diagnostyka)
		},
	},
	shared.TerminalShellPython: {
		narzedzie:    narzedziePython,
		rozszerzenie: ".py",
		argumenty:    func(sciezka string) []string { return []string{"-m", "py_compile", sciezka} },
		czytaj: func(_, diagnostyka string) []shared.TerminalLintFinding {
			return uwagiInterpretera(diagnostyka)
		},
	},
}

// analizatorRuffa jest analizą Pythona pełną — składnią oraz regułami —
// i formatuje plik zamiast strumienia, więc odczyt wyniku idzie z pliku.
var analizatorRuffa = analizatorPowloki{
	narzedzie:    narzedzieRuff,
	rozszerzenie: ".py",
	argumenty: func(sciezka string) []string {
		return []string{"check", "--output-format", "json", "--no-cache",
			"--force-exclude", sciezka}
	},
	czytaj: func(wyjscie, _ string) []shared.TerminalLintFinding {
		return uwagiRuffa(wyjscie)
	},
	formatowanie: func(sciezka string) (zewnetrzne.Narzedzie, []string) {
		return narzedzieRuff, []string{"format", "--no-cache", sciezka}
	},
	formatZPliku: true,
}

// analizatorDlaPowloki dobiera program analizy dla powłoki: Python ma
// pierwszeństwo Ruffa nad samym interpreterem, gdy Ruff jest zainstalowany
// na maszynie.
func analizatorDlaPowloki(powloka shared.TerminalShell) (analizatorPowloki, bool) {
	if powloka == shared.TerminalShellPython && zewnetrzne.Stoi(narzedzieRuff) {
		return analizatorRuffa, true
	}
	analizator, jest := analizatory[powloka]
	return analizator, jest
}

// uwagiRuffa czyta wynik `ruff check --output-format json` i przekłada
// każdy zapis na uwagę kontraktu.
func uwagiRuffa(wyjscie string) []shared.TerminalLintFinding {
	uwagi := make([]shared.TerminalLintFinding, 0, 8)
	var zapisy []wyjscieRuff
	if err := json.Unmarshal([]byte(strings.TrimSpace(wyjscie)), &zapisy); err != nil {
		return uwagi
	}
	for _, zapis := range zapisy {
		uwaga := shared.TerminalLintFinding{
			Line:     zapis.Location.Row,
			Severity: wagaUwagi(zapis.Severity),
			Message:  zapis.Message,
		}
		if zapis.Location.Column > 0 {
			uwaga.Column = wskaznikLiczby(zapis.Location.Column)
		}
		if zapis.Code != "" {
			uwaga.Rule = wskaznikTekstu(zapis.Code)
		}
		uwagi = append(uwagi, uwaga)
	}
	return uwagi
}

// SprawdzSkrypt obsługuje `terminal.script.lint`: uruchamia program analizy
// nad plikiem tymczasowym, z obszarem i zasadami izolacji pustymi celowo.
func (a *adapterTerminala) SprawdzSkrypt(ctx context.Context,
	z shared.TerminalScriptLintRequest) (shared.TerminalScriptLintResponse, error) {

	analizator, jest := analizatorDlaPowloki(z.Shell)
	if !jest {
		// Odpowiedź nazywa brak analizatora, zamiast oddać pusty wykaz uwag
		// udający treść bez zastrzeżeń.
		return shared.TerminalScriptLintResponse{
			Findings:          []shared.TerminalLintFinding{},
			AnalyzerAvailable: false,
			Analyzer: "serwer nie zna programu analizy dla powłoki " + string(z.Shell) +
				"; analizę mają powłoki bash, powershell, node i python",
		}, nil
	}
	if !zewnetrzne.Stoi(analizator.narzedzie) {
		return shared.TerminalScriptLintResponse{
			Findings:          []shared.TerminalLintFinding{},
			AnalyzerAvailable: false,
			Analyzer:          analizator.narzedzie.Nazwa + " (" + analizator.narzedzie.Program + ")",
		}, nil
	}
	if a.uruchamiacz == nil {
		return shared.TerminalScriptLintResponse{}, bladWykonaniaTerminala(
			"serwer nie ma uruchamiacza procesów, więc nie uruchomi programu analizy")
	}

	sciezka, sprzataj, err := plikDoAnalizy(z.Content, analizator.rozszerzenie)
	if err != nil {
		return shared.TerminalScriptLintResponse{}, err
	}
	defer sprzataj()

	odpowiedz := shared.TerminalScriptLintResponse{
		AnalyzerAvailable: true,
		Analyzer:          analizator.narzedzie.Nazwa + " (" + analizator.narzedzie.Program + ")",
	}
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, session.Okno{}, session.Zasady{},
		session.Obszar{}, analizator.narzedzie, analizator.argumenty(sciezka),
		filepath.Dir(sciezka), granicaCzasuAnalizy)
	// Kod wyjścia różny od zera jest wynikiem analizy, nie usterką rdzenia.
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		odpowiedz.AnalyzerAvailable = false
		odpowiedz.Findings = []shared.TerminalLintFinding{}
		return odpowiedz, nil
	}
	odpowiedz.Findings = analizator.czytaj(string(wynik.Wyjscie), wynik.Diagnostyka)
	if odpowiedz.Findings == nil {
		odpowiedz.Findings = []shared.TerminalLintFinding{}
	}

	if z.Format != nil && *z.Format && analizator.formatowanie != nil {
		narzedzie, argumenty := analizator.formatowanie(sciezka)
		if zewnetrzne.Stoi(narzedzie) {
			sformatowane, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, session.Okno{},
				session.Zasady{}, session.Obszar{}, narzedzie, argumenty,
				filepath.Dir(sciezka), granicaCzasuAnalizy)
			// Formatowanie nieudane nie unieważnia analizy; pole treści
			// sformatowanej zostaje wtedy nieobecne.
			switch {
			case err != nil:
			case analizator.formatZPliku:
				if poprawiona, odczyt := os.ReadFile(sciezka); odczyt == nil {
					tresc := string(poprawiona)
					odpowiedz.Formatted = &tresc
				}
			case len(sformatowane.Wyjscie) > 0:
				tresc := string(sformatowane.Wyjscie)
				odpowiedz.Formatted = &tresc
			}
		}
	}
	return odpowiedz, nil
}

// plikDoAnalizy odkłada treść do pliku tymczasowego z właściwym
// rozszerzeniem i oddaje drogę jego sprzątnięcia.
func plikDoAnalizy(tresc, rozszerzenie string) (string, func(), error) {
	plik, err := os.CreateTemp("", "danaco-analiza-*"+rozszerzenie)
	if err != nil {
		return "", func() {}, bladWykonaniaTerminala(
			"nie można założyć pliku tymczasowego analizy: " + err.Error())
	}
	sciezka := plik.Name()
	sprzataj := func() { _ = os.Remove(sciezka) }
	if _, err := plik.WriteString(tresc); err != nil {
		plik.Close()
		sprzataj()
		return "", func() {}, bladWykonaniaTerminala(
			"nie można zapisać treści do analizy: " + err.Error())
	}
	if err := plik.Close(); err != nil {
		sprzataj()
		return "", func() {}, bladWykonaniaTerminala(
			"nie można domknąć pliku analizy: " + err.Error())
	}
	return sciezka, sprzataj, nil
}

// uwagiShellCheck czyta wynik ShellChecka w postaci `json1` i przekłada
// każdy komentarz na uwagę kontraktu.
func uwagiShellCheck(wyjscie string) []shared.TerminalLintFinding {
	var odczyt struct {
		Comments []struct {
			Line    int    `json:"line"`
			Column  int    `json:"column"`
			Level   string `json:"level"`
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"comments"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(wyjscie)), &odczyt); err != nil {
		return nil
	}
	uwagi := make([]shared.TerminalLintFinding, 0, len(odczyt.Comments))
	for _, wpis := range odczyt.Comments {
		uwaga := shared.TerminalLintFinding{
			Line:     wpis.Line,
			Severity: wagaUwagi(wpis.Level),
			Message:  wpis.Message,
		}
		if wpis.Column > 0 {
			kolumna := wpis.Column
			uwaga.Column = &kolumna
		}
		if wpis.Code > 0 {
			regula := "SC" + strconv.Itoa(wpis.Code)
			uwaga.Rule = &regula
		}
		uwagi = append(uwagi, uwaga)
	}
	return uwagi
}

// uwagiPowerShell czyta wynik PSScriptAnalyzera. Waga przychodzi liczbą
// (0 – informacja, 1 – ostrzeżenie, 2 – błąd, 3 – błąd składni parsera).
func uwagiPowerShell(wyjscie string) []shared.TerminalLintFinding {
	tresc := strings.TrimSpace(wyjscie)
	if tresc == "" {
		return nil
	}
	var odczyt []struct {
		Line     int    `json:"Line"`
		Column   int    `json:"Column"`
		Severity int    `json:"Severity"`
		Message  string `json:"Message"`
		RuleName string `json:"RuleName"`
	}
	if err := json.Unmarshal([]byte(tresc), &odczyt); err != nil {
		return nil
	}
	uwagi := make([]shared.TerminalLintFinding, 0, len(odczyt))
	for _, wpis := range odczyt {
		waga := shared.TerminalLintSeverity(shared.TerminalLintSeverityWarning)
		switch {
		case wpis.Severity <= 0:
			waga = shared.TerminalLintSeverityInfo
		case wpis.Severity >= 2:
			waga = shared.TerminalLintSeverityError
		}
		uwaga := shared.TerminalLintFinding{Line: wpis.Line, Severity: waga, Message: wpis.Message}
		if wpis.Column > 0 {
			kolumna := wpis.Column
			uwaga.Column = &kolumna
		}
		if wpis.RuleName != "" {
			regula := wpis.RuleName
			uwaga.Rule = &regula
		}
		uwagi = append(uwagi, uwaga)
	}
	return uwagi
}

// uwagiInterpretera czyta diagnostykę interpretera sprawdzającego samą
// składnię, dającą co najwyżej jedną uwagę, zawsze błędu.
func uwagiInterpretera(diagnostyka string) []shared.TerminalLintFinding {
	tresc := strings.TrimSpace(diagnostyka)
	if tresc == "" {
		return nil
	}
	wiersze := strings.Split(tresc, "\n")
	numer := 0
	tresciwy := wiersze[len(wiersze)-1]
	for _, wiersz := range wiersze {
		if odczytany, jest := numerWierszaBledu(wiersz); jest {
			numer = odczytany
			break
		}
	}
	if numer <= 0 {
		numer = 1
	}
	return []shared.TerminalLintFinding{{
		Line:     numer,
		Severity: shared.TerminalLintSeverityError,
		Message:  strings.TrimSpace(tresciwy),
	}}
}

// numerWierszaBledu wyławia numer wiersza z komunikatu interpretera. Rozpoznaje
// obie postaci: `File "x.py", line 12` Pythona i `x.js:12` Node'a.
func numerWierszaBledu(wiersz string) (int, bool) {
	if miejsce := strings.LastIndex(wiersz, ", line "); miejsce >= 0 {
		reszta := strings.TrimSpace(wiersz[miejsce+len(", line "):])
		reszta = strings.SplitN(reszta, ",", 2)[0]
		if numer, err := strconv.Atoi(strings.TrimSpace(reszta)); err == nil {
			return numer, true
		}
	}
	dwukropek := strings.LastIndex(wiersz, ":")
	if dwukropek <= 0 || dwukropek == len(wiersz)-1 {
		return 0, false
	}
	if numer, err := strconv.Atoi(strings.TrimSpace(wiersz[dwukropek+1:])); err == nil && numer > 0 {
		return numer, true
	}
	return 0, false
}

// wagaUwagi przekłada wagę ShellChecka na słownik kontraktu. Wartość
// nierozpoznana schodzi na ostrzeżenie: uwaga o nieznanej wadze ma być
// widoczna, a nie przemilczana ani podniesiona do błędu.
func wagaUwagi(poziom string) shared.TerminalLintSeverity {
	switch strings.ToLower(strings.TrimSpace(poziom)) {
	case "error":
		return shared.TerminalLintSeverityError
	case "info", "style":
		return shared.TerminalLintSeverityInfo
	default:
		return shared.TerminalLintSeverityWarning
	}
}

// bladWykonaniaTerminala znakuje niepowodzenie czynności rdzenia przy
// uruchamianiu albo obsłudze programu analizy.
func bladWykonaniaTerminala(powod string) error {
	return protocolBladTerminala(shared.ErrorCodeInternalError, powod)
}
