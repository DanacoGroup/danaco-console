/**
 * Stan treści okien modułu Automations — wiązanie wspólnego nośnika
 * (`komponenty/stan-tresci.ts`) z przedrostkiem klas `da`, bo klasy stanu należą
 * do arkusza modułu (`automations.css`).
 *
 * Różnica wobec nośnika wspólnego: `blad()` i `pusto()` zdejmują tu wcześniejsze
 * potwierdzenie. Wspólny nośnik przestawia wyłącznie pas komunikatu, przez co
 * odmowa rdzenia stawała obok potwierdzenia czynności poprzedniej — a zielone
 * zdanie nad odmową potwierdza czynność, która się nie odbyła.
 *
 * Zmiana siedzi w module, a nie w nośniku wspólnym, bo ten sam nośnik obsługuje
 * okna pozostałych modułów.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/** Nośnik stanu treści okna Automations. */
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
