/**
 * Pasmo kart okna głównego. Karta „Centrum" przypięta na stałe, po niej po
 * jednej karcie na kartę sesji odtworzoną przez rdzeń — ten sam znacznik,
 * który czyta `karty-okna.js` (`.dn-karty-pasmo > .dn-karty-lista
 * > .dn-karta-widoku[data-karta][data-karta-rodzaj]`).
 */

import type { Session } from '../../../../shared/contract.ts';
import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciPasma {
  /** Karty sesji otwarte w środowisku, w kolejności odebranej od rdzenia. */
  sesje: Session[];
  /** Węzeł, do którego wszystkie karty prowadzą — teren nie rozdziela dziś treści karta po karcie. */
  idTresci: string;
}

function znak(rysunek: string, klasa?: string): SVGElement {
  const wezel = zeZnacznika(rysunek);
  /* `setAttribute('class', …)`, nie `classList.add` — na elemencie SVG w prawdziwym
     dokumencie `className` jest odczytywalną wyłącznie właściwością `SVGAnimatedString`. */
  if (klasa) wezel.setAttribute('class', klasa);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function kartaGlowna(idTresci: string): HTMLElement {
  return el(
    'div',
    {
      klasa: 'dn-karta-widoku dn-karta-widoku--glowna',
      role: 'tab',
      tabindex: '0',
      'aria-selected': 'true',
      'data-karta': 'centrum',
      'data-karta-rodzaj': 'centrum',
      'aria-controls': idTresci,
    },
    [
      znak(ikony.godlo, 'dn-karta-widoku-ikona'),
      el('span', { klasa: 'dn-karta-widoku-nazwa', tekst: tekst('belka.marka') }),
    ],
  );
}

function kartaSesji(sesja: Session, idTresci: string): HTMLElement {
  return el(
    'div',
    {
      klasa: 'dn-karta-widoku',
      role: 'tab',
      tabindex: '0',
      'aria-selected': 'false',
      'data-karta': sesja.id,
      'data-karta-rodzaj': 'modul',
      'aria-controls': idTresci,
    },
    [
      znak(ikony.karty, 'dn-karta-widoku-ikona'),
      el('span', { klasa: 'dn-karta-widoku-nazwa', tekst: sesja.title ?? sesja.id }),
      el('span', { klasa: 'dn-karta-widoku-zamknij dn-etykietka', 'data-etykietka': tekst('belka.zamknij'), 'aria-hidden': true }, [
        znak(ikony.zamknij),
      ]),
    ],
  );
}

export function pasmoKart(w: WlasciwosciPasma): HTMLElement {
  const lista = el(
    'div',
    { klasa: 'dn-karty-lista', role: 'tablist', 'aria-label': tekst('narzedzia.karty') },
    [kartaGlowna(w.idTresci), ...w.sesje.map((sesja) => kartaSesji(sesja, w.idTresci))],
  );

  return el('div', { klasa: 'dn-karty-pasmo', 'data-gestosc': 'zwarta' }, [lista]);
}
