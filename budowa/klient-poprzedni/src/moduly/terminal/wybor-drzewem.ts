import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';

/**
 * Wybór jednokrotny modułu Terminal osadzony na `komponenty/menu-drzewo.ts`.
 *
 * Plik nie rysuje ani jednego wiersza menu — strzałki, rozwijanie i znacznik
 * wyboru niesie mechanizm biblioteki, który oddaje wołającemu klucz pozycji
 * zamiast wykonywać wybór. Zostają tu dwie rzeczy, których mechanizm nie robi:
 *   1. pamięć wybranej wartości — `menu-drzewo` dostaje drzewo z zewnątrz przy
 *      każdym `ustaw`, więc ktoś musi wiedzieć, co jest wybrane teraz;
 *   2. przełożenie wykazów modułu (`profil-karty.ts`, `grupowanie-procesow.ts`)
 *      na `PozycjaMenu` — wykazy zostają danymi, a nie łańcuchem warunków.
 *
 * Uchwyt niesie wartość, nie nazwę nastawy: stoi na nim „Bash", a nie „Powłoka".
 * Nazwa rodzajowa idzie do `aria-label` i do podpowiedzi, bo czytnik ekranu musi
 * wiedzieć, czego wartość dotyczy.
 *
 * Wartość pusta jest wartością: `''` znaczy „bez zawężenia" w wykazie procesów
 * i ma własną pozycję („Wszystkie procesy"), inaczej z filtra nie byłoby drogi
 * powrotnej.
 */

/**
 * Jedna pozycja wykazu wyboru: wartość, nazwa widoczna, opcjonalne zdanie
 * o skutku. Trzeci człon to objaśnienie `[?]`; pominięty znaczy „ta pozycja
 * opisu nie ma".
 */
export type PozycjaWyboru = readonly [wartosc: string, nazwa: string, opis?: string];

export interface WyborDrzewem {
  /** Element montowany w pasku narzędzi okna. */
  element: HTMLElement;
  /** Wartość wybrana w tej chwili. */
  wartosc(): string;
  /** Nasłuch zmiany wyboru; wołany wyłącznie przy zmianie na inną wartość. */
  naZmiane(sluchacz: (wartosc: string) => void): void;
}

export interface OpcjeWyboru {
  /** Nazwa rodzajowa nastawy — idzie do `aria-label`, nie na uchwyt. */
  nastawa: string;
  /** Wykaz pozycji; pierwsza jest wyborem początkowym. */
  pozycje: readonly PozycjaWyboru[];
}

export function utworzWyborDrzewem(opcje: OpcjeWyboru): WyborDrzewem {
  const sluchacze = new Set<(wartosc: string) => void>();
  // Pierwsza pozycja jest wyborem początkowym — tak jak natywny `<select>` bez
  // `selected` bierze pierwszą opcję. Wykazy modułu są stałymi i puste nie bywają,
  // ale przy pustym uchwyt powie o pustce zdaniem mechanizmu zamiast przestać
  // reagować.
  let biezaca = opcje.pozycje[0]?.[0] ?? '';

  const menu = utworzMenuDrzewo({
    nastawa: opcje.nastawa,
    naWybor: (klucz) => {
      if (klucz === biezaca) return;
      biezaca = klucz;
      odrysuj();
      for (const sluchacz of [...sluchacze]) sluchacz(biezaca);
    },
  });

  /** Nazwa wybranej pozycji — to ona stoi na uchwycie. */
  function nazwaBiezacej(): string {
    const trafiona = opcje.pozycje.find(([wartosc]) => wartosc === biezaca);
    // Wybór idzie wyłącznie z pozycji drzewa, więc wartość spoza wykazu tu nie
    // trafia; gdyby trafiła, uchwyt pokaże ją samą zamiast pierwszej z brzegu
    // nazwy niezgodnej z nastawą.
    return trafiona?.[1] ?? biezaca;
  }

  function odrysuj(): void {
    const drzewo: PozycjaMenu[] = opcje.pozycje.map(([wartosc, nazwa, opis]) => ({
      rodzaj: 'wybor',
      klucz: wartosc,
      nazwa,
      ...(opis === undefined ? {} : { opis }),
      wybrany: wartosc === biezaca,
    }));
    menu.ustaw(nazwaBiezacej(), drzewo);
  }

  odrysuj();

  return {
    element: menu.element,
    wartosc: () => biezaca,
    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
    },
  };
}
