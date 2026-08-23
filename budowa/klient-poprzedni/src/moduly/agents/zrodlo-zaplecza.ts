import {
  AccessPointKind,
  Command,
  SessionConfigArea,
  type AccessPoint,
  type AdapterCapabilities,
  type Channel,
  type IdentityCategory,
  type PermissionMode,
  type ProviderTransport,
  type Window,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Zaplecze modułu Agents — rejestry, z których moduł korzysta, a których nie
 * prowadzi. Każde wywołanie sięga do cudzego obszaru kontraktu:
 *
 *   `channel.list`             — rejestr kanałów modelu dla Model Configuration.
 *                                Moduł Agents nie zakłada kanałów; wybiera
 *                                spośród wierszy rejestru rdzenia.
 *   `config.capabilities.get`  — deklaracja zdolności adaptera dostawcy dla
 *                                pól konfiguracji sesji. Okno pokazuje wprost,
 *                                którego parametru wybrany kanał nie obsłuży,
 *                                zanim Operator go wypełni.
 *   `access.point.list`        — katalog mostów MCP dla Connectors Manager.
 *                                To ten sam katalog, z którego rdzeń składa
 *                                `mcpServers` procesu modelu.
 *   `identity.category.list`   — słownik kategorii tożsamości dla Agent Buildera.
 *   `window.list` / `window.update`
 *                              — tryb uprawnień per okno komunikacji dla
 *                                Permissions Center.
 */
export interface ZrodloZaplecza {
  kanaly(): Promise<Wynik<{ channels: Channel[] }>>;
  zdolnosci(
    idKanalu: string,
    transport: ProviderTransport | '',
  ): Promise<Wynik<{ capabilities: AdapterCapabilities }>>;
  mosty(): Promise<Wynik<{ points: AccessPoint[] }>>;
  kategorieTozsamosci(): Promise<Wynik<{ categories: IdentityCategory[] }>>;
  okna(idSesji: string): Promise<Wynik<{ windows: Window[] }>>;
  ustawTrybOkna(idOkna: string, tryb: PermissionMode): Promise<Wynik<{ window: Window }>>;
}

export function utworzZrodloZaplecza(kanal: Kanal): ZrodloZaplecza {
  return {
    async kanaly() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ChannelList, { enabledOnly: true }),
        Command.ChannelList,
        (tresc) => czyTablica(tresc.channels),
      );
    },

    async zdolnosci(idKanalu, transport) {
      // Obszar `model` wystarcza: Model Configuration wypełnia wyłącznie pola
      // modelu prowadzącego i parametrów jego wywołania. Pytanie o komplet
      // obszarów przyniosłoby deklarację, której okno nie pokaże.
      const zadanie = {
        area: SessionConfigArea.Model,
        ...(idKanalu === '' ? {} : { channelId: idKanalu }),
        ...(transport === '' ? {} : { transport }),
      };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ConfigCapabilitiesGet, zadanie),
        Command.ConfigCapabilitiesGet,
        (tresc) => czyObiekt(tresc.capabilities) && czyTablica(tresc.capabilities.fields),
      );
    },

    async mosty() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AccessPointList, {
          kind: AccessPointKind.McpBridge,
          enabledOnly: true,
        }),
        Command.AccessPointList,
        (tresc) => czyTablica(tresc.points),
      );
    },

    async kategorieTozsamosci() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.IdentityCategoryList, {}),
        Command.IdentityCategoryList,
        (tresc) => czyTablica(tresc.categories),
      );
    },

    async okna(idSesji) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowList, { sessionId: idSesji }),
        Command.WindowList,
        (tresc) => czyTablica(tresc.windows),
      );
    },

    async ustawTrybOkna(idOkna, tryb) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowUpdate, { windowId: idOkna, permissionMode: tryb }),
        Command.WindowUpdate,
        (tresc) => czyObiekt(tresc.window),
      );
    },
  };
}
