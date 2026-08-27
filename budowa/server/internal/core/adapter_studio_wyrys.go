// Wyrys strony dokumentu Studia — silnik `studio.preview.render`
// i `studio.diff.visual`. Rysuje bibliotekami wkompilowanymi, nigdy
// przeglądarką bezgłową ani krojem systemowym, bo obie zależności
// wykraczają poza instalkę.
package core

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Wymiary strony A4 w pikselach przy 96 punktach na cal — rozdzielczość,
// w której przeglądarka mierzy stronę A4. To nastawa domyślna wyrysu, brana
// dla dokumentu, który własnych nastaw strony nie ma.
const (
	szerokoscStronyStudia = 794
	wysokoscStronyStudia  = 1123
	marginesStronyStudia  = 64
	rozmiarPismaStudia    = 15
	interliniaStudia      = 22
	// granicaStronWyrysu chroni przed wyrysem, który nigdy się nie kończy:
	// dokument wklejony z tysiąca stron zająłby magazyn i czas rdzenia, a
	// podglądu i tak nikt nie przewinie do końca.
	granicaStronWyrysu = 200
	// punktowNaCalWyrysu jest rozdzielczością, w której wyrys mierzy stronę.
	// Ta sama, którą przyjmuje krój (`krojWyrysuStudia`) — inna dawałaby stopień
	// pisma niezgodny z kartką.
	punktowNaCalWyrysu = 96.0
)

// geometriaStronyStudia to wymiary kartki i marginesów wyrysu w pikselach,
// wyliczone z nastaw sekcji zamiast stałych A4, bo od nich zależy łamanie
// wiersza i przez nie numeracja stron w spisie treści i w indeksie.
type geometriaStronyStudia struct {
	szerokosc     int
	wysokosc      int
	marginesGorny int
	marginesDolny int
	marginesLewy  int
	marginesPrawy int
}

// geometriaDomyslnaStudia oddaje kartkę A4 o marginesach wyrysu — nastawę
// dokumentu, który własnych nastaw strony nie ma.
func geometriaDomyslnaStudia() geometriaStronyStudia {
	return geometriaStronyStudia{
		szerokosc: szerokoscStronyStudia, wysokosc: wysokoscStronyStudia,
		marginesGorny: marginesStronyStudia, marginesDolny: marginesStronyStudia,
		marginesLewy: marginesStronyStudia, marginesPrawy: marginesStronyStudia,
	}
}

// pikseleZMilimetrowWyrysu przelicza wymiar materiału na piksele wyrysu,
// rozdzielczością `punktowNaCalWyrysu`. Wymiar niedodatni oddaje zero.
func pikseleZMilimetrowWyrysu(milimetry float64) int {
	if milimetry <= 0 {
		return 0
	}
	return int(milimetry*punktowNaCalWyrysu/milimetryNaCal + 0.5)
}

// geometriaZNastawStrony liczy geometrię wyrysu z nastaw strony. Nastawy
// niepodane biorą wartość domyślną pole po polu, a wymiar własny ma
// pierwszeństwo przed nośnikiem nazwanym.
func geometriaZNastawStrony(nastawy *shared.StudioPageSetup) geometriaStronyStudia {
	geometria := geometriaDomyslnaStudia()
	if nastawy == nil {
		return geometria
	}

	szerokoscMm, wysokoscMm := 0.0, 0.0
	if nastawy.WidthMm != nil && *nastawy.WidthMm > 0 &&
		nastawy.HeightMm != nil && *nastawy.HeightMm > 0 {

		szerokoscMm, wysokoscMm = *nastawy.WidthMm, *nastawy.HeightMm
	} else if nastawy.PageSize != nil {
		if nosnik, jest := nosnikDrukuONazwie(*nastawy.PageSize); jest {
			szerokoscMm, wysokoscMm = nosnik.SzerokoscMm, nosnik.WysokoscMm
		}
	}
	if szerokoscMm > 0 && wysokoscMm > 0 {
		if nastawy.Orientation != nil && *nastawy.Orientation == shared.StudioPageOrientationPozioma {
			szerokoscMm, wysokoscMm = wysokoscMm, szerokoscMm
		}
		geometria.szerokosc = pikseleZMilimetrowWyrysu(szerokoscMm)
		geometria.wysokosc = pikseleZMilimetrowWyrysu(wysokoscMm)
	}

	oprawa := 0
	if nastawy.GutterMm != nil && *nastawy.GutterMm > 0 {
		oprawa = pikseleZMilimetrowWyrysu(*nastawy.GutterMm)
	}
	if nastawy.MarginTop != nil && *nastawy.MarginTop > 0 {
		geometria.marginesGorny = pikseleZMilimetrowWyrysu(float64(*nastawy.MarginTop))
	}
	if nastawy.MarginBottom != nil && *nastawy.MarginBottom > 0 {
		geometria.marginesDolny = pikseleZMilimetrowWyrysu(float64(*nastawy.MarginBottom))
	}
	if nastawy.MarginLeft != nil && *nastawy.MarginLeft > 0 {
		geometria.marginesLewy = pikseleZMilimetrowWyrysu(float64(*nastawy.MarginLeft))
	}
	if nastawy.MarginRight != nil && *nastawy.MarginRight > 0 {
		geometria.marginesPrawy = pikseleZMilimetrowWyrysu(float64(*nastawy.MarginRight))
	}
	geometria.marginesLewy += oprawa

	return geometria.uzdrowiona()
}

// uzdrowiona pilnuje, żeby geometria nie zamieniła się w kartkę bez miejsca
// na litery: marginesy zjadające całą kartkę zwęża zamiast dzielić przez
// zero, a wyrys wychodzi widocznie za wąski.
func (g geometriaStronyStudia) uzdrowiona() geometriaStronyStudia {
	const najmniejszaKolumna = 32
	if g.szerokosc < najmniejszaKolumna {
		g.szerokosc = szerokoscStronyStudia
	}
	if g.wysokosc < interliniaStudia*3 {
		g.wysokosc = wysokoscStronyStudia
	}
	for g.szerokosc-g.marginesLewy-g.marginesPrawy < najmniejszaKolumna {
		if g.marginesLewy == 0 && g.marginesPrawy == 0 {
			break
		}
		g.marginesLewy /= 2
		g.marginesPrawy /= 2
	}
	for g.wysokosc-g.marginesGorny-g.marginesDolny < interliniaStudia*2 {
		if g.marginesGorny == 0 && g.marginesDolny == 0 {
			break
		}
		g.marginesGorny /= 2
		g.marginesDolny /= 2
	}
	return g
}

// szerokoscKolumny oddaje szerokość kolumny tekstu w pikselach — kartkę
// pomniejszoną o oba marginesy poziome.
func (g geometriaStronyStudia) szerokoscKolumny() int {
	szerokosc := g.szerokosc - g.marginesLewy - g.marginesPrawy
	if szerokosc < 1 {
		szerokosc = 1
	}
	return szerokosc
}

// wierszyNaStrone oddaje liczbę wierszy mieszczących się na jednej kartce.
// Odjęcie dwóch interlinii jest miejscem nagłówka i stopki, które stoją poza
// kolumną tekstu, ale kartkę zajmują.
func (g geometriaStronyStudia) wierszyNaStrone() int {
	naStrone := (g.wysokosc - g.marginesGorny - g.marginesDolny - 2*interliniaStudia) /
		interliniaStudia
	if naStrone < 1 {
		naStrone = 1
	}
	return naStrone
}

// nastawyWyrysuStudia opisuje jedną stronę wyrysu wraz z jej wycinkiem
// i geometrią kartki, nagłówkiem, stopką i znakiem wodnym.
type nastawyWyrysuStudia struct {
	naglowek  string
	stopka    string
	znakWodny string
	// stronaOd i stronaDo zawężają wyrys do wycinka; zero znaczy bez zawężenia.
	stronaOd, stronaDo int
	// geometria niesie wymiary kartki; wartość zerowa znaczy nastawę domyślną.
	geometria geometriaStronyStudia
}

// kartka oddaje geometrię nastaw albo, gdy jej nie podano, geometrię
// domyślną strony A4 o marginesach wyrysu.
func (n nastawyWyrysuStudia) kartka() geometriaStronyStudia {
	if n.geometria.szerokosc <= 0 || n.geometria.wysokosc <= 0 {
		return geometriaDomyslnaStudia()
	}
	return n.geometria
}

// kartaWyrysuStudia to jedna wyrysowana strona wraz z jej numerem w wykazie
// stron, jaki oddaje `wyrysujStronyStudia`.
type kartaWyrysuStudia struct {
	numer int
	obraz *image.RGBA
}

// bladWyrysuStudia nazywa usterkę wyrysu kodem kontraktu ErrorCodeInternalError
// wraz z opisem powodu, jednolitym dla całego pliku.
func bladWyrysuStudia(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Studio: wyrys strony — "+powod))
}

// krojWyrysuStudia składa krój pisma `gofont/goregular` o zadanym stopniu,
// przy rozdzielczości 96 punktów na cal.
func krojWyrysuStudia(stopien float64) (font.Face, error) {
	krojka, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, bladWyrysuStudia("nieczytelny krój wkompilowany: " + err.Error())
	}
	oblicze, err := opentype.NewFace(krojka, &opentype.FaceOptions{
		Size: stopien, DPI: 96, Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, bladWyrysuStudia("nie można złożyć kroju: " + err.Error())
	}
	return oblicze, nil
}

// wyrysujStronyStudia rozkłada treść na strony i rysuje każdą z nich. Łamanie
// wiersza idzie po słowach i po zmierzonej szerokości, nie po stałej liczbie
// znaków.
func wyrysujStronyStudia(tresc string, n nastawyWyrysuStudia) ([]kartaWyrysuStudia, error) {
	oblicze, err := krojWyrysuStudia(rozmiarPismaStudia)
	if err != nil {
		return nil, err
	}
	defer oblicze.Close()

	kartka := n.kartka()
	wiersze := []string{}
	for _, akapit := range strings.Split(tresc, "\n") {
		wiersze = append(wiersze, zlamWierszStudia(akapit, oblicze, kartka.szerokoscKolumny())...)
	}
	if len(wiersze) == 0 {
		wiersze = []string{""}
	}

	naStrone := kartka.wierszyNaStrone()
	stron := (len(wiersze) + naStrone - 1) / naStrone
	if stron > granicaStronWyrysu {
		stron = granicaStronWyrysu
	}

	od, do_ := 1, stron
	if n.stronaOd > 0 {
		od = n.stronaOd
	}
	if n.stronaDo > 0 && n.stronaDo < do_ {
		do_ = n.stronaDo
	}
	if od > stron {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
			"moduł Studio: wyrys od strony "+strconv.Itoa(od)+", a dokument ma stron "+
				strconv.Itoa(stron)))
	}

	karty := make([]kartaWyrysuStudia, 0, do_-od+1)
	for numer := od; numer <= do_; numer++ {
		poczatek := (numer - 1) * naStrone
		koniec := poczatek + naStrone
		if koniec > len(wiersze) {
			koniec = len(wiersze)
		}
		karty = append(karty, kartaWyrysuStudia{
			numer: numer,
			obraz: narysujStroneStudia(wiersze[poczatek:koniec], numer, stron, oblicze, n),
		})
	}
	return karty, nil
}

// zlamWierszStudia rozkłada jeden akapit na wiersze mieszczące się w kolumnie
// o podanej szerokości, mierzonej krojem `oblicze`.
func zlamWierszStudia(akapit string, oblicze font.Face, szerokosc int) []string {
	akapit = strings.ReplaceAll(akapit, "\t", "    ")
	if strings.TrimSpace(akapit) == "" {
		return []string{""}
	}
	wiersze := []string{}
	biezacy := ""
	for _, slowo := range strings.Fields(akapit) {
		proba := slowo
		if biezacy != "" {
			proba = biezacy + " " + slowo
		}
		if font.MeasureString(oblicze, proba).Ceil() <= szerokosc {
			biezacy = proba
			continue
		}
		if biezacy != "" {
			wiersze = append(wiersze, biezacy)
		}
		// Słowo dłuższe od kolumny łamie się po znakach, inaczej wyszłoby
		// poza margines.
		for font.MeasureString(oblicze, slowo).Ceil() > szerokosc {
			ciecie := len([]rune(slowo))
			for ciecie > 1 && font.MeasureString(oblicze, string([]rune(slowo)[:ciecie])).Ceil() > szerokosc {
				ciecie--
			}
			wiersze = append(wiersze, string([]rune(slowo)[:ciecie]))
			slowo = string([]rune(slowo)[ciecie:])
		}
		biezacy = slowo
	}
	if biezacy != "" {
		wiersze = append(wiersze, biezacy)
	}
	return wiersze
}

// narysujStroneStudia rysuje jedną stronę: nagłówek, treść, stopkę z numeracją
// i — gdy Operator o to poprosił — znak wodny.
func narysujStroneStudia(wiersze []string, numer, stron int, oblicze font.Face,
	n nastawyWyrysuStudia) *image.RGBA {

	kartka := n.kartka()
	strona := image.NewRGBA(image.Rect(0, 0, kartka.szerokosc, kartka.wysokosc))
	draw.Draw(strona, strona.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	// Znak wodny idzie pod treść, żeby zostać widoczny, a nie zasłaniający.
	if znak := strings.TrimSpace(n.znakWodny); znak != "" {
		rysujTekstStudia(strona, oblicze, znak, kartka.marginesLewy,
			kartka.wysokosc/2, color.RGBA{R: 226, G: 226, B: 226, A: 255})
	}

	szary := color.RGBA{R: 128, G: 128, B: 128, A: 255}
	czarny := color.RGBA{A: 255}
	if naglowek := strings.TrimSpace(n.naglowek); naglowek != "" {
		rysujTekstStudia(strona, oblicze, naglowek, kartka.marginesLewy,
			kartka.marginesGorny-interliniaStudia/2, szary)
	}

	y := kartka.marginesGorny + interliniaStudia
	for _, wiersz := range wiersze {
		rysujTekstStudia(strona, oblicze, wiersz, kartka.marginesLewy, y, czarny)
		y += interliniaStudia
	}

	stopka := strings.TrimSpace(n.stopka)
	numeracja := strconv.Itoa(numer) + " / " + strconv.Itoa(stron)
	if stopka != "" {
		numeracja = stopka + "   ·   " + numeracja
	}
	rysujTekstStudia(strona, oblicze, numeracja, kartka.marginesLewy,
		kartka.wysokosc-kartka.marginesDolny/2, szary)
	return strona
}

// rysujTekstStudia kładzie jeden wiersz tekstu w podanym miejscu na płótnie,
// zadanym krojem i barwą, bez łamania.
func rysujTekstStudia(plotno *image.RGBA, oblicze font.Face, tekst string,
	x, y int, barwa color.Color) {

	rysownik := &font.Drawer{
		Dst: plotno, Src: image.NewUniform(barwa), Face: oblicze,
		Dot: fixed.P(x, y),
	}
	rysownik.DrawString(tekst)
}

// pustaStronaStudia oddaje białą kartę o wymiarach podanej kartki. Wchodzi za
// stronę, której druga porównywana wersja nie ma, żeby miała z czym się
// porównać.
func pustaStronaStudia(kartka geometriaStronyStudia) *image.RGBA {
	strona := image.NewRGBA(image.Rect(0, 0, kartka.szerokosc, kartka.wysokosc))
	draw.Draw(strona, strona.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	return strona
}

// pngZeStronyStudia koduje stronę do PNG — formatu, który `pdfZeStronStudia`
// składa dalej w dokument PDF.
func pngZeStronyStudia(obraz image.Image) ([]byte, error) {
	var bufor bytes.Buffer
	if err := png.Encode(&bufor, obraz); err != nil {
		return nil, bladWyrysuStudia("nie można zakodować strony: " + err.Error())
	}
	return bufor.Bytes(), nil
}

// pdfZeStronStudia składa strony wyrysu w jeden dokument PDF drogą przez
// obrazy, a nie przez generator tekstu w PDF, żeby PDF pokazywał dokładnie
// to, co Operator widział w podglądzie.
func pdfZeStronStudia(strony [][]byte) ([]byte, error) {
	if len(strony) == 0 {
		return nil, bladWyrysuStudia("złożenie dokumentu bez ani jednej strony")
	}
	obrazy := make([]io.Reader, 0, len(strony))
	for _, strona := range strony {
		obrazy = append(obrazy, bytes.NewReader(strona))
	}
	nastawy := pdfcpu.DefaultImportConfig()
	var wynik bytes.Buffer
	if err := api.ImportImages(nil, &wynik, obrazy, nastawy, nastawyPdf()); err != nil {
		return nil, bladWyrysuStudia("nie można złożyć dokumentu ze stron: " + err.Error())
	}
	return wynik.Bytes(), nil
}

// ── Porównanie wizualne ─────────────────────────────────────────────────────

// bokKratkiRoznicyStudia to bok komórki, w której liczy się udział pikseli
// różniących się. Komórka mniejsza dałaby wykaz obszarów tak długi, że nikt by
// go nie przejrzał; większa scaliłaby zmianę akapitu ze zmianą marginesu.
const bokKratkiRoznicyStudia = 24

// obszaryRoznicyStudia porównuje dwie strony komórka po komórce, po jasności
// pikseli, a nie po składowych barwy, bo podgląd jest czarny na białym.
func obszaryRoznicyStudia(bazowa, docelowa image.Image, numer int) []shared.StudioVisualDiffRegion {
	prostokat := bazowa.Bounds().Intersect(docelowa.Bounds())
	obszary := []shared.StudioVisualDiffRegion{}
	for y := prostokat.Min.Y; y < prostokat.Max.Y; y += bokKratkiRoznicyStudia {
		for x := prostokat.Min.X; x < prostokat.Max.X; x += bokKratkiRoznicyStudia {
			doX := x + bokKratkiRoznicyStudia
			if doX > prostokat.Max.X {
				doX = prostokat.Max.X
			}
			doY := y + bokKratkiRoznicyStudia
			if doY > prostokat.Max.Y {
				doY = prostokat.Max.Y
			}
			rozne, wszystkie := 0, 0
			for py := y; py < doY; py++ {
				for px := x; px < doX; px++ {
					wszystkie++
					if jasnoscStudia(bazowa, px, py) != jasnoscStudia(docelowa, px, py) {
						rozne++
					}
				}
			}
			if rozne == 0 || wszystkie == 0 {
				continue
			}
			obszary = append(obszary, shared.StudioVisualDiffRegion{
				Page: numer, X: x, Y: y, Width: doX - x, Height: doY - y,
				ChangeRatio: float64(rozne) / float64(wszystkie),
			})
		}
	}
	return obszary
}

// dopasujStroneStudia sprowadza stronę do wymiarów strony porównywanej,
// filtrem CatmullRom — najbliższy sąsiad dałby na tekście szum krawędzi liter.
func dopasujStroneStudia(zrodlo image.Image, prostokat image.Rectangle) image.Image {
	if zrodlo.Bounds() == prostokat {
		return zrodlo
	}
	dopasowana := image.NewRGBA(prostokat)
	draw.CatmullRom.Scale(dopasowana, prostokat, zrodlo, zrodlo.Bounds(), draw.Src, nil)
	return dopasowana
}

// jasnoscStudia sprowadza piksel do jednej liczby w skali 0–255, ważoną sumą
// jego składowych barwy czerwonej, zielonej i niebieskiej.
func jasnoscStudia(obraz image.Image, x, y int) uint8 {
	r, g, b, _ := obraz.At(x, y).RGBA()
	return uint8((r*299 + g*587 + b*114) / 1000 >> 8)
}

// nakladkaRoznicyStudia rysuje stronę docelową z obszarami różnicy zaznaczonymi
// na czerwono. Nakładka jest osobnym obrazem, a nie zmianą strony: strona
// wyrysu ma zostać stroną wyrysu.
func nakladkaRoznicyStudia(docelowa image.Image,
	obszary []shared.StudioVisualDiffRegion) *image.RGBA {

	prostokat := docelowa.Bounds()
	nakladka := image.NewRGBA(prostokat)
	draw.Draw(nakladka, prostokat, docelowa, prostokat.Min, draw.Src)
	czerwien := image.NewUniform(color.RGBA{R: 220, G: 38, B: 38, A: 70})
	for _, obszar := range obszary {
		pole := image.Rect(obszar.X, obszar.Y, obszar.X+obszar.Width, obszar.Y+obszar.Height)
		draw.DrawMask(nakladka, pole.Intersect(prostokat), czerwien, image.Point{},
			nil, image.Point{}, draw.Over)
	}
	return nakladka
}
