/**
 * Trzy stany obowiązkowe okna modułu Translate — pusty, ładowania, błędu — dzieli jedna forma,
 * różnicowana wyłącznie treścią, z fazą trzymaną w atrybucie danych.
 */

import { elementIkony } from '../../ikony/ikony';
import { oznaczFaze, type FazaOkna } from '../../komponenty/faza-okna';

/** Zdanie stanu pustego niesie tytuł nad opisem, nazywając sytuację, w której okno stoi, oraz sposób jego zapełnienia treścią. */
export interface ZdanieStanu {
  /** Tytuł — nazywa sytuację, w której okno stoi. */
  readonly tytul: string;
  /** Opis — mówi, czym okno jest i jak je zapełnić. */
  readonly opis: string;
}

/** Powłoka okna wraz z pasem stanu udostępnia element osadzany w oknie, miejsce na treść oraz metody zmiany fazy okna. */
export interface StanOkna {
  /** Element osadzany w oknie; niesie pas stanu i treść. */
  element: HTMLElement;
  /** Miejsce na treść okna — wypełnia je widok. */
  tresc: HTMLElement;
  /** Zapowiada trwające wywołanie rdzenia. */
  ladowanie(opis: string): void;
  /** Pokazuje pustkę merytoryczną — wywołanie się udało, wyniku nie ma. */
  puste(zdanie: ZdanieStanu): void;
  /** Pokazuje odmowę albo awarię wraz z jej powodem. */
  blad(opis: string): void;
  /** Zdejmuje komunikat i odsłania treść. */
  gotowe(): void;
  /** Faza bieżąca — do sprawdzianów i do decyzji widoku. */
  faza(): FazaOkna;
}

/** Tytuł pasa w fazach, które nie są pustką, ma jedno wspólne brzmienie na trzy okna zarządców modułu tłumaczeń. */
const TYTUL_LADOWANIA = 'Rdzeń pracuje';
const TYTUL_BLEDU = 'Rdzeń odmówił';

export function utworzStanOkna(poczatkowe: ZdanieStanu): StanOkna {
  let biezaca: FazaOkna = 'puste';

  const tytul = document.createElement('p');
  tytul.className = 'dn-pusty-stan-tytul mt-stan__tytul';

  const komunikat = document.createElement('p');
  komunikat.className = 'dn-pusty-stan-opis mt-stan__opis';

  const spinner = document.createElement('span');
  spinner.className = 'dn-spinner mt-stan__spinner';
  spinner.setAttribute('role', 'status');
  spinner.setAttribute('aria-label', 'Trwa wywołanie rdzenia');
  spinner.hidden = true;
  tytul.append(spinner);

  const ikona = elementIkony('tlumacz', { rozmiar: 24, klasa: 'dn-ikona mt-stan__ikona' });

  const pas = document.createElement('div');
  pas.className = 'dn-pusty-stan mt-stan';
  pas.append(ikona, tytul, komunikat);

  const tresc = document.createElement('div');
  tresc.className = 'mt-stan__tresc';

  const element = document.createElement('div');
  element.className = 'mt-stan__powloka';
  element.append(pas, tresc);

  function ustaw(faza: FazaOkna, naglowek: string, opis: string): void {
    biezaca = faza;
    oznaczFaze(element, pas, faza);
    tytul.replaceChildren(naglowek, spinner);
    komunikat.textContent = opis;
    spinner.hidden = faza !== 'ladowanie';
  }

  ustaw('puste', poczatkowe.tytul, poczatkowe.opis);

  return {
    element,
    tresc,
    ladowanie: (opis) => ustaw('ladowanie', TYTUL_LADOWANIA, opis),
    puste: (zdanie) => ustaw('puste', zdanie.tytul, zdanie.opis),
    blad: (opis) => ustaw('blad', TYTUL_BLEDU, opis),
    gotowe: () => ustaw('gotowe', '', ''),
    faza: () => biezaca,
  };
}
