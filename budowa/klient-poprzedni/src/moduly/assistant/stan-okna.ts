import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';

/**
 * Powłoka komunikatu stanu okien modułu Assistant, czyli nośnik zdania
 * o fazie wraz z miejscem na treść okna. Wskaźnik odczytu stoi obok treści,
 * a treść pozostaje widoczna także w trakcie wywołania rdzenia, więc
 * kontrolki pozostają klikalne.
 */
export interface StanOkna {
  /** Element osadzany w oknie; niesie komunikat i treść. */
  element: HTMLElement;
  /** Miejsce na treść okna — wypełnia je widok. */
  tresc: HTMLElement;
  /** Zapowiada trwające wywołanie rdzenia; wskaźnik stoi obok treści. */
  ladowanie(opis: string): void;
  /** Pustka merytoryczna — także „jeszcze nie pytałem rdzenia". */
  puste(opis: string): void;
  /** Odmowa albo awaria wraz z powodem; komunikat zostaje w układzie. */
  blad(opis: string): void;
  /** Zdejmuje komunikat i odsłania samą treść. */
  gotowe(): void;
  /** Faza bieżąca — do sprawdzianów i do decyzji widoku. */
  faza(): FazaOkna;
}

export function utworzStanOkna(): StanOkna {
  let biezaca: FazaOkna = 'puste';

  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner';
  wskaznik.setAttribute('role', 'status');

  const zdanie = document.createElement('span');

  const pas = document.createElement('p');
  pas.className = 'dn-pusty-stan-opis ma-stan__opis';
  pas.append(wskaznik, zdanie);

  const powloka = document.createElement('div');
  powloka.className = 'dn-pusty-stan ma-stan';
  powloka.append(pas);

  const tresc = document.createElement('div');
  tresc.className = 'ma-stan__tresc';

  const element = document.createElement('div');
  element.className = 'ma-stan__powloka';
  element.append(powloka, tresc);

  function ustaw(faza: FazaOkna, opis: string): void {
    biezaca = faza;
    oznaczFaze(element, powloka, faza);
    zdanie.textContent = opis;
    wskaznik.hidden = faza !== 'ladowanie';
    wskaznik.setAttribute('aria-label', opis);
  }

  ustaw('puste', 'Jeszcze nie pytałem rdzenia.');

  return {
    element,
    tresc,
    ladowanie: (opis) => ustaw('ladowanie', opis),
    puste: (opis) => ustaw('puste', opis),
    blad: (opis) => ustaw('blad', opis),
    gotowe: () => ustaw('gotowe', ''),
    faza: () => biezaca,
  };
}
