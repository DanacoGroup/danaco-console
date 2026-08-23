package konfig

import "sync"

// zrodloPamieciowe trzyma ustawienia w pamięci procesu. Jest pełnoprawną
// implementacją interfejsu Zrodlo — obsługuje rdzeń uruchomiony bez warstwy
// trwałości oraz nakładkę ustawień ważnych do końca biegu procesu.
//
// Rozstrzygacz zbudowany bez źródła sięga po puste źródło pamięciowe, dzięki
// czemu brak trwałości nie blokuje startu, tylko daje politykę domyślną.
type zrodloPamieciowe struct {
	zamek sync.RWMutex
	wpisy map[kluczWpisu]Wpis
}

// NoweZrodloPamieciowe tworzy puste źródło.
func NoweZrodloPamieciowe() *zrodloPamieciowe {
	return &zrodloPamieciowe{wpisy: make(map[kluczWpisu]Wpis)}
}

// Ustaw zapisuje wartość ustawienia pod wskazanym adresem zasięgu, na osi
// platformy.
func (z *zrodloPamieciowe) Ustaw(poziom Poziom, kluczZasiegu, klucz, wartosc string, rodzaj Rodzaj) {
	z.ustawWOsi(Adres{Poziom: poziom, KluczZasiegu: kluczZasiegu, Os: OsPlatformy},
		klucz, wartosc, rodzaj)
}

// ustawWOsi zapisuje wartość pod adresem złożonym: poziom zasięgu razem z osią.
func (z *zrodloPamieciowe) ustawWOsi(adres Adres, klucz, wartosc string, rodzaj Rodzaj) {
	if z == nil || klucz == "" {
		return
	}
	adres.Os = OsLubPlatforma(adres.Os)
	z.zamek.Lock()
	defer z.zamek.Unlock()
	if z.wpisy == nil {
		z.wpisy = make(map[kluczWpisu]Wpis)
	}
	z.wpisy[kluczem(adres, klucz)] = Wpis{
		Poziom:       adres.Poziom,
		KluczZasiegu: adres.KluczZasiegu,
		Os:           adres.Os,
		KluczOsi:     adres.KluczOsi,
		Klucz:        klucz,
		Wartosc:      wartosc,
		Rodzaj:       RodzajLubTekst(rodzaj),
	}
}

// Usun kasuje zapis spod adresu. Skasowanie przywraca wartość poziomu szerszego,
// a w ostateczności wartość domyślną.
func (z *zrodloPamieciowe) Usun(poziom Poziom, kluczZasiegu, klucz string) {
	if z == nil {
		return
	}
	z.zamek.Lock()
	defer z.zamek.Unlock()
	delete(z.wpisy, kluczem(Adres{Poziom: poziom, KluczZasiegu: kluczZasiegu}, klucz))
}

// Wpisy zwraca zapisy spod wskazanych adresów.
func (z *zrodloPamieciowe) Wpisy(adresy []Adres) ([]Wpis, error) {
	if z == nil {
		return nil, nil
	}
	z.zamek.RLock()
	defer z.zamek.RUnlock()
	wybrane := make([]Wpis, 0, len(z.wpisy))
	for _, wpis := range z.wpisy {
		if naliscie(adresy, wpis.Adres()) {
			wybrane = append(wybrane, wpis)
		}
	}
	return wybrane, nil
}

func naliscie(adresy []Adres, szukany Adres) bool {
	for _, adres := range adresy {
		if adres == szukany {
			return true
		}
	}
	return false
}
