import {
  Command,
  type Window,
  type WindowStateGetResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Osadzenie modułu w sesji — okna komunikacji i parametry wykonania.
 *
 * Każda komenda `studio.*` wymaga identyfikatora okna (`windowId` jest polem
 * obowiązkowym `studio.document.open`), a moduł dostaje z powłoki wyłącznie
 * identyfikator sesji. Wskazanie okna poprzedza więc całą pracę modułu i nie
 * należy do żadnego okna operacyjnego z osobna — stąd własny plik, a nie
 * doklejenie do `zrodlo-studio.ts`.
 *
 * `window.state.get` jest komendą wspólną każdemu oknu operacyjnemu: niesie
 * moduł, kanał modelu, katalogi robocze, zasięg wykonania, tryb uprawnień
 * i rolę.
 */
export interface ZrodloOsadzenia {
  /** Okna komunikacji sesji — z nich bierze się `windowId` komend modułu. */
  okna(idSesji: string): Promise<Wynik<{ windows: Window[] }>>;
  /** Parametry wykonania okna wraz z konfiguracją efektywną. */
  stanOkna(idOkna: string): Promise<Wynik<WindowStateGetResponse>>;
}

export function utworzZrodloOsadzenia(kanal: Kanal): ZrodloOsadzenia {
  return {
    async okna(idSesji) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowList, { sessionId: idSesji }),
        Command.WindowList,
        (tresc) => czyTablica(tresc.windows),
      );
    },

    async stanOkna(idOkna) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowStateGet, { windowId: idOkna, includeConfig: true }),
        Command.WindowStateGet,
        (tresc) => czyObiekt(tresc.window),
      );
    },
  };
}
