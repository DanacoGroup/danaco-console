package core

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Skutek modułu Studio: czy droga od pliku do treści dokumentu jest cała.
//
// Szkoda, którą ten plik ma wykluczyć, ma w tym produkcie postać znaną:
// odpowiedź `ok` przy pustym wyniku. Rozpoznanie pisma jest na nią szczególnie
// podatne — Tesseract kończy się powodzeniem także wtedy, gdy nie odczytał ani
// jednego słowa, więc koperta udana nie mówi nic o tym, czy Operator dostał
// tekst.
//
// Dlatego żaden sprawdzian tego pliku nie kończy się na `ok`. Każdy pyta o to,
// co ZOSTAŁO: czy pozycja jest w kolejce przy kolejnym odczycie, czy tekst
// niesie słowa z obrazu, czy korekta zmieniła treść przyjmowaną do edytora
// i czy dokument założony z cyfryzacji ma tę treść po ponownym otwarciu.

// obrazZeSlowami rysuje obraz o znanej treści i oddaje jego ścieżkę.
//
// Materiał sprawdzianu powstaje na miejscu, a nie leży w drzewie: plik binarny
// w repozytorium starzeje się bez śladu, a tu chodzi o to, żeby tekst na
// obrazie i tekst oczekiwany pochodziły z jednego zapisu.
func obrazZeSlowami(t *testing.T, tresc string) string {
	t.Helper()

	rysownik, err := exec.LookPath("magick")
	if err != nil {
		t.Skipf("brak programu magick — sprawdzian rozpoznania nie ma czym narysować materiału: %v", err)
	}
	if _, err := exec.LookPath("tesseract"); err != nil {
		t.Skipf("brak programu tesseract — nie ma czym rozpoznać materiału: %v", err)
	}

	sciezka := filepath.Join(t.TempDir(), "skan.png")
	polecenie := exec.Command(rysownik, "-background", "white", "-fill", "black",
		"-pointsize", "48", "-density", "300", "label:"+tresc, sciezka)
	if wyjscie, err := polecenie.CombinedOutput(); err != nil {
		t.Skipf("nie udało się narysować materiału sprawdzianu: %v (%s)", err, wyjscie)
	}
	if dane, err := os.Stat(sciezka); err != nil || dane.Size() == 0 {
		t.Skipf("materiał sprawdzianu nie powstał albo jest pusty")
	}
	return sciezka
}

// dolozPozycje wnosi materiał do kolejki i oddaje jego pozycję.
func dolozPozycje(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno, sciezka string) shared.StudioIngestItem {
	t.Helper()

	var wynik shared.StudioIngestQueueAddResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestQueueAdd,
		shared.StudioIngestQueueAddRequest{WindowId: okno, SourcePaths: []string{sciezka}}, &wynik)
	if len(wynik.Items) != 1 {
		t.Fatalf("dołożenie jednego materiału dało %d pozycji", len(wynik.Items))
	}
	return wynik.Items[0]
}

// TestKolejkaCyfryzacjiZostajePoOdpowiedzi sprawdza wejście do rodziny: pozycja
// dołożona ma być w kolejce przy NASTĘPNYM odczycie, a nie tylko w odpowiedzi
// na dołożenie. Bez tego reszta sprawdzianów pytałaby o byt, którego nie ma.
func TestKolejkaCyfryzacjiZostajePoOdpowiedzi(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	sciezka := filepath.Join(katalog, "material.txt")
	if err := os.WriteFile(sciezka, []byte("treść materiału"), 0o644); err != nil {
		t.Fatalf("nie można założyć materiału: %v", err)
	}
	pozycja := dolozPozycje(t, zmontowany, zycie, "okno-studio-1", sciezka)

	if pozycja.State != shared.StudioIngestStateOczekuje {
		t.Errorf("pozycja świeżo dołożona ma stan %q, oczekiwano „oczekuje”", pozycja.State)
	}

	var wykaz shared.StudioIngestQueueListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestQueueList,
		shared.StudioIngestQueueListRequest{WindowId: "okno-studio-1"}, &wykaz)

	if len(wykaz.Items) != 1 || wykaz.Items[0].Id != pozycja.Id {
		t.Fatalf("kolejka po dołożeniu niesie %d pozycji, oczekiwano pozycji %s", len(wykaz.Items), pozycja.Id)
	}
	if wykaz.Items[0].SourcePath == nil || *wykaz.Items[0].SourcePath != sciezka {
		t.Errorf("pozycja w kolejce wskazuje %v, materiał leży pod %s", wykaz.Items[0].SourcePath, sciezka)
	}

	// Kolejka innego okna jest inną kolejką — bez tego wsad jednego Operatora
	// wchodziłby w pracę drugiego.
	var obce shared.StudioIngestQueueListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestQueueList,
		shared.StudioIngestQueueListRequest{WindowId: "okno-studio-2"}, &obce)
	if len(obce.Items) != 0 {
		t.Errorf("kolejka obcego okna niesie %d pozycji, oczekiwano pustki", len(obce.Items))
	}
}

// TestRozpoznanieOddajeTekstZObrazuWrazZeSlowami sprawdza rzecz, dla której ta
// rodzina istnieje: czy z obrazu wychodzi tekst. Miarą jest treść, którą
// narysowano na obrazie — nie kod odpowiedzi i nie długość napisu.
func TestRozpoznanieOddajeTekstZObrazuWrazZeSlowami(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	napis := "UMOWA RAMOWA"
	sciezka := obrazZeSlowami(t, napis)
	pozycja := dolozPozycje(t, zmontowany, zycie, "okno-studio-1", sciezka)

	var rozpoznanie shared.StudioIngestRecognizeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestRecognize,
		shared.StudioIngestRecognizeRequest{ItemId: pozycja.Id}, &rozpoznanie)

	if rozpoznanie.Item.Text == nil || strings.TrimSpace(*rozpoznanie.Item.Text) == "" {
		t.Fatal("rozpoznanie oddało odpowiedź udaną i tekst PUSTY — to jest dokładnie ten stan, " +
			"którego koperta nie odróżnia od odczytu")
	}
	odczyt := strings.ToUpper(*rozpoznanie.Item.Text)
	if !strings.Contains(odczyt, "UMOWA") {
		t.Errorf("na obrazie narysowano %q, rozpoznanie oddało %q", napis, *rozpoznanie.Item.Text)
	}

	// Słowa wraz z pewnością są warunkiem korekty rozpoznania: bez nich nie ma
	// czego poprawiać, a panel nie ma czego nałożyć na skan.
	if len(rozpoznanie.Words) == 0 {
		t.Fatal("rozpoznanie nie oddało ani jednego słowa — korekta rozpoznania nie miałaby czego poprawić")
	}
	for _, slowo := range rozpoznanie.Words {
		if slowo.Confidence <= 0 {
			t.Errorf("słowo %q niesie pewność %v — pewność niedodatnia znaczy odczyt bez pokrycia",
				slowo.Text, slowo.Confidence)
		}
		if slowo.Width <= 0 || slowo.Height <= 0 {
			t.Errorf("słowo %q niesie ramkę %dx%d — bez ramki nie ma czego nałożyć na skan",
				slowo.Text, slowo.Width, slowo.Height)
		}
	}
	if rozpoznanie.Item.Confidence == nil {
		t.Error("pozycja nie niesie pewności mimo rozpoznanych słów")
	}
	if rozpoznanie.Item.UsedOcr == nil || !*rozpoznanie.Item.UsedOcr {
		t.Error("pozycja nie przyznaje się do rozpoznania pisma — Operator ma wiedzieć, czy tekst jest pewny")
	}
	if rozpoznanie.Item.State != shared.StudioIngestStateGotowa {
		t.Errorf("pozycja po udanym rozpoznaniu ma stan %q, oczekiwano „gotowa”", rozpoznanie.Item.State)
	}
}

// TestKorektaRozpoznaniaZmieniaTrescPrzyjmowanaDoEdytora sprawdza, że poprawka
// słowa nie zostaje na warstwie tekstowej: ma zmienić tekst, który Operator za
// chwilę przyjmie. Poprawka widoczna i bezskuteczna byłaby gorsza od jej braku.
func TestKorektaRozpoznaniaZmieniaTrescPrzyjmowanaDoEdytora(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	sciezka := obrazZeSlowami(t, "UMOWA RAMOWA")
	pozycja := dolozPozycje(t, zmontowany, zycie, "okno-studio-1", sciezka)

	var rozpoznanie shared.StudioIngestRecognizeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestRecognize,
		shared.StudioIngestRecognizeRequest{ItemId: pozycja.Id}, &rozpoznanie)
	if len(rozpoznanie.Words) == 0 {
		t.Skip("rozpoznanie nie oddało słów — nie ma czego poprawiać")
	}

	przedKorekta := *rozpoznanie.Item.Text
	var korekta shared.StudioIngestCorrectionSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestCorrectionSet,
		shared.StudioIngestCorrectionSetRequest{
			ItemId: pozycja.Id, WordIndex: rozpoznanie.Words[0].Index, Text: "POROZUMIENIE",
		}, &korekta)

	if korekta.Item.Text == nil {
		t.Fatal("pozycja po korekcie nie niesie tekstu")
	}
	if *korekta.Item.Text == przedKorekta {
		t.Fatalf("korekta nie zmieniła tekstu pozycji: nadal %q", przedKorekta)
	}
	if !strings.Contains(*korekta.Item.Text, "POROZUMIENIE") {
		t.Errorf("tekst po korekcie nie niesie poprawionego słowa: %q", *korekta.Item.Text)
	}

	// Poprawka ma przeżyć odświeżenie okna — inaczej znika przy pierwszym
	// odczycie kolejki, czyli dokładnie wtedy, gdy Operator na nią patrzy.
	var wykaz shared.StudioIngestQueueListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestQueueList,
		shared.StudioIngestQueueListRequest{WindowId: "okno-studio-1"}, &wykaz)
	if len(wykaz.Items) != 1 || wykaz.Items[0].Text == nil ||
		!strings.Contains(*wykaz.Items[0].Text, "POROZUMIENIE") {
		t.Error("korekta nie przeżyła ponownego odczytu kolejki")
	}
}

// TestPrzyjeciePozycjiZakladaDokumentZTrescia zamyka drogę: czy to, co wyszło
// z cyfryzacji, jest w dokumencie i czy dokument da się z niego odtworzyć.
// Miarą jest odczyt dokumentu OSOBNYM wywołaniem — odpowiedź na przyjęcie
// mogłaby nieść treść, której baza nie przyjęła.
//
// Materiałem jest obraz, a nie plik tekstowy: rozpoznanie pisma czyta piksele,
// więc plik tekstowy podany jako materiał kończy się odmową i sprawdzian
// pomijałby się zawsze. Sprawdzian, który zawsze się pomija, niczego nie
// pilnuje.
func TestPrzyjeciePozycjiZakladaDokumentZTrescia(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	sciezka := obrazZeSlowami(t, "PROTOKOL ODBIORU")
	pozycja := dolozPozycje(t, zmontowany, zycie, "okno-studio-1", sciezka)

	var rozpoznanie shared.StudioIngestRecognizeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestRecognize,
		shared.StudioIngestRecognizeRequest{ItemId: pozycja.Id}, &rozpoznanie)
	if rozpoznanie.Item.Text == nil || strings.TrimSpace(*rozpoznanie.Item.Text) == "" {
		t.Fatal("rozpoznanie nie oddało tekstu — przyjęcie nie miałoby czego przenieść")
	}
	rozpoznanyTekst := *rozpoznanie.Item.Text

	var przyjecie shared.StudioIngestItemAcceptResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioIngestItemAccept,
		shared.StudioIngestItemAcceptRequest{
			WindowId: "okno-studio-1", ItemIds: []string{pozycja.Id},
		}, &przyjecie)

	if przyjecie.Document.Content == nil || *przyjecie.Document.Content == "" {
		t.Fatal("przyjęcie pozycji założyło dokument BEZ treści — koperta udana, wynik pusty")
	}
	if *przyjecie.Document.Content != rozpoznanyTekst {
		t.Errorf("dokument niesie %q, cyfryzacja oddała %q", *przyjecie.Document.Content, rozpoznanyTekst)
	}
	if przyjecie.Version.Id == "" {
		t.Error("przyjęcie nie założyło wersji pierwszej — historia dokumentu zaczyna się od pustki")
	}
	if przyjecie.Document.VersionId == nil || *przyjecie.Document.VersionId != przyjecie.Version.Id {
		t.Errorf("dokument wskazuje wersję bieżącą %v, a założono %s",
			przyjecie.Document.VersionId, przyjecie.Version.Id)
	}

	// Autor wersji rozstrzyga rozdział 3.6 opracowania: wersja z cyfryzacji
	// pochodzi od modelu, nie od Operatora. Kolumna autora powstała po to.
	if przyjecie.Version.Author == nil || *przyjecie.Version.Author != shared.StudioAuthorModel {
		t.Errorf("wersja z cyfryzacji niesie autora %v, oczekiwano %q",
			przyjecie.Version.Author, shared.StudioAuthorModel)
	}

	// Dokument odczytany osobno ma nieść tę samą treść. To jest miara właściwa:
	// Operator otworzy go ponownie, a nie przeczyta odpowiedź na przyjęcie.
	var otwarty shared.StudioDocumentOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentOpen,
		shared.StudioDocumentOpenRequest{
			WindowId: "okno-studio-1", DocumentId: &przyjecie.Document.Id,
		}, &otwarty)
	if otwarty.Document.Content == nil || *otwarty.Document.Content != rozpoznanyTekst {
		t.Errorf("dokument po ponownym otwarciu niesie %v, cyfryzacja oddała %q",
			otwarty.Document.Content, rozpoznanyTekst)
	}

	// Historia dokumentu ma tę wersję zawierać — inaczej Session Repository
	// pokaże pustkę tuż po cyfryzacji.
	var historia shared.StudioRepositoryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioRepositoryList,
		shared.StudioRepositoryListRequest{DocumentId: przyjecie.Document.Id}, &historia)
	if len(historia.Versions) != 1 || historia.Versions[0].Id != przyjecie.Version.Id {
		t.Errorf("historia niesie %d wersji, oczekiwano wersji %s", len(historia.Versions), przyjecie.Version.Id)
	}
}

// TestRozpoznanieNieistniejacejPozycjiOdmawiaNazywajacBrak pilnuje, żeby
// wskazanie bytu, którego nie ma, nie wracało odpowiedzią udaną z pustką.
func TestRozpoznanieNieistniejacejPozycjiOdmawiaNazywajacBrak(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandStudioIngestRecognize,
		shared.StudioIngestRecognizeRequest{ItemId: "studio-wcz-nie-ma-takiej"})
	if odpowiedz.Error == nil {
		t.Fatal("rozpoznanie pozycji nieistniejącej wróciło odpowiedzią udaną")
	}
	if odpowiedz.Error.Code != shared.ErrorCodeNotFound {
		t.Errorf("odmowa niesie kod %q, oczekiwano %q", odpowiedz.Error.Code, shared.ErrorCodeNotFound)
	}
	if !strings.Contains(odpowiedz.Error.Message, "studio-wcz-nie-ma-takiej") {
		t.Errorf("odmowa nie nazywa wskazanego bytu: %q", odpowiedz.Error.Message)
	}
}
