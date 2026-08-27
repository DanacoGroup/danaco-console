// Odpowiedzialność pliku: dokument w tłumaczeniu — `translate.document.load`,
// `translate.document.render` i `translate.document.layout.compare`.
//
// Wczytanie dokumentu robi trzy rzeczy naraz i wszystkie trzy są trwałe:
// zakłada wiersz dokumentu, zapisuje jego segmenty i wstawia treść dokumentu
// jako tekst źródłowy okna. Trzecia jest tą, dla której Operator w ogóle
// wczytuje dokument: bez niej `target.add` nie miałby czego przetłumaczyć,
// a moduł meldowałby wczytanie dokumentu, po którym okno zostaje puste.
package core

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// przedrostekDokumentuTlumaczenia znakuje identyfikator dokumentu.
const przedrostekDokumentuTlumaczenia = "dok-"

// zapasDlugosciUkladu mówi, o ile procent przekład może być dłuższy od źródła,
// zanim uznamy, że pole go nie pomieści. Piętnaście procent to zapas, który
// mieści zwykłą różnicę między językami; powyżej tekst realnie wychodzi poza
// ramkę, w której stał oryginał.
const zapasDlugosciUkladu = 115

// granicaRozpoznaniaPisma jest granicą czasu jednego wywołania Tesseracta.
// Skan wielostronicowy rozpoznaje się długo, ale program, który utknął, nie ma
// prawa trzymać żądania bez końca (`zewnetrzne.Wolaj` granicy wymaga).
const granicaRozpoznaniaPisma = 180 * time.Second

// WczytajDokument obsługuje `translate.document.load`.
func (a *adapterTlumaczenia) WczytajDokument(ctx context.Context,
	z shared.TranslateDocumentLoadRequest) (shared.TranslateDocumentLoadResponse, error) {

	sciezka := strings.TrimSpace(napisZeWskaznika(z.Path))
	if sciezka == "" {
		// Kontrakt dopuszcza wskazanie zasobu zamiast ścieżki. Rdzeń modułu
		// Translate nie ma dostępu do magazynu zasobów Designu, więc nazywa to
		// wprost, zamiast oddać pusty dokument z identyfikatorem donikąd.
		if strings.TrimSpace(napisZeWskaznika(z.AssetId)) != "" {
			return shared.TranslateDocumentLoadResponse{}, bladWskazaniaTlumaczenia(
				"moduł Translate wczytuje dokument spod ścieżki; wskazanie zasobu wymaga " +
					"wpięcia magazynu zasobów, którego ten moduł nie ma")
		}
		return shared.TranslateDocumentLoadResponse{}, bladWskazaniaTlumaczenia(
			"żądanie translate.document.load bez ścieżki dokumentu")
	}

	okno, err := a.repozytorium.Okno(ctx, z.WindowId)
	if err != nil {
		return shared.TranslateDocumentLoadResponse{}, bladNieznanegoOkna(z.WindowId, err)
	}

	format := shared.TranslationDocumentFormat("")
	if z.Format != nil {
		format = *z.Format
	}
	if strings.TrimSpace(string(format)) == "" {
		rozpoznany, jest := formatDokumentuZeSciezki(sciezka)
		if !jest {
			return shared.TranslateDocumentLoadResponse{}, bladWskazaniaTlumaczenia(
				"nie da się rozpoznać formatu pliku " + filepath.Base(sciezka) +
					" po końcówce nazwy — wskaż format wprost")
		}
		format = rozpoznany
	}

	segmenty, stron, err := segmentyZDokumentu(sciezka, format)
	if err != nil {
		return shared.TranslateDocumentLoadResponse{}, err
	}

	uzytoOcr := false
	if len(segmenty) == 0 && z.Ocr != nil && *z.Ocr {
		// Dokument bez warstwy tekstowej. Rozpoznanie pisma idzie Tesseraktem
		// z arsenału serwerowego — jedyna droga do treści skanu, i droga wskazana
		// zasadą produktu (nagłówek `*_dokument_formaty.go`).
		rozpoznane, err := a.rozpoznajPismoDokumentu(ctx, sciezka)
		if err != nil {
			return shared.TranslateDocumentLoadResponse{}, err
		}
		segmenty = segmentyZAkapitow(rozdzielAkapity(rozpoznane), "akapit")
		uzytoOcr = true
	}
	if len(segmenty) == 0 {
		return shared.TranslateDocumentLoadResponse{}, bladWskazaniaTlumaczenia(
			"dokument " + filepath.Base(sciezka) + " nie ma warstwy tekstowej; " +
				"powtórz wczytanie z rozpoznaniem pisma (`ocr`)")
	}

	dokument := dane.DokumentTlumaczenia{
		Kod:      nowyIdentyfikator(przedrostekDokumentuTlumaczenia),
		OknoID:   okno.ID,
		Sciezka:  sciezka,
		Format:   string(format),
		UzytoOcr: uzytoOcr,
	}
	if stron > 0 {
		dokument.LiczbaStron = &stron
	}
	zapisany, err := a.repozytorium.ZapiszDokument(ctx, dokument, segmenty)
	if err != nil {
		return shared.TranslateDocumentLoadResponse{}, bladTlumaczenia(err)
	}

	// Treść dokumentu staje się tekstem źródłowym okna — powód w nagłówku pliku.
	tresci := make([]string, 0, len(segmenty))
	for _, segment := range segmenty {
		tresci = append(tresci, segment.Tresc)
	}
	tekst := strings.Join(tresci, "\n\n")
	liczbaSegmentow := int64(len(segmenty))
	if _, err := a.repozytorium.ZapiszOkno(ctx, dane.OknoTlumaczenia{
		Kod:             okno.Kod,
		TekstZrodlowy:   &tekst,
		JezykZrodlowy:   okno.JezykZrodlowy,
		LiczbaSegmentow: &liczbaSegmentow,
	}); err != nil {
		return shared.TranslateDocumentLoadResponse{}, bladTlumaczenia(err)
	}
	// Podział okna idzie po segmentach dokumentu, nie po zdaniach: akapit
	// dokumentu jest jednostką, którą Operator widzi w pliku źródłowym.
	if err := a.repozytorium.UstawSegmentyOkna(ctx, okno.ID, tresci); err != nil {
		return shared.TranslateDocumentLoadResponse{}, bladTlumaczenia(err)
	}

	return shared.TranslateDocumentLoadResponse{
		Document: zlozDokumentTlumaczenia(zapisany, len(segmenty)),
		Segments: zlozSegmentyDokumentu(segmenty),
		UsedOcr:  uzytoOcr,
	}, nil
}

// rozpoznajPismoDokumentu woła Tesseracta na dokumencie bez warstwy tekstowej.
func (a *adapterTlumaczenia) rozpoznajPismoDokumentu(ctx context.Context,
	sciezka string) (string, error) {

	if a.uruchamiacz == nil {
		return "", bladWskazaniaTlumaczenia(
			"rozpoznanie pisma wymaga wpiętego portu uruchamiania procesów, którego rdzeń nie ma")
	}
	okno, zasady, obszar := a.zasiegProgramowTlumaczenia()
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		narzedzieTesseract, []string{sciezka, "stdout", "-l", jezykRozpoznaniaDomyslny},
		katalogPracySyntezy(obszar), granicaRozpoznaniaPisma)
	if err != nil {
		return "", bladWskazaniaTlumaczenia("rozpoznanie pisma nie powiodło się: " + err.Error())
	}
	if strings.TrimSpace(string(wynik.Wyjscie)) == "" {
		return "", bladWskazaniaTlumaczenia(
			"rozpoznanie pisma nie dało ani jednego znaku — dokument nie niesie czytelnego tekstu")
	}
	return string(wynik.Wyjscie), nil
}

// ZlozDokument obsługuje `translate.document.render`. Składa plik wyniku
// z treści panelu, dzieląc ją na tyle akapitów, ile miał dokument źródłowy.
func (a *adapterTlumaczenia) ZlozDokument(ctx context.Context,
	z shared.TranslateDocumentRenderRequest) (shared.TranslateDocumentRenderResponse, error) {

	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.TranslateDocumentRenderResponse{}, bladWskazaniaTlumaczenia(
			"nie ma dokumentu o identyfikatorze " + z.DocumentId)
	}
	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateDocumentRenderResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	tresc, err := trescPanelu(panel)
	if err != nil {
		return shared.TranslateDocumentRenderResponse{}, err
	}

	format := shared.TranslationDocumentFormat(dokument.Format)
	if z.Format != nil && strings.TrimSpace(string(*z.Format)) != "" {
		format = *z.Format
	}
	sciezka := strings.TrimSpace(napisZeWskaznika(z.Path))
	if sciezka == "" {
		sciezka = sciezkaWynikuDokumentu(dokument.Sciezka, panel.Jezyk, format)
	}

	akapity := rozdzielAkapity(tresc)
	if err := zapiszDokumentWyniku(sciezka, format, akapity); err != nil {
		return shared.TranslateDocumentRenderResponse{}, err
	}
	return shared.TranslateDocumentRenderResponse{
		Path:          sciezka,
		RenderedCount: len(akapity),
	}, nil
}

// sciezkaWynikuDokumentu składa ścieżkę wyniku obok materiału, z językiem
// panelu w nazwie — dwa przekłady tego samego dokumentu nie mają się nadpisać.
func sciezkaWynikuDokumentu(material, jezyk string, format shared.TranslationDocumentFormat) string {
	katalog := filepath.Dir(material)
	nazwa := strings.TrimSuffix(filepath.Base(material), filepath.Ext(material))
	return filepath.Join(katalog, nazwa+"."+jezyk+koncowkaFormatuDokumentu(format))
}

// koncowkaFormatuDokumentu daje końcówkę nazwy dla formatu wyniku.
func koncowkaFormatuDokumentu(format shared.TranslationDocumentFormat) string {
	switch format {
	case shared.TranslationDocumentFormatMarkdown:
		return ".md"
	case shared.TranslationDocumentFormatHtml:
		return ".html"
	case shared.TranslationDocumentFormatDocx:
		return ".docx"
	}
	return "." + string(format)
}

// PorownajUklad obsługuje `translate.document.layout.compare`. Zestawia segmenty
// dokumentu z odpowiadającymi im akapitami przekładu i nazywa trzy rzeczy, które
// da się z tego zestawienia stwierdzić: przepełnienie pola, przesunięcie treści
// i element pominięty.
func (a *adapterTlumaczenia) PorownajUklad(ctx context.Context,
	z shared.TranslateDocumentLayoutCompareRequest) (shared.TranslateDocumentLayoutCompareResponse, error) {

	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.TranslateDocumentLayoutCompareResponse{}, bladWskazaniaTlumaczenia(
			"nie ma dokumentu o identyfikatorze " + z.DocumentId)
	}
	panel, err := a.repozytorium.Panel(ctx, z.PanelId)
	if err != nil {
		return shared.TranslateDocumentLayoutCompareResponse{}, bladNieznanegoPanelu(z.PanelId, err)
	}
	tresc, err := trescPanelu(panel)
	if err != nil {
		return shared.TranslateDocumentLayoutCompareResponse{}, err
	}
	segmenty, err := a.repozytorium.SegmentyDokumentu(ctx, dokument.ID)
	if err != nil {
		return shared.TranslateDocumentLayoutCompareResponse{}, bladTlumaczenia(err)
	}
	akapity := rozdzielAkapity(tresc)

	odStrony := liczbaCalkowitaZeWskaznika(z.PageFrom)
	doStrony := liczbaCalkowitaZeWskaznika(z.PageTo)

	roznice := []shared.LayoutDifference{}
	for numer, segment := range segmenty {
		strona := numer + 1
		if segment.Strona != nil {
			strona = int(*segment.Strona)
		}
		if odStrony > 0 && strona < odStrony {
			continue
		}
		if doStrony > 0 && strona > doStrony {
			continue
		}
		if numer >= len(akapity) {
			roznice = append(roznice, shared.LayoutDifference{
				Page:     strona,
				Kind:     shared.LayoutDifferenceKindMissingElement,
				Severity: shared.ProofreadSeverityError,
				Detail:   "akapit obecny w źródle nie ma odpowiednika w przekładzie",
				NodePath: segment.SciezkaWezla,
			})
			continue
		}
		zrodlowy := liczbaZnakow(segment.Tresc)
		docelowy := liczbaZnakow(strings.TrimSpace(akapity[numer]))
		if zrodlowy == 0 {
			continue
		}
		if 100*docelowy/zrodlowy > zapasDlugosciUkladu {
			roznice = append(roznice, shared.LayoutDifference{
				Page:     strona,
				Kind:     shared.LayoutDifferenceKindOverflow,
				Severity: shared.ProofreadSeverityWarning,
				Detail: "przekład dłuższy od źródła o więcej niż " +
					"piętnaście procent — tekst nie zmieści się w polu oryginału",
				NodePath: segment.SciezkaWezla,
			})
		}
	}
	// Akapity nadmiarowe po stronie przekładu przesuwają treść dalej niż
	// w oryginale — to jest przesunięcie strony, nie brak.
	if len(akapity) > len(segmenty) {
		roznice = append(roznice, shared.LayoutDifference{
			Page:     len(segmenty),
			Kind:     shared.LayoutDifferenceKindPageShift,
			Severity: shared.ProofreadSeverityHint,
			Detail:   "przekład ma więcej akapitów niż źródło — dalsza treść przesunie się w składzie",
		})
	}

	porownanych := len(segmenty)
	if dokument.LiczbaStron != nil {
		porownanych = int(*dokument.LiczbaStron)
	}
	return shared.TranslateDocumentLayoutCompareResponse{
		Differences:   roznice,
		ComparedPages: porownanych,
	}, nil
}

// zlozDokumentTlumaczenia przekłada wiersz dokumentu na byt kontraktu.
func zlozDokumentTlumaczenia(dokument dane.DokumentTlumaczenia, segmentow int) shared.TranslationDocument {
	byt := shared.TranslationDocument{
		Id:           dokument.Kod,
		WindowId:     dokument.OknoKod,
		Path:         dokument.Sciezka,
		Format:       shared.TranslationDocumentFormat(dokument.Format),
		SegmentCount: segmentow,
		UsedOcr:      dokument.UzytoOcr,
		UpdatedAt:    dokument.Zaktualizowano,
	}
	if dokument.LiczbaStron != nil {
		stron := int(*dokument.LiczbaStron)
		byt.PageCount = &stron
	}
	return byt
}

// zlozSegmentyDokumentu przekłada wiersze segmentów na byty kontraktu.
func zlozSegmentyDokumentu(segmenty []dane.SegmentDokumentu) []shared.DocumentSegment {
	wykaz := make([]shared.DocumentSegment, 0, len(segmenty))
	for _, segment := range segmenty {
		byt := shared.DocumentSegment{
			Index:    int(segment.Kolejnosc),
			Text:     segment.Tresc,
			NodePath: segment.SciezkaWezla,
			Style:    segment.Styl,
		}
		if segment.Strona != nil {
			strona := int(*segment.Strona)
			byt.Page = &strona
		}
		wykaz = append(wykaz, byt)
	}
	return wykaz
}
