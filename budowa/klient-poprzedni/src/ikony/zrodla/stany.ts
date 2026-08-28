// Źródła rysunków grupy stanów, wczytywane surowym importem `?raw`, z nazwami
// i kolejnością wziętymi wprost z `ikony/manifest.json`. Plik nie przetwarza
// treści rysunku, tylko wystawia ją w postaci, w jakiej leży na dysku.

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

/**
 * Zbiór stanów i sygnalizacji, obejmujący potwierdzenie, błąd, ostrzeżenie,
 * informację, ochronę oraz czas, wystawiony jako mapa nazwy ikony na treść
 * rysunku, po którą sięga rejestr ikon przy składaniu pełnego zestawu.
 */
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
