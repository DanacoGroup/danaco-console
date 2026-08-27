//go:build ignore

// Wypisuje strumień tokenów pliku Go bez komentarzy — po jednym na wiersz,
// z literałami. Dowód „kod nietknięty" dla terenów komentarzowych: wynik przed
// zmianą i po niej musi być bajtowo identyczny. Strumień tokenów jest
// niewrażliwy na białe znaki, a literały łańcuchowe są pojedynczymi tokenami,
// więc dowód nie kłamie przy pustych wierszach wewnątrz surowych łańcuchów.
// Dyrektywy (//go:…) są komentarzami — porównuje się je osobno, wedle rejestru.
// Tryb `-gestosc` mierzy tym samym skanerem gęstość komentarza: wypisuje
// wiersze, znaki komentarza i granicę 250/1000 — skaner odróżnia komentarz od
// łańcucha, więc `https://` w literale nie zawyża pomiaru.
// Użycie: go build -o bezkom narzedzia/zrodlo-bez-komentarzy.go && bezkom [-gestosc] <plik.go>
package main

import (
	"bytes"
	"fmt"
	"go/scanner"
	"go/token"
	"os"
	"regexp"
	"strings"
)

// Z granicy wyłączone są dyrektywy (`//go:`). Zakazane jest odsyłanie
// czytelnika do opracowań — dokumentów wyjaśniających i zwrotów kierujących
// gdzie indziej. Nazywanie artefaktów technicznych, z którymi kod pracuje
// (biblioteka, arkusz stylu, program, migracja, plik nastaw, plik źródłowy),
// jest treścią i zakazu nie narusza. Tryb -gestosc liczy odesłania osobno
// jako POWOLANIA; plik z odesłaniem nie przechodzi niezależnie od gęstości.
var wzorOpracowania = regexp.MustCompile(`\S*\.md\b|(^|\s)(docs|prowadzenie)/`)
var wzorOdeslania = regexp.MustCompile(`(?i)(^|\s)(patrz|zob\.|uzasadnienie:|szczegoly w|szczegóły w|opisane w|opisano w|wiecej w|więcej w)\b`)

func dyrektywa(lit string) bool {
	return strings.HasPrefix(lit, "//go:")
}

func powolanie(lit string) bool {
	if dyrektywa(lit) {
		return false
	}
	t := lit
	if strings.HasPrefix(t, "//") {
		t = strings.TrimPrefix(t, "//")
	}
	return wzorOpracowania.MatchString(t) || wzorOdeslania.MatchString(t)
}

func main() {
	gestosc := len(os.Args) == 3 && os.Args[1] == "-gestosc"
	if len(os.Args) != 2 && !gestosc {
		fmt.Fprintln(os.Stderr, "użycie: zrodlo-bez-komentarzy [-gestosc] <plik.go>")
		os.Exit(2)
	}
	if gestosc {
		os.Args = os.Args[1:]
	}
	tresc, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "odczyt nie przeszedł:", err)
		os.Exit(1)
	}
	fset := token.NewFileSet()
	plik := fset.AddFile(os.Args[1], fset.Base(), len(tresc))
	var s scanner.Scanner
	blad := func(pos token.Position, msg string) {
		fmt.Fprintf(os.Stderr, "rozbiór nie przeszedł: %s: %s\n", pos, msg)
		os.Exit(1)
	}
	if gestosc {
		s.Init(plik, tresc, blad, scanner.ScanComments)
		znaki := 0
		powolania := 0
		for {
			_, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			if tok == token.COMMENT && !dyrektywa(lit) {
				znaki += len(lit)
				if powolanie(lit) {
					powolania++
				}
			}
		}
		wiersze := bytes.Count(tresc, []byte("\n")) + 1
		// Granica schodkowa (dopowiedzenie trzecie pozycji 18): 250 znaków
		// dla każdego pliku, plus 250 za każdy pełny tysiąc wierszy.
		granica := 250
		if wiersze >= 1000 {
			granica = (wiersze / 1000) * 250
		}
		stan := "W-GRANICY"
		if znaki > granica {
			stan = "PONAD"
		}
		if powolania > 0 {
			stan = "POWOLANIE"
		}
		fmt.Printf("%s\twiersze=%d\tznaki=%d\tgranica=%d\tpowolania=%d\t%s\n", os.Args[1], wiersze, znaki, granica, powolania, stan)
		return
	}
	s.Init(plik, tresc, blad, 0) // bez scanner.ScanComments — komentarze pominięte
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\n", tok, lit)
	}
}
