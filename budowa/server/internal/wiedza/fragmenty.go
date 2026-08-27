// Podział dokumentu na fragmenty, czyli na jednostki wchodzące do wskaźnika
// i wracające w odpowiedzi jako cytat. Wektor liczony jest z fragmentu, więc
// granica fragmentu jest granicą znaczenia. Jednostką podziału jest akapit,
// a w jego braku zdanie.
package wiedza

import "strings"

const (
	// minimalnaDlugoscFragmentu — poniżej tej granicy fragment przestaje nieść
	// kontekst i wektory zaczynają się zlewać.
	minimalnaDlugoscFragmentu = 200
	// granicaFragmentu — twardy sufit długości przytoczenia, licznik znaków
	// tekstu cytatu wracającego Operatorowi.
	granicaFragmentu = 1100
)

// Fragment to jedna jednostka wskaźnika: kawałek treści wraz z jego miejscem
// w dokumencie źródłowym, licznikiem kolejności i tekstem cytatu.
type Fragment struct {
	// Kolejnosc — numer fragmentu w dokumencie, liczony od zera.
	Kolejnosc int
	// Tresc — sam tekst fragmentu, ten, który wróci Operatorowi jako cytat.
	Tresc string
}

// Podziel dzieli treść dokumentu na fragmenty według reguły z nagłówka.
//
// Treść pusta daje wykaz pusty, a nie jeden fragment pusty: dokument bez tekstu
// nie ma czego wnieść do wskaźnika, a wektor z pustki byłby trafieniem, które
// pasuje do wszystkiego.
func Podziel(tresc string, dlugosc int) []Fragment {
	if dlugosc < minimalnaDlugoscFragmentu {
		dlugosc = dlugoscFragmentuDomyslna
	}
	if dlugosc > granicaFragmentu {
		dlugosc = granicaFragmentu
	}

	czesci := naCzesci(tresc, dlugosc)
	fragmenty := make([]Fragment, 0, len(czesci))
	poprzednieZdanie := ""
	for _, czesc := range czesci {
		pelna := czesc
		if poprzednieZdanie != "" {
			pelna = poprzednieZdanie + " " + czesc
		}
		fragmenty = append(fragmenty, Fragment{Kolejnosc: len(fragmenty), Tresc: pelna})
		poprzednieZdanie = ostatnieZdanie(czesc, dlugosc/3)
	}
	return fragmenty
}

// naCzesci składa akapity i zdania w kawałki nieprzekraczające granicy
// długości fragmentu wskaźnika znaczenia.
func naCzesci(tresc string, dlugosc int) []string {
	czesci := []string{}
	biezacy := strings.Builder{}

	domknij := func() {
		if tekst := strings.TrimSpace(biezacy.String()); tekst != "" {
			czesci = append(czesci, tekst)
		}
		biezacy.Reset()
	}
	dolacz := func(kawalek string) {
		if biezacy.Len()+len(kawalek) > dlugosc && biezacy.Len() > 0 {
			domknij()
		}
		if biezacy.Len() > 0 {
			biezacy.WriteString(" ")
		}
		biezacy.WriteString(kawalek)
	}

	for _, akapit := range akapity(tresc) {
		if len(akapit) <= dlugosc {
			dolacz(akapit)
			continue
		}
		// Akapit dłuższy od granicy rozkłada się na zdania.
		for _, zdanie := range zdania(akapit) {
			for len(zdanie) > dlugosc {
				ciecie := ostatniaSpacja(zdanie, dlugosc)
				dolacz(strings.TrimSpace(zdanie[:ciecie]))
				zdanie = strings.TrimSpace(zdanie[ciecie:])
			}
			if zdanie != "" {
				dolacz(zdanie)
			}
		}
	}
	domknij()
	return czesci
}

// akapity rozdziela tekst pustym wierszem i odrzuca wiersze puste,
// zwracając wykaz oczyszczony z bieli.
func akapity(tresc string) []string {
	surowe := strings.Split(strings.ReplaceAll(tresc, "\r\n", "\n"), "\n\n")
	wynik := make([]string, 0, len(surowe))
	for _, akapit := range surowe {
		if oczyszczony := strings.TrimSpace(akapit); oczyszczony != "" {
			wynik = append(wynik, oczyszczony)
		}
	}
	return wynik
}

// zdania dzieli akapit na zdania po znakach kończących wypowiedzenie. Znak
// kończący zostaje przy zdaniu, bo bez kropki cytat wyglądałby na urwany.
func zdania(akapit string) []string {
	wynik := []string{}
	poczatek := 0
	for i, znak := range akapit {
		if znak != '.' && znak != '!' && znak != '?' && znak != '\n' {
			continue
		}
		if zdanie := strings.TrimSpace(akapit[poczatek : i+1]); zdanie != "" {
			wynik = append(wynik, zdanie)
		}
		poczatek = i + 1
	}
	if ogon := strings.TrimSpace(akapit[poczatek:]); ogon != "" {
		wynik = append(wynik, ogon)
	}
	return wynik
}

// ostatnieZdanie oddaje końcówkę kawałka przeznaczoną na zakładkę. Zdanie
// dłuższe niż `granica` zostaje przycięte od przodu, żeby zakładka nie
// wypełniła sobą fragmentu.
func ostatnieZdanie(czesc string, granica int) string {
	lista := zdania(czesc)
	if len(lista) == 0 {
		return ""
	}
	ostatnie := lista[len(lista)-1]
	if len(ostatnie) <= granica {
		return ostatnie
	}
	poczatek := len(ostatnie) - granica
	if miejsce := strings.IndexByte(ostatnie[poczatek:], ' '); miejsce >= 0 {
		return strings.TrimSpace(ostatnie[poczatek+miejsce:])
	}
	for poczatek < len(ostatnie) && ostatnie[poczatek]&0xC0 == 0x80 {
		poczatek++
	}
	return strings.TrimSpace(ostatnie[poczatek:])
}

// ostatniaSpacja wskazuje miejsce cięcia: ostatnia spacja przed granicą, a gdy
// jej nie ma — sama granica przesunięta do najbliższego początku znaku UTF-8,
// żeby cięcie nie rozłupało litery na bajty.
func ostatniaSpacja(tekst string, granica int) int {
	if granica >= len(tekst) {
		return len(tekst)
	}
	if miejsce := strings.LastIndex(tekst[:granica], " "); miejsce > 0 {
		return miejsce
	}
	for granica > 0 && tekst[granica]&0xC0 == 0x80 {
		granica--
	}
	return granica
}
