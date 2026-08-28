import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { KOD_MODULU } from './etykiety-designu';
import { utworzModulDesign } from './modul-design';

/**
 * Punkt zbiorczy katalogu modułu Design. Powłoka zna stąd jedną rzecz: opis
 * modułu, w którym kod stoi wewnątrz opisu, więc moduł sam mówi, którym modułem
 * jest, a rozjazd z nazwą w rejestrze jest niemożliwy.
 */
export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  utworzWidok(kanal) {
    const modul = utworzModulDesign(kanal);
    return {
      element: modul.element,
      wczytaj: (idSesji) => modul.wczytaj(idSesji),
    };
  },
};

export { utworzModulDesign } from './modul-design';
export type { ModulDesign } from './modul-design';
export { utworzStanDesignu } from './stan-designu';
export type { StanDesignu } from './stan-designu';
