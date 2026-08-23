import { describe, expect, it } from 'vitest';

import {
  Command,
  WorkspaceCalendarSpan,
  WorkspaceEntityKind,
  WorkspaceProjectStatus,
  WorkspaceTaskStatus,
  type WorkspaceBoard,
  type WorkspaceNoteNode,
  type WorkspaceTask,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czynnosciPlanowania, kartyKolumny, postepProjektu } from './czynnosci-planowania';
import { czynnosciWiedzy, glebokoscWezla, nazwyOdnosnikow } from './czynnosci-wiedzy';

/**
 * Sprawdziany warstwy klienckiej obszaru planowania i wiedzy modułu Workspace.
 *
 * Główny z nich pilnuje jednej rzeczy: czy KAŻDA komenda kontraktu z tego
 * obszaru ma drogę z okna do rdzenia. Wykaz komend nie jest tu przepisany —
 * powstaje z wywołań czynności, a porównywany jest ze stałymi kontraktu, więc
 * komenda dołożona do kontraktu i pominięta w oknie zostanie tu nazwana.
 *
 * Pozostałe sprawdziany dotyczą bytów rozstrzygalnych bez rdzenia: kolejności
 * kart w kolumnie tablicy, paska postępu, wcięcia drzewa stron i odnośników
 * wyjętych z treści notatki.
 */

/** Kanał próbny: zapamiętuje nazwy komend i oddaje odpowiedź pustą, lecz poprawną. */
function kanalProbny(odpowiedzi: Record<string, unknown> = {}): {
  kanal: Kanal;
  wyslane: string[];
} {
  const wyslane: string[] = [];
  const kanal = {
    wyslij(komenda: string, _zadanie: unknown, przyWyniku?: (wynik: Wynik<unknown>) => void) {
      wyslane.push(komenda);
      przyWyniku?.({ udany: true, wynik: odpowiedzi[komenda] ?? {} });
      return `zadanie-${wyslane.length}`;
    },
    naZdarzenie: () => () => undefined,
    naDowolny: () => () => undefined,
    sesja: () => ({}) as ReturnType<Kanal['sesja']>,
    dziennikNieznanych: () => ({}) as ReturnType<Kanal['dziennikNieznanych']>,
  } as unknown as Kanal;
  return { kanal, wyslane };
}

/** Zadanie przykładowe w kształcie, w którym niesie je kontrakt. */
function zadanie(id: string, dodatki: Partial<WorkspaceTask> = {}): WorkspaceTask {
  return {
    id,
    projectId: 'sprawa-2026-114',
    title: `Zadanie ${id}`,
    status: WorkspaceTaskStatus.Todo,
    createdAt: 0,
    updatedAt: 0,
    ...dodatki,
  };
}

describe('czynności planowania projektu', () => {
  it('wysyła dwanaście komend huba planowania pod nazwami z kontraktu', async () => {
    const { kanal, wyslane } = kanalProbny({
      [Command.WorkspaceTaskList]: { tasks: [], total: 0 },
      [Command.WorkspaceBoardGet]: { board: { projectId: 'p', columns: [], tasks: [] } },
    });
    const czynnosci = czynnosciPlanowania(kanal);

    await czynnosci.zalozZadanie({ projectId: 'p', title: 'Nowe' });
    await czynnosci.zmienZadanie({ taskId: 'z1', status: WorkspaceTaskStatus.Done });
    await czynnosci.usunZadanie({ taskId: 'z1' });
    await czynnosci.zadania({ projectId: 'p' });
    await czynnosci.przeniesZadanie({ taskId: 'z1' });
    await czynnosci.zalozZaleznosc({
      projectId: 'p',
      predecessorTaskId: 'z1',
      successorTaskId: 'z2',
    });
    await czynnosci.zniesZaleznosc({ dependencyId: 'd1' });
    await czynnosci.tablica({ projectId: 'p' });
    await czynnosci.harmonogram({ projectId: 'p' });
    await czynnosci.kalendarz({
      projectId: 'p',
      span: WorkspaceCalendarSpan.Month,
      anchorAt: 0,
    });
    await czynnosci.wciagnijKalendarz({ projectId: 'p', content: 'BEGIN:VCALENDAR' });
    await czynnosci.zapiszKalendarz({ projectId: 'p' });

    expect(wyslane).toEqual([
      Command.WorkspaceTaskCreate,
      Command.WorkspaceTaskUpdate,
      Command.WorkspaceTaskDelete,
      Command.WorkspaceTaskList,
      Command.WorkspaceTaskMove,
      Command.WorkspaceTaskDependencySet,
      Command.WorkspaceTaskDependencyRemove,
      Command.WorkspaceBoardGet,
      Command.WorkspaceScheduleGet,
      Command.WorkspaceCalendarGet,
      Command.WorkspaceCalendarImport,
      Command.WorkspaceCalendarExport,
    ]);
  });

  it('układa karty kolumny po kluczu porządkowym, a kartę bez klucza stawia na końcu', () => {
    const tablica: WorkspaceBoard = {
      projectId: 'p',
      columns: [],
      tasks: [
        zadanie('c', { boardColumnId: 'kolumna-todo', rank: 'n' }),
        zadanie('a', { boardColumnId: 'kolumna-todo', rank: 'd' }),
        zadanie('bez', { boardColumnId: 'kolumna-todo' }),
        zadanie('obca', { boardColumnId: 'kolumna-done', rank: 'a' }),
      ],
    };

    expect(kartyKolumny(tablica, 'kolumna-todo').map((karta) => karta.id)).toEqual([
      'a',
      'c',
      'bez',
    ]);
  });

  it('liczy postęp z zadań ukończonych, a pustemu projektowi daje zero', () => {
    expect(postepProjektu([])).toBe(0);
    expect(
      postepProjektu([
        zadanie('1', { status: WorkspaceTaskStatus.Done }),
        zadanie('2', { status: WorkspaceTaskStatus.InProgress }),
        zadanie('3', { status: WorkspaceTaskStatus.Done }),
        zadanie('4', { status: WorkspaceTaskStatus.Todo }),
      ]),
    ).toBe(50);
  });
});

describe('czynności wiedzy i współpracy projektu', () => {
  it('wysyła dwadzieścia jeden komend obszaru pod nazwami z kontraktu', async () => {
    const { kanal, wyslane } = kanalProbny({
      [Command.WorkspaceNoteGet]: { note: {} },
      [Command.WorkspaceNoteList]: { notes: [], total: 0 },
      [Command.WorkspaceNoteTreeGet]: { nodes: [] },
      [Command.WorkspaceNoteBacklinkList]: { backlinks: [], total: 0 },
      [Command.WorkspaceKnowledgeGraphGet]: { graph: { projectId: 'p', nodes: [], edges: [] } },
      [Command.WorkspaceSearchProject]: { hits: [], total: 0 },
      [Command.WorkspaceActivityList]: { entries: [], total: 0 },
      [Command.WorkspaceCommentList]: { comments: [], total: 0 },
      [Command.WorkspaceProjectStatusSet]: { project: {} },
      [Command.WorkspaceInstructionsVersionList]: { versions: [], total: 0 },
    });
    const czynnosci = czynnosciWiedzy(kanal);

    await czynnosci.zapiszNotatke({ projectId: 'p', title: 'T', content: '' });
    await czynnosci.notatka({ noteId: 'n1' });
    await czynnosci.notatki({ projectId: 'p' });
    await czynnosci.usunNotatke({ noteId: 'n1' });
    await czynnosci.drzewoNotatek({ projectId: 'p' });
    await czynnosci.odnosnikiWsteczne({ noteId: 'n1' });
    await czynnosci.grafWiedzy({ projectId: 'p' });
    await czynnosci.kanwa({ projectId: 'p' });
    await czynnosci.zapiszKanwe({ projectId: 'p', scene: {} });
    await czynnosci.wydobadzTekst({ projectId: 'p', fileId: 'plik.txt' });
    await czynnosci.duplikaty({ projectId: 'p' });
    await czynnosci.scalDuplikaty({ projectId: 'p', keepFileId: 'a', mergedFileIds: ['b'] });
    await czynnosci.szukaj({ projectId: 'p', query: 'fraza' });
    await czynnosci.osCzasu({ projectId: 'p' });
    await czynnosci.dolozKomentarz({
      projectId: 'p',
      targetKind: WorkspaceEntityKind.Task,
      targetId: 'z1',
      content: 'treść',
    });
    await czynnosci.komentarze({ projectId: 'p' });
    await czynnosci.usunKomentarz({ commentId: 'k1' });
    await czynnosci.ustawStanProjektu({
      projectId: 'p',
      status: WorkspaceProjectStatus.Paused,
    });
    await czynnosci.odlaczAgenta({ projectId: 'p', agentId: 'analityk' });
    await czynnosci.wersjeInstrukcji({ projectId: 'p' });
    await czynnosci.przywrocInstrukcje({ projectId: 'p', versionId: 'w1' });

    expect(new Set(wyslane)).toEqual(
      new Set([
        Command.WorkspaceNoteSave,
        Command.WorkspaceNoteGet,
        Command.WorkspaceNoteList,
        Command.WorkspaceNoteDelete,
        Command.WorkspaceNoteTreeGet,
        Command.WorkspaceNoteBacklinkList,
        Command.WorkspaceKnowledgeGraphGet,
        Command.WorkspaceCanvasGet,
        Command.WorkspaceCanvasSave,
        Command.WorkspaceLibraryTextExtract,
        Command.WorkspaceLibraryDuplicateList,
        Command.WorkspaceLibraryDuplicateMerge,
        Command.WorkspaceSearchProject,
        Command.WorkspaceActivityList,
        Command.WorkspaceCommentAdd,
        Command.WorkspaceCommentList,
        Command.WorkspaceCommentDelete,
        Command.WorkspaceProjectStatusSet,
        Command.WorkspaceAgentUnassign,
        Command.WorkspaceInstructionsVersionList,
        Command.WorkspaceInstructionsVersionRestore,
      ]),
    );
    expect(wyslane).toHaveLength(21);
  });

  it('liczy wcięcie węzła drzewa stron z łańcucha przodków', () => {
    const wezly: WorkspaceNoteNode[] = [
      { noteId: 'korzen', title: 'Korzeń', order: 0, childCount: 1 },
      { noteId: 'srodek', title: 'Środek', parentNoteId: 'korzen', order: 0, childCount: 1 },
      { noteId: 'lisc', title: 'Liść', parentNoteId: 'srodek', order: 0, childCount: 0 },
    ];

    expect(glebokoscWezla(wezly, wezly[0] as WorkspaceNoteNode)).toBe(0);
    expect(glebokoscWezla(wezly, wezly[2] as WorkspaceNoteNode)).toBe(2);
  });

  it('wyjmuje nazwy odnośników z treści notatki, każdą raz', () => {
    const tresc = 'Zapisano w [[Ustalenia]] oraz w [[Protokół]] i znowu w [[ustalenia]].';

    expect(nazwyOdnosnikow(tresc)).toEqual(['Ustalenia', 'Protokół']);
    expect(nazwyOdnosnikow('Notatka bez odnośników.')).toEqual([]);
  });
});
