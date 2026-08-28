import { ConfigScope } from '../../../shared/contract';
import { ETYKIETY_ZASIEGU } from './katalog-izolacji';

/** Zasięg czynny okna Punktów Izolacji — poziom, na którym reguła izolacji ma obowiązywać, wraz z bytem. */
export interface StanZasiegu {
  /** Poziom czynny; przed pierwszym wyborem — globalny, warstwa bazowa platformy. */
  zasieg(): ConfigScope;
  /** Identyfikator bytu poziomu; napis pusty znaczy „bez bytu". */
  bytZasiegu(): string;
  /** Byt w postaci pola żądania: `undefined` zamiast napisu pustego. */
  bytDoZadania(): string | undefined;
  /** Ustawia poziom i byt; powiadamia słuchaczy tylko przy zmianie. */
  ustaw(zasieg: ConfigScope, bytZasiegu?: string): void;
  /** Zdanie o zasięgu czynnym — jedno brzmienie dla wszystkich obszarów. */
  opis(): string;
  /** Subskrypcja zmiany zasięgu; zwraca odsubskrybowanie. */
  naZmiane(sluchacz: (zasieg: ConfigScope, bytZasiegu: string) => void): () => void;
}

export function utworzStanZasiegu(
  poczatkowy: ConfigScope = ConfigScope.Global,
  bytPoczatkowy = '',
): StanZasiegu {
  let czynny: ConfigScope = poczatkowy;
  let byt = bytPoczatkowy;
  const sluchacze = new Set<(zasieg: ConfigScope, bytZasiegu: string) => void>();

  return {
    zasieg: () => czynny,
    bytZasiegu: () => byt,
    bytDoZadania: () => (byt === '' ? undefined : byt),

    ustaw(zasieg, bytZasiegu = '') {
      if (czynny === zasieg && byt === bytZasiegu) return;
      czynny = zasieg;
      byt = bytZasiegu;
      for (const sluchacz of [...sluchacze]) sluchacz(czynny, byt);
    },

    opis() {
      const nazwa = ETYKIETY_ZASIEGU[czynny] ?? czynny;
      return byt === ''
        ? `poziom „${nazwa}", bez identyfikatora bytu`
        : `poziom „${nazwa}", byt ${byt}`;
    },

    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },
  };
}
