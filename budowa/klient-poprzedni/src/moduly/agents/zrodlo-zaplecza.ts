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
 * Zaplecze modułu Agents, czyli rejestry, z których moduł korzysta, a których
 * nie prowadzi. Każde wywołanie sięga do cudzego obszaru kontraktu: rejestru
 * kanałów, zdolności adaptera, katalogu mostów, kategorii tożsamości oraz okien.
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
      // Obszar `model` wystarcza: okno wypełnia wyłącznie pola modelu.
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
