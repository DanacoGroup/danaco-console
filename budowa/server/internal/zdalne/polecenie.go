// Odpowiedzialność pliku: przekład polecenia procesu okna na jedno wywołanie
// SSH — komenda zdalna z katalogiem i środowiskiem oraz wiersz argumentów
// transportu. Czysty przekład napisów, bez uruchamiania i bez bazy.
package zdalne

import (
	"fmt"
	"strconv"
	"strings"
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

// opcjeSSH obowiązują każde połączenie toru: tryb wsadowy, przyjęcie nowego klucza i limit czasu połączenia.
var opcjeSSH = []string{
	"-o", "BatchMode=yes",
	"-o", "StrictHostKeyChecking=accept-new",
	"-o", "ConnectTimeout=10",
}

// Funkcja zbudujUruchomienie składa wywołanie SSH prowadzące polecenie procesu na wskazany host zdalny.
func zbudujUruchomienie(sciezkaSSH string, h Host, p Polecenie) Uruchomienie {
	argumenty := append([]string{}, opcjeSSH...)
	if h.Port != 0 && h.Port != 22 {
		argumenty = append(argumenty, "-p", strconv.Itoa(h.Port))
	}
	if uzytkownik := strings.TrimSpace(h.Uzytkownik); uzytkownik != "" {
		argumenty = append(argumenty, "-l", uzytkownik)
	}
	argumenty = append(argumenty, h.AdresPolaczenia(), "--", komendaZdalna(p))
	return Uruchomienie{Program: sciezkaSSH, Argumenty: argumenty, Host: h}
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
