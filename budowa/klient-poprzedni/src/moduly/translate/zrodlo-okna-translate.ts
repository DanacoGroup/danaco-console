import {
  Command,
  type WindowListResponse,
  type WindowStateGetResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Dwie komendy wspólne każdemu oknu operacyjnemu jako bytowi rdzenia: wykaz
 * okien sesji oraz odczyt parametrów wykonania jednego okna.
 */
export interface ZrodloOknaTranslate {
  /** Okna komunikacji sesji — wśród nich okno modułu Translate. */
  okna(idSesji: string): Promise<Wynik<WindowListResponse>>;
  /** Parametry wykonania okna: moduł, tryb uprawnień, zasięg, rola, stan procesu. */
  stan(idOkna: string): Promise<Wynik<WindowStateGetResponse>>;
}

export function utworzZrodloOknaTranslate(kanal: Kanal): ZrodloOknaTranslate {
  return {
    async okna(idSesji) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowList, { sessionId: idSesji }),
        Command.WindowList,
        (tresc) => czyTablica(tresc.windows),
      );
    },

    async stan(idOkna) {
      // Konfiguracji efektywnej okna moduł nie pobiera, bo żadne z jego okien jej nie pokazuje.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowStateGet, { windowId: idOkna }),
        Command.WindowStateGet,
        (tresc) => czyObiekt(tresc.window),
      );
    },
  };
}
