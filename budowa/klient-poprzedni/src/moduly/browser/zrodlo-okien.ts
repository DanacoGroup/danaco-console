import {
  Command,
  type Window,
  type WindowListResponse,
  type WindowStateGetResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { KOD_MODULU } from './etykiety-browser';

/**
 * Okno przeglądarki, do którego moduł adresuje swoje komendy: odczyt okien
 * sesji komendą window.list, wskazanie tego okna, które rdzeń przypisał
 * modułowi Browser, oraz odczyt jego stanu komendą window.state.get.
 */
export interface ZrodloOkien {
  /** Okna wskazanej sesji w kolejności nadanej przez rdzeń. */
  okna(idSesji: string): Promise<Wynik<WindowListResponse>>;
  /** Stan okna wraz z parametrami wykonania. */
  stan(idOkna: string): Promise<Wynik<WindowStateGetResponse>>;
}

export function utworzZrodloOkien(kanal: Kanal): ZrodloOkien {
  return {
    async okna(idSesji) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowList, { sessionId: idSesji }),
        Command.WindowList,
        (tresc) => czyTablica(tresc.windows),
      );
    },

    async stan(idOkna) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowStateGet, { windowId: idOkna }),
        Command.WindowStateGet,
        (tresc) => czyObiekt(tresc.window),
      );
    },
  };
}

/**
 * Okno modułu Browser spośród okien sesji, wskazane po zgodności pola moduleId
 * z kodem modułu. Gdy takiego okna nie ma, funkcja nie podstawia pierwszego
 * z brzegu, tylko oddaje wartość pustą do nazwania braku.
 */
export function oknoModulu(okna: readonly Window[]): Window | null {
  return okna.find((okno) => okno.moduleId === KOD_MODULU) ?? null;
}

/**
 * Zdanie o oknie przeglądarki pokazywane w pasie stanu modułu: tytuł okna wraz
 * z identyfikatorem w nawiasie, a przy tytule pustym sam identyfikator okna.
 */
export function opisOkna(okno: Window): string {
  const tytul = (okno.title ?? '').trim();
  return tytul === '' ? okno.id : `${tytul} (${okno.id})`;
}
