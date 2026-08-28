import {
  Command,
  EventType,
  type StudioAgentsConflictsListRequest,
  type StudioAgentsConflictsListResponse,
  type StudioAgentsSettingsGetRequest,
  type StudioAgentsSettingsGetResponse,
  type StudioAgentsSlotsListRequest,
  type StudioAgentsSlotsListResponse,
  type StudioBatchRunRequest,
  type StudioBatchRunResponse,
  type StudioChainListRequest,
  type StudioChainListResponse,
  type StudioChainProgressedEvent,
  type StudioChainRunRequest,
  type StudioChainRunResponse,
  type StudioChainSaveRequest,
  type StudioChainSaveResponse,
  type StudioOperationListRequest,
  type StudioOperationListResponse,
  type StudioPlanCreateRequest,
  type StudioPlanCreateResponse,
  type StudioPlanGetRequest,
  type StudioPlanGetResponse,
  type StudioPlanRunRequest,
  type StudioPlanRunResponse,
  type StudioPlanStopRequest,
  type StudioPlanStopResponse,
  type StudioPlanTaskUpdateRequest,
  type StudioPlanTaskUpdateResponse,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolajUczciwie } from './odmowa-rdzenia';

/** Wywołania okna pętli wykonawczej, dwanaście komend trzech rodzin: rozkład zadań, kontrola pracy wykonawców i łańcuchy operacji. */
export interface ZrodloPetli {
  /** `studio.plan.create` — rozkłada zlecenie dokumentowe na zadania. */
  rozlozZlecenie(zadanie: StudioPlanCreateRequest): Promise<Wynik<StudioPlanCreateResponse>>;
  /** `studio.plan.get` — rozkład wraz ze stanem każdego zadania. */
  rozklad(zadanie: StudioPlanGetRequest): Promise<Wynik<StudioPlanGetResponse>>;
  /** `studio.plan.task.update` — stan zadania, wykonawca, wynik. */
  przestawZadanie(
    zadanie: StudioPlanTaskUpdateRequest,
  ): Promise<Wynik<StudioPlanTaskUpdateResponse>>;
  /** `studio.plan.run` — puszcza pętlę; przy wyłączonej nastawie odmawia. */
  puscPetle(zadanie: StudioPlanRunRequest): Promise<Wynik<StudioPlanRunResponse>>;
  /** `studio.plan.stop` — zatrzymuje przebieg. */
  zatrzymaj(zadanie: StudioPlanStopRequest): Promise<Wynik<StudioPlanStopResponse>>;

  /** `studio.agents.settings.get` — czy pętla i praca wielu agentów są czynne. */
  nastawy(
    zadanie: StudioAgentsSettingsGetRequest,
  ): Promise<Wynik<StudioAgentsSettingsGetResponse>>;
  /** `studio.agents.slots.list` — kto pracuje nad którym fragmentem. */
  obsada(zadanie: StudioAgentsSlotsListRequest): Promise<Wynik<StudioAgentsSlotsListResponse>>;
  /** `studio.agents.conflicts.list` — spięcia wraz z odłożonym brzmieniem. */
  spiecia(
    zadanie: StudioAgentsConflictsListRequest,
  ): Promise<Wynik<StudioAgentsConflictsListResponse>>;

  /** `studio.chain.save` — zapisuje łańcuch operacji. */
  zapiszLancuch(zadanie: StudioChainSaveRequest): Promise<Wynik<StudioChainSaveResponse>>;
  /** `studio.chain.list` — łańcuchy dostępne Operatorowi. */
  lancuchy(zadanie: StudioChainListRequest): Promise<Wynik<StudioChainListResponse>>;
  /** `studio.chain.run` — rejestruje przebieg łańcucha na dokumencie. */
  uruchomLancuch(zadanie: StudioChainRunRequest): Promise<Wynik<StudioChainRunResponse>>;
  /** `studio.batch.run` — ta sama operacja na wskazanych dokumentach. */
  uruchomWsad(zadanie: StudioBatchRunRequest): Promise<Wynik<StudioBatchRunResponse>>;
  /** `studio.operation.list` — operacje własne Operatora do warsztatu łańcucha. */
  operacjeWlasne(
    zadanie: StudioOperationListRequest,
  ): Promise<Wynik<StudioOperationListResponse>>;

  /** Subskrypcja `studio.chain.progressed` — postęp kroku przebiegu. */
  naPostep(sluchacz: (tresc: StudioChainProgressedEvent) => void): Odsubskrybuj;
}

export function utworzZrodloPetli(kanal: Kanal): ZrodloPetli {
  return {
    async rozlozZlecenie(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPlanCreate, zadanie),
        Command.StudioPlanCreate,
        (tresc) => czyObiekt(tresc.plan),
      );
    },

    async rozklad(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPlanGet, zadanie),
        Command.StudioPlanGet,
        (tresc) => czyObiekt(tresc.plan),
      );
    },

    async przestawZadanie(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPlanTaskUpdate, zadanie),
        Command.StudioPlanTaskUpdate,
        (tresc) => czyObiekt(tresc.task) && czyObiekt(tresc.plan),
      );
    },

    async puscPetle(zadanie) {
      // Odmowa z powodem jest odpowiedzią poprawną: sprawdzian pilnuje obecności pola, nie jego wartości.
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPlanRun, zadanie),
        Command.StudioPlanRun,
        (tresc) => czyObiekt(tresc.plan) && typeof tresc.started === 'boolean',
      );
    },

    async zatrzymaj(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioPlanStop, zadanie),
        Command.StudioPlanStop,
        (tresc) => czyObiekt(tresc.plan) && typeof tresc.stoppedTasks === 'number',
      );
    },

    async nastawy(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioAgentsSettingsGet, zadanie),
        Command.StudioAgentsSettingsGet,
        (tresc) => czyObiekt(tresc.settings),
      );
    },

    async obsada(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioAgentsSlotsList, zadanie),
        Command.StudioAgentsSlotsList,
        (tresc) => czyTablica(tresc.slots),
      );
    },

    async spiecia(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioAgentsConflictsList, zadanie),
        Command.StudioAgentsConflictsList,
        (tresc) => czyTablica(tresc.conflicts),
      );
    },

    async zapiszLancuch(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioChainSave, zadanie),
        Command.StudioChainSave,
        (tresc) => czyObiekt(tresc.chain),
      );
    },

    async lancuchy(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioChainList, zadanie),
        Command.StudioChainList,
        (tresc) => czyTablica(tresc.chains),
      );
    },

    async uruchomLancuch(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioChainRun, zadanie),
        Command.StudioChainRun,
        (tresc) => typeof tresc.runId === 'string' && typeof tresc.steps === 'number',
      );
    },

    async uruchomWsad(zadanie) {
      // Wsad odrzucający wszystkie dokumenty ma kształt poprawny: bilans jest tu wynikiem wymaganym.
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioBatchRun, zadanie),
        Command.StudioBatchRun,
        (tresc) => typeof tresc.runId === 'string' && typeof tresc.accepted === 'number',
      );
    },

    async operacjeWlasne(zadanie) {
      return sprawdzKsztalt(
        await wywolajUczciwie(kanal, Command.StudioOperationList, zadanie),
        Command.StudioOperationList,
        (tresc) => czyTablica(tresc.operations),
      );
    },

    naPostep(sluchacz) {
      return kanal.naZdarzenie(EventType.StudioChainProgressed, (tresc) => sluchacz(tresc));
    },
  };
}
