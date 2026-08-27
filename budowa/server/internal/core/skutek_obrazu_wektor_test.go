package core

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"

	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// Sprawdziany złożenia obrazu i wektoryzacji dekodują wynik z powrotem i mierzą go pikselem.

// obrazJednolity składa PNG wypełniony jedną barwą — materiał, w którym każda
// zmiana piksela jest zmianą widoczną i policzalną.
func obrazJednolity(t *testing.T, szerokosc, wysokosc int, barwa color.NRGBA) []byte {
	t.Helper()

	plotno := image.NewNRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			plotno.Set(x, y, barwa)
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć obrazu sprawdzianu: %v", err)
	}
	return bufor.Bytes()
}

// wniesObraz wnosi obraz do magazynu zasobów komendą design.asset.upload i oddaje
// identyfikator zasobu, gotowy do dalszego wskazania w komendach obrazu.
func wniesObraz(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	bajty []byte, nazwa string) string {
	t.Helper()

	var wynik shared.DesignAssetUploadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandDesignAssetUpload,
		shared.DesignAssetUploadRequest{
			WindowId:      "okno-sprawdzianu",
			Name:          wskaznik(nazwa),
			Kind:          shared.DesignAssetKindImage,
			Format:        wskaznik("png"),
			ContentBase64: wskaznik(wBase64(bajty)),
		}, &wynik)
	return wynik.Asset.Id
}

// obrazZMagazynu dekoduje wynik leżący pod odwołaniem zasobu, oddając obraz
// gotowy do pomiaru pikselem zamiast surowych bajtów pliku.
func obrazZMagazynu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	katalog, kod string) image.Image {
	t.Helper()

	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, kod)
	obraz, err := odczytajObrazPliku(sciezka)
	if err != nil {
		t.Fatalf("za wynikiem nie leży obraz, który da się zdekodować: %v", err)
	}
	return obraz
}

// TestZlozenieKladzieNakladkeWWskazanymMiejscu wykazuje skutek złożenia
// pikselem: poza nakładką zostaje barwa podstawy, a pod nakładką pojawia się
// barwa nakładki.
func TestZlozenieKladzieNakladkeWWskazanymMiejscu(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	bialy := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	czerwony := color.NRGBA{R: 255, G: 0, B: 0, A: 255}

	podstawa := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 40, 40, bialy), "podstawa")
	nakladka := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 10, 10, czerwony), "nakładka")

	var wynik shared.ImageComposeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageCompose,
		shared.ImageComposeRequest{
			WindowId:       wskaznik("okno-sprawdzianu"),
			BaseAssetId:    podstawa,
			OverlayAssetId: nakladka,
			X:              wskaznik(5),
			Y:              wskaznik(5),
		}, &wynik)

	obraz := obrazZMagazynu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	if obraz.Bounds().Dx() != 40 || obraz.Bounds().Dy() != 40 {
		t.Fatalf("złożenie zmieniło wymiary podstawy na %dx%d",
			obraz.Bounds().Dx(), obraz.Bounds().Dy())
	}

	// Piksel pod nakładką ma być czerwony.
	r, g, b, _ := obraz.At(8, 8).RGBA()
	if r>>8 < 200 || g>>8 > 60 || b>>8 > 60 {
		t.Fatalf("piksel pod nakładką ma barwę (%d,%d,%d) — nakładka nie legła",
			r>>8, g>>8, b>>8)
	}
	// Piksel poza nakładką ma zostać biały.
	r, g, b, _ = obraz.At(30, 30).RGBA()
	if r>>8 < 200 || g>>8 < 200 || b>>8 < 200 {
		t.Fatalf("piksel poza nakładką ma barwę (%d,%d,%d) — nakładka rozlała się "+
			"poza swoje miejsce", r>>8, g>>8, b>>8)
	}
}

// TestKrycieNakladkiDajeBarwePosrednia wykazuje, że krycie jest liczone, a nie
// przemilczane: przy połowie krycia czerwień na bieli daje róż, a nie czerwień.
func TestKrycieNakladkiDajeBarwePosrednia(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	bialy := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	czerwony := color.NRGBA{R: 255, G: 0, B: 0, A: 255}

	podstawa := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 20, 20, bialy), "podstawa")
	nakladka := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 20, 20, czerwony), "nakładka")

	var wynik shared.ImageComposeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageCompose,
		shared.ImageComposeRequest{
			WindowId:       wskaznik("okno-sprawdzianu"),
			BaseAssetId:    podstawa,
			OverlayAssetId: nakladka,
			Opacity:        wskaznik(50),
		}, &wynik)

	obraz := obrazZMagazynu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	_, g, b, _ := obraz.At(10, 10).RGBA()
	// Składowe mają wylądować pośrodku między bielą a czerwienią; przedział jest szeroki.
	if g>>8 < 100 || g>>8 > 160 || b>>8 < 100 || b>>8 > 160 {
		t.Fatalf("przy pięćdziesięcioprocentowym kryciu piksel ma składowe "+
			"(g=%d, b=%d) — krycie nie zostało policzone", g>>8, b>>8)
	}
}

// TestMieszanieMnozeniemPrzyciemniaObraz wykazuje, że tryb mieszania naprawdę
// zmienia arytmetykę: mnożenie czerwieni przez biel zostawia czerwień, a
// mnożenie przez czerń daje czerń.
func TestMieszanieMnozeniemPrzyciemniaObraz(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	czerwony := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	czarny := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	podstawa := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 16, 16, czerwony), "podstawa")
	nakladka := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 16, 16, czarny), "nakładka")

	var wynik shared.ImageComposeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageCompose,
		shared.ImageComposeRequest{
			WindowId:       wskaznik("okno-sprawdzianu"),
			BaseAssetId:    podstawa,
			OverlayAssetId: nakladka,
			BlendMode:      wskaznik(shared.ImageBlendMode(shared.ImageBlendModeMultiply)),
		}, &wynik)

	obraz := obrazZMagazynu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	r, g, b, _ := obraz.At(8, 8).RGBA()
	if r>>8 > 10 || g>>8 > 10 || b>>8 > 10 {
		t.Fatalf("mnożenie przez czerń dało barwę (%d,%d,%d) zamiast czerni",
			r>>8, g>>8, b>>8)
	}
}

// TestWektoryzacjaDajeSciezkiZObrysu wykazuje, że za wynikiem leży dokument SVG
// z prawdziwymi ścieżkami — a liczba w odpowiedzi zgadza się z liczbą ścieżek
// w dokumencie.
func TestWektoryzacjaDajeSciezkiZObrysu(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	// Materiał: biały kwadrat z czarnym prostokątem pośrodku — jeden kształt
	// o znanym konturze.
	plotno := image.NewNRGBA(image.Rect(0, 0, 40, 40))
	for y := 0; y < 40; y++ {
		for x := 0; x < 40; x++ {
			plotno.Set(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	for y := 12; y < 28; y++ {
		for x := 12; x < 28; x++ {
			plotno.Set(x, y, color.NRGBA{A: 255})
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć materiału: %v", err)
	}
	kod := wniesObraz(t, zmontowany, zycie, bufor.Bytes(), "kwadrat")

	var wynik shared.ImageVectorizeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageVectorize,
		shared.ImageVectorizeRequest{
			WindowId: wskaznik("okno-sprawdzianu"),
			AssetId:  wskaznik(kod),
			Mode:     wskaznik(shared.ImageVectorizeMode(shared.ImageVectorizeModeOutline)),
		}, &wynik)

	if wynik.Paths == nil || *wynik.Paths < 1 {
		t.Fatalf("wektoryzacja melduje %v ścieżek dla obrazu z kształtem", wynik.Paths)
	}

	sciezka := sciezkaZasobuSprawdzianu(t, zmontowany, zycie, katalog, wynik.Asset.Id)
	dokument := string(bajtyPlikuWektoryzacji(t, sciezka))

	if !strings.HasPrefix(strings.TrimSpace(dokument), "<svg") {
		t.Fatalf("za wynikiem nie leży dokument SVG; początek: %.80s", dokument)
	}
	if !strings.Contains(dokument, `viewBox="0 0 40 40"`) {
		t.Fatal("dokument nie niesie obszaru roboczego o wymiarach materiału")
	}
	sciezekWDokumencie := strings.Count(dokument, "<path ")
	if sciezekWDokumencie != *wynik.Paths {
		t.Fatalf("odpowiedź melduje %d ścieżek, a w dokumencie jest ich %d",
			*wynik.Paths, sciezekWDokumencie)
	}
	if !strings.Contains(dokument, ` d="M`) {
		t.Fatal("ścieżki dokumentu nie mają danych geometrycznych — pusty korpus SVG")
	}
}

// TestWektoryzacjaObrazuJednolitegoOdmawia wykazuje sprawdzian przeciwny: obraz
// bez kształtu nie daje pustego dokumentu podanego jako wektoryzacja, tylko
// odmowę nazywającą powód.
func TestWektoryzacjaObrazuJednolitegoOdmawia(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	bialy := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	kod := wniesObraz(t, zmontowany, zycie, obrazJednolity(t, 24, 24, bialy), "pusty")

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandImageVectorize,
		shared.ImageVectorizeRequest{
			WindowId: wskaznik("okno-sprawdzianu"), AssetId: wskaznik(kod),
		})
	if odmowa.Message == "" {
		t.Fatal("odmowa wektoryzacji nie mówi, dlaczego obrys niczego nie znalazł")
	}
}

// bajtyPlikuWektoryzacji czyta dokument wyniku wprost z dysku. Sprawdzian
// skutku czyta bajty, a nie odpowiedź: dokument pusty przechodzi każdy warunek
// istnienia i niczego nie niesie.
func bajtyPlikuWektoryzacji(t *testing.T, sciezka string) []byte {
	t.Helper()

	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		t.Fatalf("nie można odczytać wyniku spod %s: %v", sciezka, err)
	}
	if len(bajty) == 0 {
		t.Fatalf("wynik pod %s ma zerową długość", sciezka)
	}
	return bajty
}

// TestRozkladNaWarstwyDajeOsobneZasobyZPrzezroczystoscia wykazuje skutek rozkładu:
// każdy obiekt wychodzi osobnym zasobem, a nie kopią całego obrazu.
func TestRozkladNaWarstwyDajeOsobneZasobyZPrzezroczystoscia(t *testing.T) {
	if !zewnetrzne.Stoi(narzedzieWycinaniaTla()) {
		t.Skip("rozkład na warstwy stoi na sieci segmentującej (rembg), której nie ma na tej maszynie")
	}
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	// Materiał: dwa rozdzielone kwadraty na jednolitym tle, plan pierwszy sieci.
	plotno := image.NewNRGBA(image.Rect(0, 0, 120, 60))
	for y := 0; y < 60; y++ {
		for x := 0; x < 120; x++ {
			plotno.Set(x, y, color.NRGBA{R: 250, G: 250, B: 250, A: 255})
		}
	}
	for y := 15; y < 45; y++ {
		for x := 10; x < 40; x++ {
			plotno.Set(x, y, color.NRGBA{R: 10, G: 10, B: 200, A: 255})
		}
		for x := 80; x < 110; x++ {
			plotno.Set(x, y, color.NRGBA{R: 200, G: 10, B: 10, A: 255})
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		t.Fatalf("nie można złożyć materiału: %v", err)
	}
	kod := wniesObraz(t, zmontowany, zycie, bufor.Bytes(), "dwa-obiekty")

	var wynik shared.ImageLayersSplitResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandImageLayersSplit,
		shared.ImageLayersSplitRequest{
			WindowId: wskaznik("okno-sprawdzianu"), AssetId: wskaznik(kod),
		}, &wynik)

	if wynik.Layers != len(wynik.Assets) {
		t.Fatalf("odpowiedź melduje %d warstw i %d zasobów", wynik.Layers, len(wynik.Assets))
	}
	if wynik.Layers < 1 {
		t.Fatal("rozkład nie oddał ani jednej warstwy")
	}

	for numer, zasob := range wynik.Assets {
		warstwa := obrazZMagazynu(t, zmontowany, zycie, katalog, zasob.Id)
		// Warstwa ma być WYCINKIEM, a nie całym kadrem: rozkład, który oddaje
		// obraz źródłowy, jest atrapą.
		if warstwa.Bounds().Dx() >= 120 && warstwa.Bounds().Dy() >= 60 {
			t.Fatalf("warstwa %d ma wymiary całego materiału (%dx%d) — to nie jest rozkład",
				numer, warstwa.Bounds().Dx(), warstwa.Bounds().Dy())
		}
		// Warstwa ma nieść przezroczystość: choć jeden piksel poza obiektem.
		przezroczysty := false
		for y := warstwa.Bounds().Min.Y; y < warstwa.Bounds().Max.Y && !przezroczysty; y++ {
			for x := warstwa.Bounds().Min.X; x < warstwa.Bounds().Max.X; x++ {
				if _, _, _, alfa := warstwa.At(x, y).RGBA(); alfa == 0 {
					przezroczysty = true
					break
				}
			}
		}
		if !przezroczysty && warstwa.Bounds().Dx()*warstwa.Bounds().Dy() > 4 {
			t.Fatalf("warstwa %d nie ma ani jednego piksela przezroczystego", numer)
		}
	}
}
