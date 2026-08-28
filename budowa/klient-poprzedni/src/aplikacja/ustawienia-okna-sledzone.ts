import type { Window } from '../../../shared/contract';

/**
 * Śledzenie ustawień okna między zdarzeniami: zdarzenie `window.changed` niesie
 * okno po zmianie, bez pola nazywającego zmianę i bez pola nazywającego jej
 * sprawcę, więc plik odtwarza jedno i drugie porównaniem migawek pól `Window`.
 */
export interface MigawkaUstawien {
  modul: string;
  kanal: string;
  agent: string;
  rola: string;
  zasieg: string;
  tryb: string;
}

export function migawka(okno: Window): MigawkaUstawien {
  return {
    modul: okno.moduleId,
    kanal: okno.modelChannelId,
    agent: okno.agentId ?? '',
    rola: okno.windowRole,
    zasieg: okno.executionEnv,
    tryb: okno.permissionMode,
  };
}

/**
 * Odcisk stanu okna jest kluczem porównania zdarzenia z własną odpowiedzią,
 * więc wchodzą do niego wyłącznie pola, którymi da się sterować komendą:
 * identyfikator okna oraz jego ustawienia.
 */
export function odciskOkna(okno: Window): string {
  const u = migawka(okno);
  return [okno.id, u.modul, u.kanal, u.agent, u.rola, u.zasieg, u.tryb].join('|');
}

/**
 * Zdania o ustawieniach, które naprawdę się zmieniły.
 *
 * Modułu w tym wykazie nie ma: ma własny rodzaj posunięcia i własny skutek
 * (podążanie nawigacji), więc meldowany tu byłby po raz drugi.
 */
export function roznice(przed: MigawkaUstawien, po: MigawkaUstawien): string[] {
  const opisy: Array<[keyof MigawkaUstawien, string]> = [
    ['kanal', 'kanał modelu'],
    ['agent', 'eksperta okna'],
    ['rola', 'rolę okna'],
    ['zasieg', 'zasięg wykonania'],
    ['tryb', 'tryb uprawnień'],
  ];
  const zdania: string[] = [];
  for (const [pole, nazwa] of opisy) {
    if (przed[pole] === po[pole]) continue;
    const wartosc = po[pole] === '' ? 'brak' : po[pole];
    zdania.push(`Ustawiony ${nazwa}: ${wartosc}`);
  }
  return zdania;
}
