/**
 * Składnik — kolumna tożsamości. Lewa strefa okna: sygnet z pulsującą
 * kropką, nazwa, motto i trzy zdania o tym, czym program jest.
 */

import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, wykazZalet, zeZnacznika } from '../narzedzia.ts';

/** Odsłona kolumny; nazwa wskazuje gałąź katalogu z mottem i zdaniami właściwymi tej odsłonie tego okna. */
export type OdslonaKolumny = 'uruchamianie' | 'dostep' | 'przygotowanie';

export interface WlasciwosciKolumny {
  odslona: OdslonaKolumny;
}

/** Znaki trzech zdań, po jednym na każde zdanie, ułożone w kolejności zgodnej z kolejnością zdań w kolumnie. */
const ZNAKI_ZALET: NazwaZnaku[] = ['srodowiska', 'modele', 'tarcza'];

export function kolumnaTozsamosci(w: WlasciwosciKolumny): HTMLElement {
  const godlo = zeZnacznika(ikony.godlo);
  godlo.setAttribute('class', 'we-marka-godlo');
  godlo.setAttribute('role', 'img');
  godlo.setAttribute('aria-label', tekst('marka.nazwa'));

  const glowa = el('div', {}, [
    godlo,
    el('p', { klasa: 'we-marka-nazwa', tekst: tekst('marka.nazwa') }),
    el('p', { klasa: 'we-marka-motto', tekst: tekst(`marka.${w.odslona}.motto`) }),
  ]);

  if (w.odslona === 'przygotowanie') {
    return el('div', { klasa: 'we-marka' }, [
      glowa,
      // Pole puste — rysunek wnosi składnik biblioteki, nie znacznik okna.
      el('div', {
        klasa: 'pg-bryla-pole dn-powloki',
        role: 'img',
        'aria-label': tekst('przygotowanie.obszarPowlok'),
        dane: { powloki: true },
      }),
    ]);
  }

  const zalety = wykazZalet(`marka.${w.odslona}.zalety`).map((zaleta, i) => {
    const znak = zeZnacznika(ikony[ZNAKI_ZALET[i] ?? 'tarcza']);
    znak.setAttribute('aria-hidden', 'true');
    return el('li', { klasa: 'we-marka-poz' }, [
      znak,
      el('span', {}, [el('b', { tekst: zaleta.glowa }), zaleta.tresc]),
    ]);
  });

  return el('div', { klasa: 'we-marka' }, [glowa, el('ul', { klasa: 'we-marka-lista' }, zalety)]);
}

/**
 * Nota w wierszu pasa działań — osobny węzeł, bo siada w innym wierszu siatki
 * okna. Wydawca w oknach przed wejściem, stan przywracania w oknie
 * przygotowania: w obu miejscach to rzecz, która ma stać przy dnie i nie
 * należy do żadnej odsłony.
 */
export function notaPasa(klucz: string): HTMLElement {
  return el('p', { klasa: 'we-marka-stopka', tekst: tekst(klucz) });
}

export function notaWydawcy(wersjaKlienta: string): HTMLElement {
  return el('p', { klasa: 'we-marka-stopka' }, [
    tekst('marka.wydawca'),
    el('br'),
    tekst('marka.wydanie', { wersja: wersjaKlienta }),
    el('br'),
    tekst('marka.wsparcie'),
  ]);
}
