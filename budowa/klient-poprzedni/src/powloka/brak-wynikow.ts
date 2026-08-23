import { elementIkony, type NazwaIkony } from '../ikony/ikony';

/**
 * Brak wyników — wspólna treść pustki dla wykazów z filtrem i dla wyszukiwarki.
 *
 * Wygląd niesie klasa `.dn-pusty-stan` biblioteki, więc plik nie wnosi ani
 * jednej reguły arkusza; dokłada wyłącznie treść, która w ręcznie składanych
 * pustkach rozjeżdża się między wykazami.
 *
 * Cztery przypadki i cztery różne zdania:
 *   1. `pusto`      — wykaz naprawdę nie ma pozycji,
 *   2. `bezTrafien` — pozycje są, przyciął je filtr,
 *   3. `wOdczycie`  — rdzeń jeszcze nie odpowiedział,
 *   4. `odmowa`     — rdzeń odpowiedział odmową.
 * „Brak sesji" powiedziane w trzecim przypadku jest nieprawdą o stanie rdzenia,
 * bo cisza nie jest orzeczeniem; w drugim jest nieprawdą o danych, bo dane są,
 * a zasłonił je filtr, który da się zdjąć.
 *
 * Trzy pierwsze zdania mają jeden kształt: co · dlaczego · czym to zmienić.
 * Czwarte go nie ma, bo treść odmowy należy do rdzenia i idzie dosłownie —
 * parafraza odmowy jest atrapą odmowy.
 *
 * Nośnik nie zna wykazu, nad którym stoi, nie pyta rdzenia i nie rozstrzyga,
 * który przypadek zachodzi. Wie to wołający.
 */

/** Który z czterech przypadków pustki pokazuje w tej chwili nośnik. */
export type RodzajPustki = 'pusto' | 'bez-trafien' | 'w-odczycie' | 'odmowa';

export interface BrakWynikow {
  /** Nośnik montowany w miejscu wykazu; wołający sam go pokazuje i chowa. */
  element: HTMLElement;
  /** Który przypadek stoi na ekranie; pusty napis znaczy „jeszcze żaden". */
  rodzaj(): RodzajPustki | '';
  /**
   * Wykaz naprawdę nie ma pozycji — i to jest stan poprawny, nie usterka.
   *
   * @param co czego nie ma, w dopełniaczu („otwartych sesji").
   * @param coZmieni jedna czynność, po której pozycje się pojawią.
   */
  pusto(co: string, coZmieni: string): void;
  /**
   * Filtr przyciął wykaz do zera; zdanie nie mówi „brak danych", bo dane są.
   *
   * @param fraza treść pola, cytowana dosłownie, żeby było widać, czym wykaz
   *   został przycięty.
   * @param coPrzeszukano zakres przeszukania, wymieniony wprost.
   */
  bezTrafien(fraza: string, coPrzeszukano: string): void;
  /**
   * Rdzeń nie odpowiedział; zdanie nie używa słowa „brak", bo cisza nie jest
   * orzeczeniem o pustce.
   */
  wOdczycie(co: string): void;
  /** Odmowa rdzenia — jego zdanie idzie dosłownie, bez parafrazy. */
  odmowa(zdanieRdzenia: string): void;
}

/** Ikona przypadku; rozróżnienie nigdy nie stoi na samej barwie. */
const IKONY: Readonly<Record<RodzajPustki, NazwaIkony>> = {
  pusto: 'info',
  'bez-trafien': 'szukaj',
  'w-odczycie': 'zegar',
  odmowa: 'ostrzezenie',
};

/**
 * Zdanie zastępcze, gdy rdzeń odmówił bez treści.
 *
 * Pusta odmowa zdarza się przy zerwaniu transportu. Pusty akapit zostawiłby
 * tytuł bez powodu, a dopisanie rdzeniowi zdania, którego nie powiedział, byłoby
 * zmyśleniem; zostaje zdanie mówiące wprost, że powodu nie było.
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
