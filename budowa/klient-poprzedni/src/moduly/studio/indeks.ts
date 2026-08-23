import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { utworzModulStudio } from './modul-studio';
import { KOD_MODULU } from './zrodlo-akcji-studio';

/**
 * Moduł Studio — punkt zbiorczy katalogu.
 *
 * Powłoka zna stąd jedną rzecz: opis modułu. Kod stoi wewnątrz opisu, nie
 * w mapie rejestru — moduł sam mówi, którym jest modułem, więc rozjazd między
 * nazwą w rejestrze a rzeczywistością jest niemożliwy.
 *
 * Moduł nie osadza się sam w dokumencie i nie zna powłoki: oddaje element,
 * a warstwa składająca decyduje, gdzie go postawić. Dzięki temu te same okna
 * wchodzą i w obszar roboczy powłoki, i w stanowisko sprawdzianu.
 *
 * `zamknij` oddaje `rozlacz()` modułu, czyli odpięcie subskrypcji
 * `studio.document.changed`. Pole jest w umowie `aplikacja/rejestr-modulow.ts`
 * jako nieobowiązkowe; woła je dziś stanowisko sprawdzianu, a powłoka zawoła,
 * gdy dostanie granicę życia widoku.
 */
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
