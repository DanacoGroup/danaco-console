package core

import (
	"archive/zip"
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"danacoconsole/shared"
)

// Skutek dobudowanych rodzin modułu Translate.
//
// Wzorzec szkody, którego pilnuje ten plik, ma w tym produkcie precedens:
// komenda meldowała `status: ok` z pustym wynikiem, a za odpowiedzią nie leżało
// nic. Dlatego żaden sprawdzian tutaj nie kończy się na udanej odpowiedzi.
// Każdy schodzi niżej, do jednego z dwóch miejsc, w których skutek albo jest,
// albo go nie ma:
//
//   - do bazy — drugim, niezależnym połączeniem do tego samego pliku SQLite,
//     zapytaniem SQL wprost, z pominięciem całej warstwy adapterów;
//   - na dysk — otwarciem pliku, który komenda miała wytworzyć, i odczytaniem
//     jego treści (nie samego istnienia).
//
// Żaden sprawdzian nie woła modelu ani programu zewnętrznego, więc wszystkie
// wypadają tak samo u Operatora, jak na maszynie budującej.

// bazaSprawdzianu otwiera drugie połączenie do bazy stanowiska. Odczyt idzie
// nim, a nie przez rdzeń: gdyby szedł przez rdzeń, sprawdzian mierzyłby zgodność
// adaptera z samym sobą.
func bazaSprawdzianuTlumaczen(t *testing.T, katalog string) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", filepath.ToSlash(filepath.Join(katalog, "dane.sqlite"))+
		"?_pragma=busy_timeout(5000)")
	if err != nil {
		t.Fatalf("nie można otworzyć bazy obocznej: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// policzWierszeTlumaczenia liczy wiersze zapytaniem wprost.
func policzWierszeTlumaczenia(t *testing.T, db *sql.DB, zapytanie string, argumenty ...any) int {
	t.Helper()

	var ile int
	if err := db.QueryRow(zapytanie, argumenty...).Scan(&ile); err != nil {
		t.Fatalf("nie można policzyć wierszy (%s): %v", zapytanie, err)
	}
	return ile
}

// napisZBazy odczytuje jedną wartość tekstową zapytaniem wprost.
func napisZBazy(t *testing.T, db *sql.DB, zapytanie string, argumenty ...any) string {
	t.Helper()

	var wartosc sql.NullString
	if err := db.QueryRow(zapytanie, argumenty...).Scan(&wartosc); err != nil {
		t.Fatalf("nie można odczytać wartości (%s): %v", zapytanie, err)
	}
	return wartosc.String
}

// zalozOknoZrodlowe zakłada okno tłumaczenia z podanym tekstem źródłowym.
func zalozOknoZrodlowe(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno, tekst string) {
	t.Helper()

	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSourceSet,
		shared.TranslateSourceSetRequest{WindowId: okno, Text: tekst}, nil)
}

// zalozPanelBezModelu zakłada panel języka docelowego z treścią podaną wprost.
// Droga jest ta sama, którą idzie korekta Operatora (`translation.set`), więc
// sprawdzian nie potrzebuje ani kanału modelu, ani sieci.
func zalozPanelBezModelu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	db *sql.DB, okno, jezyk, tresc string) string {
	t.Helper()

	// Panel zakłada import XLIFF: jedyna droga rdzenia, która zakłada panel
	// wskazanego języka z gotową treścią i bez wołania modelu.
	sciezka := filepath.Join(t.TempDir(), "material.xlf")
	jednostki := make([]jednostkaXliff, 0)
	for _, akapit := range rozdzielAkapity(tresc) {
		jednostki = append(jednostki, jednostkaXliff{Zrodlo: akapit, Cel: akapit})
	}
	if err := os.WriteFile(sciezka, zlozXliff("polski", jezyk, jednostki), 0o600); err != nil {
		t.Fatalf("nie można zapisać materiału XLIFF: %v", err)
	}
	var wynik shared.TranslateXliffImportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateXliffImport,
		shared.TranslateXliffImportRequest{WindowId: okno, Path: sciezka}, &wynik)

	kod := ""
	for _, panel := range wynik.Panels {
		if panel.Language == jezyk {
			kod = panel.Id
		}
	}
	if kod == "" {
		t.Fatalf("import XLIFF nie założył panelu języka %s", jezyk)
	}
	// Treść panelu ustawiamy wprost, żeby sprawdzian pracował na dokładnie tym
	// tekście, który zadeklarował.
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateTranslationSet,
		shared.TranslateTranslationSetRequest{PanelId: kod, Text: tresc}, nil)

	if wBazie := napisZBazy(t, db,
		`SELECT tresc FROM panel_tlumaczenia WHERE identyfikator_zewnetrzny = ?`, kod); wBazie != tresc {
		t.Fatalf("baza niesie treść panelu %q, a zapisano %q", wBazie, tresc)
	}
	return kod
}

// TestPamiecTlumaczenLezyWBazieIWPliku wykazuje skutek rodziny `memory.*`:
// para wniesiona komendą realnie leży w tabeli, a eksport zostawia plik TMX,
// z którego tę samą parę da się odczytać z powrotem.
func TestPamiecTlumaczenLezyWBazieIWPliku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	var zapisany shared.TranslateMemorySetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateMemorySet,
		shared.TranslateMemorySetRequest{
			Language:      "angielski",
			SourceSegment: "Faktura wymaga podpisu.",
			TargetSegment: "The invoice requires a signature.",
			Project:       wskaznik("umowy"),
		}, &zapisany)

	wBazie := napisZBazy(t, db,
		`SELECT segment_docelowy FROM pamiec_tlumaczen WHERE identyfikator_zewnetrzny = ?`,
		zapisany.Entry.Id)
	if wBazie != "The invoice requires a signature." {
		t.Fatalf("za parą pamięci nie leży wiersz o oczekiwanej treści: %q", wBazie)
	}

	sciezka := filepath.Join(t.TempDir(), "pamiec.tmx")
	var wydana shared.TranslateMemoryExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateMemoryExport,
		shared.TranslateMemoryExportRequest{Path: sciezka}, &wydana)
	if wydana.ExportedCount != 1 {
		t.Fatalf("eksport zgłosił %d par, a w pamięci jest jedna", wydana.ExportedCount)
	}

	// Plik czytamy z powrotem tą samą drogą, którą czyta go import — czyli
	// mierzymy, że wynik jest prawdziwym TMX, a nie napisem o TMX.
	pary, err := wczytajParyWymiany(sciezka)
	if err != nil {
		t.Fatalf("wynik eksportu nie jest czytelnym plikiem wymiany: %v", err)
	}
	if len(pary) != 1 || pary[0].SegmentDocelowy != "The invoice requires a signature." {
		t.Fatalf("plik wyniku niesie %d par o treści %+v", len(pary), pary)
	}

	// Usunięcie ma zdejmować wiersz, nie samą odpowiedź.
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateMemoryDelete,
		shared.TranslateMemoryDeleteRequest{EntryId: zapisany.Entry.Id}, nil)
	if ile := policzWierszeTlumaczenia(t, db, `SELECT COUNT(*) FROM pamiec_tlumaczen`); ile != 0 {
		t.Fatalf("po usunięciu w tabeli pamięci zostało %d wierszy", ile)
	}
}

// TestTlumaczenieWstepneWypelniaPanelZPamieci wykazuje, że tłumaczenie wstępne
// realnie zmienia treść panelu w bazie, a nie tylko liczbę w odpowiedzi.
func TestTlumaczenieWstepneWypelniaPanelZPamieci(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	zalozOknoZrodlowe(t, zmontowany, zycie, "okno-pamieci", "Faktura wymaga podpisu.")
	kodPanelu := zalozPanelBezModelu(t, zmontowany, zycie, db, "okno-pamieci", "angielski",
		"Faktura wymaga podpisu.")

	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateMemorySet,
		shared.TranslateMemorySetRequest{
			Language:      "angielski",
			SourceSegment: "Faktura wymaga podpisu.",
			TargetSegment: "The invoice requires a signature.",
		}, nil)

	var wynik shared.TranslateMemoryPretranslateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateMemoryPretranslate,
		shared.TranslateMemoryPretranslateRequest{
			WindowId:  "okno-pamieci",
			OnlyEmpty: wskaznik(false),
		}, &wynik)
	if wynik.FilledCount == 0 {
		t.Fatal("tłumaczenie wstępne nie wypełniło ani jednego panelu")
	}

	wBazie := napisZBazy(t, db,
		`SELECT tresc FROM panel_tlumaczenia WHERE identyfikator_zewnetrzny = ?`, kodPanelu)
	if !strings.Contains(wBazie, "The invoice requires a signature.") {
		t.Fatalf("panel w bazie niesie %q — pary pamięci w nim nie ma", wBazie)
	}
}

// TestScalanieIPodzialSegmentowZostajeWBazie wykazuje, że zmiana podziału jest
// trwała: po scaleniu i podziale w tabeli segmentów leży dokładnie ten wykaz,
// który komenda oddała.
func TestScalanieIPodzialSegmentowZostajeWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	zalozOknoZrodlowe(t, zmontowany, zycie, "okno-segmentow",
		"Pierwsze zdanie. Drugie zdanie. Trzecie zdanie.")

	var scalone shared.TranslateSegmentMergeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSegmentMerge,
		shared.TranslateSegmentMergeRequest{
			WindowId:       "okno-segmentow",
			SegmentIndexes: []string{"0", "1"},
		}, &scalone)
	if len(scalone.Segments) != 2 {
		t.Fatalf("po scaleniu dwóch z trzech zdań zostało %d segmentów", len(scalone.Segments))
	}
	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM segment_okna_tlumaczenia`); ile != 2 {
		t.Fatalf("w tabeli segmentów leży %d wierszy, a scalenie dało dwa", ile)
	}
	pierwszy := napisZBazy(t, db,
		`SELECT tresc FROM segment_okna_tlumaczenia WHERE kolejnosc = 0`)
	if !strings.Contains(pierwszy, "Pierwsze zdanie.") || !strings.Contains(pierwszy, "Drugie zdanie.") {
		t.Fatalf("scalony segment w bazie niesie %q", pierwszy)
	}

	var podzielone shared.TranslateSegmentSplitResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSegmentSplit,
		shared.TranslateSegmentSplitRequest{
			WindowId:     "okno-segmentow",
			SegmentIndex: 0,
			Offset:       16,
		}, &podzielone)
	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM segment_okna_tlumaczenia`); ile != 3 {
		t.Fatalf("po podziale w tabeli segmentów leży %d wierszy zamiast trzech", ile)
	}
}

// TestKorektaZmieniaTrescPaneluWBazie wykazuje skutek korekty: ustalenie leży
// w tabeli, a jego zastosowanie realnie poprawia treść panelu.
func TestKorektaZmieniaTrescPaneluWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	zalozOknoZrodlowe(t, zmontowany, zycie, "okno-korekty", "Zdanie źródłowe.")
	kodPanelu := zalozPanelBezModelu(t, zmontowany, zycie, db, "okno-korekty", "angielski",
		"Sentence with  double spacing .")

	var korekta shared.TranslateProofreadRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateProofreadRun,
		shared.TranslateProofreadRunRequest{PanelId: kodPanelu}, &korekta)
	if len(korekta.Findings) == 0 {
		t.Fatal("korekta nie zgłosiła ani jednego ustalenia wobec treści z podwójnym odstępem")
	}
	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM ustalenie_korekty`); ile != len(korekta.Findings) {
		t.Fatalf("w tabeli ustaleń leży %d wierszy, a odpowiedź niosła %d",
			ile, len(korekta.Findings))
	}
	if len(korekta.Readability) == 0 {
		t.Fatal("korekta nie policzyła ani jednej miary czytelności")
	}

	doZastosowania := ""
	for _, ustalenie := range korekta.Findings {
		if ustalenie.Suggestion != nil && *ustalenie.Suggestion != "" {
			doZastosowania = ustalenie.Id
			break
		}
	}
	if doZastosowania == "" {
		t.Fatal("żadne ustalenie nie niesie propozycji poprawki")
	}

	var zastosowana shared.TranslateProofreadApplyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateProofreadApply,
		shared.TranslateProofreadApplyRequest{
			PanelId:    kodPanelu,
			FindingIds: []string{doZastosowania},
		}, &zastosowana)

	wBazie := napisZBazy(t, db,
		`SELECT tresc FROM panel_tlumaczenia WHERE identyfikator_zewnetrzny = ?`, kodPanelu)
	if wBazie == "Sentence with  double spacing ." {
		t.Fatal("po zastosowaniu poprawki treść panelu w bazie się nie zmieniła")
	}
	rozstrzygniete := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM ustalenie_korekty WHERE zastosowano IS NOT NULL`)
	if rozstrzygniete != 1 {
		t.Fatalf("rozstrzygniętych ustaleń w bazie: %d, a zastosowano jedno", rozstrzygniete)
	}
}

// TestObiegZatwierdzenZostawiaWierszIMigawke wykazuje, że zatwierdzenie zmienia
// dwa miejsca naraz: dokłada krok obiegu i przestawia migawkę panelu.
func TestObiegZatwierdzenZostawiaWierszIMigawke(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	zalozOknoZrodlowe(t, zmontowany, zycie, "okno-obiegu", "Zdanie źródłowe.")
	kodPanelu := zalozPanelBezModelu(t, zmontowany, zycie, db, "okno-obiegu", "angielski",
		"Source sentence.")

	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateApprovalSet,
		shared.TranslateApprovalSetRequest{
			PanelId: kodPanelu,
			Stage:   shared.ApprovalStageApproved,
			Note:    wskaznik("sprawdzone"),
		}, nil)

	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM zatwierdzenie_panelu`); ile != 1 {
		t.Fatalf("w tabeli obiegu leży %d wierszy zamiast jednego", ile)
	}
	etap := napisZBazy(t, db,
		`SELECT etap_zatwierdzenia FROM panel_tlumaczenia WHERE identyfikator_zewnetrzny = ?`,
		kodPanelu)
	if etap != string(shared.ApprovalStageApproved) {
		t.Fatalf("migawka panelu niesie etap %q zamiast `approved`", etap)
	}

	var wykaz shared.TranslateApprovalListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateApprovalList,
		shared.TranslateApprovalListRequest{PanelId: &kodPanelu}, &wykaz)
	if len(wykaz.Records) != 1 {
		t.Fatalf("wykaz obiegu oddał %d kroków zamiast jednego", len(wykaz.Records))
	}
}

// TestProfilKontroliJakosciZyjeWBazie wykazuje pełny cykl profilu: zapis,
// odczyt i usunięcie mierzone wierszami tabel, nie odpowiedziami.
func TestProfilKontroliJakosciZyjeWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	var zapisany shared.TranslateQaProfileSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateQaProfileSet,
		shared.TranslateQaProfileSetRequest{
			Name: "profil umów",
			Checks: []shared.QaProfileCheck{
				{Kind: shared.TranslationIssueKindNumber, Severity: shared.ProofreadSeverityError, Enabled: true},
				{Kind: shared.TranslationIssueKindCurrency, Severity: shared.ProofreadSeverityWarning, Enabled: true},
			},
		}, &zapisany)

	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM profil_qa_kontrola`); ile != 2 {
		t.Fatalf("kontroli profilu w bazie: %d, a podano dwie", ile)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateQaProfileDelete,
		shared.TranslateQaProfileDeleteRequest{ProfileId: zapisany.Profile.Id}, nil)
	if ile := policzWierszeTlumaczenia(t, db, `SELECT COUNT(*) FROM profil_qa`); ile != 0 {
		t.Fatalf("po usunięciu profilu w bazie zostało %d wierszy", ile)
	}
	// Kontrole schodzą razem z profilem — klucz obcy kaskadowy ma działać
	// naprawdę, a nie tylko stać w migracji.
	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM profil_qa_kontrola`); ile != 0 {
		t.Fatalf("po usunięciu profilu zostało %d jego kontroli", ile)
	}
}

// TestDokumentWczytanyLezyWBazieAWynikNaDysku wykazuje obie strony pracy na
// dokumencie: wczytanie zostawia segmenty w bazie, złożenie zostawia plik,
// którego treść da się odczytać.
func TestDokumentWczytanyLezyWBazieAWynikNaDysku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	material := filepath.Join(t.TempDir(), "umowa.md")
	if err := os.WriteFile(material,
		[]byte("Pierwszy akapit umowy.\n\nDrugi akapit umowy.\n"), 0o600); err != nil {
		t.Fatalf("nie można zapisać materiału: %v", err)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSourceSet,
		shared.TranslateSourceSetRequest{WindowId: "okno-dokumentu", Text: "wstęp"}, nil)

	var wczytany shared.TranslateDocumentLoadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateDocumentLoad,
		shared.TranslateDocumentLoadRequest{WindowId: "okno-dokumentu", Path: &material}, &wczytany)
	if len(wczytany.Segments) != 2 {
		t.Fatalf("dokument dał %d segmentów zamiast dwóch", len(wczytany.Segments))
	}
	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM segment_dokumentu_tlumaczenia`); ile != 2 {
		t.Fatalf("w tabeli segmentów dokumentu leży %d wierszy", ile)
	}

	kodPanelu := zalozPanelBezModelu(t, zmontowany, zycie, db, "okno-dokumentu", "angielski",
		"First paragraph.\n\nSecond paragraph.")

	wynikSciezka := filepath.Join(t.TempDir(), "umowa-en.md")
	var zlozony shared.TranslateDocumentRenderResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateDocumentRender,
		shared.TranslateDocumentRenderRequest{
			DocumentId: wczytany.Document.Id,
			PanelId:    kodPanelu,
			Path:       &wynikSciezka,
		}, &zlozony)

	tresc, err := os.ReadFile(zlozony.Path)
	if err != nil {
		t.Fatalf("za wynikiem złożenia nie leży plik: %v", err)
	}
	if !strings.Contains(string(tresc), "First paragraph.") {
		t.Fatalf("plik wyniku niesie %q — przekładu w nim nie ma", string(tresc))
	}
}

// TestZasobLokalizacyjnyPrzechodziPrzezBazeIPlik wykazuje, że klucze realnie
// lądują w bazie, a wydanie zostawia plik, z którego da się je odczytać.
func TestZasobLokalizacyjnyPrzechodziPrzezBazeIPlik(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	material := filepath.Join(t.TempDir(), "pl.json")
	if err := os.WriteFile(material,
		[]byte(`{"okno.tytul":"Ustawienia","okno.zapisz":"Zapisz {liczba} zmian"}`), 0o600); err != nil {
		t.Fatalf("nie można zapisać zasobu: %v", err)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSourceSet,
		shared.TranslateSourceSetRequest{WindowId: "okno-zasobu", Text: "wstęp"}, nil)

	var wczytany shared.TranslateResourceImportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateResourceImport,
		shared.TranslateResourceImportRequest{WindowId: "okno-zasobu", Path: material}, &wczytany)
	if wczytany.Resource.KeyCount != 2 {
		t.Fatalf("zasób zgłosił %d kluczy zamiast dwóch", wczytany.Resource.KeyCount)
	}
	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM klucz_lokalizacji`); ile != 2 {
		t.Fatalf("w tabeli kluczy leży %d wierszy", ile)
	}
	// Znacznik podstawienia ma być wypisany, bo na nim stoi kontrola jakości.
	znaczniki := napisZBazy(t, db,
		`SELECT znaczniki FROM klucz_lokalizacji WHERE klucz = 'okno.zapisz'`)
	if !strings.Contains(znaczniki, "{liczba}") {
		t.Fatalf("klucz nie ma wypisanego znacznika podstawienia: %q", znaczniki)
	}

	kodPanelu := zalozPanelBezModelu(t, zmontowany, zycie, db, "okno-zasobu", "angielski",
		"Settings\n\nSave {liczba} changes")

	wynikSciezka := filepath.Join(t.TempDir(), "en.json")
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateResourceExport,
		shared.TranslateResourceExportRequest{
			ResourceId: wczytany.Resource.Id,
			PanelId:    kodPanelu,
			Path:       wynikSciezka,
		}, nil)

	klucze, err := wczytajKluczeZasobu(wynikSciezka, shared.LocalizationResourceFormatJson)
	if err != nil {
		t.Fatalf("wynik wydania nie jest czytelnym zasobem: %v", err)
	}
	if len(klucze) != 2 {
		t.Fatalf("plik wyniku niesie %d kluczy", len(klucze))
	}

	// Formy mnogie mają realnie wejść do wiersza klucza.
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateResourcePluralApply,
		shared.TranslateResourcePluralApplyRequest{
			ResourceId: wczytany.Resource.Id,
			PanelId:    kodPanelu,
		}, nil)
	formy := napisZBazy(t, db,
		`SELECT formy_mnogie FROM klucz_lokalizacji WHERE klucz = 'okno.zapisz'`)
	if !strings.Contains(formy, "other") {
		t.Fatalf("klucz nie dostał form mnogich: %q", formy)
	}

	// Kontekst klucza również jest zapisem, nie odpowiedzią.
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateResourceKeyContextSet,
		shared.TranslateResourceKeyContextSetRequest{
			ResourceId: wczytany.Resource.Id,
			Key:        "okno.tytul",
			Context:    wskaznik("nagłówek okna ustawień"),
		}, nil)
	if kontekst := napisZBazy(t, db,
		`SELECT kontekst FROM klucz_lokalizacji WHERE klucz = 'okno.tytul'`); kontekst == "" {
		t.Fatal("kontekst klucza nie zapisał się w bazie")
	}
}

// TestNapisyPrzechodzaPrzezBazeIPlik wykazuje skutek rodziny napisów: kwestie
// leżą w bazie, a wydany plik SRT da się przeczytać z powrotem.
func TestNapisyPrzechodzaPrzezBazeIPlik(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	material := filepath.Join(t.TempDir(), "material.srt")
	if err := os.WriteFile(material, []byte(
		"1\n00:00:01,000 --> 00:00:03,000\nPierwsza kwestia.\n\n"+
			"2\n00:00:04,000 --> 00:00:06,000\nDruga kwestia.\n"), 0o600); err != nil {
		t.Fatalf("nie można zapisać napisów: %v", err)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSourceSet,
		shared.TranslateSourceSetRequest{WindowId: "okno-napisow", Text: "wstęp"}, nil)

	var wczytane shared.TranslateSubtitleImportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSubtitleImport,
		shared.TranslateSubtitleImportRequest{WindowId: "okno-napisow", Path: material}, &wczytane)
	if wczytane.ImportedCount != 2 {
		t.Fatalf("import wniósł %d kwestii zamiast dwóch", wczytane.ImportedCount)
	}
	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM kwestia_napisow WHERE panel_id IS NULL`); ile != 2 {
		t.Fatalf("w tabeli kwestii leży %d wierszy materiału źródłowego", ile)
	}

	kodPanelu := zalozPanelBezModelu(t, zmontowany, zycie, db, "okno-napisow", "angielski",
		"First cue.\n\nSecond cue.")

	wynikSciezka := filepath.Join(t.TempDir(), "wynik.srt")
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSubtitleExport,
		shared.TranslateSubtitleExportRequest{
			PanelId: kodPanelu,
			Path:    wynikSciezka,
			Format:  shared.SubtitleFormatSrt,
		}, nil)

	bajty, err := os.ReadFile(wynikSciezka)
	if err != nil {
		t.Fatalf("za wydaniem napisów nie leży plik: %v", err)
	}
	kwestie := kwestieZTekstu(string(bajty))
	if len(kwestie) != 2 || kwestie[0].Tresc != "First cue." {
		t.Fatalf("plik napisów niesie %d kwestii: %+v", len(kwestie), kwestie)
	}
	if kwestie[0].PoczatekMs != 1000 || kwestie[0].KoniecMs != 3000 {
		t.Fatalf("taktowanie przekładu rozjechało się ze źródłem: %+v", kwestie[0])
	}

	// Kontrola taktowania ma mierzyć realne kwestie, nie oddawać pustki.
	var zastrzezenia shared.TranslateSubtitleTimingCheckResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSubtitleTimingCheck,
		shared.TranslateSubtitleTimingCheckRequest{
			PanelId:        kodPanelu,
			CharsPerSecond: wskaznik(1),
		}, &zastrzezenia)
	if len(zastrzezenia.Issues) == 0 {
		t.Fatal("przy progu jednego znaku na sekundę kontrola nie zgłosiła ani jednego zastrzeżenia")
	}

	// Scenariusz dubbingu ma mieć tyle kwestii, ile ma materiał.
	var scenariusz shared.TranslateDubbingScriptBuildResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateDubbingScriptBuild,
		shared.TranslateDubbingScriptBuildRequest{PanelId: kodPanelu}, &scenariusz)
	if len(scenariusz.Lines) != 2 {
		t.Fatalf("scenariusz dubbingu ma %d kwestii zamiast dwóch", len(scenariusz.Lines))
	}
}

// TestPakietPrzekazaniaJestArchiwumZTrescia wykazuje, że pakiet to prawdziwe
// archiwum z wpisami, a jego zwrot realnie zmienia panel w bazie.
func TestPakietPrzekazaniaJestArchiwumZTrescia(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	zalozOknoZrodlowe(t, zmontowany, zycie, "okno-przekazania", "Zdanie do przekładu.")
	kodPanelu := zalozPanelBezModelu(t, zmontowany, zycie, db, "okno-przekazania", "angielski",
		"Sentence to translate.")

	sciezkaPakietu := filepath.Join(t.TempDir(), "pakiet.zip")
	var zlozony shared.TranslateHandoffBuildResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateHandoffBuild,
		shared.TranslateHandoffBuildRequest{
			WindowId: "okno-przekazania",
			Contents: []shared.HandoffContent{
				shared.HandoffContentXliff, shared.HandoffContentInstructions,
			},
			Instructions: wskaznik("proszę o rejestr formalny"),
			Path:         &sciezkaPakietu,
		}, &zlozony)

	archiwum, err := zip.OpenReader(zlozony.Path)
	if err != nil {
		t.Fatalf("za pakietem nie leży archiwum: %v", err)
	}
	defer archiwum.Close()
	nazwy := []string{}
	for _, wpis := range archiwum.File {
		nazwy = append(nazwy, wpis.Name)
	}
	if len(nazwy) < 2 {
		t.Fatalf("archiwum ma wpisy %v — miało nieść XLIFF i instrukcje", nazwy)
	}

	// Zwrot wykonawcy: ten sam XLIFF z poprawionym przekładem.
	zwrot := filepath.Join(t.TempDir(), "zwrot.xlf")
	if err := os.WriteFile(zwrot, zlozXliff("polski", "angielski", []jednostkaXliff{
		{Zrodlo: "Zdanie do przekładu.", Cel: "A sentence corrected by the vendor."},
	}), 0o600); err != nil {
		t.Fatalf("nie można zapisać zwrotu: %v", err)
	}
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateHandoffReceive,
		shared.TranslateHandoffReceiveRequest{
			WindowId:  "okno-przekazania",
			Path:      zwrot,
			PackageId: &zlozony.Package.Id,
		}, nil)

	wBazie := napisZBazy(t, db,
		`SELECT tresc FROM panel_tlumaczenia WHERE identyfikator_zewnetrzny = ?`, kodPanelu)
	if !strings.Contains(wBazie, "corrected by the vendor") {
		t.Fatalf("zwrot nie wszedł do panelu — baza niesie %q", wBazie)
	}
	if stan := napisZBazy(t, db,
		`SELECT stan FROM pakiet_przekazania WHERE identyfikator_zewnetrzny = ?`,
		zlozony.Package.Id); stan != string(shared.HandoffStatusReturned) {
		t.Fatalf("pakiet po zwrocie ma stan %q", stan)
	}
}

// TestWytworLezyWMagazynieIWBibliotece wykazuje skutek `artifact.publish`:
// za identyfikatorem pliku leży wiersz biblioteki, a za nim bajty na dysku.
func TestWytworLezyWMagazynieIWBibliotece(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	zalozOknoZrodlowe(t, zmontowany, zycie, "okno-wytworu", "Zdanie źródłowe.")
	zalozPanelBezModelu(t, zmontowany, zycie, db, "okno-wytworu", "angielski", "Source sentence.")

	var wydany shared.TranslateArtifactPublishResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateArtifactPublish,
		shared.TranslateArtifactPublishRequest{
			WindowId: "okno-wytworu",
			Kind:     shared.TranslationArtifactKindTargetFile,
		}, &wydany)

	odwolanie := napisZBazy(t, db,
		`SELECT tresc_odwolanie FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`,
		wydany.FileId)
	if odwolanie == "" {
		t.Fatal("wiersz pliku biblioteki nie niesie odwołania do bajtów")
	}
	bajty, err := os.ReadFile(odwolanie)
	if err != nil {
		t.Fatalf("pod odwołaniem wiersza nie leży plik: %v", err)
	}
	if !strings.Contains(string(bajty), "Source sentence.") {
		t.Fatalf("wytwór niesie %q — przekładu w nim nie ma", string(bajty))
	}
}

// TestPrzebiegPakietowyZostawiaPozycje wykazuje, że przebieg pakietowy zakłada
// pozycje i wykonuje pracę, którą da się zmierzyć w bazie.
func TestPrzebiegPakietowyZostawiaPozycje(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	zalozOknoZrodlowe(t, zmontowany, zycie, "okno-pakietu", "wstęp")
	zalozPanelBezModelu(t, zmontowany, zycie, db, "okno-pakietu", "angielski",
		"Amount and placeholder are gone.")
	// Tekst źródłowy ustawiamy PO panelu: import XLIFF wnosi własne źródło, więc
	// kolejność odwrotna zostawiłaby okno z tekstem identycznym jak przekład,
	// a kontrola jakości nie miałaby czego zgłosić.
	zalozOknoZrodlowe(t, zmontowany, zycie, "okno-pakietu", "Kwota 100 zł i {znacznik}.")

	var przebieg shared.TranslateBatchRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateBatchRun,
		shared.TranslateBatchRunRequest{
			WindowId:   "okno-pakietu",
			Operations: []shared.BatchOperationKind{shared.BatchOperationKindQualityCheck},
		}, &przebieg)
	if przebieg.ItemCount != 1 {
		t.Fatalf("przebieg ma %d pozycji zamiast jednej", przebieg.ItemCount)
	}
	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM pozycja_pakietu_tlumaczenia WHERE stan = 'done'`); ile != 1 {
		t.Fatalf("wykonanych pozycji w bazie: %d", ile)
	}
	// Kontrola jakości pakietu ma zostawić niezgodności — inaczej przebieg
	// zapisałby „zrobione” bez ani jednego skutku.
	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM panel_tlumaczenia_niezgodnosc`); ile == 0 {
		t.Fatal("przebieg zapisał pozycję jako wykonaną, a niezgodności nie ma ani jednej")
	}
}

// TestNastawyModuluLezaWBazie wykazuje, że nastawy — reguły segmentacji, profile
// silników i polityki — są zapisem, a nie odpowiedzią.
func TestNastawyModuluLezaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	db := bazaSprawdzianuTlumaczen(t, katalog)

	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateSegmentationRulesSet,
		shared.TranslateSegmentationRulesSetRequest{
			Name:     "polskie skróty",
			Language: wskaznik("polski"),
			Rules: []shared.SegmentationRule{
				{Order: 1, BeforeBreak: "np\\.", AfterBreak: " ", Breaks: false},
			},
		}, nil)
	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM regula_segmentacji`); ile != 1 {
		t.Fatalf("reguł segmentacji w bazie: %d", ile)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateEngineProfileSet,
		shared.TranslateEngineProfileSetRequest{
			Name:       "profil prawniczy",
			Domain:     wskaznik("prawo"),
			ChannelIds: []string{"kanal-a", "kanal-b"},
		}, nil)
	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT COUNT(*) FROM kanal_profilu_silnika`); ile != 2 {
		t.Fatalf("kanałów profilu w bazie: %d", ile)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslatePivotPolicySet,
		shared.TranslatePivotPolicySetRequest{
			DefaultPivot: wskaznik("angielski"),
			Pairs: []shared.PivotPair{
				{SourceLanguage: "polski", TargetLanguage: "japoński", PivotLanguage: "angielski"},
			},
		}, nil)
	if ile := policzWierszeTlumaczenia(t, db, `SELECT COUNT(*) FROM para_pivota`); ile != 1 {
		t.Fatalf("par pivota w bazie: %d", ile)
	}
	var odczytana shared.TranslatePivotPolicyGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslatePivotPolicyGet,
		shared.TranslatePivotPolicyGetRequest{}, &odczytana)
	if len(odczytana.Policy.Pairs) != 1 {
		t.Fatalf("odczyt polityki oddał %d par", len(odczytana.Policy.Pairs))
	}

	zalozOknoZrodlowe(t, zmontowany, zycie, "okno-polityki", "Zdanie źródłowe.")
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateMemoryPolicySet,
		shared.TranslateMemoryPolicySetRequest{
			WindowId:  "okno-polityki",
			Threshold: wskaznik(90),
		}, nil)
	if ile := policzWierszeTlumaczenia(t, db,
		`SELECT prog FROM polityka_pamieci_okna`); ile != 90 {
		t.Fatalf("próg polityki w bazie: %d", ile)
	}
}

// TestWykazKrokowWskazujeKomendyRejestru wykazuje, że wykaz kroków nie wymienia
// komend, których rdzeń nie obsługuje — inaczej okno prowadziłoby Operatora do
// czynności, której nie ma.
func TestWykazKrokowWskazujeKomendyRejestru(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var kroki shared.TranslateStepListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandTranslateStepList,
		shared.TranslateStepListRequest{}, &kroki)
	if len(kroki.Steps) == 0 {
		t.Fatal("wykaz kroków modułu jest pusty")
	}
	for _, krok := range kroki.Steps {
		if _, jest := zmontowany.Rdzen.rejestr.Obsluga(shared.MessageType(krok.Command)); !jest {
			t.Fatalf("krok %q wskazuje komendę %s, której rejestr nie obsługuje",
				krok.Name, krok.Command)
		}
	}
}
