import type { FazaOdczytu, StanKont } from './stan-kont';

/**
 * Stany odczytu sekcji modeli: pas ładowania oraz komunikat blokowy błędu,
 * osadzony pod paskiem osi i nad zakładkami. Zawartość zakładek zostaje na
 * miejscu i pozostaje czynna także przy niepowodzeniu odczytu rejestru kont.
 */
export interface StanyOdczytu {
  /** Pas stanów osadzany w sekcji modeli. */
  element: HTMLElement;
  /** Nanosi fazę odczytu ze stanu rejestru kont. */
  odswiez(): void;
}

export function utworzStanyOdczytu(stan: StanKont, ponow: () => void): StanyOdczytu {
  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner';
  wskaznik.setAttribute('role', 'status');
  wskaznik.setAttribute('aria-label', 'Czytam rejestr kont modeli i kont code CLI');

  const zdanie = document.createElement('span');
  zdanie.textContent = 'Czytam rejestr kont…';

  const odczyt = document.createElement('p');
  odczyt.className = 'dm-stany__odczyt';
  odczyt.append(wskaznik, zdanie);

  const powod = document.createElement('span');

  const ponowienie = document.createElement('button');
  ponowienie.type = 'button';
  ponowienie.className = 'dn-btn dn-btn--zarys dn-btn--sm';
  ponowienie.textContent = 'Spróbuj ponownie';
  ponowienie.addEventListener('click', ponow);

  const blad = document.createElement('p');
  blad.className = 'dm-stany__blad';
  blad.append(powod, ponowienie);

  const element = document.createElement('div');
  element.className = 'dm-stany';
  element.append(odczyt, blad);

  function nanies(faza: FazaOdczytu): void {
    odczyt.hidden = faza !== 'odczyt';
    blad.hidden = faza !== 'blad';
    element.hidden = faza === 'spoczynek' || faza === 'gotowe';
    powod.textContent =
      faza === 'blad' ? `Rejestr nie dotarł. ${stan.powodNiepowodzenia()} ` : '';
  }

  nanies(stan.faza());

  return {
    element,
    odswiez: () => nanies(stan.faza()),
  };
}
