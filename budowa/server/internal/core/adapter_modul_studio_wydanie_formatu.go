// Odpowiedzialność pliku: wydanie dokumentu Studia do formatu docelowego —
// txt, md, docx, odt, pdf, html i rtf — pojedynczo i wsadowo, wraz z wykazem
// cech pominiętych, które format docelowy z postaci dokumentu nie przenosi.
package core

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"context"

	"danacoconsole/shared"
)

// WydajDoFormatu obsługuje `studio.document.export.format`: wydanie jednego
// dokumentu, z narzuceniem postaci szablonu, gdy żądanie szablon wskazało.
func (a *adapterStudia) WydajDoFormatu(ctx context.Context,
	z shared.StudioDocumentExportFormatRequest) (shared.StudioDocumentExportFormatResponse, error) {

	if strings.TrimSpace(z.DocumentId) == "" {
		return shared.StudioDocumentExportFormatResponse{}, bladWskazaniaStudio(
			"wydanie dokumentu bez wskazania dokumentu")
	}
	format, err := wydanieFormatDocelowy(z.Format)
	if err != nil {
		return shared.StudioDocumentExportFormatResponse{}, err
	}
	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioDocumentExportFormatResponse{}, err
	}
	// Szablon narzucony wydaniu podmienia postać, nie treść dokumentu.
	if kod := strings.TrimSpace(wartoscTekstu(z.TemplateId)); kod != "" {
		if err := a.wydaniePrzyjmijPostacSzablonu(ctx, stan, kod); err != nil {
			return shared.StudioDocumentExportFormatResponse{}, err
		}
	}
	wynik, err := a.wydanieDokumentu(ctx, stan, format, z.Path, z.ProfileId,
		wydanieNastawyOdmowy(z.IncludeComments, z.IncludeTrackedChanges))
	if err != nil {
		return shared.StudioDocumentExportFormatResponse{}, err
	}
	return shared.StudioDocumentExportFormatResponse{Result: wynik}, nil
}

// WydajWsadowo obsługuje `studio.document.export.batch`. Odmowa jednego
// dokumentu nie wstrzymuje pozostałych; odmowa całego wsadu zapada dopiero,
// gdy nie wydał się ani jeden dokument.
func (a *adapterStudia) WydajWsadowo(ctx context.Context,
	z shared.StudioDocumentExportBatchRequest) (shared.StudioDocumentExportBatchResponse, error) {

	if len(z.DocumentIds) == 0 {
		return shared.StudioDocumentExportBatchResponse{}, bladWskazaniaStudio(
			"wydanie wsadowe bez ani jednego dokumentu")
	}
	format, err := wydanieFormatDocelowy(z.Format)
	if err != nil {
		return shared.StudioDocumentExportBatchResponse{}, err
	}

	odpowiedz := shared.StudioDocumentExportBatchResponse{
		Results:  []shared.StudioExportResult{},
		Failures: []shared.StudioSkippedItem{},
	}
	for _, kod := range z.DocumentIds {
		czysty := strings.TrimSpace(kod)
		if czysty == "" {
			odpowiedz.Failed++
			odpowiedz.Failures = append(odpowiedz.Failures, shared.StudioSkippedItem{
				Reason: "wskazanie dokumentu przyszło puste",
			})
			continue
		}
		stan, err := a.postacWczytaj(ctx, czysty)
		if err != nil {
			odpowiedz.Failed++
			odpowiedz.Failures = append(odpowiedz.Failures, shared.StudioSkippedItem{
				Reason: "dokumentu nie dało się wczytać",
				Detail: wejscieWskaznikTekstu(czysty + ": " + err.Error()),
			})
			continue
		}
		var sciezka *string
		if katalog := strings.TrimSpace(wartoscTekstu(z.Directory)); katalog != "" {
			nazwa := wydanieNazwaPliku(stan, format)
			pelna := filepath.Join(katalog, nazwa)
			sciezka = &pelna
		}
		wynik, err := a.wydanieDokumentu(ctx, stan, format, sciezka, z.ProfileId, nil)
		if err != nil {
			odpowiedz.Failed++
			odpowiedz.Failures = append(odpowiedz.Failures, shared.StudioSkippedItem{
				Reason: "wydania dokumentu nie dało się złożyć",
				Detail: wejscieWskaznikTekstu(czysty + ": " + err.Error()),
			})
			continue
		}
		odpowiedz.Succeeded++
		odpowiedz.Results = append(odpowiedz.Results, wynik)
	}

	if odpowiedz.Succeeded == 0 {
		powody := make([]string, 0, len(odpowiedz.Failures))
		for _, odmowa := range odpowiedz.Failures {
			powody = append(powody, odmowa.Reason+" — "+wartoscTekstu(odmowa.Detail))
		}
		return shared.StudioDocumentExportBatchResponse{}, wejscieBladBraku(
			"wydanie wsadowe nie wydało ani jednego dokumentu: " + strings.Join(powody, "; "))
	}
	return odpowiedz, nil
}

// wydanieFormatDocelowy sprawdza format wydania wykazem wartości kontraktu
// i oddaje go przycięty i znormalizowany do małych liter, albo odmowę.
func wydanieFormatDocelowy(format shared.StudioExportFormat) (shared.StudioExportFormat, error) {
	czysty := shared.StudioExportFormat(strings.TrimSpace(strings.ToLower(string(format))))
	if czysty == "" {
		return "", bladWskazaniaStudio("wydanie bez wskazania formatu")
	}
	for _, znany := range shared.WartosciStudioExportFormat() {
		if czysty == znany {
			return czysty, nil
		}
	}
	return "", bladWskazaniaStudio("format wydania " + string(format) +
		" nie jest formatem, do którego rdzeń wydaje; rdzeń wydaje txt, md, docx, odt, " +
		"pdf, html i rtf")
}

// wydanieNastawyOdmowy zbiera pozycje pominięte, które wynikają z nastaw
// żądania, a nie z formatu: komentarze i zmiany śledzone nie wchodzą do
// wydania żadnego formatu, bo rachunek składa plik z postaci, nie z adiustacji.
func wydanieNastawyOdmowy(komentarze, zmiany *bool) []shared.StudioSkippedItem {
	pominiete := []shared.StudioSkippedItem{}
	if komentarze != nil && *komentarze {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "komentarze nie weszły do wydania",
			Detail: wejscieWskaznikTekstu("rachunek wydania składa plik z postaci dokumentu; " +
				"komentarze wydaje paczka redakcyjna (studio.package.export)"),
		})
	}
	if zmiany != nil && *zmiany {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "zmiany śledzone nie weszły do wydania jako znaczniki adiustacji",
			Detail: wejscieWskaznikTekstu("wydany plik niesie treść PO zmianach; wykaz zmian " +
				"do przyjęcia oddaje studio.change.list"),
		})
	}
	return pominiete
}

// wydanieZdejmijZnakowanie zdejmuje z postaci warstwę znakowania sesji —
// wyróżnienia tła fragmentów — przed złożeniem pliku i oddaje wykaz tego,
// co zdjęto; styl nazwany z wyróżnieniem zostaje nietknięty.
func wydanieZdejmijZnakowanie(postac *shared.StudioDocumentForm) []shared.StudioSkippedItem {
	zdjetych := 0
	for i := range postac.Blocks {
		for j := range postac.Blocks[i].Runs {
			zdjetych += wydanieZdejmijWyroznienie(postac.Blocks[i].Runs[j].Format)
		}
	}
	for i := range postac.Tables {
		for j := range postac.Tables[i].Cells {
			zdjetych += wydanieZdejmijWyroznienie(postac.Tables[i].Cells[j].Character)
		}
	}
	if zdjetych == 0 {
		return nil
	}
	return []shared.StudioSkippedItem{{
		Reason: "wyróżnienia tła nie weszły do pliku — warstwa znakowania żyje " +
			"w sesji pracy, a nie w piśmie wysyłanym na zewnątrz",
		Detail: wejscieWskaznikTekstu(strconv.Itoa(zdjetych) +
			" wyróżnionych fragmentów; w dokumencie zostają nietknięte"),
	}}
}

// wydanieZdejmijWyroznienie zdejmuje wyróżnienie z jednej postaci znaku i mówi,
// czy było co zdejmować.
func wydanieZdejmijWyroznienie(postac *shared.StudioCharacterFormat) int {
	if postac == nil || postac.HighlightColor == nil {
		return 0
	}
	if strings.TrimSpace(*postac.HighlightColor) == "" {
		postac.HighlightColor = nil
		return 0
	}
	postac.HighlightColor = nil
	return 1
}

// wydaniePrzyjmijPostacSzablonu narzuca dokumentowi postać wzorcową szablonu na
// czas wydania: arkusz stylów, nastawy strony, nagłówek i stopkę. Zmiana idzie
// na kopii postaci wydania i nie zapisuje się do dokumentu zastanego w bazie.
func (a *adapterStudia) wydaniePrzyjmijPostacSzablonu(ctx context.Context, stan *stanPostaci,
	kodSzablonu string) error {

	szablon, err := a.szablonPismaSzczegol(ctx, kodSzablonu)
	if err != nil {
		return err
	}
	if szablon.Form == nil {
		return wejscieBladBraku("szablon „" + szablon.Name +
			"” nie niesie postaci wzorcowej — nie ma czego narzucić wydaniu")
	}
	if szablon.Form.PageSetup != nil {
		stan.forma.PageSetup = szablon.Form.PageSetup
	}
	if len(szablon.Form.Styles) > 0 {
		stan.forma.Styles = szablon.Form.Styles
	}
	if len(szablon.Form.Sections) > 0 && len(stan.forma.Sections) > 0 {
		// Nagłówek i stopka wiszą na sekcji — przejmuje je sekcja pierwsza.
		stan.forma.Sections[0].HeadersFooters = szablon.Form.Sections[0].HeadersFooters
		stan.forma.Sections[0].PageSetup = szablon.Form.Sections[0].PageSetup
	}
	return nil
}

// wydanieNazwaPliku składa nazwę pliku wydania z tytułu dokumentu, oczyszczoną
// ze znaków niedozwolonych w nazwach plików systemu, wraz z rozszerzeniem formatu.
func wydanieNazwaPliku(stan *stanPostaci, format shared.StudioExportFormat) string {
	nazwa := strings.TrimSpace(wartoscTekstu(stan.dokument.Tytul))
	if nazwa == "" {
		nazwa = stan.dokument.Kod
	}
	nazwa = strings.Map(func(znak rune) rune {
		switch znak {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			return '-'
		}
		return znak
	}, nazwa)
	return nazwa + "." + string(format)
}

// wydanieDokumentu składa plik wydania, odkłada bajty do magazynu zasobów
// i oddaje wynik wraz z wykazem cech pominiętych; jedna droga dla wydania
// pojedynczego, wsadowego i zapisu pod nową nazwą.
func (a *adapterStudia) wydanieDokumentu(ctx context.Context, stan *stanPostaci,
	format shared.StudioExportFormat, sciezka, kodProfilu *string,
	dodatkowePominiecia []shared.StudioSkippedItem) (shared.StudioExportResult, error) {

	tresc := wartoscTekstu(stan.dokument.Tresc)
	tytul := strings.TrimSpace(wartoscTekstu(stan.dokument.Tytul))
	if tytul == "" {
		tytul = stan.dokument.Kod
	}

	// Warstwa znakowania sesji schodzi przed złożeniem pliku, dla każdego formatu.
	dodatkowePominiecia = append(dodatkowePominiecia,
		wydanieZdejmijZnakowanie(&stan.forma)...)

	var bajty []byte
	var pominiete []shared.StudioSkippedItem
	var err error

	switch format {
	case shared.StudioExportFormatTxt:
		bajty, pominiete = wydanieTekstem(&stan.forma, tresc)
	case shared.StudioExportFormatMd:
		bajty, pominiete = wydanieMarkdownem(&stan.forma, tresc)
	case shared.StudioExportFormatHtml:
		bajty, pominiete = wydanieHtmlem(&stan.forma, tresc, tytul)
	case shared.StudioExportFormatRtf:
		bajty, pominiete = wydanieRtfem(&stan.forma, tresc)
	case shared.StudioExportFormatDocx:
		bajty, pominiete, err = wejscieZlozOoxml(&stan.forma, tresc, tytul, false)
	case shared.StudioExportFormatOdt:
		bajty, pominiete, err = wejscieZlozOdf(&stan.forma, tresc, tytul, false)
	case shared.StudioExportFormatPdf:
		bajty, pominiete, err = a.wydaniePdfem(ctx, stan, kodProfilu)
	default:
		return shared.StudioExportResult{}, bladWskazaniaStudio(
			"format wydania " + string(format) + " nie ma rachunku w rdzeniu")
	}
	if err != nil {
		return shared.StudioExportResult{}, err
	}
	if len(bajty) == 0 {
		// Dokument bez treści i bez postaci nie ma czego wydać — odmowa, nie plik pusty.
		return shared.StudioExportResult{}, wejscieBladBraku("dokument " +
			stan.dokument.Kod + " jest pusty — nie ma czego wydać do formatu " +
			string(format))
	}
	pominiete = append(pominiete, dodatkowePominiecia...)

	wynik := shared.StudioExportResult{
		DocumentId:      stan.dokument.Kod,
		Format:          format,
		Bytes:           wejscieWskaznikDlugi(int64(len(bajty))),
		DroppedFeatures: pominiete,
	}
	wynik.Note = wejscieWskaznikTekstu(wydanieZdanieBilansu(format, pominiete))

	zasob, err := a.odlozTrescStudia(ctx, bajty, wydanieNazwaPliku(stan, format),
		string(format), stan.dokument.Okno)
	if err != nil {
		return shared.StudioExportResult{}, err
	}
	wynik.AssetId = wejscieWskaznikTekstu(zasob.Id)

	if wskazana := strings.TrimSpace(wartoscTekstu(sciezka)); wskazana != "" {
		if err := wydanieZapiszPlik(wskazana, bajty); err != nil {
			return shared.StudioExportResult{}, err
		}
		wynik.Path = wejscieWskaznikTekstu(wskazana)
	}
	return wynik, nil
}

// wydanieZapiszPlik zapisuje bajty wydania na maszynie serwera.
//
// Jedna droga zapisu do pliku dla wydania dokumentu i dla oddania szablonu:
// dwie dawałyby dwa różne zachowania przy braku katalogu.
func wydanieZapiszPlik(sciezka string, bajty []byte) error {
	if err := os.MkdirAll(filepath.Dir(sciezka), 0o700); err != nil {
		return wejscieBladZaplecza("katalogu " + filepath.Dir(sciezka) +
			" nie da się założyć na maszynie serwera: " + err.Error())
	}
	if err := os.WriteFile(sciezka, bajty, 0o600); err != nil {
		return wejscieBladZaplecza("pliku " + sciezka +
			" nie da się zapisać na maszynie serwera: " + err.Error())
	}
	return nil
}

// wydanieZdanieBilansu składa zdanie o tym, co odpadło przy wydaniu do formatu,
// albo zdanie, że nic nie odpadło, gdy wykaz cech pominiętych jest pusty.
func wydanieZdanieBilansu(format shared.StudioExportFormat,
	pominiete []shared.StudioSkippedItem) string {

	if len(pominiete) == 0 {
		return "wydanie do " + string(format) + " przeniosło całą postać dokumentu; " +
			"nic nie odpadło"
	}
	powody := make([]string, 0, len(pominiete))
	for _, pozycja := range pominiete {
		powody = append(powody, pozycja.Reason)
	}
	return "wydanie do " + string(format) + " pominęło " + strconv.Itoa(len(pominiete)) +
		" cech dokumentu: " + strings.Join(powody, "; ")
}

// ── Wydanie tekstem czystym ─────────────────────────────────────────────────

// wydanieTekstem składa plik tekstowy. Postać schodzi CAŁA i odpowiedź mówi to
// wprost — z rozmiarem straty, nie zdaniem ogólnym.
func wydanieTekstem(postac *shared.StudioDocumentForm, tresc string) ([]byte,
	[]shared.StudioSkippedItem) {

	pominiete := wydaniePostacSchodzi(postac, "tekst czysty")
	var budowa strings.Builder
	for _, blok := range postac.Blocks {
		switch blok.Kind {
		case wejscieRodzajBlokuTabela:
			// Tabela wychodzi wierszami z komórkami rozdzielonymi tabulatorem.
			budowa.WriteString(wydanieTabelaTekstem(postac, blok.TableId))
		case wejscieRodzajBlokuObiekt:
			if opis := wydanieOpisObiektu(postac, blok.ObjectId); opis != "" {
				budowa.WriteString("[" + opis + "]\n")
			}
		case wejscieRodzajBlokuPodzial:
			budowa.WriteString("\n")
		default:
			budowa.WriteString(postacTekstBloku(blok) + "\n")
		}
	}
	wynik := budowa.String()
	if strings.TrimSpace(wynik) == "" {
		wynik = tresc
	}
	if len(postac.Apparatus) > 0 {
		// Przypisy i podpisy wychodzą wykazem na końcu, nie w miejscu odsyłacza.
		budowaAparatu := strings.Builder{}
		budowaAparatu.WriteString("\n")
		for _, element := range postac.Apparatus {
			etykieta := strings.TrimSpace(wartoscTekstu(element.Number) + " " +
				wartoscTekstu(element.Label))
			tekst := strings.TrimSpace(wartoscTekstu(element.Text))
			if etykieta == "" && tekst == "" {
				continue
			}
			budowaAparatu.WriteString(strings.TrimSpace(etykieta+" "+tekst) + "\n")
		}
		wynik += budowaAparatu.String()
	}
	return []byte(wynik), pominiete
}

// wydaniePostacSchodzi składa wykaz cech, które w formacie tekstowym odpadają,
// z liczbami policzonymi z postaci dokumentu, a nie oszacowanymi.
func wydaniePostacSchodzi(postac *shared.StudioDocumentForm,
	nazwaFormatu string) []shared.StudioSkippedItem {

	pominiete := []shared.StudioSkippedItem{{
		Reason: "postać znaku i akapitu nie weszła — " + nazwaFormatu +
			" nie niesie krojów, wcięć ani wyrównania",
	}}
	if len(postac.Styles) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "arkusz stylów nazwanych nie wszedł",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(postac.Styles)) + " stylów"),
		})
	}
	if postac.PageSetup != nil {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "nastawy strony nie weszły — format nośnika, marginesy, nagłówek i stopka",
		})
	}
	if len(postac.Tables) > 0 {
		// Treść komórek wychodzi, siatka tabeli nie — powiedziane z liczbą.
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "siatka tabel nie weszła — komórki wyszły wierszami rozdzielonymi " +
				"tabulatorem, bez obramowania, szerokości kolumn i scaleń",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(postac.Tables)) + " tabel"),
		})
	}
	if len(postac.Objects) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "obrazy i kształty wyszły tekstem zastępczym",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(postac.Objects)) + " obiektów"),
		})
	}
	if len(postac.Fields) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "pola dokumentu wyszły wartością, bez możliwości odświeżenia",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(postac.Fields)) + " pól"),
		})
	}
	pominiete = append(pominiete, wydanieAparatSchodzi(postac, nazwaFormatu)...)
	return pominiete
}

// wydanieAparatSchodzi nazywa stratę na aparacie dokumentu, policzoną po
// rodzaju elementu (przypisy, spisy, odsyłacze, podpisy), a nie jednym workiem.
func wydanieAparatSchodzi(postac *shared.StudioDocumentForm,
	nazwaFormatu string) []shared.StudioSkippedItem {

	if len(postac.Apparatus) == 0 {
		return nil
	}
	przypisow, spisow, odsylaczy, podpisow, innych := 0, 0, 0, 0, 0
	for _, element := range postac.Apparatus {
		switch element.Kind {
		case shared.StudioApparatusKindFootnote, shared.StudioApparatusKindEndnote:
			przypisow++
		case shared.StudioApparatusKindToc, shared.StudioApparatusKindFigureIndex,
			shared.StudioApparatusKindTableIndex, shared.StudioApparatusKindIndex:
			spisow++
		case shared.StudioApparatusKindHyperlink, shared.StudioApparatusKindCrossReference,
			shared.StudioApparatusKindBookmark:
			odsylaczy++
		case shared.StudioApparatusKindCaption:
			podpisow++
		default:
			innych++
		}
	}

	pominiete := make([]shared.StudioSkippedItem, 0, 4)
	if przypisow > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "przypisy zeszły ze swoich miejsc — " + nazwaFormatu +
				" nie niesie przypisu pod stroną, więc wyszły wykazem na końcu, " +
				"bez odsyłacza w treści i bez numeracji odświeżalnej",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(przypisow) + " przypisów"),
		})
	}
	if spisow > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "spisy wyszły stanem na chwilę wydania — " + nazwaFormatu +
				" nie niesie spisu odświeżalnego ani numerów stron",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(spisow) + " spisów"),
		})
	}
	if odsylaczy > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "odsyłacze, odwołania wzajemne i zakładki nie weszły jako " +
				"powiązania — zostało samo brzmienie",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(odsylaczy) + " odsyłaczy"),
		})
	}
	if podpisow > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "podpisy ilustracji i tabel straciły numerację odświeżalną",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(podpisow) + " podpisów"),
		})
	}
	if innych > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "pozostałe elementy aparatu dokumentu wyszły treścią, bez powiązań",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(innych) + " elementów"),
		})
	}
	return pominiete
}

// wydanieTabelaTekstem składa tabelę wierszami rozdzielonymi tabulatorem, wraz
// z podpisem tabeli, gdy dokument go niesie.
func wydanieTabelaTekstem(postac *shared.StudioDocumentForm, kodTabeli *string) string {
	tabela := wydanieTabela(postac, kodTabeli)
	if tabela == nil {
		return ""
	}
	var budowa strings.Builder
	for wiersz := 0; wiersz < tabela.Rows; wiersz++ {
		komorki := make([]string, 0, tabela.Columns)
		for kolumna := 0; kolumna < tabela.Columns; kolumna++ {
			komorka := ooxmlSzukajKomorki(tabela, wiersz, kolumna)
			if komorka == nil {
				komorki = append(komorki, "")
				continue
			}
			komorki = append(komorki, strings.ReplaceAll(
				strings.TrimSpace(wartoscTekstu(komorka.Text)), "\n", " "))
		}
		budowa.WriteString(strings.Join(komorki, "\t") + "\n")
	}
	if tabela.Caption != nil && strings.TrimSpace(*tabela.Caption) != "" {
		budowa.WriteString(strings.TrimSpace(*tabela.Caption) + "\n")
	}
	return budowa.String()
}

// wydanieTabela wyszukuje tabelę postaci po jej identyfikatorze, albo oddaje
// nic, gdy wskazanie jest puste albo tabeli o tym identyfikatorze nie ma.
func wydanieTabela(postac *shared.StudioDocumentForm,
	kodTabeli *string) *shared.StudioDocumentTable {

	if kodTabeli == nil {
		return nil
	}
	for i := range postac.Tables {
		if postac.Tables[i].Id == *kodTabeli {
			return &postac.Tables[i]
		}
	}
	return nil
}

// wydanieOpisObiektu oddaje tekst zastępczy obiektu osadzonego, wraz z jego
// podpisem, gdy dokument podpis niesie.
func wydanieOpisObiektu(postac *shared.StudioDocumentForm, kodObiektu *string) string {
	if kodObiektu == nil {
		return ""
	}
	for _, obiekt := range postac.Objects {
		if obiekt.Id != *kodObiektu {
			continue
		}
		opis := strings.TrimSpace(wartoscTekstu(obiekt.AltText))
		if opis == "" {
			opis = "obiekt osadzony"
		}
		if obiekt.Caption != nil && strings.TrimSpace(*obiekt.Caption) != "" {
			opis += ": " + strings.TrimSpace(*obiekt.Caption)
		}
		return opis
	}
	return ""
}

// ── Wydanie markdownem ──────────────────────────────────────────────────────

// wydanieMarkdownem składa plik markdown: style nazwane na znaczniki, tabele na
// tabele markdown; komórka scalona wychodzi w wykazie pominiętych z liczbą scaleń.
func wydanieMarkdownem(postac *shared.StudioDocumentForm, tresc string) ([]byte,
	[]shared.StudioSkippedItem) {

	pominiete := []shared.StudioSkippedItem{}
	if postac.PageSetup != nil {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "nastawy strony nie weszły — markdown nie niesie nośnika, marginesów, " +
				"nagłówka ani stopki",
		})
	}
	if len(postac.Styles) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "arkusz stylów wyszedł znacznikami markdown, bez własnych postaci stylów",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(postac.Styles)) + " stylów; " +
				"nagłówki wyszły kratkami, cytat znakiem większości"),
		})
	}
	pominiete = append(pominiete, wydanieAparatSchodzi(postac, "markdown")...)

	var budowa strings.Builder
	for _, blok := range postac.Blocks {
		switch blok.Kind {
		case wejscieRodzajBlokuTabela:
			tekst, scalen := wydanieTabelaMarkdownem(postac, blok.TableId)
			budowa.WriteString(tekst)
			if scalen > 0 {
				pominiete = append(pominiete, shared.StudioSkippedItem{
					Reason: "komórki scalone nie weszły — markdown scaleń nie niesie",
					Detail: wejscieWskaznikTekstu(strconv.Itoa(scalen) +
						" scaleń rozpisano na komórki osobne"),
				})
			}
		case wejscieRodzajBlokuObiekt:
			if opis := wydanieOpisObiektu(postac, blok.ObjectId); opis != "" {
				budowa.WriteString("![" + opis + "]()\n\n")
			}
		case wejscieRodzajBlokuPodzial:
			budowa.WriteString("---\n\n")
		default:
			budowa.WriteString(wydanieAkapitMarkdownem(blok))
		}
	}
	wynik := budowa.String()
	if strings.TrimSpace(wynik) == "" {
		wynik = tresc
	}
	if len(postac.Apparatus) > 0 {
		wynik += wydanieAparatMarkdownem(postac)
	}
	if len(postac.Objects) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "obrazy wyszły odsyłaczem bez adresu — bajty leżą w magazynie zasobów " +
				"rdzenia, a plik markdown ich nie niesie",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(postac.Objects)) + " obiektów"),
		})
	}
	return []byte(wynik), pominiete
}

// wydanieAkapitMarkdownem składa akapit wraz z postacią znaku, rozstrzygając
// znacznik nagłówka po poziomie konspektu i pozycję listy po jej identyfikatorze.
func wydanieAkapitMarkdownem(blok shared.StudioDocumentBlock) string {
	tekst := strings.Builder{}
	for _, fragment := range blok.Runs {
		tekst.WriteString(wydanieFragmentMarkdownem(fragment))
	}
	tresc := tekst.String()
	if strings.TrimSpace(tresc) == "" {
		return "\n"
	}

	poziom := 0
	stylNazwany := ""
	listaKod := ""
	if blok.Paragraph != nil {
		if blok.Paragraph.OutlineLevel != nil {
			poziom = *blok.Paragraph.OutlineLevel
		}
		stylNazwany = strings.TrimSpace(wartoscTekstu(blok.Paragraph.StyleName))
		listaKod = strings.TrimSpace(wartoscTekstu(blok.Paragraph.ListId))
	}
	if poziom == 0 && strings.HasPrefix(stylNazwany, "naglowek-") {
		if numer, err := strconv.Atoi(strings.TrimPrefix(stylNazwany, "naglowek-")); err == nil {
			poziom = numer
		}
	}
	switch {
	case poziom > 0:
		if poziom > 6 {
			poziom = 6
		}
		return strings.Repeat("#", poziom) + " " + tresc + "\n\n"
	case stylNazwany == wejscieStylCytat:
		return "> " + tresc + "\n\n"
	case listaKod != "":
		wciecie := 0
		if blok.Paragraph.ListLevel != nil && *blok.Paragraph.ListLevel > 1 {
			wciecie = (*blok.Paragraph.ListLevel - 1) * 2
		}
		return strings.Repeat(" ", wciecie) + "- " + tresc + "\n"
	}
	return tresc + "\n\n"
}

// wydanieFragmentMarkdownem otacza fragment znacznikami postaci znaku: pogrubienie,
// kursywę i przekreślenie, jedyne odmiany, które markdown niesie znacznikiem.
func wydanieFragmentMarkdownem(fragment shared.StudioDocumentRun) string {
	tekst := fragment.Text
	if tekst == "" || fragment.Format == nil {
		return tekst
	}
	if fragment.Format.Bold != nil && *fragment.Format.Bold {
		tekst = "**" + tekst + "**"
	}
	if fragment.Format.Italic != nil && *fragment.Format.Italic {
		tekst = "*" + tekst + "*"
	}
	if fragment.Format.Strikethrough != nil && *fragment.Format.Strikethrough {
		tekst = "~~" + tekst + "~~"
	}
	return tekst
}

// wydanieTabelaMarkdownem składa tabelę markdown wraz z podpisem i oddaje
// liczbę scaleń, których format markdown nie niesie i rozpisuje na komórki osobne.
func wydanieTabelaMarkdownem(postac *shared.StudioDocumentForm,
	kodTabeli *string) (string, int) {

	tabela := wydanieTabela(postac, kodTabeli)
	if tabela == nil {
		return "", 0
	}
	scalen := 0
	var budowa strings.Builder
	for wiersz := 0; wiersz < tabela.Rows; wiersz++ {
		komorki := make([]string, 0, tabela.Columns)
		for kolumna := 0; kolumna < tabela.Columns; kolumna++ {
			komorka := ooxmlSzukajKomorki(tabela, wiersz, kolumna)
			if komorka == nil {
				komorki = append(komorki, "")
				continue
			}
			if komorka.ColumnSpan != nil && *komorka.ColumnSpan > 1 ||
				komorka.RowSpan != nil && *komorka.RowSpan > 1 {
				scalen++
			}
			tekst := strings.ReplaceAll(
				strings.TrimSpace(wartoscTekstu(komorka.Text)), "\n", " ")
			komorki = append(komorki, strings.ReplaceAll(tekst, "|", "\\|"))
		}
		budowa.WriteString("| " + strings.Join(komorki, " | ") + " |\n")
		if wiersz == 0 {
			rozdzielacze := make([]string, tabela.Columns)
			for i := range rozdzielacze {
				rozdzielacze[i] = "---"
			}
			budowa.WriteString("| " + strings.Join(rozdzielacze, " | ") + " |\n")
		}
	}
	budowa.WriteString("\n")
	if tabela.Caption != nil && strings.TrimSpace(*tabela.Caption) != "" {
		budowa.WriteString("*" + strings.TrimSpace(*tabela.Caption) + "*\n\n")
	}
	return budowa.String(), scalen
}

// wydanieAparatMarkdownem składa aparat dokumentu wykazem na końcu pliku, wraz
// z adresem docelowym elementów, które adres niosą.
func wydanieAparatMarkdownem(postac *shared.StudioDocumentForm) string {
	var budowa strings.Builder
	budowa.WriteString("\n---\n\n")
	for _, element := range postac.Apparatus {
		etykieta := strings.TrimSpace(wartoscTekstu(element.Number) + " " +
			wartoscTekstu(element.Label))
		tekst := strings.TrimSpace(wartoscTekstu(element.Text))
		if etykieta == "" && tekst == "" {
			continue
		}
		budowa.WriteString("- " + strings.TrimSpace(etykieta+" "+tekst))
		if element.TargetUrl != nil && strings.TrimSpace(*element.TargetUrl) != "" {
			budowa.WriteString(" <" + strings.TrimSpace(*element.TargetUrl) + ">")
		}
		budowa.WriteString("\n")
	}
	return budowa.String()
}

// ── Wydanie HTML ────────────────────────────────────────────────────────────

// wydanieHtmlem składa plik HTML wraz z arkuszem stylów: style nazwane wychodzą
// zasadami arkusza, a nie postacią wpisaną przy każdym akapicie — inaczej zmiana
// stylu w dokumencie nie miałaby w wydanym pliku ani jednego odpowiednika.
func wydanieHtmlem(postac *shared.StudioDocumentForm, tresc, tytul string) ([]byte,
	[]shared.StudioSkippedItem) {

	pominiete := []shared.StudioSkippedItem{}
	if postac.PageSetup != nil {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "nastawy nośnika wyszły regułą druku, nie układem strony",
			Detail: wejscieWskaznikTekstu("HTML nie ma stron: marginesy i format nośnika " +
				"wyszły zasadą @page, którą wykonuje dopiero druk przeglądarki"),
		})
	}
	if len(postac.Objects) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "obrazy wyszły znacznikiem bez adresu — bajty leżą w magazynie zasobów",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(postac.Objects)) + " obiektów"),
		})
	}

	var budowa strings.Builder
	budowa.WriteString("<!doctype html>\n<html lang=\"pl\">\n<head>\n")
	budowa.WriteString("<meta charset=\"utf-8\">\n<title>" + ooxmlZabezpiecz(tytul) +
		"</title>\n<style>\n")
	budowa.WriteString(wydanieArkuszHtml(postac))
	budowa.WriteString("</style>\n</head>\n<body>\n")

	for _, blok := range postac.Blocks {
		switch blok.Kind {
		case wejscieRodzajBlokuTabela:
			budowa.WriteString(wydanieTabelaHtml(postac, blok.TableId))
		case wejscieRodzajBlokuObiekt:
			if opis := wydanieOpisObiektu(postac, blok.ObjectId); opis != "" {
				budowa.WriteString("<figure><figcaption>" + ooxmlZabezpiecz(opis) +
					"</figcaption></figure>\n")
			}
		case wejscieRodzajBlokuPodzial:
			budowa.WriteString("<hr class=\"podzial-strony\">\n")
		default:
			budowa.WriteString(wydanieAkapitHtml(blok))
		}
	}
	if len(postac.Apparatus) > 0 {
		budowa.WriteString("<section class=\"aparat-dokumentu\">\n")
		for _, element := range postac.Apparatus {
			etykieta := strings.TrimSpace(wartoscTekstu(element.Number) + " " +
				wartoscTekstu(element.Label) + " " + wartoscTekstu(element.Text))
			if etykieta == "" {
				continue
			}
			budowa.WriteString("<p>" + ooxmlZabezpiecz(etykieta) + "</p>\n")
		}
		budowa.WriteString("</section>\n")
	}
	budowa.WriteString("</body>\n</html>\n")

	wynik := budowa.String()
	if len(postac.Blocks) == 0 && strings.TrimSpace(tresc) != "" {
		wynik = strings.Replace(wynik, "<body>\n", "<body>\n<pre>"+
			ooxmlZabezpiecz(tresc)+"</pre>\n", 1)
	}
	return []byte(wynik), pominiete
}

// wydanieArkuszHtml składa arkusz stylów pliku HTML z arkusza dokumentu, wraz
// z regułą strony druku i klasą na każdy styl nazwany, który niesie cechę.
func wydanieArkuszHtml(postac *shared.StudioDocumentForm) string {
	var budowa strings.Builder
	if postac.PageSetup != nil {
		nastawy := postac.PageSetup
		budowa.WriteString("@page { ")
		if nastawy.PageSize != nil {
			budowa.WriteString("size: " + wydanieRozmiarHtml(nastawy) + "; ")
		}
		budowa.WriteString("margin: " + wydanieMarginesHtml(nastawy.MarginTop) + " " +
			wydanieMarginesHtml(nastawy.MarginRight) + " " +
			wydanieMarginesHtml(nastawy.MarginBottom) + " " +
			wydanieMarginesHtml(nastawy.MarginLeft) + "; }\n")
	}
	budowa.WriteString("body { font-family: 'Times New Roman', serif; }\n")
	budowa.WriteString(".podzial-strony { break-after: page; border: 0; }\n")
	for _, styl := range postac.Styles {
		zasady := wydanieZasadyStyluHtml(styl)
		if zasady == "" {
			continue
		}
		budowa.WriteString("." + wydanieKlasaStyluHtml(styl.Name) + " { " + zasady + "}\n")
	}
	return budowa.String()
}

// wydanieRozmiarHtml składa rozmiar nośnika reguły druku, z A4 pionowym jako
// nastawą domyślną, gdy dokument nastaw strony nie niesie.
func wydanieRozmiarHtml(nastawy *shared.StudioPageSetup) string {
	nazwa := strings.TrimSpace(wartoscTekstu(nastawy.PageSize))
	if nazwa == "" {
		nazwa = "A4"
	}
	kierunek := "portrait"
	if nastawy.Orientation != nil && *nastawy.Orientation == shared.StudioPageOrientationPozioma {
		kierunek = "landscape"
	}
	return nazwa + " " + kierunek
}

// wydanieMarginesHtml składa margines reguły druku w milimetrach, oddając
// dwadzieścia pięć milimetrów, gdy dokument tego marginesu nie ustawił.
func wydanieMarginesHtml(wskazanie *int) string {
	if wskazanie == nil {
		return "25mm"
	}
	return strconv.Itoa(*wskazanie) + "mm"
}

// wydanieKlasaStyluHtml składa nazwę klasy arkusza z nazwy stylu, sprowadzoną
// do małych liter i cyfr, ze spacją zamienioną na łącznik.
func wydanieKlasaStyluHtml(nazwa string) string {
	return "styl-" + strings.Map(func(znak rune) rune {
		switch {
		case znak >= 'a' && znak <= 'z', znak >= '0' && znak <= '9', znak == '-':
			return znak
		case znak >= 'A' && znak <= 'Z':
			return znak + 32
		case znak == ' ' || znak == '_':
			return '-'
		}
		return -1
	}, nazwa)
}

// wydanieZasadyStyluHtml składa zasady arkusza z postaci stylu nazwanego: krój,
// stopień, wagę, kursywę i barwę znaku, wyrównanie i odstępy akapitu.
func wydanieZasadyStyluHtml(styl shared.StudioNamedStyle) string {
	var budowa strings.Builder
	if znak := styl.Character; znak != nil {
		if znak.FontFamily != nil && *znak.FontFamily != "" {
			budowa.WriteString("font-family: '" + *znak.FontFamily + "'; ")
		}
		if znak.FontSizePt != nil {
			budowa.WriteString("font-size: " +
				strconv.FormatFloat(*znak.FontSizePt, 'f', 1, 64) + "pt; ")
		}
		if znak.Bold != nil && *znak.Bold {
			budowa.WriteString("font-weight: bold; ")
		}
		if znak.Italic != nil && *znak.Italic {
			budowa.WriteString("font-style: italic; ")
		}
		if znak.Color != nil && *znak.Color != "" {
			budowa.WriteString("color: #" + ooxmlBezKrzyzyka(*znak.Color) + "; ")
		}
	}
	if akapit := styl.Paragraph; akapit != nil {
		if akapit.Align != nil {
			budowa.WriteString("text-align: " + wydanieWyrownanieHtml(*akapit.Align) + "; ")
		}
		if akapit.IndentLeftMm != nil {
			budowa.WriteString("margin-left: " +
				strconv.FormatFloat(*akapit.IndentLeftMm, 'f', 1, 64) + "mm; ")
		}
		if akapit.FirstLineIndentMm != nil {
			budowa.WriteString("text-indent: " +
				strconv.FormatFloat(*akapit.FirstLineIndentMm, 'f', 1, 64) + "mm; ")
		}
		if akapit.SpaceBeforePt != nil {
			budowa.WriteString("margin-top: " +
				strconv.FormatFloat(*akapit.SpaceBeforePt, 'f', 1, 64) + "pt; ")
		}
		if akapit.SpaceAfterPt != nil {
			budowa.WriteString("margin-bottom: " +
				strconv.FormatFloat(*akapit.SpaceAfterPt, 'f', 1, 64) + "pt; ")
		}
		if akapit.LineSpacingRule != nil && akapit.LineSpacingValue != nil &&
			*akapit.LineSpacingRule == shared.StudioLineSpacingRuleMultiple {
			budowa.WriteString("line-height: " +
				strconv.FormatFloat(*akapit.LineSpacingValue, 'f', 2, 64) + "; ")
		}
	}
	return budowa.String()
}

// wydanieWyrownanieHtml przekłada wyrównanie kontraktu na zasadę arkusza CSS,
// oddając wyrównanie do lewej jako wartość domyślną.
func wydanieWyrownanieHtml(wyrownanie shared.StudioTextAlign) string {
	switch wyrownanie {
	case shared.StudioTextAlignRight:
		return "right"
	case shared.StudioTextAlignCenter:
		return "center"
	case shared.StudioTextAlignJustify:
		return "justify"
	default:
		return "left"
	}
}

// wydanieAkapitHtml składa akapit HTML wraz z klasą stylu nazwanego, wychodząc
// znacznikiem nagłówka, gdy blok niesie poziom konspektu.
func wydanieAkapitHtml(blok shared.StudioDocumentBlock) string {
	tekst := strings.Builder{}
	for _, fragment := range blok.Runs {
		tekst.WriteString(wydanieFragmentHtml(fragment))
	}
	tresc := tekst.String()
	if strings.TrimSpace(tresc) == "" {
		return "<p></p>\n"
	}

	klasa, poziom := "", 0
	if blok.Paragraph != nil {
		if nazwa := strings.TrimSpace(wartoscTekstu(blok.Paragraph.StyleName)); nazwa != "" {
			klasa = " class=\"" + wydanieKlasaStyluHtml(nazwa) + "\""
		}
		if blok.Paragraph.OutlineLevel != nil {
			poziom = *blok.Paragraph.OutlineLevel
		}
	}
	if poziom > 0 {
		if poziom > 6 {
			poziom = 6
		}
		znacznik := "h" + strconv.Itoa(poziom)
		return "<" + znacznik + klasa + ">" + tresc + "</" + znacznik + ">\n"
	}
	return "<p" + klasa + ">" + tresc + "</p>\n"
}

// wydanieFragmentHtml składa fragment tekstu wraz z postacią znaku, otaczając
// tekst znacznikami zagnieżdżonymi zgodnie z cechami, które fragment niesie.
func wydanieFragmentHtml(fragment shared.StudioDocumentRun) string {
	tekst := ooxmlZabezpiecz(fragment.Text)
	if tekst == "" || fragment.Format == nil {
		return tekst
	}
	postac := fragment.Format
	if postac.Bold != nil && *postac.Bold {
		tekst = "<strong>" + tekst + "</strong>"
	}
	if postac.Italic != nil && *postac.Italic {
		tekst = "<em>" + tekst + "</em>"
	}
	if postac.Underline != nil && *postac.Underline != shared.StudioUnderlineStyleNone {
		tekst = "<u>" + tekst + "</u>"
	}
	if postac.Strikethrough != nil && *postac.Strikethrough {
		tekst = "<s>" + tekst + "</s>"
	}
	if postac.Superscript != nil && *postac.Superscript {
		tekst = "<sup>" + tekst + "</sup>"
	}
	if postac.Subscript != nil && *postac.Subscript {
		tekst = "<sub>" + tekst + "</sub>"
	}
	styl := ""
	if postac.Color != nil && *postac.Color != "" {
		styl += "color: #" + ooxmlBezKrzyzyka(*postac.Color) + ";"
	}
	if postac.HighlightColor != nil && *postac.HighlightColor != "" {
		styl += "background-color: #" + ooxmlBezKrzyzyka(*postac.HighlightColor) + ";"
	}
	if postac.FontSizePt != nil {
		styl += "font-size: " + strconv.FormatFloat(*postac.FontSizePt, 'f', 1, 64) + "pt;"
	}
	if styl != "" {
		tekst = "<span style=\"" + styl + "\">" + tekst + "</span>"
	}
	return tekst
}

// wydanieTabelaHtml składa tabelę HTML wraz ze scaleniami i szerokościami kolumn.
//
// HTML niesie scalenia wprost (`colspan`, `rowspan`), więc tu NIE MA straty —
// i dlatego ta droga nie dokłada pozycji do wykazu pominiętych.
func wydanieTabelaHtml(postac *shared.StudioDocumentForm, kodTabeli *string) string {
	tabela := wydanieTabela(postac, kodTabeli)
	if tabela == nil {
		return ""
	}
	var budowa strings.Builder
	budowa.WriteString("<table>\n")
	if len(tabela.ColumnWidthsMm) > 0 {
		budowa.WriteString("<colgroup>")
		for _, szerokosc := range tabela.ColumnWidthsMm {
			budowa.WriteString("<col style=\"width: " +
				strconv.FormatFloat(szerokosc, 'f', 1, 64) + "mm\">")
		}
		budowa.WriteString("</colgroup>\n")
	}
	wierszyNaglowka := 0
	if tabela.HeaderRows != nil {
		wierszyNaglowka = *tabela.HeaderRows
	}
	for wiersz := 0; wiersz < tabela.Rows; wiersz++ {
		budowa.WriteString("<tr>")
		for kolumna := 0; kolumna < tabela.Columns; kolumna++ {
			komorka := ooxmlSzukajKomorki(tabela, wiersz, kolumna)
			if komorka == nil {
				budowa.WriteString("<td></td>")
				continue
			}
			if komorka.Merged != nil && *komorka.Merged {
				continue
			}
			znacznik := "td"
			if wiersz < wierszyNaglowka {
				znacznik = "th"
			}
			atrybuty := ""
			if komorka.ColumnSpan != nil && *komorka.ColumnSpan > 1 {
				atrybuty += " colspan=\"" + strconv.Itoa(*komorka.ColumnSpan) + "\""
			}
			if komorka.RowSpan != nil && *komorka.RowSpan > 1 {
				atrybuty += " rowspan=\"" + strconv.Itoa(*komorka.RowSpan) + "\""
			}
			budowa.WriteString("<" + znacznik + atrybuty + ">" +
				ooxmlZabezpiecz(strings.TrimSpace(wartoscTekstu(komorka.Text))) +
				"</" + znacznik + ">")
		}
		budowa.WriteString("</tr>\n")
	}
	budowa.WriteString("</table>\n")
	if tabela.Caption != nil && strings.TrimSpace(*tabela.Caption) != "" {
		budowa.WriteString("<p class=\"podpis-tabeli\">" +
			ooxmlZabezpiecz(strings.TrimSpace(*tabela.Caption)) + "</p>\n")
	}
	return budowa.String()
}

// ── Wydanie RTF ─────────────────────────────────────────────────────────────

// wydanieRtfem składa plik RTF rachunkiem własnym: postać znaku, akapit,
// wyrównanie i tabelę; obrazy osadzone i aparat odświeżalny idą do pominiętych.
func wydanieRtfem(postac *shared.StudioDocumentForm, tresc string) ([]byte,
	[]shared.StudioSkippedItem) {

	pominiete := []shared.StudioSkippedItem{}
	if len(postac.Objects) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "obrazy wyszły tekstem zastępczym — osadzenia bajtów obrazu w RTF ten " +
				"rachunek nie składa",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(postac.Objects)) + " obiektów"),
		})
	}
	if len(postac.Apparatus) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "aparat dokumentu wyszedł treścią, bez pól odświeżalnych",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(postac.Apparatus)) + " elementów"),
		})
	}
	if len(postac.Styles) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "arkusz stylów wyszedł postacią wpisaną przy akapitach, bez tablicy " +
				"stylów RTF",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(postac.Styles)) + " stylów; " +
				"zmiana stylu w edytorze Operatora nie przestawi już tego pliku"),
		})
	}

	var budowa strings.Builder
	budowa.WriteString(`{\rtf1\ansi\ansicpg1250\deff0` +
		`{\fonttbl{\f0\froman Times New Roman;}}` + "\n")
	for _, blok := range postac.Blocks {
		switch blok.Kind {
		case wejscieRodzajBlokuTabela:
			budowa.WriteString(wydanieTabelaRtf(postac, blok.TableId))
		case wejscieRodzajBlokuObiekt:
			if opis := wydanieOpisObiektu(postac, blok.ObjectId); opis != "" {
				budowa.WriteString(`\pard ` + wydanieTekstRtf("["+opis+"]") + `\par` + "\n")
			}
		case wejscieRodzajBlokuPodzial:
			budowa.WriteString(`\page` + "\n")
		default:
			budowa.WriteString(wydanieAkapitRtf(blok))
		}
	}
	if len(postac.Blocks) == 0 && strings.TrimSpace(tresc) != "" {
		for _, akapit := range strings.Split(tresc, "\n") {
			budowa.WriteString(`\pard ` + wydanieTekstRtf(akapit) + `\par` + "\n")
		}
	}
	budowa.WriteString("}\n")
	return []byte(budowa.String()), pominiete
}

// wydanieAkapitRtf składa akapit RTF wraz z wyrównaniem, wcięciami, odstępami
// akapitowymi oraz postacią znaku fragmentów, które akapit niesie.
func wydanieAkapitRtf(blok shared.StudioDocumentBlock) string {
	var budowa strings.Builder
	budowa.WriteString(`\pard`)
	if akapit := blok.Paragraph; akapit != nil {
		if akapit.Align != nil {
			switch *akapit.Align {
			case shared.StudioTextAlignCenter:
				budowa.WriteString(`\qc`)
			case shared.StudioTextAlignRight:
				budowa.WriteString(`\qr`)
			case shared.StudioTextAlignJustify:
				budowa.WriteString(`\qj`)
			default:
				budowa.WriteString(`\ql`)
			}
		}
		if akapit.IndentLeftMm != nil {
			budowa.WriteString(`\li` + strconv.Itoa(ooxmlMilimetryNaTwipy(*akapit.IndentLeftMm)))
		}
		if akapit.FirstLineIndentMm != nil {
			budowa.WriteString(`\fi` +
				strconv.Itoa(ooxmlMilimetryNaTwipy(*akapit.FirstLineIndentMm)))
		}
		if akapit.SpaceBeforePt != nil {
			budowa.WriteString(`\sb` + strconv.Itoa(int(*akapit.SpaceBeforePt*20)))
		}
		if akapit.SpaceAfterPt != nil {
			budowa.WriteString(`\sa` + strconv.Itoa(int(*akapit.SpaceAfterPt*20)))
		}
	}
	budowa.WriteString(" ")
	for _, fragment := range blok.Runs {
		budowa.WriteString(wydanieFragmentRtf(fragment))
	}
	budowa.WriteString(`\par` + "\n")
	return budowa.String()
}

// wydanieFragmentRtf składa fragment RTF wraz z postacią znaku, sterowaniami
// grubości, kursywy, podkreślenia, przekreślenia i stopnia pisma w półpunktach.
func wydanieFragmentRtf(fragment shared.StudioDocumentRun) string {
	if fragment.Text == "" {
		return ""
	}
	if fragment.Format == nil {
		return wydanieTekstRtf(fragment.Text)
	}
	var budowa strings.Builder
	budowa.WriteString("{")
	postac := fragment.Format
	if postac.Bold != nil && *postac.Bold {
		budowa.WriteString(`\b`)
	}
	if postac.Italic != nil && *postac.Italic {
		budowa.WriteString(`\i`)
	}
	if postac.Underline != nil && *postac.Underline != shared.StudioUnderlineStyleNone {
		budowa.WriteString(`\ul`)
	}
	if postac.Strikethrough != nil && *postac.Strikethrough {
		budowa.WriteString(`\strike`)
	}
	if postac.Superscript != nil && *postac.Superscript {
		budowa.WriteString(`\super`)
	}
	if postac.Subscript != nil && *postac.Subscript {
		budowa.WriteString(`\sub`)
	}
	if postac.FontSizePt != nil {
		// RTF liczy stopień pisma w PÓŁPUNKTACH — tak samo jak OOXML.
		budowa.WriteString(`\fs` + strconv.Itoa(int(*postac.FontSizePt*2)))
	}
	budowa.WriteString(" ")
	budowa.WriteString(wydanieTekstRtf(fragment.Text))
	budowa.WriteString("}")
	return budowa.String()
}

// wydanieTabelaRtf składa tabelę RTF wraz z granicami komórek liczonymi
// z szerokości kolumn postaci, w twipach, jednostce miary sterowań RTF.
func wydanieTabelaRtf(postac *shared.StudioDocumentForm, kodTabeli *string) string {
	tabela := wydanieTabela(postac, kodTabeli)
	if tabela == nil {
		return ""
	}
	szerokosci := tabela.ColumnWidthsMm
	if len(szerokosci) != tabela.Columns {
		// Bez szerokości z postaci kolumny dzielą obszar pisania równo.
		szerokosci = make([]float64, tabela.Columns)
		for i := range szerokosci {
			szerokosci[i] = 160.0 / float64(wydanieWieksza(tabela.Columns, 1))
		}
	}
	var budowa strings.Builder
	for wiersz := 0; wiersz < tabela.Rows; wiersz++ {
		budowa.WriteString(`\trowd`)
		granica := 0.0
		for _, szerokosc := range szerokosci {
			granica += szerokosc
			budowa.WriteString(`\cellx` + strconv.Itoa(ooxmlMilimetryNaTwipy(granica)))
		}
		for kolumna := 0; kolumna < tabela.Columns; kolumna++ {
			komorka := ooxmlSzukajKomorki(tabela, wiersz, kolumna)
			tekst := ""
			if komorka != nil {
				tekst = strings.ReplaceAll(
					strings.TrimSpace(wartoscTekstu(komorka.Text)), "\n", " ")
			}
			budowa.WriteString(`\pard\intbl ` + wydanieTekstRtf(tekst) + `\cell`)
		}
		budowa.WriteString(`\row` + "\n")
	}
	return budowa.String()
}

// wydanieWieksza oddaje większą z dwóch liczb całkowitych, bez sięgania po
// funkcję ogólną biblioteki wzorcowej dla tego jednego porównania.
func wydanieWieksza(pierwsza, druga int) int {
	if pierwsza > druga {
		return pierwsza
	}
	return druga
}

// wydanieTekstRtf zapisuje tekst zapisem RTF.
//
// Znaki poza ASCII idą jako `\uNNNN` wraz ze znakiem zastępczym: tak stanowi
// gramatyka RTF i tylko tak pismo polskie otworzy się poprawnie w czytniku,
// który kroju uniwersalnego nie ma.
func wydanieTekstRtf(tekst string) string {
	var budowa strings.Builder
	for _, znak := range tekst {
		switch {
		case znak == '\\' || znak == '{' || znak == '}':
			budowa.WriteString("\\" + string(znak))
		case znak == '\n':
			budowa.WriteString(`\line `)
		case znak == '\t':
			budowa.WriteString(`\tab `)
		case znak < 128:
			budowa.WriteRune(znak)
		default:
			budowa.WriteString(`\u` + strconv.Itoa(int(znak)) + "?")
		}
	}
	return budowa.String()
}

// ── Wydanie PDF ─────────────────────────────────────────────────────────────

// wydaniePdfem składa dokument PDF biblioteką `pdfcpu`, przez profil wydania,
// który niesie nastawy paginacji i stopki; profil niewskazany nie jest odmową.
func (a *adapterStudia) wydaniePdfem(ctx context.Context, stan *stanPostaci,
	kodProfilu *string) ([]byte, []shared.StudioSkippedItem, error) {

	pominiete := []shared.StudioSkippedItem{{
		Reason: "postać znaku nie weszła w pełni — rachunek PDF składa dokument krojem " +
			"jednym, bez krojów dokumentu",
		Detail: wejscieWskaznikTekstu("naprawa dla pisma wzorcowego: wydać do docx albo odt, " +
			"gdzie postać przechodzi w całości"),
	}}
	if len(wydanieTabeleDokumentu(&stan.forma)) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "tabele wyszły wierszami tekstu z kolumnami rozdzielonymi odstępem",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(stan.forma.Tables)) +
				" tabel; siatka i scalenia nie weszły"),
		})
	}
	if len(stan.forma.Objects) > 0 {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "obrazy wyszły tekstem zastępczym",
			Detail: wejscieWskaznikTekstu(strconv.Itoa(len(stan.forma.Objects)) +
				" obiektów"),
		})
	}

	tekst, _ := wydanieTekstem(&stan.forma, wartoscTekstu(stan.dokument.Tresc))
	tresc := string(tekst)

	stopka := ""
	if kod := strings.TrimSpace(wartoscTekstu(kodProfilu)); kod != "" {
		profil, err := a.repozytorium.ProfilWydania(ctx, kod)
		if err != nil {
			if wejscieBrakWiersza(err) {
				return nil, pominiete, wejscieBladBraku("profil wydania nie istnieje: " + kod)
			}
			return nil, pominiete, bladStudio(err)
		}
		stopka = profil.Nazwa
	} else {
		pominiete = append(pominiete, shared.StudioSkippedItem{
			Reason: "wydanie poszło bez profilu — paginacja i stopka są domyślne",
			Detail: wejscieWskaznikTekstu("naprawa: wskazać profil wydania polem profileId"),
		})
	}
	if stan.forma.PageSetup != nil && stan.forma.PageSetup.Footer != nil {
		stopka = strings.TrimSpace(*stan.forma.PageSetup.Footer)
	}
	if stopka != "" {
		tresc += "\n\n" + stopka
	}

	bajty, err := dokumentPdfZTekstu(tresc)
	if err != nil {
		return nil, pominiete, err
	}
	return bajty, pominiete, nil
}

// wydanieTabeleDokumentu oddaje tabele dokumentu — osobna nazwa, żeby warunek
// czytał się jak zdanie, a nie jak rachunek długości.
func wydanieTabeleDokumentu(forma *shared.StudioDocumentForm) []shared.StudioDocumentTable {
	return forma.Tables
}
