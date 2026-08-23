/**
 * Kody siedmiu okien operacyjnych modułu Research.
 *
 * Jedna odpowiedzialność: związanie widoku z katalogiem rdzenia. Kody pochodzą
 * z macierzy `okno_operacyjne_modul` (wiersze modułu `research`, zakładane
 * w `server/internal/store/migracja_031_okna_modulow.sql`). Trzymane w jednym
 * miejscu sprawiają, że dopisanie okna po stronie rdzenia rozjeżdża się
 * z klientem widocznie, a nie w siedmiu plikach naraz.
 *
 * Dwa kody — `discovery-panel` i `reading-view` — wiersza w katalogu rdzenia
 * jeszcze NIE mają: migracje 030 i 031 zakładają dla modułu Research pięć okien,
 * a opracowanie modułu (`docs/moduly/research.md`, rozdz. 2) wylicza siedem
 * własnych obok dwóch wspólnych platformie. Kody stoją tu w brzmieniu, w jakim
 * mają wejść do katalogu — dzięki temu `data-okno` obu okien jest już dziś tym
 * samym napisem, którym będzie po dobudowie rdzenia, i nawigacja
 * wewnątrzmodułowa nie wymaga zmiany.
 */

/** Kod okna operacyjnego w katalogu rdzenia (`okno_operacyjne.kod`). */
export const KODY_OKIEN = {
  workspace: 'research-workspace',
  odkrywanie: 'discovery-panel',
  zrodla: 'sources-manager',
  lektura: 'reading-view',
  ustalenia: 'findings-panel',
  raport: 'report-builder',
  eksport: 'export-panel',
} as const;

/** Kod modułu w katalogu rdzenia — kolumna `modul.kod`. */
export const KOD_MODULU = 'research';

/**
 * Okna modułu, których katalog rdzenia jeszcze nie zna.
 *
 * Wykaz jest zaporą, nie zgodą: okna są zbudowane i czynne, a ich kody wchodzą
 * do `window.action` tak samo jak kody okien znanych. Rdzeń odmówi im na
 * katalogu okien, zanim dojdzie do katalogu akcji, i to okno nazywa Operatorowi
 * wprost, zamiast pokazywać puste miejsce. Wykaz znika, gdy migracje 030 i 031
 * dostaną wiersze obu okien.
 */
export const OKNA_SPOZA_KATALOGU: readonly string[] = [
  KODY_OKIEN.odkrywanie,
  KODY_OKIEN.lektura,
];
