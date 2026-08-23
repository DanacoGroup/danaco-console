/**
 * Stan treści okien modułu Workspace — wiązanie wspólnego bytu z przedrostkiem `dw-`.
 *
 * Nośnik stanu stoi w `komponenty/stan-tresci.ts`; tutaj ustalany jest wyłącznie
 * przedrostek klas CSS właściwy dla okien Workspace. Stan błędu niesie kod
 * i komunikat z kontraktu, aby okno po niepowodzeniu zapytania nie wyglądało
 * jak okno z pustą listą — pustka jest osobnym, poprawnym stanem.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/** Nośnik stanu treści okna Workspace. */
export function utworzStanTresci(): StanTresci {
  return utworzWspolny('dw');
}
