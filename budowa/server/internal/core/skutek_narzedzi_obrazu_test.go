// Skutek czterech czynności obrazu modelu (`image.inspect`, `image.transform`,
// `image.adjust`, `image.convert`) mierzy piksele wyniku wkompilowanym
// rachunkiem, nie pola odpowiedzi ani obecność programu.
package core

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"

	"golang.org/x/image/webp"

	"danacoconsole/shared"
)

// obrazPolowaNaPolowe składa PNG podzielony pionowo na dwie barwy. Materiał do
// mierzenia kadru i obrotu: każda z dwóch połówek jest rozpoznawalna po barwie,
// więc widać nie tylko wymiary wyniku, ale i to, CO w nim zostało.
func obrazPolowaNaPolowe(t *testing.T, szerokosc, wysokosc int,
	lewa, prawa color.NRGBA) []byte {
	t.Helper()

	plotno := image.NewNRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			if x < szerokosc/2 {
				plotno.Set(x, y, lewa)
			} else {
				plotno.Set(x, y, prawa)
			}
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	return bufor.Bytes()
}

// barwaPunktu oddaje trzy składowe punktu w zakresie 0..255 — postać, w której
// da się je porównać z barwą wniesioną.
func barwaPunktu(obraz image.Image, x, y int) (int, int, int) {
	r, g, b, _ := obraz.At(x, y).RGBA()
	return int(r >> 8), int(g >> 8), int(b >> 8)
}

// TestSkalowanieZachowujeProporcjeIDajeZadanySzerokosc mierzy skutek
// `image.transform` wymiarami wyniku: przy podanej samej szerokości wysokość ma
// wyjść z proporcji, a nie zostać nietknięta.
func TestSkalowanieZachowujeProporcjeIDajeZadanySzerokosc(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	bialy := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	zrodlo := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 40, 20, bialy), "źródło")

	var wynik shared.ImageTransformResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageTransform,
		shared.ImageTransformRequest{
			WindowId:  wskaznik("okno-sprawdzianu"),
			AssetId:   wskaznik(zrodlo),
			Operation: shared.ImageTransformKindResize,
			Width:     wskaznik(20),
		}, &wynik)

	obraz := obrazZMagazynu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	if obraz.Bounds().Dx() != 20 || obraz.Bounds().Dy() != 10 {
		t.Fatalf("skalowanie do szerokości 20 dało %dx%d, a proporcja 2:1 wymaga 20x10",
			obraz.Bounds().Dx(), obraz.Bounds().Dy())
	}
}

// TestKadrowanieWycinaWskazanaPolowe mierzy kadr barwą: z obrazu pół czerwonego,
// pół białego kadr lewej połowy ma być czerwony w całości.
func TestKadrowanieWycinaWskazanaPolowe(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	czerwony := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	bialy := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	zrodlo := wniesObraz(t, zmontowany, zycie,
		obrazPolowaNaPolowe(t, 40, 20, czerwony, bialy), "połowa na połowę")

	var wynik shared.ImageTransformResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageTransform,
		shared.ImageTransformRequest{
			WindowId:  wskaznik("okno-sprawdzianu"),
			AssetId:   wskaznik(zrodlo),
			Operation: shared.ImageTransformKindCrop,
			Width:     wskaznik(20),
			Height:    wskaznik(20),
		}, &wynik)

	obraz := obrazZMagazynu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	if obraz.Bounds().Dx() != 20 || obraz.Bounds().Dy() != 20 {
		t.Fatalf("kadr 20x20 dał %dx%d", obraz.Bounds().Dx(), obraz.Bounds().Dy())
	}
	granice := obraz.Bounds()
	for _, punkt := range [][2]int{{0, 0}, {19, 19}, {10, 10}} {
		r, g, b := barwaPunktu(obraz, granice.Min.X+punkt[0], granice.Min.Y+punkt[1])
		if r < 200 || g > 60 || b > 60 {
			t.Fatalf("punkt (%d,%d) kadru ma barwę (%d,%d,%d) — kadr wyciął nie tę połowę",
				punkt[0], punkt[1], r, g, b)
		}
	}
}

// TestObrotIdzieZgodnieZeWskazowkamiZegara pilnuje kierunku obrotu: mierzy, gdzie
// po obrocie wylądowała czerwona połowa, bo pomiar samych boków by go nie zauważył.
func TestObrotIdzieZgodnieZeWskazowkamiZegara(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	czerwony := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	bialy := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	zrodlo := wniesObraz(t, zmontowany, zycie,
		obrazPolowaNaPolowe(t, 40, 20, czerwony, bialy), "połowa na połowę")

	var wynik shared.ImageTransformResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageTransform,
		shared.ImageTransformRequest{
			WindowId:  wskaznik("okno-sprawdzianu"),
			AssetId:   wskaznik(zrodlo),
			Operation: shared.ImageTransformKindRotate,
			Degrees:   wskaznik(90),
		}, &wynik)

	obraz := obrazZMagazynu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	if obraz.Bounds().Dx() != 20 || obraz.Bounds().Dy() != 40 {
		t.Fatalf("obrót o 90 stopni dał %dx%d, a z 40x20 ma wyjść 20x40",
			obraz.Bounds().Dx(), obraz.Bounds().Dy())
	}
	granice := obraz.Bounds()
	r, g, b := barwaPunktu(obraz, granice.Min.X+10, granice.Min.Y+5)
	if r < 200 || g > 60 || b > 60 {
		t.Fatalf("górna część obrotu ma barwę (%d,%d,%d), a ma być czerwona — "+
			"obrót poszedł w drugą stronę", r, g, b)
	}
	r, g, b = barwaPunktu(obraz, granice.Min.X+10, granice.Max.Y-5)
	if r < 200 || g < 200 || b < 200 {
		t.Fatalf("dolna część obrotu ma barwę (%d,%d,%d), a ma być biała", r, g, b)
	}
}

// TestSkalaSzarosciZrownujeSkladowe mierzy retusz `image.adjust`: czerwień po
// zamianie na skalę szarości ma mieć trzy składowe równe.
func TestSkalaSzarosciZrownujeSkladowe(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	czerwony := color.NRGBA{R: 220, G: 30, B: 30, A: 255}
	zrodlo := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 16, 16, czerwony), "czerwień")

	var wynik shared.ImageAdjustResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageAdjust,
		shared.ImageAdjustRequest{
			WindowId:  wskaznik("okno-sprawdzianu"),
			AssetId:   wskaznik(zrodlo),
			Operation: shared.ImageAdjustKindGrayscale,
		}, &wynik)

	obraz := obrazZMagazynu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	r, g, b := barwaPunktu(obraz, obraz.Bounds().Min.X+8, obraz.Bounds().Min.Y+8)
	if r != g || g != b {
		t.Fatalf("po skali szarości punkt ma barwę (%d,%d,%d) — składowe się różnią", r, g, b)
	}
	if r == 0 || r == 255 {
		t.Fatalf("po skali szarości punkt ma jasność %d — czerwień o jasności skrajnej "+
			"znaczy, że rachunek nie policzył luminancji, a wygasił obraz", r)
	}
}

// TestOdszumienieUsuwaPunktOdstajacy mierzy odszumienie tam, gdzie widać różnicę
// między medianą a rozmyciem: pojedynczy czarny punkt na białym tle ma zniknąć
// BEZ ROZMAZANIA po sąsiedztwie. Rozmycie zostawiłoby w tym miejscu szarość.
func TestOdszumienieUsuwaPunktOdstajacy(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	bialy := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	plotno := image.NewNRGBA(image.Rect(0, 0, 16, 16))
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			plotno.Set(x, y, bialy)
		}
	}
	plotno.Set(8, 8, color.NRGBA{A: 255})
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	zrodlo := wniesObraz(t, zmontowany, zycie, bufor.Bytes(), "szum")

	var wynik shared.ImageAdjustResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageAdjust,
		shared.ImageAdjustRequest{
			WindowId:  wskaznik("okno-sprawdzianu"),
			AssetId:   wskaznik(zrodlo),
			Operation: shared.ImageAdjustKindDenoise,
		}, &wynik)

	obraz := obrazZMagazynu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	granice := obraz.Bounds()
	r, g, b := barwaPunktu(obraz, granice.Min.X+8, granice.Min.Y+8)
	if r < 250 || g < 250 || b < 250 {
		t.Fatalf("punkt odstający po odszumieniu ma barwę (%d,%d,%d) — mediana miała "+
			"go zastąpić barwą sąsiedztwa, czyli bielą", r, g, b)
	}
	r, g, b = barwaPunktu(obraz, granice.Min.X+9, granice.Min.Y+9)
	if r < 250 || g < 250 || b < 250 {
		t.Fatalf("sąsiad punktu odstającego ma barwę (%d,%d,%d) — szum został rozmazany "+
			"po sąsiedztwie, a nie usunięty", r, g, b)
	}
}

// TestKonwersjaDoJpegDajeCzytelnyJpeg mierzy `image.convert`: wynik ma być
// plikiem, który dekoduje się jako JPEG, o wymiarach źródła.
func TestKonwersjaDoJpegDajeCzytelnyJpeg(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	bialy := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	zrodlo := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 24, 12, bialy), "źródło")

	var wynik shared.ImageConvertResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageConvert,
		shared.ImageConvertRequest{
			WindowId: wskaznik("okno-sprawdzianu"),
			AssetId:  wskaznik(zrodlo),
			Format:   "jpeg",
			Quality:  wskaznik(80),
		}, &wynik)

	if wynik.SizeBytes == 0 {
		t.Fatal("konwersja oddała rozmiar zero — plik bez bajtów nie jest wynikiem")
	}
	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	nastawy, format, err := rozpoznajFormatPlikuSprawdzianu(sciezka)
	if err != nil {
		t.Fatalf("za wynikiem konwersji nie leży czytelny obraz: %v", err)
	}
	if format != "jpeg" {
		t.Fatalf("konwersja do jpeg dała plik formatu %q", format)
	}
	if nastawy.Width != 24 || nastawy.Height != 12 {
		t.Fatalf("konwersja zmieniła wymiary na %dx%d", nastawy.Width, nastawy.Height)
	}
}

// TestKonwersjaDoWebpBezstratnegoNieWymagaProgramu mierzy zapis WEBP — jedyny
// format, który produkt do tej pory umiał wyłącznie zdekodować, i ma przejść
// zawsze, bez udziału programu.
func TestKonwersjaDoWebpBezstratnegoNieWymagaProgramu(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	czerwony := color.NRGBA{R: 200, G: 40, B: 40, A: 255}
	zrodlo := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 18, 9, czerwony), "źródło")

	var wynik shared.ImageConvertResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageConvert,
		shared.ImageConvertRequest{
			WindowId: wskaznik("okno-sprawdzianu"),
			AssetId:  wskaznik(zrodlo),
			Format:   "webp",
			Lossless: wskaznik(true),
		}, &wynik)

	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	plik, err := os.Open(sciezka)
	if err != nil {
		t.Fatalf("nie można otworzyć wyniku: %v", err)
	}
	defer plik.Close()

	obraz, err := webp.Decode(plik)
	if err != nil {
		t.Fatalf("wynik nie jest czytelnym plikiem WEBP: %v", err)
	}
	if obraz.Bounds().Dx() != 18 || obraz.Bounds().Dy() != 9 {
		t.Fatalf("zapis WEBP dał %dx%d, a źródło ma 18x9",
			obraz.Bounds().Dx(), obraz.Bounds().Dy())
	}
}

// TestPomiarObrazuCzytaFormatIWymiaryBezProgramu mierzy `image.inspect`: format
// i wymiary obrazu mają wyjść bez udziału programu pakietu serwera.
func TestPomiarObrazuCzytaFormatIWymiaryBezProgramu(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	bialy := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	zrodlo := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 33, 17, bialy), "źródło")

	var wynik shared.ImageInspectResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageInspect,
		shared.ImageInspectRequest{AssetId: wskaznik(zrodlo)}, &wynik)

	if wynik.Format != "png" {
		t.Fatalf("pomiar nazwał format %q, a plik jest PNG", wynik.Format)
	}
	if wynik.Width != 33 || wynik.Height != 17 {
		t.Fatalf("pomiar oddał wymiary %dx%d, a plik ma 33x17", wynik.Width, wynik.Height)
	}
	if wynik.SizeBytes == 0 {
		t.Fatal("pomiar oddał rozmiar zero — plik z magazynu ma bajty")
	}
	if wynik.ColorSpace == nil || *wynik.ColorSpace == "" {
		t.Fatal("pomiar nie nazwał przestrzeni barw, choć dekoder zna model koloru PNG")
	}
}

// TestPomiarNiepodanegoObrazuOdmawiaNazywajacBrak pilnuje, że czynność czytająca
// nie zgaduje przedmiotu pracy: żądanie bez `assetId` i bez `sourcePath` ma
// dostać odmowę nazywającą brak, a nie pomiar ostatniego zasobu.
func TestPomiarNiepodanegoObrazuOdmawiaNazywajacBrak(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandImageInspect,
		shared.ImageInspectRequest{})
	if !strings.Contains(odmowa.Message, "assetId") {
		t.Fatalf("odmowa nie nazywa brakującego wskazania: %s", odmowa.Message)
	}
}

// TestKonwersjaDoAvifBezProgramuOdmawiaNazywajacBrak mierzy kształt odmowy na
// jedynej drodze rodziny bez rachunku wkompilowanego: AVIF ma odmówić zdaniem
// nazywającym program, nie błędem wewnętrznym.
func TestKonwersjaDoAvifBezProgramuOdmawiaNazywajacBrak(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	bialy := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	zrodlo := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 16, 16, bialy), "źródło")

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandImageConvert,
		shared.ImageConvertRequest{
			WindowId: wskaznik("okno-sprawdzianu"),
			AssetId:  wskaznik(zrodlo),
			Format:   "avif",
		})
	// Trzy człony odmowy: co, program po nazwie.

	// Skąd wiadomo: nie ma na tej maszynie; jak naprawić: pakiet do zainstalowania.
	odmowaNazywa(t, odmowa, shared.ErrorCodeChannelUnavailable,
		"ImageMagick", "nie ma na tej maszynie", "naprawa: zainstalować pakiet imagemagick")
	t.Logf("odmowa: %s", odmowa.Message)
}

// rozpoznajFormatPlikuSprawdzianu czyta nagłówek pliku wynikowego i oddaje
// nazwę formatu, jaką rozpoznaje biblioteka standardowa obrazu.
func rozpoznajFormatPlikuSprawdzianu(sciezka string) (image.Config, string, error) {
	plik, err := os.Open(sciezka)
	if err != nil {
		return image.Config{}, "", err
	}
	defer plik.Close()
	return image.DecodeConfig(plik)
}
