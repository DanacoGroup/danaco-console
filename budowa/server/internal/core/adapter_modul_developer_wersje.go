// Odpowiedzialność pliku: dwie komendy historii pliku edytora —
// `developer.file.version.list` (wykaz migawek) i `developer.file.version.restore`
// (powrót do migawki).
//
// Wersja pliku nie jest commitem i go nie zastępuje. Historia repozytorium
// należy do Gita i jedzie osobną drogą; wersja jest zapisem roboczym edytora —
// powstaje przy zapisie z `createVersion`, często wielokrotnie w obrębie jednej
// zmiany, i pozwala Operatorowi wrócić do stanu, który sam nadpisał, zanim
// cokolwiek zatwierdził.
//
// Przywrócenie zakłada własną migawkę stanu bieżącego, zanim nadpisze plik.
// Bez tego kroku powrót do wersji sprzed godziny kasowałby bezpowrotnie pracę
// tej godziny — a Operator sięga po historię właśnie dlatego, że nie jest pewien,
// który stan jest lepszy.
package core

import (
	"context"
	"errors"
	"os"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// najwiecejWersjiWykazu jest domyślną głębokością historii pliku.
const najwiecejWersjiWykazu = 50

// WykazWersjiPliku obsługuje `developer.file.version.list`.
func (a *adapterDevelopera) WykazWersjiPliku(ctx context.Context,
	z shared.DeveloperFileVersionListRequest) (shared.DeveloperFileVersionListResponse, error) {

	okno, sciezka, err := a.plikOkna(z.WindowId, z.Path)
	if err != nil {
		return shared.DeveloperFileVersionListResponse{}, err
	}
	if a.repozytorium == nil {
		return shared.DeveloperFileVersionListResponse{}, bladZasobuDevelopera(
			"rdzeń nie ma miejsca na wersje plików, więc historia pliku nie istnieje")
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

// PrzywrocWersjePliku obsługuje `developer.file.version.restore`.
//
// Zapis na dysk jest domyślny, lecz nie bezwarunkowy: `writeToDisk: false` daje
// samą treść migawki, którą Code Editor wstawia do bufora bez ruszania pliku.
// To jest droga „zobacz, jak było”, odróżniona od „wróć do tego stanu”.
func (a *adapterDevelopera) PrzywrocWersjePliku(ctx context.Context,
	z shared.DeveloperFileVersionRestoreRequest) (shared.DeveloperFileVersionRestoreResponse, error) {

	okno, err := a.oknoDevelopera(z.WindowId)
	if err != nil {
		return shared.DeveloperFileVersionRestoreResponse{}, err
	}
	if a.repozytorium == nil {
		return shared.DeveloperFileVersionRestoreResponse{}, bladZasobuDevelopera(
			"rdzeń nie ma miejsca na wersje plików, więc nie ma czego przywrócić")
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
	// Wersja niesie własne okno. Przywrócenie migawki założonej w cudzym oknie
	// wyprowadziłoby zapis poza obszar tego okna — a obszar jest granicą, którą
	// moduł sprawdza przy każdej ścieżce, nie tylko przy tych z żądania.
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

// odlozStanPrzedPrzywroceniem zakłada migawkę treści, którą przywrócenie
// nadpisze. Plik, którego na dysku nie ma, nie ma czego odkładać i nie jest to
// przeszkodą — przywrócenie odtworzy go z migawki.
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
