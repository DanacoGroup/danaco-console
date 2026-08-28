// Plik obsługuje rachunek na pikselach warsztatu fotografii modułu design: kadr,
// przekształcenia, rozdzielczość, korekcje barwne, filtry, retusz, maski,
// kompozycja warstw i obrysowanie konturów.
package core

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/draw"

	"danacoconsole/shared"
)

const (
	// granicaBokuFotografiiDesignu chroni rachunek przed obrazem o boku
	// stutysięcznym: każda czynność idzie po całej powierzchni.
	granicaBokuFotografiiDesignu = 20000

	// granicaPikseliFotografiiDesignu jest granicą powierzchni WYNIKU. Obraz
	// 12000×12000 to 144 miliony punktów i cztery bajty na punkt — pół gigabajta
	// na samo płótno.
	granicaPikseliFotografiiDesignu = 160_000_000

	// promienOdszumianiaDesignu jest odchyleniem PRZESTRZENNYM filtru
	// odszumiającego przy pełnej sile — zasięgiem w pikselach, na którym punkt
	// bierze pod uwagę sąsiadów.
	promienOdszumianiaDesignu = 3.0

	// zakresOdszumianiaDesignu jest odchyleniem ZAKRESOWYM filtru odszumiającego
	// przy pełnej sile, w skali 0–255 — ta liczba czyni filtr bilateralnym, a nie
	// rozmyciem. Trzydzieści dwa stopnie dzieli szum od treści obrazu.
	zakresOdszumianiaDesignu = 32.0

	// granicaPromieniaOdszumianiaDesignu ogranicza zasięg jądra. Bez niej siła 1
	// przy dużym odchyleniu dawałaby jądro kilkudziesięciopunktowe na każdy punkt
	// obrazu, a odszumienie ma być czynnością, nie czekaniem.
	granicaPromieniaOdszumianiaDesignu = 8

	// promienWyostrzeniaDesignu jest promieniem maski wyostrzającej w pikselach,
	// mnożonym przez żądaną siłę wyostrzenia.
	promienWyostrzeniaDesignu = 1.6
)

// przytnijObrazFotografiiDesignu wycina kadr, prostuje horyzont i sprowadza do
// proporcji. Kolejność jest zamierzona: najpierw prostowanie, potem kadr, żeby
// puste naroża po obrocie zdążył usunąć kadr, a nie zostały w wyniku.
func przytnijObrazFotografiiDesignu(obraz image.Image,
	z shared.DesignPhotoCropRequest) (image.Image, error) {

	wynik := obraz
	if z.StraightenDeg != nil && *z.StraightenDeg != 0 {
		// Obrót z wypełnieniem przezroczystością, a nie czernią — puste naroża mają być widocznym brakiem.
		wynik = imaging.Rotate(wynik, -*z.StraightenDeg, color.NRGBA{})
		// Kadr wpisany w obrót: prostokąt największy, który po obrocie nie zawiera pustych naroży.
		wynik = wpiszKadrWObrotDesignu(wynik, obraz, math.Abs(*z.StraightenDeg))
	}

	granice := wynik.Bounds()
	kadr := granice
	if z.Width != nil && z.Height != nil && *z.Width > 0 && *z.Height > 0 {
		x, y := 0.0, 0.0
		if z.X != nil {
			x = *z.X
		}
		if z.Y != nil {
			y = *z.Y
		}
		kadr = image.Rect(
			granice.Min.X+int(x), granice.Min.Y+int(y),
			granice.Min.X+int(x+*z.Width), granice.Min.Y+int(y+*z.Height),
		).Intersect(granice)
		if kadr.Empty() {
			return nil, fmt.Errorf(
				"kadr %v;%v %v×%v nie ma części wspólnej z obrazem %d×%d — kadr poza obrazem nie "+
					"jest kadrem", x, y, *z.Width, *z.Height, granice.Dx(), granice.Dy())
		}
	}
	if z.AspectRatio != nil && strings.TrimSpace(*z.AspectRatio) != "" {
		proporcje, err := proporcjeKadruDesignu(*z.AspectRatio)
		if err != nil {
			return nil, err
		}
		kadr = wpiszProporcjeWKadrDesignu(kadr, proporcje)
	}
	if kadr.Dx() < 1 || kadr.Dy() < 1 {
		return nil, fmt.Errorf("kadr wyszedł o boku %d×%d — kadr o zerowym boku nie istnieje",
			kadr.Dx(), kadr.Dy())
	}
	return imaging.Crop(wynik, kadr), nil
}

// wpiszKadrWObrotDesignu wycina z obrazu obróconego prostokąt bez pustych
// naroży, zachowując proporcje źródła: prostokąt o tych proporcjach wpisany
// w obrót o kąt fi skaluje się współczynnikiem, którego mianownik jest sumą
// rzutów boków.
func wpiszKadrWObrotDesignu(obrocony, zrodlo image.Image, katStopni float64) image.Image {
	if katStopni <= 0 {
		return obrocony
	}
	kat := math.Mod(katStopni, 180) * math.Pi / 180
	if kat > math.Pi/2 {
		kat = math.Pi - kat
	}
	szerokoscZrodla := float64(zrodlo.Bounds().Dx())
	wysokoscZrodla := float64(zrodlo.Bounds().Dy())
	if szerokoscZrodla <= 0 || wysokoscZrodla <= 0 {
		return obrocony
	}
	dluzszy, krotszy := szerokoscZrodla, wysokoscZrodla
	if krotszy > dluzszy {
		dluzszy, krotszy = krotszy, dluzszy
	}
	sin, cos := math.Sin(kat), math.Cos(kat)
	// Dwa przypadki wzoru: kąt „mały" (kadr ogranicza bok krótszy) i kąt „duży".
	var szerokoscKadru, wysokoscKadru float64
	if krotszy <= 2*sin*cos*dluzszy {
		polowa := 0.5 * krotszy
		if szerokoscZrodla >= wysokoscZrodla {
			szerokoscKadru, wysokoscKadru = polowa/sin, polowa/cos
		} else {
			szerokoscKadru, wysokoscKadru = polowa/cos, polowa/sin
		}
	} else {
		mianownik := cos*cos - sin*sin
		szerokoscKadru = (szerokoscZrodla*cos - wysokoscZrodla*sin) / mianownik
		wysokoscKadru = (wysokoscZrodla*cos - szerokoscZrodla*sin) / mianownik
	}
	if szerokoscKadru <= 1 || wysokoscKadru <= 1 {
		return obrocony
	}
	granice := obrocony.Bounds()
	if int(szerokoscKadru) >= granice.Dx() && int(wysokoscKadru) >= granice.Dy() {
		return obrocony
	}
	return imaging.CropCenter(obrocony, int(szerokoscKadru), int(wysokoscKadru))
}

// proporcjeKadruDesignu rozkłada zapis proporcji („16:9") na iloraz szerokości
// do wysokości, sprawdzając poprawność obu liczb.
func proporcjeKadruDesignu(zapis string) (float64, error) {
	czesci := strings.Split(strings.TrimSpace(zapis), ":")
	if len(czesci) != 2 {
		return 0, fmt.Errorf(
			"proporcji %q nie da się odczytać — rdzeń przyjmuje zapis „szerokość:wysokość\", "+
				"na przykład 1:1, 4:5 albo 16:9", zapis)
	}
	szerokosc, err := strconv.ParseFloat(strings.TrimSpace(czesci[0]), 64)
	if err != nil || szerokosc <= 0 {
		return 0, fmt.Errorf("proporcje %q mają niepoprawną szerokość", zapis)
	}
	wysokosc, err := strconv.ParseFloat(strings.TrimSpace(czesci[1]), 64)
	if err != nil || wysokosc <= 0 {
		return 0, fmt.Errorf("proporcje %q mają niepoprawną wysokość", zapis)
	}
	return szerokosc / wysokosc, nil
}

// wpiszProporcjeWKadrDesignu zawęża kadr do żądanych proporcji, licząc OD ŚRODKA
// — tak, żeby żaden brzeg nie został uprzywilejowany bez powodu.
func wpiszProporcjeWKadrDesignu(kadr image.Rectangle, proporcje float64) image.Rectangle {
	biezace := float64(kadr.Dx()) / float64(kadr.Dy())
	if math.Abs(biezace-proporcje) < 1e-9 {
		return kadr
	}
	if biezace > proporcje {
		nowaSzerokosc := int(float64(kadr.Dy())*proporcje + 0.5)
		odsuniecie := (kadr.Dx() - nowaSzerokosc) / 2
		return image.Rect(kadr.Min.X+odsuniecie, kadr.Min.Y,
			kadr.Min.X+odsuniecie+nowaSzerokosc, kadr.Max.Y)
	}
	nowaWysokosc := int(float64(kadr.Dx())/proporcje + 0.5)
	odsuniecie := (kadr.Dy() - nowaWysokosc) / 2
	return image.Rect(kadr.Min.X, kadr.Min.Y+odsuniecie,
		kadr.Max.X, kadr.Min.Y+odsuniecie+nowaWysokosc)
}

// przeksztalcObrazFotografiiDesignu obraca, odbija i koryguje perspektywę oraz
// dystorsję obiektywu żądania, w tej kolejności zapisanej przez wołającego.
func przeksztalcObrazFotografiiDesignu(obraz image.Image,
	z shared.DesignPhotoTransformRequest) (image.Image, error) {

	wynik := obraz
	if z.RotateDeg != nil && *z.RotateDeg != 0 {
		wynik = imaging.Rotate(wynik, -*z.RotateDeg, color.NRGBA{})
	}
	if z.FlipHorizontal != nil && *z.FlipHorizontal {
		wynik = imaging.FlipH(wynik)
	}
	if z.FlipVertical != nil && *z.FlipVertical {
		wynik = imaging.FlipV(wynik)
	}
	if z.LensDistortion != nil && *z.LensDistortion != 0 {
		wynik = skorygujDystorsjeDesignu(wynik, *z.LensDistortion)
	}
	if len(z.Perspective) > 0 {
		if len(z.Perspective) != 4 {
			return nil, fmt.Errorf(
				"korekcja perspektywy wymaga DOKŁADNIE czterech naroży (podano %d) — trzy punkty "+
					"nie wyznaczają czworokąta, a pięć wyznacza dwa różne", len(z.Perspective))
		}
		poprawiony, err := skorygujPerspektyweDesignu(wynik, z.Perspective)
		if err != nil {
			return nil, err
		}
		wynik = poprawiony
	}
	return wynik, nil
}

// skorygujDystorsjeDesignu prostuje dystorsję obiektywu przekształceniem
// promieniowym. Rachunek jest odwrotnym odwzorowaniem: dla każdego punktu
// wyniku liczy się, skąd w źródle go wziąć, żeby wynik nie miał dziur.
func skorygujDystorsjeDesignu(obraz image.Image, sila float64) image.Image {
	granice := obraz.Bounds()
	szerokosc, wysokosc := granice.Dx(), granice.Dy()
	wynik := image.NewNRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	srodekX, srodekY := float64(szerokosc)/2, float64(wysokosc)/2
	promienNormujacy := math.Hypot(srodekX, srodekY)
	if promienNormujacy == 0 {
		return obraz
	}
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			dx := (float64(x) - srodekX) / promienNormujacy
			dy := (float64(y) - srodekY) / promienNormujacy
			promien := math.Hypot(dx, dy)
			// Wielomian promieniowy pierwszego stopnia wystarcza na dystorsję beczkową i poduszkową obiektywów.
			wspolczynnik := 1 + sila*promien*promien
			zrodloweX := srodekX + dx*promienNormujacy*wspolczynnik
			zrodloweY := srodekY + dy*promienNormujacy*wspolczynnik
			wynik.Set(x, y, punktDwuliniowoDesignu(obraz, zrodloweX, zrodloweY))
		}
	}
	return wynik
}

// punktDwuliniowoDesignu próbkuje obraz dwuliniowo. Punkt poza obrazem jest
// PRZEZROCZYSTY, nie czarny: brak treści ma wyglądać na brak treści.
func punktDwuliniowoDesignu(obraz image.Image, x, y float64) color.NRGBA {
	granice := obraz.Bounds()
	if x < float64(granice.Min.X)-1 || y < float64(granice.Min.Y)-1 ||
		x > float64(granice.Max.X) || y > float64(granice.Max.Y) {
		return color.NRGBA{}
	}
	x0, y0 := int(math.Floor(x)), int(math.Floor(y))
	udzialX, udzialY := x-float64(x0), y-float64(y0)
	skladowe := func(px, py int) (float64, float64, float64, float64) {
		if px < granice.Min.X || py < granice.Min.Y || px >= granice.Max.X || py >= granice.Max.Y {
			return 0, 0, 0, 0
		}
		r, g, b, a := obraz.At(px, py).RGBA()
		return float64(r), float64(g), float64(b), float64(a)
	}
	r00, g00, b00, a00 := skladowe(x0, y0)
	r10, g10, b10, a10 := skladowe(x0+1, y0)
	r01, g01, b01, a01 := skladowe(x0, y0+1)
	r11, g11, b11, a11 := skladowe(x0+1, y0+1)
	mieszaj := func(v00, v10, v01, v11 float64) float64 {
		gora := v00*(1-udzialX) + v10*udzialX
		dol := v01*(1-udzialX) + v11*udzialX
		return gora*(1-udzialY) + dol*udzialY
	}
	alfa := mieszaj(a00, a10, a01, a11)
	if alfa <= 0 {
		return color.NRGBA{}
	}
	// Składowe z RGBA() są przemnożone przez krycie, a NRGBA ich nie mnoży — stąd dzielenie przez krycie.
	return color.NRGBA{
		R: uint8(mieszaj(r00, r10, r01, r11) / alfa * 255),
		G: uint8(mieszaj(g00, g10, g01, g11) / alfa * 255),
		B: uint8(mieszaj(b00, b10, b01, b11) / alfa * 255),
		A: uint8(alfa / 65535 * 255),
	}
}

// skorygujPerspektyweDesignu wyprostowuje czworokąt wskazany narożami do
// prostokąta, homografią liczoną z czterech par punktów metodą eliminacji
// Gaussa — przekształcenie afiniczne trzech par tu nie wystarcza, bo zachowuje
// równoległość.
func skorygujPerspektyweDesignu(obraz image.Image,
	naroza []shared.DesignPhotoPoint) (image.Image, error) {

	granice := obraz.Bounds()
	szerokosc, wysokosc := float64(granice.Dx()), float64(granice.Dy())
	// Docelowy prostokąt to płótno, źródło to czworokąt — rachunek idzie od docelowego do źródłowego.
	docelowe := [4][2]float64{{0, 0}, {szerokosc, 0}, {szerokosc, wysokosc}, {0, wysokosc}}
	zrodlowe := [4][2]float64{}
	for numer, naroze := range naroza {
		zrodlowe[numer] = [2]float64{naroze.X, naroze.Y}
	}
	macierz, err := homografiaDesignu(docelowe, zrodlowe)
	if err != nil {
		return nil, err
	}

	wynik := image.NewNRGBA(image.Rect(0, 0, granice.Dx(), granice.Dy()))
	for y := 0; y < granice.Dy(); y++ {
		for x := 0; x < granice.Dx(); x++ {
			px, py := float64(x)+0.5, float64(y)+0.5
			mianownik := macierz[6]*px + macierz[7]*py + 1
			if mianownik == 0 {
				continue
			}
			zrodloweX := (macierz[0]*px + macierz[1]*py + macierz[2]) / mianownik
			zrodloweY := (macierz[3]*px + macierz[4]*py + macierz[5]) / mianownik
			wynik.Set(x, y, punktDwuliniowoDesignu(obraz,
				float64(granice.Min.X)+zrodloweX, float64(granice.Min.Y)+zrodloweY))
		}
	}
	return wynik, nil
}

// homografiaDesignu liczy osiem współczynników homografii z czterech par
// punktów. Dziewiąty współczynnik jest ustalony na jedynkę — homografia jest
// wyznaczona z dokładnością do stałej.
func homografiaDesignu(od, do [4][2]float64) ([8]float64, error) {
	var uklad [8][9]float64
	for numer := 0; numer < 4; numer++ {
		x, y := od[numer][0], od[numer][1]
		u, v := do[numer][0], do[numer][1]
		uklad[2*numer] = [9]float64{x, y, 1, 0, 0, 0, -u * x, -u * y, u}
		uklad[2*numer+1] = [9]float64{0, 0, 0, x, y, 1, -v * x, -v * y, v}
	}
	// Eliminacja Gaussa z wyborem elementu głównego: bez wyboru zero na przekątnej dzieliłoby przez zero.
	for kolumna := 0; kolumna < 8; kolumna++ {
		glowny := kolumna
		for wiersz := kolumna + 1; wiersz < 8; wiersz++ {
			if math.Abs(uklad[wiersz][kolumna]) > math.Abs(uklad[glowny][kolumna]) {
				glowny = wiersz
			}
		}
		if math.Abs(uklad[glowny][kolumna]) < 1e-12 {
			return [8]float64{}, fmt.Errorf(
				"z podanych naroży nie da się policzyć korekcji perspektywy — trzy z nich leżą " +
					"na jednej prostej albo dwa się pokrywają, a taki czworokąt nie ma wnętrza")
		}
		uklad[kolumna], uklad[glowny] = uklad[glowny], uklad[kolumna]
		for wiersz := 0; wiersz < 8; wiersz++ {
			if wiersz == kolumna {
				continue
			}
			wspolczynnik := uklad[wiersz][kolumna] / uklad[kolumna][kolumna]
			for pole := kolumna; pole < 9; pole++ {
				uklad[wiersz][pole] -= wspolczynnik * uklad[kolumna][pole]
			}
		}
	}
	var wynik [8]float64
	for numer := 0; numer < 8; numer++ {
		wynik[numer] = uklad[numer][8] / uklad[numer][numer]
	}
	return wynik, nil
}

// przeliczRozdzielczoscFotografiiDesignu przelicza rozdzielczość obrazu
// wskazanym filtrem, zachowując proporcje, gdy żądanie o to prosi.
func przeliczRozdzielczoscFotografiiDesignu(obraz image.Image, szerokosc, wysokosc int,
	filtr shared.DesignPhotoResampleFilter, zachowajProporcje bool) (image.Image, error) {

	granice := obraz.Bounds()
	if szerokosc <= 0 && wysokosc <= 0 {
		return nil, fmt.Errorf(
			"żądanie nie podało ani szerokości, ani wysokości docelowej — rdzeń nie zgaduje, " +
				"na jaki rozmiar przeliczyć")
	}
	if zachowajProporcje {
		// Zero w jednym z wymiarów znaczy dla biblioteki „dobierz z proporcji" — o to prosi keepAspectRatio.
		if szerokosc > 0 && wysokosc > 0 {
			skalaX := float64(szerokosc) / float64(granice.Dx())
			skalaY := float64(wysokosc) / float64(granice.Dy())
			skala := mniejszaDesignu(skalaX, skalaY)
			szerokosc = int(float64(granice.Dx())*skala + 0.5)
			wysokosc = int(float64(granice.Dy())*skala + 0.5)
		}
	} else if szerokosc <= 0 || wysokosc <= 0 {
		// Bez zachowania proporcji brakujący wymiar zostaje niezmieniony — jedyna liczba pewna bez zgadywania.
		if szerokosc <= 0 {
			szerokosc = granice.Dx()
		}
		if wysokosc <= 0 {
			wysokosc = granice.Dy()
		}
	}
	if err := sprawdzRozmiarFotografiiDesignu(szerokosc, wysokosc); err != nil {
		return nil, err
	}
	return imaging.Resize(obraz, szerokosc, wysokosc, filtrPrzeliczeniaDesignu(filtr)), nil
}

// filtrPrzeliczeniaDesignu przekłada filtr kontraktu na filtr biblioteki,
// biorąc Lanczosa jako wartość domyślną.
func filtrPrzeliczeniaDesignu(filtr shared.DesignPhotoResampleFilter) imaging.ResampleFilter {
	switch filtr {
	case shared.DesignPhotoResampleFilterBilinear:
		return imaging.Linear
	case shared.DesignPhotoResampleFilterNearest:
		return imaging.NearestNeighbor
	}
	return imaging.Lanczos
}

// sprawdzRozmiarFotografiiDesignu odrzuca rozmiar nie do zmieszczenia w pamięci
// przed rachunkiem, sprawdzając bok i powierzchnię wyniku.
func sprawdzRozmiarFotografiiDesignu(szerokosc, wysokosc int) error {
	if szerokosc < 1 || wysokosc < 1 {
		return fmt.Errorf("rozmiar %d×%d: obraz o niedodatnim boku nie istnieje", szerokosc, wysokosc)
	}
	if szerokosc > granicaBokuFotografiiDesignu || wysokosc > granicaBokuFotografiiDesignu {
		return fmt.Errorf("rozmiar %d×%d przekracza granicę boku %d", szerokosc, wysokosc,
			granicaBokuFotografiiDesignu)
	}
	if szerokosc*wysokosc > granicaPikseliFotografiiDesignu {
		return fmt.Errorf(
			"rozmiar %d×%d to %d punktów, a granica wyniku to %d — rdzeń nie zamawia pamięci, "+
				"której nie dostanie", szerokosc, wysokosc, szerokosc*wysokosc,
			granicaPikseliFotografiiDesignu)
	}
	return nil
}

// popraweJakoscFotografiiDesignu wykonuje auto-poziomy, auto-kontrast,
// odszumienie i wyostrzenie, oddając obraz wraz z wykazem kroków, które
// naprawdę weszły. Żądanie z samymi wartościami zerowymi dostaje obraz
// nietknięty i pusty wykaz.
func popraweJakoscFotografiiDesignu(obraz image.Image,
	z shared.DesignPhotoEnhanceRequest) (image.Image, []string) {

	wynik := obraz
	kroki := []string{}
	// Brak wskazania czegokolwiek znaczy „popraw rozsądnie": auto-poziomy i delikatne wyostrzenie.
	bezNastaw := z.AutoLevels == nil && z.AutoContrast == nil && z.Denoise == nil && z.Sharpen == nil

	if bezNastaw || (z.AutoLevels != nil && *z.AutoLevels) {
		wynik = rozciagnijHistogramDesignu(wynik)
		kroki = append(kroki, "auto-poziomy")
	}
	if z.AutoContrast != nil && *z.AutoContrast {
		// Krzywa sigmoidalna zamiast liniowego kontrastu dociąga środek, zostawiając światła i cienie.
		wynik = imaging.AdjustSigmoid(wynik, 0.5, 3)
		kroki = append(kroki, "auto-kontrast")
	}
	if z.Denoise != nil && *z.Denoise > 0 {
		sila := przytnijUlamekDesignu(*z.Denoise)
		wynik = odszumBilateralnieDesignu(wynik, sila)
		kroki = append(kroki, fmt.Sprintf("odszumienie %.2f", sila))
	}
	sila := 0.0
	if bezNastaw {
		sila = 0.4
	}
	if z.Sharpen != nil && *z.Sharpen > 0 {
		sila = przytnijUlamekDesignu(*z.Sharpen)
	}
	if sila > 0 {
		wynik = imaging.Sharpen(wynik, promienWyostrzeniaDesignu*sila)
		kroki = append(kroki, fmt.Sprintf("wyostrzenie %.2f", sila))
	}
	return wynik, kroki
}

// odszumBilateralnieDesignu odszumia obraz filtrem bilateralnym — średnią
// ważoną odległością sąsiada i różnicą jego jasności, żeby szum znikał,
// a krawędzie zostawały ostre. Filtr idzie osobno w poziomie i w pionie.
func odszumBilateralnieDesignu(obraz image.Image, sila float64) *image.NRGBA {
	odchyleniePrzestrzenne := promienOdszumianiaDesignu * sila
	odchylenieZakresu := zakresOdszumianiaDesignu * sila
	if odchyleniePrzestrzenne <= 0 || odchylenieZakresu <= 0 {
		return imaging.Clone(obraz)
	}
	promien := int(math.Ceil(2 * odchyleniePrzestrzenne))
	if promien < 1 {
		promien = 1
	}
	if promien > granicaPromieniaOdszumianiaDesignu {
		promien = granicaPromieniaOdszumianiaDesignu
	}

	// Wagi liczą się raz, do tablic — liczenie wykładnika w pętli byłoby setką milionów wywołań.
	wagiOdleglosci := make([]float64, promien+1)
	for odleglosc := 0; odleglosc <= promien; odleglosc++ {
		wagiOdleglosci[odleglosc] = math.Exp(-float64(odleglosc*odleglosc) /
			(2 * odchyleniePrzestrzenne * odchyleniePrzestrzenne))
	}
	var wagiJasnosci [256]float64
	for roznica := 0; roznica < 256; roznica++ {
		wagiJasnosci[roznica] = math.Exp(-float64(roznica*roznica) /
			(2 * odchylenieZakresu * odchylenieZakresu))
	}

	poziome := przejscieBilateralneDesignu(imaging.Clone(obraz), wagiOdleglosci, &wagiJasnosci, true)
	return przejscieBilateralneDesignu(poziome, wagiOdleglosci, &wagiJasnosci, false)
}

// przejscieBilateralneDesignu wykonuje jedno przejście filtru — w poziomie
// albo w pionie, ważąc sąsiadów tablicami wag policzonymi wcześniej.
func przejscieBilateralneDesignu(zrodlo *image.NRGBA, wagiOdleglosci []float64,
	wagiJasnosci *[256]float64, poziomo bool) *image.NRGBA {

	granice := zrodlo.Bounds()
	wynik := image.NewNRGBA(granice)
	promien := len(wagiOdleglosci) - 1
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			srodek := zrodlo.NRGBAAt(x, y)
			jasnoscSrodka := int(jasnoscNrgbaDesignu(srodek))
			sumaWag, sumaR, sumaG, sumaB := 0.0, 0.0, 0.0, 0.0
			for krok := -promien; krok <= promien; krok++ {
				sasiadX, sasiadY := x, y
				if poziomo {
					sasiadX += krok
				} else {
					sasiadY += krok
				}
				// Punkt spoza obrazu nie wchodzi do średniej: wypełnienie brzegu czernią zniekształciłoby wynik.
				if sasiadX < granice.Min.X || sasiadY < granice.Min.Y ||
					sasiadX >= granice.Max.X || sasiadY >= granice.Max.Y {
					continue
				}
				sasiad := zrodlo.NRGBAAt(sasiadX, sasiadY)
				roznica := int(jasnoscNrgbaDesignu(sasiad)) - jasnoscSrodka
				if roznica < 0 {
					roznica = -roznica
				}
				odleglosc := krok
				if odleglosc < 0 {
					odleglosc = -odleglosc
				}
				waga := wagiOdleglosci[odleglosc] * wagiJasnosci[roznica]
				sumaWag += waga
				sumaR += waga * float64(sasiad.R)
				sumaG += waga * float64(sasiad.G)
				sumaB += waga * float64(sasiad.B)
			}
			if sumaWag <= 0 {
				wynik.SetNRGBA(x, y, srodek)
				continue
			}
			wynik.SetNRGBA(x, y, color.NRGBA{
				R: przytnijSkladowaDesignu(sumaR / sumaWag),
				G: przytnijSkladowaDesignu(sumaG / sumaWag),
				B: przytnijSkladowaDesignu(sumaB / sumaWag),
				A: srodek.A,
			})
		}
	}
	return wynik
}

// jasnoscNrgbaDesignu liczy jasność punktu wagami postrzegania barw. Ta sama
// trójka wag stoi w auto-poziomach i w odczycie maski — jasność ma w całym
// warsztacie jedną definicję.
func jasnoscNrgbaDesignu(punkt color.NRGBA) uint8 {
	return uint8((int(punkt.R)*299 + int(punkt.G)*587 + int(punkt.B)*114) / 1000)
}

// rozciagnijHistogramDesignu rozciąga histogram na pełną skalę — auto-poziomy.
//
// Krańce bierze się z ODCIĘCIEM po pół procenta z każdej strony: pojedynczy
// przepalony punkt albo jeden punkt czerni rozciągnąłby skalę do niczego, a
// odcięcie ich pomija.
func rozciagnijHistogramDesignu(obraz image.Image) image.Image {
	granice := obraz.Bounds()
	histogram := [256]int{}
	punktow := 0
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			r, g, b, _ := obraz.At(x, y).RGBA()
			jasnosc := (int(r>>8)*299 + int(g>>8)*587 + int(b>>8)*114) / 1000
			histogram[jasnosc]++
			punktow++
		}
	}
	if punktow == 0 {
		return obraz
	}
	odciecie := punktow / 200
	dolny, gorny := 0, 255
	suma := 0
	for numer := 0; numer < 256; numer++ {
		suma += histogram[numer]
		if suma > odciecie {
			dolny = numer
			break
		}
	}
	suma = 0
	for numer := 255; numer >= 0; numer-- {
		suma += histogram[numer]
		if suma > odciecie {
			gorny = numer
			break
		}
	}
	if gorny-dolny < 8 {
		// Histogram węższy niż osiem poziomów jest niemal jednolity — rozciągnięcie dałoby szum.
		return obraz
	}
	skala := 255.0 / float64(gorny-dolny)
	return imaging.AdjustFunc(obraz, func(punkt color.NRGBA) color.NRGBA {
		przelicz := func(skladowa uint8) uint8 {
			wartosc := (float64(skladowa) - float64(dolny)) * skala
			if wartosc < 0 {
				wartosc = 0
			}
			if wartosc > 255 {
				wartosc = 255
			}
			return uint8(wartosc)
		}
		return color.NRGBA{R: przelicz(punkt.R), G: przelicz(punkt.G), B: przelicz(punkt.B),
			A: punkt.A}
	})
}

// skorygujBarweFotografiiDesignu wykonuje korekcje barwne żądania: ekspozycję,
// jasność, kontrast, gammę, nasycenie, temperaturę, odcień i krzywe tonalne.
func skorygujBarweFotografiiDesignu(obraz image.Image,
	z shared.DesignPhotoColorCorrectRequest) (image.Image, error) {

	wynik := obraz
	if z.Exposure != nil && *z.Exposure != 0 {
		// Ekspozycja liczy się w działkach: jedna działka to podwojenie światła, mnożnik jest potęgą dwójki.
		mnoznik := math.Pow(2, *z.Exposure)
		wynik = imaging.AdjustFunc(wynik, func(punkt color.NRGBA) color.NRGBA {
			return color.NRGBA{
				R: przytnijSkladowaDesignu(float64(punkt.R) * mnoznik),
				G: przytnijSkladowaDesignu(float64(punkt.G) * mnoznik),
				B: przytnijSkladowaDesignu(float64(punkt.B) * mnoznik),
				A: punkt.A,
			}
		})
	}
	if z.Brightness != nil && *z.Brightness != 0 {
		wynik = imaging.AdjustBrightness(wynik, przytnijZakresJednosciDesignu(*z.Brightness)*100)
	}
	if z.Contrast != nil && *z.Contrast != 0 {
		wynik = imaging.AdjustContrast(wynik, przytnijZakresJednosciDesignu(*z.Contrast)*100)
	}
	if z.Gamma != nil && *z.Gamma > 0 && *z.Gamma != 1 {
		wynik = imaging.AdjustGamma(wynik, *z.Gamma)
	}
	if z.Saturation != nil && *z.Saturation != 0 {
		wynik = imaging.AdjustSaturation(wynik, przytnijZakresJednosciDesignu(*z.Saturation)*100)
	}
	if z.Temperature != nil && *z.Temperature != 0 {
		wynik = zmienTemperatureBarwowaDesignu(wynik, przytnijZakresJednosciDesignu(*z.Temperature))
	}
	if z.HueShiftDeg != nil && *z.HueShiftDeg != 0 {
		wynik = przesunOdcienDesignu(wynik, *z.HueShiftDeg)
	}
	if z.Lightness != nil && *z.Lightness != 0 {
		wynik = zmienJasnoscPercepcyjnaDesignu(wynik, przytnijZakresJednosciDesignu(*z.Lightness))
	}
	if len(z.Curves) > 0 {
		poprawiony, err := zastosujKrzyweDesignu(wynik, z.Curves)
		if err != nil {
			return nil, err
		}
		wynik = poprawiony
	}
	return wynik, nil
}

// zmienTemperatureBarwowaDesignu przesuwa barwę w stronę ciepłą albo chłodną,
// przesunięciem kanałów czerwonego i niebieskiego w przeciwnych kierunkach —
// tak działa suwak temperatury: cieplej znaczy więcej czerwieni, zieleń
// zostaje odniesieniem.
func zmienTemperatureBarwowaDesignu(obraz image.Image, sila float64) image.Image {
	przesuniecie := sila * 40
	return imaging.AdjustFunc(obraz, func(punkt color.NRGBA) color.NRGBA {
		return color.NRGBA{
			R: przytnijSkladowaDesignu(float64(punkt.R) + przesuniecie),
			G: punkt.G,
			B: przytnijSkladowaDesignu(float64(punkt.B) - przesuniecie),
			A: punkt.A,
		}
	})
}

// przesunOdcienDesignu obraca odcień wszystkich punktów obrazu o wskazaną
// liczbę stopni w przestrzeni HSL.
func przesunOdcienDesignu(obraz image.Image, stopnie float64) image.Image {
	return imaging.AdjustFunc(obraz, func(punkt color.NRGBA) color.NRGBA {
		barwa := kolorZNrgbaDesignu(punkt)
		odcien, nasycenie, jasnosc := barwa.Hsl()
		nowa := kolorHslDesignu(znormalizujKatDesignu(odcien+stopnie), nasycenie, jasnosc)
		return nrgbaZKoloruDesignu(nowa, punkt.A)
	})
}

// zmienJasnoscPercepcyjnaDesignu zmienia jasność w przestrzeni HCL.
//
// HCL, nie RGB: dodanie stałej do składowych rozjaśnia barwy nasycone mniej niż
// szarości i przekłamuje odcień, a zmiana `L` w HCL rozjaśnia wszystko tak, jak
// widzi to oko.
func zmienJasnoscPercepcyjnaDesignu(obraz image.Image, sila float64) image.Image {
	return imaging.AdjustFunc(obraz, func(punkt color.NRGBA) color.NRGBA {
		barwa := kolorZNrgbaDesignu(punkt)
		odcien, nasycenie, jasnosc := barwa.Hcl()
		nowa := kolorHclDesignu(odcien, nasycenie, przytnijUlamekDesignu(jasnosc+sila*0.4))
		return nrgbaZKoloruDesignu(nowa, punkt.A)
	})
}

// zastosujKrzyweDesignu nakłada krzywe tonalne na kanały. Zapis krzywej jest
// wykazem punktów `{wejscie, wyjscie}` w skali 0–255 dla kanałów `rgb`, `r`,
// `g`, `b`. Między punktami wartość idzie liniowo.
func zastosujKrzyweDesignu(obraz image.Image, zapis []byte) (image.Image, error) {
	tablice, err := tabliceKrzywychDesignu(zapis)
	if err != nil {
		return nil, err
	}
	return imaging.AdjustFunc(obraz, func(punkt color.NRGBA) color.NRGBA {
		return color.NRGBA{
			R: tablice.wspolna[tablice.czerwony[punkt.R]],
			G: tablice.wspolna[tablice.zielony[punkt.G]],
			B: tablice.wspolna[tablice.niebieski[punkt.B]],
			A: punkt.A,
		}
	}), nil
}

// przytnijSkladowaDesignu wprowadza składową barwy w zakres bajtu, zaokrąglając
// do najbliższej liczby całkowitej.
func przytnijSkladowaDesignu(wartosc float64) uint8 {
	if wartosc < 0 {
		return 0
	}
	if wartosc > 255 {
		return 255
	}
	return uint8(wartosc + 0.5)
}

// przytnijZakresJednosciDesignu wprowadza wartość w zakres od -1 do 1, w którym
// żądanie zapisuje siłę korekcji.
func przytnijZakresJednosciDesignu(wartosc float64) float64 {
	if wartosc < -1 {
		return -1
	}
	if wartosc > 1 {
		return 1
	}
	return wartosc
}

// zlozWarstwyFotografiiDesignu składa warstwy rastrowe trybami mieszania,
// skalując warstwę do obszaru docelowego, gdy wymiary się różnią.
func zlozWarstwyFotografiiDesignu(plotno *image.NRGBA, warstwa image.Image,
	obszar image.Rectangle, krycie float64, tryb shared.DesignPhotoBlendMode,
	maska image.Image) {

	skalowana := warstwa
	if obszar.Dx() != warstwa.Bounds().Dx() || obszar.Dy() != warstwa.Bounds().Dy() {
		docelowe := image.NewNRGBA(image.Rect(0, 0, obszar.Dx(), obszar.Dy()))
		draw.CatmullRom.Scale(docelowe, docelowe.Bounds(), warstwa, warstwa.Bounds(),
			draw.Over, nil)
		skalowana = docelowe
	}
	for y := 0; y < obszar.Dy(); y++ {
		for x := 0; x < obszar.Dx(); x++ {
			celX, celY := obszar.Min.X+x, obszar.Min.Y+y
			if celX < plotno.Bounds().Min.X || celY < plotno.Bounds().Min.Y ||
				celX >= plotno.Bounds().Max.X || celY >= plotno.Bounds().Max.Y {
				continue
			}
			gora := nrgbaPunktuDesignu(skalowana, skalowana.Bounds().Min.X+x,
				skalowana.Bounds().Min.Y+y)
			udzial := krycie * float64(gora.A) / 255
			if maska != nil {
				udzial *= udzialMaskiDesignu(maska, x, y, obszar.Dx(), obszar.Dy())
			}
			if udzial <= 0 {
				continue
			}
			dol := plotno.NRGBAAt(celX, celY)
			plotno.SetNRGBA(celX, celY, zmieszajPunktyDesignu(dol, gora, udzial, tryb))
		}
	}
}

// zmieszajPunktyDesignu miesza punkt górny z dolnym wedle trybu mieszania
// i udziału krycia warstwy górnej.
func zmieszajPunktyDesignu(dol, gora color.NRGBA, udzial float64,
	tryb shared.DesignPhotoBlendMode) color.NRGBA {

	mieszaj := func(pod, nad uint8) uint8 {
		p, n := float64(pod)/255, float64(nad)/255
		var wynik float64
		switch tryb {
		case shared.DesignPhotoBlendModeMultiply:
			wynik = p * n
		case shared.DesignPhotoBlendModeScreen:
			wynik = 1 - (1-p)*(1-n)
		case shared.DesignPhotoBlendModeOverlay:
			if p <= 0.5 {
				wynik = 2 * p * n
			} else {
				wynik = 1 - 2*(1-p)*(1-n)
			}
		default:
			wynik = n
		}
		return przytnijSkladowaDesignu((p*(1-udzial) + wynik*udzial) * 255)
	}
	nowaAlfa := float64(dol.A)/255 + udzial*(1-float64(dol.A)/255)
	return color.NRGBA{
		R: mieszaj(dol.R, gora.R), G: mieszaj(dol.G, gora.G), B: mieszaj(dol.B, gora.B),
		A: przytnijSkladowaDesignu(nowaAlfa * 255),
	}
}

// nrgbaPunktuDesignu odczytuje punkt obrazu w postaci bez wmnożonego krycia,
// niezależnie od formatu wewnętrznego obrazu.
func nrgbaPunktuDesignu(obraz image.Image, x, y int) color.NRGBA {
	if wprost, jest := obraz.(*image.NRGBA); jest {
		return wprost.NRGBAAt(x, y)
	}
	r, g, b, a := obraz.At(x, y).RGBA()
	if a == 0 {
		return color.NRGBA{}
	}
	return color.NRGBA{
		R: uint8(r * 255 / a), G: uint8(g * 255 / a), B: uint8(b * 255 / a),
		A: uint8(a >> 8),
	}
}

// udzialMaskiDesignu odczytuje udział maski w punkcie, skalując maskę do
// obszaru warstwy i biorąc jego jasność.
func udzialMaskiDesignu(maska image.Image, x, y, szerokosc, wysokosc int) float64 {
	granice := maska.Bounds()
	if granice.Dx() == 0 || granice.Dy() == 0 || szerokosc == 0 || wysokosc == 0 {
		return 1
	}
	maskaX := granice.Min.X + x*granice.Dx()/szerokosc
	maskaY := granice.Min.Y + y*granice.Dy()/wysokosc
	r, g, b, _ := maska.At(maskaX, maskaY).RGBA()
	// Maska jest odczytywana jako jasność, nie kanał krycia — bywa czarno-biała bez przezroczystości.
	return float64(int(r>>8)*299+int(g>>8)*587+int(b>>8)*114) / 1000 / 255
}

// kolorZNrgbaDesignu i nrgbaZKoloruDesignu przekładają punkt obrazu na kolor
// biblioteki barw i z powrotem. Kanał krycia przechodzi nietknięty: korekcje
// barwne zmieniają BARWĘ, nie przezroczystość.
func kolorZNrgbaDesignu(punkt color.NRGBA) colorful.Color {
	return colorful.Color{
		R: float64(punkt.R) / 255, G: float64(punkt.G) / 255, B: float64(punkt.B) / 255,
	}
}

func nrgbaZKoloruDesignu(barwa colorful.Color, krycie uint8) color.NRGBA {
	czysta := barwa.Clamped()
	return color.NRGBA{
		R: przytnijSkladowaDesignu(czysta.R * 255),
		G: przytnijSkladowaDesignu(czysta.G * 255),
		B: przytnijSkladowaDesignu(czysta.B * 255),
		A: krycie,
	}
}

// kolorHslDesignu i kolorHclDesignu składają kolor z współrzędnych przestrzeni
// percepcyjnych — cienka warstwa nad biblioteką, żeby wołający nie musiał jej
// importować i żeby przycięcie do zakresu było w jednym miejscu.
func kolorHslDesignu(odcien, nasycenie, jasnosc float64) colorful.Color {
	return colorful.Hsl(znormalizujKatDesignu(odcien), przytnijUlamekDesignu(nasycenie),
		przytnijUlamekDesignu(jasnosc)).Clamped()
}

func kolorHclDesignu(odcien, nasycenie, jasnosc float64) colorful.Color {
	return colorful.Hcl(znormalizujKatDesignu(odcien), przytnijUlamekDesignu(nasycenie),
		przytnijUlamekDesignu(jasnosc)).Clamped()
}

// tabliceKrzywychFotografiiDesignu to cztery tablice przeglądowe krzywych
// tonalnych: wspólna dla wszystkich kanałów i po jednej na kanał. Tablice, nie
// rachunek na punkt: obraz ma miliony punktów, a wartości wejściowych jest
// dwieście pięćdziesiąt sześć.
type tabliceKrzywychFotografiiDesignu struct {
	wspolna   [256]uint8
	czerwony  [256]uint8
	zielony   [256]uint8
	niebieski [256]uint8
}

// punktKrzywejDesignu to jeden punkt krzywej tonalnej w zapisie żądania,
// z nazwami pól polskimi i angielskimi.
type punktKrzywejDesignu struct {
	Wejscie float64 `json:"wejscie"`
	Wyjscie float64 `json:"wyjscie"`
	// Nazwy angielskie przyjmuje się obok polskich: okno bywa złożone z nazw pól kontraktu, in/out.
	In  *float64 `json:"in,omitempty"`
	Out *float64 `json:"out,omitempty"`
}

// tabliceKrzywychDesignu rozkłada zapis krzywych na tablice przeglądowe, po
// jednej wspólnej i po jednej na kanał barwy.
func tabliceKrzywychDesignu(zapis []byte) (tabliceKrzywychFotografiiDesignu, error) {
	var odczytane map[string][]punktKrzywejDesignu
	if err := json.Unmarshal(zapis, &odczytane); err != nil {
		return tabliceKrzywychFotografiiDesignu{}, fmt.Errorf(
			"krzywe tonalne nie są obiektem JSON kanał → wykaz punktów {wejscie, wyjscie}: %w", err)
	}
	tablice := tabliceKrzywychFotografiiDesignu{}
	for numer := 0; numer < 256; numer++ {
		tablice.wspolna[numer] = uint8(numer)
		tablice.czerwony[numer] = uint8(numer)
		tablice.zielony[numer] = uint8(numer)
		tablice.niebieski[numer] = uint8(numer)
	}
	for kanal, punkty := range odczytane {
		tablica, err := tablicaKrzywejDesignu(punkty)
		if err != nil {
			return tabliceKrzywychFotografiiDesignu{}, fmt.Errorf("krzywa kanału %s: %w", kanal, err)
		}
		switch strings.ToLower(strings.TrimSpace(kanal)) {
		case "rgb", "wspolna", "wspólna":
			tablice.wspolna = tablica
		case "r", "czerwony", "red":
			tablice.czerwony = tablica
		case "g", "zielony", "green":
			tablice.zielony = tablica
		case "b", "niebieski", "blue":
			tablice.niebieski = tablica
		default:
			return tabliceKrzywychFotografiiDesignu{}, fmt.Errorf(
				"krzywe wskazują kanał %q, którego rdzeń nie zna; kanały: rgb, r, g, b", kanal)
		}
	}
	return tablice, nil
}

// tablicaKrzywejDesignu składa tablicę przeglądową jednej krzywej.
//
// Punkty poza zakresem 0–255 są odmową, nie przycięciem: krzywa z punktem 400
// znaczy, że wołający liczy w innej skali, a przycięcie po cichu dałoby krzywą
// inną niż zamierzona.
func tablicaKrzywejDesignu(punkty []punktKrzywejDesignu) ([256]uint8, error) {
	var tablica [256]uint8
	if len(punkty) < 2 {
		return tablica, fmt.Errorf("krzywa o mniej niż dwóch punktach nie jest krzywą")
	}
	pary := make([][2]float64, 0, len(punkty))
	for _, punkt := range punkty {
		wejscie, wyjscie := punkt.Wejscie, punkt.Wyjscie
		if punkt.In != nil {
			wejscie = *punkt.In
		}
		if punkt.Out != nil {
			wyjscie = *punkt.Out
		}
		if wejscie < 0 || wejscie > 255 || wyjscie < 0 || wyjscie > 255 {
			return tablica, fmt.Errorf(
				"punkt {%v, %v} wychodzi poza skalę 0-255", wejscie, wyjscie)
		}
		pary = append(pary, [2]float64{wejscie, wyjscie})
	}
	sort.SliceStable(pary, func(i, j int) bool { return pary[i][0] < pary[j][0] })
	for numer := 0; numer < 256; numer++ {
		tablica[numer] = przytnijSkladowaDesignu(wartoscKrzywejDesignu(pary, float64(numer)))
	}
	return tablica, nil
}

// wartoscKrzywejDesignu odczytuje wartość krzywej liniowo między punktami.
// Poza skrajnymi punktami krzywa jest STAŁA — przedłużenie jej kierunkiem
// ostatniego odcinka wyprowadzałoby wartości poza skalę.
func wartoscKrzywejDesignu(pary [][2]float64, wejscie float64) float64 {
	if wejscie <= pary[0][0] {
		return pary[0][1]
	}
	if wejscie >= pary[len(pary)-1][0] {
		return pary[len(pary)-1][1]
	}
	for numer := 1; numer < len(pary); numer++ {
		if wejscie > pary[numer][0] {
			continue
		}
		poprzedni, biezacy := pary[numer-1], pary[numer]
		szerokosc := biezacy[0] - poprzedni[0]
		if szerokosc == 0 {
			return biezacy[1]
		}
		udzial := (wejscie - poprzedni[0]) / szerokosc
		return poprzedni[1] + udzial*(biezacy[1]-poprzedni[1])
	}
	return wejscie
}
