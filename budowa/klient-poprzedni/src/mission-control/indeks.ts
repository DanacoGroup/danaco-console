/**
 * Mission Control — pulpit operacyjny. Interfejs katalogu.
 *
 * Powłoka montuje pulpit i wpina źródło danych rdzenia jednym ruchem:
 *
 *   import {
 *     pustyKomplet,
 *     utworzMissionControl,
 *     utworzZrodloPulpitu,
 *   } from './mission-control/indeks';
 *
 *   const pulpit = utworzMissionControl(pustyKomplet());
 *   const zrodlo = utworzZrodloPulpitu(kanal, (dane) => pulpit.odswiez(dane));
 *   zrodlo.uruchom();
 *
 * Komplet danych pochodzi wyłącznie z odczytów i zdarzeń rdzenia:
 * `session.list`, `window.list`, `channel.list` oraz `session.changed`,
 * `window.changed`, `queue.changed`, `progress.changed`. Do pierwszego odczytu
 * pulpit pokazuje stany puste zamiast wartości zastępczych.
 *
 * Plik nie zawiera logiki — jest wykazem tego, co katalog wystawia na zewnątrz.
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
