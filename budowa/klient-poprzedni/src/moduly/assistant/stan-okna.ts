import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';

/**
 * Powłoka komunikatu stanu okien modułu Assistant.
 *
 * Plik odpowiada wyłącznie za nośnik komunikatu stanu wraz z miejscem na treść.
 * Nazwę fazy i znakowanie DOM wnosi `komponenty/faza-okna`; tutaj zostaje
 * wyłącznie to, czym moduł różni się od reszty drzewa: wskaźnik odczytu stoi
 * obok treści, a treść zostaje widoczna także w ładowaniu, więc kontrolki
 * pozostają klikalne przez cały czas wywołania.
 *
 * „Jeszcze nie pytałem rdzenia" nie jest fazą okna, tylko fazą źródła danych
 * (`zapis-modulu.ts`, pola `pytanoO…`). Okno pokazuje ją jako `puste` z własnym
 * zdaniem.
 *
 * Stan nie zastępuje treści, tylko ją przesłania: po powrocie do fazy `gotowe`
 * wcześniejsza treść jest nietknięta, więc nieudane odświeżenie nie kasuje
 * tego, co było już na ekranie.
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
