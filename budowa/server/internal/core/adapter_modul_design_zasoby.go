// Odpowiedzialność pliku: dwie zmiany stanu zasobu wniesionego do Assets Panel —
// oznaczenie ulubionego (`design.asset.favorite.set`) i usunięcie
// (`design.asset.remove`). Metody stoją na `*adapterDesignu`
// (`adapter_modul_design.go`); wniesienie leży w `adapter_modul_design_wgranie.go`,
// etykiety w `adapter_modul_design_etykiety.go` — podział wedle
// odpowiedzialności.
//
// Usunięcie zasobu nie kasuje bajtów z magazynu. Blob leży pod sumą swojej
// zawartości, więc dwa zasoby o identycznej treści (to samo tło wniesione
// w dwóch oknach) dzielą jeden plik; skasowanie go przy usunięciu jednego
// zasobu odebrałoby treść drugiemu, który o niczym nie wie. Zliczania odwołań
// rdzeń nie prowadzi i ten moduł go nie zakłada, więc bajty usuniętego zasobu
// zostają w katalogu danych do czasu, aż magazyn dostanie sprzątanie.
package core

import (
	"context"

	"danacoconsole/shared"
)

// UstawUlubionyZasob przestawia oznaczenie ulubionego i oddaje zasób odczytany
// po zmianie — obsługuje `design.asset.favorite.set`.
//
// Zasób czytamy przed zapisem, bo komenda niesie identyfikator kontraktu,
// a oznaczenie wisi na kluczu wiersza; przy okazji ten sam odczyt odróżnia
// zasób nieznany (odmowa `not_found`, panel ma po czym poznać, że jego wykaz
// jest nieaktualny) od usterki bazy.
//
// Odpowiedź niesie stan z bazy, nie echo żądania — wzorem `UstawEtykietyZasobu`.
// Podstawienie `z.Favorite` do zasobu odczytanego przed zapisem opisywałoby
// stan, którego w bazie może nie być; stąd drugi odczyt.
//
// Ustawienie stanu, który już obowiązuje, jest drogą udaną. Kontrakt nie pyta
// „czy się zmieniło", tylko żąda stanu docelowego — odmowa za powtórzenie
// zabrałaby oknu prawo do wysłania tego, co Operator widzi na przełączniku.
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

// UsunZasob usuwa zasób z Assets Panel — obsługuje `design.asset.remove`.
//
// Zasób czytamy przed usunięciem, żeby było co rozgłosić: `design.asset.changed`
// niesie cały zasób obok rodzaju zmiany (`DesignAssetChangedEvent.Asset`), więc
// po usunięciu wiersza nie ma już z czego złożyć zdarzenia — a zdarzenie
// `deleted` z pustym zasobem powiedziałoby panelowi „coś zniknęło" bez wskazania
// czego. Odczytany zasób wraca wołającemu osobnym wyjściem, bo odpowiedź
// kontraktu niesie samo `removed` (patrz `zarejestrujDesign`).
//
// Zasób nieznany nie jest odmową. Kontrakt pyta wprost, czy zasób istniał i
// został usunięty — `removed: false` odpowiada na to pytanie prawdziwie, a
// odmowa kazałaby oknu obsługiwać błąd tam, gdzie stan docelowy (zasobu nie ma)
// już obowiązuje. Powtórzone usunięcie tego samego kodu też oddaje `false`.
func (a *adapterDesignu) UsunZasob(ctx context.Context,
	z shared.DesignAssetRemoveRequest) (shared.DesignAssetRemoveResponse, shared.DesignAsset, error) {

	if z.AssetId == "" {
		return shared.DesignAssetRemoveResponse{}, shared.DesignAsset{}, bladWskazaniaDesignu(
			"komenda design.asset.remove bez wskazania zasobu")
	}

	zasob, err := a.repozytorium.Zasob(ctx, z.AssetId)
	if err != nil {
		// Brak wiersza to `removed: false`, nie odmowa — stąd rozpoznanie po
		// błędzie braku, a nie odesłanie go dalej.
		if czyBrakZasobuDesignu(err) {
			return shared.DesignAssetRemoveResponse{Removed: false}, shared.DesignAsset{}, nil
		}
		return shared.DesignAssetRemoveResponse{}, shared.DesignAsset{}, bladDesignu(err)
	}
	// Etykiety czytamy jeszcze przed usunięciem: `etykieta_zasobu_design`
	// znika kaskadą razem z wierszem, więc po `UsunZasob`
	// zdarzenie niosłoby zasób bez etykiet, których w chwili usunięcia miał
	// pełny zestaw.
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
