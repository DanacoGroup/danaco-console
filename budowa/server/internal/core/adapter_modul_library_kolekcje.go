// Obszar kolekcji i etykiet modułu Library: `UstawEtykiety`, `UtworzKolekcje`
// i `PrzypiszDoKolekcji` na typie `*adapterBiblioteki`. Obie połowy komendy
// `library.tag.set` mają semantykę wymiany, opartą o `UstawKolekcjePliku`, a
// nie o dokładanie.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// UtworzKolekcje zakłada nową kolekcję i oddaje jej identyfikator — obsługuje
// `library.collection.create`.
func (a *adapterBiblioteki) UtworzKolekcje(ctx context.Context,
	z shared.LibraryCollectionCreateRequest) (shared.LibraryCollectionCreateResponse, error) {

	if z.Name == "" {
		return shared.LibraryCollectionCreateResponse{}, bladWskazaniaBiblioteki("nazwa kolekcji jest wymagana")
	}
	kolekcja, err := a.repozytorium.UtworzKolekcje(ctx, dane.KolekcjaBiblioteki{
		Kod: nowyIdentyfikator(przedrostekKolekcjiBiblioteki), Nazwa: z.Name, Opis: z.Description,
	})
	if err != nil {
		return shared.LibraryCollectionCreateResponse{}, bladBiblioteki(err)
	}
	return shared.LibraryCollectionCreateResponse{CollectionId: kolekcja.Kod, Name: kolekcja.Nazwa}, nil
}

// PrzypiszDoKolekcji przypisuje pliki do kolekcji i oddaje faktyczną liczbę
// przypisanych — warstwa danych pomija kody plików, których nie ma, więc
// odpowiedź mówi prawdę o skutku, nie powtarza żądania. Obsługuje
// `library.collection.assign`.
func (a *adapterBiblioteki) PrzypiszDoKolekcji(ctx context.Context,
	z shared.LibraryCollectionAssignRequest) (shared.LibraryCollectionAssignResponse, error) {

	if z.CollectionId == "" {
		return shared.LibraryCollectionAssignResponse{}, bladWskazaniaBiblioteki("komenda bez wskazania kolekcji")
	}
	przypisane, err := a.repozytorium.PrzypiszDoKolekcji(ctx, z.CollectionId, z.FileIds)
	if err != nil {
		return shared.LibraryCollectionAssignResponse{}, bladNieznanejKolekcji(z.CollectionId, err)
	}
	return shared.LibraryCollectionAssignResponse{
		CollectionId: z.CollectionId, AssignedCount: len(przypisane),
	}, nil
}

// UstawEtykiety podmienia komplet etykiet pliku i komplet jego kolekcji na
// wskazane. Obsługuje `library.tag.set`.
func (a *adapterBiblioteki) UstawEtykiety(ctx context.Context,
	z shared.LibraryTagSetRequest) (shared.LibraryTagSetResponse, error) {

	if z.FileId == "" {
		return shared.LibraryTagSetResponse{}, bladWskazaniaBiblioteki("komenda bez wskazania pliku")
	}
	// Plik sprawdzany wprost i najpierw, żeby literówka w kodzie nie wracała jako
	// odmowa kolekcji.
	if _, err := a.plik(ctx, z.FileId); err != nil {
		return shared.LibraryTagSetResponse{}, err
	}
	// Kolekcje idą pierwsze, bo mogą odmówić i wykaz pusty jest ustawieniem.
	if _, err := a.repozytorium.UstawKolekcjePliku(ctx, z.FileId, z.CollectionIds); err != nil {
		return shared.LibraryTagSetResponse{}, bladNieznanejKolekcji(wykazKolekcji(z.CollectionIds), err)
	}
	if _, err := a.repozytorium.UstawEtykiety(ctx, z.FileId, z.Tags); err != nil {
		return shared.LibraryTagSetResponse{}, bladNieznanegoPlikuBiblioteki(z.FileId, err)
	}

	wiersz, err := a.repozytorium.Plik(ctx, z.FileId)
	if err != nil {
		return shared.LibraryTagSetResponse{}, bladNieznanegoPlikuBiblioteki(z.FileId, err)
	}
	plik, err := a.zloz(ctx, wiersz)
	if err != nil {
		return shared.LibraryTagSetResponse{}, err
	}
	return shared.LibraryTagSetResponse{File: plik}, nil
}

// wykazKolekcji nazywa w odmowie cały wskazany wykaz kolekcji, bo warstwa
// danych nie oddaje, który kod zawiódł przy pierwszym nieznanym.
func wykazKolekcji(kody []string) string {
	if len(kody) == 0 {
		return "(wykaz pusty)"
	}
	return strings.Join(kody, ", ")
}

// bladNieznanejKolekcji odróżnia "kolekcji nie ma" od "odczyt się nie powiódł"
// — jedyny błąd swoisty temu obszarowi kolekcji i etykiet.
func bladNieznanejKolekcji(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Library: kolekcja nie istnieje: "+kod))
	}
	return bladBiblioteki(err)
}
