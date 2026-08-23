// Odpowiedzialność pliku: wpięcie katalogu okna konfiguracji — kategorii
// i pozycji katalogu ustawień.
//
// Katalog jest sterowany danymi: nowa pozycja okna konfiguracji to nowy wiersz
// migracji, nie nowa gałąź w rdzeniu. Rdzeń nie zna ani jednego klucza z osobna
// — zna wyłącznie sposób odczytania katalogu, więc całe okno konfiguracji
// obsługuje jeden port, a nie obsługiwacz na ustawienie.
//
// Odczyt katalogu niczego nie zmienia, więc zdarzenia zmiany tu nie ma — wartości
// zmienia rodzina `config.*` i to ona rozgłasza `config.changed`.
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
