/**
 * Składnik — tor trzech kroków odzyskiwania: adres, potwierdzenie, nowe
 * hasło. Stan niosą dwie dane, a znak kroku zrobionego rysuje arkusz stylu.
 */

import { el, wykaz, type Dziecko } from '../narzedzia.ts';

/** Znak rozdzielający kroki na torze; jest rysunkiem, nie treścią, pomijanym całkowicie przez czytnik ekranu. */
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
