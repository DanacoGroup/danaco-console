// Odpowiedzialność pliku: `image.vectorize` — zamiana bitmapy na czyste ścieżki
// wektorowe. Arytmetyka obrysu, upraszczania i ścieńczania leży
// w `adapter_narzedzia_obraz_wektor_slad.go`; tutaj jest rozstrzygnięcie, co
// właściwie obrysowujemy w każdym z trzech trybów, i złożenie wyniku w SVG.
//
// ── Trzy tryby, trzy różne pytania ──────────────────────────────────────────
//   - `outline`   — gdzie kończy się kształt. Obraz sprowadzamy do dwóch
//     wartości progiem jasności i obrysowujemy obszar ciemny. To tryb dla
//     logotypu, pieczęci, znaku.
//   - `posterize` — z ilu płaszczyzn barwnych składa się obraz. Barwy skupiamy
//     w tyle grup, ile mówi `colors`, i każdą grupę obrysowujemy osobno. To
//     tryb dla ilustracji.
//   - `centerline`— którędy biegnie kreska. Obszar ścieńczamy do linii
//     o grubości piksela i obrysowujemy ją jako kreskę. To tryb dla rysunku
//     technicznego i pisma odręcznego, gdzie obrys konturu dałby każdą kreskę
//     jako podwójną pętlę.
//
// ── Wynik jest zasobem SVG, nie obrazem rastrowym ──────────────────────────
// Zasób idzie do tego samego magazynu, co każdy inny wytwór rodziny `image.*`,
// ale jego format to `svg`. Wymiarów rdzeń przy nim nie mierzy: `image.DecodeConfig`
// nie zna SVG, a wpisanie tam wymiarów źródła podałoby liczby, których nikt nie
// zmierzył na wyniku. Wymiary niosą atrybuty samego dokumentu SVG.
package core

import (
	"context"
	"fmt"
	"image"
	"math"
	"sort"
	"strings"

	"danacoconsole/shared"
)

const (
	// granicaPikseliWektoryzacji — obrys chodzi po każdym pikselu i po każdym
	// jego sąsiedzie, więc koszt rośnie liniowo z polem, ale pamięć maski jest
	// dodatkowa. Szesnaście megapikseli to zdjęcie 4000×4000: więcej nie ma
	// sensu obrysowywać, bo wynik ma wtedy więcej wierzchołków niż źródło
	// pikseli.
	granicaPikseliWektoryzacji = 16_000_000

	// najmniejszyObrysowywanyObszar odrzuca kontury krótsze niż pięć pikseli.
	// Bez tego pojedynczy piksel szumu wchodzi do wyniku jako osobna ścieżka,
	// a skan z aparatu daje ich dziesiątki tysięcy.
	najmniejszyObrysowywanyObszar = 5

	// gornaGranicaBarwWektoryzacji zamyka liczbę płaszczyzn trybu `posterize`.
	// Powyżej kilkunastu barw wynik przestaje być wektorem, a staje się mapą
	// pikseli zapisaną jako ścieżki — cięższą od źródła i nie do edycji.
	gornaGranicaBarwWektoryzacji = 24

	// progJasnosciObrysu rozdziela obraz na ciemny (obrysowywany) i jasny (tło)
	// w trybach `outline` i `centerline`. Połowa zakresu jest wyborem
	// neutralnym; kontrakt nie ma pola na próg, a zgadywanie go z histogramu
	// dawałoby dwa różne wyniki dla dwóch skanów tego samego rysunku.
	progJasnosciObrysu = 0.5
)

// Zwektoryzuj obsługuje `image.vectorize`.
func (a *adapterNarzedziObrazu) Zwektoryzuj(ctx context.Context,
	z shared.ImageVectorizeRequest) (shared.ImageVectorizeResponse, error) {

	zrodlo, err := a.rozwiazZrodlo(ctx, z.AssetId, z.SourcePath)
	if err != nil {
		return shared.ImageVectorizeResponse{}, err
	}
	obraz, err := odczytajObrazPliku(zrodlo.sciezka)
	if err != nil {
		return shared.ImageVectorizeResponse{}, bladWskazaniaObrazu(
			"nie można odczytać obrazu do wektoryzacji: " + err.Error())
	}
	granice := obraz.Bounds()
	if granice.Dx() <= 0 || granice.Dy() <= 0 {
		return shared.ImageVectorizeResponse{}, bladWskazaniaObrazu(
			"obraz nie ma ani jednego piksela — nie ma czego obrysować")
	}
	if granice.Dx()*granice.Dy() > granicaPikseliWektoryzacji {
		return shared.ImageVectorizeResponse{}, bladWskazaniaObrazu(
			"obraz jest za duży do obrysowania (ponad szesnaście megapikseli); " +
				"naprawa: pomniejszyć go wcześniej komendą image.transform")
	}

	var tryb shared.ImageVectorizeMode = shared.ImageVectorizeModeOutline
	if z.Mode != nil && strings.TrimSpace(string(*z.Mode)) != "" {
		tryb = *z.Mode
	}
	tolerancja := tolerancjaUpraszczania(z.Simplify)

	var sciezki []sciezkaWektorowa
	switch tryb {
	case shared.ImageVectorizeModeOutline:
		sciezki = obrysKonturu(obraz, tolerancja)
	case shared.ImageVectorizeModeCenterline:
		sciezki = obrysLiniiSrodkowej(obraz, tolerancja)
	case shared.ImageVectorizeModePosterize:
		sciezki = obrysPlaszczyznBarwnych(obraz, liczbaBarwWektoryzacji(z.Colors), tolerancja)
	default:
		return shared.ImageVectorizeResponse{}, bladWskazaniaObrazu(
			"nie znam sposobu obrysu „" + string(tryb) +
				"” — rdzeń zna: outline, centerline, posterize")
	}

	if len(sciezki) == 0 {
		// Obraz jednolity nie ma konturu i to jest odpowiedź, nie awaria. Ale
		// pusty dokument SVG podany jako wektoryzacja byłby atrapą wyniku,
		// więc odmawiamy, nazywając powód.
		return shared.ImageVectorizeResponse{}, bladPrzetwarzaniaObrazu(
			"obrys nie znalazł ani jednej ścieżki — obraz jest jednolity albo " +
				"kontrast między kształtem a tłem jest zbyt mały do rozdzielenia")
	}

	dokument := zlozDokumentSvg(granice.Dx(), granice.Dy(), sciezki)
	zasob, _, err := a.odlozZasob(ctx, zrodlo, z.WindowId, []byte(dokument), "svg", "wektoryzacja")
	if err != nil {
		return shared.ImageVectorizeResponse{}, err
	}
	liczba := len(sciezki)
	return shared.ImageVectorizeResponse{Asset: zasob, Paths: &liczba}, nil
}

// sciezkaWektorowa jest jedną ścieżką wyniku wraz z jej rolą: wypełnieniem albo
// kreską. Rozdział jest istotny — linia środkowa nie ma wnętrza, więc wypełniona
// wyszłaby jako czarna plama.
type sciezkaWektorowa struct {
	punkty    []punktSladu
	barwa     string
	zamknieta bool
	kreska    bool
}

// tolerancjaUpraszczania przekłada procent kontraktu na odchylenie w pikselach.
//
// Sto procent to odchylenie o osiem pikseli — tyle, że z litery zostaje
// czworobok. Zero znaczy „nie upraszczaj": łamana zostaje taka, jak wyszła
// z obrysu.
func tolerancjaUpraszczania(procent *int) float64 {
	if procent == nil || *procent <= 0 {
		return 0
	}
	wartosc := *procent
	if wartosc > 100 {
		wartosc = 100
	}
	return float64(wartosc) * 8 / 100
}

// liczbaBarwWektoryzacji rozstrzyga, w ile płaszczyzn skupić barwy. Brak pola
// znaczy obrys jednobarwny — tak mówi kontrakt.
func liczbaBarwWektoryzacji(zadane *int) int {
	if zadane == nil || *zadane <= 1 {
		return 1
	}
	if *zadane > gornaGranicaBarwWektoryzacji {
		return gornaGranicaBarwWektoryzacji
	}
	return *zadane
}

// maskaCiemnychPikseli sprowadza obraz do dwóch wartości progiem jasności.
// Piksel przezroczysty należy do tła: obraz z wyciętym tłem ma kształt zapisany
// w kanale alfa, a nie w jasności.
func maskaCiemnychPikseli(obraz image.Image) *maskaRastrowa {
	granice := obraz.Bounds()
	maska := nowaMaskaRastrowa(granice.Dx(), granice.Dy())
	for y := 0; y < granice.Dy(); y++ {
		for x := 0; x < granice.Dx(); x++ {
			r, g, b, a := skladoweNiePrzemnozone(obraz.At(granice.Min.X+x, granice.Min.Y+y))
			if a < 0.5 {
				continue
			}
			jasnosc := 0.2126*r + 0.7152*g + 0.0722*b
			maska.ustaw(x, y, jasnosc < progJasnosciObrysu)
		}
	}
	return maska
}

// obrysKonturu obrysowuje kształt ciemny — tryb `outline`.
func obrysKonturu(obraz image.Image, tolerancja float64) []sciezkaWektorowa {
	kontury := obrysyMaski(maskaCiemnychPikseli(obraz), najmniejszyObrysowywanyObszar)
	sciezki := make([]sciezkaWektorowa, 0, len(kontury))
	for _, kontur := range kontury {
		uproszczony := uproscLamana(kontur, tolerancja)
		if len(uproszczony) < 3 {
			continue
		}
		sciezki = append(sciezki, sciezkaWektorowa{
			punkty: uproszczony, barwa: "#000000", zamknieta: true,
		})
	}
	return sciezki
}

// obrysLiniiSrodkowej ścieńcza kształt do linii i obrysowuje ją jako kreskę —
// tryb `centerline`.
func obrysLiniiSrodkowej(obraz image.Image, tolerancja float64) []sciezkaWektorowa {
	scienczona := scienczMaske(maskaCiemnychPikseli(obraz))
	kontury := obrysyMaski(scienczona, najmniejszyObrysowywanyObszar)
	sciezki := make([]sciezkaWektorowa, 0, len(kontury))
	for _, kontur := range kontury {
		uproszczony := uproscLamana(kontur, tolerancja)
		if len(uproszczony) < 2 {
			continue
		}
		sciezki = append(sciezki, sciezkaWektorowa{
			punkty: uproszczony, barwa: "#000000", kreska: true,
		})
	}
	return sciezki
}

// obrysPlaszczyznBarwnych skupia barwy obrazu i obrysowuje każdą płaszczyznę
// osobno — tryb `posterize`.
//
// Barwy skupiamy metodą k-średnich na próbce pikseli. Próbka, a nie komplet:
// dla obrazu megapikselowego przejście po wszystkich pikselach w każdej
// iteracji kosztuje sekundy, a środki skupień z próbki co dziesiąty piksel
// wychodzą praktycznie takie same.
//
// Płaszczyzny idą od najciemniejszej: w dokumencie SVG ścieżka późniejsza
// zasłania wcześniejszą, a płaszczyzna jasna bywa tłem dla ciemnej.
func obrysPlaszczyznBarwnych(obraz image.Image, barw int, tolerancja float64) []sciezkaWektorowa {
	if barw <= 1 {
		return obrysKonturu(obraz, tolerancja)
	}
	granice := obraz.Bounds()
	srodki := skupieniaBarw(obraz, barw)
	if len(srodki) == 0 {
		return nil
	}
	sort.Slice(srodki, func(i, j int) bool {
		return jasnoscSkupienia(srodki[i]) < jasnoscSkupienia(srodki[j])
	})

	sciezki := []sciezkaWektorowa{}
	for _, srodek := range srodki {
		maska := nowaMaskaRastrowa(granice.Dx(), granice.Dy())
		for y := 0; y < granice.Dy(); y++ {
			for x := 0; x < granice.Dx(); x++ {
				r, g, b, a := skladoweNiePrzemnozone(obraz.At(granice.Min.X+x, granice.Min.Y+y))
				if a < 0.5 {
					continue
				}
				if najblizszeSkupienie([3]float64{r, g, b}, srodki) == srodek {
					maska.ustaw(x, y, true)
				}
			}
		}
		if maska.liczbaPol() == 0 {
			continue
		}
		for _, kontur := range obrysyMaski(maska, najmniejszyObrysowywanyObszar) {
			uproszczony := uproscLamana(kontur, tolerancja)
			if len(uproszczony) < 3 {
				continue
			}
			sciezki = append(sciezki, sciezkaWektorowa{
				punkty: uproszczony, barwa: zapisBarwy(srodek), zamknieta: true,
			})
		}
	}
	return sciezki
}

// skupieniaBarw wyznacza środki skupień barw metodą k-średnich.
//
// Środki startowe rozkładamy równomiernie po osi jasności próbki, a nie losowo:
// losowy start dawałby dwa różne wyniki dla dwóch wywołań na tym samym obrazie,
// a wektoryzacja ma być powtarzalna.
func skupieniaBarw(obraz image.Image, barw int) [][3]float64 {
	granice := obraz.Bounds()
	krok := 1
	if granice.Dx()*granice.Dy() > 200_000 {
		krok = 1 + (granice.Dx()*granice.Dy())/200_000
	}

	probka := make([][3]float64, 0, 200_000)
	for y := 0; y < granice.Dy(); y++ {
		for x := 0; x < granice.Dx(); x += krok {
			r, g, b, a := skladoweNiePrzemnozone(obraz.At(granice.Min.X+x, granice.Min.Y+y))
			if a < 0.5 {
				continue
			}
			probka = append(probka, [3]float64{r, g, b})
		}
	}
	if len(probka) == 0 {
		return nil
	}
	if barw > len(probka) {
		barw = len(probka)
	}

	posortowana := make([][3]float64, len(probka))
	copy(posortowana, probka)
	sort.Slice(posortowana, func(i, j int) bool {
		return jasnoscSkupienia(posortowana[i]) < jasnoscSkupienia(posortowana[j])
	})
	srodki := make([][3]float64, barw)
	for numer := range srodki {
		srodki[numer] = posortowana[numer*len(posortowana)/barw]
	}

	for iteracja := 0; iteracja < 12; iteracja++ {
		sumy := make([][3]float64, barw)
		liczby := make([]int, barw)
		for _, piksel := range probka {
			numer := numerNajblizszegoSkupienia(piksel, srodki)
			sumy[numer][0] += piksel[0]
			sumy[numer][1] += piksel[1]
			sumy[numer][2] += piksel[2]
			liczby[numer]++
		}
		ruch := 0.0
		for numer := range srodki {
			if liczby[numer] == 0 {
				continue
			}
			nowy := [3]float64{
				sumy[numer][0] / float64(liczby[numer]),
				sumy[numer][1] / float64(liczby[numer]),
				sumy[numer][2] / float64(liczby[numer]),
			}
			ruch += math.Abs(nowy[0]-srodki[numer][0]) + math.Abs(nowy[1]-srodki[numer][1]) +
				math.Abs(nowy[2]-srodki[numer][2])
			srodki[numer] = nowy
		}
		if ruch < 0.001 {
			break
		}
	}
	return srodki
}

// numerNajblizszegoSkupienia zwraca indeks środka najbliższego pikselowi.
func numerNajblizszegoSkupienia(piksel [3]float64, srodki [][3]float64) int {
	najlepszy := 0
	najmniejsza := math.MaxFloat64
	for numer, srodek := range srodki {
		odleglosc := (piksel[0]-srodek[0])*(piksel[0]-srodek[0]) +
			(piksel[1]-srodek[1])*(piksel[1]-srodek[1]) +
			(piksel[2]-srodek[2])*(piksel[2]-srodek[2])
		if odleglosc < najmniejsza {
			najmniejsza = odleglosc
			najlepszy = numer
		}
	}
	return najlepszy
}

// najblizszeSkupienie zwraca środek najbliższy pikselowi.
func najblizszeSkupienie(piksel [3]float64, srodki [][3]float64) [3]float64 {
	return srodki[numerNajblizszegoSkupienia(piksel, srodki)]
}

// jasnoscSkupienia liczy jasność środka skupienia wagami luminancji.
func jasnoscSkupienia(srodek [3]float64) float64 {
	return 0.2126*srodek[0] + 0.7152*srodek[1] + 0.0722*srodek[2]
}

// zapisBarwy zamienia środek skupienia na zapis szesnastkowy dokumentu SVG.
func zapisBarwy(srodek [3]float64) string {
	return fmt.Sprintf("#%02x%02x%02x", bajtSkladowej(srodek[0]), bajtSkladowej(srodek[1]),
		bajtSkladowej(srodek[2]))
}

// zlozDokumentSvg składa gotowy dokument z wykazu ścieżek.
//
// `shape-rendering="geometricPrecision"` mówi przeglądarce, żeby nie
// zaokrąglała wierzchołków do siatki pikseli — bez tego uproszczona łamana
// wygląda w podglądzie na bardziej postrzępioną, niż jest.
func zlozDokumentSvg(szerokosc, wysokosc int, sciezki []sciezkaWektorowa) string {
	var dokument strings.Builder
	fmt.Fprintf(&dokument,
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" `+
			`shape-rendering="geometricPrecision">`+"\n", szerokosc, wysokosc, szerokosc, wysokosc)
	for _, sciezka := range sciezki {
		dane := zapisSciezki(sciezka)
		if dane == "" {
			continue
		}
		if sciezka.kreska {
			fmt.Fprintf(&dokument,
				`  <path d="%s" fill="none" stroke="%s" stroke-width="1" `+
					`stroke-linecap="round" stroke-linejoin="round"/>`+"\n", dane, sciezka.barwa)
			continue
		}
		fmt.Fprintf(&dokument, `  <path d="%s" fill="%s"/>`+"\n", dane, sciezka.barwa)
	}
	dokument.WriteString("</svg>\n")
	return dokument.String()
}

// zapisSciezki składa atrybut `d` jednej ścieżki.
func zapisSciezki(sciezka sciezkaWektorowa) string {
	if len(sciezka.punkty) < 2 {
		return ""
	}
	var zapis strings.Builder
	for numer, punkt := range sciezka.punkty {
		polecenie := "L"
		if numer == 0 {
			polecenie = "M"
		} else {
			zapis.WriteString(" ")
		}
		fmt.Fprintf(&zapis, "%s%s %s", polecenie, zapisWspolrzednej(punkt.X),
			zapisWspolrzednej(punkt.Y))
	}
	if sciezka.zamknieta {
		zapis.WriteString(" Z")
	}
	return zapis.String()
}

// zapisWspolrzednej zapisuje liczbę bez zbędnych zer — dokument SVG z setkami
// tysięcy wierzchołków rośnie o megabajty na samych ogonach dziesiętnych.
func zapisWspolrzednej(wartosc float64) string {
	return strings.TrimSuffix(strings.TrimRight(fmt.Sprintf("%.2f", wartosc), "0"), ".")
}
