// Moduł Library — odciski treści i obrazu, na których stoi rozpoznanie
// niemal-duplikatów.
//
// Dwa odciski, bo to dwa różne pytania. Dla tekstu podobieństwo znaczy „te same
// słowa w podobnych proporcjach" — liczy je miara Jaccarda na zbiorach słów.
// Dla obrazu podobieństwo znaczy „ten sam obraz mimo innej kompresji i rozmiaru"
// — liczy je odcisk percepcyjny: obraz sprowadzony do siatki 8×8 w skali szarości
// i zamieniony na 64 bity względem jasności średniej. Dwa odciski różniące się
// o kilka bitów pochodzą z tego samego zdjęcia, choćby plik ważył dziesięć razy
// mniej.
//
// Oba liczy kod wkompilowany: dekodery obrazu ze standardowej biblioteki Go
// oraz `golang.org/x/image` na WebP. Ani jednego uruchomienia programu z zewnątrz
// — rozpoznanie duplikatów ma działać u Operatora, a nie na maszynie, na której
// ktoś doinstalował ImageMagicka.
package core

import (
	"bytes"
	"image"
	"math/bits"
	"strings"

	// Rejestracja dekoderów obrazu: bez tych importów `image.Decode` zna sam
	// format, w którym akurat zapisano plik testowy.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"danacoconsole/server/internal/dane"
)

// bokSiatkiOdcisku — bok siatki, do której sprowadza się obraz. Osiem na osiem
// daje 64 bity: dość, żeby odróżnić dwa różne zdjęcia, i mało, żeby ta sama
// fotografia po przeskalowaniu dała ten sam odcisk.
const bokSiatkiOdcisku = 8

// odciskTekstuBiblioteki składa zbiór słów treści.
//
// Zbiór, nie ciąg: przestawienie akapitów nie czyni z dokumentu innego
// dokumentu, a kolejność słów rozstrzygałaby przeciwnie. Słowa krótsze niż trzy
// znaki odpadają — spójniki są w każdym tekście i podnosiłyby podobieństwo
// dowolnych dwóch.
func odciskTekstuBiblioteki(bajty []byte) map[string]struct{} {
	tekst := strings.ToLower(string(bajty))
	slowa := strings.FieldsFunc(tekst, func(znak rune) bool {
		return !(znak >= 'a' && znak <= 'z' || znak >= '0' && znak <= '9' || znak > 127)
	})
	zbior := map[string]struct{}{}
	for _, slowo := range slowa {
		if len([]rune(slowo)) < 3 {
			continue
		}
		zbior[slowo] = struct{}{}
	}
	return zbior
}

// podobienstwoZbiorowSlowBiblioteki liczy miarę Jaccarda w setnych: część wspólna do sumy.
func podobienstwoZbiorowSlowBiblioteki(pierwszy, drugi map[string]struct{}) int {
	if len(pierwszy) == 0 || len(drugi) == 0 {
		return 0
	}
	wspolne := 0
	mniejszy, wiekszy := pierwszy, drugi
	if len(wiekszy) < len(mniejszy) {
		mniejszy, wiekszy = wiekszy, mniejszy
	}
	for slowo := range mniejszy {
		if _, jest := wiekszy[slowo]; jest {
			wspolne++
		}
	}
	suma := len(pierwszy) + len(drugi) - wspolne
	if suma == 0 {
		return 0
	}
	return wspolne * 100 / suma
}

// odciskObrazuBiblioteki liczy odcisk percepcyjny obrazu.
func odciskObrazuBiblioteki(bajty []byte) (uint64, error) {
	obraz, _, err := image.Decode(bytes.NewReader(bajty))
	if err != nil {
		return 0, err
	}
	granice := obraz.Bounds()
	if granice.Dx() == 0 || granice.Dy() == 0 {
		return 0, errPustyObrazBiblioteki
	}

	// Próbkowanie punktowe zamiast uśredniania obszarów: odcisk ma być tani,
	// a różnica w trafności między jednym a drugim jest mniejsza niż próg, od
	// którego uznajemy obrazy za podobne.
	jasnosci := make([]float64, 0, bokSiatkiOdcisku*bokSiatkiOdcisku)
	suma := 0.0
	for wiersz := 0; wiersz < bokSiatkiOdcisku; wiersz++ {
		for kolumna := 0; kolumna < bokSiatkiOdcisku; kolumna++ {
			x := granice.Min.X + granice.Dx()*kolumna/bokSiatkiOdcisku
			y := granice.Min.Y + granice.Dy()*wiersz/bokSiatkiOdcisku
			czerwony, zielony, niebieski, _ := obraz.At(x, y).RGBA()
			// Wagi luminancji: oko widzi zieleń jaśniej niż błękit, więc odcisk
			// liczony na średniej arytmetycznej rozjeżdżałby się przy zmianie
			// nasycenia barw.
			jasnosc := 0.299*float64(czerwony) + 0.587*float64(zielony) + 0.114*float64(niebieski)
			jasnosci = append(jasnosci, jasnosc)
			suma += jasnosc
		}
	}
	srednia := suma / float64(len(jasnosci))

	var odcisk uint64
	for indeks, jasnosc := range jasnosci {
		if jasnosc >= srednia {
			odcisk |= 1 << uint(indeks)
		}
	}
	return odcisk, nil
}

// podobienstwoOdciskowObrazuBiblioteki liczy trafność w setnych na podstawie liczby bitów
// zgodnych — odległość Hamminga odwrócona na skalę procentową.
func podobienstwoOdciskowObrazuBiblioteki(pierwszy, drugi uint64) int {
	rozne := bits.OnesCount64(pierwszy ^ drugi)
	zgodne := 64 - rozne
	return zgodne * 100 / 64
}

// rodzajPodgladuJestObrazem mówi, czy zasób jest materiałem graficznym —
// rozpoznanie po obrazie ma sens wyłącznie dla obrazów, a rozpoznanie po tekście
// wyłącznie dla reszty.
func rodzajPodgladuJestObrazem(zasob dane.PlikBiblioteki) bool {
	if zasob.MimeType != nil && strings.HasPrefix(strings.ToLower(*zasob.MimeType), "image/") {
		return true
	}
	if zasob.MimeType != nil && *zasob.MimeType != "" {
		return false
	}
	// Bez rodzaju treści rozstrzyga rozszerzenie nazwy — ta sama droga, którą
	// idzie podgląd (`adapter_modul_library_podglad.go`).
	nazwa := strings.ToLower(zasob.Nazwa)
	for _, rozszerzenie := range []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".tif", ".tiff"} {
		if strings.HasSuffix(nazwa, rozszerzenie) {
			return true
		}
	}
	return false
}
