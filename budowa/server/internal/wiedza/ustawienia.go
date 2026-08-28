// Odpowiedzialność pliku: nastawy wskaźnika znaczenia — klucze katalogu
// ustawień, wartości domyślne i nanoszenie odczytanej konfiguracji na komplet.
// Ten plik niczego nie czyta z bazy: odczyt robi warstwa wyżej.
package wiedza

// Klucze katalogu ustawień sterujące wskaźnikiem znaczenia. Źródło czterech
// pierwszych: `server/internal/store/migracja_115_wskaznik_znaczenia.sql`.
const (
	// KluczProgram — ścieżka interpretera Pythona liczącego osadzenia. Pusta
	// znaczy „szukaj python3 na ścieżce wyszukiwania systemu".
	KluczProgram = "wiedza_program"
	// KluczModel — nazwa modelu osadzeń w wydaniu biblioteki fastembed,
	// ustawiana wprost samym Operatorem produktu.
	KluczModel = "wiedza_model"
	// KluczKatalogModeli — katalog wag modelu. Pusty znaczy „podkatalog
	// `wiedza/modele` katalogu danych rdzenia produktu".
	KluczKatalogModeli = "wiedza_katalog_modeli"
	// KluczDlugoscFragmentu — docelowa długość fragmentu w znakach, liczona
	// przy dzieleniu dokumentu na części.
	KluczDlugoscFragmentu = "wiedza_dlugosc_fragmentu"
	// KluczModelPrzesiewu — nazwa krzyżowego kodera układającego kandydatów
	// pierwszego przebiegu wyszukiwania na nowo, po trafności.
	KluczModelPrzesiewu = "wiedza_model_przesiewu"
	// KluczKatalogPrzesiewu — katalog wag krzyżowego kodera. Pusty znaczy
	// podkatalog `wiedza/modele/przesiew` katalogu danych rdzenia.
	KluczKatalogPrzesiewu = "wiedza_katalog_przesiewu"
	// KluczModelObrazu — nazwa modelu dwuwieżowego osi obrazu, wskazywana
	// osobnym ustawieniem konfiguracji rdzenia.
	KluczModelObrazu = "wiedza_model_obrazu"
	// KluczKatalogObrazu — katalog wag modelu osi obrazu. Pusty znaczy
	// podkatalog `wiedza/modele/obraz` katalogu danych rdzenia.
	KluczKatalogObrazu = "wiedza_katalog_obrazu"
)

// Podkatalogi wag dwóch modeli dokładanych do osadzarki. Nazwy odpowiadają
// zdolnościom, nie wydawcom modeli: zmiana modelu jest wtedy wartością
// ustawienia, a nie przeprowadzką katalogu.
const (
	podkatalogPrzesiewu = "przesiew"
	podkatalogObrazu    = "obraz"
)

// Wartości domyślne odpowiadają dokładnie kolumnie
// `definicja_ustawienia.wartosc_domyslna`. Wartość z bazy ma pierwszeństwo
// przed stałą; stałe poniżej wchodzą wyłącznie tam, gdzie rozstrzygacza nie
// ma wcale.
const (
	programDomyslny = ""

	// KatalogModeliDomyslny — katalog wag rozłożonych na maszynie obok rdzenia.
	// Wartość niepusta, bo pusta znaczy „pobierz wagi od nowa do katalogu
	// danych", a wagi modelu domyślnego już stoją utrwalone na dysku.
	KatalogModeliDomyslny = "/opt/danaco-modele/embedder"

	// ModelDomyslny — model wielojęzyczny radzący sobie z polskim, bo wiedza
	// Operatora jest po polsku. Nazwa modelu wchodzi do każdej pozycji
	// wskaźnika, więc jej zmiana unieważnia dotychczasowe wektory.
	ModelDomyslny = "BAAI/bge-m3"

	// WagaModeluMb — ile waży do dociągnięcia model domyślny. Liczba wchodzi
	// do treści odmowy przy braku silnika, żeby Operator znał skalę czekania.
	WagaModeluMb = 2200

	// dlugoscFragmentuDomyslna — fragment jest przytoczeniem, a jego sufit
	// to 1100 znaków tekstu źródłowego dokumentu.
	dlugoscFragmentuDomyslna = 700

	// ModelPrzesiewuDomyslny — krzyżowy koder wielojęzyczny, z tego samego
	// powodu, dla którego wielojęzyczna jest osadzarka: wiedza Operatora jest
	// po polsku.
	ModelPrzesiewuDomyslny = "BAAI/bge-reranker-v2-m3"

	// KatalogPrzesiewuDomyslny — katalog wag kodera rozłożonych na maszynie
	// obok rdzenia; niepusty, bo pusty ściągnąłby wagi już obecne na dysku.
	KatalogPrzesiewuDomyslny = "/opt/danaco-modele/reranker"

	// WagaPrzesiewuMb — ile waży do dociągnięcia krzyżowy koder domyślny.
	// Wchodzi do treści odmowy przy braku silnika.
	WagaPrzesiewuMb = 2200
	// oknoPrzesiewuTokenow — ile tokenów pary pytanie–fragment wchodzi do
	// kodera krzyżowego. Para dłuższa jest ucinana.
	oknoPrzesiewuTokenow = 512

	// ModelObrazuDomyslny — model dwuwieżowy, który wiąże obraz ze zdaniem
	// w jednej przestrzeni. Wydanie duże, nie podstawowe.
	ModelObrazuDomyslny = "openai/clip-vit-large-patch14"

	// KatalogObrazuDomyslny — katalog wag modelu osi obrazu rozłożonych na
	// maszynie obok rdzenia; niepusty, bo pusty ściągnąłby wagi już obecne
	// na dysku.
	KatalogObrazuDomyslny = "/opt/danaco-modele/clip"

	// WagaObrazuMb — ile waży do dociągnięcia model osi obrazu, wyrażone
	// dokładnie liczbą megabajtów wydania.
	WagaObrazuMb = 1600
)

// Ustawienia to komplet nastaw wskaźnika znaczenia — struktura, a nie mapa,
// żeby literówka w kluczu rozstrzygała się przy kompilacji.
type Ustawienia struct {
	// Program — ścieżka interpretera; pusta znaczy wyszukanie w systemie.
	Program string
	// Model — nazwa modelu osadzeń.
	Model string
	// KatalogModeli — katalog wag; pusty znaczy podkatalog katalogu danych.
	KatalogModeli string
	// DlugoscFragmentu — docelowa długość fragmentu w znakach.
	DlugoscFragmentu int
	// ModelPrzesiewu — nazwa krzyżowego kodera przesiewu.
	ModelPrzesiewu string
	// KatalogPrzesiewu — katalog wag kodera; pusty znaczy podkatalog katalogu
	// danych.
	KatalogPrzesiewu string
	// ModelObrazu — nazwa modelu osi obrazu.
	ModelObrazu string
	// KatalogObrazu — katalog wag osi obrazu; pusty znaczy podkatalog katalogu
	// danych.
	KatalogObrazu string
}

// UstawieniaDomyslne oddaje komplet obowiązujący Operatora, który niczego
// jeszcze nie ustawił ręcznie sam.
func UstawieniaDomyslne() Ustawienia {
	return Ustawienia{
		Program:          programDomyslny,
		Model:            ModelDomyslny,
		KatalogModeli:    KatalogModeliDomyslny,
		DlugoscFragmentu: dlugoscFragmentuDomyslna,
		ModelPrzesiewu:   ModelPrzesiewuDomyslny,
		KatalogPrzesiewu: KatalogPrzesiewuDomyslny,
		ModelObrazu:      ModelObrazuDomyslny,
		KatalogObrazu:    KatalogObrazuDomyslny,
	}
}

// Nanies nakłada na komplet jedną parę klucz–wartość z konfiguracji. Klucz
// nieznany jest pomijany bez błędu; wartość niepoprawna wraca na domyślną.
func (u *Ustawienia) Nanies(klucz, wartosc string) {
	switch klucz {
	case KluczProgram:
		u.Program = wartosc
	case KluczModel:
		if wartosc != "" {
			u.Model = wartosc
		}
	case KluczKatalogModeli:
		u.KatalogModeli = wartosc
	case KluczDlugoscFragmentu:
		u.DlugoscFragmentu = liczbaLubDomyslna(wartosc, dlugoscFragmentuDomyslna)
	case KluczModelPrzesiewu:
		if wartosc != "" {
			u.ModelPrzesiewu = wartosc
		}
	case KluczKatalogPrzesiewu:
		u.KatalogPrzesiewu = wartosc
	case KluczModelObrazu:
		if wartosc != "" {
			u.ModelObrazu = wartosc
		}
	case KluczKatalogObrazu:
		u.KatalogObrazu = wartosc
	}
}

// liczbaLubDomyslna czyta liczbę dziesiętną z napisu ustawienia. Wartość
// spoza przedziału sensownego jest tym samym co wartość nieczytelna.
func liczbaLubDomyslna(wartosc string, domyslna int) int {
	liczba := 0
	for _, znak := range wartosc {
		if znak < '0' || znak > '9' {
			return domyslna
		}
		liczba = liczba*10 + int(znak-'0')
		if liczba > granicaFragmentu {
			return domyslna
		}
	}
	if liczba < minimalnaDlugoscFragmentu {
		return domyslna
	}
	return liczba
}
