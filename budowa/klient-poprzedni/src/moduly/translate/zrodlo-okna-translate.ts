import {
  Command,
  type WindowListResponse,
  type WindowStateGetResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Okno operacyjne jako byt rdzenia — dwie komendy wspólne każdemu oknu.
 *
 * Komendy Source Panel i Translation Panels wymagają `windowId`: bez niego
 * moduł nie ma czym zaadresować ani `translate.source.set`, ani
 * `translate.target.add`. Identyfikator bierze się z rdzenia (`window.list`
 * sesji), nigdy z literału po stronie klienta.
 *
 * Te dwie komendy idą zwykłym `wywolaj`, nie drogą odmowy z
 * `wywolanie-translate.ts`: rdzeń je obsługuje, więc odpowiedź przychodzi
 * kopertą ze statusem i korelacja ją rozpoznaje.
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
      // Konfiguracji efektywnej okna moduł nie prosi: żadne z jego okien jej
      // nie pokazuje, a wykaz bez odbiorcy byłby ruchem na pokaz.
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowStateGet, { windowId: idOkna }),
        Command.WindowStateGet,
        (tresc) => czyObiekt(tresc.window),
      );
    },
  };
}
