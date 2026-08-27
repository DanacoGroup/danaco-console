/**
 * SKŁADNIK — ZGODA NA TRWAŁĄ SESJĘ.
 *
 * Pole wyboru z opisem skutku. Opis mówi, co ta zgoda daje i kiedy jej nie
 * zaznaczać — samo „pozostań zalogowany” nie niesie ani jednego, ani drugiego.
 *
 * Zgoda idzie do rdzenia polem `keepSignedIn`: sesja bramki dostaje trwanie
 * długie zamiast doby roboczej. Bramką nie jest — znosi powtarzanie logowania,
 * niczego nie blokuje.
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
