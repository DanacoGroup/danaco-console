// Pakiet obsługuje rodzinę komend `apps.package.*`: budowę pakietu z artefaktu
// wdrożenia, zapis manifestu, walidację zgodności z kontraktem rozszerzenia,
// podpis Ed25519 i publikację pozycji w rejestrze rozszerzeń organizacji.
package core

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekPakietuApp znakuje identyfikatory pakietów rozszerzenia, wraz
// z którymi rdzeń zakłada wiersz pakietu w repozytorium modułu Apps.
const przedrostekPakietuApp = "pak-"

// nazwaManifestuWPakiecieApp jest nazwą, pod którą manifest pakietu ląduje
// w archiwum obok reszty wpisów przeniesionych z artefaktu wdrożenia.
const nazwaManifestuWPakiecieApp = "manifest.json"

// wzorzecWersjiSemantycznejApp sprawdza wersję manifestu wobec postaci, którą
// kontrakt nazywa wprost wersją semantyczną: liczba główna, poboczna i łatka.
var wzorzecWersjiSemantycznejApp = regexp.MustCompile(
	`^\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)

// wzorzecIdentyfikatoraPakietuApp sprawdza tożsamość rozszerzenia: małe litery,
// cyfry, kropkę, myślnik i podkreślenie, zaczynając od litery albo cyfry.
var wzorzecIdentyfikatoraPakietuApp = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

// Kody zastrzeżeń, którymi walidator pakietu oznacza niezgodność z kontraktem
// rozszerzenia — każdy kod odpowiada jednemu rodzajowi zastrzeżenia.
const (
	kodPakietuBezManifestuApp   = "packageWithoutManifest"
	kodPakietuTozsamoscApp      = "manifestIdentifierInvalid"
	kodPakietuWersjaApp         = "manifestVersionNotSemver"
	kodPakietuRodzajApp         = "manifestKindUnknown"
	kodPakietuBezNarzedziApp    = "manifestWithoutTools"
	kodPakietuSzerokiZakresApp  = "manifestBroadPermissions"
	kodPakietuBezPodpisuApp     = "packageWithoutSignature"
	kodPakietuBezArchiwumApp    = "packageWithoutArchive"
	kodPakietuZaleznoscPustaApp = "manifestDependencyEmpty"
)

// ZbudujPakiet obsługuje `apps.package.build`: rozpakowuje artefakt wdrożenia
// i składa z niego nowe archiwum wraz z manifestem, odłożone w magazynie treści.
func (a *adapterAplikacji) ZbudujPakiet(ctx context.Context,
	z shared.AppsPackageBuildRequest) (shared.AppsPackageBuildResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.package.build")
	if err != nil {
		return shared.AppsPackageBuildResponse{}, err
	}
	if a.magazyn == nil {
		return shared.AppsPackageBuildResponse{}, bladBrakuMagazynuApp("pakowanie produktu")
	}
	format := shared.AppPackageFormat(shared.AppPackageFormatZip)
	if z.Format != nil {
		if err := sprawdzFormatPakietuApp(*z.Format); err != nil {
			return shared.AppsPackageBuildResponse{}, err
		}
		format = *z.Format
	}

	artefakt, err := a.artefaktPakowaniaApp(ctx, okno, z.ArtifactRef)
	if err != nil {
		return shared.AppsPackageBuildResponse{}, err
	}
	wejscie, err := a.wpisyArtefaktuApp(artefakt)
	if err != nil {
		return shared.AppsPackageBuildResponse{}, bladAplikacji(err)
	}

	kod := nowyIdentyfikator(przedrostekPakietuApp)
	// Manifest zaczyna od tożsamości znanej rdzeniowi: kodu pakietu i nazwy produktu.
	manifest := shared.AppPackageManifest{
		Identifier: kod,
		Name:       kod,
		Version:    "0.0.0",
		Kind:       shared.ExtensionKind(shared.ExtensionKindPlugin),
	}
	if produkt, err := a.repozytorium.ProduktApp(ctx, okno); err == nil {
		manifest.Name = produkt.Nazwa
		manifest.Description = produkt.Opis
	} else if !isBrakWierszaApp(err) {
		return shared.AppsPackageBuildResponse{}, bladAplikacji(err)
	}

	archiwum, err := zlozArchiwumPakietuApp(format, wejscie, manifest)
	if err != nil {
		return shared.AppsPackageBuildResponse{}, bladAplikacji(err)
	}
	odwolanie, rozmiar, err := a.wniesDoMagazynuApp(archiwum)
	if err != nil {
		return shared.AppsPackageBuildResponse{}, bladAplikacji(err)
	}
	tresc, err := json.Marshal(manifest)
	if err != nil {
		return shared.AppsPackageBuildResponse{}, bladAplikacji(err)
	}

	zapisany, err := a.repozytorium.ZapiszPakietApp(ctx, dane.PakietApp{
		Kod: kod, Okno: okno, Manifest: wskaznikNapisuApp(string(tresc)),
		ArtefaktOdwolanie: wskaznikNapisuApp(artefakt.Kod), Format: string(format),
		Sciezka: &odwolanie, Rozmiar: &rozmiar,
	})
	if err != nil {
		return shared.AppsPackageBuildResponse{}, bladAplikacji(err)
	}
	a.dopiszDziennikApp(ctx, okno, artefakt.WdrozenieKod, nil,
		"pakiet "+kod+" zbudowany z artefaktu "+artefakt.Kod+" ("+string(format)+", "+
			strconv.FormatInt(rozmiar, 10)+" bajtów, wpisów: "+strconv.Itoa(len(wejscie))+")")

	pakiet, err := pakietKontraktuApp(zapisany)
	if err != nil {
		return shared.AppsPackageBuildResponse{}, bladAplikacji(err)
	}
	return shared.AppsPackageBuildResponse{Package: pakiet}, nil
}

// ZapiszManifestPakietu obsługuje `apps.package.manifest.save`. Brak
// `packageId` zakłada pakiet nowy — manifest może powstać, zanim cokolwiek
// zostanie spakowane; taki pakiet jeszcze bez archiwum wykaże to walidacja.
func (a *adapterAplikacji) ZapiszManifestPakietu(ctx context.Context,
	z shared.AppsPackageManifestSaveRequest) (shared.AppsPackageManifestSaveResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.package.manifest.save")
	if err != nil {
		return shared.AppsPackageManifestSaveResponse{}, err
	}
	kod := strings.TrimSpace(wartoscTekstu(z.PackageId))
	if kod == "" {
		kod = nowyIdentyfikator(przedrostekPakietuApp)
	} else if _, err := a.pakietOknaApp(ctx, okno, kod); err != nil {
		return shared.AppsPackageManifestSaveResponse{}, err
	}

	tresc, err := json.Marshal(z.Manifest)
	if err != nil {
		return shared.AppsPackageManifestSaveResponse{}, bladAplikacji(err)
	}
	zapisany, err := a.repozytorium.ZapiszPakietApp(ctx, dane.PakietApp{
		Kod: kod, Okno: okno, Manifest: wskaznikNapisuApp(string(tresc)),
	})
	if err != nil {
		return shared.AppsPackageManifestSaveResponse{}, bladAplikacji(err)
	}
	pakiet, err := pakietKontraktuApp(zapisany)
	if err != nil {
		return shared.AppsPackageManifestSaveResponse{}, bladAplikacji(err)
	}
	return shared.AppsPackageManifestSaveResponse{Package: pakiet}, nil
}

// SprawdzPakiet obsługuje `apps.package.validate`. Zastrzeżenia zwrócone
// walidacją są ostrzeżeniami, nie bramą — nie wstrzymują dalszej publikacji.
func (a *adapterAplikacji) SprawdzPakiet(ctx context.Context,
	z shared.AppsPackageValidateRequest) (shared.AppsPackageValidateResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.package.validate")
	if err != nil {
		return shared.AppsPackageValidateResponse{}, err
	}
	kod := strings.TrimSpace(z.PackageId)
	if kod == "" {
		return shared.AppsPackageValidateResponse{}, bladWskazaniaAplikacji(
			"apps.package.validate wymaga pakietu")
	}
	wiersz, err := a.pakietOknaApp(ctx, okno, kod)
	if err != nil {
		return shared.AppsPackageValidateResponse{}, err
	}
	return shared.AppsPackageValidateResponse{
		Issues:      zastrzezeniaPakietuApp(wiersz),
		ValidatedAt: time.Now().UTC().UnixMilli(),
	}, nil
}

// PodpiszPakiet obsługuje `apps.package.sign`: liczy sumę kontrolną archiwum,
// podpisuje ją kluczem wydawcy Ed25519 z sejfu poświadczeń i od razu weryfikuje
// podpis kluczem publicznym, zanim odda wynik.
func (a *adapterAplikacji) PodpiszPakiet(ctx context.Context,
	z shared.AppsPackageSignRequest) (shared.AppsPackageSignResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.package.sign")
	if err != nil {
		return shared.AppsPackageSignResponse{}, err
	}
	kod := strings.TrimSpace(z.PackageId)
	if kod == "" {
		return shared.AppsPackageSignResponse{}, bladWskazaniaAplikacji(
			"apps.package.sign wymaga pakietu")
	}
	odwolanieKlucza := strings.TrimSpace(z.SigningKeyRef)
	if odwolanieKlucza == "" {
		return shared.AppsPackageSignResponse{}, bladWskazaniaAplikacji(
			"apps.package.sign wymaga odwołania do klucza wydawcy")
	}
	if a.sejf == nil {
		return shared.AppsPackageSignResponse{}, bladBrakuMagazynuApp("podpisanie pakietu")
	}

	wiersz, err := a.pakietOknaApp(ctx, okno, kod)
	if err != nil {
		return shared.AppsPackageSignResponse{}, err
	}
	if wiersz.Sciezka == nil || *wiersz.Sciezka == "" {
		return shared.AppsPackageSignResponse{}, bladWskazaniaAplikacji(
			"pakiet " + kod + " nie ma jeszcze archiwum — podpisać da się bajty, nie zamiar " +
				"(zbuduj go komendą apps.package.build)")
	}
	bajty, err := os.ReadFile(sciezkaWMagazynieApp(a.katalogDanych, *wiersz.Sciezka))
	if err != nil {
		return shared.AppsPackageSignResponse{}, bladAplikacji(err)
	}

	klucz, err := a.kluczWydawcyApp(ctx, odwolanieKlucza)
	if err != nil {
		return shared.AppsPackageSignResponse{}, bladAplikacji(err)
	}
	suma := sha256.Sum256(bajty)
	podpis := ed25519.Sign(klucz, suma[:])
	publiczny, ok := klucz.Public().(ed25519.PublicKey)
	if !ok {
		return shared.AppsPackageSignResponse{}, bladAplikacji(
			fmt.Errorf("moduł Apps: klucz wydawcy %q nie ma części publicznej", odwolanieKlucza))
	}
	// Weryfikacja podpisu następuje od razu po jego złożeniu.
	potwierdzony := ed25519.Verify(publiczny, suma[:], podpis)

	sumaTekstem := hex.EncodeToString(suma[:])
	algorytm := "ed25519"
	wydawca := odwolanieKlucza
	sygnatura := shared.ExtensionSignature{
		Signed:         true,
		Verified:       potwierdzony,
		TrustLevel:     shared.ExtensionTrustLevel(shared.ExtensionTrustLevelVerifiedPublisher),
		Algorithm:      &algorytm,
		ChecksumSha256: &sumaTekstem,
		Publisher:      &wydawca,
	}
	if !potwierdzony {
		szczegol := "podpis powstał, ale nie przeszedł weryfikacji kluczem publicznym wydawcy"
		sygnatura.Detail = &szczegol
		sygnatura.TrustLevel = shared.ExtensionTrustLevel(shared.ExtensionTrustLevelUnverifiedPersonal)
	}

	// Podpis wchodzi do wiersza wraz z postacią bajtową, do ponownej weryfikacji.
	zapisany := sygnaturaZPodpisemApp(sygnatura, podpis, publiczny)
	tresc, err := json.Marshal(zapisany)
	if err != nil {
		return shared.AppsPackageSignResponse{}, bladAplikacji(err)
	}
	if _, err := a.repozytorium.ZapiszPakietApp(ctx, dane.PakietApp{
		Kod: kod, Okno: okno, Podpis: wskaznikNapisuApp(string(tresc)),
	}); err != nil {
		return shared.AppsPackageSignResponse{}, bladAplikacji(err)
	}
	a.dopiszDziennikApp(ctx, okno, nil, nil,
		"pakiet "+kod+" podpisany kluczem "+odwolanieKlucza+", weryfikacja: "+
			strconv.FormatBool(potwierdzony))

	return shared.AppsPackageSignResponse{Signature: sygnatura}, nil
}

// OpublikujPakiet obsługuje `apps.package.publish`: zapisuje pakiet jako
// pozycję prywatnego rejestru organizacji, tej samej tabeli rozszerzeń, którą
// prowadzi rodzina komend `extension.*`.
func (a *adapterAplikacji) OpublikujPakiet(ctx context.Context,
	z shared.AppsPackagePublishRequest) (shared.AppsPackagePublishResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.package.publish")
	if err != nil {
		return shared.AppsPackagePublishResponse{}, err
	}
	kod := strings.TrimSpace(z.PackageId)
	if kod == "" {
		return shared.AppsPackagePublishResponse{}, bladWskazaniaAplikacji(
			"apps.package.publish wymaga pakietu")
	}
	if z.Visibility != nil {
		if err := sprawdzWidocznoscPakietuApp(*z.Visibility); err != nil {
			return shared.AppsPackagePublishResponse{}, err
		}
	}
	if a.katalogRozszerzen == nil {
		return shared.AppsPackagePublishResponse{}, bladBrakuMagazynuApp(
			"publikacja do rejestru organizacji")
	}

	wiersz, err := a.pakietOknaApp(ctx, okno, kod)
	if err != nil {
		return shared.AppsPackagePublishResponse{}, err
	}
	if wiersz.Sciezka == nil || *wiersz.Sciezka == "" {
		return shared.AppsPackagePublishResponse{}, bladWskazaniaAplikacji(
			"pakiet " + kod + " nie ma archiwum — do rejestru trafia pakiet, nie zamiar")
	}
	manifest, err := manifestPakietuApp(wiersz)
	if err != nil {
		return shared.AppsPackagePublishResponse{}, bladAplikacji(err)
	}
	if manifest == nil {
		return shared.AppsPackagePublishResponse{}, bladWskazaniaAplikacji(
			"pakiet " + kod + " nie ma manifestu — pozycja katalogu bierze z niego tożsamość, " +
				"nazwę, rodzaj i wersję")
	}

	teraz := time.Now().UTC().UnixMilli()
	// Konfiguracja niesie wyłącznie pola odróżniające pozycję od instalowanej ręcznie.
	konfiguracja, err := json.Marshal(map[string]any{
		"appsPackageId":  kod,
		"appsWindowId":   okno,
		"archiveRef":     *wiersz.Sciezka,
		"visibility":     widocznoscPakietuApp(z.Visibility),
		"releaseNotes":   wartoscTekstu(z.ReleaseNotes),
		"publishedAtMs":  teraz,
		"signaturePlain": wiersz.Podpis != nil,
	})
	if err != nil {
		return shared.AppsPackagePublishResponse{}, bladAplikacji(err)
	}

	pozycja, err := a.wpiszPozycjeKataloguApp(ctx, *manifest, konfiguracja, teraz)
	if err != nil {
		return shared.AppsPackagePublishResponse{}, err
	}

	if _, err := a.repozytorium.ZapiszPakietApp(ctx, dane.PakietApp{
		Kod: kod, Okno: okno, RozszerzenieKod: &pozycja.Id,
	}); err != nil {
		return shared.AppsPackagePublishResponse{}, bladAplikacji(err)
	}
	// Podpis pakietu przechodzi na pozycję wraz z materiałem do ponownej weryfikacji.
	if err := a.przeniesPodpisNaPozycje(ctx, wiersz, pozycja.Id, teraz); err != nil {
		return shared.AppsPackagePublishResponse{}, err
	}
	a.dopiszDziennikApp(ctx, okno, nil, nil,
		"pakiet "+kod+" opublikowany jako pozycja katalogu "+pozycja.Code)

	return shared.AppsPackagePublishResponse{Extension: pozycja, PublishedAt: teraz}, nil
}

// wpiszPozycjeKataloguApp zakłada albo odświeża pozycję katalogu powstałą
// z manifestu: pozycja o tym samym kodzie manifestu jest tą samą pozycją
// w nowszym wydaniu, więc druga publikacja podnosi wersję zamiast zakładać
// bliźniaka.
func (a *adapterAplikacji) wpiszPozycjeKataloguApp(ctx context.Context,
	manifest shared.AppPackageManifest, konfiguracja []byte, teraz int64) (shared.Extension, error) {

	rodzaj := string(manifest.Kind)
	if err := sprawdzRodzajRozszerzeniaPakietuApp(manifest.Kind); err != nil {
		return shared.Extension{}, err
	}
	opis := manifest.Description
	wersja := manifest.Version

	zastana, err := a.katalogRozszerzen.RozszerzeniePoKodzie(ctx, manifest.Identifier)
	switch {
	case err == nil:
		zmieniona, err := a.katalogRozszerzen.ZmienRozszerzenie(ctx, zastana.Identyfikator,
			dane.ZmianaRozszerzenia{
				Nazwa: &manifest.Name, Opis: opis, Wersja: &wersja,
				Zainstalowane: wartoscLogicznaApp(true), Konfiguracja: wskaznikNapisuApp(string(konfiguracja)),
			}, teraz)
		if err != nil {
			return shared.Extension{}, bladAplikacji(err)
		}
		return rozszerzenieKontraktu(zmieniona), nil
	case errors.Is(err, dane.ErrBrakWiersza):
	default:
		return shared.Extension{}, bladAplikacji(err)
	}

	zalozona, err := a.katalogRozszerzen.ZalozRozszerzenie(ctx, dane.Rozszerzenie{
		Identyfikator: nowyIdentyfikator(przedrostekRozszerzenia),
		Kod:           manifest.Identifier,
		Rodzaj:        rodzaj,
		Nazwa:         manifest.Name,
		Opis:          opis,
		Wersja:        &wersja,
		Zainstalowane: true,
		// Pozycja z publikacji wchodzi wyłączona, jak każda pozycja spoza zestawu wbudowanego.
		Wlaczone:          false,
		ZrodloPochodzenia: shared.ExtensionOriginPersonal,
		Konfiguracja:      string(konfiguracja),
		Zaktualizowano:    teraz,
	})
	if err != nil {
		return shared.Extension{}, bladAplikacji(err)
	}
	return rozszerzenieKontraktu(zalozona), nil
}

// przeniesPodpisNaPozycje odkłada przy pozycji katalogu wersję manifestu oraz,
// gdy pakiet jest podpisany, materiał podpisu do jego ponownej weryfikacji.
func (a *adapterAplikacji) przeniesPodpisNaPozycje(ctx context.Context, wiersz dane.PakietApp,
	pozycjaKod string, teraz int64) error {

	if a.katalogRozszerzen == nil {
		return nil
	}
	// Wersja manifestu wchodzi do rejestru wersji pozycji zawsze, także przy pakiecie niepodpisanym.
	if manifest, err := manifestPakietuApp(wiersz); err == nil && manifest != nil {
		_ = a.katalogRozszerzen.ZapiszWersjeRozszerzenia(ctx, dane.WersjaRozszerzenia{
			RozszerzenieKod: pozycjaKod, Wersja: manifest.Version,
			PaczkaOdwolanie: wiersz.Sciezka, Utworzono: teraz,
		})
	}
	if wiersz.Podpis == nil || *wiersz.Podpis == "" {
		return nil
	}
	var zapisany struct {
		Algorithm       *string `json:"algorithm"`
		ChecksumSha256  *string `json:"checksumSha256"`
		Publisher       *string `json:"publisher"`
		SignatureBase64 string  `json:"signatureBase64"`
		PublicKeyBase64 string  `json:"publicKeyBase64"`
	}
	if err := json.Unmarshal([]byte(*wiersz.Podpis), &zapisany); err != nil {
		return bladAplikacji(err)
	}
	if err := a.katalogRozszerzen.ZapiszPodpisRozszerzenia(ctx, dane.PodpisRozszerzenia{
		RozszerzenieKod: pozycjaKod, Algorytm: zapisany.Algorithm,
		SumaKontrolna: zapisany.ChecksumSha256, Wydawca: zapisany.Publisher,
		PodpisBase64:   wskaznikNapisuApp(zapisany.SignatureBase64),
		KluczBase64:    wskaznikNapisuApp(zapisany.PublicKeyBase64),
		Zaktualizowano: teraz,
	}); err != nil {
		return bladAplikacji(err)
	}
	return nil
}

// artefaktPakowaniaApp rozwiązuje artefakt wejściowy: wskazany albo ostatni
// z wdrożenia udanego (kontrakt mówi to wprost).
func (a *adapterAplikacji) artefaktPakowaniaApp(ctx context.Context, okno string,
	wskazanie *string) (dane.ArtefaktApp, error) {

	kod := strings.TrimSpace(wartoscTekstu(wskazanie))
	if kod != "" {
		artefakt, err := a.repozytorium.ArtefaktApp(ctx, kod)
		if err != nil {
			return dane.ArtefaktApp{}, bladNieznanegoBytuApp("artefakt budowania", kod, err)
		}
		if artefakt.Okno != okno {
			return dane.ArtefaktApp{}, bladWskazaniaAplikacji(
				"artefakt " + kod + " należy do okna " + artefakt.Okno +
					", a żądanie przyszło z okna " + okno)
		}
		return artefakt, nil
	}
	artefakt, err := a.repozytorium.OstatniArtefaktUdanegoWdrozeniaApp(ctx, okno)
	if err != nil {
		if isBrakWierszaApp(err) {
			return dane.ArtefaktApp{}, bladWskazaniaAplikacji(
				"okno " + okno + " nie ma ani jednego artefaktu z udanego wdrożenia — " +
					"wskaż artefakt wprost albo uruchom wdrożenie (apps.deployment.run)")
		}
		return dane.ArtefaktApp{}, bladAplikacji(err)
	}
	return artefakt, nil
}

// wpisyArtefaktuApp rozpakowuje archiwum artefaktu do par ścieżka→treść.
// Silnik wdrożenia składa artefakt zawsze jako archiwum `zip`, więc pakowanie
// czyta tę postać wprost, bez zgadywania formatu.
func (a *adapterAplikacji) wpisyArtefaktuApp(artefakt dane.ArtefaktApp) (map[string][]byte, error) {
	bajty, err := os.ReadFile(sciezkaWMagazynieApp(a.katalogDanych, artefakt.Sciezka))
	if err != nil {
		return nil, fmt.Errorf("moduł Apps: nie można odczytać artefaktu %s: %w", artefakt.Kod, err)
	}
	czytnik, err := zip.NewReader(bytes.NewReader(bajty), int64(len(bajty)))
	if err != nil {
		return nil, fmt.Errorf("moduł Apps: artefakt %s nie jest archiwum zip: %w", artefakt.Kod, err)
	}
	wpisy := map[string][]byte{}
	for _, plik := range czytnik.File {
		if plik.FileInfo().IsDir() {
			continue
		}
		strumien, err := plik.Open()
		if err != nil {
			return nil, fmt.Errorf("moduł Apps: nieczytelny wpis %s artefaktu %s: %w",
				plik.Name, artefakt.Kod, err)
		}
		tresc, err := io.ReadAll(strumien)
		_ = strumien.Close()
		if err != nil {
			return nil, fmt.Errorf("moduł Apps: przerwany odczyt wpisu %s artefaktu %s: %w",
				plik.Name, artefakt.Kod, err)
		}
		wpisy[plik.Name] = tresc
	}
	return wpisy, nil
}

// zlozArchiwumPakietuApp składa archiwum w żądanym formacie wraz z manifestem,
// zastępując ewentualny manifest z wpisów artefaktu manifestem zbudowanym.
func zlozArchiwumPakietuApp(format shared.AppPackageFormat, wpisy map[string][]byte,
	manifest shared.AppPackageManifest) ([]byte, error) {

	tresc, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("moduł Apps: nie można złożyć manifestu pakietu: %w", err)
	}
	pelne := map[string][]byte{nazwaManifestuWPakiecieApp: tresc}
	for nazwa, bajty := range wpisy {
		if nazwa == nazwaManifestuWPakiecieApp {
			continue
		}
		pelne[nazwa] = bajty
	}
	// Kolejność wpisów jest ustalona, aby ta sama treść dawała tę samą sumę kontrolną.
	nazwy := make([]string, 0, len(pelne))
	for nazwa := range pelne {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)

	switch format {
	case shared.AppPackageFormatZip:
		return archiwumZipApp(nazwy, pelne)
	case shared.AppPackageFormatTargz:
		return archiwumTarGzApp(nazwy, pelne)
	}
	return nil, fmt.Errorf("moduł Apps: nieobsłużony format pakietu %q", format)
}

// archiwumZipApp składa listę wpisów w archiwum formatu `zip`, zapisując je
// w kolejności podanej listy nazw.
func archiwumZipApp(nazwy []string, wpisy map[string][]byte) ([]byte, error) {
	var bufor bytes.Buffer
	zapis := zip.NewWriter(&bufor)
	for _, nazwa := range nazwy {
		strumien, err := zapis.Create(nazwa)
		if err != nil {
			return nil, fmt.Errorf("moduł Apps: nie można założyć wpisu %s: %w", nazwa, err)
		}
		if _, err := strumien.Write(wpisy[nazwa]); err != nil {
			return nil, fmt.Errorf("moduł Apps: nie można zapisać wpisu %s: %w", nazwa, err)
		}
	}
	if err := zapis.Close(); err != nil {
		return nil, fmt.Errorf("moduł Apps: nie można domknąć archiwum zip: %w", err)
	}
	return bufor.Bytes(), nil
}

// archiwumTarGzApp składa listę wpisów w archiwum formatu `tar.gz`, zapisując
// je w kolejności podanej listy nazw.
func archiwumTarGzApp(nazwy []string, wpisy map[string][]byte) ([]byte, error) {
	var bufor bytes.Buffer
	kompresja := gzip.NewWriter(&bufor)
	zapis := tar.NewWriter(kompresja)
	for _, nazwa := range nazwy {
		tresc := wpisy[nazwa]
		naglowek := &tar.Header{
			Name: nazwa, Mode: 0o600, Size: int64(len(tresc)), Typeflag: tar.TypeReg,
		}
		if err := zapis.WriteHeader(naglowek); err != nil {
			return nil, fmt.Errorf("moduł Apps: nie można zapisać nagłówka %s: %w", nazwa, err)
		}
		if _, err := zapis.Write(tresc); err != nil {
			return nil, fmt.Errorf("moduł Apps: nie można zapisać wpisu %s: %w", nazwa, err)
		}
	}
	if err := zapis.Close(); err != nil {
		return nil, fmt.Errorf("moduł Apps: nie można domknąć archiwum tar: %w", err)
	}
	if err := kompresja.Close(); err != nil {
		return nil, fmt.Errorf("moduł Apps: nie można domknąć kompresji gzip: %w", err)
	}
	return bufor.Bytes(), nil
}

// zastrzezeniaPakietuApp jest walidatorem zgodności pakietu z kontraktem
// rozszerzenia, oddającym wykaz zastrzeżeń wagi błędu, ostrzeżenia i informacji.
func zastrzezeniaPakietuApp(wiersz dane.PakietApp) []shared.AppValidationIssue {
	zastrzezenia := []shared.AppValidationIssue{}
	if wiersz.Sciezka == nil || *wiersz.Sciezka == "" {
		zastrzezenia = append(zastrzezenia, shared.AppValidationIssue{
			Severity: shared.AppValidationSeverityError,
			Code:     kodPakietuBezArchiwumApp,
			Message:  "pakiet nie ma archiwum — zbuduj go komendą apps.package.build",
		})
	}
	if wiersz.Podpis == nil || *wiersz.Podpis == "" {
		zastrzezenia = append(zastrzezenia, shared.AppValidationIssue{
			Severity: shared.AppValidationSeverityWarning,
			Code:     kodPakietuBezPodpisuApp,
			Message:  "pakiet nie jest podpisany — pozycja trafi do rejestru jako niezweryfikowana",
		})
	}

	manifest, err := manifestPakietuApp(wiersz)
	if err != nil || manifest == nil {
		return append(zastrzezenia, shared.AppValidationIssue{
			Severity: shared.AppValidationSeverityError,
			Code:     kodPakietuBezManifestuApp,
			Message:  "pakiet nie ma czytelnego manifestu",
		})
	}

	if !wzorzecIdentyfikatoraPakietuApp.MatchString(manifest.Identifier) {
		zastrzezenia = append(zastrzezenia, shared.AppValidationIssue{
			Severity: shared.AppValidationSeverityError,
			Code:     kodPakietuTozsamoscApp,
			Message: "identyfikator " + strconv.Quote(manifest.Identifier) +
				" nie jest kodem pozycji katalogu (małe litery, cyfry, kropka, myślnik, podkreślenie)",
		})
	}
	if !wzorzecWersjiSemantycznejApp.MatchString(manifest.Version) {
		zastrzezenia = append(zastrzezenia, shared.AppValidationIssue{
			Severity: shared.AppValidationSeverityWarning,
			Code:     kodPakietuWersjaApp,
			Message: "wersja " + strconv.Quote(manifest.Version) +
				" nie jest wersją semantyczną (postać X.Y.Z)",
		})
	}
	if err := sprawdzRodzajRozszerzeniaPakietuApp(manifest.Kind); err != nil {
		zastrzezenia = append(zastrzezenia, shared.AppValidationIssue{
			Severity: shared.AppValidationSeverityError,
			Code:     kodPakietuRodzajApp,
			Message: "rodzaj " + strconv.Quote(string(manifest.Kind)) +
				" nie jest rodzajem rozszerzenia (mcp, plugin, api, skill)",
		})
	}
	if len(manifest.Tools) == 0 {
		zastrzezenia = append(zastrzezenia, shared.AppValidationIssue{
			Severity: shared.AppValidationSeverityInfo,
			Code:     kodPakietuBezNarzedziApp,
			Message:  "manifest nie deklaruje ani jednego narzędzia",
		})
	}
	for _, zaleznosc := range manifest.Dependencies {
		if strings.TrimSpace(zaleznosc) != "" {
			continue
		}
		zastrzezenia = append(zastrzezenia, shared.AppValidationIssue{
			Severity: shared.AppValidationSeverityWarning,
			Code:     kodPakietuZaleznoscPustaApp,
			Message:  "manifest niesie zależność bez nazwy",
		})
	}
	// Uprawnienie bez wskazania bytu jest uprawnieniem na wszystko; skaner o tym ostrzega.
	for _, uprawnienie := range manifest.Permissions {
		if uprawnienie.Target != nil && strings.TrimSpace(*uprawnienie.Target) != "" {
			continue
		}
		zastrzezenia = append(zastrzezenia, shared.AppValidationIssue{
			Severity: shared.AppValidationSeverityWarning,
			Code:     kodPakietuSzerokiZakresApp,
			Message: "uprawnienie " + string(uprawnienie.Scope) +
				" jest deklarowane bez wskazania bytu, którego dotyczy",
		})
	}
	return zastrzezenia
}

// pakietOknaApp odczytuje pakiet po kodzie i sprawdza, że należy do okna
// produktu, z którego przyszło żądanie.
func (a *adapterAplikacji) pakietOknaApp(ctx context.Context, okno, kod string) (dane.PakietApp, error) {
	wiersz, err := a.repozytorium.PakietApp(ctx, kod)
	if err != nil {
		return dane.PakietApp{}, bladNieznanegoBytuApp("pakiet", kod, err)
	}
	if wiersz.Okno != okno {
		return dane.PakietApp{}, bladWskazaniaAplikacji(
			"pakiet " + kod + " należy do okna " + wiersz.Okno +
				", a żądanie przyszło z okna " + okno)
	}
	return wiersz, nil
}

// manifestPakietuApp odczytuje manifest z wiersza; brak manifestu daje nil bez
// błędu, bo pakiet założony samym `apps.package.build` może go jeszcze nie mieć.
func manifestPakietuApp(wiersz dane.PakietApp) (*shared.AppPackageManifest, error) {
	if wiersz.Manifest == nil || *wiersz.Manifest == "" {
		return nil, nil
	}
	var manifest shared.AppPackageManifest
	if err := json.Unmarshal([]byte(*wiersz.Manifest), &manifest); err != nil {
		return nil, fmt.Errorf("moduł Apps: nieczytelny manifest pakietu %s: %w", wiersz.Kod, err)
	}
	return &manifest, nil
}

// pakietKontraktuApp przekłada wiersz pakietu z bazy danych na kształt
// pakietu zwracany kontraktem komunikacji, wraz z odczytanym manifestem.
func pakietKontraktuApp(wiersz dane.PakietApp) (shared.AppPackage, error) {
	manifest, err := manifestPakietuApp(wiersz)
	if err != nil {
		return shared.AppPackage{}, err
	}
	pakiet := shared.AppPackage{
		Id: wiersz.Kod, WindowId: wiersz.Okno, Manifest: manifest,
		ArtifactRef: wiersz.ArtefaktOdwolanie,
		Format:      shared.AppPackageFormat(wiersz.Format),
		SizeBytes:   wiersz.Rozmiar, PublishedExtensionId: wiersz.RozszerzenieKod,
		UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}
	if wiersz.Podpis != nil && *wiersz.Podpis != "" {
		var podpis shared.ExtensionSignature
		if err := json.Unmarshal([]byte(*wiersz.Podpis), &podpis); err != nil {
			return shared.AppPackage{}, fmt.Errorf(
				"moduł Apps: nieczytelny podpis pakietu %s: %w", wiersz.Kod, err)
		}
		pakiet.Signature = &podpis
	}
	return pakiet, nil
}

// kluczWydawcyApp wydaje klucz prywatny spod klucza jawnego. Odwołanie użyte po
// raz pierwszy zakłada klucz wydawcy nowy; sejf trzyma jego ziarno w postaci
// base64, bo przechowuje wyłącznie napisy.
func (a *adapterAplikacji) kluczWydawcyApp(ctx context.Context, odwolanie string) (ed25519.PrivateKey, error) {
	const bytWSejfie = "apps.package.signingKey:"
	if zapisane, jest := a.sejf.Odczytaj(ctx, bytWSejfie+odwolanie); jest && zapisane != "" {
		ziarno, err := base64.StdEncoding.DecodeString(zapisane)
		if err != nil {
			return nil, fmt.Errorf("moduł Apps: klucz wydawcy %q w sejfie jest nieczytelny: %w",
				odwolanie, err)
		}
		if len(ziarno) != ed25519.SeedSize {
			return nil, fmt.Errorf("moduł Apps: klucz wydawcy %q ma %d bajtów zamiast %d",
				odwolanie, len(ziarno), ed25519.SeedSize)
		}
		return ed25519.NewKeyFromSeed(ziarno), nil
	}

	publiczny, prywatny, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, fmt.Errorf("moduł Apps: nie można wytworzyć klucza wydawcy %q: %w", odwolanie, err)
	}
	_ = publiczny
	if _, err := a.sejf.Zapisz(ctx, bytWSejfie+odwolanie,
		base64.StdEncoding.EncodeToString(prywatny.Seed())); err != nil {
		return nil, fmt.Errorf("moduł Apps: nie można odłożyć klucza wydawcy %q w sejfie: %w",
			odwolanie, err)
	}
	return prywatny, nil
}

// sygnaturaZPodpisemApp dokłada do sygnatury jej postać bajtową. Kontrakt
// `ExtensionSignature` nie ma pola na sam podpis — niesie wyłącznie rozstrzygnięcia
// — więc bajty jadą do kolumny osobnym kształtem, a do klienta idzie sygnatura
// kontraktu bez nich.
func sygnaturaZPodpisemApp(sygnatura shared.ExtensionSignature,
	podpis []byte, publiczny ed25519.PublicKey) map[string]any {

	zapis := map[string]any{
		"signed":          sygnatura.Signed,
		"verified":        sygnatura.Verified,
		"trustLevel":      sygnatura.TrustLevel,
		"algorithm":       sygnatura.Algorithm,
		"checksumSha256":  sygnatura.ChecksumSha256,
		"publisher":       sygnatura.Publisher,
		"signatureBase64": base64.StdEncoding.EncodeToString(podpis),
		"publicKeyBase64": base64.StdEncoding.EncodeToString(publiczny),
	}
	if sygnatura.Detail != nil {
		zapis["detail"] = *sygnatura.Detail
	}
	return zapis
}

// sciezkaWMagazynieApp składa ścieżkę na dysku z odwołania magazynu treści,
// dopisując katalog danych, gdy odwołanie nie jest ścieżką bezwzględną.
func sciezkaWMagazynieApp(katalogDanych, odwolanie string) string {
	if filepath.IsAbs(odwolanie) {
		return odwolanie
	}
	return filepath.Join(katalogDanych, filepath.FromSlash(odwolanie))
}

// widocznoscPakietuApp rozstrzyga brak wskazania widoczności na widoczność
// domyślną, którą kontrakt nazywa widocznością w obrębie organizacji.
func widocznoscPakietuApp(wskazanie *shared.AppPackageVisibility) string {
	if wskazanie == nil {
		return shared.AppPackageVisibilityOrganization
	}
	return string(*wskazanie)
}

// wartoscLogicznaApp oddaje wskaźnik na kopię podanej wartości logicznej, do
// pól kontraktu zapisu, które przyjmują wskaźnik zamiast wartości wprost.
func wartoscLogicznaApp(wartosc bool) *bool {
	kopia := wartosc
	return &kopia
}

// sprawdzFormatPakietuApp dopuszcza wyłącznie formaty pakietu wymienione
// w kontrakcie komunikacji: `zip` albo `targz`.
func sprawdzFormatPakietuApp(format shared.AppPackageFormat) error {
	switch format {
	case shared.AppPackageFormatZip, shared.AppPackageFormatTargz:
		return nil
	}
	return bladWskazaniaAplikacji("nieznany format pakietu " + strconv.Quote(string(format)) +
		" — dopuszczalne: zip, targz")
}

// sprawdzWidocznoscPakietuApp dopuszcza wyłącznie widoczności wymienione
// w kontrakcie komunikacji: organizacyjną albo ograniczoną.
func sprawdzWidocznoscPakietuApp(widocznosc shared.AppPackageVisibility) error {
	switch widocznosc {
	case shared.AppPackageVisibilityOrganization, shared.AppPackageVisibilityRestricted:
		return nil
	}
	return bladWskazaniaAplikacji("nieznana widoczność pakietu " +
		strconv.Quote(string(widocznosc)) + " — dopuszczalne: organization, restricted")
}

// sprawdzRodzajRozszerzeniaPakietuApp dopuszcza wyłącznie rodzaje rozszerzenia
// wymienione w kontrakcie komunikacji: mcp, plugin, api albo skill.
func sprawdzRodzajRozszerzeniaPakietuApp(rodzaj shared.ExtensionKind) error {
	switch rodzaj {
	case shared.ExtensionKindMcp, shared.ExtensionKindPlugin,
		shared.ExtensionKindApi, shared.ExtensionKindSkill:
		return nil
	}
	return bladWskazaniaAplikacji("nieznany rodzaj rozszerzenia " +
		strconv.Quote(string(rodzaj)) + " — dopuszczalne: mcp, plugin, api, skill")
}
