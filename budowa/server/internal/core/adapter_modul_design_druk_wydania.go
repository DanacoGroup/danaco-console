// Odpowiedzialność pliku: kodery wydań drukarskich, których biblioteka
// standardowa nie ma — TIFF i EPS — oraz barwa znacznika sklejenia kafli.
// Czynności części drukarskiej leżą w `adapter_modul_design_druk.go`.
package core

import (
	"bytes"
	"fmt"
	"image"
	"image/color"

	"golang.org/x/image/tiff"
)

// barwaZnacznikaSklejeniaDesignu oddaje barwę linii zakładki na kaflu.
// Magenta, nie czerń: czerń zlewa się z treścią na większości materiałów.
func barwaZnacznikaSklejeniaDesignu() color.RGBA {
	return color.RGBA{R: 255, G: 0, B: 255, A: 255}
}

// zakodujTiffDesignu zapisuje obraz jako TIFF z kompresją bezstratną
// Deflate, bez utraty jakości druku.
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
// rastrowym. Rozmiar strony liczy się w punktach typograficznych z wymiaru
// pikselowego przy rozdzielczości drukarskiej.
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
	// Układ PostScriptu liczy Y od dołu, a obraz od góry.
	fmt.Fprintf(&dokument, "0 %.4f translate\n", wysokoscPunktow)
	fmt.Fprintf(&dokument, "%.4f %.4f scale\n", szerokoscPunktow, -wysokoscPunktow)
	dokument.WriteString("/DeviceRGB setcolorspace\n")
	fmt.Fprintf(&dokument, "%d %d 8 [%d 0 0 %d 0 0]\n",
		szerokosc, wysokosc, szerokosc, wysokosc)
	dokument.WriteString("{currentfile 3 string readhexstring pop} bind\n")
	dokument.WriteString("false 3 colorimage\n")

	// Dane obrazu: wiersze łamane co 32 bajty, bo PostScript ma limit długości.
	const bajtowWWierszu = 32
	licznik := 0
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		for x := granice.Min.X; x < granice.Max.X; x++ {
			r, g, b, alfa := obraz.At(x, y).RGBA()
			// Przezroczystość w EPS nie istnieje: piksel zlewa się z bielą tła.
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
