/**
 * Strefa 2 — wstążka narzędziowa okna roboczego. Leży wewnątrz bryły, w jednej
 * barwie z kartą bieżącą, na pełną szerokość okna — odrębna od wstążki
 * aplikacji, która niesie wyłącznie narzędzia poziomu aplikacji (rama).
 *
 * Grupy „Operacje modułu" z prototypu ([]Uruchom operację[], []Wyślij do
 * Library[]) tu nie stoją — wołałyby zakres pracy, którego ten teren nie
 * stawia (panele). Tabliczka sesji niesie tytuł sesji bieżącej, gdy rdzeń
 * go podał; w przeciwnym razie nazwany stan pusty, nie wymyślona nazwa.
 */

import type { Session } from '../../../../../shared/contract.ts';
import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciWstazki {
  /** Karty sesji odtworzone przez rdzeń — pierwsza nazywa tabliczkę sesji. */
  sesje: Session[];
}

function znak(rysunek: NazwaZnaku): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function przyciskIkony(rysunek: NazwaZnaku, etykieta: string, atrybuty?: Record<string, string>): HTMLElement {
  return el(
    'button',
    { klasa: 'dn-nrz-btn dn-etykietka', type: 'button', 'data-etykietka': etykieta, 'aria-label': etykieta, ...atrybuty },
    [znak(rysunek)],
  );
}

function grupa(etykieta: string, dzieci: HTMLElement[]): HTMLElement {
  return el('span', { klasa: 'st-wstazka-grupa', role: 'group', 'aria-label': etykieta }, dzieci);
}

export function wstazkaOkna(w: WlasciwosciWstazki): HTMLElement {
  const nazwaSesji = w.sesje[0]?.title ?? tekst('wstazka.bezSesji');

  return el('div', { klasa: 'st-wstazka', role: 'group', 'aria-label': tekst('wstazka.etykieta') }, [
    el('span', { klasa: 'st-wstazka-sesja' }, [znak('olowek'), el('span', { tekst: nazwaSesji })]),
    grupa(tekst('wstazka.ukladEtykieta'), [
      przyciskIkony('dymek', tekst('wstazka.oknoKomunikacji'), { 'aria-pressed': 'true' }),
      przyciskIkony('podzial', tekst('wstazka.podzialPionowy')),
    ]),
    grupa(tekst('wstazka.wersjeEtykieta'), [
      przyciskIkony('zapisz', tekst('wstazka.zapiszWersje')),
      przyciskIkony('diff', tekst('wstazka.porownajWersje')),
      przyciskIkony('oko', tekst('wstazka.podgladWydruku')),
    ]),
    el('span', { klasa: 'st-wstazka-odstep' }),
    przyciskIkony('dostosuj', tekst('wstazka.dostosuj')),
  ]);
}
