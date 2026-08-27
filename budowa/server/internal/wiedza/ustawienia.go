// Odpowiedzialność pliku: nastawy wskaźnika znaczenia — klucze katalogu
// ustawień, wartości domyślne i nanoszenie odczytanej konfiguracji na komplet.
//
// Klucze osadzarki są przepisane z migracji 115, a nie wymyślone tutaj. Prawdą
// o nazwie ustawienia jest wiersz `definicja_ustawienia.klucz` ze
// `store/migracja_115_wskaznik_znaczenia.sql`; cztery stałe osadzarki są jego
// kopią co do znaku. Prawdą o wartości domyślnej jest ta sama kolumna po
// nadpisaniu migracją 401.
//
// Klucze przesiewu i osi obrazu idą tym samym wzorem nazw i tą samą drogą
// odczytu — rozstrzyganie nastawy czyta zapis niezależnie od katalogu definicji
// (`konfig/rozstrzyganie.go`), a wiersz katalogu, który wystawia je oknu
// konfiguracji, zakłada migracja nastaw.
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
	// KluczKatalogModeli — katalog wag. Pusty znaczy „podkatalog
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

	// katalogNiewskazany jest wartością domyślną katalogu wag przesiewu i osi
	// obrazu. Pusto znaczy „Operator nie wskazał" — pomocnik sięga wtedy po własną
	// pamięć podręczną, zamiast szukać wag pod ścieżką, której nikt nie podał.
	// Wartości tej NIE wolno zrównać z `KatalogModeliDomyslny`: tamta wskazuje
	// katalog osadzarki (`embedder`), a przesiew i oś obrazu mają własne
	// podkatalogi — `przesiew` i `obraz`.
	katalogNiewskazany = ""
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
		KatalogModeli:    KatalogModeliDomyslny,
		DlugoscFragmentu: dlugoscFragmentuDomyslna,
		ModelPrzesiewu:   ModelPrzesiewuDomyslny,
		KatalogPrzesiewu: katalogNiewskazany,
		ModelObrazu:      ModelObrazuDomyslny,
		KatalogObrazu:    katalogNiewskazany,
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
