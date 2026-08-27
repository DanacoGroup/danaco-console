import type { FazaOdczytu, StanDostepow } from './stan-dostepow';

/**
 * Pas stanów odczytu sekcji dostępów: wskaźnik ładowania oraz komunikat błędu
 * z ponowieniem. Stan pusty niesie sam wykaz, ponieważ mówi o zawartości
 * wykazu, a nie o przebiegu odczytu.
 */
export interface StanyOdczytu {
  /** Pas stanów osadzany pod nagłówkiem sekcji. */
  element: HTMLElement;
  /** Nanosi fazę odczytu ze stanu sekcji. */
  odswiez(): void;
}

export function utworzStanyOdczytu(stan: StanDostepow, ponow: () => void): StanyOdczytu {
  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner';
  wskaznik.setAttribute('role', 'status');
  wskaznik.setAttribute('aria-label', 'Czytam wykaz punktów dostępu i nadania okna');

  const zdanieOdczytu = document.createElement('span');
  zdanieOdczytu.textContent = 'Czytam wykaz punktów dostępu i nadania okna…';

  const odczyt = document.createElement('p');
  odczyt.className = 'dd-stany__odczyt';
  odczyt.append(wskaznik, zdanieOdczytu);

  const powod = document.createElement('span');

  const ponowienie = document.createElement('button');
  ponowienie.type = 'button';
  ponowienie.className = 'dn-btn dn-btn--zarys dn-btn--sm';
  ponowienie.textContent = 'Spróbuj ponownie';
  ponowienie.addEventListener('click', ponow);

  const blad = document.createElement('p');
  blad.className = 'dd-komunikat dd-stany__blad';
  blad.dataset.powodzenie = 'false';
  blad.append(powod, ponowienie);

  const element = document.createElement('div');
  element.className = 'dd-stany';
  element.append(odczyt, blad);

  function nanies(faza: FazaOdczytu): void {
    odczyt.hidden = faza !== 'odczyt';
    blad.hidden = faza !== 'blad';
    element.hidden = faza === 'spoczynek' || faza === 'gotowe';
    powod.textContent =
      faza === 'blad'
        ? `Wykaz nie dotarł w całości. ${stan.powodNiepowodzenia()} `
        : '';
  }

  nanies(stan.faza());

  return {
    element,
    odswiez: () => nanies(stan.faza()),
  };
}
