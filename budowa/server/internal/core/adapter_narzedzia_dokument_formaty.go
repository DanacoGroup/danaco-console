// Plik niesie słownik formatów rodziny narzędzi dokumentowych: jedną prawdę
// o tym, jak nazywa się format w kontrakcie, jak nazywa go Pandoc i jakie
// rozszerzenie nosi jego plik, wspólną dla zamiany formatu i odczytu treści.
package core

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// opisFormatuDokumentu niesie trzy nazwy jednego formatu oraz dwie zdolności,
// od których zależy dobór narzędzia rodziny.
type opisFormatuDokumentu struct {
	// pandoc jest nazwą formatu w Pandocu; pusta znaczy, że Pandoc go nie zna.
	pandoc string
	// rozszerzenie jest nazwą pliku na dysku, po której binaria rozpoznają
	// wejście bez jawnego formatu.
	rozszerzenie string
	// czytaPandoc mówi, czy Pandoc weźmie ten format jako wejście.
	czytaPandoc bool
	// piszePandoc mówi, czy Pandoc odda ten format jako wyjście.
	piszePandoc bool
	// strawnyDlaLibre mówi, czy LibreOffice otworzy plik wprost, jedyna droga PDF-u.
	strawnyDlaLibre bool
}

// formatyDokumentu jest kompletem formatów rodziny. Wykaz idzie za kontraktem,
// plus `txt`, bo tekst czysty jest naturalnym wyjściem odczytu i wejściem
// konwersji.
var formatyDokumentu = map[string]opisFormatuDokumentu{
	"markdown": {pandoc: "markdown", rozszerzenie: "md", czytaPandoc: true, piszePandoc: true},
	"html":     {pandoc: "html", rozszerzenie: "html", czytaPandoc: true, piszePandoc: true, strawnyDlaLibre: true},
	"docx":     {pandoc: "docx", rozszerzenie: "docx", czytaPandoc: true, piszePandoc: true, strawnyDlaLibre: true},
	"odt":      {pandoc: "odt", rozszerzenie: "odt", czytaPandoc: true, piszePandoc: true, strawnyDlaLibre: true},
	"rtf":      {pandoc: "rtf", rozszerzenie: "rtf", czytaPandoc: true, piszePandoc: true, strawnyDlaLibre: true},
	"epub":     {pandoc: "epub", rozszerzenie: "epub", czytaPandoc: true, piszePandoc: true},
	"csv":      {pandoc: "csv", rozszerzenie: "csv", czytaPandoc: true, piszePandoc: true, strawnyDlaLibre: true},
	// Nazwa pandokowa plain jest nazwą zapisu — czytnika o tej nazwie Pandoc nie ma.
	"txt": {pandoc: "plain", rozszerzenie: "txt", czytaPandoc: true, piszePandoc: true, strawnyDlaLibre: true},
	// PDF stoi osobno: Pandoc go nie czyta ani nie zapisuje, wymaga silnika
	// składu poza bramą rdzenia.
	"pdf": {rozszerzenie: "pdf", strawnyDlaLibre: true},
}

// obrazyDokumentu jest wykazem formatów, które są pikselami, a nie dokumentem:
// jedyną drogą do ich treści jest rozpoznanie pisma, zdjęcie kartki i skan
// wchodzą tędy.
var obrazyDokumentu = map[string]bool{
	"png": true, "jpg": true, "jpeg": true, "tif": true, "tiff": true,
	"bmp": true, "webp": true, "gif": true, "pnm": true,
}

// normalizujFormatDokumentu sprowadza wskazanie wołającego do nazwy ze
// słownika: przyjmuje wielkie litery, kropkę wiodącą, skróty i nazwy typów
// MIME.
func normalizujFormatDokumentu(wskazanie string) string {
	nazwa := strings.ToLower(strings.TrimSpace(wskazanie))
	nazwa = strings.TrimPrefix(nazwa, ".")
	if nazwa == "" {
		return ""
	}
	// Synonimy, nie zgadywanie: każda z tych nazw wskazuje dokładnie jeden format słownika.
	switch nazwa {
	case "md", "mkd", "markdown_strict", "text/markdown":
		return "markdown"
	case "htm", "xhtml", "text/html":
		return "html"
	case "text", "plain", "text/plain":
		return "txt"
	case "application/pdf":
		return "pdf"
	case "jpe":
		return "jpg"
	}
	if _, jest := formatyDokumentu[nazwa]; jest {
		return nazwa
	}
	if obrazyDokumentu[nazwa] {
		return nazwa
	}
	return ""
}

// formatZeSciezki rozpoznaje format po rozszerzeniu pliku. Pusty wynik znaczy
// „nie wiem”, a nie „format domyślny”.
func formatZeSciezki(sciezka string) string {
	return normalizujFormatDokumentu(filepath.Ext(sciezka))
}

// pierwszyFormatDokumentu oddaje pierwsze wskazanie, które da się rozpoznać.
// Kolejność argumentów ustala wołający, zwykle najpierw model, potem plik.
func pierwszyFormatDokumentu(wskazania ...string) string {
	for _, wskazanie := range wskazania {
		if nazwa := normalizujFormatDokumentu(wskazanie); nazwa != "" {
			return nazwa
		}
	}
	return ""
}

// rozszerzenieFormatuDokumentu oddaje rozszerzenie pliku dla formatu. Format
// spoza słownika oddaje własną nazwę, bo png jest i formatem, i rozszerzeniem.
func rozszerzenieFormatuDokumentu(format string) string {
	if opis, jest := formatyDokumentu[format]; jest {
		return opis.rozszerzenie
	}
	return format
}

// formatZNaglowka rozpoznaje format po pierwszych bajtach pliku, bo wynik
// konwersji leży w magazynie jako blob bez rozszerzenia, a zgadnięty format
// kończy się treścią przeczytaną nie tym słownikiem.
func formatZNaglowka(sciezka string) string {
	plik, err := os.Open(sciezka)
	if err != nil {
		return ""
	}
	defer plik.Close()

	// Cztery kilobajty mieszczą wszystkie znaczniki wraz z nazwami pierwszych pozycji ZIP.
	naglowek := make([]byte, 4096)
	odczytane, _ := plik.Read(naglowek)
	if odczytane <= 0 {
		return ""
	}
	naglowek = naglowek[:odczytane]

	switch {
	case bytes.HasPrefix(naglowek, []byte("%PDF-")):
		return "pdf"
	case bytes.HasPrefix(naglowek, []byte{0x89, 'P', 'N', 'G'}):
		return "png"
	case bytes.HasPrefix(naglowek, []byte{0xFF, 0xD8, 0xFF}):
		return "jpg"
	case bytes.HasPrefix(naglowek, []byte("GIF8")):
		return "gif"
	case bytes.HasPrefix(naglowek, []byte("BM")):
		return "bmp"
	case bytes.HasPrefix(naglowek, []byte("II*\x00")), bytes.HasPrefix(naglowek, []byte("MM\x00*")):
		return "tif"
	case bytes.HasPrefix(naglowek, []byte("RIFF")) && bytes.Contains(naglowek[:min(odczytane, 16)], []byte("WEBP")):
		return "webp"
	case bytes.HasPrefix(naglowek, []byte("{\\rtf")):
		return "rtf"
	}

	// Archiwum ZIP jest kopertą trzech formatów, sam znacznik PK nie rozstrzyga niczego.
	if bytes.HasPrefix(naglowek, []byte("PK\x03\x04")) {
		switch {
		case bytes.Contains(naglowek, []byte("mimetypeapplication/epub+zip")):
			return "epub"
		case bytes.Contains(naglowek, []byte("opendocument.text")):
			return "odt"
		case bytes.Contains(naglowek, []byte("word/")),
			bytes.Contains(naglowek, []byte("wordprocessingml")):
			return "docx"
		}
		return ""
	}

	poczatek := strings.ToLower(strings.TrimSpace(string(naglowek[:min(odczytane, 256)])))
	if strings.HasPrefix(poczatek, "<!doctype html") || strings.HasPrefix(poczatek, "<html") {
		return "html"
	}
	return ""
}

// wykazFormatowDokumentu składa posortowany wykaz nazw do treści odmowy:
// odmowa ma powiedzieć, co rodzina umie, a nie tylko czego nie umie.
func wykazFormatowDokumentu() string {
	nazwy := make([]string, 0, len(formatyDokumentu))
	for nazwa := range formatyDokumentu {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)
	return strings.Join(nazwy, ", ")
}
