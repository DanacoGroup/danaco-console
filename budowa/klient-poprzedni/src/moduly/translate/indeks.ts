import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { KOD_MODULU } from './katalog-okien-translate';
import { utworzModulTranslate } from './modul-translate';

/**
 * Moduł Translate jest punktem zbiorczym katalogu: oddaje opis modułu i element, nie osadzając się sam w dokumencie.
 */
export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  utworzWidok(kanal) {
    const modul = utworzModulTranslate(kanal);
    return {
      element: modul.element,
      wczytaj: (idSesji) => modul.wczytaj(idSesji),
      zamknij: () => modul.rozlacz(),
    };
  },
};

export { utworzModulTranslate } from './modul-translate';
export type { ModulTranslate } from './modul-translate';
