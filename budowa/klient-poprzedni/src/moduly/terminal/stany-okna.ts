/**
 * Stan treści okien modułu Terminal wiąże wspólny nośnik stanu z przedrostkiem klas dt- używanym w tym module.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/**
 * Nośnik stanu treści okna Terminal wraz z kodem kontraktu odróżniającym odmowę uprawnienia od usterki rdzenia.
 */
export function utworzStanTresci(): StanTresci {
  return utworzWspolny('dt');
}
