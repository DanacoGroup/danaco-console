/**
 * Rozstrzygnięcie chwili przekazania sterowania. Droga wejścia kończy się
 * komendą `environment.enter` — przebieg zapisuje jej wynik w `srodowisko`
 * dopiero po powodzeniu, więc obecność tego pola jest jedynym pewnym znakiem,
 * że okno robocze ma dokąd prowadzić.
 */

import type { StanPrzebiegu } from '../wejscie/przebieg.ts';

/** Czy przebieg doszedł do miejsca, w którym rama aplikacji przejmuje ekran po drodze wejścia. */
export function gotowaDoPrzekazania(stan: StanPrzebiegu): boolean {
  return stan.etap === 'przygotowanie' && stan.srodowisko !== undefined;
}
