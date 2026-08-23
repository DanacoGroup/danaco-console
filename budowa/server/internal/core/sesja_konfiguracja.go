// Odpowiedzialność pliku: model obszarów jednolitej konfiguracji sesji widziany
// od strony trwałości — wykaz obszarów, klucz zapisu obszaru oraz rozbiór
// i złożenie konfiguracji z obszarów.
//
// Obszar jest jednostką zapisu. Komendy `config.session.*` biorą wykaz obszarów,
// nie wykaz pól, więc jeden obszar to dokładnie jeden wiersz tabeli `ustawienie`
// spod adresu złożonego (poziom zasięgu i oś). Dzięki temu obszar zapisany na
// poziomie węższym przykrywa obszar poziomu szerszego w całości, a obszar spoza
// wykazu `areas` pozostaje nietknięty.
//
// Wykaz obszarów nie jest listą stałych w kodzie — powstaje z nazw pól
// kontraktu (shared.SessionConfig). Dopisanie obszaru do kontraktu wystarcza:
// rdzeń pozna go bez zmiany ani jednej gałęzi. Kolejność obszarów jest
// kolejnością pól kontraktu, więc to samo wejście daje ten sam wynik przy
// każdym wywołaniu.
//
// Rozbiór i złożenie idą przez kodowanie JSON kontraktu, a nie przez ręczne
// przypisania pól. Drugiego opisu obszarów w rdzeniu nie ma.
package core

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// przedrostekObszaru poprzedza klucz każdego obszaru konfiguracji sesji
// w tabeli ustawień. Oddziela obszary od kluczy prostych rodziny `config.*`,
// więc jeden odczyt poziomu wystarcza, by rozpoznać jedne i drugie.
const przedrostekObszaru = "sesja.konfiguracja."

// obszaryKontraktu wylicza obszary w kolejności pól shared.SessionConfig.
var obszaryKontraktu = odczytajObszaryKontraktu()

// znaneObszary rozpoznaje obszar podany w żądaniu.
var znaneObszary = zbudujZnaneObszary()

// odczytajObszaryKontraktu bierze nazwy obszarów z etykiet JSON kontraktu.
// Nazwa obszaru jest nazwą pola kontraktu — rdzeń nie trzyma jej drugi raz.
func odczytajObszaryKontraktu() []shared.SessionConfigArea {
	typ := reflect.TypeOf(shared.SessionConfig{})
	obszary := make([]shared.SessionConfigArea, 0, typ.NumField())
	for indeks := 0; indeks < typ.NumField(); indeks++ {
		nazwa, _, _ := strings.Cut(typ.Field(indeks).Tag.Get("json"), ",")
		if nazwa == "" || nazwa == "-" {
			continue
		}
		obszary = append(obszary, shared.SessionConfigArea(nazwa))
	}
	return obszary
}

// zbudujZnaneObszary składa zbiór rozpoznawanych nazw obszarów.
func zbudujZnaneObszary() map[shared.SessionConfigArea]struct{} {
	zbior := make(map[shared.SessionConfigArea]struct{}, len(obszaryKontraktu))
	for _, obszar := range obszaryKontraktu {
		zbior[obszar] = struct{}{}
	}
	return zbior
}

// kluczObszaru buduje klucz zapisu obszaru w tabeli ustawień.
func kluczObszaru(obszar shared.SessionConfigArea) string {
	return przedrostekObszaru + string(obszar)
}

// obszarKlucza rozpoznaje obszar po kluczu ustawienia. Klucz spoza rodziny
// konfiguracji sesji oraz klucz obszaru nieznanego kontraktowi są pomijane —
// odczyt poziomu nie ma prawa wywrócić się na wierszu, którego nie rozumie.
func obszarKlucza(klucz string) (shared.SessionConfigArea, bool) {
	if !strings.HasPrefix(klucz, przedrostekObszaru) {
		return "", false
	}
	obszar := shared.SessionConfigArea(strings.TrimPrefix(klucz, przedrostekObszaru))
	_, znany := znaneObszary[obszar]
	return obszar, znany
}

// obszaryZadania zawęża wykaz obszarów do żądanych, zachowując kolejność
// kontraktu. Wykaz pusty znaczy komplet obszarów; obszar nieznany
// kontraktowi jest odmową merytoryczną jednego wywołania, nie awarią.
func obszaryZadania(zadane []shared.SessionConfigArea) ([]shared.SessionConfigArea, error) {
	if len(zadane) == 0 {
		return obszaryKontraktu, nil
	}
	wybrane := make(map[shared.SessionConfigArea]struct{}, len(zadane))
	for _, obszar := range zadane {
		if _, znany := znaneObszary[obszar]; !znany {
			return nil, bladNieznanegoObszaru(obszar)
		}
		wybrane[obszar] = struct{}{}
	}
	wynik := make([]shared.SessionConfigArea, 0, len(wybrane))
	for _, obszar := range obszaryKontraktu {
		if _, jest := wybrane[obszar]; jest {
			wynik = append(wynik, obszar)
		}
	}
	return wynik, nil
}

// rozbierzKonfiguracje rozkłada konfigurację na treści obszarów. Obszar
// niewypełniony nie trafia do wyniku — brak obszaru znaczy wartość
// odziedziczoną, nie wartość pustą.
func rozbierzKonfiguracje(k shared.SessionConfig) (map[shared.SessionConfigArea]json.RawMessage, error) {
	bajty, err := json.Marshal(k)
	if err != nil {
		return nil, err
	}
	var surowe map[string]json.RawMessage
	if err := json.Unmarshal(bajty, &surowe); err != nil {
		return nil, err
	}
	obszary := make(map[shared.SessionConfigArea]json.RawMessage, len(surowe))
	for nazwa, tresc := range surowe {
		if obszar, znany := obszarKlucza(przedrostekObszaru + nazwa); znany {
			obszary[obszar] = tresc
		}
	}
	return obszary, nil
}

// zlozKonfiguracje składa konfigurację z treści obszarów. Treść nieczytelna
// przerywa złożenie, bo cicho pominięty obszar byłby konfiguracją inną niż
// zapisana — a to wprowadzałoby Operatora w błąd.
func zlozKonfiguracje(obszary map[shared.SessionConfigArea]json.RawMessage) (shared.SessionConfig, error) {
	surowe := make(map[string]json.RawMessage, len(obszary))
	for obszar, tresc := range obszary {
		surowe[string(obszar)] = tresc
	}
	bajty, err := json.Marshal(surowe)
	if err != nil {
		return shared.SessionConfig{}, err
	}
	var konfiguracja shared.SessionConfig
	if err := json.Unmarshal(bajty, &konfiguracja); err != nil {
		return shared.SessionConfig{}, err
	}
	return konfiguracja, nil
}

// chwilaZapisu przekłada znacznik czasu wiersza na milisekundy epoki. Znacznik
// nieczytelny znaczy brak czasu, nie błąd odczytu.
func chwilaZapisu(znacznik string) (int64, bool) {
	chwila, err := time.Parse(time.RFC3339Nano, znacznik)
	if err != nil {
		return 0, false
	}
	return chwila.UnixMilli(), true
}

// wskaznikTekstu przenosi napis do pola opcjonalnego kontraktu. Napis pusty
// zostaje pominięty — kontrakt opisuje go brakiem wartości, nie pustką.
func wskaznikTekstu(napis string) *string {
	if napis == "" {
		return nil
	}
	return &napis
}

// bladNieznanegoObszaru odmawia obsługi obszaru spoza kontraktu.
func bladNieznanegoObszaru(obszar shared.SessionConfigArea) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"konfiguracja sesji: obszar "+string(obszar)+" nie należy do kontraktu"))
}
