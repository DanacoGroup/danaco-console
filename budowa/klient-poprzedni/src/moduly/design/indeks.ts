import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { KOD_MODULU } from './etykiety-designu';
import { utworzModulDesign } from './modul-design';

/**
 * Moduł Design — punkt zbiorczy katalogu.
 *
 * Powłoka zna stąd jedną rzecz: opis modułu. Kod stoi wewnątrz opisu, więc
 * moduł sam mówi, którym modułem jest, a rozjazd między nazwą w rejestrze
 * a rzeczywistością jest niemożliwy.
 *
 * Moduł nie osadza się sam w dokumencie i nie zna powłoki — oddaje element,
 * a warstwa składająca decyduje, gdzie go postawić.
 *
 * Rozłączenie zostaje w module: `WidokModulu` powłoki nie ma czynności
 * odpięcia, a subskrypcje kanału odpina `rozlacz()` widoku modułu.
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
