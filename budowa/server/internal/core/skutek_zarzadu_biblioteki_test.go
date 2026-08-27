package core

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// Skutek zarządu repozytorium: czy za odpowiedzią rodziny `library.*` stoi
// zmieniony stan, a nie sam meldunek.
//
// Wzorzec szkody jest w tym produkcie udokumentowany: moduł Design meldował
// `status: ok` wraz z wykazem zasobów, za którymi nie było ani jednego bajtu.
// Dlatego ani jeden sprawdzian w tym pliku nie kończy się na odpowiedzi komendy.
// Każdy schodzi NIŻEJ niż rdzeń: otwiera plik bazy osobnym połączeniem SQL albo
// czyta bajty z magazynu treści na dysku — i pyta o to samo, co komenda
// zameldowała.
//
// Miara niezależna, nie druga komenda. Gdyby stan czytała inna komenda tego
// samego modułu, obie mogłyby mylić się zgodnie: adapter oddający wykaz z tego
// samego miejsca, w którym zapisał, potwierdziłby sam siebie.

// bazaSprawdzianuBiblioteki otwiera bazę rdzenia osobnym połączeniem — do pomiaru
// niezależnego od modułu.
func bazaSprawdzianuBiblioteki(t *testing.T, katalog string) *sql.DB {
	t.Helper()

	polaczenie, err := sql.Open("sqlite", filepath.Join(katalog, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy do pomiaru: %v", err)
	}
	t.Cleanup(func() { _ = polaczenie.Close() })
	return polaczenie
}

// liczbaWierszyBiblioteki liczy wiersze zapytaniem policzalnym — jedno miejsce dla
// wszystkich pomiarów tego pliku.
func liczbaWierszyBiblioteki(t *testing.T, baza *sql.DB, zapytanie string, argumenty ...any) int {
	t.Helper()

	var liczba int
	if err := baza.QueryRow(zapytanie, argumenty...).Scan(&liczba); err != nil {
		t.Fatalf("pomiar %q nie powiódł się: %v", zapytanie, err)
	}
	return liczba
}

// wartoscTekstowaBiblioteki odczytuje jedną kolumnę tekstową — pomiar wprost z bazy.
func wartoscTekstowaBiblioteki(t *testing.T, baza *sql.DB, zapytanie string, argumenty ...any) string {
	t.Helper()

	var wartosc sql.NullString
	if err := baza.QueryRow(zapytanie, argumenty...).Scan(&wartosc); err != nil {
		t.Fatalf("pomiar %q nie powiódł się: %v", zapytanie, err)
	}
	return wartosc.String
}

// TestOpisZasobuLezyWBaziePoZapisie mierzy `library.metadata.set`: czy opis
// Dublin Core naprawdę wszedł do tabeli opisu, wraz z polem niestandardowym.
//
// Odpowiedź komendy niesie opis po zapisie, więc sama z siebie zawsze wygląda
// pomyślnie. Pomiar idzie po wiersz w `opis_zasobu_biblioteki`.
func TestOpisZasobuLezyWBaziePoZapisie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "umowa.txt", "treść umowy z klientem")

	var zapis shared.LibraryMetadataSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryMetadataSet,
		shared.LibraryMetadataSetRequest{
			FileId: plik.Id,
			Metadata: shared.LibraryMetadata{
				FileId: plik.Id,
				Title:  wskaznik("Umowa ramowa"),
				Rights: wskaznik("wewnętrzne"),
				Custom: []byte(`{"numerSprawy":"2026/17"}`),
			},
		}, &zapis)

	tytul := wartoscTekstowaBiblioteki(t, baza, `SELECT o.tytul FROM opis_zasobu_biblioteki o
	                                   JOIN plik_biblioteki p ON p.id = o.plik_id
	                                   WHERE p.identyfikator_zewnetrzny = ?`, plik.Id)
	if tytul != "Umowa ramowa" {
		t.Errorf("baza niesie tytuł %q, komenda meldowała zapis „Umowa ramowa”", tytul)
	}
	pola := wartoscTekstowaBiblioteki(t, baza, `SELECT o.pola_niestandardowe FROM opis_zasobu_biblioteki o
	                                  JOIN plik_biblioteki p ON p.id = o.plik_id
	                                  WHERE p.identyfikator_zewnetrzny = ?`, plik.Id)
	if !strings.Contains(pola, "2026/17") {
		t.Errorf("pole niestandardowe nie doszło do bazy: %q", pola)
	}

	// Scalanie: drugi zapis bez tytułu ma tytułu nie zdejmować, a prawa wyczyścić
	// wartością pustą. To jest różnica, którą kontrakt opisuje wprost.
	var drugi shared.LibraryMetadataSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryMetadataSet,
		shared.LibraryMetadataSetRequest{
			FileId:   plik.Id,
			Metadata: shared.LibraryMetadata{FileId: plik.Id, Rights: wskaznik("")},
		}, &drugi)

	poScaleniu := wartoscTekstowaBiblioteki(t, baza, `SELECT o.tytul FROM opis_zasobu_biblioteki o
	                                        JOIN plik_biblioteki p ON p.id = o.plik_id
	                                        WHERE p.identyfikator_zewnetrzny = ?`, plik.Id)
	if poScaleniu != "Umowa ramowa" {
		t.Errorf("scalanie zdjęło tytuł, choć żądanie go nie wymieniało: %q", poScaleniu)
	}
	if drugi.Metadata.Rights != nil {
		t.Errorf("prawa miały zostać wyczyszczone wartością pustą, są: %v", *drugi.Metadata.Rights)
	}
}

// TestArchiwizacjaZdejmujeZasobZWykazuIWracaPrzywroceniem mierzy kosz
// repozytorium: stan w bazie oraz obecność zasobu w wykazie domyślnym.
func TestArchiwizacjaZdejmujeZasobZWykazuIWracaPrzywroceniem(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "stare.txt", "materiał do archiwum")

	var zarchiwizowane shared.LibraryFileArchiveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileArchive,
		shared.LibraryFileArchiveRequest{
			FileIds: []string{plik.Id}, Reason: wskaznik("zamknięte przedsięwzięcie"),
		}, &zarchiwizowane)
	if zarchiwizowane.ArchivedCount != 1 {
		t.Fatalf("komenda zameldowała %d zarchiwizowanych, wskazano jeden",
			zarchiwizowane.ArchivedCount)
	}

	stan := wartoscTekstowaBiblioteki(t, baza,
		`SELECT stan FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`, plik.Id)
	if stan != "zarchiwizowany" {
		t.Errorf("baza niesie stan %q, komenda meldowała archiwizację", stan)
	}

	// Wykaz domyślny ma zasób pominąć — inaczej kosz nie jest koszem.
	var wykaz shared.LibraryFileListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileList,
		shared.LibraryFileListRequest{}, &wykaz)
	for _, pozycja := range wykaz.Files {
		if pozycja.Id == plik.Id {
			t.Error("zasób zarchiwizowany dalej stoi w wykazie domyślnym")
		}
	}

	// Wpis dziennika audytu powstaje razem z czynnością i niesie jej powód.
	wpisy := liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM wpis_audytu_biblioteki WHERE plik_kod = ? AND czynnosc = 'archiwizacja'`,
		plik.Id)
	if wpisy == 0 {
		t.Error("dziennik audytu nie odnotował archiwizacji")
	}

	var przywrocone shared.LibraryFileRestoreResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileRestore,
		shared.LibraryFileRestoreRequest{FileIds: []string{plik.Id}}, &przywrocone)
	if przywrocone.RestoredCount != 1 {
		t.Fatalf("przywrócenie zameldowało %d zasobów", przywrocone.RestoredCount)
	}
	if stan := wartoscTekstowaBiblioteki(t, baza,
		`SELECT stan FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`, plik.Id); stan != "aktywny" {
		t.Errorf("po przywróceniu baza niesie stan %q", stan)
	}
}

// TestUsuniecieTrwaleZdejmujeWierszWersjeIWpisIndeksu mierzy jedyną drogę utraty
// danych modułu: po niej w bazie nie ma zostać ani zasób, ani jego wersje.
func TestUsuniecieTrwaleZdejmujeWierszWersjeIWpisIndeksu(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "do-usuniecia.txt", "treść skazana")

	// Bez potwierdzenia komenda ma odmówić — i zasób ma zostać nietknięty.
	odmowa := wykonajKomende(t, zmontowany, zycie, shared.CommandLibraryFileDelete,
		shared.LibraryFileDeleteRequest{FileIds: []string{plik.Id}, Confirm: false})
	if odmowa.Error == nil {
		t.Error("usunięcie bez potwierdzenia przeszło — kontrakt żąda potwierdzenia")
	}
	if liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`, plik.Id) != 1 {
		t.Fatal("zasób zniknął mimo odmowy usunięcia")
	}

	var usuniete shared.LibraryFileDeleteResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileDelete,
		shared.LibraryFileDeleteRequest{FileIds: []string{plik.Id}, Confirm: true}, &usuniete)
	if usuniete.DeletedCount != 1 || usuniete.DeletedVersions < 1 {
		t.Errorf("komenda zameldowała %d zasobów i %d wersji",
			usuniete.DeletedCount, usuniete.DeletedVersions)
	}
	if liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`, plik.Id) != 0 {
		t.Error("wiersz zasobu przeżył usunięcie trwałe")
	}
	if liczbaWierszyBiblioteki(t, baza, `SELECT COUNT(*) FROM wersja_pliku_biblioteki`) != 0 {
		t.Error("wersje przeżyły usunięcie zasobu — kaskada nie zadziałała")
	}
	// Wpis usunięcia zostaje: dziennik ma przeżyć zasób, o którym mówi.
	if liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM wpis_audytu_biblioteki WHERE plik_kod = ? AND czynnosc = 'usuniecie'`,
		plik.Id) == 0 {
		t.Error("dziennik audytu nie odnotował usunięcia trwałego")
	}
}

// TestWeryfikacjaIntegralnosciWykrywaUszkodzonaTresc mierzy `library.fixity.check`
// tam, gdzie jest miarodajna: przy treści zmienionej na dysku pod plecami
// repozytorium.
//
// To jest sprawdzian, którego okno nigdy nie zrobi — bajtów zasobu nie oddaje
// żadna komenda kontraktu.
func TestWeryfikacjaIntegralnosciWykrywaUszkodzonaTresc(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	zdrowy := wgrajPlikTekstowy(t, zmontowany, zycie, "zdrowy.txt", "treść nietknięta")
	uszkodzony := wgrajPlikTekstowy(t, zmontowany, zycie, "uszkodzony.txt", "treść przed uszkodzeniem")

	odwolanie := wartoscTekstowaBiblioteki(t, baza,
		`SELECT tresc_odwolanie FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`,
		uszkodzony.Id)
	if odwolanie == "" {
		t.Fatal("zasób nie ma odwołania do treści — nie ma czego uszkodzić")
	}
	if err := os.WriteFile(odwolanie, []byte("treść podmieniona poza rdzeniem"), 0o600); err != nil {
		t.Fatalf("nie można podmienić treści na dysku: %v", err)
	}

	var wynik shared.LibraryFixityCheckResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFixityCheck,
		shared.LibraryFixityCheckRequest{}, &wynik)

	if wynik.CheckedCount != 2 {
		t.Fatalf("sprawdzono %d zasobów, w repozytorium są dwa", wynik.CheckedCount)
	}
	if wynik.MismatchedCount != 1 {
		t.Errorf("niezgodnych: %d — jeden zasób został podmieniony na dysku",
			wynik.MismatchedCount)
	}
	for _, pozycja := range wynik.Results {
		switch pozycja.FileId {
		case zdrowy.Id:
			if !pozycja.Matched {
				t.Error("zasób nietknięty wyszedł jako niezgodny")
			}
		case uszkodzony.Id:
			if pozycja.Matched {
				t.Error("zasób podmieniony na dysku wyszedł jako zgodny — suma nie jest liczona z bajtów")
			}
			suma := sha256.Sum256([]byte("treść podmieniona poza rdzeniem"))
			if pozycja.ActualChecksum == nil || *pozycja.ActualChecksum != hex.EncodeToString(suma[:]) {
				t.Errorf("suma wyliczona %v nie jest sumą bajtów leżących na dysku",
					pozycja.ActualChecksum)
			}
		}
	}
}

// TestPaczkaMigracyjnaNiesieBajtyZasobow mierzy `library.package.export`: czy
// pod odłożoną paczką leżą bajty, czy sam meldunek o pozycjach.
func TestPaczkaMigracyjnaNiesieBajtyZasobow(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	pierwszy := wgrajPlikTekstowy(t, zmontowany, zycie, "raport.txt", "treść raportu kwartalnego")
	drugi := wgrajPlikTekstowy(t, zmontowany, zycie, "notatka.txt", "treść notatki roboczej")

	var paczka shared.LibraryPackageExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryPackageExport,
		shared.LibraryPackageExportRequest{
			Kind:    shared.LibraryPackageKindMigration,
			FileIds: []string{pierwszy.Id, drugi.Id},
		}, &paczka)

	if paczka.SizeBytes == 0 || paczka.Entries < 4 {
		t.Fatalf("paczka niesie %d bajtów i %d pozycji — zbiór miał dwa zasoby, "+
			"indeks i manifest", paczka.SizeBytes, paczka.Entries)
	}
	if paczka.ManifestChecksum == "" {
		t.Error("paczka bez sumy manifestu — po przeniesieniu nie da się jej sprawdzić")
	}

	odwolanie := wartoscTekstowaBiblioteki(t, baza,
		`SELECT tresc_odwolanie FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`,
		paczka.FileId)
	if odwolanie == "" {
		t.Fatal("zasób paczki nie ma odwołania do treści — za meldunkiem nie ma bajtów")
	}
	bajty, err := os.ReadFile(odwolanie)
	if err != nil {
		t.Fatalf("nie można odczytać paczki z magazynu: %v", err)
	}
	if int64(len(bajty)) != paczka.SizeBytes {
		t.Errorf("paczka na dysku ma %d bajtów, odpowiedź mówiła o %d",
			len(bajty), paczka.SizeBytes)
	}
	// Treść obu zasobów ma być w archiwum — pozycja bez bajtów byłaby dokładnie
	// tą szkodą, przed którą stoi ten plik.
	if !strings.Contains(string(bajty), "raport") || len(bajty) < 200 {
		t.Errorf("archiwum nie niesie nazw zasobów albo jest puste (%d bajtów)", len(bajty))
	}
}

// TestSkanowanieDuplikatowLaczyZasobyOTejSamejTresci mierzy rozpoznanie dokładne
// i scalenie: po scaleniu zasób wchłonięty ma być w archiwum, a etykiety mają
// przejść na zasób, który został.
func TestSkanowanieDuplikatowLaczyZasobyOTejSamejTresci(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	tresc := "dokładnie ta sama treść w dwóch plikach"
	docelowy := wgrajPlikTekstowy(t, zmontowany, zycie, "pierwszy.txt", tresc)
	kopia := wgrajPlikTekstowy(t, zmontowany, zycie, "kopia.txt", tresc)

	var oznaczona shared.LibraryTagSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryTagSet,
		shared.LibraryTagSetRequest{FileId: kopia.Id, Tags: []string{"rynek-x"}}, &oznaczona)

	var skan shared.LibraryDuplicateScanResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryDuplicateScan,
		shared.LibraryDuplicateScanRequest{}, &skan)
	if skan.Total != 1 || len(skan.Groups) != 1 || len(skan.Groups[0].FileIds) != 2 {
		t.Fatalf("skan oddał %d grup, oczekiwano jednej dwuelementowej", skan.Total)
	}

	var scalenie shared.LibraryDuplicateMergeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryDuplicateMerge,
		shared.LibraryDuplicateMergeRequest{
			TargetFileId: docelowy.Id, SourceFileIds: []string{kopia.Id},
		}, &scalenie)

	if scalenie.MergedCount != 1 {
		t.Errorf("scalono %d zasobów, wskazano jeden", scalenie.MergedCount)
	}
	if stan := wartoscTekstowaBiblioteki(t, baza,
		`SELECT stan FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`,
		kopia.Id); stan != "zarchiwizowany" {
		t.Errorf("zasób wchłonięty ma stan %q — miał trafić do archiwum, nie zniknąć", stan)
	}
	przeniesiona := liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM etykieta_pliku_biblioteki e
		 JOIN plik_biblioteki p ON p.id = e.plik_id
		 WHERE p.identyfikator_zewnetrzny = ? AND e.etykieta = 'rynek-x'`, docelowy.Id)
	if przeniesiona != 1 {
		t.Error("etykieta zasobu wchłoniętego nie przeszła na zasób docelowy")
	}
}

// TestSlownikEtykietZmieniaNazweWCalymRepozytorium mierzy `library.tag.update`:
// zmiana nazwy ma przejść po zasobach, nie założyć drugiej etykiety obok.
func TestSlownikEtykietZmieniaNazweWCalymRepozytorium(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	pierwszy := wgrajPlikTekstowy(t, zmontowany, zycie, "a.txt", "treść a")
	drugi := wgrajPlikTekstowy(t, zmontowany, zycie, "b.txt", "treść b")
	for _, kod := range []string{pierwszy.Id, drugi.Id} {
		var wynik shared.LibraryTagSetResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryTagSet,
			shared.LibraryTagSetRequest{FileId: kod, Tags: []string{"klient-x"}}, &wynik)
	}

	var zmiana shared.LibraryTagUpdateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryTagUpdate,
		shared.LibraryTagUpdateRequest{
			Name: "klient-x", NewName: wskaznik("klient-iks"), Color: wskaznik("zeton-zloty"),
		}, &zmiana)

	if zmiana.AffectedFiles != 2 {
		t.Errorf("zmiana dotknęła %d zasobów, etykietę nosiły dwa", zmiana.AffectedFiles)
	}
	stare := liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM etykieta_pliku_biblioteki WHERE etykieta = 'klient-x'`)
	nowe := liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM etykieta_pliku_biblioteki WHERE etykieta = 'klient-iks'`)
	if stare != 0 || nowe != 2 {
		t.Errorf("po zmianie nazwy w bazie stoi %d starych i %d nowych przypisań", stare, nowe)
	}
	if barwa := wartoscTekstowaBiblioteki(t, baza,
		`SELECT barwa FROM etykieta_slownika_biblioteki WHERE nazwa = 'klient-iks'`); barwa != "zeton-zloty" {
		t.Errorf("słownik niesie barwę %q", barwa)
	}

	// Usunięcie etykiety używanej bez potwierdzenia jest odmową — i nic wtedy
	// z zasobów nie schodzi.
	odmowa := wykonajKomende(t, zmontowany, zycie, shared.CommandLibraryTagRemove,
		shared.LibraryTagRemoveRequest{Name: "klient-iks"})
	if odmowa.Error == nil {
		t.Error("usunięcie etykiety używanej przeszło bez potwierdzenia")
	}
	if liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM etykieta_pliku_biblioteki WHERE etykieta = 'klient-iks'`) != 2 {
		t.Error("odmowa usunięcia zdjęła etykietę mimo wszystko")
	}
}

// TestRegulaKolekcjiPrzypisujeZasobySpelniajaceWarunek mierzy `library.rule.set`
// z przeliczeniem: po przebiegu zasoby mają stać w kolekcji, a nie tylko być
// policzone w odpowiedzi.
func TestRegulaKolekcjiPrzypisujeZasobySpelniajaceWarunek(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	pasujacy := wgrajPlikTekstowy(t, zmontowany, zycie, "raport-rynkowy.txt", "treść raportu")
	obcy := wgrajPlikTekstowy(t, zmontowany, zycie, "faktura.txt", "treść faktury")
	var oznaczony shared.LibraryTagSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryTagSet,
		shared.LibraryTagSetRequest{FileId: pasujacy.Id, Tags: []string{"rynek-x"}}, &oznaczony)

	var kolekcja shared.LibraryCollectionCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryCollectionCreate,
		shared.LibraryCollectionCreateRequest{Name: "Rynek X"}, &kolekcja)

	var regula shared.LibraryRuleSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryRuleSet,
		shared.LibraryRuleSetRequest{
			Rule: shared.LibraryRule{
				Kind: shared.LibraryRuleKindCollection, Name: "materiały rynku X",
				Condition:          []byte(`{"tags":["rynek-x"]}`),
				TargetCollectionId: wskaznik(kolekcja.CollectionId),
				Enabled:            true,
			},
			RunNow: wskaznik(true),
		}, &regula)

	if regula.MatchedFiles == nil || *regula.MatchedFiles != 1 {
		t.Fatalf("przeliczenie zameldowało %v trafień, warunek spełnia jeden zasób",
			regula.MatchedFiles)
	}
	przypisane := liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM przypisanie_kolekcji_biblioteki pk
		 JOIN kolekcja_biblioteki k ON k.id = pk.kolekcja_id
		 JOIN plik_biblioteki p ON p.id = pk.plik_id
		 WHERE k.identyfikator_zewnetrzny = ? AND p.identyfikator_zewnetrzny = ?`,
		kolekcja.CollectionId, pasujacy.Id)
	if przypisane != 1 {
		t.Error("zasób spełniający warunek nie stoi w kolekcji — przeliczenie było meldunkiem")
	}
	obcePrzypisanie := liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM przypisanie_kolekcji_biblioteki pk
		 JOIN plik_biblioteki p ON p.id = pk.plik_id
		 WHERE p.identyfikator_zewnetrzny = ?`, obcy.Id)
	if obcePrzypisanie != 0 {
		t.Error("zasób niespełniający warunku został przypisany")
	}

	// Usunięcie reguły ze zdjęciem zasobów zdejmuje wyłącznie przypisania reguły.
	var usuniecie shared.LibraryRuleRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryRuleRemove,
		shared.LibraryRuleRemoveRequest{RuleId: regula.Rule.Id, DetachFiles: wskaznik(true)},
		&usuniecie)
	if !usuniecie.Removed || usuniecie.DetachedFiles != 1 {
		t.Errorf("usunięcie reguły: zdjęto %d zasobów, usunięto=%v",
			usuniecie.DetachedFiles, usuniecie.Removed)
	}
	if liczbaWierszyBiblioteki(t, baza, `SELECT COUNT(*) FROM przypisanie_kolekcji_biblioteki`) != 0 {
		t.Error("przypisania reguły przeżyły jej usunięcie ze zdjęciem zasobów")
	}
}

// TestPulpitStanuLiczyCalyZbiorNieProbke mierzy `library.stats.get` wobec
// liczb policzonych osobnym zapytaniem SQL.
func TestPulpitStanuLiczyCalyZbiorNieProbke(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	wgrajPlikTekstowy(t, zmontowany, zycie, "jeden.txt", "treść jeden")
	wgrajPlikTekstowy(t, zmontowany, zycie, "dwa.txt", "treść dwa")
	trzeci := wgrajPlikTekstowy(t, zmontowany, zycie, "trzy.txt", "treść trzy")

	var zarchiwizowany shared.LibraryFileArchiveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileArchive,
		shared.LibraryFileArchiveRequest{FileIds: []string{trzeci.Id}}, &zarchiwizowany)

	var pulpit shared.LibraryStatsGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryStatsGet,
		shared.LibraryStatsGetRequest{}, &pulpit)

	czynne := liczbaWierszyBiblioteki(t, baza, `SELECT COUNT(*) FROM plik_biblioteki WHERE stan = 'aktywny'`)
	archiwalne := liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM plik_biblioteki WHERE stan = 'zarchiwizowany'`)
	if pulpit.Stats.FileCount != czynne || pulpit.Stats.ArchivedCount != archiwalne {
		t.Errorf("pulpit mówi o %d czynnych i %d archiwalnych, baza ma %d i %d",
			pulpit.Stats.FileCount, pulpit.Stats.ArchivedCount, czynne, archiwalne)
	}
	if pulpit.Stats.OrphanCount != czynne {
		t.Errorf("zasoby bez etykiety i kolekcji: pulpit %d, w bazie %d",
			pulpit.Stats.OrphanCount, czynne)
	}
	if pulpit.Stats.TotalBytes == 0 {
		t.Error("pulpit oddał zerowy rozmiar zbioru, choć zasoby mają treść")
	}
}

// TestPulpitStanuNaPustymRepozytoriumOddajeZera mierzy pulpit na rdzeniu świeżo
// założonym — czyli w stanie, w którym zastaje go Operator otwierający moduł
// Library po raz pierwszy.
//
// Zbiór pusty jest tu przypadkiem granicznym agregatów: sumowanie po zbiorze
// pustym daje w SQL wartość pustą, a pulpit ma pola liczbowe bez stanu pustego.
// Zbiór pusty ma dać ZERA, nie odmowę — zero zasobów jest wynikiem, a nie
// niepowodzeniem odczytu.
func TestPulpitStanuNaPustymRepozytoriumOddajeZera(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var pulpit shared.LibraryStatsGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryStatsGet,
		shared.LibraryStatsGetRequest{}, &pulpit)

	if pulpit.Stats.FileCount != 0 || pulpit.Stats.ArchivedCount != 0 {
		t.Errorf("puste repozytorium: pulpit mówi o %d zasobach czynnych i %d archiwalnych",
			pulpit.Stats.FileCount, pulpit.Stats.ArchivedCount)
	}
	if pulpit.Stats.TotalBytes != 0 {
		t.Errorf("puste repozytorium ma rozmiar %d bajtów", pulpit.Stats.TotalBytes)
	}
	if pulpit.Stats.MissingChecksumCount != 0 || pulpit.Stats.DuplicateCount != 0 ||
		pulpit.Stats.OrphanCount != 0 {
		t.Errorf("puste repozytorium: bez sumy kontrolnej %d, duplikatów %d, osieroconych %d",
			pulpit.Stats.MissingChecksumCount, pulpit.Stats.DuplicateCount,
			pulpit.Stats.OrphanCount)
	}
}

// TestKlasyfikacjaWytwarzaSugestieAPrzyjecieJeWykonuje mierzy pełną drogę
// sugestii: wytworzenie, zapis w bazie i skutek przyjęcia.
//
// Sam meldunek „przyjęto" nie wystarcza — przyjęcie ma nadać etykietę, więc
// pomiar schodzi po wiersz przypisania etykiety.
func TestKlasyfikacjaWytwarzaSugestieAPrzyjecieJeWykonuje(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "umowa-najmu.txt",
		"umowa najmu umowa najmu najem lokalu przy ulicy Polnej")

	var klasyfikacja shared.LibraryClassifyRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryClassifyRun,
		shared.LibraryClassifyRunRequest{FileIds: []string{plik.Id}}, &klasyfikacja)

	if klasyfikacja.ProcessedCount != 1 || len(klasyfikacja.Suggestions) == 0 {
		t.Fatalf("klasyfikacja przetworzyła %d zasobów i oddała %d sugestii",
			klasyfikacja.ProcessedCount, len(klasyfikacja.Suggestions))
	}
	if klasyfikacja.AppliedCount != 0 {
		t.Error("klasyfikacja zapisała sugestie bez pytania — domyślnie ma sugerować")
	}
	oczekujace := liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM sugestia_biblioteki WHERE plik_kod = ? AND stan = 'oczekujaca'`,
		plik.Id)
	if oczekujace != len(klasyfikacja.Suggestions) {
		t.Errorf("w bazie stoi %d sugestii oczekujących, odpowiedź niosła %d",
			oczekujace, len(klasyfikacja.Suggestions))
	}

	var etykietowa string
	for _, sugestia := range klasyfikacja.Suggestions {
		if sugestia.Kind == shared.LibrarySuggestionKindTag {
			etykietowa = sugestia.Id
			break
		}
	}
	if etykietowa == "" {
		t.Fatal("klasyfikacja nie zaproponowała ani jednej etykiety")
	}

	var decyzja shared.LibrarySuggestionApplyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibrarySuggestionApply,
		shared.LibrarySuggestionApplyRequest{SuggestionIds: []string{etykietowa}, Accept: true},
		&decyzja)
	if decyzja.AppliedCount != 1 {
		t.Fatalf("przyjęto %d sugestii", decyzja.AppliedCount)
	}
	nadane := liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM etykieta_pliku_biblioteki e
		 JOIN plik_biblioteki p ON p.id = e.plik_id
		 WHERE p.identyfikator_zewnetrzny = ?`, plik.Id)
	if nadane == 0 {
		t.Error("przyjęcie sugestii etykiety nie nadało etykiety — meldunek bez skutku")
	}
}

// TestPorownanieWskazujeRoznicaMiedzyWersjami mierzy `library.diff.compare`
// wobec treści, którą sprawdzian sam włożył do repozytorium.
func TestPorownanieWskazujeRoznicaMiedzyWersjami(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "tekst.txt",
		"wiersz pierwszy\nwiersz drugi\nwiersz trzeci")
	pierwszaWersja := *plik.VersionId

	var dolozona shared.LibraryVersionAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryVersionAdd,
		shared.LibraryVersionAddRequest{
			FileId:        plik.Id,
			ContentBase64: wskaznik(wBase64([]byte("wiersz pierwszy\nwiersz drugi zmieniony\nwiersz trzeci"))),
		}, &dolozona)

	var porownanie shared.LibraryDiffCompareResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryDiffCompare,
		shared.LibraryDiffCompareRequest{
			LeftFileId:     plik.Id,
			LeftVersionId:  wskaznik(pierwszaWersja),
			RightVersionId: wskaznik(dolozona.Version.Id),
		}, &porownanie)

	if porownanie.Diff.Identical {
		t.Fatal("porównanie orzekło zgodność treści, które się różnią")
	}
	if len(porownanie.Diff.Hunks) != 1 {
		t.Fatalf("porównanie oddało %d różnic, zmieniono jeden wiersz",
			len(porownanie.Diff.Hunks))
	}
	roznica := porownanie.Diff.Hunks[0]
	if roznica.Kind != shared.LibraryDiffKindChanged {
		t.Errorf("różnica ma rodzaj %q, wiersz został zmieniony", roznica.Kind)
	}
	if !strings.Contains(roznica.Text, "zmieniony") {
		t.Errorf("treść różnicy nie niesie zmienionego wiersza: %q", roznica.Text)
	}
}

// TestUdostepnienieMaTokenLosowyIOdwolanieDziala mierzy rodzinę udostępnień:
// token ma być niepowtarzalny, a odwołanie ma zdjąć odnośnik z wykazu czynnych.
func TestUdostepnienieMaTokenLosowyIOdwolanieDziala(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "do-wyslania.txt", "treść dla klienta")

	var pierwsze, drugie shared.LibraryShareCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryShareCreate,
		shared.LibraryShareCreateRequest{
			Scope: shared.LibraryShareScopeFile, TargetId: plik.Id,
		}, &pierwsze)
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryShareCreate,
		shared.LibraryShareCreateRequest{
			Scope: shared.LibraryShareScopeFile, TargetId: plik.Id,
			ExpiresInSeconds: wskaznik(3600),
		}, &drugie)

	if pierwsze.Share.Token == drugie.Share.Token {
		t.Error("dwa udostępnienia dostały ten sam token — token nie jest losowy")
	}
	if len(pierwsze.Share.Token) < 32 {
		t.Errorf("token ma %d znaków — za mało, żeby nie dał się zgadnąć",
			len(pierwsze.Share.Token))
	}
	if !strings.Contains(pierwsze.Url, pierwsze.Share.Token) {
		t.Error("adres udostępnienia nie niesie tokenu")
	}
	if drugie.Share.ExpiresAt == nil {
		t.Error("udostępnienie z terminem nie niesie chwili wygaśnięcia")
	}

	var odwolanie shared.LibraryShareRevokeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryShareRevoke,
		shared.LibraryShareRevokeRequest{ShareId: pierwsze.Share.Id}, &odwolanie)
	if !odwolanie.Revoked {
		t.Fatal("odwołanie nie doszło do skutku")
	}
	if liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM udostepnienie_biblioteki
		 WHERE identyfikator_zewnetrzny = ? AND odwolano IS NOT NULL`,
		pierwsze.Share.Id) != 1 {
		t.Error("odwołanie nie zostawiło znacznika czasu — wpis miał zostać")
	}

	var czynne shared.LibraryShareListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryShareList,
		shared.LibraryShareListRequest{ActiveOnly: wskaznik(true)}, &czynne)
	for _, udostepnienie := range czynne.Shares {
		if udostepnienie.Id == pierwsze.Share.Id {
			t.Error("udostępnienie odwołane stoi dalej w wykazie czynnych")
		}
	}
}

// TestNormalizacjaNazwProbnaNieZapisujeAWlasciwaZapisuje mierzy różnicę między
// przebiegiem próbnym a zapisem — obie drogi jednej komendy.
func TestNormalizacjaNazwProbnaNieZapisujeAWlasciwaZapisuje(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "Raport Roczny ĄĘŁ.txt", "treść raportu")

	var probny shared.LibraryNameNormalizeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryNameNormalize,
		shared.LibraryNameNormalizeRequest{
			FileIds: []string{plik.Id}, DryRun: wskaznik(true),
		}, &probny)

	if probny.ChangedCount != 1 || len(probny.Results) != 1 || probny.Results[0].Applied {
		t.Fatalf("przebieg próbny: zmian %d, zapisano=%v",
			probny.ChangedCount, len(probny.Results) > 0 && probny.Results[0].Applied)
	}
	if nazwa := wartoscTekstowaBiblioteki(t, baza,
		`SELECT nazwa FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`,
		plik.Id); nazwa != "Raport Roczny ĄĘŁ.txt" {
		t.Errorf("przebieg próbny zmienił nazwę w bazie na %q", nazwa)
	}

	var wlasciwy shared.LibraryNameNormalizeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryNameNormalize,
		shared.LibraryNameNormalizeRequest{FileIds: []string{plik.Id}}, &wlasciwy)

	if wlasciwy.ChangedCount != 1 || !wlasciwy.Results[0].Applied {
		t.Fatal("przebieg właściwy nie zapisał zmiany")
	}
	nazwa := wartoscTekstowaBiblioteki(t, baza,
		`SELECT nazwa FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`, plik.Id)
	if nazwa != wlasciwy.Results[0].NewName {
		t.Errorf("baza niesie nazwę %q, odpowiedź mówiła o %q", nazwa, wlasciwy.Results[0].NewName)
	}
	if strings.ContainsAny(nazwa, "ĄĘŁ ") {
		t.Errorf("nazwa po normalizacji dalej niesie znaki diakrytyczne albo spacje: %q", nazwa)
	}
}

// TestUtrwalenieBagItSkladaPakietZManifestem mierzy `library.preservation.run`:
// pod wynikiem ma leżeć pakiet, którego manifest zgadza się z jego zawartością.
func TestUtrwalenieBagItSkladaPakietZManifestem(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "akta.txt", "treść akt sprawy")

	var utrwalenie shared.LibraryPreservationRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryPreservationRun,
		shared.LibraryPreservationRunRequest{
			FileIds: []string{plik.Id}, Kind: shared.LibraryPreservationKindBagit,
		}, &utrwalenie)

	if utrwalenie.ValidCount != 1 || len(utrwalenie.Results) != 1 {
		t.Fatalf("utrwalenie: poprawnych %d, wyników %d",
			utrwalenie.ValidCount, len(utrwalenie.Results))
	}
	wynik := utrwalenie.Results[0]
	if wynik.ProducedFileId == nil {
		t.Fatal("utrwalenie nie wskazało zasobu wytworzonego")
	}
	if !strings.Contains(wynik.Report, "manifest") {
		t.Errorf("zapis walidacji nie mówi, co sprawdzono: %q", wynik.Report)
	}

	odwolanie := wartoscTekstowaBiblioteki(t, baza,
		`SELECT tresc_odwolanie FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`,
		*wynik.ProducedFileId)
	if odwolanie == "" {
		t.Fatal("pakiet nie ma odwołania do bajtów")
	}
	bajty, err := os.ReadFile(odwolanie)
	if err != nil {
		t.Fatalf("nie można odczytać pakietu: %v", err)
	}
	if !strings.Contains(string(bajty), "bagit.txt") {
		t.Error("pakiet nie niesie pliku bagit.txt — to nie jest pakiet BagIt")
	}
	// Ślad utrwalenia zostaje w bazie razem z werdyktem walidacji.
	if liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM zadanie_utrwalenia_biblioteki WHERE plik_kod = ? AND poprawne = 1`,
		plik.Id) != 1 {
		t.Error("baza nie niesie śladu utrwalenia")
	}
}

// TestDziennikAudytuNieKlamieOCzynnosciach mierzy `library.audit.list` wobec
// czynności, które sprawdzian wykonał sam.
func TestDziennikAudytuNieKlamieOCzynnosciach(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "sledzony.txt", "treść śledzona")

	var przeniesienie shared.LibraryFileMoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileMove,
		shared.LibraryFileMoveRequest{FileIds: []string{plik.Id}, Path: "/Klienci/Klient X/"},
		&przeniesienie)
	if przeniesienie.MovedCount != 1 || len(przeniesienie.Files) != 1 {
		t.Fatalf("przeniesienie objęło %d zasobów", przeniesienie.MovedCount)
	}
	if przeniesienie.Files[0].Path == nil || *przeniesienie.Files[0].Path != "Klienci/Klient X" {
		t.Errorf("zasób po przeniesieniu niesie ścieżkę %v", przeniesienie.Files[0].Path)
	}

	var dziennik shared.LibraryAuditListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryAuditList,
		shared.LibraryAuditListRequest{FileId: wskaznik(plik.Id)}, &dziennik)

	if dziennik.Total == 0 {
		t.Fatal("dziennik nie odnotował ani jednej czynności na zasobie")
	}
	zmiany := 0
	for _, wpis := range dziennik.Entries {
		if wpis.Action == shared.LibraryAuditActionChange {
			zmiany++
		}
		if wpis.Actor == "" {
			t.Error("wpis dziennika bez sprawcy")
		}
	}
	if zmiany == 0 {
		t.Error("dziennik nie odnotował przeniesienia jako zmiany")
	}

	// Zawężenie czynnością ma odsiewać: wykaz samych archiwizacji jest pusty,
	// bo zasobu nie archiwizowano.
	var archiwizacje shared.LibraryAuditListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryAuditList,
		shared.LibraryAuditListRequest{
			FileId:  wskaznik(plik.Id),
			Actions: []shared.LibraryAuditAction{shared.LibraryAuditActionArchive},
		}, &archiwizacje)
	if archiwizacje.Total != 0 {
		t.Errorf("zawężenie do archiwizacji oddało %d wpisów, zasobu nie archiwizowano",
			archiwizacje.Total)
	}
}

// TestSchematMetadanychZaklada I ZdejmujePolaNiestandardowe mierzy
// `library.schema.set`: definicja ma powstać i zniknąć, a wartości przy
// zasobach mają zostać.
func TestSchematMetadanychZakladaIZdejmujePolaNiestandardowe(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "sprawa.txt", "treść sprawy")

	var zalozenie shared.LibrarySchemaSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibrarySchemaSet,
		shared.LibrarySchemaSetRequest{Field: shared.LibraryFieldDefinition{
			Code: "numerSprawy", Label: "Numer sprawy", Kind: shared.LibraryFieldKindText,
		}}, &zalozenie)
	if zalozenie.Field.Code != "numerSprawy" {
		t.Fatalf("zapis pola oddał kod %q", zalozenie.Field.Code)
	}
	if liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM pole_schematu_biblioteki WHERE kod = 'numerSprawy'`) != 1 {
		t.Fatal("definicja pola nie doszła do bazy")
	}

	var opis shared.LibraryMetadataSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryMetadataSet,
		shared.LibraryMetadataSetRequest{
			FileId: plik.Id,
			Metadata: shared.LibraryMetadata{
				FileId: plik.Id, Custom: []byte(`{"numerSprawy":"K/2026/8"}`),
			},
		}, &opis)

	var zdjecie shared.LibrarySchemaSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibrarySchemaSet,
		shared.LibrarySchemaSetRequest{
			Field:  shared.LibraryFieldDefinition{Code: "numerSprawy", Label: "Numer sprawy"},
			Remove: wskaznik(true),
		}, &zdjecie)

	if zdjecie.AffectedFiles != 1 {
		t.Errorf("zdjęcie definicji zameldowało %d zasobów z wartością pola",
			zdjecie.AffectedFiles)
	}
	if liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM pole_schematu_biblioteki WHERE kod = 'numerSprawy'`) != 0 {
		t.Error("definicja pola przeżyła zdjęcie")
	}
	// Wartość przy zasobie zostaje — kontrakt mówi wprost, że wraca po
	// ponownym założeniu pola.
	wartosci := wartoscTekstowaBiblioteki(t, baza, `SELECT o.pola_niestandardowe FROM opis_zasobu_biblioteki o
	                                      JOIN plik_biblioteki p ON p.id = o.plik_id
	                                      WHERE p.identyfikator_zewnetrzny = ?`, plik.Id)
	if !strings.Contains(wartosci, "K/2026/8") {
		t.Errorf("zdjęcie definicji zabrało wartość zapisaną przy zasobie: %q", wartosci)
	}
}

// TestTezaurusWywoziPojeciaWrazZRelacjami mierzy wywóz SKOS: treść ma nieść
// pojęcia i krawędzie zapisane komendą relacji.
func TestTezaurusWywoziPojeciaWrazZRelacjami(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var relacja shared.LibraryThesaurusRelateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryThesaurusRelate,
		shared.LibraryThesaurusRelateRequest{
			SourceName: "rynek-x", TargetName: "rynki",
			Relation: shared.LibraryThesaurusRelationBroader,
		}, &relacja)
	if !relacja.RelationSet {
		t.Fatal("relacja tezaurusa nie stanęła")
	}

	var wywoz shared.LibraryThesaurusExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryThesaurusExport,
		shared.LibraryThesaurusExportRequest{}, &wywoz)

	if wywoz.ConceptCount != 2 || wywoz.RelationCount != 1 {
		t.Errorf("wywóz niesie %d pojęć i %d relacji, zapisano dwa pojęcia i jedną relację",
			wywoz.ConceptCount, wywoz.RelationCount)
	}
	if !strings.Contains(wywoz.Content, "skos:broader") ||
		!strings.Contains(wywoz.Content, "rynek-x") {
		t.Errorf("treść wywozu nie niesie relacji SKOS: %q", wywoz.Content)
	}
	if wywoz.Format != "turtle" {
		t.Errorf("wywóz bez wskazania serializacji oddał %q, domyślną jest turtle", wywoz.Format)
	}
}

// TestNasluchZewnetrznyZapisujeSieWrazZeZdarzeniami mierzy rodzinę nasłuchów.
func TestNasluchZewnetrznyZapisujeSieWrazZeZdarzeniami(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	var zapis shared.LibraryWebhookSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryWebhookSet,
		shared.LibraryWebhookSetRequest{Webhook: shared.LibraryWebhook{
			Url:     "https://przyklad.test/naslucha",
			Events:  []shared.LibraryWebhookEvent{shared.LibraryWebhookEventFileAdded},
			Enabled: true,
		}}, &zapis)

	if zapis.Webhook.Id == "" || len(zapis.Webhook.Events) != 1 {
		t.Fatalf("zapis nasłuchu oddał %q z %d zdarzeniami",
			zapis.Webhook.Id, len(zapis.Webhook.Events))
	}
	if liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM zdarzenie_webhooka_biblioteki z
		 JOIN webhook_biblioteki w ON w.id = z.webhook_id
		 WHERE w.identyfikator_zewnetrzny = ? AND z.zdarzenie = 'plik_dodany'`,
		zapis.Webhook.Id) != 1 {
		t.Error("zdarzenie nasłuchu nie doszło do bazy")
	}

	// Adres spoza HTTP jest odmową — nasłuch, który nie ma dokąd zgłosić,
	// byłby wpisem bez skutku.
	odmowa := wykonajKomende(t, zmontowany, zycie, shared.CommandLibraryWebhookSet,
		shared.LibraryWebhookSetRequest{Webhook: shared.LibraryWebhook{
			Url:    "ftp://przyklad.test",
			Events: []shared.LibraryWebhookEvent{shared.LibraryWebhookEventFileAdded},
		}})
	if odmowa.Error == nil {
		t.Error("nasłuch z adresem spoza HTTP został przyjęty")
	}

	var usuniecie shared.LibraryWebhookRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryWebhookRemove,
		shared.LibraryWebhookRemoveRequest{WebhookId: zapis.Webhook.Id}, &usuniecie)
	if !usuniecie.Removed {
		t.Fatal("usunięcie nasłuchu nie doszło do skutku")
	}
	if liczbaWierszyBiblioteki(t, baza, `SELECT COUNT(*) FROM webhook_biblioteki`) != 0 {
		t.Error("nasłuch przeżył usunięcie")
	}
}

// TestPolitykaRetencjiZglaszaZasobyPoTerminie mierzy raport retencji: polityka
// o okresie zerowym... nie istnieje, więc sprawdzian ustawia okres jednego dnia
// i pyta o zasoby z terminem w najbliższych dwóch dniach.
func TestPolitykaRetencjiZglaszaZasobyPoTerminie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuBiblioteki(t, katalog)

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "akta-stare.txt", "treść akt")

	var zapis shared.LibraryRetentionSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryRetentionSet,
		shared.LibraryRetentionSetRequest{Policy: shared.LibraryRetentionPolicy{
			Scope: shared.ConfigScopeGlobal, KeepDays: 1,
			Action: shared.LibraryRetentionActionReview,
		}}, &zapis)

	if zapis.Policy.Id == "" || zapis.AffectedFiles != 1 {
		t.Fatalf("polityka %q obejmuje %d zasobów, w repozytorium jest jeden",
			zapis.Policy.Id, zapis.AffectedFiles)
	}
	if liczbaWierszyBiblioteki(t, baza,
		`SELECT COUNT(*) FROM polityka_retencji_biblioteki WHERE identyfikator_zewnetrzny = ?`,
		zapis.Policy.Id) != 1 {
		t.Fatal("polityka nie doszła do bazy")
	}

	var raport shared.LibraryRetentionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryRetentionList,
		shared.LibraryRetentionListRequest{DueWithinDays: wskaznik(2)}, &raport)

	if len(raport.Policies) != 1 {
		t.Fatalf("raport niesie %d polityk", len(raport.Policies))
	}
	if len(raport.Due) != 1 || raport.Due[0].FileId != plik.Id {
		t.Errorf("raport terminów niesie %d pozycji — zasób miał się w nim znaleźć",
			len(raport.Due))
	}

	var usuniecie shared.LibraryRetentionSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryRetentionSet,
		shared.LibraryRetentionSetRequest{
			Policy: shared.LibraryRetentionPolicy{Id: zapis.Policy.Id},
			Remove: wskaznik(true),
		}, &usuniecie)
	if liczbaWierszyBiblioteki(t, baza, `SELECT COUNT(*) FROM polityka_retencji_biblioteki`) != 0 {
		t.Error("polityka przeżyła usunięcie")
	}
}

// TestWykazKolekcjiNiesieLicznikZasobow mierzy `library.collection.list`: licznik
// ma opisywać zawartość, a nie być zerem z kształtu odpowiedzi.
func TestWykazKolekcjiNiesieLicznikZasobow(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "materiał.txt", "treść materiału")

	var kolekcja shared.LibraryCollectionCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryCollectionCreate,
		shared.LibraryCollectionCreateRequest{
			Name: "Klienci", Description: wskaznik("materiały klientów"),
		}, &kolekcja)

	var przypisanie shared.LibraryCollectionAssignResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryCollectionAssign,
		shared.LibraryCollectionAssignRequest{
			CollectionId: kolekcja.CollectionId, FileIds: []string{plik.Id},
		}, &przypisanie)

	var wykaz shared.LibraryCollectionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryCollectionList,
		shared.LibraryCollectionListRequest{}, &wykaz)

	if wykaz.Total != 1 || len(wykaz.Collections) != 1 {
		t.Fatalf("wykaz niesie %d kolekcji", wykaz.Total)
	}
	if wykaz.Collections[0].FileCount != 1 {
		t.Errorf("kolekcja niesie licznik %d, przypisano jeden zasób",
			wykaz.Collections[0].FileCount)
	}
	if wykaz.Collections[0].Name != "Klienci" {
		t.Errorf("wykaz oddał nazwę %q", wykaz.Collections[0].Name)
	}
}

// TestOdczytMetadanychTechnicznychIdzieDoBajtow mierzy `library.metadata.get`
// z żądaniem metadanych osadzonych: suma i rozmiar mają pochodzić z bajtów.
func TestOdczytMetadanychTechnicznychIdzieDoBajtow(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	tresc := "treść mierzona z bajtów"
	plik := wgrajPlikTekstowy(t, zmontowany, zycie, "mierzony.txt", tresc)

	var bezTechnicznych shared.LibraryMetadataGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryMetadataGet,
		shared.LibraryMetadataGetRequest{FileId: plik.Id}, &bezTechnicznych)
	if bezTechnicznych.Metadata.Technical != nil {
		t.Error("odczyt bez żądania oddał metadane osadzone — miały wejść na wyraźne żądanie")
	}

	var zTechnicznymi shared.LibraryMetadataGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryMetadataGet,
		shared.LibraryMetadataGetRequest{
			FileId: plik.Id, IncludeTechnical: wskaznik(true),
		}, &zTechnicznymi)

	techniczne := zTechnicznymi.Metadata.Technical
	if techniczne == nil {
		t.Fatal("odczyt z żądaniem nie oddał metadanych osadzonych")
	}
	suma := sha256.Sum256([]byte(tresc))
	if techniczne.Checksum == nil || *techniczne.Checksum != hex.EncodeToString(suma[:]) {
		t.Errorf("suma z metadanych %v nie jest sumą bajtów treści", techniczne.Checksum)
	}
	if techniczne.SizeBytes == nil || *techniczne.SizeBytes != int64(len(tresc)) {
		t.Errorf("rozmiar z metadanych %v, treść ma %d bajtów", techniczne.SizeBytes, len(tresc))
	}
}

// jpegZeZnacznikiemXmp składa najmniejszy poprawny JPEG niosący pakiet XMP
// w segmencie APP1.
//
// Plik powstaje tutaj, a nie jest wnoszony jako materiał sprawdzianu, bo XMP
// musi być OSADZONY w bajtach — sprawdzian mierzy odczyt z pliku, więc materiał
// spoza pliku niczego by nie dowiódł. Zapis jest ręczny, bo koder `image/jpeg`
// nie umie wstawić własnego segmentu.
func jpegZeZnacznikiemXmp(t *testing.T, tytul string) []byte {
	t.Helper()

	var obraz bytes.Buffer
	plotno := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			plotno.Set(x, y, color.NRGBA{R: uint8(x * 30), G: uint8(y * 30), B: 90, A: 255})
		}
	}
	if err := jpeg.Encode(&obraz, plotno, nil); err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	bajty := obraz.Bytes()

	// Nagłówek przestrzeni nazw jest tym, po którym czytnik rozpoznaje pakiet
	// XMP wśród innych segmentów APP1 (EXIF używa tego samego znacznika).
	pakiet := append([]byte("http://ns.adobe.com/xap/1.0/\x00"),
		[]byte(`<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?>`+
			`<x:xmpmeta xmlns:x="adobe:ns:meta/">`+
			`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">`+
			`<rdf:Description rdf:about="" xmlns:dc="http://purl.org/dc/elements/1.1/">`+
			`<dc:title><rdf:Alt><rdf:li xml:lang="x-default">`+tytul+
			`</rdf:li></rdf:Alt></dc:title>`+
			`</rdf:Description></rdf:RDF></x:xmpmeta><?xpacket end="w"?>`)...)

	// Długość segmentu liczy się wraz z dwoma bajtami samej długości.
	dlugosc := len(pakiet) + 2
	segment := []byte{0xFF, 0xE1, byte(dlugosc >> 8), byte(dlugosc & 0xFF)}
	segment = append(segment, pakiet...)

	// Segment wchodzi zaraz za znacznikiem początku pliku (dwa bajty `SOI`).
	zlozony := make([]byte, 0, len(bajty)+len(segment))
	zlozony = append(zlozony, bajty[:2]...)
	zlozony = append(zlozony, segment...)
	zlozony = append(zlozony, bajty[2:]...)
	return zlozony
}

// TestOdczytMetadanychOsadzonychCzytaXmpProgramem mierzy pole `xmp` kontraktu:
// przed tą pracą nie wypełniała go żadna droga, więc opis zasobu milczał o XMP,
// IPTC i ID3 niezależnie od tego, co plik naprawdę niósł.
//
// Sprawdzian pomija się z nazwanym powodem na maszynie bez programu: odczyt tych
// trzech rodzin jest POSZERZENIEM opisu, a nie jego warunkiem — reszty pól
// pilnuje `TestOdczytMetadanychTechnicznychIdzieDoBajtow`, który idzie zawsze.
func TestOdczytMetadanychOsadzonychCzytaXmpProgramem(t *testing.T) {
	if !zewnetrzne.Stoi(narzedzieMetadanychBiblioteki) {
		t.Skipf("na tej maszynie nie stoi %s (%s) — IPTC, XMP i ID3 są poszerzeniem opisu, "+
			"więc bez programu nie ma czego mierzyć",
			narzedzieMetadanychBiblioteki.Nazwa, narzedzieMetadanychBiblioteki.Program)
	}

	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	const tytul = "Tytul osadzony w bajtach"
	var wgrany shared.LibraryFileUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileUpload,
		shared.LibraryFileUploadRequest{
			Name:          "zdjecie.jpg",
			MimeType:      wskaznik("image/jpeg"),
			ContentBase64: wskaznik(wBase64(jpegZeZnacznikiemXmp(t, tytul))),
		}, &wgrany)

	var odczyt shared.LibraryMetadataGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryMetadataGet,
		shared.LibraryMetadataGetRequest{
			FileId: wgrany.File.Id, IncludeTechnical: wskaznik(true),
		}, &odczyt)

	techniczne := odczyt.Metadata.Technical
	if techniczne == nil {
		t.Fatal("odczyt z żądaniem nie oddał metadanych osadzonych")
	}
	if len(techniczne.Xmp) == 0 {
		t.Fatal("pole xmp jest puste, choć plik niesie pakiet XMP w segmencie APP1")
	}
	if !strings.Contains(string(techniczne.Xmp), tytul) {
		t.Fatalf("pole xmp nie niesie tytułu osadzonego w pliku: %s", techniczne.Xmp)
	}
	t.Logf("pole xmp: %s", techniczne.Xmp)
}

// mp3ZeZnacznikiemId3 składa najmniejszy zapis MP3 niosący znacznik ID3v2
// z ramką tytułu.
//
// Zapis jest ręczny z tego samego powodu, co przy XMP wyżej: znacznik musi być
// OSADZONY w bajtach pliku, a biblioteka standardowa nie zapisuje MP3 wcale.
func mp3ZeZnacznikiemId3(t *testing.T, tytul string) []byte {
	t.Helper()

	// Ramka TIT2 (tytuł): bajt kodowania ISO-8859-1 i sama treść.
	tresc := append([]byte{0x00}, []byte(tytul)...)
	ramka := append([]byte("TIT2"), byte(len(tresc)>>24), byte(len(tresc)>>16),
		byte(len(tresc)>>8), byte(len(tresc)), 0x00, 0x00)
	ramka = append(ramka, tresc...)

	// Nagłówek ID3v2.3. Rozmiar znacznika idzie w zapisie synchsafe — siedem
	// bitów na bajt — bo tak każe sam zapis ID3v2, nie wybór tego sprawdzianu.
	rozmiar := len(ramka)
	zapis := append([]byte{'I', 'D', '3', 0x03, 0x00, 0x00,
		byte(rozmiar >> 21 & 0x7F), byte(rozmiar >> 14 & 0x7F),
		byte(rozmiar >> 7 & 0x7F), byte(rozmiar & 0x7F)}, ramka...)

	// Trzy ramki MPEG-1 Layer III ciszy (128 kb/s, 44,1 kHz) — bez nich plik
	// byłby samym znacznikiem bez nagrania i nie uchodziłby za materiał.
	cisza := append([]byte{0xFF, 0xFB, 0x90, 0x00}, make([]byte, 413)...)
	for i := 0; i < 3; i++ {
		zapis = append(zapis, cisza...)
	}
	return zapis
}

// TestOdczytMetadanychOsadzonychCzytaId3Programem mierzy pole `id3` kontraktu —
// trzecią z rodzin, które czyta program (`dopiszMetadaneOsadzone`), i jedyną
// dotyczącą nagrania, nie obrazu.
func TestOdczytMetadanychOsadzonychCzytaId3Programem(t *testing.T) {
	if !zewnetrzne.Stoi(narzedzieMetadanychBiblioteki) {
		t.Skipf("na tej maszynie nie stoi %s (%s) — IPTC, XMP i ID3 są poszerzeniem opisu, "+
			"więc bez programu nie ma czego mierzyć",
			narzedzieMetadanychBiblioteki.Nazwa, narzedzieMetadanychBiblioteki.Program)
	}

	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	const tytul = "Tytul osadzony w bajtach"
	var wgrany shared.LibraryFileUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileUpload,
		shared.LibraryFileUploadRequest{
			Name:          "nagranie.mp3",
			MimeType:      wskaznik("audio/mpeg"),
			ContentBase64: wskaznik(wBase64(mp3ZeZnacznikiemId3(t, tytul))),
		}, &wgrany)

	var odczyt shared.LibraryMetadataGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryMetadataGet,
		shared.LibraryMetadataGetRequest{
			FileId: wgrany.File.Id, IncludeTechnical: wskaznik(true),
		}, &odczyt)

	techniczne := odczyt.Metadata.Technical
	if techniczne == nil {
		t.Fatal("odczyt z żądaniem nie oddał metadanych osadzonych")
	}
	if len(techniczne.Id3) == 0 {
		t.Fatal("pole id3 jest puste, choć plik niesie znacznik ID3v2 z tytułem")
	}
	if !strings.Contains(string(techniczne.Id3), tytul) {
		t.Fatalf("pole id3 nie niesie tytułu osadzonego w pliku: %s", techniczne.Id3)
	}
	t.Logf("pole id3: %s", techniczne.Id3)
}

// TestOdczytBezProgramuZostawiaPolaOsadzonePuste mierzy drugą połowę reguły
// z nagłówka `dopiszMetadaneOsadzone`: brak programu zostawia pola IPTC, XMP
// i ID3 puste, a odczyt opisu NIE odmawia i reszta pól przychodzi w komplecie.
//
// Pustą ścieżką wyszukiwania sprawdzian czyni z tej maszyny maszynę bez
// programu, więc mierzy to zdanie wszędzie — nie tylko na cienkiej instalce.
func TestOdczytBezProgramuZostawiaPolaOsadzonePuste(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var wgrany shared.LibraryFileUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryFileUpload,
		shared.LibraryFileUploadRequest{
			Name:          "zdjecie.jpg",
			MimeType:      wskaznik("image/jpeg"),
			ContentBase64: wskaznik(wBase64(jpegZeZnacznikiemXmp(t, "Tytul osadzony w bajtach"))),
		}, &wgrany)

	var odczyt shared.LibraryMetadataGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandLibraryMetadataGet,
		shared.LibraryMetadataGetRequest{
			FileId: wgrany.File.Id, IncludeTechnical: wskaznik(true),
		}, &odczyt)

	techniczne := odczyt.Metadata.Technical
	if techniczne == nil {
		t.Fatal("odczyt bez programu nie oddał metadanych osadzonych — miał oddać opis węższy, nie żaden")
	}
	if len(techniczne.Xmp) != 0 {
		t.Fatalf("pole xmp niesie %s, choć maszyna nie ma czym go przeczytać — "+
			"wynik bez pomiaru podany jako wynik", techniczne.Xmp)
	}
	// Wymiary liczy czytnik wkompilowany, więc mają przyjść także bez programu —
	// to one dowodzą, że opis zwęził się o trzy pola, a nie wywrócił.
	if techniczne.Width == nil || *techniczne.Width != 8 ||
		techniczne.Height == nil || *techniczne.Height != 8 {
		t.Fatalf("wymiary %v×%v nie przyszły z czytnika wkompilowanego, a od programu nie zależą",
			techniczne.Width, techniczne.Height)
	}
}

// TestKazdaKomendaBibliotekiMaUchwyt jest zaporą pokrycia rodziny: kontrakt
// niesie 48 komend `library.*` i każda ma mieć uchwyt w rejestrze rdzenia.
//
// Sprawdzian pokrycia całego kontraktu stoi osobno
// (`TestRejestrPokrywaKomendyKontraktu`), ale zapora rodziny jest tu, bo to ten
// moduł ma nie zostawić komendy bez drogi.
func TestKazdaKomendaBibliotekiMaUchwyt(t *testing.T) {
	zmontowany, _, _ := zmontujDoPomiaruSkutku(t)

	obslugiwane := zbiorNazw(zmontowany.Rdzen.rejestr.Nazwy())
	brakujace := []string{}
	for _, komenda := range shared.WszystkieKomendy() {
		if !strings.HasPrefix(string(komenda), "library.") {
			continue
		}
		if !obslugiwane[komenda] {
			brakujace = append(brakujace, string(komenda))
		}
	}
	if len(brakujace) > 0 {
		t.Errorf("komendy rodziny bez uchwytu: %s", strings.Join(brakujace, ", "))
	}
}
