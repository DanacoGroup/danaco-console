import type { NazwaIkony } from '../ikony/ikony';

/**
 * Ikona pozycji spisu okien pomocniczych.
 *
 * Jedno wiązanie kodu pozycji z ikoną obsługuje oba miejsca nagłówka rozmowy:
 * wykaz w menu `⋮` i rząd skrótów złożony z samych ikon. Rozdzielenie ich dałoby
 * tej samej pozycji dwa różne rysunki w jednym nagłówku.
 *
 * Ikonę dobiera się pod czynność, nie pod wygląd narzędzia — `podglad-bash`
 * bierze `monitor`, bo pozycja jest podglądem pracy idącej w tle, a nie drugim
 * terminalem.
 *
 * O tym, które pozycje wolno otworzyć, rozstrzyga
 * `okna-pomocnicze/wytwornia-paneli.ts`. Pozycja spoza wykazu nie jest błędem
 * i nie zostaje bez rysunku: dostaje `karta-okna`, bo każda z nich jest oknem
 * obok rozmowy.
 */

/** Wiązanie kodu pozycji z ikoną zestawu. */
const IKONY: Readonly<Record<string, NazwaIkony>> = {
  'podglad-bash': 'monitor',
  // Czynnością panelu jest sięgnięcie wstecz po zapis rozmowy, nie oglądanie
  // czegokolwiek na żywo — stąd „historia", a nie „monitor" ani „lista".
  'historia-rozmowy': 'historia',
  'terminal': 'terminal',
  'przebieg-debaty': 'debata',
  'zasoby-designu': 'paleta',
  'pliki': 'folder',
  'przegladarka': 'globus',
  'artefakty': 'archiwum',
  'pliki-srodowiska': 'plik',
};

/** Ikona pozycji; pozycja spoza wykazu dostaje znak okna, nie pustkę. */
export function ikonaPanelu(kod: string): NazwaIkony {
  return IKONY[kod] ?? 'karta-okna';
}
