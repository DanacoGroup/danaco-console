import type { Window } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import { nanies, ustawieniaDomyslne, zdejmij, type UstawieniaOkna } from './klucze-ustawien';
import { obowiazujaceNieznane, type ObowiazujaceOkna } from './wartosc-obowiazujaca';

/** Stan sterowania widziany przez pojedyncze sterowanie w chwili odrysowania: okno, ustawienia i wartości obowiązujące. */
export interface MigawkaSterowania {
  /** Okno komunikacji potwierdzone przez rdzeń. */
  okno: Window;
  // Ustawienia zapisane na poziomie okna — to, co operator ustawił tutaj sam.
  ustawienia: UstawieniaOkna;
  // Ustawienia obowiązujące — to, czym naprawdę pojedzie model, po rozstrzygnięciu poziomów zasięgu.
  obowiazujace: ObowiazujaceOkna;
}

/** Stan kompletu sterowania jednego okna komunikacji, przyjmujący wyłącznie prawdę potwierdzoną przez rdzeń. */
export interface StanSterowania {
  /** Identyfikator okna, którego dotyczy komplet. */
  idOkna(): string;
  /** Bieżąca migawka stanu. */
  migawka(): MigawkaSterowania;
  /** Przyjmuje okno potwierdzone przez rdzeń. */
  przyjmijOkno(okno: Window): void;
  /** Przyjmuje pojedyncze ustawienie poziomu okna. */
  przyjmijUstawienie(klucz: string, wartosc: unknown): void;
  // Przyjmuje zdjęcie zapisu poziomu okna, czyli przywrócenie wartości domyślnej katalogu.
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
