package core

// Rozdział na warstwę jest jedną funkcją z jednym sprawdzianem, nie formatowaniem w skrypcie.

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"danacoconsole/server/internal/mowa"
)

// Znaczniki wiersza poleceń przełączające rdzeń w tryb wypisania wykazu zamiast
// pracy właściwej. Przyjmowane w obu postaciach zapisu flagi (`-` i `--`), bo
// obie są w użyciu w tym projekcie.
const (
	znacznikWykazuZaleznosci = "wykaz-zaleznosci"
	znacznikWykazuMowy       = "wykaz-mowy"
)

// Nazwy warstw prowizjonowania. Warstwa mówi, CZYM postawić program — a przy
// silniku kontenerów: że nie stawiać go bez wyraźnego żądania.
const (
	// WarstwaObowiazkowa — pakiety dystrybucji podawane programowi apt wprost oraz polecenia
	// `pip install`, dla których prowizjonowanie ma osobną gałąź.
	WarstwaObowiazkowa = "obowiazkowa-apt"
	// WarstwaObowiazkowaRecznie — pozycje obowiązkowe, których prowizjonowanie nie postawi samo:
	// podpowiedzi zapisane zdaniem oraz polecenia menedżerów spoza apt i pip. Rozbicie takiego
	// pola na spacjach dałoby programowi apt nazwy, których nie zna, i przerwało przebieg.
	WarstwaObowiazkowaRecznie = "obowiazkowa-recznie"
	// WarstwaWarsztatGo — programy dokładane przez go install, osobno od pakietów dystrybucji tego systemu.
	WarstwaWarsztatGo = "warsztat-go"
	// WarstwaWarsztatNpm — programy dokładane przez npm i -g, osobno od pakietów dystrybucji tego systemu.
	WarstwaWarsztatNpm = "warsztat-npm"
	// WarstwaSnap — programy, których dystrybucja nie ma w apt, dokładane przez menedżera snap tego systemu.
	WarstwaSnap = "snap"
	// WarstwaModelRecznie — silniki i wagi z wydań spoza repozytoriów
	// dystrybucji; kroki ręczne, bo automat na nich zawodzi.
	WarstwaModelRecznie = "model-recznie"
	// WarstwaDecyzyjna — silnik kontenerów, wstrzymany, więc prowizjonowanie go pomija, dopóki nie zażąda się go wprost.
	WarstwaDecyzyjna = "decyzyjna"
)

// ZadanoWykazZaleznosci mówi, czy w argumentach procesu stoi żądanie wypisania wykazu zależności; sprawdzane przed odczytem konfiguracji, bo wykaz nie potrzebuje bazy ani katalogu danych.
func ZadanoWykazZaleznosci(argumenty []string) bool {
	return zadanoZnacznik(argumenty, znacznikWykazuZaleznosci)
}

// ZadanoWykazMowy mówi, czy zażądano wypisania arsenału mowy w postaci nadającej się do odczytu maszynowego.
func ZadanoWykazMowy(argumenty []string) bool {
	return zadanoZnacznik(argumenty, znacznikWykazuMowy)
}

// zadanoZnacznik odpowiada, czy argumenty procesu niosą wskazany znacznik trybu wypisania wykazu zależności.
func zadanoZnacznik(argumenty []string, znacznik string) bool {
	for _, argument := range argumenty {
		if argument == "-"+znacznik || argument == "--"+znacznik {
			return true
		}
	}
	return false
}

// WarstwaZaleznosci rozstrzyga warstwę pozycji wykazu: silnik kontenerów po programie, warsztaty
// po przedrostku polecenia, a pozycje obowiązkowe po kształcie pola. Dopasowanie warsztatu Go
// bierze przedrostek, nie podnapis: „cargo install typos-cli" niesie „go install" wewnątrz
// „[car]go install", a mimo to nie jest poleceniem Go.
//
// Pole obowiązkowe idzie do warstwy podawanej programowi apt wyłącznie wtedy, gdy w całości
// składa się z nazw pakietów dystrybucji. Zdanie i polecenie obcego menedżera trafiają do
// warstwy ręcznej, ponieważ prowizjonowanie rozbija pole warstwy apt na spacjach.
func WarstwaZaleznosci(pozycja ZaleznoscZewnetrzna) string {
	program := strings.TrimSpace(pozycja.Narzedzie.Program)
	pakiet := strings.TrimSpace(pozycja.Narzedzie.Pakiet)

	if program == "docker" || program == "podman" {
		return WarstwaDecyzyjna
	}
	switch {
	case strings.HasPrefix(pakiet, "go install "):
		return WarstwaWarsztatGo
	case strings.HasPrefix(pakiet, "npm "):
		return WarstwaWarsztatNpm
	case strings.Contains(pakiet, "(snap)"):
		return WarstwaSnap
	case strings.Contains(pakiet, "github.com"),
		strings.Contains(pakiet, "środowisku pythonowym"):
		return WarstwaModelRecznie
	case strings.HasPrefix(pakiet, "pip install "):
		return WarstwaObowiazkowa
	case wykazNazwPakietow(pakiet):
		return WarstwaObowiazkowa
	default:
		return WarstwaObowiazkowaRecznie
	}
}

// nazwaPakietuDystrybucji dopasowuje nazwę pakietu wedle polityki nazw Debiana: mała litera
// albo cyfra na początku, dalej litery, cyfry, kropka, plus i minus.
var nazwaPakietuDystrybucji = regexp.MustCompile(`^[a-z0-9][a-z0-9.+-]*$`)

// wykazNazwPakietow mówi, czy pole składa się wyłącznie z nazw pakietów dystrybucji rozdzielonych
// spacją. Polecenie obcego menedżera nazwy przypomina, więc rozstrzyga drugi człon: pole, którego
// drugim członem jest „install", jest poleceniem, nie wykazem.
func wykazNazwPakietow(pole string) bool {
	czlony := strings.Fields(pole)
	if len(czlony) == 0 {
		return false
	}
	if len(czlony) > 1 && czlony[1] == "install" {
		return false
	}
	for _, czlon := range czlony {
		if !nazwaPakietuDystrybucji.MatchString(czlon) {
			return false
		}
	}
	return true
}

// WypiszWykazZaleznosci wypisuje komplet zależności zewnętrznych, po jednym wierszu na pozycję wykazu, w polach rozdzielonych znakiem tabulacji: warstwa, program, pakiet, stoi, nazwa, zakres.
func WypiszWykazZaleznosci(wyjscie io.Writer) error {
	if _, err := fmt.Fprintln(wyjscie,
		"# wykaz zależności zewnętrznych serwera — pola: "+
			"warstwa\tprogram\tpakiet\tstoi\tnazwa\tzakres"); err != nil {
		return err
	}
	for _, pozycja := range ZaleznosciZewnetrzne() {
		stoi := "nie"
		if pozycja.Stoi {
			stoi = "tak"
		}
		if _, err := fmt.Fprintf(wyjscie, "%s\t%s\t%s\t%s\t%s\t%s\n",
			WarstwaZaleznosci(pozycja),
			pozycja.Narzedzie.Program,
			pozycja.Narzedzie.Pakiet,
			stoi,
			pozycja.Narzedzie.Nazwa,
			pozycja.Zakres,
		); err != nil {
			return err
		}
	}
	return nil
}

// WypiszWykazMowy wypisuje arsenał mowy — tę jego część, której nie widać w wykazie zależności zewnętrznych, bo stawiana jest inaczej niż pakietem dystrybucji, wierszem klucz-wartość.
func WypiszWykazMowy(wyjscie io.Writer) error {
	wiersze := [][2]string{
		{"# arsenał mowy — pola: klucz", "wartość"},
		{"synteza.piper-program", silnikPiper},
		{"synteza.piper-arsenal", piperArsenalu},
		{"synteza.piper-glosy", glosyPiperaArsenalu},
		{"synteza.piper-zmienna", zmiennaPipera},
		{"synteza.piper-glosy-zmienna", zmiennaGlosowPiper},
		{"synteza.espeak-program", silnikEspeak},
		{"synteza.espeak-zmienna", zmiennaEspeaka},
		{"rozpoznanie.model-domyslny", mowa.ModelDomyslny},
		{"rozpoznanie.ustawienie-modelu", mowa.KluczModel},
		{"rozpoznanie.ustawienie-katalogu-modeli", mowa.KluczKatalogModeli},
		{"rozpoznanie.ustawienie-interpretera", mowa.KluczProgram},
	}
	wiersze = append(wiersze, wierszePomocnikaMowy()...)
	for _, wiersz := range wiersze {
		if _, err := fmt.Fprintf(wyjscie, "%s\t%s\n", wiersz[0], wiersz[1]); err != nil {
			return err
		}
	}
	return nil
}

// wierszePomocnikaMowy oddaje położenie pomocnika transkrypcji wraz z plikiem jego zależności pythonowych; gdy skryptu nie widać, wypisuje przeszukane miejsca jako wskazówkę naprawy.
func wierszePomocnikaMowy() [][2]string {
	pomocnik, err := mowa.OdnajdzPomocnika("")
	if err != nil {
		var brak *mowa.BrakPomocnika
		if errors.As(err, &brak) {
			return [][2]string{
				{"rozpoznanie.pomocnik-skrypt", ""},
				{"rozpoznanie.pomocnik-szukano", strings.Join(brak.Szukano, " ")},
			}
		}
		return [][2]string{{"rozpoznanie.pomocnik-skrypt", ""}}
	}
	return [][2]string{
		{"rozpoznanie.pomocnik-skrypt", pomocnik.Skrypt},
		{"rozpoznanie.pomocnik-interpreter", pomocnik.Program},
		{"rozpoznanie.wymagania", filepath.Join(filepath.Dir(pomocnik.Skrypt), "wymagania.txt")},
	}
}
