// Odpowiedzialność pliku: wersje pliku repozytorium wiedzy po stronie rdzenia —
// `library.version.add`, `library.version.list` i `library.version.restore`,
// zasilające Versioning Panel. Sam plik i wyszukiwanie mieszkają
// w `adapter_modul_library.go`.
//
// `version.add` jest jedyną drogą, którą historia pliku rośnie: wgranie pliku
// (`library.file.upload`) zawsze nadaje nowy kod pliku (`nowyIdentyfikator`,
// `Wgraj`), więc dwa wgrania dają dwa osobne pliki po jednej wersji każdy,
// a nie dwie wersje jednego pliku.
package core

import (
	"context"
	"fmt"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// DolozWersje obsługuje `library.version.add`. Dokłada wiersz historii i tym
// samym ruchem przestawia na niego plik macierzysty (jedna transakcja,
// `dane/library_wersje_zapis.go`), bo kontrakt obiecuje w odpowiedzi oba byty
// naraz — wersję i plik po dołożeniu.
//
// Nieznany plik jest odmową, nie cichym powodzeniem: dołożenie wersji do pliku,
// którego nie ma, nie może skończyć się utworzeniem czegokolwiek — `a.plik`
// nazywa brak kodem `not_found`, a puste wskazanie kodem walidacji.
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
	// Nowa wersja jest od tej chwili treścią bieżącą pliku, więc indeks treści
	// musi opisywać ją, a nie treść poprzednią — inaczej wyszukiwanie po treści
	// oddawałoby trafienia w słowa, których w pliku już nie ma
	// (`adapter_modul_library_indeks.go`).
	a.zaindeksujTresc(ctx, po)
	kontrakt, err := a.zloz(ctx, po)
	if err != nil {
		return shared.LibraryVersionAddResponse{}, err
	}
	return shared.LibraryVersionAddResponse{
		Version: wersjaKontraktu(po.Kod, wersja), File: kontrakt,
	}, nil
}

// trescNowejWersji rozstrzyga, co wersja niesie jako swoją treść.
//
// Żądanie z treścią albo ze ścieżką jedzie tą samą drogą co wgranie pliku
// (`trescWgrania`, `adapter_modul_library_tresc.go`): bajty lądują w magazynie
// treści rdzenia, a ścieżka bloba staje się odwołaniem — niezależnie od tego,
// czy przyszły base64, czy zostały wciągnięte spod `sourcePath`. Rozmiar i suma
// kontrolna są liczone w obu drogach, bo wersja bez sumy nie daje się sprawdzić.
// Dzięki temu historia jest odtwarzalna: każda wersja z treścią ma własne,
// zamrożone bajty na nośniku, których nie ruszy zmiana pliku źródłowego.
//
// Żądanie bez treści i bez ścieżki jest znacznikiem na treści bieżącej —
// kamieniem milowym z kontraktu (`LibraryVersionAddRequest.label`). Taka wersja
// bierze miarę i odwołanie z pliku macierzystego, zamiast wstawić pustkę:
// wersja pusta przy przywróceniu wymazałaby plikowi odwołanie do treści, więc
// „oznaczenie kamienia milowego" niszczyłoby zasób, który miało utrwalić.
//
// Nieudany zapis treści jest odmową dołożenia wersji: wiersz historii powstaje
// dopiero po utrwaleniu bajtów, inaczej Versioning Panel pokazywałby wersję, do
// której nie ma czego przywrócić.
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
// kolejnością w historii pliku („wersja 3").
//
// To nie jest dana zmyślona: numer jest policzony z wierszy leżących w bazie.
// Schemat biblioteki (`store/migracja_045_biblioteka.sql`) świadomie nie ma
// kolumny numeru — wersja jest bytem, nie licznikiem — więc kolejność jest
// wyłącznie widokiem chwili zapisu. Bez tej etykiety Versioning Panel
// pokazywałby surowy kod `wers-…`, po którym nie widać, która wersja jest
// która. Nieudany odczyt historii zostawia etykietę pustą — brak nazwy jest
// lepszy niż nazwa policzona z niczego.
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

// wersjaKontraktu składa wersję kontraktu z wiersza repozytorium.
func wersjaKontraktu(kodPliku string, wiersz dane.WersjaPlikuBiblioteki) shared.LibraryVersion {
	return shared.LibraryVersion{
		Id: wiersz.Kod, FileId: kodPliku, Label: wiersz.Etykieta, Author: wiersz.Autor,
		SizeBytes: wiersz.RozmiarBajtow, Checksum: wiersz.SumaKontrolna,
		CreatedAt: chwilaBazy(wiersz.Utworzono),
	}
}

// PrzywrocWersje obsługuje `library.version.restore`. Kod wersji jest szukany
// wyłącznie w obrębie wskazanego pliku (`dane/library_wersje.go`,
// `wersjaDoPrzywrocenia`) — repozytorium sam zwraca błąd braku wiersza, gdy
// Operator poda wersję z innego pliku, więc adapter nie musi tego sprawdzać
// osobno.
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
	// Przywrócenie zmienia treść bieżącą tak samo jak dołożenie wersji, więc
	// indeks idzie za nią tą samą drogą.
	a.zaindeksujTresc(ctx, po)
	kontrakt, err := a.zloz(ctx, po)
	if err != nil {
		return shared.LibraryVersionRestoreResponse{}, err
	}
	return shared.LibraryVersionRestoreResponse{File: kontrakt}, nil
}
