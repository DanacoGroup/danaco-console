import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { utworzModulAssistant } from './modul-assistant';

/**
 * Moduł Assistant — punkt zbiorczy katalogu.
 *
 * Powłoka zna stąd jeden opis modułu, a kod modułu stoi wewnątrz tego opisu.
 * Moduł nie osadza się sam w dokumencie: oddaje element, a warstwa składająca
 * decyduje, gdzie go postawić.
 */
export const MODUL: OpisModulu = {
  kod: 'assistant',
  utworzWidok(kanal) {
    const modul = utworzModulAssistant(kanal);
    return {
      element: modul.element,
      wczytaj: (idSesji) => modul.wczytaj(idSesji),
    };
  },
};
