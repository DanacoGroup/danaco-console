// Odpowiedzialność pliku: słownik formatów rodziny narzędzi dokumentowych —
// jedna prawda o tym, jak nazywa się format w kontrakcie, jak nazywa go Pandoc
// i jakie rozszerzenie nosi jego plik. Stoi osobno, bo z tej samej tabeli
// korzystają obie czynności: zamiana formatu i odczyt treści.
//
// Trzy nazwy jednej rzeczy — i dlatego tabela jest jedna. „markdown" kontraktu
// to `markdown` dla Pandoca i `.md` na dysku; „txt" to `plain` dla Pandoca
// i `.txt` na dysku. Trzy osobne mapy rozjechałyby się przy pierwszym
// dołożonym formacie.
//
// Format nieznany jest odmową, nie domysłem. Wołający dostaje wykaz tego, co
// rodzina umie — zgadnięty format kończy się plikiem, którego nikt nie otworzy,
// a koperta meldowałaby powodzenie.
package core

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// opisFormatuDokumentu niesie trzy nazwy jednego formatu oraz dwie zdolności,
// od których zależy dobór narzędzia.
type opisFormatuDokumentu struct {
	// pandoc jest nazwą formatu w Pandocu. Pusta znaczy „Pandoc tego nie zna" —
	// tak jest z PDF-em, którego Pandoc ani nie czyta, ani nie zapisuje bez
	// silnika składu.
	pandoc string
	// rozszerzenie jest nazwą pliku na dysku, po której binaria rozpoznają
	// wejście, gdy nie podamy formatu jawnie.
	rozszerzenie string
	// czytaPandoc mówi, czy Pandoc weźmie ten format jako wejście.
	czytaPandoc bool
	// piszePandoc mówi, czy Pandoc odda ten format jako wyjście.
	piszePandoc bool
	// strawnyDlaLibre mówi, czy LibreOffice otworzy plik wprost. Tą drogą
	// jedzie jedyne wytworzenie PDF-u, jakie ta maszyna potrafi (patrz
	// `adapter_narzedzia_dokument_konwersja.go`).
	strawnyDlaLibre bool
}

// formatyDokumentu jest kompletem formatów rodziny. Wykaz idzie za kontraktem:
// „markdown, html, docx, odt, pdf, epub, rtf, csv" — plus `txt`, bo tekst
// czysty jest naturalnym wyjściem odczytu i wejściem konwersji, a jego brak
// zmuszałby model do udawania, że notatka jest markdownem.
var formatyDokumentu = map[string]opisFormatuDokumentu{
	"markdown": {pandoc: "markdown", rozszerzenie: "md", czytaPandoc: true, piszePandoc: true},
	"html":     {pandoc: "html", rozszerzenie: "html", czytaPandoc: true, piszePandoc: true, strawnyDlaLibre: true},
	"docx":     {pandoc: "docx", rozszerzenie: "docx", czytaPandoc: true, piszePandoc: true, strawnyDlaLibre: true},
	"odt":      {pandoc: "odt", rozszerzenie: "odt", czytaPandoc: true, piszePandoc: true, strawnyDlaLibre: true},
	"rtf":      {pandoc: "rtf", rozszerzenie: "rtf", czytaPandoc: true, piszePandoc: true, strawnyDlaLibre: true},
	"epub":     {pandoc: "epub", rozszerzenie: "epub", czytaPandoc: true, piszePandoc: true},
	"csv":      {pandoc: "csv", rozszerzenie: "csv", czytaPandoc: true, piszePandoc: true, strawnyDlaLibre: true},
	// Nazwa pandokowa `plain` jest nazwą ZAPISU — czytnika o tej nazwie Pandoc
	// nie ma, więc odczyt treści pliku tekstowego idzie wprost z dysku
	// (`adapter_narzedzia_dokument_tekst.go`), a Pandoc dostaje ten format
	// wyłącznie jako cel zamiany.
	"txt": {pandoc: "plain", rozszerzenie: "txt", czytaPandoc: true, piszePandoc: true, strawnyDlaLibre: true},
	// PDF stoi osobno: Pandoc go nie czyta (nie ma z czego złożyć struktury),
	// więc `czytaPandoc` zostaje fałszem, a `piszePandoc` — także, bo zapis
	// PDF-u Pandokiem wymaga silnika składu wołanego przez niego samego, czyli
	// procesu poza bramą rdzenia. Czytaniem PDF-u zajmuje się
	// `document.text.extract`, a zapisem dwie drogi rdzenia opisane
	// w `adapter_narzedzia_dokument_konwersja.go`: skład typstem albo
	// LibreOffice.
	"pdf": {rozszerzenie: "pdf", strawnyDlaLibre: true},
}

// obrazyDokumentu jest wykazem formatów, które są pikselami, a nie dokumentem:
// jedyną drogą do ich treści jest rozpoznanie pisma. Zdjęcie kartki i skan
// wchodzą tędy.
var obrazyDokumentu = map[string]bool{
	"png": true, "jpg": true, "jpeg": true, "tif": true, "tiff": true,
	"bmp": true, "webp": true, "gif": true, "pnm": true,
}

// normalizujFormatDokumentu sprowadza wskazanie wołającego do nazwy ze
// słownika. Przyjmuje pisownię, którą naprawdę przyśle model: wielkie litery,
// kropkę wiodącą, skróty i nazwy typów MIME w części po ukośniku.
func normalizujFormatDokumentu(wskazanie string) string {
	nazwa := strings.ToLower(strings.TrimSpace(wskazanie))
	nazwa = strings.TrimPrefix(nazwa, ".")
	if nazwa == "" {
		return ""
	}
	// Synonimy, nie zgadywanie: każda z tych nazw wskazuje dokładnie jeden
	// format słownika i nie ma drugiego kandydata.
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
// „nie wiem", a nie „format domyślny".
func formatZeSciezki(sciezka string) string {
	return normalizujFormatDokumentu(filepath.Ext(sciezka))
}

// pierwszyFormatDokumentu oddaje pierwsze wskazanie, które da się rozpoznać.
// Kolejność argumentów ustala wołający — zwykle najpierw to, co powiedział
// model, potem to, co widać po pliku.
func pierwszyFormatDokumentu(wskazania ...string) string {
	for _, wskazanie := range wskazania {
		if nazwa := normalizujFormatDokumentu(wskazanie); nazwa != "" {
			return nazwa
		}
	}
	return ""
}

// rozszerzenieFormatuDokumentu oddaje rozszerzenie pliku dla formatu. Format
// spoza słownika (obraz) oddaje własną nazwę — `png` jest i formatem, i
// rozszerzeniem.
func rozszerzenieFormatuDokumentu(format string) string {
	if opis, jest := formatyDokumentu[format]; jest {
		return opis.rozszerzenie
	}
	return format
}

// formatZNaglowka rozpoznaje format po pierwszych bajtach pliku.
//
// Wynik `document.convert` leży w magazynie jako blob, którego nazwą jest suma
// SHA256 — bez kropki i bez rozszerzenia. Rozpoznanie wyłącznie po rozszerzeniu
// zrywałoby więc łańcuch najbardziej w tej rodzinie naturalny: „zamień na PDF,
// a potem przeczytaj, co wyszło".
//
// To jest odczyt, a nie domysł po nazwie: bajty `%PDF-` na początku pliku są
// definicją PDF-u, a nie poszlaką — plik, który je niesie, jest PDF-em
// niezależnie od tego, jak się nazywa. Materiał, którego nagłówek nie mówi nic
// pewnego, oddaje pustkę, bo „nie wiem" jest odpowiedzią uczciwą, a zgadnięty
// format kończy się odmową narzędzia w połowie pracy albo, gorzej, treścią
// przeczytaną nie tym słownikiem.
func formatZNaglowka(sciezka string) string {
	plik, err := os.Open(sciezka)
	if err != nil {
		return ""
	}
	defer plik.Close()

	// Cztery kilobajty: mieszczą wszystkie znaczniki poniżej wraz z nazwami
	// pierwszych pozycji archiwum ZIP, a kosztują jeden odczyt.
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

	// Archiwum ZIP jest kopertą trzech różnych formatów, więc sam znacznik `PK`
	// nie rozstrzyga niczego — rozstrzyga zawartość pierwszych pozycji.
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

// wykazFormatowDokumentu składa posortowany wykaz nazw do treści odmowy.
// Odmowa ma powiedzieć, co rodzina umie, a nie tylko czego nie umie.
func wykazFormatowDokumentu() string {
	nazwy := make([]string, 0, len(formatyDokumentu))
	for nazwa := range formatyDokumentu {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)
	return strings.Join(nazwy, ", ")
}
