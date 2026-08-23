import type { OpisModulu } from '../../aplikacja/rejestr-modulow';
import { utworzModulApps } from './modul-apps';
import { KOD_MODULU } from './zrodlo-okna-modulu';

/**
 * Moduł Apps — punkt zbiorczy katalogu.
 *
 * Jedynym eksportem widzianym z zewnątrz jest opis modułu. Moduł nie osadza się
 * sam w dokumencie: oddaje element, a warstwa składająca decyduje, gdzie go
 * postawić — dzięki temu te same okna wchodzą i w obszar roboczy powłoki,
 * i w stanowisko sprawdzianu.
 *
 * Kod modułu bierze się ze stałej przy źródle okna modułu, które tym samym
 * kodem odsiewa okna sesji.
 *
 * `zamknij` podpina `rozlacz()` złożenia: bez niego subskrypcje stanu
 * (`apps.build.changed`, `progress.changed`) i subskrypcje okien pomocniczych
 * żyją dalej po zejściu modułu ze sceny. Pole jest w `WidokModulu` opcjonalne,
 * więc kompilator braku nie zgłosi.
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
