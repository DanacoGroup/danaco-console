// Odpowiedzialność pliku: przełożenie rodzaju powłoki karty terminala na
// polecenie uruchomienia procesu oraz na zestaw zmiennych środowiska.
//
// Karta jest profilem powłoki, nie procesem. Nie ma powłoki interaktywnej
// czekającej na wiersze: każde `terminal.command.exec` startuje własny proces
// tej powłoki w katalogu i środowisku karty. Dzięki temu Process Monitor może
// zakończyć pojedynczy proces sygnałem łagodnym albo wymuszonym; polecenia
// podanego na wejście wspólnej powłoki nie da się zakończyć inaczej niż razem
// z całą kartą. Proces polecenia ma więc własny PID, własny kod wyjścia
// i własne drzewo potomstwa.
//
// Skutkiem tej decyzji jest to, że stan powłoki nie przechodzi między
// poleceniami (`cd` nie przesuwa katalogu karty).
package core

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// zmiennaCeluSSH niosła adres powłoki zdalnej karty rodzaju `ssh`, zanim
// kontrakt dostał pole `remoteTarget`. Zostaje drogą ZASTĘPCZĄ — dla kart
// założonych wcześniej i dla klienta, który jeszcze nie przestawił się na nowe
// pola. Droga główna prowadzi przez `adapter_modul_terminal_cel.go`.
const zmiennaCeluSSH = "SSH_TARGET"

// definicjaPowloki opisuje jedną powłokę: plik wykonywalny i argumenty
// poprzedzające treść polecenia.
type definicjaPowloki struct {
	program   string
	argumenty []string
}

// powloki jest wykazem sterowanym danymi: dołożenie powłoki to dołożenie
// pozycji wykazu i wartości do wyliczenia kontraktu, nie zmiana przepływu
// sterowania.
//
// Argumenty dobrane tak, żeby proces nie czytał profilu użytkownika systemu
// i nie pytał o nic interaktywnie — inaczej wynik polecenia zależałby także od
// zawartości tego profilu.
var powloki = map[shared.TerminalShell]definicjaPowloki{
	shared.TerminalShellPowershell: {
		program:   "powershell",
		argumenty: []string{"-NoLogo", "-NoProfile", "-NonInteractive", "-Command"},
	},
	shared.TerminalShellCmd: {
		program:   "cmd",
		argumenty: []string{"/D", "/C"},
	},
	shared.TerminalShellBash: {
		program:   "bash",
		argumenty: []string{"-c"},
	},
	shared.TerminalShellNode: {
		program:   "node",
		argumenty: []string{"-e"},
	},
	shared.TerminalShellPython: {
		// `-u` wyłącza buforowanie wyjścia. Bez tego Output Console dostaje
		// wszystko dopiero na końcu procesu, a okno ma pokazywać strumień.
		program:   "python",
		argumenty: []string{"-u", "-c"},
	},
	shared.TerminalShellSsh: {
		program:   "ssh",
		argumenty: []string{"-T"},
	},
}

// CzyPowlokaZnana odpowiada, czy rodzaj powłoki należy do wykazu wykonawczego.
//
// Wykazy są dwa, bo dwa są rodzaje powłok: te uruchamiane wprost na maszynie
// rdzenia (tutaj) i te sięgające do bytu poza nią — kontenera, poda, urządzenia,
// maszyny sieciowej (`adapter_modul_terminal_powloki_urzadzen.go`). Razem
// pokrywają komplet słownika kontraktu.
func CzyPowlokaZnana(powloka shared.TerminalShell) bool {
	if _, jest := powloki[powloka]; jest {
		return true
	}
	_, urzadzeniowa := powlokiUrzadzen[powloka]
	return urzadzeniowa
}

// polecenieKarty składa polecenie uruchomienia jednego przebiegu w karcie.
// Katalog i środowisko pochodzą z karty; egzekutor izolacji dostaje polecenie
// gotowe i to on rozstrzyga, czy wolno je wykonać (izolacja.go).
func polecenieKarty(karta *kartaTerminala, tresc string) (session.Polecenie, error) {
	if _, urzadzeniowa := powlokiUrzadzen[karta.powloka]; urzadzeniowa {
		program, argumenty, err := polecenieUrzadzenia(karta, tresc)
		if err != nil {
			return session.Polecenie{}, err
		}
		return session.Polecenie{
			Program:             program,
			Argumenty:           argumenty,
			Katalog:             karta.katalog,
			Srodowisko:          srodowiskoKarty(karta),
			DziedziczSrodowisko: true,
		}, nil
	}
	definicja, jest := powloki[karta.powloka]
	if !jest {
		return session.Polecenie{}, fmt.Errorf("powłoka %q nie należy do wykazu kontraktu", karta.powloka)
	}
	argumenty := append([]string{}, definicja.argumenty...)
	if karta.powloka == shared.TerminalShellSsh {
		dodatkowe, err := argumentyPowlokiZdalnej(karta)
		if err != nil {
			return session.Polecenie{}, err
		}
		argumenty = append(argumenty, dodatkowe...)
	}
	argumenty = append(argumenty, tresc)

	return session.Polecenie{
		Program:    definicja.program,
		Argumenty:  argumenty,
		Katalog:    karta.katalog,
		Srodowisko: srodowiskoKarty(karta),
		// Dziedziczenie środowiska rdzenia zostaje decyzją izolacji, nie karty:
		// gdy punkt „środowisko procesu” jest włączony, egzekutor odrzuci
		// polecenie dziedziczące. Karta prosi o dziedziczenie, bo powłoka bez
		// PATH nie znajdzie ani jednego narzędzia.
		DziedziczSrodowisko: true,
	}, nil
}

// argumentyPowlokiZdalnej składa przełączniki `ssh` wynikające z celu karty:
// port, wskazanie klucza i sam adres.
//
// Adres bierze się z pola karty, nie ze zmiennej środowiska — pole jest źródłem
// głównym od chwili, gdy kontrakt dostał `remoteTarget`
// (`adapter_modul_terminal_cel.go`). Karta bez adresu kończy się odmową
// nazywającą oba sposoby jego podania, bo `ssh` bez celu nie ruszy i tak.
//
// Adres wchodzi jako POJEDYNCZY argument, nie jako fragment wiersza powłoki:
// nie ma tu składania napisu, więc nie ma czego wstrzyknąć spacją ani średnikiem.
func argumentyPowlokiZdalnej(karta *kartaTerminala) ([]string, error) {
	cel := strings.TrimSpace(karta.celZdalny)
	if cel == "" {
		cel = strings.TrimSpace(karta.srodowisko[zmiennaCeluSSH])
	}
	if cel == "" {
		return nil, fmt.Errorf(
			"karta powłoki zdalnej nie ma adresu — podaj pole remoteTarget, wpis książki hostów "+
				"polem hostId albo zmienną środowiska %s", zmiennaCeluSSH)
	}
	argumenty := make([]string, 0, 5)
	if karta.portZdalny > 0 {
		argumenty = append(argumenty, "-p", strconv.Itoa(karta.portZdalny))
	}
	if sciezka := strings.TrimSpace(karta.kluczSciezka); sciezka != "" {
		argumenty = append(argumenty, "-i", sciezka)
	}
	return append(argumenty, cel), nil
}

// srodowiskoKarty przekłada zmienne karty na postać KLUCZ=wartość. Kolejność
// jest ustalona przez nazwę, żeby ten sam zestaw dawał ten sam wiersz —
// inaczej ślad wykonania byłby nieporównywalny między przebiegami.
func srodowiskoKarty(karta *kartaTerminala) []string {
	if len(karta.srodowisko) == 0 {
		return nil
	}
	nazwy := make([]string, 0, len(karta.srodowisko))
	for nazwa := range karta.srodowisko {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)
	wpisy := make([]string, 0, len(nazwy))
	for _, nazwa := range nazwy {
		wpisy = append(wpisy, nazwa+"="+karta.srodowisko[nazwa])
	}
	return wpisy
}

// zmienneKarty czyta zmienne środowiska z żądania kontraktu.
//
// Zmienna wskazująca sejf, a nie wartość, kończy się odmową zamiast cichym
// pominięciem: karta terminala nie ma czytnika sejfu, więc proces ruszyłby bez
// oczekiwanego poświadczenia i nic by tego nie sygnalizowało.
func zmienneKarty(zmienne []shared.EnvironmentVariable) (map[string]string, error) {
	wynik := make(map[string]string, len(zmienne))
	for _, zmienna := range zmienne {
		nazwa := strings.TrimSpace(zmienna.Name)
		if nazwa == "" || !zmienna.Enabled {
			continue
		}
		if zmienna.Value == nil {
			if wartoscTekstu(zmienna.SecretRef) == "" {
				continue
			}
			return nil, fmt.Errorf("zmienna %s wskazuje sejf, a karta terminala nie ma jego czytnika", nazwa)
		}
		wynik[nazwa] = *zmienna.Value
	}
	return wynik, nil
}
