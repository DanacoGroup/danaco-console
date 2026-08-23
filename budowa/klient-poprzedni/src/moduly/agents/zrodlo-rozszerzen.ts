import {
  Command,
  type Extension,
  type ExtensionKind,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Źródło katalogu rozszerzeń — pięć komend rodziny `extension.*`.
 *
 * Katalog stoi osobno od `zrodlo-zaplecza`, bo obie warstwy odpowiadają za co
 * innego: zaplecze niesie rejestry, z których moduł Agents tylko korzysta —
 * kanały modelu, mosty MCP, okna sesji — a katalogiem rozszerzeń moduł zarządza,
 * instaluje i odinstalowuje pozycje.
 *
 * Katalog obejmuje rodzaje wymienione w `ExtensionKind`: serwer MCP, wtyczkę,
 * integrację API i skill.
 *
 * Źródło przekazuje pola kontraktu i oddaje odpowiedź rdzenia bez dopowiedzenia;
 * o tym, czym jest instalacja pozycji, rozstrzyga rdzeń, nie ta warstwa.
 */
export interface ZrodloRozszerzen {
  /** `extension.list` — katalog w kolejności wyświetlania, opcjonalnie zawężony. */
  katalog(zawezenie?: {
    rodzaj?: ExtensionKind;
    idEksperta?: string;
    tylkoZainstalowane?: boolean;
  }): Promise<Wynik<{ extensions: Extension[] }>>;
  /** `extension.install` — instalacja pozycji katalogu. */
  zainstaluj(
    kod: string,
    rodzaj: ExtensionKind,
    zrodlo?: string,
  ): Promise<Wynik<{ extension: Extension }>>;
  /** `extension.configure` — zapis konfiguracji rozszerzenia. */
  skonfiguruj(idRozszerzenia: string, konfiguracja: unknown): Promise<Wynik<{ extension: Extension }>>;
  /** `extension.toggle` — włączenie albo wyłączenie bez odinstalowania. */
  przelacz(idRozszerzenia: string, wlaczone: boolean): Promise<Wynik<{ extension: Extension }>>;
  /** `extension.uninstall` — odinstalowanie pozycji. */
  odinstaluj(idRozszerzenia: string): Promise<Wynik<{ uninstalled: boolean }>>;
}

export function utworzZrodloRozszerzen(kanal: Kanal): ZrodloRozszerzen {
  return {
    async katalog(zawezenie = {}) {
      const zadanie = {
        ...(zawezenie.rodzaj === undefined ? {} : { kind: zawezenie.rodzaj }),
        ...(zawezenie.idEksperta === undefined ? {} : { agentId: zawezenie.idEksperta }),
        ...(zawezenie.tylkoZainstalowane === true ? { installedOnly: true } : {}),
      };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ExtensionList, zadanie),
        Command.ExtensionList,
        (tresc) => czyTablica(tresc.extensions),
      );
    },

    async zainstaluj(kod, rodzaj, zrodlo) {
      const zadanie = {
        code: kod,
        kind: rodzaj,
        ...(zrodlo === undefined || zrodlo === '' ? {} : { source: zrodlo }),
      };
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ExtensionInstall, zadanie),
        Command.ExtensionInstall,
        (tresc) => czyObiekt(tresc.extension),
      );
    },

    async skonfiguruj(idRozszerzenia, konfiguracja) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ExtensionConfigure, {
          extensionId: idRozszerzenia,
          config: konfiguracja,
        }),
        Command.ExtensionConfigure,
        (tresc) => czyObiekt(tresc.extension),
      );
    },

    async przelacz(idRozszerzenia, wlaczone) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ExtensionToggle, {
          extensionId: idRozszerzenia,
          enabled: wlaczone,
        }),
        Command.ExtensionToggle,
        (tresc) => czyObiekt(tresc.extension),
      );
    },

    async odinstaluj(idRozszerzenia) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ExtensionUninstall, { extensionId: idRozszerzenia }),
        Command.ExtensionUninstall,
        (tresc) => typeof tresc.uninstalled === 'boolean',
      );
    },
  };
}
