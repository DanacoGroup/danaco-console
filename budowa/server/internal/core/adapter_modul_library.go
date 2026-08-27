// Moduł Library obsługuje wgranie, wykaz, podgląd i wyszukiwanie pliku
// repozytorium wiedzy; typ adaptera i konstruktor stoją tutaj, metody wersji,
// kolekcji i etykiet leżą w plikach sąsiednich.
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

// Przedrostki identyfikatorów bytów modułu. Zadeklarowane tu w całości — łącznie
// z przedrostkiem kolekcji, którego ten plik nie używa — żeby pozostałe pliki
// modułu brały je stąd, zamiast deklarować po raz drugi.
const (
	przedrostekPlikuBiblioteki    = "plik-"
	przedrostekWersjiBiblioteki   = "wers-"
	przedrostekKolekcjiBiblioteki = "kolek-"
	// Byty dobudowane migracjami 180-186: opis, reguły, audyt, retencja,
	// utrwalenie, udostępnienia, nasłuchy i sugestie. Każdy wychodzi kontraktem
	// pod własnym identyfikatorem, nie pod kluczem wiersza.
	przedrostekRegulyBiblioteki        = "regula-"
	przedrostekWpisuAudytuBiblioteki   = "audyt-"
	przedrostekPolitykiBiblioteki      = "retencja-"
	przedrostekUtrwaleniaBiblioteki    = "utrw-"
	przedrostekUdostepnieniaBiblioteki = "udost-"
	przedrostekWebhookaBiblioteki      = "hook-"
	przedrostekSugestiiBiblioteki      = "sug-"
)

// adapterBiblioteki wypełnia część portu Biblioteka: łączy repozytorium
// modułu z magazynem treści na dysku, bo baza trzyma wyłącznie odwołanie do
// bajtów, nie same bajty.
type adapterBiblioteki struct {
	repozytorium dane.RepozytoriumBiblioteki
	magazyn      *magazynTresciBiblioteki
	// Uruchamiacz, rozstrzygacz i katalog to trójka arsenału; kanały zasilają
	// klasyfikację wsadową.
	kanaly       *models.Rejestr
	uruchamiacz  session.Uruchamiacz
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
}

// ZNarzedziami wpina arsenał serwerowy — ten sam wzorzec, którym idzie Studio
// (`adapter_modul_studio.go`).
func (a *adapterBiblioteki) ZNarzedziami(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterBiblioteki {

	a.uruchamiacz, a.rozstrzygacz, a.katalog = uruchamiacz, rozstrzygacz, katalog
	return a
}

// nowyAdapterBiblioteki wiąże adapter z repozytorium modułu i wpina magazyn
// treści na katalogu domyślnym, żeby konstruktor nigdy nie oddał adaptera bez
// magazynu.
func nowyAdapterBiblioteki(repozytorium dane.RepozytoriumBiblioteki) *adapterBiblioteki {
	return &adapterBiblioteki{
		repozytorium: repozytorium,
		magazyn:      nowyMagazynTresciBiblioteki(konfiguracja.KatalogDanychDomyslny()),
	}
}

// ZKatalogiemDanych przestawia magazyn treści na katalog danych wskazany
// konfiguracją procesu i sprząta magazyn synchronicznie, zanim moduł zacznie
// obsługiwać komendy.
func (a *adapterBiblioteki) ZKatalogiemDanych(katalog string) *adapterBiblioteki {
	if katalog != "" {
		a.magazyn = nowyMagazynTresciBiblioteki(katalog)
	}
	posprzatajMagazynTresci(context.Background(), a.repozytorium, a.magazyn)
	return a
}

// Wgraj obsługuje `library.file.upload`: zapisuje treść w magazynie i wiersz
// pliku odwołaniem do niej, dopiero po udanym utrwaleniu bajtów, żeby wykaz
// nie pokazywał pliku bez treści.
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

	// Pierwsza wersja towarzyszy wgraniu, żeby wykaz wersji nie zaczynał się pusty.
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
	// Indeks treści zasila się przy zapisie, nie przy wyszukiwaniu.
	a.zaindeksujTresc(ctx, plik)

	kontrakt, err := a.zloz(ctx, plik)
	if err != nil {
		return shared.LibraryFileUploadResponse{}, err
	}
	return shared.LibraryFileUploadResponse{File: kontrakt}, nil
}

// plikZBiezacaWersja odkłada nadany identyfikator wersji na wiersz pliku
// przed ponownym zapisem, który aktualizuje wskaźnik wersji bieżącej, zamiast
// zakładać drugi wiersz.
func plikZBiezacaWersja(plik dane.PlikBiblioteki, wersjaID int64) dane.PlikBiblioteki {
	plik.WersjaBiezacaID = &wersjaID
	return plik
}

// Wykaz obsługuje `library.file.list`. Fraza zawęża wykaz po nazwie pliku;
// szukanie w treści ma własną komendę (`Szukaj`) i własny indeks.
func (a *adapterBiblioteki) Wykaz(ctx context.Context,
	z shared.LibraryFileListRequest) (shared.LibraryFileListResponse, error) {

	// Wykaz domyślny pokazuje wyłącznie zasoby czynne; archiwum ma własną drogę odczytu.
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

// Szukaj obsługuje `library.file.search`: fraza trafia w nazwę albo w treść
// (indeks FTS5, `dane/library.go` → `Szukaj`). Dopasowania po znaczeniu tu nie
// ma.
func (a *adapterBiblioteki) Szukaj(ctx context.Context,
	z shared.LibraryFileSearchRequest) (shared.LibraryFileSearchResponse, error) {

	// Przycięcie frazy poprzedza sprawdzenie pustki, żeby biały znak nie
	// trafił do indeksu jako wzorzec.
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

// Podglad obsługuje `library.file.preview`: buduje się z bieżącej wersji
// pliku przy każdym żądaniu, odmawia bez odwołania do treści, a obraz i tekst
// czyta ten sam czytelnik spod tego odwołania.
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
		// Podgląd tekstowy czyta początek pliku spod odwołania i oddaje go
		// z granicą i znacznikiem obcięcia.
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
	// Odwołanie wychodzi wyłącznie przy obrazie i wyłącznie w postaci względnej.
	podglad.ImageRef = odwolanieObrazuPodgladu(rodzaj, *plik.TrescOdwolanie)
	if z.Page != nil {
		podglad.Page = z.Page
	}
	return shared.LibraryFilePreviewResponse{Preview: podglad}, nil
}

// filtrZadania przekłada zawężenia żądania kontraktu na filtr warstwy danych.
// Wspólne dla `Wykaz` i `Szukaj`, bo oba zawężają po tych samych polach.
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

// plik odnajduje plik po kodzie i nazywa brak wprost: wskazanie bytu, którego
// nie ma, jest błędem żądania, nie awarią rdzenia.
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

// zlozWiele składa listę plików kontraktu z wierszy repozytorium, wspólną
// drogą, którą idą komendy `Wykaz` i `Szukaj`.
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

// zloz składa plik kontraktu z wiersza repozytorium wraz z etykietami,
// kolekcjami czytanymi z bazy i identyfikatorem wersji bieżącej; pole `path`
// niesie ścieżkę wewnątrz repozytorium, nie ścieżkę źródłową z maszyny Operatora.
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

// idWersjiBiezacej odnajduje kod wersji bieżącej pliku wśród jego historii,
// bo pole wewnętrzne niesie klucz wiersza, a kontrakt chce kod, którego
// repozytorium nie odczytuje wprost po samym kluczu.
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

// bladBiblioteki znakuje usterkę kodem kontraktu, żeby okno modułu pokazało
// powód, a nie samo „nie udało się". Błąd, któremu kod już nadano, przechodzi
// bez zmiany; dopiero usterka bez kodu staje się usterką wewnętrzną rdzenia.
func bladBiblioteki(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaBiblioteki nazywa brak danych w żądaniu jako błąd żądania, nie
// awarię rdzenia, i znakuje go kodem odrzucenia walidacji.
func bladWskazaniaBiblioteki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "moduł Library: "+powod))
}

// bladNieznanegoPlikuBiblioteki odróżnia brak pliku o wskazanym kodzie od
// nieudanego odczytu repozytorium i znakuje każdy przypadek osobnym kodem.
func bladNieznanegoPlikuBiblioteki(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Library: plik nie istnieje: "+kod))
	}
	return bladBiblioteki(err)
}

// bladNieznanejEtykietyBiblioteki odróżnia „etykiety nie ma w słowniku" od
// awarii odczytu — słownik etykiet (`adapter_modul_library_slownik.go`).
func bladNieznanejEtykietyBiblioteki(nazwa string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Library: etykieta nie istnieje: "+nazwa))
	}
	return bladBiblioteki(err)
}

// bladNieznanejReguly odróżnia brak reguły o wskazanym kodzie od awarii
// odczytu repozytorium i znakuje każdy przypadek osobnym kodem kontraktu.
func bladNieznanejReguly(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Library: reguła nie istnieje: "+kod))
	}
	return bladBiblioteki(err)
}

// bladNieznanejSugestii odróżnia brak sugestii o wskazanym kodzie od awarii
// odczytu repozytorium i znakuje każdy przypadek osobnym kodem kontraktu.
func bladNieznanejSugestii(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Library: sugestia nie istnieje: "+kod))
	}
	return bladBiblioteki(err)
}
