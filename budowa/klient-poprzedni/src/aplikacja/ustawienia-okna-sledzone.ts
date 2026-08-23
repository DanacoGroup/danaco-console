import type { Window } from '../../../shared/contract';

/**
 * Śledzenie ustawień okna między zdarzeniami.
 *
 * Zdarzenie `window.changed` niesie okno po zmianie, bez pola „co zmieniono"
 * i bez pola „kto zmienił". Ten plik odtwarza jedno i drugie porównaniem
 * stanów, nie dopisując niczego do kontraktu:
 *
 *  • {@link migawka} i {@link roznice} mówią, co się zmieniło. Bez nich pas
 *    meldowałby „okno zmienione" przy każdym dotknięciu, zamiast nazwać dobrany
 *    model, rolę albo zasięg.
 *
 *  • {@link odciskOkna} mówi, czy zmiana jest własna: zdarzenie przynoszące
 *    dokładnie ten stan okna, co świeża odpowiedź na własną komendę, jest
 *    skutkiem tej komendy. Droga jest opisana w nagłówku `zrodlo-posuniec.ts`.
 *
 * Plik nie zna ani widoku, ani kanału — sam przekład pól okna na wartości
 * porównywalne i na zdania dla czytającego.
 */

/** Ustawienia okna śledzone między zdarzeniami — wprost z pól `Window`. */
export interface MigawkaUstawien {
  modul: string;
  kanal: string;
  agent: string;
  rola: string;
  zasieg: string;
  tryb: string;
}

export function migawka(okno: Window): MigawkaUstawien {
  return {
    modul: okno.moduleId,
    kanal: okno.modelChannelId,
    agent: okno.agentId ?? '',
    rola: okno.windowRole,
    zasieg: okno.executionEnv,
    tryb: okno.permissionMode,
  };
}

/**
 * Odcisk stanu okna — klucz porównania zdarzenia z własną odpowiedzią.
 *
 * Wchodzą tu wyłącznie pola, którymi da się sterować komendą: identyfikator
 * i ustawienia. Czasu zmiany (`updatedAt`) tu nie ma — rdzeń nadaje go osobno
 * w zdarzeniu i w odpowiedzi, więc doklejenie go sprawiłoby, że własna czynność
 * nigdy nie zrównałaby się sama ze sobą.
 */
export function odciskOkna(okno: Window): string {
  const u = migawka(okno);
  return [okno.id, u.modul, u.kanal, u.agent, u.rola, u.zasieg, u.tryb].join('|');
}

/**
 * Zdania o ustawieniach, które naprawdę się zmieniły.
 *
 * Modułu w tym wykazie nie ma: ma własny rodzaj posunięcia i własny skutek
 * (podążanie nawigacji), więc meldowany tu byłby po raz drugi.
 */
export function roznice(przed: MigawkaUstawien, po: MigawkaUstawien): string[] {
  const opisy: Array<[keyof MigawkaUstawien, string]> = [
    ['kanal', 'kanał modelu'],
    ['agent', 'eksperta okna'],
    ['rola', 'rolę okna'],
    ['zasieg', 'zasięg wykonania'],
    ['tryb', 'tryb uprawnień'],
  ];
  const zdania: string[] = [];
  for (const [pole, nazwa] of opisy) {
    if (przed[pole] === po[pole]) continue;
    const wartosc = po[pole] === '' ? 'brak' : po[pole];
    zdania.push(`Ustawiony ${nazwa}: ${wartosc}`);
  }
  return zdania;
}
