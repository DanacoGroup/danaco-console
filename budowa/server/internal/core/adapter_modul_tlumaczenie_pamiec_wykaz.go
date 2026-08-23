// Odpowiedzialność pliku: pamięć tłumaczeń jako byt Operatora — rodzina
// `translate.memory.*` poza `memory.suggest`, która leży
// w `adapter_modul_tlumaczenie_pamiec.go` i zbiera pary automatycznie.
//
// Dziesięć komend, jedna zasada wspólna: pamięć jest własnością Operatora, a nie
// pochodną paneli. Dlatego para wniesiona ręcznie (`memory.set`) i para z pliku
// wymiany (`memory.import`) nie mają panelu, a usunięcie panelu nie zabiera ze
// sobą par, które z niego kiedyś zdjęto (klucz obcy `ON DELETE SET NULL`,
// migracja 160).
//
// Wymiana idzie standardem TMX, wkompilowanym w rdzeń przez `encoding/xml`
// biblioteki Go — bez ani jednego programu zewnętrznego. Plik CSV (kolumny:
// segment źródłowy, segment docelowy, język) jest drogą drugą, bo tyle właśnie
// eksportuje większość arkuszy, w których Operator prowadzi terminologię.
package core

import (
	"context"
	"encoding/csv"
	"encoding/xml"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekWpisuPamieciTlumaczen znakuje identyfikator pary pamięci nadany przez rdzeń.
const przedrostekWpisuPamieciTlumaczen = "wpm-"

// progDopasowaniaDomyslny obowiązuje, gdy ani żądanie, ani polityka okna nie
// wskazały progu. Siedemdziesiąt pięć procent to próg, poniżej którego
// podpowiedź częściej myli, niż pomaga — para różniąca się co czwartym znakiem
// nie jest tym samym zdaniem.
const progDopasowaniaDomyslny = 75

// WykazPamieci obsługuje `translate.memory.list`.
func (a *adapterTlumaczenia) WykazPamieci(ctx context.Context,
	z shared.TranslateMemoryListRequest) (shared.TranslateMemoryListResponse, error) {

	filtr := dane.FiltrPamieciTlumaczen{
		Jezyk:   napisZeWskaznika(z.Language),
		Fraza:   napisZeWskaznika(z.Query),
		Projekt: napisZeWskaznika(z.Project),
		Limit:   liczbaCalkowitaZeWskaznika(z.Limit),
		Offset:  liczbaCalkowitaZeWskaznika(z.Offset),
	}
	if z.Scope != nil {
		filtr.Zasieg = string(*z.Scope)
	}
	wpisy, razem, err := a.repozytorium.WpisyPamieci(ctx, filtr)
	if err != nil {
		return shared.TranslateMemoryListResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateMemoryListResponse{
		Entries: zlozWpisyPamieci(wpisy),
		Total:   razem,
	}, nil
}

// UstawWpisPamieci obsługuje `translate.memory.set` — parę wniesioną wprost
// przez Operatora albo poprawkę pary zastanej.
func (a *adapterTlumaczenia) UstawWpisPamieci(ctx context.Context,
	z shared.TranslateMemorySetRequest) (shared.TranslateMemorySetResponse, error) {

	if strings.TrimSpace(z.Language) == "" {
		return shared.TranslateMemorySetResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.memory.set bez języka pary")
	}
	if strings.TrimSpace(z.SourceSegment) == "" || strings.TrimSpace(z.TargetSegment) == "" {
		return shared.TranslateMemorySetResponse{}, bladWskazaniaTlumaczenia(
			"para pamięci bez segmentu źródłowego albo docelowego — pamięć nie przyjmuje połówek")
	}

	kod := napisZeWskaznika(z.EntryId)
	if strings.TrimSpace(kod) == "" {
		kod = nowyIdentyfikator(przedrostekWpisuPamieciTlumaczen)
	}
	wpis := dane.WpisPamieciTlumaczenPelny{
		Kod:             kod,
		Jezyk:           z.Language,
		SegmentZrodlowy: z.SourceSegment,
		SegmentDocelowy: z.TargetSegment,
		Projekt:         z.Project,
		Klient:          z.Client,
		Zasieg:          "card",
	}
	// Pole `context` kontraktu niesie sąsiedztwo pary jednym napisem — rdzeń
	// zapisuje je jako kontekst poprzedzający, bo tam trafia, gdy pary rodzą się
	// z paneli, i tam go szuka dopasowanie kontekstowe.
	if z.Context != nil && strings.TrimSpace(*z.Context) != "" {
		wpis.KontekstPoprzedni = z.Context
	}
	zapisany, err := a.repozytorium.ZapiszWpisPamieci(ctx, wpis)
	if err != nil {
		return shared.TranslateMemorySetResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateMemorySetResponse{Entry: zlozWpisPamieci(zapisany)}, nil
}

// UsunWpisPamieci obsługuje `translate.memory.delete`.
func (a *adapterTlumaczenia) UsunWpisPamieci(ctx context.Context,
	z shared.TranslateMemoryDeleteRequest) (shared.TranslateMemoryDeleteResponse, error) {

	if strings.TrimSpace(z.EntryId) == "" {
		return shared.TranslateMemoryDeleteResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.memory.delete bez wskazania pary")
	}
	zeszlo, err := a.repozytorium.UsunWpisPamieci(ctx, z.EntryId)
	if err != nil {
		return shared.TranslateMemoryDeleteResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateMemoryDeleteResponse{Deleted: zeszlo}, nil
}

// ImportujPamiec obsługuje `translate.memory.import`. Czyta plik TMX albo CSV
// spod wskazanej ścieżki i wnosi z niego pary. Para już obecna w pamięci (ten
// sam język i ten sam segment źródłowy) jest pomijana, chyba że Operator
// zażądał nadpisania — stąd dwie liczby w odpowiedzi.
func (a *adapterTlumaczenia) ImportujPamiec(ctx context.Context,
	z shared.TranslateMemoryImportRequest) (shared.TranslateMemoryImportResponse, error) {

	sciezka := strings.TrimSpace(z.Path)
	if sciezka == "" {
		return shared.TranslateMemoryImportResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.memory.import bez ścieżki pliku wymiany")
	}
	pary, err := wczytajParyWymiany(sciezka)
	if err != nil {
		return shared.TranslateMemoryImportResponse{}, err
	}

	nadpisuj := z.Overwrite != nil && *z.Overwrite
	zaimportowane, pominiete := 0, 0
	for _, para := range pary {
		zastane, _, err := a.repozytorium.WpisyPamieci(ctx, dane.FiltrPamieciTlumaczen{
			Jezyk: para.Jezyk,
			Fraza: para.SegmentZrodlowy,
		})
		if err != nil {
			return shared.TranslateMemoryImportResponse{}, bladTlumaczenia(err)
		}
		istniejacy := ""
		for _, wpis := range zastane {
			if wpis.SegmentZrodlowy == para.SegmentZrodlowy {
				istniejacy = wpis.Kod
				break
			}
		}
		if istniejacy != "" && !nadpisuj {
			pominiete++
			continue
		}
		wpis := para
		wpis.Kod = istniejacy
		if wpis.Kod == "" {
			wpis.Kod = nowyIdentyfikator(przedrostekWpisuPamieciTlumaczen)
		}
		wpis.Projekt = z.Project
		if _, err := a.repozytorium.ZapiszWpisPamieci(ctx, wpis); err != nil {
			return shared.TranslateMemoryImportResponse{}, bladTlumaczenia(err)
		}
		zaimportowane++
	}
	return shared.TranslateMemoryImportResponse{
		ImportedCount: zaimportowane,
		SkippedCount:  pominiete,
	}, nil
}

// EksportujPamiec obsługuje `translate.memory.export` — wypisuje pary do pliku
// TMX (albo CSV, gdy taka jest końcówka ścieżki). Plik powstaje naprawdę:
// odpowiedź niesie ścieżkę, pod którą leży wynik, a nie ścieżkę zapowiedzianą.
func (a *adapterTlumaczenia) EksportujPamiec(ctx context.Context,
	z shared.TranslateMemoryExportRequest) (shared.TranslateMemoryExportResponse, error) {

	sciezka := strings.TrimSpace(z.Path)
	if sciezka == "" {
		return shared.TranslateMemoryExportResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.memory.export bez ścieżki pliku wyniku")
	}
	wpisy, _, err := a.repozytorium.WpisyPamieci(ctx, dane.FiltrPamieciTlumaczen{
		Jezyk:   napisZeWskaznika(z.Language),
		Projekt: napisZeWskaznika(z.Project),
	})
	if err != nil {
		return shared.TranslateMemoryExportResponse{}, bladTlumaczenia(err)
	}
	if err := zapiszParyWymiany(sciezka, wpisy); err != nil {
		return shared.TranslateMemoryExportResponse{}, err
	}
	return shared.TranslateMemoryExportResponse{
		ExportedCount: len(wpisy),
		Path:          sciezka,
	}, nil
}

// UtrzymajPamiec obsługuje `translate.memory.maintain`. Cztery rodzaje pracy
// kontraktu wykonuje naprawdę, a `dryRun` przelicza to samo bez zapisu — żeby
// Operator zobaczył zasięg czynności przed jej wykonaniem, a nie po.
func (a *adapterTlumaczenia) UtrzymajPamiec(ctx context.Context,
	z shared.TranslateMemoryMaintainRequest) (shared.TranslateMemoryMaintainResponse, error) {

	naSucho := z.DryRun != nil && *z.DryRun
	wpisy, _, err := a.repozytorium.WpisyPamieci(ctx, dane.FiltrPamieciTlumaczen{
		Jezyk:   napisZeWskaznika(z.Language),
		Projekt: napisZeWskaznika(z.Project),
	})
	if err != nil {
		return shared.TranslateMemoryMaintainResponse{}, bladTlumaczenia(err)
	}

	switch z.Kind {
	case shared.TranslationMemoryMaintenanceKindDeduplicate,
		shared.TranslationMemoryMaintenanceKindMerge:
		// Duplikat to para o tym samym języku i tej samej parze segmentów.
		// Scalanie idzie o krok dalej: zderza pary o tym samym segmencie
		// źródłowym i zostawia najświeższą, bo to ona niesie ostatnie
		// rozstrzygnięcie Operatora.
		poScaleniu := z.Kind == shared.TranslationMemoryMaintenanceKindMerge
		doUsuniecia := nadmiaroweParyPamieci(wpisy, poScaleniu)
		if naSucho {
			return shared.TranslateMemoryMaintainResponse{
				AffectedCount: len(doUsuniecia), DryRun: true}, nil
		}
		zeszlo, err := a.repozytorium.UsunWpisyPamieci(ctx, doUsuniecia)
		if err != nil {
			return shared.TranslateMemoryMaintainResponse{}, bladTlumaczenia(err)
		}
		return shared.TranslateMemoryMaintainResponse{AffectedCount: zeszlo}, nil

	case shared.TranslationMemoryMaintenanceKindReplace:
		szukane := napisZeWskaznika(z.Search)
		if strings.TrimSpace(szukane) == "" {
			return shared.TranslateMemoryMaintainResponse{}, bladWskazaniaTlumaczenia(
				"masowa podmiana bez wskazania szukanego fragmentu — rdzeń nie zgaduje, co podmienić")
		}
		zamiennik := napisZeWskaznika(z.Replacement)
		dotkniete := 0
		for _, wpis := range wpisy {
			if !strings.Contains(wpis.SegmentDocelowy, szukane) &&
				!strings.Contains(wpis.SegmentZrodlowy, szukane) {
				continue
			}
			dotkniete++
			if naSucho {
				continue
			}
			zmieniony := wpis
			zmieniony.SegmentZrodlowy = strings.ReplaceAll(wpis.SegmentZrodlowy, szukane, zamiennik)
			zmieniony.SegmentDocelowy = strings.ReplaceAll(wpis.SegmentDocelowy, szukane, zamiennik)
			if _, err := a.repozytorium.ZapiszWpisPamieci(ctx, zmieniony); err != nil {
				return shared.TranslateMemoryMaintainResponse{}, bladTlumaczenia(err)
			}
		}
		return shared.TranslateMemoryMaintainResponse{AffectedCount: dotkniete, DryRun: naSucho}, nil

	case shared.TranslationMemoryMaintenanceKindFilter:
		szukane := napisZeWskaznika(z.Search)
		if strings.TrimSpace(szukane) == "" {
			return shared.TranslateMemoryMaintainResponse{}, bladWskazaniaTlumaczenia(
				"usunięcie par bez zawężenia skasowałoby całą pamięć — zawężenie jest wymagane")
		}
		kody := []string{}
		for _, wpis := range wpisy {
			if strings.Contains(wpis.SegmentZrodlowy, szukane) ||
				strings.Contains(wpis.SegmentDocelowy, szukane) {
				kody = append(kody, wpis.Kod)
			}
		}
		if naSucho {
			return shared.TranslateMemoryMaintainResponse{AffectedCount: len(kody), DryRun: true}, nil
		}
		zeszlo, err := a.repozytorium.UsunWpisyPamieci(ctx, kody)
		if err != nil {
			return shared.TranslateMemoryMaintainResponse{}, bladTlumaczenia(err)
		}
		return shared.TranslateMemoryMaintainResponse{AffectedCount: zeszlo}, nil
	}

	return shared.TranslateMemoryMaintainResponse{}, bladWskazaniaTlumaczenia(
		"nieznany rodzaj utrzymania pamięci: " + string(z.Kind))
}

// nadmiaroweParyPamieci wskazuje pary do zdjęcia. Zostaje zawsze para
// najświeższa; kluczem porównania jest para segmentów (czyszczenie duplikatów)
// albo sam segment źródłowy (scalanie).
func nadmiaroweParyPamieci(wpisy []dane.WpisPamieciTlumaczenPelny, poZrodle bool) []string {
	najswiezsze := map[string]dane.WpisPamieciTlumaczenPelny{}
	for _, wpis := range wpisy {
		klucz := wpis.Jezyk + "\x00" + wpis.SegmentZrodlowy
		if !poZrodle {
			klucz += "\x00" + wpis.SegmentDocelowy
		}
		zastany, jest := najswiezsze[klucz]
		if !jest || wpis.Zaktualizowano > zastany.Zaktualizowano {
			najswiezsze[klucz] = wpis
		}
	}
	zostaja := map[string]struct{}{}
	for _, wpis := range najswiezsze {
		zostaja[wpis.Kod] = struct{}{}
	}
	kody := []string{}
	for _, wpis := range wpisy {
		if _, zostaje := zostaja[wpis.Kod]; !zostaje {
			kody = append(kody, wpis.Kod)
		}
	}
	sort.Strings(kody)
	return kody
}

// TlumaczWstepnie obsługuje `translate.memory.pretranslate`. Wypełnia panele
// treścią złożoną z par pamięci: segment po segmencie, biorąc dopasowanie
// powyżej progu. Segment bez dopasowania zostaje w brzmieniu źródłowym —
// tłumaczenie wstępne nie udaje, że przetłumaczyło wszystko.
func (a *adapterTlumaczenia) TlumaczWstepnie(ctx context.Context,
	z shared.TranslateMemoryPretranslateRequest) (shared.TranslateMemoryPretranslateResponse, error) {

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateMemoryPretranslateResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	segmenty, err := a.segmentyOkna(ctx, okno)
	if err != nil {
		return shared.TranslateMemoryPretranslateResponse{}, err
	}
	if len(segmenty) == 0 {
		return shared.TranslateMemoryPretranslateResponse{}, bladWskazaniaTlumaczenia(
			"okno " + okno.Kod + " nie ma tekstu źródłowego — tłumaczenie wstępne nie ma czego dopasować")
	}

	prog := liczbaCalkowitaZeWskaznika(z.Threshold)
	if prog <= 0 {
		prog = a.progPolityki(ctx, okno.ID)
	}
	tylkoPuste := z.OnlyEmpty == nil || *z.OnlyEmpty

	panele, err := a.repozytorium.Panele(ctx, okno.ID)
	if err != nil {
		return shared.TranslateMemoryPretranslateResponse{}, bladTlumaczenia(err)
	}

	wskazanyPanel := napisZeWskaznika(z.PanelId)
	wypelnione := 0
	for _, panel := range panele {
		if wskazanyPanel != "" && panel.Kod != wskazanyPanel {
			continue
		}
		if tylkoPuste && panel.Tresc != nil && strings.TrimSpace(*panel.Tresc) != "" {
			continue
		}
		pary, _, err := a.repozytorium.WpisyPamieci(ctx,
			dane.FiltrPamieciTlumaczen{Jezyk: panel.Jezyk})
		if err != nil {
			return shared.TranslateMemoryPretranslateResponse{}, bladTlumaczenia(err)
		}
		zlozone := make([]string, 0, len(segmenty))
		trafienia := 0
		for _, segment := range segmenty {
			najlepszy, wynik := najlepszaPara(pary, segment)
			if wynik >= prog && najlepszy != "" {
				zlozone = append(zlozone, najlepszy)
				trafienia++
				continue
			}
			zlozone = append(zlozone, segment)
		}
		if trafienia == 0 {
			continue
		}
		tresc := strings.Join(zlozone, " ")
		if _, err := a.repozytorium.UstawTlumaczenie(ctx, panel.Kod, &tresc, nil); err != nil {
			return shared.TranslateMemoryPretranslateResponse{}, bladTlumaczenia(err)
		}
		wypelnione++
	}

	poZmianie, err := a.repozytorium.Panele(ctx, okno.ID)
	if err != nil {
		return shared.TranslateMemoryPretranslateResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateMemoryPretranslateResponse{
		FilledCount: wypelnione,
		Panels:      zlozPaneleTlumaczenia(poZmianie),
	}, nil
}

// progPolityki bierze próg dopasowania z polityki okna, a przy jej braku —
// z progu domyślnego modułu. Brak polityki nie jest usterką, tylko stanem
// okna, którego nikt nie nastawiał.
func (a *adapterTlumaczenia) progPolityki(ctx context.Context, oknoID int64) int {
	polityka, err := a.repozytorium.PolitykaPamieci(ctx, oknoID)
	if err != nil || polityka.Prog <= 0 {
		return progDopasowaniaDomyslny
	}
	return int(polityka.Prog)
}

// najlepszaPara wybiera z pamięci parę najbliższą segmentowi i oddaje jej
// segment docelowy wraz z miarą podobieństwa w procentach.
func najlepszaPara(pary []dane.WpisPamieciTlumaczenPelny, segment string) (string, int) {
	najlepszy, najlepszyWynik := "", 0
	for _, para := range pary {
		wynik := podobienstwoSegmentow(segment, para.SegmentZrodlowy)
		if wynik > najlepszyWynik {
			najlepszy, najlepszyWynik = para.SegmentDocelowy, wynik
		}
	}
	return najlepszy, najlepszyWynik
}

// WyrownajTeksty obsługuje `translate.memory.align`. Dzieli oba teksty na
// zdania i paruje je po kolei. Miara pary (`score`) mówi, na ile podział obu
// stron jest zgodny: przy równej liczbie zdań para jest pewna, przy nierównej
// wyrównanie jest domysłem i miara to pokazuje, zamiast milczeć.
func (a *adapterTlumaczenia) WyrownajTeksty(ctx context.Context,
	z shared.TranslateMemoryAlignRequest) (shared.TranslateMemoryAlignResponse, error) {

	if strings.TrimSpace(z.SourceText) == "" || strings.TrimSpace(z.TargetText) == "" {
		return shared.TranslateMemoryAlignResponse{}, bladWskazaniaTlumaczenia(
			"wyrównanie bez obu tekstów — nie ma czego z czym parować")
	}
	if strings.TrimSpace(z.Language) == "" {
		return shared.TranslateMemoryAlignResponse{}, bladWskazaniaTlumaczenia(
			"wyrównanie bez języka par — pamięć nie przyjmuje pary bez języka")
	}

	zrodlowe := podzielNaZdania(z.SourceText)
	docelowe := podzielNaZdania(z.TargetText)
	pewnosc := 100
	if len(zrodlowe) != len(docelowe) {
		// Podział rozjechał się między stronami — para po numerze jest wtedy
		// domysłem, nie ustaleniem. Miara spada proporcjonalnie do rozjazdu.
		pewnosc = 100 * min(len(zrodlowe), len(docelowe)) / max(len(zrodlowe), len(docelowe))
	}

	pary := []shared.BitextPair{}
	wpisane := 0
	liczba := min(len(zrodlowe), len(docelowe))
	for i := 0; i < liczba; i++ {
		pary = append(pary, shared.BitextPair{
			SourceSegment: zrodlowe[i],
			TargetSegment: docelowe[i],
			Score:         pewnosc,
		})
		if z.Commit == nil || !*z.Commit {
			continue
		}
		if _, err := a.repozytorium.ZapiszWpisPamieci(ctx, dane.WpisPamieciTlumaczenPelny{
			Kod:             nowyIdentyfikator(przedrostekWpisuPamieciTlumaczen),
			Jezyk:           z.Language,
			SegmentZrodlowy: zrodlowe[i],
			SegmentDocelowy: docelowe[i],
			Zasieg:          "card",
		}); err != nil {
			return shared.TranslateMemoryAlignResponse{}, bladTlumaczenia(err)
		}
		wpisane++
	}
	return shared.TranslateMemoryAlignResponse{Pairs: pary, CommittedCount: wpisane}, nil
}

// PolitykaPamieci obsługuje `translate.memory.policy.get`. Okno bez polityki
// dostaje politykę domyślną modułu — nastawa nieustawiona jest nastawą
// domyślną, nie odmową.
func (a *adapterTlumaczenia) PolitykaPamieci(ctx context.Context,
	z shared.TranslateMemoryPolicyGetRequest) (shared.TranslateMemoryPolicyGetResponse, error) {

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateMemoryPolicyGetResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	polityka, err := a.repozytorium.PolitykaPamieci(ctx, okno.ID)
	if err != nil {
		return shared.TranslateMemoryPolicyGetResponse{
			Policy: shared.TranslationMemoryPolicy{
				WindowId:  okno.Kod,
				Scope:     shared.TranslationMemoryScopeCard,
				Threshold: progDopasowaniaDomyslny,
				UpdatedAt: okno.Zaktualizowano,
			}}, nil
	}
	return shared.TranslateMemoryPolicyGetResponse{
		Policy: zlozPolitykePamieci(okno.Kod, polityka)}, nil
}

// UstawPolitykePamieci obsługuje `translate.memory.policy.set`. Pola nieobecne
// w żądaniu zostają takie, jakie były — nastawa cząstkowa nie ma prawa zerować
// nastaw, o których żądanie milczy.
func (a *adapterTlumaczenia) UstawPolitykePamieci(ctx context.Context,
	z shared.TranslateMemoryPolicySetRequest) (shared.TranslateMemoryPolicySetResponse, error) {

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateMemoryPolicySetResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}
	polityka, err := a.repozytorium.PolitykaPamieci(ctx, okno.ID)
	if err != nil {
		polityka = dane.PolitykaPamieciOkna{
			OknoID: okno.ID,
			Zasieg: string(shared.TranslationMemoryScopeCard),
			Prog:   progDopasowaniaDomyslny,
		}
	}
	polityka.OknoID = okno.ID
	if z.Scope != nil {
		polityka.Zasieg = string(*z.Scope)
	}
	if z.Threshold != nil {
		polityka.Prog = int64(*z.Threshold)
	}
	if z.ContextMatch != nil {
		polityka.DopasowanieKontekstu = *z.ContextMatch
	}
	if z.PreTranslate != nil {
		polityka.WstepneTlumaczenie = *z.PreTranslate
	}
	zapisana, err := a.repozytorium.ZapiszPolitykePamieci(ctx, polityka)
	if err != nil {
		return shared.TranslateMemoryPolicySetResponse{}, bladTlumaczenia(err)
	}
	return shared.TranslateMemoryPolicySetResponse{
		Policy: zlozPolitykePamieci(okno.Kod, zapisana)}, nil
}

// zlozPolitykePamieci przekłada wiersz polityki na byt kontraktu.
func zlozPolitykePamieci(oknoKod string, polityka dane.PolitykaPamieciOkna) shared.TranslationMemoryPolicy {
	return shared.TranslationMemoryPolicy{
		WindowId:     oknoKod,
		Scope:        shared.TranslationMemoryScope(polityka.Zasieg),
		Threshold:    int(polityka.Prog),
		ContextMatch: polityka.DopasowanieKontekstu,
		PreTranslate: polityka.WstepneTlumaczenie,
		UpdatedAt:    polityka.Zaktualizowano,
	}
}

// zlozWpisyPamieci przekłada wiersze pamięci na byty kontraktu.
func zlozWpisyPamieci(wpisy []dane.WpisPamieciTlumaczenPelny) []shared.TranslationMemoryEntry {
	lista := make([]shared.TranslationMemoryEntry, 0, len(wpisy))
	for _, wpis := range wpisy {
		lista = append(lista, zlozWpisPamieci(wpis))
	}
	return lista
}

// zlozWpisPamieci przekłada jeden wiersz pamięci na byt kontraktu.
func zlozWpisPamieci(wpis dane.WpisPamieciTlumaczenPelny) shared.TranslationMemoryEntry {
	return shared.TranslationMemoryEntry{
		Id:              wpis.Kod,
		Language:        wpis.Jezyk,
		SourceSegment:   wpis.SegmentZrodlowy,
		TargetSegment:   wpis.SegmentDocelowy,
		PanelId:         wpis.PanelKod,
		Project:         wpis.Projekt,
		Client:          wpis.Klient,
		Author:          wpis.Autor,
		PreviousSegment: wpis.KontekstPoprzedni,
		NextSegment:     wpis.KontekstNastepny,
		CreatedAt:       wpis.Utworzono,
		UpdatedAt:       wpis.Zaktualizowano,
	}
}

// tmxDokument, tmxCialo, tmxJednostka i tmxWariant odwzorowują tyle standardu
// TMX, ile niesie para pamięci: jednostkę tłumaczeniową z dwoma wariantami
// językowymi. Pełny TMX ma nadto nagłówek z metrykami narzędzia i noty — rdzeń
// wypisuje nagłówek minimalny i czyta plik, nie wymagając niczego ponad `<tu>`.
type tmxDokument struct {
	XMLName xml.Name  `xml:"tmx"`
	Wersja  string    `xml:"version,attr"`
	Naglow  tmxNaglow `xml:"header"`
	Cialo   tmxCialo  `xml:"body"`
}

type tmxNaglow struct {
	Narzedzie     string `xml:"creationtool,attr"`
	WersjaNarzedz string `xml:"creationtoolversion,attr"`
	Segmentacja   string `xml:"segtype,attr"`
	JezykZrodlowy string `xml:"srclang,attr"`
	FormatDanych  string `xml:"o-tmf,attr"`
	Kodowanie     string `xml:"datatype,attr"`
}

type tmxCialo struct {
	Jednostki []tmxJednostka `xml:"tu"`
}

type tmxJednostka struct {
	Warianty []tmxWariant `xml:"tuv"`
}

type tmxWariant struct {
	Jezyk   string `xml:"http://www.w3.org/XML/1998/namespace lang,attr"`
	Segment string `xml:"seg"`
}

// wczytajParyWymiany czyta pary z pliku TMX albo CSV. Rozstrzyga końcówka
// nazwy: rdzeń nie zgaduje formatu po zawartości, bo plik wymiany bez rozszerzenia
// jest brakiem wskazania, nie zagadką do rozwiązania.
func wczytajParyWymiany(sciezka string) ([]dane.WpisPamieciTlumaczenPelny, error) {
	bajty, err := os.ReadFile(filepath.Clean(sciezka))
	if err != nil {
		return nil, bladPlikuTlumaczenia(sciezka, err)
	}
	switch strings.ToLower(filepath.Ext(sciezka)) {
	case ".tmx", ".xml":
		var dokument tmxDokument
		if err := xml.Unmarshal(bajty, &dokument); err != nil {
			return nil, bladWskazaniaTlumaczenia("plik " + sciezka + " nie jest czytelnym TMX: " + err.Error())
		}
		pary := []dane.WpisPamieciTlumaczenPelny{}
		for _, jednostka := range dokument.Jednostki() {
			pary = append(pary, jednostka)
		}
		return pary, nil
	case ".csv":
		czytnik := csv.NewReader(strings.NewReader(string(bajty)))
		czytnik.FieldsPerRecord = -1
		wiersze, err := czytnik.ReadAll()
		if err != nil {
			return nil, bladWskazaniaTlumaczenia("plik " + sciezka + " nie jest czytelnym CSV: " + err.Error())
		}
		pary := []dane.WpisPamieciTlumaczenPelny{}
		for _, wiersz := range wiersze {
			if len(wiersz) < 3 || strings.TrimSpace(wiersz[0]) == "" {
				continue
			}
			pary = append(pary, dane.WpisPamieciTlumaczenPelny{
				SegmentZrodlowy: wiersz[0],
				SegmentDocelowy: wiersz[1],
				Jezyk:           strings.TrimSpace(wiersz[2]),
				Zasieg:          "card",
			})
		}
		return pary, nil
	}
	return nil, bladWskazaniaTlumaczenia(
		"pamięć wymienia się plikiem TMX albo CSV; ścieżka " + sciezka + " nie ma żadnej z tych końcówek")
}

// Jednostki przekłada jednostki TMX na pary pamięci. Jednostka o mniej niż
// dwóch wariantach nie jest parą — zostaje pominięta, bo połowa pary nie
// niesie żadnego przekładu.
func (d tmxDokument) Jednostki() []dane.WpisPamieciTlumaczenPelny {
	pary := []dane.WpisPamieciTlumaczenPelny{}
	for _, jednostka := range d.Cialo.Jednostki {
		if len(jednostka.Warianty) < 2 {
			continue
		}
		pary = append(pary, dane.WpisPamieciTlumaczenPelny{
			SegmentZrodlowy: jednostka.Warianty[0].Segment,
			SegmentDocelowy: jednostka.Warianty[1].Segment,
			Jezyk:           jednostka.Warianty[1].Jezyk,
			Zasieg:          "card",
		})
	}
	return pary
}

// zapiszParyWymiany wypisuje pary do pliku TMX albo CSV.
func zapiszParyWymiany(sciezka string, wpisy []dane.WpisPamieciTlumaczenPelny) error {
	if err := os.MkdirAll(filepath.Dir(sciezka), 0o700); err != nil {
		return bladPlikuTlumaczenia(sciezka, err)
	}
	switch strings.ToLower(filepath.Ext(sciezka)) {
	case ".csv":
		var bufor strings.Builder
		pisarz := csv.NewWriter(&bufor)
		for _, wpis := range wpisy {
			if err := pisarz.Write([]string{wpis.SegmentZrodlowy, wpis.SegmentDocelowy, wpis.Jezyk}); err != nil {
				return bladTlumaczenia(err)
			}
		}
		pisarz.Flush()
		if err := os.WriteFile(sciezka, []byte(bufor.String()), 0o600); err != nil {
			return bladPlikuTlumaczenia(sciezka, err)
		}
		return nil
	default:
		dokument := tmxDokument{
			Wersja: "1.4",
			Naglow: tmxNaglow{
				Narzedzie:     "Danaco Console",
				WersjaNarzedz: shared.ProtocolVersion,
				Segmentacja:   "sentence",
				JezykZrodlowy: "*all*",
				FormatDanych:  "DanacoTM",
				Kodowanie:     "plaintext",
			},
		}
		for _, wpis := range wpisy {
			dokument.Cialo.Jednostki = append(dokument.Cialo.Jednostki, tmxJednostka{
				Warianty: []tmxWariant{
					{Jezyk: "x-source", Segment: wpis.SegmentZrodlowy},
					{Jezyk: wpis.Jezyk, Segment: wpis.SegmentDocelowy},
				},
			})
		}
		tresc, err := xml.MarshalIndent(dokument, "", "  ")
		if err != nil {
			return bladTlumaczenia(err)
		}
		tresc = append([]byte(xml.Header), tresc...)
		if err := os.WriteFile(sciezka, tresc, 0o600); err != nil {
			return bladPlikuTlumaczenia(sciezka, err)
		}
		return nil
	}
}
