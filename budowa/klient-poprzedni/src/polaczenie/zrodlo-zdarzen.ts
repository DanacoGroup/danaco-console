import type { Envelope, EventPayloadOf, EventType } from '../../../shared/contract';
import type { Odsubskrybuj } from './magistrala-zdarzen';

/**
 * Źródło komunikatów przychodzących, na którym osadzają się obserwatory. Obserwatory żyją w warstwie łączności, a kanał kontraktu w warstwie protokołu; opisując dokładnie to, czego potrzebują, obserwatory unikają importu kanału wprost i cyklu warstw.
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
