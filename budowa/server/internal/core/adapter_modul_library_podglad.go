// Odpowiedzialność pliku: moduł Library — odczyt treści pliku do podglądu
// tekstowego (`library.file.preview`). Metoda `Podglad`
// (`adapter_modul_library.go`) rozstrzyga rodzaj podglądu i składa odpowiedź
// kontraktu; tutaj leży sam odczyt treści spod odwołania pliku, wydzielony, bo
// to osobna odpowiedzialność — dostęp do nośnika.
//
// Odwołanie jest ścieżką na dysku i tylko tyle ten plik o nim zakłada.
// Wgranie przez `sourcePath` czyni ścieżkę odwołaniem do treści; wgranie przez
// `contentBase64` odkłada bajty w magazynie treści rdzenia i odwołaniem czyni
// ścieżkę bloba (`Wgraj`, `adapter_modul_library_magazyn.go`). Obie drogi kończą
// się ścieżką, więc czytelnik jest jeden i czyta ją z dysku tak, jak moduł
// Developer czyta pliki repozytorium do Code Editor
// (`adapter_modul_developer_plik.go`). Plik bez odwołania — wiersz sprzed
// wpięcia magazynu albo wersja będąca samym znacznikiem — jest odmawiany
// wcześniej (`bladBrakuTresciBiblioteki`), więc tu trafia wyłącznie plik
// z realnym odwołaniem.
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
// tekstowego, gdy żądanie nie poda własnej (`maxChars`). Podgląd ma pokazać
// początek dokumentu w oknie modułu, a nie przesłać cały plik przez gniazdo
// zdarzeń — stąd granica na znaki, nie odczyt całości.
const granicaPodgladuBiblioteki = 64 << 10

// progRozpoznaniaTekstu mówi, ile pierwszych bajtów treści wystarcza, by
// rozstrzygnąć „tekst czy nie tekst" przy pliku bez podanego typu MIME. Osiem
// kibibajtów to głowa dokumentu — dość, by trafić na bajt zerowy albo na
// niepoprawny UTF-8 formatu binarnego, i mało, by nie czytać całości.
const progRozpoznaniaTekstu = 8 << 10

// trescPodgladuBiblioteki odczytuje treść pliku spod odwołania i zwraca jej
// początek do podglądu tekstowego wraz z informacją, czy została skrócona.
// Odwołanie to ścieżka pliku na dysku (patrz `Wgraj`, `sourcePath`), więc rdzeń
// czyta ją tak, jak moduł Developer czyta pliki repozytorium. Czyta o jeden
// bajt więcej niż górna granica w bajtach — nadmiar mówi, że treść jest dłuższa
// niż podgląd, więc podgląd jest skrócony. Odwołanie, którego nie da się
// odczytać, jest odmową wprost, a nie pustą treścią udającą podgląd.
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

	// Każdy znak UTF-8 to najwyżej cztery bajty, więc tyle bajtów wystarcza na
	// `limitZnakow` znaków; jeden bajt ponad granicę rozpoznaje nadmiar treści.
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
// odwołania pliku — plik ma odwołanie, lecz rdzeń nie mógł go otworzyć albo
// wczytać (ścieżka zniknęła, brak praw, błąd nośnika). To awaria odczytu, nie
// wina Operatora, więc kod jest wewnętrzny; podgląd woli odmówić niż zmyślić
// treść.
func bladOdczytuTresciBiblioteki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError, "moduł Library: "+powod))
}

// odwolanieObrazuPodgladu rozstrzyga, co trafi do pola `imageRef` odpowiedzi
// `library.file.preview`. Dwa zawężenia sprawiają, że najczęściej nie trafia nic.
//
// Pierwsze: tylko obraz. Kontrakt nazywa to pole odwołaniem do obrazu, więc
// wypełnia je wyłącznie podgląd obrazowy. Przy `binary` czy `pdf` klient dostałby
// wskazanie, którego nie umie otworzyć, a przy okazji wyszedłby na zewnątrz układ
// katalogu danych Operatora. Podgląd binarny zostaje więc samym rodzajem.
//
// Drugie: tylko postać względna. Wychodzi ścieżka względna magazynu, ta sama
// postać, co `asset.uri` modułu Design — ścieżka bezwzględna nazywałaby katalog
// danych rdzenia, a odbiorca bywa na innej maszynie i tak by jej nie otworzył.
// Wybór postaci opisuje `odwolanieMagazynu` (`adapter_modul_library_magazyn.go`).
//
// Pustka znaczy „nie mam czego podać", a nie „obraz bez treści": odwołanie
// spoza magazynu (ścieżka źródłowa zamiast bloba) nie wychodzi wcale, bo pole
// kontraktu jest niewymagane, a ścieżki zastępczej odbiorca nie odróżniłby
// od prawdziwej.
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
// podany albo nie mówi nic znanego. Kontrakt czyni `mimeType` polem
// opcjonalnym, więc jego brak jest zwyczajny, nie wyjątkowy.
//
// Rozstrzygają bajty, nie domysł: plik tekstowy bez `mimeType` uznany z góry
// za `binary` nie pokazałby ani jednego znaku treści, choć rdzeń ma i bajty,
// i czytnik tekstu. Bajt zerowy albo ciąg niebędący poprawnym UTF-8 znaczy
// treść nietekstową i wtedy `binary` jest prawdą.
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
	// Ostatni znak głowy bywa urwany w połowie sekwencji UTF-8 — to nie czyni
	// pliku binarnym, więc niepełny ogon odcinamy przed sprawdzeniem. Sekwencja
	// ma najwyżej cztery bajty, więc urwany ogon to najwyżej trzy: dalsze
	// skracanie zjadałoby treść binarną aż do pustej, a pusta przechodzi jako
	// poprawny UTF-8 i cały plik wyszedłby tekstem.
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

// rodzajPodgladu rozstrzyga rodzaj podglądu: najpierw z typu MIME pliku, a gdy
// typu nie ma albo nic on nie mówi — z samej treści (`rodzajPodgladuTresci`).
// Stoi tutaj, przy odczycie treści podglądu, a nie przy składaniu odpowiedzi
// (`adapter_modul_library.go`), bo to jedna decyzja tego samego obszaru.
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
		// Typ podany, ale rdzeniowi nieznany (np. „application/octet-stream"
		// przysłane hurtem przez okno). Zamiast zgadywać `binary`, pytamy bajty
		// — one wiedzą lepiej niż etykieta, którą ktoś nadał na wejściu.
		return rodzajPodgladuTresci(odwolanie)
	}
}
