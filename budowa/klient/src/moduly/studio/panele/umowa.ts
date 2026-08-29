/**
 * Umowa panelu okna Studia. Każdy z siedmiu paneli mieszka w osobnym pliku tego
 * katalogu i wystawia jedną funkcję montującą — dzięki temu panele powstają
 * niezależnie od siebie i nie dzielą żadnego pliku poza tą umową.
 *
 * Panel dostaje kanał do rdzenia i identyfikator okna modułu; resztę stanu
 * czyta sam swoimi komendami. Zwraca uchwyt zdejmujący, żeby przełączenie karty
 * mogło odłączyć nasłuchy zamiast zostawiać je w tle.
 */

import type { Kanal } from '../../../protokol/kanal.ts';

export interface ZaleznosciPanelu {
  /** Kanał, którym panel woła komendy rdzenia. */
  kanal: Kanal;
  /** Okno modułu Studia założone przez powłokę; puste, gdy rdzeń go nie dał. */
  idOkna: string | null;
  /** Sesja, w której panel pracuje; pusta, gdy rdzeń jej nie założył. */
  idSesji: string | null;
}

export interface ZamontowanyPanel {
  /** Odłącza nasłuchy i czyści węzeł. */
  zdejmij(): void;
}

/** Kształt, który wystawia każdy plik panelu. */
export type MontazPanelu = (
  wezel: HTMLElement,
  zaleznosci: ZaleznosciPanelu,
) => ZamontowanyPanel;
