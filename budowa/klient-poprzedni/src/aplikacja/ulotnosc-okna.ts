import { MODUL as MODUL_AGENTS } from '../moduly/agents/indeks';
import { testowanyAgent } from '../moduly/agents/testowany-agent';
import { politykaModulu, type PolitykaUlotnosci } from '../rozmowa/ulotnosc';

/**
 * Złożenie polityki pamięci rozmowy dla modułu okna. Pamięć sesyjną modułu
 * niesie jego profil (`okno-komunikacji/profil-modulu.ts`, pole
 * `pamiecSesyjna`). Tutaj dochodzi kontekst roboczy, którego zmiana kończy
 * rozmowę; ma go moduł Agents.
 */
export function politykaOknaModulu(kod: string): PolitykaUlotnosci {
  return politykaModulu(kod, kod === MODUL_AGENTS.kod ? testowanyAgent : null);
}
