package models

import "context"

// ZrodloDefinicji dostarcza rejestrowi wiersze kanałów. Rejestr nie wie, skąd
// wiersze pochodzą — z tabeli kanal_modelu albo z innego złożenia.
type ZrodloDefinicji interface {
	Definicje(ctx context.Context) ([]Definicja, error)
}
