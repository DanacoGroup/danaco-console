import './pomocnicze.css';

import type { Kanal } from '../protokol/kanal';
import { budzetRozmowy } from './budzet-rozmowy';
import type { PanelPomocniczy } from './panel-pomocniczy';
import { oknaPomocnicze, type OpisPomocniczego, type StanPomocniczego } from './rejestr-pomocniczych';
import { wytworniaPanelu } from './wytwornia-paneli';

/**
 * Pas okien pomocniczych modułu — Developer i Diagnostics.
 *
 * Moduły inżynierskie pracują na oknach obok rozmowy: terminal, przeglądarka,
 * artefakty, pliki, podgląd w tle bash, pliki środowiska. Pas jest miejscem,
 * w którym te okna stoją.
 *
 * Pozycja niezbudowana nie znika ze sceny: dostaje wiersz z nazwą,
 * przeznaczeniem i powodem braku. Pas z samymi oknami gotowymi wyglądałby
 * na kompletny, a brak przestałby być widoczny.
 *
 * Pas nie jest drugim złożeniem modułu: nie tworzy stanu modułu, nie zna komend
 * obszaru `developer.*` ani `diagnostics.*` i nie dotyka okien operacyjnych.
 * Składa gotowe okna pomocnicze oraz spis pozostałych pozycji.
 *
 * `zamknij()` jest obowiązkowe: okno podglądu trzyma subskrypcję `stream.chunk`
 * i bez tego wywołania zostaje ona żywa po zejściu pasa ze sceny.
 */
export interface PasPomocniczych {
  /** Element osadzany w złożeniu modułu. */
  element: HTMLElement;
  /** Odświeża okna zbudowane; spis pozycji się nie zmienia. */
  odswiez(): void;
  /**
   * Podaje pasowi okno wykonania nadane przez rdzeń.
   *
   * Nie każdy moduł zna swoje okno przy montażu: Developer i Diagnostics
   * montują się przez `widokZOknaSesji`, więc okno mają od razu, a Apps montuje
   * się z samym kanałem i poznaje okno dopiero z `window.list` w `wczytaj`.
   *
   * Panele powstają na nowo, bo okno biorą przy powołaniu — tą samą drogą, co
   * `okna-rownolegle/panele-gniazda.ts`: stare panele są zamykane wraz
   * z subskrypcjami rdzenia, nowe budowane z oknem właściwym. Podmiana okna
   * bez zamknięcia zostawiłaby subskrypcję pytającą o okno poprzednie.
   *
   * To samo okno podane drugi raz nie robi nic — przebudowa panelu, który już
   * pyta o właściwe okno, byłaby zerwaniem strumienia bez powodu.
   */
  ustawOkno(okno: string): void;
  /** Zamyka subskrypcje okien zbudowanych. */
  zamknij(): void;
}

/** Zależności pasa. */
export interface OpcjePasa {
  kanal: Kanal;
  /** Kod modułu — po nim idzie spis pozycji i profil rozmowy. */
  modul: string;
  /** Nazwa modułu w etykietach dostępności. */
  nazwaModulu: string;
  /** Okno wykonania modułu; puste znaczy „rdzeń nie dał temu modułowi okna". */
  okno: string;
  /** Przedrostek klas modułu: `mdev` albo `dg`. */
  przedrostek: string;
}

/** Zdanie wiersza spisu — po jednym na każdy stan pozycji. */
const ZAPOWIEDZ_STANU: Record<StanPomocniczego, string> = {
  'zbudowane': 'stoi w pasie',
  'stoi-w-module': 'niesie je złożenie modułu',
  'w-innej-pracy': 'buduje je osobna praca',
  'bez-komendy-rdzenia': 'brak komendy w kontrakcie',
  'do-rozstrzygniecia': 'czeka na rozstrzygnięcie Właściciela',
};

export function utworzPasPomocniczych(opcje: OpcjePasa): PasPomocniczych {
  const spis = oknaPomocnicze(opcje.modul);
  let okno = opcje.okno;
  let zbudowane: PanelPomocniczy[] = [];

  const element = document.createElement('section');
  element.className = 'dnp-pas';
  element.dataset['modul'] = opcje.modul;
  element.setAttribute('aria-label', `Okna pomocnicze modułu ${opcje.nazwaModulu}`);

  const tor = document.createElement('div');
  tor.className = 'dnp-pas__tor';

  element.append(naglowekPasa(opcje, spis.length), tor, spisPozostalych(spis, opcje.nazwaModulu));

  /**
   * Powołanie paneli pozycji zbudowanych dla okna bieżącego.
   *
   * Tor bez ani jednego panelu zostaje w drzewie, ale pusty — `:empty` nie
   * zabiera miejsca, a stały element pozwala podmienić zawartość bez ruszania
   * nagłówka i spisu braków.
   */
  function zloz(): void {
    zbudowane = [];
    const panele: HTMLElement[] = [];
    for (const pozycja of spis) {
      const panel = zbudujPozycje(pozycja, { ...opcje, okno });
      if (panel === null) continue;
      zbudowane.push(panel);
      panele.push(panel.element);
    }
    tor.replaceChildren(...panele);
  }

  zloz();

  return {
    element,

    odswiez() {
      for (const panel of zbudowane) panel.odswiez();
    },

    ustawOkno(nowe) {
      if (nowe === okno) return;
      // Zamknięcie przed przebudową: panel zdjęty z drzewa bez `zamknij()`
      // zostawia subskrypcję rdzenia pytającą o okno, którego już nie ma.
      for (const panel of zbudowane) panel.zamknij();
      okno = nowe;
      zloz();
    },

    zamknij() {
      for (const panel of zbudowane) panel.zamknij();
    },
  };
}

/**
 * Buduje panel pozycji zbudowanej; dla pozostałych stanów oddaje `null` —
 * ich miejsce jest w spisie poniżej toru, nie w torze.
 *
 * Mapowanie kodu pozycji na wytwórnię stoi poza pasem (`wytwornia-paneli.ts`),
 * więc pas nie zna ani jednego kodu, a dołożenie panelu nie wymaga jego edycji.
 * Tego samego mapowania używa kolumna paneli sceny okien równoległych.
 *
 * Pozycja nazwana w rejestrze zbudowaną, dla której nie ma wytwórni, jest
 * rozjazdem między spisem a kodem — stąd wpis w dzienniku zdarzeń zamiast
 * cichego `null`, który udawałby, że pozycji w spisie nie ma.
 */
function zbudujPozycje(pozycja: OpisPomocniczego, opcje: OpcjePasa): PanelPomocniczy | null {
  if (pozycja.stan !== 'zbudowane') return null;
  const wytwornia = wytworniaPanelu(pozycja.kod);
  if (wytwornia === null) {
    console.warn('[okna pomocnicze] pozycja oznaczona jako zbudowana nie ma wytwórni', pozycja.kod);
    return null;
  }
  return wytwornia({
    kanal: opcje.kanal,
    okno: opcje.okno,
    modul: opcje.nazwaModulu,
    przedrostek: opcje.przedrostek,
  });
}

/** Nagłówek pasa: budżet okien rozmowy i liczba pozycji spisu. */
function naglowekPasa(opcje: OpcjePasa, ile: number): HTMLElement {
  const naglowek = document.createElement('header');
  naglowek.className = 'dnp-pas__naglowek';

  const tytul = document.createElement('h3');
  tytul.className = 'dn-karta-tytul dnp-pas__tytul';
  tytul.textContent = 'Okna pomocnicze';

  const budzet = document.createElement('p');
  budzet.className = 'dn-pole-opis dnp-pas__budzet';
  budzet.textContent = budzetRozmowy(opcje.modul).zdanie;

  naglowek.append(tytul, budzet);
  if (ile === 0) naglowek.append(zdanieBezSpisu(opcje.nazwaModulu));
  return naglowek;
}

/**
 * Moduł spoza rejestru okien pomocniczych. Nie pustka — zdanie wprost, że
 * spisu dla tego modułu nie ma: pusty pas czytałby się jak „ten
 * moduł okien pomocniczych nie ma", a to co innego niż „nikt ich nie spisał".
 */
function zdanieBezSpisu(nazwaModulu: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis dnp-pas__uwaga';
  element.dataset['stan'] = 'bez-spisu';
  element.textContent =
    `Dla modułu ${nazwaModulu} nie spisano okien pomocniczych. Pas jest pusty dlatego, ` +
    'że spisu nie ma — nie dlatego, że moduł okien pomocniczych mieć nie miał.';
  return element;
}

/** Spis pozycji, które w torze nie stanęły — wraz z powodem każdej. */
function spisPozostalych(
  spis: readonly OpisPomocniczego[],
  nazwaModulu: string,
): HTMLElement {
  const pozostale = spis.filter((p) => p.stan !== 'zbudowane');
  const wykaz = document.createElement('dl');
  wykaz.className = 'dnp-braki';
  wykaz.setAttribute('aria-label', `Okna pomocnicze modułu ${nazwaModulu} jeszcze niestojące w pasie`);
  for (const pozycja of pozostale) wykaz.append(...wierszBraku(pozycja));
  wykaz.hidden = pozostale.length === 0;
  return wykaz;
}

/** Jeden wiersz spisu braków: nazwa ze stanem oraz zdanie wprost. */
function wierszBraku(pozycja: OpisPomocniczego): [HTMLElement, HTMLElement] {
  const nazwa = document.createElement('dt');
  nazwa.className = 'dnp-braki__nazwa';
  nazwa.dataset['okno'] = pozycja.kod;
  nazwa.dataset['stan'] = pozycja.stan;
  nazwa.textContent = `${pozycja.nazwa} — ${ZAPOWIEDZ_STANU[pozycja.stan]}`;

  const powod = document.createElement('dd');
  powod.className = 'dnp-braki__powod';
  powod.textContent = `${pozycja.przeznaczenie} ${pozycja.wyjasnienie}`;
  return [nazwa, powod];
}
