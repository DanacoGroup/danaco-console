// Moduł Library — cykl życia zasobu: `library.file.move`,
// `library.file.archive`, `library.file.restore`, `library.file.delete`.
//
// Trzy pierwsze czynności są odwracalne i tak też działają: przeniesienie zmienia
// miejsce w strukturze, archiwizacja zdejmuje zasób z wykazu domyślnego, a
// przywrócenie oddaje go z powrotem — w każdym przypadku wiersz, wersje,
// etykiety i kolekcje zostają nietknięte.
//
// Czwarta jest jedyną drogą utraty danych w tym module i wygląda inaczej:
// wymaga potwierdzenia wprost (`confirm`), zdejmuje wiersz zasobu wraz z jego
// wersjami i zostawia po sobie wpis w dzienniku audytu, który przeżywa usunięty
// zasób (dziennik wskazuje zasób kodem, nie kluczem obcym).
//
// Bajty treści nie znikają razem z wierszem, i to jest zamierzone: ta sama treść
// bywa współdzielona przez inny zasób pod tą samą sumą kontrolną (magazyn jest
// adresowany treścią). Bloby osierocone zdejmuje obchód magazynu przy starcie
// rdzenia (`adapter_modul_library_sprzatanie.go`).
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// PrzeniesZasoby obsługuje `library.file.move`.
func (a *adapterBiblioteki) PrzeniesZasoby(ctx context.Context,
	z shared.LibraryFileMoveRequest) (shared.LibraryFileMoveResponse, error) {

	if len(z.FileIds) == 0 {
		return shared.LibraryFileMoveResponse{}, bladWskazaniaBiblioteki(
			"przeniesienie bez wskazania zasobów")
	}
	sciezka := sciezkaRepozytoriumBiblioteki(z.Path)
	if sciezka == "" {
		return shared.LibraryFileMoveResponse{}, bladWskazaniaBiblioteki(
			"przeniesienie bez ścieżki docelowej — zasób musi gdzieś stanąć")
	}
	// Każdy wskazany zasób jest sprawdzany wprost: wykaz kodów z literówką
	// przeszedłby po cichu jako „przeniesiono mniej", a Operator zobaczyłby
	// powodzenie czynności, której nie było.
	for _, kod := range z.FileIds {
		if _, err := a.plik(ctx, kod); err != nil {
			return shared.LibraryFileMoveResponse{}, err
		}
	}
	wiersze, err := a.repozytorium.UstawSciezkeRepozytorium(ctx, z.FileIds, sciezka)
	if err != nil {
		return shared.LibraryFileMoveResponse{}, bladBiblioteki(err)
	}
	pliki, err := a.zlozWiele(ctx, wiersze)
	if err != nil {
		return shared.LibraryFileMoveResponse{}, err
	}
	for _, plik := range pliki {
		a.odnotuj(ctx, shared.LibraryAuditActionChange, wskazanieBiblioteki(plik.Id),
			"przeniesienie pod ścieżkę "+sciezka)
		a.zglosNasluchom(shared.LibraryWebhookEventFileChanged, plik.Id)
	}
	return shared.LibraryFileMoveResponse{MovedCount: len(pliki), Files: pliki}, nil
}

// ZarchiwizujZasoby obsługuje `library.file.archive`.
func (a *adapterBiblioteki) ZarchiwizujZasoby(ctx context.Context,
	z shared.LibraryFileArchiveRequest) (shared.LibraryFileArchiveResponse, error) {

	pliki, err := a.przestawStan(ctx, z.FileIds, dane.StanZasobuZarchiwizowany,
		"archiwizacja bez wskazania zasobów")
	if err != nil {
		return shared.LibraryFileArchiveResponse{}, err
	}
	powod := strings.TrimSpace(wartoscTekstu(z.Reason))
	opis := "przeniesienie do archiwum"
	if powod != "" {
		opis += ": " + powod
	}
	for _, plik := range pliki {
		a.odnotuj(ctx, shared.LibraryAuditActionArchive, wskazanieBiblioteki(plik.Id), opis)
		a.zglosNasluchom(shared.LibraryWebhookEventFileArchived, plik.Id)
	}
	return shared.LibraryFileArchiveResponse{ArchivedCount: len(pliki), Files: pliki}, nil
}

// PrzywrocZasoby obsługuje `library.file.restore`.
func (a *adapterBiblioteki) PrzywrocZasoby(ctx context.Context,
	z shared.LibraryFileRestoreRequest) (shared.LibraryFileRestoreResponse, error) {

	pliki, err := a.przestawStan(ctx, z.FileIds, dane.StanZasobuCzynny,
		"przywrócenie bez wskazania zasobów")
	if err != nil {
		return shared.LibraryFileRestoreResponse{}, err
	}
	for _, plik := range pliki {
		a.odnotuj(ctx, shared.LibraryAuditActionRestore, wskazanieBiblioteki(plik.Id),
			"przywrócenie z archiwum")
		a.zglosNasluchom(shared.LibraryWebhookEventFileRestored, plik.Id)
	}
	return shared.LibraryFileRestoreResponse{RestoredCount: len(pliki), Files: pliki}, nil
}

// UsunZasoby obsługuje `library.file.delete`.
//
// Wpis audytu powstaje PRZED usunięciem: po zdjęciu wiersza nazwa zasobu
// przestaje istnieć, a dziennik ma powiedzieć, co zniknęło, nie tylko że coś
// zniknęło.
func (a *adapterBiblioteki) UsunZasoby(ctx context.Context,
	z shared.LibraryFileDeleteRequest) (shared.LibraryFileDeleteResponse, error) {

	if len(z.FileIds) == 0 {
		return shared.LibraryFileDeleteResponse{}, bladWskazaniaBiblioteki(
			"usunięcie bez wskazania zasobów")
	}
	if !z.Confirm {
		return shared.LibraryFileDeleteResponse{}, bladWskazaniaBiblioteki(
			"usunięcie trwałe jest nieodwracalne: zdejmuje zasób wraz z wersjami — " +
				"potwierdź żądanie polem confirm")
	}
	for _, kod := range z.FileIds {
		zasob, err := a.plik(ctx, kod)
		if err != nil {
			return shared.LibraryFileDeleteResponse{}, err
		}
		a.odnotuj(ctx, shared.LibraryAuditActionDelete, wskazanieBiblioteki(zasob.Kod),
			"usunięcie trwałe zasobu "+zasob.Nazwa)
	}
	usuniete, wersje, err := a.repozytorium.UsunPliki(ctx, z.FileIds)
	if err != nil {
		return shared.LibraryFileDeleteResponse{}, bladBiblioteki(err)
	}
	return shared.LibraryFileDeleteResponse{DeletedCount: usuniete, DeletedVersions: wersje}, nil
}

// przestawStan jest wspólnym ciałem archiwizacji i przywrócenia — obie czynności
// różni jedna wartość stanu.
func (a *adapterBiblioteki) przestawStan(ctx context.Context, kody []string, stan,
	powodPustki string) ([]shared.LibraryFile, error) {

	if len(kody) == 0 {
		return nil, bladWskazaniaBiblioteki(powodPustki)
	}
	for _, kod := range kody {
		if _, err := a.plik(ctx, kod); err != nil {
			return nil, err
		}
	}
	wiersze, err := a.repozytorium.UstawStanPlikow(ctx, kody, stan)
	if err != nil {
		return nil, bladBiblioteki(err)
	}
	return a.zlozWiele(ctx, wiersze)
}

// sciezkaRepozytoriumBiblioteki sprowadza ścieżkę żądania do jednej postaci:
// człony rozdzielone ukośnikiem, bez ukośnika wiodącego i zamykającego.
//
// Postać jest ważna, bo klient czyta pierwszy człon jako katalog nawigacji
// i zawęża wykaz po przedrostku (`client/src/moduly/library/wykaz-plikow.ts`).
// Dwie zapisane ścieżki różniące się samym ukośnikiem byłyby dla niego dwoma
// różnymi katalogami.
func sciezkaRepozytoriumBiblioteki(sciezka string) string {
	czlony := strings.Split(strings.TrimSpace(sciezka), "/")
	wynik := make([]string, 0, len(czlony))
	for _, czlon := range czlony {
		oczyszczony := strings.TrimSpace(czlon)
		// Człon „.." nie ma znaczenia w strukturze repozytorium — to nie jest
		// system plików. Przepuszczony, byłby zaproszeniem do czytania ścieżki
		// jak katalogu na dysku.
		if oczyszczony == "" || oczyszczony == "." || oczyszczony == ".." {
			continue
		}
		wynik = append(wynik, oczyszczony)
	}
	return strings.Join(wynik, "/")
}
