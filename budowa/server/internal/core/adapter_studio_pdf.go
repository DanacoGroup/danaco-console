// Warsztat dokumentu PDF modułu Studio: scalanie, podział, stemplowanie
// i odchudzanie na bibliotece pdfcpu wkompilowanej w rdzeń, nigdy na programie
// zewnętrznym; materiał wchodzi zasobem, wynik wychodzi nowym zasobem.
package core

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/form"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

type adapterPdfStudia struct {
	zasoby  dane.RepozytoriumDesignu
	magazyn *magazynTresciBiblioteki
}

func nowyAdapterPdfStudia() *adapterPdfStudia {
	return &adapterPdfStudia{magazyn: magazynZasobowDesignu(konfiguracja.KatalogDanychDomyslny())}
}

// ZKatalogiemDanych przestawia magazyn wyników na katalog wskazany konfiguracją,
// zamiast na katalog domyślny założony przy budowie adaptera.
func (a *adapterPdfStudia) ZKatalogiemDanych(katalog string) *adapterPdfStudia {
	if strings.TrimSpace(katalog) != "" {
		a.magazyn = magazynZasobowDesignu(katalog)
	}
	return a
}

// ZZasobami podaje magazyn zasobów: repozytorium, przez które adapter czyta
// materiał wejściowy i zapisuje wiersz zasobu wynikowego.
func (a *adapterPdfStudia) ZZasobami(r dane.RepozytoriumDesignu) *adapterPdfStudia {
	a.zasoby = r
	return a
}

// odmowaPdf buduje odmowę rodziny warsztatu, poprzedzając powód stałym
// przedrostkiem, po którym rozpoznaje się źródło komunikatu.
func odmowaPdf(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "warsztat PDF: "+powod))
}

// nastawyPdf oddaje konfigurację biblioteki ze sprawdzaniem zgodności
// wyłączonym: dokumenty zastane bywają niezgodne ze specyfikacją w szczegółach,
// których dokument otwierany przeglądarką i tak nie ujawnia.
func nastawyPdf() *model.Configuration {
	nastawy := model.NewDefaultConfiguration()
	nastawy.ValidationMode = model.ValidationRelaxed
	return nastawy
}

// trescMaterialu oddaje bajty wskazanego zasobu wraz ze sprawdzeniem, że jest
// dokumentem PDF — biblioteka odmówiłaby dopiero komunikatem o strukturze,
// nie o rodzaju pliku.
func (a *adapterPdfStudia) trescMaterialu(ctx context.Context, kod string) ([]byte, error) {
	bajty, err := a.bajtyZasobu(ctx, kod)
	if err != nil {
		return nil, err
	}
	if !bytes.HasPrefix(bajty, []byte("%PDF-")) {
		return nil, odmowaPdf(shared.ErrorCodeValidationFailed,
			"zasób "+kod+" nie jest dokumentem PDF; warsztat pracuje wyłącznie na PDF")
	}
	return bajty, nil
}

// bajtyZasobu oddaje treść zasobu bez rozstrzygania o jej formacie: pieczęć
// graficzna dokumentem nie jest, a wchodzi tą samą drogą.
func (a *adapterPdfStudia) bajtyZasobu(ctx context.Context, kod string) ([]byte, error) {
	if a.zasoby == nil {
		return nil, odmowaPdf(shared.ErrorCodeInternalError,
			"serwer nie ma wpiętego magazynu zasobów — naprawa: podpiąć repozytorium "+
				"zasobów przy składaniu serwera")
	}
	wiersz, err := a.zasoby.Zasob(ctx, strings.TrimSpace(kod))
	if err != nil {
		return nil, odmowaPdf(shared.ErrorCodeNotFound,
			"zasobu "+kod+" nie ma w magazynie serwera")
	}
	if wiersz.URI == nil || strings.TrimSpace(*wiersz.URI) == "" {
		return nil, odmowaPdf(shared.ErrorCodeNotFound,
			"zasób "+kod+" nie ma odwołania do bajtów — nie ma czego przetworzyć")
	}
	bajty, err := os.ReadFile(*wiersz.URI)
	if err != nil {
		return nil, odmowaPdf(shared.ErrorCodeInternalError,
			"nie można odczytać treści zasobu "+kod+": "+err.Error())
	}
	return bajty, nil
}

// odlozPdf utrwala wynik w magazynie pod formatem PDF i zakłada wiersz zasobu,
// wywołując odlozTresc ze stałym formatem dokumentu.
func (a *adapterPdfStudia) odlozPdf(ctx context.Context, bajty []byte,
	nazwa, okno string) (shared.DesignAsset, error) {
	return a.odlozTresc(ctx, bajty, nazwa, "pdf", okno)
}

// odlozTresc utrwala dowolną treść wyniku pod wskazanym formatem.
//
// Wyciąganie obrazów i załączników oddaje pozycje, które dokumentami nie są,
// więc format nie może być tu na stałe wpisany.
func (a *adapterPdfStudia) odlozTresc(ctx context.Context, bajty []byte,
	nazwa, format, okno string) (shared.DesignAsset, error) {

	// Suma kontrolna liczy się tutaj: magazyn przyjmuje ją gotową i rozstrzyga
	// nazwę bloba.
	skrot := sha256.Sum256(bajty)
	suma := hex.EncodeToString(skrot[:])
	odwolanie, err := a.magazyn.Zapisz(bajty, suma)
	if err != nil {
		return shared.DesignAsset{}, odmowaPdf(shared.ErrorCodeInternalError, err.Error())
	}
	return odlozWynikArsenalu(ctx, a.zasoby, wynikArsenalu{
		odwolanie: odwolanie, okno: okno, nazwa: nazwa, format: format,
	}, func(powod string) error {
		return odmowaPdf(shared.ErrorCodeInternalError, powod)
	})
}

// stronDokumentu liczy strony biblioteką pdfcpu — tą samą, która dokument
// przetwarza dalej, więc liczba stron zgadza się z wynikiem czynności.
func stronDokumentu(bajty []byte) (int, error) {
	liczba, err := api.PageCount(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		return 0, odmowaPdf(shared.ErrorCodeValidationFailed,
			"nie można odczytać liczby stron dokumentu: "+err.Error())
	}
	return liczba, nil
}

// ── Czynności ───────────────────────────────────────────────────────────────

// ScalPdf łączy wskazane dokumenty w jeden, w podanej kolejności, i oddaje
// zasób wyniku wraz z liczbą stron dokumentu scalonego.
func (a *adapterPdfStudia) ScalPdf(ctx context.Context,
	z shared.StudioPdfMergeRequest) (shared.StudioPdfMergeResponse, error) {

	if len(z.AssetIds) < 2 {
		return shared.StudioPdfMergeResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"scalanie wymaga co najmniej dwóch dokumentów; podano "+
				strconv.Itoa(len(z.AssetIds)))
	}
	materialy := make([]io.ReadSeeker, 0, len(z.AssetIds))
	for _, kod := range z.AssetIds {
		bajty, err := a.trescMaterialu(ctx, kod)
		if err != nil {
			return shared.StudioPdfMergeResponse{}, err
		}
		materialy = append(materialy, bytes.NewReader(bajty))
	}

	var wynik bytes.Buffer
	if err := api.MergeRaw(materialy, &wynik, false, nastawyPdf()); err != nil {
		return shared.StudioPdfMergeResponse{}, odmowaPdf(shared.ErrorCodeInternalError,
			"scalanie nie powiodło się: "+err.Error())
	}

	strony, err := stronDokumentu(wynik.Bytes())
	if err != nil {
		return shared.StudioPdfMergeResponse{}, err
	}
	zasob, err := a.odlozPdf(ctx, wynik.Bytes(),
		pierwszyNiepusty(wartoscTekstu(z.Name), "dokument scalony"), wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioPdfMergeResponse{}, err
	}
	return shared.StudioPdfMergeResponse{Asset: zasob, Pages: strony}, nil
}

// PodzielPdf rozdziela dokument na części wskazane zakresami stron.
//
// Brak zakresów znaczy „każda strona osobno" — to przyjęte znaczenie
// w narzędziach tej klasy.
func (a *adapterPdfStudia) PodzielPdf(ctx context.Context,
	z shared.StudioPdfSplitRequest) (shared.StudioPdfSplitResponse, error) {

	bajty, err := a.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioPdfSplitResponse{}, err
	}
	zakresy := z.Ranges
	if len(zakresy) == 0 {
		strony, err := stronDokumentu(bajty)
		if err != nil {
			return shared.StudioPdfSplitResponse{}, err
		}
		for numer := 1; numer <= strony; numer++ {
			zakresy = append(zakresy, strconv.Itoa(numer))
		}
	}

	kody := make([]string, 0, len(zakresy))
	for _, zakres := range zakresy {
		var czesc bytes.Buffer
		if err := api.Trim(bytes.NewReader(bajty), &czesc,
			[]string{zakres}, nastawyPdf()); err != nil {
			return shared.StudioPdfSplitResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
				"zakres stron "+zakres+" nie da się wydzielić: "+err.Error())
		}
		zasob, err := a.odlozPdf(ctx, czesc.Bytes(), "część "+zakres, wartoscTekstu(z.WindowId))
		if err != nil {
			return shared.StudioPdfSplitResponse{}, err
		}
		kody = append(kody, zasob.Id)
	}
	return shared.StudioPdfSplitResponse{AssetIds: kody, Parts: len(kody)}, nil
}

// OdchudzPdf zmniejsza objętość dokumentu porządkowaniem struktury i usuwaniem
// powielonych zasobów, bez przeliczania obrazów w dół — odpowiedź niesie zysk
// zmierzony, nie obietnicę.
func (a *adapterPdfStudia) OdchudzPdf(ctx context.Context,
	z shared.StudioPdfOptimizeRequest) (shared.StudioPdfOptimizeResponse, error) {

	bajty, err := a.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioPdfOptimizeResponse{}, err
	}
	var wynik bytes.Buffer
	if err := api.Optimize(bytes.NewReader(bajty), &wynik, nastawyPdf()); err != nil {
		return shared.StudioPdfOptimizeResponse{}, odmowaPdf(shared.ErrorCodeInternalError,
			"odchudzanie nie powiodło się: "+err.Error())
	}

	zasob, err := a.odlozPdf(ctx, wynik.Bytes(), "dokument odchudzony", wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioPdfOptimizeResponse{}, err
	}
	// Zysk bywa ujemny — dokument odchudzony rośnie po ponownym zapisie; liczba
	// mówi prawdę zamiast zera.
	return shared.StudioPdfOptimizeResponse{
		Asset:      zasob,
		SizeBytes:  wynik.Len(),
		SavedBytes: len(bajty) - wynik.Len(),
	}, nil
}

// UlozStrony wykonuje czynności na stronach dokumentu sekwencyjnie na tych
// samych bajtach, bo wynik pierwszej jest materiałem drugiej — inaczej
// numeracja przesunięta usunięciem trafiałaby w stronę inną, niż wskazano.
func (a *adapterPdfStudia) UlozStrony(ctx context.Context,
	z shared.StudioPdfPagesReorderRequest) (shared.StudioPdfPagesReorderResponse, error) {

	bajty, err := a.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioPdfPagesReorderResponse{}, err
	}
	if len(z.Operations) == 0 {
		return shared.StudioPdfPagesReorderResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"żądanie nie niesie ani jednej czynności — nie ma czego ułożyć")
	}

	for _, czynnosc := range z.Operations {
		strony := []string{czynnosc.Pages}
		var wynik bytes.Buffer
		switch czynnosc.Kind {
		case shared.StudioPdfPageOperationKindUsuniecie:
			err = api.RemovePages(bytes.NewReader(bajty), &wynik, strony, nastawyPdf())

		case shared.StudioPdfPageOperationKindWyodrebnienie:
			err = api.Trim(bytes.NewReader(bajty), &wynik, strony, nastawyPdf())

		case shared.StudioPdfPageOperationKindPrzeniesienie:
			// Zebranie stron w kolejności wskazanej — inaczej niż wydzielenie,
			// zachowujące kolejność dokumentu.
			err = api.Collect(bytes.NewReader(bajty), &wynik, strony, nastawyPdf())

		case shared.StudioPdfPageOperationKindObrot:
			stopnie := 90
			if czynnosc.Degrees != nil {
				stopnie = *czynnosc.Degrees
			}
			if stopnie%90 != 0 {
				return shared.StudioPdfPagesReorderResponse{}, odmowaPdf(
					shared.ErrorCodeValidationFailed,
					"obrót o "+strconv.Itoa(stopnie)+" stopni nie jest wielokrotnością 90; "+
						"dokument obraca się ćwiartkami")
			}
			err = api.Rotate(bytes.NewReader(bajty), &wynik, stopnie, strony, nastawyPdf())

		case shared.StudioPdfPageOperationKindWstawienie:
			return shared.StudioPdfPagesReorderResponse{}, odmowaPdf(
				shared.ErrorCodeValidationFailed,
				"wstawienie stron z innego dokumentu jest scalaniem z wyborem miejsca — "+
					"użyj studio.pdf.merge; ta czynność wstawiłaby strony puste")

		default:
			return shared.StudioPdfPagesReorderResponse{}, odmowaPdf(
				shared.ErrorCodeValidationFailed,
				"czynność "+string(czynnosc.Kind)+" nie ma odwzorowania w warsztacie")
		}
		if err != nil {
			return shared.StudioPdfPagesReorderResponse{}, odmowaPdf(
				shared.ErrorCodeValidationFailed,
				"czynność "+string(czynnosc.Kind)+" na stronach "+czynnosc.Pages+
					" nie powiodła się: "+err.Error())
		}
		bajty = wynik.Bytes()
	}

	liczba, err := stronDokumentu(bajty)
	if err != nil {
		return shared.StudioPdfPagesReorderResponse{}, err
	}
	zasob, err := a.odlozPdf(ctx, bajty, "dokument po zmianie stron", wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioPdfPagesReorderResponse{}, err
	}
	return shared.StudioPdfPagesReorderResponse{Asset: zasob, Pages: liczba}, nil
}

// opisPieczeci składa opis wyglądu pieczęci w postaci tekstowej, którą czyta
// biblioteka pdfcpu: położenie, skalę, krycie i obrót.
func opisPieczeci(krycie, obrot *int) string {
	czesci := []string{"scalefactor:0.7 rel", "position:c"}
	if krycie != nil {
		if *krycie < 0 || *krycie > 100 {
			return ""
		}
		czesci = append(czesci, "opacity:"+strconv.FormatFloat(float64(*krycie)/100, 'f', 2, 64))
	}
	if obrot != nil {
		czesci = append(czesci, "rotation:"+strconv.Itoa(*obrot))
	}
	return strings.Join(czesci, ", ")
}

// OstemplujPdf nakłada na strony pieczęć tekstową albo graficzną.
//
// Pieczęć idzie NAD treść, nie pod nią: pieczęć schowana pod treścią strony
// byłaby pieczęcią, której Operator nie widzi, a meldunek o powodzeniu i tak by
// przyszedł.
func (a *adapterPdfStudia) OstemplujPdf(ctx context.Context,
	z shared.StudioPdfStampRequest) (shared.StudioPdfStampResponse, error) {

	tresc := strings.TrimSpace(wartoscTekstu(z.Text))
	pieczecGraficzna := strings.TrimSpace(wartoscTekstu(z.StampAssetId))
	if tresc == "" && pieczecGraficzna == "" {
		return shared.StudioPdfStampResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"żądanie nie niesie ani treści pieczęci, ani zasobu pieczęci graficznej")
	}
	if tresc != "" && pieczecGraficzna != "" {
		return shared.StudioPdfStampResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"podano naraz treść pieczęci i pieczęć graficzną — wskaż jedną")
	}

	bajty, err := a.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioPdfStampResponse{}, err
	}
	opis := opisPieczeci(z.Opacity, z.Rotation)
	if opis == "" {
		return shared.StudioPdfStampResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"krycie pieczęci podaje się w procentach od 0 do 100")
	}

	var pieczec *model.Watermark
	if tresc != "" {
		pieczec, err = api.TextWatermark(tresc, opis, true, false, types.POINTS)
	} else {
		var obraz []byte
		obraz, err = a.bajtyZasobu(ctx, pieczecGraficzna)
		if err != nil {
			return shared.StudioPdfStampResponse{}, err
		}
		pieczec, err = api.ImageWatermarkForReader(bytes.NewReader(obraz), opis,
			true, false, types.POINTS)
	}
	if err != nil {
		return shared.StudioPdfStampResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"nie można złożyć pieczęci: "+err.Error())
	}

	var wynik bytes.Buffer
	if err := api.AddWatermarks(bytes.NewReader(bajty), &wynik,
		zakresStron(z.Pages), pieczec, nastawyPdf()); err != nil {
		return shared.StudioPdfStampResponse{}, odmowaPdf(shared.ErrorCodeInternalError,
			"stemplowanie nie powiodło się: "+err.Error())
	}
	zasob, err := a.odlozPdf(ctx, wynik.Bytes(), "dokument ostemplowany",
		wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioPdfStampResponse{}, err
	}
	return shared.StudioPdfStampResponse{Asset: zasob}, nil
}

// zakresStron zamienia wskazanie zakresu na postać przyjmowaną przez bibliotekę.
// Brak wskazania znaczy „wszystkie strony", co biblioteka rozumie jako brak
// wyboru.
func zakresStron(zakres *string) []string {
	if zakres == nil || strings.TrimSpace(*zakres) == "" {
		return nil
	}
	return []string{strings.TrimSpace(*zakres)}
}

// NumerujPdf nakłada numerację prawną Bates wraz z nagłówkiem i stopką, stronie
// po stronie, bo biblioteka nie zna wzorca „numer bieżącej strony".
func (a *adapterPdfStudia) NumerujPdf(ctx context.Context,
	z shared.StudioPdfBatesRequest) (shared.StudioPdfBatesResponse, error) {

	bajty, err := a.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioPdfBatesResponse{}, err
	}
	strony, err := stronDokumentu(bajty)
	if err != nil {
		return shared.StudioPdfBatesResponse{}, err
	}

	poczatek := 1
	if z.StartNumber != nil {
		poczatek = *z.StartNumber
	}
	cyfry := 0
	if z.Digits != nil {
		if *z.Digits < 1 || *z.Digits > 12 {
			return shared.StudioPdfBatesResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
				"liczba cyfr numeru mieści się między 1 a 12; podano "+strconv.Itoa(*z.Digits))
		}
		cyfry = *z.Digits
	}
	przedrostek := wartoscTekstu(z.Prefix)
	naglowek := strings.TrimSpace(wartoscTekstu(z.Header))
	stopka := strings.TrimSpace(wartoscTekstu(z.Footer))

	ostatni := poczatek + strony - 1
	for numer := 0; numer < strony; numer++ {
		liczba := strconv.Itoa(poczatek + numer)
		for len(liczba) < cyfry {
			liczba = "0" + liczba
		}
		strona := []string{strconv.Itoa(numer + 1)}

		nadruki := []struct{ tresc, polozenie string }{
			{przedrostek + liczba, "br"},
		}
		if naglowek != "" {
			nadruki = append(nadruki, struct{ tresc, polozenie string }{naglowek, "tc"})
		}
		if stopka != "" {
			nadruki = append(nadruki, struct{ tresc, polozenie string }{stopka, "bc"})
		}

		for _, nadruk := range nadruki {
			pieczec, err := api.TextWatermark(nadruk.tresc,
				"points:10, position:"+nadruk.polozenie+", scalefactor:1 abs, rotation:0",
				true, false, types.POINTS)
			if err != nil {
				return shared.StudioPdfBatesResponse{}, odmowaPdf(shared.ErrorCodeInternalError,
					"nie można złożyć nadruku numeracji: "+err.Error())
			}
			var wynik bytes.Buffer
			if err := api.AddWatermarks(bytes.NewReader(bajty), &wynik,
				strona, pieczec, nastawyPdf()); err != nil {
				return shared.StudioPdfBatesResponse{}, odmowaPdf(shared.ErrorCodeInternalError,
					"numeracja strony "+strona[0]+" nie powiodła się: "+err.Error())
			}
			bajty = wynik.Bytes()
		}
	}

	zasob, err := a.odlozPdf(ctx, bajty, "dokument ponumerowany", wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioPdfBatesResponse{}, err
	}
	return shared.StudioPdfBatesResponse{Asset: zasob, LastNumber: ostatni}, nil
}

// zakladkiBiblioteki przenosi drzewo zakładek kontraktu na postać biblioteki
// pdfcpu, sprawdzając po drodze tytuł i numer strony każdej zakładki.
func zakladkiBiblioteki(zrodlo []shared.StudioPdfBookmark) ([]pdfcpu.Bookmark, error) {
	wynik := make([]pdfcpu.Bookmark, 0, len(zrodlo))
	for _, zakladka := range zrodlo {
		if strings.TrimSpace(zakladka.Title) == "" {
			return nil, odmowaPdf(shared.ErrorCodeValidationFailed,
				"zakładka bez tytułu nie ma czego pokazać w drzewie nawigacji")
		}
		if zakladka.Page < 1 {
			return nil, odmowaPdf(shared.ErrorCodeValidationFailed,
				"zakładka „"+zakladka.Title+"” wskazuje stronę "+strconv.Itoa(zakladka.Page)+
					"; strony liczą się od jedynki")
		}
		dzieci, err := zakladkiBiblioteki(zakladka.Children)
		if err != nil {
			return nil, err
		}
		wynik = append(wynik, pdfcpu.Bookmark{
			Title: zakladka.Title, PageFrom: zakladka.Page, Kids: dzieci,
		})
	}
	return wynik, nil
}

// UstawZakladki zakłada drzewo nawigacji dokumentu.
//
// Drzewo podane zastępuje zastane w całości, a nie dokłada się do niego:
// komenda nazywa się „ustaw", a dokładanie zostawiłoby Operatora bez sposobu na
// usunięcie zakładki raz założonej.
func (a *adapterPdfStudia) UstawZakladki(ctx context.Context,
	z shared.StudioPdfBookmarksSetRequest) (shared.StudioPdfBookmarksSetResponse, error) {

	bajty, err := a.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioPdfBookmarksSetResponse{}, err
	}
	if len(z.Bookmarks) == 0 {
		return shared.StudioPdfBookmarksSetResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"żądanie nie niesie ani jednej zakładki — nie ma czego ustawić")
	}
	zakladki, err := zakladkiBiblioteki(z.Bookmarks)
	if err != nil {
		return shared.StudioPdfBookmarksSetResponse{}, err
	}

	var wynik bytes.Buffer
	if err := api.AddBookmarks(bytes.NewReader(bajty), &wynik, zakladki, true,
		nastawyPdf()); err != nil {
		return shared.StudioPdfBookmarksSetResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"nie można założyć zakładek: "+err.Error())
	}
	zasob, err := a.odlozPdf(ctx, wynik.Bytes(), "dokument z zakładkami",
		wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioPdfBookmarksSetResponse{}, err
	}
	return shared.StudioPdfBookmarksSetResponse{Asset: zasob}, nil
}

// WypelnijFormularz odczytuje pola formularza, a przy podanych wartościach
// wypełnia je i oddaje dokument wypełniony; żądanie bez wartości jest samym
// odczytem.
func (a *adapterPdfStudia) WypelnijFormularz(ctx context.Context,
	z shared.StudioPdfFormFillRequest) (shared.StudioPdfFormFillResponse, error) {

	bajty, err := a.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioPdfFormFillResponse{}, err
	}
	pola, err := api.FormFields(bytes.NewReader(bajty), nastawyPdf())
	if err != nil {
		return shared.StudioPdfFormFillResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"dokument nie ma czytelnego formularza: "+err.Error())
	}

	if len(z.Values) == 0 || string(z.Values) == "null" {
		return shared.StudioPdfFormFillResponse{Fields: polaFormularza(pola)}, nil
	}

	// Wypełnienie idzie przez opis wyprowadzony z dokumentu — biblioteka
	// przyjmuje wartości tylko w nim.
	opis, err := api.ExportForm(bytes.NewReader(bajty), "dokument", nastawyPdf())
	if err != nil || opis == nil || len(opis.Forms) == 0 {
		return shared.StudioPdfFormFillResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"dokument nie ma formularza, który dałoby się wypełnić")
	}
	var wartosci map[string]string
	if err := json.Unmarshal(z.Values, &wartosci); err != nil {
		return shared.StudioPdfFormFillResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"wartości pól podaje się jako zbiór par nazwa-wartość: "+err.Error())
	}
	nieznane := wpiszWartosci(&opis.Forms[0], wartosci)
	if len(nieznane) > 0 {
		return shared.StudioPdfFormFillResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"formularz nie ma pól: "+strings.Join(nieznane, ", "))
	}

	opisJson, err := json.Marshal(opis)
	if err != nil {
		return shared.StudioPdfFormFillResponse{}, odmowaPdf(shared.ErrorCodeInternalError,
			"nie można złożyć opisu wypełnienia: "+err.Error())
	}
	var wynik bytes.Buffer
	if err := api.FillForm(bytes.NewReader(bajty), bytes.NewReader(opisJson), &wynik,
		nastawyPdf()); err != nil {
		return shared.StudioPdfFormFillResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"wypełnienie formularza nie powiodło się: "+err.Error())
	}
	wypelniony := wynik.Bytes()

	// Utrwalenie zamyka pola, nie usuwa ich — usunięcie zabrałoby razem z polem
	// wpisaną treść.
	if z.Flatten != nil && *z.Flatten {
		nazwy := make([]string, 0, len(pola))
		for _, pole := range pola {
			nazwy = append(nazwy, pole.Name)
		}
		var zamkniety bytes.Buffer
		if err := api.LockFormFields(bytes.NewReader(wypelniony), &zamkniety, nazwy,
			nastawyPdf()); err != nil {
			return shared.StudioPdfFormFillResponse{}, odmowaPdf(shared.ErrorCodeInternalError,
				"nie można utrwalić wypełnienia: "+err.Error())
		}
		wypelniony = zamkniety.Bytes()
	}

	zasob, err := a.odlozPdf(ctx, wypelniony, "formularz wypełniony",
		wartoscTekstu(z.WindowId))
	if err != nil {
		return shared.StudioPdfFormFillResponse{}, err
	}
	poWypelnieniu, err := api.FormFields(bytes.NewReader(wypelniony), nastawyPdf())
	if err != nil {
		poWypelnieniu = pola
	}
	return shared.StudioPdfFormFillResponse{
		Fields: polaFormularza(poWypelnieniu), Asset: &zasob,
	}, nil
}

// wpiszWartosci wstawia wartości Operatora do opisu formularza i oddaje nazwy,
// których w formularzu nie ma.
func wpiszWartosci(formularz *form.Form, wartosci map[string]string) []string {
	nieznane := make([]string, 0)
	for nazwa, wartosc := range wartosci {
		trafione := false
		for _, pole := range formularz.TextFields {
			if pole != nil && (pole.Name == nazwa || pole.ID == nazwa) {
				pole.Value, trafione = wartosc, true
			}
		}
		for _, pole := range formularz.DateFields {
			if pole != nil && (pole.Name == nazwa || pole.ID == nazwa) {
				pole.Value, trafione = wartosc, true
			}
		}
		for _, pole := range formularz.ComboBoxes {
			if pole != nil && (pole.Name == nazwa || pole.ID == nazwa) {
				pole.Value, trafione = wartosc, true
			}
		}
		for _, pole := range formularz.ListBoxes {
			if pole != nil && (pole.Name == nazwa || pole.ID == nazwa) {
				pole.Values, trafione = []string{wartosc}, true
			}
		}
		for _, pole := range formularz.CheckBoxes {
			if pole != nil && (pole.Name == nazwa || pole.ID == nazwa) {
				pole.Value, trafione = wartosc == "true", true
			}
		}
		for _, pole := range formularz.RadioButtonGroups {
			if pole != nil && (pole.Name == nazwa || pole.ID == nazwa) {
				pole.Value, trafione = wartosc, true
			}
		}
		if !trafione {
			nieznane = append(nieznane, nazwa)
		}
	}
	sort.Strings(nieznane)
	return nieznane
}

// polaFormularza przenosi pola formularza z postaci biblioteki pdfcpu na postać
// kontraktu, niosąc nazwę, rodzaj, wartość i opcje pola.
func polaFormularza(zrodlo []form.Field) []shared.StudioPdfFormField {
	pola := make([]shared.StudioPdfFormField, 0, len(zrodlo))
	for _, pole := range zrodlo {
		przeniesione := shared.StudioPdfFormField{
			Name: pole.Name,
			Kind: pole.Typ.String(),
		}
		if pole.V != "" {
			wartosc := pole.V
			przeniesione.Value = &wartosc
		}
		if pole.Opts != "" {
			przeniesione.Options = strings.Split(pole.Opts, ",")
		}
		zamkniete := pole.Locked
		przeniesione.ReadOnly = &zamkniete
		if len(pole.Pages) > 0 {
			strona := pole.Pages[0]
			przeniesione.Page = &strona
		}
		pola = append(pola, przeniesione)
	}
	return pola
}

// WyciagnijZPdf wyjmuje z dokumentu osadzone obrazy i załączniki, zakładając
// każdemu osobny zasób magazynu; brak wskazania znaczy wyciągnięcie obu
// rodzajów.
func (a *adapterPdfStudia) WyciagnijZPdf(ctx context.Context,
	z shared.StudioPdfExtractRequest) (shared.StudioPdfExtractResponse, error) {

	bajty, err := a.trescMaterialu(ctx, z.AssetId)
	if err != nil {
		return shared.StudioPdfExtractResponse{}, err
	}
	obrazy := z.Images == nil || *z.Images
	zalaczniki := z.Attachments == nil || *z.Attachments
	if !obrazy && !zalaczniki {
		return shared.StudioPdfExtractResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
			"żądanie odrzuca i obrazy, i załączniki — nie ma czego wyciągnąć")
	}
	okno := wartoscTekstu(z.WindowId)
	kody := make([]string, 0)

	if obrazy {
		strony, err := api.ExtractImagesRaw(bytes.NewReader(bajty), nil, nastawyPdf())
		if err != nil {
			return shared.StudioPdfExtractResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
				"nie można wyciągnąć obrazów: "+err.Error())
		}
		for _, strona := range strony {
			for numer, obraz := range strona {
				tresc, err := io.ReadAll(obraz)
				if err != nil {
					return shared.StudioPdfExtractResponse{}, odmowaPdf(
						shared.ErrorCodeInternalError,
						"nie można odczytać obrazu ze strony "+strconv.Itoa(numer)+": "+err.Error())
				}
				zasob, err := a.odlozTresc(ctx, tresc,
					pierwszyNiepusty(obraz.Name, "obraz strony "+strconv.Itoa(numer)),
					pierwszyNiepusty(obraz.FileType, "bin"), okno)
				if err != nil {
					return shared.StudioPdfExtractResponse{}, err
				}
				kody = append(kody, zasob.Id)
			}
		}
	}

	if zalaczniki {
		// Biblioteka przy odczycie surowym nie dotyka katalogu wyjściowego:
		// treść wraca strumieniem.
		pozycje, err := api.ExtractAttachmentsRaw(bytes.NewReader(bajty), "", nil, nastawyPdf())
		if err != nil {
			return shared.StudioPdfExtractResponse{}, odmowaPdf(shared.ErrorCodeValidationFailed,
				"nie można wyciągnąć załączników: "+err.Error())
		}
		for _, pozycja := range pozycje {
			tresc, err := io.ReadAll(pozycja)
			if err != nil {
				return shared.StudioPdfExtractResponse{}, odmowaPdf(shared.ErrorCodeInternalError,
					"nie można odczytać załącznika "+pozycja.FileName+": "+err.Error())
			}
			format := "bin"
			if kropka := strings.LastIndex(pozycja.FileName, "."); kropka >= 0 {
				format = strings.ToLower(pozycja.FileName[kropka+1:])
			}
			zasob, err := a.odlozTresc(ctx, tresc,
				pierwszyNiepusty(pozycja.FileName, "załącznik"), format, okno)
			if err != nil {
				return shared.StudioPdfExtractResponse{}, err
			}
			kody = append(kody, zasob.Id)
		}
	}

	return shared.StudioPdfExtractResponse{AssetIds: kody, Extracted: len(kody)}, nil
}

// pierwszyNiepusty oddaje pierwsze wskazanie o niepustej treści z podanego
// ciągu, służąc jako nazwa zapasowa, gdy wskazanie pierwotne jest puste.
func pierwszyNiepusty(wskazania ...string) string {
	for _, wskazanie := range wskazania {
		if strings.TrimSpace(wskazanie) != "" {
			return wskazanie
		}
	}
	return ""
}

// ── Rejestracja ─────────────────────────────────────────────────────────────

// WarsztatPdf wypełnia rodzinę komend `studio.pdf.*` w części zbudowanej
// biblioteką pdfcpu: scalanie, podział, stemplowanie i odchudzanie.
type WarsztatPdf interface {
	ScalPdf(ctx context.Context, z shared.StudioPdfMergeRequest) (shared.StudioPdfMergeResponse, error)
	PodzielPdf(ctx context.Context, z shared.StudioPdfSplitRequest) (shared.StudioPdfSplitResponse, error)
	OdchudzPdf(ctx context.Context, z shared.StudioPdfOptimizeRequest) (shared.StudioPdfOptimizeResponse, error)
	UlozStrony(ctx context.Context, z shared.StudioPdfPagesReorderRequest) (shared.StudioPdfPagesReorderResponse, error)
	OstemplujPdf(ctx context.Context, z shared.StudioPdfStampRequest) (shared.StudioPdfStampResponse, error)
	NumerujPdf(ctx context.Context, z shared.StudioPdfBatesRequest) (shared.StudioPdfBatesResponse, error)
	UstawZakladki(ctx context.Context, z shared.StudioPdfBookmarksSetRequest) (shared.StudioPdfBookmarksSetResponse, error)
	WypelnijFormularz(ctx context.Context, z shared.StudioPdfFormFillRequest) (shared.StudioPdfFormFillResponse, error)
	WyciagnijZPdf(ctx context.Context, z shared.StudioPdfExtractRequest) (shared.StudioPdfExtractResponse, error)
}

// zarejestrujWarsztatPdf wpina czynności warsztatu dokumentu PDF do rejestru
// komend, wiążąc każdą nazwę komendy z metodą warsztatu.
func zarejestrujWarsztatPdf(r *Rejestr, w WarsztatPdf) {
	if r == nil || w == nil {
		return
	}
	r.Zarejestruj(shared.CommandStudioPdfMerge, obsluz(w.ScalPdf))
	r.Zarejestruj(shared.CommandStudioPdfSplit, obsluz(w.PodzielPdf))
	r.Zarejestruj(shared.CommandStudioPdfOptimize, obsluz(w.OdchudzPdf))
	r.Zarejestruj(shared.CommandStudioPdfPagesReorder, obsluz(w.UlozStrony))
	r.Zarejestruj(shared.CommandStudioPdfStamp, obsluz(w.OstemplujPdf))
	r.Zarejestruj(shared.CommandStudioPdfBates, obsluz(w.NumerujPdf))
	r.Zarejestruj(shared.CommandStudioPdfBookmarksSet, obsluz(w.UstawZakladki))
	r.Zarejestruj(shared.CommandStudioPdfFormFill, obsluz(w.WypelnijFormularz))
	r.Zarejestruj(shared.CommandStudioPdfExtract, obsluz(w.WyciagnijZPdf))
}
