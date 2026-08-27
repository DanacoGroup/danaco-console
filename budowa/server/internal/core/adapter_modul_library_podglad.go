// Odpowiedzialność pliku: moduł Library — odczyt treści pliku do podglądu
// tekstowego (`library.file.preview`). Tu leży sam odczyt treści spod
// odwołania pliku, wydzielony jako osobna odpowiedzialność.
package core

import (
	"bytes"
	"errors"
	"io"
	"os"
	"unicode/utf8"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// granicaPodgladuBiblioteki jest domyślną górną granicą długości podglądu
// tekstowego, gdy żądanie nie poda własnej (`maxChars`).
const granicaPodgladuBiblioteki = 64 << 10

// progRozpoznaniaTekstu mówi, ile pierwszych bajtów treści wystarcza, by
// rozstrzygnąć „tekst czy nie tekst” przy pliku bez podanego typu MIME.
const progRozpoznaniaTekstu = 8 << 10

// trescPodgladuBiblioteki odczytuje treść pliku spod odwołania i zwraca
// jej początek do podglądu tekstowego wraz z informacją, czy została
// skrócona.
func trescPodgladuBiblioteki(odwolanie string, maxZnakow *int) (string, bool, error) {
	limitZnakow := granicaPodgladuBiblioteki
	if maxZnakow != nil && *maxZnakow > 0 {
		limitZnakow = *maxZnakow
	}

	plik, err := os.Open(odwolanie)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", false, bladOdczytuTresciBiblioteki(
				"treść nie leży pod odwołaniem " + odwolanie + " — podgląd nie ma czego pokazać")
		}
		return "", false, bladOdczytuTresciBiblioteki(
			"nie można otworzyć treści spod odwołania " + odwolanie + ": " + err.Error())
	}
	defer plik.Close()

	// Każdy znak UTF-8 to najwyżej cztery bajty, więc tyle bajtów wystarcza
	// na limitZnakow znaków.
	gornaBajtow := limitZnakow*4 + 1
	bajty := make([]byte, gornaBajtow)
	n, err := io.ReadFull(plik, bajty)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", false, bladOdczytuTresciBiblioteki(
			"nie można odczytać treści spod odwołania " + odwolanie + ": " + err.Error())
	}
	bajty = bajty[:n]

	znaki := []rune(string(bajty))
	if len(znaki) > limitZnakow {
		return string(znaki[:limitZnakow]), true, nil
	}
	return string(znaki), false, nil
}

// bladOdczytuTresciBiblioteki nazywa niepowodzenie odczytu treści spod
// odwołania pliku. To awaria odczytu, nie wina Operatora, więc kod jest
// wewnętrzny.
func bladOdczytuTresciBiblioteki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError, "moduł Library: "+powod))
}

// odwolanieObrazuPodgladu rozstrzyga, co trafi do pola `imageRef`
// odpowiedzi `library.file.preview`. Dwa zawężenia sprawiają, że
// najczęściej nie trafia nic: tylko obraz, i tylko postać względna.
func odwolanieObrazuPodgladu(rodzaj shared.LibraryPreviewKind, odwolanie string) *string {
	if rodzaj != shared.LibraryPreviewKindImage {
		return nil
	}
	wzgledne := odwolanieTresciBiblioteki(odwolanie)
	if wzgledne == "" {
		return nil
	}
	return &wzgledne
}

// rodzajPodgladuTresci rozstrzyga rodzaj podglądu, gdy typ MIME nie został
// podany albo nie mówi nic znanego. Rozstrzygają bajty, nie domysł.
func rodzajPodgladuTresci(odwolanie string) shared.LibraryPreviewKind {
	plik, err := os.Open(odwolanie)
	if err != nil {
		return shared.LibraryPreviewKindBinary
	}
	defer plik.Close()

	glowa := make([]byte, progRozpoznaniaTekstu)
	n, err := io.ReadFull(plik, glowa)
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		return shared.LibraryPreviewKindBinary
	}
	glowa = glowa[:n]
	if bytes.IndexByte(glowa, 0) >= 0 {
		return shared.LibraryPreviewKindBinary
	}
	// Ostatni znak głowy bywa urwany w sekwencji UTF-8; niepełny ogon jest
	// odcinany przed sprawdzeniem.
	if n == progRozpoznaniaTekstu {
		for i := 0; i < 3 && len(glowa) > 0 && !utf8.Valid(glowa); i++ {
			glowa = glowa[:len(glowa)-1]
		}
	}
	if !utf8.Valid(glowa) {
		return shared.LibraryPreviewKindBinary
	}
	return shared.LibraryPreviewKindText
}

// rodzajPodgladu rozstrzyga rodzaj podglądu: najpierw z typu MIME pliku,
// a gdy typu nie ma albo nic on nie mówi — z samej treści.
func rodzajPodgladu(mimeType *string, odwolanie string) shared.LibraryPreviewKind {
	if mimeType == nil || *mimeType == "" {
		return rodzajPodgladuTresci(odwolanie)
	}
	switch {
	case *mimeType == "application/pdf":
		return shared.LibraryPreviewKindPdf
	case len(*mimeType) >= 5 && (*mimeType)[:5] == "text/":
		return shared.LibraryPreviewKindText
	case *mimeType == "application/json":
		return shared.LibraryPreviewKindText
	case len(*mimeType) >= 6 && (*mimeType)[:6] == "image/":
		return shared.LibraryPreviewKindImage
	default:
		// Typ podany, ale rdzeniowi nieznany. Zamiast zgadywać binary, sprawdzane
		// są bajty treści.
		return rodzajPodgladuTresci(odwolanie)
	}
}
