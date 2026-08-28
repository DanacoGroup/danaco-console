import {
  Command,
  type Agent,
  type AgentAssignment,
  type AgentAssignmentListRequest,
  type AgentConnector,
  type AgentConnectorConfigureRequest,
  type AgentPermission,
  type AgentPermissionGroup,
  type AgentPolicy,
  type AgentVersion,
  type AgentVersionSnapshot,
  type IsolationTechnicalSwitch,
  type ToolScope,
  type ToolsScopeSetRequest,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Komendy zakresu działania eksperta widziane przez okna modułu Agents opisują, co ekspertowi
 * wolno zrobić w systemie, osobno od biblioteki założeń i tożsamości.
 */
export interface ZrodloZakresuEksperta {
  /** Zdejmuje umiejętność z definicji eksperta (`agent.skill.remove`). */
  odlaczUmiejetnosc(idEksperta: string, idUmiejetnosci: string): Promise<Wynik<{ agent: Agent }>>;
  /** Konektory eksperta wraz z definicją (`agent.connector.list`). */
  konektory(idEksperta: string): Promise<Wynik<{ connectors: AgentConnector[] }>>;
  /** Odłącza konektor od eksperta (`agent.connector.remove`). */
  odlaczKonektor(idEksperta: string, idKonektora: string): Promise<Wynik<{ agent: Agent }>>;
  /** Zapisuje konfigurację instancji konektora (`agent.connector.configure`). */
  skonfigurujKonektor(
    zlecenie: ZlecenieKonfiguracjiKonektora,
  ): Promise<Wynik<{ connector: AgentConnector }>>;
  /** Tożsamość utrwalona w jednej wersji (`agent.version.get`). */
  wersja(
    idEksperta: string,
    idWersji: string,
  ): Promise<Wynik<{ version: AgentVersion; snapshot: AgentVersionSnapshot }>>;
  /** Przypisania eksperta do projektów i ról (`agent.assignment.list`). */
  przypisania(
    idEksperta?: string,
  ): Promise<Wynik<{ assignments: AgentAssignment[]; total: number }>>;
  /** Moduły zastosowania eksperta (`agent.modules.set`); pusty wykaz zdejmuje ograniczenie. */
  ustawModuly(idEksperta: string, kodyModulow: readonly string[]): Promise<Wynik<{ agent: Agent }>>;
  /** Osiem zakresów izolacji technicznej eksperta (`agent.isolation.get`). */
  izolacja(idEksperta: string): Promise<Wynik<{ switches: IsolationTechnicalSwitch[] }>>;
  /** Zapisuje wskazane zakresy izolacji (`agent.isolation.set`). */
  zapiszIzolacje(
    idEksperta: string,
    przelaczniki: readonly IsolationTechnicalSwitch[],
  ): Promise<Wynik<{ switches: IsolationTechnicalSwitch[] }>>;
  /** Polityka efektywna eksperta (`agent.policy.get`). */
  polityka(idEksperta: string, idOkna?: string): Promise<Wynik<{ policy: AgentPolicy }>>;
  /** Subagent Network eksperta i jego granica (`agent.subagent.set`). */
  ustawPodagentow(
    idEksperta: string,
    czynny: boolean,
    granica?: number,
  ): Promise<Wynik<{ agent: Agent }>>;
  /** Zdejmuje wpisy uprawnień (`agent.permission.remove`); grupa pominięta zdejmuje wszystkie. */
  zdejmijUprawnienie(
    zlecenie: ZlecenieZdjeciaUprawnienia,
  ): Promise<Wynik<{ permissions: AgentPermission[]; removed: number }>>;
  /** Zakresy narzędzi profilu asystenta (`tools.scope.list`). */
  zakresyNarzedzi(
    idProfilu?: string,
  ): Promise<Wynik<{ scopes: ToolScope[]; total: number }>>;
  /** Zapisuje zakres i limit wywołań pozycji katalogu (`tools.scope.set`). */
  zapiszZakresNarzedzia(zlecenie: ZlecenieZakresuNarzedzia): Promise<Wynik<{ scope: ToolScope }>>;
}

/**
 * Zlecenie konfiguracji instancji konektora niesie pole pominięte jako brak zmiany
 * dotychczasowej wartości zapisanej w rdzeniu dla tego eksperta.
 */
export interface ZlecenieKonfiguracjiKonektora {
  idEksperta: string;
  idKonektora: string;
  punktDostepu?: string;
  konfiguracja?: string;
  czynny?: boolean;
}

/**
 * Zlecenie zdjęcia wpisu uprawnienia z pominiętą grupą zdejmuje wszystkie grupy zakresu tego
 * eksperta naraz.
 */
export interface ZlecenieZdjeciaUprawnienia {
  idEksperta: string;
  grupa?: AgentPermissionGroup;
  zakres?: string;
}

/**
 * Zlecenie zapisu zakresu narzędzia niesie pole pominięte jako brak zmiany dotychczasowej
 * wartości zapisanej dla profilu.
 */
export interface ZlecenieZakresuNarzedzia {
  idProfilu: string;
  nazwaPozycji: string;
  dostepna?: boolean;
  potwierdzenie?: boolean;
  granicaWywolan?: number;
  oknoSekund?: number;
  uzasadnienie?: string;
}

export function utworzZrodloZakresuEksperta(kanal: Kanal): ZrodloZakresuEksperta {
  return {
    async odlaczUmiejetnosc(idEksperta, idUmiejetnosci) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentSkillRemove, {
          agentId: idEksperta,
          skillId: idUmiejetnosci.trim(),
        }),
        Command.AgentSkillRemove,
        (tresc) => czyObiekt(tresc.agent),
      );
    },

    async konektory(idEksperta) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentConnectorList, { agentId: idEksperta }),
        Command.AgentConnectorList,
        (tresc) => czyTablica(tresc.connectors),
      );
    },

    async odlaczKonektor(idEksperta, idKonektora) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentConnectorRemove, {
          agentId: idEksperta,
          connectorId: idKonektora,
        }),
        Command.AgentConnectorRemove,
        (tresc) => czyObiekt(tresc.agent),
      );
    },

    async skonfigurujKonektor(zlecenie) {
      const zadanie: AgentConnectorConfigureRequest = {
        agentId: zlecenie.idEksperta,
        connectorId: zlecenie.idKonektora,
      };
      if (zlecenie.punktDostepu !== undefined) zadanie.accessPointId = zlecenie.punktDostepu.trim();
      if (zlecenie.czynny !== undefined) zadanie.enabled = zlecenie.czynny;
      // Konfiguracja jedzie jako treść JSON; zapis niepoprawny jest odmową wywołania, nie wyjątkiem.
      if (zlecenie.konfiguracja !== undefined && zlecenie.konfiguracja.trim() !== '') {
        try {
          zadanie.config = JSON.parse(zlecenie.konfiguracja) as unknown;
        } catch (powod) {
          return {
            udany: false,
            blad: {
              code: 'validation_failed',
              message: `konfiguracja integracji nie jest poprawnym JSON: ${String(powod)}`,
              retryable: false,
            },
          };
        }
      }
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentConnectorConfigure, zadanie),
        Command.AgentConnectorConfigure,
        (tresc) => czyObiekt(tresc.connector),
      );
    },

    async wersja(idEksperta, idWersji) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentVersionGet, {
          agentId: idEksperta,
          versionId: idWersji,
        }),
        Command.AgentVersionGet,
        (tresc) => czyObiekt(tresc.snapshot),
      );
    },

    async przypisania(idEksperta) {
      const zadanie: AgentAssignmentListRequest = {};
      if (idEksperta !== undefined && idEksperta !== '') zadanie.agentId = idEksperta;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentAssignmentList, zadanie),
        Command.AgentAssignmentList,
        (tresc) => czyTablica(tresc.assignments),
      );
    },

    async ustawModuly(idEksperta, kodyModulow) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentModulesSet, {
          agentId: idEksperta,
          // Wykaz jedzie zawsze, także pusty: pusty znaczy brak ograniczenia, nie brak żądania.
          moduleCodes: [...kodyModulow],
        }),
        Command.AgentModulesSet,
        (tresc) => czyObiekt(tresc.agent),
      );
    },

    async izolacja(idEksperta) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentIsolationGet, { agentId: idEksperta }),
        Command.AgentIsolationGet,
        (tresc) => czyTablica(tresc.switches),
      );
    },

    async zapiszIzolacje(idEksperta, przelaczniki) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentIsolationSet, {
          agentId: idEksperta,
          switches: [...przelaczniki],
        }),
        Command.AgentIsolationSet,
        (tresc) => czyTablica(tresc.switches),
      );
    },

    async polityka(idEksperta, idOkna) {
      const zadanie: { agentId: string; windowId?: string } = { agentId: idEksperta };
      if (idOkna !== undefined && idOkna !== '') zadanie.windowId = idOkna;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentPolicyGet, zadanie),
        Command.AgentPolicyGet,
        (tresc) => czyObiekt(tresc.policy),
      );
    },

    async ustawPodagentow(idEksperta, czynny, granica) {
      const zadanie: { agentId: string; enabled: boolean; limit?: number } = {
        agentId: idEksperta,
        enabled: czynny,
      };
      if (granica !== undefined) zadanie.limit = granica;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentSubagentSet, zadanie),
        Command.AgentSubagentSet,
        (tresc) => czyObiekt(tresc.agent),
      );
    },

    async zdejmijUprawnienie(zlecenie) {
      const zadanie: { agentId: string; group?: AgentPermissionGroup; scope?: string } = {
        agentId: zlecenie.idEksperta,
      };
      if (zlecenie.grupa !== undefined) zadanie.group = zlecenie.grupa;
      if (zlecenie.zakres !== undefined && zlecenie.zakres !== '') zadanie.scope = zlecenie.zakres;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentPermissionRemove, zadanie),
        Command.AgentPermissionRemove,
        (tresc) => czyTablica(tresc.permissions) && typeof tresc.removed === 'number',
      );
    },

    async zakresyNarzedzi(idProfilu) {
      const zadanie: { profileId?: string } = {};
      if (idProfilu !== undefined && idProfilu !== '') zadanie.profileId = idProfilu;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ToolsScopeList, zadanie),
        Command.ToolsScopeList,
        (tresc) => czyTablica(tresc.scopes),
      );
    },

    async zapiszZakresNarzedzia(zlecenie) {
      const zadanie: ToolsScopeSetRequest = {
        profileId: zlecenie.idProfilu.trim(),
        toolName: zlecenie.nazwaPozycji.trim(),
      };
      if (zlecenie.dostepna !== undefined) zadanie.enabled = zlecenie.dostepna;
      if (zlecenie.potwierdzenie !== undefined) zadanie.confirmRequired = zlecenie.potwierdzenie;
      if (zlecenie.granicaWywolan !== undefined) zadanie.callLimit = zlecenie.granicaWywolan;
      if (zlecenie.oknoSekund !== undefined) zadanie.callWindowSeconds = zlecenie.oknoSekund;
      if (zlecenie.uzasadnienie !== undefined && zlecenie.uzasadnienie !== '') {
        zadanie.note = zlecenie.uzasadnienie;
      }
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ToolsScopeSet, zadanie),
        Command.ToolsScopeSet,
        (tresc) => czyObiekt(tresc.scope),
      );
    },
  };
}
