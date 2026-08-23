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
 * Okno przeglądarki, do którego moduł adresuje swoje komendy.
 *
 * Jedna odpowiedzialność: odczyt okien sesji (`window.list`) i wskazanie tego,
 * które rdzeń przypisał modułowi Browser, oraz odczyt jego stanu
 * (`window.state.get`).
 *
 * Każda komenda obszaru `browser.*` wymaga pola `windowId`, a moduł dostaje
 * z powłoki wyłącznie identyfikator sesji. Ustalenie okna stoi więc w jednym
 * miejscu — pytanie o nie osobno w każdym oknie modułu dałoby kilka różnych
 * okien dla jednego modułu.
 *
 * Te dwie komendy należą do obszaru okien, nie do obszaru `browser.*` — stąd
 * osobne źródło obok `zrodlo-browser.ts`.
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
 * Okno modułu Browser spośród okien sesji.
 *
 * Rdzeń przestawia moduł okna komendą `workspace.enter`, więc oknem
 * przeglądarki jest to, którego `moduleId` równa się kodowi modułu. Gdy takiego
 * nie ma, funkcja nie podstawia okna pierwszego z brzegu — oddaje `null`, żeby
 * moduł mógł brak nazwać.
 */
export function oknoModulu(okna: readonly Window[]): Window | null {
  return okna.find((okno) => okno.moduleId === KOD_MODULU) ?? null;
}

/** Zdanie o oknie przeglądarki pokazywane w pasie stanu modułu. */
export function opisOkna(okno: Window): string {
  const tytul = (okno.title ?? '').trim();
  return tytul === '' ? okno.id : `${tytul} (${okno.id})`;
}
