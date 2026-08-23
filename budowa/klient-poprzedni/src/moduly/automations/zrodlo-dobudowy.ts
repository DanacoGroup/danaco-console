/**
 * Moduł Automations widziany przez klienta — czterdzieści trzy czynności
 * dopełniające `zrodlo-automations.ts`.
 *
 * Rozdział idzie po roli, nie po objętości. Tamten plik niesie rdzeń modułu:
 * definicję, harmonogram, kolejkę, zależności i przebiegi — czyli to, czym
 * okna pracują bez otwierania jakiegokolwiek panelu. Tutaj stoją czynności
 * paneli i szuflad: wersje, szablony, zmienne, kanwa, logi, ładunki, punkty
 * wznowienia, alarmy, budżety, skarbiec, audyt, zlecenia kolejki i nadzór
 * harmonogramu.
 *
 * Każda czynność oddaje `Wynik` z CAŁĄ odpowiedzią, nie z wyciętym polem.
 * Powód jest w kontrakcie: większość tych odpowiedzi niesie więcej niż jedną
 * rzecz — wykaz i znacznik przycięcia, automatykę i numer wersji, zlecenie
 * i znacznik duplikatu, zmienne i zastrzeżenia. Wycięcie jednego pola gubiłoby
 * drugie, a okno musi pokazać oba.
 *
 * Sprawdzian kształtu pilnuje POLA OBOWIĄZKOWEGO odpowiedzi. Rdzeń, który
 * oddał kopertę powodzenia bez treści, jest dla okna odmową — inaczej widok
 * rysowałby pustkę i twierdził, że to wynik.
 */
import {
  Command,
  type AutomationAlertRuleListRequest,
  type AutomationAlertRuleListResponse,
  type AutomationAlertRuleSetRequest,
  type AutomationAlertRuleSetResponse,
  type AutomationAuditListRequest,
  type AutomationAuditListResponse,
  type AutomationExecutionBudgetSetRequest,
  type AutomationExecutionBudgetSetResponse,
  type AutomationExecutionCheckpointListRequest,
  type AutomationExecutionCheckpointListResponse,
  type AutomationExecutionLogRequest,
  type AutomationExecutionLogResponse,
  type AutomationExecutionPayloadGetRequest,
  type AutomationExecutionPayloadGetResponse,
  type AutomationExecutionReplayRequest,
  type AutomationExecutionReplayResponse,
  type AutomationExecutionResumeRequest,
  type AutomationExecutionResumeResponse,
  type AutomationExecutionStepsRequest,
  type AutomationExecutionStepsResponse,
  type AutomationSecretListRequest,
  type AutomationSecretListResponse,
  type AutomationSecretRemoveRequest,
  type AutomationSecretRemoveResponse,
  type AutomationSecretSetRequest,
  type AutomationSecretSetResponse,
  type AutomationStepLayoutSetRequest,
  type AutomationStepLayoutSetResponse,
  type AutomationStepNoteSetRequest,
  type AutomationStepNoteSetResponse,
  type AutomationTemplateApplyRequest,
  type AutomationTemplateApplyResponse,
  type AutomationTemplateListRequest,
  type AutomationTemplateListResponse,
  type AutomationTemplateSaveRequest,
  type AutomationTemplateSaveResponse,
  type AutomationWorkflowPublishRequest,
  type AutomationWorkflowPublishResponse,
  type AutomationWorkflowShareRequest,
  type AutomationWorkflowShareResponse,
  type AutomationWorkflowSimulateRequest,
  type AutomationWorkflowSimulateResponse,
  type AutomationWorkflowTagSetRequest,
  type AutomationWorkflowTagSetResponse,
  type AutomationWorkflowVariablesSetRequest,
  type AutomationWorkflowVariablesSetResponse,
  type AutomationWorkflowVersionDiffRequest,
  type AutomationWorkflowVersionDiffResponse,
  type AutomationWorkflowVersionListRequest,
  type AutomationWorkflowVersionListResponse,
  type AutomationWorkflowVersionRestoreRequest,
  type AutomationWorkflowVersionRestoreResponse,
  type QueueActionRequest,
  type QueueActionResponse,
  type QueueDeadListRequest,
  type QueueDeadListResponse,
  type QueueDepthGetRequest,
  type QueueDepthGetResponse,
  type QueueItemBranchRequest,
  type QueueItemBranchResponse,
  type QueueItemConditionRequest,
  type QueueItemConditionResponse,
  type QueueItemDelayRequest,
  type QueueItemDelayResponse,
  type QueueItemDequeueRequest,
  type QueueItemDequeueResponse,
  type QueueItemEnqueueRequest,
  type QueueItemEnqueueResponse,
  type QueueItemListRequest,
  type QueueItemListResponse,
  type QueueItemMergeRequest,
  type QueueItemMergeResponse,
  type QueueItemRouteRequest,
  type QueueItemRouteResponse,
  type QueueItemSplitRequest,
  type QueueItemSplitResponse,
  type QueuePolicySetRequest,
  type QueuePolicySetResponse,
  type ScheduleBackfillRunRequest,
  type ScheduleBackfillRunResponse,
  type ScheduleHeartbeatSetRequest,
  type ScheduleHeartbeatSetResponse,
  type ScheduleTriggerHistoryRequest,
  type ScheduleTriggerHistoryResponse,
  type ScheduleWebhookEndpointGetRequest,
  type ScheduleWebhookEndpointGetResponse,
  type ScheduleWindowSetRequest,
  type ScheduleWindowSetResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/** Czterdzieści cztery czynności dopełniające moduł Automations. */
export interface ZrodloDobudowyAutomations {
  /**
   * `queue.action` — działanie silnika kolejek na kolejce BEZ automatyki.
   *
   * Odrębna od `automation.queue.action`, choć obie posuwają ten sam silnik.
   * Tamta wskazuje automatykę i zasila kolejkę jej krokami; ta działa na
   * kolejce samej — na przykład na kolejce przeglądu ręcznego, do której
   * skierowano zlecenie po błędzie. Queue Manager musi umieć posunąć jedną
   * i drugą, bo obie stoją w tym samym wykazie kolejek.
   */
  dzialanieNaKolejce(zadanie: QueueActionRequest): Promise<Wynik<QueueActionResponse>>;
  /** `automation.workflow.simulate` — Uruchamia przebieg probny definicji bez wpiecia produkcyjnego; efekty uboczne krokow sa wstrzymane. */
  symulujPrzeplyw(zadanie: AutomationWorkflowSimulateRequest): Promise<Wynik<AutomationWorkflowSimulateResponse>>;
  /** `automation.workflow.version.list` — Zwraca wersje definicji automatyki w kolejnosci od najnowszej. */
  wersjeAutomatyki(zadanie: AutomationWorkflowVersionListRequest): Promise<Wynik<AutomationWorkflowVersionListResponse>>;
  /** `automation.workflow.version.restore` — Przywraca wczesniejsza wersje definicji jako wersje biezaca; wersja zastana zostaje w historii. */
  przywrocWersje(zadanie: AutomationWorkflowVersionRestoreRequest): Promise<Wynik<AutomationWorkflowVersionRestoreResponse>>;
  /** `automation.workflow.version.diff` — Porownuje dwie wersje definicji strukturalnie: kroki dodane, usuniete i zmienione. */
  porownajWersje(zadanie: AutomationWorkflowVersionDiffRequest): Promise<Wynik<AutomationWorkflowVersionDiffResponse>>;
  /** `automation.workflow.tag.set` — Ustala komplet etykiet automatyki; wykaz pusty zdejmuje wszystkie. */
  ustawEtykiety(zadanie: AutomationWorkflowTagSetRequest): Promise<Wynik<AutomationWorkflowTagSetResponse>>;
  /** `automation.workflow.variables.set` — Ustala zmienne przeplywu i mapowanie danych miedzy krokami. */
  ustawZmienne(zadanie: AutomationWorkflowVariablesSetRequest): Promise<Wynik<AutomationWorkflowVariablesSetResponse>>;
  /** `automation.step.note.set` — Zapisuje notatke opisowa przy kroku automatyki; tresc pusta zdejmuje notatke. */
  ustawNotatkeKroku(zadanie: AutomationStepNoteSetRequest): Promise<Wynik<AutomationStepNoteSetResponse>>;
  /** `automation.step.layout.set` — Zapisuje polozenie wezlow krokow na kanwie; bez niego uklad kanwy liczy sie z zaleznosci i nie przezywa odczytu. */
  ustawUkladKanwy(zadanie: AutomationStepLayoutSetRequest): Promise<Wynik<AutomationStepLayoutSetResponse>>;
  /** `automation.template.save` — Zapisuje definicje jako szablon przeplywu wraz z jego parametrami. */
  zapiszSzablon(zadanie: AutomationTemplateSaveRequest): Promise<Wynik<AutomationTemplateSaveResponse>>;
  /** `automation.template.list` — Zwraca biblioteke szablonow przeplywow dostepnych Operatorowi. */
  szablony(zadanie: AutomationTemplateListRequest): Promise<Wynik<AutomationTemplateListResponse>>;
  /** `automation.template.apply` — Zaklada automatyke z szablonu, podstawiajac wartosci jego parametrow. */
  zastosujSzablon(zadanie: AutomationTemplateApplyRequest): Promise<Wynik<AutomationTemplateApplyResponse>>;
  /** `automation.workflow.publish` — Rozdziela wersje robocza od opublikowanej; wykonywana produkcyjnie jest wersja opublikowana. */
  opublikujAutomatyke(zadanie: AutomationWorkflowPublishRequest): Promise<Wynik<AutomationWorkflowPublishResponse>>;
  /** `automation.workflow.share` — Udostepnia automatyke jako komponent wlasny w obrebie organizacji. */
  udostepnijAutomatyke(zadanie: AutomationWorkflowShareRequest): Promise<Wynik<AutomationWorkflowShareResponse>>;
  /** `schedule.window.set` — Ustala okna wykonania harmonogramu: przedzialy czasu, w ktorych uruchomienie nastepuje. */
  ustawOknaWykonania(zadanie: ScheduleWindowSetRequest): Promise<Wynik<ScheduleWindowSetResponse>>;
  /** `schedule.backfill.run` — Wykonuje przebiegi dla przeszlych, pominietych terminow w zadanym zakresie dat. */
  uruchomWstecznie(zadanie: ScheduleBackfillRunRequest): Promise<Wynik<ScheduleBackfillRunResponse>>;
  /** `schedule.trigger.history` — Zwraca rzeczywiste momenty wyzwolenia wraz z przyczyna: harmonogram, zdarzenie albo uruchomienie reczne. */
  historiaWyzwolen(zadanie: ScheduleTriggerHistoryRequest): Promise<Wynik<ScheduleTriggerHistoryResponse>>;
  /** `schedule.heartbeat.set` — Ustala nadzor obecnosci uruchomien: alarm, gdy oczekiwane uruchomienie nie nastapilo w oknie tolerancji. */
  ustawNadzorUruchomien(zadanie: ScheduleHeartbeatSetRequest): Promise<Wynik<ScheduleHeartbeatSetResponse>>;
  /** `schedule.webhook.endpoint.get` — Zwraca unikatowy adres wejsciowy wyzwalacza webhook wraz z kluczem podpisu HMAC. */
  adresWebhooka(zadanie: ScheduleWebhookEndpointGetRequest): Promise<Wynik<ScheduleWebhookEndpointGetResponse>>;
  /** `queue.item.enqueue` — Dokłada zlecenie do kolejki istniejacej, poza harmonogramem. */
  dodajZlecenie(zadanie: QueueItemEnqueueRequest): Promise<Wynik<QueueItemEnqueueResponse>>;
  /** `queue.item.dequeue` — Zdejmuje zlecenie z kolejki przed jego wykonaniem. */
  zdejmijZlecenie(zadanie: QueueItemDequeueRequest): Promise<Wynik<QueueItemDequeueResponse>>;
  /** `queue.item.delay` — Odklada wykonanie zlecenia o wskazany czas. */
  odlozZlecenie(zadanie: QueueItemDelayRequest): Promise<Wynik<QueueItemDelayResponse>>;
  /** `queue.item.split` — Dzieli zlecenie na podzadania; zlecenie zrodlowe zostaje zamkniete. */
  podzielZlecenie(zadanie: QueueItemSplitRequest): Promise<Wynik<QueueItemSplitResponse>>;
  /** `queue.item.merge` — Scala kilka zlecen w jedno; zlecenia zrodlowe zostaja zamkniete. */
  scalZlecenia(zadanie: QueueItemMergeRequest): Promise<Wynik<QueueItemMergeResponse>>;
  /** `queue.item.route` — Kieruje zlecenie do innej kolejki albo do innego wykonawcy. */
  skierujZlecenie(zadanie: QueueItemRouteRequest): Promise<Wynik<QueueItemRouteResponse>>;
  /** `queue.item.branch` — Rozgalezia przetwarzanie zlecenia na tory rownolegle. */
  rozgalezZlecenie(zadanie: QueueItemBranchRequest): Promise<Wynik<QueueItemBranchResponse>>;
  /** `queue.item.condition` — Ustala warunek przetworzenia zlecenia; zlecenie niespelniajace warunku zostaje pominiete. */
  uwarunkujZlecenie(zadanie: QueueItemConditionRequest): Promise<Wynik<QueueItemConditionResponse>>;
  /** `queue.item.list` — Zwraca zlecenia kolejki wraz z ich stanem, ladunkiem i dotychczasowymi probami. */
  zleceniaKolejki(zadanie: QueueItemListRequest): Promise<Wynik<QueueItemListResponse>>;
  /** `queue.policy.set` — Ustala zasieg kolejki, wspolbieznosc, przepustowosc i polityke ponawiania. */
  ustawPolitykeKolejki(zadanie: QueuePolicySetRequest): Promise<Wynik<QueuePolicySetResponse>>;
  /** `queue.dead.list` — Zwraca zlecenia trwale nieudane, przeniesione do kolejki zadan martwych. */
  zadaniaMartwe(zadanie: QueueDeadListRequest): Promise<Wynik<QueueDeadListResponse>>;
  /** `queue.depth.get` — Zwraca glebokosc kolejki w czasie: liczbe zlecen oczekujacych w kolejnych odcinkach. */
  glebokoscKolejki(zadanie: QueueDepthGetRequest): Promise<Wynik<QueueDepthGetResponse>>;
  /** `automation.execution.log` — Zwraca pelny zapis zdarzen pojedynczego uruchomienia, takze sprzed otwarcia okna. */
  dziennikPrzebiegu(zadanie: AutomationExecutionLogRequest): Promise<Wynik<AutomationExecutionLogResponse>>;
  /** `automation.execution.steps` — Zwraca stan kazdego kroku przebiegu osobno. */
  krokiPrzebiegu(zadanie: AutomationExecutionStepsRequest): Promise<Wynik<AutomationExecutionStepsResponse>>;
  /** `automation.execution.checkpoint.list` — Zwraca punkty wznowienia przebiegu. */
  punktyWznowienia(zadanie: AutomationExecutionCheckpointListRequest): Promise<Wynik<AutomationExecutionCheckpointListResponse>>;
  /** `automation.execution.resume` — Wznawia przebieg od punktu wznowienia, bez powtarzania krokow ukonczonych. */
  wznowPrzebieg(zadanie: AutomationExecutionResumeRequest): Promise<Wynik<AutomationExecutionResumeResponse>>;
  /** `automation.execution.payload.get` — Zwraca dane wejsciowe i wyjsciowe kroku przebiegu; wartosci wrazliwe sa zredagowane. */
  podgladLadunku(zadanie: AutomationExecutionPayloadGetRequest): Promise<Wynik<AutomationExecutionPayloadGetResponse>>;
  /** `automation.execution.replay` — Odtwarza przebieg z ladunkiem kroku wskazanego; zaklada przebieg nowy, nie zmienia zrodlowego. */
  odtworzPrzebieg(zadanie: AutomationExecutionReplayRequest): Promise<Wynik<AutomationExecutionReplayResponse>>;
  /** `automation.alert.rule.set` — Ustala regule alarmowania: warunek i kanaly powiadomien. */
  ustawReguleAlarmowania(zadanie: AutomationAlertRuleSetRequest): Promise<Wynik<AutomationAlertRuleSetResponse>>;
  /** `automation.alert.rule.list` — Zwraca reguly alarmowania automatyki. */
  regulyAlarmowania(zadanie: AutomationAlertRuleListRequest): Promise<Wynik<AutomationAlertRuleListResponse>>;
  /** `automation.execution.budget.set` — Ustala budzet czasu przebiegu i kroku wraz z alarmem przy jego przekroczeniu. */
  ustawBudzetyPrzebiegu(zadanie: AutomationExecutionBudgetSetRequest): Promise<Wynik<AutomationExecutionBudgetSetResponse>>;
  /** `automation.secret.set` — Zapisuje poswiadczenie w skarbcu i zwraca jego referencje; wartosc nie wraca nigdy. */
  zapiszPoswiadczenie(zadanie: AutomationSecretSetRequest): Promise<Wynik<AutomationSecretSetResponse>>;
  /** `automation.secret.list` — Zwraca referencje poswiadczen dostepnych krokom; wartosci nie wracaja. */
  poswiadczenia(zadanie: AutomationSecretListRequest): Promise<Wynik<AutomationSecretListResponse>>;
  /** `automation.secret.remove` — Usuwa poswiadczenie ze skarbca; kroki przywolujace je przestaja miec pokrycie. */
  usunPoswiadczenie(zadanie: AutomationSecretRemoveRequest): Promise<Wynik<AutomationSecretRemoveResponse>>;
  /** `automation.audit.list` — Zwraca dziennik audytu: kto i kiedy zmienil definicje, uruchomil przebieg, przerwal go albo zmienil harmonogram. */
  odczytajAudyt(zadanie: AutomationAuditListRequest): Promise<Wynik<AutomationAuditListResponse>>;
}

export function utworzZrodloDobudowy(kanal: Kanal): ZrodloDobudowyAutomations {
  return {
    dzialanieNaKolejce(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueueAction, zadanie),
        Command.QueueAction,
        (tresc) => czyObiekt(tresc.queue),
      );
    },
    symulujPrzeplyw(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationWorkflowSimulate, zadanie),
        Command.AutomationWorkflowSimulate,
        (tresc) => czyTablica(tresc.results),
      );
    },
    wersjeAutomatyki(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationWorkflowVersionList, zadanie),
        Command.AutomationWorkflowVersionList,
        (tresc) => czyTablica(tresc.versions),
      );
    },
    przywrocWersje(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationWorkflowVersionRestore, zadanie),
        Command.AutomationWorkflowVersionRestore,
        (tresc) => czyObiekt(tresc.workflow),
      );
    },
    porownajWersje(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationWorkflowVersionDiff, zadanie),
        Command.AutomationWorkflowVersionDiff,
        (tresc) => czyTablica(tresc.changes),
      );
    },
    ustawEtykiety(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationWorkflowTagSet, zadanie),
        Command.AutomationWorkflowTagSet,
        (tresc) => czyObiekt(tresc.workflow),
      );
    },
    ustawZmienne(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationWorkflowVariablesSet, zadanie),
        Command.AutomationWorkflowVariablesSet,
        (tresc) => czyTablica(tresc.variables),
      );
    },
    ustawNotatkeKroku(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationStepNoteSet, zadanie),
        Command.AutomationStepNoteSet,
        (tresc) => czyObiekt(tresc.step),
      );
    },
    ustawUkladKanwy(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationStepLayoutSet, zadanie),
        Command.AutomationStepLayoutSet,
        (tresc) => czyTablica(tresc.positions),
      );
    },
    zapiszSzablon(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationTemplateSave, zadanie),
        Command.AutomationTemplateSave,
        (tresc) => czyObiekt(tresc.template),
      );
    },
    szablony(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationTemplateList, zadanie),
        Command.AutomationTemplateList,
        (tresc) => czyTablica(tresc.templates),
      );
    },
    zastosujSzablon(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationTemplateApply, zadanie),
        Command.AutomationTemplateApply,
        (tresc) => czyObiekt(tresc.workflow),
      );
    },
    opublikujAutomatyke(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationWorkflowPublish, zadanie),
        Command.AutomationWorkflowPublish,
        (tresc) => czyObiekt(tresc.workflow),
      );
    },
    udostepnijAutomatyke(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationWorkflowShare, zadanie),
        Command.AutomationWorkflowShare,
        (tresc) => czyObiekt(tresc.workflow),
      );
    },
    ustawOknaWykonania(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.ScheduleWindowSet, zadanie),
        Command.ScheduleWindowSet,
        (tresc) => czyObiekt(tresc.schedule),
      );
    },
    uruchomWstecznie(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.ScheduleBackfillRun, zadanie),
        Command.ScheduleBackfillRun,
        (tresc) => czyTablica(tresc.queuedAt),
      );
    },
    historiaWyzwolen(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.ScheduleTriggerHistory, zadanie),
        Command.ScheduleTriggerHistory,
        (tresc) => czyTablica(tresc.entries),
      );
    },
    ustawNadzorUruchomien(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.ScheduleHeartbeatSet, zadanie),
        Command.ScheduleHeartbeatSet,
        (tresc) => czyObiekt(tresc.schedule),
      );
    },
    adresWebhooka(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.ScheduleWebhookEndpointGet, zadanie),
        Command.ScheduleWebhookEndpointGet,
        (tresc) => typeof tresc.endpointUrl === 'string',
      );
    },
    dodajZlecenie(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueueItemEnqueue, zadanie),
        Command.QueueItemEnqueue,
        (tresc) => czyObiekt(tresc.item),
      );
    },
    zdejmijZlecenie(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueueItemDequeue, zadanie),
        Command.QueueItemDequeue,
        (tresc) => typeof tresc.removed === 'boolean',
      );
    },
    odlozZlecenie(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueueItemDelay, zadanie),
        Command.QueueItemDelay,
        (tresc) => czyObiekt(tresc.item),
      );
    },
    podzielZlecenie(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueueItemSplit, zadanie),
        Command.QueueItemSplit,
        (tresc) => czyTablica(tresc.items),
      );
    },
    scalZlecenia(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueueItemMerge, zadanie),
        Command.QueueItemMerge,
        (tresc) => czyObiekt(tresc.item),
      );
    },
    skierujZlecenie(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueueItemRoute, zadanie),
        Command.QueueItemRoute,
        (tresc) => czyObiekt(tresc.item),
      );
    },
    rozgalezZlecenie(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueueItemBranch, zadanie),
        Command.QueueItemBranch,
        (tresc) => czyTablica(tresc.items),
      );
    },
    uwarunkujZlecenie(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueueItemCondition, zadanie),
        Command.QueueItemCondition,
        (tresc) => czyObiekt(tresc.item),
      );
    },
    zleceniaKolejki(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueueItemList, zadanie),
        Command.QueueItemList,
        (tresc) => czyTablica(tresc.items),
      );
    },
    ustawPolitykeKolejki(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueuePolicySet, zadanie),
        Command.QueuePolicySet,
        (tresc) => czyObiekt(tresc.queue),
      );
    },
    zadaniaMartwe(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueueDeadList, zadanie),
        Command.QueueDeadList,
        (tresc) => czyTablica(tresc.items),
      );
    },
    glebokoscKolejki(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.QueueDepthGet, zadanie),
        Command.QueueDepthGet,
        (tresc) => czyTablica(tresc.points),
      );
    },
    dziennikPrzebiegu(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationExecutionLog, zadanie),
        Command.AutomationExecutionLog,
        (tresc) => czyTablica(tresc.entries),
      );
    },
    krokiPrzebiegu(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationExecutionSteps, zadanie),
        Command.AutomationExecutionSteps,
        (tresc) => czyTablica(tresc.steps),
      );
    },
    punktyWznowienia(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationExecutionCheckpointList, zadanie),
        Command.AutomationExecutionCheckpointList,
        (tresc) => czyTablica(tresc.checkpoints),
      );
    },
    wznowPrzebieg(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationExecutionResume, zadanie),
        Command.AutomationExecutionResume,
        (tresc) => czyObiekt(tresc.execution),
      );
    },
    podgladLadunku(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationExecutionPayloadGet, zadanie),
        Command.AutomationExecutionPayloadGet,
        (tresc) => czyObiekt(tresc.input),
      );
    },
    odtworzPrzebieg(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationExecutionReplay, zadanie),
        Command.AutomationExecutionReplay,
        (tresc) => czyObiekt(tresc.execution),
      );
    },
    ustawReguleAlarmowania(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationAlertRuleSet, zadanie),
        Command.AutomationAlertRuleSet,
        (tresc) => czyObiekt(tresc.rule),
      );
    },
    regulyAlarmowania(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationAlertRuleList, zadanie),
        Command.AutomationAlertRuleList,
        (tresc) => czyTablica(tresc.rules),
      );
    },
    ustawBudzetyPrzebiegu(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationExecutionBudgetSet, zadanie),
        Command.AutomationExecutionBudgetSet,
        (tresc) => czyObiekt(tresc.workflow),
      );
    },
    zapiszPoswiadczenie(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationSecretSet, zadanie),
        Command.AutomationSecretSet,
        (tresc) => czyObiekt(tresc.secret),
      );
    },
    poswiadczenia(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationSecretList, zadanie),
        Command.AutomationSecretList,
        (tresc) => czyTablica(tresc.secrets),
      );
    },
    usunPoswiadczenie(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationSecretRemove, zadanie),
        Command.AutomationSecretRemove,
        (tresc) => typeof tresc.removed === 'boolean',
      );
    },
    odczytajAudyt(zadanie) {
      return sprawdzOdpowiedz(
        wywolaj(kanal, Command.AutomationAuditList, zadanie),
        Command.AutomationAuditList,
        (tresc) => czyTablica(tresc.entries),
      );
    },
  };
}

/**
 * Sprawdzian kształtu odpowiedzi oddawanej w całości.
 *
 * Osobny od `sprawdzKsztaltObietnicy` z `zrodlo-automations.ts`, choć robi to
 * samo: tamten jest prywatny w swoim pliku, a wyniesienie go do modułu
 * wspólnego związałoby oba pliki bez powodu — sprawdzian jest trzema wierszami,
 * a wspólny byt do utrzymania jest zawsze czymś więcej.
 */
async function sprawdzOdpowiedz<T>(
  obietnica: Promise<Wynik<T>>,
  komenda: string,
  sprawdzian: (tresc: T) => boolean,
): Promise<Wynik<T>> {
  return sprawdzKsztalt(await obietnica, komenda, sprawdzian);
}
