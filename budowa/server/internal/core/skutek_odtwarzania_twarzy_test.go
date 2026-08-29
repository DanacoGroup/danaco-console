package core

import (
	"bytes"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// materialTwarzySprawdzianu to fotografia domeny publicznej z jedną twarzą wyśrodkowaną w kadrze.
const materialTwarzySprawdzianu = "testdata/twarz_sprawdzianu.jpg"

/*
Wycinek twarzy na materiale sprawdzianu, ułamkiem szerokości i wysokości obrazu:
twarz na tej fotografii stoi wyśrodkowana, poza wycinkiem leży flaga i tło,
których przebieg twarzowy nie ma prawa dotknąć.

Granice pionowe biorą się z pomiaru, nie z oka. Przebieg twarzowy odtwarza twarz
w całości — od nakrycia głowy po żuchwę i kołnierz — a zmiany na materiale
sprawdzianu mieszczą się w prostokącie (395,110)–(810,890) obrazu 1024×1200.
Wcześniejsza granica dolna 0.58 (696 px) przecinała twarz pod nosem, więc usta,
broda i żuchwa liczyły się jako obszar poza twarzą.

Granice poziome zostają bez zmian: pomiar pokazał, że zmiany nie wychodzą poza
kolumnę twarzy, więc rozlanie na flagę dalej ten sprawdzian złapie.
*/
const (
	wycinekTwarzyLewy  = 0.22
	wycinekTwarzyGorny = 0.08
	wycinekTwarzyPrawy = 0.80
	wycinekTwarzyDolny = 0.76
)

// progZmienionychPozaWycinkiem to liczba pikseli poza wycinkiem twarzy, poniżej
// której różnicę tłumaczy się szumem potoku (piksele graniczne pociągnięte przez
// maskę wtopienia), nie zmianą rozlaną po obrazie; zmiana globalna dotyka setek
// tysięcy pikseli poza wycinkiem, więc próg zostaje daleko w tyle za taką skalą.
const progZmienionychPozaWycinkiem = 2000

// Sprawdzian odmowy przebiegu twarzowego image.upscale przy braku wag, mierzonej bez ciężkich zasobów.

// TestBrakWagTwarzyOdmawiaNazywajacPlikIDrogeNaprawy pilnuje, żeby odmowa niosła nazwę brakującego pliku, ścieżkę, pod którą ma leżeć, oraz miejsce, z którego się go bierze.
func TestBrakWagTwarzyOdmawiaNazywajacPlikIDrogeNaprawy(t *testing.T) {
	pusty := t.TempDir()

	err := sprawdzWagiTwarzy(pusty)
	if err == nil {
		t.Fatal("katalog bez wag przeszedł sprawdzenie — przebieg ruszyłby i wywrócił się " +
			"dopiero we wnętrzu pomocnika, komunikatem Pythona zamiast odmową rdzenia")
	}
	blad := protocol.BladZeZrodla(shared.ErrorCodeInternalError, err)
	if blad.Code != shared.ErrorCodeChannelUnavailable {
		t.Fatalf("brak wag dostał kod %s, a brakująca instalacja jest zapleczem "+
			"niedostępnym (%s)", blad.Code, shared.ErrorCodeChannelUnavailable)
	}
	for _, oczekiwany := range []string{wagiOdtwarzaniaTwarzy, pusty, "github.com"} {
		if !strings.Contains(blad.Message, oczekiwany) {
			t.Fatalf("odmowa nie wymienia %q; treść: %s", oczekiwany, blad.Message)
		}
	}
}

// TestPomocnikTwarzyStoiWWykazieZaleznosci pilnuje, żeby silnik przebiegu twarzowego był widoczny w sondzie startowej, zanim jego brak ujawni się dopiero po wywołaniu.
func TestPomocnikTwarzyStoiWWykazieZaleznosci(t *testing.T) {
	szukany := narzedzieOdtwarzaniaTwarzy().Program
	if szukany == "" {
		t.Fatal("silnik przebiegu twarzowego nie ma nazwy programu")
	}
	for _, pozycja := range ZaleznosciZewnetrzne() {
		if pozycja.Narzedzie.Program == szukany {
			if strings.TrimSpace(pozycja.Zakres) == "" {
				t.Fatal("pozycja wykazu nie mówi, co przestaje działać przy braku")
			}
			return
		}
	}
	t.Fatalf("programu %s nie ma w wykazie zależności zewnętrznych", szukany)
}

// wymagajSilnikaTwarzowego zatrzymuje bieg odmową, nie pominięciem: katalogWagTwarzy
// wskazuje stały katalog wdrożeniowy, a nie zmienną środowiska, więc maszyna drabiny
// odbioru ma na nim nosić wagi — ich brak jest awarią środowiska, którą sprawdzian
// zgłasza głośno, zamiast świecić zielono bez zmierzenia kryterium ani razu.
func wymagajSilnikaTwarzowego(t *testing.T) {
	t.Helper()

	if err := sprawdzWagiTwarzy(katalogWagTwarzy()); err != nil {
		t.Fatalf("wagi przebiegu twarzowego nie stoją na tej maszynie: %v", err)
	}
	if _, err := os.Stat(filepath.Join(katalogModeliPowiekszenia(), modelPowiekszeniaZdjec+".bin")); err != nil {
		t.Fatalf("wagi Real-ESRGAN nie stoją na tej maszynie: %v", err)
	}
}

// wycinekTwarzy oddaje prostokąt wycinka twarzy wewnątrz `granice`, niezależnie od
// krotności powiększenia, którym `granice` zostały pomnożone względem materiału źródłowego.
func wycinekTwarzy(granice image.Rectangle) image.Rectangle {
	szerokosc := float64(granice.Dx())
	wysokosc := float64(granice.Dy())
	return image.Rect(
		granice.Min.X+int(szerokosc*wycinekTwarzyLewy),
		granice.Min.Y+int(wysokosc*wycinekTwarzyGorny),
		granice.Min.X+int(szerokosc*wycinekTwarzyPrawy),
		granice.Min.Y+int(wysokosc*wycinekTwarzyDolny),
	)
}

// roznicaWObszarze mierzy największą różnicę składowej koloru i liczbę zauważalnie
// zmienionych pikseli między dwoma obrazami tych samych wymiarów, licząc osobno
// wewnątrz podanego wycinka i poza nim.
func roznicaWObszarze(bez, z image.Image, wycinek image.Rectangle) (
	wewnatrzNajwieksza, wewnatrzZmienionych, pozaNajwieksza, pozaZmienionych int) {

	granice := bez.Bounds()
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			r1, g1, b1, _ := bez.At(x, y).RGBA()
			r2, g2, b2, _ := z.At(x, y).RGBA()
			roznica := abs(int(r1>>8) - int(r2>>8))
			if d := abs(int(g1>>8) - int(g2>>8)); d > roznica {
				roznica = d
			}
			if d := abs(int(b1>>8) - int(b2>>8)); d > roznica {
				roznica = d
			}
			if (image.Point{X: x, Y: y}).In(wycinek) {
				if roznica > wewnatrzNajwieksza {
					wewnatrzNajwieksza = roznica
				}
				if roznica > 15 {
					wewnatrzZmienionych++
				}
			} else {
				if roznica > pozaNajwieksza {
					pozaNajwieksza = roznica
				}
				if roznica > 15 {
					pozaZmienionych++
				}
			}
		}
	}
	return wewnatrzNajwieksza, wewnatrzZmienionych, pozaNajwieksza, pozaZmienionych
}

// TestPrzebiegTwarzowyZmieniaPikseleTwarzyNaPrawdziwymZdjeciu porównuje piksele wyniku
// bez i z przebiegiem twarzowym: rozstrzyga wyłącznie zmierzona różnica w wycinku
// twarzy, nigdy samoopis pomocnika, a różnica poza wycinkiem ma zostać znikoma.
func TestPrzebiegTwarzowyZmieniaPikseleTwarzyNaPrawdziwymZdjeciu(t *testing.T) {
	wymagajSilnikaTwarzowego(t)

	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	bajty, err := os.ReadFile(materialTwarzySprawdzianu)
	if err != nil {
		t.Fatalf("nie można odczytać materiału sprawdzianu %s: %v", materialTwarzySprawdzianu, err)
	}
	zrodlo := wniesObraz(t, zmontowany, zycie, bajty, "twarz.jpg")

	var bezTwarzy shared.ImageUpscaleResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageUpscale,
		shared.ImageUpscaleRequest{
			WindowId: wskaznik("okno-sprawdzianu"),
			AssetId:  wskaznik(zrodlo),
			Scale:    wskaznik(2),
		}, &bezTwarzy)

	var zTwarza shared.ImageUpscaleResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageUpscale,
		shared.ImageUpscaleRequest{
			WindowId: wskaznik("okno-sprawdzianu"),
			AssetId:  wskaznik(zrodlo),
			Scale:    wskaznik(2),
			Faces:    wskaznik(true),
		}, &zTwarza)

	obrazBez := obrazZMagazynu(t, zmontowany, zycie, katalog, bezTwarzy.Asset.Id)
	obrazZ := obrazZMagazynu(t, zmontowany, zycie, katalog, zTwarza.Asset.Id)
	if obrazBez.Bounds() != obrazZ.Bounds() {
		t.Fatalf("wynik z przebiegiem twarzowym ma wymiary %v, a bez niego %v — powiększenie miało dać "+
			"te same wymiary, przebieg twarzowy tylko poprawia treść wycinka twarzy",
			obrazZ.Bounds(), obrazBez.Bounds())
	}

	wycinek := wycinekTwarzy(obrazBez.Bounds())
	wNajwieksza, wZmienionych, pNajwieksza, pZmienionych := roznicaWObszarze(obrazBez, obrazZ, wycinek)

	if wNajwieksza < 40 {
		t.Fatalf("największa różnica składowej koloru w wycinku twarzy to %d — GFPGAN "+
			"praktycznie nie zmienił wycinka", wNajwieksza)
	}
	if wZmienionych < 1000 {
		t.Fatalf("tylko %d pikseli wycinka twarzy zmieniło się zauważalnie — przebieg twarzowy "+
			"miał przemalować cały wycinek, nie garstkę pikseli", wZmienionych)
	}
	if pZmienionych > progZmienionychPozaWycinkiem {
		t.Fatalf("poza wycinkiem twarzy zmieniło się zauważalnie %d pikseli (największa różnica %d) "+
			"— przebieg twarzowy miał dotknąć wyłącznie wycinka twarzy, nie całego obrazu",
			pZmienionych, pNajwieksza)
	}
}

// TestRoznicaWObszarzeOdrzucaZmianeCalegoObrazu jest przeciwsprawdzianem miary, bez
// zależności od wag: materiał przyciemniony jednolicie na całej powierzchni ma dawać
// dużą różnicę też POZA wycinkiem twarzy — miara liczona na całym obrazie przepuściłaby
// taką zmianę, miara w wycinku ją odrzuca.
func TestRoznicaWObszarzeOdrzucaZmianeCalegoObrazu(t *testing.T) {
	bajty, err := os.ReadFile(materialTwarzySprawdzianu)
	if err != nil {
		t.Fatalf("nie można odczytać materiału sprawdzianu %s: %v", materialTwarzySprawdzianu, err)
	}
	oryginal, _, err := image.Decode(bytes.NewReader(bajty))
	if err != nil {
		t.Fatalf("materiał sprawdzianu nie jest obrazem: %v", err)
	}

	granice := oryginal.Bounds()
	przyciemniony := image.NewRGBA(granice)
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			r, g, b, a := oryginal.At(x, y).RGBA()
			przyciemniony.Set(x, y, color.RGBA{
				R: przyciemnijSkladowa(r), G: przyciemnijSkladowa(g), B: przyciemnijSkladowa(b),
				A: uint8(a >> 8),
			})
		}
	}

	_, _, pNajwieksza, pZmienionych := roznicaWObszarze(oryginal, przyciemniony, wycinekTwarzy(granice))
	if pZmienionych < 1000 {
		t.Fatalf("przyciemnienie całego obrazu zmieniło poza wycinkiem twarzy tylko %d pikseli "+
			"(największa różnica %d) — przeciwsprawdzian nie odróżnia zmiany globalnej od zmiany "+
			"w wycinku", pZmienionych, pNajwieksza)
	}
}

// przyciemnijSkladowa odejmuje stałą wartość od jednej składowej koloru w skali
// 0..255, przycinając do zera, żeby przeciwsprawdzian miał materiał zmieniony
// jednolicie na całej powierzchni, bez sięgania po żaden silnik.
func przyciemnijSkladowa(wartosc uint32) uint8 {
	w := int(wartosc>>8) - 60
	if w < 0 {
		w = 0
	}
	return uint8(w)
}
