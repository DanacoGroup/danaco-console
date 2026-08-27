// Pakiet obsługuje rodzinę komend `extension.*` katalogu aplikacji:
// wyszukiwarkę, kartę szczegółów, kolekcje, rejestr organizacji, sprawdzanie
// aktualizacji, przesyłkę paczki, wersjonowanie, instalację zestawu, dziennik
// cyklu życia i operacje zbiorcze.
package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów rodziny `extension.*`: kolekcji, wpisu
// historii, paczki, webhooka, odwzorowania i ramki protokołu.
const (
	przedrostekKolekcjiRozszerzen  = "kolek-"
	przedrostekHistoriiRozszerzen  = "hist-"
	przedrostekPaczkiRozszerzenia  = "pacz-"
	przedrostekWebhookaRozszerzen  = "hook-"
	przedrostekMapowaniaRozszerzen = "mapa-"
	przedrostekRamkiProtokolu      = "ramka-"
)

// granicaWykazuRozszerzen jest górną granicą liczby pozycji strony wyniku
// wyszukiwania, stosowaną przy braku jawnego wskazania w żądaniu.
const granicaWykazuRozszerzen = 100

// Szukaj obsługuje komendę `extension.search`: przeszukuje pozycje katalogu
// po nazwie, opisie, kodzie i nazwach ich narzędzi, wraz z podpowiedziami.
func (a *adapterRozszerzen) Szukaj(ctx context.Context,
	z shared.ExtensionSearchRequest) (shared.ExtensionSearchResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionSearchResponse{}, err
	}
	fraza := strings.ToLower(strings.TrimSpace(z.Query))

	filtr := dane.FiltrRozszerzen{TylkoZainstalowane: wartoscLogiczna(z.InstalledOnly)}
	if z.Kind != nil {
		rodzaj, err := rodzajRozszerzenia(string(*z.Kind))
		if err != nil {
			return shared.ExtensionSearchResponse{}, err
		}
		filtr.Rodzaj = rodzaj
	}
	wiersze, err := a.rejestr.Rozszerzenia(ctx, filtr)
	if err != nil {
		return shared.ExtensionSearchResponse{}, bladRozszerzenia(err)
	}
	if z.Origin != nil {
		pochodzenie, err := zrodloPochodzenia(z.Origin)
		if err != nil {
			return shared.ExtensionSearchResponse{}, err
		}
		zawezone := make([]dane.Rozszerzenie, 0, len(wiersze))
		for _, wiersz := range wiersze {
			if wiersz.ZrodloPochodzenia != string(pochodzenie) {
				continue
			}
			zawezone = append(zawezone, wiersz)
		}
		wiersze = zawezone
	}

	trafienia := make([]shared.Extension, 0, len(wiersze))
	podpowiedzi := map[string]struct{}{}
	for _, wiersz := range wiersze {
		narzedzia, err := a.rejestr.NarzedziaRozszerzenia(ctx, wiersz.Identyfikator, "")
		if err != nil {
			return shared.ExtensionSearchResponse{}, bladRozszerzenia(err)
		}
		if fraza == "" || pasujeDoFrazyRozszerzenia(wiersz, narzedzia, fraza) {
			trafienia = append(trafienia, rozszerzenieKontraktu(wiersz))
			podpowiedzi[wiersz.Nazwa] = struct{}{}
			for _, narzedzie := range narzedzia {
				if fraza == "" || strings.Contains(strings.ToLower(narzedzie.Nazwa), fraza) {
					podpowiedzi[narzedzie.Nazwa] = struct{}{}
				}
			}
		}
	}

	razem := len(trafienia)
	// Odsunięcie i granica liczone po zawężeniu wyników do zapytania.
	odsuniecie := 0
	if z.Offset != nil && *z.Offset > 0 {
		odsuniecie = *z.Offset
	}
	if odsuniecie > len(trafienia) {
		odsuniecie = len(trafienia)
	}
	trafienia = trafienia[odsuniecie:]
	granica := granicaWykazuRozszerzen
	if z.Limit != nil && *z.Limit > 0 {
		granica = *z.Limit
	}
	if len(trafienia) > granica {
		trafienia = trafienia[:granica]
	}

	lista := make([]string, 0, len(podpowiedzi))
	for nazwa := range podpowiedzi {
		lista = append(lista, nazwa)
	}
	sort.Strings(lista)
	if len(lista) > 20 {
		lista = lista[:20]
	}
	return shared.ExtensionSearchResponse{
		Extensions: trafienia, Total: razem, Suggestions: lista,
	}, nil
}

// pasujeDoFrazyRozszerzenia rozstrzyga trafienie po polach, które pozycja
// niesie: kod, nazwa, rodzaj, opis oraz nazwa i opis jej narzędzi.
func pasujeDoFrazyRozszerzenia(wiersz dane.Rozszerzenie,
	narzedzia []dane.NarzedzieRozszerzenia, fraza string) bool {

	pola := []string{wiersz.Kod, wiersz.Nazwa, wiersz.Rodzaj, wartoscTekstu(wiersz.Opis)}
	for _, narzedzie := range narzedzia {
		pola = append(pola, narzedzie.Nazwa, wartoscTekstu(narzedzie.Opis))
	}
	for _, pole := range pola {
		if strings.Contains(strings.ToLower(pole), fraza) {
			return true
		}
	}
	return false
}

// PobierzSzczegol obsługuje komendę `extension.detail.get`: pełną kartę
// pozycji wraz z narzędziami, uprawnieniami, wersjami i podpisem.
func (a *adapterRozszerzen) PobierzSzczegol(ctx context.Context,
	z shared.ExtensionDetailGetRequest) (shared.ExtensionDetailGetResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionDetailGetResponse{}, err
	}

	narzedzia, err := a.rejestr.NarzedziaRozszerzenia(ctx, wiersz.Identyfikator, "")
	if err != nil {
		return shared.ExtensionDetailGetResponse{}, bladRozszerzenia(err)
	}
	uprawnienia, err := a.rejestr.UprawnieniaRozszerzenia(ctx, wiersz.Identyfikator)
	if err != nil {
		return shared.ExtensionDetailGetResponse{}, bladRozszerzenia(err)
	}
	wersje, err := a.rejestr.WersjeRozszerzenia(ctx, wiersz.Identyfikator)
	if err != nil {
		return shared.ExtensionDetailGetResponse{}, bladRozszerzenia(err)
	}

	szczegol := shared.ExtensionDetail{ExtensionId: wiersz.Identyfikator}
	for _, narzedzie := range narzedzia {
		szczegol.Tools = append(szczegol.Tools, narzedzieKontraktuRozszerzen(narzedzie))
	}
	for _, uprawnienie := range uprawnienia {
		if uprawnienie.Nadane {
			continue
		}
		szczegol.Permissions = append(szczegol.Permissions, uprawnienieKontraktuRozszerzen(uprawnienie))
	}
	// Dziennik zmian składa się z wpisów wersji niosących opis zmian.
	zmiany := make([]string, 0, len(wersje))
	for _, wersja := range wersje {
		if wersja.DziennikZmian == nil || *wersja.DziennikZmian == "" {
			continue
		}
		zmiany = append(zmiany, wersja.Wersja+": "+*wersja.DziennikZmian)
	}
	if len(zmiany) > 0 {
		szczegol.Changelog = wskaznikNapisuApp(strings.Join(zmiany, "\n"))
	}
	if podpis, err := a.rejestr.PodpisRozszerzenia(ctx, wiersz.Identyfikator); err == nil {
		sygnatura := sygnaturaKontraktu(podpis, wiersz)
		szczegol.Signature = &sygnatura
		szczegol.Publisher = podpis.Wydawca
	} else if !errors.Is(err, dane.ErrBrakWiersza) {
		return shared.ExtensionDetailGetResponse{}, bladRozszerzenia(err)
	}

	// Zależności i znaczniki niesie konfiguracja pozycji, nie osobne pola.
	szczegol.Dependencies = listaZKonfiguracjiRozszerzenia(wiersz.Konfiguracja, "dependencies")
	szczegol.Tags = listaZKonfiguracjiRozszerzenia(wiersz.Konfiguracja, "tags")
	if adres := napisZKonfiguracjiRozszerzenia(wiersz.Konfiguracja, "homepageUrl"); adres != "" {
		szczegol.HomepageUrl = &adres
	}

	return shared.ExtensionDetailGetResponse{Detail: szczegol}, nil
}

// WypiszKolekcje obsługuje komendę `extension.collection.list`: wskazaną
// kolekcję albo, bez wskazania, wykaz wszystkich kolekcji kuratorskich.
func (a *adapterRozszerzen) WypiszKolekcje(ctx context.Context,
	z shared.ExtensionCollectionListRequest) (shared.ExtensionCollectionListResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionCollectionListResponse{}, err
	}
	kod := strings.TrimSpace(wartoscTekstu(z.CollectionId))
	if kod != "" {
		kolekcja, err := a.rejestr.KolekcjaRozszerzen(ctx, kod)
		if err != nil {
			return shared.ExtensionCollectionListResponse{},
				bladNieznanegoBytuRozszerzenia("kolekcja", kod, err)
		}
		return shared.ExtensionCollectionListResponse{
			Collections: []shared.ExtensionCollection{kolekcjaKontraktuRozszerzen(kolekcja)}, Total: 1,
		}, nil
	}
	wiersze, err := a.rejestr.KolekcjeRozszerzen(ctx)
	if err != nil {
		return shared.ExtensionCollectionListResponse{}, bladRozszerzenia(err)
	}
	kolekcje := make([]shared.ExtensionCollection, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kolekcje = append(kolekcje, kolekcjaKontraktuRozszerzen(wiersz))
	}
	return shared.ExtensionCollectionListResponse{Collections: kolekcje, Total: len(kolekcje)}, nil
}

// ZapiszKolekcje obsługuje komendę `extension.collection.save`: zakłada nową
// kolekcję kuratorską albo zapisuje zmianę kolekcji zastanej.
func (a *adapterRozszerzen) ZapiszKolekcje(ctx context.Context,
	z shared.ExtensionCollectionSaveRequest) (shared.ExtensionCollectionSaveResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionCollectionSaveResponse{}, err
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.ExtensionCollectionSaveResponse{}, bladWskazaniaRozszerzenia(
			"kolekcja bez nazwy")
	}
	kod := strings.TrimSpace(wartoscTekstu(z.CollectionId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekKolekcjiRozszerzen)
	} else if _, err := a.rejestr.KolekcjaRozszerzen(ctx, kod); err != nil {
		return shared.ExtensionCollectionSaveResponse{},
			bladNieznanegoBytuRozszerzenia("kolekcja", kod, err)
	}

	zapisana, err := a.rejestr.ZapiszKolekcjeRozszerzen(ctx, dane.KolekcjaRozszerzen{
		Kod: kod, Nazwa: z.Name, Opis: z.Description, OznaczenieBarwne: z.ColorTag,
		KodyPozycji: z.ExtensionIds, Zaktualizowano: a.teraz(),
	})
	if err != nil {
		return shared.ExtensionCollectionSaveResponse{}, bladRozszerzenia(err)
	}
	return shared.ExtensionCollectionSaveResponse{Collection: kolekcjaKontraktuRozszerzen(zapisana)}, nil
}

// ZastosujKolekcje obsługuje komendę `extension.collection.apply`: przestawia
// stan włączenia każdej pozycji kolekcji, do końca wykazu bez przerywania.
func (a *adapterRozszerzen) ZastosujKolekcje(ctx context.Context,
	z shared.ExtensionCollectionApplyRequest) (shared.ExtensionCollectionApplyResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionCollectionApplyResponse{}, err
	}
	kolekcja, err := a.rejestr.KolekcjaRozszerzen(ctx, strings.TrimSpace(z.CollectionId))
	if err != nil {
		return shared.ExtensionCollectionApplyResponse{},
			bladNieznanegoBytuRozszerzenia("kolekcja", z.CollectionId, err)
	}
	wlacz := true
	if z.Enable != nil {
		wlacz = *z.Enable
	}

	zastosowane := []shared.Extension{}
	odrzucone := []shared.ExtensionRejection{}
	for _, kod := range kolekcja.KodyPozycji {
		zmieniona, err := a.rejestr.ZmienRozszerzenie(ctx, kod,
			dane.ZmianaRozszerzenia{Wlaczone: wartoscLogicznaApp(wlacz)}, a.teraz())
		if err != nil {
			odrzucone = append(odrzucone, shared.ExtensionRejection{
				ExtensionId: wskaznikNapisuApp(kod), Code: kod,
				Reason: "pozycji nie da się przestawić: " + err.Error(),
			})
			continue
		}
		pozycja := rozszerzenieKontraktu(zmieniona)
		a.rozglosRozszerzenie(shared.ChangeKindUpdated, pozycja)
		a.odnotujCyklZycia(ctx, zmieniona, czynnoscWlaczenia(wlacz), nil, nil,
			"zmiana grupowa kolekcji "+kolekcja.Nazwa)
		zastosowane = append(zastosowane, pozycja)
	}
	return shared.ExtensionCollectionApplyResponse{Applied: zastosowane, Rejected: odrzucone}, nil
}

// WypiszRejestr obsługuje komendę `extension.registry.list`: pozycje
// katalogu powstałe z publikacji pakietu, czyli prywatny rejestr organizacji.
func (a *adapterRozszerzen) WypiszRejestr(ctx context.Context,
	z shared.ExtensionRegistryListRequest) (shared.ExtensionRegistryListResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionRegistryListResponse{}, err
	}
	filtr := dane.FiltrRozszerzen{}
	if z.Kind != nil {
		rodzaj, err := rodzajRozszerzenia(string(*z.Kind))
		if err != nil {
			return shared.ExtensionRegistryListResponse{}, err
		}
		filtr.Rodzaj = rodzaj
	}
	wiersze, err := a.rejestr.Rozszerzenia(ctx, filtr)
	if err != nil {
		return shared.ExtensionRegistryListResponse{}, bladRozszerzenia(err)
	}
	fraza := strings.ToLower(strings.TrimSpace(wartoscTekstu(z.Query)))

	pozycje := []shared.Extension{}
	for _, wiersz := range wiersze {
		// Pozycja rejestru to ta powstała z publikacji pakietu.
		if napisZKonfiguracjiRozszerzenia(wiersz.Konfiguracja, "appsPackageId") == "" {
			continue
		}
		if fraza != "" && !pasujeDoFrazyRozszerzenia(wiersz, nil, fraza) {
			continue
		}
		pozycje = append(pozycje, rozszerzenieKontraktu(wiersz))
	}
	// Adres rejestru jest lokalny, bo rejestr jest lokalny.
	return shared.ExtensionRegistryListResponse{
		Extensions: pozycje, Total: len(pozycje),
		RegistryUrl: wskaznikNapisuApp(korzenWytworowApp),
	}, nil
}

// SprawdzAktualizacje obsługuje komendę `extension.update.check`: zestawia
// wersję zainstalowaną z najwyższą wersją zapisaną w dzienniku wersji.
func (a *adapterRozszerzen) SprawdzAktualizacje(ctx context.Context,
	z shared.ExtensionUpdateCheckRequest) (shared.ExtensionUpdateCheckResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionUpdateCheckResponse{}, err
	}
	wiersze := []dane.Rozszerzenie{}
	kod := strings.TrimSpace(wartoscTekstu(z.ExtensionId))
	if kod != "" {
		wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, kod)
		if err != nil {
			return shared.ExtensionUpdateCheckResponse{}, err
		}
		wiersze = append(wiersze, wiersz)
	} else {
		wszystkie, err := a.rejestr.Rozszerzenia(ctx, dane.FiltrRozszerzen{TylkoZainstalowane: true})
		if err != nil {
			return shared.ExtensionUpdateCheckResponse{}, bladRozszerzenia(err)
		}
		wiersze = wszystkie
	}

	aktualizacje := []shared.ExtensionUpdate{}
	for _, wiersz := range wiersze {
		wersje, err := a.rejestr.WersjeRozszerzenia(ctx, wiersz.Identyfikator)
		if err != nil {
			return shared.ExtensionUpdateCheckResponse{}, bladRozszerzenia(err)
		}
		biezaca := wartoscTekstu(wiersz.Wersja)
		najwyzsza, dziennik := najwyzszaWersjaRozszerzenia(wersje)
		if najwyzsza == "" || najwyzsza == biezaca {
			continue
		}
		if porownajWersjeSemantyczne(najwyzsza, biezaca) <= 0 {
			continue
		}
		lamie := lamieZgodnoscSemantyczna(biezaca, najwyzsza)
		aktualizacje = append(aktualizacje, shared.ExtensionUpdate{
			ExtensionId: wiersz.Identyfikator, CurrentVersion: biezaca,
			AvailableVersion: najwyzsza, Breaking: &lamie,
			Changelog: wskaznikNapisuApp(dziennik),
		})
	}
	return shared.ExtensionUpdateCheckResponse{
		Updates: aktualizacje, CheckedAt: a.teraz(),
	}, nil
}

// PrzeslijPaczke obsługuje `extension.package.upload`: kładzie bajty w magazynie
// treści rdzenia i oddaje odwołanie, którym woła się potem instalację.
func (a *adapterRozszerzen) PrzeslijPaczke(ctx context.Context,
	z shared.ExtensionPackageUploadRequest) (shared.ExtensionPackageUploadResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionPackageUploadResponse{}, err
	}
	if a.magazyn == nil {
		return shared.ExtensionPackageUploadResponse{}, bladBrakuMagazynuApp(
			"przesłanie paczki rozszerzenia")
	}
	if strings.TrimSpace(z.FileName) == "" {
		return shared.ExtensionPackageUploadResponse{}, bladWskazaniaRozszerzenia(
			"przesyłka bez nazwy pliku")
	}
	bajty, err := base64.StdEncoding.DecodeString(z.ContentBase64)
	if err != nil {
		return shared.ExtensionPackageUploadResponse{}, bladWskazaniaRozszerzenia(
			"treść przesyłki nie jest poprawnym base64: " + err.Error())
	}
	if len(bajty) == 0 {
		return shared.ExtensionPackageUploadResponse{}, bladWskazaniaRozszerzenia(
			"przesyłka jest pusta — nie ma czego zapisać")
	}

	suma := sumaTresciApp(bajty)
	// Suma podana w żądaniu jest sprawdzana, nie przyjmowana bezkrytycznie.
	if z.ChecksumSha256 != nil && *z.ChecksumSha256 != "" &&
		!strings.EqualFold(*z.ChecksumSha256, suma) {
		return shared.ExtensionPackageUploadResponse{}, bladWskazaniaRozszerzenia(
			"suma kontrolna przesyłki nie zgadza się z jej treścią: podano " +
				*z.ChecksumSha256 + ", policzono " + suma)
	}

	sciezka, err := a.magazyn.Zapisz(bajty, suma)
	if err != nil {
		return shared.ExtensionPackageUploadResponse{}, bladRozszerzenia(err)
	}
	kod := nowyIdentyfikator(przedrostekPaczkiRozszerzenia)
	if err := a.rejestr.ZalozPaczkeRozszerzenia(ctx, dane.PaczkaRozszerzenia{
		Kod: kod, NazwaPliku: z.FileName, Sciezka: odwolanieWytworuApp(sciezka),
		Rozmiar: int64(len(bajty)), SumaKontrolna: suma, Utworzono: a.teraz(),
	}); err != nil {
		return shared.ExtensionPackageUploadResponse{}, bladRozszerzenia(err)
	}
	return shared.ExtensionPackageUploadResponse{
		UploadRef: kod, SizeBytes: int64(len(bajty)),
	}, nil
}

// PrzypnijWersje obsługuje `extension.version.pin`. Brak `version` zdejmuje
// przypięcie — kontrakt mówi to wprost.
func (a *adapterRozszerzen) PrzypnijWersje(ctx context.Context,
	z shared.ExtensionVersionPinRequest) (shared.ExtensionVersionPinResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionVersionPinResponse{}, err
	}
	wersja := strings.TrimSpace(wartoscTekstu(z.Version))
	if wersja != "" {
		// Przypiąć da się wyłącznie wersję, która istnieje w dzienniku.
		wersje, err := a.rejestr.WersjeRozszerzenia(ctx, wiersz.Identyfikator)
		if err != nil {
			return shared.ExtensionVersionPinResponse{}, bladRozszerzenia(err)
		}
		if !czyWersjaZnana(wersje, wersja, wartoscTekstu(wiersz.Wersja)) {
			return shared.ExtensionVersionPinResponse{}, bladWskazaniaRozszerzenia(
				"pozycja " + wiersz.Kod + " nie ma wersji " + wersja +
					" — przypiąć da się wyłącznie wersję zarejestrowaną")
		}
	}
	if err := a.rejestr.PrzypnijWersjeRozszerzenia(ctx, wiersz.Identyfikator, wersja, a.teraz()); err != nil {
		return shared.ExtensionVersionPinResponse{}, bladRozszerzenia(err)
	}
	odswiezona, err := a.rejestr.Rozszerzenie(ctx, wiersz.Identyfikator)
	if err != nil {
		return shared.ExtensionVersionPinResponse{}, bladRozszerzenia(err)
	}
	pozycja := rozszerzenieKontraktu(odswiezona)
	a.rozglosRozszerzenie(shared.ChangeKindUpdated, pozycja)
	a.odnotujCyklZycia(ctx, odswiezona, shared.ExtensionLifecycleActionConfigured, nil, nil,
		"przypięcie wersji: "+wersjaAlboBrak(wersja))
	return shared.ExtensionVersionPinResponse{
		Extension: pozycja, PinnedVersion: wskaznikNapisuApp(wersja),
	}, nil
}

// CofnijWersje obsługuje `extension.version.rollback`: przywraca pozycję do
// wersji zarejestrowanej wcześniej i odnotowuje to w dzienniku cyklu życia.
func (a *adapterRozszerzen) CofnijWersje(ctx context.Context,
	z shared.ExtensionVersionRollbackRequest) (shared.ExtensionVersionRollbackResponse, error) {

	wiersz, err := a.pozycjaRozszerzeniaZadania(ctx, z.ExtensionId)
	if err != nil {
		return shared.ExtensionVersionRollbackResponse{}, err
	}
	cel := strings.TrimSpace(z.TargetVersion)
	if cel == "" {
		return shared.ExtensionVersionRollbackResponse{}, bladWskazaniaRozszerzenia(
			"cofnięcie bez wskazania wersji docelowej")
	}
	wersje, err := a.rejestr.WersjeRozszerzenia(ctx, wiersz.Identyfikator)
	if err != nil {
		return shared.ExtensionVersionRollbackResponse{}, bladRozszerzenia(err)
	}
	if !czyWersjaZnana(wersje, cel, "") {
		return shared.ExtensionVersionRollbackResponse{}, bladWskazaniaRozszerzenia(
			"pozycja " + wiersz.Kod + " nigdy nie miała wersji " + cel +
				" — cofnąć da się wyłącznie do wersji zarejestrowanej")
	}
	poprzednia := wartoscTekstu(wiersz.Wersja)
	if poprzednia == cel {
		return shared.ExtensionVersionRollbackResponse{}, bladWskazaniaRozszerzenia(
			"pozycja " + wiersz.Kod + " już stoi w wersji " + cel)
	}

	zmieniona, err := a.rejestr.ZmienRozszerzenie(ctx, wiersz.Identyfikator,
		dane.ZmianaRozszerzenia{Wersja: &cel}, a.teraz())
	if err != nil {
		return shared.ExtensionVersionRollbackResponse{}, bladRozszerzenia(err)
	}
	pozycja := rozszerzenieKontraktu(zmieniona)
	a.rozglosRozszerzenie(shared.ChangeKindUpdated, pozycja)
	a.odnotujCyklZycia(ctx, zmieniona, shared.ExtensionLifecycleActionRolledBack,
		&poprzednia, &cel, "cofnięcie do wersji zarejestrowanej")
	return shared.ExtensionVersionRollbackResponse{
		Extension: pozycja, RolledBackFromVersion: poprzednia,
	}, nil
}

// ZainstalujZestaw obsługuje `extension.bundle.install`: instaluje komplet
// pozycji z jednego manifestu i oddaje osobno te, których nie dało się założyć.
func (a *adapterRozszerzen) ZainstalujZestaw(ctx context.Context,
	z shared.ExtensionBundleInstallRequest) (shared.ExtensionBundleInstallResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionBundleInstallResponse{}, err
	}
	var manifest struct {
		Extensions []struct {
			Code        string          `json:"code"`
			Kind        string          `json:"kind"`
			Name        string          `json:"name"`
			Version     string          `json:"version"`
			Origin      string          `json:"origin"`
			Config      json.RawMessage `json:"config"`
			Description string          `json:"description"`
		} `json:"extensions"`
	}
	if err := json.Unmarshal(z.Manifest, &manifest); err != nil {
		return shared.ExtensionBundleInstallResponse{}, bladWskazaniaRozszerzenia(
			"manifest zestawu nie ma kształtu {\"extensions\":[…]}: " + err.Error())
	}
	if len(manifest.Extensions) == 0 {
		return shared.ExtensionBundleInstallResponse{}, bladWskazaniaRozszerzenia(
			"manifest zestawu nie wymienia ani jednej pozycji")
	}
	wlacz := false
	if z.Enable != nil {
		wlacz = *z.Enable
	}

	zainstalowane := []shared.Extension{}
	odrzucone := []shared.ExtensionRejection{}
	for _, wpis := range manifest.Extensions {
		kod := strings.TrimSpace(wpis.Code)
		if kod == "" {
			odrzucone = append(odrzucone, shared.ExtensionRejection{
				Code: "", Reason: "pozycja manifestu bez kodu",
			})
			continue
		}
		rodzaj, err := rodzajRozszerzenia(wpis.Kind)
		if err != nil {
			odrzucone = append(odrzucone, shared.ExtensionRejection{
				Code: kod, Reason: "nieznany rodzaj pozycji: " + wpis.Kind,
			})
			continue
		}
		pochodzenie := shared.ExtensionOriginPersonal
		if wpis.Origin != "" {
			wskazane, err := zrodloPochodzenia((*shared.ExtensionOrigin)(&wpis.Origin))
			if err != nil {
				odrzucone = append(odrzucone, shared.ExtensionRejection{
					Code: kod, Reason: "nieznane pochodzenie pozycji: " + wpis.Origin,
				})
				continue
			}
			pochodzenie = string(wskazane)
		}

		nazwa := wpis.Name
		if nazwa == "" {
			nazwa = kod
		}
		konfiguracja := "{}"
		if len(wpis.Config) > 0 && json.Valid(wpis.Config) {
			konfiguracja = string(wpis.Config)
		}

		zastane, err := a.rejestr.RozszerzeniePoKodzie(ctx, kod)
		if err == nil {
			zmiana := dane.ZmianaRozszerzenia{
				Nazwa: &nazwa, Zainstalowane: wartoscLogicznaApp(true),
				Wlaczone: wartoscLogicznaApp(wlacz), Konfiguracja: &konfiguracja,
			}
			if wpis.Version != "" {
				zmiana.Wersja = &wpis.Version
			}
			if wpis.Description != "" {
				zmiana.Opis = &wpis.Description
			}
			zmieniona, err := a.rejestr.ZmienRozszerzenie(ctx, zastane.Identyfikator, zmiana, a.teraz())
			if err != nil {
				odrzucone = append(odrzucone, shared.ExtensionRejection{
					ExtensionId: wskaznikNapisuApp(zastane.Identyfikator), Code: kod,
					Reason: "nie można odświeżyć pozycji: " + err.Error(),
				})
				continue
			}
			pozycja := rozszerzenieKontraktu(zmieniona)
			a.rozglosRozszerzenie(shared.ChangeKindUpdated, pozycja)
			a.odnotujCyklZycia(ctx, zmieniona, shared.ExtensionLifecycleActionUpdated,
				zastane.Wersja, zmieniona.Wersja, "instalacja z manifestu zestawu")
			zainstalowane = append(zainstalowane, pozycja)
			continue
		}
		if !errors.Is(err, dane.ErrBrakWiersza) {
			odrzucone = append(odrzucone, shared.ExtensionRejection{
				Code: kod, Reason: "nie można odczytać pozycji: " + err.Error(),
			})
			continue
		}

		zalozona, err := a.rejestr.ZalozRozszerzenie(ctx, dane.Rozszerzenie{
			Identyfikator: nowyIdentyfikator(przedrostekRozszerzenia),
			Kod:           kod, Rodzaj: rodzaj, Nazwa: nazwa,
			Opis:          wskaznikNapisuApp(wpis.Description),
			Wersja:        wskaznikNapisuApp(wpis.Version),
			Zainstalowane: true, Wlaczone: wlacz,
			ZrodloPochodzenia: pochodzenie, Konfiguracja: konfiguracja,
			Zaktualizowano: a.teraz(),
		})
		if err != nil {
			odrzucone = append(odrzucone, shared.ExtensionRejection{
				Code: kod, Reason: "nie można założyć pozycji: " + err.Error(),
			})
			continue
		}
		pozycja := rozszerzenieKontraktu(zalozona)
		a.rozglosRozszerzenie(shared.ChangeKindCreated, pozycja)
		a.odnotujCyklZycia(ctx, zalozona, shared.ExtensionLifecycleActionInstalled,
			nil, zalozona.Wersja, "instalacja z manifestu zestawu")
		zainstalowane = append(zainstalowane, pozycja)
	}
	return shared.ExtensionBundleInstallResponse{
		Installed: zainstalowane, Rejected: odrzucone,
	}, nil
}

// WypiszHistorie obsługuje komendę `extension.history.list`: wpisy dziennika
// cyklu życia pozycji wskazanej albo, bez wskazania, wszystkich pozycji.
func (a *adapterRozszerzen) WypiszHistorie(ctx context.Context,
	z shared.ExtensionHistoryListRequest) (shared.ExtensionHistoryListResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionHistoryListResponse{}, err
	}
	kod := strings.TrimSpace(wartoscTekstu(z.ExtensionId))
	if kod != "" {
		if _, err := a.pozycjaRozszerzeniaZadania(ctx, kod); err != nil {
			return shared.ExtensionHistoryListResponse{}, err
		}
	}
	od := int64(0)
	if z.Since != nil {
		od = *z.Since
	}
	granica := 0
	if z.Limit != nil {
		granica = *z.Limit
	}
	wiersze, razem, err := a.rejestr.HistoriaRozszerzenia(ctx, kod, od, granica)
	if err != nil {
		return shared.ExtensionHistoryListResponse{}, bladRozszerzenia(err)
	}
	wpisy := make([]shared.ExtensionHistoryEntry, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wpisy = append(wpisy, shared.ExtensionHistoryEntry{
			Id: wiersz.Kod, ExtensionId: wiersz.RozszerzenieKod,
			Action:      shared.ExtensionLifecycleAction(wiersz.Czynnosc),
			FromVersion: wiersz.WersjaPrzed, ToVersion: wiersz.WersjaPo,
			Detail: wiersz.Szczegol, OccurredAt: wiersz.Zaszlo,
		})
	}
	return shared.ExtensionHistoryListResponse{Entries: wpisy, Total: razem}, nil
}

// WykonajZbiorczo obsługuje `extension.admin.bulk`: włącza, wyłącza albo
// odinstalowuje wskazane pozycje. Odinstalowanie NIE kasuje wiersza — zdejmuje
// stan zainstalowania i włączenia, tak samo jak `extension.uninstall`.
func (a *adapterRozszerzen) WykonajZbiorczo(ctx context.Context,
	z shared.ExtensionAdminBulkRequest) (shared.ExtensionAdminBulkResponse, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return shared.ExtensionAdminBulkResponse{}, err
	}
	if len(z.ExtensionIds) == 0 {
		return shared.ExtensionAdminBulkResponse{}, bladWskazaniaRozszerzenia(
			"operacja zbiorcza bez ani jednej pozycji")
	}
	var zmiana dane.ZmianaRozszerzenia
	var czynnosc shared.ExtensionLifecycleAction
	switch z.Action {
	case shared.ExtensionBulkActionEnable:
		zmiana = dane.ZmianaRozszerzenia{Wlaczone: wartoscLogicznaApp(true)}
		czynnosc = shared.ExtensionLifecycleActionEnabled
	case shared.ExtensionBulkActionDisable:
		zmiana = dane.ZmianaRozszerzenia{Wlaczone: wartoscLogicznaApp(false)}
		czynnosc = shared.ExtensionLifecycleActionDisabled
	case shared.ExtensionBulkActionUninstall:
		zmiana = dane.ZmianaRozszerzenia{
			Zainstalowane: wartoscLogicznaApp(false), Wlaczone: wartoscLogicznaApp(false),
		}
		czynnosc = shared.ExtensionLifecycleActionUninstalled
	default:
		return shared.ExtensionAdminBulkResponse{}, bladWskazaniaRozszerzenia(
			"nieznana czynność zbiorcza " + strconv.Quote(string(z.Action)) +
				" — dopuszczalne: enable, disable, uninstall")
	}

	dotkniete := []shared.Extension{}
	odrzucone := []shared.ExtensionRejection{}
	for _, kod := range z.ExtensionIds {
		zmieniona, err := a.rejestr.ZmienRozszerzenie(ctx, kod, zmiana, a.teraz())
		if err != nil {
			odrzucone = append(odrzucone, shared.ExtensionRejection{
				ExtensionId: wskaznikNapisuApp(kod), Code: kod,
				Reason: "nie można zmienić pozycji: " + err.Error(),
			})
			continue
		}
		pozycja := rozszerzenieKontraktu(zmieniona)
		a.rozglosRozszerzenie(shared.ChangeKindUpdated, pozycja)
		a.odnotujCyklZycia(ctx, zmieniona, czynnosc, nil, nil, "operacja zbiorcza rejestru")
		dotkniete = append(dotkniete, pozycja)
	}
	return shared.ExtensionAdminBulkResponse{Affected: dotkniete, Rejected: odrzucone}, nil
}

// odnotujCyklZycia dopisuje wpis dziennika. Nieudany zapis nie przewraca
// czynności, której wpis dotyczy — dziennik opisuje pracę, nie warunkuje jej.
func (a *adapterRozszerzen) odnotujCyklZycia(ctx context.Context, wiersz dane.Rozszerzenie,
	czynnosc shared.ExtensionLifecycleAction, przed, po *string, szczegol string) {

	_ = a.rejestr.DopiszHistorieRozszerzenia(ctx, dane.WpisHistoriiRozszerzenia{
		Kod:             nowyIdentyfikator(przedrostekHistoriiRozszerzen),
		RozszerzenieKod: wiersz.Identyfikator, Czynnosc: string(czynnosc),
		WersjaPrzed: przed, WersjaPo: po,
		Szczegol: wskaznikNapisuApp(szczegol), Zaszlo: a.teraz(),
	})
}

// czynnoscWlaczenia nazywa czynność cyklu życia po stronie przełącznika
// włączenia: włączona albo wyłączona, zgodnie z kontraktem historii.
func czynnoscWlaczenia(wlacz bool) shared.ExtensionLifecycleAction {
	if wlacz {
		return shared.ExtensionLifecycleActionEnabled
	}
	return shared.ExtensionLifecycleActionDisabled
}

// wersjaAlboBrak nazywa pustkę słowem „zdjęte", bo pusty łańcuch w dzienniku
// czytałby się jak brak zapisu, a nie jak zdjęcie przypięcia wersji.
func wersjaAlboBrak(wersja string) string {
	if wersja == "" {
		return "zdjęte"
	}
	return wersja
}

// czyWersjaZnana rozstrzyga, czy pozycja kiedykolwiek miała wskazaną wersję.
// Wersja bieżąca liczy się jako znana także wtedy, gdy nie ma jeszcze wiersza
// historii — pozycja zarejestrowana z wersją, a bez wydania, nadal ją ma.
func czyWersjaZnana(wersje []dane.WersjaRozszerzenia, szukana, biezaca string) bool {
	if szukana == biezaca && biezaca != "" {
		return true
	}
	for _, wersja := range wersje {
		if wersja.Wersja == szukana {
			return true
		}
	}
	return false
}

// najwyzszaWersjaRozszerzenia zwraca najwyższą wersję z wykazu wersji pozycji
// wraz z jej dziennikiem zmian, porównując wersje semantycznie.
func najwyzszaWersjaRozszerzenia(wersje []dane.WersjaRozszerzenia) (string, string) {
	najwyzsza := ""
	dziennik := ""
	for _, wersja := range wersje {
		if najwyzsza == "" || porownajWersjeSemantyczne(wersja.Wersja, najwyzsza) > 0 {
			najwyzsza = wersja.Wersja
			dziennik = wartoscTekstu(wersja.DziennikZmian)
		}
	}
	return najwyzsza, dziennik
}

// porownajWersjeSemantyczne porównuje dwie wersje po członach liczbowych.
// Wersja nieliczbowa porównuje się tekstowo — to nie jest pełny semver
// i nie udaje, że nim jest: rozstrzyga wyłącznie „nowsza czy nie".
func porownajWersjeSemantyczne(pierwsza, druga string) int {
	czlonyPierwszej := czlonyWersji(pierwsza)
	czlonyDrugiej := czlonyWersji(druga)
	for indeks := 0; indeks < 3; indeks++ {
		if czlonyPierwszej[indeks] != czlonyDrugiej[indeks] {
			if czlonyPierwszej[indeks] > czlonyDrugiej[indeks] {
				return 1
			}
			return -1
		}
	}
	return strings.Compare(pierwsza, druga)
}

// czlonyWersji rozkłada zapis wersji na trzy człony liczbowe, pomijając
// przedrostek `v` i przyrostki; brak członu w zapisie daje zero.
func czlonyWersji(wersja string) [3]int {
	var czlony [3]int
	rdzen, _, _ := strings.Cut(strings.TrimPrefix(wersja, "v"), "-")
	rdzen, _, _ = strings.Cut(rdzen, "+")
	for indeks, czesc := range strings.SplitN(rdzen, ".", 3) {
		if indeks > 2 {
			break
		}
		liczba, err := strconv.Atoi(strings.TrimSpace(czesc))
		if err != nil {
			continue
		}
		czlony[indeks] = liczba
	}
	return czlony
}

// lamieZgodnoscSemantyczna rozstrzyga, czy zmiana wersji łamie zgodność.
// Zmiana pierwszego członu łamie ją zawsze; przy wersji zerowej łamie ją także
// drugi człon, bo tak stanowi semantyka wersji rozwojowych.
func lamieZgodnoscSemantyczna(przed, po string) bool {
	stare := czlonyWersji(przed)
	nowe := czlonyWersji(po)
	if stare[0] != nowe[0] {
		return true
	}
	return nowe[0] == 0 && stare[1] != nowe[1]
}

// listaZKonfiguracjiRozszerzenia wyciąga wykaz napisów spod wskazanego klucza
// dokumentu JSON konfiguracji pozycji, albo pustkę, gdy klucza nie ma.
func listaZKonfiguracjiRozszerzenia(konfiguracja, klucz string) []string {
	if strings.TrimSpace(konfiguracja) == "" {
		return nil
	}
	var pola map[string]json.RawMessage
	if err := json.Unmarshal([]byte(konfiguracja), &pola); err != nil {
		return nil
	}
	surowe, jest := pola[klucz]
	if !jest {
		return nil
	}
	var wykaz []string
	if err := json.Unmarshal(surowe, &wykaz); err != nil {
		return nil
	}
	return wykaz
}

// napisZKonfiguracjiRozszerzenia wyciąga pojedynczy napis spod wskazanego
// klucza dokumentu JSON konfiguracji pozycji, albo pustkę, gdy klucza nie ma.
func napisZKonfiguracjiRozszerzenia(konfiguracja, klucz string) string {
	if strings.TrimSpace(konfiguracja) == "" {
		return ""
	}
	var pola map[string]json.RawMessage
	if err := json.Unmarshal([]byte(konfiguracja), &pola); err != nil {
		return ""
	}
	surowe, jest := pola[klucz]
	if !jest {
		return ""
	}
	var wartosc string
	if err := json.Unmarshal(surowe, &wartosc); err != nil {
		return ""
	}
	return wartosc
}

// kolekcjaKontraktu przekłada wiersz kolekcji z bazy danych na kształt
// kolekcji zwracany kontraktem komunikacji.
func kolekcjaKontraktuRozszerzen(wiersz dane.KolekcjaRozszerzen) shared.ExtensionCollection {
	pozycje := wiersz.KodyPozycji
	if pozycje == nil {
		pozycje = []string{}
	}
	return shared.ExtensionCollection{
		Id: wiersz.Kod, Name: wiersz.Nazwa, Description: wiersz.Opis,
		ColorTag: wiersz.OznaczenieBarwne, ExtensionIds: pozycje,
		UpdatedAt: wiersz.Zaktualizowano,
	}
}

// narzedzieKontraktu przekłada wiersz narzędzia z bazy danych na kształt
// narzędzia zwracany kontraktem komunikacji.
func narzedzieKontraktuRozszerzen(wiersz dane.NarzedzieRozszerzenia) shared.ExtensionToolEntry {
	wpis := shared.ExtensionToolEntry{
		Name: wiersz.Nazwa, Kind: shared.ExtensionToolKind(wiersz.Rodzaj),
		Description: wiersz.Opis, Uri: wiersz.Adres,
	}
	if wiersz.SchematWejscia != nil && json.Valid([]byte(*wiersz.SchematWejscia)) {
		wpis.InputSchema = json.RawMessage(*wiersz.SchematWejscia)
	}
	return wpis
}

// uprawnienieKontraktu przekłada wiersz uprawnienia z bazy danych na kształt
// uprawnienia zwracany kontraktem komunikacji.
func uprawnienieKontraktuRozszerzen(wiersz dane.UprawnienieRozszerzenia) shared.ExtensionPermission {
	uprawnienie := shared.ExtensionPermission{
		Scope:  shared.ExtensionPermissionScope(wiersz.Zakres),
		Target: wiersz.Byt, Explanation: wiersz.Objasnienie, GrantedAt: wiersz.Nadano,
	}
	if wiersz.Tryb != nil {
		tryb := shared.AccessMode(*wiersz.Tryb)
		uprawnienie.Mode = &tryb
	}
	return uprawnienie
}

// pozycjaRozszerzeniaZadania odczytuje pozycję po `Extension.id` i nazywa jej
// brak. Wszystkie komendy rodziny wskazujące pozycję idą tędy, żeby odmowa
// „nie ma takiej pozycji" brzmiała wszędzie tak samo.
func (a *adapterRozszerzen) pozycjaRozszerzeniaZadania(ctx context.Context,
	kod string) (dane.Rozszerzenie, error) {

	if err := a.sprawdzKatalog(); err != nil {
		return dane.Rozszerzenie{}, err
	}
	przyciety := strings.TrimSpace(kod)
	if przyciety == "" {
		return dane.Rozszerzenie{}, bladWskazaniaRozszerzenia("żądanie bez pozycji katalogu")
	}
	wiersz, err := a.rejestr.Rozszerzenie(ctx, przyciety)
	if err != nil {
		return dane.Rozszerzenie{}, bladNieznanegoBytuRozszerzenia("pozycja katalogu", przyciety, err)
	}
	return wiersz, nil
}

// bladNieznanegoBytuRozszerzenia odróżnia odpowiedź „bytu nie ma" od usterki
// odczytu, nazywając rodzaj bytu i jego kod w treści odmowy.
func bladNieznanegoBytuRozszerzenia(nazwaBytu, kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"katalog rozszerzeń: "+nazwaBytu+" nie istnieje: "+kod))
	}
	return bladRozszerzenia(err)
}

// sygnaturaKontraktu składa wynik weryfikacji podpisu w kształcie kontraktu.
// Poziom zaufania wynika z dwóch faktów naraz: pochodzenia pozycji i tego, czy
// podpis przeszedł weryfikację — sam napis wydawcy niczego nie potwierdza.
func sygnaturaKontraktu(podpis dane.PodpisRozszerzenia,
	wiersz dane.Rozszerzenie) shared.ExtensionSignature {

	sygnatura := shared.ExtensionSignature{
		Signed:         podpis.PodpisBase64 != nil && *podpis.PodpisBase64 != "",
		Algorithm:      podpis.Algorytm,
		ChecksumSha256: podpis.SumaKontrolna,
		Publisher:      podpis.Wydawca,
	}
	sygnatura.TrustLevel = shared.ExtensionTrustLevel(shared.ExtensionTrustLevelUnverifiedPersonal)
	if wiersz.ZrodloPochodzenia == shared.ExtensionOriginDanaco {
		sygnatura.TrustLevel = shared.ExtensionTrustLevel(shared.ExtensionTrustLevelDanacoPlugin)
	}
	return sygnatura
}
