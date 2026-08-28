/**
 * Stan treści okien modułu Roundtable — wiązanie wspólnego nośnika z przedrostkiem klas dr,
 * ze stanem błędu koniecznym przy odmowie debaty.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/** Nośnik stanu treści okna Roundtable, złożony z nośnika wspólnego przy pomocy przedrostka klas modułu. */
export function utworzStanTresci(): StanTresci {
  return utworzWspolny('dr');
}
