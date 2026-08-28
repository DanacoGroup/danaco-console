/**
 * SKŁADNIK — GŁOWA EKRANU.
 *
 * Nadtytuł mówi, gdzie stoisz; tytuł — co się dzieje; zdanie wprowadzające —
 * czego się spodziewać. Nadtytuł bywa zbędny: w oknie dostępu jego rolę pełnią
 * zakładki, więc powtarzanie go byłoby szumem.
 */

import { el, tekst, type DanePodstawienia, type Dziecko } from '../narzedzia.ts';

export interface WlasciwosciGlowy {
  /** Klucz katalogu — nadtytuł; pomijalny. */
  nadtytul?: string | null;
  /** Klucz katalogu — tytuł. */
  tytul: string;
  /** Klucz katalogu — zdanie wprowadzające; pomijalne. */
  lid?: string | null;
  /** Dane podstawiane w zdanie wprowadzające. */
  daneLidu?: DanePodstawienia;
  /** Węzły zamiast tekstu zdania — dla wyróżnionego fragmentu, na przykład adresu krojem maszynowym. */
  lidWezly?: Dziecko[] | null;
  /** Węzły między tytułem a zdaniem wprowadzającym; tor kroków należy do głowy, nie stoi obok niej. */
  poTytule?: Dziecko[] | null;
}

export function naglowekEkranu(w: WlasciwosciGlowy): HTMLElement {
  const przedZdaniem: Dziecko[] = [
    w.nadtytul ? el('p', { klasa: 'we-nadtytul', tekst: tekst(w.nadtytul) }) : null,
    el('h2', { klasa: 'we-tytul', tekst: tekst(w.tytul) }),
  ];
  const zdanie: Dziecko = w.lidWezly
    ? el('p', { klasa: 'we-lid' }, w.lidWezly)
    : w.lid
      ? el('p', { klasa: 'we-lid', tekst: tekst(w.lid, w.daneLidu) })
      : null;

  return el('div', { klasa: 'we-glowa' }, [...przedZdaniem, ...(w.poTytule ?? []), zdanie]);
}
