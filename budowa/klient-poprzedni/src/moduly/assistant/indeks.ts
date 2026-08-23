import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { utworzModulAssistant } from './modul-assistant';

/**
 * Moduł Assistant — punkt zbiorczy katalogu.
 *
 * Powłoka zna stąd jeden opis modułu. Kod modułu stoi wewnątrz opisu, bo moduł
 * sam mówi, którym jest modułem — rozjazd między nazwą w rejestrze
 * a rzeczywistością staje się przez to niemożliwy
 * (`aplikacja/rejestr-modulow.ts`).
 *
 * Moduł nie osadza się sam w dokumencie i nie zna powłoki: oddaje element,
 * a warstwa składająca decyduje, gdzie go postawić i czy wpisać go do rejestru.
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
