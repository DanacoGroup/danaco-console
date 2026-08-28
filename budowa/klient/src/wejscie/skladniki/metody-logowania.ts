/**
 * Składnik — inne metody logowania. Trzy drogi poboczne, każda ze stanem
 * dostępności; metoda nieaktywna zostaje w oknie zamiast znikać.
 */

import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

/** Klucz metody w katalogu treści, wskazujący jedną z trzech dróg pobocznych logowania wobec głównej drogi. */
export type KluczMetody = 'pin' | 'klucz' | 'email';

/** Trzy metody poboczne wraz ze znakiem każdej z nich, ułożone w kolejności, w jakiej stają w siatce metod. */
const METODY: { klucz: KluczMetody; ikona: NazwaZnaku }[] = [
  { klucz: 'pin', ikona: 'klawiatura' },
  { klucz: 'klucz', ikona: 'tarcza' },
  { klucz: 'email', ikona: 'koperta' },
];

export interface WlasciwosciMetod {
  /** Klucze metod czynnych; pozostałe zostają wygaszone. */
  dostepne: KluczMetody[];
}

export function metodyLogowania(w: WlasciwosciMetod): HTMLElement[] {
  return [
    el('p', { klasa: 'au-etykieta', tekst: tekst('dostep.logowanie.metody.naglowek') }),
    el(
      'div',
      { klasa: 'au-metody' },
      METODY.map((m) => {
        const znak = zeZnacznika(ikony[m.ikona]);
        znak.setAttribute('aria-hidden', 'true');
        const czynna = w.dostepne.includes(m.klucz);
        return el('button', { klasa: 'dn-kafel dn-kafel--wybor', type: 'button' }, [
          znak,
          el('span', { tekst: tekst(`dostep.logowanie.metody.${m.klucz}`) }),
          el('span', {
            klasa: 'dn-kafel-stan',
            tekst: tekst(`dostep.logowanie.metody.${czynna ? 'aktywna' : 'nieaktywna'}`),
          }),
        ]);
      }),
    ),
  ];
}
