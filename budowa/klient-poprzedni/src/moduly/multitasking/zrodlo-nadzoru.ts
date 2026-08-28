import {
  Command,
  EventType,
  type AodObserveAttachRequest,
  type AodObserveAttachResponse,
  type AodObserveDetachRequest,
  type AodObserveDetachResponse,
  type AodStatus,
  type AodStatusGetRequest,
  type AutomationDependency,
  type AutomationSchedule,
  type AutomationScheduleSetRequest,
  type AutomationWorkflow,
  type AutomationWorkflowListRequest,
  type OrchestrationDependencyListRequest,
  type OrchestrationDependencyRemoveRequest,
  type OrchestrationDependencyRemoveResponse,
  type OrchestrationDependencySetRequest,
  type OrchestrationDependencySetResponse,
  type OrchestrationValidateResponse,
  type ProgressChangedEvent,
  type Queue,
  type QueueLinkRequest,
  type ScheduleGetRequest,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Źródło trzech sekcji panelu orkiestracji — orkiestracji, harmonogramu
 * i nadzoru — dzieli jeden wykaz automatyk i przełącznik trybu nakładki,
 * którego kontrakt nie niesie, więc sekcja utrwala go ustawieniem sesji.
 */
export interface ZrodloNadzoru {
  /** `orchestration.dependency.list` — zależności kroków układu. */
  zaleznosci(
    zadanie: OrchestrationDependencyListRequest,
  ): Promise<Wynik<AutomationDependency[]>>;
  /** `orchestration.dependency.set` — dopisanie albo zmiana zależności. */
  zapiszZaleznosc(
    zadanie: OrchestrationDependencySetRequest,
  ): Promise<Wynik<OrchestrationDependencySetResponse>>;
  /** `orchestration.dependency.remove` — zdjęcie zależności między krokami. */
  zdejmijZaleznosc(
    zadanie: OrchestrationDependencyRemoveRequest,
  ): Promise<Wynik<OrchestrationDependencyRemoveResponse>>;
  /** `orchestration.validate` — sprawdzenie układu wraz ze ścieżką krytyczną. */
  sprawdzUklad(idUkladu: string): Promise<Wynik<OrchestrationValidateResponse>>;
  /** `automation.workflow.list` — automatyki, na których stoją trzy sekcje. */
  automatyki(zadanie: AutomationWorkflowListRequest): Promise<Wynik<AutomationWorkflow[]>>;
  /** `schedule.get` — harmonogramy; odczyt do pary z zapisem. */
  harmonogramy(zadanie: ScheduleGetRequest): Promise<Wynik<AutomationSchedule[]>>;
  /** `automation.schedule.set` — reguła czasowa pracy ciągłej. */
  zapiszHarmonogram(zadanie: AutomationScheduleSetRequest): Promise<Wynik<AutomationSchedule>>;
  /** `queue.link` — powiązanie kolejki z oknami, ekspertem, projektem, automatyką. */
  powiazKolejke(zadanie: QueueLinkRequest): Promise<Wynik<Queue>>;
  /** `aod.status.get` — stan nakładki Always On Display. */
  stanNakladki(zadanie: AodStatusGetRequest): Promise<Wynik<AodStatus>>;
  /** `aod.observe.attach` — przypięcie procesu do obserwacji nakładki. */
  przypnijDoNakladki(
    zadanie: AodObserveAttachRequest,
  ): Promise<Wynik<AodObserveAttachResponse>>;
  /** `aod.observe.detach` — odpięcie procesu od obserwacji nakładki. */
  odepnijOdNakladki(
    zadanie: AodObserveDetachRequest,
  ): Promise<Wynik<AodObserveDetachResponse>>;
  /** Subskrypcja `progress.changed` — postęp procesu na żywo. */
  naPostep(sluchacz: (tresc: ProgressChangedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloNadzoru(kanal: Kanal): ZrodloNadzoru {
  return {
    async zaleznosci(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.OrchestrationDependencyList, zadanie),
        Command.OrchestrationDependencyList,
        (tresc) => czyTablica(tresc.dependencies),
      );
      return przenies(wynik, (tresc) => tresc.dependencies);
    },

    async zapiszZaleznosc(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.OrchestrationDependencySet, zadanie),
        Command.OrchestrationDependencySet,
        (tresc) => czyTablica(tresc.dependencies) && typeof tresc.valid === 'boolean',
      );
    },

    async zdejmijZaleznosc(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.OrchestrationDependencyRemove, zadanie),
        Command.OrchestrationDependencyRemove,
        (tresc) => czyTablica(tresc.dependencies) && typeof tresc.removed === 'boolean',
      );
    },

    async sprawdzUklad(idUkladu) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.OrchestrationValidate, { workflowId: idUkladu }),
        Command.OrchestrationValidate,
        (tresc) => typeof tresc.valid === 'boolean',
      );
    },

    async automatyki(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationWorkflowList, zadanie),
        Command.AutomationWorkflowList,
        (tresc) => czyTablica(tresc.workflows),
      );
      return przenies(wynik, (tresc) => tresc.workflows);
    },

    async harmonogramy(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ScheduleGet, zadanie),
        Command.ScheduleGet,
        (tresc) => czyTablica(tresc.schedules),
      );
      return przenies(wynik, (tresc) => tresc.schedules);
    },

    async zapiszHarmonogram(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationScheduleSet, zadanie),
        Command.AutomationScheduleSet,
        (tresc) => czyObiekt(tresc.schedule),
      );
      return przenies(wynik, (tresc) => tresc.schedule);
    },

    async powiazKolejke(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.QueueLink, zadanie),
        Command.QueueLink,
        (tresc) => czyObiekt(tresc.queue),
      );
      return przenies(wynik, (tresc) => tresc.queue);
    },

    async stanNakladki(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AodStatusGet, zadanie),
        Command.AodStatusGet,
        (tresc) => czyObiekt(tresc.status),
      );
      return przenies(wynik, (tresc) => tresc.status);
    },

    async przypnijDoNakladki(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AodObserveAttach, zadanie),
        Command.AodObserveAttach,
        (tresc) => czyTablica(tresc.attachedProcessIds),
      );
    },

    async odepnijOdNakladki(zadanie) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.AodObserveDetach, zadanie),
        Command.AodObserveDetach,
        (tresc) => czyTablica(tresc.attachedProcessIds),
      );
    },

    naPostep: (sluchacz) => kanal.naZdarzenie(EventType.ProgressChanged, sluchacz),
  };
}
