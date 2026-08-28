package konfig

import "danacoconsole/shared"

// WpisKontraktu zamienia rozstrzygnięcie na wpis konfiguracji kontraktu.
// Pole scope niesie poziom, z którego pochodzi wartość; dla wartości
// domyślnej poziom pozostaje pusty, ponieważ wartość nie pochodzi z żadnego
// zapisu.
func (w Wynik) WpisKontraktu() shared.ConfigEntry {
	wpis := shared.ConfigEntry{
		Key:   w.Klucz,
		Value: KodujJSON(w.Wartosc, w.Rodzaj),
		Scope: w.Poziom,
	}
	if w.KluczZasiegu != "" {
		klucz := w.KluczZasiegu
		wpis.ScopeId = &klucz
	}
	// Oś platformy nie jedzie w kopercie: brak osi znaczy właśnie platformę.
	if OsLubPlatforma(w.Os) != OsPlatformy {
		os := OsLubPlatforma(w.Os)
		wpis.Axis = &os
		if w.KluczOsi != "" {
			kluczOsi := w.KluczOsi
			wpis.AxisId = &kluczOsi
		}
	}
	return wpis
}

// WpisyKontraktu zamienia komplet polityki efektywnej na wpisy kontraktu —
// gotowa treść odpowiedzi na komendę odczytu konfiguracji.
func (p Polityka) WpisyKontraktu() []shared.ConfigEntry {
	wpisy := make([]shared.ConfigEntry, 0, len(p.Pozycje))
	for _, pozycja := range p.Pozycje {
		wpisy = append(wpisy, pozycja.WpisKontraktu())
	}
	return wpisy
}

// ZapisZKontraktuWOsi zamienia żądanie zapisu na wpis pod adresem złożonym:
// poziom zasięgu razem z osią (pola axis i axisId żądania). Oś pusta znaczy
// platformę; oś spoza kontraktu daje false, tak samo jak poziom spoza kontraktu.
func ZapisZKontraktuWOsi(klucz string, wartosc []byte, adres Adres) (Wpis, bool) {
	if klucz == "" || !Znany(adres.Poziom) || !ZnanaOs(adres.Os) {
		return Wpis{}, false
	}
	tresc, rodzaj := DekodujJSON(wartosc)
	return Wpis{
		Poziom:       adres.Poziom,
		KluczZasiegu: adres.KluczZasiegu,
		Os:           OsLubPlatforma(adres.Os),
		KluczOsi:     adres.KluczOsi,
		Klucz:        klucz,
		Wartosc:      tresc,
		Rodzaj:       rodzaj,
	}, true
}
