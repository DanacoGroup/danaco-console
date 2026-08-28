import {
  Command,
  type Window,
  type WindowStateGetResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Osadzenie modułu w sesji — okna komunikacji i parametry wykonania; każda komenda `studio.*`
 * wymaga identyfikatora okna, którego moduł nie dostaje z powłoki wprost.
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
