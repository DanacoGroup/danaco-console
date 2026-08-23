import type { Envelope, EventPayloadOf, EventType } from '../../../shared/contract';
import type { Odsubskrybuj } from './magistrala-zdarzen';

/**
 * Źródło komunikatów przychodzących, na którym osadzają się obserwatory.
 *
 * Obserwatory żyją w warstwie łączności, a kanał kontraktu — w warstwie
 * protokołu. Gdyby obserwator importował `Kanal`, strzałka zależności
 * odwróciłaby się i powstałby cykl warstw. Zamiast tego obserwator opisuje
 * dokładnie to, czego potrzebuje: subskrypcję zdarzenia po nazwie z kontraktu
 * i podgląd całego ruchu. `Kanal` spełnia ten opis kształtem, bez ani jednej
 * dodatkowej deklaracji (jeden byt = jeden moduł).
 */
export interface ZrodloZdarzen {
  /** Subskrypcja zdarzeń jednego typu wraz z ich treścią. */
  naZdarzenie<K extends EventType>(
    zdarzenie: K,
    sluchacz: (tresc: EventPayloadOf<K>, koperta: Envelope) => void,
  ): Odsubskrybuj;
  /** Subskrypcja całego ruchu przychodzącego. */
  naDowolny(sluchacz: (koperta: Envelope) => void): Odsubskrybuj;
}
