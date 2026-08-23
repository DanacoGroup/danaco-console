import { ProcessInitiator, TerminalProcessStatus, type TerminalProcess } from '../../../../shared/contract';
import type { PozycjaWyboru } from './wybor-drzewem';

/**
 * Wyliczenia filtrów i grupowanie wykazu Process Monitora.
 *
 * Filtr to co innego niż grupowanie. Filtr stanu i inicjatora jedzie do rdzenia
 * parametrem `terminal.process.list`, bo tak stanowi kontrakt i tylko tak wynik
 * jest spójny z dziennikiem rdzenia. Grupowanie nie zmienia zbioru, wyłącznie
 * porządek wyświetlania, więc zostaje w kliencie. Różnicę niesie objaśnienie
 * przy każdej pozycji, bo wykazy karmią menu z biblioteki.
 *
 * Wartość pusta ma własną pozycję: „Wszystkie procesy" i „Każdy inicjator" to
 * `''`, czyli brak zawężenia, a nie brak wyboru. Bez niej filtr byłby drogą
 * w jedną stronę.
 */

/** Filtry z panelu akcji wykazu modułów. */
export const FILTRY: readonly PozycjaWyboru[] = [
  ['', 'Wszystkie procesy', 'Bez zawężenia stanu — rdzeń oddaje komplet rejestru.'],
  [TerminalProcessStatus.Running, 'Aktywne', 'Procesy, które w chwili odczytu biegły.'],
  [TerminalProcessStatus.Finished, 'Zakończone', 'Domknięte bez błędu.'],
  [TerminalProcessStatus.Failed, 'Zakończone błędem', 'Domknięte niezerowym kodem wyjścia.'],
  [TerminalProcessStatus.Stopped, 'Zatrzymane', 'Przerwane sygnałem — kodu wyjścia zwykle nie mają.'],
];

export const INICJATORZY: readonly PozycjaWyboru[] = [
  ['', 'Każdy inicjator', 'Bez zawężenia — Operator i model razem.'],
  [ProcessInitiator.Model, 'Uruchomione przez AI', 'Polecenia wydane przez model w turze.'],
  [ProcessInitiator.Operator, 'Uruchomione przez Operatora', 'Polecenia wydane ręcznie z okna.'],
];

export const GRUPOWANIA: readonly PozycjaWyboru[] = [
  ['brak', 'Bez grupowania', 'Jeden wykaz w kolejności stanu okna.'],
  ['karta', 'Grupuj po karcie', 'Porządek wyświetlania; do rdzenia nic nie jedzie.'],
  ['inicjator', 'Grupuj po inicjatorze', 'Porządek wyświetlania; do rdzenia nic nie jedzie.'],
  ['status', 'Grupuj po stanie', 'Porządek wyświetlania; do rdzenia nic nie jedzie.'],
];

/** Grupuje wykaz po wskazanym kluczu; „brak” daje jedną grupę bez podpisu. */
export function pogrupuj(procesy: readonly TerminalProcess[], klucz: string): Array<[string, TerminalProcess[]]> {
  if (klucz === 'brak') return [['', [...procesy]]];
  const grupy = new Map<string, TerminalProcess[]>();
  for (const proces of procesy) {
    const nazwa = nazwaGrupy(proces, klucz);
    const pozycje = grupy.get(nazwa) ?? [];
    pozycje.push(proces);
    grupy.set(nazwa, pozycje);
  }
  return [...grupy.entries()].sort(([a], [b]) => a.localeCompare(b));
}

function nazwaGrupy(proces: TerminalProcess, klucz: string): string {
  if (klucz === 'karta') return `Karta ${proces.sessionId ?? '— bez karty'}`;
  if (klucz === 'inicjator') {
    return proces.initiator === ProcessInitiator.Model ? 'Uruchomione przez AI' : 'Uruchomione przez Operatora';
  }
  return `Stan: ${proces.status}`;
}
