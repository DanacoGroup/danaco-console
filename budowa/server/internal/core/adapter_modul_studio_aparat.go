// Aparat dokumentu Studio: spis treści, spisy ilustracji i tabel, przypisy,
// podpisy, bibliografia z powołaniami, odwołania wzajemne, odsyłacze, indeks
// i zakładki, przez cztery czynności kontraktu — założenie, wykaz, odświeżenie
// i usunięcie.
package core

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// WstawElementAparatu zakłada element aparatu dokumentu
// (`studio.apparatus.insert`); numer i zebrane pozycje dostaje dopiero po
// odświeżeniu.
func (a *adapterStudia) WstawElementAparatu(ctx context.Context,
	z shared.StudioApparatusInsertRequest) (shared.StudioApparatusInsertResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioApparatusInsertResponse{}, err
	}
	autor := postacAutor(z.Author)
	if err := aparatSprawdzRodzaj(z.Kind); err != nil {
		return shared.StudioApparatusInsertResponse{}, err
	}
	if err := aparatSprawdzWymagania(z); err != nil {
		return shared.StudioApparatusInsertResponse{}, err
	}

	dlugosc := postacDlugosc(&stan.forma)
	od, do := aparatZakotwiczenie(z, dlugosc)
	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor)
	if len(odcinki) == 0 {
		return shared.StudioApparatusInsertResponse{}, bladWskazaniaStudio(
			"założenie elementu aparatu " + postacZapisZakresu(od, do) +
				" zatrzymane przez blokadę fragmentu: " + postacNazwaBlokad(pominiete))
	}

	// Cel odwołania wzajemnego i odsyłacza sprawdza się przed zapisem, inaczej byłby odsyłaczem w nikąd.
	if z.TargetId != nil && strings.TrimSpace(*z.TargetId) != "" {
		if !aparatCelIstnieje(&stan.forma, *z.TargetId) {
			return shared.StudioApparatusInsertResponse{}, bladWskazaniaStudio(
				"odwołanie do elementu „" + *z.TargetId + "”, którego dokument nie ma; " +
					"celem odwołania jest zakładka, podpis, tabela albo obiekt tego dokumentu")
		}
	}

	element := shared.StudioApparatusItem{
		Kind: z.Kind, Text: z.Text, Label: z.Label, TargetId: z.TargetId,
		TargetUrl: z.TargetUrl, LevelsFrom: z.LevelsFrom, LevelsTo: z.LevelsTo,
		CitationKey: z.CitationKey, SourceTitle: z.SourceTitle,
		SourceAuthor: z.SourceAuthor, SourceYear: z.SourceYear, SourceUrl: z.SourceUrl,
	}
	if aparatRodzajZbierany(z.Kind) {
		if element.LevelsFrom == nil {
			element.LevelsFrom = postacWskaznikLiczby(1)
		}
		if element.LevelsTo == nil {
			element.LevelsTo = postacWskaznikLiczby(3)
		}
	}
	wiersz, err := aparatDoWiersza(stan.dokument.ID, nowyIdentyfikator(przedrostekAparatuPostaci),
		element, od, do, true)
	if err != nil {
		return shared.StudioApparatusInsertResponse{}, err
	}
	skladnica, err := a.postacSkladnica()
	if err != nil {
		return shared.StudioApparatusInsertResponse{}, err
	}
	zapisany, err := skladnica.ZapiszElementAparatu(ctx, wiersz)
	if err != nil {
		return shared.StudioApparatusInsertResponse{}, bladStudio(err)
	}

	// Wstawienie przypisu znaczy jako nieświeże przypisy tego rodzaju i spisy, które go obejmują.
	if err := a.aparatZnaczNieswiezoscPokrewnych(ctx, stan, z.Kind); err != nil {
		return shared.StudioApparatusInsertResponse{}, err
	}
	if err := a.postacWczytajWiersze(ctx, stan); err != nil {
		return shared.StudioApparatusInsertResponse{}, err
	}

	bilans := shared.StudioActionBalance{
		Applied: 1, Skipped: pominiete,
		Note: postacWskaznikTekstu(aparatNazwaRodzaju(z.Kind) + " założony " +
			postacZapisZakresu(od, do) + "; element wymaga odświeżenia, żeby dostać numer " +
			"i zebrać pozycje (studio.apparatus.refresh)"),
	}
	stan.opisCzynnosci = "założenie elementu aparatu: " + aparatNazwaRodzaju(z.Kind)

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindWstawienie, shared.StudioActionKindApparatusChange, od, do, bilans)
	if err != nil {
		return shared.StudioApparatusInsertResponse{}, err
	}
	zlozony := postacZlozAparat([]dane.ElementAparatuStudia{zapisany})
	if len(zlozony) == 0 {
		return shared.StudioApparatusInsertResponse{}, postacBladZaplecza(
			"element aparatu zapisany, ale nie da się go złożyć do odpowiedzi")
	}
	return shared.StudioApparatusInsertResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Item: zlozony[0],
	}, nil
}

// WykazAparatu oddaje aparat dokumentu (`studio.apparatus.list`), z możliwością
// zawężenia do rodzaju albo do elementów nieświeżych.
func (a *adapterStudia) WykazAparatu(ctx context.Context,
	z shared.StudioApparatusListRequest) (shared.StudioApparatusListResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioApparatusListResponse{}, err
	}
	elementy := make([]shared.StudioApparatusItem, 0, len(stan.forma.Apparatus))
	for _, element := range stan.forma.Apparatus {
		if z.Kind != nil && element.Kind != *z.Kind {
			continue
		}
		if z.StaleOnly != nil && *z.StaleOnly && (element.Stale == nil || !*element.Stale) {
			continue
		}
		elementy = append(elementy, element)
	}
	return shared.StudioApparatusListResponse{Items: elementy}, nil
}

// UsunElementAparatuDokumentu usuwa element aparatu (`studio.apparatus.remove`).
//
// Przypisy pozostałe przenumerowują się: usunięcie przypisu drugiego z trzech
// ma zostawić przypisy o numerach jeden i dwa, a nie jeden i trzy.
func (a *adapterStudia) UsunElementAparatuDokumentu(ctx context.Context,
	z shared.StudioApparatusRemoveRequest) (shared.StudioApparatusRemoveResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioApparatusRemoveResponse{}, err
	}
	autor := postacAutor(z.Author)
	kod := strings.TrimSpace(z.ItemId)
	if kod == "" {
		return shared.StudioApparatusRemoveResponse{}, bladWskazaniaStudio(
			"usunięcie elementu aparatu bez wskazania elementu")
	}
	usuwany, jest := aparatZnajdz(&stan.forma, kod)
	if !jest {
		return shared.StudioApparatusRemoveResponse{}, bladWskazaniaStudio(
			"usunięcie elementu aparatu „" + kod + "”, którego dokument nie ma")
	}
	od, do := 0, 0
	if usuwany.AnchorStart != nil {
		od = *usuwany.AnchorStart
	}
	if usuwany.AnchorEnd != nil {
		do = *usuwany.AnchorEnd
	}
	if odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor); len(odcinki) == 0 {
		return shared.StudioApparatusRemoveResponse{}, bladWskazaniaStudio(
			"usunięcie elementu aparatu zatrzymane przez blokadę fragmentu: " +
				postacNazwaBlokad(pominiete))
	}

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	// Element, do którego prowadzą odwołania, nie schodzi cicho — rdzeń nazywa je w bilansie.
	zalezne := aparatOdwolaniaDoCelu(&stan.forma, kod)
	skladnica, err := a.postacSkladnica()
	if err != nil {
		return shared.StudioApparatusRemoveResponse{}, err
	}
	for _, zalezny := range zalezne {
		bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
			Reason: "odwołanie straciło cel",
			Detail: postacWskaznikTekstu(aparatNazwaRodzaju(zalezny.Kind) + " „" + zalezny.Id +
				"” prowadził do usuwanego elementu; odwołanie zostało oznaczone jako " +
				"nieświeże i wymaga wskazania nowego celu"),
		})
		if err := a.aparatZapiszNieswiezosc(ctx, stan, zalezny.Id, true); err != nil {
			return shared.StudioApparatusRemoveResponse{}, err
		}
	}

	usunieto, err := skladnica.UsunElementAparatu(ctx, kod)
	if err != nil {
		return shared.StudioApparatusRemoveResponse{}, bladStudio(err)
	}
	if !usunieto {
		return shared.StudioApparatusRemoveResponse{}, bladWskazaniaStudio(
			"elementu aparatu „" + kod + "” nie udało się usunąć — wiersza już nie ma")
	}
	// Wiersze wczytują się przed znaczeniem nieświeżości — inaczej zapis znacznika nie znajdzie wiersza.
	if err := a.postacWczytajWiersze(ctx, stan); err != nil {
		return shared.StudioApparatusRemoveResponse{}, err
	}
	if err := a.aparatZnaczNieswiezoscPokrewnych(ctx, stan, usuwany.Kind); err != nil {
		return shared.StudioApparatusRemoveResponse{}, err
	}

	bilans.Applied = 1
	bilans.Note = postacWskaznikTekstu(aparatNazwaRodzaju(usuwany.Kind) + " usunięty; " +
		"pozostałe elementy tego rodzaju wymagają odświeżenia, żeby przenumerować się " +
		"(studio.apparatus.refresh)")
	stan.opisCzynnosci = "usunięcie elementu aparatu: " + aparatNazwaRodzaju(usuwany.Kind)

	forma, bilansGotowy, _, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindUsuniecie, shared.StudioActionKindApparatusChange, od, do, bilans)
	if err != nil {
		return shared.StudioApparatusRemoveResponse{}, err
	}
	return shared.StudioApparatusRemoveResponse{
		Removed: true, Form: forma, Balance: bilansGotowy, ActionId: stan.czynnosc,
	}, nil
}

// OdswiezAparat przelicza aparat dokumentu (`studio.apparatus.refresh`)
// jednym przebiegiem: spis treści, przypisy, spisy ilustracji i tabel, indeks
// i bibliografię, bo te rzeczy zależą od siebie.
func (a *adapterStudia) OdswiezAparat(ctx context.Context,
	z shared.StudioApparatusRefreshRequest) (shared.StudioApparatusRefreshResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioApparatusRefreshResponse{}, err
	}
	autor := postacAutor(z.Author)
	if z.Kind != nil {
		if err := aparatSprawdzRodzaj(*z.Kind); err != nil {
			return shared.StudioApparatusRefreshResponse{}, err
		}
	}
	if len(stan.forma.Apparatus) == 0 {
		return shared.StudioApparatusRefreshResponse{}, bladWskazaniaStudio(
			"odświeżenie aparatu dokumentu, który nie ma ani jednego elementu aparatu — " +
				"spis treści, przypis albo indeks trzeba najpierw założyć " +
				"(studio.apparatus.insert)")
	}

	strony, stron, err := aparatStronyAkapitow(postacTekstFormy(&stan.forma), &stan.forma)
	if err != nil {
		return shared.StudioApparatusRefreshResponse{}, err
	}
	obliczone := aparatPrzeliczCalosc(&stan.forma, strony, stron)

	skladnica, err := a.postacSkladnica()
	if err != nil {
		return shared.StudioApparatusRefreshResponse{}, err
	}
	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	odswiezone := make([]shared.StudioApparatusItem, 0, len(obliczone))
	for _, element := range obliczone {
		if z.Kind != nil && element.Kind != *z.Kind {
			continue
		}
		if z.ItemId != nil && strings.TrimSpace(*z.ItemId) != "" &&
			element.Id != strings.TrimSpace(*z.ItemId) {
			continue
		}
		if !aparatRodzajOdswiezalny(element.Kind) {
			bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
				Reason: "element bez czego odświeżać",
				Detail: postacWskaznikTekstu(aparatNazwaRodzaju(element.Kind) + " „" +
					element.Id + "” nie liczy się z treści dokumentu: jego brzmienie " +
					"podaje Operator i odświeżenie nie miałoby czego przeliczyć"),
			})
			odswiezone = append(odswiezone, element)
			continue
		}
		wiersz, err := aparatDoWiersza(stan.dokument.ID, element.Id, element,
			aparatWartoscLiczby(element.AnchorStart), aparatWartoscLiczby(element.AnchorEnd), false)
		if err != nil {
			return shared.StudioApparatusRefreshResponse{}, err
		}
		zapisany, err := skladnica.ZapiszElementAparatu(ctx, wiersz)
		if err != nil {
			return shared.StudioApparatusRefreshResponse{}, bladStudio(err)
		}
		zlozony := postacZlozAparat([]dane.ElementAparatuStudia{zapisany})
		if len(zlozony) > 0 {
			odswiezone = append(odswiezone, zlozony[0])
		}
		bilans.Applied++
	}
	if bilans.Applied == 0 && len(bilans.Skipped) == 0 {
		return shared.StudioApparatusRefreshResponse{}, bladWskazaniaStudio(
			"odświeżenie aparatu nie objęło ani jednego elementu — zawężenie rodzajem " +
				"albo wskazaniem elementu nie pasuje do niczego w dokumencie")
	}
	if err := a.postacWczytajWiersze(ctx, stan); err != nil {
		return shared.StudioApparatusRefreshResponse{}, err
	}
	bilans.Note = postacWskaznikTekstu("aparat odświeżony; elementów przeliczonych " +
		strconv.Itoa(bilans.Applied) + ", stron dokumentu " + strconv.Itoa(stron))
	stan.opisCzynnosci = "odświeżenie aparatu dokumentu"

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindApparatusChange,
		0, postacDlugosc(&stan.forma), bilans)
	if err != nil {
		return shared.StudioApparatusRefreshResponse{}, err
	}
	return shared.StudioApparatusRefreshResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Items: odswiezone,
	}, nil
}

// ── Rachunek odświeżenia ────────────────────────────────────────────────────

// aparatPrzeliczCalosc liczy aparat dokumentu od nowa i oddaje elementy
// z nadanymi numerami i zebranymi pozycjami.
func aparatPrzeliczCalosc(forma *shared.StudioDocumentForm, strony []int,
	stron int) []shared.StudioApparatusItem {

	elementy := append([]shared.StudioApparatusItem(nil), forma.Apparatus...)
	sort.SliceStable(elementy, func(i, j int) bool {
		return aparatWartoscLiczby(elementy[i].AnchorStart) <
			aparatWartoscLiczby(elementy[j].AnchorStart)
	})

	// Krok pierwszy: numeracja tego, co numeruje się z miejsca w treści.
	liczniki := map[shared.StudioApparatusKind]int{}
	for i := range elementy {
		rodzaj := elementy[i].Kind
		if !aparatRodzajNumerowany(rodzaj) {
			continue
		}
		liczniki[rodzaj]++
		elementy[i].Number = postacWskaznikTekstu(strconv.Itoa(liczniki[rodzaj]))
		elementy[i].Stale = postacWskaznikPrawdy(false)
	}

	// Krok drugi: powołanie dostaje numer z bibliografii, by cytat i pozycja wykazu były tym źródłem.
	numeryZrodel := aparatNumeryZrodel(elementy)
	for i := range elementy {
		if elementy[i].Kind != shared.StudioApparatusKindCitation {
			continue
		}
		klucz := aparatKluczZrodla(elementy[i])
		if numer, jest := numeryZrodel[klucz]; jest {
			elementy[i].Number = postacWskaznikTekstu(strconv.Itoa(numer))
		}
		elementy[i].Stale = postacWskaznikPrawdy(false)
	}

	// Krok trzeci: spisy i indeks zbierają pozycje z gotowej już numeracji.
	for i := range elementy {
		switch elementy[i].Kind {
		case shared.StudioApparatusKindToc:
			elementy[i].Entries = aparatZbierzSpisTresci(forma, strony,
				aparatWartoscLiczby(elementy[i].LevelsFrom),
				aparatWartoscLiczby(elementy[i].LevelsTo))
		case shared.StudioApparatusKindFigureIndex:
			elementy[i].Entries = aparatZbierzSpisIlustracji(forma, elementy, strony)
		case shared.StudioApparatusKindTableIndex:
			elementy[i].Entries = aparatZbierzSpisTabel(forma, elementy, strony)
		case shared.StudioApparatusKindIndex:
			elementy[i].Entries = aparatZbierzIndeks(forma, elementy, strony)
		case shared.StudioApparatusKindBibliography:
			elementy[i].Entries = aparatZbierzBibliografie(elementy, numeryZrodel)
		case shared.StudioApparatusKindCrossReference:
			elementy[i].Label = aparatEtykietaOdwolania(elementy, forma, elementy[i])
		default:
			continue
		}
		elementy[i].Stale = postacWskaznikPrawdy(false)
	}
	_ = stron
	return elementy
}

// aparatZbierzSpisTresci zbiera spis treści z nagłówków dokumentu, biorąc
// poziom nagłówka z postaci skutecznej akapitu, nie z postaci własnej.
func aparatZbierzSpisTresci(forma *shared.StudioDocumentForm, strony []int,
	odPoziomu, doPoziomu int) []shared.StudioApparatusEntry {

	if odPoziomu <= 0 {
		odPoziomu = 1
	}
	if doPoziomu <= 0 {
		doPoziomu = 3
	}
	pozycje := make([]shared.StudioApparatusEntry, 0, 16)
	akapit := 0
	for _, blok := range forma.Blocks {
		if !postacBlokNiesieTekst(blok) {
			continue
		}
		wskazanie := akapit
		akapit++
		poziom := aparatPoziomNaglowka(forma, blok)
		if poziom < odPoziomu || poziom > doPoziomu {
			continue
		}
		tresc := strings.TrimSpace(postacTekstBloku(blok))
		if tresc == "" {
			continue
		}
		pozycja := shared.StudioApparatusEntry{
			Text: tresc, Level: postacWskaznikLiczby(poziom),
			TargetId: postacWskaznikTekstu(blok.Id),
		}
		if blok.RangeStart != nil {
			pozycja.AnchorOffset = postacWskaznikLiczby(*blok.RangeStart)
		}
		if wskazanie < len(strony) {
			pozycja.PageNumber = postacWskaznikLiczby(strony[wskazanie])
		}
		pozycje = append(pozycje, pozycja)
	}
	return pozycje
}

// aparatPoziomNaglowka oddaje poziom nagłówka bloku z jego postaci skutecznej;
// zero znaczy „nie nagłówek".
func aparatPoziomNaglowka(forma *shared.StudioDocumentForm, blok shared.StudioDocumentBlock) int {
	skuteczna := postacAkapitSkuteczny(forma, blok)
	if skuteczna.OutlineLevel != nil && *skuteczna.OutlineLevel > 0 {
		return *skuteczna.OutlineLevel
	}
	if blok.Kind == blokPostaciNaglowek {
		// Blok oznaczony jako nagłówek bez poziomu jest nagłówkiem poziomu pierwszego.
		return 1
	}
	return 0
}

// aparatZbierzSpisIlustracji zbiera spis ilustracji z obiektów dokumentu wraz
// z ich podpisami oraz z podpisów aparatu przypiętych do obiektu.
func aparatZbierzSpisIlustracji(forma *shared.StudioDocumentForm,
	elementy []shared.StudioApparatusItem, strony []int) []shared.StudioApparatusEntry {

	pozycje := make([]shared.StudioApparatusEntry, 0, len(forma.Objects))
	numer := 0
	for _, obiekt := range forma.Objects {
		if obiekt.Kind == shared.StudioObjectKindTextbox {
			// Pole tekstowe nie jest ilustracją i w spisie ilustracji nie ma czego robić.
			continue
		}
		podpis := wartoscTekstu(obiekt.Caption)
		if podpis == "" {
			podpis = wartoscTekstu(obiekt.AltText)
		}
		if podpis == "" {
			podpis = aparatNazwaRodzajuObiektu(obiekt.Kind)
		}
		numer++
		miejsce := aparatWartoscLiczby(obiekt.AnchorOffset)
		pozycja := shared.StudioApparatusEntry{
			Text:         "Ilustracja " + strconv.Itoa(numer) + ". " + podpis,
			TargetId:     postacWskaznikTekstu(obiekt.Id),
			AnchorOffset: postacWskaznikLiczby(miejsce),
			PageNumber:   postacWskaznikLiczby(aparatStronaMiejsca(forma, strony, miejsce)),
		}
		pozycje = append(pozycje, pozycja)
	}
	// Podpisy założone jako element aparatu i przypięte do obiektu wchodzą do spisu ilustracji na równi.
	for _, element := range elementy {
		if element.Kind != shared.StudioApparatusKindCaption {
			continue
		}
		if element.TargetId == nil || !aparatCelObiektu(forma, *element.TargetId) {
			continue
		}
		miejsce := aparatWartoscLiczby(element.AnchorStart)
		pozycje = append(pozycje, shared.StudioApparatusEntry{
			Text:         strings.TrimSpace(wartoscTekstu(element.Number) + ". " + wartoscTekstu(element.Text)),
			TargetId:     postacWskaznikTekstu(element.Id),
			AnchorOffset: postacWskaznikLiczby(miejsce),
			PageNumber:   postacWskaznikLiczby(aparatStronaMiejsca(forma, strony, miejsce)),
		})
	}
	return pozycje
}

// aparatZbierzSpisTabel zbiera spis tabel z tabel dokumentu oraz z podpisów
// aparatu przypiętych do tabeli.
func aparatZbierzSpisTabel(forma *shared.StudioDocumentForm,
	elementy []shared.StudioApparatusItem, strony []int) []shared.StudioApparatusEntry {

	pozycje := make([]shared.StudioApparatusEntry, 0, len(forma.Tables))
	numer := 0
	for _, tabela := range forma.Tables {
		numer++
		podpis := wartoscTekstu(tabela.Caption)
		if podpis == "" {
			podpis = "tabela o " + strconv.Itoa(tabela.Rows) + " wierszach i " +
				strconv.Itoa(tabela.Columns) + " kolumnach"
		}
		miejsce := tabelaMiejsceTabeli(forma, tabela.Id)
		pozycje = append(pozycje, shared.StudioApparatusEntry{
			Text:         "Tabela " + strconv.Itoa(numer) + ". " + podpis,
			TargetId:     postacWskaznikTekstu(tabela.Id),
			AnchorOffset: postacWskaznikLiczby(miejsce),
			PageNumber:   postacWskaznikLiczby(aparatStronaMiejsca(forma, strony, miejsce)),
		})
	}
	for _, element := range elementy {
		if element.Kind != shared.StudioApparatusKindCaption {
			continue
		}
		if element.TargetId == nil || !aparatCelTabeli(forma, *element.TargetId) {
			continue
		}
		miejsce := aparatWartoscLiczby(element.AnchorStart)
		pozycje = append(pozycje, shared.StudioApparatusEntry{
			Text:         strings.TrimSpace(wartoscTekstu(element.Number) + ". " + wartoscTekstu(element.Text)),
			TargetId:     postacWskaznikTekstu(element.Id),
			AnchorOffset: postacWskaznikLiczby(miejsce),
			PageNumber:   postacWskaznikLiczby(aparatStronaMiejsca(forma, strony, miejsce)),
		})
	}
	return pozycje
}

// aparatZbierzIndeks składa indeks z haseł, zbierając strony jednego hasła
// w jedną pozycję uporządkowaną alfabetycznie.
func aparatZbierzIndeks(forma *shared.StudioDocumentForm,
	elementy []shared.StudioApparatusItem, strony []int) []shared.StudioApparatusEntry {

	kolejnosc := make([]string, 0, len(elementy))
	miejsca := map[string][]int{}
	kotwice := map[string]int{}
	for _, element := range elementy {
		if element.Kind != shared.StudioApparatusKindIndexEntry {
			continue
		}
		haslo := strings.TrimSpace(wartoscTekstu(element.Text))
		if haslo == "" {
			haslo = strings.TrimSpace(wartoscTekstu(element.Label))
		}
		if haslo == "" {
			continue
		}
		if _, jest := miejsca[haslo]; !jest {
			kolejnosc = append(kolejnosc, haslo)
			kotwice[haslo] = aparatWartoscLiczby(element.AnchorStart)
		}
		miejsca[haslo] = append(miejsca[haslo],
			aparatStronaMiejsca(forma, strony, aparatWartoscLiczby(element.AnchorStart)))
	}
	sort.SliceStable(kolejnosc, func(i, j int) bool {
		return strings.ToLower(kolejnosc[i]) < strings.ToLower(kolejnosc[j])
	})
	pozycje := make([]shared.StudioApparatusEntry, 0, len(kolejnosc))
	for _, haslo := range kolejnosc {
		numery := miejsca[haslo]
		sort.Ints(numery)
		zapisy := make([]string, 0, len(numery))
		poprzedni := -1
		for _, numer := range numery {
			if numer == poprzedni {
				continue
			}
			poprzedni = numer
			zapisy = append(zapisy, strconv.Itoa(numer))
		}
		pozycja := shared.StudioApparatusEntry{
			Text:         haslo + " — " + strings.Join(zapisy, ", "),
			AnchorOffset: postacWskaznikLiczby(kotwice[haslo]),
		}
		if len(numery) > 0 {
			pozycja.PageNumber = postacWskaznikLiczby(numery[0])
		}
		pozycje = append(pozycje, pozycja)
	}
	return pozycje
}

// aparatZbierzBibliografie składa wykaz źródeł z powołań dokumentu, w porządku
// numeracji nadanej wcześniej.
func aparatZbierzBibliografie(elementy []shared.StudioApparatusItem,
	numery map[string]int) []shared.StudioApparatusEntry {

	zebrane := map[string]shared.StudioApparatusItem{}
	for _, element := range elementy {
		if element.Kind != shared.StudioApparatusKindCitation {
			continue
		}
		klucz := aparatKluczZrodla(element)
		if klucz == "" {
			continue
		}
		if _, jest := zebrane[klucz]; !jest {
			zebrane[klucz] = element
		}
	}
	klucze := make([]string, 0, len(zebrane))
	for klucz := range zebrane {
		klucze = append(klucze, klucz)
	}
	sort.Slice(klucze, func(i, j int) bool { return numery[klucze[i]] < numery[klucze[j]] })

	pozycje := make([]shared.StudioApparatusEntry, 0, len(klucze))
	for _, klucz := range klucze {
		element := zebrane[klucz]
		pozycje = append(pozycje, shared.StudioApparatusEntry{
			Text:         "[" + strconv.Itoa(numery[klucz]) + "] " + aparatZapisZrodla(element),
			TargetId:     postacWskaznikTekstu(element.Id),
			AnchorOffset: element.AnchorStart,
		})
	}
	return pozycje
}

// aparatNumeryZrodel nadaje numery źródłom bibliograficznym w porządku
// pierwszego powołania w treści — tak numeruje się przypisy prawnicze i tak
// czyta się pismo od początku do końca.
func aparatNumeryZrodel(elementy []shared.StudioApparatusItem) map[string]int {
	numery := map[string]int{}
	nastepny := 0
	for _, element := range elementy {
		if element.Kind != shared.StudioApparatusKindCitation {
			continue
		}
		klucz := aparatKluczZrodla(element)
		if klucz == "" {
			continue
		}
		if _, jest := numery[klucz]; jest {
			continue
		}
		nastepny++
		numery[klucz] = nastepny
	}
	return numery
}

// aparatKluczZrodla rozstrzyga, co jest tym samym źródłem: klucz podany
// wprost albo złożenie autora, tytułu i roku.
func aparatKluczZrodla(element shared.StudioApparatusItem) string {
	if element.CitationKey != nil && strings.TrimSpace(*element.CitationKey) != "" {
		return strings.TrimSpace(*element.CitationKey)
	}
	czesci := []string{
		strings.TrimSpace(wartoscTekstu(element.SourceAuthor)),
		strings.TrimSpace(wartoscTekstu(element.SourceTitle)),
		strings.TrimSpace(wartoscTekstu(element.SourceYear)),
	}
	klucz := strings.Trim(strings.Join(czesci, "|"), "|")
	if klucz == "" {
		return strings.TrimSpace(wartoscTekstu(element.Text))
	}
	return klucz
}

// aparatZapisZrodla składa pozycję bibliografii z pól źródła: autora, tytułu,
// roku i adresu; brak wszystkich oddaje zdanie o braku opisu.
func aparatZapisZrodla(element shared.StudioApparatusItem) string {
	czesci := make([]string, 0, 4)
	for _, pole := range []*string{element.SourceAuthor, element.SourceTitle,
		element.SourceYear, element.SourceUrl} {

		if wartosc := strings.TrimSpace(wartoscTekstu(pole)); wartosc != "" {
			czesci = append(czesci, wartosc)
		}
	}
	if len(czesci) == 0 {
		if tresc := strings.TrimSpace(wartoscTekstu(element.Text)); tresc != "" {
			return tresc
		}
		return "źródło bez opisu — powołanie wymaga tytułu albo klucza"
	}
	return strings.Join(czesci, ", ")
}

// aparatEtykietaOdwolania składa brzmienie odwołania wzajemnego: „zobacz Tabela
// 2", „zobacz przypis 4". Odwołanie do celu bez numeru bierze jego etykietę.
func aparatEtykietaOdwolania(elementy []shared.StudioApparatusItem,
	forma *shared.StudioDocumentForm, odwolanie shared.StudioApparatusItem) *string {

	if odwolanie.TargetId == nil || strings.TrimSpace(*odwolanie.TargetId) == "" {
		return odwolanie.Label
	}
	cel := strings.TrimSpace(*odwolanie.TargetId)
	for _, element := range elementy {
		if element.Id != cel {
			continue
		}
		nazwa := aparatNazwaRodzaju(element.Kind)
		if element.Number != nil && *element.Number != "" {
			return postacWskaznikTekstu("zobacz " + nazwa + " " + *element.Number)
		}
		if element.Label != nil && *element.Label != "" {
			return postacWskaznikTekstu("zobacz " + *element.Label)
		}
		return postacWskaznikTekstu("zobacz " + nazwa)
	}
	for numer, tabela := range forma.Tables {
		if tabela.Id == cel {
			return postacWskaznikTekstu("zobacz Tabela " + strconv.Itoa(numer+1))
		}
	}
	for numer, obiekt := range forma.Objects {
		if obiekt.Id == cel {
			return postacWskaznikTekstu("zobacz Ilustracja " + strconv.Itoa(numer+1))
		}
	}
	return odwolanie.Label
}

// ── Strony ──────────────────────────────────────────────────────────────────

// aparatStronyAkapitow liczy, na której stronie stoi każdy akapit dokumentu,
// i ile stron dokument ma, tym samym silnikiem, którym jedzie podgląd wydruku.
func aparatStronyAkapitow(tresc string, forma *shared.StudioDocumentForm) ([]int, int, error) {
	oblicze, err := krojWyrysuStudia(rozmiarPismaStudia)
	if err != nil {
		return nil, 0, err
	}
	defer oblicze.Close()

	domyslna := geometriaZNastawStrony(aparatNastawyDokumentu(forma))
	akapity := strings.Split(tresc, "\n")
	strony := make([]int, 0, len(akapity))

	// wiersz liczy wiersze wypełnione na stronie bieżącej; pojemność strony zmienia się razem z sekcją.
	kartka := domyslna
	naStrone := kartka.wierszyNaStrone()
	strona, wiersz := 1, 0
	poprzedniaSekcja := ""
	poczatek := 0

	for _, akapit := range akapity {
		dlugosc := len([]rune(akapit))
		sekcja, jest := aparatSekcjaMiejsca(forma, poczatek)
		if jest {
			if sekcja.Id != poprzedniaSekcja {
				kartka = geometriaZNastawStrony(aparatNastawySekcji(sekcja, forma))
				naStrone = kartka.wierszyNaStrone()
				if aparatSekcjaLamieStrone(sekcja) && (wiersz > 0 || strona > 1) {
					strona++
					wiersz = 0
				}
				poprzedniaSekcja = sekcja.Id
			}
		} else if poprzedniaSekcja != "" {
			kartka = domyslna
			naStrone = kartka.wierszyNaStrone()
			poprzedniaSekcja = ""
		}

		strony = append(strony, strona)
		wiersz += len(zlamWierszStudia(akapit, oblicze, kartka.szerokoscKolumny()))
		for wiersz >= naStrone {
			wiersz -= naStrone
			strona++
		}
		// Granica akapitu jest jednym znakiem treści, zgodnie z zakresami, którymi liczy je reszta modułu.
		poczatek += dlugosc + 1
	}

	stron := strona
	if wiersz == 0 && stron > 1 {
		// Kartka pusta na końcu nie jest stroną dokumentu — ostatni akapit domknął stronę poprzednią.
		stron--
	}
	if stron < 1 {
		stron = 1
	}
	return strony, stron, nil
}

// aparatNastawyDokumentu oddaje nastawy strony dokumentu albo brak; aparatNastawySekcji
// sięga po nie, gdy sekcja własnych nie ma.
func aparatNastawyDokumentu(forma *shared.StudioDocumentForm) *shared.StudioPageSetup {
	if forma == nil {
		return nil
	}
	return forma.PageSetup
}

// aparatNastawySekcji oddaje nastawy strony obowiązujące w sekcji: własne, gdy
// sekcja je ma, a w przeciwnym razie nastawy dokumentu. Sekcja bez własnych
// nastaw nie jest sekcją o nastawach domyślnych — jest sekcją, która dziedziczy.
func aparatNastawySekcji(sekcja shared.StudioSection,
	forma *shared.StudioDocumentForm) *shared.StudioPageSetup {

	if sekcja.PageSetup != nil {
		return sekcja.PageSetup
	}
	return aparatNastawyDokumentu(forma)
}

// aparatSekcjaMiejsca oddaje sekcję, której zakres obejmuje wskazane miejsce
// treści. Sekcje zachodzące na siebie rozstrzyga pierwsza pasująca — wykaz idzie
// w kolejności zakresów, więc pierwsza jest tą, która się zaczęła najwcześniej.
func aparatSekcjaMiejsca(forma *shared.StudioDocumentForm,
	miejsce int) (shared.StudioSection, bool) {

	if forma == nil {
		return shared.StudioSection{}, false
	}
	for _, sekcja := range forma.Sections {
		if miejsce >= sekcja.RangeStart && miejsce < sekcja.RangeEnd {
			return sekcja, true
		}
	}
	return shared.StudioSection{}, false
}

// aparatSekcjaLamieStrone mówi, czy początek sekcji przerywa stronę: `newPage`,
// `evenPage` i `oddPage` przerywają, `continuous` i `newColumn` nie.
func aparatSekcjaLamieStrone(sekcja shared.StudioSection) bool {
	if sekcja.Start == nil {
		return false
	}
	switch *sekcja.Start {
	case shared.StudioSectionStartNewPage, shared.StudioSectionStartEvenPage,
		shared.StudioSectionStartOddPage:

		return true
	default:
		return false
	}
}

// aparatStronaMiejsca oddaje numer strony miejsca w treści liczonego w znakach,
// z rachunku stron akapitów.
func aparatStronaMiejsca(forma *shared.StudioDocumentForm, strony []int, miejsce int) int {
	akapit := 0
	for _, blok := range forma.Blocks {
		if !postacBlokNiesieTekst(blok) {
			continue
		}
		if blok.RangeEnd != nil && *blok.RangeEnd >= miejsce {
			break
		}
		akapit++
	}
	if akapit < len(strony) {
		return strony[akapit]
	}
	if len(strony) > 0 {
		return strony[len(strony)-1]
	}
	return 1
}

// ── Nieświeżość ─────────────────────────────────────────────────────────────

// aparatZnaczNieswiezoscPokrewnych zapala znacznik nieświeżości na elementach,
// które zmiana wskazanego rodzaju unieważnia.
func (a *adapterStudia) aparatZnaczNieswiezoscPokrewnych(ctx context.Context,
	stan *stanPostaci, rodzaj shared.StudioApparatusKind) error {

	for _, zalezny := range aparatRodzajeZalezne(rodzaj) {
		if err := a.aparatZnaczNieswiezoscRodzaju(ctx, stan, zalezny); err != nil {
			return err
		}
	}
	return nil
}

// aparatZnaczNieswiezoscRodzaju zapala znacznik nieświeżości na wszystkich
// elementach wskazanego rodzaju, także w obszarze tabel i wstawień.
func (a *adapterStudia) aparatZnaczNieswiezoscRodzaju(ctx context.Context,
	stan *stanPostaci, rodzaj shared.StudioApparatusKind) error {

	for _, element := range stan.forma.Apparatus {
		if element.Kind != rodzaj {
			continue
		}
		if err := a.aparatZapiszNieswiezosc(ctx, stan, element.Id, true); err != nil {
			return err
		}
	}
	return nil
}

// aparatZapiszNieswiezosc zapisuje znacznik nieświeżości elementu kolumną,
// nie polem drzewa — drzewo postaci zapisuje się bez aparatu.
func (a *adapterStudia) aparatZapiszNieswiezosc(ctx context.Context, stan *stanPostaci,
	kod string, nieswiezy bool) error {

	skladnica, err := a.postacSkladnica()
	if err != nil {
		return err
	}
	wiersz, err := skladnica.ElementAparatu(ctx, kod)
	if err != nil {
		return bladStudio(err)
	}
	if wiersz.Nieswiezy == nieswiezy {
		return nil
	}
	wiersz.Nieswiezy = nieswiezy
	if _, err := skladnica.ZapiszElementAparatu(ctx, wiersz); err != nil {
		return bladStudio(err)
	}
	for i := range stan.forma.Apparatus {
		if stan.forma.Apparatus[i].Id == kod {
			stan.forma.Apparatus[i].Stale = postacWskaznikPrawdy(nieswiezy)
		}
	}
	return nil
}

// ── Drobne rachunki ─────────────────────────────────────────────────────────

// aparatZakotwiczenie ustala zakotwiczenie elementu w treści: zakres wskazany
// wprost albo punkt wstawienia.
func aparatZakotwiczenie(z shared.StudioApparatusInsertRequest, dlugosc int) (int, int) {
	if z.RangeStart != nil || z.RangeEnd != nil {
		return postacZakres(z.RangeStart, z.RangeEnd, dlugosc)
	}
	miejsce := 0
	if z.Offset != nil {
		miejsce = *z.Offset
	}
	if miejsce < 0 {
		miejsce = 0
	}
	if miejsce > dlugosc {
		miejsce = dlugosc
	}
	return miejsce, miejsce
}

// aparatZnajdz odnajduje element aparatu w postaci po jego kodzie, zwracając
// też, czy element istnieje.
func aparatZnajdz(forma *shared.StudioDocumentForm, kod string) (shared.StudioApparatusItem, bool) {
	for _, element := range forma.Apparatus {
		if element.Id == kod {
			return element, true
		}
	}
	return shared.StudioApparatusItem{}, false
}

// aparatOdwolaniaDoCelu wymienia elementy prowadzące do wskazanego celu — te,
// które usunięcie celu unieważnia.
func aparatOdwolaniaDoCelu(forma *shared.StudioDocumentForm,
	kod string) []shared.StudioApparatusItem {

	zalezne := make([]shared.StudioApparatusItem, 0, 4)
	for _, element := range forma.Apparatus {
		if element.Id == kod || element.TargetId == nil {
			continue
		}
		if strings.TrimSpace(*element.TargetId) == kod {
			zalezne = append(zalezne, element)
		}
	}
	return zalezne
}

// aparatCelIstnieje mówi, czy dokument ma byt o wskazaniu podanym jako cel:
// element aparatu, blok, obiekt albo tabelę.
func aparatCelIstnieje(forma *shared.StudioDocumentForm, kod string) bool {
	szukany := strings.TrimSpace(kod)
	if szukany == "" {
		return false
	}
	for _, element := range forma.Apparatus {
		if element.Id == szukany {
			return true
		}
	}
	for _, blok := range forma.Blocks {
		if blok.Id == szukany {
			return true
		}
	}
	return aparatCelObiektu(forma, szukany) || aparatCelTabeli(forma, szukany)
}

// aparatCelObiektu mówi, czy wskazanie należy do obiektu dokumentu, jednego
// ze źródeł spisu ilustracji.
func aparatCelObiektu(forma *shared.StudioDocumentForm, kod string) bool {
	for _, obiekt := range forma.Objects {
		if obiekt.Id == strings.TrimSpace(kod) {
			return true
		}
	}
	return false
}

// aparatCelTabeli mówi, czy wskazanie należy do tabeli dokumentu, jednego ze
// źródeł spisu tabel obok podpisów aparatu.
func aparatCelTabeli(forma *shared.StudioDocumentForm, kod string) bool {
	for _, tabela := range forma.Tables {
		if tabela.Id == strings.TrimSpace(kod) {
			return true
		}
	}
	return false
}

// aparatWartoscLiczby odczytuje wskaźnik na liczbę, znosząc brak przez
// oddanie zera zamiast wskaźnika pustego.
func aparatWartoscLiczby(wskazanie *int) int {
	if wskazanie == nil {
		return 0
	}
	return *wskazanie
}

// aparatSprawdzRodzaj odrzuca rodzaj elementu aparatu, którego kontrakt nie
// zna, nazywając wykaz rodzajów znanych.
func aparatSprawdzRodzaj(rodzaj shared.StudioApparatusKind) error {
	for _, znany := range shared.WartosciStudioApparatusKind() {
		if rodzaj == znany {
			return nil
		}
	}
	nazwy := make([]string, 0, 13)
	for _, znany := range shared.WartosciStudioApparatusKind() {
		nazwy = append(nazwy, string(znany))
	}
	return bladWskazaniaStudio("element aparatu rodzaju „" + string(rodzaj) +
		"”, którego serwer nie zna; wykaz: " + strings.Join(nazwy, ", "))
}

// aparatSprawdzWymagania pilnuje, żeby element aparatu wchodził z tym, bez
// czego jest pusty.
//
// Element bez brzmienia i bez celu byłby dokładnie tym, czego zlecenie zakazuje:
// odpowiedzią `ok` z niczym w środku. Odmowa nazywa, czego brakuje.
func aparatSprawdzWymagania(z shared.StudioApparatusInsertRequest) error {
	pusteBrzmienie := z.Text == nil || strings.TrimSpace(*z.Text) == ""
	switch z.Kind {
	case shared.StudioApparatusKindFootnote, shared.StudioApparatusKindEndnote:
		if pusteBrzmienie {
			return bladWskazaniaStudio("przypis bez brzmienia — przypis wymaga treści (pole text)")
		}
	case shared.StudioApparatusKindCaption:
		if pusteBrzmienie {
			return bladWskazaniaStudio("podpis bez brzmienia — podpis wymaga treści (pole text)")
		}
	case shared.StudioApparatusKindIndexEntry:
		if pusteBrzmienie && (z.Label == nil || strings.TrimSpace(*z.Label) == "") {
			return bladWskazaniaStudio("hasło indeksu bez brzmienia — hasło wymaga treści " +
				"(pole text) albo etykiety (pole label)")
		}
	case shared.StudioApparatusKindBookmark:
		if z.Label == nil || strings.TrimSpace(*z.Label) == "" {
			return bladWskazaniaStudio("zakładka bez nazwy — zakładka wymaga etykiety " +
				"(pole label), inaczej Operator nie ma po czym do niej wrócić")
		}
	case shared.StudioApparatusKindHyperlink:
		if z.TargetUrl == nil || strings.TrimSpace(*z.TargetUrl) == "" {
			if z.TargetId == nil || strings.TrimSpace(*z.TargetId) == "" {
				return bladWskazaniaStudio("odsyłacz bez celu — odsyłacz wymaga adresu " +
					"(pole targetUrl) albo wskazania elementu dokumentu (pole targetId)")
			}
		}
	case shared.StudioApparatusKindCrossReference:
		if z.TargetId == nil || strings.TrimSpace(*z.TargetId) == "" {
			return bladWskazaniaStudio("odwołanie wzajemne bez celu — odwołanie wymaga " +
				"wskazania elementu dokumentu (pole targetId)")
		}
	case shared.StudioApparatusKindCitation:
		if (z.CitationKey == nil || strings.TrimSpace(*z.CitationKey) == "") &&
			(z.SourceTitle == nil || strings.TrimSpace(*z.SourceTitle) == "") {
			return bladWskazaniaStudio("powołanie bez źródła — powołanie wymaga klucza " +
				"(pole citationKey) albo tytułu źródła (pole sourceTitle)")
		}
	}
	return nil
}

// aparatRodzajZbierany mówi, czy element zbiera pozycje z dokumentu, jak spis
// treści, indeks albo bibliografia.
func aparatRodzajZbierany(rodzaj shared.StudioApparatusKind) bool {
	switch rodzaj {
	case shared.StudioApparatusKindToc, shared.StudioApparatusKindFigureIndex,
		shared.StudioApparatusKindTableIndex, shared.StudioApparatusKindIndex,
		shared.StudioApparatusKindBibliography:
		return true
	default:
		return false
	}
}

// aparatRodzajNumerowany mówi, czy element dostaje numer z porządku treści,
// jak przypis, podpis albo powołanie.
func aparatRodzajNumerowany(rodzaj shared.StudioApparatusKind) bool {
	switch rodzaj {
	case shared.StudioApparatusKindFootnote, shared.StudioApparatusKindEndnote,
		shared.StudioApparatusKindCaption:
		return true
	default:
		return false
	}
}

// aparatRodzajOdswiezalny mówi, czy odświeżenie ma co przeliczyć.
//
// Zakładka i odsyłacz do adresu nie liczą się z niczego — ich brzmienie podaje
// Operator. Rdzeń mówi to wprost w bilansie, zamiast udawać odświeżenie.
func aparatRodzajOdswiezalny(rodzaj shared.StudioApparatusKind) bool {
	switch rodzaj {
	case shared.StudioApparatusKindBookmark, shared.StudioApparatusKindHyperlink:
		return false
	default:
		return true
	}
}

// aparatRodzajeZalezne wymienia rodzaje, które zmiana wskazanego rodzaju
// unieważnia i każe znaczyć jako nieświeże.
func aparatRodzajeZalezne(rodzaj shared.StudioApparatusKind) []shared.StudioApparatusKind {
	switch rodzaj {
	case shared.StudioApparatusKindFootnote:
		return []shared.StudioApparatusKind{shared.StudioApparatusKindFootnote}
	case shared.StudioApparatusKindEndnote:
		return []shared.StudioApparatusKind{shared.StudioApparatusKindEndnote}
	case shared.StudioApparatusKindCaption:
		return []shared.StudioApparatusKind{shared.StudioApparatusKindCaption,
			shared.StudioApparatusKindFigureIndex, shared.StudioApparatusKindTableIndex}
	case shared.StudioApparatusKindIndexEntry:
		return []shared.StudioApparatusKind{shared.StudioApparatusKindIndex}
	case shared.StudioApparatusKindCitation:
		return []shared.StudioApparatusKind{shared.StudioApparatusKindCitation,
			shared.StudioApparatusKindBibliography}
	default:
		return nil
	}
}

// aparatNazwaRodzaju nazywa rodzaj elementu pełnym słowem — odmowa i bilans
// mówią Operatorowi, o co chodzi, bez zaglądania do kontraktu.
func aparatNazwaRodzaju(rodzaj shared.StudioApparatusKind) string {
	switch rodzaj {
	case shared.StudioApparatusKindToc:
		return "spis treści"
	case shared.StudioApparatusKindFigureIndex:
		return "spis ilustracji"
	case shared.StudioApparatusKindTableIndex:
		return "spis tabel"
	case shared.StudioApparatusKindFootnote:
		return "przypis dolny"
	case shared.StudioApparatusKindEndnote:
		return "przypis końcowy"
	case shared.StudioApparatusKindCaption:
		return "podpis"
	case shared.StudioApparatusKindBookmark:
		return "zakładka"
	case shared.StudioApparatusKindCrossReference:
		return "odwołanie wzajemne"
	case shared.StudioApparatusKindHyperlink:
		return "odsyłacz"
	case shared.StudioApparatusKindCitation:
		return "powołanie"
	case shared.StudioApparatusKindBibliography:
		return "bibliografia"
	case shared.StudioApparatusKindIndexEntry:
		return "hasło indeksu"
	case shared.StudioApparatusKindIndex:
		return "indeks"
	default:
		return "element aparatu"
	}
}

// aparatNazwaRodzajuObiektu nazywa rodzaj obiektu na potrzeby spisu ilustracji,
// gdy obiekt jest bez podpisu.
func aparatNazwaRodzajuObiektu(rodzaj shared.StudioObjectKind) string {
	switch rodzaj {
	case shared.StudioObjectKindImage:
		return "obraz bez podpisu"
	case shared.StudioObjectKindShape:
		return "kształt bez podpisu"
	case shared.StudioObjectKindIcon:
		return "ikona bez podpisu"
	case shared.StudioObjectKindLogo:
		return "logo bez podpisu"
	case shared.StudioObjectKindChart:
		return "wykres bez podpisu"
	default:
		return "obiekt bez podpisu"
	}
}

// aparatDoWiersza przekłada element aparatu na wiersz warstwy danych: pola
// z kolumną idą kolumnami, reszta polem JSON.
func aparatDoWiersza(dokumentID int64, kod string, element shared.StudioApparatusItem,
	od, do int, nieswiezy bool) (dane.ElementAparatuStudia, error) {

	dodatkowe := shared.StudioApparatusItem{
		LevelsFrom: element.LevelsFrom, LevelsTo: element.LevelsTo,
		Entries: element.Entries, CitationKey: element.CitationKey,
		SourceTitle: element.SourceTitle, SourceAuthor: element.SourceAuthor,
		SourceYear: element.SourceYear, SourceUrl: element.SourceUrl,
	}
	zapis, err := json.Marshal(dodatkowe)
	if err != nil {
		return dane.ElementAparatuStudia{}, postacBladZaplecza(
			"elementu aparatu nie da się zapisać: " + err.Error())
	}
	wiersz := dane.ElementAparatuStudia{
		Kod: kod, DokumentID: dokumentID, Rodzaj: string(element.Kind),
		KotwicaOd: int64(od), KotwicaDo: int64(do),
		Numer: element.Number, Etykieta: element.Label, Tresc: element.Text,
		CelKod: element.TargetId, CelAdres: element.TargetUrl,
		Kolejnosc: int64(od), Nieswiezy: nieswiezy,
		DaneJSON: postacWskaznikTekstu(string(zapis)),
	}
	return wiersz, nil
}
