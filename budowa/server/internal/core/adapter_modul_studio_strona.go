// Nastawy strony i sekcje dokumentu: format nośnika, orientacja, marginesy,
// kolumny, nagłówek i stopka, numeracja stron, znak wodny, podziały
// wstawiane w treść, nadruk koperty, wykaz formatów nośnika oraz tabulatory
// linijki.
package core

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"danacoconsole/shared"
)

// ── Nastawy strony ──────────────────────────────────────────────────────────

// NastawyStrony oddaje nastawy strony dokumentu albo sekcji — obsługuje
// studio.page.setup.get. Nastawy sekcji wychodzą już nałożone na nastawy
// dokumentu, bo Operator pyta o nośnik załącznika, nie o pola nadpisane.
func (a *adapterStudia) NastawyStrony(ctx context.Context,
	z shared.StudioPageSetupGetRequest) (shared.StudioPageSetupGetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioPageSetupGetResponse{}, err
	}
	nastawy := shared.StudioPageSetup{}
	if stan.forma.PageSetup != nil {
		nastawy = *stan.forma.PageSetup
	} else {
		nastawy = *postacNastawyDomyslne()
	}
	if z.SectionId != nil && strings.TrimSpace(*z.SectionId) != "" {
		sekcja := stronaSekcjaPoKodzie(&stan.forma, *z.SectionId)
		if sekcja == nil {
			return shared.StudioPageSetupGetResponse{}, bladWskazaniaStudio(
				"sekcji „" + strings.TrimSpace(*z.SectionId) + "” dokument " +
					stan.dokument.Kod + " nie ma; wykaz sekcji oddaje komenda studio.section.list")
		}
		nastawy = stronaNastawySkuteczne(&nastawy, sekcja.PageSetup)
	}
	return shared.StudioPageSetupGetResponse{
		PageSetup: nastawy,
		Sections:  stan.forma.Sections,
	}, nil
}

// stronaNastawySkuteczne nakłada nastawy sekcji na nastawy dokumentu. Pole, na
// które sekcja nic nie mówi, zostaje z dokumentu — to jest sens dziedziczenia
// i bez niego zmiana marginesu w dokumencie omijałaby sekcje.
func stronaNastawySkuteczne(dokumentu, sekcji *shared.StudioPageSetup) shared.StudioPageSetup {
	wynik := shared.StudioPageSetup{}
	if dokumentu != nil {
		wynik = *dokumentu
	}
	if sekcji == nil {
		return wynik
	}
	if sekcji.PageSize != nil {
		wynik.PageSize = sekcji.PageSize
		// Nazwa nośnika sekcji unieważnia wymiar własny odziedziczony po
		// dokumencie.
		wynik.WidthMm, wynik.HeightMm = sekcji.WidthMm, sekcji.HeightMm
	}
	if sekcji.WidthMm != nil {
		wynik.WidthMm = sekcji.WidthMm
	}
	if sekcji.HeightMm != nil {
		wynik.HeightMm = sekcji.HeightMm
	}
	if sekcji.PaperKind != nil {
		wynik.PaperKind = sekcji.PaperKind
	}
	if sekcji.Orientation != nil {
		wynik.Orientation = sekcji.Orientation
	}
	if sekcji.MarginTop != nil {
		wynik.MarginTop = sekcji.MarginTop
	}
	if sekcji.MarginBottom != nil {
		wynik.MarginBottom = sekcji.MarginBottom
	}
	if sekcji.MarginLeft != nil {
		wynik.MarginLeft = sekcji.MarginLeft
	}
	if sekcji.MarginRight != nil {
		wynik.MarginRight = sekcji.MarginRight
	}
	if sekcji.MarginPreset != nil {
		wynik.MarginPreset = sekcji.MarginPreset
	}
	if sekcji.GutterMm != nil {
		wynik.GutterMm = sekcji.GutterMm
	}
	if sekcji.MirrorMargins != nil {
		wynik.MirrorMargins = sekcji.MirrorMargins
	}
	if sekcji.Columns != nil {
		wynik.Columns = sekcji.Columns
	}
	if sekcji.ColumnGapMm != nil {
		wynik.ColumnGapMm = sekcji.ColumnGapMm
	}
	if sekcji.ColumnRule != nil {
		wynik.ColumnRule = sekcji.ColumnRule
	}
	if sekcji.Envelope != nil {
		wynik.Envelope = sekcji.Envelope
	}
	return wynik
}

// UstawNastawyStrony przestawia nośnik, orientację, marginesy i kolumny
// dokumentu albo pojedynczej sekcji — obsługuje studio.page.setup.set i
// oddaje bilans tego, co się nie zmieściło po zmianie nośnika.
func (a *adapterStudia) UstawNastawyStrony(ctx context.Context,
	z shared.StudioPageSetupSetRequest) (shared.StudioPageSetupSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioPageSetupSetResponse{}, err
	}
	autor := postacAutor(z.Author)

	naSekcji := z.SectionId != nil && strings.TrimSpace(*z.SectionId) != ""
	var sekcja *shared.StudioSection
	zastane := stan.forma.PageSetup
	if naSekcji {
		if sekcja, err = a.stronaSekcjaZadania(ctx, stan, z.SectionId); err != nil {
			return shared.StudioPageSetupSetResponse{}, err
		}
		zastane = sekcja.PageSetup
	}

	nastawy, zmian, err := stronaScalNastawy(zastane, z)
	if err != nil {
		return shared.StudioPageSetupSetResponse{}, err
	}
	if zmian == 0 {
		return shared.StudioPageSetupSetResponse{}, bladWskazaniaStudio(
			"nastawy strony bez ani jednej rzeczy do przestawienia — żądanie nie niosło " +
				"nośnika, orientacji, marginesów ani kolumn")
	}

	// Bilans liczy się na nastawie skutecznej, bo to ona rozstrzyga, ile
	// miejsca treść dostanie.
	skuteczne := nastawy
	if naSekcji {
		skuteczne = stronaNastawySkuteczne(stan.forma.PageSetup, &nastawy)
	}
	pominiete := stronaBilansUkladu(&stan.forma, &skuteczne)

	zakresOd, zakresDo := 0, postacDlugosc(&stan.forma)
	opis := "nastawy strony dokumentu przestawione"
	if naSekcji {
		sekcja.PageSetup = &nastawy
		if err := a.stronaZapiszSekcje(ctx, stan.dokument.ID, *sekcja); err != nil {
			return shared.StudioPageSetupSetResponse{}, err
		}
		zakresOd, zakresDo = sekcja.RangeStart, sekcja.RangeEnd
		opis = "nastawy strony sekcji " + sekcja.Id + " przestawione"
	} else {
		stan.forma.PageSetup = &nastawy
	}
	stan.opisCzynnosci = opis

	bilans := shared.StudioActionBalance{
		Applied: zmian,
		Skipped: pominiete,
		Note:    postacWskaznikTekstu(opis + ": " + stronaOpisNosnika(&skuteczne)),
	}
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindPageChange,
		zakresOd, zakresDo, bilans)
	if err != nil {
		return shared.StudioPageSetupSetResponse{}, err
	}
	return shared.StudioPageSetupSetResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, PageSetup: skuteczne,
	}, nil
}

// stronaOpisNosnika opisuje nośnik zdaniem, które Operator czyta bez
// zaglądania do pola po polu, składając nazwę, kierunek, wymiary i
// szerokość kolumny tekstu.
func stronaOpisNosnika(nastawy *shared.StudioPageSetup) string {
	nazwa := "własny"
	if nastawy != nil && nastawy.PageSize != nil && strings.TrimSpace(*nastawy.PageSize) != "" {
		nazwa = *nastawy.PageSize
	}
	szerokosc, wysokosc := stronaWymiary(nastawy)
	kierunek := string(shared.StudioPageOrientationPionowa)
	if nastawy != nil && nastawy.Orientation != nil {
		kierunek = string(*nastawy.Orientation)
	}
	return nazwa + " " + kierunek + ", " + stronaZapisMiary(szerokosc) + " na " +
		stronaZapisMiary(wysokosc) + " mm, kolumna tekstu " +
		stronaZapisMiary(stronaSzerokoscUzytkowa(nastawy)) + " mm"
}

// NosnikiStrony oddaje wykaz formatów nośnika — obsługuje
// studio.page.paper.list, czerpiąc ze wspólnego miejsca rdzenia, z którego
// bierze go też wykaz Designu. Zawężenie rodzajem nieznanym jest odmową
// nazwaną.
func (a *adapterStudia) NosnikiStrony(_ context.Context,
	z shared.StudioPagePaperListRequest) (shared.StudioPagePaperListResponse, error) {

	wykaz := stronaNosniki()
	if z.Kind == nil || strings.TrimSpace(string(*z.Kind)) == "" {
		return shared.StudioPagePaperListResponse{Papers: wykaz}, nil
	}
	rodzajZnany := false
	for _, wartosc := range shared.WartosciStudioPaperKind() {
		if wartosc == *z.Kind {
			rodzajZnany = true
			break
		}
	}
	if !rodzajZnany {
		nazwy := make([]string, 0, 3)
		for _, wartosc := range shared.WartosciStudioPaperKind() {
			nazwy = append(nazwy, string(wartosc))
		}
		return shared.StudioPagePaperListResponse{}, bladWskazaniaStudio(
			"rodzaju nośnika „" + string(*z.Kind) + "” kontrakt nie zna; rodzaje: " +
				strings.Join(nazwy, ", "))
	}
	wybrane := make([]shared.StudioPaperFormat, 0, len(wykaz))
	for _, nosnik := range wykaz {
		if nosnik.Kind == *z.Kind {
			wybrane = append(wybrane, nosnik)
		}
	}
	if len(wybrane) == 0 {
		return shared.StudioPagePaperListResponse{}, bladWskazaniaStudio(
			"wykaz nośników rdzenia nie niesie ani jednej pozycji rodzaju „" +
				string(*z.Kind) + "”. Rodzaj `custom` nie stoi w wykazie z zamysłu: " +
				"wymiar własny podaje się polami widthMm i heightMm komendy " +
				"studio.page.setup.set")
	}
	return shared.StudioPagePaperListResponse{Papers: wybrane}, nil
}

// ── Nagłówek, stopka, numeracja, znak wodny ─────────────────────────────────

// NaglowkiIStopki oddaje nagłówki i stopki wedle zasięgu — obsługuje
// studio.page.headerfooter.get, czytając je z sekcji wskazanej albo z
// sekcji pierwszej dokumentu.
func (a *adapterStudia) NaglowkiIStopki(ctx context.Context,
	z shared.StudioPageHeaderfooterGetRequest) (shared.StudioPageHeaderfooterGetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioPageHeaderfooterGetResponse{}, err
	}
	var wykaz []shared.StudioHeaderFooter
	if z.SectionId != nil && strings.TrimSpace(*z.SectionId) != "" {
		sekcja := stronaSekcjaPoKodzie(&stan.forma, *z.SectionId)
		if sekcja == nil {
			return shared.StudioPageHeaderfooterGetResponse{}, bladWskazaniaStudio(
				"sekcji „" + strings.TrimSpace(*z.SectionId) + "” dokument " +
					stan.dokument.Kod + " nie ma")
		}
		wykaz = sekcja.HeadersFooters
	} else if len(stan.forma.Sections) > 0 {
		sort.SliceStable(stan.forma.Sections, func(i, j int) bool {
			return stan.forma.Sections[i].Index < stan.forma.Sections[j].Index
		})
		wykaz = stan.forma.Sections[0].HeadersFooters
	}
	// Dokument sprzed tej dobudowy trzyma nagłówek w jednym polu nastaw
	// strony.
	if len(wykaz) == 0 {
		wykaz = stronaNaglowkiZNastaw(stan.forma.PageSetup)
	}
	if wykaz == nil {
		wykaz = []shared.StudioHeaderFooter{}
	}
	return shared.StudioPageHeaderfooterGetResponse{HeadersFooters: wykaz}, nil
}

// UstawNaglowekIStopke ustawia nagłówek i stopkę zasięgu — obsługuje
// studio.page.headerfooter.set. Strony zwykłe, pierwsza strona i strony
// parzyste są trzema osobnymi nagłówkami jednej sekcji.
func (a *adapterStudia) UstawNaglowekIStopke(ctx context.Context,
	z shared.StudioPageHeaderfooterSetRequest) (shared.StudioPageHeaderfooterSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioPageHeaderfooterSetResponse{}, err
	}
	zasieg, err := stronaZasiegNaglowka(z.Scope)
	if err != nil {
		return shared.StudioPageHeaderfooterSetResponse{}, err
	}
	autor := postacAutor(z.Author)
	sekcja, err := a.stronaSekcjaZadania(ctx, stan, z.SectionId)
	if err != nil {
		return shared.StudioPageHeaderfooterSetResponse{}, err
	}
	zastane := sekcja.HeadersFooters
	if len(zastane) == 0 {
		zastane = stronaNaglowkiZNastaw(stan.forma.PageSetup)
	}
	wykaz, zmian := stronaScalNaglowek(zastane, z, zasieg)
	if zmian == 0 {
		return shared.StudioPageHeaderfooterSetResponse{}, bladWskazaniaStudio(
			"nagłówek bez ani jednej rzeczy do ustawienia — żądanie nie niosło treści " +
				"nagłówka, treści stopki, odległości od krawędzi ani przejęcia z sekcji " +
				"poprzedniej")
	}
	sekcja.HeadersFooters = wykaz
	if err := a.stronaZapiszSekcje(ctx, stan.dokument.ID, *sekcja); err != nil {
		return shared.StudioPageHeaderfooterSetResponse{}, err
	}
	stan.opisCzynnosci = "nagłówek i stopka sekcji " + sekcja.Id + " zasięgu " +
		string(zasieg) + " ustawione"

	bilans := shared.StudioActionBalance{
		Applied: zmian,
		Skipped: []shared.StudioSkippedItem{},
		Note:    postacWskaznikTekstu(stan.opisCzynnosci),
	}
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindPageChange,
		sekcja.RangeStart, sekcja.RangeEnd, bilans)
	if err != nil {
		return shared.StudioPageHeaderfooterSetResponse{}, err
	}
	return shared.StudioPageHeaderfooterSetResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, HeadersFooters: wykaz,
	}, nil
}

// UstawNumeracjeStron ustawia numerację stron sekcji — obsługuje
// studio.page.numbering.set i przenosi przełącznik numeracji do nastaw
// strony na potrzeby podglądu wydruku.
func (a *adapterStudia) UstawNumeracjeStron(ctx context.Context,
	z shared.StudioPageNumberingSetRequest) (shared.StudioPageNumberingSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioPageNumberingSetResponse{}, err
	}
	autor := postacAutor(z.Author)
	sekcja, err := a.stronaSekcjaZadania(ctx, stan, z.SectionId)
	if err != nil {
		return shared.StudioPageNumberingSetResponse{}, err
	}
	numeracja, err := stronaScalNumeracje(sekcja.Numbering, z)
	if err != nil {
		return shared.StudioPageNumberingSetResponse{}, err
	}
	sekcja.Numbering = &numeracja
	if err := a.stronaZapiszSekcje(ctx, stan.dokument.ID, *sekcja); err != nil {
		return shared.StudioPageNumberingSetResponse{}, err
	}
	// Przełącznik numeracji w nastawach strony idzie razem z numeracją
	// sekcji.
	if stan.forma.PageSetup == nil {
		stan.forma.PageSetup = postacNastawyDomyslne()
	}
	stan.forma.PageSetup.PageNumbers = postacWskaznikPrawdy(numeracja.Enabled)

	stan.opisCzynnosci = "numeracja stron sekcji " + sekcja.Id + " ustawiona"
	bilans := shared.StudioActionBalance{
		Applied: 1,
		Skipped: []shared.StudioSkippedItem{},
		Note:    postacWskaznikTekstu(stan.opisCzynnosci),
	}
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindPageChange,
		sekcja.RangeStart, sekcja.RangeEnd, bilans)
	if err != nil {
		return shared.StudioPageNumberingSetResponse{}, err
	}
	return shared.StudioPageNumberingSetResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Numbering: numeracja,
	}, nil
}

// UstawZnakWodny ustawia znak wodny dokumentu albo sekcji — obsługuje
// studio.page.watermark.set, zapisując go w nastawach sekcji wskazanej
// żądaniem.
func (a *adapterStudia) UstawZnakWodny(ctx context.Context,
	z shared.StudioPageWatermarkSetRequest) (shared.StudioPageWatermarkSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioPageWatermarkSetResponse{}, err
	}
	autor := postacAutor(z.Author)
	sekcja, err := a.stronaSekcjaZadania(ctx, stan, z.SectionId)
	if err != nil {
		return shared.StudioPageWatermarkSetResponse{}, err
	}
	znak, err := stronaScalZnakWodny(sekcja.Watermark, z)
	if err != nil {
		return shared.StudioPageWatermarkSetResponse{}, err
	}
	sekcja.Watermark = &znak
	if err := a.stronaZapiszSekcje(ctx, stan.dokument.ID, *sekcja); err != nil {
		return shared.StudioPageWatermarkSetResponse{}, err
	}
	stan.opisCzynnosci = "znak wodny sekcji " + sekcja.Id + " ustawiony"
	bilans := shared.StudioActionBalance{
		Applied: 1,
		Skipped: []shared.StudioSkippedItem{},
		Note:    postacWskaznikTekstu(stan.opisCzynnosci),
	}
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindPageChange,
		sekcja.RangeStart, sekcja.RangeEnd, bilans)
	if err != nil {
		return shared.StudioPageWatermarkSetResponse{}, err
	}
	return shared.StudioPageWatermarkSetResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Watermark: znak,
	}, nil
}

// UstawNadrukKoperty ustawia nadruk koperty — obsługuje
// studio.page.envelope.set. Nadruk na nośniku niekopertowym wraca odmową
// nazwaną, a koperta bez nastaw adresata i nadawcy jest samym rozmiarem,
// nie funkcją.
func (a *adapterStudia) UstawNadrukKoperty(ctx context.Context,
	z shared.StudioPageEnvelopeSetRequest) (shared.StudioPageEnvelopeSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioPageEnvelopeSetResponse{}, err
	}
	autor := postacAutor(z.Author)

	naSekcji := z.SectionId != nil && strings.TrimSpace(*z.SectionId) != ""
	var sekcja *shared.StudioSection
	zastane := stan.forma.PageSetup
	if naSekcji {
		if sekcja, err = a.stronaSekcjaZadania(ctx, stan, z.SectionId); err != nil {
			return shared.StudioPageEnvelopeSetResponse{}, err
		}
		zastane = sekcja.PageSetup
	}
	skuteczne := shared.StudioPageSetup{}
	if naSekcji {
		skuteczne = stronaNastawySkuteczne(stan.forma.PageSetup, zastane)
	} else if zastane != nil {
		skuteczne = *zastane
	}
	if skuteczne.PaperKind == nil || *skuteczne.PaperKind != shared.StudioPaperKindEnvelope {
		nazwa := "arkusz"
		if skuteczne.PageSize != nil {
			nazwa = *skuteczne.PageSize
		}
		return shared.StudioPageEnvelopeSetResponse{}, bladWskazaniaStudio(
			"nadruk koperty na nośniku „" + nazwa + "”, który kopertą nie jest. " +
				"Najpierw ustaw nośnik kopertowy komendą studio.page.setup.set — " +
				"koperty w wykazie: DL, C4, C5, C6")
	}

	koperta := shared.StudioEnvelopeSetup{}
	if zastane != nil && zastane.Envelope != nil {
		koperta = *zastane.Envelope
	}
	zmian := 0
	if z.Recipient != nil {
		koperta.Recipient = postacWskaznikTekstu(*z.Recipient)
		zmian++
	}
	if z.Sender != nil {
		koperta.Sender = postacWskaznikTekstu(*z.Sender)
		zmian++
	}
	if z.IncludeSender != nil {
		koperta.IncludeSender = postacWskaznikPrawdy(*z.IncludeSender)
		zmian++
	}
	for _, para := range []struct {
		wartosc *float64
		zapis   **float64
	}{
		{z.RecipientXMm, &koperta.RecipientXMm},
		{z.RecipientYMm, &koperta.RecipientYMm},
		{z.SenderXMm, &koperta.SenderXMm},
		{z.SenderYMm, &koperta.SenderYMm},
	} {
		if para.wartosc == nil {
			continue
		}
		if *para.wartosc < 0 {
			return shared.StudioPageEnvelopeSetResponse{}, bladWskazaniaStudio(
				"położenie na kopercie liczy się od jej krawędzi i nie może być ujemne")
		}
		*para.zapis = postacWskaznikMiary(*para.wartosc)
		zmian++
	}
	if zmian == 0 {
		return shared.StudioPageEnvelopeSetResponse{}, bladWskazaniaStudio(
			"nadruk koperty bez ani jednej rzeczy do ustawienia — żądanie nie niosło " +
				"adresata, nadawcy ani ich położenia")
	}

	// Położenie niepodane dostaje nastawę z normy: adresat w prawej dolnej
	// ćwiartce.
	szerokosc, wysokosc := stronaWymiary(&skuteczne)
	if koperta.RecipientXMm == nil {
		koperta.RecipientXMm = postacWskaznikMiary(szerokosc * 0.45)
	}
	if koperta.RecipientYMm == nil {
		koperta.RecipientYMm = postacWskaznikMiary(wysokosc * 0.55)
	}
	if koperta.SenderXMm == nil {
		koperta.SenderXMm = postacWskaznikMiary(15)
	}
	if koperta.SenderYMm == nil {
		koperta.SenderYMm = postacWskaznikMiary(15)
	}
	if koperta.IncludeSender == nil {
		koperta.IncludeSender = postacWskaznikPrawdy(koperta.Sender != nil)
	}

	zakresOd, zakresDo := 0, postacDlugosc(&stan.forma)
	if naSekcji {
		if sekcja.PageSetup == nil {
			nastawy := shared.StudioPageSetup{}
			sekcja.PageSetup = &nastawy
		}
		sekcja.PageSetup.Envelope = &koperta
		if err := a.stronaZapiszSekcje(ctx, stan.dokument.ID, *sekcja); err != nil {
			return shared.StudioPageEnvelopeSetResponse{}, err
		}
		zakresOd, zakresDo = sekcja.RangeStart, sekcja.RangeEnd
		stan.opisCzynnosci = "nadruk koperty sekcji " + sekcja.Id + " ustawiony"
	} else {
		if stan.forma.PageSetup == nil {
			stan.forma.PageSetup = postacNastawyDomyslne()
		}
		stan.forma.PageSetup.Envelope = &koperta
		stan.opisCzynnosci = "nadruk koperty dokumentu ustawiony"
	}

	bilans := shared.StudioActionBalance{
		Applied: zmian,
		Skipped: []shared.StudioSkippedItem{},
		Note:    postacWskaznikTekstu(stan.opisCzynnosci),
	}
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindPageChange,
		zakresOd, zakresDo, bilans)
	if err != nil {
		return shared.StudioPageEnvelopeSetResponse{}, err
	}
	return shared.StudioPageEnvelopeSetResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Envelope: koperta,
	}, nil
}

// ── Podział wstawiany w treść ───────────────────────────────────────────────

// WstawPodzial wstawia podział strony, kolumny, sekcji albo wiersza —
// obsługuje studio.page.break.insert. Podział jest blokiem nietekstowym;
// podział sekcji zakłada nową sekcję od miejsca podziału do końca
// dokumentu.
func (a *adapterStudia) WstawPodzial(ctx context.Context,
	z shared.StudioPageBreakInsertRequest) (shared.StudioPageBreakInsertResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioPageBreakInsertResponse{}, err
	}
	rodzajZnany := false
	for _, wartosc := range shared.WartosciStudioBreakKind() {
		if wartosc == z.Kind {
			rodzajZnany = true
			break
		}
	}
	if !rodzajZnany {
		nazwy := make([]string, 0, 4)
		for _, wartosc := range shared.WartosciStudioBreakKind() {
			nazwy = append(nazwy, string(wartosc))
		}
		return shared.StudioPageBreakInsertResponse{}, bladWskazaniaStudio(
			"rodzaju podziału „" + string(z.Kind) + "” kontrakt nie zna; rodzaje: " +
				strings.Join(nazwy, ", "))
	}
	autor := postacAutor(z.Author)
	dlugosc := postacDlugosc(&stan.forma)
	miejsce, _ := postacZakres(&z.Offset, &z.Offset, dlugosc)

	// Blokada obowiązuje także tu: podział pod blokadą rozerwałby chroniony
	// fragment.
	if blokady := postacBlokadyZakresu(&stan.forma, miejsce, miejsce); len(blokady) > 0 {
		for _, blokada := range blokady {
			if blokada.Scope == shared.StudioLockScopeEveryone || autor == shared.StudioAuthorModel {
				return shared.StudioPageBreakInsertResponse{}, bladWskazaniaStudio(
					"podziału nie da się wstawić " + postacZapisZakresu(miejsce, miejsce) +
						": blokada „" + blokada.Name + "” nie przepuszcza tej czynności")
			}
		}
	}

	wskazanie, roznica := stronaRozdzielAkapit(&stan.forma, miejsce)
	rodzaj := z.Kind
	podzial := shared.StudioDocumentBlock{
		Id:        nowyIdentyfikator(przedrostekBlokuPostaci),
		Kind:      blokPostaciPodzial,
		BreakKind: &rodzaj,
	}
	bloki := make([]shared.StudioDocumentBlock, 0, len(stan.forma.Blocks)+1)
	bloki = append(bloki, stan.forma.Blocks[:wskazanie]...)
	bloki = append(bloki, podzial)
	bloki = append(bloki, stan.forma.Blocks[wskazanie:]...)
	stan.forma.Blocks = bloki
	if roznica != 0 {
		postacPrzesunZakotwiczenia(&stan.forma, miejsce, miejsce, roznica)
	}
	postacPrzeliczZakresy(&stan.forma)

	bilans := shared.StudioActionBalance{Applied: 1, Skipped: []shared.StudioSkippedItem{}}
	opis := "podział „" + string(z.Kind) + "” wstawiony " + postacZapisZakresu(miejsce, miejsce)

	if z.Kind == shared.StudioBreakKindSection {
		nowa, err := a.stronaPodzielSekcje(ctx, stan, miejsce, z.SectionStart)
		if err != nil {
			return shared.StudioPageBreakInsertResponse{}, err
		}
		bilans.Applied++
		opis += "; sekcja " + nowa + " założona od tego miejsca do końca dokumentu"
	}
	stan.opisCzynnosci = opis
	bilans.Note = postacWskaznikTekstu(opis)

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindWstawienie, shared.StudioActionKindPageChange,
		miejsce, miejsce, bilans)
	if err != nil {
		return shared.StudioPageBreakInsertResponse{}, err
	}
	return shared.StudioPageBreakInsertResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana, ActionId: stan.czynnosc,
	}, nil
}

// stronaRozdzielAkapit rozdziela akapit w miejscu podziału i oddaje
// wskazanie bloku, przed którym podział ma stanąć, oraz różnicę długości
// treści; podział na granicy akapitu niczego nie rozdziela.
func stronaRozdzielAkapit(forma *shared.StudioDocumentForm, miejsce int) (int, int) {
	postacRozetnij(forma, miejsce)
	for i := range forma.Blocks {
		blok := &forma.Blocks[i]
		if !postacBlokNiesieTekst(*blok) {
			continue
		}
		start, koniec := 0, 0
		if blok.RangeStart != nil {
			start = *blok.RangeStart
		}
		if blok.RangeEnd != nil {
			koniec = *blok.RangeEnd
		}
		if miejsce < start || miejsce > koniec {
			continue
		}
		if miejsce == start {
			return i, 0
		}
		if miejsce == koniec {
			return i + 1, 0
		}
		// Rozdzielenie: fragmenty do miejsca zostają, dalsze idą do bloku
		// nowego, kopią.
		zostaja := make([]shared.StudioDocumentRun, 0, len(blok.Runs))
		dalsze := make([]shared.StudioDocumentRun, 0, len(blok.Runs))
		for _, run := range blok.Runs {
			granica := 0
			if run.RangeEnd != nil {
				granica = *run.RangeEnd
			}
			if granica <= miejsce {
				zostaja = append(zostaja, run)
			} else {
				dalsze = append(dalsze, run)
			}
		}
		nowy := shared.StudioDocumentBlock{
			Id:        nowyIdentyfikator(przedrostekBlokuPostaci),
			Kind:      blok.Kind,
			SectionId: blok.SectionId,
			Paragraph: postacKopiaAkapitu(blok.Paragraph),
			Runs:      dalsze,
		}
		if len(zostaja) == 0 {
			zostaja = append(zostaja, shared.StudioDocumentRun{Text: ""})
		}
		if len(nowy.Runs) == 0 {
			nowy.Runs = []shared.StudioDocumentRun{{Text: ""}}
		}
		blok.Runs = zostaja
		bloki := make([]shared.StudioDocumentBlock, 0, len(forma.Blocks)+1)
		bloki = append(bloki, forma.Blocks[:i+1]...)
		bloki = append(bloki, nowy)
		bloki = append(bloki, forma.Blocks[i+1:]...)
		forma.Blocks = bloki
		postacPrzeliczZakresy(forma)
		return i + 1, 1
	}
	return len(forma.Blocks), 0
}

// stronaPodzielSekcje zakłada sekcję od miejsca podziału do końca dokumentu
// i skraca sekcję, w której to miejsce leży.
func (a *adapterStudia) stronaPodzielSekcje(ctx context.Context, stan *stanPostaci,
	miejsce int, rozpoczecie *shared.StudioSectionStart) (string, error) {

	dlugosc := postacDlugosc(&stan.forma)
	poprzednia, err := a.stronaSekcjaZadania(ctx, stan, nil)
	if err != nil {
		return "", err
	}
	for i := range stan.forma.Sections {
		sekcja := &stan.forma.Sections[i]
		if sekcja.RangeStart <= miejsce && miejsce <= sekcja.RangeEnd {
			poprzednia = sekcja
			break
		}
	}
	koniecPoprzedniej := poprzednia.RangeEnd
	if koniecPoprzedniej < miejsce {
		koniecPoprzedniej = dlugosc
	}
	nowa := shared.StudioSection{
		Id:         nowyIdentyfikator(przedrostekSekcjiPostaci),
		Index:      poprzednia.Index + 1,
		RangeStart: miejsce,
		RangeEnd:   koniecPoprzedniej,
	}
	start := shared.StudioSectionStart(shared.StudioSectionStartNewPage)
	if rozpoczecie != nil && strings.TrimSpace(string(*rozpoczecie)) != "" {
		znany := false
		for _, wartosc := range shared.WartosciStudioSectionStart() {
			if wartosc == *rozpoczecie {
				znany = true
				break
			}
		}
		if !znany {
			nazwy := make([]string, 0, 5)
			for _, wartosc := range shared.WartosciStudioSectionStart() {
				nazwy = append(nazwy, string(wartosc))
			}
			return "", bladWskazaniaStudio("sposobu rozpoczęcia sekcji „" +
				string(*rozpoczecie) + "” kontrakt nie zna; sposoby: " +
				strings.Join(nazwy, ", "))
		}
		start = *rozpoczecie
	}
	nowa.Start = &start

	poprzednia.RangeEnd = miejsce
	if err := a.stronaZapiszSekcje(ctx, stan.dokument.ID, *poprzednia); err != nil {
		return "", err
	}
	// Sekcje za miejscem podziału przesuwają się o jedno oczko numeracji.
	for i := range stan.forma.Sections {
		sekcja := &stan.forma.Sections[i]
		if sekcja.Id == poprzednia.Id || sekcja.Index < nowa.Index {
			continue
		}
		sekcja.Index++
		if err := a.stronaZapiszSekcje(ctx, stan.dokument.ID, *sekcja); err != nil {
			return "", err
		}
	}
	if err := a.stronaZapiszSekcje(ctx, stan.dokument.ID, nowa); err != nil {
		return "", err
	}
	if err := a.stronaOdswierzSekcje(ctx, stan); err != nil {
		return "", err
	}
	return nowa.Id, nil
}

// ── Sekcje ──────────────────────────────────────────────────────────────────

// SekcjeDokumentu oddaje sekcje wraz z nastawami, nagłówkami i numeracją —
// obsługuje studio.section.list. Dokument bez ani jednej sekcji zapisanej
// oddaje sekcję jedną, obejmującą całą treść, bez zakładania jej w bazie.
func (a *adapterStudia) SekcjeDokumentu(ctx context.Context,
	z shared.StudioSectionListRequest) (shared.StudioSectionListResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioSectionListResponse{}, err
	}
	if len(stan.forma.Sections) > 0 {
		sort.SliceStable(stan.forma.Sections, func(i, j int) bool {
			return stan.forma.Sections[i].Index < stan.forma.Sections[j].Index
		})
		return shared.StudioSectionListResponse{Sections: stan.forma.Sections}, nil
	}
	rozpoczecie := shared.StudioSectionStart(shared.StudioSectionStartContinuous)
	domyslna := shared.StudioSection{
		Id:             "",
		Index:          0,
		Title:          postacWskaznikTekstu(stronaNazwaSekcjiPierwszej),
		RangeStart:     0,
		RangeEnd:       postacDlugosc(&stan.forma),
		Start:          &rozpoczecie,
		PageSetup:      stan.forma.PageSetup,
		HeadersFooters: stronaNaglowkiZNastaw(stan.forma.PageSetup),
	}
	return shared.StudioSectionListResponse{
		Sections: []shared.StudioSection{domyslna},
	}, nil
}

// ZapiszSekcjeDokumentu zakłada sekcję o własnych nastawach albo zmienia zastaną
// (`studio.section.save`).
func (a *adapterStudia) ZapiszSekcjeDokumentu(ctx context.Context,
	z shared.StudioSectionSaveRequest) (shared.StudioSectionSaveResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioSectionSaveResponse{}, err
	}
	autor := postacAutor(z.Author)
	dlugosc := postacDlugosc(&stan.forma)

	var sekcja shared.StudioSection
	nowa := true
	if z.SectionId != nil && strings.TrimSpace(*z.SectionId) != "" {
		zastana := stronaSekcjaPoKodzie(&stan.forma, *z.SectionId)
		if zastana == nil {
			return shared.StudioSectionSaveResponse{}, bladWskazaniaStudio(
				"sekcji „" + strings.TrimSpace(*z.SectionId) + "” dokument " +
					stan.dokument.Kod + " nie ma; sekcję nową zakłada się bez pola sectionId")
		}
		sekcja = *zastana
		nowa = false
	} else {
		sekcja = shared.StudioSection{
			Id:         nowyIdentyfikator(przedrostekSekcjiPostaci),
			Index:      len(stan.forma.Sections),
			RangeStart: 0,
			RangeEnd:   dlugosc,
		}
		rozpoczecie := shared.StudioSectionStart(shared.StudioSectionStartNewPage)
		sekcja.Start = &rozpoczecie
	}

	zmian := 0
	if z.Title != nil {
		sekcja.Title = postacWskaznikTekstu(strings.TrimSpace(*z.Title))
		zmian++
	}
	if z.RangeStart != nil {
		sekcja.RangeStart = *z.RangeStart
		zmian++
	}
	if z.RangeEnd != nil {
		sekcja.RangeEnd = *z.RangeEnd
		zmian++
	}
	if sekcja.RangeStart > sekcja.RangeEnd {
		sekcja.RangeStart, sekcja.RangeEnd = sekcja.RangeEnd, sekcja.RangeStart
	}
	if sekcja.RangeStart < 0 {
		sekcja.RangeStart = 0
	}
	if sekcja.RangeEnd > dlugosc {
		sekcja.RangeEnd = dlugosc
	}
	if z.Start != nil && strings.TrimSpace(string(*z.Start)) != "" {
		znany := false
		for _, wartosc := range shared.WartosciStudioSectionStart() {
			if wartosc == *z.Start {
				znany = true
				break
			}
		}
		if !znany {
			nazwy := make([]string, 0, 5)
			for _, wartosc := range shared.WartosciStudioSectionStart() {
				nazwy = append(nazwy, string(wartosc))
			}
			return shared.StudioSectionSaveResponse{}, bladWskazaniaStudio(
				"sposobu rozpoczęcia sekcji „" + string(*z.Start) + "” kontrakt nie zna; " +
					"sposoby: " + strings.Join(nazwy, ", "))
		}
		start := *z.Start
		sekcja.Start = &start
		zmian++
	}
	if len(z.PageSetup) > 0 {
		var nastawy shared.StudioPageSetup
		if err := json.Unmarshal(z.PageSetup, &nastawy); err != nil {
			return shared.StudioSectionSaveResponse{}, bladWskazaniaStudio(
				"nastawy strony sekcji są nieczytelne: " + err.Error())
		}
		if nastawy.PageSize != nil && strings.TrimSpace(*nastawy.PageSize) != "" {
			if _, jest := stronaNosnik(*nastawy.PageSize); !jest &&
				nastawy.WidthMm == nil && nastawy.HeightMm == nil {

				return shared.StudioSectionSaveResponse{}, bladWskazaniaStudio(
					"nośnika „" + strings.TrimSpace(*nastawy.PageSize) + "” rdzeń nie zna. " +
						"Nośniki znane: " + strings.Join(stronaNazwyNosnikow(), ", ") +
						". Format własny podaje się polami widthMm i heightMm")
			}
		}
		sekcja.PageSetup = &nastawy
		zmian++
	}
	if zmian == 0 && !nowa {
		return shared.StudioSectionSaveResponse{}, bladWskazaniaStudio(
			"zapis sekcji bez ani jednej rzeczy do zapisania — żądanie nie niosło nazwy, " +
				"zakresu, sposobu rozpoczęcia ani nastaw strony")
	}
	if sekcja.Title == nil && nowa {
		sekcja.Title = postacWskaznikTekstu("Sekcja " + postacZapisLiczby(sekcja.Index+1))
	}

	if err := a.stronaZapiszSekcje(ctx, stan.dokument.ID, sekcja); err != nil {
		return shared.StudioSectionSaveResponse{}, err
	}
	if err := a.stronaOdswierzSekcje(ctx, stan); err != nil {
		return shared.StudioSectionSaveResponse{}, err
	}
	zapisana := stronaSekcjaPoKodzie(&stan.forma, sekcja.Id)
	if zapisana == nil {
		return shared.StudioSectionSaveResponse{}, postacBladZaplecza(
			"sekcja " + sekcja.Id + " zapisana, ale nie wróciła z warstwy danych")
	}

	czynnosc := "zmieniona"
	if nowa {
		czynnosc = "założona"
	}
	stan.opisCzynnosci = "sekcja " + sekcja.Id + " " + czynnosc + " na zakresie " +
		postacZapisZakresu(sekcja.RangeStart, sekcja.RangeEnd)
	bilans := shared.StudioActionBalance{
		Applied: zmian + 1,
		Skipped: stronaBilansUkladu(&stan.forma,
			postacWskaznikNastaw(stronaNastawySkuteczne(stan.forma.PageSetup, zapisana.PageSetup))),
		Note: postacWskaznikTekstu(stan.opisCzynnosci),
	}
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindPageChange,
		sekcja.RangeStart, sekcja.RangeEnd, bilans)
	if err != nil {
		return shared.StudioSectionSaveResponse{}, err
	}
	wynikowa := stronaSekcjaPoKodzie(&forma, sekcja.Id)
	if wynikowa == nil {
		wynikowa = zapisana
	}
	return shared.StudioSectionSaveResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Section: *wynikowa,
	}, nil
}

// postacWskaznikNastaw oddaje wskaźnik na nastawy strony — rachunek bilansu
// bierze wskaźnik, a adresu wartości zwróconej z funkcji wziąć nie można.
func postacWskaznikNastaw(nastawy shared.StudioPageSetup) *shared.StudioPageSetup {
	kopia := nastawy
	return &kopia
}

// UsunSekcjeDokumentu usuwa sekcję — obsługuje studio.section.delete.
// Treść sekcji nie ginie: przechodzi do sekcji poprzedniej wraz z jej
// nastawami. Usunięcie sekcji jedynej jest odmową nazwaną.
func (a *adapterStudia) UsunSekcjeDokumentu(ctx context.Context,
	z shared.StudioSectionDeleteRequest) (shared.StudioSectionDeleteResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioSectionDeleteResponse{}, err
	}
	kod := strings.TrimSpace(z.SectionId)
	if kod == "" {
		return shared.StudioSectionDeleteResponse{}, bladWskazaniaStudio(
			"usunięcie sekcji bez wskazania sekcji")
	}
	usuwana := stronaSekcjaPoKodzie(&stan.forma, kod)
	if usuwana == nil {
		return shared.StudioSectionDeleteResponse{}, bladWskazaniaStudio(
			"sekcji „" + kod + "” dokument " + stan.dokument.Kod + " nie ma")
	}
	if len(stan.forma.Sections) == 1 {
		return shared.StudioSectionDeleteResponse{}, bladWskazaniaStudio(
			"sekcji „" + kod + "” nie da się usunąć: jest jedyną sekcją dokumentu " +
				stan.dokument.Kod + ", a dokument bez sekcji nie ma gdzie trzymać nośnika, " +
				"nagłówka ani numeracji. Zmień jej nastawy komendą studio.section.save")
	}
	autor := postacAutor(z.Author)
	skladnica, err := a.postacSkladnica()
	if err != nil {
		return shared.StudioSectionDeleteResponse{}, err
	}

	sort.SliceStable(stan.forma.Sections, func(i, j int) bool {
		return stan.forma.Sections[i].Index < stan.forma.Sections[j].Index
	})
	wskazanie := -1
	for i := range stan.forma.Sections {
		if stan.forma.Sections[i].Id == kod {
			wskazanie = i
			break
		}
	}
	zakresOd, zakresDo := usuwana.RangeStart, usuwana.RangeEnd
	przejela := ""
	switch {
	case wskazanie > 0:
		poprzednia := &stan.forma.Sections[wskazanie-1]
		if zakresDo > poprzednia.RangeEnd {
			poprzednia.RangeEnd = zakresDo
		}
		if err := a.stronaZapiszSekcje(ctx, stan.dokument.ID, *poprzednia); err != nil {
			return shared.StudioSectionDeleteResponse{}, err
		}
		przejela = poprzednia.Id
	case len(stan.forma.Sections) > 1:
		// Sekcja pierwsza usuwana: treść przechodzi do sekcji następnej, bo
		// poprzedniej nie ma.
		nastepna := &stan.forma.Sections[1]
		if zakresOd < nastepna.RangeStart {
			nastepna.RangeStart = zakresOd
		}
		if err := a.stronaZapiszSekcje(ctx, stan.dokument.ID, *nastepna); err != nil {
			return shared.StudioSectionDeleteResponse{}, err
		}
		przejela = nastepna.Id
	}

	usunieta, err := skladnica.UsunSekcje(ctx, kod)
	if err != nil {
		return shared.StudioSectionDeleteResponse{}, bladStudio(err)
	}
	if !usunieta {
		return shared.StudioSectionDeleteResponse{}, postacBladZaplecza(
			"sekcja " + kod + " nie dała się usunąć, choć stoi w postaci dokumentu")
	}
	// Blok sekcji usuniętej przechodzi do sekcji przejmującej, nie zostaje
	// bez sekcji.
	for i := range stan.forma.Blocks {
		if stan.forma.Blocks[i].SectionId != nil && *stan.forma.Blocks[i].SectionId == kod {
			if przejela == "" {
				stan.forma.Blocks[i].SectionId = nil
				continue
			}
			stan.forma.Blocks[i].SectionId = postacWskaznikTekstu(przejela)
		}
	}
	if err := a.stronaOdswierzSekcje(ctx, stan); err != nil {
		return shared.StudioSectionDeleteResponse{}, err
	}

	nota := "sekcja " + kod + " usunięta; jej treść " +
		postacZapisZakresu(zakresOd, zakresDo) + " przeszła do sekcji " + przejela
	if przejela == "" {
		nota = "sekcja " + kod + " usunięta"
	}
	stan.opisCzynnosci = nota
	bilans := shared.StudioActionBalance{
		Applied: 1,
		Skipped: []shared.StudioSkippedItem{},
		Note:    postacWskaznikTekstu(nota),
	}
	forma, bilansGotowy, _, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindPageChange,
		zakresOd, zakresDo, bilans)
	if err != nil {
		return shared.StudioSectionDeleteResponse{}, err
	}
	return shared.StudioSectionDeleteResponse{
		Deleted: true, Form: forma, Balance: bilansGotowy,
	}, nil
}

// ── Tabulatory linijki ──────────────────────────────────────────────────────

// UstawTabulatorLinijki zakłada, przestawia albo zdejmuje tabulator —
// obsługuje studio.ruler.tabstop.set. Tabulator jest cechą akapitu, więc
// czynność obejmuje wszystkie akapity, które zaznaczenie dotyka.
func (a *adapterStudia) UstawTabulatorLinijki(ctx context.Context,
	z shared.StudioRulerTabstopSetRequest) (shared.StudioRulerTabstopSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioRulerTabstopSetResponse{}, err
	}
	if z.PositionMm < 0 {
		return shared.StudioRulerTabstopSetResponse{}, bladWskazaniaStudio(
			"tabulator liczy się od lewego marginesu i nie stoi na wartości ujemnej")
	}
	zdejmowanie := z.Remove != nil && *z.Remove
	rodzaj := shared.StudioTabKind(shared.StudioTabKindLeft)
	if z.Kind != nil && strings.TrimSpace(string(*z.Kind)) != "" {
		znany := false
		for _, wartosc := range shared.WartosciStudioTabKind() {
			if wartosc == *z.Kind {
				znany = true
				break
			}
		}
		if !znany {
			nazwy := make([]string, 0, 5)
			for _, wartosc := range shared.WartosciStudioTabKind() {
				nazwy = append(nazwy, string(wartosc))
			}
			return shared.StudioRulerTabstopSetResponse{}, bladWskazaniaStudio(
				"rodzaju tabulatora „" + string(*z.Kind) + "” kontrakt nie zna; rodzaje: " +
					strings.Join(nazwy, ", "))
		}
		rodzaj = *z.Kind
	}
	autor := postacAutor(z.Author)
	od, do := postacZakres(z.RangeStart, z.RangeEnd, postacDlugosc(&stan.forma))
	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor)

	bilans := shared.StudioActionBalance{Skipped: pominiete}
	tabulatory := []shared.StudioTabStop{}
	for _, odcinek := range odcinki {
		for _, wskazanie := range postacBlokiZakresu(&stan.forma, odcinek[0], odcinek[1]) {
			blok := &stan.forma.Blocks[wskazanie]
			zastane := []shared.StudioTabStop{}
			if blok.Paragraph != nil && blok.Paragraph.TabStops != nil {
				zastane = append(zastane, blok.Paragraph.TabStops...)
			}
			nowe, zmieniono := stronaTabulatory(zastane, z.PositionMm, rodzaj, z.Leader, zdejmowanie)
			if !zmieniono {
				continue
			}
			blok.Paragraph = postacScalAkapit(blok.Paragraph,
				shared.StudioParagraphFormat{TabStops: nowe})
			tabulatory = nowe
			bilans.Applied++
		}
	}
	if bilans.Applied == 0 {
		if zdejmowanie {
			return shared.StudioRulerTabstopSetResponse{}, bladWskazaniaStudio(
				"tabulatora na położeniu " + stronaZapisMiary(z.PositionMm) +
					" mm nie ma w ani jednym akapicie zakresu " + postacZapisZakresu(od, do) +
					" — nie ma czego zdjąć")
		}
		return shared.StudioRulerTabstopSetResponse{}, bladWskazaniaStudio(
			"tabulator nie miał na czym stanąć: zakres " + postacZapisZakresu(od, do) +
				" nie obejmuje ani jednego akapitu wolnego od blokady")
	}
	stan.opisCzynnosci = "tabulator na położeniu " + stronaZapisMiary(z.PositionMm) + " mm " +
		map[bool]string{true: "zdjęty", false: "ustawiony"}[zdejmowanie] + " " +
		postacZapisZakresu(od, do)
	bilans.Note = postacWskaznikTekstu(stan.opisCzynnosci)

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindPageChange, od, do, bilans)
	if err != nil {
		return shared.StudioRulerTabstopSetResponse{}, err
	}
	return shared.StudioRulerTabstopSetResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, TabStops: tabulatory,
	}, nil
}

// stronaTabulatory zakłada, przestawia albo zdejmuje tabulator w wykazie
// akapitu i oddaje, czy wykaz się zmienił; tabulator na tym samym położeniu
// jest tym samym tabulatorem w tolerancji jednej dziesiątej milimetra.
func stronaTabulatory(zastane []shared.StudioTabStop, polozenie float64,
	rodzaj shared.StudioTabKind, znak *shared.StudioTabLeader,
	zdejmowanie bool) ([]shared.StudioTabStop, bool) {

	const tolerancja = 0.1
	wynik := make([]shared.StudioTabStop, 0, len(zastane)+1)
	znaleziony := false
	for _, tabulator := range zastane {
		roznica := tabulator.PositionMm - polozenie
		if roznica < 0 {
			roznica = -roznica
		}
		if roznica > tolerancja {
			wynik = append(wynik, tabulator)
			continue
		}
		znaleziony = true
		if zdejmowanie {
			continue
		}
		tabulator.Kind = rodzaj
		if znak != nil {
			wiodacy := *znak
			tabulator.Leader = &wiodacy
		}
		wynik = append(wynik, tabulator)
	}
	if zdejmowanie {
		return wynik, znaleziony
	}
	if !znaleziony {
		tabulator := shared.StudioTabStop{PositionMm: polozenie, Kind: rodzaj}
		if znak != nil {
			wiodacy := *znak
			tabulator.Leader = &wiodacy
		}
		wynik = append(wynik, tabulator)
	}
	sort.SliceStable(wynik, func(i, j int) bool {
		return wynik[i].PositionMm < wynik[j].PositionMm
	})
	return wynik, true
}
