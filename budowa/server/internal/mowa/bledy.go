// Odpowiedzialność pliku: trzy typowane odmowy silnika mowy — brak
// pomocnika, brak silnika w Pythonie, brak użytecznego nagrania; każda
// niesie co, dlaczego i naprawę.
package mowa

import (
	"errors"
	"os/exec"
	"strings"
)

// BrakPomocnika jest odmową: nie ma interpretera albo nie ma skryptu
// transkrypcji; bez pary interpreter+skrypt nie ma czego uruchomić.
type BrakPomocnika struct {
	// Szukano wylicza ścieżki, pod którymi pomocnika nie było — wykaz idzie
	// do komunikatu.
	Szukano []string
	// Powod niesie błąd źródłowy ostatniego sprawdzenia (os.Stat, exec.LookPath).
	Powod error
}

// Error mówi wprost, czego brakuje, gdzie tego szukano i czym to naprawić,
// w kolejności zrozumiałej dla Operatora czytającego zdanie.
func (b *BrakPomocnika) Error() string {
	komunikat := "Brak Pythona (python_helper) do transkrypcji." +
		" Silnik mowy potrzebuje interpretera i skryptu pomocniczego transkrypcja.py"
	if len(b.Szukano) > 0 {
		komunikat += "; szukano w: " + strings.Join(b.Szukano, ", ")
	}
	komunikat += "; naprawa: dołożyć katalog pomocniki/transkrypcja obok binarium serwera" +
		" i zainstalować Pythona 3 na ścieżce wyszukiwania systemu"
	if b.Powod != nil {
		komunikat += " (" + b.Powod.Error() + ")"
	}
	return komunikat
}

// Unwrap oddaje błąd źródłowy sprawdzenia ścieżki, zgodnie z umową
// errors.Unwrap obowiązującą pozostałe odmowy pakietu.
func (b *BrakPomocnika) Unwrap() error { return b.Powod }

// BrakInterpretera jest odmową: skrypt pomocnika leży na miejscu, ale nie
// ma czym go uruchomić — interpretera Pythona nie ma na ścieżce
// wyszukiwania.
type BrakInterpretera struct {
	// Program to nazwa, pod którą interpretera szukano.
	Program string
	// Powod niesie błąd źródłowy uruchomienia.
	Powod error
}

// Error nazywa brakujący interpreter i podaje naprawę właściwą temu
// brakowi, odróżnioną od naprawy brakującego modułu.
func (b *BrakInterpretera) Error() string {
	nazwa := strings.TrimSpace(b.Program)
	if nazwa == "" {
		nazwa = "python"
	}
	komunikat := "Interpretera Pythona nie ma na ścieżce wyszukiwania systemu." +
		" Skrypt pomocnika transkrypcji jest na miejscu, ale nie ma czym go uruchomić" +
		"; instalacja jest niepełna — brakuje Pythona 3 pod nazwą " + nazwa +
		" na ścieżce wyszukiwania, a wraz z nim biblioteki faster-whisper"
	if b.Powod != nil {
		komunikat += " (" + b.Powod.Error() + ")"
	}
	return komunikat
}

// Unwrap oddaje błąd źródłowy uruchomienia, zgodnie z umową errors.Unwrap
// obowiązującą pozostałe odmowy pakietu.
func (b *BrakInterpretera) Unwrap() error { return b.Powod }

// odmowaUruchomienia rozstrzyga, KTÓREGO ogniwa zabrakło, gdy uruchomienie
// pomocnika nie doszło do skutku, zamiast zgadywać z treści wyjścia.
func odmowaUruchomienia(program string, err error, wynik Wynik) error {
	if czyBrakInterpretera(err) {
		return &BrakInterpretera{Program: program, Powod: err}
	}
	return &BrakSilnika{Powod: powodZUruchomienia(err, wynik)}
}

// czyBrakInterpretera rozpoznaje brak programu na ścieżce wyszukiwania po
// błędzie źródłowym, nie po treści wyjścia procesu.
func czyBrakInterpretera(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, exec.ErrNotFound) {
		return true
	}
	return strings.Contains(err.Error(), "executable file not found")
}

// BrakSilnika jest odmową: interpreter Pythona jest, ale nie ma w nim
// modułu faster-whisper; odrębny przypadek od BrakPomocnika.
type BrakSilnika struct {
	// Powod niesie to, co pomocnik powiedział na wyjściu diagnostycznym.
	Powod string
}

// Error nazywa brakujący moduł i podaje polecenie, którym Operator go
// dołoży, bez twierdzenia o interpreterze, którego nie sprawdzono.
func (b *BrakSilnika) Error() string {
	komunikat := "Silnik faster-whisper niedostępny w Pythonie." +
		" Pomocnik transkrypcji nie doszedł do rozpoznania" +
		"; najczęstsza przyczyna to brak biblioteki faster-whisper w tym samym" +
		" interpreterze, który uruchamia pomocnika — instalacja serwera jest wtedy" +
		" niepełna. Drugą możliwością jest interpreter Pythona 3 poza zasięgiem"
	if strings.TrimSpace(b.Powod) != "" {
		komunikat += " (" + strings.TrimSpace(b.Powod) + ")"
	}
	return komunikat
}

// Unwrap nie ma czego oddać: powód przychodzi z wyjścia diagnostycznego
// cudzego procesu jako tekst, nie jako błąd Go.
func (b *BrakSilnika) Unwrap() error { return nil }

// BrakNagrania jest odmową: odnośnik nie wskazuje na plik, który da się
// przepisać; jedyna z trzech odmów wywołana daną przysłaną przez klienta.
type BrakNagrania struct {
	// Sciezka to odnośnik tak, jak przyszedł, bez niego Operator nie wie,
	// który plik rdzeń otwierał.
	Sciezka string
	// Powod nazywa konkretne niespełnione oczekiwanie.
	Powod string
}

// Error nazywa ścieżkę, powód i wykaz przyjmowanych formatów, żeby
// Operator wiedział, co poprawić w przysłanym odnośniku.
func (b *BrakNagrania) Error() string {
	sciezka := strings.TrimSpace(b.Sciezka)
	if sciezka == "" {
		sciezka = "(odnośnik pusty)"
	}
	return "nagranie " + sciezka + " nie nadaje się do transkrypcji: " + b.Powod +
		"; naprawa: wskazać istniejący, niepusty plik na dysku Operatora o rozszerzeniu " +
		strings.Join(FormatyNagran, ", ")
}

// Unwrap nie ma czego oddać: sprawdzenia nagrania rozstrzygają się na
// wyniku os.Stat i na wykazie formatów, nie na cudzym błędzie.
func (b *BrakNagrania) Unwrap() error { return nil }
