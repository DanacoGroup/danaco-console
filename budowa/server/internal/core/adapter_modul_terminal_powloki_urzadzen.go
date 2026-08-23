// Odpowiedzialność pliku: cztery powłoki, które nie są programem uruchamianym
// wprost na maszynie rdzenia — kontener, pod, port szeregowy i sesja Telnet.
//
// ── Dlaczego osobny plik ────────────────────────────────────────────────────
// Powłoki z `adapter_modul_terminal_powloki.go` są jednym wierszem wykazu:
// program i argumenty poprzedzające treść. Te cztery wymagają czegoś więcej —
// wskazania CELU, którym jest kontener, pod, urządzenie albo adres sieciowy —
// więc ich argumenty składa się z pól karty, a nie z samej stałej.
//
// ── Skąd wiadomo, że pole karty jest wypełnione ─────────────────────────────
// Odmowa idzie przy SKŁADANIU polecenia, a nie przy otwarciu karty. Karta jest
// profilem powłoki i wolno ją otworzyć z pustym wskazaniem, tak samo jak kartę
// zdalną bez adresu; dopiero polecenie musi wiedzieć, gdzie się wykonać.
//
// ── Czego te powłoki NIE robią ──────────────────────────────────────────────
// Nie zarządzają kontenerem, podem ani urządzeniem: nie zakładają, nie usuwają
// i nie zmieniają ich stanu. Wykonują polecenie w bycie, który już istnieje.
// Kontenerami zarządza moduł Developer własną drogą (biblioteka Dockera), a nie
// karta terminala.
package core

import (
	"fmt"
	"strconv"
	"strings"

	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// Programy powłok urządzeniowych. Deklaracje stoją tu, a sonda startowa bierze
// je stąd (`zaleznosci_zewnetrzne.go`) — nazwa wypisana drugi raz rozjechałaby
// się z pierwszą przy pierwszej zmianie pakietu.
var (
	narzedzieDocker = zewnetrzne.Narzedzie{
		Nazwa: "Docker (klient wiersza poleceń)", Program: "docker", Pakiet: "docker.io",
	}
	narzedzieKubectl = zewnetrzne.Narzedzie{
		Nazwa: "kubectl", Program: "kubectl", Pakiet: "kubectl (snap)",
	}
	narzedziePicocom = zewnetrzne.Narzedzie{
		Nazwa: "picocom", Program: "picocom", Pakiet: "picocom",
	}
	narzedzieTelnet = zewnetrzne.Narzedzie{
		Nazwa: "Telnet (klient)", Program: "telnet", Pakiet: "telnet",
	}
)

// domyslnaPredkoscPortu obowiązuje, gdy karta nie poda własnej. 115200 bodów to
// wartość, z którą pracuje dziś przytłaczająca większość konsol szeregowych
// sprzętu sieciowego i płytek deweloperskich.
const domyslnaPredkoscPortu = 115200

// powlokiUrzadzen wiąże cztery powłoki z programem i sposobem złożenia jego
// argumentów. Wykaz jest danymi, tak samo jak `powloki`: dołożenie powłoki to
// dołożenie pozycji, nie nowa gałąź warunku.
var powlokiUrzadzen = map[shared.TerminalShell]struct {
	narzedzie zewnetrzne.Narzedzie
	// argumenty składa wiersz uruchomienia dla karty wraz z treścią polecenia.
	argumenty func(karta *kartaTerminala, tresc string) ([]string, error)
}{
	shared.TerminalShellContainer: {
		narzedzie: narzedzieDocker,
		argumenty: func(karta *kartaTerminala, tresc string) ([]string, error) {
			kontener := strings.TrimSpace(karta.kontener)
			if kontener == "" {
				return nil, fmt.Errorf(
					"karta powłoki kontenera nie wskazuje kontenera — podaj pole containerRef.containerId")
			}
			// `-i` bez `-t`: proces rdzenia nie ma terminala znakowego, a `-t`
			// zażądałoby go i uruchomienie odpadłoby na maszynie bez konsoli.
			argumenty := []string{"exec", "-i"}
			if katalog := strings.TrimSpace(karta.katalog); katalog != "" {
				argumenty = append(argumenty, "--workdir", katalog)
			}
			return append(argumenty, kontener, "sh", "-c", tresc), nil
		},
	},
	shared.TerminalShellPod: {
		narzedzie: narzedzieKubectl,
		argumenty: func(karta *kartaTerminala, tresc string) ([]string, error) {
			pod := strings.TrimSpace(karta.pod)
			if pod == "" {
				return nil, fmt.Errorf(
					"karta powłoki poda nie wskazuje poda — podaj pole containerRef.podName")
			}
			argumenty := make([]string, 0, 10)
			if kontekst := strings.TrimSpace(karta.kontekstKlastra); kontekst != "" {
				argumenty = append(argumenty, "--context", kontekst)
			}
			if przestrzen := strings.TrimSpace(karta.przestrzenNazw); przestrzen != "" {
				argumenty = append(argumenty, "--namespace", przestrzen)
			}
			argumenty = append(argumenty, "exec", "-i", pod)
			// Kontener wewnątrz poda wskazuje się wtedy, gdy pod ma ich więcej
			// niż jeden; bez wskazania `kubectl` bierze pierwszy i mówi o tym sam.
			if kontener := strings.TrimSpace(karta.kontener); kontener != "" {
				argumenty = append(argumenty, "--container", kontener)
			}
			return append(argumenty, "--", "sh", "-c", tresc), nil
		},
	},
	shared.TerminalShellSerial: {
		narzedzie: narzedziePicocom,
		argumenty: func(karta *kartaTerminala, tresc string) ([]string, error) {
			urzadzenie := strings.TrimSpace(karta.urzadzenie)
			if urzadzenie == "" {
				return nil, fmt.Errorf(
					"karta konsoli szeregowej nie wskazuje urządzenia — podaj pole serialDevice")
			}
			predkosc := karta.predkoscPortu
			if predkosc <= 0 {
				predkosc = domyslnaPredkoscPortu
			}
			// `--exit-after` domyka sesję po chwili ciszy na łączu. Bez tego
			// picocom trzymałby port otwarty bez końca, a polecenie karty nigdy
			// nie miałoby kodu wyjścia — proces bez końca to proces, którego
			// Process Monitor nie ma jak domknąć.
			return []string{
				"--baud", strconv.Itoa(predkosc),
				"--exit-after", "2000",
				"--quiet",
				"--initstring", tresc + "\r",
				urzadzenie,
			}, nil
		},
	},
	shared.TerminalShellTelnet: {
		narzedzie: narzedzieTelnet,
		argumenty: func(karta *kartaTerminala, _ string) ([]string, error) {
			cel := strings.TrimSpace(karta.celZdalny)
			if cel == "" {
				cel = strings.TrimSpace(karta.srodowisko[zmiennaCeluSSH])
			}
			if cel == "" {
				return nil, fmt.Errorf(
					"karta sesji Telnet nie ma adresu — podaj pole remoteTarget")
			}
			// Telnet nie przyjmuje polecenia argumentem: rozmawia strumieniem.
			// Adres bierze się z pola karty, a nazwa użytkownika — jeśli wpisana
			// w postaci `użytkownik@host` — odpada, bo Telnet pyta o nią sam,
			// wewnątrz sesji.
			if miejsce := strings.LastIndex(cel, "@"); miejsce >= 0 {
				cel = cel[miejsce+1:]
			}
			argumenty := []string{cel}
			if karta.portZdalny > 0 {
				argumenty = append(argumenty, strconv.Itoa(karta.portZdalny))
			}
			return argumenty, nil
		},
	},
}

// polecenieUrzadzenia składa polecenie dla powłoki urządzeniowej.
//
// Program odnajduje się TERAZ, a nie przy uruchomieniu: `session.Polecenie`
// niesie ścieżkę, a program dołożony do pakietu produktu leży poza ścieżką
// wyszukiwania systemu i po samej nazwie by nie wystartował. Brak programu
// kończy się odmową nazywającą go wraz z pakietem — tą samą, którą oddaje cały
// arsenał (`zewnetrzne.BrakNarzedzia`).
func polecenieUrzadzenia(karta *kartaTerminala, tresc string) (string, []string, error) {
	definicja, jest := powlokiUrzadzen[karta.powloka]
	if !jest {
		return "", nil, fmt.Errorf("powłoka %q nie należy do wykazu powłok urządzeniowych", karta.powloka)
	}
	sciezka, stoi := zewnetrzne.Odnajdz(definicja.narzedzie)
	if !stoi {
		return "", nil, &zewnetrzne.BrakNarzedzia{Narzedzie: definicja.narzedzie}
	}
	argumenty, err := definicja.argumenty(karta, tresc)
	if err != nil {
		return "", nil, err
	}
	return sciezka, argumenty, nil
}

// wskazanieUrzadzenia przepisuje do karty wskazanie kontenera, poda albo portu
// szeregowego z żądania otwarcia karty.
//
// Sprawdzenie kompletności wskazania leży przy składaniu polecenia, nie tutaj:
// karta jest profilem powłoki i wolno ją otworzyć z polem pustym, tak samo jak
// kartę zdalną bez adresu.
func wskazanieUrzadzenia(karta *kartaTerminala, z shared.TerminalSessionOpenRequest) {
	if z.ContainerRef != nil {
		karta.kontener = strings.TrimSpace(wartoscTekstu(z.ContainerRef.ContainerId))
		karta.pod = strings.TrimSpace(wartoscTekstu(z.ContainerRef.PodName))
		karta.przestrzenNazw = strings.TrimSpace(wartoscTekstu(z.ContainerRef.Namespace))
		karta.kontekstKlastra = strings.TrimSpace(wartoscTekstu(z.ContainerRef.Context))
		// Pod ma własne pole kontenera; gdy je podano, ono rozstrzyga.
		if nazwa := strings.TrimSpace(wartoscTekstu(z.ContainerRef.ContainerName)); nazwa != "" {
			karta.kontener = nazwa
		}
	}
	karta.urzadzenie = strings.TrimSpace(wartoscTekstu(z.SerialDevice))
	if z.SerialBaudRate != nil && *z.SerialBaudRate > 0 {
		karta.predkoscPortu = *z.SerialBaudRate
	}
}
