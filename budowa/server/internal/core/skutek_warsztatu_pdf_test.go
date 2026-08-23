package core

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"danacoconsole/shared"
)

// Skutek warsztatu dokumentu: czy za wynikiem leży prawdziwy PDF.
//
// Wzorzec szkody, którego pilnuje ten plik, ma w tym produkcie precedens
// w module Design: komenda meldowała `status: ok` z wykazem zasobów, za którymi
// nie było ani jednego bajtu. Dlatego żaden sprawdzian tutaj nie kończy się na
// sprawdzeniu, że odpowiedź jest udana: każdy schodzi po odwołaniu do magazynu,
// otwiera plik i liczy jego strony z bajtów, a nie z odpowiedzi komendy.
//
// Sprawdzian nie pomija się przy braku żadnego programu, bo warsztat PDF nie
// uruchamia ani jednego procesu: pracuje biblioteką wkompilowaną w rdzeń.
// Materiał powstaje tą samą biblioteką i to jest świadome — mierzone jest
// działanie warsztatu, nie zgodność dwóch bibliotek między sobą.

// pdfProbny wytwarza dokument o wskazanej liczbie stron i oddaje jego bajty.
func pdfProbny(t *testing.T, strony int) []byte {
	t.Helper()

	var opis strings.Builder
	opis.WriteString(`{"pages":{`)
	for numer := 1; numer <= strony; numer++ {
		if numer > 1 {
			opis.WriteString(",")
		}
		fmt.Fprintf(&opis, `"%d":{"content":{"text":[{"value":"Strona %d",`+
			`"font":{"name":"Helvetica","size":24},"position":[0.5,0.5]}]}}`, numer, numer)
	}
	opis.WriteString(`}}`)

	var dokument bytes.Buffer
	if err := api.Create(nil, strings.NewReader(opis.String()), &dokument, nastawyPdf()); err != nil {
		t.Fatalf("nie można wytworzyć materiału próbnego: %v", err)
	}
	return dokument.Bytes()
}

// stronWyniku liczy strony dokumentu leżącego pod wskazaną ścieżką magazynu.
func stronWyniku(t *testing.T, sciezka string) int {
	t.Helper()

	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać wyniku: %v", err)
	}
	liczba, err := api.PageCount(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		t.Fatalf("nie można policzyć stron wyniku: %v", err)
	}
	return liczba
}

// wniesPdf wnosi dokument do magazynu zasobów i oddaje jego identyfikator.
func wniesPdf(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	bajty []byte, nazwa string) string {
	t.Helper()

	var wynik shared.DesignAssetUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetUpload,
		shared.DesignAssetUploadRequest{
			WindowId:      "okno-sprawdzianu",
			Name:          wskaznik(nazwa),
			Kind:          shared.DesignAssetKindDocument,
			Format:        wskaznik("pdf"),
			ContentBase64: wskaznik(wBase64(bajty)),
		}, &wynik)
	return wynik.Asset.Id
}

// sciezkaZasobuSprawdzianu schodzi po odwołaniu zasobu do pliku w magazynie.
func sciezkaZasobuSprawdzianu(t *testing.T, zmontowany *Zmontowany,
	zycie context.Context, katalog, kod string) string {
	t.Helper()

	var wykaz shared.DesignAssetListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetList,
		shared.DesignAssetListRequest{WindowId: wskaznik("okno-sprawdzianu")}, &wykaz)
	for _, zasob := range wykaz.Assets {
		if zasob.Id != kod {
			continue
		}
		if zasob.Uri == nil || *zasob.Uri == "" {
			t.Fatalf("zasób %s nie niesie odwołania do bajtów", kod)
		}
		if filepath.IsAbs(*zasob.Uri) {
			return *zasob.Uri
		}
		return filepath.Join(katalog, *zasob.Uri)
	}
	t.Fatalf("zasobu %s nie ma w wykazie okna", kod)
	return ""
}

// TestScalaniePdfDajePlikOSumieStron wykazuje skutek: za wynikiem leży dokument,
// który ma tyle stron, ile miały materiały razem.
func TestScalaniePdfDajePlikOSumieStron(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	kodPierwszego := wniesPdf(t, zmontowany, zycie, pdfProbny(t, 2), "pierwszy")
	kodDrugiego := wniesPdf(t, zmontowany, zycie, pdfProbny(t, 3), "drugi")

	var wynik shared.StudioPdfMergeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPdfMerge,
		shared.StudioPdfMergeRequest{
			AssetIds: []string{kodPierwszego, kodDrugiego},
			WindowId: wskaznik("okno-sprawdzianu"),
		}, &wynik)

	if wynik.Pages != 5 {
		t.Fatalf("scalony dokument melduje %d stron, a materiały miały 2 + 3", wynik.Pages)
	}
	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	if strony := stronWyniku(t, sciezka); strony != 5 {
		t.Fatalf("za wynikiem leży dokument o %d stronach, a odpowiedź meldowała 5", strony)
	}
}

// TestPodzialPdfDajeTyleCzesciIleZakresow wykazuje, że każda część jest osobnym
// dokumentem o właściwej liczbie stron, a nie kopią całości.
func TestPodzialPdfDajeTyleCzesciIleZakresow(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	kod := wniesPdf(t, zmontowany, zycie, pdfProbny(t, 4), "caly")

	var wynik shared.StudioPdfSplitResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPdfSplit,
		shared.StudioPdfSplitRequest{
			AssetId:  kod,
			Ranges:   []string{"1-2", "3-4"},
			WindowId: wskaznik("okno-sprawdzianu"),
		}, &wynik)

	if wynik.Parts != 2 || len(wynik.AssetIds) != 2 {
		t.Fatalf("podział oddał %d części przy dwóch zakresach", wynik.Parts)
	}
	for _, kodCzesci := range wynik.AssetIds {
		sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, kodCzesci)
		if strony := stronWyniku(t, sciezka); strony != 2 {
			t.Fatalf("część %s ma %d stron, a zakres wskazywał dwie", kodCzesci, strony)
		}
	}
}

// TestPodzialBezZakresowDzieliNaPojedynczeStrony pilnuje znaczenia przyjętego
// dla żądania bez zakresów: każda strona osobno, a nie cały dokument jako
// jedna część.
func TestPodzialBezZakresowDzieliNaPojedynczeStrony(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	kod := wniesPdf(t, zmontowany, zycie, pdfProbny(t, 3), "trzy")

	var wynik shared.StudioPdfSplitResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPdfSplit,
		shared.StudioPdfSplitRequest{AssetId: kod, WindowId: wskaznik("okno-sprawdzianu")},
		&wynik)

	if wynik.Parts != 3 {
		t.Fatalf("podział bez zakresów oddał %d części przy dokumencie o trzech stronach",
			wynik.Parts)
	}
	for _, kodCzesci := range wynik.AssetIds {
		sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, kodCzesci)
		if strony := stronWyniku(t, sciezka); strony != 1 {
			t.Fatalf("część %s ma %d stron, a miała nieść jedną", kodCzesci, strony)
		}
	}
}

// TestPrzestawienieStronZostawiaDokumentOZadanejDlugosci wykazuje skutek
// wyodrębnienia: za wynikiem leży dokument o liczbie stron wskazanej działaniem,
// a nie kopia materiału.
func TestPrzestawienieStronZostawiaDokumentOZadanejDlugosci(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	kod := wniesPdf(t, zmontowany, zycie, pdfProbny(t, 5), "piec")

	var wynik shared.StudioPdfPagesReorderResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioPdfPagesReorder,
		shared.StudioPdfPagesReorderRequest{
			AssetId:  kod,
			WindowId: wskaznik("okno-sprawdzianu"),
			Operations: []shared.StudioPdfPageOperation{{
				Kind:  shared.StudioPdfPageOperationKindWyodrebnienie,
				Pages: "2-3",
			}},
		}, &wynik)

	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	if strony := stronWyniku(t, sciezka); strony != 2 {
		t.Fatalf("po wyodrębnieniu dwóch stron wynik ma %d stron", strony)
	}
}

// TestWarsztatOdmawiaMaterialowiNiebedacemuDokumentem pilnuje odmowy, która
// mówi Operatorowi, co podał: bez niej biblioteka odmówiłaby komunikatem
// o strukturze pliku, z którego nie wynika, że materiałem był obraz.
func TestWarsztatOdmawiaMaterialowiNiebedacemuDokumentem(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var wgrany shared.DesignAssetUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetUpload,
		shared.DesignAssetUploadRequest{
			WindowId:      "okno-sprawdzianu",
			Name:          wskaznik("nie-dokument"),
			Kind:          shared.DesignAssetKindImage,
			Format:        wskaznik("txt"),
			ContentBase64: wskaznik(wBase64([]byte("to nie jest dokument"))),
		}, &wgrany)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandStudioPdfSplit,
		shared.StudioPdfSplitRequest{AssetId: wgrany.Asset.Id})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("odmowa ma nieść kod validation_failed, a niesie %s", odmowa.Code)
	}
	if !strings.Contains(odmowa.Message, "PDF") {
		t.Fatalf("odmowa nie mówi, że materiał nie jest dokumentem PDF: %q", odmowa.Message)
	}
}
