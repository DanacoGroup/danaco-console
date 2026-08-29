// Odpowiedzialność pliku: osiem czynności wejścia do edytora modułu Studio,
// od założenia dokumentu pustego po wniesienie pliku, obrazu, fragmentu
// z Biblioteki i ze strony sieci, każda z zapisaną postacią i bilansem.
package core

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// ── Nowy dokument ───────────────────────────────────────────────────────────

// ZalozDokument obsługuje `studio.document.create` — zakłada dokument pusty,
// czyli nową stronę gotową do pisania, z arkuszem stylów i nastawami strony
// domyślnymi albo przejętymi z szablonu, wraz z pustym akapitem, na którym
// staje kursor.
func (a *adapterStudia) ZalozDokument(ctx context.Context,
	z shared.StudioDocumentCreateRequest) (shared.StudioDocumentCreateResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.StudioDocumentCreateResponse{}, bladWskazaniaStudio(
			"założenie dokumentu bez okna, w którym ma stanąć")
	}
	autor, err := wejscieAutorCzynnosci(z.Author)
	if err != nil {
		return shared.StudioDocumentCreateResponse{}, err
	}

	format := shared.StudioDocumentFormat(shared.StudioDocumentFormatTxt)
	if z.Format != nil && strings.TrimSpace(string(*z.Format)) != "" {
		format = *z.Format
	}
	tytul := z.Title

	var szablon *shared.StudioTemplateDetail
	if kod := strings.TrimSpace(wartoscTekstu(z.TemplateId)); kod != "" {
		wczytany, err := a.szablonPismaSzczegol(ctx, kod)
		if err != nil {
			return shared.StudioDocumentCreateResponse{}, err
		}
		szablon = &wczytany
		if tytul == nil {
			nazwa := wczytany.Name
			tytul = &nazwa
		}
		if strings.TrimSpace(string(wczytany.Format)) != "" {
			format = wczytany.Format
		}
	}

	dokument, err := a.wejscieDokumentAlboNowy(ctx, nil, z.WindowId, tytul, format)
	if err != nil {
		return shared.StudioDocumentCreateResponse{}, err
	}

	stan := &stanPostaci{dokument: dokument}
	if szablon != nil && szablon.Form != nil {
		stan.forma = *szablon.Form
		stan.forma.DocumentId = dokument.Kod
		wejsciePrzepiszIdentyfikatory(&stan.forma)
		stan.opisCzynnosci = "założenie dokumentu z szablonu „" + szablon.Name + "”"
	} else {
		nazwaNosnika := strings.TrimSpace(wartoscTekstu(z.PaperName))
		stan.forma = wejscieNowaPostac(dokument.Kod, nazwaNosnika, z.Orientation)
		stan.opisCzynnosci = "założenie dokumentu pustego"
	}
	// Nastawy nośnika wskazane wprost mają pierwszeństwo nad nastawami szablonu.
	if nazwa := strings.TrimSpace(wartoscTekstu(z.PaperName)); nazwa != "" && stan.forma.PageSetup != nil {
		stan.forma.PageSetup.PageSize = wejscieWskaznikTekstu(nazwa)
	}
	if z.Orientation != nil && stan.forma.PageSetup != nil {
		stan.forma.PageSetup.Orientation = z.Orientation
	}

	if err := a.wejscieUtrwalPostac(ctx, stan); err != nil {
		return shared.StudioDocumentCreateResponse{}, err
	}
	if szablon != nil {
		// Blokady wzorcowe szablonu idą do dokumentu z niego zakładanego.
		if _, err := a.blokadaPrzenieSzablon(ctx, stan.dokument.ID, szablon.Id); err != nil {
			return shared.StudioDocumentCreateResponse{}, err
		}
	}
	if err := a.wejscieOdlozCzynnosc(ctx, stan, autor); err != nil {
		return shared.StudioDocumentCreateResponse{}, err
	}

	return shared.StudioDocumentCreateResponse{
		Document: a.zlozDokument(stan.dokument),
		Form:     stan.forma,
	}, nil
}

// wejscieOdlozCzynnosc dopisuje czynność wejścia do odwracalnego dziennika
// dokumentu, zawsze rodzajem `importChange`, po którym wpis rozpoznaje się
// w wykazie jako wniesienie, odróżnione od zwykłej zmiany treści.
func (a *adapterStudia) wejscieOdlozCzynnosc(ctx context.Context, stan *stanPostaci,
	autor shared.StudioAuthor) error {

	dlugosc := postacDlugosc(&stan.forma)
	// Wykonawca wchodzi z podpisu żądania, tak samo jak przy zmianie postaci.
	return a.postacOdlozCzynnosc(ctx, stan, postacWykonawca(ctx, autor),
		shared.StudioActionKindImportChange, 0, dlugosc, nil,
		shared.StudioActionBalance{Applied: 1})
}

// wejsciePrzepiszIdentyfikatory nadaje nowe identyfikatory wszystkim bytom
// postaci przenoszonej między dokumentami, żeby kopia i oryginał nie dzieliły
// kluczy wierszy jednoznacznych w całej bazie.
func wejsciePrzepiszIdentyfikatory(forma *shared.StudioDocumentForm) {
	nowe := map[string]string{}
	przepisz := func(stary, przedrostek string) string {
		if stary == "" {
			return nowyIdentyfikator(przedrostek)
		}
		if nowy, jest := nowe[stary]; jest {
			return nowy
		}
		nowy := nowyIdentyfikator(przedrostek)
		nowe[stary] = nowy
		return nowy
	}

	for i := range forma.Sections {
		forma.Sections[i].Id = przepisz(forma.Sections[i].Id, przedrostekSekcjiStudia)
	}
	for i := range forma.Tables {
		forma.Tables[i].Id = przepisz(forma.Tables[i].Id, przedrostekTabeliStudia)
	}
	for i := range forma.Objects {
		forma.Objects[i].Id = przepisz(forma.Objects[i].Id, przedrostekObiektuStudia)
	}
	for i := range forma.Lists {
		forma.Lists[i].Id = przepisz(forma.Lists[i].Id, przedrostekListyWejscia)
	}
	for i := range forma.Apparatus {
		forma.Apparatus[i].Id = przepisz(forma.Apparatus[i].Id, przedrostekAparatuWejscia)
	}
	for i := range forma.Fields {
		forma.Fields[i].Id = przepisz(forma.Fields[i].Id, przedrostekPolaPostaci)
	}
	for i := range forma.Blocks {
		forma.Blocks[i].Id = nowyIdentyfikator(przedrostekBlokuStudia)
		if forma.Blocks[i].SectionId != nil {
			if nowy, jest := nowe[*forma.Blocks[i].SectionId]; jest {
				forma.Blocks[i].SectionId = &nowy
			}
		}
		if forma.Blocks[i].TableId != nil {
			if nowy, jest := nowe[*forma.Blocks[i].TableId]; jest {
				forma.Blocks[i].TableId = &nowy
			}
		}
		if forma.Blocks[i].ObjectId != nil {
			if nowy, jest := nowe[*forma.Blocks[i].ObjectId]; jest {
				forma.Blocks[i].ObjectId = &nowy
			}
		}
	}
	for i := range forma.Blocks {
		if forma.Blocks[i].Paragraph == nil || forma.Blocks[i].Paragraph.ListId == nil {
			continue
		}
		if nowy, jest := nowe[*forma.Blocks[i].Paragraph.ListId]; jest {
			forma.Blocks[i].Paragraph.ListId = &nowy
		}
	}
	// Blokady fragmentów nie dostają tu nowych identyfikatorów, tylko przy przeniesieniu.
	forma.Locks = nil
	forma.Revision = nil
	forma.UpdatedAt = nil
}

// ── Wniesienie pliku do edytora ─────────────────────────────────────────────

// WniesPlikDoEdytora obsługuje `studio.document.import.file`: format
// rozpoznaje się po zawartości pliku, nie po rozszerzeniu, a zapis znaków
// ustala się przed rozbiorem treści.
func (a *adapterStudia) WniesPlikDoEdytora(ctx context.Context,
	z shared.StudioDocumentImportFileRequest) (shared.StudioDocumentImportFileResponse, error) {

	autor, err := wejscieAutorCzynnosci(z.Author)
	if err != nil {
		return shared.StudioDocumentImportFileResponse{}, err
	}
	bajty, nazwaPliku, err := a.wejscieBajtyZrodla(ctx, wejscieWskazanieZrodla{
		Sciezka:        z.Path,
		ZasobKod:       z.AssetId,
		PlikBiblioteki: z.LibraryFileId,
		BajtyBase64:    z.BytesBase64,
		Czynnosc:       "wniesienie pliku do edytora",
	})
	if err != nil {
		return shared.StudioDocumentImportFileResponse{}, err
	}
	format, err := wejscieRozpoznajFormat(nazwaPliku, bajty, z.Format)
	if err != nil {
		return shared.StudioDocumentImportFileResponse{}, err
	}

	tytul := wejscieNazwaZPliku(nazwaPliku)
	dokument, err := a.wejscieDokumentAlboNowy(ctx, z.DocumentId, z.WindowId, tytul,
		wejscieFormatDokumentu(format))
	if err != nil {
		return shared.StudioDocumentImportFileResponse{}, err
	}

	postac, tresc, bilans, err := a.wejscieRozbierzPlik(dokument.Kod, bajty, format, z.Encoding)
	if err != nil {
		return shared.StudioDocumentImportFileResponse{}, err
	}

	stan := &stanPostaci{dokument: dokument}
	odKursora, doKursora := 0, 0

	if z.InsertAtOffset != nil && strings.TrimSpace(wartoscTekstu(z.DocumentId)) != "" {
		// Wniesienie w miejsce kursora do dokumentu istniejącego, bez postaci pliku.
		zastany, err := a.postacWczytaj(ctx, dokument.Kod)
		if err != nil {
			return shared.StudioDocumentImportFileResponse{}, err
		}
		stan = zastany
		dlugosc := postacDlugosc(&stan.forma)
		odKursora = *z.InsertAtOffset
		if odKursora < 0 {
			odKursora = 0
		}
		if odKursora > dlugosc {
			odKursora = dlugosc
		}
		doKursora = odKursora

		odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, odKursora, doKursora, autor)
		if len(odcinki) == 0 {
			return shared.StudioDocumentImportFileResponse{}, bladWskazaniaStudio(
				"wniesienie pliku zatrzymane w całości przez blokadę fragmentu: " +
					postacNazwaBlokad(pominiete))
		}
		autorWpisu := autor
		postacZamienTresc(&stan.forma, odKursora, doKursora, tresc, nil, &autorWpisu)
		doKursora = odKursora + len([]rune(tresc))
		bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
			Reason: "wniesienie w miejsce kursora wniosło treść z postacią znaku, bez " +
				"arkusza stylów i nastaw strony pliku",
			Detail: wejscieWskaznikTekstu("podmiana arkusza stylów dokumentu Operatora " +
				"nie była tym, o co prosił; naprawa: wnieść plik jako dokument osobny"),
		})
		bilans.Skipped = append(bilans.Skipped, pominiete...)
	} else {
		stan.forma = postac
		stan.forma.DocumentId = dokument.Kod
		stan.tekstPrzed = wartoscTekstu(dokument.Tresc)
		doKursora = len([]rune(tresc))
	}
	stan.opisCzynnosci = "wniesienie pliku " + strings.TrimSpace(nazwaPliku) +
		" (" + string(format) + ") do edytora"

	if err := a.wejscieUtrwalPostac(ctx, stan); err != nil {
		return shared.StudioDocumentImportFileResponse{}, err
	}
	if err := a.wejscieOdlozCzynnosc(ctx, stan, autor); err != nil {
		return shared.StudioDocumentImportFileResponse{}, err
	}

	tytulZrodla := strings.TrimSpace(nazwaPliku)
	pochodzenie, err := a.wejscieOdlozPochodzenie(ctx, stan.dokument,
		shared.StudioProvenanceKindImportedFile, odKursora, doKursora,
		z.Path, z.LibraryFileId, &tytulZrodla, autor, z.AgentId, z.AgentName)
	if err != nil {
		return shared.StudioDocumentImportFileResponse{}, err
	}

	return shared.StudioDocumentImportFileResponse{
		Document:   a.zlozDokument(stan.dokument),
		Form:       stan.forma,
		Balance:    bilans,
		Provenance: pochodzenie,
	}, nil
}

// wejscieRozbierzPlik kieruje bajty do rachunku właściwego dla formatu.
//
// Jedno miejsce rozstrzygnięcia, nie dwa: wniesienie dokumentu i wniesienie
// szablonu czytają plik tą samą drogą, bo `.dotx` różni się od `.docx`
// przeznaczeniem, a nie rozbiorem.
func (a *adapterStudia) wejscieRozbierzPlik(kodDokumentu string, bajty []byte,
	format shared.StudioImportFormat,
	wskazanieZapisu *string) (shared.StudioDocumentForm, string,
	shared.StudioImportBalance, error) {

	switch format {
	case shared.StudioImportFormatDocx, shared.StudioImportFormatDotx:
		return wejscieCzytajOoxml(kodDokumentu, bajty, format)
	case shared.StudioImportFormatOdt, shared.StudioImportFormatOtt:
		return wejscieCzytajOdf(kodDokumentu, bajty, format)
	case shared.StudioImportFormatPdf:
		odzyskane, err := wejscieCzytajPdf(kodDokumentu, bajty, "", true, false)
		if err != nil {
			return shared.StudioDocumentForm{}, "", odzyskane.Bilans, err
		}
		if odzyskane.SamSkan {
			return shared.StudioDocumentForm{}, "", odzyskane.Bilans, wejscieBladBraku(
				"dokument PDF nie ma warstwy tekstowej — to skan. Wniesienie pliku nie ma " +
					"czego wnieść do edytora; naprawa: `studio.document.import.pdf`, " +
					"który kieruje skan na rozpoznanie pisma i oddaje pozycję jego kolejki")
		}
		return odzyskane.Postac, odzyskane.Tresc, odzyskane.Bilans, nil
	}

	// Formaty tekstowe: zapis znaków rozpoznaje się PRZED rozbiorem treści.
	tekst, zapis, err := wejscieRozpoznajZapisZnakow(bajty, wskazanieZapisu)
	if err != nil {
		return shared.StudioDocumentForm{}, "", shared.StudioImportBalance{Format: format}, err
	}
	bilans := shared.StudioImportBalance{
		Format:   format,
		Encoding: wejscieWskaznikTekstu(zapis),
	}

	var postac shared.StudioDocumentForm
	var pominiete []shared.StudioSkippedItem
	switch format {
	case shared.StudioImportFormatMd:
		postac, pominiete = wejsciePostacZMarkdown(kodDokumentu, tekst)
	case shared.StudioImportFormatHtml:
		postac, pominiete = wejsciePostacZHtml(kodDokumentu, tekst)
	case shared.StudioImportFormatRtf:
		postac, pominiete = wejsciePostacZRtf(kodDokumentu, tekst)
	default:
		postac = wejsciePostacZTekstu(kodDokumentu, tekst)
		pominiete = []shared.StudioSkippedItem{{
			Reason: "plik tekstowy nie niesie postaci — dokument dostał arkusz stylów " +
				"domyślny platformy",
		}}
	}
	bilans.Skipped = append(bilans.Skipped, pominiete...)
	bilans.ParagraphsRecovered = wejscieWskaznikCalkowity(len(postac.Blocks))
	bilans.TablesRecognized = wejscieWskaznikCalkowity(len(postac.Tables))
	bilans.StylesRecovered = wejscieWskaznikCalkowity(len(postac.Styles))
	bilans.SectionsRecovered = wejscieWskaznikCalkowity(len(postac.Sections))
	bilans.ImagesEmbedded = wejscieWskaznikCalkowity(len(postac.Objects))
	tresc := wejscieTrescZPostaci(&postac)
	wejscieUzycieStylow(&postac)
	bilans.Note = wejscieWskaznikTekstu("plik wniesiony zapisem znaków " + zapis + "; " +
		ooxmlZdanieBilansu(&bilans))
	return postac, tresc, bilans, nil
}

// ── PDF na dokument edytowalny ──────────────────────────────────────────────

// WniesPdfDoEdytora obsługuje `studio.document.import.pdf`; PDF ze samych
// skanów nie udaje konwersji — wchodzi do kolejki rozpoznania pisma, a odpowiedź
// oddaje pozycję tej kolejki wraz z bilansem mówiącym, że warstwy tekstowej nie było.
func (a *adapterStudia) WniesPdfDoEdytora(ctx context.Context,
	z shared.StudioDocumentImportPdfRequest) (shared.StudioDocumentImportPdfResponse, error) {

	autor, err := wejscieAutorCzynnosci(z.Author)
	if err != nil {
		return shared.StudioDocumentImportPdfResponse{}, err
	}
	bajty, nazwaPliku, err := a.wejscieBajtyZrodla(ctx, wejscieWskazanieZrodla{
		Sciezka:        z.Path,
		ZasobKod:       z.AssetId,
		PlikBiblioteki: z.LibraryFileId,
		BajtyBase64:    z.BytesBase64,
		Czynnosc:       "zamiana PDF na dokument edytowalny",
	})
	if err != nil {
		return shared.StudioDocumentImportPdfResponse{}, err
	}

	odzyskacTabele := z.RecoverTables == nil || *z.RecoverTables
	osadzacObrazy := z.RecoverImages == nil || *z.RecoverImages
	odzyskane, err := wejscieCzytajPdf("", bajty, wartoscTekstu(z.Pages),
		odzyskacTabele, osadzacObrazy)
	if err != nil {
		return shared.StudioDocumentImportPdfResponse{}, err
	}

	tytul := wejscieNazwaZPliku(nazwaPliku)
	dokument, err := a.wejscieDokumentAlboNowy(ctx, z.DocumentId, z.WindowId, tytul,
		shared.StudioDocumentFormatTxt)
	if err != nil {
		return shared.StudioDocumentImportPdfResponse{}, err
	}

	if odzyskane.SamSkan {
		if strings.TrimSpace(z.WindowId) == "" {
			return shared.StudioDocumentImportPdfResponse{}, bladWskazaniaStudio(
				"skierowanie skanu na rozpoznanie pisma wymaga okna, do którego kolejka należy")
		}
		wskazanie := strings.TrimSpace(nazwaPliku)
		pozycja, err := a.repozytorium.ZapiszPozycjeWczytywania(ctx, dane.PozycjaWczytywania{
			Kod:             nowyIdentyfikator(przedrostekPozycjiWczytywania),
			Okno:            z.WindowId,
			SciezkaZrodlowa: &wskazanie,
			ZasobID:         z.AssetId,
			// Stan „oczekuje": rozpoznanie pisma jeszcze nie zaszło.
			Stan: "oczekuje",
		})
		if err != nil {
			return shared.StudioDocumentImportPdfResponse{}, bladStudio(err)
		}
		stan := &stanPostaci{dokument: dokument}
		stan.forma = shared.StudioDocumentForm{DocumentId: dokument.Kod}
		if zastany, err := a.postacWczytaj(ctx, dokument.Kod); err == nil {
			stan = zastany
		}
		return shared.StudioDocumentImportPdfResponse{
			Document:     a.zlozDokument(stan.dokument),
			Form:         stan.forma,
			Balance:      odzyskane.Bilans,
			IngestItemId: wejscieWskaznikTekstu(pozycja.Kod),
		}, nil
	}

	// Bajty obrazów wyjętych ze stron idą do magazynu zasobów rdzenia.
	odzyskane.Postac.DocumentId = dokument.Kod
	if len(odzyskane.Obrazy) > 0 {
		for _, obraz := range odzyskane.Obrazy {
			zasob, err := a.odlozObrazStudia(ctx, obraz.Bajty,
				"obraz z PDF strona "+strconv.Itoa(obraz.NumerStrony)+"."+obraz.Format,
				dokument.Okno, obraz.Szerokosc, obraz.Wysokosc)
			if err != nil {
				// Odzyskanie tekstu nie przepada — bilans mówi, że obrazy nie weszły.
				odzyskane.Bilans.Skipped = append(odzyskane.Bilans.Skipped,
					shared.StudioSkippedItem{
						Reason: "obrazu nie dało się odłożyć w magazynie zasobów serwera",
						Detail: wejscieWskaznikTekstu(err.Error()),
					})
				odzyskane.Bilans.ImagesEmbedded = wejscieWskaznikCalkowity(0)
				odzyskane.Bilans.ImagesSkipped = wejscieWskaznikCalkowity(
					len(odzyskane.Obrazy))
				break
			}
			for i := range odzyskane.Postac.Objects {
				if odzyskane.Postac.Objects[i].Id == obraz.ObiektKod {
					odzyskane.Postac.Objects[i].AssetId = wejscieWskaznikTekstu(zasob.Id)
					odzyskane.Postac.Objects[i].Source = wejscieWskaznikZrodlaObiektu(
						shared.StudioObjectSourceCoreAsset)
				}
			}
		}
	}

	stan := &stanPostaci{dokument: dokument, forma: odzyskane.Postac}
	stan.tekstPrzed = wartoscTekstu(dokument.Tresc)
	stan.opisCzynnosci = "zamiana PDF " + strings.TrimSpace(nazwaPliku) +
		" na dokument edytowalny"
	if err := a.wejscieUtrwalPostac(ctx, stan); err != nil {
		return shared.StudioDocumentImportPdfResponse{}, err
	}
	if err := a.wejscieOdlozCzynnosc(ctx, stan, autor); err != nil {
		return shared.StudioDocumentImportPdfResponse{}, err
	}
	tytulZrodla := strings.TrimSpace(nazwaPliku)
	if _, err := a.wejscieOdlozPochodzenie(ctx, stan.dokument,
		shared.StudioProvenanceKindImportedFile, 0, postacDlugosc(&stan.forma),
		z.Path, z.LibraryFileId, &tytulZrodla, autor, z.AgentId, z.AgentName); err != nil {
		return shared.StudioDocumentImportPdfResponse{}, err
	}

	odpowiedz := shared.StudioDocumentImportPdfResponse{
		Document: a.zlozDokument(stan.dokument),
		Form:     stan.forma,
		Balance:  odzyskane.Bilans,
	}
	return odpowiedz, nil
}

// ── Wniesienie obrazu ───────────────────────────────────────────────────────

// WniesObraz obsługuje `studio.document.image.import` — wnosi obraz wprost
// w miejsce kursora, z pliku, z magazynu zasobów rdzenia, z modułu Design albo
// z bazy zdjęciowej.
func (a *adapterStudia) WniesObraz(ctx context.Context,
	z shared.StudioDocumentImageImportRequest) (shared.StudioDocumentImageImportResponse, error) {

	autor, err := wejscieAutorCzynnosci(z.Author)
	if err != nil {
		return shared.StudioDocumentImageImportResponse{}, err
	}
	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDocumentImageImportResponse{}, err
	}

	dlugosc := postacDlugosc(&stan.forma)
	miejsce := z.Offset
	if miejsce < 0 {
		miejsce = 0
	}
	if miejsce > dlugosc {
		miejsce = dlugosc
	}
	// Blokada obowiązuje przed dotknięciem treści, nie po nim.
	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, miejsce, miejsce, autor)
	if len(odcinki) == 0 {
		return shared.StudioDocumentImageImportResponse{}, bladWskazaniaStudio(
			"wniesienie obrazu zatrzymane przez blokadę fragmentu: " +
				postacNazwaBlokad(pominiete))
	}

	obiekt, pochodzenieRodzaj, adresZrodla, err := a.wejscieObiektObrazu(ctx, stan.dokument, z)
	if err != nil {
		return shared.StudioDocumentImageImportResponse{}, err
	}
	obiekt.AnchorOffset = wejscieWskaznikCalkowity(miejsce)
	obiekt.Anchor = wejscieWskaznikZakotwiczenia(shared.StudioAnchorKindCharacter)
	if z.WidthMm != nil {
		obiekt.WidthMm = z.WidthMm
	}
	if z.HeightMm != nil {
		obiekt.HeightMm = z.HeightMm
	}
	if z.AltText != nil {
		obiekt.AltText = z.AltText
	}
	if z.Caption != nil {
		obiekt.Caption = z.Caption
	}

	stan.forma.Objects = append(stan.forma.Objects, obiekt)
	stan.forma.Blocks = append(stan.forma.Blocks, shared.StudioDocumentBlock{
		Id:        nowyIdentyfikator(przedrostekBlokuStudia),
		Kind:      wejscieRodzajBlokuObiekt,
		SectionId: wejscieSekcjaMiejsca(&stan.forma, miejsce),
		ObjectId:  wejscieWskaznikTekstu(obiekt.Id),
	})
	stan.opisCzynnosci = "wniesienie obrazu w miejsce " + strconv.Itoa(miejsce)

	bilans := shared.StudioActionBalance{Applied: 1, Skipped: pominiete}
	bilans.Note = wejscieWskaznikTekstu("obraz osadzony w dokumencie w miejscu " +
		strconv.Itoa(miejsce))

	if err := a.wejscieUtrwalPostac(ctx, stan); err != nil {
		return shared.StudioDocumentImageImportResponse{}, err
	}
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKind(shared.StudioChangeKindWstawienie),
		shared.StudioActionKindImportChange, miejsce, miejsce, bilans)
	if err != nil {
		return shared.StudioDocumentImageImportResponse{}, err
	}

	pochodzenie, err := a.wejscieOdlozPochodzenie(ctx, stan.dokument, pochodzenieRodzaj,
		miejsce, miejsce, adresZrodla, z.LibraryFileId, obiekt.AltText, autor,
		z.AgentId, z.AgentName)
	if err != nil {
		return shared.StudioDocumentImageImportResponse{}, err
	}

	return shared.StudioDocumentImageImportResponse{
		Form:       forma,
		Balance:    bilansGotowy,
		Change:     zmiana,
		ActionId:   stan.czynnosc,
		Object:     obiekt,
		Provenance: pochodzenie,
	}, nil
}

// wejscieObiektObrazu składa obiekt obrazu wedle wskazanego źródła i odkłada
// bajty w magazynie zasobów rdzenia, gdy przyszły z zewnątrz; źródło modułu
// Design przyjmuje kod węzła do zapisu pochodzenia, a bajty bierze z zasobu.
func (a *adapterStudia) wejscieObiektObrazu(ctx context.Context, dokument dane.DokumentStudia,
	z shared.StudioDocumentImageImportRequest) (shared.StudioDocumentObject,
	shared.StudioProvenanceKind, *string, error) {

	obiekt := shared.StudioDocumentObject{
		Id:            nowyIdentyfikator(przedrostekObiektuStudia),
		Kind:          shared.StudioObjectKindImage,
		Source:        wejscieWskaznikZrodlaObiektu(z.Source),
		DesignNodeId:  z.DesignNodeId,
		LibraryFileId: z.LibraryFileId,
	}
	rodzajPochodzenia := shared.StudioProvenanceKind(shared.StudioProvenanceKindImportedFile)
	var adres *string

	switch z.Source {
	case shared.StudioObjectSourceCoreAsset:
		kod := strings.TrimSpace(wartoscTekstu(z.AssetId))
		if kod == "" {
			return obiekt, rodzajPochodzenia, nil, bladWskazaniaStudio(
				"wniesienie obrazu z magazynu zasobów bez wskazania zasobu")
		}
		if _, err := a.bajtyZasobuStudia(ctx, kod); err != nil {
			return obiekt, rodzajPochodzenia, nil, err
		}
		obiekt.AssetId = wejscieWskaznikTekstu(kod)
		return obiekt, rodzajPochodzenia, nil, nil

	case shared.StudioObjectSourceDesignModule:
		if strings.TrimSpace(wartoscTekstu(z.DesignNodeId)) == "" {
			return obiekt, rodzajPochodzenia, nil, bladWskazaniaStudio(
				"wniesienie obrazu z modułu Design bez wskazania węzła")
		}
		kod := strings.TrimSpace(wartoscTekstu(z.AssetId))
		if kod == "" {
			return obiekt, rodzajPochodzenia, nil, bladWskazaniaStudio(
				"wniesienie węzła modułu Design wymaga jego wyrysu w magazynie zasobów: " +
					"węzeł wektorowy nie jest obrazem, dopóki nie zostanie wyrysowany. " +
					"Naprawa: wydać go komendą design.vector.export i podać kod zasobu " +
					"w polu assetId obok kodu węzła")
		}
		if _, err := a.bajtyZasobuStudia(ctx, kod); err != nil {
			return obiekt, rodzajPochodzenia, nil, err
		}
		obiekt.AssetId = wejscieWskaznikTekstu(kod)
		return obiekt, rodzajPochodzenia, nil, nil

	case shared.StudioObjectSourceWeb, shared.StudioObjectSourcePhotoBank:
		adresZrodla := strings.TrimSpace(wartoscTekstu(z.Path))
		if adresZrodla == "" {
			adresZrodla = strings.TrimSpace(wartoscTekstu(z.PhotoBankId))
		}
		if adresZrodla == "" {
			return obiekt, rodzajPochodzenia, nil, bladWskazaniaStudio(
				"wniesienie obrazu z bazy zdjęciowej bez wskazania jego adresu")
		}
		bajty, err := a.wejscieObrazZeSieci(ctx, adresZrodla)
		if err != nil {
			return obiekt, rodzajPochodzenia, nil, err
		}
		zasob, err := a.odlozTrescStudia(ctx, bajty, wejscieNazwaObrazuZAdresu(adresZrodla),
			"png", dokument.Okno)
		if err != nil {
			return obiekt, rodzajPochodzenia, nil, err
		}
		obiekt.AssetId = wejscieWskaznikTekstu(zasob.Id)
		obiekt.SourceUrl = wejscieWskaznikTekstu(adresZrodla)
		return obiekt, shared.StudioProvenanceKindWeb, &adresZrodla, nil
	}

	// Plik, plik Biblioteki albo bajty wprost — jedna droga do bajtów.
	bajty, nazwa, err := a.wejscieBajtyZrodla(ctx, wejscieWskazanieZrodla{
		Sciezka:        z.Path,
		ZasobKod:       z.AssetId,
		PlikBiblioteki: z.LibraryFileId,
		BajtyBase64:    z.BytesBase64,
		Czynnosc:       "wniesienie obrazu",
	})
	if err != nil {
		return obiekt, rodzajPochodzenia, nil, err
	}
	format := strings.TrimPrefix(strings.ToLower(wejscieRozszerzenieObrazu(nazwa)), ".")
	zasob, err := a.odlozTrescStudia(ctx, bajty, nazwa, format, dokument.Okno)
	if err != nil {
		return obiekt, rodzajPochodzenia, nil, err
	}
	obiekt.AssetId = wejscieWskaznikTekstu(zasob.Id)
	if obiekt.AltText == nil && strings.TrimSpace(nazwa) != "" {
		obiekt.AltText = wejscieWskaznikTekstu(strings.TrimSpace(nazwa))
	}
	if strings.TrimSpace(wartoscTekstu(z.LibraryFileId)) != "" {
		rodzajPochodzenia = shared.StudioProvenanceKindLibraryFile
	}
	adres = z.Path
	return obiekt, rodzajPochodzenia, adres, nil
}

// wejscieRozszerzenieObrazu wyjmuje rozszerzenie nazwy obrazu; brak znaczy PNG,
// bo taki jest zapis wyrysów platformy.
func wejscieRozszerzenieObrazu(nazwa string) string {
	wskazanie := strings.LastIndex(nazwa, ".")
	if wskazanie < 0 || wskazanie+1 >= len(nazwa) {
		return "png"
	}
	return nazwa[wskazanie+1:]
}

// wejscieNazwaObrazuZAdresu składa nazwę zasobu z ostatniego człona ścieżki
// adresu obrazu, albo z nazwy zapasowej, gdy adres ścieżki nie niesie.
func wejscieNazwaObrazuZAdresu(adres string) string {
	if rozbity, err := url.Parse(adres); err == nil {
		if nazwa := strings.Trim(rozbity.Path, "/"); nazwa != "" {
			czlony := strings.Split(nazwa, "/")
			return czlony[len(czlony)-1]
		}
	}
	return "obraz ze strony"
}

// wejscieObrazZeSieci pobiera obraz z adresu. Idzie tą samą drogą, którą
// kolejka wczytywania pobiera strony (`net/http`, bez silnika przeglądarki).
func (a *adapterStudia) wejscieObrazZeSieci(ctx context.Context, adres string) ([]byte, error) {
	rozbity, err := url.Parse(strings.TrimSpace(adres))
	if err != nil || rozbity.Host == "" {
		return nil, bladWskazaniaStudio("adres " + adres + " nie jest adresem obrazu")
	}
	if rozbity.Scheme != "http" && rozbity.Scheme != "https" {
		return nil, bladWskazaniaStudio("adres " + adres + " ma schemat " + rozbity.Scheme +
			"; serwer sięga wyłącznie po http i https")
	}
	bajty, _, err := pobierzStroneStudia(ctx, rozbity.String())
	if err != nil {
		return nil, err
	}
	if len(bajty) == 0 {
		return nil, wejscieBladBraku("adres " + adres + " nie oddał żadnych bajtów")
	}
	return bajty, nil
}

// wejscieSekcjaMiejsca oddaje sekcję, w której stoi wskazane miejsce treści,
// albo sekcję ostatnią, gdy miejsce wykracza poza zasięgi wszystkich sekcji.
func wejscieSekcjaMiejsca(forma *shared.StudioDocumentForm, miejsce int) *string {
	for _, sekcja := range forma.Sections {
		if miejsce >= sekcja.RangeStart && miejsce <= sekcja.RangeEnd {
			return wejscieWskaznikTekstu(sekcja.Id)
		}
	}
	if len(forma.Sections) > 0 {
		return wejscieWskaznikTekstu(forma.Sections[len(forma.Sections)-1].Id)
	}
	return nil
}

// wejscieWskaznikZakotwiczenia oddaje wskaźnik na kopię wskazanego rodzaju
// zakotwiczenia obiektu, do pola kontraktu, które wskaźnika wymaga.
func wejscieWskaznikZakotwiczenia(wartosc shared.StudioAnchorKind) *shared.StudioAnchorKind {
	kopia := wartosc
	return &kopia
}

// ── Zapis pod nową nazwą ────────────────────────────────────────────────────

// ZapiszDokumentPodNazwa obsługuje `studio.document.save.as` — zapisuje dokument
// pod nową nazwą albo do wskazanego pliku, wraz z całą postacią przez zapis
// obszaru postaci: arkusz stylów, nastawy strony, sekcje, tabele i obiekty.
func (a *adapterStudia) ZapiszDokumentPodNazwa(ctx context.Context,
	z shared.StudioDocumentSaveAsRequest) (shared.StudioDocumentSaveAsResponse, error) {

	autor, err := wejscieAutorCzynnosci(z.Author)
	if err != nil {
		return shared.StudioDocumentSaveAsResponse{}, err
	}
	if z.Title == nil && z.Path == nil {
		return shared.StudioDocumentSaveAsResponse{}, bladWskazaniaStudio(
			"zapis pod nową nazwą bez nazwy i bez ścieżki pliku — nie ma czego zmienić " +
				"ani gdzie zapisać")
	}
	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDocumentSaveAsResponse{}, err
	}
	if z.Title != nil && strings.TrimSpace(*z.Title) != "" {
		nazwa := strings.TrimSpace(*z.Title)
		stan.dokument.Tytul = &nazwa
	}
	stan.opisCzynnosci = "zapis dokumentu pod nazwą " +
		strings.TrimSpace(wartoscTekstu(stan.dokument.Tytul))

	if err := a.wejscieUtrwalPostac(ctx, stan); err != nil {
		return shared.StudioDocumentSaveAsResponse{}, err
	}

	odpowiedz := shared.StudioDocumentSaveAsResponse{
		Document: a.zlozDokument(stan.dokument),
		Form:     stan.forma,
	}

	if z.Path != nil && strings.TrimSpace(*z.Path) != "" {
		format := shared.StudioExportFormat(shared.StudioExportFormatDocx)
		if z.Format != nil && strings.TrimSpace(string(*z.Format)) != "" {
			format = *z.Format
		}
		wynik, err := a.wydanieDokumentu(ctx, stan, format, z.Path, nil, nil)
		if err != nil {
			return shared.StudioDocumentSaveAsResponse{}, err
		}
		odpowiedz.Export = &wynik
	}

	if z.CreateVersion != nil && *z.CreateVersion {
		dokument, err := a.zalozWersjeDokumentu(ctx, stan.dokument,
			wartoscTekstu(stan.dokument.Tresc), autor, nil)
		if err != nil {
			return shared.StudioDocumentSaveAsResponse{}, err
		}
		stan.dokument = dokument
		odpowiedz.Document = a.zlozDokument(dokument)
	}
	if err := a.wejscieOdlozCzynnosc(ctx, stan, autor); err != nil {
		return shared.StudioDocumentSaveAsResponse{}, err
	}
	return odpowiedz, nil
}

// ── Kopia dokumentu ─────────────────────────────────────────────────────────

// SkopiujDokument obsługuje `studio.document.copy` — zakłada kopię dokumentu,
// osobną od oryginału, wraz z całą postacią i, wedle jawnego wyboru w żądaniu,
// z historią wersji, blokadami fragmentów, znakowaniem i komentarzami.
func (a *adapterStudia) SkopiujDokument(ctx context.Context,
	z shared.StudioDocumentCopyRequest) (shared.StudioDocumentCopyResponse, error) {

	autor, err := wejscieAutorCzynnosci(z.Author)
	if err != nil {
		return shared.StudioDocumentCopyResponse{}, err
	}
	zrodlo, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDocumentCopyResponse{}, err
	}

	okno := strings.TrimSpace(wartoscTekstu(z.WindowId))
	if okno == "" {
		okno = zrodlo.dokument.Okno
	}
	tytul := z.Title
	if tytul == nil {
		nazwa := strings.TrimSpace(wartoscTekstu(zrodlo.dokument.Tytul))
		if nazwa == "" {
			nazwa = "dokument"
		}
		kopia := nazwa + " — kopia"
		tytul = &kopia
	}

	tresc := wartoscTekstu(zrodlo.dokument.Tresc)
	nowy, err := a.repozytorium.ZapiszDokument(ctx, dane.DokumentStudia{
		Kod:                nowyIdentyfikator(przedrostekDokumentuStudio),
		Okno:               okno,
		Tytul:              tytul,
		Format:             zrodlo.dokument.Format,
		Tresc:              &tresc,
		PlikRepozytoriumID: zrodlo.dokument.PlikRepozytoriumID,
	})
	if err != nil {
		return shared.StudioDocumentCopyResponse{}, bladStudio(err)
	}

	stan := &stanPostaci{dokument: nowy, forma: zrodlo.forma}
	stan.forma.DocumentId = nowy.Kod
	wejsciePrzepiszIdentyfikatory(&stan.forma)
	stan.opisCzynnosci = "kopia dokumentu " + zrodlo.dokument.Kod
	if err := a.wejscieUtrwalPostac(ctx, stan); err != nil {
		return shared.StudioDocumentCopyResponse{}, err
	}

	przeniesioneWersje := 0
	if z.IncludeVersions != nil && *z.IncludeVersions {
		wersje, err := a.repozytorium.Wersje(ctx, zrodlo.dokument.ID)
		if err != nil {
			return shared.StudioDocumentCopyResponse{}, bladStudio(err)
		}
		for _, wersja := range wersje {
			trescWersji, err := a.trescZOdwolania(wersja.Tresc, wersja.TrescOdwolanie)
			if err != nil {
				return shared.StudioDocumentCopyResponse{}, err
			}
			if _, err := a.repozytorium.ZapiszWersje(ctx, nowy.ID, dane.WersjaDokumentu{
				Kod:          nowyIdentyfikator(przedrostekWersjiStudio),
				Etykieta:     wersja.Etykieta,
				Podsumowanie: wersja.Podsumowanie,
				Tresc:        &trescWersji,
				Autor:        wersja.Autor,
				KamienMilowy: wersja.KamienMilowy,
			}); err != nil {
				return shared.StudioDocumentCopyResponse{}, bladStudio(err)
			}
			przeniesioneWersje++
		}
	}

	// Blokady fragmentów przechodzą domyślnie; znakowanie i komentarze — nie.
	if z.IncludeLocks == nil || *z.IncludeLocks {
		if err := a.wejsciePrzeniesBlokady(ctx, zrodlo.dokument.ID, nowy.ID, nowy.Kod); err != nil {
			return shared.StudioDocumentCopyResponse{}, err
		}
	}
	if z.IncludeMarkup != nil && *z.IncludeMarkup {
		if err := a.wejsciePrzeniesZnakowanie(ctx, zrodlo.dokument.ID, nowy.ID); err != nil {
			return shared.StudioDocumentCopyResponse{}, err
		}
	}
	if err := a.wejscieOdlozCzynnosc(ctx, stan, autor); err != nil {
		return shared.StudioDocumentCopyResponse{}, err
	}

	return shared.StudioDocumentCopyResponse{
		Document:       a.zlozDokument(stan.dokument),
		Form:           stan.forma,
		VersionsCopied: przeniesioneWersje,
	}, nil
}

// wejsciePrzeniesBlokady przenosi blokady fragmentów do kopii dokumentu, każdą
// pod nowym kodem, wskazującą dokument docelowy.
func (a *adapterStudia) wejsciePrzeniesBlokady(ctx context.Context,
	zrodloID, celID int64, celKod string) error {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return err
	}
	blokady, err := skladnica.BlokadyFragmentow(ctx, zrodloID)
	if err != nil {
		return bladStudio(err)
	}
	for _, blokada := range blokady {
		blokada.ID = 0
		blokada.Kod = nowyIdentyfikator(przedrostekBlokadyStudia)
		blokada.DokumentKod = celKod
		if _, err := skladnica.ZapiszBlokadeFragmentu(ctx, celID, blokada); err != nil {
			return bladStudio(err)
		}
	}
	return nil
}

// wejsciePrzeniesZnakowanie przenosi znakowanie i komentarze do kopii dokumentu,
// każdy wpis pod nowym kodem, wskazujący dokument docelowy.
func (a *adapterStudia) wejsciePrzeniesZnakowanie(ctx context.Context,
	zrodloID, celID int64) error {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return err
	}
	znakowania, err := skladnica.Znakowania(ctx, zrodloID)
	if err != nil {
		return bladStudio(err)
	}
	for _, znakowanie := range znakowania {
		znakowanie.ID = 0
		znakowanie.Kod = nowyIdentyfikator(przedrostekZnakowaniaStudia)
		if _, err := skladnica.ZapiszZnakowanie(ctx, celID, znakowanie); err != nil {
			return bladStudio(err)
		}
	}
	return nil
}

// ── Wniesienie z Biblioteki ─────────────────────────────────────────────────

// WniesZBiblioteki obsługuje `studio.insert.from.library` — wnosi plik, wzór,
// załącznik albo obraz z Biblioteki wprost do dokumentu w miejsce kursora, nie
// do kolejki wczytywania, wraz z obowiązkowym zapisem pochodzenia.
func (a *adapterStudia) WniesZBiblioteki(ctx context.Context,
	z shared.StudioInsertFromLibraryRequest) (shared.StudioInsertFromLibraryResponse, error) {

	autor, err := wejscieAutorCzynnosci(z.Author)
	if err != nil {
		return shared.StudioInsertFromLibraryResponse{}, err
	}
	if strings.TrimSpace(z.LibraryFileId) == "" {
		return shared.StudioInsertFromLibraryResponse{}, bladWskazaniaStudio(
			"wniesienie z Biblioteki bez wskazania pliku")
	}
	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioInsertFromLibraryResponse{}, err
	}
	kodPliku := strings.TrimSpace(z.LibraryFileId)
	bajty, nazwa, err := a.wejscieBajtyZrodla(ctx, wejscieWskazanieZrodla{
		PlikBiblioteki: &kodPliku,
		Czynnosc:       "wniesienie z Biblioteki",
	})
	if err != nil {
		return shared.StudioInsertFromLibraryResponse{}, err
	}

	dlugosc := postacDlugosc(&stan.forma)
	miejsce := z.Offset
	if miejsce < 0 {
		miejsce = 0
	}
	if miejsce > dlugosc {
		miejsce = dlugosc
	}
	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, miejsce, miejsce, autor)
	if len(odcinki) == 0 {
		return shared.StudioInsertFromLibraryResponse{}, bladWskazaniaStudio(
			"wniesienie z Biblioteki zatrzymane przez blokadę fragmentu: " +
				postacNazwaBlokad(pominiete))
	}

	bilans := shared.StudioActionBalance{Applied: 1, Skipped: pominiete}
	koniec := miejsce

	if z.AsObject != nil && *z.AsObject || z.AsAttachment != nil && *z.AsAttachment {
		// Plik wchodzi obiektem osadzonym, zakotwiczonym w miejscu kursora.
		zasob, err := a.odlozTrescStudia(ctx, bajty, nazwa,
			wejscieRozszerzenieObrazu(nazwa), stan.dokument.Okno)
		if err != nil {
			return shared.StudioInsertFromLibraryResponse{}, err
		}
		rodzaj := shared.StudioObjectKind(shared.StudioObjectKindImage)
		if z.AsAttachment != nil && *z.AsAttachment {
			rodzaj = shared.StudioObjectKindTextbox
			bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
				Reason: "załącznik wszedł obiektem osadzonym wskazującym zasób, bo " +
					"kontrakt postaci dokumentu nie zna osobnego rodzaju „załącznik”",
				Detail: wejscieWskaznikTekstu("plik " + nazwa + " jest dostępny pod " +
					"zasobem " + zasob.Id),
			})
		}
		obiekt := shared.StudioDocumentObject{
			Id:            nowyIdentyfikator(przedrostekObiektuStudia),
			Kind:          rodzaj,
			Source:        wejscieWskaznikZrodlaObiektu(shared.StudioObjectSourceLibraryFile),
			AssetId:       wejscieWskaznikTekstu(zasob.Id),
			LibraryFileId: &kodPliku,
			AltText:       wejscieWskaznikTekstu(nazwa),
			AnchorOffset:  wejscieWskaznikCalkowity(miejsce),
			Anchor:        wejscieWskaznikZakotwiczenia(shared.StudioAnchorKindCharacter),
		}
		stan.forma.Objects = append(stan.forma.Objects, obiekt)
		stan.forma.Blocks = append(stan.forma.Blocks, shared.StudioDocumentBlock{
			Id:        nowyIdentyfikator(przedrostekBlokuStudia),
			Kind:      wejscieRodzajBlokuObiekt,
			SectionId: wejscieSekcjaMiejsca(&stan.forma, miejsce),
			ObjectId:  wejscieWskaznikTekstu(obiekt.Id),
		})
	} else {
		format, err := wejscieRozpoznajFormat(nazwa, bajty, nil)
		if err != nil {
			return shared.StudioInsertFromLibraryResponse{}, err
		}
		_, tresc, bilansPliku, err := a.wejscieRozbierzPlik(stan.dokument.Kod, bajty, format, nil)
		if err != nil {
			return shared.StudioInsertFromLibraryResponse{}, err
		}
		bilans.Skipped = append(bilans.Skipped, bilansPliku.Skipped...)
		znaki := []rune(tresc)
		od, do := 0, len(znaki)
		if z.RangeStart != nil && *z.RangeStart > 0 && *z.RangeStart < len(znaki) {
			od = *z.RangeStart
		}
		if z.RangeEnd != nil && *z.RangeEnd > od && *z.RangeEnd <= len(znaki) {
			do = *z.RangeEnd
		}
		wnoszona := string(znaki[od:do])
		autorWpisu := autor
		postacZamienTresc(&stan.forma, miejsce, miejsce, wnoszona, nil, &autorWpisu)
		koniec = miejsce + len([]rune(wnoszona))
	}

	stan.opisCzynnosci = "wniesienie z Biblioteki: " + nazwa
	bilans.Note = wejscieWskaznikTekstu("plik Biblioteki " + nazwa +
		" wniesiony do dokumentu w miejscu " + strconv.Itoa(miejsce))

	if err := a.wejscieUtrwalPostac(ctx, stan); err != nil {
		return shared.StudioInsertFromLibraryResponse{}, err
	}
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKind(shared.StudioChangeKindWstawienie),
		shared.StudioActionKindImportChange, miejsce, koniec, bilans)
	if err != nil {
		return shared.StudioInsertFromLibraryResponse{}, err
	}

	tytulZrodla := nazwa
	pochodzenie, err := a.wejscieOdlozPochodzenie(ctx, stan.dokument,
		shared.StudioProvenanceKindLibraryFile, miejsce, koniec, nil, &kodPliku,
		&tytulZrodla, autor, z.AgentId, z.AgentName)
	if err != nil {
		return shared.StudioInsertFromLibraryResponse{}, err
	}

	return shared.StudioInsertFromLibraryResponse{
		Form:       forma,
		Balance:    bilansGotowy,
		Change:     zmiana,
		ActionId:   stan.czynnosc,
		Provenance: *pochodzenie,
	}, nil
}

// ── Wniesienie ze strony sieci ──────────────────────────────────────────────

// WniesZeSieci obsługuje `studio.insert.from.web` — wnosi fragment albo obraz ze
// strony sieci wprost do dokumentu w miejsce kursora, wraz z zapisem pochodzenia;
// fragment wskazany w oknie przeglądarki ma pierwszeństwo nad pobraniem strony.
func (a *adapterStudia) WniesZeSieci(ctx context.Context,
	z shared.StudioInsertFromWebRequest) (shared.StudioInsertFromWebResponse, error) {

	autor, err := wejscieAutorCzynnosci(z.Author)
	if err != nil {
		return shared.StudioInsertFromWebResponse{}, err
	}
	adres, err := url.Parse(strings.TrimSpace(z.Url))
	if err != nil || adres.Host == "" {
		return shared.StudioInsertFromWebResponse{}, bladWskazaniaStudio(
			"adres " + z.Url + " nie jest adresem strony")
	}
	if adres.Scheme != "http" && adres.Scheme != "https" {
		return shared.StudioInsertFromWebResponse{}, bladWskazaniaStudio(
			"adres " + z.Url + " ma schemat " + adres.Scheme +
				"; serwer sięga wyłącznie po http i https")
	}
	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioInsertFromWebResponse{}, err
	}

	dlugosc := postacDlugosc(&stan.forma)
	miejsce := z.Offset
	if miejsce < 0 {
		miejsce = 0
	}
	if miejsce > dlugosc {
		miejsce = dlugosc
	}
	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, miejsce, miejsce, autor)
	if len(odcinki) == 0 {
		return shared.StudioInsertFromWebResponse{}, bladWskazaniaStudio(
			"wniesienie ze strony zatrzymane przez blokadę fragmentu: " +
				postacNazwaBlokad(pominiete))
	}

	bilans := shared.StudioActionBalance{Applied: 1, Skipped: pominiete}
	koniec := miejsce
	tytulZrodla := adres.Host

	switch {
	case strings.TrimSpace(wartoscTekstu(z.ImageUrl)) != "":
		bajty, err := a.wejscieObrazZeSieci(ctx, wartoscTekstu(z.ImageUrl))
		if err != nil {
			return shared.StudioInsertFromWebResponse{}, err
		}
		nazwa := wejscieNazwaObrazuZAdresu(wartoscTekstu(z.ImageUrl))
		zasob, err := a.odlozTrescStudia(ctx, bajty, nazwa,
			wejscieRozszerzenieObrazu(nazwa), stan.dokument.Okno)
		if err != nil {
			return shared.StudioInsertFromWebResponse{}, err
		}
		obiekt := shared.StudioDocumentObject{
			Id:           nowyIdentyfikator(przedrostekObiektuStudia),
			Kind:         shared.StudioObjectKindImage,
			Source:       wejscieWskaznikZrodlaObiektu(shared.StudioObjectSourceWeb),
			AssetId:      wejscieWskaznikTekstu(zasob.Id),
			SourceUrl:    z.ImageUrl,
			AltText:      wejscieWskaznikTekstu(nazwa),
			AnchorOffset: wejscieWskaznikCalkowity(miejsce),
			Anchor:       wejscieWskaznikZakotwiczenia(shared.StudioAnchorKindCharacter),
		}
		stan.forma.Objects = append(stan.forma.Objects, obiekt)
		stan.forma.Blocks = append(stan.forma.Blocks, shared.StudioDocumentBlock{
			Id:        nowyIdentyfikator(przedrostekBlokuStudia),
			Kind:      wejscieRodzajBlokuObiekt,
			SectionId: wejscieSekcjaMiejsca(&stan.forma, miejsce),
			ObjectId:  wejscieWskaznikTekstu(obiekt.Id),
		})
		bilans.Note = wejscieWskaznikTekstu("obraz ze strony " + adres.Host +
			" osadzony w dokumencie")

	case strings.TrimSpace(wartoscTekstu(z.Text)) != "":
		fragment := strings.TrimSpace(*z.Text)
		autorWpisu := autor
		postacZamienTresc(&stan.forma, miejsce, miejsce, fragment, nil, &autorWpisu)
		koniec = miejsce + len([]rune(fragment))
		bilans.Note = wejscieWskaznikTekstu("fragment wskazany w oknie przeglądarki " +
			"wniesiony do dokumentu wraz z zapisem pochodzenia")

	default:
		bajty, typTresci, err := pobierzStroneStudia(ctx, adres.String())
		if err != nil {
			return shared.StudioInsertFromWebResponse{}, err
		}
		tresc := string(bajty)
		if strings.Contains(strings.ToLower(typTresci), "html") {
			tytulZrodla = wejscieTytulStrony(tresc, adres.Host)
			tresc = tekstZeStronyStudia(tresc)
		}
		tresc = strings.TrimSpace(tresc)
		if tresc == "" {
			return shared.StudioInsertFromWebResponse{}, wejscieBladBraku(
				"strona " + adres.String() + " nie oddała treści tekstowej. Serwer pobiera " +
					"stronę bez silnika przeglądarki, więc strona zbudowana wyłącznie " +
					"skryptem jej tu nie ma; naprawa: wskazać fragment w oknie " +
					"przeglądarki albo sięgnąć po migawkę komendą browser.snapshot.get")
		}
		if strings.TrimSpace(wartoscTekstu(z.Selector)) != "" {
			bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
				Reason: "wskazanie fragmentu strony nie zawęziło treści",
				Detail: wejscieWskaznikTekstu("serwer pobiera stronę bez silnika " +
					"przeglądarki, więc wskazania CSS nie wykonuje; weszła treść całej " +
					"strony. Naprawa: podać fragment polem text z okna przeglądarki"),
			})
		}
		autorWpisu := autor
		postacZamienTresc(&stan.forma, miejsce, miejsce, tresc, nil, &autorWpisu)
		koniec = miejsce + len([]rune(tresc))
		bilans.Note = wejscieWskaznikTekstu("treść strony " + adres.Host +
			" wniesiona do dokumentu wraz z zapisem pochodzenia")
	}

	stan.opisCzynnosci = "wniesienie ze strony " + adres.Host
	if err := a.wejscieUtrwalPostac(ctx, stan); err != nil {
		return shared.StudioInsertFromWebResponse{}, err
	}
	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKind(shared.StudioChangeKindWstawienie),
		shared.StudioActionKindImportChange, miejsce, koniec, bilans)
	if err != nil {
		return shared.StudioInsertFromWebResponse{}, err
	}

	adresZrodla := adres.String()
	pochodzenie, err := a.wejscieOdlozPochodzenie(ctx, stan.dokument,
		shared.StudioProvenanceKindWeb, miejsce, koniec, &adresZrodla, nil,
		&tytulZrodla, autor, z.AgentId, z.AgentName)
	if err != nil {
		return shared.StudioInsertFromWebResponse{}, err
	}

	return shared.StudioInsertFromWebResponse{
		Form:       forma,
		Balance:    bilansGotowy,
		Change:     zmiana,
		ActionId:   stan.czynnosc,
		Provenance: *pochodzenie,
	}, nil
}

// wejscieTytulStrony wyjmuje tytuł strony z jej treści. Tytuł idzie do zapisu
// pochodzenia — bez niego wykaz źródeł pisma niósłby same adresy.
func wejscieTytulStrony(tresc, zapasowy string) string {
	nizej := strings.ToLower(tresc)
	poczatek := strings.Index(nizej, "<title")
	if poczatek < 0 {
		return zapasowy
	}
	poczatek = strings.Index(nizej[poczatek:], ">")
	if poczatek < 0 {
		return zapasowy
	}
	reszta := tresc[poczatek+1:]
	koniec := strings.Index(strings.ToLower(reszta), "</title>")
	if koniec < 0 {
		return zapasowy
	}
	tytul := strings.TrimSpace(reszta[:koniec])
	if tytul == "" {
		return zapasowy
	}
	return tytul
}
