// Historia wersji dokumentu Studio (Repository Panel) — wykaz i przywracanie.
// Dokument i jego zapis leżą w `adapter_modul_studio.go`; ten plik implementuje
// wyłącznie metody na tym samym `adapterStudia`, zadeklarowanym tam.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Wersje zwraca historię wersji dokumentu od najnowszej. `limit` z żądania
// przycina wykaz po stronie adaptera — warstwa danych oddaje całą historię,
// bo Repository Panel bywa też otwierany bez granicy (przewijanie w dół).
func (a *adapterStudia) Wersje(ctx context.Context,
	z shared.StudioRepositoryListRequest) (shared.StudioRepositoryListResponse, error) {

	if z.DocumentId == "" {
		return shared.StudioRepositoryListResponse{}, bladWskazaniaStudio("wykaz wersji bez wskazania dokumentu")
	}
	dokument, err := a.repozytorium.Dokument(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioRepositoryListResponse{}, bladNieznanegoDokumentu(z.DocumentId, err)
	}
	wiersze, err := a.repozytorium.Wersje(ctx, dokument.ID)
	if err != nil {
		return shared.StudioRepositoryListResponse{}, bladStudio(err)
	}
	if z.Limit != nil && *z.Limit >= 0 && *z.Limit < len(wiersze) {
		wiersze = wiersze[:*z.Limit]
	}
	wersje := make([]shared.StudioVersion, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wersje = append(wersje, a.zlozWersje(wiersz))
	}
	return shared.StudioRepositoryListResponse{Versions: wersje}, nil
}

// PrzywrocWersje przywraca treść wcześniejszej wersji jako bieżącą, bez
// usuwania wersji nowszych; transakcja `dane.PrzywrocWersje` przestawia
// `wersja_biezaca_id` na wersję wskazaną w żądaniu.
func (a *adapterStudia) PrzywrocWersje(ctx context.Context,
	z shared.StudioRepositoryRestoreRequest) (shared.StudioRepositoryRestoreResponse, error) {

	if z.DocumentId == "" || z.VersionId == "" {
		return shared.StudioRepositoryRestoreResponse{}, bladWskazaniaStudio(
			"przywrócenie wersji wymaga dokumentu i wersji")
	}
	dokument, err := a.repozytorium.PrzywrocWersje(ctx, z.DocumentId, z.VersionId)
	if err != nil {
		return shared.StudioRepositoryRestoreResponse{}, bladNieznanegoDokumentuAlboWersji(
			z.DocumentId, z.VersionId, err)
	}
	return shared.StudioRepositoryRestoreResponse{Document: a.zlozDokument(dokument)}, nil
}

// bladNieznanegoDokumentuAlboWersji nazywa brak bytu przy przywracaniu.
// `dane.PrzywrocWersje` oddaje ten sam `ErrBrakWiersza` zarówno dla braku
// dokumentu, jak i dla braku wersji, więc komunikat wymienia oba
// identyfikatory naraz.
func bladNieznanegoDokumentuAlboWersji(kodDokumentu, kodWersji string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Studio: dokument "+kodDokumentu+" albo wersja "+kodWersji+" nie istnieje"))
	}
	return bladStudio(err)
}
