/**
 * Stan treści okien modułu Diagnostics — wiązanie wspólnego nośnika z przedrostkiem `dg-`.
 *
 * Sam nośnik stoi w `komponenty/stan-tresci.ts`; tutaj podstawiany jest wyłącznie
 * przedrostek klas. Stan błędu jest rozłączny ze stanem pustki i w tym module ma
 * znaczenie szczególne: port `Diagnostyka` jest w rdzeniu odbiorcą odmów wykonania
 * komend, więc Errors Panel pokazujący po nieudanym odczycie pusty wykaz ukrywa
 * odmowy dwa razy — cudze i swoją własną. Pustka bywa zaś poprawna: instalacja bez
 * ani jednego błędu w zakresie czasu nie jest usterką.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/** Nośnik stanu treści okna Diagnostics. */
export function utworzStanTresci(): StanTresci {
  return utworzWspolny('dg');
}
