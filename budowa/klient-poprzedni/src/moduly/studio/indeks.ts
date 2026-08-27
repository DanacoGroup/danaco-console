import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { utworzModulStudio } from './modul-studio';
import { KOD_MODULU } from './zrodlo-akcji-studio';

/** Moduł Studio, punkt zbiorczy katalogu: powłoka zna stąd tylko opis modułu, a kod stoi wewnątrz opisu, nie w mapie rejestru. */
export const MODUL: OpisModulu = {
  kod: KOD_MODULU,

  utworzWidok(kanal) {
    const modul = utworzModulStudio(kanal);
    return {
      element: modul.element,
      wczytaj: (idSesji) => modul.wczytaj(idSesji),
      zamknij: () => modul.rozlacz(),
    };
  },
};

export { utworzModulStudio } from './modul-studio';
export type { ModulStudio } from './modul-studio';
