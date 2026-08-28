/**
 * Składnik — pasek postępu z głową. Nazwa tego, co trwa, wartość liczbą
 * i tor; wartość stoi w danej na znaczniku, nie w atrybucie stylu.
 */

import { el, tekst } from '../narzedzia.ts';

export interface WlasciwosciPaska {
  /** Klucz katalogu — co się dzieje. */
  etykieta: string;
  /** Klucz katalogu — nazwa toru dla czytnika ekranu. */
  opisPaska: string;
  /** Wartość od zera do stu. */
  wartosc: number;
}

export function pasekPostepu(w: WlasciwosciPaska): HTMLElement {
  return el('div', { klasa: 'pg-postep' }, [
    el('div', { klasa: 'pg-postep-wiersz' }, [
      el('span', { tekst: tekst(w.etykieta) }),
      el('b', { tekst: `${w.wartosc}%` }),
    ]),
    el('div', { klasa: 'dn-postep' }, [
      el(
        'div',
        {
          klasa: 'dn-postep-tor',
          role: 'progressbar',
          'aria-valuenow': String(w.wartosc),
          'aria-valuemin': '0',
          'aria-valuemax': '100',
          'aria-label': tekst(w.opisPaska),
        },
        [el('div', { klasa: 'dn-postep-wartosc', dane: { wartosc: w.wartosc } })],
      ),
    ]),
  ]);
}
