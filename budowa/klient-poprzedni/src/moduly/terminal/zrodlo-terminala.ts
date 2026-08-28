import {
  Command,
  EventType,
  type AutomationSchedule,
  type AutomationScheduleSetRequest,
  type AutomationWorkflow,
  type AutomationWorkflowSaveRequest,
  type Envelope,
  type Queue,
  type QueueActionRequest,
  type QueueListRequest,
  type ScheduleGetRequest,
  type StreamChunkEvent,
  type TerminalCommandExecRequest,
  type TerminalOutputReadRequest,
  type TerminalOutputReadResponse,
  type TerminalProcess,
  type TerminalProcessChangedEvent,
  type TerminalProcessKillRequest,
  type TerminalProcessListRequest,
  type TerminalSession,
  type TerminalSessionOpenRequest,
  type TerminalFileReadRequest,
  type TerminalFileReadResponse,
  type TerminalHost,
  type TerminalHostListRequest,
  type TerminalHostRemoveRequest,
  type TerminalHostSaveRequest,
  type TerminalHostSaveResponse,
  type TerminalKeyGenerateRequest,
  type TerminalKeyImportRequest,
  type TerminalKeyRemoveRequest,
  type TerminalKeyRemoveResponse,
  type TerminalProcessSuspendRequest,
  type TerminalProcessSuspendResponse,
  type TerminalScript,
  type TerminalScriptLintRequest,
  type TerminalScriptLintResponse,
  type TerminalScriptListRequest,
  type TerminalScriptRemoveRequest,
  type TerminalScriptSaveRequest,
  type TerminalScriptSaveResponse,
  type TerminalSessionCloseRequest,
  type TerminalSessionCloseResponse,
  type TerminalSessionListRequest,
  type TerminalSshKey,
  type TerminalTunnel,
  type TerminalTunnelCloseRequest,
  type TerminalTunnelListRequest,
  type TerminalTunnelOpenRequest,
  type TerminalWatch,
  type TerminalWatchListRequest,
  type TerminalWatchStartRequest,
  type TerminalWatchStopRequest,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/** Moduł Terminal widziany przez klienta: pięć komend obszaru terminal.*, pięć komend harmonogramu i kolejek oraz dwie subskrypcje zdarzeń na żywo. */
export interface ZrodloTerminala {
  /** `terminal.session.open` — otwarcie karty powłoki. */
  otworzKarte(zadanie: TerminalSessionOpenRequest): Promise<Wynik<TerminalSession>>;
  /** `terminal.command.exec` — uruchomienie polecenia w karcie. */
  wykonaj(zadanie: TerminalCommandExecRequest): Promise<Wynik<TerminalProcess>>;
  /** `terminal.process.list` — procesy rejestru rdzenia. */
  procesy(zadanie: TerminalProcessListRequest): Promise<Wynik<TerminalProcess[]>>;
  /** `terminal.process.kill` — zakończenie procesu. */
  zakoncz(zadanie: TerminalProcessKillRequest): Promise<Wynik<TerminalProcess>>;
  /** `terminal.output.read` — wyjście jednego procesu z rejestru rdzenia. */
  odczytajWyjscie(zadanie: TerminalOutputReadRequest): Promise<Wynik<TerminalOutputReadResponse>>;
  // Odczyt harmonogramów — do pary z zapisem niżej, którym okno zakłada plan.
  harmonogramy(zadanie: ScheduleGetRequest): Promise<Wynik<AutomationSchedule[]>>;
  // Zapis automatyki z zadaniem powłoki objętym harmonogramem — terminal nie ma rodziny harmonogramu.
  zapiszAutomatyke(zadanie: AutomationWorkflowSaveRequest): Promise<Wynik<AutomationWorkflow>>;
  // Cykliczność zapisanej automatyki — zapis jest planem, nie budzikiem odpalającym go samodzielnie.
  ustawHarmonogram(zadanie: AutomationScheduleSetRequest): Promise<Wynik<AutomationSchedule>>;
  // Kolejki silnika pętli obsługujące okno modułu — zawężenie idzie oknem, nie identyfikatorem karty.
  kolejki(zadanie: QueueListRequest): Promise<Wynik<Queue[]>>;
  /** `queue.action` — sześć działań silnika kolejek na kolejce okna. */
  dzialanieKolejki(zadanie: QueueActionRequest): Promise<Wynik<Queue>>;
  /** `terminal.session.close` — zamknięcie karty powłoki w rdzeniu. */
  zamknijKarte(zadanie: TerminalSessionCloseRequest): Promise<Wynik<TerminalSessionCloseResponse>>;
  // Karty powłoki znane rdzeniowi — wykaz sięga dalej niż pamięć widoku po ponownym podłączeniu.
  karty(zadanie: TerminalSessionListRequest): Promise<Wynik<TerminalSession[]>>;
  /** `terminal.file.read` — treść pliku z katalogu roboczego karty. */
  odczytajPlik(zadanie: TerminalFileReadRequest): Promise<Wynik<TerminalFileReadResponse>>;
  /** `terminal.process.suspend` — wstrzymanie albo wznowienie procesu. */
  wstrzymaj(zadanie: TerminalProcessSuspendRequest): Promise<Wynik<TerminalProcessSuspendResponse>>;

  /** `terminal.host.save` — zapis wpisu książki hostów w rdzeniu. */
  zapiszHosta(zadanie: TerminalHostSaveRequest): Promise<Wynik<TerminalHostSaveResponse>>;
  /** `terminal.host.list` — książka hostów Operatora. */
  hosty(zadanie: TerminalHostListRequest): Promise<Wynik<TerminalHost[]>>;
  /** `terminal.host.remove` — zdjęcie wpisu z książki hostów. */
  usunHosta(zadanie: TerminalHostRemoveRequest): Promise<Wynik<boolean>>;

  /** `terminal.script.save` — zapis pozycji biblioteki jako kolejnej wersji. */
  zapiszSkrypt(zadanie: TerminalScriptSaveRequest): Promise<Wynik<TerminalScriptSaveResponse>>;
  /** `terminal.script.list` — pozycje biblioteki skryptów. */
  skrypty(zadanie: TerminalScriptListRequest): Promise<Wynik<TerminalScript[]>>;
  /** `terminal.script.remove` — usunięcie pozycji wraz z jej wersjami. */
  usunSkrypt(zadanie: TerminalScriptRemoveRequest): Promise<Wynik<boolean>>;
  // Analiza statyczna i formatowanie treści — pole analyzerAvailable mówi, czy sprawdzenie się odbyło.
  sprawdzSkrypt(zadanie: TerminalScriptLintRequest): Promise<Wynik<TerminalScriptLintResponse>>;

  /** `terminal.tunnel.open` — założenie przekierowania portu. */
  otworzTunel(zadanie: TerminalTunnelOpenRequest): Promise<Wynik<TerminalTunnel>>;
  /** `terminal.tunnel.list` — przekierowania portów wraz ze stanem. */
  tunele(zadanie: TerminalTunnelListRequest): Promise<Wynik<TerminalTunnel[]>>;
  /** `terminal.tunnel.close` — zamknięcie przekierowania portu. */
  zamknijTunel(zadanie: TerminalTunnelCloseRequest): Promise<Wynik<TerminalTunnel>>;

  /** `terminal.key.generate` — wytworzenie pary kluczy na maszynie rdzenia. */
  wytworzKlucz(zadanie: TerminalKeyGenerateRequest): Promise<Wynik<TerminalSshKey>>;
  /** `terminal.key.import` — wciągnięcie do wykazu klucza leżącego na maszynie rdzenia. */
  wciagnijKlucz(zadanie: TerminalKeyImportRequest): Promise<Wynik<TerminalSshKey>>;
  /** `terminal.key.list` — wykaz kluczy SSH wraz z odciskami. */
  klucze(): Promise<Wynik<TerminalSshKey[]>>;
  /** `terminal.key.remove` — zdjęcie klucza z wykazu. */
  usunKlucz(zadanie: TerminalKeyRemoveRequest): Promise<Wynik<TerminalKeyRemoveResponse>>;

  /** `terminal.watch.start` — obserwacja plików wyzwalająca polecenie karty. */
  zalozObserwacje(zadanie: TerminalWatchStartRequest): Promise<Wynik<TerminalWatch>>;
  /** `terminal.watch.stop` — zatrzymanie obserwacji plików. */
  zatrzymajObserwacje(zadanie: TerminalWatchStopRequest): Promise<Wynik<TerminalWatch>>;
  /** `terminal.watch.list` — obserwacje wraz z licznikiem wyzwoleń. */
  obserwacje(zadanie: TerminalWatchListRequest): Promise<Wynik<TerminalWatch[]>>;

  /** Subskrypcja `terminal.process.changed` — Process Monitor na żywo. */
  naZmianeProcesu(sluchacz: (tresc: TerminalProcessChangedEvent) => void): Odsubskrybuj;
  /** Subskrypcja `stream.chunk` — wyjście procesów do Output Console. */
  naFragmentWyjscia(sluchacz: (tresc: StreamChunkEvent, koperta: Envelope) => void): Odsubskrybuj;
}

export function utworzZrodloTerminala(kanal: Kanal): ZrodloTerminala {
  return {
    async otworzKarte(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalSessionOpen, zadanie),
        Command.TerminalSessionOpen,
        (tresc) => czyObiekt(tresc.session),
      );
      return przenies(wynik, (tresc) => tresc.session);
    },

    async wykonaj(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalCommandExec, zadanie),
        Command.TerminalCommandExec,
        (tresc) => czyObiekt(tresc.process),
      );
      return przenies(wynik, (tresc) => tresc.process);
    },

    async procesy(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalProcessList, zadanie),
        Command.TerminalProcessList,
        (tresc) => czyTablica(tresc.processes),
      );
      return przenies(wynik, (tresc) => tresc.processes);
    },

    async zakoncz(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalProcessKill, zadanie),
        Command.TerminalProcessKill,
        (tresc) => czyObiekt(tresc.process),
      );
      return przenies(wynik, (tresc) => tresc.process);
    },

    // Odpowiedź jest całą treścią, nie opakowaniem — oba strumienie są sprawdzane, bo są obowiązkowe.
    async odczytajWyjscie(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalOutputRead, zadanie),
        Command.TerminalOutputRead,
        (tresc) => czyTekst(tresc.stdout) && czyTekst(tresc.stderr),
      );
    },

    async harmonogramy(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ScheduleGet, zadanie),
        Command.ScheduleGet,
        (tresc) => czyTablica(tresc.schedules),
      );
      return przenies(wynik, (tresc) => tresc.schedules);
    },

    async zapiszAutomatyke(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationWorkflowSave, zadanie),
        Command.AutomationWorkflowSave,
        (tresc) => czyObiekt(tresc.workflow),
      );
      return przenies(wynik, (tresc) => tresc.workflow);
    },

    async ustawHarmonogram(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationScheduleSet, zadanie),
        Command.AutomationScheduleSet,
        (tresc) => czyObiekt(tresc.schedule),
      );
      return przenies(wynik, (tresc) => tresc.schedule);
    },

    async kolejki(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.QueueList, zadanie),
        Command.QueueList,
        (tresc) => czyTablica(tresc.queues),
      );
      return przenies(wynik, (tresc) => tresc.queues);
    },

    async dzialanieKolejki(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.QueueAction, zadanie),
        Command.QueueAction,
        (tresc) => czyObiekt(tresc.queue),
      );
      return przenies(wynik, (tresc) => tresc.queue);
    },

    async zamknijKarte(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalSessionClose, zadanie),
        Command.TerminalSessionClose,
        (tresc) => czyObiekt(tresc.session),
      );
    },

    async karty(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalSessionList, zadanie),
        Command.TerminalSessionList,
        (tresc) => czyTablica(tresc.sessions),
      );
      return przenies(wynik, (tresc) => tresc.sessions);
    },

    async odczytajPlik(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalFileRead, zadanie),
        Command.TerminalFileRead,
        (tresc) => czyTekst(tresc.content) && czyTekst(tresc.path),
      );
    },

    async wstrzymaj(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalProcessSuspend, zadanie),
        Command.TerminalProcessSuspend,
        (tresc) => czyObiekt(tresc.process),
      );
    },

    async zapiszHosta(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalHostSave, zadanie),
        Command.TerminalHostSave,
        (tresc) => czyObiekt(tresc.host),
      );
    },

    async hosty(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalHostList, zadanie),
        Command.TerminalHostList,
        (tresc) => czyTablica(tresc.hosts),
      );
      return przenies(wynik, (tresc) => tresc.hosts);
    },

    // Pole removed jest obowiązkowe i niesie prawdę o skutku — fałsz nie jest błędem.
    async usunHosta(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalHostRemove, zadanie),
        Command.TerminalHostRemove,
        (tresc) => typeof tresc.removed === 'boolean',
      );
      return przenies(wynik, (tresc) => tresc.removed);
    },

    async zapiszSkrypt(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalScriptSave, zadanie),
        Command.TerminalScriptSave,
        (tresc) => czyObiekt(tresc.script),
      );
    },

    async skrypty(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalScriptList, zadanie),
        Command.TerminalScriptList,
        (tresc) => czyTablica(tresc.scripts),
      );
      return przenies(wynik, (tresc) => tresc.scripts);
    },

    async usunSkrypt(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalScriptRemove, zadanie),
        Command.TerminalScriptRemove,
        (tresc) => typeof tresc.removed === 'boolean',
      );
      return przenies(wynik, (tresc) => tresc.removed);
    },

    async sprawdzSkrypt(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalScriptLint, zadanie),
        Command.TerminalScriptLint,
        (tresc) => czyTablica(tresc.findings) && typeof tresc.analyzerAvailable === 'boolean',
      );
    },

    async otworzTunel(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalTunnelOpen, zadanie),
        Command.TerminalTunnelOpen,
        (tresc) => czyObiekt(tresc.tunnel),
      );
      return przenies(wynik, (tresc) => tresc.tunnel);
    },

    async tunele(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalTunnelList, zadanie),
        Command.TerminalTunnelList,
        (tresc) => czyTablica(tresc.tunnels),
      );
      return przenies(wynik, (tresc) => tresc.tunnels);
    },

    async zamknijTunel(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalTunnelClose, zadanie),
        Command.TerminalTunnelClose,
        (tresc) => czyObiekt(tresc.tunnel),
      );
      return przenies(wynik, (tresc) => tresc.tunnel);
    },

    async wytworzKlucz(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalKeyGenerate, zadanie),
        Command.TerminalKeyGenerate,
        (tresc) => czyObiekt(tresc.key),
      );
      return przenies(wynik, (tresc) => tresc.key);
    },

    async wciagnijKlucz(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalKeyImport, zadanie),
        Command.TerminalKeyImport,
        (tresc) => czyObiekt(tresc.key),
      );
      return przenies(wynik, (tresc) => tresc.key);
    },

    async klucze() {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalKeyList, {}),
        Command.TerminalKeyList,
        (tresc) => czyTablica(tresc.keys),
      );
      return przenies(wynik, (tresc) => tresc.keys);
    },

    async usunKlucz(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalKeyRemove, zadanie),
        Command.TerminalKeyRemove,
        (tresc) => typeof tresc.removed === 'boolean',
      );
    },

    async zalozObserwacje(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalWatchStart, zadanie),
        Command.TerminalWatchStart,
        (tresc) => czyObiekt(tresc.watch),
      );
      return przenies(wynik, (tresc) => tresc.watch);
    },

    async zatrzymajObserwacje(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalWatchStop, zadanie),
        Command.TerminalWatchStop,
        (tresc) => czyObiekt(tresc.watch),
      );
      return przenies(wynik, (tresc) => tresc.watch);
    },

    async obserwacje(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.TerminalWatchList, zadanie),
        Command.TerminalWatchList,
        (tresc) => czyTablica(tresc.watches),
      );
      return przenies(wynik, (tresc) => tresc.watches);
    },

    naZmianeProcesu(sluchacz) {
      return kanal.naZdarzenie(EventType.TerminalProcessChanged, (tresc) => sluchacz(tresc));
    },

    naFragmentWyjscia(sluchacz) {
      return kanal.naZdarzenie(EventType.StreamChunk, (tresc, koperta) => sluchacz(tresc, koperta));
    },
  };
}
