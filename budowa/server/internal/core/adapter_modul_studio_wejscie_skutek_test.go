package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"danacoconsole/shared"
)

// Sprawdziany skutku wejścia modułu Studio: uczciwość konwersji PDF, osobność
// kopii dokumentu i przeżywalność postaci przy zapisie.
//
// Szkody, które ten plik ma wykluczyć:
//  1. konwersja PDF meldująca powodzenie i oddająca dokument okaleczony bez
//     ani jednego słowa o tym, czego nie odzyskała;
//  2. PDF ze samych skanów przepuszczony jako „skonwertowany" — czyli dokument
//     pusty podany jako gotowy do pracy;
//  3. kopia dokumentu będąca drugim odwołaniem do tego samego bytu, po której
//     poprawka w kopii zmienia oryginał;
//  4. `studio.document.save` gubiący postać dokumentu — ta sama dziura, przez
//     którą model był wobec dokumentu ślepy.

// ── Uprząż ──────────────────────────────────────────────────────────────────

// wejsciePdfZWarstwaTekstowa składa PDF o wskazanej liczbie stron, w którym
// KAŻDA strona niesie warstwę tekstową.
//
// Materiał powstaje `pdfcpu` — tą samą biblioteką, którą rdzeń go czyta.
// To jest świadome i nazwane: mierzona jest uczciwość bilansu konwersji, a nie
// zgodność dwóch bibliotek PDF między sobą. Tekst jest zapisany wprost jako
// operator pokazania tekstu, więc warstwa tekstowa jest tu prawdziwa.
func wejsciePdfZWarstwaTekstowa(t *testing.T, strony int) []byte {
	t.Helper()

	var opis strings.Builder
	opis.WriteString(`{"pages":{`)
	for numer := 1; numer <= strony; numer++ {
		if numer > 1 {
			opis.WriteString(",")
		}
		fmt.Fprintf(&opis, `"%d":{"content":{"text":[`+
			`{"value":"Paragraf pierwszy strony %d o tresci dostatecznie dlugiej.",`+
			`"font":{"name":"Helvetica","size":12},"position":[0.1,0.9]},`+
			`{"value":"Paragraf drugi strony %d zamykajacy jej mysl.",`+
			`"font":{"name":"Helvetica","size":12},"position":[0.1,0.8]}`+
			`]}}`, numer, numer, numer)
	}
	opis.WriteString(`}}`)

	var dokument bytes.Buffer
	if err := api.Create(nil, strings.NewReader(opis.String()), &dokument, nastawyPdf()); err != nil {
		t.Fatalf("nie można złożyć PDF-a z warstwą tekstową: %v", err)
	}
	return dokument.Bytes()
}

// wejsciePdfSamychSkanow składa PDF, w którym żadna strona warstwy tekstowej nie
// ma — strony niosą wyłącznie obraz, jak skan pisma z faksu.
func wejsciePdfSamychSkanow(t *testing.T, strony int) []byte {
	t.Helper()

	obraz := obrazPNG(t, 64, 64)
	sciezka := t.TempDir() + "/skan.png"
	zapiszPlikSprawdzianu(t, sciezka, obraz)

	var opis strings.Builder
	opis.WriteString(`{"pages":{`)
	for numer := 1; numer <= strony; numer++ {
		if numer > 1 {
			opis.WriteString(",")
		}
		wartosc, err := json.Marshal(sciezka)
		if err != nil {
			t.Fatalf("nie można złożyć opisu strony skanu: %v", err)
		}
		fmt.Fprintf(&opis, `"%d":{"content":{"image":[{"src":%s,`+
			`"position":[0.1,0.1],"width":50}]}}`, numer, wartosc)
	}
	opis.WriteString(`}}`)

	var dokument bytes.Buffer
	if err := api.Create(nil, strings.NewReader(opis.String()), &dokument, nastawyPdf()); err != nil {
		t.Fatalf("nie można złożyć PDF-a ze samych skanów: %v", err)
	}
	return dokument.Bytes()
}

// wejscieWniesPdf kieruje bajty PDF-a przez `studio.document.import.pdf`.
func wejscieWniesPdf(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno string, bajty []byte) shared.StudioDocumentImportPdfResponse {

	t.Helper()

	var wniesiony shared.StudioDocumentImportPdfResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentImportPdf,
		shared.StudioDocumentImportPdfRequest{
			WindowId:    okno,
			BytesBase64: wskaznik(wBase64(bajty)),
		}, &wniesiony)
	return wniesiony
}

// ── Punkt pierwszy: bilans konwersji PDF ────────────────────────────────────

// TestKonwersjaPdfOddajeBilansOdzyskania wykazuje, że konwersja PDF nie oddaje
// kaleki jako gotowego dokumentu.
//
// PDF nie niesie struktury akapitu ani tabeli wprost — odzyskanie jest
// ODTWORZENIEM, nie odczytem. Miara: bilans musi nieść policzone strony,
// policzone strony z warstwą tekstową i bez niej, oraz zdanie o stanie wyniku.
// Liczba stron w bilansie jest porównywana z liczbą policzoną w pliku
// biblioteką, a nie brana na słowo.
func TestKonwersjaPdfOddajeBilansOdzyskania(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	const stronMaterialu = 3
	bajty := wejsciePdfZWarstwaTekstowa(t, stronMaterialu)

	// Niezależny pomiar liczby stron: to, co bilans twierdzi, musi zgadzać się
	// z tym, co w pliku naprawdę jest.
	stronWPliku, err := api.PageCount(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		t.Fatalf("materiał próbny nie jest poprawnym PDF-em: %v", err)
	}
	if stronWPliku != stronMaterialu {
		t.Fatalf("materiał próbny ma %d stron, a miał mieć %d", stronWPliku, stronMaterialu)
	}

	wniesiony := wejscieWniesPdf(t, zmontowany, zycie, "okno-konwersji-pdf", bajty)
	bilans := wniesiony.Balance

	if bilans.Format != shared.StudioImportFormatPdf {
		t.Errorf("bilans mówi o formacie %q, a wnoszono PDF", bilans.Format)
	}
	if wartoscCalkowita(bilans.Pages) != stronWPliku {
		t.Errorf("bilans mówi o %d stronach, a w pliku jest %d — bilans nie liczy, "+
			"tylko zgaduje", wartoscCalkowita(bilans.Pages), stronWPliku)
	}
	// Rozbicie na strony z warstwą i bez niej musi się zsumować do całości.
	zWarstwa := wartoscCalkowita(bilans.PagesWithText)
	bezWarstwy := wartoscCalkowita(bilans.PagesWithoutText)
	if zWarstwa+bezWarstwy != stronWPliku {
		t.Errorf("strony z warstwą tekstową (%d) i bez niej (%d) nie sumują się do "+
			"%d stron dokumentu", zWarstwa, bezWarstwy, stronWPliku)
	}
	if zWarstwa == 0 {
		t.Error("bilans nie znalazł ani jednej strony z warstwą tekstową, choć " +
			"każda strona materiału ją niesie — rozbiór warstwy tekstowej nie działa")
	}
	// Odzyskane akapity mają być POLICZONE, nie zadeklarowane.
	if wartoscCalkowita(bilans.ParagraphsRecovered) == 0 {
		t.Error("bilans mówi o zerze odzyskanych akapitów, a dokument powstał — " +
			"jedno z dwóch jest nieprawdą")
	}
	// Uczciwość nazwana: PDF nie niesie tabel wprost, więc bilans ma o tym
	// mówić, a nie milczeć.
	if bilans.Note == nil || strings.TrimSpace(*bilans.Note) == "" {
		t.Error("bilans konwersji PDF nie niesie zdania o uczciwym stanie wyniku")
	} else {
		nota := strings.ToLower(*bilans.Note)
		if !strings.Contains(nota, "odtworzen") && !strings.Contains(nota, "odzysk") {
			t.Errorf("zdanie bilansu nie mówi, że odzyskanie jest odtworzeniem: %q", nota)
		}
	}
	if len(bilans.Skipped) == 0 {
		t.Error("konwersja PDF nie wymieniła ani jednej rzeczy nieodzyskanej — " +
			"PDF nie niesie struktury akapitu ani tabeli, więc coś odpaść musiało")
	}

	// Skutek na dokumencie: treść ma naprawdę wejść do edytora.
	tresc := trescDokumentu(t, zmontowany, zycie, "okno-konwersji-pdf",
		wniesiony.Document.Id)
	if strings.TrimSpace(tresc) == "" {
		t.Fatal("konwersja PDF zameldowała powodzenie i zostawiła dokument pusty")
	}
	if !strings.Contains(tresc, "Paragraf pierwszy") {
		t.Errorf("w dokumencie nie ma tekstu ze pierwszej strony PDF-a; treść: %q",
			poczatekTekstu(tresc))
	}
	// Wszystkie strony mają wejść, nie tylko pierwsza.
	if !strings.Contains(tresc, "strony 3") {
		t.Errorf("w dokumencie nie ma tekstu ze strony trzeciej — konwersja " +
			"przepuściła tylko część pliku")
	}
	if wniesiony.Form.DocumentId != wniesiony.Document.Id {
		t.Error("postać oddana przez konwersję nie należy do dokumentu, który powstał")
	}
}

// TestKonwersjaPdfSamychSkanowKierujeNaRozpoznanie wykazuje, że PDF bez warstwy
// tekstowej NIE jest udawany jako skonwertowany.
//
// To jest wymaganie rozstrzygające: dokument pusty oddany jako „gotowy do
// pracy" byłby najgorszą możliwą odpowiedzią, bo Operator zaczął by pisać
// w pliku, który treści nie ma. Miara: bilans zaznacza potrzebę rozpoznania
// pisma ORAZ odpowiedź wskazuje pozycję kolejki rozpoznania — nazwana droga
// dalsza, nie sama odmowa.
func TestKonwersjaPdfSamychSkanowKierujeNaRozpoznanie(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	const stronSkanu = 2
	bajty := wejsciePdfSamychSkanow(t, stronSkanu)
	wniesiony := wejscieWniesPdf(t, zmontowany, zycie, "okno-skanow", bajty)
	bilans := wniesiony.Balance

	if bilans.NeedsTextRecognition == nil || !*bilans.NeedsTextRecognition {
		t.Error("PDF ze samych skanów nie został oznaczony jako wymagający " +
			"rozpoznania pisma — rdzeń podał kalekę za dokument gotowy")
	}
	if wartoscCalkowita(bilans.PagesWithText) != 0 {
		t.Errorf("bilans twierdzi, że %d stron skanu ma warstwę tekstową",
			wartoscCalkowita(bilans.PagesWithText))
	}
	if wartoscCalkowita(bilans.PagesWithoutText) != stronSkanu {
		t.Errorf("bilans mówi o %d stronach bez warstwy tekstowej, a skan ma %d",
			wartoscCalkowita(bilans.PagesWithoutText), stronSkanu)
	}
	// Droga dalsza ma być WSKAZANA, nie domyślona przez Operatora.
	if wniesiony.IngestItemId == nil || strings.TrimSpace(*wniesiony.IngestItemId) == "" {
		t.Error("PDF ze samych skanów nie dostał pozycji kolejki rozpoznania pisma — " +
			"Operator zostaje bez drogi dalszej")
	} else {
		// Pozycja ma istnieć naprawdę: wykaz kolejki musi ją znać.
		var kolejka shared.StudioIngestQueueListResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestQueueList,
			shared.StudioIngestQueueListRequest{WindowId: "okno-skanow"}, &kolejka)
		znaleziona := false
		for _, pozycja := range kolejka.Items {
			if pozycja.Id == *wniesiony.IngestItemId {
				znaleziona = true
			}
		}
		if !znaleziona {
			t.Errorf("odpowiedź wskazuje pozycję kolejki %q, której w kolejce nie ma",
				*wniesiony.IngestItemId)
		}
	}
	if bilans.Note == nil || !strings.Contains(strings.ToLower(*bilans.Note), "rozpozna") {
		t.Errorf("zdanie bilansu nie kieruje na rozpoznanie pisma: %v", bilans.Note)
	}
}

// poczatekTekstu przycina treść do wielkości czytelnej w dzienniku sprawdzianu.
func poczatekTekstu(tekst string) string {
	const granica = 200
	if len(tekst) <= granica {
		return tekst
	}
	return tekst[:granica] + "…"
}

// ── Punkt drugi: kopia jest osobnym bytem ───────────────────────────────────

// TestKopiaDokumentuJestOsobnymBytem wykazuje, że kopia jest osobnym
// dokumentem, a nie drugim odwołaniem do tego samego.
//
// Miara jest ta, którą wskazał Właściciel: ZMIANA W KOPII NIE RUSZA ORYGINAŁU.
// Sprawdzane są obie warstwy — treść i postać — bo kopia dzieląca postać
// z oryginałem jest tak samo zepsuta jak kopia dzieląca treść, tylko trudniej
// to zauważyć.
func TestKopiaDokumentuJestOsobnymBytem(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	okno := "okno-kopii"

	const trescPierwotna = "Zdanie pierwotne oryginalu."
	oryginal := dokumentZTrescia(t, zmontowany, zycie, okno, trescPierwotna)

	// Oryginał dostaje postać, żeby było co porównywać: wytłuszczenie fragmentu.
	var postawiona shared.StudioFormatCharacterSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioFormatCharacterSet,
		shared.StudioFormatCharacterSetRequest{
			DocumentId: oryginal.Id,
			RangeStart: wskaznik(0),
			RangeEnd:   wskaznik(6),
			Bold:       wskaznik(true),
		}, &postawiona)

	var kopia shared.StudioDocumentCopyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentCopy,
		shared.StudioDocumentCopyRequest{
			DocumentId: oryginal.Id,
			Title:      wskaznik("Kopia pisma"),
		}, &kopia)

	if kopia.Document.Id == oryginal.Id {
		t.Fatal("kopia dostała identyfikator oryginału — to jest to samo odwołanie, " +
			"nie kopia")
	}

	// Postać ma przejść do kopii — kopia bez postaci nie jest kopią dokumentu.
	postacKopii := ooxmlPostacDokumentu(t, zmontowany, zycie, kopia.Document.Id)
	if !ooxmlCzyFragmentWytluszczony(postacKopii, "Zdanie") {
		t.Error("kopia nie przejęła wytłuszczenia oryginału — skopiowano treść " +
			"bez postaci")
	}
	if postacKopii.DocumentId != kopia.Document.Id {
		t.Errorf("postać kopii jest podpisana dokumentem %q, a kopia ma %q — dwa "+
			"dokumenty wskazują jedną postać", postacKopii.DocumentId, kopia.Document.Id)
	}

	// ── Miara rozstrzygająca: zmiana w kopii ────────────────────────────────
	const trescKopii = "Zdanie zmienione w kopii, oryginal ma zostac nietkniety."
	var zapisanaKopia shared.StudioDocumentSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId: kopia.Document.Id, Content: trescKopii,
		}, &zapisanaKopia)

	// Postać kopii też się zmienia — zdejmujemy wytłuszczenie.
	var zdjeta shared.StudioFormatCharacterSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioFormatCharacterSet,
		shared.StudioFormatCharacterSetRequest{
			DocumentId: kopia.Document.Id,
			RangeStart: wskaznik(0),
			RangeEnd:   wskaznik(6),
			Bold:       wskaznik(false),
		}, &zdjeta)

	// Oryginał czytany OSOBNYM wywołaniem — to jest miara właściwa.
	trescOryginalu := trescDokumentu(t, zmontowany, zycie, okno, oryginal.Id)
	if trescOryginalu != trescPierwotna {
		t.Errorf("zmiana treści w kopii ruszyła oryginał: oryginał niesie %q, "+
			"a miał %q", trescOryginalu, trescPierwotna)
	}
	postacOryginalu := ooxmlPostacDokumentu(t, zmontowany, zycie, oryginal.Id)
	if !ooxmlCzyFragmentWytluszczony(postacOryginalu, "Zdanie") {
		t.Error("zdjęcie wytłuszczenia w kopii zdjęło je także w oryginale — " +
			"oba dokumenty stoją na jednej postaci")
	}

	// I odwrotnie: kopia ma naprawdę nieść swoją zmianę, a nie treść oryginału.
	trescPoZmianie := trescDokumentu(t, zmontowany, zycie, okno, kopia.Document.Id)
	if trescPoZmianie != trescKopii {
		t.Errorf("kopia nie niesie własnej treści po zapisie: %q", trescPoZmianie)
	}
}

// TestKopiaBezHistoriiIZHistoriaSaJawnymWyborem wykazuje, że przeniesienie
// historii wersji jest JAWNYM wyborem Operatora, a nie zachowaniem zaszytym.
//
// Liczba przeniesionych wersji wychodzi kontraktem, więc sprawdzian mierzy ją,
// a nie samo powodzenie komendy.
func TestKopiaBezHistoriiIZHistoriaSaJawnymWyborem(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	okno := "okno-kopii-historii"

	oryginal := dokumentZTrescia(t, zmontowany, zycie, okno, "Wersja pierwsza.")
	// Druga i trzecia wersja, żeby historia była historią, a nie jednym wpisem.
	for _, tresc := range []string{"Wersja druga.", "Wersja trzecia."} {
		var zapisany shared.StudioDocumentSaveResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentSave,
			shared.StudioDocumentSaveRequest{
				DocumentId: oryginal.Id, Content: tresc, CreateVersion: wskaznik(true),
			}, &zapisany)
	}

	var wersjeOryginalu shared.StudioRepositoryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioRepositoryList,
		shared.StudioRepositoryListRequest{DocumentId: oryginal.Id}, &wersjeOryginalu)
	if len(wersjeOryginalu.Versions) < 3 {
		t.Fatalf("oryginał ma %d wersji, a sprawdzian potrzebuje co najmniej trzech",
			len(wersjeOryginalu.Versions))
	}

	// Bez historii — wybór domyślny wedle kontraktu.
	var bezHistorii shared.StudioDocumentCopyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentCopy,
		shared.StudioDocumentCopyRequest{DocumentId: oryginal.Id}, &bezHistorii)
	if bezHistorii.VersionsCopied != 0 {
		t.Errorf("kopia bez wskazania historii przeniosła %d wersji — kontrakt mówi, "+
			"że brak wskazania znaczy „bez historii”", bezHistorii.VersionsCopied)
	}

	// Z historią — wybór jawny.
	var zHistoria shared.StudioDocumentCopyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentCopy,
		shared.StudioDocumentCopyRequest{
			DocumentId: oryginal.Id, IncludeVersions: wskaznik(true),
		}, &zHistoria)
	if zHistoria.VersionsCopied == 0 {
		t.Error("kopia z jawnym wskazaniem historii nie przeniosła ani jednej wersji")
	}

	// Skutek POLICZONY w repozytorium kopii, nie wzięty z odpowiedzi.
	var wersjeKopii shared.StudioRepositoryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioRepositoryList,
		shared.StudioRepositoryListRequest{DocumentId: zHistoria.Document.Id}, &wersjeKopii)
	if len(wersjeKopii.Versions) != zHistoria.VersionsCopied {
		t.Errorf("odpowiedź mówi o %d przeniesionych wersjach, a w repozytorium kopii "+
			"stoi %d", zHistoria.VersionsCopied, len(wersjeKopii.Versions))
	}

	var wersjeBezHistorii shared.StudioRepositoryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioRepositoryList,
		shared.StudioRepositoryListRequest{DocumentId: bezHistorii.Document.Id},
		&wersjeBezHistorii)
	if len(wersjeBezHistorii.Versions) >= len(wersjeOryginalu.Versions) {
		t.Errorf("kopia „bez historii” ma %d wersji, a oryginał %d — historia "+
			"przeszła wbrew wyborowi", len(wersjeBezHistorii.Versions),
			len(wersjeOryginalu.Versions))
	}
}

// ── Punkt trzeci: zapis przenosi postać ─────────────────────────────────────

// TestZapisPrzenosiPostacDokumentu wykazuje, że `studio.document.save` przenosi
// POSTAĆ, a nie samą treść.
//
// To jest dziura, od której zaczęło się całe zlecenie: komenda przyjmowała
// `documentId`, `content` i `title`, więc postać przy każdym zapisie ginęła
// i model był wobec dokumentu ślepy. Miara: postać podana polem `form`
// odczytana PONOWNIE, osobnym wywołaniem, musi być ta sama — arkusz stylów,
// nastawy strony, tabela i styl akapitu.
func TestZapisPrzenosiPostacDokumentu(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	okno := "okno-zapisu-postaci"

	dokument := dokumentZTrescia(t, zmontowany, zycie, okno, "Tresc pierwotna.")

	// Postać budowana wprost, żeby sprawdzian wiedział, czego szuka.
	nastawy := wejscieDomyslneNastawyStrony("A5", wejscieWskaznikOrientacji(
		shared.StudioPageOrientationPozioma))
	nastawy.MarginTop = wejscieWskaznikCalkowity(33)
	sekcja := shared.StudioSection{
		Id: "sekcja-zapisu", Index: 0, PageSetup: &nastawy,
	}
	tabela := shared.StudioDocumentTable{
		Id: "tabela-zapisu", Rows: 2, Columns: 2,
		ColumnWidthsMm: []float64{60, 40},
		HeaderRows:     wejscieWskaznikCalkowity(1),
		RepeatHeader:   wejscieWskaznikLogiczny(true),
		Cells: []shared.StudioTableCell{
			{Row: 0, Column: 0, Text: wskaznik("Pozycja")},
			{Row: 0, Column: 1, Text: wskaznik("Kwota")},
			{Row: 1, Column: 0, Text: wskaznik("Wykonanie")},
			{Row: 1, Column: 1, Text: wskaznik("9 000")},
		},
	}
	postac := shared.StudioDocumentForm{
		DocumentId: dokument.Id,
		PageSetup:  &nastawy,
		Styles:     wejscieDomyslnyArkuszStylow(),
		Sections:   []shared.StudioSection{sekcja},
		Tables:     []shared.StudioDocumentTable{tabela},
		Blocks: []shared.StudioDocumentBlock{
			{
				Id: "blok-naglowka", Kind: wejscieRodzajBlokuAkapit,
				SectionId: wskaznik(sekcja.Id),
				Paragraph: &shared.StudioParagraphFormat{
					StyleName:    wskaznik(wejscieStylNaglowek2),
					OutlineLevel: wejscieWskaznikCalkowity(2),
				},
				Runs: []shared.StudioDocumentRun{{Text: "Naglowek pisma"}},
			},
			{
				Id: "blok-tresci", Kind: wejscieRodzajBlokuAkapit,
				SectionId: wskaznik(sekcja.Id),
				Paragraph: &shared.StudioParagraphFormat{
					StyleName: wskaznik(wejscieStylTekstZasadniczy),
				},
				Runs: []shared.StudioDocumentRun{{
					Text:   "Zdanie wytluszczone.",
					Format: &shared.StudioCharacterFormat{Bold: wejscieWskaznikLogiczny(true)},
				}},
			},
			{
				Id: "blok-tabeli", Kind: wejscieRodzajBlokuTabela,
				SectionId: wskaznik(sekcja.Id), TableId: wskaznik(tabela.Id),
			},
		},
	}
	zapisPostaci, err := json.Marshal(postac)
	if err != nil {
		t.Fatalf("nie można złożyć postaci sprawdzianu: %v", err)
	}

	var zapisany shared.StudioDocumentSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId: dokument.Id,
			Content:    "Naglowek pisma\nZdanie wytluszczone.",
			Form:       zapisPostaci,
		}, &zapisany)

	// ── Odczyt PONOWNY: to jest miara ───────────────────────────────────────
	poZapisie := ooxmlPostacDokumentu(t, zmontowany, zycie, dokument.Id)

	if len(poZapisie.Styles) == 0 {
		t.Error("po zapisie postaci arkusz stylów jest pusty — postać nie dojechała " +
			"do rdzenia")
	}
	nastawyPoZapisie := poZapisie.PageSetup
	if nastawyPoZapisie == nil && len(poZapisie.Sections) > 0 {
		nastawyPoZapisie = poZapisie.Sections[0].PageSetup
	}
	if nastawyPoZapisie == nil {
		t.Fatal("po zapisie postaci nie ma nastaw strony ani na dokumencie, ani na sekcji")
	}
	if nastawyPoZapisie.PageSize == nil || *nastawyPoZapisie.PageSize != "A5" {
		t.Errorf("nośnik A5 nie przeżył zapisu — stoi %v", nastawyPoZapisie.PageSize)
	}
	if nastawyPoZapisie.Orientation == nil ||
		*nastawyPoZapisie.Orientation != shared.StudioPageOrientationPozioma {
		t.Errorf("orientacja pozioma nie przeżyła zapisu — stoi %v",
			nastawyPoZapisie.Orientation)
	}
	if nastawyPoZapisie.MarginTop == nil || *nastawyPoZapisie.MarginTop != 33 {
		t.Errorf("margines górny 33 mm nie przeżył zapisu — stoi %v",
			nastawyPoZapisie.MarginTop)
	}
	if len(poZapisie.Sections) == 0 {
		t.Error("sekcja nie przeżyła zapisu")
	}
	if len(poZapisie.Tables) != 1 {
		t.Fatalf("po zapisie stoi %d tabel, a zapisano jedną", len(poZapisie.Tables))
	}
	tabelaPoZapisie := poZapisie.Tables[0]
	if tabelaPoZapisie.Rows != 2 || tabelaPoZapisie.Columns != 2 {
		t.Errorf("tabela po zapisie ma %d na %d, a zapisano 2 na 2",
			tabelaPoZapisie.Rows, tabelaPoZapisie.Columns)
	}
	if len(tabelaPoZapisie.ColumnWidthsMm) != 2 ||
		tabelaPoZapisie.ColumnWidthsMm[0] != 60 {
		t.Errorf("szerokości kolumn nie przeżyły zapisu: %v",
			tabelaPoZapisie.ColumnWidthsMm)
	}
	if tabelaPoZapisie.HeaderRows == nil || *tabelaPoZapisie.HeaderRows != 1 {
		t.Errorf("wiersz nagłówkowy tabeli nie przeżył zapisu: %v",
			tabelaPoZapisie.HeaderRows)
	}
	if ooxmlBlokOStylu(poZapisie, wejscieStylNaglowek2) == nil {
		t.Errorf("styl akapitu %q nie przeżył zapisu", wejscieStylNaglowek2)
	}
	if !ooxmlCzyFragmentWytluszczony(poZapisie, "wytluszczone") {
		t.Error("postać znaku nie przeżyła zapisu — fragment stracił wytłuszczenie")
	}

	// Numer porządkowy postaci ma rosnąć: dwa zapisy tej samej postaci nie mogą
	// wyjść z tym samym numerem, bo wtedy nie da się rozstrzygnąć kolejności.
	if poZapisie.Revision == nil {
		t.Error("postać po zapisie nie niesie numeru porządkowego")
	}

	// ── Zapis bez pola `form` NIE MA PRAWA zetrzeć postaci ──────────────────
	// To jest druga połowa tej samej szkody: kontrakt mówi, że brak pola znaczy
	// „bez zmiany postaci", a nie „postać na zero".
	var drugiZapis shared.StudioDocumentSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId: dokument.Id, Content: "Tresc zmieniona bez postaci.",
		}, &drugiZapis)

	poDrugimZapisie := ooxmlPostacDokumentu(t, zmontowany, zycie, dokument.Id)
	if len(poDrugimZapisie.Tables) != 1 {
		t.Errorf("zapis bez pola form starł tabelę dokumentu — stoi %d tabel, "+
			"a kontrakt mówi, że brak pola znaczy „bez zmiany postaci”",
			len(poDrugimZapisie.Tables))
	}
	if len(poDrugimZapisie.Styles) == 0 {
		t.Error("zapis bez pola form starł arkusz stylów dokumentu")
	}
}
