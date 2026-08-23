package core

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"danacoconsole/shared"
)

// Skutek modułu Apps: czy za odpowiedzią stoi zapis, plik albo stojący serwer.
//
// Wzorzec szkody, którego pilnuje ten plik, wystąpił w tym produkcie: moduł
// Design meldował `status: ok` z wykazem zasobów, za którymi nie było ani
// jednego bajtu. Dlatego żaden sprawdzian tutaj nie kończy się na tym, że
// odpowiedź jest udana. Każdy schodzi niżej niż koperta:
//   - do bazy DRUGIM połączeniem, otwartym niezależnie od rdzenia, i liczy
//     wiersze albo czyta wartości kolumn,
//   - do pliku w magazynie treści, otwiera go i sprawdza, czy to naprawdę jest
//     archiwum, obraz albo dokument, za który się podaje,
//   - do sieci, wołając adres, który rdzeń wypuścił jako `previewUrl`.
//
// Sprawdziany nie pomijają się przy braku żadnego programu: cały moduł pracuje
// bibliotekami wkompilowanymi w rdzeń — `archive/zip`, `image/png`,
// `crypto/ed25519`, `net/http` — więc mierzą to, co ma stać u Operatora na
// cienkiej instalce.

// oknoSprawdzianuApps — okno, w którym stoją wszystkie sprawdziany tego pliku.
const oknoSprawdzianuApps = "apps.product-builder/sprawdzian"

// bazaSprawdzianuApps otwiera drugie połączenie z bazą rdzenia. To ono jest
// miarą: wartość odczytana tędy nie pochodzi od tego samego kodu, który ją
// zapisywał.
func bazaSprawdzianuApps(t *testing.T, katalog string) *sql.DB {
	t.Helper()

	baza, err := sql.Open("sqlite", filepath.Join(katalog, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy do pomiaru: %v", err)
	}
	t.Cleanup(func() { _ = baza.Close() })
	return baza
}

// wierszyApps liczy wiersze spełniające warunek — jedna miara na wszystkie tabele.
func wierszyApps(t *testing.T, baza *sql.DB, zapytanie string, argumenty ...any) int {
	t.Helper()

	var ile int
	if err := baza.QueryRow(zapytanie, argumenty...).Scan(&ile); err != nil {
		t.Fatalf("nie można policzyć wierszy (%s): %v", zapytanie, err)
	}
	return ile
}

// tekstZBazyApps odczytuje jedną wartość tekstową.
func tekstZBazyApps(t *testing.T, baza *sql.DB, zapytanie string, argumenty ...any) string {
	t.Helper()

	var wartosc string
	if err := baza.QueryRow(zapytanie, argumenty...).Scan(&wartosc); err != nil {
		t.Fatalf("nie można odczytać wartości (%s): %v", zapytanie, err)
	}
	return wartosc
}

// oknoModuluSprawdzianu zakłada sesję i okno modułu Apps, i oddaje jego
// identyfikator.
//
// Prawdziwe okno jest tu nieodzowne przy wdrożeniu: `apps.deployment.run` bierze
// treść z przestrzeni roboczej okna i odmawia oknu, którego rdzeń nie zna.
// Pozostałe komendy obszaru okna w rejestrze nie wymagają — pracują na wierszach
// znakowanych jego identyfikatorem — więc reszta sprawdzianów posługuje się
// nazwą stałą i mierzy to samo taniej.
func oknoModuluSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context) string {
	t.Helper()

	var sesja shared.SessionCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandSessionCreate,
		shared.SessionCreateRequest{Title: wskaznik("sprawdzian modułu Apps")}, &sesja)

	var okno shared.WindowCreateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandWindowCreate, shared.WindowCreateRequest{
		SessionId:      sesja.Session.Id,
		ModuleId:       "apps",
		ModelChannelId: "kanal-sprawdzianu",
		ExecutionEnv:   shared.ExecutionEnvLocal,
		PermissionMode: shared.PermissionModeAuto,
		WindowRole:     shared.WindowRoleStandalone,
	}, &okno)

	return okno.Window.Id
}

// zapiszPlikWarsztatuSprawdzianu kładzie plik w warstwie warsztatu okna.
func zapiszPlikWarsztatuSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	warstwa shared.AppWorkspaceLayer, sciezka, tresc string) {
	t.Helper()

	zapiszPlikWarsztatuOkna(t, zmontowany, zycie, oknoSprawdzianuApps, warstwa, sciezka, tresc)
}

// zapiszPlikWarsztatuOkna kładzie plik we wskazanym oknie — sprawdziany
// wdrożenia pracują na oknie założonym naprawdę, nie na nazwie stałej.
func zapiszPlikWarsztatuOkna(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	okno string, warstwa shared.AppWorkspaceLayer, sciezka, tresc string) {
	t.Helper()

	var wynik shared.AppsWorkspaceUpdateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsWorkspaceUpdate,
		shared.AppsWorkspaceUpdateRequest{
			WindowId: okno, Layer: warstwa, Path: sciezka, Content: tresc,
		}, &wynik)
}

// zdefiniujUkladSprawdzianu zapisuje architekturę o wskazanych komponentach.
func zdefiniujUkladSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	komponenty []shared.AppComponent) shared.AppArchitecture {
	t.Helper()

	var wynik shared.AppsArchitectureDefineResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsArchitectureDefine,
		shared.AppsArchitectureDefineRequest{
			WindowId: oknoSprawdzianuApps, Name: wskaznik("Portal sprawdzianu"),
			Components: komponenty,
		}, &wynik)
	return wynik.Architecture
}

// ── Product Builder ─────────────────────────────────────────────────────────

// TestZapisProduktuZostajeWBazie wykazuje skutek: po `apps.product.save`
// w tabeli leży wiersz o tej nazwie, a `apps.product.get` czyta ten sam.
func TestZapisProduktuZostajeWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	var zapis shared.AppsProductSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsProductSave,
		shared.AppsProductSaveRequest{
			WindowId: oknoSprawdzianuApps, Name: "Portal klienta",
			Description: wskaznik("produkt sprawdzianu"),
			Platforms: []shared.AppProductPlatform{
				shared.AppProductPlatformWeb, shared.AppProductPlatformMobile,
			},
		}, &zapis)

	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM produkt_apps WHERE okno = ? AND nazwa = 'Portal klienta'`,
		oknoSprawdzianuApps); ile != 1 {
		t.Fatalf("za odpowiedzią produktu nie ma wiersza w bazie: %d", ile)
	}
	platformy := tekstZBazyApps(t, baza,
		`SELECT platformy FROM produkt_apps WHERE okno = ?`, oknoSprawdzianuApps)
	if !strings.Contains(platformy, "web") || !strings.Contains(platformy, "mobile") {
		t.Fatalf("kolumna platform nie niesie obu platform: %q", platformy)
	}

	// Drugi zapis nie zakłada drugiego produktu — jeden produkt na okno.
	var drugi shared.AppsProductSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsProductSave,
		shared.AppsProductSaveRequest{WindowId: oknoSprawdzianuApps, Name: "Portal klienta v2"},
		&drugi)
	if ile := wierszyApps(t, baza, `SELECT COUNT(*) FROM produkt_apps WHERE okno = ?`,
		oknoSprawdzianuApps); ile != 1 {
		t.Fatalf("drugi zapis założył %d wierszy produktu zamiast nadpisać jeden", ile)
	}
	if drugi.Product.Id != zapis.Product.Id {
		t.Fatalf("tożsamość produktu zmieniła się przy zapisie: %q → %q",
			zapis.Product.Id, drugi.Product.Id)
	}

	var odczyt shared.AppsProductGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsProductGet,
		shared.AppsProductGetRequest{WindowId: oknoSprawdzianuApps}, &odczyt)
	if odczyt.Product == nil || odczyt.Product.Name != "Portal klienta v2" {
		t.Fatalf("odczyt nie oddał produktu zapisanego jako ostatni: %+v", odczyt.Product)
	}
}

// TestEtapyIKamienieMiloweZostajaWBazie wykazuje, że oś etapów i kamienie
// milowe mają pokrycie w wierszach, a usunięcie kamienia naprawdę je kasuje.
func TestEtapyIKamienieMiloweZostajaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	var etap shared.AppsStageSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsStageSave,
		shared.AppsStageSaveRequest{
			WindowId: oknoSprawdzianuApps, Name: wskaznik("Architektura"),
			Order: wskaznik(1), OwnerAgentId: wskaznik("ekspert-1"),
		}, &etap)

	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM etap_apps WHERE identyfikator_zewnetrzny = ? AND wykonawca = 'ekspert-1'`,
		etap.Stage.Id); ile != 1 {
		t.Fatalf("etapu nie ma w bazie albo nie ma wykonawcy: %d", ile)
	}

	// Pusty łańcuch zdejmuje przypisanie — kontrakt mówi to wprost.
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsStageSave,
		shared.AppsStageSaveRequest{
			WindowId: oknoSprawdzianuApps, StageId: wskaznik(etap.Stage.Id),
			Status:       wskaznik(shared.AppStageStatus(shared.AppStageStatusActive)),
			OwnerAgentId: wskaznik(""),
		}, &etap)
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM etap_apps WHERE identyfikator_zewnetrzny = ?
		   AND wykonawca IS NULL AND stan = 'active'`, etap.Stage.Id); ile != 1 {
		t.Fatal("pusty ownerAgentId nie zdjął przypisania albo stan nie wszedł do bazy")
	}

	var kamien shared.AppsMilestoneSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsMilestoneSave,
		shared.AppsMilestoneSaveRequest{
			WindowId: oknoSprawdzianuApps, Name: "Wydanie 1.0",
			DueAt: wskaznik(int64(1767225600000)), StageIds: []string{etap.Stage.Id},
		}, &kamien)

	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM kamien_milowy_etap_apps z
		   JOIN kamien_milowy_apps k ON k.id = z.kamien_id
		  WHERE k.identyfikator_zewnetrzny = ? AND z.etap_kod = ?`,
		kamien.Milestone.Id, etap.Stage.Id); ile != 1 {
		t.Fatal("związku kamienia milowego z etapem nie ma w tabeli złącznikowej")
	}

	var usuniecie shared.AppsMilestoneDeleteResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsMilestoneDelete,
		shared.AppsMilestoneDeleteRequest{
			WindowId: oknoSprawdzianuApps, MilestoneId: kamien.Milestone.Id,
		}, &usuniecie)
	if !usuniecie.Deleted {
		t.Fatal("usunięcie kamienia milowego zameldowało brak skutku")
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM kamien_milowy_apps WHERE identyfikator_zewnetrzny = ?`,
		kamien.Milestone.Id); ile != 0 {
		t.Fatal("kamień milowy zameldowany jako usunięty leży dalej w bazie")
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM kamien_milowy_etap_apps`); ile != 0 {
		t.Fatal("kaskada nie zdjęła związku kamienia z etapem")
	}
}

// TestOsCzasuSkladaSieZPieciuZrodel wykazuje, że oś czasu nie jest osobną
// tabelą, tylko widokiem pracy — i że pokazuje zdarzenia, które naprawdę zaszły.
func TestOsCzasuSkladaSieZPieciuZrodel(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	zdefiniujUkladSprawdzianu(t, zmontowany, zycie, []shared.AppComponent{
		{Id: "front", Name: "Interfejs", Kind: shared.AppComponentKindFrontend},
	})
	zapiszPlikWarsztatuSprawdzianu(t, zmontowany, zycie,
		shared.AppWorkspaceLayerFrontend, "index.html", "<h1>Portal</h1>")
	var etap shared.AppsStageSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsStageSave,
		shared.AppsStageSaveRequest{WindowId: oknoSprawdzianuApps, Name: wskaznik("Frontend")},
		&etap)

	var os shared.AppsTimelineListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsTimelineList,
		shared.AppsTimelineListRequest{WindowId: oknoSprawdzianuApps}, &os)

	rodzaje := map[shared.AppTimelineKind]int{}
	for _, wpis := range os.Entries {
		rodzaje[wpis.Kind]++
	}
	for _, oczekiwany := range []shared.AppTimelineKind{
		shared.AppTimelineKindArchitecture,
		shared.AppTimelineKindWorkspace,
		shared.AppTimelineKindStage,
	} {
		if rodzaje[oczekiwany] == 0 {
			t.Fatalf("oś czasu nie widzi zdarzeń rodzaju %s: %+v", oczekiwany, rodzaje)
		}
	}
	if os.Total != len(os.Entries) {
		t.Fatalf("total (%d) rozjeżdża się z liczbą oddanych zdarzeń (%d)",
			os.Total, len(os.Entries))
	}
}

// TestPowiazaniaProduktuLiczaRzeczywistyStan wykazuje, że powiązanie czynne ma
// pokrycie w danych, a nie w deklaracji: przed zapisem pliku warsztatu
// powiązanie z Developerem jest nieczynne, po zapisie — czynne.
func TestPowiazaniaProduktuLiczaRzeczywistyStan(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var przed shared.AppsProductLinkListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsProductLinkList,
		shared.AppsProductLinkListRequest{WindowId: oknoSprawdzianuApps}, &przed)
	if czyPowiazanieCzynneApps(przed.Links, "developer") {
		t.Fatal("powiązanie z Developerem jest czynne, choć w oknie nie ma ani jednego pliku")
	}

	zapiszPlikWarsztatuSprawdzianu(t, zmontowany, zycie,
		shared.AppWorkspaceLayerBackend, "serwer.go", "package main")

	var po shared.AppsProductLinkListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsProductLinkList,
		shared.AppsProductLinkListRequest{WindowId: oknoSprawdzianuApps}, &po)
	if !czyPowiazanieCzynneApps(po.Links, "developer") {
		t.Fatal("powiązanie z Developerem nie zauważyło pliku warsztatu")
	}
}

// czyPowiazanieCzynneApps odczytuje stan jednego powiązania z wykazu.
func czyPowiazanieCzynneApps(powiazania []shared.AppProductLink, kod string) bool {
	for _, powiazanie := range powiazania {
		if powiazanie.ModuleCode == kod {
			return powiazanie.Enabled
		}
	}
	return false
}

// ── Architecture Designer ───────────────────────────────────────────────────

// TestWalidacjaUkladuWykrywaOsamotnienieINieznanaZaleznosc wykazuje, że
// walidator liczy zastrzeżenia z układu i odkłada je w bazie.
func TestWalidacjaUkladuWykrywaOsamotnienieINieznanaZaleznosc(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	zdefiniujUkladSprawdzianu(t, zmontowany, zycie, []shared.AppComponent{
		{Id: "front", Name: "Interfejs", Kind: shared.AppComponentKindFrontend},
		{Id: "back", Name: "Usługi", Kind: shared.AppComponentKindBackend,
			DependsOn: []string{"widmo"}},
	})

	var wynik shared.AppsArchitectureValidateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsArchitectureValidate,
		shared.AppsArchitectureValidateRequest{WindowId: oknoSprawdzianuApps}, &wynik)

	kody := map[string]int{}
	for _, zastrzezenie := range wynik.Issues {
		kody[zastrzezenie.Code]++
	}
	if kody[kodZastrzezeniaNieznanyApp] == 0 {
		t.Fatalf("walidator nie zauważył zależności od komponentu spoza układu: %+v", wynik.Issues)
	}
	if kody[kodZastrzezeniaOsamotnionyApp] == 0 {
		t.Fatalf("walidator nie zauważył komponentu bez zależności: %+v", wynik.Issues)
	}

	zapisane := tekstZBazyApps(t, baza,
		`SELECT IFNULL(zastrzezenia_walidacji,'') FROM architektura_apps WHERE okno = ?`,
		oknoSprawdzianuApps)
	if !strings.Contains(zapisane, kodZastrzezeniaNieznanyApp) {
		t.Fatalf("zastrzeżenia nie wróciły do kolumny bazy: %q", zapisane)
	}
}

// TestWalidacjaWykrywaCyklZaleznosci wykazuje, że walidator naprawdę przechodzi
// graf, a nie sprawdza samych końców krawędzi.
func TestWalidacjaWykrywaCyklZaleznosci(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	zdefiniujUkladSprawdzianu(t, zmontowany, zycie, []shared.AppComponent{
		{Id: "a", Name: "A", Kind: shared.AppComponentKindService, DependsOn: []string{"c"}},
		{Id: "b", Name: "B", Kind: shared.AppComponentKindService, DependsOn: []string{"a"}},
		{Id: "c", Name: "C", Kind: shared.AppComponentKindService, DependsOn: []string{"b"}},
	})

	var wynik shared.AppsArchitectureValidateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsArchitectureValidate,
		shared.AppsArchitectureValidateRequest{WindowId: oknoSprawdzianuApps}, &wynik)

	for _, zastrzezenie := range wynik.Issues {
		if zastrzezenie.Code == kodZastrzezeniaCyklApp {
			return
		}
	}
	t.Fatalf("walidator nie zauważył cyklu a→b→c→a: %+v", wynik.Issues)
}

// TestHistoriaWersjiRosnieZKazdymZapisem wykazuje, że wersje nie są liczbą
// w kolumnie, tylko wierszami, które da się wyliczyć.
func TestHistoriaWersjiRosnieZKazdymZapisem(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	pierwsza := zdefiniujUkladSprawdzianu(t, zmontowany, zycie, []shared.AppComponent{
		{Id: "front", Name: "Interfejs", Kind: shared.AppComponentKindFrontend},
	})
	var drugi shared.AppsArchitectureDefineResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsArchitectureDefine,
		shared.AppsArchitectureDefineRequest{
			WindowId: oknoSprawdzianuApps, ArchitectureId: wskaznik(pierwsza.Id),
			Components: []shared.AppComponent{
				{Id: "front", Name: "Interfejs", Kind: shared.AppComponentKindFrontend},
				{Id: "back", Name: "Usługi", Kind: shared.AppComponentKindBackend},
			},
		}, &drugi)

	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM wersja_architektury_apps w
		   JOIN architektura_apps a ON a.id = w.architektura_id
		  WHERE a.identyfikator_zewnetrzny = ?`, pierwsza.Id); ile != 2 {
		t.Fatalf("historia wersji ma %d wierszy zamiast dwóch", ile)
	}

	var wersje shared.AppsArchitectureVersionListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsArchitectureVersionList,
		shared.AppsArchitectureVersionListRequest{WindowId: oknoSprawdzianuApps}, &wersje)
	if wersje.Total != 2 {
		t.Fatalf("odczyt historii oddał %d wersji zamiast dwóch", wersje.Total)
	}
	if wersje.Versions[0].ComponentCount != 2 {
		t.Fatalf("najnowsza wersja liczy %d komponentów zamiast dwóch",
			wersje.Versions[0].ComponentCount)
	}
}

// TestAdnotacjaPrzypinaSieDoJednegoMiejsca wykazuje zapis notatki i odmowę przy
// przypięciu sprzecznym.
func TestAdnotacjaPrzypinaSieDoJednegoMiejsca(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	zdefiniujUkladSprawdzianu(t, zmontowany, zycie, []shared.AppComponent{
		{Id: "back", Name: "Usługi", Kind: shared.AppComponentKindBackend},
	})

	var zapis shared.AppsArchitectureAnnotationSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsArchitectureAnnotationSave,
		shared.AppsArchitectureAnnotationSaveRequest{
			WindowId: oknoSprawdzianuApps, ComponentId: wskaznik("back"),
			Text: "kolejka zamówień idzie osobnym kanałem",
		}, &zapis)

	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM adnotacja_architektury_apps
		  WHERE identyfikator_zewnetrzny = ? AND komponent_kod = 'back'`,
		zapis.Annotation.Id); ile != 1 {
		t.Fatal("adnotacji nie ma w bazie albo nie jest przypięta do komponentu")
	}

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandAppsArchitectureAnnotationSave,
		shared.AppsArchitectureAnnotationSaveRequest{
			WindowId: oknoSprawdzianuApps, ComponentId: wskaznik("back"),
			DependencyFrom: wskaznik("back"), DependencyTo: wskaznik("front"),
			Text: "przypięcie sprzeczne",
		})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("odmowa sprzecznego przypięcia ma kod %q zamiast validation_failed", odmowa.Code)
	}
}

// TestEksportUkladuDajePlikWKazdymFormacie wykazuje najtwardszy skutek tego
// obszaru: za odwołaniem leży plik, który naprawdę jest tym, za co się podaje.
func TestEksportUkladuDajePlikWKazdymFormacie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)

	zdefiniujUkladSprawdzianu(t, zmontowany, zycie, []shared.AppComponent{
		{Id: "front", Name: "Interfejs", Kind: shared.AppComponentKindFrontend},
		{Id: "back", Name: "Usługi", Kind: shared.AppComponentKindBackend,
			DependsOn: []string{"front"}},
	})

	for _, przypadek := range []struct {
		format  shared.AppExportFormat
		sprawdz func(t *testing.T, bajty []byte)
	}{
		{shared.AppExportFormatMermaid, func(t *testing.T, bajty []byte) {
			tresc := string(bajty)
			if !strings.HasPrefix(tresc, "graph TD") || !strings.Contains(tresc, "-->") {
				t.Fatalf("plik mermaid nie jest diagramem: %q", tresc)
			}
		}},
		{shared.AppExportFormatMarkdown, func(t *testing.T, bajty []byte) {
			tresc := string(bajty)
			if !strings.Contains(tresc, "## Komponenty") || !strings.Contains(tresc, "Interfejs") {
				t.Fatalf("dokument markdown nie opisuje komponentów: %q", tresc)
			}
		}},
		{shared.AppExportFormatSvg, func(t *testing.T, bajty []byte) {
			tresc := string(bajty)
			if !strings.Contains(tresc, "<svg") || !strings.Contains(tresc, "<rect") {
				t.Fatalf("plik svg nie jest rysunkiem: %q", tresc)
			}
		}},
		{shared.AppExportFormatPng, func(t *testing.T, bajty []byte) {
			// Najtwardsza miara: dekoder PNG czyta nagłówek i piksele. Napis
			// udający obraz przeszedłby sprawdzenie rozmiaru, ale nie to.
			obraz, err := png.Decode(bytes.NewReader(bajty))
			if err != nil {
				t.Fatalf("plik png nie jest obrazem: %v", err)
			}
			if obraz.Bounds().Dx() <= 0 || obraz.Bounds().Dy() <= 0 {
				t.Fatalf("obraz ma zerowe wymiary: %v", obraz.Bounds())
			}
		}},
	} {
		t.Run(string(przypadek.format), func(t *testing.T) {
			var wynik shared.AppsArchitectureExportResponse
			wykonajUdana(t, zmontowany, zycie, shared.CommandAppsArchitectureExport,
				shared.AppsArchitectureExportRequest{
					WindowId: oknoSprawdzianuApps, Format: przypadek.format,
				}, &wynik)

			bajty := bajtyPodOdwolaniem(t, katalog, wynik.ArtifactRef)
			if int64(len(bajty)) != wynik.SizeBytes {
				t.Fatalf("rdzeń zameldował %d bajtów, a plik ma %d",
					wynik.SizeBytes, len(bajty))
			}
			przypadek.sprawdz(t, bajty)
		})
	}
}

// ── Warsztaty: podgląd, trasy, motyw, punkty końcowe, schemat ───────────────

// TestPodgladOddajeTrescPodWlasnymAdresem wykazuje skutek najtrudniejszy do
// udania: pod adresem z odpowiedzi naprawdę stoi serwer i oddaje plik warsztatu.
func TestPodgladOddajeTrescPodWlasnymAdresem(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	zapiszPlikWarsztatuSprawdzianu(t, zmontowany, zycie,
		shared.AppWorkspaceLayerFrontend, "index.html",
		"<!doctype html><title>Portal</title><p>zamówienia</p>")

	var start shared.AppsPreviewStartResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPreviewStart,
		shared.AppsPreviewStartRequest{WindowId: oknoSprawdzianuApps}, &start)
	if start.Status != shared.AppPreviewStatusRunning || start.PreviewUrl == "" {
		t.Fatalf("podgląd zameldował stan %q pod adresem %q", start.Status, start.PreviewUrl)
	}

	odpowiedz, err := http.Get(start.PreviewUrl)
	if err != nil {
		t.Fatalf("pod adresem %s nikt nie odpowiada: %v", start.PreviewUrl, err)
	}
	tresc, err := io.ReadAll(odpowiedz.Body)
	_ = odpowiedz.Body.Close()
	if err != nil {
		t.Fatalf("nie można odczytać odpowiedzi podglądu: %v", err)
	}
	if !strings.Contains(string(tresc), "zamówienia") {
		t.Fatalf("podgląd oddał treść, której nie ma w warsztacie: %q", tresc)
	}

	if stan := tekstZBazyApps(t, baza, `SELECT stan FROM podglad_apps WHERE okno = ?`,
		oknoSprawdzianuApps); stan != string(shared.AppPreviewStatusRunning) {
		t.Fatalf("stan podglądu w bazie to %q zamiast running", stan)
	}

	var stop shared.AppsPreviewStopResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPreviewStop,
		shared.AppsPreviewStopRequest{WindowId: oknoSprawdzianuApps}, &stop)
	if !stop.Stopped {
		t.Fatal("zatrzymanie podglądu zameldowało brak skutku")
	}
	// Zatrzymanie ma zamknąć nasłuch, nie tylko przestawić kolumnę.
	klient := &http.Client{Timeout: 2 * time.Second}
	if _, err := klient.Get(start.PreviewUrl); err == nil {
		t.Fatal("po zatrzymaniu podglądu adres dalej odpowiada — nasłuch nie został zamknięty")
	}
	if stan := tekstZBazyApps(t, baza, `SELECT stan FROM podglad_apps WHERE okno = ?`,
		oknoSprawdzianuApps); stan != string(shared.AppPreviewStatusStopped) {
		t.Fatalf("stan podglądu w bazie to %q zamiast stopped", stan)
	}

	// Dziennik ma pokrycie w obu czynnościach.
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM wiersz_dziennika_apps WHERE okno = ? AND tresc LIKE 'podgląd%'`,
		oknoSprawdzianuApps); ile < 2 {
		t.Fatalf("dziennik zapisał %d wierszy podglądu zamiast dwóch", ile)
	}
}

// TestTrasyCzytaneSaZPlikowWarstwyInterfejsu wykazuje, że mapa routingu jest
// odczytem pracy, a nie osobnym rejestrem.
func TestTrasyCzytaneSaZPlikowWarstwyInterfejsu(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	var puste shared.AppsRouteListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsRouteList,
		shared.AppsRouteListRequest{WindowId: oknoSprawdzianuApps}, &puste)
	if puste.Total != 0 {
		t.Fatalf("okno bez plików ma %d tras", puste.Total)
	}

	zapiszPlikWarsztatuSprawdzianu(t, zmontowany, zycie,
		shared.AppWorkspaceLayerFrontend, "trasy.ts",
		"export const trasy = [\n"+
			"  { path: '/zamowienia', component: 'ListaZamowien' },\n"+
			"  { path: '/klienci', component: 'ListaKlientow' },\n"+
			"];\n")
	zapiszPlikWarsztatuSprawdzianu(t, zmontowany, zycie,
		shared.AppWorkspaceLayerFrontend, "app.tsx",
		`<Route path="/faktury" element={<Faktury/>} />`)

	var trasy shared.AppsRouteListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsRouteList,
		shared.AppsRouteListRequest{WindowId: oknoSprawdzianuApps}, &trasy)
	if trasy.Total != 3 {
		t.Fatalf("mapa routingu ma %d tras zamiast trzech: %+v", trasy.Total, trasy.Routes)
	}
	sciezki := map[string]string{}
	for _, trasa := range trasy.Routes {
		widok := ""
		if trasa.ViewName != nil {
			widok = *trasa.ViewName
		}
		sciezki[trasa.Path] = widok
	}
	if sciezki["/zamowienia"] != "ListaZamowien" {
		t.Fatalf("trasa /zamowienia nie ma widoku: %q", sciezki["/zamowienia"])
	}
	if _, jest := sciezki["/faktury"]; !jest {
		t.Fatalf("trasa ze znacznika Route nie weszła do mapy: %+v", sciezki)
	}
}

// TestMotywZostajeWBazieIWracaTenSam wykazuje przechowanie surowego JSON-a.
func TestMotywZostajeWBazieIWracaTenSam(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	motyw := jsonSurowy(t, map[string]any{"kolorGlowny": "#c8a24a", "krok": 8})
	var zapis shared.AppsThemeSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsThemeSet,
		shared.AppsThemeSetRequest{WindowId: oknoSprawdzianuApps, Theme: motyw}, &zapis)

	zapisany := tekstZBazyApps(t, baza, `SELECT tresc FROM motyw_apps WHERE okno = ?`,
		oknoSprawdzianuApps)
	if !strings.Contains(zapisany, "#c8a24a") {
		t.Fatalf("motywu nie ma w bazie: %q", zapisany)
	}

	var odczyt shared.AppsThemeGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsThemeGet,
		shared.AppsThemeGetRequest{WindowId: oknoSprawdzianuApps}, &odczyt)
	var wczytany map[string]any
	if err := json.Unmarshal(odczyt.Theme, &wczytany); err != nil {
		t.Fatalf("odczytany motyw nie jest JSON-em: %v", err)
	}
	if wczytany["kolorGlowny"] != "#c8a24a" {
		t.Fatalf("odczyt oddał inny motyw niż zapisano: %+v", wczytany)
	}

	// Motyw pusty jest odmową, nie zapisem pustki: kontrakt ma pole jako
	// wymagane, a zapis `{}` skasowałby Operatorowi motyw bez jego żądania.
	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandAppsThemeSet,
		shared.AppsThemeSetRequest{WindowId: oknoSprawdzianuApps})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("motyw pusty wrócił kodem %q", odmowa.Code)
	}
}

// TestPunktyKoncoweCzytaneSaZKontraktuApi wykazuje, że eksplorator punktów
// końcowych opisuje architekturę, a stan punktu wynika z pracy w backendzie.
func TestPunktyKoncoweCzytaneSaZKontraktuApi(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	zdefiniujUkladSprawdzianu(t, zmontowany, zycie, []shared.AppComponent{
		{Id: "back", Name: "Usługi", Kind: shared.AppComponentKindBackend,
			ApiContract: wskaznik("GET /zamowienia — wykaz zamówień\nPOST /klienci — nowy klient")},
	})
	zapiszPlikWarsztatuSprawdzianu(t, zmontowany, zycie,
		shared.AppWorkspaceLayerBackend, "trasy.go",
		`mux.HandleFunc("/zamowienia", wykazZamowien)`)

	var punkty shared.AppsEndpointListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsEndpointList,
		shared.AppsEndpointListRequest{WindowId: oknoSprawdzianuApps}, &punkty)
	if punkty.Total != 2 {
		t.Fatalf("eksplorator ma %d punktów zamiast dwóch: %+v", punkty.Total, punkty.Endpoints)
	}

	stany := map[string]shared.AppEndpointStatus{}
	for _, punkt := range punkty.Endpoints {
		if punkt.Status != nil {
			stany[punkt.Path] = *punkt.Status
		}
	}
	if stany["/zamowienia"] != shared.AppEndpointStatusImplemented {
		t.Fatalf("punkt obecny w kodzie backendu ma stan %q zamiast implemented",
			stany["/zamowienia"])
	}
	if stany["/klienci"] != shared.AppEndpointStatusDraft {
		t.Fatalf("punkt nieobecny w kodzie ma stan %q zamiast draft", stany["/klienci"])
	}
}

// TestZapytanieProbneMierzyPrawdziwaOdpowiedz wykazuje, że proba idzie po sieci:
// sprawdzian podnosi własny serwer i oczekuje jego kodu odpowiedzi.
func TestZapytanieProbneMierzyPrawdziwaOdpowiedz(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	serwer := httptest.NewServer(http.HandlerFunc(
		func(odpowiedz http.ResponseWriter, zadanie *http.Request) {
			if zadanie.URL.Path != "/zamowienia" {
				odpowiedz.WriteHeader(http.StatusNotFound)
				return
			}
			odpowiedz.Header().Set("X-Sprawdzian", "tak")
			odpowiedz.WriteHeader(http.StatusCreated)
			_, _ = odpowiedz.Write([]byte(`{"zamowien":3}`))
		}))
	t.Cleanup(serwer.Close)

	var domena shared.AppsDeploymentDomainSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsDeploymentDomainSet,
		shared.AppsDeploymentDomainSetRequest{
			WindowId: oknoSprawdzianuApps, Environment: shared.AppDeployEnvironmentDev,
			Domain: strings.TrimPrefix(serwer.URL, "http://"),
		}, &domena)
	if !domena.Verified {
		t.Fatalf("domena %q na pętli zwrotnej nie rozwiązała się", domena.Domain)
	}

	var proba shared.AppsEndpointProbeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsEndpointProbe,
		shared.AppsEndpointProbeRequest{
			WindowId: oknoSprawdzianuApps, Method: shared.AppEndpointMethodGet,
			Path: "/zamowienia", Environment: wskaznik(
				shared.AppDeployEnvironment(shared.AppDeployEnvironmentDev)),
		}, &proba)

	if proba.Result.StatusCode != http.StatusCreated {
		t.Fatalf("zapytanie próbne oddało kod %d zamiast 201: %+v",
			proba.Result.StatusCode, proba.Result)
	}
	if proba.Result.Body == nil || !strings.Contains(*proba.Result.Body, `"zamowien":3`) {
		t.Fatalf("zapytanie próbne nie przyniosło treści odpowiedzi: %+v", proba.Result)
	}
	var naglowki map[string]string
	if err := json.Unmarshal(proba.Result.Headers, &naglowki); err != nil {
		t.Fatalf("nagłówki odpowiedzi nie są obiektem: %v", err)
	}
	if naglowki["X-Sprawdzian"] != "tak" {
		t.Fatalf("nagłówki odpowiedzi nie przeszły: %+v", naglowki)
	}

	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM wiersz_dziennika_apps WHERE okno = ? AND tresc LIKE 'zapytanie próbne%'`,
		oknoSprawdzianuApps); ile == 0 {
		t.Fatal("zapytanie próbne nie zostawiło śladu w dzienniku")
	}
}

// TestSchematCzytanyJestZPolecenCreateTable wykazuje, że podgląd schematu
// opisuje kod backendu, a nie osobną tabelę.
func TestSchematCzytanyJestZPolecenCreateTable(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	zapiszPlikWarsztatuSprawdzianu(t, zmontowany, zycie,
		shared.AppWorkspaceLayerBackend, "schemat.sql",
		"CREATE TABLE klient (\n"+
			"  id INTEGER PRIMARY KEY,\n"+
			"  nazwa TEXT NOT NULL,\n"+
			"  opis TEXT\n"+
			");\n"+
			"CREATE TABLE zamowienie (\n"+
			"  id INTEGER PRIMARY KEY,\n"+
			"  klient_id INTEGER NOT NULL REFERENCES klient(id),\n"+
			"  kwota DECIMAL(10,2) NOT NULL\n"+
			");\n")

	var schemat shared.AppsSchemaGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsSchemaGet,
		shared.AppsSchemaGetRequest{WindowId: oknoSprawdzianuApps}, &schemat)
	if schemat.Schema == nil {
		t.Fatal("schemat nie wrócił, choć backend niesie dwa polecenia CREATE TABLE")
	}
	if len(schemat.Schema.Tables) != 2 {
		t.Fatalf("schemat ma %d tabel zamiast dwóch: %+v",
			len(schemat.Schema.Tables), schemat.Schema.Tables)
	}

	tabele := map[string]shared.AppSchemaTable{}
	for _, tabela := range schemat.Schema.Tables {
		tabele[tabela.Name] = tabela
	}
	klient, jest := tabele["klient"]
	if !jest || len(klient.Columns) != 3 {
		t.Fatalf("tabela klient ma %d kolumn zamiast trzech: %+v", len(klient.Columns), klient)
	}
	if klient.Columns[0].Name != "id" || klient.Columns[0].PrimaryKey == nil ||
		!*klient.Columns[0].PrimaryKey {
		t.Fatalf("kolumna klucza pierwotnego nie została rozpoznana: %+v", klient.Columns[0])
	}
	if klient.Columns[1].Nullable == nil || *klient.Columns[1].Nullable {
		t.Fatalf("kolumna NOT NULL została opisana jako dopuszczająca brak: %+v", klient.Columns[1])
	}
	zamowienie := tabele["zamowienie"]
	if len(zamowienie.References) != 1 || zamowienie.References[0] != "klient" {
		t.Fatalf("więz obcy nie został odczytany: %+v", zamowienie.References)
	}
	// Przecinek w DECIMAL(10,2) nie ma prawa rozbić definicji kolumny.
	if len(zamowienie.Columns) != 3 {
		t.Fatalf("tabela zamowienie ma %d kolumn zamiast trzech: %+v",
			len(zamowienie.Columns), zamowienie.Columns)
	}
}

// ── Deployment Panel ────────────────────────────────────────────────────────

// TestSrodowiskaZmienneISkalowanieZostajaWBazie wykazuje trwałość nastaw panelu.
func TestSrodowiskaZmienneISkalowanieZostajaWBazie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	var srodowiska shared.AppsEnvironmentListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsEnvironmentList,
		shared.AppsEnvironmentListRequest{WindowId: oknoSprawdzianuApps}, &srodowiska)
	if srodowiska.Total != 3 {
		t.Fatalf("wykaz środowisk ma %d pozycji zamiast trzech", srodowiska.Total)
	}
	if ile := wierszyApps(t, baza, `SELECT COUNT(*) FROM srodowisko_apps WHERE okno = ?`,
		oknoSprawdzianuApps); ile != 3 {
		t.Fatalf("w bazie stoi %d środowisk zamiast trzech", ile)
	}
	// Drugi wykaz nie zakłada wierszy po raz drugi.
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsEnvironmentList,
		shared.AppsEnvironmentListRequest{WindowId: oknoSprawdzianuApps}, &srodowiska)
	if ile := wierszyApps(t, baza, `SELECT COUNT(*) FROM srodowisko_apps WHERE okno = ?`,
		oknoSprawdzianuApps); ile != 3 {
		t.Fatalf("drugi wykaz podniósł liczbę środowisk do %d", ile)
	}

	var zmienna shared.AppsEnvironmentVariableSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsEnvironmentVariableSet,
		shared.AppsEnvironmentVariableSetRequest{
			WindowId: oknoSprawdzianuApps, Environment: shared.AppDeployEnvironmentProduction,
			Name: "BAZA_HASLO", SecretRef: wskaznik("sejf:baza/haslo"),
		}, &zmienna)
	if zmienna.Variable.Value != nil {
		t.Fatal("zmienna sekretna wróciła z wartością jawną")
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM zmienna_srodowiska_apps
		  WHERE okno = ? AND nazwa = 'BAZA_HASLO' AND wartosc IS NULL
		    AND odwolanie_sekretu = 'sejf:baza/haslo'`, oknoSprawdzianuApps); ile != 1 {
		t.Fatal("zmiennej sekretnej nie ma w bazie albo trzyma wartość jawną")
	}

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandAppsEnvironmentVariableSet,
		shared.AppsEnvironmentVariableSetRequest{
			WindowId: oknoSprawdzianuApps, Environment: shared.AppDeployEnvironmentProduction,
			Name: "OBIE", Value: wskaznik("jawna"), SecretRef: wskaznik("sejf:x"),
		})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("zmienna z dwiema wartościami wróciła kodem %q", odmowa.Code)
	}

	var skalowanie shared.AppsDeploymentScaleSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsDeploymentScaleSet,
		shared.AppsDeploymentScaleSetRequest{
			WindowId: oknoSprawdzianuApps, Environment: shared.AppDeployEnvironmentProduction,
			MinInstances: wskaznik(2), MaxInstances: wskaznik(8),
		}, &skalowanie)
	if !skalowanie.Applied || skalowanie.EffectiveInstances == nil ||
		*skalowanie.EffectiveInstances != 2 {
		t.Fatalf("nastawa skalowania oddała %+v", skalowanie)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM skalowanie_apps
		  WHERE okno = ? AND srodowisko = 'production' AND min_instancji = 2 AND maks_instancji = 8`,
		oknoSprawdzianuApps); ile != 1 {
		t.Fatal("nastawy skalowania nie ma w bazie")
	}

	odmowaGranic := wykonajOdmowna(t, zmontowany, zycie, shared.CommandAppsDeploymentScaleSet,
		shared.AppsDeploymentScaleSetRequest{
			WindowId: oknoSprawdzianuApps, Environment: shared.AppDeployEnvironmentDev,
			MinInstances: wskaznik(9), MaxInstances: wskaznik(2),
		})
	if odmowaGranic.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("odwrócone granice skalowania wróciły kodem %q", odmowaGranic.Code)
	}
}

// TestKondycjaMierzyProduktIZapisujeSprawdzenie wykazuje, że dostępność bierze
// się z zapytania, a udział — z wierszy dziennika kondycji.
func TestKondycjaMierzyProduktIZapisujeSprawdzenie(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	serwer := httptest.NewServer(http.HandlerFunc(
		func(odpowiedz http.ResponseWriter, _ *http.Request) {
			odpowiedz.WriteHeader(http.StatusOK)
		}))
	t.Cleanup(serwer.Close)

	var domena shared.AppsDeploymentDomainSetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsDeploymentDomainSet,
		shared.AppsDeploymentDomainSetRequest{
			WindowId: oknoSprawdzianuApps, Environment: shared.AppDeployEnvironmentStaging,
			Domain: strings.TrimPrefix(serwer.URL, "http://"),
		}, &domena)

	var kondycja shared.AppsDeploymentHealthGetResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsDeploymentHealthGet,
		shared.AppsDeploymentHealthGetRequest{
			WindowId: oknoSprawdzianuApps,
			Environment: wskaznik(
				shared.AppDeployEnvironment(shared.AppDeployEnvironmentStaging)),
		}, &kondycja)

	if !kondycja.Health.Available {
		t.Fatalf("produkt odpowiadający kodem 200 uznano za niedostępny: %+v", kondycja.Health)
	}
	if kondycja.Health.AvailabilityPercent == nil ||
		*kondycja.Health.AvailabilityPercent != "100.00" {
		t.Fatalf("udział dostępności po jednym udanym sprawdzeniu: %+v",
			kondycja.Health.AvailabilityPercent)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM kondycja_wdrozenia_apps WHERE okno = ? AND dostepna = 1`,
		oknoSprawdzianuApps); ile != 1 {
		t.Fatalf("dziennik kondycji ma %d wierszy zamiast jednego", ile)
	}

	// Serwer padnięty: drugie sprawdzenie ma zejść do 50% i zapisać wiersz.
	serwer.Close()
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsDeploymentHealthGet,
		shared.AppsDeploymentHealthGetRequest{
			WindowId: oknoSprawdzianuApps,
			Environment: wskaznik(
				shared.AppDeployEnvironment(shared.AppDeployEnvironmentStaging)),
		}, &kondycja)
	if kondycja.Health.Available {
		t.Fatal("produkt, który nie odpowiada, uznano za dostępny")
	}
	if kondycja.Health.AvailabilityPercent == nil ||
		*kondycja.Health.AvailabilityPercent != "50.00" {
		t.Fatalf("udział po jednym udanym i jednym nieudanym sprawdzeniu: %+v",
			kondycja.Health.AvailabilityPercent)
	}
}

// ── Wdrożenie, artefakt, dziennik ───────────────────────────────────────────

// przeprowadzWdrozenieSprawdzianu uruchamia wdrożenie i czeka, aż silnik
// wykonania dobiegnie. Czekanie idzie po bazie, nie po zegarze: przebieg biegnie
// w gorutynie rdzenia, a uśpienie na stałą liczbę milisekund byłoby sprawdzianem
// szybkości maszyny, nie skutku.
func przeprowadzWdrozenieSprawdzianu(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	baza *sql.DB, okno string) string {
	t.Helper()

	var wynik shared.AppsDeploymentRunResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsDeploymentRun,
		shared.AppsDeploymentRunRequest{
			WindowId: okno, Environment: shared.AppDeployEnvironmentStaging,
			Version: wskaznik("1.0.0"),
		}, &wynik)

	const prob = 400
	for proba := 0; proba < prob; proba++ {
		stan := tekstZBazyApps(t, baza, `SELECT stan FROM wdrozenie_apps WHERE kod = ?`,
			wynik.Deployment.Id)
		if stan == string(shared.AppDeployStatusSucceeded) {
			return wynik.Deployment.Id
		}
		if stan == string(shared.AppDeployStatusFailed) {
			t.Fatalf("wdrożenie %s zakończyło się niepowodzeniem", wynik.Deployment.Id)
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("wdrożenie %s nie doszło do stanu końcowego", wynik.Deployment.Id)
	return ""
}

// TestWdrozenieZostawiaArtefaktZBajtamiIDziennik wykazuje, że po udanym
// przebiegu istnieje archiwum, które da się otworzyć i w którym leży treść
// warsztatu, oraz że dziennik przebiegu ma wiersze.
func TestWdrozenieZostawiaArtefaktZBajtamiIDziennik(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	okno := oknoModuluSprawdzianu(t, zmontowany, zycie)
	zapiszPlikWarsztatuOkna(t, zmontowany, zycie, okno,
		shared.AppWorkspaceLayerFrontend, "index.html", "<h1>Portal klienta</h1>")
	zapiszPlikWarsztatuOkna(t, zmontowany, zycie, okno,
		shared.AppWorkspaceLayerBackend, "serwer.go", "package main // usługi")

	kodWdrozenia := przeprowadzWdrozenieSprawdzianu(t, zmontowany, zycie, baza, okno)

	var artefakty shared.AppsArtifactListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsArtifactList,
		shared.AppsArtifactListRequest{
			WindowId: okno, DeploymentId: wskaznik(kodWdrozenia),
		}, &artefakty)
	if artefakty.Total != 1 {
		t.Fatalf("po wdrożeniu jest %d artefaktów zamiast jednego", artefakty.Total)
	}
	artefakt := artefakty.Artifacts[0]

	// Za wierszem artefaktu ma leżeć archiwum, które naprawdę się otwiera.
	bajty := bajtyPodOdwolaniem(t, katalog, artefakt.Path)
	if artefakt.SizeBytes == nil || int64(len(bajty)) != *artefakt.SizeBytes {
		t.Fatalf("rozmiar w wierszu (%v) rozjeżdża się z plikiem (%d)",
			artefakt.SizeBytes, len(bajty))
	}
	if artefakt.ChecksumSha256 == nil || *artefakt.ChecksumSha256 != sumaSha256(bajty) {
		t.Fatalf("suma kontrolna w wierszu nie opisuje bajtów pliku")
	}
	czytnik, err := zip.NewReader(bytes.NewReader(bajty), int64(len(bajty)))
	if err != nil {
		t.Fatalf("artefakt nie jest archiwum zip: %v", err)
	}
	wpisy := map[string]string{}
	for _, plik := range czytnik.File {
		strumien, err := plik.Open()
		if err != nil {
			t.Fatalf("nie można otworzyć wpisu %s: %v", plik.Name, err)
		}
		tresc, err := io.ReadAll(strumien)
		_ = strumien.Close()
		if err != nil {
			t.Fatalf("nie można odczytać wpisu %s: %v", plik.Name, err)
		}
		wpisy[plik.Name] = string(tresc)
	}
	if wpisy["frontend/index.html"] != "<h1>Portal klienta</h1>" {
		t.Fatalf("w archiwum nie ma treści pliku frontendu: %+v", wpisy)
	}
	if wpisy["backend/serwer.go"] != "package main // usługi" {
		t.Fatalf("w archiwum nie ma treści pliku backendu: %+v", wpisy)
	}

	var dziennik shared.AppsDeploymentLogReadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsDeploymentLogRead,
		shared.AppsDeploymentLogReadRequest{
			WindowId: okno, DeploymentId: kodWdrozenia,
		}, &dziennik)
	if dziennik.Total < 3 {
		t.Fatalf("dziennik przebiegu ma %d wierszy: %+v", dziennik.Total, dziennik.Lines)
	}
	if dziennik.Streaming {
		t.Fatal("odczyt dziennika obiecuje strumień, którego rdzeń nie nadaje")
	}
	if !strings.Contains(strings.Join(dziennik.Lines, "\n"), "artefakt") {
		t.Fatalf("dziennik nie odnotował złożenia artefaktu: %+v", dziennik.Lines)
	}

	var uslugi shared.AppsServiceLogReadResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsServiceLogRead,
		shared.AppsServiceLogReadRequest{WindowId: okno}, &uslugi)
	if uslugi.Total < dziennik.Total {
		t.Fatalf("dziennik usług okna (%d) jest krótszy niż dziennik jednego przebiegu (%d)",
			uslugi.Total, dziennik.Total)
	}
}

// ── Publisher Panel ─────────────────────────────────────────────────────────

// TestPakowaniePodpisIPublikacjaDajaPlikPodpisIPozycje jest sprawdzianem
// całego przebiegu wydawniczego: archiwum na dysku, podpis dający się
// zweryfikować kluczem publicznym i pozycja w tabeli katalogu rozszerzeń.
func TestPakowaniePodpisIPublikacjaDajaPlikPodpisIPozycje(t *testing.T) {
	zmontowany, zycie, katalog := zmontujDoPomiaruSkutku(t)
	baza := bazaSprawdzianuApps(t, katalog)

	okno := oknoModuluSprawdzianu(t, zmontowany, zycie)
	zapiszPlikWarsztatuOkna(t, zmontowany, zycie, okno,
		shared.AppWorkspaceLayerFrontend, "index.html", "<h1>Portal klienta</h1>")
	przeprowadzWdrozenieSprawdzianu(t, zmontowany, zycie, baza, okno)

	var pakiet shared.AppsPackageBuildResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPackageBuild,
		shared.AppsPackageBuildRequest{WindowId: okno}, &pakiet)
	if pakiet.Package.SizeBytes == nil || *pakiet.Package.SizeBytes == 0 {
		t.Fatalf("pakiet zameldowany bez rozmiaru: %+v", pakiet.Package)
	}

	sciezka := tekstZBazyApps(t, baza,
		`SELECT IFNULL(sciezka,'') FROM pakiet_apps WHERE identyfikator_zewnetrzny = ?`,
		pakiet.Package.Id)
	archiwum := bajtyPodOdwolaniem(t, katalog, sciezka)
	czytnik, err := zip.NewReader(bytes.NewReader(archiwum), int64(len(archiwum)))
	if err != nil {
		t.Fatalf("pakiet nie jest archiwum zip: %v", err)
	}
	maManifest := false
	maTrescProduktu := false
	for _, plik := range czytnik.File {
		if plik.Name == nazwaManifestuWPakiecieApp {
			maManifest = true
		}
		if plik.Name == "frontend/index.html" {
			maTrescProduktu = true
		}
	}
	if !maManifest || !maTrescProduktu {
		t.Fatalf("archiwum pakietu nie niesie manifestu i treści produktu: manifest=%v tresc=%v",
			maManifest, maTrescProduktu)
	}

	// Manifest wypełniony przez Operatora.
	var poManifescie shared.AppsPackageManifestSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPackageManifestSave,
		shared.AppsPackageManifestSaveRequest{
			WindowId: okno, PackageId: wskaznik(pakiet.Package.Id),
			Manifest: shared.AppPackageManifest{
				Identifier: "portal-klienta", Name: "Portal klienta", Version: "1.2.0",
				Kind: shared.ExtensionKindPlugin,
				Tools: []shared.ExtensionToolEntry{
					{Name: "portal.zamowienia", Kind: shared.ExtensionToolKindTool},
				},
			},
		}, &poManifescie)
	if poManifescie.Package.Manifest == nil ||
		poManifescie.Package.Manifest.Identifier != "portal-klienta" {
		t.Fatalf("manifest nie wrócił po zapisie: %+v", poManifescie.Package.Manifest)
	}

	var walidacja shared.AppsPackageValidateResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPackageValidate,
		shared.AppsPackageValidateRequest{
			WindowId: okno, PackageId: pakiet.Package.Id,
		}, &walidacja)
	kody := map[string]int{}
	for _, zastrzezenie := range walidacja.Issues {
		kody[zastrzezenie.Code]++
	}
	if kody[kodPakietuBezPodpisuApp] == 0 {
		t.Fatalf("walidator nie zauważył braku podpisu: %+v", walidacja.Issues)
	}
	if kody[kodPakietuTozsamoscApp] > 0 || kody[kodPakietuWersjaApp] > 0 {
		t.Fatalf("walidator ma zastrzeżenia do poprawnego manifestu: %+v", walidacja.Issues)
	}

	var podpis shared.AppsPackageSignResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPackageSign,
		shared.AppsPackageSignRequest{
			WindowId: okno, PackageId: pakiet.Package.Id,
			SigningKeyRef: "sejf:wydawca/danaco",
		}, &podpis)
	if !podpis.Signature.Signed || !podpis.Signature.Verified {
		t.Fatalf("podpis nie przeszedł weryfikacji: %+v", podpis.Signature)
	}
	if podpis.Signature.ChecksumSha256 == nil ||
		*podpis.Signature.ChecksumSha256 != sumaSha256(archiwum) {
		t.Fatal("suma kontrolna podpisu nie opisuje bajtów archiwum")
	}

	// Najtwardsza miara podpisu: sprawdzian sam weryfikuje go kluczem
	// publicznym odłożonym w wierszu — nie wierzy polu `verified`.
	surowy := tekstZBazyApps(t, baza,
		`SELECT IFNULL(podpis,'') FROM pakiet_apps WHERE identyfikator_zewnetrzny = ?`,
		pakiet.Package.Id)
	var zapisany struct {
		SignatureBase64 string `json:"signatureBase64"`
		PublicKeyBase64 string `json:"publicKeyBase64"`
	}
	if err := json.Unmarshal([]byte(surowy), &zapisany); err != nil {
		t.Fatalf("podpis w bazie jest nieczytelny: %v", err)
	}
	bajtyPodpisu, err := base64.StdEncoding.DecodeString(zapisany.SignatureBase64)
	if err != nil {
		t.Fatalf("podpis nie jest base64: %v", err)
	}
	bajtyKlucza, err := base64.StdEncoding.DecodeString(zapisany.PublicKeyBase64)
	if err != nil {
		t.Fatalf("klucz publiczny nie jest base64: %v", err)
	}
	suma := sha256.Sum256(archiwum)
	if !ed25519.Verify(ed25519.PublicKey(bajtyKlucza), suma[:], bajtyPodpisu) {
		t.Fatal("podpis odłożony w bazie nie weryfikuje się kluczem publicznym wydawcy")
	}
	if hex.EncodeToString(suma[:]) != *podpis.Signature.ChecksumSha256 {
		t.Fatal("suma kontrolna liczona niezależnie różni się od oddanej przez rdzeń")
	}

	// Klucz wydawcy leży w sejfie, a nie w bazie modułu.
	sejf, err := os.ReadFile(filepath.Join(katalog, "poswiadczenia.json"))
	if err == nil && !strings.Contains(string(sejf), "wydawca/danaco") {
		t.Fatal("klucz wydawcy nie trafił do sejfu poświadczeń")
	}

	var publikacja shared.AppsPackagePublishResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPackagePublish,
		shared.AppsPackagePublishRequest{
			WindowId: okno, PackageId: pakiet.Package.Id,
			ReleaseNotes: wskaznik("pierwsze wydanie"),
		}, &publikacja)

	if publikacja.Extension.Code != "portal-klienta" {
		t.Fatalf("pozycja katalogu ma kod %q zamiast portal-klienta", publikacja.Extension.Code)
	}
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM rozszerzenie WHERE kod = 'portal-klienta' AND zainstalowane = 1`,
	); ile != 1 {
		t.Fatal("pozycji opublikowanej nie ma w tabeli katalogu rozszerzeń")
	}
	if wersja := tekstZBazyApps(t, baza,
		`SELECT IFNULL(wersja,'') FROM rozszerzenie WHERE kod = 'portal-klienta'`); wersja != "1.2.0" {
		t.Fatalf("pozycja katalogu ma wersję %q zamiast 1.2.0", wersja)
	}
	if kodPakietu := tekstZBazyApps(t, baza,
		`SELECT IFNULL(rozszerzenie_kod,'') FROM pakiet_apps WHERE identyfikator_zewnetrzny = ?`,
		pakiet.Package.Id); kodPakietu == "" {
		t.Fatal("pakiet nie zapamiętał pozycji, która z niego powstała")
	}

	// Druga publikacja tego samego kodu podnosi wydanie, nie zakłada bliźniaka.
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPackageManifestSave,
		shared.AppsPackageManifestSaveRequest{
			WindowId: okno, PackageId: wskaznik(pakiet.Package.Id),
			Manifest: shared.AppPackageManifest{
				Identifier: "portal-klienta", Name: "Portal klienta", Version: "1.3.0",
				Kind: shared.ExtensionKindPlugin,
			},
		}, &poManifescie)
	wykonajUdana(t, zmontowany, zycie, shared.CommandAppsPackagePublish,
		shared.AppsPackagePublishRequest{
			WindowId: okno, PackageId: pakiet.Package.Id,
		}, &publikacja)
	if ile := wierszyApps(t, baza,
		`SELECT COUNT(*) FROM rozszerzenie WHERE kod = 'portal-klienta'`); ile != 1 {
		t.Fatalf("druga publikacja założyła %d pozycji katalogu zamiast podnieść jedną", ile)
	}
	if wersja := tekstZBazyApps(t, baza,
		`SELECT IFNULL(wersja,'') FROM rozszerzenie WHERE kod = 'portal-klienta'`); wersja != "1.3.0" {
		t.Fatalf("druga publikacja zostawiła w katalogu wersję %q", wersja)
	}
}

// TestPakowanieBezArtefaktuOdmawiaZPowodem wykazuje, że pakiet nie powstaje
// znikąd: okno bez udanego wdrożenia dostaje odmowę, a nie pusty pakiet.
func TestPakowanieBezArtefaktuOdmawiaZPowodem(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	odmowa := wykonajOdmowna(t, zmontowany, zycie, shared.CommandAppsPackageBuild,
		shared.AppsPackageBuildRequest{WindowId: oknoSprawdzianuApps})
	if odmowa.Code != shared.ErrorCodeValidationFailed {
		t.Fatalf("pakowanie bez artefaktu wróciło kodem %q", odmowa.Code)
	}
	if !strings.Contains(odmowa.Message, "artefakt") {
		t.Fatalf("odmowa nie nazywa braku: %q", odmowa.Message)
	}
}
