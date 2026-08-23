/**
 * Okno operacyjne modułu Apps: rama biblioteczna wraz z przesłoną stanu.
 *
 * Obudowa okna (nagłówek, plakietka roli, pas akcji, ciało) należy do
 * `komponenty/rama-okna.ts`, a znakowanie fazy do `komponenty/faza-okna.ts`.
 * Ten plik splata oba byty: rama nie zna cyklu życia danych, których nie
 * pobiera, a pięć okien modułu ma jedną przesłonę stanu, nie pięć wariantów
 * tego samego.
 *
 * Stan nie kasuje treści, tylko ją przesłania: powrót do fazy `gotowe` odsłania
 * to, co Operator już widział. Nieudane odświeżenie nie zabiera więc wyniku
 * poprzedniego.
 *
 * Wskaźnik odczytu stoi obok komunikatu, nigdy zamiast kontrolki: treść okna
 * zostaje na miejscu i pozostaje klikalna. Wspólna `oznaczFaze` tego nie
 * przesądza — moduły różnią się tu między sobą.
 *
 * Wygląd w całości z biblioteki `komponenty/` (`dn-*`) i żetonów `motyw/`;
 * plik nie zna ani jednej barwy i ani jednego odstępu.
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
