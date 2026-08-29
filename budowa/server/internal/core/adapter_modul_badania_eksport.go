// Odpowiedzialność pliku: Export Panel — złożenie dokumentu raportu,
// wytworzenie pliku w formacie docelowym, podgląd przed eksportem, historia
// eksportów, szablony eksportu i udostępnienie odnośnika.
package core

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// dokumentRaportuBadania niesie złożony dokument raportu wraz z jego postacią
// tekstową, użyteczną przy formatach niosących sam tekst.
type dokumentRaportuBadania struct {
	tytul    string
	markdown string
	wiersze  [][]string
}

// WyeksportujRaport obsługuje komendę eksportu raportu badania: wytwarza
// plik, odkłada go w magazynie i dopiero potem zapisuje ślad zlecenia, który
// wskazuje bajty faktycznie zapisane na dysku.
func (a *adapterBadan) WyeksportujRaport(ctx context.Context,
	z shared.ResearchReportExportRequest) (shared.ResearchReportExportResponse, error) {

	if z.ReportId == "" {
		return shared.ResearchReportExportResponse{}, bladWskazaniaBadan("report.export bez raportu")
	}
	zawartosc := zawartoscEksportuBadania(z.Content)
	if z.TemplateId != nil && strings.TrimSpace(*z.TemplateId) != "" {
		szablon, err := a.szablonEksportuBadania(ctx, *z.TemplateId)
		if err != nil {
			return shared.ResearchReportExportResponse{}, err
		}
		zawartosc = szablon
	}
	dokument, err := a.zlozDokumentRaportuBadania(ctx, z.ReportId, zawartosc)
	if err != nil {
		return shared.ResearchReportExportResponse{}, err
	}
	bajty, err := a.bajtyEksportuBadania(ctx, dokument, z.Format)
	if err != nil {
		return shared.ResearchReportExportResponse{}, err
	}
	sciezka, rozmiar, err := a.odlozMaterialBadania(bajty)
	if err != nil {
		return shared.ResearchReportExportResponse{}, err
	}

	cel := shared.ResearchExportTarget(shared.ResearchExportTargetDownload)
	if z.Target != nil && *z.Target != "" {
		cel = *z.Target
	} else if z.ToLibrary != nil && *z.ToLibrary {
		cel = shared.ResearchExportTarget(shared.ResearchExportTargetLibrary)
	}
	slad := dane.EksportRaportu{
		Kod: nowyIdentyfikator(przedrostekEksportuRaportu), RaportKod: z.ReportId,
		Format: string(z.Format), Cel: string(cel), SciezkaDocelowa: z.TargetPath,
		SciezkaWyniku: &sciezka, RozmiarBajtow: &rozmiar,
	}
	zapisany, err := a.repozytorium.ZapiszEksport(ctx, slad)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ResearchReportExportResponse{}, bladNieznanegoRaportuBadania(z.ReportId)
	}
	if err != nil {
		return shared.ResearchReportExportResponse{}, bladBadan(err)
	}
	return shared.ResearchReportExportResponse{
		Format: z.Format, Path: &sciezka, SizeBytes: &rozmiar,
		ExportId: &zapisany.Kod, Target: &cel,
	}, nil
}

// zawartoscEksportuBadania rozstrzyga skład dokumentu. Brak wskazania w żądaniu
// znaczy pełny skład, czyli dokument ze wszystkimi sekcjami włączonymi.
func zawartoscEksportuBadania(wskazana *shared.ResearchExportContent) shared.ResearchExportContent {
	wlaczone := true
	pelna := shared.ResearchExportContent{
		TitlePage: &wlaczone, TableOfContents: &wlaczone, Bibliography: &wlaczone,
		Footer: &wlaczone, PrismaDiagram: &wlaczone, EvidenceTable: &wlaczone,
	}
	if wskazana == nil {
		return pelna
	}
	return *wskazana
}

// wlaczoneBadania rozstrzyga pole wyboru składu dokumentu; brak wskazania
// w żądaniu znaczy, że dana część jest włączona domyślnie.
func wlaczoneBadania(pole *bool) bool {
	return pole == nil || *pole
}

// zlozDokumentRaportuBadania składa dokument raportu z jego sekcji, bibliografii,
// przesiewu i tabeli dowodów.
func (a *adapterBadan) zlozDokumentRaportuBadania(ctx context.Context, kodRaportu string,
	zawartosc shared.ResearchExportContent) (dokumentRaportuBadania, error) {

	raport, err := a.repozytorium.Raport(ctx, kodRaportu)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return dokumentRaportuBadania{}, bladNieznanegoRaportuBadania(kodRaportu)
	}
	if err != nil {
		return dokumentRaportuBadania{}, bladBadan(err)
	}
	sekcje, err := a.repozytorium.Sekcje(ctx, raport.ID)
	if err != nil {
		return dokumentRaportuBadania{}, bladBadan(err)
	}

	dokument := dokumentRaportuBadania{tytul: raport.Tytul}
	if strings.TrimSpace(dokument.tytul) == "" {
		dokument.tytul = "Raport badawczy"
	}
	var tresc strings.Builder
	if wlaczoneBadania(zawartosc.TitlePage) {
		tresc.WriteString("# " + dokument.tytul + "\n\n")
	}
	if wlaczoneBadania(zawartosc.TableOfContents) && len(sekcje) > 0 {
		tresc.WriteString("## Spis treści\n\n")
		for numer, sekcja := range sekcje {
			tresc.WriteString(strconv.Itoa(numer+1) + ". " + sekcja.Tytul + "\n")
		}
		tresc.WriteString("\n")
	}
	for _, sekcja := range sekcje {
		tresc.WriteString("## " + sekcja.Tytul + "\n\n")
		if sekcja.Tresc != nil && strings.TrimSpace(*sekcja.Tresc) != "" {
			tresc.WriteString(strings.TrimSpace(*sekcja.Tresc) + "\n\n")
		}
	}

	if wlaczoneBadania(zawartosc.PrismaDiagram) {
		liczniki, err := a.licznikiPrismaBadania(ctx, raport.Okno)
		if err != nil {
			return dokumentRaportuBadania{}, err
		}
		tresc.WriteString("## Przesiew PRISMA\n\n")
		for _, wiersz := range []string{
			opisLicznikaBadania("zidentyfikowane", liczniki.Identified),
			opisLicznikaBadania("duplikaty zdjęte", liczniki.DuplicatesRemoved),
			opisLicznikaBadania("przesiane", liczniki.Screened),
			opisLicznikaBadania("wyłączone", liczniki.Excluded),
			opisLicznikaBadania("włączone", liczniki.Included),
		} {
			tresc.WriteString("- " + wiersz + "\n")
		}
		tresc.WriteString("\n")
	}

	if wlaczoneBadania(zawartosc.EvidenceTable) {
		naglowki, wiersze, err := a.zawartoscBlokuBadania(ctx, raport.Okno,
			shared.ResearchReportInsertRequest{Kind: shared.ResearchBlockKindEvidenceTable})
		if err == nil && len(wiersze) > 0 {
			dokument.wiersze = append([][]string{naglowki}, wiersze...)
			tresc.WriteString("## Tabela dowodów\n\n")
			tresc.WriteString("| " + strings.Join(naglowki, " | ") + " |\n")
			tresc.WriteString("|" + strings.Repeat(" --- |", len(naglowki)) + "\n")
			for _, wiersz := range wiersze {
				tresc.WriteString("| " + strings.Join(wiersz, " | ") + " |\n")
			}
			tresc.WriteString("\n")
		}
	}

	if wlaczoneBadania(zawartosc.Bibliography) {
		styl := "apa"
		if zawartosc.CitationStyleId != nil && strings.TrimSpace(*zawartosc.CitationStyleId) != "" {
			styl = *zawartosc.CitationStyleId
		}
		bibliografia, err := a.ZlozBibliografie(ctx, shared.ResearchReportBibliographyRequest{
			ReportId: kodRaportu, StyleId: styl, Scope: shared.ResearchBibliographyScopeAll,
		})
		if err == nil && len(bibliografia.Entries) > 0 {
			tresc.WriteString("## Bibliografia\n\n")
			for _, pozycja := range bibliografia.Entries {
				tresc.WriteString("- " + pozycja + "\n")
			}
			tresc.WriteString("\n")
		}
	}

	if wlaczoneBadania(zawartosc.Footer) {
		tresc.WriteString("---\n\n")
		tresc.WriteString("Raport złożony " + time.Now().Format("2006-01-02") +
			" · sekcji: " + strconv.Itoa(len(sekcje)) + "\n")
	}
	dokument.markdown = tresc.String()
	return dokument, nil
}

// bajtyEksportuBadania wytwarza plik w formacie docelowym, kierując skład do
// procedury właściwej dla tego formatu.
func (a *adapterBadan) bajtyEksportuBadania(ctx context.Context, dokument dokumentRaportuBadania,
	format shared.ExportFormat) ([]byte, error) {

	switch format {
	case shared.ExportFormatMarkdown:
		return []byte(dokument.markdown), nil
	case shared.ExportFormatTxt:
		return []byte(bezZnacznikowMarkdownBadania(dokument.markdown)), nil
	case shared.ExportFormatHtml:
		return []byte(dokumentHtmlBadania(dokument)), nil
	case shared.ExportFormatLatex:
		return []byte(dokumentLatexBadania(dokument)), nil
	case shared.ExportFormatXlsx:
		return arkuszXlsxBadania(dokument)
	case shared.ExportFormatPdf, shared.ExportFormatDocx, shared.ExportFormatPptx:
		return a.bajtyPrzezArsenalBadania(ctx, dokument, format)
	default:
		return nil, bladWskazaniaBadan("format eksportu „" + string(format) +
			"” nie jest jednym z ośmiu formatów kontraktu")
	}
}

// bajtyPrzezArsenalBadania oddaje skład formatu biurowego portowi arsenału
// dokumentowego, który wytwarza plik PDF, DOCX albo PPTX.
func (a *adapterBadan) bajtyPrzezArsenalBadania(ctx context.Context,
	dokument dokumentRaportuBadania, format shared.ExportFormat) ([]byte, error) {

	if a.dokumenty == nil {
		return nil, protokolBladBadania(shared.ErrorCodeInternalError,
			"arsenał dokumentowy nie jest wpięty — naprawa: podpiąć port dokumentów "+
				"przy składaniu serwera")
	}
	zrodlowy := "markdown"
	tresc := dokument.markdown
	wynik, err := a.dokumenty.Przeksztalc(ctx, shared.DocumentConvertRequest{
		Content: &tresc, FromFormat: &zrodlowy, ToFormat: string(format),
	})
	if err != nil {
		return nil, err
	}
	if wynik.Asset.Uri == nil || strings.TrimSpace(*wynik.Asset.Uri) == "" {
		return nil, protokolBladBadania(shared.ErrorCodeInternalError,
			"arsenał oddał zasób bez odwołania do bajtów — eksport nie ma czego zapisać")
	}
	return a.odczytajPlikBadania(*wynik.Asset.Uri)
}

// bezZnacznikowMarkdownBadania sprowadza Markdown do czystego tekstu,
// usuwając znaczniki formatowania z treści sekcji.
func bezZnacznikowMarkdownBadania(tekst string) string {
	wiersze := strings.Split(tekst, "\n")
	czyste := make([]string, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wiersz = strings.TrimLeft(wiersz, "#")
		wiersz = strings.ReplaceAll(wiersz, "**", "")
		czyste = append(czyste, strings.TrimSpace(wiersz))
	}
	return strings.Join(czyste, "\n")
}

// dokumentHtmlBadania składa samodzielny plik HTML raportu, ze stylem
// i treścią osadzonymi w jednym dokumencie.
func dokumentHtmlBadania(dokument dokumentRaportuBadania) string {
	var wynik strings.Builder
	wynik.WriteString("<!doctype html>\n<html lang=\"pl\">\n<head>\n")
	wynik.WriteString("<meta charset=\"utf-8\">\n<title>")
	xml.EscapeText(&wynik, []byte(dokument.tytul))
	wynik.WriteString("</title>\n</head>\n<body>\n")
	for _, wiersz := range strings.Split(dokument.markdown, "\n") {
		czysty := strings.TrimSpace(wiersz)
		if czysty == "" {
			continue
		}
		poziom := len(czysty) - len(strings.TrimLeft(czysty, "#"))
		if poziom > 0 {
			znacznik := "h" + strconv.Itoa(poziom)
			wynik.WriteString("<" + znacznik + ">")
			xml.EscapeText(&wynik, []byte(strings.TrimSpace(czysty[poziom:])))
			wynik.WriteString("</" + znacznik + ">\n")
			continue
		}
		wynik.WriteString("<p>")
		xml.EscapeText(&wynik, []byte(czysty))
		wynik.WriteString("</p>\n")
	}
	wynik.WriteString("</body>\n</html>\n")
	return wynik.String()
}

// dokumentLatexBadania składa źródło LaTeX raportu, gotowe do złożenia
// w dokument przez program zewnętrzny.
func dokumentLatexBadania(dokument dokumentRaportuBadania) string {
	var wynik strings.Builder
	wynik.WriteString("\\documentclass[11pt]{article}\n")
	wynik.WriteString("\\usepackage[utf8]{inputenc}\n\\usepackage[polish]{babel}\n")
	wynik.WriteString("\\title{" + oslonLatexBadania(dokument.tytul) + "}\n")
	wynik.WriteString("\\begin{document}\n\\maketitle\n")
	for _, wiersz := range strings.Split(dokument.markdown, "\n") {
		czysty := strings.TrimSpace(wiersz)
		if czysty == "" {
			continue
		}
		poziom := len(czysty) - len(strings.TrimLeft(czysty, "#"))
		if poziom == 1 {
			continue
		}
		if poziom >= 2 {
			wynik.WriteString("\\section{" + oslonLatexBadania(strings.TrimSpace(czysty[poziom:])) + "}\n")
			continue
		}
		wynik.WriteString(oslonLatexBadania(czysty) + "\n\n")
	}
	wynik.WriteString("\\end{document}\n")
	return wynik.String()
}

// oslonLatexBadania osłania znaki, które w LaTeX-u mają znaczenie sterujące,
// żeby treść raportu nie złamała składu.
func oslonLatexBadania(tekst string) string {
	zamiennik := strings.NewReplacer(
		"\\", "\\textbackslash{}", "&", "\\&", "%", "\\%", "$", "\\$",
		"#", "\\#", "_", "\\_", "{", "\\{", "}", "\\}", "~", "\\textasciitilde{}",
		"^", "\\textasciicircum{}")
	return zamiennik.Replace(tekst)
}

// arkuszXlsxBadania składa arkusz OOXML z tabeli dowodów raportu.
//
// Arkusz powstaje tutaj, a nie programem zewnętrznym, bo XLSX jest spakowanym
// XML-em: cztery pliki w archiwum wystarczą do otwarcia w programie biurowym.
func arkuszXlsxBadania(dokument dokumentRaportuBadania) ([]byte, error) {
	wiersze := dokument.wiersze
	if len(wiersze) == 0 {
		wiersze = [][]string{{"raport"}, {dokument.tytul}}
	}
	var dane bytes.Buffer
	dane.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)
	for numer, wiersz := range wiersze {
		dane.WriteString(`<row r="` + strconv.Itoa(numer+1) + `">`)
		for kolumna, komorka := range wiersz {
			dane.WriteString(`<c r="` + oznaczenieKolumnyBadania(kolumna) + strconv.Itoa(numer+1) +
				`" t="inlineStr"><is><t>`)
			xml.EscapeText(&dane, []byte(komorka))
			dane.WriteString(`</t></is></c>`)
		}
		dane.WriteString(`</row>`)
	}
	dane.WriteString(`</sheetData></worksheet>`)

	var archiwum bytes.Buffer
	zapis := zip.NewWriter(&archiwum)
	pliki := []struct{ nazwa, tresc string }{
		{"[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">` +
			`<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>` +
			`<Default Extension="xml" ContentType="application/xml"/>` +
			`<Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>` +
			`<Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>` +
			`</Types>`},
		{"_rels/.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>` +
			`</Relationships>`},
		{"xl/workbook.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" ` +
			`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">` +
			`<sheets><sheet name="Raport" sheetId="1" r:id="rId1"/></sheets></workbook>`},
		{"xl/_rels/workbook.xml.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
			`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
			`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>` +
			`</Relationships>`},
		{"xl/worksheets/sheet1.xml", dane.String()},
	}
	for _, plik := range pliki {
		ujscie, err := zapis.Create(plik.nazwa)
		if err != nil {
			return nil, bladBadan(err)
		}
		if _, err := ujscie.Write([]byte(plik.tresc)); err != nil {
			return nil, bladBadan(err)
		}
	}
	if err := zapis.Close(); err != nil {
		return nil, bladBadan(err)
	}
	return archiwum.Bytes(), nil
}

// oznaczenieKolumnyBadania przekłada numer kolumny arkusza na jej literowe
// oznaczenie, zgodne z zapisem adresów OOXML.
func oznaczenieKolumnyBadania(numer int) string {
	oznaczenie := ""
	for numer >= 0 {
		oznaczenie = string(rune('A'+numer%26)) + oznaczenie
		numer = numer/26 - 1
	}
	return oznaczenie
}

// ── Podgląd, historia, szablony, udostępnienie ─────────────────────────────

// PodejrzyjEksport obsługuje komendę podglądu eksportu, oddając skład
// dokumentu bez wytwarzania pliku.
func (a *adapterBadan) PodejrzyjEksport(ctx context.Context,
	z shared.ResearchExportPreviewRequest) (shared.ResearchExportPreviewResponse, error) {

	if z.ReportId == "" {
		return shared.ResearchExportPreviewResponse{}, bladWskazaniaBadan("export.preview bez raportu")
	}
	dokument, err := a.zlozDokumentRaportuBadania(ctx, z.ReportId, zawartoscEksportuBadania(z.Content))
	if err != nil {
		return shared.ResearchExportPreviewResponse{}, err
	}
	podglad := dokument.markdown
	obciety := false
	if len([]rune(podglad)) > granicaLekturyBadania {
		podglad = skrocDoBadania(podglad, granicaLekturyBadania)
		obciety = true
	}
	return shared.ResearchExportPreviewResponse{Preview: shared.LibraryPreview{
		FileId: z.ReportId, Kind: shared.LibraryPreviewKindText,
		Text: &podglad, Truncated: &obciety,
	}}, nil
}

// WypiszEksporty obsługuje komendę odczytu historii eksportów raportu
// badania, uporządkowanej od najnowszego.
func (a *adapterBadan) WypiszEksporty(ctx context.Context,
	z shared.ResearchExportListRequest) (shared.ResearchExportListResponse, error) {

	if (z.ReportId == nil || *z.ReportId == "") && (z.WindowId == nil || *z.WindowId == "") {
		return shared.ResearchExportListResponse{}, bladWskazaniaBadan(
			"export.list bez raportu i bez okna — wskaż, czyja historia ma być pokazana")
	}
	limit := 0
	if z.Limit != nil {
		limit = *z.Limit
	}
	eksporty, err := a.repozytorium.Eksporty(ctx, wartoscTekstu(z.ReportId),
		wartoscTekstu(z.WindowId), limit)
	if err != nil {
		return shared.ResearchExportListResponse{}, bladBadan(err)
	}
	przelozone := make([]shared.ResearchExportRecord, 0, len(eksporty))
	for _, eksport := range eksporty {
		cel := shared.ResearchExportTarget(eksport.Cel)
		przelozone = append(przelozone, shared.ResearchExportRecord{
			Id: eksport.Kod, ReportId: eksport.RaportKod,
			Format: shared.ExportFormat(eksport.Format), Target: cel,
			Path: eksport.SciezkaWyniku, LibraryFileId: eksport.PlikBibliotekiID,
			SizeBytes: eksport.RozmiarBajtow, CreatedAt: chwilaBazy(eksport.Utworzono),
		})
	}
	return shared.ResearchExportListResponse{Exports: przelozone}, nil
}

// ZapiszSzablonEksportu obsługuje komendę zapisu szablonu eksportu, utrwalając
// skład dokumentu pod nazwaną pozycją.
func (a *adapterBadan) ZapiszSzablonEksportu(ctx context.Context,
	z shared.ResearchExportTemplateSetRequest) (shared.ResearchExportTemplateSetResponse, error) {

	if strings.TrimSpace(z.Name) == "" {
		return shared.ResearchExportTemplateSetResponse{},
			bladWskazaniaBadan("export.template.set bez nazwy szablonu")
	}
	if z.Format == "" {
		return shared.ResearchExportTemplateSetResponse{},
			bladWskazaniaBadan("export.template.set bez formatu")
	}
	kod := nowyIdentyfikator(przedrostekSzablonuBadania)
	if z.TemplateId != nil && *z.TemplateId != "" {
		kod = *z.TemplateId
	}
	surowa, err := json.Marshal(z.Content)
	if err != nil {
		return shared.ResearchExportTemplateSetResponse{}, bladBadan(err)
	}
	var cel *string
	if z.Target != nil && *z.Target != "" {
		wartosc := string(*z.Target)
		cel = &wartosc
	}
	zapisany, err := a.repozytorium.ZapiszSzablonEksportu(ctx, dane.SzablonEksportuBadania{
		Kod: kod, Nazwa: z.Name, Format: string(z.Format), Cel: cel, Zawartosc: string(surowa),
	})
	if err != nil {
		return shared.ResearchExportTemplateSetResponse{}, bladBadan(err)
	}
	return shared.ResearchExportTemplateSetResponse{
		Template: zlozSzablonEksportuBadania(zapisany),
	}, nil
}

// WypiszSzablonyEksportu obsługuje komendę odczytu szablonów eksportu
// zapisanych dla badania, uporządkowanych po nazwie.
func (a *adapterBadan) WypiszSzablonyEksportu(ctx context.Context,
	_ shared.ResearchExportTemplateListRequest) (shared.ResearchExportTemplateListResponse, error) {

	szablony, err := a.repozytorium.SzablonyEksportu(ctx)
	if err != nil {
		return shared.ResearchExportTemplateListResponse{}, bladBadan(err)
	}
	przelozone := make([]shared.ResearchExportTemplate, 0, len(szablony))
	for _, szablon := range szablony {
		przelozone = append(przelozone, zlozSzablonEksportuBadania(szablon))
	}
	return shared.ResearchExportTemplateListResponse{Templates: przelozone}, nil
}

// zlozSzablonEksportuBadania przekłada wiersz szablonu eksportu na byt
// kontraktu, niosący nazwę i skład dokumentu.
func zlozSzablonEksportuBadania(s dane.SzablonEksportuBadania) shared.ResearchExportTemplate {
	szablon := shared.ResearchExportTemplate{
		Id: s.Kod, Name: s.Nazwa, Format: shared.ExportFormat(s.Format),
	}
	if s.Cel != nil {
		cel := shared.ResearchExportTarget(*s.Cel)
		szablon.Target = &cel
	}
	var zawartosc shared.ResearchExportContent
	if err := json.Unmarshal([]byte(s.Zawartosc), &zawartosc); err == nil {
		szablon.Content = zawartosc
	}
	return szablon
}

// szablonEksportuBadania odczytuje skład dokumentu zapisany w szablonie,
// rozbierając zapis kolumn na pole wyboru.
func (a *adapterBadan) szablonEksportuBadania(ctx context.Context,
	kod string) (shared.ResearchExportContent, error) {

	szablony, err := a.repozytorium.SzablonyEksportu(ctx)
	if err != nil {
		return shared.ResearchExportContent{}, bladBadan(err)
	}
	for _, szablon := range szablony {
		if szablon.Kod == kod {
			return zlozSzablonEksportuBadania(szablon).Content, nil
		}
	}
	return shared.ResearchExportContent{}, protokolBladBadania(shared.ErrorCodeNotFound,
		"szablon eksportu nie istnieje: "+kod)
}

// UdostepnijRaport obsługuje komendę udostępnienia raportu badania. Odnośnik
// wskazuje plik wytworzony w magazynie rdzenia i otwiera dokument, który
// realnie leży na tej maszynie, a nie adres usługi udostępniania.
func (a *adapterBadan) UdostepnijRaport(ctx context.Context,
	z shared.ResearchExportShareRequest) (shared.ResearchExportShareResponse, error) {

	if z.ReportId == "" {
		return shared.ResearchExportShareResponse{}, bladWskazaniaBadan("export.share bez raportu")
	}
	eksporty, err := a.repozytorium.Eksporty(ctx, z.ReportId, "", 1)
	if err != nil {
		return shared.ResearchExportShareResponse{}, bladBadan(err)
	}
	if len(eksporty) == 0 || eksporty[0].SciezkaWyniku == nil {
		return shared.ResearchExportShareResponse{}, protokolBladBadania(shared.ErrorCodeNotFound,
			"raport "+z.ReportId+" nie ma jeszcze wytworzonego pliku — naprawa: "+
				"wykonać research.report.export przed udostępnieniem")
	}
	adres := "file://" + filepath.ToSlash(*eksporty[0].SciezkaWyniku)
	odpowiedz := shared.ResearchExportShareResponse{Url: adres}
	var wygasa *string
	if z.ExpiresInMinutes != nil && *z.ExpiresInMinutes > 0 {
		chwila := time.Now().Add(time.Duration(*z.ExpiresInMinutes) * time.Minute).UTC()
		tekst := chwila.Format("2006-01-02T15:04:05.000Z")
		wygasa = &tekst
		znacznik := chwila.UnixMilli()
		odpowiedz.ExpiresAt = &znacznik
	}
	if err := a.repozytorium.ZapiszUdostepnienie(ctx, z.ReportId,
		nowyIdentyfikator(przedrostekUdostepnieniaBadania), adres, wygasa); err != nil {
		return shared.ResearchExportShareResponse{}, bladBadan(err)
	}
	return odpowiedz, nil
}
