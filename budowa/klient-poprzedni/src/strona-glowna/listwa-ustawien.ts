import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { utworzNaglowekStrefy } from './naglowek-strefy';
import {
  KODY_STREFY_TRZECIEJ,
  POZYCJE_USTAWIEN,
  type PozycjaUstawienia,
} from './pozycje-ustawien';
import { utworzSygnalWyboru, type SluchaczWyboru } from './sygnal-wyboru';

/** Strefa trzecia jest listwą ustawień: zwartym paskiem narzędziowym o jednej odpowiedzialności. */

/**
 * Ikona segmentu ustawienia.
 *
 * Ciąg glifów maleje przez trzy strefy: godło środowiska 24, ikona kafla modułu
 * 18, ikona segmentu ustawienia 16 — strefa o najniższej wadze niesie znak
 * najlżejszy.
 */
const ROZMIAR_IKONY = 16;

/** Nazwa strefy trzeciej i jej wyjaśnienie pochodzą z ustalonej makiety, bez własnej parafrazy nazwy strefy. */
const ETYKIETA = 'Strefa 3 · Listwa ustawień';
const WYJASNIENIE =
  'Okno konfiguracji platformy oraz dwie funkcje globalne — Mobile i Always On Display.';

export interface ListwaUstawien {
  element: HTMLElement;
  naWybor(sluchacz: SluchaczWyboru<PozycjaUstawienia>): void;
  /** Miejsce na formularz zakładania — wpina je warstwa, która ma kanał. */
  przybornik: HTMLElement;
  /** Piąty segment belki — „Dodaj nowy". */
  naDodanie(sluchacz: () => void): void;
}

export function utworzListweUstawien(): ListwaUstawien {
  const sygnal = utworzSygnalWyboru<PozycjaUstawienia>();

  const element = document.createElement('section');
  element.className = 'dn-strona__strefa dn-strona__strefa--ustawienia';
  element.setAttribute('aria-label', ETYKIETA);

  const belka = document.createElement('div');
  belka.className = 'dn-strona__belka';

  // Trzy segmenty wykazu pomijają pozostałe pozycje, będące zakresami okna konfiguracji.
  for (const pozycja of POZYCJE_USTAWIEN) {
    if (!KODY_STREFY_TRZECIEJ.includes(pozycja.kod)) continue;
    belka.append(utworzSegment(pozycja.nazwa, pozycja.ikona, pozycja.wyjasnienie, () => sygnal.nadaj(pozycja), pozycja.kod));
  }

  const sluchacze = new Set<() => void>();
  belka.append(
    utworzSegment(
      'Dodaj nowy',
      'plus',
      'Zakłada komponent własny — automatykę, eksperta albo projekt.',
      () => {
        for (const sluchacz of sluchacze) sluchacz();
      },
      'dodaj',
    ),
  );

  const przybornik = document.createElement('div');
  przybornik.className = 'dn-strona__przybornik';

  element.append(utworzNaglowekStrefy(ETYKIETA, WYJASNIENIE), belka, przybornik);

  return {
    element,
    przybornik,
    naWybor: sygnal.sluchaj,
    naDodanie: (sluchacz) => {
      sluchacze.add(sluchacz);
    },
  };
}

/** Jeden segment belki niesie ikonę i nazwę; rozwinięcie nazwy czyta technologia wspomagająca i dymek, nie sama belka. */
function utworzSegment(
  nazwaPozycji: string,
  ikonaPozycji: NazwaIkony,
  wyjasnienie: string,
  przyWyborze: () => void,
  kod: string,
): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-strona__segment';
  element.dataset['ustawienie'] = kod;
  element.title = wyjasnienie;
  element.setAttribute('aria-description', wyjasnienie);

  const ikona = document.createElement('span');
  ikona.className = 'dn-strona__segment-ikona';
  ikona.append(elementIkony(ikonaPozycji, { rozmiar: ROZMIAR_IKONY }));

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-strona__segment-nazwa';
  nazwa.textContent = nazwaPozycji;

  element.append(ikona, nazwa);
  element.addEventListener('click', przyWyborze);
  return element;
}
