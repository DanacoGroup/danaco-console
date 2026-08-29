// Odpowiedzialność pliku: kroje pisma modułu Design — katalog krojów,
// wczytanie kroju, glify, szerokość tekstu i zamiana tekstu na kontury.
package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gobolditalic"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/gofont/gosmallcaps"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"

	"github.com/tdewolff/canvas"

	"danacoconsole/shared"
)

const (
	// katalogiKrojowSerwera wylicza miejsca, w których leżą kroje pisma na
	// systemach, na jakich stoi serwer produktu; odczyt jest wyłącznie
	// odczytem.
	katalogKrojowSystemowy = "/usr/share/fonts"
	katalogKrojowLokalny   = "/usr/local/share/fonts"

	// granicaGlifowOdpowiedziDesignu chroni odpowiedź przed wykazem
	// dwudziestotysięcznym; kroje pełne CJK mają tyle glifów.
	granicaGlifowOdpowiedziDesignu = 512

	// domyslnyTekstProbnyDesignu jest tekstem podglądu, gdy Operator własnego
	// nie podał; zdanie polskie z ogonkami dla sprawdzenia kroju.
	domyslnyTekstProbnyDesignu = "Zażółć gęślą jaźń — ĄĆĘŁŃÓŚŹŻ 0123456789"
)

// krojeWkompilowaneDesignu to kroje leżące w binarium serwera; rodzina Go
// jest tu w całości, bo to jedyna rodzina, którą produkt ma NA PEWNO.
var krojeWkompilowaneDesignu = map[string][]byte{
	"Go Regular":     goregular.TTF,
	"Go Medium":      gomedium.TTF,
	"Go Bold":        gobold.TTF,
	"Go Italic":      goitalic.TTF,
	"Go Bold Italic": gobolditalic.TTF,
	"Go Mono":        gomono.TTF,
	"Go Mono Bold":   gomonobold.TTF,
	"Go Smallcaps":   gosmallcaps.TTF,
}

// pamiecKrojowSerweraDesignu trzyma wynik przeglądu katalogów krojów;
// przegląd idzie raz na proces, katalog nie zmienia się w trakcie pracy.
var (
	razKrojowSerweraDesignu sync.Once
	krojeSerweraDesignu     map[string]string
)

// przegladajKrojeSerweraDesignu składa wykaz krojów tej maszyny: nazwa
// rodziny na ścieżkę pliku, wziętą z nazwy pliku, nie z tablicy name.
func przegladajKrojeSerweraDesignu() map[string]string {
	razKrojowSerweraDesignu.Do(func() {
		krojeSerweraDesignu = map[string]string{}
		for _, korzen := range []string{katalogKrojowSystemowy, katalogKrojowLokalny} {
			_ = filepath.WalkDir(korzen, func(sciezka string, wpis os.DirEntry, err error) error {
				if err != nil || wpis.IsDir() {
					return nil
				}
				rozszerzenie := strings.ToLower(filepath.Ext(wpis.Name()))
				if rozszerzenie != ".ttf" && rozszerzenie != ".otf" {
					return nil
				}
				nazwa := nazwaKrojuZPlikuDesignu(wpis.Name())
				if _, zajete := krojeSerweraDesignu[nazwa]; !zajete {
					krojeSerweraDesignu[nazwa] = sciezka
				}
				return nil
			})
		}
	})
	return krojeSerweraDesignu
}

// nazwaKrojuZPlikuDesignu wyciąga nazwę rodziny z nazwy pliku, na przykład
// NimbusSans-Regular.otf daje NimbusSans-Regular.
func nazwaKrojuZPlikuDesignu(plik string) string {
	return strings.TrimSuffix(plik, filepath.Ext(plik))
}

// kluczKrojuDesignu sprowadza nazwę kroju do postaci porównywalnej: bez
// spacji, dywizów i podkreśleń, małymi literami.
func kluczKrojuDesignu(nazwa string) string {
	zamiana := strings.NewReplacer(" ", "", "-", "", "_", "", ".", "")
	return strings.ToLower(zamiana.Replace(strings.TrimSpace(nazwa)))
}

// krojDesignu wczytuje krój po nazwie i mówi, skąd pochodzi; pierwszeństwo
// mają kroje wkompilowane, dostępne na każdej maszynie.
func krojDesignu(nazwa string) (*sfnt.Font, string, bool, error) {
	klucz := kluczKrojuDesignu(nazwa)
	if klucz == "" {
		return nil, "", false, fmt.Errorf("nazwa kroju jest pusta")
	}
	for nazwaKroju, bajty := range krojeWkompilowaneDesignu {
		if kluczKrojuDesignu(nazwaKroju) != klucz {
			continue
		}
		krojWczytany, err := sfnt.Parse(bajty)
		if err != nil {
			return nil, "", false, fmt.Errorf(
				"krój %s jest wkompilowany, ale jego zapis nie daje się rozłożyć: %w", nazwaKroju, err)
		}
		return krojWczytany, nazwaKroju, true, nil
	}
	for nazwaKroju, sciezka := range przegladajKrojeSerweraDesignu() {
		if kluczKrojuDesignu(nazwaKroju) != klucz {
			continue
		}
		bajty, err := os.ReadFile(sciezka)
		if err != nil {
			return nil, "", false, fmt.Errorf("kroju %s nie da się odczytać z tej maszyny: %w",
				nazwaKroju, err)
		}
		krojWczytany, err := sfnt.Parse(bajty)
		if err != nil {
			return nil, "", false, fmt.Errorf(
				"plik kroju %s leży na tej maszynie, ale nie jest krojem, który serwer rozkłada: %w",
				nazwaKroju, err)
		}
		return krojWczytany, nazwaKroju, false, nil
	}
	return nil, "", false, fmt.Errorf(
		"kroju %q serwer nie ma ani wkompilowanego, ani na tej maszynie — podglądu nie złoży krojem "+
			"zastępczym, bo wyglądałby jak prawdziwy; kroje wkompilowane: %s",
		nazwa, strings.Join(nazwyKrojowWkompilowanychDesignu(), ", "))
}

// nazwyKrojowWkompilowanychDesignu oddaje nazwy krojów z binarium
// w kolejności alfabetycznej, do treści odmowy.
func nazwyKrojowWkompilowanychDesignu() []string {
	nazwy := make([]string, 0, len(krojeWkompilowaneDesignu))
	for nazwa := range krojeWkompilowaneDesignu {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)
	return nazwy
}

// nazwyKrojowDesignu oddaje komplet krojów, jakimi rdzeń dysponuje:
// najpierw wkompilowane, potem odnalezione na maszynie serwera.
func nazwyKrojowDesignu() []string {
	nazwy := nazwyKrojowWkompilowanychDesignu()
	serwera := make([]string, 0, len(przegladajKrojeSerweraDesignu()))
	for nazwa := range przegladajKrojeSerweraDesignu() {
		serwera = append(serwera, nazwa)
	}
	sort.Strings(serwera)
	return append(nazwy, serwera...)
}

// sciezkaTekstuDesignu składa kontury tekstu jako ścieżkę biblioteki,
// ułożoną od punktu x, y linii pisma, w jednostkach ekranu.
func sciezkaTekstuDesignu(krojWczytany *sfnt.Font, tekst string, rozmiar, x, y float64) (*canvas.Path, error) {
	if strings.TrimSpace(tekst) == "" {
		return nil, fmt.Errorf("tekst jest pusty — nie ma czego zamienić w kontury")
	}
	if rozmiar <= 0 {
		return nil, fmt.Errorf("rozmiar pisma %v: pismo o niedodatnim rozmiarze nie ma konturu", rozmiar)
	}
	ppem := fixed.Int26_6(rozmiar * 64)
	bufor := &sfnt.Buffer{}
	sciezka := &canvas.Path{}
	pioro := x
	for _, znak := range tekst {
		numer, err := krojWczytany.GlyphIndex(bufor, znak)
		if err != nil {
			return nil, fmt.Errorf("krój nie oddał numeru glifu dla znaku %q: %w", znak, err)
		}
		if numer == 0 {
			// Glif zerowy znaczy krój tego znaku nie ma; odmowa jest tu
			// właściwa, nie pusty prostokąt.
			return nil, fmt.Errorf("krój nie ma glifu dla znaku %q — konturów nie da się złożyć "+
				"bez tej litery, a prostokąt zastępczy nie jest literą", znak)
		}
		odcinki, err := krojWczytany.LoadGlyph(bufor, numer, ppem, nil)
		if err != nil {
			return nil, fmt.Errorf("krój nie oddał konturu glifu znaku %q: %w", znak, err)
		}
		dopiszGlifDoSciezkiDesignu(sciezka, odcinki, pioro, y)
		postep, err := krojWczytany.GlyphAdvance(bufor, numer, ppem, 0)
		if err != nil {
			return nil, fmt.Errorf("krój nie oddał szerokości glifu znaku %q: %w", znak, err)
		}
		pioro += float64(postep) / 64
	}
	if sciezka.Empty() {
		return nil, fmt.Errorf("kontury tekstu wyszły puste — tekst złożony z samych odstępów " +
			"nie ma konturu")
	}
	return sciezka, nil
}

// dopiszGlifDoSciezkiDesignu przenosi odcinki jednego glifu na ścieżkę,
// przesunięte do bieżącego położenia pióra rysunku.
func dopiszGlifDoSciezkiDesignu(sciezka *canvas.Path, odcinki sfnt.Segments, pioro, linia float64) {
	punkt := func(p fixed.Point26_6) (float64, float64) {
		return pioro + float64(p.X)/64, linia + float64(p.Y)/64
	}
	for _, odcinek := range odcinki {
		switch odcinek.Op {
		case sfnt.SegmentOpMoveTo:
			x, y := punkt(odcinek.Args[0])
			sciezka.MoveTo(x, y)
		case sfnt.SegmentOpLineTo:
			x, y := punkt(odcinek.Args[0])
			sciezka.LineTo(x, y)
		case sfnt.SegmentOpQuadTo:
			x1, y1 := punkt(odcinek.Args[0])
			x2, y2 := punkt(odcinek.Args[1])
			sciezka.QuadTo(x1, y1, x2, y2)
		case sfnt.SegmentOpCubeTo:
			x1, y1 := punkt(odcinek.Args[0])
			x2, y2 := punkt(odcinek.Args[1])
			x3, y3 := punkt(odcinek.Args[2])
			sciezka.CubeTo(x1, y1, x2, y2, x3, y3)
		}
	}
	// Kontur glifu jest zamknięty z natury: litera jest obszarem, nie kreską.
	if !sciezka.Empty() && !sciezka.Closed() {
		sciezka.Close()
	}
}

// szerokoscTekstuDesignu mierzy postęp pióra dla całego tekstu; znak,
// którego krój nie ma, jest pomijany w pomiarze.
func szerokoscTekstuDesignu(krojWczytany *sfnt.Font, tekst string, rozmiar float64) (float64, int) {
	ppem := fixed.Int26_6(rozmiar * 64)
	bufor := &sfnt.Buffer{}
	szerokosc := 0.0
	brakujacych := 0
	for _, znak := range tekst {
		numer, err := krojWczytany.GlyphIndex(bufor, znak)
		if err != nil || numer == 0 {
			brakujacych++
			continue
		}
		postep, err := krojWczytany.GlyphAdvance(bufor, numer, ppem, 0)
		if err != nil {
			brakujacych++
			continue
		}
		szerokosc += float64(postep) / 64
	}
	return szerokosc, brakujacych
}

// glifyKrojuDesignu oddaje glify kroju z zakresu punktów kodowych wraz
// z ich konturem SVG — obsługuje design.font.glyphs.get.
func glifyKrojuDesignu(krojWczytany *sfnt.Font, od, do, granica int) []shared.DesignGlyph {
	ppem := fixed.Int26_6(krojWczytany.UnitsPerEm())
	bufor := &sfnt.Buffer{}
	glify := make([]shared.DesignGlyph, 0, granica)
	for punkt := od; punkt <= do && len(glify) < granica; punkt++ {
		numer, err := krojWczytany.GlyphIndex(bufor, rune(punkt))
		if err != nil || numer == 0 {
			continue
		}
		glif := shared.DesignGlyph{Codepoint: punkt}
		if nazwa, err := krojWczytany.GlyphName(bufor, numer); err == nil && nazwa != "" {
			wartosc := nazwa
			glif.Name = &wartosc
		}
		if postep, err := krojWczytany.GlyphAdvance(bufor, numer, ppem, 0); err == nil {
			wartosc := float64(postep) / 64
			glif.Advance = &wartosc
		}
		if odcinki, err := krojWczytany.LoadGlyph(bufor, numer, ppem, nil); err == nil {
			kontur := &canvas.Path{}
			dopiszGlifDoSciezkiDesignu(kontur, odcinki, 0, 0)
			if !kontur.Empty() {
				zapis := kontur.ToSVG()
				glif.Svg = &zapis
			}
		}
		glify = append(glify, glif)
	}
	return glify
}

// liczbaGlifowKrojuDesignu oddaje liczbę glifów kroju — pomiar, nie
// oszacowanie, policzony przejściem po tablicy cmap.
func liczbaGlifowKrojuDesignu(krojWczytany *sfnt.Font) int {
	return krojWczytany.NumGlyphs()
}
