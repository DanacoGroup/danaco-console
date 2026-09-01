package core

import (
	"database/sql"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// poczatekTekstuBadania przycina treść do wielkości czytelnej w komunikacie
// niepowodzenia — plik eksportu bywa dłuższy niż dziennik sprawdzianu.
func poczatekTekstuBadania(tekst string, granica int) string {
	if len(tekst) <= granica {
		return tekst
	}
	return tekst[:granica]
}

// bazaBadania otwiera plik bazy sprawdzianu do pomiaru niezależnego, drugim
// połączeniem obok tego, którym pracuje rdzeń.
func bazaBadania(t *testing.T, katalogDanych string) *sql.DB {
	t.Helper()

	baza, err := sql.Open("sqlite", filepath.Join(katalogDanych, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy do pomiaru: %v", err)
	}
	t.Cleanup(func() { _ = baza.Close() })
	return baza
}

// policzBadania oddaje pojedynczą liczbę zwróconą przez zapytanie pomiarowe,
// żeby sprawdziany nie powtarzały tego samego odczytu.
func policzBadania(t *testing.T, baza *sql.DB, zapytanie string, argumenty ...any) int {
	t.Helper()

	var liczba int
	if err := baza.QueryRow(zapytanie, argumenty...).Scan(&liczba); err != nil {
		t.Fatalf("pomiar %q nie powiódł się: %v", zapytanie, err)
	}
	return liczba
}

// TestSkutekKatalogowaniaZrodelBadania mierzy, czy katalogowanie źródeł zostawia
// wiersze: samo źródło, jego etykiety oraz załącznik z bajtami na dysku.
func TestSkutekKatalogowaniaZrodelBadania(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaBadania(t, katalogDanych)
	okno := "okno-badania-katalog"

	var dodane shared.ResearchSourceAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchSourceAdd,
		shared.ResearchSourceAddRequest{
			WindowId: okno, Title: "Raport rynkowy 2024",
			Kind: shared.ResearchSourceKindDocument,
		}, &dodane)

	if liczba := policzBadania(t, baza,
		`SELECT COUNT(*) FROM zrodlo_badania WHERE identyfikator_zewnetrzny = ?`,
		dodane.Source.Id); liczba != 1 {
		t.Fatalf("źródło zameldowane, a w bazie go nie ma (wierszy: %d)", liczba)
	}

	// Etykiety: skutkiem jest wiązanie w tabeli katalogu, nie pole odpowiedzi.
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchSourceTag,
		shared.ResearchSourceTagRequest{
			SourceId: dodane.Source.Id, Tags: []string{"rynek", "2024"},
		}, nil)
	if liczba := policzBadania(t, baza, `SELECT COUNT(*) FROM etykieta_zrodla_badania e
	    JOIN zrodlo_badania z ON z.id = e.zrodlo_id
	    WHERE z.identyfikator_zewnetrzny = ?`, dodane.Source.Id); liczba != 2 {
		t.Fatalf("etykiety zameldowane, a w bazie leży ich %d zamiast 2", liczba)
	}

	// Załącznik: skutkiem są BAJTY w magazynie, a nie wiersz o pliku.
	tresc := []byte("pełny tekst raportu rynkowego")
	zakodowana := base64.StdEncoding.EncodeToString(tresc)
	var zalacznik shared.ResearchSourceAttachmentAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchSourceAttachmentAdd,
		shared.ResearchSourceAttachmentAddRequest{
			SourceId: dodane.Source.Id, Kind: shared.ResearchAttachmentKindFulltext,
			ContentBase64: &zakodowana,
		}, &zalacznik)

	var sciezka string
	err := baza.QueryRow(`SELECT za.sciezka FROM zalacznik_zrodla_badania za
	    WHERE za.identyfikator_zewnetrzny = ?`, zalacznik.Attachment.Id).Scan(&sciezka)
	if err != nil {
		t.Fatalf("załącznik zameldowany, a w bazie nie ma jego wiersza: %v", err)
	}
	leżące, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("wiersz załącznika wskazuje plik, którego nie ma: %v", err)
	}
	if string(leżące) != string(tresc) {
		t.Fatalf("w magazynie leży inna treść niż wgrana: %q", string(leżące))
	}

	// Wykaz braków kompletności ma mówić prawdę o tym, czego brak.
	var wykaz shared.ResearchSourceAttachmentListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchSourceAttachmentList,
		shared.ResearchSourceAttachmentListRequest{SourceId: dodane.Source.Id}, &wykaz)
	if len(wykaz.Attachments) != 1 {
		t.Fatalf("wykaz załączników oddał %d pozycji zamiast jednej", len(wykaz.Attachments))
	}
	if len(wykaz.MissingKinds) != 1 || wykaz.MissingKinds[0] != string(shared.ResearchAttachmentKindSnapshot) {
		t.Fatalf("wykaz braków nie wskazał brakującej migawki: %v", wykaz.MissingKinds)
	}

	// Zdjęcie źródła ma zdjąć wiersz, a nie tylko zameldować zdjęcie.
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchSourceRemove,
		shared.ResearchSourceRemoveRequest{SourceId: dodane.Source.Id}, nil)
	if liczba := policzBadania(t, baza,
		`SELECT COUNT(*) FROM zrodlo_badania WHERE identyfikator_zewnetrzny = ?`,
		dodane.Source.Id); liczba != 0 {
		t.Fatalf("źródło zdjęte, a w bazie nadal leży (wierszy: %d)", liczba)
	}
}

// TestSkutekKodowaniaISprzecznosciBadania mierzy kodowanie jakościowe i detektor
// rozbieżności liczbowych: książka kodów, wiązania kodów oraz wiersz sprzeczności
// wraz z jej rozstrzygnięciem.
func TestSkutekKodowaniaISprzecznosciBadania(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaBadania(t, katalogDanych)
	okno := "okno-badania-analiza"

	var zrodlo shared.ResearchSourceAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchSourceAdd,
		shared.ResearchSourceAddRequest{
			WindowId: okno, Title: "Analiza segmentu X", Kind: shared.ResearchSourceKindWeb,
		}, &zrodlo)

	var pierwsze, drugie shared.ResearchFindingAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchFindingAdd,
		shared.ResearchFindingAddRequest{
			WindowId: okno, Content: "Segment X rośnie o 12 procent rocznie od 2023 roku",
			SourceIds: []string{zrodlo.Source.Id},
		}, &pierwsze)
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchFindingAdd,
		shared.ResearchFindingAddRequest{
			WindowId: okno, Content: "Segment X rośnie o 30 procent rocznie od 2023 roku",
			SourceIds: []string{zrodlo.Source.Id},
		}, &drugie)

	// Kodowanie: nowa nazwa kodu ma założyć pozycję książki kodów i wiązanie.
	var kodowanie shared.ResearchFindingCodeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchFindingCode,
		shared.ResearchFindingCodeRequest{
			FindingId: pierwsze.Finding.Id, NewCodeNames: []string{"wzrost"},
		}, &kodowanie)
	if liczba := policzBadania(t, baza,
		`SELECT COUNT(*) FROM kod_badania WHERE okno = ?`, okno); liczba != 1 {
		t.Fatalf("kod zameldowany, a książka kodów ma %d pozycji zamiast jednej", liczba)
	}
	if liczba := policzBadania(t, baza, `SELECT COUNT(*) FROM kod_ustalenia_badania ku
	    JOIN ustalenie_badania u ON u.id = ku.ustalenie_id
	    WHERE u.identyfikator_zewnetrzny = ?`, pierwsze.Finding.Id); liczba != 1 {
		t.Fatalf("kod przypisany, a wiązania w bazie nie ma (wierszy: %d)", liczba)
	}

	// Macierz kod × źródło ma liczyć z wiązań, nie z niczego.
	var macierz shared.ResearchFindingMatrixResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchFindingMatrix,
		shared.ResearchFindingMatrixRequest{WindowId: okno}, &macierz)
	if len(macierz.Cells) != 1 || macierz.Cells[0].Count != 1 {
		t.Fatalf("macierz kodowania nie policzyła jedynego wystąpienia: %+v", macierz.Cells)
	}

	// Sprzeczność: dwie liczby o tym samym przedmiocie mają zostawić wiersz.
	var sprzecznosci shared.ResearchFindingContradictionsResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchFindingContradictions,
		shared.ResearchFindingContradictionsRequest{WindowId: okno}, &sprzecznosci)
	if len(sprzecznosci.Contradictions) == 0 {
		t.Fatalf("detektor nie wskazał rozbieżności 12 wobec 30 przy tym samym przedmiocie")
	}
	if liczba := policzBadania(t, baza,
		`SELECT COUNT(*) FROM sprzecznosc_badania WHERE okno = ?`, okno); liczba == 0 {
		t.Fatalf("sprzeczność zameldowana, a w bazie jej nie ma")
	}

	kodSprzecznosci := sprzecznosci.Contradictions[0].Id
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchContradictionResolve,
		shared.ResearchContradictionResolveRequest{
			ContradictionId: kodSprzecznosci, Rationale: "źródło wtórne podało dane za inny okres",
		}, nil)
	if liczba := policzBadania(t, baza,
		`SELECT rozstrzygnieta FROM sprzecznosc_badania WHERE identyfikator_zewnetrzny = ?`,
		kodSprzecznosci); liczba != 1 {
		t.Fatalf("rozstrzygnięcie zameldowane, a wiersz sprzeczności go nie niesie")
	}

	// Ślad prowenancji ma narastać przy zmianie ustalenia.
	waga := shared.ResearchFindingWeight(shared.ResearchFindingWeightKey)
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchFindingUpdate,
		shared.ResearchFindingUpdateRequest{FindingId: pierwsze.Finding.Id, Weight: &waga}, nil)
	if liczba := policzBadania(t, baza,
		`SELECT COUNT(*) FROM prowenancja_badania WHERE ustalenie_kod = ?`,
		pierwsze.Finding.Id); liczba == 0 {
		t.Fatalf("zmiana ustalenia nie zostawiła śladu prowenancji")
	}
	if liczba := policzBadania(t, baza,
		`SELECT COUNT(*) FROM ustalenie_badania WHERE identyfikator_zewnetrzny = ? AND waga = 'key'`,
		pierwsze.Finding.Id); liczba != 1 {
		t.Fatalf("waga ustalenia zameldowana, a w bazie jej nie ma")
	}

	// Scalenie ustaleń ma zdjąć ustalenie scalone, a nie tylko je przepiąć.
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchFindingMerge,
		shared.ResearchFindingMergeRequest{
			TargetFindingId: pierwsze.Finding.Id, MergedFindingIds: []string{drugie.Finding.Id},
		}, nil)
	if liczba := policzBadania(t, baza,
		`SELECT COUNT(*) FROM ustalenie_badania WHERE identyfikator_zewnetrzny = ?`,
		drugie.Finding.Id); liczba != 0 {
		t.Fatalf("ustalenie scalone nadal leży w bazie")
	}
}

// TestSkutekPrzestrzeniIPokryciaBadania mierzy pytania badawcze, ich pokrycie
// źródłami, notatkę roboczą i wykaz luk.
func TestSkutekPrzestrzeniIPokryciaBadania(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaBadania(t, katalogDanych)
	okno := "okno-badania-przestrzen"

	odbiorca := "zarząd"
	var przestrzen shared.ResearchWorkspaceSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchWorkspaceSet,
		shared.ResearchWorkspaceSetRequest{
			Scope:  "Analiza konkurencyjna — segment X",
			Stages: []string{"rozpoznanie", "zbieranie", "analiza"},
			Questions: []shared.ResearchQuestion{
				{Id: "", Text: "Jak rośnie segment X?"},
				{Id: "", Text: "Kto jest liderem segmentu?"},
			},
			Audience: &odbiorca,
		}, &przestrzen)

	if liczba := policzBadania(t, baza, `SELECT COUNT(*) FROM pytanie_badania`); liczba != 2 {
		t.Fatalf("pytania badawcze zameldowane, a w bazie leży ich %d zamiast 2", liczba)
	}
	if liczba := policzBadania(t, baza,
		`SELECT COUNT(*) FROM przestrzen_badania WHERE odbiorca = 'zarząd'`); liczba != 1 {
		t.Fatalf("odbiorca raportu zameldowany, a w bazie go nie ma")
	}
	if len(przestrzen.Questions) != 2 {
		t.Fatalf("odpowiedź nie oddała zapisanych pytań: %+v", przestrzen.Questions)
	}

	// Pokrycie: pytanie bez źródła ma być luką, pytanie ze źródłem — pokryte.
	pierwszePytanie := przestrzen.Questions[0].Id
	var zrodlo shared.ResearchSourceAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchSourceAdd,
		shared.ResearchSourceAddRequest{
			WindowId: okno, Title: "Dane o wzroście segmentu",
			Kind:        shared.ResearchSourceKindWeb,
			QuestionIds: []string{pierwszePytanie},
		}, &zrodlo)
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchSourceUpdate,
		shared.ResearchSourceUpdateRequest{
			SourceId: zrodlo.Source.Id, QuestionIds: []string{pierwszePytanie},
		}, nil)
	if liczba := policzBadania(t, baza, `SELECT COUNT(*) FROM pytanie_zrodla_badania p
	    JOIN zrodlo_badania z ON z.id = p.zrodlo_id
	    WHERE z.identyfikator_zewnetrzny = ?`, zrodlo.Source.Id); liczba != 1 {
		t.Fatalf("przypisanie źródła do pytania nie zostawiło wiersza (wierszy: %d)", liczba)
	}

	var pokrycie shared.ResearchWorkspaceCoverageResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchWorkspaceCoverage,
		shared.ResearchWorkspaceCoverageRequest{WindowId: okno}, &pokrycie)
	if pokrycie.UncoveredCount != 1 {
		t.Fatalf("pokrycie policzone błędnie: niepokrytych %d zamiast 1", pokrycie.UncoveredCount)
	}

	var luki shared.ResearchGapFindResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchGapFind,
		shared.ResearchGapFindRequest{WindowId: okno}, &luki)
	if len(luki.Gaps) != 1 {
		t.Fatalf("wykaz luk oddał %d pozycji zamiast jednej", len(luki.Gaps))
	}

	// Notatka robocza ma wylądować w wierszu przestrzeni.
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchWorkspaceNoteSet,
		shared.ResearchWorkspaceNoteSetRequest{Content: "hipoteza: wzrost napędza segment premium"}, nil)
	if liczba := policzBadania(t, baza,
		`SELECT COUNT(*) FROM przestrzen_badania WHERE notatka LIKE 'hipoteza%'`); liczba != 1 {
		t.Fatalf("notatka robocza zameldowana, a w bazie jej nie ma")
	}
}

// TestSkutekRaportuIEksportuBadania mierzy najdalej idący skutek modułu: plik
// eksportu leżący w magazynie i niosący treść raportu.
func TestSkutekRaportuIEksportuBadania(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaBadania(t, katalogDanych)
	okno := "okno-badania-raport"

	var zrodlo shared.ResearchSourceAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchSourceAdd,
		shared.ResearchSourceAddRequest{
			WindowId: okno, Title: "Rocznik statystyczny 2024", Kind: shared.ResearchSourceKindWeb,
		}, &zrodlo)

	var ustalenie shared.ResearchFindingAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchFindingAdd,
		shared.ResearchFindingAddRequest{
			WindowId: okno, Content: "Udział segmentu X wyniósł 18 procent",
			SourceIds: []string{zrodlo.Source.Id},
		}, &ustalenie)

	tytul := "Analiza konkurencyjna — segment X"
	tresc := "Segment X rośnie szybciej niż rynek."
	var raport shared.ResearchReportBuildResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchReportBuild,
		shared.ResearchReportBuildRequest{
			WindowId: okno, Title: &tytul,
			Sections: []shared.ResearchReportSection{
				{Title: "Kontekst", Content: &tresc, FindingIds: []string{ustalenie.Finding.Id}},
			},
		}, &raport)

	if liczba := policzBadania(t, baza, `SELECT COUNT(*) FROM sekcja_raportu_badania s
	    JOIN raport_badania r ON r.id = s.raport_id
	    WHERE r.identyfikator_zewnetrzny = ?`, raport.Report.Id); liczba != 1 {
		t.Fatalf("sekcja raportu zameldowana, a w bazie leży jej %d zamiast jednej", liczba)
	}

	// Odczyt raportu bez wskazania identyfikatora ma oddać raport bieżący.
	var odczytany shared.ResearchReportGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchReportGet,
		shared.ResearchReportGetRequest{WindowId: okno}, &odczytany)
	if odczytany.Report == nil || odczytany.Report.Id != raport.Report.Id {
		t.Fatalf("odczyt raportu bieżącego nie oddał zbudowanego raportu: %+v", odczytany.Report)
	}

	// Wstawka tabeli dowodów ma zostawić wiersz bloku z policzonymi danymi.
	var blok shared.ResearchReportInsertResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchReportInsert,
		shared.ResearchReportInsertRequest{
			ReportId: raport.Report.Id, SectionId: raport.Report.Sections[0].Id,
			Kind: shared.ResearchBlockKindEvidenceTable,
		}, &blok)
	if liczba := policzBadania(t, baza,
		`SELECT COUNT(*) FROM blok_raportu_badania WHERE identyfikator_zewnetrzny = ?`,
		blok.Block.Id); liczba != 1 {
		t.Fatalf("wstawka zameldowana, a w bazie nie ma jej wiersza")
	}
	// Wstawka zakłada wersję raportu — panel wersji ma mieć do czego wrócić.
	if liczba := policzBadania(t, baza, `SELECT COUNT(*) FROM wersja_raportu_badania w
	    JOIN raport_badania r ON r.id = w.raport_id
	    WHERE r.identyfikator_zewnetrzny = ?`, raport.Report.Id); liczba == 0 {
		t.Fatalf("zmiana raportu nie założyła ani jednej wersji")
	}

	// Eksport: skutkiem jest PLIK, a nie wiersz o pliku.
	var eksport shared.ResearchReportExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchReportExport,
		shared.ResearchReportExportRequest{
			ReportId: raport.Report.Id, Format: shared.ExportFormatMarkdown,
		}, &eksport)
	if eksport.Path == nil || strings.TrimSpace(*eksport.Path) == "" {
		t.Fatalf("eksport zameldowany bez ścieżki wyniku — to jest wzorzec szkody, "+
			"którego ten sprawdzian pilnuje: %+v", eksport)
	}
	bajty, err := os.ReadFile(*eksport.Path)
	if err != nil {
		t.Fatalf("eksport wskazuje plik, którego nie ma: %v", err)
	}
	if len(bajty) == 0 {
		t.Fatalf("plik eksportu leży, ale jest pusty")
	}
	if !strings.Contains(string(bajty), tytul) || !strings.Contains(string(bajty), "Kontekst") {
		t.Fatalf("plik eksportu nie niesie treści raportu; początek: %q",
			poczatekTekstuBadania(string(bajty), 200))
	}

	// Historia eksportów ma czytać z bazy to samo, co zostało zapisane.
	var historia shared.ResearchExportListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchExportList,
		shared.ResearchExportListRequest{ReportId: &raport.Report.Id}, &historia)
	if len(historia.Exports) != 1 || historia.Exports[0].Path == nil {
		t.Fatalf("historia eksportów nie oddała wytworzonego pliku: %+v", historia.Exports)
	}
	if historia.Exports[0].SizeBytes == nil || *historia.Exports[0].SizeBytes != int64(len(bajty)) {
		t.Fatalf("rozmiar zapisany w historii nie zgadza się z plikiem na dysku")
	}

	// Udostępnienie ma wskazywać plik, który naprawdę leży.
	var udostepnienie shared.ResearchExportShareResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchExportShare,
		shared.ResearchExportShareRequest{ReportId: raport.Report.Id}, &udostepnienie)
	if !strings.HasPrefix(udostepnienie.Url, "file://") {
		t.Fatalf("udostępnienie oddało odnośnik, który do niczego nie prowadzi: %q", udostepnienie.Url)
	}
	if liczba := policzBadania(t, baza, `SELECT COUNT(*) FROM udostepnienie_raportu_badania u
	    JOIN raport_badania r ON r.id = u.raport_id
	    WHERE r.identyfikator_zewnetrzny = ?`, raport.Report.Id); liczba != 1 {
		t.Fatalf("udostępnienie zameldowane, a w bazie nie ma jego wiersza")
	}

	// Podgląd przed eksportem ma nieść treść, a nie pustą kopertę.
	var podglad shared.ResearchExportPreviewResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchExportPreview,
		shared.ResearchExportPreviewRequest{
			ReportId: raport.Report.Id, Format: shared.ExportFormatMarkdown,
		}, &podglad)
	if podglad.Preview.Text == nil || !strings.Contains(*podglad.Preview.Text, "Kontekst") {
		t.Fatalf("podgląd eksportu nie niesie treści raportu")
	}
}

// TestSkutekAdnotacjiIWypisowBadania mierzy adnotacje lektury oraz wypisy z nich
// złożone — skutkiem jest wiersz adnotacji, nie sam cytat w odpowiedzi.
func TestSkutekAdnotacjiIWypisowBadania(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaBadania(t, katalogDanych)
	okno := "okno-badania-lektura"

	var zrodlo shared.ResearchSourceAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchSourceAdd,
		shared.ResearchSourceAddRequest{
			WindowId: okno, Title: "Wywiad branżowy", Kind: shared.ResearchSourceKindWeb,
		}, &zrodlo)

	strona := 12
	cytat := "Segment X rośnie o 12% rocznie od 2023 roku"
	var adnotacja shared.ResearchAnnotationAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchAnnotationAdd,
		shared.ResearchAnnotationAddRequest{
			SourceId: zrodlo.Source.Id, Kind: shared.ResearchAnnotationKindHighlight,
			Anchor: shared.ResearchAnchor{Kind: shared.ResearchAnchorKindPage, Page: &strona},
			Quote:  &cytat,
		}, &adnotacja)

	if liczba := policzBadania(t, baza,
		`SELECT kotwica_strona FROM adnotacja_badania WHERE identyfikator_zewnetrzny = ?`,
		adnotacja.Annotation.Id); liczba != strona {
		t.Fatalf("kotwica adnotacji nie zapisała numeru strony (zapisano: %d)", liczba)
	}

	var wypisy shared.ResearchExcerptListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchExcerptList,
		shared.ResearchExcerptListRequest{WindowId: okno}, &wypisy)
	if wypisy.Total != 1 || wypisy.Excerpts[0].Quote != cytat {
		t.Fatalf("wypisy nie zebrały zapisanego podświetlenia: %+v", wypisy.Excerpts)
	}
	if wypisy.Excerpts[0].SourceTitle != "Wywiad branżowy" {
		t.Fatalf("wypis nie niesie tytułu źródła — odnośnik prowadzi donikąd")
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchAnnotationRemove,
		shared.ResearchAnnotationRemoveRequest{AnnotationId: adnotacja.Annotation.Id}, nil)
	if liczba := policzBadania(t, baza,
		`SELECT COUNT(*) FROM adnotacja_badania WHERE identyfikator_zewnetrzny = ?`,
		adnotacja.Annotation.Id); liczba != 0 {
		t.Fatalf("adnotacja zdjęta, a w bazie nadal leży")
	}
}

// TestSkutekPrzesiewuPrismaBadania mierzy przesiew: odrzucenie pozycji ma zmienić
// liczniki liczone z bazy, a nie tylko wrócić w odpowiedzi.
func TestSkutekPrzesiewuPrismaBadania(t *testing.T) {
	zmontowany, zycie, katalogDanych := zmontujDoPomiaruSkutku(t)
	baza := bazaBadania(t, katalogDanych)
	okno := "okno-badania-prisma"

	// Wynik odkrycia wchodzi tu wprost do bazy: sprawdzian mierzy przesiew,
	// nie łącze wyszukiwania.
	_, err := baza.Exec(`INSERT INTO wynik_odkrycia_badania (klucz, okno, tytul, dostawca)
	    VALUES ('crossref:10.1/a', ?, 'Praca A', 'crossref'),
	           ('crossref:10.1/b', ?, 'Praca B', 'crossref')`, okno, okno)
	if err != nil {
		t.Fatalf("nie można przygotować wyników odkrycia: %v", err)
	}

	var przesiew shared.ResearchDiscoveryRejectResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchDiscoveryReject,
		shared.ResearchDiscoveryRejectRequest{
			WindowId: okno, ResultKeys: []string{"crossref:10.1/b"},
			Reason: "poza zakresem tematycznym",
		}, &przesiew)

	if przesiew.Rejected != 1 {
		t.Fatalf("odrzucono %d pozycji zamiast jednej", przesiew.Rejected)
	}
	if liczba := policzBadania(t, baza,
		`SELECT COUNT(*) FROM wynik_odkrycia_badania WHERE okno = ? AND odrzucony = 1`,
		okno); liczba != 1 {
		t.Fatalf("odrzucenie zameldowane, a w bazie nie ma jego śladu")
	}

	var liczniki shared.ResearchPrismaGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandResearchPrismaGet,
		shared.ResearchPrismaGetRequest{WindowId: okno}, &liczniki)
	if liczniki.Counts.Identified != 2 || liczniki.Counts.Excluded != 1 {
		t.Fatalf("liczniki PRISMA nie zgadzają się z bazą: %+v", liczniki.Counts)
	}
	if len(liczniki.Counts.ExclusionReasons) != 1 {
		t.Fatalf("liczniki PRISMA nie niosą powodu wyłączenia: %+v", liczniki.Counts.ExclusionReasons)
	}
}
