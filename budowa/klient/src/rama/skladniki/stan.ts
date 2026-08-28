/**
 * Strefa 3 — pas stanu. Środowisko bieżące, liczba kart sesji odtworzonych
 * przez rdzeń i motyw czynny — trzy rzeczy, które przebieg już zna po
 * wejściu do środowiska, bez pytania rdzenia o nic dodatkowego.
 */

import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciStanu {
  srodowisko: string;
  liczbaSesji: number;
  motywCiemny: boolean;
}

function znak(rysunek: string): SVGElement {
  const wezel = zeZnacznika(rysunek);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

function pozycja(rysunek: string, dzieci: (Node | string)[]): HTMLElement {
  return el('span', { klasa: 'dn-stan-poz' }, [znak(rysunek), ...dzieci]);
}

export function stan(w: WlasciwosciStanu): HTMLElement {
  return el(
    'div',
    { klasa: 'dn-stan', role: 'status', 'aria-label': tekst('stan.etykieta') },
    [
      pozycja(ikony.godlo, [
        el('b', { tekst: tekst('stan.srodowisko') }),
        ` · ${w.srodowisko}`,
      ]),
      el('span', { klasa: 'dn-stan-sep', 'aria-hidden': true }),
      pozycja(ikony.karty, [
        `${tekst('stan.sesje')}: `,
        el('b', { 'data-stan-sesje': true, tekst: String(w.liczbaSesji) }),
      ]),
      el(
        'span',
        { klasa: 'dn-stan-poz dn-stan-poz--prawa' },
        [
          znak(w.motywCiemny ? ikony.ksiezyc : ikony.slonce),
          `${tekst('stan.widok')} `,
          el('b', {
            'data-stan-motyw': true,
            tekst: tekst(w.motywCiemny ? 'stan.motywCiemny' : 'stan.motywJasny'),
          }),
        ],
      ),
    ],
  );
}
