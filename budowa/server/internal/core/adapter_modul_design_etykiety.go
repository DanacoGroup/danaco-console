// Odpowiedzialność pliku: nadawanie etykiet zasobom panelu zasobów —
// `design.asset.tag.set` na typie `*adapterDesignu`. Pusty zestaw etykiet
// jest wartością, nie brakiem żądania: `tags: []` znaczy zdejmij wszystkie.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// UstawEtykietyZasobu podmienia komplet etykiet zasobu na nadesłany —
// obsługuje `design.asset.tag.set`, jedną transakcją warstwy danych.
func (a *adapterDesignu) UstawEtykietyZasobu(ctx context.Context,
	z shared.DesignAssetTagSetRequest) (shared.DesignAssetTagSetResponse, error) {

	if z.AssetId == "" {
		return shared.DesignAssetTagSetResponse{}, bladWskazaniaDesignu("komenda bez wskazania zasobu")
	}
	zasob, err := a.repozytorium.Zasob(ctx, z.AssetId)
	if err != nil {
		return shared.DesignAssetTagSetResponse{}, bladNieznanegoZasobuDesignu(z.AssetId, err)
	}
	if err := a.repozytorium.UstawEtykietyZasobu(ctx, zasob.ID, z.Tags); err != nil {
		return shared.DesignAssetTagSetResponse{}, bladDesignu(err)
	}
	etykiety, err := a.repozytorium.EtykietyZasobu(ctx, zasob.ID)
	if err != nil {
		return shared.DesignAssetTagSetResponse{}, bladDesignu(err)
	}
	return shared.DesignAssetTagSetResponse{Asset: zasobKontraktu(zasob, etykiety)}, nil
}

// czyBrakZasobuDesignu odróżnia „zasobu nie ma” od usterki odczytu bez
// zamieniania tego w błąd, regułą warstwy danych (`ErrBrakWiersza`).
func czyBrakZasobuDesignu(err error) bool {
	return errors.Is(err, dane.ErrBrakWiersza)
}

// bladNieznanegoZasobuDesignu odróżnia „zasobu nie ma” od „odczyt się nie
// powiódł”. Assets Panel ma po tym poznać, że wskazanie w wykazie jest
// nieaktualne, zamiast czytać usterkę rdzenia tam, gdzie jej nie było.
func bladNieznanegoZasobuDesignu(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Design: zasób nie istnieje: "+kod))
	}
	return bladDesignu(err)
}
