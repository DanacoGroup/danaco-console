import type { FazaOdczytu } from './stan-konfiguracji';

/**
 * Pas stanów odczytu okna konfiguracji: wskaźnik ładowania oraz komunikat
 * blokowy niepowodzenia z przyciskiem ponowienia. Pas leży nad stopką okna,
 * obok formularza, nigdy zamiast niego, ponieważ pola pozostają edytowalne
 * przez cały czas odczytu.
 */
export interface StanyOdczytu {
  /** Pas stanów osadzany między ciałem okna a jego stopką. */
  element: HTMLElement;
  /** Nanosi fazę odczytu ze stanu okna. */
  odswiez(): void;
}

/**
 * Tyle stanu, ile pas naprawdę czyta: faza odczytu i powód niepowodzenia. Pas
 * nie zna ani katalogu ustawień, ani obszarów sesji, dzięki czemu obsługuje oba
 * odczyty okna konfiguracji, `config.get` oraz `config.effective.get`, bez
 * rozgałęzienia.
 */
export interface ZrodloFazyOdczytu {
  faza(): FazaOdczytu;
  powodNiepowodzenia(): string;
}

/**
 * Zdania pasa, podstawiane przy tworzeniu; domyślne mówią o katalogu ustawień.
 * Wywołanie dla konfiguracji obowiązującej podaje własną parę, ponieważ obraz
 * jest ten sam, a nazwa czytanego zbioru inna.
 */
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
