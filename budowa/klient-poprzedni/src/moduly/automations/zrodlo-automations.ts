import {
  Command,
  EventType,
  type AutomationExecution,
  type AutomationExecutionStatusEvent,
  type AutomationExecutionSubscribeRequest,
  type AutomationOrchestratorDefineRequest,
  type AutomationOrchestratorDefineResponse,
  type AutomationQueueActionRequest,
  type AutomationSchedule,
  type AutomationScheduleSetRequest,
  type AutomationScheduleSetResponse,
  type AutomationWorkflow,
  type AutomationWorkflowListRequest,
  type AutomationWorkflowSaveRequest,
  SessionConfigArea,
  type AutomationDependency,
  type AutomationLinkChangedEvent,
  type ConfigEffectiveGetRequest,
  type OrchestrationChangedEvent,
  type OrchestrationDependencyListRequest,
  type OrchestrationDependencyRemoveRequest,
  type OrchestrationDependencyRemoveResponse,
  type OrchestrationDependencySetRequest,
  type OrchestrationDependencySetResponse,
  type OrchestrationValidateResponse,
  type OrchestrationCompensationSetRequest,
  type OrchestrationCompensationSetResponse,
  type OrchestrationGateSetRequest,
  type OrchestrationGateSetResponse,
  type OrchestrationGroupSetRequest,
  type OrchestrationGroupSetResponse,
  type OrchestrationMultitaskingLinkRequest,
  type OrchestrationMultitaskingLinkResponse,
  type Queue,
  type QueueChangedEvent,
  type QueueCreateRequest,
  type QueueLinkRequest,
  type QueueListRequest,
  type ScheduleGetRequest,
  type SessionConfigEffective,
  type Window,
  type WindowListRequest,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import {
  utworzZrodloDobudowy,
  type ZrodloDobudowyAutomations,
} from './zrodlo-dobudowy';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Moduł Automations widziany przez klienta — komendy obszaru `automation.*`
 * wraz z komendami kolejek, harmonogramów i orkiestracji, których wprost żądają
 * okna modułu.
 *
 * Każda czynność oddaje `Wynik`, nie samą treść: okna mają obowiązkowy stan
 * błędu, więc źródło nie połyka odmowy ani nie podstawia pustego wykazu w jej
 * miejsce — widok musi odróżnić „nic nie ma” od „nie udało się zapytać”.
 *
 * Po stronie klienta nie stoi drugi silnik kolejek: Queue Manager zakłada
 * kolejkę komendą `queue.create` i posuwa ją komendą `automation.queue.action`,
 * która po stronie rdzenia jedzie tym samym adapterem kolejek co pętla sesyjna.
 */
export interface ZrodloAutomations extends ZrodloDobudowyAutomations {
  /** `automation.workflow.save` — zapis definicji automatyki wraz z krokami. */
  zapiszAutomatyke(zadanie: AutomationWorkflowSaveRequest): Promise<Wynik<AutomationWorkflow>>;
  /** `automation.workflow.list` — wykaz automatyk. */
  automatyki(zadanie: AutomationWorkflowListRequest): Promise<Wynik<AutomationWorkflow[]>>;
  /** `automation.schedule.set` — cykliczność i wyzwalacze automatyki. */
  ustawHarmonogram(
    zadanie: AutomationScheduleSetRequest,
  ): Promise<Wynik<AutomationScheduleSetResponse>>;
  /** `automation.orchestrator.define` — układ zależności i jego ocena. */
  ustawZaleznosci(
    zadanie: AutomationOrchestratorDefineRequest,
  ): Promise<Wynik<AutomationOrchestratorDefineResponse>>;
  /** `automation.execution.subscribe` — przebiegi wraz z założeniem obserwacji. */
  przebiegi(zadanie: AutomationExecutionSubscribeRequest): Promise<
    Wynik<{ executions: AutomationExecution[]; subscribed: boolean }>
  >;
  /** `queue.create` — założenie kolejki wykonującej automatykę. */
  zalozKolejke(zadanie: QueueCreateRequest): Promise<Wynik<Queue>>;
  /** `automation.queue.action` — działanie silnika kolejek na kolejce automatyki. */
  dzialanieKolejki(zadanie: AutomationQueueActionRequest): Promise<Wynik<Queue>>;
  /** `queue.list` — wykaz kolejek jednego silnika. */
  wykazKolejek(zadanie: QueueListRequest): Promise<Wynik<Queue[]>>;
  /** `queue.link` — powiązanie kolejki z oknami, ekspertem, projektem albo automatyką. */
  powiazKolejke(zadanie: QueueLinkRequest): Promise<Wynik<Queue>>;
  /** `schedule.get` — odczyt harmonogramów do pary z `automation.schedule.set`. */
  odczytajHarmonogramy(zadanie: ScheduleGetRequest): Promise<Wynik<AutomationSchedule[]>>;
  /** `orchestration.dependency.remove` — usunięcie pojedynczej zależności układu. */
  usunZaleznosc(
    zadanie: OrchestrationDependencyRemoveRequest,
  ): Promise<Wynik<OrchestrationDependencyRemoveResponse>>;
  /**
   * `orchestration.dependency.list` — zależności układu, opcjonalnie zawężone
   * do jednego kroku.
   *
   * Odrębna od `automation.orchestrator.define`, choć obie mówią o tym samym
   * grafie: `define` zapisuje układ w całości i oddaje go przy okazji, a `list`
   * wyłącznie czyta. Sam odczyt jest potrzebny, bo okno musi pokazać układ,
   * którego nikt wcześniej nie zapisywał z tego okna.
   */
  wykazZaleznosci(
    zadanie: OrchestrationDependencyListRequest,
  ): Promise<Wynik<{ dependencies: AutomationDependency[] }>>;
  /**
   * `orchestration.dependency.set` — zapis albo zmiana pojedynczej zależności.
   *
   * Zapis nie jest odmawiany: układ z cyklem zostaje zapisany, a zastrzeżenia
   * wracają w wyniku. Ocena układu nie jest więc oceną samej czynności — okno
   * rozdziela te dwa zdania.
   */
  zapiszZaleznosc(
    zadanie: OrchestrationDependencySetRequest,
  ): Promise<Wynik<OrchestrationDependencySetResponse>>;
  /** `orchestration.validate` — cykle, kroki osierocone i ścieżka krytyczna. */
  sprawdzUklad(idUkladu: string): Promise<Wynik<OrchestrationValidateResponse>>;

  /**
   * Cztery dopełnienia układu, których krawędź nie wyraża.
   *
   * Krawędź mówi, że kroki się schodzą. Bramka mówi, KIEDY tory scalają się
   * w jednym kroku; grupa mówi, że zbiór kroków biegnie razem; kompensacja
   * mówi, co zrobić, gdy przebieg pękł w pół; spięcie mówi, że kolejki
   * automatyki prowadzi silnik środowiska MultitaskingAI.
   */
  ustawBramke(zadanie: OrchestrationGateSetRequest): Promise<Wynik<OrchestrationGateSetResponse>>;
  ustawGrupe(zadanie: OrchestrationGroupSetRequest): Promise<Wynik<OrchestrationGroupSetResponse>>;
  ustawKompensacje(
    zadanie: OrchestrationCompensationSetRequest,
  ): Promise<Wynik<OrchestrationCompensationSetResponse>>;
  spnijZMultitaskingiem(
    zadanie: OrchestrationMultitaskingLinkRequest,
  ): Promise<Wynik<OrchestrationMultitaskingLinkResponse>>;
  /**
   * `window.list` — okna komunikacji platformy.
   *
   * Moduł sięga po nie w jednej sprawie: przekazanie scenariusza z innego
   * modułu (`context.transfer`) zakłada okno modułu Automations i to w jego
   * konfiguracji mieszka przeniesiony komplet. Bez wykazu okien nie ma jak
   * dojść do tego, komu przekazano.
   */
  oknaKomunikacji(zadanie: WindowListRequest): Promise<Wynik<Window[]>>;
  /**
   * `config.effective.get` zawężone do obszaru kontekstu rozmowy — nośnika
   * przeniesionego kompletu (`SessionConfigConversationContext.transferredContext`).
   */
  kontekstOkna(zadanie: ConfigEffectiveGetRequest): Promise<Wynik<SessionConfigEffective>>;
  /** Subskrypcja `automation.execution.status` — Execution Monitor na żywo. */
  naStanPrzebiegu(sluchacz: (tresc: AutomationExecutionStatusEvent) => void): Odsubskrybuj;
  /** Subskrypcja `queue.changed` — Queue Manager na żywo. */
  naZmianeKolejki(sluchacz: (tresc: QueueChangedEvent) => void): Odsubskrybuj;
  /**
   * Subskrypcja `automation.link.changed` — zmiana powiązania automatyki
   * z bytem wyzwalającym.
   *
   * Zdarzenie przychodzi także wtedy, gdy scenariusz wpiął tu inny moduł, więc
   * to ono jest sygnałem, że wykaz automatyk okna jest już nieaktualny.
   */
  naZmianePowiazania(sluchacz: (tresc: AutomationLinkChangedEvent) => void): Odsubskrybuj;
  /** Subskrypcja `orchestration.changed` — układ zmieniony poza tym oknem. */
  naZmianeUkladu(sluchacz: (tresc: OrchestrationChangedEvent) => void): Odsubskrybuj;
}

/**
 * Obszar konfiguracji sesji niosący przeniesiony komplet. Stała stoi tutaj, bo
 * wskazuje ją i odczyt kontekstu okna, i sprawdzian tego odczytu.
 */
export const OBSZAR_KONTEKSTU_ROZMOWY = SessionConfigArea.ConversationContext;

export function utworzZrodloAutomations(kanal: Kanal): ZrodloAutomations {
  return {
    // Czynności paneli i szuflad stoją w `zrodlo-dobudowy.ts` — jedno źródło
    // widziane przez okna, dwa pliki wedle roli czynności.
    ...utworzZrodloDobudowy(kanal),

    async zapiszAutomatyke(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationWorkflowSave, zadanie),
        Command.AutomationWorkflowSave,
        (tresc) => czyObiekt(tresc.workflow),
      );
      return przenies(wynik, (tresc) => tresc.workflow);
    },

    async automatyki(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationWorkflowList, zadanie),
        Command.AutomationWorkflowList,
        (tresc) => czyTablica(tresc.workflows),
      );
      return przenies(wynik, (tresc) => tresc.workflows);
    },

    ustawHarmonogram(zadanie) {
      return sprawdzKsztaltObietnicy(
        wywolaj(kanal, Command.AutomationScheduleSet, zadanie),
        Command.AutomationScheduleSet,
        (tresc) => czyObiekt(tresc.schedule),
      );
    },

    ustawZaleznosci(zadanie) {
      return sprawdzKsztaltObietnicy(
        wywolaj(kanal, Command.AutomationOrchestratorDefine, zadanie),
        Command.AutomationOrchestratorDefine,
        (tresc) => czyTablica(tresc.dependencies),
      );
    },

    async przebiegi(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationExecutionSubscribe, zadanie),
        Command.AutomationExecutionSubscribe,
        (tresc) => czyTablica(tresc.executions),
      );
      return przenies(wynik, (tresc) => ({
        executions: tresc.executions,
        subscribed: tresc.subscribed,
      }));
    },

    async zalozKolejke(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.QueueCreate, zadanie),
        Command.QueueCreate,
        (tresc) => czyObiekt(tresc.queue),
      );
      return przenies(wynik, (tresc) => tresc.queue);
    },

    async dzialanieKolejki(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.AutomationQueueAction, zadanie),
        Command.AutomationQueueAction,
        (tresc) => czyObiekt(tresc.queue),
      );
      return przenies(wynik, (tresc) => tresc.queue);
    },

    async wykazKolejek(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.QueueList, zadanie),
        Command.QueueList,
        (tresc) => czyTablica(tresc.queues),
      );
      return przenies(wynik, (tresc) => tresc.queues);
    },

    async powiazKolejke(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.QueueLink, zadanie),
        Command.QueueLink,
        (tresc) => czyObiekt(tresc.queue),
      );
      return przenies(wynik, (tresc) => tresc.queue);
    },

    async odczytajHarmonogramy(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ScheduleGet, zadanie),
        Command.ScheduleGet,
        (tresc) => czyTablica(tresc.schedules),
      );
      return przenies(wynik, (tresc) => tresc.schedules);
    },

    usunZaleznosc(zadanie) {
      return sprawdzKsztaltObietnicy(
        wywolaj(kanal, Command.OrchestrationDependencyRemove, zadanie),
        Command.OrchestrationDependencyRemove,
        (tresc) => czyTablica(tresc.dependencies),
      );
    },

    wykazZaleznosci(zadanie) {
      return sprawdzKsztaltObietnicy(
        wywolaj(kanal, Command.OrchestrationDependencyList, zadanie),
        Command.OrchestrationDependencyList,
        (tresc) => czyTablica(tresc.dependencies),
      );
    },

    zapiszZaleznosc(zadanie) {
      return sprawdzKsztaltObietnicy(
        wywolaj(kanal, Command.OrchestrationDependencySet, zadanie),
        Command.OrchestrationDependencySet,
        (tresc) => czyTablica(tresc.dependencies) && typeof tresc.valid === 'boolean',
      );
    },

    sprawdzUklad(idUkladu) {
      return sprawdzKsztaltObietnicy(
        wywolaj(kanal, Command.OrchestrationValidate, { workflowId: idUkladu }),
        Command.OrchestrationValidate,
        (tresc) => typeof tresc.valid === 'boolean',
      );
    },

    ustawBramke(zadanie) {
      return sprawdzKsztaltObietnicy(
        wywolaj(kanal, Command.OrchestrationGateSet, zadanie),
        Command.OrchestrationGateSet,
        (tresc) => czyTablica(tresc.gates) && typeof tresc.valid === 'boolean',
      );
    },

    ustawGrupe(zadanie) {
      return sprawdzKsztaltObietnicy(
        wywolaj(kanal, Command.OrchestrationGroupSet, zadanie),
        Command.OrchestrationGroupSet,
        (tresc) => czyTablica(tresc.groups),
      );
    },

    ustawKompensacje(zadanie) {
      return sprawdzKsztaltObietnicy(
        wywolaj(kanal, Command.OrchestrationCompensationSet, zadanie),
        Command.OrchestrationCompensationSet,
        (tresc) => czyTablica(tresc.compensations),
      );
    },

    spnijZMultitaskingiem(zadanie) {
      return sprawdzKsztaltObietnicy(
        wywolaj(kanal, Command.OrchestrationMultitaskingLink, zadanie),
        Command.OrchestrationMultitaskingLink,
        (tresc) => typeof tresc.linked === 'boolean',
      );
    },

    async oknaKomunikacji(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WindowList, zadanie),
        Command.WindowList,
        (tresc) => czyTablica(tresc.windows),
      );
      return przenies(wynik, (tresc) => tresc.windows);
    },

    async kontekstOkna(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ConfigEffectiveGet, zadanie),
        Command.ConfigEffectiveGet,
        (tresc) => czyObiekt(tresc.effective),
      );
      return przenies(wynik, (tresc) => tresc.effective);
    },

    naStanPrzebiegu(sluchacz) {
      return kanal.naZdarzenie(EventType.AutomationExecutionStatus, (tresc) => sluchacz(tresc));
    },

    naZmianeKolejki(sluchacz) {
      return kanal.naZdarzenie(EventType.QueueChanged, (tresc) => sluchacz(tresc));
    },

    naZmianePowiazania(sluchacz) {
      return kanal.naZdarzenie(EventType.AutomationLinkChanged, (tresc) => sluchacz(tresc));
    },

    naZmianeUkladu(sluchacz) {
      return kanal.naZdarzenie(EventType.OrchestrationChanged, (tresc) => sluchacz(tresc));
    },
  };
}

/** Sprawdzian kształtu odpowiedzi oddawanej w całości, bez wycinania pola. */
async function sprawdzKsztaltObietnicy<T>(
  obietnica: Promise<Wynik<T>>,
  komenda: string,
  sprawdzian: (tresc: T) => boolean,
): Promise<Wynik<T>> {
  return sprawdzKsztalt(await obietnica, komenda, sprawdzian);
}
