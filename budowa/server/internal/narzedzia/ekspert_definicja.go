// Odpowiedzialność pliku: odczyt definicji eksperta z rdzenia.
//
// ── DLACZEGO PYTAMY RDZEŃ, A NIE BAZĘ ───────────────────────────────────────
// Serwer narzędzi jest dla rdzenia zwykłym urządzeniem: ta sama koperta, to samo
// gniazdo, co okno interfejsu (`polaczenie.go`). Drugiego wejścia do
// danych eksperta tu nie ma i nie powstaje — sięgnięcie po sterownik
// bazy z procesu modelu byłoby obejściem rdzenia, a nie skrótem.
//
// ── DLACZEGO `agent.list`, SKORO SZUKAMY JEDNEGO ────────────────────────────
// Kontrakt nie ma komendy `agent.get`, a `AgentListRequest` nie filtruje po
// identyfikatorze — pola takiego po prostu nie niesie. Bierzemy więc wykaz
// i dopasowujemy po `Agent.Id`, i tak to jest tu nazwane zamiast schowane:
// gdyby `agent.get` kiedyś powstał, ta funkcja jest jedynym miejscem do zmiany.
// Komenda jest w `shared.KomendyNarzedzi`, więc droga istnieje bez żadnej
// zmiany kontraktu.
//
// ── ODCZYT JEST LENIWY ──────────────────────────────────────────────────────
// Gniazdo do rdzenia zestawia się przy pierwszym użyciu, z zamysłem: proces
// modelu uruchamia serwery MCP na starcie rozmowy, a rdzeń bywa wtedy jeszcze
// niegotowy (`polaczenie.go`). Definicji nie da się więc mieć w chwili startu
// procesu i nie udajemy, że się da — czyta się ją przy pierwszym `tools/list`.
package narzedzia

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync/atomic"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// DefinicjaEksperta niesie to i tylko to, czego dobór narzędzi potrzebuje:
// kod eksperta oraz kody, którymi wskazał swoje wyposażenie.
//
// Pola `Umiejetnosci` i `Konektory` to `Agent.SkillIds` i `Agent.ConnectorIds`
// kontraktu, przeniesione bez zmiany znaczenia. Reszta eksperta — warstwy,
// model, uprawnienia — do doboru narzędzi nie należy i nie jest tu kopiowana.
type DefinicjaEksperta struct {
	// Kod jest identyfikatorem eksperta z pola `Agent.Id`.
	Kod string
	// Umiejetnosci niesie `Agent.SkillIds`.
	Umiejetnosci []string
	// Konektory niesie `Agent.ConnectorIds`.
	Konektory []string
}

// Kody zwraca komplet kodów wskazanych przez eksperta, umiejętności przed
// konektorami, bez powtórzeń i bez kodów pustych.
//
// Jedna lista, bo dobór narzędzi nie rozróżnia pochodzenia kodu: rozpoznanie
// idzie po tym, czy kod nazywa narzędzie albo grupę (`ekspert_wykaz.go`), a nie
// po tym, w którym polu eksperta go zapisano.
func (d DefinicjaEksperta) Kody() []string {
	widziane := make(map[string]bool, len(d.Umiejetnosci)+len(d.Konektory))
	kody := make([]string, 0, len(d.Umiejetnosci)+len(d.Konektory))
	for _, zrodlo := range [][]string{d.Umiejetnosci, d.Konektory} {
		for _, kod := range zrodlo {
			if kod == "" || widziane[kod] {
				continue
			}
			widziane[kod] = true
			kody = append(kody, kod)
		}
	}
	return kody
}

// licznikOdczytow nadaje identyfikatory żądaniom odczytu definicji. Osobny od
// licznika rozdzielni, bo odczyt bywa robiony bez niej (test, narzędzie
// diagnostyczne) — identyfikator ma być jednoznaczny, a nie wspólny.
var licznikOdczytow atomic.Uint64

// OdczytajEksperta pyta rdzeń o eksperta o wskazanym kodzie.
//
// Błąd znaczy „nie wiadomo", a nie „ekspert nie ma narzędzi": rdzeń nie
// odpowiedział, odmówił, albo w jego wykazie tego kodu nie ma. Rozstrzygnięcie,
// co z niewiedzą zrobić, należy do składania wykazu (`ekspert_wykaz.go`) —
// tutaj jest wyłącznie odczyt i wyłącznie prawda o tym, co rdzeń powiedział.
func OdczytajEksperta(kontekst context.Context, rdzen Rdzen, kod string) (DefinicjaEksperta, error) {
	if rdzen == nil {
		return DefinicjaEksperta{}, fmt.Errorf("ekspert %q: rdzenia nie ma po drugiej stronie", kod)
	}
	if kod == "" {
		return DefinicjaEksperta{}, fmt.Errorf("odczyt eksperta bez kodu — nie ma o kogo zapytać")
	}
	identyfikator := "narzedzia-ekspert-" + strconv.FormatUint(licznikOdczytow.Add(1), 10)
	zadanie, err := protocol.NowaKoperta(shared.CommandAgentList, identyfikator, "", shared.AgentListRequest{})
	if err != nil {
		return DefinicjaEksperta{}, fmt.Errorf("ekspert %q: żądanie %s nie daje się złożyć: %w",
			kod, shared.CommandAgentList, err)
	}
	odpowiedz, err := rdzen.Wykonaj(kontekst, zadanie)
	if err != nil {
		return DefinicjaEksperta{}, fmt.Errorf("ekspert %q: rdzeń nie odpowiedział na %s: %w",
			kod, shared.CommandAgentList, err)
	}
	wykaz, err := wykazEkspertow(kod, odpowiedz)
	if err != nil {
		return DefinicjaEksperta{}, err
	}
	return dopasuj(kod, wykaz)
}

// wykazEkspertow wyjmuje z odpowiedzi rdzenia wykaz ekspertów albo opis odmowy.
func wykazEkspertow(kod string, odpowiedz protocol.Koperta) ([]shared.Agent, error) {
	if odpowiedz.Error != nil {
		return nil, fmt.Errorf("ekspert %q: rdzeń odmówił na %s: %s",
			kod, shared.CommandAgentList, protocol.Opis(*odpowiedz.Error))
	}
	if odpowiedz.Status != nil && *odpowiedz.Status != shared.EnvelopeStatusOk {
		return nil, fmt.Errorf("ekspert %q: rdzeń odpowiedział na %s stanem %s bez opisu błędu",
			kod, shared.CommandAgentList, *odpowiedz.Status)
	}
	var tresc shared.AgentListResponse
	if err := json.Unmarshal(odpowiedz.Payload, &tresc); err != nil {
		return nil, fmt.Errorf("ekspert %q: treść wyniku %s nieczytelna: %w",
			kod, shared.CommandAgentList, err)
	}
	return tresc.Agents, nil
}

// dopasuj odszukuje eksperta po kodzie.
//
// Kod nieznany rdzeniowi jest faktem, nie pustką: ekspert bywa kasowany
// niezależnie od okien, w których pracował (`session/okno.go`), więc okno może
// nieść kod, którego biblioteka już nie ma. Taki stan wraca błędem mówiącym,
// ilu ekspertów rdzeń zna — inaczej „nie znaleziono" byłoby nieodróżnialne od
// „rdzeń oddał wykaz pusty".
func dopasuj(kod string, wykaz []shared.Agent) (DefinicjaEksperta, error) {
	for _, ekspert := range wykaz {
		if ekspert.Id != kod {
			continue
		}
		return DefinicjaEksperta{
			Kod:          ekspert.Id,
			Umiejetnosci: ekspert.SkillIds,
			Konektory:    ekspert.ConnectorIds,
		}, nil
	}
	return DefinicjaEksperta{}, fmt.Errorf(
		"ekspert %q: rdzeń go nie zna — w wykazie %s stoi %d ekspertów, żaden o tym kodzie",
		kod, shared.CommandAgentList, len(wykaz))
}
