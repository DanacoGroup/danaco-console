// Odpowiedzialność pliku: obszar słownika modułu Translate — `UstawTerminy`
// (`glossary.set`), `ZastosujSlownik` (`glossary.apply`) i
// `WystapieniaTerminow` (`glossary.occurrences`). Termin oznaczony
// `NieTlumaczyc` zostaje w podmianie nietknięty.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// UstawTerminy zakłada nowy termin słownika albo aktualizuje istniejący po
// `TermId`. Obsługuje `glossary.set`.
func (a *adapterTlumaczenia) UstawTerminy(ctx context.Context,
	z shared.TranslateGlossarySetRequest) (shared.TranslateGlossarySetResponse, error) {

	if z.Source == "" {
		return shared.TranslateGlossarySetResponse{}, bladWskazaniaTlumaczenia("komenda bez treści źródłowej terminu")
	}
	if z.Language == "" {
		return shared.TranslateGlossarySetResponse{}, bladWskazaniaTlumaczenia("komenda bez języka odpowiednika terminu")
	}

	kod := ""
	if z.TermId != nil && *z.TermId != "" {
		kod = *z.TermId
	} else {
		kod = nowyIdentyfikator(przedrostekTerminuSlownika)
	}

	zapisane, err := a.repozytorium.ZapiszTerminy(ctx, []dane.TerminSlownika{{
		Kod:          kod,
		Zrodlo:       z.Source,
		Jezyk:        z.Language,
		Cel:          z.Target,
		NieTlumaczyc: logicznaZWskaznika(z.DoNotTranslate),
		Uwaga:        z.Note,
	}})
	if err != nil {
		return shared.TranslateGlossarySetResponse{}, bladTlumaczenia(err)
	}
	if len(zapisane) == 0 {
		return shared.TranslateGlossarySetResponse{}, bladTlumaczenia(errors.New("zapis terminu słownika nie oddał wiersza"))
	}
	return shared.TranslateGlossarySetResponse{Term: zlozTerminSlownika(zapisane[0])}, nil
}

// ZastosujSlownik ujednolica terminologię treści panelu wedle słownika —
// podmiana tekstu, operacja mechaniczna, bez udziału silnika tłumaczenia.
// Obsługuje `glossary.apply`. Bez wskazanego `PanelId` obejmuje komplet paneli.
func (a *adapterTlumaczenia) ZastosujSlownik(ctx context.Context,
	z shared.TranslateGlossaryApplyRequest) (shared.TranslateGlossaryApplyResponse, error) {

	terminy, err := a.repozytorium.Terminy(ctx)
	if err != nil {
		return shared.TranslateGlossaryApplyResponse{}, bladTlumaczenia(err)
	}

	panele, err := a.paneleDoZastosowania(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateGlossaryApplyResponse{}, err
	}

	razem := 0
	for _, panel := range panele {
		zmienione, err := a.zastosujWPanelu(ctx, panel, terminy)
		if err != nil {
			return shared.TranslateGlossaryApplyResponse{}, err
		}
		razem += zmienione
	}
	return shared.TranslateGlossaryApplyResponse{ChangedCount: razem}, nil
}

// paneleDoZastosowania rozstrzyga zakres komendy.
//
// Pusty `panelId` znaczy „wszystkie panele" — tak stanowi kontrakt, a zbiór
// paneli daje `WszystkiePanele` repozytorium.
func (a *adapterTlumaczenia) paneleDoZastosowania(ctx context.Context,
	kodPanelu *string) ([]dane.PanelTlumaczenia, error) {

	if kodPanelu == nil || *kodPanelu == "" {
		panele, err := a.repozytorium.WszystkiePanele(ctx)
		if err != nil {
			return nil, bladTlumaczenia(err)
		}
		return panele, nil
	}
	panel, err := a.repozytorium.Panel(ctx, *kodPanelu)
	if err != nil {
		return nil, bladNieznanegoPanelu(*kodPanelu, err)
	}
	return []dane.PanelTlumaczenia{panel}, nil
}

// zastosujWPanelu podmienia terminy w jednym panelu i oddaje liczbę
// faktycznie podmienionych wystąpień. Panel bez treści w bazie daje zero.
func (a *adapterTlumaczenia) zastosujWPanelu(ctx context.Context,
	panel dane.PanelTlumaczenia, terminy []dane.TerminSlownika) (int, error) {

	if panel.Tresc == nil || *panel.Tresc == "" {
		return 0, nil
	}
	tresc, zmienione := zastosujTerminySlownika(*panel.Tresc, terminy)
	if zmienione == 0 {
		return 0, nil
	}
	if _, err := a.repozytorium.UstawTlumaczenie(ctx, panel.Kod, &tresc, panel.TrescOdwolanie); err != nil {
		return 0, bladNieznanegoPanelu(panel.Kod, err)
	}
	return zmienione, nil
}

// zastosujTerminySlownika podmienia w treści każdy termin słownika na jego
// odpowiednik docelowy, z wyjątkiem terminów oznaczonych `NieTlumaczyc`.
func zastosujTerminySlownika(tresc string, terminy []dane.TerminSlownika) (string, int) {
	zmienione := 0
	for _, termin := range terminy {
		if termin.NieTlumaczyc || termin.Zrodlo == "" || termin.Cel == nil || *termin.Cel == "" {
			continue
		}
		wystapien := strings.Count(tresc, termin.Zrodlo)
		if wystapien == 0 {
			continue
		}
		tresc = strings.ReplaceAll(tresc, termin.Zrodlo, *termin.Cel)
		zmienione += wystapien
	}
	return tresc, zmienione
}

// WystapieniaTerminow liczy wystąpienia terminu w locie, bez własnej
// tabeli. Obsługuje `glossary.occurrences` i przeszukuje komplet paneli.
func (a *adapterTlumaczenia) WystapieniaTerminow(ctx context.Context,
	z shared.TranslateGlossaryOccurrencesRequest) (shared.TranslateGlossaryOccurrencesResponse, error) {

	if z.Term == "" {
		return shared.TranslateGlossaryOccurrencesResponse{}, bladWskazaniaTlumaczenia("komenda bez wskazania terminu")
	}
	panele, err := a.repozytorium.WszystkiePanele(ctx)
	if err != nil {
		return shared.TranslateGlossaryOccurrencesResponse{}, bladTlumaczenia(err)
	}

	wystapienia := []string{}
	for _, panel := range panele {
		if panel.Tresc == nil {
			continue
		}
		wystapienia = append(wystapienia, otoczeniaTerminu(*panel.Tresc, z.Term)...)
	}
	return shared.TranslateGlossaryOccurrencesResponse{Term: z.Term, Occurrences: wystapienia}, nil
}

// otoczeniaTerminu wycina fragment treści wokół każdego trafienia terminu.
// Szerokość otoczenia jest stała i niewielka — chodzi o pokazanie kontekstu
// Operatorowi, nie o oddanie całej treści panelu przy każdym trafieniu.
func otoczeniaTerminu(tresc, termin string) []string {
	const otoczenie = 40
	wynik := []string{}
	przesuniecie := 0
	for {
		wzgledna := strings.Index(tresc[przesuniecie:], termin)
		if wzgledna < 0 {
			return wynik
		}
		pozycja := przesuniecie + wzgledna
		od := pozycja - otoczenie
		if od < 0 {
			od = 0
		}
		do_ := pozycja + len(termin) + otoczenie
		if do_ > len(tresc) {
			do_ = len(tresc)
		}
		wynik = append(wynik, tresc[od:do_])
		przesuniecie = pozycja + len(termin)
	}
}

// zlozTerminSlownika przekłada wiersz słownika na kontraktowy `GlossaryTerm`,
// w tym stan i dziedzinę doszłe migracją 162.
func zlozTerminSlownika(termin dane.TerminSlownika) shared.GlossaryTerm {
	var nieTlumaczyc *bool
	if termin.NieTlumaczyc {
		wartosc := true
		nieTlumaczyc = &wartosc
	}
	byt := shared.GlossaryTerm{
		Id:             termin.Kod,
		Source:         termin.Zrodlo,
		Target:         termin.Cel,
		Language:       termin.Jezyk,
		DoNotTranslate: nieTlumaczyc,
		Note:           termin.Uwaga,
		UpdatedAt:      termin.Zaktualizowano,
		// Dziedzina i stan doszły z migracją 162; termin bez oznaczenia wychodzi
		// bez stanu.
		Domain: termin.Dziedzina,
	}
	if termin.Stan != nil && strings.TrimSpace(*termin.Stan) != "" {
		stan := shared.GlossaryTermStatus(*termin.Stan)
		byt.Status = &stan
	}
	return byt
}

// logicznaZWskaznika oddaje wartość spod wskaźnika logicznego albo `false`,
// gdy żądanie go nie niosło — `DoNotTranslate` w kontrakcie jest opcjonalne.
func logicznaZWskaznika(wskazanie *bool) bool {
	if wskazanie == nil {
		return false
	}
	return *wskazanie
}
