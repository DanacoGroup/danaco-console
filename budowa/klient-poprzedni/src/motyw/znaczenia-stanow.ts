import { ProgressStatus } from '../../../shared/contract';
import type { NazwaIkony } from '../ikony/zrodla-ikon';

/**
 * Znaczenia stanów — jedno miejsce prawdy o tym, czym stan jest widziany.
 *
 * Stan nigdy nie jest sygnalizowany samym kolorem: katalog wiąże każdy stan
 * z ikoną i etykietą, więc odczyt nie zależy od barwy. Jeden zapis znaczenia
 * dla wszystkich widoków zastępuje rozsypane mapy stanu.
 *
 * Barw tu nie ma — mieszkają w `stany.css`, a tu stoi wyłącznie nazwa rodziny,
 * czyli wskazanie, którego wariantu biblioteki użyć. Cztery rodziny (`sukces`,
 * `ostrzezenie`, `blad`, `informacja`) plus brak rodziny (`neutralna`)
 * wystarczają, bo rozróżnienie niesie ikona i etykieta, nie kolejna barwa.
 *
 * Granica wiedzy: plik mówi, jak stan ma być pokazany. Nie mówi, jaki stan
 * jest — to rozstrzyga `shared/contract`.
 */

/** Rodzina barw stanu — nazwy wariantów `stany.css`; `neutralna` to brak rodziny. */
export type RodzinaStanu = 'sukces' | 'ostrzezenie' | 'blad' | 'informacja' | 'neutralna';

/**
 * Stan znaczący, czyli taki, który użytkownik ma odróżnić od innego.
 *
 * `do-weryfikacji` i `przyjety` nie mają odpowiednika w kontrakcie — patrz
 * `STANY_BEZ_ODPOWIEDNIKA_W_POSTEPIE`. `do-weryfikacji` siedzi na rodzinie
 * `ostrzezenie`, czyli dokładnie tej samej barwie co `wstrzymany`: dwa różne
 * stany, jedna barwa, więc bez ikony i bez etykiety byłyby na ekranie nie do
 * rozróżnienia.
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

/** Znak stanu: rodzina barw, ikona, etykieta i informacja o pracy trwającej. */
export interface ZnakStanu {
  /** Rodzina barw — sama nigdy nie wystarcza za odczyt stanu. */
  rodzina: RodzinaStanu;
  /** Ikona z biblioteki — pierwszy znak niebędący kolorem. */
  ikona: NazwaIkony;
  /** Etykieta po polsku — drugi znak niebędący kolorem. */
  etykieta: string;
  /**
   * Czy stan oznacza pracę trwającą teraz.
   *
   * Prawda wprowadza `.dn-spinner` obok etykiety — nigdy zamiast niej — oraz
   * tętno kropki.
   */
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
  // Dzieli rodzinę z „zakończonym" — „przyjęty" to zatwierdzenie wyniku,
  // „zakończony" to sam koniec pracy. Dlatego ptaszek w kole, nie sam ptaszek.
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
 * Stany katalogu, które nie mają odpowiednika w `ProgressStatus` — wymienione
 * wprost, nie przemilczane.
 *
 * Trzy pozycje, dwa różne powody:
 *
 *   `anulowany` — kontrakt go zna, tyle że pod innym wyliczeniem
 *   (`AssistantActionStatus.Cancelled`). Katalog motywu nie mapuje tego
 *   wyliczenia, bo moduł Assistant prowadzi własny wykaz plakietek
 *   (`moduly/assistant/etykiety-assistant.ts`) i oba wykazy się różnią:
 *   „w toku" idzie tam rodziną `sygnal`, a w parze okien rodziną `informacja`.
 *   Dołożenie tu drugiego zapisu bez usunięcia tamtego dałoby trzecią prawdę
 *   zamiast jednej.
 *
 *   `do-weryfikacji`, `przyjety` — kontrakt ich nie zna. Pozycja kolejki ma
 *   wyłącznie `QueueStatus` (idle, running, paused, stopped, done),
 *   postęp — `ProgressStatus`.
 *
 * Stan wniesiony do kontraktu wypada stąd i wchodzi do mapy.
 */
export const STANY_BEZ_ODPOWIEDNIKA_W_POSTEPIE: readonly ZnaczenieStanu[] = [
  'do-weryfikacji',
  'przyjety',
  'anulowany',
];

/** Znak stanu procesu oddanego przez rdzeń — jedno wejście dla widoków. */
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
 * Klasa wariantu kropki dla znaku stanu; pusty łańcuch dla kropki sygnałowej.
 *
 * Praca trwająca bierze tętno, nie barwę rodziny — tętno jest jedynym ruchem
 * ciągłym interfejsu i jest zastrzeżone właśnie dla pracy w tle
 * (`komponenty/plakietka.css`). Przy `prefers-reduced-motion` tętno zamiera,
 * a jego znaczenie przejmuje pierścień statyczny — również w arkuszu.
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
      // Kropka bazowa jest kropką sygnału (`--dn-kropka`), a informacja to
      // rodzina sygnału — osobnego wariantu nie ma.
      return '';
  }
}
