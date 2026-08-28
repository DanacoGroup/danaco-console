package injection

import "strconv"

// argumentyStale są szkieletem każdego wywołania kanału: tryb bezinterakcyjny,
// wejście i wyjście strumieniem JSON-lines oraz pełne wyjście zdarzeń.
var argumentyStale = []string{
	"-p",
	"--input-format", "stream-json",
	"--output-format", "stream-json",
	"--verbose",
	"--include-hook-events",
}

// Argumenty buduje wiersz argumentów programu `claude`; przełącznik pojawia
// się wyłącznie wtedy, gdy odpowiadające mu ustawienie jest wypełnione.
func Argumenty(u Ustawienia, n Nakladka) []string {
	argv := make([]string, 0, len(argumentyStale)+2*len(u.Katalogi)+16)
	argv = append(argv, argumentyStale...)
	argv = dopisz(argv, "--model", u.Model)
	argv = dopisz(argv, "--fallback-model", u.ModelZapasowy)
	argv = dopisz(argv, "--effort", u.Naklad)
	argv = dopisz(argv, "--permission-mode", string(u.TrybUprawnien))
	for _, katalog := range u.Katalogi {
		argv = dopisz(argv, "--add-dir", katalog)
	}
	argv = dopisz(argv, "--settings", u.PlikUstawien)
	for _, konfiguracja := range u.KonfiguracjaMCP {
		argv = dopisz(argv, "--mcp-config", konfiguracja)
	}
	argv = dopisz(argv, "--resume", u.Wznowienie)
	argv = dopisz(argv, "--max-budget-usd", kwotaPulapu(u.PulapKosztuUSD))
	if przelacznik, wartosc, jest := n.Przelacznik(); jest {
		argv = append(argv, przelacznik, wartosc)
	}
	return argv
}

// kwotaPulapu zamienia pułap kosztu na wartość przełącznika --max-budget-usd,
// z dokładnością do dwóch miejsc po przecinku.
func kwotaPulapu(pulap float64) string {
	if pulap <= 0 {
		return ""
	}
	return strconv.FormatFloat(pulap, 'f', 2, 64)
}

// dopisz dokłada parę przełącznik–wartość do wiersza argumentów, gdy
// przekazana wartość jest niepusta.
func dopisz(argv []string, przelacznik, wartosc string) []string {
	if wartosc == "" {
		return argv
	}
	return append(argv, przelacznik, wartosc)
}
