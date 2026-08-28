// Odpowiedzialność pliku: terminologia modułu Translate —
// `translate.glossary.list` i `translate.term.extract`. Wyjmowanie
// kandydatów idzie miarą częstości, nie modelem.
package core

import (
	"context"
	"sort"
	"strings"
	"unicode"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// najmniejszaDlugoscKandydata odcina słowa jednoliterowe i dwuliterowe. Termin
// krótszy niż trzy znaki jest zwykle spójnikiem albo przyimkiem.
const najmniejszaDlugoscKandydata = 3

// domyslnaLiczbaKandydatow ogranicza wykaz, gdy Operator nie wskazał własnego
// pułapu — inaczej długi dokument oddałby setki pozycji, z których nikt nie
// przejrzy nawet połowy.
const domyslnaLiczbaKandydatow = 50

// WykazTerminow obsługuje `translate.glossary.list` i oddaje wykaz
// terminów zawężony wskazaniem panelu.
func (a *adapterTlumaczenia) WykazTerminow(ctx context.Context,
	z shared.TranslateGlossaryListRequest) (shared.TranslateGlossaryListResponse, error) {

	filtr := dane.FiltrTerminow{
		Jezyk:     napisZeWskaznika(z.Language),
		Dziedzina: napisZeWskaznika(z.Domain),
		Fraza:     napisZeWskaznika(z.Query),
		Limit:     liczbaCalkowitaZeWskaznika(z.Limit),
		Offset:    liczbaCalkowitaZeWskaznika(z.Offset),
	}
	if z.Status != nil {
		filtr.Stan = string(*z.Status)
	}
	terminy, razem, err := a.repozytorium.TerminyZawezone(ctx, filtr)
	if err != nil {
		return shared.TranslateGlossaryListResponse{}, bladTlumaczenia(err)
	}
	wykaz := make([]shared.GlossaryTerm, 0, len(terminy))
	for _, termin := range terminy {
		wykaz = append(wykaz, zlozTerminSlownika(termin))
	}
	return shared.TranslateGlossaryListResponse{Terms: wykaz, Total: razem}, nil
}

// WyjmijTerminy obsługuje `translate.term.extract`. Liczy częstość słów
// i zbitek dwuwyrazowych w tekście źródłowym okna, odrzuca słowa funkcyjne
// i oddaje kandydatów wraz z jednym zdaniem kontekstu.
func (a *adapterTlumaczenia) WyjmijTerminy(ctx context.Context,
	z shared.TranslateTermExtractRequest) (shared.TranslateTermExtractResponse, error) {

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateTermExtractResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	if okno.TekstZrodlowy == nil || strings.TrimSpace(*okno.TekstZrodlowy) == "" {
		return shared.TranslateTermExtractResponse{}, bladWskazaniaTlumaczenia(
			"okno " + okno.Kod + " nie ma tekstu źródłowego — nie ma z czego wyjąć terminów")
	}
	tekst := *okno.TekstZrodlowy

	najmniejszaCzestosc := liczbaCalkowitaZeWskaznika(z.MinFrequency)
	if najmniejszaCzestosc <= 0 {
		najmniejszaCzestosc = 2
	}
	pulap := liczbaCalkowitaZeWskaznika(z.MaxCandidates)
	if pulap <= 0 {
		pulap = domyslnaLiczbaKandydatow
	}

	// Terminy już w słowniku są znane; Operator z `excludeKnown` dostaje
	// wyłącznie nierozstrzygnięte.
	znane := map[string]struct{}{}
	terminy, err := a.repozytorium.Terminy(ctx)
	if err != nil {
		return shared.TranslateTermExtractResponse{}, bladTlumaczenia(err)
	}
	for _, termin := range terminy {
		znane[strings.ToLower(termin.Zrodlo)] = struct{}{}
	}
	pomijajZnane := z.ExcludeKnown != nil && *z.ExcludeKnown

	czestosci := policzKandydatow(tekst)
	zdania := podzielNaZdania(tekst)

	kandydaci := []shared.TermCandidate{}
	for fraza, ile := range czestosci {
		if ile < najmniejszaCzestosc {
			continue
		}
		_, jestZnany := znane[fraza]
		if jestZnany && pomijajZnane {
			continue
		}
		kandydaci = append(kandydaci, shared.TermCandidate{
			Source:    fraza,
			Frequency: ile,
			Context:   wskaznikNapisu(zdanieZFraza(zdania, fraza)),
			Known:     jestZnany,
		})
	}
	// Porządek: najczęstsze najpierw, przy równej częstości alfabetycznie,
	// ten sam przy każdym wywołaniu.
	sort.Slice(kandydaci, func(i, j int) bool {
		if kandydaci[i].Frequency != kandydaci[j].Frequency {
			return kandydaci[i].Frequency > kandydaci[j].Frequency
		}
		return kandydaci[i].Source < kandydaci[j].Source
	})
	if len(kandydaci) > pulap {
		kandydaci = kandydaci[:pulap]
	}
	return shared.TranslateTermExtractResponse{Candidates: kandydaci}, nil
}

// policzKandydatow liczy wystąpienia słów i zbitek dwuwyrazowych
// w tekście źródłowym, do wyboru kandydatów.
func policzKandydatow(tekst string) map[string]int {
	slowa := strings.FieldsFunc(strings.ToLower(tekst), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-'
	})
	czestosci := map[string]int{}
	for i, slowo := range slowa {
		if liczbaZnakow(slowo) >= najmniejszaDlugoscKandydata && !slowoFunkcyjne(slowo) {
			czestosci[slowo]++
		}
		if i+1 >= len(slowa) {
			continue
		}
		nastepne := slowa[i+1]
		if liczbaZnakow(slowo) < najmniejszaDlugoscKandydata ||
			liczbaZnakow(nastepne) < najmniejszaDlugoscKandydata {
			continue
		}
		if slowoFunkcyjne(slowo) || slowoFunkcyjne(nastepne) {
			continue
		}
		czestosci[slowo+" "+nastepne]++
	}
	return czestosci
}

// slowaFunkcyjne to zamknięty wykaz słów, które są częste w każdym tekście
// i nie są terminami w żadnym.
var slowaFunkcyjne = map[string]struct{}{
	"oraz": {}, "albo": {}, "lecz": {}, "jest": {}, "sie": {}, "się": {}, "nie": {},
	"tego": {}, "tym": {}, "przez": {}, "dla": {}, "jako": {}, "przy": {}, "pod": {},
	"nad": {}, "bez": {}, "gdy": {}, "aby": {},
	"the": {}, "and": {}, "for": {}, "with": {}, "that": {}, "this": {}, "from": {},
	"are": {}, "was": {}, "were": {}, "have": {}, "has": {}, "not": {},
}

// slowoFunkcyjne mówi, czy słowo należy do wykazu `slowaFunkcyjne`
// powyżej, zamkniętego dla polskiego i angielskiego.
func slowoFunkcyjne(slowo string) bool {
	_, jest := slowaFunkcyjne[slowo]
	return jest
}

// zdanieZFraza wskazuje pierwsze zdanie, w którym fraza wystąpiła — kontekst,
// bez którego kandydat jest samym napisem.
func zdanieZFraza(zdania []string, fraza string) string {
	for _, zdanie := range zdania {
		if strings.Contains(strings.ToLower(zdanie), fraza) {
			return zdanie
		}
	}
	return ""
}
