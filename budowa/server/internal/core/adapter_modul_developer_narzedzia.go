// Odpowiedzialność pliku: programy warsztatu deweloperskiego używane w module
// Developer, droga ich wywołania i komenda developer.toolchain.check.
package core

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"

	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// Programy warsztatu deweloperskiego. Nazwa czytelna wchodzi do treści odmowy,
// a podpowiedź instalacyjna mówi Operatorowi serwera, czym brak uzupełnić.
var (
	narzedzieGofmt = zewnetrzne.Narzedzie{
		Nazwa: "gofmt", Program: "gofmt", Pakiet: "golang"}
	narzedzieGoimports = zewnetrzne.Narzedzie{
		Nazwa: "goimports", Program: "goimports",
		Pakiet: "go install golang.org/x/tools/cmd/goimports@latest"}
	narzedzieGopls = zewnetrzne.Narzedzie{
		Nazwa: "gopls", Program: "gopls",
		Pakiet: "go install golang.org/x/tools/gopls@latest"}
	narzedzieGolangciLint = zewnetrzne.Narzedzie{
		Nazwa: "golangci-lint", Program: "golangci-lint",
		Pakiet: "go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"}
	narzedzieStaticcheck = zewnetrzne.Narzedzie{
		Nazwa: "staticcheck", Program: "staticcheck",
		Pakiet: "go install honnef.co/go/tools/cmd/staticcheck@latest"}
	narzedzieDelve = zewnetrzne.Narzedzie{
		Nazwa: "Delve", Program: "dlv",
		Pakiet: "go install github.com/go-delve/delve/cmd/dlv@latest"}
	narzedziePrettier = zewnetrzne.Narzedzie{
		Nazwa: "Prettier", Program: "prettier", Pakiet: "npm i -g prettier"}
	narzedzieEslint = zewnetrzne.Narzedzie{
		Nazwa: "ESLint", Program: "eslint", Pakiet: "npm i -g eslint"}
	narzedzieSilnikaKontenerow = zewnetrzne.Narzedzie{
		Nazwa: "Docker", Program: "docker", Pakiet: "docker.io albo podman"}
	narzedzieRuff = zewnetrzne.Narzedzie{
		Nazwa: "Ruff", Program: "ruff", Pakiet: "pip install ruff"}
	narzedzieSemgrep = zewnetrzne.Narzedzie{
		Nazwa: "Semgrep", Program: "semgrep", Pakiet: "pip install semgrep"}
	narzedzieAstGrep = zewnetrzne.Narzedzie{
		Nazwa: "ast-grep", Program: "ast-grep", Pakiet: "npm i -g @ast-grep/cli"}
	narzedzieJscpd = zewnetrzne.Narzedzie{
		Nazwa: "jscpd", Program: "jscpd", Pakiet: "npm i -g jscpd"}
	narzedzieDupl = zewnetrzne.Narzedzie{
		Nazwa: "dupl", Program: "dupl",
		Pakiet: "go install github.com/mibk/dupl@latest"}
	narzedzieTypos = zewnetrzne.Narzedzie{
		Nazwa: "typos", Program: "typos", Pakiet: "cargo install typos-cli"}
	narzedzieStylelint = zewnetrzne.Narzedzie{
		Nazwa: "Stylelint", Program: "stylelint", Pakiet: "npm i -g stylelint"}
	// narzedzieSerweraTypeScript jest zadeklarowane, lecz rdzeń nie ma czym go
	// zapytać pojedynczym wywołaniem — program mówi wyłącznie sesją LSP, nie
	// wierszem poleceń; powód stoi przy serwerJezykaPliku.
	narzedzieSerweraTypeScript = zewnetrzne.Narzedzie{
		Nazwa: "serwer języka TypeScript", Program: "typescript-language-server",
		Pakiet: "npm i -g typescript-language-server typescript"}
)

// narzedziaWarsztatuDevelopera oddaje komplet programów warsztatu w kolejności
// ustalonej. Kolejność jest ustalona, żeby dwa kolejne odczyty dawały ten sam
// wykaz — Dev Tools ma się różnić wtedy, gdy zmienił się stan serwera.
func narzedziaWarsztatuDevelopera() []zewnetrzne.Narzedzie {
	return []zewnetrzne.Narzedzie{
		narzedzieAstGrep,
		narzedzieDelve,
		narzedzieDupl,
		narzedzieEslint,
		narzedzieGofmt,
		narzedzieGoimports,
		narzedzieGolangciLint,
		narzedzieGopls,
		narzedzieJscpd,
		narzedziePrettier,
		narzedzieRuff,
		narzedzieSemgrep,
		narzedzieSerweraTypeScript,
		narzedzieSilnikaKontenerow,
		narzedzieStaticcheck,
		narzedzieStylelint,
		narzedzieTypos,
	}
}

// narzedzieWarsztatuPoNazwie odnajduje deklarację po nazwie programu albo po
// nazwie czytelnej. Klient pyta nazwą programu (`gopls`), lecz Operator czyta
// nazwę czytelną (`Delve`) — obie mają trafiać w ten sam wiersz.
func narzedzieWarsztatuPoNazwie(wskazanie string) (zewnetrzne.Narzedzie, bool) {
	szukane := strings.ToLower(strings.TrimSpace(wskazanie))
	if szukane == "" {
		return zewnetrzne.Narzedzie{}, false
	}
	for _, narzedzie := range narzedziaWarsztatuDevelopera() {
		if strings.ToLower(narzedzie.Program) == szukane ||
			strings.ToLower(narzedzie.Nazwa) == szukane {
			return narzedzie, true
		}
	}
	// Program spoza wykazu nie jest odmową — odpowiedź nie ma jest odpowiedzią, nie usterką.
	return zewnetrzne.Narzedzie{Nazwa: wskazanie, Program: wskazanie}, true
}

// SprawdzWarsztat obsługuje developer.toolchain.check. Pusty wykaz w żądaniu
// znaczy powiedz o wszystkim, co znasz — Dev Tools otwiera zakładkę bez
// wiedzy o tym, czego szukać, i to rdzeń ma ten wykaz.
func (a *adapterDevelopera) SprawdzWarsztat(_ context.Context,
	z shared.DeveloperToolchainCheckRequest) (shared.DeveloperToolchainCheckResponse, error) {

	pytane := make([]zewnetrzne.Narzedzie, 0, len(z.Programs))
	for _, program := range z.Programs {
		if narzedzie, jest := narzedzieWarsztatuPoNazwie(program); jest {
			pytane = append(pytane, narzedzie)
		}
	}
	if len(pytane) == 0 {
		pytane = narzedziaWarsztatuDevelopera()
	}

	wykaz := make([]shared.ToolchainProgram, 0, len(pytane))
	for _, narzedzie := range pytane {
		sciezka, stoi := zewnetrzne.Odnajdz(narzedzie)
		pozycja := shared.ToolchainProgram{Program: narzedzie.Program, Present: stoi}
		if stoi {
			pozycja.Path = wskaznikTekstu(sciezka)
			if wersja := wersjaProgramuWarsztatu(sciezka); wersja != "" {
				pozycja.Version = wskaznikTekstu(wersja)
			}
		}
		wykaz = append(wykaz, pozycja)
	}
	return shared.DeveloperToolchainCheckResponse{Programs: wykaz}, nil
}

// wersjaProgramuWarsztatu pyta program o jego wersję. Jedyne miejsce modułu,
// w którym wywołanie idzie poza port uruchamiacza, bo pytanie o wersję nie
// należy do żadnego okna.
func wersjaProgramuWarsztatu(sciezka string) string {
	ctx, przerwij := context.WithTimeout(context.Background(), 3*time.Second)
	defer przerwij()

	for _, parametr := range []string{"--version", "version", "-version"} {
		wyjscie, err := exec.CommandContext(ctx, sciezka, parametr).CombinedOutput()
		if err != nil {
			continue
		}
		if wersja := pierwszyWierszWersji(string(wyjscie)); wersja != "" {
			return wersja
		}
	}
	return ""
}

// pierwszyWierszWersji bierze pierwszy niepusty wiersz odpowiedzi programu, pomijając wiersze puste na wejściu.
func pierwszyWierszWersji(wyjscie string) string {
	for _, wiersz := range strings.Split(wyjscie, "\n") {
		if tresc := strings.TrimSpace(wiersz); tresc != "" {
			return tresc
		}
	}
	return ""
}

// wolajNarzedzieWarsztatu uruchamia program warsztatu w obszarze okna, portem
// session.Uruchamiacz, tą samą bramą izolacji co budowanie i git.
func (a *adapterDevelopera) wolajNarzedzieWarsztatu(ctx context.Context, okno session.Okno,
	narzedzie zewnetrzne.Narzedzie, argumenty []string, katalog string,
	limit time.Duration) (zewnetrzne.Wynik, error) {

	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfigKontekstOkna(okno))
	}
	return zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, a.obszarDevelopera(okno),
		narzedzie, argumenty, katalog, limit)
}

// brakNarzedziaWarsztatu rozpoznaje odmowę nie ma czym. Rodziny funkcji
// rozróżniają ją od awarii osobnym polem kontraktu, bo brak programu nie
// jest odmową całej komendy — odpowiedź wraca pusta i jawnie to mówi.
func brakNarzedziaWarsztatu(err error) bool {
	var brak *zewnetrzne.BrakNarzedzia
	return err != nil && errors.As(err, &brak)
}
