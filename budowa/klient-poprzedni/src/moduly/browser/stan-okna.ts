/**
 * Stany okna modułu Browser: pusty, ładowania, błędu i gotowy. Nośnik trzyma
 * stan wraz z komunikatem; widok dostaje `tresc` i wypełnia ją swoim.
 *
 * Stany są rozdzielone, bo „jeszcze nie pytałem", „pytam" i „rdzeń odmówił" to
 * trzy różne sytuacje; zlanie ich w jedno kazałoby Operatorowi zgadywać, czy
 * czekać, czy działać (`dostepy/stany-odczytu.ts`).
 *
 * Stan nie kasuje treści, tylko ją przesłania — nieudane odświeżenie zostawia
 * to, co Operator już widział, a „Spróbuj ponownie" nie zdejmuje komunikatu
 * błędu.
 *
 * Nazwa fazy i jej znakowanie pochodzą z `komponenty/faza-okna`. Wartość
 * trafia do `data-faza`, po którym sięgają arkusze i sprawdziany — własny
 * zestaw nazw w module znaczyłby, że ten sam stan okna nazywa się gdzie
 * indziej inaczej.
 *
 * Opis okna jest argumentem fabryki, bo stan pusty ma nazwać, czym okno jest
 * i zachęcić do pierwszej czynności, zamiast pokazywać komunikat ogólny.
 */

import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';

/** Zdanie okna o sobie samym — treść stanu pustego przed pierwszym odczytem. */
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
    // Treść zostaje widoczna także w błędzie: odmowa odświeżenia nie kasuje
    // tego, co Operator już przeczytał. Znika wyłącznie na czas odczytu.
    tresc.hidden = faza === 'ladowanie';
  }

  ponowienie.addEventListener('click', () => ponow?.());

  function pusteZOpisu(): void {
    ustaw('puste', opisOkna.tytul, opisOkna.opis);
  }

  // Między złożeniem okna a pierwszym odświeżeniem Operator widzi zdanie,
  // które i tak zobaczyłby przy pustym wyniku.
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
