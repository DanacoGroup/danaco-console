package core

import (
	"context"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
)

// wykonanie niesie parametry wywołania modelu rozstrzygnięte dla okna: nakład
// rozumowania i model zapasowy. Pole puste znaczy „bez wskazania" — kanał nie
// dopisze wtedy przełącznika, zamiast podać dostawcy pusty argument.
type wykonanie struct {
	// Naklad trafia do --effort: low, medium, high, xhigh, max.
	Naklad string
	// ModelZapasowy trafia do --fallback-model i jest identyfikatorem modelu, nie kodem kanału.
	ModelZapasowy string
	// PulapKosztuUSD trafia do --max-budget-usd, zero znaczy brak pułapu, jak napis pusty wyżej.
	PulapKosztuUSD float64
}

// ParametryWykonania podaje parametry wywołania obowiązujące okno rozmowy
// w danej chwili trwania sesji.
type ParametryWykonania interface {
	Ustal(kontekst context.Context, okno session.Okno) wykonanie
}

// wykonanieZKonfiguracji rozstrzyga parametry po ośmiu poziomach zasięgu i tłumaczy
// kod kanału zapasowego na identyfikator modelu.
type wykonanieZKonfiguracji struct {
	rozstrzygacz *konfig.Rozstrzygacz
	kanaly       dane.RepozytoriumKanalow
}

// noweWykonanieZKonfiguracji składa port. Brak rozstrzygacza daje port pusty —
// rozmowa biegnie dalej, tylko bez parametrów z konfiguracji.
func noweWykonanieZKonfiguracji(r *konfig.Rozstrzygacz, kanaly dane.RepozytoriumKanalow) ParametryWykonania {
	if r == nil {
		return nil
	}
	return &wykonanieZKonfiguracji{rozstrzygacz: r, kanaly: kanaly}
}

// Ustal rozstrzyga parametry w kontekście okna.
//
// Kontekst obejmuje okno i kartę sesji, bo okno jest poziomem najwęższym, oraz
// kanał modelu jako oś rozstrzygania — ten sam nakład może być ustawiony
// inaczej dla różnych modeli.
func (p *wykonanieZKonfiguracji) Ustal(kontekst context.Context, okno session.Okno) wykonanie {
	if p == nil || p.rozstrzygacz == nil {
		return wykonanie{}
	}
	zasieg := konfig.Kontekst{
		Okno:           okno.Id,
		KartaSesji:     okno.IdSesji,
		Model:          okno.KanalModelu,
		KontoOperatora: dane.KontoOperatora(kontekst),
	}
	return wykonanie{
		Naklad:         p.wartosc(zasieg, konfig.KluczNakladRozumowania),
		ModelZapasowy:  p.modelKanalu(kontekst, p.wartosc(zasieg, konfig.KluczKanalZapasowy)),
		PulapKosztuUSD: p.kwota(zasieg, konfig.KluczPulapKosztu),
	}
}

// kwota odczytuje nastawę liczbową jako kwotę w dolarach.
//
// Wartość nieczytelna jako liczba albo ujemna daje zero, czyli „bez pułapu" —
// a nie odmowę wywołania.
func (p *wykonanieZKonfiguracji) kwota(zasieg konfig.Kontekst, klucz string) float64 {
	surowa := p.wartosc(zasieg, klucz)
	if surowa == "" {
		return 0
	}
	liczba, err := strconv.ParseFloat(surowa, 64)
	if err != nil || liczba <= 0 {
		return 0
	}
	return liczba
}

// wartosc zwraca rozstrzygniętą wartość klucza bez otoczki niosącej poziom
// zasięgu, z którego pochodzi.
func (p *wykonanieZKonfiguracji) wartosc(zasieg konfig.Kontekst, klucz string) string {
	return strings.TrimSpace(p.rozstrzygacz.Rozstrzygnij(zasieg, klucz).Wartosc)
}

// modelKanalu odwzorowuje kod kanału na identyfikator modelu u dostawcy: kanał
// nieznany albo bez identyfikatora modelu daje napis pusty, stan poprawny,
// nie powód przerwania tury.
func (p *wykonanieZKonfiguracji) modelKanalu(kontekst context.Context, kod string) string {
	if kod == "" || p.kanaly == nil {
		return ""
	}
	kanal, err := p.kanaly.PobierzPoKodzie(kontekst, kod)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(kanal.IdentyfikatorModelu)
}
