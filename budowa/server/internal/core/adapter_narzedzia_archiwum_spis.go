// Odpowiedzialność pliku: czytanie spisu archiwum (`7z l -slt`) i wyrok o tym,
// czy wolno je rozpakować. Zaplecze rodziny leży w
// `adapter_narzedzia_archiwum.go`, samo rozpakowanie w
// `adapter_narzedzia_archiwum_rozpakowanie.go`.
//
// ── Ścieżki wychodzące poza katalog docelowy ───────────────────────────────
// Spis sprawdza się przed rozpakowaniem, a nie ścieżki po nim, i to jest wybór,
// nie skrót. Sprawdzenie po rozpakowaniu przychodzi za późno z definicji:
// pozycja `../../.ssh/authorized_keys` jest już wtedy zapisana, a wykrycie
// nadpisania nie odkręca nadpisania. Odmowa ma paść, zanim poleci pierwszy bajt.
//
// Spis czyta to samo binarium, które będzie rozpakowywać, i w tym samym
// przebiegu wywołań — oglądane jest więc dokładnie to, co `7z` z tego archiwum
// odczytuje, a nie własne wyobrażenie o formacie zip. Własny czytnik nagłówków
// zip w Go byłby drugą prawdą o zawartości archiwum, a rozjazd między czytnikiem
// sprawdzającym a rozpakowującym jest klasycznym sposobem obejścia takiej
// kontroli.
//
// `7z x` przy rozpakowaniu sam obcina człon `..` i zapisuje pozycję wewnątrz
// katalogu docelowego. Ta obrona binarium jest prawdziwa, ale nie jest tą, na
// której stoi rdzeń, z dwóch powodów:
//   - jest cudzą własnością i cudzą wersją — archiwum wychodzące poza katalog
//     ma zostać odmówione na każdej maszynie i przy każdym `7z`, a nie
//     rozpakowane inaczej, niż zapowiada;
//   - ciche obcięcie członu jest zmianą znaczenia bez powiedzenia o tym.
//     Operator dostałby drzewo inne niż to, które archiwum opisuje, i nie
//     dowiedziałby się o tym. Nazwanie braku jest uczciwsze niż naprawienie
//     archiwum po cichu.
//
// Pewność bierze się więc z trzech warstw naraz: spis odmawia przed zapisem,
// rozpakowanie idzie do kwarantanny (pustego katalogu, w którym nie ma czego
// nadpisać), a przejście po wyniku sprawdza, co naprawdę powstało — dopiero
// potem treść wchodzi do katalogu Operatora.
//
// ── Dowiązania są odmawiane, nie rozpakowywane ─────────────────────────────
// Pozycja będąca dowiązaniem symbolicznym wychodzi poza katalog docelowy inaczej
// niż członem `..`: sama w sobie jest niewinna (`link.txt`), ale wskazuje na
// zewnątrz, a kolejna pozycja archiwum zapisuje przez nią — i zapis ląduje tam,
// gdzie wskazuje dowiązanie. `7z l -slt` pokazuje na warstwie `tar` pole
// `Symbolic Link = …`, więc dowiązanie widać w spisie bez zgadywania.
package core

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"
)

// spisArchiwum pyta `7z` o zawartość archiwum i oddaje ją rozebraną na pozycje.
//
// Spis jest odczytem i niczego nie zapisuje — dlatego wolno go zrobić przed
// wyrokiem, na materiale jeszcze nieocenionym. To jest cały fundament obrony:
// gdyby poznanie zawartości archiwum wymagało jego rozpakowania, odmowa zawsze
// przychodziłaby po szkodzie.
//
// Archiwum puste w odczycie jest odmową, a nie pustym wynikiem: `7z`, który nie
// umiał odczytać spisu (archiwum uszkodzone, zaszyfrowany nagłówek), zostawia
// rdzeń bez wiedzy o zawartości — a rozpakowanie czegoś nieobejrzanego omijałoby
// wszystkie sprawdzenia niżej.
func (a *adapterNarzedziArchiwum) spisArchiwum(ctx context.Context,
	sciezka string) ([]pozycjaArchiwum, error) {

	wyjscie, err := a.wolaj7z(ctx, []string{"l", "-slt", "-bd", "-y", "--", sciezka},
		filepath.Dir(sciezka))
	if err != nil {
		return nil, err
	}
	pozycje := czytajSpisArchiwum(wyjscie)
	if len(pozycje) == 0 {
		return nil, bladWskazaniaArchiwum("nie da się odczytać spisu archiwum " +
			nazwaBezKatalogow(sciezka) + " — archiwum jest puste, uszkodzone albo ma zaszyfrowany nagłówek; " +
			"rdzeń nie rozpakowuje archiwum, którego zawartości nie obejrzał")
	}
	return pozycje, nil
}

// pozycjaArchiwum jest jedną pozycją spisu — tyle, ile potrzeba do wyroku.
type pozycjaArchiwum struct {
	// sciezka jest ścieżką pozycji wewnątrz archiwum, tak jak zapisał ją twórca.
	sciezka string
	// katalog mówi, czy pozycja jest katalogiem (`Folder = +`).
	katalog bool
	// rozmiar jest rozmiarem po rozpakowaniu, deklarowanym przez archiwum.
	rozmiar int64
	// dowiazanie niesie cel dowiązania symbolicznego albo twardego; puste
	// znaczy „pozycja jest zwykłym plikiem albo katalogiem".
	dowiazanie string
}

// czytajSpisArchiwum rozbiera wyjście `7z l -slt` na pozycje.
//
// Format `-slt` jest prosty i dlatego wybrany: po linii `----------` idą
// akapity oddzielone pustą linią, a każdy niesie pola `Klucz = wartość`.
// Postać zwykła (`7z l`) jest tabelą kolumnową z nazwą uciętą do szerokości
// kolumny — ścieżki nie da się z niej odczytać wiernie, a wyrok o ścieżce
// odczytanej niewiernie nie jest wart nic.
func czytajSpisArchiwum(wyjscie string) []pozycjaArchiwum {
	var pozycje []pozycjaArchiwum
	var biezaca pozycjaArchiwum
	var wSpisie, rozpoczeta bool

	domknij := func() {
		if rozpoczeta {
			pozycje = append(pozycje, biezaca)
		}
		biezaca = pozycjaArchiwum{}
		rozpoczeta = false
	}

	for _, linia := range strings.Split(wyjscie, "\n") {
		linia = strings.TrimRight(linia, "\r")
		// Nagłówek `7z` też niesie pole `Path` (nazwę samego archiwum), więc
		// pozycje liczymy dopiero od kreski oddzielającej nagłówek od spisu.
		if !wSpisie {
			if strings.HasPrefix(strings.TrimSpace(linia), "----------") {
				wSpisie = true
			}
			continue
		}
		if strings.TrimSpace(linia) == "" {
			domknij()
			continue
		}
		klucz, wartosc, jest := strings.Cut(linia, " = ")
		if !jest {
			continue
		}
		klucz = strings.TrimSpace(klucz)
		wartosc = strings.TrimSpace(wartosc)

		switch klucz {
		case "Path":
			// Nowe `Path` otwiera nową pozycję, także gdy poprzedniej nie
			// zamknęła pusta linia.
			domknij()
			biezaca.sciezka = wartosc
			rozpoczeta = true
		case "Folder":
			biezaca.katalog = wartosc == "+"
		case "Size":
			if liczba, err := strconv.ParseInt(wartosc, 10, 64); err == nil {
				biezaca.rozmiar = liczba
			}
		case "Symbolic Link", "Hard Link":
			if wartosc != "" {
				biezaca.dowiazanie = wartosc
			}
		case "Attributes", "Mode":
			// Tryb uniksowy zaczynający się od `l` to dowiązanie także tam,
			// gdzie osobnego pola `Symbolic Link` nie ma (zip z rozszerzeniem
			// uniksowym). Drugie źródło tej samej prawdy, bo pierwsze bywa puste.
			if strings.HasPrefix(wartosc, "l") && biezaca.dowiazanie == "" {
				biezaca.dowiazanie = "(tryb " + wartosc + ")"
			}
		}
	}
	domknij()
	return pozycje
}

// sprawdzSpisArchiwum wydaje wyrok: czy wolno rozpakować to archiwum.
//
// Kolejność sprawdzeń jest od najgroźniejszego do najbardziej ilościowego, żeby
// odmowa nazwała powód najpoważniejszy, a nie pierwszy napotkany. Archiwum
// nadpisujące cudze pliki jest czymś innym niż archiwum po prostu wielkim
// i Operator ma to przeczytać wprost.
func sprawdzSpisArchiwum(pozycje []pozycjaArchiwum) error {
	if len(pozycje) == 0 {
		return bladWskazaniaArchiwum("archiwum nie ma ani jednej pozycji — " +
			"nie ma czego rozpakować")
	}
	if len(pozycje) > granicaRozpakowaniaPozycji {
		return bladWskazaniaArchiwum("archiwum ma " + strconv.Itoa(len(pozycje)) +
			" pozycji, a granica wynosi " + strconv.Itoa(granicaRozpakowaniaPozycji) +
			" — archiwum o takiej liczbie pozycji wyczerpuje i-węzły nośnika, " +
			"nie zajmując go prawie wcale w bajtach; " +
			"naprawa: rozpakować je poza produktem albo podzielić na części")
	}

	for _, pozycja := range pozycje {
		if err := sprawdzSciezkePozycji(pozycja); err != nil {
			return err
		}
	}

	var suma int64
	for _, pozycja := range pozycje {
		if pozycja.katalog {
			continue
		}
		suma += pozycja.rozmiar
		if suma > granicaRozpakowaniaBajty {
			return bladPrzekroczonejGranicyRozmiaru(suma)
		}
	}
	return nil
}

// sprawdzSciezkePozycji jest sprawdzeniem jednej pozycji spisu.
//
// Cztery postacie ucieczki poza katalog docelowy, wszystkie odmawiane:
//   - ścieżka bezwzględna uniksowa (`/etc/cron.d/cokolwiek`),
//   - ścieżka bezwzględna windowsowa (`C:\...`) i sieciowa (`\\serwer\...`),
//   - człon `..` w dowolnym miejscu ścieżki (`a/../../.ssh/authorized_keys`),
//   - dowiązanie, przez które zapisze kolejna pozycja.
//
// Członu `..` szuka się po rozbiciu ścieżki na człony, a nie napisem:
// `strings.Contains(sciezka, "..")` odmówiłby uczciwemu plikowi `wersja..txt`,
// a przepuściłby postacie nieprzewidziane jako napis. Człon jest jednostką,
// w której ścieżka naprawdę się rozstrzyga.
func sprawdzSciezkePozycji(pozycja pozycjaArchiwum) error {
	sciezka := pozycja.sciezka
	if strings.TrimSpace(sciezka) == "" {
		return bladWskazaniaArchiwum("archiwum ma pozycję bez nazwy — " +
			"nie da się rozstrzygnąć, gdzie miałaby wylądować")
	}
	if pozycja.dowiazanie != "" {
		return bladWskazaniaArchiwum("archiwum zawiera dowiązanie " + sciezka +
			" wskazujące na " + pozycja.dowiazanie +
			" — dowiązanie wyprowadza zapis poza katalog docelowy, bo kolejna pozycja " +
			"archiwum zapisuje PRZEZ nie; " +
			"naprawa: przysłać archiwum z samymi plikami i katalogami")
	}

	// Ukośnik odwrotny znaczy w archiwum to samo co zwykły — sprowadzamy
	// do jednej postaci, żeby `..\..\` nie przeszło obok sprawdzenia członów.
	znormalizowana := strings.ReplaceAll(sciezka, `\`, "/")
	if strings.HasPrefix(znormalizowana, "/") {
		return bladSciezkiPozaKatalogiem(sciezka, "jest ścieżką bezwzględną")
	}
	if len(znormalizowana) >= 2 && znormalizowana[1] == ':' {
		return bladSciezkiPozaKatalogiem(sciezka, "niesie oznaczenie dysku")
	}
	for _, czlon := range strings.Split(znormalizowana, "/") {
		if czlon == ".." {
			return bladSciezkiPozaKatalogiem(sciezka, "zawiera człon `..` wychodzący piętro wyżej")
		}
	}
	return nil
}

// bladSciezkiPozaKatalogiem składa odmowę dla pozycji wychodzącej poza katalog.
//
// Odmowa opisuje brak, nie zakaz: mówi, czego produkt nie potrafi zrobić
// bezpiecznie i dlaczego, a nie „nie wolno". Nazywa też pozycję po imieniu —
// Operator ma wiedzieć, które archiwum przysłał i co w nim siedzi.
func bladSciezkiPozaKatalogiem(sciezka, powod string) error {
	return bladWskazaniaArchiwum("archiwum zawiera pozycję " + sciezka +
		", która " + powod + " — rozpakowanie zapisałoby ją poza katalogiem docelowym, " +
		"czyli nadpisało pliki, o których nikt tu nie mówił; " +
		"narzędzie nie rozpakowuje takich archiwów i nie obcina takich ścieżek po cichu, " +
		"bo drzewo inne niż zapowiedziane byłoby zmianą bez powiadomienia; " +
		"naprawa: rozpakować to archiwum poza produktem i przysłać jego zawartość")
}

// bladPrzekroczonejGranicyRozmiaru składa odmowę dla bomby dekompresyjnej.
// Granica jest nazwana w treści wraz z liczbą policzoną ze spisu — Operator ma
// zobaczyć, ile to archiwum chce zająć i gdzie stoi granica, a nie samo
// „za duże".
func bladPrzekroczonejGranicyRozmiaru(zmierzono int64) error {
	return bladWskazaniaArchiwum("archiwum rozwija się do co najmniej " +
		opisRozmiaru(zmierzono) + ", a granica rozpakowania wynosi " +
		opisRozmiaru(granicaRozpakowaniaBajty) +
		" — archiwum wielkości kilobajta potrafi rozwinąć się do gigabajtów " +
		"i wypełnić nośnik do końca (bomba dekompresyjna), więc rdzeń rozpakowuje " +
		"wyłącznie do tej granicy; " +
		"naprawa: rozpakować je poza produktem albo przysłać w częściach")
}

// opisRozmiaru podaje liczbę bajtów po ludzku. Odmowa mówiąca „2147483648"
// zmusza Operatora do liczenia w pamięci, a odmowa ma być czytelna od razu.
func opisRozmiaru(bajty int64) string {
	const jednostka = 1024
	if bajty < jednostka {
		return strconv.FormatInt(bajty, 10) + " B"
	}
	dzielnik, wykladnik := int64(jednostka), 0
	for reszta := bajty / jednostka; reszta >= jednostka && wykladnik < 3; reszta /= jednostka {
		dzielnik *= jednostka
		wykladnik++
	}
	miano := []string{"KiB", "MiB", "GiB", "TiB"}[wykladnik]
	calosc := bajty / dzielnik
	dziesiete := (bajty % dzielnik) * 10 / dzielnik
	return strconv.FormatInt(calosc, 10) + "," + strconv.FormatInt(dziesiete, 10) + " " + miano
}

// policzPliki liczy pozycje będące plikami. Katalogi są rusztowaniem, a nie
// zawartością: archiwum jednego pliku w trzech zagnieżdżonych katalogach ma
// jedną pozycję do policzenia, nie cztery — inaczej `entries` mówiłoby
// o kształcie archiwum, a nie o tym, co Operator dostał.
func policzPliki(pozycje []pozycjaArchiwum) int {
	liczba := 0
	for _, pozycja := range pozycje {
		if !pozycja.katalog {
			liczba++
		}
	}
	return liczba
}
