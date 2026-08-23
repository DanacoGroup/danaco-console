// Źródła SVG — grupa: stany.
// Nazwy i kolejność wprost z `ikony/manifest.json` (pozycje 30–40).

import ptaszek from '../svg/ptaszek.svg?raw';
import ptaszekKolo from '../svg/ptaszek-kolo.svg?raw';
import blad from '../svg/blad.svg?raw';
import ostrzezenie from '../svg/ostrzezenie.svg?raw';
import info from '../svg/info.svg?raw';
import klodka from '../svg/klodka.svg?raw';
import tarcza from '../svg/tarcza.svg?raw';
import zegar from '../svg/zegar.svg?raw';
import historia from '../svg/historia.svg?raw';
import dzwonek from '../svg/dzwonek.svg?raw';
import aktywnosc from '../svg/aktywnosc.svg?raw';

/** Stany i sygnalizacja — potwierdzenie, błąd, ochrona, czas. */
export const STANY = {
  'ptaszek': ptaszek,
  'ptaszek-kolo': ptaszekKolo,
  'blad': blad,
  'ostrzezenie': ostrzezenie,
  'info': info,
  'klodka': klodka,
  'tarcza': tarcza,
  'zegar': zegar,
  'historia': historia,
  'dzwonek': dzwonek,
  'aktywnosc': aktywnosc,
} as const;
