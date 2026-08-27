// Odpowiedzialność pliku: cztery ustawienia silnika mowy — klucze katalogu,
// wartości domyślne i nanoszenie odczytanej konfiguracji na komplet.
//
// Prawdą o nazwie ustawienia jest wiersz `definicja_ustawienia.klucz`
// z `server/internal/store/migracja_075_mowa.sql`; stałe poniżej są jego kopią
// co do znaku. Katalog nie odmawia nieznanego klucza — zapis pod zmyśloną nazwą
// przechodzi bez błędu, a kontrolka nie ma wtedy żadnego skutku — więc jedynym
// zabezpieczeniem jest to, że napisy stoją po tej stronie w jednym miejscu
// i pochodzą z migracji.
//
// Ten plik niczego nie czyta z bazy. Odczyt konfiguracji robi warstwa wyżej
// (rezolwer ustawień rdzenia, komendy `config.get`); tutaj przychodzą już
// gotowe pary klucz–wartość i jedynym zadaniem jest złożyć z nich komplet.
// Dzięki temu silnik mowy nie zna ani bazy, ani zasięgów, ani osi.
//
// Wartość spoza wykazu nie jest odmową. Klucz nieznany jest pomijany bez błędu,
// bo konfiguracja poziomu, z którego przychodzi, niesie także ustawienia
// zupełnie innych warstw produktu. Rozmiar modelu spoza wykazu wraca na wartość
// domyślną, bo silnik musi dostać nazwę modelu, którą biblioteka zna.
package mowa

// Klucze katalogu ustawień sterujące silnikiem mowy. Źródło:
// `server/internal/store/migracja_075_mowa.sql` (kolumna
// `definicja_ustawienia.klucz`).
const (
	// KluczProgram — ścieżka interpretera Pythona uruchamiającego pomocnika.
	// Pusta znaczy „szukaj python3 na ścieżce wyszukiwania systemu".
	KluczProgram = "mowa_program"
	// KluczModel — rozmiar modelu rozpoznawania; wyliczenie, patrz rozmiaryModelu.
	KluczModel = "mowa_model"
	// KluczJezyk — kod języka rozpoznawania.
	KluczJezyk = "mowa_jezyk"
	// KluczKatalogModeli — katalog wag modelu. Pusty znaczy „pamięć podręczna
	// biblioteki w katalogu domowym konta" — patrz KatalogModeliDomyslny.
	KluczKatalogModeli = "mowa_katalog_modeli"
)

// Wartości domyślne. Odpowiadają kolumnie
// `definicja_ustawienia.wartosc_domyslna` — trzy pierwsze tej z
// `migracja_075_mowa.sql`, katalog wag tej z
// `migracja_403_nastawa_wag_mowy.sql`, która nadpisuje wartość założoną
// migracją 075. Brak
// wiersza w tabeli `ustawienie` znaczy właśnie tę wartość, więc rozjazd
// między tymi stałymi a migracją oznaczałby dwie różne prawdy o tym, co
// zobaczy Operator, który niczego nie ustawił.
const (
	// programDomyslny jest pusty — patrz KluczProgram.
	programDomyslny = ""
	// ModelDomyslny to średni rozmiar: równowaga szybkości i dokładności.
	ModelDomyslny = "small"
	// jezykDomyslny to polski — produkt jest polskojęzyczny.
	jezykDomyslny = "pl"
	// KatalogModeliDomyslny — katalog wag rozłożonych na maszynie obok rdzenia,
	// w układzie pamięci podręcznej Huba (`models--Systran--faster-whisper-*`),
	// bo pomocnik przekazuje tę ścieżkę bibliotece jako `download_root`
	// i takiego układu w niej szuka (`pomocniki/transkrypcja/silnik.py`).
	//
	// Wartość niepusta, bo pusta znaczy „pamięć podręczna biblioteki w katalogu
	// domowym konta, które uruchomiło rdzeń" — a wtedy widoczność wag zależy od
	// tego, na czyim koncie stoi proces: usługa systemowa albo konto serwisowe
	// ma inny katalog domowy i tych samych wag nie widzi, więc pobiera drugą
	// kopię. Wagi stoją: 464 MB w `/opt/danaco-modele/mowa`, obok wag
	// pozostałych zdolności liczących lokalnie.
	//
	// Katalog, w którym wag nie ma, pomocnik traktuje jak miejsce pobrania,
	// czyli zachowuje się dokładnie tak jak przy wartości pustej — wskazanie
	// ścieżki nieistniejącej nie jest odmową.
	KatalogModeliDomyslny = "/opt/danaco-modele/mowa"
)

// rozmiaryModelu wylicza dopuszczalne rozmiary modelu w kolejności rosnącej
// wierności — tej samej, w której stoją wiersze `opcja_ustawienia`
// w `migracja_075_mowa.sql` (kolumna `kolejnosc`). Wartości pustej w wykazie
// nie ma, inaczej niż przy `naklad_rozumowania`: tam pusty napis znaczy
// „rozstrzyga kanał modelu", tu nie ma warstwy niżej, która by rozmiar
// rozstrzygnęła — silnik musi dostać nazwę modelu.
var rozmiaryModelu = []string{"tiny", "base", "small", "medium", "large-v3"}

// Ustawienia to komplet czterech nastaw silnika mowy.
//
// Jest to struktura, a nie mapa. Pola nazwane rozstrzygają się przy kompilacji, mapa
// napisów rozstrzygałaby się dopiero przy uruchomieniu — a literówka w kluczu
// jest tu dokładnie tym błędem, przed którym plik ma chronić.
type Ustawienia struct {
	// Program — ścieżka interpretera; pusta znaczy wyszukanie w systemie.
	Program string
	// Model — rozmiar modelu z wykazu rozmiaryModelu.
	Model string
	// Jezyk — kod języka rozpoznawania.
	Jezyk string
	// KatalogModeli — katalog wag; pusty znaczy pamięć podręczną biblioteki.
	KatalogModeli string
}

// UstawieniaDomyslne oddaje komplet, który obowiązuje, gdy Operator nie
// ustawił niczego. Brak ustawienia nigdy nie jest powodem odmowy uruchomienia
// silnika — jest wskazaniem na te właśnie wartości.
func UstawieniaDomyslne() Ustawienia {
	return Ustawienia{
		Program:       programDomyslny,
		Model:         ModelDomyslny,
		Jezyk:         jezykDomyslny,
		KatalogModeli: KatalogModeliDomyslny,
	}
}

// Nanies nanosi na komplet pojedynczy wpis odczytany z konfiguracji i oddaje
// komplet po zmianie.
//
// Klucz spoza kompletu jest pomijany bez błędu. Konfiguracja
// zasięgu, z którego przychodzą wpisy, niesie ustawienia całego produktu —
// kanał modelu, izolację, motyw. Odmowa na każdy nieswój klucz zamieniłaby
// zwykły odczyt konfiguracji w awarię silnika mowy.
//
// Komplet wraca przez wartość, nie przez wskaźnik: wołający składa go
// pętlą po wpisach i nie ma powodu, by dwie takie pętle biegnące obok siebie
// pisały po tej samej strukturze.
func Nanies(u Ustawienia, klucz string, wartosc any) Ustawienia {
	switch klucz {
	case KluczProgram:
		u.Program = tekstUstawienia(wartosc, programDomyslny)
	case KluczModel:
		u.Model = rozmiarModelu(tekstUstawienia(wartosc, ModelDomyslny))
	case KluczJezyk:
		u.Jezyk = tekstUstawienia(wartosc, jezykDomyslny)
	case KluczKatalogModeli:
		u.KatalogModeli = tekstUstawienia(wartosc, KatalogModeliDomyslny)
	}
	return u
}

// rozmiarModelu sprowadza rozmiar do wykazu rozmiaryModelu; wartość spoza
// wykazu wraca jako ModelDomyslny.
//
// SPROWADZENIE, A NIE ODMOWA. Rozmiar nieznany bibliotece kończyłby się
// błędem procesu Pythona — komunikatem cudzego narzędzia, z którego Operator
// nie wyczyta, że winna jest jedna literówka w ustawieniu. Sprowadzenie do
// wartości domyślnej daje transkrypcję zamiast odmowy, a wykaz dopuszczalnych
// wartości Operator i tak widzi w oknie konfiguracji (opcje pozycji katalogu).
func rozmiarModelu(wartosc string) string {
	for _, rozmiar := range rozmiaryModelu {
		if wartosc == rozmiar {
			return wartosc
		}
	}
	return ModelDomyslny
}

// tekstUstawienia wyjmuje z wartości konfiguracji napis. Kształt inny niż napis znaczy
// „ustawienia nie ma" i oddaje wartość domyślną pola, a nie pustkę: pustka
// w polu Jezyk znaczyłaby co innego niż brak wpisu (rozpoznanie automatyczne
// zamiast polskiego), więc każdy kształt niesie tu własną wartość zastępczą.
func tekstUstawienia(wartosc any, domyslna string) string {
	napis, jest := wartosc.(string)
	if !jest {
		return domyslna
	}
	return napis
}
