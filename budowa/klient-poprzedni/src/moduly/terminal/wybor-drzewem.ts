import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';

/**
 * Wybór jednokrotny modułu Terminal osadzony na bibliotece menu drzewa.
 */

/** Jedna pozycja wykazu wyboru: wartość, nazwa widoczna i opcjonalne zdanie o skutku, jaki wybór tej pozycji wywołuje w oknie. */
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
  // Pierwsza pozycja jest wyborem początkowym, jak natywny select bez atrybutu selected.
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
    // Wybór idzie wyłącznie z pozycji drzewa; wartość spoza wykazu tutaj nie trafia.
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
