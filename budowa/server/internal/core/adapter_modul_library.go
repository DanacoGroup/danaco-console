// Moduł Library obsługuje wgranie, wykaz, podgląd i wyszukiwanie pliku repozytorium
// wiedzy; metody wersji, kolekcji i etykiet leżą w plikach sąsiednich.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

const (
	przedrostekPlikuBiblioteki    = "plik-"
	przedrostekWersjiBiblioteki   = "wers-"
	przedrostekKolekcjiBiblioteki = "kolek-"
	// Byty migracji 180-186 wychodzą kontraktem pod własnym identyfikatorem, nie pod kluczem wiersza.
	przedrostekRegulyBiblioteki        = "regula-"
	przedrostekWpisuAudytuBiblioteki   = "audyt-"
	przedrostekPolitykiBiblioteki      = "retencja-"
	przedrostekUtrwaleniaBiblioteki    = "utrw-"
	przedrostekUdostepnieniaBiblioteki = "udost-"
	przedrostekWebhookaBiblioteki      = "hook-"
	przedrostekSugestiiBiblioteki      = "sug-"
)

type adapterBiblioteki struct {
	repozytorium dane.RepozytoriumBiblioteki
	magazyn      *magazynTresciBiblioteki
	kanaly       *models.Rejestr
	uruchamiacz  session.Uruchamiacz
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
}

func (a *adapterBiblioteki) ZNarzedziami(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterBiblioteki {

	a.uruchamiacz, a.rozstrzygacz, a.katalog = uruchamiacz, rozstrzygacz, katalog
	return a
}

func nowyAdapterBiblioteki(repozytorium dane.RepozytoriumBiblioteki) *adapterBiblioteki {
	return &adapterBiblioteki{
		repozytorium: repozytorium,
		magazyn:      nowyMagazynTresciBiblioteki(konfiguracja.KatalogDanychDomyslny()),
	}
}

func (a *adapterBiblioteki) ZKatalogiemDanych(katalog string) *adapterBiblioteki {
	if katalog != "" {
		a.magazyn = nowyMagazynTresciBiblioteki(katalog)
	}
	posprzatajMagazynTresci(context.Background(), a.repozytorium, a.magazyn)
	return a
}

// Wiersz pliku powstaje po utrwaleniu bajtów, żeby wykaz nie pokazał pliku bez treści.
func (a *adapterBiblioteki) Wgraj(ctx context.Context,
	z shared.LibraryFileUploadRequest) (shared.LibraryFileUploadResponse, error) {

	if z.Name == "" {
		return shared.LibraryFileUploadResponse{}, bladWskazaniaBiblioteki("żądanie bez nazwy pliku")
	}

	trescOdwolanie, rozmiar, sumaKontrolna, err := a.trescWgrania(z.ContentBase64, z.SourcePath)
	if err != nil {
		return shared.LibraryFileUploadResponse{}, err
	}

	kod := nowyIdentyfikator(przedrostekPlikuBiblioteki)
	plik, err := a.repozytorium.ZapiszPlik(ctx, dane.PlikBiblioteki{
		Kod: kod, Nazwa: z.Name, MimeType: z.MimeType, RozmiarBajtow: rozmiar,
		Sciezka: z.SourcePath, ProjektID: z.ProjectId, ModulZrodlowyID: z.SourceModuleId,
		SumaKontrolna: sumaKontrolna, TrescOdwolanie: trescOdwolanie,
	})
	if err != nil {
		return shared.LibraryFileUploadResponse{}, bladBiblioteki(err)
	}

	wersja, err := a.repozytorium.ZapiszWersje(ctx, plik.ID, dane.WersjaPlikuBiblioteki{
		Kod: nowyIdentyfikator(przedrostekWersjiBiblioteki), PlikID: plik.ID,
		RozmiarBajtow: rozmiar, SumaKontrolna: sumaKontrolna, TrescOdwolanie: trescOdwolanie,
	})
	if err != nil {
		return shared.LibraryFileUploadResponse{}, bladBiblioteki(err)
	}
	wersjaID := wersja.ID
	plik, err = a.repozytorium.ZapiszPlik(ctx, plikZBiezacaWersja(plik, wersjaID))
	if err != nil {
		return shared.LibraryFileUploadResponse{}, bladBiblioteki(err)
	}

	if len(z.Tags) > 0 {
		if _, err := a.repozytorium.UstawEtykiety(ctx, plik.Kod, z.Tags); err != nil {
			return shared.LibraryFileUploadResponse{}, bladBiblioteki(err)
		}
	}
	a.zaindeksujTresc(ctx, plik)

	kontrakt, err := a.zloz(ctx, plik)
	if err != nil {
		return shared.LibraryFileUploadResponse{}, err
	}
	return shared.LibraryFileUploadResponse{File: kontrakt}, nil
}

func plikZBiezacaWersja(plik dane.PlikBiblioteki, wersjaID int64) dane.PlikBiblioteki {
	plik.WersjaBiezacaID = &wersjaID
	return plik
}

func (a *adapterBiblioteki) Wykaz(ctx context.Context,
	z shared.LibraryFileListRequest) (shared.LibraryFileListResponse, error) {

	// Wykaz domyślny pokazuje wyłącznie zasoby czynne; archiwum ma własną komendę.
	stan := dane.StanZasobuCzynny
	filtr := filtrZadania(z.Query, z.Tags, z.CollectionId, z.ProjectId, z.Limit, z.Offset)
	filtr.Stan = &stan
	wiersze, lacznie, err := a.repozytorium.Pliki(ctx, filtr)
	if err != nil {
		return shared.LibraryFileListResponse{}, bladBiblioteki(err)
	}
	pliki, err := a.zlozWiele(ctx, wiersze)
	if err != nil {
		return shared.LibraryFileListResponse{}, err
	}
	return shared.LibraryFileListResponse{Files: pliki, Total: &lacznie}, nil
}

func (a *adapterBiblioteki) Szukaj(ctx context.Context,
	z shared.LibraryFileSearchRequest) (shared.LibraryFileSearchResponse, error) {

	fraza := strings.TrimSpace(z.Query)
	if fraza == "" {
		return shared.LibraryFileSearchResponse{}, bladWskazaniaBiblioteki("żądanie bez frazy wyszukiwania")
	}
	filtr := filtrZadania(nil, nil, nil, nil, z.Limit, nil)
	wiersze, lacznie, err := a.repozytorium.Szukaj(ctx, fraza, filtr)
	if err != nil {
		return shared.LibraryFileSearchResponse{}, bladBiblioteki(err)
	}
	pliki, err := a.zlozWiele(ctx, wiersze)
	if err != nil {
		return shared.LibraryFileSearchResponse{}, err
	}
	return shared.LibraryFileSearchResponse{Files: pliki, Total: lacznie}, nil
}

func (a *adapterBiblioteki) Podglad(ctx context.Context,
	z shared.LibraryFilePreviewRequest) (shared.LibraryFilePreviewResponse, error) {

	plik, err := a.plik(ctx, z.FileId)
	if err != nil {
		return shared.LibraryFilePreviewResponse{}, err
	}
	if plik.TrescOdwolanie == nil || *plik.TrescOdwolanie == "" {
		return shared.LibraryFilePreviewResponse{}, bladBrakuTresciBiblioteki(plik.Kod)
	}

	rodzaj := rodzajPodgladu(plik.MimeType, *plik.TrescOdwolanie)
	podglad := shared.LibraryPreview{FileId: plik.Kod, Kind: rodzaj}
	if rodzaj == shared.LibraryPreviewKindText {
		tekst, skrocono, err := trescPodgladuBiblioteki(*plik.TrescOdwolanie, z.MaxChars)
		if err != nil {
			return shared.LibraryFilePreviewResponse{}, err
		}
		podglad.Text = &tekst
		if skrocono {
			podglad.Truncated = &skrocono
		}
		if z.Page != nil {
			podglad.Page = z.Page
		}
		return shared.LibraryFilePreviewResponse{Preview: podglad}, nil
	}
	// Odwołanie obrazu wychodzi wyłącznie w postaci względnej.
	podglad.ImageRef = odwolanieObrazuPodgladu(rodzaj, *plik.TrescOdwolanie)
	if z.Page != nil {
		podglad.Page = z.Page
	}
	return shared.LibraryFilePreviewResponse{Preview: podglad}, nil
}

func filtrZadania(fraza *string, etykiety []string, kolekcjaID, projektID *string,
	limit, offset *int) dane.FiltrPlikow {

	filtr := dane.FiltrPlikow{Fraza: fraza, Etykiety: etykiety, KolekcjaKod: kolekcjaID, ProjektID: projektID}
	if limit != nil {
		filtr.Limit = *limit
	}
	if offset != nil {
		filtr.Offset = *offset
	}
	return filtr
}

func (a *adapterBiblioteki) plik(ctx context.Context, kod string) (dane.PlikBiblioteki, error) {
	if kod == "" {
		return dane.PlikBiblioteki{}, bladWskazaniaBiblioteki("komenda bez wskazania pliku")
	}
	plik, err := a.repozytorium.Plik(ctx, kod)
	if err != nil {
		return dane.PlikBiblioteki{}, bladNieznanegoPlikuBiblioteki(kod, err)
	}
	return plik, nil
}

func (a *adapterBiblioteki) zlozWiele(ctx context.Context, wiersze []dane.PlikBiblioteki) ([]shared.LibraryFile, error) {
	pliki := make([]shared.LibraryFile, 0, len(wiersze))
	for _, wiersz := range wiersze {
		plik, err := a.zloz(ctx, wiersz)
		if err != nil {
			return nil, err
		}
		pliki = append(pliki, plik)
	}
	return pliki, nil
}

// Pole `path` kontraktu niesie ścieżkę wewnątrz repozytorium, nie ścieżkę z maszyny Operatora.
func (a *adapterBiblioteki) zloz(ctx context.Context, wiersz dane.PlikBiblioteki) (shared.LibraryFile, error) {
	etykiety, err := a.repozytorium.Etykiety(ctx, wiersz.ID)
	if err != nil {
		return shared.LibraryFile{}, bladBiblioteki(err)
	}
	kolekcje, err := a.repozytorium.KolekcjePliku(ctx, wiersz.ID)
	if err != nil {
		return shared.LibraryFile{}, bladBiblioteki(err)
	}
	return shared.LibraryFile{
		Id: wiersz.Kod, Name: wiersz.Nazwa, Path: wiersz.SciezkaRepozytorium,
		MimeType:  wiersz.MimeType,
		SizeBytes: wiersz.RozmiarBajtow, SourceModuleId: wiersz.ModulZrodlowyID, ProjectId: wiersz.ProjektID,
		Tags: etykiety, CollectionIds: kolekcje, VersionId: a.idWersjiBiezacej(ctx, wiersz),
		Checksum:  wiersz.SumaKontrolna,
		CreatedAt: chwilaBazy(wiersz.Utworzono), UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}, nil
}

func (a *adapterBiblioteki) idWersjiBiezacej(ctx context.Context, wiersz dane.PlikBiblioteki) *string {
	if wiersz.WersjaBiezacaID == nil {
		return nil
	}
	wersje, err := a.repozytorium.Wersje(ctx, wiersz.ID)
	if err != nil {
		return nil
	}
	for _, wersja := range wersje {
		if wersja.ID == *wiersz.WersjaBiezacaID {
			kod := wersja.Kod
			return &kod
		}
	}
	return nil
}

// Kolizja wiersza wychodzi kodem conflict (wzór adapter_modul_auth.go):
// internal_error klient ponawia.
func bladBiblioteki(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, dane.ErrKolizjaWiersza) {
		return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeConflict, err))
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

func bladWskazaniaBiblioteki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "moduł Library: "+powod))
}

func bladNieznanegoPlikuBiblioteki(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Library: plik nie istnieje: "+kod))
	}
	return bladBiblioteki(err)
}

func bladNieznanejEtykietyBiblioteki(nazwa string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Library: etykieta nie istnieje: "+nazwa))
	}
	return bladBiblioteki(err)
}

func bladNieznanejReguly(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Library: reguła nie istnieje: "+kod))
	}
	return bladBiblioteki(err)
}

func bladNieznanejSugestii(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Library: sugestia nie istnieje: "+kod))
	}
	return bladBiblioteki(err)
}
