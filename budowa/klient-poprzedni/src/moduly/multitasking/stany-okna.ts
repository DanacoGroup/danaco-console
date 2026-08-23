/**
 * Stan treści okien ról modułu Multitasking — wiązanie wspólnego nośnika
 * z przedrostkiem klas `dm-`.
 *
 * Sam nośnik stoi w `komponenty/stan-tresci.ts`; tutaj zostaje wyłącznie wybór
 * przedrostka, dzięki czemu okna ról pokazują ładowanie, odmowę i pustkę
 * jednakowo. Okno, które po odmowie odczytu pokazałoby wyzerowany licznik
 * zamiast błędu, sugerowałoby zatrzymany bieg tam, gdzie bieg trwa.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/** Nośnik stanu treści okna roli. */
export function utworzStanTresci(): StanTresci {
  return utworzWspolny('dm');
}
