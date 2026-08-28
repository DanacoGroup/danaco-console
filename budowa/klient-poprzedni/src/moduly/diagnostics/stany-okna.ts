/**
 * Stan treści okien modułu Diagnostics — wiązanie wspólnego nośnika z przedrostkiem `dg-`.
 *
 * Sam nośnik stoi w `komponenty/stan-tresci.ts`. Stan błędu jest tu rozłączny
 * ze stanem pustki, bo moduł pokazuje właśnie odmowy wykonania komend.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/**
 * Nośnik stanu treści okna Diagnostics, na klasach z przedrostkiem `dg-`.
 *
 * Okno dostaje stąd jeden element i przez niego pokazuje ładowanie, odmowę
 * odczytu albo pusty wykaz błędów.
 */
export function utworzStanTresci(): StanTresci {
  return utworzWspolny('dg');
}
