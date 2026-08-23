// Odpowiedzialność pliku: rozstrzygnięcie odnośnika nagrania (AudioRef) w plik
// na dysku, który pomocnik da radę przepisać.
//
// AudioRef jest ścieżką pliku na dysku Operatora, nie treścią zakodowaną i nie
// identyfikatorem w składnicy: rdzeń nie kopiuje bajtów. Nagranie zostaje tam,
// gdzie je nagrano, a pomocnik otwiera je w miejscu. Przyjmowanie treści
// przepuszczałoby każde nagranie przez pamięć procesu i przez gniazdo —
// kilkadziesiąt megabajtów na jedno zdanie — a rdzeń musiałby je gdzieś odłożyć,
// czyli prowadzić drugą składnicę plików obok tej, którą Operator już ma.
//
// Czego ten plik nie robi: nie otwiera deskryptora, nie czyta ani jednego bajtu
// treści i nie sprawdza, czy plik jest naprawdę dźwiękiem. Sprawdzenie
// rozszerzenia jest bramką na oczywiste pomyłki (odnośnik na dokument), nie
// rozpoznaniem formatu — rozpoznaje go silnik, bo tylko on wie, co potrafi
// zdekodować.
package mowa

import (
	"os"
	"path/filepath"
	"strings"
)

// FormatyNagran wylicza rozszerzenia przyjmowane od klienta: trzy pierwsze to
// typowe zapisy z dyktafonu, czwarte to format, w którym nagrywa przeglądarka.
var FormatyNagran = []string{".wav", ".ogg", ".m4a", ".webm"}

// Nagranie to rozstrzygnięty odnośnik: ścieżka bezwzględna i rozmiar w bajtach.
// Rozmiar jedzie razem ze ścieżką, bo wołający decyduje na jego podstawie
// o granicy czasu transkrypcji, a drugie odpytanie dysku dałoby inną wartość niż
// ta, na której podjęto decyzję o dopuszczeniu pliku.
type Nagranie struct {
	// Sciezka — bezwzględna ścieżka pliku na dysku Operatora.
	Sciezka string
	// Rozmiar — rozmiar pliku w bajtach w chwili sprawdzenia.
	Rozmiar int64
}

// OtworzNagranie rozstrzyga odnośnik i sprawdza, czy da się go przepisać.
//
// Każda odmowa jest *BrakNagrania nazywającym ścieżkę oraz wykaz przyjmowanych
// formatów — komunikat „nie da się” bez tych dwóch rzeczy nie prowadzi do
// naprawy, bo klient nie wie, czy zawinił plik, czy jego rodzaj.
//
// Kolejność sprawdzeń idzie od najtańszego do najdroższego i od najbardziej
// ogólnego do najbardziej szczegółowego: pusty odnośnik, rozwinięcie ścieżki,
// istnienie, rodzaj wpisu, rozmiar, rozszerzenie. Rozszerzenie sprawdzane jest
// ostatnie: gdy pliku nie ma, zdanie o nieprzyjmowanym formacie byłoby
// odpowiedzią na niezadane pytanie.
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
			Powod: "pliku nie ma albo rdzeń nie ma do niego dostępu: " + err.Error()}
	}
	if opis.IsDir() {
		return Nagranie{}, &BrakNagrania{Sciezka: sciezka,
			Powod: "wskazanie prowadzi do katalogu, a miało prowadzić do pliku nagrania"}
	}
	if opis.Size() <= 0 {
		// Plik pusty przechodzi przez silnik bez błędu i oddaje pustą
		// transkrypcję, którą wołający wziąłby za udaną — a to byłaby atrapa
		// wyniku. Zatrzymujemy go tutaj, póki wiadomo dlaczego.
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

// rozszerzenie oddaje rozszerzenie ścieżki sprowadzone do małych liter.
func rozszerzenie(sciezka string) string {
	return strings.ToLower(filepath.Ext(sciezka))
}

// rozwinSciezke doprowadza odnośnik do postaci bezwzględnej.
//
// Tylda rozwijana jest tutaj, bo powłoka rozwija ją sama, a rdzeń polecenia od
// klienta nie dostaje przez powłokę — ścieżka „~/nagrania/x.wav” trafiłaby do
// os.Stat dosłownie i dała odmowę „pliku nie ma” dla pliku, który jest.
// Ścieżka względna liczy się od katalogu bieżącego rdzenia; to droga zapasowa
// dla wywołań ręcznych, klient ma przysyłać ścieżkę bezwzględną.
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
