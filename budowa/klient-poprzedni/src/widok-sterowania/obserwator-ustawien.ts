import type { Window } from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import {
  utworzStanSterowania,
  type MigawkaSterowania,
} from '../sterowanie/stan-sterowania';
import { utworzZmianeOkna } from '../sterowanie/zmiana-okna';
import { utworzZmianeUstawienia } from '../sterowanie/zmiana-ustawienia';

/** Obserwator odczytuje bieżące ustawienia okna na potrzeby podsumowania widoku, nie wywołując żadnej zmiany i nie dopisując reguł protokołu ponad katalog sterowania. */
export interface ObserwatorUstawien {
  /** Bieżąca migawka ustawień okna. */
  migawka(): MigawkaSterowania;
  /** Subskrypcja zmian potwierdzonych przez rdzeń. */
  naZmiane(sluchacz: (migawka: MigawkaSterowania) => void): Odsubskrybuj;
  /** Odłącza subskrypcje kanału. */
  rozlacz(): void;
}

/** Odbiorca komunikatów o losie zmiany ustawień; obserwator wyłącznie czyta, więc komunikat pomija bez działania. */
const POMIN_KOMUNIKAT = (): void => {};

export function utworzObserwatorUstawien(
  kanal: Kanal,
  okno: Window,
): ObserwatorUstawien {
  const stan = utworzStanSterowania(okno);

  // Obie drogi zmiany podpinają subskrypcje zdarzeń zmiany okna i konfiguracji.
  const zmianaOkna = utworzZmianeOkna(kanal, stan, POMIN_KOMUNIKAT);
  const zmianaUstawienia = utworzZmianeUstawienia(kanal, stan, POMIN_KOMUNIKAT);

  // Odczyt ustawień poziomu okna uzupełnia pola nieobecne w zdarzeniu aktualizacji okna.
  zmianaUstawienia.wczytaj();

  return {
    migawka: stan.migawka,
    naZmiane: stan.naZmiane,
    rozlacz() {
      zmianaOkna.rozlacz();
      zmianaUstawienia.rozlacz();
    },
  };
}
