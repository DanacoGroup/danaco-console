import { ProgressStatus } from '../../../shared/contract';
import type { NazwaIkony } from '../ikony/zrodla-ikon';

/**
 * Katalog znaczeń stanów wiąże każdy stan z rodziną barw, ikoną i etykietą, tak by odczyt nigdy nie zależał od samego koloru, a nazwa rodziny barw wskazuje jedynie, którego wariantu biblioteki stylów użyć.
 */
export type RodzinaStanu = 'sukces' | 'ostrzezenie' | 'blad' | 'informacja' | 'neutralna';

/**
 * Stan znaczący to taki, który użytkownik ma odróżnić od innego; dwa stany dzielące tę samą rodzinę barw różnią się zawsze ikoną i etykietą, nigdy samą barwą.
 */
export type ZnaczenieStanu =
  | 'oczekuje'
  | 'biegnie'
  | 'wstrzymany'
  | 'zatrzymany'
  | 'do-weryfikacji'
  | 'przyjety'
  | 'zakonczony'
  | 'bledny'
  | 'anulowany';

/** Znak stanu widziany przez interfejs: rodzina barw, ikona z biblioteki, etykieta po polsku oraz informacja o pracy trwającej teraz. */
export interface ZnakStanu {
  /** Rodzina barw — sama nigdy nie wystarcza za odczyt stanu. */
  rodzina: RodzinaStanu;
  /** Ikona z biblioteki — pierwszy znak niebędący kolorem. */
  ikona: NazwaIkony;
  /** Etykieta po polsku — drugi znak niebędący kolorem. */
  etykieta: string;
  /** Czy stan oznacza pracę trwającą teraz; prawda dokłada wskaźnik obok etykiety oraz tętno kropki. */
  wTrakcie: boolean;
}

/**
 * Katalog znaczeń.
 *
 * Dwie reguły, których wykaz pilnuje:
 *   1. każdy stan ma ikonę i etykietę, więc odczyt nigdy nie zależy od barwy;
 *   2. dwa stany dzielące rodzinę barw mają różne ikony i różne etykiety —
 *      inaczej barwa byłaby jedynym rozróżnieniem.
 */
export const ZNACZENIA_STANOW: Readonly<Record<ZnaczenieStanu, ZnakStanu>> = {
  oczekuje: { rodzina: 'neutralna', ikona: 'zegar', etykieta: 'Oczekuje', wTrakcie: false },
  // Rodzina i ikona zgodne z parą okien (`okna-rownolegle/stan-pary.ts`).
  biegnie: { rodzina: 'informacja', ikona: 'uruchom', etykieta: 'W toku', wTrakcie: true },
  wstrzymany: { rodzina: 'ostrzezenie', ikona: 'wstrzymaj', etykieta: 'Wstrzymany', wTrakcie: false },
  zatrzymany: { rodzina: 'neutralna', ikona: 'zatrzymaj', etykieta: 'Zatrzymany', wTrakcie: false },
  // Dzieli rodzinę ze „wstrzymanym" — dlatego ikona musi być inna.
  'do-weryfikacji': { rodzina: 'ostrzezenie', ikona: 'walidator', etykieta: 'Do weryfikacji', wTrakcie: false },
  // Dzieli rodzinę z zakończonym, lecz oznacza zatwierdzenie wyniku — stąd ptaszek w kole.
  przyjety: { rodzina: 'sukces', ikona: 'ptaszek-kolo', etykieta: 'Przyjęty', wTrakcie: false },
  zakonczony: { rodzina: 'sukces', ikona: 'ptaszek', etykieta: 'Zakończony', wTrakcie: false },
  bledny: { rodzina: 'blad', ikona: 'blad', etykieta: 'Błąd', wTrakcie: false },
  anulowany: { rodzina: 'neutralna', ikona: 'zamknij', etykieta: 'Anulowany', wTrakcie: false },
};

/**
 * Stan procesu z kontraktu → znaczenie.
 *
 * Wykaz jest pełny z zamysłu: dopisanie stanu w `shared/contract` przerywa
 * kompilację tutaj, zamiast wypuścić na ekran stan bez ikony i bez etykiety.
 */
export const ZNACZENIE_Z_POSTEPU: Readonly<Record<ProgressStatus, ZnaczenieStanu>> = {
  [ProgressStatus.Pending]: 'oczekuje',
  [ProgressStatus.Running]: 'biegnie',
  [ProgressStatus.Paused]: 'wstrzymany',
  [ProgressStatus.Stopped]: 'zatrzymany',
  [ProgressStatus.Done]: 'zakonczony',
  [ProgressStatus.Failed]: 'bledny',
};

/**
 * Stany katalogu bez odpowiednika w postępie kontraktu są wymienione wprost, nie przemilczane, bo stan wniesiony później do kontraktu wypada stąd i wchodzi do mapy.
 */
export const STANY_BEZ_ODPOWIEDNIKA_W_POSTEPIE: readonly ZnaczenieStanu[] = [
  'do-weryfikacji',
  'przyjety',
  'anulowany',
];

/** Znak stanu procesu oddanego przez rdzeń, obliczony jednym wejściem wspólnym dla wszystkich widoków interfejsu. */
export function znakStanuPostepu(stan: ProgressStatus): ZnakStanu {
  return ZNACZENIA_STANOW[ZNACZENIE_Z_POSTEPU[stan]];
}

/**
 * Klasa wariantu plakietki dla rodziny; pusty łańcuch dla `neutralna`.
 *
 * Nazwy pełne, nie sklejane w czasie działania: sklejka `dn-plakietka--${x}`
 * byłaby niewidoczna dla wyszukiwania klas po drzewie.
 */
export function wariantPlakietki(rodzina: RodzinaStanu): string {
  switch (rodzina) {
    case 'sukces':
      return 'dn-plakietka--sukces';
    case 'ostrzezenie':
      return 'dn-plakietka--ostrzezenie';
    case 'blad':
      return 'dn-plakietka--blad';
    case 'informacja':
      return 'dn-plakietka--informacja';
    case 'neutralna':
      // Biblioteka nie ma wariantu neutralnego — plakietka bazowa jest nim sama.
      return '';
  }
}

/**
 * Klasa wariantu kropki dla znaku stanu — pusty łańcuch dla kropki sygnałowej, bo praca trwająca bierze tętno, nie barwę rodziny.
 */
export function wariantKropki(znak: ZnakStanu): string {
  if (znak.wTrakcie) return 'dn-kropka--tetno';
  switch (znak.rodzina) {
    case 'sukces':
      return 'dn-kropka--sukces';
    case 'ostrzezenie':
      return 'dn-kropka--ostrzezenie';
    case 'blad':
      return 'dn-kropka--blad';
    case 'neutralna':
      return 'dn-kropka--neutralna';
    case 'informacja':
      // Kropka bazowa jest kropką sygnału, a informacja to rodzina sygnału — osobnego wariantu nie ma.
      return '';
  }
}
