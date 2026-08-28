import type { Window } from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import type { KomunikatZmiany } from '../sterowanie/komunikat-zmiany';
import { utworzStanSterowania } from '../sterowanie/stan-sterowania';
import { utworzZmianeOkna } from '../sterowanie/zmiana-okna';
import { utworzZmianeUstawienia } from '../sterowanie/zmiana-ustawienia';
import type { StanPrzelacznika, ZaleznosciPrzelacznika } from './przelacznik-srodowiska';

/**
 * Port przełącznika środowiska do kompletu sterowania okna podpina migawkę, subskrypcję zmian i wysyłkę pod stan rdzenia, tak jak katalog akcji zamienia kanał na źródło pozycji panelu.
 */
export interface ZrodloSrodowiska extends ZaleznosciPrzelacznika {
  /** Odłącza subskrypcje kanału i stanu. */
  rozlacz(): void;
}

/**
 * Nazwa zmiany w komunikacie o jej losie — ta sama, którą wysyła lista wyboru
 * w szufladzie (`sterowanie/srodowisko-wykonania.ts`). Odmowa rdzenia dociera
 * więc do operatora tym samym zdaniem, niezależnie od tego, którym z dwóch
 * widoków jej dotknął.
 */
const NAZWA = 'Środowisko wykonania';

export function utworzZrodloSrodowiska(kanal: Kanal, okno: Window): ZrodloSrodowiska {
  const stan = utworzStanSterowania(okno);

  /** Kolejka wywołań czekających na los zmiany — przekład między wywołaniem zwrotnym a obietnicą. */
  const oczekujace: ((komunikat: KomunikatZmiany) => void)[] = [];

  const zmiana = utworzZmianeOkna(kanal, stan, (komunikat) => {
    oczekujace.shift()?.(komunikat);
  });

  // Ustawienia okna czyta port, ale ich nie zapisuje — nazwę hosta nadaje pole w szufladzie sterowania.
  const ustawienia = utworzZmianeUstawienia(kanal, stan, () => undefined);
  ustawienia.wczytaj();

  const odsubskrybuj: Odsubskrybuj[] = [];

  return {
    migawka(): StanPrzelacznika {
      const biezaca = stan.migawka();
      return {
        srodowisko: biezaca.okno.executionEnv,
        host: biezaca.ustawienia.hostWykonania,
      };
    },

    naZmiane(sluchacz) {
      odsubskrybuj.push(stan.naZmiane(() => sluchacz()));
    },

    /** Odmowa rdzenia dochodzi dosłownie, tym samym zdaniem co pasek komunikatów przy zmianie z szuflady. */
    zastosuj(srodowisko) {
      return new Promise<void>((spelnij, odrzuc) => {
        oczekujace.push((komunikat) => {
          if (komunikat.udany) spelnij();
          else odrzuc(new Error(komunikat.tresc));
        });
        zmiana.zastosuj(NAZWA, { executionEnv: srodowisko });
      });
    },

    rozlacz() {
      for (const zdejmij of odsubskrybuj.splice(0)) zdejmij();
      zmiana.rozlacz();
      ustawienia.rozlacz();
    },
  };
}
