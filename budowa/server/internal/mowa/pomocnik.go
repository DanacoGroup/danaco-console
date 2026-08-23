// Odpowiedzialność pliku: odnalezienie pary interpreter+skrypt, którą rdzeń
// wywoła, żeby przepisać nagranie na tekst.
//
// Wzorzec odnajdywania jest ten sam, co w `narzedzia/wpiecie.go`: najpierw obok
// binarium rdzenia, potem droga zapasowa, a gdy zawiodą obie — typowany błąd
// zamiast ścieżki. Ścieżka zmyślona jest gorsza od odmowy, bo odmowę da się
// powiedzieć Operatorowi, a zmyślona ścieżka wraca dopiero jako niezrozumiały
// błąd uruchomienia procesu.
//
// Ten plik niczego nie uruchamia. Składa jedynie nazwę programu i wykaz
// argumentów; start procesu należy wyłącznie do portu session.Uruchamiacz.
// Jedyny wyjątek to `exec.LookPath` — ono nie uruchamia procesu, tylko
// przegląda ścieżkę wyszukiwania systemu, dokładnie tak jak robi to
// `narzedzia/wpiecie.go`.
package mowa

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	// nazwaSkryptu jest nazwą skryptu pomocniczego transkrypcji.
	nazwaSkryptu = "transkrypcja.py"
	// interpreterPreferowany to interpreter szukany, gdy Operator nie wskazał
	// własnego. Trójka jawnie, bo goła nazwa `python` na wielu systemach wciąż
	// wskazuje wydanie drugie, w którym pomocnik się nie uruchomi.
	interpreterPreferowany = "python3"
	// interpreterZapasowy wchodzi tam, gdzie `python3` nie istnieje jako osobne
	// polecenie (typowo Windows i część obrazów kontenerowych).
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

// OdnajdzPomocnika wskazuje parę interpreter+skrypt albo mówi, czemu jej nie ma.
//
// Skryptu szuka obok binarium, które właśnie pracuje: rdzeń i katalog
// `pomocniki` wychodzą z jednego pakowania i stoją w jednym katalogu. Droga
// zapasowa to ta sama ścieżka względna liczona od katalogu bieżącego rdzenia —
// obsługuje uruchomienie z `go run`, gdzie binarium stoi w katalogu tymczasowym
// i obok niego nie ma niczego z produktu.
//
// Interpreter: wskazanie Operatora wchodzi wprost, bez sprawdzania na dysku,
// bo może być nazwą do rozwinięcia przez system albo dowiązaniem środowiska
// wirtualnego, a odmowa na podstawie własnego sprawdzenia unieważniałaby to
// ustawienie. Gdy wskazania nie ma, szukamy `python3`, a gdy i tego nie ma —
// zostaje `python`. Ta ostatnia wartość jest zgadywana i może nie istnieć;
// odmowa przyjdzie wtedy z uruchomienia procesu, bo tylko ono zna prawdę
// o wykonywalności.
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

// odnajdzInterpreter rozstrzyga, który program uruchomi skrypt.
func odnajdzInterpreter(program string) string {
	if program != "" {
		return program
	}
	// LookPath, nie uruchomienie: pytamy system o położenie pliku, nie startujemy
	// procesu. To ta sama droga, którą idzie `narzedzia/wpiecie.go`.
	if zeSciezki, err := exec.LookPath(interpreterPreferowany); err == nil {
		return zeSciezki
	}
	return interpreterZapasowy
}

// plikUzyteczny odpowiada, czy pod ścieżką stoi plik, a nie katalog i nie nic.
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

// Argumenty składa wiersz wywołania pomocnika w jednym miejscu.
//
// `-X utf8` idzie zawsze i nie jest opcją wołającego: pomocnik oddaje tekst
// transkrypcji na standardowe wyjście, a Python bez tego przełącznika koduje je
// według ustawień regionalnych systemu. Na polskim Windowsie znaczy to stronę
// kodową 1250, w której transkrypcja rozpada się na krzaki, zanim rdzeń zdąży ją
// odczytać. Wymuszenie UTF-8 w jednym miejscu jest jedyną obroną, która nie
// zależy od tego, kto pomocnika woła.
func (p Pomocnik) Argumenty(dalsze ...string) []string {
	argumenty := make([]string, 0, 3+len(dalsze))
	argumenty = append(argumenty, "-X", "utf8", p.Skrypt)
	return append(argumenty, dalsze...)
}
