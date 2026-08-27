// Plik przekłada wiersz rejestru kont na strukturę kontraktu oraz nakłada na
// wiersz pola żądania zmiany, dotykając poświadczenia wyłącznie znacznikiem
// jego obecności.
package core

import (
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// kontoKontraktu przekłada wiersz rejestru kont na strukturę Account
// kontraktu, z polami identyfikatora, dostawcy i stanu aktywności.
func kontoKontraktu(k dane.Konto) shared.Account {
	return shared.Account{
		Id:            strconv.FormatInt(k.ID, 10),
		Name:          k.Nazwa,
		Kind:          k.Rodzaj,
		Provider:      k.Dostawca,
		ExternalId:    k.IdentyfikatorZewnetrzny,
		DefaultModel:  k.ModelDomyslny,
		BaseUrl:       k.AdresBazowy,
		ConfigDir:     k.KatalogKonfiguracji,
		HasCredential: k.MaPoswiadczenie,
		IsDefault:     k.Domyslne,
		Enabled:       k.Aktywne,
		Order:         k.Kolejnosc,
		CreatedAt:     chwilaBazy(k.Utworzono),
		UpdatedAt:     chwilaBazy(k.Zaktualizowano),
	}
}

// zastosujZmianeKonta nakłada na wiersz pola wskazane w żądaniu; pole
// pominięte zostaje bez zmiany, a rodzaj konta pozostaje niezmienny, bo
// rozstrzyga o puli rotacji.
func zastosujZmianeKonta(konto *dane.Konto, z shared.AccountUpdateRequest) {
	if z.Name != nil {
		konto.Nazwa = *z.Name
	}
	if z.Provider != nil {
		konto.Dostawca = *z.Provider
	}
	if z.ExternalId != nil {
		konto.IdentyfikatorZewnetrzny = z.ExternalId
	}
	if z.DefaultModel != nil {
		konto.ModelDomyslny = z.DefaultModel
	}
	if z.BaseUrl != nil {
		konto.AdresBazowy = z.BaseUrl
	}
	if z.ConfigDir != nil {
		konto.KatalogKonfiguracji = z.ConfigDir
	}
	if z.Enabled != nil {
		konto.Aktywne = *z.Enabled
	}
}

// numeryTekstem przekłada numery wierszy na identyfikatory kontraktu. Wykaz
// pusty daje brak wartości, a nie tablicę pustą — pole jest opcjonalne.
func numeryTekstem(numery []int64) []string {
	if len(numery) == 0 {
		return nil
	}
	teksty := make([]string, 0, len(numery))
	for _, numer := range numery {
		teksty = append(teksty, strconv.FormatInt(numer, 10))
	}
	return teksty
}
