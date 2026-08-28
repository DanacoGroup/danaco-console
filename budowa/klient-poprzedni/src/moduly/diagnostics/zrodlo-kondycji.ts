import {
  Command,
  EventType,
  type AlertRuleListResponse,
  type AlertRuleRemoveResponse,
  type AlertRuleSaveRequest,
  type AlertRuleSaveResponse,
  type AlertTriggerAcknowledgeRequest,
  type AlertTriggerAcknowledgeResponse,
  type AlertTriggeredEvent,
  type AlertTriggerListResponse,
  type HealthProbeListResponse,
  type HealthProbeRemoveResponse,
  type HealthProbeRunResponse,
  type HealthProbeSaveRequest,
  type HealthProbeSaveResponse,
  type HealthResultListResponse,
  type HealthUptimeGetResponse,
  type ProvenanceCallReplayResponse,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyLogiczna, czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Interfejs udostępnia dwanaście komend trzech rodzin opisujących stan produktu:
 * sondy kondycji, reguły wyzwalania z rejestrem wyzwoleń oraz powtórzenie
 * wywołania modelu, kosztowne jak nowe wywołanie kanału.
 */
export interface ZrodloKondycji {
  /** `health.probe.save` — definicja sondy; brak identyfikatora zakłada nową. */
  zapiszSonde(zadanie: HealthProbeSaveRequest): Promise<Wynik<HealthProbeSaveResponse>>;
  /** `health.probe.list` — definicje wraz z wynikiem ostatniego przebiegu. */
  sondy(): Promise<Wynik<HealthProbeListResponse>>;
  /** `health.probe.remove` — usunięcie sondy wraz z serią jej wyników. */
  usunSonde(idSondy: string): Promise<Wynik<HealthProbeRemoveResponse>>;
  /** `health.probe.run` — pomiar na żądanie, poza odstępem sondy. */
  wykonajSonde(idSondy: string): Promise<Wynik<HealthProbeRunResponse>>;
  /** `health.result.list` — seria pomiarów, z której liczy się dostępność. */
  wyniki(idSondy: string): Promise<Wynik<HealthResultListResponse>>;
  /** `health.uptime.get` — dostępność i budżet błędów wyliczone z serii. */
  dostepnosc(idSondy: string): Promise<Wynik<HealthUptimeGetResponse>>;
  /** `alert.rule.save` — reguła wyzwalania; brak identyfikatora zakłada nową. */
  zapiszRegule(zadanie: AlertRuleSaveRequest): Promise<Wynik<AlertRuleSaveResponse>>;
  /** `alert.rule.list` — reguły wraz z licznikiem ich wyzwoleń. */
  reguly(): Promise<Wynik<AlertRuleListResponse>>;
  /** `alert.rule.remove` — usunięcie reguły wraz z jej rejestrem wyzwoleń. */
  usunRegule(idReguly: string): Promise<Wynik<AlertRuleRemoveResponse>>;
  /** `alert.trigger.list` — rejestr wyzwoleń przeliczony w chwili odpowiedzi. */
  wyzwolenia(idReguly: string): Promise<Wynik<AlertTriggerListResponse>>;
  /** `alert.trigger.acknowledge` — potwierdzenie obsłużenia wyzwolenia. */
  potwierdzWyzwolenie(
    idWyzwolenia: string,
    notatka: string,
  ): Promise<Wynik<AlertTriggerAcknowledgeResponse>>;
  /** `provenance.call.replay` — powtórzenie wywołania modelu. */
  powtorzWywolanie(idWywolania: string): Promise<Wynik<ProvenanceCallReplayResponse>>;
  /** Subskrypcja `alert.triggered` — droga alertu do okna. */
  naWyzwolenie(sluchacz: (tresc: AlertTriggeredEvent) => void): Odsubskrybuj;
}

export function utworzZrodloKondycji(kanal: Kanal): ZrodloKondycji {
  return {
    async zapiszSonde(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.HealthProbeSave, zadanie),
        Command.HealthProbeSave,
        (tresc) => czyObiekt(tresc.probe) && czyLogiczna(tresc.created),
      );
    },

    async sondy() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.HealthProbeList, {}),
        Command.HealthProbeList,
        (tresc) => czyTablica(tresc.probes),
      );
    },

    async usunSonde(idSondy) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.HealthProbeRemove, { probeId: idSondy }),
        Command.HealthProbeRemove,
        (tresc) => typeof tresc.probeId === 'string',
      );
    },

    async wykonajSonde(idSondy) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.HealthProbeRun, { probeId: idSondy }),
        Command.HealthProbeRun,
        (tresc) => czyObiekt(tresc.result) && czyObiekt(tresc.probe),
      );
    },

    async wyniki(idSondy) {
      const zadanie: Record<string, unknown> = {};
      if (idSondy !== '') zadanie['probeId'] = idSondy;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.HealthResultList, zadanie),
        Command.HealthResultList,
        (tresc) => czyTablica(tresc.results),
      );
    },

    async dostepnosc(idSondy) {
      const zadanie: Record<string, unknown> = {};
      if (idSondy !== '') zadanie['probeId'] = idSondy;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.HealthUptimeGet, zadanie),
        Command.HealthUptimeGet,
        (tresc) => czyTablica(tresc.uptimes),
      );
    },

    async zapiszRegule(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AlertRuleSave, zadanie),
        Command.AlertRuleSave,
        (tresc) => czyObiekt(tresc.rule) && czyLogiczna(tresc.created),
      );
    },

    async reguly() {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AlertRuleList, {}),
        Command.AlertRuleList,
        (tresc) => czyTablica(tresc.rules),
      );
    },

    async usunRegule(idReguly) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AlertRuleRemove, { ruleId: idReguly }),
        Command.AlertRuleRemove,
        (tresc) => typeof tresc.ruleId === 'string',
      );
    },

    async wyzwolenia(idReguly) {
      const zadanie: Record<string, unknown> = {};
      if (idReguly !== '') zadanie['ruleId'] = idReguly;
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AlertTriggerList, zadanie),
        Command.AlertTriggerList,
        (tresc) => czyTablica(tresc.triggers),
      );
    },

    async potwierdzWyzwolenie(idWyzwolenia, notatka) {
      const zadanie: AlertTriggerAcknowledgeRequest = { triggerId: idWyzwolenia };
      if (notatka.trim() !== '') zadanie.note = notatka.trim();
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AlertTriggerAcknowledge, zadanie),
        Command.AlertTriggerAcknowledge,
        (tresc) => czyObiekt(tresc.trigger),
      );
    },

    async powtorzWywolanie(idWywolania) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ProvenanceCallReplay, { callId: idWywolania }),
        Command.ProvenanceCallReplay,
        (tresc) => czyObiekt(tresc.replay),
      );
    },

    naWyzwolenie(sluchacz) {
      return kanal.naZdarzenie(EventType.AlertTriggered, (tresc) => sluchacz(tresc));
    },
  };
}
