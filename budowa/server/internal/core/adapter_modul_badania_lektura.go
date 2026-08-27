// Pakiet obsługuje widok lektury badania: otwarcie źródła, adnotacje i wypisy,
// wydobycie tabel i twierdzeń, rozpoznanie pisma oraz rozmowę opartą na
// korpusie źródeł komendą `research.corpus.ask`, zakotwiczoną w ich treści.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// granicaLekturyBadania jest domyślną liczbą znaków jednej strony lektury,
// gdy żądanie otwarcia źródła nie poda własnej granicy stronicowania.
const granicaLekturyBadania = 20000

// OtworzDoLektury obsługuje `research.reading.open`: oddaje tekst źródła
// stronicowany, wydobywając go z załącznika, gdy jeszcze go nie ma.
func (a *adapterBadan) OtworzDoLektury(ctx context.Context,
	z shared.ResearchReadingOpenRequest) (shared.ResearchReadingOpenResponse, error) {

	if z.SourceId == "" {
		return shared.ResearchReadingOpenResponse{}, bladWskazaniaBadan("reading.open bez źródła")
	}
	tekst, warstwa, err := a.trescDoLekturyBadania(ctx, z.SourceId, false)
	if err != nil {
		return shared.ResearchReadingOpenResponse{}, err
	}

	granica := granicaLekturyBadania
	if z.MaxChars != nil && *z.MaxChars > 0 {
		granica = *z.MaxChars
	}
	znaki := []rune(tekst)
	stron := (len(znaki) + granica - 1) / granica
	if stron == 0 {
		stron = 1
	}
	strona := 1
	if z.Page != nil && *z.Page > 0 {
		strona = *z.Page
	}
	if strona > stron {
		strona = stron
	}
	od := (strona - 1) * granica
	do := od + granica
	if od > len(znaki) {
		od = len(znaki)
	}
	if do > len(znaki) {
		do = len(znaki)
	}
	wycinek := string(znaki[od:do])
	obciety := do < len(znaki)

	return shared.ResearchReadingOpenResponse{Content: shared.ResearchReadingContent{
		SourceId: z.SourceId, Text: &wycinek, Page: &strona, PageCount: &stron,
		Truncated: &obciety, HasTextLayer: &warstwa,
	}}, nil
}

// trescDoLekturyBadania oddaje tekst źródła, wydobywając go z załącznika, gdy
// jeszcze go nie ma. `wymusRozpoznanie` żąda rozpoznania pisma nawet wtedy, gdy
// dokument ma warstwę tekstową — tego chce `research.source.ocr`.
func (a *adapterBadan) trescDoLekturyBadania(ctx context.Context, kodZrodla string,
	wymusRozpoznanie bool) (string, bool, error) {

	tresc, err := a.repozytorium.TrescZrodlaBadania(ctx, kodZrodla)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return "", false, bladNieznanegoZrodlaBadania(kodZrodla)
	}
	if err != nil {
		return "", false, bladBadan(err)
	}
	if !wymusRozpoznanie && tresc.Tekst != nil && strings.TrimSpace(*tresc.Tekst) != "" {
		return *tresc.Tekst, tresc.WarstwaTekstu, nil
	}

	sciezka, err := a.sciezkaMaterialuBadania(ctx, kodZrodla)
	if err != nil {
		return "", false, err
	}
	if a.dokumenty == nil {
		return "", false, protokolBladBadania(shared.ErrorCodeInternalError,
			"arsenał dokumentowy nie jest wpięty — naprawa: podpiąć port dokumentów "+
				"przy składaniu rdzenia")
	}
	wymuszone := wymusRozpoznanie
	wynik, err := a.dokumenty.WyciagnijTekst(ctx, shared.DocumentTextExtractRequest{
		SourcePath: &sciezka, ForceOcr: &wymuszone,
	})
	if err != nil {
		return "", false, err
	}
	warstwa := !wynik.UsedOcr
	stron := int64(0)
	if wynik.Pages != nil {
		stron = int64(*wynik.Pages)
	}
	tekst := wynik.Text
	zapis := dane.TrescZrodlaBadania{Tekst: &tekst, WarstwaTekstu: warstwa}
	if stron > 0 {
		zapis.Stron = &stron
	}
	if err := a.repozytorium.UstawTrescZrodla(ctx, kodZrodla, zapis); err != nil {
		return "", false, bladBadan(err)
	}
	return tekst, warstwa, nil
}

// sciezkaMaterialuBadania wskazuje plik, z którego da się wydobyć treść źródła:
// załącznik pełnotekstowy albo migawkę. Brak obu to odmowa nazwana — Operator
// ma wiedzieć, że najpierw trzeba dołożyć pełny tekst.
func (a *adapterBadan) sciezkaMaterialuBadania(ctx context.Context, kodZrodla string) (string, error) {
	zalaczniki, err := a.repozytorium.Zalaczniki(ctx, kodZrodla)
	if err != nil {
		return "", bladBadan(err)
	}
	for _, rodzaj := range []string{
		string(shared.ResearchAttachmentKindFulltext),
		string(shared.ResearchAttachmentKindSnapshot),
		string(shared.ResearchAttachmentKindData),
	} {
		for _, zalacznik := range zalaczniki {
			if zalacznik.Rodzaj == rodzaj && zalacznik.Sciezka != nil &&
				strings.TrimSpace(*zalacznik.Sciezka) != "" {
				return *zalacznik.Sciezka, nil
			}
		}
	}
	return "", protokolBladBadania(shared.ErrorCodeNotFound,
		"źródło "+kodZrodla+" nie ma ani treści, ani załącznika z bajtami — "+
			"naprawa: dołożyć pełny tekst przez research.source.attachment.add "+
			"albo pozyskać stronę przez research.source.capture")
}

// ── Adnotacje i wypisy ─────────────────────────────────────────────────────

// DodajAdnotacje obsługuje `research.annotation.add`: zapisuje adnotację albo
// wypis przy źródle, wraz z kotwicą wskazującą miejsce w tekście.
func (a *adapterBadan) DodajAdnotacje(ctx context.Context,
	z shared.ResearchAnnotationAddRequest) (shared.ResearchAnnotationAddResponse, error) {

	if z.SourceId == "" {
		return shared.ResearchAnnotationAddResponse{}, bladWskazaniaBadan("annotation.add bez źródła")
	}
	if z.Kind == "" {
		return shared.ResearchAnnotationAddResponse{}, bladWskazaniaBadan("annotation.add bez rodzaju")
	}
	kod := nowyIdentyfikator(przedrostekAdnotacjiBadania)
	if z.AnnotationId != nil && *z.AnnotationId != "" {
		kod = *z.AnnotationId
	}
	zapisana, err := a.repozytorium.ZapiszAdnotacje(ctx, dane.AdnotacjaBadania{
		Kod: kod, ZrodloKod: z.SourceId, Rodzaj: string(z.Kind),
		Cytat: z.Quote, Komentarz: z.Comment, Kolor: z.Color,
		Kotwica: kotwicaDoBazyBadania(z.Anchor),
	})
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchAnnotationAddResponse{}, bladNieznanegoZrodlaBadania(z.SourceId)
	}
	if err != nil {
		return shared.ResearchAnnotationAddResponse{}, bladBadan(err)
	}
	return shared.ResearchAnnotationAddResponse{Annotation: zlozAdnotacjeBadania(zapisana)}, nil
}

// WypiszAdnotacje obsługuje `research.annotation.list`: wykaz adnotacji
// i wypisów zapisanych przy źródle.
func (a *adapterBadan) WypiszAdnotacje(ctx context.Context,
	z shared.ResearchAnnotationListRequest) (shared.ResearchAnnotationListResponse, error) {

	if (z.SourceId == nil || *z.SourceId == "") && (z.WindowId == nil || *z.WindowId == "") {
		return shared.ResearchAnnotationListResponse{}, bladWskazaniaBadan(
			"annotation.list bez źródła i bez okna — wskaż, czego dotyczy przegląd")
	}
	rodzaj := ""
	if z.Kind != nil {
		rodzaj = string(*z.Kind)
	}
	limit := 0
	if z.Limit != nil {
		limit = *z.Limit
	}
	adnotacje, err := a.repozytorium.Adnotacje(ctx, wartoscTekstu(z.SourceId),
		wartoscTekstu(z.WindowId), rodzaj, limit)
	if err != nil {
		return shared.ResearchAnnotationListResponse{}, bladBadan(err)
	}
	przelozone := make([]shared.ResearchAnnotation, 0, len(adnotacje))
	for _, adnotacja := range adnotacje {
		przelozone = append(przelozone, zlozAdnotacjeBadania(adnotacja))
	}
	return shared.ResearchAnnotationListResponse{Annotations: przelozone}, nil
}

// UsunAdnotacje obsługuje `research.annotation.remove`: usuwa adnotację albo
// wypis wskazany identyfikatorem.
func (a *adapterBadan) UsunAdnotacje(ctx context.Context,
	z shared.ResearchAnnotationRemoveRequest) (shared.ResearchAnnotationRemoveResponse, error) {

	if z.AnnotationId == "" {
		return shared.ResearchAnnotationRemoveResponse{}, bladWskazaniaBadan("annotation.remove bez adnotacji")
	}
	err := a.repozytorium.UsunAdnotacje(ctx, z.AnnotationId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchAnnotationRemoveResponse{}, protokolBladBadania(
			shared.ErrorCodeNotFound, "adnotacja nie istnieje: "+z.AnnotationId)
	}
	if err != nil {
		return shared.ResearchAnnotationRemoveResponse{}, bladBadan(err)
	}
	return shared.ResearchAnnotationRemoveResponse{AnnotationId: z.AnnotationId}, nil
}

// WypiszWypisy obsługuje `research.excerpt.list`. Wypis jest adnotacją niosącą
// cytat — zakładka bez cytatu wypisem nie jest i na listę nie wchodzi.
func (a *adapterBadan) WypiszWypisy(ctx context.Context,
	z shared.ResearchExcerptListRequest) (shared.ResearchExcerptListResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchExcerptListResponse{}, bladWskazaniaBadan("excerpt.list bez okna badania")
	}
	limit := 0
	if z.Limit != nil {
		limit = *z.Limit
	}
	adnotacje, err := a.repozytorium.Adnotacje(ctx, "", z.WindowId, "", limit)
	if err != nil {
		return shared.ResearchExcerptListResponse{}, bladBadan(err)
	}

	wypisy := []shared.ResearchExcerpt{}
	for _, adnotacja := range adnotacje {
		if adnotacja.Cytat == nil || strings.TrimSpace(*adnotacja.Cytat) == "" {
			continue
		}
		if len(z.SourceIds) > 0 && !zawieraNapisBadania(z.SourceIds, adnotacja.ZrodloKod) {
			continue
		}
		if z.Tag != nil && strings.TrimSpace(*z.Tag) != "" {
			katalog, err := a.repozytorium.KatalogZrodlaBadania(ctx, adnotacja.ZrodloKod)
			if err != nil {
				return shared.ResearchExcerptListResponse{}, bladBadan(err)
			}
			if !zawieraNapisBadania(katalog.Etykiety, *z.Tag) {
				continue
			}
		}
		wypisy = append(wypisy, wypisZAdnotacjiBadania(ctx, a, adnotacja))
	}
	return shared.ResearchExcerptListResponse{Excerpts: wypisy, Total: len(wypisy)}, nil
}

// wypisZAdnotacjiBadania składa wypis kontraktu z adnotacji wraz z etykietami
// źródła — po nich Findings Panel grupuje wypisy tematycznie.
func wypisZAdnotacjiBadania(ctx context.Context, a *adapterBadan,
	adnotacja dane.AdnotacjaBadania) shared.ResearchExcerpt {

	katalog, _ := a.repozytorium.KatalogZrodlaBadania(ctx, adnotacja.ZrodloKod)
	cytat := ""
	if adnotacja.Cytat != nil {
		cytat = *adnotacja.Cytat
	}
	return shared.ResearchExcerpt{
		AnnotationId: adnotacja.Kod, SourceId: adnotacja.ZrodloKod,
		SourceTitle: adnotacja.ZrodloTytul, Quote: cytat,
		Anchor: kotwicaZBazyBadania(adnotacja.Kotwica), Tags: katalog.Etykiety,
	}
}

// zlozAdnotacjeBadania przekłada wiersz repozytorium na byt adnotacji
// zwracany kontraktem komunikacji, wraz z odczytaną kotwicą.
func zlozAdnotacjeBadania(a dane.AdnotacjaBadania) shared.ResearchAnnotation {
	return shared.ResearchAnnotation{
		Id: a.Kod, SourceId: a.ZrodloKod, Kind: shared.ResearchAnnotationKind(a.Rodzaj),
		Quote: a.Cytat, Comment: a.Komentarz, Color: a.Kolor,
		Anchor: kotwicaZBazyBadania(a.Kotwica), FindingId: a.UstalenieKod,
		CreatedAt: chwilaBazy(a.Utworzono),
	}
}

// kotwicaDoBazyBadania przekłada kotwicę kontraktu na postać zapisywaną
// w wierszu bazy danych, jako dokument JSON.
func kotwicaDoBazyBadania(k shared.ResearchAnchor) dane.KotwicaBadania {
	kotwica := dane.KotwicaBadania{Rodzaj: string(k.Kind), Selektor: k.Selector, CzasMs: k.TimestampMs}
	if k.Page != nil {
		strona := int64(*k.Page)
		kotwica.Strona = &strona
	}
	if k.OffsetStart != nil {
		od := int64(*k.OffsetStart)
		kotwica.Od = &od
	}
	if k.OffsetEnd != nil {
		do := int64(*k.OffsetEnd)
		kotwica.Do = &do
	}
	return kotwica
}

// kotwicaZBazyBadania przekłada zapisany dokument JSON kotwicy z wiersza
// bazy danych z powrotem na byt kotwicy kontraktu komunikacji.
func kotwicaZBazyBadania(k dane.KotwicaBadania) shared.ResearchAnchor {
	rodzaj := k.Rodzaj
	if rodzaj == "" {
		rodzaj = string(shared.ResearchAnchorKindPage)
	}
	kotwica := shared.ResearchAnchor{
		Kind: shared.ResearchAnchorKind(rodzaj), Selector: k.Selektor, TimestampMs: k.CzasMs,
	}
	if k.Strona != nil {
		strona := int(*k.Strona)
		kotwica.Page = &strona
	}
	if k.Od != nil {
		od := int(*k.Od)
		kotwica.OffsetStart = &od
	}
	if k.Do != nil {
		do := int(*k.Do)
		kotwica.OffsetEnd = &do
	}
	return kotwica
}

// ── Operacje na treści źródła ──────────────────────────────────────────────

// StreszczZrodlo obsługuje `research.source.summarize`. Streszczenie powstaje
// modelem z treści źródła i zostaje zapisane przy źródle — druga prośba
// o streszczenie tej samej pozycji nie pali kolejnego wywołania modelu na darmo.
func (a *adapterBadan) StreszczZrodlo(ctx context.Context,
	z shared.ResearchSourceSummarizeRequest) (shared.ResearchSourceSummarizeResponse, error) {

	if z.SourceId == "" {
		return shared.ResearchSourceSummarizeResponse{}, bladWskazaniaBadan("source.summarize bez źródła")
	}
	tekst, _, err := a.trescDoLekturyBadania(ctx, z.SourceId, false)
	if err != nil {
		return shared.ResearchSourceSummarizeResponse{}, err
	}
	kanal := wartoscTekstu(z.ModelChannelId)
	if kanal == "" {
		kanal, err = a.domyslnyKanalBadania()
		if err != nil {
			return shared.ResearchSourceSummarizeResponse{}, err
		}
	}
	odpowiedz, err := a.zapytajModel(ctx, z.SourceId, kanal, polecenieStreszczeniaZrodlaBadania(tekst))
	if err != nil {
		return shared.ResearchSourceSummarizeResponse{}, err
	}

	streszczenie := streszczenieZOdpowiedziBadania(odpowiedz)
	if err := a.repozytorium.ZapiszStreszczenieZrodla(ctx, z.SourceId, streszczenie); err != nil {
		return shared.ResearchSourceSummarizeResponse{}, bladBadan(err)
	}
	return shared.ResearchSourceSummarizeResponse{
		Summary: shared.ResearchSummary{
			SourceId: z.SourceId, Abstract: streszczenie.Abstrakt, Theses: streszczenie.Tezy,
			Methodology: streszczenie.Metodologia, Conclusions: streszczenie.Wnioski,
		},
		FromModel: true,
	}, nil
}

// polecenieStreszczeniaZrodlaBadania wiąże model treścią źródła i żąda
// odpowiedzi w postaci, którą da się rozłożyć na pola kontraktu.
func polecenieStreszczeniaZrodlaBadania(tekst string) string {
	return "Streść poniższe źródło badawcze. Odpowiedz WYŁĄCZNIE obiektem JSON o polach: " +
		"\"abstrakt\" (jedno zdanie), \"tezy\" (lista zdań), \"metodologia\" (jedno zdanie albo pusty napis), " +
		"\"wnioski\" (lista zdań). Opieraj się wyłącznie na podanej treści — nie dodawaj faktów spoza niej.\n\n" +
		"Treść źródła:\n" + skrocDoBadania(tekst, 24000)
}

// streszczenieZOdpowiedziBadania rozkłada odpowiedź modelu na pola streszczenia.
// Odpowiedź, której nie da się rozłożyć, nie przepada: wchodzi w całości jako
// abstrakt roboczy, bo model powiedział coś o źródle, tylko innym kształtem.
func streszczenieZOdpowiedziBadania(odpowiedz string) dane.StreszczenieZrodlaBadania {
	var rozlozone struct {
		Abstrakt    string   `json:"abstrakt"`
		Tezy        []string `json:"tezy"`
		Metodologia string   `json:"metodologia"`
		Wnioski     []string `json:"wnioski"`
	}
	if err := json.Unmarshal([]byte(wnetrzeJsonBadania(odpowiedz)), &rozlozone); err != nil {
		calosc := strings.TrimSpace(odpowiedz)
		return dane.StreszczenieZrodlaBadania{Abstrakt: &calosc, Tezy: []string{}, Wnioski: []string{}}
	}
	streszczenie := dane.StreszczenieZrodlaBadania{Tezy: rozlozone.Tezy, Wnioski: rozlozone.Wnioski}
	if strings.TrimSpace(rozlozone.Abstrakt) != "" {
		abstrakt := rozlozone.Abstrakt
		streszczenie.Abstrakt = &abstrakt
	}
	if strings.TrimSpace(rozlozone.Metodologia) != "" {
		metodologia := rozlozone.Metodologia
		streszczenie.Metodologia = &metodologia
	}
	if streszczenie.Tezy == nil {
		streszczenie.Tezy = []string{}
	}
	if streszczenie.Wnioski == nil {
		streszczenie.Wnioski = []string{}
	}
	return streszczenie
}

// WyodrebnijTabele obsługuje `research.source.extractTable`: wykrywa tabele
// w tekście źródła heurystyką siatki znaków, do decyzji o ich utrwaleniu.
func (a *adapterBadan) WyodrebnijTabele(ctx context.Context,
	z shared.ResearchSourceExtractTableRequest) (shared.ResearchSourceExtractTableResponse, error) {

	if z.SourceId == "" {
		return shared.ResearchSourceExtractTableResponse{},
			bladWskazaniaBadan("source.extractTable bez źródła")
	}
	tekst, _, err := a.trescDoLekturyBadania(ctx, z.SourceId, false)
	if err != nil {
		return shared.ResearchSourceExtractTableResponse{}, err
	}
	tabele := tabeleZTekstuBadania(tekst)
	if z.TableIndex != nil {
		if *z.TableIndex < 0 || *z.TableIndex >= len(tabele) {
			return shared.ResearchSourceExtractTableResponse{},
				protokolBladBadania(shared.ErrorCodeNotFound,
					"w treści źródła "+z.SourceId+" nie ma tabeli o numerze "+
						strconv.Itoa(*z.TableIndex))
		}
		tabele = tabele[*z.TableIndex : *z.TableIndex+1]
	}

	przelozone := make([]shared.ResearchTable, 0, len(tabele))
	doZapisu := make([]dane.TabelaZrodlaBadania, 0, len(tabele))
	for _, tabela := range tabele {
		wiersze, err := json.Marshal(tabela.wiersze)
		if err != nil {
			return shared.ResearchSourceExtractTableResponse{}, bladBadan(err)
		}
		kod := nowyIdentyfikator(przedrostekTabeliBadania)
		przelozone = append(przelozone, shared.ResearchTable{
			Id: kod, SourceId: z.SourceId, Headers: tabela.naglowki, Rows: wiersze,
		})
		doZapisu = append(doZapisu, dane.TabelaZrodlaBadania{
			Kod: kod, ZrodloKod: z.SourceId, Naglowki: tabela.naglowki, Wiersze: string(wiersze),
		})
	}
	if z.Persist != nil && *z.Persist {
		if err := a.repozytorium.ZapiszTabeleZrodla(ctx, doZapisu); err != nil {
			return shared.ResearchSourceExtractTableResponse{}, bladBadan(err)
		}
	}
	return shared.ResearchSourceExtractTableResponse{Tables: przelozone}, nil
}

// tabelaTekstowaBadania jest jedną tabelą wykrytą w tekście źródła, wraz
// z zakresem wierszy tekstu, z których ją wykryto.
type tabelaTekstowaBadania struct {
	naglowki []string
	wiersze  [][]string
}

// tabeleZTekstuBadania wykrywa tabele w tekście po powtarzalnym rozkładzie
// separatorów. Wiersz osamotniony tabelą nie jest — potrzeba co najmniej
// nagłówka i jednego wiersza danych, inaczej każdy adres z myślnikami byłby tabelą.
func tabeleZTekstuBadania(tekst string) []tabelaTekstowaBadania {
	tabele := []tabelaTekstowaBadania{}
	biezaca := [][]string{}
	kolumn := 0

	domkniecie := func() {
		if len(biezaca) >= 2 {
			tabele = append(tabele, tabelaTekstowaBadania{
				naglowki: biezaca[0], wiersze: biezaca[1:],
			})
		}
		biezaca = [][]string{}
		kolumn = 0
	}

	for _, wiersz := range strings.Split(tekst, "\n") {
		komorki := komorkiWierszaBadania(wiersz)
		if len(komorki) < 2 {
			domkniecie()
			continue
		}
		if kolumn != 0 && len(komorki) != kolumn {
			domkniecie()
		}
		kolumn = len(komorki)
		biezaca = append(biezaca, komorki)
	}
	domkniecie()
	return tabele
}

// komorkiWierszaBadania rozkłada wiersz na komórki po separatorze pionowym,
// tabulatorze albo po co najmniej dwóch spacjach.
func komorkiWierszaBadania(wiersz string) []string {
	czysty := strings.TrimSpace(wiersz)
	if czysty == "" {
		return nil
	}
	var surowe []string
	switch {
	case strings.Contains(czysty, "|"):
		surowe = strings.Split(strings.Trim(czysty, "|"), "|")
	case strings.Contains(czysty, "\t"):
		surowe = strings.Split(czysty, "\t")
	default:
		surowe = strings.Split(czysty, "  ")
	}
	komorki := []string{}
	for _, komorka := range surowe {
		komorka = strings.TrimSpace(komorka)
		if komorka != "" {
			komorki = append(komorki, komorka)
		}
	}
	return komorki
}

// WyodrebnijTwierdzenia obsługuje `research.source.extractClaims`: prosi
// model o twierdzenia weryfikowalne wydobyte z treści źródła.
func (a *adapterBadan) WyodrebnijTwierdzenia(ctx context.Context,
	z shared.ResearchSourceExtractClaimsRequest) (shared.ResearchSourceExtractClaimsResponse, error) {

	if z.SourceId == "" {
		return shared.ResearchSourceExtractClaimsResponse{},
			bladWskazaniaBadan("source.extractClaims bez źródła")
	}
	tekst, _, err := a.trescDoLekturyBadania(ctx, z.SourceId, false)
	if err != nil {
		return shared.ResearchSourceExtractClaimsResponse{}, err
	}
	kanal, err := a.domyslnyKanalBadania()
	if err != nil {
		return shared.ResearchSourceExtractClaimsResponse{}, err
	}
	rodzaje := z.Kinds
	if len(rodzaje) == 0 {
		rodzaje = []string{"fakt", "dana liczbowa", "podmiot"}
	}
	polecenie := "Wydobądź z poniższego źródła kluczowe twierdzenia rodzajów: " +
		strings.Join(rodzaje, ", ") + ". Odpowiedz WYŁĄCZNIE tablicą JSON obiektów o polach " +
		"\"tresc\" i \"rodzaj\". Cytuj wyłącznie to, co w źródle stoi.\n\nTreść źródła:\n" +
		skrocDoBadania(tekst, 24000)

	odpowiedz, err := a.zapytajModel(ctx, z.SourceId, kanal, polecenie)
	if err != nil {
		return shared.ResearchSourceExtractClaimsResponse{}, err
	}
	var rozlozone []struct {
		Tresc  string `json:"tresc"`
		Rodzaj string `json:"rodzaj"`
	}
	if err := json.Unmarshal([]byte(wnetrzeJsonBadania(odpowiedz)), &rozlozone); err != nil {
		return shared.ResearchSourceExtractClaimsResponse{},
			protokolBladBadania(shared.ErrorCodeInternalError,
				"model nie oddał twierdzeń w postaci, którą da się rozłożyć: "+err.Error())
	}
	twierdzenia := make([]shared.ResearchClaim, 0, len(rozlozone))
	for _, wpis := range rozlozone {
		if strings.TrimSpace(wpis.Tresc) == "" {
			continue
		}
		twierdzenia = append(twierdzenia, shared.ResearchClaim{
			Text: wpis.Tresc, Kind: wpis.Rodzaj,
		})
	}
	return shared.ResearchSourceExtractClaimsResponse{Claims: twierdzenia}, nil
}

// RozpoznajPismoZrodla obsługuje `research.source.ocr`. Skutkiem jest treść
// źródła zapisana w bazie — po nim `reading.open` i `corpus.ask` widzą skan tak
// samo jak dokument tekstowy.
func (a *adapterBadan) RozpoznajPismoZrodla(ctx context.Context,
	z shared.ResearchSourceOcrRequest) (shared.ResearchSourceOcrResponse, error) {

	if z.SourceId == "" {
		return shared.ResearchSourceOcrResponse{}, bladWskazaniaBadan("source.ocr bez źródła")
	}
	if _, _, err := a.trescDoLekturyBadania(ctx, z.SourceId, true); err != nil {
		return shared.ResearchSourceOcrResponse{}, err
	}
	tresc, err := a.repozytorium.TrescZrodlaBadania(ctx, z.SourceId)
	if err != nil {
		return shared.ResearchSourceOcrResponse{}, bladBadan(err)
	}
	stron := 1
	if tresc.Stron != nil && *tresc.Stron > 0 {
		stron = int(*tresc.Stron)
	}
	rozpoznane := tresc.Tekst != nil && strings.TrimSpace(*tresc.Tekst) != ""
	return shared.ResearchSourceOcrResponse{PagesProcessed: stron, Indexed: rozpoznane}, nil
}

// ── Rozmowa oparta na korpusie ─────────────────────────────────────────────

// ZapytajKorpus obsługuje `research.corpus.ask`: odpowiada modelem
// wyłącznie z treści źródeł wskazanego korpusu, zakotwiczoną w nich.
func (a *adapterBadan) ZapytajKorpus(ctx context.Context,
	z shared.ResearchCorpusAskRequest) (shared.ResearchCorpusAskResponse, error) {

	if z.WindowId == "" {
		return shared.ResearchCorpusAskResponse{}, bladWskazaniaBadan("corpus.ask bez okna badania")
	}
	if strings.TrimSpace(z.Question) == "" {
		return shared.ResearchCorpusAskResponse{}, bladWskazaniaBadan("corpus.ask bez pytania")
	}
	zrodla, err := a.repozytorium.Zrodla(ctx, z.WindowId)
	if err != nil {
		return shared.ResearchCorpusAskResponse{}, bladBadan(err)
	}

	limit := 5
	if z.Limit != nil && *z.Limit > 0 {
		limit = *z.Limit
	}
	trafienia := a.trafieniaKorpusuBadania(ctx, zrodla, z.SourceIds, z.Question, limit)
	wymagaKotwic := z.RequireGrounding == nil || *z.RequireGrounding
	if len(trafienia) == 0 {
		if wymagaKotwic {
			return shared.ResearchCorpusAskResponse{}, protokolBladBadania(shared.ErrorCodeNotFound,
				"w korpusie badania nie ma treści pasującej do pytania — odpowiedź bez "+
					"zakotwiczenia w źródle jest w module Research odmówiona; naprawa: "+
					"dołożyć źródła albo wyłączyć wymóg zakotwiczenia (requireGrounding)")
		}
		return shared.ResearchCorpusAskResponse{Answer: shared.ResearchCorpusAnswer{
			Text: "", Citations: []shared.ResearchExcerpt{}, Grounded: false,
		}}, nil
	}

	kanal, err := a.domyslnyKanalBadania()
	if err != nil {
		return shared.ResearchCorpusAskResponse{}, err
	}
	odpowiedz, err := a.zapytajModel(ctx, z.WindowId, kanal,
		polecenieKorpusuBadania(z.Question, trafienia))
	if err != nil {
		return shared.ResearchCorpusAskResponse{}, err
	}

	cytaty := make([]shared.ResearchExcerpt, 0, len(trafienia))
	for _, trafienie := range trafienia {
		cytaty = append(cytaty, shared.ResearchExcerpt{
			AnnotationId: "", SourceId: trafienie.kodZrodla, SourceTitle: trafienie.tytul,
			Quote:  trafienie.fragment,
			Anchor: shared.ResearchAnchor{Kind: shared.ResearchAnchorKindPage},
		})
	}
	return shared.ResearchCorpusAskResponse{Answer: shared.ResearchCorpusAnswer{
		Text: odpowiedz, Citations: cytaty, Grounded: true,
	}}, nil
}

// trafienieKorpusuBadania jest jednym fragmentem korpusu dopasowanym do
// pytania, wraz ze źródłem, z którego fragment pochodzi.
type trafienieKorpusuBadania struct {
	kodZrodla string
	tytul     string
	fragment  string
	punkty    int
}

// trafieniaKorpusuBadania wybiera fragmenty źródeł najbliższe pytaniu miarą
// pokrycia słów, bez magazynu wektorów podobieństwa semantycznego.
func (a *adapterBadan) trafieniaKorpusuBadania(ctx context.Context, zrodla []dane.ZrodloBadania,
	wskazane []string, pytanie string, limit int) []trafienieKorpusuBadania {

	slowa := zbiorSlowBadania(pytanie)
	trafienia := []trafienieKorpusuBadania{}
	for _, zrodlo := range zrodla {
		if len(wskazane) > 0 && !zawieraNapisBadania(wskazane, zrodlo.Kod) {
			continue
		}
		tresc, err := a.repozytorium.TrescZrodlaBadania(ctx, zrodlo.Kod)
		if err != nil || tresc.Tekst == nil {
			continue
		}
		for _, akapit := range akapityBadania(*tresc.Tekst) {
			punkty := 0
			wAkapicie := zbiorSlowBadania(akapit)
			for slowo := range slowa {
				if wAkapicie[slowo] {
					punkty++
				}
			}
			if punkty == 0 {
				continue
			}
			trafienia = append(trafienia, trafienieKorpusuBadania{
				kodZrodla: zrodlo.Kod, tytul: zrodlo.Tytul,
				fragment: skrocDoBadania(akapit, 1200), punkty: punkty,
			})
		}
	}
	sort.SliceStable(trafienia, func(i, j int) bool {
		return trafienia[i].punkty > trafienia[j].punkty
	})
	if len(trafienia) > limit {
		trafienia = trafienia[:limit]
	}
	return trafienia
}

// akapityBadania rozkłada treść źródła na akapity nadające się na cytat
// w odpowiedzi opartej na korpusie, pomijając akapity zbyt krótkie.
func akapityBadania(tekst string) []string {
	akapity := []string{}
	for _, kawalek := range strings.Split(tekst, "\n") {
		kawalek = strings.TrimSpace(kawalek)
		if len(kawalek) >= 40 {
			akapity = append(akapity, kawalek)
		}
	}
	return akapity
}

// polecenieKorpusuBadania wiąże model wyłącznie podanymi fragmentami korpusu,
// zabraniając mu odpowiedzi z wiedzy spoza treści źródeł.
func polecenieKorpusuBadania(pytanie string, trafienia []trafienieKorpusuBadania) string {
	var polecenie strings.Builder
	polecenie.WriteString("Odpowiedz na pytanie badawcze WYŁĄCZNIE na podstawie poniższych ")
	polecenie.WriteString("fragmentów źródeł. Jeżeli fragmenty nie wystarczają, napisz to wprost ")
	polecenie.WriteString("zamiast uzupełniać odpowiedź wiedzą spoza nich. Wskazuj źródła numerem ")
	polecenie.WriteString("w nawiasie kwadratowym.\n\nPytanie: ")
	polecenie.WriteString(pytanie)
	polecenie.WriteString("\n\nFragmenty:\n")
	for numer, trafienie := range trafienia {
		polecenie.WriteString("[")
		polecenie.WriteString(strconv.Itoa(numer + 1))
		polecenie.WriteString("] ")
		polecenie.WriteString(trafienie.tytul)
		polecenie.WriteString(": ")
		polecenie.WriteString(trafienie.fragment)
		polecenie.WriteString("\n")
	}
	return polecenie.String()
}

// ── Pomocniki wspólne ──────────────────────────────────────────────────────

// skrocDoBadania przycina tekst do wskazanej liczby znaków. Model ma okno
// kontekstu, a treść źródła bywa od niego większa — przycięcie jest uczciwsze
// niż odmowa, bo streszczenie początku źródła to nadal streszczenie źródła.
func skrocDoBadania(tekst string, granica int) string {
	znaki := []rune(tekst)
	if len(znaki) <= granica {
		return tekst
	}
	return string(znaki[:granica]) + "\n[…treść przycięta do " + strconv.Itoa(granica) + " znaków…]"
}

// jsonModeluBadania rozkłada odpowiedź modelu na wskazaną strukturę, wyławiając
// wpierw sam JSON z tego, co model dopisał wokół niego.
func jsonModeluBadania(odpowiedz string, cel any) error {
	return json.Unmarshal([]byte(wnetrzeJsonBadania(odpowiedz)), cel)
}

// wnetrzeJsonBadania wyławia obiekt albo tablicę JSON z odpowiedzi modelu.
// Model bywa uprzejmy i dokłada zdanie wprowadzające albo ogrodzenie ```json —
// rozkładanie surowej odpowiedzi wywracałoby się na tej uprzejmości.
func wnetrzeJsonBadania(odpowiedz string) string {
	czysta := strings.TrimSpace(odpowiedz)
	czysta = strings.TrimPrefix(czysta, "```json")
	czysta = strings.TrimPrefix(czysta, "```")
	czysta = strings.TrimSuffix(czysta, "```")
	czysta = strings.TrimSpace(czysta)

	poczatek := strings.IndexAny(czysta, "[{")
	if poczatek < 0 {
		return czysta
	}
	koniec := strings.LastIndexAny(czysta, "]}")
	if koniec < poczatek {
		return czysta
	}
	return czysta[poczatek : koniec+1]
}
