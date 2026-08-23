package konfiguracja

// Wczytaj ustala konfigurację rdzenia w kolejności warstw:
// wartości domyślne → zmienne środowiska → argumenty wywołania.
// Odczyt środowiska jest parametrem, więc wczytanie daje się sprawdzić bez zmiany procesu.
// Zwrócenie flag.ErrHelp oznacza żądanie pomocy, nie błąd konfiguracji.
func Wczytaj(argumenty []string, odczytSrodowiska func(string) string) (Konfiguracja, error) {
	kon := Domyslna()
	if err := zastosujSrodowisko(&kon, odczytSrodowiska); err != nil {
		return Konfiguracja{}, err
	}
	if err := zastosujArgumenty(&kon, argumenty); err != nil {
		return Konfiguracja{}, err
	}
	if err := kon.Sprawdz(); err != nil {
		return Konfiguracja{}, err
	}
	return kon, nil
}
