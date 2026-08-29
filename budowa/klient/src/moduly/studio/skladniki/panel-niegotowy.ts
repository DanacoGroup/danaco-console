/**
 * Strefa 5 — zaczep panelu w paśmie kart. Powłoka wstawia go pod każdą kartą
 * i to on niesie identyfikator, rolę i powiązanie z kartą; panel, który ma
 * swojego wykonawcę, wchodzi w ten węzeł przy pierwszym wejściu na kartę.
 * Karta bez wykonawcy zostaje przy tej treści: nazwany stan pusty, nie pusty
 * prostokąt i nie treść przykładowa przeniesiona z prototypu.
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
