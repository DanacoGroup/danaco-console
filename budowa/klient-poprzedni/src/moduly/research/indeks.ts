import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { utworzModulResearch, type ModulResearch } from './modul-research';

/**
 * Moduł Research — punkt zbiorczy katalogu.
 *
 * Powłoka zna stąd jedną rzecz: opis modułu. Kod modułu stoi wewnątrz opisu,
 * bo moduł sam mówi, którym jest modułem — rozjazd między nazwą w rejestrze
 * a rzeczywistością staje się przez to niemożliwy.
 *
 * Moduł nie osadza się sam w dokumencie i nie zna powłoki: oddaje element,
 * a warstwa składająca decyduje, gdzie go postawić. Dzięki temu te same okna
 * wchodzą i w obszar roboczy powłoki, i w stanowisko sprawdzianu.
 *
 * `zamknij` podpina `rozlacz()` modułu: bez niego subskrypcja stanu badania
 * i kreator raportu żyją dalej po zejściu modułu ze sceny. Pole jest
 * w `WidokModulu` opcjonalne, więc kompilator braku nie zgłosi.
 */
export const MODUL: OpisModulu = {
  kod: 'research',

  utworzWidok(kanal) {
    const modul: ModulResearch = utworzModulResearch(kanal);
    return {
      element: modul.element,
      wczytaj: (idSesji) => modul.wczytaj(idSesji),
      zamknij: () => modul.rozlacz(),
    };
  },
};

export { utworzModulResearch } from './modul-research';
export type { ModulResearch } from './modul-research';
export { utworzStanBadania } from './stan-badania';
export type { StanBadania } from './stan-badania';
