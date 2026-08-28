/**
 * Rama aplikacji — narzędzia pod składnikami. Ten sam kształt co w drodze
 * wejścia: budowanie węzła i sięganie po łańcuch z własnego katalogu treści.
 */

import { tresci, type WezelTresci } from './tresci.ts';

/** Wartość atrybutu węzła: null i false pomijają atrybut całkowicie, true stawia go pusty, bez wartości. */
export type WartoscAtrybutu = string | number | boolean | null | undefined;

/** Atrybuty węzła przekazywane budowie; nazwy własne — klasa, tekst — reszta idzie wprost jako atrybut HTML. */
export interface AtrybutyWezla {
  klasa?: string | null;
  tekst?: string | null;
  [nazwa: string]: WartoscAtrybutu | null | undefined;
}

/** Dziecko węzła: sam węzeł, łańcuch tekstu albo nic, gdy warunek składnika nie chce go wcale wstawiać. */
export type Dziecko = Node | string | null | undefined | false;

/**
 * Budowa węzła. Atrybuty rozpoznawane po nazwie: klasa, tekst, reszta wprost
 * jako atrybut. Wartość null albo false pomija atrybut.
 */
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
  // Wpis diagnostyczny dla wykonawcy, nie tekst okna: pusty znacznik znaku jest usterką zestawu.
  if (pierwszy === null) throw new Error('[rama] pusty znacznik znaku');
  return pierwszy as SVGElement;
}

/**
 * Sięgnięcie po treść ścieżką kluczy: `tekst('stan.srodowisko')`.
 *
 * Brak klucza nie jest sytuacją do obsłużenia po cichu — wraca sama ścieżka
 * w nawiasach kątowych, żeby usterka katalogu rzucała się w oczy w oknie.
 */
export function tekst(sciezka: string): string {
  let biezacy: unknown = tresci;
  for (const czlon of sciezka.split('.')) {
    if (biezacy === null || typeof biezacy !== 'object') return `⟨${sciezka}⟩`;
    biezacy = (biezacy as Record<string, WezelTresci>)[czlon];
  }
  if (typeof biezacy !== 'string') return `⟨${sciezka}⟩`;
  return biezacy;
}
