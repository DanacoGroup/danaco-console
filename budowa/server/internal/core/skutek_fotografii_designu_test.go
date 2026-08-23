package core

import (
	"bytes"
	"compress/zlib"
	"context"
	"image"
	"image/color"
	"image/png"
	"io"
	"testing"

	"danacoconsole/shared"
)

// Skutek warsztatu fotografii i części drukarskiej modułu Design — mierzony
// w PIKSELACH i w STRONACH, nie w kopercie.
//
// Żaden sprawdzian tego pliku nie kończy się na tym, że odpowiedź jest udana.
// Każdy schodzi po odwołaniu zasobu do magazynu, rozkłada plik i pyta go
// o rzeczy, których koperta nie zna: ile ma pikseli, czy niesie kanał krycia
// i jakie ma w nim wartości, o ile przesunęła się średnia jasność, ile kafli
// powstało i o jakich wymiarach, ile stron ma wydany plik PDF.
//
// Powód jest zapisany w historii tego produktu: `design.asset.generate` meldował
// kiedyś `status: ok` z wykazem zasobów, za którymi nie było ani jednego bajtu.
// Sprawdzian zaglądający w `status` świecił wtedy zielono. Odtąd sprawdzian
// obszaru Design mierzy SKUTEK.

// wniesObrazSprawdzianu wnosi obraz do magazynu okna i oddaje zasób kontraktu.
//
// Droga jest drogą Operatora (`design.asset.upload` z treścią w base64), a nie
// zapisem wprost do bazy: sprawdzian ma mierzyć to, co dzieje się w produkcie,
// a nie stan, który sam sobie ustawił obok produktu.
func wniesObrazSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno string, bajty []byte) shared.DesignAsset {

	t.Helper()

	var odpowiedz shared.DesignAssetUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetUpload,
		shared.DesignAssetUploadRequest{
			WindowId:      okno,
			Name:          wskaznik("materiał sprawdzianu"),
			Kind:          shared.DesignAssetKindImage,
			ContentBase64: wskaznik(wBase64(bajty)),
		}, &odpowiedz)
	if odpowiedz.Asset.Uri == nil {
		t.Fatal("wniesiony zasób nie ma odwołania do treści")
	}
	return odpowiedz.Asset
}

// obrazZMagazynu rozkłada plik leżący pod odwołaniem zasobu.
//
// To jest sedno pomiaru skutku: obraz przychodzi z DYSKU, nie z odpowiedzi
// komendy. Odpowiedź może mówić o wymiarach cokolwiek — plik mówi prawdę.
func obrazZMagazynuFotografii(t *testing.T, katalog string, zasob shared.DesignAsset) image.Image {
	t.Helper()

	if zasob.Uri == nil {
		t.Fatal("zasób bez odwołania — nie ma czego zmierzyć")
	}
	bajty := bajtyPodOdwolaniem(t, katalog, *zasob.Uri)
	obraz, _, err := image.Decode(bytes.NewReader(bajty))
	if err != nil {
		t.Fatalf("treść pod odwołaniem %q nie jest obrazem: %v", *zasob.Uri, err)
	}
	return obraz
}

// obrazJednolityPNG składa obraz o jednej barwie — materiał, na którym odcięcie
// tła ma zadziałanie jednoznaczne i sprawdzalne.
func obrazJednolityPNG(t *testing.T, szerokosc, wysokosc int, barwa color.RGBA) []byte {
	t.Helper()

	plotno := image.NewRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			plotno.SetRGBA(x, y, barwa)
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	return bufor.Bytes()
}

// obrazZPrzedmiotemPNG składa obraz o jednolitym tle z prostokątem w środku —
// materiał, na którym odcięcie tła MA co zostawić.
func obrazZPrzedmiotemPNG(t *testing.T, bok int) []byte {
	t.Helper()

	plotno := image.NewRGBA(image.Rect(0, 0, bok, bok))
	tlo := color.RGBA{R: 0xf8, G: 0xf8, B: 0xf8, A: 0xff}
	przedmiot := color.RGBA{R: 0x10, G: 0x20, B: 0x80, A: 0xff}
	for y := 0; y < bok; y++ {
		for x := 0; x < bok; x++ {
			if x > bok/4 && x < 3*bok/4 && y > bok/4 && y < 3*bok/4 {
				plotno.SetRGBA(x, y, przedmiot)
				continue
			}
			plotno.SetRGBA(x, y, tlo)
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	return bufor.Bytes()
}

// sredniaJasnoscPliku liczy średnią jasność obrazu odczytanego z dysku.
func sredniaJasnoscPliku(obraz image.Image) float64 {
	granice := obraz.Bounds()
	suma, punktow := 0.0, 0
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			r, g, b, _ := obraz.At(x, y).RGBA()
			suma += float64(int(r>>8)*299+int(g>>8)*587+int(b>>8)*114) / 1000
			punktow++
		}
	}
	if punktow == 0 {
		return 0
	}
	return suma / float64(punktow)
}

// TestPowiekszenieZdjeciaMaZmierzonaRozdzielczoscWPliku mierzy PIKSELE wyniku,
// nie pole `width` odpowiedzi: powiększenie, które melduje 4× i zapisuje plik
// o rozmiarze źródła, jest dokładnie tą szkodą, przed którą broni ten plik.
func TestPowiekszenieZdjeciaMaZmierzonaRozdzielczoscWPliku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	zrodlo := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-fotografii",
		obrazPNG(t, 20, 15))

	var wynik shared.DesignPhotoUpscaleResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoUpscale,
		shared.DesignPhotoUpscaleRequest{AssetId: zrodlo.Id, Factor: 4}, &wynik)

	obraz := obrazZMagazynuFotografii(t, katalog, wynik.Asset)
	granice := obraz.Bounds()
	if granice.Dx() != 80 || granice.Dy() != 60 {
		t.Errorf("plik po powiększeniu ×4 ma %d×%d pikseli, a źródło 20×15 daje 80×60",
			granice.Dx(), granice.Dy())
	}
	// Odpowiedź ma mówić o tym samym, co plik. Rozjazd znaczy, że jedno z dwojga
	// kłamie, i wtedy nie da się powiedzieć które.
	if wynik.Width != granice.Dx() || wynik.Height != granice.Dy() {
		t.Errorf("odpowiedź podaje %d×%d, a plik ma %d×%d",
			wynik.Width, wynik.Height, granice.Dx(), granice.Dy())
	}
	if wynik.ComputedBy == "" {
		t.Error("odpowiedź nie mówi, którą drogą rdzeń policzył wynik — pole computedBy jest puste")
	}
	// Wariant wskazuje źródło: bez tego łańcuch edycji nie ma jak wrócić do
	// zdjęcia, które Operator wniósł.
	if wynik.Asset.VariantOfAssetId == nil || *wynik.Asset.VariantOfAssetId != zrodlo.Id {
		t.Errorf("wynik nie wskazuje źródła jako wariantu (%v), a źródłem jest %s",
			wynik.Asset.VariantOfAssetId, zrodlo.Id)
	}
}

// TestOdcieciaTlaZostawiaKanalKryciaWPliku mierzy OBECNOŚĆ i TREŚĆ kanału krycia
// w pliku: punkt tła ma mieć krycie zerowe, punkt przedmiotu — pełne. Odpowiedź
// `hasAlpha: true` nad plikiem bez przezroczystości byłaby kopertą bez skutku.
func TestOdcieciaTlaZostawiaKanalKryciaWPliku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	zrodlo := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-fotografii",
		obrazZPrzedmiotemPNG(t, 40))

	var wynik shared.DesignPhotoBackgroundRemoveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoBackgroundRemove,
		shared.DesignPhotoBackgroundRemoveRequest{AssetId: zrodlo.Id}, &wynik)

	if !wynik.HasAlpha {
		t.Fatal("odpowiedź mówi, że wynik nie ma kanału krycia — odcięcie tła bez kanału krycia " +
			"nie jest odcięciem tła")
	}
	if wynik.TransparentShare <= 0 || wynik.TransparentShare >= 1 {
		t.Errorf("zmierzony udział punktów przezroczystych to %v; obraz z przedmiotem na tle "+
			"ma mieć go między zerem a jednością", wynik.TransparentShare)
	}

	obraz := obrazZMagazynuFotografii(t, katalog, wynik.Asset)
	// Naroże jest tłem — po odcięciu ma być przezroczyste.
	_, _, _, kryciaNaroza := obraz.At(1, 1).RGBA()
	if kryciaNaroza != 0 {
		t.Errorf("punkt tła (1;1) ma krycie %d, a po odcięciu tła ma mieć zero", kryciaNaroza)
	}
	// Środek jest przedmiotem — po odcięciu ma zostać widoczny.
	_, _, _, kryciaSrodka := obraz.At(20, 20).RGBA()
	if kryciaSrodka == 0 {
		t.Error("punkt przedmiotu (20;20) zniknął razem z tłem — wynik jest pustym płótnem")
	}

	// Udział zmierzony w pliku ma się zgadzać z tym, co powiedziała odpowiedź.
	// Pomiar niezależny jest sednem: liczba w odpowiedzi mogła powstać z czegoś
	// innego niż z pikseli, które wyszły.
	granice := obraz.Bounds()
	przezroczystych := 0
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			if _, _, _, krycie := obraz.At(x, y).RGBA(); krycie == 0 {
				przezroczystych++
			}
		}
	}
	zmierzony := float64(przezroczystych) / float64(granice.Dx()*granice.Dy())
	if roznica := zmierzony - wynik.TransparentShare; roznica > 0.05 || roznica < -0.05 {
		t.Errorf("plik ma %.3f punktów przezroczystych, a odpowiedź mówi %.3f",
			zmierzony, wynik.TransparentShare)
	}
}

// TestKorekcjaBarwyPrzesuwaHistogramWPliku mierzy PRZESUNIĘCIE średniej jasności
// pliku wynikowego wobec źródła — a nie to, że komenda się udała. Korekcja
// meldująca powodzenie i zapisująca plik identyczny ze źródłem jest pokrętłem
// podłączonym donikąd.
func TestKorekcjaBarwyPrzesuwaHistogramWPliku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	// Szarość średnia: rozjaśnienie ma na niej gdzie zadziałać w obie strony.
	zrodloBajty := obrazJednolityPNG(t, 24, 24, color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff})
	zrodlo := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-fotografii", zrodloBajty)

	var wynik shared.DesignPhotoColorCorrectResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoColorCorrect,
		shared.DesignPhotoColorCorrectRequest{
			AssetId:    zrodlo.Id,
			Brightness: wskaznik(0.3),
		}, &wynik)

	przed := sredniaJasnoscPliku(obrazZMagazynuFotografii(t, katalog, zrodlo))
	po := sredniaJasnoscPliku(obrazZMagazynuFotografii(t, katalog, wynik.Asset))
	if po <= przed {
		t.Errorf("po rozjaśnieniu średnia jasność pliku to %.2f, a przed była %.2f — "+
			"korekcja nie ruszyła pikseli", po, przed)
	}
	// Odpowiedź podaje przesunięcie; ma się zgadzać z pomiarem na plikach.
	if zmierzone := po - przed; zmierzone-wynik.HistogramShift > 2 ||
		wynik.HistogramShift-zmierzone > 2 {

		t.Errorf("odpowiedź podaje przesunięcie %.2f, a pliki różnią się o %.2f",
			wynik.HistogramShift, zmierzone)
	}
	// Źródło ZOSTAJE nietknięte — to jest warunek łańcucha edycji.
	if przed < 120 || przed > 140 {
		t.Errorf("średnia jasność ŹRÓDŁA to %.2f, a wniesiono szarość 0x80 (128) — "+
			"obróbka nadpisała oryginał", przed)
	}
}

// TestPodzialWielkoformatowyOddajeKafleOZmierzonychWymiarach mierzy liczbę kafli
// i ich wymiary, a każdy kafel schodzi do magazynu po własne bajty. Wykaz kafli
// bez plików byłby spisem prostokątów, których drukarnia nie dostanie.
func TestPodzialWielkoformatowyOddajeKafleOZmierzonychWymiarach(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	zrodlo := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-druku", obrazPNG(t, 240, 120))

	var wynik shared.DesignLargeformatTileResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignLargeformatTile,
		shared.DesignLargeformatTileRequest{
			AssetId:        zrodlo.Id,
			TileWidthMm:    100,
			TileHeightMm:   100,
			TargetWidthMm:  wskaznik(300.0),
			TargetHeightMm: wskaznik(200.0),
		}, &wynik)

	// Materiał 300×200 mm na kafle 100×100 mm bez zakładki to trzy kolumny i dwa
	// wiersze — liczba, którą da się policzyć w głowie i sprawdzić w wyniku.
	if wynik.Columns != 3 || wynik.Rows != 2 {
		t.Errorf("podział dał %d kolumn i %d wierszy; 300×200 mm na kafle 100×100 mm daje 3×2",
			wynik.Columns, wynik.Rows)
	}
	if len(wynik.Tiles) != 6 {
		t.Fatalf("podział oddał %d kafli, a układ 3×2 ma sześć", len(wynik.Tiles))
	}
	if wynik.EffectiveDpi == nil || *wynik.EffectiveDpi <= 0 {
		t.Error("odpowiedź nie podaje zmierzonej rozdzielczości skutecznej")
	}
	for _, kafel := range wynik.Tiles {
		if kafel.AssetId == nil {
			t.Fatalf("kafel %d-%d nie ma zasobu — spis prostokątów bez plików",
				kafel.Row, kafel.Column)
		}
		var wykaz shared.DesignAssetListResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetList,
			shared.DesignAssetListRequest{WindowId: wskaznik("okno-druku")}, &wykaz)
		znaleziony := false
		for _, zasob := range wykaz.Assets {
			if zasob.Id != *kafel.AssetId {
				continue
			}
			znaleziony = true
			obraz := obrazZMagazynuFotografii(t, katalog, zasob)
			if obraz.Bounds().Dx() < 1 || obraz.Bounds().Dy() < 1 {
				t.Errorf("kafel %d-%d ma plik o zerowym boku", kafel.Row, kafel.Column)
			}
		}
		if !znaleziony {
			t.Errorf("kafla %d-%d nie ma w wykazie zasobów okna", kafel.Row, kafel.Column)
		}
		if kafel.WidthMm <= 0 || kafel.HeightMm <= 0 {
			t.Errorf("kafel %d-%d ma wymiar %v×%v mm", kafel.Row, kafel.Column,
				kafel.WidthMm, kafel.HeightMm)
		}
	}
}

// TestWydaniePublikacjiMaZmierzonaLiczbeStronWPliku mierzy LICZBĘ STRON w pliku
// PDF, nie pole odpowiedzi: publikacja czterostronicowa wydana jako jedna strona
// wraca z drukarni jako jedna strona, a odpowiedź mówiłaby wtedy „cztery".
func TestWydaniePublikacjiMaZmierzonaLiczbeStronWPliku(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	// Materiał stron: jeden obraz na każdą, żeby każda strona miała bajty.
	zasob := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-publikacji",
		obrazPNG(t, 60, 84))

	warstwa := func(numer int) []shared.DesignBoardLayer {
		return []shared.DesignBoardLayer{{
			Id:      "",
			AssetId: wskaznik(zasob.Id),
			X:       wskaznik(0.0),
			Y:       wskaznik(0.0),
			Width:   wskaznik(210.0),
			Height:  wskaznik(297.0),
			Note:    wskaznik("strona " + string(rune('0'+numer))),
		}}
	}

	var zapis shared.DesignTemplateSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignTemplateSave,
		shared.DesignTemplateSaveRequest{
			WindowId: "okno-publikacji",
			Name:     "broszura sprawdzianu",
			Kind:     shared.DesignTemplateKindPrint,
			Width:    210,
			Height:   297,
			Pages: []shared.DesignTemplatePage{
				{Number: 1, Name: wskaznik("okładka"), Layers: warstwa(1)},
				{Number: 2, Layers: warstwa(2)},
				{Number: 3, Layers: warstwa(3)},
				{Number: 4, Name: wskaznik("okładka tylna"), Layers: warstwa(4)},
			},
		}, &zapis)

	// Kontrola przeddrukowa pominięta świadomie: materiałem stron jest obraz
	// 60×84 px rozciągnięty na arkusz A4, więc rozdzielczość skuteczna jest niska
	// i kontrola słusznie by ją zablokowała. Sprawdzian mierzy tu LICZBĘ STRON,
	// a pominięcie wraca w odpowiedzi i jest niżej sprawdzone.
	var wydanie shared.DesignPrintExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPrintExport,
		shared.DesignPrintExportRequest{
			Format:        "pdf",
			TemplateId:    wskaznik(zapis.Template.Id),
			SkipPreflight: wskaznik(true),
			WindowId:      wskaznik("okno-publikacji"),
			Binding:       wskaznik(shared.DesignPrintBinding(shared.DesignPrintBindingZeszytowa)),
		}, &wydanie)

	if !wydanie.PreflightSkipped {
		t.Error("odpowiedź nie mówi, że kontrolę przeddrukową pominięto — pominięcie ma wracać")
	}
	if wydanie.PageCount == nil || *wydanie.PageCount != 4 {
		t.Errorf("odpowiedź podaje %v stron, a publikacja ma cztery", wydanie.PageCount)
	}

	// POMIAR NIEZALEŻNY: liczba stron odczytana z bajtów pliku PDF.
	bajty := bajtyPodOdwolaniem(t, katalog, *wydanie.Asset.Uri)
	if stron := liczbaStronPdfSprawdzianu(bajty); stron != 4 {
		t.Errorf("plik PDF ma %d stron, a publikacja miała cztery (rozmiar pliku: %d B)",
			stron, len(bajty))
	}
}

// liczbaStronPdfSprawdzianu liczy strony pliku PDF, licząc obiekty typu `/Page`
// w treści pliku ORAZ w jego strumieniach obiektów po rozpakowaniu.
//
// Rachunek jest własny i celowo nie sięga po bibliotekę, którą rdzeń plik złożył.
// Gdyby liczył strony tą samą biblioteką, mierzyłby zgodność biblioteki z samą
// sobą — a to jest właśnie ten rodzaj sprawdzianu, który przepuszcza szkodę:
// błąd w składaniu pliku i błąd w jego odczycie zniosłyby się wzajemnie.
//
// Rozpakowanie strumieni jest konieczne, bo `pdfcpu` zapisuje katalog obiektów
// w STRUMIENIACH OBIEKTÓW skompresowanych metodą Flate (PDF 1.5 i wyżej). Napisu
// `/Type /Page` nie ma wtedy w pliku wprost. Pakiet `compress/zlib` biblioteki
// standardowej wystarcza — to ta sama kompresja.
//
// Wzorzec bez ukośnika po `Page` odróżnia stronę od drzewa stron (`/Type /Pages`),
// które w pliku występuje raz.
func liczbaStronPdfSprawdzianu(bajty []byte) int {
	stron := zlicznikStronWTresciSprawdzianu(bajty)
	for _, strumien := range strumienieFlateSprawdzianu(bajty) {
		stron += zlicznikStronWTresciSprawdzianu(strumien)
	}
	return stron
}

// zlicznikStronWTresciSprawdzianu liczy wystąpienia typu `/Page` w treści.
func zlicznikStronWTresciSprawdzianu(tresc []byte) int {
	stron := 0
	for _, wzor := range [][]byte{[]byte("/Type /Page"), []byte("/Type/Page")} {
		reszta := tresc
		for {
			numer := bytes.Index(reszta, wzor)
			if numer < 0 {
				break
			}
			po := reszta[numer+len(wzor):]
			// Znak następny rozstrzyga: litera `s` znaczy `/Pages`, czyli drzewo
			// stron, a nie stronę.
			if len(po) > 0 && po[0] != 's' {
				stron++
			}
			reszta = po
		}
	}
	return stron
}

// strumienieFlateSprawdzianu wyciąga z pliku PDF treść strumieni rozpakowywalnych
// metodą Flate.
//
// Rozbiór jest prymitywny i taki ma być: szuka par `stream` / `endstream`
// i próbuje rozpakować każdą. Strumień, którego nie da się rozpakować (obraz JPEG,
// treść nieskompresowana), jest po prostu pomijany — sprawdzian szuka katalogu
// obiektów, a nie wszystkiego, co w pliku leży.
func strumienieFlateSprawdzianu(bajty []byte) [][]byte {
	strumienie := [][]byte{}
	reszta := bajty
	for {
		poczatek := bytes.Index(reszta, []byte("stream"))
		if poczatek < 0 {
			return strumienie
		}
		tresc := reszta[poczatek+len("stream"):]
		// Po słowie `stream` idzie koniec wiersza — CRLF albo LF.
		tresc = bytes.TrimPrefix(tresc, []byte("\r"))
		tresc = bytes.TrimPrefix(tresc, []byte("\n"))
		koniec := bytes.Index(tresc, []byte("endstream"))
		if koniec < 0 {
			return strumienie
		}
		if rozpakowany, err := rozpakujFlateSprawdzianu(tresc[:koniec]); err == nil {
			strumienie = append(strumienie, rozpakowany)
		}
		reszta = tresc[koniec+len("endstream"):]
	}
}

// rozpakujFlateSprawdzianu rozpakowuje strumień metodą Flate.
func rozpakujFlateSprawdzianu(tresc []byte) ([]byte, error) {
	czytnik, err := zlib.NewReader(bytes.NewReader(tresc))
	if err != nil {
		return nil, err
	}
	defer czytnik.Close()
	return io.ReadAll(czytnik)
}

// TestKontrolaPrzeddrukowaOdmawiaWydaniaPrzyWadzieOWadzeBledu pilnuje
// rozstrzygnięcia, które kosztuje nakład: plik nie do druku NIE wychodzi jako
// gotowy do druku.
//
// Materiałem jest obraz o rozdzielczości skutecznej rażąco poniżej progu —
// wtedy kontrola ma dać wadę o wadze błędu, a wydanie ma ODMÓWIĆ. Sprawdzian
// mierzy tu odmowę, bo odmowa jest tu funkcją.
func TestKontrolaPrzeddrukowaOdmawiaWydaniaPrzyWadzieOWadzeBledu(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	zasob := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-druku", obrazPNG(t, 20, 20))

	// Kompozycja o wymiarach arkusza A4 w milimetrach: warstwa 210 mm szerokości
	// z obrazu 20 px daje około 2 dpi, czyli grubo poniżej połowy progu 300 dpi.
	var kompozycja shared.DesignBoardUpdateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignBoardUpdate,
		shared.DesignBoardUpdateRequest{
			WindowId: "okno-druku",
			Name:     wskaznik("arkusz sprawdzianu"),
			Layers: []shared.DesignBoardLayer{{
				AssetId: wskaznik(zasob.Id),
				X:       wskaznik(0.0),
				Y:       wskaznik(0.0),
				Width:   wskaznik(210.0),
				Height:  wskaznik(297.0),
			}},
		}, &kompozycja)

	var kontrola shared.DesignPrintPreflightResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPrintPreflight,
		shared.DesignPrintPreflightRequest{BoardId: wskaznik(kompozycja.Board.Id)}, &kontrola)

	if kontrola.Errors == 0 {
		t.Fatalf("kontrola nie znalazła ani jednej wady o wadze błędu na materiale o około 2 dpi "+
			"(ostrzeżeń: %d, zastrzeżeń łącznie: %d)", kontrola.Warnings, len(kontrola.Issues))
	}
	if kontrola.Ready {
		t.Error("kontrola mówi, że materiał jest gotowy do druku, mając wadę o wadze błędu")
	}
	// Zastrzeżenia idą w kolejności wagi — pierwsze ma być błędem.
	if len(kontrola.Issues) > 0 &&
		kontrola.Issues[0].Severity != shared.DesignPreflightSeverityBlad {
		t.Errorf("pierwsze zastrzeżenie ma wagę %q, a wykaz ma iść od wad o wadze błędu",
			kontrola.Issues[0].Severity)
	}

	// Wydanie bez pominięcia kontroli ma ODMÓWIĆ. Odmowa jest tu wynikiem
	// oczekiwanym: plik nie do druku wydany jako gotowy do druku jedzie do
	// drukarni i kosztuje nakład.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandDesignPrintExport,
		shared.DesignPrintExportRequest{
			Format:   "pdf",
			BoardId:  wskaznik(kompozycja.Board.Id),
			WindowId: wskaznik("okno-druku"),
		})
	if odmowa.Message == "" {
		t.Error("odmowa wydania nie nazywa powodu — Operator ma dostać zdanie, nie sam kod")
	}

	// Pominięcie kontroli jest jawnym wyborem Operatora i ma się udać — inaczej
	// funkcja byłaby uprzejmą odmową bez drogi wyjścia.
	var zPominieciem shared.DesignPrintExportResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPrintExport,
		shared.DesignPrintExportRequest{
			Format:        "pdf",
			BoardId:       wskaznik(kompozycja.Board.Id),
			WindowId:      wskaznik("okno-druku"),
			SkipPreflight: wskaznik(true),
		}, &zPominieciem)
	if !zPominieciem.PreflightSkipped {
		t.Error("wydanie z pominiętą kontrolą nie mówi o pominięciu — nikt potem nie powie, " +
			"że nie wiedział")
	}
}

// TestLancuchEdycjiOddajeCzynnosciWKolejnosciWykonania mierzy, czy łańcuch
// edycji naprawdę zapisuje CZYNNOŚĆ i jej nastawy — bez tego
// `design.photo.history.get` byłby wykazem obrazków bez słowa o tym, co je różni.
func TestLancuchEdycjiOddajeCzynnosciWKolejnosciWykonania(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	zrodlo := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-fotografii", obrazPNG(t, 32, 32))

	var poKadrze shared.DesignPhotoCropResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoCrop,
		shared.DesignPhotoCropRequest{
			AssetId: zrodlo.Id,
			X:       wskaznik(4.0),
			Y:       wskaznik(4.0),
			Width:   wskaznik(16.0),
			Height:  wskaznik(16.0),
		}, &poKadrze)

	var poFiltrze shared.DesignPhotoFilterApplyResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoFilterApply,
		shared.DesignPhotoFilterApplyRequest{
			AssetId: poKadrze.Asset.Id,
			Filter:  shared.DesignPhotoFilter(shared.DesignPhotoFilterMonochrome),
			Amount:  wskaznik(1.0),
		}, &poFiltrze)

	var lancuch shared.DesignPhotoHistoryGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoHistoryGet,
		shared.DesignPhotoHistoryGetRequest{AssetId: poFiltrze.Asset.Id}, &lancuch)

	if len(lancuch.Edits) != 2 {
		t.Fatalf("łańcuch ma %d ogniw, a wykonano dwie czynności (kadr, filtr)", len(lancuch.Edits))
	}
	// Ogniwa idą od NAJSTARSZEGO — tak opisuje pole kontrakt i tak czyta się
	// historię pracy.
	if lancuch.Edits[0].Command != shared.CommandDesignPhotoCrop {
		t.Errorf("pierwsze ogniwo to %q, a pierwszą czynnością był kadr", lancuch.Edits[0].Command)
	}
	if lancuch.Edits[1].Command != shared.CommandDesignPhotoFilterApply {
		t.Errorf("drugie ogniwo to %q, a drugą czynnością był filtr", lancuch.Edits[1].Command)
	}
	if len(lancuch.Edits[0].Settings) == 0 {
		t.Error("ogniwo kadru nie niesie nastaw — bez nich nie da się powtórzyć tej czynności")
	}
	if lancuch.Edits[0].SourceAssetId == nil || *lancuch.Edits[0].SourceAssetId != zrodlo.Id {
		t.Errorf("pierwsze ogniwo wskazuje źródło %v, a obróbka szła z %s",
			lancuch.Edits[0].SourceAssetId, zrodlo.Id)
	}
}

// TestMetadaneZasobuMierzaPlikNieWiersz pilnuje, żeby `design.photo.metadata.get`
// czytał PLIK, a nie odsyłał tego, co wpisano do wiersza przy wniesieniu.
func TestMetadaneZasobuMierzaPlikNieWiersz(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	bajty := obrazPNG(t, 37, 19)
	zasob := wniesObrazSprawdzianu(t, zmontowany, zycie, "okno-fotografii", bajty)

	var metadane shared.DesignPhotoMetadataGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignPhotoMetadataGet,
		shared.DesignPhotoMetadataGetRequest{AssetId: zasob.Id}, &metadane)

	if metadane.Metadata.Width == nil || *metadane.Metadata.Width != 37 {
		t.Errorf("zmierzona szerokość to %v, a plik ma 37 px", metadane.Metadata.Width)
	}
	if metadane.Metadata.Height == nil || *metadane.Metadata.Height != 19 {
		t.Errorf("zmierzona wysokość to %v, a plik ma 19 px", metadane.Metadata.Height)
	}
	if metadane.Metadata.SizeBytes == nil || *metadane.Metadata.SizeBytes != len(bajty) {
		t.Errorf("zmierzona wielkość to %v bajtów, a plik ma %d", metadane.Metadata.SizeBytes,
			len(bajty))
	}
	if metadane.Metadata.Format == nil || *metadane.Metadata.Format != "png" {
		t.Errorf("zmierzony format to %v, a plik jest png", metadane.Metadata.Format)
	}
	// Liczba odczytanych pól EXIF wchodzi ZAWSZE, także zerowa: zero znaczy „plik
	// EXIF-u nie ma" i jest odpowiedzią, nie brakiem odpowiedzi.
	if metadane.Metadata.ExifFieldsRead == nil {
		t.Error("odpowiedź nie podaje liczby odczytanych pól EXIF — zero jest tu odpowiedzią")
	}
}
