// Odpowiedzialność pliku: przekład polecenia procesu okna na jedno wywołanie
// SSH — komenda zdalna z katalogiem i środowiskiem, wiersz argumentów
// transportu oraz plik known_hosts, z którego ssh bierze klucz hosta.
package zdalne

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"danacoconsole/server/internal/konfiguracja"
)

// Polecenie opisuje proces, który ma ruszyć na hoście zdalnym. Kształt
// odpowiada session.Polecenie, ale pakiet nie importuje warstwy sesji —
// zależność biegnie od kanału do toru, nigdy odwrotnie.
type Polecenie struct {
	// Program to plik wykonywalny na hoście zdalnym: ścieżka albo nazwa rozstrzygana na miejscu.
	Program string
	// Argumenty wiersza poleceń, bez nazwy programu.
	Argumenty []string
	// Katalog uruchomienia na hoście zdalnym; pusty oznacza katalog logowania.
	Katalog string
	// Srodowisko w postaci KLUCZ=wartość, ustawiane dla procesu zdalnego.
	Srodowisko []string
}

// Uruchomienie to gotowy wiersz poleceń transportu — program i argumenty,
// które spawner platformy (injection.Wystartuj) uruchamia na maszynie rdzenia.
type Uruchomienie struct {
	Program   string
	Argumenty []string
	// Host, do którego tor prowadzi — dla dziennika i komunikatów.
	Host Host
}

// opcjeSSH obowiązują każde połączenie toru: tryb wsadowy, wymóg klucza hosta
// znanego z wykazu i limit czasu połączenia. Klucza nieznanego tor nie
// przyjmuje — maszyną zdalną steruje wtedy ten, kto odpowiedział, a nie ten,
// z kim Operator wydał zgodę na rozmowę.
var opcjeSSH = []string{
	"-o", "BatchMode=yes",
	"-o", "StrictHostKeyChecking=yes",
	"-o", "ConnectTimeout=10",
}

// zbudujUruchomienie składa wywołanie SSH prowadzące polecenie procesu na
// wskazany host zdalny. Brak klucza hosta jest odmową — tą samą, którą oddaje
// tor plików — a nie wierszem poleceń, który ssh odrzuci dopiero przy łączeniu.
func zbudujUruchomienie(sciezkaSSH string, h Host, p Polecenie) (Uruchomienie, error) {
	plikKluczy, err := zapiszZnaneHosty(h)
	if err != nil {
		return Uruchomienie{}, err
	}
	argumenty := append([]string{}, opcjeSSH...)
	argumenty = append(argumenty, "-o", "UserKnownHostsFile="+plikKluczy)
	if h.Port != 0 && h.Port != 22 {
		argumenty = append(argumenty, "-p", strconv.Itoa(h.Port))
	}
	if uzytkownik := strings.TrimSpace(h.Uzytkownik); uzytkownik != "" {
		argumenty = append(argumenty, "-l", uzytkownik)
	}
	argumenty = append(argumenty, h.AdresPolaczenia(), "--", komendaZdalna(p))
	return Uruchomienie{Program: sciezkaSSH, Argumenty: argumenty, Host: h}, nil
}

// Funkcja komendaZdalna składa komendę wykonywaną przez powłokę logowania hosta, z katalogiem i środowiskiem.
func komendaZdalna(p Polecenie) string {
	czesci := make([]string, 0, 3+len(p.Srodowisko)+len(p.Argumenty))
	czesci = append(czesci, "exec")
	if len(p.Srodowisko) > 0 {
		czesci = append(czesci, "env")
		for _, wpis := range p.Srodowisko {
			czesci = append(czesci, cytuj(wpis))
		}
	}
	czesci = append(czesci, cytuj(p.Program))
	for _, argument := range p.Argumenty {
		czesci = append(czesci, cytuj(argument))
	}
	komenda := strings.Join(czesci, " ")
	if katalog := strings.TrimSpace(p.Katalog); katalog != "" {
		komenda = fmt.Sprintf("cd %s && %s", cytuj(katalog), komenda)
	}
	return komenda
}

// cytuj ujmuje napis w apostrofy POSIX; apostrof w treści przechodzi przez
// sekwencję '\” — jedyną, którą powłoka zgodna z POSIX czyta jednoznacznie.
func cytuj(tresc string) string {
	return "'" + strings.ReplaceAll(tresc, "'", `'\''`) + "'"
}

// zapiszZnaneHosty składa plik known_hosts hosta z klucza niesionego wierszem
// wykazu (Host.KluczHosta) i zwraca jego ścieżkę. Przy odmowie plik jest
// usunięty, więc ssh przy StrictHostKeyChecking=yes odrzuca połączenie zamiast
// przyjąć klucz podstawiony w locie. Odcisk sam w sobie by nie wystarczył: ssh
// porównuje klucz, nie jego skrót.
func zapiszZnaneHosty(h Host) (string, error) {
	sciezka := sciezkaZnanychHostow(h)
	klucz := strings.TrimSpace(h.KluczHosta)
	if klucz == "" {
		_ = os.Remove(sciezka)
		return sciezka, odmowaKluczaHosta(h.Nazwa)
	}
	if err := os.MkdirAll(filepath.Dir(sciezka), 0o700); err != nil {
		_ = os.Remove(sciezka)
		return sciezka, fmt.Errorf("zdalne: połączenie z hostem %q nie zostało nawiązane, "+
			"bo katalog na plik known_hosts nie powstał: %w", h.Nazwa, err)
	}
	if err := os.WriteFile(sciezka, []byte(wierszeZnanychHostow(h, klucz)), 0o600); err != nil {
		_ = os.Remove(sciezka)
		return sciezka, fmt.Errorf("zdalne: połączenie z hostem %q nie zostało nawiązane, "+
			"bo pliku known_hosts nie udało się zapisać: %w", h.Nazwa, err)
	}
	return sciezka, nil
}

// sciezkaZnanychHostow wskazuje plik known_hosts jednego hosta. Plik stoi
// w katalogu danych rdzenia, bo tor biegnie na koncie procesu rdzenia, a nie
// na koncie Operatora, i jest osobny dla każdego wiersza wykazu, żeby dwa
// równoczesne połączenia nie pisały po jednym pliku.
func sciezkaZnanychHostow(h Host) string {
	return filepath.Join(konfiguracja.KatalogDanychDomyslny(), "zdalne",
		fmt.Sprintf("known_hosts-%d", h.Id))
}

// wierszeZnanychHostow składa treść pliku known_hosts. Wiersze są dwa, bo ssh
// szuka hosta pod samym adresem przy porcie 22, a pod adresem w nawiasach
// kwadratowych z portem przy każdym innym.
func wierszeZnanychHostow(h Host, klucz string) string {
	adres := h.AdresPolaczenia()
	return fmt.Sprintf("%s %s\n[%s]:%d %s\n", adres, klucz, adres, port(h), klucz)
}

// odmowaKluczaHosta nazywa brak klucza publicznego hosta — trójczęściowo, jak
// pozostałe odmowy pakietu: co się nie stało, dlaczego i co to zdejmuje.
func odmowaKluczaHosta(nazwa string) error {
	return fmt.Errorf("zdalne: połączenie z hostem %q nie zostało nawiązane, bo wiersz tego "+
		"hosta nie niesie klucza publicznego, a bez niego ssh nie odróżni maszyny Operatora "+
		"od maszyny podstawionej; klucz wpisuje Operator instrukcją Danaco: UPDATE host_zdalny "+
		"SET klucz_hosta = '<typ> <klucz base64>' WHERE nazwa = '%s' — wartość bierze "+
		"z ssh-keyscan uruchomionego na maszynie rdzenia (kolumna z migracji 409)", nazwa, nazwa)
}
