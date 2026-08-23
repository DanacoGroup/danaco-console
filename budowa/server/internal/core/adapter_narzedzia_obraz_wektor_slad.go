// Odpowiedzialność pliku: zamiana maski rastrowej na ścieżki — obrys konturu,
// upraszczanie łamanej i ścieńczanie do linii środkowej. To czysta arytmetyka
// na tablicy pikseli; wołający (`adapter_narzedzia_obraz_wektor.go`) rozstrzyga,
// skąd maska pochodzi i co z gotowymi ścieżkami zrobić.
//
// ── Dlaczego własna arytmetyka, a nie program zewnętrzny ────────────────────
// Zamiana rastra na ścieżki bywa robiona programem `potrace`. Program ten nie
// stoi na serwerze, a instalka Operatora nie niesie żadnego programu — czynność
// oparta na nim byłaby u odbiorcy odmową, a nie funkcją. Obrys konturu,
// upraszczanie Ramera–Douglasa–Peuckera i ścieńczanie Zhanga–Suena to
// algorytmy opisane i skończone; wkompilowane w binarium działają wszędzie tam,
// gdzie działa rdzeń.
//
// Wynikiem jest łamana, nie krzywa Béziera. To rozstrzygnięcie, nie brak:
// łamana po uproszczeniu opisuje kontur wiernie i przewidywalnie, a
// dopasowanie krzywych wprowadza odchylenie, którego Operator nie kontroluje
// żadnym polem kontraktu. Pole `simplify` steruje właśnie odchyleniem łamanej
// i mówi wprost, ile go wolno.
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

// nowaMaskaRastrowa zakłada pustą maskę o zadanych wymiarach.
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

// ustaw wpisuje przynależność piksela.
func (m *maskaRastrowa) ustaw(x, y int, wartosc bool) {
	if x < 0 || y < 0 || x >= m.szerokosc || y >= m.wysokosc {
		return
	}
	m.pola[y*m.szerokosc+x] = wartosc
}

// liczbaPol zwraca liczbę pikseli należących do obszaru.
func (m *maskaRastrowa) liczbaPol() int {
	liczba := 0
	for _, pole := range m.pola {
		if pole {
			liczba++
		}
	}
	return liczba
}

// punktSladu jest wierzchołkiem łamanej w układzie obrazu.
type punktSladu struct {
	X float64
	Y float64
}

// kierunkiObrysu to osiem kierunków sąsiedztwa w kolejności zgodnej z ruchem
// wskazówek zegara, zaczynając od wschodu. Kolejność ma znaczenie: śledzenie
// konturu metodą Moore'a chodzi po sąsiadach właśnie w tym porządku i to on
// rozstrzyga, że kontur zewnętrzny wychodzi zgodnie z ruchem wskazówek.
var kierunkiObrysu = [8][2]int{
	{1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}, {-1, -1}, {0, -1}, {1, -1},
}

// obrysyMaski zwraca kontury wszystkich spójnych obszarów maski.
//
// Metoda jest klasycznym śledzeniem sąsiedztwa Moore'a: znajdź piksel brzegowy,
// obejdź obszar dookoła, wróć do punktu wyjścia. Piksele już objęte konturem
// znakujemy, żeby ten sam obszar nie dał dwóch identycznych ścieżek — bez tego
// każdy piksel brzegu byłby początkiem osobnego obejścia.
//
// Obszary mniejsze niż `najmniejszyObszar` pomijamy: pojedyncze piksele szumu
// dałyby setki ścieżek o wielkości kropki, przez które wynik jest cięższy od
// źródła i nie do otwarcia w edytorze wektorowym.
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
			// Piksel wewnętrzny (otoczony ze wszystkich stron) nie zaczyna
			// konturu — kontur zaczyna się na brzegu obszaru.
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

// obejdzObszar prowadzi jedno obejście konturu od wskazanego piksela brzegowego.
//
// Granica liczby kroków chroni przed obrazem, którego kontur z jakiegoś powodu
// nie domyka się w punkcie wyjścia: pętla bez granicy zawiesiłaby żądanie na
// zawsze, a odmowa po granicy jest odpowiedzią.
func obejdzObszar(maska, odwiedzone *maskaRastrowa, startX, startY int) []punktSladu {
	kontur := []punktSladu{}
	biezacyX, biezacyY := startX, startY
	kierunek := 0
	maksimum := maska.szerokosc*maska.wysokosc*4 + 16

	for krok := 0; krok < maksimum; krok++ {
		kontur = append(kontur, punktSladu{X: float64(biezacyX), Y: float64(biezacyY)})
		odwiedzone.ustaw(biezacyX, biezacyY, true)

		// Szukamy następnego piksela obszaru, obchodząc sąsiadów od kierunku
		// „w tył i w lewo" — tak, żeby obejście trzymało się brzegu.
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
// do punktu daje odległość od tego punktu — bez tego przypadku dzielilibyśmy
// przez zero na konturze zamkniętym w jednym pikselu.
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
// algorytm Zhanga–Suena. Zasila obrys linii środkowej (`centerline`): rysunek
// kreskowy obrysowany po konturze dałby każdą kreskę jako podwójną pętlę,
// a obrysowany po linii środkowej — jako jedną kreskę.
//
// Algorytm chodzi naprzemiennie dwoma podprzebiegami, aż przestanie cokolwiek
// zdejmować. Granica przebiegów chroni przed układem, który oscyluje.
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

// czyZdejmowalny rozstrzyga warunki Zhanga–Suena dla jednego piksela.
//
// Sąsiedzi liczeni są w kolejności zegarowej od północy. Warunki są cztery:
// liczba sąsiadów mieści się w 2..6 (piksel nie jest ani końcem, ani wnętrzem),
// przejść z tła do obszaru jest dokładnie jedno (zdjęcie nie rozerwie linii),
// oraz dwie pary sąsiadów zależne od podprzebiegu — to one na przemian ścinają
// obszar z dwóch przeciwnych stron, żeby linia wyszła pośrodku, a nie przy
// jednej krawędzi.
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
