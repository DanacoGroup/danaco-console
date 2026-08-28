/**
 * Strefa 5 — zaczep paneli: treść jednej karty jeszcze nie postawionej.
 * Sześć z siedmiu kart tego terenu (wszystkie poza Studio Editor) wchodzą
 * osobnym zakresem prac — tu stoi nazwany stan pusty, nie pusty prostokąt
 * i nie treść przykładowa przeniesiona z prototypu.
 */

import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';
import type { DefinicjaKarty } from './definicje.ts';

export function panelNiegotowy(k: DefinicjaKarty): HTMLElement {
  const znakPanelu = zeZnacznika(ikony.pusto);
  znakPanelu.setAttribute('aria-hidden', 'true');

  return el(
    'section',
    { klasa: 'sta-okno', id: `panel-${k.kod}`, role: 'tabpanel', 'aria-labelledby': `karta-${k.kod}`, hidden: true },
    [
      el('header', { klasa: 'sta-okno-belka' }, [
        el('span', { klasa: 'sta-okno-tytul' }, [el('b', { tekst: k.nazwa })]),
      ]),
      el('div', { klasa: 'sta-okno-tresc' }, [
        el('div', { klasa: 'dn-pusty-stan' }, [
          znakPanelu,
          el('span', { klasa: 'dn-pusty-stan-tytul', tekst: tekst('panelNiegotowy.tytul') }),
          el('span', { klasa: 'dn-pusty-stan-opis', tekst: tekst('panelNiegotowy.opis') }),
        ]),
      ]),
    ],
  );
}
