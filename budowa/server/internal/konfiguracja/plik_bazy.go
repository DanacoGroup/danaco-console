package konfiguracja

import "path/filepath"

// nazwaPlikuBazy to nazwa jedynego pliku trwałości rdzenia, zapisywanego
// wewnątrz katalogu danych produktu.
const nazwaPlikuBazy = "danaco-console.db"

// SciezkaBazy wskazuje plik bazy wewnątrz katalogu danych.
// Zmiana katalogu danych przenosi całą trwałość, bez osobnego przełącznika.
func (k Konfiguracja) SciezkaBazy() string {
	return filepath.Join(k.KatalogDanych, nazwaPlikuBazy)
}
