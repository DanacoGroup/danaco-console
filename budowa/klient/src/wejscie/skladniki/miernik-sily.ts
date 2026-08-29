/**
 * Składnik — miernik siły hasła. Cztery odcinki toru, zdanie o wyniku
 * i cztery warunki; znak warunku różni się kształtem, nie samą barwą.
 */

import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';
import { DLUGOSC_HASLA } from '../przebieg.ts';

/** Cztery warunki, w kolejności wyświetlania, oceniane osobno przez przebieg i pokazywane tym miernikiem. */
export const WARUNKI_HASLA = ['dlugosc', 'wielkosc', 'cyfra', 'znak'] as const;

/**
 * Warunki zalecane, nie wymagane. Ich brak nie wstrzymuje rejestracji, więc
 * miernik nie zaznacza ich na czerwono — inaczej Operator szukałby usterki
 * w haśle, które program przyjmie.
 */
const WARUNKI_ZALECANE = new Set<string>(['znak']);

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
      {
        klasa: 'au-sila-warunek',
        dane: { warunek: nazwa, spelniony: 'nie', zalecany: WARUNKI_ZALECANE.has(nazwa) ? 'tak' : null },
      },
      [
        kolko,
        ptaszek,
        tekst(`dostep.sila.warunki.${nazwa}`, { znaki: DLUGOSC_HASLA }),
        WARUNKI_ZALECANE.has(nazwa)
          ? el('span', { klasa: 'dn-meta', tekst: tekst('dostep.sila.zalecany') })
          : null,
      ],
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
