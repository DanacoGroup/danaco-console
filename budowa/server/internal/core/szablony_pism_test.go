package core

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Skutek WARSZTATU SZABLONÓW: czy szablon niesie NARAZ arkusz stylów, nastawy
// strony, nagłówek, stopkę i pola do wypełnienia, czy szablonu fabrycznego nie
// da się usunąć i czy pole wymagane bez wartości nie znika po cichu z pisma.
//
// Szkody, które ten plik ma wykluczyć:
//  1. szablon zapisujący samą treść — po nim dokument z wzoru pisma wychodzi bez
//     papieru firmowego, czyli wzór nie jest wzorem;
//  2. usunięcie szablonu fabrycznego, po którym wykazu nie da się odtworzyć bez
//     ponownego wdrożenia;
//  3. pole wymagane bez wartości usunięte z treści — pismo wygląda na kompletne,
//     a nie jest;
//  4. wypełnienie pól zamianą w napisie treści, po której postać wzorcowa pisma
//     (kroje, wcięcia, granice akapitów) przepada.

// TestSzablonZDokumentuNiesiePostacIPola mierzy wymaganie Właściciela wprost:
// szablon niesie arkusz stylów, nastawy strony, nagłówek, stopkę i pola NARAZ.
func TestSzablonZDokumentuNiesiePostacIPola(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	postac := wejscieDokumentWzorcowy("")
	// Znacznik pola w treści pisma wzorcowego — tak Operator zaznacza miejsce
	// do podstawienia.
	postac.Blocks[2].Runs[0].Text = "Adresat: {{adresat}}, dnia {{data}}."
	tresc := wejscieTrescZPostaci(&postac)
	plik, _, err := wejscieZlozOoxml(&postac, tresc, "Wzór pisma", false)
	if err != nil {
		t.Fatalf("złożenie pliku wzorcowego odmówiło: %v", err)
	}
	dokument := wydanieWniesDokument(t, adapter, zycie, plik, "okno-szablon")

	pola, err := json.Marshal([]shared.StudioTemplateFieldSpec{
		{
			Name:     "adresat",
			Label:    "Adresat pisma",
			Kind:     shared.StudioTemplateFieldKindText,
			Required: true,
		},
		{
			Name:  "data",
			Label: "Data pisma",
			Kind:  shared.StudioTemplateFieldKindDate,
		},
	})
	if err != nil {
		t.Fatalf("nie można złożyć wykazu pól: %v", err)
	}

	zapisany, err := adapter.ZapiszSzablonPisma(zycie, shared.StudioTemplateSaveRequest{
		DocumentId: &dokument,
		Name:       "Pismo urzędowe Danaco",
		Category:   wskaznik("pisma"),
		Fields:     pola,
	})
	if err != nil {
		t.Fatalf("zapis szablonu z dokumentu odmówił: %v", err)
	}

	// Miara właściwa: odczyt szablonu OSOBNYM wywołaniem, nie z odpowiedzi.
	odczytany, err := adapter.szablonPismaSzczegol(zycie, zapisany.Template.Id)
	if err != nil {
		t.Fatalf("odczyt szablonu odmówił: %v", err)
	}
	if odczytany.Form == nil {
		t.Fatal("szablon nie niesie postaci wzorcowej — wzór pisma bez arkusza stylów " +
			"i nastaw strony nie jest wzorem")
	}
	if odczytany.Form.PageSetup == nil || odczytany.Form.PageSetup.MarginLeft == nil {
		t.Fatalf("szablon nie przejął nastaw strony: %#v", odczytany.Form.PageSetup)
	}
	stylWlasny := false
	for _, styl := range odczytany.Form.Styles {
		if styl.Name == "podstawa prawna" {
			stylWlasny = true
		}
	}
	if !stylWlasny {
		t.Fatalf("szablon nie przejął arkusza stylów; stylów: %d",
			len(odczytany.Form.Styles))
	}
	naglowek := false
	for _, sekcja := range odczytany.Form.Sections {
		for _, wpis := range sekcja.HeadersFooters {
			if wpis.HeaderText != nil && strings.Contains(*wpis.HeaderText, "Danaco") {
				naglowek = true
			}
		}
	}
	if !naglowek {
		t.Fatal("szablon nie przejął nagłówka strony")
	}
	if len(odczytany.Form.Tables) != 1 {
		t.Fatalf("szablon nie przejął tabeli wzorcowej; tabel: %d",
			len(odczytany.Form.Tables))
	}

	wykaz, err := adapter.PolaSzablonu(zycie, shared.StudioTemplateFieldListRequest{
		TemplateId: zapisany.Template.Id,
	})
	if err != nil {
		t.Fatalf("wykaz pól szablonu odmówił: %v", err)
	}
	if len(wykaz.Fields) != 2 {
		t.Fatalf("szablon niesie %d pól, a wskazano dwa: %#v", len(wykaz.Fields),
			wykaz.Fields)
	}
	for _, pole := range wykaz.Fields {
		if pole.Name == "adresat" {
			if !pole.Required {
				t.Fatal("pole wymagane wyszło jako nieobowiązkowe")
			}
			if pole.AnchorOffset == nil {
				t.Fatal("pole nie ma miejsca w treści — okno nie miałoby gdzie postawić " +
					"kursora przy wypełnianiu")
			}
		}
	}
}

// TestSzablonFabrycznegoNieDaSieUsunac mierzy odmowę NAZWANĄ, wzorem
// `studio.operation.delete` dla operacji fabrycznych.
func TestSzablonFabrycznegoNieDaSieUsunac(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	skladnica, err := adapter.wejscieSkladnica()
	if err != nil {
		t.Fatalf("repozytorium warsztatu szablonów nie jest dostępne: %v", err)
	}
	wykaz, err := skladnica.SzablonyWarsztatu(zycie, "")
	if err != nil {
		t.Fatalf("wykaz szablonów odmówił: %v", err)
	}
	fabryczny := ""
	for _, szablon := range wykaz {
		if szablon.Fabryczny {
			fabryczny = szablon.Kod
			break
		}
	}
	if fabryczny == "" {
		t.Fatal("baza nie niesie ani jednego szablonu fabrycznego — wykaz fabryczny " +
			"z migracji 131 miał je wnieść")
	}

	_, err = adapter.UsunSzablonPisma(zycie, shared.StudioTemplateDeleteRequest{
		TemplateId: fabryczny,
	})
	if err == nil {
		t.Fatal("szablon fabryczny dał się usunąć — wykazu fabrycznego nie da się " +
			"odtworzyć bez ponownego wdrożenia")
	}
	if !strings.Contains(err.Error(), "fabryczny") {
		t.Fatalf("odmowa nie nazywa powodu: %v", err)
	}

	// Szablon własny usuwa się bez przeszkód — odmowa dotyczy fabrycznego, nie
	// warsztatu jako takiego.
	wlasny, err := adapter.ZapiszSzablonPisma(zycie, shared.StudioTemplateSaveRequest{
		Name: "Wzór własny do usunięcia",
	})
	if err != nil {
		t.Fatalf("zapis szablonu własnego odmówił: %v", err)
	}
	usuniety, err := adapter.UsunSzablonPisma(zycie, shared.StudioTemplateDeleteRequest{
		TemplateId: wlasny.Template.Id,
	})
	if err != nil {
		t.Fatalf("usunięcie szablonu własnego odmówiło: %v", err)
	}
	if !usuniety.Deleted {
		t.Fatal("usunięcie szablonu własnego oddało fałsz bez odmowy — to jest cisza")
	}
	if _, err := adapter.szablonPismaSzczegol(zycie, wlasny.Template.Id); err == nil {
		t.Fatal("szablon usunięty nadal się odczytuje — wiersz nie zszedł")
	}
}

// TestSzablonPoleZmianaIUsuniecie mierzy warsztat pól: dołożenie, zmianę
// i zdjęcie pola tą samą drogą.
func TestSzablonPoleZmianaIUsuniecie(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	szablon, err := adapter.ZapiszSzablonPisma(zycie, shared.StudioTemplateSaveRequest{
		Name: "Wzór z polami",
	})
	if err != nil {
		t.Fatalf("zapis szablonu odmówił: %v", err)
	}

	ustawiony, err := adapter.UstawPoleSzablonu(zycie, shared.StudioTemplateFieldSetRequest{
		TemplateId:   szablon.Template.Id,
		Name:         "sygnatura",
		Label:        wskaznik("Sygnatura sprawy"),
		Kind:         wskaznik(shared.StudioTemplateFieldKind(shared.StudioTemplateFieldKindText)),
		Required:     wskaznik(true),
		DefaultValue: wskaznik("bez sygnatury"),
	})
	if err != nil {
		t.Fatalf("ustawienie pola szablonu odmówiło: %v", err)
	}
	if len(ustawiony.Fields) != 1 || ustawiony.Fields[0].Name != "sygnatura" {
		t.Fatalf("pole nie weszło do wykazu: %#v", ustawiony.Fields)
	}

	// Pole wyboru bez wykazu wartości jest odmawiane — nie byłoby z czego wybrać.
	_, err = adapter.UstawPoleSzablonu(zycie, shared.StudioTemplateFieldSetRequest{
		TemplateId: szablon.Template.Id,
		Name:       "rodzaj pisma",
		Kind:       wskaznik(shared.StudioTemplateFieldKind(shared.StudioTemplateFieldKindChoice)),
	})
	if err == nil {
		t.Fatal("pole wyboru bez wykazu wartości weszło — nie byłoby z czego wybrać")
	}

	zdjety, err := adapter.UstawPoleSzablonu(zycie, shared.StudioTemplateFieldSetRequest{
		TemplateId: szablon.Template.Id,
		Name:       "sygnatura",
		Remove:     wskaznik(true),
	})
	if err != nil {
		t.Fatalf("zdjęcie pola szablonu odmówiło: %v", err)
	}
	if len(zdjety.Fields) != 0 {
		t.Fatalf("pole nie zszedło z wykazu: %#v", zdjety.Fields)
	}
}

// TestSzablonWypelnieniePolNieUkrywaBrakow mierzy zasadę: pole wymagane bez
// wartości zostaje w treści WIDOCZNE, a odpowiedź je wymienia.
func TestSzablonWypelnieniePolNieUkrywaBrakow(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	postac := wejscieDokumentWzorcowy("")
	postac.Blocks[2].Runs[0].Text = "Adresat: {{adresat}}. Kwota: {{kwota}} złotych."
	tresc := wejscieTrescZPostaci(&postac)
	plik, _, err := wejscieZlozOoxml(&postac, tresc, "Wzór pisma", false)
	if err != nil {
		t.Fatalf("złożenie pliku wzorcowego odmówiło: %v", err)
	}
	dokument := wydanieWniesDokument(t, adapter, zycie, plik, "okno-wypelnienie")

	pola, err := json.Marshal([]shared.StudioTemplateFieldSpec{
		{Name: "adresat", Label: "Adresat", Kind: shared.StudioTemplateFieldKindText,
			Required: true},
		{Name: "kwota", Label: "Kwota", Kind: shared.StudioTemplateFieldKindNumber,
			Required: true},
	})
	if err != nil {
		t.Fatalf("nie można złożyć wykazu pól: %v", err)
	}
	szablon, err := adapter.ZapiszSzablonPisma(zycie, shared.StudioTemplateSaveRequest{
		DocumentId: &dokument,
		Name:       "Wzór z polami wymaganymi",
		Fields:     pola,
	})
	if err != nil {
		t.Fatalf("zapis szablonu odmówił: %v", err)
	}

	wartosci, err := json.Marshal(map[string]any{"kwota": 4200})
	if err != nil {
		t.Fatalf("nie można złożyć wartości pól: %v", err)
	}
	wypelniony, err := adapter.WypelnijSzablon(zycie, shared.StudioTemplateFillRequest{
		TemplateId: szablon.Template.Id,
		WindowId:   wskaznik("okno-wypelnienie"),
		Values:     wartosci,
		Title:      wskaznik("Pismo do Kowalskiego"),
	})
	if err != nil {
		t.Fatalf("wypełnienie szablonu odmówiło: %v", err)
	}
	if len(wypelniony.MissingRequired) != 1 || wypelniony.MissingRequired[0] != "adresat" {
		t.Fatalf("odpowiedź nie wymienia pola wymaganego bez wartości: %#v",
			wypelniony.MissingRequired)
	}
	if wypelniony.Balance.SkippedCount == 0 {
		t.Fatal("bilans wypełnienia nie policzył pól niewypełnionych")
	}

	// Miara właściwa: treść dokumentu odczytana Z BAZY.
	odczytany, err := adapter.repozytorium.Dokument(zycie, wypelniony.Document.Id)
	if err != nil {
		t.Fatalf("odczyt dokumentu wypełnionego odmówił: %v", err)
	}
	trescGotowa := wartoscTekstu(odczytany.Tresc)
	if !strings.Contains(trescGotowa, "4200") {
		t.Fatalf("wartość liczbowa nie weszła do pisma: %q", trescGotowa)
	}
	if !strings.Contains(trescGotowa, "{{adresat}}") {
		t.Fatalf("znacznik pola wymaganego zniknął z pisma — dokument z pustym miejscem "+
			"wygląda na kompletny, a nie jest: %q", trescGotowa)
	}

	// Postać wzorcowa MA przejść do dokumentu wypełnionego — to jest sens
	// szablonu. Zamiana w napisie treści by ją zgubiła.
	postacGotowa, err := adapter.PostacDokumentu(zycie, shared.StudioDocumentFormGetRequest{
		DocumentId: wypelniony.Document.Id,
	})
	if err != nil {
		t.Fatalf("odczyt postaci dokumentu wypełnionego odmówił: %v", err)
	}
	stylWlasny := false
	for _, styl := range postacGotowa.Form.Styles {
		if styl.Name == "podstawa prawna" {
			stylWlasny = true
		}
	}
	if !stylWlasny {
		t.Fatalf("dokument z szablonu nie dostał arkusza stylów wzorcowego; stylów: %d",
			len(postacGotowa.Form.Styles))
	}
	if len(postacGotowa.Form.Tables) != 1 {
		t.Fatalf("dokument z szablonu nie dostał tabeli wzorcowej; tabel: %d",
			len(postacGotowa.Form.Tables))
	}
}

// TestSzablonObiegPlikuOperatora mierzy dwie czynności naraz: oddanie szablonu do
// pliku i wniesienie go z powrotem. Pola mają się odtworzyć ze znaczników treści,
// bo żaden format biurowy nie niesie wykazu pól platformy.
func TestSzablonObiegPlikuOperatora(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	postac := wejscieDokumentWzorcowy("")
	postac.Blocks[2].Runs[0].Text = "Adresat: {{adresat}}."
	tresc := wejscieTrescZPostaci(&postac)
	plik, _, err := wejscieZlozOoxml(&postac, tresc, "Wzór pisma", false)
	if err != nil {
		t.Fatalf("złożenie pliku wzorcowego odmówiło: %v", err)
	}
	dokument := wydanieWniesDokument(t, adapter, zycie, plik, "okno-obieg-szablonu")

	szablon, err := adapter.ZapiszSzablonPisma(zycie, shared.StudioTemplateSaveRequest{
		DocumentId: &dokument,
		Name:       "Wzór do oddania",
	})
	if err != nil {
		t.Fatalf("zapis szablonu odmówił: %v", err)
	}

	oddany, err := adapter.OddajSzablonDoPliku(zycie, shared.StudioTemplateExportRequest{
		TemplateId: szablon.Template.Id,
	})
	if err != nil {
		t.Fatalf("oddanie szablonu do pliku odmówiło: %v", err)
	}
	if oddany.Result.AssetId == nil {
		t.Fatal("oddanie szablonu nie odłożyło pliku w magazynie zasobów")
	}
	nazwanePola := false
	for _, cecha := range oddany.Result.DroppedFeatures {
		if strings.Contains(cecha.Reason, "pola do wypełnienia") {
			nazwanePola = true
		}
	}
	if !nazwanePola {
		t.Fatalf("oddanie szablonu nie powiedziało, że wykaz pól nie wchodzi do pliku: %#v",
			oddany.Result.DroppedFeatures)
	}

	bajty, err := adapter.bajtyZasobuStudia(zycie, *oddany.Result.AssetId)
	if err != nil {
		t.Fatalf("odczyt pliku szablonu odmówił: %v", err)
	}
	wniesiony, err := adapter.WniesSzablonZPliku(zycie, shared.StudioTemplateImportRequest{
		Name:        wskaznik("Wzór wniesiony z pliku"),
		BytesBase64: wskaznik(base64.StdEncoding.EncodeToString(bajty)),
		Path:        wskaznik("wzor.dotx"),
	})
	if err != nil {
		t.Fatalf("wniesienie szablonu z pliku odmówiło: %v", err)
	}
	if wniesiony.Template.Form == nil || len(wniesiony.Template.Form.Styles) == 0 {
		t.Fatal("szablon wniesiony z pliku nie przejął arkusza stylów — a to jest " +
			"wprost wymaganie Właściciela")
	}
	if wniesiony.Template.Form.PageSetup == nil {
		t.Fatal("szablon wniesiony z pliku nie przejął nastaw strony")
	}
	polaOdtworzone := false
	for _, pole := range wniesiony.Template.Fields {
		if pole.Name == "adresat" {
			polaOdtworzone = true
		}
	}
	if !polaOdtworzone {
		t.Fatalf("pola do wypełnienia nie odtworzyły się ze znaczników treści: %#v",
			wniesiony.Template.Fields)
	}
	if wniesiony.Balance.Note == nil ||
		!strings.Contains(*wniesiony.Balance.Note, "pól do wypełnienia") {
		t.Fatalf("bilans wniesienia szablonu nie mówi, co przejęto: %v",
			wniesiony.Balance.Note)
	}
}

// TestSzablonZakladaDokumentZPostacia mierzy założenie dokumentu Z SZABLONU:
// nowa strona ma stanąć na papierze firmowym, nie na domyślnym.
func TestSzablonZakladaDokumentZPostacia(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	postac := wejscieDokumentWzorcowy("")
	tresc := wejscieTrescZPostaci(&postac)
	plik, _, err := wejscieZlozOoxml(&postac, tresc, "Wzór pisma", false)
	if err != nil {
		t.Fatalf("złożenie pliku wzorcowego odmówiło: %v", err)
	}
	dokument := wydanieWniesDokument(t, adapter, zycie, plik, "okno-z-szablonu")

	szablon, err := adapter.ZapiszSzablonPisma(zycie, shared.StudioTemplateSaveRequest{
		DocumentId: &dokument,
		Name:       "Papier firmowy Danaco",
	})
	if err != nil {
		t.Fatalf("zapis szablonu odmówił: %v", err)
	}

	zalozony, err := adapter.ZalozDokument(zycie, shared.StudioDocumentCreateRequest{
		WindowId:   "okno-z-szablonu",
		TemplateId: &szablon.Template.Id,
	})
	if err != nil {
		t.Fatalf("założenie dokumentu z szablonu odmówiło: %v", err)
	}
	odczytana, err := adapter.PostacDokumentu(zycie, shared.StudioDocumentFormGetRequest{
		DocumentId: zalozony.Document.Id,
	})
	if err != nil {
		t.Fatalf("odczyt postaci dokumentu z szablonu odmówił: %v", err)
	}
	stylWlasny := false
	for _, styl := range odczytana.Form.Styles {
		if styl.Name == "podstawa prawna" {
			stylWlasny = true
		}
	}
	if !stylWlasny {
		t.Fatalf("dokument z szablonu nie dostał arkusza stylów wzorca; stylów: %d",
			len(odczytana.Form.Styles))
	}
	if zalozony.Document.Title == nil || *zalozony.Document.Title != "Papier firmowy Danaco" {
		t.Fatalf("dokument z szablonu nie przejął nazwy wzorca: %v", zalozony.Document.Title)
	}
}
