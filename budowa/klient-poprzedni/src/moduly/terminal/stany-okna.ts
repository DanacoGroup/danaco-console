/**
 * Stan treści okien modułu Terminal — wiązanie wspólnego bytu z przedrostkiem `dt-`.
 *
 * Sam nośnik stoi w `komponenty/stan-tresci.ts`; tutaj zostaje wyłącznie wybór
 * przedrostka klas.
 *
 * Okno bez stanu błędu jest niegotowe, a w terminalu kosztuje to więcej niż
 * gdzie indziej: okno, które po odmowie rdzenia pokazuje pustą listę procesów,
 * mówi Operatorowi „nic nie biegnie” wtedy, gdy biegnie kompilacja, o którą nie
 * udało się zapytać. Kod kontraktu w zdaniu błędu odróżnia odmowę uprawnienia
 * od usterki rdzenia.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/** Nośnik stanu treści okna Terminal. */
export function utworzStanTresci(): StanTresci {
  return utworzWspolny('dt');
}
