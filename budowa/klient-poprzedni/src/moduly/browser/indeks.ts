import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { KOD_MODULU } from './etykiety-browser';
import { utworzModulBrowser } from './modul-browser';

/**
 * Punkt zbiorczy katalogu modułu Browser. Powłoka zna stąd jedną rzecz: opis
 * modułu, w którym kod stoi wewnątrz opisu, a nie w mapie rejestru, dlatego
 * rozjazd między nazwą w rejestrze a rzeczywistością jest niemożliwy.
 */
export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  // Wytwórnia oddaje moduł wprost: `ModulBrowser` niesie `element`, `wczytaj`
  // oraz `rozlacz`.
  utworzWidok: utworzModulBrowser,
};

export { utworzModulBrowser } from './modul-browser';
export type { ModulBrowser } from './modul-browser';
