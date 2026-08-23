import type { FazaOdczytu } from './stan-konfiguracji';

/**
 * Stany odczytu okna konfiguracji: ładowanie i komunikat blokowy błędu.
 *
 * Modal niesie komunikat blokowy nad stopką z akcjami, a wskaźnik odczytu obok
 * niego — nigdy zamiast formularza, bo pola mają pozostać edytowalne przez
 * cały czas.
 *
 * Komunikat nie znika po naciśnięciu: „Spróbuj ponownie" wyzwala odczyt
 * i zostawia zdanie na miejscu, dopóki sytuacja nie ustanie — dopiero udany
 * odczyt je zdejmuje.
 *
 * Stan pusty katalogu należy do panelu kategorii, nie tutaj: mówi o katalogu,
 * a nie o odczycie.
 */
export interface StanyOdczytu {
  /** Pas stanów osadzany między ciałem okna a jego stopką. */
  element: HTMLElement;
  /** Nanosi fazę odczytu ze stanu okna. */
  odswiez(): void;
}

/**
 * Tyle stanu, ile pas naprawdę czyta.
 *
 * Pas nie zna ani katalogu, ani obszarów sesji — pyta wyłącznie o fazę odczytu
 * i o powód niepowodzenia. Zawężenie zależności do tych dwóch czynności czyni
 * z pasa jeden byt dla obu odczytów okna konfiguracji: katalogu ustawień
 * (`config.get`) i konfiguracji obowiązującej (`config.effective.get`).
 * `StanKonfiguracji` spełnia ten kształt bez żadnej zmiany.
 */
export interface ZrodloFazyOdczytu {
  faza(): FazaOdczytu;
  powodNiepowodzenia(): string;
}

/** Zdania pasa; domyślne mówią o katalogu ustawień. */
export interface ZdaniaOdczytu {
  /** Zdanie towarzyszące wskaźnikowi odczytu. */
  odczyt: string;
  /** Początek komunikatu blokowego; powód rdzenia idzie po nim. */
  blad: string;
}

const ZDANIA_KATALOGU: ZdaniaOdczytu = {
  odczyt: 'Czytam katalog ustawień i zapisy konfiguracji…',
  blad: 'Katalog nie dotarł w całości; pola pozostają edytowalne.',
};

export function utworzStanyOdczytu(
  stan: ZrodloFazyOdczytu,
  ponow: () => void,
  zdania: ZdaniaOdczytu = ZDANIA_KATALOGU,
): StanyOdczytu {
  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner';
  wskaznik.setAttribute('role', 'status');
  wskaznik.setAttribute('aria-label', zdania.odczyt);

  const zdanie = document.createElement('span');
  zdanie.textContent = zdania.odczyt;

  const odczyt = document.createElement('p');
  odczyt.className = 'dk-stany__odczyt';
  odczyt.append(wskaznik, zdanie);

  const powod = document.createElement('span');

  const ponowienie = document.createElement('button');
  ponowienie.type = 'button';
  ponowienie.className = 'dn-btn dn-btn--zarys dn-btn--sm';
  ponowienie.textContent = 'Spróbuj ponownie';
  ponowienie.addEventListener('click', ponow);

  const blad = document.createElement('p');
  blad.className = 'dk-stany__blad';
  blad.append(powod, ponowienie);

  const element = document.createElement('div');
  element.className = 'dk-stany';
  element.append(odczyt, blad);

  function nanies(faza: FazaOdczytu): void {
    odczyt.hidden = faza !== 'odczyt';
    blad.hidden = faza !== 'blad';
    element.hidden = faza === 'spoczynek' || faza === 'gotowe';
    powod.textContent =
      faza === 'blad'
        ? `${zdania.blad} ${stan.powodNiepowodzenia()} `
        : '';
  }

  nanies(stan.faza());

  return {
    element,
    odswiez: () => nanies(stan.faza()),
  };
}
