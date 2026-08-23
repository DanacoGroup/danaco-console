import type { ToolCatalogEntry } from '../../../shared/contract';
import { utworzMenuDrzewo, type LiscMenu } from '../komponenty/menu-drzewo';
import { utworzRejestrOstatnich } from './ostatnio-uzyte';
import {
  czyPowolanieNarzedzia,
  type DolozenieNarzedzia,
  type ZrodloWykazuUkosnika,
} from './zrodlo-wykazu-ukosnika';

/**
 * Wykaz po ukośniku — obsada mechanizmu z biblioteki, nie nowy byt.
 *
 * Plik nie rysuje ani jednego wiersza, nie filtruje, nie szereguje i nie
 * wytłuszcza. Składa dane w płaską tablicę `PozycjaMenu` i podaje ją
 * mechanizmowi `komponenty/menu-drzewo.ts`, który to wszystko już niesie.
 * Wykaz po ukośniku jest jego kolejną obsadą, obok sterów paska zlecenia —
 * z jedną różnicą: pracuje w trybie `bezUchwytu`, bo wyzwalaczem jest znak
 * w polu wpisywania, a nie kliknięcie w przycisk.
 *
 * Zachowania wiążące i miejsce, w którym każde stoi:
 *   - wyzwalacz „/" — `pole-wysylki.ts` woła `otworz`, ten `rozwin()`;
 *   - filtr w tym samym polu — `bezUchwytu: true` gasi pole mechanizmu, frazę
 *     podaje pole wypowiedzi przez `ustawFraze`;
 *   - otwarcie w górę — `kierunek: 'gora'`;
 *   - przewijanie z paskiem — mechanizm (`menu-drzewo.css`);
 *   - pierwsza pozycja wyróżniona — mechanizm wyróżnia ją po każdym
 *     odrysowaniu, Enter idzie przez `wybierzWyrozniona()`;
 *   - filtr od pierwszego znaku — `ustawFraze` przy każdym znaku;
 *   - dopasowanie w środku napisu — mechanizm (`ocenaPola` przez `indexOf`);
 *   - wytłuszczenie trafień — mechanizm (`wstawTrafienia`);
 *   - dwie nazwy — `nazwa` i `nazwaKrotka` pozycji;
 *   - porządek według trafności — mechanizm; ostatnio użyte dokłada ten plik.
 *
 * Do tego pliku należy więc jedno: „pozycje ostatnio użyte bliżej wierzchu".
 * Rejestru sięgania po pozycje nie ma ani w kontrakcie, ani w bibliotece, więc
 * prowadzi go `ostatnio-uzyte.ts`.
 *
 * Dwie oceny trafności się nie biją. Ten plik nie ocenia niczego: podaje wykaz
 * ułożony, ostatnio używanymi do przodu, a resztę rozstrzyga mechanizm. Jego
 * sortowanie jest stabilne i przy remisie oceny oddaje kolejność wejścia, czyli
 * tę ustawioną tutaj. Świeżość podnosi więc pozycję wśród równie trafnych
 * i ani o krok wyżej.
 */

/** Napis na uchwycie; w trybie `bezUchwytu` nie trafia na ekran, tylko do ARIA. */
const NASTAWA = 'Komendy po ukośniku';

export interface OpcjeWykazuUkosnika {
  /** Skąd biorą się pozycje (`tools.catalog.list`). */
  zrodlo: ZrodloWykazuUkosnika;
  /** Czym wykonać wybór narzędzia (`session.tool.attach`); brak = bez dokładania. */
  dolozenie?: DolozenieNarzedzia;
  /** Gdzie postawić zdanie dla Operatora — wpis automatyzacji w wątku. */
  naKomunikat?: (zdanie: string) => void;
  /** Co zrobić z pozycją rodzaju `action`; brak = wykaz powie, że wykonawcy nie ma. */
  naAkcje?: (pozycja: ToolCatalogEntry) => void;
}

export interface WykazPoUkosniku {
  /** Element montowany nad polem wypowiedzi. */
  element: HTMLElement;
  /** Czy wykaz jest w tej chwili rozwinięty. */
  otwarty(): boolean;
  /** Otwiera wykaz i sięga po pozycje do rdzenia. */
  otworz(): void;
  /** Podaje frazę z pola wypowiedzi — to jest całe filtrowanie. */
  ustawFraze(fraza: string): void;
  /** Przesuwa wyróżnienie o krok, nie ruszając ogniska. */
  przesun(krok: number): void;
  /** Zatwierdza wyróżnioną pozycję; fałsz znaczy „nie było czego zatwierdzić". */
  zatwierdz(): boolean;
  /** Zwija wykaz. */
  zamknij(): void;
  /**
   * Podpina słuchacza zmiany stanu wykazu.
   *
   * Pole wypowiedzi nie ma jak samo zauważyć, że wykaz zwinął się bez jego
   * udziału — mechanizm zamyka się także po kliknięciu poza nim i po wyborze
   * pozycji myszą. Bez tego sygnału podpowiedź „wpisz, by zawęzić" zostawałaby
   * pod polem, przy którym nic już nie stoi.
   */
  naZmiane(sluchacz: (otwarty: boolean) => void): void;
}

export function utworzWykazPoUkosniku(opcje: OpcjeWykazuUkosnika): WykazPoUkosniku {
  const rejestr = utworzRejestrOstatnich();
  const sluchacze: ((otwarty: boolean) => void)[] = [];

  /** Komplet pozycji z rdzenia — czytany na otwarcie, przesiewany przez mechanizm. */
  let wszystkie: readonly ToolCatalogEntry[] = [];
  /** Powód pustego wykazu; pusty napis znaczy „pozycje są". */
  let powod = '';
  let czynny = false;

  const menu = utworzMenuDrzewo({
    nastawa: NASTAWA,
    // Trzy nastawienia trybu płaskiego: wyzwalacz z klawiatury (bez uchwytu),
    // wykaz nad polem (w górę) i jeden opis zamiast ściany opisów przy setkach
    // pozycji.
    bezUchwytu: true,
    kierunek: 'gora',
    opisTylkoPrzyWyroznionej: true,
    naWybor: (klucz) => wykonaj(klucz),
  });
  const element = document.createElement('div');
  element.className = 'dc-ukosnik';
  element.append(menu.element);

  function oglosStan(otwarty: boolean): void {
    for (const sluchacz of sluchacze) sluchacz(otwarty);
  }

  /**
   * Wykaz ułożony świeżością — jedyne, co ten plik wnosi do porządku.
   *
   * Ostatnio użyte idą na przód w kolejności użycia, reszta zachowuje kolejność
   * z rdzenia (kontrakt: alfabetycznie w obrębie grupy). Mechanizm przy czynnej
   * frazie przestawi to według trafności, a przy remisie zostawi tak, jak tu.
   */
  function ulozone(): ToolCatalogEntry[] {
    const swieze = rejestr.klucze();
    if (swieze.length === 0) return [...wszystkie];
    const naPrzod: ToolCatalogEntry[] = [];
    for (const klucz of swieze) {
      const pozycja = wszystkie.find((p) => p.name === klucz);
      if (pozycja !== undefined) naPrzod.push(pozycja);
    }
    return [...naPrzod, ...wszystkie.filter((p) => !swieze.includes(p.name))];
  }

  function odrysuj(): void {
    menu.ustaw(NASTAWA, ulozone().map(liscWykazu));
  }

  /**
   * Jedna pozycja wykazu jako liść mechanizmu.
   *
   * `wybrany` mówi „ta pozycja jest już dołożona do sesji", a nie „ta jest
   * wyróżniona": wyróżnieniem zarządza mechanizm sam. Haczyk przy pozycji
   * dołożonej jest znacznikiem bieżącego wyboru.
   */
  function liscWykazu(pozycja: ToolCatalogEntry): LiscMenu {
    return {
      rodzaj: 'wybor',
      klucz: pozycja.name,
      nazwa: pozycja.name,
      // Druga nazwa wchodzi tylko wtedy, gdy naprawdę jest druga: skrót równy
      // nazwie pełnej narysowałby „schedule (schedule)".
      ...(pozycja.shortName === '' || pozycja.shortName === pozycja.name
        ? {}
        : { nazwaKrotka: pozycja.shortName }),
      opis: opisWiersza(pozycja),
      wybrany: pozycja.attached === true,
    };
  }

  /**
   * Opis pełnym zdaniem — po co to jest i kiedy użyć.
   *
   * Przed opisem z rdzenia stoi zdanie o skutku wyboru, bo skutek jest różny
   * dla dwóch rodzajów wpisu w jednym wykazie i z samej nazwy go nie widać:
   * jedno powiększa zestaw narzędzi modelu, drugie wykonuje czynność aplikacji.
   */
  function opisWiersza(pozycja: ToolCatalogEntry): string {
    const czlony: string[] = [];
    if (czyPowolanieNarzedzia(pozycja)) {
      czlony.push(
        pozycja.attached === true
          ? 'Narzędzie już dołożone do tej sesji.'
          : 'Powołanie narzędzia — poszerza zestaw modelu na czas sesji.',
      );
    } else {
      czlony.push('Komenda akcji — wykonuje czynność aplikacji, zestawu narzędzi nie tyka.');
    }
    if (pozycja.description !== '') czlony.push(pozycja.description);
    return czlony.join(' ');
  }

  /** Zdanie do wątku; bez odbiorcy nie ma dokąd pójść. */
  function zglos(zdanie: string): void {
    if (zdanie === '') return;
    opcje.naKomunikat?.(zdanie);
  }

  /** Wybór pozycji — tu rozchodzą się dwa rodzaje wpisu. */
  function wykonaj(klucz: string): void {
    const pozycja = wszystkie.find((p) => p.name === klucz);
    zamknij();
    if (pozycja === undefined) return;
    rejestr.zapamietaj(klucz);

    if (!czyPowolanieNarzedzia(pozycja)) {
      if (opcje.naAkcje !== undefined) {
        opcje.naAkcje(pozycja);
        return;
      }
      // Wykonania nie udajemy: komenda akcji ma w kontrakcie wskazaną komendę
      // do wywołania, ale jej kształtu żądania ten plik nie zna, a wysłanie
      // pustej treści byłoby zgadywaniem. Operator dostaje zdanie o tym, co się
      // nie stało, zamiast ciszy nieodróżnialnej od powodzenia.
      zglos(
        `Komenda akcji „${pozycja.name}" nie ma jeszcze wykonawcy po stronie klienta` +
          `${pozycja.command === undefined ? '' : ` (kontrakt wskazuje „${pozycja.command}")`}. ` +
          'Zestaw narzędzi modelu pozostał nietknięty.',
      );
      return;
    }

    if (opcje.dolozenie === undefined) {
      zglos(
        `Narzędzia „${pozycja.name}" nie dołożono: to okno nie dostało drogi ` +
          'do komendy session.tool.attach.',
      );
      return;
    }
    if (!pozycja.attachable) {
      zglos(
        `Pozycja „${pozycja.name}" jest w wykazie, ale rdzeń oznaczył ją jako ` +
          'niedołączalną — zestaw narzędzi sesji pozostał bez zmian.',
      );
      return;
    }
    opcje.dolozenie(pozycja, (zdanie) => zglos(zdanie));
  }

  function zamknij(): void {
    if (!czynny) return;
    czynny = false;
    menu.zwin();
    oglosStan(false);
  }

  /** Czy mechanizm ma wykaz rozwinięty — jego własny, wystawiony stan. */
  function mechanizmRozwiniety(): boolean {
    return menu.element.dataset['otwarte'] === 'tak';
  }

  /**
   * Godzi stan obsady ze stanem mechanizmu.
   *
   * Mechanizm zwija się także bez udziału obsady — po kliknięciu poza wykazem
   * zamyka się sam i nikogo o tym nie zawiadamia. Bez pogodzenia obsada
   * sterowałaby dalej
   * wykazem, którego na ekranie już nie ma: strzałki przesuwałyby niewidoczne
   * wyróżnienie, a podpowiedź „wpisz, by zawęzić" stałaby pod polem bez wykazu.
   */
  function pogodzStan(): boolean {
    if (czynny && !mechanizmRozwiniety()) zamknij();
    return czynny;
  }

  return {
    element,

    otwarty: () => czynny && mechanizmRozwiniety(),

    otworz() {
      czynny = true;
      // Wykaz rozwija się najpierw, z tym, co już jest: otwarcie ma być
      // natychmiastowe, a odpowiedź rdzenia przychodzi, kiedy przyjdzie.
      // Fraza wraca do pustej, bo ukośnik właśnie stanął sam.
      menu.ustawFraze('');
      odrysuj();
      menu.rozwin();
      oglosStan(true);

      opcje.zrodlo((pozycje, powodPustki) => {
        wszystkie = pozycje;
        powod = powodPustki;
        if (pozycje.length === 0 && powod !== '') zglos(powod);
        if (czynny) odrysuj();
      });
    },

    ustawFraze(nowa) {
      if (!pogodzStan()) return;
      menu.ustawFraze(nowa);
    },

    przesun(krok) {
      if (!pogodzStan()) return;
      menu.przesunWyroznienie(krok);
    },

    /**
     * Enter przy otwartym wykazie.
     *
     * Oddaje `false`, gdy mechanizm nie miał czego wybrać — wtedy pole
     * wypowiedzi robi swoje zwykłe zadanie zamiast niczego. Zdanie o powodzie
     * idzie do wątku, żeby Enter bez skutku nie był ciszą.
     */
    zatwierdz() {
      if (!pogodzStan()) return false;
      if (menu.wybierzWyrozniona()) return true;
      zglos(
        powod !== ''
          ? powod
          : 'Wykaz po ukośniku nie ma w tej chwili ani jednej pozycji pasującej do ' +
            'wpisanej frazy — nie było czego zatwierdzić.',
      );
      zamknij();
      return false;
    },

    zamknij,

    naZmiane(sluchacz) {
      sluchacze.push(sluchacz);
    },
  };
}
