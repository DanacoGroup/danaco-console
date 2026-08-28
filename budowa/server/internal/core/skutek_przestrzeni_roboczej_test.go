// Sprawdziany tego pliku mierzą skutek modułu Workspace niezależnie od
// odpowiedzi komendy: własnym zapytaniem SQL do pliku bazy albo odczytem
// pliku na dysku.
package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"danacoconsole/shared"
)

// projektSprawdzianuWorkspace jest kodem projektu wspólnym sprawdzianom pliku,
// żeby każdy sprawdzian pracował na tym samym, niezależnym od reszty zakresie
// danych.
const projektSprawdzianuWorkspace = "projekt-sprawdzianu"

// bazaWorkspace otwiera drugie połączenie do pliku bazy sprawdzianu. Miara ma
// iść obok rdzenia, nie przez jego repozytoria.
func bazaWorkspace(t *testing.T, katalog string) *sql.DB {
	t.Helper()

	polaczenie, err := sql.Open("sqlite", filepath.Join(katalog, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy do pomiaru: %v", err)
	}
	t.Cleanup(func() { _ = polaczenie.Close() })
	return polaczenie
}

// liczbaWorkspace odczytuje jedną liczbę własnym zapytaniem, z pominięciem
// warstwy adaptera, którą sprawdzian bada.
func liczbaWorkspace(t *testing.T, baza *sql.DB, zapytanie string, argumenty ...any) int {
	t.Helper()

	var wynik int
	if err := baza.QueryRow(zapytanie, argumenty...).Scan(&wynik); err != nil {
		t.Fatalf("zapytanie pomiaru %q nie powiodło się: %v", zapytanie, err)
	}
	return wynik
}

// napisWorkspace odczytuje jedną wartość tekstową własnym zapytaniem,
// z pominięciem warstwy adaptera, którą sprawdzian bada.
func napisWorkspace(t *testing.T, baza *sql.DB, zapytanie string, argumenty ...any) string {
	t.Helper()

	var wynik string
	if err := baza.QueryRow(zapytanie, argumenty...).Scan(&wynik); err != nil {
		t.Fatalf("zapytanie pomiaru %q nie powiodło się: %v", zapytanie, err)
	}
	return wynik
}

// zalozZadanieSprawdzianu zakłada zadanie komendą warsztatu i oddaje jego
// identyfikator wraz z pozostałymi polami odpowiedzi.
func zalozZadanieSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	tytul string, poczatek, termin int64) shared.WorkspaceTask {
	t.Helper()

	zadanie := shared.WorkspaceTaskCreateRequest{
		ProjectId: projektSprawdzianuWorkspace, Title: tytul,
	}
	if poczatek != 0 {
		zadanie.StartAt = wskaznik(poczatek)
	}
	if termin != 0 {
		zadanie.DueAt = wskaznik(termin)
	}
	var wynik shared.WorkspaceTaskCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceTaskCreate, zadanie, &wynik)
	return wynik.Task
}

// katalogProjektuSprawdzianu przestawia podstawę katalogu roboczego na katalog
// sprawdzianu i oddaje ścieżkę katalogu bibliotecznego projektu. Bez tego pliki
// projektu powstawałyby obok binarium sprawdzianu.
func katalogProjektuSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	katalog string) string {
	t.Helper()

	podstawa := filepath.Join(katalog, "warsztat")
	wartosc, err := json.Marshal(podstawa)
	if err != nil {
		t.Fatalf("nie można złożyć wartości ustawienia: %v", err)
	}
	wykonajUdana(t, zmontowany, zycie, shared.CommandConfigSet, shared.ConfigSetRequest{
		Key: KluczKatalogRoboczyPodstawa, Value: wartosc, Scope: shared.ConfigScopeGlobal,
	}, nil)

	sciezka := SciezkaSesji(podstawa, "", "projekt-"+projektSprawdzianuWorkspace)
	if err := os.MkdirAll(sciezka, 0o750); err != nil {
		t.Fatalf("nie można założyć katalogu projektu: %v", err)
	}
	return sciezka
}

// TestZadanieProjektuZostajeWBazie wykazuje skutek: po `workspace.task.create`
// wiersz zadania leży w bazie, a `workspace.task.move` naprawdę przestawia jego
// stan i kolumnę — mierzone własnym zapytaniem, nie odpowiedzią komendy.
func TestZadanieProjektuZostajeWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaWorkspace(t, katalog)

	zadanie := zalozZadanieSprawdzianu(t, zmontowany, zycie, "Przygotować raport", 0, 0)

	tytul := napisWorkspace(t, baza,
		`SELECT tytul FROM zadanie_projektu WHERE identyfikator_zewnetrzny = ?`, zadanie.Id)
	if tytul != "Przygotować raport" {
		t.Fatalf("w bazie leży zadanie o tytule %q, a założono „Przygotować raport”", tytul)
	}
	stan := napisWorkspace(t, baza,
		`SELECT stan FROM zadanie_projektu WHERE identyfikator_zewnetrzny = ?`, zadanie.Id)
	if stan != shared.WorkspaceTaskStatusTodo {
		t.Fatalf("zadanie nowe stoi w stanie %q zamiast todo", stan)
	}

	var przeniesione shared.WorkspaceTaskMoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceTaskMove,
		shared.WorkspaceTaskMoveRequest{
			TaskId:        zadanie.Id,
			BoardColumnId: wskaznik("kolumna-" + shared.WorkspaceTaskStatusInProgress),
		}, &przeniesione)

	stanPoPrzeniesieniu := napisWorkspace(t, baza,
		`SELECT stan FROM zadanie_projektu WHERE identyfikator_zewnetrzny = ?`, zadanie.Id)
	if stanPoPrzeniesieniu != shared.WorkspaceTaskStatusInProgress {
		t.Fatalf("po przeniesieniu karty baza trzyma stan %q zamiast inProgress", stanPoPrzeniesieniu)
	}
	kolumna := napisWorkspace(t, baza,
		`SELECT kolumna_tablicy FROM zadanie_projektu WHERE identyfikator_zewnetrzny = ?`, zadanie.Id)
	if !strings.HasSuffix(kolumna, shared.WorkspaceTaskStatusInProgress) {
		t.Fatalf("karta stoi w kolumnie %q, a przeniesiono ją do kolumny stanu inProgress", kolumna)
	}

	var wykaz shared.WorkspaceTaskListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceTaskList,
		shared.WorkspaceTaskListRequest{ProjectId: projektSprawdzianuWorkspace}, &wykaz)
	if wykaz.Total != 1 {
		t.Fatalf("wykaz zadań mówi o %d pozycjach, a w projekcie stoi jedno zadanie", wykaz.Total)
	}

	var usuniete shared.WorkspaceTaskDeleteResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceTaskDelete,
		shared.WorkspaceTaskDeleteRequest{TaskId: zadanie.Id}, &usuniete)
	if !usuniete.Deleted {
		t.Fatal("rdzeń zameldował, że zadania nie było, choć przed chwilą je zakładał")
	}
	pozostale := liczbaWorkspace(t, baza,
		`SELECT count(*) FROM zadanie_projektu WHERE identyfikator_zewnetrzny = ?`, zadanie.Id)
	if pozostale != 0 {
		t.Fatalf("po usunięciu w bazie zostało %d wierszy zadania", pozostale)
	}
}

// TestZaleznoscPrzesuwaTerminNastepnika wykazuje skutek zależności: termin
// zadania następującego naprawdę przesuwa się w bazie, a zależność domykająca
// cykl nie zostaje zapisana.
func TestZaleznoscPrzesuwaTerminNastepnika(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaWorkspace(t, katalog)

	const godzina = int64(3_600_000)
	pierwsze := zalozZadanieSprawdzianu(t, zmontowany, zycie, "Zebranie danych",
		10*godzina, 20*godzina)
	drugie := zalozZadanieSprawdzianu(t, zmontowany, zycie, "Analiza", 12*godzina, 14*godzina)

	var zaleznosc shared.WorkspaceTaskDependencySetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceTaskDependencySet,
		shared.WorkspaceTaskDependencySetRequest{
			ProjectId:         projektSprawdzianuWorkspace,
			PredecessorTaskId: pierwsze.Id, SuccessorTaskId: drugie.Id,
		}, &zaleznosc)

	poczatek := liczbaWorkspace(t, baza,
		`SELECT poczatek_ms FROM zadanie_projektu WHERE identyfikator_zewnetrzny = ?`, drugie.Id)
	if int64(poczatek) < 20*godzina {
		t.Fatalf("następnik zaczyna się w %d, a poprzednik kończy się w %d — "+
			"zależność nie przesunęła terminu", poczatek, 20*godzina)
	}

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandWorkspaceTaskDependencySet,
		shared.WorkspaceTaskDependencySetRequest{
			ProjectId:         projektSprawdzianuWorkspace,
			PredecessorTaskId: drugie.Id, SuccessorTaskId: pierwsze.Id,
		})
	if !strings.Contains(odmowa.Message, "cykl") {
		t.Fatalf("odmowa nie nazywa cyklu: %s", odmowa.Message)
	}
	krawedzie := liczbaWorkspace(t, baza, `SELECT count(*) FROM zaleznosc_zadan_projektu`)
	if krawedzie != 1 {
		t.Fatalf("w bazie leży %d krawędzi zależności, a zapisana miała zostać jedna", krawedzie)
	}

	var harmonogram shared.WorkspaceScheduleGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceScheduleGet,
		shared.WorkspaceScheduleGetRequest{ProjectId: projektSprawdzianuWorkspace}, &harmonogram)
	if len(harmonogram.Bars) != 2 {
		t.Fatalf("harmonogram niesie %d słupków, a zadania z granicami czasu są dwa",
			len(harmonogram.Bars))
	}

	var zniesienie shared.WorkspaceTaskDependencyRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceTaskDependencyRemove,
		shared.WorkspaceTaskDependencyRemoveRequest{DependencyId: zaleznosc.Dependency.Id},
		&zniesienie)
	if !zniesienie.Removed {
		t.Fatal("zniesienie zależności zameldowało brak krawędzi, którą przed chwilą założono")
	}
	if pozostale := liczbaWorkspace(t, baza,
		`SELECT count(*) FROM zaleznosc_zadan_projektu`); pozostale != 0 {
		t.Fatalf("po zniesieniu w bazie zostało %d krawędzi", pozostale)
	}
}

// TestTablicaProjektuMaKolumnyIKarty wykazuje skutek: kolumny tablicy powstają
// w bazie przy pierwszym wejściu, a karta stoi w kolumnie swojego stanu.
func TestTablicaProjektuMaKolumnyIKarty(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaWorkspace(t, katalog)

	zalozZadanieSprawdzianu(t, zmontowany, zycie, "Karta pierwsza", 0, 0)

	var tablica shared.WorkspaceBoardGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceBoardGet,
		shared.WorkspaceBoardGetRequest{ProjectId: projektSprawdzianuWorkspace}, &tablica)

	kolumny := liczbaWorkspace(t, baza, `SELECT count(*) FROM kolumna_tablicy_projektu`)
	if kolumny != len(kolumnyWyjscioweWorkspace) {
		t.Fatalf("w bazie leży %d kolumn tablicy, a zestaw wyjściowy liczy %d",
			kolumny, len(kolumnyWyjscioweWorkspace))
	}
	if len(tablica.Board.Tasks) != 1 {
		t.Fatalf("tablica niesie %d kart, a w projekcie stoi jedna", len(tablica.Board.Tasks))
	}
	if tablica.Board.Tasks[0].BoardColumnId == nil ||
		*tablica.Board.Tasks[0].BoardColumnId == "" {
		t.Fatal("karta tablicy nie wskazuje żadnej kolumny — nie ma jej gdzie narysować")
	}
}

// TestNotatkaZapisujeOdnosnikiIOdnosnikiWsteczne wykazuje skutek notatek:
// treść leży w bazie, odnośniki `[[nazwa]]` mają własne wiersze, a panel „co
// linkuje tutaj" znajduje stronę linkującą.
func TestNotatkaZapisujeOdnosnikiIOdnosnikiWsteczne(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaWorkspace(t, katalog)

	var docelowa shared.WorkspaceNoteSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceNoteSave,
		shared.WorkspaceNoteSaveRequest{
			ProjectId: projektSprawdzianuWorkspace, Title: "Ustalenia",
			Content: "# Ustalenia\n\nTreść strony docelowej.",
		}, &docelowa)

	var linkujaca shared.WorkspaceNoteSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceNoteSave,
		shared.WorkspaceNoteSaveRequest{
			ProjectId: projektSprawdzianuWorkspace, Title: "Protokół spotkania",
			Content:      "Zapisano w [[Ustalenia]] oraz w [[Strona bez treści]].",
			ParentNoteId: wskaznik(docelowa.Note.Id),
		}, &linkujaca)

	if len(linkujaca.MissingNames) != 1 || linkujaca.MissingNames[0] != "Strona bez treści" {
		t.Fatalf("zapis nie nazwał odnośnika bez strony: %+v", linkujaca.MissingNames)
	}
	tresc := napisWorkspace(t, baza,
		`SELECT tresc FROM notatka_projektu WHERE identyfikator_zewnetrzny = ?`, linkujaca.Note.Id)
	if !strings.Contains(tresc, "[[Ustalenia]]") {
		t.Fatalf("w bazie leży treść notatki bez odnośnika: %q", tresc)
	}
	odnosniki := liczbaWorkspace(t, baza,
		`SELECT count(*) FROM odnosnik_notatki_projektu WHERE notatka_zrodlowa = ?`, linkujaca.Note.Id)
	if odnosniki != 2 {
		t.Fatalf("w bazie leży %d odnośników notatki, a treść niesie dwa", odnosniki)
	}
	domkniete := liczbaWorkspace(t, baza,
		`SELECT count(*) FROM odnosnik_notatki_projektu
		 WHERE notatka_zrodlowa = ? AND notatka_docelowa = ?`, linkujaca.Note.Id, docelowa.Note.Id)
	if domkniete != 1 {
		t.Fatal("odnośnik do istniejącej strony nie został z nią związany w bazie")
	}

	var wsteczne shared.WorkspaceNoteBacklinkListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceNoteBacklinkList,
		shared.WorkspaceNoteBacklinkListRequest{NoteId: docelowa.Note.Id}, &wsteczne)
	if wsteczne.Total != 1 || wsteczne.Backlinks[0].SourceNoteId != linkujaca.Note.Id {
		t.Fatalf("panel „co linkuje tutaj” nie znalazł strony linkującej: %+v", wsteczne)
	}

	var drzewo shared.WorkspaceNoteTreeGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceNoteTreeGet,
		shared.WorkspaceNoteTreeGetRequest{ProjectId: projektSprawdzianuWorkspace}, &drzewo)
	if len(drzewo.Nodes) != 2 {
		t.Fatalf("drzewo stron niesie %d węzłów, a stron są dwie", len(drzewo.Nodes))
	}

	var graf shared.WorkspaceKnowledgeGraphGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceKnowledgeGraphGet,
		shared.WorkspaceKnowledgeGraphGetRequest{ProjectId: projektSprawdzianuWorkspace}, &graf)
	if len(graf.Graph.Edges) == 0 {
		t.Fatal("graf wiedzy nie ma ani jednej krawędzi, choć strona linkuje do strony")
	}

	var usuniecie shared.WorkspaceNoteDeleteResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceNoteDelete,
		shared.WorkspaceNoteDeleteRequest{NoteId: linkujaca.Note.Id}, &usuniecie)
	if pozostale := liczbaWorkspace(t, baza,
		`SELECT count(*) FROM notatka_projektu`); pozostale != 1 {
		t.Fatalf("po usunięciu notatki w bazie zostało %d stron zamiast jednej", pozostale)
	}
}

// TestKalendarzWciagaIcalIOddajePlik wykazuje skutek kalendarza: wciągnięcie
// zakłada zadania w bazie, a zapis oddaje treść, którą własny czytnik odczyta
// z powrotem.
func TestKalendarzWciagaIcalIOddajePlik(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaWorkspace(t, katalog)

	plik := strings.Join([]string{
		"BEGIN:VCALENDAR", "VERSION:2.0",
		"BEGIN:VEVENT", "UID:spotkanie-1", "SUMMARY:Spotkanie z klientem",
		"DTSTART:20260901T090000Z", "DTEND:20260901T100000Z", "END:VEVENT",
		"BEGIN:VEVENT", "UID:bez-tytulu", "DTSTART:20260902T090000Z", "END:VEVENT",
		"END:VCALENDAR",
	}, "\r\n")

	var wciagniecie shared.WorkspaceCalendarImportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceCalendarImport,
		shared.WorkspaceCalendarImportRequest{
			ProjectId: projektSprawdzianuWorkspace, Content: plik,
			FileName: wskaznik("kalendarz.ics"),
		}, &wciagniecie)

	if wciagniecie.Imported != 1 {
		t.Fatalf("wciągnięto %d pozycji, a plik niesie jedno wydarzenie z tytułem",
			wciagniecie.Imported)
	}
	if wciagniecie.Skipped != 1 || len(wciagniecie.SkippedReasons) != 1 {
		t.Fatalf("pominięcie wydarzenia bez tytułu nie wróciło powodem: %+v", wciagniecie)
	}
	zadania := liczbaWorkspace(t, baza,
		`SELECT count(*) FROM zadanie_projektu WHERE tytul = 'Spotkanie z klientem'`)
	if zadania != 1 {
		t.Fatalf("w bazie leży %d zadań z wciągniętego kalendarza", zadania)
	}

	var zapis shared.WorkspaceCalendarExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceCalendarExport,
		shared.WorkspaceCalendarExportRequest{ProjectId: projektSprawdzianuWorkspace}, &zapis)
	if zapis.Exported != 1 {
		t.Fatalf("zapis kalendarza mówi o %d pozycjach, a projekt ma jedną", zapis.Exported)
	}
	odczytane, _ := rozbierzIcalWorkspace(zapis.Content)
	if len(odczytane) != 1 || odczytane[0].Tytul != "Spotkanie z klientem" {
		t.Fatalf("zapisany plik iCal nie daje się odczytać z powrotem: %+v", odczytane)
	}

	var siatka shared.WorkspaceCalendarGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceCalendarGet,
		shared.WorkspaceCalendarGetRequest{
			ProjectId: projektSprawdzianuWorkspace,
			Span:      shared.WorkspaceCalendarSpanMonth,
			AnchorAt:  odczytane[0].PoczatekMs,
		}, &siatka)
	if len(siatka.Entries) != 1 {
		t.Fatalf("siatka miesiąca niesie %d pozycji, a wydarzenie jest jedno",
			len(siatka.Entries))
	}
}

// TestWyciagTekstuWchodziDoWyszukiwania wykazuje skutek: treść pliku projektu
// trafia do wskaźnika w bazie i znajduje się wyszukiwaniem po słowie, którego
// nie ma w nazwie pliku.
func TestWyciagTekstuWchodziDoWyszukiwania(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaWorkspace(t, katalog)
	katalogProjektu := katalogProjektuSprawdzianu(t, zmontowany, zycie, katalog)

	zapiszPlikSprawdzianu(t, filepath.Join(katalogProjektu, "notatka-terenowa.txt"),
		[]byte("Ustalono wynagrodzenie ryczałtowe w wysokości stu tysięcy."))

	var wydobycie shared.WorkspaceLibraryTextExtractResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceLibraryTextExtract,
		shared.WorkspaceLibraryTextExtractRequest{
			ProjectId: projektSprawdzianuWorkspace, FileId: "notatka-terenowa.txt",
		}, &wydobycie)

	if wydobycie.Extraction.CharacterCount == 0 {
		t.Fatal("wydobycie zameldowało zero znaków — pliku nie przeczytano")
	}
	tresc := napisWorkspace(t, baza,
		`SELECT tresc FROM wyciag_tekstu_projektu WHERE plik = ?`, "notatka-terenowa.txt")
	if !strings.Contains(tresc, "ryczałtowe") {
		t.Fatalf("wskaźnik treści w bazie nie niesie treści pliku: %q", tresc)
	}

	var wyszukanie shared.WorkspaceSearchProjectResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceSearchProject,
		shared.WorkspaceSearchProjectRequest{
			ProjectId: projektSprawdzianuWorkspace, Query: "ryczałtowe",
		}, &wyszukanie)
	if wyszukanie.Total == 0 {
		t.Fatal("wyszukiwanie nie znalazło słowa, które leży w treści pliku projektu")
	}
	if wyszukanie.Hits[0].Snippet == nil || *wyszukanie.Hits[0].Snippet == "" {
		t.Fatal("trafienie w treść wróciło bez fragmentu — okno nie ma czego podświetlić")
	}
}

// TestDuplikatyScalajaSieDoJednegoPliku wykazuje skutek scalenia: plik scalony
// znika z dysku, a plik zachowany na nim zostaje.
func TestDuplikatyScalajaSieDoJednegoPliku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	katalogProjektu := katalogProjektuSprawdzianu(t, zmontowany, zycie, katalog)

	tresc := []byte("Umowa ramowa — treść identyczna w obu plikach.")
	zapiszPlikSprawdzianu(t, filepath.Join(katalogProjektu, "umowa.txt"), tresc)
	zapiszPlikSprawdzianu(t, filepath.Join(katalogProjektu, "umowa-kopia.txt"), tresc)
	zapiszPlikSprawdzianu(t, filepath.Join(katalogProjektu, "inny.txt"), []byte("Co innego."))

	var duplikaty shared.WorkspaceLibraryDuplicateListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceLibraryDuplicateList,
		shared.WorkspaceLibraryDuplicateListRequest{ProjectId: projektSprawdzianuWorkspace},
		&duplikaty)
	if duplikaty.Total != 1 {
		t.Fatalf("wykryto %d grup duplikatów, a identyczne pliki są dwa", duplikaty.Total)
	}

	var scalenie shared.WorkspaceLibraryDuplicateMergeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceLibraryDuplicateMerge,
		shared.WorkspaceLibraryDuplicateMergeRequest{
			ProjectId: projektSprawdzianuWorkspace, KeepFileId: "umowa.txt",
			MergedFileIds: []string{"umowa-kopia.txt"},
		}, &scalenie)
	if scalenie.Merged != 1 {
		t.Fatalf("scalono %d plików, a wskazano jeden", scalenie.Merged)
	}
	if _, err := os.Stat(filepath.Join(katalogProjektu, "umowa-kopia.txt")); err == nil {
		t.Fatal("plik scalony nadal leży na dysku — scalenie zameldowało skutek, którego nie ma")
	}
	if _, err := os.Stat(filepath.Join(katalogProjektu, "umowa.txt")); err != nil {
		t.Fatalf("plik zachowany zniknął z dysku: %v", err)
	}
}

// TestKomentarzeIOsCzasuZostajaWBazie wykazuje skutek współpracy: komentarz ma
// wiersz w bazie, oś czasu zapisuje zdarzenia czynności, a usunięcie wątku
// zabiera także odpowiedzi.
func TestKomentarzeIOsCzasuZostajaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaWorkspace(t, katalog)

	zadanie := zalozZadanieSprawdzianu(t, zmontowany, zycie, "Zadanie z komentarzem", 0, 0)

	var komentarz shared.WorkspaceCommentAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceCommentAdd,
		shared.WorkspaceCommentAddRequest{
			ProjectId:  projektSprawdzianuWorkspace,
			TargetKind: shared.WorkspaceEntityKindTask, TargetId: zadanie.Id,
			Content: "Proszę o korektę @analityk",
		}, &komentarz)

	if len(komentarz.MentionedIds) != 1 || komentarz.MentionedIds[0] != "analityk" {
		t.Fatalf("przywołanie w treści nie zostało rozpoznane: %+v", komentarz.MentionedIds)
	}
	tresc := napisWorkspace(t, baza,
		`SELECT tresc FROM komentarz_projektu WHERE identyfikator_zewnetrzny = ?`,
		komentarz.Comment.Id)
	if !strings.Contains(tresc, "korektę") {
		t.Fatalf("w bazie leży komentarz o treści %q", tresc)
	}

	var odpowiedz shared.WorkspaceCommentAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceCommentAdd,
		shared.WorkspaceCommentAddRequest{
			ProjectId:  projektSprawdzianuWorkspace,
			TargetKind: shared.WorkspaceEntityKindTask, TargetId: zadanie.Id,
			Content: "Poprawione.", ParentCommentId: wskaznik(komentarz.Comment.Id),
		}, &odpowiedz)

	var wykaz shared.WorkspaceCommentListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceCommentList,
		shared.WorkspaceCommentListRequest{
			ProjectId: projektSprawdzianuWorkspace,
			TargetId:  wskaznik(zadanie.Id),
		}, &wykaz)
	if wykaz.Total != 2 {
		t.Fatalf("wątek niesie %d komentarzy, a zapisano dwa", wykaz.Total)
	}

	var os shared.WorkspaceActivityListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceActivityList,
		shared.WorkspaceActivityListRequest{ProjectId: projektSprawdzianuWorkspace}, &os)
	zdarzenia := liczbaWorkspace(t, baza, `SELECT count(*) FROM zdarzenie_projektu`)
	if zdarzenia == 0 || os.Total == 0 {
		t.Fatalf("oś czasu jest pusta mimo trzech czynności w projekcie (baza: %d, wykaz: %d)",
			zdarzenia, os.Total)
	}

	var usuniecie shared.WorkspaceCommentDeleteResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceCommentDelete,
		shared.WorkspaceCommentDeleteRequest{CommentId: komentarz.Comment.Id}, &usuniecie)
	if pozostale := liczbaWorkspace(t, baza,
		`SELECT count(*) FROM komentarz_projektu`); pozostale != 0 {
		t.Fatalf("po usunięciu wątku w bazie zostało %d komentarzy", pozostale)
	}
	_ = odpowiedz
}

// TestTablicaWizualnaTrzymaScene wykazuje skutek tablicy wizualnej: scena leży
// w bazie w tej postaci, w której przyszła, i wraca odczytem.
func TestTablicaWizualnaTrzymaScene(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaWorkspace(t, katalog)

	scena := jsonSurowy(t, map[string]any{
		"wezly": []map[string]any{{"id": "k1", "tekst": "Cel projektu", "x": 10, "y": 20}},
	})
	var zapis shared.WorkspaceCanvasSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceCanvasSave,
		shared.WorkspaceCanvasSaveRequest{
			ProjectId: projektSprawdzianuWorkspace, Name: wskaznik("Mapa myśli"), Scene: scena,
		}, &zapis)

	zapisana := napisWorkspace(t, baza,
		`SELECT scena FROM tablica_wizualna_projektu WHERE identyfikator_zewnetrzny = ?`,
		zapis.Canvas.Id)
	if !strings.Contains(zapisana, "Cel projektu") {
		t.Fatalf("w bazie leży scena bez treści kartki: %q", zapisana)
	}

	var odczyt shared.WorkspaceCanvasGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceCanvasGet,
		shared.WorkspaceCanvasGetRequest{ProjectId: projektSprawdzianuWorkspace}, &odczyt)
	if odczyt.Canvas == nil || odczyt.Canvas.Name != "Mapa myśli" {
		t.Fatalf("odczyt tablicy nie oddał zapisanej tablicy: %+v", odczyt)
	}
	if !strings.Contains(string(odczyt.Canvas.Scene), "Cel projektu") {
		t.Fatal("odczyt oddał tablicę bez sceny — płótno wróciło puste")
	}
}

// TestWersjeInstrukcjiIPrzywrocenie wykazuje skutek historii instrukcji: każdy
// zapis odkłada wersję w bazie, a przywrócenie zakłada wersję nową o treści
// wskazanej, zamiast przepisywać historię.
func TestWersjeInstrukcjiIPrzywrocenie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaWorkspace(t, katalog)

	for _, tresc := range []string{"Instrukcja pierwsza.", "Instrukcja druga."} {
		wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceInstructionsSet,
			shared.WorkspaceInstructionsSetRequest{
				ProjectId: projektSprawdzianuWorkspace, Content: tresc,
			}, nil)
	}
	if wersje := liczbaWorkspace(t, baza,
		`SELECT count(*) FROM wersja_instrukcji_projektu`); wersje != 2 {
		t.Fatalf("w bazie leży %d wersji instrukcji, a zapisów były dwa", wersje)
	}

	var wykaz shared.WorkspaceInstructionsVersionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceInstructionsVersionList,
		shared.WorkspaceInstructionsVersionListRequest{ProjectId: projektSprawdzianuWorkspace},
		&wykaz)
	if wykaz.Total != 2 || wykaz.Versions[0].Content != "Instrukcja druga." {
		t.Fatalf("wykaz wersji nie zaczyna się od najnowszej: %+v", wykaz)
	}

	najstarsza := wykaz.Versions[1]
	var przywrocenie shared.WorkspaceInstructionsVersionRestoreResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceInstructionsVersionRestore,
		shared.WorkspaceInstructionsVersionRestoreRequest{
			ProjectId: projektSprawdzianuWorkspace, VersionId: najstarsza.Id,
		}, &przywrocenie)

	if przywrocenie.Instructions.Content != "Instrukcja pierwsza." {
		t.Fatalf("po przywróceniu obowiązuje treść %q", przywrocenie.Instructions.Content)
	}
	zrodlo := napisWorkspace(t, baza,
		`SELECT przywrocono_z FROM wersja_instrukcji_projektu WHERE identyfikator_zewnetrzny = ?`,
		przywrocenie.Version.Id)
	if zrodlo != najstarsza.Id {
		t.Fatalf("wersja z przywrócenia nie wskazuje źródła (%q zamiast %q)", zrodlo, najstarsza.Id)
	}
	obowiazujaca := napisWorkspace(t, baza,
		`SELECT wartosc FROM ustawienie WHERE klucz = ?`, kluczInstrukcjiProjektu)
	if obowiazujaca != "Instrukcja pierwsza." {
		t.Fatalf("ustawienie instrukcji trzyma %q, a przywrócono „Instrukcja pierwsza.”",
			obowiazujaca)
	}
}

// TestStanProjektuIOdlaczenieEksperta wykazuje skutek dwóch czynności na
// projekcie: stan zapisuje się w kolumnie projektu, a odłączenie eksperta
// naprawdę zdejmuje wiersz przypisania.
func TestStanProjektuIOdlaczenieEksperta(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaWorkspace(t, katalog)

	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceProjectStatusSet,
		shared.WorkspaceProjectStatusSetRequest{
			ProjectId: projektSprawdzianuWorkspace,
			Status:    shared.WorkspaceProjectStatusPaused,
		}, nil)

	stan := napisWorkspace(t, baza, `SELECT stan FROM projekt WHERE kod = ?`,
		projektSprawdzianuWorkspace)
	if stan != shared.WorkspaceProjectStatusPaused {
		t.Fatalf("projekt stoi w stanie %q zamiast paused", stan)
	}

	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceAgentAssign,
		shared.WorkspaceAgentAssignRequest{
			ProjectId: projektSprawdzianuWorkspace, AgentId: "analityk",
			DefaultExecutor: wskaznik(true),
		}, nil)
	if przypisania := liczbaWorkspace(t, baza,
		`SELECT count(*) FROM przypisanie_agenta_projektu`); przypisania != 1 {
		t.Fatalf("po przypisaniu w bazie leży %d wierszy", przypisania)
	}

	var odlaczenie shared.WorkspaceAgentUnassignResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceAgentUnassign,
		shared.WorkspaceAgentUnassignRequest{
			ProjectId: projektSprawdzianuWorkspace, AgentId: "analityk",
		}, &odlaczenie)
	if !odlaczenie.Unassigned {
		t.Fatal("odłączenie zameldowało brak przypisania, które przed chwilą powstało")
	}
	if przypisania := liczbaWorkspace(t, baza,
		`SELECT count(*) FROM przypisanie_agenta_projektu`); przypisania != 0 {
		t.Fatalf("po odłączeniu w bazie zostało %d wierszy przypisania", przypisania)
	}

	var powtorne shared.WorkspaceAgentUnassignResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWorkspaceAgentUnassign,
		shared.WorkspaceAgentUnassignRequest{
			ProjectId: projektSprawdzianuWorkspace, AgentId: "analityk",
		}, &powtorne)
	if powtorne.Unassigned {
		t.Fatal("powtórne odłączenie zameldowało zniesienie przypisania, którego nie było")
	}
}
