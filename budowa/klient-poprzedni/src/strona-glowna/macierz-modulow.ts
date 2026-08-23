import { KnownModuleIds, type Environment } from '../../../shared/contract';
import { kodZnany, type KodSrodowiska } from './pozycje-srodowisk';

/**
 * Macierz widoczności modułów po stronie klienta — odwzorowanie macierzy
 * `srodowisko_modul` z rdzenia, bez własnego wykazu par wpisanego w klienta.
 *
 * Dane niesie `environment.list` z `includeModules: true`: każde środowisko ma
 * pole `moduleCodes` z kodami modułów widocznych w jego bocznej nawigacji,
 * w kolejności wyświetlania. Osobnego zapytania o macierz nie ma.
 *
 * Macierz niczego nie blokuje. Zanim rdzeń odpowie, i gdy moduł nie stoi
 * w żadnym środowisku, `srodowiskoModulu` zwraca `null` — co znaczy „nie wiem,
 * dokąd", a nie „nie wolno". Co z tym zrobić, rozstrzyga miejsce wyboru.
 */

/** Środowisko sprowadzone do tego, co macierz o nim mówi. */
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
  /**
   * Środowisko, w którym moduł pracuje — albo `null`, gdy macierz jeszcze nie
   * przyszła lub moduł nie stoi w żadnej bocznej nawigacji.
   */
  srodowiskoModulu(kodModulu: string): KodSrodowiska | null;
  /**
   * Środowisko, od którego zaczyna się praca, gdy macierz nic nie wskazuje.
   *
   * Pierwsza karta wg kolejności z rdzenia, a przed pierwszą odpowiedzią —
   * pierwszy kod z kontraktu. Nie ma tu nazwy wpisanej ręcznie: `KnownModuleIds`
   * jest tym samym źródłem prawdy, z którego żyje strefa pierwsza.
   */
  srodowiskoPoczatkowe(): KodSrodowiska;
}

export function utworzMacierzModulow(): MacierzModulow {
  let wiersze: readonly WierszMacierzy[] = [];

  return {
    ustawZeSrodowisk(srodowiska) {
      // Wykaz pusty zostawia macierz poprzednią. Odmowa rdzenia i cisza mają
      // nie kasować tego, co już wiadomo — kafel prowadzący gdziekolwiek jest
      // lepszy niż kafel, który nagle przestał wiedzieć.
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
      // Remis: moduł bywa widoczny w kilku środowiskach naraz, a kontrakt nie
      // niesie środowiska macierzystego. Pierwszeństwo ma środowisko o tym samym
      // kodzie co moduł, a gdy takiego nie ma — pierwsze wg kolejności kart.
      // Reguła liczy się z danych; żadnego kodu nie ma tu wpisanego.
      const wlasne = kandydaci.find((wiersz) => wiersz.kod === kodModulu);
      return (wlasne ?? kandydaci[0]).kod;
    },

    srodowiskoPoczatkowe() {
      return wiersze[0]?.kod ?? KnownModuleIds[0];
    },
  };
}
