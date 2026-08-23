import {
  Command,
  type WorkspaceActivityEntry,
  type WorkspaceActivityListRequest,
  type WorkspaceAgentUnassignRequest,
  type WorkspaceAgentUnassignResponse,
  type WorkspaceBacklink,
  type WorkspaceCanvasGetRequest,
  type WorkspaceCanvasGetResponse,
  type WorkspaceCanvasSaveRequest,
  type WorkspaceCanvasSaveResponse,
  type WorkspaceComment,
  type WorkspaceCommentAddRequest,
  type WorkspaceCommentAddResponse,
  type WorkspaceCommentDeleteRequest,
  type WorkspaceCommentDeleteResponse,
  type WorkspaceCommentListRequest,
  type WorkspaceInstructionsVersion,
  type WorkspaceInstructionsVersionListRequest,
  type WorkspaceInstructionsVersionRestoreRequest,
  type WorkspaceInstructionsVersionRestoreResponse,
  type WorkspaceKnowledgeGraph,
  type WorkspaceKnowledgeGraphGetRequest,
  type WorkspaceLibraryDuplicateListRequest,
  type WorkspaceLibraryDuplicateListResponse,
  type WorkspaceLibraryDuplicateMergeRequest,
  type WorkspaceLibraryDuplicateMergeResponse,
  type WorkspaceLibraryTextExtractRequest,
  type WorkspaceLibraryTextExtractResponse,
  type WorkspaceNote,
  type WorkspaceNoteBacklinkListRequest,
  type WorkspaceNoteDeleteRequest,
  type WorkspaceNoteDeleteResponse,
  type WorkspaceNoteGetRequest,
  type WorkspaceNoteListRequest,
  type WorkspaceNoteNode,
  type WorkspaceNoteSaveRequest,
  type WorkspaceNoteSaveResponse,
  type WorkspaceNoteTreeGetRequest,
  type WorkspaceProject,
  type WorkspaceProjectStatusSetRequest,
  type WorkspaceSearchHit,
  type WorkspaceSearchProjectRequest,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { przenies } from '../../protokol/wynik-czastkowy';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Wiedza i współpraca projektu widziane przez klienta — notatki i wiki, graf
 * odnośników, tablica wizualna, materiały biblioteki projektu, wyszukiwanie,
 * oś czasu, komentarze oraz trzy czynności na samym projekcie: stan,
 * odłączenie eksperta i historia instrukcji.
 *
 * Graf tej rodziny rysuje sieć ODNOŚNIKÓW — kto na kogo wskazuje wprost.
 * Wyszukiwanie po znaczeniu prowadzi rodzina `knowledge.*` i nie ma go tutaj,
 * żeby okno nie obiecywało podobieństwa tam, gdzie dostaje dopasowanie po
 * słowach.
 */
export interface CzynnosciWiedzy {
  /** `workspace.note.save` — zapis notatki wraz z przeliczeniem odnośników. */
  zapiszNotatke(zadanie: WorkspaceNoteSaveRequest): Promise<Wynik<WorkspaceNoteSaveResponse>>;
  /** `workspace.note.get` — notatka wraz z treścią. */
  notatka(zadanie: WorkspaceNoteGetRequest): Promise<Wynik<WorkspaceNote>>;
  /** `workspace.note.list` — wykaz notatek projektu. */
  notatki(zadanie: WorkspaceNoteListRequest): Promise<Wynik<WorkspaceNote[]>>;
  /** `workspace.note.delete` — usunięcie notatki. */
  usunNotatke(zadanie: WorkspaceNoteDeleteRequest): Promise<Wynik<WorkspaceNoteDeleteResponse>>;
  /** `workspace.note.tree.get` — drzewo stron do nawigacji zakładki. */
  drzewoNotatek(zadanie: WorkspaceNoteTreeGetRequest): Promise<Wynik<WorkspaceNoteNode[]>>;
  /** `workspace.note.backlink.list` — panel „co linkuje tutaj”. */
  odnosnikiWsteczne(
    zadanie: WorkspaceNoteBacklinkListRequest,
  ): Promise<Wynik<WorkspaceBacklink[]>>;
  /** `workspace.knowledge.graph.get` — sieć powiązań bytów projektu. */
  grafWiedzy(zadanie: WorkspaceKnowledgeGraphGetRequest): Promise<Wynik<WorkspaceKnowledgeGraph>>;
  /** `workspace.canvas.get` — tablica wizualna projektu. */
  kanwa(zadanie: WorkspaceCanvasGetRequest): Promise<Wynik<WorkspaceCanvasGetResponse>>;
  /** `workspace.canvas.save` — zapis sceny tablicy wizualnej. */
  zapiszKanwe(zadanie: WorkspaceCanvasSaveRequest): Promise<Wynik<WorkspaceCanvasSaveResponse>>;
  /** `workspace.library.text.extract` — wydobycie treści pliku do wskaźnika. */
  wydobadzTekst(
    zadanie: WorkspaceLibraryTextExtractRequest,
  ): Promise<Wynik<WorkspaceLibraryTextExtractResponse>>;
  /** `workspace.library.duplicate.list` — grupy plików o treści identycznej. */
  duplikaty(
    zadanie: WorkspaceLibraryDuplicateListRequest,
  ): Promise<Wynik<WorkspaceLibraryDuplicateListResponse>>;
  /** `workspace.library.duplicate.merge` — scalenie duplikatów w jeden plik. */
  scalDuplikaty(
    zadanie: WorkspaceLibraryDuplicateMergeRequest,
  ): Promise<Wynik<WorkspaceLibraryDuplicateMergeResponse>>;
  /** `workspace.search.project` — wyszukiwanie po słowach w całym projekcie. */
  szukaj(zadanie: WorkspaceSearchProjectRequest): Promise<Wynik<WorkspaceSearchHit[]>>;
  /** `workspace.activity.list` — oś czasu aktywności projektu. */
  osCzasu(zadanie: WorkspaceActivityListRequest): Promise<Wynik<WorkspaceActivityEntry[]>>;
  /** `workspace.comment.add` — komentarz przy bycie projektu. */
  dolozKomentarz(
    zadanie: WorkspaceCommentAddRequest,
  ): Promise<Wynik<WorkspaceCommentAddResponse>>;
  /** `workspace.comment.list` — wątek komentarzy bytu albo projektu. */
  komentarze(zadanie: WorkspaceCommentListRequest): Promise<Wynik<WorkspaceComment[]>>;
  /** `workspace.comment.delete` — usunięcie komentarza wraz z odpowiedziami. */
  usunKomentarz(
    zadanie: WorkspaceCommentDeleteRequest,
  ): Promise<Wynik<WorkspaceCommentDeleteResponse>>;
  /** `workspace.project.status.set` — stan projektu (aktywny · wstrzymany · zarchiwizowany). */
  ustawStanProjektu(
    zadanie: WorkspaceProjectStatusSetRequest,
  ): Promise<Wynik<WorkspaceProject>>;
  /** `workspace.agent.unassign` — odłączenie eksperta od projektu. */
  odlaczAgenta(
    zadanie: WorkspaceAgentUnassignRequest,
  ): Promise<Wynik<WorkspaceAgentUnassignResponse>>;
  /** `workspace.instructions.version.list` — historia instrukcji systemowych. */
  wersjeInstrukcji(
    zadanie: WorkspaceInstructionsVersionListRequest,
  ): Promise<Wynik<WorkspaceInstructionsVersion[]>>;
  /** `workspace.instructions.version.restore` — przywrócenie wersji instrukcji. */
  przywrocInstrukcje(
    zadanie: WorkspaceInstructionsVersionRestoreRequest,
  ): Promise<Wynik<WorkspaceInstructionsVersionRestoreResponse>>;
}

export function czynnosciWiedzy(kanal: Kanal): CzynnosciWiedzy {
  return {
    zapiszNotatke: (zadanie) => wywolaj(kanal, Command.WorkspaceNoteSave, zadanie),

    async notatka(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceNoteGet, zadanie),
        Command.WorkspaceNoteGet,
        (tresc) => czyObiekt(tresc.note),
      );
      return przenies(wynik, (tresc) => tresc.note);
    },

    async notatki(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceNoteList, zadanie),
        Command.WorkspaceNoteList,
        (tresc) => czyTablica(tresc.notes),
      );
      return przenies(wynik, (tresc) => tresc.notes);
    },

    usunNotatke: (zadanie) => wywolaj(kanal, Command.WorkspaceNoteDelete, zadanie),

    async drzewoNotatek(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceNoteTreeGet, zadanie),
        Command.WorkspaceNoteTreeGet,
        (tresc) => czyTablica(tresc.nodes),
      );
      return przenies(wynik, (tresc) => tresc.nodes);
    },

    async odnosnikiWsteczne(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceNoteBacklinkList, zadanie),
        Command.WorkspaceNoteBacklinkList,
        (tresc) => czyTablica(tresc.backlinks),
      );
      return przenies(wynik, (tresc) => tresc.backlinks);
    },

    async grafWiedzy(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceKnowledgeGraphGet, zadanie),
        Command.WorkspaceKnowledgeGraphGet,
        (tresc) => czyObiekt(tresc.graph),
      );
      return przenies(wynik, (tresc) => tresc.graph);
    },

    kanwa: (zadanie) => wywolaj(kanal, Command.WorkspaceCanvasGet, zadanie),
    zapiszKanwe: (zadanie) => wywolaj(kanal, Command.WorkspaceCanvasSave, zadanie),
    wydobadzTekst: (zadanie) => wywolaj(kanal, Command.WorkspaceLibraryTextExtract, zadanie),
    duplikaty: (zadanie) => wywolaj(kanal, Command.WorkspaceLibraryDuplicateList, zadanie),
    scalDuplikaty: (zadanie) => wywolaj(kanal, Command.WorkspaceLibraryDuplicateMerge, zadanie),

    async szukaj(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceSearchProject, zadanie),
        Command.WorkspaceSearchProject,
        (tresc) => czyTablica(tresc.hits),
      );
      return przenies(wynik, (tresc) => tresc.hits);
    },

    async osCzasu(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceActivityList, zadanie),
        Command.WorkspaceActivityList,
        (tresc) => czyTablica(tresc.entries),
      );
      return przenies(wynik, (tresc) => tresc.entries);
    },

    dolozKomentarz: (zadanie) => wywolaj(kanal, Command.WorkspaceCommentAdd, zadanie),

    async komentarze(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceCommentList, zadanie),
        Command.WorkspaceCommentList,
        (tresc) => czyTablica(tresc.comments),
      );
      return przenies(wynik, (tresc) => tresc.comments);
    },

    usunKomentarz: (zadanie) => wywolaj(kanal, Command.WorkspaceCommentDelete, zadanie),

    async ustawStanProjektu(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceProjectStatusSet, zadanie),
        Command.WorkspaceProjectStatusSet,
        (tresc) => czyObiekt(tresc.project),
      );
      return przenies(wynik, (tresc) => tresc.project);
    },

    odlaczAgenta: (zadanie) => wywolaj(kanal, Command.WorkspaceAgentUnassign, zadanie),

    async wersjeInstrukcji(zadanie) {
      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.WorkspaceInstructionsVersionList, zadanie),
        Command.WorkspaceInstructionsVersionList,
        (tresc) => czyTablica(tresc.versions),
      );
      return przenies(wynik, (tresc) => tresc.versions);
    },

    przywrocInstrukcje: (zadanie) =>
      wywolaj(kanal, Command.WorkspaceInstructionsVersionRestore, zadanie),
  };
}

/**
 * Wcięcie węzła drzewa stron — liczba przodków strony.
 *
 * Drzewo przychodzi z rdzenia jako wykaz w kolejności obchodu, a nie jako
 * zagnieżdżenie: głębokość liczy się tu, przy rysowaniu, żeby odpowiedź nie
 * niosła tej samej informacji dwa razy.
 */
export function glebokoscWezla(
  wezly: readonly WorkspaceNoteNode[],
  wezel: WorkspaceNoteNode,
): number {
  const wedlugKodu = new Map(wezly.map((pozycja) => [pozycja.noteId, pozycja]));
  let glebokosc = 0;
  let rodzic = wezel.parentNoteId ?? '';
  while (rodzic !== '' && glebokosc <= wezly.length) {
    const nadrzedny = wedlugKodu.get(rodzic);
    if (nadrzedny === undefined) break;
    glebokosc += 1;
    rodzic = nadrzedny.parentNoteId ?? '';
  }
  return glebokosc;
}

/**
 * Nazwy odnośników `[[nazwa]]` wyjęte z treści notatki.
 *
 * Okno liczy je przed zapisem, żeby pokazać Operatorowi, do czego strona
 * linkuje, zanim rdzeń odpowie. Po zapisie prawdą jest odpowiedź rdzenia
 * (`linkedNames`, `missingNames`) — ta funkcja jest podglądem, nie drugim
 * źródłem prawdy.
 */
export function nazwyOdnosnikow(tresc: string): string[] {
  const nazwy: string[] = [];
  const widziane = new Set<string>();
  for (const trafienie of tresc.matchAll(/\[\[([^\]]+)\]\]/g)) {
    const nazwa = (trafienie[1] ?? '').trim();
    if (nazwa === '' || widziane.has(nazwa.toLowerCase())) continue;
    widziane.add(nazwa.toLowerCase());
    nazwy.push(nazwa);
  }
  return nazwy;
}
