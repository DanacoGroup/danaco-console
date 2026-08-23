/**
 * Trzy stany obowiązkowe okna operacyjnego: puste, ładowanie, błąd.
 *
 * Rdzeń odpowiada na każde wywołanie — także odmową — a odmowa jest pokazywana
 * w oknie, w którym Operator ją wywołał, nie wyłącznie w konsoli przeglądarki.
 *
 * Stan nie zastępuje treści, tylko ją przesłania: gdy okno wraca do stanu
 * `gotowe`, wcześniejsza treść jest nietknięta. Dzięki temu nieudane odświeżenie
 * nie kasuje tego, co Operator już widział.
 *
 * Zestaw faz i znakowanie powłoki są wspólne dla wszystkich modułów
 * (`komponenty/faza-okna`) — tu zostaje to, czym Agents się różni: własne klasy
 * `da-stan*` i chowanie treści na czas ładowania.
 */

import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';

/** Powłoka okna wraz z komunikatem stanu. */
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
