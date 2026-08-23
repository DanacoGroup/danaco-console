package core

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"danacoconsole/shared"
)

// Sprawdziany skutku wejścia i wyjścia modułu Studio: czy postać dokumentu
// przechodzi PRZEZ PLIK tam i z powrotem.
//
// ── Dlaczego materiał jest składany ręcznie, a nie własnym składaczem ────────
// Gdyby plik próbny powstawał `wejscieZlozOoxml`, sprawdzian mierzyłby zgodność
// składacza z własnym rozbiorem — czyli że rdzeń czyta to, co sam napisał.
// Taki pomiar przechodzi także wtedy, gdy oba końce mylą się w ten sam sposób,
// a plik jest dla Worda nieczytelny. Dlatego materiał wejściowy jest tu
// WPISANY WPROST: archiwum ZIP ze składnikami XML w postaci, w jakiej wychodzą
// z pakietu biurowego. Rozbiór ma zdać egzamin z cudzego pliku, nie ze swojego.
//
// ── Dlaczego porównanie idzie PO ODCZYCIE, nie po bajtach ───────────────────
// Ten sam dokument da się zapisać na wiele poprawnych sposobów: inna kolejność
// węzłów, inne nazwy stylów automatycznych, inne zaokrąglenie twipów. Bajt
// w bajt nie zgodzi się nigdy i nie ma się zgodzić. Miarą jest to, czy po
// wczytaniu wyniku POSTAĆ jest ta sama: nazwa stylu akapitu, orientacja
// i marginesy sekcji, wymiary tabeli, jej wiersz nagłówkowy, scalenie komórek
// i szerokości kolumn.
//
// Szkody, które ten plik ma wykluczyć:
//  1. wczytanie `.docx`, po którym w postaci stoi sam tekst, a styl, sekcja
//     i tabela przepadły;
//  2. wydanie `.docx`, które zapisuje treść i gubi to, co przyszło z pliku —
//     czyli obieg gubiący postać dokładnie tam, gdzie Właściciel go sprawdza;
//  3. wydanie do formatu uboższego, które o stracie milczy;
//  4. wydanie wielostronicowe oddające jedną stronę, choć odpowiedź mówi inaczej;
//  5. plik wyjściowy niosący warstwę znakowania sesji — komentarze i wyróżnienia
//     w piśmie wysłanym na zewnątrz.

// ── Uprząż wspólna ──────────────────────────────────────────────────────────

// ooxmlNaglowekXml jest nagłówkiem składnika XML archiwum biurowego.
const ooxmlNaglowekXml = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`

// ooxmlZlozArchiwumProbne składa archiwum ZIP ze wskazanych składników.
// Składnik wskazany jako pierwszy idzie bez kompresji — tego wymaga
// OpenDocument od swojego `mimetype`.
func ooxmlZlozArchiwumProbne(t *testing.T, skladniki map[string]string, pierwszy string) []byte {
	t.Helper()

	var bufor bytes.Buffer
	archiwum := zip.NewWriter(&bufor)
	zapisz := func(nazwa string, bezKompresji bool) {
		var wpis interface{ Write([]byte) (int, error) }
		var err error
		if bezKompresji {
			wpis, err = archiwum.CreateHeader(&zip.FileHeader{Name: nazwa, Method: zip.Store})
		} else {
			wpis, err = archiwum.Create(nazwa)
		}
		if err != nil {
			t.Fatalf("nie można założyć składnika %s archiwum próbnego: %v", nazwa, err)
		}
		if _, err := wpis.Write([]byte(skladniki[nazwa])); err != nil {
			t.Fatalf("nie można zapisać składnika %s archiwum próbnego: %v", nazwa, err)
		}
	}
	if pierwszy != "" {
		zapisz(pierwszy, true)
	}
	for nazwa := range skladniki {
		if nazwa == pierwszy {
			continue
		}
		zapisz(nazwa, false)
	}
	if err := archiwum.Close(); err != nil {
		t.Fatalf("nie można domknąć archiwum próbnego: %v", err)
	}
	return bufor.Bytes()
}

// ooxmlDocxProbny składa dokument `.docx` w postaci, w jakiej wychodzi z Worda.
//
// Materiał niesie wszystko, co sprawdzian mierzy, i nic ponad to:
//   - arkusz stylów z nagłówkiem poziomu pierwszego i tekstem zasadniczym,
//   - akapit nagłówkowy o stylu `Heading1` oraz akapit zasadniczy z fragmentem
//     wytłuszczonym (postać znaku ma przejść razem ze stylem akapitu),
//   - tabelę dwa na trzy z wierszem nagłówkowym powtarzanym, scaleniem dwóch
//     kolumn w wierszu drugim i jawną siatką szerokości,
//   - sekcję poziomą A4 o marginesach 30/20/25/15 mm.
func ooxmlDocxProbny(t *testing.T) []byte {
	t.Helper()

	dokument := ooxmlNaglowekXml +
		`<w:document xmlns:w="` + ooxmlPrzestrzenGlowna + `" xmlns:r="` +
		ooxmlPrzestrzenPowiazania + `"><w:body>` +

		// Nagłówek poziomu pierwszego.
		`<w:p><w:pPr><w:pStyle w:val="Heading1"/><w:outlineLvl w:val="0"/></w:pPr>` +
		`<w:r><w:t>Umowa o dzieło</w:t></w:r></w:p>` +

		// Akapit zasadniczy: fragment zwykły i fragment wytłuszczony.
		`<w:p><w:pPr><w:pStyle w:val="Normal"/><w:jc w:val="both"/>` +
		`<w:ind w:firstLine="708"/></w:pPr>` +
		`<w:r><w:t xml:space="preserve">Strony ustalaja </w:t></w:r>` +
		`<w:r><w:rPr><w:b/></w:rPr><w:t>zakres prac</w:t></w:r>` +
		`<w:r><w:t> na warunkach ponizszych.</w:t></w:r></w:p>` +

		// Tabela: wiersz nagłówkowy powtarzany, w drugim wierszu scalenie
		// dwóch pierwszych kolumn.
		`<w:tbl><w:tblPr><w:tblStyle w:val="Siatka"/>` +
		`<w:tblW w:w="9070" w:type="dxa"/>` +
		`<w:tblBorders><w:top w:val="single" w:sz="8" w:color="000000"/>` +
		`<w:left w:val="single" w:sz="8" w:color="000000"/>` +
		`<w:bottom w:val="single" w:sz="8" w:color="000000"/>` +
		`<w:right w:val="single" w:sz="8" w:color="000000"/></w:tblBorders></w:tblPr>` +
		`<w:tblGrid><w:gridCol w:w="3023"/><w:gridCol w:w="3023"/>` +
		`<w:gridCol w:w="3024"/></w:tblGrid>` +
		`<w:tr><w:trPr><w:tblHeader/></w:trPr>` +
		`<w:tc><w:tcPr/><w:p><w:r><w:t>Pozycja</w:t></w:r></w:p></w:tc>` +
		`<w:tc><w:tcPr/><w:p><w:r><w:t>Termin</w:t></w:r></w:p></w:tc>` +
		`<w:tc><w:tcPr/><w:p><w:r><w:t>Kwota</w:t></w:r></w:p></w:tc></w:tr>` +
		`<w:tr>` +
		`<w:tc><w:tcPr><w:gridSpan w:val="2"/></w:tcPr>` +
		`<w:p><w:r><w:t>Projekt i wykonanie</w:t></w:r></w:p></w:tc>` +
		`<w:tc><w:tcPr/><w:p><w:r><w:t>12 000</w:t></w:r></w:p></w:tc></w:tr>` +
		`</w:tbl>` +

		// Nastawy sekcji: A4 pozioma, marginesy niesymetryczne.
		`<w:sectPr><w:pgSz w:w="16838" w:h="11906" w:orient="landscape"/>` +
		`<w:pgMar w:top="1701" w:right="850" w:bottom="1417" w:left="1134"/>` +
		`</w:sectPr>` +
		`</w:body></w:document>`

	style := ooxmlNaglowekXml +
		`<w:styles xmlns:w="` + ooxmlPrzestrzenGlowna + `">` +
		`<w:style w:type="paragraph" w:styleId="Normal" w:default="1">` +
		`<w:name w:val="Normal"/>` +
		`<w:rPr><w:rFonts w:ascii="Times New Roman"/><w:sz w:val="24"/></w:rPr>` +
		`</w:style>` +
		`<w:style w:type="paragraph" w:styleId="Heading1"><w:name w:val="heading 1"/>` +
		`<w:basedOn w:val="Normal"/>` +
		`<w:pPr><w:outlineLvl w:val="0"/><w:keepNext/></w:pPr>` +
		`<w:rPr><w:b/><w:sz w:val="36"/></w:rPr></w:style>` +
		`</w:styles>`

	return ooxmlZlozArchiwumProbne(t, map[string]string{
		"[Content_Types].xml": ooxmlNaglowekXml +
			`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
			`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
			`<Default Extension="xml" ContentType="application/xml"/>` +
			`<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>` +
			`</Types>`,
		"_rels/.rels": ooxmlNaglowekXml +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>` +
			`</Relationships>`,
		ooxmlSkladnikDokumentu: dokument,
		ooxmlSkladnikStylow:    style,
	}, "")
}

// ooxmlOdtProbny składa dokument `.odt` w postaci, w jakiej wychodzi
// z LibreOffice: `mimetype` pierwszy i nieskompresowany, treść w `content.xml`,
// arkusz stylów i układ strony w `styles.xml`.
//
// Materiał niesie te same rzeczy, co próbny `.docx`, żeby oba sprawdziany
// mierzyły to samo: nagłówek o stylu nazwanym, akapit z fragmentem
// wytłuszczonym, tabelę trzykolumnową ze scaleniem i wierszem nagłówkowym,
// oraz sekcję poziomą o marginesach niesymetrycznych.
func ooxmlOdtProbny(t *testing.T) []byte {
	t.Helper()

	przestrzenie := ` xmlns:office="urn:oasis:names:tc:opendocument:xmlns:office:1.0"` +
		` xmlns:text="urn:oasis:names:tc:opendocument:xmlns:text:1.0"` +
		` xmlns:table="urn:oasis:names:tc:opendocument:xmlns:table:1.0"` +
		` xmlns:style="urn:oasis:names:tc:opendocument:xmlns:style:1.0"` +
		` xmlns:fo="urn:oasis:names:tc:opendocument:xmlns:xsl-fo-compatible:1.0"`

	tresc := ooxmlNaglowekXml + `<office:document-content` + przestrzenie + `>` +
		`<office:automatic-styles>` +
		`<style:style style:name="T-wytluszczony" style:family="text">` +
		`<style:text-properties fo:font-weight="bold"/></style:style>` +
		`<style:style style:name="Kolumna-jedna" style:family="table-column">` +
		`<style:table-column-properties style:column-width="53.3mm"/></style:style>` +
		`</office:automatic-styles>` +
		`<office:body><office:text>` +

		`<text:h text:style-name="Naglowek-pierwszy" text:outline-level="1">` +
		`Umowa o dzielo</text:h>` +

		`<text:p text:style-name="Tekst-zasadniczy">Strony ustalaja ` +
		`<text:span text:style-name="T-wytluszczony">zakres prac</text:span>` +
		` na warunkach ponizszych.</text:p>` +

		`<table:table table:name="Tabela-pozycji">` +
		`<table:table-column table:style-name="Kolumna-jedna" table:number-columns-repeated="3"/>` +
		`<table:table-header-rows><table:table-row>` +
		`<table:table-cell office:value-type="string"><text:p>Pozycja</text:p></table:table-cell>` +
		`<table:table-cell office:value-type="string"><text:p>Termin</text:p></table:table-cell>` +
		`<table:table-cell office:value-type="string"><text:p>Kwota</text:p></table:table-cell>` +
		`</table:table-row></table:table-header-rows>` +
		`<table:table-row>` +
		`<table:table-cell table:number-columns-spanned="2" office:value-type="string">` +
		`<text:p>Projekt i wykonanie</text:p></table:table-cell>` +
		`<table:covered-table-cell/>` +
		`<table:table-cell office:value-type="string"><text:p>12 000</text:p></table:table-cell>` +
		`</table:table-row>` +
		`</table:table>` +

		`</office:text></office:body></office:document-content>`

	style := ooxmlNaglowekXml + `<office:document-styles` + przestrzenie + `>` +
		`<office:styles>` +
		`<style:style style:name="Tekst-zasadniczy" style:family="paragraph">` +
		`<style:text-properties fo:font-family="Times New Roman" fo:font-size="12pt"/>` +
		`<style:paragraph-properties fo:text-align="justify"/></style:style>` +
		`<style:style style:name="Naglowek-pierwszy" style:family="paragraph"` +
		` style:parent-style-name="Tekst-zasadniczy">` +
		`<style:text-properties fo:font-weight="bold" fo:font-size="18pt"/>` +
		`<style:paragraph-properties fo:keep-with-next="always"/></style:style>` +
		`</office:styles>` +
		`<office:automatic-styles>` +
		`<style:page-layout style:name="Uklad-pierwszy">` +
		`<style:page-layout-properties fo:page-width="297mm" fo:page-height="210mm"` +
		` style:print-orientation="landscape" fo:margin-top="30mm" fo:margin-bottom="25mm"` +
		` fo:margin-left="20mm" fo:margin-right="15mm"/>` +
		`</style:page-layout></office:automatic-styles>` +
		`<office:master-styles>` +
		`<style:master-page style:name="Standard" style:page-layout-name="Uklad-pierwszy"/>` +
		`</office:master-styles>` +
		`</office:document-styles>`

	return ooxmlZlozArchiwumProbne(t, map[string]string{
		odfSkladnikRodzaju: "application/vnd.oasis.opendocument.text",
		odfSkladnikManifestu: ooxmlNaglowekXml +
			`<manifest:manifest xmlns:manifest="urn:oasis:names:tc:opendocument:xmlns:manifest:1.0">` +
			`<manifest:file-entry manifest:full-path="/" manifest:media-type="application/vnd.oasis.opendocument.text"/>` +
			`<manifest:file-entry manifest:full-path="content.xml" manifest:media-type="text/xml"/>` +
			`<manifest:file-entry manifest:full-path="styles.xml" manifest:media-type="text/xml"/>` +
			`</manifest:manifest>`,
		odfSkladnikTresci: tresc,
		odfSkladnikStylow: style,
	}, odfSkladnikRodzaju)
}

// ooxmlWniesPlik wnosi bajty pliku do edytora i oddaje odpowiedź wraz
// z bilansem.
func ooxmlWniesPlik(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno string, bajty []byte,
	format shared.StudioImportFormat) shared.StudioDocumentImportFileResponse {

	t.Helper()

	var wniesiony shared.StudioDocumentImportFileResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentImportFile,
		shared.StudioDocumentImportFileRequest{
			WindowId:    okno,
			Format:      &format,
			BytesBase64: wskaznik(wBase64(bajty)),
		}, &wniesiony)
	return wniesiony
}

// ooxmlPostacDokumentu odczytuje postać dokumentu OSOBNYM wywołaniem.
//
// To jest miara właściwa: postać wzięta z odpowiedzi komendy zapisującej mówi,
// co rdzeń zamierzał zapisać, a nie co naprawdę leży w bazie. Operator otworzy
// dokument ponownie i zobaczy to drugie.
func ooxmlPostacDokumentu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	kod string) shared.StudioDocumentForm {

	t.Helper()

	var postac shared.StudioDocumentFormGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentFormGet,
		shared.StudioDocumentFormGetRequest{DocumentId: kod}, &postac)
	return postac.Form
}

// ooxmlWydajDoPliku wydaje dokument do wskazanego formatu i pliku, oddając
// wynik wraz z bajtami leżącymi pod ścieżką. Bajty czyta się z DYSKU, a nie
// z odpowiedzi — odpowiedź niesie rozmiar, a sprawdzian mierzy plik.
func ooxmlWydajDoPliku(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	kod string, format shared.StudioExportFormat, sciezka string,
	nastaw func(*shared.StudioDocumentExportFormatRequest)) (shared.StudioExportResult, []byte) {

	t.Helper()

	zadanie := shared.StudioDocumentExportFormatRequest{
		DocumentId: kod, Format: format, Path: wskaznik(sciezka),
	}
	if nastaw != nil {
		nastaw(&zadanie)
	}
	var wydane shared.StudioDocumentExportFormatResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentExportFormat,
		zadanie, &wydane)

	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("wydanie do %s zameldowało powodzenie, a pliku %s nie ma: %v",
			format, sciezka, err)
	}
	if len(bajty) == 0 {
		t.Fatalf("wydanie do %s zostawiło plik zerowej długości", format)
	}
	return wydane.Result, bajty
}

// ooxmlBlokOStylu szuka bloku akapitu o wskazanym stylu nazwanym.
func ooxmlBlokOStylu(postac shared.StudioDocumentForm, styl string) *shared.StudioDocumentBlock {
	for i := range postac.Blocks {
		blok := &postac.Blocks[i]
		if blok.Paragraph == nil || blok.Paragraph.StyleName == nil {
			continue
		}
		if *blok.Paragraph.StyleName == styl {
			return blok
		}
	}
	return nil
}

// ooxmlCzyFragmentWytluszczony rozstrzyga, czy w postaci stoi fragment
// o wskazanej treści i postaci wytłuszczonej.
func ooxmlCzyFragmentWytluszczony(postac shared.StudioDocumentForm, fragment string) bool {
	for i := range postac.Blocks {
		for j := range postac.Blocks[i].Runs {
			bieg := &postac.Blocks[i].Runs[j]
			if !strings.Contains(bieg.Text, fragment) {
				continue
			}
			if bieg.Format != nil && bieg.Format.Bold != nil && *bieg.Format.Bold {
				return true
			}
		}
	}
	return false
}

// ooxmlSprawdzPostacProbna mierzy postać wobec materiału próbnego. Ta sama
// miara przechodzi po wczytaniu pliku wejściowego i po wczytaniu pliku
// wydanego — i to jest cały sens tego sprawdzianu: obieg ma nie gubić niczego,
// co wymieniono niżej.
func ooxmlSprawdzPostacProbna(t *testing.T, postac shared.StudioDocumentForm,
	etap, stylNaglowka, stylZasadniczy string) {

	t.Helper()

	// ── Styl nazwany ────────────────────────────────────────────────────────
	if ooxmlBlokOStylu(postac, stylNaglowka) == nil {
		nazwy := make([]string, 0, len(postac.Blocks))
		for i := range postac.Blocks {
			if postac.Blocks[i].Paragraph != nil && postac.Blocks[i].Paragraph.StyleName != nil {
				nazwy = append(nazwy, *postac.Blocks[i].Paragraph.StyleName)
			}
		}
		t.Errorf("%s: nie ma akapitu o stylu %q — style w postaci: %v",
			etap, stylNaglowka, nazwy)
	}
	if ooxmlBlokOStylu(postac, stylZasadniczy) == nil {
		t.Errorf("%s: nie ma akapitu o stylu %q", etap, stylZasadniczy)
	}
	stoiWArkuszu := false
	for _, styl := range postac.Styles {
		if styl.Name == stylNaglowka {
			stoiWArkuszu = true
		}
	}
	if !stoiWArkuszu {
		t.Errorf("%s: styl %q jest używany przez akapit, a nie stoi w arkuszu stylów "+
			"— dokument odwołuje się do stylu, którego nie ma", etap, stylNaglowka)
	}

	// ── Postać znaku wewnątrz akapitu ───────────────────────────────────────
	if !ooxmlCzyFragmentWytluszczony(postac, "zakres prac") {
		t.Errorf("%s: fragment „zakres prac” stracił wytłuszczenie — postać znaku "+
			"nie przeszła razem ze stylem akapitu", etap)
	}

	// ── Sekcja: orientacja i marginesy ──────────────────────────────────────
	nastawy := postac.PageSetup
	if nastawy == nil && len(postac.Sections) > 0 {
		nastawy = postac.Sections[len(postac.Sections)-1].PageSetup
	}
	if nastawy == nil {
		t.Fatalf("%s: postać nie niesie nastaw strony ani na dokumencie, ani na sekcji", etap)
	}
	if nastawy.Orientation == nil || *nastawy.Orientation != shared.StudioPageOrientationPozioma {
		t.Errorf("%s: orientacja pozioma przepadła — w postaci stoi %v",
			etap, nastawy.Orientation)
	}
	// Marginesy niesymetryczne są tu miarą celową: gdyby rachunek podstawiał
	// nastawę domyślną, wszystkie cztery wyszłyby równe.
	if nastawy.MarginTop == nil || *nastawy.MarginTop != 30 {
		t.Errorf("%s: margines górny nie przeszedł — oczekiwano 30 mm, stoi %v",
			etap, nastawy.MarginTop)
	}
	if nastawy.MarginRight == nil || *nastawy.MarginRight != 15 {
		t.Errorf("%s: margines prawy nie przeszedł — oczekiwano 15 mm, stoi %v",
			etap, nastawy.MarginRight)
	}
	if len(postac.Sections) == 0 {
		t.Errorf("%s: postać nie niesie ani jednej sekcji", etap)
	}

	// ── Tabela ──────────────────────────────────────────────────────────────
	if len(postac.Tables) != 1 {
		t.Fatalf("%s: oczekiwano jednej tabeli, jest %d", etap, len(postac.Tables))
	}
	tabela := postac.Tables[0]
	if tabela.Rows != 2 || tabela.Columns != 3 {
		t.Errorf("%s: tabela ma %d na %d, a materiał niósł 2 na 3",
			etap, tabela.Rows, tabela.Columns)
	}
	if tabela.HeaderRows == nil || *tabela.HeaderRows != 1 {
		t.Errorf("%s: wiersz nagłówkowy tabeli przepadł — headerRows=%v",
			etap, tabela.HeaderRows)
	}
	if tabela.RepeatHeader == nil || !*tabela.RepeatHeader {
		t.Errorf("%s: powtarzanie wiersza nagłówkowego przepadło", etap)
	}
	// Szerokości POLICZONE, nie zerowe: kolumna zerowej szerokości jest
	// kolumną niewidoczną w pakiecie biurowym.
	if len(tabela.ColumnWidthsMm) != 3 {
		t.Errorf("%s: tabela niesie %d szerokości kolumn, a ma trzy kolumny",
			etap, len(tabela.ColumnWidthsMm))
	}
	for numer, szerokosc := range tabela.ColumnWidthsMm {
		if szerokosc <= 0 {
			t.Errorf("%s: kolumna %d ma szerokość %v — kolumna zerowej szerokości "+
				"jest w dokumencie niewidoczna", etap, numer, szerokosc)
		}
	}
	// Scalenie dwóch kolumn w wierszu drugim.
	scalajaca := ooxmlSzukajKomorki(&tabela, 1, 0)
	if scalajaca == nil || scalajaca.ColumnSpan == nil || *scalajaca.ColumnSpan != 2 {
		t.Errorf("%s: scalenie dwóch kolumn w wierszu drugim przepadło — komórka=%+v",
			etap, scalajaca)
	}
	wchlonieta := ooxmlSzukajKomorki(&tabela, 1, 1)
	if wchlonieta == nil || wchlonieta.Merged == nil || !*wchlonieta.Merged {
		t.Errorf("%s: komórka wchłonięta scaleniem nie jest oznaczona jako scalona "+
			"— po zapisie tabela rozjedzie się o kolumnę", etap)
	}
	// Treść komórek ma zostać treścią, nie zniknąć razem ze strukturą.
	if naglowek := ooxmlSzukajKomorki(&tabela, 0, 0); naglowek == nil ||
		naglowek.Text == nil || !strings.Contains(*naglowek.Text, "Pozycja") {
		t.Errorf("%s: treść pierwszej komórki nagłówka przepadła", etap)
	}
}

// ── Punkt pierwszy: obieg `.docx` ───────────────────────────────────────────

// TestDocxPoOdczycieZachowujeStylSekcjeITabele wykazuje, że dokument wczytany
// z `.docx` i oddany z powrotem do `.docx` zachowuje styl nazwany, sekcję
// i tabelę — mierzone PO ODCZYCIE obu plików, nie po ich bajtach.
func TestDocxPoOdczycieZachowujeStylSekcjeITabele(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	okno := "okno-obiegu-docx"

	// Etap pierwszy: cudzy plik wchodzi do edytora.
	wniesiony := ooxmlWniesPlik(t, zmontowany, zycie, okno,
		ooxmlDocxProbny(t), shared.StudioImportFormatDocx)

	poWczytaniu := ooxmlPostacDokumentu(t, zmontowany, zycie, wniesiony.Document.Id)
	ooxmlSprawdzPostacProbna(t, poWczytaniu, "po wczytaniu .docx",
		wejscieStylNaglowek1, wejscieStylTekstZasadniczy)

	// Bilans ma być prawdziwy, nie ozdobny.
	if wartoscCalkowita(wniesiony.Balance.TablesRecognized) != 1 {
		t.Errorf("bilans mówi o %d rozpoznanych tabelach, a tabela jest jedna",
			wartoscCalkowita(wniesiony.Balance.TablesRecognized))
	}
	if wartoscCalkowita(wniesiony.Balance.StylesRecovered) < 2 {
		t.Errorf("bilans mówi o %d przejętych stylach, a arkusz pliku niósł dwa",
			wartoscCalkowita(wniesiony.Balance.StylesRecovered))
	}

	// Etap drugi: dokument wychodzi z powrotem do `.docx`.
	sciezka := filepath.Join(t.TempDir(), "oddany.docx")
	_, bajty := ooxmlWydajDoPliku(t, zmontowany, zycie, wniesiony.Document.Id,
		shared.StudioExportFormatDocx, sciezka, nil)

	// Wydany plik ma być archiwum biurowym, a nie czymkolwiek o tej nazwie.
	skladniki, err := wejscieOtworzArchiwum(bajty)
	if err != nil {
		t.Fatalf("wydany .docx nie jest archiwum biurowym: %v", err)
	}
	for _, wymagany := range []string{
		"[Content_Types].xml", "_rels/.rels", ooxmlSkladnikDokumentu, ooxmlSkladnikStylow,
	} {
		if len(skladniki[wymagany]) == 0 {
			t.Errorf("wydany .docx nie niesie składnika %s — Word odmówi otwarcia", wymagany)
		}
	}

	// Etap trzeci: wydany plik wraca do edytora jako dokument osobny. Porównanie
	// idzie po TEJ SAMEJ mierze, co po wczytaniu pliku wejściowego.
	wrocony := ooxmlWniesPlik(t, zmontowany, zycie, okno+"-powrot",
		bajty, shared.StudioImportFormatDocx)
	if wrocony.Document.Id == wniesiony.Document.Id {
		t.Fatal("plik wydany wrócił do tego samego dokumentu — sprawdzian mierzyłby " +
			"postać, której nie przepuszczono przez plik")
	}
	poObiegu := ooxmlPostacDokumentu(t, zmontowany, zycie, wrocony.Document.Id)
	ooxmlSprawdzPostacProbna(t, poObiegu, "po obiegu przez wydany .docx",
		wejscieStylNaglowek1, wejscieStylTekstZasadniczy)

	// Treść ma przejść w całości — postać bez treści byłaby postacią pustego
	// dokumentu i przeszłaby każdy sprawdzian struktury.
	//
	// Mierzona jest treść AKAPITÓW. Treść komórek tabeli w treści płaskiej nie
	// stoi: płaski odczyt dokumentu składa się z akapitów, a tabela jest
	// strukturą i jej brzmienie stoi w postaci (sprawdzone wyżej, w komórkach).
	// Znaki diakrytyczne są tu miarą osobną i celową: „dzieło" przepuszczone
	// przez dwa archiwa i dwa rozbiory XML musi wrócić z literą „ł", a nie jako
	// „dzielo" ani jako znak zastępczy.
	trescPoObiegu := trescDokumentu(t, zmontowany, zycie, okno+"-powrot", wrocony.Document.Id)
	for _, oczekiwany := range []string{"Umowa o dzieło", "zakres prac", "ponizszych"} {
		if !strings.Contains(trescPoObiegu, oczekiwany) {
			t.Errorf("po obiegu przez .docx w treści nie ma %q; treść: %q",
				oczekiwany, trescPoObiegu)
		}
	}
}

// ── Punkt drugi: obieg `.odt` ───────────────────────────────────────────────

// TestOdtPoOdczycieZachowujeStylSekcjeITabele mierzy to samo, co sprawdzian
// `.docx`, na dokumencie OpenDocument.
func TestOdtPoOdczycieZachowujeStylSekcjeITabele(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	okno := "okno-obiegu-odt"

	wniesiony := ooxmlWniesPlik(t, zmontowany, zycie, okno,
		ooxmlOdtProbny(t), shared.StudioImportFormatOdt)

	// Styl ODF nosi nazwę własną pliku, a nie nazwę arkusza platformy — plik
	// LibreOffice nie ma powodu nazywać stylu tak, jak nazywa go ten produkt.
	poWczytaniu := ooxmlPostacDokumentu(t, zmontowany, zycie, wniesiony.Document.Id)
	ooxmlSprawdzPostacProbna(t, poWczytaniu, "po wczytaniu .odt",
		"Naglowek-pierwszy", "Tekst-zasadniczy")

	sciezka := filepath.Join(t.TempDir(), "oddany.odt")
	_, bajty := ooxmlWydajDoPliku(t, zmontowany, zycie, wniesiony.Document.Id,
		shared.StudioExportFormatOdt, sciezka, nil)

	// `mimetype` musi być składnikiem PIERWSZYM i NIESKOMPRESOWANYM — tym
	// OpenDocument rozpoznaje swoje pliki. Sprawdzenie idzie po nagłówkach
	// archiwum, nie po treści.
	archiwum, err := zip.NewReader(bytes.NewReader(bajty), int64(len(bajty)))
	if err != nil {
		t.Fatalf("wydany .odt nie jest archiwum: %v", err)
	}
	if len(archiwum.File) == 0 {
		t.Fatal("wydany .odt jest archiwum pustym")
	}
	if archiwum.File[0].Name != odfSkladnikRodzaju {
		t.Errorf("pierwszym składnikiem wydanego .odt jest %q, a ma być %q — "+
			"LibreOffice nie rozpozna rodzaju dokumentu",
			archiwum.File[0].Name, odfSkladnikRodzaju)
	}
	if archiwum.File[0].Method != zip.Store {
		t.Error("składnik mimetype wydanego .odt jest skompresowany — OpenDocument " +
			"wymaga go bez kompresji")
	}

	wrocony := ooxmlWniesPlik(t, zmontowany, zycie, okno+"-powrot",
		bajty, shared.StudioImportFormatOdt)
	poObiegu := ooxmlPostacDokumentu(t, zmontowany, zycie, wrocony.Document.Id)
	ooxmlSprawdzPostacProbna(t, poObiegu, "po obiegu przez wydany .odt",
		"Naglowek-pierwszy", "Tekst-zasadniczy")

	trescPoObiegu := trescDokumentu(t, zmontowany, zycie, okno+"-powrot", wrocony.Document.Id)
	for _, oczekiwany := range []string{"Umowa o dzielo", "zakres prac", "ponizszych"} {
		if !strings.Contains(trescPoObiegu, oczekiwany) {
			t.Errorf("po obiegu przez .odt w treści nie ma %q; treść: %q",
				oczekiwany, trescPoObiegu)
		}
	}
}

// ── Punkt trzeci: bilans cech pominiętych ───────────────────────────────────

// TestWydanieUbozszeOddajeWykazCechPominietych wykazuje, że wydanie do formatu,
// który czegoś nie niesie, MÓWI to wprost.
//
// Format uboższy niż dokument jest normalną sytuacją; przemilczenie straty nie
// jest. Miara: dokument z tabelą o scalonych komórkach i z przypisem wydany do
// `txt` i do `md` musi oddać wykaz niepusty, a w wykazie ma stać nazwa rzeczy,
// która odpadła — nie samo „coś odpadło".
func TestWydanieUbozszeOddajeWykazCechPominietych(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	okno := "okno-bilansu-wydania"

	// Dokument powstaje z `.docx`, żeby miał tabelę ze scaleniem naprawdę,
	// a nie tabelę zadeklarowaną w żądaniu.
	wniesiony := ooxmlWniesPlik(t, zmontowany, zycie, okno,
		ooxmlDocxProbny(t), shared.StudioImportFormatDocx)

	// Przypis dolny wchodzi drogą aparatu dokumentu — to jest cecha, której
	// tekst czysty nie niesie z natury.
	var przypis shared.StudioApparatusInsertResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioApparatusInsert,
		shared.StudioApparatusInsertRequest{
			DocumentId: wniesiony.Document.Id,
			Kind:       shared.StudioApparatusKindFootnote,
			Offset:     wskaznik(5),
			Text:       wskaznik("Podstawa: umowa z dnia pierwszego."),
		}, &przypis)

	katalog := t.TempDir()

	for _, przypadek := range []struct {
		format shared.StudioExportFormat
		nazwa  string
	}{
		{shared.StudioExportFormatTxt, "wydanie.txt"},
		{shared.StudioExportFormatMd, "wydanie.md"},
	} {
		wynik, bajty := ooxmlWydajDoPliku(t, zmontowany, zycie, wniesiony.Document.Id,
			przypadek.format, filepath.Join(katalog, przypadek.nazwa), nil)

		if len(wynik.DroppedFeatures) == 0 {
			t.Errorf("wydanie do %s nie oddało ani jednej cechy pominiętej, choć "+
				"dokument niesie tabelę ze scaleniem i przypis — przemilczenie "+
				"straty jest zakazane", przypadek.format)
			continue
		}
		if wynik.Note == nil || strings.TrimSpace(*wynik.Note) == "" {
			t.Errorf("wydanie do %s oddało wykaz cech pominiętych bez zdania "+
				"o tym, co odpadło", przypadek.format)
		}
		// Wykaz ma NAZYWAĆ rzeczy, nie liczyć je bezimiennie.
		zebrane := ""
		for _, pozycja := range wynik.DroppedFeatures {
			zebrane += " " + strings.ToLower(pozycja.Reason)
			if pozycja.Detail != nil {
				zebrane += " " + strings.ToLower(*pozycja.Detail)
			}
			if strings.TrimSpace(pozycja.Reason) == "" {
				t.Errorf("wydanie do %s oddało pozycję pominięcia bez powodu",
					przypadek.format)
			}
		}
		if !strings.Contains(zebrane, "przypis") {
			t.Errorf("wydanie do %s nie nazwało straty przypisu; wykaz: %s",
				przypadek.format, zebrane)
		}
		if przypadek.format == shared.StudioExportFormatTxt &&
			!strings.Contains(zebrane, "tabel") {
			t.Errorf("wydanie do txt nie nazwało straty na tabeli; wykaz: %s", zebrane)
		}
		if przypadek.format == shared.StudioExportFormatMd &&
			!strings.Contains(zebrane, "scal") {
			t.Errorf("wydanie do md nie nazwało straty na scalonych komórkach — "+
				"markdown scalenia nie wyraża; wykaz: %s", zebrane)
		}
		// Treść ma być w pliku niezależnie od strat na postaci.
		if !strings.Contains(string(bajty), "Umowa o dzieło") {
			t.Errorf("wydanie do %s zgubiło treść dokumentu", przypadek.format)
		}
	}
}

// ── Punkt czwarty: liczba stron wydania ─────────────────────────────────────

// TestWydaniePdfDajePlikOWlasciwejLiczbieStron wykazuje, że wydanie
// wielostronicowe daje plik o właściwej liczbie stron.
//
// Strony LICZONE są w pliku biblioteką `pdfcpu`, a nie brane z odpowiedzi.
// Miara jest odporna na zmianę wysokości kartki: ta sama treść wydana na A4
// i na A5 nie może dać tej samej liczby stron, a wydanie treści na wiele stron
// nie może dać jednej.
func TestWydaniePdfDajePlikOWlasciwejLiczbieStron(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	okno := "okno-stron-wydania"

	// Treść dostatecznie długa, żeby nie zmieściła się na jednej kartce.
	akapity := make([]string, 0, 160)
	for i := 0; i < 160; i++ {
		akapity = append(akapity,
			"Akapit o treści dostatecznie długiej, żeby zajął pełny wiersz na kartce pisma.")
	}
	dokument := dokumentZTrescia(t, zmontowany, zycie, okno, strings.Join(akapity, "\n"))

	katalog := t.TempDir()
	sciezka := filepath.Join(katalog, "wielostronicowy.pdf")
	_, bajty := ooxmlWydajDoPliku(t, zmontowany, zycie, dokument.Id,
		shared.StudioExportFormatPdf, sciezka, nil)

	stron, err := api.PageCount(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		t.Fatalf("wydany PDF nie daje się policzyć — nie jest poprawnym PDF-em: %v", err)
	}
	if stron < 2 {
		t.Fatalf("wydanie stu sześćdziesięciu akapitów dało %d stron — treść "+
			"nie zmieściłaby się na jednej kartce, więc plik jest obcięty", stron)
	}

	// Druga miara: wydanie połowy treści ma dać mniej stron. Bez tego
	// sprawdzian przechodziłby także wtedy, gdyby rachunek zawsze dawał
	// stałą liczbę stron większą od jednego.
	krotszy := dokumentZTrescia(t, zmontowany, zycie, okno,
		strings.Join(akapity[:40], "\n"))
	_, bajtyKrotsze := ooxmlWydajDoPliku(t, zmontowany, zycie, krotszy.Id,
		shared.StudioExportFormatPdf, filepath.Join(katalog, "krotszy.pdf"), nil)

	stronKrotszych, err := api.PageCount(bytes.NewReader(bajtyKrotsze), nastawyPdf())
	if err != nil {
		t.Fatalf("krótszy PDF nie daje się policzyć: %v", err)
	}
	if stronKrotszych >= stron {
		t.Errorf("dokument czterokrotnie krótszy dał %d stron, a dłuższy %d — "+
			"rachunek stron nie zależy od treści", stronKrotszych, stron)
	}
}

// ── Punkt piąty: plik wyjściowy bez warstwy znakowania ──────────────────────

// TestPlikWyjsciowyJestCzystyZWarstwyZnakowania wykazuje rozstrzygnięcie
// Właściciela: znakowanie żyje w sesji, nie w pliku wysyłanym na zewnątrz.
//
// Miara jest dwustronna, bo jedna strona nie wystarcza:
//   - w BAJTACH pliku nie ma brzmienia komentarza (gdyby był, adresat pisma
//     przeczytałby uwagi redakcyjne);
//   - PO ODCZYCIE pliku nie ma wyróżnienia tła (gdyby było, pismo wyszłoby
//     w kolorowych plamach roboczych).
//
// Sprawdzane jest też wskazanie `includeComments`: rdzeń komentarzy do pliku
// nie wpisuje, więc żądanie ich wydania ma być NAZWANE w wykazie cech
// pominiętych, a nie przemilczane albo spełnione.
func TestPlikWyjsciowyJestCzystyZWarstwyZnakowania(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	okno := "okno-czystego-wydania"

	wniesiony := ooxmlWniesPlik(t, zmontowany, zycie, okno,
		ooxmlDocxProbny(t), shared.StudioImportFormatDocx)

	const brzmienieKomentarza = "TU-WYMAGA-ZRODLA-uwaga-redakcyjna"
	var komentarz shared.StudioCommentAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioCommentAdd,
		shared.StudioCommentAddRequest{
			DocumentId:     wniesiony.Document.Id,
			Body:           brzmienieKomentarza,
			SelectionStart: wskaznik(0),
			SelectionEnd:   wskaznik(5),
		}, &komentarz)

	// Wyróżnienie tła jest cechą postaci dokumentu, więc idzie tą samą drogą,
	// co reszta formatowania.
	var wyroznienie shared.StudioFormatCharacterSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioFormatCharacterSet,
		shared.StudioFormatCharacterSetRequest{
			DocumentId:     wniesiony.Document.Id,
			RangeStart:     wskaznik(0),
			RangeEnd:       wskaznik(5),
			HighlightColor: wskaznik("#FFFF00"),
		}, &wyroznienie)

	// Wyróżnienie musi naprawdę stanąć w dokumencie — inaczej sprawdzian
	// mierzyłby brak, którego nikt nie założył.
	wDokumencie := ooxmlPostacDokumentu(t, zmontowany, zycie, wniesiony.Document.Id)
	if !ooxmlCzyStoiWyroznienie(wDokumencie) {
		t.Fatal("wyróżnienie nie stanęło w dokumencie — sprawdzian czystości pliku " +
			"mierzyłby brak, którego nie założono")
	}

	sciezka := filepath.Join(t.TempDir(), "czysty.docx")
	wynik, bajty := ooxmlWydajDoPliku(t, zmontowany, zycie, wniesiony.Document.Id,
		shared.StudioExportFormatDocx, sciezka,
		func(z *shared.StudioDocumentExportFormatRequest) {
			z.IncludeComments = wskaznik(true)
		})

	// ── Miara pierwsza: bajty pliku ─────────────────────────────────────────
	if bytes.Contains(bajty, []byte(brzmienieKomentarza)) {
		t.Error("brzmienie komentarza redakcyjnego stoi w bajtach wydanego .docx — " +
			"adresat pisma przeczytałby uwagi robocze")
	}

	// ── Miara druga: postać po odczycie wydanego pliku ──────────────────────
	wrocony := ooxmlWniesPlik(t, zmontowany, zycie, okno+"-powrot",
		bajty, shared.StudioImportFormatDocx)
	poObiegu := ooxmlPostacDokumentu(t, zmontowany, zycie, wrocony.Document.Id)
	if ooxmlCzyStoiWyroznienie(poObiegu) {
		t.Error("wyróżnienie tła przeszło do wydanego .docx — warstwa znakowania " +
			"ma żyć w sesji, nie w pliku")
	}

	// ── Miara trzecia: odmowa nazwana, nie cisza ─────────────────────────────
	nazwane := false
	for _, pozycja := range wynik.DroppedFeatures {
		opis := strings.ToLower(pozycja.Reason)
		if pozycja.Detail != nil {
			opis += " " + strings.ToLower(*pozycja.Detail)
		}
		if strings.Contains(opis, "komentarz") || strings.Contains(opis, "znakowan") {
			nazwane = true
		}
	}
	if !nazwane {
		t.Errorf("wydanie z includeComments nie nazwało w wykazie, że komentarze "+
			"do pliku nie wchodzą; wykaz: %+v", wynik.DroppedFeatures)
	}

	// Sama treść pisma ma zostać nietknięta przez czyszczenie znakowania.
	trescPoObiegu := trescDokumentu(t, zmontowany, zycie, okno+"-powrot", wrocony.Document.Id)
	if !strings.Contains(trescPoObiegu, "Umowa o dzieło") {
		t.Errorf("czyszczenie warstwy znakowania zjadło treść dokumentu; treść: %q",
			trescPoObiegu)
	}
}

// ooxmlCzyStoiWyroznienie rozstrzyga, czy w postaci stoi choć jeden fragment
// z wyróżnieniem tła.
func ooxmlCzyStoiWyroznienie(postac shared.StudioDocumentForm) bool {
	for i := range postac.Blocks {
		for j := range postac.Blocks[i].Runs {
			postacZnaku := postac.Blocks[i].Runs[j].Format
			if postacZnaku == nil || postacZnaku.HighlightColor == nil {
				continue
			}
			if strings.TrimSpace(*postacZnaku.HighlightColor) != "" {
				return true
			}
		}
	}
	return false
}
