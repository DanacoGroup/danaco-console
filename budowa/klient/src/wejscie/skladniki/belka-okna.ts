/**
 * SKŁADNIK — BELKA OKNA SYSTEMOWEGO.
 *
 * Okno przed uwierzytelnieniem nie ma ramy aplikacji: ani szyny nawigacji, ani
 * wstążki, ani paska stanu. Ma wyłącznie belkę — znak, tytuł i trzy kontrolki
 * okna. Znak w belce dorównuje wielkością kontrolkom, więc niesie kropkę
 * w barwie sygnału: przygaszony monochromat czyta się jak brakująca ikona.
 */

import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciBelki {
  /** Klucz katalogu — tytuł okna w belce. */
  tytul: string;
}

function kontrolka(znak: NazwaZnaku, etykieta: string): HTMLElement {
  const rysunek = zeZnacznika(ikony[znak]);
  rysunek.setAttribute('aria-hidden', 'true');
  return el(
    'button',
    { klasa: 'we-belka-btn', type: 'button', 'aria-label': tekst(etykieta) },
    [rysunek],
  );
}

export function belkaOkna(w: WlasciwosciBelki): HTMLElement {
  const znak = zeZnacznika(ikony.godlo);
  znak.setAttribute('class', 'we-belka-znak');
  znak.setAttribute('aria-hidden', 'true');

  return el('div', { klasa: 'we-belka' }, [
    znak,
    el('p', { klasa: 'we-belka-tytul', tekst: tekst(w.tytul) }),
    el('div', { klasa: 'we-belka-sterowanie' }, [
      kontrolka('zwin', 'okno.zwin'),
      kontrolka('rozwin', 'okno.rozwin'),
      kontrolka('zamknij', 'okno.zamknij'),
    ]),
  ]);
}
