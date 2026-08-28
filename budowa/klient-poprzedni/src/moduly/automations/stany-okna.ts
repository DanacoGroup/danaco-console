/**
 * Stan treści okien modułu Automations — wiązanie wspólnego nośnika
 * (`komponenty/stan-tresci.ts`) z przedrostkiem klas `da` arkusza
 * `automations.css`. Metody `blad()` i `pusto()` zdejmują tu wcześniejsze
 * potwierdzenie.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/**
 * Nośnik stanu treści okna Automations. Zwraca nośnik wspólny z przedrostkiem
 * `da`, w którym `blad()` i `pusto()` najpierw gaszą pas potwierdzenia,
 * a dopiero potem stawiają własny komunikat.
 */
export function utworzStanTresci(): StanTresci {
  const wspolny = utworzWspolny('da');
  return {
    ...wspolny,

    blad(zdanie, powod) {
      wspolny.potwierdzenie('', false);
      wspolny.blad(zdanie, powod);
    },

    pusto(zdanie) {
      wspolny.potwierdzenie('', false);
      wspolny.pusto(zdanie);
    },
  };
}
