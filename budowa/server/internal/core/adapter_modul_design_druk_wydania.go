// Odpowiedzialność pliku: kodery wydań drukarskich, których biblioteka
// standardowa nie ma — TIFF i EPS — oraz barwa znacznika sklejenia kafli.
// Czynności części drukarskiej leżą w `adapter_modul_design_druk.go`.
//
// ── TIFF przez bibliotekę wkompilowaną ──────────────────────────────────────
// Zapis TIFF idzie przez `golang.org/x/image/tiff` — tę samą bibliotekę, którą
// rdzeń TIFF CZYTA. Kompresja jest bezstratna (Deflate): materiał drukarski
// przepuszczony przez kompresję stratną wraca z drukarni z widocznymi artefaktami
// na płaskich plamach, a to jest dokładnie ten materiał, dla którego istnieje
// TIFF.
//
// ── EPS rdzeń pisze sam, i mówi dlaczego ────────────────────────────────────
// Kodera EPS z obrazem rastrowym nie ma ani w bibliotece standardowej, ani
// w `tdewolff/canvas` (jej wydanie PostScript rysuje ŚCIEŻKI, nie osadza
// pikseli). Programy do rasteryzacji i przekształceń obrazu, którymi zwykle się
// to robi, leżą poza instalką Operatora — nie wolno ich nazwać nawet
// w komentarzu, żeby nikt nie wziął nazwy za wskazówkę, i nie wolno od nich
// zależeć, bo u Operatora byłyby odmową. Zapis stoi więc tutaj i jest
// wkompilowany.
//
// EPS niesie obraz operatorem `colorimage` z danymi szesnastkowymi. Zapis
// szesnastkowy jest dwa razy dłuższy od binarnego, ale jest CZYSTYM tekstem —
// przechodzi przez każdy strumień, każdą bramkę pocztową i każdy system
// drukarski, także taki, który psuje bajty ósmego bitu. Dla materiału, który
// jedzie do obcej drukarni, ta pewność jest warta dwukrotności rozmiaru.
package core

import (
	"bytes"
	"fmt"
	"image"
	"image/color"

	"golang.org/x/image/tiff"
)

// barwaZnacznikaSklejeniaDesignu oddaje barwę linii zakładki na kaflu.
//
// Magenta, nie czerń: czerń zlewa się z treścią na większości materiałów, a
// magenta w druku wielkoformatowym jest barwą, której na materiałach użytkowych
// prawie nie ma — więc znacznik jest widoczny i nie da się go pomylić z rysunkiem.
func barwaZnacznikaSklejeniaDesignu() color.RGBA {
	return color.RGBA{R: 255, G: 0, B: 255, A: 255}
}

// zakodujTiffDesignu zapisuje obraz jako TIFF z kompresją bezstratną.
func zakodujTiffDesignu(obraz image.Image) ([]byte, error) {
	var bufor bytes.Buffer
	if err := tiff.Encode(&bufor, obraz, &tiff.Options{
		Compression: tiff.Deflate, Predictor: true,
	}); err != nil {
		return nil, fmt.Errorf("nie można zapisać wydania jako tiff: %w", err)
	}
	return bufor.Bytes(), nil
}

// zakodujEpsDesignu zapisuje obraz jako dokument EPS z osadzonym obrazem
// rastrowym.
//
// Rozmiar strony liczy się w punktach typograficznych z wymiaru pikselowego przy
// rozdzielczości drukarskiej: EPS nie ma pola na rozdzielczość, więc jedyną
// drogą podania fizycznego rozmiaru jest pole `BoundingBox`.
func zakodujEpsDesignu(obraz image.Image) ([]byte, error) {
	granice := obraz.Bounds()
	szerokosc, wysokosc := granice.Dx(), granice.Dy()
	if szerokosc < 1 || wysokosc < 1 {
		return nil, fmt.Errorf("obraz o boku %d×%d nie istnieje", szerokosc, wysokosc)
	}
	szerokoscPunktow := float64(szerokosc) / float64(rozdzielczoscDrukuDomyslna) * punktyNaCal
	wysokoscPunktow := float64(wysokosc) / float64(rozdzielczoscDrukuDomyslna) * punktyNaCal

	var dokument bytes.Buffer
	dokument.WriteString("%!PS-Adobe-3.0 EPSF-3.0\n")
	fmt.Fprintf(&dokument, "%%%%BoundingBox: 0 0 %d %d\n",
		int(szerokoscPunktow+0.5), int(wysokoscPunktow+0.5))
	fmt.Fprintf(&dokument, "%%%%HiResBoundingBox: 0 0 %.4f %.4f\n",
		szerokoscPunktow, wysokoscPunktow)
	dokument.WriteString("%%Creator: Danaco Console, moduł Design\n")
	dokument.WriteString("%%LanguageLevel: 2\n")
	dokument.WriteString("%%EndComments\n")
	dokument.WriteString("gsave\n")
	// Układ PostScriptu liczy Y od dołu, a obraz od góry: skala ujemna w Y
	// i przesunięcie na wysokość strony odwracają go raz, w jednym miejscu.
	fmt.Fprintf(&dokument, "0 %.4f translate\n", wysokoscPunktow)
	fmt.Fprintf(&dokument, "%.4f %.4f scale\n", szerokoscPunktow, -wysokoscPunktow)
	dokument.WriteString("/DeviceRGB setcolorspace\n")
	fmt.Fprintf(&dokument, "%d %d 8 [%d 0 0 %d 0 0]\n",
		szerokosc, wysokosc, szerokosc, wysokosc)
	dokument.WriteString("{currentfile 3 string readhexstring pop} bind\n")
	dokument.WriteString("false 3 colorimage\n")

	// Dane obrazu: trzy składowe na piksel, zapis szesnastkowy, wiersze łamane
	// co 32 bajty. Łamanie jest wymagane — PostScript nie gwarantuje obsługi
	// wiersza dłuższego niż 255 znaków.
	const bajtowWWierszu = 32
	licznik := 0
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			r, g, b, alfa := obraz.At(x, y).RGBA()
			// Przezroczystość w EPS nie istnieje: piksel częściowo przezroczysty
			// zlewa się z BIELĄ, bo takie jest tło nośnika. Zostawienie samych
			// składowych dałoby na wydruku ciemną obwódkę wokół każdej krawędzi.
			if alfa < 65535 {
				r, g, b = zlejZBielaDesignu(r, alfa), zlejZBielaDesignu(g, alfa),
					zlejZBielaDesignu(b, alfa)
			}
			fmt.Fprintf(&dokument, "%02x%02x%02x", r>>8, g>>8, b>>8)
			licznik += 3
			if licznik >= bajtowWWierszu {
				dokument.WriteString("\n")
				licznik = 0
			}
		}
	}
	if licznik > 0 {
		dokument.WriteString("\n")
	}
	dokument.WriteString("grestore\n")
	dokument.WriteString("%%EOF\n")
	return dokument.Bytes(), nil
}

// zlejZBielaDesignu zlewa składową z krycim wmnożonym z białym tłem.
//
// Składowe `color.RGBA` są PRZEMNOŻONE przez krycie, więc dołożenie białego tła
// to dodanie tego, czego kryciu brakuje do pełni.
func zlejZBielaDesignu(skladowa, alfa uint32) uint32 {
	return skladowa + (65535 - alfa)
}
