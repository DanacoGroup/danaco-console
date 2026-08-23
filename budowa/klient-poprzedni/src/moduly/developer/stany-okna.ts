/**
 * Stan treści okien modułu Developer — wiązanie wspólnego nośnika z przedrostkiem `mdev-`.
 *
 * Sam nośnik stoi w `komponenty/stan-tresci.ts`; tutaj podstawiany jest wyłącznie
 * przedrostek klas. Stan błędu jest rozłączny ze stanem pustki: okno, które po
 * odmowie pokazuje puste drzewo, mówi „katalog jest pusty” zamiast „nie udało się
 * odczytać drzewa”, a pustka bywa poprawna — repozytorium bez zmian do zatwierdzenia
 * nie jest usterką.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/** Nośnik stanu treści okna Developer. */
export function utworzStanTresci(): StanTresci {
  return utworzWspolny('mdev');
}
