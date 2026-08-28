// Źródła SVG — grupa: nawigacja.
// Nazwy i kolejność wprost z `ikony/manifest.json` (pozycje 1–13).
// Każdy import wciąga treść pliku znacznikiem `?raw`, więc wykaz niesie
// znaczniki SVG, a nie ścieżki do plików.

import dom from '../svg/dom.svg?raw';
import menu from '../svg/menu.svg?raw';
import strzalkaLewo from '../svg/strzalka-lewo.svg?raw';
import strzalkaPrawo from '../svg/strzalka-prawo.svg?raw';
import grotDol from '../svg/grot-dol.svg?raw';
import grotGora from '../svg/grot-gora.svg?raw';
import grotPrawo from '../svg/grot-prawo.svg?raw';
import wiecej from '../svg/wiecej.svg?raw';
import linkZewnetrzny from '../svg/link-zewnetrzny.svg?raw';
import szukaj from '../svg/szukaj.svg?raw';
import filtr from '../svg/filtr.svg?raw';
import zamknij from '../svg/zamknij.svg?raw';
import kartaOkna from '../svg/karta-okna.svg?raw';

/**
 * Nawigacja i orientacja — ruch po widokach, karty, okna. Kluczem wykazu jest
 * nazwa ikony, wartością treść pliku SVG. Zapis `as const` utrwala zbiór nazw
 * w typie, więc odwołanie do nazwy spoza grupy nie przechodzi budowy.
 */
export const NAWIGACJA = {
  'dom': dom,
  'menu': menu,
  'strzalka-lewo': strzalkaLewo,
  'strzalka-prawo': strzalkaPrawo,
  'grot-dol': grotDol,
  'grot-gora': grotGora,
  'grot-prawo': grotPrawo,
  'wiecej': wiecej,
  'link-zewnetrzny': linkZewnetrzny,
  'szukaj': szukaj,
  'filtr': filtr,
  'zamknij': zamknij,
  'karta-okna': kartaOkna,
} as const;
