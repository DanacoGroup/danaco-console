/**
 * Nośnik stanu treści okien modułu Workspace, oparty na implementacji wspólnej
 * i posługujący się przedrostkiem klas CSS „dw-" właściwym dla tych okien.
 */
import { utworzStanTresci as utworzWspolny, type StanTresci } from '../../komponenty/stan-tresci';

export type { StanTresci };

/** Tworzy nośnik stanu treści okna modułu Workspace na bazie implementacji wspólnej, z przedrostkiem klas CSS „dw-". */
export function utworzStanTresci(): StanTresci {
  return utworzWspolny('dw');
}
