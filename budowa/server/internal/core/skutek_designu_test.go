package core

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"image"
	"image/png"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	_ "modernc.org/sqlite"

	"danacoconsole/server/internal/transport"
	"danacoconsole/shared"
)

// Skutek modułu Design: czy za odpowiedzią naprawdę coś zostaje.
//
// ── Dlaczego akurat ten moduł mierzy się tak surowo ─────────────────────────
// To jest moduł z precedensem szkody. `design.asset.list` meldował kiedyś
// `status: ok` z wykazem zasobów, za którymi nie było ANI JEDNEGO BAJTU.
// Koperta była udana i zarazem kłamała, a klient czyta kopertę, nie komentarz
// w kodzie. Od tamtej pory obowiązuje w produkcie zasada, której ten plik jest
// wzorcowym zastosowaniem:
//
//	Sprawdzian schodzi po odwołaniu do magazynu albo do bazy i mierzy
//	NIEZALEŻNIE. Nigdy nie kończy się na treści odpowiedzi.
//
// Stąd dwa narzędzia pomiaru, oba omijające rdzeń:
//
//  1. Wydanie zasobu jest DEKODOWANE Z POWROTEM JAKO OBRAZ i mierzone co do
//     wymiarów po skali. Odpowiedź mówiąca „oto png @2x" nie jest dowodem;
//     dowodem jest `image.Decode`, któremu te bajty wystarczą, i szerokość,
//     która wyszła dwa razy większa.
//
//  2. Kolekcja, wersja kompozycji, adnotacja i zestaw żetonów są odczytywane
//     DRUGIM, NIEZALEŻNYM POŁĄCZENIEM SQLite do pliku bazy stanowiska. Rdzeń
//     nie bierze w tym odczycie udziału: gdyby zapis nie doszedł do pliku,
//     odpowiedź komendy i tak wyglądałaby tak samo.
//
// Sumę kontrolną treści sprawdzian liczy z BAJTÓW, nie przepisuje jej
// z odpowiedzi — zgodność tych dwóch wartości jest tu mierzona, a nie założona.

// ── Uprząż pomiaru ──────────────────────────────────────────────────────────

// polaczenieOboczneDesignu otwiera drugie połączenie do pliku bazy stanowiska.
// Osobne połączenie, nie pula rdzenia: pomiar ma powiedzieć, co leży W PLIKU,
// a nie co pamięta proces, który to zapisywał.
func polaczenieOboczneDesignu(t *testing.T, katalogDanych string) *sql.DB {
	t.Helper()

	sciezka := filepath.ToSlash(filepath.Join(katalogDanych, "dane.sqlite"))
	oboczne, err := sql.Open("sqlite", sciezka+"?_pragma=busy_timeout(5000)")
	if err != nil {
		t.Fatalf("nie można otworzyć obocznego połączenia do bazy: %v", err)
	}
	t.Cleanup(func() { _ = oboczne.Close() })
	if err := oboczne.Ping(); err != nil {
		t.Fatalf("oboczne połączenie do bazy nie odpowiada: %v", err)
	}
	return oboczne
}

// liczbaObocznaDesignu odczytuje jedną liczbę zapytaniem na obocznym
// połączeniu.
func liczbaObocznaDesignu(t *testing.T, oboczne *sql.DB, zapytanie string, argumenty ...any) int {
	t.Helper()

	var wynik int
	if err := oboczne.QueryRow(zapytanie, argumenty...).Scan(&wynik); err != nil {
		t.Fatalf("oboczny odczyt liczby nie powiódł się (%s): %v", zapytanie, err)
	}
	return wynik
}

// tekstObocznyDesignu odczytuje jeden tekst zapytaniem na obocznym połączeniu.
func tekstObocznyDesignu(t *testing.T, oboczne *sql.DB, zapytanie string, argumenty ...any) string {
	t.Helper()

	var wynik sql.NullString
	if err := oboczne.QueryRow(zapytanie, argumenty...).Scan(&wynik); err != nil {
		t.Fatalf("oboczny odczyt tekstu nie powiódł się (%s): %v", zapytanie, err)
	}
	return wynik.String
}

// wniesObrazDesignu wnosi obraz o zadanych wymiarach i oddaje identyfikator
// zasobu wraz z bajtami, które poszły do rdzenia.
func wniesObrazDesignu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno string, szerokosc, wysokosc int) (string, []byte) {
	t.Helper()

	bajty := obrazPNG(t, szerokosc, wysokosc)
	var wynik shared.DesignAssetUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetUpload,
		shared.DesignAssetUploadRequest{
			WindowId:      okno,
			Name:          wskaznik("materiał sprawdzianu"),
			Kind:          shared.DesignAssetKindImage,
			ContentBase64: wskaznik(wBase64(bajty)),
		}, &wynik)
	if wynik.Asset.Id == "" {
		t.Fatal("wniesienie oddało zasób bez identyfikatora")
	}
	return wynik.Asset.Id, bajty
}

// wymiaryWydaniaDesignu rozkłada treść wydania z powrotem na obraz i oddaje
// jego wymiary. To jest właściwy pomiar wydania: bajty, których `image.Decode`
// nie przyjmuje, nie są obrazem, choćby odpowiedź nazywała je `image/png`.
func wymiaryWydaniaDesignu(t *testing.T, trescBase64 string) (int, int) {
	t.Helper()

	bajty, err := base64.StdEncoding.DecodeString(trescBase64)
	if err != nil {
		t.Fatalf("treść wydania nie jest poprawnym base64: %v", err)
	}
	if len(bajty) == 0 {
		t.Fatal("treść wydania jest pusta — odpowiedź udana bez ani jednego bajtu")
	}
	opis, _, err := image.DecodeConfig(bytes.NewReader(bajty))
	if err != nil {
		t.Fatalf("treści wydania nie da się rozłożyć jako obrazu (%d bajtów): %v", len(bajty), err)
	}
	return opis.Width, opis.Height
}

// zycieZGniazdemDesignu dokłada do kontekstu tożsamość połączenia — taką, jaką
// nadaje transport oknu, które się przywitało.
//
// Zgłoszenie obecności bez rozpoznanego klienta jest odmawiane, i słusznie:
// kursor bez tożsamości nie da się odróżnić od cudzego ani zdjąć przy odejściu.
// Uprząż woła rdzeń z pominięciem gniazda, więc tożsamość trzeba tu podstawić —
// inaczej sprawdzian mierzyłby wywołanie, które w produkcie nie zachodzi.
func zycieZGniazdemDesignu(zycie context.Context, klient string) context.Context {
	return zPolaczeniem(zycie, transport.Tozsamosc{
		IdPolaczenia: "gniazdo-sprawdzianu",
		IdKlienta:    klient,
	})
}

// stronDokumentuDesignu liczy strony dokumentu z jego BAJTÓW. Wydanie pdf,
// którego biblioteka nie umie otworzyć, nie jest dokumentem — choćby zaczynało
// się właściwym nagłówkiem.
func stronDokumentuDesignu(t *testing.T, bajty []byte) int {
	t.Helper()

	liczba, err := api.PageCount(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		t.Fatalf("wydania pdf nie da się otworzyć jako dokumentu (%d bajtów): %v", len(bajty), err)
	}
	return liczba
}

// ── Treść zasobu ────────────────────────────────────────────────────────────

// TestTrescZasobuDesignuOddajeBajtyZgodneZSumaZBajtow sprawdza to, czego
// precedens szkody nie sprawdzał: czy za zasobem leży treść, i czy suma
// kontrolna z odpowiedzi zgadza się z sumą policzoną z bajtów.
func TestTrescZasobuDesignuOddajeBajtyZgodneZSumaZBajtow(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	kod, zrodlo := wniesObrazDesignu(t, zmontowany, zycie, "okno-tresci", 40, 24)

	var wynik shared.DesignAssetContentGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetContentGet,
		shared.DesignAssetContentGetRequest{AssetId: kod}, &wynik)

	if wynik.ContentBase64 == nil || *wynik.ContentBase64 == "" {
		t.Fatal("komenda oddała zasób bez treści — dokładnie ta szkoda, której moduł ma precedens")
	}
	bajty, err := base64.StdEncoding.DecodeString(*wynik.ContentBase64)
	if err != nil {
		t.Fatalf("treść zasobu nie jest poprawnym base64: %v", err)
	}
	if !bytes.Equal(bajty, zrodlo) {
		t.Errorf("treść oddana (%d bajtów) różni się od wniesionej (%d bajtów)",
			len(bajty), len(zrodlo))
	}
	if wynik.SizeBytes != len(bajty) {
		t.Errorf("odpowiedź podaje %d bajtów, a treść ma %d", wynik.SizeBytes, len(bajty))
	}
	// Suma liczona TU, z bajtów — nie przepisana z odpowiedzi.
	if zmierzona := sumaSha256(bajty); wynik.Checksum != zmierzona {
		t.Errorf("suma kontrolna odpowiedzi %q nie zgadza się z policzoną z bajtów %q",
			wynik.Checksum, zmierzona)
	}
	if _, _, err := image.DecodeConfig(bytes.NewReader(bajty)); err != nil {
		t.Errorf("oddana treść nie rozkłada się jako obraz: %v", err)
	}
}

// TestTrescZasobuDesignuOdmawiaPonadGranicaZamiastUcinac pilnuje, że
// przekroczenie granicy wołającego kończy się odmową NAZYWAJĄCĄ ZMIERZONĄ
// WIELKOŚĆ, a nie treścią uciętą podaną jako całość.
func TestTrescZasobuDesignuOdmawiaPonadGranicaZamiastUcinac(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	kod, zrodlo := wniesObrazDesignu(t, zmontowany, zycie, "okno-granicy", 64, 64)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignAssetContentGet,
		shared.DesignAssetContentGetRequest{AssetId: kod, MaxBytes: wskaznik(len(zrodlo) - 1)})
	if !strings.Contains(odmowa.Message, "bajtów") {
		t.Errorf("odmowa nie nazywa zmierzonej wielkości: %q", odmowa.Message)
	}
}

// TestTrescZasobuDesignuOdmawiaZasobuNieznanego pilnuje, że rdzeń nie oddaje
// bajtów zastępczych za zasób, którego nie zna.
func TestTrescZasobuDesignuOdmawiaZasobuNieznanego(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignAssetContentGet,
		shared.DesignAssetContentGetRequest{AssetId: "zasob-ktorego-nie-ma"})
	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Errorf("odmowa za zasób nieznany ma kod %q, a ma mieć %q",
			odmowa.Code, shared.ErrorCodeNotFound)
	}
}

// ── Wydania zasobu ──────────────────────────────────────────────────────────

// TestWydanieZasobuDesignuMaZmierzoneWymiaryPoSkali jest sprawdzianem
// wzorcowym: wydanie wraca, zostaje rozłożone z powrotem na obraz i zmierzone.
func TestWydanieZasobuDesignuMaZmierzoneWymiaryPoSkali(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	kod, _ := wniesObrazDesignu(t, zmontowany, zycie, "okno-wydania", 40, 20)

	var naturalne shared.DesignAssetExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetExport,
		shared.DesignAssetExportRequest{AssetId: kod, Format: "png"}, &naturalne)
	if szerokosc, wysokosc := wymiaryWydaniaDesignu(t, naturalne.ContentBase64); szerokosc != 40 ||
		wysokosc != 20 {
		t.Errorf("wydanie w skali naturalnej ma %d×%d, a materiał 40×20", szerokosc, wysokosc)
	}
	if naturalne.MediaType != "image/png" {
		t.Errorf("wydanie png podaje typ treści %q", naturalne.MediaType)
	}
	if !strings.HasSuffix(naturalne.FileName, ".png") {
		t.Errorf("proponowana nazwa %q nie kończy się rozszerzeniem formatu", naturalne.FileName)
	}

	var podwojne shared.DesignAssetExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetExport,
		shared.DesignAssetExportRequest{AssetId: kod, Format: "png", Scale: wskaznik(2.0)}, &podwojne)
	szerokosc, wysokosc := wymiaryWydaniaDesignu(t, podwojne.ContentBase64)
	if szerokosc != 80 || wysokosc != 40 {
		t.Errorf("wydanie @2x ma %d×%d, a ma mieć 80×40", szerokosc, wysokosc)
	}
	if !strings.Contains(podwojne.FileName, "@2x") {
		t.Errorf("nazwa wydania gęstości %q nie niesie wariantu", podwojne.FileName)
	}
	if podwojne.SizeBytes == 0 {
		t.Error("wydanie melduje zerową wielkość")
	}
}

// TestWydanieZasobuDesignuSkladaJpegIkoneIDokument sprawdza pozostałe formaty
// wydania biblioteką wkompilowaną — każdy mierzony po treści, nie po nazwie.
func TestWydanieZasobuDesignuSkladaJpegIkoneIDokument(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	kod, _ := wniesObrazDesignu(t, zmontowany, zycie, "okno-formatow", 32, 32)

	var jpeg shared.DesignAssetExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetExport,
		shared.DesignAssetExportRequest{AssetId: kod, Format: "jpg", Quality: wskaznik(70)}, &jpeg)
	if szerokosc, wysokosc := wymiaryWydaniaDesignu(t, jpeg.ContentBase64); szerokosc != 32 ||
		wysokosc != 32 {
		t.Errorf("wydanie jpeg ma %d×%d, a materiał 32×32", szerokosc, wysokosc)
	}
	if jpeg.MediaType != "image/jpeg" {
		t.Errorf("wydanie jpeg podaje typ treści %q", jpeg.MediaType)
	}

	var ikona shared.DesignAssetExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetExport,
		shared.DesignAssetExportRequest{AssetId: kod, Format: "ico"}, &ikona)
	bajtyIkony, err := base64.StdEncoding.DecodeString(ikona.ContentBase64)
	if err != nil {
		t.Fatalf("treść ikony nie jest poprawnym base64: %v", err)
	}
	// Nagłówek ikony: pole zastrzeżone 0, rodzaj 1, liczba wpisów 1.
	if len(bajtyIkony) < 22 || bajtyIkony[0] != 0 || bajtyIkony[1] != 0 ||
		bajtyIkony[2] != 1 || bajtyIkony[3] != 0 {
		t.Fatalf("wydanie ico nie ma nagłówka ikony (%d bajtów)", len(bajtyIkony))
	}
	// Zawartość wpisu leży od bajtu 22 i jest obrazem PNG — rozkładamy ją,
	// zamiast wierzyć nagłówkowi.
	if _, err := png.Decode(bytes.NewReader(bajtyIkony[22:])); err != nil {
		t.Errorf("zawartość ikony nie rozkłada się jako obraz: %v", err)
	}

	var dokument shared.DesignAssetExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetExport,
		shared.DesignAssetExportRequest{AssetId: kod, Format: "pdf"}, &dokument)
	bajtyDokumentu, err := base64.StdEncoding.DecodeString(dokument.ContentBase64)
	if err != nil {
		t.Fatalf("treść dokumentu nie jest poprawnym base64: %v", err)
	}
	if !bytes.HasPrefix(bajtyDokumentu, []byte("%PDF-")) {
		t.Error("wydanie pdf nie zaczyna się nagłówkiem dokumentu")
	}
	if stron := stronDokumentuDesignu(t, bajtyDokumentu); stron != 1 {
		t.Errorf("wydanie pdf ma %d stron, a ma mieć jedną", stron)
	}
}

// TestWydanieZasobuDesignuOdmawiaFormatuBezBiblioteki pilnuje zasady
// bezwzględnej produktu: format, którego nie da się złożyć biblioteką
// wkompilowaną, kończy się ODMOWĄ WYMIENIAJĄCĄ formaty obsługiwane — nigdy
// cichym PNG pod cudzą nazwą.
func TestWydanieZasobuDesignuOdmawiaFormatuBezBiblioteki(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	kod, _ := wniesObrazDesignu(t, zmontowany, zycie, "okno-odmowy", 16, 16)

	for _, format := range []string{"avif", "webp"} {
		odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignAssetExport,
			shared.DesignAssetExportRequest{AssetId: kod, Format: format})
		if !strings.Contains(odmowa.Message, formatyWydaniaDesignu) {
			t.Errorf("odmowa formatu %s nie wymienia formatów obsługiwanych: %q", format, odmowa.Message)
		}
	}
}

// TestWydaniePartiaDesignuNiesieBilansZamiastCiszy sprawdza, że odrzucenie
// jednego zasobu nie wstrzymuje pozostałych i że wykaz odrzuconych naprawdę
// wraca — krótsza lista bez słowa byłaby ciszą.
func TestWydaniePartiaDesignuNiesieBilansZamiastCiszy(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	pierwszy, _ := wniesObrazDesignu(t, zmontowany, zycie, "okno-partii", 20, 10)
	drugi, _ := wniesObrazDesignu(t, zmontowany, zycie, "okno-partii", 30, 15)

	var wynik shared.DesignAssetExportBatchResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetExportBatch,
		shared.DesignAssetExportBatchRequest{
			AssetIds: []string{pierwszy, "zasob-widmo", drugi},
			Format:   "png",
			Scales:   []float64{1, 2},
		}, &wynik)

	if wynik.Exported != 2 {
		t.Errorf("partia melduje %d wydanych zasobów, a wydać się miały dwa", wynik.Exported)
	}
	if len(wynik.Files) != 4 {
		t.Fatalf("partia oddała %d wydań, a dwa zasoby w dwóch skalach dają cztery", len(wynik.Files))
	}
	if len(wynik.Rejected) != 1 || wynik.Rejected[0].AssetId != "zasob-widmo" {
		t.Errorf("partia nie zameldowała odrzucenia zasobu widma: %+v", wynik.Rejected)
	}
	if wynik.Rejected[0].Reason == "" {
		t.Error("odrzucenie w partii nie niesie powodu")
	}
	// Każde wydanie rozkładamy z powrotem: partia melduje cztery pliki, więc
	// cztery mają być obrazami.
	for _, plik := range wynik.Files {
		szerokosc, wysokosc := wymiaryWydaniaDesignu(t, plik.ContentBase64)
		if plik.Scale == nil {
			t.Fatalf("wydanie %s nie niesie krotności skali", plik.FileName)
		}
		oczekiwana := map[string][2]int{
			pierwszy: {20, 10},
			drugi:    {30, 15},
		}[plik.AssetId]
		chcianaSzerokosc := int(float64(oczekiwana[0])**plik.Scale + 0.5)
		chcianaWysokosc := int(float64(oczekiwana[1])**plik.Scale + 0.5)
		if szerokosc != chcianaSzerokosc || wysokosc != chcianaWysokosc {
			t.Errorf("wydanie %s w skali %v ma %d×%d, a ma mieć %d×%d",
				plik.FileName, *plik.Scale, szerokosc, wysokosc, chcianaSzerokosc, chcianaWysokosc)
		}
	}
}

// ── Kolekcje ────────────────────────────────────────────────────────────────

// TestKolekcjaDesignuLezyWBazieStanowiska mierzy kolekcję drugim połączeniem
// do pliku bazy: odpowiedź komendy nie jest tu dowodem na nic.
func TestKolekcjaDesignuLezyWBazieStanowiska(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneDesignu(t, katalog)
	pierwszy, _ := wniesObrazDesignu(t, zmontowany, zycie, "okno-kolekcji", 8, 8)
	drugi, _ := wniesObrazDesignu(t, zmontowany, zycie, "okno-kolekcji", 8, 8)

	var zalozona shared.DesignCollectionCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignCollectionCreate,
		shared.DesignCollectionCreateRequest{
			WindowId:    "okno-kolekcji",
			Name:        "Kampania Q3",
			Description: wskaznik("zestaw do arkusza brandingowego"),
		}, &zalozona)

	if nazwa := tekstObocznyDesignu(t, oboczne,
		`SELECT nazwa FROM kolekcja_design WHERE identyfikator_zewnetrzny = ?`,
		zalozona.Collection.Id); nazwa != "Kampania Q3" {
		t.Errorf("w bazie stanowiska kolekcja nazywa się %q, a założono „Kampania Q3”", nazwa)
	}

	var przypisana shared.DesignCollectionAssignResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignCollectionAssign,
		shared.DesignCollectionAssignRequest{
			CollectionId: zalozona.Collection.Id,
			AssetIds:     []string{pierwszy, drugi},
		}, &przypisana)
	if przypisana.Changed != 2 {
		t.Errorf("przypisanie melduje %d zmian, a dołożono dwa zasoby", przypisana.Changed)
	}
	if przypisana.Collection.AssetCount != 2 {
		t.Errorf("kolekcja melduje %d zasobów, a ma dwa", przypisana.Collection.AssetCount)
	}

	// Pomiar niezależny: ile pozycji leży w pliku bazy.
	pozycji := liczbaObocznaDesignu(t, oboczne,
		`SELECT COUNT(*) FROM pozycja_kolekcji_design p
		   JOIN kolekcja_design k ON k.id = p.kolekcja_id
		  WHERE k.identyfikator_zewnetrzny = ?`, zalozona.Collection.Id)
	if pozycji != 2 {
		t.Errorf("w bazie stanowiska leży %d przypisań, a odpowiedź meldowała dwa", pozycji)
	}

	// Powtórzone dołożenie nie zmienia niczego — `changed` ma to powiedzieć.
	var powtorzona shared.DesignCollectionAssignResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignCollectionAssign,
		shared.DesignCollectionAssignRequest{
			CollectionId: zalozona.Collection.Id,
			AssetIds:     []string{pierwszy},
		}, &powtorzona)
	if powtorzona.Changed != 0 {
		t.Errorf("powtórzone dołożenie melduje %d zmian, a nie zmieniło niczego", powtorzona.Changed)
	}

	// Zdjęcie zasobu z kolekcji musi zejść do pliku bazy.
	var zdjeta shared.DesignCollectionAssignResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignCollectionAssign,
		shared.DesignCollectionAssignRequest{
			CollectionId: zalozona.Collection.Id,
			AssetIds:     []string{drugi},
			Remove:       wskaznik(true),
		}, &zdjeta)
	poZdjeciu := liczbaObocznaDesignu(t, oboczne,
		`SELECT COUNT(*) FROM pozycja_kolekcji_design p
		   JOIN kolekcja_design k ON k.id = p.kolekcja_id
		  WHERE k.identyfikator_zewnetrzny = ?`, zalozona.Collection.Id)
	if poZdjeciu != 1 {
		t.Errorf("po zdjęciu w bazie leży %d przypisań, a ma leżeć jedno", poZdjeciu)
	}

	var wykaz shared.DesignCollectionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignCollectionList,
		shared.DesignCollectionListRequest{WindowId: "okno-kolekcji", AssetId: &pierwszy}, &wykaz)
	if wykaz.Total != 1 {
		t.Errorf("zawężenie do kolekcji z zasobem oddało %d kolekcji, a ma jedną", wykaz.Total)
	}
}

// TestKolekcjaDesignuOdmawiaZasobuWidma pilnuje, żeby kolekcja nie zapełniła
// się kodami bez treści — czyli tą samą szkodą, którą moduł ma w historii.
func TestKolekcjaDesignuOdmawiaZasobuWidma(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneDesignu(t, katalog)

	var zalozona shared.DesignCollectionCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignCollectionCreate,
		shared.DesignCollectionCreateRequest{WindowId: "okno-widma", Name: "Widma"}, &zalozona)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignCollectionAssign,
		shared.DesignCollectionAssignRequest{
			CollectionId: zalozona.Collection.Id,
			AssetIds:     []string{"zasob-widmo"},
		})
	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Errorf("odmowa za zasób widmo ma kod %q, a ma mieć %q", odmowa.Code, shared.ErrorCodeNotFound)
	}
	if pozycji := liczbaObocznaDesignu(t, oboczne,
		`SELECT COUNT(*) FROM pozycja_kolekcji_design`); pozycji != 0 {
		t.Errorf("po odmowie w bazie leży %d przypisań, a nie ma leżeć żadne", pozycji)
	}
}

// ── Szablony i historia promptów ────────────────────────────────────────────

// TestSzablonPromptuDesignuLezyWBazieStanowiska mierzy szablon obocznym
// połączeniem i sprawdza, że nadpisanie nie zakłada drugiego wiersza.
func TestSzablonPromptuDesignuLezyWBazieStanowiska(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneDesignu(t, katalog)

	var zapisany shared.DesignPromptTemplateSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPromptTemplateSave,
		shared.DesignPromptTemplateSaveRequest{
			WindowId: "okno-szablonow",
			Name:     "Bohater strony",
			Prompt: shared.DesignPrompt{
				Subject:     "ilustracja bohatera strony",
				Style:       wskaznik("złoty, płaski"),
				AspectRatio: wskaznik("16:9"),
				Variants:    wskaznik(4),
			},
		}, &zapisany)

	if temat := tekstObocznyDesignu(t, oboczne,
		`SELECT temat FROM szablon_promptu_design WHERE identyfikator_zewnetrzny = ?`,
		zapisany.Template.Id); temat != "ilustracja bohatera strony" {
		t.Errorf("w bazie stanowiska szablon niesie temat %q", temat)
	}

	var nadpisany shared.DesignPromptTemplateSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPromptTemplateSave,
		shared.DesignPromptTemplateSaveRequest{
			WindowId:   "okno-szablonow",
			Name:       "Bohater strony — wersja druga",
			TemplateId: &zapisany.Template.Id,
			Prompt:     shared.DesignPrompt{Subject: "ilustracja bohatera, wariant chłodny"},
		}, &nadpisany)

	if nadpisany.Template.Id != zapisany.Template.Id {
		t.Errorf("nadpisanie oddało szablon %q zamiast %q", nadpisany.Template.Id, zapisany.Template.Id)
	}
	if ile := liczbaObocznaDesignu(t, oboczne,
		`SELECT COUNT(*) FROM szablon_promptu_design WHERE okno = ?`, "okno-szablonow"); ile != 1 {
		t.Errorf("po nadpisaniu w bazie leży %d szablonów, a ma leżeć jeden", ile)
	}

	var wykaz shared.DesignPromptTemplateListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPromptTemplateList,
		shared.DesignPromptTemplateListRequest{WindowId: "okno-szablonow"}, &wykaz)
	if wykaz.Total != 1 || len(wykaz.Templates) != 1 {
		t.Fatalf("wykaz szablonów oddał %d pozycji, a ma jedną", wykaz.Total)
	}
	if wykaz.Templates[0].Prompt.Subject != "ilustracja bohatera, wariant chłodny" {
		t.Errorf("wykaz oddaje temat sprzed nadpisania: %q", wykaz.Templates[0].Prompt.Subject)
	}

	// Szablon spoza okna nie da się nadpisać cudzym kodem.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignPromptTemplateSave,
		shared.DesignPromptTemplateSaveRequest{
			WindowId:   "okno-obce",
			Name:       "podszycie",
			TemplateId: &zapisany.Template.Id,
			Prompt:     shared.DesignPrompt{Subject: "cokolwiek"},
		})
	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Errorf("nadpisanie szablonu z cudzego okna nie zostało odmówione: %q", odmowa.Code)
	}
}

// TestHistoriaPromptowDesignuOddajePustyWykazBezWydan pilnuje, że historia
// okna bez ani jednego wydania jest wykazem pustym, a nie odmową ani wykazem
// promptów cudzych okien.
func TestHistoriaPromptowDesignuOddajePustyWykazBezWydan(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var wykaz shared.DesignPromptHistoryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPromptHistoryList,
		shared.DesignPromptHistoryListRequest{WindowId: "okno-bez-wydan"}, &wykaz)
	if wykaz.Total != 0 || len(wykaz.Prompts) != 0 {
		t.Errorf("historia okna bez wydań oddała %d promptów", wykaz.Total)
	}
}

// ── Wersje kompozycji ───────────────────────────────────────────────────────

// zalozKompozycjeDesignu zakłada kompozycję z warstwą wskazującą zasób
// i oddaje jej identyfikator.
func zalozKompozycjeDesignu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno, zasob string) string {
	t.Helper()

	var wynik shared.DesignBoardUpdateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardUpdate,
		shared.DesignBoardUpdateRequest{
			WindowId: okno,
			Name:     wskaznik("Strona główna — wariant A"),
			Layers: []shared.DesignBoardLayer{{
				AssetId: &zasob,
				X:       wskaznik(0.0), Y: wskaznik(0.0),
				Width: wskaznik(40.0), Height: wskaznik(20.0),
			}},
		}, &wynik)
	return wynik.Board.Id
}

// TestWersjaKompozycjiDesignuLezyWBazieIWracaPrzyPrzywroceniu mierzy
// wersjonowanie obocznym połączeniem i sprawdza obietnicę kontraktu:
// przywrócenie zakłada wersję z układu sprzed przywrócenia.
func TestWersjaKompozycjiDesignuLezyWBazieIWracaPrzyPrzywroceniu(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneDesignu(t, katalog)
	zasob, _ := wniesObrazDesignu(t, zmontowany, zycie, "okno-wersji", 40, 20)
	plansza := zalozKompozycjeDesignu(t, zmontowany, zycie, "okno-wersji", zasob)

	var wersja shared.DesignBoardVersionSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardVersionSave,
		shared.DesignBoardVersionSaveRequest{
			BoardId: plansza,
			Name:    wskaznik("układ z jedną warstwą"),
			Note:    wskaznik("przed przebudową"),
		}, &wersja)
	if wersja.Version.LayerCount != 1 {
		t.Errorf("wersja melduje %d warstw, a układ miał jedną", wersja.Version.LayerCount)
	}

	// Pomiar niezależny: migawka warstw naprawdę leży w pliku bazy.
	warstwWersji := liczbaObocznaDesignu(t, oboczne,
		`SELECT COUNT(*) FROM warstwa_wersji_kompozycji_design w
		   JOIN wersja_kompozycji_design v ON v.id = w.wersja_id
		  WHERE v.identyfikator_zewnetrzny = ?`, wersja.Version.Id)
	if warstwWersji != 1 {
		t.Errorf("w bazie stanowiska migawka ma %d warstw, a ma mieć jedną", warstwWersji)
	}

	// Płótno zostaje wyczyszczone — stan poprawny, do którego wersja ma być
	// drogą powrotną.
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardUpdate,
		shared.DesignBoardUpdateRequest{WindowId: "okno-wersji", BoardId: &plansza,
			Layers: []shared.DesignBoardLayer{}}, nil)
	if warstw := liczbaObocznaDesignu(t, oboczne,
		`SELECT COUNT(*) FROM warstwa_kompozycji_design w
		   JOIN kompozycja_design k ON k.id = w.kompozycja_id
		  WHERE k.identyfikator_zewnetrzny = ?`, plansza); warstw != 0 {
		t.Fatalf("po wyczyszczeniu płótna w bazie zostało %d warstw", warstw)
	}

	var przywrocona shared.DesignBoardVersionRestoreResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardVersionRestore,
		shared.DesignBoardVersionRestoreRequest{VersionId: wersja.Version.Id}, &przywrocona)
	if len(przywrocona.Board.Layers) != 1 {
		t.Errorf("po przywróceniu kompozycja ma %d warstw, a ma mieć jedną",
			len(przywrocona.Board.Layers))
	}
	if przywrocona.SupersededVersion == nil {
		t.Fatal("przywrócenie nie odłożyło wersji z układu sprzed przywrócenia — " +
			"cofnięcie skasowałoby stan, który Operator właśnie porzucił")
	}
	if przywrocona.SupersededVersion.LayerCount != 0 {
		t.Errorf("wersja odłożona melduje %d warstw, a płótno było puste",
			przywrocona.SupersededVersion.LayerCount)
	}

	// Pomiar niezależny: warstwa wróciła do tabeli warstw żywych.
	if warstw := liczbaObocznaDesignu(t, oboczne,
		`SELECT COUNT(*) FROM warstwa_kompozycji_design w
		   JOIN kompozycja_design k ON k.id = w.kompozycja_id
		  WHERE k.identyfikator_zewnetrzny = ?`, plansza); warstw != 1 {
		t.Errorf("po przywróceniu w bazie leży %d warstw, a ma leżeć jedna", warstw)
	}

	var wykaz shared.DesignBoardVersionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardVersionList,
		shared.DesignBoardVersionListRequest{BoardId: plansza}, &wykaz)
	if wykaz.Total != 2 {
		t.Errorf("wykaz wersji ma %d pozycji, a powstały dwie", wykaz.Total)
	}
	for _, pozycja := range wykaz.Versions {
		if len(pozycja.Layers) != 0 {
			t.Errorf("wykaz wersji niesie warstwy wersji %s — miał ich nie nieść", pozycja.Id)
		}
	}
}

// ── Wyrys kompozycji ────────────────────────────────────────────────────────

// TestWyrysKompozycjiDesignuJestObrazemOZmierzonychWymiarach rozkłada wyrys
// z powrotem na obraz — wyrys, którego nie da się otworzyć, nie jest wyrysem.
func TestWyrysKompozycjiDesignuJestObrazemOZmierzonychWymiarach(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	zasob, _ := wniesObrazDesignu(t, zmontowany, zycie, "okno-wyrysu", 40, 20)
	plansza := zalozKompozycjeDesignu(t, zmontowany, zycie, "okno-wyrysu", zasob)

	var wyrys shared.DesignBoardExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardExport,
		shared.DesignBoardExportRequest{BoardId: plansza, Format: "png"}, &wyrys)
	if szerokosc, wysokosc := wymiaryWydaniaDesignu(t, wyrys.ContentBase64); szerokosc != 40 ||
		wysokosc != 20 {
		t.Errorf("wyrys ma %d×%d, a jedyna warstwa zajmuje 40×20", szerokosc, wysokosc)
	}

	var podwojny shared.DesignBoardExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardExport,
		shared.DesignBoardExportRequest{BoardId: plansza, Format: "png", Scale: wskaznik(2.0)},
		&podwojny)
	if szerokosc, wysokosc := wymiaryWydaniaDesignu(t, podwojny.ContentBase64); szerokosc != 80 ||
		wysokosc != 40 {
		t.Errorf("wyrys @2x ma %d×%d, a ma mieć 80×40", szerokosc, wysokosc)
	}

	var kadr shared.DesignBoardExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardExport,
		shared.DesignBoardExportRequest{BoardId: plansza, Format: "png",
			Region: &shared.DesignBoardRegion{X: 0, Y: 0, Width: 10, Height: 5}}, &kadr)
	if szerokosc, wysokosc := wymiaryWydaniaDesignu(t, kadr.ContentBase64); szerokosc != 10 ||
		wysokosc != 5 {
		t.Errorf("wyrys obszaru ma %d×%d, a obszar to 10×5", szerokosc, wysokosc)
	}

	var dokument shared.DesignBoardExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardExport,
		shared.DesignBoardExportRequest{BoardId: plansza, Format: "pdf"}, &dokument)
	bajtyDokumentu, err := base64.StdEncoding.DecodeString(dokument.ContentBase64)
	if err != nil {
		t.Fatalf("wyrys pdf nie jest poprawnym base64: %v", err)
	}
	if stron := stronDokumentuDesignu(t, bajtyDokumentu); stron != 1 {
		t.Errorf("wyrys pdf ma %d stron, a ma mieć jedną", stron)
	}

	var wektor shared.DesignBoardExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardExport,
		shared.DesignBoardExportRequest{BoardId: plansza, Format: "svg"}, &wektor)
	bajtyWektora, err := base64.StdEncoding.DecodeString(wektor.ContentBase64)
	if err != nil {
		t.Fatalf("wyrys svg nie jest poprawnym base64: %v", err)
	}
	tresc := string(bajtyWektora)
	if !strings.Contains(tresc, "<svg") || !strings.Contains(tresc, "</svg>") {
		t.Error("wyrys svg nie jest dokumentem svg")
	}
	// Treść osadzona, nie dowiązana: dokument wskazujący ścieżkę tej maszyny
	// byłby u odbiorcy pustą ramką.
	if !strings.Contains(tresc, "data:image/png;base64,") {
		t.Error("wyrys svg nie niesie osadzonej treści warstwy")
	}
	if strings.Contains(tresc, katalogMagazynuWSvgDesignu) {
		t.Error("wyrys svg wynosi ścieżkę magazynu rdzenia")
	}
}

// katalogMagazynuWSvgDesignu jest fragmentem ścieżki magazynu — wyrys nie ma
// prawa go nieść.
const katalogMagazynuWSvgDesignu = "design/zasoby"

// TestWyrysKompozycjiDesignuOdmawiaGdyNieMaCzegoWyrysowac pilnuje, że pusta
// tablica kończy się odmową, a nie przezroczystym prostokątem podanym jako plik.
func TestWyrysKompozycjiDesignuOdmawiaGdyNieMaCzegoWyrysowac(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var pusta shared.DesignBoardUpdateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardUpdate,
		shared.DesignBoardUpdateRequest{WindowId: "okno-pustej", Layers: []shared.DesignBoardLayer{}},
		&pusta)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignBoardExport,
		shared.DesignBoardExportRequest{BoardId: pusta.Board.Id, Format: "png"})
	if !strings.Contains(odmowa.Message, "ani jednej warstwy") {
		t.Errorf("odmowa wyrysu pustej tablicy nie nazywa braku: %q", odmowa.Message)
	}
}

// ── Adnotacje ───────────────────────────────────────────────────────────────

// TestAdnotacjaDesignuLezyWBazieIUkladaSieWWatek mierzy adnotacje obocznym
// połączeniem i sprawdza, że wątek naprawdę powstaje.
func TestAdnotacjaDesignuLezyWBazieIUkladaSieWWatek(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneDesignu(t, katalog)
	zasob, _ := wniesObrazDesignu(t, zmontowany, zycie, "okno-adnotacji", 8, 8)
	plansza := zalozKompozycjeDesignu(t, zmontowany, zycie, "okno-adnotacji", zasob)

	var pierwsza shared.DesignAnnotationSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAnnotationSet,
		shared.DesignAnnotationSetRequest{BoardId: plansza, Text: "margines za wąski"}, &pierwsza)

	if tresc := tekstObocznyDesignu(t, oboczne,
		`SELECT tresc FROM adnotacja_kompozycji_design WHERE identyfikator_zewnetrzny = ?`,
		pierwsza.Annotation.Id); tresc != "margines za wąski" {
		t.Errorf("w bazie stanowiska adnotacja niesie treść %q", tresc)
	}

	var odpowiedz shared.DesignAnnotationSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAnnotationSet,
		shared.DesignAnnotationSetRequest{
			BoardId:  plansza,
			ParentId: &pierwsza.Annotation.Id,
			Text:     "poprawione, proszę sprawdzić",
		}, &odpowiedz)
	if odpowiedz.Annotation.ParentId == nil || *odpowiedz.Annotation.ParentId != pierwsza.Annotation.Id {
		t.Error("odpowiedź w wątku nie wskazuje adnotacji nadrzędnej")
	}

	// Zamknięcie pierwszej: zawężenie do niezamkniętych ma ją pominąć,
	// a odpowiedź otwartą zostawić.
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAnnotationSet,
		shared.DesignAnnotationSetRequest{
			BoardId:      plansza,
			AnnotationId: &pierwsza.Annotation.Id,
			Text:         "margines za wąski",
			Resolved:     wskaznik(true),
		}, nil)

	if zamknietych := liczbaObocznaDesignu(t, oboczne,
		`SELECT COUNT(*) FROM adnotacja_kompozycji_design WHERE zamknieta = 1`); zamknietych != 1 {
		t.Errorf("w bazie stanowiska zamkniętych adnotacji jest %d, a ma być jedna", zamknietych)
	}

	var wszystkie shared.DesignAnnotationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAnnotationList,
		shared.DesignAnnotationListRequest{BoardId: plansza}, &wszystkie)
	if wszystkie.Total != 2 {
		t.Errorf("wykaz adnotacji ma %d pozycji, a powstały dwie", wszystkie.Total)
	}

	var otwarte shared.DesignAnnotationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAnnotationList,
		shared.DesignAnnotationListRequest{BoardId: plansza, OpenOnly: wskaznik(true)}, &otwarte)
	if otwarte.Total != 1 {
		t.Errorf("wykaz niezamkniętych ma %d pozycji, a ma mieć jedną", otwarte.Total)
	}

	// Adnotacja nadrzędna spoza kompozycji jest odmową — wątek rozpięty między
	// tablicami nie jest wątkiem.
	obca := "adnotacja-obca"
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignAnnotationSet,
		shared.DesignAnnotationSetRequest{BoardId: plansza, ParentId: &obca, Text: "wisząca"})
	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Errorf("odmowa za adnotację nadrzędną widmo ma kod %q", odmowa.Code)
	}
}

// ── Zestawy żetonów ─────────────────────────────────────────────────────────

// TestZestawZetonowDesignuLezyWBazieIWychodziDoKodu mierzy zestaw obocznym
// połączeniem i sprawdza, że wydanie do kodu naprawdę niesie role.
func TestZestawZetonowDesignuLezyWBazieIWychodziDoKodu(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneDesignu(t, katalog)

	var zapisany shared.DesignTokensetSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignTokensetSave,
		shared.DesignTokensetSaveRequest{
			WindowId: "okno-zetonow",
			Name:     "System marki",
			Theme:    wskaznik("ciemny"),
			Tokens: []shared.DesignToken{
				{Name: "sygnal", Kind: shared.DesignTokenKindColor, Value: "#c8a24a"},
				{Name: "od-4", Kind: shared.DesignTokenKindDimension, Value: "16px"},
				{Name: "czas-1", Kind: shared.DesignTokenKindDuration, Value: "120ms"},
			},
		}, &zapisany)

	if zapisany.TokenSet.TokenCount != 3 {
		t.Errorf("zestaw melduje %d żetonów, a zapisano trzy", zapisany.TokenSet.TokenCount)
	}
	// Pomiar niezależny: żetony leżą w pliku bazy stanowiska.
	if ile := liczbaObocznaDesignu(t, oboczne,
		`SELECT COUNT(*) FROM zeton_design t
		   JOIN zestaw_zetonow_design z ON z.id = t.zestaw_id
		  WHERE z.identyfikator_zewnetrzny = ?`, zapisany.TokenSet.Id); ile != 3 {
		t.Errorf("w bazie stanowiska leży %d żetonów, a odpowiedź meldowała trzy", ile)
	}
	if wartosc := tekstObocznyDesignu(t, oboczne,
		`SELECT t.wartosc FROM zeton_design t
		   JOIN zestaw_zetonow_design z ON z.id = t.zestaw_id
		  WHERE z.identyfikator_zewnetrzny = ? AND t.nazwa = ?`,
		zapisany.TokenSet.Id, "sygnal"); wartosc != "#c8a24a" {
		t.Errorf("w bazie stanowiska rola „sygnal” ma wartość %q", wartosc)
	}

	// Zapis jest pełny: zestaw nadpisany dwoma żetonami ma mieć dwa, nie pięć.
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignTokensetSave,
		shared.DesignTokensetSaveRequest{
			WindowId:   "okno-zetonow",
			Name:       "System marki",
			TokenSetId: &zapisany.TokenSet.Id,
			Tokens: []shared.DesignToken{
				{Name: "sygnal", Kind: shared.DesignTokenKindColor, Value: "#d4b25c"},
				{Name: "od-4", Kind: shared.DesignTokenKindDimension, Value: "16px"},
			},
		}, nil)
	if ile := liczbaObocznaDesignu(t, oboczne,
		`SELECT COUNT(*) FROM zeton_design t
		   JOIN zestaw_zetonow_design z ON z.id = t.zestaw_id
		  WHERE z.identyfikator_zewnetrzny = ?`, zapisany.TokenSet.Id); ile != 2 {
		t.Errorf("po nadpisaniu w bazie leży %d żetonów, a nadesłano dwa", ile)
	}

	var wykaz shared.DesignTokensetListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignTokensetList,
		shared.DesignTokensetListRequest{WindowId: "okno-zetonow"}, &wykaz)
	if wykaz.Total != 1 || len(wykaz.TokenSets[0].Tokens) != 2 {
		t.Fatalf("wykaz zestawów oddał %d zestawów o %d żetonach",
			wykaz.Total, len(wykaz.TokenSets[0].Tokens))
	}

	// Wydanie do kodu: każda postać ma nieść role, nie sam nagłówek.
	for _, postac := range shared.WartosciDesignTokenTarget() {
		var wydanie shared.DesignTokensetExportResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandDesignTokensetExport,
			shared.DesignTokensetExportRequest{TokenSetId: zapisany.TokenSet.Id, Target: postac},
			&wydanie)
		if !strings.Contains(wydanie.Content, "#d4b25c") {
			t.Errorf("wydanie %s nie niesie wartości roli: %q", postac, wydanie.Content)
		}
		if wydanie.FileName == "" {
			t.Errorf("wydanie %s nie proponuje nazwy pliku", postac)
		}
	}
}

// TestWczytanieZetonowDesignuNiesieBilansNieznanychRol sprawdza, że import
// wnosi role i wymienia te, których system produktu nie zna — bilans zamiast
// ciszy, a rola nieznana wchodzi do zestawu, zamiast zginąć.
func TestWczytanieZetonowDesignuNiesieBilansNieznanychRol(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneDesignu(t, katalog)

	zapis := `{"sygnal":"#c8a24a","marka-klienta":"#0055ff","od-4":"16px"}`
	var wczytany shared.DesignTokensetImportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignTokensetImport,
		shared.DesignTokensetImportRequest{
			WindowId:      "okno-importu",
			Name:          "System klienta",
			ContentBase64: wBase64([]byte(zapis)),
		}, &wczytany)

	if wczytany.TokenSet.TokenCount != 3 {
		t.Errorf("wczytany zestaw ma %d żetonów, a zapis niósł trzy role",
			wczytany.TokenSet.TokenCount)
	}
	if ile := liczbaObocznaDesignu(t, oboczne,
		`SELECT COUNT(*) FROM zeton_design t
		   JOIN zestaw_zetonow_design z ON z.id = t.zestaw_id
		  WHERE z.identyfikator_zewnetrzny = ?`, wczytany.TokenSet.Id); ile != 3 {
		t.Errorf("w bazie stanowiska leży %d wczytanych żetonów", ile)
	}
	if len(wczytany.UnknownNames) != 1 || wczytany.UnknownNames[0] != "marka-klienta" {
		t.Errorf("import nie zameldował roli nieznanej systemowi: %+v", wczytany.UnknownNames)
	}

	// Zapis w postaci zmiennych CSS wchodzi tą samą drogą.
	var zCss shared.DesignTokensetImportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignTokensetImport,
		shared.DesignTokensetImportRequest{
			WindowId:      "okno-importu",
			Name:          "System z arkusza",
			ContentBase64: wBase64([]byte(":root {\n  --dn-sygnal: #c8a24a;\n  --dn-od-4: 16px;\n}\n")),
		}, &zCss)
	if zCss.TokenSet.TokenCount != 2 {
		t.Errorf("wczytanie z arkusza CSS dało %d żetonów, a arkusz niósł dwa",
			zCss.TokenSet.TokenCount)
	}

	// Postać, której rdzeń nie czyta, jest odmową — nie zestawem z przypadkowych
	// napisów.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignTokensetImport,
		shared.DesignTokensetImportRequest{
			WindowId:      "okno-importu",
			Name:          "coś",
			ContentBase64: wBase64([]byte("to nie jest zapis żetonów")),
		})
	if odmowa.Message == "" {
		t.Error("odmowa postaci nierozpoznanej nie niesie zdania")
	}
}

// TestPrzewodnikStyluDesignuZakladaZasobZBajtami schodzi po odwołaniu zasobu do
// magazynu i czyta bajty przewodnika z dysku — przewodnik zapowiedziany bez
// bajtów byłby dokładnie tą szkodą, którą moduł ma w historii.
func TestPrzewodnikStyluDesignuZakladaZasobZBajtami(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	var zestaw shared.DesignTokensetSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignTokensetSave,
		shared.DesignTokensetSaveRequest{
			WindowId: "okno-przewodnika",
			Name:     "System marki",
			Tokens: []shared.DesignToken{
				{Name: "sygnal", Kind: shared.DesignTokenKindColor, Value: "#c8a24a"},
			},
		}, &zestaw)

	var kolekcja shared.DesignCollectionCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignCollectionCreate,
		shared.DesignCollectionCreateRequest{WindowId: "okno-przewodnika", Name: "Wydania"},
		&kolekcja)

	var wydany shared.DesignStyleguidePublishResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignStyleguidePublish,
		shared.DesignStyleguidePublishRequest{
			TokenSetId:     zestaw.TokenSet.Id,
			TargetModuleId: "library",
			CollectionId:   &kolekcja.Collection.Id,
		}, &wydany)

	if !wydany.Published || wydany.AssetId == "" {
		t.Fatalf("wydanie przewodnika melduje %v z zasobem %q", wydany.Published, wydany.AssetId)
	}

	// Pomiar niezależny numer jeden: bajty leżą pod odwołaniem w magazynie.
	var wykazZasobow shared.DesignAssetListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetList,
		shared.DesignAssetListRequest{WindowId: wskaznik("okno-przewodnika")}, &wykazZasobow)
	odwolanie := ""
	for _, zasob := range wykazZasobow.Assets {
		if zasob.Id == wydany.AssetId && zasob.Uri != nil {
			odwolanie = *zasob.Uri
		}
	}
	if odwolanie == "" {
		t.Fatal("zasób przewodnika nie ma odwołania do treści")
	}
	bajty := bajtyPodOdwolaniem(t, katalog, odwolanie)
	if !bytes.Contains(bajty, []byte("#c8a24a")) {
		t.Error("przewodnik leżący w magazynie nie niesie wartości roli")
	}
	if !bytes.Contains(bajty, []byte("<html")) {
		t.Error("przewodnik leżący w magazynie nie jest dokumentem")
	}

	// Pomiar niezależny numer dwa: treść wraca tą samą sumą, którą ma z bajtów.
	var tresc shared.DesignAssetContentGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetContentGet,
		shared.DesignAssetContentGetRequest{AssetId: wydany.AssetId}, &tresc)
	if tresc.Checksum != sumaSha256(bajty) {
		t.Errorf("suma treści przewodnika %q nie zgadza się z policzoną z bajtów %q",
			tresc.Checksum, sumaSha256(bajty))
	}

	// Przewodnik trafił do wskazanej kolekcji.
	var wykaz shared.DesignCollectionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignCollectionList,
		shared.DesignCollectionListRequest{WindowId: "okno-przewodnika", AssetId: &wydany.AssetId},
		&wykaz)
	if wykaz.Total != 1 {
		t.Errorf("przewodnika nie ma w kolekcji, do której go wydano (%d kolekcji)", wykaz.Total)
	}
}

// ── Obecność ────────────────────────────────────────────────────────────────

// TestObecnoscDesignuJestUlotnaINieZostawiaWiersza pilnuje obietnicy kontraktu:
// zgłoszenie obecności nie zapisuje się w bazie.
func TestObecnoscDesignuJestUlotnaINieZostawiaWiersza(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	oboczne := polaczenieOboczneDesignu(t, katalog)
	zasob, _ := wniesObrazDesignu(t, zmontowany, zycie, "okno-obecnosci", 8, 8)
	plansza := zalozKompozycjeDesignu(t, zmontowany, zycie, "okno-obecnosci", zasob)

	var obecni shared.DesignPresenceReportResponse
	wykonajUdana(t, zmontowany, zycieZGniazdemDesignu(zycie, "klient-sprawdzianu"),
		shared.CommandDesignPresenceReport,
		shared.DesignPresenceReportRequest{
			BoardId: plansza, X: wskaznik(12.5), Y: wskaznik(30.0),
		}, &obecni)
	if len(obecni.Participants) != 1 || obecni.Participants[0].ClientId != "klient-sprawdzianu" {
		t.Errorf("zgłoszenie oddało obecnych %+v — miało oddać jednego, tego zgłaszającego",
			obecni.Participants)
	}

	// Pomiar niezależny: w schemacie nie ma tabeli obecności i nie ma jej mieć.
	var tabel int
	if err := oboczne.QueryRow(
		`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name LIKE '%obecnos%'`,
	).Scan(&tabel); err != nil {
		t.Fatalf("oboczny odczyt schematu nie powiódł się: %v", err)
	}
	if tabel != 0 {
		t.Errorf("obecność ma tabelę w schemacie (%d) — a jest bytem ulotnym", tabel)
	}

	odmowa := wykonajOdmowna(t, zmontowany, zycieZGniazdemDesignu(zycie, "klient-sprawdzianu"),
		shared.CommandDesignPresenceReport,
		shared.DesignPresenceReportRequest{BoardId: "plansza-widmo"})
	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Errorf("obecność na kompozycji widmie nie została odmówiona: %q", odmowa.Code)
	}

	// Zgłoszenie bez rozpoznanego klienta jest odmawiane — kursor bez
	// tożsamości nie da się ani odróżnić od cudzego, ani zdjąć przy odejściu.
	bezKlienta := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignPresenceReport,
		shared.DesignPresenceReportRequest{BoardId: plansza})
	if bezKlienta.Code != shared.ErrorCodeValidationFailed {
		t.Errorf("zgłoszenie bez klienta ma kod %q, a ma być odmową wskazania", bezKlienta.Code)
	}
}

// TestRejestrObecnosciDesignuZdejmujeZgloszeniaPrzeterminowane mierzy sam
// rejestr: kursor po kimś, kto zamknął kartę bez zgłoszenia odejścia, ma zniknąć.
func TestRejestrObecnosciDesignuZdejmujeZgloszeniaPrzeterminowane(t *testing.T) {
	rejestr := nowyRejestrObecnosciDesignu()
	teraz := int64(1_000_000)
	rejestr.terazTestem = func() time.Time { return time.UnixMilli(teraz) }

	rejestr.zglos("plansza-1", shared.DesignPresence{ClientId: "duch",
		SeenAt: teraz - terminObecnosciDesignu.Milliseconds() - 1}, false)
	obecni := rejestr.zglos("plansza-1", shared.DesignPresence{ClientId: "zywy", SeenAt: teraz}, false)

	if len(obecni) != 1 || obecni[0].ClientId != "zywy" {
		t.Errorf("rejestr oddał obecnych %+v — kursor po kliencie sprzed terminu miał zniknąć", obecni)
	}

	po := rejestr.zglos("plansza-1", shared.DesignPresence{ClientId: "zywy", SeenAt: teraz}, true)
	if len(po) != 0 {
		t.Errorf("po zdjęciu obecności rejestr oddał %d obecnych", len(po))
	}
}
