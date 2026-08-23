// Odpowiedzialność pliku: segmentacja modułu Translate —
// `translate.segmentation.rules.list`, `translate.segmentation.rules.set`,
// `translate.segment.merge` i `translate.segment.split`.
//
// Scalanie i podział zmieniają trwały podział okna (tabela
// `segment_okna_tlumaczenia`, migracja 161), a nie sam wynik odpowiedzi.
// Pierwsze wywołanie utrwala bieżący podział w całości, dopiero potem zmienia
// w nim jedną rzecz — inaczej numer segmentu z żądania wskazywałby po chwili na
// inny segment, niż widział Operator.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekZestawuSegmentacji znakuje identyfikator zestawu reguł.
const przedrostekZestawuSegmentacji = "seg-"

// WykazRegulSegmentacji obsługuje `translate.segmentation.rules.list`.
func (a *adapterTlumaczenia) WykazRegulSegmentacji(ctx context.Context,
	z shared.TranslateSegmentationRulesListRequest) (shared.TranslateSegmentationRulesListResponse, error) {

	zestawy, err := a.repozytorium.ZestawyRegulSegmentacji(ctx, napisZeWskaznika(z.Language))
	if err != nil {
		return shared.TranslateSegmentationRulesListResponse{}, bladTlumaczenia(err)
	}
	wykaz := make([]shared.SegmentationRuleset, 0, len(zestawy))
	for _, zestaw := range zestawy {
		wykaz = append(wykaz, zlozZestawSegmentacji(zestaw))
	}
	return shared.TranslateSegmentationRulesListResponse{Rulesets: wykaz}, nil
}

// UstawRegulySegmentacji obsługuje `translate.segmentation.rules.set`.
func (a *adapterTlumaczenia) UstawRegulySegmentacji(ctx context.Context,
	z shared.TranslateSegmentationRulesSetRequest) (shared.TranslateSegmentationRulesSetResponse, error) {

	if strings.TrimSpace(z.Name) == "" {
		return shared.TranslateSegmentationRulesSetResponse{}, bladWskazaniaTlumaczenia(
			"zestaw reguł segmentacji bez nazwy — Operator nie odróżni go potem od innego")
	}
	kod := napisZeWskaznika(z.RulesetId)
	if strings.TrimSpace(kod) == "" {
		kod = nowyIdentyfikator(przedrostekZestawuSegmentacji)
	}
	zestaw := dane.ZestawRegulSegmentacji{
		Kod:   kod,
		Nazwa: z.Name,
		Jezyk: z.Language,
		Srx:   z.Srx,
	}
	for numer, regula := range z.Rules {
		kolejnosc := int64(numer)
		if regula.Order > 0 {
			kolejnosc = int64(regula.Order)
		}
		zestaw.Reguly = append(zestaw.Reguly, dane.RegulaSegmentacji{
			Kolejnosc: kolejnosc,
			Przed:     regula.BeforeBreak,
			Po:        regula.AfterBreak,
			Lamie:     regula.Breaks,
		})
	}
	zapisany, err := a.repozytorium.ZapiszZestawRegulSegmentacji(ctx, zestaw)
	if err != nil {
		return shared.TranslateSegmentationRulesSetResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateSegmentationRulesSetResponse{
		Ruleset: zlozZestawSegmentacji(zapisany)}, nil
}

// zlozZestawSegmentacji przekłada wiersze zestawu na byt kontraktu.
func zlozZestawSegmentacji(zestaw dane.ZestawRegulSegmentacji) shared.SegmentationRuleset {
	reguly := make([]shared.SegmentationRule, 0, len(zestaw.Reguly))
	for _, regula := range zestaw.Reguly {
		reguly = append(reguly, shared.SegmentationRule{
			Order:       int(regula.Kolejnosc),
			BeforeBreak: regula.Przed,
			AfterBreak:  regula.Po,
			Breaks:      regula.Lamie,
		})
	}
	return shared.SegmentationRuleset{
		Id:        zestaw.Kod,
		Name:      zestaw.Nazwa,
		Language:  zestaw.Jezyk,
		Rules:     reguly,
		UpdatedAt: zestaw.Zaktualizowano,
	}
}

// ScalSegmenty obsługuje `translate.segment.merge`. Scala wskazane segmenty
// w jeden, w kolejności ich numerów, i utrwala nowy podział okna.
func (a *adapterTlumaczenia) ScalSegmenty(ctx context.Context,
	z shared.TranslateSegmentMergeRequest) (shared.TranslateSegmentMergeResponse, error) {

	okno, segmenty, err := a.oknoISegmenty(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateSegmentMergeResponse{}, err
	}
	if len(z.SegmentIndexes) < 2 {
		return shared.TranslateSegmentMergeResponse{}, bladWskazaniaTlumaczenia(
			"scalenie wymaga co najmniej dwóch segmentów — jeden segment jest już scalony")
	}
	numery, err := numerySegmentow(z.SegmentIndexes, len(segmenty))
	if err != nil {
		return shared.TranslateSegmentMergeResponse{}, err
	}

	doScalenia := map[int]struct{}{}
	for _, numer := range numery {
		doScalenia[numer] = struct{}{}
	}
	nowe := []string{}
	scalony := []string{}
	for numer, segment := range segmenty {
		if _, wchodzi := doScalenia[numer]; wchodzi {
			scalony = append(scalony, segment)
			// Scalony segment wchodzi w miejsce pierwszego ze scalanych — tam,
			// gdzie Operator go widzi, a nie na końcu wykazu.
			if numer == numery[0] {
				nowe = append(nowe, "")
			}
			continue
		}
		nowe = append(nowe, segment)
	}
	for i, segment := range nowe {
		if segment == "" {
			nowe[i] = strings.Join(scalony, " ")
			break
		}
	}

	if err := a.repozytorium.UstawSegmentyOkna(ctx, okno.ID, nowe); err != nil {
		return shared.TranslateSegmentMergeResponse{}, bladTlumaczenia(err)
	}
	panele, err := a.repozytorium.Panele(ctx, okno.ID)
	if err != nil {
		return shared.TranslateSegmentMergeResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateSegmentMergeResponse{
		Segments: nowe,
		Panels:   zlozPaneleTlumaczenia(panele),
	}, nil
}

// PodzielSegment obsługuje `translate.segment.split`. Rozcina wskazany segment
// w podanym miejscu i utrwala nowy podział okna.
func (a *adapterTlumaczenia) PodzielSegment(ctx context.Context,
	z shared.TranslateSegmentSplitRequest) (shared.TranslateSegmentSplitResponse, error) {

	okno, segmenty, err := a.oknoISegmenty(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateSegmentSplitResponse{}, err
	}
	if z.SegmentIndex < 0 || z.SegmentIndex >= len(segmenty) {
		return shared.TranslateSegmentSplitResponse{}, bladWskazaniaTlumaczenia(
			"okno nie ma segmentu o wskazanym numerze")
	}
	segment := []rune(segmenty[z.SegmentIndex])
	if z.Offset <= 0 || z.Offset >= len(segment) {
		return shared.TranslateSegmentSplitResponse{}, bladWskazaniaTlumaczenia(
			"miejsce podziału leży poza segmentem — podział na jego krańcu nie dałby dwóch segmentów")
	}

	nowe := make([]string, 0, len(segmenty)+1)
	nowe = append(nowe, segmenty[:z.SegmentIndex]...)
	nowe = append(nowe,
		strings.TrimSpace(string(segment[:z.Offset])),
		strings.TrimSpace(string(segment[z.Offset:])))
	nowe = append(nowe, segmenty[z.SegmentIndex+1:]...)

	if err := a.repozytorium.UstawSegmentyOkna(ctx, okno.ID, nowe); err != nil {
		return shared.TranslateSegmentSplitResponse{}, bladTlumaczenia(err)
	}
	panele, err := a.repozytorium.Panele(ctx, okno.ID)
	if err != nil {
		return shared.TranslateSegmentSplitResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateSegmentSplitResponse{
		Segments: nowe,
		Panels:   zlozPaneleTlumaczenia(panele),
	}, nil
}

// oknoISegmenty odczytuje okno i jego bieżący podział — wspólne wejście obu
// komend zmiany segmentów.
func (a *adapterTlumaczenia) oknoISegmenty(ctx context.Context,
	kodOkna string) (dane.OknoTlumaczenia, []string, error) {

	okno, err := a.repozytorium.Okno(ctx, kodOkna)
	if err != nil {
		return dane.OknoTlumaczenia{}, nil, bladNieznanegoOkna(kodOkna, err)
	}
	segmenty, err := a.segmentyOkna(ctx, okno)
	if err != nil {
		return dane.OknoTlumaczenia{}, nil, err
	}
	if len(segmenty) == 0 {
		return dane.OknoTlumaczenia{}, nil, bladWskazaniaTlumaczenia(
			"okno " + kodOkna + " nie ma tekstu źródłowego — nie ma czego dzielić ani scalać")
	}
	return okno, segmenty, nil
}

// numerySegmentow sprowadza wskazania kontraktu (napisy) do numerów i pilnuje,
// żeby wszystkie wskazywały segment istniejący. Kontrakt niesie je napisami,
// bo wykaz jest `string[]`; rdzeń nie zgaduje, co znaczy wskazanie nieliczbowe.
func numerySegmentow(wskazania []string, ile int) ([]int, error) {
	numery := make([]int, 0, len(wskazania))
	for _, wskazanie := range wskazania {
		numer, err := liczbaZeWskazania(wskazanie)
		if err != nil {
			return nil, bladWskazaniaTlumaczenia(
				"wskazanie segmentu " + wskazanie + " nie jest numerem")
		}
		if numer < 0 || numer >= ile {
			return nil, bladWskazaniaTlumaczenia(
				"okno nie ma segmentu o numerze " + wskazanie)
		}
		numery = append(numery, numer)
	}
	// Kolejność scalania idzie po numerach, nie po kolejności wskazania —
	// scalenie „3, 1" ma dać ten sam tekst, co „1, 3".
	for i := 1; i < len(numery); i++ {
		for j := i; j > 0 && numery[j] < numery[j-1]; j-- {
			numery[j], numery[j-1] = numery[j-1], numery[j]
		}
	}
	return numery, nil
}
