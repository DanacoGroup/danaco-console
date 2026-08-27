// Sprawdziany tego pliku mierzą skutek piętnastu czynności dobudowanych do
// modułu Studio: nie czy odpowiedź jest zgodna z kontraktem, lecz czy za nią
// coś zostało.
package core

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// oknoSprawdzianuStudia jest oknem, w którego imieniu idą wszystkie żądania
// tego pliku. Ta sama nazwa, którą bierze `sciezkaZasobuSprawdzianu`: zasób
// odłożony do jednego okna, a szukany w drugim, nie znalazłby się nigdy.
const oknoSprawdzianuStudia = "okno-sprawdzianu"

// polaczenieOboczneStudia otwiera drugie połączenie do pliku bazy sprawdzianu,
// drogą, której mierzony kod nie kontroluje.
func polaczenieOboczneStudia(t *testing.T, katalog string) *sql.DB {
	t.Helper()

	sciezka := filepath.ToSlash(filepath.Join(katalog, "dane.sqlite"))
	baza, err := sql.Open("sqlite", sciezka+"?_pragma=busy_timeout(5000)")
	if err != nil {
		t.Fatalf("nie można otworzyć obocznego połączenia do bazy: %v", err)
	}
	t.Cleanup(func() { _ = baza.Close() })
	return baza
}

// dokumentSprawdzianuStudia zakłada dokument komendami warsztatu i oddaje jego
// identyfikator do dalszej pracy sprawdzianu.
func dokumentSprawdzianuStudia(t *testing.T, zmontowany *Zmontowany,
	zycie context.Context, tresc string) string {
	t.Helper()

	var otwarcie shared.StudioDocumentOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentOpen,
		shared.StudioDocumentOpenRequest{WindowId: oknoSprawdzianuStudia}, &otwarcie)

	var zapis shared.StudioDocumentSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId: otwarcie.Document.Id,
			Content:    tresc,
			Title:      wskaznik("dokument sprawdzianu"),
		}, &zapis)
	return zapis.Document.Id
}

// wersjaSprawdzianuStudia zapisuje treść jako nową wersję dokumentu i oddaje
// jej kod do dalszej pracy sprawdzianu.
func wersjaSprawdzianuStudia(t *testing.T, zmontowany *Zmontowany,
	zycie context.Context, kodDokumentu, tresc string) string {
	t.Helper()

	var zapis shared.StudioDocumentSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId: kodDokumentu, Content: tresc, CreateVersion: wskaznik(true),
		}, &zapis)
	if zapis.Version == nil {
		t.Fatal("zapis z żądaniem wersji nie oddał wersji")
	}
	return zapis.Version.Id
}

// wpisyArchiwumSprawdzianu rozpakowuje archiwum ZIP i oddaje jego zawartość
// jako mapę nazwy pliku na jego bajty.
func wpisyArchiwumSprawdzianu(t *testing.T, bajty []byte) map[string][]byte {
	t.Helper()

	czytnik, err := zip.NewReader(bytes.NewReader(bajty), int64(len(bajty)))
	if err != nil {
		t.Fatalf("wynik nie jest archiwum ZIP: %v", err)
	}
	zawartosc := map[string][]byte{}
	for _, plik := range czytnik.File {
		strumien, err := plik.Open()
		if err != nil {
			t.Fatalf("nie można otworzyć wpisu %s: %v", plik.Name, err)
		}
		var bufor bytes.Buffer
		if _, err := bufor.ReadFrom(strumien); err != nil {
			t.Fatalf("nie można odczytać wpisu %s: %v", plik.Name, err)
		}
		_ = strumien.Close()
		zawartosc[plik.Name] = bufor.Bytes()
	}
	return zawartosc
}

// TestWydanieRepozytoriumDajeArchiwumZTresciaWersji rozpakowuje archiwum
// i sprawdza, że pliki wersji niosą TĘ treść, którą Operator zapisał.
func TestWydanieRepozytoriumDajeArchiwumZTresciaWersji(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	kodDokumentu := dokumentSprawdzianuStudia(t, zmontowany, zycie, "wersja pierwsza")
	kodPierwszej := wersjaSprawdzianuStudia(t, zmontowany, zycie, kodDokumentu, "wersja pierwsza")
	kodDrugiej := wersjaSprawdzianuStudia(t, zmontowany, zycie, kodDokumentu, "wersja druga, dłuższa")

	var wynik shared.StudioRepositoryExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioRepositoryExport,
		shared.StudioRepositoryExportRequest{
			DocumentId: kodDokumentu, WindowId: wskaznik(oknoSprawdzianuStudia),
		}, &wynik)

	if wynik.Entries != 2 {
		t.Fatalf("archiwum melduje %d wersji, a dokument ma dwie", wynik.Entries)
	}
	if wynik.Asset.Uri == nil {
		t.Fatal("archiwum nie niesie odwołania do bajtów")
	}
	bajty := bajtyPodOdwolaniem(t, katalog, *wynik.Asset.Uri)
	if len(bajty) != wynik.SizeBytes {
		t.Fatalf("odpowiedź melduje %d bajtów, a plik ma %d", wynik.SizeBytes, len(bajty))
	}

	wpisy := wpisyArchiwumSprawdzianu(t, bajty)
	if string(wpisy["wersje/"+kodPierwszej+".txt"]) != "wersja pierwsza" {
		t.Fatalf("wpis pierwszej wersji niesie %q zamiast treści zapisanej",
			string(wpisy["wersje/"+kodPierwszej+".txt"]))
	}
	if string(wpisy["wersje/"+kodDrugiej+".txt"]) != "wersja druga, dłuższa" {
		t.Fatalf("wpis drugiej wersji niesie %q zamiast treści zapisanej",
			string(wpisy["wersje/"+kodDrugiej+".txt"]))
	}

	surowy, jest := wpisy["manifest.json"]
	if !jest {
		t.Fatal("archiwum nie niesie manifestu, choć wydanie ma go nieść zawsze")
	}
	var manifest manifestWydaniaStudia
	if err := json.Unmarshal(surowy, &manifest); err != nil {
		t.Fatalf("manifest jest nieczytelny: %v", err)
	}
	if len(manifest.Wersje) != 2 || manifest.Dokument != kodDokumentu {
		t.Fatalf("manifest opisuje %d wersji dokumentu %q", len(manifest.Wersje), manifest.Dokument)
	}
}

// TestPaczkaRedakcyjnaNiesieDokumentHistorieIRaport pilnuje, żeby paczka
// przekazania naprawdę miała wszystkie cztery człony, a nie tylko je meldowała.
func TestPaczkaRedakcyjnaNiesieDokumentHistorieIRaport(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	kodDokumentu := dokumentSprawdzianuStudia(t, zmontowany, zycie, "materiał wyjściowy")
	wersjaSprawdzianuStudia(t, zmontowany, zycie, kodDokumentu, "materiał wyjściowy")
	wersjaSprawdzianuStudia(t, zmontowany, zycie, kodDokumentu, "materiał po redakcji")

	var wynik shared.StudioPackageExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPackageExport,
		shared.StudioPackageExportRequest{
			DocumentId:     kodDokumentu,
			DocumentFormat: wskaznik("markdown"),
			WindowId:       wskaznik(oknoSprawdzianuStudia),
		}, &wynik)

	if wynik.Asset.Uri == nil {
		t.Fatal("paczka nie niesie odwołania do bajtów")
	}
	wpisy := wpisyArchiwumSprawdzianu(t, bajtyPodOdwolaniem(t, katalog, *wynik.Asset.Uri))

	if string(wpisy["dokument.markdown"]) != "materiał po redakcji" {
		t.Fatalf("dokument finalny paczki niesie %q", string(wpisy["dokument.markdown"]))
	}
	if _, jest := wpisy["raport-zmian.md"]; !jest {
		t.Fatal("paczka nie niesie raportu zmian, choć nie był wyłączony")
	}
	if _, jest := wpisy["adnotacje.json"]; !jest {
		t.Fatal("paczka nie niesie adnotacji, choć nie były wyłączone")
	}
	historii := 0
	for nazwa := range wpisy {
		if strings.HasPrefix(nazwa, "wersje/") {
			historii++
		}
	}
	if historii != 2 {
		t.Fatalf("paczka niesie %d plików historii, a dokument ma dwie wersje", historii)
	}
	if wynik.Entries != len(wpisy) {
		t.Fatalf("paczka melduje %d pozycji, a archiwum ma %d", wynik.Entries, len(wpisy))
	}
}

// TestRaportRoznicyWDocxJestOtwieralnymDokumentem sprawdza, że raport wydany
// jako DOCX jest naprawdę paczką biurową z częścią treści, a nie plikiem
// tekstowym pod cudzym rozszerzeniem.
func TestRaportRoznicyWDocxJestOtwieralnymDokumentem(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	kodDokumentu := dokumentSprawdzianuStudia(t, zmontowany, zycie, "zdanie pierwsze")
	kodBazy := wersjaSprawdzianuStudia(t, zmontowany, zycie, kodDokumentu, "zdanie pierwsze")
	kodCelu := wersjaSprawdzianuStudia(t, zmontowany, zycie, kodDokumentu, "zdanie pierwsze zmienione")

	var wynik shared.StudioDiffReportExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDiffReportExport,
		shared.StudioDiffReportExportRequest{
			DocumentId:      kodDokumentu,
			BaseVersionId:   kodBazy,
			TargetVersionId: &kodCelu,
			Format:          "docx",
			WindowId:        wskaznik(oknoSprawdzianuStudia),
		}, &wynik)

	if wynik.Asset.Uri == nil {
		t.Fatal("raport nie niesie odwołania do bajtów")
	}
	wpisy := wpisyArchiwumSprawdzianu(t, bajtyPodOdwolaniem(t, katalog, *wynik.Asset.Uri))
	tresc, jest := wpisy["word/document.xml"]
	if !jest {
		t.Fatal("plik DOCX nie ma części `word/document.xml` — nie otworzy go żaden edytor")
	}
	if _, jest := wpisy["[Content_Types].xml"]; !jest {
		t.Fatal("plik DOCX nie ma wykazu typów zawartości")
	}
	if !strings.Contains(string(tresc), "Raport zmian") {
		t.Fatal("część treści DOCX nie niesie raportu zmian")
	}
	if !strings.Contains(string(tresc), "zdanie pierwsze zmienione") {
		t.Fatal("raport nie pokazuje treści po zmianie, choć wersje się różnią")
	}
}

// TestWyrysPodgladuDajeStroneDajacaSieOdkodowac wykazuje skutek wyrysu: pod
// odwołaniem leżą bajty, które dekoder PNG czyta jako obraz o wymiarach strony.
func TestWyrysPodgladuDajeStroneDajacaSieOdkodowac(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	tresc := strings.Repeat("Akapit podglądu z polskimi znakami: żółć, ćma, ężąś.\n\n", 40)
	kodDokumentu := dokumentSprawdzianuStudia(t, zmontowany, zycie, tresc)

	var wynik shared.StudioPreviewRenderResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPreviewRender,
		shared.StudioPreviewRenderRequest{
			DocumentId: kodDokumentu, Format: "png",
			WindowId: wskaznik(oknoSprawdzianuStudia),
		}, &wynik)

	if wynik.Pages < 1 || len(wynik.PageAssetIds) != wynik.Pages {
		t.Fatalf("podgląd melduje %d stron przy %d zasobach", wynik.Pages, len(wynik.PageAssetIds))
	}
	for _, kod := range wynik.PageAssetIds {
		bajty := bajtyPlikuSprawdzianu(t,
			sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, kod))
		obraz, err := png.Decode(bytes.NewReader(bajty))
		if err != nil {
			t.Fatalf("zasób %s nie daje się odkodować jako obraz: %v", kod, err)
		}
		if obraz.Bounds().Dx() != szerokoscStronyStudia ||
			obraz.Bounds().Dy() != wysokoscStronyStudia {
			t.Fatalf("strona %s ma wymiary %v, a strona wyrysu ma %dx%d",
				kod, obraz.Bounds(), szerokoscStronyStudia, wysokoscStronyStudia)
		}
		if bialaStronaSprawdzianu(obraz) {
			t.Fatalf("strona %s jest w całości biała — wyrys nic na niej nie postawił", kod)
		}
	}
}

// bajtyPlikuSprawdzianu czyta plik spod ścieżki bezwzględnej i przerywa
// sprawdzian, gdy jest pusty. Plik zerowej długości przechodzi każdy sprawdzian
// istnienia i nie pokazuje nic.
func bajtyPlikuSprawdzianu(t *testing.T, sciezka string) []byte {
	t.Helper()

	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać %s: %v", sciezka, err)
	}
	if len(bajty) == 0 {
		t.Fatalf("plik %s ma zerową długość", sciezka)
	}
	return bajty
}

// bialaStronaSprawdzianu mówi, czy na stronie nie ma ani jednego ciemnego
// piksela. Strona biała przechodzi każdy sprawdzian formatu i nie pokazuje nic.
func bialaStronaSprawdzianu(obraz image.Image) bool {
	prostokat := obraz.Bounds()
	for y := prostokat.Min.Y; y < prostokat.Max.Y; y++ {
		for x := prostokat.Min.X; x < prostokat.Max.X; x++ {
			if jasnoscStudia(obraz, x, y) < 200 {
				return false
			}
		}
	}
	return true
}

// TestRoznicaWizualnaWskazujeObszarIDajeNakladke wykazuje, że porównanie
// wizualne widzi zmianę i odkłada nakładkę, którą da się odkodować.
func TestRoznicaWizualnaWskazujeObszarIDajeNakladke(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	kodDokumentu := dokumentSprawdzianuStudia(t, zmontowany, zycie, "treść wyjściowa")
	kodBazy := wersjaSprawdzianuStudia(t, zmontowany, zycie, kodDokumentu, "treść wyjściowa")
	kodCelu := wersjaSprawdzianuStudia(t, zmontowany, zycie, kodDokumentu,
		"treść wyjściowa, a po niej zdanie dołożone w redakcji")

	var wynik shared.StudioDiffVisualResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDiffVisual,
		shared.StudioDiffVisualRequest{
			DocumentId: kodDokumentu, BaseVersionId: kodBazy, TargetVersionId: kodCelu,
			WindowId: wskaznik(oknoSprawdzianuStudia),
		}, &wynik)

	if len(wynik.Regions) == 0 {
		t.Fatal("porównanie wizualne nie wskazało ani jednego obszaru, choć wersje różnią się treścią")
	}
	for _, obszar := range wynik.Regions {
		if obszar.ChangeRatio <= 0 || obszar.Width <= 0 || obszar.Height <= 0 {
			t.Fatalf("obszar %v nie ma ani wymiarów, ani udziału zmiany", obszar)
		}
	}
	if len(wynik.OverlayAssetIds) == 0 {
		t.Fatal("porównanie nie odłożyło ani jednej nakładki, choć obszary wskazało")
	}
	for _, kod := range wynik.OverlayAssetIds {
		bajty := bajtyPlikuSprawdzianu(t,
			sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, kod))
		if _, err := png.Decode(bytes.NewReader(bajty)); err != nil {
			t.Fatalf("nakładka %s nie daje się odkodować jako obraz: %v", kod, err)
		}
	}
}

// TestGalazIOdwolanieLezaWBazie odczytuje oba byty drugim połączeniem do pliku
// bazy — rdzeń, który melduje zapis bez wiersza, nie ma tu jak się obronić.
func TestGalazIOdwolanieLezaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneStudia(t, katalog)

	kodDokumentu := dokumentSprawdzianuStudia(t, zmontowany, zycie, "pień dokumentu")
	kodStartowej := wersjaSprawdzianuStudia(t, zmontowany, zycie, kodDokumentu, "pień dokumentu")

	var galaz shared.StudioBranchCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioBranchCreate,
		shared.StudioBranchCreateRequest{
			DocumentId: kodDokumentu, FromVersionId: kodStartowej, Name: "redakcja klienta",
		}, &galaz)

	var nazwa, wersjaStartowa string
	err := oboczne.QueryRow(
		`SELECT nazwa, wersja_startowa_id FROM galaz_studio WHERE identyfikator_zewnetrzny = ?`,
		galaz.Branch.Id).Scan(&nazwa, &wersjaStartowa)
	if err != nil {
		t.Fatalf("gałęzi %s nie ma w bazie, choć komenda zameldowała jej założenie: %v",
			galaz.Branch.Id, err)
	}
	if nazwa != "redakcja klienta" || wersjaStartowa != kodStartowej {
		t.Fatalf("wiersz gałęzi niesie nazwę %q i wersję startową %q", nazwa, wersjaStartowa)
	}

	var wykaz shared.StudioBranchListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioBranchList,
		shared.StudioBranchListRequest{DocumentId: kodDokumentu}, &wykaz)
	if len(wykaz.Branches) != 1 || wykaz.Branches[0].Id != galaz.Branch.Id {
		t.Fatalf("wykaz gałęzi oddał %d pozycji", len(wykaz.Branches))
	}

	var odwolanie shared.StudioVersionReferenceCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioVersionReferenceCreate,
		shared.StudioVersionReferenceCreateRequest{VersionId: kodStartowej}, &odwolanie)

	if !strings.Contains(odwolanie.Reference, kodStartowej) {
		t.Fatalf("odwołanie %q nie wskazuje wersji, do której miało powstać", odwolanie.Reference)
	}
	var wierszy int
	if err := oboczne.QueryRow(
		`SELECT COUNT(*) FROM odwolanie_wersji_studio WHERE wersja_id = ? AND zasieg = ?`,
		kodStartowej, string(odwolanie.Scope)).Scan(&wierszy); err != nil {
		t.Fatalf("nie można policzyć odwołań w bazie: %v", err)
	}
	if wierszy != 1 {
		t.Fatalf("w bazie leży %d odwołań do wersji, a komenda zameldowała jedno", wierszy)
	}
}

// TestScalenieGalezinOddajeKonfliktZamiastRozstrzygacGo pilnuje zasady:
// zmiana obu stron w tym samym miejscu wraca konfliktem, a nie cichym wyborem.
func TestScalenieGalezinOddajeKonfliktZamiastRozstrzygacGo(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	kodDokumentu := dokumentSprawdzianuStudia(t, zmontowany, zycie, "wiersz A\nwiersz B\nwiersz C")
	kodStartowej := wersjaSprawdzianuStudia(t, zmontowany, zycie, kodDokumentu,
		"wiersz A\nwiersz B\nwiersz C")

	var pierwsza, druga shared.StudioBranchCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioBranchCreate,
		shared.StudioBranchCreateRequest{
			DocumentId: kodDokumentu, FromVersionId: kodStartowej, Name: "redakcja pierwsza",
		}, &pierwsza)
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioBranchCreate,
		shared.StudioBranchCreateRequest{
			DocumentId: kodDokumentu, FromVersionId: kodStartowej, Name: "redakcja druga",
		}, &druga)

	// Obie gałęzie ruszają ten sam wiersz, i to inaczej — to jest konflikt
	// z definicji.
	przestawCzoloGalezi(t, zmontowany, zycie, kodDokumentu, pierwsza.Branch.Id,
		"wiersz A\nwiersz B od pierwszej\nwiersz C")
	przestawCzoloGalezi(t, zmontowany, zycie, kodDokumentu, druga.Branch.Id,
		"wiersz A\nwiersz B od drugiej\nwiersz C")

	var scalenie shared.StudioBranchMergeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioBranchMerge,
		shared.StudioBranchMergeRequest{
			SourceBranchId: pierwsza.Branch.Id, TargetBranchId: druga.Branch.Id,
		}, &scalenie)

	if scalenie.Merged || len(scalenie.Conflicts) == 0 {
		t.Fatalf("scalenie zameldowało %v przy %d konfliktach — rozstrzygnęło za Operatora",
			scalenie.Merged, len(scalenie.Conflicts))
	}
	if scalenie.Conflicts[0].Source == scalenie.Conflicts[0].Target {
		t.Fatal("konflikt niesie tę samą treść po obu stronach — nie jest konfliktem")
	}

	// Rozstrzygnięcie podane z góry ma scalenie domknąć i zapisać wybraną treść.
	var domkniete shared.StudioBranchMergeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioBranchMerge,
		shared.StudioBranchMergeRequest{
			SourceBranchId: pierwsza.Branch.Id, TargetBranchId: druga.Branch.Id,
			Resolutions: []shared.StudioMergeResolution{
				{Index: 1, Side: shared.StudioMergeSideScalana},
			},
		}, &domkniete)

	if !domkniete.Merged || domkniete.Document == nil {
		t.Fatalf("scalenie z rozstrzygnięciem nie doszło do skutku: %+v", domkniete)
	}
	if domkniete.Document.Content == nil ||
		!strings.Contains(*domkniete.Document.Content, "wiersz B od pierwszej") {
		t.Fatal("dokument po scaleniu nie niesie treści wybranej rozstrzygnięciem")
	}
}

// przestawCzoloGalezi dopisuje wersję na gałęzi i przestawia na nią jej czoło,
// drogą warstwy danych, bo kontrakt Studia nie ma dziś polecenia zapisu na
// gałęzi.
func przestawCzoloGalezi(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	kodDokumentu, kodGalezi, tresc string) {
	t.Helper()

	repozytorium := zmontowany.dane.Studio
	dokument, err := repozytorium.Dokument(zycie, kodDokumentu)
	if err != nil {
		t.Fatalf("nie można odczytać dokumentu %s: %v", kodDokumentu, err)
	}
	galaz, err := repozytorium.Galaz(zycie, kodGalezi)
	if err != nil {
		t.Fatalf("nie można odczytać gałęzi %s: %v", kodGalezi, err)
	}
	wersja, err := repozytorium.ZapiszWersje(zycie, dokument.ID, dane.WersjaDokumentu{
		Kod: nowyIdentyfikator(przedrostekWersjiStudio), Tresc: &tresc, GalazKod: &kodGalezi,
	})
	if err != nil {
		t.Fatalf("nie można zapisać wersji na gałęzi %s: %v", kodGalezi, err)
	}
	galaz.WersjaBiezacaID = &wersja.Kod
	if _, err := repozytorium.ZapiszGalaz(zycie, dokument.ID, galaz); err != nil {
		t.Fatalf("nie można przestawić czoła gałęzi %s: %v", kodGalezi, err)
	}
}

// TestWsadOdrzucaDokumentBezKanaluAleNieWstrzymujePozostalych pilnuje reguły
// z rozdz. 4.5: odmowa jednego dokumentu nie przerywa wsadu.
func TestWsadOdrzucaDokumentBezKanaluAleNieWstrzymujePozostalych(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	pierwszy := dokumentSprawdzianuStudia(t, zmontowany, zycie, "pierwszy")
	drugi := dokumentSprawdzianuStudia(t, zmontowany, zycie, "drugi")

	var wynik shared.StudioBatchRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioBatchRun,
		shared.StudioBatchRunRequest{
			WindowId:    oknoSprawdzianuStudia,
			DocumentIds: []string{pierwszy, drugi, "studio-dok-nie-ma-takiego"},
			ActionId:    "skroc",
		}, &wynik)

	if wynik.RunId == "" {
		t.Fatal("wsad nie oddał identyfikatora przebiegu")
	}
	if len(wynik.Rejected) != 3 {
		t.Fatalf("wsad odrzucił %d dokumentów; bez kanału modelu odmówić mają wszystkie trzy",
			len(wynik.Rejected))
	}
	for _, odrzucony := range wynik.Rejected {
		if strings.TrimSpace(odrzucony.Reason) == "" {
			t.Fatalf("dokument %s odrzucony bez podania powodu", odrzucony.DocumentId)
		}
	}
}

// TestWyszukiwanieZnaczenioweMowiKtoraDrogaPoszlo pilnuje, żeby czynność
// oddała wynik I nazwała drogę — cisza zamiast wyniku jest niedopuszczalna,
// a wynik bez nazwy drogi nie mówi, ile jest wart.
func TestWyszukiwanieZnaczenioweMowiKtoraDrogaPoszlo(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	tresc := "Umowa o dzieło z terminem wykonania.\n\n" +
		"Kara umowna za zwłokę wykonawcy.\n\n" +
		"Przepis kuchenny na żurek śląski."
	kodDokumentu := dokumentSprawdzianuStudia(t, zmontowany, zycie, tresc)

	var wynik shared.StudioSearchSemanticResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioSearchSemantic,
		shared.StudioSearchSemanticRequest{
			DocumentId: kodDokumentu, Query: "kara umowna za zwłokę",
		}, &wynik)

	if wynik.Mode != drogaSemantykiMiara && wynik.Mode != drogaSemantykiModelem {
		t.Fatalf("odpowiedź nie nazywa drogi liczenia: %q", wynik.Mode)
	}
	if len(wynik.Matches) == 0 {
		t.Fatal("wyszukiwanie nie oddało ani jednego fragmentu, choć dokument zawiera zapytanie wprost")
	}
	if !strings.Contains(wynik.Matches[0].Text, "Kara umowna") {
		t.Fatalf("najbliższym fragmentem jest %q, a nie ten, który mówi o karze umownej",
			wynik.Matches[0].Text)
	}
	if wynik.Matches[0].Score <= 0 {
		t.Fatal("najbliższy fragment ma bliskość zerową")
	}
}

// TestSkanowanieZUrzadzeniaOdmawiaNazwanie pilnuje zasady bezwzględnej:
// czynność bez drogi ma odmówić, nazywając brak, a nie oddać pusty wykaz
// pozycji.
func TestSkanowanieZUrzadzeniaOdmawiaNazwanie(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	if !zewnetrzne.Stoi(narzedzieSkanera) {
		t.Skipf("pomiar niewykonany: na tej maszynie nie ma programu %s (%s), więc"+
			" skanowanie odmawia brakiem warstwy, a nie brakiem urządzenia",
			narzedzieSkanera.Nazwa, narzedzieSkanera.Program)
	}
	var wykaz shared.StudioIngestDeviceListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestDeviceList,
		shared.StudioIngestDeviceListRequest{}, &wykaz)
	if len(wykaz.Devices) > 0 {
		t.Skipf("pomiar niewykonany: warstwa skanera widzi %d urządzeń, więc skanowanie"+
			" nie odmawia brakiem urządzenia", len(wykaz.Devices))
	}

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandStudioIngestDeviceScan,
		shared.StudioIngestDeviceScanRequest{WindowId: oknoSprawdzianuStudia})

	// Brak urządzenia jest stanem maszyny, nie usterką rdzenia — kod odmowy ma
	// to odzwierciedlać.
	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Errorf("odmowa niesie kod %q, oczekiwany %q — brak urządzenia nie jest usterką"+
			" rdzenia; treść: %s", odmowa.Code, shared.ErrorCodeNotFound, odmowa.Message)
	}
	komunikat := strings.ToLower(odmowa.Message)
	for _, slowo := range []string{"operator", "scanimage", "queue.add"} {
		if !strings.Contains(komunikat, slowo) {
			t.Fatalf("odmowa nie nazywa braku słowem %q: %s", slowo, komunikat)
		}
	}
}

// TestOsadzenieZasobuWstawiaOdwolanieWTresc wykazuje skutek osadzenia: treść
// dokumentu naprawdę niesie odwołanie do zasobu, i to w miejscu wskazanym.
func TestOsadzenieZasobuWstawiaOdwolanieWTresc(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	kodDokumentu := dokumentSprawdzianuStudia(t, zmontowany, zycie, "PRZED|PO")

	var wgranie shared.DesignAssetUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetUpload,
		shared.DesignAssetUploadRequest{
			WindowId:      oknoSprawdzianuStudia,
			Name:          wskaznik("rysunek"),
			Kind:          shared.DesignAssetKindImage,
			Format:        wskaznik("png"),
			ContentBase64: wskaznik(wBase64(obrazPNG(t, 16, 16))),
		}, &wgranie)

	var wynik shared.StudioAssetEmbedResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioAssetEmbed,
		shared.StudioAssetEmbedRequest{
			DocumentId: kodDokumentu, AssetId: wgranie.Asset.Id,
			Position: wskaznik(5), Caption: wskaznik("Rysunek 1"),
		}, &wynik)

	if wynik.Document.Content == nil {
		t.Fatal("dokument po osadzeniu nie niesie treści")
	}
	tresc := *wynik.Document.Content
	if !strings.Contains(tresc, "danaco://zasob/"+wgranie.Asset.Id) {
		t.Fatalf("treść nie niesie odwołania do zasobu: %q", tresc)
	}
	if !strings.Contains(tresc, "Rysunek 1") {
		t.Fatalf("treść nie niesie podpisu: %q", tresc)
	}
	if !strings.HasPrefix(tresc, "PRZED") || !strings.HasSuffix(tresc, "|PO") {
		t.Fatalf("osadzenie wstawiło zasób w innym miejscu niż wskazane: %q", tresc)
	}
}
