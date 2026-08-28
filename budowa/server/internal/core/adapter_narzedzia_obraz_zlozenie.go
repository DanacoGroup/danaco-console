// Plik obsługuje `image.compose`: składa obraz na obrazie wedle trybu
// mieszania, krycia, położenia i skali, korzystając wyłącznie z bibliotek
// wkompilowanych, bez procesu potomnego, i zawsze zwraca PNG.
package core

import (
	"bytes"
	"context"
	"image"
	"image/draw"
	"image/png"
	"math"
	"strings"

	rysowanie "golang.org/x/image/draw"

	"danacoconsole/shared"
)

// granicaPikseliZlozenia chroni przed obrazem, którego złożenie zjadłoby całą
// pamięć maszyny. Sto milionów pikseli to ponad 380 MB w postaci RGBA — więcej,
// niż ma sens dla znaku wodnego, a mniej, niż wywraca serwer.
const granicaPikseliZlozenia = 100_000_000

// Zloz obsługuje `image.compose`. Kolejność jest zamierzona: najpierw oba
// źródła, potem skala nakładki, potem mieszanie. Nakładka wychodząca poza
// obszar podstawy nie jest odmową, jest przycinana jak każda inna warstwa
// graficzna.
func (a *adapterNarzedziObrazu) Zloz(ctx context.Context,
	z shared.ImageComposeRequest) (shared.ImageComposeResponse, error) {

	if strings.TrimSpace(z.BaseAssetId) == "" || strings.TrimSpace(z.OverlayAssetId) == "" {
		return shared.ImageComposeResponse{}, bladWskazaniaObrazu(
			"złożenie wymaga dwóch zasobów: baseAssetId (podstawa) i overlayAssetId (nakładka)")
	}

	zrodloPodstawy, err := a.zrodloZZasobu(ctx, strings.TrimSpace(z.BaseAssetId))
	if err != nil {
		return shared.ImageComposeResponse{}, err
	}
	zrodloNakladki, err := a.zrodloZZasobu(ctx, strings.TrimSpace(z.OverlayAssetId))
	if err != nil {
		return shared.ImageComposeResponse{}, err
	}

	podstawa, err := odczytajObrazZlozenia(zrodloPodstawy.sciezka, "podstawy")
	if err != nil {
		return shared.ImageComposeResponse{}, err
	}
	nakladka, err := odczytajObrazZlozenia(zrodloNakladki.sciezka, "nakładki")
	if err != nil {
		return shared.ImageComposeResponse{}, err
	}

	nakladka = przeskalujNakladke(nakladka, z.Scale)

	// Podstawa wchodzi do bufora RGBA, bo dekodowany obraz bywa niemalowalny.
	plotno := image.NewRGBA(podstawa.Bounds())
	draw.Draw(plotno, plotno.Bounds(), podstawa, podstawa.Bounds().Min, draw.Src)

	krycie := krycieZlozenia(z.Opacity)
	var tryb shared.ImageBlendMode = shared.ImageBlendModeNormal
	if z.BlendMode != nil && strings.TrimSpace(string(*z.BlendMode)) != "" {
		tryb = *z.BlendMode
	}
	przesuniecie := image.Pt(wartoscLiczby(z.X), wartoscLiczby(z.Y))

	if err := zmieszajNakladke(plotno, nakladka, przesuniecie, krycie, tryb); err != nil {
		return shared.ImageComposeResponse{}, err
	}

	var bufor bytes.Buffer
	if err := png.Encode(&bufor, plotno); err != nil {
		return shared.ImageComposeResponse{}, bladPrzetwarzaniaObrazu(
			"nie można zapisać złożonego obrazu: " + err.Error())
	}

	zasob, _, err := a.odlozZasob(ctx, zrodloPodstawy, z.WindowId, bufor.Bytes(), "png", "złożenie")
	if err != nil {
		return shared.ImageComposeResponse{}, err
	}
	return shared.ImageComposeResponse{Asset: zasob}, nil
}

// odczytajObrazZlozenia dekoduje plik do postaci, po której da się rysować.
// Rola wchodzi do treści odmowy: bez niej Operator z dwoma zasobami w żądaniu
// nie wie, który z nich nie jest obrazem.
func odczytajObrazZlozenia(sciezka, rola string) (image.Image, error) {
	obraz, err := odczytajObrazPliku(sciezka)
	if err != nil {
		return nil, bladWskazaniaObrazu("nie można odczytać obrazu " + rola + ": " + err.Error())
	}
	rozmiar := obraz.Bounds().Dx() * obraz.Bounds().Dy()
	if rozmiar <= 0 {
		return nil, bladWskazaniaObrazu("obraz " + rola + " nie ma ani jednego piksela")
	}
	if rozmiar > granicaPikseliZlozenia {
		return nil, bladWskazaniaObrazu("obraz " + rola +
			" jest za duży do złożenia w pamięci rdzenia (ponad sto milionów pikseli)")
	}
	return obraz, nil
}

// przeskalujNakladke zmienia rozmiar nakładki krotnością z żądania. Brak
// krotności i krotność 1 zostawiają nakładkę nietkniętą, a filtr CatmullRom
// jest brany zarówno przy pomniejszaniu, jak i powiększaniu.
func przeskalujNakladke(nakladka image.Image, krotnosc *float64) image.Image {
	if krotnosc == nil || *krotnosc <= 0 || math.Abs(*krotnosc-1) < 0.0001 {
		return nakladka
	}
	granice := nakladka.Bounds()
	szerokosc := int(math.Round(float64(granice.Dx()) * *krotnosc))
	wysokosc := int(math.Round(float64(granice.Dy()) * *krotnosc))
	if szerokosc < 1 || wysokosc < 1 {
		// Zerowa wielkość nakładki nie jest odmową, niczego nie zmienia.
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}
	if szerokosc*wysokosc > granicaPikseliZlozenia {
		return nakladka
	}
	cel := image.NewRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	rysowanie.CatmullRom.Scale(cel, cel.Bounds(), nakladka, granice, rysowanie.Over, nil)
	return cel
}

// krycieZlozenia przekłada procent kontraktu na mnożnik. Brak pola znaczy pełne
// krycie, tak jak stanowi kontrakt.
func krycieZlozenia(procent *int) float64 {
	if procent == nil {
		return 1
	}
	wartosc := float64(*procent) / 100
	switch {
	case wartosc < 0:
		return 0
	case wartosc > 1:
		return 1
	default:
		return wartosc
	}
}

// zmieszajNakladke rysuje nakładkę na płótnie piksel po pikselu, własną pętlą
// zamiast `draw.Draw`, bo tryby mieszania inne niż zwykły nie mają odpowiednika
// w bibliotece standardowej. Barwy liczy bez wstępnego mnożenia przez alfę.
func zmieszajNakladke(plotno *image.RGBA, nakladka image.Image, przesuniecie image.Point,
	krycie float64, tryb shared.ImageBlendMode) error {

	mieszanie, err := funkcjaMieszania(tryb)
	if err != nil {
		return err
	}
	granice := nakladka.Bounds()
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			celX := przesuniecie.X + (x - granice.Min.X)
			celY := przesuniecie.Y + (y - granice.Min.Y)
			if !(image.Point{X: celX, Y: celY}).In(plotno.Bounds()) {
				continue
			}
			gornaR, gornaG, gornaB, gornaA := skladoweNiePrzemnozone(nakladka.At(x, y))
			udzial := gornaA * krycie
			if udzial <= 0 {
				continue
			}
			dolnaR, dolnaG, dolnaB, dolnaA := skladoweNiePrzemnozone(plotno.At(celX, celY))

			wynikR := dolnaR + (mieszanie(dolnaR, gornaR)-dolnaR)*udzial
			wynikG := dolnaG + (mieszanie(dolnaG, gornaG)-dolnaG)*udzial
			wynikB := dolnaB + (mieszanie(dolnaB, gornaB)-dolnaB)*udzial
			wynikA := dolnaA + (1-dolnaA)*udzial

			plotno.Set(celX, celY, kolorZeSkladowych(wynikR, wynikG, wynikB, wynikA))
		}
	}
	return nil
}

// funkcjaMieszania oddaje działanie jednego trybu na parze składowych z zakresu
// 0..1. Tryb spoza wyliczenia kontraktu kończy się odmową nazywającą go wprost.
func funkcjaMieszania(tryb shared.ImageBlendMode) (func(dolna, gorna float64) float64, error) {
	switch tryb {
	case shared.ImageBlendModeNormal:
		return func(_, gorna float64) float64 { return gorna }, nil
	case shared.ImageBlendModeMultiply:
		return func(dolna, gorna float64) float64 { return dolna * gorna }, nil
	case shared.ImageBlendModeScreen:
		return func(dolna, gorna float64) float64 { return 1 - (1-dolna)*(1-gorna) }, nil
	case shared.ImageBlendModeOverlay:
		return func(dolna, gorna float64) float64 {
			if dolna <= 0.5 {
				return 2 * dolna * gorna
			}
			return 1 - 2*(1-dolna)*(1-gorna)
		}, nil
	default:
		return nil, bladWskazaniaObrazu("nie znam trybu mieszania „" + string(tryb) +
			"” — rdzeń zna: normal, multiply, screen, overlay")
	}
}
