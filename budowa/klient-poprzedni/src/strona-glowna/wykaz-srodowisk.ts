import { POZYCJE_SRODOWISK, type KodSrodowiska, type PozycjaSrodowiska } from './pozycje-srodowisk';

/**
 * Wykaz środowisk strony głównej — jedno źródło nazwy dla całego ekranu:
 * i dla kart środowisk, i dla wierszy sesji w tle.
 *
 * Wykaz jest obiektem, a nie funkcją, bo nazwę trzeba odczytać w chwili
 * rysowania wiersza, a nie w chwili importu modułu — odpowiedź rdzenia
 * przychodzi później niż pierwsza klatka. `okno-strona-glowna.ts` spina go
 * z `strefa-srodowisk.ustawWykaz`, więc zasilenie kart bez zasilenia wykazu
 * jest niewykonalne.
 *
 * Do pierwszej odpowiedzi rdzenia wykaz niesie `POZYCJE_SRODOWISK` — tę samą
 * stałą, z której rysują się karty. Jedna treść zastana, nie dwie.
 */
export interface WykazSrodowisk {
  /** Bieżący wykaz — z rdzenia, a przed jego odpowiedzią zastany. */
  pozycje(): readonly PozycjaSrodowiska[];
  /**
   * Podmienia wykaz na ten z rdzenia. Wykaz pusty nie podmienia niczego:
   * odmowa i cisza nie są wykazem zerowym.
   */
  ustaw(pozycje: readonly PozycjaSrodowiska[]): void;
  /**
   * Nazwa środowiska o podanym kodzie.
   *
   * Kod spoza wykazu wraca dosłownie, a nie jako `undefined`: rdzeń mógł oddać
   * środowisko, o którym klient nie wie, a wiersz sesji ma wtedy pokazać kod
   * zamiast przemilczeć środowisko. Wykaz jest informacyjny, nie jest bramą.
   */
  nazwa(kod: string | undefined): string | undefined;
  /** Kod środowiska, o ile wykaz go niesie; inaczej `null`. */
  kod(kod: string | undefined): KodSrodowiska | null;
}

export function utworzWykazSrodowisk(
  zastane: readonly PozycjaSrodowiska[] = POZYCJE_SRODOWISK,
): WykazSrodowisk {
  let biezace = zastane;

  function znajdz(kod: string | undefined): PozycjaSrodowiska | undefined {
    if (kod === undefined || kod === '') return undefined;
    return biezace.find((pozycja) => pozycja.kod === kod);
  }

  return {
    pozycje: () => biezace,

    ustaw(pozycje) {
      if (pozycje.length === 0) return;
      biezace = pozycje;
    },

    nazwa(kod) {
      if (kod === undefined || kod === '') return undefined;
      return znajdz(kod)?.nazwa ?? kod;
    },

    kod(kod) {
      return znajdz(kod)?.kod ?? null;
    },
  };
}
