/**
 * SKŁADNIK — TOR TRZECH KROKÓW ODZYSKIWANIA.
 *
 * Adres → potwierdzenie → nowe hasło. Stan niosą dwie dane: `data-biezacy`
 * i `data-zrobiony`. Znak kroku zrobionego rysuje arkusz przez `::after` —
 * numer ustępuje ptaszkowi, więc stan nie stoi na samej barwie (WCAG 1.4.1).
 * Wstawianie tu rysunku dublowałoby ten sam znak.
 *
 * Strzałki między krokami są rysunkiem, nie treścią — czytnik ekranu je pomija.
 */

import { el, wykaz, type Dziecko } from '../narzedzia.ts';

/** Znak rozdzielający kroki; rysunek, nie treść. */
const STRZALKA = '→';

export interface WlasciwosciKrokow {
  /** Numer kroku bieżącego, liczony od jedynki. */
  biezacy: number;
}

export function krokiOdzyskiwania(w: WlasciwosciKrokow): HTMLElement {
  const nazwy = wykaz('dostep.odzyskiwanie.kroki');
  const dzieci: Dziecko[] = [];
  nazwy.forEach((nazwa, i) => {
    const numer = i + 1;
    dzieci.push(
      el(
        'span',
        {
          klasa: 'au-krok',
          dane: {
            biezacy: numer === w.biezacy ? 'tak' : 'nie',
            zrobiony: numer < w.biezacy ? 'tak' : null,
          },
        },
        [el('span', { klasa: 'au-krok-nr', tekst: String(numer) }), nazwa],
      ),
    );
    if (numer < nazwy.length) {
      dzieci.push(
        el('span', { klasa: 'au-krok-strzalka', 'aria-hidden': 'true', tekst: STRZALKA }),
      );
    }
  });
  return el('div', { klasa: 'au-kroki' }, dzieci);
}
