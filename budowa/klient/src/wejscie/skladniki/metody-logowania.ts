/**
 * SKŁADNIK — INNE METODY LOGOWANIA.
 *
 * Trzy drogi poboczne, każda ze stanem dostępności. Metoda nieaktywna zostaje
 * w oknie zamiast znikać: pusta lista dróg pobocznych nie mówi, że są jakieś,
 * a lista z wygaszonymi pozycjami mówi, co można włączyć w ustawieniach.
 *
 * Składnik zwraca WYKAZ węzłów, nie jeden węzeł: etykieta i siatka metod są
 * rodzeństwem w kolumnie panelu. Opakowanie ich w pudełko wprowadziłoby
 * dodatkowy poziom, przez który odstęp kolumny liczyłby się raz zamiast dwa.
 */

import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

/** Klucz metody w katalogu treści. */
export type KluczMetody = 'pin' | 'klucz' | 'email';

/** Trzy metody poboczne wraz ze znakiem każdej z nich. */
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
        return el('button', { klasa: 'au-metoda', type: 'button' }, [
          znak,
          el('span', { tekst: tekst(`dostep.logowanie.metody.${m.klucz}`) }),
          el('span', {
            klasa: 'au-metoda-stan',
            tekst: tekst(`dostep.logowanie.metody.${czynna ? 'aktywna' : 'nieaktywna'}`),
          }),
        ]);
      }),
    ),
  ];
}
