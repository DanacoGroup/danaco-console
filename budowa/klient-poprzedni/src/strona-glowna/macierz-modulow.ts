import { KnownModuleIds, type Environment } from '../../../shared/contract';
import { kodZnany, type KodSrodowiska } from './pozycje-srodowisk';

/** Macierz widoczności modułów po stronie klienta odwzorowuje macierz rdzenia bez wykazu par środowisk. */

/** Środowisko sprowadzone do tego, co macierz widoczności modułów mówi o nim: kod, kolejność i lista modułów. */
interface WierszMacierzy {
  kod: KodSrodowiska;
  /** Kolejność karty środowiska z kontraktu, liczona od 1. */
  kolejnosc: number;
  /** Kody modułów widocznych w bocznej nawigacji tego środowiska. */
  moduly: readonly string[];
}

export interface MacierzModulow {
  /** Przyjmuje wykaz środowisk z rdzenia; wykaz pusty niczego nie kasuje. */
  ustawZeSrodowisk(srodowiska: readonly Environment[]): void;
  /** Środowisko, w którym moduł pracuje, albo `null`, gdy macierz jeszcze nie przyszła. */
  srodowiskoModulu(kodModulu: string): KodSrodowiska | null;
  /** Środowisko, od którego zaczyna się praca, gdy macierz nic nie wskazuje. */
  srodowiskoPoczatkowe(): KodSrodowiska;
}

export function utworzMacierzModulow(): MacierzModulow {
  let wiersze: readonly WierszMacierzy[] = [];

  return {
    ustawZeSrodowisk(srodowiska) {
      // Wykaz pusty zostawia macierz poprzednią, zamiast kasować to, co już wiadomo.
      if (srodowiska.length === 0) return;
      wiersze = srodowiska
        .filter((srodowisko) => kodZnany(srodowisko.code))
        .map((srodowisko) => ({
          kod: srodowisko.code as KodSrodowiska,
          kolejnosc: srodowisko.order,
          moduly: srodowisko.moduleCodes ?? [],
        }))
        .sort((a, b) => a.kolejnosc - b.kolejnosc);
    },

    srodowiskoModulu(kodModulu) {
      const kandydaci = wiersze.filter((wiersz) => wiersz.moduly.includes(kodModulu));
      if (kandydaci.length === 0) return null;
      // Remis rozstrzyga środowisko o tym samym kodzie co moduł, inaczej pierwsze wg kolejności.
      const wlasne = kandydaci.find((wiersz) => wiersz.kod === kodModulu);
      return (wlasne ?? kandydaci[0]).kod;
    },

    srodowiskoPoczatkowe() {
      return wiersze[0]?.kod ?? KnownModuleIds[0];
    },
  };
}
