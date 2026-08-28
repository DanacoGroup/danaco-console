// Plik definiuje port macierzy widoczności modułów w środowiskach; kontrakt
// nie niesie dziś komend macierzy, a dane docierają do klienta polami bytów
// nawigacji.
package core

import "context"

// Macierz udostępnia macierz widoczności modułów w środowiskach; nawigacja
// czyta odwzorowanie przez ten port zamiast sięgać po repozytorium wprost.
type Macierz interface {
	// KodySrodowisk zwraca odwzorowanie modułu na kody środowisk, w
	// kolejności kart środowisk.
	KodySrodowisk(ctx context.Context) (map[int64][]string, error)
}

// Asercja rozjazdu portu Macierz z adapterem adapterMacierzy sprawdzana jest
// na etapie kompilacji pliku.
var _ Macierz = (*adapterMacierzy)(nil)
