/**
 * Trzy stany obowiązkowe okna operacyjnego: puste, ładowanie i błąd. Zestaw faz
 * oraz znakowanie powłoki są wspólne dla wszystkich modułów, a tutaj zostaje to,
 * czym Agents się różni: własne klasy `da-stan*` i chowanie treści.
 */

import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';

/**
 * Powłoka okna wraz z komunikatem stanu: element osadzany w oknie, miejsce na
 * treść oraz czynności przestawiające fazę i odczytujące fazę bieżącą.
 */
export interface StanOkna {
  /** Element osadzany w oknie; niesie komunikat i treść. */
  element: HTMLElement;
  /** Miejsce na treść okna — wypełnia je widok. */
  tresc: HTMLElement;
  /** Zapowiada trwające wywołanie rdzenia. */
  ladowanie(opis: string): void;
  /** Pokazuje pustkę merytoryczną — wywołanie się udało, wyniku nie ma. */
  puste(opis: string): void;
  /** Pokazuje odmowę albo awarię wraz z jej powodem. */
  blad(opis: string): void;
  /** Zdejmuje komunikat i odsłania treść. */
  gotowe(): void;
  /** Faza bieżąca — do sprawdzianów i do decyzji widoku. */
  faza(): FazaOkna;
}

export function utworzStanOkna(): StanOkna {
  let biezaca: FazaOkna = 'puste';

  const komunikat = document.createElement('p');
  komunikat.className = 'dn-pusty-stan-opis da-stan__opis';

  const powloka = document.createElement('div');
  powloka.className = 'dn-pusty-stan da-stan';
  powloka.append(komunikat);

  const tresc = document.createElement('div');
  tresc.className = 'da-stan__tresc';

  const element = document.createElement('div');
  element.className = 'da-stan__powloka';
  element.append(powloka, tresc);

  function ustaw(faza: FazaOkna, opis: string): void {
    biezaca = faza;
    oznaczFaze(element, powloka, faza);
    komunikat.textContent = opis;
    tresc.hidden = faza === 'ladowanie';
  }

  ustaw('puste', 'Brak danych.');

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
