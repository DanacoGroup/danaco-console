package core

// Rozdział na warstwę jest jedną funkcją z jednym sprawdzianem, nie formatowaniem w skrypcie.

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
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
	// WarstwaObowiazkowa — pakiety dystrybucji, bez których moduły odmawiają działania na tej maszynie
	// budującej. Niesie też podpowiedzi zdaniem albo poleceniem pip/cargo — arsenal-serwera.sh
	// rozpoznaje ich postać i wyprowadza z nich krok apt, pip albo krok ręczny.
	WarstwaObowiazkowa = "obowiazkowa-apt"
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
// i kroki ręczne po przedrostku podpowiedzi, reszta do warstwy obowiązkowej. Dopasowanie warsztatu
// Go bierze przedrostek, nie podnapis: „cargo install typos-cli" niesie „go install" wewnątrz
// „[car]go install", a mimo to nie jest poleceniem Go.
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
	default:
		return WarstwaObowiazkowa
	}
}

// WypiszWykazZaleznosci wypisuje komplet zależności zewnętrznych, po jednym wierszu na pozycję wykazu, w polach rozdzielonych znakiem tabulacji: warstwa, program, pakiet, stoi, nazwa, zakres.
func WypiszWykazZaleznosci(wyjscie io.Writer) error {
	if _, err := fmt.Fprintln(wyjscie,
		"# wykaz zależności zewnętrznych rdzenia — pola: "+
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
