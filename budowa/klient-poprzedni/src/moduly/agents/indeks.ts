import type { Kanal } from '../../protokol/kanal';
import type { OpisModulu } from '../rejestracja';
import { utworzModulAgents as wytworzAgents } from './modul-agents';
/**
 * Moduł Agents — punkt zbiorczy katalogu.
 *
 * Powłoka zna stąd jedną czynność:
 *
 *   import { utworzModulAgents } from '../moduly/agents/indeks';
 *   const modul = utworzModulAgents(rdzen.kanal);
 *   obszarRoboczy.pokaz(modul.element);
 *   await modul.wczytaj(idSesji);
 *
 * Moduł nie osadza się sam w dokumencie i nie zna powłoki — oddaje element,
 * a warstwa składająca decyduje, gdzie go postawić. Dzięki temu te same okna
 * wchodzą i w obszar roboczy powłoki, i w podgląd sprawdzianu.
 */
export { utworzModulAgents } from './modul-agents';
export type { ModulAgents } from './modul-agents';
export { utworzStanAgentow } from './stan-agentow';
export type { StanAgentow } from './stan-agentow';
export { utworzZrodloAgentow } from './zrodlo-agentow';
export type { ZrodloAgentow } from './zrodlo-agentow';
export { utworzZrodloZaplecza } from './zrodlo-zaplecza';
export type { ZrodloZaplecza } from './zrodlo-zaplecza';
/**
 * Testowany ekspert — kontekst roboczy czatu testowego tego modułu.
 *
 * Wystawiony w interfejsie katalogu, bo sięga po niego powłoka, nie moduł:
 * okno rozmowy stoi na scenie sesji, obok widoku modułu, i to ono musi
 * wiedzieć, kiedy testowany agent się zmienił (`aplikacja/ulotnosc-okna.ts`).
 */
export { testowanyAgent, type RozglosTestowanego } from './testowany-agent';

/**
 * Samoopisujący się moduł dla rejestru powłoki. Agents ma kształt widoku od
 * początku — oddaje element i `wczytaj`, więc przejścia nie potrzebuje.
 */
export const MODUL: OpisModulu = {
  kod: 'agents',
  utworzWidok: (kanal: Kanal) => wytworzAgents(kanal),
};
