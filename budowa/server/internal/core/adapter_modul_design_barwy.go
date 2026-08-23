// Odpowiedzialność pliku: rachunek barwy modułu Design — rozpoznanie zapisu
// wejściowego, przeliczenie między przestrzeniami i współczynnik kontrastu
// wedle WCAG 2.1. Czynności kontraktu stoją w `adapter_modul_design_kolor.go`
// i `_kolor_obraz.go`; tutaj leży sam rachunek, żeby miał jedno miejsce i jedną
// prawdę.
//
// ── Rachunek jest wkompilowany, nie wołany ──────────────────────────────────
// Konwersje przestrzeni percepcyjnych i harmonie idą przez
// `lucasb-eyer/go-colorful`, bibliotekę Go wkompilowaną w binarium. Programu
// zewnętrznego do przekształceń obrazu w tej drodze nie ma i mieć nie będzie:
// funkcja zależna od programu spoza instalki jest u Operatora odmową, nie
// funkcją.
//
// ── Luminancja liczy się wzorem WCAG, nie „jasnością" z biblioteki ──────────
// Współczynnik kontrastu WCAG 2.1 stoi na luminancji względnej liczonej ze
// składowych sRGB po zdjęciu gamma (`LinearRgb`), z wagami 0.2126 / 0.7152 /
// 0.0722. Jasność Lab (`L*`) jest inną wielkością i dałaby inne liczby — para
// czerni na bieli ma dawać dokładnie 21, bo taki jest kres tej skali, i po tej
// liczbie sprawdzian poznaje, że rachunek jest ten, o który chodzi.
//
// ── CMYK jest przeliczeniem wprost, i mówimy to wprost ──────────────────────
// Przeliczenie RGB → CMYK bez profilu ICC jest przeliczeniem naiwnym: oddaje
// wartości, których drukarnia użyje jako punktu wyjścia, a nie barwę
// rozdzieloną pod konkretną maszynę. Rdzeń nie udaje, że zna profil, którego
// nie dostał — pole `cmyk` niesie przeliczenie wprost i tyle.
package core

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"

	"danacoconsole/shared"
)

// barwyNazwaneDesignu to próbki nazwane rozpoznawane przy wejściu i oddawane
// w polu `name`, gdy barwa trafia w nie dokładnie.
//
// Wykaz jest krótki z zamysłu: to są nazwy, które mają w produkcie znaczenie
// (barwy podstawowe CSS poziomu 1 plus czerń, biel i szarości). Wciągnięcie
// pełnej listy X11 dałoby nazwy w rodzaju „papayawhip", których nikt w module
// nie wpisze, a każda z nich musiałaby być tu utrzymywana.
var barwyNazwaneDesignu = map[string]string{
	"black":   "#000000",
	"white":   "#ffffff",
	"gray":    "#808080",
	"grey":    "#808080",
	"silver":  "#c0c0c0",
	"red":     "#ff0000",
	"lime":    "#00ff00",
	"green":   "#008000",
	"blue":    "#0000ff",
	"yellow":  "#ffff00",
	"cyan":    "#00ffff",
	"aqua":    "#00ffff",
	"magenta": "#ff00ff",
	"fuchsia": "#ff00ff",
	"maroon":  "#800000",
	"olive":   "#808000",
	"navy":    "#000080",
	"teal":    "#008080",
	"purple":  "#800080",
	"orange":  "#ffa500",
}

// rozpoznajBarweDesignu przekłada dowolny zapis barwy na kolor biblioteki.
//
// Rozpoznaje: `#rgb`, `#rrggbb`, `rgb(...)`, `hsl(...)`, `lab(...)`,
// `cmyk(...)` oraz próbkę nazwaną. Zapis nierozpoznany jest ODMOWĄ, nie barwą
// domyślną: podstawienie czerni za tekst, którego rdzeń nie zrozumiał, dałoby
// paletę zbudowaną wokół barwy, której nikt nie wskazał.
func rozpoznajBarweDesignu(zapis string) (colorful.Color, error) {
	tekst := strings.ToLower(strings.TrimSpace(zapis))
	if tekst == "" {
		return colorful.Color{}, fmt.Errorf("zapis barwy jest pusty")
	}
	if hex, nazwana := barwyNazwaneDesignu[tekst]; nazwana {
		tekst = hex
	}
	if strings.HasPrefix(tekst, "#") {
		return rozpoznajBarweSzesnastkowoDesignu(tekst)
	}
	otwarcie := strings.IndexByte(tekst, '(')
	if otwarcie > 0 && strings.HasSuffix(tekst, ")") {
		przestrzen := strings.TrimSpace(tekst[:otwarcie])
		liczby, err := liczbyZapisuBarwyDesignu(tekst[otwarcie+1 : len(tekst)-1])
		if err != nil {
			return colorful.Color{}, err
		}
		return barwaZPrzestrzeniDesignu(przestrzen, liczby)
	}
	// Zapis szesnastkowy bez krzyżyka jest częstym skrótem Operatora i tak samo
	// jednoznacznym — rozpoznajemy go, zamiast odmawiać za znak.
	if len(tekst) == 3 || len(tekst) == 6 {
		if _, err := strconv.ParseUint(tekst, 16, 32); err == nil {
			return rozpoznajBarweSzesnastkowoDesignu("#" + tekst)
		}
	}
	return colorful.Color{}, fmt.Errorf(
		"zapisu %q nie da się odczytać jako barwy; rdzeń przyjmuje #rrggbb, rgb(), hsl(), "+
			"lab(), cmyk() oraz próbki nazwane", zapis)
}

// rozpoznajBarweSzesnastkowoDesignu rozwija skrót `#rgb` do `#rrggbb`
// i oddaje kolor. Skrót trzyznakowy jest zapisem CSS, więc odmowa za niego
// byłaby odmową za zapis poprawny.
func rozpoznajBarweSzesnastkowoDesignu(tekst string) (colorful.Color, error) {
	if len(tekst) == 4 {
		tekst = "#" + strings.Repeat(string(tekst[1]), 2) +
			strings.Repeat(string(tekst[2]), 2) + strings.Repeat(string(tekst[3]), 2)
	}
	barwa, err := colorful.Hex(tekst)
	if err != nil {
		return colorful.Color{}, fmt.Errorf("zapis %q nie jest poprawną barwą szesnastkową", tekst)
	}
	return barwa, nil
}

// liczbyZapisuBarwyDesignu rozkłada wnętrze nawiasu na liczby. Rozdzielnikiem
// bywa przecinek albo spacja (składnia CSS Color 4), a odsetek zamienia się na
// ułamek od razu — inaczej `hsl(210, 50%, 40%)` i `hsl(210 0.5 0.4)` znaczyłyby
// co innego, choć to ten sam zapis.
func liczbyZapisuBarwyDesignu(wnetrze string) ([]float64, error) {
	pola := strings.FieldsFunc(wnetrze, func(znak rune) bool {
		return znak == ',' || znak == ' ' || znak == '/' || znak == '\t'
	})
	liczby := make([]float64, 0, len(pola))
	for _, pole := range pola {
		pole = strings.TrimSpace(pole)
		if pole == "" {
			continue
		}
		odsetek := strings.HasSuffix(pole, "%")
		wartosc, err := strconv.ParseFloat(strings.TrimSuffix(pole, "%"), 64)
		if err != nil {
			return nil, fmt.Errorf("składowa %q zapisu barwy nie jest liczbą", pole)
		}
		if odsetek {
			wartosc /= 100
		}
		liczby = append(liczby, wartosc)
	}
	return liczby, nil
}

// barwaZPrzestrzeniDesignu składa kolor z liczb odczytanych dla wskazanej
// przestrzeni.
func barwaZPrzestrzeniDesignu(przestrzen string, liczby []float64) (colorful.Color, error) {
	switch przestrzen {
	case "rgb", "rgba":
		if len(liczby) < 3 {
			return colorful.Color{}, fmt.Errorf("zapis rgb() wymaga trzech składowych")
		}
		// Składowe rgb() bywają podane jako 0–255 albo jako odsetek zamieniony
		// wyżej na ułamek. Rozróżnia je kres: wartość powyżej jedynki nie może
		// być ułamkiem sRGB.
		skala := 255.0
		if liczby[0] <= 1 && liczby[1] <= 1 && liczby[2] <= 1 {
			skala = 1.0
		}
		return colorful.Color{
			R: przytnijUlamekDesignu(liczby[0] / skala),
			G: przytnijUlamekDesignu(liczby[1] / skala),
			B: przytnijUlamekDesignu(liczby[2] / skala),
		}, nil
	case "hsl", "hsla":
		if len(liczby) < 3 {
			return colorful.Color{}, fmt.Errorf("zapis hsl() wymaga trzech składowych")
		}
		return colorful.Hsl(znormalizujKatDesignu(liczby[0]),
			przytnijUlamekDesignu(liczby[1]), przytnijUlamekDesignu(liczby[2])).Clamped(), nil
	case "lab":
		if len(liczby) < 3 {
			return colorful.Color{}, fmt.Errorf("zapis lab() wymaga trzech składowych")
		}
		// L* w zapisie CSS idzie 0–100, a biblioteka liczy go 0–1. Ułamek
		// podany wprost też jest tu poprawny i wchodzi bez skalowania.
		jasnosc := liczby[0]
		if jasnosc > 1 {
			jasnosc /= 100
		}
		return colorful.Lab(jasnosc, liczby[1]/100, liczby[2]/100).Clamped(), nil
	case "cmyk":
		if len(liczby) < 4 {
			return colorful.Color{}, fmt.Errorf("zapis cmyk() wymaga czterech składowych")
		}
		return barwaZCmykDesignu(liczby[0], liczby[1], liczby[2], liczby[3]), nil
	}
	return colorful.Color{}, fmt.Errorf("przestrzeń %q nie jest znana rdzeniowi", przestrzen)
}

// barwaZCmykDesignu przelicza CMYK na sRGB wprost — bez profilu ICC, bo żadnego
// nie dostaliśmy (nagłówek pliku).
func barwaZCmykDesignu(c, m, y, k float64) colorful.Color {
	if c > 1 || m > 1 || y > 1 || k > 1 {
		c, m, y, k = c/100, m/100, y/100, k/100
	}
	return colorful.Color{
		R: przytnijUlamekDesignu((1 - c) * (1 - k)),
		G: przytnijUlamekDesignu((1 - m) * (1 - k)),
		B: przytnijUlamekDesignu((1 - y) * (1 - k)),
	}
}

// cmykBarwyDesignu przelicza sRGB na CMYK wprost i oddaje zapis tekstowy.
func cmykBarwyDesignu(barwa colorful.Color) string {
	k := 1 - math.Max(barwa.R, math.Max(barwa.G, barwa.B))
	if k >= 1 {
		return "cmyk(0%, 0%, 0%, 100%)"
	}
	c := (1 - barwa.R - k) / (1 - k)
	m := (1 - barwa.G - k) / (1 - k)
	y := (1 - barwa.B - k) / (1 - k)
	return fmt.Sprintf("cmyk(%.0f%%, %.0f%%, %.0f%%, %.0f%%)", c*100, m*100, y*100, k*100)
}

// barwaKontraktuDesignu składa `DesignColorValue` — barwę wyrażoną we
// wszystkich przestrzeniach naraz, tak jak żąda `design.color.convert`.
func barwaKontraktuDesignu(barwa colorful.Color) shared.DesignColorValue {
	czysta := barwa.Clamped()
	h, s, l := czysta.Hsl()
	jasnosc, a, b := czysta.Lab()
	zapis := shared.DesignColorValue{
		Hex: czysta.Hex(),
		Rgb: fmt.Sprintf("rgb(%d, %d, %d)", zaokraglijSkladowaDesignu(czysta.R),
			zaokraglijSkladowaDesignu(czysta.G), zaokraglijSkladowaDesignu(czysta.B)),
		Hsl:  fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", znormalizujKatDesignu(h), s*100, l*100),
		Lab:  fmt.Sprintf("lab(%.1f%% %.1f %.1f)", jasnosc*100, a*100, b*100),
		Cmyk: cmykBarwyDesignu(czysta),
	}
	// Nazwa wchodzi wyłącznie przy trafieniu dokładnym. Nazwa „najbliższa"
	// mówiłaby o barwie, której w żądaniu nie było — a pole jest niewymagane
	// właśnie po to, żeby brak nazwy dało się powiedzieć wprost.
	for nazwa, hex := range barwyNazwaneDesignu {
		if hex == zapis.Hex {
			wartosc := nazwa
			zapis.Name = &wartosc
			break
		}
	}
	return zapis
}

// barwaRgbaDesignu przekłada barwę na składowe z kanałem krycia — postać, którą
// przyjmują biblioteki wyrysu (`image/draw`, `tdewolff/canvas`).
//
// Składowe są PRZEMNOŻONE przez krycie, bo `color.RGBA` biblioteki standardowej
// jest formatem z krycim wmnożonym. Wartości niepomnożone dawałyby przy kryciu
// częściowym barwę jaśniejszą, niż wskazano, i wyrys nie zgadzałby się
// z podglądem w oknie.
func barwaRgbaDesignu(barwa colorful.Color, krycie *float64) color.RGBA {
	czysta := barwa.Clamped()
	kanal := 1.0
	if krycie != nil {
		kanal = przytnijUlamekDesignu(*krycie)
	}
	return color.RGBA{
		R: uint8(math.Round(czysta.R * kanal * 255)),
		G: uint8(math.Round(czysta.G * kanal * 255)),
		B: uint8(math.Round(czysta.B * kanal * 255)),
		A: uint8(math.Round(kanal * 255)),
	}
}

// luminancjaWcagDesignu liczy luminancję względną wedle WCAG 2.1: składowe
// sRGB po zdjęciu gamma, z wagami 0.2126 / 0.7152 / 0.0722.
func luminancjaWcagDesignu(barwa colorful.Color) float64 {
	r, g, b := barwa.Clamped().LinearRgb()
	return 0.2126*r + 0.7152*g + 0.0722*b
}

// kontrastWcagDesignu liczy współczynnik kontrastu pary barw. Kres skali to 21
// (czerń na bieli) i 1 (barwa na samej sobie) — po tych dwóch liczbach poznaje
// się, że rachunek jest ten właściwy.
func kontrastWcagDesignu(pierwszy, drugi colorful.Color) float64 {
	jasniejszy := luminancjaWcagDesignu(pierwszy)
	ciemniejszy := luminancjaWcagDesignu(drugi)
	if jasniejszy < ciemniejszy {
		jasniejszy, ciemniejszy = ciemniejszy, jasniejszy
	}
	return (jasniejszy + 0.05) / (ciemniejszy + 0.05)
}

// progiKontrastuDesignu to progi WCAG 2.1 dla tekstu zwykłego i dużego.
const (
	progKontrastuAA       = 4.5
	progKontrastuAAA      = 7.0
	progKontrastuDuzegoAA = 3.0
	// Tekst duży wedle WCAG: od 18 punktów (24 px) albo od 14 punktów (18.66 px)
	// przy pogrubieniu.
	rozmiarTekstuDuzegoDesignu       = 24.0
	rozmiarTekstuDuzegoPogrubionego  = 18.66
	miejscaPoPrzecinkuKontrastuWcag  = 2
	dzielnikZaokragleniaKontrastuAAA = 100.0
)

// wynikKontrastuDesignu składa `DesignContrastResult` dla pary barw wraz
// z oceną progów.
//
// `passesLargeAA` liczy się zawsze, a nie tylko wtedy, gdy wołający podał
// rozmiar pisma: pole odpowiada na pytanie „czy ta para nadaje się na nagłówek",
// które ma sens także bez wskazania rozmiaru. Rozmiar i pogrubienie
// rozstrzygają natomiast o `passesAA` — para na 4.0 jest zgodna dla tekstu
// dużego i niezgodna dla zwykłego, więc jedna odpowiedź na oba przypadki
// byłaby nieprawdziwa dla jednego z nich.
func wynikKontrastuDesignu(pierwszyZapis, drugiZapis string, pierwszy, drugi colorful.Color,
	rozmiar *float64, pogrubienie *bool) shared.DesignContrastResult {

	wspolczynnik := math.Round(kontrastWcagDesignu(pierwszy, drugi)*
		dzielnikZaokragleniaKontrastuAAA) / dzielnikZaokragleniaKontrastuAAA

	duzy := czyTekstDuzyDesignu(rozmiar, pogrubienie)
	progAA := progKontrastuAA
	progAAA := progKontrastuAAA
	if duzy {
		progAA = progKontrastuDuzegoAA
		progAAA = progKontrastuAA
	}
	spelniaDuzy := wspolczynnik >= progKontrastuDuzegoAA
	return shared.DesignContrastResult{
		Foreground:    pierwszyZapis,
		Background:    drugiZapis,
		Ratio:         wspolczynnik,
		PassesAA:      wspolczynnik >= progAA,
		PassesAAA:     wspolczynnik >= progAAA,
		PassesLargeAA: &spelniaDuzy,
	}
}

// czyTekstDuzyDesignu rozstrzyga, czy para dotyczy tekstu dużego w rozumieniu
// WCAG. Brak wskazania rozmiaru znaczy tekst zwykły — próg łagodniejszy
// przyznany domyślnie byłby oceną wystawioną na wyrost.
func czyTekstDuzyDesignu(rozmiar *float64, pogrubienie *bool) bool {
	if rozmiar == nil || *rozmiar <= 0 {
		return false
	}
	if pogrubienie != nil && *pogrubienie {
		return *rozmiar >= rozmiarTekstuDuzegoPogrubionego
	}
	return *rozmiar >= rozmiarTekstuDuzegoDesignu
}

// przytnijUlamekDesignu wprowadza wartość w przedział 0–1.
func przytnijUlamekDesignu(wartosc float64) float64 {
	if wartosc < 0 {
		return 0
	}
	if wartosc > 1 {
		return 1
	}
	return wartosc
}

// znormalizujKatDesignu wprowadza kąt odcienia w przedział 0–360.
func znormalizujKatDesignu(kat float64) float64 {
	wynik := math.Mod(kat, 360)
	if wynik < 0 {
		wynik += 360
	}
	return wynik
}

// zaokraglijSkladowaDesignu przekłada ułamek sRGB na całkowitą składową 0–255.
func zaokraglijSkladowaDesignu(wartosc float64) int {
	return int(math.Round(przytnijUlamekDesignu(wartosc) * 255))
}
