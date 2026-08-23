import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { KOD_MODULU } from './etykiety-browser';
import { utworzModulBrowser } from './modul-browser';

/**
 * Moduł Browser — punkt zbiorczy katalogu.
 *
 * Powłoka zna stąd jedną rzecz: opis modułu. Kod stoi wewnątrz opisu, nie
 * w mapie rejestru, więc rozjazd między nazwą w rejestrze a rzeczywistością
 * jest niemożliwy.
 *
 * Moduł nie osadza się sam w dokumencie: oddaje element, a warstwa składająca
 * decyduje, gdzie go postawić. Dzięki temu te same okna wchodzą i w obszar
 * roboczy powłoki, i w podgląd sprawdzianu.
 */
export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  // Wytwórnia oddaje moduł wprost: `ModulBrowser` niesie `element` i `wczytaj`
  // żądane przez powłokę oraz `rozlacz`, po które sięga sprawdzian modułu.
  utworzWidok: utworzModulBrowser,
};

export { utworzModulBrowser } from './modul-browser';
export type { ModulBrowser } from './modul-browser';
