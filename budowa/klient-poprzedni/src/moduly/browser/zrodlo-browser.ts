import {
  Command,
  EventType,
  type BrowserNavigateRequest,
  type BrowserNavigateResponse,
  type BrowserNoteAddRequest,
  type BrowserNoteAddResponse,
  type BrowserNoteListRequest,
  type BrowserNoteListResponse,
  type BrowserPageChangedEvent,
  type BrowserSnapshotGetRequest,
  type BrowserSnapshotGetResponse,
  type BrowserSourceAddRequest,
  type BrowserSourceAddResponse,
  type BrowserSourceListRequest,
  type BrowserSourceListResponse,
  type BrowserTabOpenRequest,
  type BrowserTabOpenResponse,
  type BrowserTabListRequest,
  type BrowserTabListResponse,
  type BrowserTabUpdateRequest,
  type BrowserTabUpdateResponse,
  type BrowserTabCloseRequest,
  type BrowserTabCloseResponse,
  type BrowserTabGroupSetRequest,
  type BrowserTabGroupSetResponse,
  type BrowserWorkspaceSaveRequest,
  type BrowserWorkspaceSaveResponse,
  type BrowserWorkspaceListRequest,
  type BrowserWorkspaceListResponse,
  type BrowserWorkspaceOpenRequest,
  type BrowserWorkspaceOpenResponse,
  type BrowserWorkspaceRemoveRequest,
  type BrowserWorkspaceRemoveResponse,
  type BrowserMonitorAddRequest,
  type BrowserMonitorAddResponse,
  type BrowserMonitorListRequest,
  type BrowserMonitorListResponse,
  type BrowserMonitorCheckRequest,
  type BrowserMonitorCheckResponse,
  type BrowserMonitorRemoveRequest,
  type BrowserMonitorRemoveResponse,
  type BrowserFeedSubscribeRequest,
  type BrowserFeedSubscribeResponse,
  type BrowserFeedListRequest,
  type BrowserFeedListResponse,
  type BrowserFeedRemoveRequest,
  type BrowserFeedRemoveResponse,
  type BrowserReadlistAddRequest,
  type BrowserReadlistAddResponse,
  type BrowserReadlistListRequest,
  type BrowserReadlistListResponse,
  type BrowserReadlistRemoveRequest,
  type BrowserReadlistRemoveResponse,
  type BrowserBookmarkAddRequest,
  type BrowserBookmarkAddResponse,
  type BrowserBookmarkListRequest,
  type BrowserBookmarkListResponse,
  type BrowserBookmarkRemoveRequest,
  type BrowserBookmarkRemoveResponse,
  type BrowserSourceRemoveRequest,
  type BrowserSourceRemoveResponse,
  type BrowserNoteUpdateRequest,
  type BrowserNoteUpdateResponse,
  type BrowserSourceGroupSetRequest,
  type BrowserSourceGroupSetResponse,
  type BrowserSourceGroupListRequest,
  type BrowserSourceGroupListResponse,
  type BrowserNoteThreadSetRequest,
  type BrowserNoteThreadSetResponse,
  type BrowserNoteThreadListRequest,
  type BrowserNoteThreadListResponse,
  type BrowserDomInspectRequest,
  type BrowserDomInspectResponse,
  type BrowserNetworkHarRequest,
  type BrowserNetworkHarResponse,
  type BrowserConsoleReadRequest,
  type BrowserConsoleReadResponse,
  type BrowserDeviceEmulateRequest,
  type BrowserDeviceEmulateResponse,
  type BrowserScrollRequest,
  type BrowserScrollResponse,
  type BrowserScreenshotCaptureRequest,
  type BrowserScreenshotCaptureResponse,
  type BrowserSnapshotScreenshotGetRequest,
  type BrowserSnapshotScreenshotGetResponse,
  type BrowserArtifactAddRequest,
  type BrowserArtifactAddResponse,
  type BrowserDownloadListRequest,
  type BrowserDownloadListResponse,
  type BrowserDownloadControlRequest,
  type BrowserDownloadControlResponse,
  type BrowserMacroRecordRequest,
  type BrowserMacroRecordResponse,
  type BrowserExecutorLimitsSetRequest,
  type BrowserExecutorLimitsSetResponse,
  type BrowserExecutorLimitsGetRequest,
  type BrowserExecutorLimitsGetResponse,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Sześć komend obszaru `browser.*` widzianych przez okna modułu.
 *
 * Plik jest warstwą wywołań wraz ze sprawdzianem kształtu odpowiedzi. Źródło
 * nie ma stanu i nie buduje elementów — zebrane źródła i notatki mieszkają
 * w `stan-przegladania.ts`, żeby trzy okna modułu patrzyły na jeden zbiór,
 * a nie na trzy kopie.
 *
 * `source.list` i `note.list` czytają dokładnie to, co moduł zapisał przez
 * `browser.source.add` i `browser.note.add`, i nic ponad to.
 *
 * Odmowa wykazu ma dwa znaczenia, które okno rozróżnia. `not_found` znaczy
 * „rdzeń nie zna tego okna przeglądania", a nie „nic nie zebrano" — wykaz
 * pusty przychodzi ze statusem udanym i pustą tablicą.
 *
 * Żadne wywołanie nie rzuca wyjątkiem i nie odrzuca obietnicy: niepowodzenie
 * wraca polem `blad` wyniku. Nazwy komend biorą się wyłącznie ze
 * stałych kontraktu.
 */
export interface ZrodloBrowser {
  przejdz(zadanie: BrowserNavigateRequest): Promise<Wynik<BrowserNavigateResponse>>;
  migawka(zadanie: BrowserSnapshotGetRequest): Promise<Wynik<BrowserSnapshotGetResponse>>;
  dodajZrodlo(zadanie: BrowserSourceAddRequest): Promise<Wynik<BrowserSourceAddResponse>>;
  dodajNotatke(zadanie: BrowserNoteAddRequest): Promise<Wynik<BrowserNoteAddResponse>>;
  /** `browser.source.list`; źródła okna od najnowszego. */
  wykazZrodel(zadanie: BrowserSourceListRequest): Promise<Wynik<BrowserSourceListResponse>>;
  /** `browser.note.list`; notatki okna od najświeższej. */
  wykazNotatek(zadanie: BrowserNoteListRequest): Promise<Wynik<BrowserNoteListResponse>>;

  // ── Rodziny dołożone ponad migawkę, źródło i notatkę ─────────────────────
  // Każda pozycja jest jedną komendą kontraktu. Warstwa nie ma stanu i nie
  // buduje elementów — sprawdza wyłącznie kształt odpowiedzi rdzenia.
  /** `BrowserTabOpen`. */
  otworzKarte(zadanie: BrowserTabOpenRequest): Promise<Wynik<BrowserTabOpenResponse>>;
  /** `BrowserTabList`. */
  wykazKart(zadanie: BrowserTabListRequest): Promise<Wynik<BrowserTabListResponse>>;
  /** `BrowserTabUpdate`. */
  zmienKarte(zadanie: BrowserTabUpdateRequest): Promise<Wynik<BrowserTabUpdateResponse>>;
  /** `BrowserTabClose`. */
  zamknijKarte(zadanie: BrowserTabCloseRequest): Promise<Wynik<BrowserTabCloseResponse>>;
  /** `BrowserTabGroupSet`. */
  ustawGrupeKart(zadanie: BrowserTabGroupSetRequest): Promise<Wynik<BrowserTabGroupSetResponse>>;
  /** `BrowserWorkspaceSave`. */
  zapiszPrzestrzen(zadanie: BrowserWorkspaceSaveRequest): Promise<Wynik<BrowserWorkspaceSaveResponse>>;
  /** `BrowserWorkspaceList`. */
  wykazPrzestrzeni(zadanie: BrowserWorkspaceListRequest): Promise<Wynik<BrowserWorkspaceListResponse>>;
  /** `BrowserWorkspaceOpen`. */
  otworzPrzestrzen(zadanie: BrowserWorkspaceOpenRequest): Promise<Wynik<BrowserWorkspaceOpenResponse>>;
  /** `BrowserWorkspaceRemove`. */
  usunPrzestrzen(zadanie: BrowserWorkspaceRemoveRequest): Promise<Wynik<BrowserWorkspaceRemoveResponse>>;
  /** `BrowserMonitorAdd`. */
  zalozMonitor(zadanie: BrowserMonitorAddRequest): Promise<Wynik<BrowserMonitorAddResponse>>;
  /** `BrowserMonitorList`. */
  wykazMonitorow(zadanie: BrowserMonitorListRequest): Promise<Wynik<BrowserMonitorListResponse>>;
  /** `BrowserMonitorCheck`. */
  sprawdzMonitor(zadanie: BrowserMonitorCheckRequest): Promise<Wynik<BrowserMonitorCheckResponse>>;
  /** `BrowserMonitorRemove`. */
  zdejmijMonitor(zadanie: BrowserMonitorRemoveRequest): Promise<Wynik<BrowserMonitorRemoveResponse>>;
  /** `BrowserFeedSubscribe`. */
  subskrybujKanal(zadanie: BrowserFeedSubscribeRequest): Promise<Wynik<BrowserFeedSubscribeResponse>>;
  /** `BrowserFeedList`. */
  wykazKanalow(zadanie: BrowserFeedListRequest): Promise<Wynik<BrowserFeedListResponse>>;
  /** `BrowserFeedRemove`. */
  zdejmijKanal(zadanie: BrowserFeedRemoveRequest): Promise<Wynik<BrowserFeedRemoveResponse>>;
  /** `BrowserReadlistAdd`. */
  odlozDoCzytania(zadanie: BrowserReadlistAddRequest): Promise<Wynik<BrowserReadlistAddResponse>>;
  /** `BrowserReadlistList`. */
  wykazCzytania(zadanie: BrowserReadlistListRequest): Promise<Wynik<BrowserReadlistListResponse>>;
  /** `BrowserReadlistRemove`. */
  zdejmijZCzytania(zadanie: BrowserReadlistRemoveRequest): Promise<Wynik<BrowserReadlistRemoveResponse>>;
  /** `BrowserBookmarkAdd`. */
  dodajZakladke(zadanie: BrowserBookmarkAddRequest): Promise<Wynik<BrowserBookmarkAddResponse>>;
  /** `BrowserBookmarkList`. */
  wykazZakladek(zadanie: BrowserBookmarkListRequest): Promise<Wynik<BrowserBookmarkListResponse>>;
  /** `BrowserBookmarkRemove`. */
  usunZakladke(zadanie: BrowserBookmarkRemoveRequest): Promise<Wynik<BrowserBookmarkRemoveResponse>>;
  /** `BrowserSourceRemove`. */
  usunZrodlo(zadanie: BrowserSourceRemoveRequest): Promise<Wynik<BrowserSourceRemoveResponse>>;
  /** `BrowserNoteUpdate`. */
  zmienNotatke(zadanie: BrowserNoteUpdateRequest): Promise<Wynik<BrowserNoteUpdateResponse>>;
  /** `BrowserSourceGroupSet`. */
  ustawZestawZrodel(zadanie: BrowserSourceGroupSetRequest): Promise<Wynik<BrowserSourceGroupSetResponse>>;
  /** `BrowserSourceGroupList`. */
  wykazZestawowZrodel(zadanie: BrowserSourceGroupListRequest): Promise<Wynik<BrowserSourceGroupListResponse>>;
  /** `BrowserNoteThreadSet`. */
  ustawWatekNotatek(zadanie: BrowserNoteThreadSetRequest): Promise<Wynik<BrowserNoteThreadSetResponse>>;
  /** `BrowserNoteThreadList`. */
  wykazWatkowNotatek(zadanie: BrowserNoteThreadListRequest): Promise<Wynik<BrowserNoteThreadListResponse>>;
  /** `BrowserDomInspect`. */
  zbadajDrzewo(zadanie: BrowserDomInspectRequest): Promise<Wynik<BrowserDomInspectResponse>>;
  /** `BrowserNetworkHar`. */
  rejestrSieciowy(zadanie: BrowserNetworkHarRequest): Promise<Wynik<BrowserNetworkHarResponse>>;
  /** `BrowserConsoleRead`. */
  odczytajKonsole(zadanie: BrowserConsoleReadRequest): Promise<Wynik<BrowserConsoleReadResponse>>;
  /** `BrowserDeviceEmulate`. */
  emulujUrzadzenie(zadanie: BrowserDeviceEmulateRequest): Promise<Wynik<BrowserDeviceEmulateResponse>>;
  /** `BrowserScroll`. */
  przewin(zadanie: BrowserScrollRequest): Promise<Wynik<BrowserScrollResponse>>;
  /** `BrowserScreenshotCapture`. */
  wykonajZrzut(zadanie: BrowserScreenshotCaptureRequest): Promise<Wynik<BrowserScreenshotCaptureResponse>>;
  /** `BrowserSnapshotScreenshotGet`. */
  odczytajZrzut(zadanie: BrowserSnapshotScreenshotGetRequest): Promise<Wynik<BrowserSnapshotScreenshotGetResponse>>;
  /** `BrowserArtifactAdd`. */
  dodajWytwor(zadanie: BrowserArtifactAddRequest): Promise<Wynik<BrowserArtifactAddResponse>>;
  /** `BrowserDownloadList`. */
  wykazPobran(zadanie: BrowserDownloadListRequest): Promise<Wynik<BrowserDownloadListResponse>>;
  /** `BrowserDownloadControl`. */
  sterujPobraniem(zadanie: BrowserDownloadControlRequest): Promise<Wynik<BrowserDownloadControlResponse>>;
  /** `BrowserMacroRecord`. */
  nagrywajMakro(zadanie: BrowserMacroRecordRequest): Promise<Wynik<BrowserMacroRecordResponse>>;
  /** `BrowserExecutorLimitsSet`. */
  ustawGraniceWykonawcy(zadanie: BrowserExecutorLimitsSetRequest): Promise<Wynik<BrowserExecutorLimitsSetResponse>>;
  /** `BrowserExecutorLimitsGet`. */
  odczytajGraniceWykonawcy(zadanie: BrowserExecutorLimitsGetRequest): Promise<Wynik<BrowserExecutorLimitsGetResponse>>;
  /** Subskrypcja `browser.page.changed` — jednego ze zdarzeń obszaru. */
  naZmianeStrony(sluchacz: (tresc: BrowserPageChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloBrowser(kanal: Kanal): ZrodloBrowser {
  return {
    async przejdz(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserNavigate, zadanie),
        Command.BrowserNavigate,
        (tresc) => czyObiekt(tresc.snapshot),
      );
    },

    async migawka(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserSnapshotGet, zadanie),
        Command.BrowserSnapshotGet,
        (tresc) => czyObiekt(tresc.snapshot),
      );
    },

    async dodajZrodlo(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserSourceAdd, zadanie),
        Command.BrowserSourceAdd,
        (tresc) => czyObiekt(tresc.source),
      );
    },

    async dodajNotatke(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserNoteAdd, zadanie),
        Command.BrowserNoteAdd,
        (tresc) => czyObiekt(tresc.note),
      );
    },

    async wykazZrodel(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserSourceList, zadanie),
        Command.BrowserSourceList,
        (tresc) => czyTablica(tresc.sources),
      );
    },

    async wykazNotatek(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserNoteList, zadanie),
        Command.BrowserNoteList,
        (tresc) => czyTablica(tresc.notes),
      );
    },


    async otworzKarte(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserTabOpen, zadanie),
        Command.BrowserTabOpen,
        (tresc) => czyObiekt(tresc.tab),
      );
    },

    async wykazKart(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserTabList, zadanie),
        Command.BrowserTabList,
        (tresc) => czyTablica(tresc.tabs),
      );
    },

    async zmienKarte(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserTabUpdate, zadanie),
        Command.BrowserTabUpdate,
        (tresc) => czyObiekt(tresc.tab),
      );
    },

    async zamknijKarte(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserTabClose, zadanie),
        Command.BrowserTabClose,
        () => true,
      );
    },

    async ustawGrupeKart(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserTabGroupSet, zadanie),
        Command.BrowserTabGroupSet,
        (tresc) => czyObiekt(tresc.group),
      );
    },

    async zapiszPrzestrzen(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserWorkspaceSave, zadanie),
        Command.BrowserWorkspaceSave,
        (tresc) => czyObiekt(tresc.workspace),
      );
    },

    async wykazPrzestrzeni(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserWorkspaceList, zadanie),
        Command.BrowserWorkspaceList,
        (tresc) => czyTablica(tresc.workspaces),
      );
    },

    async otworzPrzestrzen(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserWorkspaceOpen, zadanie),
        Command.BrowserWorkspaceOpen,
        (tresc) => czyObiekt(tresc.workspace),
      );
    },

    async usunPrzestrzen(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserWorkspaceRemove, zadanie),
        Command.BrowserWorkspaceRemove,
        () => true,
      );
    },

    async zalozMonitor(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserMonitorAdd, zadanie),
        Command.BrowserMonitorAdd,
        (tresc) => czyObiekt(tresc.monitor),
      );
    },

    async wykazMonitorow(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserMonitorList, zadanie),
        Command.BrowserMonitorList,
        (tresc) => czyTablica(tresc.monitors),
      );
    },

    async sprawdzMonitor(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserMonitorCheck, zadanie),
        Command.BrowserMonitorCheck,
        (tresc) => czyObiekt(tresc.monitor),
      );
    },

    async zdejmijMonitor(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserMonitorRemove, zadanie),
        Command.BrowserMonitorRemove,
        () => true,
      );
    },

    async subskrybujKanal(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserFeedSubscribe, zadanie),
        Command.BrowserFeedSubscribe,
        (tresc) => czyObiekt(tresc.feed),
      );
    },

    async wykazKanalow(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserFeedList, zadanie),
        Command.BrowserFeedList,
        (tresc) => czyTablica(tresc.feeds),
      );
    },

    async zdejmijKanal(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserFeedRemove, zadanie),
        Command.BrowserFeedRemove,
        () => true,
      );
    },

    async odlozDoCzytania(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserReadlistAdd, zadanie),
        Command.BrowserReadlistAdd,
        (tresc) => czyObiekt(tresc.item),
      );
    },

    async wykazCzytania(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserReadlistList, zadanie),
        Command.BrowserReadlistList,
        (tresc) => czyTablica(tresc.items),
      );
    },

    async zdejmijZCzytania(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserReadlistRemove, zadanie),
        Command.BrowserReadlistRemove,
        () => true,
      );
    },

    async dodajZakladke(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserBookmarkAdd, zadanie),
        Command.BrowserBookmarkAdd,
        (tresc) => czyObiekt(tresc.bookmark),
      );
    },

    async wykazZakladek(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserBookmarkList, zadanie),
        Command.BrowserBookmarkList,
        (tresc) => czyTablica(tresc.bookmarks),
      );
    },

    async usunZakladke(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserBookmarkRemove, zadanie),
        Command.BrowserBookmarkRemove,
        () => true,
      );
    },

    async usunZrodlo(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserSourceRemove, zadanie),
        Command.BrowserSourceRemove,
        () => true,
      );
    },

    async zmienNotatke(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserNoteUpdate, zadanie),
        Command.BrowserNoteUpdate,
        (tresc) => czyObiekt(tresc.note),
      );
    },

    async ustawZestawZrodel(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserSourceGroupSet, zadanie),
        Command.BrowserSourceGroupSet,
        (tresc) => czyObiekt(tresc.group),
      );
    },

    async wykazZestawowZrodel(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserSourceGroupList, zadanie),
        Command.BrowserSourceGroupList,
        (tresc) => czyTablica(tresc.groups),
      );
    },

    async ustawWatekNotatek(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserNoteThreadSet, zadanie),
        Command.BrowserNoteThreadSet,
        (tresc) => czyObiekt(tresc.thread),
      );
    },

    async wykazWatkowNotatek(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserNoteThreadList, zadanie),
        Command.BrowserNoteThreadList,
        (tresc) => czyTablica(tresc.threads),
      );
    },

    async zbadajDrzewo(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserDomInspect, zadanie),
        Command.BrowserDomInspect,
        (tresc) => czyTablica(tresc.nodes),
      );
    },

    async rejestrSieciowy(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserNetworkHar, zadanie),
        Command.BrowserNetworkHar,
        (tresc) => czyTablica(tresc.entries),
      );
    },

    async odczytajKonsole(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserConsoleRead, zadanie),
        Command.BrowserConsoleRead,
        (tresc) => czyTablica(tresc.entries),
      );
    },

    async emulujUrzadzenie(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserDeviceEmulate, zadanie),
        Command.BrowserDeviceEmulate,
        (tresc) => czyObiekt(tresc.metrics),
      );
    },

    async przewin(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserScroll, zadanie),
        Command.BrowserScroll,
        (tresc) => czyObiekt(tresc.snapshot),
      );
    },

    async wykonajZrzut(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserScreenshotCapture, zadanie),
        Command.BrowserScreenshotCapture,
        (tresc) => czyObiekt(tresc.screenshot),
      );
    },

    async odczytajZrzut(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserSnapshotScreenshotGet, zadanie),
        Command.BrowserSnapshotScreenshotGet,
        (tresc) => czyObiekt(tresc.screenshot),
      );
    },

    async dodajWytwor(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserArtifactAdd, zadanie),
        Command.BrowserArtifactAdd,
        (tresc) => czyObiekt(tresc.artifact),
      );
    },

    async wykazPobran(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserDownloadList, zadanie),
        Command.BrowserDownloadList,
        (tresc) => czyTablica(tresc.downloads),
      );
    },

    async sterujPobraniem(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserDownloadControl, zadanie),
        Command.BrowserDownloadControl,
        (tresc) => czyObiekt(tresc.download),
      );
    },

    async nagrywajMakro(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserMacroRecord, zadanie),
        Command.BrowserMacroRecord,
        (tresc) => czyObiekt(tresc.macro),
      );
    },

    async ustawGraniceWykonawcy(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserExecutorLimitsSet, zadanie),
        Command.BrowserExecutorLimitsSet,
        (tresc) => czyObiekt(tresc.limits),
      );
    },

    async odczytajGraniceWykonawcy(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.BrowserExecutorLimitsGet, zadanie),
        Command.BrowserExecutorLimitsGet,
        (tresc) => czyObiekt(tresc.limits),
      );
    },

    naZmianeStrony(sluchacz) {
      return kanal.naZdarzenie(EventType.BrowserPageChanged, (tresc) => sluchacz(tresc));
    },
  };
}
