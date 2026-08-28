import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { utworzRameOkna, type RolaOkna } from '../../komponenty/rama-okna';
import { utworzStanOknaStudio, type StanOknaStudio } from './stan-okna-studio';

/** Okno operacyjne Studio: biblioteczna rama okna złączona z pasem stanu modułu, wraz z wymaganym dymkiem objaśnienia. */
export interface OknoStudio {
  /** Sekcja osadzana w pasie układu modułu. */
  element: HTMLElement;
  /** Pas stanu wraz z miejscem na treść okna. */
  stan: StanOknaStudio;
  /** Pasek nagłówka — miejsce na sterowanie okna. */
  pasek: HTMLElement;
}

/** Kod okna w katalogu rdzenia oraz jego tytuł, rola i objaśnienie wyświetlane operatorowi przy otwarciu. */
export interface OpisOkna {
  kod: string;
  tytul: string;
  rola: RolaOkna;
  objasnienie: string;
}

export function utworzOknoStudio(opis: OpisOkna): OknoStudio {
  const stan = utworzStanOknaStudio();
  const dymek = utworzDymekObjasnienia(opis.objasnienie, {
    powloka: 'ms-dymek',
    znak: 'ms-dymek__znak',
  });

  const rama = utworzRameOkna({
    tytul: opis.tytul,
    rola: opis.rola,
    kod: opis.kod,
    modul: 'Studio',
    dodatkiNaglowka: [dymek],
    // Przedrostek ms jest uchwytem reguły znaczącej okno wskazane z paska zaznaczenia edytora.
    przedrostek: 'ms',
  });
  rama.cialo.append(stan.element);

  return { element: rama.element, stan, pasek: rama.pasek };
}
