package models

import "context"

// Kanal jest adapterem jednego kanału modelu. Wszystkie kanały — bez sieci, procesowy CLI, sieciowy
// API — realizują ten sam interfejs i nadają ten sam strumień fragmentów, sięgając po siebie wyłącznie przez rejestr.
type Kanal interface {
	// Kod zwraca kod kanału z wiersza rejestru.
	Kod() string
	// Definicja zwraca wiersz rejestru, z którego kanał powstał.
	Definicja() Definicja
	// Wyslij wykonuje jedno zapytanie i nadaje strumień fragmentów do ujścia.
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

// Fabryki jest zbiorem fabryk adapterów kanałów, kluczowanym wartością danych zapisaną w wierszu rejestru.
type Fabryki map[string]Fabryka

// FabrykiWbudowane zwraca fabryki adapterów mieszkających w tym pakiecie: echo bez sieci, generyczny api
// parametryzowany adresem i modelem oraz obrazy oddające bajty obrazu zamiast tekstu.
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
	// AdapterEcho — kanał bez sieci i bez procesu zewnętrznego; odsyła treść zapytania z powrotem jako odpowiedź testową.
	AdapterEcho = "echo"
	// AdapterAPI — generyczny kanał sieciowy, którego dostawcę i zachowanie wyznaczają wyłącznie parametry wiersza rejestru.
	AdapterAPI = "api"
	// AdapterObrazy — kanał sieciowy oddający bajty obrazu zamiast tekstu; kształt żądania i odpowiedzi jak w OpenAI Images, wiersz zakłada się komendą channel.add z parametrem adapter obrazy.
	AdapterObrazy = "obrazy"
	// AdapterCLI — kanał główny rozmowy; adapter tego kanału realizuje pakiet internal/injection tej budowy.
	AdapterCLI = "cli"
)
