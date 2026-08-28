/**
 * Strefa 1 — belka tytułowa. Znak i nazwa produktu po lewej, tytuł bieżącego
 * widoku pośrodku. Sterowanie oknem systemowym należy do powłoki Tauri —
 * poza tym terenem — więc belka go nie niesie.
 */

import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciBelki {
  /** Nazwa środowiska, do którego przebieg wszedł — jedyny fragment tytułu, który nie stoi w katalogu treści. */
  srodowisko: string;
}

/** Belka wraz z odrębnym uchwytem do węzła tytułu — wywołujący zmienia tytuł bez ponownego odpytywania drzewa. */
export interface Belka {
  wezel: HTMLElement;
  tytul: HTMLElement;
}

export function belka(w: WlasciwosciBelki): Belka {
  const znak = zeZnacznika(ikony.godlo);
  znak.setAttribute('aria-hidden', 'true');

  const tytul = el('span', {
    klasa: 'dn-belka-tytul',
    tekst: `${tekst('belka.marka')} ${tekst('belka.separator')} ${w.srodowisko}`,
    'data-belka-tytul': true,
  });

  const wezel = el('header', { klasa: 'dn-belka' }, [
    el('span', { klasa: 'dn-belka-marka' }, [
      znak,
      el('span', { klasa: 'dn-belka-nazwa', tekst: tekst('belka.marka') }),
    ]),
    tytul,
  ]);

  return { wezel, tytul };
}
