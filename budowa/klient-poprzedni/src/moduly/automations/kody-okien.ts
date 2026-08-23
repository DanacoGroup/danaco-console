/**
 * Kody okien operacyjnych modułu Automations — jedno miejsce, jeden zapis.
 *
 * Kod okna jest nagi, bez przedrostka modułu: rejestr rdzenia zna
 * `workflow-builder`, `scheduler`, `queue-manager`, `orchestrator`
 * i `execution-monitor` (`migracja_030_rejestr_okien_operacyjnych.sql`),
 * a powiązanie tych kodów z modułem `automations` niesie
 * `migracja_031_okna_modulow.sql`.
 *
 * Kody stoją w jednym miejscu, bo rdzeń nie sprawdza identyfikatora okna
 * podawanego w `automation.execution.subscribe` — zapamiętuje go jako klucz
 * obserwacji. Kod z przedrostkiem nie wywołałby więc odmowy, tylko trafiłby
 * do `Queue.windowIds` i do telemetrii postępu jako okno nieistniejące.
 */

/** Pięć okien, które moduł buduje — kolejność jak w układzie modułu. */
export const KODY_OKIEN = {
  workflowBuilder: 'workflow-builder',
  orchestrator: 'orchestrator',
  scheduler: 'scheduler',
  queueManager: 'queue-manager',
  executionMonitor: 'execution-monitor',
} as const;

/** Kod modułu w rejestrze rdzenia — po nim `module.list` wskazuje jego okna. */
export const KOD_MODULU = 'automations';

/**
 * Znacznik układu wykazu gotowych pętli — wyłącznie po stronie klienta.
 *
 * Stoi osobno, a nie w `KODY_OKIEN`, bo rejestr okien operacyjnych rdzenia tego
 * kodu nie zna. Wartość służy jedynie znacznikowi `data-okno` w układzie modułu
 * (skok nawigacji, sprawdziany widoku) i nie trafia do rdzenia w żadnym żądaniu:
 * okno wykazu nie zakłada obserwacji telemetrii i nie podaje `windowIds` przy
 * zakładaniu kolejki.
 */
export const KOD_WYKAZU_PETLI = 'wykaz-petli';
