import {
  Command,
  type AutomationExecutionSubscribeRequest,
  type AutomationExecutionSubscribeResponse,
  type AutomationScheduleSetRequest,
  type AutomationScheduleSetResponse,
  type AutomationWorkflowListRequest,
  type AutomationWorkflowListResponse,
  type AutomationWorkflowSaveRequest,
  type AutomationWorkflowSaveResponse,
  type ScheduleGetRequest,
  type ScheduleGetResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Komendy, którymi Automation Studio prowadzi scenariusze przeglądania.
 *
 * Obszar `browser.*` nie niesie ani jednej komendy makra: kontrakt trzyma
 * powtarzalne przebiegi w jednym miejscu — jako automatyki platformy
 * (`AutomationWorkflow`). Scenariusz nagrany przy przeglądaniu jest właśnie
 * taką automatyką, więc okno zapisuje go tą samą drogą, którą czyta go moduł
 * Automations. Drugi, własny magazyn scenariuszy w module Browser byłby drugą
 * prawdą o tym samym bycie.
 *
 * Stąd też jedno nazewnictwo dla jednego pojęcia: `AutomationWorkflow` jest
 * w całym produkcie **automatyką**, `AutomationStep` — **krokiem**,
 * `AutomationSchedule` — **harmonogramem**, `AutomationExecution` —
 * **przebiegiem**. Nazwa okna (Automation Studio) pochodzi z opracowania
 * modułu i pozostaje nietknięta.
 *
 * Plik jest warstwą wywołań wraz ze sprawdzianem kształtu odpowiedzi: bez
 * stanu i bez elementów widoku. Żadne wywołanie nie rzuca wyjątkiem —
 * niepowodzenie wraca polem `blad` wyniku.
 */
export interface ZrodloAutomatyk {
  /** `automation.workflow.list` — automatyki Operatora. */
  wykaz(zadanie: AutomationWorkflowListRequest): Promise<Wynik<AutomationWorkflowListResponse>>;
  /** `automation.workflow.save` — zapis definicji wraz z krokami. */
  zapisz(zadanie: AutomationWorkflowSaveRequest): Promise<Wynik<AutomationWorkflowSaveResponse>>;
  /** `automation.schedule.set` — cykliczność i wyzwalacze automatyki. */
  ustawHarmonogram(
    zadanie: AutomationScheduleSetRequest,
  ): Promise<Wynik<AutomationScheduleSetResponse>>;
  /** `schedule.get` — odczyt harmonogramów do pary z zapisem. */
  harmonogramy(zadanie: ScheduleGetRequest): Promise<Wynik<ScheduleGetResponse>>;
  /** `automation.execution.subscribe` — stan przebiegów wraz z założeniem obserwacji. */
  przebiegi(
    zadanie: AutomationExecutionSubscribeRequest,
  ): Promise<Wynik<AutomationExecutionSubscribeResponse>>;
}

export function utworzZrodloAutomatyk(kanal: Kanal): ZrodloAutomatyk {
  return {
    async wykaz(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationWorkflowList, zadanie),
        Command.AutomationWorkflowList,
        (tresc) => czyTablica(tresc.workflows),
      );
    },

    async zapisz(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationWorkflowSave, zadanie),
        Command.AutomationWorkflowSave,
        (tresc) => czyObiekt(tresc.workflow),
      );
    },

    async ustawHarmonogram(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationScheduleSet, zadanie),
        Command.AutomationScheduleSet,
        (tresc) => czyObiekt(tresc.schedule),
      );
    },

    async harmonogramy(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.ScheduleGet, zadanie),
        Command.ScheduleGet,
        (tresc) => czyTablica(tresc.schedules),
      );
    },

    async przebiegi(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationExecutionSubscribe, zadanie),
        Command.AutomationExecutionSubscribe,
        (tresc) => czyTablica(tresc.executions),
      );
    },
  };
}
