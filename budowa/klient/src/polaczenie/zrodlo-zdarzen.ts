import type { Envelope, EventPayloadOf, EventType } from '../../../shared/contract.ts';
import type { Odsubskrybuj } from './magistrala-zdarzen.ts';

/**
 * Źródło komunikatów przychodzących, na którym osadzają się obserwatory.
 * Obserwatory żyją w warstwie połączenia, a kanał kontraktu — w warstwie
 * protokołu. Import kanału odwróciłby strzałkę zależności i stworzyłby
 * cykl warstw.
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
