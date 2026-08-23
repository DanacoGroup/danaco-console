package injection

import "strconv"

// argumentyStale są szkieletem każdego wywołania kanału: tryb bezinterakcyjny,
// wejście i wyjście strumieniem JSON-lines oraz pełne wyjście zdarzeń. Bez
// --verbose program nie wypuszcza zdarzeń pośrednich, więc strumień
// przestałby być strumieniem.
//
// --include-hook-events stoi tu z tego samego powodu co --verbose: stan zmieniają
// zdarzenia wykonawcze, nie słowo modelu, a bez tego przełącznika zdarzenia
// zaczepów nie wchodzą na strumień i rdzeń nie ma czym odróżnić „zaczep
// skonfigurowany" od „zaczep zadziałał".
var argumentyStale = []string{
	"-p",
	"--input-format", "stream-json",
	"--output-format", "stream-json",
	"--verbose",
	"--include-hook-events",
}

// Argumenty buduje wiersz argumentów programu `claude`. Przełącznik pojawia
// się wyłącznie wtedy, gdy odpowiadające mu ustawienie jest wypełnione — brak
// ustawienia nie jest błędem, tylko brakiem przełącznika.
//
// Katalogi robocze są listą, więc --add-dir powtarza się tyle razy, ile jest
// katalogów okna. Tak samo --mcp-config.
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

// kwotaPulapu zamienia pułap kosztu na wartość przełącznika --max-budget-usd.
//
// Pułap niedodatni daje napis pusty, a napis pusty nie dokłada przełącznika
// (dopisz niżej) — brak nastawy jest brakiem ograniczenia, nie ograniczeniem
// zerowym. Rozróżnienie jest istotne: --max-budget-usd 0 kazałby programowi
// przerwać turę, zanim ta cokolwiek zrobi.
//
// Zapis idzie z dokładnością do dwóch miejsc — do centa. Notacja wykładnicza,
// którą %v dałoby dla kwot bardzo małych albo bardzo dużych, nie jest tu
// wartością liczbową dla programu, tylko napisem, którego może nie przyjąć.
func kwotaPulapu(pulap float64) string {
	if pulap <= 0 {
		return ""
	}
	return strconv.FormatFloat(pulap, 'f', 2, 64)
}

// dopisz dokłada parę przełącznik–wartość, gdy wartość jest niepusta.
func dopisz(argv []string, przelacznik, wartosc string) []string {
	if wartosc == "" {
		return argv
	}
	return append(argv, przelacznik, wartosc)
}
