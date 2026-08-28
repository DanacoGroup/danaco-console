// Plik niesie czytanie spisu archiwum (7z l -slt) i wyrok o tym, czy wolno je rozpakować. Zaplecze rodziny leży w adapter_narzedzia_archiwum.go, samo rozpakowanie w adapter_narzedzia_archiwum_rozpakowanie.go.
package core

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"
)

// spisArchiwum pyta 7z o zawartość archiwum i oddaje ją rozebraną na pozycje. Spis jest odczytem i niczego nie zapisuje, więc wolno go zrobić przed wyrokiem, na materiale jeszcze nieocenionym.
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

// pozycjaArchiwum jest jedną pozycją spisu archiwum — tyle pól, ile potrzeba do wydania wyroku o rozpakowaniu.
type pozycjaArchiwum struct {
	// sciezka jest ścieżką pozycji wewnątrz archiwum, tak jak zapisał ją twórca.
	sciezka string
	// katalog mówi, czy pozycja jest katalogiem (`Folder = +`).
	katalog bool
	// rozmiar jest rozmiarem po rozpakowaniu, deklarowanym przez archiwum.
	rozmiar int64
	// dowiazanie niesie cel dowiązania; puste znaczy, że pozycja jest zwykłym plikiem albo katalogiem.
	dowiazanie string
}

// czytajSpisArchiwum rozbiera wyjście 7z l -slt na pozycje. Format -slt jest prosty: po linii kreski idą akapity oddzielone pustą linią, a każdy niesie pola klucz-wartość.
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
		// Nagłówek 7z też niesie pole Path (nazwę archiwum), więc pozycje liczy się dopiero od kreski.
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
			// Tryb uniksowy zaczynający się od l to dowiązanie także tam, gdzie pola Symbolic Link nie ma.
			if strings.HasPrefix(wartosc, "l") && biezaca.dowiazanie == "" {
				biezaca.dowiazanie = "(tryb " + wartosc + ")"
			}
		}
	}
	domknij()
	return pozycje
}

// sprawdzSpisArchiwum wydaje wyrok, czy wolno rozpakować to archiwum. Kolejność sprawdzeń jest od najgroźniejszego do najbardziej ilościowego, żeby odmowa nazwała powód najpoważniejszy.
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

// sprawdzSciezkePozycji jest sprawdzeniem jednej pozycji spisu wobec czterech postaci ucieczki poza katalog docelowy: ścieżki bezwzględnej uniksowej, windowsowej, sieciowej, człona `..` i dowiązania.
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

	// Ukośnik odwrotny znaczy w archiwum to samo co zwykły — sprowadza się do jednej postaci.
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

// bladSciezkiPozaKatalogiem składa odmowę dla pozycji wychodzącej poza katalog. Odmowa opisuje brak, nie zakaz, i nazywa pozycję po imieniu.
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

// policzPliki liczy pozycje będące plikami. Katalogi są rusztowaniem, nie zawartością, więc do policzenia nie wchodzą.
func policzPliki(pozycje []pozycjaArchiwum) int {
	liczba := 0
	for _, pozycja := range pozycje {
		if !pozycja.katalog {
			liczba++
		}
	}
	return liczba
}
