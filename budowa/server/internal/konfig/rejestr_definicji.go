package konfig

import "sync"

// Definicja opisuje jedno ustawienie: co ustawia, wartość domyślną
// i objaśnienie pokazywane operatorowi. Miejsce pozycji w oknie konfiguracji
// i dopuszczalne adresy wypełnia katalog ustawień z bazy; definicja
// zbudowana w kodzie zostawia je puste.
type Definicja struct {
	Klucz       string // ustawienie.klucz
	Domyslna    string // wartość obowiązująca przy braku zapisu na każdym z ośmiu poziomów
	Rodzaj      Rodzaj // postać wartości
	Objasnienie string // opis dla Operatora

	Kategoria        string   // kategoria_ustawien.kod
	Nazwa            string   // etykieta pola formularza
	DozwoloneZasiegi []Poziom // poziomy, na których wolno zapisać klucz
	DozwoloneOsie    []Os     // osie, dla których wolno zapisać klucz
	Kolejnosc        int      // kolejność pola w kategorii
	Aktywna          bool     // czy wiersz katalogu jest czynny
}

// Rejestr zna definicje ustawień. Nie jest bramą: klucz spoza rejestru nadal
// daje się odczytać i zapisać, tylko bez wartości domyślnej i objaśnienia.
type Rejestr struct {
	zamek     sync.RWMutex
	wgKlucza  map[string]Definicja
	kolejnosc []string
}

// NowyRejestr buduje rejestr z podanych definicji przez wywołanie Dodaj,
// zachowując ich kolejność dodania.
func NowyRejestr(definicje ...Definicja) *Rejestr {
	rejestr := &Rejestr{wgKlucza: make(map[string]Definicja, len(definicje))}
	rejestr.Dodaj(definicje...)
	return rejestr
}

// Dodaj rozszerza rejestr o kolejne definicje. Ponowne podanie klucza nadpisuje
// definicję bez zmiany miejsca w kolejności.
func (r *Rejestr) Dodaj(definicje ...Definicja) {
	if r == nil {
		return
	}
	r.zamek.Lock()
	defer r.zamek.Unlock()
	if r.wgKlucza == nil {
		r.wgKlucza = make(map[string]Definicja, len(definicje))
	}
	for _, definicja := range definicje {
		if definicja.Klucz == "" {
			continue
		}
		definicja.Rodzaj = RodzajLubTekst(definicja.Rodzaj)
		if _, jest := r.wgKlucza[definicja.Klucz]; !jest {
			r.kolejnosc = append(r.kolejnosc, definicja.Klucz)
		}
		r.wgKlucza[definicja.Klucz] = definicja
	}
}

// Definicja zwraca definicję zarejestrowaną pod danym kluczem oraz
// informację, czy klucz w ogóle istnieje w rejestrze.
func (r *Rejestr) Definicja(klucz string) (Definicja, bool) {
	if r == nil {
		return Definicja{}, false
	}
	r.zamek.RLock()
	defer r.zamek.RUnlock()
	definicja, jest := r.wgKlucza[klucz]
	return definicja, jest
}

// Klucze zwraca wszystkie klucze zarejestrowane w rejestrze, w kolejności,
// w jakiej zostały do niego dodane.
func (r *Rejestr) Klucze() []string {
	if r == nil {
		return nil
	}
	r.zamek.RLock()
	defer r.zamek.RUnlock()
	kopia := make([]string, len(r.kolejnosc))
	copy(kopia, r.kolejnosc)
	return kopia
}

// Definicje zwraca komplet zarejestrowanych definicji w kolejności, w jakiej
// zostały dodane do rejestru.
func (r *Rejestr) Definicje() []Definicja {
	if r == nil {
		return nil
	}
	r.zamek.RLock()
	defer r.zamek.RUnlock()
	komplet := make([]Definicja, 0, len(r.kolejnosc))
	for _, klucz := range r.kolejnosc {
		komplet = append(komplet, r.wgKlucza[klucz])
	}
	return komplet
}

// Liczba zwraca liczbę zdefiniowanych w rejestrze ustawień, wykorzystywaną
// przy diagnostyce startu rdzenia.
func (r *Rejestr) Liczba() int {
	if r == nil {
		return 0
	}
	r.zamek.RLock()
	defer r.zamek.RUnlock()
	return len(r.kolejnosc)
}

// RejestrWbudowany zwraca komplet ustawień znanych platformie w chwili
// wydania: parametry wykonania okna komunikacji oraz jedenaście punktów
// izolacji. Rejestr jest zbiorem otwartym, rozszerzenie o kolejne ustawienie
// nie wymaga zmiany rozstrzygania.
func RejestrWbudowany() *Rejestr {
	rejestr := NowyRejestr(definicjeWykonania()...)
	rejestr.Dodaj(definicjeIzolacji()...)
	rejestr.Dodaj(definicjeAplikacji()...)
	return rejestr
}
