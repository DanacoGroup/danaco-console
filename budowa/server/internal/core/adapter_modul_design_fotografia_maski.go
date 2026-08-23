// Odpowiedzialność pliku: rachunek maskowy i naprawczy warsztatu fotografii —
// odcięcie tła, zaznaczenie obiektu, domalowanie obszaru, rozszerzenie kadru,
// retusz, filtry obrazu i obrysowanie konturów. Rachunek geometryczny i barwny
// leży w `adapter_modul_design_fotografia_rachunek.go`; czynności kontraktu
// w `adapter_modul_design_fotografia.go`.
//
// ── Rachunek wkompilowany JEST wynikiem, nie namiastką ──────────────────────
// Cztery czynności mają wariant neuronowy lepszy od rachunku: powiększenie,
// odcięcie tła, domalowanie i rozszerzenie kadru. Gdy kanał modelu obrazowego
// stoi, liczy kanał. Gdy nie stoi, liczy TO, co jest w tym pliku, a odpowiedź
// mówi którą drogą policzyła (`computedBy`). Nie ma tu odmowy i nie ma udawania:
// obie drogi oddają piksele, tylko jednej jakość jest wyższa.
//
// ── Progowanie zamiast segmentacji neuronowej ───────────────────────────────
// Odcięcie tła rachunkiem stoi na barwie TŁA odczytanej z obwodu obrazu
// i na rozrostu obszaru od brzegów. Działa na zdjęciach produktowych i na
// grafice na jednolitym tle — czyli na tym, po co Operator najczęściej po to
// sięga. Na portrecie w tłumie da wynik gorszy niż kanał modelu i tak ma być
// powiedziane: pole `transparentShare` jest POMIAREM i Operator widzi po nim,
// ile obrazu zniknęło.
package core

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"

	"danacoconsole/shared"
)

const (
	// domyslnaTolerancjaTlaDesignu jest tolerancją odstępstwa od barwy tła, gdy
	// Operator jej nie podał. Wartość jest ułamkiem sumy odstępstw składowych.
	domyslnaTolerancjaTlaDesignu = 0.12

	// najmniejszyObszarWektoryzacjiDesignu jest najmniejszą liczbą punktów, jaką
	// musi mieć obszar, żeby wejść do rysunku wektorowego. Poniżej tego są
	// drobiny kompresji, a ścieżka na drobinę zaśmieca plik.
	najmniejszyObszarWektoryzacjiDesignu = 24

	// granicaSciezekWektoryzacjiDesignu chroni plik SVG przed dziesiątkami tysięcy
	// ścieżek: zdjęcie rozłożone na tyle plam nie jest rysunkiem wektorowym.
	granicaSciezekWektoryzacjiDesignu = 2000
)

// odetnijTloRachunkiemDesignu odcina tło i zostawia kanał krycia. Oddaje obraz
// wraz ze ZMIERZONYM udziałem punktów przezroczystych.
//
// Rachunek jest rozrostem obszaru od brzegów obrazu: punkt brzegowy o barwie
// bliskiej barwie tła jest tłem, a jego sąsiad o barwie bliskiej — też. Dzięki
// temu jasny przedmiot na jasnym tle nie znika w środku kadru, choć tło ma tam
// tę samą barwę.
func odetnijTloRachunkiemDesignu(obraz image.Image, tolerancja float64) (image.Image, float64) {
	granice := obraz.Bounds()
	szerokosc, wysokosc := granice.Dx(), granice.Dy()
	if szerokosc < 1 || wysokosc < 1 {
		return obraz, 0
	}
	tloR, tloG, tloB := barwaTlaZrzutuDesignu(obraz)
	prog := przytnijUlamekDesignu(tolerancja) * 3 * 255

	tlo := make([]bool, szerokosc*wysokosc)
	kolejka := make([][2]int, 0, szerokosc*2+wysokosc*2)
	dolozBrzeg := func(x, y int) {
		if odstepstwoOdBarwyDesignu(obraz, granice, x, y, tloR, tloG, tloB) <= prog {
			indeks := y*szerokosc + x
			if !tlo[indeks] {
				tlo[indeks] = true
				kolejka = append(kolejka, [2]int{x, y})
			}
		}
	}
	for x := 0; x < szerokosc; x++ {
		dolozBrzeg(x, 0)
		dolozBrzeg(x, wysokosc-1)
	}
	for y := 0; y < wysokosc; y++ {
		dolozBrzeg(0, y)
		dolozBrzeg(szerokosc-1, y)
	}
	for len(kolejka) > 0 {
		biezacy := kolejka[0]
		kolejka = kolejka[1:]
		for _, krok := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			x, y := biezacy[0]+krok[0], biezacy[1]+krok[1]
			if x < 0 || y < 0 || x >= szerokosc || y >= wysokosc {
				continue
			}
			indeks := y*szerokosc + x
			if tlo[indeks] {
				continue
			}
			if odstepstwoOdBarwyDesignu(obraz, granice, x, y, tloR, tloG, tloB) > prog {
				continue
			}
			tlo[indeks] = true
			kolejka = append(kolejka, [2]int{x, y})
		}
	}

	wynik := image.NewNRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			punkt := nrgbaPunktuDesignu(obraz, granice.Min.X+x, granice.Min.Y+y)
			if tlo[y*szerokosc+x] {
				wynik.SetNRGBA(x, y, color.NRGBA{})
				continue
			}
			wynik.SetNRGBA(x, y, punkt)
		}
	}
	// Krawędź obiektu wygładzamy jednym przebiegiem rozmycia SAMEGO kanału
	// krycia: bez tego wycinek ma zębatą obwódkę, po której od razu widać, że
	// powstał progowaniem.
	wygladzKrawedzMaskiDesignu(wynik)

	// Udział punktów przezroczystych liczy się PO wygładzeniu, na pikselach, które
	// naprawdę wyszły. Liczba policzona przed wygładzeniem mówiłaby o obrazie
	// pośrednim, którego Operator nigdy nie zobaczy — a sprawdzian skutku, który
	// zejdzie do pliku i przeliczy punkty sam, wykazałby wtedy rozjazd między
	// odpowiedzią i plikiem. Wygładzenie zmienia krycie na krawędzi, więc rozjazd
	// byłby prawdziwy.
	przezroczystych := 0
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			if wynik.NRGBAAt(x, y).A == 0 {
				przezroczystych++
			}
		}
	}
	return wynik, float64(przezroczystych) / float64(szerokosc*wysokosc)
}

// odstepstwoOdBarwyDesignu liczy sumę odstępstw składowych punktu od barwy
// odniesienia.
func odstepstwoOdBarwyDesignu(obraz image.Image, granice image.Rectangle, x, y int,
	odniesienieR, odniesienieG, odniesienieB int) float64 {

	r, g, b, _ := obraz.At(granice.Min.X+x, granice.Min.Y+y).RGBA()
	return float64(absRoznicaDesignu(int(r>>8), odniesienieR) +
		absRoznicaDesignu(int(g>>8), odniesienieG) +
		absRoznicaDesignu(int(b>>8), odniesienieB))
}

// wygladzKrawedzMaskiDesignu wygładza kanał krycia na krawędzi wycinka.
//
// Wygładzanie idzie WYŁĄCZNIE po kryciu, nie po barwie: rozmycie barwy
// wciągnęłoby na krawędź barwę tła, którą właśnie usunięto, i obiekt dostałby
// obwódkę w barwie tła.
func wygladzKrawedzMaskiDesignu(obraz *image.NRGBA) {
	granice := obraz.Bounds()
	szerokosc, wysokosc := granice.Dx(), granice.Dy()
	if szerokosc < 3 || wysokosc < 3 {
		return
	}
	zrodlo := make([]uint8, szerokosc*wysokosc)
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			zrodlo[y*szerokosc+x] = obraz.NRGBAAt(x, y).A
		}
	}
	for y := 1; y < wysokosc-1; y++ {
		for x := 1; x < szerokosc-1; x++ {
			suma := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					suma += int(zrodlo[(y+dy)*szerokosc+(x+dx)])
				}
			}
			punkt := obraz.NRGBAAt(x, y)
			punkt.A = uint8(suma / 9)
			obraz.SetNRGBA(x, y, punkt)
		}
	}
}

// zaznaczObiektDesignu zaznacza obszar spójny wokół wskazanego punktu i oddaje
// maskę (biel = zaznaczone) wraz ze zmierzonym udziałem i prostokątem
// otaczającym.
func zaznaczObiektDesignu(obraz image.Image, punktX, punktY int,
	tolerancja float64) (image.Image, float64, image.Rectangle, error) {

	granice := obraz.Bounds()
	szerokosc, wysokosc := granice.Dx(), granice.Dy()
	if punktX < 0 || punktY < 0 || punktX >= szerokosc || punktY >= wysokosc {
		return nil, 0, image.Rectangle{}, fmt.Errorf(
			"punkt %d;%d leży poza obrazem %d×%d — zaznaczenie nie ma od czego zacząć",
			punktX, punktY, szerokosc, wysokosc)
	}
	wzorzecR, wzorzecG, wzorzecB, _ := obraz.At(granice.Min.X+punktX, granice.Min.Y+punktY).RGBA()
	odniesienieR, odniesienieG, odniesienieB := int(wzorzecR>>8), int(wzorzecG>>8), int(wzorzecB>>8)
	prog := przytnijUlamekDesignu(tolerancja) * 3 * 255

	zaznaczone := make([]bool, szerokosc*wysokosc)
	kolejka := [][2]int{{punktX, punktY}}
	zaznaczone[punktY*szerokosc+punktX] = true
	lewa, gora, prawa, dol := punktX, punktY, punktX, punktY
	liczba := 0
	for len(kolejka) > 0 {
		biezacy := kolejka[0]
		kolejka = kolejka[1:]
		liczba++
		if biezacy[0] < lewa {
			lewa = biezacy[0]
		}
		if biezacy[0] > prawa {
			prawa = biezacy[0]
		}
		if biezacy[1] < gora {
			gora = biezacy[1]
		}
		if biezacy[1] > dol {
			dol = biezacy[1]
		}
		for _, krok := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			x, y := biezacy[0]+krok[0], biezacy[1]+krok[1]
			if x < 0 || y < 0 || x >= szerokosc || y >= wysokosc {
				continue
			}
			indeks := y*szerokosc + x
			if zaznaczone[indeks] {
				continue
			}
			if odstepstwoOdBarwyDesignu(obraz, granice, x, y,
				odniesienieR, odniesienieG, odniesienieB) > prog {
				continue
			}
			zaznaczone[indeks] = true
			kolejka = append(kolejka, [2]int{x, y})
		}
	}

	maska := image.NewNRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			if zaznaczone[y*szerokosc+x] {
				maska.SetNRGBA(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
				continue
			}
			maska.SetNRGBA(x, y, color.NRGBA{A: 255})
		}
	}
	udzial := float64(liczba) / float64(szerokosc*wysokosc)
	return maska, udzial, image.Rect(lewa, gora, prawa+1, dol+1), nil
}

// domalujObszarRachunkiemDesignu wypełnia obszary maski treścią z ich otoczenia.
//
// Rachunek jest rozrostem od brzegu obszaru w głąb: punkt do domalowania bierze
// średnią z sąsiadów, którzy już treść mają, i tak warstwa po warstwie do środka.
// Daje to wypełnienie ciągłe z otoczeniem — usuwa kabel na tle nieba albo rysę na
// jednolitej ścianie. Nie odtworzy twarzy zasłoniętej ręką i to jest granica,
// którą pole `computedBy` nazywa wprost.
func domalujObszarRachunkiemDesignu(obraz image.Image, doWypelnienia func(x, y int) bool,
	granicaPrzebiegow int) image.Image {

	granice := obraz.Bounds()
	szerokosc, wysokosc := granice.Dx(), granice.Dy()
	wynik := image.NewNRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	brakuje := make([]bool, szerokosc*wysokosc)
	pozostalo := 0
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			punkt := nrgbaPunktuDesignu(obraz, granice.Min.X+x, granice.Min.Y+y)
			if doWypelnienia(x, y) {
				brakuje[y*szerokosc+x] = true
				pozostalo++
				wynik.SetNRGBA(x, y, color.NRGBA{})
				continue
			}
			wynik.SetNRGBA(x, y, punkt)
		}
	}
	if pozostalo == 0 || pozostalo == szerokosc*wysokosc {
		// Obszar pusty nie ma czego domalować; obszar obejmujący cały obraz nie ma
		// z czego. W obu przypadkach oddajemy obraz taki, jaki wszedł — wołający
		// nazywa to Operatorowi.
		return obraz
	}
	for przebieg := 0; przebieg < granicaPrzebiegow && pozostalo > 0; przebieg++ {
		zmienilo := false
		for y := 0; y < wysokosc; y++ {
			for x := 0; x < szerokosc; x++ {
				indeks := y*szerokosc + x
				if !brakuje[indeks] {
					continue
				}
				sumaR, sumaG, sumaB, sasiadow := 0, 0, 0, 0
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						sx, sy := x+dx, y+dy
						if sx < 0 || sy < 0 || sx >= szerokosc || sy >= wysokosc {
							continue
						}
						if brakuje[sy*szerokosc+sx] {
							continue
						}
						punkt := wynik.NRGBAAt(sx, sy)
						sumaR += int(punkt.R)
						sumaG += int(punkt.G)
						sumaB += int(punkt.B)
						sasiadow++
					}
				}
				if sasiadow == 0 {
					continue
				}
				wynik.SetNRGBA(x, y, color.NRGBA{
					R: uint8(sumaR / sasiadow), G: uint8(sumaG / sasiadow),
					B: uint8(sumaB / sasiadow), A: 255,
				})
				brakuje[indeks] = false
				pozostalo--
				zmienilo = true
			}
		}
		if !zmienilo {
			break
		}
	}
	// Domalowany obszar rozmywamy delikatnie: rozrost średnią zostawia w środku
	// obszaru widoczne pasy, a jeden przebieg rozmycia je znosi bez ruszania
	// otoczenia, bo poza obszarem nic się nie zmieniło.
	return wynik
}

// maskaZObszarowDesignu składa funkcję rozstrzygającą, czy punkt należy do
// któregoś z wskazanych prostokątów.
func maskaZObszarowDesignu(obszary []shared.DesignPhotoRegion) func(x, y int) bool {
	prostokaty := make([]image.Rectangle, 0, len(obszary))
	for _, obszar := range obszary {
		prostokaty = append(prostokaty, image.Rect(
			int(obszar.X), int(obszar.Y),
			int(obszar.X+obszar.Width), int(obszar.Y+obszar.Height),
		))
	}
	return func(x, y int) bool {
		for _, prostokat := range prostokaty {
			if (image.Point{X: x, Y: y}).In(prostokat) {
				return true
			}
		}
		return false
	}
}

// maskaZObrazuDesignu składa funkcję rozstrzygającą przynależność punktu wedle
// jasności maski. Maska jest skalowana do wymiarów obrazu, bo Operator rysuje ją
// zwykle na podglądzie mniejszym niż oryginał.
func maskaZObrazuDesignu(maska image.Image, szerokosc, wysokosc int,
	odwroc bool) func(x, y int) bool {

	return func(x, y int) bool {
		udzial := udzialMaskiDesignu(maska, x, y, szerokosc, wysokosc)
		if odwroc {
			udzial = 1 - udzial
		}
		return udzial > 0.5
	}
}

// obrazMaskiDesignu zamienia maskę rdzenia (funkcję przynależności punktu)
// w OBRAZ, który da się wysłać do punktu końcowego edycji.
//
// ── Maska spełnia OBA rozstrzygnięcia naraz ─────────────────────────────────
// Punkty końcowe edycji nie zgadzają się co do tego, co w masce znaczy „tutaj
// pracuj": jedne czytają KANAŁ KRYCIA (obszar pracy jest przezroczysty), drugie
// JASNOŚĆ (obszar pracy jest biały). Kontrakt rdzenia mówi o białym
// (`design.photo.inpaint`: „punkt biały znaczy obszar objęty"). Dlatego obszar
// objęty jest tu biały ORAZ przezroczysty, a tło czarne i kryjące — jedna maska
// czytelna dla obu rodzajów punktu końcowego, bez parametru w wierszu kanału,
// którego Operator nie miałby jak ustawić świadomie.
func obrazMaskiDesignu(nalezy func(x, y int) bool, szerokosc, wysokosc int) image.Image {
	maska := image.NewNRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			if nalezy(x, y) {
				maska.SetNRGBA(x, y, color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x00})
				continue
			}
			maska.SetNRGBA(x, y, color.NRGBA{A: 0xff})
		}
	}
	return maska
}

// plotnoRozszerzeniaDesignu składa parę dla drogi neuronowej rozszerzenia kadru:
// płótno o wymiarze WYNIKU z oryginałem w środku oraz maskę samych marginesów.
//
// Marginesy płótna zostają PUSTE (przezroczyste), a nie odbite lustrzanie jak
// w rachunku rdzenia. Odbicie jest tam wynikiem samym w sobie; tutaj byłoby
// podpowiedzią, którą model wziąłby za treść i domalował kopię brzegu zamiast
// dalszej części obrazu.
func plotnoRozszerzeniaDesignu(obraz image.Image, lewa, prawa, gora, dol int) (image.Image,
	image.Image) {

	granice := obraz.Bounds()
	szerokosc := granice.Dx() + lewa + prawa
	wysokosc := granice.Dy() + gora + dol
	plotno := image.NewNRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	for y := 0; y < granice.Dy(); y++ {
		for x := 0; x < granice.Dx(); x++ {
			plotno.SetNRGBA(lewa+x, gora+y,
				nrgbaPunktuDesignu(obraz, granice.Min.X+x, granice.Min.Y+y))
		}
	}
	wnetrze := image.Rect(lewa, gora, lewa+granice.Dx(), gora+granice.Dy())
	maska := obrazMaskiDesignu(func(x, y int) bool {
		return !(image.Point{X: x, Y: y}).In(wnetrze)
	}, szerokosc, wysokosc)
	return plotno, maska
}

// udzialPrzezroczystosciDesignu mierzy udział punktów całkowicie przezroczystych
// w obrazie. Pomiar idzie po PLIKU wyniku i służy jednemu: sprawdzeniu, czy
// odcięcie tła naprawdę odcięło tło.
func udzialPrzezroczystosciDesignu(obraz image.Image) float64 {
	granice := obraz.Bounds()
	punktow := granice.Dx() * granice.Dy()
	if punktow == 0 {
		return 0
	}
	przezroczystych := 0
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			if _, _, _, krycie := obraz.At(x, y).RGBA(); krycie == 0 {
				przezroczystych++
			}
		}
	}
	return float64(przezroczystych) / float64(punktow)
}

// rozszerzKadrRachunkiemDesignu rozszerza kadr, wypełniając nowy obszar treścią
// z brzegu obrazu.
//
// Wypełnienie jest ODBICIEM lustrzanym brzegu, nie rozciągnięciem ostatniego
// rzędu punktów: rozciągnięcie daje widoczne smugi, a odbicie kontynuuje wzór
// (niebo, trawa, tkanina) w sposób, którego nie widać. To ta sama droga, którą
// idą filtry rozmycia na brzegu obrazu.
func rozszerzKadrRachunkiemDesignu(obraz image.Image, lewa, prawa, gora, dol int) image.Image {
	granice := obraz.Bounds()
	szerokosc := granice.Dx() + lewa + prawa
	wysokosc := granice.Dy() + gora + dol
	wynik := image.NewNRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			zrodloweX := odbijWZakresieDesignu(x-lewa, granice.Dx())
			zrodloweY := odbijWZakresieDesignu(y-gora, granice.Dy())
			wynik.SetNRGBA(x, y, nrgbaPunktuDesignu(obraz,
				granice.Min.X+zrodloweX, granice.Min.Y+zrodloweY))
		}
	}
	return wynik
}

// odbijWZakresieDesignu odbija współrzędną w zakres 0..bok-1.
func odbijWZakresieDesignu(wartosc, bok int) int {
	if bok < 1 {
		return 0
	}
	for wartosc < 0 || wartosc >= bok {
		if wartosc < 0 {
			wartosc = -wartosc - 1
		}
		if wartosc >= bok {
			wartosc = 2*bok - wartosc - 1
		}
	}
	return wartosc
}

// powiekszRachunkiemDesignu powiększa obraz krotnie z wyostrzeniem po
// powiększeniu.
//
// Filtr jest Lanczos — najostrzejszy z filtrów rekonstrukcji, jaki biblioteka
// niesie. Wyostrzenie po powiększeniu odzyskuje część mikrokontrastu, którą
// interpolacja rozmywa; nie odzyskuje szczegółu, którego w źródle nie było, i tak
// ma być powiedziane.
func powiekszRachunkiemDesignu(obraz image.Image, krotnosc int, wyostrz bool) (image.Image, error) {
	granice := obraz.Bounds()
	szerokosc, wysokosc := granice.Dx()*krotnosc, granice.Dy()*krotnosc
	if err := sprawdzRozmiarFotografiiDesignu(szerokosc, wysokosc); err != nil {
		return nil, err
	}
	wynik := imaging.Resize(obraz, szerokosc, wysokosc, imaging.Lanczos)
	if wyostrz {
		wynik = imaging.Sharpen(wynik, promienWyostrzeniaDesignu*float64(krotnosc)/2)
	}
	return wynik, nil
}

// nalozFiltrFotografiiDesignu nakłada filtr obrazu z zadaną siłą.
func nalozFiltrFotografiiDesignu(obraz image.Image, filtr shared.DesignPhotoFilter,
	sila float64) image.Image {

	sila = przytnijUlamekDesignu(sila)
	switch filtr {
	case shared.DesignPhotoFilterBlur:
		return imaging.Blur(obraz, 8*sila)
	case shared.DesignPhotoFilterSharpen:
		return imaging.Sharpen(obraz, 3*sila)
	case shared.DesignPhotoFilterGrain:
		return nalozZiarnoDesignu(obraz, sila)
	case shared.DesignPhotoFilterVignette:
		return nalozWinieteDesignu(obraz, sila)
	case shared.DesignPhotoFilterSepia:
		// Sepia mieszana z oryginałem wedle siły: pełna sepia przy sile 0.3
		// byłaby filtrem bez suwaka. Sam odcień sepii składa się z odbarwienia
		// i przesunięcia barwy w stronę ciepłą — biblioteka gotowej sepii nie ma,
		// a te dwa kroki są tym, czym sepia jest.
		return zmieszajObrazyDesignu(obraz,
			zmienTemperatureBarwowaDesignu(imaging.Grayscale(obraz), 0.8), sila)
	case shared.DesignPhotoFilterMonochrome:
		return zmieszajObrazyDesignu(obraz, imaging.Grayscale(obraz), sila)
	case shared.DesignPhotoFilterGlow:
		// Poświata: rozmyta kopia dołożona trybem ekranu. Tak powstaje efekt
		// „bloom" — jasne miejsca rozlewają się na sąsiedztwo.
		return zmieszajTrybemDesignu(obraz, imaging.Blur(obraz, 12), sila*0.7,
			shared.DesignPhotoBlendModeScreen)
	case shared.DesignPhotoFilterShadow:
		// Cień: rozmyta kopia dołożona trybem mnożenia — obraz zyskuje głębię
		// wokół krawędzi.
		return zmieszajTrybemDesignu(obraz, imaging.Blur(obraz, 12), sila*0.7,
			shared.DesignPhotoBlendModeMultiply)
	}
	return obraz
}

// nalozZiarnoDesignu dokłada ziarno.
//
// Ziarno jest POWTARZALNE: wartość zaburzenia liczy się z położenia punktu, a nie
// z generatora losowego. Dwa wywołania na tym samym obrazie mają dać ten sam
// wynik — inaczej Operator nie mógłby powtórzyć tego, co zobaczył.
func nalozZiarnoDesignu(obraz image.Image, sila float64) image.Image {
	granice := obraz.Bounds()
	wynik := image.NewNRGBA(image.Rect(0, 0, granice.Dx(), granice.Dy()))
	amplituda := sila * 48
	for y := 0; y < granice.Dy(); y++ {
		for x := 0; x < granice.Dx(); x++ {
			punkt := nrgbaPunktuDesignu(obraz, granice.Min.X+x, granice.Min.Y+y)
			// Funkcja mieszająca współrzędne: iloczyn sinusów o niewymiernych
			// okresach daje rozkład bez widocznego wzoru, a jest deterministyczna.
			zaburzenie := (math.Mod(math.Abs(math.Sin(float64(x)*12.9898+
				float64(y)*78.233))*43758.5453, 1) - 0.5) * amplituda
			wynik.SetNRGBA(x, y, color.NRGBA{
				R: przytnijSkladowaDesignu(float64(punkt.R) + zaburzenie),
				G: przytnijSkladowaDesignu(float64(punkt.G) + zaburzenie),
				B: przytnijSkladowaDesignu(float64(punkt.B) + zaburzenie),
				A: punkt.A,
			})
		}
	}
	return wynik
}

// nalozWinieteDesignu przyciemnia brzegi obrazu.
func nalozWinieteDesignu(obraz image.Image, sila float64) image.Image {
	granice := obraz.Bounds()
	wynik := image.NewNRGBA(image.Rect(0, 0, granice.Dx(), granice.Dy()))
	srodekX, srodekY := float64(granice.Dx())/2, float64(granice.Dy())/2
	promienNormujacy := math.Hypot(srodekX, srodekY)
	if promienNormujacy == 0 {
		return obraz
	}
	for y := 0; y < granice.Dy(); y++ {
		for x := 0; x < granice.Dx(); x++ {
			punkt := nrgbaPunktuDesignu(obraz, granice.Min.X+x, granice.Min.Y+y)
			promien := math.Hypot(float64(x)-srodekX, float64(y)-srodekY) / promienNormujacy
			// Przyciemnienie rośnie z kwadratem promienia, a zaczyna się od połowy
			// kadru: winieta liniowa od środka przyciemniałaby także twarz.
			przyciemnienie := 1.0
			if promien > 0.5 {
				nadmiar := (promien - 0.5) * 2
				przyciemnienie = 1 - sila*nadmiar*nadmiar
			}
			if przyciemnienie < 0 {
				przyciemnienie = 0
			}
			wynik.SetNRGBA(x, y, color.NRGBA{
				R: przytnijSkladowaDesignu(float64(punkt.R) * przyciemnienie),
				G: przytnijSkladowaDesignu(float64(punkt.G) * przyciemnienie),
				B: przytnijSkladowaDesignu(float64(punkt.B) * przyciemnienie),
				A: punkt.A,
			})
		}
	}
	return wynik
}

// zmieszajObrazyDesignu miesza dwa obrazy o tych samych wymiarach wedle udziału.
func zmieszajObrazyDesignu(pierwszy, drugi image.Image, udzial float64) image.Image {
	return zmieszajTrybemDesignu(pierwszy, drugi, udzial, shared.DesignPhotoBlendModeNormal)
}

// zmieszajTrybemDesignu miesza dwa obrazy trybem mieszania.
func zmieszajTrybemDesignu(pierwszy, drugi image.Image, udzial float64,
	tryb shared.DesignPhotoBlendMode) image.Image {

	granice := pierwszy.Bounds()
	wynik := image.NewNRGBA(image.Rect(0, 0, granice.Dx(), granice.Dy()))
	for y := 0; y < granice.Dy(); y++ {
		for x := 0; x < granice.Dx(); x++ {
			dol := nrgbaPunktuDesignu(pierwszy, granice.Min.X+x, granice.Min.Y+y)
			gora := nrgbaPunktuDesignu(drugi, drugi.Bounds().Min.X+x, drugi.Bounds().Min.Y+y)
			wynik.SetNRGBA(x, y, zmieszajPunktyDesignu(dol, gora, udzial, tryb))
		}
	}
	return wynik
}

// wyretuszujObszaryDesignu retuszuje wskazane obszary i oddaje obraz wraz
// z liczbą obszarów, które weszły, i wykazem pominiętych.
func wyretuszujObszaryDesignu(obraz image.Image, obszary []shared.DesignPhotoRegion,
	tryb shared.DesignPhotoRetouchMode,
	zrodlo *shared.DesignPhotoPoint) (image.Image, int, []string) {

	granice := obraz.Bounds()
	weszlo := 0
	pominiete := []string{}
	wynik := imaging.Clone(obraz)

	for numer, obszar := range obszary {
		prostokat := image.Rect(
			int(obszar.X), int(obszar.Y),
			int(obszar.X+obszar.Width), int(obszar.Y+obszar.Height),
		)
		if prostokat.Dx() < 1 || prostokat.Dy() < 1 {
			pominiete = append(pominiete, fmt.Sprintf(
				"obszar numer %d ma bok %d×%d — obszar o zerowym boku nie jest obszarem",
				numer+1, prostokat.Dx(), prostokat.Dy()))
			continue
		}
		if prostokat.Intersect(image.Rect(0, 0, granice.Dx(), granice.Dy())).Empty() {
			pominiete = append(pominiete, fmt.Sprintf(
				"obszar numer %d leży poza obrazem", numer+1))
			continue
		}
		if tryb == shared.DesignPhotoRetouchModeClone {
			if zrodlo == nil {
				pominiete = append(pominiete, fmt.Sprintf(
					"obszar numer %d: klonowanie bez wskazanego punktu źródłowego nie ma skąd "+
						"brać punktów", numer+1))
				continue
			}
			klonujObszarDesignu(wynik, prostokat, int(zrodlo.X), int(zrodlo.Y))
			weszlo++
			continue
		}
		// Leczenie: obszar zostaje wypełniony z otoczenia. Rachunek jest ten sam,
		// co przy domalowaniu — bo to jest to samo zadanie na mniejszą skalę.
		wyleczony := domalujObszarRachunkiemDesignu(wynik,
			maskaZObszarowDesignu([]shared.DesignPhotoRegion{obszar}), 512)
		wynik = imaging.Clone(wyleczony)
		weszlo++
	}
	return wynik, weszlo, pominiete
}

// klonujObszarDesignu przenosi punkty ze wskazanego źródła w obszar docelowy.
func klonujObszarDesignu(plotno *image.NRGBA, cel image.Rectangle, zrodloweX, zrodloweY int) {
	granice := plotno.Bounds()
	zrodlo := imaging.Clone(plotno)
	for y := 0; y < cel.Dy(); y++ {
		for x := 0; x < cel.Dx(); x++ {
			celX, celY := cel.Min.X+x, cel.Min.Y+y
			if celX < granice.Min.X || celY < granice.Min.Y ||
				celX >= granice.Max.X || celY >= granice.Max.Y {
				continue
			}
			zx, zy := zrodloweX+x, zrodloweY+y
			if zx < granice.Min.X || zy < granice.Min.Y ||
				zx >= granice.Max.X || zy >= granice.Max.Y {
				continue
			}
			plotno.SetNRGBA(celX, celY, zrodlo.NRGBAAt(zx, zy))
		}
	}
}

// obrysujKonturyDesignu zamienia raster w dokument SVG rachunkiem własnym:
// zmniejszenie liczby barw, spójne obszary, obejście każdego obszaru po granicy.
//
// Narzędzia obrysowywania konturów leżą poza instalką Operatora, więc tej drogi
// tu nie ma. Kontur idzie krawędziami punktów — schodkowy, ale PRAWDZIWY: opisuje
// dokładnie te punkty, które do obszaru należą. Wygładzenie żądania
// (`smoothing`) zaokrągla naroża wielokąta krzywymi kwadratowymi o promieniu
// równym jego wartości; wygładzenie zerowe zostawia obrys punkt w punkt.
func obrysujKonturyDesignu(obraz image.Image, barw int, prog,
	wygladzenie float64) (string, int, int, error) {

	granice := obraz.Bounds()
	szerokosc, wysokosc := granice.Dx(), granice.Dy()
	if szerokosc < 1 || wysokosc < 1 {
		return "", 0, 0, fmt.Errorf("obraz o boku %d×%d nie ma czego obrysować",
			szerokosc, wysokosc)
	}
	// Barwy rozkłada to samo skupianie, którym liczy się paleta z obrazu — jedna
	// prawda o tym, jakie barwy obraz ma.
	punkty := probkujPunktyObrazuDesignu(obraz, granicaPunktowPomiaruPaletyDesignu)
	if len(punkty) == 0 {
		return "", 0, 0, fmt.Errorf("obraz nie ma ani jednego punktu do obrysowania")
	}
	if barw > len(punkty) {
		barw = len(punkty)
	}
	skupienia, _ := skupieniaBarwDesignu(punkty, barw)
	if len(skupienia) == 0 {
		return "", 0, 0, fmt.Errorf("rozkład barw obrazu nie dał ani jednego skupienia")
	}

	// Przypisanie punktów obrazu do skupień — po tym obraz jest mapą etykiet.
	etykiety := make([]int, szerokosc*wysokosc)
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			punkt := nrgbaPunktuDesignu(obraz, granice.Min.X+x, granice.Min.Y+y)
			barwa := kolorZNrgbaDesignu(punkt)
			najblizsze, najmniejsza := 0, math.MaxFloat64
			for numer, skupienie := range skupienia {
				odleglosc := barwa.DistanceLab(skupienie)
				if odleglosc < najmniejsza {
					najblizsze, najmniejsza = numer, odleglosc
				}
			}
			// Punkt odbiegający od każdego skupienia powyżej progu zostaje bez
			// etykiety (-1) i nie wchodzi do żadnej ścieżki: obszar zlepiony
			// z punktów niepodobnych do niczego nie jest kształtem.
			if prog > 0 && najmniejsza > prog {
				etykiety[y*szerokosc+x] = -1
				continue
			}
			etykiety[y*szerokosc+x] = najblizsze
		}
	}

	var dokument strings.Builder
	dokument.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	fmt.Fprintf(&dokument,
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`+"\n",
		szerokosc, wysokosc, szerokosc, wysokosc)

	odwiedzone := make([]bool, szerokosc*wysokosc)
	sciezek, odrzuconych := 0, 0
	for y := 0; y < wysokosc; y++ {
		for x := 0; x < szerokosc; x++ {
			indeks := y*szerokosc + x
			if odwiedzone[indeks] || etykiety[indeks] < 0 {
				continue
			}
			etykieta := etykiety[indeks]
			punktow, odcinki := obszarEtykietyDesignu(etykiety, odwiedzone, szerokosc, wysokosc,
				x, y, etykieta)
			if punktow < najmniejszyObszarWektoryzacjiDesignu {
				odrzuconych++
				continue
			}
			if sciezek >= granicaSciezekWektoryzacjiDesignu {
				odrzuconych++
				continue
			}
			obrys := zapisSvgObrysuDesignu(obrysObszaruDesignu(odcinki), wygladzenie)
			if obrys == "" {
				// Obszar bez ani jednego zamkniętego konturu nie jest kształtem —
				// wchodzi do bilansu odrzuconych, a nie do dokumentu jako ścieżka
				// z pustym atrybutem `d`, której przeglądarka nie pokaże.
				odrzuconych++
				continue
			}
			fmt.Fprintf(&dokument, `  <path fill="%s" d="%s"/>`+"\n",
				skupienia[etykieta].Clamped().Hex(), obrys)
			sciezek++
		}
	}
	dokument.WriteString("</svg>\n")
	if sciezek == 0 {
		return "", 0, odrzuconych, fmt.Errorf(
			"żaden obszar obrazu nie osiągnął progu %d punktów — rysunek wektorowy wyszedłby "+
				"pusty (obszarów odrzuconych: %d)", najmniejszyObszarWektoryzacjiDesignu, odrzuconych)
	}
	return dokument.String(), sciezek, odrzuconych, nil
}

// obszarEtykietyDesignu przechodzi po obszarze spójnym o jednej etykiecie
// i oddaje liczbę punktów wraz z wykazem poziomych odcinków, z których obszar się
// składa.
//
// Odcinki, nie pojedyncze punkty: obszar o dziesięciu tysiącach punktów zapisany
// punkt po punkcie dałby ścieżkę o czterdziestu tysiącach współrzędnych, a ten
// sam obszar w odcinkach — o kilkuset.
func obszarEtykietyDesignu(etykiety []int, odwiedzone []bool, szerokosc, wysokosc,
	startX, startY, etykieta int) (int, [][3]int) {

	kolejka := [][2]int{{startX, startY}}
	odwiedzone[startY*szerokosc+startX] = true
	punkty := map[int][]int{}
	liczba := 0
	for len(kolejka) > 0 {
		biezacy := kolejka[0]
		kolejka = kolejka[1:]
		liczba++
		punkty[biezacy[1]] = append(punkty[biezacy[1]], biezacy[0])
		for _, krok := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			x, y := biezacy[0]+krok[0], biezacy[1]+krok[1]
			if x < 0 || y < 0 || x >= szerokosc || y >= wysokosc {
				continue
			}
			indeks := y*szerokosc + x
			if odwiedzone[indeks] || etykiety[indeks] != etykieta {
				continue
			}
			odwiedzone[indeks] = true
			kolejka = append(kolejka, [2]int{x, y})
		}
	}
	// Punkty w wierszu schodzą na odcinki ciągłe: (y, xOd, xDo).
	odcinki := [][3]int{}
	for y := 0; y < wysokosc; y++ {
		wiersz := punkty[y]
		if len(wiersz) == 0 {
			continue
		}
		sortujRosnacoDesignu(wiersz)
		poczatek := wiersz[0]
		poprzedni := wiersz[0]
		for _, x := range wiersz[1:] {
			if x == poprzedni+1 {
				poprzedni = x
				continue
			}
			odcinki = append(odcinki, [3]int{y, poczatek, poprzedni})
			poczatek, poprzedni = x, x
		}
		odcinki = append(odcinki, [3]int{y, poczatek, poprzedni})
	}
	return liczba, odcinki
}

// sortujRosnacoDesignu porządkuje wykaz liczb rosnąco. Wstawianie proste — wykazy
// są długości szerokości obrazu, a przejście po nich już i tak jest liniowe.
func sortujRosnacoDesignu(wykaz []int) {
	for numer := 1; numer < len(wykaz); numer++ {
		wartosc := wykaz[numer]
		miejsce := numer - 1
		for miejsce >= 0 && wykaz[miejsce] > wartosc {
			wykaz[miejsce+1] = wykaz[miejsce]
			miejsce--
		}
		wykaz[miejsce+1] = wartosc
	}
}

// obrysObszaruDesignu obchodzi obszar po jego GRANICY i oddaje wielokąty
// obrysu — po jednym na każdy zamknięty kontur obszaru.
//
// ── Dlaczego obchodzenie granicy, a nie wykaz prostokątów ───────────────────
// Obszar dałoby się zapisać jako wykaz prostokątów o wysokości jednego punktu
// i tak było zapisywany wcześniej. Zapis jest wtedy prawdziwy, ale nie jest
// KSZTAŁTEM: nie ma w nim naroży, więc nie ma czego wygładzić, a pole
// `smoothing` żądania nie miałoby na czym pracować. Obrys granicą daje wielokąt
// o narożach — ten sam kształt punkt w punkt, a przy tym poddający się
// zaokrągleniu.
//
// ── Kierunek obchodzenia rozstrzyga o dziurach ──────────────────────────────
// Krawędzie składa się tak, że obszar zostaje po LEWEJ stronie kierunku marszu.
// Kontur zewnętrzny wychodzi wtedy zgodnie z ruchem wskazówek zegara, a dziura
// w obszarze — przeciwnie. Dzięki temu domyślna reguła wypełniania SVG (niezerowa)
// wycina dziury sama, bez wskazywania jej w atrybucie.
//
// ── Naroże, w którym schodzą się dwie krawędzie po skosie ───────────────────
// Obszar jest spójny czterokierunkowo, ale dwa jego ramiona mogą stykać się
// narożem. W takim wierzchołku wychodzą dwie krawędzie i wybór między nimi
// rozstrzyga o tym, czy obrys się nie przecina. Marsz trzyma się ściany:
// najpierw skręt w prawo, potem prosto, potem w lewo — reguła znana i dająca
// kontury nieprzecinające się.
// ── Rachunek idzie po PUNKTACH obszaru, nie po jego prostokącie otaczającym ──
// Obszar wężowaty — ukośna kreska przez cały obraz — ma prostokąt otaczający
// wielkości obrazu i kilkaset punktów. Przejście po prostokącie kosztowałoby przy
// każdym takim obszarze tyle, ile cały obraz, a obrysowanie liczy do dwóch
// tysięcy obszarów. Dlatego przynależność punktu rozstrzygają ODCINKI (wyszukanie
// połówkowe w wierszu), a nie tablica wielkości prostokąta.
func obrysObszaruDesignu(odcinki [][3]int) [][]image.Point {
	if len(odcinki) == 0 {
		return nil
	}
	// Odcinki w wierszach: (xOd, xDo), w kolejności rosnącej — tak je oddaje
	// przejście po obszarze.
	wiersze := map[int][][2]int{}
	for _, odcinek := range odcinki {
		wiersze[odcinek[0]] = append(wiersze[odcinek[0]], [2]int{odcinek[1], odcinek[2]})
	}
	nalezy := func(x, y int) bool {
		wykaz := wiersze[y]
		od, do := 0, len(wykaz)-1
		for od <= do {
			srodek := (od + do) / 2
			switch {
			case x < wykaz[srodek][0]:
				do = srodek - 1
			case x > wykaz[srodek][1]:
				od = srodek + 1
			default:
				return true
			}
		}
		return false
	}

	// Krawędzie granicy: dla każdego punktu obszaru te jego boki, za którymi
	// obszaru już nie ma. Współrzędne są NAROŻAMI punktów, nie punktami — bok
	// punktu (0;0) od góry biegnie od naroża (0;0) do naroża (1;0).
	type krawedzObrysuDesignu struct {
		poczatek image.Point
		kierunek image.Point
	}
	wychodzace := map[image.Point][]image.Point{}
	// Kolejność zakładania krawędzi jest kolejnością wierszy i punktów w wierszu,
	// więc jest ta sama przy każdym wywołaniu. Od niej zależy kolejność konturów
	// w pliku: przejście po mapie dawałoby dwa różne pliki z jednego obrazu.
	kolejnosc := []krawedzObrysuDesignu{}
	dolozKrawedz := func(poczatek, kierunek image.Point) {
		wychodzace[poczatek] = append(wychodzace[poczatek], kierunek)
		kolejnosc = append(kolejnosc, krawedzObrysuDesignu{poczatek, kierunek})
	}
	for _, odcinek := range odcinki {
		y := odcinek[0]
		for x := odcinek[1]; x <= odcinek[2]; x++ {
			if !nalezy(x, y-1) {
				dolozKrawedz(image.Pt(x, y), image.Pt(1, 0))
			}
			if !nalezy(x+1, y) {
				dolozKrawedz(image.Pt(x+1, y), image.Pt(0, 1))
			}
			if !nalezy(x, y+1) {
				dolozKrawedz(image.Pt(x+1, y+1), image.Pt(-1, 0))
			}
			if !nalezy(x-1, y) {
				dolozKrawedz(image.Pt(x, y+1), image.Pt(0, -1))
			}
		}
	}

	zuzyte := map[krawedzObrysuDesignu]bool{}
	wPrawo := func(kierunek image.Point) image.Point { return image.Pt(-kierunek.Y, kierunek.X) }
	wLewo := func(kierunek image.Point) image.Point { return image.Pt(kierunek.Y, -kierunek.X) }
	wolna := func(punkt, kierunek image.Point) bool {
		for _, wyjscie := range wychodzace[punkt] {
			if wyjscie == kierunek && !zuzyte[krawedzObrysuDesignu{punkt, kierunek}] {
				return true
			}
		}
		return false
	}

	kontury := [][]image.Point{}
	for _, pierwsza := range kolejnosc {
		if zuzyte[pierwsza] {
			continue
		}
		poczatek := pierwsza.poczatek
		wierzcholki := []image.Point{poczatek}
		biezacy, kierunek := poczatek, pierwsza.kierunek
		for {
			zuzyte[krawedzObrysuDesignu{biezacy, kierunek}] = true
			biezacy = biezacy.Add(kierunek)
			if biezacy == poczatek {
				break
			}
			// Marsz trzymający się ściany (nagłówek funkcji).
			nastepny := image.Point{}
			znaleziony := false
			for _, proba := range []image.Point{wPrawo(kierunek), kierunek, wLewo(kierunek)} {
				if wolna(biezacy, proba) {
					nastepny, znaleziony = proba, true
					break
				}
			}
			if !znaleziony {
				break
			}
			if nastepny != kierunek {
				// Wierzchołek zapisuje się dopiero przy ZMIANIE kierunku: prosty
				// odcinek złożony z dziesięciu krawędzi jest jednym bokiem wielokąta,
				// nie dziesięcioma.
				wierzcholki = append(wierzcholki, biezacy)
			}
			kierunek = nastepny
		}
		if len(wierzcholki) >= 3 {
			kontury = append(kontury, wierzcholki)
		}
	}
	return kontury
}

// zapisSvgObrysuDesignu składa treść `d` ścieżki z wielokątów obrysu.
//
// Wygładzenie ZAOKRĄGLA naroża: bok skraca się z obu stron naroża o promień
// wygładzenia, a łuk między skróconymi końcami idzie krzywą kwadratową, której
// punktem sterującym jest samo naroże. Promień przycina się do połowy krótszego
// z boków schodzących się w narożu — inaczej dwa sąsiednie zaokrąglenia zjadłyby
// ten sam bok i kształt zawinąłby się na siebie.
//
// Wygładzenie zerowe daje wielokąt bez krzywych, czyli obrys punkt w punkt taki,
// jakie są krawędzie punktów obszaru. To rozstrzygnięcie: obrysowanie bez
// wskazania wygładzenia ma być POMIAREM obrazu, a nie kształtem upiększonym
// o liczbę, której nikt nie podał.
func zapisSvgObrysuDesignu(kontury [][]image.Point, wygladzenie float64) string {
	var zapis strings.Builder
	for _, wierzcholki := range kontury {
		if len(wierzcholki) < 3 {
			continue
		}
		if wygladzenie <= 0 {
			fmt.Fprintf(&zapis, "M%d %d", wierzcholki[0].X, wierzcholki[0].Y)
			for _, wierzcholek := range wierzcholki[1:] {
				fmt.Fprintf(&zapis, "L%d %d", wierzcholek.X, wierzcholek.Y)
			}
			zapis.WriteString("z")
			continue
		}
		zapiszZaokraglonyKonturDesignu(&zapis, wierzcholki, wygladzenie)
	}
	return zapis.String()
}

// zapiszZaokraglonyKonturDesignu dopisuje jeden kontur z zaokrąglonymi narożami.
func zapiszZaokraglonyKonturDesignu(zapis *strings.Builder, wierzcholki []image.Point,
	wygladzenie float64) {

	ile := len(wierzcholki)
	// Wejście i wyjście każdego naroża: punkty na bokach, odsunięte od naroża
	// o promień przycięty do połowy boku.
	wejscia := make([][2]float64, ile)
	wyjscia := make([][2]float64, ile)
	for numer := 0; numer < ile; numer++ {
		poprzedni := wierzcholki[(numer-1+ile)%ile]
		naroze := wierzcholki[numer]
		nastepny := wierzcholki[(numer+1)%ile]
		wejscia[numer] = punktNaBokuDesignu(naroze, poprzedni, wygladzenie)
		wyjscia[numer] = punktNaBokuDesignu(naroze, nastepny, wygladzenie)
	}
	fmt.Fprintf(zapis, "M%s %s", liczbaObrysuDesignu(wejscia[0][0]),
		liczbaObrysuDesignu(wejscia[0][1]))
	for numer := 0; numer < ile; numer++ {
		naroze := wierzcholki[numer]
		if numer > 0 {
			fmt.Fprintf(zapis, "L%s %s", liczbaObrysuDesignu(wejscia[numer][0]),
				liczbaObrysuDesignu(wejscia[numer][1]))
		}
		fmt.Fprintf(zapis, "Q%d %d %s %s", naroze.X, naroze.Y,
			liczbaObrysuDesignu(wyjscia[numer][0]), liczbaObrysuDesignu(wyjscia[numer][1]))
	}
	zapis.WriteString("z")
}

// punktNaBokuDesignu oddaje punkt leżący na boku od naroża w stronę sąsiada,
// odsunięty od naroża o promień wygładzenia przycięty do połowy boku.
func punktNaBokuDesignu(naroze, sasiad image.Point, wygladzenie float64) [2]float64 {
	roznicaX := float64(sasiad.X - naroze.X)
	roznicaY := float64(sasiad.Y - naroze.Y)
	dlugosc := math.Hypot(roznicaX, roznicaY)
	if dlugosc <= 0 {
		return [2]float64{float64(naroze.X), float64(naroze.Y)}
	}
	promien := wygladzenie
	if promien > dlugosc/2 {
		promien = dlugosc / 2
	}
	udzial := promien / dlugosc
	return [2]float64{
		float64(naroze.X) + roznicaX*udzial,
		float64(naroze.Y) + roznicaY*udzial,
	}
}

// liczbaObrysuDesignu zapisuje współrzędną obrysu zaokrągloną do setnej części
// punktu obrazu i bez zer na końcu. Poniżej setnej nie ma czego zapisywać,
// a wygładzenie podane jako jedna trzecia rozdmuchałoby plik o kilkanaście cyfr
// przy każdej współrzędnej.
func liczbaObrysuDesignu(wartosc float64) string {
	return strconv.FormatFloat(math.Round(wartosc*100)/100, 'f', -1, 64)
}
