import type {
  Channel,
  Environment,
  ProgressChangedEvent,
  Queue,
  Session,
  SessionPresence,
  Window,
} from '../../../shared/contract';

/**
 * Stan surowy pulpitu — dokładnie to, co rdzeń oddał odczytami i zdarzeniami.
 *
 * Plik ustala kształt i wartość początkową wsadu złożenia. `zrodlo-pulpitu.ts`
 * zapisuje tu odpowiedzi i zdarzenia, `zlozenie-danych.ts` czyta i przekłada na
 * `DanePulpitu`. Stan początkowy jest pusty, a `odczytano` mówi, czy rdzeń już
 * się odezwał.
 */
export interface StanZrodla {
  /** Czy odczyt z rdzenia już nadszedł (odpowiedź albo zdarzenie). */
  odczytano: boolean;
  sesje: Map<string, Session>;
  obecnosc: Map<string, SessionPresence>;
  okna: Map<string, Window>;
  kanaly: Channel[];
  kolejki: Map<string, Queue>;
  procesy: Map<string, ProgressChangedEvent>;
  /**
   * Środowiska platformy z odczytu `environment.list`, kluczowane kodem.
   * Nazwa kolumny matrycy pochodzi stąd, nie z kopii katalogu w pulpicie.
   */
  srodowiska: Map<string, Environment>;
}

/** Stan początkowy — nic nie odczytano, wszystkie wykazy puste. */
export function pustyStanZrodla(): StanZrodla {
  return {
    odczytano: false,
    sesje: new Map(),
    obecnosc: new Map(),
    okna: new Map(),
    kanaly: [],
    kolejki: new Map(),
    procesy: new Map(),
    srodowiska: new Map(),
  };
}
