// Odpowiedzialność pliku: trzy typowane odmowy silnika mowy — brak pomocnika,
// brak silnika w Pythonie, brak użytecznego nagrania.
//
// Odmowa jest osobnym typem, a nie napisem, bo obsługiwacz komendy kontraktu
// rozróżnia te trzy przypadki przez errors.As: każdy dostaje inny kod błędu
// i inną podpowiedź naprawy. Brak Pythona to niedokończona instalacja produktu,
// brak silnika to niedoinstalowana zależność Pythona, a brak nagrania to zła
// dana przysłana przez klienta.
//
// Każdy komunikat ma trzy części: co odmówiło, dlaczego i czym to naprawić.
//
// ŻADNA ODMOWA NIE NAZYWA BRAKU, KTÓREGO NIE ZMIERZONO. Łańcuch transkrypcji ma
// trzy ogniwa dokładane osobno i naprawiane osobno: interpreter Pythona,
// biblioteka `faster-whisper` w tym interpreterze, wagi modelu na dysku. Odmowa,
// która przypisuje brak niewłaściwemu ogniwu, prowadzi Operatora do naprawy
// bezskutecznej — instalacji biblioteki do interpretera, którego nie ma. Dlatego
// rozpoznanie ma tu wartość „nie wiem" i przy niej odmowa oddaje sam zmierzony
// powód, zamiast zgadywać rozpoznanie i naprawę.
package mowa

import (
	"errors"
	"os/exec"
	"strings"
)

// BrakPomocnika jest odmową: nie ma interpretera albo nie ma skryptu
// transkrypcji. Bez pary interpreter+skrypt nie ma czego uruchomić.
type BrakPomocnika struct {
	// Szukano wylicza ścieżki, pod którymi pomocnika nie było. Wykaz idzie do
	// komunikatu dosłownie: bez niego meldunek nie prowadzi do naprawy, bo
	// Operator nie wie, gdzie plik dołożyć.
	Szukano []string
	// Powod niesie błąd źródłowy ostatniego sprawdzenia (os.Stat, exec.LookPath).
	Powod error
}

// Error mówi wprost, czego brakuje, gdzie tego szukano i czym to naprawić.
func (b *BrakPomocnika) Error() string {
	komunikat := "Brak Pythona (python_helper) do transkrypcji." +
		" Silnik mowy potrzebuje interpretera i skryptu pomocniczego transkrypcja.py"
	if len(b.Szukano) > 0 {
		komunikat += "; szukano w: " + strings.Join(b.Szukano, ", ")
	}
	komunikat += "; naprawa: dołożyć katalog pomocniki/transkrypcja obok binarium rdzenia" +
		" i zainstalować Pythona 3 na ścieżce wyszukiwania systemu"
	if b.Powod != nil {
		komunikat += " (" + b.Powod.Error() + ")"
	}
	return komunikat
}

// Unwrap oddaje błąd źródłowy sprawdzenia ścieżki.
func (b *BrakPomocnika) Unwrap() error { return b.Powod }

// BrakInterpretera jest odmową: skrypt pomocnika leży na miejscu, ale nie ma
// czym go uruchomić — interpretera Pythona nie ma na ścieżce wyszukiwania.
//
// Typ osobny od BrakSilnika, bo naprawa jest inna i pomylenie ich prowadzi
// Operatora donikąd: instalowanie biblioteki do interpretera, którego nie ma,
// kończy się drugim komunikatem o tym samym braku. Osobny też od BrakPomocnika,
// bo tam brakuje CZĘŚCI PRODUKTU (katalogu ze skryptem), a tu brakuje programu,
// który produkt zastaje na maszynie.
//
// Rozpoznanie idzie po błędzie źródłowym uruchomienia (exec.ErrNotFound), nie
// po zgadywaniu z treści wyjścia — mierzone, nie domniemane.
type BrakInterpretera struct {
	// Program to nazwa, pod którą interpretera szukano.
	Program string
	// Powod niesie błąd źródłowy uruchomienia.
	Powod error
}

// Error nazywa brakujący interpreter i podaje naprawę właściwą temu brakowi.
func (b *BrakInterpretera) Error() string {
	nazwa := strings.TrimSpace(b.Program)
	if nazwa == "" {
		nazwa = "python"
	}
	komunikat := "Interpretera Pythona nie ma na ścieżce wyszukiwania systemu." +
		" Skrypt pomocnika transkrypcji jest na miejscu, ale nie ma czym go uruchomić" +
		"; naprawa: zainstalować Pythona 3 i udostępnić go pod nazwą " + nazwa +
		" na ścieżce wyszukiwania. Biblioteki faster-whisper NIE instaluj teraz —" +
		" bez interpretera nie ma do czego jej doinstalować"
	if b.Powod != nil {
		komunikat += " (" + b.Powod.Error() + ")"
	}
	return komunikat
}

// Unwrap oddaje błąd źródłowy uruchomienia.
func (b *BrakInterpretera) Unwrap() error { return b.Powod }

// odmowaUruchomienia rozstrzyga, KTÓREGO ogniwa zabrakło, gdy uruchomienie
// pomocnika nie doszło do skutku.
//
// Bez tego rozstrzygnięcia każde niepowodzenie uruchomienia szło jako brak
// silnika — a więc odmowa twierdziła „interpreter odnaleziony" także wtedy, gdy
// w tym samym zdaniu, w nawiasie, stało `executable file not found`. Operator
// czyta zdanie główne i instaluje bibliotekę do interpretera, którego nie ma.
//
// Rozpoznanie idzie po błędzie źródłowym, nie po treści wyjścia: `exec.ErrNotFound`
// jest odpowiedzią systemu, a nie zgadywaniem z napisu. Dopasowanie po tekście
// zostaje jako druga droga, bo błąd bywa owinięty przez warstwę uruchamiania
// i wtedy nie niesie już sygnału typowanego.
func odmowaUruchomienia(program string, err error, wynik Wynik) error {
	if czyBrakInterpretera(err) {
		return &BrakInterpretera{Program: program, Powod: err}
	}
	return &BrakSilnika{Powod: powodZUruchomienia(err, wynik)}
}

// czyBrakInterpretera rozpoznaje brak programu na ścieżce wyszukiwania.
func czyBrakInterpretera(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, exec.ErrNotFound) {
		return true
	}
	return strings.Contains(err.Error(), "executable file not found")
}

// BrakSilnika jest odmową: interpreter Pythona jest, ale nie ma w nim modułu
// faster-whisper. To odrębny przypadek od BrakPomocnika, bo naprawa jest inna —
// tu niczego nie brakuje w produkcie, brakuje pakietu w środowisku Pythona.
type BrakSilnika struct {
	// Powod niesie to, co pomocnik powiedział na wyjściu diagnostycznym.
	Powod string
}

// Error nazywa brakujący moduł i podaje polecenie, którym Operator go dołoży.
//
// Zdanie NIE twierdzi, że interpreter został odnaleziony. Twierdziło tak
// wcześniej i było to twierdzenie niezmierzone: pomocnik bywa nieuruchomiony
// z wielu powodów, a odmowa niosła wtedy w nawiasie prawdę przeciwną do zdania
// głównego („executable file not found") i prowadziła Operatora do instalowania
// biblioteki dla interpretera, którego nie ma. Brak interpretera rozpoznany
// wprost ma własną odmowę (BrakInterpretera); tutaj zostaje przypadek, w którym
// uruchomienie nie powiodło się z powodu nierozstrzygniętego — więc odmowa mówi
// o najczęstszej przyczynie i o tym, czym ją sprawdzić, zamiast orzekać.
func (b *BrakSilnika) Error() string {
	komunikat := "Silnik faster-whisper niedostępny w Pythonie." +
		" Pomocnik transkrypcji nie doszedł do rozpoznania" +
		"; najczęstsza przyczyna to brak biblioteki — naprawa: pip install faster-whisper" +
		" w tym samym interpreterze, który uruchamia pomocnika." +
		" Jeżeli to nie pomoże, sprawdź, czy interpreter Pythona 3 jest osiągalny," +
		" uruchamiając pomocnika ręcznie"
	if strings.TrimSpace(b.Powod) != "" {
		komunikat += " (" + strings.TrimSpace(b.Powod) + ")"
	}
	return komunikat
}

// Unwrap nie ma czego oddać: powód przychodzi z wyjścia diagnostycznego cudzego
// procesu jako tekst, nie jako błąd Go. Metoda istnieje, żeby wszystkie trzy
// odmowy pakietu dały się obsłużyć jednakowo, i zwraca nil zgodnie z umową
// `errors.Unwrap` — nil znaczy „łańcuch kończy się tutaj”, a nie usterkę.
func (b *BrakSilnika) Unwrap() error { return nil }

// BrakNagrania jest odmową: odnośnik nie wskazuje na plik, który da się
// przepisać. Osobny typ, bo to jedyna z trzech odmów, którą wywołuje dana
// przysłana przez klienta, a nie stan instalacji — kod błędu kontraktu jest
// wtedy „zły argument”, nie „usterka rdzenia”.
type BrakNagrania struct {
	// Sciezka to odnośnik tak, jak przyszedł — bez niego Operator nie wie,
	// który plik rdzeń próbował otworzyć.
	Sciezka string
	// Powod nazywa konkretne niespełnione oczekiwanie.
	Powod string
}

// Error nazywa ścieżkę, powód i wykaz przyjmowanych formatów.
func (b *BrakNagrania) Error() string {
	sciezka := strings.TrimSpace(b.Sciezka)
	if sciezka == "" {
		sciezka = "(odnośnik pusty)"
	}
	return "nagranie " + sciezka + " nie nadaje się do transkrypcji: " + b.Powod +
		"; naprawa: wskazać istniejący, niepusty plik na dysku Operatora o rozszerzeniu " +
		strings.Join(FormatyNagran, ", ")
}

// Unwrap nie ma czego oddać: sprawdzenia nagrania rozstrzygają się na wyniku
// os.Stat i na wykazie formatów, a nie na cudzym błędzie do przekazania dalej.
func (b *BrakNagrania) Unwrap() error { return nil }
