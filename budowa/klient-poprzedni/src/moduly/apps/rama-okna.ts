/**
 * Okno operacyjne modułu Apps, czyli rama biblioteczna spleciona z przesłoną
 * stanu. Rama wnosi nagłówek, pas akcji i ciało, a przesłona znakuje fazę
 * odczytu wspólną dla pięciu okien modułu i przykrywa treść zamiast ją
 * kasować.
 */

import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';
import { utworzRameOkna, type RolaOkna } from '../../komponenty/rama-okna';

export interface RamaApps {
  /** Element osadzany w przestrzeni modułu. */
  element: HTMLElement;
  /** Miejsce na treść okna — wypełnia je plik okna. */
  tresc: HTMLElement;
  /** Panel akcji okna: wykaz braków kontraktu i narzędzia kontekstowe. */
  akcje: HTMLElement;
  ladowanie(opis: string): void;
  puste(opis: string): void;
  blad(opis: string): void;
  gotowe(): void;
  faza(): FazaOkna;
}

export function utworzRameApps(kodOkna: string, tytul: string, rola: RolaOkna): RamaApps {
  let biezaca: FazaOkna = 'puste';
  const rama = utworzRameOkna({ tytul, rola, kod: kodOkna, modul: 'Apps', przedrostek: 'mp' });

  const komunikat = document.createElement('p');
  komunikat.className = 'dn-pusty-stan-opis mp-stan__opis';

  const wskaznik = document.createElement('span');
  wskaznik.className = 'dn-spinner';
  wskaznik.setAttribute('role', 'status');
  wskaznik.setAttribute('aria-label', `Okno ${tytul} czeka na odpowiedź rdzenia`);

  const pas = document.createElement('div');
  pas.className = 'dn-pusty-stan mp-stan';
  pas.append(wskaznik, komunikat);

  const tresc = document.createElement('div');
  tresc.className = 'mp-okno__tresc';
  rama.cialo.append(pas, tresc);

  function ustaw(faza: FazaOkna, opis: string): void {
    biezaca = faza;
    oznaczFaze(rama.element, pas, faza);
    komunikat.textContent = opis;
    wskaznik.hidden = faza !== 'ladowanie';
  }

  ustaw('puste', 'Brak danych.');

  return {
    element: rama.element,
    tresc,
    akcje: rama.akcje,
    ladowanie: (opis) => ustaw('ladowanie', opis),
    puste: (opis) => ustaw('puste', opis),
    blad: (opis) => ustaw('blad', opis),
    gotowe: () => ustaw('gotowe', ''),
    faza: () => biezaca,
  };
}
