// Plik obsługuje `image.layers.split`: rozkłada obraz na obiekty, z których każdy wychodzi osobnym zasobem PNG z przezroczystością, a bez silnika segmentacji odmawia zamiast zwracać atrapę rozkładu.
package core

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"path/filepath"
	"strconv"
	"strings"

	"danacoconsole/shared"
)

const (
	// progAlfyWarstwy rozstrzyga, który piksel należy do obiektu: alfa poniżej
	// progu liczy się jako tło, powyżej jako obiekt, a próg służy wyłącznie
	// wyznaczeniu obszarów spójnych.
	progAlfyWarstwy = 128

	// najmniejszaWarstwa odrzuca obszary mniejsze niż sto pikseli. Sieć zostawia
	// przy krawędziach pojedyncze wysepki, które jako osobne zasoby zaśmieciłyby
	// panel dziesiątkami kafelków wielkości kropki.
	najmniejszaWarstwa = 100

	// domyslnaGranicaWarstw jest brana, gdy żądanie nie stawia własnej.
	// Dwadzieścia warstw to więcej, niż ma sensownie obraz, a mniej, niż
	// wywraca panel.
	domyslnaGranicaWarstw = 20
)

// RozlozNaWarstwy obsługuje `image.layers.split`, rozkładając obraz na warstwy
// obiektów zapisywane jako osobne zasoby PNG z przezroczystością wokół
// wyciętego kształtu.
func (a *adapterNarzedziObrazuModelu) RozlozNaWarstwy(ctx context.Context,
	z shared.ImageLayersSplitRequest) (shared.ImageLayersSplitResponse, error) {

	nazwaModelu, opis, err := rozstrzygnijModelWycinania(z.Model)
	if err != nil {
		return shared.ImageLayersSplitResponse{}, err
	}
	if a.wspolne == nil {
		return shared.ImageLayersSplitResponse{}, bladZapleczaModeluObrazu(
			"zaplecze narzędzi obrazu nie jest wpięte — nie ma czym rozwiązać źródła")
	}
	zrodlo, err := a.wspolne.rozwiazZrodlo(ctx, z.AssetId, z.SourcePath)
	if err != nil {
		return shared.ImageLayersSplitResponse{}, err
	}
	if err := sprawdzWagi(filepath.Join(katalogWagWycinania(), opis.plik),
		nazwaModelu, opis.waga, opis.skad); err != nil {
		return shared.ImageLayersSplitResponse{}, err
	}

	pracownia, err := przygotujPracownie(zrodlo.sciezka)
	if err != nil {
		return shared.ImageLayersSplitResponse{}, err
	}
	defer pracownia.sprzatnij()

	argumenty := []string{"i", "-m", nazwaModelu, pracownia.wejscie, pracownia.wyjscie}
	if err := a.wolajSilnik(ctx, narzedzieWycinaniaTla(), argumenty, granicaWycinaniaTla); err != nil {
		return shared.ImageLayersSplitResponse{}, err
	}

	rozdzielony, err := odczytajObrazPliku(pracownia.wyjscie)
	if err != nil {
		return shared.ImageLayersSplitResponse{}, bladPrzetwarzaniaModeluObrazu(
			"silnik segmentacji zostawił plik, którego nie da się odczytać jako obrazu: " + err.Error())
	}

	obszary := obszarySpojneWarstw(rozdzielony)
	if len(obszary) == 0 {
		return shared.ImageLayersSplitResponse{}, bladPrzetwarzaniaModeluObrazu(
			"silnik segmentacji nie znalazł na tym obrazie ani jednego obiektu — " +
				"cały kadr został uznany za tło; rdzeń nie oddaje wtedy obrazu " +
				"źródłowego jako rzekomej jednej warstwy")
	}

	granica := domyslnaGranicaWarstw
	if z.MaxLayers != nil && *z.MaxLayers > 0 {
		granica = *z.MaxLayers
	}
	if len(obszary) > granica {
		obszary = obszary[:granica]
	}

	zasoby := make([]shared.DesignAsset, 0, len(obszary))
	for numer, obszar := range obszary {
		bajty, err := wytnijWarstwe(rozdzielony, obszar)
		if err != nil {
			return shared.ImageLayersSplitResponse{}, err
		}
		zasob, _, err := a.wspolne.odlozZasob(ctx, zrodlo, z.WindowId, bajty, "png",
			"warstwa "+strconv.Itoa(numer+1))
		if err != nil {
			return shared.ImageLayersSplitResponse{}, err
		}
		zasoby = append(zasoby, zasob)
	}
	return shared.ImageLayersSplitResponse{Assets: zasoby, Layers: len(zasoby)}, nil
}

// obszarWarstwy jest jednym obiektem: prostokątem obejmującym i przynależnością
// pikseli, trzymaną osobno, bo obiekty bywają wklęsłe i ich prostokąty się
// nachodzą.
type obszarWarstwy struct {
	prostokat image.Rectangle
	nalezy    *maskaRastrowa
	pole      int
}

// obszarySpojneWarstw dzieli plan pierwszy na obszary spójne, od największego,
// licząc spójność w sąsiedztwie ośmiu pikseli, tak by obiekt przewężony do
// linii ukośnej pozostał jednym kształtem.
func obszarySpojneWarstw(obraz image.Image) []obszarWarstwy {
	granice := obraz.Bounds()
	szerokosc, wysokosc := granice.Dx(), granice.Dy()
	if szerokosc <= 0 || wysokosc <= 0 {
		return nil
	}
	plan := nowaMaskaRastrowa(szerokosc, wysokosc)
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			_, _, _, alfa := obraz.At(granice.Min.X+x, granice.Min.Y+y).RGBA()
			plan.ustaw(x, y, alfa>>8 >= progAlfyWarstwy)
		}
	}

	odwiedzone := nowaMaskaRastrowa(szerokosc, wysokosc)
	obszary := []obszarWarstwy{}
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			if !plan.wewnatrz(x, y) || odwiedzone.wewnatrz(x, y) {
				continue
			}
			nalezy := nowaMaskaRastrowa(szerokosc, wysokosc)
			minX, minY, maxX, maxY := x, y, x, y
			pole := 0
			stos := [][2]int{{x, y}}
			odwiedzone.ustaw(x, y, true)
			for len(stos) > 0 {
				punkt := stos[len(stos)-1]
				stos = stos[:len(stos)-1]
				nalezy.ustaw(punkt[0], punkt[1], true)
				pole++
				if punkt[0] < minX {
					minX = punkt[0]
				}
				if punkt[0] > maxX {
					maxX = punkt[0]
				}
				if punkt[1] < minY {
					minY = punkt[1]
				}
				if punkt[1] > maxY {
					maxY = punkt[1]
				}
				for _, kierunek := range kierunkiObrysu {
					sasiadX := punkt[0] + kierunek[0]
					sasiadY := punkt[1] + kierunek[1]
					if plan.wewnatrz(sasiadX, sasiadY) && !odwiedzone.wewnatrz(sasiadX, sasiadY) {
						odwiedzone.ustaw(sasiadX, sasiadY, true)
						stos = append(stos, [2]int{sasiadX, sasiadY})
					}
				}
			}
			if pole < najmniejszaWarstwa {
				continue
			}
			obszary = append(obszary, obszarWarstwy{
				prostokat: image.Rect(minX, minY, maxX+1, maxY+1),
				nalezy:    nalezy,
				pole:      pole,
			})
		}
	}

	// Od największego: gdy granica warstw utnie wykaz, ma uciąć drobiazgi,
	// a nie główny obiekt kadru.
	for i := 1; i < len(obszary); i++ {
		for j := i; j > 0 && obszary[j].pole > obszary[j-1].pole; j-- {
			obszary[j], obszary[j-1] = obszary[j-1], obszary[j]
		}
	}
	return obszary
}

// wytnijWarstwe zapisuje jeden obiekt jako PNG przycięty do jego prostokąta.
// Piksele spoza obszaru zostają przezroczyste, także gdy leżą w prostokącie —
// dwa obiekty o zachodzących prostokątach nie mają prawa nieść siebie nawzajem.
func wytnijWarstwe(obraz image.Image, obszar obszarWarstwy) ([]byte, error) {
	granice := obraz.Bounds()
	warstwa := image.NewNRGBA(image.Rect(0, 0, obszar.prostokat.Dx(), obszar.prostokat.Dy()))
	for y := obszar.prostokat.Min.Y; y < obszar.prostokat.Max.Y; y++ {
		for x := obszar.prostokat.Min.X; x < obszar.prostokat.Max.X; x++ {
			if !obszar.nalezy.wewnatrz(x, y) {
				continue
			}
			r, g, b, a := skladoweNiePrzemnozone(obraz.At(granice.Min.X+x, granice.Min.Y+y))
			warstwa.Set(x-obszar.prostokat.Min.X, y-obszar.prostokat.Min.Y,
				kolorZeSkladowych(r, g, b, a))
		}
	}
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, warstwa); err != nil {
		return nil, bladPrzetwarzaniaModeluObrazu("nie można zapisać warstwy: " + err.Error())
	}
	return bufor.Bytes(), nil
}

// nazwaWarstwy składa czytelną nazwę zasobu warstwy. Trzymana osobno, bo tę
// samą nazwę widzi Assets Panel i wykaz zasobów okna.
func nazwaWarstwy(zrodlo string, numer int) string {
	podstawa := strings.TrimSpace(zrodlo)
	if podstawa == "" {
		podstawa = "obraz"
	}
	return podstawa + " — warstwa " + strconv.Itoa(numer)
}
