// Odpowiedzialność pliku: dwie zmiany stanu zasobu wniesionego do panelu
// zasobów — oznaczenie ulubionego i usunięcie. Usunięcie zasobu nie kasuje
// bajtów z magazynu.
package core

import (
	"context"

	"danacoconsole/shared"
)

// UstawUlubionyZasob przestawia oznaczenie ulubionego i oddaje zasób odczytany
// po zmianie — obsługuje `design.asset.favorite.set`. Zasób czyta się przed
// zapisem, a odpowiedź niesie stan z bazy, nie echo żądania.
func (a *adapterDesignu) UstawUlubionyZasob(ctx context.Context,
	z shared.DesignAssetFavoriteSetRequest) (shared.DesignAssetFavoriteSetResponse, error) {

	if z.AssetId == "" {
		return shared.DesignAssetFavoriteSetResponse{}, bladWskazaniaDesignu(
			"komenda design.asset.favorite.set bez wskazania zasobu")
	}
	zasob, err := a.repozytorium.Zasob(ctx, z.AssetId)
	if err != nil {
		return shared.DesignAssetFavoriteSetResponse{}, bladNieznanegoZasobuDesignu(z.AssetId, err)
	}
	if err := a.repozytorium.UstawUlubionyZasobu(ctx, zasob.ID, z.Favorite); err != nil {
		return shared.DesignAssetFavoriteSetResponse{}, bladDesignu(err)
	}

	poZmianie, err := a.repozytorium.Zasob(ctx, z.AssetId)
	if err != nil {
		return shared.DesignAssetFavoriteSetResponse{}, bladNieznanegoZasobuDesignu(z.AssetId, err)
	}
	etykiety, err := a.repozytorium.EtykietyZasobu(ctx, poZmianie.ID)
	if err != nil {
		return shared.DesignAssetFavoriteSetResponse{}, bladDesignu(err)
	}
	return shared.DesignAssetFavoriteSetResponse{Asset: zasobKontraktu(poZmianie, etykiety)}, nil
}

// UsunZasob usuwa zasób z panelu zasobów — obsługuje `design.asset.remove`.
// Zasób czyta się przed usunięciem, żeby było co rozgłosić, a wraca wołającemu
// osobnym wyjściem, bo odpowiedź kontraktu niesie samo `removed`. Zasób
// nieznany nie jest odmową.
func (a *adapterDesignu) UsunZasob(ctx context.Context,
	z shared.DesignAssetRemoveRequest) (shared.DesignAssetRemoveResponse, shared.DesignAsset, error) {

	if z.AssetId == "" {
		return shared.DesignAssetRemoveResponse{}, shared.DesignAsset{}, bladWskazaniaDesignu(
			"komenda design.asset.remove bez wskazania zasobu")
	}

	zasob, err := a.repozytorium.Zasob(ctx, z.AssetId)
	if err != nil {
		// Brak wiersza to `removed: false`, nie odmowa.
		if czyBrakZasobuDesignu(err) {
			return shared.DesignAssetRemoveResponse{Removed: false}, shared.DesignAsset{}, nil
		}
		return shared.DesignAssetRemoveResponse{}, shared.DesignAsset{}, bladDesignu(err)
	}
	// Etykiety czyta się jeszcze przed usunięciem, bo znikają kaskadą.
	etykiety, err := a.repozytorium.EtykietyZasobu(ctx, zasob.ID)
	if err != nil {
		return shared.DesignAssetRemoveResponse{}, shared.DesignAsset{}, bladDesignu(err)
	}

	usuniety, err := a.repozytorium.UsunZasob(ctx, z.AssetId)
	if err != nil {
		return shared.DesignAssetRemoveResponse{}, shared.DesignAsset{}, bladDesignu(err)
	}
	return shared.DesignAssetRemoveResponse{Removed: usuniety}, zasobKontraktu(zasob, etykiety), nil
}
