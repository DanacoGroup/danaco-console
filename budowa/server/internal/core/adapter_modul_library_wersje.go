// Odpowiedzialność pliku: wersje pliku repozytorium wiedzy po stronie rdzenia —
// library.version.add, library.version.list i library.version.restore,
// zasilające Versioning Panel. Sam plik i wyszukiwanie mieszkają
// w adapter_modul_library.go.
package core

import (
	"context"
	"fmt"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// DolozWersje obsługuje library.version.add. Dokłada wiersz historii i tym
// samym ruchem przestawia na niego plik macierzysty w jednej transakcji, bo
// kontrakt obiecuje w odpowiedzi oba byty naraz.
func (a *adapterBiblioteki) DolozWersje(ctx context.Context,
	z shared.LibraryVersionAddRequest) (shared.LibraryVersionAddResponse, error) {

	plik, err := a.plik(ctx, z.FileId)
	if err != nil {
		return shared.LibraryVersionAddResponse{}, err
	}
	nowa, err := a.trescNowejWersji(plik, z)
	if err != nil {
		return shared.LibraryVersionAddResponse{}, err
	}
	nowa.Kod = nowyIdentyfikator(przedrostekWersjiBiblioteki)
	nowa.PlikID = plik.ID
	nowa.Autor = z.Author
	nowa.Etykieta = a.etykietaWersji(ctx, plik, z.Label)

	wersja, po, err := a.repozytorium.DolozWersje(ctx, plik.ID, nowa)
	if err != nil {
		return shared.LibraryVersionAddResponse{}, bladBiblioteki(err)
	}
	// Nowa wersja jest od tej chwili treścią bieżącą pliku, więc indeks musi opisywać ją, nie poprzednią.
	a.zaindeksujTresc(ctx, po)
	kontrakt, err := a.zloz(ctx, po)
	if err != nil {
		return shared.LibraryVersionAddResponse{}, err
	}
	return shared.LibraryVersionAddResponse{
		Version: wersjaKontraktu(po.Kod, wersja), File: kontrakt,
	}, nil
}

// trescNowejWersji rozstrzyga, co wersja niesie jako swoją treść. Żądanie bez
// treści i bez ścieżki jest znacznikiem na treści bieżącej — kamieniem
// milowym, biorącym miarę i odwołanie z pliku macierzystego.
func (a *adapterBiblioteki) trescNowejWersji(plik dane.PlikBiblioteki,
	z shared.LibraryVersionAddRequest) (dane.WersjaPlikuBiblioteki, error) {

	if bezWartosci(z.ContentBase64) && bezWartosci(z.SourcePath) {
		return dane.WersjaPlikuBiblioteki{
			RozmiarBajtow: plik.RozmiarBajtow, SumaKontrolna: plik.SumaKontrolna,
			TrescOdwolanie: plik.TrescOdwolanie,
		}, nil
	}
	odwolanie, rozmiar, suma, err := a.trescWgrania(z.ContentBase64, z.SourcePath)
	if err != nil {
		return dane.WersjaPlikuBiblioteki{}, err
	}
	return dane.WersjaPlikuBiblioteki{
		RozmiarBajtow: rozmiar, SumaKontrolna: suma, TrescOdwolanie: odwolanie,
	}, nil
}

// bezWartosci mówi, że pole opcjonalne kontraktu nic nie wnosi — nie ma go
// wcale albo jest puste. Jedno i drugie znaczy dla Operatora to samo.
func bezWartosci(pole *string) bool {
	return pole == nil || *pole == ""
}

// etykietaWersji oddaje etykietę żądania, a przy jej braku nazywa wersję jej
// kolejnością w historii pliku. Numer jest policzony z wierszy w bazie, bo
// schemat biblioteki świadomie nie ma kolumny numeru — wersja jest bytem,
// nie licznikiem.
func (a *adapterBiblioteki) etykietaWersji(ctx context.Context,
	plik dane.PlikBiblioteki, zadana *string) *string {

	if !bezWartosci(zadana) {
		return zadana
	}
	wersje, err := a.repozytorium.Wersje(ctx, plik.ID)
	if err != nil {
		return nil
	}
	etykieta := fmt.Sprintf("wersja %d", len(wersje)+1)
	return &etykieta
}

// Wersje obsługuje `library.version.list`. Zwraca historię pliku od
// najnowszej — tak oddaje ją repozytorium (`dane/library_wersje.go`).
func (a *adapterBiblioteki) Wersje(ctx context.Context,
	z shared.LibraryVersionListRequest) (shared.LibraryVersionListResponse, error) {

	plik, err := a.plik(ctx, z.FileId)
	if err != nil {
		return shared.LibraryVersionListResponse{}, err
	}
	wiersze, err := a.repozytorium.Wersje(ctx, plik.ID)
	if err != nil {
		return shared.LibraryVersionListResponse{}, bladBiblioteki(err)
	}
	wiersze = ograniczWersje(wiersze, z.Limit)

	wersje := make([]shared.LibraryVersion, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wersje = append(wersje, wersjaKontraktu(plik.Kod, wiersz))
	}
	return shared.LibraryVersionListResponse{Versions: wersje}, nil
}

// ograniczWersje obcina wykaz do żądanej granicy. Repozytorium oddaje całą
// historię — obcięcie jest sprawą kontraktu (`limit` jest opcjonalny), nie
// warstwy danych, żeby ta sama metoda `Wersje` służyła też `PrzywrocWersje`,
// której limit nie dotyczy.
func ograniczWersje(wiersze []dane.WersjaPlikuBiblioteki, limit *int) []dane.WersjaPlikuBiblioteki {
	if limit == nil || *limit <= 0 || *limit >= len(wiersze) {
		return wiersze
	}
	return wiersze[:*limit]
}

// wersjaKontraktu składa wersję kontraktu z wiersza repozytorium historii pliku wraz z kodem pliku macierzystego.
func wersjaKontraktu(kodPliku string, wiersz dane.WersjaPlikuBiblioteki) shared.LibraryVersion {
	return shared.LibraryVersion{
		Id: wiersz.Kod, FileId: kodPliku, Label: wiersz.Etykieta, Author: wiersz.Autor,
		SizeBytes: wiersz.RozmiarBajtow, Checksum: wiersz.SumaKontrolna,
		CreatedAt: chwilaBazy(wiersz.Utworzono),
	}
}

// PrzywrocWersje obsługuje library.version.restore. Kod wersji jest szukany
// wyłącznie w obrębie wskazanego pliku — repozytorium samo zwraca błąd braku
// wiersza, gdy Operator poda wersję z innego pliku.
func (a *adapterBiblioteki) PrzywrocWersje(ctx context.Context,
	z shared.LibraryVersionRestoreRequest) (shared.LibraryVersionRestoreResponse, error) {

	plik, err := a.plik(ctx, z.FileId)
	if err != nil {
		return shared.LibraryVersionRestoreResponse{}, err
	}
	if z.VersionId == "" {
		return shared.LibraryVersionRestoreResponse{}, bladWskazaniaBiblioteki(
			"komenda przywrócenia wersji bez wskazania wersji")
	}
	po, err := a.repozytorium.PrzywrocWersje(ctx, plik.ID, z.VersionId)
	if err != nil {
		return shared.LibraryVersionRestoreResponse{}, bladNieznanegoPlikuBiblioteki(z.VersionId, err)
	}
	// Przywrócenie zmienia treść bieżącą tak samo jak dołożenie wersji, indeks idzie tą samą drogą.
	a.zaindeksujTresc(ctx, po)
	kontrakt, err := a.zloz(ctx, po)
	if err != nil {
		return shared.LibraryVersionRestoreResponse{}, err
	}
	return shared.LibraryVersionRestoreResponse{File: kontrakt}, nil
}
