import { elementIkony, type NazwaIkony } from '../ikony/ikony';

// Brak wyników to wspólna treść pustki dla wykazów z filtrem i dla wyszukiwarki.

/** Który z czterech przypadków pustki pokazuje w tej chwili nośnik, wybierany przez wołający wykaz albo wyszukiwarkę. */
export type RodzajPustki = 'pusto' | 'bez-trafien' | 'w-odczycie' | 'odmowa';

export interface BrakWynikow {
  /** Nośnik montowany w miejscu wykazu; wołający sam go pokazuje i chowa. */
  element: HTMLElement;
  /** Który przypadek stoi na ekranie; pusty napis znaczy „jeszcze żaden". */
  rodzaj(): RodzajPustki | '';
  /** Wykaz naprawdę nie ma pozycji — stan poprawny, nie usterka, z czynnością do pojawienia pozycji. */
  pusto(co: string, coZmieni: string): void;
  /** Filtr przyciął wykaz do zera; zdanie cytuje frazę i zakres przeszukania, nie mówi brak danych. */
  bezTrafien(fraza: string, coPrzeszukano: string): void;
  /**
   * Rdzeń nie odpowiedział; zdanie nie używa słowa „brak", bo cisza nie jest
   * orzeczeniem o pustce.
   */
  wOdczycie(co: string): void;
  /** Odmowa rdzenia — jego zdanie idzie dosłownie, bez parafrazy. */
  odmowa(zdanieRdzenia: string): void;
}

/** Ikona przypadku pustki; rozróżnienie między czterema stanami nigdy nie stoi na samej barwie tła ikony. */
const IKONY: Readonly<Record<RodzajPustki, NazwaIkony>> = {
  pusto: 'info',
  'bez-trafien': 'szukaj',
  'w-odczycie': 'zegar',
  odmowa: 'ostrzezenie',
};

/**
 * Zdanie zastępcze, gdy rdzeń odmówił bez treści: pusta odmowa zdarza się przy zerwaniu transportu, a dopisanie rdzeniowi zdania, którego nie powiedział, byłoby zmyśleniem.
 */
const ODMOWA_BEZ_TRESCI =
  'Rdzeń odmówił bez podania powodu — odpowiedź przyszła pusta. '
  + 'Najczęściej znaczy to zerwane połączenie, nie zakaz.';

export function utworzBrakWynikow(): BrakWynikow {
  let biezacy: RodzajPustki | '' = '';

  const element = document.createElement('div');
  element.className = 'dn-pusty-stan';
  // `status`, nie `alert`: pustka nie jest błędem i nie ma przerywać czytania.
  element.setAttribute('role', 'status');

  const tytul = document.createElement('p');
  tytul.className = 'dn-pusty-stan-tytul';

  const opis = document.createElement('p');
  opis.className = 'dn-pusty-stan-opis';

  element.append(tytul, opis);

  /** Jedno miejsce nadające postać — cztery przypadki różnią się treścią. */
  function ubierz(rodzaj: RodzajPustki, naglowek: string, zdanie: string): void {
    biezacy = rodzaj;
    element.dataset['pustka'] = rodzaj;
    tytul.textContent = naglowek;
    opis.textContent = zdanie;
    element.replaceChildren(elementIkony(IKONY[rodzaj], { rozmiar: 28 }), tytul, opis);
  }

  return {
    element,

    rodzaj: () => biezacy,

    pusto(co, coZmieni) {
      ubierz(
        'pusto',
        `Brak ${co}.`,
        `Wykaz jest pusty i to jest stan poprawny — nie ma tu jeszcze ani jednej pozycji. ${coZmieni}`,
      );
    },

    bezTrafien(fraza, coPrzeszukano) {
      ubierz(
        'bez-trafien',
        `Fraza „${fraza}" nie pasuje do niczego.`,
        `Przeszukano: ${coPrzeszukano}. Pozycje są — przyciął je filtr, `
        + 'więc wyczyszczenie pola pokaże cały wykaz z powrotem.',
      );
    },

    wOdczycie(co) {
      ubierz(
        'w-odczycie',
        `Odczyt ${co} w toku…`,
        'Rdzeń jeszcze nie odpowiedział, więc nie wiadomo, czy pozycje są, czy ich nie ma. '
        + 'Wykaz wypełni się sam, gdy odpowiedź przyjdzie.',
      );
    },

    odmowa(zdanieRdzenia) {
      const tresc = zdanieRdzenia.trim();
      ubierz('odmowa', 'Rdzeń odmówił wykazu.', tresc === '' ? ODMOWA_BEZ_TRESCI : tresc);
    },
  };
}
