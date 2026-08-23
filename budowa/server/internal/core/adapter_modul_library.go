// Moduł Library — wgranie, wykaz, podgląd i wyszukiwanie pliku repozytorium
// wiedzy. Wersje mają własny plik (`adapter_modul_library_wersje.go`), kolekcje
// i etykiety swoje (`adapter_modul_library_kolekcje.go`,
// `adapter_modul_library_uchwyty.go` — tam też stoi port `Biblioteka`
// i `zarejestrujBiblioteke`). Typ adaptera i konstruktor stoją tutaj, bo metody
// rozłożone po tych plikach wiszą na jednym `*adapterBiblioteki`.
//
// `library.file.search` dopasowuje frazę do nazwy pliku albo do jego treści —
// ta druga przez indeks pełnotekstowy FTS5 (`dane/biblioteka_indeks_tresci.go`),
// zasilany przy każdym zapisie treści (`adapter_modul_library_indeks.go`).
// Bajty leżą poza bazą (`tresc_odwolanie`), indeks jest ich odtwarzalnym
// wyciągiem tekstowym; dopasowanie jest trafieniem w słowo, nie w znaczenie.
// `library.file.list` zawęża wykaz po nazwie i nie szuka w dokumentach.
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

// adapterBiblioteki wypełnia część portu Biblioteka. Zależności są dwie:
// wspólne repozytorium modułu, rozłożone po stronie danych na kilka plików wedle
// odpowiedzialności, ale niosące jeden typ, oraz magazyn treści na dysku
// (`adapter_modul_library_magazyn.go`) — baza trzyma wyłącznie odwołanie do
// treści, więc same bajty muszą mieć gdzie leżeć.
type adapterBiblioteki struct {
	repozytorium dane.RepozytoriumBiblioteki
	magazyn      *magazynTresciBiblioteki
	// Uruchamiacz, rozstrzygacz i katalog roboczy są trójką arsenału: jedyną
	// drogą, którą rdzeń woła binarium spoza siebie
	// (`zewnetrzne.Wolaj`). Moduł używa jej dokładnie w jednym miejscu — pomiar
	// czasu trwania nagrania przy odczycie metadanych osadzonych
	// (`adapter_modul_library_technika.go`). Brak trójki nie wyłącza modułu:
	// pole czasu trwania po prostu nie wchodzi do odpowiedzi.
	// Kanały modelu zasilają klasyfikację wsadową
	// (`adapter_modul_library_sugestie.go`). Ten sam rejestr, którym jedzie okno
	// rozmowy — drugiego silnika modelu w rdzeniu nie ma.
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
// treści oparty o katalog danych rdzenia — tak samo, jak `nowyAdapterKont`
// wpina sejf poświadczeń (`handlers_konta_adapter.go`). Bez magazynu plik wgrany
// samym `contentBase64` nie ma odwołania i podgląd go odmawia.
//
// Katalog jest tu domyślny (`konfiguracja.KatalogDanychDomyslny`);
// obowiązujący — przestawiony przełącznikiem
// `-dane` albo zmienną `DANACO_KATALOG_DANYCH` — zna wyłącznie montaż
// (`Montaz.Konfiguracja.KatalogDanych`) i podaje go przez `ZKatalogiemDanych`
// przy składaniu portu (`montaz_porty.go`). Wartość domyślna stoi tu dla
// wywołania bez montażu, żeby konstruktor nigdy nie oddał adaptera bez magazynu.
func nowyAdapterBiblioteki(repozytorium dane.RepozytoriumBiblioteki) *adapterBiblioteki {
	return &adapterBiblioteki{
		repozytorium: repozytorium,
		magazyn:      nowyMagazynTresciBiblioteki(konfiguracja.KatalogDanychDomyslny()),
	}
}

// ZKatalogiemDanych przestawia magazyn treści na katalog danych wskazany
// konfiguracją procesu. Wołane przez montaż, który jako jedyny zna rzeczywisty
// katalog; bez wywołania magazyn stoi na katalogu domyślnym.
//
// Tu wypada też sprzątanie magazynu (`adapter_modul_library_sprzatanie.go`): to
// jedyna chwila startu rdzenia, w której moduł zna już swój prawdziwy katalog
// danych i jeszcze nie obsługuje żadnej komendy — tak samo sprząta kosz sesji
// (`trwalosc_kosza.go`). Sprzątanie idzie synchronicznie: obchód katalogu jest
// tani wobec odtworzenia stanu, a rdzeń przyjmujący wgrania w trakcie
// przemiatania widziałby wykaz żywych odwołań sprzed nich.
func (a *adapterBiblioteki) ZKatalogiemDanych(katalog string) *adapterBiblioteki {
	if katalog != "" {
		a.magazyn = nowyMagazynTresciBiblioteki(katalog)
	}
	posprzatajMagazynTresci(context.Background(), a.repozytorium, a.magazyn)
	return a
}

// Wgraj obsługuje `library.file.upload`. Schemat trzyma treść poza bazą
// (`tresc_odwolanie`), więc wgranie kończy się odwołaniem do treści, nigdy
// treścią w wierszu.
//
// Obie drogi wgrania prowadzą do jednego magazynu: treść przysłana base64
// i treść wciągnięta spod `sourcePath` lądują w magazynie treści rdzenia
// (`adapter_modul_library_tresc.go`, `adapter_modul_library_magazyn.go`),
// a odwołaniem jest ścieżka bloba. Ścieżka źródłowa zostaje przy pliku
// w kolumnie `sciezka` — mówi, skąd plik przyszedł, i nie jest wskaźnikiem na
// treść żywą, którą ktoś z zewnątrz mógłby nadpisać po wgraniu.
//
// Nieudany zapis treści odmawia całego wgrania: wiersz pliku powstaje dopiero po
// utrwaleniu bajtów, inaczej wykaz pokazywałby plik, za którym nie ma nic.
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

	// Pierwsza wersja towarzyszy wgraniu — Versioning Panel ma się od czego
	// zacząć, zamiast pokazywać pusty wykaz aż do pierwszej zmiany.
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
	// Indeks treści zasila się przy zapisie, nie przy wyszukiwaniu
	// (`adapter_modul_library_indeks.go`) — bez tego `library.file.search`
	// wraca do dopasowywania samej nazwy pliku.
	a.zaindeksujTresc(ctx, plik)

	kontrakt, err := a.zloz(ctx, plik)
	if err != nil {
		return shared.LibraryFileUploadResponse{}, err
	}
	return shared.LibraryFileUploadResponse{File: kontrakt}, nil
}

// plikZBiezacaWersja odkłada nadany identyfikator wersji na wiersz pliku
// przed ponownym zapisem — `ZapiszPlik` nadpisuje przez `identyfikator_zewnetrzny`
// (ON CONFLICT), więc drugie wywołanie ze zmienionym `WersjaBiezacaID` po
// prostu aktualizuje wskaźnik, nie zakłada drugiego wiersza.
func plikZBiezacaWersja(plik dane.PlikBiblioteki, wersjaID int64) dane.PlikBiblioteki {
	plik.WersjaBiezacaID = &wersjaID
	return plik
}

// Wykaz obsługuje `library.file.list`. Fraza zawęża wykaz po nazwie pliku;
// szukanie w treści ma własną komendę (`Szukaj`) i własny indeks.
func (a *adapterBiblioteki) Wykaz(ctx context.Context,
	z shared.LibraryFileListRequest) (shared.LibraryFileListResponse, error) {

	// Wykaz domyślny pokazuje zasoby czynne. Archiwum ma własną drogę
	// (`library.file.restore` po kodzie zasobu i pulpit stanu), a wykaz, który
	// pokazuje zarchiwizowane obok czynnych, znosi sens kosza repozytorium.
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

	// Sprawdzenie pustki patrzy na tę samą, przyciętą postać frazy, która idzie
	// do indeksu FTS5: fraza z samych spacji nie jest równa "", więc bez
	// przycięcia zeszłaby do FTS5 jako pusty wzorzec MATCH i wróciłaby błędem
	// SQL zamiast odmową żądania.
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

// Podglad obsługuje `library.file.preview`. Podgląd nie ma własnej tabeli —
// buduje się z bieżącej wersji pliku przy każdym żądaniu. Plik bez
// `tresc_odwolanie` wraca odmową; gdy odwołanie jest, podgląd tekstowy czyta
// spod niego treść (`trescPodgladuBiblioteki`), a obraz oddaje odwołanie w polu
// `imageRef`.
//
// Podgląd nie rozróżnia dróg wgrania: odwołanie jest ścieżką na dysku niezależnie
// od tego, czy wskazuje `sourcePath`, czy bloba w magazynie treści rdzenia
// (`Wgraj`, `adapter_modul_library_magazyn.go`), więc czytelnik jest jeden.
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
		// Podgląd tekstowy czyta treść spod odwołania: kontrakt ma dla niego
		// pole `text` wraz z granicą `maxChars` i znacznikiem `truncated`, więc
		// rdzeń odczytuje początek pliku i oddaje go jako treść. Odmowa idzie
		// dopiero wtedy, gdy odwołania nie da się odczytać
		// (`bladOdczytuTresciBiblioteki`).
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
	// Odwołanie wychodzi wyłącznie przy obrazie i wyłącznie w postaci względnej;
	// oba zawężenia stoją przy `odwolanieObrazuPodgladu`
	// (`adapter_modul_library_podglad.go`).
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

// zlozWiele składa listę plików kontraktu z wierszy repozytorium — wspólne
// dla `Wykaz` i `Szukaj`.
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
// kolekcjami i identyfikatorem wersji bieżącej.
//
// Kolekcje czyta się z bazy (`KolekcjePliku`,
// `dane/biblioteka_kolekcje_pliku.go`), nie z żądania — dzięki temu każda
// odpowiedź niosąca plik mówi tę samą przynależność, niezależnie od drogi.
//
// Pole `path` kontraktu niesie ścieżkę WEWNĄTRZ repozytorium
// (`sciezka_repozytorium`, migracja 180), którą nadaje `library.file.move` — nie
// ścieżkę źródłową z maszyny Operatora. Ta druga zostaje w kolumnie `sciezka`
// jako prowenancja wewnętrzna i z rdzenia nie wychodzi: wyniosłaby na zewnątrz
// układ cudzego dysku wraz z nazwami katalogów.
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

// idWersjiBiezacej odnajduje kod wersji bieżącej pliku wśród jego historii.
// `PlikBiblioteki.WersjaBiezacaID` niesie klucz wewnętrzny wiersza, a kontrakt
// chce kod (`LibraryFile.versionId`) — jedyny sposób przełożenia jednego na
// drugie to przejrzeć wykaz wersji tego pliku (repozytorium nie ma odczytu
// wersji po samym kluczu).
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

// bladWskazaniaBiblioteki nazywa brak danych w żądaniu — błąd żądania, nie rdzenia.
func bladWskazaniaBiblioteki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "moduł Library: "+powod))
}

// bladNieznanegoPlikuBiblioteki odróżnia „pliku nie ma" od „odczyt się nie powiódł".
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

// bladNieznanejReguly odróżnia „reguły nie ma" od awarii odczytu.
func bladNieznanejReguly(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Library: reguła nie istnieje: "+kod))
	}
	return bladBiblioteki(err)
}

// bladNieznanejSugestii odróżnia „sugestii nie ma" od awarii odczytu.
func bladNieznanejSugestii(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Library: sugestia nie istnieje: "+kod))
	}
	return bladBiblioteki(err)
}
