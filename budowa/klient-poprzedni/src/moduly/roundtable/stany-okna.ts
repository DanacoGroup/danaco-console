/**
 * Stan treści okien modułu Roundtable — wiązanie wspólnego nośnika
 * z `komponenty/stan-tresci.ts` z przedrostkiem klas `dr-`.
 *
 * Stan błędu jest tu konieczny: debata odmawia z powodów zwyczajnych — kanał
 * uczestnika nieczynny, tura zamknięta, stanowiska jeszcze nie ma — a okno
 * pokazujące wtedy pustą listę wypowiedzi mówiłoby „nikt nic nie powiedział"
 * zamiast „nie udało się zapytać".
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/** Nośnik stanu treści okna Roundtable. */
export function utworzStanTresci(): StanTresci {
  return utworzWspolny('dr');
}
