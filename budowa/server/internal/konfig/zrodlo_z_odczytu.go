package konfig

import "errors"

// odczytAdresu zwraca ustawienia zapisane pod adresem złożonym: poziom zasięgu
// razem z osią rozstrzygania. Taki kształt ma odczyt warstwy trwałości —
// metoda ListaOsi repozytorium konfiguracji.
type odczytAdresu func(adres Adres) ([]Wpis, error)

// zrodloZOdczytu składa źródło całej przestrzeni poziom × oś z funkcji odczytu.
type zrodloZOdczytu struct {
	odczytOsi odczytAdresu
}

// NoweZrodloZOdczytuOsi buduje źródło nad funkcją czytającą adres złożony —
// pełna przestrzeń poziom × oś.
func NoweZrodloZOdczytuOsi(odczyt odczytAdresu) *zrodloZOdczytu {
	return &zrodloZOdczytu{odczytOsi: odczyt}
}

// Wpisy odczytuje kolejno adresy kontekstu.
//
// Niepowodzenie odczytu jednego poziomu nie przerywa pozostałych: zebrane wpisy
// wracają w komplecie, a błędy idą obok nich jako informacja diagnostyczna.
// Rozstrzyganie wykona się na tym, co udało się odczytać.
func (z *zrodloZOdczytu) Wpisy(adresy []Adres) ([]Wpis, error) {
	if z == nil || z.odczytOsi == nil {
		return nil, nil
	}
	zebrane := make([]Wpis, 0, len(adresy))
	var bledy []error
	for _, adres := range adresy {
		wpisy, err := z.wpisyAdresu(adres)
		if err != nil {
			bledy = append(bledy, err)
		}
		for _, wpis := range wpisy {
			wpis.Poziom = adres.Poziom
			wpis.KluczZasiegu = adres.KluczZasiegu
			wpis.Os = OsLubPlatforma(adres.Os)
			wpis.KluczOsi = adres.KluczOsi
			zebrane = append(zebrane, wpis)
		}
	}
	return zebrane, errors.Join(bledy...)
}

// wpisyAdresu czyta jeden adres złożony: poziom zasięgu razem z osią.
func (z *zrodloZOdczytu) wpisyAdresu(adres Adres) ([]Wpis, error) {
	return z.odczytOsi(adres)
}
