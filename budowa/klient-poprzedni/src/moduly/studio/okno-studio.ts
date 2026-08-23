import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { utworzRameOkna, type RolaOkna } from '../../komponenty/rama-okna';
import { utworzStanOknaStudio, type StanOknaStudio } from './stan-okna-studio';

/**
 * Okno operacyjne Studio: biblioteczna rama plus pas stanu modułu.
 *
 * Obudowa okna — nagłówek, plakietka roli, gniazda akcji i ciało — jest wspólna
 * całemu drzewu (`komponenty/rama-okna`). Rama nie zna faz okna, bo nie pobiera
 * danych, których cykl życia miałaby pokazywać. Studio potrzebuje jednego
 * i drugiego naraz, więc tutaj zostaje samo zszycie: rama z biblioteki, pas
 * stanu z modułu, jedno wywołanie zamiast powtarzania go w pięciu oknach.
 *
 * Dymek objaśnienia jest wymagany, nie opcjonalny: każde okno modułu mówi, po co
 * jest i którą komendą działa, zanim cokolwiek zostanie naciśnięte.
 */

export interface OknoStudio {
  /** Sekcja osadzana w pasie układu modułu. */
  element: HTMLElement;
  /** Pas stanu wraz z miejscem na treść okna. */
  stan: StanOknaStudio;
  /** Pasek nagłówka — miejsce na sterowanie okna. */
  pasek: HTMLElement;
}

/** Kod okna w katalogu rdzenia (`okno_operacyjne.kod`) i jego opis dla Operatora. */
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
    // Przedrostek `ms` nie wnosi wyglądu — jest uchwytem reguły
    // `.ms-okno[data-ognisko='tak']`, którą Tools Panel znaczy okno wskazane
    // z paska zaznaczenia edytora.
    przedrostek: 'ms',
  });
  rama.cialo.append(stan.element);

  return { element: rama.element, stan, pasek: rama.pasek };
}
