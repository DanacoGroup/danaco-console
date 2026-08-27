// Plik obsługuje wykres (`design.chart.render`) i schemat
// (`design.diagram.render`) modułu design. Część drukarska leży
// w `adapter_modul_design_druk.go`.
package core

import (
	"bytes"
	"context"
	"fmt"
	"image/color"
	"math"
	"sort"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/tdewolff/canvas"
	renderRaster "github.com/tdewolff/canvas/renderers/rasterizer"
	renderSvg "github.com/tdewolff/canvas/renderers/svg"

	"danacoconsole/shared"
)

const (
	// domyslnaSzerokoscWykresuDesignu i domyslnaWysokoscWykresuDesignu to płótno
	// wykresu, gdy Operator wymiarów nie podał. Proporcja 16:10 jest tą, w której
	// wykres wchodzi do prezentacji bez przycinania.
	domyslnaSzerokoscWykresuDesignu = 800.0
	domyslnaWysokoscWykresuDesignu  = 500.0

	// granicaBokuWyrysuDanychDesignu chroni rachunek przed płótnem o boku
	// stutysięcznym, ograniczając bok wyrysu do dwudziestu tysięcy jednostek płótna.
	granicaBokuWyrysuDanychDesignu = 20000.0

	// granicaPunktowSeriiDesignu chroni wykres przed serią milionową: wykres
	// z milionem słupków nie jest wykresem.
	granicaPunktowSeriiDesignu = 5000

	// granicaWezlowSchematuDesignu chroni schemat przed grafem tysiącznym, ograniczając
	// liczbę węzłów do pięciuset.
	granicaWezlowSchematuDesignu = 500

	// rozdzielczoscWyrysuDanychDesignu jest rozdzielczością rasteryzacji wykresu
	// i schematu: dwa punkty na jednostkę płótna dają obraz ostry na ekranie
	// o podwyższonej gęstości.
	rozdzielczoscWyrysuDanychDesignu = 2.0

	// krojWyrysuDanychDesignu jest krojem podpisów. Krój WKOMPILOWANY, żeby wykres
	// wyszedł identycznie na każdej maszynie.
	krojWyrysuDanychDesignu = "Go Regular"

	// rozmiarPodpisuDanychDesignu i rozmiarTytuluDanychDesignu to rozmiary pisma
	// podpisów i tytułu w jednostkach płótna.
	rozmiarPodpisuDanychDesignu = 11.0
	rozmiarTytuluDanychDesignu  = 18.0
)

// barwySeriiDesignu to barwy domyślne serii wykresu i węzłów schematu. Wykaz
// jest ośmioelementowy, dobrany tak, żeby dał się rozróżnić także przy wadzie
// widzenia barw. Serię dziewiątą barwi ta sama lista od początku.
var barwySeriiDesignu = []string{
	"#1f6feb", "#f0883e", "#2da44e", "#cf222e",
	"#8250df", "#0e7490", "#bf8700", "#6e7781",
}

// WyrysujWykres składa wykres z serii danych i oddaje go jako zasób SVG albo
// PNG — obsługuje `design.chart.render`.
func (a *adapterDesignu) WyrysujWykres(ctx context.Context,
	z shared.DesignChartRenderRequest) (shared.DesignChartRenderResponse, error) {

	if err := sprawdzWyliczenieDesignu("design.chart.render", "kind", z.Kind,
		shared.WartosciDesignChartKind()); err != nil {
		return shared.DesignChartRenderResponse{}, err
	}
	if len(z.Series) == 0 {
		return shared.DesignChartRenderResponse{}, bladWskazaniaDesignu(
			"komenda design.chart.render bez ani jednej serii danych: wykres bez liczb nie ma " +
				"czego pokazać, a rdzeń nie wymyśla danych")
	}
	punktow := 0
	for numer, seria := range z.Series {
		if strings.TrimSpace(seria.Name) == "" {
			return shared.DesignChartRenderResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"seria numer %d bez nazwy: bez niej legenda nie ma czego wypisać", numer+1))
		}
		if len(seria.Values) == 0 {
			return shared.DesignChartRenderResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"seria %s nie ma ani jednej wartości", seria.Name))
		}
		punktow += len(seria.Values)
	}
	if punktow > granicaPunktowSeriiDesignu {
		return shared.DesignChartRenderResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"serie niosą %d wartości, a granica wyrysu to %d — wykres z tyloma punktami byłby "+
				"plamą", punktow, granicaPunktowSeriiDesignu))
	}
	format, err := formatWyrysuDanychDesignu("design.chart.render", z.Format)
	if err != nil {
		return shared.DesignChartRenderResponse{}, err
	}
	szerokosc, wysokosc, err := plotnoWyrysuDanychDesignu("design.chart.render", z.Width, z.Height)
	if err != nil {
		return shared.DesignChartRenderResponse{}, err
	}

	barwy, err := a.barwyWyrysuDanychDesignu(ctx, z.TokenSetId, len(z.Series))
	if err != nil {
		return shared.DesignChartRenderResponse{}, err
	}
	legenda := z.Legend == nil || *z.Legend

	plotno := canvas.New(szerokosc, wysokosc)
	kontekst := canvas.NewContext(plotno)
	if err := narysujWykresDesignu(kontekst, z, szerokosc, wysokosc, barwy, legenda); err != nil {
		return shared.DesignChartRenderResponse{}, bladWydaniaDesignu(err.Error())
	}

	bajty, typTresci, err := wydajPlotnoDanychDesignu(plotno, szerokosc, wysokosc, format)
	if err != nil {
		return shared.DesignChartRenderResponse{}, bladWydaniaDesignu(err.Error())
	}
	nazwa := "wykres " + string(z.Kind)
	if z.Title != nil && strings.TrimSpace(*z.Title) != "" {
		nazwa = strings.TrimSpace(*z.Title)
	}
	okno := oknoWytworu(z.WindowId, "")
	if strings.TrimSpace(okno) == "" {
		return shared.DesignChartRenderResponse{}, bladWskazaniaDesignu(
			"komenda design.chart.render bez wskazania okna: wykres jest zasobem magazynu " +
				"i musi mieć okno, w którym stanie")
	}
	zasob, err := a.zalozZasobZBajtowDesignu(ctx, okno, nazwa,
		rodzajWyrysuDanychDesignu(format), format, bajty)
	if err != nil {
		return shared.DesignChartRenderResponse{}, err
	}
	return shared.DesignChartRenderResponse{
		Asset: zasobWytworzonyKontraktu(zasob), MediaType: typTresci,
	}, nil
}

// narysujWykresDesignu rysuje wykres na płótnie: tło, tytuł, dane wybranym
// rodzajem wykresu oraz legendę, w marginesach obliczonych dla podpisów.
func narysujWykresDesignu(kontekst *canvas.Context, z shared.DesignChartRenderRequest,
	szerokosc, wysokosc float64, barwy []string, legenda bool) error {

	// Marginesy: lewy na podpisy wartości, dolny na kategorie, górny na tytuł, prawy na legendę.
	marginesLewy, marginesPrawy := 64.0, 24.0
	marginesGorny, marginesDolny := 24.0, 48.0
	if z.Title != nil && strings.TrimSpace(*z.Title) != "" {
		marginesGorny = 56.0
	}
	if legenda {
		marginesPrawy = 160.0
	}
	polePola := szerokosc - marginesLewy - marginesPrawy
	poleWysokosc := wysokosc - marginesGorny - marginesDolny
	if polePola <= 0 || poleWysokosc <= 0 {
		return fmt.Errorf("płótno %g×%g jest za małe na wykres z marginesami podpisów",
			szerokosc, wysokosc)
	}

	tlo := canvas.Rectangle(szerokosc, wysokosc)
	kontekst.RenderPath(tlo, stylWypelnieniaDanychDesignu("#ffffff"), canvas.Identity)

	if z.Title != nil && strings.TrimSpace(*z.Title) != "" {
		if err := napisWyrysuDanychDesignu(kontekst, strings.TrimSpace(*z.Title),
			marginesLewy, wysokosc-marginesGorny/2, rozmiarTytuluDanychDesignu, "#1f2328"); err != nil {
			return err
		}
	}

	switch z.Kind {
	case shared.DesignChartKindPie, shared.DesignChartKindDonut:
		narysujKoloweDesignu(kontekst, z, marginesLewy, marginesDolny, polePola, poleWysokosc, barwy)
	case shared.DesignChartKindScatter:
		if err := narysujPunktoweDesignu(kontekst, z, marginesLewy, marginesDolny,
			polePola, poleWysokosc, barwy); err != nil {
			return err
		}
	case shared.DesignChartKindLine, shared.DesignChartKindArea:
		if err := narysujLinioweDesignu(kontekst, z, marginesLewy, marginesDolny,
			polePola, poleWysokosc, barwy); err != nil {
			return err
		}
	default:
		if err := narysujSlupkoweDesignu(kontekst, z, marginesLewy, marginesDolny,
			polePola, poleWysokosc, barwy); err != nil {
			return err
		}
	}

	if legenda {
		if err := narysujLegendeDesignu(kontekst, z.Series, barwy,
			szerokosc-marginesPrawy+16, wysokosc-marginesGorny); err != nil {
			return err
		}
	}
	return nil
}

// zakresWartosciDesignu oddaje najmniejszą i największą wartość serii wraz
// z informacją, czy wykres jest skumulowany. Zakres bierze zero jako jeden
// z krańców, gdy wszystkie wartości są dodatnie, żeby nie przekłamać proporcji.
func zakresWartosciDesignu(z shared.DesignChartRenderRequest) (float64, float64) {
	skumulowany := z.Kind == shared.DesignChartKindStackedBar
	najmniejsza, najwieksza := math.MaxFloat64, -math.MaxFloat64
	if skumulowany {
		dlugosc := 0
		for _, seria := range z.Series {
			if len(seria.Values) > dlugosc {
				dlugosc = len(seria.Values)
			}
		}
		for numer := 0; numer < dlugosc; numer++ {
			suma := 0.0
			for _, seria := range z.Series {
				if numer < len(seria.Values) {
					suma += seria.Values[numer]
				}
			}
			najmniejsza = mniejszaDesignu(najmniejsza, suma)
			najwieksza = wiekszaDesignu(najwieksza, suma)
		}
	} else {
		for _, seria := range z.Series {
			for _, wartosc := range seria.Values {
				najmniejsza = mniejszaDesignu(najmniejsza, wartosc)
				najwieksza = wiekszaDesignu(najwieksza, wartosc)
			}
		}
	}
	if najmniejsza > 0 {
		najmniejsza = 0
	}
	if najwieksza < 0 {
		najwieksza = 0
	}
	if najwieksza == najmniejsza {
		// Wartości równe dałyby zakres zerowy i dzielenie przez zero — rozszerza się go o jedność.
		najwieksza = najmniejsza + 1
	}
	return najmniejsza, najwieksza
}

// narysujSlupkoweDesignu rysuje wykres słupkowy, kolumnowy albo skumulowany,
// dzieląc każdą grupę kategorii między słupki serii.
func narysujSlupkoweDesignu(kontekst *canvas.Context, z shared.DesignChartRenderRequest,
	lewa, dol, szerokosc, wysokosc float64, barwy []string) error {

	najmniejsza, najwieksza := zakresWartosciDesignu(z)
	poziomy := z.Kind == shared.DesignChartKindBar
	skumulowany := z.Kind == shared.DesignChartKindStackedBar

	kategorii := 0
	for _, seria := range z.Series {
		if len(seria.Values) > kategorii {
			kategorii = len(seria.Values)
		}
	}
	if kategorii == 0 {
		return fmt.Errorf("serie nie mają ani jednej kategorii")
	}
	narysujOsieDesignu(kontekst, lewa, dol, szerokosc, wysokosc)

	// Grupa jest miejscem jednej kategorii; słupki dzielą ją między sobą, a skumulowany — stoi na sobie.
	dlugoscOsi := szerokosc
	if poziomy {
		dlugoscOsi = wysokosc
	}
	grupa := dlugoscOsi / float64(kategorii)
	odstepGrupy := grupa * 0.2
	slupkow := len(z.Series)
	if skumulowany {
		slupkow = 1
	}
	bokSlupka := (grupa - odstepGrupy) / float64(slupkow)

	for numerKategorii := 0; numerKategorii < kategorii; numerKategorii++ {
		podstawaDodatnia, podstawaUjemna := 0.0, 0.0
		for numerSerii, seria := range z.Series {
			if numerKategorii >= len(seria.Values) {
				continue
			}
			wartosc := seria.Values[numerKategorii]
			poczatek := podstawaDodatnia
			if wartosc < 0 {
				poczatek = podstawaUjemna
			}
			if !skumulowany {
				poczatek = 0
			}
			odKrancaA := udzialWZakresieDesignu(poczatek, najmniejsza, najwieksza)
			odKrancaB := udzialWZakresieDesignu(poczatek+wartosc, najmniejsza, najwieksza)
			if skumulowany {
				if wartosc < 0 {
					podstawaUjemna += wartosc
				} else {
					podstawaDodatnia += wartosc
				}
			}

			przesuniecie := odstepGrupy/2 + float64(numerSerii)*bokSlupka
			if skumulowany {
				przesuniecie = odstepGrupy / 2
			}
			barwa := barwy[numerSerii%len(barwy)]
			if poziomy {
				y := dol + float64(numerKategorii)*grupa + przesuniecie
				x1 := lewa + odKrancaA*szerokosc
				x2 := lewa + odKrancaB*szerokosc
				narysujProstokatDanychDesignu(kontekst, mniejszaDesignu(x1, x2), y,
					math.Abs(x2-x1), bokSlupka*0.92, barwa)
				continue
			}
			x := lewa + float64(numerKategorii)*grupa + przesuniecie
			y1 := dol + odKrancaA*wysokosc
			y2 := dol + odKrancaB*wysokosc
			narysujProstokatDanychDesignu(kontekst, x, mniejszaDesignu(y1, y2),
				bokSlupka*0.92, math.Abs(y2-y1), barwa)
		}
	}
	return narysujPodpisyKategoriiDesignu(kontekst, z.Categories, kategorii,
		lewa, dol, szerokosc, wysokosc, poziomy)
}

// narysujLinioweDesignu rysuje wykres liniowy albo warstwowy, łącząc punkty
// serii odcinkami na polu danych.
func narysujLinioweDesignu(kontekst *canvas.Context, z shared.DesignChartRenderRequest,
	lewa, dol, szerokosc, wysokosc float64, barwy []string) error {

	najmniejsza, najwieksza := zakresWartosciDesignu(z)
	narysujOsieDesignu(kontekst, lewa, dol, szerokosc, wysokosc)
	kategorii := 0
	for _, seria := range z.Series {
		if len(seria.Values) > kategorii {
			kategorii = len(seria.Values)
		}
	}
	for numerSerii, seria := range z.Series {
		barwa := barwy[numerSerii%len(barwy)]
		if seria.Color != nil && strings.TrimSpace(*seria.Color) != "" {
			barwa = *seria.Color
		}
		linia := &canvas.Path{}
		for numer, wartosc := range seria.Values {
			x := lewa + udzialPozycjiDesignu(numer, len(seria.Values))*szerokosc
			y := dol + udzialWZakresieDesignu(wartosc, najmniejsza, najwieksza)*wysokosc
			if numer == 0 {
				linia.MoveTo(x, y)
				continue
			}
			linia.LineTo(x, y)
		}
		if z.Kind == shared.DesignChartKindArea && len(seria.Values) > 1 {
			// Wykres warstwowy domyka się do linii zera, bo to obszar POD krzywą
			// jest jego treścią.
			zero := dol + udzialWZakresieDesignu(0, najmniejsza, najwieksza)*wysokosc
			linia.LineTo(lewa+szerokosc, zero)
			linia.LineTo(lewa, zero)
			linia.Close()
			polprzezroczysta := 0.35
			kontekst.RenderPath(linia, canvas.Style{
				Fill:     canvas.Paint{Color: barwaRgbaZapisuDesignu(barwa, &polprzezroczysta)},
				FillRule: canvas.NonZero,
			}, canvas.Identity)
			continue
		}
		kontekst.RenderPath(linia, zObrysemDanychDesignu(canvas.Style{
			Stroke:      canvas.Paint{Color: barwaRgbaZapisuDesignu(barwa, nil)},
			StrokeWidth: 2,
		}), canvas.Identity)
	}
	return narysujPodpisyKategoriiDesignu(kontekst, z.Categories, kategorii,
		lewa, dol, szerokosc, wysokosc, false)
}

// narysujPunktoweDesignu rysuje wykres punktowy, stawiając jeden znacznik dla
// każdej wartości serii na polu danych.
func narysujPunktoweDesignu(kontekst *canvas.Context, z shared.DesignChartRenderRequest,
	lewa, dol, szerokosc, wysokosc float64, barwy []string) error {

	najmniejsza, najwieksza := zakresWartosciDesignu(z)
	narysujOsieDesignu(kontekst, lewa, dol, szerokosc, wysokosc)
	for numerSerii, seria := range z.Series {
		barwa := barwy[numerSerii%len(barwy)]
		if seria.Color != nil && strings.TrimSpace(*seria.Color) != "" {
			barwa = *seria.Color
		}
		for numer, wartosc := range seria.Values {
			x := lewa + udzialPozycjiDesignu(numer, len(seria.Values))*szerokosc
			y := dol + udzialWZakresieDesignu(wartosc, najmniejsza, najwieksza)*wysokosc
			punkt := canvas.Circle(4).Translate(x, y)
			kontekst.RenderPath(punkt, stylWypelnieniaDanychDesignu(barwa), canvas.Identity)
		}
	}
	return nil
}

// narysujKoloweDesignu rysuje wykres kołowy albo pierścieniowy.
//
// Wykres kołowy pokazuje UDZIAŁY, więc bierze pierwszą serię i jej wartości
// bezwzględne: udział ujemny nie istnieje, a wartość ujemna w kole byłaby
// wycinkiem o negatywnym kącie.
func narysujKoloweDesignu(kontekst *canvas.Context, z shared.DesignChartRenderRequest,
	lewa, dol, szerokosc, wysokosc float64, barwy []string) {

	wartosci := z.Series[0].Values
	suma := 0.0
	for _, wartosc := range wartosci {
		suma += math.Abs(wartosc)
	}
	if suma == 0 {
		return
	}
	promien := mniejszaDesignu(szerokosc, wysokosc) / 2 * 0.9
	srodekX, srodekY := lewa+szerokosc/2, dol+wysokosc/2
	promienWewnetrzny := 0.0
	if z.Kind == shared.DesignChartKindDonut {
		promienWewnetrzny = promien * 0.55
	}

	kat := math.Pi / 2 // pierwszy wycinek zaczyna się u góry
	for numer, wartosc := range wartosci {
		udzial := math.Abs(wartosc) / suma
		koniec := kat - udzial*2*math.Pi
		wycinek := wycinekKolaDesignu(srodekX, srodekY, promien, promienWewnetrzny, kat, koniec)
		barwa := barwy[numer%len(barwy)]
		kontekst.RenderPath(wycinek, stylWypelnieniaDanychDesignu(barwa), canvas.Identity)
		kat = koniec
	}
}

// wycinekKolaDesignu składa ścieżkę wycinka koła albo pierścienia. Łuk idzie
// odcinkami o kroku poniżej stopnia, nierozróżnialnymi od łuku na oko.
func wycinekKolaDesignu(srodekX, srodekY, promien, promienWewnetrzny,
	od, do float64) *canvas.Path {

	const krok = math.Pi / 180
	sciezka := &canvas.Path{}
	kroki := int(math.Abs(do-od)/krok) + 2
	for numer := 0; numer <= kroki; numer++ {
		kat := od + (do-od)*float64(numer)/float64(kroki)
		x := srodekX + promien*math.Cos(kat)
		y := srodekY + promien*math.Sin(kat)
		if numer == 0 {
			sciezka.MoveTo(x, y)
			continue
		}
		sciezka.LineTo(x, y)
	}
	if promienWewnetrzny <= 0 {
		sciezka.LineTo(srodekX, srodekY)
		sciezka.Close()
		return sciezka
	}
	for numer := kroki; numer >= 0; numer-- {
		kat := od + (do-od)*float64(numer)/float64(kroki)
		sciezka.LineTo(srodekX+promienWewnetrzny*math.Cos(kat),
			srodekY+promienWewnetrzny*math.Sin(kat))
	}
	sciezka.Close()
	return sciezka
}

// narysujOsieDesignu rysuje osie pola danych: pionową po lewej i poziomą u
// dołu, jako jedną łamaną linię.
func narysujOsieDesignu(kontekst *canvas.Context, lewa, dol, szerokosc, wysokosc float64) {
	osie := &canvas.Path{}
	osie.MoveTo(lewa, dol+wysokosc)
	osie.LineTo(lewa, dol)
	osie.LineTo(lewa+szerokosc, dol)
	kontekst.RenderPath(osie, zObrysemDanychDesignu(canvas.Style{
		Stroke: canvas.Paint{Color: barwaRgbaZapisuDesignu("#8c959f", nil)}, StrokeWidth: 1,
	}), canvas.Identity)
}

// narysujProstokatDanychDesignu rysuje jeden słupek wykresu jako wypełniony
// prostokąt o wskazanych wymiarach.
func narysujProstokatDanychDesignu(kontekst *canvas.Context, x, y, szerokosc, wysokosc float64,
	barwa string) {

	if szerokosc <= 0 || wysokosc <= 0 {
		// Słupek zerowy nie jest rysowany, ale jego miejsce w grupie zostaje.
		return
	}
	prostokat := canvas.Rectangle(szerokosc, wysokosc).Translate(x, y)
	kontekst.RenderPath(prostokat, stylWypelnieniaDanychDesignu(barwa), canvas.Identity)
}

// narysujPodpisyKategoriiDesignu podpisuje kategorie osi. Podpisy wchodzą co
// n-tą kategorię, gdy nie mieszczą się obok siebie — krok liczy się z szerokości
// pola, nie jest zmyślany.
func narysujPodpisyKategoriiDesignu(kontekst *canvas.Context, kategorie []string, ile int,
	lewa, dol, szerokosc, wysokosc float64, poziomy bool) error {

	if len(kategorie) == 0 {
		return nil
	}
	// Szacowana szerokość podpisu: rozmiar pisma razy liczba znaków, z zapasem.
	najdluzszy := 0
	for _, podpis := range kategorie {
		if len([]rune(podpis)) > najdluzszy {
			najdluzszy = len([]rune(podpis))
		}
	}
	miejsceNaPodpis := szerokosc / float64(ile)
	if poziomy {
		miejsceNaPodpis = wysokosc / float64(ile)
	}
	krok := 1
	szerokoscPodpisu := float64(najdluzszy) * rozmiarPodpisuDanychDesignu * 0.6
	for miejsceNaPodpis*float64(krok) < szerokoscPodpisu && krok < ile {
		krok++
	}
	for numer, podpis := range kategorie {
		if numer >= ile || numer%krok != 0 || strings.TrimSpace(podpis) == "" {
			continue
		}
		if poziomy {
			y := dol + (float64(numer)+0.5)*wysokosc/float64(ile)
			if err := napisWyrysuDanychDesignu(kontekst, podpis, 8, y,
				rozmiarPodpisuDanychDesignu, "#57606a"); err != nil {
				return err
			}
			continue
		}
		x := lewa + (float64(numer)+0.5)*szerokosc/float64(ile) -
			float64(len([]rune(podpis)))*rozmiarPodpisuDanychDesignu*0.28
		if err := napisWyrysuDanychDesignu(kontekst, podpis, x, dol-18,
			rozmiarPodpisuDanychDesignu, "#57606a"); err != nil {
			return err
		}
	}
	return nil
}

// narysujLegendeDesignu rysuje legendę serii: próbkę barwy i nazwę dla każdej
// serii, w kolumnie od góry.
func narysujLegendeDesignu(kontekst *canvas.Context, serie []shared.DesignChartSeries,
	barwy []string, lewa, gora float64) error {

	y := gora
	for numer, seria := range serie {
		barwa := barwy[numer%len(barwy)]
		if seria.Color != nil && strings.TrimSpace(*seria.Color) != "" {
			barwa = *seria.Color
		}
		probka := canvas.Rectangle(10, 10).Translate(lewa, y-10)
		kontekst.RenderPath(probka, stylWypelnieniaDanychDesignu(barwa), canvas.Identity)
		if err := napisWyrysuDanychDesignu(kontekst, seria.Name, lewa+16, y-9,
			rozmiarPodpisuDanychDesignu, "#1f2328"); err != nil {
			return err
		}
		y -= 20
		if y < 20 {
			// Legenda dłuższa niż płótno urywa się i mówi to wprost, zamiast wyglądać jak wykaz kompletny.
			return napisWyrysuDanychDesignu(kontekst,
				fmt.Sprintf("… i %d dalszych serii", len(serie)-numer-1),
				lewa, y, rozmiarPodpisuDanychDesignu, "#57606a")
		}
	}
	return nil
}

// napisWyrysuDanychDesignu rysuje napis jako kontury glifów kroju wkompilowanego,
// niezależnie od krojów zainstalowanych na maszynie odbiorcy.
func napisWyrysuDanychDesignu(kontekst *canvas.Context, tresc string, x, y, rozmiar float64,
	barwa string) error {

	if strings.TrimSpace(tresc) == "" {
		return nil
	}
	krojWczytany, _, _, err := krojDesignu(krojWyrysuDanychDesignu)
	if err != nil {
		return fmt.Errorf("podpisów wykresu nie da się złożyć: %w", err)
	}
	// Kontury tekstu mają oś Y w dół, a płótno biblioteki w górę — stąd odbicie względem linii pisma.
	kontury, err := sciezkaTekstuDesignu(krojWczytany, tresc, rozmiar, 0, 0)
	if err != nil {
		// Znak spoza kroju nie kończy wyrysu — napis wypada, wykres zostaje.
		return nil
	}
	kontekst.RenderPath(kontury.Transform(canvas.Identity.Translate(x, y).Scale(1, -1)),
		stylWypelnieniaDanychDesignu(barwa), canvas.Identity)
	return nil
}

// stylWypelnieniaDanychDesignu składa styl wypełnienia jednolitego dla wskazanej
// barwy płótna, bez obrysu.
func stylWypelnieniaDanychDesignu(barwa string) canvas.Style {
	return canvas.Style{
		Fill: canvas.Paint{Color: barwaRgbaZapisuDesignu(barwa, nil)}, FillRule: canvas.NonZero,
	}
}

// zObrysemDanychDesignu dopełnia styl zakończeniem i złączeniem kreski: wydawca
// SVG biblioteki nie ma dla nich wartości domyślnej, a styl z obrysem bez nich
// przerywa mu wykonanie.
func zObrysemDanychDesignu(styl canvas.Style) canvas.Style {
	if styl.StrokeCapper == nil {
		styl.StrokeCapper = canvas.RoundCap
	}
	if styl.StrokeJoiner == nil {
		styl.StrokeJoiner = canvas.RoundJoin
	}
	return styl
}

// barwaRgbaZapisuDesignu rozpoznaje zapis barwy i oddaje składowe. Zapis
// nieczytelny daje czerń — tu, w wyrysie wewnętrznym, gdzie barwy pochodzą
// z wykazu rdzenia albo z żetonów już sprawdzonych, a nie wprost z żądania.
func barwaRgbaZapisuDesignu(zapis string, krycie *float64) color.RGBA {
	barwa, err := rozpoznajBarweDesignu(zapis)
	if err != nil {
		barwa = colorful.Color{}
	}
	return barwaRgbaDesignu(barwa, krycie)
}

// udzialWZakresieDesignu przekłada wartość na ułamek wysokości pola danych,
// względem krańców zakresu serii.
func udzialWZakresieDesignu(wartosc, najmniejsza, najwieksza float64) float64 {
	return (wartosc - najmniejsza) / (najwieksza - najmniejsza)
}

// udzialPozycjiDesignu przekłada numer punktu na ułamek szerokości pola,
// rozkładając punkty równomiernie.
func udzialPozycjiDesignu(numer, ile int) float64 {
	if ile <= 1 {
		return 0
	}
	return float64(numer) / float64(ile-1)
}

// barwyWyrysuDanychDesignu oddaje barwy wyrysu: z zestawu żetonów, gdy wskazany,
// albo z wykazu rdzenia. Zestaw bez ani jednego żetonu barwnego jest odmową —
// zestaw narzuca barwy, po to jest wskazywany.
func (a *adapterDesignu) barwyWyrysuDanychDesignu(ctx context.Context, zestaw *string,
	ile int) ([]string, error) {

	if zestaw == nil || strings.TrimSpace(*zestaw) == "" {
		return barwySeriiDesignu, nil
	}
	wiersz, err := a.repozytorium.ZestawZetonowDesignuPoKodzie(ctx, strings.TrimSpace(*zestaw))
	if err != nil {
		if czyBrakZasobuDesignu(err) {
			return nil, bladNieznanegoBytuDesignu(
				"zestawu żetonów " + *zestaw + " nie ma w tym rdzeniu")
		}
		return nil, bladDesignu(err)
	}
	zetony, err := a.repozytorium.ZetonyZestawuDesignu(ctx, wiersz.ID)
	if err != nil {
		return nil, bladDesignu(err)
	}
	barwy := []string{}
	for _, zeton := range zetony {
		if zeton.Rodzaj != shared.DesignTokenKindColor {
			continue
		}
		if _, err := rozpoznajBarweDesignu(zeton.Wartosc); err != nil {
			continue
		}
		barwy = append(barwy, zeton.Wartosc)
	}
	if len(barwy) == 0 {
		return nil, bladWskazaniaDesignu(fmt.Sprintf(
			"zestaw żetonów %s nie ma ani jednego żetonu barwnego, a został wskazany jako "+
				"narzucający barwy — rdzeń nie podstawi za niego własnych", wiersz.Kod))
	}
	if len(barwy) < ile {
		// Barw mniej niż serii: wykaz powtarza się od początku, zamiast dobierać barwy spoza zestawu.
		sort.SliceStable(barwy, func(i, j int) bool { return barwy[i] < barwy[j] })
	}
	return barwy, nil
}

// formatWyrysuDanychDesignu rozstrzyga format wyrysu. Brak wskazania bierze SVG:
// wykres wektorowy skaluje się bez utraty ostrości, a to jest jego naturalna
// postać.
func formatWyrysuDanychDesignu(komenda string, wskazany *string) (string, error) {
	if wskazany == nil || strings.TrimSpace(*wskazany) == "" {
		return "svg", nil
	}
	format := normalizujFormatWydaniaDesignu(*wskazany)
	if format != "svg" && format != "png" {
		return "", bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s z formatem %q: wyrys wychodzi jako svg albo png", komenda, *wskazany))
	}
	return format, nil
}

// rodzajWyrysuDanychDesignu nazywa rodzaj zasobu wyrysu wynikowego: wektorowy
// dla SVG, obrazowy dla PNG.
func rodzajWyrysuDanychDesignu(format string) shared.DesignAssetKind {
	if format == "svg" {
		return shared.DesignAssetKindVector
	}
	return shared.DesignAssetKindImage
}

// plotnoWyrysuDanychDesignu rozstrzyga wymiary płótna: wskazane w żądaniu,
// a bez wskazania — domyślne, w granicy dopuszczalnego boku.
func plotnoWyrysuDanychDesignu(komenda string, szerokosc, wysokosc *float64) (float64, float64, error) {
	wynikSzerokosc := domyslnaSzerokoscWykresuDesignu
	wynikWysokosc := domyslnaWysokoscWykresuDesignu
	if szerokosc != nil {
		wynikSzerokosc = *szerokosc
	}
	if wysokosc != nil {
		wynikWysokosc = *wysokosc
	}
	if wynikSzerokosc <= 0 || wynikWysokosc <= 0 {
		return 0, 0, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s z płótnem %g×%g: płótno o niedodatnim boku nie istnieje",
			komenda, wynikSzerokosc, wynikWysokosc))
	}
	if wynikSzerokosc > granicaBokuWyrysuDanychDesignu ||
		wynikWysokosc > granicaBokuWyrysuDanychDesignu {

		return 0, 0, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s z płótnem %g×%g: granica boku wyrysu to %g",
			komenda, wynikSzerokosc, wynikWysokosc, granicaBokuWyrysuDanychDesignu))
	}
	return wynikSzerokosc, wynikWysokosc, nil
}

// wydajPlotnoDanychDesignu wydaje płótno jako SVG albo PNG — jedno płótno,
// dwa wydania tej samej treści.
func wydajPlotnoDanychDesignu(plotno *canvas.Canvas, szerokosc, wysokosc float64,
	format string) ([]byte, string, error) {

	if format == "svg" {
		var bufor bytes.Buffer
		wydawca := renderSvg.New(&bufor, szerokosc, wysokosc, nil)
		plotno.RenderTo(wydawca)
		if err := wydawca.Close(); err != nil {
			return nil, "", fmt.Errorf("nie można zapisać wyrysu jako svg: %w", err)
		}
		return bufor.Bytes(), typTresciWydaniaDesignu("svg"), nil
	}
	obraz := renderRaster.Draw(plotno, canvas.Resolution(rozdzielczoscWyrysuDanychDesignu),
		canvas.DefaultColorSpace)
	bajty, typTresci, err := zakodujObrazDesignu(obraz, "png", nil)
	if err != nil {
		return nil, "", err
	}
	return bajty, typTresci, nil
}

// WyrysujSchemat składa schemat z węzłów i połączeń, układa go warstwami
// i oddaje jako zasób — obsługuje `design.diagram.render`.
func (a *adapterDesignu) WyrysujSchemat(ctx context.Context,
	z shared.DesignDiagramRenderRequest) (shared.DesignDiagramRenderResponse, error) {

	if err := sprawdzWyliczenieDesignu("design.diagram.render", "kind", z.Kind,
		shared.WartosciDesignDiagramKind()); err != nil {
		return shared.DesignDiagramRenderResponse{}, err
	}
	if z.Direction != nil {
		if err := sprawdzWyliczenieDesignu("design.diagram.render", "direction", *z.Direction,
			shared.WartosciDesignLayoutDirection()); err != nil {
			return shared.DesignDiagramRenderResponse{}, err
		}
	}
	if len(z.Nodes) == 0 {
		return shared.DesignDiagramRenderResponse{}, bladWskazaniaDesignu(
			"komenda design.diagram.render bez ani jednego węzła: schemat bez węzłów nie ma " +
				"czego pokazać")
	}
	if len(z.Nodes) > granicaWezlowSchematuDesignu {
		return shared.DesignDiagramRenderResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"schemat z %d węzłami przekracza granicę %d — układ wyszedłby nieczytelny",
			len(z.Nodes), granicaWezlowSchematuDesignu))
	}
	znane := map[string]bool{}
	for numer, wezel := range z.Nodes {
		if strings.TrimSpace(wezel.Id) == "" {
			return shared.DesignDiagramRenderResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"węzeł numer %d bez identyfikatora: połączenia nie miałyby czym go wskazać", numer+1))
		}
		if znane[wezel.Id] {
			return shared.DesignDiagramRenderResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"węzeł %s występuje dwa razy — dwa węzły o tym samym identyfikatorze nie dają "+
					"się rozróżnić w połączeniach", wezel.Id))
		}
		znane[wezel.Id] = true
		if wezel.Shape != nil {
			if err := sprawdzWyliczenieDesignu("design.diagram.render",
				fmt.Sprintf("nodes[%d].shape", numer), *wezel.Shape,
				shared.WartosciDesignShapeKind()); err != nil {
				return shared.DesignDiagramRenderResponse{}, err
			}
		}
	}
	format, err := formatWyrysuDanychDesignu("design.diagram.render", z.Format)
	if err != nil {
		return shared.DesignDiagramRenderResponse{}, err
	}
	okno := oknoWytworu(z.WindowId, "")
	if strings.TrimSpace(okno) == "" {
		return shared.DesignDiagramRenderResponse{}, bladWskazaniaDesignu(
			"komenda design.diagram.render bez wskazania okna: schemat jest zasobem magazynu " +
				"i musi mieć okno, w którym stanie")
	}
	barwy, err := a.barwyWyrysuDanychDesignu(ctx, z.TokenSetId, len(z.Nodes))
	if err != nil {
		return shared.DesignDiagramRenderResponse{}, err
	}

	polozenia, nieumieszczone, szerokosc, wysokosc := ulozSchematDesignu(z)
	plotno := canvas.New(szerokosc, wysokosc)
	kontekst := canvas.NewContext(plotno)
	kontekst.RenderPath(canvas.Rectangle(szerokosc, wysokosc),
		stylWypelnieniaDanychDesignu("#ffffff"), canvas.Identity)
	if err := narysujSchematDesignu(kontekst, z, polozenia, barwy, wysokosc); err != nil {
		return shared.DesignDiagramRenderResponse{}, bladWydaniaDesignu(err.Error())
	}

	bajty, typTresci, err := wydajPlotnoDanychDesignu(plotno, szerokosc, wysokosc, format)
	if err != nil {
		return shared.DesignDiagramRenderResponse{}, bladWydaniaDesignu(err.Error())
	}
	nazwa := "schemat " + string(z.Kind)
	if z.Title != nil && strings.TrimSpace(*z.Title) != "" {
		nazwa = strings.TrimSpace(*z.Title)
	}
	zasob, err := a.zalozZasobZBajtowDesignu(ctx, okno, nazwa,
		rodzajWyrysuDanychDesignu(format), format, bajty)
	if err != nil {
		return shared.DesignDiagramRenderResponse{}, err
	}
	return shared.DesignDiagramRenderResponse{
		Asset: zasobWytworzonyKontraktu(zasob), MediaType: typTresci,
		UnplacedNodeIds: uporzadkujBilansDesignu(nieumieszczone),
	}, nil
}

// polozenieWezlaSchematuDesignu to prostokąt jednego węzła na płótnie: położenie
// lewego górnego rogu wraz z szerokością i wysokością.
type polozenieWezlaSchematuDesignu struct {
	X         float64
	Y         float64
	Szerokosc float64
	Wysokosc  float64
}

// ulozSchematDesignu liczy położenia węzłów i oddaje wykaz nieumieszczonych
// wraz z wymiarami płótna. Układ jest warstwowy: węzeł bez nadrzędnego stoi
// w warstwie zerowej, węzeł z nadrzędnym — warstwę niżej. Nieumieszczony wraca
// w bilansie.
func ulozSchematDesignu(z shared.DesignDiagramRenderRequest) (
	map[string]polozenieWezlaSchematuDesignu, []string, float64, float64) {

	const (
		szerokoscWezla = 160.0
		wysokoscWezla  = 56.0
		odstepPoziomy  = 40.0
		odstepPionowy  = 72.0
		marginesPlotna = 32.0
	)
	istnieje := map[string]bool{}
	for _, wezel := range z.Nodes {
		istnieje[wezel.Id] = true
	}

	// Warstwa węzła: liczba kroków do korzenia. Cykl zatrzymuje się na granicy liczby węzłów.
	warstwa := map[string]int{}
	nieumieszczone := []string{}
	nadrzedny := map[string]string{}
	for _, wezel := range z.Nodes {
		if wezel.ParentId != nil && strings.TrimSpace(*wezel.ParentId) != "" {
			nadrzedny[wezel.Id] = strings.TrimSpace(*wezel.ParentId)
		}
	}
	for _, wezel := range z.Nodes {
		poziom, biezacy := 0, wezel.Id
		zapetlony := false
		for krokow := 0; krokow <= len(z.Nodes); krokow++ {
			rodzic, maRodzica := nadrzedny[biezacy]
			if !maRodzica {
				break
			}
			if !istnieje[rodzic] {
				zapetlony = true
				break
			}
			poziom++
			biezacy = rodzic
			if krokow == len(z.Nodes) {
				zapetlony = true
			}
		}
		if zapetlony {
			nieumieszczone = append(nieumieszczone, wezel.Id)
			continue
		}
		warstwa[wezel.Id] = poziom
	}
	if z.Kind == shared.DesignDiagramKindTimeline {
		// Osi czasu nadrzędność nie porządkuje — porządkuje ją kolejność wykazu.
		for _, wezel := range z.Nodes {
			if _, jest := warstwa[wezel.Id]; jest {
				warstwa[wezel.Id] = 0
			}
		}
	}

	// Kolejność w warstwie to kolejność wykazu żądania, nie porządek alfabetyczny.
	wWarstwie := map[int][]string{}
	najwyzszaWarstwa := 0
	for _, wezel := range z.Nodes {
		poziom, jest := warstwa[wezel.Id]
		if !jest {
			continue
		}
		wWarstwie[poziom] = append(wWarstwie[poziom], wezel.Id)
		if poziom > najwyzszaWarstwa {
			najwyzszaWarstwa = poziom
		}
	}
	najszersza := 0
	for _, wykaz := range wWarstwie {
		if len(wykaz) > najszersza {
			najszersza = len(wykaz)
		}
	}
	pionowo := z.Direction == nil || *z.Direction == shared.DesignLayoutDirectionVertical

	szerokoscPlotna := 2*marginesPlotna +
		float64(najszersza)*(szerokoscWezla+odstepPoziomy) - odstepPoziomy
	wysokoscPlotna := 2*marginesPlotna +
		float64(najwyzszaWarstwa+1)*(wysokoscWezla+odstepPionowy) - odstepPionowy
	if !pionowo {
		szerokoscPlotna = 2*marginesPlotna +
			float64(najwyzszaWarstwa+1)*(szerokoscWezla+odstepPoziomy) - odstepPoziomy
		wysokoscPlotna = 2*marginesPlotna +
			float64(najszersza)*(wysokoscWezla+odstepPionowy) - odstepPionowy
	}
	if szerokoscPlotna < 200 {
		szerokoscPlotna = 200
	}
	if wysokoscPlotna < 120 {
		wysokoscPlotna = 120
	}

	polozenia := make(map[string]polozenieWezlaSchematuDesignu, len(z.Nodes))
	for poziom := 0; poziom <= najwyzszaWarstwa; poziom++ {
		wykaz := wWarstwie[poziom]
		for numer, kod := range wykaz {
			if pionowo {
				polozenia[kod] = polozenieWezlaSchematuDesignu{
					X:         marginesPlotna + float64(numer)*(szerokoscWezla+odstepPoziomy),
					Y:         marginesPlotna + float64(poziom)*(wysokoscWezla+odstepPionowy),
					Szerokosc: szerokoscWezla, Wysokosc: wysokoscWezla,
				}
				continue
			}
			polozenia[kod] = polozenieWezlaSchematuDesignu{
				X:         marginesPlotna + float64(poziom)*(szerokoscWezla+odstepPoziomy),
				Y:         marginesPlotna + float64(numer)*(wysokoscWezla+odstepPionowy),
				Szerokosc: szerokoscWezla, Wysokosc: wysokoscWezla,
			}
		}
	}
	return polozenia, nieumieszczone, szerokoscPlotna, wysokoscPlotna
}

// narysujSchematDesignu rysuje węzły i połączenia schematu na wskazanych
// położeniach, połączenia pod węzłami.
func narysujSchematDesignu(kontekst *canvas.Context, z shared.DesignDiagramRenderRequest,
	polozenia map[string]polozenieWezlaSchematuDesignu, barwy []string,
	wysokoscPlotna float64) error {

	// Połączenia idą przed węzłami, żeby linie schodziły pod nimi, a nie nad podpisami.
	for _, polaczenie := range z.Edges {
		od, maOd := polozenia[strings.TrimSpace(polaczenie.FromId)]
		do, maDo := polozenia[strings.TrimSpace(polaczenie.ToId)]
		if !maOd || !maDo {
			// Połączenie do węzła nieumieszczonego nie jest rysowane — węzeł już wrócił w unplacedNodeIds.
			continue
		}
		linia := &canvas.Path{}
		linia.MoveTo(od.X+od.Szerokosc/2, wysokoscPlotna-(od.Y+od.Wysokosc))
		linia.LineTo(do.X+do.Szerokosc/2, wysokoscPlotna-do.Y)
		kontekst.RenderPath(linia, zObrysemDanychDesignu(canvas.Style{
			Stroke: canvas.Paint{Color: barwaRgbaZapisuDesignu("#57606a", nil)}, StrokeWidth: 1.5,
		}), canvas.Identity)
		if polaczenie.Label != nil && strings.TrimSpace(*polaczenie.Label) != "" {
			srodekX := (od.X + do.X + od.Szerokosc) / 2
			srodekY := wysokoscPlotna - (od.Y+od.Wysokosc+do.Y)/2
			if err := napisWyrysuDanychDesignu(kontekst, strings.TrimSpace(*polaczenie.Label),
				srodekX, srodekY, rozmiarPodpisuDanychDesignu*0.9, "#57606a"); err != nil {
				return err
			}
		}
	}

	for numer, wezel := range z.Nodes {
		polozenie, jest := polozenia[wezel.Id]
		if !jest {
			continue
		}
		barwa := barwy[numer%len(barwy)]
		if wezel.Color != nil && strings.TrimSpace(*wezel.Color) != "" {
			barwa = *wezel.Color
		}
		// Oś Y płótna rośnie w górę, układ liczy od góry — każdy węzeł odbija się względem wysokości płótna.
		y := wysokoscPlotna - polozenie.Y - polozenie.Wysokosc
		ksztaltWezla := ksztaltWezlaSchematuDesignu(wezel.Shape, polozenie.Szerokosc,
			polozenie.Wysokosc).Translate(polozenie.X, y)
		kontekst.RenderPath(ksztaltWezla, zObrysemDanychDesignu(canvas.Style{
			Fill:        canvas.Paint{Color: barwaRgbaZapisuDesignu(barwa, wskazUlamekDesignu(0.16))},
			Stroke:      canvas.Paint{Color: barwaRgbaZapisuDesignu(barwa, nil)},
			StrokeWidth: 1.5, FillRule: canvas.NonZero,
		}), canvas.Identity)
		if err := napisWyrysuDanychDesignu(kontekst, wezel.Label,
			polozenie.X+12, y+polozenie.Wysokosc/2-4,
			rozmiarPodpisuDanychDesignu, "#1f2328"); err != nil {
			return err
		}
	}
	return nil
}

// ksztaltWezlaSchematuDesignu składa ksztaltWezla węzła. Brak wskazania bierze
// prostokąt — ksztaltWezla, w którym podpis mieści się najlepiej.
func ksztaltWezlaSchematuDesignu(rodzaj *shared.DesignShapeKind,
	szerokosc, wysokosc float64) *canvas.Path {

	if rodzaj == nil {
		return canvas.RoundedRectangle(szerokosc, wysokosc, 6)
	}
	switch *rodzaj {
	case shared.DesignShapeKindEllipse:
		return canvas.Ellipse(szerokosc/2, wysokosc/2).Translate(szerokosc/2, wysokosc/2)
	case shared.DesignShapeKindPolygon:
		return canvas.RegularPolygon(6, mniejszaDesignu(szerokosc, wysokosc)/2, true).
			Translate(szerokosc/2, wysokosc/2)
	case shared.DesignShapeKindStar:
		promien := mniejszaDesignu(szerokosc, wysokosc) / 2
		return canvas.StarPolygon(5, promien, promien*0.45, true).
			Translate(szerokosc/2, wysokosc/2)
	case shared.DesignShapeKindLine:
		linia := &canvas.Path{}
		linia.MoveTo(0, wysokosc/2)
		linia.LineTo(szerokosc, wysokosc/2)
		return linia
	}
	return canvas.RoundedRectangle(szerokosc, wysokosc, 6)
}

// wskazUlamekDesignu oddaje wskaźnik na ułamek — pola krycia są wskaźnikami,
// a literału adresu wziąć nie można.
func wskazUlamekDesignu(wartosc float64) *float64 {
	return &wartosc
}
