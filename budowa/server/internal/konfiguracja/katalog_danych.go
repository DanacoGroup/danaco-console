package konfiguracja

import (
	"os"
	"path/filepath"
)

// nazwaKatalogu to nazwa katalogu danych produktu w profilu użytkownika,
// wspólna dla wszystkich platform.
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

// prawaKatalogu odcinają grupę i pozostałych. W katalogu danych leży plik bazy
// z rozmowami i skrótami sesji bramki oraz sejf poświadczeń, więc prawo wejścia
// ma wyłącznie właściciel procesu rdzenia.
const prawaKatalogu = 0o700

// PrzygotujKatalogDanych zakłada katalog danych wraz z brakującymi katalogami nadrzędnymi
// i zawęża prawa samego katalogu danych — także wtedy, gdy zastał go szerzej otwartym.
func PrzygotujKatalogDanych(katalog string) error {
	if err := os.MkdirAll(katalog, prawaKatalogu); err != nil {
		return err
	}
	return os.Chmod(katalog, prawaKatalogu)
}
