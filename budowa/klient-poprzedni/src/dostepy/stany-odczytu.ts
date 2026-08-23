import type { FazaOdczytu, StanDostepow } from './stan-dostepow';

/**
 * Pas stanów odczytu sekcji dostępów: wskaźnik ładowania i komunikat błędu.
 * Stan pusty niesie sam wykaz, bo mówi o wykazie, nie o odczycie.
 *
 * Pusty wykaz punktów znaczy co innego, gdy rdzeń jeszcze nie odpowiedział, co
 * innego, gdy odmówił, i co innego, gdy rejestr jest pusty. Zlanie tych
 * przypadków w jeden kazałoby Operatorowi zgadywać, czy czekać, czy działać.
 *
 * „Spróbuj ponownie" wyzwala odczyt i zostawia komunikat na miejscu, dopóki
 * sytuacja nie ustanie; sekcja pozostaje przez cały czas czynna.
 *
 * Biblioteka `komponenty/` nie ma jeszcze reguł komunikatu blokowego, więc
 * komunikat stoi na klasie komunikatu tej sekcji.
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
