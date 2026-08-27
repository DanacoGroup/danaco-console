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

// Skutek dogniecenia zapisu przy `image.convert`
// (`adapter_narzedzia_obraz_kompresja.go`).
//
// ── Dlaczego te sprawdziany pomijają się przy braku programu ────────────────
// Dogniecenie jest ULEPSZENIEM, nie warunkiem: przy braku programu konwersja ma
// oddać ten sam obraz zapisany dłuższym strumieniem. Sprawdzian mierzący zysk na
// rozmiarze mierzy więc obecność programu tak samo jak jego pracę — i na maszynie
// bez niego musi się pominąć z NAZWANYM powodem, zamiast zawieść albo, gorzej,
// przejść na zielono nie zmierzywszy niczego.
//
// Osobny sprawdzian niżej idzie ZAWSZE i mierzy rzecz odwrotną: że konwersja
// udaje się niezależnie od tego, czy program stoi. To jest ta połowa reguły,
// która ma być prawdziwa na instalce Operatora.

// obrazDoDogniecenia składa PNG, który da się skrócić: gradient o łagodnym
// przebiegu ma silne predykcje międzywierszowe, a koder wkompilowany bierze
// jeden filtr i jeden przebieg deflate, więc zostawia po sobie zapas.
//
// Obraz jednolity nie nadałby się do pomiaru: koder Go zapisuje go już blisko
// granicy i zysk bywa zerowy, co czytałoby się jak brak dogniecenia.
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

// TestKonwersjaDoPngDogniataZapisProgramem mierzy najtwardszy skutek tej pracy:
// plik leżący za odwołaniem jest KRÓTSZY niż ten sam obraz zapisany samym
// koderem wkompilowanym — i nadal jest tym samym obrazem.
//
// Porównanie idzie z zapisem kodera Go policzonym tu na miejscu, a nie ze stałą
// liczbą bajtów: stała rozjechałaby się przy pierwszej zmianie biblioteki
// i zaczęłaby mierzyć jej wydanie zamiast pracy programu.
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

	// Krótszy zapis, który przestał być tym obrazem, nie jest ulepszeniem.
	// Dogniatanie bezstratne ma zostawić KAŻDY punkt taki, jaki był.
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

	// Wynik ma zostać czytelnym WEBP-em o wymiarach źródła — dogniecenie, po
	// którym pliku nie da się odczytać, byłoby stratą, nie zyskiem.
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

// TestKonwersjaStratnaSprowadzaPngDoPalety mierzy pngquant — jedyny z czterech
// programów, który PIKSELE ZMIENIA, i dlatego jedyny wołany wyłącznie na
// wyraźne żądanie zapisu stratnego.
//
// Miarą jest model barw wyniku, a nie sama liczba bajtów: plik palety niesie
// `color.Palette`, a zapis pełnobarwny — `NRGBA`. Model palety dowodzi, że wynik
// programu został wzięty, bo koder wkompilowany palety nie zapisuje.
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

// obrazRozsypanejPalety składa PNG z dwustu barw rozrzuconych bez ładu.
//
// Taki obraz jest materiałem WŁAŚCIWYM do mierzenia palety: zapis pełnobarwny
// nie ma czego przewidzieć i płaci trzy bajty za punkt, a paleta mieści
// wszystkie barwy co do jednej. Gradient nadałby się gorzej — pngquant odmawia
// sprowadzenia płynnego przejścia do palety, bo nie zmieściłby się w progu
// jakości, i wtedy sprawdzian mierzyłby odmowę zamiast pracy.
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

// TestKonwersjaUdajeSieNiezaleznieOdProgramuDogniatajacego jest drugą połową
// reguły i dlatego NIE POMIJA SIĘ nigdy: konwersja ma dać czytelny obraz także
// wtedy, gdy żaden program dogniatający nie stoi.
//
// Sprawdzian mierzy to, co widzi Operator na cienkiej instalce — i wytwarza ją
// tu na miejscu: pusta ścieżka wyszukiwania czyni z tej maszyny maszynę bez
// programów, więc zdanie z nazwy sprawdzianu jest mierzone wszędzie, a nie
// tylko tam, gdzie programów akurat nie doinstalowano. Gdyby dogniatanie
// kiedykolwiek zaczęło odmawiać przy braku programu, ten sprawdzian zawiedzie.
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
