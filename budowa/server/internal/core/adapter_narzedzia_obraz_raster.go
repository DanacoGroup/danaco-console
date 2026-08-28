// Plik podaje wspólne czytanie i liczenie na rastrze dla czynności rodziny image, które pracują biblioteką wkompilowaną w rdzeń.
package core

import (
	"image"
	"image/color"
	"os"

	// Dekodery rejestrowane importem pobocznym — potrzebne `image.Decode`.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"
)

// odczytajObrazPliku dekoduje plik do obrazu. Odmowa wraca surowa: nazywa ją
// wołający, bo to on wie, czym ten plik jest w jego czynności.
func odczytajObrazPliku(sciezka string) (image.Image, error) {
	plik, err := os.Open(sciezka)
	if err != nil {
		return nil, err
	}
	defer plik.Close()

	obraz, _, err := image.Decode(plik)
	if err != nil {
		return nil, err
	}
	return obraz, nil
}

// skladoweNiePrzemnozone rozkłada kolor na cztery składowe z zakresu 0..1, zdejmując wstępne przemnożenie przez alfę koloru.
func skladoweNiePrzemnozone(kolor color.Color) (r, g, b, a float64) {
	czerwona, zielona, niebieska, alfa := kolor.RGBA()
	if alfa == 0 {
		return 0, 0, 0, 0
	}
	skala := float64(alfa)
	return float64(czerwona) / skala, float64(zielona) / skala, float64(niebieska) / skala,
		float64(alfa) / 65535
}

// kolorZeSkladowych składa kolor z czterech składowych 0..1 z powrotem do
// postaci ośmiobitowej z alfą niewstępnie przemnożoną.
func kolorZeSkladowych(r, g, b, a float64) color.NRGBA {
	return color.NRGBA{
		R: bajtSkladowej(r), G: bajtSkladowej(g), B: bajtSkladowej(b), A: bajtSkladowej(a),
	}
}

// bajtSkladowej przycina składową do zakresu i zamienia na bajt. Przycięcie
// jest konieczne, bo `screen` i `overlay` potrafią wyjść nieznacznie poza
// jedynkę na zaokrągleniach.
func bajtSkladowej(wartosc float64) uint8 {
	switch {
	case wartosc <= 0:
		return 0
	case wartosc >= 1:
		return 255
	default:
		return uint8(wartosc*255 + 0.5)
	}
}
