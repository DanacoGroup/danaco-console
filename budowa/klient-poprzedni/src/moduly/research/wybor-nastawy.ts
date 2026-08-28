import { utworzMenuDrzewo, type PozycjaMenu } from '../../komponenty/menu-drzewo';

/** Jedna pozycja wyboru — para klucz-nazwa z opcjonalnym opisem — wśród obsady bibliotecznego komponentu menu-drzewo dla trzech nastaw modułu Research. */
export interface PozycjaNastawy {
  /** Wartość jadąca do rdzenia — napis wyliczenia kontraktu. */
  wartosc: string;
  /** Nazwa widoczna w wykazie i na uchwycie po wybraniu. */
  etykieta: string;
  /** Zdanie mówiące, co ta pozycja znaczy; pominięte znaczy „bez opisu". */
  opis?: string;
}

export interface WyborNastawy {
  /** Element montowany w formularzu okna. */
  element: HTMLElement;
  /** Wartość bieżąca — ta sama, którą niesie uchwyt. */
  wartosc(): string;
  /** Ustawia wybór z zewnątrz; klucz spoza wykazu nie jest wybierany i oddaje odmowę wyboru. */
  ustawWartosc(klucz: string): boolean;
}

/**
 * @param nastawa nazwa rodzajowa wyboru, idąca do etykiety dostępności, nie na ekran
 * @param pozycje wykaz stały, którego pierwsza pozycja jest wyborem początkowym
 * @param naZmiane wołane wyłącznie po wyborze Operatora, nie po ustawieniu z kodu
 */
export function utworzWyborNastawy(
  nastawa: string,
  pozycje: readonly PozycjaNastawy[],
  naZmiane?: (wartosc: string) => void,
): WyborNastawy {
  let wybrany = pozycje[0]?.wartosc ?? '';

  const menu = utworzMenuDrzewo({
    nastawa,
    naWybor: (klucz) => {
      wybrany = klucz;
      odrysuj();
      naZmiane?.(klucz);
    },
  });

  function odrysuj(): void {
    const biezaca = pozycje.find((pozycja) => pozycja.wartosc === wybrany);
    const drzewo: PozycjaMenu[] = pozycje.map((pozycja) => ({
      rodzaj: 'wybor',
      klucz: pozycja.wartosc,
      nazwa: pozycja.etykieta,
      ...(pozycja.opis === undefined ? {} : { opis: pozycja.opis }),
      wybrany: pozycja.wartosc === wybrany,
    }));
    // Uchwyt niesie wartość bieżącą; wykaz pusty mówi to zdaniem, a nie pustym
    // przyciskiem bez znaczenia.
    menu.ustaw(biezaca?.etykieta ?? 'brak pozycji do wyboru', drzewo);
  }

  odrysuj();

  return {
    element: menu.element,
    wartosc: () => wybrany,

    ustawWartosc(klucz) {
      if (!pozycje.some((pozycja) => pozycja.wartosc === klucz)) return false;
      wybrany = klucz;
      odrysuj();
      return true;
    },
  };
}

/**
 * Wiersz formularza: podpis nad sterem, w kształcie pola biblioteki, niezwiązany z uchwytem
 * etykietą, bo uchwyt jest przyciskiem.
 */
export function wierszNastawy(etykieta: string, wybor: WyborNastawy): HTMLElement {
  const podpis = document.createElement('span');
  podpis.className = 'dn-pole-etykieta';
  podpis.textContent = etykieta;

  const element = document.createElement('div');
  element.className = 'dn-pole mr-nastawa';
  element.append(podpis, wybor.element);
  return element;
}
