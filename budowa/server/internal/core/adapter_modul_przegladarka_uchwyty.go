// Odpowiedzialność pliku: port modułu Browser i wpięcie jego sześciu komend
// do rejestru (`navigate`, `snapshot.get`, `source.add`, `note.add`,
// `source.list`, `note.list`).
//
// Nawigacja rozgłasza `browser.page.changed`. Kontrakt niesie to zdarzenie
// (`shared.EventBrowserPageChanged`), a klient je subskrybuje, więc `navigate`
// po udanym pobraniu strony rozgłasza migawkę po zmianie — treść widoczną
// jednocześnie Operatorowi i modelowi. Rozgłoszenie jedzie tym samym emiterem
// rdzenia, co pozostałe zmiany obszarów (wzór: Automatyki).
//
// Źródło i notatka nie rozgłaszają zdarzeń. W odróżnieniu od nawigacji,
// `browser.source.add` i `browser.note.add` odkładają wynik wprost do
// odpowiedzi — kontrakt nie niesie dla nich żadnego zdarzenia domenowego
// (`shared/contract.go` nie ma `browser.source.changed` ani podobnego), więc
// port ich nie wymyśla.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Przegladarka jest portem modułu Browser.
type Przegladarka interface {
	Nawiguj(ctx context.Context, z shared.BrowserNavigateRequest) (shared.BrowserNavigateResponse, error)
	Migawka(ctx context.Context, z shared.BrowserSnapshotGetRequest) (shared.BrowserSnapshotGetResponse, error)
	DodajZrodlo(ctx context.Context, z shared.BrowserSourceAddRequest) (shared.BrowserSourceAddResponse, error)
	DodajNotatke(ctx context.Context, z shared.BrowserNoteAddRequest) (shared.BrowserNoteAddResponse, error)
	WykazZrodel(ctx context.Context, z shared.BrowserSourceListRequest) (shared.BrowserSourceListResponse, error)
	WykazNotatek(ctx context.Context, z shared.BrowserNoteListRequest) (shared.BrowserNoteListResponse, error)

	// Karty, grupy kart i przestrzenie robocze.
	OtworzKarte(ctx context.Context, z shared.BrowserTabOpenRequest) (shared.BrowserTabOpenResponse, error)
	WykazKart(ctx context.Context, z shared.BrowserTabListRequest) (shared.BrowserTabListResponse, error)
	ZmienKarte(ctx context.Context, z shared.BrowserTabUpdateRequest) (shared.BrowserTabUpdateResponse, error)
	ZamknijKarte(ctx context.Context, z shared.BrowserTabCloseRequest) (shared.BrowserTabCloseResponse, error)
	UstawGrupeKart(ctx context.Context, z shared.BrowserTabGroupSetRequest) (shared.BrowserTabGroupSetResponse, error)
	ZapiszPrzestrzen(ctx context.Context, z shared.BrowserWorkspaceSaveRequest) (shared.BrowserWorkspaceSaveResponse, error)
	WykazPrzestrzeni(ctx context.Context, z shared.BrowserWorkspaceListRequest) (shared.BrowserWorkspaceListResponse, error)
	OtworzPrzestrzen(ctx context.Context, z shared.BrowserWorkspaceOpenRequest) (shared.BrowserWorkspaceOpenResponse, error)
	UsunPrzestrzen(ctx context.Context, z shared.BrowserWorkspaceRemoveRequest) (shared.BrowserWorkspaceRemoveResponse, error)

	// Monitory zmian strony.
	ZalozMonitor(ctx context.Context, z shared.BrowserMonitorAddRequest) (shared.BrowserMonitorAddResponse, error)
	WykazMonitorow(ctx context.Context, z shared.BrowserMonitorListRequest) (shared.BrowserMonitorListResponse, error)
	SprawdzMonitor(ctx context.Context, z shared.BrowserMonitorCheckRequest) (shared.BrowserMonitorCheckResponse, error)
	ZdejmijMonitor(ctx context.Context, z shared.BrowserMonitorRemoveRequest) (shared.BrowserMonitorRemoveResponse, error)

	// Kanały i kolejka czytania.
	SubskrybujKanal(ctx context.Context, z shared.BrowserFeedSubscribeRequest) (shared.BrowserFeedSubscribeResponse, error)
	WykazKanalow(ctx context.Context, z shared.BrowserFeedListRequest) (shared.BrowserFeedListResponse, error)
	ZdejmijKanal(ctx context.Context, z shared.BrowserFeedRemoveRequest) (shared.BrowserFeedRemoveResponse, error)
	OdlozDoCzytania(ctx context.Context, z shared.BrowserReadlistAddRequest) (shared.BrowserReadlistAddResponse, error)
	WykazCzytania(ctx context.Context, z shared.BrowserReadlistListRequest) (shared.BrowserReadlistListResponse, error)
	ZdejmijZCzytania(ctx context.Context, z shared.BrowserReadlistRemoveRequest) (shared.BrowserReadlistRemoveResponse, error)

	// Zakładki.
	DodajZakladke(ctx context.Context, z shared.BrowserBookmarkAddRequest) (shared.BrowserBookmarkAddResponse, error)
	WykazZakladek(ctx context.Context, z shared.BrowserBookmarkListRequest) (shared.BrowserBookmarkListResponse, error)
	UsunZakladke(ctx context.Context, z shared.BrowserBookmarkRemoveRequest) (shared.BrowserBookmarkRemoveResponse, error)

	// Porządkowanie zebranego materiału.
	UsunZrodlo(ctx context.Context, z shared.BrowserSourceRemoveRequest) (shared.BrowserSourceRemoveResponse, error)
	ZmienNotatke(ctx context.Context, z shared.BrowserNoteUpdateRequest) (shared.BrowserNoteUpdateResponse, error)
	UstawZestawZrodel(ctx context.Context, z shared.BrowserSourceGroupSetRequest) (shared.BrowserSourceGroupSetResponse, error)
	WykazZestawowZrodel(ctx context.Context, z shared.BrowserSourceGroupListRequest) (shared.BrowserSourceGroupListResponse, error)
	UstawWatekNotatek(ctx context.Context, z shared.BrowserNoteThreadSetRequest) (shared.BrowserNoteThreadSetResponse, error)
	WykazWatkowNotatek(ctx context.Context, z shared.BrowserNoteThreadListRequest) (shared.BrowserNoteThreadListResponse, error)

	// Narzędzia inspekcyjne i praca na stronie uruchomionej.
	ZbadajDrzewo(ctx context.Context, z shared.BrowserDomInspectRequest) (shared.BrowserDomInspectResponse, error)
	RejestrSieciowy(ctx context.Context, z shared.BrowserNetworkHarRequest) (shared.BrowserNetworkHarResponse, error)
	OdczytajKonsole(ctx context.Context, z shared.BrowserConsoleReadRequest) (shared.BrowserConsoleReadResponse, error)
	EmulujUrzadzenie(ctx context.Context, z shared.BrowserDeviceEmulateRequest) (shared.BrowserDeviceEmulateResponse, error)
	Przewin(ctx context.Context, z shared.BrowserScrollRequest) (shared.BrowserScrollResponse, error)

	// Materiał wizualny i wytwory sesji.
	WykonajZrzut(ctx context.Context, z shared.BrowserScreenshotCaptureRequest) (shared.BrowserScreenshotCaptureResponse, error)
	OdczytajZrzut(ctx context.Context, z shared.BrowserSnapshotScreenshotGetRequest) (shared.BrowserSnapshotScreenshotGetResponse, error)
	DodajWytwor(ctx context.Context, z shared.BrowserArtifactAddRequest) (shared.BrowserArtifactAddResponse, error)

	// Pobrania, makra i granice Wykonawcy.
	WykazPobran(ctx context.Context, z shared.BrowserDownloadListRequest) (shared.BrowserDownloadListResponse, error)
	SterujPobraniem(ctx context.Context, z shared.BrowserDownloadControlRequest) (shared.BrowserDownloadControlResponse, error)
	NagrywajMakro(ctx context.Context, z shared.BrowserMacroRecordRequest) (shared.BrowserMacroRecordResponse, error)
	UstawGraniceWykonawcy(ctx context.Context, z shared.BrowserExecutorLimitsSetRequest) (shared.BrowserExecutorLimitsSetResponse, error)
	OdczytajGraniceWykonawcy(ctx context.Context, z shared.BrowserExecutorLimitsGetRequest) (shared.BrowserExecutorLimitsGetResponse, error)
}

// zarejestrujPrzegladarke wpina sześć komend modułu Browser.
func zarejestrujPrzegladarke(r *Rejestr, m Przegladarka, e *emiter) {
	if r == nil || m == nil {
		return
	}

	r.Zarejestruj(shared.CommandBrowserNavigate,
		obsluz(func(ctx context.Context, z shared.BrowserNavigateRequest) (shared.BrowserNavigateResponse, error) {
			odpowiedz, err := m.Nawiguj(ctx, z)
			if err == nil {
				e.stronaPrzegladarki(odpowiedz.Snapshot, powodPrzejscia)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandBrowserSnapshotGet, obsluz(m.Migawka))
	r.Zarejestruj(shared.CommandBrowserSourceAdd, obsluz(m.DodajZrodlo))
	r.Zarejestruj(shared.CommandBrowserNoteAdd, obsluz(m.DodajNotatke))
	// Dwa wykazy dopisane do tej samej funkcji, nie do własnej: rejestr ma
	// jedno miejsce wiążące nazwy komend modułu z metodami portu.
	r.Zarejestruj(shared.CommandBrowserSourceList, obsluz(m.WykazZrodel))
	r.Zarejestruj(shared.CommandBrowserNoteList, obsluz(m.WykazNotatek))

	// Rząd kart i przestrzenie robocze. Trzy z tych komend rozgłaszają
	// `browser.tab.changed`: otwarcie, zmiana stanu i zamknięcie karty są
	// zmianami widocznymi w oknie każdego klienta patrzącego na tę sesję.
	r.Zarejestruj(shared.CommandBrowserTabOpen,
		obsluz(func(ctx context.Context, z shared.BrowserTabOpenRequest) (shared.BrowserTabOpenResponse, error) {
			odpowiedz, err := m.OtworzKarte(ctx, z)
			if err == nil {
				e.kartaPrzegladarki(odpowiedz.Tab, shared.ChangeKindCreated)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandBrowserTabList, obsluz(m.WykazKart))
	r.Zarejestruj(shared.CommandBrowserTabUpdate,
		obsluz(func(ctx context.Context, z shared.BrowserTabUpdateRequest) (shared.BrowserTabUpdateResponse, error) {
			odpowiedz, err := m.ZmienKarte(ctx, z)
			if err == nil {
				e.kartaPrzegladarki(odpowiedz.Tab, shared.ChangeKindUpdated)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandBrowserTabClose, obsluz(m.ZamknijKarte))
	r.Zarejestruj(shared.CommandBrowserTabGroupSet, obsluz(m.UstawGrupeKart))
	r.Zarejestruj(shared.CommandBrowserWorkspaceSave, obsluz(m.ZapiszPrzestrzen))
	r.Zarejestruj(shared.CommandBrowserWorkspaceList, obsluz(m.WykazPrzestrzeni))
	r.Zarejestruj(shared.CommandBrowserWorkspaceOpen, obsluz(m.OtworzPrzestrzen))
	r.Zarejestruj(shared.CommandBrowserWorkspaceRemove, obsluz(m.UsunPrzestrzen))

	// Monitory. Sprawdzenie, które wykryło zmianę, rozgłasza
	// `browser.monitor.changed` — Capture & Monitor Panel dowiaduje się o niej
	// bez odpytywania, a alert jest treścią zdarzenia, nie odpowiedzi.
	r.Zarejestruj(shared.CommandBrowserMonitorAdd, obsluz(m.ZalozMonitor))
	r.Zarejestruj(shared.CommandBrowserMonitorList, obsluz(m.WykazMonitorow))
	r.Zarejestruj(shared.CommandBrowserMonitorCheck,
		obsluz(func(ctx context.Context, z shared.BrowserMonitorCheckRequest) (shared.BrowserMonitorCheckResponse, error) {
			odpowiedz, err := m.SprawdzMonitor(ctx, z)
			if err == nil && odpowiedz.Changed {
				e.monitorPrzegladarki(odpowiedz.Monitor, odpowiedz.Diff)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandBrowserMonitorRemove, obsluz(m.ZdejmijMonitor))

	// Kanały i kolejka czytania.
	r.Zarejestruj(shared.CommandBrowserFeedSubscribe, obsluz(m.SubskrybujKanal))
	r.Zarejestruj(shared.CommandBrowserFeedList, obsluz(m.WykazKanalow))
	r.Zarejestruj(shared.CommandBrowserFeedRemove, obsluz(m.ZdejmijKanal))
	r.Zarejestruj(shared.CommandBrowserReadlistAdd, obsluz(m.OdlozDoCzytania))
	r.Zarejestruj(shared.CommandBrowserReadlistList, obsluz(m.WykazCzytania))
	r.Zarejestruj(shared.CommandBrowserReadlistRemove, obsluz(m.ZdejmijZCzytania))

	// Zakładki.
	r.Zarejestruj(shared.CommandBrowserBookmarkAdd, obsluz(m.DodajZakladke))
	r.Zarejestruj(shared.CommandBrowserBookmarkList, obsluz(m.WykazZakladek))
	r.Zarejestruj(shared.CommandBrowserBookmarkRemove, obsluz(m.UsunZakladke))

	// Porządkowanie materiału.
	r.Zarejestruj(shared.CommandBrowserSourceRemove, obsluz(m.UsunZrodlo))
	r.Zarejestruj(shared.CommandBrowserNoteUpdate, obsluz(m.ZmienNotatke))
	r.Zarejestruj(shared.CommandBrowserSourceGroupSet, obsluz(m.UstawZestawZrodel))
	r.Zarejestruj(shared.CommandBrowserSourceGroupList, obsluz(m.WykazZestawowZrodel))
	r.Zarejestruj(shared.CommandBrowserNoteThreadSet, obsluz(m.UstawWatekNotatek))
	r.Zarejestruj(shared.CommandBrowserNoteThreadList, obsluz(m.WykazWatkowNotatek))

	// Narzędzia inspekcyjne. Emulacja i przewinięcie zmieniają wspólny podgląd,
	// więc rozgłaszają `browser.page.changed` z powodem `interaction` — to nie
	// jest przejście pod nowy adres, tylko zmiana stanu strony już otwartej.
	r.Zarejestruj(shared.CommandBrowserDomInspect, obsluz(m.ZbadajDrzewo))
	r.Zarejestruj(shared.CommandBrowserNetworkHar, obsluz(m.RejestrSieciowy))
	r.Zarejestruj(shared.CommandBrowserConsoleRead, obsluz(m.OdczytajKonsole))
	r.Zarejestruj(shared.CommandBrowserDeviceEmulate,
		obsluz(func(ctx context.Context, z shared.BrowserDeviceEmulateRequest) (shared.BrowserDeviceEmulateResponse, error) {
			odpowiedz, err := m.EmulujUrzadzenie(ctx, z)
			if err == nil && odpowiedz.Snapshot != nil {
				e.stronaPrzegladarki(*odpowiedz.Snapshot, powodInterakcji)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandBrowserScroll,
		obsluz(func(ctx context.Context, z shared.BrowserScrollRequest) (shared.BrowserScrollResponse, error) {
			odpowiedz, err := m.Przewin(ctx, z)
			if err == nil {
				e.stronaPrzegladarki(odpowiedz.Snapshot, powodInterakcji)
			}
			return odpowiedz, err
		}))

	// Materiał wizualny i wytwory.
	r.Zarejestruj(shared.CommandBrowserScreenshotCapture, obsluz(m.WykonajZrzut))
	r.Zarejestruj(shared.CommandBrowserSnapshotScreenshotGet, obsluz(m.OdczytajZrzut))
	r.Zarejestruj(shared.CommandBrowserArtifactAdd, obsluz(m.DodajWytwor))

	// Pobrania, makra i granice Wykonawcy. Sterowanie pobraniem rozgłasza
	// `browser.download.changed`: postęp i stan pobrania są tym, co menedżer
	// pokazuje na żywo.
	r.Zarejestruj(shared.CommandBrowserDownloadList, obsluz(m.WykazPobran))
	r.Zarejestruj(shared.CommandBrowserDownloadControl,
		obsluz(func(ctx context.Context, z shared.BrowserDownloadControlRequest) (shared.BrowserDownloadControlResponse, error) {
			odpowiedz, err := m.SterujPobraniem(ctx, z)
			if err == nil {
				e.pobraniePrzegladarki(odpowiedz.Download)
			}
			return odpowiedz, err
		}))
	r.Zarejestruj(shared.CommandBrowserMacroRecord, obsluz(m.NagrywajMakro))
	r.Zarejestruj(shared.CommandBrowserExecutorLimitsSet, obsluz(m.UstawGraniceWykonawcy))
	r.Zarejestruj(shared.CommandBrowserExecutorLimitsGet, obsluz(m.OdczytajGraniceWykonawcy))
}

// powodInterakcji nazywa zmianę widoku strony już wczytanej — w odróżnieniu od
// `powodPrzejscia`, którym jedzie wejście pod nowy adres. Kontrakt
// `BrowserPageChangedEvent.reason` rozróżnia te dwa źródła zmiany.
const powodInterakcji = "interaction"

// kartaPrzegladarki rozgłasza `browser.tab.changed`. Rodzaj zmiany jest polem
// kontraktu (`ChangeKind`), a nie dowolnym napisem: klient odsiewa po nim
// założenie karty od zmiany jej stanu.
func (e *emiter) kartaPrzegladarki(karta shared.BrowserTab, rodzaj shared.ChangeKind) {
	kopia := karta
	tresc := shared.BrowserTabChangedEvent{WindowId: karta.WindowId, Change: rodzaj, Tab: &kopia}
	e.wyslij(shared.EventBrowserTabChanged, "", tresc)
}

// monitorPrzegladarki rozgłasza `browser.monitor.changed` — wykrytą zmianę
// pilnowanej strony wraz z jej miarą. Zdarzenie idzie wyłącznie przy zmianie
// naprawdę wykrytej; rozgłoszenie przy każdym sprawdzeniu byłoby alarmem bez
// zdarzenia.
func (e *emiter) monitorPrzegladarki(monitor shared.BrowserMonitor, roznica *shared.BrowserContentDiff) {
	tresc := shared.BrowserMonitorChangedEvent{
		MonitorId: monitor.Id, WindowId: monitor.WindowId, Monitor: monitor,
	}
	if roznica != nil {
		tresc.Diff = *roznica
	}
	e.wyslij(shared.EventBrowserMonitorChanged, "", tresc)
}

// pobraniePrzegladarki rozgłasza `browser.download.changed`.
func (e *emiter) pobraniePrzegladarki(pobranie shared.BrowserDownload) {
	tresc := shared.BrowserDownloadChangedEvent{
		DownloadId: pobranie.Id, WindowId: pobranie.WindowId, Download: pobranie,
	}
	e.wyslij(shared.EventBrowserDownloadChanged, "", tresc)
}

// powodPrzejscia nazywa powód migawki rozgłaszanej po `browser.navigate`:
// przejście pod adres, nie interakcja Operatora ze stroną już wczytaną.
// Kontrakt `BrowserPageChangedEvent.reason` rozróżnia te dwa źródła zmiany.
const powodPrzejscia = "navigation"

// stronaPrzegladarki rozgłasza `browser.page.changed` — migawkę strony po
// zmianie. Metoda emitera per moduł (wzór: `przebiegAutomatyki`), zadeklarowana
// tu, a nie w `zdarzenia.go`, bo to obszar Browser nazywa własne zdarzenie.
// Zdarzenie niesie okno migawki i jedzie bez wskazania sesji — okno
// przeglądarki nie jest bytem karty sesji.
func (e *emiter) stronaPrzegladarki(migawka shared.BrowserSnapshot, powod string) {
	tresc := shared.BrowserPageChangedEvent{WindowId: migawka.WindowId, Snapshot: migawka}
	if powod != "" {
		p := powod
		tresc.Reason = &p
	}
	e.wyslij(shared.EventBrowserPageChanged, "", tresc)
}
