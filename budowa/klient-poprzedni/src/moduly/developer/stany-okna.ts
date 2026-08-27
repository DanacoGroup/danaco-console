/**
 * Stan treści okien modułu Developer — wiązanie wspólnego nośnika z przedrostkiem `mdev-`.
 *
 * Sam nośnik stoi w `komponenty/stan-tresci.ts`, a tutaj podstawiany jest
 * przedrostek klas. Stan błędu jest rozłączny ze stanem pustki, bo pustka
 * bywa poprawna.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/**
 * Nośnik stanu treści okna Developer, zbudowany na klasach z przedrostkiem `mdev-`.
 *
 * Okno dostaje stąd jeden element i przez niego pokazuje ładowanie, odmowę
 * odczytu albo pustkę drzewa.
 */
export function utworzStanTresci(): StanTresci {
  return utworzWspolny('mdev');
}
