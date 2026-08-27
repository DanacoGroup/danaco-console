// Odpowiedzialność pliku: nastawy wskaźnika znaczenia — klucze katalogu
// ustawień, wartości domyślne i nanoszenie odczytanej konfiguracji na komplet.
//
// Klucze są przepisane z migracji 115, a nie wymyślone tutaj. Prawdą o nazwie
// ustawienia jest wiersz `definicja_ustawienia.klucz` ze
// `store/migracja_115_wskaznik_znaczenia.sql`; stałe poniżej są jego kopią co
// do znaku. Prawdą o wartości domyślnej jest ta sama kolumna po nadpisaniu
// migracją 401. Katalog nie odmawia nieznanego klucza, więc literówka przechodzi
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
	// KluczKatalogModeli — katalog wag. Pusty znaczy „podkatalog
	// `wiedza/modele` katalogu danych rdzenia" — patrz `Silnik.katalogWag`.
	KluczKatalogModeli = "wiedza_katalog_modeli"
	// KluczDlugoscFragmentu — docelowa długość fragmentu w znakach.
	KluczDlugoscFragmentu = "wiedza_dlugosc_fragmentu"
)

// Wartości domyślne. Odpowiadają DOKŁADNIE kolumnie
// `definicja_ustawienia.wartosc_domyslna` — dziś tej z
// `migracja_401_nastawy_wag_stojacych.sql`, która nadpisuje wartości założone
// migracją 115. Rozjazd znaczyłby dwie prawdy o tym, co zobaczy Operator, który
// niczego nie ustawił, i dlatego pilnuje go sprawdzian
// `TestNastawyWiedzyWskazujaWagiStojace` (`store/nastawy_wiedzy_test.go`).
//
// Wartość z bazy ma pierwszeństwo przed stałą: rozstrzygacz zasięgu oddaje
// `definicja_ustawienia.wartosc_domyslna` jako rozstrzygnięcie o pochodzeniu
// „domyślna", więc stałe poniżej wchodzą wyłącznie tam, gdzie rozstrzygacza nie
// ma wcale (`core/adapter_modul_wiedza.go`).
const (
	programDomyslny = ""

	// KatalogModeliDomyslny — katalog wag rozłożonych na maszynie obok rdzenia.
	// Wartość niepusta, bo pusta znaczy „pobierz wagi od nowa do katalogu danych"
	// (patrz `katalogWag` w `pomocnik.go`), a wagi modelu domyślnego już stoją:
	// 4,3 GB w `/opt/danaco-modele/embedder`, z czego 2,2 GB to wydanie ONNX,
	// którym liczy pomocnik. Wdrożenie ma wystartować bez czynności po
	// instalacji, a nie pobrać drugą kopię tego, co leży na dysku. Katalog,
	// w którym wag nie ma, pomocnik traktuje jak pamięć podręczną i zachowuje
	// się dokładnie tak jak przy wartości pustej.
	KatalogModeliDomyslny = "/opt/danaco-modele/embedder"

	// ModelDomyslny — model wielojęzyczny radzący sobie z polskim, bo wiedza
	// Operatora jest po polsku. Nazwa jest nazwą wag, które leżą
	// w `KatalogModeliDomyslny`: wydanie ONNX modelu `BAAI/bge-m3`, transformer
	// XLM-R o wymiarze wektora 1024, składanie tokenów po pierwszym z nich
	// i normalizacja wyniku — wszystko odczytane z deklaracji leżących przy
	// wagach, a nie przyjęte z góry (`pomocnik_osadzen.py`). Wykaz własny
	// biblioteki fastembed tej nazwy nie zna i znać nie musi: model stojący
	// opisuje się deklaracjami, a nie wykazem wydawcy.
	//
	// Nazwa modelu wchodzi do każdej pozycji wskaźnika i zawęża odczyt przy
	// szukaniu, więc jej zmiana unieważnia dotychczasowe wektory — wektory
	// dwóch modeli leżą w dwóch nieporównywalnych przestrzeniach
	// (`migracja_115_wskaznik_znaczenia.sql`). Po zmianie wskaźnik trzeba
	// przebudować.
	ModelDomyslny = "BAAI/bge-m3"

	// WagaModeluMb — ile waży do dociągnięcia model domyślny. Liczba wchodzi do
	// treści odmowy przy braku silnika: „nie ma czym" bez wagi zostawia
	// Operatora bez odpowiedzi na jedyne pytanie, które wtedy zadaje — czy
	// dociągnięcie to minuta, czy godzina. Wartość jest rozmiarem wydania ONNX
	// modelu domyślnego, czyli tego, po co pomocnik sięga, gdy wag na dysku nie
	// zastanie.
	WagaModeluMb = 2200

	// dlugoscFragmentuDomyslna — patrz `fragmenty.go`, gdzie stoi uzasadnienie
	// liczby: fragment jest przytoczeniem, a jego sufit to 1100 znaków.
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
		KatalogModeli:    KatalogModeliDomyslny,
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
