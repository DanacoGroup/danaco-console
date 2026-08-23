import type { Window } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import { nanies, ustawieniaDomyslne, zdejmij, type UstawieniaOkna } from './klucze-ustawien';
import { obowiazujaceNieznane, type ObowiazujaceOkna } from './wartosc-obowiazujaca';

/** Stan sterowania widziany przez pojedyncze sterowanie w chwili odrysowania. */
export interface MigawkaSterowania {
  /** Okno komunikacji potwierdzone przez rdzeń. */
  okno: Window;
  /**
   * Ustawienia zapisane na poziomie okna — to, co Operator ustawił tutaj.
   * Menu i suwaki pokazują tę wartość jako swoje położenie.
   */
  ustawienia: UstawieniaOkna;
  /**
   * Ustawienia obowiązujące — to, czym naprawdę pojedzie model, po
   * rozstrzygnięciu poziomów zasięgu. To druga prawda, nie ta sama: nakład
   * zapisany globalnie obowiązuje okno, w którym nikt go nie ustawił.
   * Etykieta steru niesie tę wartość (`wartosc-obowiazujaca.ts`).
   */
  obowiazujace: ObowiazujaceOkna;
}

/**
 * Stan kompletu sterowania **jednego** okna komunikacji.
 *
 * Egzemplarz powstaje osobno dla każdego okna i domyka się na jego
 * identyfikatorze. Nie ma tu ani jednej zmiennej na poziomie modułu, więc dwa
 * okna otwarte obok siebie mają dwa niezależne stany: zmiana w jednym nie ma
 * żadnej drogi, którą mogłaby dosięgnąć drugiego.
 *
 * Stan przyjmuje wyłącznie prawdę potwierdzoną przez rdzeń — wynik komendy
 * albo zdarzenie zmiany. Komunikat dotyczący innego okna jest pomijany; nie
 * jest to brama wykonania, lecz kierowanie ruchu do właściwego adresata.
 */
export interface StanSterowania {
  /** Identyfikator okna, którego dotyczy komplet. */
  idOkna(): string;
  /** Bieżąca migawka stanu. */
  migawka(): MigawkaSterowania;
  /** Przyjmuje okno potwierdzone przez rdzeń. */
  przyjmijOkno(okno: Window): void;
  /** Przyjmuje pojedyncze ustawienie poziomu okna. */
  przyjmijUstawienie(klucz: string, wartosc: unknown): void;
  /**
   * Przyjmuje zdjęcie zapisu poziomu okna — przywrócenie wartości domyślnej
   * katalogu (`config.reset`, zdarzenie `config.changed` o rodzaju `deleted`).
   * Bez tej drogi wpis skasowany przez rdzeń wchodziłby tu jak świeży zapis.
   */
  przyjmijUsuniecie(klucz: string): void;
  /** Przyjmuje komplet nastaw obowiązujących oddany przez rdzeń. */
  przyjmijObowiazujace(komplet: ObowiazujaceOkna): void;
  /** Ogłasza bieżącą migawkę bez zmiany treści — powrót widoku do stanu potwierdzonego. */
  odswiez(): void;
  /** Subskrypcja zmian stanu. */
  naZmiane(sluchacz: (migawka: MigawkaSterowania) => void): Odsubskrybuj;
}

export function utworzStanSterowania(poczatkowe: Window): StanSterowania {
  const zmiany = utworzMagistrale<MigawkaSterowania>();
  const identyfikator = poczatkowe.id;
  let okno = poczatkowe;
  let ustawienia = ustawieniaDomyslne();
  let obowiazujace = obowiazujaceNieznane();

  function migawka(): MigawkaSterowania {
    return { okno, ustawienia, obowiazujace };
  }

  return {
    idOkna: () => identyfikator,

    migawka,

    przyjmijOkno(nowe) {
      if (nowe.id !== identyfikator) return;
      okno = nowe;
      zmiany.oglos(migawka());
    },

    przyjmijUstawienie(klucz, wartosc) {
      ustawienia = nanies(ustawienia, klucz, wartosc);
      zmiany.oglos(migawka());
    },

    przyjmijUsuniecie(klucz) {
      ustawienia = zdejmij(ustawienia, klucz);
      zmiany.oglos(migawka());
    },

    przyjmijObowiazujace(komplet) {
      obowiazujace = komplet;
      zmiany.oglos(migawka());
    },

    odswiez() {
      zmiany.oglos(migawka());
    },

    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}
