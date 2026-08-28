import { Command, type Module } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';

/**
 * Katalog modułów rdzenia — czytająca pamięć podręczna nad odczytem wykazu modułów z
 * rdzenia aplikacji.
 */
export interface RejestrModulow {
  /** Moduły znane w tej chwili; pusta lista, dopóki rdzeń nie odpowie. */
  moduly(): Module[];
  /** Czy rdzeń odpowiedział na `module.list` choć raz. */
  odpowiedzOtrzymana(): boolean;
  /** Zamawia wykaz z rdzenia. */
  odswiez(): void;
  /** Subskrypcja zmian wykazu. */
  naZmiane(sluchacz: (moduly: Module[]) => void): Odsubskrybuj;
}

export function utworzRejestrModulow(kanal: Kanal): RejestrModulow {
  const zmiany = utworzMagistrale<Module[]>();
  let wykaz: Module[] = [];
  let odpowiedziano = false;

  return {
    moduly: () => wykaz,

    odpowiedzOtrzymana: () => odpowiedziano,

    odswiez() {
      kanal.wyslij(Command.ModuleList, {}, (wynik) => {
        // Odpowiedź odnotowujemy także przy odmowie, żeby widok pokazał stan pusty zamiast
        // czekania.
        odpowiedziano = true;
        const odebrane = wynik.wynik?.modules;
        wykaz = wynik.udany && odebrane !== undefined ? odebrane : wykaz;
        zmiany.oglos(wykaz);
      });
    },

    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}
