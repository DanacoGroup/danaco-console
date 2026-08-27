// Odpowiedzialność pliku: zamiana maski rastrowej na ścieżki — obrys konturu,
// upraszczanie łamanej i ścieńczanie do linii środkowej, czystą arytmetyką na
// tablicy pikseli, bez programu zewnętrznego typu `potrace`.
package core

import (
	"math"
)

// maskaRastrowa jest tablicą logiczną o wymiarach obrazu: prawda znaczy piksel
// należący do obrysowywanego obszaru.
type maskaRastrowa struct {
	szerokosc int
	wysokosc  int
	pola      []bool
}

// nowaMaskaRastrowa zakłada pustą maskę o zadanych wymiarach, gotową do
// wypełnienia przynależnością pikseli przed śledzeniem konturu.
func nowaMaskaRastrowa(szerokosc, wysokosc int) *maskaRastrowa {
	return &maskaRastrowa{
		szerokosc: szerokosc,
		wysokosc:  wysokosc,
		pola:      make([]bool, szerokosc*wysokosc),
	}
}

// wewnatrz mówi, czy piksel należy do obszaru. Współrzędna poza obrazem daje
// fałsz, dzięki czemu obrys nie musi osobno pilnować krawędzi.
func (m *maskaRastrowa) wewnatrz(x, y int) bool {
	if x < 0 || y < 0 || x >= m.szerokosc || y >= m.wysokosc {
		return false
	}
	return m.pola[y*m.szerokosc+x]
}

// ustaw wpisuje przynależność piksela do obszaru, milcząc na współrzędnej
// spoza obrazu zamiast wywracać się paniką indeksu.
func (m *maskaRastrowa) ustaw(x, y int, wartosc bool) {
	if x < 0 || y < 0 || x >= m.szerokosc || y >= m.wysokosc {
		return
	}
	m.pola[y*m.szerokosc+x] = wartosc
}

// liczbaPol zwraca liczbę pikseli należących do obszaru, do porównań progu
// wielkości bez osobnego przechodzenia po całej masce.
func (m *maskaRastrowa) liczbaPol() int {
	liczba := 0
	for _, pole := range m.pola {
		if pole {
			liczba++
		}
	}
	return liczba
}

// punktSladu jest wierzchołkiem łamanej w układzie obrazu, we współrzędnych
// zmiennoprzecinkowych, żeby uproszczenie mogło przesuwać punkty swobodnie.
type punktSladu struct {
	X float64
	Y float64
}

// kierunkiObrysu to osiem kierunków sąsiedztwa w kolejności zgodnej z ruchem
// wskazówek zegara, zaczynając od wschodu — kolejność, w jakiej metoda Moore'a
// obchodzi sąsiadów.
var kierunkiObrysu = [8][2]int{
	{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1},
}

// obrysyMaski zwraca kontury wszystkich spójnych obszarów maski, metodą
// śledzenia sąsiedztwa Moore'a, pomijając obszary mniejsze niż
// `najmniejszyObszar`.
func obrysyMaski(maska *maskaRastrowa, najmniejszyObszar int) [][]punktSladu {
	if maska == nil {
		return nil
	}
	odwiedzone := nowaMaskaRastrowa(maska.szerokosc, maska.wysokosc)
	kontury := [][]punktSladu{}

	for y := 0; y < maska.wysokosc; y++ {
		for x := 0; x < maska.szerokosc; x++ {
			if !maska.wewnatrz(x, y) || odwiedzone.wewnatrz(x, y) {
				continue
			}
			// Piksel wewnętrzny (otoczony ze wszystkich stron) nie zaczyna konturu.

			// Kontur zaczyna się na brzegu obszaru.
			if maska.wewnatrz(x-1, y) && maska.wewnatrz(x+1, y) &&
				maska.wewnatrz(x, y-1) && maska.wewnatrz(x, y+1) {
				continue
			}
			kontur := obejdzObszar(maska, odwiedzone, x, y)
			if len(kontur) >= 3 && len(kontur) >= najmniejszyObszar {
				kontury = append(kontury, kontur)
			}
		}
	}
	return kontury
}

// obejdzObszar prowadzi jedno obejście konturu od wskazanego piksela
// brzegowego. Granica liczby kroków chroni przed konturem, który się nie
// domyka: bez niej pętla zawiesiłaby żądanie na zawsze.
func obejdzObszar(maska, odwiedzone *maskaRastrowa, startX, startY int) []punktSladu {
	kontur := []punktSladu{}
	biezacyX, biezacyY := startX, startY
	kierunek := 0
	maksimum := maska.szerokosc*maska.wysokosc*4 + 16

	for krok := 0; krok < maksimum; krok++ {
		kontur = append(kontur, punktSladu{X: float64(biezacyX), Y: float64(biezacyY)})
		odwiedzone.ustaw(biezacyX, biezacyY, true)

		// Następny piksel obszaru szuka się od kierunku „w tył i w lewo".

		// Tak, żeby obejście trzymało się brzegu obszaru.
		znaleziono := false
		poczatek := (kierunek + 6) % 8
		for obrot := 0; obrot < 8; obrot++ {
			sprawdzany := (poczatek + obrot) % 8
			nastepnyX := biezacyX + kierunkiObrysu[sprawdzany][0]
			nastepnyY := biezacyY + kierunkiObrysu[sprawdzany][1]
			if maska.wewnatrz(nastepnyX, nastepnyY) {
				biezacyX, biezacyY = nastepnyX, nastepnyY
				kierunek = sprawdzany
				znaleziono = true
				break
			}
		}
		if !znaleziono {
			// Piksel samotny — kontur jest nim samym.
			break
		}
		if biezacyX == startX && biezacyY == startY {
			break
		}
	}
	return kontur
}

// uproscLamana odrzuca wierzchołki, których usunięcie odchyla łamaną mniej niż
// o `tolerancja` — algorytm Ramera–Douglasa–Peuckera.
//
// Tolerancja zero zwraca łamaną nietkniętą: „nie upraszczaj" ma znaczyć
// dokładnie to, a nie „upraszczaj minimalnie".
func uproscLamana(punkty []punktSladu, tolerancja float64) []punktSladu {
	if len(punkty) < 3 || tolerancja <= 0 {
		return punkty
	}
	zachowane := make([]bool, len(punkty))
	zachowane[0] = true
	zachowane[len(punkty)-1] = true
	uproscOdcinek(punkty, 0, len(punkty)-1, tolerancja, zachowane)

	wynik := make([]punktSladu, 0, len(punkty))
	for numer, punkt := range punkty {
		if zachowane[numer] {
			wynik = append(wynik, punkt)
		}
	}
	return wynik
}

// uproscOdcinek zaznacza wierzchołki, które muszą zostać, żeby odchylenie
// łamanej nie przekroczyło tolerancji.
func uproscOdcinek(punkty []punktSladu, poczatek, koniec int, tolerancja float64, zachowane []bool) {
	if koniec <= poczatek+1 {
		return
	}
	najdalszy := -1
	najwieksza := 0.0
	for numer := poczatek + 1; numer < koniec; numer++ {
		odleglosc := odlegloscOdOdcinka(punkty[numer], punkty[poczatek], punkty[koniec])
		if odleglosc > najwieksza {
			najwieksza = odleglosc
			najdalszy = numer
		}
	}
	if najdalszy < 0 || najwieksza <= tolerancja {
		return
	}
	zachowane[najdalszy] = true
	uproscOdcinek(punkty, poczatek, najdalszy, tolerancja, zachowane)
	uproscOdcinek(punkty, najdalszy, koniec, tolerancja, zachowane)
}

// odlegloscOdOdcinka liczy odległość punktu od odcinka. Odcinek zdegenerowany
// do punktu daje odległość od tego punktu, zamiast dzielenia przez zero na
// konturze zamkniętym w jednym pikselu.
func odlegloscOdOdcinka(punkt, poczatek, koniec punktSladu) float64 {
	dx := koniec.X - poczatek.X
	dy := koniec.Y - poczatek.Y
	if dx == 0 && dy == 0 {
		return math.Hypot(punkt.X-poczatek.X, punkt.Y-poczatek.Y)
	}
	udzial := ((punkt.X-poczatek.X)*dx + (punkt.Y-poczatek.Y)*dy) / (dx*dx + dy*dy)
	switch {
	case udzial < 0:
		udzial = 0
	case udzial > 1:
		udzial = 1
	}
	rzutX := poczatek.X + udzial*dx
	rzutY := poczatek.Y + udzial*dy
	return math.Hypot(punkt.X-rzutX, punkt.Y-rzutY)
}

// scienczMaske sprowadza obszar do linii o grubości jednego piksela —
// algorytm Zhanga–Suena — i zasila obrys linii środkowej (`centerline`).
// Granica przebiegów chroni przed układem, który oscyluje.
func scienczMaske(maska *maskaRastrowa) *maskaRastrowa {
	if maska == nil {
		return nil
	}
	praca := nowaMaskaRastrowa(maska.szerokosc, maska.wysokosc)
	copy(praca.pola, maska.pola)

	for przebieg := 0; przebieg < 200; przebieg++ {
		zdjeto := false
		for podprzebieg := 0; podprzebieg < 2; podprzebieg++ {
			doZdjecia := [][2]int{}
			for y := 1; y < praca.wysokosc-1; y++ {
				for x := 1; x < praca.szerokosc-1; x++ {
					if !praca.wewnatrz(x, y) {
						continue
					}
					if czyZdejmowalny(praca, x, y, podprzebieg) {
						doZdjecia = append(doZdjecia, [2]int{x, y})
					}
				}
			}
			for _, punkt := range doZdjecia {
				praca.ustaw(punkt[0], punkt[1], false)
				zdjeto = true
			}
		}
		if !zdjeto {
			break
		}
	}
	return praca
}

// czyZdejmowalny rozstrzyga cztery warunki Zhanga–Suena dla jednego piksela,
// licząc sąsiadów w kolejności zegarowej od północy.
func czyZdejmowalny(maska *maskaRastrowa, x, y, podprzebieg int) bool {
	polnoc := maska.wewnatrz(x, y-1)
	polnocnyWschod := maska.wewnatrz(x+1, y-1)
	wschod := maska.wewnatrz(x+1, y)
	poludniowyWschod := maska.wewnatrz(x+1, y+1)
	poludnie := maska.wewnatrz(x, y+1)
	poludniowyZachod := maska.wewnatrz(x-1, y+1)
	zachod := maska.wewnatrz(x-1, y)
	polnocnyZachod := maska.wewnatrz(x-1, y-1)

	sasiedzi := []bool{polnoc, polnocnyWschod, wschod, poludniowyWschod,
		poludnie, poludniowyZachod, zachod, polnocnyZachod}

	liczba := 0
	for _, sasiad := range sasiedzi {
		if sasiad {
			liczba++
		}
	}
	if liczba < 2 || liczba > 6 {
		return false
	}

	przejscia := 0
	for numer := range sasiedzi {
		nastepny := sasiedzi[(numer+1)%len(sasiedzi)]
		if !sasiedzi[numer] && nastepny {
			przejscia++
		}
	}
	if przejscia != 1 {
		return false
	}

	if podprzebieg == 0 {
		if polnoc && wschod && poludnie {
			return false
		}
		if wschod && poludnie && zachod {
			return false
		}
		return true
	}
	if polnoc && wschod && zachod {
		return false
	}
	if polnoc && poludnie && zachod {
		return false
	}
	return true
}
