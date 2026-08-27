import type { NazwaIkony } from '../ikony/ikony';

/**
 * Ikona pozycji spisu okien pomocniczych łączy kod pozycji z ikoną zestawu obsługującą oba miejsca nagłówka rozmowy — wykaz menu i rząd skrótów — dobraną pod czynność, nie pod wygląd narzędzia.
 */
const IKONY: Readonly<Record<string, NazwaIkony>> = {
  'podglad-bash': 'monitor',
  // Czynnością panelu jest sięgnięcie wstecz po zapis rozmowy, nie oglądanie czegokolwiek na żywo.
  'historia-rozmowy': 'historia',
  'terminal': 'terminal',
  'przebieg-debaty': 'debata',
  'zasoby-designu': 'paleta',
  'pliki': 'folder',
  'przegladarka': 'globus',
  'artefakty': 'archiwum',
  'pliki-srodowiska': 'plik',
};

/** Ikona pozycji spisu okien pomocniczych; pozycja spoza wykazu dostaje znak okna, nie pustkę wizualną. */
export function ikonaPanelu(kod: string): NazwaIkony {
  return IKONY[kod] ?? 'karta-okna';
}
