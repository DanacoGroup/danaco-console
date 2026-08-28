import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';

/**
 * Obsada bibliotecznego komponentu drzewa menu dla wyborów lokalnych okna:
 * podaje mechanizmowi dane i nie odtwarza formy wyboru.
 */
/** Jedna pozycja wyboru — para klucz–nazwa wraz z opcjonalnym zdaniem opisu wyjaśniającym jej znaczenie. */
export interface PozycjaWyboruMenu {
  /** Klucz oddawany oknu; pusty jest poprawną wartością („bez wskazania"). */
  wartosc: string;
  /** Nazwa widoczna w wykazie i na uchwycie po wybraniu. */
  etykieta: string;
  /** Zdanie mówiące, co ta pozycja znaczy; pominięte znaczy „bez opisu". */
  opis?: string;
}

export interface WyborZMenu {
  /** Element montowany w formularzu albo w pasku narzędzi okna. */
  element: HTMLElement;
  /** Klucz pozycji wybranej; pusty łańcuch, gdy wybrano pozycję o pustym kluczu. */
  wartosc(): string;
  /** Ustawia wybór z zewnątrz; klucz spoza wykazu nie jest wybierany i oddaje odmowę. */
  ustawWartosc(klucz: string): boolean;
  /** Wymienia pozycje; wybór zostaje, o ile nadal istnieje w nowym wykazie pozycji. */
  ustawPozycje(pozycje: readonly PozycjaWyboruMenu[]): void;
  /** Zgłasza słuchacza wyboru dokonanego przez Operatora — i wyłącznie takiego, nie zmiany z kodu. */
  naZmiane(sluchacz: (klucz: string) => void): void;
}

/**
 * Napis stawiany na uchwycie wyboru, gdy wykaz pozycji jest pusty i rozwinięcie
 * menu nie miałoby czego pokazać Operatorowi.
 */
const UCHWYT_PUSTEGO = 'brak pozycji do wyboru';

/**
 * @param nastawa rodzajowa nazwa wyboru; wchodzi do atrybutu opisującego.
 * @param poczatkowe wykaz startowy; pusty znaczy wykaz wchodzący później.
 */
export function utworzWyborZMenu(
  nastawa: string,
  poczatkowe: readonly PozycjaWyboruMenu[] = [],
): WyborZMenu {
  let pozycje: readonly PozycjaWyboruMenu[] = [];
  let wybrany = '';
  let jestWybor = false;
  const sluchacze = new Set<(klucz: string) => void>();

  const menu = utworzMenuDrzewo({
    nastawa,
    naWybor: (klucz) => {
      wybrany = klucz;
      jestWybor = true;
      odrysuj();
      for (const sluchacz of [...sluchacze]) sluchacz(klucz);
    },
  });

  /** Pozycja wybrana albo `undefined`, gdy wykaz jej nie niesie. */
  function wybrana(): PozycjaWyboruMenu | undefined {
    if (!jestWybor) return undefined;
    return pozycje.find((pozycja) => pozycja.wartosc === wybrany);
  }

  function odrysuj(): void {
    const biezaca = wybrana();
    const drzewo: PozycjaMenu[] = pozycje.map((pozycja) => ({
      rodzaj: 'wybor',
      klucz: pozycja.wartosc,
      nazwa: pozycja.etykieta,
      ...(pozycja.opis === undefined ? {} : { opis: pozycja.opis }),
      wybrany: biezaca !== undefined && pozycja.wartosc === biezaca.wartosc,
    }));
    menu.ustaw(biezaca?.etykieta ?? UCHWYT_PUSTEGO, drzewo);
  }

  function ustawPozycje(nowe: readonly PozycjaWyboruMenu[]): void {
    pozycje = [...nowe];
    // Wybór utrzymany, o ile przetrwał wymianę wykazu.
    if (!jestWybor || !pozycje.some((pozycja) => pozycja.wartosc === wybrany)) {
      const pierwsza = pozycje[0];
      wybrany = pierwsza?.wartosc ?? '';
      jestWybor = pierwsza !== undefined;
    }
    odrysuj();
  }

  ustawPozycje(poczatkowe);

  return {
    element: menu.element,

    wartosc: () => (wybrana() === undefined ? '' : wybrany),

    ustawWartosc(klucz) {
      if (!pozycje.some((pozycja) => pozycja.wartosc === klucz)) return false;
      wybrany = klucz;
      jestWybor = true;
      odrysuj();
      return true;
    },

    ustawPozycje,

    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
    },
  };
}

/**
 * Wiersz formularza z podpisem nad wyborem w kształcie biblioteki komponentów;
 * nazwę nastawy niesie atrybut opisujący uchwytu, a podpis stoi dla oka.
 */
export function wierszWyboru(etykieta: string, wybor: WyborZMenu): HTMLElement {
  const podpis = document.createElement('span');
  podpis.className = 'dn-pole-etykieta';
  podpis.textContent = etykieta;

  const element = document.createElement('div');
  element.className = 'dn-pole dm-pole mp-pole-wyboru';
  element.append(podpis, wybor.element);
  return element;
}
