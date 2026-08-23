// Odpowiedzialność pliku: kroje pisma modułu Design — katalog krojów, jakimi
// rdzeń NAPRAWDĘ dysponuje, wczytanie kroju, glify, szerokość tekstu i zamiana
// tekstu na kontury. Czytają to `design.font.preview`,
// `design.font.glyphs.get`, `design.font.pair.suggest`,
// `design.vector.text.path` i wydania drukarskie.
//
// ── Katalog jest POMIAREM, nie zapowiedzią ──────────────────────────────────
// Kroje są dwojakie i różnią się tym, czego Operator może być pewien:
//
//   - kroje WKOMPILOWANE w binarium (rodzina Go z `golang.org/x/image/font/gofont`)
//     — są zawsze, na każdej maszynie, bo leżą w pliku wykonywalnym serwera;
//   - kroje SERWERA — pliki `.ttf`/`.otf` leżące w katalogach krojów tej
//     maszyny; odczytywane, nie doinstalowywane.
//
// Pole `available` kontraktu mówi prawdę o tym rozróżnieniu. Kroju, którego nie
// ma, rdzeń nie podstawia innym: podgląd złożony krojem zastępczym wygląda
// identycznie jak podgląd prawdziwy i Operator wybrałby typografię, której
// u siebie nie zobaczy.
//
// ── Rachunek jest wkompilowany, nie wołany ──────────────────────────────────
// Kontury glifów czyta `golang.org/x/image/font/sfnt` — biblioteka Go
// wkompilowana w binarium, rozkładająca zarówno kontury TrueType, jak
// i PostScript (CFF). Programy do przeglądania i rozkładania krojów tu nie
// wchodzą: instalka Operatora ich nie niesie.
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
	// systemach, na jakich stoi serwer produktu. Odczyt jest wyłącznie odczytem
	// katalogu — rdzeń krojów nie instaluje i nie kopiuje.
	katalogKrojowSystemowy = "/usr/share/fonts"
	katalogKrojowLokalny   = "/usr/local/share/fonts"

	// granicaGlifowOdpowiedziDesignu chroni odpowiedź przed wykazem
	// dwudziestotysięcznym: kroje pełne CJK mają tyle glifów, a każdy niesie
	// kontur SVG.
	granicaGlifowOdpowiedziDesignu = 512

	// domyslnyTekstProbnyDesignu jest tekstem podglądu, gdy Operator własnego
	// nie podał. Zdanie polskie z ogonkami, bo produkt jest polski i to na
	// ogonkach poznaje się, czy krój je w ogóle ma.
	domyslnyTekstProbnyDesignu = "Zażółć gęślą jaźń — ĄĆĘŁŃÓŚŹŻ 0123456789"
)

// krojeWkompilowaneDesignu to kroje leżące w binarium serwera. Rodzina Go jest
// tu w całości, bo to jedyna rodzina, którą produkt ma NA PEWNO — u Operatora
// bez ani jednego kroju w systemie podgląd i wydanie nadal działają.
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

// pamiecKrojowSerweraDesignu trzyma wynik przeglądu katalogów krojów. Przegląd
// idzie raz na proces: katalog krojów systemu nie zmienia się w trakcie pracy
// serwera, a odczyt stu plików przy każdym podglądzie byłby kosztem bez zysku.
var (
	razKrojowSerweraDesignu sync.Once
	krojeSerweraDesignu     map[string]string
)

// przegladajKrojeSerweraDesignu składa wykaz krojów tej maszyny: nazwa rodziny
// → ścieżka pliku. Nazwa bierze się z pliku, nie z tablicy `name` kroju —
// otwarcie stu plików tylko po nazwę kosztowałoby tyle, co złożenie stu
// podglądów. Nazwa z tablicy wchodzi dopiero przy wczytaniu jednego kroju.
//
// Katalog nieistniejący nie jest awarią: na maszynie bez pakietu krojów go po
// prostu nie ma, a kroje wkompilowane zostają.
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

// nazwaKrojuZPlikuDesignu wyciąga nazwę rodziny z nazwy pliku:
// `NimbusSans-Regular.otf` → `NimbusSans-Regular`.
func nazwaKrojuZPlikuDesignu(plik string) string {
	return strings.TrimSuffix(plik, filepath.Ext(plik))
}

// kluczKrojuDesignu sprowadza nazwę kroju do postaci porównywalnej: bez spacji,
// dywizów i podkreśleń, małymi literami. Operator wpisuje „Go Mono", „go-mono"
// i „GoMono" mając na myśli ten sam krój, a odmowa za znak rozdzielający byłaby
// odmową za zapis poprawny.
func kluczKrojuDesignu(nazwa string) string {
	zamiana := strings.NewReplacer(" ", "", "-", "", "_", "", ".", "")
	return strings.ToLower(zamiana.Replace(strings.TrimSpace(nazwa)))
}

// krojDesignu wczytuje krój po nazwie i mówi, skąd pochodzi.
//
// Pierwszeństwo mają kroje wkompilowane: są na każdej maszynie, więc wydanie
// nimi złożone wygląda tak samo u Operatora i na serwerze. Krój serwera wchodzi
// dopiero wtedy, gdy nazwa nie trafia w żaden wkompilowany.
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
				"plik kroju %s leży na tej maszynie, ale nie jest krojem, który rdzeń rozkłada: %w",
				nazwaKroju, err)
		}
		return krojWczytany, nazwaKroju, false, nil
	}
	return nil, "", false, fmt.Errorf(
		"kroju %q rdzeń nie ma ani wkompilowanego, ani na tej maszynie — podglądu nie złoży krojem "+
			"zastępczym, bo wyglądałby jak prawdziwy; kroje wkompilowane: %s",
		nazwa, strings.Join(nazwyKrojowWkompilowanychDesignu(), ", "))
}

// nazwyKrojowWkompilowanychDesignu oddaje nazwy krojów z binarium w kolejności
// alfabetycznej — wykaz wchodzi w treść odmowy, żeby Operator miał dokąd pójść.
func nazwyKrojowWkompilowanychDesignu() []string {
	nazwy := make([]string, 0, len(krojeWkompilowaneDesignu))
	for nazwa := range krojeWkompilowaneDesignu {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)
	return nazwy
}

// nazwyKrojowDesignu oddaje komplet krojów, jakimi rdzeń dysponuje: najpierw
// wkompilowane, potem odnalezione na maszynie.
func nazwyKrojowDesignu() []string {
	nazwy := nazwyKrojowWkompilowanychDesignu()
	serwera := make([]string, 0, len(przegladajKrojeSerweraDesignu()))
	for nazwa := range przegladajKrojeSerweraDesignu() {
		serwera = append(serwera, nazwa)
	}
	sort.Strings(serwera)
	return append(nazwy, serwera...)
}

// sciezkaTekstuDesignu składa kontury tekstu jako ścieżkę biblioteki, ułożoną
// od punktu (x, y) linii pisma.
//
// Glify jadą przez `sfnt`, w jednostkach ekranu przy zadanym rozmiarze pisma.
// Oś Y rośnie w dół — tak samo jak w kompozycji Design Board — więc kontur nie
// jest tu nigdzie odwracany i tekst nie wychodzi do góry nogami.
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
			// Glif zerowy znaczy „krój tego znaku nie ma". Odmowa jest tu
			// właściwa: kontury z pustym prostokątem w miejscu litery byłyby
			// napisem, którego Operator nie zamawiał.
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
// przesunięte do bieżącego położenia pióra.
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

// szerokoscTekstuDesignu mierzy postęp pióra dla całego tekstu. Znak, którego
// krój nie ma, jest pomijany w pomiarze i wraca liczbą pominiętych — pomiar ma
// się udać nawet dla tekstu z jednym znakiem spoza kroju, bo służy do ułożenia
// podglądu, nie do wydania.
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

// glifyKrojuDesignu oddaje glify kroju z zakresu punktów kodowych wraz z ich
// konturem SVG — obsługuje rachunek `design.font.glyphs.get`.
//
// Kontur wchodzi do odpowiedzi, bo bez niego wykaz glifów byłby wykazem liczb.
// Znak, którego krój nie ma, nie wchodzi do wykazu w ogóle: kontrakt pyta
// o glify KROJU, a nie o zakres, którego Operator się spodziewał.
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

// liczbaGlifowKrojuDesignu oddaje liczbę glifów kroju — pomiar, nie oszacowanie.
func liczbaGlifowKrojuDesignu(krojWczytany *sfnt.Font) int {
	return krojWczytany.NumGlyphs()
}
