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
 * Plik udostępnia klientowi komendy obszaru workspace wywoływane dziś przez okna modułu wraz z czynnościami sąsiednimi przeniesionymi z innych modułów, a każda czynność oddaje wynik odróżniający brak danych od nieudanego zapytania.
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
  /** Usunięcie wpisu pamięci założonego zapisem instrukcji; inna rodzina komend niż pamięć projektu. */
  usunWpisPamieci(zadanie: MemoryDeleteRequest): Promise<Wynik<MemoryDeleteResponse>>;
  /** Wpisy pamięci widoczne w zasięgu karty sesji przy włączonych poziomach, nie pamięć samego projektu. */
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
