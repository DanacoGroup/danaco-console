// Odpowiedzialność pliku: obszar słownika modułu Translate — `UstawTerminy`
// (`glossary.set`), `ZastosujSlownik` (`glossary.apply`) i `WystapieniaTerminow`
// (`glossary.occurrences`) na typie `*adapterTlumaczenia`, zadeklarowanym
// w `core/adapter_modul_tlumaczenie.go`. Ten plik nie deklaruje ani typu
// adaptera, ani konstruktora, ani przedrostków identyfikatorów.
//
// Termin oznaczony `NieTlumaczyc` zostaje w podmianie nietknięty — to sens tej
// flagi w `migracja_054_slownik_tlumaczenia.sql`.
//
// Słownik nie żyje wyłącznie tutaj. `glossary.apply` jest narzędziem
// naprawczym: ujednolica terminologię treści, która już powstała. Właściwym
// miejscem słownika jest sam przekład — terminy i zakazy Operatora wchodzą do
// polecenia dla modelu przy `target.add`
// (`adapter_modul_tlumaczenie_polecenia.go`), a `zastosujTerminySlownika` z tego
// pliku przechodzi jeszcze po wyniku modelu jako siatka bezpieczeństwa. Ta sama
// funkcja w dwóch zastosowaniach, nie dwie kopie zasady.
//
// Zakres obu komend bez wskazania panelu bierze się z `WszystkiePanele`
// repozytorium: `glossary.apply` z pustym `PanelId` przechodzi po komplecie
// paneli, a `glossary.occurrences` — które w kontrakcie niesie samo `Term` —
// przeszukuje ten sam komplet.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// UstawTerminy zakłada nowy termin słownika albo aktualizuje istniejący po
// `TermId`. Obsługuje `glossary.set`. Zapis idzie przez `ZapiszTerminy`
// (jedna transakcja, choć niesiona jest tu zawsze lista jednoelementowa) —
// warstwa danych i tak nie ma osobnej drogi INSERT/UPDATE (`dane/slownik.go`),
// więc adapter nie dubluje tego rozróżnienia.
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

// zastosujWPanelu podmienia terminy w jednym panelu i oddaje liczbę faktycznie
// podmienionych wystąpień. Panel bez treści w bazie daje zero — treść żyjąca
// wyłącznie poza bazą (`TrescOdwolanie`) nie jest czytana, bo adapter nie sięga
// po pliki spoza repozytorium.
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
// odpowiednik docelowy — z wyjątkiem terminów oznaczonych `NieTlumaczyc`,
// które zostają nietknięte. Oddaje
// zmienioną treść i liczbę faktycznie podmienionych wystąpień, żeby
// `ChangedCount` mówił prawdę o skutku, nie o liczbie terminów w słowniku.
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

// WystapieniaTerminow liczy wystąpienia terminu w locie, bez własnej tabeli
// (`migracja_054_slownik_tlumaczenia.sql`). Obsługuje `glossary.occurrences`.
//
// Zakresem są wszystkie panele: kontrakt nie niesie wskazania okna ani panelu,
// a słownik jest jeden na instalację, więc jedynym uczciwym odczytaniem jest
// przeszukanie kompletu paneli.
//
// Otoczenie wystąpienia to wycinek treści wokół trafienia — kontrakt chce
// „wystąpień wraz z otoczeniem", a nie samych pozycji.
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

// zlozTerminSlownika przekłada wiersz słownika na kontraktowy `GlossaryTerm`.
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
		// Dziedzina i stan doszły z migracją 162 pod `translate.glossary.list`.
		// Termin, którego nikt nie oznaczył, wychodzi bez stanu — brak oznaczenia
		// nie jest stanem `candidate` ani żadnym innym.
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
