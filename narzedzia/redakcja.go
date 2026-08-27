//go:build ignore

// Narzędzie mierzy zgodność plików z wymaganiami redakcyjnymi budowy: długość
// nagłówków, długość komentarzy w treści kodu, obecność wyrażeń zakazanych oraz
// niezmienność kodu po redakcji komentarza.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	naglowekDolem = 100
	naglowekGora  = 250
	punktGora     = 100
)

var wyrazeniaZakazane = []struct {
	nazwa  string
	wzorem *regexp.Regexp
}{
	{"odeslanie", regexp.MustCompile(`(?i)\S*\.md\b|(^|[\s(])(docs|prowadzenie)/|\b(patrz|zob\.|uzasadnienie:|szczeg[oó]{1}[lł]y w|opisan[eo] w|wi[eę]cej w)\b`)},
	{"oznaczenie", regexp.MustCompile(`\b[A-Z]{1,5}-\d{1,3}\b|(?i)\b(pozycja|dopowiedzenie)\s+\w*\d`)},
	{"odwolanie-osobowe", regexp.MustCompile(`(?i)\b(w[lł]a[sś]cicie|prowadz[aą]c)\w*`)},
	{"skrot", regexp.MustCompile(`(?i)(^|\s)(np\.|tj\.|itp\.|itd\.|m\.in\.|tzn\.|tzw\.|ok\.|ww\.)`)},
}

type komentarz struct {
	tresc    string
	naglowek bool
	wiersz   int
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "użycie: redakcja -budowa|-jezyk|-tokeny <plik>")
		os.Exit(2)
	}
	tryb, sciezka := os.Args[1], os.Args[2]
	tresc, err := os.ReadFile(sciezka)
	if err != nil {
		fmt.Fprintln(os.Stderr, "odczyt nie przeszedł:", err)
		os.Exit(1)
	}
	switch tryb {
	case "-tokeny":
		wypiszTokeny(sciezka, tresc)
	case "-budowa":
		wypiszBudowe(sciezka, tresc)
	case "-jezyk":
		wypiszJezyk(sciezka, tresc)
	case "-bezkomentarzy":
		wypiszBezKomentarzy(sciezka, tresc)
	default:
		fmt.Fprintln(os.Stderr, "nieznany tryb:", tryb)
		os.Exit(2)
	}
}

// wypiszTokeny wypisuje strumień jednostek składniowych pliku Go z pominięciem
// komentarzy. Porównanie wyniku przed redakcją i po niej rozstrzyga, czy kod
// pozostał nietknięty, ponieważ strumień nie zależy od odstępów ani od wierszy.
func wypiszTokeny(sciezka string, tresc []byte) {
	fset := token.NewFileSet()
	plik := fset.AddFile(sciezka, fset.Base(), len(tresc))
	var s scanner.Scanner
	s.Init(plik, tresc, func(token.Position, string) {}, 0)
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			return
		}
		fmt.Printf("%s\t%s\n", tok, lit)
	}
}

// wypiszBezKomentarzy podaje treść pliku pozbawioną komentarzy i pustych wierszy.
// Zgodność wyniku przed redakcją i po niej dowodzi, że zmieniono wyłącznie
// komentarz; działa w każdym języku, także tam, gdzie rozbiór składniowy nie sięga.
func wypiszBezKomentarzy(sciezka string, tresc []byte) {
	znaczniki := znacznikiJezyka(filepath.Ext(sciezka))
	for _, w := range strings.Split(string(tresc), "\n") {
		przyciety := strings.TrimSpace(w)
		if przyciety == "" {
			continue
		}
		komentarzem := false
		for _, z := range znaczniki {
			if strings.HasPrefix(przyciety, z) {
				komentarzem = true
				break
			}
		}
		if !komentarzem {
			fmt.Println(przyciety)
		}
	}
}

// wypiszBudowe podaje liczbę nagłówków krótszych i dłuższych od wymaganych oraz
// liczbę komentarzy w treści kodu przekraczających długość punktu. Dokument
// tekstowy zwraca poza zasięgiem, ponieważ miara dotyczy komentarza, nie prozy.
func wypiszBudowe(sciezka string, tresc []byte) {
	if dokument(sciezka) {
		fmt.Printf("%s\tnaglowki=0\tkrotkie=0\tdlugie=0\tpunkty=0\tPOZA-ZASIEGIEM\n", sciezka)
		return
	}
	komentarze := zbierz(sciezka, tresc)
	var krotkie, dlugie, punkty, naglowki int
	for _, k := range komentarze {
		dlugosc := utf8.RuneCountInString(k.tresc)
		if k.naglowek {
			naglowki++
			if dlugosc < naglowekDolem {
				krotkie++
			} else if dlugosc > naglowekGora {
				dlugie++
			}
			continue
		}
		if dlugosc > punktGora {
			punkty++
		}
	}
	stan := "ZGODNY"
	if naglowki == 0 || krotkie > 0 || dlugie > 0 || punkty > 0 {
		stan = "NIEZGODNY"
	}
	fmt.Printf("%s\tnaglowki=%d\tkrotkie=%d\tdlugie=%d\tpunkty=%d\t%s\n",
		sciezka, naglowki, krotkie, dlugie, punkty, stan)
}

// wypiszJezyk wykazuje wyrażenia zakazane w komentarzach plików kodu oraz w całej
// treści dokumentów tekstowych. Dokumentowi nie liczy odesłań, ponieważ wskazanie
// dokumentu źródłowego jest tam treścią, a nie odesłaniem zamiast treści.
func wypiszJezyk(sciezka string, tresc []byte) {
	var doZbadania []komentarz
	if strings.EqualFold(filepath.Ext(sciezka), ".md") {
		for i, w := range strings.Split(string(tresc), "\n") {
			doZbadania = append(doZbadania, komentarz{tresc: w, wiersz: i + 1})
		}
	} else {
		doZbadania = zbierz(sciezka, tresc)
	}
	naruszen := 0
	for _, k := range doZbadania {
		for _, z := range wyrazeniaZakazane {
			if z.nazwa == "odeslanie" && dokument(sciezka) {
				continue
			}
			if m := z.wzorem.FindString(k.tresc); m != "" {
				naruszen++
				fmt.Printf("%s\t%d\t%s\t%s\n", sciezka, k.wiersz, z.nazwa, strings.TrimSpace(m))
			}
		}
	}
	if naruszen == 0 {
		fmt.Printf("%s\t0\tzgodny\t\n", sciezka)
	}
}

// zbierz wydobywa komentarze pliku i rozstrzyga, które z nich są nagłówkami.
// Pliki Go rozbiera składniowo, ponieważ przypisanie komentarza do deklaracji
// jest tam jednoznaczne; pozostałe języki rozstrzyga położenie komentarza.
func zbierz(sciezka string, tresc []byte) []komentarz {
	if strings.EqualFold(filepath.Ext(sciezka), ".go") {
		return zbierzGo(sciezka, tresc)
	}
	return zbierzOgolnie(sciezka, tresc)
}

func zbierzGo(sciezka string, tresc []byte) []komentarz {
	fset := token.NewFileSet()
	plik, err := parser.ParseFile(fset, sciezka, tresc, parser.ParseComments)
	if err != nil {
		return zbierzOgolnie(sciezka, tresc)
	}
	czyNaglowek := map[token.Pos]bool{}
	if plik.Doc != nil {
		czyNaglowek[plik.Doc.Pos()] = true
	}
	for _, d := range plik.Decls {
		switch v := d.(type) {
		case *ast.GenDecl:
			if v.Doc != nil {
				czyNaglowek[v.Doc.Pos()] = true
			}
			for _, s := range v.Specs {
				switch t := s.(type) {
				case *ast.TypeSpec:
					if t.Doc != nil {
						czyNaglowek[t.Doc.Pos()] = true
					}
				case *ast.ValueSpec:
					if t.Doc != nil {
						czyNaglowek[t.Doc.Pos()] = true
					}
				}
			}
		case *ast.FuncDecl:
			if v.Doc != nil {
				czyNaglowek[v.Doc.Pos()] = true
			}
		}
	}
	var wynik []komentarz
	for _, g := range plik.Comments {
		if strings.HasPrefix(g.List[0].Text, "//go:") {
			continue
		}
		wynik = append(wynik, komentarz{
			tresc:    oczysc(g.Text()),
			naglowek: czyNaglowek[g.Pos()],
			wiersz:   fset.Position(g.Pos()).Line,
		})
	}
	return wynik
}

func zbierzOgolnie(sciezka string, tresc []byte) []komentarz {
	wiersze := strings.Split(string(tresc), "\n")
	znaczniki := znacznikiJezyka(filepath.Ext(sciezka))
	var wynik []komentarz
	var blok []string
	poczatek := 0
	domknij := func(nastepny string) {
		if len(blok) == 0 {
			return
		}
		naglowek := poczatek == 1 || (nastepny != "" && !strings.HasPrefix(nastepny, " ") && !strings.HasPrefix(nastepny, "\t"))
		wynik = append(wynik, komentarz{tresc: oczysc(strings.Join(blok, " ")), naglowek: naglowek, wiersz: poczatek})
		blok = nil
	}
	for i, w := range wiersze {
		przyciety := strings.TrimSpace(w)
		czyKomentarz := false
		for _, z := range znaczniki {
			if strings.HasPrefix(przyciety, z) {
				czyKomentarz = true
				przyciety = strings.TrimPrefix(przyciety, z)
				break
			}
		}
		if czyKomentarz {
			if len(blok) == 0 {
				poczatek = i + 1
			}
			blok = append(blok, przyciety)
			continue
		}
		domknij(w)
	}
	domknij("")
	return wynik
}

// dokument rozstrzyga, czy plik jest opracowaniem tekstowym. Miara nagłówka
// i punktu dotyczy komentarza w kodzie, a w dokumencie znak wyliczenia wypadałby
// za komentarz, więc opracowanie ocenia się odczytem i miarą języka.
func dokument(sciezka string) bool {
	return strings.EqualFold(filepath.Ext(sciezka), ".md")
}

func znacznikiJezyka(rozszerzenie string) []string {
	switch strings.ToLower(rozszerzenie) {
	case ".py", ".sh", ".yml", ".yaml", ".toml":
		return []string{"#"}
	case ".sql":
		return []string{"--"}
	case ".html", ".htm", ".xml", ".svg":
		return []string{"<!--"}
	case ".css":
		return []string{"/*", "*"}
	default:
		return []string{"//", "/*", "*"}
	}
}

// oczysc sprowadza treść komentarza do samych słów: zdejmuje znaczniki końca
// komentarza, nadmiarowe odstępy i znaki ozdobne, aby długość mierzyła treść,
// a nie sposób jej obramowania.
func oczysc(s string) string {
	s = strings.NewReplacer("*/", " ", "-->", " ", "/*", " ", "<!--", " ").Replace(s)
	s = strings.Trim(s, " \t\n─═—-*#")
	return strings.Join(strings.Fields(s), " ")
}
