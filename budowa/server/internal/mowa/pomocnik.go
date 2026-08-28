// Odpowiedzialność pliku: odnalezienie pary interpreter+skrypt, którą rdzeń
// wywoła, żeby przepisać nagranie na tekst; plik niczego nie uruchamia.
package mowa

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	// nazwaSkryptu jest nazwą skryptu pomocniczego transkrypcji, uruchamianego
	// przez znaleziony interpreter.
	nazwaSkryptu = "transkrypcja.py"
	// interpreterPreferowany to interpreter szukany, gdy Operator nie wskazał
	// własnego programu do uruchomienia skryptu.
	interpreterPreferowany = "python3"
	// interpreterZapasowy wchodzi tam, gdzie `python3` nie istnieje jako osobne
	// polecenie w systemie operacyjnym.
	interpreterZapasowy = "python"
)

// katalogSkryptu to ścieżka względna pomocnika wewnątrz pakietu produktu.
// Wykaz członów, a nie gotowy napis z ukośnikami, bo separator ścieżki różni się
// między systemami, a filepath.Join zna ten właściwy.
var katalogSkryptu = []string{"pomocniki", "transkrypcja", nazwaSkryptu}

// Pomocnik to para, bez której transkrypcja nie ruszy: interpreter i skrypt.
// Trzymane razem, bo osobno nie znaczą nic — sam interpreter nie wie, co ma
// zrobić, a sam skrypt nie jest wykonywalny.
type Pomocnik struct {
	// Program — plik wykonywalny interpretera.
	Program string
	// Skrypt — bezwzględna ścieżka skryptu transkrypcji.
	Skrypt string
}

// OdnajdzPomocnika wskazuje parę interpreter+skrypt albo mówi, czemu jej nie
// ma; skryptu szuka obok binarium, które właśnie pracuje, potem drogą
// zapasową liczoną od katalogu bieżącego rdzenia.
func OdnajdzPomocnika(program string) (Pomocnik, error) {
	skrypt, szukano, err := odnajdzSkrypt()
	if err != nil {
		return Pomocnik{}, &BrakPomocnika{Szukano: szukano, Powod: err}
	}
	return Pomocnik{Program: odnajdzInterpreter(program), Skrypt: skrypt}, nil
}

// odnajdzSkrypt przechodzi obie drogi i oddaje wykaz sprawdzonych miejsc.
// Wykaz liczy się przy odmowie: jest tam jedyną wskazówką naprawy.
func odnajdzSkrypt() (string, []string, error) {
	var szukano []string
	var ostatni error

	if biezace, err := os.Executable(); err == nil {
		obok := filepath.Join(append([]string{filepath.Dir(biezace)}, katalogSkryptu...)...)
		szukano = append(szukano, obok)
		if err := plikUzyteczny(obok); err == nil {
			return obok, szukano, nil
		} else {
			ostatni = err
		}
	}

	if katalog, err := os.Getwd(); err == nil {
		wBiezacym := filepath.Join(append([]string{katalog}, katalogSkryptu...)...)
		szukano = append(szukano, wBiezacym)
		if err := plikUzyteczny(wBiezacym); err == nil {
			return wBiezacym, szukano, nil
		} else {
			ostatni = err
		}
	}

	return "", szukano, ostatni
}

// odnajdzInterpreter rozstrzyga, który program uruchomi skrypt; wskazanie
// Operatora wchodzi wprost, bez sprawdzania na dysku.
func odnajdzInterpreter(program string) string {
	if program != "" {
		return program
	}
	// LookPath, nie uruchomienie: pytanie systemu o położenie pliku interpretera.
	if zeSciezki, err := exec.LookPath(interpreterPreferowany); err == nil {
		return zeSciezki
	}
	return interpreterZapasowy
}

// plikUzyteczny odpowiada, czy pod ścieżką stoi plik, a nie katalog i nie
// nic, zanim rdzeń spróbuje go uruchomić.
func plikUzyteczny(sciezka string) error {
	opis, err := os.Stat(sciezka)
	if err != nil {
		return err
	}
	if opis.IsDir() {
		return fmt.Errorf("%s jest katalogiem, a miał być plikiem skryptu", sciezka)
	}
	return nil
}

// Argumenty składa wiersz wywołania pomocnika w jednym miejscu; `-X utf8`
// idzie zawsze i nie jest opcją wołającego.
func (p Pomocnik) Argumenty(dalsze ...string) []string {
	argumenty := make([]string, 0, 3+len(dalsze))
	argumenty = append(argumenty, "-X", "utf8", p.Skrypt)
	return append(argumenty, dalsze...)
}
