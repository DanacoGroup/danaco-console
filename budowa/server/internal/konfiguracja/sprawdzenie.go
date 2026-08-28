package konfiguracja

import (
	"errors"
	"fmt"
	"strings"
)

// Sprawdz weryfikuje spójność ustawień po nałożeniu wszystkich warstw: roli,
// portu i pozostałych pól konfiguracji.
func (k Konfiguracja) Sprawdz() error {
	if _, err := RolaZTekstu(string(k.Rola)); err != nil {
		return err
	}
	if k.Port < 1 || k.Port > 65535 {
		return fmt.Errorf("port %d poza zakresem 1-65535", k.Port)
	}
	if strings.TrimSpace(k.KatalogDanych) == "" {
		return errors.New("katalog danych nie może być pusty")
	}
	return nil
}
