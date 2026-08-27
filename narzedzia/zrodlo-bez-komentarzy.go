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
)

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
		for {
			_, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			if tok == token.COMMENT {
				znaki += len(lit)
			}
		}
		wiersze := bytes.Count(tresc, []byte("\n")) + 1
		granica := wiersze * 250 / 1000
		stan := "W-GRANICY"
		if znaki > granica {
			stan = "PONAD"
		}
		fmt.Printf("%s\twiersze=%d\tznaki=%d\tgranica=%d\t%s\n", os.Args[1], wiersze, znaki, granica, stan)
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
