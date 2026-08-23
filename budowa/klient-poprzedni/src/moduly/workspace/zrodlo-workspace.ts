import {
  Command,
  EventType,
  type Agent,
  type ConfigEntry,
  type ConfigScope,
  type LibraryFile,
  type WorkspaceAgentAssignRequest,
  type WorkspaceAgentAssignResponse,
  type WorkspaceContextGetRequest,
  type WorkspaceContextSetRequest,
  type WorkspaceContextSetResponse,
  type WorkspaceDashboard,
  type WorkspaceInstructionsSetRequest,
  type WorkspaceInstructionsSetResponse,
  type WorkspaceLibraryListRequest,
  type MemoryDeleteRequest,
  type MemoryDeleteResponse,
  type MemoryListRequest,
  type MemorySetRequest,
  type MemorySetResponse,
  type MemoryToggleRequest,
  type MemoryToggleResponse,
  type WorkspaceMemoryEntry,
  type WorkspaceProjectChangedEvent,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { przenies } from '../../protokol/wynik-czastkowy';
import { KLUCZ_INSTRUKCJI } from './wynik-czastkowy';
import { czynnosciSasiednie, type CzynnosciSasiednie } from './zrodlo-sasiadow';

/**
 * Moduł Workspace widziany przez klienta — komendy obszaru `workspace.*`, które
 * okna modułu dziś wywołują, wraz z czynnościami sąsiednimi.
 *
 * Obszar liczy w kontrakcie więcej komend, niż stoi w tym pliku: zadania,
 * tablica, harmonogram, kalendarz, notatki, graf wiedzy, oś czasu i komentarze
 * mają nazwy, lecz nie mają jeszcze obsługi w rdzeniu ani okien w tym module.
 * Nieobecność w tym pliku znaczy „niezbudowane”, nie „nieistniejące
 * w kontrakcie” — o pokryciu każdej z nich orzeka `braki-kontraktu.ts`, pytając
 * rdzeń o jego własny wykaz.
 *
 * Każda czynność oddaje `Wynik`, nie samą treść. Okna modułu mają obowiązkowy
 * stan błędu, więc źródło nie połyka niepowodzenia ani nie zwraca pustej listy
 * w jego miejsce — widok ma odróżnić „nic nie ma” od „nie udało się zapytać”.
 * Stąd brak `?? []` w całym pliku.
 *
 * Czynności spoza obszaru — wgranie pliku, wersje, etykieta, uprawnienie
 * eksperta, przeniesienie kontekstu — leżą w `zrodlo-sasiadow.ts` i wchodzą tu
 * rozszerzeniem, bo należą do innych modułów.
 */
export interface ZrodloWorkspace extends CzynnosciSasiednie {
  /** `workspace.dashboard.get` — zestawienie stanu projektu. */
  pulpit(idProjektu: string): Promise<Wynik<WorkspaceDashboard>>;
  /** `workspace.instructions.set` — zapis instrukcji; wynik niesie warstwę obowiązującą. */
  zapiszInstrukcje(
    zadanie: WorkspaceInstructionsSetRequest,
  ): Promise<Wynik<WorkspaceInstructionsSetResponse>>;
  /** `config.get` — zapis instrukcji na jednym poziomie zasięgu, bez rozstrzygania. */
  instrukcjeWarstwy(zasieg: ConfigScope, bytZasiegu: string): Promise<Wynik<ConfigEntry[]>>;
  /** `workspace.context.get` — pamięć projektu. */
  wpisyPamieci(zadanie: WorkspaceContextGetRequest): Promise<Wynik<WorkspaceMemoryEntry[]>>;
  /** `workspace.context.set` — zapis albo zmiana wpisu pamięci. */
  zapiszWpisPamieci(
    zadanie: WorkspaceContextSetRequest,
  ): Promise<Wynik<WorkspaceContextSetResponse>>;
  /**
   * `memory.delete` — usunięcie wpisu pamięci wskazanego identyfikatorem.
   *
   * Osobna rodzina komend niż `workspace.context.*`, ale ten sam byt:
   * `memory.delete` kasuje wpis założony przez `workspace.context.set`, po czym
   * znika on z `context.get`. `memory.detach` nie jest tu wystawiony, bo
   * kontrakt oznacza jego znaczenie jako nieustalone.
   */
  usunWpisPamieci(zadanie: MemoryDeleteRequest): Promise<Wynik<MemoryDeleteResponse>>;
  /**
   * `memory.list` — wpisy pamięci widoczne w zasięgu karty sesji.
   *
   * To nie jest drugie źródło prawdy dla `workspace.context.get`. Tamta komenda
   * pyta o pamięć projektu; ta pyta o to, co widzi karta sesji przy włączonych
   * poziomach — a poziomy przestawia `memory.toggle`. Bez tej pary poziomów
   * pamięci nie da się przestawić z widoku.
   */
  pamiecSesji(zadanie: MemoryListRequest): Promise<Wynik<WorkspaceMemoryEntry[]>>;
  /** `memory.set` — zapis ustalenia w zasięgu innym niż projekt. */
  zapiszPamiec(zadanie: MemorySetRequest): Promise<Wynik<MemorySetResponse>>;
  /** `memory.toggle` — poziomy pamięci włączone dla karty sesji i zapis pamięci. */
  przestawPamiec(zadanie: MemoryToggleRequest): Promise<Wynik<MemoryToggleResponse>>;
  /** `workspace.library.list` — pliki biblioteki projektu. */
  biblioteka(zadanie: WorkspaceLibraryListRequest): Promise<Wynik<LibraryFile[]>>;
  /** `workspace.agent.assign` — przypisanie eksperta do projektu. */
  przypiszAgenta(
    zadanie: WorkspaceAgentAssignRequest,
  ): Promise<Wynik<WorkspaceAgentAssignResponse>>;
  /** `agent.list` — biblioteka ekspertów modułu Agents. */
  agenci(): Promise<Wynik<Agent[]>>;
  /** Subskrypcja `workspace.project.changed`. */
  naZmianeProjektu(sluchacz: (tresc: WorkspaceProjectChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloWorkspace(kanal: Kanal): ZrodloWorkspace {
  return {
    ...czynnosciSasiednie(kanal),

    async pulpit(idProjektu) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceDashboardGet, { projectId: idProjektu }),
        Command.WorkspaceDashboardGet,
        (tresc) => czyObiekt(tresc.dashboard),
      );
      return przenies(wynik, (tresc) => tresc.dashboard);
    },

    zapiszInstrukcje(zadanie) {
      return wywolaj(kanal, Command.WorkspaceInstructionsSet, zadanie);
    },

    async instrukcjeWarstwy(zasieg, bytZasiegu) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ConfigGet, {
          key: KLUCZ_INSTRUKCJI,
          scope: zasieg,
          scopeId: bytZasiegu,
        }),
        Command.ConfigGet,
        (tresc) => czyTablica(tresc.entries),
      );
      return przenies(wynik, (tresc) => tresc.entries);
    },

    async wpisyPamieci(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceContextGet, zadanie),
        Command.WorkspaceContextGet,
        (tresc) => czyTablica(tresc.entries),
      );
      return przenies(wynik, (tresc) => tresc.entries);
    },

    zapiszWpisPamieci(zadanie) {
      return wywolaj(kanal, Command.WorkspaceContextSet, zadanie);
    },
    usunWpisPamieci(zadanie) {
      return wywolaj(kanal, Command.MemoryDelete, zadanie);
    },

    async pamiecSesji(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.MemoryList, zadanie),
        Command.MemoryList,
        (tresc) => czyTablica(tresc.entries),
      );
      return przenies(wynik, (tresc) => tresc.entries);
    },

    zapiszPamiec(zadanie) {
      return wywolaj(kanal, Command.MemorySet, zadanie);
    },

    przestawPamiec(zadanie) {
      return wywolaj(kanal, Command.MemoryToggle, zadanie);
    },

    async biblioteka(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceLibraryList, zadanie),
        Command.WorkspaceLibraryList,
        (tresc) => czyTablica(tresc.files),
      );
      return przenies(wynik, (tresc) => tresc.files);
    },

    przypiszAgenta(zadanie) {
      return wywolaj(kanal, Command.WorkspaceAgentAssign, zadanie);
    },

    async agenci() {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentList, {}),
        Command.AgentList,
        (tresc) => czyTablica(tresc.agents),
      );
      return przenies(wynik, (tresc) => tresc.agents);
    },

    naZmianeProjektu(sluchacz) {
      return kanal.naZdarzenie(EventType.WorkspaceProjectChanged, (tresc) => sluchacz(tresc));
    },
  };
}
