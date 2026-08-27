// Plik składa krój ikonowy TrueType (postać `webfont`) z konturów ikon katalogu —
// obsługuje `design.icon.sprite.build`. Postać `sprite-svg` składa
// `adapter_modul_design_ikony.go`, katalog wzorów leży
// w `adapter_modul_design_ikony_katalog.go`.
package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/tdewolff/canvas"
)

const (
	// jednostekNaFiretDesignu jest liczbą jednostek kroju na firet. 1024 dzieli
	// się bez reszty przez potęgi dwójki, więc skalowanie ikony z siatki 24 nie
	// wprowadza błędu narastającego.
	jednostekNaFiretDesignu = 1024

	// poczatekObszaruPrywatnegoDesignu jest pierwszym punktem kodowym obszaru
	// prywatnego Unicode (Private Use Area).
	poczatekObszaruPrywatnegoDesignu = 0xE000

	// granicaGlifowKrojuIkonowegoDesignu wynika z obszaru prywatnego pierwszego
	// poziomu: od U+E000 do U+F8FF jest 6400 miejsc. Powyżej krój musiałby wejść
	// w płaszczyzny dodatkowe, a tam `cmap` formatu 4 nie sięga.
	granicaGlifowKrojuIkonowegoDesignu = 6400

	// tolerancjaSplaszczeniaKrojuDesignu jest dopuszczalnym odejściem konturu od
	// krzywej, w jednostkach kroju. Poniżej jednostki różnicy nie widać.
	tolerancjaSplaszczeniaKrojuDesignu = 1.0

	// tolerancjaObrysuKrojuDesignu jest tolerancją liczenia obrysu kreski, podaną
	// w jednostkach siatki ikony — obrys powstaje przed przeskalowaniem do
	// jednostek kroju, na poziomie jednej czwartej jednostki kroju.
	tolerancjaObrysuKrojuDesignu = 0.25 * siatkaIkonyKataloguDesignu /
		jednostekNaFiretDesignu

	// najmniejszyOdstepPunktowKrojuDesignu jest najmniejszą odległością dwóch
	// kolejnych punktów łamanej, w jednostkach kroju. Poniżej pół jednostki punkty
	// zlewają się po zaokrągleniu.
	najmniejszyOdstepPunktowKrojuDesignu = 0.5

	// katNarozaKrojuDesignu jest kątem załamania łamanej, od którego wierzchołek
	// uznaje się za naroże — dwadzieścia pięć stopni leży wyraźnie między
	// załamaniem spłaszczonego łuku a prawdziwym narożem ikony.
	katNarozaKrojuDesignu = 25.0
)

// glifKrojuIkonowegoDesignu to jeden glif gotowy do zapisania: kontury w
// jednostkach kroju, postęp firetowy i prostokąt graniczny.
type glifKrojuIkonowegoDesignu struct {
	// Kontury niosą punkty konturu w jednostkach kroju, oś Y w górę.
	Kontury [][]punktKrojuDesignu
	Postep  int
	XMin    int
	YMin    int
	XMax    int
	YMax    int
}

// punktKrojuDesignu to jeden punkt konturu w jednostkach kroju. NaKonturze
// rozstrzyga, czy punkt leży na rysunku, czy jest punktem sterującym krzywej
// kwadratowej.
type punktKrojuDesignu struct {
	X          int
	Y          int
	NaKonturze bool
}

// zlozKrojIkonowyDesignu składa krój TrueType z ikon podanych jako pary (nazwa,
// treść SVG) i oddaje bajty pliku wraz z liczbą glifów. Ikona bez ani jednej
// ścieżki jest odmową całego pakietu, nie glifem pustym.
func zlozKrojIkonowyDesignu(nazwaKroju string, ikony []ikonaDoKrojuDesignu) ([]byte, int, error) {
	if len(ikony) == 0 {
		return nil, 0, fmt.Errorf("krój ikonowy bez ani jednej ikony nie ma czego nieść")
	}
	if len(ikony) > granicaGlifowKrojuIkonowegoDesignu {
		return nil, 0, fmt.Errorf(
			"krój ikonowy z %d ikonami przekracza obszar prywatny Unicode pierwszego poziomu "+
				"(%d miejsc od U+E000) — podziel pakiet na dwa",
			len(ikony), granicaGlifowKrojuIkonowegoDesignu)
	}

	// Glif zerowy jest `.notdef` i musi istnieć — czytnik bierze go dla znaku, którego krój nie ma.
	glify := []glifKrojuIkonowegoDesignu{{Postep: jednostekNaFiretDesignu}}
	punkty := make([]rune, 0, len(ikony))
	for numer, ikona := range ikony {
		glif, err := glifZIkonyDesignu(ikona)
		if err != nil {
			return nil, 0, fmt.Errorf("ikona %s: %w", ikona.Nazwa, err)
		}
		glify = append(glify, glif)
		punkty = append(punkty, rune(poczatekObszaruPrywatnegoDesignu+numer))
	}
	return zapiszKrojTrueTypeDesignu(nazwaKroju, glify, punkty), len(ikony), nil
}

// ikonaDoKrojuDesignu to jedna ikona wchodząca do kroju: nazwa, ścieżki SVG,
// siatka źródłowa i grubość obrysu, z których powstaje glif.
type ikonaDoKrojuDesignu struct {
	Nazwa   string
	Sciezki []string
	Siatka  int
	Grubosc float64
}

// glifZIkonyDesignu zamienia ścieżki ikony na kontury glifu przez obrys kreski
// skalowany do jednostek kroju.
func glifZIkonyDesignu(ikona ikonaDoKrojuDesignu) (glifKrojuIkonowegoDesignu, error) {
	if len(ikona.Sciezki) == 0 {
		return glifKrojuIkonowegoDesignu{}, fmt.Errorf(
			"treść SVG nie niesie ani jednej ścieżki (atrybutu d) — glif wyszedłby pusty")
	}
	grubosc := ikona.Grubosc
	if grubosc <= 0 {
		grubosc = gruboscObrysuIkonyDomyslna
	}
	skala := float64(jednostekNaFiretDesignu) / float64(siatkaIkonyKataloguDesignu)

	obszar := &canvas.Path{}
	for _, zapis := range ikona.Sciezki {
		sciezka, err := canvas.ParseSVGPath(zapis)
		if err != nil {
			return glifKrojuIkonowegoDesignu{}, fmt.Errorf(
				"ścieżki %q nie da się odczytać jako ścieżki SVG: %w", zapis, err)
		}
		// Kreska w obszar: obrys kreski staje się konturem glifu.
		obrys := sciezka.Stroke(grubosc, canvas.RoundCap, canvas.RoundJoin,
			tolerancjaObrysuKrojuDesignu)
		obszar = obszar.Append(obrys)
	}
	if obszar.Empty() {
		return glifKrojuIkonowegoDesignu{}, fmt.Errorf(
			"obrys ścieżek ikony wyszedł pusty — glif nie miałby konturu")
	}

	// Oś Y kroju rośnie w górę, siatki ikony — w dół; odbicie i skalowanie idą jedną macierzą.
	przeksztalcenie := canvas.Identity.Scale(skala, -skala).
		Translate(0, -float64(siatkaIkonyKataloguDesignu))
	// Spłaszczenie gwarantuje, że ścieżka niesie wyłącznie odcinki — dopasowanie niżej zakłada łamaną.
	lamany := obszar.Transform(przeksztalcenie).
		Flatten(tolerancjaSplaszczeniaKrojuDesignu)

	// Postęp jest zawsze firetem niezależnie od siatki ikony — inaczej glify różniłyby się szerokością.
	glif := glifKrojuIkonowegoDesignu{Postep: jednostekNaFiretDesignu}
	for _, lamana := range lamaneKonturowKrojuDesignu(lamany) {
		kontur := dopasujKrzyweKonturuDesignu(lamana, tolerancjaSplaszczeniaKrojuDesignu)
		if len(kontur) < 3 {
			continue
		}
		glif.Kontury = append(glif.Kontury, kontur)
	}
	if len(glif.Kontury) == 0 {
		return glifKrojuIkonowegoDesignu{}, fmt.Errorf(
			"po dopasowaniu krzywych nie został ani jeden kontur — kształt jest zbyt drobny " +
				"na siatce kroju")
	}

	glif.XMin, glif.YMin = math.MaxInt32, math.MaxInt32
	glif.XMax, glif.YMax = math.MinInt32, math.MinInt32
	for _, kontur := range glif.Kontury {
		for _, punkt := range kontur {
			if punkt.X < glif.XMin {
				glif.XMin = punkt.X
			}
			if punkt.Y < glif.YMin {
				glif.YMin = punkt.Y
			}
			if punkt.X > glif.XMax {
				glif.XMax = punkt.X
			}
			if punkt.Y > glif.YMax {
				glif.YMax = punkt.Y
			}
		}
	}
	return glif, nil
}

// lamaneKonturowKrojuDesignu rozbiera ścieżkę na łamane zamknięte, po jednej na
// kontur. Punkt leżący bliżej poprzedniego niż najmniejszy odstęp nie wchodzi —
// kontur z punktami zlewającymi się po zaokrągleniu bywa odrzucany przez
// rasteryzatory.
func lamaneKonturowKrojuDesignu(sciezka *canvas.Path) [][]canvas.Point {
	kontury := [][]canvas.Point{}
	biezacy := []canvas.Point{}
	dolozPunkt := func(punkt canvas.Point) {
		if len(biezacy) > 0 {
			poprzedni := biezacy[len(biezacy)-1]
			if math.Hypot(punkt.X-poprzedni.X, punkt.Y-poprzedni.Y) <
				najmniejszyOdstepPunktowKrojuDesignu {
				return
			}
		}
		biezacy = append(biezacy, punkt)
	}
	zamknijKontur := func() {
		if len(biezacy) > 2 {
			pierwszy, ostatni := biezacy[0], biezacy[len(biezacy)-1]
			// Punkt domykający równy pierwszemu jest w TrueType zbędny — kontur
			// zamyka się sam.
			if math.Hypot(ostatni.X-pierwszy.X, ostatni.Y-pierwszy.Y) <
				najmniejszyOdstepPunktowKrojuDesignu {
				biezacy = biezacy[:len(biezacy)-1]
			}
		}
		if len(biezacy) > 2 {
			kontury = append(kontury, biezacy)
		}
		biezacy = nil
	}

	skaner := sciezka.Scanner()
	for skaner.Scan() {
		switch skaner.Cmd() {
		case canvas.MoveToCmd:
			zamknijKontur()
			biezacy = []canvas.Point{skaner.End()}
		case canvas.CloseCmd:
			zamknijKontur()
		default:
			dolozPunkt(skaner.End())
		}
	}
	zamknijKontur()
	return kontury
}

// dopasujKrzyweKonturuDesignu zamienia łamaną w kontur glifu: ciągi wierzchołków
// leżące na jednej krzywej schodzą do krzywych kwadratowych, ciągi leżące na
// prostej — do odcinków. Bieg zachłanny kończy się na każdym rozpoznanym narożu.
func dopasujKrzyweKonturuDesignu(lamana []canvas.Point, tolerancja float64) []punktKrojuDesignu {
	ile := len(lamana)
	if ile < 3 {
		return nil
	}
	naroza := narozaLamanejKrojuDesignu(lamana)
	for numer, naroze := range naroza {
		if !naroze {
			continue
		}
		if numer > 0 {
			obrocona := make([]canvas.Point, 0, ile)
			obrocona = append(obrocona, lamana[numer:]...)
			obrocona = append(obrocona, lamana[:numer]...)
			lamana = obrocona
			naroza = narozaLamanejKrojuDesignu(lamana)
		}
		break
	}

	// Łamana rozszerzona o powrót do punktu pierwszego — ostatni bieg konturu kończy się właśnie tam.
	rozszerzona := make([]canvas.Point, 0, ile+1)
	rozszerzona = append(rozszerzona, lamana...)
	rozszerzona = append(rozszerzona, lamana[0])

	kontur := []punktKrojuDesignu{punktKrojuNaKonturzeDesignu(lamana[0])}
	poczatek := 0
	for poczatek < ile {
		koniec := poczatek + 1
		najlepszy := koniec
		jestKrzywa := false
		sterujacy := canvas.Point{}
		for koniec <= ile {
			if odejscieOdOdcinkaKrojuDesignu(rozszerzona[poczatek:koniec+1]) <= tolerancja {
				najlepszy, jestKrzywa = koniec, false
			} else if punkt, odejscie := dopasujKrzywaKwadratowaDesignu(
				rozszerzona[poczatek : koniec+1]); odejscie <= tolerancja {

				najlepszy, jestKrzywa, sterujacy = koniec, true, punkt
			} else {
				break
			}
			if koniec == ile || naroza[koniec] {
				break
			}
			koniec++
		}
		if jestKrzywa {
			kontur = append(kontur, punktKrojuSterujacyDesignu(sterujacy))
		}
		if najlepszy < ile {
			kontur = append(kontur, punktKrojuNaKonturzeDesignu(rozszerzona[najlepszy]))
		}
		poczatek = najlepszy
	}
	return bezPowtorzenKonturuKrojuDesignu(kontur)
}

// narozaLamanejKrojuDesignu wskazuje wierzchołki, w których łamana załamuje się
// mocniej niż o kąt naroża.
func narozaLamanejKrojuDesignu(lamana []canvas.Point) []bool {
	ile := len(lamana)
	naroza := make([]bool, ile)
	granica := math.Cos(katNarozaKrojuDesignu * math.Pi / 180)
	for numer := 0; numer < ile; numer++ {
		poprzedni := lamana[(numer-1+ile)%ile]
		biezacy := lamana[numer]
		nastepny := lamana[(numer+1)%ile]
		wejscieX, wejscieY := biezacy.X-poprzedni.X, biezacy.Y-poprzedni.Y
		wyjscieX, wyjscieY := nastepny.X-biezacy.X, nastepny.Y-biezacy.Y
		dlugoscWejscia := math.Hypot(wejscieX, wejscieY)
		dlugoscWyjscia := math.Hypot(wyjscieX, wyjscieY)
		if dlugoscWejscia <= 0 || dlugoscWyjscia <= 0 {
			continue
		}
		cosinus := (wejscieX*wyjscieX + wejscieY*wyjscieY) / (dlugoscWejscia * dlugoscWyjscia)
		naroza[numer] = cosinus < granica
	}
	return naroza
}

// odejscieOdOdcinkaKrojuDesignu oddaje największe odejście wierzchołków
// wewnętrznych od odcinka łączącego krańce ciągu.
func odejscieOdOdcinkaKrojuDesignu(ciag []canvas.Point) float64 {
	if len(ciag) < 3 {
		return 0
	}
	poczatek, koniec := ciag[0], ciag[len(ciag)-1]
	roznicaX, roznicaY := koniec.X-poczatek.X, koniec.Y-poczatek.Y
	dlugosc := math.Hypot(roznicaX, roznicaY)
	najwieksze := 0.0
	for _, punkt := range ciag[1 : len(ciag)-1] {
		odejscie := 0.0
		if dlugosc <= 0 {
			odejscie = math.Hypot(punkt.X-poczatek.X, punkt.Y-poczatek.Y)
		} else {
			// Odległość punktu od prostej: iloczyn wektorowy podzielony przez długość.
			odejscie = math.Abs((punkt.X-poczatek.X)*roznicaY-(punkt.Y-poczatek.Y)*roznicaX) /
				dlugosc
		}
		if odejscie > najwieksze {
			najwieksze = odejscie
		}
	}
	return najwieksze
}

// dopasujKrzywaKwadratowaDesignu liczy punkt sterujący krzywej kwadratowej
// przechodzącej przez ustalone krańce ciągu i najbliższej jego wierzchołkom,
// metodą najmniejszych kwadratów, wraz z największym odejściem od nich.
func dopasujKrzywaKwadratowaDesignu(ciag []canvas.Point) (canvas.Point, float64) {
	if len(ciag) < 3 {
		return ciag[0], 0
	}
	poczatek, koniec := ciag[0], ciag[len(ciag)-1]
	// Położenia wierzchołków na krzywej: udział w łącznej długości łamanej.
	dlugosci := make([]float64, len(ciag))
	for numer := 1; numer < len(ciag); numer++ {
		dlugosci[numer] = dlugosci[numer-1] +
			math.Hypot(ciag[numer].X-ciag[numer-1].X, ciag[numer].Y-ciag[numer-1].Y)
	}
	calosc := dlugosci[len(ciag)-1]
	if calosc <= 0 {
		return poczatek, math.MaxFloat64
	}

	licznikX, licznikY, mianownik := 0.0, 0.0, 0.0
	for numer := 1; numer < len(ciag)-1; numer++ {
		polozenie := dlugosci[numer] / calosc
		udzialSrodka := 2 * polozenie * (1 - polozenie)
		if udzialSrodka <= 0 {
			continue
		}
		resztaX := ciag[numer].X - (1-polozenie)*(1-polozenie)*poczatek.X -
			polozenie*polozenie*koniec.X
		resztaY := ciag[numer].Y - (1-polozenie)*(1-polozenie)*poczatek.Y -
			polozenie*polozenie*koniec.Y
		licznikX += udzialSrodka * resztaX
		licznikY += udzialSrodka * resztaY
		mianownik += udzialSrodka * udzialSrodka
	}
	if mianownik <= 0 {
		return poczatek, math.MaxFloat64
	}
	sterujacy := canvas.Point{X: licznikX / mianownik, Y: licznikY / mianownik}

	// Odejście liczy się po dopasowaniu — najmniejsze kwadraty minimalizują sumę, nie odejście pojedyncze.
	najwieksze := 0.0
	for numer := 1; numer < len(ciag)-1; numer++ {
		polozenie := dlugosci[numer] / calosc
		krzywaX := (1-polozenie)*(1-polozenie)*poczatek.X +
			2*polozenie*(1-polozenie)*sterujacy.X + polozenie*polozenie*koniec.X
		krzywaY := (1-polozenie)*(1-polozenie)*poczatek.Y +
			2*polozenie*(1-polozenie)*sterujacy.Y + polozenie*polozenie*koniec.Y
		odejscie := math.Hypot(ciag[numer].X-krzywaX, ciag[numer].Y-krzywaY)
		if odejscie > najwieksze {
			najwieksze = odejscie
		}
	}
	return sterujacy, najwieksze
}

// punktKrojuNaKonturzeDesignu i punktKrojuSterujacyDesignu zaokrąglają punkt do
// jednostek kroju. Jednostki kroju są całkowite — tabela `glyf` nie ma innej
// postaci współrzędnej.
func punktKrojuNaKonturzeDesignu(punkt canvas.Point) punktKrojuDesignu {
	return punktKrojuDesignu{
		X: int(math.Round(punkt.X)), Y: int(math.Round(punkt.Y)), NaKonturze: true,
	}
}

func punktKrojuSterujacyDesignu(punkt canvas.Point) punktKrojuDesignu {
	return punktKrojuDesignu{X: int(math.Round(punkt.X)), Y: int(math.Round(punkt.Y))}
}

// bezPowtorzenKonturuKrojuDesignu zdejmuje punkty konturu, które po zaokrągleniu
// do jednostek kroju zlały się z poprzednim. Dwa punkty konturu o tych samych
// współrzędnych bywają odrzucane przez rasteryzatory, a rysunku nie zmieniają.
func bezPowtorzenKonturuKrojuDesignu(kontur []punktKrojuDesignu) []punktKrojuDesignu {
	wynik := make([]punktKrojuDesignu, 0, len(kontur))
	for _, punkt := range kontur {
		if len(wynik) > 0 && wynik[len(wynik)-1] == punkt {
			continue
		}
		wynik = append(wynik, punkt)
	}
	// Punkt ostatni równy pierwszemu jest zbędny tak samo jak w łamanej: kontur
	// domyka się sam.
	if len(wynik) > 2 && wynik[len(wynik)-1] == wynik[0] {
		wynik = wynik[:len(wynik)-1]
	}
	return wynik
}

// zapiszKrojTrueTypeDesignu składa bajty pliku kroju z gotowych glifów: tabele
// `head`, `hhea`, `maxp`, `OS/2`, `hmtx`, `cmap`, `loca`, `glyf`, `name` i
// `post`, w katalogu ułożonym alfabetycznie po nazwie.
func zapiszKrojTrueTypeDesignu(nazwa string, glify []glifKrojuIkonowegoDesignu,
	punkty []rune) []byte {

	glyf, loca := zapiszGlifyDesignu(glify)
	tabele := map[string][]byte{
		"head": zapiszTabeleHeadDesignu(glify),
		"hhea": zapiszTabeleHheaDesignu(glify),
		"maxp": zapiszTabeleMaxpDesignu(glify),
		"OS/2": zapiszTabeleOs2Designu(glify, punkty),
		"hmtx": zapiszTabeleHmtxDesignu(glify),
		"cmap": zapiszTabeleCmapDesignu(punkty),
		"loca": loca,
		"glyf": glyf,
		"name": zapiszTabeleNameDesignu(nazwa),
		"post": zapiszTabelePostDesignu(),
	}

	nazwyTabel := make([]string, 0, len(tabele))
	for nazwaTabeli := range tabele {
		nazwyTabel = append(nazwyTabel, nazwaTabeli)
	}
	sort.Strings(nazwyTabel)

	naglowek := &bytes.Buffer{}
	_ = binary.Write(naglowek, binary.BigEndian, uint32(0x00010000)) // wersja: TrueType
	_ = binary.Write(naglowek, binary.BigEndian, uint16(len(tabele)))
	// Trójka pól wyszukiwania binarnego: potęga dwójki, jej podwojenie w bajtach wpisu i reszta.
	potega := uint16(1)
	log := uint16(0)
	for potega*2 <= uint16(len(tabele)) {
		potega *= 2
		log++
	}
	_ = binary.Write(naglowek, binary.BigEndian, potega*16)
	_ = binary.Write(naglowek, binary.BigEndian, log)
	_ = binary.Write(naglowek, binary.BigEndian, uint16(len(tabele))*16-potega*16)

	przesuniecie := uint32(12 + 16*len(tabele))
	katalog := &bytes.Buffer{}
	tresc := &bytes.Buffer{}
	for _, nazwaTabeli := range nazwyTabel {
		dane := tabele[nazwaTabeli]
		katalog.WriteString(dopelnijNazweTabeliDesignu(nazwaTabeli))
		_ = binary.Write(katalog, binary.BigEndian, sumaKontrolnaTabeliDesignu(dane))
		_ = binary.Write(katalog, binary.BigEndian, przesuniecie)
		_ = binary.Write(katalog, binary.BigEndian, uint32(len(dane)))
		tresc.Write(dane)
		// Wyrównanie do czterech bajtów: format wymaga, żeby każda tabela
		// zaczynała się na granicy słowa.
		for tresc.Len()%4 != 0 {
			tresc.WriteByte(0)
		}
		przesuniecie = uint32(12 + 16*len(tabele) + tresc.Len())
	}

	plik := &bytes.Buffer{}
	plik.Write(naglowek.Bytes())
	plik.Write(katalog.Bytes())
	plik.Write(tresc.Bytes())
	bajty := plik.Bytes()

	// checkSumAdjustment liczy się z sumy całego pliku, więc wpisuje się po złożeniu, w tabeli head.
	przesuniecieHead := uint32(0)
	for numer, nazwaTabeli := range nazwyTabel {
		if nazwaTabeli != "head" {
			continue
		}
		wpis := 12 + 16*numer + 8
		przesuniecieHead = binary.BigEndian.Uint32(bajty[wpis : wpis+4])
	}
	if przesuniecieHead > 0 {
		suma := uint32(0)
		for numer := 0; numer+3 < len(bajty); numer += 4 {
			suma += binary.BigEndian.Uint32(bajty[numer : numer+4])
		}
		binary.BigEndian.PutUint32(bajty[przesuniecieHead+8:przesuniecieHead+12],
			0xB1B0AFBA-suma)
	}
	return bajty
}

// dopelnijNazweTabeliDesignu wyrównuje nazwę tabeli do czterech znaków spacjami,
// zgodnie z zapisem katalogu tabel formatu TrueType.
func dopelnijNazweTabeliDesignu(nazwa string) string {
	for len(nazwa) < 4 {
		nazwa += " "
	}
	return nazwa[:4]
}

// sumaKontrolnaTabeliDesignu liczy sumę kontrolną tabeli wedle formatu: suma
// słów czterobajtowych z dopełnieniem zerami.
func sumaKontrolnaTabeliDesignu(dane []byte) uint32 {
	suma := uint32(0)
	for numer := 0; numer < len(dane); numer += 4 {
		slowo := uint32(0)
		for bajt := 0; bajt < 4; bajt++ {
			slowo <<= 8
			if numer+bajt < len(dane) {
				slowo |= uint32(dane[numer+bajt])
			}
		}
		suma += slowo
	}
	return suma
}

// zapiszGlifyDesignu składa tabele `glyf` i `loca` w postaci długiej —
// przesunięcia czterobajtowe bez warunku parzystej długości glifu, kosztem
// czterech bajtów na glif.
func zapiszGlifyDesignu(glify []glifKrojuIkonowegoDesignu) ([]byte, []byte) {
	glyf := &bytes.Buffer{}
	loca := &bytes.Buffer{}
	for _, glif := range glify {
		_ = binary.Write(loca, binary.BigEndian, uint32(glyf.Len()))
		if len(glif.Kontury) == 0 {
			// Glif pusty zajmuje zero bajtów — tak stanowi format i tak zapisuje
			// się `.notdef` oraz odstęp.
			continue
		}
		_ = binary.Write(glyf, binary.BigEndian, int16(len(glif.Kontury)))
		_ = binary.Write(glyf, binary.BigEndian, int16(glif.XMin))
		_ = binary.Write(glyf, binary.BigEndian, int16(glif.YMin))
		_ = binary.Write(glyf, binary.BigEndian, int16(glif.XMax))
		_ = binary.Write(glyf, binary.BigEndian, int16(glif.YMax))
		koniec := -1
		for _, kontur := range glif.Kontury {
			koniec += len(kontur)
			_ = binary.Write(glyf, binary.BigEndian, uint16(koniec))
		}
		_ = binary.Write(glyf, binary.BigEndian, uint16(0)) // bez instrukcji hintingu
		// Znacznik: bit zerowy mówi, czy punkt leży na konturze. Zapis bez powtórzeń i skrótów jednobajtowych.
		for _, kontur := range glif.Kontury {
			for _, punkt := range kontur {
				if punkt.NaKonturze {
					glyf.WriteByte(0x01)
					continue
				}
				glyf.WriteByte(0x00)
			}
		}
		poprzedni := 0
		for _, kontur := range glif.Kontury {
			for _, punkt := range kontur {
				_ = binary.Write(glyf, binary.BigEndian, int16(punkt.X-poprzedni))
				poprzedni = punkt.X
			}
		}
		poprzedni = 0
		for _, kontur := range glif.Kontury {
			for _, punkt := range kontur {
				_ = binary.Write(glyf, binary.BigEndian, int16(punkt.Y-poprzedni))
				poprzedni = punkt.Y
			}
		}
		for glyf.Len()%4 != 0 {
			glyf.WriteByte(0)
		}
	}
	_ = binary.Write(loca, binary.BigEndian, uint32(glyf.Len()))
	return glyf.Bytes(), loca.Bytes()
}

// zapiszTabeleHeadDesignu składa tabelę `head` formatu TrueType, niosącą wersję
// kroju, skalę jednostek i prostokąt graniczny wszystkich glifów.
func zapiszTabeleHeadDesignu(glify []glifKrojuIkonowegoDesignu) []byte {
	xMin, yMin, xMax, yMax := graniceKrojuDesignu(glify)
	bufor := &bytes.Buffer{}
	_ = binary.Write(bufor, binary.BigEndian, uint32(0x00010000)) // wersja tabeli
	_ = binary.Write(bufor, binary.BigEndian, uint32(0x00010000)) // wersja kroju
	_ = binary.Write(bufor, binary.BigEndian, uint32(0))          // checkSumAdjustment — wpisywane po złożeniu
	_ = binary.Write(bufor, binary.BigEndian, uint32(0x5F0F3CF5)) // magicNumber
	_ = binary.Write(bufor, binary.BigEndian, uint16(0x000B))     // flagi: bazowa w zerze, lsb, skala całkowita
	_ = binary.Write(bufor, binary.BigEndian, uint16(jednostekNaFiretDesignu))
	// Znacznik czasu liczy się od 1904 roku; wpisujemy zero, żeby pakiety nie różniły się bajtami.
	_ = binary.Write(bufor, binary.BigEndian, int64(0))
	_ = binary.Write(bufor, binary.BigEndian, int64(0))
	_ = binary.Write(bufor, binary.BigEndian, int16(xMin))
	_ = binary.Write(bufor, binary.BigEndian, int16(yMin))
	_ = binary.Write(bufor, binary.BigEndian, int16(xMax))
	_ = binary.Write(bufor, binary.BigEndian, int16(yMax))
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // macStyle
	_ = binary.Write(bufor, binary.BigEndian, uint16(8)) // lowestRecPPEM
	_ = binary.Write(bufor, binary.BigEndian, int16(2))  // fontDirectionHint
	_ = binary.Write(bufor, binary.BigEndian, int16(1))  // indexToLocFormat: postać długa
	_ = binary.Write(bufor, binary.BigEndian, int16(0))  // glyphDataFormat
	return bufor.Bytes()
}

// zapiszTabeleHheaDesignu składa tabelę `hhea` formatu TrueType, niosącą metryki
// poziome kroju wymagane do rozstawiania glifów w wierszu.
func zapiszTabeleHheaDesignu(glify []glifKrojuIkonowegoDesignu) []byte {
	xMin, _, xMax, _ := graniceKrojuDesignu(glify)
	najwiekszyPostep := 0
	for _, glif := range glify {
		if glif.Postep > najwiekszyPostep {
			najwiekszyPostep = glif.Postep
		}
	}
	bufor := &bytes.Buffer{}
	_ = binary.Write(bufor, binary.BigEndian, uint32(0x00010000))
	_ = binary.Write(bufor, binary.BigEndian, int16(jednostekNaFiretDesignu)) // ascender
	_ = binary.Write(bufor, binary.BigEndian, int16(0))                       // descender
	_ = binary.Write(bufor, binary.BigEndian, int16(0))                       // lineGap
	_ = binary.Write(bufor, binary.BigEndian, uint16(najwiekszyPostep))       // advanceWidthMax
	_ = binary.Write(bufor, binary.BigEndian, int16(xMin))                    // minLeftSideBearing
	_ = binary.Write(bufor, binary.BigEndian, int16(najwiekszyPostep-xMax))   // minRightSideBearing
	_ = binary.Write(bufor, binary.BigEndian, int16(xMax))                    // xMaxExtent
	_ = binary.Write(bufor, binary.BigEndian, int16(1))                       // caretSlopeRise
	_ = binary.Write(bufor, binary.BigEndian, int16(0))                       // caretSlopeRun
	_ = binary.Write(bufor, binary.BigEndian, int16(0))                       // caretOffset
	for numer := 0; numer < 4; numer++ {
		_ = binary.Write(bufor, binary.BigEndian, int16(0)) // pola zastrzeżone
	}
	_ = binary.Write(bufor, binary.BigEndian, int16(0)) // metricDataFormat
	_ = binary.Write(bufor, binary.BigEndian, uint16(len(glify)))
	return bufor.Bytes()
}

// zapiszTabeleMaxpDesignu składa tabelę `maxp` formatu TrueType, niosącą liczbę
// glifów i graniczne rozmiary konturu wymagane przez czytniki kroju.
func zapiszTabeleMaxpDesignu(glify []glifKrojuIkonowegoDesignu) []byte {
	najwiecejPunktow, najwiecejKonturow := 0, 0
	for _, glif := range glify {
		punktow := 0
		for _, kontur := range glif.Kontury {
			punktow += len(kontur)
		}
		if punktow > najwiecejPunktow {
			najwiecejPunktow = punktow
		}
		if len(glif.Kontury) > najwiecejKonturow {
			najwiecejKonturow = len(glif.Kontury)
		}
	}
	// Tabela `maxp` wersji 1.0 ma dokładnie 32 bajty — inna długość jest błędem, pola stoją tu z nazwami.
	bufor := &bytes.Buffer{}
	_ = binary.Write(bufor, binary.BigEndian, uint32(0x00010000))        // wersja 1.0
	_ = binary.Write(bufor, binary.BigEndian, uint16(len(glify)))        // numGlyphs
	_ = binary.Write(bufor, binary.BigEndian, uint16(najwiecejPunktow))  // maxPoints
	_ = binary.Write(bufor, binary.BigEndian, uint16(najwiecejKonturow)) // maxContours
	// Glifów złożonych krój ikonowy nie ma — każda ikona jest własnym konturem, stąd zera przy kompozycji.
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // maxCompositePoints
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // maxCompositeContours
	// Strefa: jedna, bo krój bez instrukcji nie używa strefy pomocniczej.
	_ = binary.Write(bufor, binary.BigEndian, uint16(1)) // maxZones
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // maxTwilightPoints
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // maxStorage
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // maxFunctionDefs
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // maxInstructionDefs
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // maxStackElements
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // maxSizeOfInstructions
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // maxComponentElements
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // maxComponentDepth
	return bufor.Bytes()
}

// zapiszTabeleOs2Designu składa tabelę `OS/2` w wersji 4 — bez niej krój nie
// wczyta się w systemie Windows ani w przeglądarce na nim.
func zapiszTabeleOs2Designu(glify []glifKrojuIkonowegoDesignu, punkty []rune) []byte {
	sredniPostep := 0
	if len(glify) > 0 {
		suma := 0
		for _, glif := range glify {
			suma += glif.Postep
		}
		sredniPostep = suma / len(glify)
	}
	pierwszy, ostatni := uint16(0xFFFF), uint16(0)
	for _, punkt := range punkty {
		if uint16(punkt) < pierwszy {
			pierwszy = uint16(punkt)
		}
		if uint16(punkt) > ostatni {
			ostatni = uint16(punkt)
		}
	}
	bufor := &bytes.Buffer{}
	_ = binary.Write(bufor, binary.BigEndian, uint16(4)) // wersja tabeli
	_ = binary.Write(bufor, binary.BigEndian, int16(sredniPostep))
	_ = binary.Write(bufor, binary.BigEndian, uint16(400)) // usWeightClass: zwykła
	_ = binary.Write(bufor, binary.BigEndian, uint16(5))   // usWidthClass: średnia
	_ = binary.Write(bufor, binary.BigEndian, uint16(0))   // fsType: bez ograniczeń osadzania
	// Pola indeksów i przekreślenia są wymagane, choć krój ich nie używa — zera bywają błędem.
	for _, wartosc := range []int16{650, 600, 0, 0, 650, 600, 0, 0, 50, 512} {
		_ = binary.Write(bufor, binary.BigEndian, wartosc)
	}
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // sFamilyClass
	bufor.Write([]byte{2, 0, 5, 3, 0, 0, 0, 0, 0, 0})    // panose: bezszeryfowa
	for numer := 0; numer < 4; numer++ {
		_ = binary.Write(bufor, binary.BigEndian, uint32(0)) // ulUnicodeRange
	}
	bufor.WriteString("DNCO")                               // achVendID
	_ = binary.Write(bufor, binary.BigEndian, uint16(0x40)) // fsSelection: zwykła
	_ = binary.Write(bufor, binary.BigEndian, pierwszy)
	_ = binary.Write(bufor, binary.BigEndian, ostatni)
	_ = binary.Write(bufor, binary.BigEndian, int16(jednostekNaFiretDesignu))   // sTypoAscender
	_ = binary.Write(bufor, binary.BigEndian, int16(0))                         // sTypoDescender
	_ = binary.Write(bufor, binary.BigEndian, int16(0))                         // sTypoLineGap
	_ = binary.Write(bufor, binary.BigEndian, uint16(jednostekNaFiretDesignu))  // usWinAscent
	_ = binary.Write(bufor, binary.BigEndian, uint16(0))                        // usWinDescent
	_ = binary.Write(bufor, binary.BigEndian, uint32(1))                        // ulCodePageRange1
	_ = binary.Write(bufor, binary.BigEndian, uint32(0))                        // ulCodePageRange2
	_ = binary.Write(bufor, binary.BigEndian, int16(jednostekNaFiretDesignu/2)) // sxHeight
	_ = binary.Write(bufor, binary.BigEndian, int16(jednostekNaFiretDesignu))   // sCapHeight
	_ = binary.Write(bufor, binary.BigEndian, uint16(0))                        // usDefaultChar
	_ = binary.Write(bufor, binary.BigEndian, uint16(0x20))                     // usBreakChar
	_ = binary.Write(bufor, binary.BigEndian, uint16(1))                        // usMaxContext
	return bufor.Bytes()
}

// zapiszTabeleHmtxDesignu składa tabelę `hmtx` formatu TrueType, niosącą postęp
// i lewy odstęp boczny każdego glifu kroju.
func zapiszTabeleHmtxDesignu(glify []glifKrojuIkonowegoDesignu) []byte {
	bufor := &bytes.Buffer{}
	for _, glif := range glify {
		_ = binary.Write(bufor, binary.BigEndian, uint16(glif.Postep))
		_ = binary.Write(bufor, binary.BigEndian, int16(glif.XMin))
	}
	return bufor.Bytes()
}

// zapiszTabeleCmapDesignu składa tabelę `cmap` z jednym podziałem formatu 4,
// który pokrywa płaszczyznę podstawową Unicode wraz z obszarem prywatnym
// pierwszego poziomu.
func zapiszTabeleCmapDesignu(punkty []rune) []byte {
	posortowane := make([]rune, len(punkty))
	copy(posortowane, punkty)
	sort.Slice(posortowane, func(i, j int) bool { return posortowane[i] < posortowane[j] })

	// Przedziały spójne: kolejny punkt kodowy o kolejnym numerze glifu wchodzi do
	// tego samego przedziału.
	type przedzialDesignu struct {
		poczatek uint16
		koniec   uint16
		glif     uint16
	}
	przedzialy := []przedzialDesignu{}
	for numer, punkt := range posortowane {
		glif := uint16(numer + 1) // glif zerowy to `.notdef`
		if len(przedzialy) > 0 {
			ostatni := &przedzialy[len(przedzialy)-1]
			if uint16(punkt) == ostatni.koniec+1 &&
				glif == ostatni.glif+(ostatni.koniec-ostatni.poczatek)+1 {
				ostatni.koniec = uint16(punkt)
				continue
			}
		}
		przedzialy = append(przedzialy, przedzialDesignu{
			poczatek: uint16(punkt), koniec: uint16(punkt), glif: glif,
		})
	}
	// Przedział domykający 0xFFFF jest w formacie wymagany.
	przedzialy = append(przedzialy, przedzialDesignu{poczatek: 0xFFFF, koniec: 0xFFFF, glif: 0})

	segmentow := len(przedzialy)
	podzial := &bytes.Buffer{}
	_ = binary.Write(podzial, binary.BigEndian, uint16(4)) // format
	_ = binary.Write(podzial, binary.BigEndian, uint16(16+8*segmentow))
	_ = binary.Write(podzial, binary.BigEndian, uint16(0)) // language
	_ = binary.Write(podzial, binary.BigEndian, uint16(2*segmentow))
	potega := uint16(1)
	log := uint16(0)
	for potega*2 <= uint16(segmentow) {
		potega *= 2
		log++
	}
	_ = binary.Write(podzial, binary.BigEndian, potega*2)
	_ = binary.Write(podzial, binary.BigEndian, log)
	_ = binary.Write(podzial, binary.BigEndian, uint16(2*segmentow)-potega*2)
	for _, przedzial := range przedzialy {
		_ = binary.Write(podzial, binary.BigEndian, przedzial.koniec)
	}
	_ = binary.Write(podzial, binary.BigEndian, uint16(0)) // reservedPad
	for _, przedzial := range przedzialy {
		_ = binary.Write(podzial, binary.BigEndian, przedzial.poczatek)
	}
	for _, przedzial := range przedzialy {
		// idDelta liczy się modulo 65536: numer glifu jest sumą punktu kodowego
		// i tej różnicy.
		_ = binary.Write(podzial, binary.BigEndian, int16(przedzial.glif-przedzial.poczatek))
	}
	for range przedzialy {
		_ = binary.Write(podzial, binary.BigEndian, uint16(0)) // idRangeOffset
	}

	bufor := &bytes.Buffer{}
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // wersja
	_ = binary.Write(bufor, binary.BigEndian, uint16(1)) // liczba podziałów
	_ = binary.Write(bufor, binary.BigEndian, uint16(3)) // platformID: Windows
	_ = binary.Write(bufor, binary.BigEndian, uint16(1)) // encodingID: BMP
	_ = binary.Write(bufor, binary.BigEndian, uint32(12))
	bufor.Write(podzial.Bytes())
	return bufor.Bytes()
}

// zapiszTabeleNameDesignu składa tabelę `name` z zapisami wymaganymi przez
// systemy i przeglądarki: nazwą rodziny, podrodziny, kroju pełnego i wersji.
func zapiszTabeleNameDesignu(nazwa string) []byte {
	nazwa = strings.TrimSpace(nazwa)
	if nazwa == "" {
		nazwa = "ikony"
	}
	// Nazwa postscriptowa nie może mieć odstępów ani znaków spoza zakresu ASCII drukowalnego.
	postscriptowa := strings.Map(func(znak rune) rune {
		if znak >= 'A' && znak <= 'Z' || znak >= 'a' && znak <= 'z' ||
			znak >= '0' && znak <= '9' || znak == '-' {
			return znak
		}
		return -1
	}, nazwa)
	if postscriptowa == "" {
		postscriptowa = "ikony"
	}
	zapisy := []struct {
		numer uint16
		tresc string
	}{
		{0, "Krój ikonowy złożony przez rdzeń Danaco Console"},
		{1, nazwa},
		{2, "Regular"},
		{3, postscriptowa + "-1-0"},
		{4, nazwa},
		{5, "Version 1.0"},
		{6, postscriptowa},
	}

	tresc := &bytes.Buffer{}
	katalog := &bytes.Buffer{}
	for _, zapis := range zapisy {
		// Kodowanie UTF-16BE, platforma Windows — jedyne, które czyta każdy system.
		bajty := &bytes.Buffer{}
		for _, znak := range zapis.tresc {
			_ = binary.Write(bajty, binary.BigEndian, uint16(znak))
		}
		_ = binary.Write(katalog, binary.BigEndian, uint16(3))      // platformID
		_ = binary.Write(katalog, binary.BigEndian, uint16(1))      // encodingID
		_ = binary.Write(katalog, binary.BigEndian, uint16(0x0409)) // languageID
		_ = binary.Write(katalog, binary.BigEndian, zapis.numer)
		_ = binary.Write(katalog, binary.BigEndian, uint16(bajty.Len()))
		_ = binary.Write(katalog, binary.BigEndian, uint16(tresc.Len()))
		tresc.Write(bajty.Bytes())
	}

	bufor := &bytes.Buffer{}
	_ = binary.Write(bufor, binary.BigEndian, uint16(0)) // format
	_ = binary.Write(bufor, binary.BigEndian, uint16(len(zapisy)))
	_ = binary.Write(bufor, binary.BigEndian, uint16(6+12*len(zapisy)))
	bufor.Write(katalog.Bytes())
	bufor.Write(tresc.Bytes())
	return bufor.Bytes()
}

// zapiszTabelePostDesignu składa tabelę `post` w wersji 3.0 — bez nazw glifów.
// Nazwy glifów w kroju ikonowym niosłyby te same napisy, które już są nazwami
// ikon w bazie, a wersja 3.0 jest w formacie przewidziana dokładnie na ten
// przypadek.
func zapiszTabelePostDesignu() []byte {
	bufor := &bytes.Buffer{}
	_ = binary.Write(bufor, binary.BigEndian, uint32(0x00030000)) // wersja 3.0
	_ = binary.Write(bufor, binary.BigEndian, uint32(0))          // italicAngle
	_ = binary.Write(bufor, binary.BigEndian, int16(0))           // underlinePosition
	_ = binary.Write(bufor, binary.BigEndian, int16(50))          // underlineThickness
	_ = binary.Write(bufor, binary.BigEndian, uint32(0))          // isFixedPitch
	for numer := 0; numer < 4; numer++ {
		_ = binary.Write(bufor, binary.BigEndian, uint32(0)) // pola pamięci VM
	}
	return bufor.Bytes()
}

// graniceKrojuDesignu oddaje prostokąt obejmujący wszystkie glify kroju, liczony
// z granic każdego glifu z osobna.
func graniceKrojuDesignu(glify []glifKrojuIkonowegoDesignu) (int, int, int, int) {
	xMin, yMin := math.MaxInt32, math.MaxInt32
	xMax, yMax := math.MinInt32, math.MinInt32
	maKontury := false
	for _, glif := range glify {
		if len(glif.Kontury) == 0 {
			continue
		}
		maKontury = true
		if glif.XMin < xMin {
			xMin = glif.XMin
		}
		if glif.YMin < yMin {
			yMin = glif.YMin
		}
		if glif.XMax > xMax {
			xMax = glif.XMax
		}
		if glif.YMax > yMax {
			yMax = glif.YMax
		}
	}
	if !maKontury {
		return 0, 0, jednostekNaFiretDesignu, jednostekNaFiretDesignu
	}
	return xMin, yMin, xMax, yMax
}
