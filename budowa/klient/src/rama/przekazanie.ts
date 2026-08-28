/**
 * Rozstrzygnięcie i wykonanie chwili przekazania sterowania. Droga wejścia
 * kończy się komendą `environment.enter` — przebieg zapisuje jej wynik w
 * `srodowisko` dopiero po powodzeniu, więc obecność tego pola jest jedynym
 * pewnym znakiem, że rama aplikacji ma dokąd prowadzić.
 */

import type { Environment } from '../../../shared/contract.ts';
import type { StanPrzebiegu } from '../wejscie/przebieg.ts';
import { zamontujRame } from './montaz.ts';

/** Stan przebiegu, w którym środowisko jest już znane — warunek konieczny wykonania przekazania. */
type StanZeSrodowiskiem = StanPrzebiegu & { srodowisko: Environment };

/** Czy przebieg doszedł do miejsca, w którym rama aplikacji przejmuje ekran po drodze wejścia. */
export function gotowaDoPrzekazania(stan: StanPrzebiegu): stan is StanZeSrodowiskiem {
  return stan.etap === 'przygotowanie' && stan.srodowisko !== undefined;
}

/** Węzły dokumentu i działanie, których przekazanie potrzebuje od wywołującego. */
export interface ZaleznosciPrzekazania {
  dokument: Document;
  zdejmijOknoWejscia: () => void;
}

/**
 * Wykonuje przekazanie sterowania ramie: zdejmuje scenę drogi wejścia,
 * montuje ramę i odkrywa jej miejsce. Brak węzła montażu w dokumencie jest
 * odmową nazwaną w dzienniku, nie cichym zaniechaniem — wywołujący rozpoznaje
 * niepowodzenie po zwróconej wartości `false` i może spróbować przy kolejnej
 * zmianie stanu.
 */
export function wykonajPrzekazanie(stan: StanZeSrodowiskiem, zaleznosci: ZaleznosciPrzekazania): boolean {
  const scenaWejscia = zaleznosci.dokument.querySelector('[data-wejscie]');
  const miejsceRamy = zaleznosci.dokument.querySelector('[data-rama-aplikacji]');
  if (scenaWejscia === null || miejsceRamy === null) {
    console.error('[rama] przekazanie odrzucone: brak węzła montażu w dokumencie');
    return false;
  }

  zaleznosci.zdejmijOknoWejscia();
  (scenaWejscia as HTMLElement).hidden = true;
  zamontujRame({
    miejsce: miejsceRamy as HTMLElement,
    srodowisko: stan.srodowisko,
    moduly: stan.moduly,
    sesje: stan.sesje,
  });
  (miejsceRamy as HTMLElement).hidden = false;
  return true;
}
