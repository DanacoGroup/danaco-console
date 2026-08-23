package konfiguracja

import (
	"os"
	"path/filepath"
)

// nazwaKatalogu to nazwa katalogu danych produktu w profilu użytkownika.
const nazwaKatalogu = "DanacoConsole"

// KatalogDanychDomyslny wskazuje %LOCALAPPDATA%\DanacoConsole.
// Poza Windows sięga po katalog konfiguracyjny użytkownika, aby start nigdy nie zależał
// od obecności jednej zmiennej środowiska.
func KatalogDanychDomyslny() string {
	if lokalny := os.Getenv("LOCALAPPDATA"); lokalny != "" {
		return filepath.Join(lokalny, nazwaKatalogu)
	}
	if bazowy, err := os.UserConfigDir(); err == nil {
		return filepath.Join(bazowy, nazwaKatalogu)
	}
	return filepath.Join(".", nazwaKatalogu)
}

// PrzygotujKatalogDanych zakłada katalog danych wraz z brakującymi katalogami nadrzędnymi.
// Istniejący katalog nie jest zmieniany.
func PrzygotujKatalogDanych(katalog string) error {
	return os.MkdirAll(katalog, 0o755)
}
