// Pakiet stdio prowadzi protokół MCP po strumieniach procesu: jedno wywołanie
// JSON-RPC 2.0 w wierszu wejścia, jedna odpowiedź w wierszu wyjścia. Pakiet nie
// zna komend kontraktu ani nazw narzędzi, pyta o nie katalog.
package stdio

import (
	"context"
	"encoding/json"

	"danacoconsole/server/internal/narzedzia"
)

const (
	// wersjaJsonRpc jest wersją protokołu JSON-RPC przenoszącego wywołania między
	// procesem modelu a serwerem.
	wersjaJsonRpc = "2.0"
	// wersjaMcp jest wersją protokołu MCP, którą serwer ogłasza procesowi modelu
	// w treści odpowiedzi powitalnej.
	wersjaMcp = "2025-06-18"
	// wersjaSerwera odpowiada wersji produktu. Pakiet nie sięga po tamtą stałą,
	// bo import rdzenia wciągnąłby do binarium serwera narzędzi całą trwałość
	// wraz z bazą.
	wersjaSerwera = "1.0"
	// nazwaSerwera jest nazwą serwera narzędzi widzianą przez proces modelu
	// w odpowiedzi powitalnej protokołu.
	nazwaSerwera = narzedzia.KluczWpisu
)

// Kody błędów JSON-RPC 2.0 używane przez serwer. Odmowa narzędzia nie jest
// błędem protokołu — wraca wynikiem oznaczonym jako błędny, żeby model mógł ją
// przeczytać i poprawić.
const (
	kodZlyKomunikat = -32600
	kodBrakMetody   = -32601
	kodZleParametry = -32602
)

// Katalog jest jedynym, czego warstwa protokołu potrzebuje od rozdzielni.
// Interfejs stoi po stronie odbiorcy, więc pakiet nie narzuca rozdzielni
// żadnego kształtu poza tymi dwiema czynnościami.
type Katalog interface {
	// Narzedzia zwraca wykaz narzędzi podawany modelowi; kontekst pozwala
	// złożeniu wykazu pytać rdzeń.
	Narzedzia(kontekst context.Context) []narzedzia.Narzedzie
	// Wywolaj wykonuje jedno narzędzie; zwrócony błąd jest treścią dla modelu.
	Wywolaj(kontekst context.Context, nazwa string, argumenty map[string]any) (string, error)
}

// zadanie jest wywołaniem JSON-RPC. Identyfikator zostaje surowy, bo protokół
// dopuszcza napis i liczbę, a odpowiedź ma go powtórzyć bez zmiany kształtu.
// Zawiadomienie (notification) nie ma identyfikatora i nie dostaje odpowiedzi.
type zadanie struct {
	JsonRpc string          `json:"jsonrpc"`
	Id      json.RawMessage `json:"id,omitempty"`
	Metoda  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// odpowiedz jest wynikiem wywołania JSON-RPC zwracanym procesowi modelu
// w jednym wierszu wyjścia serwera.
type odpowiedz struct {
	JsonRpc string          `json:"jsonrpc"`
	Id      json.RawMessage `json:"id"`
	Wynik   any             `json:"result,omitempty"`
	Blad    *bladProtokolu  `json:"error,omitempty"`
}

// bladProtokolu opisuje usterkę samego wywołania: nieczytelny komunikat, brak
// metody, parametry nie do odczytania.
type bladProtokolu struct {
	Kod       int    `json:"code"`
	Komunikat string `json:"message"`
}

// wynikiem pakuje wynik wykonanej metody w odpowiedź JSON-RPC gotową do
// zapisania w wierszu wyjścia serwera.
func wynikiem(id json.RawMessage, wynik any) odpowiedz {
	return odpowiedz{JsonRpc: wersjaJsonRpc, Id: id, Wynik: wynik}
}

// bledem pakuje usterkę wywołania w odpowiedź JSON-RPC gotową do zapisania
// procesowi modelu jednym wierszem.
func bledem(id json.RawMessage, kod int, komunikat string) odpowiedz {
	return odpowiedz{JsonRpc: wersjaJsonRpc, Id: id, Blad: &bladProtokolu{Kod: kod, Komunikat: komunikat}}
}
