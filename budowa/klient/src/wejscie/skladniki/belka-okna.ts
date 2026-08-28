/**
 * Składnik — belka okna systemowego. Okno przed uwierzytelnieniem nie ma
 * ramy aplikacji: ma wyłącznie belkę — znak, tytuł i trzy kontrolki okna.
 */

import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciBelki {
  /** Klucz katalogu — tytuł okna w belce. */
  tytul: string;
}

function kontrolka(znak: NazwaZnaku, etykieta: string, odmiana?: string): HTMLElement {
  const rysunek = zeZnacznika(ikony[znak]);
  rysunek.setAttribute('aria-hidden', 'true');
  const klasa = odmiana === undefined
    ? 'dn-okno-wejsciowe-belka-btn'
    : `dn-okno-wejsciowe-belka-btn ${odmiana}`;
  return el('button', { klasa, type: 'button', 'aria-label': tekst(etykieta) }, [rysunek]);
}

export function belkaOkna(w: WlasciwosciBelki): HTMLElement {
  const znak = zeZnacznika(ikony.godlo);
  znak.setAttribute('class', 'dn-okno-wejsciowe-belka-znak');
  znak.setAttribute('aria-hidden', 'true');

  /* Uchwyt rozpiera belkę między tytułem a sterowaniem — biblioteka nie ma
     grupy kontrolek, bo pozycję prawą nadaje właśnie ten pusty człon. */
  return el('div', { klasa: 'dn-okno-wejsciowe-belka' }, [
    znak,
    el('p', { klasa: 'dn-okno-wejsciowe-belka-tytul', tekst: tekst(w.tytul) }),
    el('div', { klasa: 'dn-okno-wejsciowe-belka-uchwyt' }),
    kontrolka('zwin', 'okno.zwin'),
    kontrolka('rozwin', 'okno.rozwin'),
    kontrolka('zamknij', 'okno.zamknij', 'dn-okno-wejsciowe-belka-btn--zamknij'),
  ]);
}
