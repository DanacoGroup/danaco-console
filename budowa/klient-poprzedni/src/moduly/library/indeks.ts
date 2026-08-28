import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { utworzModulLibrary } from './modul-library';

/**
 * Punkt zbiorczy katalogu modułu Library. Powłoka zna stąd jeden byt: opis
 * modułu wraz z kodem katalogu rdzenia i wytwórnią widoku. Kod `library`
 * odpowiada kolumnie `modul.kod` w rdzeniu.
 */
export const MODUL: OpisModulu = {
  kod: 'library',
  utworzWidok(kanal) {
    const modul = utworzModulLibrary(kanal);
    return {
      element: modul.element,
      wczytaj: (idSesji) => modul.wczytaj(idSesji),
    };
  },
};

export { utworzModulLibrary } from './modul-library';
export type { ModulLibrary } from './modul-library';
export { utworzStanBiblioteki } from './stan-biblioteki';
export type { StanBiblioteki } from './stan-biblioteki';
export { utworzZrodloBiblioteki } from './zrodlo-biblioteki';
export type { ZrodloBiblioteki } from './zrodlo-biblioteki';
