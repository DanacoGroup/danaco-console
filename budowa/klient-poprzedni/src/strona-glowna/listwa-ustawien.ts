import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { utworzNaglowekStrefy } from './naglowek-strefy';
import {
  KODY_STREFY_TRZECIEJ,
  POZYCJE_USTAWIEN,
  type PozycjaUstawienia,
} from './pozycje-ustawien';
import { utworzSygnalWyboru, type SluchaczWyboru } from './sygnal-wyboru';

/**
 * Strefa trzecia — listwa ustawień: zwarty pasek narzędziowy.
 *
 * Jedna odpowiedzialność: pas wejść pomocniczych rozpięty na całą szerokość.
 * Segmenty dzieli delikatny separator wewnątrz jednej powierzchni, więc pas
 * czyta się jako jeden byt o kilku wejściach, a nie jako trzecia siatka kart —
 * opracowanie mówi wprost, że pozycje listwy nie mają formy kart ani kafli
 * (rozdz. 3.4).
 *
 * Segmenty ustawień są trzy — tyle wymienia tabela zawartości listwy: Okno
 * konfiguracji, Mobile, Always On Display (wykaz w `pozycje-ustawien.ts`).
 *
 * Czwarty segment nie jest ustawieniem: „Dodaj nowy" otwiera formularz
 * zakładania komponentu własnego. Zgłasza to osobnym wywołaniem zwrotnym, nie
 * przez wykaz ustawień — pozycja, która ustawieniem nie jest, nie udaje jego
 * kodu. Wedle opracowania zakładanie komponentu należy do strefy drugiej;
 * przeniesienie segmentu wymaga zmiany montażu w `aplikacja/`, więc zostaje
 * zgłoszone, a nie wykonane z tego katalogu.
 *
 * Waga wizualna strefy jest najniższa z trzech: segment ma wysokość kontrolki
 * (32 px) i niesie ikonę oraz nazwę, bez wezwania do działania i bez metadanych.
 *
 * Żadna pozycja nie jest wyszarzona ani pozbawiona klikalności; wykaz skutków
 * mieszka w `aplikacja/akcje-ustawien.ts`.
 */

/**
 * Ikona segmentu ustawienia.
 *
 * Ciąg glifów maleje przez trzy strefy: godło środowiska 24, ikona kafla modułu
 * 18, ikona segmentu ustawienia 16 — strefa o najniższej wadze niesie znak
 * najlżejszy.
 */
const ROZMIAR_IKONY = 16;

/** Nazwa strefy z opracowania (rozdz. 3.1, 3.4) i z makiety — bez parafrazy. */
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

  // Trzy segmenty wykazu: pozostałe pozycje są zakresami Okna konfiguracji
  // i dublowałyby jego segment (powód w `pozycje-ustawien.ts`).
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

/**
 * Jeden segment belki: ikona i nazwa.
 *
 * Rozwinięcie nazwy idzie w dymek, nie na belkę: wyjaśnienia bywają długie
 * i postawione w pasie rozsadziłyby go do wysokości strefy drugiej, odwracając
 * hierarchię wag. W segmencie stoi nazwa, a rozwinięcie czyta technologia
 * wspomagająca i dymek.
 */
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
