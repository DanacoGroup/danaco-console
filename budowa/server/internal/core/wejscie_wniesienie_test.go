package core

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"

	"golang.org/x/text/encoding/charmap"

	"danacoconsole/shared"
)

// Skutek WEJŚCIA do edytora: czy plik Operatora wchodzi wprost do dokumentu
// z zachowaną postacią, czy zapis znaków jest rozpoznawany i czy kopia dokumentu
// jest osobnym bytem.
//
// Miara jest zawsze taka sama i nie jest kopertą odpowiedzi: po czynności
// dokument czyta się PONOWNIE (`PostacDokumentu`), bo Operator otworzy go
// ponownie, a nie przeczyta odpowiedzi komendy.
//
// Szkody, które ten plik ma wykluczyć:
//  1. plik wniesiony jako treść płaska, z arkuszem stylów i tabelą zgubionymi
//     po drodze — czyli postać ginąca przy wniesieniu;
//  2. plik w stronie kodowej innej niż UTF-8 wczytany jako krzaczki, bez ani
//     jednego zdania o tym w bilansie;
//  3. kopia dokumentu założona jako drugie odwołanie do tego samego bytu, po
//     której zmiana w kopii rusza oryginał;
//  4. wniesienie oddane jako udane, a bez zapisu pochodzenia — po tygodniu nikt
//     nie odtworzy, na czym pismo się opiera.

// wejscieUprzazSprawdzianu składa adapter modułu Studia wraz z bazą.
func wejscieUprzazSprawdzianu(t *testing.T) (*adapterStudia, context.Context) {
	t.Helper()
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	adapter := nowyAdapterStudia(zmontowany.dane.Studio).
		ZZasobami(zmontowany.dane.Design, katalog)
	return adapter, zycie
}

// wejscieDokumentWzorcowy składa postać pisma wzorcowego: dwa akapity, nagłówek,
// styl własny i tabelę o policzonych szerokościach kolumn.
func wejscieDokumentWzorcowy(kod string) shared.StudioDocumentForm {
	nastawy := wejscieDomyslneNastawyStrony("A4", nil)
	sekcja := shared.StudioSection{
		Id:        "studio-sek-wzorcowa",
		Index:     0,
		Start:     wejscieWskaznikPoczatkuSekcji(shared.StudioSectionStartContinuous),
		PageSetup: &nastawy,
		HeadersFooters: []shared.StudioHeaderFooter{{
			Scope:      shared.StudioHeaderScopeDefault,
			HeaderText: wejscieWskaznikTekstu("Danaco Group — pismo wzorcowe"),
			FooterText: wejscieWskaznikTekstu("strona"),
		}},
	}
	arkusz := append(wejscieDomyslnyArkuszStylow(), shared.StudioNamedStyle{
		Name:        "podstawa prawna",
		DisplayName: wejscieWskaznikTekstu("Podstawa prawna"),
		Kind:        shared.StudioStyleKindParagraph,
		BasedOn:     wejscieWskaznikTekstu(wejscieStylTekstZasadniczy),
		Character: &shared.StudioCharacterFormat{
			Bold:       wejscieWskaznikLogiczny(true),
			FontSizePt: wejscieWskaznikRzeczywisty(11),
		},
	})
	tabela := shared.StudioDocumentTable{
		Id:             "studio-tab-wzorcowa",
		Rows:           2,
		Columns:        2,
		ColumnWidthsMm: []float64{80, 80},
		WidthMm:        wejscieWskaznikRzeczywisty(160),
		HeaderRows:     wejscieWskaznikCalkowity(1),
		Cells: []shared.StudioTableCell{
			{Row: 0, Column: 0, Text: wejscieWskaznikTekstu("Pozycja")},
			{Row: 0, Column: 1, Text: wejscieWskaznikTekstu("Kwota")},
			{Row: 1, Column: 0, Text: wejscieWskaznikTekstu("Wynagrodzenie")},
			{Row: 1, Column: 1, Text: wejscieWskaznikTekstu("4 200,00")},
		},
	}
	naglowek := wejscieBlokAkapitu(sekcja.Id, "Umowa najmu", wejscieStylNaglowek1, nil)
	naglowek.Paragraph.OutlineLevel = wejscieWskaznikCalkowity(1)
	podstawa := wejscieBlokAkapitu(sekcja.Id,
		"Zawarta na podstawie przepisów prawa cywilnego.", "podstawa prawna", nil)
	tresc := wejscieBlokAkapitu(sekcja.Id,
		"Strony ustalają wynagrodzenie zgodnie z tabelą.", wejscieStylTekstZasadniczy, nil)

	return shared.StudioDocumentForm{
		DocumentId: kod,
		PageSetup:  &nastawy,
		Styles:     arkusz,
		Sections:   []shared.StudioSection{sekcja},
		Blocks: []shared.StudioDocumentBlock{naglowek, podstawa, tresc, {
			Id:        "studio-blok-tabela",
			Kind:      wejscieRodzajBlokuTabela,
			SectionId: wejscieWskaznikTekstu(sekcja.Id),
			TableId:   wejscieWskaznikTekstu(tabela.Id),
		}},
		Tables: []shared.StudioDocumentTable{tabela},
	}
}

// TestWejscieDocxZachowujePostacPoWniesieniu mierzy obieg pełny: postać pisma
// złożona do `.docx`, wniesiona z powrotem i odczytana Z BAZY.
//
// Porównanie idzie po ODCZYCIE, nie po bajtach — tak stanowi zlecenie: dokument
// wczytany z `.docx` i oddany z powrotem zachowuje styl, sekcje i tabele.
func TestWejscieDocxZachowujePostacPoWniesieniu(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	wzorcowa := wejscieDokumentWzorcowy("")
	tresc := wejscieTrescZPostaci(&wzorcowa)
	bajty, _, err := wejscieZlozOoxml(&wzorcowa, tresc, "Umowa najmu", false)
	if err != nil {
		t.Fatalf("złożenie pliku docx odmówiło: %v", err)
	}

	zapis := base64.StdEncoding.EncodeToString(bajty)
	wniesiony, err := adapter.WniesPlikDoEdytora(zycie, shared.StudioDocumentImportFileRequest{
		WindowId:    "okno-wejscie",
		BytesBase64: &zapis,
		Path:        wskaznik("umowa najmu.docx"),
	})
	if err != nil {
		t.Fatalf("wniesienie pliku docx odmówiło: %v", err)
	}
	if wniesiony.Balance.Format != shared.StudioImportFormatDocx {
		t.Fatalf("bilans nazwał format %q, a plik był docx", wniesiony.Balance.Format)
	}
	if wniesiony.Provenance == nil || wniesiony.Provenance.Kind !=
		shared.StudioProvenanceKindImportedFile {
		t.Fatal("wniesienie pliku nie odłożyło zapisu pochodzenia — bez niego nikt " +
			"nie odtworzy, na czym pismo się opiera")
	}
	if wniesiony.Document.Title == nil || *wniesiony.Document.Title != "umowa najmu" {
		t.Fatalf("dokument nie przejął nazwy z pliku: %v", wniesiony.Document.Title)
	}

	// Miara właściwa: odczyt Z BAZY, nie z odpowiedzi.
	odczytana, err := adapter.PostacDokumentu(zycie, shared.StudioDocumentFormGetRequest{
		DocumentId: wniesiony.Document.Id,
	})
	if err != nil {
		t.Fatalf("odczyt postaci po wniesieniu odmówił: %v", err)
	}
	postac := odczytana.Form

	if postac.PageSetup == nil || postac.PageSetup.MarginLeft == nil ||
		*postac.PageSetup.MarginLeft != 25 {
		t.Fatalf("nastawy strony nie przeszły przez wniesienie: %#v", postac.PageSetup)
	}
	if len(postac.Sections) == 0 {
		t.Fatal("sekcje nie przeszły przez wniesienie")
	}
	stylWlasny := false
	for _, styl := range postac.Styles {
		if styl.Name == "podstawa prawna" {
			stylWlasny = true
			if styl.Character == nil || styl.Character.Bold == nil || !*styl.Character.Bold {
				t.Fatal("styl własny przeszedł bez swojej postaci znaku")
			}
		}
	}
	if !stylWlasny {
		t.Fatalf("arkusz stylów zgubił styl własny; przeszło %d stylów", len(postac.Styles))
	}
	if len(postac.Tables) != 1 {
		t.Fatalf("tabela nie przeszła przez wniesienie; tabel: %d", len(postac.Tables))
	}
	tabela := postac.Tables[0]
	if tabela.Rows != 2 || tabela.Columns != 2 {
		t.Fatalf("tabela przeszła w rozmiarze %dx%d, a była 2x2", tabela.Rows, tabela.Columns)
	}
	if len(tabela.ColumnWidthsMm) != 2 || tabela.ColumnWidthsMm[0] <= 0 {
		t.Fatalf("szerokości kolumn tabeli są zerowe po wniesieniu: %v",
			tabela.ColumnWidthsMm)
	}
	if !strings.Contains(wejscieTrescZPostaci(&postac), "Umowa najmu") {
		t.Fatalf("treść nie przeszła przez obieg docx: %q", wejscieTrescZPostaci(&postac))
	}
}

// TestWejscieOdtZachowujePostac mierzy obieg OpenDocument: postać złożona do
// `.odt` i odczytana z powrotem.
func TestWejscieOdtZachowujePostac(t *testing.T) {
	wzorcowa := wejscieDokumentWzorcowy("")
	tresc := wejscieTrescZPostaci(&wzorcowa)
	bajty, _, err := wejscieZlozOdf(&wzorcowa, tresc, "Umowa najmu", false)
	if err != nil {
		t.Fatalf("złożenie pliku odt odmówiło: %v", err)
	}

	format, err := wejscieRozpoznajFormat("umowa.odt", bajty, nil)
	if err != nil {
		t.Fatalf("rozpoznanie formatu odmówiło: %v", err)
	}
	if format != shared.StudioImportFormatOdt {
		t.Fatalf("rozpoznano format %q, a plik był odt", format)
	}

	postac, trescOdczytana, bilans, err := wejscieCzytajOdf("studio-dok-1", bajty, format)
	if err != nil {
		t.Fatalf("odczyt pliku odt odmówił: %v", err)
	}
	if !strings.Contains(trescOdczytana, "Umowa najmu") {
		t.Fatalf("treść nie przeszła przez obieg odt: %q", trescOdczytana)
	}
	if postac.PageSetup == nil || postac.PageSetup.MarginTop == nil ||
		*postac.PageSetup.MarginTop != 25 {
		t.Fatalf("nastawy strony nie przeszły przez obieg odt: %#v", postac.PageSetup)
	}
	stylWlasny := false
	for _, styl := range postac.Styles {
		if styl.Name == "podstawa prawna" {
			stylWlasny = true
		}
	}
	if !stylWlasny {
		t.Fatalf("arkusz stylów nie przeszedł przez obieg odt; stylów: %d", len(postac.Styles))
	}
	if len(postac.Tables) != 1 || postac.Tables[0].Columns != 2 {
		t.Fatalf("tabela nie przeszła przez obieg odt: %#v", postac.Tables)
	}
	if len(postac.Tables[0].ColumnWidthsMm) != 2 ||
		postac.Tables[0].ColumnWidthsMm[0] <= 0 {
		t.Fatalf("szerokości kolumn po obiegu odt są zerowe: %v",
			postac.Tables[0].ColumnWidthsMm)
	}
	naglowekPrzeszedl := false
	for _, sekcja := range postac.Sections {
		for _, wpis := range sekcja.HeadersFooters {
			if wpis.HeaderText != nil && strings.Contains(*wpis.HeaderText, "Danaco") {
				naglowekPrzeszedl = true
			}
		}
	}
	if !naglowekPrzeszedl {
		t.Fatal("nagłówek strony nie przeszedł przez obieg odt — szablon ma nieść " +
			"nagłówek i stopkę razem z arkuszem stylów")
	}
	if bilans.StylesRecovered == nil || *bilans.StylesRecovered == 0 {
		t.Fatal("bilans wniesienia nie policzył przejętych stylów")
	}
}

// TestWejscieRozpoznajeStroneKodowa mierzy, że pismo w Windows-1250 wchodzi jako
// polskie litery, a bilans NAZYWA rozpoznany zapis znaków.
func TestWejscieRozpoznajeStroneKodowa(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	zdanie := "Zażółć gęślą jaźń — pismo urzędowe."
	przelozone, err := charmap.Windows1250.NewEncoder().Bytes([]byte(zdanie))
	if err != nil {
		t.Fatalf("nie można złożyć pliku sprawdzianu: %v", err)
	}
	zapis := base64.StdEncoding.EncodeToString(przelozone)

	wniesiony, err := adapter.WniesPlikDoEdytora(zycie, shared.StudioDocumentImportFileRequest{
		WindowId:    "okno-kodowanie",
		BytesBase64: &zapis,
		Path:        wskaznik("pismo.txt"),
	})
	if err != nil {
		t.Fatalf("wniesienie pliku tekstowego odmówiło: %v", err)
	}
	if wniesiony.Balance.Encoding == nil {
		t.Fatal("bilans nie nazwał rozpoznanego zapisu znaków — Operator ma wiedzieć, " +
			"jak rdzeń odczytał jego plik")
	}
	if !strings.Contains(strings.ToLower(*wniesiony.Balance.Encoding), "1250") {
		t.Fatalf("rozpoznano zapis %q, a plik był w Windows-1250", *wniesiony.Balance.Encoding)
	}

	odczytany, err := adapter.repozytorium.Dokument(zycie, wniesiony.Document.Id)
	if err != nil {
		t.Fatalf("odczyt dokumentu z bazy odmówił: %v", err)
	}
	if !strings.Contains(wartoscTekstu(odczytany.Tresc), "Zażółć gęślą jaźń") {
		t.Fatalf("polskie litery nie przeszły przez wniesienie: %q",
			wartoscTekstu(odczytany.Tresc))
	}
}

// TestWejscieKopiaJestOsobnymDokumentem mierzy to, o co prosi zlecenie wprost:
// zmiana w kopii NIE rusza oryginału.
func TestWejscieKopiaJestOsobnymDokumentem(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	zalozony, err := adapter.ZalozDokument(zycie, shared.StudioDocumentCreateRequest{
		WindowId: "okno-kopia",
		Title:    wskaznik("Pismo pierwotne"),
	})
	if err != nil {
		t.Fatalf("założenie dokumentu odmówiło: %v", err)
	}
	if _, err := adapter.ZmienTresc(zycie, shared.StudioTextEditRequest{
		DocumentId: zalozony.Document.Id,
		RangeStart: 0,
		RangeEnd:   0,
		Text:       "Treść oryginału.",
	}); err != nil {
		t.Fatalf("zapis treści oryginału odmówił: %v", err)
	}

	kopia, err := adapter.SkopiujDokument(zycie, shared.StudioDocumentCopyRequest{
		DocumentId: zalozony.Document.Id,
	})
	if err != nil {
		t.Fatalf("kopia dokumentu odmówiła: %v", err)
	}
	if kopia.Document.Id == zalozony.Document.Id {
		t.Fatal("kopia dostała identyfikator oryginału — to nie kopia, a drugie " +
			"odwołanie do tego samego bytu")
	}

	if _, err := adapter.ZmienTresc(zycie, shared.StudioTextEditRequest{
		DocumentId: kopia.Document.Id,
		RangeStart: 0,
		RangeEnd:   0,
		Text:       "Zmiana wyłącznie w kopii. ",
	}); err != nil {
		t.Fatalf("zmiana treści kopii odmówiła: %v", err)
	}

	oryginal, err := adapter.repozytorium.Dokument(zycie, zalozony.Document.Id)
	if err != nil {
		t.Fatalf("odczyt oryginału odmówił: %v", err)
	}
	if strings.Contains(wartoscTekstu(oryginal.Tresc), "wyłącznie w kopii") {
		t.Fatalf("zmiana w kopii ruszyła oryginał; treść oryginału: %q",
			wartoscTekstu(oryginal.Tresc))
	}
	postacKopii, err := adapter.PostacDokumentu(zycie, shared.StudioDocumentFormGetRequest{
		DocumentId: kopia.Document.Id,
	})
	if err != nil {
		t.Fatalf("odczyt postaci kopii odmówił: %v", err)
	}
	if len(postacKopii.Form.Sections) == 0 {
		t.Fatal("kopia nie dostała sekcji — postać nie przeszła do kopii")
	}
	for _, sekcja := range postacKopii.Form.Sections {
		for _, oryginalna := range zalozony.Form.Sections {
			if sekcja.Id == oryginalna.Id {
				t.Fatalf("sekcja kopii ma identyfikator sekcji oryginału (%s) — jeden "+
					"wiersz widziany z dwóch dokumentów", sekcja.Id)
			}
		}
	}
}

// TestWejscieObrazWchodziWMiejsceKursora mierzy, że obraz wnosi się wprost do
// dokumentu, z bajtami odłożonymi w magazynie zasobów rdzenia.
func TestWejscieObrazWchodziWMiejsceKursora(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	zalozony, err := adapter.ZalozDokument(zycie, shared.StudioDocumentCreateRequest{
		WindowId: "okno-obraz",
		Title:    wskaznik("Pismo z logo"),
	})
	if err != nil {
		t.Fatalf("założenie dokumentu odmówiło: %v", err)
	}
	// Najmniejszy poprawny PNG — jeden piksel. Bajty są tu treścią sprawdzianu,
	// bo mierzymy odłożenie ich w magazynie, a nie rozbiór obrazu.
	png, err := base64.StdEncoding.DecodeString(
		"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8DwHwAFAAH/q842iQAAAABJRU5ErkJggg==")
	if err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	zapis := base64.StdEncoding.EncodeToString(png)

	wniesiony, err := adapter.WniesObraz(zycie, shared.StudioDocumentImageImportRequest{
		DocumentId:  zalozony.Document.Id,
		Offset:      0,
		Source:      shared.StudioObjectSourceFile,
		BytesBase64: &zapis,
		Path:        wskaznik("logo.png"),
		AltText:     wskaznik("logo Danaco Group"),
	})
	if err != nil {
		t.Fatalf("wniesienie obrazu odmówiło: %v", err)
	}
	if wniesiony.Object.AssetId == nil {
		t.Fatal("obraz wszedł bez zasobu magazynu — okno nie miałoby czego pokazać")
	}
	if wniesiony.Balance.Applied != 1 {
		t.Fatalf("bilans czynności mówi o %d zmianach, a obraz wszedł raz",
			wniesiony.Balance.Applied)
	}

	odczytana, err := adapter.PostacDokumentu(zycie, shared.StudioDocumentFormGetRequest{
		DocumentId: zalozony.Document.Id,
	})
	if err != nil {
		t.Fatalf("odczyt postaci po wniesieniu obrazu odmówił: %v", err)
	}
	znaleziony := false
	for _, obiekt := range odczytana.Form.Objects {
		if obiekt.Id == wniesiony.Object.Id && obiekt.AssetId != nil {
			znaleziony = true
		}
	}
	if !znaleziony {
		t.Fatalf("obraz nie stoi w postaci dokumentu po ponownym odczycie; obiektów: %d",
			len(odczytana.Form.Objects))
	}
}

// TestWejsciePdfNiesieBilansOdzyskanego mierzy zasadę rozstrzygającą: odzyskanie
// z PDF jest ODTWORZENIEM i odpowiedź ma to powiedzieć liczbami.
func TestWejsciePdfNiesieBilansOdzyskanego(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	bajty, err := dokumentPdfZTekstu(strings.Join([]string{
		"Umowa najmu",
		"Strony ustalaja wynagrodzenie miesieczne.",
		"Termin zaplaty przypada na piaty dzien miesiaca.",
	}, "\n"))
	if err != nil {
		t.Fatalf("nie można złożyć dokumentu PDF sprawdzianu: %v", err)
	}
	zapis := base64.StdEncoding.EncodeToString(bajty)

	wynik, err := adapter.WniesPdfDoEdytora(zycie, shared.StudioDocumentImportPdfRequest{
		WindowId:    "okno-pdf",
		BytesBase64: &zapis,
		Path:        wskaznik("umowa.pdf"),
	})
	if err != nil {
		t.Fatalf("zamiana PDF na dokument edytowalny odmówiła: %v", err)
	}
	bilans := wynik.Balance
	if bilans.Pages == nil || *bilans.Pages == 0 {
		t.Fatal("bilans nie policzył stron dokumentu PDF")
	}
	if bilans.PagesWithText == nil || *bilans.PagesWithText == 0 {
		t.Fatalf("bilans mówi, że żadna strona nie miała warstwy tekstowej, a dokument " +
			"powstał z tekstu")
	}
	if bilans.ParagraphsRecovered == nil || *bilans.ParagraphsRecovered == 0 {
		t.Fatal("bilans nie policzył odtworzonych akapitów")
	}
	if bilans.Note == nil || !strings.Contains(*bilans.Note, "ODTWORZENIEM") {
		t.Fatalf("bilans nie mówi wprost, że odzyskanie jest odtworzeniem: %v", bilans.Note)
	}
	if len(bilans.Skipped) == 0 {
		t.Fatal("bilans nie wymienił ani jednej cechy, której PDF nie niesie — " +
			"cisza w tym miejscu jest zakazana")
	}

	odczytany, err := adapter.repozytorium.Dokument(zycie, wynik.Document.Id)
	if err != nil {
		t.Fatalf("odczyt dokumentu po zamianie odmówił: %v", err)
	}
	if !strings.Contains(wartoscTekstu(odczytany.Tresc), "Umowa najmu") {
		t.Fatalf("treść PDF nie weszła do dokumentu: %q", wartoscTekstu(odczytany.Tresc))
	}
}
