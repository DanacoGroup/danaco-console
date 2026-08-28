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
 * `zrodlo-pulpitu.ts` zapisuje tu odpowiedzi i zdarzenia, `zlozenie-danych.ts`
 * czyta je i przekłada na `DanePulpitu`. Stan początkowy jest pusty.
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
  /** Środowiska platformy z odczytu `environment.list`, kluczowane kodem. */
  srodowiska: Map<string, Environment>;
}

/**
 * Stan początkowy — nic nie odczytano, wszystkie wykazy puste. Złożenie sięga po
 * tę wartość przy budowie źródła, a odpowiedzi i zdarzenia rdzenia dopisują
 * pozycje do zwróconych tu map.
 */
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
