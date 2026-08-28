/**
 * Definicje siedmiu kart okna roboczego Studia — jedno źródło dla pasma kart,
 * dla zaczepu paneli i dla stanu początkowego przełączania. Karta `editor`
 * jest jedyną z rzeczywistym przekrojem do rdzenia (`dokument.ts`); pozostałe
 * sześć wchodzi osobnym zakresem prac i tu niosą wyłącznie nazwę i znak.
 */

import type { NazwaZnaku } from '../ikony.ts';
import { tresci } from '../tresci.ts';

export interface DefinicjaKarty {
  /** Kod karty — rdzeń identyfikatorów `karta-${kod}` / `panel-${kod}`. */
  kod: string;
  nazwa: string;
  ikona: NazwaZnaku;
}

export const KARTA_EDITOR = 'editor';

export const karty: DefinicjaKarty[] = [
  { kod: KARTA_EDITOR, nazwa: tresci.karty.editor, ikona: 'olowek' },
  { kod: 'tools', nazwa: tresci.karty.tools, ikona: 'klucz' },
  { kod: 'diff', nazwa: tresci.karty.diff, ikona: 'diff' },
  { kod: 'repo', nazwa: tresci.karty.repo, ikona: 'historia' },
  { kod: 'preview', nazwa: tresci.karty.preview, ikona: 'oko' },
  { kod: 'pliki', nazwa: tresci.karty.pliki, ikona: 'plik' },
  { kod: 'plan', nazwa: tresci.karty.plan, ikona: 'plan' },
];
