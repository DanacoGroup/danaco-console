import type { ToolCatalogEntry } from '../../../shared/contract';
import { utworzMenuDrzewo, type LiscMenu } from '../komponenty/menu-drzewo';
import { utworzRejestrOstatnich } from './ostatnio-uzyte';
import {
  czyPowolanieNarzedzia,
  type DolozenieNarzedzia,
  type ZrodloWykazuUkosnika,
} from './zrodlo-wykazu-ukosnika';

/**
 * Wykaz po ukośniku — obsada mechanizmu z biblioteki, pracująca w trybie bez uchwytu.
 */

/**
 * Napis na uchwycie; w trybie bez uchwytu nie trafia na ekran, tylko do warstwy
 * dostępności ARIA czytnika.
 */
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
  // Podpina słuchacza zmiany stanu wykazu, bo pole wypowiedzi nie zauważy samo jego
  // zwinięcia.
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
    // Trzy nastawienia trybu płaskiego: wyzwalacz z klawiatury, wykaz nad polem, opis zamiast
    // wielu.
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

  // Wykaz ułożony świeżością — jedyne, co ten plik wnosi do porządku pozycji.
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

  // Jedna pozycja wykazu jako liść mechanizmu; wybrany znaczy dołożenie, nie wyróżnienie.
  function liscWykazu(pozycja: ToolCatalogEntry): LiscMenu {
    return {
      rodzaj: 'wybor',
      klucz: pozycja.name,
      nazwa: pozycja.name,
      // Druga nazwa wchodzi tylko wtedy, gdy naprawdę jest inna niż nazwa pełna.
      ...(pozycja.shortName === '' || pozycja.shortName === pozycja.name
        ? {}
        : { nazwaKrotka: pozycja.shortName }),
      opis: opisWiersza(pozycja),
      wybrany: pozycja.attached === true,
    };
  }

  // Opis pełnym zdaniem tego, po co pozycja jest i kiedy jej użyć, ze zdaniem o skutku
  // wyboru.
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
      // Wykonania nie udajemy: bez znanego kształtu żądania Operator dostaje zdanie o
      // niepowodzeniu.
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

  // Godzi stan obsady ze stanem mechanizmu, który zwija się także sam, bez udziału obsady.
  function pogodzStan(): boolean {
    if (czynny && !mechanizmRozwiniety()) zamknij();
    return czynny;
  }

  return {
    element,

    otwarty: () => czynny && mechanizmRozwiniety(),

    otworz() {
      czynny = true;
      // Wykaz rozwija się najpierw z tym, co już jest; odpowiedź rdzenia przychodzi później.
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

    // Enter przy otwartym wykazie oddaje fałsz, gdy nie było czego wybrać, ze zdaniem w wątku.
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
