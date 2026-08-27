// Odpowiedzialność pliku: nazwane wersje kompozycji Design Board —
// `design.board.version.save`, `design.board.version.list`,
// `design.board.version.restore`.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekWersjiKompozycjiDesign znakuje identyfikatory zewnętrzne wersji
// kompozycji Design Board Panel.
const przedrostekWersjiKompozycjiDesign = "wersja-planszy-"

// ZapiszWersjeKompozycji utrwala bieżący układ kompozycji jako nazwaną wersję —
// obsługuje `design.board.version.save`. Migawka schodzi z bazy, nie z
// żądania.
func (a *adapterDesignu) ZapiszWersjeKompozycji(ctx context.Context,
	z shared.DesignBoardVersionSaveRequest) (shared.DesignBoardVersionSaveResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignBoardVersionSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.board.version.save bez wskazania kompozycji")
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignBoardVersionSaveResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}

	wersja, err := a.odlozWersjeKompozycjiDesignu(ctx, kompozycja, z.Name, z.Note)
	if err != nil {
		return shared.DesignBoardVersionSaveResponse{}, err
	}
	return shared.DesignBoardVersionSaveResponse{Version: wersja}, nil
}

// WersjeKompozycji zwraca wersje kompozycji od najświeższej — obsługuje
// `design.board.version.list`. Warstw wykaz nie niesie, stąd `LayerCount`
// z wiersza wersji.
func (a *adapterDesignu) WersjeKompozycji(ctx context.Context,
	z shared.DesignBoardVersionListRequest) (shared.DesignBoardVersionListResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignBoardVersionListResponse{}, bladWskazaniaDesignu(
			"komenda design.board.version.list bez wskazania kompozycji")
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignBoardVersionListResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}
	wiersze, razem, err := a.repozytorium.WersjeKompozycjiDesignu(ctx, kompozycja.ID,
		wartoscLiczby(z.Limit))
	if err != nil {
		return shared.DesignBoardVersionListResponse{}, bladDesignu(err)
	}
	wersje := make([]shared.DesignBoardVersion, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wersje = append(wersje, wersjaKompozycjiKontraktuDesignu(wiersz, kompozycja.Kod, nil))
	}
	return shared.DesignBoardVersionListResponse{Versions: wersje, Total: razem}, nil
}

// PrzywrocWersjeKompozycji przywraca układ kompozycji z wersji — obsługuje
// `design.board.version.restore`. Kolejność ma znaczenie: najpierw wersja
// z układu sprzed przywrócenia, dopiero potem nadpisanie warstw.
func (a *adapterDesignu) PrzywrocWersjeKompozycji(ctx context.Context,
	z shared.DesignBoardVersionRestoreRequest) (shared.DesignBoardVersionRestoreResponse, error) {

	if strings.TrimSpace(z.VersionId) == "" {
		return shared.DesignBoardVersionRestoreResponse{}, bladWskazaniaDesignu(
			"komenda design.board.version.restore bez wskazania wersji")
	}
	wersja, err := a.repozytorium.WersjaKompozycjiDesignuPoKodzie(ctx, strings.TrimSpace(z.VersionId))
	if err != nil {
		if czyBrakZasobuDesignu(err) {
			return shared.DesignBoardVersionRestoreResponse{}, bladNieznanegoBytuDesignu(
				"wersji kompozycji " + z.VersionId + " nie ma w module Design")
		}
		return shared.DesignBoardVersionRestoreResponse{}, bladDesignu(err)
	}
	kompozycja, err := a.repozytorium.KompozycjaDesignuPoKluczu(ctx, wersja.KompozycjaID)
	if err != nil {
		return shared.DesignBoardVersionRestoreResponse{}, bladDesignu(err)
	}

	odlozona, err := a.odlozWersjeKompozycjiDesignu(ctx, kompozycja,
		wskaznikTekstuDesignu("układ sprzed przywrócenia wersji "+wersja.Kod), nil)
	if err != nil {
		return shared.DesignBoardVersionRestoreResponse{}, err
	}

	warstwy, err := a.repozytorium.WarstwyWersjiKompozycjiDesignu(ctx, wersja.ID)
	if err != nil {
		return shared.DesignBoardVersionRestoreResponse{}, bladDesignu(err)
	}
	zapisana, err := a.repozytorium.ZapiszKompozycje(ctx, dane.KompozycjaDesignu{
		Kod: kompozycja.Kod, Okno: kompozycja.Okno, Nazwa: kompozycja.Nazwa,
	}, warstwy)
	if err != nil {
		return shared.DesignBoardVersionRestoreResponse{}, bladDesignu(err)
	}
	board, err := a.zlozBoard(ctx, zapisana)
	if err != nil {
		return shared.DesignBoardVersionRestoreResponse{}, bladDesignu(err)
	}
	return shared.DesignBoardVersionRestoreResponse{Board: board, SupersededVersion: &odlozona}, nil
}

// odlozWersjeKompozycjiDesignu zdejmuje migawkę bieżącego układu kompozycji
// i utrwala ją jako wersję. Wspólna droga zapisu i przywrócenia. Nazwa pusta
// bierze znacznik czasu.
func (a *adapterDesignu) odlozWersjeKompozycjiDesignu(ctx context.Context,
	kompozycja dane.KompozycjaDesignu, nazwa, uzasadnienie *string) (shared.DesignBoardVersion, error) {

	warstwy, err := a.repozytorium.Warstwy(ctx, kompozycja.ID)
	if err != nil {
		return shared.DesignBoardVersion{}, bladDesignu(err)
	}
	zapisana, err := a.repozytorium.ZapiszWersjeKompozycjiDesignu(ctx, dane.WersjaKompozycjiDesignu{
		Kod:          nowyIdentyfikator(przedrostekWersjiKompozycjiDesign),
		KompozycjaID: kompozycja.ID,
		Nazwa:        nazwa,
		Uzasadnienie: uzasadnienie,
	}, warstwy)
	if err != nil {
		return shared.DesignBoardVersion{}, bladDesignu(err)
	}
	return wersjaKompozycjiKontraktuDesignu(zapisana, kompozycja.Kod, nil), nil
}

// wersjaKompozycjiKontraktuDesignu składa `DesignBoardVersion` kontraktu.
// Warstwy wchodzą wyłącznie tam, gdzie wołający ich chce — wykaz ich nie niesie.
func wersjaKompozycjiKontraktuDesignu(w dane.WersjaKompozycjiDesignu, kodKompozycji string,
	warstwy []shared.DesignBoardLayer) shared.DesignBoardVersion {

	return shared.DesignBoardVersion{
		Id:         w.Kod,
		BoardId:    kodKompozycji,
		Name:       w.Nazwa,
		Note:       w.Uzasadnienie,
		Layers:     warstwy,
		LayerCount: w.LiczbaWarstw,
		CreatedAt:  chwilaBazy(w.Utworzono),
	}
}

// bladNieznanejKompozycjiDesignu nazywa kompozycję, której rdzeń nie zna —
// zapis wskazujący ją byłby odmową.
func bladNieznanejKompozycjiDesignu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu("kompozycji " + kod + " nie ma w module Design")
	}
	return bladDesignu(err)
}

// wskaznikTekstuDesignu oddaje wskaźnik na tekst — pola opcjonalne kontraktu są
// wskaźnikami, a adresu literału wziąć nie można.
func wskaznikTekstuDesignu(wartosc string) *string {
	return &wartosc
}
