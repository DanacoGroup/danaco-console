import {
  Command,
  type WorkspaceBoard,
  type WorkspaceBoardGetRequest,
  type WorkspaceCalendarExportRequest,
  type WorkspaceCalendarExportResponse,
  type WorkspaceCalendarGetRequest,
  type WorkspaceCalendarGetResponse,
  type WorkspaceCalendarImportRequest,
  type WorkspaceCalendarImportResponse,
  type WorkspaceScheduleGetRequest,
  type WorkspaceScheduleGetResponse,
  type WorkspaceTask,
  type WorkspaceTaskCreateRequest,
  type WorkspaceTaskCreateResponse,
  type WorkspaceTaskDeleteRequest,
  type WorkspaceTaskDeleteResponse,
  type WorkspaceTaskDependencyRemoveRequest,
  type WorkspaceTaskDependencyRemoveResponse,
  type WorkspaceTaskDependencySetRequest,
  type WorkspaceTaskDependencySetResponse,
  type WorkspaceTaskListRequest,
  type WorkspaceTaskMoveRequest,
  type WorkspaceTaskMoveResponse,
  type WorkspaceTaskUpdateRequest,
  type WorkspaceTaskUpdateResponse,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Hub planowania projektu widziany przez klienta grupuje dwanaście komend rodziny zadań, tablicy,
 * osi czasu i kalendarza; każda czynność oddaje wynik, nie samą treść.
 */
export interface CzynnosciPlanowania {
  /** `workspace.task.create` — założenie zadania projektu. */
  zalozZadanie(zadanie: WorkspaceTaskCreateRequest): Promise<Wynik<WorkspaceTaskCreateResponse>>;
  /** `workspace.task.update` — zmiana pól zadania; pominięte zostają bez zmiany. */
  zmienZadanie(zadanie: WorkspaceTaskUpdateRequest): Promise<Wynik<WorkspaceTaskUpdateResponse>>;
  /** `workspace.task.delete` — usunięcie zadania wraz z podzadaniami. */
  usunZadanie(zadanie: WorkspaceTaskDeleteRequest): Promise<Wynik<WorkspaceTaskDeleteResponse>>;
  /** `workspace.task.list` — wykaz zadań widoku listy. */
  zadania(zadanie: WorkspaceTaskListRequest): Promise<Wynik<WorkspaceTask[]>>;
  /** `workspace.task.move` — przeciągnięcie karty na tablicy. */
  przeniesZadanie(zadanie: WorkspaceTaskMoveRequest): Promise<Wynik<WorkspaceTaskMoveResponse>>;
  /** `workspace.task.dependency.set` — założenie zależności między zadaniami. */
  zalozZaleznosc(
    zadanie: WorkspaceTaskDependencySetRequest,
  ): Promise<Wynik<WorkspaceTaskDependencySetResponse>>;
  /** `workspace.task.dependency.remove` — zniesienie zależności. */
  zniesZaleznosc(
    zadanie: WorkspaceTaskDependencyRemoveRequest,
  ): Promise<Wynik<WorkspaceTaskDependencyRemoveResponse>>;
  /** `workspace.board.get` — tablica kanban wraz z kartami. */
  tablica(zadanie: WorkspaceBoardGetRequest): Promise<Wynik<WorkspaceBoard>>;
  /** `workspace.schedule.get` — słupki osi czasu i ścieżka krytyczna. */
  harmonogram(
    zadanie: WorkspaceScheduleGetRequest,
  ): Promise<Wynik<WorkspaceScheduleGetResponse>>;
  /** `workspace.calendar.get` — pozycje kalendarza w siatce. */
  kalendarz(zadanie: WorkspaceCalendarGetRequest): Promise<Wynik<WorkspaceCalendarGetResponse>>;
  /** `workspace.calendar.import` — wciągnięcie pliku iCal do projektu. */
  wciagnijKalendarz(
    zadanie: WorkspaceCalendarImportRequest,
  ): Promise<Wynik<WorkspaceCalendarImportResponse>>;
  /** `workspace.calendar.export` — zapis kalendarza projektu jako iCal. */
  zapiszKalendarz(
    zadanie: WorkspaceCalendarExportRequest,
  ): Promise<Wynik<WorkspaceCalendarExportResponse>>;
}

export function czynnosciPlanowania(kanal: Kanal): CzynnosciPlanowania {
  return {
    zalozZadanie: (zadanie) => wywolaj(kanal, Command.WorkspaceTaskCreate, zadanie),
    zmienZadanie: (zadanie) => wywolaj(kanal, Command.WorkspaceTaskUpdate, zadanie),
    usunZadanie: (zadanie) => wywolaj(kanal, Command.WorkspaceTaskDelete, zadanie),

    async zadania(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceTaskList, zadanie),
        Command.WorkspaceTaskList,
        (tresc) => czyTablica(tresc.tasks),
      );
      return przenies(wynik, (tresc) => tresc.tasks);
    },

    przeniesZadanie: (zadanie) => wywolaj(kanal, Command.WorkspaceTaskMove, zadanie),
    zalozZaleznosc: (zadanie) => wywolaj(kanal, Command.WorkspaceTaskDependencySet, zadanie),
    zniesZaleznosc: (zadanie) => wywolaj(kanal, Command.WorkspaceTaskDependencyRemove, zadanie),

    async tablica(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceBoardGet, zadanie),
        Command.WorkspaceBoardGet,
        (tresc) => czyObiekt(tresc.board),
      );
      return przenies(wynik, (tresc) => tresc.board);
    },

    harmonogram: (zadanie) => wywolaj(kanal, Command.WorkspaceScheduleGet, zadanie),
    kalendarz: (zadanie) => wywolaj(kanal, Command.WorkspaceCalendarGet, zadanie),
    wciagnijKalendarz: (zadanie) => wywolaj(kanal, Command.WorkspaceCalendarImport, zadanie),
    zapiszKalendarz: (zadanie) => wywolaj(kanal, Command.WorkspaceCalendarExport, zadanie),
  };
}

/**
 * Karty jednej kolumny tablicy w kolejności kluczy porządkowych: klucz jest napisem układanym
 * przez rdzeń przy przeciąganiu, a karta bez klucza idzie na koniec kolumny.
 */
export function kartyKolumny(tablica: WorkspaceBoard, idKolumny: string): WorkspaceTask[] {
  return tablica.tasks
    .filter((karta) => (karta.boardColumnId ?? '') === idKolumny)
    .sort((lewa, prawa) => (lewa.rank ?? '￿').localeCompare(prawa.rank ?? '￿'));
}

/**
 * Postęp projektu w procentach to udział zadań ukończonych w zadaniach ogółem; projekt bez zadań
 * ma postęp zerowy, nie stuprocentowy.
 */
export function postepProjektu(zadania: readonly WorkspaceTask[]): number {
  if (zadania.length === 0) return 0;
  const ukonczone = zadania.filter((zadanie) => zadanie.status === 'done').length;
  return Math.round((ukonczone / zadania.length) * 100);
}
