// Plik wpina katalog okna konfiguracji, kategorie i pozycje katalogu
// ustawień, sterowany danymi bez zdarzenia zmiany.
package core

import (
	"context"

	"danacoconsole/shared"
)

// KatalogUstawien jest portem katalogu okna konfiguracji. Mówi
// wyłącznie typami kontraktu; wiersze tabel leżą po drugiej stronie adaptera.
type KatalogUstawien interface {
	// Kategorie zwraca kategorie okna konfiguracji w kolejności wyświetlania.
	Kategorie(ctx context.Context, z shared.SettingsCategoryListRequest) (shared.SettingsCategoryListResponse, error)
	// Definicje zwraca pozycje katalogu wraz z metadanymi pola formularza.
	Definicje(ctx context.Context, z shared.SettingsDefinitionListRequest) (shared.SettingsDefinitionListResponse, error)
}

// zarejestrujKatalogUstawien wpina obie komendy odczytu katalogu.
//
// Brak portu nie rejestruje niczego: komendy odpowiedzą wtedy `settings.unknown`,
// a pozostałe domeny pracują bez zmian.
func zarejestrujKatalogUstawien(r *Rejestr, katalog KatalogUstawien) {
	if r == nil || katalog == nil {
		return
	}
	r.Zarejestruj(shared.CommandSettingsCategoryList, obsluz(katalog.Kategorie))
	r.Zarejestruj(shared.CommandSettingsDefinitionList, obsluz(katalog.Definicje))
}
