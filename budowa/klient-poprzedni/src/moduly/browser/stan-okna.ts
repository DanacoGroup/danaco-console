/**
 * Stany okna modułu Browser: pusty, ładowania, błędu i gotowy. Nośnik trzyma stan
 * wraz z komunikatem, a widok dostaje osobny element treści i wypełnia go swoim.
 * Stan nie kasuje treści, tylko ją przesłania.
 */

import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';

/**
 * Zdanie okna o sobie samym — treść stanu pustego przed pierwszym odczytem. Opis jest
 * argumentem fabryki, ponieważ stan pusty ma nazwać, czym okno jest, i zachęcić do
 * pierwszej czynności, zamiast pokazywać komunikat ogólny.
 */
export interface OpisOkna {
  /** Tytuł stanu pustego. */
  tytul: string;
  /** Czym okno jest i jak je zapełnić — jedno zdanie, nie kod błędu. */
  opis: string;
}

export interface StanOkna {
  /** Element osadzany w oknie; niesie komunikat i treść. */
  element: HTMLElement;
  /** Miejsce na treść okna — wypełnia je widok. */
  tresc: HTMLElement;
  /** Zapowiada trwające wywołanie rdzenia. */
  ladowanie(opis: string): void;
  /** Pustka merytoryczna: wywołanie się udało, wyniku nie ma. */
  puste(tytul: string, opis: string): void;
  /** Odmowa rdzenia albo awaria wraz z jej powodem. */
  blad(opis: string): void;
  /** Zdejmuje komunikat i odsłania treść. */
  gotowe(): void;
  /** Faza bieżąca — do decyzji widoku i do sprawdzianu. */
  faza(): FazaOkna;
  /** Podpina „Spróbuj ponownie"; przycisk widać wyłącznie w stanie błędu. */
  ustawPonowienie(ponow: () => void): void;
  /** Wraca do stanu pustego opisem, z którym okno powstało. */
  pusteZOpisu(): void;
}

export function utworzStanOkna(opisOkna: OpisOkna): StanOkna {
  let biezaca: FazaOkna = 'puste';

  const spinner = document.createElement('span');
  spinner.className = 'dn-spinner mb-stan__spinner';
  spinner.setAttribute('role', 'status');
  spinner.setAttribute('aria-label', 'Odczyt z rdzenia w toku');

  const tytul = document.createElement('p');
  tytul.className = 'dn-pusty-stan-tytul mb-stan__tytul';

  const opis = document.createElement('p');
  opis.className = 'dn-pusty-stan-opis mb-stan__opis';

  const ponowienie = document.createElement('button');
  ponowienie.type = 'button';
  ponowienie.className = 'dn-btn dn-btn--sm dn-btn--zarys mb-stan__ponow';
  ponowienie.textContent = 'Spróbuj ponownie';
  ponowienie.hidden = true;

  const powloka = document.createElement('div');
  powloka.className = 'dn-pusty-stan mb-stan';
  powloka.append(spinner, tytul, opis, ponowienie);

  const tresc = document.createElement('div');
  tresc.className = 'mb-stan__tresc';

  const element = document.createElement('div');
  element.className = 'mb-stan__powloka';
  element.append(powloka, tresc);

  let ponow: (() => void) | null = null;

  function ustaw(faza: FazaOkna, naglowek: string, zdanie: string): void {
    biezaca = faza;
    oznaczFaze(element, powloka, faza);
    tytul.textContent = naglowek;
    tytul.hidden = naglowek === '';
    opis.textContent = zdanie;
    spinner.hidden = faza !== 'ladowanie';
    ponowienie.hidden = faza !== 'blad' || ponow === null;
    // Treść zostaje widoczna także w błędzie; znika wyłącznie na czas odczytu.
    tresc.hidden = faza === 'ladowanie';
  }

  ponowienie.addEventListener('click', () => ponow?.());

  function pusteZOpisu(): void {
    ustaw('puste', opisOkna.tytul, opisOkna.opis);
  }

  // Do pierwszego odświeżenia widoczne jest zdanie stanu pustego.
  pusteZOpisu();

  return {
    element,
    tresc,
    ladowanie: (zdanie) => ustaw('ladowanie', '', zdanie),
    puste: (naglowek, zdanie) => ustaw('puste', naglowek, zdanie),
    blad: (zdanie) => ustaw('blad', 'Rdzeń odmówił', zdanie),
    gotowe: () => ustaw('gotowe', '', ''),
    faza: () => biezaca,
    ustawPonowienie(nowe) {
      ponow = nowe;
      ponowienie.hidden = biezaca !== 'blad';
    },
    pusteZOpisu,
  };
}
