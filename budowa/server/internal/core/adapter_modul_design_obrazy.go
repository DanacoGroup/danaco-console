// Odpowiedzialność pliku: warsztat obrazu modułu Design — odczyt bajtów
// zasobu, skalowanie i przekodowanie do formatu wydania. Z niego korzystają
// `design.asset.export`, `design.asset.export.batch` i wyrys kompozycji
// (`adapter_modul_design_wyrys.go`).
//
// ── Biblioteka wkompilowana, nigdy program zewnętrzny ───────────────────────
// Nie ma tu ani jednego uruchomienia procesu i mieć nie będzie. Skalowanie robi
// `golang.org/x/image/draw` filtrem `CatmullRom`, zapis PNG i JPEG — biblioteka
// standardowa, dokument — `pdfcpu`, ikonę — koder w tym pliku. Powód jest ten
// sam, co w warsztacie PDF modułu Studio (`adapter_studio_pdf.go`): funkcja
// zależna od programu, którego instalka nie niesie, jest u Operatora odmową,
// a nie funkcją, a sprawdzian na maszynie deweloperskiej świeciłby przy niej
// zielono.
//
// ── Formaty, których biblioteka nie umie, są ODMAWIANE ──────────────────────
// WEBP wchodzi (dekoder `golang.org/x/image/webp`), ale nie wychodzi: enkodera
// WEBP w Go bez zależności zewnętrznej nie ma. AVIF nie wchodzi i nie wychodzi.
// Wydanie takiego formatu kończy się odmową wymieniającą formaty obsługiwane —
// cichy PNG pod nazwą `.avif` byłby plikiem, który Operator wyśle dalej jako
// AVIF i który odbiorcy nie otworzy się tam, gdzie miał się otworzyć.
//
// ── SVG przechodzi bez rasteryzacji ─────────────────────────────────────────
// Zasób wektorowy wydany jako SVG to te same bajty, które leżą w magazynie:
// rasteryzacja odebrałaby mu jedyną własność, dla której jest wektorem. Skala
// nie ma wtedy zastosowania i nie jest po cichu stosowana.
package core

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	// Dekodery wejścia. Import pusty — używamy wyłącznie skutku rejestracji,
	// tak przewiduje pakiet `image`.
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	_ "image/gif"

	"golang.org/x/image/draw"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// formatyWydaniaDesignu wylicza formaty, w których rdzeń NAPRAWDĘ potrafi
	// wydać zasób. Wykaz wchodzi wprost w treść każdej odmowy formatu, żeby
	// Operator dostał drogę wyjścia, a nie samo „nie".
	formatyWydaniaDesignu = "png, jpeg, svg, pdf, ico"

	// granicaSkaliWydaniaDesignu chroni rdzeń przed żądaniem, które zamawia
	// obraz nie do zmieszczenia w pamięci. Nie jest to polityka jakości, tylko
	// granica rachunku: skala tysiąckrotna z obrazu 4000×4000 to bilion pikseli.
	granicaSkaliWydaniaDesignu = 16.0
)

// obrazZasobuDesignu czyta bajty zasobu spod ścieżki magazynu i rozkłada je na
// obraz. Odmowa nazywa, czego zabrakło: pliku pod odwołaniem albo formatu,
// którego dekoder nie zna.
func obrazZasobuDesignu(sciezka string) (image.Image, error) {
	plik, err := os.Open(sciezka)
	if err != nil {
		return nil, fmt.Errorf("treści zasobu nie ma pod jego odwołaniem w magazynie: %w", err)
	}
	defer plik.Close()

	obraz, _, err := image.Decode(plik)
	if err != nil {
		return nil, fmt.Errorf("treść zasobu nie jest obrazem, który rdzeń potrafi rozłożyć "+
			"(rozkłada png, jpeg, gif, webp, bmp, tiff): %w", err)
	}
	return obraz, nil
}

// przeskalujObrazDesignu oddaje obraz w zadanej krotności skali. Skala pusta
// albo równa jedności oddaje obraz bez dotknięcia — przepuszczenie go przez
// filtr „dla porządku" kosztowałoby jakość bez żadnego zysku.
//
// Filtr jest `CatmullRom`: przy pomniejszaniu nie zostawia schodków, przy
// powiększaniu nie rozmywa krawędzi tak, jak dwuliniowy. Miniatura ikony
// wydana filtrem najbliższego sąsiada wygląda na uszkodzoną, a to jest
// najczęstsze wydanie tego modułu.
func przeskalujObrazDesignu(obraz image.Image, skala float64) image.Image {
	if skala <= 0 || skala == 1 {
		return obraz
	}
	granice := obraz.Bounds()
	szerokosc := int(float64(granice.Dx())*skala + 0.5)
	wysokosc := int(float64(granice.Dy())*skala + 0.5)
	if szerokosc < 1 {
		szerokosc = 1
	}
	if wysokosc < 1 {
		wysokosc = 1
	}
	plotno := image.NewRGBA(image.Rect(0, 0, szerokosc, wysokosc))
	draw.CatmullRom.Scale(plotno, plotno.Bounds(), obraz, granice, draw.Over, nil)
	return plotno
}

// zakodujObrazDesignu składa bajty wydania w żądanym formacie wraz z typem
// treści wedle IANA.
//
// Jakość dotyczy wyłącznie JPEG — PNG jest bezstratny, a udawanie, że suwak
// jakości nim rusza, byłoby pokrętłem podłączonym donikąd.
func zakodujObrazDesignu(obraz image.Image, format string, jakosc *int) ([]byte, string, error) {
	var bufor bytes.Buffer
	switch normalizujFormatWydaniaDesignu(format) {
	case "png":
		if err := png.Encode(&bufor, obraz); err != nil {
			return nil, "", fmt.Errorf("nie można zapisać wydania jako png: %w", err)
		}
		return bufor.Bytes(), "image/png", nil
	case "jpeg":
		nastawy := &jpeg.Options{Quality: jakoscWydaniaDesignu(jakosc)}
		if err := jpeg.Encode(&bufor, obraz, nastawy); err != nil {
			return nil, "", fmt.Errorf("nie można zapisać wydania jako jpeg: %w", err)
		}
		return bufor.Bytes(), "image/jpeg", nil
	case "ico":
		bajty, err := zakodujIkoneDesignu(obraz)
		if err != nil {
			return nil, "", err
		}
		return bajty, "image/vnd.microsoft.icon", nil
	case "pdf":
		bajty, err := zakodujDokumentDesignu(obraz)
		if err != nil {
			return nil, "", err
		}
		return bajty, "application/pdf", nil
	}
	return nil, "", fmt.Errorf("formatu %q rdzeń nie wydaje", format)
}

// jakoscWydaniaDesignu przycina jakość kompresji do zakresu kontraktu (1-100).
// Brak wskazania bierze wartość domyślną biblioteki — nie zero, które dałoby
// obraz nie do obejrzenia.
func jakoscWydaniaDesignu(jakosc *int) int {
	if jakosc == nil {
		return jpeg.DefaultQuality
	}
	if *jakosc < 1 {
		return 1
	}
	if *jakosc > 100 {
		return 100
	}
	return *jakosc
}

// normalizujFormatWydaniaDesignu sprowadza zapis formatu do jednej postaci:
// `JPG`, `jpg` i `jpeg` są tym samym formatem, a Operator wpisuje je zamiennie.
func normalizujFormatWydaniaDesignu(format string) string {
	nazwa := strings.ToLower(strings.TrimSpace(format))
	nazwa = strings.TrimPrefix(nazwa, ".")
	if nazwa == "jpg" {
		return "jpeg"
	}
	if nazwa == "tif" {
		return "tiff"
	}
	return nazwa
}

// rozszerzenieWydaniaDesignu oddaje rozszerzenie pliku dla formatu wydania.
// `jpeg` wychodzi jako `.jpg`, bo tego oczekuje system plików Operatora.
func rozszerzenieWydaniaDesignu(format string) string {
	if normalizujFormatWydaniaDesignu(format) == "jpeg" {
		return "jpg"
	}
	return normalizujFormatWydaniaDesignu(format)
}

// odmowaFormatuWydaniaDesignu składa zdanie odmowy formatu, którego rdzeń nie
// wydaje — wraz z wykazem formatów obsługiwanych, żeby Operator miał dokąd
// pójść. WEBP i AVIF dostają zdanie własne, bo ich brak ma powód, którego nie
// widać z wykazu: WEBP rdzeń CZYTA, więc „nie znam" byłoby nieprawdą.
func odmowaFormatuWydaniaDesignu(komenda, format string) error {
	nazwa := normalizujFormatWydaniaDesignu(format)
	switch nazwa {
	case "webp":
		return bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s z formatem wydania webp: rdzeń CZYTA webp jako materiał, ale go nie zapisuje — "+
				"zapis webp wymagałby programu spoza instalki, a produkt takich nie używa; "+
				"formaty wydania: %s", komenda, formatyWydaniaDesignu))
	case "avif":
		return bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s z formatem wydania avif: ani odczytu, ani zapisu avif rdzeń nie ma — "+
				"biblioteki wkompilowanej dla tego formatu nie ma, a wydanie png pod nazwą .avif "+
				"byłoby plikiem, który u odbiorcy się nie otworzy; formaty wydania: %s",
			komenda, formatyWydaniaDesignu))
	}
	return bladWskazaniaDesignu(fmt.Sprintf(
		"komenda %s z formatem wydania %q, którego rdzeń nie zapisuje; formaty wydania: %s",
		komenda, format, formatyWydaniaDesignu))
}

// sprawdzSkaleWydaniaDesignu odrzuca skalę bezsensowną PRZED wczytaniem obrazu.
// Skala niedodatnia dałaby obraz o zerowym boku, a skala nad granicą — żądanie
// pamięci, którego rdzeń nie zaspokoi; jedno i drugie jest pomyłką wołającego,
// nie awarią rdzenia.
func sprawdzSkaleWydaniaDesignu(komenda string, skala *float64) error {
	if skala == nil {
		return nil
	}
	if *skala <= 0 {
		return bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s ze skalą %v: skala niedodatnia dałaby obraz o zerowym boku", komenda, *skala))
	}
	if *skala > granicaSkaliWydaniaDesignu {
		return bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s ze skalą %v: granica skali wydania to %v — powyżej niej rdzeń zamawiałby "+
				"pamięć, której nie dostanie", komenda, *skala, granicaSkaliWydaniaDesignu))
	}
	return nil
}

// zakodujIkoneDesignu składa plik ICO z jednego obrazu.
//
// Zawartością wpisu jest PNG, nie mapa bitowa DIB: postać PNG w ikonie
// przyjmują wszystkie systemy od Windows Vista i wszystkie przeglądarki, a DIB
// wymagałby własnej maski przezroczystości zapisanej odwróconymi wierszami —
// czyli drugiego kodera obrazu w tym pliku.
//
// Bok ponad 256 pikseli jest odmawiany, a nie przycinany po cichu: pole
// szerokości w katalogu ikony ma JEDEN bajt, zero znaczy w nim 256, więc obraz
// większy zapisałby się jako ikona o boku wziętym z reszty z dzielenia.
func zakodujIkoneDesignu(obraz image.Image) ([]byte, error) {
	granice := obraz.Bounds()
	szerokosc, wysokosc := granice.Dx(), granice.Dy()
	if szerokosc > 256 || wysokosc > 256 {
		return nil, fmt.Errorf("ikona o boku %d×%d nie zmieści się w formacie ico "+
			"(największa ikona ma bok 256) — zmniejsz wydanie skalą", szerokosc, wysokosc)
	}
	if szerokosc < 1 || wysokosc < 1 {
		return nil, fmt.Errorf("ikona o boku %d×%d nie istnieje", szerokosc, wysokosc)
	}

	var zawartosc bytes.Buffer
	if err := png.Encode(&zawartosc, obraz); err != nil {
		return nil, fmt.Errorf("nie można zapisać zawartości ikony: %w", err)
	}

	var plik bytes.Buffer
	// Nagłówek: pole zastrzeżone (0), rodzaj (1 = ikona), liczba wpisów.
	_ = binary.Write(&plik, binary.LittleEndian, uint16(0))
	_ = binary.Write(&plik, binary.LittleEndian, uint16(1))
	_ = binary.Write(&plik, binary.LittleEndian, uint16(1))
	// Wpis katalogu: bok 256 zapisuje się zerem — tak stanowi format.
	plik.WriteByte(byte(szerokosc % 256))
	plik.WriteByte(byte(wysokosc % 256))
	plik.WriteByte(0)                                        // barwy palety: zero, bo obraz jest pełnobarwny
	plik.WriteByte(0)                                        // pole zastrzeżone
	_ = binary.Write(&plik, binary.LittleEndian, uint16(1))  // płaszczyzny barw
	_ = binary.Write(&plik, binary.LittleEndian, uint16(32)) // bity na piksel
	_ = binary.Write(&plik, binary.LittleEndian, uint32(zawartosc.Len()))
	// Przesunięcie zawartości: nagłówek (6) plus jeden wpis katalogu (16).
	_ = binary.Write(&plik, binary.LittleEndian, uint32(22))
	plik.Write(zawartosc.Bytes())
	return plik.Bytes(), nil
}

// zakodujDokumentDesignu osadza obraz w dokumencie PDF biblioteką `pdfcpu` —
// tą samą, którą pracuje warsztat dokumentu modułu Studio.
//
// Droga wiedzie przez PNG w pamięci, nie przez plik pośredni: `pdfcpu` przyjmuje
// strumień, a plik pośredni byłby trzecim miejscem, w którym ta sama treść żyje.
func zakodujDokumentDesignu(obraz image.Image) ([]byte, error) {
	var zrodlo bytes.Buffer
	if err := png.Encode(&zrodlo, obraz); err != nil {
		return nil, fmt.Errorf("nie można przygotować obrazu do osadzenia w dokumencie: %w", err)
	}
	var dokument bytes.Buffer
	if err := api.ImportImages(nil, &dokument, []io.Reader{bytes.NewReader(zrodlo.Bytes())},
		nil, nastawyPdf()); err != nil {
		return nil, fmt.Errorf("nie można osadzić obrazu w dokumencie: %w", err)
	}
	return dokument.Bytes(), nil
}

// typTresciWydaniaDesignu oddaje typ treści wedle IANA dla formatu wydania.
// Jedna prawda dla wszystkich wołających — także dla dróg, które bajtów nie
// kodują (SVG przechodzi bez dotknięcia).
func typTresciWydaniaDesignu(format string) string {
	switch normalizujFormatWydaniaDesignu(format) {
	case "png":
		return "image/png"
	case "jpeg":
		return "image/jpeg"
	case "svg":
		return "image/svg+xml"
	case "pdf":
		return "application/pdf"
	case "ico":
		return "image/vnd.microsoft.icon"
	}
	return "application/octet-stream"
}

// wBaza64Designu koduje bajty wydania w postaci, w której niesie je kontrakt
// (`contentBase64`). Jedno miejsce dla wszystkich dróg wydania, żeby zapis
// (standardowy, z dopełnieniem) nie rozjechał się między rodzinami komend.
func wBaza64Designu(bajty []byte) string {
	return base64.StdEncoding.EncodeToString(bajty)
}

// bladWydaniaDesignu nazywa niepowodzenie samego wydania: obraz nie do
// rozłożenia, koder odmówił, dokument nie powstał. Kod jest wewnętrzny, bo to
// awaria pracy rdzenia, a nie pomyłka w żądaniu.
func bladWydaniaDesignu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError, "moduł Design: "+powod))
}
