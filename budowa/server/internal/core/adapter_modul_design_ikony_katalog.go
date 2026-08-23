// Odpowiedzialność pliku: katalog ikon WKOMPILOWANY w binarium rdzenia —
// wzory, których `design.icon.library.search` szuka, `design.icon.set` używa
// jako punktu wyjścia, a `design.icon.generate` bierze, gdy pojęcie trafia
// w któryś z nich. Ikony WŁASNE Operatora leżą w bazie
// (`dane/design_ikony.go`); tutaj leży to, co produkt ma bez ani jednego zapisu
// w bazie i bez ani jednego pobrania z sieci.
//
// ── Dlaczego katalog jest w kodzie, a nie w bazie ani w pliku ───────────────
// Instalka Operatora jest cienka, a arsenał stoi na serwerze wkompilowany
// w binarium. Katalog ikon w bazie musiałby być zasiany migracją i dałby się
// skasować; katalog w plikach obok binarium wymagałby, żeby wdrożenie kopiowało
// katalog danych. Katalog w kodzie jest zawsze i jest jeden.
//
// ── Jeden styl, jedna siatka ────────────────────────────────────────────────
// Wszystkie wzory są na siatce 24 i rysowane KONTUREM (obrys, nie wypełnienie),
// o grubości dwóch jednostek. Zestaw mieszający obrys z wypełnieniem wygląda
// w oknie jak zestaw zebrany z dwóch źródeł — a to jest dokładnie to, czemu
// zestaw ikon ma zapobiegać. Grubość obrysu jest cechą ikony
// (`strokeWidth`), więc Operator zmienia ją i dostaje tę samą ikonę cieńszą,
// a nie inną ikonę.
//
// ── Zapis jest ścieżką, nie gotowym dokumentem ──────────────────────────────
// Wzór niesie samą treść `d` ścieżek. Dokument SVG składa się z niej na wyjściu
// (`svgIkonyKataloguDesignu`) razem z grubością obrysu, którą Operator wskazał —
// gdyby wzory były gotowymi dokumentami, zmiana grubości wymagałaby przepisania
// czterdziestu napisów.
package core

import (
	"fmt"
	"sort"
	"strings"

	"danacoconsole/shared"
)

const (
	// siatkaIkonyKataloguDesignu jest siatką, na której powstały wszystkie wzory.
	// Ikona żądana na innej siatce skaluje się z tej — jedno źródło kształtu.
	siatkaIkonyKataloguDesignu = 24

	// gruboscObrysuIkonyDomyslna jest grubością obrysu wzorów katalogu.
	gruboscObrysuIkonyDomyslna = 2.0

	// zestawIkonKataloguDesignu jest nazwą zestawu wkompilowanego. Nazwa wchodzi
	// do pola `set` każdej ikony katalogu, żeby Operator poznał, że ma przed sobą
	// wzór rdzenia, a nie własną ikonę.
	zestawIkonKataloguDesignu = "rdzen-24"
)

// wzorIkonyDesignu to jeden wzór katalogu wkompilowanego.
type wzorIkonyDesignu struct {
	Nazwa    string
	Etykiety []string
	// Sciezki niosą treść atrybutu `d` — po jednej na podścieżkę rysunku.
	Sciezki []string
}

// katalogIkonDesignu to komplet wzorów wkompilowanych w binarium.
//
// Wykaz jest krótki z zamysłu i obejmuje pojęcia, które w narzędziu pracy
// naprawdę występują: nawigacja, praca z plikami, stany, komunikacja, media,
// dane. Wciągnięcie kilku tysięcy ikon dałoby nazwy, których nikt nie wpisze,
// a każda musiałaby być tu utrzymywana.
var katalogIkonDesignu = []wzorIkonyDesignu{
	{"strzalka-w-gore", []string{"strzałka", "góra", "wyżej", "nawigacja"},
		[]string{"M12 20 L12 4", "M5 11 L12 4 L19 11"}},
	{"strzalka-w-dol", []string{"strzałka", "dół", "niżej", "nawigacja"},
		[]string{"M12 4 L12 20", "M5 13 L12 20 L19 13"}},
	{"strzalka-w-lewo", []string{"strzałka", "lewo", "wstecz", "nawigacja"},
		[]string{"M20 12 L4 12", "M11 5 L4 12 L11 19"}},
	{"strzalka-w-prawo", []string{"strzałka", "prawo", "dalej", "nawigacja"},
		[]string{"M4 12 L20 12", "M13 5 L20 12 L13 19"}},
	{"plus", []string{"plus", "dodaj", "nowy", "załóż"},
		[]string{"M12 5 L12 19", "M5 12 L19 12"}},
	{"minus", []string{"minus", "odejmij", "zmniejsz"},
		[]string{"M5 12 L19 12"}},
	{"krzyzyk", []string{"zamknij", "anuluj", "usuń", "krzyżyk"},
		[]string{"M6 6 L18 18", "M18 6 L6 18"}},
	{"ptaszek", []string{"gotowe", "zatwierdź", "tak", "ptaszek"},
		[]string{"M4 13 L9 18 L20 6"}},
	{"kolo", []string{"koło", "okrąg", "kropka", "stan"},
		[]string{"M12 3 A9 9 0 1 0 12 21 A9 9 0 1 0 12 3"}},
	{"kwadrat", []string{"kwadrat", "prostokąt", "ramka", "stop"},
		[]string{"M4 4 L20 4 L20 20 L4 20 Z"}},
	{"trojkat", []string{"trójkąt", "ostrzeżenie", "uwaga"},
		[]string{"M12 4 L21 20 L3 20 Z"}},
	{"gwiazda", []string{"gwiazda", "ulubione", "ocena"},
		[]string{"M12 3 L15 9.5 L22 10.4 L17 15.3 L18.2 22 L12 18.7 L5.8 22 L7 15.3 " +
			"L2 10.4 L9 9.5 Z"}},
	{"serce", []string{"serce", "polubienie", "ulubione"},
		[]string{"M12 20 L5 13 A4 4 0 0 1 12 8 A4 4 0 0 1 19 13 Z"}},
	{"dom", []string{"dom", "start", "początek", "nawigacja"},
		[]string{"M3 11 L12 4 L21 11", "M5 11 L5 20 L19 20 L19 11"}},
	{"lupa", []string{"szukaj", "wyszukiwanie", "lupa", "filtr"},
		[]string{"M10 4 A6 6 0 1 0 10 16 A6 6 0 1 0 10 4", "M14.5 14.5 L20 20"}},
	{"kolo-zebate", []string{"ustawienia", "nastawy", "konfiguracja", "koło zębate"},
		[]string{"M12 8 A4 4 0 1 0 12 16 A4 4 0 1 0 12 8",
			"M12 2 L12 5", "M12 19 L12 22", "M2 12 L5 12", "M19 12 L22 12",
			"M5 5 L7 7", "M17 17 L19 19", "M19 5 L17 7", "M7 17 L5 19"}},
	{"uzytkownik", []string{"użytkownik", "osoba", "konto", "profil"},
		[]string{"M12 4 A4 4 0 1 0 12 12 A4 4 0 1 0 12 4", "M4 21 A8 8 0 0 1 20 21"}},
	{"ludzie", []string{"ludzie", "zespół", "współpraca", "grupa"},
		[]string{"M9 4 A3.5 3.5 0 1 0 9 11 A3.5 3.5 0 1 0 9 4",
			"M2 20 A7 7 0 0 1 16 20", "M16 5 A3 3 0 0 1 16 11", "M17 13 A7 7 0 0 1 22 20"}},
	{"koperta", []string{"poczta", "wiadomość", "koperta", "e-mail"},
		[]string{"M3 6 L21 6 L21 18 L3 18 Z", "M3 6 L12 13 L21 6"}},
	{"dzwonek", []string{"powiadomienie", "alarm", "dzwonek"},
		[]string{"M12 3 A6 6 0 0 1 18 9 L18 15 L20 18 L4 18 L6 15 L6 9 A6 6 0 0 1 12 3 Z",
			"M10 18 A2 2 0 0 0 14 18"}},
	{"rozmowa", []string{"rozmowa", "czat", "komentarz", "dymek"},
		[]string{"M4 5 L20 5 L20 16 L12 16 L7 20 L7 16 L4 16 Z"}},
	{"dokument", []string{"plik", "dokument", "kartka"},
		[]string{"M6 3 L14 3 L19 8 L19 21 L6 21 Z", "M14 3 L14 8 L19 8"}},
	{"folder", []string{"folder", "katalog", "teczka"},
		[]string{"M3 6 L10 6 L12 9 L21 9 L21 19 L3 19 Z"}},
	{"pobierz", []string{"pobierz", "zapisz", "wydanie", "eksport"},
		[]string{"M12 3 L12 15", "M6 10 L12 16 L18 10", "M4 20 L20 20"}},
	{"wgraj", []string{"wgraj", "wnieś", "import", "wczytaj"},
		[]string{"M12 17 L12 5", "M6 11 L12 5 L18 11", "M4 20 L20 20"}},
	{"kosz", []string{"kosz", "usuń", "wyrzuć"},
		[]string{"M4 7 L20 7", "M9 7 L9 4 L15 4 L15 7", "M6 7 L7 21 L17 21 L18 7",
			"M10 11 L10 18", "M14 11 L14 18"}},
	{"olowek", []string{"edytuj", "ołówek", "popraw", "zmień"},
		[]string{"M4 20 L4 16 L16 4 L20 8 L8 20 Z", "M14 6 L18 10"}},
	{"kopiuj", []string{"kopiuj", "duplikat", "powtórz"},
		[]string{"M8 3 L18 3 L18 15 L8 15 Z", "M6 8 L6 21 L16 21"}},
	{"klucz", []string{"klucz", "hasło", "poświadczenie", "dostęp"},
		[]string{"M8 8 A4 4 0 1 0 8 16 A4 4 0 1 0 8 8", "M11.5 12 L21 12",
			"M17 12 L17 16", "M20 12 L20 15"}},
	{"zamek", []string{"zamek", "blokada", "bezpieczeństwo", "prywatne"},
		[]string{"M5 11 L19 11 L19 21 L5 21 Z", "M8 11 L8 7 A4 4 0 0 1 16 7 L16 11"}},
	{"oko", []string{"podgląd", "oko", "widok", "pokaż"},
		[]string{"M2 12 A12 8 0 0 1 22 12 A12 8 0 0 1 2 12 Z",
			"M12 9 A3 3 0 1 0 12 15 A3 3 0 1 0 12 9"}},
	{"zegar", []string{"zegar", "czas", "historia", "harmonogram"},
		[]string{"M12 3 A9 9 0 1 0 12 21 A9 9 0 1 0 12 3", "M12 7 L12 12 L16 14"}},
	{"kalendarz", []string{"kalendarz", "data", "termin"},
		[]string{"M4 6 L20 6 L20 20 L4 20 Z", "M4 10 L20 10", "M8 3 L8 6", "M16 3 L16 6"}},
	{"wykres-slupkowy", []string{"wykres", "słupki", "dane", "statystyki"},
		[]string{"M4 20 L4 12", "M10 20 L10 6", "M16 20 L16 10", "M3 21 L21 21"}},
	{"wykres-liniowy", []string{"wykres", "linia", "trend", "wzrost"},
		[]string{"M3 17 L9 11 L13 14 L21 5", "M3 21 L21 21"}},
	{"obraz", []string{"obraz", "grafika", "zdjęcie", "zasób"},
		[]string{"M3 5 L21 5 L21 19 L3 19 Z", "M3 16 L9 10 L14 15 L17 12 L21 16",
			"M15.5 8.5 A1.5 1.5 0 1 0 15.5 8.6"}},
	{"aparat", []string{"aparat", "fotografia", "zdjęcie"},
		[]string{"M3 8 L7 8 L9 5 L15 5 L17 8 L21 8 L21 19 L3 19 Z",
			"M12 9.5 A4 4 0 1 0 12 17.5 A4 4 0 1 0 12 9.5"}},
	{"pedzel", []string{"pędzel", "malowanie", "retusz", "warsztat"},
		[]string{"M6 21 A4 4 0 0 0 10 17 L10 14 L14 14 L14 17 A4 4 0 0 1 10 21 Z",
			"M10 14 L20 4 L22 6 L14 14"}},
	{"kropla", []string{"kropla", "barwa", "próbnik", "nasycenie"},
		[]string{"M12 3 L18 11 A6 6 0 1 1 6 11 Z"}},
	{"paleta", []string{"paleta", "barwy", "kolor", "harmonia"},
		[]string{"M12 3 A9 9 0 1 0 12 21 A3 3 0 0 0 12 15 A3 3 0 0 1 12 9 A9 9 0 0 0 12 3 Z",
			"M8 8 A1 1 0 1 0 8 8.1", "M7 13 A1 1 0 1 0 7 13.1", "M11 6 A1 1 0 1 0 11 6.1"}},
	{"wektor", []string{"wektor", "węzły", "ścieżka", "pióro"},
		[]string{"M4 20 C4 10 14 4 20 4", "M2 18 L6 18 L6 22 L2 22 Z",
			"M18 2 L22 2 L22 6 L18 6 Z"}},
	{"drukarka", []string{"druk", "drukarka", "wydruk", "nośnik"},
		[]string{"M7 8 L7 3 L17 3 L17 8", "M3 8 L21 8 L21 16 L17 16 L17 21 L7 21 L7 16 L3 16 Z"}},
	{"warstwy", []string{"warstwy", "kompozycja", "stos"},
		[]string{"M12 3 L21 8 L12 13 L3 8 Z", "M3 12 L12 17 L21 12", "M3 16 L12 21 L21 16"}},
	{"siatka", []string{"siatka", "układ", "kolumny", "kafle"},
		[]string{"M4 4 L20 4 L20 20 L4 20 Z", "M12 4 L12 20", "M4 12 L20 12"}},
	{"odswiez", []string{"odśwież", "ponów", "regeneruj", "obrót"},
		[]string{"M20 12 A8 8 0 1 1 12 4", "M12 1 L12 7 L18 7"}},
	{"informacja", []string{"informacja", "pomoc", "opis"},
		[]string{"M12 3 A9 9 0 1 0 12 21 A9 9 0 1 0 12 3", "M12 11 L12 16",
			"M12 7.5 A0.8 0.8 0 1 0 12 7.6"}},
}

// ikonaKataloguDesignu odnajduje wzór po nazwie.
func ikonaKataloguDesignu(nazwa string) (wzorIkonyDesignu, bool) {
	szukana := strings.ToLower(strings.TrimSpace(nazwa))
	for _, wzor := range katalogIkonDesignu {
		if strings.ToLower(wzor.Nazwa) == szukana {
			return wzor, true
		}
	}
	return wzorIkonyDesignu{}, false
}

// wzorDlaPojeciaDesignu szuka wzoru odpowiadającego pojęciu — po nazwie i po
// etykietach. Dopasowanie jest po zawieraniu w obie strony, bo Operator pisze
// „strzałka w prawo" i „prawo", mając na myśli to samo.
func wzorDlaPojeciaDesignu(pojecie string) (wzorIkonyDesignu, bool) {
	szukane := strings.ToLower(strings.TrimSpace(pojecie))
	if szukane == "" {
		return wzorIkonyDesignu{}, false
	}
	if wzor, jest := ikonaKataloguDesignu(szukane); jest {
		return wzor, true
	}
	for _, wzor := range katalogIkonDesignu {
		if strings.Contains(strings.ToLower(wzor.Nazwa), szukane) {
			return wzor, true
		}
		for _, etykieta := range wzor.Etykiety {
			if strings.Contains(szukane, strings.ToLower(etykieta)) ||
				strings.Contains(strings.ToLower(etykieta), szukane) {
				return wzor, true
			}
		}
	}
	return wzorIkonyDesignu{}, false
}

// zestawyIkonDesignu oddaje zestawy obecne w katalogu rdzenia. Wykaz jest
// jednoelementowy i taki wchodzi do odpowiedzi — pusty byłby ciszą o tym, że
// rdzeń w ogóle ma ikony.
func zestawyIkonDesignu() []string {
	return []string{zestawIkonKataloguDesignu}
}

// svgIkonyKataloguDesignu składa dokument SVG wzoru na żądanej siatce
// i z żądaną grubością obrysu.
//
// Skalowanie idzie polem `viewBox`, nie przeliczaniem współrzędnych: wzór ma
// jedną prawdę o kształcie, a rozmiar dokumentu jest jego oprawą. Dzięki temu
// ikona na siatce 16 i na siatce 48 to ten sam rysunek, a nie dwa zaokrąglone
// inaczej.
func svgIkonyKataloguDesignu(wzor wzorIkonyDesignu, siatka int, grubosc float64) string {
	if siatka <= 0 {
		siatka = siatkaIkonyKataloguDesignu
	}
	if grubosc <= 0 {
		grubosc = gruboscObrysuIkonyDomyslna
	}
	var dokument strings.Builder
	fmt.Fprintf(&dokument,
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" `+
			`fill="none" stroke="currentColor" stroke-width="%g" stroke-linecap="round" `+
			`stroke-linejoin="round">`,
		siatka, siatka, siatkaIkonyKataloguDesignu, siatkaIkonyKataloguDesignu, grubosc)
	for _, sciezka := range wzor.Sciezki {
		fmt.Fprintf(&dokument, `<path d="%s"/>`, sciezka)
	}
	dokument.WriteString(`</svg>`)
	return dokument.String()
}

// ikonaKataloguKontraktuDesignu składa `DesignIcon` kontraktu ze wzoru.
func ikonaKataloguKontraktuDesignu(wzor wzorIkonyDesignu, siatka int,
	grubosc float64) shared.DesignIcon {

	zestaw := zestawIkonKataloguDesignu
	svg := svgIkonyKataloguDesignu(wzor, siatka, grubosc)
	if siatka <= 0 {
		siatka = siatkaIkonyKataloguDesignu
	}
	if grubosc <= 0 {
		grubosc = gruboscObrysuIkonyDomyslna
	}
	etykiety := make([]string, len(wzor.Etykiety))
	copy(etykiety, wzor.Etykiety)
	sort.Strings(etykiety)
	return shared.DesignIcon{
		Id: zestawIkonKataloguDesignu + "/" + wzor.Nazwa, Name: wzor.Nazwa, Set: &zestaw,
		Tags: etykiety, Svg: &svg, GridSize: &siatka, StrokeWidth: &grubosc,
	}
}

// sciezkiZDokumentuSvgDesignu wyciąga treść atrybutów `d` z dokumentu SVG.
//
// Rozbiór jest celowo wąski: rdzeń czyta ŚCIEŻKI, bo tylko z nich składa się
// kontur ikony w pakiecie i w kroju. Dokument z prostokątami i okręgami zamiast
// ścieżek oddaje wykaz pusty, a wołający nazywa to wprost — cichy pakiet
// z pustymi glifami byłby plikiem, w którym nie widać nic.
func sciezkiZDokumentuSvgDesignu(dokument string) []string {
	sciezki := []string{}
	reszta := dokument
	for {
		poczatek := strings.Index(reszta, ` d="`)
		if poczatek < 0 {
			return sciezki
		}
		reszta = reszta[poczatek+4:]
		koniec := strings.IndexByte(reszta, '"')
		if koniec < 0 {
			return sciezki
		}
		wartosc := strings.TrimSpace(reszta[:koniec])
		if wartosc != "" {
			sciezki = append(sciezki, wartosc)
		}
		reszta = reszta[koniec+1:]
	}
}
