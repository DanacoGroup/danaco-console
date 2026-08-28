package protocol

import "danacoconsole/shared"

// RejestrKomend zna zbiór dopuszczalnych nazw komend i zdarzeń, wstrzyknięty
// wyłącznie ze stałych wytworzonych do pakietu shared, bez literałów własnych.
type RejestrKomend struct {
	znane map[shared.MessageType]struct{}
}

// NowyRejestrKomend buduje rejestr komend z nazw pochodzących z pakietu
// shared, gotowy do odpytania metodą zna.
func NowyRejestrKomend(nazwy ...shared.MessageType) *RejestrKomend {
	r := &RejestrKomend{znane: make(map[shared.MessageType]struct{}, len(nazwy))}
	r.Dodaj(nazwy...)
	return r
}

// Dodaj rozszerza rejestr o kolejne nazwy — używane, gdy kontrakt składa się
// z kilku grup stałych (komendy, zdarzenia, akcje rejestru).
func (r *RejestrKomend) Dodaj(nazwy ...shared.MessageType) {
	if r == nil {
		return
	}
	if r.znane == nil {
		r.znane = make(map[shared.MessageType]struct{}, len(nazwy))
	}
	for _, nazwa := range nazwy {
		if nazwa == "" {
			continue
		}
		r.znane[nazwa] = struct{}{}
	}
}

// zna odpowiada, czy typ należy do zbioru wstrzykniętego z kontraktu. Rejestr
// pusty albo nieustawiony nie zna niczego — każdy typ trafi wtedy na ścieżkę
// `*.unknown`, co jest zachowaniem fail-open, nie blokadą.
func (r *RejestrKomend) zna(typ shared.MessageType) bool {
	if r == nil || r.znane == nil {
		return false
	}
	_, jest := r.znane[typ]
	return jest
}

// Liczba zwraca rozmiar rejestru komend znanych, do diagnostyki startu
// rdzenia i porównania z kontraktem.
func (r *RejestrKomend) Liczba() int {
	if r == nil {
		return 0
	}
	return len(r.znane)
}

// Rozpoznaj zwraca nazwę do dalszego kierowania oraz informację, czy typ jest
// znany. Nieznany typ nie jest błędem zrywającym — zamienia się na zdarzenie
// `*.unknown` swojego obszaru, wskazane przez kontrakt.
func (r *RejestrKomend) Rozpoznaj(typ shared.MessageType) (shared.MessageType, bool) {
	if r.zna(typ) {
		return typ, true
	}
	return shared.ZdarzenieNieznanej(typ), false
}
