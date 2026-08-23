package models

import "context"

// Kanal jest adapterem jednego kanału modelu. Wszystkie kanały — bez sieci,
// procesowy CLI, sieciowy API — realizują ten sam interfejs i nadają ten sam
// strumień fragmentów. Rdzeń nie zna żadnego adaptera z osobna,
// sięga po nie wyłącznie przez rejestr.
//
// Adapter kanału głównego (Claude Code CLI) realizuje ten interfejs w pakiecie
// internal/injection — tam mieszka budowa argv, pula kont i parser strumienia.
type Kanal interface {
	// Kod zwraca kod kanału z wiersza rejestru.
	Kod() string
	// Definicja zwraca wiersz rejestru, z którego kanał powstał.
	Definicja() Definicja
	// Wyslij wykonuje jedno zapytanie i nadaje strumień fragmentów do ujścia.
	// Pierwszym nadanym fragmentem jest prowenancja wywołania.
	// Zwrócony błąd dotyczy wyłącznie tego wywołania. Adapter błędu
	// nie zamienia sam na fragment — robi to Rejestr.Wyslij, dzięki czemu
	// strumień nie niesie dwóch fragmentów błędu o tej samej przyczynie.
	Wyslij(ctx context.Context, z Zapytanie, u Ujscie) error
}

// Zamykalny realizuje kanał, który trzyma zasób wymagający zwolnienia — proces,
// połączenie, pulę. Rejestr zamyka takie kanały przy odświeżeniu i wygaszeniu.
type Zamykalny interface {
	Zamknij() error
}

// Fabryka buduje adapter z wiersza rejestru kanałów. Klucz, pod którym fabryka
// stoi w rejestrze, jest wartością danych (kolumna rodzaj_kanalu albo parametr
// "adapter"), nie nazwą typu — dopisanie kanału to dopisanie wiersza, nie zmiana
// kodu.
type Fabryka func(d Definicja) (Kanal, error)

// Fabryki jest zbiorem fabryk adapterów, kluczowanym wartością danych.
type Fabryki map[string]Fabryka

// FabrykiWbudowane zwraca fabryki adapterów mieszkających w tym pakiecie:
// echo (bez sieci, do testów i pracy lokalnej) oraz generyczny api
// (parametryzowany adresem, modelem i odwołaniem do klucza — bez zaszytego
// dostawcy) oraz obrazy — kanał sieciowy oddający bajty obrazu zamiast tekstu.
// Adapter kanału głównego dokłada pakiet internal/injection przez
// Rejestr.UstawFabryke.
func FabrykiWbudowane() Fabryki {
	return Fabryki{
		AdapterEcho:   NowyKanalEcho,
		AdapterAPI:    NowyKanalAPI,
		AdapterObrazy: NowyKanalObrazow,
	}
}

// Klucze adapterów wbudowanych. Są wartościami danych zapisywanymi w wierszu
// rejestru, nie nazwami typów.
const (
	// AdapterEcho — kanał bez sieci i bez procesu; odsyła treść zapytania.
	AdapterEcho = "echo"
	// AdapterAPI — generyczny kanał sieciowy; dostawcę wyznaczają parametry.
	AdapterAPI = "api"
	// AdapterObrazy — kanał sieciowy oddający bajty obrazu (fragment `image`)
	// zamiast tekstu; kształt żądania i odpowiedzi jak w OpenAI Images. Stoi
	// obok AdapterAPI, a nie zamiast niego: rodzaj kanału (`api`) mówi, jak
	// kanał rozmawia, a adapter — co oddaje. Wiersz zakłada się tym samym
	// `channel.add`, dokładając do konfiguracji "adapter":"obrazy".
	AdapterObrazy = "obrazy"
	// AdapterCLI — kanał główny; adapter realizuje internal/injection.
	AdapterCLI = "cli"
)
