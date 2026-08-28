/**
 * Składnik — zgoda na trwałą sesję. Pole wyboru z opisem skutku: co ta
 * zgoda daje i kiedy jej nie zaznaczać.
 */

import { el, tekst } from '../narzedzia.ts';

export interface WlasciwosciSesji {
  /** Identyfikator kontrolki; przebieg odczytuje z niej zgodę. */
  id: string;
}

export function poleSesji(w: WlasciwosciSesji): HTMLElement {
  return el('div', { klasa: 'au-opcje' }, [
    el('label', { klasa: 'au-opcja' }, [
      el('input', { klasa: 'dn-check', type: 'checkbox', id: w.id }),
      el('span', {}, [
        tekst('dostep.sesja.etykieta'),
        el('small', { tekst: tekst('dostep.sesja.opis') }),
      ]),
    ]),
  ]);
}
