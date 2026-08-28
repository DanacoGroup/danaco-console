// Plik obsługuje przełożenie rodzaju powłoki karty terminala na polecenie uruchomienia procesu oraz na zestaw zmiennych środowiska. Karta jest profilem powłoki, nie procesem: każde polecenie startuje własny proces.
package core

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// zmiennaCeluSSH niosła adres powłoki zdalnej karty rodzaju `ssh`, zanim kontrakt dostał pole `remoteTarget`. Zostaje drogą zastępczą dla kart założonych wcześniej. Droga główna prowadzi przez `adapter_modul_terminal_cel.go`.
const zmiennaCeluSSH = "SSH_TARGET"

// definicjaPowloki opisuje jedną powłokę: plik wykonywalny i argumenty poprzedzające treść polecenia uruchomienia procesu w karcie.
type definicjaPowloki struct {
	program   string
	argumenty []string
}

// powloki jest wykazem sterowanym danymi: dołożenie powłoki to dołożenie pozycji wykazu, nie zmiana przepływu sterowania. Argumenty dobrane tak, żeby proces nie czytał profilu użytkownika systemu i nie pytał o nic interaktywnie.
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
		// `-u` wyłącza buforowanie wyjścia, żeby Output Console pokazywał strumień, nie zwał na końcu.
		program:   "python",
		argumenty: []string{"-u", "-c"},
	},
	shared.TerminalShellSsh: {
		program:   "ssh",
		argumenty: []string{"-T"},
	},
}

// CzyPowlokaZnana odpowiada, czy rodzaj powłoki należy do wykazu wykonawczego. Wykazy są dwa: te uruchamiane wprost na maszynie rdzenia i te sięgające do bytu poza nią (`adapter_modul_terminal_powloki_urzadzen.go`).
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
		// Dziedziczenie środowiska jest decyzją izolacji, nie karty; bez PATH powłoka nie działa.
		DziedziczSrodowisko: true,
	}, nil
}

// argumentyPowlokiZdalnej składa przełączniki `ssh` wynikające z celu karty: port, wskazanie klucza i sam adres. Adres bierze się z pola karty, nie ze zmiennej środowiska. Adres wchodzi jako pojedynczy argument, nie jako fragment wiersza powłoki.
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

// zmienneKarty czyta zmienne środowiska z żądania kontraktu. Zmienna wskazująca sejf, a nie wartość, kończy się odmową zamiast cichym pominięciem: karta terminala nie ma czytnika sejfu.
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
