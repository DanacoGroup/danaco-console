// Plik wypełnia port KatalogUstawien repozytorium warstwy danych, zawężając
// katalog w rdzeniu do wskazań żądania.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Zgodność adaptera z portem sprawdzana jest przy kompilacji, a nie dopiero
// przy składaniu rdzenia, więc pomyłka typu ujawnia się od razu.
var _ KatalogUstawien = (*adapterKatalogUstawien)(nil)

// adapterKatalogUstawien wypełnia port KatalogUstawien tabelami katalogu,
// kategorii i pozycji, przez repozytorium.
type adapterKatalogUstawien struct {
	repozytorium dane.RepozytoriumKatalogUstawien
}

// nowyAdapterKatalogUstawien wiąże port z repozytorium katalogu, jedynym
// źródłem kategorii i pozycji katalogu.
func nowyAdapterKatalogUstawien(repozytorium dane.RepozytoriumKatalogUstawien) *adapterKatalogUstawien {
	return &adapterKatalogUstawien{repozytorium: repozytorium}
}

// Kategorie zwraca kategorie okna konfiguracji, zawężone kategorią nadrzędną,
// w kolejności wyświetlania.
func (a *adapterKatalogUstawien) Kategorie(ctx context.Context,
	z shared.SettingsCategoryListRequest) (shared.SettingsCategoryListResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.SettingsCategoryListResponse{}, bladBrakuKatalogu("ustawień")
	}
	wiersze, err := a.repozytorium.Kategorie(ctx, tylkoAktywneKatalogu(z.IncludeDisabled))
	if err != nil {
		return shared.SettingsCategoryListResponse{}, err
	}
	kategorie := make([]shared.SettingCategory, 0, len(wiersze))
	for _, kategoria := range wiersze {
		if !pasujeRodzic(kategoria, z.ParentId) {
			continue
		}
		kategorie = append(kategorie, kategoria)
	}
	return shared.SettingsCategoryListResponse{Categories: kategorie}, nil
}

// Definicje zwraca pozycje katalogu zawężone kategorią, kluczem, poziomem
// zasięgu oraz osią rozstrzygania.
func (a *adapterKatalogUstawien) Definicje(ctx context.Context,
	z shared.SettingsDefinitionListRequest) (shared.SettingsDefinitionListResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.SettingsDefinitionListResponse{}, bladBrakuKatalogu("ustawień")
	}
	wiersze, err := a.repozytorium.Definicje(ctx, tylkoAktywneKatalogu(z.IncludeDisabled))
	if err != nil {
		return shared.SettingsDefinitionListResponse{}, err
	}
	pozycje := make([]shared.SettingDefinition, 0, len(wiersze))
	for _, pozycja := range wiersze {
		if pasujeDefinicja(pozycja, z) {
			pozycje = append(pozycje, pozycja)
		}
	}
	return shared.SettingsDefinitionListResponse{Definitions: pozycje, Total: len(pozycje)}, nil
}

// tylkoAktywneKatalogu przekłada pole `includeDisabled` żądania na warunek
// odczytu. Pole pominięte znaczy „bez wierszy wyłączonych" — okno konfiguracji
// pokazuje domyślnie to, co obowiązuje.
func tylkoAktywneKatalogu(zWylaczonymi *bool) bool {
	return zWylaczonymi == nil || !*zWylaczonymi
}

// pasujeRodzic rozstrzyga zawężenie kategorią nadrzędną. Wskazanie puste
// zwraca katalog w całości, zgodnie z opisem pola kontraktu.
func pasujeRodzic(kategoria shared.SettingCategory, rodzic *string) bool {
	if rodzic == nil || *rodzic == "" {
		return true
	}
	return kategoria.ParentId != nil && *kategoria.ParentId == *rodzic
}

// pasujeDefinicja sprawdza pozycję katalogu wobec kompletu zawężeń żądania:
// kategorii, klucza i poziomu.
func pasujeDefinicja(pozycja shared.SettingDefinition, z shared.SettingsDefinitionListRequest) bool {
	if z.CategoryId != nil && *z.CategoryId != "" && pozycja.CategoryId != *z.CategoryId {
		return false
	}
	if z.Key != nil && *z.Key != "" && pozycja.Key != *z.Key {
		return false
	}
	if z.Scope != nil && *z.Scope != "" && !dopuszczaZasieg(pozycja, *z.Scope) {
		return false
	}
	if z.Axis != nil && *z.Axis != "" && !dopuszczaOs(pozycja, *z.Axis) {
		return false
	}
	return true
}

// dopuszczaZasieg odpowiada, czy klucz wolno zapisać na wskazanym poziomie
// zasięgu. Pozycja bez wykazu poziomów nie dopuszcza żadnego —
// katalog wnosi wykaz wprost, więc jego brak jest brakiem zgody, nie zgodą
// domniemaną.
func dopuszczaZasieg(pozycja shared.SettingDefinition, poziom shared.ConfigScope) bool {
	for _, dozwolony := range pozycja.AllowedScopes {
		if dozwolony == poziom {
			return true
		}
	}
	return false
}

// dopuszczaOs odpowiada, czy klucz wolno zapisać dla wskazanej osi. Pusty wykaz
// osi znaczy wyłącznie oś platformy — tak stanowi opis pola `allowedAxes`
// kontraktu.
func dopuszczaOs(pozycja shared.SettingDefinition, os shared.ConfigAxis) bool {
	if len(pozycja.AllowedAxes) == 0 {
		return os == shared.ConfigAxisPlatform
	}
	for _, dozwolona := range pozycja.AllowedAxes {
		if dozwolona == os {
			return true
		}
	}
	return false
}
