import { POZYCJE_SRODOWISK, type KodSrodowiska, type PozycjaSrodowiska } from './pozycje-srodowisk';

/** Wykaz środowisk strony głównej jest jednym źródłem nazwy dla całego ekranu: dla kart środowisk i dla wierszy sesji w tle, zasilanym wykazem z rdzenia. */
export interface WykazSrodowisk {
  /** Bieżący wykaz — z rdzenia, a przed jego odpowiedzią zastany. */
  pozycje(): readonly PozycjaSrodowiska[];
  /** Podmienia wykaz na ten z rdzenia; pusty nie podmienia niczego — odmowa nie jest wykazem zerowym. */
  ustaw(pozycje: readonly PozycjaSrodowiska[]): void;
  /** Nazwa środowiska o podanym kodzie; kod spoza wykazu wraca dosłownie, nie jako brak wartości. */
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
