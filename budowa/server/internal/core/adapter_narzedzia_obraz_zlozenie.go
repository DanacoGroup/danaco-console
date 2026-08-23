// Odpowiedzialność pliku: `image.compose` — złożenie obrazu na obrazie wedle
// trybu mieszania, krycia, położenia i skali.
//
// ── Biblioteka wkompilowana, nie program zewnętrzny ─────────────────────────
// Cała czynność stoi na `image`, `image/png`, `image/jpeg` i
// `golang.org/x/image/draw` — bibliotekach wkompilowanych w binarium rdzenia.
// Nie startuje tu ani jeden proces potomny. Złożenie dwóch rastrów to
// przejście po pikselach, a nie praca, do której potrzeba silnika: wołanie
// ImageMagicka odebrałoby tej komendzie prawo do działania na maszynie, która
// go nie ma, i to bez żadnego zysku.
//
// Jedna komenda zamyka cztery funkcje opracowania: znak wodny, branding
// wsadowy, osadzenie w ramce urządzenia i warstwy rastrowe. Wszystkie cztery
// potrzebują dokładnie tego samego — obrazu położonego na obrazie.
//
// Wynik jest zawsze PNG. Podstawa bywa fotografią bez kanału alfa, ale nakładka
// z przezroczystością wnosi go do wyniku; zapis w JPEG-u zamieniłby
// przezroczystość na czarny albo biały prostokąt. PNG nie traci też jakości
// przy powtórnym składaniu, a branding wsadowy jest właśnie powtarzaniem.
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

// Zloz obsługuje `image.compose`.
//
// Kolejność jest zamierzona: najpierw oba źródła (bo bez któregoś nie ma czego
// składać), potem skala nakładki, potem mieszanie. Nakładka wychodząca poza
// obszar podstawy NIE jest odmową — znak wodny wypuszczony za krawędź jest
// przycinany, tak jak przycięłaby go każda inna warstwa graficzna. Odmowa
// kazałaby Operatorowi liczyć piksele, zamiast przesunąć nakładkę.
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

	// Podstawa wchodzi do bufora RGBA, bo rysowanie idzie po niej wielokrotnie,
	// a obraz zdekodowany z PNG bywa `image.Paletted` albo `image.YCbCr`, po
	// których biblioteka standardowa nie umie rysować.
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
// Rola („podstawy", „nakładki") wchodzi do treści odmowy: bez niej Operator
// z dwoma zasobami w żądaniu nie wie, który z nich nie jest obrazem.
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
// krotności i krotność 1 zostawiają nakładkę nietkniętą — przeskalowanie
// jeden do jednego byłoby przepróbkowaniem bez powodu, a każde przepróbkowanie
// kosztuje ostrość.
//
// Filtr `CatmullRom` bierzemy zarówno przy pomniejszaniu, jak i powiększaniu:
// znak wodny pomniejszany najbliższym sąsiadem rozsypuje się na schodki
// widoczne gołym okiem.
func przeskalujNakladke(nakladka image.Image, krotnosc *float64) image.Image {
	if krotnosc == nil || *krotnosc <= 0 || math.Abs(*krotnosc-1) < 0.0001 {
		return nakladka
	}
	granice := nakladka.Bounds()
	szerokosc := int(math.Round(float64(granice.Dx()) * *krotnosc))
	wysokosc := int(math.Round(float64(granice.Dy()) * *krotnosc))
	if szerokosc < 1 || wysokosc < 1 {
		// Krotność, po której nie zostaje ani jeden piksel, nie jest odmową:
		// nakładka o zerowej wielkości niczego nie zmienia, a odmowa
		// zatrzymałaby wsad z jednym źle policzonym wierszem.
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
// krycie — tak mówi kontrakt („brak bierze pełne").
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

// zmieszajNakladke rysuje nakładkę na płótnie piksel po pikselu.
//
// Własna pętla, a nie `draw.Draw`, bo tryby mieszania inne niż zwykły nie mają
// odpowiednika w bibliotece standardowej, a mieszanie „prawie zwykłe" byłoby
// trybem, którego nazwa mówi co innego niż skutek.
//
// Barwy liczymy bez wstępnego mnożenia przez alfę (`NRGBA`), bo tryby
// mieszania są zdefiniowane na barwie własnej piksela; praca na barwie już
// przemnożonej przez alfę dałaby mnożenie ciemniejsze przy każdej
// półprzezroczystości.
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
// 0..1. Tryb spoza wyliczenia kontraktu kończy się odmową nazywającą go wprost:
// gałąź domyślna „mieszaj zwykle" oddałaby złożenie, które wygląda poprawnie
// i nie jest tym, o co proszono.
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
