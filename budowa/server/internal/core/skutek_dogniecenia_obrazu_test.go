package core

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"testing"

	"golang.org/x/image/webp"

	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// Sprawdziany dogniecenia pomijają się bez programu zewnętrznego; sprawdzian niżej idzie zawsze.

// obrazDoDogniecenia składa PNG z gradientem o łagodnym przebiegu, który koder wkompilowany zapisuje z zapasem miejsca do skrócenia przez dogniatanie.
func obrazDoDogniecenia(t *testing.T, szerokosc, wysokosc int) []byte {
	t.Helper()

	plotno := image.NewNRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			plotno.Set(x, y, color.NRGBA{
				R: uint8((x*7 + y*3) % 256),
				G: uint8((y * 5) % 256),
				B: uint8((x * 11) % 256),
				A: 255,
			})
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	return bufor.Bytes()
}

// pominBezProgramu pomija sprawdzian, gdy program dogniatający nie stoi —
// nazywając, którego brakuje. Pominięcie widoczne w wyniku biegu jest tu
// informacją: mówi, że maszyna nie ma czym zmierzyć tej połowy reguły.
func pominBezProgramu(t *testing.T, narzedzie zewnetrzne.Narzedzie) {
	t.Helper()

	if !zewnetrzne.Stoi(narzedzie) {
		t.Skipf("na tej maszynie nie stoi %s (%s) — dogniecenie zapisu jest ulepszeniem, "+
			"więc bez programu nie ma czego mierzyć; konwersji pilnuje "+
			"TestKonwersjaUdajeSieNiezaleznieOdProgramuDogniatajacego",
			narzedzie.Nazwa, narzedzie.Program)
	}
}

// TestKonwersjaDoPngDogniataZapisProgramem sprawdza, że plik zapisany po dognieceniu jest krótszy niż ten sam obraz zapisany samym koderem wkompilowanym i pozostaje tym samym obrazem.
func TestKonwersjaDoPngDogniataZapisProgramem(t *testing.T) {
	pominBezProgramu(t, narzedzieOptipng())

	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	zrodloBajty := obrazDoDogniecenia(t, 400, 300)
	zrodlo := wniesObraz(t, zmontowany, zycie, zrodloBajty, "źródło")

	var wynik shared.ImageConvertResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageConvert,
		shared.ImageConvertRequest{
			WindowId: wskaznik("okno-sprawdzianu"),
			AssetId:  wskaznik(zrodlo),
			Format:   "png",
			Lossless: wskaznik(true),
		}, &wynik)

	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	zapisany, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać wyniku konwersji: %v", err)
	}

	// Miarą jest zapis kodera wkompilowanego dla TEGO SAMEGO obrazu.
	odczytany, err := png.Decode(bytes.NewReader(zrodloBajty))
	if err != nil {
		t.Fatalf("nie można odczytać obrazu sprawdzianu: %v", err)
	}
	var samKoder bytes.Buffer
	if err := png.Encode(&samKoder, odczytany); err != nil {
		t.Fatalf("nie można zapisać obrazu koderem wkompilowanym: %v", err)
	}

	if len(zapisany) >= samKoder.Len() {
		t.Fatalf("zapis po dogniataniu ma %d bajtów, a sam koder wkompilowany daje %d — "+
			"program nie skrócił zapisu", len(zapisany), samKoder.Len())
	}
	t.Logf("koder wkompilowany %d B, po %s %d B", samKoder.Len(),
		narzedzieOptipng().Program, len(zapisany))

	// Dogniatanie bezstratne ma zostawić każdy punkt obrazu bez zmian.
	poDognieceniu, err := png.Decode(bytes.NewReader(zapisany))
	if err != nil {
		t.Fatalf("plik po dogniataniu nie jest czytelnym obrazem: %v", err)
	}
	if poDognieceniu.Bounds() != odczytany.Bounds() {
		t.Fatalf("dogniatanie zmieniło wymiary z %v na %v",
			odczytany.Bounds(), poDognieceniu.Bounds())
	}
	for _, punkt := range []image.Point{{X: 0, Y: 0}, {X: 399, Y: 299}, {X: 137, Y: 201}} {
		czerwonyA, zielonyA, niebieskiA := barwaPunktu(odczytany, punkt.X, punkt.Y)
		czerwonyB, zielonyB, niebieskiB := barwaPunktu(poDognieceniu, punkt.X, punkt.Y)
		if czerwonyA != czerwonyB || zielonyA != zielonyB || niebieskiA != niebieskiB {
			t.Fatalf("dogniatanie zmieniło punkt %v z (%d,%d,%d) na (%d,%d,%d) — "+
				"zapis bezstratny nie ma prawa ruszyć pikseli", punkt,
				czerwonyA, zielonyA, niebieskiA, czerwonyB, zielonyB, niebieskiB)
		}
	}
}

// TestKonwersjaDoWebpDogniataZapisProgramem mierzy tę samą rzecz dla WEBP-a,
// którego koder wkompilowany (`nativewebp`) nie stroi predyktorów.
func TestKonwersjaDoWebpDogniataZapisProgramem(t *testing.T) {
	pominBezProgramu(t, narzedzieCwebp())

	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	zrodlo := wniesObraz(t, zmontowany, zycie, obrazDoDogniecenia(t, 400, 300), "źródło")

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

	// Wynik ma pozostać czytelnym plikiem WEBP o wymiarach źródła.
	obraz, err := webp.Decode(plik)
	if err != nil {
		t.Fatalf("wynik nie jest czytelnym plikiem WEBP: %v", err)
	}
	if obraz.Bounds().Dx() != 400 || obraz.Bounds().Dy() != 300 {
		t.Fatalf("zapis WEBP dał %dx%d, a źródło ma 400x300",
			obraz.Bounds().Dx(), obraz.Bounds().Dy())
	}
	if wynik.SizeBytes == 0 {
		t.Fatal("konwersja oddała rozmiar zero")
	}
	t.Logf("zapis WEBP po %s: %d B", narzedzieCwebp().Program, wynik.SizeBytes)
}

// TestKonwersjaDoJpegDogniataZapisProgramem mierzy jpegoptim: przeliczenie
// tablic Huffmana skraca plik, nie ruszając ani jednego współczynnika DCT.
func TestKonwersjaDoJpegDogniataZapisProgramem(t *testing.T) {
	pominBezProgramu(t, narzedzieJpegoptim())

	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	zrodloBajty := obrazDoDogniecenia(t, 400, 300)
	zrodlo := wniesObraz(t, zmontowany, zycie, zrodloBajty, "źródło")

	const jakosc = 90
	var wynik shared.ImageConvertResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageConvert,
		shared.ImageConvertRequest{
			WindowId: wskaznik("okno-sprawdzianu"),
			AssetId:  wskaznik(zrodlo),
			Format:   "jpeg",
			Quality:  wskaznik(jakosc),
		}, &wynik)

	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	zapisany, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać wyniku konwersji: %v", err)
	}

	// Miarą jest zapis kodera wkompilowanego przy TEJ SAMEJ jakości.
	odczytany, err := png.Decode(bytes.NewReader(zrodloBajty))
	if err != nil {
		t.Fatalf("nie można odczytać obrazu sprawdzianu: %v", err)
	}
	var samKoder bytes.Buffer
	if err := jpeg.Encode(&samKoder, odczytany, &jpeg.Options{Quality: jakosc}); err != nil {
		t.Fatalf("nie można zapisać obrazu koderem wkompilowanym: %v", err)
	}
	if len(zapisany) >= samKoder.Len() {
		t.Fatalf("zapis po dogniataniu ma %d bajtów, a sam koder wkompilowany daje %d — "+
			"program nie skrócił zapisu", len(zapisany), samKoder.Len())
	}
	if _, err := jpeg.Decode(bytes.NewReader(zapisany)); err != nil {
		t.Fatalf("plik po dogniataniu nie jest czytelnym JPEG-iem: %v", err)
	}
	t.Logf("koder wkompilowany %d B, po %s %d B", samKoder.Len(),
		narzedzieJpegoptim().Program, len(zapisany))
}

// TestKonwersjaStratnaSprowadzaPngDoPalety sprawdza, że zapis stratny niesie model barw palety, którego koder wkompilowany sam nie tworzy.
func TestKonwersjaStratnaSprowadzaPngDoPalety(t *testing.T) {
	pominBezProgramu(t, narzedziePngquant())

	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	zrodlo := wniesObraz(t, zmontowany, zycie, obrazRozsypanejPalety(t, 400, 300), "źródło")

	var stratny shared.ImageConvertResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageConvert,
		shared.ImageConvertRequest{
			WindowId: wskaznik("okno-sprawdzianu"),
			AssetId:  wskaznik(zrodlo),
			Format:   "png",
			Lossless: wskaznik(false),
		}, &stratny)

	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, stratny.Asset.Id)
	plik, err := os.Open(sciezka)
	if err != nil {
		t.Fatalf("nie można otworzyć wyniku: %v", err)
	}
	defer plik.Close()

	nastawy, err := png.DecodeConfig(plik)
	if err != nil {
		t.Fatalf("wynik nie jest czytelnym plikiem PNG: %v", err)
	}
	if nastawy.Width != 400 || nastawy.Height != 300 {
		t.Fatalf("zapis stratny dał %dx%d, a źródło ma 400x300", nastawy.Width, nastawy.Height)
	}
	if _, paleta := nastawy.ColorModel.(color.Palette); !paleta {
		t.Fatalf("wynik zapisu stratnego ma model barw %T, a nie paletę — "+
			"pngquant albo nie wszedł, albo jego wynik został odrzucony", nastawy.ColorModel)
	}
	t.Logf("zapis stratny PNG po %s: %d B, model barw paleta",
		narzedziePngquant().Program, stratny.SizeBytes)
}

// obrazRozsypanejPalety składa PNG z dwustu barw rozrzuconych bez ładu, materiał właściwy do mierzenia zapisu palety kolorów.
func obrazRozsypanejPalety(t *testing.T, szerokosc, wysokosc int) []byte {
	t.Helper()

	losowy := rand.New(rand.NewSource(7))
	paleta := make([]color.NRGBA, 200)
	for i := range paleta {
		paleta[i] = color.NRGBA{
			R: uint8(losowy.Intn(256)), G: uint8(losowy.Intn(256)),
			B: uint8(losowy.Intn(256)), A: 255,
		}
	}
	plotno := image.NewNRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			plotno.Set(x, y, paleta[losowy.Intn(len(paleta))])
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	return bufor.Bytes()
}

// TestKonwersjaUdajeSieNiezaleznieOdProgramuDogniatajacego sprawdza, że konwersja daje czytelny obraz nawet bez programu dogniatającego zewnętrznego.
func TestKonwersjaUdajeSieNiezaleznieOdProgramuDogniatajacego(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	zrodlo := wniesObraz(t, zmontowany, zycie, obrazDoDogniecenia(t, 40, 20), "źródło")

	var wynik shared.ImageConvertResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageConvert,
		shared.ImageConvertRequest{
			WindowId: wskaznik("okno-sprawdzianu"),
			AssetId:  wskaznik(zrodlo),
			Format:   "png",
			Lossless: wskaznik(true),
		}, &wynik)

	if wynik.SizeBytes == 0 {
		t.Fatal("konwersja oddała rozmiar zero — plik bez bajtów nie jest wynikiem")
	}
	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	nastawy, format, err := rozpoznajFormatPlikuSprawdzianu(sciezka)
	if err != nil {
		t.Fatalf("za wynikiem konwersji nie leży czytelny obraz: %v", err)
	}
	if format != "png" {
		t.Fatalf("konwersja do png dała plik formatu %q", format)
	}
	if nastawy.Width != 40 || nastawy.Height != 20 {
		t.Fatalf("konwersja zmieniła wymiary na %dx%d", nastawy.Width, nastawy.Height)
	}
}
