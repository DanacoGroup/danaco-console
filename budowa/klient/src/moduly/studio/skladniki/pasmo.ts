/**
 * Strefa 1 — pasmo siedmiu kart okna roboczego. Znak modułu w narożniku,
 * wykaz kart `role="tablist"` i sterowanie pasma poza wykazem — rola
 * `tablist` nie przyjmuje innych dzieci. Przełączanie kart wiąże `montaz.ts`,
 * bo dotyczy też widoczności paneli w zaczepie, poza tym składnikiem.
 */

import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';
import { KARTA_EDITOR, karty, type DefinicjaKarty } from './definicje.ts';

function znak(rysunek: NazwaZnaku, klasa?: string): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  if (klasa !== undefined) wezel.setAttribute('class', klasa);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function przyciskIkony(rysunek: NazwaZnaku, etykieta: string, atrybuty?: Record<string, string>): HTMLElement {
  return el(
    'button',
    { klasa: 'dn-btn-ikona dn-etykietka', type: 'button', 'data-etykietka': etykieta, 'aria-label': etykieta, ...atrybuty },
    [znak(rysunek)],
  );
}

function karta(k: DefinicjaKarty): HTMLElement {
  const robocza = k.kod === KARTA_EDITOR;
  return el(
    'div',
    {
      klasa: robocza ? 'dn-karta-widoku dn-karta dn-karta--robocza' : 'dn-karta-widoku dn-karta',
      role: 'tab',
      id: `karta-${k.kod}`,
      'aria-controls': `panel-${k.kod}`,
      'aria-selected': robocza ? 'true' : 'false',
      tabindex: robocza ? '0' : '-1',
      'data-karta': k.kod,
    },
    [znak(k.ikona, 'dn-karta-widoku-ikona'), el('span', { klasa: 'dn-karta-widoku-nazwa', tekst: k.nazwa })],
  );
}

export function pasmoKart(): HTMLElement {
  const lista = el(
    'div',
    { klasa: 'st-karty', role: 'tablist', 'aria-label': tekst('pasmo.etykietaKart') },
    karty.map(karta),
  );

  const sterowanie = el('span', { klasa: 'st-pasmo-sterowanie' }, [
    przyciskIkony('plus', tekst('pasmo.nowe')),
    przyciskIkony('maksymalizuj', tekst('pasmo.maksymalizuj'), { 'aria-pressed': 'false' }),
  ]);

  return el('div', { klasa: 'dn-karty-pasmo st-pasmo', 'data-gestosc': 'ciasna' }, [
    el('span', { klasa: 'st-pasmo-znak', 'aria-hidden': 'true' }, [znak('olowek')]),
    lista,
    sterowanie,
  ]);
}
