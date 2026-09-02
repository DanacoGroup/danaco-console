package konfig

import "sync"

// zrodloPamieciowe trzyma ustawienia w pamięci procesu i jest pełnoprawną
// implementacją interfejsu Zrodlo, obsługującą rdzeń uruchomiony bez
// warstwy trwałości oraz nakładkę ustawień ważnych do końca biegu procesu.
type zrodloPamieciowe struct {
	zamek sync.RWMutex
	wpisy map[kluczWpisu]Wpis
}

// NoweZrodloPamieciowe tworzy puste źródło pamięciowe gotowe do zapisu
// i odczytu ustawień w pamięci procesu.
func NoweZrodloPamieciowe() *zrodloPamieciowe {
	return &zrodloPamieciowe{wpisy: make(map[kluczWpisu]Wpis)}
}

// Ustaw zapisuje wartość ustawienia pod wskazanym adresem zasięgu, na osi
// platformy, wywołując metodę ustawWOsi.
func (z *zrodloPamieciowe) Ustaw(poziom Poziom, kluczZasiegu, klucz, wartosc string, rodzaj Rodzaj) {
	z.ustawWOsi(Adres{Poziom: poziom, KluczZasiegu: kluczZasiegu, Os: OsPlatformy},
		klucz, wartosc, rodzaj)
}

// ustawWOsi zapisuje wartość pod adresem złożonym: poziom zasięgu razem
// z osią, zastępując wcześniejszy zapis pod tym samym adresem.
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

// Wpisy zwraca zapisy spod adresów wskazanych na liście, pomijając wpisy
// zapisane pod adresami spoza niej.
func (z *zrodloPamieciowe) Wpisy(_ int64, adresy []Adres) ([]Wpis, error) {
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
