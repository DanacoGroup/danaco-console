/**
 * Pasek narzędzi. Stały pas czynności okna nad płótnem roboczym: przełącznik
 * szyny, nawigacja historii, odświeżenie widoku i pole wyszukiwania —
 * znaczniki `dn-narzedzia-pas` / `dn-narzedzia` / `dn-szukaj` czytane przez
 * `rama.js`.
 */

import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

function znak(rysunek: string): SVGElement {
  const wezel = zeZnacznika(rysunek);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function przycisk(rysunek: keyof typeof ikony, etykieta: string, atrybuty?: Record<string, string>): HTMLElement {
  return el(
    'button',
    { klasa: 'dn-nrz-btn dn-etykietka', type: 'button', 'data-etykietka': etykieta, 'aria-label': etykieta, ...atrybuty },
    [znak(ikony[rysunek])],
  );
}

function grupa(dzieci: HTMLElement[], dodatkowaKlasa?: string): HTMLElement {
  return el('div', { klasa: dodatkowaKlasa ? `dn-narzedzia-grupa ${dodatkowaKlasa}` : 'dn-narzedzia-grupa' }, dzieci);
}

export function pasekNarzedzi(): HTMLElement {
  const pole = el('div', { klasa: 'dn-szukaj dn-narzedzia-szukaj' }, [
    znak(ikony.szukaj),
    el('input', { type: 'search', placeholder: tekst('narzedzia.szukaj'), 'aria-label': tekst('narzedzia.szukaj') }),
  ]);

  const pas = el('div', { klasa: 'dn-narzedzia' }, [
    grupa([przycisk('zwinPanel', tekst('narzedzia.zwinPanel'), { 'aria-pressed': 'false', 'data-zwin-szyne': '' })]),
    grupa([
      przycisk('strzalkaLewo', tekst('narzedzia.wstecz')),
      przycisk('strzalkaPrawo', tekst('narzedzia.naprzod')),
    ]),
    grupa([przycisk('odswiez', tekst('narzedzia.odswiez'))]),
    pole,
  ]);

  return el('div', { klasa: 'dn-narzedzia-pas' }, [pas]);
}
