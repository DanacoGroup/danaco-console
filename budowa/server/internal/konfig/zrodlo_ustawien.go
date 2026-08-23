package konfig

// Wpis to jedno ustawienie zapisane pod konkretnym adresem zasięgu.
// Odpowiada wierszowi tabeli `ustawienie`.
type Wpis struct {
	Poziom       Poziom
	KluczZasiegu string
	// Os i KluczOsi wskazują oś rozstrzygania. Oś pusta znaczy platformę,
	// więc wpis zbudowany bez wskazania osi obowiązuje na osi platformy.
	Os       Os
	KluczOsi string
	Klucz    string
	Wartosc  string
	Rodzaj   Rodzaj
}

// Adres zwraca miejsce zapisu wpisu.
func (w Wpis) Adres() Adres {
	return Adres{
		Poziom:       w.Poziom,
		KluczZasiegu: w.KluczZasiegu,
		Os:           OsLubPlatforma(w.Os),
		KluczOsi:     w.KluczOsi,
	}
}

// Zrodlo dostarcza zapisane ustawienia dla wskazanych adresów zasięgu.
//
// Pakiet konfig nie sięga do bazy danych samodzielnie — łączy się z warstwą
// trwałości wyłącznie przez ten interfejs. Implementacja czytająca
// tabelę `ustawienie` należy do warstwy repozytoriów; implementacja pamięciowa
// z tego pakietu obsługuje pracę bez trwałości i sprawdzenia.
//
// Zwrócony błąd nie zatrzymuje rozstrzygania: wpisy zwrócone mimo błędu wchodzą
// do rozstrzygnięcia, a ustawienia bez zapisu schodzą na wartości domyślne.
// Błąd jest widoczny w polityce efektywnej jako informacja
// diagnostyczna, nie jako odmowa. Implementacja może więc zwrócić wynik
// częściowy razem z błędem.
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
