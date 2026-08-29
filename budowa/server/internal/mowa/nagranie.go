// Odpowiedzialność pliku: rozstrzygnięcie odnośnika nagrania (AudioRef) w plik
// na dysku, który pomocnik da radę przepisać; rdzeń nie kopiuje bajtów ani nie
// czyta treści.
package mowa

import (
	"os"
	"path/filepath"
	"strings"
)

// FormatyNagran wylicza rozszerzenia przyjmowane od klienta: trzy pierwsze to
// typowe zapisy z dyktafonu, czwarte to format, w którym nagrywa przeglądarka.
var FormatyNagran = []string{".wav", ".ogg", ".m4a", ".webm"}

// Nagranie to rozstrzygnięty odnośnik: ścieżka bezwzględna i rozmiar w bajtach;
// rozmiar jedzie razem ze ścieżką, bo wołający decyduje na jego podstawie
// o granicy czasu transkrypcji.
type Nagranie struct {
	// Sciezka — bezwzględna ścieżka pliku na dysku Operatora.
	Sciezka string
	// Rozmiar — rozmiar pliku w bajtach w chwili sprawdzenia.
	Rozmiar int64
}

// OtworzNagranie rozstrzyga odnośnik i sprawdza, czy da się go przepisać;
// każda odmowa jest *BrakNagrania nazywającym ścieżkę oraz wykaz przyjmowanych
// formatów.
func OtworzNagranie(odnosnik string) (Nagranie, error) {
	odnosnik = strings.TrimSpace(odnosnik)
	if odnosnik == "" {
		return Nagranie{}, &BrakNagrania{Sciezka: odnosnik,
			Powod: "odnośnik nagrania jest pusty, więc nie ma czego otworzyć"}
	}

	sciezka, err := rozwinSciezke(odnosnik)
	if err != nil {
		return Nagranie{}, &BrakNagrania{Sciezka: odnosnik,
			Powod: "ścieżki nie da się rozwinąć do postaci bezwzględnej: " + err.Error()}
	}

	opis, err := os.Stat(sciezka)
	if err != nil {
		return Nagranie{}, &BrakNagrania{Sciezka: sciezka,
			Powod: "pliku nie ma albo serwer nie ma do niego dostępu: " + err.Error()}
	}
	if opis.IsDir() {
		return Nagranie{}, &BrakNagrania{Sciezka: sciezka,
			Powod: "wskazanie prowadzi do katalogu, a miało prowadzić do pliku nagrania"}
	}
	if opis.Size() <= 0 {
		// Plik pusty przechodzi przez silnik bez błędu i oddaje pustą
		// transkrypcję wziętą za udany wynik.
		return Nagranie{}, &BrakNagrania{Sciezka: sciezka,
			Powod: "plik ma zerowy rozmiar, więc nie niesie dźwięku do przepisania"}
	}
	if !FormatPrzyjmowany(sciezka) {
		return Nagranie{}, &BrakNagrania{Sciezka: sciezka,
			Powod: "rozszerzenie " + rozszerzenie(sciezka) + " nie jest przyjmowane"}
	}

	return Nagranie{Sciezka: sciezka, Rozmiar: opis.Size()}, nil
}

// FormatPrzyjmowany odpowiada, czy rozszerzenie ścieżki stoi w wykazie.
// Porównanie bez względu na wielkość liter, bo dyktafony zapisują `.WAV`
// równie często jak `.wav`.
func FormatPrzyjmowany(sciezka string) bool {
	rodzaj := rozszerzenie(sciezka)
	for _, przyjmowane := range FormatyNagran {
		if rodzaj == przyjmowane {
			return true
		}
	}
	return false
}

// rozszerzenie oddaje rozszerzenie ścieżki sprowadzone do małych liter, do
// porównania niezależnego od wielkości liter.
func rozszerzenie(sciezka string) string {
	return strings.ToLower(filepath.Ext(sciezka))
}

// rozwinSciezke doprowadza odnośnik do postaci bezwzględnej; tylda rozwijana
// jest tutaj, bo powłoka rozwija ją sama, a rdzeń polecenia od klienta nie
// dostaje przez powłokę.
func rozwinSciezke(odnosnik string) (string, error) {
	if odnosnik == "~" || strings.HasPrefix(odnosnik, "~/") || strings.HasPrefix(odnosnik, `~\`) {
		dom, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		odnosnik = filepath.Join(dom, strings.TrimLeft(odnosnik[1:], `/\`))
	}
	if filepath.IsAbs(odnosnik) {
		return filepath.Clean(odnosnik), nil
	}
	return filepath.Abs(odnosnik)
}
