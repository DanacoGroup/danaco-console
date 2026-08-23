package core

import "danacoconsole/shared"

// zarejestrujAkcje wpina katalog akcji sterowany danymi: panel akcji
// i narzędzia kanału modelu czerpią z jednego rejestru budowanego
// z wierszy tabeli, nie z listy zaszytej w kodzie. Rdzeń nie zna ani jednej
// akcji — zna wyłącznie sposób jej odczytania.
//
// Nazwa komendy jest parametrem, a nie literałem w tym pliku: podaje ją punkt
// składania stałą pakietu shared (`shared.CommandActionList`), dokładnie jak
// przy pozostałych domenach.
//
// Pusta nazwa albo brak portu nie rejestruje niczego. Komenda odpowie wtedy
// `*.unknown`, a pozostałe domeny pracują bez zmian.
func zarejestrujAkcje(r *Rejestr, akcje Akcje, nazwa shared.MessageType) {
	if r == nil || akcje == nil || nazwa == "" {
		return
	}
	r.Zarejestruj(nazwa, obsluz(akcje.Wykaz))
}
