import { ConfigScope } from '../../../shared/contract';
import { ETYKIETY_ZASIEGU } from './katalog-izolacji';

/**
 * Zasięg czynny okna Punktów Izolacji — poziom, na którym reguła izolacji ma
 * obowiązywać, wraz z identyfikatorem bytu tego poziomu.
 *
 * Stan jest wspólny dla całego okna z tego samego powodu, dla którego wspólna
 * jest warstwa (`stan-warstwy.ts`): poziom rozstrzyga, który zapis czyta i pisze
 * macierz izolacji, dokąd trafia przypisanie profilu i czego dotyczy podgląd
 * polityki efektywnej. Osobny selektor w każdym z tych miejsc pokazywałby obok
 * siebie wartości z trzech różnych poziomów pod jedną nazwą „izolacja".
 *
 * Poziom globalny jest warstwą bazową i jedynym, który nie potrzebuje bytu —
 * każdy węższy wskazuje byt (kod środowiska, identyfikator projektu, sesji,
 * nazwę roli). Byt pusty przy poziomie węższym nie jest błędem klienta: żądanie
 * idzie bez pola `scopeId`, a odmowę — jeżeli rdzeń bytu wymaga — nazywa on sam.
 *
 * Plik nie woła rdzenia i nie buduje ani jednego elementu widoku. Selektor stoi
 * w `panel-zasiegu.ts`, czytelnicy — w obszarach.
 */
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
