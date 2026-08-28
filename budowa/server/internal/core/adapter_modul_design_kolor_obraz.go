// Odpowiedzialność pliku: dwie czynności barwy, które czytają piksele —
// wyciągnięcie palety z obrazu i symulacja wady widzenia barw. Czynności
// liczone bez obrazu leżą w `adapter_modul_design_kolor.go`.
package core

import (
	"context"
	"fmt"
	"image"
	"math"
	"sort"
	"strings"

	"github.com/lucasb-eyer/go-colorful"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// granicaPunktowPomiaruPaletyDesignu jest liczbą punktów próbki palety,
	// powyżej której rachunek próbkuje obraz zamiast liczyć każdy piksel.
	granicaPunktowPomiaruPaletyDesignu = 40000

	// obrotowSkupianiaPaletyDesignu jest liczbą przebiegów k-średnich, zapasem
	// ponad liczbę potrzebną do ustabilizowania skupień.
	obrotowSkupianiaPaletyDesignu = 20
)

// WyciagnijPalete liczy barwy dominujące obrazu skupianiem k-średnich —
// obsługuje `design.color.palette.extract`.
func (a *adapterDesignu) WyciagnijPalete(ctx context.Context,
	z shared.DesignColorPaletteExtractRequest) (shared.DesignColorPaletteExtractResponse, error) {

	if strings.TrimSpace(z.AssetId) == "" {
		return shared.DesignColorPaletteExtractResponse{}, bladWskazaniaDesignu(
			"komenda design.color.palette.extract bez wskazania zasobu")
	}
	ile := domyslnaLiczbaBarwPaletyDesignu
	if z.Count != nil {
		if *z.Count < 1 || *z.Count > granicaBarwPaletyDesignu {
			return shared.DesignColorPaletteExtractResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.color.palette.extract z liczbą barw %d: rdzeń liczy palety od 1 "+
					"do %d barw", *z.Count, granicaBarwPaletyDesignu))
		}
		ile = *z.Count
	}

	obraz, _, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.color.palette.extract", z.AssetId)
	if err != nil {
		return shared.DesignColorPaletteExtractResponse{}, err
	}

	punkty := probkujPunktyObrazuDesignu(obraz, granicaPunktowPomiaruPaletyDesignu)
	if len(punkty) == 0 {
		return shared.DesignColorPaletteExtractResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"zasób %s nie ma ani jednego punktu do zmierzenia — obraz o zerowej powierzchni nie "+
				"ma barw dominujących", strings.TrimSpace(z.AssetId)))
	}
	if ile > len(punkty) {
		ile = len(punkty)
	}

	skupienia, udzialy := skupieniaBarwDesignu(punkty, ile)
	barwy := make([]shared.DesignPaletteColor, 0, len(skupienia))
	for numer, skupienie := range skupienia {
		udzial := udzialy[numer]
		wpis := shared.DesignPaletteColor{
			Hex: skupienie.Clamped().Hex(), Share: &udzial,
		}
		if nazwa := nazwaBarwyDesignu(wpis.Hex); nazwa != nil {
			wpis.Name = nazwa
		}
		barwy = append(barwy, wpis)
	}
	// Barwy w kolejności udziału, tak jak opisuje pole kontrakt.
	sort.SliceStable(barwy, func(i, j int) bool {
		return *barwy[i].Share > *barwy[j].Share
	})
	zmierzonych := len(punkty)
	return shared.DesignColorPaletteExtractResponse{
		Colors: barwy, PixelsSampled: &zmierzonych,
	}, nil
}

// probkujPunktyObrazuDesignu bierze z obrazu równomierną próbę punktów
// w przestrzeni Lab. Punkty całkowicie przezroczyste nie wchodzą do próby.
func probkujPunktyObrazuDesignu(obraz image.Image, granica int) []colorful.Color {
	granice := obraz.Bounds()
	szerokosc, wysokosc := granice.Dx(), granice.Dy()
	if szerokosc <= 0 || wysokosc <= 0 {
		return nil
	}
	// Krok próbkowania liczony z powierzchni obrazu wobec granicy punktów.
	krok := 1
	for (szerokosc/krok)*(wysokosc/krok) > granica {
		krok++
	}
	punkty := make([]colorful.Color, 0, granica)
	for y := granice.Min.Y; y < granice.Max.Y; y += krok {
		for x := granice.Min.X; x < granice.Max.X; x += krok {
			r, g, b, alfa := obraz.At(x, y).RGBA()
			if alfa == 0 {
				continue
			}
			punkty = append(punkty, colorful.Color{
				R: float64(r) / 65535, G: float64(g) / 65535, B: float64(b) / 65535,
			})
		}
	}
	return punkty
}

// skupieniaBarwDesignu liczy k-średnich w przestrzeni Lab i oddaje środki
// skupień wraz z udziałem punktów w każdym. Zasiew jest rozłożony po
// posortowanej próbie, nie losowy, żeby wynik był powtarzalny.
func skupieniaBarwDesignu(punkty []colorful.Color, ile int) ([]colorful.Color, []float64) {
	if ile < 1 {
		ile = 1
	}
	rodzaj := make([]colorful.Color, len(punkty))
	copy(rodzaj, punkty)
	// Porządek po jasności percepcyjnej: zasiew obejmuje pełny zakres barw.
	sort.SliceStable(rodzaj, func(i, j int) bool {
		pierwsza, _, _ := rodzaj[i].Lab()
		druga, _, _ := rodzaj[j].Lab()
		return pierwsza < druga
	})
	srodki := make([]colorful.Color, 0, ile)
	for numer := 0; numer < ile; numer++ {
		indeks := numer * len(rodzaj) / ile
		srodki = append(srodki, rodzaj[indeks])
	}

	przypisania := make([]int, len(punkty))
	for obrot := 0; obrot < obrotowSkupianiaPaletyDesignu; obrot++ {
		zmienilo := false
		for numer, punkt := range punkty {
			najblizsze, najmniejsza := 0, math.MaxFloat64
			for numerSrodka, srodek := range srodki {
				odleglosc := punkt.DistanceLab(srodek)
				if odleglosc < najmniejsza {
					najblizsze, najmniejsza = numerSrodka, odleglosc
				}
			}
			if przypisania[numer] != najblizsze {
				przypisania[numer] = najblizsze
				zmienilo = true
			}
		}
		// Nowe środki liczone jako średnia w Lab, nie w sRGB.
		sumaL := make([]float64, len(srodki))
		sumaA := make([]float64, len(srodki))
		sumaB := make([]float64, len(srodki))
		liczba := make([]int, len(srodki))
		for numer, punkt := range punkty {
			jasnosc, a, b := punkt.Lab()
			skupienie := przypisania[numer]
			sumaL[skupienie] += jasnosc
			sumaA[skupienie] += a
			sumaB[skupienie] += b
			liczba[skupienie]++
		}
		for numer := range srodki {
			if liczba[numer] == 0 {
				continue
			}
			srodki[numer] = colorful.Lab(sumaL[numer]/float64(liczba[numer]),
				sumaA[numer]/float64(liczba[numer]), sumaB[numer]/float64(liczba[numer])).Clamped()
		}
		if !zmienilo {
			break
		}
	}

	liczba := make([]int, len(srodki))
	for _, skupienie := range przypisania {
		liczba[skupienie]++
	}
	// Skupienie puste nie wchodzi do palety.
	wynikSrodki := make([]colorful.Color, 0, len(srodki))
	wynikUdzialy := make([]float64, 0, len(srodki))
	for numer, srodek := range srodki {
		if liczba[numer] == 0 {
			continue
		}
		wynikSrodki = append(wynikSrodki, srodek)
		wynikUdzialy = append(wynikUdzialy, float64(liczba[numer])/float64(len(punkty)))
	}
	return wynikSrodki, wynikUdzialy
}

// SymulujWidzenie wydaje obraz przepuszczony przez symulację wady widzenia barw —
// obsługuje `design.color.vision.simulate`.
func (a *adapterDesignu) SymulujWidzenie(ctx context.Context,
	z shared.DesignColorVisionSimulateRequest) (shared.DesignColorVisionSimulateResponse, error) {

	if strings.TrimSpace(z.AssetId) == "" {
		return shared.DesignColorVisionSimulateResponse{}, bladWskazaniaDesignu(
			"komenda design.color.vision.simulate bez wskazania zasobu")
	}
	if err := sprawdzWyliczenieDesignu("design.color.vision.simulate", "vision", z.Vision,
		shared.WartosciDesignColorVision()); err != nil {
		return shared.DesignColorVisionSimulateResponse{}, err
	}
	obraz, zrodlo, err := a.obrazZasobuPoKodzieDesignu(ctx, "design.color.vision.simulate", z.AssetId)
	if err != nil {
		return shared.DesignColorVisionSimulateResponse{}, err
	}

	wynik := przepuscObrazPrzezWadeDesignu(obraz, z.Vision)
	bajty, _, err := zakodujObrazDesignu(wynik, "png", nil)
	if err != nil {
		return shared.DesignColorVisionSimulateResponse{}, bladWydaniaDesignu(err.Error())
	}

	nazwa := nazwaZasobuDesignu(zrodlo) + " — " + string(z.Vision)
	zapisany, err := a.zalozZasobZBajtowDesignu(ctx, oknoWytworu(z.WindowId, zrodlo.Okno),
		nazwa, shared.DesignAssetKindImage, "png", bajty)
	if err != nil {
		return shared.DesignColorVisionSimulateResponse{}, err
	}
	// Wariant wskazuje źródło: obraz do porównania przed i po symulacji.
	zapisany.WariantZasobuID = &zrodlo.Kod
	zapisany, err = a.repozytorium.ZapiszZasob(ctx, zapisany)
	if err != nil {
		return shared.DesignColorVisionSimulateResponse{}, bladDesignu(err)
	}
	return shared.DesignColorVisionSimulateResponse{
		Asset: zasobWytworzonyKontraktu(zapisany),
	}, nil
}

// macierzeWadWidzeniaDesignu to macierze przejścia LMS dla trzech wad
// dichromatycznych, wedle Brettela-Vienota-Mollona.
var macierzeWadWidzeniaDesignu = map[shared.DesignColorVision][9]float64{
	// Protanopia — brak czopka długofalowego (L).
	shared.DesignColorVisionProtanopia: {
		0.0, 2.02344, -2.52581,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
	},
	// Deuteranopia — brak czopka średniofalowego (M).
	shared.DesignColorVisionDeuteranopia: {
		1.0, 0.0, 0.0,
		0.494207, 0.0, 1.24827,
		0.0, 0.0, 1.0,
	},
	// Tritanopia — brak czopka krótkofalowego (S).
	shared.DesignColorVisionTritanopia: {
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		-0.395913, 0.801109, 0.0,
	},
}

// macierzSrgbNaLmsDesignu i macierzLmsNaSrgbDesignu to przejścia między sRGB
// liniowym i przestrzenią czopków (Hunt-Pointer-Estevez znormalizowany do D65).
var (
	macierzSrgbNaLmsDesignu = [9]float64{
		0.31399022, 0.63951294, 0.04649755,
		0.15537241, 0.75789446, 0.08670142,
		0.01775239, 0.10944209, 0.87256922,
	}
	macierzLmsNaSrgbDesignu = [9]float64{
		5.47221206, -4.6419601, 0.16963708,
		-1.1252419, 2.29317094, -0.1678952,
		0.02980165, -0.19318073, 1.16364789,
	}
)

// przepuscObrazPrzezWadeDesignu liczy obraz widziany przez wadę widzenia
// barw, punkt po punkcie, bez próbkowania oszczędzającego rachunek.
func przepuscObrazPrzezWadeDesignu(obraz image.Image,
	wada shared.DesignColorVision) image.Image {

	granice := obraz.Bounds()
	wynik := image.NewRGBA(image.Rect(0, 0, granice.Dx(), granice.Dy()))
	macierz, dichromatyczna := macierzeWadWidzeniaDesignu[wada]
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			r, g, b, alfa := obraz.At(x, y).RGBA()
			zrodlowa := colorful.Color{
				R: float64(r) / 65535, G: float64(g) / 65535, B: float64(b) / 65535,
			}
			var docelowa colorful.Color
			if !dichromatyczna {
				// Achromatopsja: luminancja względna wedle wag WCAG.
				luminancja := luminancjaWcagDesignu(zrodlowa)
				szara := gammaSrgbDesignu(luminancja)
				docelowa = colorful.Color{R: szara, G: szara, B: szara}
			} else {
				docelowa = przepuscBarwePrzezMacierzDesignu(zrodlowa, macierz)
			}
			// Krycie źródła zostaje: symulacja zmienia barwę, nie przezroczystość.
			kanal := float64(alfa) / 65535
			wynik.SetRGBA(x-granice.Min.X, y-granice.Min.Y, barwaRgbaDesignu(docelowa, &kanal))
		}
	}
	return wynik
}

// przepuscBarwePrzezMacierzDesignu przenosi barwę do przestrzeni czopków,
// stosuje macierz wady i wraca do sRGB.
func przepuscBarwePrzezMacierzDesignu(barwa colorful.Color, wada [9]float64) colorful.Color {
	r, g, b := barwa.Clamped().LinearRgb()
	l, m, s := pomnozMacierzaDesignu(macierzSrgbNaLmsDesignu, r, g, b)
	lWada, mWada, sWada := pomnozMacierzaDesignu(wada, l, m, s)
	rWynik, gWynik, bWynik := pomnozMacierzaDesignu(macierzLmsNaSrgbDesignu, lWada, mWada, sWada)
	return colorful.LinearRgb(przytnijUlamekDesignu(rWynik), przytnijUlamekDesignu(gWynik),
		przytnijUlamekDesignu(bWynik))
}

// pomnozMacierzaDesignu mnoży wektor trzech składowych przez macierz 3×3
// zapisaną wierszami, do przejść między przestrzeniami barwy.
func pomnozMacierzaDesignu(macierz [9]float64, pierwsza, druga, trzecia float64) (float64, float64, float64) {
	return macierz[0]*pierwsza + macierz[1]*druga + macierz[2]*trzecia,
		macierz[3]*pierwsza + macierz[4]*druga + macierz[5]*trzecia,
		macierz[6]*pierwsza + macierz[7]*druga + macierz[8]*trzecia
}

// gammaSrgbDesignu zakłada gamma sRGB na składową liniową — odwrotność
// zdejmowania gamma, którym liczy się luminancja WCAG.
func gammaSrgbDesignu(liniowa float64) float64 {
	liniowa = przytnijUlamekDesignu(liniowa)
	if liniowa <= 0.0031308 {
		return liniowa * 12.92
	}
	return 1.055*math.Pow(liniowa, 1/2.4) - 0.055
}

// obrazZasobuPoKodzieDesignu odczytuje obraz zasobu wraz z jego wierszem —
// droga wspólna dla wszystkich czynności modułu, które czytają piksele.
// Zasób bez treści w magazynie jest odmową, nie pustym obrazem.
func (a *adapterDesignu) obrazZasobuPoKodzieDesignu(ctx context.Context, komenda,
	kod string) (image.Image, dane.ZasobDesignu, error) {

	zasob, err := a.repozytorium.Zasob(ctx, strings.TrimSpace(kod))
	if err != nil {
		return nil, dane.ZasobDesignu{}, bladNieznanegoZasobuDesignu(kod, err)
	}
	if zasob.URI == nil || strings.TrimSpace(*zasob.URI) == "" {
		return nil, dane.ZasobDesignu{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s: zasób %s nie ma treści w magazynie — nie ma czego zmierzyć",
			komenda, zasob.Kod))
	}
	obraz, err := obrazZasobuDesignu(*zasob.URI)
	if err != nil {
		return nil, dane.ZasobDesignu{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s: zasób %s nie jest obrazem, który rdzeń potrafi rozłożyć: %s",
			komenda, zasob.Kod, err.Error()))
	}
	return obraz, zasob, nil
}

// nazwaZasobuDesignu oddaje nazwę zasobu do podpisania wytworu — nadaną przez
// Operatora albo jego identyfikator.
func nazwaZasobuDesignu(zasob dane.ZasobDesignu) string {
	if zasob.Nazwa != nil && strings.TrimSpace(*zasob.Nazwa) != "" {
		return strings.TrimSpace(*zasob.Nazwa)
	}
	return zasob.Kod
}
