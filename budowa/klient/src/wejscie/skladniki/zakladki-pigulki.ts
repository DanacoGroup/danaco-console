/**
 * Składnik — zakładki logowania i rejestracji. Dwie drogi równorzędne;
 * wygląd wnosi wariant biblioteki, stan wybrany niesie podświetlenie, nie
 * wypełnienie.
 */

import { el, tekst, wykaz } from '../narzedzia.ts';

/** Dwie drogi równorzędne, w kolejności zakładek katalogu, pomiędzy którymi przełącza się ten składnik. */
export type Pigulka = 'logowanie' | 'rejestracja';

const CELE: Pigulka[] = ['logowanie', 'rejestracja'];

export interface WlasciwosciZakladek {
  wybrana: Pigulka;
}

export function zakladkiPigulki(w: WlasciwosciZakladek): HTMLElement {
  const napisy = wykaz('dostep.zakladki');
  return el(
    'div',
    {
      klasa: 'dn-zakladki dn-zakladki--pigulki dn-zakladki--wybor',
      role: 'tablist',
      'aria-label': tekst('dostep.logowanie.tytul'),
    },
    napisy.map((napis, i) => {
      const cel = CELE[i] ?? 'logowanie';
      return el('button', {
        klasa: 'dn-zakladka',
        type: 'button',
        role: 'tab',
        'aria-selected': cel === w.wybrana ? 'true' : 'false',
        'aria-controls': `s-${cel}`,
        tekst: napis,
        dane: { idz: cel },
      });
    }),
  );
}
