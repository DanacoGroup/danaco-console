/**
 * Wiersze wykazu obszarów i trzy pasy poboczne panelu konfiguracji
 * obowiązującej. Plik składa wiersz obszaru z pól kontraktu, buduje treść
 * wykazu wraz ze zdaniem o jego pustce oraz ubiera wiersze, katalog i braki.
 */
import type { SessionConfigArea } from '../../../shared/contract';
import { przelacznik } from '../modele/kontrolki-formularza-braki';
import {
  czyZapisWlasny,
  nazwaObszaru,
  opisKatalogu,
  opisPochodzenia,
} from './obszary-sesji';
import { stanPusty } from './panel-kategorii';
import type { StanObszarowSesji } from './stan-obszarow-sesji';

/**
 * Wiersz jednego obszaru: nazwa, plakietka pochodzenia oraz pole wyboru do zapisu.
 * Plakietka mówi, skąd wartość obszaru pochodzi, więc Operator widzi zapis własny
 * sesji obok wartości odziedziczonej, zanim cokolwiek zaznaczy.
 */
export interface WierszObszaru {
  element: HTMLElement;
  obszar: SessionConfigArea;
  znak: HTMLElement;
  wybor: HTMLInputElement;
}

export function zlozWiersz(obszar: SessionConfigArea): WierszObszaru {
  const wybor = przelacznik(`Obejmij obszar „${nazwaObszaru(obszar)}" zapisem`);
  wybor.classList.add('dk-obszary__wybor');

  const nazwa = document.createElement('span');
  nazwa.className = 'dk-obszary__nazwa';
  nazwa.textContent = nazwaObszaru(obszar);

  const znak = document.createElement('span');
  znak.className = 'dn-plakietka dk-obszary__znak';

  const element = document.createElement('li');
  element.className = 'dk-obszary__wiersz';
  element.dataset['obszar'] = obszar;
  element.append(wybor, nazwa, znak);
  return { element, obszar, znak, wybor };
}

/**
 * Treść wykazu: wiersze obszarów albo stan pusty. Przy braku konfiguracji
 * obowiązującej wchodzi stan pusty ze zdaniem właściwym dla fazy odczytu.
 */
export function trescWykazu(
  stan: StanObszarowSesji,
  wiersze: readonly WierszObszaru[],
): HTMLElement[] {
  if (stan.obowiazujaca() !== null) return wiersze.map((wiersz) => wiersz.element);
  return [
    stanPusty('Brak konfiguracji obowiązującej', zdanieBezObszarow(stan), 'dk-obszary__pusto'),
  ];
}

/**
 * Zdanie stanu pustego dobrane do fazy odczytu. Brak wykazu przed pierwszym pytaniem
 * i brak wykazu po odpowiedzi rdzenia to dwie różne rzeczy, więc każda faza dostaje
 * własne zdanie zamiast wspólnego komunikatu.
 */
function zdanieBezObszarow(stan: StanObszarowSesji): string {
  switch (stan.faza()) {
    case 'spoczynek':
      return 'Okno nie pytało jeszcze rdzenia o konfigurację obowiązującą.';
    case 'odczyt':
      return 'Rdzeń rozstrzyga obszary. Okno pozostaje otwarte i czynne.';
    case 'blad':
      return 'Rdzeń odmówił rozstrzygnięcia — powód nad wykazem. Odczyt można powtórzyć.';
    case 'gotowe':
      return 'Rdzeń odpowiedział, lecz nie oddał konfiguracji obowiązującej.';
  }
}

/**
 * Plakietka pochodzenia przy każdym obszarze. Odróżnia obszar zapisany na tym
 * poziomie od odziedziczonego z poziomu szerszego lub wziętego z innego rejestru.
 */
export function ubierzWiersze(
  stan: StanObszarowSesji,
  wiersze: readonly WierszObszaru[],
): void {
  const zapisane = new Set(stan.zapisane());
  for (const wiersz of wiersze) {
    const pochodzenie = stan.pochodzenie(wiersz.obszar);
    const wlasny = pochodzenie !== null && czyZapisWlasny(pochodzenie);
    const opis =
      pochodzenie === null ? 'rdzeń nie mówił o tym obszarze' : opisPochodzenia(pochodzenie);

    wiersz.znak.textContent = opis;
    wiersz.znak.className = wlasny
      ? 'dn-plakietka dn-plakietka--sygnal dk-obszary__znak'
      : 'dn-plakietka dk-obszary__znak';
    wiersz.znak.title = `Wartość obszaru pochodzi z: ${opis}`;
    wiersz.element.dataset['wlasny'] = String(wlasny);
    wiersz.element.dataset['utrwalony'] = String(zapisane.has(wiersz.obszar));
  }
}

/**
 * Rozejście katalogu roboczego wypisane wprost, a nie pomijane milczeniem. Katalog
 * inny niż obowiązujący zmienia znaczenie każdej ścieżki w sesji, więc pas poboczny
 * nazywa różnicę zamiast zostawiać ją do odkrycia.
 */
export function ubierzKatalog(stan: StanObszarowSesji, katalog: HTMLElement): void {
  const obowiazujaca = stan.obowiazujaca();
  if (obowiazujaca === null) {
    katalog.hidden = true;
    katalog.textContent = '';
    return;
  }
  const rozejscie = obowiazujaca.workingDirectory;
  katalog.hidden = false;
  katalog.textContent = `Katalog roboczy — ${opisKatalogu(rozejscie)}`;
  katalog.dataset['rozejscie'] = String(rozejscie.degraded);
}

/**
 * Pola, których adapter dostawcy nie obsłuży. Wypełnienie takiego pola niczego
 * nie wstrzymuje, więc wykaz jest ostrzeżeniem, nie blokadą.
 */
export function ubierzBraki(stan: StanObszarowSesji, braki: HTMLElement): void {
  const pola = stan.nieobsluzone();
  braki.hidden = pola.length === 0;
  braki.replaceChildren(
    ...pola.map((pole) => {
      const wiersz = document.createElement('li');
      wiersz.className = 'dk-obszary__brak';
      wiersz.textContent =
        `${pole.fieldPath} — ${pole.support}: ${pole.reason ?? 'adapter nie podał powodu'}`;
      return wiersz;
    }),
  );
}
