/**
 * Stała zawiera kody siedmiu okien operacyjnych modułu Research, zgodne z wierszami zakładanymi przez migrację migracja_031_okna_modulow.sql.
 */
export const KODY_OKIEN = {
  workspace: 'research-workspace',
  odkrywanie: 'discovery-panel',
  zrodla: 'sources-manager',
  lektura: 'reading-view',
  ustalenia: 'findings-panel',
  raport: 'report-builder',
  eksport: 'export-panel',
} as const;

/** Stała zawiera kod modułu Research w katalogu rdzenia, odpowiadający kolumnie kod tabeli modułów bazy danych. */
export const KOD_MODULU = 'research';

/**
 * Stała wylicza kody okien modułu Research nieobecne jeszcze w katalogu rdzenia, dla których rdzeń odmawia dostępu przed dobudowaniem migracji.
 */
export const OKNA_SPOZA_KATALOGU: readonly string[] = [
  KODY_OKIEN.odkrywanie,
  KODY_OKIEN.lektura,
];
