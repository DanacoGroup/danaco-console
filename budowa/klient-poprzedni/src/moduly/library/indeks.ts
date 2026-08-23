import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { utworzModulLibrary } from './modul-library';

/**
 * Moduł Library — punkt zbiorczy katalogu.
 *
 * Powłoka zna stąd jeden byt: opis modułu wraz z kodem katalogu rdzenia
 * i wytwórnią widoku. Moduł nie osadza się sam w dokumencie i nie zna powłoki —
 * oddaje element, a warstwa składająca decyduje, gdzie go postawić.
 *
 * Kod `library` odpowiada kolumnie `modul.kod` w rdzeniu; wpis do rejestru
 * wiąże widok z modułem po tym kodzie.
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
