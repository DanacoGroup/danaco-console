import './obszary-sesji.css';

import { przyciskAkcji } from '../modele/kontrolki-formularza-braki';
import { opisAdresu } from './adres-ustawienia';
import { nazwaObszaru, OBSZARY_SESJI } from './obszary-sesji';
import type { StanObszarowSesji } from './stan-obszarow-sesji';
import { utworzStanyOdczytu } from './stany-odczytu';
import {
  trescWykazu,
  ubierzBraki,
  ubierzKatalog,
  ubierzWiersze,
  zlozWiersz,
  type WierszObszaru,
} from './wiersze-obszarow-sesji';

/**
 * Panel konfiguracji obowiązującej: okno komend `config.effective.get`
 * i `config.session.set`. Zasięg bierze z punktu widzenia ustawionego paskiem
 * okna, a ładowanie, odmowę i pustkę zaciąga z pasa `stany-odczytu`.
 */
export interface PanelObszarowSesji {
  /** Sekcja osadzana w ciele okna konfiguracji. */
  element: HTMLElement;
  /** Nanosi stan na zbudowany wykaz; wykazu nie przebudowuje. */
  odswiez(): void;
}

const ZDANIA_OBSZAROW = {
  odczyt: 'Rozstrzygam osiemnaście obszarów konfiguracji sesji…',
  blad: 'Konfiguracja obowiązująca nie dotarła; wykaz obszarów zostaje pusty.',
};

export function utworzPanelObszarowSesji(
  stan: StanObszarowSesji,
  ponow: () => void,
): PanelObszarowSesji {
  const adres = akapit('dk-obszary__opis');
  const katalog = akapit('dk-obszary__katalog');
  const odpowiedz = akapit('dk-obszary__odpowiedz');
  odpowiedz.hidden = true;

  const wiersze = OBSZARY_SESJI.map(zlozWiersz);
  const wykaz = document.createElement('ul');
  wykaz.className = 'dk-obszary__wykaz';
  wykaz.setAttribute('aria-label', 'Obszary konfiguracji sesji');

  const braki = document.createElement('ul');
  braki.className = 'dk-obszary__braki';

  const stany = utworzStanyOdczytu(stan, ponow, ZDANIA_OBSZAROW);
  const utrwal = przyciskAkcji('Utrwal wybrane obszary na tym poziomie', 'dn-btn dn-btn--zarys');
  utrwal.addEventListener('click', () => void wykonajZapis(stan, wiersze, odpowiedz));

  const element = document.createElement('section');
  element.className = 'dk-obszary';
  element.append(
    zlozNaglowek(adres),
    stany.element,
    wykaz,
    katalog,
    braki,
    zlozStopke(utrwal, odpowiedz),
  );

  function odswiez(): void {
    stany.odswiez();
    adres.textContent = zdanieAdresu(stan);
    wykaz.replaceChildren(...trescWykazu(stan, wiersze));
    ubierzWiersze(stan, wiersze);
    ubierzKatalog(stan, katalog);
    ubierzBraki(stan, braki);
  }

  odswiez();
  return { element, odswiez };
}

function akapit(klasa: string): HTMLElement {
  const element = document.createElement('p');
  element.className = klasa;
  return element;
}

/**
 * Składa nagłówek sekcji: tytuł, zdanie o ośmiu poziomach zasięgu i trzech
 * osiach rozstrzygania oraz akapit adresu przekazany przez wywołującego.
 */
function zlozNaglowek(adres: HTMLElement): HTMLElement {
  const tytul = document.createElement('h3');
  tytul.className = 'dk-obszary__tytul';
  tytul.textContent = 'Konfiguracja obowiązująca';

  const wyjasnienie = akapit('dn-pole-opis');
  wyjasnienie.textContent =
    'Obszary rozstrzygnięte przez rdzeń po ośmiu poziomach zasięgu i trzech osiach. ' +
    'Plakietka mówi, z którego rejestru i z którego poziomu przyszła wartość obszaru.';

  const element = document.createElement('header');
  element.className = 'dk-obszary__naglowek';
  element.append(tytul, wyjasnienie, adres);
  return element;
}

function zlozStopke(utrwal: HTMLElement, odpowiedz: HTMLElement): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dk-obszary__stopka';
  element.append(utrwal, odpowiedz);
  return element;
}

/**
 * Składa zdanie o adresie z opisu punktu widzenia stanu: ten sam adres służy
 * odczytowi konfiguracji obowiązującej i przyjmuje utrwalenie wybranych obszarów.
 */
function zdanieAdresu(stan: StanObszarowSesji): string {
  return (
    `Liczone dla: ${opisAdresu(stan.punkt())}. ` +
    'Utrwalenie zapisze wybrane obszary pod tym samym adresem.'
  );
}

/**
 * Zapisuje wybrane obszary komendą `config.session.set` i zawsze zwraca zdanie
 * odpowiedzi. Wybór pusty nie idzie do rdzenia, bo nie zapisałby niczego.
 */
async function wykonajZapis(
  stan: StanObszarowSesji,
  wiersze: readonly WierszObszaru[],
  odpowiedz: HTMLElement,
): Promise<void> {
  const wybrane = wiersze.filter((wiersz) => wiersz.wybor.checked).map((wiersz) => wiersz.obszar);
  if (wybrane.length === 0) {
    pokazOdpowiedz(odpowiedz, 'Nie wskazano ani jednego obszaru do utrwalenia.', false);
    return;
  }

  const wynik = await stan.utrwal(wybrane);
  if (!wynik.udany) {
    pokazOdpowiedz(odpowiedz, wynik.blad?.message ?? 'Rdzeń odmówił bez podania powodu.', false);
    return;
  }
  const zapisane = wynik.wynik?.storedAreas ?? [];
  pokazOdpowiedz(
    odpowiedz,
    `Obszary zapisane na tym poziomie: ${zapisane.map((obszar) => nazwaObszaru(obszar)).join(', ')}`,
    true,
  );
}

function pokazOdpowiedz(element: HTMLElement, tresc: string, powodzenie: boolean): void {
  element.textContent = tresc;
  element.hidden = false;
  element.dataset['powodzenie'] = String(powodzenie);
}
