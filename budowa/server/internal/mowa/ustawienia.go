// Odpowiedzialność pliku: cztery ustawienia silnika mowy — klucze katalogu,
// wartości domyślne i nanoszenie odczytanej konfiguracji na komplet. Ten plik
// niczego nie czyta z bazy — gotowe pary klucz-wartość przychodzą z warstwy wyżej.
package mowa

// Klucze katalogu ustawień sterujące silnikiem mowy. Źródło:
// `server/internal/store/migracja_075_mowa.sql` (kolumna
// `definicja_ustawienia.klucz`).
const (
	// KluczProgram — ścieżka interpretera Pythona uruchamiającego pomocnika.
	// Pusta znaczy „szukaj python3 na ścieżce wyszukiwania systemu".
	KluczProgram = "mowa_program"
	// KluczModel — rozmiar modelu rozpoznawania; dozwolony wykaz w
	// rozmiaryModelu, poza wykazem wraca ModelDomyslny.
	KluczModel = "mowa_model"
	// KluczJezyk — kod języka rozpoznawania mowy, dwuliterowy zapis ISO 639-1,
	// na przykład pl dla polskiego.
	KluczJezyk = "mowa_jezyk"
	// KluczKatalogModeli — katalog wag modelu. Pusty znaczy pamięć podręczna
	// biblioteki w katalogu domowym konta.
	KluczKatalogModeli = "mowa_katalog_modeli"
)

// Wartości domyślne. Odpowiadają kolumnie definicja_ustawienia.wartosc_domyslna.
// Brak wiersza w tabeli ustawienie znaczy właśnie tę wartość.
const (
	// programDomyslny jest pusty; wyszukanie interpretera idzie automatycznie
	// po ścieżce wyszukiwania systemu operacyjnego.
	programDomyslny = ""
	// ModelDomyslny to średni rozmiar: równowaga szybkości i dokładności
	// rozpoznawania mowy przez bibliotekę.
	ModelDomyslny = "small"
	// jezykDomyslny to polski — produkt jest polskojęzyczny domyślnie, dla
	// Operatora, który nie ustawił niczego.
	jezykDomyslny = "pl"
	// KatalogModeliDomyslny — katalog wag rozłożonych na maszynie obok rdzenia.
	// Wartość niepusta, bo pusta znaczyłaby katalog domowy konta uruchamiającego.
	KatalogModeliDomyslny = "/opt/danaco-modele/mowa"
)

// rozmiaryModelu wylicza dopuszczalne rozmiary modelu w kolejności rosnącej
// wierności. Wartości pustej w wykazie nie ma — silnik musi dostać nazwę
// modelu, bo nie ma warstwy niżej, która by rozmiar rozstrzygnęła.
var rozmiaryModelu = []string{"tiny", "base", "small", "medium", "large-v3"}

// Ustawienia to komplet czterech nastaw silnika mowy. Jest to struktura,
// a nie mapa: pola nazwane rozstrzygają się przy kompilacji.
type Ustawienia struct {
	// Program — ścieżka interpretera; pusta znaczy wyszukanie w systemie.
	Program string
	// Model — rozmiar modelu z wykazu rozmiaryModelu.
	Model string
	// Jezyk — kod języka rozpoznawania, dwuliterowy ISO 639-1.
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
// komplet po zmianie. Klucz spoza kompletu jest pomijany bez błędu. Komplet
// wraca przez wartość, nie przez wskaźnik.
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
// wykazu wraca jako ModelDomyslny, zamiast odmowy.
func rozmiarModelu(wartosc string) string {
	for _, rozmiar := range rozmiaryModelu {
		if wartosc == rozmiar {
			return wartosc
		}
	}
	return ModelDomyslny
}

// tekstUstawienia wyjmuje z wartości konfiguracji napis. Kształt inny niż
// napis znaczy brak ustawienia i oddaje wartość domyślną pola, nie pustkę.
func tekstUstawienia(wartosc any, domyslna string) string {
	napis, jest := wartosc.(string)
	if !jest {
		return domyslna
	}
	return napis
}
