/**
 * Moduł Studio — narzędzia pod składnikami. Ten sam kształt co w `rama/narzedzia.ts`
 * i w drodze wejścia: budowanie węzła i sięganie po łańcuch z własnego katalogu treści.
 */

import { tresci, type WezelTresci } from './tresci.ts';

export type WartoscAtrybutu = string | number | boolean | null | undefined;

export interface AtrybutyWezla {
  klasa?: string | null;
  tekst?: string | null;
  [nazwa: string]: WartoscAtrybutu | null | undefined;
}

export type Dziecko = Node | string | null | undefined | false;

export function el(znacznik: string, atrybuty?: AtrybutyWezla, dzieci?: Dziecko[]): HTMLElement {
  const wezel = document.createElement(znacznik);
  for (const [nazwa, wartosc] of Object.entries(atrybuty ?? {})) {
    if (wartosc === null || wartosc === undefined || wartosc === false) continue;
    if (nazwa === 'klasa') wezel.className = String(wartosc);
    else if (nazwa === 'tekst') wezel.textContent = String(wartosc);
    else if (wartosc === true) wezel.setAttribute(nazwa, '');
    else wezel.setAttribute(nazwa, String(wartosc));
  }
  for (const dziecko of dzieci ?? []) {
    if (dziecko === null || dziecko === undefined || dziecko === false) continue;
    wezel.appendChild(typeof dziecko === 'string' ? document.createTextNode(dziecko) : dziecko);
  }
  return wezel;
}

/**
 * Węzeł zbudowany ze znacznika. Używany wyłącznie dla znaków z zestawu — nigdy
 * dla danych z zewnątrz, bo znacznik z zewnątrz wykonałby się jak treść okna.
 */
export function zeZnacznika(znacznik: string): SVGElement {
  const szablon = document.createElement('template');
  szablon.innerHTML = znacznik.trim();
  const pierwszy = szablon.content.firstElementChild;
  if (pierwszy === null) throw new Error('[moduly/studio] pusty znacznik znaku');
  return pierwszy as SVGElement;
}

/** Sięgnięcie po treść ścieżką kluczy: `tekst('panel.tytul')`. */
export function tekst(sciezka: string): string {
  let biezacy: unknown = tresci;
  for (const czlon of sciezka.split('.')) {
    if (biezacy === null || typeof biezacy !== 'object') return `⟨${sciezka}⟩`;
    biezacy = (biezacy as Record<string, WezelTresci>)[czlon];
  }
  if (typeof biezacy !== 'string') return `⟨${sciezka}⟩`;
  return biezacy;
}
