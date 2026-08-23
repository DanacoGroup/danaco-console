import {
  Command,
  ErrorCode,
  EventType,
  type Agent,
  type AgentChangedEvent,
  type AgentConnector,
  type AgentConnectorAddRequest,
  type AgentCreateRequest,
  type AgentLayerRemoveRequest,
  type AgentLayerSetRequest,
  type AgentListRequest,
  type AgentModelSetRequest,
  type AgentPermission,
  type AgentPermissionSetRequest,
  type AgentPlugin,
  type AgentPluginAddRequest,
  type AgentUpdateRequest,
  IdentityMode,
  type ErrorInfo,
  type IdentityLayer,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import type {
  TozsamoscEksperta,
  ZlecenieKonektora,
  ZlecenieModelu,
  ZlecenieUprawnienia,
  ZlecenieWarstwy,
  ZlecenieWtyczki,
  ZmianaEksperta,
} from './zlecenia-agentow';

/**
 * Komendy obszaru `agent.*` widziane przez okna modułu Agents.
 *
 * Źródło nie ma własnego stanu i niczego nie pamięta — jest wyłącznie warstwą
 * wywołań i sprawdzianu kształtu odpowiedzi. Stan biblioteki ekspertów mieszka
 * w `stan-agentow.ts`, żeby pięć okien modułu patrzyło na jeden zbiór, a nie na
 * pięć osobnych kopii.
 *
 * Żadne wywołanie nie rzuca wyjątkiem ani nie odrzuca obietnicy: niepowodzenie
 * wraca polem `blad` wyniku, a okno pokazuje je w swoim stanie błędu. Dotyczy to
 * także treści JSON wpisanej przez Operatora — niepoprawny zapis jest odmową
 * wywołania, nie wyjątkiem wywracającym widok.
 */
export interface ZrodloAgentow {
  wykaz(fraza: string, tylkoCzynne: boolean): Promise<Wynik<{ agents: Agent[]; total?: number }>>;
  utworz(tozsamosc: TozsamoscEksperta): Promise<Wynik<{ agent: Agent }>>;
  zmien(idEksperta: string, zmiana: ZmianaEksperta): Promise<Wynik<{ agent: Agent }>>;
  usun(idEksperta: string): Promise<Wynik<{ agentId: string; deleted: boolean }>>;
  ustawModel(zlecenie: ZlecenieModelu): Promise<Wynik<{ agent: Agent }>>;
  dodajUmiejetnosc(idEksperta: string, idUmiejetnosci: string): Promise<Wynik<{ agent: Agent }>>;
  dodajKonektor(zlecenie: ZlecenieKonektora): Promise<Wynik<{ connector: AgentConnector }>>;
  ustawUprawnienie(
    zlecenie: ZlecenieUprawnienia,
  ): Promise<Wynik<{ permissions: AgentPermission[] }>>;
  /** Zapisuje treść jednej warstwy promptu i jej stan czynności. */
  zapiszWarstwe(zlecenie: ZlecenieWarstwy): Promise<Wynik<{ agent: Agent }>>;
  /** Usuwa treść warstwy; brak warstwy znaczy prompt bez niej, nie błąd. */
  usunWarstwe(idEksperta: string, warstwa: IdentityLayer): Promise<Wynik<{ agent: Agent }>>;
  dodajWtyczke(zlecenie: ZlecenieWtyczki): Promise<Wynik<{ plugin: AgentPlugin }>>;
  odlaczWtyczke(idEksperta: string, idWtyczki: string): Promise<Wynik<{ agent: Agent }>>;
  /**
   * Wtyczki eksperta wraz z definicją (`agent.plugin.list`).
   *
   * `Agent.pluginIds` niesie same identyfikatory, więc bez tej czynności wykaz
   * nie zna ani nazwy nadanej przez Operatora, ani wersji, ani źródła.
   */
  wtyczki(idEksperta: string): Promise<Wynik<{ plugins: AgentPlugin[] }>>;
  /** Subskrypcja zdarzenia `agent.changed` — jedynego zdarzenia obszaru. */
  naZmiane(sluchacz: (tresc: AgentChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloAgentow(kanal: Kanal): ZrodloAgentow {
  return {
    async wykaz(fraza, tylkoCzynne) {
      const zadanie: AgentListRequest = { enabledOnly: tylkoCzynne };
      if (fraza.trim() !== '') zadanie.query = fraza.trim();
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentList, zadanie),
        Command.AgentList,
        (tresc) => czyTablica(tresc.agents),
      );
    },

    async utworz(tozsamosc) {
      const zadanie: AgentCreateRequest = { name: tozsamosc.nazwa.trim() };
      if (tozsamosc.opis.trim() !== '') zadanie.description = tozsamosc.opis.trim();
      if (tozsamosc.instrukcje !== '') zadanie.systemPrompt = tozsamosc.instrukcje;
      if (tozsamosc.kanal !== '') zadanie.channelId = tozsamosc.kanal;
      if (tozsamosc.model.trim() !== '') zadanie.model = tozsamosc.model.trim();
      if (tozsamosc.imie.trim() !== '') zadanie.displayName = tozsamosc.imie.trim();
      if (tozsamosc.favikon.trim() !== '') zadanie.favicon = tozsamosc.favikon.trim();
      // Tryb idzie wyłącznie przy odstępstwie. Pominięcie pola znaczy w rdzeniu
      // `DOLACZ` — stan domyślny — więc wysyłanie go zawsze nie zmieniałoby
      // niczego poza tym, że zakładanie eksperta wyglądałoby na wybór między
      // dwiema równorzędnymi wartościami.
      if (tozsamosc.zastepuje) zadanie.mode = IdentityMode.ZASTAP;
      zadanie.visibility = tozsamosc.widocznosc;
      // Poziomy pamięci jadą zawsze, także puste: pusty zbiór jest w kontrakcie
      // jedynym zapisem wyłączenia pamięci, a pominięcie pola znaczy co innego
      // — „bez zmiany”.
      zadanie.memoryLevels = [...tozsamosc.poziomyPamieci];
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentCreate, zadanie),
        Command.AgentCreate,
        (tresc) => czyObiekt(tresc.agent),
      );
    },

    async zmien(idEksperta, zmiana) {
      const zadanie: AgentUpdateRequest = { agentId: idEksperta };
      if (zmiana.nazwa !== undefined) zadanie.name = zmiana.nazwa;
      if (zmiana.opis !== undefined) zadanie.description = zmiana.opis;
      if (zmiana.instrukcje !== undefined) zadanie.systemPrompt = zmiana.instrukcje;
      if (zmiana.czynny !== undefined) zadanie.enabled = zmiana.czynny;
      if (zmiana.imie !== undefined) zadanie.displayName = zmiana.imie;
      if (zmiana.favikon !== undefined) zadanie.favicon = zmiana.favikon;
      // Przy zmianie tryb idzie w obie strony, inaczej niż przy zakładaniu:
      // cofnięcie odstępstwa to podanie `DOLACZ` wprost. Pominięcie pola
      // zostawiłoby w rdzeniu wartość starą i odznaczenie kontrolki nie doszłoby
      // do skutku.
      if (zmiana.zastepuje !== undefined) {
        zadanie.mode = zmiana.zastepuje ? IdentityMode.ZASTAP : IdentityMode.DOLACZ;
      }
      if (zmiana.widocznosc !== undefined) zadanie.visibility = zmiana.widocznosc;
      if (zmiana.poziomyPamieci !== undefined) zadanie.memoryLevels = [...zmiana.poziomyPamieci];
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentUpdate, zadanie),
        Command.AgentUpdate,
        (tresc) => czyObiekt(tresc.agent),
      );
    },

    async usun(idEksperta) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentDelete, { agentId: idEksperta }),
        Command.AgentDelete,
        (tresc) => typeof tresc.deleted === 'boolean',
      );
    },

    async ustawModel(zlecenie) {
      const parametry = odczytajJSON(zlecenie.parametry, 'parametry wywołania');
      if (!parametry.udany) return { udany: false, blad: parametry.blad };
      const zadanie: AgentModelSetRequest = {
        agentId: zlecenie.idEksperta,
        channelId: zlecenie.kanal,
      };
      if (zlecenie.model.trim() !== '') zadanie.model = zlecenie.model.trim();
      if (zlecenie.transport !== '') zadanie.transport = zlecenie.transport;
      if (parametry.wartosc !== undefined) zadanie.parameters = parametry.wartosc;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentModelSet, zadanie),
        Command.AgentModelSet,
        (tresc) => czyObiekt(tresc.agent),
      );
    },

    async dodajUmiejetnosc(idEksperta, idUmiejetnosci) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentSkillAdd, {
          agentId: idEksperta,
          skillId: idUmiejetnosci.trim(),
        }),
        Command.AgentSkillAdd,
        (tresc) => czyObiekt(tresc.agent),
      );
    },

    async dodajKonektor(zlecenie) {
      const konfiguracja = odczytajJSON(zlecenie.konfiguracja, 'konfiguracja integracji');
      if (!konfiguracja.udany) return { udany: false, blad: konfiguracja.blad };
      const zadanie: AgentConnectorAddRequest = {
        agentId: zlecenie.idEksperta,
        name: zlecenie.nazwa.trim(),
        kind: zlecenie.rodzaj,
      };
      if (zlecenie.idPunktu !== '') zadanie.accessPointId = zlecenie.idPunktu;
      if (konfiguracja.wartosc !== undefined) zadanie.config = konfiguracja.wartosc;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentConnectorAdd, zadanie),
        Command.AgentConnectorAdd,
        (tresc) => czyObiekt(tresc.connector),
      );
    },

    async ustawUprawnienie(zlecenie) {
      const zadanie: AgentPermissionSetRequest = {
        agentId: zlecenie.idEksperta,
        group: zlecenie.grupa,
        granted: zlecenie.przyznane,
      };
      if (zlecenie.zakres.trim() !== '') zadanie.scope = zlecenie.zakres.trim();
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentPermissionSet, zadanie),
        Command.AgentPermissionSet,
        (tresc) => czyTablica(tresc.permissions),
      );
    },

    /**
     * Zapis warstwy podaje stan czynności zawsze, choć kontrakt ma to pole
     * opcjonalne. Pominięcie znaczyłoby w rdzeniu „warstwa czynna", więc
     * wyłączenie warstwy nie doszłoby do skutku.
     *
     * Trybu podania żądanie nie niesie, bo kontrakt go nie ma: warstwa eksperta
     * dopisuje się do promptu systemowego jako zakres użytkownika i nigdy go nie
     * zastępuje (`warstwy-promptu.ts`).
     */
    async zapiszWarstwe(zlecenie) {
      const zadanie: AgentLayerSetRequest = {
        agentId: zlecenie.idEksperta,
        layer: zlecenie.warstwa,
        content: zlecenie.tresc,
        enabled: zlecenie.czynna,
      };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentLayerSet, zadanie),
        Command.AgentLayerSet,
        (tresc) => czyObiekt(tresc.agent),
      );
    },

    async usunWarstwe(idEksperta, warstwa) {
      const zadanie: AgentLayerRemoveRequest = { agentId: idEksperta, layer: warstwa };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentLayerRemove, zadanie),
        Command.AgentLayerRemove,
        (tresc) => czyObiekt(tresc.agent),
      );
    },

    async dodajWtyczke(zlecenie) {
      const zadanie: AgentPluginAddRequest = {
        agentId: zlecenie.idEksperta,
        name: zlecenie.nazwa.trim(),
      };
      if (zlecenie.zrodlo.trim() !== '') zadanie.source = zlecenie.zrodlo.trim();
      if (zlecenie.wersja.trim() !== '') zadanie.version = zlecenie.wersja.trim();
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentPluginAdd, zadanie),
        Command.AgentPluginAdd,
        (tresc) => czyObiekt(tresc.plugin),
      );
    },

    async wtyczki(idEksperta) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentPluginList, { agentId: idEksperta }),
        Command.AgentPluginList,
        (tresc) => czyTablica(tresc.plugins),
      );
    },

    async odlaczWtyczke(idEksperta, idWtyczki) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AgentPluginRemove, {
          agentId: idEksperta,
          pluginId: idWtyczki,
        }),
        Command.AgentPluginRemove,
        (tresc) => czyObiekt(tresc.agent),
      );
    },

    naZmiane(sluchacz) {
      return kanal.naZdarzenie(EventType.AgentChanged, (tresc) => sluchacz(tresc));
    },
  };
}

/** Wynik odczytu treści JSON wpisanej w oknie. */
interface OdczytJSON {
  udany: boolean;
  wartosc?: unknown;
  blad?: ErrorInfo;
}

/** Odczyt JSON wpisanego przez Operatora; pusty tekst znaczy brak wartości. */
function odczytajJSON(tekst: string, nazwaPola: string): OdczytJSON {
  if (tekst.trim() === '') return { udany: true };
  try {
    return { udany: true, wartosc: JSON.parse(tekst) as unknown };
  } catch (blad) {
    return {
      udany: false,
      blad: {
        code: ErrorCode.ValidationFailed,
        message: `Pole „${nazwaPola}" nie jest poprawnym JSON: ${String(blad)}`,
        retryable: false,
      },
    };
  }
}
