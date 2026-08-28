package konfig

// Wpis to jedno ustawienie zapisane pod konkretnym adresem zasięgu.
// Odpowiada wierszowi tabeli `ustawienie`.
type Wpis struct {
	Poziom       Poziom
	KluczZasiegu string
	// Os i KluczOsi wskazują oś rozstrzygania; oś pusta oznacza platformę.
	Os       Os
	KluczOsi string
	Klucz    string
	Wartosc  string
	Rodzaj   Rodzaj
}

// Adres zwraca miejsce zapisu wpisu jako parę poziomu zasięgu i osi wraz
// z ich bytami, gotową do porównań.
func (w Wpis) Adres() Adres {
	return Adres{
		Poziom:       w.Poziom,
		KluczZasiegu: w.KluczZasiegu,
		Os:           OsLubPlatforma(w.Os),
		KluczOsi:     w.KluczOsi,
	}
}

// Zrodlo dostarcza zapisane ustawienia dla wskazanych adresów zasięgu.
// Pakiet konfig nie sięga do bazy samodzielnie, tylko łączy się z warstwą
// trwałości przez ten interfejs. Zwrócony błąd nie zatrzymuje
// rozstrzygania.
type Zrodlo interface {
	Wpisy(adresy []Adres) ([]Wpis, error)
}

// kluczWpisu jednoznacznie wskazuje wiersz ustawienia — odpowiednik więzu
// UNIQUE(poziom_zasiegu_id, klucz_zasiegu, os, klucz_osi, klucz).
type kluczWpisu struct {
	adres Adres
	klucz string
}

func kluczem(adres Adres, klucz string) kluczWpisu {
	adres.Os = OsLubPlatforma(adres.Os)
	return kluczWpisu{adres: adres, klucz: klucz}
}
