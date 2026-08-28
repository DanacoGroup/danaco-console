import type { Kanal } from '../../protokol/kanal';
import type { OpisModulu } from '../rejestracja';
import { utworzModulAgents as wytworzAgents } from './modul-agents';
/**
 * Punkt zbiorczy katalogu modułu Agents: wystawia wytwórnię widoku oraz typy
 * stanu i źródeł danych. Moduł nie osadza się sam w dokumencie i nie zna
 * powłoki — oddaje element, a warstwa składająca rozstrzyga, gdzie go
 * postawić.
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
 * Testowany ekspert, czyli kontekst roboczy czatu testowego modułu.
 * Wystawiony w interfejsie katalogu, ponieważ sięga po niego powłoka: okno
 * rozmowy stoi na scenie sesji obok widoku modułu i to ono musi wiedzieć,
 * kiedy testowany agent się zmienił.
 */
export { testowanyAgent, type RozglosTestowanego } from './testowany-agent';

/**
 * Samoopisujący się moduł dla rejestru powłoki. Agents ma kształt widoku,
 * ponieważ oddaje element oraz czynność wczytania, więc warstwy przejściowej
 * nie potrzebuje.
 */
export const MODUL: OpisModulu = {
  kod: 'agents',
  utworzWidok: (kanal: Kanal) => wytworzAgents(kanal),
};
