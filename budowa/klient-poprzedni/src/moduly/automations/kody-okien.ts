/**
 * Kody okien operacyjnych modułu Automations zapisane w jednym miejscu. Kod
 * okna jest nagi, bez przedrostka modułu, ponieważ dokładnie w takim zapisie
 * zna go rejestr okien operacyjnych rdzenia.
 */

/**
 * Pięć okien operacyjnych, które moduł buduje. Kolejność pozycji odpowiada
 * kolejności w układzie modułu, a wartości są kodami znanymi rejestrowi rdzenia.
 */
export const KODY_OKIEN = {
  workflowBuilder: 'workflow-builder',
  orchestrator: 'orchestrator',
  scheduler: 'scheduler',
  queueManager: 'queue-manager',
  executionMonitor: 'execution-monitor',
} as const;

/**
 * Kod modułu w rejestrze rdzenia. Odczyt `module.list` wskazuje po nim okna
 * operacyjne przypisane modułowi, więc zapis musi być zgodny z rejestrem.
 */
export const KOD_MODULU = 'automations';

/**
 * Znacznik układu wykazu gotowych pętli, używany wyłącznie po stronie klienta.
 * Stoi osobno od kodów okien, ponieważ rejestr okien operacyjnych rdzenia tego
 * kodu nie zna, a wartość nie trafia do rdzenia w żadnym żądaniu.
 */
export const KOD_WYKAZU_PETLI = 'wykaz-petli';
