/**
 * Składnik — fraza nawigacyjna. Zdanie nad pasem działań: mówi, co będzie
 * dalej, albo stawia pytanie z wplecioną czynnością.
 */

import { el, tekst } from '../narzedzia.ts';

export interface WlasciwosciFrazy {
  /** Klucz katalogu — całe zdanie, bez odsyłacza. */
  klucz?: string | null;
  /** Klucz katalogu — część przed odsyłaczem. */
  pytanie?: string | null;
  /** Klucz katalogu — napis odsyłacza. */
  czynnosc?: string | null;
  /** Nazwa odsłony, do której odsyłacz prowadzi. */
  cel?: string | null;
}

export function frazaNawigacyjna(w: WlasciwosciFrazy): HTMLElement {
  if (w.klucz) {
    return el('p', { klasa: 'we-fraza' }, [el('span', { tekst: tekst(w.klucz) })]);
  }
  return el('p', { klasa: 'we-fraza' }, [
    w.pytanie ? el('span', { klasa: 'we-fraza-pytanie', tekst: tekst(w.pytanie) }) : null,
    w.pytanie ? ' ' : null,
    el('a', {
      klasa: 'au-link',
      href: `#s-${w.cel ?? ''}`,
      tekst: w.czynnosc ? tekst(w.czynnosc) : '',
      dane: { idz: w.cel ?? null },
    }),
  ]);
}
