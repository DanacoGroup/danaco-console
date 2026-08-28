import { QueueStatus, WindowRole } from '../../../shared/contract';
import type { NazwaIkony } from '../ikony/ikony';

/**
 * Stan pary koordynator-wykonawca ustala, w jakim położeniu jest pętla, oraz sprowadza stan kolejki z kontraktu do stanu widocznego na scenie, bez własnego katalogu nazw kolejki.
 */
export type StanPary =
  /** Brak dwóch okien w rolach koordynatora i wykonawcy. */
  | 'brak-pary'
  /** Para ustawiona, pętla jeszcze nie ruszyła. */
  | 'gotowa'
  /** Wykonawca pracuje nad zleceniem. */
  | 'wykonawca-pracuje'
  /** Wykonawca skończył turę i wybudził koordynatora. */
  | 'koordynator-wybudzony'
  /** Kolejka wstrzymana — decyzja należy do operatora. */
  | 'kolejka-wstrzymana'
  /** Chwila przekazania zlecenia z okna koordynatora do okna wykonawcy. */
  | 'przekazanie';

/** Wygląd stanu pary: ikona, wariant plakietki i kropka — nigdy sama barwa, żeby odczyt nie zależał od koloru. */
export interface WygladStanu {
  ikona: NazwaIkony;
  /** Klasa wariantu `.dn-plakietka--*`; pusta dla plakietki neutralnej. */
  wariantPlakietki: string;
  /** Klasa wariantu `.dn-kropka--*`; pusta dla kropki sygnałowej. */
  wariantKropki: string;
  /** Czy stan oznacza pracę trwającą teraz; prawda wprowadza wskaźnik ładowania obok etykiety. */
  wTrakcie: boolean;
}

/**
 * Wygląd stanu pary.
 *
 * Każdy stan niesie ikonę i etykietę, więc odczyt nie zależy od rozróżnienia
 * barw. Akcent sygnałowy występuje wyłącznie na plakietce wielkości pigułki,
 * nigdy na powierzchni.
 */
export function wygladStanu(stan: StanPary): WygladStanu {
  switch (stan) {
    case 'wykonawca-pracuje':
      // Tętno jest zastrzeżone dla pracy trwającej w tle (`plakietka.css`).
      return znak('uruchom', 'dn-plakietka--informacja', 'dn-kropka--tetno', true);
    case 'koordynator-wybudzony':
      return znak('dzwonek', 'dn-plakietka--sygnal', '', false);
    case 'kolejka-wstrzymana':
      return znak('zatrzymaj', 'dn-plakietka--ostrzezenie', 'dn-kropka--ostrzezenie', false);
    case 'przekazanie':
      return znak('wyslij', 'dn-plakietka--sukces', 'dn-kropka--sukces', false);
    case 'gotowa':
      return znak('zegar', '', 'dn-kropka--neutralna', false);
    case 'brak-pary':
      return znak('info', '', 'dn-kropka--neutralna', false);
  }
}

/** Skrót zapisu wyglądu stanu pary — cztery pola złożone w jednym wierszu gałęzi, bez powtarzania nazw pól. */
function znak(
  ikona: NazwaIkony,
  wariantPlakietki: string,
  wariantKropki: string,
  wTrakcie: boolean,
): WygladStanu {
  return { ikona, wariantPlakietki, wariantKropki, wTrakcie };
}

/**
 * Stan pary wyprowadzony ze stanu kolejki rdzenia.
 *
 * Kolejka wyczerpana oznacza, że wykonawca zamknął turę — a zamknięcie tury
 * wybudza koordynatora. Wstrzymanie i zatrzymanie prowadzą do tego
 * samego obrazu na scenie: pętla stoi i czeka na człowieka.
 */
export function stanZKolejki(status: QueueStatus): StanPary {
  switch (status) {
    case QueueStatus.Running:
      return 'wykonawca-pracuje';
    case QueueStatus.Paused:
    case QueueStatus.Stopped:
      return 'kolejka-wstrzymana';
    case QueueStatus.Done:
      return 'koordynator-wybudzony';
    case QueueStatus.Idle:
      return 'gotowa';
  }
}

/**
 * Stan gniazda w oderwaniu od pary — dla okna samodzielnego i dla okna,
 * które w pętli nie uczestniczy.
 */
export function czyRolaWPetli(rola: WindowRole): boolean {
  return rola === WindowRole.Coordinator || rola === WindowRole.Executor;
}
