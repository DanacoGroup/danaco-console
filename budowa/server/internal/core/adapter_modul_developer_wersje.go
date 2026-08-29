// Odpowiedzialność pliku: dwie komendy historii pliku edytora — wykaz migawek roboczych
// `developer.file.version.list` i powrót do migawki `developer.file.version.restore`.
package core

import (
	"context"
	"errors"
	"os"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// najwiecejWersjiWykazu jest domyślną głębokością historii pliku edytora Code Editor modułu Developer platformy.
const najwiecejWersjiWykazu = 50

// WykazWersjiPliku obsługuje `developer.file.version.list` i zwraca migawki wskazanego pliku okna Developer.
func (a *adapterDevelopera) WykazWersjiPliku(ctx context.Context,
	z shared.DeveloperFileVersionListRequest) (shared.DeveloperFileVersionListResponse, error) {

	okno, sciezka, err := a.plikOkna(z.WindowId, z.Path)
	if err != nil {
		return shared.DeveloperFileVersionListResponse{}, err
	}
	if a.repozytorium == nil {
		return shared.DeveloperFileVersionListResponse{}, bladZasobuDevelopera(
			"serwer nie ma miejsca na wersje plików, więc historia pliku nie istnieje")
	}

	limit := najwiecejWersjiWykazu
	if z.Limit != nil && *z.Limit > 0 {
		limit = *z.Limit
	}
	wiersze, err := a.repozytorium.Wersje(ctx, okno.Id, sciezka, limit)
	if err != nil {
		return shared.DeveloperFileVersionListResponse{}, bladWykonaniaDevelopera(
			"nie można odczytać historii pliku " + z.Path + ": " + err.Error())
	}

	wersje := make([]shared.DeveloperFileVersion, 0, len(wiersze))
	for _, wiersz := range wiersze {
		rozmiar := wiersz.Rozmiar
		wersje = append(wersje, shared.DeveloperFileVersion{
			Id:        wiersz.Kod,
			Path:      wiersz.Sciezka,
			SizeBytes: &rozmiar,
			CreatedAt: chwilaBazy(wiersz.Utworzono),
		})
	}
	razem := len(wersje)
	return shared.DeveloperFileVersionListResponse{Versions: wersje, Total: &razem}, nil
}

// PrzywrocWersjePliku obsługuje `developer.file.version.restore`; zapis na dysk jest domyślny,
// lecz nie bezwarunkowy — pole zapisu na dysk może dać samą treść migawki bez ruszania pliku.
func (a *adapterDevelopera) PrzywrocWersjePliku(ctx context.Context,
	z shared.DeveloperFileVersionRestoreRequest) (shared.DeveloperFileVersionRestoreResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperFileVersionRestoreResponse{}, err
	}
	if a.repozytorium == nil {
		return shared.DeveloperFileVersionRestoreResponse{}, bladZasobuDevelopera(
			"serwer nie ma miejsca na wersje plików, więc nie ma czego przywrócić")
	}
	kod := strings.TrimSpace(z.VersionId)
	if kod == "" {
		return shared.DeveloperFileVersionRestoreResponse{}, bladZadaniaDevelopera(
			"przywrócenie wymaga wskazania wersji")
	}

	wersja, err := a.repozytorium.WersjaPoKodzie(ctx, kod)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.DeveloperFileVersionRestoreResponse{}, bladZasobuDevelopera(
			"nie ma wersji pliku o identyfikatorze " + kod)
	}
	if err != nil {
		return shared.DeveloperFileVersionRestoreResponse{}, bladWykonaniaDevelopera(
			"nie można odczytać wersji " + kod + ": " + err.Error())
	}
	// Wersja niesie własne okno; przywrócenie migawki cudzego okna wyprowadziłoby zapis poza obszar.
	if wersja.OknoKod != okno.Id {
		return shared.DeveloperFileVersionRestoreResponse{}, bladDostepuDevelopera(
			"wersja " + kod + " należy do innego okna niż " + okno.Id)
	}
	sciezka, err := sciezkaWObszarze(a.korzenieOkna(okno), wersja.Sciezka, plikIstnieje)
	if err != nil {
		return shared.DeveloperFileVersionRestoreResponse{}, err
	}

	zapisz := z.WriteToDisk == nil || *z.WriteToDisk
	if zapisz {
		if err := sprawdzZmianeSystemu(okno.TrybUprawnien, "przywrócenie wersji pliku"); err != nil {
			return shared.DeveloperFileVersionRestoreResponse{}, err
		}
		if err := a.odlozStanPrzedPrzywroceniem(ctx, okno.Id, sciezka); err != nil {
			return shared.DeveloperFileVersionRestoreResponse{}, err
		}
		if err := os.WriteFile(sciezka, []byte(wersja.Tresc), 0o644); err != nil {
			return shared.DeveloperFileVersionRestoreResponse{}, bladWykonaniaDevelopera(
				"nie można przywrócić pliku " + wersja.Sciezka + ": " + err.Error())
		}
	}

	plik := shared.DeveloperFile{
		Path:      sciezka,
		Content:   wskaznikTekstu(wersja.Tresc),
		Language:  jezykPliku(sciezka),
		VersionId: wskaznikTekstu(wersja.Kod),
		UpdatedAt: chwilaBazy(wersja.Utworzono),
	}
	rozmiar := wersja.Rozmiar
	plik.SizeBytes = &rozmiar
	if opis, err := os.Stat(sciezka); err == nil {
		plik.UpdatedAt = opis.ModTime().UTC().UnixMilli()
	}
	return shared.DeveloperFileVersionRestoreResponse{File: plik}, nil
}

// odlozStanPrzedPrzywroceniem zakłada migawkę treści, którą przywrócenie nadpisze; plik
// nieobecny na dysku nie ma czego odkładać, a to nie jest przeszkodą.
func (a *adapterDevelopera) odlozStanPrzedPrzywroceniem(ctx context.Context,
	oknoKod, sciezka string) error {

	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		return nil
	}
	if err := a.repozytorium.ZapiszWersje(ctx, dane.WersjaPliku{
		Kod:     nowyIdentyfikator(przedrostekWersjiPliku),
		OknoKod: oknoKod,
		Sciezka: sciezka,
		Tresc:   string(bajty),
		Rozmiar: int64(len(bajty)),
	}); err != nil {
		return bladWykonaniaDevelopera(
			"nie można odłożyć stanu sprzed przywrócenia pliku " + sciezka + ": " + err.Error())
	}
	return nil
}
