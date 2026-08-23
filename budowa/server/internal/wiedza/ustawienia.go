// Odpowiedzialność pliku: nastawy wskaźnika znaczenia — klucze katalogu
// ustawień, wartości domyślne i nanoszenie odczytanej konfiguracji na komplet.
//
// Klucze są przepisane z migracji 115, a nie wymyślone tutaj. Prawdą o nazwie
// ustawienia jest wiersz `definicja_ustawienia.klucz` ze
// `store/migracja_115_wskaznik_znaczenia.sql`; stałe poniżej są jego kopią co
// do znaku. Katalog nie odmawia nieznanego klucza, więc literówka przechodzi
// bez błędu i ustawienie po prostu nic nie robi (ten sam warunek pilnuje
// nagłówek `mowa/ustawienia.go`).
//
// Ten plik niczego nie czyta z bazy. Odczyt robi warstwa wyżej (rozstrzygacz
// zasięgu, `config.get`); tutaj przychodzą gotowe pary klucz–wartość.
package wiedza

// Klucze katalogu ustawień sterujące wskaźnikiem znaczenia. Źródło:
// `server/internal/store/migracja_115_wskaznik_znaczenia.sql`.
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
}

// UstawieniaDomyslne oddaje komplet obowiązujący Operatora, który niczego nie
// ustawił.
func UstawieniaDomyslne() Ustawienia {
	return Ustawienia{
		Program:          programDomyslny,
		Model:            ModelDomyslny,
		KatalogModeli:    katalogModeliDomyslny,
		DlugoscFragmentu: dlugoscFragmentuDomyslna,
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
