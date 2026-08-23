import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { KOD_MODULU } from './katalog-okien-translate';
import { utworzModulTranslate } from './modul-translate';

/**
 * Moduł Translate — punkt zbiorczy katalogu.
 *
 * Powłoka zna stąd jedną rzecz: opis modułu. Kod `translate` pochodzi z kolumny
 * `modul.kod` rdzenia, nie z literału wymyślonego w kliencie — rozjazd między
 * nazwą w mapie a rzeczywistością staje się przez to niemożliwy. Kod stoi
 * w `katalog-okien-translate.ts` razem z kodami okien, bo obie rzeczy są tym
 * samym: oznaczeniami, którymi moduł zgłasza się rdzeniowi.
 *
 * Moduł nie osadza się sam w dokumencie i nie zna powłoki: oddaje element,
 * a warstwa składająca decyduje, gdzie go postawić. Dzięki temu te same okna
 * wchodzą i w obszar roboczy powłoki, i w stanowisko sprawdzianu.
 *
 * Rozłączenie zdejmuje trzy rzeczy naraz — subskrypcję stanu, nasłuch skrótów
 * klawiszowych i wpis wspólnego katalogu okien — więc widok wystawia je jako
 * `zamknij` umowy powłoki, a nie zostawia wywołującemu do złożenia.
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
