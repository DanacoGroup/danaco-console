// Odpowiedzialność pliku: przekład wiersza tabeli `konto` na strukturę Konto.
// Wartości słownikowe idą przez przekład z `konta_rotacja.go`, nie przez
// literały.
//
// Kolumna `poswiadczenie_odwolanie` dochodzi tutaj wyłącznie jako znacznik 0/1
// policzony w zapytaniu — struktura nie ma pola na jej treść, więc odczyt
// katalogu nie ma czym wynieść odwołania.
package dane

import (
	"database/sql"
	"fmt"
)

// zbierzKonta odczytuje wszystkie wiersze wyniku.
func zbierzKonta(wiersze *sql.Rows) ([]Konto, error) {
	lista := []Konto{}
	for wiersze.Next() {
		konto, err := odczytajKonto(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, konto)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt katalogu kont: %w", err)
	}
	return lista, nil
}

// odczytajKonto składa strukturę z jednego wiersza wyniku.
func odczytajKonto(wiersz skaner) (Konto, error) {
	var konto Konto
	var rodzaj, stan string
	var identyfikator, model, adres, katalog, wyczerpaneDo sql.NullString
	var poswiadczenie, domyslne, aktywne int
	err := wiersz.Scan(&konto.ID, &konto.Nazwa, &rodzaj, &konto.Dostawca, &identyfikator,
		&model, &adres, &katalog, &poswiadczenie, &stan, &wyczerpaneDo, &domyslne,
		&aktywne, &konto.Kolejnosc, &konto.Utworzono, &konto.Zaktualizowano)
	if err != nil {
		return Konto{}, err
	}
	konto.IdentyfikatorZewnetrzny = tekstZKolumny(identyfikator)
	konto.ModelDomyslny = tekstZKolumny(model)
	konto.AdresBazowy = tekstZKolumny(adres)
	konto.KatalogKonfiguracji = tekstZKolumny(katalog)
	konto.WyczerpaneDo = tekstZKolumny(wyczerpaneDo)
	konto.MaPoswiadczenie = poswiadczenie == 1
	konto.Domyslne = domyslne == 1
	konto.Aktywne = aktywne == 1
	if konto.Rodzaj, err = rodzajKontaZBazy(rodzaj); err != nil {
		return Konto{}, err
	}
	if konto.Stan, err = stanKontaZBazy(stan); err != nil {
		return Konto{}, err
	}
	return konto, nil
}
