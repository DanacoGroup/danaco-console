// Plik odczytuje definicję eksperta z rdzenia komendą agent.list, zamiast
// sięgać po bazę wprost, i czyta ją leniwie, przy pierwszym żądaniu
// tools/list.
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

// OdczytajEksperta pyta rdzeń o eksperta o wskazanym kodzie; błąd znaczy nie
// wiadomo, nie ekspert nie ma narzędzi.
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

// wykazEkspertow wyjmuje z odpowiedzi rdzenia na żądanie agent.list wykaz
// ekspertów albo opis odmowy rdzenia.
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

// dopasuj odszukuje eksperta po kodzie; kod nieznany rdzeniowi jest faktem,
// nie pustką, i wraca błędem z liczbą ekspertów.
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
