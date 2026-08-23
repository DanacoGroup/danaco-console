package core

import (
	"database/sql"
	"strings"
	"testing"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Sprawdziany SKUTKU autozapisu i kopii zapasowych.
//
// ── Szkoda, którą ten plik ma wykluczyć w pierwszej kolejności ───────────────
// Wskaźnik „zapisano" pokazany, gdy zapis się NIE UDAŁ, jest najgorszym możliwym
// błędem tego modułu: Operator zamknie okno i straci pracę. Dlatego sprawdzian
// mierzy zapis NIEUDANY, a nie tylko udany — i mierzy go na prawdziwym
// niepowodzeniu zapisu, nie na atrapie.
//
// Niepowodzenie wywołujemy, zabierając rdzeniowi tabelę, w którą zapisuje postać
// dokumentu. To jest niepowodzenie tej samej klasy, co awaria dysku albo
// uszkodzony plik bazy: zapis wchodzi tą samą drogą i pada w tym samym miejscu.
// Atrapa repozytorium mierzyłaby wtedy atrapę.

// autozapisSkutekTresc jest treścią, której zapis samoczynny ma nie zgubić.
const autozapisSkutekTresc = "Pismo w toku redakcji.\nAkapit niezapisany."

// TestAutozapisSkutekNieudanyZapisNiePokazujeZapisano jest sprawdzianem
// uczciwości zapisu: po nieudanym zapisie samoczynnym odpowiedź mówi wprost, że
// się nie udało, wskaźnik NIE stoi na „zapisano", a praca zostaje w kopii.
func TestAutozapisSkutekNieudanyZapisNiePokazujeZapisano(t *testing.T) {
	uprzaz := blokadaZmontuj(t)

	// Autozapis włączony jawnie — jak przełącznikiem na pasku szybkiego dostępu.
	odstep := 60
	var nastawy shared.StudioAutosaveSetResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioAutosaveSet,
		shared.StudioAutosaveSetRequest{
			DocumentId: &uprzaz.dokument, Enabled: true, IntervalSeconds: &odstep,
		}, &nastawy)

	// Zapis udany na wejściu — żeby było widać różnicę między „stało" i „padło".
	tresc := autozapisSkutekTresc
	var udany shared.StudioAutosaveRunResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioAutosaveRun,
		shared.StudioAutosaveRunRequest{DocumentId: uprzaz.dokument, Content: &tresc}, &udany)
	if !udany.Saved {
		t.Fatalf("zapis samoczynny nie doszedł do skutku na zdrowej bazie: %+v", udany.FailureReason)
	}

	// ── Prawdziwa awaria ZAPISU, przy zdrowym ODCZYCIE ───────────────────────
	// Wyzwalacz odrzuca każdy zapis postaci, a odczyt zostawia nietknięty. To jest
	// właśnie ten stan, o który idzie wymaganie: rdzeń CZYTA dokument, przyjmuje
	// pracę, a utrwalić jej nie potrafi. Zabranie całej tabeli mierzyłoby co
	// innego — tam pada już odczyt i nie ma nawet czego zapisać.
	if _, err := uprzaz.baza.Exec(
		`CREATE TRIGGER awaria_zapisu_postaci
		 BEFORE UPDATE ON postac_dokumentu_studio
		 BEGIN SELECT RAISE(ABORT, 'awaria nosnika w sprawdzianie'); END`); err != nil {
		t.Fatalf("nie można wywołać awarii zapisu: %v", err)
	}
	t.Cleanup(func() {
		_, _ = uprzaz.baza.Exec(`DROP TRIGGER IF EXISTS awaria_zapisu_postaci`)
	})

	nowaTresc := autozapisSkutekTresc + "\nAkapit dopisany po awarii."
	var nieudany shared.StudioAutosaveRunResponse
	odpowiedz := wykonajKomende(t, uprzaz.zmontowany, uprzaz.zycie,
		shared.CommandStudioAutosaveRun,
		shared.StudioAutosaveRunRequest{DocumentId: uprzaz.dokument, Content: &nowaTresc})

	// Odmowa też jest dopuszczalna — byle NIE „zapisano". Niedopuszczalne jest
	// jedno: odpowiedź udana mówiąca, że zapis stanął.
	if odpowiedz.Error == nil {
		if err := protocol.LadunekDo(odpowiedz, &nieudany); err != nil {
			t.Fatalf("nieczytelny ładunek odpowiedzi: %v", err)
		}
		if nieudany.Saved {
			t.Fatal("WSKAŹNIK POKAZAŁ „ZAPISANO”, A ZAPIS PADŁ — Operator zamknie " +
				"okno i straci pracę")
		}
		if nieudany.FailureReason == nil || *nieudany.FailureReason == "" {
			t.Error("nieudany zapis samoczynny nie jest NAZWANY — odpowiedź nie mówi, " +
				"co się stało")
		}
		if nieudany.Backup == nil {
			t.Error("nieudany zapis nie oddał kopii zapasowej, w której praca została")
		}
	}

	// ── Pomiar niezależny: praca ZOSTAŁA w kopii ────────────────────────────
	var trescKopii sql.NullString
	var niezapisane int
	err := uprzaz.baza.QueryRow(
		`SELECT k.tresc, k.zmiany_niezapisane FROM kopia_zapasowa_studio k
		 JOIN dokument_studio d ON d.id = k.dokument_id
		 WHERE d.identyfikator_zewnetrzny = ? AND k.zmiany_niezapisane = 1
		 ORDER BY k.id DESC LIMIT 1`, uprzaz.dokument).Scan(&trescKopii, &niezapisane)
	if err != nil {
		t.Fatalf("po nieudanym zapisie nie ma kopii niosącej zmiany niezapisane: %v", err)
	}
	if !strings.Contains(trescKopii.String, "Akapit dopisany po awarii.") {
		t.Errorf("kopia zapasowa nie niesie pracy, której zapis nie utrwalił; treść: %q",
			trescKopii.String)
	}

	// ── Pomiar niezależny: wskaźnik stanu w bazie ───────────────────────────
	var nieudanyWBazie int
	var ostatniZapis sql.NullString
	err = uprzaz.baza.QueryRow(
		`SELECT n.ostatni_zapis_nieudany, n.ostatni_zapis FROM nastawa_pracy_studio n
		 JOIN dokument_studio d ON d.id = n.dokument_id
		 WHERE d.identyfikator_zewnetrzny = ?`, uprzaz.dokument).Scan(&nieudanyWBazie, &ostatniZapis)
	if err != nil {
		t.Fatalf("nie można odczytać wskaźnika stanu autozapisu: %v", err)
	}
	if nieudanyWBazie != 1 {
		t.Error("wskaźnik „ostatni zapis nieudany” nie stoi po awarii zapisu — " +
			"okno pokaże Operatorowi stan zdrowy")
	}
}

// TestAutozapisSkutekWersjaIdzieOsobnymSzeregiem mierzy, że zapisy samoczynne
// NIE zaśmiecają historii Operatora: idą osobnym szeregiem, odróżnialnym
// w wykazie.
func TestAutozapisSkutekWersjaIdzieOsobnymSzeregiem(t *testing.T) {
	uprzaz := blokadaZmontuj(t)

	// Trzy zapisy samoczynne.
	for _, tresc := range []string{"Wersja robocza jeden.", "Wersja robocza dwa.",
		"Wersja robocza trzy."} {

		zapis := tresc
		var wynik shared.StudioAutosaveRunResponse
		wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioAutosaveRun,
			shared.StudioAutosaveRunRequest{DocumentId: uprzaz.dokument, Content: &zapis}, &wynik)
		if wynik.Version == nil {
			t.Fatal("zapis samoczynny nie odłożył wersji")
		}
	}

	// Pomiar niezależny: wersje szeregu autozapisu i szeregu Operatora osobno.
	var autozapisu, operatora int
	err := uprzaz.baza.QueryRow(
		`SELECT
		     SUM(CASE WHEN w.szereg = 'autosave' THEN 1 ELSE 0 END),
		     SUM(CASE WHEN w.szereg = 'operator' THEN 1 ELSE 0 END)
		 FROM wersja_dokumentu_studio w
		 JOIN dokument_studio d ON d.id = w.dokument_id
		 WHERE d.identyfikator_zewnetrzny = ?`, uprzaz.dokument).Scan(&autozapisu, &operatora)
	if err != nil {
		t.Fatalf("nie można policzyć wersji w szeregach: %v", err)
	}
	if autozapisu != 3 {
		t.Errorf("szereg autozapisu niesie %d wersji, a zapisów samoczynnych było 3", autozapisu)
	}
	if operatora != 1 {
		t.Errorf("szereg Operatora niesie %d wersji, a Operator zapisał raz — "+
			"zapisy samoczynne zaśmieciły jego historię", operatora)
	}

	// Wykaz szeregów oddaje oba liczniki i wskazuje wersję założycielską.
	var wykaz shared.StudioVersionSeriesListResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioVersionSeriesList,
		shared.StudioVersionSeriesListRequest{DocumentId: uprzaz.dokument}, &wykaz)
	if wykaz.AutosaveCount != 3 || wykaz.OperatorCount != 1 {
		t.Errorf("wykaz szeregów liczy autozapis=%d operator=%d, a w bazie stoi 3 i 1",
			wykaz.AutosaveCount, wykaz.OperatorCount)
	}

	// Zawężenie do szeregu Operatora nie pokazuje zapisów samoczynnych — na tym
	// stoi przełącznik „pokaż także zapisy samoczynne" w wykazie historii.
	//
	// Mierzone LICZBĄ pozycji, a nie polem szeregu w pozycji: `StudioVersion`
	// w kontrakcie pola `series` nie ma (odcinek kontraktu ma je dołożyć), więc
	// odróżnialność wychodzi dziś zawężeniem wykazu i dwoma licznikami. To jest
	// prawdziwy stan i tak jest zgłoszony — sprawdzian mierzy to, co jest.
	szereg := shared.StudioVersionSeries(shared.StudioVersionSeriesOperator)
	var tylkoOperator shared.StudioVersionSeriesListResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioVersionSeriesList,
		shared.StudioVersionSeriesListRequest{DocumentId: uprzaz.dokument, Series: &szereg},
		&tylkoOperator)
	if len(tylkoOperator.Versions) != 1 {
		t.Errorf("wykaz zawężony do szeregu Operatora niesie %d pozycji, a Operator "+
			"zapisał raz — zapisy samoczynne przeciekły do jego historii",
			len(tylkoOperator.Versions))
	}

	// I odwrotnie: zawężenie do szeregu autozapisu oddaje wyłącznie zapisy samoczynne.
	szeregAutozapisu := shared.StudioVersionSeries(shared.StudioVersionSeriesAutosave)
	var tylkoAutozapis shared.StudioVersionSeriesListResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioVersionSeriesList,
		shared.StudioVersionSeriesListRequest{DocumentId: uprzaz.dokument, Series: &szeregAutozapisu},
		&tylkoAutozapis)
	if len(tylkoAutozapis.Versions) != 3 {
		t.Errorf("wykaz zawężony do szeregu autozapisu niesie %d pozycji, a zapisów "+
			"samoczynnych było 3", len(tylkoAutozapis.Versions))
	}
}

// TestKopiaSkutekPrzywrocenieDoNowegoDokumentu mierzy, że przywrócenie kopii
// zapasowej DO NOWEGO dokumentu nie kasuje tego, co jest — a nowy dokument jest
// osobnym bytem, nie drugim odwołaniem do tego samego.
func TestKopiaSkutekPrzywrocenieDoNowegoDokumentu(t *testing.T) {
	uprzaz := blokadaZmontuj(t)

	// Kopia niosąca treść wcześniejszą.
	trescKopii := "Brzmienie sprzed pomyłki."
	var zalozenie shared.StudioBackupCreateResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioBackupCreate,
		shared.StudioBackupCreateRequest{
			DocumentId: uprzaz.dokument, Content: &trescKopii,
		}, &zalozenie)
	if zalozenie.Backup.Id == "" {
		t.Fatal("kopia zapasowa założona bez identyfikatora")
	}

	// Dokument pierwotny idzie dalej swoją drogą.
	trescBiezaca := "Brzmienie bieżące, którego nie wolno stracić."
	var zapis shared.StudioDocumentSaveResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{DocumentId: uprzaz.dokument, Content: trescBiezaca}, &zapis)

	doNowego := true
	nazwa := "Przywrócone z kopii"
	var przywrocenie shared.StudioBackupRestoreResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioBackupRestore,
		shared.StudioBackupRestoreRequest{
			BackupId: zalozenie.Backup.Id, AsNewDocument: &doNowego, Title: &nazwa,
		}, &przywrocenie)

	if przywrocenie.Document.Id == uprzaz.dokument {
		t.Fatal("przywrócenie DO NOWEGO dokumentu nadpisało dokument pierwotny — " +
			"przywracanie skasowało to, co było")
	}

	// Pomiar niezależny: dwa osobne wiersze, każdy o swojej treści.
	if biezaca := uprzaz.blokadaTrescZBazy(t); !strings.Contains(biezaca, "Brzmienie bieżące") {
		t.Errorf("dokument pierwotny stracił swoją treść: %q", biezaca)
	}
	var trescNowego sql.NullString
	err := uprzaz.baza.QueryRow(
		`SELECT tresc FROM dokument_studio WHERE identyfikator_zewnetrzny = ?`,
		przywrocenie.Document.Id).Scan(&trescNowego)
	if err != nil {
		t.Fatalf("nowy dokument nie stanął w bazie: %v", err)
	}
	if !strings.Contains(trescNowego.String, "Brzmienie sprzed pomyłki.") {
		t.Errorf("nowy dokument nie niesie treści kopii: %q", trescNowego.String)
	}
}

// TestKopiaSkutekWykazZglaszaNiezapisane mierzy przywrócenie po nagłym
// zamknięciu: rdzeń SAM zgłasza, że ma niezapisany dokument, a nie czeka, aż
// Operator się domyśli.
func TestKopiaSkutekWykazZglaszaNiezapisane(t *testing.T) {
	uprzaz := blokadaZmontuj(t)

	// Kopia niosąca zmiany niezapisane powstaje przy zapisie samoczynnym —
	// dokładnie tak, jak powstałaby przed nagłym zamknięciem okna.
	tresc := "Praca, której Operator nie zapisał."
	var zapis shared.StudioAutosaveRunResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioAutosaveRun,
		shared.StudioAutosaveRunRequest{DocumentId: uprzaz.dokument, Content: &tresc}, &zapis)

	// Kopię niosącą zmiany niezapisane zakładamy wprost — zapis udany zeruje ten
	// wskaźnik, a mierzymy stan po nagłym zamknięciu, czyli bez zapisu udanego.
	trescNiezapisana := "Akapit, który nie doszedł na dysk."
	if _, err := uprzaz.baza.Exec(
		`INSERT INTO kopia_zapasowa_studio
		     (identyfikator_zewnetrzny, dokument_id, powod, tresc, rozmiar_bajtow,
		      udalo_sie, zmiany_niezapisane)
		 SELECT 'studio-kop-sprawdzian', d.id, 'event', ?, ?, 1, 1
		 FROM dokument_studio d WHERE d.identyfikator_zewnetrzny = ?`,
		trescNiezapisana, len(trescNiezapisana), uprzaz.dokument); err != nil {
		t.Fatalf("nie można odłożyć kopii niezapisanej: %v", err)
	}

	tylkoNiezapisane := true
	var wykaz shared.StudioBackupListResponse
	wykonajUdana(t, uprzaz.zmontowany, uprzaz.zycie, shared.CommandStudioBackupList,
		shared.StudioBackupListRequest{
			DocumentId: &uprzaz.dokument, UnsavedOnly: &tylkoNiezapisane,
		}, &wykaz)

	if wykaz.UnsavedCount == 0 {
		t.Fatal("wykaz kopii nie zgłasza ani jednej kopii ze zmianami niezapisanymi — " +
			"Studio nie zgłosi „mam niezapisany dokument z godziny X”")
	}
	znaleziona := false
	for _, kopia := range wykaz.Backups {
		if kopia.Id == "studio-kop-sprawdzian" {
			znaleziona = true
			if kopia.UnsavedChanges == nil || !*kopia.UnsavedChanges {
				t.Error("kopia nie jest oznaczona jako niosąca zmiany niezapisane")
			}
		}
	}
	if !znaleziona {
		t.Errorf("kopia niosąca zmiany niezapisane nie wyszła wykazem: %+v", wykaz.Backups)
	}
}
