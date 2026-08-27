/**
 * Katalog pulpitu operacyjnego Mission Control: wykaz tego, co warstwa wystawia
 * na zewnątrz, bez własnej logiki. Komplet danych pochodzi wyłącznie z odczytów
 * i zdarzeń rdzenia, a do pierwszego odczytu pulpit pokazuje stany puste.
 */
export { utworzMissionControl, type MissionControl } from './mission-control';

export { utworzZrodloPulpitu, type ZrodloPulpitu } from './zrodlo-pulpitu';

export { pustyKomplet } from './zlozenie-danych';

export {
  RodzajRelacji,
  RodzajUtworzenia,
  ZrodloDanych,
  type AgentZespolu,
  type AktywnoscAI,
  type DanePulpitu,
  type IdSrodowiska,
  type KanalOperacyjny,
  type KolejkaPulpitu,
  type KolumnaSrodowiska,
  type PasDecyzji,
  type ProcesWTle,
  type RelacjaSesji,
  type SesjaMatrycy,
} from './model-danych';

export {
  DZIALANIA_TRANSPORTU,
  type WejscieDoSesji,
  type ZamiarDecyzji,
  type ZamiarKolejki,
  type ZamiarUtworzenia,
} from './zdarzenia-pulpitu';

export type { Odpiecie } from './nadajnik';
