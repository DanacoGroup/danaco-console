/**
 * SKŁADNIK — MIERNIK SIŁY HASŁA.
 *
 * Cztery odcinki toru, zdanie o wyniku i cztery warunki. Znak warunku różni
 * się KSZTAŁTEM, nie samą barwą: niespełniony niesie puste kółko, spełniony
 * ptaszka. Ptaszek w każdym stanie czytał się jako „zrobione”, a barwa jako
 * jedyna różnica łamie WCAG 1.4.1.
 *
 * Miernik wiąże się z polem przez `data-sila-dla`, nie przez sąsiedztwo
 * w drzewie — sąsiedztwo bywa różne w różnych oknach i cicho się rozjeżdża.
 * Regułę oceny niesie przebieg; miernik jej nie powtarza.
 */

import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';
import { DLUGOSC_HASLA } from '../przebieg.ts';

/** Cztery warunki, w kolejności wyświetlania. */
export const WARUNKI_HASLA = ['dlugosc', 'wielkosc', 'cyfra', 'znak'] as const;

export interface WlasciwosciMiernika {
  /** Identyfikator pola hasła, które ten miernik ocenia. */
  dla: string;
}

export function miernikSily(w: WlasciwosciMiernika): HTMLElement {
  const odcinki = WARUNKI_HASLA.map(() => el('span', { klasa: 'au-sila-odc' }));

  const warunki = WARUNKI_HASLA.map((nazwa) => {
    const kolko = zeZnacznika(ikony.kolko);
    kolko.setAttribute('class', 'ik-brak');
    kolko.setAttribute('aria-hidden', 'true');
    const ptaszek = zeZnacznika(ikony.ptaszek);
    ptaszek.setAttribute('class', 'ik-jest');
    ptaszek.setAttribute('aria-hidden', 'true');
    return el(
      'span',
      { klasa: 'au-sila-warunek', dane: { warunek: nazwa, spelniony: 'nie' } },
      [kolko, ptaszek, tekst(`dostep.sila.warunki.${nazwa}`, { znaki: DLUGOSC_HASLA })],
    );
  });

  return el('div', { klasa: 'au-sila', dane: { stopien: '0', 'sila-dla': w.dla } }, [
    el('div', { klasa: 'au-sila-tor', 'aria-hidden': 'true' }, odcinki),
    el('span', {
      klasa: 'au-sila-opis',
      id: `${w.dla}-sila-opis`,
      'aria-live': 'polite',
      tekst: tekst('dostep.sila.puste'),
      dane: { 'sila-opis': true },
    }),
    el('div', { klasa: 'au-sila-warunki' }, warunki),
  ]);
}
