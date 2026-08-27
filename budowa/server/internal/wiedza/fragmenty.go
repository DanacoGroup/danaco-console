// Podział dokumentu na fragmenty, czyli na jednostki wchodzące do wskaźnika
// i wracające w odpowiedzi jako cytat. Wektor liczony jest z fragmentu, więc
// granica fragmentu jest granicą znaczenia.
//
// Jednostką podziału jest akapit, a w jego braku zdanie: akapit to myśl
// wydzielona przez autora tekstu, natomiast cięcie co N znaków rozłupuje zdanie
// w połowie słowa i daje cytat bezużyteczny. Tekst rozpada się najpierw na
// akapity (pusty wiersz); akapit krótszy niż docelowa długość doklejany jest do
// sąsiada, dopóki mieści się w granicy, bo akapit jednozdaniowy osadzony osobno
// niesie za mało kontekstu, żeby dać się odróżnić od innego jednozdaniowego.
// Akapit dłuższy od granicy rozpada się na zdania (kropka, wykrzyknik, pytajnik,
// koniec wiersza), a zdanie dłuższe od granicy cięte jest po ostatniej spacji
// przed granicą, nigdy w środku słowa.
//
// Każdy fragment poza pierwszym zaczyna się od ostatniego zdania fragmentu
// poprzedniego. Bez tej zakładki zdanie stojące na styku dwóch fragmentów traci
// połowę kontekstu po każdej stronie granicy; koszt zakładki to około jednej
// piątej więcej wektorów.
//
// Granica 1100 znaków wynika z okna modelu, przy którym te liczby powstały:
// `paraphrase-multilingual-mpnet-base-v2` obcinał wejście na 384 tokenach
// podziału XLM-R, a polszczyzna kosztuje w nim około trzech znaków na token.
// Fragment dłuższy niż okno zostaje obcięty po cichu, a cytat byłby wtedy
// dłuższy niż to, co model przeczytał. Model domyślny jest dziś inny — `bge-m3`
// przyjmuje 8192 tokeny (`max_position_embeddings` w opisie jego wag) — więc
// okno przestało być ciasne i granica przestała być jego odwzorowaniem.
// Zostaje jednak nietknięta, bo jest zarazem granicą cytatu: fragment jest tym,
// co wraca Operatorowi i modelowi jako przytoczenie, a przytoczenie na kilka
// tysięcy znaków przestaje być przytoczeniem. Docelowe 700 znaków zostawia pod
// tą granicą zapas na zakładkę i na słowa łamane na kilka tokenów.
package wiedza

import "strings"

const (
	// minimalnaDlugoscFragmentu — poniżej tej granicy fragment przestaje nieść
	// kontekst i wektory zaczynają się zlewać.
	minimalnaDlugoscFragmentu = 200
	// granicaFragmentu — twardy sufit długości przytoczenia (patrz nagłówek).
	granicaFragmentu = 1100
)

// Fragment to jedna jednostka wskaźnika: kawałek treści wraz z jego miejscem
// w dokumencie źródłowym.
type Fragment struct {
	// Kolejnosc — numer fragmentu w dokumencie, liczony od zera. Wchodzi do
	// tożsamości wiersza wskaźnika, żeby powtórne indeksowanie tego samego
	// dokumentu nadpisywało fragmenty, a nie dokładało ich drugi komplet.
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

// naCzesci składa akapity i zdania w kawałki nieprzekraczające granicy.
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
		// Akapit dłuższy od granicy: rozkładamy go na zdania, a zdanie dłuższe
		// od granicy — na kawałki cięte po ostatniej spacji.
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

// akapity rozdziela tekst pustym wierszem i odrzuca wiersze puste.
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

// zdania dzieli akapit na zdania po znakach kończących wypowiedzenie.
//
// Znak kończący zostaje przy zdaniu, bo bez kropki cytat wyglądałby na urwany.
// Skróty pisane z kropką („ul.", „art.") rozdzielą tu zdanie w miejscu, które
// zdaniem nie jest; skutek jest kosmetyczny — fragment o jedno zdanie krótszy —
// a obroną byłby dopiero słownik skrótów polszczyzny utrzymywany w rdzeniu.
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

// ostatnieZdanie oddaje końcówkę kawałka przeznaczoną na zakładkę.
// Zdanie dłuższe niż `granica` zostaje przycięte od przodu, żeby zakładka nie
// wypełniła sobą fragmentu. Cięcie pada na pierwszej spacji za granicą liczoną
// od końca, a gdy spacji tam nie ma — na najbliższym początku znaku UTF-8, więc
// zakładka nie zaczyna się od rozłupanej litery.
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
