/**
 * Stan treści okien ról modułu Multitasking — wiązanie wspólnego nośnika
 * z przedrostkiem klas `dm-`.
 *
 * Sam nośnik stoi w `komponenty/stan-tresci.ts`, a tutaj zostaje wyłącznie
 * wybór przedrostka, dzięki czemu okna ról pokazują te stany jednakowo.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/**
 * Nośnik stanu treści okna roli, zbudowany na klasach z przedrostkiem `dm-`.
 *
 * Okno roli dostaje stąd jeden element i przez niego pokazuje ładowanie,
 * odmowę odczytu albo pustkę wykazu.
 */
export function utworzStanTresci(): StanTresci {
  return utworzWspolny('dm');
}
