// Źródła SVG — grupa: czynnosci.
// Nazwy i kolejność wprost z `ikony/manifest.json` (pozycje 14–29).
// Każdy import wciąga treść pliku znacznikiem `?raw`, więc wykaz niesie
// znaczniki SVG, a nie ścieżki do plików.

import plus from '../svg/plus.svg?raw';
import olowek from '../svg/olowek.svg?raw';
import kosz from '../svg/kosz.svg?raw';
import pobierz from '../svg/pobierz.svg?raw';
import wgraj from '../svg/wgraj.svg?raw';
import wyslij from '../svg/wyslij.svg?raw';
import odswiez from '../svg/odswiez.svg?raw';
import odpowiedz from '../svg/odpowiedz.svg?raw';
import uruchom from '../svg/uruchom.svg?raw';
import zatrzymaj from '../svg/zatrzymaj.svg?raw';
import wstrzymaj from '../svg/wstrzymaj.svg?raw';
import spinacz from '../svg/spinacz.svg?raw';
import ustawienia from '../svg/ustawienia.svg?raw';
import oko from '../svg/oko.svg?raw';
import gwiazdka from '../svg/gwiazdka.svg?raw';
import kopiuj from '../svg/kopiuj.svg?raw';

/**
 * Czynności operatora — tworzenie, zmiana, przesył, sterowanie. Kluczem wykazu
 * jest nazwa ikony, wartością treść pliku SVG. Zapis `as const` utrwala zbiór
 * nazw w typie, więc odwołanie do nazwy spoza grupy nie przechodzi budowy.
 */
export const CZYNNOSCI = {
  'plus': plus,
  'olowek': olowek,
  'kosz': kosz,
  'pobierz': pobierz,
  'wgraj': wgraj,
  'wyslij': wyslij,
  'odswiez': odswiez,
  'odpowiedz': odpowiedz,
  'uruchom': uruchom,
  'zatrzymaj': zatrzymaj,
  'wstrzymaj': wstrzymaj,
  'spinacz': spinacz,
  'ustawienia': ustawienia,
  'oko': oko,
  'gwiazdka': gwiazdka,
  'kopiuj': kopiuj,
} as const;
