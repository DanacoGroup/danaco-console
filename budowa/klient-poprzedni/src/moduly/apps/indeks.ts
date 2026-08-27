import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { utworzModulApps } from './modul-apps';
import { KOD_MODULU } from './zrodlo-okna-modulu';

/**
 * Punkt zbiorczy katalogu modułu Apps. Jedynym eksportem widzianym z zewnątrz
 * jest opis modułu, którego kod bierze się ze stałej przy źródle okna modułu,
 * odsiewającego tym samym kodem okna sesji.
 */
export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  utworzWidok(kanal) {
    const modul = utworzModulApps(kanal);
    return {
      element: modul.element,
      wczytaj: (idSesji) => modul.wczytaj(idSesji),
      zamknij: () => modul.rozlacz(),
    };
  },
};
