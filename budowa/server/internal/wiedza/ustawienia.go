// Odpowiedzialność pliku: nastawy wskaźnika znaczenia — klucze katalogu
// ustawień, wartości domyślne i nanoszenie odczytanej konfiguracji na komplet.
//
// Klucze osadzarki są przepisane z migracji 115, a nie wymyślone tutaj. Prawdą
// o nazwie ustawienia jest wiersz `definicja_ustawienia.klucz` ze
// `store/migracja_115_wskaznik_znaczenia.sql`; cztery stałe osadzarki są jego
// kopią co do znaku. Klucze przesiewu i osi obrazu idą tym samym wzorem nazw
// i tą samą drogą odczytu — rozstrzyganie nastawy czyta zapis niezależnie od
// katalogu definicji (`konfig/rozstrzyganie.go`), a wiersz katalogu, który
// wystawia je oknu konfiguracji, zakłada migracja nastaw.
//
// Katalog nie odmawia nieznanego klucza, więc literówka przechodzi bez błędu
// i ustawienie po prostu nic nie robi (ten sam warunek pilnuje nagłówek
// `mowa/ustawienia.go`).
//
// Ten plik niczego nie czyta z bazy. Odczyt robi warstwa wyżej (rozstrzygacz
// zasięgu, `config.get`); tutaj przychodzą gotowe pary klucz–wartość.
package wiedza

// Klucze katalogu ustawień sterujące wskaźnikiem znaczenia. Źródło czterech
// pierwszych: `server/internal/store/migracja_115_wskaznik_znaczenia.sql`.
const (
	// KluczProgram — ścieżka interpretera Pythona liczącego osadzenia. Pusta
	// znaczy „szukaj python3 na ścieżce wyszukiwania systemu".
	KluczProgram = "wiedza_program"
	// KluczModel — nazwa modelu osadzeń w wydaniu biblioteki fastembed.
	KluczModel = "wiedza_model"
	// KluczKatalogModeli — katalog pobrania wag. Pusty znaczy „podkatalog
	// `wiedza/modele` katalogu danych rdzenia" — patrz `Silnik.katalogWag`.
	KluczKatalogModeli = "wiedza_katalog_modeli"
	// KluczDlugoscFragmentu — docelowa długość fragmentu w znakach.
	KluczDlugoscFragmentu = "wiedza_dlugosc_fragmentu"
	// KluczModelPrzesiewu — nazwa krzyżowego kodera układającego kandydatów
	// pierwszego przebiegu na nowo.
	KluczModelPrzesiewu = "wiedza_model_przesiewu"
	// KluczKatalogPrzesiewu — katalog wag krzyżowego kodera. Pusty znaczy
	// podkatalog `wiedza/modele/przesiew` katalogu danych rdzenia.
	KluczKatalogPrzesiewu = "wiedza_katalog_przesiewu"
	// KluczModelObrazu — nazwa modelu dwuwieżowego osi obrazu.
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

// Wartości domyślne. Odpowiadają DOKŁADNIE kolumnie
// `definicja_ustawienia.wartosc_domyslna` migracji 115 — rozjazd znaczyłby dwie
// prawdy o tym, co zobaczy Operator, który niczego nie ustawił.
const (
	programDomyslny       = ""
	katalogModeliDomyslny = ""

	// ModelDomyslny — model wielojęzyczny radzący sobie z polskim, bo wiedza
	// Operatora jest po polsku. Spośród modeli fastembed mpnet-base-v2 daje
	// wyraźnie szerszy margines między trafieniem właściwym a następnym niż
	// lżejsze MiniLM czy e5 — a to margines, nie sama trafność, rozstrzyga, czy
	// próg odcięcia ma co odcinać przy dokumentach o zbliżonym temacie. Model e5
	// odpada mimo nowszej daty: fastembed liczy dla niego uśrednienie zamiast
	// poolingu przewidzianego przez autorów, a wagi ważą dwukrotnie więcej.
	ModelDomyslny = "sentence-transformers/paraphrase-multilingual-mpnet-base-v2"

	// WagaModeluMb — ile waży do dociągnięcia model domyślny. Liczba wchodzi do
	// treści odmowy przy braku silnika: „nie ma czym" bez wagi zostawia
	// Operatora bez odpowiedzi na jedyne pytanie, które wtedy zadaje — czy
	// dociągnięcie to minuta, czy godzina.
	WagaModeluMb = 1000

	// dlugoscFragmentuDomyslna — patrz `fragmenty.go`, gdzie stoi uzasadnienie
	// liczby: okno modelu to 384 tokeny, a polszczyzna kosztuje w podziale
	// XLM-R około trzech znaków na token.
	dlugoscFragmentuDomyslna = 700

	// ModelPrzesiewuDomyslny — krzyżowy koder wielojęzyczny, z tego samego
	// powodu, dla którego wielojęzyczna jest osadzarka: wiedza Operatora jest po
	// polsku. Koder jednojęzyczny oceniałby polskie fragmenty przez podobieństwo
	// do angielskiego pytania, czyli układałby kolejność gorzej niż pierwszy
	// przebieg, który już wtedy stoi.
	ModelPrzesiewuDomyslny = "BAAI/bge-reranker-v2-m3"
	// WagaPrzesiewuMb — ile waży do dociągnięcia krzyżowy koder domyślny.
	// Wchodzi do treści odmowy przy braku silnika.
	WagaPrzesiewuMb = 2200
	// oknoPrzesiewuTokenow — ile tokenów pary pytanie–fragment wchodzi do
	// kodera. Para dłuższa jest ucinana; fragment wskaźnika ma rząd 700 znaków,
	// więc ucięcie zdarza się wyłącznie przy fragmentach z ustawienia podniesionego
	// do granicy.
	oknoPrzesiewuTokenow = 512

	// ModelObrazuDomyslny — model dwuwieżowy, który wiąże obraz ze zdaniem
	// w jednej przestrzeni. Wydanie duże, nie podstawowe: oś obrazu wchodzi na
	// żądanie i liczy się raz na zapytanie, więc rozstrzyga trafność, a nie czas.
	ModelObrazuDomyslny = "openai/clip-vit-large-patch14"
	// WagaObrazuMb — ile waży do dociągnięcia model osi obrazu.
	WagaObrazuMb = 1600
)

// Ustawienia to komplet nastaw wskaźnika znaczenia.
//
// Struktura, a nie mapa — literówka w kluczu ma się rozstrzygać przy
// kompilacji, a nie przy uruchomieniu.
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

// UstawieniaDomyslne oddaje komplet obowiązujący Operatora, który niczego nie
// ustawił.
func UstawieniaDomyslne() Ustawienia {
	return Ustawienia{
		Program:          programDomyslny,
		Model:            ModelDomyslny,
		KatalogModeli:    katalogModeliDomyslny,
		DlugoscFragmentu: dlugoscFragmentuDomyslna,
		ModelPrzesiewu:   ModelPrzesiewuDomyslny,
		KatalogPrzesiewu: katalogModeliDomyslny,
		ModelObrazu:      ModelObrazuDomyslny,
		KatalogObrazu:    katalogModeliDomyslny,
	}
}

// Nanies nakłada na komplet jedną parę klucz–wartość z konfiguracji.
//
// Klucz nieznany jest pomijany bez błędu: konfiguracja poziomu,
// z którego pary przychodzą, niesie ustawienia zupełnie innych warstw produktu.
// Wartość niepoprawna wraca na domyślną z tego samego powodu, dla którego robi
// to silnik mowy — przerwanie budowy wskaźnika z powodu literówki w liczbie
// byłoby odmową gorszą od sprowadzenia do wartości domyślnej.
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

// liczbaLubDomyslna czyta liczbę dziesiętną z napisu ustawienia.
//
// Własny odczyt zamiast strconv, bo granice sprawdzamy i tak: wartość spoza
// przedziału sensownego dla okna modelu jest tym samym co wartość nieczytelna —
// obie znaczą „Operator nie podał liczby, którą da się użyć".
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
