// Plik mierzy skutek wydania do formatu: czy plik naprawdę powstaje, czy
// postać dochodzi tam, gdzie format ją niesie, i czy strata jest nazwana tam,
// gdzie format jej nie niesie.
package core

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"

	"danacoconsole/shared"
)


// TestWydanieTekstemNazywaZgubionaTabele mierzy zasadę rozstrzygającą wydania:
// format uboższy niż dokument jest normalny, przemilczenie straty — nie.
func TestWydanieTekstemNazywaZgubionaTabele(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	postac := wejscieDokumentWzorcowy("")
	tresc := wejscieTrescZPostaci(&postac)
	bajty, _, err := wejscieZlozOoxml(&postac, tresc, "Umowa najmu", false)
	if err != nil {
		t.Fatalf("złożenie pliku wejściowego odmówiło: %v", err)
	}
	dokument := wydanieWniesDokument(t, adapter, zycie, bajty, "okno-wydanie")

	wynik, err := adapter.WydajDoFormatu(zycie, shared.StudioDocumentExportFormatRequest{
		DocumentId: dokument,
		Format:     shared.StudioExportFormatTxt,
	})
	if err != nil {
		t.Fatalf("wydanie do txt odmówiło: %v", err)
	}
	if wynik.Result.AssetId == nil {
		t.Fatal("wydanie nie odłożyło bajtów w magazynie zasobów — cienka instalka " +
			"nie miałaby po co sięgnąć")
	}
	if wynik.Result.Bytes == nil || *wynik.Result.Bytes == 0 {
		t.Fatal("wydanie oddało zero bajtów — plik pusty nie jest wydaniem")
	}

	nazwana := false
	for _, cecha := range wynik.Result.DroppedFeatures {
		if strings.Contains(cecha.Reason, "tabel") {
			nazwana = true
			if cecha.Detail == nil || !strings.Contains(*cecha.Detail, "1 tabel") {
				t.Fatalf("strata tabeli nazwana bez liczby: %v", cecha.Detail)
			}
		}
	}
	if !nazwana {
		t.Fatalf("wydanie do txt nie nazwało zgubionej siatki tabeli; wykaz: %#v",
			wynik.Result.DroppedFeatures)
	}
	if wynik.Result.Note == nil || !strings.Contains(*wynik.Result.Note, "pominęło") {
		t.Fatalf("wynik nie niesie zdania o tym, co odpadło: %v", wynik.Result.Note)
	}

	// Treść komórek MA wyjść — strata dotyczy siatki, nie danych.
	odczytane, err := adapter.bajtyZasobuStudia(zycie, *wynik.Result.AssetId)
	if err != nil {
		t.Fatalf("odczyt wydanego pliku odmówił: %v", err)
	}
	plik := string(odczytane)
	if !strings.Contains(plik, "Wynagrodzenie") || !strings.Contains(plik, "4 200,00") {
		t.Fatalf("treść komórek tabeli nie weszła do pliku txt: %q", plik)
	}
}

// wydanieWniesDokument wnosi plik do edytora poleceniem otwarcia dokumentu
// i oddaje kod dokumentu przydzielony przez edytor, gotowy do dalszych wydań.
func wydanieWniesDokument(t *testing.T, adapter *adapterStudia, zycie context.Context,
	bajty []byte, okno string) string {
	t.Helper()
	zapis := base64.StdEncoding.EncodeToString(bajty)
	wniesiony, err := adapter.WniesPlikDoEdytora(zycie, shared.StudioDocumentImportFileRequest{
		WindowId:    okno,
		BytesBase64: &zapis,
		Path:        wskaznik("umowa.docx"),
	})
	if err != nil {
		t.Fatalf("wniesienie pliku odmówiło: %v", err)
	}
	return wniesiony.Document.Id
}

// TestWydanieDocxPrzenosiPostac mierzy, że wydanie do formatu bogatego NIE gubi
// postaci: plik wydany i wczytany z powrotem niesie styl, sekcje i tabelę.
func TestWydanieDocxPrzenosiPostac(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	postac := wejscieDokumentWzorcowy("")
	tresc := wejscieTrescZPostaci(&postac)
	wejsciowy, _, err := wejscieZlozOoxml(&postac, tresc, "Umowa najmu", false)
	if err != nil {
		t.Fatalf("złożenie pliku wejściowego odmówiło: %v", err)
	}
	dokument := wydanieWniesDokument(t, adapter, zycie, wejsciowy, "okno-wydanie-docx")

	wynik, err := adapter.WydajDoFormatu(zycie, shared.StudioDocumentExportFormatRequest{
		DocumentId: dokument,
		Format:     shared.StudioExportFormatDocx,
	})
	if err != nil {
		t.Fatalf("wydanie do docx odmówiło: %v", err)
	}
	bajty, err := adapter.bajtyZasobuStudia(zycie, *wynik.Result.AssetId)
	if err != nil {
		t.Fatalf("odczyt wydanego pliku odmówił: %v", err)
	}

	odczytana, _, bilans, err := wejscieCzytajOoxml("studio-dok-po-wydaniu", bajty,
		shared.StudioImportFormatDocx)
	if err != nil {
		t.Fatalf("odczyt wydanego docx odmówił: %v", err)
	}
	if len(odczytana.Tables) != 1 {
		t.Fatalf("tabela nie przeszła przez wydanie do docx; tabel: %d",
			len(odczytana.Tables))
	}
	if len(odczytana.Sections) == 0 {
		t.Fatal("sekcje nie przeszły przez wydanie do docx")
	}
	stylWlasny := false
	for _, styl := range odczytana.Styles {
		if styl.Name == "podstawa prawna" {
			stylWlasny = true
		}
	}
	if !stylWlasny {
		t.Fatalf("styl własny nie przeszedł przez wydanie do docx; stylów: %d",
			len(odczytana.Styles))
	}
	if bilans.StylesRecovered == nil || *bilans.StylesRecovered == 0 {
		t.Fatal("bilans odczytu wydanego pliku nie policzył stylów")
	}
}

// TestWydanieMarkdownemNazywaScalenia mierzy przykład, który zlecenie wymienia
// wprost: tabela o scalonych komórkach w markdown.
func TestWydanieMarkdownemNazywaScalenia(t *testing.T) {
	postac := wejscieDokumentWzorcowy("studio-dok-md")
	postac.Tables[0].Cells[0].ColumnSpan = wejscieWskaznikCalkowity(2)

	bajty, pominiete := wydanieMarkdownem(&postac, wejscieTrescZPostaci(&postac))
	plik := string(bajty)
	if !strings.Contains(plik, "# Umowa najmu") {
		t.Fatalf("nagłówek nie wyszedł znacznikiem markdown: %q", plik)
	}
	if !strings.Contains(plik, "| Pozycja |") {
		t.Fatalf("tabela nie wyszła tabelą markdown: %q", plik)
	}
	nazwane := false
	for _, cecha := range pominiete {
		if strings.Contains(cecha.Reason, "scalone") {
			nazwane = true
		}
	}
	if !nazwane {
		t.Fatalf("wydanie do markdown nie nazwało scaleń, których format nie niesie: %#v",
			pominiete)
	}
}

// TestWydanieHtmlemNiesieArkuszStylow mierzy, że style nazwane wychodzą arkuszem,
// a nie postacią wpisaną przy każdym akapicie.
func TestWydanieHtmlemNiesieArkuszStylow(t *testing.T) {
	postac := wejscieDokumentWzorcowy("studio-dok-html")
	bajty, _ := wydanieHtmlem(&postac, wejscieTrescZPostaci(&postac), "Umowa najmu")
	plik := string(bajty)

	if !strings.Contains(plik, ".styl-podstawa-prawna {") {
		t.Fatalf("styl własny nie wyszedł zasadą arkusza: %q", plik)
	}
	if !strings.Contains(plik, "class=\"styl-podstawa-prawna\"") {
		t.Fatalf("akapit nie wskazuje klasy swojego stylu: %q", plik)
	}
	if !strings.Contains(plik, "<table>") || !strings.Contains(plik, "colspan") &&
		strings.Contains(plik, "rowspan") {
		t.Fatalf("tabela nie wyszła tabelą HTML: %q", plik)
	}
	if !strings.Contains(plik, "@page") {
		t.Fatalf("nastawy nośnika nie wyszły regułą druku: %q", plik)
	}
}

// TestWydanieRtfemJestWykonalne mierzy rozstrzygnięcie zlecenia o RTF: format
// wchodzi do wykazu, jeśli rachunek własny go składa bez nadmiaru.
func TestWydanieRtfemJestWykonalne(t *testing.T) {
	postac := wejscieDokumentWzorcowy("studio-dok-rtf")
	bajty, pominiete := wydanieRtfem(&postac, wejscieTrescZPostaci(&postac))
	plik := string(bajty)

	if !strings.HasPrefix(plik, `{\rtf1`) {
		t.Fatalf("plik nie zaczyna się nagłówkiem RTF: %q", plik[:wydanieWieksza(len(plik), 40)])
	}
	if !strings.Contains(plik, `\cellx`) {
		t.Fatalf("tabela nie wyszła siatką RTF: %q", plik)
	}
	if !strings.Contains(plik, `\u380?`) && !strings.Contains(plik, `\u261?`) {
		t.Fatalf("polskie litery nie wyszły zapisem \\uNNNN: %q", plik)
	}
	if len(pominiete) == 0 {
		t.Fatal("wydanie do RTF nie nazwało ani jednej cechy pominiętej — arkusz stylów " +
			"i aparat nie przechodzą, więc cisza jest tu nieprawdą")
	}
}

// TestWydanieWsadoweNieWstrzymujeSie mierzy zasadę wsadu: odmowa jednego
// dokumentu nie wstrzymuje pozostałych, a każda odmowa jest nazwana.
func TestWydanieWsadoweNieWstrzymujeSie(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	pierwszy, err := adapter.ZalozDokument(zycie, shared.StudioDocumentCreateRequest{
		WindowId: "okno-wsad", Title: wskaznik("Pismo pierwsze"),
	})
	if err != nil {
		t.Fatalf("założenie dokumentu odmówiło: %v", err)
	}
	if _, err := adapter.ZmienTresc(zycie, shared.StudioTextEditRequest{
		DocumentId: pierwszy.Document.Id, RangeStart: 0, RangeEnd: 0,
		Text: "Treść pisma pierwszego.",
	}); err != nil {
		t.Fatalf("zapis treści odmówił: %v", err)
	}
	drugi, err := adapter.ZalozDokument(zycie, shared.StudioDocumentCreateRequest{
		WindowId: "okno-wsad", Title: wskaznik("Pismo drugie"),
	})
	if err != nil {
		t.Fatalf("założenie dokumentu odmówiło: %v", err)
	}
	if _, err := adapter.ZmienTresc(zycie, shared.StudioTextEditRequest{
		DocumentId: drugi.Document.Id, RangeStart: 0, RangeEnd: 0,
		Text: "Treść pisma drugiego.",
	}); err != nil {
		t.Fatalf("zapis treści odmówił: %v", err)
	}

	wynik, err := adapter.WydajWsadowo(zycie, shared.StudioDocumentExportBatchRequest{
		DocumentIds: []string{pierwszy.Document.Id, "studio-dok-nie-ma", drugi.Document.Id},
		Format:      shared.StudioExportFormatTxt,
	})
	if err != nil {
		t.Fatalf("wydanie wsadowe odmówiło w całości: %v", err)
	}
	if wynik.Succeeded != 2 {
		t.Fatalf("wsad wydał %d dokumentów, a dwa były do wydania", wynik.Succeeded)
	}
	if wynik.Failed != 1 || len(wynik.Failures) != 1 {
		t.Fatalf("wsad nie nazwał odmowy dokumentu nieistniejącego: %#v", wynik.Failures)
	}
	if wynik.Failures[0].Detail == nil ||
		!strings.Contains(*wynik.Failures[0].Detail, "studio-dok-nie-ma") {
		t.Fatalf("odmowa nie mówi, którego dokumentu dotyczy: %#v", wynik.Failures[0])
	}
	for _, wydany := range wynik.Results {
		if wydany.AssetId == nil {
			t.Fatalf("dokument %s wydany bez zasobu magazynu", wydany.DocumentId)
		}
	}
}

// TestWydanieOdmawiaFormatuNieznanego mierzy, że format spoza wykazu kontraktu
// kończy się odmową NAZWANĄ, a nie plikiem o obcym rozszerzeniu.
func TestWydanieOdmawiaFormatuNieznanego(t *testing.T) {
	adapter, zycie := wejscieUprzazSprawdzianu(t)

	zalozony, err := adapter.ZalozDokument(zycie, shared.StudioDocumentCreateRequest{
		WindowId: "okno-format", Title: wskaznik("Pismo"),
	})
	if err != nil {
		t.Fatalf("założenie dokumentu odmówiło: %v", err)
	}
	_, err = adapter.WydajDoFormatu(zycie, shared.StudioDocumentExportFormatRequest{
		DocumentId: zalozony.Document.Id,
		Format:     shared.StudioExportFormat("pages"),
	})
	if err == nil {
		t.Fatal("wydanie do formatu spoza kontraktu przeszło — a nie ma czym go złożyć")
	}
	if !strings.Contains(err.Error(), "txt, md, docx, odt") {
		t.Fatalf("odmowa nie wymienia formatów, które rdzeń wydaje: %v", err)
	}
}
