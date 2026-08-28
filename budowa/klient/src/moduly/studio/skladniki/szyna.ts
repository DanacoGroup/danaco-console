/**
 * Strefa 3 — szyna dokumentów sesji. Rdzeń nie wystawia dziś wykazu dokumentów
 * sesji ani grupowania po projekcie, więc lista stoi nazwanym stanem pustym
 * zamiast wartości zmyślonej po stronie klienta — ten sam chwyt co w
 * `rama/skladniki/panel-sesje.ts` dla zakładki „Projekty".
 */

import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

function znak(rysunek: NazwaZnaku): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

export function szynaDokumentow(): HTMLElement {
  const naglowek = el('div', { klasa: 'st-szyna-naglowek' }, [
    el('button', { klasa: 'dn-btn dn-btn--zarys dn-btn--sm', type: 'button', tekst: tekst('szyna.nowyDokument') }),
    el(
      'button',
      { klasa: 'dn-btn-ikona dn-etykietka', type: 'button', 'data-etykietka': tekst('szyna.filtry'), 'aria-label': tekst('szyna.filtry') },
      [znak('filtr')],
    ),
  ]);

  const lista = el('div', { klasa: 'st-szyna-lista' }, [
    el('p', { klasa: 'dn-pusty-stan dn-pusty-stan--zwarty', tekst: tekst('szyna.brakDokumentow') }),
  ]);

  return el('nav', { klasa: 'st-szyna', 'aria-label': tekst('szyna.etykieta') }, [naglowek, lista]);
}
